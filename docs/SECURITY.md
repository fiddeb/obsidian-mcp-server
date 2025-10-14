# Säkerhetsguide för Obsidian MCP Server

## 🔒 Säkerhetsnivåer

### Nivå 1: Grundläggande (Standard)
MCP-servern lyssnar endast på localhost och är inte tillgänglig från nätverket.

**Konfiguration:**
```yaml
mcp:
  host: "localhost"  # Endast lokala anslutningar
  port: 8080

security:
  allowed_ips:
    - "127.0.0.1"
    - "::1"
```

### Nivå 2: Med autentisering
Kräver en Bearer token för alla anrop.

**Konfiguration:**
```yaml
security:
  enable_auth: true
  auth_token: "din-hemliga-token-här"  # Generera med: openssl rand -hex 32
```

**Användning:**
```bash
curl -H "Authorization: Bearer din-hemliga-token-här" \
  http://localhost:8080/mcp/tools/list
```

### Nivå 3: Med rate limiting
Begränsar antal anrop per minut.

**Konfiguration:**
```yaml
security:
  enable_rate_limit: true
  rate_limit: 60  # Max 60 anrop/minut per IP
```

### Nivå 4: Full säkerhet
Kombination av alla säkerhetsåtgärder.

**Konfiguration:**
```yaml
mcp:
  host: "localhost"
  port: 8080

security:
  enable_auth: true
  auth_token: "your-strong-secret-token"
  allowed_ips:
    - "127.0.0.1"
    - "::1"
  enable_rate_limit: true
  rate_limit: 100
  enable_cors: false
```

## 🛡️ Säkerhetsåtgärder som är implementerade

### 1. IP-vitlista
```yaml
security:
  allowed_ips:
    - "127.0.0.1"      # IPv4 localhost
    - "::1"            # IPv6 localhost
    - "192.168.1.0/24" # Hela subnettet (CIDR-notation)
```

