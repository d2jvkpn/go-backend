package api

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/d2jvkpn/go-backend/internal"
	"github.com/d2jvkpn/go-backend/internal/settings"

	"github.com/spf13/cobra"
)

func NewCmd(name string) (command *cobra.Command) {
	command = &cobra.Command{
		Use:   name,
		Short: "api service",

		Run: func(cmd *cobra.Command, args []string) {
			run(args)
		},
	}

	return command
}

func run(args []string) {
	var (
		fSet         *flag.FlagSet
		release      bool
		app_name     string
		httpAddr     string
		internalAddr string
		config       string
		count        int

		err    error
		errch  chan error
		quit   chan os.Signal
		logger *slog.Logger
	)

	// 1. setup project
	logger = slog.New(slog.NewJSONHandler(os.Stderr, nil))

	defer func() {
		if err != nil {
			logger.Error("exit", "error", err)
			os.Exit(1)
		}
	}()

	// fmt.Println("~~~", args)

	app_name = settings.Project.GetString("app_name")

	fSet = flag.NewFlagSet("api", flag.ExitOnError)

	fSet.BoolVar(&release, "release", false, "run in release mode")
	fSet.StringVar(&config, "config", "configs/local.yaml", "configuration file(yaml) path")
	fSet.StringVar(&httpAddr, "http.addr", ":9011", "http listening address")
	fSet.StringVar(&internalAddr, "internal.addr", ":9019", "internal listening address")

	fSet.Usage = func() {
		output := flag.CommandLine.Output()
		fmt.Fprintf(output, "api:\n")
		fSet.PrintDefaults()
	}

	fSet.Parse(args)

	// 2. configuration
	err = settings.Load(
		config,
		map[string]any{
			"release":       release,
			"http_addr":     httpAddr,
			"internal_addr": internalAddr,
		},
	)
	if err != nil {
		err = fmt.Errorf("settings.Load: %w", err)
		return
	}

	// 3.
	if err = internal.Load(release); err != nil {
		err = fmt.Errorf("internal.Load: %w", err)
		return
	}

	// 4. up
	if errch, err = internal.Run(httpAddr, internalAddr); err != nil {
		err = fmt.Errorf("internal.Run: %w", err)
		return
	}

	logger.Info(
		fmt.Sprintf("%s is up", app_name),
		"config", config,
		"release", release,
		"app_version", settings.Meta["app_version"],
	)

	// 5. exit
	quit = make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // linux: syscall.SIGUSR2

	syncErrors := func(count int) {
		for i := 0; i < count; i++ {
			err = errors.Join(err, <-errch)
		}
	}

	count = cap(errch)

	select {
	case err = <-errch:
		logger.Error("... received from channel errch")
		count -= 1
	case sig := <-quit:
		logger.Info("... received from channel quit", "signal", sig.String())
	}

	err = errors.Join(err, internal.Shutdown())
	syncErrors(count)
}
