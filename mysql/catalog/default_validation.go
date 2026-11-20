package catalog

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"math/big"
	"math/bits"
	"strconv"
	"strings"
	"time"
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
	truncateFrac := c.isTimeTruncateFractionalSQLMode()
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
		if !validIntegerLiteralDefault(dt, unsigned, kind, val, strict) {
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
		if kind == litString && !validStringForColumnCharset(val, col.Charset) {
			return errInvalidDefault(col.Name)
		}
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
		if !validBitLiteralDefault(col.ColumnType, kind, val, unquotedNum, strict) {
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
		if !validDateOrDatetimeLiteralDefault(dt, temporalColumnFSP(col.ColumnType), kind, val, unquotedNum, noZeroDate, noZeroInDate, allowInvalidDates, truncateFrac, sessionTZ) {
			return errInvalidDefault(col.Name)
		}
	case "time":
		if !validTimeLiteralDefault(temporalColumnFSP(col.ColumnType), kind, val, unquotedNum, truncateFrac) {
			return errInvalidDefault(col.Name)
		}
	}
	normalizeColumnStoredDefault(col, expr, strict, truncateFrac, sessionTZ)
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

func evalBinaryNumericUint64(isHex bool, val string, unsigned, strict bool) (uint64, bool) {
	raw, ok := decodeLiteralRawBytes(val, isHex)
	if !ok {
		return 0, false
	}
	if len(raw) > 8 {
		if strict {
			return 0, false
		}
		if unsigned {
			return math.MaxUint64, true
		}
		return math.MaxInt64, true
	}
	var u uint64
	for _, b := range raw {
		u = (u << 8) | uint64(b)
	}
	if !unsigned && u > math.MaxInt64 {
		if strict {
			return 0, false
		}
		return math.MaxInt64, true
	}
	return u, true
}

func evalIntegerLiteralBigInt(kind literalDefaultKind, val string, unsigned, strict bool) (*big.Int, bool) {
	switch kind {
	case litBit, litHex:
		u, ok := evalBinaryNumericUint64(kind == litHex, val, unsigned, strict)
		if !ok {
			return nil, false
		}
		return new(big.Int).SetUint64(u), true
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

func validIntegerLiteralDefault(dt string, unsigned bool, kind literalDefaultKind, val string, strict bool) bool {
	minB, maxB := intBoundsForType(dt, unsigned)
	n, ok := evalIntegerLiteralBigInt(kind, val, unsigned, strict)
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
		u, ok := evalBinaryNumericUint64(kind == litHex, val, unsigned, strict)
		if !ok {
			return "", false
		}
		n := new(big.Int).SetUint64(u)
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

func coerceNumericColumnDefault(col *Column, expr nodes.ExprNode, strict, truncateFrac bool, sessionTZ string) (string, bool) {
	if col == nil || expr == nil {
		return "", false
	}
	if _, isParen := expr.(*nodes.ParenExpr); isParen {
		return "", false
	}
	if isNullDefault(expr) {
		return "", false
	}
	kind, val, unquotedNum := extractDefaultLiteral(expr)
	if kind == litNone {
		return "", false
	}
	unsigned := strings.Contains(strings.ToLower(col.ColumnType), "unsigned")
	dt := strings.ToLower(col.DataType)
	switch dt {
	case "date", "datetime", "timestamp":
		return coerceDateOrDatetimeDefault(dt, temporalColumnFSP(col.ColumnType), kind, val, unquotedNum, truncateFrac, sessionTZ)
	case "time":
		return coerceTimeDefault(temporalColumnFSP(col.ColumnType), kind, val, unquotedNum, truncateFrac)
	case "year":
		if kind == litHex || kind == litBit {
			return coerceBinaryYearDefault(kind == litHex, val)
		}
		return "", false
	case "enum":
		if kind == litHex || kind == litBit {
			return coerceBinaryEnumDefault(col, kind == litHex, val, strict)
		}
		return "", false
	case "set":
		if kind == litHex || kind == litBit {
			return coerceBinarySetDefault(col, kind == litHex, val)
		}
		return "", false
	case "bit":
		if !strict && (kind == litBit || kind == litHex) {
			if raw, ok := decodeLiteralRawBytes(val, kind == litHex); ok && len(raw) > 8 {
				return "b'" + strings.Repeat("1", 64) + "'", true
			}
		}
		return "", false
	case "tinyint", "smallint", "mediumint", "int", "integer", "bigint", "boolean", "bool":
		n, ok := evalIntegerLiteralBigInt(kind, val, unsigned, strict)
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
			u, ok := evalBinaryNumericUint64(kind == litHex, val, unsigned, strict)
			if !ok {
				return "", false
			}
			f = float64(u)
		} else {
			var err error
			f, err = strconv.ParseFloat(strings.TrimSpace(val), 64)
			if err != nil || math.IsInf(f, 0) || math.IsNaN(f) {
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
	if math.IsInf(float64(f), 0) || math.IsNaN(float64(f)) {
		return ""
	}
	if f == 0 {
		return "0"
	}
	neg := f < 0
	avail := 12
	av := float64(f)
	if neg {
		avail = 11
		av = -av
	}
	sci := strconv.FormatFloat(av, 'e', 5, 32)
	ePos := strings.IndexByte(sci, 'e')
	if ePos < 0 {
		return ""
	}
	mant := strings.ReplaceAll(sci[:ePos], ".", "")
	mant = strings.TrimRight(mant, "0")
	if mant == "" {
		mant = "0"
	}
	exp, _ := strconv.Atoi(sci[ePos+1:])
	decpt := exp + 1
	var res string
	if decpt >= 1 && decpt <= avail {
		if len(mant) <= decpt {
			res = mant + strings.Repeat("0", decpt-len(mant))
		} else {
			res = mant[:decpt] + "." + mant[decpt:]
		}
	} else if decpt <= 0 && 2+(-decpt)+len(mant) <= avail {
		res = "0." + strings.Repeat("0", -decpt) + mant
	} else if len(mant) == 1 {
		res = mant + "e" + strconv.Itoa(exp)
	} else {
		res = mant[:1] + "." + mant[1:] + "e" + strconv.Itoa(exp)
	}
	if neg {
		return "-" + res
	}
	return res
}

func normalizeColumnStoredDefault(col *Column, expr nodes.ExprNode, strict, truncateFrac bool, sessionTZ string) {
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
	if coerced, ok := coerceNumericColumnDefault(col, expr, strict, truncateFrac, sessionTZ); ok {
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
		// MySQL left-pads binary string literals with 0x00 bytes to a multiple of
		// mbminlen (2 bytes for ucs2/utf16/utf16le).
		if len(raw)%2 != 0 {
			raw = append([]byte{0}, raw...)
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
		// Validate UTF-16 surrogate pairing directly rather than via utf16.Decode,
		// which maps both unpaired surrogates and a valid U+FFFD to utf8.RuneError.
		charCount := 0
		for i := 0; i < len(u16); i++ {
			w := u16[i]
			if w >= 0xD800 && w <= 0xDBFF {
				if i+1 >= len(u16) || u16[i+1] < 0xDC00 || u16[i+1] > 0xDFFF {
					return false
				}
				i++
			} else if w >= 0xDC00 && w <= 0xDFFF {
				return false
			}
			charCount++
		}
		return charCount <= maxLen
	case "utf32":
		// MySQL left-pads binary string literals with 0x00 bytes to a multiple of
		// mbminlen (4 bytes for utf32).
		if rem := len(raw) % 4; rem != 0 {
			raw = append(make([]byte, 4-rem), raw...)
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

func validStringForColumnCharset(val, charset string) bool {
	if !utf8.ValidString(val) {
		return false
	}
	switch strings.ToLower(strings.TrimSpace(charset)) {
	case "ascii":
		for i := 0; i < len(val); i++ {
			if val[i] >= 0x80 {
				return false
			}
		}
	case "utf8", "utf8mb3", "ucs2":
		for _, r := range val {
			if r > 0xFFFF {
				return false
			}
		}
	case "latin1":
		for _, r := range val {
			if !isMySQLLatin1Rune(r) {
				return false
			}
		}
	}
	return true
}

func isMySQLLatin1Rune(r rune) bool {
	if r >= 0 && r <= 0xFF {
		return true
	}
	switch r {
	case 0x20AC, 0x201A, 0x0192, 0x201E, 0x2026, 0x2020, 0x2021, 0x02C6,
		0x2030, 0x0160, 0x2039, 0x0152, 0x017D, 0x2018, 0x2019, 0x201C,
		0x201D, 0x2022, 0x2013, 0x2014, 0x02DC, 0x2122, 0x0161, 0x203A,
		0x0153, 0x017E, 0x0178:
		return true
	}
	return false
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
		u, ok := evalBinaryNumericUint64(kind == litHex, val, unsigned, strict)
		if !ok {
			return false
		}
		n := new(big.Int).SetUint64(u)
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
	isFloat32 := dt == "float" && !(pt.hasLen && !pt.hasScale && pt.length >= 25)
	if kind == litBit || kind == litHex {
		u, ok := evalBinaryNumericUint64(kind == litHex, val, unsigned, strict)
		if !ok {
			return false
		}
		f := float64(u)
		if math.IsInf(f, 0) || math.IsNaN(f) {
			return false
		}
		if isFloat32 && math.IsInf(float64(float32(f)), 0) {
			return false
		}
		if hasMD {
			return validDecimalLiteralDefault(colType, unsigned, litNumber, strconv.FormatUint(u, 10), strict)
		}
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

func validBitLiteralDefault(colType string, kind literalDefaultKind, val string, unquotedNum, strict bool) bool {
	width := bitColumnWidth(colType)
	switch kind {
	case litBit, litHex:
		u, ok := evalBinaryNumericUint64(kind == litHex, val, true, strict)
		if !ok {
			return false
		}
		return bits.Len64(u) <= width
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

func evalBinaryYearUint(isHex bool, val string) (uint64, bool) {
	raw, ok := decodeLiteralRawBytes(val, isHex)
	if !ok || len(raw) > 8 {
		return 0, false
	}
	var u uint64
	for _, b := range raw {
		u = (u << 8) | uint64(b)
	}
	return u, true
}

func validYearLiteralDefault(kind literalDefaultKind, val string, unquotedNum bool) bool {
	if kind == litBit || kind == litHex {
		u, ok := evalBinaryYearUint(kind == litHex, val)
		if !ok {
			return false
		}
		return u == 0 || (u >= 1 && u <= 99) || (u >= 1901 && u <= 2155)
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

func coerceBinaryYearDefault(isHex bool, val string) (string, bool) {
	u, ok := evalBinaryYearUint(isHex, val)
	if !ok {
		return "", false
	}
	switch {
	case u == 0:
		return "'0000'", true
	case u >= 1 && u <= 69:
		return fmt.Sprintf("'%04d'", 2000+u), true
	case u >= 70 && u <= 99:
		return fmt.Sprintf("'%04d'", 1900+u), true
	case u >= 1901 && u <= 2155:
		return fmt.Sprintf("'%04d'", u), true
	default:
		return "", false
	}
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
	if ePos < 0 {
		return raw
	}
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
	if width <= 0 || math.IsInf(f, 0) || math.IsNaN(f) {
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
	if ePos < 0 {
		return "", false
	}
	mant := strings.ReplaceAll(sci[:ePos], ".", "")
	mant = strings.TrimRight(mant, "0")
	if mant == "" {
		mant = "0"
	}
	exp, _ := strconv.Atoi(sci[ePos+1:])
	decpt := exp + 1

	if decpt >= 1 && decpt <= 15 {
		if decpt < avail {
			if len(mant) <= decpt {
				res := mant + strings.Repeat("0", decpt-len(mant))
				if neg {
					return "-" + res, true
				}
				return res, true
			}
			availFrac := avail - decpt - 1
			if availFrac >= len(mant)-decpt {
				res := mant[:decpt] + "." + mant[decpt:]
				if neg {
					return "-" + res, true
				}
				return res, true
			}
			if availFrac < 0 {
				availFrac = 0
			}
			return formatGcvtFixed(av, neg, availFrac), true
		}
		if decpt == avail {
			if len(mant) <= decpt {
				res := mant + strings.Repeat("0", decpt-len(mant))
				if neg {
					return "-" + res, true
				}
				return res, true
			}
			if !(allNinePrefix(mant, decpt) && mant[decpt] >= '5') {
				return formatGcvtFixed(av, neg, 0), true
			}
		}
	} else if decpt <= 0 && decpt >= -5 && avail >= 3-decpt {
		if avail-2 >= len(mant)-decpt {
			res := "0." + strings.Repeat("0", -decpt) + mant
			if neg {
				return "-" + res, true
			}
			return res, true
		}
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
		if ndigits >= len(mant) {
			var res string
			if len(mant) == 1 {
				res = mant + "e" + strconv.Itoa(exp)
			} else {
				res = mant[:1] + "." + mant[1:] + "e" + strconv.Itoa(exp)
			}
			if neg {
				return "-" + res, true
			}
			return res, true
		}
		if !(allNinePrefix(mant, ndigits) && mant[ndigits] >= '5' && avail-(1+len(strconv.Itoa(exp+1)))-dot < 1) {
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
	if kind == litHex || kind == litBit {
		raw, ok := decodeLiteralRawBytes(val, kind == litHex)
		if !ok || !utf8.Valid(raw) {
			return false
		}
		return enumContainsStringMember(members, string(raw), col, strict)
	}
	if kind != litString {
		return false
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
	if kind == litHex || kind == litBit {
		raw, ok := decodeLiteralRawBytes(val, kind == litHex)
		if !ok || !utf8.Valid(raw) {
			return false
		}
		return setContainsStringMembers(members, string(raw), col)
	}
	if kind != litString {
		return false
	}
	return setContainsStringMembers(members, val, col)
}

func coerceBinaryEnumDefault(col *Column, isHex bool, val string, strict bool) (string, bool) {
	raw, ok := decodeLiteralRawBytes(val, isHex)
	if !ok || !utf8.Valid(raw) {
		return "", false
	}
	members := enumMembersFromColumnType(col.ColumnType)
	target := collationComparisonKey(strings.TrimRight(string(raw), " "), col.Collation, col.Charset)
	for _, m := range members {
		if collationComparisonKey(strings.TrimRight(m, " "), col.Collation, col.Charset) == target {
			return nodeToSQL(&nodes.StringLit{Value: m}), true
		}
	}
	if target == "" && !strict {
		return "''", true
	}
	return "", false
}

func coerceBinarySetDefault(col *Column, isHex bool, val string) (string, bool) {
	raw, ok := decodeLiteralRawBytes(val, isHex)
	if !ok || !utf8.Valid(raw) {
		return "", false
	}
	s := string(raw)
	if s == "" {
		return "''", true
	}
	members := enumMembersFromColumnType(col.ColumnType)
	chosen := make([]bool, len(members))
	for _, part := range strings.Split(s, ",") {
		pk := collationComparisonKey(strings.TrimRight(part, " "), col.Collation, col.Charset)
		matched := false
		for i, m := range members {
			if collationComparisonKey(strings.TrimRight(m, " "), col.Collation, col.Charset) == pk {
				chosen[i] = true
				matched = true
				break
			}
		}
		if !matched {
			return "", false
		}
	}
	var out []string
	for i, m := range members {
		if chosen[i] {
			out = append(out, m)
		}
	}
	return nodeToSQL(&nodes.StringLit{Value: strings.Join(out, ",")}), true
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

func temporalColumnFSP(colType string) int {
	pt := parseColumnType(colType)
	if pt.hasLen && pt.length >= 0 && pt.length <= 6 {
		return pt.length
	}
	return 0
}

func pow10Int(n int) int {
	p := 1
	for i := 0; i < n; i++ {
		p *= 10
	}
	return p
}

func convertSessionLocalToUTC(year, month, day, hour, min, sec, nsec int, sessionTZ string) (time.Time, bool) {
	tz := strings.TrimSpace(sessionTZ)
	if tz == "" || strings.EqualFold(tz, "SYSTEM") || tz == "+00:00" || tz == "-00:00" || strings.EqualFold(tz, "UTC") {
		return time.Date(year, time.Month(month), day, hour, min, sec, nsec, time.UTC), true
	}
	if len(tz) == 6 && (tz[0] == '+' || tz[0] == '-') && tz[3] == ':' {
		hh := int(tz[1]-'0')*10 + int(tz[2]-'0')
		mm := int(tz[4]-'0')*10 + int(tz[5]-'0')
		off := hh*60 + mm
		if tz[0] == '-' {
			off = -off
		}
		return time.Date(year, time.Month(month), day, hour, min, sec, nsec, time.UTC).Add(-time.Duration(off) * time.Minute), true
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return time.Date(year, time.Month(month), day, hour, min, sec, nsec, time.UTC), true
	}
	localT := time.Date(year, time.Month(month), day, hour, min, sec, nsec, loc)
	// Reject spring-forward DST gap times that do not exist in loc.
	if localT.Year() != year || localT.Month() != time.Month(month) || localT.Day() != day ||
		localT.Hour() != hour || localT.Minute() != min || localT.Second() != sec {
		return time.Time{}, false
	}
	return localT.UTC(), true
}

type parsedMySQLDateTime struct {
	year, month, day int
	hour, min, sec   int
	micros           int
	allZero          bool
	hasTZ            bool
	tzOffsetMin      int
}

func roundOrTruncateFractionRatToFSP(frac *big.Rat, fsp int, truncate bool) (micros int, carrySec bool) {
	if frac == nil || frac.Sign() <= 0 {
		return 0, false
	}
	scale := pow10Int(fsp)
	scaled := new(big.Rat).Mul(frac, big.NewRat(int64(scale), 1))
	if truncate {
		ticks := int(new(big.Int).Quo(scaled.Num(), scaled.Denom()).Int64())
		return ticks * pow10Int(6-fsp), false
	}
	ticks := int(roundRatHalfAwayFromZero(scaled).Int64())
	if ticks >= scale {
		return 0, true
	}
	return ticks * pow10Int(6-fsp), false
}

func roundOrTruncateFractionDigitsToFSP(fracDigits string, fsp int, truncate bool) (micros int, carrySec bool) {
	trimmed := strings.TrimRight(fracDigits, "0")
	if trimmed == "" {
		return 0, false
	}
	r, ok := new(big.Rat).SetString("0." + trimmed)
	if !ok {
		return 0, false
	}
	return roundOrTruncateFractionRatToFSP(r, fsp, truncate)
}

func addOneSecondToDateTime(year, month, day, hour, min, sec int) (int, int, int, int, int, int, bool) {
	if month >= 1 && month <= 12 && day >= 1 && day <= daysInMonth(year, month) {
		t := time.Date(year, time.Month(month), day, hour, min, sec, 0, time.UTC).Add(time.Second)
		if t.Year() > 9999 {
			return 0, 0, 0, 0, 0, 0, false
		}
		return t.Year(), int(t.Month()), t.Day(), t.Hour(), t.Minute(), t.Second(), true
	}
	sec++
	if sec < 60 {
		return year, month, day, hour, min, sec, true
	}
	sec = 0
	min++
	if min < 60 {
		return year, month, day, hour, min, sec, true
	}
	min = 0
	hour++
	if hour < 24 {
		return year, month, day, hour, min, sec, true
	}
	hour = 0
	day++
	if day <= 31 {
		return year, month, day, hour, min, sec, true
	}
	return 0, 0, 0, 0, 0, 0, false
}

// parseMySQLNumberDateTime implements MySQL 8.0's number_to_datetime /
// my_decimal_to_datetime rules for unquoted numeric temporal literals:
//   - 0 -> 0000-00-00 00:00:00
//   - 101..999999 -> left-padded to 6 digits YYMMDD (2-digit year)
//   - 1000000..99999999 -> left-padded to 8 digits YYYYMMDD (4-digit year)
//   - 100000000..999999999999 -> left-padded to 12 digits YYMMDDHHMMSS (2-digit year)
//   - 1000000000000..99999999999999 -> left-padded to 14 digits YYYYMMDDHHMMSS (4-digit year)
//
// Date-only values (< 100000000) reject non-zero fractional seconds.
func parseMySQLNumberDateTime(dt string, fsp int, val string, truncateFrac bool) (parsedMySQLDateTime, bool) {
	s := strings.TrimSpace(val)
	if !isStrictNumberSyntax(s) {
		return parsedMySQLDateTime{}, false
	}
	r, ok := new(big.Rat).SetString(s)
	if !ok || r.Sign() < 0 {
		return parsedMySQLDateTime{}, false
	}
	intPart := new(big.Int).Quo(r.Num(), r.Denom())
	rem := new(big.Int).Rem(r.Num(), r.Denom())
	if rem.Sign() != 0 && intPart.Cmp(big.NewInt(100000000)) < 0 {
		return parsedMySQLDateTime{}, false
	}
	if !intPart.IsUint64() {
		return parsedMySQLDateTime{}, false
	}
	nr := intPart.Uint64()
	if nr == 0 {
		return parsedMySQLDateTime{allZero: true}, true
	}
	if nr < 101 || nr >= 100000000000000 {
		return parsedMySQLDateTime{}, false
	}
	var out parsedMySQLDateTime
	switch {
	case nr <= 999999:
		d := fmt.Sprintf("%06d", nr)
		y, _ := strconv.Atoi(d[0:2])
		out.month, _ = strconv.Atoi(d[2:4])
		out.day, _ = strconv.Atoi(d[4:6])
		out.year = expandTwoDigitYear(y, out.month, out.day)
	case nr <= 99999999:
		d := fmt.Sprintf("%08d", nr)
		out.year, _ = strconv.Atoi(d[0:4])
		out.month, _ = strconv.Atoi(d[4:6])
		out.day, _ = strconv.Atoi(d[6:8])
	case nr <= 999999999999:
		d := fmt.Sprintf("%012d", nr)
		y, _ := strconv.Atoi(d[0:2])
		out.month, _ = strconv.Atoi(d[2:4])
		out.day, _ = strconv.Atoi(d[4:6])
		out.year = expandTwoDigitYear(y, out.month, out.day)
		out.hour, _ = strconv.Atoi(d[6:8])
		out.min, _ = strconv.Atoi(d[8:10])
		out.sec, _ = strconv.Atoi(d[10:12])
	default:
		d := fmt.Sprintf("%014d", nr)
		out.year, _ = strconv.Atoi(d[0:4])
		out.month, _ = strconv.Atoi(d[4:6])
		out.day, _ = strconv.Atoi(d[6:8])
		out.hour, _ = strconv.Atoi(d[8:10])
		out.min, _ = strconv.Atoi(d[10:12])
		out.sec, _ = strconv.Atoi(d[12:14])
	}
	if out.month > 12 || out.day > 31 || out.hour > 23 || out.min > 59 || out.sec > 59 {
		return parsedMySQLDateTime{}, false
	}
	if dt != "date" && rem.Sign() > 0 {
		frac := new(big.Rat).SetFrac(rem, r.Denom())
		micros, carry := roundOrTruncateFractionRatToFSP(frac, fsp, truncateFrac)
		out.micros = micros
		if carry {
			var carryOK bool
			out.year, out.month, out.day, out.hour, out.min, out.sec, carryOK =
				addOneSecondToDateTime(out.year, out.month, out.day, out.hour, out.min, out.sec)
			if !carryOK {
				return parsedMySQLDateTime{}, false
			}
		}
	}
	return out, true
}

func parseMySQLStringDateTime(dt string, fsp int, s string, truncateFrac bool) (parsedMySQLDateTime, bool) {
	s = strings.TrimSpace(s)
	if s == "" || !isASCIIDigit(s[0]) {
		return parsedMySQLDateTime{}, false
	}
	var out parsedMySQLDateTime
	base := s
	// Check for an attached timezone suffix `[+-]HH:MM` immediately following a
	// delimited 6-field date-and-time value (no whitespace before the sign).
	// MySQL 8.0.46 accepts `-14:00` through `+14:00` and rejects `-00:00` and
	// compact digit-only prefixes (`20200101123456+00:00`).
	if len(base) >= 7 {
		tz := base[len(base)-6:]
		if (tz[0] == '+' || tz[0] == '-') && isASCIIDigit(tz[1]) && isASCIIDigit(tz[2]) &&
			tz[3] == ':' && isASCIIDigit(tz[4]) && isASCIIDigit(tz[5]) && isASCIIDigit(base[len(base)-7]) {
			prefix := base[:len(base)-6]
			prefixNoFrac, _ := splitDateTimeTrailingFraction(prefix)
			allPrefixDigits := len(prefixNoFrac) > 0
			for i := 0; i < len(prefixNoFrac); i++ {
				if !isASCIIDigit(prefixNoFrac[i]) {
					allPrefixDigits = false
					break
				}
			}
			if !allPrefixDigits && countDigitRuns(prefixNoFrac) == 6 {
				tzH := int(tz[1]-'0')*10 + int(tz[2]-'0')
				tzM := int(tz[4]-'0')*10 + int(tz[5]-'0')
				if tz[0] == '-' && tzH == 0 && tzM == 0 {
					return parsedMySQLDateTime{}, false
				}
				if tzM > 59 || tzH*60+tzM > 14*60 {
					return parsedMySQLDateTime{}, false
				}
				out.hasTZ = true
				out.tzOffsetMin = tzH*60 + tzM
				if tz[0] == '-' {
					out.tzOffsetMin = -out.tzOffsetMin
				}
				base = prefix
			}
		}
	}

	var parts []string
	var seps []string // seps[k] precedes parts[k]
	i := 0
	for i < len(base) {
		sepStart := i
		for i < len(base) && !isASCIIDigit(base[i]) {
			i++
		}
		sep := base[sepStart:i]
		if i >= len(base) {
			// Reject trailing non-space characters such as 'Z' or trailing '.'.
			if strings.TrimSpace(sep) != "" {
				return parsedMySQLDateTime{}, false
			}
			break
		}
		digStart := i
		for i < len(base) && isASCIIDigit(base[i]) {
			i++
		}
		seps = append(seps, sep)
		parts = append(parts, base[digStart:i])
	}
	if len(parts) == 0 {
		return parsedMySQLDateTime{}, false
	}

	// Validate separators: reject ASCII letters except a single 'T'/'t' between
	// the date part and time part.
	p0 := parts[0]
	n0 := len(p0)
	tSepIdx := 3
	if n0 == 5 || n0 == 6 || n0 == 8 {
		tSepIdx = 1
	}
	for idx, sep := range seps {
		for k := 0; k < len(sep); k++ {
			ch := sep[k]
			if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') {
				if idx == tSepIdx && (sep == "T" || sep == "t") {
					continue
				}
				return parsedMySQLDateTime{}, false
			}
		}
	}

	var fracDigits string
	if n0 >= 5 {
		if n0 == 13 || n0 > 14 {
			return parsedMySQLDateTime{}, false
		}
		filled := 0
		switch n0 {
		case 8:
			out.year, _ = strconv.Atoi(p0[0:4])
			out.month, _ = strconv.Atoi(p0[4:6])
			out.day, _ = strconv.Atoi(p0[6:8])
			filled = 3
		case 14:
			out.year, _ = strconv.Atoi(p0[0:4])
			out.month, _ = strconv.Atoi(p0[4:6])
			out.day, _ = strconv.Atoi(p0[6:8])
			out.hour, _ = strconv.Atoi(p0[8:10])
			out.min, _ = strconv.Atoi(p0[10:12])
			out.sec, _ = strconv.Atoi(p0[12:14])
			filled = 6
		default:
			yy, _ := strconv.Atoi(p0[0:2])
			out.month, _ = strconv.Atoi(p0[2:4])
			dEnd := 6
			if n0 < 6 {
				dEnd = n0
			}
			out.day, _ = strconv.Atoi(p0[4:dEnd])
			out.year = expandTwoDigitYear(yy, out.month, out.day)
			filled = 3
			if n0 >= 7 {
				hEnd := 8
				if n0 < 8 {
					hEnd = n0
				}
				out.hour, _ = strconv.Atoi(p0[6:hEnd])
				filled = 4
			}
			if n0 >= 9 {
				mEnd := 10
				if n0 < 10 {
					mEnd = n0
				}
				out.min, _ = strconv.Atoi(p0[8:mEnd])
				filled = 5
			}
			if n0 >= 11 {
				sEnd := 12
				if n0 < 12 {
					sEnd = n0
				}
				out.sec, _ = strconv.Atoi(p0[10:sEnd])
				filled = 6
			}
		}
		for idx := 1; idx < len(parts); idx++ {
			if filled < 6 {
				v, err := strconv.Atoi(parts[idx])
				if err != nil {
					return parsedMySQLDateTime{}, false
				}
				switch filled {
				case 3:
					out.hour = v
				case 4:
					out.min = v
				case 5:
					out.sec = v
				}
				filled++
			} else if filled == 6 && seps[idx] == "." {
				fracDigits = parts[idx]
				filled = 7
			} else {
				return parsedMySQLDateTime{}, false
			}
		}
	} else {
		if len(parts) < 3 || len(parts) > 7 {
			return parsedMySQLDateTime{}, false
		}
		y, err1 := strconv.Atoi(parts[0])
		m, err2 := strconv.Atoi(parts[1])
		d, err3 := strconv.Atoi(parts[2])
		if err1 != nil || err2 != nil || err3 != nil {
			return parsedMySQLDateTime{}, false
		}
		if len(parts[0]) <= 2 {
			y = expandTwoDigitYear(y, m, d)
		}
		out.year, out.month, out.day = y, m, d
		if len(parts) >= 4 {
			v, err := strconv.Atoi(parts[3])
			if err != nil {
				return parsedMySQLDateTime{}, false
			}
			out.hour = v
		}
		if len(parts) >= 5 {
			v, err := strconv.Atoi(parts[4])
			if err != nil {
				return parsedMySQLDateTime{}, false
			}
			out.min = v
		}
		if len(parts) >= 6 {
			v, err := strconv.Atoi(parts[5])
			if err != nil {
				return parsedMySQLDateTime{}, false
			}
			out.sec = v
		}
		if len(parts) == 7 {
			if seps[6] != "." {
				return parsedMySQLDateTime{}, false
			}
			fracDigits = parts[6]
		}
	}

	out.allZero = out.year == 0 && out.month == 0 && out.day == 0 &&
		out.hour == 0 && out.min == 0 && out.sec == 0 && strings.TrimLeft(fracDigits, "0") == ""
	if out.allZero {
		return out, true
	}
	if out.year > 9999 || out.month > 12 || out.day > 31 || out.hour > 23 || out.min > 59 || out.sec > 59 {
		return parsedMySQLDateTime{}, false
	}
	if dt != "date" && fracDigits != "" {
		// MySQL truncates string microseconds beyond 6 digits before rounding or
		// truncating to fsp.
		if len(fracDigits) > 6 {
			fracDigits = fracDigits[:6]
		}
		micros, carry := roundOrTruncateFractionDigitsToFSP(fracDigits, fsp, truncateFrac)
		out.micros = micros
		if carry {
			var carryOK bool
			out.year, out.month, out.day, out.hour, out.min, out.sec, carryOK =
				addOneSecondToDateTime(out.year, out.month, out.day, out.hour, out.min, out.sec)
			if !carryOK {
				return parsedMySQLDateTime{}, false
			}
		}
	}
	return out, true
}

func splitDateTimeTrailingFraction(s string) (string, string) {
	dot := strings.LastIndexByte(s, '.')
	if dot < 0 {
		return s, ""
	}
	frac := s[dot+1:]
	if len(frac) == 0 {
		return s, ""
	}
	for i := 0; i < len(frac); i++ {
		if !isASCIIDigit(frac[i]) {
			return s, ""
		}
	}
	if strings.ContainsAny(s[:dot], " Tt:") || len(s[:dot]) >= 11 {
		return s[:dot], frac
	}
	return s, ""
}

func validDateOrDatetimeLiteralDefault(dt string, fsp int, kind literalDefaultKind, val string, unquotedNum, noZeroDate, noZeroInDate, allowInvalidDates, truncateFrac bool, sessionTZ string) bool {
	var p parsedMySQLDateTime
	var ok bool
	switch {
	case kind == litNumber && unquotedNum:
		p, ok = parseMySQLNumberDateTime(dt, fsp, val, truncateFrac)
	case kind == litString:
		p, ok = parseMySQLStringDateTime(dt, fsp, val, truncateFrac)
	case kind == litHex || kind == litBit:
		raw, decOK := decodeLiteralRawBytes(val, kind == litHex)
		if !decOK {
			return false
		}
		p, ok = parseMySQLStringDateTime(dt, fsp, string(raw), truncateFrac)
	default:
		return false
	}
	if !ok {
		return false
	}
	if p.allZero {
		if p.hasTZ {
			return false
		}
		return !noZeroDate
	}
	if p.month == 0 || p.day == 0 {
		if dt == "timestamp" || p.hasTZ || noZeroInDate {
			return false
		}
		return true
	}
	if (dt == "timestamp" || p.hasTZ || !allowInvalidDates) && p.day > daysInMonth(p.year, p.month) {
		return false
	}
	minUTC := time.Date(1970, time.January, 1, 0, 0, 1, 0, time.UTC)
	maxUTC := time.Date(2038, time.January, 19, 3, 14, 7, 999999000, time.UTC)
	nsec := p.micros * 1000
	if p.hasTZ {
		utc := time.Date(p.year, time.Month(p.month), p.day, p.hour, p.min, p.sec, nsec, time.UTC).Add(-time.Duration(p.tzOffsetMin) * time.Minute)
		if utc.Year() < 0 || utc.Year() > 9999 {
			return false
		}
		if dt == "timestamp" && (utc.Before(minUTC) || utc.After(maxUTC)) {
			return false
		}
		localT := convertUTCToSessionLocal(utc, sessionTZ)
		if localT.Year() < 0 || localT.Year() > 9999 {
			return false
		}
		return true
	}
	if dt == "timestamp" {
		utc, validLocal := convertSessionLocalToUTC(p.year, p.month, p.day, p.hour, p.min, p.sec, nsec, sessionTZ)
		if !validLocal || utc.Before(minUTC) || utc.After(maxUTC) {
			return false
		}
	}
	return true
}

func convertUTCToSessionLocal(utc time.Time, sessionTZ string) time.Time {
	tz := strings.TrimSpace(sessionTZ)
	if tz == "" || strings.EqualFold(tz, "SYSTEM") || tz == "+00:00" || tz == "-00:00" || strings.EqualFold(tz, "UTC") {
		return utc.UTC()
	}
	if len(tz) == 6 && (tz[0] == '+' || tz[0] == '-') && tz[3] == ':' {
		hh := int(tz[1]-'0')*10 + int(tz[2]-'0')
		mm := int(tz[4]-'0')*10 + int(tz[5]-'0')
		off := hh*60 + mm
		if tz[0] == '-' {
			off = -off
		}
		return utc.UTC().Add(time.Duration(off) * time.Minute)
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return utc.UTC()
	}
	return utc.In(loc)
}

func resolveTimezoneDateTimeToSession(p parsedMySQLDateTime, sessionTZ string) (parsedMySQLDateTime, bool) {
	if !p.hasTZ {
		return p, true
	}
	nsec := p.micros * 1000
	utc := time.Date(p.year, time.Month(p.month), p.day, p.hour, p.min, p.sec, nsec, time.UTC).Add(-time.Duration(p.tzOffsetMin) * time.Minute)
	if utc.Year() < 0 || utc.Year() > 9999 {
		return parsedMySQLDateTime{}, false
	}
	localT := convertUTCToSessionLocal(utc, sessionTZ)
	if localT.Year() < 0 || localT.Year() > 9999 {
		return parsedMySQLDateTime{}, false
	}
	return parsedMySQLDateTime{
		year:   localT.Year(),
		month:  int(localT.Month()),
		day:    localT.Day(),
		hour:   localT.Hour(),
		min:    localT.Minute(),
		sec:    localT.Second(),
		micros: localT.Nanosecond() / 1000,
	}, true
}

func formatParsedDateOrDatetime(dt string, fsp int, p parsedMySQLDateTime) string {
	if dt == "date" {
		if p.allZero {
			return "'0000-00-00'"
		}
		return fmt.Sprintf("'%04d-%02d-%02d'", p.year, p.month, p.day)
	}
	fracSuffix := ""
	if fsp > 0 {
		fracSuffix = fmt.Sprintf(".%0*d", fsp, p.micros/pow10Int(6-fsp))
	}
	if p.allZero {
		return fmt.Sprintf("'0000-00-00 00:00:00%s'", fracSuffix)
	}
	return fmt.Sprintf("'%04d-%02d-%02d %02d:%02d:%02d%s'",
		p.year, p.month, p.day, p.hour, p.min, p.sec, fracSuffix)
}

func coerceDateOrDatetimeDefault(dt string, fsp int, kind literalDefaultKind, val string, unquotedNum, truncateFrac bool, sessionTZ string) (string, bool) {
	switch {
	case kind == litNumber && unquotedNum:
		p, ok := parseMySQLNumberDateTime(dt, fsp, val, truncateFrac)
		if !ok {
			return "", false
		}
		return formatParsedDateOrDatetime(dt, fsp, p), true
	case kind == litHex || kind == litBit:
		raw, decOK := decodeLiteralRawBytes(val, kind == litHex)
		if !decOK {
			return "", false
		}
		p, ok := parseMySQLStringDateTime(dt, fsp, string(raw), truncateFrac)
		if !ok {
			return "", false
		}
		p, ok = resolveTimezoneDateTimeToSession(p, sessionTZ)
		if !ok {
			return "", false
		}
		return formatParsedDateOrDatetime(dt, fsp, p), true
	case kind == litString:
		s := strings.TrimSpace(val)
		p, ok := parseMySQLStringDateTime(dt, fsp, s, truncateFrac)
		if !ok {
			return "", false
		}
		if p.hasTZ {
			p, ok = resolveTimezoneDateTimeToSession(p, sessionTZ)
			if !ok {
				return "", false
			}
			return formatParsedDateOrDatetime(dt, fsp, p), true
		}
		if dt == "date" {
			return "", false
		}
		_, frac := splitDateTimeTrailingFraction(s)
		if len(frac) <= fsp {
			return "", false
		}
		return formatParsedDateOrDatetime(dt, fsp, p), true
	default:
		return "", false
	}
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

func expandTwoDigitYear(y, m, d int) int {
	if y == 0 && m == 0 && d == 0 {
		return 0
	}
	if y >= 70 {
		return 1900 + y
	}
	return 2000 + y
}

func parseTimeValue(fsp int, s string, unquotedNum, truncateFrac bool) (neg bool, totalHours, mm, ss, micros int, ok bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return false, 0, 0, 0, 0, false
	}
	if unquotedNum {
		if !isStrictNumberSyntax(s) {
			return false, 0, 0, 0, 0, false
		}
		r, ratOK := new(big.Rat).SetString(s)
		if !ratOK {
			return false, 0, 0, 0, 0, false
		}
		neg = r.Sign() < 0
		absR := new(big.Rat).Abs(r)
		intPart := new(big.Int).Quo(absR.Num(), absR.Denom())
		rem := new(big.Int).Rem(absR.Num(), absR.Denom())
		if !intPart.IsUint64() {
			return false, 0, 0, 0, 0, false
		}
		nr := intPart.Uint64()
		if !neg && nr >= 100000000 {
			p, dtOK := parseMySQLNumberDateTime("datetime", fsp, s, truncateFrac)
			if !dtOK || (!p.allZero && p.day > daysInMonth(p.year, p.month)) {
				return false, 0, 0, 0, 0, false
			}
			return false, p.hour, p.min, p.sec, p.micros, true
		}
		ss = int(nr % 100)
		mm = int((nr / 100) % 100)
		totalHours = int(nr / 10000)
		if mm >= 60 || ss >= 60 || totalHours > 838 {
			return false, 0, 0, 0, 0, false
		}
		if rem.Sign() > 0 {
			// MySQL's number_to_time enforces TIME_MAX_VALUE (838:59:59.000000)
			// on the input microseconds before TIME_TRUNCATE_FRACTIONAL applies.
			frac := new(big.Rat).SetFrac(rem, absR.Denom())
			rawMicros, rawCarry := roundOrTruncateFractionRatToFSP(frac, 6, true)
			if totalHours == 838 && mm == 59 && ss == 59 && (rawMicros > 0 || rawCarry) {
				return false, 0, 0, 0, 0, false
			}
			var carry bool
			micros, carry = roundOrTruncateFractionRatToFSP(frac, fsp, truncateFrac)
			if carry {
				ss++
				if ss >= 60 {
					ss = 0
					mm++
					if mm >= 60 {
						mm = 0
						totalHours++
					}
				}
			}
		}
		if totalHours > 838 || (totalHours == 838 && mm == 59 && ss == 59 && micros > 0) {
			return false, 0, 0, 0, 0, false
		}
		if totalHours == 0 && mm == 0 && ss == 0 && micros == 0 {
			neg = false
		}
		return neg, totalHours, mm, ss, micros, true
	}

	orig := s
	if strings.HasPrefix(s, "-") {
		neg = true
		s = strings.TrimSpace(s[1:])
	} else if strings.HasPrefix(s, "+") {
		s = strings.TrimSpace(s[1:])
	}
	if s == "" {
		return false, 0, 0, 0, 0, false
	}
	var fracDigits string
	if dot := strings.IndexByte(s, '.'); dot >= 0 {
		fracDigits = s[dot+1:]
		if fracDigits == "" || len(fracDigits) > 6 {
			return false, 0, 0, 0, 0, false
		}
		for i := 0; i < len(fracDigits); i++ {
			if !isASCIIDigit(fracDigits[i]) {
				return false, 0, 0, 0, 0, false
			}
		}
		s = s[:dot]
	}
	hasRawNonZeroFrac := strings.TrimRight(fracDigits, "0") != ""
	plainTime := len(s) > 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if (c < '0' || c > '9') && c != ':' && c != ' ' {
			plainTime = false
			break
		}
	}
	if plainTime {
		days := 0
		hasSpace := false
		if sp := strings.IndexByte(s, ' '); sp >= 0 {
			hasSpace = true
			d, err := strconv.Atoi(strings.TrimSpace(s[:sp]))
			if err == nil && d <= 34 {
				days = d
				s = strings.TrimSpace(s[sp+1:])
			} else {
				plainTime = false
			}
		}
		if plainTime && strings.Contains(s, ":") {
			fields := strings.Split(s, ":")
			if len(fields) >= 2 && len(fields) <= 3 {
				hh, err1 := strconv.Atoi(fields[0])
				m, err2 := strconv.Atoi(fields[1])
				sec := 0
				var err3 error
				if len(fields) == 3 {
					sec, err3 = strconv.Atoi(fields[2])
				}
				if err1 == nil && err2 == nil && err3 == nil {
					th := days*24 + hh
					if m < 60 && sec < 60 && th <= 838 {
						if th == 838 && m == 59 && sec == 59 && hasRawNonZeroFrac {
							return false, 0, 0, 0, 0, false
						}
						mc, carry := roundOrTruncateFractionDigitsToFSP(fracDigits, fsp, truncateFrac)
						if carry {
							sec++
							if sec >= 60 {
								sec = 0
								m++
								if m >= 60 {
									m = 0
									th++
								}
							}
						}
						if th > 838 || (th == 838 && m == 59 && sec == 59 && mc > 0) {
							return false, 0, 0, 0, 0, false
						}
						if th == 0 && m == 0 && sec == 0 && mc == 0 {
							neg = false
						}
						return neg, th, m, sec, mc, true
					}
				}
			}
		} else if plainTime && !hasSpace {
			if n, err := strconv.Atoi(s); err == nil {
				sec := n % 100
				m := (n / 100) % 100
				th := n / 10000
				if m < 60 && sec < 60 && th <= 838 {
					if th == 838 && m == 59 && sec == 59 && hasRawNonZeroFrac {
						return false, 0, 0, 0, 0, false
					}
					mc, carry := roundOrTruncateFractionDigitsToFSP(fracDigits, fsp, truncateFrac)
					if carry {
						sec++
						if sec >= 60 {
							sec = 0
							m++
							if m >= 60 {
								m = 0
								th++
							}
						}
					}
					if th > 838 || (th == 838 && m == 59 && sec == 59 && mc > 0) {
						return false, 0, 0, 0, 0, false
					}
					if th == 0 && m == 0 && sec == 0 && mc == 0 {
						neg = false
					}
					return neg, th, m, sec, mc, true
				}
			}
		}
	}
	if !neg {
		if p, dtOK := parseMySQLStringDateTime("datetime", fsp, orig, truncateFrac); dtOK && (!p.allZero && p.month >= 1 && p.day >= 1 && p.day <= daysInMonth(p.year, p.month)) {
			return false, p.hour, p.min, p.sec, p.micros, true
		}
	}
	return false, 0, 0, 0, 0, false
}

func validTimeLiteralDefault(fsp int, kind literalDefaultKind, val string, unquotedNum, truncateFrac bool) bool {
	switch kind {
	case litNumber:
		_, _, _, _, _, ok := parseTimeValue(fsp, val, unquotedNum, truncateFrac)
		return ok
	case litString:
		_, _, _, _, _, ok := parseTimeValue(fsp, val, false, truncateFrac)
		return ok
	case litHex, litBit:
		raw, decOK := decodeLiteralRawBytes(val, kind == litHex)
		if !decOK {
			return false
		}
		_, _, _, _, _, ok := parseTimeValue(fsp, string(raw), false, truncateFrac)
		return ok
	default:
		return false
	}
}

func formatParsedTime(fsp int, neg bool, th, mm, ss, micros int) string {
	sign := ""
	if neg {
		sign = "-"
	}
	fracSuffix := ""
	if fsp > 0 {
		fracSuffix = fmt.Sprintf(".%0*d", fsp, micros/pow10Int(6-fsp))
	}
	return fmt.Sprintf("'%s%02d:%02d:%02d%s'", sign, th, mm, ss, fracSuffix)
}

func coerceTimeDefault(fsp int, kind literalDefaultKind, val string, unquotedNum, truncateFrac bool) (string, bool) {
	switch {
	case kind == litNumber && unquotedNum:
		neg, th, mm, ss, micros, ok := parseTimeValue(fsp, val, true, truncateFrac)
		if !ok {
			return "", false
		}
		return formatParsedTime(fsp, neg, th, mm, ss, micros), true
	case kind == litHex || kind == litBit:
		raw, decOK := decodeLiteralRawBytes(val, kind == litHex)
		if !decOK {
			return "", false
		}
		neg, th, mm, ss, micros, ok := parseTimeValue(fsp, string(raw), false, truncateFrac)
		if !ok {
			return "", false
		}
		return formatParsedTime(fsp, neg, th, mm, ss, micros), true
	case kind == litString:
		s := strings.TrimSpace(val)
		dot := strings.LastIndexByte(s, '.')
		if dot < 0 || len(s[dot+1:]) <= fsp {
			return "", false
		}
		neg, th, mm, ss, micros, ok := parseTimeValue(fsp, s, false, truncateFrac)
		if !ok {
			return "", false
		}
		return formatParsedTime(fsp, neg, th, mm, ss, micros), true
	default:
		return "", false
	}
}
