package parser

import (
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"
	_ "time/tzdata"
)

// splitMultipartIdent splits a multipart identifier string on '.' outside
// backtick-quoted segments (preserving escaped backticks).
func splitMultipartIdent(s string) []string {
	var parts []string
	start := 0
	inBacktick := false
	for i := 0; i < len(s); i++ {
		if s[i] == '`' {
			if inBacktick && i+1 < len(s) && s[i+1] == '`' {
				i++
				continue
			}
			inBacktick = !inBacktick
		} else if s[i] == '.' && !inBacktick {
			parts = append(parts, s[start:i])
			start = i + 1
		}
	}
	parts = append(parts, s[start:])
	return parts
}

// isQualifiedIdent reports whether a parsed multipart identifier string
// contains a qualification dot outside backtick-quoted segments.
func isQualifiedIdent(s string) bool {
	return len(splitMultipartIdent(s)) > 1
}

// exactIdentName unquotes a backtick-quoted identifier (unescaping “) while
// preserving case, matching Spark v4.0.5's case-sensitive CTE and property checks.
func exactIdentName(s string) string {
	if len(s) >= 2 && s[0] == '`' && s[len(s)-1] == '`' {
		return strings.ReplaceAll(s[1:len(s)-1], "``", "`")
	}
	return s
}

// unquoteStringLit decodes a Spark single-quoted, double-quoted, or raw string
// literal (R'...' / R"...") into its logical value, matching Spark v4.0.5
// ParserUtils.unescapeSQLString (including \uXXXX, \UXXXXXXXX, and 3-digit
// octal \[0-3][0-7][0-7] escapes).
func unquoteStringLit(s string) string {
	if len(s) >= 3 && (s[0] == 'R' || s[0] == 'r') && (s[1] == '\'' || s[1] == '"') && s[len(s)-1] == s[1] {
		return s[2 : len(s)-1]
	}
	if len(s) < 2 {
		return s
	}
	q := s[0]
	if (q != '\'' && q != '"') || s[len(s)-1] != q {
		return s
	}
	inner := s[1 : len(s)-1]
	if !strings.ContainsRune(inner, '\\') {
		return inner
	}
	var b strings.Builder
	b.Grow(len(inner))
	for i := 0; i < len(inner); i++ {
		if inner[i] == '\\' && i+1 < len(inner) {
			i++
			switch inner[i] {
			case '0':
				if i+2 < len(inner) && isOctalDigit(inner[i+1]) && isOctalDigit(inner[i+2]) {
					val := (inner[i]-'0')*64 + (inner[i+1]-'0')*8 + (inner[i+2] - '0')
					b.WriteByte(val)
					i += 2
				} else {
					b.WriteByte(0)
				}
			case '1', '2', '3':
				if i+2 < len(inner) && isOctalDigit(inner[i+1]) && isOctalDigit(inner[i+2]) {
					val := (inner[i]-'0')*64 + (inner[i+1]-'0')*8 + (inner[i+2] - '0')
					b.WriteByte(val)
					i += 2
				} else {
					b.WriteByte(inner[i])
				}
			case 'u':
				if i+4 < len(inner) {
					if v, err := strconv.ParseUint(inner[i+1:i+5], 16, 32); err == nil {
						b.WriteRune(rune(v))
						i += 4
						continue
					}
				}
				b.WriteByte('u')
			case 'U':
				if i+8 < len(inner) {
					if v, err := strconv.ParseUint(inner[i+1:i+9], 16, 32); err == nil && v <= 0x10FFFF {
						b.WriteRune(rune(v))
						i += 8
						continue
					}
				}
				b.WriteByte('U')
			case 'n':
				b.WriteByte('\n')
			case 'r':
				b.WriteByte('\r')
			case 't':
				b.WriteByte('\t')
			case 'b':
				b.WriteByte('\b')
			case 'f':
				b.WriteByte('\f')
			case 'Z':
				b.WriteByte(0x1A)
			case '%', '_':
				b.WriteByte('\\')
				b.WriteByte(inner[i])
			default:
				b.WriteByte(inner[i])
			}
			continue
		}
		b.WriteByte(inner[i])
	}
	return b.String()
}

