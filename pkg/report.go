package pkg

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"github.com/olekukonko/tablewriter"
)

type TemplateData struct {
	Label           string
	SendRate        float64
	SealRate        float64
	AvgSendLatency  string
	AvgSealLatency  string
	AvgLatency      string
	MinLatency      string
	MaxLatency      string
	TotalTx         int
	SuccessfulTx    int
	FailedTx        int
	Network         string
	Round           Round
}

func PrintStatsTable(stats TransactionStats) {
	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader([]string{"Metric", "Value"})

	// avgSendLatency := fmt.Sprintf("%.1f ms", stats.AverageSendLatency.Seconds()*1000)
	// avgSealLatency := fmt.Sprintf("%.1f ms", stats.AverageSealLatency.Seconds()*1000)
	averageLatency := fmt.Sprintf("%.1f ms", stats.AverageLatency.Seconds()*1000)
	minLatency := fmt.Sprintf("%.1f ms", stats.MinLatency.Seconds()*1000)
	maxLatency := fmt.Sprintf("%.1f ms", stats.MaxLatency.Seconds()*1000)

	table.SetAlignment(tablewriter.ALIGN_CENTER)

	table.Append([]string{"Send Rate (tps)", fmt.Sprintf("%.2f", stats.SendRate)})
	table.Append([]string{"Seal Rate (tps)", fmt.Sprintf("%.2f", stats.SealRate)})
	// table.Append([]string{"Average Send Latency", avgSendLatency})
	// table.Append([]string{"Average Seal Latency", avgSealLatency})
	table.Append([]string{"Minimum Network Latency", minLatency})
	table.Append([]string{"Maximum Network Latency", maxLatency})
	table.Append([]string{"Average Network Latency", averageLatency})
	table.Append([]string{"Total Transactions", fmt.Sprintf("%d", stats.TotalTx)})
	table.Append([]string{"Successful Transactions", fmt.Sprintf("%d", stats.SuccessfulTx)})
	table.Append([]string{"Failed Transactions", fmt.Sprintf("%d", stats.FailedTx)})

	table.Render()
}

func PrintSummary(allStats []TransactionStats, rounds []Round) {
    table := tablewriter.NewWriter(os.Stdout)
    table.SetHeader([]string{"Name", "Send Rate (tps)", "Seal Rate", "Max Latency", "Min Latency", "Avg Latency", "Successful Transactions", "Failed Transactions"})

    for i, stats := range allStats {
        // avgSendLatency := fmt.Sprintf("%.1f ms", stats.AverageSendLatency.Seconds()*1000)
        // avgSealLatency := fmt.Sprintf("%.1f ms", stats.AverageSealLatency.Seconds()*1000)
        averageLatency := fmt.Sprintf("%.1f ms", stats.AverageLatency.Seconds()*1000)
        minLatency := fmt.Sprintf("%.1f ms", stats.MinLatency.Seconds()*1000)
        maxLatency := fmt.Sprintf("%.1f ms", stats.MaxLatency.Seconds()*1000)

        table.Append([]string{
            rounds[i].Label,
            fmt.Sprintf("%.2f", stats.SendRate),
            fmt.Sprintf("%.2f", stats.SealRate),
            maxLatency,
            minLatency,
            averageLatency,
            fmt.Sprintf("%d", stats.SuccessfulTx),
            fmt.Sprintf("%d", stats.FailedTx),
        })
    }

    table.Render()
}


