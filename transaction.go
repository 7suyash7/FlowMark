package main

import (
    "context"
    "time"
    "fmt"

    "github.com/onflow/cadence"
    "github.com/onflow/flow-go-sdk"
    "github.com/onflow/flow-go-sdk/access/http"
    "github.com/onflow/flow-go-sdk/crypto"
)

func SendTransaction(ctx context.Context, client *http.Client, senderAccount *flow.Account, sequenceNumber uint64) (time.Duration, string, flow.Identifier) {
    tx := flow.NewTransaction()
    var recipientAddressHex = LoadEnvVar("RECIPIENT_ADDRESS")
    var senderPrivateKeyHex = LoadEnvVar("SENDER_PRIVATE_KEY")

    script := `
    import FungibleToken from 0x9a0766d93b6608b7
    import FlowToken from 0x7e60df042a9c0868

    transaction(amount: UFix64, recipient: Address) {
        let sentVault: @FungibleToken.Vault
        prepare(signer: AuthAccount) {
        let vaultRef = signer.borrow<&FlowToken.Vault>(from: /storage/flowTokenVault)
            ?? panic("failed to borrow reference to sender vault")

        self.sentVault <- vaultRef.withdraw(amount: amount)
        }

        execute {
        let receiverRef =  getAccount(recipient)
            .getCapability(/public/flowTokenReceiver)
            .borrow<&{FungibleToken.Receiver}>()
            ?? panic("failed to borrow reference to recipient vault")

        receiverRef.deposit(from: <-self.sentVault)
        }
    }
    `
    tx.SetScript([]byte(script))
    tx.SetGasLimit(100)

    latestBlock, err := client.GetLatestBlockHeader(ctx, true)
    if err != nil {
        panic(err)
    }
    tx.SetReferenceBlockID(latestBlock.ID)

    tx.SetProposalKey(senderAccount.Address, senderAccount.Keys[0].Index, sequenceNumber)
    tx.SetPayer(senderAccount.Address)
    tx.AddAuthorizer(senderAccount.Address)

    amount, err := cadence.NewUFix64("1.234")
    if err != nil {
        panic(err)
    }

    if err = tx.AddArgument(amount); err != nil {
        panic(err)
    }

    recipient := cadence.NewAddress(flow.HexToAddress(recipientAddressHex))

    err = tx.AddArgument(recipient)
    if err != nil {
        panic(err)
    }

    sigAlgo := crypto.ECDSA_P256
    hashAlgo := crypto.SHA3_256
    privateKey, err := crypto.DecodePrivateKeyHex(sigAlgo, senderPrivateKeyHex)
    if err != nil {
        panic(err)
    }

    signer, err := crypto.NewInMemorySigner(privateKey, hashAlgo)
    if err != nil {
        panic(err)
    }

    if err = tx.SignEnvelope(senderAccount.Address, senderAccount.Keys[0].Index, signer); err != nil {
        panic(err)
    }

    sendStartTime := time.Now()

    if err = client.SendTransaction(ctx, *tx); err != nil {
        panic(err)
    }

    sendEndTime := time.Now()

    latency := sendEndTime.Sub(sendStartTime)

    txHex := tx.ID().Hex()
    fmt.Printf("hex: %s \n", txHex)
    return latency, txHex, tx.ID()
}
