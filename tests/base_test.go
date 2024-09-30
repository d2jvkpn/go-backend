package tests

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/d2jvkpn/gotk"
	"github.com/spf13/viper"
)

var (
	_TestFlags  *flag.FlagSet
	_TestCtx    context.Context
	_TestConfig *viper.Viper
)

func TestMain(m *testing.M) {
	var (
		config string
		err    error
	)

	_TestFlags = flag.NewFlagSet("tests", flag.ExitOnError)
	flag.Parse() // must do

	_TestFlags.StringVar(&config, "config", "../configs/local.yaml", "config filepath")

	_TestFlags.Parse(flag.Args())
	fmt.Printf("==> load config: %q\n", config)

	defer func() {
		if err != nil {
			fmt.Printf("!!! TestMain: %v\n", err)
			os.Exit(1)
		}
	}()

	_TestCtx = context.TODO()

	_TestConfig, err = gotk.LoadYamlConfig(config, "config")
	if err != nil {
		return
	}

	m.Run()
}

// go test -- --config=../configs/local.yaml grpc -addr=localhost:9016
func TestClients(t *testing.T) {
	var (
		cmd  string
		args []string
	)

	if args = _TestFlags.Args(); len(args) == 0 {
		return
	}
	fmt.Println("==> args:", args)
	cmd = args[0]

	switch cmd {
	case "grpc":
		testGrpcClient(args[1:])
	default:
		t.Fatalf("unkonwn command: %s", cmd)
	}
}
