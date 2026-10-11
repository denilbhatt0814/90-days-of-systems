# 90 Days of Systems

### Build. Break. Measure. Understand. Share.

**Duration:** \~90 days\
**Cadence:** \~14 hours/week\
**Primary language:** Go\
**Learning mode:** Theory → Build → Break → Measure → Understand → Share\
**Primary outcome:** Become a stronger systems/performance engineer capable of investigating unfamiliar systems, identifying bottlenecks, reasoning about distributed failures, and making evidence-backed engineering decisions.

---

# 0. Why This Exists

This is not a generic "learn distributed systems" course.

The objective is to develop a specific engineering capability:

> **Enter an unfamiliar system, understand how it behaves, identify what limits it, prove the cause with evidence, improve it, and explain the trade-offs.**

The 90-day journey is therefore built around engineering investigations rather than passive learning.

We will repeatedly practice:

```text
Observe
   ↓
Form hypotheses
   ↓
Measure
   ↓
Isolate
   ↓
Experiment
   ↓
Fix
   ↓
Benchmark
   ↓
Explain
```

The goal is not to know every technology.

The goal is to develop **engineering judgment**.

---

# 1. North-Star Outcome

By the end of these 90 days, I should be able to look at a system like:

```text
                    ┌──────────────┐
                    │    Client    │
                    └──────┬───────┘
                           │
                    ┌──────▼───────┐
                    │ Load Balancer│
                    └──────┬───────┘
                           │
                    ┌──────▼───────┐
                    │   Go API     │
                    └──┬───┬───┬───┘
                       │   │   │
                 ┌─────┘   │   └──────┐
                 ▼          ▼          ▼
              Redis       Postgres   External API
                 │
                 ▼
               Queue
                 │
                 ▼
              Workers
```

and systematically ask:

- Where is the latency coming from?
- Where is the CPU going?
- Where is memory going?
- What happens when traffic increases 10×?
- What happens when a dependency becomes slow?
- What happens when a worker dies?
- What happens when a message is delivered twice?
- What happens when the DB becomes the bottleneck?
- What happens when Redis fails?
- What happens when the external API times out?
- Where does backpressure occur?
- What state must be durable?
- What can be retried?
- What must be idempotent?
- What can be asynchronous?
- What is the actual bottleneck?
- What evidence proves it?

---

# 2. The Engineering Philosophy

## 2.1 Evidence over intuition

We are allowed to have hypotheses.

We are not allowed to confuse hypotheses with facts.

Bad:

> CPU is high, therefore GC is the problem.

Better:

> CPU is high. I suspect GC may be contributing. Let's profile CPU and allocation behaviour.

---

## 2.2 Bottleneck before optimization

Before changing architecture:

```text
What is limiting throughput?
```

Before increasing resources:

```text
What resource is saturated?
```

Before adding caching:

```text
Is the database actually the bottleneck?
```

Before increasing concurrency:

```text
Is the workload actually parallelizable?
```

Before introducing Kubernetes:

```text
What problem does Kubernetes solve here?
```

---

## 2.3 Understand the mechanism

We don't want:

> "Redis is faster than Postgres."

We want:

> "This workload repeatedly reads relatively stable objects, and database access is contributing X ms/request. Local/Redis caching reduces repeated database work, but introduces stale-data and invalidation trade-offs."

---

## 2.4 Technology follows the problem

We will not learn:

- Kafka because Kafka is popular
- Kubernetes because Kubernetes appears in job descriptions
- Rust because Rust is fast
- another cloud certification because certifications look good

We learn a technology when the problem requires it.

---

# 3. Weekly Operating Model

Every week follows the same structure.

## Monday — Learn

Understand the minimum theory required.

## Tuesday — Build

Create a small implementation or experiment.

## Wednesday — Break

Change conditions deliberately.

Introduce:

- more load
- contention
- failures
- latency
- resource limits
- duplicates
- crashes
- malformed input
- dependency failures

## Thursday — Measure

Collect evidence.

Examples:

- throughput
- latency
- CPU
- memory
- allocations
- GC
- goroutines
- queue depth
- connection pool usage
- DB latency
- error rate

## Friday — Investigate

Turn observations into hypotheses and test them.

## Saturday — Deep Work

Spend the largest block of time implementing, debugging, profiling or improving the system.

## Sunday — Document & Reflect

Produce the week's artifact and technical write-up.

---

# 4. Weekly Time Budget

Target: **14 hours/week**

