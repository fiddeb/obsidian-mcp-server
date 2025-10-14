<!-- Workspace-specific instructions for GitHub Copilot -->

## Project: Obsidian MCP Server in Go

A Model Context Protocol (MCP) server written in Go that integrates with Obsidian via the Local REST API plugin.

## Code Structure

The project follows Go best practices with package separation:

```
obsidian-mcp/
├── main.go              # Entry point, MCP server setup and handlers
├── api/
│   └── obsidian.go     # ObsidianAPI client for REST API integration
├── security/
│   └── security.go     # Security middleware, rate limiting, validation
├── config.yaml         # Configuration file
└── go.mod              # Go module definition
```

## Development Guidelines

### General Principles
- Never attempt to compile the project automatically
- User will handle all terminal commands and compilation
- Be concise and minimize use of emojis in responses
- Focus on code quality and Go best practices

### Documentation
- All project documentation must be stored in Obsidian vault
- Documentation path: `Dev/Obsidian MCP server/` in user's vault
- Use MCP server tools to create/update documentation when available
- README.md in project root is the only exception (keep in repo)
- Do not create or maintain a `docs/` directory in the project

### Using MCP Server for Documentation
When the MCP server is running (http://localhost:8080), use it to:
- Create documentation in Obsidian: `create_note` tool
- Update existing docs: `update_note` tool
- Search documentation: `search_notes` tool
- Path for all docs: `Dev/Obsidian MCP server/{filename}.md`

Example MCP tool usage:
```json
{
  "name": "create_note",
  "arguments": {
    "path": "Dev/Obsidian MCP server/Feature Documentation.md",
    "content": "# Feature Name\n\nDescription..."
  }
}
```

### Code Style
- Follow standard Go conventions
- Use meaningful variable and function names
- Packages: lowercase, single word when possible
- Exported functions/types: PascalCase
- Private functions/types: camelCase
- Error handling: always check errors, return descriptive error messages

### Project-Specific Rules
- Main package: Entry point and HTTP handlers only
- API package: All Obsidian REST API interactions
- Security package: Middleware, validation, rate limiting
- Configuration: Load from config.yaml or environment variables
- Never expose API tokens or sensitive data in code

## Available MCP Tools

When server is running, these tools are available:
- `get_note` - Retrieve note content by path
- `create_note` - Create new note with content
- `update_note` - Update existing note
- `delete_note` - Delete a note
- `list_notes` - List all notes (optionally filter by folder)
- `search_notes` - Full-text search across notes
- `get_vault_info` - Get vault statistics and info

## Testing Workflow

User will handle:
- Compilation: `go build -o obsidian-mcp-server`
- Running server: `./obsidian-mcp-server` or VS Code task
- API testing: curl commands or MCP client

Do not attempt to:
- Run terminal commands
- Compile the project
- Execute the binary
- Run tests automatically

## Communication Style

- Be direct and technical
- Minimize emojis (use sparingly, only when truly helpful)
- Provide code examples when explaining concepts
- Focus on Go best practices and idiomatic code
- When suggesting changes, explain the reasoning

## Current Project Status

The MCP server is production-ready with:
- 7 working tools (get_note, create_note, update_note, delete_note, list_notes, search_notes, get_vault_info)
- Security features (rate limiting, IP filtering, input validation)
- VS Code integration via MCP extension
- Comprehensive error handling
- Automatic .md file extension normalization