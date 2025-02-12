package main

import (
	"fmt"
	"os"

	"github.com/ckm54/env"
)

var bindAddress = env.String("BIND_ADDRESS", true, "", "bind address for server, i.e. localhost")
var bindPort = env.Int("BIND_PORT", true, 0, "bind port for the server, i.e. 9090")

func main() {
	err := env.Parse()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("BIND_ADDRESS", bindAddress)
	fmt.Println("BIND_PORT", bindPort)
}
