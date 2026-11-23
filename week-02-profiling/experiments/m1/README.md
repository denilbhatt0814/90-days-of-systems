# M1 reproduction and configuration

Run commands from the repository root. Build and start the server as described in the [Week 2 README](../../README.md). Install `hey` separately if needed (`go install github.com/rakyll/hey@latest`); its original version and the original machine configuration are unknown. Use separate terminals for the server, load and profile capture. Do not overwrite recorded evidence: the examples below write disposable captures to `/tmp`.

The commands recovered from the recorded benchmark outputs are:

```sh
$(go env GOPATH)/bin/hey -z 25s -c 4 'http://localhost:8080/cpu?n=30000'
$(go env GOPATH)/bin/hey -z 25s -c 4 'http://localhost:8080/wait?ms=100'
```

The original notes also record a distinct goroutine experiment: `hey -z 60s -c 10 'http://localhost:8080/wait?ms=2000'`. Run this separately; it is not the 100 ms benchmark. Its notes show a 15-second goroutine delta and a snapshot. Capture these while that load is active:

```sh
curl -o /tmp/m1-wait-delta.prof 'http://127.0.0.1:6060/debug/pprof/goroutine?seconds=15'
curl -o /tmp/m1-wait-snapshot.prof 'http://127.0.0.1:6060/debug/pprof/goroutine'
curl -o /tmp/m1-wait-goroutines.txt 'http://127.0.0.1:6060/debug/pprof/goroutine?debug=2'
```

For a fresh CPU capture, start a load run long enough to overlap the entire capture. The following is an illustrative 30-second capture, not a claim about the original capture's timing or load:

```sh
curl -o /tmp/m1-cpu.prof 'http://127.0.0.1:6060/debug/pprof/profile?seconds=30'
curl -o /tmp/m1-heap.prof 'http://127.0.0.1:6060/debug/pprof/heap?gc=1'
go tool pprof -top /tmp/m1-cpu.prof
go tool pprof -sample_index=inuse_space -top /tmp/m1-heap.prof
go tool pprof -sample_index=alloc_space -top /tmp/m1-heap.prof
go tool pprof -top /tmp/m1-wait-snapshot.prof
```

`gc=1` requests garbage collection before the heap snapshot. The original capture timing, exact overlap between benchmark and profiles, machine characteristics, Go and hey versions, and configuration for the before/after heap captures were not recorded. These commands show how to repeat the method, not guarantee identical results.

## Visualize

```sh
go tool pprof -http=127.0.0.1:8081 week-02-profiling/results/m1/profiles/cpu.cpu.prof
```

Open the local UI and choose **View → Flame Graph**; Graphviz may be needed for other views. The original notes supplied the built server binary as an additional argument; rebuild it and place `week-02-profiling/bin/server` before the profile path if symbols need help. A rebuild may not match the historical binary exactly. Preserve the existing screenshot; any newly exported or annotated visualization is a new artifact, not a replacement for the original.
