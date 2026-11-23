# M1 retrospective

## Observations

The CPU profile identified `isPrime` as the dominant sampled hotspot. The waiting workload had noticeable latency with little sampled CPU; the goroutine snapshot supplied complementary evidence that handlers were sleeping. Heap retained and cumulative allocation views answered different questions.

## Mistakes and limitations

The original evidence was mixed in one directory, and the capture timeline, machine configuration, tool versions and profile-to-load correspondence were not fully recorded. That makes exact reproduction and causal comparisons weaker. The benchmark commands survive in their output; the separate long-wait command survives in the working notes. The screenshot is preserved without assuming an independently recorded annotation.

## Lessons

Record commands, versions, hardware, concurrency, workload, timestamps and capture overlap at collection time. Keep raw evidence alongside its interpretation and retain original figures. Compare matching workloads and use CPU, goroutine and heap profiles together: a slow response alone does not identify CPU work, and a single snapshot cannot describe an entire run. Maintain this unchanged baseline for later experiments.
