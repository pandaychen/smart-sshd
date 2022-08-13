package main

import (
	"sshd/conf"
	"sshd/sshserver"
)

func InitConfig() {
	conf.InitConfigAbpath("conf", "sshd", "yaml")
	conf.SSHDConfigInit()
}

func startSshdServer() {
	s := sshserver.NewSmartSshdServer(&conf.GSshdConfig)
	s.StartSshdStart()
}

func main() {
	InitConfig()
	startSshdServer()
}
