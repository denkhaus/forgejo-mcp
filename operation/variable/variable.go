// Package variable provides MCP tools for managing action variables at repo, org, and user levels.
package variable

import (
	"context"
	"fmt"

	"codeberg.org/goern/forgejo-mcp/v2/operation/params"
	"codeberg.org/goern/forgejo-mcp/v2/pkg/forgejo"
	"codeberg.org/goern/forgejo-mcp/v2/pkg/log"
	"codeberg.org/goern/forgejo-mcp/v2/pkg/to"

	forgejo_sdk "codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
	forgejo_models "codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Tool names for action variable operations.
const (
	ListRepoActionVariablesToolName  = "list_repo_action_variables"
	CreateRepoActionVariableToolName = "create_repo_action_variable"
	GetRepoActionVariableToolName    = "get_repo_action_variable"
	UpdateRepoActionVariableToolName = "update_repo_action_variable"
	DeleteRepoActionVariableToolName = "delete_repo_action_variable"
	ListOrgActionVariablesToolName   = "list_org_action_variables"
	CreateOrgActionVariableToolName  = "create_org_action_variable"
	GetOrgActionVariableToolName     = "get_org_action_variable"
	UpdateOrgActionVariableToolName  = "update_org_action_variable"
	DeleteOrgActionVariableToolName  = "delete_org_action_variable"
	ListUserActionVariablesToolName  = "list_user_action_variables"
	CreateUserActionVariableToolName = "create_user_action_variable"
	GetUserActionVariableToolName    = "get_user_action_variable"
	UpdateUserActionVariableToolName = "update_user_action_variable"
	DeleteUserActionVariableToolName = "delete_user_action_variable"
)

// MCP tool definitions for action variable operations.
var (
	// Repo Level Variables
	ListRepoActionVariablesTool = mcp.NewTool(
		ListRepoActionVariablesToolName,
		mcp.WithDescription("List action variables for a repository"),
		mcp.WithString("owner", mcp.Required(), mcp.Description(params.Owner)),
		mcp.WithString("repo", mcp.Required(), mcp.Description(params.Repo)),
		mcp.WithNumber("page", mcp.Description(params.Page), mcp.DefaultNumber(1)),
		mcp.WithNumber("limit", mcp.Description(params.Limit), mcp.DefaultNumber(20)),
	)

	CreateRepoActionVariableTool = mcp.NewTool(
		CreateRepoActionVariableToolName,
		mcp.WithDescription("Create an action variable in a repository"),
		mcp.WithString("owner", mcp.Required(), mcp.Description(params.Owner)),
		mcp.WithString("repo", mcp.Required(), mcp.Description(params.Repo)),
		mcp.WithString("name", mcp.Required(), mcp.Description("Variable name")),
		mcp.WithString("value", mcp.Required(), mcp.Description("Variable value")),
	)

	GetRepoActionVariableTool = mcp.NewTool(
		GetRepoActionVariableToolName,
		mcp.WithDescription("Get an action variable from a repository"),
		mcp.WithString("owner", mcp.Required(), mcp.Description(params.Owner)),
		mcp.WithString("repo", mcp.Required(), mcp.Description(params.Repo)),
		mcp.WithString("name", mcp.Required(), mcp.Description("Variable name")),
	)

	UpdateRepoActionVariableTool = mcp.NewTool(
		UpdateRepoActionVariableToolName,
		mcp.WithDescription("Update an action variable in a repository"),
		mcp.WithString("owner", mcp.Required(), mcp.Description(params.Owner)),
		mcp.WithString("repo", mcp.Required(), mcp.Description(params.Repo)),
		mcp.WithString("name", mcp.Required(), mcp.Description("Variable name")),
		mcp.WithString("value", mcp.Required(), mcp.Description("New variable value")),
	)

	DeleteRepoActionVariableTool = mcp.NewTool(
		DeleteRepoActionVariableToolName,
		mcp.WithDescription("Delete an action variable from a repository"),
		mcp.WithString("owner", mcp.Required(), mcp.Description(params.Owner)),
		mcp.WithString("repo", mcp.Required(), mcp.Description(params.Repo)),
		mcp.WithString("name", mcp.Required(), mcp.Description("Variable name")),
	)

	// Org Level Variables
	ListOrgActionVariablesTool = mcp.NewTool(
		ListOrgActionVariablesToolName,
		mcp.WithDescription("List action variables for an organization"),
		mcp.WithString("org", mcp.Required(), mcp.Description("Organization name")),
		mcp.WithNumber("page", mcp.Description(params.Page), mcp.DefaultNumber(1)),
		mcp.WithNumber("limit", mcp.Description(params.Limit), mcp.DefaultNumber(20)),
	)

	CreateOrgActionVariableTool = mcp.NewTool(
		CreateOrgActionVariableToolName,
		mcp.WithDescription("Create an action variable in an organization"),
		mcp.WithString("org", mcp.Required(), mcp.Description("Organization name")),
		mcp.WithString("name", mcp.Required(), mcp.Description("Variable name")),
		mcp.WithString("value", mcp.Required(), mcp.Description("Variable value")),
	)

	GetOrgActionVariableTool = mcp.NewTool(
		GetOrgActionVariableToolName,
		mcp.WithDescription("Get an action variable from an organization"),
		mcp.WithString("org", mcp.Required(), mcp.Description("Organization name")),
		mcp.WithString("name", mcp.Required(), mcp.Description("Variable name")),
	)

	UpdateOrgActionVariableTool = mcp.NewTool(
		UpdateOrgActionVariableToolName,
		mcp.WithDescription("Update an action variable in an organization"),
		mcp.WithString("org", mcp.Required(), mcp.Description("Organization name")),
		mcp.WithString("name", mcp.Required(), mcp.Description("Variable name")),
		mcp.WithString("value", mcp.Required(), mcp.Description("New variable value")),
	)

	DeleteOrgActionVariableTool = mcp.NewTool(
		DeleteOrgActionVariableToolName,
		mcp.WithDescription("Delete an action variable from an organization"),
		mcp.WithString("org", mcp.Required(), mcp.Description("Organization name")),
		mcp.WithString("name", mcp.Required(), mcp.Description("Variable name")),
	)

	// User Level Variables
	ListUserActionVariablesTool = mcp.NewTool(
		ListUserActionVariablesToolName,
		mcp.WithDescription("List action variables for the current user"),
		mcp.WithNumber("page", mcp.Description(params.Page), mcp.DefaultNumber(1)),
		mcp.WithNumber("limit", mcp.Description(params.Limit), mcp.DefaultNumber(20)),
	)

	CreateUserActionVariableTool = mcp.NewTool(
		CreateUserActionVariableToolName,
		mcp.WithDescription("Create an action variable for the current user"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Variable name")),
		mcp.WithString("value", mcp.Required(), mcp.Description("Variable value")),
	)

	GetUserActionVariableTool = mcp.NewTool(
		GetUserActionVariableToolName,
		mcp.WithDescription("Get an action variable for the current user"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Variable name")),
	)

	UpdateUserActionVariableTool = mcp.NewTool(
		UpdateUserActionVariableToolName,
		mcp.WithDescription("Update an action variable for the current user"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Variable name")),
		mcp.WithString("value", mcp.Required(), mcp.Description("New variable value")),
	)

	DeleteUserActionVariableTool = mcp.NewTool(
		DeleteUserActionVariableToolName,
		mcp.WithDescription("Delete an action variable for the current user"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Variable name")),
	)
)

// RegisterTool registers all variable-related tools with the MCP server.
func RegisterTool(s *server.MCPServer) {
	// Repo Level
	s.AddTool(ListRepoActionVariablesTool, ListRepoActionVariablesFn)
	s.AddTool(CreateRepoActionVariableTool, CreateRepoActionVariableFn)
	s.AddTool(GetRepoActionVariableTool, GetRepoActionVariableFn)
	s.AddTool(UpdateRepoActionVariableTool, UpdateRepoActionVariableFn)
	s.AddTool(DeleteRepoActionVariableTool, DeleteRepoActionVariableFn)

	// Org Level
	s.AddTool(ListOrgActionVariablesTool, ListOrgActionVariablesFn)
	s.AddTool(CreateOrgActionVariableTool, CreateOrgActionVariableFn)
	s.AddTool(GetOrgActionVariableTool, GetOrgActionVariableFn)
	s.AddTool(UpdateOrgActionVariableTool, UpdateOrgActionVariableFn)
	s.AddTool(DeleteOrgActionVariableTool, DeleteOrgActionVariableFn)

	// User Level
	s.AddTool(ListUserActionVariablesTool, ListUserActionVariablesFn)
	s.AddTool(CreateUserActionVariableTool, CreateUserActionVariableFn)
	s.AddTool(GetUserActionVariableTool, GetUserActionVariableFn)
	s.AddTool(UpdateUserActionVariableTool, UpdateUserActionVariableFn)
	s.AddTool(DeleteUserActionVariableTool, DeleteUserActionVariableFn)
}

// Repo Level Functions

// ListRepoActionVariablesFn lists action variables for a repository.
func ListRepoActionVariablesFn(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called ListRepoActionVariablesFn")
	owner, err := req.RequireString("owner")
	if err != nil {
		return to.ErrorResult(err)
	}
	repo, err := req.RequireString("repo")
	if err != nil {
		return to.ErrorResult(err)
	}
	page := req.GetFloat("page", 1)
	limit := req.GetFloat("limit", 20)

	variables, _, err := forgejo.Client().ListRepoActionVariables(owner, repo,
		forgejo_sdk.ListRepoActionVariablesOption{
			ListOptions: forgejo_sdk.ListOptions{
				Page:     int(page),
				PageSize: int(limit),
			},
		})
	if err != nil {
		return to.ErrorResult(fmt.Errorf("list repo action variables err: %v", err))
	}
	return to.TextResult(variables)
}

// CreateRepoActionVariableFn creates an action variable in a repository.
func CreateRepoActionVariableFn(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called CreateRepoActionVariableFn")
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
	value, err := req.RequireString("value")
	if err != nil {
		return to.ErrorResult(err)
	}

	_, err = forgejo.Client().CreateRepoActionVariable(owner, repo, name,
		forgejo_models.CreateVariableOption{
			Value: &value,
		})
	if err != nil {
		return to.ErrorResult(fmt.Errorf("create repo action variable err: %v", err))
	}
	return to.TextResult(fmt.Sprintf("Variable '%s' created successfully in %s/%s", name, owner, repo))
}

// GetRepoActionVariableFn gets an action variable from a repository.
func GetRepoActionVariableFn(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called GetRepoActionVariableFn")
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

	variable, _, err := forgejo.Client().GetRepoActionVariable(owner, repo, name)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get repo action variable err: %v", err))
	}
	return to.TextResult(variable)
}

// UpdateRepoActionVariableFn updates an action variable in a repository.
func UpdateRepoActionVariableFn(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called UpdateRepoActionVariableFn")
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
	value, err := req.RequireString("value")
	if err != nil {
		return to.ErrorResult(err)
	}

	_, err = forgejo.Client().UpdateRepoActionVariable(owner, repo, name,
		forgejo_models.UpdateVariableOption{
			Value: &value,
		})
	if err != nil {
		return to.ErrorResult(fmt.Errorf("update repo action variable err: %v", err))
	}
	return to.TextResult(fmt.Sprintf("Variable '%s' updated successfully in %s/%s", name, owner, repo))
}

// DeleteRepoActionVariableFn deletes an action variable from a repository.
func DeleteRepoActionVariableFn(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called DeleteRepoActionVariableFn")
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

	_, err = forgejo.Client().DeleteRepoActionVariable(owner, repo, name)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("delete repo action variable err: %v", err))
	}
	return to.TextResult(fmt.Sprintf("Variable '%s' deleted successfully from %s/%s", name, owner, repo))
}

// Org Level Functions

// ListOrgActionVariablesFn lists action variables for an organization.
func ListOrgActionVariablesFn(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called ListOrgActionVariablesFn")
	org, err := req.RequireString("org")
	if err != nil {
		return to.ErrorResult(err)
	}
	page := req.GetFloat("page", 1)
	limit := req.GetFloat("limit", 20)

	variables, _, err := forgejo.Client().ListOrgActionVariables(org,
		forgejo_sdk.ListOrgActionVariablesOption{
			ListOptions: forgejo_sdk.ListOptions{
				Page:     int(page),
				PageSize: int(limit),
			},
		})
	if err != nil {
		return to.ErrorResult(fmt.Errorf("list org action variables err: %v", err))
	}
	return to.TextResult(variables)
}

// CreateOrgActionVariableFn creates an action variable in an organization.
func CreateOrgActionVariableFn(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called CreateOrgActionVariableFn")
	org, err := req.RequireString("org")
	if err != nil {
		return to.ErrorResult(err)
	}
	name, err := req.RequireString("name")
	if err != nil {
		return to.ErrorResult(err)
	}
	value, err := req.RequireString("value")
	if err != nil {
		return to.ErrorResult(err)
	}

	_, err = forgejo.Client().CreateOrgActionVariable(org, name,
		forgejo_models.CreateVariableOption{
			Value: &value,
		})
	if err != nil {
		return to.ErrorResult(fmt.Errorf("create org action variable err: %v", err))
	}
	return to.TextResult(fmt.Sprintf("Variable '%s' created successfully in org %s", name, org))
}

// GetOrgActionVariableFn gets an action variable from an organization.
func GetOrgActionVariableFn(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called GetOrgActionVariableFn")
	org, err := req.RequireString("org")
	if err != nil {
		return to.ErrorResult(err)
	}
	name, err := req.RequireString("name")
	if err != nil {
		return to.ErrorResult(err)
	}

	variable, _, err := forgejo.Client().GetOrgActionVariable(org, name)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get org action variable err: %v", err))
	}
	return to.TextResult(variable)
}

// UpdateOrgActionVariableFn updates an action variable in an organization.
func UpdateOrgActionVariableFn(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called UpdateOrgActionVariableFn")
	org, err := req.RequireString("org")
	if err != nil {
		return to.ErrorResult(err)
	}
	name, err := req.RequireString("name")
	if err != nil {
		return to.ErrorResult(err)
	}
	value, err := req.RequireString("value")
	if err != nil {
		return to.ErrorResult(err)
	}

	_, err = forgejo.Client().UpdateOrgActionVariable(org, name,
		forgejo_models.UpdateVariableOption{
			Value: &value,
		})
	if err != nil {
		return to.ErrorResult(fmt.Errorf("update org action variable err: %v", err))
	}
	return to.TextResult(fmt.Sprintf("Variable '%s' updated successfully in org %s", name, org))
}

// DeleteOrgActionVariableFn deletes an action variable from an organization.
func DeleteOrgActionVariableFn(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called DeleteOrgActionVariableFn")
	org, err := req.RequireString("org")
	if err != nil {
		return to.ErrorResult(err)
	}
	name, err := req.RequireString("name")
	if err != nil {
		return to.ErrorResult(err)
	}

	_, err = forgejo.Client().DeleteOrgActionVariable(org, name)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("delete org action variable err: %v", err))
	}
	return to.TextResult(fmt.Sprintf("Variable '%s' deleted successfully from org %s", name, org))
}

