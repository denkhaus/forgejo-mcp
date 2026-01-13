package runner

import (
	"github.com/mark3labs/mcp-go/server"
)

// RegisterTool registers all runner-related tools with the MCP server
func RegisterTool(s *server.MCPServer) {
	// Action Run Tools
	s.AddTool(ListActionRunsTool, ListActionRunsFn)
	s.AddTool(GetActionRunTool, GetActionRunFn)

	// Runner Job Tools (repo level)
	s.AddTool(SearchRunnerJobsTool, SearchRunnerJobsFn)

	// Runner Job Tools (org/user level)
	s.AddTool(SearchOrgRunnerJobsTool, SearchOrgRunnerJobsFn)
	s.AddTool(SearchUserRunnerJobsTool, SearchUserRunnerJobsFn)

	// Registration Token Tools
	s.AddTool(GetRepoRunnerRegistrationTokenTool, GetRepoRunnerRegistrationTokenFn)
}
