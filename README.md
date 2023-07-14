# FlowMark

FlowMark is a benchmarking tool designed to evaluate the performance of the Flow Blockchain. This tool sends multiple transactions to the Flow Blockchain and gathers statistics such as send rate, seal rate, average send latency, average seal latency, minimum network latency, maximum network latency, and average network latency.


## Prequisities

Make sure you have installed all of the following prerequisites on your development machine:

 - Golang - **[Download and install](https://go.dev/doc/install)** Golang and ensure that the **`$GOPATH`** environment variable has been set.

## Cloning the GitHub Repository

The easiest way to get started with the FlowMark tool is to clone the repository from Github and build it locally. You can clone the repository by running the following command:
```
    git clone https://github.com/7suyash7/FlowMark.git
```

## Setting up the Environment Variables

The FlowMark tool requires several environment variables to be set in order to function correctly. These variables can be set in a **`.env`** file in the project's root directory. An example **`.env`** file looks like this:
Make sure to rename the from **`.env.example`** to **`.env`**
```
    NO_OF_TRANSACTION=5
    RECIPIENT_ADDRESS="xxxxxxxxxxxxxxxx"
    SENDER_ADDRESS="xxxxxxxxxxxxxxxx"
    SENDER_PRIVATE_KEY="xxxxxxxxxxxxxxxxxxxxxxx"
```
You can also setup the config through the command line, instructions for that are mentioned below.

## Building 

You can build the application by running the following command in the project's root directory:
```
    make build
```
This will compile the application and create a binary name **FlowMark**

## Running the Application

There are two ways to run the binary.

1. Running the binary when you have setup the **`.env`** yourself.
 ```
    make run
```
 2. Setting up the **`.env`** through the command 
```
     ./FlowMark start
```
The application supports command-line options to override the environment variables. Here is how it can be used:
```
    ./FlowMark start --sender-address <senderAddress> --receiver-address <receiverAddress> --numTransaction <numberOfTransactions> --network <testnet> --sender-priv-address <senderAccountPrivKey>
```

## Help and Manual

You can display the help manual by running the following command:
```
    ./FlowMark help
```
To display the current configuration, use the following command:
```
    ./FlowMark config
```

## Generating Reports
At the end of benchmark test, the tool generates a detailed report file **(`report.html`)**. This report contains the performance metrics for the test run. The results are also displayed in the terminal after the benchmark is finished.
# FlowMark

FlowMark is a benchmarking tool designed to evaluate the performance of the Flow Blockchain. This tool sends multiple transactions to the Flow Blockchain and gathers statistics such as send rate, seal rate, average send latency, average seal latency, minimum network latency, maximum network latency, and average network latency.


## Prequisities

Make sure you have installed all of the following prerequisites on your development machine:

 - Golang - **[Download and install](https://go.dev/doc/install)** Golang and ensure that the **`$GOPATH`** environment variable has been set.

## Cloning the GitHub Repository

The easiest way to get started with the FlowMark tool is to clone the repository from Github and build it locally. You can clone the repository by running the following command:
```
    git clone https://github.com/7suyash7/FlowMark.git
```

## Setting up the Environment Variables

The FlowMark tool requires several environment variables to be set in order to function correctly. These variables can be set in a **`.env`** file in the project's root directory. An example **`.env`** file looks like this:
Make sure to rename the from **`.env.example`** to **`.env`**
```
    NO_OF_TRANSACTION=5
    RECIPIENT_ADDRESS="xxxxxxxxxxxxxxxx"
    SENDER_ADDRESS="xxxxxxxxxxxxxxxx"
    SENDER_PRIVATE_KEY="xxxxxxxxxxxxxxxxxxxxxxx"
```
You can also setup the config through the command line, instructions for that are mentioned below.

## Building 

You can build the application by running the following command in the project's root directory:
```
    make build
```
This will compile the application and create a binary name **FlowMark**

## Running the Application

There are two ways to run the binary.

1. Running the binary when you have the **`.env`** setup yourself.
 ```
    make run
```

2. Running the binary using command-line options. Here is how it can be used:
```
    ./FlowMark start --sender-address <senderAddress> --receiver-address <receiverAddress> --numTransaction <numberOfTransactions> --network <testnet> --sender-priv-address <senderAccountPrivKey>
```

## Help and Manual

You can display the help manual by running the following command:
```
    ./FlowMark help
```
To display the current configuration, use the following command:
```
    ./FlowMark config
```

## Generating Reports
At the end of benchmark test, the tool generates a detailed report file **(`report.html`)**. This report contains the performance metrics for the test run. The results are also displayed in the terminal after the benchmark is finished.