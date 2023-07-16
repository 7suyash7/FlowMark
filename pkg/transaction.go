package pkg

import (
    "context"
    "time"
    "fmt"
    "strings"
    "log"

    "github.com/onflow/cadence"
    "github.com/onflow/flow-go-sdk"
    "github.com/onflow/flow-go-sdk/access/http"
    "github.com/onflow/flow-go-sdk/crypto"
)



func SendTransaction(ctx context.Context, client *http.Client, senderAccount *flow.Account, sequenceNumber uint64, keyID int) (time.Duration, time.Duration, string, flow.Identifier, bool) {
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

    var latestBlock *flow.BlockHeader
    var fetchErr error

    maxRetryAttempts := 5 // adjust this as needed
    retryInterval := time.Second * 1 // adjust this as needed

    for attempts := 0; attempts < maxRetryAttempts; attempts++ {
        latestBlock, fetchErr = client.GetLatestBlockHeader(ctx, true)
        if fetchErr == nil {
            // We successfully fetched the block, no need to retry anymore.
            break
        }
        // Optionally log the error and the retry attempt.
        log.Printf("Failed to fetch the block (attempt %d): %v", attempts+1, fetchErr)
        // Pause for a while before next retry.
        time.Sleep(retryInterval)
    }
    if fetchErr != nil {
        // We failed to fetch the block even after retrying.
        // Print the error instead of panicking.
        fmt.Printf("Error fetching the block: %v\n", fetchErr)
        return 0, 0, "", flow.Identifier{}, false
    }
    
    tx.SetReferenceBlockID(latestBlock.ID)

    tx.SetProposalKey(senderAccount.Address, senderAccount.Keys[keyID].Index, sequenceNumber)
    tx.SetPayer(senderAccount.Address)
    tx.AddAuthorizer(senderAccount.Address)

    amount, err := cadence.NewUFix64("1.234")
    if err != nil {
        fmt.Printf("Error creating UFix64 amount: %v\n", err)
        return 0, 0, "", flow.Identifier{}, false
    }

    if err = tx.AddArgument(amount); err != nil {
        fmt.Printf("Error adding amount argument: %v\n", err)
        return 0, 0, "", flow.Identifier{}, false
    }

    recipient := cadence.NewAddress(flow.HexToAddress(recipientAddressHex))

    err = tx.AddArgument(recipient)
    if err != nil {
        fmt.Printf("Error adding recipient argument: %v\n", err)
        return 0, 0, "", flow.Identifier{}, false
    }

    sigAlgo := crypto.ECDSA_P256
    hashAlgo := crypto.SHA3_256
    privateKey, err := crypto.DecodePrivateKeyHex(sigAlgo, senderPrivateKeyHex)
    if err != nil {
        fmt.Printf("Error decoding private key: %v\n", err)
        return 0, 0, "", flow.Identifier{}, false
    }

    signer, err := crypto.NewInMemorySigner(privateKey, hashAlgo)
    if err != nil {
        fmt.Printf("Error creating signer: %v\n", err)
        return 0, 0, "", flow.Identifier{}, false
    }

    if err = tx.SignEnvelope(senderAccount.Address, senderAccount.Keys[0].Index, signer); err != nil {
        fmt.Printf("Error signing envelope: %v\n", err)
        return 0, 0, "", flow.Identifier{}, false
    }
    if keyID != 0 {
        if err = tx.SignEnvelope(senderAccount.Address, senderAccount.Keys[keyID].Index, signer); err != nil {
            fmt.Printf("Error signing envelope: %v\n", err)
            return 0, 0, "", flow.Identifier{}, false
        }
    }

    startTime := time.Now()
    if err = client.SendTransaction(ctx, *tx); err != nil {
        fmt.Printf("Error sending transaction: %v\n", err)
        return 0, 0, "", flow.Identifier{}, false
    }
    txEndTime := time.Now()
    txHex := tx.ID().Hex()
    fmt.Printf("hex: %s \n", txHex)

    WaitForSeal(ctx, client, tx.ID())
    sealEndTime := time.Now()

    txLatency := txEndTime.Sub(startTime)
    sealLatency := sealEndTime.Sub(startTime)

    return txLatency, sealLatency, txHex, tx.ID(), true
}

func AddKeys(ctx context.Context, client *http.Client, senderAccount *flow.Account, sequenceNumber uint64, numOfKeysToAdd int) error {
	tx := flow.NewTransaction()
	var senderPrivateKeyHex = LoadEnvVar("SENDER_PRIVATE_KEY")
	publicKeyHex := strings.TrimPrefix(fmt.Sprintf("%+v", senderAccount.Keys[0].PublicKey), "0x")

	script := `
		transaction(publicKey: String, numOfKeysToAdd: Int) {
			prepare(signer: AuthAccount) {
				let bytes = publicKey.decodeHex()
				let key = PublicKey(
					publicKey: bytes,
					signatureAlgorithm: SignatureAlgorithm.ECDSA_P256
				)

				var counter = 0
				while counter < numOfKeysToAdd {
					counter = counter + 1
					signer.keys.add(
						publicKey: key,
						hashAlgorithm: HashAlgorithm.SHA3_256,
						weight: 0.0
					)
				}
			}

			execute {

			}
		}
	`
	tx.SetScript([]byte(script))
	tx.SetGasLimit(100)

	latestBlock, err := client.GetLatestBlockHeader(ctx, true)
	if err != nil {
		return fmt.Errorf("failed to get latest block header: %w", err)
	}
	tx.SetReferenceBlockID(latestBlock.ID)

	tx.SetProposalKey(senderAccount.Address, senderAccount.Keys[0].Index, sequenceNumber)
	tx.SetPayer(senderAccount.Address)
	tx.AddAuthorizer(senderAccount.Address)

	cadencePubKeyHex, err := cadence.NewValue(publicKeyHex)
	cadenceKeysToAdd := cadence.NewInt(numOfKeysToAdd)

	tx.AddArgument(cadencePubKeyHex)

	tx.AddArgument(cadenceKeysToAdd)

	sigAlgo := crypto.ECDSA_P256
	hashAlgo := crypto.SHA3_256
	privateKey, err := crypto.DecodePrivateKeyHex(sigAlgo, senderPrivateKeyHex)
	if err != nil {
		return fmt.Errorf("failed to decode private key: %w", err)
	}

	signer, err := crypto.NewInMemorySigner(privateKey, hashAlgo)
	if err != nil {
		return fmt.Errorf("failed to create signer: %w", err)
	}

	err = tx.SignEnvelope(senderAccount.Address, senderAccount.Keys[0].Index, signer)
	if err != nil {
		return fmt.Errorf("failed to sign transaction envelope: %w", err)
	}

	err = client.SendTransaction(ctx, *tx)
	if err != nil {
		return fmt.Errorf("failed to send transaction: %w", err)
	}

	txHex := tx.ID().Hex()
	fmt.Printf("%d Keys generated, Hex: %s \n", numOfKeysToAdd, txHex)
    time.Sleep(10 * time.Second)
	return nil
}

func WaitForSeal(ctx context.Context, client *http.Client, txID flow.Identifier) {
	for {
		result, err := client.GetTransactionResult(ctx, txID)
		if err != nil {
			// log.Printf("Failed to get transaction result for %s: %v", txID, err)
			continue
		} else if result.Status == flow.TransactionStatusSealed {
			if result.Error != nil {
				log.Printf("Transaction %s sealed with error: %v", txID, result.Error)
                break
			} else {
                // successful transaction
                break
			}
            
		}
		// Sleep for a while before checking again.
		time.Sleep(1 * time.Second)
	}
}
