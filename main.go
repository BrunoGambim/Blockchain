package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"gambim.com/blockchain/blockchain"
)

type CommandLine struct{}

func (cli *CommandLine) PrintUsage() {
	fmt.Println("Usage:")
	fmt.Println("getbalance -address ADDRESS - Get the balance")
	fmt.Println("createblockchain -address ADDRESS - Creates a blockchain")
	fmt.Println("printchain - Prints the block in the chain")
	fmt.Println("send -from FROM -to TO -amount AMOUNT - Send amount")
}

func (cli *CommandLine) ValidateArgs() {
	if len(os.Args) < 2 {
		cli.PrintUsage()
		runtime.Goexit()
	}
}

func (cli *CommandLine) PrintChain() {
	chain := blockchain.ContinueBlockchain("")
	defer chain.Database.Close()

	iterator := chain.Iterator()
	for len(iterator.IteratorHash) != 0 {
		block := iterator.Next()
		fmt.Printf("hash:%x  prev_hash:%x\n", block.Hash, block.PrevHash)
	}
}

func (cli *CommandLine) createBlockchain(address string) {
	chain := blockchain.InitBlockchain(address)
	defer chain.Database.Close()

	fmt.Println("Finished!")
}

func (cli *CommandLine) getbalance(address string) {
	chain := blockchain.ContinueBlockchain("")
	defer chain.Database.Close()

	balance := 0
	UTXOs := chain.FindUnspentTransactionsOutputs(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of %s: %d\n", address, balance)
}

func (cli *CommandLine) send(from string, to string, amount int) {
	chain := blockchain.ContinueBlockchain("")
	defer chain.Database.Close()

	tx := blockchain.NewTransaction(from, to, amount, chain)
	chain.AddBlock([]*blockchain.Transaction{tx})

	fmt.Printf("Success!")
}

func (cli *CommandLine) run() {
	cli.ValidateArgs()

	fmt.Println("getbalance -address ADDRESS - Get the balance")
	fmt.Println("createblockchain -address ADDRESS - Creates a blockchain")
	fmt.Println("printchain - Prints the block in the chain")
	fmt.Println("send -from FROM -to TO -amount AMOUNT - Send amount")

	getBalanceCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "The wallet address")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The wallet address")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			runtime.Goexit()
		}
		cli.getbalance(*getBalanceAddress)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			runtime.Goexit()
		}
		cli.createBlockchain(*createBlockchainAddress)
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount == 0 {
			sendCmd.Usage()
			runtime.Goexit()
		}
		cli.send(*sendFrom, *sendTo, *sendAmount)
	}

	if printChainCmd.Parsed() {
		cli.PrintChain()
	}
}

func main() {
	defer os.Exit(0)
	cli := CommandLine{}
	cli.run()
}
