package crons

import (
	// "flag"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func NewCmd(name string) (command *cobra.Command) {
	var (
		config string
		fSet   *pflag.FlagSet
	)

	command = &cobra.Command{
		Use:   name,
		Short: "cron deamon",

		Run: func(cmd *cobra.Command, args []string) {
			run(config)
		},
	}

	fSet = command.Flags()
	fSet.StringVar(&config, "config", "configs/crons.yaml", "configuration file(yaml)")

	return command
}

func run(config string) {
	fmt.Printf("!!! TODO: %s\n", config)
}
