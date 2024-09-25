package settings

import (
	// "fmt"
	"bytes"
	"embed"

	"github.com/d2jvkpn/gotk"
	"github.com/spf13/viper"
)

var (
	Meta       map[string]any
	Migrations embed.FS

	Project *viper.Viper
	Config  *viper.Viper
)

func init() {
	Project = viper.New()
	Config = viper.New()

}

func Setup(bts []byte, migrations embed.FS) (err error) {
	Migrations = migrations

	Project.SetConfigType("yaml")

	// _Project.ReadConfig(strings.NewReader(str))
	if err = Project.ReadConfig(bytes.NewReader(bts)); err != nil {
		return err
	}

	Meta = gotk.BuildInfo()
	Meta["app_name"] = Project.GetString("app_name")
	Meta["app_version"] = Project.GetString("app_version")

	return nil
}

func Load(config string, mp map[string]any) (err error) {
	return nil
}
