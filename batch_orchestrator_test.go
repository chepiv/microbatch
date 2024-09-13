package microbatch

import (
	"sync"
	"testing"
	"time"
)

type MockBatchProcessor struct {
	mu               sync.Mutex
	processedBatches [][]Job
}

func (bp *MockBatchProcessor) Process(batch []Job) []JobResult {
	bp.mu.Lock()
	defer bp.mu.Unlock()
	bp.processedBatches = append(bp.processedBatches, batch)
	var results []JobResult

	return results
}

func TestBatchOrchestrator(t *testing.T) {
	tests := []struct {
		name            string
		jobs            []Job
		batchSize       int
		batchFreq       time.Duration
		delay           bool
		delayTime       time.Duration
		expectedBatches int
	}{
		{
			name:            "Process jobs based on batch size",
			jobs:            createJobs(10),
			batchSize:       5,
			batchFreq:       100 * time.Millisecond, // Large frequency, only batch size matters
			expectedBatches: 2,
		},
		{
			name:            "Process jobs based on batch frequency",
			jobs:            createJobs(10),
			batchSize:       15,
			batchFreq:       10 * time.Millisecond, // Frequency is small, forcing processing based on time
			delay:           true,
			delayTime:       2 * time.Millisecond,
			expectedBatches: 3, // 2 batches of 4 on ticker and 1 batch on shutdown
		},
		{
			name:            "Process fewer jobs than batch size",
			jobs:            createJobs(3),
			batchSize:       10,
			batchFreq:       100 * time.Millisecond,
			expectedBatches: 1,
		},
		{
			name:            "No jobs submitted",
			jobs:            createJobs(0),
			batchSize:       5,
			batchFreq:       1 * time.Second,
			expectedBatches: 0,
		},
		{
			name:            "Process jobs when batch size is an exact multiple",
			jobs:            createJobs(20),
			batchSize:       5,
			batchFreq:       500 * time.Millisecond, // Large frequency, batch size triggers processing
			expectedBatches: 4,
		},
		{
			name:            "Batch frequency triggers before batch size",
			jobs:            createJobs(7),
			batchSize:       10,
			batchFreq:       200 * time.Millisecond, // Small frequency forces processing based on time
			expectedBatches: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor := &MockBatchProcessor{}
			orchestrator := NewBatchOrchestrator(processor, tt.batchSize, tt.batchFreq)

			for _, job := range tt.jobs {
				if tt.delay {
					time.Sleep(tt.delayTime)
				}
				orchestrator.Submit(job)
			}

			orchestrator.Shutdown()

			if len(processor.processedBatches) != tt.expectedBatches {
				t.Errorf("expected %d batches to be processed, but got %d", tt.expectedBatches, len(processor.processedBatches))
			}

			totalJobsProcessed := 0
			for _, batch := range processor.processedBatches {
				totalJobsProcessed += len(batch)
			}
			if totalJobsProcessed != len(tt.jobs) {
				t.Errorf("expected %d jobs to be processed, but got %d", len(tt.jobs), totalJobsProcessed)
			}
		})
	}
}

func createJobs(num int) []Job {
	jobs := make([]Job, num)
	for i := 0; i < num; i++ {
		jobs[i] = NewJob(i)
	}
	return jobs
}