func isOctalDigit(c byte) bool {
	return c >= '0' && c <= '7'
}

// decodePropertyKey returns the decoded key name for a property key token or
// multipart identifier, matching Spark v4.0.5 AstBuilder.visitPropertyKeyValues.
func decodePropertyKey(raw string, isString bool) string {
	if isString {
		return unquoteStringLit(raw)
	}
	return exactIdentName(raw)
}

// validateSQLState checks Spark v4.0.5 AstBuilder.assertSqlState:
// exactly 5 ASCII alphanumeric characters, not starting with "00", "01", or "XX".
func validateSQLState(val string) bool {
	if len(val) != 5 {
		return false
	}
	for i := 0; i < 5; i++ {
		c := val[i]
		isAlphaNum := (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')
		if !isAlphaNum {
			return false
		}
	}
	prefix := val[:2]
	if prefix == "00" || prefix == "01" || prefix == "XX" {
		return false
	}
	return true
}

// isValidHiveSerDeFormat reports whether format matches one of Spark v4.0.5's
// HiveSerDe.serdeMap keys (used by CREATE TABLE LIKE ... STORED AS and
// INSERT OVERWRITE DIRECTORY ... STORED AS).
func isValidHiveSerDeFormat(name string) bool {
	switch strings.ToLower(exactIdentName(name)) {
	case "rcfile", "textfile", "orc", "parquet", "avro", "sequencefile":
		return true
	}
	return false
}

// isValidFunctionResourceType reports whether name is JAR, FILE, or ARCHIVE
// (case-insensitively), matching Spark v4.0.5 SparkSqlParser.visitCreateFunction.
func isValidFunctionResourceType(name string) bool {
	switch strings.ToUpper(exactIdentName(name)) {
	case "JAR", "FILE", "ARCHIVE":
		return true
	}
	return false
}

// isValidManageResourceType reports whether name is a valid resource type for
// Spark's ADD / LIST resource commands.
func isValidManageResourceType(name string) bool {
	switch strings.ToUpper(exactIdentName(name)) {
	case "JAR", "JARS", "FILE", "FILES", "ARCHIVE", "ARCHIVES":
		return true
	}
	return false
}

// isValidShowFunctionsScope reports whether name is USER, SYSTEM, or ALL
// (case-insensitively), matching Spark v4.0.5 visitShowFunctions.
func isValidShowFunctionsScope(name string) bool {
	switch strings.ToUpper(exactIdentName(name)) {
	case "USER", "SYSTEM", "ALL":
		return true
	}
	return false
}

// sparkICULocales holds the exact 143 ICU locale collation bases from Apache Spark
// 4.0.5's CollationFactory.getICULocaleNames() (stored uppercased for matching).
var sparkICULocales = map[string]struct{}{
	"UNICODE": {}, "AF": {}, "AM": {}, "AR": {}, "AR_SAU": {}, "AS": {}, "AZ": {},
	"BE": {}, "BG": {}, "BN": {}, "BO": {}, "BR": {}, "BS": {}, "BS_CYRL": {},
	"CA": {}, "CEB": {}, "CHR": {}, "CS": {}, "CY": {}, "DA": {}, "DE": {},
	"DE_AUT": {}, "DSB": {}, "DZ": {}, "EE": {}, "EL": {}, "EN": {}, "EN_USA": {},
	"EO": {}, "ES": {}, "ET": {}, "FA": {}, "FA_AFG": {}, "FF": {}, "FF_ADLM": {},
	"FI": {}, "FIL": {}, "FO": {}, "FR": {}, "FR_CAN": {}, "FY": {}, "GA": {},
	"GL": {}, "GU": {}, "HA": {}, "HAW": {}, "HE": {}, "HE_ISR": {}, "HI": {},
	"HR": {}, "HSB": {}, "HU": {}, "HY": {}, "ID": {}, "ID_IDN": {}, "IG": {},
	"IS": {}, "IT": {}, "JA": {}, "KA": {}, "KK": {}, "KL": {}, "KM": {},
	"KN": {}, "KO": {}, "KOK": {}, "KU": {}, "KY": {}, "LB": {}, "LIJ": {},
	"LKT": {}, "LN": {}, "LO": {}, "LT": {}, "LV": {}, "MK": {}, "ML": {},
	"MN": {}, "MR": {}, "MS": {}, "MT": {}, "MY": {}, "NB": {}, "NB_NOR": {},
	"NE": {}, "NL": {}, "NN": {}, "NO": {}, "NSO": {}, "OM": {}, "OR": {},
	"PA": {}, "PA_GURU": {}, "PA_GURU_IND": {}, "PL": {}, "PS": {}, "PT": {},
	"RO": {}, "RU": {}, "SA": {}, "SE": {}, "SI": {}, "SK": {}, "SL": {},
	"SMN": {}, "SQ": {}, "SR": {}, "SR_CYRL": {}, "SR_CYRL_BIH": {},
	"SR_CYRL_MNE": {}, "SR_CYRL_SRB": {}, "SR_LATN": {}, "SR_LATN_BIH": {},
	"SR_LATN_SRB": {}, "ST": {}, "SV": {}, "SW": {}, "TA": {}, "TE": {},
	"TH": {}, "TK": {}, "TN": {}, "TO": {}, "TR": {}, "UG": {}, "UK": {},
	"UR": {}, "UZ": {}, "VI": {}, "WAE": {}, "WO": {}, "XH": {}, "YI": {},
	"YO": {}, "ZH": {}, "ZH_HANS": {}, "ZH_HANS_CHN": {}, "ZH_HANS_SGP": {},
	"ZH_HANT": {}, "ZH_HANT_HKG": {}, "ZH_HANT_MAC": {}, "ZH_HANT_TWN": {},
	"ZU": {},
}

// isTrimSparkCollation reports whether a collation name uses a TRIM/RTRIM
// specifier, which Spark v4.0.5 rejects at parse time in both expression
// COLLATE clauses and data type COLLATE clauses (UNSUPPORTED_FEATURE.TRIM_COLLATION).
func isTrimSparkCollation(raw string) bool {
	parts := splitMultipartIdent(raw)
	last := strings.ToUpper(exactIdentName(parts[len(parts)-1]))
	return last == "TRIM" || last == "RTRIM" ||
		strings.HasSuffix(last, "_TRIM") || strings.HasSuffix(last, "_RTRIM")
}

// isValidSparkCollation validates a collation name (optionally qualified with
// SYSTEM.BUILTIN.) against Apache Spark 4.0.5's CollationFactory, rejecting
// unknown names and unsupported TRIM/RTRIM collations.
func isValidSparkCollation(raw string) bool {
	parts := splitMultipartIdent(raw)
	for i := range parts {
		parts[i] = strings.ToUpper(exactIdentName(parts[i]))
	}
	var name string
	switch len(parts) {
	case 1:
		name = parts[0]
	case 2:
		if parts[0] != "BUILTIN" {
			return false
		}
		name = parts[1]
	case 3:
		if parts[0] != "SYSTEM" || parts[1] != "BUILTIN" {
			return false
		}
		name = parts[2]
	default:
		return false
	}
	if name == "UTF8_BINARY" || name == "UTF8_LCASE" {
		return true
	}
	// Strip optional _CI/_CS and _AI/_AS modifiers (at most one case and one accent specifier).
	hasCase, hasAccent := false, false
	for {
		switch {
		case !hasCase && (strings.HasSuffix(name, "_CI") || strings.HasSuffix(name, "_CS")):
			hasCase = true
			name = name[:len(name)-3]
		case !hasAccent && (strings.HasSuffix(name, "_AI") || strings.HasSuffix(name, "_AS")):
			hasAccent = true
			name = name[:len(name)-3]
		default:
			_, ok := sparkICULocales[name]
			return ok
		}
	}
}

// validateNumericLiteral checks Spark v4.0.5's parse-time range and precision
// bounds for numeric literals (Y/S/L/F/D/BD/exponent/large integer).
func validateNumericLiteral(kind TokenKind, text string, negative bool) bool {
	switch kind {
	case tokTinyInt:
		digits := text[:len(text)-1]
		if negative {
			digits = "-" + digits
		}
		v, err := strconv.ParseInt(digits, 10, 8)
		return err == nil && v >= math.MinInt8 && v <= math.MaxInt8
	case tokSmallInt:
		digits := text[:len(text)-1]
		if negative {
			digits = "-" + digits
		}
		v, err := strconv.ParseInt(digits, 10, 16)
		return err == nil && v >= math.MinInt16 && v <= math.MaxInt16
	case tokBigInt:
		digits := text[:len(text)-1]
		if negative {
			digits = "-" + digits
		}
		_, err := strconv.ParseInt(digits, 10, 64)
		return err == nil
	case tokFloat:
		return decimalMagnitudeWithinBound(text[:len(text)-1], sparkMaxFloatRat)
	case tokDouble:
		return decimalMagnitudeWithinBound(text[:len(text)-1], sparkMaxDoubleRat)
	case tokExponent:
		return decimalMagnitudeWithinBound(text, sparkMaxDoubleRat)
	case tokBigDecimal:
		num := text[:len(text)-2]
		return decimalPrecisionWithinLimit(num, 38)
	case tokInt:
		trimmed := strings.TrimLeft(text, "0")
		return len(trimmed) <= 38
	case tokDecimal:
		return decimalPrecisionWithinLimit(text, 38)
	}
	return true
}

// Spark v4.0.5 AstBuilder.numericLiteral compares the exact decimal literal
// against BigDecimal(Float.MaxValue.toDouble) = 3.4028234663852886E38 and
// BigDecimal(Double.MaxValue) = 1.7976931348623157E308 before float conversion.
var (
	sparkMaxFloatRat, _  = new(big.Rat).SetString("3.4028234663852886e38")
	sparkMaxDoubleRat, _ = new(big.Rat).SetString("1.7976931348623157e308")
)

func decimalMagnitudeWithinBound(s string, maxBound *big.Rat) bool {
	mantissa := s
	exp := 0
	if idx := strings.IndexAny(s, "eE"); idx >= 0 {
		mantissa = s[:idx]
		e, err := strconv.Atoi(s[idx+1:])
		if err != nil {
			// Overflowing exponent: only valid if mantissa is zero.
			return isZeroMantissa(mantissa)
		}
		exp = e
	}
	if isZeroMantissa(mantissa) {
		return true
	}
	// Quick exponent guard before arbitrary-precision big.Rat construction.
	if exp > 400 {
		return false
	}
	if exp < -400 {
		return true
	}
	r, ok := new(big.Rat).SetString(s)
	if !ok {
		return false
	}
	if r.Sign() < 0 {
		r.Neg(r)
	}
	return r.Cmp(maxBound) <= 0
}

func isZeroMantissa(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= '1' && s[i] <= '9' {
			return false
		}
	}
	return true
}

