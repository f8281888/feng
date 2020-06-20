package nodecmd

import (
	"feng/business/node"
	"feng/internal/log"
	"fmt"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.AppName = "node-start"
		log.AppLog().Infof("node start")
		fmt.Println("node start")
		node.Start()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
