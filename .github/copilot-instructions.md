<!-- Use this file to provide workspace-specific custom instructions to Copilot. For more details, visit https://code.visualstudio.com/docs/copilot/copilot-customization#_use-a-githubcopilotinstructionsmd-file -->
- [x] Verify that the copilot-instructions.md file in the .github directory is created. ✅ Created

- [x] Clarify Project Requirements ✅ MCP server in Go for Obsidian REST API integration

- [x] Scaffold the Project ✅ Created main.go, obsidian_api.go, go.mod, config.yaml, README.md, .gitignore

- [x] Customize the Project ✅ Complete MCP server with all Obsidian operations

- [x] Install Required Extensions ✅ No extensions needed

- [x] Compile the Project ✅ Successfully built binary

- [x] Create and Run Task ✅ Created VS Code task and server is running

- [x] Launch the Project ✅ Server is running and responding

- [x] Ensure Documentation is Complete ✅ README and setup instructions complete

## Project: Obsidian MCP Server in Go
This is an MCP (Model Context Protocol) server written in Go that integrates with Obsidian via the Local REST API plugin.

### Setup Instructions
1. Install Obsidian Local REST API plugin
2. Configure API token in config.yaml or environment variables
3. Run: `go run .` or use VS Code task "Run Obsidian MCP Server"
4. Server runs on http://localhost:8080

### Available Tools
- get_note, create_note, update_note, delete_note
- list_notes, search_notes, get_vault_info, create_folder

### Testing
- Health check: `curl http://localhost:8080/health`
- List tools: `curl -X POST http://localhost:8080/mcp/tools/list`