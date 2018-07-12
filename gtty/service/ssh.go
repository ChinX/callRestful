package service

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

type Server struct {
	Name     string `json:"name"`
	Host     string `json:"ip"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}

func (s *Server) Connection() {
	config := &ssh.ClientConfig{
		User: s.User,
		Auth: []ssh.AuthMethod{ssh.Password(s.Password)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	addr := s.Host + ":" + strconv.Itoa(s.Port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		fmt.Println("connect tcp error: ", err)
		return
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		fmt.Println("create ssh session error: ", err)
		return
	}

	defer session.Close()

	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		fmt.Println("create terminal raw error: ", err)
		return
	}

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	defer terminal.Restore(fd, oldState)

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := session.RequestPty("xterm-256color", 200, 40, modes); err != nil {
		fmt.Println("create pty error: ", err)
		return
	}

	err = session.Shell()
	if err != nil {
		fmt.Println("exec shell command error: ", err)
		return
	}

	err = session.Wait()
	if err != nil {
		//fmt.Println("执行Wait出错:", err)
		return
	}
}
