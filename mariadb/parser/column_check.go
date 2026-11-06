package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

// What MariaDB refuses of a column while it parses. The action for a field definition runs
// Column_definition::check as the statement parses, so what would be a later check on
// another engine is a parse error here, with a code of its own: 1063 for AUTO_INCREMENT on
// a type that does not count, 1067 for a default the column cannot have, 1294 for ON UPDATE
// on a type that has no clock. A routine's types pass the length checks of refuseTypeLimits
// the same way, and CAST has its own in castTypeRefused.
//
// Ref: Column_definition::check in MariaDB 11.8 sql/field.cc, and sql/sql_type.cc for the
// attributes per type.

// startsVcolAttribute reports what vcol_attribute lists, which is all a generated column
// takes after its expression.
func (p *Parser) startsVcolAttribute() bool {
	switch p.cur.Type {
	case kwUNIQUE, kwCOMMENT, kwCHECK, kwCONSTRAINT:
		return true
	}
	return strings.EqualFold(p.cur.Str, "INVISIBLE") && p.cur.Type != tokSCONST
}

// refuseColumnAttributes holds a column's attributes to its type and to each other.
func refuseColumnAttributes(col *nodes.ColumnDef, pos int) error {
	name := ""
	if col.TypeName != nil {
		name = col.TypeName.Name
	}
	refuse := func(msg string) error {
		return &ParseError{Message: msg, Position: pos}
	}
	if col.AutoIncrement && name != "" {
		switch name {
		case "TINYINT", "SMALLINT", "MEDIUMINT", "INT", "BIGINT", "FLOAT", "DOUBLE", "BOOL", "BOOLEAN":
		default:
			return refuse("Incorrect column specifier for column '" + col.Name + "'") // 1063
		}
	}
	if col.OnUpdate != nil && name != "" {
		if name != "TIMESTAMP" && name != "DATETIME" {
			return refuse("Invalid ON UPDATE clause for '" + col.Name + "' column") // 1294
		}
		if call, ok := col.OnUpdate.(*nodes.FuncCallExpr); ok && len(call.Args) == 1 {
			if lit, ok := call.Args[0].(*nodes.IntLit); ok && int(lit.Value) < col.TypeName.Length {
				return refuse("Invalid ON UPDATE clause for '" + col.Name + "' column") // 1294
			}
		}
	}
	notNull := false
	keyed := false
	var def nodes.ExprNode
	for _, c := range col.Constraints {
		switch c.Type {
		case nodes.ColConstrNotNull:
			notNull = true
		case nodes.ColConstrNull:
			notNull = false
		case nodes.ColConstrPrimaryKey:
			keyed = true
		case nodes.ColConstrDefault:
			def = c.Expr
		}
	}
	if def == nil {
		return nil
	}
	if isNullDefault(def) {
		if notNull || keyed {
			return refuse("Invalid default value for '" + col.Name + "'") // 1067
		}
		return nil
	}
	if col.AutoIncrement {
		return refuse("Invalid default value for '" + col.Name + "'") // 1067
	}
	return nil
}

func isNullDefault(def nodes.ExprNode) bool {
	if paren, ok := def.(*nodes.ParenExpr); ok {
		def = paren.Expr
	}
	_, ok := def.(*nodes.NullLit)
	return ok
}

// refuseColumnCollationConflict is COLLATE as an attribute after a collation already
// written, beside the type or as another attribute: the same one again is fine, and two
// collations are 1302, ER_CONFLICTING_DECLARATIONS. A name the server would have to look
// up answers first with its own 1273, which reads.
func refuseColumnCollationConflict(col *nodes.ColumnDef, name string, pos int) error {
	prior := ""
	if col.TypeName != nil {
		prior = col.TypeName.Collate
	}
	for _, c := range col.Constraints {
		if c.Type == nodes.ColConstrCollate {
			prior = c.Name
		}
	}
	if prior == "" || strings.EqualFold(prior, name) {
		return nil
	}
	if collationCharset(prior) == "" || collationCharset(name) == "" {
		return nil
	}
	return &ParseError{Message: "Conflicting declarations: 'COLLATE " + prior + "' and 'COLLATE " + name + "'", Position: pos}
}

// refuseTypeLimits holds a type's length and precision to what MariaDB allows, which the
// parser's actions check for a column, a routine's parameter, return value and variable
// alike: 1439 for a display width or length, 1074 for a CHAR or BINARY past 255, 1426 for a
// precision, 1425 for a scale and 1427 for a scale past its precision, 1063 for FLOAT(p) past
// a double's 53.
//
// Ref: Column_definition::fix_attributes_* in MariaDB 11.8 sql/sql_type.cc and sql/field.cc.
func refuseTypeLimits(dt *nodes.DataType, hasScale bool, pos int) error {
	refuse := func(msg string) error {
		return &ParseError{Message: msg, Position: pos}
	}
	const maxLength = 4294967295
	switch dt.Name {
	case "TINYINT", "SMALLINT", "MEDIUMINT", "INT", "BIGINT":
		if dt.Length > 255 {
			return refuse("Display width out of range")
		}
	case "BIT":
		if dt.Length > 64 {
			return refuse("Display width out of range")
		}
	case "CHAR", "BINARY":
		if dt.Length > 255 {
			return refuse("Column length too big")
		}
	case "VARCHAR", "VARBINARY", "TEXT", "BLOB":
		if dt.Length > maxLength {
			return refuse("Display width out of range")
		}
	case "DECIMAL", "NUMERIC":
		switch {
		case dt.Length > 65:
			return refuse("Too-big precision")
		case dt.Scale > 38:
			return refuse("Too big scale")
		case hasScale && dt.Scale > dt.Length:
			return refuse("For float(M,D), double(M,D) or decimal(M,D), M must be >= D")
		}
	case "FLOAT", "DOUBLE":
		switch {
		case !hasScale && dt.Length > 53:
			return refuse("Incorrect column specifier")
		case hasScale && dt.Length > 255:
			return refuse("Display width out of range")
		case hasScale && (dt.Scale > 30 || dt.Scale > dt.Length):
			return refuse("For float(M,D), double(M,D) or decimal(M,D), M must be >= D")
		}
	case "TIME", "DATETIME", "TIMESTAMP":
		if dt.Length > 6 {
			return refuse("Too-big precision")
		}
	}
	return nil
}
