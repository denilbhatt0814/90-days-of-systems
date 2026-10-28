package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"sync"
	"time"
	"uuid"
)

type dependencies struct {
	httpClient             *http.Client
	requestStatsAggregator chan requestStats
}

type serverStats struct {
	ActiveProcessingCounter    int64 `json:"activeProcessingCounter"`
	MaxActiveProcessingCounter int64 `json:"maxActiveProcessingCounter"`
}

type runRecord struct {
	runID           string
	workers         int
	gomaxprocs      int
	jobs            int
	timestampUTC    string
	elapsedMS       int64
	throughputRPS   float64
	successes       uint64
	failures        uint64
	serverMaxActive int64
	p50MS           float64
	p95MS           float64
	p99MS           float64
}

type requestStats struct {
	requestID  string
	latencyUS  int64
	statusCode int
	success    bool
	err        string
}

func externalServiceCall(deps dependencies) bool {
	requestID := uuid.New().String()
	startTime := time.Now()
	record := requestStats{requestID: requestID}
	defer func() {
		record.latencyUS = time.Since(startTime).Microseconds()
		deps.requestStatsAggregator <- record
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/work", nil)
	if err != nil {
		record.err = err.Error()
		log.Printf("Failed to create request requestID=%s error=%v", requestID, err)
		return false
	}
	request.Header.Set("x-trace-id", requestID)

	resp, err := deps.httpClient.Do(request)
	if err != nil {
		record.err = err.Error()
		log.Printf("Service call failed requestID=%s error=%v", requestID, err)
		return false
	}
	defer resp.Body.Close()
	record.statusCode = resp.StatusCode

	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil {
		record.err = err.Error()
		log.Printf("Failed to read response requestID=%s error=%v", requestID, err)
		return false
	}

	if resp.StatusCode != http.StatusOK {
		record.err = http.StatusText(resp.StatusCode)
		log.Printf("Unexpected status code: %d requestID=%s\n", resp.StatusCode, requestID)
		return false
	}

	record.success = true
	// log.Printf("Service call success requestID=%s latency=%v\n", requestID, time.Since(startTime))
	return true
}

func ioWork(deps dependencies) uint64 {
	if externalServiceCall(deps) {
		return 1
	}
	return 0
}

func runExperiment(jobs, workers int, deps dependencies) (uint64, time.Duration) {
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
				results <- ioWork(deps)
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

func readServerStats(client *http.Client) (serverStats, error) {
	request, err := http.NewRequest(http.MethodGet, "http://localhost:8080/stats", nil)
	if err != nil {
		return serverStats{}, err
	}

	resp, err := client.Do(request)
	if err != nil {
		return serverStats{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return serverStats{}, io.ErrUnexpectedEOF
	}

	var stats serverStats
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return serverStats{}, err
	}
	return stats, nil
}

func resetServerStats(client *http.Client) error {
	request, err := http.NewRequest(http.MethodPost, "http://localhost:8080/stats/reset", nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return io.ErrUnexpectedEOF
	}
	return nil
}

func resultsFilePath() string {
	_, sourceFile, _, ok := runtime.Caller(0)
	if ok {
		return filepath.Join(filepath.Dir(sourceFile), "..", "results", "runs.csv")
	}
	return filepath.Join("results", "runs.csv")
}

func writeRunRecord(record runRecord) error {
	path := resultsFilePath()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	fileInfo, err := os.Stat(path)
	fileExists := err == nil && fileInfo.Size() > 0
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	if !fileExists {
		if err := writer.Write([]string{
			"run_id", "workers", "gomaxprocs", "jobs", "timestamp_utc",
			"elapsed_ms", "throughput_rps", "successes", "failures",
			"server_max_active", "p50_ms", "p95_ms", "p99_ms",
		}); err != nil {
			return err
		}
	}

	if err := writer.Write([]string{
		record.runID,
		strconv.Itoa(record.workers),
		strconv.Itoa(record.gomaxprocs),
		strconv.Itoa(record.jobs),
		record.timestampUTC,
		strconv.FormatInt(record.elapsedMS, 10),
		strconv.FormatFloat(record.throughputRPS, 'f', 3, 64),
		strconv.FormatUint(record.successes, 10),
		strconv.FormatUint(record.failures, 10),
		strconv.FormatInt(record.serverMaxActive, 10),
		strconv.FormatFloat(record.p50MS, 'f', 3, 64),
		strconv.FormatFloat(record.p95MS, 'f', 3, 64),
		strconv.FormatFloat(record.p99MS, 'f', 3, 64),
	}); err != nil {
		return err
	}
	writer.Flush()
	return writer.Error()
}

func requestLatencyFilePath() string {
	_, sourceFile, _, ok := runtime.Caller(0)
	if ok {
		return filepath.Join(filepath.Dir(sourceFile), "..", "results", "request-latency.csv")
	}
	return filepath.Join("results", "request-latency.csv")
}

func saveRequestLatencyRecords(runID string, records []requestStats) error {
	path := requestLatencyFilePath()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	fileInfo, err := os.Stat(path)
	fileExists := err == nil && fileInfo.Size() > 0
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	if !fileExists {
		if err := writer.Write([]string{
			"run_id", "request_id", "latency_us", "status_code", "success",
			"error",
		}); err != nil {
			return err
		}
	}

	for _, record := range records {
		if err := writer.Write([]string{
			runID,
			record.requestID,
			strconv.FormatInt(record.latencyUS, 10),
			strconv.Itoa(record.statusCode),
			strconv.FormatBool(record.success),
			record.err,
		}); err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}

func percentileMS(records []requestStats, percentile float64) float64 {
	if len(records) == 0 {
		return 0
	}

	latencies := make([]int64, len(records))
	for i, record := range records {
		latencies[i] = record.latencyUS
	}
	sort.Slice(latencies, func(i, j int) bool {
		return latencies[i] < latencies[j]
	})

	index := int((percentile / 100) * float64(len(latencies)-1))
	return float64(latencies[index]) / 1000
}

func collectRequestStats(requestStatsAggregator <-chan requestStats, jobs int) []requestStats {
	requestStats := make([]requestStats, jobs)
	for i := range requestStats {
		requestStats[i] = <-requestStatsAggregator
	}
	return requestStats
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
	deps := dependencies{
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 100,
			},
		},
		requestStatsAggregator: make(chan requestStats, jobs),
	}

	var workers int
	flag.IntVar(&workers, "workers", 1, "Number of workers in workerpool")
	flag.Parse()

	runID := uuid.New().String()
	runTimestamp := time.Now().UTC()
	log.Printf("STARTING RUN (I/O Workload) runID=%v", runID)
	gomaxprocs := runtime.GOMAXPROCS(0)
	log.Printf("CPUs: %d | GOMAXPROCS: %d\n", runtime.NumCPU(), gomaxprocs)

	if err := resetServerStats(deps.httpClient); err != nil {
		log.Fatalf("failed to reset server stats: %v", err)
	}

	checksum, elapsed := runExperiment(jobs, workers, deps)
	requestStats := collectRequestStats(deps.requestStatsAggregator, jobs)

	stats, err := readServerStats(deps.httpClient)
	if err != nil {
		log.Fatalf("failed to read server stats: %v", err)
	}

	failures := uint64(jobs) - checksum
	log.Printf("workers=%2d duration=%s checksum=%d\n", workers, elapsed, checksum)

	record := runRecord{
		runID:           runID,
		workers:         workers,
		gomaxprocs:      gomaxprocs,
		jobs:            jobs,
		timestampUTC:    runTimestamp.Format(time.RFC3339Nano),
		elapsedMS:       elapsed.Milliseconds(),
		throughputRPS:   float64(jobs) / elapsed.Seconds(),
		successes:       checksum,
		failures:        failures,
		serverMaxActive: stats.MaxActiveProcessingCounter,
		p50MS:           percentileMS(requestStats, 50),
		p95MS:           percentileMS(requestStats, 95),
		p99MS:           percentileMS(requestStats, 99),
	}
	if err := writeRunRecord(record); err != nil {
		log.Printf("failed to write run results: %v", err)
	}

	if err := saveRequestLatencyRecords(runID, requestStats); err != nil {
		log.Printf("failed to write request latency results: %v", err)
	}

	log.Printf("RUN ENDED (I/O Workload) runID=%v", runID)
}
