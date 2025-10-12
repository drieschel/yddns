package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/drieschel/yddns/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh ip addresses for dynamic dns domains",
	Long:  `Refresh ip addresses for dynamic dns domains`,
	Run: func(cmd *cobra.Command, args []string) {
		var domains []internal.Domain
		err := viper.UnmarshalKey("domain", &domains)
		if err != nil {
			panic(err)
		}

		var client = internal.NewClient(domains, &http.Client{})

		for _, domain := range client.Domains {
			err = client.Refresh(domain)
			if err != nil {
				fmt.Printf("an error occured when refreshing %s: %s\n", domain.Domain, err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(refreshCmd)

	viper.SetConfigName("config")       // name of config file (without extension)
	viper.SetConfigType("toml")         // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/yddns")   // path to look for the config file in
	viper.AddConfigPath("$HOME/.yddns") // call multiple times to add many search paths
	viper.AddConfigPath(".")            // optionally look for config in the working directory
	err := viper.ReadInConfig()         // Find and read the config file
	if err != nil {                     // Handle errors reading the config file
		log.Fatalf("an error occurred when reading config: %w", err)
	}
}
