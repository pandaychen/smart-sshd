package sshserver

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sshd/conf"
	"sshd/httpauth"
	"sshd/sshauth"
	"sshd/sshbash"
	"sshd/system"

	"github.com/creack/pty" //fix bugs
	glssh "github.com/gliderlabs/ssh"
	zaplog "github.com/pandaychen/goes-wrapper/zaplog"
	"go.uber.org/zap"
)

const (
	defaultShell = "/bin/bash"
)

//sshd 配置加载
type ServerOption func(svr *SmartSshdServer) error

type SmartSshdServer struct {
	ListenAddr       string
	HostkeyFile      string
	AuthorizekeyFile string
	KnownhostsFile   string
	CaAuthorityFile  string
	CaAuthorityMap   map[string]string //multi cakey
	LocalUsers       *system.LocalUserStore

	HttpauthClient *httpauth.HttpAuthClient
	Logger         *zap.Logger

	PasswordAuthMethods            []sshauth.PasswordAuthCallback
	PublickeyAuthMethods           []sshauth.PublicKeyAuthCallback
	KeyboardinteractiveAuthMethods []sshauth.KeyboadInteractiveAuthCallback
}

func NewSmartSshdServer(config *conf.SshdConfig, options ...ServerOption) *SmartSshdServer {
	//init logger
	var (
		op []ServerOption
	)
	logger, _ := zaplog.ZapLoggerInit("sshd")

	//init local users
	users, err := system.NewLocalUserStore(config.SystemUserfilePath)
	if err != nil {
		panic(err)
	}

	svr := &SmartSshdServer{
		ListenAddr:       config.Addr,
		HostkeyFile:      config.HostkeyfilePath,
		KnownhostsFile:   config.KnownhostsFilePath,
		CaAuthorityFile:  config.CaAuthorityFile,
		AuthorizekeyFile: config.AuthorizekeyFile,
		Logger:           logger,
		LocalUsers:       users,
	}

	baseauth := &sshauth.SshAuthBase{
		Logger: logger,
	}

	passwd_auth := sshauth.NewPasswordAuth(baseauth)
	op = append(op, WithSSHAuthMethod(passwd_auth))

	pubkey_auth, err := sshauth.NewPublickeyAuth("/root/.ssh/authorized_keys", baseauth)
	if err == nil {
		op = append(op, WithSSHAuthMethod(pubkey_auth))
	} else {
		logger.Error("NewPublickeyAuth error", zap.Any("errmsg", err))
	}
	options = append(options, op...)

	for _, opt := range options {
		if err := opt(svr); err != nil {
			svr.Logger.Error("[NewSmartSshdServer]set config error", zap.String("errmsg", err.Error()))
		}
	}

	return svr
}

func (s *SmartSshdServer) StartSshdStart() error {
	var opts []glssh.Option

	//set gliderlabs/ssh auth options
	opts = append(opts,
		glssh.PasswordAuth(s.PasswordAuthCallback),
		glssh.PublicKeyAuth(s.PublicKeyCallback),
		glssh.KeyboardInteractiveAuth(s.KeyboadInteractiveAuthCallback),
	)

	//opts必须为glssh.Option
	return glssh.ListenAndServe(s.ListenAddr, s.SessionHandlerCallback, opts...)
}

// Handler is a callback for handling established SSH sessions -- type Handler func(Session)
func (s *SmartSshdServer) SessionHandlerCallback(session glssh.Session) {
	s.Logger.Info("[SessionHandlerCallback]start user session", zap.String("user", session.User()), zap.Any("raddr", session.RemoteAddr()), zap.String("sessionid", session.Context().(glssh.Context).SessionID()))

	user, err := s.LocalUsers.GetUser(session.User())
	if err != nil {
		s.Logger.Error("[SessionHandlerCallback]get user error", zap.String("errmsg", err.Error()))
		//sess.Exit(1)
		return
	}

	s.Logger.Info("[SessionHandlerCallback]", zap.String("user", session.User()), zap.Any("raddr", session.RemoteAddr()), zap.String("sessionid", session.Context().(glssh.Context).SessionID()), zap.String("shell", user.Shell))

	if user.Shell == "" {
		user.Shell = defaultShell
	}

	bash_cmd, err := sshbash.NewBashCommand(user)
	if err != nil {
		return
	}

	// bind shell and tty
	err = s.BindBashAndPty(session, bash_cmd)
	if err != nil {
		s.Logger.Info("[SessionHandlerCallback]BindSessionBashAndPty error", zap.String("errmsg", err.Error()), zap.String("user", session.User()), zap.Any("raddr", session.RemoteAddr()), zap.String("sessionid", session.Context().(glssh.Context).SessionID()), zap.String("shell", user.Shell))
		session.Exit(1)
		return
	}

	//exit
	session.Exit(0)
}

// 将bash的输入输出与tty结合，实现交互式操作环境
func (s *SmartSshdServer) BindBashAndPty(session glssh.Session, ssh_bash *sshbash.BashCommand) error {
	var (
		err error
	)
	ptyReq, winCh, isPty := session.Pty()
	if !isPty {
		return errors.New("create pty error")
	}

	//TODO：bash profile bugs
	ssh_bash.Cmd.Env = append(os.Environ(), fmt.Sprintf("TERM=%s", ptyReq.Term))
	//ssh_bash.Cmd.Env = append(ssh_bash.Cmd.Env, fmt.Sprintf("TERM=%s", ptyReq.Term))
	bashfd, err := pty.Start(ssh_bash.Cmd)
	if err != nil {
		s.Logger.Error("[BindSessionBashAndPty]pty start error", zap.String("errmsg", err.Error()))
		return err
	}

	ssh_bash.SetBashfd(bashfd)

	//monitor windows size change message
	go func() {
		for win := range winCh {
			ssh_bash.SetWindowsize(win.Width, win.Height)
		}
	}()

	//create recorder for each session

	//pipe session to bash and visa-versa
	go func() {
		io.Copy(session, bashfd)
	}()

	go func() {
		io.Copy(bashfd, session)
	}()

	err = ssh_bash.Cmd.Wait()

	return err
}
