package cmd

import (
	"fmt"
	"strings"

	"github.com/drieschel/yddns/internal/config"
	"github.com/spf13/cobra"
)

// domainCmd represents the domain command
var domainCmd = &cobra.Command{
	Use:   "domain [refresh url or template name]",
	Short: "Refresh a domain via cli",
	Long:  `Refresh a single domain via command"`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Arg 0: %s\n", args[0])
	},
}

func init() {
	refreshCmd.AddCommand(domainCmd)

	domainCmd.Flags().String(strings.Replace(config.DomainKeyAuthMethod, "_", "-", -1), config.DefaultAuthMethod, "Set authentication method used for the service")
	domainCmd.Flags().String(strings.Replace(config.DomainKeyDomainName, "_", "-", -1), "", "Set name of the domain in the refresh URL [<domain>]")
	domainCmd.Flags().String(strings.Replace(config.DomainKeyHost, "_", "-", -1), "", "Set host name of the service [<host>]")
	domainCmd.Flags().String(strings.Replace(config.DomainKeyIp4Address, "_", "-", -1), "", "Set IPv4 address instead determining via wan request [<ip4>]")
	domainCmd.Flags().String(strings.Replace(config.DomainKeyIp6Address, "_", "-", -1), "", "Set IPv6 address instead determining via wan request [<ip6>]")
	domainCmd.Flags().String(strings.Replace(config.DomainKeyIp6HostId, "_", "-", -1), "", "Set IPv6 host id/interface id and use prefix + host id")
	domainCmd.Flags().String(strings.Replace(config.DomainKeyPassword, "_", "-", -1), "", "Set password used to authenticate [<password>]")
	domainCmd.Flags().String(strings.Replace(config.DomainKeyProtocol, "_", "-", -1), config.DefaultProtocol, "Set protocol in the refresh URL [<protocol>]")
	domainCmd.Flags().String(strings.Replace(config.DomainKeyRefreshUrl, "_", "-", -1), "", "Set URL to refresh the domain including placeholders")
	domainCmd.Flags().String(strings.Replace(config.DomainKeyRequestMethod, "_", "-", -1), config.DefaultRequestMethod, "Set request method of the service")
	domainCmd.Flags().String(strings.Replace(config.DomainKeyUserAgent, "_", "-", -1), "", "Set user agent in refresh requests")
	domainCmd.Flags().String(strings.Replace(config.DomainKeyUsername, "_", "-", -1), "", "Set username used to authenticate [<username>]")
}
