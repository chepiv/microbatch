
# Microbatch Library

The Microbatch library provides functionality for submitting and processing jobs in batches. It allows users to define their own job processing logic by implementing the `BatchProcessor` interface, which the library uses to process jobs in configurable batches based on size and frequency.

## Features
- Submit jobs that are processed in batches.
- Configure the batch size (number of jobs per batch).
- Configure the batch frequency (time interval for processing batches).
- Graceful shutdown that ensures all submitted jobs are processed.

## Usage

1. **Implement a BatchProcessor:**
   
   Users need to define how to process a batch of jobs by implementing the `BatchProcessor` interface.

   ```go
   type MyBatchProcessor struct {}

   func (bp *MyBatchProcessor) Process(batch []Job) []JobResult {
       var results []JobResult
       for _, job := range batch {
           // Process the job data...
           results = append(results, JobResult{JobID: job.ID, Error: nil, TimeProcessed: time.Now()})
       }
       return results
   }
   ```

2. **Create a BatchOrchestrator:**

   The `BatchOrchestrator` handles job submission and batching.

   ```go
   processor := &MyBatchProcessor{}
   batchOrchestrator := NewBatchOrchestrator(processor, 10, 2*time.Second)
   ```

3. **Submit Jobs:**

   Submit jobs to the `BatchOrchestrator`.

   ```go
   job := Job{ID: "job1", TimeSubmitted: time.Now(), Data: "some data"}
   batchOrchestrator.Submit(job)
   ```

4. **Shutdown:**

   Ensure all jobs are processed before shutting down.

   ```go
   batchOrchestrator.Shutdown()
   ```

## Installation

```bash
go get github.com/chepiv/microbatch
```