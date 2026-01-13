// Package runner provides MCP tools for runner-related operations.
// This file contains org/user level runner jobs search tools.
package runner

import (
	"context"
	"fmt"
	"strings"

	"codeberg.org/goern/forgejo-mcp/v2/pkg/forgejo"
	"codeberg.org/goern/forgejo-mcp/v2/pkg/log"
	"codeberg.org/goern/forgejo-mcp/v2/pkg/to"

	forgejo_sdk "codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
	"github.com/mark3labs/mcp-go/mcp"
)

// Tool names for org/user level runner job operations.
const (
	SearchOrgRunnerJobsToolName  = "search_org_runner_jobs"
	SearchUserRunnerJobsToolName = "search_user_runner_jobs"
)

// MCP tool definitions for org/user level runner jobs search.
var (
	// SearchOrgRunnerJobsTool searches for organization's runner jobs
	SearchOrgRunnerJobsTool = mcp.NewTool(
		SearchOrgRunnerJobsToolName,
		mcp.WithDescription("Search for organization's action runner jobs according to filter conditions"),
		mcp.WithString("org", mcp.Required(), mcp.Description("Organization name")),
		mcp.WithString("labels", mcp.Description("Filter by job labels (comma-separated list)")),
	)

	// SearchUserRunnerJobsTool searches for user's runner jobs
	SearchUserRunnerJobsTool = mcp.NewTool(
		SearchUserRunnerJobsToolName,
		mcp.WithDescription("Search for user's action runner jobs according to filter conditions"),
		mcp.WithString("labels", mcp.Description("Filter by job labels (comma-separated list)")),
	)
)

// SearchOrgRunnerJobsFn searches for organization's action runner jobs according to filter conditions.
func SearchOrgRunnerJobsFn(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called SearchOrgRunnerJobsFn")
	org, err := req.RequireString("org")
	if err != nil {
		return to.ErrorResult(err)
	}
	labelsStr := req.GetString("labels", "")

	opt := forgejo_sdk.SearchOrgRunnerJobsOption{}
	if labelsStr != "" {
		opt.Labels = strings.Split(labelsStr, ",")
		// Trim whitespace from each label
		for i := range opt.Labels {
			opt.Labels[i] = strings.TrimSpace(opt.Labels[i])
		}
	}

	jobs, _, err := forgejo.Client().SearchOrgRunnerJobs(org, opt)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("search org runner jobs err: %v", err))
	}
	return to.TextResult(jobs)
}

// SearchUserRunnerJobsFn searches for user's action runner jobs according to filter conditions.
func SearchUserRunnerJobsFn(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called SearchUserRunnerJobsFn")
	labelsStr := req.GetString("labels", "")

	opt := forgejo_sdk.SearchUserRunnerJobsOption{}
	if labelsStr != "" {
		opt.Labels = strings.Split(labelsStr, ",")
		// Trim whitespace from each label
		for i := range opt.Labels {
			opt.Labels[i] = strings.TrimSpace(opt.Labels[i])
		}
	}

	jobs, _, err := forgejo.Client().SearchUserRunnerJobs(opt)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("search user runner jobs err: %v", err))
	}
	return to.TextResult(jobs)
}
