package crons

import (
	"errors"
	"fmt"

	"github.com/d2jvkpn/go-backend/pkg/infra"

	"github.com/d2jvkpn/gotk"
	_ "github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func Load(project *viper.Viper) (err error) {
	var (
		appName string
		config  *viper.Viper
	)

	// 1. Log
	appName = project.GetString("app_name")
	config, err = gotk.LoadYamlConfig(project.GetString("meta.config"), "config")
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			Exit()
		}
	}()

	// 2. databases: postgres, redis
	err = gotk.ConcRunErr(
		func() (err error) {
			_GORM_PG, _DB, err = infra.PgConnect(config.Sub("postgres"), true)
			return err
		},
		func() (err error) {
			_Redis, err = infra.NewRedisClient(config.Sub("redis"))
			return err
		},
	)
	if err != nil {
		return err
	}

	// 3.
	fmt.Printf("==> TODO crons: %s\n", appName)

	return err
}

func Exit() (err error) {
	var e error

	joinErr := func(e error) {
		err = errors.Join(err, e)
	}

	// 1. stop crons
	// TODO:

	// 2. close databases: postgres and redis
	e = gotk.ConcRunErr(
		func() error {
			if _Redis == nil {
				return nil
			}
			return _Redis.Close()
		},
		func() error {
			if _DB == nil {
				return nil
			}
			return _DB.Close()
		},
	)
	if e != nil {
		_Logger.Error("close databases", zap.String("error", e.Error()))
		joinErr(e)
	}

	// 6. close logger
	if _Logger != nil {
		if err == nil {
			_Logger.Info("exit")
		} else {
			_Logger.Error("exit", zap.String("error", e.Error()))
		}

		if e = _Logger.Down(); e != nil {
			joinErr(e)
		}
	}

	return err
}
