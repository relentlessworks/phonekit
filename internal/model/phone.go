package model

import (
	"regexp"
	"strings"
)

// PhoneNumber represents a parsed phone number.
type PhoneNumber struct {
	Raw          string // original input
	CountryCode  string // e.g. "1", "44"
	NationalNum  string // significant national number (digits only)
	AreaCode     string // area/region code if extractable
	Subscriber   string // subscriber number
	Extension    string // extension if present
	E164         string // E.164 format: +CCNNNNNNN
	International string // international format: +CC NNN NNN NNNN
	National     string // national format with typical spacing
	RFC3966      string // RFC 3966 tel: URI format
	Country      string // country name (best-effort)
	CountryISO   string // ISO 3166-1 alpha-2 code
	LineType     string // mobile, fixed, tollfree, premium, unknown
	Valid        bool   // basic validity
}

var (
	// stripNonDigits removes everything except digits
	nonDigitRe = regexp.MustCompile(`[^\d]`)
	// detect extension: ext=123, x123, ;ext=123
	extRe = regexp.MustCompile(`(?i)(?:ext|x|;ext=)\s*(\d+)`)
)

// Parse parses a phone number string and returns structured info.
// It does a best-effort parse without a full phone library.
// It handles common formats: +1 (555) 123-4567, 44 20 7946 0958, etc.
func Parse(input string) (*PhoneNumber, error) {
	if input == "" {
		return nil, ErrEmptyNumber
	}

	raw := strings.TrimSpace(input)

	// Extract extension before stripping
	ext := ""
	if m := extRe.FindStringSubmatch(raw); len(m) > 1 {
		ext = m[1]
	}

	// Strip extension from the raw for digit extraction
	cleaned := extRe.ReplaceAllString(raw, "")

	// Check for leading + to determine if country code is present
	hasPlus := strings.HasPrefix(strings.TrimSpace(cleaned), "+")

	// Extract all digits
	digits := nonDigitRe.ReplaceAllString(cleaned, "")

	if len(digits) < 4 {
		return nil, ErrTooShort
	}
	if len(digits) > 15 {
		return nil, ErrTooLong
	}

	pn := &PhoneNumber{
		Raw:       raw,
		Extension: ext,
	}

	// Determine country code
	remaining := digits
	if hasPlus {
		// Try to match country codes (1-3 digits)
		cc, iso, country, rest, ok := matchCountryCode(digits)
		if ok {
			pn.CountryCode = cc
			pn.CountryISO = iso
			pn.Country = country
			remaining = rest
		} else {
			// Can't determine country code, use as-is
			pn.NationalNum = digits
		}
	} else {
		// No + prefix — assume it's a national number
		// Try to guess country code for common patterns
		cc, iso, country, rest, ok := guessCountryCode(digits)
		if ok {
			pn.CountryCode = cc
			pn.CountryISO = iso
			pn.Country = country
			remaining = rest
		} else {
			pn.NationalNum = digits
		}
	}

	if pn.NationalNum == "" {
		pn.NationalNum = remaining
	}

	// Split national number into area code and subscriber
	splitNational(pn)

	// Determine line type
	pn.LineType = detectLineType(pn)

	// Generate formatted versions
	pn.E164 = formatE164(pn)
	pn.International = formatInternational(pn)
	pn.National = formatNational(pn)
	pn.RFC3966 = formatRFC3966(pn)

	// Basic validation
	pn.Valid = len(pn.NationalNum) >= 4 && len(pn.NationalNum) <= 14

	return pn, nil
}

// splitNational attempts to split the national number into area code and subscriber.
func splitNational(pn *PhoneNumber) {
	n := pn.NationalNum
	switch pn.CountryISO {
	case "US", "CA":
		if len(n) >= 10 {
			pn.AreaCode = n[:3]
			pn.Subscriber = n[3:]
		} else if len(n) == 7 {
			pn.Subscriber = n
		}
	case "GB":
		if len(n) >= 10 && strings.HasPrefix(n, "0") {
			n = n[1:]
		}
		if len(n) >= 9 {
			// UK area codes vary: 2 digits (London 20), 3, 4, 5 digits
			if strings.HasPrefix(n, "1") || strings.HasPrefix(n, "2") {
				pn.AreaCode = n[:2]
				pn.Subscriber = n[2:]
			} else {
				pn.AreaCode = n[:3]
				pn.Subscriber = n[3:]
			}
		}
	case "DE":
		if len(n) >= 6 {
			// German area codes are 2-5 digits
			if strings.HasPrefix(n, "0") {
				n = n[1:]
			}
			pn.AreaCode = n[:3]
			pn.Subscriber = n[3:]
		}
	case "FR":
		if len(n) >= 9 {
			pn.AreaCode = n[:1]
			pn.Subscriber = n[1:]
		}
	case "AU":
		if len(n) >= 9 {
			pn.AreaCode = n[:1]
			pn.Subscriber = n[1:]
		}
	case "JP":
		if len(n) >= 9 {
			pn.AreaCode = n[:2]
			pn.Subscriber = n[2:]
		}
	case "IN":
		if len(n) >= 10 {
			pn.AreaCode = n[:4]
			pn.Subscriber = n[4:]
		}
	default:
		if len(n) >= 6 {
			pn.AreaCode = n[:3]
			pn.Subscriber = n[3:]
		} else {
			pn.Subscriber = n
		}
	}
}

