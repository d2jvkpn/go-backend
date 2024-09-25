package main

import (
	"embed"
	"fmt"
	"log/slog"
	"os"

	"github.com/d2jvkpn/go-backend/bin/api"
	"github.com/d2jvkpn/go-backend/bin/crons"
	"github.com/d2jvkpn/go-backend/internal/settings"

	"github.com/d2jvkpn/gotk"
	"github.com/spf13/cobra"
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
		err        error
		logger     *slog.Logger
		command    *cobra.Command
		showConfig *cobra.Command
		showBuild  *cobra.Command
	)

	logger = slog.New(slog.NewJSONHandler(os.Stderr, nil))

	defer func() {
		if err != nil {
			logger.Error("exit", "error", err)
			os.Exit(1)
		}
	}()

	if err = settings.Setup(_Project, _Migrations); err != nil {
		err = fmt.Errorf("settings.Setup: %w", err)
		return
	}

	showConfig = &cobra.Command{
		Use:   "show-config",
		Short: "show configurations",

		Run: func(cmd *cobra.Command, args []string) {
			var arg string

			if len(args) == 0 {
				fmt.Fprintf(os.Stderr, "arg required: api | crons | swagger\n")
				os.Exit(1)
			} else {
				arg = args[0]
			}

			switch arg {
			case "api":
				fmt.Printf("%s\n", settings.Project.GetString("api_config"))
			case "crons":
				fmt.Printf("%s\n", settings.Project.GetString("crons_config"))
			case "swagger":
				fmt.Printf("%s\n", settings.Project.GetString("swagger_config"))
			default:
				fmt.Fprintf(os.Stderr, "arg required: api | crons | swagger\n")
				os.Exit(1)
			}
		},
	}

	showBuild = &cobra.Command{
		Use:   "show-build",
		Short: "show build information",

		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s\n", gotk.BuildInfoText(settings.Meta))
		},
	}

	command = &cobra.Command{Use: settings.Project.GetString("app_name")}

	command.AddCommand(api.NewCmd("api"))
	command.AddCommand(crons.NewCmd("crons"))
	command.AddCommand(showConfig)
	command.AddCommand(showBuild)

	command.Execute()
}
