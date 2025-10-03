package cmd

import (
	"log"

	"github.com/drieschel/dddns/internal"
	"github.com/spf13/cobra"
)

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Refresh ip address(es) for a domain",
	Long:  `Ping a dyn dns server for a given domain`,
	Run: func(cmd *cobra.Command, args []string) {
		ipVersions, err := cmd.Flags().GetIntSlice("ip-version")
		if err != nil {
			log.Fatal(err)
		}

		for _, ipVersion := range ipVersions {
			internal.ValidateIpVersion(ipVersion)
		}
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pingCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	pingCmd.Flags().IntSliceP("ip-version", "v", []int{4}, "Supported ip versions (4, 6)")
	pingCmd.Flags().StringP("domain", "d", "", "domain name")
	pingCmd.Flags().StringP("user", "u", "", "User for authentication")
	pingCmd.Flags().StringP("password", "p", "", "Password for authentication")
}
