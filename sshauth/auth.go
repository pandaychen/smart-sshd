package sshauth

import (
	"sshd/httpauth"

	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"

	//"github.com/pkg/errors"
	"go.uber.org/zap"
)

type SshAuthBase struct {
	HttpauthClient *httpauth.HttpAuthClient
	Logger         *zap.Logger
}

type PasswordAuthCallback interface {
	// PasswordHandler is a callback for performing password authentication.
	Auth(ctx ssh.Context, password string) bool
}

type PublicKeyAuthCallback interface {
	// PublicKeyHandler is a callback for performing public key authentication.
	Auth(ctx ssh.Context, key ssh.PublicKey) bool
}

type KeyboadInteractiveAuthCallback interface {
	// KeyboardInteractiveHandler is a callback for performing keyboard-interactive authentication.
	Auth(ctx ssh.Context, challenger gossh.KeyboardInteractiveChallenge) bool
}
