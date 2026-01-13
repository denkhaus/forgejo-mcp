// Package action provides MCP tools for managing action secrets and workflows.
package action

import (
	"context"
	"encoding/json"
	"fmt"

	"codeberg.org/goern/forgejo-mcp/v2/operation/params"
	"codeberg.org/goern/forgejo-mcp/v2/pkg/forgejo"
	"codeberg.org/goern/forgejo-mcp/v2/pkg/log"
	"codeberg.org/goern/forgejo-mcp/v2/pkg/to"

	forgejo_sdk "codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Tool names for action operations.
const (
	UpdateRepoActionSecretToolName = "update_repo_action_secret"
	DeleteRepoActionSecretToolName = "delete_repo_action_secret"
	UpdateOrgActionSecretToolName  = "update_org_action_secret"
	DeleteOrgActionSecretToolName  = "delete_org_action_secret"
	UpdateUserActionSecretToolName = "update_user_action_secret"
	DeleteUserActionSecretToolName = "delete_user_action_secret"
	WorkflowDispatchToolName       = "workflow_dispatch"
)

// MCP tool definitions for action operations.
var (
	// Repo Level Secrets
	UpdateRepoActionSecretTool = mcp.NewTool(
		UpdateRepoActionSecretToolName,
		mcp.WithDescription("Update an action secret in a repository"),
		mcp.WithString("owner", mcp.Required(), mcp.Description(params.Owner)),
		mcp.WithString("repo", mcp.Required(), mcp.Description(params.Repo)),
		mcp.WithString("name", mcp.Required(), mcp.Description("Secret name")),
		mcp.WithString("data", mcp.Required(), mcp.Description("New secret data")),
	)

	DeleteRepoActionSecretTool = mcp.NewTool(
		DeleteRepoActionSecretToolName,
		mcp.WithDescription("Delete an action secret from a repository"),
		mcp.WithString("owner", mcp.Required(), mcp.Description(params.Owner)),
		mcp.WithString("repo", mcp.Required(), mcp.Description(params.Repo)),
		mcp.WithString("name", mcp.Required(), mcp.Description("Secret name")),
	)

	// Org Level Secrets
	UpdateOrgActionSecretTool = mcp.NewTool(
		UpdateOrgActionSecretToolName,
		mcp.WithDescription("Update an action secret in an organization"),
		mcp.WithString("org", mcp.Required(), mcp.Description("Organization name")),
		mcp.WithString("name", mcp.Required(), mcp.Description("Secret name")),
		mcp.WithString("data", mcp.Required(), mcp.Description("New secret data")),
	)

	DeleteOrgActionSecretTool = mcp.NewTool(
		DeleteOrgActionSecretToolName,
		mcp.WithDescription("Delete an action secret from an organization"),
		mcp.WithString("org", mcp.Required(), mcp.Description("Organization name")),
		mcp.WithString("name", mcp.Required(), mcp.Description("Secret name")),
	)

	// User Level Secrets
	UpdateUserActionSecretTool = mcp.NewTool(
		UpdateUserActionSecretToolName,
		mcp.WithDescription("Update an action secret for the current user"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Secret name")),
		mcp.WithString("data", mcp.Required(), mcp.Description("New secret data")),
	)

	DeleteUserActionSecretTool = mcp.NewTool(
		DeleteUserActionSecretToolName,
		mcp.WithDescription("Delete an action secret for the current user"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Secret name")),
	)

	// Workflow Operations
	WorkflowDispatchTool = mcp.NewTool(
		WorkflowDispatchToolName,
		mcp.WithDescription("Trigger a workflow run via workflow_dispatch event"),
		mcp.WithString("owner", mcp.Required(), mcp.Description(params.Owner)),
		mcp.WithString("repo", mcp.Required(), mcp.Description(params.Repo)),
		mcp.WithString("workflow_filename", mcp.Required(), mcp.Description("Workflow filename (e.g., '.github/workflows/main.yml')")),
		mcp.WithString("ref", mcp.Description("Branch or tag to run the workflow on (default: default branch)")),
		mcp.WithString("inputs", mcp.Description("Input parameters for the workflow as JSON string (optional)")),
	)
)

// RegisterTool registers all action-related tools with the MCP server.
func RegisterTool(s *server.MCPServer) {
	// Repo Level
	s.AddTool(UpdateRepoActionSecretTool, UpdateRepoActionSecretFn)
	s.AddTool(DeleteRepoActionSecretTool, DeleteRepoActionSecretFn)

	// Org Level
	s.AddTool(UpdateOrgActionSecretTool, UpdateOrgActionSecretFn)
	s.AddTool(DeleteOrgActionSecretTool, DeleteOrgActionSecretFn)

	// User Level
	s.AddTool(UpdateUserActionSecretTool, UpdateUserActionSecretFn)
	s.AddTool(DeleteUserActionSecretTool, DeleteUserActionSecretFn)

	// Workflow
	s.AddTool(WorkflowDispatchTool, WorkflowDispatchFn)
}

// Repo Level Functions

// UpdateRepoActionSecretFn updates an action secret in a repository.
func UpdateRepoActionSecretFn(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called UpdateRepoActionSecretFn")
	owner, err := req.RequireString("owner")
	if err != nil {
		return to.ErrorResult(err)
	}
	repo, err := req.RequireString("repo")
	if err != nil {
		return to.ErrorResult(err)
	}
	name, err := req.RequireString("name")
	if err != nil {
		return to.ErrorResult(err)
	}
	data, err := req.RequireString("data")
	if err != nil {
		return to.ErrorResult(err)
	}

	_, err = forgejo.Client().UpdateRepoActionSecret(owner, repo, name,
		forgejo_sdk.UpdateRepoActionSecretOption{
			Value: data,
		})
	if err != nil {
		return to.ErrorResult(fmt.Errorf("update repo action secret err: %v", err))
	}
	return to.TextResult(fmt.Sprintf("Secret '%s' updated successfully in %s/%s", name, owner, repo))
}

// DeleteRepoActionSecretFn deletes an action secret from a repository.
func DeleteRepoActionSecretFn(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called DeleteRepoActionSecretFn")
	owner, err := req.RequireString("owner")
	if err != nil {
		return to.ErrorResult(err)
	}
	repo, err := req.RequireString("repo")
	if err != nil {
		return to.ErrorResult(err)
	}
	name, err := req.RequireString("name")
	if err != nil {
		return to.ErrorResult(err)
	}

	_, err = forgejo.Client().DeleteRepoActionSecret(owner, repo, name)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("delete repo action secret err: %v", err))
	}
	return to.TextResult(fmt.Sprintf("Secret '%s' deleted successfully from %s/%s", name, owner, repo))
}

// Org Level Functions

// UpdateOrgActionSecretFn updates an action secret in an organization.
func UpdateOrgActionSecretFn(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called UpdateOrgActionSecretFn")
	org, err := req.RequireString("org")
	if err != nil {
		return to.ErrorResult(err)
	}
	name, err := req.RequireString("name")
	if err != nil {
		return to.ErrorResult(err)
	}
	data, err := req.RequireString("data")
	if err != nil {
		return to.ErrorResult(err)
	}

	_, err = forgejo.Client().UpdateOrgActionSecret(org, name,
		forgejo_sdk.UpdateOrgActionSecretOption{
			Value: data,
		})
	if err != nil {
		return to.ErrorResult(fmt.Errorf("update org action secret err: %v", err))
	}
	return to.TextResult(fmt.Sprintf("Secret '%s' updated successfully in org %s", name, org))
}

// DeleteOrgActionSecretFn deletes an action secret from an organization.
func DeleteOrgActionSecretFn(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called DeleteOrgActionSecretFn")
	org, err := req.RequireString("org")
	if err != nil {
		return to.ErrorResult(err)
	}
	name, err := req.RequireString("name")
	if err != nil {
		return to.ErrorResult(err)
	}

	_, err = forgejo.Client().DeleteOrgActionSecret(org, name)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("delete org action secret err: %v", err))
	}
	return to.TextResult(fmt.Sprintf("Secret '%s' deleted successfully from org %s", name, org))
}

