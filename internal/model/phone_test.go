package model

import (
	"strings"
	"testing"
)

func TestParseUSNumber(t *testing.T) {
	pn, err := Parse("+1 (415) 555-1234")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pn.CountryCode != "1" {
		t.Errorf("expected country code 1, got %s", pn.CountryCode)
	}
	if pn.CountryISO != "US" {
		t.Errorf("expected ISO US, got %s", pn.CountryISO)
	}
	if pn.NationalNum != "4155551234" {
		t.Errorf("expected national 4155551234, got %s", pn.NationalNum)
	}
	if pn.AreaCode != "415" {
		t.Errorf("expected area 415, got %s", pn.AreaCode)
	}
	if pn.E164 != "+14155551234" {
		t.Errorf("expected E164 +14155551234, got %s", pn.E164)
	}
	if !pn.Valid {
		t.Error("expected valid=true")
	}
}

func TestParseUKNumber(t *testing.T) {
	pn, err := Parse("+44 20 7946 0958")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pn.CountryCode != "44" {
		t.Errorf("expected country code 44, got %s", pn.CountryCode)
	}
	if pn.CountryISO != "GB" {
		t.Errorf("expected ISO GB, got %s", pn.CountryISO)
	}
	if pn.NationalNum != "2079460958" {
		t.Errorf("expected national 2079460958, got %s", pn.NationalNum)
	}
	if pn.E164 != "+442079460958" {
		t.Errorf("expected E164 +442079460958, got %s", pn.E164)
	}
}

func TestParseGermanNumber(t *testing.T) {
	pn, err := Parse("+49 30 1234567")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pn.CountryCode != "49" {
		t.Errorf("expected country code 49, got %s", pn.CountryCode)
	}
	if pn.CountryISO != "DE" {
		t.Errorf("expected ISO DE, got %s", pn.CountryISO)
	}
	if pn.NationalNum != "301234567" {
		t.Errorf("expected national 301234567, got %s", pn.NationalNum)
	}
}

func TestParseFrenchNumber(t *testing.T) {
	pn, err := Parse("+33 1 23 45 67 89")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pn.CountryCode != "33" {
		t.Errorf("expected country code 33, got %s", pn.CountryCode)
	}
	if pn.CountryISO != "FR" {
		t.Errorf("expected ISO FR, got %s", pn.CountryISO)
	}
	if pn.NationalNum != "123456789" {
		t.Errorf("expected national 123456789, got %s", pn.NationalNum)
	}
}

func TestParseJapaneseNumber(t *testing.T) {
	pn, err := Parse("+81 3 1234 5678")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pn.CountryCode != "81" {
		t.Errorf("expected country code 81, got %s", pn.CountryCode)
	}
	if pn.CountryISO != "JP" {
		t.Errorf("expected ISO JP, got %s", pn.CountryISO)
	}
}

func TestParseUSTollFree(t *testing.T) {
	pn, err := Parse("+1 800 555 1234")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pn.LineType != "tollfree" {
		t.Errorf("expected line type tollfree, got %s", pn.LineType)
	}
}

func TestParseUSPremium(t *testing.T) {
	pn, err := Parse("+1 900 555 1234")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pn.LineType != "premium" {
		t.Errorf("expected line type premium, got %s", pn.LineType)
	}
}

func TestParseUKMobile(t *testing.T) {
	pn, err := Parse("+44 7700 900123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pn.LineType != "mobile" {
		t.Errorf("expected line type mobile, got %s", pn.LineType)
	}
}

func TestParseGermanMobile(t *testing.T) {
	pn, err := Parse("+49 151 23456789")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pn.LineType != "mobile" {
		t.Errorf("expected line type mobile, got %s", pn.LineType)
	}
}

func TestParseFrenchMobile(t *testing.T) {
	pn, err := Parse("+33 6 12 34 56 78")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pn.LineType != "mobile" {
		t.Errorf("expected line type mobile, got %s", pn.LineType)
	}
}

