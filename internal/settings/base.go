package settings

import (
	// "fmt"
	"bytes"

	"github.com/d2jvkpn/gotk"
	"github.com/spf13/viper"
)

func LoadProject(bts []byte) (project *viper.Viper, err error) {
	var meta map[string]any

	project = viper.New()
	project.SetConfigType("yaml")

	// _Project.ReadConfig(strings.NewReader(str))
	if err = project.ReadConfig(bytes.NewReader(bts)); err != nil {
		return nil, err
	}

	meta = gotk.BuildInfo()
	meta["app_name"] = project.GetString("app_name")
	meta["app_version"] = project.GetString("app_version")
	project.Set("meta", meta)

	return project, nil
}
