package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/sqot0/packshift/internal/config"
	"github.com/sqot0/packshift/internal/crypto"
	"github.com/sqot0/packshift/internal/prompt"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Packshift configuration",
	Long: `Initialize Packshift by interactively setting up FTP or SFTP configuration and defining path mappings.
This command prompts for server details, credentials, and mappings between local directories and remote paths on the server.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("   ___           _     __ _     _  __ _   \n" +
			"  / _ \\__ _  ___| | __/ _\\ |__ (_)/ _| |_ \n" +
			" / /_)/ _` |/ __| |/ /\\ \\| '_ \\| | |_| __|\n" +
			"/ ___/ (_| | (__|   < _\\ \\ | | | |  _| |_ \n" +
			"\\/    \\__,_|\\___|_|\\_\\\\__/_| |_|_|_|  \\__|")

		// Prompt for FTP config
		host, err := prompt.Text("FTP Host", "")
		if err != nil {
			log.Fatal(err)
		}

		portStr, err := prompt.Text("FTP Port", "21")
		if err != nil {
			log.Fatal(err)
		}
		port, err := strconv.Atoi(portStr)
		if err != nil {
			log.Fatal(err)
		}

		username, err := prompt.Text("FTP Username", "")
		if err != nil {
			log.Fatal(err)
		}

		password, err := prompt.Password("FTP Password")
		if err != nil {
			log.Fatal(err)
		}

		protocol, err := prompt.Select("Choose protocol", []string{"FTP", "SFTP"})
		if err != nil {
			log.Fatal(err)
		}

		ssl := protocol == 1

		encrypted, err := crypto.Encrypt(password)
		if err != nil {
			log.Fatal(err)
		}

		ftpConfig := config.FTPConfig{
			Host:     host,
			Port:     port,
			Username: username,
			Password: encrypted,
			SSL:      ssl,
		}

		// Prompt for path mappings
		pathMappings := make(map[string]string)
		for {
			local, err := prompt.Text("Local path (relative to current directory, leave empty to finish)", "")
			if err != nil {
				log.Fatal(err)
			}
			if local == "" {
				break
			}

			remote, err := prompt.Text("Remote path on FTP server", "")
			if err != nil {
				log.Fatal(err)
			}

			pathMappings[local] = remote

			addMore, err := prompt.Confirm("Add another path mapping?")
			if err != nil {
				log.Fatal(err)
			}
			if !addMore {
				break
			}
		}

		// Get current working directory
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		// Initialize config
		err = config.Init(cwd, &ftpConfig, pathMappings)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Packshift initialized successfully!")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
