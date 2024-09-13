package main

import (
	"fmt"
	"github.com/chepiv/microbatch"
	"time"
)

type FrequencyBatchProcessor struct {
	batchCount int
}

func (bp *FrequencyBatchProcessor) Process(batch []microbatch.Job) []microbatch.JobResult {
	bp.batchCount++
	results := make([]microbatch.JobResult, len(batch))
	for i, job := range batch {
		time.Sleep(10 * time.Millisecond)
		results[i] = microbatch.JobResult{JobID: job.ID, Error: nil, TimeProcessed: time.Now()}
	}
	return results
}

func main() {
	// Create batch processor and orchestrator with a batch size of 10 and frequency of 100ms
	processor := &FrequencyBatchProcessor{}
	orchestrator := microbatch.NewBatchOrchestrator(processor, 10, 100*time.Millisecond)

	// Create 15 jobs (fewer than two full batches, to demonstrate frequency processing)
	jobs := make([]microbatch.Job, 15)
	for i := 0; i < 15; i++ {
		jobs[i] = microbatch.NewJob(fmt.Sprintf("job #%d", i+1))
	}

	// Submit jobs to the orchestrator with a delay to force batch frequency processing
	for i, job := range jobs {
		orchestrator.Submit(job)
		if i == 5 {
			// Introduce a delay to demonstrate batch frequency behavior
			time.Sleep(150 * time.Millisecond)
		}
	}

	// Shutdown the orchestrator
	orchestrator.Shutdown()

	// Show how many batches were processed
	fmt.Printf("Total batches processed: %d\n", processor.batchCount)
	fmt.Printf("Total processed jobs: %d\n", len(orchestrator.GetResults()))
}
