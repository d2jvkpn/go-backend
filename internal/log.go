package internal

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/d2jvkpn/go-backend/internal/settings"

	"github.com/d2jvkpn/gotk"
	"go.uber.org/zap"
)

func SetupLog(appName string, release bool) (err error) {
	log_file := filepath.Join("logs", appName+".log")

	if release {
		settings.Logger, err = gotk.NewZapLogger(log_file, zap.InfoLevel, 1024)
	} else {
		settings.Logger, err = gotk.NewZapLogger("", zap.DebugLevel, 1024)
	}
	if err != nil {
		return fmt.Errorf("NewLogger: %w", err)
	}

	_SLogger = slog.New(slog.NewJSONHandler(os.Stderr, nil))
	_Logger = settings.Logger.Named("internal")

	return nil
}
