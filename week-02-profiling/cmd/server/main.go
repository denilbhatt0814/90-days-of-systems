package main

import (
	"encoding/json"
	"log"
	"net/http"
	_ "net/http/pprof"
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

		if value > 100000 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("n must be a valid int between 0 and 100000"))
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

type WaitResponse struct {
	DurationMs float64 `json:"duration_ms"`
}

func handleWaitTask(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	msWaitDurationStr := r.URL.Query().Get("ms")

	// parsing
	msWaitDuration := 100
	if msWaitDurationStr != "" {
		value, err := strconv.Atoi(msWaitDurationStr)
		if err != nil || value <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("ms must be a valid int above 0"))
			return
		}
		if value > 10000 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("ms must be a valid int between 0 and 10000"))
			return
		}
		msWaitDuration = value
	}

	time.Sleep(time.Duration(msWaitDuration) * time.Millisecond)
	elapsed := time.Since(start)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(WaitResponse{
		DurationMs: elapsed.Seconds() * 1000,
	})
}

func main() {

	go func() {
		log.Println("Profiling server listening on 127.0.0.1:6060")
		err := http.ListenAndServe("127.0.0.1:6060", http.DefaultServeMux)
		if err != nil {
			log.Printf("Profiling server failed: %v", err)
		}
	}()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /cpu", handleCPUTask)
	mux.HandleFunc("GET /wait", handleWaitTask)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Starting server on :8080")
	log.Fatal(server.ListenAndServe())
}
