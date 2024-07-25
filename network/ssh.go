package network

import (
	"golang.org/x/crypto/ssh"
	"log"
)

func ConnectToSshClient(url, username, password string) (*ssh.Client, error) {
	sshCfg := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// connect to ssh
	conn, err := ssh.Dial("tcp", url, sshCfg)
	if err != nil {
		log.Fatal(err)
	}

	return conn, nil
}
