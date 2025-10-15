# Obsidian MCP Server

A Model Context Protocol (MCP) server written in Go that integrates with Obsidian via the Local REST API plugin. This server uses the official MCP Go SDK and communicates via standard input/output (stdio) following the MCP protocol specification.

> **AI-Developed**: This project was developed with the assistance of GitHub Copilot and Claude Sonnet 4.5 (via VS Code GitHub Copilot Chat).

## Features

- **Note Management**: Create, read, update, and delete notes
- **Search**: Search through your vault for specific text
- **Vault Information**: Get an overview of your vault and its contents
- **Official MCP SDK**: Built with the official [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk)
- **Stdio Transport**: Process-based communication via stdin/stdout
- **MCP Compatible**: Fully compatible with Model Context Protocol 2024-11-05
- **Path Normalization**: Automatic `.md` extension handling for all operations
- **VS Code Integration**: Works seamlessly with VS Code's MCP extension

## Prerequisites

1. **Obsidian** with the **Local REST API** plugin installed and enabled
2. **Go 1.21** or later
3. Configured API token for Local REST API

## Installation

1. Clone or download this project
2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Build the server:
   ```bash
   go build -o obsidian-mcp-server
   ```

## Configuration

The server is configured via environment variables when started by an MCP client (like VS Code or Claude Desktop).

### Required Environment Variables
- `OBSIDIAN_API_TOKEN`: Your Obsidian Local REST API token
- `OBSIDIAN_API_BASE_URL`: Base URL for Obsidian API (default: `http://localhost:27123`)

### Getting Your API Token

1. Open Obsidian
2. Go to Settings → Community Plugins
3. Find "Local REST API" plugin
4. Copy the API key from the plugin settings

## Usage

The MCP server is designed to be started by an MCP client (like VS Code or Claude Desktop) and communicates via stdin/stdout. You typically don't run it manually.

### VS Code Integration

1. **Install the MCP extension** in VS Code

2. **Configure the server** in `.vscode/mcp.json`:
```json
{
  "servers": {
    "obsidian": {
      "command": "/absolute/path/to/obsidian-mcp-server",
      "env": {
        "OBSIDIAN_API_TOKEN": "your-api-token-here",
        "OBSIDIAN_API_BASE_URL": "http://localhost:27123"
      }
    }
  }
}
```

3. **Reload VS Code** - The MCP extension will automatically start the server

4. **Use the tools** via GitHub Copilot Chat or MCP commands

### Claude Desktop Integration

1. **Build the server** (if not already done):
```bash
go build -o obsidian-mcp-server
```

2. **Configure Claude Desktop** by editing `~/Library/Application Support/Claude/claude_desktop_config.json`:
```json
{
  "mcpServers": {
    "obsidian": {
      "command": "/absolute/path/to/obsidian-mcp-server",
      "env": {
        "OBSIDIAN_API_TOKEN": "your-api-token-here",
        "OBSIDIAN_API_BASE_URL": "http://localhost:27123"
      }
    }
  }
}
```

3. **Restart Claude Desktop** - The server will start automatically when needed

### Manual Testing (Advanced)

For development/testing, you can run the server manually and interact via stdin/stdout:

```bash
# Set environment variables
export OBSIDIAN_API_TOKEN="your-token"
export OBSIDIAN_API_BASE_URL="http://localhost:27123"

# Run the server
./obsidian-mcp-server
```

Then send JSON-RPC messages via stdin (see MCP Protocol Examples below).

### Available Tools

1. **get_note** - Get the content of a note
   - Parameter: `path` (path to the note, `.md` extension optional)

2. **create_note** - Create a new note
   - Parameters: `path` (path, `.md` extension optional), `content` (note content)

3. **update_note** - Update an existing note
   - Parameters: `path` (path, `.md` extension optional), `content` (new content)

4. **delete_note** - Delete a note
   - Parameter: `path` (path to the note, `.md` extension optional)

5. **list_notes** - List all notes in the vault
   - Parameter: `folder` (optional, filter by folder path)

6. **search_notes** - Search for notes containing text
   - Parameter: `query` (search query string)

7. **get_vault_info** - Get vault statistics and information
   - No parameters

**Note:** All tools automatically normalize paths by adding the `.md` extension if not present. You can use `"testfile"` or `"testfile.md"` - both work identically.

