// Package issue provides MCP tools for managing Forgejo issue dependencies.
package issue

import (
	"context"
	"fmt"

	"codeberg.org/goern/forgejo-mcp/v2/operation/params"
	"codeberg.org/goern/forgejo-mcp/v2/pkg/forgejo"
	"codeberg.org/goern/forgejo-mcp/v2/pkg/log"
	"codeberg.org/goern/forgejo-mcp/v2/pkg/to"

	forgejo_sdk "codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const (
	// ListIssueDependenciesToolName is the name of the tool for listing issue dependencies
	ListIssueDependenciesToolName = "list_issue_dependencies"
	// ListBlockedIssuesToolName is the name of the tool for listing blocked issues
	ListBlockedIssuesToolName = "list_blocked_issues"
	// CreateIssueDependencyToolName is the name of the tool for creating issue dependencies
	CreateIssueDependencyToolName = "create_issue_dependency"
	// RemoveIssueDependencyToolName is the name of the tool for removing issue dependencies
	RemoveIssueDependencyToolName = "remove_issue_dependency"
)

var (
	// ListIssueDependenciesTool lists all issues that block the specified issue
	ListIssueDependenciesTool = mcp.NewTool(
		ListIssueDependenciesToolName,
		mcp.WithDescription("List all issues that block the specified issue (dependencies)"),
		mcp.WithString("owner", mcp.Required(), mcp.Description(params.Owner)),
		mcp.WithString("repo", mcp.Required(), mcp.Description(params.Repo)),
		mcp.WithNumber("index", mcp.Required(), mcp.Description(params.IssueIndex)),
		mcp.WithNumber("page", mcp.Description(params.Page), mcp.DefaultNumber(1)),
		mcp.WithNumber("limit", mcp.Description(params.Limit), mcp.DefaultNumber(50)),
	)

	// ListBlockedIssuesTool lists all issues that are blocked by the specified issue
	ListBlockedIssuesTool = mcp.NewTool(
		ListBlockedIssuesToolName,
		mcp.WithDescription("List all issues that are blocked by the specified issue"),
		mcp.WithString("owner", mcp.Required(), mcp.Description(params.Owner)),
		mcp.WithString("repo", mcp.Required(), mcp.Description(params.Repo)),
		mcp.WithNumber("index", mcp.Required(), mcp.Description(params.IssueIndex)),
	)

	// CreateIssueDependencyTool creates a dependency relationship between issues
	CreateIssueDependencyTool = mcp.NewTool(
		CreateIssueDependencyToolName,
		mcp.WithDescription("Create a dependency relationship - the specified issue will be blocked by new_dependency"),
		mcp.WithString("owner", mcp.Required(), mcp.Description(params.Owner)),
		mcp.WithString("repo", mcp.Required(), mcp.Description(params.Repo)),
		mcp.WithNumber("index", mcp.Required(), mcp.Description("Issue index (the issue that will be blocked)")),
		mcp.WithNumber("new_dependency", mcp.Required(), mcp.Description(params.NewDependency)),
	)

	// RemoveIssueDependencyTool removes a dependency relationship between issues
	RemoveIssueDependencyTool = mcp.NewTool(
		RemoveIssueDependencyToolName,
		mcp.WithDescription("Remove a dependency relationship between issues"),
		mcp.WithString("owner", mcp.Required(), mcp.Description(params.Owner)),
		mcp.WithString("repo", mcp.Required(), mcp.Description(params.Repo)),
		mcp.WithNumber("index", mcp.Required(), mcp.Description("Issue index (the issue that is blocked)")),
		mcp.WithNumber("dependency", mcp.Required(), mcp.Description(params.Dependency)),
	)
)

// RegisterDependencyTools registers all dependency-related MCP tools with the server.
func RegisterDependencyTools(s *server.MCPServer) {
	s.AddTool(ListIssueDependenciesTool, ListIssueDependenciesFn)
	s.AddTool(ListBlockedIssuesTool, ListBlockedIssuesFn)
	s.AddTool(CreateIssueDependencyTool, CreateIssueDependencyFn)
	s.AddTool(RemoveIssueDependencyTool, RemoveIssueDependencyFn)
}

