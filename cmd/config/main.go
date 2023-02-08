package config

import (
	"github.com/spf13/cobra"
)

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration operation",
}

func init() {
	ConfigCmd.AddCommand(DefaultCmd)
}
