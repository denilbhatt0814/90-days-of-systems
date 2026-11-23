# M1: Profiling foundation

**Complete.** The baseline contrasts deliberately inefficient prime counting with timer-based waiting. CPU samples describe where the process executed on CPU, while request latency includes time spent waiting. The two workloads and their parameter values differ, so their throughput and latency are descriptive observations, not a controlled speed comparison.

| Recorded metric | CPU workload | Waiting workload |
| --- | ---: | ---: |
| Sampled CPU time | 29.99 s | 0.15 s |
| Throughput | 12.92 RPS | 39.54 RPS |
| p95 latency | 372.9 ms | 101.5 ms |

The [CPU benchmark output](../results/m1/benchmarks/hey-cpu-prof-out-1.txt) records `n=30000`, 25 seconds and concurrency 4 (12.9273 RPS, 327 HTTP 200 responses). The [wait output](../results/m1/benchmarks/hey-wait-prof-out-1.txt) records `ms=100`, the same duration and concurrency (39.5413 RPS, 992 HTTP 200 responses). Profile capture timing and exact overlap with these runs are unknown. CPU sample totals are accumulated CPU time, not request latency or necessarily wall-clock duration.

## Evidence and interpretation

- The [CPU profile](../results/m1/profiles/cpu.cpu.prof) attributes 83.79% of samples flat to `isPrime`, supporting the conclusion that prime checking dominates sampled CPU execution in that capture. Flat means samples directly in the function, excluding its callees.
- The [idle CPU profile](../results/m1/profiles/idle.cpu.prof) recorded zero samples. That does not prove the process used absolutely no CPU; sampling can miss small amounts of activity.
- The [waiting CPU profile](../results/m1/profiles/wait.cpu.prof) contains only 0.15 s of CPU samples despite request latency. The [goroutine snapshot](../results/m1/profiles/wait.goroutine.snapshot.prof) recorded 24 goroutines, 10 sleeping in `handleWaitTask`. The separate [text dump](../results/m1/profiles/wait.goroutines.txt) shows those sleeping stacks; the notes describe a separate 2000 ms, concurrency-10 run. A snapshot is one instant, not a complete history.
- The [baseline heap profile](../results/m1/profiles/baseline.heap.prof) records approximately 512.05 KiB `inuse_space` and 5990.45 KiB `alloc_space`. These are sampled estimates: retained live heap and cumulative allocations respectively, not interchangeable measures or per-request allocation counts. The [before and after profiles](../results/m1/profiles/) are also preserved; their workload and timing are unknown.

## Figure annotation

The [original CPU flame graph screenshot](../results/m1/figures/cpu-prof-flamegraph.png) is preserved as captured. Read width as the relative share of sampled CPU in a stack; `isPrime` is the prominent hotspot, consistent with its 83.79% flat share. A CPU flame graph does not show time spent sleeping. See the [reproduction guide](../experiments/m1/README.md#visualize) for regenerating the view from the raw profile.

## Limits

The hardware, CPU limits, original tool versions, capture scheduling and exact relation between each profile and benchmark are unknown. The recorded runs do not establish a general capacity limit or improvement. Preserve this baseline before changing the workload; later experiments need matched conditions and independently recorded parameters.
