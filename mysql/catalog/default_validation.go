package catalog

import (
	"encoding/hex"
	"math"
	"math/big"
	"math/bits"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	nodes "github.com/bytebase/omni/mysql/ast"
)

// columnDefEffectiveDefaultExpr returns the last DEFAULT expression declared on
// colDef, matching the order createTable and buildColumnFromDef apply them.
func columnDefEffectiveDefaultExpr(colDef *nodes.ColumnDef) nodes.ExprNode {
	if colDef == nil {
		return nil
	}
	eff := colDef.DefaultValue
	for _, cc := range colDef.Constraints {
		if cc != nil && cc.Type == nodes.ColConstrDefault && cc.Expr != nil {
			eff = cc.Expr
		}
	}
	return eff
}

// validateAutoIncrementKeyed checks MySQL's post-DDL AUTO_INCREMENT invariant
// (1075, measured on mysql:8.0.46): a table may hold at most one AUTO_INCREMENT
// column, and an index must cover it — as the leading non-functional key part on
// InnoDB/MEMORY/default, or as any non-functional key part on MyISAM.
func validateAutoIncrementKeyed(tbl *Table) error {
	var autoCol *Column
	for _, col := range tbl.Columns {
		if col == nil || !col.AutoIncrement {
			continue
		}
		if autoCol != nil {
			return errIncorrectAutoIncrementDefinition()
		}
		autoCol = col
	}
	if autoCol == nil {
		return nil
	}
	myisam := strings.EqualFold(tbl.Engine, "MyISAM")
	for _, idx := range tbl.Indexes {
		if idx == nil || len(idx.Columns) == 0 {
			continue
		}
		if myisam {
			for _, ic := range idx.Columns {
				if ic != nil && ic.Expr == "" && strings.EqualFold(ic.Name, autoCol.Name) {
					return nil
				}
			}
		} else {
			ic := idx.Columns[0]
			if ic != nil && ic.Expr == "" && strings.EqualFold(ic.Name, autoCol.Name) {
				return nil
			}
		}
	}
	return errIncorrectAutoIncrementDefinition()
}

