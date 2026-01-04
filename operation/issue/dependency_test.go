package issue

import (
	"context"
	"testing"

	forgejo_sdk "codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
)

func TestListIssueDependenciesFn(t *testing.T) {
	tests := []struct {
		name          string
		requestParams map[string]interface{}
		wantErr       bool
		errContains   string
	}{
		{
			name: "happy path - dependencies exist",
			requestParams: map[string]interface{}{
				"owner": "testowner",
				"repo":  "testrepo",
				"index": float64(1),
				"page":  float64(1),
				"limit": float64(50),
			},
			wantErr: false,
		},
		{
			name: "empty dependencies",
			requestParams: map[string]interface{}{
				"owner": "testowner",
				"repo":  "testrepo",
				"index": float64(999),
				"page":  float64(1),
				"limit": float64(50),
			},
			wantErr: false,
		},
		{
			name: "missing required param - owner",
			requestParams: map[string]interface{}{
				"repo":  "testrepo",
				"index": float64(1),
			},
			wantErr:     true,
			errContains: "owner",
		},
		{
			name: "missing required param - repo",
			requestParams: map[string]interface{}{
				"owner": "testowner",
				"index": float64(1),
			},
			wantErr:     true,
			errContains: "repo",
		},
		{
			name: "missing required param - index",
			requestParams: map[string]interface{}{
				"owner": "testowner",
				"repo":  "testrepo",
			},
			wantErr:     true,
			errContains: "index",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: tt.requestParams,
				},
			}

			result, err := ListIssueDependenciesFn(context.Background(), req)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				// Note: This will fail with actual API call if not mocked
				// In a real test environment, you would mock the forgejo.Client()
				// For now, we just verify the function signature works
				assert.NotNil(t, result)
			}
		})
	}
}

func TestListBlockedIssuesFn(t *testing.T) {
	tests := []struct {
		name          string
		requestParams map[string]interface{}
		wantErr       bool
		errContains   string
	}{
		{
			name: "happy path - blocked issues exist",
			requestParams: map[string]interface{}{
				"owner": "testowner",
				"repo":  "testrepo",
				"index": float64(1),
			},
			wantErr: false,
		},
		{
			name: "empty blocked issues",
			requestParams: map[string]interface{}{
				"owner": "testowner",
				"repo":  "testrepo",
				"index": float64(999),
			},
			wantErr: false,
		},
		{
			name: "missing required param - owner",
			requestParams: map[string]interface{}{
				"repo":  "testrepo",
				"index": float64(1),
			},
			wantErr:     true,
			errContains: "owner",
		},
		{
			name: "missing required param - repo",
			requestParams: map[string]interface{}{
				"owner": "testowner",
				"index": float64(1),
			},
			wantErr:     true,
			errContains: "repo",
		},
		{
			name: "missing required param - index",
			requestParams: map[string]interface{}{
				"owner": "testowner",
				"repo":  "testrepo",
			},
			wantErr:     true,
			errContains: "index",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: tt.requestParams,
				},
			}

			result, err := ListBlockedIssuesFn(context.Background(), req)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NotNil(t, result)
			}
		})
	}
}

func TestCreateIssueDependencyFn(t *testing.T) {
	tests := []struct {
		name          string
		requestParams map[string]interface{}
		wantErr       bool
		errContains   string
	}{
		{
			name: "happy path - dependency created",
			requestParams: map[string]interface{}{
				"owner":          "testowner",
				"repo":           "testrepo",
				"index":          float64(5),
				"new_dependency": float64(1),
			},
			wantErr: false,
		},
		{
			name: "self-dependency prevention",
			requestParams: map[string]interface{}{
				"owner":          "testowner",
				"repo":           "testrepo",
				"index":          float64(5),
				"new_dependency": float64(5),
			},
			wantErr:     true,
			errContains: "self-dependency not allowed",
		},
		{
			name: "missing required param - owner",
			requestParams: map[string]interface{}{
				"repo":           "testrepo",
				"index":          float64(5),
				"new_dependency": float64(1),
			},
			wantErr:     true,
			errContains: "owner",
		},
		{
			name: "missing required param - repo",
			requestParams: map[string]interface{}{
				"owner":          "testowner",
				"index":          float64(5),
				"new_dependency": float64(1),
			},
			wantErr:     true,
			errContains: "repo",
		},
		{
			name: "missing required param - index",
			requestParams: map[string]interface{}{
				"owner":          "testowner",
				"repo":           "testrepo",
				"new_dependency": float64(1),
			},
			wantErr:     true,
			errContains: "index",
		},
		{
			name: "missing required param - new_dependency",
			requestParams: map[string]interface{}{
				"owner": "testowner",
				"repo":  "testrepo",
				"index": float64(5),
			},
			wantErr:     true,
			errContains: "new_dependency",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: tt.requestParams,
				},
			}

			result, err := CreateIssueDependencyFn(context.Background(), req)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NotNil(t, result)
			}
		})
	}
}

