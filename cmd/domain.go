package cmd

import (
	"log"
	"net/http"
	"strings"

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
)

var domainCmd = &cobra.Command{
	Use:   "domain [refresh-url | :template-name]",
	Short: "Refresh a domain config via cli",
	Long:  `Refresh a single domain config via command"`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		domain := createDomain(cmd, args[0])
		cfg := config.NewConfig(version)

		err := cfg.PrepareDomain(domain)
		if err != nil {
			log.Fatal(err)
		}

		client := client.NewClient(&http.Client{})

		response, err := client.Refresh(domain)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("provider responded \"%s\"", response)
	},
}

func init() {
	refreshCmd.AddCommand(domainCmd)

	domainCmd.Flags().String(flagAuthMethod, config.DefaultAuthMethod, "Set authentication method used for the service")
	domainCmd.Flags().String(flagDomainName, "", "Set name of the domain in the refresh URL [<domain>]")
	domainCmd.Flags().String(flagHost, "", "Set host name of the service [<host>]")
	domainCmd.Flags().String(flagIp4Address, "", "Set IPv4 address instead determining via wan request [<ip4>]")
	domainCmd.Flags().String(flagIp6Address, "", "Set IPv6 address instead determining via wan request [<ip6>]")
	domainCmd.Flags().String(flagIp6HostId, "", "Set IPv6 host id/interface id and use prefix + host id")
	domainCmd.Flags().String(flagPassword, "", "Set password used to authenticate [<password>]")
	domainCmd.Flags().String(flagProtocol, config.DefaultProtocol, "Set protocol in the refresh URL [<protocol>]")
	domainCmd.Flags().String(flagRequestMethod, config.DefaultRequestMethod, "Set request method of the service")
	domainCmd.Flags().String(flagUserAgent, "", "Set user agent in refresh requests")
	domainCmd.Flags().String(flagUsername, "", "Set username used to authenticate [<username>]")
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
