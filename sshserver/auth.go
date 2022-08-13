package sshserver

import (
	"sshd/pkg/hash"

	glssh "github.com/gliderlabs/ssh"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
)

//password auth
func (s *SmartSshdServer) PasswordAuthCallback(ctx glssh.Context, password string) bool {
	for _, a := range s.PasswordAuthMethods {
		if a.Auth(ctx, password) {
			s.Logger.Info("[PasswordAuthCallback]auth succ", zap.String("password", hash.GetMd5Str(password)))
			return true
		}
	}

	//auth error
	return false
}

//public key Or cacert auth
func (s *SmartSshdServer) PublicKeyCallback(ctx glssh.Context, key glssh.PublicKey) bool {
	for _, a := range s.PublickeyAuthMethods {
		if a.Auth(ctx, key) {
			s.Logger.Info("[PublicKeyCallback]auth succ")
			return true
		}
	}
	return false
}

// kb auth
func (s *SmartSshdServer) KeyboadInteractiveAuthCallback(ctx glssh.Context, challenger ssh.KeyboardInteractiveChallenge) bool {
	for _, a := range s.KeyboardinteractiveAuthMethods {
		if a.Auth(ctx, challenger) {
			s.Logger.Info("[KeyboadInteractiveAuthCallback]auth succ")
			return true
		}
	}
	return false
}
