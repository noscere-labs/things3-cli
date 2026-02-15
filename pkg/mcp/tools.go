package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	gomcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/yourusername/things3-cli/pkg/things"
)

func executeTool(client *things.Client, action string, params map[string]string, opts things.ExecuteOptions) (*gomcp.CallToolResult, error) {
	callback, err := client.Execute(action, params, opts)
	if err != nil {
		return &gomcp.CallToolResult{
			Content: []gomcp.Content{&gomcp.TextContent{Text: fmt.Sprintf("Error: %v", err)}},
			IsError: true,
		}, nil
	}

	result := things.NormalizeResponse(action, callback)
	data, _ := json.MarshalIndent(result, "", "  ")
	return &gomcp.CallToolResult{
		Content: []gomcp.Content{&gomcp.TextContent{Text: string(data)}},
	}, nil
}

func setIfNonEmpty(params map[string]string, key, value string) {
	if value != "" {
		params[key] = value
	}
}

type AddInput struct {
	Title          string `json:"title,omitempty" jsonschema:"description=To-do title"`
	Titles         string `json:"titles,omitempty" jsonschema:"description=Newline-separated list of to-do titles (for batch creation)"`
	Notes          string `json:"notes,omitempty" jsonschema:"description=Notes for the to-do"`
	When           string `json:"when,omitempty" jsonschema:"description=When to schedule: today, tonight, anytime, someday, or YYYY-MM-DD"`
	Deadline       string `json:"deadline,omitempty" jsonschema:"description=Deadline date (YYYY-MM-DD)"`
	Tags           string `json:"tags,omitempty" jsonschema:"description=Comma-separated tags"`
	List           string `json:"list,omitempty" jsonschema:"description=List name or project title"`
	ListID         string `json:"list_id,omitempty" jsonschema:"description=List or project ID"`
	Heading        string `json:"heading,omitempty" jsonschema:"description=Heading title"`
	HeadingID      string `json:"heading_id,omitempty" jsonschema:"description=Heading ID"`
	ChecklistItems string `json:"checklist_items,omitempty" jsonschema:"description=Newline-separated checklist items"`
	Completed      bool   `json:"completed,omitempty" jsonschema:"description=Mark as completed"`
	Canceled       bool   `json:"canceled,omitempty" jsonschema:"description=Mark as canceled"`
	ShowQuickEntry bool   `json:"show_quick_entry,omitempty" jsonschema:"description=Show quick entry after adding"`
	Reveal         bool   `json:"reveal,omitempty" jsonschema:"description=Reveal the created to-do in Things"`
	CreationDate   string `json:"creation_date,omitempty" jsonschema:"description=Creation date (ISO 8601)"`
	CompletionDate string `json:"completion_date,omitempty" jsonschema:"description=Completion date (ISO 8601)"`
}

type AddProjectInput struct {
	Title          string `json:"title,omitempty" jsonschema:"description=Project title"`
	Notes          string `json:"notes,omitempty" jsonschema:"description=Project notes"`
	When           string `json:"when,omitempty" jsonschema:"description=When to schedule: today, tonight, anytime, someday, or YYYY-MM-DD"`
	Deadline       string `json:"deadline,omitempty" jsonschema:"description=Deadline date (YYYY-MM-DD)"`
	Tags           string `json:"tags,omitempty" jsonschema:"description=Comma-separated tags"`
	Area           string `json:"area,omitempty" jsonschema:"description=Area name"`
	AreaID         string `json:"area_id,omitempty" jsonschema:"description=Area ID"`
	ToDos          string `json:"to_dos,omitempty" jsonschema:"description=Newline-separated to-do titles for the project"`
	Completed      bool   `json:"completed,omitempty" jsonschema:"description=Mark as completed"`
	Canceled       bool   `json:"canceled,omitempty" jsonschema:"description=Mark as canceled"`
	Reveal         bool   `json:"reveal,omitempty" jsonschema:"description=Reveal the created project in Things"`
	CreationDate   string `json:"creation_date,omitempty" jsonschema:"description=Creation date (ISO 8601)"`
	CompletionDate string `json:"completion_date,omitempty" jsonschema:"description=Completion date (ISO 8601)"`
}

