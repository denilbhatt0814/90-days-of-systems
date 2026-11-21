package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

// We shall assume finding isPrime is a highly CPU intensive task
func isPrime(n uint) bool {
	if n <= 1 {
		return false
	}

	for i := uint(2); i < n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func countPrimeToN(N uint) uint {

	var count uint = 0
	for i := uint(2); i <= N; i++ {
		if isPrime(i) {
			count++
		}
	}

	return count
}

type CPUResponse struct {
	Result     uint    `json:"result"`
	DurationMs float64 `json:"duration_ms"`
}

func handleCPUTask(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	nStr := r.URL.Query().Get("n")

	// parsing
	n := 10000
	if nStr != "" {
		value, err := strconv.Atoi(nStr)
		if err != nil || value <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("n must be a valid int above 0"))
			return
		}
		n = value
	}

	result := countPrimeToN(uint(n))
	elapsed := time.Since(start)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(CPUResponse{
		Result:     result,
		DurationMs: elapsed.Seconds() * 1000,
	})
	if err != nil {
		return
	}
}

func handleWaitTask(w http.ResponseWriter, r *http.Request) {
	// start := time.Now()
	msWaitDurationStr := r.URL.Query().Get("ms")

	// parsing
	msWaitDuration := 100
	if msWaitDurationStr != "" {
		value, err := strconv.Atoi(msWaitDurationStr)
		if err != nil || value <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("n must be a valid int above 0"))
			return
		}
		msWaitDuration = value
	}

	time.Sleep(time.Duration(msWaitDuration) * time.Millisecond)
	// elapsed := time.Since(start)
	w.WriteHeader(http.StatusOK)
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /cpu", handleCPUTask)
	mux.HandleFunc("GET /wait", handleWaitTask)

	server := &http.Server{
		Addr:    ":6060",
		Handler: mux,
	}

	log.Println("Starting server on :6060")
	log.Fatal(server.ListenAndServe())
}
