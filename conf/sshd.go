package conf

import "errors"

type SshdConfig struct {
	Addr               string
	SessionStorePath   string
	KnownhostsFilePath string
	HostkeyfilePath    string
	AuthorizekeyFile   string
	CaAuthorityFile    string
	IsDebug            bool
	SystemUserfilePath string
}

//global
var GSshdConfig SshdConfig

func SSHDConfigInit() {
	Config := vipers.Use("server")
	if Config == nil {
		panic(errors.New("find sshd config error"))
		return
	}
	Subconfig := Config.Use("sshd")
	if Subconfig == nil {
		panic(errors.New("find sshd config error"))
		return
	}

	GSshdConfig.Addr = Subconfig.GetString("listen")
	GSshdConfig.SystemUserfilePath = Subconfig.MustString("userfilepath", "/etc/passwd")
}
