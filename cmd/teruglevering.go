package cmd

import (
	"strings"

	"github.com/PrivatePuffin/shem/pkg/power"
	"github.com/PrivatePuffin/shem/pkg/prices"
	"github.com/spf13/cobra"
)

var terugleveringLongHelp = strings.TrimSpace(`
Fetches expected teruglevering for today and/or tomorrow if available

`)

var terugleveringCmd = &cobra.Command{
	Use:     "teruglevering",
	Short:   "Prints terugleveringrmation about the shem binary",
	Long:    terugleveringLongHelp,
	Example: "shem teruglevering",
	Run: func(cmd *cobra.Command, args []string) {
		prices.Fetch()
		power.Render()
	},
}

func init() {
	RootCmd.AddCommand(terugleveringCmd)
}
