package handler

import (
	"fmt"
	"net"
	"log"
	"time"

	"github.com/chinx/rsa/model"
	"golang.org/x/crypto/ssh"
)

func ConnectSSH(remote *model.Remote) error {
	addr := fmt.Sprintf("%s:%d", remote.Host, remote.Port)

	conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
	if err != nil {
		log.Println(0)
		return err
	}

	defer conn.Close()

	c, chans, reqs, err := ssh.NewClientConn(conn, addr, &ssh.ClientConfig{
		User: remote.User,
		Auth: []ssh.AuthMethod{ssh.Password(remote.Passwd)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	})
	if err != nil {
		log.Println(1)
		return err
	}

	client := ssh.NewClient(c, chans, reqs)
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		log.Println(2)
		return err
	}
	session.Close()
	return nil
}
