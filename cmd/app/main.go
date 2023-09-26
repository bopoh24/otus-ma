package main

import (
	"log"
	"net/http"
)

func main() {
	log.Printf("starting server on port 8000")
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"status": "OK"}`))
		if err != nil {
			log.Printf("error writing response body: %v", err)
		} else {
			log.Printf("responded to health check request")
		}
	})
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	http.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	log.Fatal(http.ListenAndServe(":8000", nil))
}
