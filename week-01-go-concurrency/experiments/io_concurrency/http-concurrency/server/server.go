package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
	"uuid"
)

type WorkerPool struct {
	poolSize      int
	jobQueue      chan *job
	waitOnClosing *sync.WaitGroup
}

func (w *WorkerPool) addJob(j *job) {
	w.jobQueue <- j
}

func (w *WorkerPool) close() {
	log.Printf("stopping worker pool...")
	close(w.jobQueue)
	w.waitOnClosing.Wait()
}

func (w *WorkerPool) init(waitForStart *sync.WaitGroup) {
	var wg sync.WaitGroup
	w.waitOnClosing.Add(1)
	wg.Add(w.poolSize)

	// spinning up workers in the workerpool
	for i := 0; i < w.poolSize; i++ {
		go func() {
			defer wg.Done()

			for job := range w.jobQueue {
				doJob(job)
			}
		}()
	}

	log.Printf("Worker pool started poolSize=%d\n", w.poolSize)
	waitForStart.Done()
	wg.Wait()
	log.Printf("Worker pool stopped\n")
	w.waitOnClosing.Done()
}

type workers struct {
	workerPool           WorkerPool
	perJobProcessingTime time.Duration
}

type job struct {
	jobId          string
	processingTime time.Duration
	resultChannel  chan struct{}
}

func doJob(j *job) {
	currentActive := activeProcessingCounter.Add(1)
	updateMaxActive(currentActive)
	defer activeProcessingCounter.Add(-1)
	// log.Printf("Starting Job jobId=%s activeProcessingCounter=%d\n", j.jobId, counter)

	time.Sleep(j.processingTime)

	j.resultChannel <- struct{}{}
	// log.Printf("Job Completed jobId=%s\n", j.jobId)
}

func (workers *workers) workHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	traceId := r.Header.Get("x-trace-id")

	// we need to lock down the processing slots to 8
	// we will need to have buffer channel/worker pool with 8 workers
	// new requests simply adds a job to the pool
	// job is to sleep 50ms and return
	jobId := traceId
	if traceId == "" {
		jobId = uuid.New().String()
	}
	requestLifecycleMu.Lock()
	inFlightRequestCounter.Add(1)
	requestLifecycleMu.Unlock()
	// log.Printf("Processing Request id=%s inFlightCounter=%v\n", jobId, currentInflight)
	// create a new job
	job := &job{
		jobId:          jobId,
		processingTime: workers.perJobProcessingTime,
		resultChannel:  make(chan struct{}, 1),
	}

	workers.workerPool.addJob(job)

	<-job.resultChannel
	// log.Printf("Request completed id=%s\n", jobId)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("work done!"))

	completedJobs.Add(1)
	inFlightRequestCounter.Add(-1)
}

type Stats struct {
	InFlightRequestCounter     int64 `json:"inFlightRequestCounter"`
	ActiveProcessingCounter    int64 `json:"activeProcessingCounter"`
	MaxActiveProcessingCounter int64 `json:"maxActiveProcessingCounter"`
	CompletedJobs              int64 `json:"completedJobs"`
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	response := Stats{
		InFlightRequestCounter:     inFlightRequestCounter.Load(),
		ActiveProcessingCounter:    activeProcessingCounter.Load(),
		MaxActiveProcessingCounter: maxActiveProcessingCounter.Load(),
		CompletedJobs:              completedJobs.Load(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}

func resetStatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	requestLifecycleMu.Lock()
	defer requestLifecycleMu.Unlock()

	if inFlightRequestCounter.Load() != 0 || activeProcessingCounter.Load() != 0 {
		http.Error(w, "cannot reset stats while requests are in flight", http.StatusConflict)
		return
	}

	maxActiveProcessingCounter.Store(0)
	w.WriteHeader(http.StatusNoContent)
}

func setupLogging() {
	file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v\n", err)
	}
	defer file.Close()

	// Combine the file and terminal outputs
	multiWriter := io.MultiWriter(os.Stdout, file)
	log.SetOutput(multiWriter)
}

var inFlightRequestCounter atomic.Int64
var activeProcessingCounter atomic.Int64
var maxActiveProcessingCounter atomic.Int64
var completedJobs atomic.Int64
var requestLifecycleMu sync.Mutex

func updateMaxActive(current int64) {
	for {
		previous := maxActiveProcessingCounter.Load()
		if current <= previous {
			return
		}
		if maxActiveProcessingCounter.CompareAndSwap(previous, current) {
			return
		}
	}
}

func main() {
	setupLogging()

	inFlightRequestCounter.Store(0)
	activeProcessingCounter.Store(0)
	maxActiveProcessingCounter.Store(0)
	completedJobs.Store(0)

	concurrencyLimit := 8
	w := &workers{
		workerPool: WorkerPool{
			poolSize:      concurrencyLimit,
			jobQueue:      make(chan *job),
			waitOnClosing: &sync.WaitGroup{},
		},
		perJobProcessingTime: 50 * time.Millisecond,
	}

	waitForStart := &sync.WaitGroup{}
	waitForStart.Add(1)
	go w.workerPool.init(waitForStart)
	waitForStart.Wait()

	// setting up a router
	mux := http.NewServeMux()
	// route handler registration
	mux.HandleFunc("/work", w.workHandler)
	mux.HandleFunc("/stats", statsHandler)
	mux.HandleFunc("/stats/reset", resetStatsHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// channel to listen on OS signals
	serverErrors := make(chan error, 1)

	go func() {
		log.Println("Server is running at http://localhost:8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	// channel to listen for interruption signals
	shutdownSignals := make(chan os.Signal, 1)
	signal.Notify(shutdownSignals, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Printf("Server forced to shutdown due to error: %v\n", err)
		return
	case sig := <-shutdownSignals:
		log.Printf("\nRecieved signal: %v. Starting graceful shutdown...\n", sig)

		// create a context w/ timeout for shutdown phase
		// giving time to active requests to complete
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Graceful shutdown failed: %v. Forcing close.\n", err)
			server.Close() // force close on timeout expiry
		}

	}
	w.workerPool.close()

	log.Println("Server cleanly stopped.")
}
