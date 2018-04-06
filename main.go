package main

var blockchain *BlockChain

func main() {
	blockchain = &BlockChain{}

	block := CreateNewBlock(Block{}, PayloadData{FirstName: "s", LastName: "f"})
	blockchain.AddBlock(block)
}
