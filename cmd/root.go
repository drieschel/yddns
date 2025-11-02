package cmd

import (
	"os"

	"github.com/drieschel/yddns/internal/config"
	"github.com/spf13/cobra"
)

var (
	version = config.DefaultAppVersion
	rootCmd = &cobra.Command{
		Version: version,
		Use:     "yddns",
		Short:   "A flexible and lightweight dyndns client",
		Long:    `Drieschel's flexible and lightweight dyndns client`,
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
