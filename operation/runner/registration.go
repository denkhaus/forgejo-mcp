package runner

import (
	"context"
	"fmt"

	"codeberg.org/goern/forgejo-mcp/v2/operation/params"
	"codeberg.org/goern/forgejo-mcp/v2/pkg/forgejo"
	"codeberg.org/goern/forgejo-mcp/v2/pkg/log"
	"codeberg.org/goern/forgejo-mcp/v2/pkg/to"

	"github.com/mark3labs/mcp-go/mcp"
)

const (
	GetRepoRunnerRegistrationTokenToolName = "get_repo_runner_registration_token"
)

var (
	GetRepoRunnerRegistrationTokenTool = mcp.NewTool(
		GetRepoRunnerRegistrationTokenToolName,
		mcp.WithDescription("Get a repository's runner registration token for setting up runners"),
		mcp.WithString("owner", mcp.Required(), mcp.Description(params.Owner)),
		mcp.WithString("repo", mcp.Required(), mcp.Description(params.Repo)),
	)
)

func GetRepoRunnerRegistrationTokenFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called GetRepoRunnerRegistrationTokenFn")

	owner, err := req.RequireString("owner")
	if err != nil {
		return to.ErrorResult(err)
	}
	repo, err := req.RequireString("repo")
	if err != nil {
		return to.ErrorResult(err)
	}

	token, _, err := forgejo.Client().GetRepoRunnerRegistrationToken(owner, repo)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get repo runner registration token err: %v", err))
	}

	return to.TextResult(token)
}
