package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/mysql/ast"
)

const (
	functionalIndexBaseName   = "functional_index"
	maxFunctionalColumnLength = 64
)

func hasFunctionalIndexColumns(idxCols []*IndexColumn) bool {
	for _, idxCol := range idxCols {
		if idxCol.Expr != "" {
			return true
		}
	}
	return false
}

func allocFunctionalIndexName(tbl *Table) string {
	if !indexNameExists(tbl, functionalIndexBaseName) {
		return functionalIndexBaseName
	}
	for suffix := 2; ; suffix++ {
		name := fmt.Sprintf("%s_%d", functionalIndexBaseName, suffix)
		if !indexNameExists(tbl, name) {
			return name
		}
	}
}

func synthesizeFunctionalIndexColumns(tbl *Table, idx *Index) error {
	for _, idxCol := range idx.Columns {
		if idxCol.Expr == "" {
			continue
		}
		if err := validateFunctionalIndexExpr(tbl, idx.Name, idxCol.ExprNode); err != nil {
			return err
		}
	}
	for part, idxCol := range idx.Columns {
		if idxCol.Expr == "" {
			continue
		}
		hiddenName := makeFunctionalIndexColumnName(tbl, idx.Name, part)
		idxCol.Name = hiddenName
		hiddenCol := columnFromFunctionalIndexExpr(tbl, hiddenName, idxCol)
		tbl.Columns = append(tbl.Columns, &Column{
			Position:   hiddenCol.Position,
			Name:       hiddenCol.Name,
			DataType:   hiddenCol.DataType,
			ColumnType: hiddenCol.ColumnType,
			Nullable:   hiddenCol.Nullable,
			Charset:    hiddenCol.Charset,
			Collation:  hiddenCol.Collation,
			Generated:  hiddenCol.Generated,
			Hidden:     hiddenCol.Hidden,
		})
		tbl.colByName[toLower(hiddenName)] = len(tbl.Columns) - 1
	}
	return nil
}

func columnFromFunctionalIndexExpr(tbl *Table, hiddenName string, idxCol *IndexColumn) *Column {
	rt := inferFunctionalIndexExprType(tbl, idxCol.ExprNode)
	col := columnFromResolvedType(rt)
	col.Position = len(tbl.Columns) + 1
	col.Name = hiddenName
	col.Nullable = true
	col.Generated = &GeneratedColumnInfo{
		Expr:   idxCol.Expr,
		Stored: false,
	}
	col.Hidden = ColumnHiddenSystem
	return col
}

func inferFunctionalIndexExprType(tbl *Table, expr nodes.ExprNode) *ResolvedType {
	switch e := expr.(type) {
	case *nodes.ParenExpr:
		return inferFunctionalIndexExprType(tbl, e.Expr)
	case *nodes.ColumnRef:
		if col := tbl.GetColumn(e.Column); col != nil {
			return resolvedTypeFromColumn(col)
		}
	case *nodes.BinaryExpr:
		switch e.Op {
		case nodes.BinOpAdd, nodes.BinOpSub, nodes.BinOpMul, nodes.BinOpDivInt, nodes.BinOpMod:
			return &ResolvedType{BaseType: BaseTypeBigInt}
		case nodes.BinOpDiv:
			return &ResolvedType{BaseType: BaseTypeDecimal}
		case nodes.BinOpJsonExtract:
			return &ResolvedType{BaseType: BaseTypeJSON}
		case nodes.BinOpJsonUnquote:
			return &ResolvedType{BaseType: BaseTypeText}
		}
	case *nodes.FuncCallExpr:
		name := strings.ToLower(e.Name)
		switch name {
		case "lower", "upper":
			if len(e.Args) > 0 {
				rt := inferFunctionalIndexExprType(tbl, e.Args[0])
				if rt != nil && isResolvedStringType(rt) {
					cp := *rt
					return &cp
				}
			}
			return &ResolvedType{BaseType: BaseTypeVarchar}
		default:
			args := make([]AnalyzedExpr, 0, len(e.Args))
			for _, arg := range e.Args {
				args = append(args, &ConstExprQ{Type: inferFunctionalIndexExprType(tbl, arg)})
			}
			return functionReturnType(name, args)
		}
	case *nodes.CastExpr:
		if e.TypeName != nil {
			return dataTypeToResolvedType(e.TypeName)
		}
	}
	return &ResolvedType{BaseType: BaseTypeVarchar, Length: 255}
}

