package pkg

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/olekukonko/tablewriter"
)

type TemplateData struct {
	SendRate       float64
	SealRate       float64
	AvgSendLatency string
	AvgSealLatency string
	AvgLatency 	   string
	MinLatency     string
	MaxLatency     string
	Throughput     float64
	TxHexes        []string
	Network        string
    Config         Configuration
	TotalTx        int
	SuccessfulTx   int
	FailedTx       int
}

type Configuration struct {
	Network            string
	NumTransactions    string
	RecipientAddress   string
	SenderAddress      string
	SenderPrivateKey   string
}



func PrintStatsTable(stats TransactionStats) {
	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader([]string{"Metric", "Value"})

	avgSendLatency := fmt.Sprintf("%.1f ms", stats.AverageSendLatency.Seconds()*1000)
	avgSealLatency := fmt.Sprintf("%.1f ms", stats.AverageSealLatency.Seconds()*1000)
	averageLatency := fmt.Sprintf("%.1f ms", stats.AverageLatency.Seconds()*1000)
	minLatency := fmt.Sprintf("%.1f ms", stats.MinLatency.Seconds()*1000)
	maxLatency := fmt.Sprintf("%.1f ms", stats.MaxLatency.Seconds()*1000)

	table.SetAlignment(tablewriter.ALIGN_CENTER)

	table.Append([]string{"Send Rate (tps)", fmt.Sprintf("%.2f", stats.SendRate)})
	table.Append([]string{"Seal Rate (tps)", fmt.Sprintf("%.2f", stats.SealRate)})
	table.Append([]string{"Average Send Latency", avgSendLatency})
	table.Append([]string{"Average Seal Latency", avgSealLatency})
	table.Append([]string{"Minimum Network Latency", minLatency})
	table.Append([]string{"Maximum Network Latency", maxLatency})
	table.Append([]string{"Average Network Latency", averageLatency})
	table.Append([]string{"Total Transactions", fmt.Sprintf("%d", stats.TotalTx)})
	table.Append([]string{"Successful Transactions", fmt.Sprintf("%d", stats.SuccessfulTx)})
	table.Append([]string{"Failed Transactions", fmt.Sprintf("%d", stats.FailedTx)})

	table.Render()
}


