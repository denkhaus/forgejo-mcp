package runner

import (
	"context"
	"fmt"

	"codeberg.org/goern/forgejo-mcp/v2/operation/params"
	"codeberg.org/goern/forgejo-mcp/v2/pkg/forgejo"
	"codeberg.org/goern/forgejo-mcp/v2/pkg/log"
	"codeberg.org/goern/forgejo-mcp/v2/pkg/to"
	forgejo_sdk "codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"

	"github.com/mark3labs/mcp-go/mcp"
)

const (
	SearchRunnerJobsToolName = "search_runner_jobs"
)

var (
	SearchRunnerJobsTool = mcp.NewTool(
		SearchRunnerJobsToolName,
		mcp.WithDescription("Search for repository's runner jobs according to filters"),
		mcp.WithString("owner", mcp.Required(), mcp.Description(params.Owner)),
		mcp.WithString("repo", mcp.Required(), mcp.Description(params.Repo)),
		mcp.WithString("labels", mcp.Description("Filter by job labels (comma-separated list)")),
	)
)

func SearchRunnerJobsFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called SearchRunnerJobsFn")

	owner, err := req.RequireString("owner")
	if err != nil {
		return to.ErrorResult(err)
	}
	repo, err := req.RequireString("repo")
	if err != nil {
		return to.ErrorResult(err)
	}

	opt := forgejo_sdk.SearchRunnerJobsOption{}

	// Optional filters
	if labels := req.GetString("labels", ""); labels != "" {
		opt.Labels = parseStringList(labels)
	}

	jobs, _, err := forgejo.Client().SearchRunnerJobs(owner, repo, opt)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("search runner jobs err: %v", err))
	}

	return to.TextResult(jobs)
}
