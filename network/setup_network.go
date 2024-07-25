package network

import (
	"Tools/config"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func SetupNetwork(cfg *config.Config) (*ssh.Client, *sftp.Client, error) {

	url := fmt.Sprintf("%s:%s", cfg.FTPConfig.Host, cfg.FTPConfig.Port)

	sshClient, err := ConnectToSshClient(url, cfg.FTPConfig.Username, cfg.FTPConfig.Password)
	if err != nil {
		return nil, nil, err
	}

	sftpClient, err := ConnectToSftp(sshClient, url, cfg.FTPConfig.Username, cfg.FTPConfig.Password)
	if err != nil {
		return nil, nil, err
	}

	return sshClient, sftpClient, nil
}
