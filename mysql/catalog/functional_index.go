package catalog

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

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

func synthesizeFunctionalIndexColumns(tbl *Table, idx *Index, ignoreSpace bool) error {
	return synthesizeFunctionalIndexColumnsForStmt(tbl, idx, nil, true, ignoreSpace)
}

func synthesizeFunctionalIndexColumnsForStmt(tbl *Table, idx *Index, leaves map[string]bool, checkRefs, ignoreSpace bool) error {
	for _, idxCol := range idx.Columns {
		if idxCol.Expr == "" {
			continue
		}
		if checkRefs {
			if err := validateFunctionalIndexColumnRefs(tbl, leaves, idxCol.ExprNode); err != nil {
				return err
			}
		}
		if idx.Fulltext || idx.Spatial {
			return validateFulltextSpatialFunctionalExpr(tbl, idx.Name, idxCol.ExprNode, idx.Spatial, ignoreSpace)
		}
		if err := validateFunctionalIndexExpr(tbl, idx.Name, idxCol.ExprNode, ignoreSpace); err != nil {
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

func stripTableQualifiersFromTableExprs(tbl *Table, oldDB, oldTable string, cloneAST bool, preExprs map[nodes.ExprNode]bool) {
	if tbl == nil || oldTable == "" {
		return
	}
	for _, col := range tbl.Columns {
		if col == nil {
			continue
		}
		if col.DefaultExprNode != nil {
			oldNode := col.DefaultExprNode
			rebound, sql, changed := rebindTableExprQualifiers(oldNode, oldDB, oldTable, cloneAST)
			if changed {
				col.DefaultExprNode = rebound
				if sql != "" {
					col.Default = &sql
				}
				if preExprs != nil && preExprs[oldNode] {
					delete(preExprs, oldNode)
					preExprs[rebound] = true
				}
			}
		}
		if col.Generated != nil && col.Generated.Expr != "" {
			if parsed, err := parseExpr(col.Generated.Expr); err == nil && parsed != nil {
				_, sql, changed := rebindTableExprQualifiers(parsed, oldDB, oldTable, false)
				if changed && sql != "" {
					col.Generated.Expr = sql
				}
			}
		}
	}
	for _, idx := range tbl.Indexes {
		if idx == nil {
			continue
		}
		for _, ic := range idx.Columns {
			if ic == nil || ic.ExprNode == nil {
				continue
			}
			rebound, sql, changed := rebindTableExprQualifiers(ic.ExprNode, oldDB, oldTable, cloneAST)
			if changed {
				ic.ExprNode = rebound
				if sql != "" {
					ic.Expr = sql
				}
			}
		}
	}
}

func rebindTableExprQualifiers(expr nodes.ExprNode, oldDB, oldTable string, cloneAST bool) (nodes.ExprNode, string, bool) {
	if expr == nil {
		return nil, "", false
	}
	target := expr
	if cloneAST {
		if cloned, err := parseExpr(nodeToSQL(expr)); err == nil && cloned != nil {
			target = cloned
		}
	}
	changed := false
	nodes.Inspect(target, func(n nodes.Node) bool {
		cr, ok := n.(*nodes.ColumnRef)
		if !ok || cr.Table == "" {
			return true
		}
		if (cr.Schema == "" || eqFoldStr(cr.Schema, oldDB)) && eqFoldStr(cr.Table, oldTable) {
			cr.Schema = ""
			cr.Table = ""
			changed = true
		}
		return true
	})
	if !changed {
		return target, "", false
	}
	return target, nodeToSQL(target), true
}

func recheckFunctionalIndexesOnFinalTable(tbl *Table, added []*Index, ignoreSpace bool) error {
	for _, idx := range added {
		if idx == nil {
			continue
		}
		for _, idxCol := range idx.Columns {
			if idxCol.Expr == "" {
				continue
			}
			if err := validateFunctionalIndexColumnRefs(tbl, nil, idxCol.ExprNode); err != nil {
				return err
			}
			if err := validateFunctionalIndexExpr(tbl, idx.Name, idxCol.ExprNode, ignoreSpace); err != nil {
				return err
			}
			if col := tbl.GetColumn(idxCol.Name); col != nil && col.Hidden == ColumnHiddenSystem {
				resolved := resolvedFunctionalIndexColumnSpec(tbl, inferFunctionalIndexExprType(tbl, idxCol.ExprNode))
				col.DataType = resolved.DataType
				col.ColumnType = resolved.ColumnType
				col.Charset = resolved.Charset
				col.Collation = resolved.Collation
			}
		}
	}
	return nil
}

func constraintHasFunctionalKeyPart(con *nodes.Constraint) bool {
	if con == nil {
		return false
	}
	for _, ic := range con.IndexColumns {
		if ic != nil && ic.Functional {
			return true
		}
	}
	return false
}

// validateNoColumnFunctionalParts refuses a functional key part that wraps a
// bare column reference: mysql:8.0.46 answers 3762 before checking index name
// collisions (1061) and before resolving whether the column exists (1054).
func validateNoColumnFunctionalParts(parts []*nodes.IndexColumn) error {
	for _, ic := range parts {
		if ic != nil && ic.Functional {
			if _, ok := peelFunctionalIndexParens(ic.Expr).(*nodes.ColumnRef); ok {
				return errFunctionalIndexOnField()
			}
		}
	}
	return nil
}

func validateFunctionalIndexColumnRefs(tbl *Table, leaves map[string]bool, expr nodes.ExprNode) error {
	var err error
	nodes.Inspect(expr, func(n nodes.Node) bool {
		if err != nil {
			return false
		}
		cr, ok := n.(*nodes.ColumnRef)
		if !ok || cr.Column == "" {
			return true
		}
		schemaOK := cr.Schema == "" || (tbl.Database != nil && strings.EqualFold(cr.Schema, tbl.Database.Name))
		tableOK := cr.Table == "" || strings.EqualFold(cr.Table, tbl.Name)
		if !schemaOK || !tableOK {
			err = errNoSuchColumn(cr.Column, "functional index")
			return false
		}
		if leaves != nil {
			// A functional index's hidden column is no name an expression
			// resolves, on ALTER as on CREATE: mysql:8.0.46 answers 1054.
			if !leaves[toLower(cr.Column)] {
				err = errNoSuchColumn(cr.Column, "functional index")
				return false
			}
			if col := tbl.GetColumn(cr.Column); col != nil && col.Hidden == ColumnHiddenSystem {
				err = errNoSuchColumn(cr.Column, "functional index")
				return false
			}
		} else {
			col := tbl.GetColumn(cr.Column)
			if col == nil || col.Hidden == ColumnHiddenSystem {
				err = errNoSuchColumn(cr.Column, "functional index")
				return false
			}
		}
		return true
	})
	return err
}

func validateDroppedOrRenamedFunctionalIndexDeps(tbl *Table, origTableName string, droppedOrRenamed []string, beforeIndexes map[*Index]bool) error {
	if len(droppedOrRenamed) == 0 {
		return nil
	}
	droppedSet := make(map[string]bool, len(droppedOrRenamed))
	for _, name := range droppedOrRenamed {
		droppedSet[toLower(name)] = true
	}
	for _, idx := range tbl.Indexes {
		if idx == nil || !beforeIndexes[idx] {
			continue
		}
		for _, ic := range idx.Columns {
			if ic == nil || ic.ExprNode == nil {
				continue
			}
			var depCol string
			nodes.Inspect(ic.ExprNode, func(n nodes.Node) bool {
				if depCol != "" {
					return false
				}
				if cr, ok := n.(*nodes.ColumnRef); ok && cr.Column != "" && cr.Schema == "" &&
					(cr.Table == "" || strings.EqualFold(cr.Table, origTableName) || strings.EqualFold(cr.Table, tbl.Name)) {
					if droppedSet[toLower(cr.Column)] {
						depCol = cr.Column
						return false
					}
				}
				return true
			})
			if depCol != "" {
				return errFunctionalIndexColumnDep(depCol)
			}
		}
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

func resolvedFunctionalIndexColumnSpec(tbl *Table, rt *ResolvedType) *Column {
	col := columnFromResolvedType(rt)
	if rt != nil && isResolvedStringType(rt) && col.Charset == "" && tbl != nil {
		col.Charset = tbl.Charset
		col.Collation = tbl.Collation
	}
	return col
}

func columnFromFunctionalIndexExpr(tbl *Table, hiddenName string, idxCol *IndexColumn) *Column {
	col := resolvedFunctionalIndexColumnSpec(tbl, inferFunctionalIndexExprType(tbl, idxCol.ExprNode))
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
	case *nodes.IntLit, *nodes.BoolLit:
		return &ResolvedType{BaseType: BaseTypeBigInt}
	case *nodes.FloatLit:
		// 1.5 is DECIMAL(2,1) and 1e0 a DOUBLE.
		if strings.ContainsAny(e.Value, "eE") {
			return &ResolvedType{BaseType: BaseTypeDouble}
		}
		digits := strings.TrimLeft(strings.ReplaceAll(strings.TrimLeft(e.Value, "+-"), ".", ""), "0")
		scale := 0
		if dot := strings.IndexByte(e.Value, '.'); dot >= 0 {
			scale = len(e.Value) - dot - 1
		}
		return &ResolvedType{BaseType: BaseTypeDecimal, Precision: max(len(digits), scale, 1), Scale: scale}
	case *nodes.StringLit:
		return &ResolvedType{BaseType: BaseTypeVarchar, Length: max(utf8.RuneCountInString(e.Value), 1), Charset: e.Charset}
	case *nodes.TemporalLit:
		switch strings.ToUpper(e.Type) {
		case "DATE":
			return &ResolvedType{BaseType: BaseTypeDate}
		case "TIME":
			return &ResolvedType{BaseType: BaseTypeTime}
		default:
			return &ResolvedType{BaseType: BaseTypeDateTime}
		}
	case *nodes.HexLit, *nodes.BitLit:
		// A hex or bit literal is a binary string of its bytes.
		val, isHex := "", false
		if h, ok := e.(*nodes.HexLit); ok {
			val, isHex = cleanHexDigits(h.Value), true
		} else {
			val = cleanBitDigits(e.(*nodes.BitLit).Value)
		}
		if b, ok := decodeLiteralRawBytes(val, isHex); ok {
			return &ResolvedType{BaseType: BaseTypeVarBinary, Length: max(len(b), 1)}
		}
	case *nodes.UnaryExpr:
		switch e.Op {
		case nodes.UnaryNot, nodes.UnaryBitNot:
			return &ResolvedType{BaseType: BaseTypeBigInt}
		case nodes.UnaryMinus, nodes.UnaryPlus:
			// Negation keeps a number's type and makes a string a double:
			// mysql:8.0.46 takes KEY ((-a)) and KEY ((-s)) on a COMPACT table.
			rt := inferFunctionalIndexExprType(tbl, e.Operand)
			switch {
			case rt == nil || !isResolvedNumericType(rt):
				return &ResolvedType{BaseType: BaseTypeDouble}
			case rt.BaseType == BaseTypeDecimal || rt.BaseType == BaseTypeFloat || rt.BaseType == BaseTypeDouble:
				cp := *rt
				return &cp
			}
			return &ResolvedType{BaseType: BaseTypeBigInt}
		}
	case *nodes.IsExpr, *nodes.InExpr, *nodes.BetweenExpr, *nodes.LikeExpr:
		// A predicate is an integer truth value: mysql:8.0.46 takes
		// KEY ((a IS NULL)), ((a IN (1,2))) and ((s LIKE 'x%')) on COMPACT.
		return &ResolvedType{BaseType: BaseTypeBigInt}
	case *nodes.BinaryExpr:
		switch e.Op {
		case nodes.BinOpEq, nodes.BinOpNe, nodes.BinOpLt, nodes.BinOpGt, nodes.BinOpLe, nodes.BinOpGe,
			nodes.BinOpNullSafeEq, nodes.BinOpAnd, nodes.BinOpOr, nodes.BinOpXor,
			nodes.BinOpRegexp, nodes.BinOpLikeEscape, nodes.BinOpSoundsLike:
			// A comparison or logical result is an integer truth value:
			// mysql:8.0.46 takes KEY ((a = 1)) on a COMPACT table.
			return &ResolvedType{BaseType: BaseTypeBigInt}
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
		if rt := inferFunctionalIndexFuncType(tbl, e); rt != nil {
			return rt
		}
		args := make([]AnalyzedExpr, 0, len(e.Args))
		for _, arg := range e.Args {
			args = append(args, &ConstExprQ{Type: inferFunctionalIndexExprType(tbl, arg)})
		}
		return functionReturnType(strings.ToLower(e.Name), args)
	case *nodes.CastExpr:
		if e.TypeName != nil {
			return dataTypeToResolvedType(e.TypeName)
		}
	}
	return &ResolvedType{BaseType: BaseTypeVarchar, Length: 255, widthUnknown: true}
}

// functionalTypeWidthKnown reports a functional key part type whose width the
// catalog inferred, rather than a placeholder.
func functionalTypeWidthKnown(rt *ResolvedType) bool {
	if rt == nil || rt.widthUnknown {
		return false
	}
	switch rt.BaseType {
	case BaseTypeChar, BaseTypeVarchar, BaseTypeBinary, BaseTypeVarBinary:
		return rt.Length > 0
	}
	return true
}

func extractNonNegIntLit(expr nodes.ExprNode) (int, bool) {
	lit, ok := peelFunctionalIndexParens(expr).(*nodes.IntLit)
	if !ok || lit.Value < 0 {
		return 0, false
	}
	return int(lit.Value), true
}

func exprCharsetRank(cs string) int {
	c := strings.ToLower(strings.TrimSpace(cs))
	switch {
	case c == "":
		return 0
	case c == "utf8mb4" || c == "utf16" || c == "utf16le" || c == "utf32":
		return 50
	case c == "gb18030":
		return 45
	case c == "utf8" || c == "utf8mb3" || c == "ucs2":
		return 40
	case charsetMaxBytes(c) == 3:
		return 30
	case charsetMaxBytes(c) == 2:
		return 20
	case c == "ascii":
		return 5
	default:
		return 10
	}
}

func aggregateExprCharset(a, b string) string {
	if exprCharsetRank(b) > exprCharsetRank(a) {
		return b
	}
	return a
}

func inferStringExprLenAndCharset(tbl *Table, expr nodes.ExprNode) (int, bool, string) {
	switch e := peelFunctionalIndexParens(expr).(type) {
	case *nodes.StringLit:
		return utf8.RuneCountInString(e.Value), true, ""
	case *nodes.IntLit:
		return len(strconv.FormatInt(e.Value, 10)), true, ""
	default:
		rt := inferFunctionalIndexExprType(tbl, e)
		if rt != nil && rt.Length > 0 {
			return rt.Length, true, rt.Charset
		}
		return 0, false, ""
	}
}

func inferBoundedStringExprType(tbl *Table, args []nodes.ExprNode, limitIdx int) *ResolvedType {
	if len(args) < 2 {
		return nil
	}
	var maxLen int
	var charset string
	base := BaseTypeVarchar
	if arg0 := inferFunctionalIndexExprType(tbl, args[0]); arg0 != nil {
		maxLen = arg0.Length
		charset = arg0.Charset
		if isResolvedBinaryType(arg0) {
			base = BaseTypeVarBinary
		}
	}
	if len(args) > limitIdx {
		if k, ok := extractNonNegIntLit(args[limitIdx]); ok && k > 0 {
			if maxLen == 0 || k < maxLen {
				maxLen = k
			}
		}
	}
	if maxLen > 0 {
		return &ResolvedType{BaseType: base, Length: maxLen, Charset: charset}
	}
	return nil
}

func inferFunctionalIndexFuncType(tbl *Table, e *nodes.FuncCallExpr) *ResolvedType {
	name := strings.ToLower(e.Name)
	switch name {
	case "lower", "lcase", "upper", "ucase", "trim", "ltrim", "rtrim", "reverse":
		if len(e.Args) > 0 {
			rt := inferFunctionalIndexExprType(tbl, e.Args[0])
			if rt != nil && (isResolvedStringType(rt) || isResolvedBinaryType(rt)) {
				cp := *rt
				return &cp
			}
		}
		return &ResolvedType{BaseType: BaseTypeVarchar}
	case "concat", "concat_ws":
		if len(e.Args) == 0 || (name == "concat_ws" && len(e.Args) < 2) {
			return nil
		}
		sepLen := 0
		sepOK := true
		charset := ""
		start := 0
		if name == "concat_ws" {
			sepLen, sepOK, charset = inferStringExprLenAndCharset(tbl, e.Args[0])
			start = 1
		}
		if !sepOK {
			return nil
		}
		base := BaseTypeVarchar
		for _, arg := range e.Args {
			at := inferFunctionalIndexExprType(tbl, arg)
			if isFunctionalIndexLOBType(at) {
				return at
			}
			if isResolvedBinaryType(at) {
				base = BaseTypeVarBinary
			}
		}
		// A binary CONCAT_WS counts a character separator at utf8mb4's 4
		// bytes per character too: mysql:8.0.46 takes
		// CONCAT_WS('x', varbinary_254, varbinary_254) and answers 3757 at 255.
		if name == "concat_ws" && base == BaseTypeVarBinary && !isResolvedBinaryType(inferFunctionalIndexExprType(tbl, e.Args[0])) {
			w := "utf8mb4"
			if charset != "" {
				w = charset
			}
			sepLen *= max(charsetMaxBytes(w), 1)
		}
		totalLen := 0
		for i := start; i < len(e.Args); i++ {
			argLen, ok, cs := inferStringExprLenAndCharset(tbl, e.Args[i])
			if !ok {
				return nil
			}
			// A binary CONCAT counts a character argument at its charset's
			// widest character, a literal at utf8mb4's: mysql:8.0.46 refuses
			// KEY ((CONCAT(tinyblob, tinyblob, 'x'))), 255 + 255 + 4 bytes.
			if base == BaseTypeVarBinary && !isResolvedBinaryType(inferFunctionalIndexExprType(tbl, e.Args[i])) {
				w := "utf8mb4"
				if cs != "" {
					w = cs
				}
				argLen *= max(charsetMaxBytes(w), 1)
			}
			totalLen += argLen
			charset = aggregateExprCharset(charset, cs)
			if name == "concat_ws" && i > start {
				totalLen += sepLen
			}
		}
		if totalLen < 1 {
			totalLen = 1
		}
		return &ResolvedType{BaseType: base, Length: totalLen, Charset: charset}
	case "substring", "substr", "mid":
		return inferBoundedStringExprType(tbl, e.Args, 2)
	case "left", "right":
		return inferBoundedStringExprType(tbl, e.Args, 1)
	case "lpad", "rpad":
		if len(e.Args) < 2 {
			return nil
		}
		if k, ok := extractNonNegIntLit(e.Args[1]); ok && k > 0 {
			var charset string
			base := BaseTypeVarchar
			if arg0 := inferFunctionalIndexExprType(tbl, e.Args[0]); arg0 != nil {
				charset = arg0.Charset
				// A binary string pads to a binary string: mysql:8.0.46 takes
				// RPAD(varbinary_1, 512, X'00') on a COMPACT table.
				if isResolvedBinaryType(arg0) {
					base = BaseTypeVarBinary
				}
			}
			return &ResolvedType{BaseType: base, Length: k, Charset: charset}
		}
	case "ascii", "ord", "bit_length", "char_length", "character_length", "length", "octet_length",
		"crc32", "bit_count", "sign", "instr", "locate", "position", "field", "find_in_set",
		"strcmp", "interval", "inet_aton", "is_ipv4", "is_ipv6", "is_ipv4_compat", "is_ipv4_mapped",
		"is_uuid", "json_length", "json_depth", "json_valid", "json_contains", "json_contains_path",
		"json_overlaps", "regexp_instr", "regexp_like", "year", "month", "day", "dayofmonth",
		"dayofweek", "dayofyear", "weekday", "week", "weekofyear", "yearweek", "quarter",
		"hour", "minute", "second", "microsecond", "to_days", "to_seconds", "datediff",
		"period_add", "period_diff", "timestampdiff", "time_to_sec", "isnull":
		// An integer result: mysql:8.0.46 holds KEY ((ASCII(c1)), ..., (ASCII(c13)))
		// to 13 small integer parts, not to placeholder strings.
		return &ResolvedType{BaseType: BaseTypeBigInt}
	case "space":
		// SPACE(n) is n characters, and SPACE of anything but a constant
		// is as long as a LOB: mysql:8.0.46 takes KEY ((SPACE(512))) and
		// answers 3757 for SPACE(513) and SPACE(a).
		if len(e.Args) != 1 {
			return nil
		}
		k, ok := extractNonNegIntLit(e.Args[0])
		if !ok || k > 1<<30 {
			return &ResolvedType{BaseType: BaseTypeLongText}
		}
		return &ResolvedType{BaseType: BaseTypeVarchar, Length: max(k, 1)}
	case "replace":
		// REPLACE's width, as mysql:8.0.46 measures it, is the subject's
		// times the replacement's length less one, at least the subject's:
		// REPLACE(varchar_100, 'x', 'yyy') is varchar(200) and with 'yy' or
		// 'y' varchar(100), whatever the search string. A binary subject counts
		// a character replacement at utf8mb4's 4 bytes per character.
		if len(e.Args) != 3 {
			return nil
		}
		arg0 := inferFunctionalIndexExprType(tbl, e.Args[0])
		if arg0 == nil || !functionalTypeWidthKnown(arg0) || isFunctionalIndexLOBType(arg0) {
			return arg0
		}
		replLen, ok, _ := inferStringExprLenAndCharset(tbl, e.Args[2])
		if !ok {
			if _, isStr := peelFunctionalIndexParens(e.Args[2]).(*nodes.StringLit); !isStr {
				return nil
			}
			replLen = 0
		}
		base := BaseTypeVarchar
		if isResolvedBinaryType(arg0) {
			base = BaseTypeVarBinary
			if !isResolvedBinaryType(inferFunctionalIndexExprType(tbl, e.Args[2])) {
				replLen *= 4
			}
		}
		n := arg0.Length
		if replLen > 2 {
			n = arg0.Length * (replLen - 1)
		}
		return &ResolvedType{BaseType: base, Length: n, Charset: arg0.Charset}
	case "repeat":
		if len(e.Args) < 2 {
			return nil
		}
		argLen, argOK, charset := inferStringExprLenAndCharset(tbl, e.Args[0])
		if k, ok := extractNonNegIntLit(e.Args[1]); ok && argOK && k > 0 {
			// REPEAT keeps a binary argument binary, and its length saturates
			// past any limit an index checks: mysql:8.0.46 takes
			// CONCAT(varbinary_1, REPEAT(X'00', 128)) and refuses
			// REPEAT(varchar_2, 9223372036854775807) with 3757.
			base := BaseTypeVarchar
			if isResolvedBinaryType(inferFunctionalIndexExprType(tbl, e.Args[0])) {
				base = BaseTypeVarBinary
			}
			const lengthCap = 1 << 30
			n := lengthCap
			if argLen <= lengthCap/k {
				n = max(argLen*k, 1)
			}
			return &ResolvedType{BaseType: base, Length: n, Charset: charset}
		}
	}
	return nil
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
	case "binary", "varbinary":
		rt.BaseType = BaseTypeVarBinary
		if n, ok := fixedLengthStringColumnLimit(col); ok {
			rt.Length = n
		}
	// A TINYTEXT or TINYBLOB holds 255 bytes, which no function result over
	// it turns into a LOB: mysql:8.0.46 takes KEY ((REVERSE(tinytext_col))).
	case "tinyblob":
		rt.BaseType = BaseTypeVarBinary
		rt.Length = 255
	case "tinytext":
		rt.BaseType = BaseTypeVarchar
		rt.Length = 255 / max(charsetMaxBytes(col.Charset), 1)
	case "text", "mediumtext", "longtext":
		rt.BaseType = BaseTypeText
	case "blob", "mediumblob", "longblob":
		rt.BaseType = BaseTypeBlob
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

func validateFunctionalIndexExpr(tbl *Table, indexName string, expr nodes.ExprNode, ignoreSpace bool) error {
	if _, ok := peelFunctionalIndexParens(expr).(*nodes.ColumnRef); ok {
		return errFunctionalIndexOnField()
	}
	if err := functionalIndexExprRefusal(indexName, expr, ignoreSpace); err != nil {
		return err
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
			Code:     3757,
			SQLState: "HY000",
			Message:  "Cannot create a functional index on an expression that returns a BLOB or TEXT. Please consider using CAST.",
		}
	}
	return nil
}

// validateFulltextSpatialFunctionalExpr refuses a functional key part of a
// FULLTEXT (3759) or SPATIAL (3760) index. mysql:8.0.46 judges the expression
// first, as for any functional key part, and its type only for a BLOB or TEXT
// result: FULLTEXT ((NOW())) is 3758, FULLTEXT ((CONCAT(a, REPEAT(a, 10000))))
// is 3757, and FULLTEXT ((CAST(a AS JSON))) is 3759, not 3753.
func validateFulltextSpatialFunctionalExpr(tbl *Table, indexName string, expr nodes.ExprNode, spatial, ignoreSpace bool) error {
	if _, ok := peelFunctionalIndexParens(expr).(*nodes.ColumnRef); ok {
		return errFunctionalIndexOnField()
	}
	if err := functionalIndexExprRefusal(indexName, expr, ignoreSpace); err != nil {
		return err
	}
	if rt := inferFunctionalIndexExprType(tbl, expr); isFunctionalIndexLOBType(rt) && !isFunctionalIndexJSONOrGeometryType(rt) {
		return &Error{
			Code:     3757,
			SQLState: "HY000",
			Message:  "Cannot create a functional index on an expression that returns a BLOB or TEXT. Please consider using CAST.",
		}
	}
	if spatial {
		return errSpatialFunctionalIndex()
	}
	return errFulltextFunctionalIndex()
}

// validateFulltextSpatialFunctionalParts judges the functional key parts of a
// FULLTEXT or SPATIAL key a table definition declares, after its plain key parts.
func validateFulltextSpatialFunctionalParts(tbl *Table, leaves map[string]bool, checkRefs bool, indexName string, idxCols []*IndexColumn, spatial, ignoreSpace bool) error {
	if indexName == "" {
		indexName = "functional_index"
	}
	for _, ic := range idxCols {
		if ic == nil || ic.Expr == "" {
			continue
		}
		if checkRefs {
			if err := validateFunctionalIndexColumnRefs(tbl, leaves, ic.ExprNode); err != nil {
				return err
			}
		}
		return validateFulltextSpatialFunctionalExpr(tbl, indexName, ic.ExprNode, spatial, ignoreSpace)
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

// functionalIndexExprRefusal answers what mysql:8.0.46 answers for the first
// node of a functional key part it refuses, or nil. It judges in post-order, a
// function's arguments before the function: KEY ((@x + NOW())) and
// KEY ((SLEEP(@x))) are 3772, KEY ((NOW() + @x)) is 3758. Aggregates (1111)
// and window functions (3593) are judged only after the whole expression
// passes: KEY ((COUNT(a) + NOW())) is 3758, KEY ((COUNT(a))) is 1111.
func functionalIndexExprRefusal(indexName string, expr nodes.ExprNode, ignoreSpace bool) error {
	v := &functionalIndexJudge{ignoreSpace: ignoreSpace}
	nodes.Walk(v, expr)
	switch v.code {
	case ErrDefaultValExprRefVariable:
		return errDefaultValExprVariable(indexName)
	case 3758:
		return &Error{
			Code:     3758,
			SQLState: "HY000",
			Message:  fmt.Sprintf("Expression of functional index '%s' contains a disallowed function.", indexName),
		}
	}
	var err error
	nodes.Inspect(expr, func(n nodes.Node) bool {
		fc, ok := n.(*nodes.FuncCallExpr)
		if err != nil || !ok {
			return err == nil
		}
		if fc.Over != nil {
			err = &Error{Code: 3593, SQLState: "HY000",
				Message: fmt.Sprintf("You cannot use the window function '%s' in this context.'", strings.ToLower(fc.Name))}
		} else if functionalIndexAggregates[strings.ToLower(fc.Name)] {
			err = &Error{Code: 1111, SQLState: "HY000", Message: "Invalid use of group function"}
		}
		return err == nil
	})
	return err
}

type functionalIndexJudge struct {
	stack       []nodes.Node
	code        int
	ignoreSpace bool
}

func (v *functionalIndexJudge) Visit(n nodes.Node) nodes.Visitor {
	if v.code != 0 {
		return nil
	}
	if n == nil {
		top := v.stack[len(v.stack)-1]
		v.stack = v.stack[:len(v.stack)-1]
		switch e := top.(type) {
		case *nodes.VariableRef:
			v.code = ErrDefaultValExprRefVariable
		case *nodes.DefaultExpr, *nodes.MatchExpr:
			v.code = 3758
		case *nodes.FuncCallExpr:
			if isFunctionalIndexDisallowedFunc(e, v.ignoreSpace) {
				v.code = 3758
			}
		}
		return nil
	}
	switch e := n.(type) {
	case *nodes.SubqueryExpr, *nodes.ExistsExpr:
		v.code = 3758
		return nil
	case *nodes.InExpr:
		if e.Select != nil {
			v.code = 3758
			return nil
		}
	}
	v.stack = append(v.stack, n)
	return v
}

// isFunctionalIndexDisallowedFunc reports a call a functional index refuses
// with 3758: a stored function or UDF (a qualified or unknown name, or a spaced
// gated keyword such as `position (a)`), a built-in DEFAULT (expr) refuses too
// (LAST_INSERT_ID, VERSION, VALUES, ...), and the built-ins DEFAULT takes but a
// generated column does not. Aggregates and window functions are judged apart.
func isFunctionalIndexDisallowedFunc(e *nodes.FuncCallExpr, ignoreSpace bool) bool {
	// A spaced call names a stored function unless IGNORE_SPACE makes it the
	// built-in: mysql:8.0.46 takes KEY ((SUBSTRING (a,1,1))) under
	// IGNORE_SPACE or ANSI and answers 3758 without it.
	if e.Schema != "" || (e.SpacedParen && !ignoreSpace) {
		return true
	}
	name := strings.ToLower(e.Name)
	if e.Over != nil || functionalIndexAggregates[name] {
		return false
	}
	name = canonicalDefaultBuiltinName(name)
	return functionalIndexDisallowedBuiltins[name] || !allowedDefaultExprBuiltinFuncs[name]
}

// functionalIndexAggregates are the group functions a functional key part
// answers 1111 for. ANY_VALUE is not one: KEY ((ANY_VALUE(a))) is taken.
var functionalIndexAggregates = map[string]bool{
	"count": true, "sum": true, "avg": true, "min": true, "max": true,
	"group_concat": true, "json_arrayagg": true, "json_objectagg": true,
	"bit_and": true, "bit_or": true, "bit_xor": true,
	"std": true, "stddev": true, "stddev_pop": true, "stddev_samp": true,
	"var_pop": true, "var_samp": true, "variance": true, "st_collect": true,
}

// functionalIndexDisallowedBuiltins are the built-ins allowedDefaultExprBuiltinFuncs
// takes that mysql:8.0.46 refuses in a functional key part with 3758: every
// name of that list was probed as KEY ((fn(...))) and these answered 3758.
var functionalIndexDisallowedBuiltins = map[string]bool{
	"connection_id": true, "create_digest": true, "curdate": true, "current_date": true,
	"current_role": true, "current_time": true, "current_timestamp": true, "current_user": true,
	"curtime": true, "database": true, "geomcollfromtext": true, "geomcollfromwkb": true,
	"geometrycollectionfromtext": true, "geometrycollectionfromwkb": true, "geomfromtext": true,
	"geomfromwkb": true, "linefromtext": true, "linefromwkb": true, "localtime": true,
	"localtimestamp": true, "mlinefromtext": true, "mlinefromwkb": true, "mpointfromtext": true,
	"mpointfromwkb": true, "mpolyfromtext": true, "mpolyfromwkb": true, "now": true,
	"pointfromtext": true, "pointfromwkb": true, "polyfromtext": true, "polyfromwkb": true,
	"rand": true, "roles_graphml": true, "schema": true, "session_user": true,
	"statement_digest": true, "statement_digest_text": true, "sysdate": true, "system_user": true,
	"unix_timestamp": true, "user": true, "utc_date": true, "utc_time": true,
	"utc_timestamp": true, "uuid": true, "uuid_short": true,
}

func isFunctionalIndexJSONOrGeometryType(rt *ResolvedType) bool {
	switch rt.BaseType {
	case BaseTypeJSON, BaseTypeGeometry, BaseTypePoint, BaseTypeLineString, BaseTypePolygon,
		BaseTypeMultiPoint, BaseTypeMultiLineString, BaseTypeMultiPolygon, BaseTypeGeometryCollection:
		return true
	}
	return false
}

// isResolvedNumericType reports an integer, fixed-point or floating-point type.
func isResolvedNumericType(rt *ResolvedType) bool {
	switch rt.BaseType {
	case BaseTypeTinyInt, BaseTypeSmallInt, BaseTypeMediumInt, BaseTypeInt, BaseTypeBigInt,
		BaseTypeDecimal, BaseTypeFloat, BaseTypeDouble:
		return true
	}
	return false
}

// isResolvedBinaryType reports a binary string type a function result keeps.
func isResolvedBinaryType(rt *ResolvedType) bool {
	if rt == nil {
		return false
	}
	switch rt.BaseType {
	case BaseTypeBinary, BaseTypeVarBinary, BaseTypeTinyBlob, BaseTypeBlob, BaseTypeMediumBlob, BaseTypeLongBlob:
		return true
	}
	return false
}

// isFunctionalIndexLOBType reports a result a functional key part refuses with
// 3757: a TEXT, BLOB, JSON or geometry type, or a string longer than 512
// characters (bytes for a binary string), which MySQL makes a TEXT or BLOB:
// mysql:8.0.46 takes KEY ((LEFT(v, 512))) and refuses KEY ((LEFT(v, 513))),
// KEY ((CAST(i AS CHAR(513)))) and KEY ((REVERSE(varbinary_4000_col))).
func isFunctionalIndexLOBType(rt *ResolvedType) bool {
	if rt == nil {
		return false
	}
	switch rt.BaseType {
	case BaseTypeChar, BaseTypeVarchar, BaseTypeBinary, BaseTypeVarBinary:
		return rt.Length > 512
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