| Activity                |   Hours |
| ----------------------- | ------: |
| Theory / reading        |      3h |
| Implementation          |      4h |
| Experiments / debugging |      3h |
| Measurement / analysis  |      2h |
| Documentation / blog    |      1h |
| Weekly retrospective    |      1h |
| **Total**               | **14h** |

This is a target, not a prison.

If an experiment becomes unusually interesting, engineering work takes priority over rigid time allocation.

---

# 5. Definition of Done

A week is **not complete because the videos/readings are finished.**

A week is complete when:

1. I understand the relevant concepts.
2. I built something or ran a meaningful experiment.
3. I deliberately changed conditions or broke something.
4. I collected evidence.
5. I can explain what happened.
6. I can explain **why** it happened.
7. I documented the investigation.
8. I produced a tangible artifact.
9. I can answer the week's central question without relying on memorized definitions.

---

# 6. 90-Day Curriculum

---

# PHASE I — UNDERSTAND THE MACHINE

## Week 1 — Do More Goroutines Actually Make Go Programs Faster?

### Central Question

> **When does increasing concurrency improve performance, and when does it make a Go program slower?**

### Why This Matters

Concurrency is one of the most important tools in Go, but blindly adding goroutines is a common performance mistake.

We want to understand what actually happens underneath.

### Concepts

- Concurrency vs parallelism
- CPU-bound vs I/O-bound workloads
- Goroutines
- OS threads
- Go scheduler
- `GOMAXPROCS`
- Context switching
- Worker pools
- Synchronization
- Contention
- Atomic operations
- Mutexes
- Channels
- Goroutine overhead

### Experiment

Build a CPU-bound workload.

Compare:

```text
Sequential
1 worker
2 workers
4 workers
8 workers
16 workers
32 workers
64 workers
```

Then vary `GOMAXPROCS`.

Measure:

- execution time
- throughput
- CPU utilization
- memory
- allocations
- goroutine count

Then create an I/O-bound workload and repeat.

### Break It

Find a case where:

> More goroutines make the program slower.

### Artifact

- Go benchmark
- benchmark table
- performance graph
- explanation of results

### Blog

**"I Thought More Goroutines Would Make My Go Program Faster. I Was Wrong."**

### Definition of Done

I can explain:

> Why 100 goroutines may improve an I/O workload but make a CPU-bound workload slower.

---

# Week 2 — Where Does the Time Actually Go?

## Central Question

> **When a Go service becomes slow, how do we determine what it is actually spending time doing?**

### Concepts

- Latency
- Throughput
- CPU utilization
- CPU saturation
- Memory
- Heap
- Stack
- Allocation
- Escape analysis
- Garbage collection
- GC pressure
- Flame graphs
- `pprof`
- CPU profiles
- Heap profiles
- Allocation profiles
- Goroutine profiles
- Mutex/block profiles

### Experiment

Create several intentionally problematic programs:

1. CPU-heavy
2. Allocation-heavy
3. GC-heavy
4. Lock-contention-heavy
5. Goroutine-heavy

Profile each.

### Key Investigation

Given:

```text
p95 latency: 100ms → 2s
CPU: 40% → 95%
Memory: unchanged
```

determine possible causes and prove/disprove them.

### Break It

Create a service where latency increases without a large increase in traffic.

Use profiling to find the cause.

### Artifact

Annotated flame graphs + profiling report.

### Blog

**"CPU Is at 95%. But What Is the CPU Actually Doing?"**

### Definition of Done

Given high CPU, I can systematically determine whether the CPU is being consumed by:

- business logic
- JSON serialization
- allocations/GC
- locks
- logging
- compression
- scheduler overhead
- inefficient algorithms
- something else

without guessing.

---

# Week 3 — Follow a Request Across the Network

## Central Question

> **Where does latency come from when a request travels through multiple network components?**

### Concepts

- DNS
- IP
- TCP
- SYN/SYN-ACK/ACK
- Connection establishment
- TCP connection reuse
- Keep-alive
- TLS
- HTTP/1.1
- HTTP/2
- Connection pools
- Load balancers
- Timeouts
- Retries
- Network latency
- Packet loss
- Head-of-line blocking

### Experiment

Build:

```text
Client → Go API → Go API → DB/mock dependency
```

Measure:

- new connection per request
- connection reuse
- keep-alive
- different concurrency levels
- artificial network latency

### Break It

Introduce:

- 50ms latency
- 200ms latency
- slow connections
- timeouts
- dropped requests

