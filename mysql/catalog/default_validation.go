package catalog

import (
	"bytes"
	"encoding/hex"
	"errors"
	"math"
	"math/big"
	"math/bits"
	"strconv"
	"strings"
	"time"
	"unicode/utf16"
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

// validateDroppedOrRenamedDefaultExprDeps enforces MySQL 3766 when an ALTER
// TABLE drops or renames a column that a surviving pre-existing expression
// DEFAULT still references (measured on mysql:8.0.46).
func validateDroppedOrRenamedDefaultExprDeps(tbl *Table, origTableName string, droppedOrRenamed []string, preExprs map[nodes.ExprNode]bool) error {
	for _, target := range droppedOrRenamed {
		for _, col := range tbl.Columns {
			if col == nil || !col.DefaultIsExpr || col.DefaultExprNode == nil || !preExprs[col.DefaultExprNode] {
				continue
			}
			if exprReferencesColumn(col.DefaultExprNode, tbl.Name, target) {
				return errDependentByDefaultValExpr(target, origTableName)
			}
		}
	}
	return nil
}

func exprReferencesColumn(expr nodes.ExprNode, tableName, colName string) bool {
	found := false
	nodes.Inspect(expr, func(n nodes.Node) bool {
		if found {
			return false
		}
		if cr, ok := n.(*nodes.ColumnRef); ok && cr.Column != "" {
			if cr.Schema == "" && (cr.Table == "" || strings.EqualFold(cr.Table, tableName)) && strings.EqualFold(cr.Column, colName) {
				found = true
				return false
			}
		}
		return true
	})
	return found
}

// validateEnumSetDuplicates checks an ENUM or SET member list for duplicates
// under the column's effective collation when strict sql_mode is active (1291,
// measured on mysql:8.0.46). Trailing spaces are stripped before comparison, and
// MySQL reports the first trimmed spelling.
func validateEnumSetDuplicates(colName, dataType, collation, charset string, values []string, strict bool) error {
	if !strict || !isEnumSetType(dataType) || len(values) == 0 {
		return nil
	}
	seen := make(map[string]string, len(values))
	for _, raw := range values {
		trimmed := strings.TrimRight(raw, " ")
		key := collationComparisonKey(trimmed, collation, charset)
		if first, dup := seen[key]; dup {
			return errDuplicatedValueInType(colName, first, strings.ToUpper(dataType))
		}
		seen[key] = trimmed
	}
	return nil
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
// the column's type and length/range (1067, measured on mysql:8.0.46), and
// stores the coerced value MySQL records in the data dictionary.
func (c *Catalog) validateColumnLiteralDefault(col *Column, expr nodes.ExprNode) error {
	strict := c.isStrictSQLMode()
	noZeroDate := c.isNoZeroDateSQLMode()
	noZeroInDate := c.isNoZeroInDateSQLMode()
	allowInvalidDates := c.isAllowInvalidDatesSQLMode()
	sessionTZ := c.session.TimeZone
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
	if kind == litNumber && strings.ContainsAny(val, "eE") {
		s := strings.TrimSpace(val)
		if isStrictNumberSyntax(s) {
			if _, err := strconv.ParseFloat(s, 64); errors.Is(err, strconv.ErrRange) {
				return errIllegalValueForType("double", s)
			}
		}
	}

	unsigned := strings.Contains(strings.ToLower(col.ColumnType), "unsigned")

	switch dt {
	case "tinyint", "smallint", "mediumint", "int", "integer", "bigint", "boolean", "bool":
		if !validIntegerLiteralDefault(dt, unsigned, kind, val) {
			return errInvalidDefault(col.Name)
		}
	case "decimal", "numeric", "dec", "fixed":
		if !validDecimalLiteralDefault(col.ColumnType, unsigned, kind, val, strict) {
			return errInvalidDefault(col.Name)
		}
	case "float", "double", "real":
		if !validFloatLiteralDefault(dt, col.ColumnType, unsigned, kind, val, strict) {
			return errInvalidDefault(col.Name)
		}
	case "char", "varchar":
		if maxLen, ok := fixedLengthStringColumnLimit(col); ok {
			switch kind {
			case litString:
				// MySQL truncates excess trailing U+0020 spaces on CHAR and VARCHAR
				// defaults rather than rejecting the column definition with 1067.
				if utf8.RuneCountInString(strings.TrimRight(val, " ")) > maxLen {
					return errInvalidDefault(col.Name)
				}
			case litNumber:
				if !validNumericDefaultForStringLen(val, maxLen) {
					return errInvalidDefault(col.Name)
				}
			case litHex:
				if !validBinaryBytesForCharColumn(val, true, col.Charset, maxLen) {
					return errInvalidDefault(col.Name)
				}
			case litBit:
				if !validBinaryBytesForCharColumn(val, false, col.Charset, maxLen) {
					return errInvalidDefault(col.Name)
				}
			}
		}
	case "binary", "varbinary":
		if maxLen, ok := fixedLengthStringColumnLimit(col); ok {
			switch kind {
			case litString:
				if len(val) > maxLen {
					return errInvalidDefault(col.Name)
				}
			case litNumber:
				if !validNumericDefaultForStringLen(val, maxLen) {
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
		if !validDateOrDatetimeLiteralDefault(dt, kind, val, unquotedNum, noZeroDate, noZeroInDate, allowInvalidDates, sessionTZ) {
			return errInvalidDefault(col.Name)
		}
	case "time":
		if !validTimeLiteralDefault(kind, val, unquotedNum) {
			return errInvalidDefault(col.Name)
		}
	}
	normalizeColumnStoredDefault(col, expr, strict)
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

func evalIntegerLiteralBigInt(kind literalDefaultKind, val string) (*big.Int, bool) {
	switch kind {
	case litBit:
		trimmed := strings.TrimLeft(val, "0")
		if len(trimmed) > 64 {
			return nil, false
		}
		if trimmed == "" {
			return big.NewInt(0), true
		}
		return new(big.Int).SetString(trimmed, 2)
	case litHex:
		trimmed := strings.TrimLeft(val, "0")
		if len(trimmed) > 16 {
			return nil, false
		}
		if trimmed == "" {
			return big.NewInt(0), true
		}
		return new(big.Int).SetString(trimmed, 16)
	default:
		s := strings.TrimSpace(val)
		if !isStrictNumberSyntax(s) {
			return nil, false
		}
		r, ok := new(big.Rat).SetString(s)
		if !ok {
			return nil, false
		}
		return roundRatHalfAwayFromZero(r), true
	}
}

func validIntegerLiteralDefault(dt string, unsigned bool, kind literalDefaultKind, val string) bool {
	minB, maxB := intBoundsForType(dt, unsigned)
	n, ok := evalIntegerLiteralBigInt(kind, val)
	if !ok {
		return false
	}
	return n.Cmp(minB) >= 0 && n.Cmp(maxB) <= 0
}

func formatDecimalLiteralDefault(colType string, unsigned bool, kind literalDefaultKind, val string, strict bool) (string, bool) {
	pt := parseColumnType(colType)
	m, d := 10, 0
	if pt.hasLen {
		m = pt.length
	}
	if pt.hasScale && pt.scale >= 0 {
		d = pt.scale
	}
	if kind == litBit || kind == litHex {
		raw, ok := decodeLiteralRawBytes(val, kind == litHex)
		if !ok {
			return "", false
		}
		n := new(big.Int).SetBytes(raw)
		intDigits := m - d
		if intDigits < 0 {
			intDigits = 0
		}
		maxInt := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(intDigits)), nil)
		if n.Cmp(maxInt) >= 0 {
			return "", false
		}
		if d == 0 {
			return n.String(), true
		}
		return n.String() + "." + strings.Repeat("0", d), true
	}
	s := strings.TrimSpace(val)
	if !isStrictNumberSyntax(s) {
		if strict {
			return "", false
		}
		s = leadingNumericPrefix(s)
	}
	r, ok := new(big.Rat).SetString(s)
	if !ok {
		return "", false
	}
	if unsigned && r.Sign() < 0 {
		return "", false
	}
	scalePow := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(d)), nil)
	scaled := new(big.Rat).Mul(new(big.Rat).Abs(r), new(big.Rat).SetInt(scalePow))
	roundedScaled := roundRatHalfAwayFromZero(scaled)
	maxDigits := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(m)), nil)
	if roundedScaled.Cmp(maxDigits) >= 0 {
		return "", false
	}
	neg := r.Sign() < 0 && roundedScaled.Sign() > 0
	digits := roundedScaled.String()
	if d == 0 {
		if neg {
			return "-" + digits, true
		}
		return digits, true
	}
	if len(digits) <= d {
		digits = strings.Repeat("0", d+1-len(digits)) + digits
	}
	formatted := digits[:len(digits)-d] + "." + digits[len(digits)-d:]
	if neg {
		return "-" + formatted, true
	}
	return formatted, true
}

