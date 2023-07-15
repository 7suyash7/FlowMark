package pkg

import (
	"context"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/access/http"
)

func InitializeClient(host string) (*http.Client, error) {
	return http.NewClient(host)
}

func GetAccount(ctx context.Context, client *http.Client, address flow.Address) (*flow.Account, error) {
	return client.GetAccount(ctx, address)
}

func GetInitialSequenceNumber(account *flow.Account) uint64 {
	return account.Keys[0].SequenceNumber
}

func GetSequenceNumber(account *flow.Account, keyIndex int) (uint64, int) {
	for keyIndex >= len(account.Keys) {
		keyIndex -= len(account.Keys)
	}
	// NOTE: added returning keyIndex here, cause it's already being calculated here,
	// 		 Or we would run into runtime error in sendTransaction error while setting Proposal Key
	//		 and while signing the envelope.
	return account.Keys[keyIndex].SequenceNumber, keyIndex
}