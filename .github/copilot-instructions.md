<!-- Workspace-specific instructions for GitHub Copilot -->

## Project: Obsidian MCP Server in Go

A Model Context Protocol (MCP) server written in Go that integrates with Obsidian via the Local REST API plugin.

## Code Structure

The project follows Go best practices with package separation:

```
obsidian-mcp/
├── main.go              # Entry point, MCP server setup (stdio) and tool handlers
├── api/
│   └── obsidian.go     # ObsidianAPI client for REST API integration
├── security/
│   └── security.go     # Path validation and content sanitization
├── docs/
│   └── USER_GUIDE.md   # Comprehensive user guide (English)
├── LICENSE              # MIT License
└── go.mod              # Go module definition (includes MCP SDK)
```

## Development Guidelines

### General Principles
- Be concise and minimize use of emojis in responses
- Focus on code quality and Go best practices

### Documentation

**Language:**
- Obsidian documentation: Swedish with technical terms in English
- Repository docs (`README.md`, `docs/`): English
- Code comments: English
- Commit messages: English

**Structure:**
- Comprehensive documentation: Obsidian vault at `Dev/Obsidian MCP server/`
- Quick reference/guides: `docs/` folder in repository (English)
- Main README: Repository root (English)

**Obsidian Documentation (Swedish):**
- Write in Swedish but keep technical terms in English
- Example: "MCP-servern använder stdio transport" (not "standard in/ut")
- Keep code, commands, and API terms in English
- Path: `Dev/Obsidian MCP server/{filename}.md`

**Repository Documentation (English):**
- `docs/USER_GUIDE.md` - Comprehensive user guide
- `README.md` - Project overview and quick start
- Keep concise and focused


### Using MCP Server for Documentation

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
- Main package: MCP server setup (stdio transport) and tool handlers
- API package: All Obsidian REST API interactions
- Security package: Path validation and content sanitization only
- Configuration: Via environment variables (set by MCP client)
- Never expose API tokens or sensitive data in code
- No HTTP server code (uses stdio transport)

## Available MCP Tools

When server is running, these tools are available:
- `get_note` - Retrieve note content by path
- `create_note` - Create new note with content
- `update_note` - Update existing note
- `delete_note` - Delete a note
- `list_notes` - List all notes (optionally filter by folder)
- `search_notes` - Full-text search across notes
- `get_vault_info` - Get vault statistics and info

## Architecture

### Stdio Transport
- Server communicates via stdin/stdout
- MCP clients start server as subprocess
- JSON-RPC messages over stdio
- Follows MCP Protocol 2024-11-05

### Security Model
- No network exposure (stdio only)
- Path validation prevents directory traversal
- Content sanitization removes dangerous characters
- No rate limiting needed (single process per client)

## Testing Workflow

User will handle:
- MCP client starts server automatically
- Manual testing: `./obsidian-mcp-server` with stdin/stdout


## Communication Style

- Be direct and technical
- Minimize emojis (use sparingly, only when truly helpful)
- Provide code examples when explaining concepts
- Focus on Go best practices and idiomatic code
- When suggesting changes, explain the reasoning

## Current Project Status

The MCP server is production-ready with:
- 7 working tools (get_note, create_note, update_note, delete_note, list_notes, search_notes, get_vault_info)
- Stdio transport (official MCP Go SDK v1.0.0)
- Security features (path validation, content sanitization)
- VS Code and Claude Desktop integration
- Comprehensive error handling
- Automatic .md file extension normalization