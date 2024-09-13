package microbatch

import (
	"github.com/google/uuid"
	"time"
)

// Job represents a unit of work to be processed.
// ID: A unique identifier for the job.
// TimeSubmitted: The time the job was submitted to the BatchOrchestrator.
// Data: The data associated with the job that will be processed.
type Job struct {
	ID            uuid.UUID   // Unique identifier for the job
	TimeSubmitted time.Time   // Timestamp of when the job was submitted
	Data          interface{} // Arbitrary data to be processed
}

// NewJob creates a new Job with a unique UUID and the provided data.
func NewJob(data interface{}) Job {
	return Job{
		ID:            uuid.New(), // Generate a new UUID
		TimeSubmitted: time.Now(), // Set current time as submission time
		Data:          data,
	}
}
