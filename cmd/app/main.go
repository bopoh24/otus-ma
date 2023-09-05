package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/health/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"status": "OK"}`))
		if err != nil {
			log.Printf("error writing response body: %v", err)
		} else {
			log.Printf("responded to health check request")
		}
	})
	log.Fatal(http.ListenAndServe(":8000", nil))
}
