package sshbash

import (
	"os"
	"os/exec"
	sys "sshd/system"
	"syscall"
	"unsafe"
)

// 带pty的bash结构封装
type BashCommand struct {
	Shell  string
	Cmd    *exec.Cmd
	Bashfd *os.File //bash的fd
	User   *sys.User
}

func NewBashCommand(user *sys.User) (*BashCommand, error) {
	bashcmd := &BashCommand{
		Shell: user.Shell,
		User:  user,
	}

	//init shell environment
	bashcmd.Cmd = exec.Command(user.Shell)
	bashcmd.Cmd.Args = []string{
		user.Shell,
	}

	bashcmd.Cmd.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{
			Uid:         user.UID,
			Gid:         user.GID,
			NoSetGroups: true,
		},
	}

	return bashcmd, nil
}

func (b *BashCommand) SetBashfd(fd *os.File) {
	b.Bashfd = fd
}

func (b *BashCommand) SetWindowsize(width, height int) {
	syscall.Syscall(syscall.SYS_IOCTL, b.Bashfd.Fd(), uintptr(syscall.TIOCSWINSZ),
		uintptr(unsafe.Pointer(&struct{ h, w, x, y uint16 }{uint16(height), uint16(width), 0, 0})))
}