// detectLineType attempts to determine the line type from the number pattern.
func detectLineType(pn *PhoneNumber) string {
	n := pn.NationalNum
	cc := pn.CountryCode

	// Toll-free numbers
	if cc == "1" && len(n) == 10 {
		ac := n[:3]
		if ac == "800" || ac == "888" || ac == "877" || ac == "866" || ac == "855" || ac == "844" || ac == "833" {
			return "tollfree"
		}
		// Premium rate
		if ac == "900" {
			return "premium"
		}
		// Mobile: US mobile numbers don't have distinct prefixes, but we can check area codes
		// This is a simplification
		return "fixed_or_mobile"
	}

	// UK mobile numbers start with 7
	if cc == "44" && strings.HasPrefix(n, "7") {
		return "mobile"
	}
	if cc == "44" && strings.HasPrefix(n, "80") {
		return "tollfree"
	}
	if cc == "44" && strings.HasPrefix(n, "9") {
		return "premium"
	}

	// German mobile: 15, 16, 17 prefixes
	if cc == "49" && (strings.HasPrefix(n, "15") || strings.HasPrefix(n, "16") || strings.HasPrefix(n, "17")) {
		return "mobile"
	}

	// French mobile: 6 or 7
	if cc == "33" && (strings.HasPrefix(n, "6") || strings.HasPrefix(n, "7")) {
		return "mobile"
	}

	// Australian mobile: 4
	if cc == "61" && strings.HasPrefix(n, "4") {
		return "mobile"
	}

	// Japanese mobile: 70, 80, 90
	if cc == "81" && (strings.HasPrefix(n, "70") || strings.HasPrefix(n, "80") || strings.HasPrefix(n, "90")) {
		return "mobile"
	}

	// Indian mobile: 6, 7, 8, 9
	if cc == "91" && len(n) == 10 && (strings.HasPrefix(n, "6") || strings.HasPrefix(n, "7") || strings.HasPrefix(n, "8") || strings.HasPrefix(n, "9")) {
		return "mobile"
	}

	return "unknown"
}

// formatE164 returns +CC followed by national number
func formatE164(pn *PhoneNumber) string {
	if pn.CountryCode != "" {
		return "+" + pn.CountryCode + pn.NationalNum
	}
	return "+" + pn.NationalNum
}

// formatInternational returns +CC with grouped national number
func formatInternational(pn *PhoneNumber) string {
	if pn.CountryCode == "" {
		return "+" + groupDigits(pn.NationalNum)
	}
	return "+" + pn.CountryCode + " " + groupDigits(pn.NationalNum)
}

// formatNational returns the national format
func formatNational(pn *PhoneNumber) string {
	n := pn.NationalNum
	switch pn.CountryISO {
	case "US", "CA":
		if len(n) == 10 {
			return "(" + n[:3] + ") " + n[3:6] + "-" + n[6:]
		}
		if len(n) == 7 {
			return n[:3] + "-" + n[3:]
		}
	case "GB":
		// UK national format: 0XX XXXX XXXX or 0XXX XXX XXXX
		if len(n) >= 9 {
			if strings.HasPrefix(n, "1") || strings.HasPrefix(n, "2") {
				// 2-digit area code (London, etc.)
				ac := n[:2]
				sub := n[2:]
				if len(sub) >= 4 {
					return "0" + ac + " " + sub[:4] + " " + sub[4:]
				}
				return "0" + ac + " " + sub
			}
			// 3-digit area code
			ac := n[:3]
			sub := n[3:]
			if len(sub) >= 3 {
				return "0" + ac + " " + sub[:3] + " " + sub[3:]
			}
			return "0" + ac + " " + sub
		}
	case "DE":
		if len(n) >= 6 {
			return "0" + n[:3] + " " + n[3:]
		}
	case "FR":
		if len(n) == 9 {
			return "0" + n[:1] + " " + n[1:3] + " " + n[3:5] + " " + n[5:7] + " " + n[7:9]
		}
	case "AU":
		if len(n) == 9 {
			return "0" + n[:1] + " " + n[1:4] + " " + n[4:]
		}
	case "JP":
		if len(n) >= 9 {
			return "0" + n[:2] + "-" + n[2:]
		}
	case "IN":
		if len(n) == 10 {
			return n[:5] + " " + n[5:]
		}
	}
	return groupDigits(n)
}

// formatRFC3966 returns the tel: URI format
func formatRFC3966(pn *PhoneNumber) string {
	s := "tel:" + formatE164(pn)
	if pn.Extension != "" {
		s += ";ext=" + pn.Extension
	}
	return s
}

// groupDigits groups digits in pairs from the right
func groupDigits(s string) string {
	if len(s) <= 2 {
		return s
	}
	// Group from the right in pairs
	var result []byte
	remainder := len(s) % 2
	if remainder == 1 {
		result = append(result, s[0])
		s = s[1:]
	}
	for i := 0; i < len(s); i += 2 {
		if i > 0 || remainder == 1 {
			result = append(result, ' ')
		}
		if i+2 <= len(s) {
			result = append(result, s[i], s[i+1])
		} else {
			result = append(result, s[i])
		}
	}
	return string(result)
}

// Validate checks if a phone number is valid (basic check).
func Validate(input string) bool {
	pn, err := Parse(input)
	if err != nil {
		return false
	}
	return pn.Valid
}

// Format converts a phone number to the specified format.
func Format(input, format string) (string, error) {
	pn, err := Parse(input)
	if err != nil {
		return "", err
	}
	switch strings.ToLower(format) {
	case "e164", "e.164":
		return pn.E164, nil
	case "international", "intl":
		return pn.International, nil
	case "national":
		return pn.National, nil
	case "rfc3966", "tel", "uri":
		return pn.RFC3966, nil
	default:
		return "", ErrInvalidFormat
	}
}

// GetInfo returns all available info about a phone number.
func GetInfo(input string) (*PhoneNumber, error) {
	return Parse(input)
}
