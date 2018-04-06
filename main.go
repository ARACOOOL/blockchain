package main

import (
	"fmt"
	"time"
)

func main() {
	d := Payload{Data: PayloadData{FirstName: "s", LastName: "f"}}
	d.Sign()
	block := &Block{Timestamp: time.Now().Unix(), PreviousHash: []byte(`0000000000000000000000000000000000000000000000000000000000000000`), Nonce: 23423, Payload: d}
	block.CreateHash()
	fmt.Println(block)
}
