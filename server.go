package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func StartServer() {
	httpPort := os.Getenv("PORT")
	log.Println("HTTP Server Listening on port :", httpPort)

	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/", handleWriteBlock).Methods("POST")

	s := &http.Server{
		Addr:           ":" + httpPort,
		Handler:        muxRouter,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

// handleWriteBlock adds the new block to our chain
func handleWriteBlock(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	var m PayloadData

	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(writer, request, http.StatusBadRequest, request.Body)
		return
	}
	defer request.Body.Close()

	lastBlock := blockchain.GetLastBlock()
	newBlock := CreateNewBlock(lastBlock, m)

	blockchain.AddBlock(newBlock)

	respondWithJSON(writer, request, http.StatusCreated, newBlock)
}

// handleGetBlockchain returns the list of block
func handleGetBlockchain(writer http.ResponseWriter, request *http.Request) {
	bytes, err := json.MarshalIndent(blockchain.blocks, "", "  ")
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	io.WriteString(writer, string(bytes))
}

// respondWithJSON returns JSON response
func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}
