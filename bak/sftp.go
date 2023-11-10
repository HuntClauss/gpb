package bak

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func NewSFTPClient(user, password, host string) (*sftp.Client, error) {
	conf := &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
	}

	sshConn, err := ssh.Dial("tcp", host, conf)
	if err != nil {
		return nil, fmt.Errorf("cannot dial ssh: %w", err)
	}

	client, err := sftp.NewClient(sshConn)
	if err != nil {
		return nil, fmt.Errorf("cannot create sftp client: %w", err)
	}

	return client, err
}
