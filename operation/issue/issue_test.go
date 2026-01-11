package issue

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
)

// TestReplaceIssueLabelsTool verifies the tool definition is correctly configured
func TestReplaceIssueLabelsTool(t *testing.T) {
	tool := ReplaceIssueLabelsTool

	assert.Equal(t, "replace_issue_labels", tool.Name)
	assert.NotNil(t, tool.Description)

	// Check required parameters
	params := tool.InputSchema.Properties
	assert.Contains(t, params, "owner")
	assert.Contains(t, params, "repo")
	assert.Contains(t, params, "index")
	assert.Contains(t, params, "labels")

	// Verify required fields
	assert.Contains(t, tool.InputSchema.Required, "owner")
	assert.Contains(t, tool.InputSchema.Required, "repo")
	assert.Contains(t, tool.InputSchema.Required, "index")
	assert.Contains(t, tool.InputSchema.Required, "labels")
}

// TestDeleteIssueLabelTool verifies the tool definition is correctly configured
func TestDeleteIssueLabelTool(t *testing.T) {
	tool := DeleteIssueLabelTool

	assert.Equal(t, "delete_issue_label", tool.Name)
	assert.NotNil(t, tool.Description)

	// Check required parameters
	params := tool.InputSchema.Properties
	assert.Contains(t, params, "owner")
	assert.Contains(t, params, "repo")
	assert.Contains(t, params, "index")
	assert.Contains(t, params, "id")

	// Verify required fields
	assert.Contains(t, tool.InputSchema.Required, "owner")
	assert.Contains(t, tool.InputSchema.Required, "repo")
	assert.Contains(t, tool.InputSchema.Required, "index")
	assert.Contains(t, tool.InputSchema.Required, "id")
}

// TestAddIssueLabelsTool verifies the tool definition is correctly configured
func TestAddIssueLabelsTool(t *testing.T) {
	tool := AddIssueLabelsTools

	assert.Equal(t, "add_issue_labels", tool.Name)
	assert.NotNil(t, tool.Description)

	// Check required parameters
	params := tool.InputSchema.Properties
	assert.Contains(t, params, "owner")
	assert.Contains(t, params, "repo")
	assert.Contains(t, params, "index")
	assert.Contains(t, params, "labels")

	// Verify required fields
	assert.Contains(t, tool.InputSchema.Required, "owner")
	assert.Contains(t, tool.InputSchema.Required, "repo")
	assert.Contains(t, tool.InputSchema.Required, "index")
	assert.Contains(t, tool.InputSchema.Required, "labels")
}

// TestReplaceIssueLabelsFn_MissingRequiredParams tests error handling for missing required parameters
func TestReplaceIssueLabelsFn_MissingRequiredParams(t *testing.T) {
	tests := []struct {
		name     string
		args     map[string]interface{}
		wantErr  bool
		errField string
	}{
		{
			name: "missing owner",
			args: map[string]interface{}{
				"repo":   "test-repo",
				"index":  float64(1),
				"labels": "123,456",
			},
			wantErr:  true,
			errField: "owner",
		},
		{
			name: "missing repo",
			args: map[string]interface{}{
				"owner":  "test-owner",
				"index":  float64(1),
				"labels": "123,456",
			},
			wantErr:  true,
			errField: "repo",
		},
		{
			name: "missing index",
			args: map[string]interface{}{
				"owner":  "test-owner",
				"repo":   "test-repo",
				"labels": "123,456",
			},
			wantErr:  true,
			errField: "index",
		},
		{
			name: "missing labels",
			args: map[string]interface{}{
				"owner": "test-owner",
				"repo":  "test-repo",
				"index": float64(1),
			},
			wantErr:  true,
			errField: "labels",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: tt.args,
				},
			}

			result, err := ReplaceIssueLabelsFn(context.TODO(), req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), tt.errField)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

