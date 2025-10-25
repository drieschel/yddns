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
	flagNameInterval     = "interval"
	flagNamePeriodically = "periodically"
)

var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh ip addresses for dynamic dns domains",
	Long:  `Refresh ip addresses for dynamic dns domains`,
	Run: func(cmd *cobra.Command, args []string) {
		domains := internal.Domains{}
		err := viper.UnmarshalKey(internal.ConfigKeyDomains, &domains.List, viper.DecodeHook(internal.CreateDomainConfigDefaultsHookFunc()))
		if err != nil {
			log.Fatal(err)
		}

		interval, err := cmd.Flags().GetInt(flagNameInterval)
		if err != nil {
			log.Fatal(err)
		}

		periodically, err := cmd.Flags().GetBool(flagNamePeriodically)
		if err != nil {
			log.Fatal(err)
		}

		var client = internal.NewClient(&http.Client{})

		for {
			client.Clear()
			for _, domain := range domains.List {
				err = client.Refresh(&domain)
				if err != nil {
					log.Printf("An error occured when refreshing %s: %s\n", domain.DomainName, err)
				} else {
					log.Printf("%s successfully refreshed", domain.DomainName)
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
	refreshCmd.Flags().IntP(flagNameInterval, "i", viper.GetInt(internal.ConfigKeyRefreshInterval), "Define refresh interval in seconds")
	refreshCmd.Flags().BoolP(flagNamePeriodically, "p", false, "Execute refresh periodically")

	viper.SetDefault(internal.ConfigKeyRefreshInterval, internal.DefaultValueRefreshInterval)
}