Observe how latency propagates.

### Artifact

Request lifecycle diagram showing where time is spent.

### Blog

**"A 200ms Dependency Can Turn Into Seconds of API Latency."**

### Definition of Done

I can break down a request into:

```text
DNS
TCP
TLS
HTTP
Application
Dependency
Response
```

and reason about the latency contribution of each.

---

# Week 4 — How Do We Know What Is Happening?

## Central Question

> **How do you investigate a production system when you don't know where the problem is?**

### Concepts

- Observability
- Metrics
- Logs
- Traces
- RED
- USE
- Golden signals
- Histograms
- Percentiles
- p50/p95/p99
- Cardinality
- Correlation
- Structured logging
- Profiling
- Alerting

### Experiment

Instrument a Go service with:

- request count
- error count
- request latency
- CPU
- memory
- goroutines
- GC metrics
- dependency latency
- DB latency

### Investigation Scenario

Simulate:

```text
Normal:
p95 = 150ms

Incident:
p95 = 2.5s
```

Determine where the additional latency came from.

### Break It

Remove one observability signal at a time.

Ask:

> What can I still infer?

### Artifact

A complete observability dashboard + incident investigation.

### Blog

**"Where Did Those 2.5 Seconds Go?"**

### Definition of Done

I can approach an unfamiliar performance incident using:

```text
Symptom
→ Scope
→ Correlation
→ Hypothesis
→ Evidence
→ Root cause
```

rather than immediately changing infrastructure.

---

# PHASE II — UNDERSTAND DISTRIBUTED FAILURE

# Week 5 — What Happens When Work Is Asynchronous?

## Central Question

> **How do queues and workers turn one synchronous system into a distributed workflow?**

### Concepts

- Producer
- Queue
- Consumer
- Worker
- Acknowledgement
- Visibility timeout
- Polling
- Long polling
- Queue depth
- Worker concurrency
- Throughput
- Backpressure
- At-most-once
- At-least-once
- Duplicate delivery
- Poison messages
- Dead-letter queues

### Experiment

Build:

```text
Producer
   ↓
Queue
   ↓
Workers
   ↓
Database
```

Implement:

- retries
- acknowledgement
- worker concurrency
- failed jobs
- DLQ

### Break It

Kill workers:

- before processing
- during processing
- after processing but before acknowledgement

### Questions

What happens?

When does work get duplicated?

When is work lost?

### Artifact

Queue/worker implementation + failure matrix.

### Blog

**"Your Queue Doesn't Guarantee Your Job Runs Exactly Once."**

### Definition of Done

I can explain precisely why:

> **At-least-once delivery does not mean exactly-once execution.**

---

# Week 6 — What Happens When Things Fail Between Systems?

## Central Question

> **How do you safely handle failures when you don't know whether the other system received your request?**

### Concepts

- Timeouts
- Retries
- Exponential backoff
- Jitter
- Idempotency
- Idempotency keys
- Deduplication
- Circuit breakers
- Bulkheads
- Rate limiting
- Backpressure
- Outbox pattern
- Reconciliation
- Partial failure
- Uncertain outcomes

### Experiment

Build:

```text
Worker
   ↓
External API
```

Simulate:

```text
Request succeeds
Request times out
Response is lost
Server crashes
Dependency becomes slow
```

### Critical Scenario

```text
Worker → Payment Provider
             ↓
          SUCCESS

Worker crashes before receiving response.
```

What happens when the message is retried?

### Artifact

Failure-state matrix:

| Failure | Retry? | Idempotency required? | Risk |
| ------- | ------ | --------------------- | ---- |

### Blog

**"The Request Timed Out. Did It Actually Fail?"**

### Definition of Done

I understand why distributed systems cannot safely infer:

> "I didn't receive success, therefore the operation didn't happen."

And I can design for that uncertainty.

---

# Week 7 — Who Owns the Truth?

## Central Question

> **How do databases, caches and replicas behave when different components don't agree immediately?**

### Concepts

- Transactions
- ACID
- Isolation levels
- Locks
- MVCC
- Read/write consistency
- Replication
- Read replicas
- Eventual consistency
- Strong consistency
- Caching
- Cache-aside
- Write-through
- Write-behind
- Cache invalidation
- TTL
- Hot keys
- Stampede
- Single-flight
- CAP theorem

### Experiment

Build:

```text
API
 ↓
Redis
 ↓
Postgres
```

