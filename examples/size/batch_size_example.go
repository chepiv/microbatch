package main

import (
	"fmt"
	"github.com/chepiv/microbatch"
	"time"
)

type SizeBatchProcessor struct {
	batchCount int
}

func (bp *SizeBatchProcessor) Process(batch []microbatch.Job) []microbatch.JobResult {
	bp.batchCount++
	results := make([]microbatch.JobResult, len(batch))
	for i, job := range batch {
		time.Sleep(10 * time.Millisecond) // Simulate work
		results[i] = microbatch.JobResult{JobID: job.ID, Error: nil, TimeProcessed: time.Now()}
	}
	return results
}

func main() {
	processor := &SizeBatchProcessor{}
	orchestrator := microbatch.NewBatchOrchestrator(processor, 10, 5*time.Second) // Large frequency to focus on batch size

	jobs := make([]microbatch.Job, 100)
	for i := 0; i < 100; i++ {
		jobs[i] = microbatch.NewJob(fmt.Sprintf("job #%d", i+1))
	}

	for _, job := range jobs {
		orchestrator.Submit(job)
	}

	orchestrator.Shutdown()

	fmt.Printf("Total batches processed: %d\n", processor.batchCount)
	fmt.Printf("Total processed jobs: %d\n", len(orchestrator.GetResults()))
}
