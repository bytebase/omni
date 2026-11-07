package catalog

import (
	"fmt"
	"strconv"
	"strings"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// TiDB builds expression-index key parts as hidden generated expressions, then
// checks their result types before installing the index. Keep the parsed node:
// catalog validation must not reparse a rendered expression.
// https://github.com/pingcap/tidb/blob/v8.5.5/pkg/ddl/create_table.go
// https://github.com/pingcap/tidb/blob/v8.5.5/pkg/ddl/index.go
func (c *Catalog) prepareTiDBIndexExpressions(tbl *Table, idx *Index) error {
	if idx.Vector {
		return nil
	}
	for _, part := range idx.Columns {
		if part.ExprNode == nil {
			continue // Metadata-loaded indexes already passed the engine's checks.
		}
		if tiDBIllegalGeneratedExpr(part.ExprNode) {
			return &Error{Code: 3758, SQLState: "HY000", Message: fmt.Sprintf("Expression of expression index '%s' contains a disallowed function", idx.Name)}
		}
		var aggregate, unsafe bool
		var unknownColumn string
		nodes.Inspect(part.ExprNode, func(node nodes.Node) bool {
			switch n := node.(type) {
			case *nodes.FuncCallExpr:
				aggregate = aggregate || isAggregateFunc(n.Name)
				unsafe = unsafe || !tiDBGAIndexFunction(n.Name)
			case *nodes.ExtractExpr, *nodes.PositionExpr:
				unsafe = true
			case *nodes.ColumnRef:
				if tbl.GetColumn(n.Column) == nil || n.Table != "" && !eqFoldStr(n.Table, tbl.Name) {
					unknownColumn = n.Column
				}
			}
			return true
		})
		if aggregate {
			return &Error{Code: 1111, SQLState: "HY000", Message: "Invalid use of group function"}
		}
		if unsafe && !c.allowExpressionIndex {
			return &Error{Code: 8200, SQLState: "HY000", Message: "Unsupported creating expression index containing unsafe functions without allow-expression-index in config"}
		}
		if unknownColumn != "" {
			return errNoSuchColumn(unknownColumn, "expression")
		}
		expr := part.ExprNode
		for {
			paren, ok := expr.(*nodes.ParenExpr)
			if !ok {
				break
			}
			expr = paren.Expr
		}
		if _, bare := expr.(*nodes.ColumnRef); bare {
			return &Error{Code: 3762, SQLState: "HY000", Message: "Expression index on a column is not supported. Consider using a regular index instead"}
		}
		part.ExprType = tiDBIndexExpressionType(part.ExprNode, tbl)
		if cast, ok := expr.(*nodes.CastExpr); !ok || !cast.Array {
			if part.ExprType != nil {
				switch part.ExprType.BaseType {
				case BaseTypeJSON, BaseTypeGeometry:
					return &Error{Code: 3753, SQLState: "HY000", Message: "Cannot create an expression index on a function that returns a JSON or GEOMETRY value"}
				case BaseTypeTinyText, BaseTypeText, BaseTypeMediumText, BaseTypeLongText, BaseTypeTinyBlob, BaseTypeBlob, BaseTypeMediumBlob, BaseTypeLongBlob:
					return &Error{Code: 3757, SQLState: "HY000", Message: "Cannot create an expression index on an expression that returns a BLOB or TEXT. Please consider using CAST"}
				}
			}
			if typ := part.ExprType; typ != nil && typ.BaseType == BaseTypeVarchar {
				width := 4
				charset := typ.Charset
				if charset == "" {
					charset = tbl.Charset
				}
				switch strings.ToLower(charset) {
				case "ascii", "latin1", "binary":
					width = 1
				case "utf8", "utf8mb3":
					width = 3
				case "gbk":
					width = 2
				}
				if typ.Length > 3072/width {
					return &Error{Code: 1071, SQLState: "42000", Message: "Specified key was too long; max key length is 3072 bytes"}
				}
			}
		}
		part.Expr = storedTiDBExpressionSQL(part.ExprNode, c.defaultCharset, true)
	}
	return c.prepareTiDBIndexLengths(tbl, idx)
}

// getIndexColumnLength/buildIndexColumns in TiDB count encoded key bytes,
// including all parts of a composite index, rather than character counts.
// https://github.com/pingcap/tidb/blob/v8.5.5/pkg/ddl/index.go
func (c *Catalog) prepareTiDBIndexLengths(tbl *Table, idx *Index) error {
	bytes := 0
	for _, part := range idx.Columns {
		if part.Name == "" {
			continue
		}
		col := tbl.GetColumn(part.Name)
		if col == nil {
			return &Error{Code: 1072, SQLState: "42000", Message: "column does not exist: " + part.Name}
		}
		typeName := strings.ToLower(col.DataType)
		length, width := 0, 1
		prefixable := false
		switch typeName {
		case "json":
			return &Error{Code: 3152, SQLState: "42000", Message: fmt.Sprintf("JSON column '%s' supports indexing only via generated columns on a specified JSON path.", col.Name)}
		case "tinytext", "text", "mediumtext", "longtext", "tinyblob", "blob", "mediumblob", "longblob":
			if part.Length == 0 {
				return &Error{Code: 1170, SQLState: "42000", Message: fmt.Sprintf("BLOB/TEXT column '%s' used in key specification without a key length", col.Name)}
			}
			prefixable, length = true, part.Length
		case "char", "varchar", "binary", "varbinary":
			prefixable = true
			if _, after, ok := strings.Cut(col.ColumnType, "("); ok {
				if value, _, ok := strings.Cut(after, ")"); ok {
					if parsed, err := strconv.Atoi(value); err == nil {
						length = parsed
					}
				}
			}
			if part.Length > length {
				return &Error{Code: 1089, SQLState: "HY000", Message: "Incorrect prefix key"}
			}
			if part.Length == length {
				part.Length = 0
			} else if part.Length > 0 {
				length = part.Length
			}
		case "tinyint", "year":
			length = 1
		case "smallint", "enum":
			length = 2
		case "mediumint", "date", "time":
			length = 3
		case "int", "integer", "float", "timestamp":
			length = 4
		case "bigint", "double", "datetime", "set":
			length = 8
		}
		if part.Length > 0 && !prefixable {
			return &Error{Code: 1089, SQLState: "HY000", Message: "Incorrect prefix key"}
		}
		if prefixable {
			width = tiDBCharsetWidth(col.Charset)
		}
		bytes += length * width
		if bytes > 3072 {
			if c.strictDDL || idx.Unique || idx.Primary || len(idx.Columns) > 1 {
				return &Error{Code: 1071, SQLState: "42000", Message: "Specified key was too long; max key length is 3072 bytes"}
			}
			part.Length = 3072 / width
		}
	}
	return nil
}

func tiDBCharsetWidth(charset string) int {
	switch strings.ToLower(charset) {
	case "utf8mb4":
		return 4
	case "utf8", "utf8mb3":
		return 3
	case "gbk":
		return 2
	}
	return 1
}

// Result types needed by TiDB's expression-index restrictions. Unknown function
// results stay unknown rather than being guessed to be an indexable scalar.
func tiDBIndexExpressionType(expr nodes.ExprNode, tbl *Table) *ResolvedType {
	switch n := expr.(type) {
	case *nodes.ParenExpr:
		return tiDBIndexExpressionType(n.Expr, tbl)
	case *nodes.CastExpr:
		if n.TypeName != nil {
			return dataTypeToResolvedType(n.TypeName)
		}
	case *nodes.ColumnRef:
		if col := tbl.GetColumn(n.Column); col != nil {
			switch strings.ToLower(col.DataType) {
			case "tinytext", "text", "mediumtext", "longtext", "tinyblob", "blob", "mediumblob", "longblob":
				length := 65535
				if strings.HasPrefix(col.DataType, "tiny") {
					length = 255
				} else if strings.HasPrefix(col.DataType, "medium") {
					length = 16777215
				} else if strings.HasPrefix(col.DataType, "long") {
					length = 4294967295
				}
				return &ResolvedType{BaseType: BaseTypeLongBlob, Length: length, Charset: col.Charset}
			case "char", "varchar":
				length := 0
				if _, after, ok := strings.Cut(col.ColumnType, "("); ok {
					if value, _, ok := strings.Cut(after, ")"); ok {
						if parsed, err := strconv.Atoi(value); err == nil {
							length = parsed
						}
					}
				}
				return &ResolvedType{BaseType: BaseTypeVarchar, Length: length, Charset: col.Charset}
			case "float", "double", "decimal", "json":
				return dataTypeToResolvedType(&nodes.DataType{Name: col.DataType})
			default:
				return &ResolvedType{BaseType: BaseTypeBigInt}
			}
		}
	case *nodes.BinaryExpr:
		switch n.Op {
		case nodes.BinOpJsonExtract:
			return &ResolvedType{BaseType: BaseTypeJSON}
		case nodes.BinOpJsonUnquote:
			return &ResolvedType{BaseType: BaseTypeLongText}
		case nodes.BinOpAdd, nodes.BinOpSub, nodes.BinOpMul, nodes.BinOpDiv, nodes.BinOpMod:
			base := BaseTypeBigInt
			if n.Op == nodes.BinOpDiv {
				base = BaseTypeDecimal
			}
			for _, arg := range []nodes.ExprNode{n.Left, n.Right} {
				if typ := tiDBIndexExpressionType(arg, tbl); typ != nil {
					switch typ.BaseType {
					case BaseTypeFloat, BaseTypeDouble, BaseTypeVarchar, BaseTypeLongText, BaseTypeJSON:
						return &ResolvedType{BaseType: BaseTypeDouble}
					case BaseTypeDecimal:
						base = BaseTypeDecimal
					}
				}
			}
			return &ResolvedType{BaseType: base}
		}
	case *nodes.FuncCallExpr:
		if strings.EqualFold(n.Name, "json_unquote") {
			return &ResolvedType{BaseType: BaseTypeLongText}
		}
		switch strings.ToLower(n.Name) {
		case "lower", "upper", "reverse":
			if len(n.Args) == 1 {
				if arg := tiDBIndexExpressionType(n.Args[0], tbl); arg != nil {
					return &ResolvedType{BaseType: BaseTypeVarchar, Length: arg.Length, Charset: arg.Charset}
				}
			}
		case "md5":
			return &ResolvedType{BaseType: BaseTypeVarchar, Length: 32}
		}
		return functionReturnType(n.Name, nil)
	case *nodes.ExtractExpr, *nodes.PositionExpr:
		return &ResolvedType{BaseType: BaseTypeBigInt}
	case *nodes.IntLit:
		return &ResolvedType{BaseType: BaseTypeBigInt}
	case *nodes.FloatLit:
		return &ResolvedType{BaseType: BaseTypeDecimal}
	case *nodes.UnaryExpr:
		return tiDBIndexExpressionType(n.Operand, tbl)
	}
	return nil
}

// https://github.com/pingcap/tidb/blob/v8.5.5/pkg/sessionctx/variable/varsutil.go
func tiDBGAIndexFunction(name string) bool {
	switch strings.ToLower(name) {
	case "lower", "upper", "md5", "reverse", "vitess_hash", "tidb_shard",
		"json_type", "json_extract", "json_unquote", "json_array", "json_object",
		"json_set", "json_insert", "json_replace", "json_remove", "json_contains",
		"json_contains_path", "json_valid", "json_array_append", "json_array_insert",
		"json_merge_patch", "json_merge_preserve", "json_pretty", "json_quote",
		"json_schema_valid", "json_search", "json_storage_size", "json_depth",
		"json_keys", "json_length":
		return true
	}
	return false
}
