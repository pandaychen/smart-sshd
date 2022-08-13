package sshauth

import (
	"bufio"
	"os"
	"sync"

	"github.com/gliderlabs/ssh"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

//only support fixed path key file auth
type PublickeyAuth struct {
	*SshAuthBase
	sync.RWMutex
	AuthorizekeyFilePath string
	PubkeyMap            map[string]ssh.PublicKey
}

//
func NewPublickeyAuth(auth_file string, base *SshAuthBase) (*PublickeyAuth, error) {
	p := &PublickeyAuth{
		SshAuthBase:          base,
		AuthorizekeyFilePath: auth_file,
		PubkeyMap:            make(map[string]ssh.PublicKey),
	}

	//INIT reading  file at once
	f, err := os.Open(p.AuthorizekeyFilePath)
	if err != nil {
		p.Logger.Error("[NewPublickeyAuth]open AuthorizekeyFilePath error", zap.String("errmsg", err.Error()), zap.String("path", p.AuthorizekeyFilePath))
		return nil, errors.Wrap(err, "open authorizedkey file error")
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		pk, _, _, _, err := ssh.ParseAuthorizedKey(scanner.Bytes())
		if err != nil {
			p.Logger.Error("[NewPublickeyAuth]parse AuthorizekeyFilePath line error", zap.String("errmsg", err.Error()), zap.String("path", p.AuthorizekeyFilePath))
			continue
		}
		p.PubkeyMap[string(pk.Marshal())] = pk
	}

	//TODO:Auto update pkeys if authorizekeyfile updated

	return p, nil
}

/*
实现pubkey的认证方法：https://github.com/gliderlabs/ssh/blob/master/ssh.go#L39
// PublicKeyHandler is a callback for performing public key authentication.
type PublicKeyHandler func(ctx Context, key PublicKey) bool
*/
func (p PublickeyAuth) Auth(ctx ssh.Context, key ssh.PublicKey) bool {
	return true
}