// ListIssueDependenciesFn lists all issues that block the specified issue.
func ListIssueDependenciesFn(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called ListIssueDependenciesFn")
	owner, err := req.RequireString("owner")
	if err != nil {
		return to.ErrorResult(err)
	}
	repo, err := req.RequireString("repo")
	if err != nil {
		return to.ErrorResult(err)
	}
	index, err := req.RequireFloat("index")
	if err != nil {
		return to.ErrorResult(err)
	}
	page := req.GetFloat("page", 1)
	limit := req.GetFloat("limit", 50)

	opt := forgejo_sdk.ListDependenciesOptions{
		ListOptions: forgejo_sdk.ListOptions{
			Page:     int(page),
			PageSize: int(limit),
		},
	}

	dependencies, _, err := forgejo.Client().ListIssueDependencies(owner, repo, int64(index), opt)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("list issue dependencies err: %v", err))
	}
	return to.TextResult(dependencies)
}

// ListBlockedIssuesFn lists all issues that are blocked by the specified issue.
func ListBlockedIssuesFn(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called ListBlockedIssuesFn")
	owner, err := req.RequireString("owner")
	if err != nil {
		return to.ErrorResult(err)
	}
	repo, err := req.RequireString("repo")
	if err != nil {
		return to.ErrorResult(err)
	}
	index, err := req.RequireFloat("index")
	if err != nil {
		return to.ErrorResult(err)
	}

	issues, _, err := forgejo.Client().ListBlockedIssues(owner, repo, int64(index))
	if err != nil {
		return to.ErrorResult(fmt.Errorf("list blocked issues err: %v", err))
	}
	return to.TextResult(issues)
}

// CreateIssueDependencyFn creates a dependency relationship between issues.
func CreateIssueDependencyFn(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called CreateIssueDependencyFn")
	owner, err := req.RequireString("owner")
	if err != nil {
		return to.ErrorResult(err)
	}
	repo, err := req.RequireString("repo")
	if err != nil {
		return to.ErrorResult(err)
	}
	index, err := req.RequireFloat("index")
	if err != nil {
		return to.ErrorResult(err)
	}
	newDependency, err := req.RequireFloat("new_dependency")
	if err != nil {
		return to.ErrorResult(err)
	}

	// Client-side validation: prevent self-dependency
	if int64(index) == int64(newDependency) {
		return to.ErrorResult(fmt.Errorf("self-dependency not allowed: issue cannot depend on itself"))
	}

	opt := forgejo_sdk.CreateIssueDependencyOption{
		NewDependency: int64(newDependency),
	}

	_, err = forgejo.Client().CreateIssueDependency(owner, repo, int64(index), opt)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("create issue dependency err: %v", err))
	}

	// Fetch updated issue to return
	issue, _, err := forgejo.Client().GetIssue(owner, repo, int64(index))
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get updated issue err: %v", err))
	}
	return to.TextResult(issue)
}

// RemoveIssueDependencyFn removes a dependency relationship between issues.
func RemoveIssueDependencyFn(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called RemoveIssueDependencyFn")
	owner, err := req.RequireString("owner")
	if err != nil {
		return to.ErrorResult(err)
	}
	repo, err := req.RequireString("repo")
	if err != nil {
		return to.ErrorResult(err)
	}
	index, err := req.RequireFloat("index")
	if err != nil {
		return to.ErrorResult(err)
	}
	dependency, err := req.RequireFloat("dependency")
	if err != nil {
		return to.ErrorResult(err)
	}

	_, err = forgejo.Client().RemoveIssueDependency(owner, repo, int64(index), int64(dependency))
	if err != nil {
		return to.ErrorResult(fmt.Errorf("remove issue dependency err: %v", err))
	}

	// Fetch updated issue to return
	issue, _, err := forgejo.Client().GetIssue(owner, repo, int64(index))
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get updated issue err: %v", err))
	}
	return to.TextResult(issue)
}