// TestReplaceIssueLabelsFn_InvalidLabelIDs tests error handling for invalid label formats
func TestReplaceIssueLabelsFn_InvalidLabelIDs(t *testing.T) {
	tests := []struct {
		name        string
		labels      string
		wantErr     bool
		errContains string
		skip        bool // Skip tests that require Forgejo API
	}{
		{
			name:        "empty label",
			labels:      "",
			wantErr:     true,
			errContains: "labels cannot be empty",
		},
		{
			name:        "whitespace only",
			labels:      "   ",
			wantErr:     true,
			errContains: "labels cannot be empty",
		},
		{
			name:        "empty label in list",
			labels:      "123,",
			wantErr:     true,
			errContains: "empty label not allowed",
			skip:        true, // Requires Forgejo API
		},
		{
			name:        "zero ID",
			labels:      "0",
			wantErr:     true,
			errContains: "ID must be positive",
			skip:        true, // Requires Forgejo API
		},
		{
			name:        "negative ID",
			labels:      "-1",
			wantErr:     true,
			errContains: "ID must be positive",
			skip:        true, // Requires Forgejo API
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skip("Skipping test that requires Forgejo API")
			}

			req := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]interface{}{
						"owner":  "test-owner",
						"repo":   "test-repo",
						"index":  float64(1),
						"labels": tt.labels,
					},
				},
			}

			result, err := ReplaceIssueLabelsFn(context.TODO(), req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

// TestDeleteIssueLabelFn_MissingRequiredParams tests error handling for missing required parameters
func TestDeleteIssueLabelFn_MissingRequiredParams(t *testing.T) {
	tests := []struct {
		name     string
		args     map[string]interface{}
		wantErr  bool
		errField string
	}{
		{
			name: "missing owner",
			args: map[string]interface{}{
				"repo":  "test-repo",
				"index": float64(1),
				"id":    float64(123),
			},
			wantErr:  true,
			errField: "owner",
		},
		{
			name: "missing repo",
			args: map[string]interface{}{
				"owner": "test-owner",
				"index": float64(1),
				"id":    float64(123),
			},
			wantErr:  true,
			errField: "repo",
		},
		{
			name: "missing index",
			args: map[string]interface{}{
				"owner": "test-owner",
				"repo":  "test-repo",
				"id":    float64(123),
			},
			wantErr:  true,
			errField: "index",
		},
		{
			name: "missing id",
			args: map[string]interface{}{
				"owner": "test-owner",
				"repo":  "test-repo",
				"index": float64(1),
			},
			wantErr:  true,
			errField: "id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: tt.args,
				},
			}

			result, err := DeleteIssueLabelFn(context.TODO(), req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), tt.errField)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

// TestAddIssueLabelsFn_InvalidLabelIDs tests error handling for invalid label formats in add operation
func TestAddIssueLabelsFn_InvalidLabelIDs(t *testing.T) {
	tests := []struct {
		name        string
		labels      string
		wantErr     bool
		errContains string
		skip        bool // Skip tests that require Forgejo API
	}{
		{
			name:        "empty label",
			labels:      "",
			wantErr:     true,
			errContains: "labels cannot be empty",
		},
		{
			name:        "whitespace only",
			labels:      "   ",
			wantErr:     true,
			errContains: "labels cannot be empty",
		},
		{
			name:        "empty label in list",
			labels:      "123,",
			wantErr:     true,
			errContains: "empty label not allowed",
			skip:        true, // Requires Forgejo API
		},
		{
			name:        "zero ID",
			labels:      "0",
			wantErr:     true,
			errContains: "ID must be positive",
			skip:        true, // Requires Forgejo API
		},
		{
			name:        "negative ID",
			labels:      "-1",
			wantErr:     true,
			errContains: "ID must be positive",
			skip:        true, // Requires Forgejo API
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skip("Skipping test that requires Forgejo API")
			}

			req := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]interface{}{
						"owner":  "test-owner",
						"repo":   "test-repo",
						"index":  float64(1),
						"labels": tt.labels,
					},
				},
			}

			result, err := AddIssueLabelsFn(context.TODO(), req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}
