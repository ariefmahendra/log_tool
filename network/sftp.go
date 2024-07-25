package network

import (
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func ConnectToSftp(conn *ssh.Client, url string, username string, password string) (*sftp.Client, error) {
	client, err := sftp.NewClient(conn)
	if err != nil {
		return nil, err
	}

	return client, nil
}
