package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

func validateExpressionSemantics(node nodes.Node) error {
	var err error
	nodes.Inspect(node, func(n nodes.Node) bool {
		if err != nil {
			return false
		}
		fn, ok := n.(*nodes.FuncCallExpr)
		if !ok {
			return true
		}
		name := strings.ToLower(fn.Name)
		switch name {
		case "now", "current_timestamp", "localtime", "localtimestamp", "sysdate",
			"utc_timestamp", "curtime", "current_time", "utc_time":
			if fsp, ok := firstIntArg(fn); ok && fsp > 6 {
				err = errTooBigPrecision(fsp, name)
			}
		case "curdate", "current_date", "utc_date":
			if len(fn.Args) > 0 {
				err = errWrongArguments(strings.ToUpper(name))
			}
		}
		return err == nil
	})
	return err
}

func validateColumnDefSemantics(colDef *nodes.ColumnDef) error {
	if colDef == nil || colDef.TypeName == nil {
		return nil
	}
	typeName := normalizedColumnDataType(colDef.TypeName)
	switch typeName {
	case "geometry", "point", "linestring", "polygon", "multipoint", "multilinestring", "multipolygon", "geometrycollection":
		if isLiteralDefault(colDef.DefaultValue) || hasLiteralColumnDefaultConstraint(colDef) {
			return &Error{Code: 1416, SQLState: "22003", Message: "Cannot get geometry object from data you send to the GEOMETRY field"}
		}
	}

	switch typeName {
	case "datetime", "timestamp", "time":
		if colDef.TypeName.Length > 6 {
			return errTooBigPrecision(colDef.TypeName.Length, typeName)
		}
	case "decimal", "numeric", "dec", "fixed":
		if colDef.TypeName.Length > 65 {
			return &Error{Code: 1426, SQLState: "42000", Message: fmt.Sprintf("Too big precision specified for '%s'. Maximum is 65", colDef.Name)}
		}
		if colDef.TypeName.Scale > 38 {
			return &Error{Code: 1425, SQLState: "42000", Message: fmt.Sprintf("Too big scale specified for '%s'. Maximum is 38", colDef.Name)}
		}
		if colDef.TypeName.Scale > colDef.TypeName.Length && colDef.TypeName.Length > 0 {
			return &Error{Code: 1427, SQLState: "42000", Message: "For float(M,D), double(M,D) or decimal(M,D), M must be >= D"}
		}
	case "varchar":
		if colDef.TypeName.Length == 0 {
			return &Error{Code: 1064, SQLState: "42000", Message: "VARCHAR requires a length"}
		}
	}
	if colDef.Generated != nil && (colDef.DefaultValue != nil || hasColumnDefaultConstraint(colDef)) {
		return errInvalidDefault(colDef.Name)
	}
	if colDef.Generated != nil {
		if err := validateGeneratedExpr(colDef.Name, colDef.Generated.Expr, !colDef.Generated.Stored); err != nil {
			return err
		}
		if !colDef.Generated.Stored && columnDefHasPrimaryKey(colDef) {
			return errPrimaryKeyOnGeneratedColumn(colDef.Name)
		}
	}
	if !columnDefNullable(colDef) && hasNullDefault(colDef) {
		return errInvalidDefault(colDef.Name)
	}
	if err := validateTemporalExprForColumn(colDef.Name, typeName, colDef.TypeName.Length, colDef.DefaultValue, true); err != nil {
		return err
	}
	if err := validateTemporalExprForColumn(colDef.Name, typeName, colDef.TypeName.Length, colDef.OnUpdate, false); err != nil {
		return err
	}
	for _, cc := range colDef.Constraints {
		if cc.Type == nodes.ColConstrDefault {
			if err := validateTemporalExprForColumn(colDef.Name, typeName, colDef.TypeName.Length, cc.Expr, true); err != nil {
				return err
			}
		}
	}
	return nil
}

func generatedExprNondeterministic(expr nodes.ExprNode) bool {
	invalid := false
	nodes.Inspect(expr, func(node nodes.Node) bool {
		if fn, ok := node.(*nodes.FuncCallExpr); ok && isNondeterministicFunc(fn.Name) {
			invalid = true
		}
		return !invalid
	})
	return invalid
}