func GenerateReport(stats TransactionStats) {
	avgSendLatency := fmt.Sprintf("%.1f ms", stats.AverageSendLatency.Seconds()*1000)
	avgSealLatency := fmt.Sprintf("%.1f ms", stats.AverageSealLatency.Seconds()*1000)
	averageLatency := fmt.Sprintf("%.1f ms", stats.AverageLatency.Seconds()*1000)
	minLatency := fmt.Sprintf("%.1f ms", stats.MinLatency.Seconds()*1000)
	maxLatency := fmt.Sprintf("%.1f ms", stats.MaxLatency.Seconds()*1000)

	// Load configuration data
	config := Configuration{
		Network:            LoadEnvVar("NETWORK"),
		NumTransactions:    LoadEnvVar("NO_OF_TRANSACTION"),
		RecipientAddress:   LoadEnvVar("RECIPIENT_ADDRESS"),
		SenderAddress:      LoadEnvVar("SENDER_ADDRESS"),
		SenderPrivateKey:   LoadEnvVar("SENDER_PRIVATE_KEY"),
	}

	// Add config to template data
	data := TemplateData{
		SendRate:     stats.SendRate,
		SealRate:     stats.SealRate,
		AvgSendLatency: avgSendLatency,
		AvgSealLatency: avgSealLatency,
		AvgLatency:   averageLatency,
		MinLatency:   minLatency,
		MaxLatency:   maxLatency,
		TxHexes:      stats.TxHexes,
		TotalTx:	  stats.TotalTx,
		SuccessfulTx: stats.SuccessfulTx,
		FailedTx: 	  stats.TotalTx - stats.SuccessfulTx,
		Network:      stats.Network,
		Config:       config,
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
    <title>Flowmark Blockchain Testnet Benchmark Report</title>
		<style>
			body {
				font-family: Arial, sans-serif;
				background-color: #f4f4f4;
				color: #333;
				padding: 10px;
            display: flex; /* Added to use Flexbox */
			}

        h1,
        h2 {
				color: #333;
			}

        // table {
        //     width: 100%;
        //     margin-top: 20px;
        //     border-collapse: collapse;
        // }

        // th,
        // td {
        //     padding: 10px;
        //     border: 1px solid #ddd;
        //     text-align: left;
        // }
			table {
				width: 100%;
				margin-top: 20px;
				border-collapse: collapse;
			table-layout: fixed; /* add this line */
			word-wrap: break-word; /* add this line */
			}
		
		table, th, td {
			overflow: auto; /* add this line */
				padding: 10px;
				border: 1px solid #ddd;
				text-align: left;
			}

			th {
				background-color: #4CAF50;
				color: white;
			}

			tr:nth-child(even) {
				background-color: #f2f2f2;
			}

			a {
				color: #4CAF50;
			}

			a:hover {
				color: #45a049;
			}

        .config-container {
            width: 20%; /* This means that the config-container will take up 20% of the body width */
            margin-right: 10px;
            padding: 20px;
				background-color: white;
            border-radius: 5px;
            box-shadow: 0px 0px 10px 0px rgba(0, 0, 0, 0.1);
        }

        .main-container {
            flex-grow: 1; /* This means that the main-container will take up the remaining space */
				padding: 20px;
            background-color: white;
				border-radius: 5px;
            box-shadow: 0px 0px 10px 0px rgba(0, 0, 0, 0.1);
			}

			.section {
				border: 1px solid #ddd;
				padding: 10px;
				margin: 10px 0;
				border-radius: 5px;
			}

        .project-name {
            font-size: 24px;
            font-weight: bold;
            margin-bottom: 10px;
        }
		</style>
	</head>

	<body>
		<div class="main-container">
			<div class="column">
		<div class="section">
			<h2>Summary Table</h2>
			<table>
						<tr>
							<th>Metric</th>
							<th>Value</th>
						</tr>
						<tr>
							<td>Send Rate</td>
							<td>{{.SendRate}}</td>
						</tr>
						<tr>
							<td>Seal Rate (tps)</td>
							<td>{{.SealRate}}</td>
						</tr>
						<tr>
							<td>Average Send Latency</td>
							<td>{{.AvgSendLatency}}</td>
						</tr>
						<tr>
							<td>Average Seal Latency</td>
							<td>{{.AvgSealLatency}}</td>
						</tr>
						<tr>
							<td>Minimum Network Latency</td>
							<td>{{.MinLatency}}</td>
						</tr>
						<tr>
							<td>Maximum Network Latency</td>
							<td>{{.MaxLatency}}</td>
						</tr>
						<tr>
							<td>Average Network Latency</td>
							<td>{{.AvgLatency}}</td>
						</tr>
						<tr>
							<td>Total Transactions</td>
							<td>{{.TotalTx}}</td>
						</tr>
						<tr>
							<td>Successful Transactions</td>
							<td>{{.SuccessfulTx}}</td>
						</tr>
						<tr>
							<td>Failed Transactions</td>
							<td>{{.FailedTx}}</td>
						</tr>
			</table>
		</div>
		<div class="container">
					<h1 class="project-name">Flowmark Blockchain Testnet Benchmark Report</h1>
					<div class="section">
						<h2>Seal Rate</h2>
						<p>This represents the rate at which transactions are being sealed (or confirmed) on the blockchain network. It's a measure of how quickly transactions are being processed by the network.</p>
						<p>Result: {{.SealRate}}</p>
					</div>
			<div class="section">
				<h2>Send Rate</h2>
						<p>This represents the rate at which transactions are being sent to the blockchain network. It's a measure of how quickly transactions are being dispatched from the originating point.</p>
				<p>Result: {{.SendRate}}</p>
			</div>
			<div class="section">
						<h2>Average Send Latency</h2>
						<p>This represents the average amount of time it takes for a transaction to be sent to the blockchain network. It gives an idea of how long it takes for transactions to leave the originating point.</p>
						<p>Result: {{.AvgSendLatency}}</p>
			</div>
			<div class="section">
						<h2>Average Seal Latency</h2>
						<p>This represents the average amount of time it takes for a transaction to be sealed (or confirmed) on the blockchain network. It gives an idea of how long it takes for transactions to be processed by the network.</p>
						<p>Result: {{.AvgSealLatency}}</p>
			</div>			
			<div class="section">
						<h2>Minimum Network Latency</h2>
						<p>This refers to the fastest recorded time that it took for a transaction to be processed by the network. It represents the shortest observed duration between when a transaction was sent and when it was confirmed on the blockchain. This measure can give an indication of the best-case performance of the network.</p>
				<p>Result: {{.MinLatency}}</p>
			</div>
			<div class="section">
						<h2>Maximum Network Latency</h2>
						<p>This refers to the longest recorded time that it took for a transaction to be processed by the network. It represents the longest observed duration between when a transaction was sent and when it was confirmed on the blockchain. This measure can give an indication of the worst-case performance of the network. It's particularly important for understanding the upper bounds of delay that transactions might experience under current network conditions.</p>
				<p>Result: {{.MaxLatency}}</p>
			</div>
			<div class="section">
						<h2>Average Network Latency</h2>
						<p>This represents the average amount of time it takes for a transaction to be processed, from the moment it's sent to when it's confirmed on the blockchain network. It's a measure of overall network performance.</p>
						<p>Result: {{.AvgLatency}}</p>
					</div>
					<div class="section">
				<h2>Transactions</h2>
				{{if eq .Network "testnet"}}
					{{range $index, $hex := .TxHexes}}
						<p>{{add $index 1}}}. <a href="https://testnet.flowscan.org/transaction/{{$hex}}">{{$hex}}</a></p>
					{{end}}
				{{else if eq .Network "mainnet"}}
					{{range $index, $hex := .TxHexes}}
						<p>{{add $index 1}}}. <a href="https://flowscan.org/transaction/{{$hex}}">{{$hex}}</a></p>
					{{end}}
				{{end}}
			</div>
		</div>
			</div>
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
