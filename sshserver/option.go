package sshserver

import (
	"errors"
	"sshd/sshauth"
)

func WithAddress(addr string) ServerOption {
	return func(s *SmartSshdServer) error {
		s.ListenAddr = addr
		return nil
	}
}

func WithCAauthority(ca_cert_path string) ServerOption {
	return func(s *SmartSshdServer) error {
		s.CaAuthorityFile = ca_cert_path
		return nil
	}
}

func WithHostKeyFile(hostkeyfile string) ServerOption {
	return func(s *SmartSshdServer) error {
		s.HostkeyFile = hostkeyfile
		return nil
	}
}

//设置认证方法
func WithSSHAuthMethod(ssh_auth_method interface{}) ServerOption {
	return func(s *SmartSshdServer) error {
		switch ssh_auth_method.(type) {
		case *sshauth.PublickeyAuth:
			md := ssh_auth_method.(*sshauth.PublickeyAuth)
			s.PublickeyAuthMethods = append(s.PublickeyAuthMethods, md)
		case *sshauth.PasswordAuth:
			md := ssh_auth_method.(*sshauth.PasswordAuth)
			s.PasswordAuthMethods = append(s.PasswordAuthMethods, md)
		default:
			return errors.New("unknown ssh auth method")
		}
		return nil
	}
}
