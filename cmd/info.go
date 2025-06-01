package cmd

import (
	"strings"

	"github.com/PrivatePuffin/public/shem/pkg/info"
	"github.com/spf13/cobra"
)

var infoLongHelp = strings.TrimSpace(`
shem is a tool to help you easily deploy and maintain a Talos Kubernetes Cluster.


Workflow:
  Create talconfig.yaml file defining your nodes information like so:

 Available commands
  > shem init
  > shem genconfig

`)

var infoCmd = &cobra.Command{
	Use:     "info",
	Short:   "Prints information about the shem binary",
	Long:    infoLongHelp,
	Example: "shem info",
	Run: func(cmd *cobra.Command, args []string) {
		info.NewInfo().Print()
	},
}

func init() {
	RootCmd.AddCommand(infoCmd)
}
