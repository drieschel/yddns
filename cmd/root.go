package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagNameConfigFile = "config-file"
)

var (
	version = "dev"
	rootCmd = &cobra.Command{
		Version: version,
		Use:     "yddns",
		Short:   "A flexible and lightweight dyndns client",
		Long:    `Drieschel's flexible and lightweight dyndns client`,
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP(flagNameConfigFile, "c", "", "Override default config using absolute file path")

	cobra.OnInitialize(initConfig)
}

func initConfig() {
	configFile, err := rootCmd.PersistentFlags().GetString(flagNameConfigFile)
	if err != nil {
		log.Fatal(err)
	}

	viper.SupportedExts = []string{"toml", "json", "yaml", "yml"}
	viper.SetConfigFile(configFile)

	if configFile == "" {
		viper.SetConfigName("config")
		viper.AddConfigPath("/etc/yddns")
		viper.AddConfigPath("$HOME/.yddns")
		viper.AddConfigPath("./")
	}

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
}
