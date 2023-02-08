package main

import (
	_ "embed"

	"cli-tpl/cmd"
	embedutil "cli-tpl/pkg/embed"
)

//go:embed etc/default.yaml
var defaultConfig string

func main() {
	embedutil.DefaultConfig = defaultConfig
	cmd.Execute()
}
