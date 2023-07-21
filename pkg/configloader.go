package pkg

import (
	"fmt"
	"log"
	"path/filepath"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type RateControl struct {
	TxNumber int `yaml:"txNumber"`
	Tps      int `yaml:"tps"`
}

type Round struct {
	Label        string      `yaml:"label"`
	Description  string      `yaml:"description"`
	RateControl  RateControl `yaml:"rateControl"`
}

type Workers struct {
	Number int `yaml:"number"`
}

type Test struct {
	Network     string   `yaml:"network"`
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Workers     Workers  `yaml:"workers"`
	Rounds      []Round  `yaml:"rounds"`
}

type Benchmark struct {
	Test Test `yaml:"test"`
}

type Transaction struct {
	ScriptPath       string `yaml:"scriptPath"`
	GasLimit         uint64 `yaml:"gasLimit"`
	ScriptArguments  []struct {
		Name  string `yaml:"name"`
		Type  string `yaml:"type"`
		Value string `yaml:"value"`
	} `yaml:"scriptArguments"`
	Payer struct {
		UseSameAccount bool   `yaml:"useSameAccount"`
		Address        string `yaml:"address"`
		PrivateKey     string `yaml:"privateKey"`
	} `yaml:"payer"`
	Proposer struct {
		UseSameAccount   bool   `yaml:"useSameAccount"`
		Address          string `yaml:"address"`
		PrivateKey       string `yaml:"privateKey"`
		ProposerKeyIndex int    `yaml:"proposerKeyIndex"`
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
