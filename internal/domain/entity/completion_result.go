package entity

// CompletionResult represents the result of processing forest completions.
type CompletionResult struct {
	ProcessedCount int
	FailedCount    int
	Errors         []error
}

func NewCompletionResult() *CompletionResult {
	return &CompletionResult{
		Errors: make([]error, 0),
	}
}

func (r *CompletionResult) HasErrors() bool {
	return len(r.Errors) > 0
}

func (r *CompletionResult) AddError(err error) {
	r.Errors = append(r.Errors, err)
	r.FailedCount++
}

func (r *CompletionResult) IncrementProcessed() {
	r.ProcessedCount++
}
