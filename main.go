package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"obsidian-mcp/api"
	"obsidian-mcp/security"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

// Config represents the server configuration
type Config struct {
	ObsidianAPI struct {
		BaseURL string `yaml:"base_url"`
		Token   string `yaml:"token"`
		Port    int    `yaml:"port"`
	} `yaml:"obsidian_api"`
	MCP struct {
		Port        int    `yaml:"port"`
		Host        string `yaml:"host"`
		Description string `yaml:"description"`
	} `yaml:"mcp"`
	Security security.SecurityConfig `yaml:"security"`
}

// MCPServer represents the MCP server
type MCPServer struct {
	config      Config
	obsidianAPI *api.ObsidianAPI
	rateLimiter *security.RateLimiter
}

// MCPResponse represents a standard MCP response
type MCPResponse struct {
	Jsonrpc string      `json:"jsonrpc"`
	ID      interface{} `json:"id,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

// MCPError represents an MCP error
type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// MCPRequest represents an MCP request
type MCPRequest struct {
	Jsonrpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
	ID      interface{} `json:"id,omitempty"`
}

func loadConfig() (Config, error) {
	var config Config

	// Default configuration
	config.ObsidianAPI.BaseURL = "http://localhost:27123"
	config.ObsidianAPI.Port = 27123
	config.MCP.Port = 8080
	config.MCP.Host = "localhost"
	config.MCP.Description = "Obsidian MCP Server - Access and manage your Obsidian vault"

	// Try to load from config file
	configPath := "config.yaml"
	if _, err := os.Stat(configPath); err == nil {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return config, err
		}

		err = yaml.Unmarshal(data, &config)
		if err != nil {
			return config, err
		}
	}

	// Override with environment variables if set
	if token := os.Getenv("OBSIDIAN_API_TOKEN"); token != "" {
		config.ObsidianAPI.Token = token
	}
	if baseURL := os.Getenv("OBSIDIAN_API_BASE_URL"); baseURL != "" {
		config.ObsidianAPI.BaseURL = baseURL
	}

	return config, nil
}

func (s *MCPServer) handleInitialize(w http.ResponseWriter, r *http.Request) {
	var req MCPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendError(w, req.ID, -32700, "Parse error")
		return
	}

	result := map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities": map[string]interface{}{
			"tools": map[string]interface{}{
				"listChanged": false,
			},
			"resources": map[string]interface{}{
				"subscribe":   false,
				"listChanged": false,
			},
		},
		"serverInfo": map[string]interface{}{
			"name":    "obsidian-mcp-server",
			"version": "1.0.0",
		},
	}

	s.sendResponse(w, req.ID, result)
}

func (s *MCPServer) handleListTools(w http.ResponseWriter, r *http.Request) {
	var req MCPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendError(w, req.ID, -32700, "Parse error")
		return
	}

	tools := []map[string]interface{}{
		{
			"name":        "get_note",
			"description": "Get the content of a note by its path",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"path": map[string]interface{}{
						"type":        "string",
						"description": "Path to the note (relative to vault root)",
					},
				},
				"required": []string{"path"},
			},
		},
		{
			"name":        "create_note",
			"description": "Create a new note with the specified content",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"path": map[string]interface{}{
						"type":        "string",
						"description": "Path where the note should be created",
					},
					"content": map[string]interface{}{
						"type":        "string",
						"description": "Content of the note",
					},
				},
				"required": []string{"path", "content"},
			},
		},
		{
			"name":        "update_note",
			"description": "Update an existing note's content",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"path": map[string]interface{}{
						"type":        "string",
						"description": "Path to the note to update",
					},
					"content": map[string]interface{}{
						"type":        "string",
						"description": "New content for the note",
					},
				},
				"required": []string{"path", "content"},
			},
		},
		{
			"name":        "delete_note",
			"description": "Delete a note",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"path": map[string]interface{}{
						"type":        "string",
						"description": "Path to the note to delete",
					},
				},
				"required": []string{"path"},
			},
		},
		{
			"name":        "list_notes",
			"description": "List all notes in the vault",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"folder": map[string]interface{}{
						"type":        "string",
						"description": "Optional folder to filter by",
					},
				},
			},
		},
		{
			"name":        "search_notes",
			"description": "Search for notes containing specific text",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query": map[string]interface{}{
						"type":        "string",
						"description": "Search query",
					},
				},
				"required": []string{"query"},
			},
		},
		{
			"name":        "get_vault_info",
			"description": "Get information about the vault",
			"inputSchema": map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
	}

	result := map[string]interface{}{
		"tools": tools,
	}

	s.sendResponse(w, req.ID, result)
}

func (s *MCPServer) handleCallTool(w http.ResponseWriter, r *http.Request) {
	var req MCPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendError(w, req.ID, -32700, "Parse error")
		security.AuditLog(r, "parse_error", "failed")
		return
	}

	params, ok := req.Params.(map[string]interface{})
	if !ok {
		s.sendError(w, req.ID, -32602, "Invalid params")
		security.AuditLog(r, "invalid_params", "failed")
		return
	}

	toolName, ok := params["name"].(string)
	if !ok {
		s.sendError(w, req.ID, -32602, "Missing tool name")
		security.AuditLog(r, "missing_tool_name", "failed")
		return
	}

	arguments, ok := params["arguments"].(map[string]interface{})
	if !ok {
		arguments = make(map[string]interface{})
	}

	// Validate path if present
	if path, ok := arguments["path"].(string); ok {
		if err := security.ValidatePath(path); err != nil {
			s.sendError(w, req.ID, -32602, fmt.Sprintf("Invalid path: %s", err.Error()))
			security.AuditLog(r, fmt.Sprintf("invalid_path_%s", toolName), "failed")
			return
		}
	}

	// Sanitize content if present
	if content, ok := arguments["content"].(string); ok {
		arguments["content"] = security.SanitizeContent(content)
	}

	result, err := s.executeTool(toolName, arguments)
	if err != nil {
		s.sendError(w, req.ID, -32603, err.Error())
		security.AuditLog(r, toolName, "failed")
		return
	}

	s.sendResponse(w, req.ID, map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": result,
			},
		},
	})
	security.AuditLog(r, toolName, "success")
}

func (s *MCPServer) executeTool(toolName string, arguments map[string]interface{}) (string, error) {
	switch toolName {
	case "get_note":
		path, ok := arguments["path"].(string)
		if !ok {
			return "", fmt.Errorf("missing or invalid path parameter")
		}
		return s.obsidianAPI.GetNote(path)

	case "create_note":
		path, ok := arguments["path"].(string)
		if !ok {
			return "", fmt.Errorf("missing or invalid path parameter")
		}
		content, ok := arguments["content"].(string)
		if !ok {
			return "", fmt.Errorf("missing or invalid content parameter")
		}
		return s.obsidianAPI.CreateNote(path, content)

	case "update_note":
		path, ok := arguments["path"].(string)
		if !ok {
			return "", fmt.Errorf("missing or invalid path parameter")
		}
		content, ok := arguments["content"].(string)
		if !ok {
			return "", fmt.Errorf("missing or invalid content parameter")
		}
		return s.obsidianAPI.UpdateNote(path, content)

	case "delete_note":
		path, ok := arguments["path"].(string)
		if !ok {
			return "", fmt.Errorf("missing or invalid path parameter")
		}
		return s.obsidianAPI.DeleteNote(path)

	case "list_notes":
		folder := ""
		if f, ok := arguments["folder"].(string); ok {
			folder = f
		}
		return s.obsidianAPI.ListNotes(folder)

	case "search_notes":
		query, ok := arguments["query"].(string)
		if !ok {
			return "", fmt.Errorf("missing or invalid query parameter")
		}
		return s.obsidianAPI.SearchNotes(query)

	case "get_vault_info":
		return s.obsidianAPI.GetVaultInfo()

	default:
		return "", fmt.Errorf("unknown tool: %s", toolName)
	}
}

func (s *MCPServer) sendResponse(w http.ResponseWriter, id interface{}, result interface{}) {
	response := MCPResponse{
		Jsonrpc: "2.0",
		ID:      id,
		Result:  result,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *MCPServer) sendError(w http.ResponseWriter, id interface{}, code int, message string) {
	response := MCPResponse{
		Jsonrpc: "2.0",
		ID:      id,
		Error: &MCPError{
			Code:    code,
			Message: message,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(response)
}

func (s *MCPServer) handleMCPRequest(w http.ResponseWriter, r *http.Request) {
	var req MCPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendError(w, nil, -32700, "Parse error")
		return
	}

	// Handle notifications (requests without ID) - these don't need a response
	if req.ID == nil {
		// Notifications like "notifications/initialized" don't require a response
		w.WriteHeader(http.StatusOK)
		return
	}

	// Route based on method
	switch req.Method {
	case "initialize":
		// Re-encode the request for the handler
		body, _ := json.Marshal(req)
		r.Body = io.NopCloser(bytes.NewBuffer(body))
		s.handleInitialize(w, r)
	case "tools/list":
		body, _ := json.Marshal(req)
		r.Body = io.NopCloser(bytes.NewBuffer(body))
		s.handleListTools(w, r)
	case "tools/call":
		body, _ := json.Marshal(req)
		r.Body = io.NopCloser(bytes.NewBuffer(body))
		s.handleCallTool(w, r)
	default:
		s.sendError(w, req.ID, -32601, fmt.Sprintf("Method not found: %s", req.Method))
	}
}

func (s *MCPServer) setupRoutes() *mux.Router {
	router := mux.NewRouter()

	// Create security middleware wrapper
	securityMw := security.Middleware(s.config.Security, s.rateLimiter)

	// Apply security middleware to all routes
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			securityMw(func(w http.ResponseWriter, r *http.Request) {
				security.ValidateInputMiddleware(next.ServeHTTP)(w, r)
			})(w, r)
		})
	})

	// MCP endpoints - support both root and /mcp/ prefix for compatibility
	router.HandleFunc("/", s.handleMCPRequest).Methods("POST")
	router.HandleFunc("/mcp/initialize", s.handleInitialize).Methods("POST")
	router.HandleFunc("/mcp/tools/list", s.handleListTools).Methods("POST")
	router.HandleFunc("/mcp/tools/call", s.handleCallTool).Methods("POST")

	// Health check (no security for monitoring)
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	}).Methods("GET")

	return router
}

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	obsidianAPI := api.NewObsidianAPI(config.ObsidianAPI.BaseURL, config.ObsidianAPI.Token)

	server := &MCPServer{
		config:      config,
		obsidianAPI: obsidianAPI,
		rateLimiter: security.NewRateLimiter(),
	}

	router := server.setupRoutes()

	addr := fmt.Sprintf("%s:%d", config.MCP.Host, config.MCP.Port)
	log.Printf("Starting MCP server on %s", addr)
	log.Printf("Obsidian API endpoint: %s", config.ObsidianAPI.BaseURL)

	// Log security settings
	if config.Security.EnableAuth {
		log.Printf("🔒 Authentication: ENABLED")
	}
	if len(config.Security.AllowedIPs) > 0 {
		log.Printf("🔒 IP Whitelist: %v", config.Security.AllowedIPs)
	}
	if config.Security.EnableRateLimit {
		log.Printf("🔒 Rate Limit: %d req/min", config.Security.RateLimit)
	}

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