func decimalPrecisionWithinLimit(s string, maxPrec int) bool {
	mantissa := s
	exp := 0
	if idx := strings.IndexAny(s, "eE"); idx >= 0 {
		mantissa = s[:idx]
		e, err := strconv.Atoi(s[idx+1:])
		if err != nil {
			return false
		}
		exp = e
	}
	intPart, fracPart, hasDot := strings.Cut(mantissa, ".")
	if !hasDot {
		fracPart = ""
	}
	digits := strings.TrimLeft(intPart+fracPart, "0")
	scale := len(fracPart) - exp
	if len(digits) == 0 {
		return scale <= maxPrec && -scale <= maxPrec
	}
	// Spark computes BigDecimal precision as max(precision, scale) when scale >= 0,
	// or precision + (-scale) when scale < 0 (unscaled trailing zeros).
	prec := len(digits)
	if scale < 0 {
		prec += -scale
	} else if scale > prec {
		prec = scale
	}
	return prec <= maxPrec && scale <= maxPrec
}

// validateTypedLiteral checks Spark v4.0.5 AstBuilder.visitTypeConstructor
// parse-time validation for X, DATE, TIMESTAMP, TIMESTAMP_LTZ, TIMESTAMP_NTZ,
// and INTERVAL typed literals.
func validateTypedLiteral(kind TokenKind, rawLit string) bool {
	val := unquoteStringLit(rawLit)
	switch kind {
	case kwBINARY_HEX:
		for i := 0; i < len(val); i++ {
			c := val[i]
			if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
				return false
			}
		}
		return true
	case kwDATE:
		return isValidSparkDate(val)
	case kwTIMESTAMP, kwTIMESTAMP_LTZ, kwTIMESTAMP_NTZ:
		return isValidSparkTimestamp(val)
	case kwINTERVAL:
		return isValidSparkIntervalString(val)
	}
	return false
}

