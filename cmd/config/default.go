package config

import (
	"fmt"

	"github.com/spf13/cobra"

	embedutil "cli-tpl/pkg/embed"
)

var DefaultCmd = &cobra.Command{
	Use:   "default",
	Short: "Display default configuration",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(embedutil.DefaultConfig)
	},
}
