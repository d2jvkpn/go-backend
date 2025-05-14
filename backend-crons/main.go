package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend-crons/internal"

	"github.com/d2jvkpn/gotk"
	"github.com/spf13/viper"
)

var (
	//go:embed project.yaml
	_Project []byte
)

func main() {
	var (
		config  string
		err     error
		logger  *slog.Logger
		project *viper.Viper
	)

	// 1. setup
	if project, err = gotk.ProjectFromBytes(_Project); err != nil {
		err = fmt.Errorf("Failed to load project.yaml: %w", err)
		return
	}

	flag.StringVar(&config, "config", "configs/backend-crons.local.yaml", "configuration file(yaml)")

	flag.Usage = func() {
		output := flag.CommandLine.Output()
		fmt.Fprintf(output, "Usage:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	//logger = slog.New(slog.NewJSONHandler(
	//	os.Stderr, &slog.HandlerOptions{AddSource: true},
	//).WithGroup("api"))
	logger = slog.New(slog.NewJSONHandler(os.Stderr, nil))

	defer func() {
		if err != nil {
			logger.Error("Crons exit", "error", err)
			os.Exit(1)
		} else {
			logger.Info("Crons exit")
		}
	}()

	// 2. configuration
	updateMeta(
		project,
		map[string]any{
			"config":     config,
			"command":    "crons",
			"startup_at": time.Now().Format(gotk.RFC3339Milli),
		},
	)

	// 3. load
	if err = internal.Load(project); err != nil {
		return
	}

	// 4. up
	if err = internal.Run(project); err != nil {
		return
	}

	logger.Info(
		"Crons is up",
		"config", config,
		"app_version", project.GetString("meta.app_version"),
	)

	// 5. exit
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // linux: syscall.SIGUSR2

	sig := <-quit
	logger.Info("... received from channel quit", "signal", sig.String())

	err = internal.Exit()
}

func updateMeta(project *viper.Viper, mp map[string]any) {
	meta := project.GetStringMap("meta")

	for k, v := range mp {
		meta[k] = v
	}
}