// isValidSparkConditionName checks Spark v4.0.5 AstBuilder.assertConditionName:
// decoded condition name must be non-empty ASCII [A-Za-z0-9_]+ regardless of quoting.
func isValidSparkConditionName(raw string) bool {
	if isQualifiedIdent(raw) {
		return false
	}
	decoded := exactIdentName(raw)
	if decoded == "" {
		return false
	}
	for i := 0; i < len(decoded); i++ {
		c := decoded[i]
		ok := (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '_'
		if !ok {
			return false
		}
	}
	return true
}

func isSpecialSparkDateTime(s string) bool {
	switch strings.ToLower(s) {
	case "epoch", "now", "today", "tomorrow", "yesterday":
		return true
	}
	return false
}

func isValidSparkDate(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	if isSpecialSparkDateTime(s) {
		return true
	}
	// Strip optional time suffix starting with ' ' or 'T'.
	if idx := strings.IndexAny(s, " Tt"); idx >= 0 {
		s = s[:idx]
	}
	return isValidSparkNumericDate(s)
}

func isValidSparkNumericDate(s string) bool {
	// Optional leading sign for +/-YYYY.
	if strings.HasPrefix(s, "+") || strings.HasPrefix(s, "-") {
		s = s[1:]
	}
	parts := strings.Split(s, "-")
	if len(parts) < 1 || len(parts) > 3 {
		return false
	}
	year, ok := parseDigits(parts[0], 4, 6)
	if !ok {
		return false
	}
	if len(parts) == 1 {
		return true
	}
	month, ok := parseDigits(parts[1], 1, 2)
	if !ok || month < 1 || month > 12 {
		return false
	}
	if len(parts) == 2 {
		return true
	}
	day, ok := parseDigits(parts[2], 1, 2)
	if !ok || day < 1 || day > daysInMonth(year, month) {
		return false
	}
	return true
}

func isValidSparkTimestamp(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	if isSpecialSparkDateTime(s) {
		return true
	}
	datePart := s
	timePart := ""
	if idx := strings.IndexAny(s, " Tt"); idx >= 0 {
		datePart = s[:idx]
		timePart = strings.TrimSpace(s[idx+1:])
	}
	if !isValidSparkNumericDate(datePart) {
		return false
	}
	if timePart == "" {
		return true
	}
	// Separate the clock prefix (digits, ':', '.') from any trailing timezone
	// suffix (e.g. `Z`, `+08:00`, `UTC+08:00`, ` America/Los_Angeles`, ` PST`).
	clockEnd := 0
	for clockEnd < len(timePart) {
		c := timePart[clockEnd]
		if (c >= '0' && c <= '9') || c == ':' || c == '.' {
			clockEnd++
		} else {
			break
		}
	}
	clockPart := timePart[:clockEnd]
	tzPart := strings.TrimSpace(timePart[clockEnd:])
	if clockPart == "" {
		return false
	}
	if tzPart != "" && !isValidSparkTimezone(tzPart) {
		return false
	}
	// Optional fractional seconds.
	if idx := strings.IndexByte(clockPart, '.'); idx >= 0 {
		frac := clockPart[idx+1:]
		if len(frac) == 0 || len(frac) > 9 {
			return false
		}
		for i := 0; i < len(frac); i++ {
			if frac[i] < '0' || frac[i] > '9' {
				return false
			}
		}
		clockPart = clockPart[:idx]
	}
	tparts := strings.Split(clockPart, ":")
	if len(tparts) < 1 || len(tparts) > 3 {
		return false
	}
	h, ok := parseDigits(tparts[0], 1, 2)
	if !ok || h < 0 || h > 23 {
		return false
	}
	if len(tparts) >= 2 {
		m, ok := parseDigits(tparts[1], 1, 2)
		if !ok || m < 0 || m > 59 {
			return false
		}
	}
	if len(tparts) == 3 {
		sec, ok := parseDigits(tparts[2], 1, 2)
		if !ok || sec < 0 || sec > 59 {
			return false
		}
	}
	return true
}

// sparkShortZoneIDs mirrors Java's java.time.ZoneId.SHORT_IDS accepted by
// Apache Spark's SparkDateTimeUtils.getZoneId.
var sparkShortZoneIDs = map[string]struct{}{
	"ACT": {}, "AET": {}, "AGT": {}, "ART": {}, "AST": {}, "BET": {}, "BST": {},
	"CAT": {}, "CNT": {}, "CST": {}, "CTT": {}, "EAT": {}, "ECT": {}, "EST": {},
	"HST": {}, "IET": {}, "IST": {}, "JST": {}, "MIT": {}, "MST": {}, "NET": {},
	"NST": {}, "PLT": {}, "PNT": {}, "PRT": {}, "PST": {}, "SST": {}, "VST": {},
}

func isValidSparkTimezone(tz string) bool {
	if tz == "" {
		return false
	}
	if tz == "Z" || tz == "z" {
		return true
	}
	// Prefix UTC / GMT / UT with optional offset (e.g. `UTC`, `UTC+08:00`, `GMT-5`).
	for _, prefix := range []string{"UTC", "GMT", "UT"} {
		if strings.HasPrefix(tz, prefix) {
			rest := tz[len(prefix):]
			if rest == "" {
				return true
			}
			return isValidSparkZoneOffset(rest)
		}
	}
	if tz[0] == '+' || tz[0] == '-' {
		return isValidSparkZoneOffset(tz)
	}
	if _, ok := sparkShortZoneIDs[tz]; ok {
		return true
	}
	// Reject Go's special "Local" sentinel and check IANA zone names.
	if strings.EqualFold(tz, "Local") {
		return false
	}
	_, err := time.LoadLocation(tz)
	return err == nil
}

// isValidSparkZoneOffset validates Java ZoneOffset / Spark timezone offsets:
// ±H, ±HH, ±H:MM, ±HH:MM, ±HHMM, ±HH:MM:SS, ±HHMMSS within [-18:00:00, +18:00:00].
func isValidSparkZoneOffset(s string) bool {
	if len(s) < 2 || (s[0] != '+' && s[0] != '-') {
		return false
	}
	body := s[1:]
	var h, m, sec int
	var ok bool
	if strings.ContainsRune(body, ':') {
		parts := strings.Split(body, ":")
		if len(parts) < 2 || len(parts) > 3 {
			return false
		}
		h, ok = parseDigits(parts[0], 1, 2)
		if !ok {
			return false
		}
		m, ok = parseDigits(parts[1], 2, 2)
		if !ok || m > 59 {
			return false
		}
		if len(parts) == 3 {
			sec, ok = parseDigits(parts[2], 2, 2)
			if !ok || sec > 59 {
				return false
			}
		}
	} else {
		switch len(body) {
		case 1, 2:
			h, ok = parseDigits(body, 1, 2)
			if !ok {
				return false
			}
		case 4:
			h, ok = parseDigits(body[:2], 2, 2)
			if !ok {
				return false
			}
			m, ok = parseDigits(body[2:], 2, 2)
			if !ok || m > 59 {
				return false
			}
		case 6:
			h, ok = parseDigits(body[:2], 2, 2)
			if !ok {
				return false
			}
			m, ok = parseDigits(body[2:4], 2, 2)
			if !ok || m > 59 {
				return false
			}
			sec, ok = parseDigits(body[4:], 2, 2)
			if !ok || sec > 59 {
				return false
			}
		default:
			return false
		}
	}
	total := h*3600 + m*60 + sec
	return total <= 18*3600
}

func parseDigits(s string, minLen, maxLen int) (int, bool) {
	if len(s) < minLen || len(s) > maxLen {
		return 0, false
	}
	n := 0
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return 0, false
		}
		n = n*10 + int(s[i]-'0')
	}
	return n, true
}