### 2. Path Validation
Automatiskt skydd mot:
- Directory traversal (`../../../etc/passwd`)
- Absoluta sökvägar (`/etc/passwd`)
- Farliga tecken (`` ` ``, `$`, `|`, `;`, `&`)
- Null bytes (`\x00`)

### 3. Input Sanitization
- Begränsar request body till 1MB
- Validerar Content-Type
- Tar bort farliga tecken från innehåll

### 4. Rate Limiting
- Per IP-adress
- Konfigurerbar gräns per minut
- Automatisk rensning av gamla requests

### 5. CORS-skydd
```yaml
security:
  enable_cors: true
  allowed_origins:
    - "http://localhost:3000"
    - "http://127.0.0.1:5173"
```

### 6. Audit Logging
Aktivera med miljövariabel:
```bash
export ENABLE_AUDIT_LOG=true
go run .
```

Loggar till stderr i JSON-format:
```json
{
  "timestamp": "2025-10-14T12:00:00Z",
  "ip": "127.0.0.1",
  "action": "create_note",
  "result": "success",
  "path": "/mcp/tools/call"
}
```

## 🔑 Generera säkra tokens

### För autentisering
```bash
# Generera en stark token
openssl rand -hex 32

# Eller med base64
openssl rand -base64 32
```

Lägg sedan till i config.yaml:
```yaml
security:
  enable_auth: true
  auth_token: "b5f8d9e2c4a1f3e7d8c9b2a5f6e1d4c3a7b8e9f2d5c1a6e3f8d9c2b7a4e1f6"
```

## 🚨 Bästa praxis

### 1. Kör endast på localhost
```yaml
mcp:
  host: "localhost"  # INTE "0.0.0.0"
```

### 2. Använd alltid IP-vitlista
```yaml
security:
  allowed_ips:
    - "127.0.0.1"
    - "::1"
```

### 3. Aktivera autentisering för produktion
```yaml
security:
  enable_auth: true
  auth_token: "generated-secure-token"
```

### 4. Förvara tokens säkert
- Använd miljövariabler istället för config-filer
- Lägg ALDRIG till tokens i version control
- Rotera tokens regelbundet

```bash
export MCP_AUTH_TOKEN="your-secret-token"
go run .
```

### 5. Aktivera rate limiting
```yaml
security:
  enable_rate_limit: true
  rate_limit: 100  # Anpassa efter behov
```

### 6. Övervaka via audit logs
```bash
ENABLE_AUDIT_LOG=true go run . 2>&1 | tee mcp-audit.log
```

## 🌐 Nätverksexponering (VARNING)

Om du **måste** exponera servern i nätverket:

### Alternativ 1: Reverse Proxy (Rekommenderat)
Använd nginx eller Caddy med HTTPS:

```nginx
server {
    listen 443 ssl;
    server_name mcp.example.com;
    
    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;
    
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

### Alternativ 2: SSH Tunnel
```bash
# På remote server
ssh -L 8080:localhost:8080 user@remote-server

# Nu kan du accessa via localhost:8080 säkert
```

### Alternativ 3: VPN
- Använd Tailscale, WireGuard eller liknande
- Exponera endast inom VPN-nätverket

## 🔍 Validering av sökvägar

Servern validerar automatiskt alla sökvägar:

**✅ Godkända sökvägar:**
```
"Projekt/README.md"
"Daily/2025-10-14.md"
"Möten/Team meeting.md"
```

**❌ Blockerade sökvägar:**
```
"../../../etc/passwd"          # Directory traversal
"/absolute/path/file.md"       # Absolut sökväg
"file`dangerous`.md"           # Farliga tecken
"path;rm -rf /"               # Command injection
```

## 🧪 Testa säkerheten

### 1. Test IP-restriktion
```bash
# Detta ska INTE fungera från annan maskin
curl http://server-ip:8080/health
# Förväntat: Connection refused eller Forbidden
```

### 2. Test autentisering
```bash
# Utan token
curl http://localhost:8080/mcp/tools/list
# Förväntat: 401 Unauthorized

# Med korrekt token
curl -H "Authorization: Bearer your-token" \
  http://localhost:8080/mcp/tools/list
# Förväntat: 200 OK
```

### 3. Test rate limiting
```bash
# Skicka många requests
for i in {1..100}; do
  curl http://localhost:8080/health
done
# Förväntat: 429 Too Many Requests efter gränsen
```

### 4. Test path validation
```bash
curl -X POST http://localhost:8080/mcp/tools/call \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "params": {
      "name": "get_note",
      "arguments": {"path": "../../../etc/passwd"}
    }
  }'
# Förväntat: Error om ogiltig sökväg
```

## 📋 Säkerhetschecklista

- [ ] Servern lyssnar endast på localhost
- [ ] IP-vitlista är konfigurerad
- [ ] Autentisering är aktiverad för produktion
- [ ] Rate limiting är aktiverat
- [ ] Audit logging är aktiverat
- [ ] Tokens förvaras säkert (ej i git)
- [ ] HTTPS används för nätverksexponering
- [ ] Regelbunden övervakning av audit logs
- [ ] Backup av vault finns
- [ ] Säkerhetsuppdateringar installeras

## 🆘 Vid säkerhetsincident

1. **Stoppa servern omedelbart**
   ```bash
   pkill -f obsidian-mcp-server
   ```

2. **Rotera alla tokens**
   ```bash
   # Generera ny token
   openssl rand -hex 32
   # Uppdatera config.yaml och Obsidian API token
   ```

3. **Granska audit logs**
   ```bash
   grep "AUDIT" mcp-audit.log | jq .
   ```

4. **Kontrollera vault för oväntade ändringar**
   ```bash
   # Om du har git i din vault
   cd ~/vault
   git log --all --oneline
   git diff
   ```

5. **Rapportera till utvecklare**
   - Öppna ett issue på GitHub
   - Inkludera relevanta logs (ta bort känslig info)

## 📚 Ytterligare resurser

- [OWASP API Security Top 10](https://owasp.org/www-project-api-security/)
- [Go Security Best Practices](https://go.dev/doc/security/best-practices)
- [Rate Limiting Strategies](https://cloud.google.com/architecture/rate-limiting-strategies)

---

**Kom ihåg:** Säkerhet är en process, inte en produkt. Granska och uppdatera säkerhetsinställningarna regelbundet!