func coerceNumericColumnDefault(col *Column, expr nodes.ExprNode, strict bool) (string, bool) {
	if col == nil || expr == nil {
		return "", false
	}
	if _, isParen := expr.(*nodes.ParenExpr); isParen {
		return "", false
	}
	if isNullDefault(expr) {
		return "", false
	}
	kind, val, _ := extractDefaultLiteral(expr)
	if kind == litNone {
		return "", false
	}
	unsigned := strings.Contains(strings.ToLower(col.ColumnType), "unsigned")
	dt := strings.ToLower(col.DataType)
	switch dt {
	case "tinyint", "smallint", "mediumint", "int", "integer", "bigint", "boolean", "bool":
		n, ok := evalIntegerLiteralBigInt(kind, val)
		if !ok {
			return "", false
		}
		minB, maxB := intBoundsForType(dt, unsigned)
		if n.Cmp(minB) < 0 || n.Cmp(maxB) > 0 {
			return "", false
		}
		return n.String(), true
	case "decimal", "numeric", "dec", "fixed":
		return formatDecimalLiteralDefault(col.ColumnType, unsigned, kind, val, strict)
	case "float", "double", "real":
		if !validFloatLiteralDefault(dt, col.ColumnType, unsigned, kind, val, strict) {
			return "", false
		}
		pt := parseColumnType(col.ColumnType)
		hasMD := pt.hasLen && pt.hasScale && pt.scale >= 0
		isFloat32 := dt == "float" && !(pt.hasLen && !pt.hasScale && pt.length >= 25)
		var f float64
		if kind == litBit || kind == litHex {
			raw, ok := decodeLiteralRawBytes(val, kind == litHex)
			if !ok {
				return "", false
			}
			f, _ = new(big.Float).SetInt(new(big.Int).SetBytes(raw)).Float64()
		} else {
			var err error
			f, err = strconv.ParseFloat(strings.TrimSpace(val), 64)
			if err != nil {
				return "", false
			}
		}
		if hasMD {
			bitSize := 64
			if isFloat32 {
				bitSize = 32
				f = float64(float32(f))
			}
			formatted := strconv.FormatFloat(math.Abs(f), 'f', pt.scale, bitSize)
			if f < 0 && strings.Trim(formatted, "0.") != "" {
				return "-" + formatted, true
			}
			return formatted, true
		}
		if isFloat32 {
			return myGcvtFloat32(float32(f)), true
		}
		res, ok := myGcvtFormat(f, 325)
		return res, ok
	}
	return "", false
}

