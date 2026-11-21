package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

// We shall assume finding isPrime is a highly CPU intensive task
func isPrime(n uint64) bool {
	if n <= 1 {
		return false
	}

	for i := uint64(2); i < n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func countPrimeToN(N uint64) uint64 {

	var count uint64 = 0
	for i := uint64(2); i <= N; i++ {
		if isPrime(i) {
			count++
		}
	}

	return count
}

func handleCPUTask(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	nStr := r.URL.Query().Get("n")

	// parsing
	n := 10000
	if nStr != "" {
		n, err := strconv.Atoi(nStr)
		if err != nil || n <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("n must be a valid int above 0"))
			return
		}
		n = n
	}

	result := countPrimeToN(uint64(n))
	elapsed := time.Since(start)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(struct {
		result      uint64
		duration_ms int64
	}{
		result:      result,
		duration_ms: elapsed.Milliseconds(),
	})
	if err != nil {
		return
	}
}

func handleWaitTask(w http.ResponseWriter, r *http.Request) {

}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /cpu", handleCPUTask)
	mux.HandleFunc("GET /wait", handleWaitTask)

	server := &http.Server{
		Addr:    ":6060",
		Handler: mux,
	}

	log.Println("Starting server on :8080")
	log.Fatal(server.ListenAndServe())
}
