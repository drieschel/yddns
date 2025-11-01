package cmd

import (
	"log"
	"os"

	"github.com/drieschel/yddns/internal/config"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

const (
	flagNameConfigFile = "config-file"
)

var (
	fs      = afero.NewOsFs()
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

	config.FilePath = configFile
}
