package sshauth

import (
	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
)

//KeyboardAuth is instance of PublicKeyAuthCallback
type KeyboardAuth struct {
	*SshAuthBase
}

func NewKeyboardAuth(base *SshAuthBase) *KeyboardAuth {
	return &KeyboardAuth{
		SshAuthBase: base,
	}
}

func (p KeyboardAuth) Auth(ctx ssh.Context, challenger gossh.KeyboardInteractiveChallenge) bool {

	return true
}