func myGcvtFloat32(f float32) string {
	if f == 0 {
		return "0"
	}
	raw := strconv.FormatFloat(float64(f), 'g', 6, 32)
	ePos := strings.IndexByte(raw, 'e')
	if ePos < 0 {
		return raw
	}
	expVal, _ := strconv.Atoi(raw[ePos+1:])
	return raw[:ePos] + "e" + strconv.Itoa(expVal)
}

func normalizeColumnStoredDefault(col *Column, expr nodes.ExprNode, strict bool) {
	if col == nil || expr == nil || col.Default == nil {
		return
	}
	if _, isParen := expr.(*nodes.ParenExpr); isParen {
		return
	}
	if sl, ok := expr.(*nodes.StringLit); ok {
		if normVal, changed := normalizeCharDefaultStringLiteral(col, sl.Value); changed {
			cp := *sl
			cp.Charset = ""
			cp.Value = normVal
			v := nodeToSQL(&cp)
			col.Default = &v
			return
		}
	}
	if coerced, ok := coerceNumericDefaultForCharColumn(col, expr); ok {
		v := nodeToSQL(&nodes.StringLit{Value: coerced})
		col.Default = &v
		return
	}
	if coerced, ok := coerceNumericColumnDefault(col, expr, strict); ok {
		col.Default = &coerced
	}
}

func decodeLiteralRawBytes(val string, isHex bool) ([]byte, bool) {
	if val == "" {
		return nil, true
	}
	if isHex {
		raw := val
		if len(raw)%2 == 1 {
			raw = "0" + raw
		}
		b, err := hex.DecodeString(raw)
		if err != nil {
			return nil, false
		}
		return b, true
	}
	pad := (8 - len(val)%8) % 8
	bitsStr := strings.Repeat("0", pad) + val
	out := make([]byte, len(bitsStr)/8)
	for i := 0; i < len(out); i++ {
		var v byte
		for j := 0; j < 8; j++ {
			ch := bitsStr[i*8+j]
			if ch != '0' && ch != '1' {
				return nil, false
			}
			v = (v << 1) | (ch - '0')
		}
		out[i] = v
	}
	return out, true
}

func validBinaryBytesForCharColumn(val string, isHex bool, charset string, maxLen int) bool {
	raw, ok := decodeLiteralRawBytes(val, isHex)
	if !ok {
		return false
	}
	cs := strings.ToLower(strings.TrimSpace(charset))
	switch cs {
	case "", "utf8mb4", "utf8", "utf8mb3":
		if !utf8.Valid(raw) {
			return false
		}
		s := string(raw)
		if cs == "utf8" || cs == "utf8mb3" {
			for _, r := range s {
				if r > 0xFFFF {
					return false
				}
			}
		}
		return utf8.RuneCountInString(strings.TrimRight(s, " ")) <= maxLen
	case "ascii":
		for _, b := range raw {
			if b >= 0x80 {
				return false
			}
		}
		return len(bytes.TrimRight(raw, " ")) <= maxLen
	case "ucs2", "utf16", "utf16le":
		if len(raw)%2 != 0 {
			return false
		}
		u16 := make([]uint16, len(raw)/2)
		for i := 0; i < len(u16); i++ {
			if cs == "utf16le" {
				u16[i] = uint16(raw[2*i]) | uint16(raw[2*i+1])<<8
			} else {
				u16[i] = uint16(raw[2*i])<<8 | uint16(raw[2*i+1])
			}
		}
		for len(u16) > 0 && u16[len(u16)-1] == 0x0020 {
			u16 = u16[:len(u16)-1]
		}
		if cs == "ucs2" {
			return len(u16) <= maxLen
		}
		runes := utf16.Decode(u16)
		for _, r := range runes {
			if r == utf8.RuneError {
				return false
			}
		}
		return len(runes) <= maxLen
	case "utf32":
		if len(raw)%4 != 0 {
			return false
		}
		n := len(raw) / 4
		runes := make([]rune, n)
		for i := 0; i < n; i++ {
			cp := uint32(raw[4*i])<<24 | uint32(raw[4*i+1])<<16 | uint32(raw[4*i+2])<<8 | uint32(raw[4*i+3])
			if cp > 0x10FFFF || (cp >= 0xD800 && cp <= 0xDFFF) {
				return false
			}
			runes[i] = rune(cp)
		}
		return utf8.RuneCountInString(strings.TrimRight(string(runes), " ")) <= maxLen
	default:
		return len(bytes.TrimRight(raw, " ")) <= maxLen
	}
}