func TestParseIndianMobile(t *testing.T) {
	pn, err := Parse("+91 98765 43210")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pn.CountryCode != "91" {
		t.Errorf("expected country code 91, got %s", pn.CountryCode)
	}
	if pn.LineType != "mobile" {
		t.Errorf("expected line type mobile, got %s", pn.LineType)
	}
}

func TestParseWithExtension(t *testing.T) {
	pn, err := Parse("+1 415 555 1234 ext 99")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pn.Extension != "99" {
		t.Errorf("expected extension 99, got %s", pn.Extension)
	}
	if pn.RFC3966 != "tel:+14155551234;ext=99" {
		t.Errorf("expected RFC3966 tel:+1415551234;ext=99, got %s", pn.RFC3966)
	}
}

func TestParseNoPlusPrefix(t *testing.T) {
	pn, err := Parse("14155551234")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pn.CountryCode != "1" {
		t.Errorf("expected country code 1, got %s", pn.CountryCode)
	}
}

func TestParseInternationalPrefix00(t *testing.T) {
	pn, err := Parse("00442079460958")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pn.CountryCode != "44" {
		t.Errorf("expected country code 44, got %s", pn.CountryCode)
	}
}

func TestParseEmpty(t *testing.T) {
	_, err := Parse("")
	if err != ErrEmptyNumber {
		t.Errorf("expected ErrEmptyNumber, got %v", err)
	}
}

func TestParseTooShort(t *testing.T) {
	_, err := Parse("+1")
	if err != ErrTooShort {
		t.Errorf("expected ErrTooShort, got %v", err)
	}
}

func TestParseTooLong(t *testing.T) {
	_, err := Parse("+1 234 567 890 123 456 7")
	if err != ErrTooLong {
		t.Errorf("expected ErrTooLong, got %v", err)
	}
}

func TestValidate(t *testing.T) {
	if !Validate("+14155551234") {
		t.Error("expected +14155551234 to be valid")
	}
	if Validate("") {
		t.Error("expected empty to be invalid")
	}
}

func TestFormatE164(t *testing.T) {
	result, err := Format("+1 (415) 555-1234", "e164")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "+14155551234" {
		t.Errorf("expected +14155551234, got %s", result)
	}
}

func TestFormatInternational(t *testing.T) {
	result, err := Format("+14155551234", "international")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "+1 41 55 55 12 34" {
		t.Errorf("expected +1 41 55 55 12 34, got %s", result)
	}
}

func TestFormatNational(t *testing.T) {
	result, err := Format("+14155551234", "national")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "(415) 555-1234" {
		t.Errorf("expected (415) 555-1234, got %s", result)
	}
}

func TestFormatRFC3966(t *testing.T) {
	result, err := Format("+14155551234", "rfc3966")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "tel:+14155551234" {
		t.Errorf("expected tel:+14155551234, got %s", result)
	}
}

func TestFormatInvalid(t *testing.T) {
	_, err := Format("+14155551234", "invalid")
	if err != ErrInvalidFormat {
		t.Errorf("expected ErrInvalidFormat, got %v", err)
	}
}

func TestListCountries(t *testing.T) {
	countries := ListCountries()
	if len(countries) == 0 {
		t.Error("expected non-empty country list")
	}
}

func TestLookupCountry(t *testing.T) {
	info, ok := LookupCountry("44")
	if !ok {
		t.Error("expected to find country code 44")
	}
	if info.ISO != "GB" {
		t.Errorf("expected ISO GB, got %s", info.ISO)
	}
}

func TestLookupCountryNotFound(t *testing.T) {
	_, ok := LookupCountry("999")
	if ok {
		t.Error("expected not to find country code 999")
	}
}

func TestFormatUKNational(t *testing.T) {
	result, err := Format("+442079460958", "national")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "020 7946 0958" {
		t.Errorf("expected 020 7946 0958, got %s", result)
	}
}

func TestFormatUKMobileNational(t *testing.T) {
	result, err := Format("+447700900123", "national")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// UK mobile: 07XXX XXXXXX
	if !strings.HasPrefix(result, "077") {
		t.Errorf("expected 077 prefix, got %s", result)
	}
}