// User Level Functions

// UpdateUserActionSecretFn updates an action secret for the current user.
func UpdateUserActionSecretFn(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called UpdateUserActionSecretFn")
	name, err := req.RequireString("name")
	if err != nil {
		return to.ErrorResult(err)
	}
	data, err := req.RequireString("data")
	if err != nil {
		return to.ErrorResult(err)
	}

	_, err = forgejo.Client().UpdateUserActionSecret(name,
		forgejo_sdk.UpdateUserActionSecretOption{
			Value: data,
		})
	if err != nil {
		return to.ErrorResult(fmt.Errorf("update user action secret err: %v", err))
	}
	return to.TextResult(fmt.Sprintf("Secret '%s' updated successfully for current user", name))
}

// DeleteUserActionSecretFn deletes an action secret for the current user.
func DeleteUserActionSecretFn(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called DeleteUserActionSecretFn")
	name, err := req.RequireString("name")
	if err != nil {
		return to.ErrorResult(err)
	}

	_, err = forgejo.Client().DeleteUserActionSecret(name)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("delete user action secret err: %v", err))
	}
	return to.TextResult(fmt.Sprintf("Secret '%s' deleted successfully for current user", name))
}

// Workflow Functions

// WorkflowDispatchFn triggers a workflow run via workflow_dispatch event.
func WorkflowDispatchFn(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called WorkflowDispatchFn")
	owner, err := req.RequireString("owner")
	if err != nil {
		return to.ErrorResult(err)
	}
	repo, err := req.RequireString("repo")
	if err != nil {
		return to.ErrorResult(err)
	}
	workflowFilename, err := req.RequireString("workflow_filename")
	if err != nil {
		return to.ErrorResult(err)
	}
	ref := req.GetString("ref", "")
	inputsStr := req.GetString("inputs", "")

	opt := forgejo_sdk.WorkflowDispatchOption{
		Ref: ref,
	}

	// Parse inputs if provided
	if inputsStr != "" {
		var inputs map[string]any
		if err := json.Unmarshal([]byte(inputsStr), &inputs); err != nil {
			return to.ErrorResult(fmt.Errorf("invalid inputs JSON: %v", err))
		}
		opt.Inputs = inputs
	}

	run, _, err := forgejo.Client().WorkflowDispatch(owner, repo, workflowFilename, opt)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("workflow dispatch err: %v", err))
	}
	return to.TextResult(run)
}
