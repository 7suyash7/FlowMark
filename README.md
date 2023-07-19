
# FlowMark

[Pitch Deck](https://docs.google.com/presentation/d/17dlC4BCtM9YbDP4K2p4MJBlBShgg3x2FuhRq-WNAGFA/edit?usp=sharing)

Welcome to the Flow Blockchain Benchmarking Tool, a comprehensive utility designed to evaluate and measure the performance of the Flow Blockchain. This tool is an essential asset for developers, researchers, and enthusiasts who are interested in understanding the performance characteristics of the Flow Blockchain under various conditions.

Benchmarking is a crucial process in the world of blockchain technology. It provides a systematic and quantitative evaluation of the system's performance, including aspects like transaction speed, seal time, smart contract execution efficiency, network latency, etv. Benchmarking helps identify bottlenecks, areas for improvement, and provides a basis for comparison with other blockchain systems.

Our Flow Blockchain Benchmarking Tool is designed to provide these insights by simulating various workloads and network conditions, and then measuring how the Flow Blockchain performs under these circumstances. The tool provides a range of metrics and can be customized to focus on specific areas of interest. Whether you're a developer looking to optimize your dApp, a researcher studying blockchain performance, or simply a blockchain enthusiast, this tool provides valuable insights into the workings of the Flow Blockchain.

## Table of Contents
- [FlowMark](#flowmark)
  * [Prequisities and Dependencies](#prequisities-and-dependencies)
  * [Cloning the GitHub Repository](#cloning-the-github-repository)
  * [Setting up the Benchmark settings](#setting-up-the-benchmark-settings)
    + [- Test](#--test)
    + [- Workers](#--workers)
    + [- Rounds](#--rounds)
  * [Setting up the settings for Transactions](#setting-up-the-settings-for-transactions)
  * [Building and Running the Benchmark](#building-and-running-the-benchmark)
- [ADD HTML SCREENSHOT HERE](#add-html-screenshot-here)
  * [Understanding the Metrics](#understanding-the-metrics)
  * [How it Works (Architecture)](#how-it-works--architecture-)
  * [Features and Benefits for the Flow Blockchain Ecosystem](#features-and-benefits-for-the-flow-blockchain-ecosystem)
    + [Research Tool](#research-tool)
    + [Identifying Bottlenecks](#identifying-bottlenecks)
    + [Performance Overview](#performance-overview)
    + [Future Development](#future-development)

## Prequisities and Dependencies

Make sure you have installed all of the following prerequisites on your development machine:

 - Golang : **[Download and install](https://go.dev/doc/install)** Golang and ensure that the **`$GOPATH`** environment variable has been set.
 - Flow-Cli : **[Download and install](https://developers.flow.com/tooling/flow-cli/install)** Flow-Cli for running the local emulator and account generation.  

## Cloning the GitHub Repository

The easiest way to get started with the FlowMark tool is to clone the repository from Github and build it locally. You can clone the repository by running the following command:
```
    git clone https://github.com/7suyash7/FlowMark.git
```

## Setting up the Benchmark settings
The **`benchmarkConfig.yaml`** file is the heart of the Flow Blockchain Benchmarking Tool. It allows you to define the parameters of your benchmark tests, including the network to be tested, the type of test, the number of workers, and the specifics of each round of testing. 

The default example included in the repository uses the emulator network to run the benchmarks. Hence, make sure that the emulator is running using the following command:
```
flow emulator start
```
This will give you a default account that can be used as the sender address. 
To generate the receiver address you can run the following command in a new terminal.
```
flow accounts create
```
Make sure to select the **`emulator`** option on the terminal to generate an account on the emulator network.
Here's a breakdown of each section and how you can configure it:
### - Test
 - **network**: This field specifies the network on which the benchmark will be run. In the example, it's set to **"emulator"**, but it could be **"mainnet"** or **"testnet"**.

 - **name**: This is the name of the test. It's a string that should briefly describe the test being performed. In this case, it's **"Test"**.

 - **description**: This field provides a more detailed explanation of what the test is doing. Here, it's set to "To benchmark transferring tokens between accounts."

### - Workers
 - **number**: This field specifies the number of workers that will be used to perform the test. Workers are essentially concurrent threads that execute the transactions. In the example, it's set to 1, but you can increase this number to simulate higher loads.

### - Rounds

Each round represents a different set of transactions that will be executed as part of the test. You can define multiple rounds with different parameters. Each round has the following fields:

 - **label**: This is a brief description of the round. It should give an idea of what the round is testing. For example, **"50 txns with 1tps"** means this round will execute 50 transactions at a rate of 1 transaction per second.

 - **description**: This field provides a more detailed explanation of what the round is doing to provide more context

 - **rateControl**: This section defines the specifics of the transactions that will be executed during the round.

 - **txNumber**: This is the total number of transactions that will be executed during the round. In the first round of the example, it's set to **50**.

 - **tps**: This is the rate at which transactions will be executed, measured in transactions per second. In the first round of the example, it's set to **1**.

By adjusting these parameters, you can create a wide variety of tests to benchmark the Flow Blockchain under different conditions. Remember to save your changes to the **`benchmarkConfig.yaml`** file before running the benchmark tool.

## Setting up the settings for Transactions

Note: Currently there's a bug that will cause a transaction to fail sometimes when using custom scripts, this happens because the arguments don't load in order all the time. It's fixable but we didn't have enought time :D
  This is the error that might show up -
    * transaction execute failed: [Error Code: 1101] cadence runtime error: Execution failed:
  error: invalid argument at index 0: expected value of type `UFix64`


The **`transactionConfig.yaml`** file is where you define the specifics of the transactions that will be executed during the benchmark tests. This includes the path to the Cadence script that will be executed, the gas limit for the transactions, the arguments passed to the script, and the details of the accounts involved in the transaction. The accounts that we generated in the previous section can be used here.
Here's a breakdown of each section:

 - **scriptPath**: This is the path to the Cadence script that will be executed during the benchmark test. Replace "/path/to/script.cdc" with the actual path to your script.

 - **gasLimit**: This is the maximum amount of gas that can be used by each transaction. In the example, it's set to 100000.

 - **scriptArguments**: This section defines the arguments that will be passed to the Cadence script. Each argument has a type and a value. The type should be one of the 13 supported script argument types, and the value should be a valid value for that type. In the example, two arguments are defined: amount and recipient.

 - **payer, proposer, authorizer**: These sections define the accounts that will be used for the transaction. Each account has an address and a privateKey. If useSameAccount is set to true, the same account will be used for all three roles. Currently, we only support using the same account for the payer, proposer, and authorizer (single party, multiple signatures). In the future, we will be adding support for multiple parties with multiple signatures, multiple parties with two authorizers, and multiple parties in general.
 
The tool currently supports 13 script argument types:
1. String
2. Address
3. Boolean
4. Fix64
5. UFix64
6. UInt8
7. UInt16
8. UInt32
9. UInt64
10. Int8
11. Int16
12. Int32
13. Int64

Support is being added for more script arguments types to allow for a wide variety of scripts to be used in benchmark tests.

 - Please note that currently, we only support addresses generated using the signature algorithm: **`ECDSA_P256`** and Hash Algorithm: **`SHA3_256`**. Future updates will include support for **`ECDSA_secp256k1`** and **`SHA2_256`**.
 
 - Remember to replace all placeholder values (marked with "xxxxxxxx") with your actual data, and to save your changes to the **`transactionConfig.yaml`** file before running the benchmark tool.

## Building and Running the Benchmark
Building the application is a straightforward process. This step compiles the application and creates a binary named FlowMark. To do this, you need to run a specific command in the project's root directory.

Here are the steps:
1. Open a terminal window.
2. Navigate to the project's root directory. If your project is located at **`path/to/project`**, you can do this by running: 
```
cd /path/to/project
```
3. Once you're in the projects root directory, run the following command:
```
make build
```
This command triggers the build process. It may take a few moments to complete. Once it's done, you'll find the **`FlowMark`** binary in the project's root directory. This binary is the executable form of the application, and you can run it to start the benchmark tests.
4. Running the benchmark can be done by running the following command in the root directory:
```
make run
```
This command starts the benchmark tests. The results for each round will be displayed as they're completed, and a summary table will be shown at the end.

The summary table contains the following metrics:
1. **Name**: The label of each round.
2. **Send Rate**: The rate at which transactions were sent.
3. **Seal Rate**: The rate at which transactions were sealed.
4. **Maximum Network Latency**: The longest network latency throughout the benchmark.
5. **Minimum Network Latency**: The shortest network latency calculated throughout the benchmark.
6. **Average Network Latency**: The average time of the network latency throughout the benchmark.
7. **Successful Transactions**: The number of transactions that were successfully sealed.
8. **Failed Transactions**: The number of transactions that failed to be sealed.

These metrics provide a comprehensive overview of the performance of the Flow Blockchain under the conditions defined in your **`benchmarkConfig.yaml`** and **`transactionConfig.yaml files`**.

The Results are displayed in the terminal and also a **`report.html`** is generated in the project root directory that can be opened from the terminal and it displays these metrics in a more accurate and extensive format.

!(html page)[./img/HTMLpage.jpeg]

## Understanding the Metrics
The benchmarking tool provides a range of metrics that offer insights into the performance of the Flow Blockchain under different conditions. Here's what each metric means:

1. **Name**: This is the name of the test. It helps you identify the specific test that was run, especially when you're running multiple tests in a single benchmarking session.

2. **Send Rate**: This is the rate at which transactions were sent to the network, measured in transactions per second (tps). It indicates the load that was applied to the network during the test.

3. **Seal Rate**: This is the rate at which transactions were sealed (i.e., finalized) by the network, also measured in transactions per second. It indicates how quickly the network was able to process the transactions.

4. **Max Latency**: This is the longest network latency observed for a transaction, measured from the time it was sent from the client to the time it was received by the network. It's measured in milliseconds (ms) and provides an indication of the worst-case network latency during the test.

5. **Min Latency**: This is the shortest network latency observed for a transaction, measured in the same way as the max latency. It provides an indication of the best-case network latency during the test.

6. **Avg Latency**: This is the average network latency observed for transactions. It's calculated by adding up the latency for each transaction and dividing by the total number of transactions. It provides a general indication of the network latency during the test.

7. **Successful Transactions**: This is the number of transactions that were successfully sealed by the network. It provides an indication of the reliability of the network under the conditions of the test.

8. **Failed Transactions**: This is the number of transactions that failed to be sealed by the network. It also provides an indication of the reliability of the network.

These metrics together provide a comprehensive overview of the performance and reliability of the Flow Blockchain under the conditions of the test. By adjusting the parameters of the test, you can use these metrics to understand how the network behaves under different loads and conditions.

## How it Works (Architecture)

![Picture of architecture](./img/architecture.jpg)

## Features and Benefits for the Flow Blockchain Ecosystem
The Flow Blockchain Benchmarking Tool is a powerful utility that offers a range of features designed to help understand and optimize the performance of the Flow Blockchain. Here's how it can benefit the Flow Blockchain ecosystem:

### Research Tool
The benchmarking tool serves as a valuable research instrument for academics, developers, and blockchain enthusiasts. It allows for systematic and quantitative evaluation of the Flow Blockchain's performance under various conditions. This can lead to a deeper understanding of the blockchain's behavior, its strengths and weaknesses, and how it responds to different loads and network conditions.

By providing a clear and detailed picture of the blockchain's performance, the tool can facilitate the writing of research papers about the Flow Blockchain. These papers can contribute to the broader academic discourse on blockchain technology, providing valuable insights for other researchers and developers.

For example, benchmarking tools have been instrumental in the production of research papers on other blockchains. Various research papers and studies have been written using
a variety of Blockchain Benchmarking tools.
1. ["Measuring performances and footprint of blockchains with BCTMark"](https://link.springer.com/article/10.1007/s10586-021-03441-x)
2. ["Performance Evaluation of Ethereum Private and Testnet Networks Using Hyperledger Caliper"](https://ieeexplore.ieee.org/abstract/document/9562684)
3. ["Performance Analysis of a Hyperledger Fabric Blockchain Framework: Throughput, Latency and Scalability"](https://ieeexplore.ieee.org/abstract/document/8946222)
4. ["A framework for automating deployment and evaluation of blockchain networks
Author links open overlay panel"](https://www.sciencedirect.com/science/article/abs/pii/S1084804522001102)

### Identifying Bottlenecks
The benchmarking tool can help identify bottlenecks in the Flow Blockchain. By running tests under different conditions, you can see where performance issues arise and where the system struggles to keep up. This information can be invaluable for developers working on the Flow Blockchain, as it can guide them in optimizing the system and improving its performance.

### Performance Overview
The tool provides a comprehensive overview of the Flow Blockchain's performance. It measures a range of metrics, including send rate, seal rate, network latency, and the number of successful and failed transactions. These metrics provide a clear picture of how the blockchain performs under different conditions, which can be useful for anyone using or developing for the Flow Blockchain.

### Future Development
The benchmarking tool can also guide future development of the Flow Blockchain. By identifying areas where the blockchain performs well and areas where it struggles, it can help developers and researchers focus their efforts where they're most needed. This can lead to a more efficient and effective blockchain, benefiting everyone who uses it.

In conclusion, the Flow Blockchain Benchmarking Tool is a powerful utility that can provide valuable insights into the performance of the Flow Blockchain. Whether you're a researcher studying blockchain technology, a developer working on the Flow Blockchain, or simply a blockchain enthusiast, this tool can provide you with the information you need to understand and optimize the Flow Blockchain.