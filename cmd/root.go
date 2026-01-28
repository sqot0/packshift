package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "packshift",
	Short: "Packshift - A tool for deploying files to FTP/SFTP servers",
	Long: `Packshift is a command-line tool that simplifies deploying files to FTP or SFTP servers.
It allows you to configure FTP connections and path mappings, then upload files from your local directories to remote servers.

Use 'packshift init' to set up your configuration, and 'packshift deploy' to upload files.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
