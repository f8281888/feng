package clicmd

import (
	"feng/business/cli"
	"feng/internal/log"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.AppName = "cli-get"
		cli.Start()
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
