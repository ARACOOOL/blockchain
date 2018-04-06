package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Block struct {
	Timestamp    int64
	Hash         []byte
	PreviousHash []byte
	Nonce        int
	Payload      Payload
}

func (b *Block) CreateHash() {
	cr := sha256.New()
	sum := b.Payload.ToString()
	sum = append(b.PreviousHash, byte(b.Nonce))

	cr.Write(sum)
	b.Hash = cr.Sum(nil)
}

type PayloadData struct {
	FirstName string
	LastName  string
}

type Payload struct {
	Data      PayloadData
	Node      []byte
	Signature []byte
}

func (p Payload) ToString() []byte {
	b, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}

	return b
}

func (p *Payload) Sign() {
	jsonByte, err := json.Marshal(p.Data)
	if err != nil {
		panic(err)
	}

	cr := sha256.New()
	cr.Write(jsonByte)
	d := cr.Sum(nil)

	prvFile, err := os.OpenFile("keys/prv", os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal("Could not open private key file")
	}

	keyContent, err := ioutil.ReadAll(prvFile)
	if err != nil {
		log.Fatal("Could not read private key file")
	}

	block, _ := pem.Decode(keyContent)
	pk, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	fmt.Println(pk.N.String())
	p.Signature, err = rsa.SignPKCS1v15(rand.Reader, pk, crypto.SHA256, d)
}
