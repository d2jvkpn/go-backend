package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/d2jvkpn/go-backend/bin"

	"github.com/d2jvkpn/gotk"
	"github.com/spf13/cobra"
	// "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	//go:embed project.yaml
	_Project []byte

	//go:embed migrations/*.sql
	_Migrations embed.FS
)

// ./target/main api -- --release
func main() {
	var (
		err error
		// fSet   *pflag.FlagSet
		project *viper.Viper

		command *cobra.Command
	)

	defer func() {
		if err != nil {
			fmt.Fprintf(os.Stderr, "main exit: %s\n", err)
			os.Exit(1)
		}
	}()

	if project, err = gotk.ProjectFromBytes(_Project); err != nil {
		err = fmt.Errorf("settings.LoadProject: %w", err)
		return
	}

	command = &cobra.Command{
		Use: project.GetString("meta.app_name"),
	}
	/*
		fSet = command.Flags()
		fSet.StringVar(&config, "config", "configs/local.yaml", "configuration file(yaml)")
		fSet.BoolVar(&release, "release", false, "run in release mode")

		if err = fSet.Parse(os.Args[1:]); err != nil {
			fmt.Println("~~~ error:", err)
		}
		command.SetArgs(fSet.Args())
	*/

	showCmd := &cobra.Command{
		Use:   "show",
		Short: "show build information(build) and configuration(api, crons, swagger)",
		Run: func(cmd *cobra.Command, args []string) {
			errMsg := "required: build | api | crons | swagger\n"

			if len(args) == 0 {
				fmt.Fprintf(os.Stderr, errMsg)
				os.Exit(1)
			}

			switch args[0] {
			case "build":
				fmt.Printf("%s\n", gotk.BuildInfoText(project.GetStringMap("meta")))
			case "api":
				fmt.Printf("%s\n", project.GetString("api_config"))
			case "crons":
				fmt.Printf("%s\n", project.GetString("crons_config"))
			case "swagger":
				fmt.Printf("%s\n", project.GetString("swagger_config"))
			default:
				fmt.Fprintf(os.Stderr, errMsg)
				os.Exit(1)
			}
		},
	}

	apiCmd := &cobra.Command{
		Use:   "api",
		Short: "api service",

		Run: func(cmd *cobra.Command, args []string) {
			// fmt.Println("~~~ api args:", args)
			bin.RunApi(project, args, _Migrations)
		},
	}

	cronsCmd := &cobra.Command{
		Use:   "crons",
		Short: "cron deamon",

		Run: func(cmd *cobra.Command, args []string) {
			bin.RunCrons(args)
		},
	}

	swaggerCmd := &cobra.Command{
		Use:   "swagger",
		Short: "swagger service",

		Run: func(cmd *cobra.Command, args []string) {
			bin.RunBin("swagger", args)
		},
	}

	command.AddCommand(showCmd)
	command.AddCommand(apiCmd)
	command.AddCommand(cronsCmd)
	command.AddCommand(swaggerCmd)

	command.Execute()
}
