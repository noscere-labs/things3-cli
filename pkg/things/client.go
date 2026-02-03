package things

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/yourusername/things3-cli/pkg/util"
)

// Client handles communication with Things via the URL scheme.
type Client struct {
	AuthToken    string
	CallbackPort int
	timeout      time.Duration
}

// ExecuteOptions controls how actions are executed.
type ExecuteOptions struct {
	RequiresAuth      bool
	UseAuthIfAvailable bool
}

// CallbackError represents an error returned via the callback URL.
type CallbackError struct {
	Code     string
	Message  string
	Callback map[string]string
}

func (e *CallbackError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("%s (%s)", e.Message, e.Code)
	}
	return e.Message
}

// NewClient creates a new Things client with default settings.
func NewClient() (*Client, error) {
	token, err := util.GetAuthToken()
	if err != nil {
		token = ""
	}

	config, err := util.LoadConfig()
	if err != nil {
		config = util.DefaultConfig()
	}

	return &Client{
		AuthToken:    token,
		CallbackPort: config.CallbackPort,
		timeout:      time.Duration(config.CallbackTimeoutSeconds) * time.Second,
	}, nil
}

// buildThingsURL constructs a Things URL scheme invocation.
func (c *Client) buildThingsURL(action string, params map[string]string) string {
	baseURL := fmt.Sprintf("things:///%s", action)
	queryStr := util.EncodeParams(params)
	if queryStr == "" {
		return baseURL
	}
	return baseURL + "?" + queryStr
}

// Execute runs the given Things action and returns the callback response.
func (c *Client) Execute(action string, params map[string]string, opts ExecuteOptions) (map[string]string, error) {
	if params == nil {
		params = make(map[string]string)
	}

	if params["auth-token"] == "" && (opts.RequiresAuth || opts.UseAuthIfAvailable) {
		if c.AuthToken != "" {
			params["auth-token"] = c.AuthToken
		} else if opts.RequiresAuth {
			return nil, fmt.Errorf("auth token required (set with things config set-token or THINGS_AUTH_TOKEN)")
		}
	}

	port := c.CallbackPort
	if !IsPortAvailable(port) {
		alt := FindAvailablePort(port + 1)
		if alt < 0 {
			return nil, fmt.Errorf("no available callback port found")
		}
		port = alt
	}

	params["x-success"] = fmt.Sprintf("http://localhost:%d/callback?result=success", port)
	params["x-error"] = fmt.Sprintf("http://localhost:%d/callback?result=error", port)

	callbackServer := NewCallbackServer(port)
	if err := callbackServer.Start(); err != nil {
		return nil, fmt.Errorf("failed to start callback server: %w", err)
	}
	defer callbackServer.Stop()

	thingsURL := c.buildThingsURL(action, params)
	cmd := exec.Command("open", thingsURL)
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to execute Things URL: %w", err)
	}

	response, err := callbackServer.WaitForResponse(c.timeout)
	if err != nil {
		return nil, err
	}

	if response["result"] == "error" {
		code := response["errorCode"]
		message := response["errorMessage"]
		if message == "" {
			message = "Things returned an error"
		}
		return response, &CallbackError{Code: code, Message: message, Callback: response}
	}

	return response, nil
}

// NormalizeResponse produces a structured result from a callback response.
func NormalizeResponse(action string, callback map[string]string) ActionResult {
	result := ActionResult{Action: action}
	if len(callback) == 0 {
		return result
	}

	cleaned := make(map[string]string)
	for key, value := range callback {
		if key == "result" {
			continue
		}
		cleaned[key] = value
	}
	result.Callback = cleaned

	if ids := callback["x-things-ids"]; ids != "" {
		var parsed []string
		if err := json.Unmarshal([]byte(ids), &parsed); err == nil {
			result.ThingsIDs = parsed
		}
	}

	if id := callback["x-things-id"]; id != "" {
		if strings.Contains(id, ",") {
			for _, part := range strings.Split(id, ",") {
				trimmed := strings.TrimSpace(part)
				if trimmed != "" {
					result.ThingsIDs = append(result.ThingsIDs, trimmed)
				}
			}
		} else {
			result.ThingsID = id
		}
	}

	if callback["x-things-scheme-version"] != "" {
		result.ThingsSchemeVersion = callback["x-things-scheme-version"]
	}
	if callback["x-things-client-version"] != "" {
		result.ThingsClientVersion = callback["x-things-client-version"]
	}

	return result
}
