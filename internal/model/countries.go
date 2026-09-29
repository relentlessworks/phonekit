package model

import "strings"

// CountryInfo holds basic country info for phone number matching.
type CountryInfo struct {
	Code     string // country calling code
	ISO      string // ISO 3166-1 alpha-2
	Name     string // country name
	MinLen   int    // minimum national number length
	MaxLen   int    // maximum national number length
}

// countryCodes is a lookup table of country calling codes.
// Ordered by code length (longest first) for prefix matching.
var countryCodes = []CountryInfo{
	{"1", "US", "United States", 10, 10},
	{"1", "CA", "Canada", 10, 10},
	{"7", "RU", "Russia", 10, 10},
	{"20", "EG", "Egypt", 9, 10},
	{"27", "ZA", "South Africa", 9, 9},
	{"30", "GR", "Greece", 10, 10},
	{"31", "NL", "Netherlands", 9, 9},
	{"32", "BE", "Belgium", 9, 9},
	{"33", "FR", "France", 9, 9},
	{"34", "ES", "Spain", 9, 9},
	{"36", "HU", "Hungary", 8, 9},
	{"39", "IT", "Italy", 9, 11},
	{"40", "RO", "Romania", 9, 9},
	{"41", "CH", "Switzerland", 9, 9},
	{"43", "AT", "Austria", 10, 11},
	{"44", "GB", "United Kingdom", 9, 10},
	{"45", "DK", "Denmark", 8, 8},
	{"46", "SE", "Sweden", 7, 10},
	{"47", "NO", "Norway", 8, 8},
	{"48", "PL", "Poland", 9, 9},
	{"49", "DE", "Germany", 6, 11},
	{"51", "PE", "Peru", 9, 9},
	{"52", "MX", "Mexico", 10, 10},
	{"53", "CU", "Cuba", 8, 8},
	{"54", "AR", "Argentina", 10, 11},
	{"55", "BR", "Brazil", 10, 11},
	{"56", "CL", "Chile", 9, 9},
	{"57", "CO", "Colombia", 10, 10},
	{"58", "VE", "Venezuela", 10, 10},
	{"60", "MY", "Malaysia", 9, 10},
	{"61", "AU", "Australia", 9, 9},
	{"62", "ID", "Indonesia", 9, 10},
	{"63", "PH", "Philippines", 10, 10},
	{"64", "NZ", "New Zealand", 8, 10},
	{"65", "SG", "Singapore", 8, 8},
	{"66", "TH", "Thailand", 9, 9},
	{"81", "JP", "Japan", 9, 9},
	{"82", "KR", "South Korea", 9, 10},
	{"84", "VN", "Vietnam", 9, 10},
	{"86", "CN", "China", 11, 11},
	{"90", "TR", "Turkey", 10, 10},
	{"91", "IN", "India", 10, 10},
	{"92", "PK", "Pakistan", 10, 10},
	{"93", "AF", "Afghanistan", 9, 9},
	{"94", "LK", "Sri Lanka", 9, 10},
	{"95", "MM", "Myanmar", 8, 10},
	{"98", "IR", "Iran", 10, 10},
	{"212", "MA", "Morocco", 9, 9},
	{"213", "DZ", "Algeria", 9, 9},
	{"216", "TN", "Tunisia", 8, 8},
	{"234", "NG", "Nigeria", 10, 10},
	{"254", "KE", "Kenya", 9, 9},
	{"351", "PT", "Portugal", 9, 9},
	{"352", "LU", "Luxembourg", 8, 8},
	{"353", "IE", "Ireland", 9, 9},
	{"354", "IS", "Iceland", 7, 7},
	{"358", "FI", "Finland", 9, 10},
	{"359", "BG", "Bulgaria", 9, 9},
	{"370", "LT", "Lithuania", 8, 8},
	{"371", "LV", "Latvia", 8, 8},
	{"372", "EE", "Estonia", 7, 8},
	{"374", "AM", "Armenia", 8, 8},
	{"375", "BY", "Belarus", 9, 9},
	{"380", "UA", "Ukraine", 9, 9},
	{"381", "RS", "Serbia", 8, 9},
	{"385", "HR", "Croatia", 9, 9},
	{"386", "SI", "Slovenia", 8, 8},
	{"420", "CZ", "Czech Republic", 9, 9},
	{"421", "SK", "Slovakia", 9, 9},
	{"852", "HK", "Hong Kong", 8, 8},
	{"853", "MO", "Macau", 8, 8},
	{"855", "KH", "Cambodia", 8, 9},
	{"856", "LA", "Laos", 9, 10},
	{"880", "BD", "Bangladesh", 10, 10},
	{"886", "TW", "Taiwan", 9, 9},
	{"971", "AE", "United Arab Emirates", 9, 9},
	{"972", "IL", "Israel", 9, 9},
	{"974", "QA", "Qatar", 8, 8},
	{"977", "NP", "Nepal", 10, 10},
}

// matchCountryCode tries to match a country code from the digit string.
// Returns: countryCode, iso, countryName, remainingDigits, found
func matchCountryCode(digits string) (string, string, string, string, bool) {
	// Try 3-digit codes first, then 2-digit, then 1-digit
	for _, l := range []int{3, 2, 1} {
		if len(digits) < l {
			continue
		}
		prefix := digits[:l]
		rest := digits[l:]
		for _, c := range countryCodes {
			if c.Code == prefix {
				// Check length constraints
				natLen := len(rest)
				if natLen >= c.MinLen && natLen <= c.MaxLen {
					return c.Code, c.ISO, c.Name, rest, true
				}
			}
		}
	}
	return "", "", "", "", false
}

// guessCountryCode tries to guess the country code when no + prefix is present.
// This is less reliable and should be used cautiously.
func guessCountryCode(digits string) (string, string, string, string, bool) {
	// If the number starts with 00 (international prefix), strip it
	if strings.HasPrefix(digits, "00") && len(digits) > 4 {
		stripped := digits[2:]
		if cc, iso, country, rest, ok := matchCountryCode(stripped); ok {
			return cc, iso, country, rest, true
		}
	}

	// If the number starts with 011 (US international prefix), strip it
	if strings.HasPrefix(digits, "011") && len(digits) > 4 {
		stripped := digits[3:]
		if cc, iso, country, rest, ok := matchCountryCode(stripped); ok {
			return cc, iso, country, rest, true
		}
	}

	// Try to match with common patterns
	// US/CA: 11 digits starting with 1
	if len(digits) == 11 && strings.HasPrefix(digits, "1") {
		rest := digits[1:]
		if len(rest) == 10 {
			return "1", "US", "United States", rest, true
		}
	}

	// Try direct country code match
	if cc, iso, country, rest, ok := matchCountryCode(digits); ok {
		return cc, iso, country, rest, true
	}

	return "", "", "", "", false
}

// ListCountries returns all known country codes.
func ListCountries() []CountryInfo {
	return countryCodes
}

// LookupCountry returns country info for a given country code.
func LookupCountry(code string) (CountryInfo, bool) {
	for _, c := range countryCodes {
		if c.Code == code {
			return c, true
		}
	}
	return CountryInfo{}, false
}
