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
		if !IsBlockValid(lastBlock, block) {
			log.Println("Block is invalid")
			return
		}
	}

	chain.blocks = append(chain.blocks, block)
}

func (chain *BlockChain) GetLastBlock() Block {
	return chain.blocks[len(chain.blocks)-1]
}