func columnDefHasPrimaryKey(colDef *nodes.ColumnDef) bool {
	for _, cc := range colDef.Constraints {
		if cc.Type == nodes.ColConstrPrimaryKey {
			return true
		}
	}
	return false
}

func hasExplicitNullPrimaryKey(colDef *nodes.ColumnDef) bool {
	hasNull := false
	hasPK := false
	for _, cc := range colDef.Constraints {
		switch cc.Type {
		case nodes.ColConstrNull:
			hasNull = true
		case nodes.ColConstrPrimaryKey:
			hasPK = true
		}
	}
	return hasNull && hasPK
}

func columnDefNullable(colDef *nodes.ColumnDef) bool {
	nullable := true
	if colDef.AutoIncrement {
		nullable = false
	}
	for _, cc := range colDef.Constraints {
		switch cc.Type {
		case nodes.ColConstrNotNull, nodes.ColConstrPrimaryKey, nodes.ColConstrAutoIncrement:
			nullable = false
		case nodes.ColConstrNull:
			nullable = true
		}
	}
	return nullable
}

func validateGeneratedExpr(colName string, expr nodes.ExprNode, allowNondeterministic bool) error {
	var invalid bool
	nodes.Inspect(expr, func(n nodes.Node) bool {
		if invalid {
			return false
		}
		switch x := n.(type) {
		case *nodes.SubqueryExpr:
			invalid = true
		case *nodes.VariableRef:
			invalid = true
		case *nodes.FuncCallExpr:
			if !allowNondeterministic && isNondeterministicFunc(x.Name) {
				invalid = true
			}
		}
		return !invalid
	})
	if invalid {
		return &Error{Code: 1901, SQLState: "HY000", Message: "Expression of generated column '" + colName + "' contains a disallowed function"}
	}
	return nil
}

func errPrimaryKeyOnGeneratedColumn(colName string) error {
	return &Error{Code: 3106, SQLState: "HY000", Message: "Defining a virtual generated column as primary key is not supported"}
}

func validatePrimaryKeyColumns(tbl *Table, cols []string, parts ...*IndexColumn) error {
	for _, colName := range cols {
		col := tbl.GetColumn(colName)
		if col != nil && isTextBlobType(col.DataType) {
			prefixed := false
			for _, part := range parts {
				if strings.EqualFold(part.Name, colName) && part.Length > 0 {
					prefixed = true
				}
			}
			if !prefixed {
				return &Error{Code: 1170, SQLState: "42000", Message: "BLOB/TEXT column '" + colName + "' used in key specification without a key length"}
			}
		}
		if col != nil && col.Generated != nil && !col.Generated.Stored {
			return errPrimaryKeyOnGeneratedColumn(colName)
		}
	}
	return nil
}

func validateColumnCheckReferences(colName, conName string, expr nodes.ExprNode) error {
	// MariaDB allows another column of the same table in an inline CHECK.
	return validateCheckExpr(conName, expr)
}

func validateCheckColumnReferences(tbl *Table, expr nodes.ExprNode) error {
	var err error
	nodes.Inspect(expr, func(node nodes.Node) bool {
		if ref, ok := node.(*nodes.ColumnRef); ok && tbl.GetColumn(ref.Column) == nil {
			err = errNoSuchColumn(ref.Column, "CHECK")
		}
		return err == nil
	})
	return err
}

func hasColumnDefaultConstraint(colDef *nodes.ColumnDef) bool {
	for _, cc := range colDef.Constraints {
		if cc.Type == nodes.ColConstrDefault && cc.Expr != nil {
			return true
		}
	}
	return false
}

func hasNullDefault(colDef *nodes.ColumnDef) bool {
	if isNullDefault(colDef.DefaultValue) {
		return true
	}
	for _, cc := range colDef.Constraints {
		if cc.Type == nodes.ColConstrDefault && isNullDefault(cc.Expr) {
			return true
		}
	}
	return false
}

