# phonekit

Agentic-first phone number parsing, validation, and formatting service. Parse, validate, and format phone numbers in E.164, national, and international formats. Extract country codes, detect line types, and get location info. Plain text API, agent-driven, single Go binary.

## Quick Start

```bash
# Build
make build

# Run
./phonekit

# Parse a number
curl "http://localhost:8080/parse?number=%2B14155551234"

# Validate
curl "http://localhost:8080/validate?number=%2B14155551234"

# Format
curl "http://localhost:8080/format?number=%2B14155551234&format=e164"

# Get info
curl "http://localhost:8080/info?number=%2B442079460958"

# List countries
curl "http://localhost:8080/countries"

# Lookup country
curl "http://localhost:8080/country/44"
```

## API Reference

### GET /parse?number=PHONE
Parse a phone number and return structured info.

**Output:** `raw=+14155551234 country_code=1 national=4155551234 area=415 subscriber=5551234 e164=+14155551234 international=+1 415 555 1234 national_fmt=(415) 555-1234 rfc3966=tel:+14155551234 country=United States iso=US line_type=fixed_or_mobile valid=true`

### GET /validate?number=PHONE
Check if a phone number is valid (basic check).

**Output:** `number=+14155551234 valid=true`

### GET /format?number=PHONE&format=FORMAT
Convert a phone number to a specific format.

**Formats:** `e164`, `international`, `national`, `rfc3966`

**Output:** `number=+14155551234 format=e164 result=+14155551234`

### GET /info?number=PHONE
Get detailed info about a phone number.

### GET /countries
List all known country calling codes.

**Output:** One line per country: `code=1 iso=US name=United States min=10 max=10`

### GET /country/CODE
Lookup a country by its calling code.

**Output:** `code=44 iso=GB name=United Kingdom min=9 max=10`

### POST /mcp
MCP JSON-RPC 2.0 endpoint. Tools: `parse`, `validate`, `format`, `info`, `list_countries`, `lookup_country`.

### GET /help
Self-documenting operating manual.

## Response Format

- **Plain text by default** — one labeled, grepable line per record
- **JSON on demand** — `Accept: application/json` or `?format=json`
- **Instructive errors** — `error: message | hint: what to do next`

## Configuration

| Flag | Env | Default | Description |
|------|-----|---------|-------------|
| `-addr` | `PHONEKIT_ADDR` | `:8080` | Listen address |
| `-secret` | `PHONEKIT_SECRET` | (auto) | Auth token secret |

## Build

```bash
make build    # CGO_ENABLED=0, single binary
make test     # go test -race
make vet      # go vet
```

## Design

- **No database** — pure stateless computation
- **No auth required** for basic operations
- **Single Go binary** — CGO_ENABLED=0, zero external dependencies
- **80+ country calling codes** with ISO codes and national number length ranges
- **Line type detection** — mobile, fixed, tollfree, premium (best-effort based on number patterns)
- **Multiple formats** — E.164, international, national, RFC 3966 (tel: URI)

## License

MIT
