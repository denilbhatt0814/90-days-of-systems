# Do More Goroutines Actually Make Go Programs Faster?

This is a writing scaffold, not a publish-ready article.

## Proposed narrative

1. Open with the question and the distinction between concurrency, parallelism, and capacity.
2. Describe the observation that motivated the experiments.
3. Explain the CPU-bound setup, including worker counts and `GOMAXPROCS`.
4. Show the CPU result and explain where the speedup flattened.
5. Introduce the simulated-I/O setup and what overlapping waits demonstrate.
6. Build the HTTP server/client experiment as a downstream-capacity test.
7. Show throughput and latency as client concurrency crosses the server's eight-worker limit.
8. State the evidence-backed answer and its boundaries.
9. Close with the next investigation and personal reflection.

## Charts to add

- TODO: Chart CPU duration by worker count for each `GOMAXPROCS` setting.
- TODO: Chart simulated-I/O duration by worker count.
- TODO: Chart HTTP throughput and p50/p95/p99 latency by client worker count.
- TODO: Annotate the HTTP chart at the server's eight active processing slots.

## Findings to develop

- TODO: Explain the CPU plateau using the recorded logs.
- TODO: Explain why simulated waits scale differently from CPU work.
- TODO: Explain why HTTP throughput plateaus while latency continues to rise.
- TODO: Cite the raw CSVs and describe the 24-run, 7,200-request dataset.

## Personal reflections

- TODO: Author input — what I expected before running the experiments.
- TODO: Author input — what surprised me.
- TODO: Author input — what I would do differently next time.
