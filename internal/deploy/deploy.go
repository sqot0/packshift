package deploy

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/sqot0/packshift/internal/config"
	"github.com/sqot0/packshift/internal/crypto"
)

type IFTPClient interface {
	Connect() error
	UploadFile(localPath string, remotePath string) error
	Disconnect() error
}

func Run(cfg *config.Config) error {
	log.Println("Running deploy")

	var ftpClient IFTPClient

	decryptedPassword, err := crypto.Decrypt(cfg.FTPConfig.Password)
	if err != nil {
		return err
	}
	cfg.FTPConfig.Password = decryptedPassword

	if cfg.FTPConfig.SSL {
		ftpClient = NewSFTPClient(&cfg.FTPConfig)
	} else {
		ftpClient = NewFTPClient(&cfg.FTPConfig)
	}

	if err := ftpClient.Connect(); err != nil {
		return err
	}
	defer ftpClient.Disconnect()

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, 8)
	errCh := make(chan error, 1)

	for localPath, remotePath := range cfg.PathMappings {
		absoluteLocalPath := filepath.Join(cwd, localPath)

		err := filepath.WalkDir(absoluteLocalPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}

			relativePath, err := filepath.Rel(absoluteLocalPath, path)
			if err != nil {
				return err
			}
			remoteFilePath := filepath.Join(remotePath, relativePath)
			remoteFilePath = strings.ReplaceAll(remoteFilePath, "\\", "/")

			wg.Add(1)
			go func(local, remote string) {
				defer wg.Done()

				sem <- struct{}{}
				defer func() { <-sem }()

				log.Printf("Uploading %s to %s\n", local, remote)
				if err := ftpClient.UploadFile(local, remote); err != nil {
					select {
					case errCh <- err:
					default:
					}
				}
			}(path, remoteFilePath)

			return nil
		})

		if err != nil {
			return err
		}
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	if err := <-errCh; err != nil {
		return err
	}

	return nil
}
