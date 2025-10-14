# Obsidian MCP Server

En Model Context Protocol (MCP) server skriven i Go som integrerar med Obsidian via Local REST API-pluginet. Denna server låter dig komma åt och hantera din Obsidian vault programmatiskt genom MCP-protokollet.

## Funktioner

- **Anteckningshantering**: Skapa, läsa, uppdatera och ta bort anteckningar
- **Sök**: Sök genom din vault efter specifik text
- **Mapphantering**: Skapa mappar och organisera ditt innehåll
- **Vault-information**: Få översikt över din vault och dess innehåll
- **MCP-kompatibel**: Fullt kompatibel med Model Context Protocol 2024-11-05
- **Säkerhetsfunktioner**: Rate limiting, IP-filtrering, och token-baserad autentisering
- **VS Code-integration**: Fungerar sömlöst med VS Code's MCP-extension

## Förutsättningar

1. **Obsidian** med **Local REST API**-pluginet installerat och aktiverat
2. **Go 1.21** eller senare
3. Konfigurerad API-token för Local REST API

## Installation

1. Klona eller ladda ner detta projekt
2. Installera beroenden:
   ```bash
   go mod tidy
   ```

3. Konfigurera din `config.yaml` fil eller använd miljövariabler:
   ```yaml
   obsidian_api:
     base_url: "http://localhost:27123"
     token: "din-api-token-här"
     port: 27123
   
   mcp:
     host: "localhost"
     port: 8080
     description: "Obsidian MCP Server"
   ```

## Konfiguration

### Via config.yaml
Redigera `config.yaml` filen med dina inställningar:

```yaml
obsidian_api:
  base_url: "http://localhost:27123"
  token: "din-api-token-här"
  port: 27123

mcp:
  host: "localhost"
  port: 8080
  description: "Obsidian MCP Server"

security:
  enable_auth: true
  api_token: "din-säkra-token"
  enable_rate_limit: true
  rate_limit: 100  # requests per minut
  allowed_ips:
    - "127.0.0.1"
    - "::1"
```

### Via miljövariabler
Du kan också använda miljövariabler:
- `OBSIDIAN_API_TOKEN`: Din Obsidian Local REST API token
- `OBSIDIAN_API_BASE_URL`: Base URL för Obsidian API (standard: http://localhost:27123)

## Användning

### Starta servern

#### Via Go
```bash
go run .
```

#### Via kompilerad binary
```bash
# Kompilera först
go build -o obsidian-mcp-server

# Kör sedan
./obsidian-mcp-server
```

#### Via VS Code Task
1. Öppna projektet i VS Code
2. Tryck `Cmd+Shift+P` (Mac) eller `Ctrl+Shift+P` (Windows/Linux)
3. Välj "Tasks: Run Task"
4. Välj "Run Obsidian MCP Server"

Servern startar på `http://localhost:8080` (eller den port du konfigurerat).

### Integration med VS Code

För att använda MCP-servern i VS Code:

1. **Installera MCP-extension** i VS Code (om inte redan installerad)

2. **Konfigurera MCP-servern** i `.vscode/mcp.json`:
```json
{
  "mcpServers": {
    "obsidian": {
      "command": "/absolut/sökväg/till/obsidian-mcp-server"
    }
  }
}
```

3. **Starta servern** (antingen via VS Code task eller manuellt)

4. **Använd verktygen** via Copilot Chat eller MCP-kommandopaletten

### Tillgängliga verktyg

1. **get_note** - Hämta innehållet i en anteckning
   - Parameter: `path` (sökväg till anteckningen)

2. **create_note** - Skapa en ny anteckning
   - Parametrar: `path` (sökväg), `content` (innehåll)
   - *Not: Mappar skapas automatiskt när du skapar en fil i dem*

3. **update_note** - Uppdatera en befintlig anteckning
   - Parametrar: `path` (sökväg), `content` (nytt innehåll)

4. **delete_note** - Ta bort en anteckning
   - Parameter: `path` (sökväg till anteckningen)

5. **list_notes** - Lista alla anteckningar
   - Parameter: `folder` (valfri, filtrera på mapp)

6. **search_notes** - Sök efter anteckningar
   - Parameter: `query` (sökfråga)

7. **get_vault_info** - Få information om vault
   - Inga parametrar

### API-endpoints

MCP-servern stöder både root-path och `/mcp/`-prefix för kompatibilitet:

- `POST /` - Universal MCP endpoint (hanterar alla MCP-requests)
- `POST /mcp/initialize` - Initialisera MCP-sessionen
- `POST /mcp/tools/list` - Lista tillgängliga verktyg
- `POST /mcp/tools/call` - Anropa ett verktyg
- `GET /health` - Hälsokontroll (ingen autentisering krävs)

**Not:** Root-path (`/`) hanterar automatiskt JSON-RPC-notifikationer (som `notifications/initialized`) korrekt enligt MCP-protokollet.

## Exempel på MCP-anrop

### Initialisera
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

### Lista verktyg
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/list"
}
```

### Skapa en anteckning
```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "tools/call",
  "params": {
    "name": "create_note",
    "arguments": {
      "path": "Min nya anteckning.md",
      "content": "# Min nya anteckning\n\nDetta är innehållet i min anteckning."
    }
  }
}
```

## Utveckling

### Projektstruktur
```
obsidian-mcp/
├── main.go              # Huvudservern och MCP-hantering
├── obsidian_api.go      # Obsidian REST API-klient
├── security.go          # Säkerhetsfunktioner (rate limiting, IP-filtrering)
├── config.yaml          # Konfigurationsfil
├── go.mod               # Go-modulberoenden
├── go.sum               # Go-modulchecksumor
├── .vscode/
│   ├── tasks.json       # VS Code-tasks för att köra servern
│   └── mcp.json         # MCP-serverkonfiguration för VS Code
├── .gitignore           # Git ignore-fil
└── README.md            # Denna fil
```

### Bygga för produktion
```bash
go build -o obsidian-mcp-server .
```

### Testa anslutningen
```bash
curl http://localhost:8080/health
```

## Felsökning

### Vanliga problem

1. **"connection refused"**: 
   - Kontrollera att Obsidian är igång och Local REST API-pluginet är aktiverat
   - Verifiera att porten (standard 27123) är korrekt i konfigurationen

2. **"unauthorized"**: 
   - Verifiera din API-token i konfigurationen
   - Kontrollera att tokenen matchar den i Obsidian Local REST API-inställningarna

3. **"not found"**: 
   - Kontrollera att anteckningssökvägarna är korrekta (relativt till vault-roten)
   - Använd `.md`-filändelsen i sökvägen

4. **"Method not found: notifications/initialized"**:
   - Detta har nu åtgärdats - servern hanterar notifikationer korrekt
   - Om problemet kvarstår, säkerställ att du kör den senaste versionen

5. **Rate limit-fel**:
   - Justera `rate_limit` i `config.yaml` om du behöver fler requests per minut
   - Eller inaktivera rate limiting genom att sätta `enable_rate_limit: false`

6. **IP-blockeringsfel**:
   - Lägg till din IP-adress i `allowed_ips`-listan i `config.yaml`
   - Eller inaktivera IP-filtrering genom att lämna `allowed_ips` tom

### Loggar

Servern loggar all aktivitet till stdout. Kontrollera loggarna för detaljerad felsökningsinformation:
- Säkerhetsåtgärder (autentisering, rate limiting, IP-filtrering)
- MCP-requests och svar
- Obsidian API-anrop och resultat

## Bidrag

Bidrag är välkomna! Skapa gärna issues eller pull requests.

## Licens

MIT License