// validateTableDefaultExpressions validates every expression DEFAULT `(expr)` on
// the final table (measured on mysql:8.0.46): subqueries (3769), variables
// (3772), disallowed functions (3770), unknown columns (1054), auto-increment
// column references (3768), and self/forward references to columns that are
// generated or have an expression default (3767).
func validateTableDefaultExpressions(tbl *Table) error {
	for colIdx, col := range tbl.Columns {
		if col == nil || !col.DefaultIsExpr || col.DefaultExprNode == nil {
			continue
		}
		var err error
		nodes.Inspect(col.DefaultExprNode, func(n nodes.Node) bool {
			if err != nil {
				return false
			}
			switch x := n.(type) {
			case *nodes.SubqueryExpr, *nodes.ExistsExpr:
				err = errDefaultValExprSubquery(col.Name)
				return false
			case *nodes.InExpr:
				if x.Select != nil {
					err = errDefaultValExprSubquery(col.Name)
					return false
				}
			case *nodes.VariableRef:
				err = errDefaultValExprVariable(col.Name)
				return false
			case *nodes.DefaultExpr:
				err = &Error{Code: 3770, SQLState: "HY000",
					Message: "Default value expression of column '" + col.Name + "' contains a disallowed function."}
				return false
			case *nodes.FuncCallExpr:
				if x.SpacedParen {
					err = &Error{Code: 3770, SQLState: "HY000",
						Message: "Default value expression of column '" + col.Name + "' contains a disallowed function."}
					return false
				}
			case *nodes.ColumnRef:
				if x.Column == "" {
					return true
				}
				if x.Schema != "" || (x.Table != "" && !strings.EqualFold(x.Table, tbl.Name)) {
					err = errNoSuchColumn(x.Column, "default value expression")
					return false
				}
				refIdx, ok := tbl.colByName[toLower(x.Column)]
				if !ok || tbl.Columns[refIdx].Hidden == ColumnHiddenSystem {
					err = errNoSuchColumn(x.Column, "default value expression")
					return false
				}
				refCol := tbl.Columns[refIdx]
				if refCol.AutoIncrement {
					err = errDefaultValExprRefAutoInc(col.Name)
					return false
				}
				if refIdx >= colIdx && (refCol.Generated != nil || refCol.DefaultIsExpr) {
					err = errDefaultValExprForwardRef(col.Name)
					return false
				}
			}
			return true
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// validateEnumSetDuplicates checks an ENUM or SET member list for duplicates
// under the column's effective collation when strict sql_mode is active (1291,
// measured on mysql:8.0.46). Trailing spaces are stripped before comparison, and
// MySQL reports the first trimmed spelling.
func validateEnumSetDuplicates(colName, dataType, collation, charset string, values []string, strict bool) error {
	if !strict || !isEnumSetType(dataType) || len(values) == 0 {
		return nil
	}
	caseSensitive, turkish := collationSensitivity(collation, charset)
	seen := make(map[string]string, len(values))
	for _, raw := range values {
		trimmed := strings.TrimRight(raw, " ")
		key := enumCollationKey(trimmed, caseSensitive, turkish)
		if first, dup := seen[key]; dup {
			return errDuplicatedValueInType(colName, first, strings.ToUpper(dataType))
		}
		seen[key] = trimmed
	}
	return nil
}

func collationSensitivity(collation, charset string) (caseSensitive, turkish bool) {
	coll := strings.ToLower(strings.TrimSpace(collation))
	if coll == "" && charset != "" {
		coll = defaultCollationForCharset[toLower(charset)]
	}
	if coll == "" {
		coll = "utf8mb4_0900_ai_ci"
	}
	caseSensitive = coll == "binary" || strings.HasSuffix(coll, "_bin") || strings.HasSuffix(coll, "_cs")
	turkish = strings.Contains(coll, "_tr_") || strings.Contains(coll, "_turkish_")
	return caseSensitive, turkish
}

func enumCollationKey(s string, caseSensitive, turkish bool) string {
	if caseSensitive {
		return s
	}
	if turkish {
		return strings.ToLowerSpecial(unicode.TurkishCase, s)
	}
	return strings.ToLower(s)
}

func enumMembersFromColumnType(colType string) []string {
	open := strings.IndexByte(colType, '(')
	closeIdx := strings.LastIndexByte(colType, ')')
	if open < 0 || closeIdx <= open {
		return nil
	}
	return splitEnumMembers(colType[open+1 : closeIdx])
}

type literalDefaultKind int

const (
	litNone literalDefaultKind = iota
	litString
	litNumber
	litBit
	litHex
)

func extractDefaultLiteral(expr nodes.ExprNode) (literalDefaultKind, string, bool) {
	switch e := expr.(type) {
	case *nodes.StringLit:
		return litString, e.Value, false
	case *nodes.IntLit:
		if s, ok := exactIntLiteral(e, ""); ok {
			return litNumber, s, true
		}
		return litNumber, strconv.FormatInt(e.Value, 10), true
	case *nodes.FloatLit:
		return litNumber, e.Value, true
	case *nodes.UnaryExpr:
		if e.Op == nodes.UnaryMinus || e.Op == nodes.UnaryPlus {
			sign := ""
			if e.Op == nodes.UnaryMinus {
				sign = "-"
			}
			switch op := e.Operand.(type) {
			case *nodes.IntLit:
				if s, ok := exactIntLiteral(op, sign); ok {
					return litNumber, s, true
				}
				return litNumber, sign + strconv.FormatInt(op.Value, 10), true
			case *nodes.FloatLit:
				return litNumber, sign + op.Value, true
			}
		}
	case *nodes.BoolLit:
		if e.Value {
			return litNumber, "1", true
		}
		return litNumber, "0", true
	case *nodes.BitLit:
		return litBit, cleanBitDigits(e.Value), false
	case *nodes.HexLit:
		return litHex, cleanHexDigits(e.Value), false
	}
	return litNone, "", false
}

func cleanBitDigits(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "0b") || strings.HasPrefix(s, "0B") {
		return s[2:]
	}
	if len(s) >= 3 && (s[0] == 'b' || s[0] == 'B') && s[1] == '\'' && s[len(s)-1] == '\'' {
		return s[2 : len(s)-1]
	}
	return s
}

func cleanHexDigits(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
		return s[2:]
	}
	if len(s) >= 3 && (s[0] == 'x' || s[0] == 'X') && s[1] == '\'' && s[len(s)-1] == '\'' {
		return s[2 : len(s)-1]
	}
	return s
}

// validateColumnLiteralDefault checks an unparenthesized literal DEFAULT against
// the column's type and length/range (1067, measured on mysql:8.0.46).
func validateColumnLiteralDefault(col *Column, expr nodes.ExprNode, strict, noZeroDate, noZeroInDate bool) error {
	if col == nil || expr == nil {
		return nil
	}
	if _, isParen := expr.(*nodes.ParenExpr); isParen {
		return nil
	}
	if isNullDefault(expr) {
		return nil
	}
	if _, isTemporalFn := unwrapTemporalFunc(expr); isTemporalFn {
		return nil
	}
	dt := strings.ToLower(col.DataType)
	if isLiteralDefaultForbiddenType(dt) {
		return nil
	}

	kind, val, unquotedNum := extractDefaultLiteral(expr)
	if kind == litNone {
		return nil
	}

	unsigned := strings.Contains(strings.ToLower(col.ColumnType), "unsigned")

	switch dt {
	case "tinyint", "smallint", "mediumint", "int", "integer", "bigint", "boolean", "bool":
		if !validIntegerLiteralDefault(dt, unsigned, kind, val) {
			return errInvalidDefault(col.Name)
		}
	case "decimal", "numeric", "dec", "fixed":
		if !validDecimalLiteralDefault(col.ColumnType, unsigned, kind, val) {
			return errInvalidDefault(col.Name)
		}
	case "float", "double", "real":
		if !validFloatLiteralDefault(dt, unsigned, kind, val) {
			return errInvalidDefault(col.Name)
		}
	case "char", "varchar":
		if maxLen, ok := fixedLengthStringColumnLimit(col); ok {
			switch kind {
			case litString, litNumber:
				if utf8.RuneCountInString(val) > maxLen {
					return errInvalidDefault(col.Name)
				}
			}
		}
	case "binary", "varbinary":
		if maxLen, ok := fixedLengthStringColumnLimit(col); ok {
			switch kind {
			case litString, litNumber:
				if len(val) > maxLen {
					return errInvalidDefault(col.Name)
				}
			case litHex:
				if (len(val)+1)/2 > maxLen {
					return errInvalidDefault(col.Name)
				}
			case litBit:
				if (len(val)+7)/8 > maxLen {
					return errInvalidDefault(col.Name)
				}
			}
		}
	case "bit":
		if !validBitLiteralDefault(col.ColumnType, kind, val, unquotedNum) {
			return errInvalidDefault(col.Name)
		}
	case "year":
		if !validYearLiteralDefault(kind, val, unquotedNum) {
			return errInvalidDefault(col.Name)
		}
	case "enum":
		if !validEnumLiteralDefault(col, kind, val, unquotedNum, strict) {
			return errInvalidDefault(col.Name)
		}
	case "set":
		if !validSetLiteralDefault(col, kind, val, unquotedNum) {
			return errInvalidDefault(col.Name)
		}
	case "date", "datetime", "timestamp":
		if !validDateOrDatetimeLiteralDefault(dt, kind, val, unquotedNum, strict, noZeroDate, noZeroInDate) {
			return errInvalidDefault(col.Name)
		}
	case "time":
		if !validTimeLiteralDefault(kind, val, unquotedNum) {
			return errInvalidDefault(col.Name)
		}
	}
	return nil
}

func isStrictNumberSyntax(s string) bool {
	if s == "" {
		return false
	}
	i := 0
	if s[i] == '+' || s[i] == '-' {
		i++
	}
	if i >= len(s) {
		return false
	}
	digits := 0
	for i < len(s) && s[i] >= '0' && s[i] <= '9' {
		digits++
		i++
	}
	if i < len(s) && s[i] == '.' {
		i++
		for i < len(s) && s[i] >= '0' && s[i] <= '9' {
			digits++
			i++
		}
	}
	if digits == 0 {
		return false
	}
	if i < len(s) && (s[i] == 'e' || s[i] == 'E') {
		i++
		if i < len(s) && (s[i] == '+' || s[i] == '-') {
			i++
		}
		expDigits := 0
		for i < len(s) && s[i] >= '0' && s[i] <= '9' {
			expDigits++
			i++
		}
		if expDigits == 0 {
			return false
		}
	}
	return i == len(s)
}

func roundRatHalfAwayFromZero(r *big.Rat) *big.Int {
	sign := r.Sign()
	if sign == 0 {
		return big.NewInt(0)
	}
	abs := new(big.Rat).Abs(r)
	num := new(big.Int).Set(abs.Num())
	den := abs.Denom()
	q, rem := new(big.Int).QuoRem(num, den, new(big.Int))
	twiceRem := new(big.Int).Lsh(rem, 1)
	if twiceRem.Cmp(den) >= 0 {
		q.Add(q, big.NewInt(1))
	}
	if sign < 0 {
		q.Neg(q)
	}
	return q
}

func intBoundsForType(dt string, unsigned bool) (*big.Int, *big.Int) {
	switch dt {
	case "tinyint", "boolean", "bool":
		if unsigned {
			return big.NewInt(0), big.NewInt(255)
		}
		return big.NewInt(-128), big.NewInt(127)
	case "smallint":
		if unsigned {
			return big.NewInt(0), big.NewInt(65535)
		}
		return big.NewInt(-32768), big.NewInt(32767)
	case "mediumint":
		if unsigned {
			return big.NewInt(0), big.NewInt(16777215)
		}
		return big.NewInt(-8388608), big.NewInt(8388607)
	case "int", "integer":
		if unsigned {
			return big.NewInt(0), big.NewInt(4294967295)
		}
		return big.NewInt(-2147483648), big.NewInt(2147483647)
	default: // bigint
		if unsigned {
			maxU := new(big.Int).SetUint64(math.MaxUint64)
			return big.NewInt(0), maxU
		}
		return big.NewInt(math.MinInt64), big.NewInt(math.MaxInt64)
	}
}

func validIntegerLiteralDefault(dt string, unsigned bool, kind literalDefaultKind, val string) bool {
	minB, maxB := intBoundsForType(dt, unsigned)
	var n *big.Int
	switch kind {
	case litBit:
		trimmed := strings.TrimLeft(val, "0")
		if len(trimmed) > 64 {
			return false
		}
		if trimmed == "" {
			n = big.NewInt(0)
		} else {
			var ok bool
			n, ok = new(big.Int).SetString(trimmed, 2)
			if !ok {
				return false
			}
		}
	case litHex:
		trimmed := strings.TrimLeft(val, "0")
		if len(trimmed) > 16 {
			return false
		}
		if trimmed == "" {
			n = big.NewInt(0)
		} else {
			var ok bool
			n, ok = new(big.Int).SetString(trimmed, 16)
			if !ok {
				return false
			}
		}
	default:
		s := strings.TrimSpace(val)
		if !isStrictNumberSyntax(s) {
			return false
		}
		r, ok := new(big.Rat).SetString(s)
		if !ok {
			return false
		}
		n = roundRatHalfAwayFromZero(r)
	}
	return n.Cmp(minB) >= 0 && n.Cmp(maxB) <= 0
}

func validDecimalLiteralDefault(colType string, unsigned bool, kind literalDefaultKind, val string) bool {
	if kind == litBit || kind == litHex {
		return true
	}
	s := strings.TrimSpace(val)
	if !isStrictNumberSyntax(s) {
		return false
	}
	r, ok := new(big.Rat).SetString(s)
	if !ok {
		return false
	}
	if unsigned && r.Sign() < 0 {
		return false
	}
	pt := parseColumnType(colType)
	m, d := 10, 0
	if pt.hasLen {
		m = pt.length
	}
	if pt.hasScale && pt.scale >= 0 {
		d = pt.scale
	}
	scalePow := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(d)), nil)
	scaled := new(big.Rat).Mul(new(big.Rat).Abs(r), new(big.Rat).SetInt(scalePow))
	roundedScaled := roundRatHalfAwayFromZero(scaled)
	maxDigits := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(m)), nil)
	return roundedScaled.Cmp(maxDigits) < 0
}

