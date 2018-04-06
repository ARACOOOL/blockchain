package main

import (
	"log"
	"sync"
)

type BlockChain struct {
	sync.Mutex
	blocks []Block
}

func (chain *BlockChain) AddBlock(block Block) {
	chain.Lock()
	defer chain.Unlock()

	if len(chain.blocks) > 0 {
		lastBlock := chain.blocks[len(chain.blocks)-1]
		if !IsValidBlock(lastBlock, block) {
			log.Fatal("Block is invalid")
			return
		}
	}

	chain.blocks = append(chain.blocks, block)
}
