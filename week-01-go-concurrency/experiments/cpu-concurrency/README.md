# CPU-bound concurrency

## Purpose

Measure how a CPU-bound prime-number workload changes with worker count and `GOMAXPROCS`, using both a sequential baseline and worker pools.

## Setup and prerequisites

From the repository root, use Go 1.27 or newer:

```sh
cd week-01-go-concurrency/experiments/cpu-concurrency
go build ./sequential ./workerpool ./workerpool-isolate-workers
```

The root module makes the separate experiment directories buildable together, and Go's standard library supplies UUIDs for run IDs. No external services are required.

## Run commands

```sh
cd sequential && go run .
cd ../workerpool && GOMAXPROCS=4 go run .
cd ../workerpool-isolate-workers && GOMAXPROCS=4 go run . -workers 4
```

## Configuration

All variants use 300 jobs and 10,000 prime-search iterations per job. `workerpool` sweeps workers `{1,2,4,6,8,16,32,64}` in one process. `workerpool-isolate-workers` accepts `-workers` and runs one setting per process. Set `GOMAXPROCS` in the environment to isolate scheduler parallelism.

## Metrics collected

The programs print and append CPU count, `GOMAXPROCS`, worker count, checksum, and elapsed duration. The checksum should remain `1720918800` for the recorded workload.

## Results location

Historical logs remain beside the source files, including older runs from both 2-CPU and 16-CPU environments. New runs append to `app.log` in the current program directory.

## Known limitations

The prime test is intentionally inefficient and synthetic. Results are sensitive to CPU model, thermal state, background load, Go version, and the difference between changing `GOMAXPROCS` in-process versus starting a fresh process.