func validFloatLiteralDefault(dt string, unsigned bool, kind literalDefaultKind, val string) bool {
	if kind == litBit || kind == litHex {
		return true
	}
	s := strings.TrimSpace(val)
	if !isStrictNumberSyntax(s) {
		return false
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil || math.IsInf(f, 0) || math.IsNaN(f) {
		return false
	}
	if unsigned && f < 0 {
		return false
	}
	if dt == "float" && math.Abs(f) > math.MaxFloat32 {
		return false
	}
	return true
}

func bitColumnWidth(colType string) int {
	pt := parseColumnType(colType)
	if pt.hasLen && pt.length > 0 {
		return pt.length
	}
	return 1
}

func validBitLiteralDefault(colType string, kind literalDefaultKind, val string, unquotedNum bool) bool {
	width := bitColumnWidth(colType)
	switch kind {
	case litBit:
		trimmed := strings.TrimLeft(val, "0")
		return len(trimmed) <= width
	case litHex:
		raw := val
		if len(raw)%2 == 1 {
			raw = "0" + raw
		}
		b, err := hex.DecodeString(raw)
		if err != nil {
			return false
		}
		return new(big.Int).SetBytes(b).BitLen() <= width
	case litNumber:
		if unquotedNum {
			s := strings.TrimSpace(val)
			if strings.HasPrefix(s, "-") {
				return false
			}
			s = strings.TrimPrefix(s, "+")
			u, err := strconv.ParseUint(s, 10, 64)
			if err != nil {
				return false
			}
			return bits.Len64(u) <= width
		}
	}
	// Quoted string on BIT(N): MySQL interprets raw bytes as big-endian unsigned value.
	return new(big.Int).SetBytes([]byte(val)).BitLen() <= width
}

func validYearLiteralDefault(kind literalDefaultKind, val string, unquotedNum bool) bool {
	if kind == litBit || kind == litHex {
		return false
	}
	s := strings.TrimSpace(val)
	if s == "" {
		return false
	}
	if unquotedNum {
		n, err := strconv.Atoi(s)
		if err != nil {
			return false
		}
		return n == 0 || (n >= 1 && n <= 99) || (n >= 1901 && n <= 2155)
	}
	if len(s) > 4 {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return false
	}
	if n == 0 {
		return true
	}
	return (len(s) <= 2 && n >= 1 && n <= 99) || (n >= 1901 && n <= 2155)
}

func validEnumLiteralDefault(col *Column, kind literalDefaultKind, val string, unquotedNum, strict bool) bool {
	members := enumMembersFromColumnType(col.ColumnType)
	if unquotedNum {
		n, err := strconv.Atoi(strings.TrimSpace(val))
		if err != nil || n < 0 || n > len(members) {
			return false
		}
		if n == 0 && strict {
			return false
		}
		return true
	}
	if kind != litString {
		return true
	}
	caseSensitive, turkish := collationSensitivity(col.Collation, col.Charset)
	target := enumCollationKey(strings.TrimRight(val, " "), caseSensitive, turkish)
	for _, m := range members {
		if enumCollationKey(strings.TrimRight(m, " "), caseSensitive, turkish) == target {
			return true
		}
	}
	if target == "" && !strict {
		return true
	}
	return false
}

func validSetLiteralDefault(col *Column, kind literalDefaultKind, val string, unquotedNum bool) bool {
	members := enumMembersFromColumnType(col.ColumnType)
	if unquotedNum {
		n, err := strconv.ParseInt(strings.TrimSpace(val), 10, 64)
		if err != nil || n < 0 {
			return false
		}
		if len(members) < 63 && n >= (int64(1)<<len(members)) {
			return false
		}
		return true
	}
	if kind != litString {
		return true
	}
	if val == "" {
		return true
	}
	caseSensitive, turkish := collationSensitivity(col.Collation, col.Charset)
	allowed := make(map[string]bool, len(members))
	for _, m := range members {
		allowed[enumCollationKey(strings.TrimRight(m, " "), caseSensitive, turkish)] = true
	}
	for _, part := range strings.Split(val, ",") {
		if !allowed[enumCollationKey(strings.TrimRight(part, " "), caseSensitive, turkish)] {
			return false
		}
	}
	return true
}

func validDateOrDatetimeLiteralDefault(dt string, kind literalDefaultKind, val string, unquotedNum, strict, noZeroDate, noZeroInDate bool) bool {
	if kind != litString && kind != litNumber {
		return false
	}
	s := strings.TrimSpace(val)
	if s == "" {
		return false
	}
	year, month, day, hour, min, sec, ok := parseMySQLDateTimeParts(s, unquotedNum)
	if !ok {
		return false
	}
	allZero := year == 0 && month == 0 && day == 0 && hour == 0 && min == 0 && sec == 0
	if allZero {
		return !noZeroDate
	}
	if month > 12 || day > 31 || hour > 23 || min > 59 || sec > 59 {
		return false
	}
	if month == 0 || day == 0 {
		if dt == "timestamp" || noZeroInDate {
			return false
		}
		return true
	}
	if day > daysInMonth(year, month) {
		return false
	}
	if dt == "timestamp" {
		if year < 1970 || year > 2038 {
			return false
		}
		if year == 1970 && month == 1 && day == 1 && hour == 0 && min == 0 && sec == 0 {
			return false
		}
		if year == 2038 && (month > 1 || day > 19 || (day == 19 && (hour > 3 || (hour == 3 && (min > 14 || (min == 14 && sec > 7)))))) {
			return false
		}
	}
	return true
}

func parseMySQLDateTimeParts(s string, unquotedNum bool) (year, month, day, hour, min, sec int, ok bool) {
	if unquotedNum && s == "0" {
		return 0, 0, 0, 0, 0, 0, true
	}
	base := s
	if dot := strings.LastIndexByte(base, '.'); dot >= 0 {
		frac := base[dot+1:]
		allDigits := len(frac) > 0
		for i := 0; i < len(frac); i++ {
			if frac[i] < '0' || frac[i] > '9' {
				allDigits = false
				break
			}
		}
		// Only strip fractional seconds when there is a time part (or compact digits).
		if allDigits && (strings.ContainsAny(base[:dot], " Tt:") || len(base[:dot]) >= 12) {
			base = base[:dot]
		}
	}
	// Compact all-digit forms: YYMMDD (6), YYYYMMDD (8), YYMMDDHHMMSS (12), YYYYMMDDHHMMSS (14).
	allDigits := len(base) > 0
	for i := 0; i < len(base); i++ {
		if base[i] < '0' || base[i] > '9' {
			allDigits = false
			break
		}
	}
	if allDigits {
		switch len(base) {
		case 6:
			y, _ := strconv.Atoi(base[0:2])
			m, _ := strconv.Atoi(base[2:4])
			d, _ := strconv.Atoi(base[4:6])
			return expandTwoDigitYear(y, m, d), m, d, 0, 0, 0, true
		case 8:
			y, _ := strconv.Atoi(base[0:4])
			m, _ := strconv.Atoi(base[4:6])
			d, _ := strconv.Atoi(base[6:8])
			return y, m, d, 0, 0, 0, true
		case 12:
			y, _ := strconv.Atoi(base[0:2])
			m, _ := strconv.Atoi(base[2:4])
			d, _ := strconv.Atoi(base[4:6])
			hh, _ := strconv.Atoi(base[6:8])
			mm, _ := strconv.Atoi(base[8:10])
			ss, _ := strconv.Atoi(base[10:12])
			return expandTwoDigitYear(y, m, d), m, d, hh, mm, ss, true
		case 14:
			y, _ := strconv.Atoi(base[0:4])
			m, _ := strconv.Atoi(base[4:6])
			d, _ := strconv.Atoi(base[6:8])
			hh, _ := strconv.Atoi(base[8:10])
			mm, _ := strconv.Atoi(base[10:12])
			ss, _ := strconv.Atoi(base[12:14])
			return y, m, d, hh, mm, ss, true
		default:
			return 0, 0, 0, 0, 0, 0, false
		}
	}

	// Tokenize digit runs separated by punctuation / 'T' / space (and timezone offset +/-HH:MM).
	// Reject any alphabetic characters other than a single 'T'/'t' between date and time, or trailing 'Z'/'z'.
	cleaned := strings.TrimSuffix(strings.TrimSuffix(s, "Z"), "z")
	var parts []string
	i := 0
	for i < len(cleaned) {
		ch := cleaned[i]
		if ch >= '0' && ch <= '9' {
			j := i
			for j < len(cleaned) && cleaned[j] >= '0' && cleaned[j] <= '9' {
				j++
			}
			parts = append(parts, cleaned[i:j])
			i = j
			continue
		}
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') {
			if (ch == 'T' || ch == 't') && len(parts) == 3 {
				i++
				continue
			}
			return 0, 0, 0, 0, 0, 0, false
		}
		i++
	}
	if len(parts) < 3 || len(parts) > 9 {
		return 0, 0, 0, 0, 0, 0, false
	}
	y, err1 := strconv.Atoi(parts[0])
	m, err2 := strconv.Atoi(parts[1])
	d, err3 := strconv.Atoi(parts[2])
	if err1 != nil || err2 != nil || err3 != nil || len(parts[0]) > 4 {
		return 0, 0, 0, 0, 0, 0, false
	}
	if len(parts[0]) <= 2 {
		y = expandTwoDigitYear(y, m, d)
	}
	if len(parts) >= 4 {
		hour, _ = strconv.Atoi(parts[3])
	}
	if len(parts) >= 5 {
		min, _ = strconv.Atoi(parts[4])
	}
	if len(parts) >= 6 {
		sec, _ = strconv.Atoi(parts[5])
	}
	return y, m, d, hour, min, sec, true
}

