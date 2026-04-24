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

func synthesizeFunctionalIndexColumns(tbl *Table, idx *Index) {
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