Implement:

1. cache-aside
2. write-through
3. stale data
4. cache failure
5. cache stampede

### Break It

Cause:

- Redis failure
- stale cache
- concurrent cache misses
- DB slowdown

### Artifact

Consistency + caching decision matrix.

### Blog

**"Caching Didn't Make My System Faster. It Made the Failure Modes More Interesting."**

### Definition of Done

I can answer:

> When should this data be cached?

and explain the consistency trade-off rather than simply saying:

> "Redis is faster."

---

# Week 8 — What Actually Breaks at 10× Scale?

## Central Question

> **How do you predict the first bottleneck before scaling a system?**

### Concepts

- Horizontal scaling
- Vertical scaling
- Capacity
- Little's Law
- Throughput
- Concurrency
- Queueing
- Saturation
- Bottleneck analysis
- Connection limits
- DB connection pools
- CPU limits
- Memory limits
- Rate limits
- Hotspots
- Load shedding
- Backpressure

### Experiment

Take a system and test:

```text
100 RPS
500 RPS
1,000 RPS
5,000 RPS
10,000 RPS
```

Measure:

- p50
- p95
- p99
- CPU
- memory
- DB latency
- connections
- queue depth
- errors

### Key Question

What breaks first?

### Then:

Fix only that bottleneck.

Run the experiment again.

### Artifact

Capacity curve:

```text
Load → Throughput
Load → Latency
Load → CPU
Load → Errors
```

### Blog

**"We Scaled a Go Service From 1K to 10K RPS. Here's What Broke First."**

### Definition of Done

Given a system and workload, I can identify likely bottlenecks and design experiments to validate them.

---

# PHASE III — BUILD A REAL SYSTEM

# Week 9 — Build the Distributed Job Platform

## Central Question

> **Can we build a reliable asynchronous processing system from scratch?**

### System

Build a Go-based job-processing platform.

```text
                ┌─────────────┐
                │   Client    │
                └──────┬──────┘
                       │
                ┌──────▼──────┐
                │ Job API     │
                └──────┬──────┘
                       │
                ┌──────▼──────┐
                │ PostgreSQL  │
                └──────┬──────┘
                       │
                ┌──────▼──────┐
                │    Queue    │
                └──────┬──────┘
                       │
              ┌────────┼────────┐
              ▼        ▼        ▼
           Worker   Worker   Worker
              │        │        │
              └────────┼────────┘
                       │
                ┌──────▼──────┐
                │ PostgreSQL  │
                └─────────────┘
```

### Features

- Job submission
- Job IDs
- Queueing
- Workers
- Job status
- Retry
- Failure handling
- DLQ
- Idempotency
- Graceful shutdown

### Artifact

Working system + architecture diagram.

### Blog

**"I Built a Distributed Job Processing System in Go."**

### Definition of Done

A job can:

```text
Submit
→ Persist
→ Queue
→ Execute
→ Complete
```

with observable state throughout the lifecycle.

---

# Week 10 — Make the System Production-Like

## Central Question

> **What does "reliable" actually mean for a distributed worker system?**

### Add

- Structured logging
- Metrics
- Tracing
- Retry policy
- Exponential backoff
- Jitter
- Idempotency
- Worker heartbeats
- Timeouts
- Graceful shutdown
- Dead-letter handling
- Persistent job state
- Recovery

### Failure Scenarios

Test:

```text
Worker crash
DB unavailable
Queue unavailable
Duplicate message
Slow job
Poison job
External dependency timeout
Worker shutdown
Network interruption
```

### Artifact

Failure/recovery matrix.

### Blog

**"I Tried to Break My Distributed Worker System."**

### Definition of Done

The system has explicit behaviour for major failure modes.

---

# PHASE IV — BREAK THE SYSTEM

# Week 11 — Become the Incident

## Central Question

> **Can I diagnose a system when the obvious explanation is wrong?**

This week is intentionally different.

We will create incidents.

### Scenario 1

```text
Traffic +20%
CPU +80%
Latency +400%
```

Find the cause.

### Scenario 2

```text
DB healthy
Redis healthy
External API healthy
Go CPU 95%
```

Find the cause.

### Scenario 3

```text
Queue depth increasing
CPU only 50%
```

Find the cause.

### Scenario 4

```text
Workers increasing
Throughput decreasing
```

Find the cause.

### Scenario 5

```text
Requests timing out
External API shows no increase
```

Find the cause.

### Scenario 6