type UpdateInput struct {
	ID                    string `json:"id" jsonschema:"description=To-do ID (required),required"`
	Title                 string `json:"title,omitempty" jsonschema:"description=Updated title"`
	Notes                 string `json:"notes,omitempty" jsonschema:"description=Replace notes"`
	PrependNotes          string `json:"prepend_notes,omitempty" jsonschema:"description=Prepend to notes"`
	AppendNotes           string `json:"append_notes,omitempty" jsonschema:"description=Append to notes"`
	When                  string `json:"when,omitempty" jsonschema:"description=Update schedule"`
	Deadline              string `json:"deadline,omitempty" jsonschema:"description=Update deadline"`
	Tags                  string `json:"tags,omitempty" jsonschema:"description=Replace tags (comma-separated)"`
	AddTags               string `json:"add_tags,omitempty" jsonschema:"description=Add tags (comma-separated)"`
	ChecklistItems        string `json:"checklist_items,omitempty" jsonschema:"description=Replace checklist items (newline-separated)"`
	PrependChecklistItems string `json:"prepend_checklist_items,omitempty" jsonschema:"description=Prepend checklist items (newline-separated)"`
	AppendChecklistItems  string `json:"append_checklist_items,omitempty" jsonschema:"description=Append checklist items (newline-separated)"`
	List                  string `json:"list,omitempty" jsonschema:"description=Move to list by name"`
	ListID                string `json:"list_id,omitempty" jsonschema:"description=Move to list by ID"`
	Heading               string `json:"heading,omitempty" jsonschema:"description=Move to heading by name"`
	HeadingID             string `json:"heading_id,omitempty" jsonschema:"description=Move to heading by ID"`
	Completed             bool   `json:"completed,omitempty" jsonschema:"description=Mark as completed"`
	Canceled              bool   `json:"canceled,omitempty" jsonschema:"description=Mark as canceled"`
	Reveal                bool   `json:"reveal,omitempty" jsonschema:"description=Reveal the updated to-do"`
	Duplicate             bool   `json:"duplicate,omitempty" jsonschema:"description=Duplicate the to-do"`
	CreationDate          string `json:"creation_date,omitempty" jsonschema:"description=Set creation date (ISO 8601)"`
	CompletionDate        string `json:"completion_date,omitempty" jsonschema:"description=Set completion date (ISO 8601)"`
}

type UpdateProjectInput struct {
	ID             string `json:"id" jsonschema:"description=Project ID (required),required"`
	Title          string `json:"title,omitempty" jsonschema:"description=Updated title"`
	Notes          string `json:"notes,omitempty" jsonschema:"description=Replace notes"`
	PrependNotes   string `json:"prepend_notes,omitempty" jsonschema:"description=Prepend to notes"`
	AppendNotes    string `json:"append_notes,omitempty" jsonschema:"description=Append to notes"`
	When           string `json:"when,omitempty" jsonschema:"description=Update schedule"`
	Deadline       string `json:"deadline,omitempty" jsonschema:"description=Update deadline"`
	Tags           string `json:"tags,omitempty" jsonschema:"description=Replace tags (comma-separated)"`
	AddTags        string `json:"add_tags,omitempty" jsonschema:"description=Add tags (comma-separated)"`
	Area           string `json:"area,omitempty" jsonschema:"description=Move to area by name"`
	AreaID         string `json:"area_id,omitempty" jsonschema:"description=Move to area by ID"`
	Completed      bool   `json:"completed,omitempty" jsonschema:"description=Mark as completed"`
	Canceled       bool   `json:"canceled,omitempty" jsonschema:"description=Mark as canceled"`
	Reveal         bool   `json:"reveal,omitempty" jsonschema:"description=Reveal the updated project"`
	Duplicate      bool   `json:"duplicate,omitempty" jsonschema:"description=Duplicate the project"`
	CreationDate   string `json:"creation_date,omitempty" jsonschema:"description=Set creation date (ISO 8601)"`
	CompletionDate string `json:"completion_date,omitempty" jsonschema:"description=Set completion date (ISO 8601)"`
}

type ShowInput struct {
	ID    string `json:"id,omitempty" jsonschema:"description=Item ID to show"`
	Query string `json:"query,omitempty" jsonschema:"description=List query: Inbox, Today, Upcoming, Anytime, Someday, Logbook"`
}

type SearchInput struct {
	Query string `json:"query" jsonschema:"description=Search query,required"`
}

type JSONInput struct {
	Data   string `json:"data" jsonschema:"description=JSON payload string for Things batch operations,required"`
	Reveal bool   `json:"reveal,omitempty" jsonschema:"description=Reveal created items"`
}

type VersionInput struct{}