// normalizeCharDefaultStringLiteral applies MySQL's trailing-space coercion when
// storing a quoted string default on CHAR or VARCHAR: CHAR strips all trailing
// U+0020 spaces, and VARCHAR strips only trailing U+0020 spaces past maxLen.
func normalizeCharDefaultStringLiteral(col *Column, val string) (string, bool) {
	if col == nil {
		return val, false
	}
	switch strings.ToLower(col.DataType) {
	case "char":
		trimmed := strings.TrimRight(val, " ")
		return trimmed, trimmed != val
	case "varchar":
		maxLen, ok := fixedLengthStringColumnLimit(col)
		if !ok {
			return val, false
		}
		runes := []rune(val)
		if len(runes) <= maxLen {
			return val, false
		}
		for len(runes) > maxLen && runes[len(runes)-1] == ' ' {
			runes = runes[:len(runes)-1]
		}
		trimmed := string(runes)
		return trimmed, trimmed != val
	default:
		return val, false
	}
}

func leadingNumericPrefix(s string) string {
	i := 0
	if i < len(s) && (s[i] == '+' || s[i] == '-') {
		i++
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
		return "0"
	}
	if i < len(s) && (s[i] == 'e' || s[i] == 'E') {
		j := i + 1
		if j < len(s) && (s[j] == '+' || s[j] == '-') {
			j++
		}
		expDigits := 0
		for j < len(s) && s[j] >= '0' && s[j] <= '9' {
			expDigits++
			j++
		}
		if expDigits > 0 {
			i = j
		}
	}
	return s[:i]
}

func validDecimalLiteralDefault(colType string, unsigned bool, kind literalDefaultKind, val string, strict bool) bool {
	pt := parseColumnType(colType)
	m, d := 10, 0
	if pt.hasLen {
		m = pt.length
	}
	if pt.hasScale && pt.scale >= 0 {
		d = pt.scale
	}
	if kind == litBit || kind == litHex {
		raw, ok := decodeLiteralRawBytes(val, kind == litHex)
		if !ok {
			return false
		}
		n := new(big.Int).SetBytes(raw)
		intDigits := m - d
		if intDigits < 0 {
			intDigits = 0
		}
		maxInt := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(intDigits)), nil)
		return n.Cmp(maxInt) < 0
	}
	s := strings.TrimSpace(val)
	if !isStrictNumberSyntax(s) {
		if strict {
			return false
		}
		s = leadingNumericPrefix(s)
	}
	r, ok := new(big.Rat).SetString(s)
	if !ok {
		return false
	}
	if unsigned && r.Sign() < 0 {
		return false
	}
	scalePow := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(d)), nil)
	scaled := new(big.Rat).Mul(new(big.Rat).Abs(r), new(big.Rat).SetInt(scalePow))
	roundedScaled := roundRatHalfAwayFromZero(scaled)
	maxDigits := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(m)), nil)
	return roundedScaled.Cmp(maxDigits) < 0
}

