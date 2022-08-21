package client

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	blockchain "gambim.com/blockchain/chain"
	"gambim.com/blockchain/wallet"
)

type CommandLine struct{}

func (cli *CommandLine) PrintUsage() {
	fmt.Println("Usage:")
	fmt.Println("getbalance -address ADDRESS - Get the balance")
	fmt.Println("createblockchain -address ADDRESS - Creates a blockchain")
	fmt.Println("printchain - Prints the block in the chain")
	fmt.Println("send -from FROM -to TO -amount AMOUNT - Send amount")
	fmt.Println("createwallet - Creates a new Wallet")
	fmt.Println("listaddresses - List the addresses in our wallet file")
}

func (cli *CommandLine) ValidateArgs() {
	if len(os.Args) < 2 {
		cli.PrintUsage()
		runtime.Goexit()
	}
}

func (cli *CommandLine) listWallets() {
	wallets, _ := wallet.CreateWallets()
	addresses := wallets.GetAllAddresses()

	for _, address := range addresses {
		fmt.Println(address)
	}
}

func (cli *CommandLine) createNewWalletCmd() {
	wallets, _ := wallet.CreateWallets()
	address := wallets.AddWallet()
	wallets.SaveFile()
	fmt.Printf("New address is: %s\n", address)
}

func (cli *CommandLine) printChain() {
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

func (cli *CommandLine) Run() {
	cli.ValidateArgs()

	getBalanceCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	listAddressesCmd := flag.NewFlagSet("listaddresses", flag.ExitOnError)

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
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "listaddresses":
		err := listAddressesCmd.Parse(os.Args[2:])
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
		cli.printChain()
	}

	if createWalletCmd.Parsed() {
		cli.createNewWalletCmd()
	}

	if listAddressesCmd.Parsed() {
		cli.listWallets()
	}
}
