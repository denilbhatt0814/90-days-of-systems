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