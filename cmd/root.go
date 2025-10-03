package cmd

import (
	"fmt"
	"os"

	"github.com/drieschel/dddns/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dddns",
	Short: "A lightweight dyndns client",
	Long:  `Drieschel's lightweight dyndns client`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		var clients []internal.Client
		values := viper.UnmarshalKey("domain", &clients)

		fmt.Printf("%+v\n", values)
	},
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
	viper.SetConfigName("config")       // name of config file (without extension)
	viper.SetConfigType("toml")         // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/dddns")   // path to look for the config file in
	viper.AddConfigPath("$HOME/.dddns") // call multiple times to add many search paths
	viper.AddConfigPath(".")            // optionally look for config in the working directory
	err := viper.ReadInConfig()         // Find and read the config file
	if err != nil {                     // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}
