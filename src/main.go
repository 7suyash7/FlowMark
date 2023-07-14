package main

import (
	"time"
	"context"
	"math"
	"strconv"
	"fmt"
	"log"
	"os"
	"flag"
	"strings"
	. "github.com/7suyash7/FlowMark/pkg"

	"github.com/onflow/flow-go-sdk/access/http"
	"github.com/onflow/flow-go-sdk"
	"github.com/joho/godotenv"
)

func main() {
	args := os.Args[1:]

	// Check if there are flags specified in front of the binary
	var hasFlags bool
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			hasFlags = true
			break
		}
	}

	if len(args) > 0 && args[0] == "start" {
		runBenchmark()
	} else if len(args) > 0 && args[0] == "help" {
		displayManual()
	} else if len(os.Args) > 1 && os.Args[1] == "config" {
		displayConfiguration()
		return
	} else if hasFlags {
		checkFlags()
	} else {
		displayHelp()
	}
}

func displayHelp() {
	fmt.Println("Usage: ./binary start")
	fmt.Println("Please type 'help' to learn the available commands.")
}

func displayManual() {
	fmt.Println("=== BENCHMARK MANUAL ===")
	fmt.Println("This is the manual for the benchmark tool.")
	fmt.Println("To run the benchmark, use the following command:")
	fmt.Println("./binary start")
	fmt.Println()
	fmt.Println("Command-line options:")
	fmt.Println("start                  - Run the benchmark")
	fmt.Println("help                   - Show this manual")
	fmt.Println("config                 - Display the configuration")
	fmt.Println("Options for benchmark:")
	fmt.Println("--sender-address       - Set the sender address")
	fmt.Println("--receiver-address     - Set the receiver address")
	fmt.Println("--numTransaction       - Set the number of transactions")
	fmt.Println("--network              - Set the network (emulator, testnet, mainnet)")
	fmt.Println("--sender-priv-address  - Set the sender private key")
	fmt.Println()
	fmt.Println("Example usage:")
	fmt.Println("./binary start --sender-address bdb89318be61241e --receiver-address 1ba7234d25ebb0c0 --numTransaction 10 --network testnet --sender-priv-address dd4ccf9ef501eee0ee0690550342e7c09e0e9d997d926f7a959e6f3b05b1c81a")
}

func displayConfiguration() {
	envFile := ".env"

	env, err := godotenv.Read(envFile)
	if err != nil {
		log.Fatalf("Error reading .env file: %v", err)
	}

	fmt.Println("Configuration:")
	for key, value := range env {
		fmt.Printf("%s=%s\n", key, value)
	}
}



func checkFlags() {
	envFile := ".env"

	senderAddressFlag := flag.String("sender-address", "", "Sender address")
	receiverAddressFlag := flag.String("recipient-address", "", "Receiver address")
	numTransactionsFlag := flag.Int("numTransaction", 0, "Number of transactions")
	networkFlag := flag.String("network", "", "Network (emulator, testnet, mainnet)")
	senderPrivateKeyFlag := flag.String("sender-priv-address", "", "Sender private key")

	flag.Parse()

	env, err := godotenv.Read(envFile)
	if err != nil {
		log.Fatalf("Error reading .env file: %v", err)
	}

	// Check if each field is empty in the .env file and remind the user if it's empty
	if env["SENDER_ADDRESS"] == "" && *senderAddressFlag == "" {
		fmt.Println("WARNING: SENDER_ADDRESS is empty in .env file")
	}
	if env["SENDER_PRIVATE_KEY"] == "" && *senderPrivateKeyFlag == "" {
		fmt.Println("WARNING: SENDER_PRIVATE_KEY is empty in .env file")
	}
	if env["RECIPIENT_ADDRESS"] == "" && *receiverAddressFlag == "" {
		fmt.Println("WARNING: RECIPIENT_ADDRESS is empty in .env file")
	}
	if env["NO_OF_TRANSACTION"] == "" && *numTransactionsFlag == 0 {
		fmt.Println("WARNING: NO_OF_TRANSACTION is empty in .env file")
	}
	if env["NETWORK"] == "" && *networkFlag == "" {
		fmt.Println("WARNING: NETWORK is empty in .env file")
	}

	// Update the environment variables if provided via flags
	if *senderAddressFlag != "" {
		env["SENDER_ADDRESS"] = *senderAddressFlag
	}
	if *senderPrivateKeyFlag != "" {
		env["SENDER_PRIVATE_KEY"] = *senderPrivateKeyFlag
	}
	if *receiverAddressFlag != "" {
		env["RECIPIENT_ADDRESS"] = *receiverAddressFlag
	}
	if *numTransactionsFlag != 0 {
		env["NO_OF_TRANSACTION"] = strconv.Itoa(*numTransactionsFlag)
	}
	if *networkFlag != "" {
		env["NETWORK"] = *networkFlag
	}

	err = godotenv.Write(env, envFile)
	if err != nil {
		log.Fatalf("Error writing .env file: %v", err)
	}
}

