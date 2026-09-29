package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/relentlessworks/phonekit/internal/model"
)

// MCPRequest is a JSON-RPC 2.0 request.
type MCPRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

// MCPResponse is a JSON-RPC 2.0 response.
type MCPResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Result  interface{}     `json:"result,omitempty"`
	Error   *MCPError       `json:"error,omitempty"`
}

// MCPError is a JSON-RPC 2.0 error.
type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// MCPTool represents a tool in tools/list.
type MCPTool struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	InputSchema MCPInputSchema    `json:"inputSchema"`
}

// MCPInputSchema describes the input schema for a tool.
type MCPInputSchema struct {
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
	Required   []string               `json:"required,omitempty"`
}

// mcp handles POST /mcp
func (h *Handler) mcp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, r, http.StatusMethodNotAllowed, "method not allowed", "use POST for MCP endpoint")
		return
	}

	var req MCPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(MCPResponse{
			JSONRPC: "2.0",
			Error:   &MCPError{Code: -32700, Message: "parse error"},
		})
		return
	}

	resp := h.handleMCP(&req)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) handleMCP(req *MCPRequest) MCPResponse {
	switch req.Method {
	case "initialize":
		return MCPResponse{
			JSONRPC: "2.0",
			ID:     req.ID,
			Result: map[string]interface{}{
				"protocolVersion": "2024-11-05",
				"capabilities": map[string]interface{}{
					"tools": map[string]interface{}{},
				},
				"serverInfo": map[string]string{
					"name":    "phonekit",
					"version": "0.1.0",
				},
			},
		}

	case "tools/list":
		return MCPResponse{
			JSONRPC: "2.0",
			ID:     req.ID,
			Result: map[string]interface{}{
				"tools": mcpTools(),
			},
		}

	case "tools/call":
		return h.handleToolCall(req)

	default:
		return MCPResponse{
			JSONRPC: "2.0",
			ID:     req.ID,
			Error:  &MCPError{Code: -32601, Message: "method not found: " + req.Method},
		}
	}
}

func mcpTools() []MCPTool {
	return []MCPTool{
		{
			Name:        "parse",
			Description: "Parse a phone number and return structured info (country code, national number, area code, subscriber, E.164, international, national, RFC 3966 formats, country, line type, validity).",
			InputSchema: MCPInputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"number": map[string]string{
						"type":        "string",
						"description": "Phone number to parse, e.g. +14155551234",
					},
				},
				Required: []string{"number"},
			},
		},
		{
			Name:        "validate",
			Description: "Check if a phone number is valid (basic check).",
			InputSchema: MCPInputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"number": map[string]string{
						"type":        "string",
						"description": "Phone number to validate",
					},
				},
				Required: []string{"number"},
			},
		},
		{
			Name:        "format",
			Description: "Convert a phone number to a specific format (e164, international, national, rfc3966).",
			InputSchema: MCPInputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"number": map[string]string{
						"type":        "string",
						"description": "Phone number to format",
					},
					"format": map[string]string{
						"type":        "string",
						"description": "Target format: e164, international, national, rfc3966",
					},
				},
				Required: []string{"number", "format"},
			},
		},
		{
			Name:        "info",
			Description: "Get detailed info about a phone number (same as parse).",
			InputSchema: MCPInputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"number": map[string]string{
						"type":        "string",
						"description": "Phone number to get info for",
					},
				},
				Required: []string{"number"},
			},
		},
		{
			Name:        "list_countries",
			Description: "List all known country calling codes with ISO codes, names, and national number length ranges.",
			InputSchema: MCPInputSchema{
				Type:       "object",
				Properties: map[string]interface{}{},
			},
		},
		{
			Name:        "lookup_country",
			Description: "Lookup a country by its calling code.",
			InputSchema: MCPInputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"code": map[string]string{
						"type":        "string",
						"description": "Country calling code, e.g. 1, 44, 49",
					},
				},
				Required: []string{"code"},
			},
		},
	}
}

