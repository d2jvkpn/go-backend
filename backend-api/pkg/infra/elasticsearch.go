package infra

import (
	"io/ioutil"

	es "github.com/elastic/go-elasticsearch/v9"
	"github.com/spf13/viper"
)

func NewEsClient(vp *viper.Viper) (client *es.Client, err error) {
	var config es.Config

	config.Addresses = vp.GetStringSlice("addresses")
	config.Username = vp.GetString("username")
	config.Password = vp.GetString("password")
	config.APIKey = vp.GetString("api_key")

	if fp := vp.GetString("ca_cert_file"); fp != "" {
		if config.CACert, err = ioutil.ReadFile(fp); err != nil {
			return nil, err
		}
	}

	return es.NewClient(config)
}
