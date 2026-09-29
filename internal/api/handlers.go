package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/relentlessworks/phonekit/internal/model"
)

// Handler holds dependencies for HTTP handlers.
type Handler struct {
	secret string
}

// NewHandler creates a new API handler.
func NewHandler(secret string) *Handler {
	return &Handler{secret: secret}
}

// Register sets up all routes on the given mux.
func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/help", h.help)
	mux.HandleFunc("/.well-known/agent.md", h.help)
	mux.HandleFunc("/parse", h.parse)
	mux.HandleFunc("/validate", h.validate)
	mux.HandleFunc("/format", h.format)
	mux.HandleFunc("/info", h.info)
	mux.HandleFunc("/countries", h.countries)
	mux.HandleFunc("/country/", h.countryLookup)
	mux.HandleFunc("/mcp", h.mcp)
	mux.HandleFunc("/", h.root)
}

// wantsJSON checks if the client wants JSON output.
func wantsJSON(r *http.Request) bool {
	if r.URL.Query().Get("format") == "json" {
		return true
	}
	return strings.Contains(r.Header.Get("Accept"), "application/json")
}

// writeError writes an instructive error response.
func writeError(w http.ResponseWriter, r *http.Request, status int, msg, hint string) {
	if wantsJSON(r) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(map[string]string{
			"error": msg,
			"hint":  hint,
		})
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	fmt.Fprintf(w, "error: %s | hint: %s\n", msg, hint)
}

// writeRecord writes a plain text record or JSON.
func writeRecord(w http.ResponseWriter, r *http.Request, text string, data interface{}) {
	if wantsJSON(r) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintln(w, text)
}

// root handles GET / with service info.
func (h *Handler) root(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if wantsJSON(r) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"service": "phonekit",
			"version": "0.1.0",
			"status":  "ok",
			"help":    "GET /help for operating manual",
		})
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintln(w, "service=phonekit version=0.1.0 status=ok")
	fmt.Fprintln(w, "hint: GET /help for operating manual")
}

// help returns the self-documenting operating manual.
func (h *Handler) help(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, helpText)
}

const helpText = `# phonekit — Agentic-First Phone Number Utility

Parse, validate, format, and get info about phone numbers.

## Endpoints

### GET /parse?number=PHONE
Parse a phone number and return structured info.
Example: GET /parse?number=%2B14155551234
Output: raw=+14155551234 country_code=1 national=4155551234 area=415 subscriber=5551234 e164=+14155551234 international=+1 415 555 1234 national=(415) 555-1234 rfc3966=tel:+14155551234 country=United States iso=US line_type=fixed_or_mobile valid=true

### GET /validate?number=PHONE
Check if a phone number is valid (basic check).
Example: GET /validate?number=%2B14155551234
Output: number=+14155551234 valid=true

### GET /format?number=PHONE&format=FORMAT
Convert a phone number to a specific format.
Formats: e164, international, national, rfc3966
Example: GET /format?number=%2B14155551234&format=e164
Output: number=+14155551234 format=e164 result=+14155551234

### GET /info?number=PHONE
Get detailed info about a phone number (same as parse but focused on metadata).
Example: GET /info?number=%2B442079460958
Output: raw=+442079460958 country_code=44 national=2079460958 area=20 subscriber=79460958 e164=+442079460958 international=+44 20 79 46 09 58 country=United Kingdom iso=GB line_type=fixed valid=true

### GET /countries
List all known country calling codes.
Output: one line per country: code=1 iso=US name=United States min=10 max=10

### GET /country/CODE
Lookup a country by its calling code.
Example: GET /country/44
Output: code=44 iso=GB name=United Kingdom min=9 max=10

### POST /mcp
MCP JSON-RPC 2.0 endpoint. Tools: parse, validate, format, info, list_countries, lookup_country.

## Response Format
- Plain text by default (one labeled line per record)
- JSON via Accept: application/json or ?format=json
- Errors: error: message | hint: what to do next

## Notes
- No database needed — pure stateless computation
- No auth required for basic operations
- Country code detection is best-effort without a full phone library
- Line type detection is approximate based on number patterns
`