func resolvedTypeFromColumn(col *Column) *ResolvedType {
	rt := &ResolvedType{
		Charset:   col.Charset,
		Collation: col.Collation,
	}
	switch strings.ToLower(col.DataType) {
	case "tinyint":
		rt.BaseType = BaseTypeTinyInt
	case "smallint":
		rt.BaseType = BaseTypeSmallInt
	case "mediumint":
		rt.BaseType = BaseTypeMediumInt
	case "int", "integer":
		rt.BaseType = BaseTypeInt
	case "bigint":
		rt.BaseType = BaseTypeBigInt
	case "varchar":
		rt.BaseType = BaseTypeVarchar
		if n, ok := fixedLengthStringColumnLimit(col); ok {
			rt.Length = n
		}
	case "char":
		rt.BaseType = BaseTypeChar
		if n, ok := fixedLengthStringColumnLimit(col); ok {
			rt.Length = n
		}
	case "json":
		rt.BaseType = BaseTypeJSON
	case "text", "tinytext", "mediumtext", "longtext":
		rt.BaseType = BaseTypeText
	default:
		rt.BaseType = BaseTypeVarchar
	}
	if strings.Contains(strings.ToLower(col.ColumnType), "unsigned") {
		rt.Unsigned = true
	}
	return rt
}

func columnFromResolvedType(rt *ResolvedType) *Column {
	if rt == nil {
		rt = &ResolvedType{BaseType: BaseTypeVarchar, Length: 255}
	}
	col := &Column{Nullable: true, Charset: rt.Charset, Collation: rt.Collation}
	switch rt.BaseType {
	case BaseTypeBigInt:
		col.DataType = "bigint"
		col.ColumnType = "bigint"
		if rt.Unsigned {
			col.ColumnType = "bigint unsigned"
		}
	case BaseTypeInt:
		col.DataType = "int"
		col.ColumnType = "int"
		if rt.Unsigned {
			col.ColumnType = "int unsigned"
		}
	case BaseTypeChar:
		col.DataType = "char"
		if rt.Length > 0 {
			col.ColumnType = fmt.Sprintf("char(%d)", rt.Length)
		} else {
			col.ColumnType = "char(1)"
		}
	case BaseTypeVarchar:
		col.DataType = "varchar"
		if rt.Length > 0 {
			col.ColumnType = fmt.Sprintf("varchar(%d)", rt.Length)
		} else {
			col.ColumnType = "varchar(255)"
		}
	case BaseTypeJSON:
		col.DataType = "json"
		col.ColumnType = "json"
	case BaseTypeText, BaseTypeTinyText, BaseTypeMediumText, BaseTypeLongText:
		col.DataType = "text"
		col.ColumnType = "text"
	case BaseTypeDecimal:
		col.DataType = "decimal"
		if rt.Precision > 0 {
			col.ColumnType = fmt.Sprintf("decimal(%d,%d)", rt.Precision, rt.Scale)
		} else {
			col.ColumnType = "decimal"
		}
	default:
		col.DataType = "varchar"
		col.ColumnType = "varchar(255)"
	}
	return col
}

func isResolvedStringType(rt *ResolvedType) bool {
	switch rt.BaseType {
	case BaseTypeChar, BaseTypeVarchar, BaseTypeTinyText, BaseTypeText, BaseTypeMediumText, BaseTypeLongText:
		return true
	default:
		return false
	}
}

func validateFunctionalIndexExpr(tbl *Table, indexName string, expr nodes.ExprNode) error {
	if _, ok := peelFunctionalIndexParens(expr).(*nodes.ColumnRef); ok {
		return &Error{
			Code:     3756,
			SQLState: "HY000",
			Message:  fmt.Sprintf("Functional index on a column is not supported. Consider using a regular index instead. Index '%s'.", indexName),
		}
	}
	if containsDisallowedFunctionalIndexExpr(expr) {
		return &Error{
			Code:     3757,
			SQLState: "HY000",
			Message:  fmt.Sprintf("Expression of functional index '%s' contains a disallowed function.", indexName),
		}
	}
	if isFunctionalIndexLOBType(inferFunctionalIndexExprType(tbl, expr)) {
		return &Error{
			Code:     3754,
			SQLState: "HY000",
			Message:  "Cannot create a functional index on an expression that returns a BLOB or TEXT. Please consider using CAST.",
		}
	}
	return nil
}

