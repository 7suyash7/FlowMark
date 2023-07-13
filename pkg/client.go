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
