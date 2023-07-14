package pkg

import (
	"time"
)

type TransactionStats struct {
	SealRate       float64
	AverageLatency time.Duration
	MinLatency     time.Duration
	MaxLatency     time.Duration
	Throughput     float64
	TxHexes        []string
	TotalTx        int
	SuccessfulTx   int
	FailedTx       int
	Network        string
}

func NewTransactionStats() TransactionStats {
	return TransactionStats{
		TxHexes: make([]string, 0),
	}
}

func UpdateStats(stats TransactionStats, latency time.Duration, txHex string) TransactionStats {
	stats.TxHexes = append(stats.TxHexes, txHex)
	return stats
}

func FinalizeStats(stats TransactionStats, startTime time.Time, endTime time.Time, totalLatency time.Duration, minLatency time.Duration, maxLatency time.Duration, numTransactions int, successfulTransactions int, Network string) TransactionStats {
	duration := endTime.Sub(startTime)
	sealRate := float64(numTransactions) / duration.Seconds()
	avgLatency := totalLatency / time.Duration(numTransactions)
	throughput := 1 / avgLatency.Seconds()

	stats.SealRate = sealRate
	stats.AverageLatency = avgLatency
	stats.MinLatency = minLatency
	stats.MaxLatency = maxLatency
	stats.Throughput = throughput
	stats.TotalTx = numTransactions
	stats.SuccessfulTx = successfulTransactions
	stats.FailedTx = numTransactions - successfulTransactions
	stats.Network = Network

	return stats
}