func validFloatLiteralDefault(dt, colType string, unsigned bool, kind literalDefaultKind, val string, strict bool) bool {
	pt := parseColumnType(colType)
	hasMD := pt.hasLen && pt.hasScale && pt.scale >= 0
	if kind == litBit || kind == litHex {
		if !hasMD {
			return true
		}
		raw, ok := decodeLiteralRawBytes(val, kind == litHex)
		if !ok {
			return false
		}
		n := new(big.Int).SetBytes(raw)
		intDigits := pt.length - pt.scale
		if intDigits < 0 {
			intDigits = 0
		}
		maxInt := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(intDigits)), nil)
		return n.Cmp(maxInt) < 0
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
	isFloat32 := dt == "float" && !(pt.hasLen && !pt.hasScale && pt.length >= 25)
	if isFloat32 && math.Abs(f) > math.MaxFloat32 {
		return false
	}
	if hasMD {
		return validDecimalLiteralDefault(colType, unsigned, kind, val, strict)
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

func coercedNumericStringForCharLen(val string, maxLen int) (string, bool) {
	if maxLen <= 0 {
		return "", false
	}
	s := strings.TrimSpace(val)
	if !isStrictNumberSyntax(s) {
		if len(s) <= maxLen {
			return s, true
		}
		return "", false
	}
	if strings.ContainsAny(s, "eE") {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil || math.IsInf(f, 0) || math.IsNaN(f) {
			return "", false
		}
		return myGcvtFormat(f, maxLen)
	}
	exact := coercedExactNumericString(s)
	if len(exact) > maxLen {
		return "", false
	}
	return exact, true
}

func validNumericDefaultForStringLen(val string, maxLen int) bool {
	_, ok := coercedNumericStringForCharLen(val, maxLen)
	return ok
}

// coerceNumericDefaultForCharColumn returns the string MySQL stores in the data
// dictionary when an unparenthesized numeric literal default is declared on a
// CHAR, VARCHAR, BINARY, or VARBINARY column (e.g. `CHAR(1) DEFAULT 01` or
// `BINARY(1) DEFAULT 01` -> `'1'`, `VARBINARY(2) DEFAULT 007` -> `'7'`).
func coerceNumericDefaultForCharColumn(col *Column, expr nodes.ExprNode) (string, bool) {
	if col == nil {
		return "", false
	}
	dt := strings.ToLower(col.DataType)
	if dt != "char" && dt != "varchar" && dt != "binary" && dt != "varbinary" {
		return "", false
	}
	kind, val, unquotedNum := extractDefaultLiteral(expr)
	if kind != litNumber || !unquotedNum {
		return "", false
	}
	maxLen, ok := fixedLengthStringColumnLimit(col)
	if !ok {
		maxLen = 325
	}
	coerced, ok := coercedNumericStringForCharLen(val, maxLen)
	if !ok {
		return "", false
	}
	if norm, changed := normalizeCharDefaultStringLiteral(col, coerced); changed {
		return norm, true
	}
	return coerced, true
}

// coercedExactNumericString normalizes an exact integer or fixed-point decimal
// literal to the string Item_int / Item_decimal::val_str produces in MySQL
// (e.g. `01` -> `"1"`, `-0` -> `"0"`, `+1` -> `"1"`, `.5` -> `"0.5"`, `1.50` -> `"1.50"`).
func coercedExactNumericString(s string) string {
	neg := false
	if strings.HasPrefix(s, "-") {
		neg = true
		s = s[1:]
	} else if strings.HasPrefix(s, "+") {
		s = s[1:]
	}
	dot := strings.IndexByte(s, '.')
	if dot < 0 {
		intPart := strings.TrimLeft(s, "0")
		if intPart == "" {
			return "0"
		}
		if neg {
			return "-" + intPart
		}
		return intPart
	}
	intPart := strings.TrimLeft(s[:dot], "0")
	fracPart := s[dot+1:]
	allZero := intPart == "" && strings.TrimLeft(fracPart, "0") == ""
	if intPart == "" {
		intPart = "0"
	}
	if fracPart == "" {
		if allZero || !neg {
			return intPart
		}
		return "-" + intPart
	}
	if neg && !allZero {
		return "-" + intPart + "." + fracPart
	}
	return intPart + "." + fracPart
}

func allNinePrefix(digits string, n int) bool {
	for i := 0; i < n && i < len(digits); i++ {
		if digits[i] != '9' {
			return false
		}
	}
	return true
}

func formatGcvtFixed(av float64, neg bool, fracDigits int) string {
	raw := strconv.FormatFloat(av, 'f', fracDigits, 64)
	if strings.IndexByte(raw, '.') >= 0 {
		raw = strings.TrimRight(raw, "0")
		raw = strings.TrimSuffix(raw, ".")
	}
	if raw == "0" || !neg {
		return raw
	}
	return "-" + raw
}

func formatGcvtSci(av float64, neg bool, ndigits int) string {
	raw := strconv.FormatFloat(av, 'e', ndigits-1, 64)
	ePos := strings.IndexByte(raw, 'e')
	mPart := raw[:ePos]
	if strings.IndexByte(mPart, '.') >= 0 {
		mPart = strings.TrimRight(mPart, "0")
		mPart = strings.TrimSuffix(mPart, ".")
	}
	expVal, _ := strconv.Atoi(raw[ePos+1:])
	res := mPart + "e" + strconv.Itoa(expVal)
	if neg {
		return "-" + res
	}
	return res
}

// myGcvtFormat reproduces MySQL's my_gcvt(f, MY_GCVT_ARG_DOUBLE, width, ...)
// string output and reports whether f fits within width bytes without overflow.
func myGcvtFormat(f float64, width int) (string, bool) {
	if width <= 0 {
		return "", false
	}
	if f == 0 {
		return "0", true
	}
	neg := f < 0
	sign := 0
	av := f
	if neg {
		sign = 1
		av = -f
	}
	avail := width - sign
	if avail <= 0 {
		return "", false
	}
	sci := strconv.FormatFloat(av, 'e', -1, 64)
	ePos := strings.IndexByte(sci, 'e')
	mant := strings.ReplaceAll(sci[:ePos], ".", "")
	mant = strings.TrimRight(mant, "0")
	if mant == "" {
		mant = "0"
	}
	exp, _ := strconv.Atoi(sci[ePos+1:])
	decpt := exp + 1

	if decpt >= 1 {
		if decpt < avail {
			fracDigits := avail - decpt - 1
			if fracDigits < 0 {
				fracDigits = 0
			}
			return formatGcvtFixed(av, neg, fracDigits), true
		}
		if decpt == avail {
			if len(mant) <= decpt || !(allNinePrefix(mant, decpt) && mant[decpt] >= '5') {
				return formatGcvtFixed(av, neg, 0), true
			}
		}
	} else if decpt >= -4 && avail >= 3-decpt {
		return formatGcvtFixed(av, neg, avail-2), true
	}

	// Scientific notation path (`d[.ddd]e[+-]exp`). When len(mant) > 1, my_gcvt
	// reserves 1 byte for '.' when computing available significant digits.
	expLen := 1 + len(strconv.Itoa(exp))
	dot := 0
	if len(mant) > 1 {
		dot = 1
	}
	ndigits := avail - expLen - dot
	if ndigits >= 1 {
		if !(ndigits < len(mant) && allNinePrefix(mant, ndigits) && mant[ndigits] >= '5' && avail-(1+len(strconv.Itoa(exp+1)))-dot < 1) {
			return formatGcvtSci(av, neg, ndigits), true
		}
	}

	// Sub-unit fixed fallback ('0') when decpt in {-2, -1, 0} and avail >= 2.
	if decpt <= 0 && decpt >= -2 && avail >= 2 {
		return formatGcvtFixed(av, neg, avail-2), true
	}
	return "", false
}

func coercedDefaultLiteralString(s string) (string, bool) {
	if !isStrictNumberSyntax(s) {
		return "", false
	}
	if strings.ContainsAny(s, "eE") {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil || math.IsInf(f, 0) || math.IsNaN(f) {
			return "", false
		}
		return strconv.FormatFloat(f, 'g', -1, 64), true
	}
	return coercedExactNumericString(s), true
}

// parseMySQLIntOrUintLiteral reports whether MySQL's lexer produces Item_int or
// Item_uint for s (an integer token with no '.' or exponent fitting in int64 or
// uint64), returning its sign and uint64 magnitude when non-negative.
func parseMySQLIntOrUintLiteral(s string) (isIntToken bool, nonNeg bool, mag uint64) {
	if strings.ContainsAny(s, ".eE") {
		return false, false, 0
	}
	if strings.HasPrefix(s, "-") {
		digits := strings.TrimLeft(s[1:], "0")
		if digits == "" {
			return true, true, 0
		}
		if _, err := strconv.ParseInt(s, 10, 64); err == nil {
			return true, false, 0
		}
		// Negative integer below MinInt64 is lexed as DECIMAL_NUM by MySQL.
		return false, false, 0
	}
	u, err := strconv.ParseUint(strings.TrimPrefix(s, "+"), 10, 64)
	if err != nil {
		// Positive integer above MaxUint64 is lexed as DECIMAL_NUM by MySQL.
		return false, false, 0
	}
	return true, true, u
}

func enumContainsStringMember(members []string, val string, col *Column, strict bool) bool {
	target := collationComparisonKey(strings.TrimRight(val, " "), col.Collation, col.Charset)
	for _, m := range members {
		if collationComparisonKey(strings.TrimRight(m, " "), col.Collation, col.Charset) == target {
			return true
		}
	}
	if target == "" && !strict {
		return true
	}
	return false
}

func setContainsStringMembers(members []string, val string, col *Column) bool {
	if val == "" {
		return true
	}
	allowed := make(map[string]bool, len(members))
	for _, m := range members {
		allowed[collationComparisonKey(strings.TrimRight(m, " "), col.Collation, col.Charset)] = true
	}
	for _, part := range strings.Split(val, ",") {
		if !allowed[collationComparisonKey(strings.TrimRight(part, " "), col.Collation, col.Charset)] {
			return false
		}
	}
	return true
}

func validEnumLiteralDefault(col *Column, kind literalDefaultKind, val string, unquotedNum, strict bool) bool {
	members := enumMembersFromColumnType(col.ColumnType)
	if unquotedNum {
		s := strings.TrimSpace(val)
		coerced, ok := coercedDefaultLiteralString(s)
		if !ok {
			return false
		}
		if isIntTok, nonNeg, u := parseMySQLIntOrUintLiteral(s); isIntTok {
			if !nonNeg || u > uint64(len(members)) {
				return false
			}
			if u == 0 && strict {
				return false
			}
		}
		return enumContainsStringMember(members, coerced, col, strict)
	}
	if kind != litString {
		return true
	}
	return enumContainsStringMember(members, val, col, strict)
}

func validSetLiteralDefault(col *Column, kind literalDefaultKind, val string, unquotedNum bool) bool {
	members := enumMembersFromColumnType(col.ColumnType)
	if unquotedNum {
		s := strings.TrimSpace(val)
		coerced, ok := coercedDefaultLiteralString(s)
		if !ok {
			return false
		}
		if isIntTok, nonNeg, u := parseMySQLIntOrUintLiteral(s); isIntTok {
			if !nonNeg {
				return false
			}
			if len(members) < 64 && u >= (uint64(1)<<len(members)) {
				return false
			}
		}
		return setContainsStringMembers(members, coerced, col)
	}
	if kind != litString {
		return true
	}
	return setContainsStringMembers(members, val, col)
}

// validateAndNormalizeSessionTimeZone validates a SET time_zone value per MySQL
// 8.0.46 (`SYSTEM`, `[+-]HH:MM` in `[-13:59, +14:00]`, or a known IANA zone).
func validateAndNormalizeSessionTimeZone(raw string) (string, error) {
	tz := strings.TrimSpace(raw)
	if tz == "" || strings.EqualFold(tz, "Local") {
		return "", errUnknownTimeZone(raw)
	}
	if strings.EqualFold(tz, "SYSTEM") {
		return "SYSTEM", nil
	}
	if tz[0] == '+' || tz[0] == '-' {
		if len(tz) != 6 || !isASCIIDigit(tz[1]) || !isASCIIDigit(tz[2]) ||
			tz[3] != ':' || !isASCIIDigit(tz[4]) || !isASCIIDigit(tz[5]) {
			return "", errUnknownTimeZone(raw)
		}
		hh := int(tz[1]-'0')*10 + int(tz[2]-'0')
		mm := int(tz[4]-'0')*10 + int(tz[5]-'0')
		total := hh*60 + mm
		if mm > 59 {
			return "", errUnknownTimeZone(raw)
		}
		if tz[0] == '-' {
			if total == 0 {
				return "+00:00", nil
			}
			if total > 13*60+59 {
				return "", errUnknownTimeZone(raw)
			}
			return tz, nil
		}
		if total > 14*60 {
			return "", errUnknownTimeZone(raw)
		}
		return tz, nil
	}
	if _, err := time.LoadLocation(tz); err != nil {
		return "", errUnknownTimeZone(raw)
	}
	return tz, nil
}

func convertSessionLocalToUTC(year, month, day, hour, min, sec int, sessionTZ string) (time.Time, bool) {
	tz := strings.TrimSpace(sessionTZ)
	if tz == "" || strings.EqualFold(tz, "SYSTEM") || tz == "+00:00" || tz == "-00:00" || strings.EqualFold(tz, "UTC") {
		return time.Date(year, time.Month(month), day, hour, min, sec, 0, time.UTC), true
	}
	if len(tz) == 6 && (tz[0] == '+' || tz[0] == '-') && tz[3] == ':' {
		hh := int(tz[1]-'0')*10 + int(tz[2]-'0')
		mm := int(tz[4]-'0')*10 + int(tz[5]-'0')
		off := hh*60 + mm
		if tz[0] == '-' {
			off = -off
		}
		return time.Date(year, time.Month(month), day, hour, min, sec, 0, time.UTC).Add(-time.Duration(off) * time.Minute), true
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return time.Date(year, time.Month(month), day, hour, min, sec, 0, time.UTC), true
	}
	localT := time.Date(year, time.Month(month), day, hour, min, sec, 0, loc)
	// Reject spring-forward DST gap times that do not exist in loc.
	if localT.Year() != year || localT.Month() != time.Month(month) || localT.Day() != day ||
		localT.Hour() != hour || localT.Minute() != min || localT.Second() != sec {
		return time.Time{}, false
	}
	return localT.UTC(), true
}

func validDateOrDatetimeLiteralDefault(dt string, kind literalDefaultKind, val string, unquotedNum, noZeroDate, noZeroInDate, allowInvalidDates bool, sessionTZ string) bool {
	if kind != litString && kind != litNumber {
		return false
	}
	s := strings.TrimSpace(val)
	if s == "" {
		return false
	}
	year, month, day, hour, min, sec, hasTZ, tzOffsetMin, ok := parseMySQLDateTimeParts(s, unquotedNum)
	if !ok {
		return false
	}
	allZero := year == 0 && month == 0 && day == 0 && hour == 0 && min == 0 && sec == 0
	if allZero {
		if hasTZ {
			return false
		}
		return !noZeroDate
	}
	if month > 12 || day > 31 || hour > 23 || min > 59 || sec > 59 {
		return false
	}
	if month == 0 || day == 0 {
		if dt == "timestamp" || hasTZ || noZeroInDate {
			return false
		}
		return true
	}
	if (dt == "timestamp" || hasTZ || !allowInvalidDates) && day > daysInMonth(year, month) {
		return false
	}
	minUTC := time.Date(1970, time.January, 1, 0, 0, 1, 0, time.UTC)
	maxUTC := time.Date(2038, time.January, 19, 3, 14, 7, 0, time.UTC)
	if hasTZ {
		utc := time.Date(year, time.Month(month), day, hour, min, sec, 0, time.UTC).Add(-time.Duration(tzOffsetMin) * time.Minute)
		if utc.Year() < 0 || utc.Year() > 9999 {
			return false
		}
		if dt == "timestamp" && (utc.Before(minUTC) || utc.After(maxUTC)) {
			return false
		}
		return true
	}
	if dt == "timestamp" {
		utc, validLocal := convertSessionLocalToUTC(year, month, day, hour, min, sec, sessionTZ)
		if !validLocal || utc.Before(minUTC) || utc.After(maxUTC) {
			return false
		}
	}
	return true
}

func isASCIIDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func countDigitRuns(s string) int {
	n := 0
	for i := 0; i < len(s); {
		if isASCIIDigit(s[i]) {
			n++
			for i < len(s) && isASCIIDigit(s[i]) {
				i++
			}
		} else {
			i++
		}
	}
	return n
}

func stripTrailingFraction(s string) string {
	dot := strings.LastIndexByte(s, '.')
	if dot < 0 {
		return s
	}
	frac := s[dot+1:]
	if len(frac) == 0 {
		return s
	}
	for i := 0; i < len(frac); i++ {
		if !isASCIIDigit(frac[i]) {
			return s
		}
	}
	if strings.ContainsAny(s[:dot], " Tt:") || len(s[:dot]) >= 12 {
		return s[:dot]
	}
	return s
}

func parseMySQLDateTimeParts(s string, unquotedNum bool) (year, month, day, hour, min, sec int, hasTZ bool, tzOffsetMin int, ok bool) {
	if unquotedNum && s == "0" {
		return 0, 0, 0, 0, 0, 0, false, 0, true
	}
	base := s
	// Check for an attached timezone suffix `[+-]HH:MM` immediately following a
	// full date-and-time value (no whitespace before the sign). MySQL 8.0.46
	// accepts `-14:00` through `+14:00` and rejects `-00:00`.
	if !unquotedNum && len(base) >= 7 {
		tz := base[len(base)-6:]
		if (tz[0] == '+' || tz[0] == '-') && isASCIIDigit(tz[1]) && isASCIIDigit(tz[2]) &&
			tz[3] == ':' && isASCIIDigit(tz[4]) && isASCIIDigit(tz[5]) && isASCIIDigit(base[len(base)-7]) {
			prefixNoFrac := stripTrailingFraction(base[:len(base)-6])
			allPrefixDigits := len(prefixNoFrac) > 0
			for i := 0; i < len(prefixNoFrac); i++ {
				if !isASCIIDigit(prefixNoFrac[i]) {
					allPrefixDigits = false
					break
				}
			}
			isFullDateTime := (allPrefixDigits && (len(prefixNoFrac) == 12 || len(prefixNoFrac) == 14)) ||
				(!allPrefixDigits && countDigitRuns(prefixNoFrac) == 6)
			if isFullDateTime {
				tzH := int(tz[1]-'0')*10 + int(tz[2]-'0')
				tzM := int(tz[4]-'0')*10 + int(tz[5]-'0')
				if tz[0] == '-' && tzH == 0 && tzM == 0 {
					return 0, 0, 0, 0, 0, 0, false, 0, false
				}
				if tzM > 59 || tzH*60+tzM > 14*60 {
					return 0, 0, 0, 0, 0, 0, false, 0, false
				}
				hasTZ = true
				tzOffsetMin = tzH*60 + tzM
				if tz[0] == '-' {
					tzOffsetMin = -tzOffsetMin
				}
				base = prefixNoFrac
			}
		}
	}
	if !hasTZ {
		base = stripTrailingFraction(base)
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
			return expandTwoDigitYear(y, m, d), m, d, 0, 0, 0, hasTZ, tzOffsetMin, true
		case 8:
			y, _ := strconv.Atoi(base[0:4])
			m, _ := strconv.Atoi(base[4:6])
			d, _ := strconv.Atoi(base[6:8])
			return y, m, d, 0, 0, 0, hasTZ, tzOffsetMin, true
		case 12:
			y, _ := strconv.Atoi(base[0:2])
			m, _ := strconv.Atoi(base[2:4])
			d, _ := strconv.Atoi(base[4:6])
			hh, _ := strconv.Atoi(base[6:8])
			mm, _ := strconv.Atoi(base[8:10])
			ss, _ := strconv.Atoi(base[10:12])
			return expandTwoDigitYear(y, m, d), m, d, hh, mm, ss, hasTZ, tzOffsetMin, true
		case 14:
			y, _ := strconv.Atoi(base[0:4])
			m, _ := strconv.Atoi(base[4:6])
			d, _ := strconv.Atoi(base[6:8])
			hh, _ := strconv.Atoi(base[8:10])
			mm, _ := strconv.Atoi(base[10:12])
			ss, _ := strconv.Atoi(base[12:14])
			return y, m, d, hh, mm, ss, hasTZ, tzOffsetMin, true
		default:
			return 0, 0, 0, 0, 0, 0, false, 0, false
		}
	}

	// Tokenize digit runs separated by punctuation / 'T' / space.
	// Reject any alphabetic characters other than a single 'T'/'t' between date and time.
	var parts []string
	i := 0
	for i < len(base) {
		ch := base[i]
		if ch >= '0' && ch <= '9' {
			j := i
			for j < len(base) && base[j] >= '0' && base[j] <= '9' {
				j++
			}
			parts = append(parts, base[i:j])
			i = j
			continue
		}
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') {
			if (ch == 'T' || ch == 't') && len(parts) == 3 {
				i++
				continue
			}
			return 0, 0, 0, 0, 0, 0, false, 0, false
		}
		i++
	}
	if len(parts) < 3 || len(parts) > 6 {
		return 0, 0, 0, 0, 0, 0, false, 0, false
	}
	y, err1 := strconv.Atoi(parts[0])
	m, err2 := strconv.Atoi(parts[1])
	d, err3 := strconv.Atoi(parts[2])
	if err1 != nil || err2 != nil || err3 != nil || len(parts[0]) > 4 {
		return 0, 0, 0, 0, 0, 0, false, 0, false
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
	return y, m, d, hour, min, sec, hasTZ, tzOffsetMin, true
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
