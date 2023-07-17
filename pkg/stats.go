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
	MinSealLatency        time.Duration
	MaxSealLatency        time.Duration
	benchmarkTime		  time.Duration
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

func UpdateStats(stats TransactionStats, txHex string) TransactionStats {
	stats.TxHexes = append(stats.TxHexes, txHex)
	return stats
}

func FinalizeStats(stats TransactionStats, startTime time.Time, endTime time.Time, totalSendLatency time.Duration, totalSealLatency time.Duration, minLatency time.Duration, maxLatency time.Duration, numTransactions int, successfulTransactions int, Network string) TransactionStats {
	// duration := endTime.Sub(startTime)
	sealRate := float64(numTransactions) / totalSealLatency.Seconds()
	avgSendLatency := totalSendLatency / time.Duration(successfulTransactions)
	avgSealLatency := totalSealLatency / time.Duration(successfulTransactions)
	benchmarkTime := endTime.Sub(startTime)
	sendRate := float64(numTransactions) / benchmarkTime.Seconds()
	averageLatency := (totalSendLatency + totalSealLatency) / time.Duration(successfulTransactions)

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
	stats.benchmarkTime = benchmarkTime

	return stats
}
