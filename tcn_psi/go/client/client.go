package client

import (
	"fmt"
	"github.com/openmined/psi/client"
)

func Client() {
	psiClient, err := client.Create()
	if err == nil {
		fmt.Println("client loaded")
		psiClient.Destroy()
	}
}
