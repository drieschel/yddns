/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"

	"github.com/drieschel/dddns/internal"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dddns",
	Short: "A lightweight dyndns client",
	Long:  `Drieschel's lightweight dyndns client`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dddns.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().IntSliceP("ip-version", "v", []int{4}, "Supported ip versions (4, 6)")
	rootCmd.Flags().StringP("domain", "d", "", "Domain name")
	rootCmd.Flags().StringP("user", "u", "", "User for authentication")
	rootCmd.Flags().StringP("password", "p", "", "Password for authentication")
}
