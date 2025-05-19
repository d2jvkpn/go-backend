package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"backend-api/internal"
	"backend-api/pkg/utils"

	"github.com/d2jvkpn/gotk"
	"github.com/spf13/viper"
)

var (
	//go:embed project.yaml
	_Project []byte
)

func main() {
	var (
		release      bool
		config       string
		httpAddr     string
		internalAddr string
		grpcAddr     string
		err          error

		project *viper.Viper

		errCh  chan error
		logger *slog.Logger
	)

	// 1. setup
	if project, err = gotk.ProjectFromBytes(_Project); err != nil {
		err = fmt.Errorf("Failed to load project.yaml: %w", err)
		return
	}

	flag.BoolVar(&release, "release", false, "run in release mode")
	flag.StringVar(&config, "config", "configs/backend-api.local.yaml", "configuration file(yaml)")

	flag.StringVar(&httpAddr, "http.addr", ":4011", "http listening address")
	flag.StringVar(&internalAddr, "internal.addr", ":4021", "internal listening address")

	flag.StringVar(&grpcAddr, "grpc.addr", ":4031", "grpc listening address")

	flag.Usage = func() {
		output := flag.CommandLine.Output()
		fmt.Fprintf(output, "Usage:\n")
		flag.PrintDefaults()
		fmt.Printf("\n\nConfig:\n```yaml\n%s\n```\n", project.GetString("config"))
	}

	flag.Parse()

	// logger = slog.New(slog.NewJSONHandler(os.Stderr, nil))
	if release {
		logger = utils.NewJSONLogger(os.Stderr, slog.LevelInfo)
	} else {
		logger = utils.NewJSONLogger(os.Stderr, slog.LevelDebug)
	}

	defer func() {
		if err != nil {
			logger.Error("Exit", "error", err)
			os.Exit(1)
		} else {
			logger.Info("Exit")
		}
	}()

	// 2. configuration
	updateMeta(
		project,
		map[string]any{
			"config":        config,
			"release":       release,
			"http_addr":     httpAddr, // don't use http.addr as key here
			"internal_addr": internalAddr,
			"grpc_addr":     grpcAddr,
			"startup_at":    time.Now().Format(gotk.RFC3339Milli),
		},
	)

	// 3. load
	if err = internal.Load(project); err != nil {
		err = fmt.Errorf("Faild to load: %w", err)
		return
	}

	// 4. up
	if errCh, err = internal.Run(project); err != nil {
		err = fmt.Errorf("Failed to run: %w", err)
		return
	}

	logger.Info(
		fmt.Sprintf("Service is up"),
		"release", release,
		"app_version", project.GetString("meta.app_version"),
		"config", config,
		"http_addr", httpAddr,
		"internal_addr", internalAddr,
		"grpc_addr", grpcAddr,
	)

	// 5. exit
	err = gotk.ExitChan(errCh, internal.Exit)
}

func updateMeta(project *viper.Viper, mp map[string]any) {
	meta := project.GetStringMap("meta")

	for k, v := range mp {
		meta[k] = v
	}
}
