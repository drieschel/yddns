package cmd

import (
	"log"
	"net/http"
	"time"

	"github.com/drieschel/yddns/internal/cache"
	"github.com/drieschel/yddns/internal/client"
	"github.com/drieschel/yddns/internal/config"
	"github.com/spf13/cobra"
)

const (
	flagNameConfigFile      = "config-file"
	flagNameRefreshInterval = "refresh-interval"
	flagNamePeriodically    = "periodically"
)

var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh ip addresses for dynamic dns domains",
	Long:  `Refresh ip addresses for dynamic dns domains`,
	Run: func(cmd *cobra.Command, args []string) {
		refreshInterval, err := cmd.Flags().GetInt(flagNameRefreshInterval)
		if err != nil {
			log.Fatal(err)
		}

		periodically, err := cmd.Flags().GetBool(flagNamePeriodically)
		if err != nil {
			log.Fatal(err)
		}

		cfg := config.NewFileConfig(version)

		domains, err := cfg.PrepareAndGetDomains()
		if err != nil {
			log.Fatal(err)
		}

		if refreshInterval == 0 {
			refreshInterval = cfg.RefreshInterval
		}

		cacheExpirySeconds := cache.ExpirySecondsDefault
		if periodically {
			cacheExpirySeconds += refreshInterval
		}

		client := client.NewClient(cfg.CreateFileCache(cacheExpirySeconds), &http.Client{})

		for {
			client.Clear()
			for _, domain := range domains {
				response, err := client.Refresh(domain)
				if err != nil {
					log.Printf("An error occurred when refreshing %s: %s\n", domain.DomainName, err)
				} else {
					log.Printf("%s: %s", domain.DomainName, response)
				}
			}

			if !periodically {
				break
			}

			time.Sleep(time.Duration(refreshInterval) * time.Second)
		}
	},
}

func init() {
	rootCmd.AddCommand(refreshCmd)
	refreshCmd.Flags().IntP(flagNameRefreshInterval, "i", 0, "Define refresh interval in seconds")
	refreshCmd.Flags().BoolP(flagNamePeriodically, "p", false, "Refresh periodically")
	refreshCmd.Flags().StringP(flagNameConfigFile, "c", "", "Override default config using absolute file path")

	cobra.OnInitialize(initConfig)
}

func initConfig() {
	configFile, err := refreshCmd.Flags().GetString(flagNameConfigFile)
	if err != nil {
		log.Fatal(err)
	}

	config.FilePath = configFile
}