## MCP Protocol Examples

These examples show the JSON-RPC messages for manual testing via stdin/stdout.

### Initialize
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": {
    "protocolVersion": "2024-11-05",
    "capabilities": {},
    "clientInfo": {
      "name": "example-client",
      "version": "1.0.0"
    }
  }
}
```

### List Tools
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/list"
}
```

### Create a Note
```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "tools/call",
  "params": {
    "name": "create_note",
    "arguments": {
      "path": "My New Note.md",
      "content": "# My New Note\n\nThis is the content of my note."
    }
  }
}
```

## Development

### Project Structure
```
obsidian-mcp/
├── main.go              # MCP server with stdio transport and tool handlers
├── api/
│   └── obsidian.go     # Obsidian REST API client
├── security/
│   └── security.go     # Path validation and content sanitization
├── go.mod               # Go module dependencies (includes MCP SDK)
├── go.sum               # Go module checksums
├── .vscode/
│   └── mcp.json         # MCP server configuration for VS Code
├── .gitignore           # Git ignore file
└── README.md            # This file
```

### Architecture

The server follows a clean architecture:

- **main.go**: MCP server setup using official SDK
  - Tool registration with input/output schemas
  - StdioTransport for stdin/stdout communication
  - Context-based API client injection
  
- **api/obsidian.go**: Obsidian REST API integration
  - HTTP client for Local REST API
  - Path normalization (automatic `.md` extension)
  - CRUD operations for notes
  
- **security/security.go**: Security utilities
  - Path validation (prevent directory traversal)
  - Content sanitization

### Building for Production
```bash
go build -o obsidian-mcp-server .
```

The resulting binary is self-contained and can be deployed anywhere. Just ensure the environment variables are set when the MCP client starts it.

## Troubleshooting

### Common Issues

1. **"connection refused" or "dial tcp [::1]:27123: connect: connection refused"**: 
   - Check that Obsidian is running
   - Verify the Local REST API plugin is enabled in Obsidian
   - Confirm the port (default 27123) matches your Local REST API settings

2. **"unauthorized"**: 
   - Verify your `OBSIDIAN_API_TOKEN` matches the token in Obsidian Local REST API settings
   - Check that the environment variable is correctly set in your MCP client configuration

3. **Server doesn't start in VS Code**:
   - Check the MCP extension output panel for errors
   - Verify the `command` path in `.vscode/mcp.json` is absolute and correct
   - Ensure the binary is built (`go build -o obsidian-mcp-server`)
   - Reload VS Code after configuration changes

4. **Duplicate files (e.g., "testfile" and "testfile.md")**:
   - This has been fixed in the latest version
   - All operations now normalize paths consistently
   - Rebuild the server: `go build -o obsidian-mcp-server`
   - Restart your MCP client

5. **Path issues**:
   - Use forward slashes `/` even on Windows
   - Paths are relative to vault root
   - `.md` extension is optional - it's added automatically

### Debugging

Check the MCP client's output panel/logs:
- **VS Code**: Open Output panel → Select "MCP: Obsidian" from dropdown
- **Claude Desktop**: Check Console logs in Developer Tools

The server logs all activity to stderr, which the MCP client captures.

## Contributing

Contributions are welcome! Feel free to create issues or pull requests.

## Development

This project was developed with:
- **GitHub Copilot** - AI-assisted code generation and editing
- **Claude Sonnet 4.5** - Architecture, problem-solving, and documentation (via GitHub Copilot Chat)

### Technology Stack

- **Go 1.21+** - Modern, efficient language with great concurrency support
- **MCP Go SDK** - Official SDK for Model Context Protocol
- **Obsidian Local REST API** - Integration with Obsidian vault
- **StdioTransport** - Standard MCP communication pattern

### Key Design Decisions

1. **Stdio over HTTP**: Following MCP best practices, the server uses stdin/stdout for communication rather than HTTP
2. **Official SDK**: Using the official MCP Go SDK ensures protocol compliance and reduces maintenance
3. **Path Normalization**: Automatic `.md` extension handling for better UX
4. **Context-Based Injection**: API client passed via context for clean handler signatures

The project follows Go best practices with package separation:
- `main` - MCP server setup and tool handlers
- `api/` - Obsidian REST API client
- `security/` - Validation and sanitization utilities

## License

MIT License