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
	fmt.Println("reindexutxo - Rebuilds the UTXO set")
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

func (cli *CommandLine) reindexutxo() {
	chain := blockchain.ContinueBlockchain("")
	defer chain.Database.Close()
	utxoSet := blockchain.NewUTXOSet(chain)
	utxoSet.Reindex()

	count := utxoSet.CountTransactions()
	fmt.Printf("Done! There are %d transactions in the UTXO set.\n", count)
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
		for _, tx := range block.Transactions {
			fmt.Println(tx.String())
		}
		fmt.Println()
	}
}

func (cli *CommandLine) createBlockchain(address string) {
	if !wallet.ValidateAddress(address) {
		log.Panic("Invalid Address")
	}

	chain := blockchain.InitBlockchain(address)
	defer chain.Database.Close()

	fmt.Println("Finished!")
}

func (cli *CommandLine) getbalance(address string) {
	if !wallet.ValidateAddress(address) {
		log.Panic("Invalid Address")
	}

	chain := blockchain.ContinueBlockchain("")
	defer chain.Database.Close()
	utxoSet := blockchain.NewUTXOSet(chain)

	balance := 0
	publicKeyFullHash := wallet.Base58Decode([]byte(address))
	publicKeyHash := publicKeyFullHash[1 : len(publicKeyFullHash)-wallet.GetChecksumLength()]
	UTXOs := utxoSet.FindUnspentTransactionOutputs(publicKeyHash)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of %s: %d\n", address, balance)
}

func (cli *CommandLine) send(from string, to string, amount int) {
	if !wallet.ValidateAddress(to) {
		log.Panic("Invalid Address")
	}
	if !wallet.ValidateAddress(from) {
		log.Panic("Invalid Address")
	}

	chain := blockchain.ContinueBlockchain(from)
	defer chain.Database.Close()
	utxoSet := blockchain.NewUTXOSet(chain)

	tx := blockchain.NewTransaction(from, to, amount, utxoSet)
	block := chain.AddBlock([]*blockchain.Transaction{tx})
	utxoSet.Update(block)

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
	reindexCmd := flag.NewFlagSet("reindexutxo", flag.ExitOnError)

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
	case "reindexutxo":
		err := reindexCmd.Parse(os.Args[2:])
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

	if reindexCmd.Parsed() {
		cli.reindexutxo()
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
