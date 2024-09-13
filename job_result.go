package microbatch

import (
	"github.com/google/uuid"
	"time"
)

// JobResult represents the result of a processed job.
type JobResult struct {
	JobID         uuid.UUID // Job identifier
	Error         error     // Error encountered, if any
	TimeProcessed time.Time // When the job was processed
}