```text
Memory stable
GC CPU increasing
```

Find the cause.

### Method

For every incident:

```text
Observation
↓
Possible explanations
↓
Evidence required
↓
Experiment
↓
Root cause
↓
Mitigation
↓
Permanent fix
```

### Artifact

Six mini postmortems.

### Blog

**"I Broke My Own Distributed System Six Different Ways."**

### Definition of Done

I stop treating symptoms as root causes.

---

# PHASE V — PERFORMANCE ENGINEERING

# Week 12 — How Fast Can We Make It?

## Central Question

> **How do we turn a working system into a measured, optimized system?**

### Concepts

- Benchmarking
- Profiling
- CPU efficiency
- Allocation rate
- Memory efficiency
- GC tuning
- Concurrency tuning
- Connection pooling
- Batching
- Serialization
- Algorithmic complexity
- Lock contention
- Throughput vs latency
- Tail latency

### Experiment

Establish a baseline:

```text
Throughput:
Latency:
CPU:
Memory:
Allocations:
GC:
DB:
Queue:
```

Then optimize one variable at a time.

Examples:

```text
Worker count
Connection pool
Batch size
Serialization
Caching
Allocation patterns
Concurrency
```

### Rule

**One change → one measurement.**

No "I changed five things and it got faster."

### Artifact

Before/after benchmark report.

### Blog

**"Making a Go Distributed System Faster: One Bottleneck at a Time."**

### Definition of Done

I can demonstrate at least several optimizations with measured evidence and explain their trade-offs.

---

# Week 13 — The Systems Engineer Case Study

## Central Question

> **Can I communicate an engineering investigation at a professional level?**

This is the synthesis week.

Choose the strongest investigation from the previous 12 weeks.

Turn it into a serious engineering case study.

### Structure

```text
1. Problem
2. System architecture
3. Initial symptoms
4. Baseline metrics
5. Initial hypotheses
6. Experiments
7. Evidence
8. Root cause
9. Fix
10. Before/after measurements
11. Trade-offs
12. What we got wrong
13. What we learned
14. What we would do in production
```

### Deliverables

- Architecture diagram
- Benchmark graphs
- Profiling evidence
- Code
- Incident timeline
- Final technical article
- README
- Lessons learned

### Blog

Potential titles:

**"How I Diagnosed and Fixed a 5× Latency Regression in a Distributed Go System"**

or

**"What 90 Days of Breaking Systems Taught Me About Performance Engineering"**

### Definition of Done

A senior engineer unfamiliar with the project should be able to read the case study and understand:

- what happened
- why it happened
- how we proved it
- what we changed
- what improved
- what trade-offs remain

---

# 7. The Final 90-Day Deliverable

The final GitHub repository should tell a coherent story.

```text
90-days-of-systems/
│
├── README.md
│
├── 00-SYLLABUS.md
│
├── weekly-retrospectives/
│
├── week-01-go-concurrency/
│   ├── README.md
│   ├── experiments/
│   ├── benchmarks/
│   └── blog.md
│
├── week-02-profiling/
├── week-03-networking/
├── week-04-observability/
├── week-05-queues/
├── week-06-distributed-failures/
├── week-07-consistency/
├── week-08-scaling/
├── week-09-job-platform/
├── week-10-reliability/
├── week-11-failure-lab/
├── week-12-performance/
│
└── final-case-study/
```

---

# 8. What We Should Have at the End

## Knowledge

Strong working understanding of:

- Go concurrency
- Go runtime behaviour
- CPU and memory
- profiling
- networking
- observability
- queues
- distributed failure
- retries
- idempotency
- consistency
- caching
- databases
- scaling
- backpressure
- performance engineering

## Engineering Artifacts

Target:

- 12+ experiments
- 10–13 technical write-ups
- multiple benchmark datasets
- multiple profiling sessions
- flame graphs
- architecture diagrams
- failure simulations
- one serious distributed Go project
- one polished engineering case study

## Professional Evidence

The repository should demonstrate:

> **I don't just know system-design vocabulary. I can investigate and engineer systems.**

---

# 9. Weekly Scorecard

At the end of every week:

```text
# Week X Retrospective

## Central Question

What were we trying to answer?

## Prediction

What did I think would happen?

## Observation

What actually happened?

## Explanation

Why did it happen?

## Evidence

What measurements support the conclusion?

## Experiment

What did I change?

## Result

What changed after the experiment?

## Artifact

What did I build?

## Failure

What did I deliberately break?

## Biggest Lesson

What do I understand now that I didn't understand before?

## Remaining Questions

What don't I understand yet?

## Blog

Published: YES / NO

## Definition of Done

Complete: YES / NO

## Confidence

__/10
```

