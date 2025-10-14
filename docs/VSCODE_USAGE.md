# Använda Obsidian MCP Server med VS Code

## 📌 Viktigt att veta

Den nuvarande MCP-servern fungerar redan perfekt via HTTP API och körs på `http://localhost:8080`. 

För VS Code/GitHub Copilot finns det några alternativ:

## Alternativ 1: Använd HTTP API direkt (Rekommenderat för nu)

VS Code och GitHub Copilot kan inte använda MCP-servrar direkt ännu, men du kan:

### 1. Håll servern igång

```bash
# Kör i en terminal
cd /Users/faar/Documents/Src/github/fiddeb/Obsidian_mcp
go run .
```

Eller använd VS Code task (som redan är konfigurerad):
- Tryck `Cmd+Shift+P`
- Välj "Tasks: Run Task"
- Välj "Run Obsidian MCP Server"

### 2. Fråga Copilot att använda HTTP API:et

När du vill arbeta med Obsidian, fråga Copilot:

```
"Kan du använda curl för att skapa en ny anteckning i Obsidian via 
http://localhost:8080/mcp/tools/call? Anteckningen ska heta 'Test.md' 
och innehålla 'Hello World'"
```

Copilot kommer då generera curl-kommandon som detta:

```bash
curl -X POST http://localhost:8080/mcp/tools/call \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "tools/call",
    "params": {
      "name": "create_note",
      "arguments": {
        "path": "Test.md",
        "content": "Hello World"
      }
    }
  }'
```

## Alternativ 2: VS Code MCP Extension (Framtida)

Microsoft arbetar på MCP-stöd i VS Code. När det kommer kan du:

1. Installera MCP-extensionen
2. Konfigurera i `.vscode/settings.json`:

```json
{
  "mcp.servers": {
    "obsidian": {
      "url": "http://localhost:8080",
      "type": "http"
    }
  }
}
```

## Alternativ 3: Claude Desktop (Bäst för MCP)

Om du vill använda fullständig MCP-funktionalitet:

1. Installera Claude Desktop
2. Konfigurera i `~/Library/Application Support/Claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "obsidian": {
      "command": "go",
      "args": ["run", "."],
      "cwd": "/Users/faar/Documents/Src/github/fiddeb/Obsidian_mcp",
      "env": {
        "OBSIDIAN_API_TOKEN": "5060e7700b373ccd30eb08293d51f936808c7eb1dd650f6762d594ae97aafdcc",
        "OBSIDIAN_API_BASE_URL": "http://127.0.0.1:27123"
      }
    }
  }
}
```

3. Starta om Claude Desktop

## 🎯 Praktiska exempel för VS Code

### Skapa en VS Code snippet för vanliga operationer

Skapa `.vscode/obsidian.code-snippets`:

```json
{
  "Create Obsidian Note": {
    "prefix": "obs-create",
    "body": [
      "curl -X POST http://localhost:8080/mcp/tools/call \\\\",
      "  -H \"Content-Type: application/json\" \\\\",
      "  -d '{",
      "    \"jsonrpc\": \"2.0\",",
      "    \"id\": 1,",
      "    \"method\": \"tools/call\",",
      "    \"params\": {",
      "      \"name\": \"create_note\",",
      "      \"arguments\": {",
      "        \"path\": \"${1:filename}.md\",",
      "        \"content\": \"${2:content}\"",
      "      }",
      "    }",
      "  }'"
    ]
  },
  "Search Obsidian": {
    "prefix": "obs-search",
    "body": [
      "curl -X POST http://localhost:8080/mcp/tools/call \\\\",
      "  -H \"Content-Type: application/json\" \\\\",
      "  -d '{",
      "    \"jsonrpc\": \"2.0\",",
      "    \"id\": 1,",
      "    \"params\": {",
      "      \"name\": \"search_notes\",",
      "      \"arguments\": {\"query\": \"${1:search term}\"}",
      "    }",
      "  }'"
    ]
  }
}
```

### Skapa ett shell script

Skapa `obsidian-helper.sh`:

