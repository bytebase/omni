package catalog

import (
	"math/big"
	"regexp"
	"strings"
	"unicode/utf8"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// OceanBase reads a column's definition before it looks at the column's name, so
// each rule below answers before ADD COLUMN's 1060 for a name the table has.
// TestOracle_OceanBaseColumnDefinition holds them to the engine.

// checkColumnDefinition refuses what OceanBase refuses in a column definition by
// itself: a literal default the type cannot hold (1067), a column in an expression
// default (1235), an SRID on a type that is not spatial (1221), and a column CHECK
// that is a bare column, not a condition (3812).
func checkColumnDefinition(colDef *nodes.ColumnDef) error {
	if colDef == nil || colDef.TypeName == nil {
		return nil
	}
	defaults := []nodes.ExprNode{colDef.DefaultValue}
	references := false
	for _, cc := range colDef.Constraints {
		switch cc.Type {
		case nodes.ColConstrDefault:
			defaults = append(defaults, cc.Expr)
		case nodes.ColConstrReferences:
			references = true
		case nodes.ColConstrCheck:
			// A CHECK after REFERENCES goes with it and is ignored: ADD a INT
			// REFERENCES t (a) CHECK (a) is taken, ADD a INT CHECK (a) is 3812.
			if !references && isBareColumn(cc.Expr) {
				return errNonBooleanCheck(cc.Name)
			}
		}
	}
	for _, d := range defaults {
		if err := checkDefault(colDef.Name, colDef.TypeName, d); err != nil {
			return err
		}
	}
	if colDef.TypeName.HasSRID && !isSpatialTypeName(normalizedColumnDataType(colDef.TypeName)) {
		return &Error{Code: 1221, SQLState: "HY000", Message: "Incorrect usage of srid."}
	}
	return nil
}

// checkDefault refuses a default a column of type dt cannot have. dt is nil for a
// column whose type is known only by its catalog form, which checks the expression
// rules alone.
func checkDefault(colName string, dt *nodes.DataType, expr nodes.ExprNode) error {
	if expr == nil {
		return nil
	}
	if paren, ok := expr.(*nodes.ParenExpr); ok {
		// An expression default may not read a column, existing or not.
		if hasColumnRef(paren) {
			return errNotSupportedDefault()
		}
		return nil
	}
	if dt == nil {
		return nil
	}
	if !literalDefaultFits(dt, expr) {
		return errInvalidDefault(colName)
	}
	return nil
}

func errNotSupportedDefault() error {
	return &Error{Code: 1235, SQLState: "0A000", Message: "Column reference in default value expression not supported"}
}

func errNonBooleanCheck(name string) error {
	return &Error{Code: 3812, SQLState: "HY000", Message: "An expression of non-boolean type specified to a check constraint '" + name + "'"}
}

func hasColumnRef(expr nodes.Node) bool {
	found := false
	nodes.Inspect(expr, func(n nodes.Node) bool {
		if _, ok := n.(*nodes.ColumnRef); ok {
			found = true
		}
		return !found
	})
	return found
}

func isBareColumn(expr nodes.ExprNode) bool {
	for {
		switch e := expr.(type) {
		case *nodes.ParenExpr:
			expr = e.Expr
		case *nodes.ColumnRef:
			return true
		default:
			return false
		}
	}
}

func isSpatialTypeName(name string) bool {
	switch strings.ToLower(name) {
	case "geometry", "point", "linestring", "polygon", "multipoint", "multilinestring",
		"multipolygon", "geometrycollection", "geomcollection":
		return true
	}
	return false
}

var (
	numericText  = regexp.MustCompile(`^[+-]?(\d+\.?\d*|\.\d+)([eE][+-]?\d+)?$`)
	dateText     = regexp.MustCompile(`^(\d{2,4})[^\d](\d{1,2})[^\d](\d{1,2})(?:[ T](\d{1,2})[^\d](\d{1,2})(?:[^\d](\d{1,2})(?:\.\d+)?)?)?$`)
	dateDigits   = regexp.MustCompile(`^(\d{4})(\d{2})(\d{2})(?:(\d{2})(\d{2})(\d{2}))?$`)
	timeText     = regexp.MustCompile(`^-?(?:\d+ )?\d{1,3}(?::(\d{1,2})(?::(\d{1,2})(?:\.\d+)?)?)?$`)
	timeDigits   = regexp.MustCompile(`^-?\d+(?:\.\d+)?$`)
	temporalFunc = map[string]bool{"current_timestamp": true, "now": true, "localtime": true, "localtimestamp": true}
)

// literalDefaultFits reports whether a column of type dt can hold the literal
// default expr in strict mode: a number in range for a numeric type (a string
// default is read as a number, '1.5' and ' 1' alike, and ” or '0x10' is none), a
// string no longer than a CHAR or BINARY's length, a member of an ENUM or SET, a
// date or time a temporal type can read, and CURRENT_TIMESTAMP only on DATETIME and
// TIMESTAMP. Literals of other kinds (TRUE, b'1', 0x10) and types not listed (YEAR,
// which takes 's') are taken, as OceanBase takes them.
func literalDefaultFits(dt *nodes.DataType, expr nodes.ExprNode) bool {
	typeName := strings.ToLower(normalizedColumnDataType(dt))
	if fn, ok := expr.(*nodes.FuncCallExpr); ok && temporalFunc[strings.ToLower(fn.Name)] {
		return typeName == "datetime" || typeName == "timestamp"
	}
	text, isString, ok := literalText(expr)
	if !ok {
		return true
	}
	switch typeName {
	case "tinyint", "smallint", "mediumint", "int", "integer", "bigint":
		n, ok := literalNumber(text)
		if !ok {
			return false
		}
		return integerInRange(n, typeName, dt.Unsigned || dt.Zerofill)
	case "decimal", "numeric", "dec", "fixed":
		n, ok := literalNumber(text)
		if !ok {
			return false
		}
		precision, scale := dt.Length, dt.Scale
		if precision == 0 {
			precision = 10
		}
		// The integer part may have precision - scale digits once rounded to scale.
		limit := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(precision-scale)), nil))
		half := new(big.Float).Quo(big.NewFloat(0.5), new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(scale)), nil)))
		abs := new(big.Float).Abs(n)
		return abs.Add(abs, half).Cmp(limit) < 0
	case "float", "double", "real":
		_, ok := literalNumber(text)
		return ok
	case "bit":
		n, ok := literalNumber(text)
		if !ok || isString {
			return true
		}
		bits := dt.Length
		if bits == 0 {
			bits = 1
		}
		limit := new(big.Float).SetInt(new(big.Int).Lsh(big.NewInt(1), uint(bits)))
		return n.Sign() >= 0 && n.Cmp(limit) < 0
	case "char", "binary":
		if !isString {
			return true
		}
		n := dt.Length
		if !dt.HasLength {
			n = 1
		}
		if typeName == "char" {
			return utf8.RuneCountInString(text) <= n
		}
		return len(text) <= n
	case "varchar":
		return !isString || utf8.RuneCountInString(text) <= dt.Length
	case "varbinary":
		return !isString || len(text) <= dt.Length
	case "enum":
		if !isString {
			return true
		}
		return memberFold(dt.EnumValues, text)
	case "set":
		if !isString || text == "" {
			return true
		}
		for _, part := range strings.Split(text, ",") {
			if !memberFold(dt.EnumValues, part) {
				return false
			}
		}
		return true
	case "date", "datetime", "timestamp":
		if !isString {
			return true
		}
		return validDateText(text)
	case "time":
		if !isString {
			return true
		}
		return validTimeText(text)
	}
	return true
}

