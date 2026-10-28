# HTTP concurrency and downstream saturation

## Purpose

Measure what happens when a concurrent HTTP client sends work to a local server whose processing capacity is deliberately limited to eight workers.

## Setup and prerequisites

Go 1.27 or newer is required. Run the server and client on the same machine. The server binds `localhost:8080`; ensure that port is available.

## Run commands

Terminal 1:

```sh
cd week-01-go-concurrency/experiments/http-concurrency/server
go run .
```

Terminal 2, with the server running:

```sh
cd week-01-go-concurrency/experiments/http-concurrency/client
GOMAXPROCS=4 go run . -workers 1
GOMAXPROCS=4 go run . -workers 8
GOMAXPROCS=4 go run . -workers 64
```

Repeat the client command for each worker count you want to compare. Stop the server with Ctrl-C after the runs.

## Configuration

The client sends 300 requests per run and accepts `-workers` as its concurrency limit. The historical data uses client workers `{1,4,8,16,32,64}` and `GOMAXPROCS=4`, with four repetitions per setting. The server uses an eight-worker pool and sleeps 50 ms per accepted job. The HTTP client and server timeouts are 5 seconds.

## Metrics collected

The client records elapsed time, throughput, successes, failures, server maximum active processing, and p50/p95/p99 latency per run. It also writes one request-level row with run ID, request ID, latency in microseconds, status code, success flag, and error text. The server exposes `/stats` and `/stats/reset` for benchmark coordination.

## Results location

The raw run summaries are [`results/runs.csv`](./results/runs.csv). The raw request records are [`results/request-latency.csv`](./results/request-latency.csv). Existing application logs beside the server, client, and HTTP experiment are retained as historical artifacts.

## Known limitations

This is localhost traffic with fixed 50 ms work, a fixed eight-slot server pool, and no real downstream dependency. The benchmark does not test network loss, server errors, timeouts, retries, or multiple server instances. Its percentile calculation is a direct sorted-sample index, and the four repetitions are sequential rather than randomized.