// User Level Functions

// ListUserActionVariablesFn lists action variables for the current user.
func ListUserActionVariablesFn(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called ListUserActionVariablesFn")
	page := req.GetFloat("page", 1)
	limit := req.GetFloat("limit", 20)

	variables, _, err := forgejo.Client().ListUserActionVariables(
		forgejo_sdk.ListUserActionVariablesOption{
			ListOptions: forgejo_sdk.ListOptions{
				Page:     int(page),
				PageSize: int(limit),
			},
		})
	if err != nil {
		return to.ErrorResult(fmt.Errorf("list user action variables err: %v", err))
	}
	return to.TextResult(variables)
}

// CreateUserActionVariableFn creates an action variable for the current user.
func CreateUserActionVariableFn(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called CreateUserActionVariableFn")
	name, err := req.RequireString("name")
	if err != nil {
		return to.ErrorResult(err)
	}
	value, err := req.RequireString("value")
	if err != nil {
		return to.ErrorResult(err)
	}

	_, err = forgejo.Client().CreateUserActionVariable(name,
		forgejo_models.CreateVariableOption{
			Value: &value,
		})
	if err != nil {
		return to.ErrorResult(fmt.Errorf("create user action variable err: %v", err))
	}
	return to.TextResult(fmt.Sprintf("Variable '%s' created successfully for current user", name))
}

