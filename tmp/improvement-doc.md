### Issue 2 — The server doesn't actually track maximum active processing
Your `/stats` endpoint exposes the current active-processing count, not the maximum observed. 

server(1).go

The client compensates by polling `/stats` every millisecond and remembering the highest value. 

client(1).go

That works as an approximation, but it has two drawbacks:

- Polling can miss short-lived peaks.
- Hundreds or thousands of extra HTTP requests can interfere with the benchmark itself.
Better approach: track `maxActive` inside the server using an atomic compare-and-swap operation whenever a worker starts processing a job.

For example:

```
var maxActiveProcessingCounter atomic.Int64

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
```

Then update it inside `doJob`:

```
current := activeProcessingCounter.Add(1)
updateMaxActive(current)

defer activeProcessingCounter.Add(-1)
time.Sleep(j.processingTime)
```

Expose `maxActiveProcessingCounter` through `/stats`.

The client can then read the server's statistics once before and once after a run, rather than polling continuously. For a per-run maximum, either reset the statistic safely between runs or restart the server for each run.


Improvements:
A. Make your latency measurements more precise: record microseconds instead of ms
B. add a few columns to runs.csv:
- timestamp_utc
- throughput_rps
- p50_ms, p95_ms, p99_ms