func expandTwoDigitYear(y, m, d int) int {
	if y == 0 && m == 0 && d == 0 {
		return 0
	}
	if y >= 70 {
		return 1900 + y
	}
	return 2000 + y
}

func validTimeLiteralDefault(kind literalDefaultKind, val string, unquotedNum bool) bool {
	if kind != litString && kind != litNumber {
		return false
	}
	s := strings.TrimSpace(val)
	if s == "" {
		return false
	}
	if strings.HasPrefix(s, "-") || strings.HasPrefix(s, "+") {
		s = strings.TrimSpace(s[1:])
	}
	if s == "" {
		return false
	}
	if dot := strings.IndexByte(s, '.'); dot >= 0 {
		frac := s[dot+1:]
		if frac == "" || len(frac) > 6 {
			return false
		}
		for i := 0; i < len(frac); i++ {
			if frac[i] < '0' || frac[i] > '9' {
				return false
			}
		}
		s = s[:dot]
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if (c < '0' || c > '9') && c != ':' && c != ' ' {
			return false
		}
	}
	days := 0
	if sp := strings.IndexByte(s, ' '); sp >= 0 {
		d, err := strconv.Atoi(strings.TrimSpace(s[:sp]))
		if err != nil || d > 34 {
			return false
		}
		days = d
		s = strings.TrimSpace(s[sp+1:])
	}
	if strings.Contains(s, ":") {
		fields := strings.Split(s, ":")
		if len(fields) < 2 || len(fields) > 3 {
			return false
		}
		hh, err1 := strconv.Atoi(fields[0])
		mm, err2 := strconv.Atoi(fields[1])
		ss := 0
		var err3 error
		if len(fields) == 3 {
			ss, err3 = strconv.Atoi(fields[2])
		}
		if err1 != nil || err2 != nil || err3 != nil {
			return false
		}
		totalHours := days*24 + hh
		return mm < 60 && ss < 60 && totalHours <= 838
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return false
	}
	ss := n % 100
	mm := (n / 100) % 100
	hh := n/10000 + days*24
	return mm < 60 && ss < 60 && hh <= 838
}
