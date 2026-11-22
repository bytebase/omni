package parser

import (
	"math"
	"math/big"
	"strconv"
	"strings"
)

// exactIdentName unquotes a backtick-quoted identifier (unescaping “) while
// preserving case, matching Spark v4.0.5's case-sensitive CTE and property checks.
func exactIdentName(s string) string {
	if len(s) >= 2 && s[0] == '`' && s[len(s)-1] == '`' {
		return strings.ReplaceAll(s[1:len(s)-1], "``", "`")
	}
	return s
}

// unquoteStringLit decodes a Spark single-quoted, double-quoted, or raw string
// literal (R'...' / R"...") into its logical value.
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
			case '\\', '\'', '"':
				b.WriteByte(inner[i])
			default:
				b.WriteByte('\\')
				b.WriteByte(inner[i])
			}
			continue
		}
		b.WriteByte(inner[i])
	}
	return b.String()
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
	parts := strings.Split(raw, ".")
	last := strings.ToUpper(exactIdentName(parts[len(parts)-1]))
	return last == "TRIM" || last == "RTRIM" ||
		strings.HasSuffix(last, "_TRIM") || strings.HasSuffix(last, "_RTRIM")
}

// isValidSparkCollation validates a collation name (optionally qualified with
// SYSTEM.BUILTIN.) against Apache Spark 4.0.5's CollationFactory, rejecting
// unknown names and unsupported TRIM/RTRIM collations.
func isValidSparkCollation(raw string) bool {
	parts := strings.Split(raw, ".")
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
		num := text[:len(text)-1]
		if negative {
			num = "-" + num
		}
		v, err := strconv.ParseFloat(num, 32)
		if err != nil || math.IsInf(v, 0) || math.IsNaN(v) {
			return false
		}
		// Spark also rejects float32 boundary values whose BigDecimal conversion
		// rounds to Infinity (e.g. 3.4028235e38F).
		f32 := float32(v)
		return !math.IsInf(float64(f32), 0) && math.Abs(v) < 3.4028235e38
	case tokDouble:
		num := text[:len(text)-1]
		if negative {
			num = "-" + num
		}
		v, err := strconv.ParseFloat(num, 64)
		return err == nil && !math.IsInf(v, 0) && !math.IsNaN(v)
	case tokExponent:
		num := text
		if negative {
			num = "-" + num
		}
		v, err := strconv.ParseFloat(num, 64)
		return err == nil && !math.IsInf(v, 0) && !math.IsNaN(v)
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

func isValidSparkDate(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	// Strip optional time suffix starting with ' ' or 'T'.
	if idx := strings.IndexAny(s, " Tt"); idx >= 0 {
		s = s[:idx]
	}
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
	datePart := s
	timePart := ""
	if idx := strings.IndexAny(s, " Tt"); idx >= 0 {
		datePart = s[:idx]
		timePart = strings.TrimSpace(s[idx+1:])
	}
	if !isValidSparkDate(datePart) {
		return false
	}
	if timePart == "" {
		return true
	}
	// Strip optional trailing timezone (Z, +HH:MM, -HH:MM, +HH, -HH, or named zone after space).
	if idx := strings.IndexByte(timePart, ' '); idx >= 0 {
		timePart = timePart[:idx]
	}
	if strings.HasSuffix(timePart, "Z") || strings.HasSuffix(timePart, "z") {
		timePart = timePart[:len(timePart)-1]
	} else if idx := strings.IndexAny(timePart, "+-"); idx > 0 {
		timePart = timePart[:idx]
	}
	// Optional fractional seconds.
	if idx := strings.IndexByte(timePart, '.'); idx >= 0 {
		frac := timePart[idx+1:]
		if len(frac) == 0 || len(frac) > 9 {
			return false
		}
		for i := 0; i < len(frac); i++ {
			if frac[i] < '0' || frac[i] > '9' {
				return false
			}
		}
		timePart = timePart[:idx]
	}
	tparts := strings.Split(timePart, ":")
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
// and returns its total duration in seconds (for timezone bounds checking).
func parseUnitToUnitIntervalValue(val string, fromUnit, toUnit TokenKind, outerSign int64) (int64, bool) {
	s := strings.TrimSpace(val)
	if s == "" {
		return 0, false
	}
	sign := outerSign
	if s[0] == '+' {
		s = s[1:]
	} else if s[0] == '-' {
		sign = -sign
		s = s[1:]
	}
	switch {
	case fromUnit == kwYEAR && toUnit == kwMONTH:
		parts := strings.Split(s, "-")
		if len(parts) != 2 {
			return 0, false
		}
		y, err1 := strconv.ParseInt(parts[0], 10, 64)
		m, err2 := strconv.ParseInt(parts[1], 10, 64)
		if err1 != nil || err2 != nil || m < 0 || m > 11 {
			return 0, false
		}
		return sign * (y*365*86400 + m*30*86400), true

	case fromUnit == kwDAY && toUnit == kwHOUR:
		parts := strings.Split(s, " ")
		if len(parts) != 2 {
			return 0, false
		}
		d, err1 := strconv.ParseInt(parts[0], 10, 64)
		h, err2 := strconv.ParseInt(parts[1], 10, 64)
		if err1 != nil || err2 != nil || h < 0 || h > 23 {
			return 0, false
		}
		return sign * (d*86400 + h*3600), true

	case fromUnit == kwDAY && toUnit == kwMINUTE:
		parts := strings.Split(s, " ")
		if len(parts) != 2 {
			return 0, false
		}
		d, err1 := strconv.ParseInt(parts[0], 10, 64)
		hm := strings.Split(parts[1], ":")
		if err1 != nil || len(hm) != 2 {
			return 0, false
		}
		h, err2 := strconv.ParseInt(hm[0], 10, 64)
		m, err3 := strconv.ParseInt(hm[1], 10, 64)
		if err2 != nil || err3 != nil || h < 0 || h > 23 || m < 0 || m > 59 {
			return 0, false
		}
		return sign * (d*86400 + h*3600 + m*60), true

	case fromUnit == kwDAY && toUnit == kwSECOND:
		parts := strings.Split(s, " ")
		if len(parts) != 2 {
			return 0, false
		}
		d, err1 := strconv.ParseInt(parts[0], 10, 64)
		hms := strings.Split(parts[1], ":")
		if err1 != nil || len(hms) != 3 {
			return 0, false
		}
		h, err2 := strconv.ParseInt(hms[0], 10, 64)
		m, err3 := strconv.ParseInt(hms[1], 10, 64)
		sec, ok := parseSecondsWithOptionalFrac(hms[2])
		if err2 != nil || err3 != nil || !ok || h < 0 || h > 23 || m < 0 || m > 59 {
			return 0, false
		}
		return sign * (d*86400 + h*3600 + m*60 + sec), true

	case fromUnit == kwHOUR && toUnit == kwMINUTE:
		hm := strings.Split(s, ":")
		if len(hm) != 2 {
			return 0, false
		}
		h, err1 := strconv.ParseInt(hm[0], 10, 64)
		m, err2 := strconv.ParseInt(hm[1], 10, 64)
		if err1 != nil || err2 != nil || m < 0 || m > 59 {
			return 0, false
		}
		return sign * (h*3600 + m*60), true

	case fromUnit == kwHOUR && toUnit == kwSECOND:
		hms := strings.Split(s, ":")
		if len(hms) != 3 {
			return 0, false
		}
		h, err1 := strconv.ParseInt(hms[0], 10, 64)
		m, err2 := strconv.ParseInt(hms[1], 10, 64)
		sec, ok := parseSecondsWithOptionalFrac(hms[2])
		if err1 != nil || err2 != nil || !ok || m < 0 || m > 59 {
			return 0, false
		}
		return sign * (h*3600 + m*60 + sec), true

	case fromUnit == kwMINUTE && toUnit == kwSECOND:
		ms := strings.Split(s, ":")
		if len(ms) != 2 {
			return 0, false
		}
		m, err1 := strconv.ParseInt(ms[0], 10, 64)
		sec, ok := parseSecondsWithOptionalFrac(ms[1])
		if err1 != nil || !ok {
			return 0, false
		}
		return sign * (m*60 + sec), true
	}
	return 0, false
}

func parseSecondsWithOptionalFrac(s string) (int64, bool) {
	secPart, fracPart, hasDot := strings.Cut(s, ".")
	sec, err := strconv.ParseInt(secPart, 10, 64)
	if err != nil || sec < 0 || sec > 59 {
		return 0, false
	}
	if hasDot {
		if len(fracPart) == 0 || len(fracPart) > 9 {
			return 0, false
		}
		for i := 0; i < len(fracPart); i++ {
			if fracPart[i] < '0' || fracPart[i] > '9' {
				return 0, false
			}
		}
	}
	return sec, true
}
