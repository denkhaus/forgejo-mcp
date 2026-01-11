package runner

import (
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
)

// TestListActionRunsTool verifies the tool definition is correctly configured
func TestListActionRunsTool(t *testing.T) {
	tool := ListActionRunsTool

	assert.Equal(t, "list_action_runs", tool.Name)
	assert.NotNil(t, tool.Description)

	// Check required parameters
	params := tool.InputSchema.Properties
	assert.Contains(t, params, "owner")
	assert.Contains(t, params, "repo")

	// Verify required fields
	assert.Contains(t, tool.InputSchema.Required, "owner")
	assert.Contains(t, tool.InputSchema.Required, "repo")

	// Check optional parameters
	assert.Contains(t, params, "page")
	assert.Contains(t, params, "limit")
	assert.Contains(t, params, "status")
	assert.Contains(t, params, "events")
	assert.Contains(t, params, "run_number")
	assert.Contains(t, params, "head_sha")
}

// TestGetActionRunTool verifies the tool definition is correctly configured
func TestGetActionRunTool(t *testing.T) {
	tool := GetActionRunTool

	assert.Equal(t, "get_action_run", tool.Name)
	assert.NotNil(t, tool.Description)

	// Check required parameters
	params := tool.InputSchema.Properties
	assert.Contains(t, params, "owner")
	assert.Contains(t, params, "repo")
	assert.Contains(t, params, "run_id")

	// Verify required fields
	assert.Contains(t, tool.InputSchema.Required, "owner")
	assert.Contains(t, tool.InputSchema.Required, "repo")
	assert.Contains(t, tool.InputSchema.Required, "run_id")
}

// TestSearchRunnerJobsTool verifies the tool definition is correctly configured
func TestSearchRunnerJobsTool(t *testing.T) {
	tool := SearchRunnerJobsTool

	assert.Equal(t, "search_runner_jobs", tool.Name)
	assert.NotNil(t, tool.Description)

	// Check required parameters
	params := tool.InputSchema.Properties
	assert.Contains(t, params, "owner")
	assert.Contains(t, params, "repo")

	// Verify required fields
	assert.Contains(t, tool.InputSchema.Required, "owner")
	assert.Contains(t, tool.InputSchema.Required, "repo")

	// Check optional parameters
	assert.Contains(t, params, "labels")
}

// TestGetRepoRunnerRegistrationTokenTool verifies the tool definition is correctly configured
func TestGetRepoRunnerRegistrationTokenTool(t *testing.T) {
	tool := GetRepoRunnerRegistrationTokenTool

	assert.Equal(t, "get_repo_runner_registration_token", tool.Name)
	assert.NotNil(t, tool.Description)

	// Check required parameters
	params := tool.InputSchema.Properties
	assert.Contains(t, params, "owner")
	assert.Contains(t, params, "repo")

	// Verify required fields
	assert.Contains(t, tool.InputSchema.Required, "owner")
	assert.Contains(t, tool.InputSchema.Required, "repo")
}

// TestListActionRunsFn_MissingOwner tests error handling when owner is missing
func TestListActionRunsFn_MissingOwner(t *testing.T) {
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"repo": "test-repo",
			},
		},
	}

	result, err := ListActionRunsFn(nil, req)
	assert.Error(t, err)
	assert.Nil(t, result)
}

// TestListActionRunsFn_MissingRepo tests error handling when repo is missing
func TestListActionRunsFn_MissingRepo(t *testing.T) {
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"owner": "test-owner",
			},
		},
	}

	result, err := ListActionRunsFn(nil, req)
	assert.Error(t, err)
	assert.Nil(t, result)
}

// TestGetActionRunFn_MissingRequiredParams tests error handling for missing required parameters
func TestGetActionRunFn_MissingRequiredParams(t *testing.T) {
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
				"run_id": float64(123),
			},
			wantErr:  true,
			errField: "owner",
		},
		{
			name: "missing repo",
			args: map[string]interface{}{
				"owner":  "test-owner",
				"run_id": float64(123),
			},
			wantErr:  true,
			errField: "repo",
		},
		{
			name: "missing run_id",
			args: map[string]interface{}{
				"owner": "test-owner",
				"repo":  "test-repo",
			},
			wantErr:  true,
			errField: "run_id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: tt.args,
				},
			}

			result, err := GetActionRunFn(nil, req)
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

// TestSearchRunnerJobsFn_MissingRequiredParams tests error handling for missing required parameters
func TestSearchRunnerJobsFn_MissingRequiredParams(t *testing.T) {
	tests := []struct {
		name     string
		args     map[string]interface{}
		wantErr  bool
		errField string
	}{
		{
			name: "missing owner",
			args: map[string]interface{}{
				"repo": "test-repo",
			},
			wantErr:  true,
			errField: "owner",
		},
		{
			name: "missing repo",
			args: map[string]interface{}{
				"owner": "test-owner",
			},
			wantErr:  true,
			errField: "repo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: tt.args,
				},
			}

			result, err := SearchRunnerJobsFn(nil, req)
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

// TestGetRepoRunnerRegistrationTokenFn_MissingRequiredParams tests error handling for missing required parameters
func TestGetRepoRunnerRegistrationTokenFn_MissingRequiredParams(t *testing.T) {
	tests := []struct {
		name     string
		args     map[string]interface{}
		wantErr  bool
		errField string
	}{
		{
			name: "missing owner",
			args: map[string]interface{}{
				"repo": "test-repo",
			},
			wantErr:  true,
			errField: "owner",
		},
		{
			name: "missing repo",
			args: map[string]interface{}{
				"owner": "test-owner",
			},
			wantErr:  true,
			errField: "repo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: tt.args,
				},
			}

			result, err := GetRepoRunnerRegistrationTokenFn(nil, req)
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

// TestParseStringList tests the parseStringList helper function
func TestParseStringList(t *testing.T) {
	tests := []struct {
		name string
		input string
		want []string
	}{
		{
			name: "single value",
			input: "success",
			want: []string{"success"},
		},
		{
			name: "multiple values",
			input: "success,failure,running",
			want: []string{"success", "failure", "running"},
		},
		{
			name: "values with spaces",
			input: "success, failure, running",
			want: []string{"success", "failure", "running"},
		},
		{
			name: "values with extra spaces",
			input: "  success  ,  failure  ,  running  ",
			want: []string{"success", "failure", "running"},
		},
		{
			name: "empty string",
			input: "",
			want: nil,
		},
		{
			name: "only commas",
			input: ",,",
			want: []string{},
		},
		{
			name: "trailing comma",
			input: "success,failure,",
			want: []string{"success", "failure"},
		},
		{
			name: "leading comma",
			input: ",success,failure",
			want: []string{"success", "failure"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseStringList(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
