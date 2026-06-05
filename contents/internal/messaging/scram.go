package messaging

import (
	"crypto/sha256"
	"crypto/sha512"
	"hash"

	"github.com/xdg-go/scram"
)

// SHA256 and SHA512 are hash generator functions for SCRAM authentication.
var (
	SHA256 scram.HashGeneratorFcn = func() hash.Hash { return sha256.New() }
	SHA512 scram.HashGeneratorFcn = func() hash.Hash { return sha512.New() }
)

// XDGSCRAMClient implements sarama.SCRAMClient using the xdg-go/scram package.
type XDGSCRAMClient struct {
	*scram.Client
	*scram.ClientConversation
	scram.HashGeneratorFcn
}

func (x *XDGSCRAMClient) Begin(userName, password, authzID string) error {
	var err error
	x.Client, err = x.HashGeneratorFcn.NewClient(userName, password, authzID)
	if err != nil {
		return err
	}
	x.ClientConversation = x.Client.NewConversation()
	return nil
}

func (x *XDGSCRAMClient) Step(challenge string) (string, error) {
	return x.ClientConversation.Step(challenge)
}

func (x *XDGSCRAMClient) Done() bool {
	return x.ClientConversation.Done()
}
