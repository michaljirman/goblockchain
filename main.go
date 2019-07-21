package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/michaljirman/goblockchain/config"
	"github.com/michaljirman/goblockchain/logger"

	"github.com/rs/zerolog/log"

	"github.com/michaljirman/goblockchain/blockchain"
)

type CommandLine struct {
	blockchain *blockchain.BlockChain
}

func (cli *CommandLine) printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" add -block BLOCK_DATA - add a block to the chain")
	fmt.Println(" print - Prints the blocks in the chain")
}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}

func (cli *CommandLine) addBlock(data string) {
	cli.blockchain.AddBlock(data)
	log.Info().Msg("a new block added")
}

func (cli *CommandLine) printChain() {
	iter := cli.blockchain.Iterator()
	for {
		block := iter.Next()
		pow := blockchain.NewProof(block)
		log.Info().Hex("prev_hash", block.PrevHash).
			Str("data", fmt.Sprintf("%s", block.Data)).
			Hex("hash", block.Hash).
			Bool("pow", pow.Validate()).
			Msg("block information")
		log.Info().Msg("")

		if len(block.PrevHash) == 0 {
			break
		}
	}
}

func (cli *CommandLine) run() {
	cli.validateArgs()
	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	addBlockData := addBlockCmd.String("block", "", "Block data")

	switch os.Args[1] {
	case "add":
		err := addBlockCmd.Parse(os.Args[2:])
		blockchain.HandleError(err)
	case "print":
		err := printChainCmd.Parse(os.Args[2:])
		blockchain.HandleError(err)
	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			runtime.Goexit()
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}

//PRETTY_LOGGING=TRUE go run main.go
func main() {
	cfg, err := config.GetConfigFromEnv()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config from env")
	}

	// setup global logger settings
	err = logger.SetupLogger(&cfg.Log)
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to setup logger")
	}

	log.Debug().Msg("test")
	defer os.Exit(0)
	chain := blockchain.InitBlockChain()
	defer chain.Database.Close()

	cli := CommandLine{chain}
	cli.run()
}
