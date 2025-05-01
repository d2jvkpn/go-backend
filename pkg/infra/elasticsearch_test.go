package infra

import (
	"fmt"
	"testing"

	"github.com/d2jvkpn/gotk"
	es8 "github.com/elastic/go-elasticsearch/v8"
	"github.com/spf13/viper"
)

func TestES(t *testing.T) {
	var (
		err    error
		vp     *viper.Viper
		client *es8.Client
	)

	if vp, err = gotk.LoadYamlConfig("../../configs/local.yaml", "local"); err != nil {
		t.Fatal(err)
	}

	if client, err = NewEsClient(vp); err != nil {
		t.Fatal(err)
	}

	fmt.Println(client)
}
