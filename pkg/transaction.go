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
	"io/ioutil"
	"strconv"
)




func createCadenceValue(t string, value string) (cadence.Value, error) {
	switch t {
	case "Int8":
		integer, err := strconv.ParseInt(value, 10, 8)
		if err != nil {
			return nil, fmt.Errorf("error converting string to int8: %w", err)
		}
		intValue := cadence.NewInt8(int8(integer))
		return intValue, nil
	case "UInt8":
		uinteger, err := strconv.ParseUint(value, 10, 8)
		if err != nil {
			return nil, fmt.Errorf("error converting string to uint8: %w", err)
		}
		uintValue := cadence.NewUInt8(uint8(uinteger))
		return uintValue, nil
	case "Int16":
		integer, err := strconv.ParseInt(value, 10, 16)
		if err != nil {
			return nil, fmt.Errorf("error converting string to int16: %w", err)
		}
		intValue := cadence.NewInt16(int16(integer))
		return intValue, nil
	case "UInt16":
		uinteger, err := strconv.ParseUint(value, 10, 16)
		if err != nil {
			return nil, fmt.Errorf("error converting string to uint16: %w", err)
		}
		uintValue := cadence.NewUInt16(uint16(uinteger))
		return uintValue, nil
	case "Int32":
		integer, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("error converting string to int32: %w", err)
		}
		intValue := cadence.NewInt32(int32(integer))
		return intValue, nil
	case "UInt32":
		uinteger, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("error converting string to uint32: %w", err)
		}
		uintValue := cadence.NewUInt32(uint32(uinteger))
		return uintValue, nil
	case "Int64":
		integer, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error converting string to int64: %w", err)
		}
		intValue := cadence.NewInt64(integer)
		return intValue, nil
	case "UInt64":
		uinteger, err := strconv.ParseUint(value, 10, 64)
        if err != nil {
			return nil, fmt.Errorf("error converting string to uint64: %w", err)
        }
		uintValue := cadence.NewUInt64(uinteger)
		return uintValue, nil

	case "Address":
		addressValue := cadence.BytesToAddress(flow.HexToAddress(value).Bytes())
		return addressValue, nil
	case "String":
		stringValue := cadence.String(value)
		return stringValue, nil
	case "Bool":
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return nil, fmt.Errorf("error loading the boolean, make sure %s is a boolean", value)
		}
		boolCadenceValue := cadence.Bool(boolValue)
		return boolCadenceValue, nil
	case "Fix64":
		fix64Value, err := cadence.NewFix64(value)
		if err != nil {
			return nil, err
		}
		return fix64Value, nil
	case "UFix64":
		ufix64Value, err := cadence.NewUFix64(value)
		if err != nil {
			return nil, err
		}
		return ufix64Value, nil
	}

	return nil, fmt.Errorf("unsupported type: %s", t)
}

func SendTransaction(ctx context.Context, client *http.Client, senderAccount *flow.Account, sequenceNumber uint64, keyID int, transaction Transaction) (time.Duration, time.Duration, string, flow.Identifier, time.Time, bool) {
    tx := flow.NewTransaction()
    transactionsss, err := LoadTransactionConfig()
        if err != nil {
            log.Fatalf("Failed to load transaction configuration: %v", err)
        }
    var senderPrivateKeyHex = transactionsss.Payer.PrivateKey

	script, err := ioutil.ReadFile(transaction.ScriptPath)
	if err != nil {
		panic("script path wrong!")
    }

    tx.SetScript([]byte(script))
    tx.SetGasLimit(transaction.GasLimit)
    

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
        return 0, 0, "", flow.Identifier{}, time.Time{}, false
    }
    
    tx.SetReferenceBlockID(latestBlock.ID)

    tx.SetProposalKey(senderAccount.Address, senderAccount.Keys[keyID].Index, sequenceNumber)
    tx.SetPayer(senderAccount.Address)
    tx.AddAuthorizer(senderAccount.Address)

	for _, arg := range transaction.ScriptArguments {
		argType := arg.Type
		argValue := arg.Value
		
		cadenceValue, err := createCadenceValue(argType, argValue)
			if err != nil {
					fmt.Println("Error creating Cadence value:", err)
					continue
			}

		tx.AddArgument(cadenceValue)
    }

    sigAlgo := crypto.ECDSA_P256
    hashAlgo := crypto.SHA3_256
    privateKey, err := crypto.DecodePrivateKeyHex(sigAlgo, senderPrivateKeyHex)
    if err != nil {
        fmt.Printf("Error decoding private key: %v\n", err)
        return 0, 0, "", flow.Identifier{}, time.Time{}, false
    }

    signer, err := crypto.NewInMemorySigner(privateKey, hashAlgo)
    if err != nil {
        fmt.Printf("Error creating signer: %v\n", err)
        return 0, 0, "", flow.Identifier{}, time.Time{}, false
    }

    if err = tx.SignEnvelope(senderAccount.Address, senderAccount.Keys[0].Index, signer); err != nil {
        fmt.Printf("Error signing envelope: %v\n", err)
        return 0, 0, "", flow.Identifier{}, time.Time{}, false
    }
    if keyID != 0 {
        if err = tx.SignEnvelope(senderAccount.Address, senderAccount.Keys[keyID].Index, signer); err != nil {
            fmt.Printf("Error signing envelope: %v\n", err)
            return 0, 0, "", flow.Identifier{}, time.Time{}, false
        }
    }

    startTime := time.Now()
    if err = client.SendTransaction(ctx, *tx); err != nil {
        fmt.Printf("Error sending transaction: %v\n", err)
        return 0, 0, "", flow.Identifier{}, time.Time{}, false
    }
    txEndTime := time.Now()
    txHex := tx.ID().Hex()

    WaitForSeal(ctx, client, tx.ID())
    sealEndTime := time.Now()

    txLatency := txEndTime.Sub(startTime)
    sealLatency := sealEndTime.Sub(startTime)

    return txLatency, sealLatency, txHex, tx.ID(), txEndTime, true
}

func AddKeys(ctx context.Context, client *http.Client, senderAccount *flow.Account, sequenceNumber uint64, numOfKeysToAdd int) error {
	tx := flow.NewTransaction()
    transaction, err := LoadTransactionConfig()
	if err != nil {
		log.Fatalf("Failed to load transaction configuration: %v", err)
	}
	var senderPrivateKeyHex = transaction.Payer.PrivateKey
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
	tx.SetGasLimit(9999)

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