func makeAddHandler(client *things.Client) func(context.Context, *gomcp.CallToolRequest, AddInput) (*gomcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *gomcp.CallToolRequest, input AddInput) (*gomcp.CallToolResult, any, error) {
		params := make(map[string]string)
		if input.Titles != "" {
			params["titles"] = input.Titles
		} else {
			setIfNonEmpty(params, "title", input.Title)
		}
		setIfNonEmpty(params, "notes", input.Notes)
		setIfNonEmpty(params, "when", input.When)
		setIfNonEmpty(params, "deadline", input.Deadline)
		setIfNonEmpty(params, "tags", input.Tags)
		setIfNonEmpty(params, "list", input.List)
		setIfNonEmpty(params, "list-id", input.ListID)
		setIfNonEmpty(params, "heading", input.Heading)
		setIfNonEmpty(params, "heading-id", input.HeadingID)
		setIfNonEmpty(params, "checklist-items", input.ChecklistItems)
		setIfNonEmpty(params, "creation-date", input.CreationDate)
		setIfNonEmpty(params, "completion-date", input.CompletionDate)
		if input.Completed {
			params["completed"] = "true"
		}
		if input.Canceled {
			params["canceled"] = "true"
		}
		if input.ShowQuickEntry {
			params["show-quick-entry"] = "true"
		}
		if input.Reveal {
			params["reveal"] = "true"
		}
		result, err := executeTool(client, "add", params, things.ExecuteOptions{})
		return result, nil, err
	}
}

func makeAddProjectHandler(client *things.Client) func(context.Context, *gomcp.CallToolRequest, AddProjectInput) (*gomcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *gomcp.CallToolRequest, input AddProjectInput) (*gomcp.CallToolResult, any, error) {
		params := make(map[string]string)
		setIfNonEmpty(params, "title", input.Title)
		setIfNonEmpty(params, "notes", input.Notes)
		setIfNonEmpty(params, "when", input.When)
		setIfNonEmpty(params, "deadline", input.Deadline)
		setIfNonEmpty(params, "tags", input.Tags)
		setIfNonEmpty(params, "area", input.Area)
		setIfNonEmpty(params, "area-id", input.AreaID)
		setIfNonEmpty(params, "to-dos", input.ToDos)
		setIfNonEmpty(params, "creation-date", input.CreationDate)
		setIfNonEmpty(params, "completion-date", input.CompletionDate)
		if input.Completed {
			params["completed"] = "true"
		}
		if input.Canceled {
			params["canceled"] = "true"
		}
		if input.Reveal {
			params["reveal"] = "true"
		}
		result, err := executeTool(client, "add-project", params, things.ExecuteOptions{})
		return result, nil, err
	}
}

func makeUpdateHandler(client *things.Client) func(context.Context, *gomcp.CallToolRequest, UpdateInput) (*gomcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *gomcp.CallToolRequest, input UpdateInput) (*gomcp.CallToolResult, any, error) {
		if input.ID == "" {
			return &gomcp.CallToolResult{
				Content: []gomcp.Content{&gomcp.TextContent{Text: "Error: id is required"}},
				IsError: true,
			}, nil, nil
		}
		params := map[string]string{"id": input.ID}
		setIfNonEmpty(params, "title", input.Title)
		setIfNonEmpty(params, "notes", input.Notes)
		setIfNonEmpty(params, "prepend-notes", input.PrependNotes)
		setIfNonEmpty(params, "append-notes", input.AppendNotes)
		setIfNonEmpty(params, "when", input.When)
		setIfNonEmpty(params, "deadline", input.Deadline)
		setIfNonEmpty(params, "tags", input.Tags)
		setIfNonEmpty(params, "add-tags", input.AddTags)
		setIfNonEmpty(params, "checklist-items", input.ChecklistItems)
		setIfNonEmpty(params, "prepend-checklist-items", input.PrependChecklistItems)
		setIfNonEmpty(params, "append-checklist-items", input.AppendChecklistItems)
		setIfNonEmpty(params, "list", input.List)
		setIfNonEmpty(params, "list-id", input.ListID)
		setIfNonEmpty(params, "heading", input.Heading)
		setIfNonEmpty(params, "heading-id", input.HeadingID)
		setIfNonEmpty(params, "creation-date", input.CreationDate)
		setIfNonEmpty(params, "completion-date", input.CompletionDate)
		if input.Completed {
			params["completed"] = "true"
		}
		if input.Canceled {
			params["canceled"] = "true"
		}
		if input.Reveal {
			params["reveal"] = "true"
		}
		if input.Duplicate {
			params["duplicate"] = "true"
		}
		result, err := executeTool(client, "update", params, things.ExecuteOptions{RequiresAuth: true, UseAuthIfAvailable: true})
		return result, nil, err
	}
}

