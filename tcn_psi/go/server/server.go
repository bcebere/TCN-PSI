package server

import (
	"errors"
	psiserver "github.com/openmined/psi/server"
	"github.com/openmined/tcn-psi/tcn"
)

//TCNServer context for the server side of a TCN-Private Set Intersection-Cardinality protocol.
type TCNServer struct {
	context *psiserver.PsiServer
}

//CreateWithNewKey creates and returns a new server instance with a fresh private key.
//
//Returns an error if any crypto operations fail.
func CreateWithNewKey() (*TCNServer, error) {
	tcnServer := new(TCNServer)

	psiServer, err := psiserver.CreateWithNewKey(false)
	if err != nil {
		return nil, err
	}
	tcnServer.context = psiServer
	return tcnServer, nil
}

//CreateFromKey creates and returns a new server instance with the provided private key.
//
//Returns an error if any crypto operations fail.
func CreateFromKey(key []byte) (*TCNServer, error) {
	tcnServer := new(TCNServer)

	psiServer, err := psiserver.CreateFromKey(key, false)
	if err != nil {
		return nil, err
	}
	tcnServer.context = psiServer
	return tcnServer, nil
}

//CreateSetupMessage creates a setup message from the server's dataset to be sent to the
//client.
//
//Returns an error if the context is invalid or if the encryption fails.
func (s *TCNServer) CreateSetupMessage(fpr float64, inputCount int64, reports []*tcn.SignedReport) (string, error) {
	if s.context == nil {
		return "", errors.New("invalid context")
	}

	contacts := []string{}
	for idx := range reports {
		candidates, err := reports[idx].Report.TemporaryContactNumbers()
		if err != nil {
			return "", err
		}
		for jdx := range candidates {
			contacts = append(contacts, candidates[jdx].ToString())
		}
	}
	return s.context.CreateSetupMessage(fpr, inputCount, contacts)
}

//ProcessRequest processes a client query and returns the corresponding server response to
//be sent to the client.
//
//Returns an error if the context is invalid.
func (s *TCNServer) ProcessRequest(request string) (string, error) {
	if s.context == nil {
		return "", errors.New("invalid context")
	}
	return s.context.ProcessRequest(request)
}

//GetPrivateKeyBytes returns this instance's private key. This key should only be used to
//create other server instances. DO NOT SEND THIS KEY TO ANY OTHER PARTY!
func (s *TCNServer) GetPrivateKeyBytes() ([]byte, error) {
	if s.context == nil {
		return nil, errors.New("invalid context")
	}

	return s.context.GetPrivateKeyBytes()
}

//Version of the library.
func (s *TCNServer) Version() string {
	return s.context.Version()
}
