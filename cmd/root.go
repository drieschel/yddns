package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	FlagNameConfigFile = "config-file"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "yddns",
	Short: "a flexible and lightweight dyndns client",
	Long:  `drieschel's flexible and lightweight dyndns client`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP(FlagNameConfigFile, "c", "", "Override default config using absolute file path")

	cobra.OnInitialize(initConfig)
}

func initConfig() {
	configFile, err := rootCmd.PersistentFlags().GetString(FlagNameConfigFile)
	if err != nil {
		log.Fatal(err)
	}

	viper.SetConfigFile(configFile)

	if configFile == "" {
		viper.SupportedExts = []string{"toml", "json", "yaml", "yml"}
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
