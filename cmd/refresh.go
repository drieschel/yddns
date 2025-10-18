package cmd

import (
	"log"
	"net/http"
	"time"

	"github.com/drieschel/yddns/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const DEFAULT_VALUE_REFRESH_INTERVAL = 600
const CONFIG_KEY_DOMAINS = "domain"
const CONFIG_KEY_REFRESH_INTERVAL = "refresh_interval"
const FLAG_NAME_PERIODICALLY = "periodically"
const FLAG_NAME_INTERVAL = "interval"

var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh ip addresses for dynamic dns domains",
	Long:  `Refresh ip addresses for dynamic dns domains`,
	Run: func(cmd *cobra.Command, args []string) {
		domains := internal.Domains{}
		err := viper.UnmarshalKey(CONFIG_KEY_DOMAINS, &domains.List)
		if err != nil {
			log.Fatal(err)
		}

		interval, err := cmd.Flags().GetInt(FLAG_NAME_INTERVAL)
		if err != nil {
			log.Fatal(err)
		}

		periodically, err := cmd.Flags().GetBool(FLAG_NAME_PERIODICALLY)
		if err != nil {
			log.Fatal(err)
		}

		var client = internal.NewClient(&http.Client{})

		for {
			for _, domain := range domains.List {
				err = client.Refresh(domain)
				if err != nil {
					log.Printf("an error occured when refreshing %s: %s\n", domain.Name, err)
				} else {
					log.Printf("%s successfully refreshed", domain.Name)
				}
			}

			if !periodically {
				break
			}

			time.Sleep(time.Duration(interval) * time.Second)
		}
	},
}

func init() {
	rootCmd.AddCommand(refreshCmd)

	viper.SetConfigName("config")       // name of config file (without extension)
	viper.SetConfigType("toml")         // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/yddns")   // path to look for the config file in
	viper.AddConfigPath("$HOME/.yddns") // call multiple times to add many search paths
	viper.AddConfigPath("./")           // optionally look for config in the working directory
	viper.SetDefault(CONFIG_KEY_REFRESH_INTERVAL, DEFAULT_VALUE_REFRESH_INTERVAL)
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Fatal(err)
	}

	refreshCmd.Flags().IntP(FLAG_NAME_INTERVAL, "i", viper.GetInt(CONFIG_KEY_REFRESH_INTERVAL), "refresh interval in seconds")
	refreshCmd.Flags().BoolP(FLAG_NAME_PERIODICALLY, "p", false, "refresh periodically")
}