```bash
#!/bin/bash

# Hjälpfunktioner för Obsidian MCP

obsidian_create() {
    curl -X POST http://localhost:8080/mcp/tools/call \
      -H "Content-Type: application/json" \
      -d "{
        \"jsonrpc\": \"2.0\",
        \"id\": 1,
        \"method\": \"tools/call\",
        \"params\": {
          \"name\": \"create_note\",
          \"arguments\": {
            \"path\": \"$1\",
            \"content\": \"$2\"
          }
        }
      }"
}

obsidian_search() {
    curl -X POST http://localhost:8080/mcp/tools/call \
      -H "Content-Type: application/json" \
      -d "{
        \"jsonrpc\": \"2.0\",
        \"id\": 1,
        \"params\": {
          \"name\": \"search_notes\",
          \"arguments\": {\"query\": \"$1\"}
        }
      }" | jq -r '.result.content[0].text'
}

# Exempel:
# obsidian_create "Test.md" "# Hello World"
# obsidian_search "projekt"
```

Gör den körbar:
```bash
chmod +x obsidian-helper.sh
```

## 💡 Tips för VS Code

### 1. Använd VS Code Tasks

Lägg till i `.vscode/tasks.json`:

```json
{
  "version": "2.0.0",
  "tasks": [
    {
      "label": "Obsidian: Create Daily Note",
      "type": "shell",
      "command": "curl -X POST http://localhost:8080/mcp/tools/call -H 'Content-Type: application/json' -d '{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"tools/call\",\"params\":{\"name\":\"create_note\",\"arguments\":{\"path\":\"Daily/$(date +%Y-%m-%d).md\",\"content\":\"# $(date +%Y-%m-%d)\\n\\n\"}}}'"
    },
    {
      "label": "Obsidian: Search",
      "type": "shell",
      "command": "curl -s -X POST http://localhost:8080/mcp/tools/call -H 'Content-Type: application/json' -d '{\"jsonrpc\":\"2.0\",\"id\":1,\"params\":{\"name\":\"search_notes\",\"arguments\":{\"query\":\"${input:searchQuery}\"}}}' | jq -r '.result.content[0].text'",
      "problemMatcher": []
    }
  ],
  "inputs": [
    {
      "id": "searchQuery",
      "type": "promptString",
      "description": "Vad vill du söka efter?"
    }
  ]
}
```

### 2. Keybindings

Lägg till i `keybindings.json`:

```json
[
  {
    "key": "cmd+shift+o cmd+shift+s",
    "command": "workbench.action.tasks.runTask",
    "args": "Obsidian: Search"
  },
  {
    "key": "cmd+shift+o cmd+shift+d",
    "command": "workbench.action.tasks.runTask", 
    "args": "Obsidian: Create Daily Note"
  }
]
```

## 🚀 Snabbstart

1. **Starta servern** (om den inte redan körs):
   ```bash
   go run .
   ```

2. **I VS Code terminal**, kör:
   ```bash
   # Skapa anteckning
   curl -X POST http://localhost:8080/mcp/tools/call \
     -H "Content-Type: application/json" \
     -d '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"create_note","arguments":{"path":"Test från VS Code.md","content":"# Hej från VS Code!"}}}'
   
   # Sök
   curl -X POST http://localhost:8080/mcp/tools/call \
     -H "Content-Type: application/json" \
     -d '{"jsonrpc":"2.0","id":1,"params":{"name":"search_notes","arguments":{"query":"VS Code"}}}'
   ```

3. **Fråga GitHub Copilot**:
   - "Använd Obsidian API på localhost:8080 för att skapa en dagboksanteckning"
   - "Sök i Obsidian efter 'projekt' via MCP API"

## 📚 Tillgängliga verktyg

- `get_note` - Hämta anteckning
- `create_note` - Skapa anteckning
- `update_note` - Uppdatera anteckning  
- `delete_note` - Ta bort anteckning
- `list_notes` - Lista anteckningar
- `search_notes` - Sök med kontext
- `get_vault_info` - Vault-info
- `create_folder` - Skapa mapp

Se `MCP-server/README.md` i din Obsidian vault för fullständig dokumentation!
