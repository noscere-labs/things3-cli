package mcp

import (
	"fmt"
	"log"
	"net/http"

	gomcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/yourusername/things3-cli/pkg/things"
)

func NewThingsServer() (*gomcp.Server, error) {
	client, err := things.NewClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create Things client: %w", err)
	}

	server := gomcp.NewServer(
		&gomcp.Implementation{
			Name:    "things3",
			Version: "1.0.0",
		},
		nil,
	)

	gomcp.AddTool(server, &gomcp.Tool{
		Name:        "things_add",
		Description: "Add a new to-do in Things 3. Supports title, notes, tags, scheduling, checklist items, and more.",
	}, makeAddHandler(client))

	gomcp.AddTool(server, &gomcp.Tool{
		Name:        "things_add_project",
		Description: "Add a new project in Things 3. Supports title, notes, tags, area, and initial to-dos.",
	}, makeAddProjectHandler(client))

	gomcp.AddTool(server, &gomcp.Tool{
		Name:        "things_update",
		Description: "Update an existing to-do in Things 3 by ID. Requires an auth token to be configured.",
	}, makeUpdateHandler(client))

	gomcp.AddTool(server, &gomcp.Tool{
		Name:        "things_update_project",
		Description: "Update an existing project in Things 3 by ID. Requires an auth token to be configured.",
	}, makeUpdateProjectHandler(client))

	gomcp.AddTool(server, &gomcp.Tool{
		Name:        "things_show",
		Description: "Show a list or specific item in Things 3. Use query for lists (Inbox, Today, Upcoming, Anytime, Someday, Logbook) or id for a specific item.",
	}, makeShowHandler(client))

	gomcp.AddTool(server, &gomcp.Tool{
		Name:        "things_search",
		Description: "Search for items in Things 3 using a text query.",
	}, makeSearchHandler(client))

	gomcp.AddTool(server, &gomcp.Tool{
		Name:        "things_json",
		Description: "Send a JSON payload to Things 3 for batch creation or updates. See Things URL scheme docs for payload format.",
	}, makeJSONHandler(client))

	gomcp.AddTool(server, &gomcp.Tool{
		Name:        "things_version",
		Description: "Get the Things URL scheme version and client version.",
	}, makeVersionHandler(client))

	return server, nil
}

func Serve(port int) error {
	server, err := NewThingsServer()
	if err != nil {
		return err
	}

	handler := gomcp.NewStreamableHTTPHandler(func(r *http.Request) *gomcp.Server {
		return server
	}, nil)

	addr := fmt.Sprintf(":%d", port)
	log.Printf("Things MCP server listening on http://localhost:%d/mcp", port)

	mux := http.NewServeMux()
	mux.Handle("/mcp", handler)

	return http.ListenAndServe(addr, mux)
}