// parse handles GET /parse
func (h *Handler) parse(w http.ResponseWriter, r *http.Request) {
	number := r.URL.Query().Get("number")
	if number == "" {
		writeError(w, r, http.StatusBadRequest, "missing number parameter", "provide a phone number via ?number= e.g. /parse?number=%2B14155551234")
		return
	}

	pn, err := model.Parse(number)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error(), "provide a valid phone number with at least 4 digits and at most 15 digits")
		return
	}

	if wantsJSON(r) {
		writeRecord(w, r, "", pn)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "raw=%s country_code=%s national=%s area=%s subscriber=%s e164=%s international=%s national_fmt=%s rfc3966=%s country=%s iso=%s line_type=%s valid=%t\n",
		pn.Raw, pn.CountryCode, pn.NationalNum, pn.AreaCode, pn.Subscriber,
		pn.E164, pn.International, pn.National, pn.RFC3966,
		pn.Country, pn.CountryISO, pn.LineType, pn.Valid)
}

// validate handles GET /validate
func (h *Handler) validate(w http.ResponseWriter, r *http.Request) {
	number := r.URL.Query().Get("number")
	if number == "" {
		writeError(w, r, http.StatusBadRequest, "missing number parameter", "provide a phone number via ?number= e.g. /validate?number=%2B14155551234")
		return
	}

	pn, err := model.Parse(number)
	valid := false
	if err == nil {
		valid = pn.Valid
	}

	if wantsJSON(r) {
		writeRecord(w, r, "", map[string]interface{}{
			"number": number,
			"valid":  valid,
		})
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "number=%s valid=%t\n", number, valid)
}

// format handles GET /format
func (h *Handler) format(w http.ResponseWriter, r *http.Request) {
	number := r.URL.Query().Get("number")
	if number == "" {
		writeError(w, r, http.StatusBadRequest, "missing number parameter", "provide a phone number via ?number= e.g. /format?number=%2B14155551234&format=e164")
		return
	}

	format := r.URL.Query().Get("format")
	if format == "" {
		writeError(w, r, http.StatusBadRequest, "missing format parameter", "provide a format via ?format= (e164, international, national, rfc3966)")
		return
	}

	result, err := model.Format(number, format)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error(), "use one of: e164, international, national, rfc3966")
		return
	}

	if wantsJSON(r) {
		writeRecord(w, r, "", map[string]string{
			"number": number,
			"format": format,
			"result": result,
		})
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "number=%s format=%s result=%s\n", number, format, result)
}

// info handles GET /info
func (h *Handler) info(w http.ResponseWriter, r *http.Request) {
	number := r.URL.Query().Get("number")
	if number == "" {
		writeError(w, r, http.StatusBadRequest, "missing number parameter", "provide a phone number via ?number= e.g. /info?number=%2B14155551234")
		return
	}

	pn, err := model.Parse(number)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error(), "provide a valid phone number with at least 4 digits and at most 15 digits")
		return
	}

	if wantsJSON(r) {
		writeRecord(w, r, "", pn)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "raw=%s country_code=%s national=%s area=%s subscriber=%s e164=%s international=%s national_fmt=%s rfc3966=%s country=%s iso=%s line_type=%s valid=%t\n",
		pn.Raw, pn.CountryCode, pn.NationalNum, pn.AreaCode, pn.Subscriber,
		pn.E164, pn.International, pn.National, pn.RFC3966,
		pn.Country, pn.CountryISO, pn.LineType, pn.Valid)
}

// countries handles GET /countries
func (h *Handler) countries(w http.ResponseWriter, r *http.Request) {
	countries := model.ListCountries()

	if wantsJSON(r) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(countries)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	for _, c := range countries {
		fmt.Fprintf(w, "code=%s iso=%s name=%s min=%d max=%d\n", c.Code, c.ISO, c.Name, c.MinLen, c.MaxLen)
	}
}

// countryLookup handles GET /country/CODE
func (h *Handler) countryLookup(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/country/")
	if code == "" {
		writeError(w, r, http.StatusBadRequest, "missing country code", "provide a calling code in the path e.g. /country/44")
		return
	}

	info, ok := model.LookupCountry(code)
	if !ok {
		writeError(w, r, http.StatusNotFound, "country code not found", "try GET /countries to list all known country calling codes")
		return
	}

	if wantsJSON(r) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(info)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "code=%s iso=%s name=%s min=%d max=%d\n", info.Code, info.ISO, info.Name, info.MinLen, info.MaxLen)
}
