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

func externalServiceCall(ioTime time.Duration) {
	time.Sleep(ioTime)
}

func ioWork(ioTime time.Duration) uint64 {
	externalServiceCall(ioTime)
	return 1
}

func runExperiment(jobs, workers int, ioTime time.Duration) (uint64, time.Duration) {
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
				results <- ioWork(ioTime)
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
	const ioTime = 20 * time.Millisecond

	var workers int
	flag.IntVar(&workers, "workers", 1, "Number of workers in workerpool")
	flag.Parse()

	runID := uuid.New().String()
	log.Printf("STARTING RUN (I/O Workload) runID=%v", runID)
	log.Printf("CPUs: %d | GOMAXPROCS: %d\n", runtime.NumCPU(), runtime.GOMAXPROCS(0))

	checksum, elapsed := runExperiment(jobs, workers, ioTime)
	log.Printf("workers=%2d duration=%s checksum=%d\n", workers, elapsed, checksum)

	log.Printf("RUN ENDED (I/O Workload) runID=%v", runID)
}
