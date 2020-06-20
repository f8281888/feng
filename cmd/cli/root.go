package clicmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "feng",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var cfgPath string
var logPath string

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgPath, "config-path", "../config", "config path")
	rootCmd.PersistentFlags().StringVar(&logPath, "log-path", "../logs", "log path")
	cobra.OnInitialize(initConfig)
	_ = viper.BindPFlag("log-path", rootCmd.PersistentFlags().Lookup("log-path"))
	_ = viper.BindPFlag("config-path", rootCmd.PersistentFlags().Lookup("config-path"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
}