func GenerateReport(stats TransactionStats, round Round) {
    avgSendLatency := fmt.Sprintf("%.1f ms", stats.AverageSendLatency.Seconds()*1000)
    avgSealLatency := fmt.Sprintf("%.1f ms", stats.AverageSealLatency.Seconds()*1000)
    averageLatency := fmt.Sprintf("%.1f ms", stats.AverageLatency.Seconds()*1000)
    minLatency := fmt.Sprintf("%.1f ms", stats.MinLatency.Seconds()*1000)
    maxLatency := fmt.Sprintf("%.1f ms", stats.MaxLatency.Seconds()*1000)

    // Load configuration data
    data := TemplateData{
        Label:           round.Label,
        SendRate:        stats.SendRate,
        SealRate:        stats.SealRate,
        AvgSendLatency:  avgSendLatency,
        AvgSealLatency:  avgSealLatency,
        AvgLatency:      averageLatency,
        MinLatency:      minLatency,
        MaxLatency:      maxLatency,
        TotalTx:         stats.TotalTx,
        SuccessfulTx:    stats.SuccessfulTx,
        FailedTx:        stats.TotalTx - stats.SuccessfulTx,
        Network:         stats.Network,
        Round:           round,
    }

    filename := "report.html"
    file, err := os.Create(filename)
    if err != nil {
        panic(err)
    }
    defer file.Close()

	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
	}

	tmpl := template.Must(template.New("report").Funcs(funcMap).Parse(`
	<!DOCTYPE html>
	<html>
	
	<head>
	  <link href="https://fonts.googleapis.com/css2?family=Roboto:wght@300;500&display=swap" rel="stylesheet">
	  <style>
		body {
		  font-family: 'Roboto', sans-serif;
		  color: #333;
		  line-height: 1.5;
		  padding: 20px;
		}
	
		.container {
		  display: flex;
		  flex-wrap: wrap;
		}
	
		.config {
		  width: 100%;
		  border-right: 2px solid #ddd;
		  padding-right: 20px;
		  box-sizing: border-box;
		}
	
		.config .border-bottom {
		  border-bottom: 2px solid #ddd;
		  padding-bottom: 10px;
		  margin-bottom: 10px;
		}
	
		.config a {
		  color: #007BFF;
		  text-decoration: none;
		}
	
		.config a:hover {
		  text-decoration: underline;
		}
	
		.summary {
		  width: 100%;
		  padding-left: 20px;
		  box-sizing: border-box;
		}
	
		.summary table {
		  width: 100%;
		  border-collapse: collapse;
		}
	
		.summary th,
		.summary td {
		  border: 1px solid #ddd;
		  padding: 10px;
		}
	
		.summary tr:nth-child(even) {
		  background-color: #f2f2f2;
		}
	
		.summary th {
		  background-color: #f2f2f2;
		  color: #333;
		}
	
		.summary h3 {
		  margin-bottom: 5px;
		}
	
		.summary p {
		  margin-bottom: 20px;
		}
	
		@media (min-width: 600px) {
		  .config {
			width: 25%;
		  }
	
		  .summary {
			width: 75%;
		  }
		}
	  </style>
	</head>
	
	<body>
	  <div class="container">
		<div class="config">
		  <h2>Configuration</h2>
		  
		  <p><a href="#SendRate">SendRate</a>: {{.SendRate}}</p>
		  <p><a href="#SealRate">SealRate</a>: {{.SealRate}}</p>
		  <p><a href="#AvgSendLatency">Average Send Latency</a>: {{.AvgSendLatency}}</p>
		  <p><a href="#AvgSealLatency">Average Seal Latency</a>: {{.AvgSealLatency}}</p>
		  <p><a href="#MinLatency">Minimum Latency</a>: {{.MinLatency}}</p>
		  <p><a href="#MaxLatency">Maximum Latency</a>: {{.MaxLatency}}</p>
		  <p><a href="#AvgLatency">Average Latency</a>: {{.AvgLatency}}</p>
		  <p><a href="#TotalTx">Total Transactions</a>: {{.TotalTx}}</p>
		  <p><a href="#SuccessfulTx">Successful Transactions</a>: {{.SuccessfulTx}}</p>
		  <p><a href="#FailedTx">Failed Transactions</a>: {{.FailedTx}}</p>
		</div>
		<div class="summary">
		  <h2>Summary</h2>
		  <table>
			<tr>
			  <th>SendRate</th>
			  <td>{{.SendRate}}</td>
			</tr>
			<tr>
			  <th>SealRate</th>
			  <td>{{.SealRate}}</td>
			</tr>
			<tr>
			  <th>AvgSendLatency</th>
			  <td>{{.AvgSendLatency}}</td>
			</tr>
			<tr>
			  <th>AvgSealLatency</th>
			  <td>{{.AvgSealLatency}}</td>
			</tr>
			<tr>
			  <th>MinLatency</th>
			  <td>{{.MinLatency}}</td>
			</tr>
			<tr>
			  <th>MaxLatency</th>
			  <td>{{.MaxLatency}}</td>
			</tr>
			<tr>
			  <th>AvgLatency</th>
			  <td>{{.AvgLatency}}</td>
			</tr>
			<tr>
			  <th>TotalTx</th>
			  <td>{{.TotalTx}}</td>
			</tr>
			<tr>
			  <th>SuccessfulTx</th>
			  <td>{{.SuccessfulTx}}</td>
			</tr>
			<tr>
			  <th>FailedTx</th>
			  <td>{{.FailedTx}}</td>
			</tr>
		  </table>
		  <h3 id="SendRate">Send Rate</h3>
		  <p>This metric measures the rate at which transactions are sent. Result: {{.SendRate}}</p>
		  <h3 id="SealRate">Seal Rate</h3>
		  <p>This metric measures the rate at which transactions are sealed. Result: {{.SealRate}}</p>
		  <h3 id="AvgSendLatency">Average Send Latency</h3>
		  <p>This metric measures the average latency for sending transactions. Result: {{.AvgSendLatency}}</p>
		  <h3 id="AvgSealLatency">Average Seal Latency</h3>
		  <p>This metric measures the rate at which transactions are sent. Result: {{.SendRate}}</p>
		  <h3 id="MinLatency">Minimum Latency</h3>
		  <p>This metric measures the rate at which transactions are sealed. Result: {{.SealRate}}</p>
		  <h3 id="MaxLatency">Maximum Latency</h3>
		  <p>This metric measures the average latency for sending transactions. Result: {{.AvgSendLatency}}</p>
		  <h3 id="AvgLatency">Average Latency</h3>
		  <p>This metric measures the rate at which transactions are sent. Result: {{.SendRate}}</p>
		  <h3 id="TotalTx">Total Transactions</h3>
		  <p>This metric measures the rate at which transactions are sealed. Result: {{.SealRate}}</p>
		  <h3 id="SuccessfulTx">Successful Transactions</h3>
		  <p>This metric measures the average latency for sending transactions. Result: {{.AvgSendLatency}}</p>
		  <h3 id="FailedTx">Failed Transactions</h3>
		  <p>This metric measures the average latency for sending transactions. Result: {{.AvgSendLatency}}</p>
		</div>
	  </div>
	</body>
	
	</html>
	`))

	if err := tmpl.Execute(file, data); err != nil {
		panic(err)
	}

	absPath, err := filepath.Abs(filename)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Benchmark Complete!\n")
	fmt.Printf("For more information, check out report at file://%s\n", absPath)
}

// removed configuration code from html
// <p>Network: {{.Config.Network}}</p>
		//   <p>NumTransactions: {{.Config.NumTransactions}}</p>
		//   <p>TransactionsPerSecond: {{.Config.TransactionsPerSecond}}</p>
		//   <p>RecipientAddress: {{.Config.RecipientAddress}}</p>
		//   <p class="border-bottom">SenderAddress: {{.Config.SenderAddress}}</p>