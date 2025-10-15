# Obsidian MCP Server - User Guide

A comprehensive guide to installing, configuring, and using the Obsidian MCP Server.

## Table of Contents

- [What is the Obsidian MCP Server?](#what-is-the-obsidian-mcp-server)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
  - [Getting Your API Token](#getting-your-api-token)
  - [VS Code Setup](#vs-code-setup)
  - [Claude Desktop Setup](#claude-desktop-setup)
- [Available Tools](#available-tools)
- [Usage Examples](#usage-examples)
  - [Using with VS Code GitHub Copilot](#using-with-vs-code-github-copilot)
  - [Using with Claude Desktop](#using-with-claude-desktop)
- [Advanced Usage](#advanced-usage)
- [Security Features](#security-features)
- [Troubleshooting](#troubleshooting)
- [Best Practices](#best-practices)

---

## What is the Obsidian MCP Server?

The Obsidian MCP Server is a **Model Context Protocol (MCP)** server that integrates your Obsidian vault with AI assistants like GitHub Copilot and Claude Desktop. It allows AI to:

- Read notes from your vault
- Create new notes
- Update existing notes
- Search through your knowledge base
- Delete notes
- Get vault statistics

**Key Features:**

- ✅ **Secure**: Process-based communication via stdio (no network exposure)
- ✅ **Fast**: Written in Go for high performance
- ✅ **Standard**: Uses the official MCP Go SDK
- ✅ **Smart**: Automatic `.md` extension handling
- ✅ **Safe**: Built-in path validation to prevent directory traversal

**Architecture Overview:**

```
AI Assistant (VS Code/Claude)
         ↓
   MCP Client
         ↓
  Obsidian MCP Server (stdio)
         ↓
  Obsidian Local REST API
         ↓
   Your Vault
```

---

## Prerequisites

Before you can use the Obsidian MCP Server, you need:

### 1. Obsidian with Local REST API Plugin

**Install the plugin:**

1. Open Obsidian
2. Go to **Settings** → **Community Plugins**
3. Click **Browse** and search for "Local REST API"
4. Install and **Enable** the plugin
5. The plugin runs on port **27123** by default

### 2. Go 1.21 or Later

**Check if Go is installed:**

```bash
go version
```

**If not installed, download from:** https://go.dev/download/

### 3. MCP Client

Choose one or both:

- **VS Code** with the MCP extension (for GitHub Copilot integration)
- **Claude Desktop** (for direct Claude integration)

---

## Installation

### Step 1: Clone or Download the Project

```bash
git clone https://github.com/fiddeb/obsidian-mcp-server.git
cd obsidian-mcp-server
```

### Step 2: Install Dependencies

```bash
go mod tidy
```

This downloads all required Go packages, including the official MCP SDK.

### Step 3: Build the Server

```bash
go build -o obsidian-mcp-server
```

This creates an executable binary called `obsidian-mcp-server` (or `obsidian-mcp-server.exe` on Windows).

**Verify the build:**

```bash
# macOS/Linux
ls -la obsidian-mcp-server

# Windows
dir obsidian-mcp-server.exe
```

---

## Configuration

The server is configured differently depending on which MCP client you're using.

### Getting Your API Token

All configurations require your Obsidian Local REST API token:

1. Open **Obsidian**
2. Go to **Settings** → **Community Plugins** → **Local REST API**
3. Copy the **API Key** shown in the plugin settings
4. Keep this token secure - treat it like a password!

### VS Code Setup

#### Step 1: Install the MCP Extension

Search for "MCP" in VS Code extensions and install the official MCP extension.

#### Step 2: Create MCP Configuration

Create a file at `.vscode/mcp.json` in your workspace:

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

**Important:**

- Replace `/absolute/path/to/obsidian-mcp-server` with the actual path
  - Example (macOS): `/Users/yourname/projects/obsidian-mcp-server/obsidian-mcp-server`
  - Example (Windows): `C:\\Users\\yourname\\projects\\obsidian-mcp-server\\obsidian-mcp-server.exe`
- Replace `your-api-token-here` with the token from Obsidian
- **Use absolute paths** - relative paths won't work

#### Step 3: Add to .gitignore

**Important:** Don't commit your token to version control!

Add to `.gitignore`:

```
.vscode/mcp.json
```

#### Step 4: Reload VS Code

Press `Cmd+Shift+P` (macOS) or `Ctrl+Shift+P` (Windows/Linux), then select:

```
Developer: Reload Window
```

The MCP extension will automatically start the server when needed.

### Claude Desktop Setup

#### Step 1: Locate Claude Desktop Config

**macOS:**
```
~/Library/Application Support/Claude/claude_desktop_config.json
```

**Windows:**
```
%APPDATA%\Claude\claude_desktop_config.json
```

**Linux:**
```
~/.config/Claude/claude_desktop_config.json
```

#### Step 2: Edit Configuration

Open `claude_desktop_config.json` and add:

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

**Important:**
- Use absolute path to the binary
- Replace with your actual API token
- On Windows, use double backslashes in paths: `C:\\Users\\...`

#### Step 3: Restart Claude Desktop

Completely quit and restart Claude Desktop for changes to take effect.

---

## Available Tools

The server provides 7 tools for interacting with your Obsidian vault:

### 1. `get_note`

**Description:** Get the content of a note

**Parameters:**
- `path` (string): Path to the note (e.g., `"Daily/2025-10-15.md"` or `"Daily/2025-10-15"`)

**Returns:** The note's content as markdown text

**Example:**
```json
{
  "path": "Projects/MCP Integration.md"
}
```

### 2. `create_note`

**Description:** Create a new note

**Parameters:**
- `path` (string): Where to create the note
- `content` (string): The note's content

**Returns:** Confirmation message

**Example:**
```json
{
  "path": "Meeting Notes/Team Sync.md",
  "content": "# Team Sync\n\n## Attendees\n- Alice\n- Bob\n\n## Topics\n- Q1 Planning"
}
```

### 3. `update_note`

**Description:** Update an existing note

**Parameters:**
- `path` (string): Path to the note
- `content` (string): New content (replaces existing)

**Returns:** Confirmation message

**Example:**
```json
{
  "path": "TODO.md",
  "content": "# TODO\n\n- [x] Set up MCP server\n- [ ] Write documentation"
}
```

### 4. `delete_note`

**Description:** Delete a note

**Parameters:**
- `path` (string): Path to the note to delete

**Returns:** Confirmation message

**Example:**
```json
{
  "path": "Old Draft.md"
}
```

### 5. `list_notes`

**Description:** List all notes in the vault or a specific folder

**Parameters:**
- `folder` (string, optional): Folder to filter by

**Returns:** List of note paths

**Example (all notes):**
```json
{}
```

**Example (specific folder):**
```json
{
  "folder": "Projects"
}
```

### 6. `search_notes`

**Description:** Search for notes containing specific text

**Parameters:**
- `query` (string): Search query

**Returns:** Search results with filename and context

**Example:**
```json
{
  "query": "MCP integration"
}
```

### 7. `get_vault_info`

**Description:** Get vault statistics and information

**Parameters:** None

**Returns:** Vault info including:
- Authentication status
- Obsidian version
- Number of notes and folders

**Example:**
```json
{}
```

---

## Usage Examples

### Using with VS Code GitHub Copilot

Once configured, you can interact with your Obsidian vault directly through GitHub Copilot Chat.

#### Example 1: Create a Daily Note

**You ask:**
```
@workspace Create a daily note for today in my Obsidian vault
```

**Copilot will:**
1. Use `create_note` to create `Daily/2025-10-15.md`
2. Add standard daily note template
3. Confirm creation

#### Example 2: Search Your Notes

**You ask:**
```
@workspace Search my Obsidian vault for notes about "machine learning"
```

**Copilot will:**
1. Use `search_notes` with query "machine learning"
2. Return matching notes with context
3. Optionally summarize findings

#### Example 3: Update a Note

**You ask:**
```
@workspace Add this meeting summary to my Meeting Notes/Team Sync.md:
- Discussed Q1 goals
- Assigned tasks
- Next meeting: 2025-10-22
```

**Copilot will:**
1. Use `get_note` to read current content
2. Use `update_note` to append the new content
3. Confirm update

#### Example 4: Get Vault Statistics

**You ask:**
```
@workspace How many notes do I have in my Obsidian vault?
```

**Copilot will:**
1. Use `get_vault_info`
2. Report the number of notes and folders

#### Example 5: List Notes in a Folder

**You ask:**
```
@workspace List all my project notes
```

**Copilot will:**
1. Use `list_notes` with folder "Projects"
2. Display the list

### Using with Claude Desktop

Claude Desktop integration works similarly, but you interact through Claude's chat interface.

#### Example 1: Research Assistant

**You:**
```
Can you search my Obsidian vault for notes about "Python async programming" and summarize what I've learned?
```

**Claude will:**
1. Use `search_notes` to find relevant notes
2. Use `get_note` to read the full content
3. Provide a summary of your learnings

#### Example 2: Note Organization

**You:**
```
I want to create a new note called "Python Best Practices" with sections for:
- Code Style
- Error Handling
- Testing
- Performance

Can you create this in my Projects folder?
```

**Claude will:**
1. Use `create_note` with the structured template
2. Confirm creation

#### Example 3: Knowledge Management

**You:**
```
Search my vault for all notes about "MCP" and create an index note that links to them
```

**Claude will:**
1. Use `search_notes` to find MCP-related notes
2. Use `create_note` to create an index with wikilinks
3. Organize by topic/date

---

## Advanced Usage

### Custom VS Code Tasks

You can create VS Code tasks for common operations. Add to `.vscode/tasks.json`:

```json
{
  "version": "2.0.0",
  "tasks": [
    {
      "label": "Build Obsidian MCP Server",
      "type": "shell",
      "command": "go",
      "args": ["build", "-o", "obsidian-mcp-server"],
      "group": {
        "kind": "build",
        "isDefault": true
      },
      "problemMatcher": ["$go"]
    },
    {
      "label": "Test MCP Server Manually",
      "type": "shell",
      "command": "${workspaceFolder}/obsidian-mcp-server",
      "isBackground": true,
      "options": {
        "env": {
          "OBSIDIAN_API_TOKEN": "${input:apiToken}",
          "OBSIDIAN_API_BASE_URL": "http://localhost:27123"
        }
      },
      "problemMatcher": []
    }
  ],
  "inputs": [
    {
      "id": "apiToken",
      "type": "promptString",
      "description": "Enter your Obsidian API token",
      "password": true
    }
  ]
}
```

### Manual Testing via stdin/stdout

For development and debugging, you can test the server manually:

```bash
# Set environment variables
export OBSIDIAN_API_TOKEN="your-token"
export OBSIDIAN_API_BASE_URL="http://localhost:27123"

# Start the server
./obsidian-mcp-server

# In another terminal, send JSON-RPC messages:
echo '{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": {
    "protocolVersion": "2024-11-05",
    "capabilities": {},
    "clientInfo": {
      "name": "test",
      "version": "1.0"
    }
  }
}' | ./obsidian-mcp-server
```

### Building for Different Platforms

Build binaries for different operating systems:

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o obsidian-mcp-server-linux

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o obsidian-mcp-server-mac-intel

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o obsidian-mcp-server-mac-arm

# Windows
GOOS=windows GOARCH=amd64 go build -o obsidian-mcp-server.exe
```

---

## Security Features

The Obsidian MCP Server includes several security features to protect your vault:

### 1. Path Validation

**Blocks dangerous paths:**
- Absolute paths (e.g., `/etc/passwd`)
- Directory traversal (e.g., `../../sensitive/file`)
- Dangerous characters (`` ` ``, `$`, `|`, `;`, `~`)

**Example blocked paths:**
```
/etc/passwd             ❌ Absolute path
../../../secret.md      ❌ Directory traversal
file`command`.md        ❌ Command injection attempt
$(malicious).md         ❌ Variable expansion attempt
```

**Example allowed paths:**
```
Daily/2025-10-15.md     ✅ Normal note
Projects/MCP/notes.md   ✅ Nested folder
Meeting Notes.md        ✅ Root level note
```

### 2. Content Sanitization

- Removes null bytes from content
- Prevents binary data injection
- Ensures clean text storage

### 3. No Network Exposure

- Uses stdio (stdin/stdout) for communication
- No HTTP server running
- No open network ports
- Process-based isolation

### 4. Token Security

**Best practices:**

✅ **DO:**
- Store tokens in MCP configuration files
- Add configuration files to `.gitignore`
- Use environment variables
- Rotate tokens periodically

❌ **DON'T:**
- Hardcode tokens in source code
- Commit tokens to version control
- Share tokens in chat/email
- Use the same token across multiple services

### 5. API Permissions

The Obsidian Local REST API token has full access to your vault. Consider:

- Using a dedicated vault for MCP integration (if you have sensitive data)
- Regularly reviewing which notes are being accessed
- Revoking and regenerating tokens if compromised

---

## Troubleshooting

### Common Issues and Solutions

#### 1. "connection refused" or "dial tcp [::1]:27123: connect: connection refused"

**Problem:** The server can't connect to Obsidian's Local REST API.

**Solutions:**
1. **Check Obsidian is running**
   ```bash
   # Test the API directly
   curl -H "Authorization: Bearer YOUR_TOKEN" http://localhost:27123/
   ```

2. **Verify Local REST API plugin is enabled**
   - Obsidian → Settings → Community Plugins → Local REST API
   - Should show "Enabled"

3. **Check the port number**
   - Default is 27123
   - Verify in Local REST API plugin settings
   - Update `OBSIDIAN_API_BASE_URL` if different

#### 2. "unauthorized" or 401 errors

**Problem:** Invalid or incorrect API token.

**Solutions:**
1. **Get fresh token from Obsidian**
   - Settings → Community Plugins → Local REST API
   - Copy the API key

2. **Update configuration**
   - VS Code: `.vscode/mcp.json`
   - Claude Desktop: `claude_desktop_config.json`
   - Check for extra spaces or quotes

3. **Restart the MCP client**
   - VS Code: Reload window
   - Claude Desktop: Quit and restart

#### 3. Server doesn't start in VS Code

**Problem:** MCP extension can't start the server.

**Solutions:**
1. **Verify the binary path**
   ```bash
   # Check if file exists
   ls -la /absolute/path/to/obsidian-mcp-server
   ```

2. **Rebuild the server**
   ```bash
   cd /path/to/obsidian-mcp-server
   go build -o obsidian-mcp-server
   ```

3. **Check permissions**
   ```bash
   # Make it executable (macOS/Linux)
   chmod +x obsidian-mcp-server
   ```

4. **Check MCP extension logs**
   - VS Code → Output panel
   - Select "MCP: Obsidian" from dropdown

5. **Reload VS Code**
   - Cmd+Shift+P → "Developer: Reload Window"

#### 4. Path issues (file not found)

**Problem:** Notes can't be found or created.

**Solutions:**
1. **Use vault-relative paths**
   ```
   ✅ Projects/Note.md
   ✅ Daily/2025-10-15.md
   ❌ /Users/name/Vault/Projects/Note.md
   ```

2. **Extension is optional**
   ```
   ✅ Projects/Note
   ✅ Projects/Note.md
   (Both work - .md is added automatically)
   ```

3. **Use forward slashes**
   ```
   ✅ Projects/Subfolder/Note.md
   ❌ Projects\Subfolder\Note.md (even on Windows)
   ```

#### 5. Duplicate files created

**Problem:** Both "note" and "note.md" exist.

**Solution:** This was a bug in earlier versions. Update to the latest version:

```bash
git pull
go build -o obsidian-mcp-server
```

Then restart your MCP client.

#### 6. Changes not visible in Obsidian

**Problem:** Files created/updated but don't show in Obsidian.

**Solutions:**
1. **Refresh Obsidian**
   - Click on another note and back
   - Or restart Obsidian

2. **Check vault path**
   - Ensure Obsidian is using the correct vault
   - Settings → About → Vault path

3. **Verify file was created**
   ```bash
   # Check in file system
   ls -la /path/to/vault/Projects/
   ```

---

## Best Practices

### 1. Organize Your Notes

Use consistent folder structures for better AI interaction:

```
vault/
├── Daily/              # Daily notes
├── Projects/           # Project documentation
├── Meetings/           # Meeting notes
├── Archive/            # Old/completed notes
└── Templates/          # Note templates
```

### 2. Use Descriptive Filenames

**Good:**
```
Projects/MCP Integration Guide.md
Meetings/Team Sync 2025-10-15.md
Daily/2025-10-15.md
```

**Avoid:**
```
note1.md
temp.md
untitled.md
```

### 3. Leverage Templates

Create template notes that AI can use:

**Template note:** `Templates/Daily Note.md`
```markdown
# {{date}}

## Tasks
- [ ] 

## Notes


## Reflections

```

**Ask AI:**
```
Create today's daily note using the Daily Note template
```

### 4. Regular Backups

The server has full write access to your vault. Protect against accidents:

- Use Obsidian Sync or Git for version control
- Regular backups to external storage
- Test restore procedures

### 5. Explicit Instructions to AI

Be specific in your requests:

**Good:**
```
Create a note in Projects/Web App called "Architecture Overview" 
with sections for Frontend, Backend, and Database
```

**Less effective:**
```
Make a note about the web app
```

### 6. Verify Important Changes

For critical notes, verify AI changes:

```
Update my TODO list and show me what you changed
```

### 7. Use Wikilinks for Connections

Encourage AI to create connections:

```
Create a note about Python decorators and link it to 
my existing Python notes
```

AI can create wikilinks like `[[Python Basics]]` to connect notes.

---

## Workflow Examples

### Daily Note Automation

**Morning routine:**
```
Create today's daily note with:
- Weather for [location]
- My top 3 priorities
- Upcoming meetings from my calendar
```

### Meeting Notes

**Before meeting:**
```
Create a meeting note for "Q1 Planning" with sections for:
- Attendees
- Agenda
- Discussion
- Action Items
```

**After meeting:**
```
Update the Q1 Planning meeting note with:
[paste your raw notes]
Format it nicely with proper headings and action items
```

### Research and Learning

**Start research:**
```
Search my vault for everything about "machine learning" 
and create a learning roadmap based on what I already know
```

**After learning:**
```
Create a note about [topic] and link it to related notes 
I already have
```

### Project Documentation

**New project:**
```
Create a project folder "My New App" with notes for:
- README
- Architecture
- Tasks
- Meeting Notes
Link them together with wikilinks
```

### Knowledge Management

**Weekly review:**
```
Search for notes created this week and create a 
weekly summary note with links to important content
```

---

## Additional Resources

### Documentation

- [MCP Specification](https://modelcontextprotocol.io/)
- [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk)
- [Obsidian Local REST API](https://github.com/coddingtonbear/obsidian-local-rest-api)

### Support

- **Issues:** [GitHub Issues](https://github.com/fiddeb/obsidian-mcp-server/issues)
- **Discussions:** [GitHub Discussions](https://github.com/fiddeb/obsidian-mcp-server/discussions)

### Development

This project was developed with assistance from:
- **GitHub Copilot** - AI-assisted code generation
- **Claude Sonnet 4.5** - Architecture and documentation

### Contributing

Contributions are welcome! Please feel free to submit:
- Bug reports
- Feature requests
- Pull requests
- Documentation improvements

---

## FAQ

**Q: Does this work offline?**
A: The MCP server works offline (no internet needed), but requires Obsidian to be running locally.

**Q: Can multiple AI assistants use the same server?**
A: Yes, but only one at a time. Each MCP client starts its own instance of the server.

**Q: Is my data sent to the cloud?**
A: No. The MCP server communicates locally via stdio. However, AI assistants (like Claude or GitHub Copilot) may send note content to their servers for processing.

**Q: Can I use this with my own AI models?**
A: Yes! Any MCP-compatible client can use this server. You could integrate it with local LLMs or custom AI tools.

**Q: What happens if Obsidian isn't running?**
A: The server will fail to connect to the Local REST API and return connection errors. Obsidian must be running for the server to work.

**Q: Can I run multiple servers for different vaults?**
A: Yes. Create separate configurations with different command paths and tokens, one for each vault.

**Q: How do I update the server?**
A: Pull the latest code and rebuild:
```bash
cd obsidian-mcp-server
git pull
go build -o obsidian-mcp-server
# Restart your MCP client
```

---

**Version:** 2.0.0  
**Last Updated:** October 2025  
**Status:** Production Ready

Happy note-taking with AI assistance! 🚀