func hasLiteralColumnDefaultConstraint(colDef *nodes.ColumnDef) bool {
	for _, cc := range colDef.Constraints {
		if cc.Type == nodes.ColConstrDefault && isLiteralDefault(cc.Expr) {
			return true
		}
	}
	return false
}

func isLiteralDefault(expr nodes.ExprNode) bool {
	switch e := expr.(type) {
	case *nodes.StringLit, *nodes.IntLit, *nodes.FloatLit, *nodes.BitLit, *nodes.BoolLit:
		return true
	case *nodes.ParenExpr:
		return isLiteralDefault(e.Expr)
	default:
		return false
	}
}

func isNullDefault(expr nodes.ExprNode) bool {
	switch e := expr.(type) {
	case *nodes.NullLit:
		return true
	case *nodes.ParenExpr:
		return isNullDefault(e.Expr)
	default:
		return false
	}
}

func isLiteralDefaultForbiddenType(typeName string) bool {
	switch strings.ToLower(typeName) {
	case "blob", "tinyblob", "mediumblob", "longblob",
		"text", "tinytext", "mediumtext", "longtext",
		"json", "geometry", "point", "linestring", "polygon",
		"multipoint", "multilinestring", "multipolygon", "geometrycollection", "geomcollection":
		return true
	default:
		return false
	}
}

func fmtCantHaveDefault(colName string) string {
	return "BLOB, TEXT, GEOMETRY or JSON column '" + colName + "' can't have a default value"
}

func validateTemporalExprForColumn(colName, typeName string, colFsp int, expr nodes.ExprNode, isDefault bool) error {
	if expr == nil {
		return nil
	}
	if err := validateExpressionSemantics(expr); err != nil {
		return err
	}
	_, ok := unwrapTemporalFunc(expr)
	if !ok {
		return nil
	}
	if !isDefault && typeName != "datetime" && typeName != "timestamp" {
		return errInvalidOnUpdate(colName)
	}

	return nil
}

func validateCheckExpr(name string, expr nodes.ExprNode) error {
	var invalid bool
	exprSQL := strings.ToLower(nodeToSQL(expr))
	if strings.Contains(exprSQL, "select ") || strings.Contains(exprSQL, "(select") {
		return errCheckConstraintNotAllowed(name)
	}
	nodes.Inspect(expr, func(n nodes.Node) bool {
		if invalid {
			return false
		}
		switch x := n.(type) {
		case *nodes.SubqueryExpr:
			invalid = true
		case *nodes.VariableRef:
			invalid = true
		case *nodes.FuncCallExpr:
			if isNondeterministicFunc(x.Name) {
				invalid = true
			}
		}
		return !invalid
	})
	if invalid {
		return errCheckConstraintNotAllowed(name)
	}
	return nil
}

func isNondeterministicFunc(name string) bool {
	switch strings.ToLower(name) {
	case "now", "current_timestamp", "localtime", "localtimestamp", "sysdate",
		"utc_timestamp", "curtime", "current_time", "utc_time", "curdate",
		"current_date", "utc_date", "rand", "uuid", "uuid_short", "connection_id":
		return true
	default:
		return false
	}
}

func firstIntArg(fn *nodes.FuncCallExpr) (int, bool) {
	if fn == nil || len(fn.Args) == 0 {
		return 0, false
	}
	if lit, ok := fn.Args[0].(*nodes.IntLit); ok {
		return int(lit.Value), true
	}
	return 0, false
}

func temporalFuncFSP(fn *nodes.FuncCallExpr) (int, bool) {
	if fn == nil {
		return 0, false
	}
	name := strings.ToLower(fn.Name)
	switch name {
	case "now", "current_timestamp", "localtime", "localtimestamp", "sysdate",
		"utc_timestamp":
		if fsp, ok := firstIntArg(fn); ok {
			return fsp, true
		}
		return 0, true
	default:
		return 0, false
	}
}

func unwrapTemporalFunc(expr nodes.ExprNode) (*nodes.FuncCallExpr, bool) {
	switch e := expr.(type) {
	case *nodes.FuncCallExpr:
		return e, true
	case *nodes.ParenExpr:
		return unwrapTemporalFunc(e.Expr)
	default:
		return nil, false
	}
}
