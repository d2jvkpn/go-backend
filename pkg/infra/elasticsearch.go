package infra

import (
	"io/ioutil"

	es8 "github.com/elastic/go-elasticsearch/v8"
	"github.com/spf13/viper"
)

func NewEsClient(vp *viper.Viper) (client *es8.Client, err error) {
	var config es8.Config

	config.Addresses = vp.GetStringSlice("addresses")
	config.Username = vp.GetString("username")
	config.Password = vp.GetString("password")
	config.APIKey = vp.GetString("api_key")

	if fp := vp.GetString("ca_cert_file"); fp != "" {
		if config.CACert, err = ioutil.ReadFile(fp); err != nil {
			return nil, err
		}
	}

	return es8.NewClient(config)
}
