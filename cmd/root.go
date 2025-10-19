package cmd

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	FlagNameConfigFile      = "config-file"
	FlagNameConfigExtension = "config-ext"
)

var SupportedConfigExtensions = []string{"json", "toml", "yaml", "yml"}

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
	configExtensionsString := fmt.Sprintf("\"%s\"", strings.Join(SupportedConfigExtensions, "\", \""))
	rootCmd.Flags().StringP(FlagNameConfigFile, "c", "", "Override default config using absolute file path")
	rootCmd.Flags().StringP(FlagNameConfigExtension, "e", "toml", fmt.Sprintf("Change default config extension with a supported one: %s", configExtensionsString))

	configFile, err := rootCmd.Flags().GetString(FlagNameConfigFile)
	if err != nil {
		log.Fatal(err)
	}

	viper.SetConfigFile(configFile)
	if configFile == "" {
		var configExtension string
		configExtension, err = rootCmd.Flags().GetString(FlagNameConfigExtension)
		if err != nil {
			log.Fatal(err)
		}

		if !slices.Contains(SupportedConfigExtensions, configExtension) {
			log.Fatalf("Config extension \"%s\" not supported. Valid extensions are %s", configExtension, configExtensionsString)
		}

		viper.SetConfigName("config")
		viper.SetConfigType("toml")
		viper.AddConfigPath("/etc/yddns")
		viper.AddConfigPath("$HOME/.yddns")
		viper.AddConfigPath("./")
	}

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
}
