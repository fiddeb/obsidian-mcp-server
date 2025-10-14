# Konfigurera Obsidian MCP Server för Claude Desktop

## Steg 1: Bygg MCP-servern

Först behöver vi bygga en binär av servern:

```bash
cd /Users/faar/Documents/Src/github/fiddeb/Obsidian_mcp
go build -o obsidian-mcp-server .
```

## Steg 2: Hitta Claude Desktop config

Claude Desktop's MCP-konfiguration finns här:

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

## Steg 3: Redigera konfigurationen

Öppna `claude_desktop_config.json` och lägg till Obsidian MCP-servern:

### Alternativ A: Använd Go direkt (för utveckling)

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

### Alternativ B: Använd kompilerad binär (rekommenderat)

```json
{
  "mcpServers": {
    "obsidian": {
      "command": "/Users/faar/Documents/Src/github/fiddeb/Obsidian_mcp/obsidian-mcp-server",
      "args": [],
      "env": {
        "OBSIDIAN_API_TOKEN": "5060e7700b373ccd30eb08293d51f936808c7eb1dd650f6762d594ae97aafdcc",
        "OBSIDIAN_API_BASE_URL": "http://127.0.0.1:27123"
      }
    }
  }
}
```

## Steg 4: Starta om Claude Desktop

Efter att du har sparat konfigurationen, starta om Claude Desktop helt:

1. Stäng Claude Desktop
2. Starta Claude Desktop igen
3. MCP-servern ska nu vara tillgänglig

## Verifiera installationen

I Claude Desktop, du bör nu kunna:

1. Se Obsidian som en tillgänglig server
2. Använda verktygen genom att fråga Claude om att arbeta med dina Obsidian-anteckningar
3. Till exempel: "Kan du skapa en ny anteckning i Obsidian som heter 'Dagens tankar'?"

## Felsökning

### Problem: Server startar inte

Kontrollera att:
- Obsidian körs och Local REST API är aktiverat
- API-token är korrekt
- Sökvägen till binären eller projektet är korrekt

### Problem: Kan inte se servern i Claude

1. Kontrollera Claude Desktop logs:
   ```bash
   # macOS
   tail -f ~/Library/Logs/Claude/mcp*.log
   ```

2. Verifiera att config-filen är valid JSON
3. Kontrollera att inga andra MCP-servrar har konflikter

## För GitHub Copilot i VS Code

GitHub Copilot stödjer för närvarande inte MCP-servrar direkt via konfiguration. Men du kan:

1. Använda HTTP API:et direkt
2. Vänta på framtida stöd för MCP i VS Code
3. Använda Claude Desktop för MCP-funktionalitet

## Användning

När servern är konfigurerad kan du i Claude Desktop fråga:

- "Skapa en ny anteckning om dagens möte"
- "Sök efter alla anteckningar om projekt"
- "Visa mig innehållet i min dagboksanteckning"
- "Lista alla anteckningar i mappen Projekt"

Claude kommer automatiskt använda Obsidian MCP-servern för att utföra dessa uppgifter!
