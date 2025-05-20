package main

import (
	"embed"
	"flag"
	"fmt"
	"os"

	"backend-utils/internal"

	"github.com/d2jvkpn/gotk"
	"github.com/spf13/viper"
)

var (
	//go:embed migrations/*.sql
	_Migrations embed.FS
)

func main() {
	var (
		config    string
		directory string
		dsn       string
		err       error

		vp *viper.Viper
	)

	flag.StringVar(&config, "config", "configs/local.yaml", "file path of config")

	flag.StringVar(
		&directory, "directory", "",
		"the directory where the SQL migration files are located, if not set, the embedded SQL files will be used.",
	)

	flag.Parse()

	defer func() {
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		} else {
			fmt.Fprintf(os.Stderr, "Exit\n")
		}
	}()

	if vp, err = gotk.LoadYamlConfig(config, "config"); err != nil {
		return
	}

	if dsn = vp.GetString("postgres.dsn"); dsn == "" {
		err = fmt.Errorf("postgres.dsn isn't set in config")
		return
	}

	if directory == "" {
		err = internal.MigratePgFs(dsn, _Migrations, "migrations")
	} else {
		err = internal.MigratePgDir(dsn, directory)
	}
	if err != nil {
		return
	}
}
