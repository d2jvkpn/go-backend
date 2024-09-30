package bin

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/d2jvkpn/go-backend/internal/crons"

	"github.com/d2jvkpn/gotk"
	"github.com/spf13/viper"
)

func RunCrons(project *viper.Viper, args []string) {
	var (
		fSet   *flag.FlagSet
		config string
		err    error
		logger *slog.Logger
	)

	// 1. setup
	// fmt.Println("~~~", args)
	fSet = flag.NewFlagSet("crons", flag.ExitOnError)

	fSet.StringVar(&config, "config", "configs/local.yaml", "configuration file(yaml)")

	fSet.Usage = func() {
		output := flag.CommandLine.Output()
		fmt.Fprintf(output, "Usage crons:\n")
		fSet.PrintDefaults()
	}

	if err = fSet.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "crons exit: %s\n", err)
		os.Exit(1)
		return
	}

	//logger = slog.New(slog.NewJSONHandler(
	//	os.Stderr, &slog.HandlerOptions{AddSource: true},
	//).WithGroup("api"))
	logger = slog.New(slog.NewJSONHandler(os.Stderr, nil))

	defer func() {
		if err != nil {
			logger.Error("crons exit", "error", err)
			os.Exit(1)
		} else {
			logger.Info("crons exit")
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

	// 3.
	if err = crons.Load(project); err != nil {
		return
	}

	if err = crons.Run(project); err != nil {
		return
	}

	// 4. exit
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // linux: syscall.SIGUSR2

	sig := <-quit
	logger.Info("... received from channel quit", "signal", sig.String())

	err = crons.Exit()
}
