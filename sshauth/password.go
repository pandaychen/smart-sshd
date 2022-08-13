package sshauth

import (
	"github.com/gliderlabs/ssh"
)

type PasswordAuth struct {
	*SshAuthBase
}

func NewPasswordAuth(base *SshAuthBase) *PasswordAuth {
	return &PasswordAuth{
		SshAuthBase: base,
	}
}

/*
sshserver/option.go:44:34: cannot use ssh_auth_method.(sshauth.PasswordAuth) (type sshauth.PasswordAuth) as type sshauth.PasswordAuthCallback in append:
        sshauth.PasswordAuth does not implement sshauth.PasswordAuthCallback (Auth method has pointer receiver)
*/
//func (p *PasswordAuth) SSHAuth(ctx ssh.Context, password string) bool
func (p PasswordAuth) Auth(ctx ssh.Context, password string) bool {
	//verify password with apicalls

	return true
}
