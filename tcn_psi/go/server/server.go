package main

import (
	"fmt"
	"github.com/openmined/psi/server"
)

func Server() {
	psiServer, err := server.CreateWithNewKey()
	if err == nil {
		fmt.Println("server loaded")
		psiServer.Destroy()
	}
}
