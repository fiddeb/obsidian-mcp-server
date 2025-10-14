# Obsidian MCP Server - Dokumentation

En Model Context Protocol (MCP) server skriven i Go som integrerar med Obsidian via Local REST API-pluginet.

## 🚀 Snabbstart

1. Installera Obsidian Local REST API plugin
2. Konfigurera API-token i config.yaml
3. Kör: `go run .`
4. Server körs på http://localhost:8080

## 📋 Tillgängliga Verktyg

### 1. get_note
Hämtar innehåll från en anteckning.

**Parametrar:**
- `path` (string): Sökväg till anteckningen

**Exempel:**
```json
{
  "name": "get_note",
  "arguments": {
    "path": "MCP-server/README.md"
  }
}
```

### 2. create_note
Skapar en ny anteckning.

**Parametrar:**
- `path` (string): Sökväg där anteckningen ska skapas
- `content` (string): Innehåll i anteckningen

**Exempel:**
```json
{
  "name": "create_note",
  "arguments": {
    "path": "Dagbok/2025-10-14.md",
    "content": "# Min dagbok\\n\\nIdag lärde jag mig om MCP!"
  }
}
```

### 3. update_note
Uppdaterar befintlig anteckning.

**Parametrar:**
- `path` (string): Sökväg till anteckningen
- `content` (string): Nytt innehåll

### 4. delete_note
Tar bort en anteckning.

**Parametrar:**
- `path` (string): Sökväg till anteckningen

### 5. list_notes
Listar alla anteckningar i vault eller specifik mapp.

**Parametrar:**
- `folder` (string, optional): Mapp att filtrera på

### 6. search_notes
Söker efter anteckningar med kontext.

**Parametrar:**
- `query` (string): Sökfråga

**Returnerar:**
- Filnamn
- Kontext runt varje match
- Relevanspoäng

### 7. get_vault_info
Hämtar information om vault.

**Returnerar:**
- Autentiseringsstatus
- Plugin-version
- Obsidian-version
- Antal anteckningar

### 8. create_folder
Skapar ny mapp i vault.

**Parametrar:**
- `path` (string): Sökväg till mappen

## 💡 Användningsexempel

### Skapa en anteckning via curl

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
        "path": "Min anteckning.md",
        "content": "# Hello World\\n\\nDetta är min första anteckning via MCP!"
      }
    }
  }'
```

### Söka efter anteckningar

```bash
curl -X POST http://localhost:8080/mcp/tools/call \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 2,
    "method": "tools/call",
    "params": {
      "name": "search_notes",
      "arguments": {
        "query": "projekt"
      }
    }
  }'
```

### Lista alla anteckningar

```bash
curl -X POST http://localhost:8080/mcp/tools/call \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 3,
    "method": "tools/call",
    "params": {
      "name": "list_notes",
      "arguments": {}
    }
  }'
```

## ⚙️ Konfiguration

### Via config.yaml

```yaml
obsidian_api:
  base_url: "http://127.0.0.1:27123"
  token: "din-api-token-här"
  port: 27123

mcp:
  host: "localhost"
  port: 8080
  description: "Obsidian MCP Server"
```

### Via miljövariabler

```bash
export OBSIDIAN_API_TOKEN="din-api-token"
export OBSIDIAN_API_BASE_URL="http://127.0.0.1:27123"
```

### Hitta din API-token

1. Öppna Obsidian
2. Gå till **Settings → Community Plugins → Local REST API**
3. Kopiera API-token som visas

## 🔧 API Endpoints

### MCP-protokoll

- `POST /mcp/initialize` - Initialisera MCP-sessionen
- `POST /mcp/tools/list` - Lista tillgängliga verktyg
- `POST /mcp/tools/call` - Anropa ett verktyg

### Utility

- `GET /health` - Hälsokontroll (returnerar `{"status":"healthy"}`)

## 🐛 Felsökning

### Problem: "connection refused"

**Orsak:** Obsidian körs inte eller Local REST API är inte aktiverat.

**Lösning:**
1. Starta Obsidian
2. Aktivera Local REST API plugin
3. Verifiera att plugin körs på port 27123

### Problem: "unauthorized"

**Orsak:** Felaktig API-token.

**Lösning:**
1. Kontrollera token i config.yaml
2. Kopiera token från Obsidian Settings → Local REST API
3. Se till att det inte finns extra mellanslag

### Problem: "not found"

**Orsak:** Felaktig filsökväg.

**Lösning:**
- Använd relativa sökvägar från vault root
- Inkludera .md extension
- Exempel: `"MCP-server/README.md"` inte `"/Users/.../vault/MCP-server/README.md"`

### Problem: Port 8080 används

**Lösning:**
```bash
# Döda process på port 8080
lsof -ti:8080 | xargs kill -9

# Eller ändra port i config.yaml
```

## 📦 Installation

### Förutsättningar

- Go 1.21 eller senare
- Obsidian med Local REST API plugin

### Installationssteg

```bash
# 1. Klona repository
git clone <repository-url>
cd Obsidian_mcp

# 2. Installera dependencies
go mod tidy

# 3. Konfigurera
# Redigera config.yaml med din API-token

# 4. Starta server
go run .

# Eller bygg binary
go build -o obsidian-mcp-server .
./obsidian-mcp-server
```

## 🏗️ Projektstruktur

```
Obsidian_mcp/
├── main.go              # MCP server & handlers
├── obsidian_api.go      # Obsidian REST API client
├── config.yaml          # Konfiguration
├── go.mod               # Go dependencies
├── .gitignore          # Git ignore rules
└── README.md           # Projektdokumentation
```

## 🔒 Säkerhet

### Rekommendationer

1. **API-token**: Förvara ALDRIG token i version control
2. **Nätverk**: Kör endast på localhost om inte nödvändigt
3. **Backup**: Ha regelbundna backups av vault

### Säker konfiguration

```bash
# Använd miljövariabler istället
export OBSIDIAN_API_TOKEN="secret-token"
go run .
```

## 🚀 Utveckling

### Bygg för olika plattformar

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o obsidian-mcp-server-linux .

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o obsidian-mcp-server-mac-intel .

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o obsidian-mcp-server-mac-arm .

# Windows
GOOS=windows GOARCH=amd64 go build -o obsidian-mcp-server.exe .
```

### Kör i bakgrunden

```bash
# macOS/Linux
nohup ./obsidian-mcp-server > server.log 2>&1 &

# Stoppa
pkill -f obsidian-mcp-server
```

## 📝 MCP Protocol Exempel

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
      "name": "my-client",
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

### Call Tool

```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "tools/call",
  "params": {
    "name": "get_note",
    "arguments": {
      "path": "test.md"
    }
  }
}
```

## 📊 Funktioner som används från Obsidian API

- **GET /vault/{path}** - Hämta filinnehåll
- **PUT /vault/{path}** - Skapa/uppdatera fil
- **DELETE /vault/{path}** - Ta bort fil
- **GET /vault/** - Lista filer
- **POST /search/simple/** - Söka med kontext
- **GET /** - Vault-information

## 🎯 Användningsfall

1. **Automatisering**: Skapa dagboksanteckningar automatiskt
2. **Integration**: Koppla Obsidian till andra system
3. **Sökning**: Avancerad sökning med programmatisk åtkomst
4. **Backup**: Automatiska exports av anteckningar
5. **AI-integration**: Använd med LLM för att analysera/generera innehåll

## 📄 Licens

MIT License

## 🤝 Bidrag

Bidrag välkomnas! Öppna issues eller pull requests.

---

**Version:** 1.0.0  
**Datum:** 2025-10-14  
**Status:** ✅ Production Ready
