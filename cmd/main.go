package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Hello World!")

	server := http.NewServeMux()
	server.HandleFunc("/api/ping", ping)

	port := ":8888"

	if err := http.ListenAndServe(port, server); err != nil {
		log.Fatalf("Couldn't setup server %v", err)
	}
}

func ping(w http.ResponseWriter, r *http.Request) {
	var pong struct {
		Pong string `json:"pong"`
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pong)
}