// NOTE - USE THIS TO LIKE SETUP WHEN SOMEONE FIRST STARTS
func promptField(fieldName string) {
	fmt.Printf("Please specify the value for %s: ", fieldName)
	var value string
	_, err := fmt.Scanln(&value)
	if err != nil {
		log.Fatalf("Error reading user input: %v", err)
	}
	os.Setenv(fieldName, value)
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
			} else {
				log.Printf("Transaction %s sealed successfully", txID)
			}
			break
		}

		// Sleep for a while before checking again.
		time.Sleep(1 * time.Second)
	}
}


func runBenchmark() {
	numTransactionsStr := LoadEnvVar("NO_OF_TRANSACTION")
	numTransactions, err := strconv.Atoi(numTransactionsStr)
	if err != nil {
    	log.Fatalf("Error converting NO_OF_TRANSACTION to int: %v", err)
	}
	startTime := time.Now()
	var totalSendLatency time.Duration
	var totalSealLatency time.Duration
	maxLatency := time.Duration(0)
	minLatency := time.Duration(math.MaxInt64)

	ctx := context.Background()
	network := LoadEnvVar("NETWORK")
	var client *http.Client

	switch network {
		case "emulator":
			client, err = InitializeClient(http.EmulatorHost)
		case "testnet":
			client, err = InitializeClient(http.TestnetHost)
		case "mainnet":
			client, err = InitializeClient(http.MainnetHost)
		default:
			panic("No Network Selected! Select mainnet, testnet, or emulator in .env under the network variable")
	}

	if err != nil {
		panic(err)
	}

	var senderAddressHex = LoadEnvVar("SENDER_ADDRESS")
	senderAddress := flow.HexToAddress(senderAddressHex)

	stats := NewTransactionStats()
    transactionIDs := make([]flow.Identifier, 0, numTransactions)	

	for i := 0; i < numTransactions; i++ {
		// Get the latest account info for this address
		senderAccount, err := client.GetAccountAtLatestBlock(ctx, senderAddress)
		if err != nil {
		panic("failed to fetch proposer account")
		}

		// Get the latest sequence number for this key
		sequenceNumber := senderAccount.Keys[0].SequenceNumber

		sendStartTime := time.Now()
		latency, txHex, txID := SendTransaction(ctx, client, senderAccount, sequenceNumber)
		sendEndTime := time.Now()
		totalSendLatency += sendEndTime.Sub(sendStartTime)
		
		WaitForSeal(ctx, client, txID)
		sealEndTime := time.Now()
		totalSealLatency += sealEndTime.Sub(sendStartTime)

		sequenceNumber++
        transactionIDs = append(transactionIDs, txID)
		
		stats = UpdateStats(stats, latency, txHex)

		if latency > maxLatency {
			maxLatency = latency
		}

		if latency < minLatency {
			minLatency = latency
		}
	}

	endTime := time.Now()

	time.Sleep(10 * time.Second)
	fmt.Printf("Generating results...\n")

    successfulTransactions := 0
	for _, txID := range transactionIDs {
		result, err := client.GetTransactionResult(ctx, txID)
		if err != nil {
			log.Printf("Failed to get transaction result for %s: %v", txID, err)
			continue
		}

		if result.Status == flow.TransactionStatusSealed && result.Error == nil {
			successfulTransactions++
		}
	}

	stats = FinalizeStats(stats, startTime, endTime, totalSendLatency, totalSealLatency, minLatency, maxLatency, numTransactions, successfulTransactions, network)

	PrintStatsTable(stats)
	GenerateReport(stats)
}
