package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"obsidian-mcp/api"
	"obsidian-mcp/security"

	"github.com/modelcontextprotocol/go-sdk/mcp"
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
		Description string `yaml:"description"`
	} `yaml:"mcp"`
}

// contextKey type for context values
type contextKey string

const apiKey contextKey = "api"

// Tool Input/Output types

type GetNoteInput struct {
	Path string `json:"path" jsonschema:"description:Path to the note file"`
}

type CreateNoteInput struct {
	Path    string `json:"path" jsonschema:"description:Path where the note should be created"`
	Content string `json:"content" jsonschema:"description:Content of the note"`
}

type UpdateNoteInput struct {
	Path    string `json:"path" jsonschema:"description:Path to the note to update"`
	Content string `json:"content" jsonschema:"description:New content for the note"`
}

type DeleteNoteInput struct {
	Path string `json:"path" jsonschema:"description:Path to the note to delete"`
}

type ListNotesInput struct {
	Folder string `json:"folder,omitempty" jsonschema:"description:Optional folder to filter by"`
}

type SearchNotesInput struct {
	Query string `json:"query" jsonschema:"description:Search query"`
}

type VaultInfoInput struct {
	// No parameters needed
}

// Tool Output types

type NoteContentOutput struct {
	Content string `json:"content" jsonschema:"description:Content of the note"`
}

type MessageOutput struct {
	Message string `json:"message" jsonschema:"description:Operation result message"`
}

type NotesListOutput struct {
	Notes string `json:"notes" jsonschema:"description:List of notes"`
}

type SearchResultOutput struct {
	Results string `json:"results" jsonschema:"description:Search results"`
}

type VaultInfoOutput struct {
	Info string `json:"info" jsonschema:"description:Vault information"`
}

// Tool handlers

func GetNote(ctx context.Context, req *mcp.CallToolRequest, input GetNoteInput) (*mcp.CallToolResult, NoteContentOutput, error) {
	obsidianAPI := ctx.Value(apiKey).(*api.ObsidianAPI)

	if err := security.ValidatePath(input.Path); err != nil {
		return nil, NoteContentOutput{}, fmt.Errorf("invalid path: %v", err)
	}

	content, err := obsidianAPI.GetNote(input.Path)
	if err != nil {
		return nil, NoteContentOutput{}, fmt.Errorf("failed to get note: %v", err)
	}

	return nil, NoteContentOutput{Content: content}, nil
}

func CreateNote(ctx context.Context, req *mcp.CallToolRequest, input CreateNoteInput) (*mcp.CallToolResult, MessageOutput, error) {
	obsidianAPI := ctx.Value(apiKey).(*api.ObsidianAPI)

	if err := security.ValidatePath(input.Path); err != nil {
		return nil, MessageOutput{}, fmt.Errorf("invalid path: %v", err)
	}

	sanitizedContent := security.SanitizeContent(input.Content)
	msg, err := obsidianAPI.CreateNote(input.Path, sanitizedContent)
	if err != nil {
		return nil, MessageOutput{}, fmt.Errorf("failed to create note: %v", err)
	}

	return nil, MessageOutput{Message: msg}, nil
}

func UpdateNote(ctx context.Context, req *mcp.CallToolRequest, input UpdateNoteInput) (*mcp.CallToolResult, MessageOutput, error) {
	obsidianAPI := ctx.Value(apiKey).(*api.ObsidianAPI)

	if err := security.ValidatePath(input.Path); err != nil {
		return nil, MessageOutput{}, fmt.Errorf("invalid path: %v", err)
	}

	sanitizedContent := security.SanitizeContent(input.Content)
	msg, err := obsidianAPI.UpdateNote(input.Path, sanitizedContent)
	if err != nil {
		return nil, MessageOutput{}, fmt.Errorf("failed to update note: %v", err)
	}

	return nil, MessageOutput{Message: msg}, nil
}

func DeleteNote(ctx context.Context, req *mcp.CallToolRequest, input DeleteNoteInput) (*mcp.CallToolResult, MessageOutput, error) {
	obsidianAPI := ctx.Value(apiKey).(*api.ObsidianAPI)

	if err := security.ValidatePath(input.Path); err != nil {
		return nil, MessageOutput{}, fmt.Errorf("invalid path: %v", err)
	}

	msg, err := obsidianAPI.DeleteNote(input.Path)
	if err != nil {
		return nil, MessageOutput{}, fmt.Errorf("failed to delete note: %v", err)
	}

	return nil, MessageOutput{Message: msg}, nil
}

func ListNotes(ctx context.Context, req *mcp.CallToolRequest, input ListNotesInput) (*mcp.CallToolResult, NotesListOutput, error) {
	obsidianAPI := ctx.Value(apiKey).(*api.ObsidianAPI)

	if input.Folder != "" {
		if err := security.ValidatePath(input.Folder); err != nil {
			return nil, NotesListOutput{}, fmt.Errorf("invalid folder path: %v", err)
		}
	}

	result, err := obsidianAPI.ListNotes(input.Folder)
	if err != nil {
		return nil, NotesListOutput{}, fmt.Errorf("failed to list notes: %v", err)
	}

	return nil, NotesListOutput{Notes: result}, nil
}

func SearchNotes(ctx context.Context, req *mcp.CallToolRequest, input SearchNotesInput) (*mcp.CallToolResult, SearchResultOutput, error) {
	obsidianAPI := ctx.Value(apiKey).(*api.ObsidianAPI)

	result, err := obsidianAPI.SearchNotes(input.Query)
	if err != nil {
		return nil, SearchResultOutput{}, fmt.Errorf("failed to search notes: %v", err)
	}

	return nil, SearchResultOutput{Results: result}, nil
}

func GetVaultInfo(ctx context.Context, req *mcp.CallToolRequest, input VaultInfoInput) (*mcp.CallToolResult, VaultInfoOutput, error) {
	obsidianAPI := ctx.Value(apiKey).(*api.ObsidianAPI)

	result, err := obsidianAPI.GetVaultInfo()
	if err != nil {
		return nil, VaultInfoOutput{}, fmt.Errorf("failed to get vault info: %v", err)
	}

	return nil, VaultInfoOutput{Info: result}, nil
}

func loadConfig() (Config, error) {
	var config Config

	// Default configuration
	config.ObsidianAPI.BaseURL = "http://localhost:27123"
	config.ObsidianAPI.Port = 27123
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

func main() {
	// Load configuration
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create Obsidian API client
	obsidianAPI := api.NewObsidianAPI(config.ObsidianAPI.BaseURL, config.ObsidianAPI.Token)

	// Create context with API client
	ctx := context.WithValue(context.Background(), apiKey, obsidianAPI)

	// Create MCP server
	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "obsidian-mcp-server",
			Version: "1.0.2",
		},
		nil,
	)

	// Register all tools
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_note",
		Description: "Get the content of a note by its path",
	}, GetNote)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_note",
		Description: "Create a new note with the specified path and content",
	}, CreateNote)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "update_note",
		Description: "Update an existing note with new content",
	}, UpdateNote)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete_note",
		Description: "Delete a note by its path",
	}, DeleteNote)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_notes",
		Description: "List all notes in the vault or in a specific folder",
	}, ListNotes)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_notes",
		Description: "Search for notes containing the specified query",
	}, SearchNotes)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_vault_info",
		Description: "Get information about the vault (authentication status, version, statistics)",
	}, GetVaultInfo)

	// Run server over stdio
	log.Println("Starting Obsidian MCP Server with stdio transport...")
	if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