// literalText is a literal's text, whether it is a quoted string, and whether expr
// is a literal this check reads at all.
func literalText(expr nodes.ExprNode) (string, bool, bool) {
	switch e := expr.(type) {
	case *nodes.StringLit:
		return e.Value, true, true
	case *nodes.IntLit:
		if e.Text != "" {
			return e.Text, false, numericText.MatchString(e.Text)
		}
		return big.NewInt(e.Value).String(), false, true
	case *nodes.FloatLit:
		return e.Value, false, true
	case *nodes.UnaryExpr:
		if e.Op == nodes.UnaryMinus {
			if text, isString, ok := literalText(e.Operand); ok && !isString {
				return "-" + text, false, true
			}
		}
	}
	return "", false, false
}

// numberPrec is the mantissa a default's number is read and compared with: wide
// enough to hold any integer bound exactly, and fixed so a value's own size never
// sets it.
const numberPrec = 256

func literalNumber(text string) (*big.Float, bool) {
	text = strings.TrimSpace(text)
	if !numericText.MatchString(text) {
		return nil, false
	}
	n, _, err := big.ParseFloat(text, 10, numberPrec, big.ToNearestEven)
	return n, err == nil
}

// integerInRange reports whether n, rounded half away from zero, fits the type.
// The bounds are compared as floats and n is never materialized as an integer: a
// literal of a few bytes, '"'"'1e10000000'"'"', expands to megabytes as one.
func integerInRange(n *big.Float, typeName string, unsigned bool) bool {
	bits := map[string]uint{"tinyint": 8, "smallint": 16, "mediumint": 24, "int": 32, "integer": 32, "bigint": 64}[typeName]
	var lo, hi *big.Float
	if unsigned {
		lo = new(big.Float)
		hi = new(big.Float).SetInt(new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), bits), big.NewInt(1)))
	} else {
		hi = new(big.Float).SetInt(new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), bits-1), big.NewInt(1)))
		lo = new(big.Float).Neg(new(big.Float).SetInt(new(big.Int).Lsh(big.NewInt(1), bits-1)))
	}
	// Rounding half away from zero holds every value within half a unit of a bound,
	// so the bounds move by half, not n: adding to n allocates a mantissa as wide as
	// the gap between the exponents, megabytes for a literal like '1e100000000',
	// while the bounds are small and Cmp allocates nothing.
	// SetInt gives a bound only as many bits as it needs, so the half would round
	// away at a bound of 64 bits: -2^63 - 0.5 came back as -2^63 and refused
	// BIGINT DEFAULT '-9223372036854775808'. The sum is taken at the precision
	// literalNumber reads a value with, where 65 bits are exact.
	half := big.NewFloat(0.5)
	upper := new(big.Float).SetPrec(numberPrec).Add(hi, half)
	lower := new(big.Float).SetPrec(numberPrec).Sub(lo, half)
	return n.Cmp(upper) < 0 && n.Cmp(lower) > 0
}

