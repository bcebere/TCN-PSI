package main

import (
	"github.com/openmined/tcn-psi/client"
	"github.com/openmined/tcn-psi/server"
)

func main() {
	client.Client()
	server.Server()
}