func daysInMonth(year, month int) int {
	switch month {
	case 1, 3, 5, 7, 8, 10, 12:
		return 31
	case 4, 6, 9, 11:
		return 30
	case 2:
		if (year%4 == 0 && year%100 != 0) || year%400 == 0 {
			return 29
		}
		return 28
	}
	return 0
}

func isValidSparkIntervalString(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	// Strip optional leading "interval " prefix (Spark allows INTERVAL 'interval 1 day').
	fields := strings.Fields(s)
	if len(fields) > 0 && strings.EqualFold(fields[0], "interval") {
		fields = fields[1:]
	}
	if len(fields) == 0 || len(fields)%2 != 0 {
		return false
	}
	for i := 0; i < len(fields); i += 2 {
		numStr, unitStr := fields[i], strings.ToUpper(fields[i+1])
		unitKind, ok := lookupKeyword(unitStr), true
		if !isSupportedMultiUnit(unitKind) {
			return false
		}
		if !isValidMultiUnitNumber(numStr, unitKind) {
			return false
		}
		_ = ok
	}
	return true
}

func isSupportedMultiUnit(k TokenKind) bool {
	switch k {
	case kwMICROSECOND, kwMICROSECONDS,
		kwMILLISECOND, kwMILLISECONDS,
		kwSECOND, kwSECONDS,
		kwMINUTE, kwMINUTES,
		kwHOUR, kwHOURS,
		kwDAY, kwDAYS,
		kwWEEK, kwWEEKS,
		kwMONTH, kwMONTHS,
		kwYEAR, kwYEARS:
		return true
	}
	return false
}