---

# 10. Rules for Learning in Public

The blog is **not** the objective.

Engineering is the objective.

The blog exists to force clear thinking.

### Prefer

> "I expected X. I observed Y. I initially thought Z. The profile showed something else. Here's what actually happened."

### Avoid

> "10 Things You Should Know About Go Concurrency."

We want **engineering investigations**, not generic tutorials.

The best posts may include mistakes.

A mistake followed by a good investigation is valuable evidence of engineering growth.

---

# 11. Rules for Experiments

Every meaningful experiment should have:

### Hypothesis

> I believe X because Y.

### Controlled Variable

What are we changing?

### Measurements

What will we observe?

### Baseline

What happens before the change?

### Experiment

What happens after the change?

### Interpretation

Why did the result occur?

### Trade-off

What did we gain and what did we sacrifice?

---

# 12. What We Will NOT Optimize For

During these 90 days we will deliberately avoid optimizing for:

- certifications
- number of technologies learned
- number of GitHub repositories
- number of blog posts
- LeetCode volume
- Kubernetes expertise
- cloud-service memorization
- framework knowledge
- buzzwords
- tutorial completion

We optimize for:

> **Systems intuition + measurement + engineering judgment.**

---

# 13. Career Connection

The longer-term progression we are building toward is:

```text
Backend Engineer
       ↓
Strong Backend Engineer
       ↓
Systems / Performance Engineer
       ↓
Senior Systems Engineer
       ↓
Systems Architect / Specialist
       ↓
High-value Problem Solver
       ↓
Independent Consultant / Specialist
```

The economic thesis is:

```text
Sell time
   ↓
Sell engineering ability
   ↓
Sell specialized expertise
   ↓
Sell judgment
   ↓
Sell outcomes
```

The 90 days are not expected to magically transform compensation overnight.

They are intended to build the **technical foundation and evidence of capability** required for the next career jump.

---

# 14. The Most Important Rule

If reality contradicts the syllabus:

> **Reality wins.**

If Week 3's networking experiment exposes a deeper systems question, we investigate it.

If Week 6's idempotency experiment reveals a database consistency issue, we follow it.

If the job-processing project exposes a major concurrency problem, we may spend additional time understanding it.

The syllabus is our **map**, not a prison.

The goal is not to finish every checkbox.

The goal is to become a better systems engineer.

---

# 15. Final Definition of Done

At the end of the 90 days, I should be able to receive a problem like:

> "Our service used to handle 5,000 RPS at 150ms p95. During today's sale it handles 30,000 RPS at 5 seconds p95 with 8% errors. Different teams blame Redis, Postgres, Kubernetes and an external API."

And my instinct should NOT be:

> "Let's add Redis."

or:

> "Let's move to Kubernetes."

or:

> "Let's increase the ECS instances."

Instead, my instinct should be:

```text
What changed?
       ↓
What does the data show?
       ↓
Where is the latency?
       ↓
What resource is saturated?
       ↓
What are the plausible hypotheses?
       ↓
What evidence distinguishes them?
       ↓
Let's measure.
       ↓
Now let's change one thing.
       ↓
Did the evidence improve?
       ↓
What are the trade-offs?
```

That is the engineer this 90-day program is designed to build.

---

# 16. The Contract

For the next \~90 days:

**I will spend approximately 14 hours/week.**

**I will prioritize depth over breadth.**

**I will build instead of merely watching.**

**I will break things deliberately.**

**I will measure before claiming.**

**I will distinguish symptoms from causes.**

**I will document what I learn.**

**I will publish what is worth sharing.**

**I will allow experiments to challenge my assumptions.**

And above all:

> **I will learn to understand systems, not merely learn technologies.**

---

# 17. Starting Point

## Week 1

### Question

> **Do more goroutines actually make Go programs faster?**

### First principle

> **Never assume concurrency improves performance. Measure it.**

### First experiment

Build the simplest CPU-bound workload possible.

Run it:

```text
1 worker
2 workers
4 workers
8 workers
16 workers
32 workers
64 workers
```

Measure the result.

Then we investigate **why the curve looks the way it does**.

---

**Status:** Curriculum defined\
**Next action:** Start Week 1 — build the baseline experiment.
