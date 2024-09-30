package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/d2jvkpn/go-backend/bin"

	"github.com/d2jvkpn/gotk"
	"github.com/spf13/viper"
)

var (
	//go:embed project.yaml
	_Project []byte

	//go:embed migrations/*.sql
	_Migrations embed.FS
)

func main() {
	var (
		err     error
		project *viper.Viper
		command *gotk.Command
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
	if project.GetString("app_name") == "" || project.GetString("app_version") == "" {
		err = fmt.Errorf("either app_name or app_version is unset in project.yaml")
		return
	}

	command = gotk.NewCommand(project.GetString("meta.app_name"), project)

	command.AddCmd(
		"config",
		"show configuration(api, crons, swagger)",
		func(args []string) {
			errMsg := "required: api | crons | swagger\n"

			if len(args) == 0 {
				fmt.Fprintf(os.Stderr, errMsg)
				os.Exit(1)
			}

			switch args[0] {
			case "api", "crons", "swagger":
				fmt.Printf("%s\n", project.GetString(args[0]+"_config"))
			default:
				fmt.Fprintf(os.Stderr, errMsg)
				os.Exit(1)
			}
		},
	)

	command.AddCmd(
		"api",
		"api service",
		func(args []string) { bin.RunApi(project, args, _Migrations) },
	)

	command.AddCmd(
		"crons",
		"cron deamon",
		func(args []string) { bin.RunCrons(project, args) },
	)

	command.AddCmd(
		"swagger",
		"swagger service",
		func(args []string) { bin.RunBin("swagger", args) },
	)

	command.Execute(os.Args[1:])
}
