package httpauth

import "net/http"

type HttpAuthClient struct {
	Client    *http.Client
	AppCode   string
	AppSecert string
}

// TODO：借助于接口，实现sshd认证的外部认证方式
func NewHttpAuthClient() *HttpAuthClient {
	client := &HttpAuthClient{}

	return client
}
