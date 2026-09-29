package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRoot(t *testing.T) {
	h := NewHandler("test")
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "service=phonekit") {
		t.Errorf("expected service=phonekit in response, got %s", w.Body.String())
	}
}

func TestHelp(t *testing.T) {
	h := NewHandler("test")
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("GET", "/help", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "phonekit") {
		t.Error("expected phonekit in help text")
	}
	if !strings.Contains(w.Body.String(), "/parse") {
		t.Error("expected /parse in help text")
	}
}

func TestParseEndpoint(t *testing.T) {
	h := NewHandler("test")
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("GET", "/parse?number=%2B14155551234", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "country_code=1") {
		t.Errorf("expected country_code=1, got %s", body)
	}
	if !strings.Contains(body, "iso=US") {
		t.Errorf("expected iso=US, got %s", body)
	}
	if !strings.Contains(body, "e164=+14155551234") {
		t.Errorf("expected e164=+14155551234, got %s", body)
	}
}

func TestParseEndpointJSON(t *testing.T) {
	h := NewHandler("test")
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("GET", "/parse?number=%2B14155551234&format=json", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
		var result map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatalf("failed to parse JSON: %v", err)
		}
		if result["country_code"] != "1" {
			t.Errorf("expected country_code 1, got %v", result["country_code"])
		}
	}
}

func TestParseMissingNumber(t *testing.T) {
	h := NewHandler("test")
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("GET", "/parse", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "missing number") {
		t.Errorf("expected missing number error, got %s", w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "hint:") {
		t.Error("expected hint in error response")
	}
}

func TestValidateEndpoint(t *testing.T) {
	h := NewHandler("test")
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("GET", "/validate?number=%2B14155551234", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "valid=true") {
		t.Errorf("expected valid=true, got %s", w.Body.String())
	}
}

func TestValidateInvalidNumber(t *testing.T) {
	h := NewHandler("test")
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("GET", "/validate?number=abc", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "valid=false") {
		t.Errorf("expected valid=false, got %s", w.Body.String())
	}
}

func TestFormatEndpoint(t *testing.T) {
	h := NewHandler("test")
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("GET", "/format?number=%2B14155551234&format=e164", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "result=+14155551234") {
		t.Errorf("expected result=+14155551234, got %s", w.Body.String())
	}
}

func TestFormatMissingFormat(t *testing.T) {
	h := NewHandler("test")
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("GET", "/format?number=%2B14155551234", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestInfoEndpoint(t *testing.T) {
	h := NewHandler("test")
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("GET", "/info?number=%2B442079460958", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "country_code=44") {
		t.Errorf("expected country_code=44, got %s", body)
	}
	if !strings.Contains(body, "iso=GB") {
		t.Errorf("expected iso=GB, got %s", body)
	}
}

func TestCountriesEndpoint(t *testing.T) {
	h := NewHandler("test")
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("GET", "/countries", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "code=1") {
		t.Errorf("expected code=1 in response, got %s", body)
	}
	if !strings.Contains(body, "code=44") {
		t.Errorf("expected code=44 in response, got %s", body)
	}
}

func TestCountryLookupEndpoint(t *testing.T) {
	h := NewHandler("test")
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("GET", "/country/44", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "iso=GB") {
		t.Errorf("expected iso=GB, got %s", w.Body.String())
	}
}

func TestCountryLookupNotFound(t *testing.T) {
	h := NewHandler("test")
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("GET", "/country/999", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestMCPInitialize(t *testing.T) {
	h := NewHandler("test")
	mux := http.NewServeMux()
	h.Register(mux)

	body := `{"jsonrpc":"2.0","id":1,"method":"initialize"}`
	req := httptest.NewRequest("POST", "/mcp", strings.NewReader(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}
	if resp["jsonrpc"] != "2.0" {
		t.Errorf("expected jsonrpc 2.0, got %v", resp["jsonrpc"])
	}
}

func TestMCPToolsList(t *testing.T) {
	h := NewHandler("test")
	mux := http.NewServeMux()
	h.Register(mux)

	body := `{"jsonrpc":"2.0","id":2,"method":"tools/list"}`
	req := httptest.NewRequest("POST", "/mcp", strings.NewReader(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}
	result, ok := resp["result"].(map[string]interface{})
	if !ok {
		t.Fatal("expected result in response")
	}
	tools, ok := result["tools"].([]interface{})
	if !ok {
		t.Fatal("expected tools array in result")
	}
	if len(tools) < 6 {
		t.Errorf("expected at least 6 tools, got %d", len(tools))
	}
}

func TestMCPToolCallParse(t *testing.T) {
	h := NewHandler("test")
	mux := http.NewServeMux()
	h.Register(mux)

	body := `{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"parse","arguments":{"number":"+14155551234"}}}`
	req := httptest.NewRequest("POST", "/mcp", strings.NewReader(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}
	result, ok := resp["result"].(map[string]interface{})
	if !ok {
		t.Fatal("expected result in response")
	}
	content, ok := result["content"].([]interface{})
	if !ok || len(content) == 0 {
		t.Fatal("expected content array in result")
	}
}

func TestMCPToolCallValidate(t *testing.T) {
	h := NewHandler("test")
	mux := http.NewServeMux()
	h.Register(mux)

	body := `{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"validate","arguments":{"number":"+14155551234"}}}`
	req := httptest.NewRequest("POST", "/mcp", strings.NewReader(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}
	if _, ok := resp["result"]; !ok {
		t.Error("expected result in response")
	}
}

func TestMCPToolCallFormat(t *testing.T) {
	h := NewHandler("test")
	mux := http.NewServeMux()
	h.Register(mux)

	body := `{"jsonrpc":"2.0","id":5,"method":"tools/call","params":{"name":"format","arguments":{"number":"+14155551234","format":"e164"}}}`
	req := httptest.NewRequest("POST", "/mcp", strings.NewReader(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}
	if _, ok := resp["result"]; !ok {
		t.Error("expected result in response")
	}
}

func TestMCPToolCallListCountries(t *testing.T) {
	h := NewHandler("test")
	mux := http.NewServeMux()
	h.Register(mux)

	body := `{"jsonrpc":"2.0","id":6,"method":"tools/call","params":{"name":"list_countries","arguments":{}}}`
	req := httptest.NewRequest("POST", "/mcp", strings.NewReader(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}
	if _, ok := resp["result"]; !ok {
		t.Error("expected result in response")
	}
}

func TestMCPToolCallLookupCountry(t *testing.T) {
	h := NewHandler("test")
	mux := http.NewServeMux()
	h.Register(mux)

	body := `{"jsonrpc":"2.0","id":7,"method":"tools/call","params":{"name":"lookup_country","arguments":{"code":"44"}}}`
	req := httptest.NewRequest("POST", "/mcp", strings.NewReader(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}
	if _, ok := resp["result"]; !ok {
		t.Error("expected result in response")
	}
}

func TestMCPMethodNotFound(t *testing.T) {
	h := NewHandler("test")
	mux := http.NewServeMux()
	h.Register(mux)

	body := `{"jsonrpc":"2.0","id":8,"method":"unknown_method"}`
	req := httptest.NewRequest("POST", "/mcp", strings.NewReader(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}
	if _, ok := resp["error"]; !ok {
		t.Error("expected error in response")
	}
}
