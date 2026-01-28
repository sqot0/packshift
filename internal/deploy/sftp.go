package deploy

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/sftp"
	"github.com/sqot0/packshift/internal/config"
	"golang.org/x/crypto/ssh"
)

type SFTPClient struct {
	Host, User, Password string
	Port                 int
	sshClient            *ssh.Client
	sftpClient           *sftp.Client
}

func (c *SFTPClient) Connect() error {
	config := &ssh.ClientConfig{
		User: c.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(c.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	sshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return err
	}
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		sshClient.Close()
		return err
	}
	c.sshClient = sshClient
	c.sftpClient = sftpClient
	return nil
}

func (c *SFTPClient) UploadFile(localPath string, remotePath string) error {
	localFile, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer localFile.Close()
	dir := filepath.Dir(remotePath)
	if err := c.sftpClient.MkdirAll(dir); err != nil {
		return err
	}
	remoteFile, err := c.sftpClient.Create(remotePath)
	if err != nil {
		return err
	}
	defer remoteFile.Close()
	_, err = io.Copy(remoteFile, localFile)
	return err
}

func (c *SFTPClient) Disconnect() error {
	if c.sftpClient != nil {
		c.sftpClient.Close()
		c.sftpClient = nil
	}
	if c.sshClient != nil {
		err := c.sshClient.Close()
		c.sshClient = nil
		return err
	}
	return nil
}

func NewSFTPClient(ftpConfig *config.FTPConfig) *SFTPClient {
	return &SFTPClient{
		Host:     ftpConfig.Host,
		Port:     ftpConfig.Port,
		User:     ftpConfig.Username,
		Password: ftpConfig.Password,
	}
}
