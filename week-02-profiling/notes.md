go install github.com/rakyll/hey@latest

$(go env GOPATH)/bin/hey \
  -z 60s \
  -c 10 \
  "http://localhost:8080/wait?ms=2000"

PPROF for a duration
curl -o week-02-profiling/results/profiles/wait.goroutine.prof \             
  "http://127.0.0.1:6060/debug/pprof/goroutine?seconds=15"


PPROF snapshot
curl -o week-02-profiling/results/profiles/wait.goroutine.snapshot.prof \
  "http://127.0.0.1:6060/debug/pprof/goroutine"

PPROF ANALYSIS
go tool pprof week-02-profiling/results/profiles/wait.goroutine.snapshot.prof                                  

FLAME GRAPH
go tool pprof -http=127.0.0.1:8081 week-02-profiling/bin/server week-02-profiling/results/profiles/cpu.cpu.prof

Heap Profile baseline:
curl -o week-02-profiling/results/profiles/baseline.heap.prof \
  "http://127.0.0.1:6060/debug/pprof/heap?gc=1"

gc=1 parameter requests a garbage collection before capturing the profile

Inspect allocated versus retained memory
go tool pprof \
  -sample_index=inuse_space \
  week-02-profiling/results/profiles/baseline.heap.prof

inuse_space:	Estimated heap bytes still live at the time of profiling


go tool pprof \
  -sample_index=alloc_space \
  week-02-profiling/results/profiles/baseline.heap.prof

alloc_space:	Estimated cumulative heap bytes allocated over the process's lifetime

