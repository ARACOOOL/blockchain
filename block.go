package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
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
	sum := b.Payload.ToBytes()
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

func (p Payload) ToBytes() []byte {
	b, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}

	return b
}

func (p *Payload) Sign() {
	cr := sha256.New()
	cr.Write(p.ToBytes())
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

	p.Node = x509.MarshalPKCS1PublicKey(&pk.PublicKey)
	p.Signature, err = rsa.SignPKCS1v15(rand.Reader, pk, crypto.SHA256, d)
}
