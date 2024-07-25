package config

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
)

type FTPConfig struct {
	Host          string
	Port          string
	Username      string
	Password      string
	BufferSize    string
	DefaultFolder string
}

type Config struct {
	FTPConfig
}

func (c *Config) SetupConfig() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	//config FTP SERVER
	c.FTPConfig = FTPConfig{
		Host:          os.Getenv("SFTP_HOST"),
		Port:          os.Getenv("SFTP_PORT"),
		Username:      os.Getenv("SFTP_USERNAME"),
		Password:      os.Getenv("SFTP_PASSWORD"),
		BufferSize:    os.Getenv("BUFFER_SIZE"),
		DefaultFolder: os.Getenv("DEFAULT_DIRECTORY"),
	}

	if c.FTPConfig.Port == "" || c.FTPConfig.Host == "" || c.FTPConfig.Username == "" || c.FTPConfig.Password == "" || c.FTPConfig.BufferSize == "" || c.FTPConfig.DefaultFolder == "" {
		return errors.New("env variable is not set")
	}
	return nil
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := cfg.SetupConfig(); err != nil {
		return nil, err
	}

	return cfg, nil
}
