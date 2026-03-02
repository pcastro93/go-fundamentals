package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type HelloResponse struct {
	CurrentTime int64  `json:"current_time"`
	Status      string `json:"status"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
			return
		}
		_, err := fmt.Fprintln(w, "pong")
		if err != nil {
			fmt.Println("Error writing response")
		}
	})
	mux.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
			return
		}
		defer r.Body.Close()
		bodyBytes, _ := io.ReadAll(r.Body)
		body := string(bodyBytes)
		_, err := fmt.Fprintln(w, body)
		if err != nil {
			fmt.Println("Error writing response")
		}
	})
	mux.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
			return
		}
		json.NewEncoder(w).Encode(HelloResponse{
			CurrentTime: time.Now().Unix(),
			Status:      "OK",
		})
	})
	fmt.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