// memberFold reports whether v is one of an ENUM's or SET's members. The members
// are decoded, as parseStringValueList leaves them, so a member that holds a quote
// is compared as it stands: ENUM('"'"'a') has the one member '"'"'a.
func memberFold(members []string, v string) bool {
	for _, m := range members {
		if strings.EqualFold(m, v) {
			return true
		}
	}
	return false
}

// validDateText reports whether a DATE, DATETIME or TIMESTAMP reads text: delimited
// or digits only, month to 12, day to 31, and a time within the day. A zero date is a
// date.
func validDateText(text string) bool {
	var parts []string
	if m := dateText.FindStringSubmatch(text); m != nil {
		parts = m[1:]
	} else if m := dateDigits.FindStringSubmatch(text); m != nil {
		parts = m[1:]
	} else {
		return false
	}
	limits := []int{9999, 12, 31, 23, 59, 59}
	for i, p := range parts {
		if p == "" {
			continue
		}
		v := 0
		for _, c := range p {
			v = v*10 + int(c-'0')
		}
		if v > limits[i] {
			return false
		}
	}
	return true
}

// validTimeText reports whether TIME reads text: [-][D ]H[:MM[:SS[.f]]] with minutes
// and seconds under 60, or digits.
func validTimeText(text string) bool {
	if timeDigits.MatchString(text) {
		return true
	}
	m := timeText.FindStringSubmatch(text)
	if m == nil {
		return false
	}
	for _, p := range m[1:] {
		if len(p) == 2 && p[0] >= '6' {
			return false
		}
	}
	return true
}

// checkAddedGeneratedCheck refuses a generated column with a CHECK in ALTER TABLE
// ... ADD COLUMN (1235), before anything else of the column is read.
func checkAddedGeneratedCheck(colDef *nodes.ColumnDef) error {
	if colDef == nil || colDef.Generated == nil {
		return nil
	}
	for _, cc := range colDef.Constraints {
		if cc.Type == nodes.ColConstrCheck {
			return &Error{Code: 1235, SQLState: "0A000", Message: "Add column-level check constraint while altering table not supported"}
		}
	}
	return nil
}

// checkAddedColumn refuses, before the name is looked up, what an ALTER TABLE ...
// ADD COLUMN may not add to tbl: a generated column reading a column tbl lacks, itself
// included (1054), and a PRIMARY KEY beside tbl's (1068).
func checkAddedColumn(tbl *Table, colDef *nodes.ColumnDef) error {
	if colDef.Generated != nil {
		for _, name := range exprColumnNames(colDef.Generated.Expr) {
			if tbl.GetColumn(name) == nil {
				return &Error{Code: ErrNoSuchColumn, SQLState: sqlState(ErrNoSuchColumn),
					Message: "Unknown column '" + name + "' in 'generated column function'"}
			}
		}
	}
	if columnDefHasPrimaryKey(colDef) {
		for _, idx := range tbl.Indexes {
			if idx.Primary {
				return errMultiplePriKey()
			}
		}
	}
	return nil
}
