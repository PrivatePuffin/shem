package cmd

import (
	"strings"

	"github.com/PrivatePuffin/shem/pkg/prices"
	"github.com/spf13/cobra"
)

var pricesLongHelp = strings.TrimSpace(`
Fetches expected prices for today and/or tomorrow if available

`)

var pricesCmd = &cobra.Command{
	Use:     "prices",
	Short:   "Prints pricesrmation about the shem binary",
	Long:    pricesLongHelp,
	Example: "shem prices",
	Run: func(cmd *cobra.Command, args []string) {
		prices.Fetch()
		prices.Render()
	},
}

func init() {
	RootCmd.AddCommand(pricesCmd)
}
