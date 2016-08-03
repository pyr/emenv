package main

import (
	"github.com/pyr/emenv"
	"flag"
	"fmt"
	"os"
)

func main() {

	cfg := flag.String("c", os.ExpandEnv("${PWD}/Emenv"), "configuration path")
	flag.Parse()

	env, err := emenv.LoadEnv(*cfg)
	if err != nil {
		panic(err)
	}

	switch {
	case flag.Arg(0) == "sync":
		err = env.Sync()
	case flag.Arg(0) == "install":
		err = env.Install()
	default:
		fmt.Printf("unknown command: %s\n", flag.Arg(0))
		os.Exit(1)
	}
	if err != nil {
		panic(err)
	}
	os.Exit(0)
}
