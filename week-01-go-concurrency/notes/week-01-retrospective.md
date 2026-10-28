# Week 1 retrospective

## What I set out to learn

The central question was: **Do more goroutines actually make Go programs faster?** The experiments separated CPU-bound work, waiting work, and a client/server system with a constrained downstream worker pool.

## What I built

- A sequential and worker-pool prime-number benchmark.
- A `time.Sleep`-based worker pool for simulated I/O.
- A local HTTP server with eight processing workers and 50 ms simulated processing time.
- An HTTP client that varies its worker count and records run-level and request-level results.

## Initial hypotheses

- CPU-bound work should improve until available CPU parallelism is used; additional workers should then provide little benefit.
- Waiting work should overlap as worker count increases, with `GOMAXPROCS` less important while goroutines are blocked.
- HTTP throughput should stop improving once the server's eight processing workers are saturated, while queueing should increase request latency.

## What the measurements showed

The CPU logs support a `GOMAXPROCS`-dependent ceiling: in the 16-CPU run at `GOMAXPROCS=8`, duration improved from about 5.3 s at one worker to about 0.85 s at 8 workers, then flattened through 64 workers. At `GOMAXPROCS=1`, worker count did not materially change the roughly 5.35 s duration.

The simulated-I/O sweep improved from about 6.09 s at one worker to about 102 ms at 64 workers for 300 × 20 ms waits.

The HTTP dataset contains 24 runs and 7,200 request records. Each of six client worker counts has four repetitions and 300 records per run. All requests succeeded. Throughput rose from about 19.6 requests/second at one client worker to about 154–156 requests/second at 8 workers and stayed there through 64 workers. Maximum server processing concurrency stayed at 8. Median latency increased from about 51 ms at 8 client workers to about 405 ms at 64.

## What surprised me

TODO: Author input

## What I understand better now

The measurements distinguish goroutine concurrency from CPU parallelism and from downstream capacity. More goroutines can overlap waits, but they cannot make a one-processor CPU execute more CPU work simultaneously, and they cannot increase a server's fixed eight processing slots.

## Limitations and unresolved questions

These are synthetic, local workloads. The CPU and sleep models omit real application behavior; the HTTP service has fixed work and no network or dependency failures. The historical CPU logs were collected on different reported CPU counts, and the HTTP runs were sequential and all successful. The experiment does not establish behavior under real remote I/O, variable service time, retries, or multiple server instances.

## What I would investigate next

TODO: Author input