func makeUpdateProjectHandler(client *things.Client) func(context.Context, *gomcp.CallToolRequest, UpdateProjectInput) (*gomcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *gomcp.CallToolRequest, input UpdateProjectInput) (*gomcp.CallToolResult, any, error) {
		if input.ID == "" {
			return &gomcp.CallToolResult{
				Content: []gomcp.Content{&gomcp.TextContent{Text: "Error: id is required"}},
				IsError: true,
			}, nil, nil
		}
		params := map[string]string{"id": input.ID}
		setIfNonEmpty(params, "title", input.Title)
		setIfNonEmpty(params, "notes", input.Notes)
		setIfNonEmpty(params, "prepend-notes", input.PrependNotes)
		setIfNonEmpty(params, "append-notes", input.AppendNotes)
		setIfNonEmpty(params, "when", input.When)
		setIfNonEmpty(params, "deadline", input.Deadline)
		setIfNonEmpty(params, "tags", input.Tags)
		setIfNonEmpty(params, "add-tags", input.AddTags)
		setIfNonEmpty(params, "area", input.Area)
		setIfNonEmpty(params, "area-id", input.AreaID)
		setIfNonEmpty(params, "creation-date", input.CreationDate)
		setIfNonEmpty(params, "completion-date", input.CompletionDate)
		if input.Completed {
			params["completed"] = "true"
		}
		if input.Canceled {
			params["canceled"] = "true"
		}
		if input.Reveal {
			params["reveal"] = "true"
		}
		if input.Duplicate {
			params["duplicate"] = "true"
		}
		result, err := executeTool(client, "update-project", params, things.ExecuteOptions{RequiresAuth: true, UseAuthIfAvailable: true})
		return result, nil, err
	}
}

func makeShowHandler(client *things.Client) func(context.Context, *gomcp.CallToolRequest, ShowInput) (*gomcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *gomcp.CallToolRequest, input ShowInput) (*gomcp.CallToolResult, any, error) {
		params := make(map[string]string)
		setIfNonEmpty(params, "id", input.ID)
		setIfNonEmpty(params, "query", input.Query)
		if len(params) == 0 {
			return &gomcp.CallToolResult{
				Content: []gomcp.Content{&gomcp.TextContent{Text: "Error: provide id or query"}},
				IsError: true,
			}, nil, nil
		}
		result, err := executeTool(client, "show", params, things.ExecuteOptions{})
		return result, nil, err
	}
}

func makeSearchHandler(client *things.Client) func(context.Context, *gomcp.CallToolRequest, SearchInput) (*gomcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *gomcp.CallToolRequest, input SearchInput) (*gomcp.CallToolResult, any, error) {
		if input.Query == "" {
			return &gomcp.CallToolResult{
				Content: []gomcp.Content{&gomcp.TextContent{Text: "Error: query is required"}},
				IsError: true,
			}, nil, nil
		}
		params := map[string]string{"query": input.Query}
		result, err := executeTool(client, "search", params, things.ExecuteOptions{})
		return result, nil, err
	}
}

func makeJSONHandler(client *things.Client) func(context.Context, *gomcp.CallToolRequest, JSONInput) (*gomcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *gomcp.CallToolRequest, input JSONInput) (*gomcp.CallToolResult, any, error) {
		if input.Data == "" {
			return &gomcp.CallToolResult{
				Content: []gomcp.Content{&gomcp.TextContent{Text: "Error: data is required"}},
				IsError: true,
			}, nil, nil
		}
		params := map[string]string{"data": input.Data}
		if input.Reveal {
			params["reveal"] = "true"
		}
		result, err := executeTool(client, "json", params, things.ExecuteOptions{UseAuthIfAvailable: true})
		return result, nil, err
	}
}

func makeVersionHandler(client *things.Client) func(context.Context, *gomcp.CallToolRequest, VersionInput) (*gomcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *gomcp.CallToolRequest, input VersionInput) (*gomcp.CallToolResult, any, error) {
		result, err := executeTool(client, "version", map[string]string{}, things.ExecuteOptions{})
		return result, nil, err
	}
}
