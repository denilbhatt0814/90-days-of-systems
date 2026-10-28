# Week 1 — Do More Goroutines Actually Make Go Programs Faster?

This week investigates when increasing Go concurrency improves completion time or throughput, and when the limiting resource is somewhere else. The experiments follow the repository's working loop: observe, hypothesize, measure, isolate, experiment, and explain.

## Research question

Do more goroutines make a Go program faster? The answer depends on the workload, the number of CPUs available through `GOMAXPROCS`, and downstream capacity.

## Experiments

| Experiment | Question | Source | Historical evidence |
| --- | --- | --- | --- |
| CPU-bound concurrency | How do worker count and `GOMAXPROCS` affect CPU work? | [`experiments/cpu-concurrency/`](./experiments/cpu-concurrency/) | Worker-pool logs and checksums |
| Simulated I/O | Can waiting work overlap as worker count increases? | [`experiments/simulated-io/`](./experiments/simulated-io/) | [`GOMAXPROCS_4_VARIABLE_WORKERS.log`](./experiments/simulated-io/GOMAXPROCS_4_VARIABLE_WORKERS.log) |
| HTTP concurrency and saturation | What happens when a concurrent client meets an eight-worker server? | [`experiments/http-concurrency/`](./experiments/http-concurrency/) | [`results/runs.csv`](./experiments/http-concurrency/results/runs.csv) and [`results/request-latency.csv`](./experiments/http-concurrency/results/request-latency.csv) |

## Prerequisites

- Go 1.27 or newer. The recorded runs used Go `go1.27.2` on Linux/amd64; hardware and scheduler behavior will change the timings.
- A shell with `go`, `gofmt`, and `GOMAXPROCS` available.
- The HTTP experiment needs permission to bind localhost port `8080`.

The repository has a root `go.mod` and uses Go 1.27's standard-library `uuid` package, so it does not need network access or third-party modules.

## Reproduce the experiments

Run commands from each program directory. Each program appends application output to a local `app.log`; those logs are ignored for future runs, while historical tracked logs are retained.

### CPU-bound concurrency

```sh
cd week-01-go-concurrency/experiments/cpu-concurrency/sequential
go run .

cd ../workerpool
GOMAXPROCS=4 go run .

cd ../workerpool-isolate-workers
GOMAXPROCS=4 go run . -workers 4
```

The worker-pool sweep uses `workers={1,2,4,6,8,16,32,64}`, `jobs=300`, and `iterationsPerJob=10000`. The isolated variant runs one worker count per process, making it useful for changing `GOMAXPROCS` between runs. The sequential baseline uses the same 300 jobs and 10,000 iterations per job.

### Simulated I/O

```sh
cd week-01-go-concurrency/experiments/simulated-io
GOMAXPROCS=4 go run . -workers 1
GOMAXPROCS=4 go run . -workers 64
```

This program performs 300 jobs, each waiting for 20 ms with `time.Sleep`. `-workers` controls the worker-pool size; `GOMAXPROCS` is held at 4 in the historical sweep.

### HTTP concurrency and downstream saturation

Start the server in one terminal:

```sh
cd week-01-go-concurrency/experiments/http-concurrency/server
go run .
```

Run the client in another terminal, repeating the command for the worker counts of interest:

```sh
cd week-01-go-concurrency/experiments/http-concurrency/client
GOMAXPROCS=4 go run . -workers 1
GOMAXPROCS=4 go run . -workers 4
GOMAXPROCS=4 go run . -workers 8
GOMAXPROCS=4 go run . -workers 16
GOMAXPROCS=4 go run . -workers 32
GOMAXPROCS=4 go run . -workers 64
```

The server has an eight-worker processing pool and sleeps for 50 ms per job. The client sends 300 requests per run and uses `-workers` as its concurrency limit. It resets and reads the server's `/stats` endpoints around each run.

## Measurements and results

The CPU and simulated-I/O programs print duration and a checksum; the checksum is a simple guard that all jobs completed. The HTTP client records elapsed time, throughput, success/failure counts, maximum observed server processing concurrency, and p50/p95/p99 request latency in [`runs.csv`](./experiments/http-concurrency/results/runs.csv). It records one row per request in [`request-latency.csv`](./experiments/http-concurrency/results/request-latency.csv), including latency in microseconds, status code, success flag, and error text.

The preserved HTTP data contains 24 runs: four repetitions each at client worker counts 1, 4, 8, 16, 32, and 64. Every run has 300 request records, for 7,200 records total. All 7,200 requests returned HTTP 200 and were marked successful. The server reached a maximum of 8 active processing jobs for client worker counts 8 and above.

The HTTP results show the saturation shape clearly: throughput rises from about 19.6 requests/second at one client worker to about 154–156 requests/second at 8 workers, then stays around that level through 64 workers. Median latency is about 51 ms at 8 workers, then about 101 ms, 203 ms, and 405 ms at 16, 32, and 64 workers respectively. More client concurrency therefore adds queueing delay after the eight server processing slots are occupied.

The CPU logs show a different limit. In the 16-CPU historical run with `GOMAXPROCS=8`, duration falls from about 5.3 s at one worker to about 0.85 s at 8 workers, with little further improvement at 16–64 workers. At `GOMAXPROCS=1`, changing worker count leaves duration near 5.35 s. The older 2-CPU logs are retained as historical measurements and are not directly comparable to the 16-CPU runs.

The simulated-I/O sweep goes from about 6.09 s at one worker to about 102 ms at 64 workers for 300 × 20 ms waits. This is expected overlap of blocked waits, not evidence that the CPU performed 64 units of work in parallel.

## Limitations

- These are small local experiments, not production workload models. Timings depend on hardware, OS scheduling, Go version, background load, and localhost networking.
- The CPU workload uses a deliberately simple prime test and repeats identical jobs; it does not represent every CPU-bound workload.
- `time.Sleep` models waiting but not real disk, network, database, or external-service behavior.
- The HTTP server has fixed 50 ms work, a fixed pool of 8, no real downstream service, and no network contention. The client measures end-to-end localhost request latency, including queueing.
- The HTTP client records all request outcomes in the latency file, but the saved runs contain no failures; timeout and error behavior remains untested in the preserved dataset.
- The percentile implementation uses the sorted sample at `int((p/100)*(n-1))`; it is consistent with the recorded summaries but is only one percentile convention.

## Closing artifacts

- [Week 1 retrospective](./notes/week-01-retrospective.md) — evidence-backed notes with author-input placeholders.
- [Blog draft outline](./blog/draft.md) — a writing scaffold, not a finished article.