func isValidMultiUnitNumber(s string, unit TokenKind) bool {
	if s == "" {
		return false
	}
	if s[0] == '+' || s[0] == '-' {
		s = s[1:]
		if s == "" {
			return false
		}
	}
	if dot := strings.IndexByte(s, '.'); dot >= 0 {
		if unit != kwSECOND && unit != kwSECONDS {
			return false
		}
		left, right := s[:dot], s[dot+1:]
		if left == "" && right == "" {
			return false
		}
		for i := 0; i < len(left); i++ {
			if left[i] < '0' || left[i] > '9' {
				return false
			}
		}
		for i := 0; i < len(right); i++ {
			if right[i] < '0' || right[i] > '9' {
				return false
			}
		}
		return true
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	_, ok := new(big.Int).SetString(s, 10)
	return ok
}

// parseUnitToUnitIntervalValue validates Spark's from-to interval string literal
// and returns its duration in (months, microseconds) for timezone bounds checking.
func parseUnitToUnitIntervalValue(val string, fromUnit, toUnit TokenKind, outerSign int64) (months int64, micros int64, ok bool) {
	s := strings.TrimSpace(val)
	if s == "" {
		return 0, 0, false
	}
	sign := outerSign
	if s[0] == '+' {
		s = s[1:]
	} else if s[0] == '-' {
		sign = -sign
		s = s[1:]
	}
	const secUS int64 = 1_000_000
	switch {
	case fromUnit == kwYEAR && toUnit == kwMONTH:
		parts := strings.Split(s, "-")
		if len(parts) != 2 {
			return 0, 0, false
		}
		y, err1 := strconv.ParseInt(parts[0], 10, 64)
		m, err2 := strconv.ParseInt(parts[1], 10, 64)
		if err1 != nil || err2 != nil || m < 0 || m > 11 {
			return 0, 0, false
		}
		return sign * (y*12 + m), 0, true

	case fromUnit == kwDAY && toUnit == kwHOUR:
		parts := strings.Split(s, " ")
		if len(parts) != 2 {
			return 0, 0, false
		}
		d, err1 := strconv.ParseInt(parts[0], 10, 64)
		h, err2 := strconv.ParseInt(parts[1], 10, 64)
		if err1 != nil || err2 != nil || h < 0 || h > 23 {
			return 0, 0, false
		}
		return 0, sign * (d*86400 + h*3600) * secUS, true

	case fromUnit == kwDAY && toUnit == kwMINUTE:
		parts := strings.Split(s, " ")
		if len(parts) != 2 {
			return 0, 0, false
		}
		d, err1 := strconv.ParseInt(parts[0], 10, 64)
		hm := strings.Split(parts[1], ":")
		if err1 != nil || len(hm) != 2 {
			return 0, 0, false
		}
		h, err2 := strconv.ParseInt(hm[0], 10, 64)
		m, err3 := strconv.ParseInt(hm[1], 10, 64)
		if err2 != nil || err3 != nil || h < 0 || h > 23 || m < 0 || m > 59 {
			return 0, 0, false
		}
		return 0, sign * (d*86400 + h*3600 + m*60) * secUS, true

	case fromUnit == kwDAY && toUnit == kwSECOND:
		parts := strings.Split(s, " ")
		if len(parts) != 2 {
			return 0, 0, false
		}
		d, err1 := strconv.ParseInt(parts[0], 10, 64)
		hms := strings.Split(parts[1], ":")
		if err1 != nil || len(hms) != 3 {
			return 0, 0, false
		}
		h, err2 := strconv.ParseInt(hms[0], 10, 64)
		m, err3 := strconv.ParseInt(hms[1], 10, 64)
		secMicros, ok := parseSecondsMicros(hms[2])
		if err2 != nil || err3 != nil || !ok || h < 0 || h > 23 || m < 0 || m > 59 {
			return 0, 0, false
		}
		return 0, sign * ((d*86400+h*3600+m*60)*secUS + secMicros), true

	case fromUnit == kwHOUR && toUnit == kwMINUTE:
		hm := strings.Split(s, ":")
		if len(hm) != 2 {
			return 0, 0, false
		}
		h, err1 := strconv.ParseInt(hm[0], 10, 64)
		m, err2 := strconv.ParseInt(hm[1], 10, 64)
		if err1 != nil || err2 != nil || m < 0 || m > 59 {
			return 0, 0, false
		}
		return 0, sign * (h*3600 + m*60) * secUS, true

	case fromUnit == kwHOUR && toUnit == kwSECOND:
		hms := strings.Split(s, ":")
		if len(hms) != 3 {
			return 0, 0, false
		}
		h, err1 := strconv.ParseInt(hms[0], 10, 64)
		m, err2 := strconv.ParseInt(hms[1], 10, 64)
		secMicros, ok := parseSecondsMicros(hms[2])
		if err1 != nil || err2 != nil || !ok || m < 0 || m > 59 {
			return 0, 0, false
		}
		return 0, sign * ((h*3600+m*60)*secUS + secMicros), true

	case fromUnit == kwMINUTE && toUnit == kwSECOND:
		ms := strings.Split(s, ":")
		if len(ms) != 2 {
			return 0, 0, false
		}
		m, err1 := strconv.ParseInt(ms[0], 10, 64)
		secMicros, ok := parseSecondsMicros(ms[1])
		if err1 != nil || !ok {
			return 0, 0, false
		}
		return 0, sign * (m*60*secUS + secMicros), true
	}
	return 0, 0, false
}

// parseSecondsMicros parses SS[.fffffffff] (0 <= SS <= 59) and returns total microseconds.
func parseSecondsMicros(s string) (int64, bool) {
	secPart, fracPart, hasDot := strings.Cut(s, ".")
	sec, err := strconv.ParseInt(secPart, 10, 64)
	if err != nil || sec < 0 || sec > 59 {
		return 0, false
	}
	var fracMicros int64
	if hasDot {
		if len(fracPart) == 0 || len(fracPart) > 9 {
			return 0, false
		}
		for i := 0; i < len(fracPart); i++ {
			if fracPart[i] < '0' || fracPart[i] > '9' {
				return 0, false
			}
		}
		fracMicros = parseFracToMicros(fracPart)
	}
	return sec*1_000_000 + fracMicros, true
}

// parseFracToMicros converts a decimal fraction string (digits after '.') to
// microseconds (truncated/padded to 6 digits, flagging any nonzero sub-microsecond
// digits as 1 microsecond so non-whole-second checks reject them).
func parseFracToMicros(frac string) int64 {
	var us int64
	for i := 0; i < 6; i++ {
		us *= 10
		if i < len(frac) {
			us += int64(frac[i] - '0')
		}
	}
	if us == 0 {
		for i := 6; i < len(frac); i++ {
			if frac[i] != '0' {
				return 1
			}
		}
	}
	return us
}
