package pkg

import (
	"fmt"
	"log"
	"path/filepath"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)


type Benchmark struct {
	NumOfTransactions int `yaml:"numOfTransactions"`
	Tps               int `yaml:"tps"`
	Network           string `yaml:"network"`
	KeyGeneration struct {
		Enabled       bool `yaml:"enabled"`
		KeysGenerated int  `yaml:"keysGenerated"`
	} `yaml:"keyGeneration"`
	Concurrency struct {
		Enabled         bool `yaml:"enabled"`
		MaxConcurrency  int  `yaml:"maxConcurrency"`
		Buffer          bool `yaml:"buffer"`
		KeyGeneration struct {
			Enabled       bool `yaml:"enabled"`
			KeysGenerated int  `yaml:"keysGenerated"`
		} `yaml:"keyGeneration"`
	} `yaml:"concurrency"`
	ReportingAndOutput struct {
		GenerateReport  bool `yaml:"generateReport"`
		PrintStatsTable bool `yaml:"printStatsTable"`
	} `yaml:"reportingAndOutput"`
	RetrySettings struct {
		Enabled          bool `yaml:"enabled"`
		MaxRetryAttempts int  `yaml:"maxRetryAttempts"`
		RetryInterval    int  `yaml:"retryInterval"`
	} `yaml:"retrySettings"`
	LoggingAndErrorHandling struct {
		EnableLogging     bool `yaml:"enableLogging"`
		EnableErrorHandling bool `yaml:"enableErrorHandling"`
	} `yaml:"loggingAndErrorHandling"`
}

type Transaction struct {
	ScriptPath       string `yaml:"scriptPath"`
	GasLimit         int    `yaml:"gasLimit"`
	ScriptArguments  struct {
		Amount struct {
			Type  string `yaml:"type"`
			Value string `yaml:"value"`
		} `yaml:"amount"`
		Recipient struct {
			Type  string `yaml:"type"`
			Value string `yaml:"value"`
		} `yaml:"recipient"`
	} `yaml:"scriptArguments"`
	Payer struct {
		UseSameAccount bool   `yaml:"useSameAccount"`
		Address        string `yaml:"address"`
		PrivateKey     string `yaml:"privateKey"`
	} `yaml:"payer"`
	Proposer struct {
		UseSameAccount bool   `yaml:"useSameAccount"`
		Address        string `yaml:"address"`
		PrivateKey     string `yaml:"privateKey"`
	} `yaml:"proposer"`
	Authorizer struct {
		UseSameAccount bool   `yaml:"useSameAccount"`
		Address        string `yaml:"address"`
		PrivateKey     string `yaml:"privateKey"`
	} `yaml:"authorizer"`
}

func LoadBenchmarkConfig() (*Benchmark, error) {
	absPath, err := filepath.Abs("benchmarkConfig.yaml")
	if err != nil {
		log.Fatalf("Failed to get absolute path: %v", err)
	}

	data, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML file: %w", err)
	}

	var benchmark Benchmark
	err = yaml.Unmarshal(data, &benchmark)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return &benchmark, nil
}


func LoadTransactionConfig() (*Transaction, error) {
	absPath, err := filepath.Abs("transactionConfig.yaml")
	if err != nil {
		log.Fatalf("Failed to get absolute path: %v", err)
	}
	data, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML file: %w", err)
	}

	var transaction Transaction
	err = yaml.Unmarshal(data, &transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return &transaction, nil
}
