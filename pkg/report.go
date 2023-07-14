package pkg

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/olekukonko/tablewriter"
)

type TemplateData struct {
	SealRate       float64
	AvgLatency     string
	MinLatency     string
	MaxLatency     string
	Throughput     float64
	TxHexes        []string
	Network        string
}

func PrintStatsTable(stats TransactionStats) {
	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader([]string{"Metric", "Value"})

	avgLatency := fmt.Sprintf("%.1f ms", stats.AverageLatency.Seconds()*1000)
	minLatency := fmt.Sprintf("%.1f ms", stats.MinLatency.Seconds()*1000)
	maxLatency := fmt.Sprintf("%.1f ms", stats.MaxLatency.Seconds()*1000)

	table.SetAlignment(tablewriter.ALIGN_CENTER)

	table.Append([]string{"Seal Rate (tps)", fmt.Sprintf("%.2f", stats.SealRate)})
	table.Append([]string{"Average Latency", avgLatency})
	table.Append([]string{"Minimum Latency", minLatency})
	table.Append([]string{"Maximum Latency", maxLatency})
	table.Append([]string{"Throughput (tps)", fmt.Sprintf("%.2f", stats.Throughput)})
	table.Append([]string{"Total Transactions", fmt.Sprintf("%d", stats.TotalTx)})
	table.Append([]string{"Successful Transactions", fmt.Sprintf("%d", stats.SuccessfulTx)})
	table.Append([]string{"Failed Transactions", fmt.Sprintf("%d", stats.FailedTx)})

	table.Render()
}

func GenerateReport(stats TransactionStats) {
	avgLatency := fmt.Sprintf("%.1f ms", stats.AverageLatency.Seconds()*1000)
	minLatency := fmt.Sprintf("%.1f ms", stats.MinLatency.Seconds()*1000)
	maxLatency := fmt.Sprintf("%.1f ms", stats.MaxLatency.Seconds()*1000)

	data := TemplateData{
		SealRate:   stats.SealRate,
		AvgLatency: avgLatency,
		MinLatency: minLatency,
		MaxLatency: maxLatency,
		Throughput: stats.Throughput,
		TxHexes:    stats.TxHexes,
		Network:    stats.Network,
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
        }

        h1,
        h2 {
            color: #333;
        }

        table {
            width: 100%;
            margin-top: 20px;
            border-collapse: collapse;
        }

        th,
        td {
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

        .container {
            background-color: white;
            padding: 20px;
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
    <div class="section">
        <h2>Summary Table</h2>
        <table>
            <tr>
                <th>Metric</th>
                <th>Value</th>
            </tr>
            <tr>
                <td>Seal Rate (tps)</td>
                <td>{{.SealRate}}</td>
            </tr>
            <tr>
                <td>Minimum Latency</td>
                <td>{{.MinLatency}}</td>
            </tr>
            <tr>
                <td>Maximum Latency</td>
                <td>{{.MaxLatency}}</td>
            </tr>
            <tr>
                <td>Average Latency</td>
                <td>{{.AvgLatency}}</td>
            </tr>
            <tr>
                <td>Throughput (tps)</td>
                <td>{{.Throughput}}</td>
            </tr>
        </table>
    </div>
    <div class="container">
        <h1 class="project-name">Flowmark Blockchain Testnet Benchmark Report</h1>
        <div class="section">
            <h2>Seal Rate</h2>
            <p>Seal rate represents the rate at which blocks are sealed or mined on the blockchain network. It is calculated by dividing the total number of sealed blocks by the duration of the benchmarking period.</p>
            <p>Result: {{.SealRate}}</p>
        </div>
        <div class="section">
            <h2>Throughput</h2>
            <p>Throughput refers to the number of successful transactions processed per unit of time. It is calculated by dividing the total number of successful transactions by the duration of the benchmarking period.</p>
            <p>Result: {{.Throughput}}</p>
        </div>
        <div class="section">
            <h2>Average Latency</h2>
            <p>Latency refers to the time delay for transaction processing and confirmation on the blockchain network. Average latency is calculated by taking the average of all transaction latencies.</p>
            <p>Result: {{.AvgLatency}}</p>
        </div>
        <div class="section">
            <h2>Minimum Latency</h2>
            <p>Result: {{.MinLatency}}</p>
        </div>
        <div class="section">
            <h2>Maximum Latency</h2>
            <p>Result: {{.MaxLatency}}</p>
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
