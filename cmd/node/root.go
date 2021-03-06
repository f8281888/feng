/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package nodecmd

import (
	"feng/config"
	"feng/internal/log"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

var cfgPath string
var logPath string
var dataPath string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nodefeng",
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
		log.Assert(err.Error())
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgPath, "config-path", "../config", "config path")
	rootCmd.PersistentFlags().StringVar(&dataPath, "data-path", "../data", "data path")
	rootCmd.PersistentFlags().StringVar(&logPath, "log-path", "../logs", "log path")
	cobra.OnInitialize(initConfig)
	_ = viper.BindPFlag("log-path", rootCmd.PersistentFlags().Lookup("log-path"))
	_ = viper.BindPFlag("config-path", rootCmd.PersistentFlags().Lookup("config-path"))
	_ = viper.BindPFlag("data-path", rootCmd.PersistentFlags().Lookup("data-path"))

}

func initConfig() {
	config.InitConfig(cfgPath, "node", &config.NodeConf)
}
