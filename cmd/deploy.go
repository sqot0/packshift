package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/sqot0/packshift/internal/config"
	"github.com/sqot0/packshift/internal/deploy"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy files to the remote server",
	Long: `Deploy all configured files from local paths to the remote server using the specified FTP or SFTP connection.
Files are uploaded concurrently for efficiency.`,
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		cfg, err := config.Load(cwd)
		if err != nil {
			log.Fatal(err)
		}

		if err := deploy.Run(cfg); err != nil {
			log.Fatal(err)
		}

		log.Println("All files deployed to server")
	}}

func init() {
	rootCmd.AddCommand(deployCmd)
}
