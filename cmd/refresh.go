package cmd

import (
	"log"
	"net/http"
	"time"

	"github.com/drieschel/yddns/internal/client"
	"github.com/drieschel/yddns/internal/config"
	"github.com/spf13/cobra"
)

const (
	flagConfigFile      = "config-file"
	flagRefreshInterval = "refresh-interval"
	flagPeriodically    = "periodically"
)

var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh ip addresses for dynamic dns domains",
	Long:  `Refresh ip addresses for dynamic dns domains`,
	Run: func(cmd *cobra.Command, args []string) {
		refreshInterval, err := cmd.Flags().GetInt(flagRefreshInterval)
		if err != nil {
			log.Fatal(err)
		}

		periodically, err := cmd.Flags().GetBool(flagPeriodically)
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

		cacheLifetime := cfg.CacheModifiedExpirySeconds
		if periodically && cacheLifetime > 0 && cacheLifetime < (refreshInterval+2*len(domains)) {
			cfg.CacheModifiedExpirySeconds = refreshInterval + 2*len(domains)
		}

		client := client.NewClient(cfg.CreateFileCache(), &http.Client{})

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
	refreshCmd.Flags().StringP(flagConfigFile, "c", "", "Override default config using absolute file path")
	refreshCmd.Flags().BoolP(flagPeriodically, "p", false, "Refresh periodically")
	refreshCmd.Flags().IntP(flagRefreshInterval, "i", 0, "Define refresh interval in seconds")

	cobra.OnInitialize(initConfig)
}

func initConfig() {
	configFile, err := refreshCmd.Flags().GetString(flagConfigFile)
	if err != nil {
		log.Fatal(err)
	}

	config.FilePath = configFile
}
