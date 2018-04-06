package main

import (
	"fmt"
	"time"
)

func main() {
	d := Payload{Data: PayloadData{FirstName: "s", LastName: "f"}, Node: []byte(`047d40e7d9981a27ae94d6cd76909092423fcb0d63c280fac76b7008262528c62b641006fe17c03df032e73c8ae825de3d1a9c9a32162ad6c68a6611897f989592`)}
	d.Sign()
	block := &Block{Timestamp: time.Now().Unix(), PreviousHash: []byte(`0000000000000000000000000000000000000000000000000000000000000000`), Nonce: 23423, Payload: d}
	block.CreateHash()
	fmt.Println(block)
}
