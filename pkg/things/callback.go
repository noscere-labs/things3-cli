package things

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"
)

// CallbackServer handles receiving x-callback-url responses from Things via HTTP
// Things will request our local server with response parameters
// after completing an action.
type CallbackServer struct {
	Port     int
	server   *http.Server
	response chan map[string]string
	mu       sync.Mutex
	started  bool
}

// NewCallbackServer creates a new callback server instance
func NewCallbackServer(port int) *CallbackServer {
	return &CallbackServer{
		Port:     port,
		response: make(chan map[string]string, 1),
	}
}

// Start begins listening for x-callback responses
func (s *CallbackServer) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.started {
		return fmt.Errorf("callback server already started")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		params := make(map[string]string)
		for key, values := range r.URL.Query() {
			if len(values) > 0 {
				params[key] = values[0]
			}
		}

		select {
		case s.response <- params:
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Things CLI</title>
<style>
body{font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,sans-serif;
     display:flex;align-items:center;justify-content:center;height:100vh;
     margin:0;background:#f5f5f5;color:#333;}
.msg{text-align:center;padding:2rem;background:white;border-radius:8px;
     box-shadow:0 2px 8px rgba(0,0,0,0.1);}
h1{margin:0 0 0.5rem;font-size:1.5rem;color:#2563eb;}
p{margin:0;font-size:0.9rem;color:#666;}
</style>
</head>
<body>
<div class="msg">
<h1>✓ Success</h1>
<p>Things CLI callback received. You can close this tab.</p>
</div>
<script>
setTimeout(function(){window.close();},500);
setTimeout(function(){document.body.innerHTML='<div class="msg"><h1>✓ Success</h1><p>You can close this tab now.</p></div>';},600);
</script>
</body>
</html>`))
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to process response"))
		}
	})

	s.server = &http.Server{
		Addr:         fmt.Sprintf("localhost:%d", s.Port),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	ready := make(chan struct{})
	go func() {
		listener, err := net.Listen("tcp", s.server.Addr)
		if err != nil {
			close(ready)
			return
		}
		close(ready)

		if err := s.server.Serve(listener); err != nil && err != http.ErrServerClosed {
			// The callback may have already been received.
		}
	}()

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

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown callback server: %w", err)
	}

	s.started = false
	return nil
}

// WaitForResponse blocks until a response is received from Things or timeout occurs
func (s *CallbackServer) WaitForResponse(timeout time.Duration) (map[string]string, error) {
	select {
	case response := <-s.response:
		return response, nil
	case <-time.After(timeout):
		return nil, fmt.Errorf("callback timeout: no response from Things within %v", timeout)
	}
}

// IsPortAvailable checks if the given port is available for listening
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
func FindAvailablePort(startPort int) int {
	for port := startPort; port < startPort+100; port++ {
		if IsPortAvailable(port) {
			return port
		}
	}
	return -1
}
