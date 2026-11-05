package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/oceanbase/ast"
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
		idxCol.Expr = normalizeFunctionalIndexExprSQL(idxCol.Expr, tbl.Charset)
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

func removeFunctionalIndexHiddenColumns(tbl *Table, idx *Index) {
	remove := make(map[string]bool)
	for _, idxCol := range idx.Columns {
		if idxCol.Expr == "" || idxCol.Name == "" {
			continue
		}
		remove[toLower(idxCol.Name)] = true
	}
	if len(remove) == 0 {
		return
	}

	cols := tbl.Columns[:0]
	for _, col := range tbl.Columns {
		if remove[toLower(col.Name)] && col.Hidden == ColumnHiddenSystem {
			continue
		}
		cols = append(cols, col)
	}
	tbl.Columns = cols
	rebuildColIndex(tbl)
}

func renameFunctionalIndexHiddenColumns(tbl *Table, idx *Index, newIndexName string) {
	renamed := false
	for part, idxCol := range idx.Columns {
		if idxCol.Expr == "" || idxCol.Name == "" {
			continue
		}
		col := tbl.GetColumn(idxCol.Name)
		if col == nil || col.Hidden != ColumnHiddenSystem {
			continue
		}
		newName := makeFunctionalIndexColumnName(tbl, newIndexName, part)
		col.Name = newName
		idxCol.Name = newName
		renamed = true
	}
	if renamed {
		rebuildColIndex(tbl)
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
		// JSON_VALUE answers its RETURNING type, and char(512) with none: mysql:8.0.46 gives
		// KEY ((JSON_VALUE(a,'$.x' RETURNING SIGNED))) a bigint hidden column, and one with
		// no RETURNING a varchar(512).
		if e.JSONValueArm {
			if e.JSONValue != nil && e.JSONValue.Returning != nil {
				return jsonValueIndexType(e.JSONValue.Returning)
			}
			return &ResolvedType{BaseType: BaseTypeVarchar, Length: 512}
		}
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

// jsonValueIndexType is the hidden column a functional index over JSON_VALUE ... RETURNING gets:
// a CHAR(n) is a varchar(n), and a temporal type keeps its precision.
//
// A string longer than 512 characters, which is what one with no length is, is a TEXT or a BLOB,
// and an index over it is refused: mysql:8.0.46 takes RETURNING CHAR(512) and answers 3757 for
// CHAR(513) and for CHAR.
func jsonValueIndexType(dt *nodes.DataType) *ResolvedType {
	if strings.EqualFold(dt.Name, "year") {
		return &ResolvedType{BaseType: BaseTypeYear}
	}
	rt := dataTypeToResolvedType(dt)
	// CHAR ... CHARSET binary is a binary string: mysql:8.0.46 gives RETURNING CHAR(3) CHARSET
	// binary a varbinary(3) hidden column.
	if rt.BaseType == BaseTypeChar && strings.EqualFold(dt.Charset, "binary") {
		rt.BaseType = BaseTypeBinary
	}
	switch rt.BaseType {
	case BaseTypeChar, BaseTypeBinary:
		long := !dt.HasLength || dt.Length > 512
		switch {
		case long && rt.BaseType == BaseTypeChar:
			rt.BaseType = BaseTypeText
		case long:
			rt.BaseType = BaseTypeBlob
		case rt.BaseType == BaseTypeChar:
			rt.BaseType = BaseTypeVarchar
		default:
			rt.BaseType = BaseTypeVarBinary
		}
		rt.Length = dt.Length
	case BaseTypeDateTime, BaseTypeTime:
		rt.Precision = dt.Length
	case BaseTypeDecimal:
		if rt.Precision == 0 {
			rt.Precision = 10
		}
	}
	return rt
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
	case BaseTypeVarBinary:
		col.DataType = "varbinary"
		col.ColumnType = fmt.Sprintf("varbinary(%d)", rt.Length)
	case BaseTypeDate:
		col.DataType, col.ColumnType = "date", "date"
	case BaseTypeDateTime, BaseTypeTime:
		col.DataType = "datetime"
		if rt.BaseType == BaseTypeTime {
			col.DataType = "time"
		}
		col.ColumnType = col.DataType
		if rt.Precision > 0 {
			col.ColumnType = fmt.Sprintf("%s(%d)", col.DataType, rt.Precision)
		}
	case BaseTypeDouble:
		col.DataType, col.ColumnType = "double", "double"
	case BaseTypeFloat:
		col.DataType, col.ColumnType = "float", "float"
	case BaseTypeYear:
		col.DataType, col.ColumnType = "year", "year"
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
	// A JSON or geometry result has a code of its own: mysql:8.0.46 answers 3753 for
	// KEY ((JSON_VALUE(a,'$.x' RETURNING JSON))) and KEY ((CAST(a->'$.x' AS JSON))).
	if rt := inferFunctionalIndexExprType(tbl, expr); rt != nil && isFunctionalIndexJSONOrGeometryType(rt) {
		return &Error{
			Code:     3753,
			SQLState: "HY000",
			Message:  "Cannot create a functional index on a function that returns a JSON or GEOMETRY value.",
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
	case *nodes.PositionExpr:
		// KEY ((position('x' in rand()))) is rejected exactly like
		// KEY ((rand())), and position over a column is allowed (MySQL 8.0,
		// whose functional-index rules this validator follows).
		return containsDisallowedFunctionalIndexExpr(e.Substr) ||
			containsDisallowedFunctionalIndexExpr(e.Str)
	case *nodes.IsExpr:
		return containsDisallowedFunctionalIndexExpr(e.Expr)
	case *nodes.WeightStringExpr:
		// KEY ((weight_string(now() as char(4)))) is rejected by the engine
		// exactly like KEY ((now())) — error 3758 (oracle 8.0.32).
		return containsDisallowedFunctionalIndexExpr(e.Expr)
	case *nodes.ExtractExpr:
		// KEY ((extract(day from now()))) → 3758; extract over a column is
		// allowed (oracle 8.0.32).
		return containsDisallowedFunctionalIndexExpr(e.Expr)
	case *nodes.IntervalExpr:
		// KEY ((a + interval (day(now())) day)) → 3758 (oracle 8.0.32).
		return containsDisallowedFunctionalIndexExpr(e.Value)
	}
	return false
}

func isFunctionalIndexJSONOrGeometryType(rt *ResolvedType) bool {
	switch rt.BaseType {
	case BaseTypeJSON, BaseTypeGeometry, BaseTypePoint, BaseTypeLineString, BaseTypePolygon,
		BaseTypeMultiPoint, BaseTypeMultiLineString, BaseTypeMultiPolygon, BaseTypeGeometryCollection:
		return true
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

func normalizeFunctionalIndexExprSQL(sql, charset string) string {
	if charset == "" {
		charset = "utf8mb4"
	}
	return strings.ReplaceAll(sql, ",'$", ",_"+charset+"'$")
}
