package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var thisversion string

var RootCmd = &cobra.Command{
	Use:           "shem",
	Short:         "A tool to help with creating Talos cluster",
	Long:          infoLongHelp,
	SilenceUsage:  true,
	SilenceErrors: true,
	Version:       thisversion,
}

func init() {
}

func Execute() error {
	// Parse only the persistent flags (like --cluster) before executing any command
	RootCmd.PersistentFlags().Parse(os.Args[1:])

	// Execute the root command and all subcommands
	return RootCmd.Execute()
}