// GetUserActionVariableFn gets an action variable for the current user.
func GetUserActionVariableFn(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called GetUserActionVariableFn")
	name, err := req.RequireString("name")
	if err != nil {
		return to.ErrorResult(err)
	}

	variable, _, err := forgejo.Client().GetUserActionVariable(name)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get user action variable err: %v", err))
	}
	return to.TextResult(variable)
}

// UpdateUserActionVariableFn updates an action variable for the current user.
func UpdateUserActionVariableFn(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called UpdateUserActionVariableFn")
	name, err := req.RequireString("name")
	if err != nil {
		return to.ErrorResult(err)
	}
	value, err := req.RequireString("value")
	if err != nil {
		return to.ErrorResult(err)
	}

	_, err = forgejo.Client().UpdateUserActionVariable(name,
		forgejo_models.UpdateVariableOption{
			Value: &value,
		})
	if err != nil {
		return to.ErrorResult(fmt.Errorf("update user action variable err: %v", err))
	}
	return to.TextResult(fmt.Sprintf("Variable '%s' updated successfully for current user", name))
}

// DeleteUserActionVariableFn deletes an action variable for the current user.
func DeleteUserActionVariableFn(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called DeleteUserActionVariableFn")
	name, err := req.RequireString("name")
	if err != nil {
		return to.ErrorResult(err)
	}

	_, err = forgejo.Client().DeleteUserActionVariable(name)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("delete user action variable err: %v", err))
	}
	return to.TextResult(fmt.Sprintf("Variable '%s' deleted successfully for current user", name))
}
