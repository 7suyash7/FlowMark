package pkg

import (
	"time"
)

type TransactionStats struct {
	SendRate          float64
	SealRate          float64
	AverageSendLatency time.Duration
	AverageSealLatency time.Duration
	MinLatency        time.Duration
	MaxLatency        time.Duration
	AverageLatency	  time.Duration
	SendThroughput    float64
	SealThroughput    float64
	TxHexes           []string
	TotalTx           int
	SuccessfulTx      int
	FailedTx          int
	Network           string
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

func FinalizeStats(stats TransactionStats, startTime time.Time, endTime time.Time, totalSendLatency time.Duration, totalSealLatency time.Duration, minLatency time.Duration, maxLatency time.Duration, numTransactions int, successfulTransactions int, Network string) TransactionStats {
	// duration := endTime.Sub(startTime)
	sendRate := float64(numTransactions) / totalSendLatency.Seconds()
	sealRate := float64(numTransactions) / totalSealLatency.Seconds()
	avgSendLatency := totalSendLatency / time.Duration(numTransactions)
	avgSealLatency := totalSealLatency / time.Duration(numTransactions)
	averageLatency := (minLatency + maxLatency / 2)
	// averageLatency := (totalSendLatency + totalSealLatency) / time.Duration(numTransactions)

	stats.SendRate = sendRate
	stats.SealRate = sealRate
	stats.AverageSendLatency = avgSendLatency
	stats.AverageSealLatency = avgSealLatency
	stats.MinLatency = minLatency
	stats.MaxLatency = maxLatency
	stats.TotalTx = numTransactions
	stats.SuccessfulTx = successfulTransactions
	stats.FailedTx = numTransactions - successfulTransactions
	stats.Network = Network
	stats.AverageLatency = averageLatency

	return stats
}
