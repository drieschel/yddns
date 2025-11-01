package cmd

import (
	"log"
	"net/http"
	"time"

	"github.com/drieschel/yddns/internal/client"
	"github.com/drieschel/yddns/internal/config"
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
		interval, err := cmd.Flags().GetInt(flagNameInterval)
		if err != nil {
			log.Fatal(err)
		}

		periodically, err := cmd.Flags().GetBool(flagNamePeriodically)
		if err != nil {
			log.Fatal(err)
		}

		cfg := config.New(version, fs)

		domains, err := cfg.PrepareAndGetDomains()
		if err != nil {
			log.Fatal(err)
		}

		var client = client.NewClient(&http.Client{})

		for {
			client.Clear()
			for _, domain := range domains {
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
	refreshCmd.Flags().IntP(flagNameInterval, "i", viper.GetInt(config.KeyRefreshInterval), "Define refresh interval in seconds")
	refreshCmd.Flags().BoolP(flagNamePeriodically, "p", false, "Refresh periodically")

	viper.SetDefault(config.KeyRefreshInterval, config.DefaultValueRefreshInterval)
}
