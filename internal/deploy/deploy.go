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

	decryptedPassword, err := crypto.Decrypt(cfg.FTPConfig.Password)
	if err != nil {
		return err
	}
	cfg.FTPConfig.Password = decryptedPassword

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	type job struct {
		local  string
		remote string
	}

	jobs := make(chan job)
	errCh := make(chan error, 1)

	workerCount := 4

	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			var ftpClient IFTPClient
			if cfg.FTPConfig.SSL {
				ftpClient = NewSFTPClient(&cfg.FTPConfig)
			} else {
				ftpClient = NewFTPClient(&cfg.FTPConfig)
			}

			if err := ftpClient.Connect(); err != nil {
				select {
				case errCh <- err:
				default:
				}
				return
			}
			defer ftpClient.Disconnect()

			for j := range jobs {
				log.Printf("Uploading %s -> %s\n", j.local, j.remote)
				if err := ftpClient.UploadFile(j.local, j.remote); err != nil {
					select {
					case errCh <- err:
					default:
					}
					return
				}
			}
		}()
	}

	go func() {
		defer close(jobs)

		for localPath, remotePath := range cfg.PathMappings {
			absoluteLocalPath := filepath.Join(cwd, localPath)

			filepath.WalkDir(absoluteLocalPath, func(path string, d fs.DirEntry, err error) error {
				if err != nil || d.IsDir() {
					return err
				}

				relativePath, err := filepath.Rel(absoluteLocalPath, path)
				if err != nil {
					return err
				}

				remoteFilePath := filepath.Join(remotePath, relativePath)
				remoteFilePath = strings.ReplaceAll(remoteFilePath, "\\", "/")

				jobs <- job{
					local:  path,
					remote: remoteFilePath,
				}
				return nil
			})
		}
	}()

	go func() {
		wg.Wait()
		close(errCh)
	}()

	if err := <-errCh; err != nil {
		return err
	}

	return nil
}
