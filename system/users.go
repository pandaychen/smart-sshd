package system

import (
	"bufio"
	"errors"
	"os"
	"sshd/util"
	"strings"
	"sync"
)

type User struct {
	Name     string
	UID      uint32
	GID      uint32
	Homepath string
	Shell    string
}

type LocalUserStore struct {
	Lock     *sync.RWMutex
	Path     string
	Usersmap map[string]*User
}

func NewLocalUserStore(path string) (*LocalUserStore, error) {
	usermap := &LocalUserStore{
		Path:     path,
		Lock:     new(sync.RWMutex),
		Usersmap: make(map[string]*User),
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines := strings.Split(scanner.Text(), ":")
		if len(lines) < 7 {
			continue
		}
		uid, err := util.ParseUint32(lines[2])
		if err != nil {
			continue
		}
		gid, err := util.ParseUint32(lines[3])
		if err != nil {
			continue
		}
		user := &User{
			Name:     lines[0],
			UID:      uid,
			GID:      gid,
			Homepath: lines[5],
			Shell:    lines[6],
		}
		usermap.Usersmap[lines[0]] = user
	}

	return usermap, nil
}

func (s *LocalUserStore) GetUser(name string) (*User, error) {
	s.Lock.RLock()
	defer s.Lock.RUnlock()
	_, exists := s.Usersmap[name]
	if exists {
		return s.Usersmap[name], nil
	} else {
		return nil, errors.New("find user error")
	}
}
