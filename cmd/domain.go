package cmd

import (
	"log"
	"net/http"
	"strings"

	"github.com/drieschel/yddns/internal/cache"
	"github.com/drieschel/yddns/internal/client"
	"github.com/drieschel/yddns/internal/config"
	"github.com/spf13/cobra"
)

var (
	flagAuthMethod    = createFlagName(config.KeyAuthMethod)
	flagDomainName    = createFlagName(config.KeyDomainName)
	flagHost          = createFlagName(config.KeyHost)
	flagIp4Address    = createFlagName(config.KeyIp4Address)
	flagIp6Address    = createFlagName(config.KeyIp6Address)
	flagIp6HostId     = createFlagName(config.KeyIp6HostId)
	flagPassword      = createFlagName(config.KeyPassword)
	flagProtocol      = createFlagName(config.KeyProtocol)
	flagRequestMethod = createFlagName(config.KeyRequestMethod)
	flagUserAgent     = createFlagName(config.KeyUserAgent)
	flagUsername      = createFlagName(config.KeyUsername)

	flagCacheCreatedLifetime  = createFlagName(config.KeyCacheCreatedExpirySeconds)
	flagCacheModifiedLifetime = createFlagName(config.KeyCacheModifiedExpirySeconds)
)

var domainCmd = &cobra.Command{
	Use:   "domain [refresh-url | :template-name]",
	Short: "Refresh a domain config via cli",
	Long:  `Refresh a single domain config via command"`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := createConfig(cmd)
		domain := createDomain(cmd, args[0])

		err := cfg.PrepareDomain(domain)
		if err != nil {
			log.Fatal(err)
		}

		client := client.NewClient(cfg.CreateFileCache(), &http.Client{})

		response, err := client.Refresh(domain)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("provider responded \"%s\"", response)
	},
}

func init() {
	refreshCmd.AddCommand(domainCmd)

	domainCmd.Flags().String(flagUsername, "", "Set username used for authentication [<username>]")
	domainCmd.Flags().String(flagPassword, "", "Set password used for authentication [<password>]")
	domainCmd.Flags().String(flagDomainName, "", "Set your dns domain [<domain>]")
	domainCmd.Flags().String(flagIp4Address, "", "Set IPv4 address instead determining via wan request [<ip4>]")
	domainCmd.Flags().String(flagIp6Address, "", "Set IPv6 address instead determining via wan request [<ip6>]")
	domainCmd.Flags().String(flagIp6HostId, "", "Set IPv6 host id/interface id and use prefix + host id in the refresh url [<ip6>")
	domainCmd.Flags().String(flagHost, "", "Set host name of the service in the refresh url [<host>]")
	domainCmd.Flags().String(flagProtocol, config.DefaultProtocol, "Set protocol in the refresh url [<protocol>]")
	domainCmd.Flags().String(flagAuthMethod, config.DefaultAuthMethod, "Set authentication method in refresh requests")
	domainCmd.Flags().String(flagRequestMethod, config.DefaultRequestMethod, "Set request method in refresh requests")
	domainCmd.Flags().String(flagUserAgent, "", "Set user agent in refresh requests")
	domainCmd.Flags().Int(flagCacheModifiedLifetime, cache.ModifiedExpirySecondsDefault, "Set relative domain configuration cache lifetime in seconds [0 is disabled]")
	domainCmd.Flags().Int(flagCacheCreatedLifetime, cache.CreatedExpirySecondsDefault, "Set max domain configuration cache lifetime in seconds [0 is disabled]")
	domainCmd.Flags().SortFlags = false
}

func createConfig(cmd *cobra.Command) *config.Config {
	cfg := config.NewConfig(version)
	cfg.CacheCreatedExpirySeconds, _ = cmd.Flags().GetInt(flagCacheCreatedLifetime)
	cfg.CacheModifiedExpirySeconds, _ = cmd.Flags().GetInt(flagCacheModifiedLifetime)

	return cfg
}

func createDomain(cmd *cobra.Command, refreshUrl string) *config.Domain {
	authMethod, _ := cmd.Flags().GetString(flagAuthMethod)
	domainName, _ := cmd.Flags().GetString(flagDomainName)
	host, _ := cmd.Flags().GetString(flagHost)
	ip4Address, _ := cmd.Flags().GetString(flagIp4Address)
	ip6Address, _ := cmd.Flags().GetString(flagIp6Address)
	ip6HostId, _ := cmd.Flags().GetString(flagIp6HostId)
	password, _ := cmd.Flags().GetString(flagPassword)
	protocol, _ := cmd.Flags().GetString(flagProtocol)
	requestMethod, _ := cmd.Flags().GetString(flagRequestMethod)
	userAgent, _ := cmd.Flags().GetString(flagUserAgent)
	username, _ := cmd.Flags().GetString(flagUsername)

	template := &config.Template{
		AuthMethod:    authMethod,
		Host:          host,
		Protocol:      protocol,
		RefreshUrl:    refreshUrl,
		RequestMethod: requestMethod,
		UserAgent:     userAgent,
	}

	return &config.Domain{
		AuthUser:     username,
		AuthPassword: password,
		DomainName:   domainName,
		Ip4Address:   ip4Address,
		Ip6Address:   ip6Address,
		Ip6HostId:    ip6HostId,
		Template:     *template,
	}
}

func createFlagName(key string) string {
	return strings.Replace(key, "_", "-", -1)
}
