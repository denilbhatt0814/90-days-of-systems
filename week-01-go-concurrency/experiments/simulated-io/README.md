# Simulated I/O concurrency

## Purpose

Show how worker-pool concurrency overlaps waiting work when each job is modeled as a 20 ms `time.Sleep`.

## Setup and prerequisites

From the repository root:

```sh
cd week-01-go-concurrency/experiments/simulated-io
go build .
```

Go 1.27 or newer is required. No external services are required.

## Run commands

```sh
GOMAXPROCS=4 go run . -workers 1
GOMAXPROCS=4 go run . -workers 8
GOMAXPROCS=4 go run . -workers 64
```

## Configuration

The fixed workload is 300 jobs, each with a 20 ms simulated external-service wait. `-workers` sets the worker-pool size. The historical sweep held `GOMAXPROCS=4` and tested workers `{1,2,4,8,16,32,64}`.

## Metrics collected

Each run records CPU count, `GOMAXPROCS`, worker count, elapsed duration, and a completion checksum of 300.

## Results location

The preserved sweep is [`GOMAXPROCS_4_VARIABLE_WORKERS.log`](./GOMAXPROCS_4_VARIABLE_WORKERS.log). New output is appended to `app.log` beside the program.

## Known limitations

`time.Sleep` is only a model of waiting. It does not include real I/O variability, bandwidth, connection pools, remote queues, retries, or failures. The result demonstrates overlap of waits rather than a general performance guarantee.