func (h *Handler) handleToolCall(req *MCPRequest) MCPResponse {
	var params struct {
		Name      string            `json:"name"`
		Arguments map[string]string `json:"arguments"`
	}
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return MCPResponse{
			JSONRPC: "2.0",
			ID:     req.ID,
			Error:  &MCPError{Code: -32602, Message: "invalid params"},
		}
	}

	switch params.Name {
	case "parse", "info":
		number := params.Arguments["number"]
		if number == "" {
			return mcpError(req.ID, -32602, "missing number argument")
		}
		pn, err := model.Parse(number)
		if err != nil {
			return mcpError(req.ID, -32603, err.Error())
		}
		return MCPResponse{
			JSONRPC: "2.0",
			ID:     req.ID,
			Result: map[string]interface{}{
				"content": []map[string]string{
					{"type": "text", "text": fmt.Sprintf("raw=%s country_code=%s national=%s area=%s subscriber=%s e164=%s international=%s national_fmt=%s rfc3966=%s country=%s iso=%s line_type=%s valid=%t",
						pn.Raw, pn.CountryCode, pn.NationalNum, pn.AreaCode, pn.Subscriber,
						pn.E164, pn.International, pn.National, pn.RFC3966,
						pn.Country, pn.CountryISO, pn.LineType, pn.Valid)},
				},
			},
		}

	case "validate":
		number := params.Arguments["number"]
		if number == "" {
			return mcpError(req.ID, -32602, "missing number argument")
		}
		pn, err := model.Parse(number)
		valid := false
		if err == nil {
			valid = pn.Valid
		}
		return MCPResponse{
			JSONRPC: "2.0",
			ID:     req.ID,
			Result: map[string]interface{}{
				"content": []map[string]string{
					{"type": "text", "text": fmt.Sprintf("number=%s valid=%t", number, valid)},
				},
			},
		}

	case "format":
		number := params.Arguments["number"]
		format := params.Arguments["format"]
		if number == "" {
			return mcpError(req.ID, -32602, "missing number argument")
		}
		if format == "" {
			return mcpError(req.ID, -32602, "missing format argument")
		}
		result, err := model.Format(number, format)
		if err != nil {
			return mcpError(req.ID, -32603, err.Error())
		}
		return MCPResponse{
			JSONRPC: "2.0",
			ID:     req.ID,
			Result: map[string]interface{}{
				"content": []map[string]string{
					{"type": "text", "text": fmt.Sprintf("number=%s format=%s result=%s", number, format, result)},
				},
			},
		}

	case "list_countries":
		countries := model.ListCountries()
		var lines []string
		for _, c := range countries {
			lines = append(lines, fmt.Sprintf("code=%s iso=%s name=%s min=%d max=%d", c.Code, c.ISO, c.Name, c.MinLen, c.MaxLen))
		}
		return MCPResponse{
			JSONRPC: "2.0",
			ID:     req.ID,
			Result: map[string]interface{}{
				"content": []map[string]string{
					{"type": "text", "text": strings.Join(lines, "\n")},
				},
			},
		}

	case "lookup_country":
		code := params.Arguments["code"]
		if code == "" {
			return mcpError(req.ID, -32602, "missing code argument")
		}
		info, ok := model.LookupCountry(code)
		if !ok {
			return mcpError(req.ID, -32603, "country code not found")
		}
		return MCPResponse{
			JSONRPC: "2.0",
			ID:     req.ID,
			Result: map[string]interface{}{
				"content": []map[string]string{
					{"type": "text", "text": fmt.Sprintf("code=%s iso=%s name=%s min=%d max=%d", info.Code, info.ISO, info.Name, info.MinLen, info.MaxLen)},
				},
			},
		}

	default:
		return mcpError(req.ID, -32602, "unknown tool: "+params.Name)
	}
}

func mcpError(id json.RawMessage, code int, msg string) MCPResponse {
	return MCPResponse{
		JSONRPC: "2.0",
		ID:     id,
		Error:  &MCPError{Code: code, Message: msg},
	}
}
