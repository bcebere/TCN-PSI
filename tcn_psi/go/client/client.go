package client

import (
	"errors"
	psiclient "github.com/openmined/psi/client"
)

//TCNClient context for the client side of a TCN-Private Set Intersection-Cardinality protocol.
type TCNClient struct {
	context *psiclient.PsiClient
}

//Create returns a new TCN-PSI client
func Create() (*TCNClient, error) {
	tcnClient := new(TCNClient)

	psiClient, err := psiclient.Create()
	if err != nil {
		return nil, err
	}
	tcnClient.context = psiClient
	return tcnClient, nil
}

//CreateRequest generates a request message to be sent to the server.
//
//Returns an error if the context is invalid or if the encryption fails.
func (c *TCNClient) CreateRequest(rawInput []string) (string, error) {
	if c.context == nil {
		return "", errors.New("invalid context")
	}

	return c.context.CreateRequest(rawInput)
}

//ProcessResponse processes the server's response and returns the PSI cardinality.
//
//Returns an error if the context is invalid,  if any input messages are malformed or if decryption fails.
func (c *TCNClient) ProcessResponse(serverSetup, serverResponse string) (int64, error) {
	if c.context == nil {
		return 0, errors.New("invalid context")
	}

	return c.context.ProcessResponse(serverSetup, serverResponse)
}

//Version of the library.
func (c *TCNClient) Version() string {
	return c.context.Version()
}
