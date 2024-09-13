package microbatch

import (
	"sync"
	"time"
)

// BatchProcessor defines the contract for processing a batch of jobs.
// Process processes batch of jobs and returns slice of results.
type BatchProcessor interface {
	Process(batch []Job) []JobResult
}

// BatchOrchestrator manages submission and batch processing of jobs.
type BatchOrchestrator struct {
	processor BatchProcessor
	jobQueue  chan Job
	batchSize int
	batchFreq time.Duration
	wg        sync.WaitGroup
	results   []JobResult // Store the results internally
	mu        sync.Mutex  // Mutex to protect the results slice
}

// NewBatchOrchestrator creates a new BatchOrchestrator instance and starts listening for incoming jobs to process.
// BatchOrchestrator can be configured to send batches for processing either based on batch size or batch frequency.
func NewBatchOrchestrator(processor BatchProcessor, batchSize int, batchFreq time.Duration) *BatchOrchestrator {
	orchestrator := &BatchOrchestrator{
		processor: processor,
		jobQueue:  make(chan Job),
		batchSize: batchSize,
		batchFreq: batchFreq,
		results:   []JobResult{},
	}

	go orchestrator.run() // Start the batching process.
	return orchestrator
}

// Submit adds a job to the orchestrator for processing.
func (bo *BatchOrchestrator) Submit(job Job) {
	bo.wg.Add(1)
	bo.jobQueue <- job
}

// Shutdown gracefully shuts down the orchestrator, ensuring all jobs are processed.
func (bo *BatchOrchestrator) Shutdown() {
	close(bo.jobQueue)
	bo.wg.Wait()
}

// GetResults returns the results of all processed jobs.
func (bo *BatchOrchestrator) GetResults() []JobResult {
	bo.mu.Lock()
	defer bo.mu.Unlock()
	return bo.results
}

// run processes jobs in batches based on size or time frequency.
func (bo *BatchOrchestrator) run() {
	ticker := time.NewTicker(bo.batchFreq)
	defer ticker.Stop()

	var batch []Job

	for {
		select {
		case job, open := <-bo.jobQueue:
			// process left over jobs after channel closure
			if !open {
				bo.processBatch(batch)
				return
			}
			batch = append(batch, job)
			if len(batch) >= bo.batchSize {
				bo.processBatch(batch)
				batch = nil
			}
		case <-ticker.C:
			if len(batch) > 0 {
				bo.processBatch(batch)
				batch = nil
			}
		}
	}
}

// processBatch sends the batch of jobs to the BatchProcessor for processing.
func (bo *BatchOrchestrator) processBatch(batch []Job) {
	if len(batch) == 0 {
		return
	}
	results := bo.processor.Process(batch)

	bo.mu.Lock()
	bo.results = append(bo.results, results...)
	bo.mu.Unlock()

	for range batch {
		bo.wg.Done() // Mark job as processed
	}
}
