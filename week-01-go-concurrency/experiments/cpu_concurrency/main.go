package main

import (
	"fmt"
	"runtime"
	"time"
)

// We shall assume finding isPrime is a highly CPU intensive task
func isPrime(n uint64) bool {
	if n <= 1 {
		return false
	}

	for i := uint64(2); i < n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func cpuWork(iterations uint64) uint64 {
	// our job is to find sum of prime numbers till n = iterations
	var sum uint64 = 0
	for i := uint64(0); i < iterations; i++ {
		if isPrime(i) {
			sum += i
		}
	}

	return sum
}

func main() {
	const jobs = 300
	const iterationsPerJob = 10000

	// just to verify each job produces same expected results
	var checksum uint64

	start := time.Now()
	for i := 0; i < jobs; i++ {
		checksum += cpuWork(iterationsPerJob)
	}
	elapsed := time.Since(start)

	fmt.Printf("CPUs reported by Go: %d\n", runtime.NumCPU())
	fmt.Printf("GOMAXPROCS: %d\n", runtime.GOMAXPROCS(0)) // returns the previous setting of max usable CPUs
	fmt.Printf("Jobs: %d\n", jobs)
	fmt.Printf("Iterations per job: %d\n", iterationsPerJob)
	fmt.Printf("Checksum: %d\n", checksum)
	fmt.Printf("Duration: %s\n", elapsed)
}