func TestRemoveIssueDependencyFn(t *testing.T) {
	tests := []struct {
		name          string
		requestParams map[string]interface{}
		wantErr       bool
		errContains   string
	}{
		{
			name: "happy path - dependency removed",
			requestParams: map[string]interface{}{
				"owner":      "testowner",
				"repo":       "testrepo",
				"index":      float64(5),
				"dependency": float64(1),
			},
			wantErr: false,
		},
		{
			name: "missing required param - owner",
			requestParams: map[string]interface{}{
				"repo":       "testrepo",
				"index":      float64(5),
				"dependency": float64(1),
			},
			wantErr:     true,
			errContains: "owner",
		},
		{
			name: "missing required param - repo",
			requestParams: map[string]interface{}{
				"owner":      "testowner",
				"index":      float64(5),
				"dependency": float64(1),
			},
			wantErr:     true,
			errContains: "repo",
		},
		{
			name: "missing required param - index",
			requestParams: map[string]interface{}{
				"owner":      "testowner",
				"repo":       "testrepo",
				"dependency": float64(1),
			},
			wantErr:     true,
			errContains: "index",
		},
		{
			name: "missing required param - dependency",
			requestParams: map[string]interface{}{
				"owner": "testowner",
				"repo":  "testrepo",
				"index": float64(5),
			},
			wantErr:     true,
			errContains: "dependency",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: tt.requestParams,
				},
			}

			result, err := RemoveIssueDependencyFn(context.Background(), req)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NotNil(t, result)
			}
		})
	}
}

// Test tool definitions
func TestDependencyToolsAreDefined(t *testing.T) {
	t.Run("ListIssueDependenciesTool is defined", func(t *testing.T) {
		assert.Equal(t, "list_issue_dependencies", ListIssueDependenciesTool.Name)
		assert.NotNil(t, ListIssueDependenciesTool)
	})

	t.Run("ListBlockedIssuesTool is defined", func(t *testing.T) {
		assert.Equal(t, "list_blocked_issues", ListBlockedIssuesTool.Name)
		assert.NotNil(t, ListBlockedIssuesTool)
	})

	t.Run("CreateIssueDependencyTool is defined", func(t *testing.T) {
		assert.Equal(t, "create_issue_dependency", CreateIssueDependencyTool.Name)
		assert.NotNil(t, CreateIssueDependencyTool)
	})

	t.Run("RemoveIssueDependencyTool is defined", func(t *testing.T) {
		assert.Equal(t, "remove_issue_dependency", RemoveIssueDependencyTool.Name)
		assert.NotNil(t, RemoveIssueDependencyTool)
	})
}

// Test SDK option validation
func TestCreateIssueDependencyOptionValidation(t *testing.T) {
	tests := []struct {
		name      string
		option    forgejo_sdk.CreateIssueDependencyOption
		wantErr   bool
		errString string
	}{
		{
			name: "valid option",
			option: forgejo_sdk.CreateIssueDependencyOption{
				NewDependency: 5,
			},
			wantErr: false,
		},
		{
			name: "invalid - zero dependency",
			option: forgejo_sdk.CreateIssueDependencyOption{
				NewDependency: 0,
			},
			wantErr:   true,
			errString: "newDependency must be a positive issue number",
		},
		{
			name: "invalid - negative dependency",
			option: forgejo_sdk.CreateIssueDependencyOption{
				NewDependency: -1,
			},
			wantErr:   true,
			errString: "newDependency must be a positive issue number",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.option.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errString != "" {
					assert.Contains(t, err.Error(), tt.errString)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Benchmark tests for performance
func BenchmarkListIssueDependenciesFn(b *testing.B) {
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"owner": "testowner",
				"repo":  "testrepo",
				"index": float64(1),
				"page":  float64(1),
				"limit": float64(50),
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ListIssueDependenciesFn(context.Background(), req)
	}
}

func BenchmarkCreateIssueDependencyFn(b *testing.B) {
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"owner":          "testowner",
				"repo":           "testrepo",
				"index":          float64(5),
				"new_dependency": float64(1),
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = CreateIssueDependencyFn(context.Background(), req)
	}
}

// Example usage tests (for documentation)
func ExampleListIssueDependenciesFn() {
	// This example demonstrates how to list issue dependencies
	ctx := context.Background()
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"owner": "myorg",
				"repo":  "myrepo",
				"index": float64(42),
				"page":  float64(1),
				"limit": float64(50),
			},
		},
	}

	result, err := ListIssueDependenciesFn(ctx, req)
	if err != nil {
		// Handle error
		return
	}
	_ = result // Use result
}

func ExampleCreateIssueDependencyFn() {
	// This example demonstrates how to create an issue dependency
	// Issue #5 will be blocked by issue #1
	ctx := context.Background()
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"owner":          "myorg",
				"repo":           "myrepo",
				"index":          float64(5), // Issue to be blocked
				"new_dependency": float64(1), // Blocking issue
			},
		},
	}

	result, err := CreateIssueDependencyFn(ctx, req)
	if err != nil {
		// Handle error (e.g., self-dependency, circular dependency)
		return
	}
	_ = result // Use result
}
