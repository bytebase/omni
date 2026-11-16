package catalog

import (
	"fmt"
	"strconv"
	"strings"
)

func normalizeDataTypeText(text string) string {
	if text == "" {
		return ""
	}

	normalized := upperOutsideQuotes(text)

	switch {
	case strings.HasPrefix(normalized, "VARCHAR2"):
		if strings.HasPrefix(normalized, "VARCHAR2(") && strings.HasSuffix(normalized, ")") {
			start := strings.Index(normalized, "(") + 1
			end := strings.LastIndex(normalized, ")")
			if start < end {
				sizeStr := normalized[start:end]
				if size := parseInt(sizeStr); size > 0 {
					return fmt.Sprintf("VARCHAR2(%d BYTE)", size)
				}
			}
		}
		return normalized
	case strings.HasPrefix(normalized, "CHAR("):
		if strings.HasSuffix(normalized, ")") {
			start := strings.Index(normalized, "(") + 1
			end := strings.LastIndex(normalized, ")")
			if start < end {
				sizeStr := normalized[start:end]
				if size := parseInt(sizeStr); size > 0 {
					return fmt.Sprintf("CHAR(%d BYTE)", size)
				}
			}
		}
		return normalized
	case strings.HasPrefix(normalized, "TIMESTAMP"):
		if strings.Contains(normalized, "WITH LOCAL TIME ZONE") || strings.Contains(normalized, "WITHLOCALTIMEZONE") {
			if normalized == "TIMESTAMP WITH LOCAL TIME ZONE" {
				return "TIMESTAMP(6) WITH LOCAL TIME ZONE"
			}
			if strings.Contains(normalized, "(") && strings.Contains(normalized, ")") {
				return strings.ReplaceAll(normalized, "WITHLOCALTIMEZONE", " WITH LOCAL TIME ZONE")
			}
			return "TIMESTAMP(6) WITH LOCAL TIME ZONE"
		}
		if strings.Contains(normalized, "WITH TIME ZONE") || strings.Contains(normalized, "WITHTIMEZONE") {
			if normalized == "TIMESTAMP WITH TIME ZONE" {
				return "TIMESTAMP(6) WITH TIME ZONE"
			}
			if strings.Contains(normalized, "(") && strings.Contains(normalized, ")") {
				return strings.ReplaceAll(normalized, "WITHTIMEZONE", " WITH TIME ZONE")
			}
			return "TIMESTAMP(6) WITH TIME ZONE"
		}
		if normalized == "TIMESTAMP" {
			return "TIMESTAMP(6)"
		}
		return normalized
	case strings.HasPrefix(normalized, "INTERVAL"):
		if strings.Contains(normalized, "YEAR TO MONTH") || strings.Contains(normalized, "YEARTOMONTH") {
			if normalized == "INTERVAL YEAR TO MONTH" {
				return "INTERVAL YEAR(2) TO MONTH"
			}
			if strings.Contains(normalized, "(") && strings.Contains(normalized, ")") {
				result := strings.ReplaceAll(normalized, "YEARTOMONTH", " YEAR TO MONTH")
				return strings.ReplaceAll(result, "INTERVALYEAR", "INTERVAL YEAR")
			}
			return "INTERVAL YEAR(2) TO MONTH"
		}
		if strings.Contains(normalized, "DAY TO SECOND") || strings.Contains(normalized, "DAYTOSECOND") {
			if normalized == "INTERVAL DAY TO SECOND" {
				return "INTERVAL DAY(2) TO SECOND(6)"
			}
			if strings.Contains(normalized, "(") && strings.Contains(normalized, ")") {
				result := strings.ReplaceAll(normalized, "DAYTOSECOND", " DAY TO SECOND")
				return strings.ReplaceAll(result, "INTERVALDAY", "INTERVAL DAY")
			}
			return "INTERVAL DAY(2) TO SECOND(6)"
		}
		return normalized
	case strings.HasPrefix(normalized, "LONGRAW"):
		return "LONG RAW"
	case strings.HasPrefix(normalized, "DOUBLE"):
		return "DOUBLE PRECISION"
	case normalized == "UROWID":
		return "UROWID(4000)"
	case normalized == "FLOAT":
		return "FLOAT(126)"
	case strings.HasPrefix(normalized, "NVARCHAR2"):
		// National character types declare character counts. Sync reads
		// CHAR_LENGTH, so doubling here would make SQL-derived metadata and a
		// synced snapshot disagree.
		return normalized
	case strings.HasPrefix(normalized, "NCHAR"):
		return normalized
	default:
		return normalized
	}
}

// upperOutsideQuotes folds unquoted text the way Oracle folds identifiers,
// leaving a quoted identifier's spelling alone.
func upperOutsideQuotes(text string) string {
	var out strings.Builder
	var quote bool
	for i := 0; i < len(text); i++ {
		c := text[i]
		if quote {
			out.WriteByte(c)
			if c == '"' {
				if i+1 < len(text) && text[i+1] == '"' {
					out.WriteByte(text[i+1])
					i++
				} else {
					quote = false
				}
			}
			continue
		}
		if c == '"' {
			quote = true
			out.WriteByte(c)
			continue
		}
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		out.WriteByte(c)
	}
	return out.String()
}

func parseInt(s string) int {
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	return 0
}
