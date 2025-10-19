package cmd

import (
	"log"
	"net/http"
	"time"

	"github.com/drieschel/yddns/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	DefaultValueRefreshInterval = 600
	ConfigKeyDomains            = "domain"
	ConfigKeyRefreshInterval    = "refresh_interval"
	FlagNamePeriodically        = "periodically"
	FlagNameInterval            = "interval"
)

var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh ip addresses for dynamic dns domains",
	Long:  `Refresh ip addresses for dynamic dns domains`,
	Run: func(cmd *cobra.Command, args []string) {
		domains := internal.Domains{}
		err := viper.UnmarshalKey(ConfigKeyDomains, &domains.List)
		if err != nil {
			log.Fatal(err)
		}

		interval, err := cmd.Flags().GetInt(FlagNameInterval)
		if err != nil {
			log.Fatal(err)
		}

		periodically, err := cmd.Flags().GetBool(FlagNamePeriodically)
		if err != nil {
			log.Fatal(err)
		}

		var client = internal.NewClient(&http.Client{})

		for {
			client.Clear()
			for _, domain := range domains.List {
				err = client.Refresh(domain)
				if err != nil {
					log.Printf("An error occured when refreshing %s: %s\n", domain.Name, err)
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
	refreshCmd.Flags().IntP(FlagNameInterval, "i", viper.GetInt(ConfigKeyRefreshInterval), "Define refresh interval in seconds")
	refreshCmd.Flags().BoolP(FlagNamePeriodically, "p", false, "Execute refresh periodically")

	viper.SetDefault(ConfigKeyRefreshInterval, DefaultValueRefreshInterval)
}
