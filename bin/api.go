package bin

import (
	"embed"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/d2jvkpn/go-backend/internal"
	"github.com/d2jvkpn/go-backend/pkg/utils"

	"github.com/d2jvkpn/gotk"
	"github.com/spf13/viper"
)

func RunApi(project *viper.Viper, args []string, migrations embed.FS) {
	var (
		fSet         *flag.FlagSet
		release      bool
		config       string
		httpAddr     string
		internalAddr string
		grpcAddr     string

		err    error
		errCh  chan error
		logger *slog.Logger
	)

	// 1. setup
	// fmt.Println("~~~", args)
	fSet = flag.NewFlagSet("api", flag.ExitOnError)

	fSet.BoolVar(&release, "release", false, "run in release mode")
	fSet.StringVar(&config, "config", "configs/local.yaml", "configuration file(yaml)")

	fSet.StringVar(&httpAddr, "http.addr", ":9011", "http listening address")
	fSet.StringVar(&internalAddr, "internal.addr", ":9015", "internal listening address")
	fSet.StringVar(&grpcAddr, "grpc.addr", ":9016", "grpc listening address")

	fSet.Usage = func() {
		output := flag.CommandLine.Output()
		fmt.Fprintf(output, "Usage api:\n")
		fSet.PrintDefaults()
	}

	if err = fSet.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "API exit: %s\n", err)
		os.Exit(1)
		return
	}

	// logger = slog.New(slog.NewJSONHandler(os.Stderr, nil))
	if release {
		logger = utils.NewJSONLogger(os.Stderr, slog.LevelInfo)
	} else {
		logger = utils.NewJSONLogger(os.Stderr, slog.LevelDebug)
	}

	defer func() {
		if err != nil {
			logger.Error("API exit", "error", err)
			os.Exit(1)
		} else {
			logger.Info("API exit")
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
			"command":       "api",
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
		fmt.Sprintf("API is up"),
		"config", config,
		"release", release,
		"app_version", project.GetString("meta.app_version"),
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
