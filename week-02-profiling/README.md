# Week 2: Where does the time actually go?

**Status: in progress. M1 profiling foundation: complete.** The week asks how to determine what a slow Go service is actually spending time doing. M1 establishes a baseline using CPU, heap and goroutine profiles alongside request throughput and latency. Further Week 2 investigations remain open.

## Directory map

- `cmd/server/` — unchanged HTTP service source.
- `experiments/m1/` — experiment configuration and reproduction guide.
- `results/m1/profiles/` — recorded raw profiles and goroutine dump.
- `results/m1/benchmarks/` — original load-generator output and commands.
- `results/m1/figures/` — original CPU flame graph screenshot.
- `reports/` — findings and figure interpretation.
- `retrospectives/` — observations, mistakes and lessons.
- `notes.md` — original working notes, with relocated paths updated.

## Build and run

From the repository root, with the Go version required by `go.mod`:

```sh
mkdir -p week-02-profiling/bin
go build -o week-02-profiling/bin/server ./week-02-profiling/cmd/server
./week-02-profiling/bin/server
```

The application listens on `:8080`; pprof listens only on `127.0.0.1:6060`. For example, `curl 'http://localhost:8080/cpu?n=30000'` counts primes, and `curl 'http://localhost:8080/wait?ms=100'` waits on a timer. Defaults are 10000 and 100 respectively. Both parameters must be positive integers; maximums are 100000 and 10000 respectively.

Inspect a recorded profile or capture a new one while load is running:

```sh
go tool pprof -top week-02-profiling/results/m1/profiles/cpu.cpu.prof
curl -o /tmp/m1-cpu.prof 'http://127.0.0.1:6060/debug/pprof/profile?seconds=30'
go tool pprof -http=127.0.0.1:8081 /tmp/m1-cpu.prof
```

See the [reproduction guide](experiments/m1/README.md), [M1 report](reports/M1_PROFILING_FOUNDATION.md), and [retrospective](retrospectives/M1_RETROSPECTIVE.md). The original screenshot is retained; the report annotates its interpretation.