func peelFunctionalIndexParens(expr nodes.ExprNode) nodes.ExprNode {
	for {
		p, ok := expr.(*nodes.ParenExpr)
		if !ok {
			return expr
		}
		expr = p.Expr
	}
}

func containsDisallowedFunctionalIndexExpr(expr nodes.ExprNode) bool {
	switch e := expr.(type) {
	case nil:
		return false
	case *nodes.FuncCallExpr:
		switch strings.ToLower(e.Name) {
		case "rand", "uuid", "sleep", "now", "current_timestamp", "sysdate":
			return true
		}
		for _, arg := range e.Args {
			if containsDisallowedFunctionalIndexExpr(arg) {
				return true
			}
		}
	case *nodes.DefaultExpr, *nodes.SubqueryExpr, *nodes.ExistsExpr, *nodes.VariableRef:
		return true
	case *nodes.ParenExpr:
		return containsDisallowedFunctionalIndexExpr(e.Expr)
	case *nodes.BinaryExpr:
		return containsDisallowedFunctionalIndexExpr(e.Left) || containsDisallowedFunctionalIndexExpr(e.Right)
	case *nodes.UnaryExpr:
		return containsDisallowedFunctionalIndexExpr(e.Operand)
	case *nodes.CastExpr:
		return containsDisallowedFunctionalIndexExpr(e.Expr)
	case *nodes.CaseExpr:
		if containsDisallowedFunctionalIndexExpr(e.Operand) || containsDisallowedFunctionalIndexExpr(e.Default) {
			return true
		}
		for _, when := range e.Whens {
			if containsDisallowedFunctionalIndexExpr(when.Cond) || containsDisallowedFunctionalIndexExpr(when.Result) {
				return true
			}
		}
	case *nodes.InExpr:
		if e.Select != nil {
			return true
		}
		if containsDisallowedFunctionalIndexExpr(e.Expr) {
			return true
		}
		for _, item := range e.List {
			if containsDisallowedFunctionalIndexExpr(item) {
				return true
			}
		}
	case *nodes.BetweenExpr:
		return containsDisallowedFunctionalIndexExpr(e.Expr) ||
			containsDisallowedFunctionalIndexExpr(e.Low) ||
			containsDisallowedFunctionalIndexExpr(e.High)
	case *nodes.IsExpr:
		return containsDisallowedFunctionalIndexExpr(e.Expr)
	}
	return false
}

func isFunctionalIndexLOBType(rt *ResolvedType) bool {
	if rt == nil {
		return false
	}
	switch rt.BaseType {
	case BaseTypeJSON, BaseTypeTinyText, BaseTypeText, BaseTypeMediumText, BaseTypeLongText,
		BaseTypeTinyBlob, BaseTypeBlob, BaseTypeMediumBlob, BaseTypeLongBlob,
		BaseTypeGeometry, BaseTypePoint, BaseTypeLineString, BaseTypePolygon,
		BaseTypeMultiPoint, BaseTypeMultiLineString, BaseTypeMultiPolygon, BaseTypeGeometryCollection:
		return true
	default:
		return false
	}
}

func makeFunctionalIndexColumnName(tbl *Table, indexName string, part int) string {
	for count := 0; ; count++ {
		suffix := fmt.Sprintf("!%d!%d", part, count)
		prefix := "!hidden!" + indexName
		maxPrefixLen := maxFunctionalColumnLength - len(suffix)
		if len(prefix) > maxPrefixLen {
			prefix = prefix[:maxPrefixLen]
		}
		name := prefix + suffix
		if tbl.GetColumn(name) == nil {
			return name
		}
	}
}

func defaultIndexName(tbl *Table, cols []string, idxCols []*IndexColumn) string {
	if len(cols) > 0 {
		return allocIndexName(tbl, cols[0])
	}
	if hasFunctionalIndexColumns(idxCols) {
		return allocFunctionalIndexName(tbl)
	}
	return ""
}

func isFunctionalIndex(idx *Index) bool {
	for _, idxCol := range idx.Columns {
		if idxCol.Expr != "" {
			return true
		}
	}
	return false
}

func normalizeFunctionalIndexName(name string) string {
	return strings.TrimSpace(name)
}
