package bear

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"
)

// CallbackServer handles receiving x-callback-url responses from Bear via HTTP
// Bear will POST/GET to our local server with response parameters
type CallbackServer struct {
	Port     int                   // Port to listen on (e.g., 8765)
	server   *http.Server          // HTTP server instance
	response chan map[string]string // Channel to pass response back to caller
	mu       sync.Mutex            // Mutex to protect state
	started  bool                  // Whether server has been started
}

// NewCallbackServer creates a new callback server instance
// port: The port to listen on (e.g., 8765)
func NewCallbackServer(port int) *CallbackServer {
	return &CallbackServer{
		Port:     port,
		response: make(chan map[string]string, 1),
	}
}

// Start begins listening for x-callback responses
// It sets up an HTTP server that will receive responses from Bear
func (s *CallbackServer) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.started {
		return fmt.Errorf("callback server already started")
	}

	// Setup HTTP handlers for the callback endpoint
	mux := http.NewServeMux()

	// The callback handler receives the x-callback-url response
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		// Parse all query parameters from the callback URL
		params := make(map[string]string)
		for key, values := range r.URL.Query() {
			if len(values) > 0 {
				params[key] = values[0]
			}
		}

		// Send response back to the waiting caller
		select {
		case s.response <- params:
			// Successfully sent to channel
			// Return minimal HTML that tries to close and shows friendly message
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Bear CLI</title>
<style>
body{font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,sans-serif;
     display:flex;align-items:center;justify-content:center;height:100vh;
     margin:0;background:#f5f5f5;color:#333;}
.msg{text-align:center;padding:2rem;background:white;border-radius:8px;
     box-shadow:0 2px 8px rgba(0,0,0,0.1);}
h1{margin:0 0 0.5rem;font-size:1.5rem;color:#059669;}
p{margin:0;font-size:0.9rem;color:#666;}
</style>
</head>
<body>
<div class="msg">
<h1>✓ Success</h1>
<p>Bear CLI callback received. This window will close automatically.</p>
</div>
<script>
setTimeout(function(){window.close();},500);
setTimeout(function(){document.body.innerHTML='<div class="msg"><h1>✓ Success</h1><p>You can close this tab now.</p></div>';},600);
</script>
</body>
</html>`))
		default:
			// Channel not ready (shouldn't happen)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to process response"))
		}
	})

	// Create the HTTP server
	s.server = &http.Server{
		Addr:         fmt.Sprintf("localhost:%d", s.Port),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	// Channel to signal server is ready
	ready := make(chan struct{})

	// Start listening in a goroutine so we don't block
	go func() {
		// Create listener first to ensure port is bound
		listener, err := net.Listen("tcp", s.server.Addr)
		if err != nil {
			close(ready)
			return
		}
		close(ready) // Signal that we're ready to accept connections

		if err := s.server.Serve(listener); err != nil && err != http.ErrServerClosed {
			// Server error - callback might have already been received
		}
	}()

	// Wait for server to be ready
	<-ready

	s.started = true
	return nil
}

// Stop shuts down the callback server
func (s *CallbackServer) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.started || s.server == nil {
		return nil
	}

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown callback server: %w", err)
	}

	s.started = false
	return nil
}

// WaitForResponse blocks until a response is received from Bear or timeout occurs
// timeout: Maximum time to wait for a response
// Returns the parsed query parameters from the x-success URL
func (s *CallbackServer) WaitForResponse(timeout time.Duration) (map[string]string, error) {
	select {
	case response := <-s.response:
		return response, nil
	case <-time.After(timeout):
		// No response received within the timeout period
		return nil, fmt.Errorf("callback timeout: no response from Bear within %v", timeout)
	}
}

// IsPortAvailable checks if the given port is available for listening
// This is useful for selecting an alternative port if the default is in use
func IsPortAvailable(port int) bool {
	addr := fmt.Sprintf("localhost:%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return false
	}
	listener.Close()
	return true
}

// FindAvailablePort finds an available port starting from the given port
// Useful as a fallback if the default port is already in use
func FindAvailablePort(startPort int) int {
	for port := startPort; port < startPort+100; port++ {
		if IsPortAvailable(port) {
			return port
		}
	}
	return -1 // No available port found
}
