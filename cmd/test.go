package cmd

import (
	"strings"

	"github.com/PrivatePuffin/shem/pkg/power"
	"github.com/PrivatePuffin/shem/pkg/prices"
	"github.com/spf13/cobra"
)

var testLongHelp = strings.TrimSpace(`
Fetches expected test for today and/or tomorrow if available

`)

var testCmd = &cobra.Command{
	Use:     "test",
	Short:   "Prints testrmation about the shem binary",
	Long:    testLongHelp,
	Example: "shem test",
	Run: func(cmd *cobra.Command, args []string) {
		prices.Fetch()
		prices.Render()
		power.Render()
	},
}

func init() {
	RootCmd.AddCommand(testCmd)
}
