package deploy

import (
	"fmt"
	"os"

	"github.com/jlaffaye/ftp"
	"github.com/sqot0/packshift/internal/config"
)

type FTPClient struct {
	Host, User, Password string
	Port                 int
	conn                 *ftp.ServerConn
}

func (c *FTPClient) Connect() error {
	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	conn, err := ftp.Dial(addr)
	if err != nil {
		return err
	}
	err = conn.Login(c.User, c.Password)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *FTPClient) UploadFile(localPath string, remotePath string) error {
	file, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer file.Close()
	return c.conn.Stor(remotePath, file)
}

func (c *FTPClient) Disconnect() error {
	if c.conn != nil {
		err := c.conn.Quit()
		c.conn = nil
		return err
	}
	return nil
}

func NewFTPClient(ftpConfig *config.FTPConfig) *FTPClient {
	return &FTPClient{
		Host:     ftpConfig.Host,
		Port:     ftpConfig.Port,
		User:     ftpConfig.Username,
		Password: ftpConfig.Password,
	}
}
