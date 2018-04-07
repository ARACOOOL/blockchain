package main

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"io/ioutil"
	"log"
	mathRand "math/rand"
	"os"
	"time"
)

const NonceMax  = 99999999

type Block struct {
	Timestamp    int64
	Hash         []byte
	PreviousHash []byte
	Nonce        int
	Payload      Payload
}

func CreateHash(b Block) []byte {
	cr := sha256.New()
	sum := b.Payload.ToBytes()
	sum = append(sum, byte(b.Nonce))
	sum = append(sum, b.PreviousHash...)
	sum = append(sum, b.Payload.Node...)
	sum = append(sum, b.Payload.Signature...)

	cr.Write(sum)
	return cr.Sum(nil)
}

func CreateNewBlock(block Block, data PayloadData) Block {
	payload := &Payload{Data: data}
	payload.Sign()

	newBlock := Block{Timestamp: time.Now().UnixNano(), PreviousHash: block.Hash, Nonce: mathRand.Intn(NonceMax), Payload: *payload}
	newBlock.Hash = CreateHash(newBlock)

	return newBlock
}

func IsBlockValid(oldBlock Block, newBlock Block) bool {
	if bytes.Compare(oldBlock.Hash, newBlock.PreviousHash) != 0 {
		log.Println("Previous hash invalid")
		return false
	}

	if bytes.Compare(CreateHash(newBlock), newBlock.Hash) != 0 {
		log.Println("New block hash is invalid")
		return false
	}

	cr := sha256.New()
	cr.Write(newBlock.Payload.ToBytes())
	d := cr.Sum(nil)

	pubKey, err := x509.ParsePKCS1PublicKey(newBlock.Payload.Node)
	if err != nil {
		log.Println("Public is invalid")
		return false
	}

	err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, d, newBlock.Payload.Signature)
	if err != nil {
		log.Println("Signature is invalid")
		return false
	}

	return true
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
	b, err := json.Marshal(p.Data)
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
