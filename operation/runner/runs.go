package runner

import (
	"context"
	"fmt"
	"strings"

	"codeberg.org/goern/forgejo-mcp/v2/operation/params"
	"codeberg.org/goern/forgejo-mcp/v2/pkg/forgejo"
	"codeberg.org/goern/forgejo-mcp/v2/pkg/log"
	"codeberg.org/goern/forgejo-mcp/v2/pkg/to"
	forgejo_sdk "codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"

	"github.com/mark3labs/mcp-go/mcp"
)

const (
	ListActionRunsToolName = "list_action_runs"
	GetActionRunToolName   = "get_action_run"
)

var (
	ListActionRunsTool = mcp.NewTool(
		ListActionRunsToolName,
		mcp.WithDescription("List a repository's action runs"),
		mcp.WithString("owner", mcp.Required(), mcp.Description(params.Owner)),
		mcp.WithString("repo", mcp.Required(), mcp.Description(params.Repo)),
		mcp.WithNumber("page", mcp.Description(params.Page), mcp.DefaultNumber(1), mcp.Min(1)),
		mcp.WithNumber("limit", mcp.Description(params.Limit), mcp.DefaultNumber(50), mcp.Min(1), mcp.Max(100)),
		mcp.WithString("status", mcp.Description("Filter by run status (unknown, waiting, running, success, failure, cancelled, skipped, blocked)")),
		mcp.WithString("events", mcp.Description("Filter by workflow event (e.g., push, pull_request, workflow_dispatch)")),
		mcp.WithNumber("run_number", mcp.Description("Filter by run number")),
		mcp.WithString("head_sha", mcp.Description("Filter by head commit SHA")),
	)

	GetActionRunTool = mcp.NewTool(
		GetActionRunToolName,
		mcp.WithDescription("Get details of a specific action run"),
		mcp.WithString("owner", mcp.Required(), mcp.Description(params.Owner)),
		mcp.WithString("repo", mcp.Required(), mcp.Description(params.Repo)),
		mcp.WithNumber("run_id", mcp.Required(), mcp.Description("Action run ID")),
	)
)

func ListActionRunsFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called ListActionRunsFn")

	owner, err := req.RequireString("owner")
	if err != nil {
		return to.ErrorResult(err)
	}
	repo, err := req.RequireString("repo")
	if err != nil {
		return to.ErrorResult(err)
	}

	page := int(req.GetFloat("page", 1))
	limit := int(req.GetFloat("limit", 50))

	opt := forgejo_sdk.ListActionRunsOption{
		ListOptions: forgejo_sdk.ListOptions{
			Page:     page,
			PageSize: limit,
		},
	}

	// Optional filters
	if status := req.GetString("status", ""); status != "" {
		opt.Status = parseStringList(status)
	}
	if events := req.GetString("events", ""); events != "" {
		opt.Events = parseStringList(events)
	}
	if runNumber := req.GetFloat("run_number", 0); runNumber > 0 {
		opt.RunNumber = int64(runNumber)
	}
	if headSHA := req.GetString("head_sha", ""); headSHA != "" {
		opt.HeadSHA = headSHA
	}

	runs, _, err := forgejo.Client().ListActionRuns(owner, repo, opt)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("list action runs err: %v", err))
	}

	return to.TextResult(runs)
}

func GetActionRunFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called GetActionRunFn")

	owner, err := req.RequireString("owner")
	if err != nil {
		return to.ErrorResult(err)
	}
	repo, err := req.RequireString("repo")
	if err != nil {
		return to.ErrorResult(err)
	}
	runID, err := req.RequireFloat("run_id")
	if err != nil {
		return to.ErrorResult(err)
	}

	run, _, err := forgejo.Client().GetActionRun(owner, repo, int64(runID))
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get action run err: %v", err))
	}

	return to.TextResult(run)
}

// parseStringList parses a comma-separated string into a slice
func parseStringList(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
