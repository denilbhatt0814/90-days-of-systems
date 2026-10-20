package main

import (
	"flag"
	"io"
	"log"
	"os"
	"runtime"
	"sync"
	"time"
	"uuid"
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

func cpuWork(iterations uint64) uint64 {
	// our job is to find sum of prime numbers till n = iterations
	var sum uint64 = 0
	for i := uint64(0); i < iterations; i++ {
		if isPrime(i) {
			sum += i
		}
	}

	return sum
}

func runExperiment(jobs, workers int, iterations uint64) (uint64, time.Duration) {
	jobQueue := make(chan struct{}, jobs)
	results := make(chan uint64, jobs)

	// appending jobs to the queue
	for i := 0; i < jobs; i++ {
		jobQueue <- struct{}{}
	}
	close(jobQueue)

	var wg sync.WaitGroup
	wg.Add(workers)

	start := time.Now()
	// spinning up workers in the workerpool
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()

			for range jobQueue {
				results <- cpuWork(iterations)
			}
		}()
	}

	wg.Wait()
	elapsed := time.Since(start)
	close(results)

	var checksum uint64
	for result := range results {
		checksum += result
	}
	return checksum, elapsed
}

func main() {
	file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer file.Close()

	// Combine the file and terminal outputs
	multiWriter := io.MultiWriter(os.Stdout, file)
	log.SetOutput(multiWriter)

	const jobs = 300
	const iterationsPerJob = 10000

	var workers int
	flag.IntVar(&workers, "workers", 1, "Number of workers in workerpool")
	flag.Parse()

	runID := uuid.New().String()
	log.Printf("STARTING RUN runID=%v", runID)
	log.Printf("CPUs: %d | GOMAXPROCS: %d\n", runtime.NumCPU(), runtime.GOMAXPROCS(0))

	checksum, elapsed := runExperiment(jobs, workers, iterationsPerJob)
	log.Printf("workers=%2d duration=%s checksum=%d\n", workers, elapsed, checksum)

	log.Printf("RUN ENDED runID=%v", runID)
}
