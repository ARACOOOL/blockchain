package main

import (
	"log"

	"github.com/joho/godotenv"
)

var blockchain *BlockChain

func main() {
	blockchain = &BlockChain{}

	genesisBlock := CreateNewBlock(Block{}, PayloadData{})
	blockchain.AddBlock(genesisBlock)

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	StartServer()
}
