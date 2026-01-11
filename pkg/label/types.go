package label

import "fmt"

// ResolvedLabel represents a label with both name and ID
type ResolvedLabel struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// LabelResolutionError provides detailed error information
type LabelResolutionError struct {
	Input       string
	Reason      string
	Suggestions []string
}

// Error returns the error message
func (e *LabelResolutionError) Error() string {
	if len(e.Suggestions) > 0 {
		return fmt.Sprintf("label '%s' %s. Did you mean: %s?", e.Input, e.Reason, formatSuggestions(e.Suggestions))
	}
	return fmt.Sprintf("label '%s' %s", e.Input, e.Reason)
}

// NewResolutionError creates a new LabelResolutionError
func NewResolutionError(input, reason string, suggestions []string) *LabelResolutionError {
	return &LabelResolutionError{
		Input:       input,
		Reason:      reason,
		Suggestions: suggestions,
	}
}

// formatSuggestions formats suggestions for error messages
func formatSuggestions(suggestions []string) string {
	if len(suggestions) == 0 {
		return ""
	}

	result := ""
	for i, s := range suggestions {
		if i > 0 {
			result += ", "
		}
		result += s
	}
	return result
}

// ResolutionResult contains the results of a label resolution operation
type ResolutionResult struct {
	LabelIDs       []int64
	ResolvedLabels []ResolvedLabel
}
