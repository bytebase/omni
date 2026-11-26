package review

import (
	"strconv"
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/mysql/parser"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
)

// tablePartitionInfo tracks what the walk knows about a table's partitioning.
type tablePartitionInfo struct {
	partitioned bool
	hasType     bool
	partType    ast.PartitionType
	names       map[string]bool // lower-case partition names; nil if unknown
	count       int             // number of partitions; 0 if unknown
	cols        map[string]bool // lower-case columns referenced by partition function; nil if unknown
	// single is set when each partition's bound is one value: RANGE or
	// LIST of an expression, or of one column.
	single bool
	// subparts is set when the table has subpartitions, whose names
	// TRUNCATE and EXCHANGE PARTITION take too; names holds none of them.
	subparts bool
	// bounds holds each partition's name and bound in order, as
	// INFORMATION_SCHEMA.PARTITIONS describes them; nil if unknown.
	bounds []partBound
}

// partBound is one partition's lower-case name and its bound.
type partBound struct {
	name, value string
}

func (p *tablePartitionInfo) clone() *tablePartitionInfo {
	if p == nil {
		return nil
	}
	cp := *p
	cp.bounds = append([]partBound(nil), p.bounds...)
	if p.names != nil {
		cp.names = make(map[string]bool, len(p.names))
		for k, v := range p.names {
			cp.names[k] = v
		}
	}
	if p.cols != nil {
		cp.cols = make(map[string]bool, len(p.cols))
		for k, v := range p.cols {
			cp.cols[k] = v
		}
	}
	return &cp
}

// indexMeta records what the walk knows about one index on a table.
type indexMeta struct {
	unique      bool
	primary     bool
	fulltext    bool
	spatial     bool
	cols        []string // leading non-functional column keys (lower-case)
	colPrefixes []bool   // parallel to cols: true if key part has a prefix length > 0
}

func (m *indexMeta) clone() *indexMeta {
	if m == nil {
		return nil
	}
	cp := *m
	if m.cols != nil {
		cp.cols = append([]string(nil), m.cols...)
	}
	if m.colPrefixes != nil {
		cp.colPrefixes = append([]bool(nil), m.colPrefixes...)
	}
	return &cp
}

func (m *indexMeta) anyPrefix() bool {
	if m == nil {
		return false
	}
	for _, p := range m.colPrefixes {
		if p {
			return true
		}
	}
	return false
}

func indexMetaFromParts(unique, primary, fulltext, spatial bool, parts []keyPartSpec) *indexMeta {
	im := &indexMeta{
		unique:   unique,
		primary:  primary,
		fulltext: fulltext,
		spatial:  spatial,
	}
	for _, p := range parts {
		if p.functional {
			break
		}
		im.cols = append(im.cols, strings.ToLower(p.name))
		im.colPrefixes = append(im.colPrefixes, p.length > 0)
	}
	return im
}

func partitionProtoType(t metadata.TablePartitionMetadata_Type) (ast.PartitionType, bool) {
	switch t {
	case metadata.TablePartitionMetadata_RANGE, metadata.TablePartitionMetadata_RANGE_COLUMNS:
		return ast.PartitionRange, true
	case metadata.TablePartitionMetadata_LIST, metadata.TablePartitionMetadata_LIST_COLUMNS:
		return ast.PartitionList, true
	case metadata.TablePartitionMetadata_HASH, metadata.TablePartitionMetadata_LINEAR_HASH:
		return ast.PartitionHash, true
	case metadata.TablePartitionMetadata_KEY, metadata.TablePartitionMetadata_LINEAR_KEY:
		return ast.PartitionKey, true
	}
	return 0, false
}

// extractPartitionExprCols extracts column names referenced in a synced
// partition expression string by parsing it as a MySQL expression list.
func extractPartitionExprCols(expr string) []string {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return nil
	}
	list, err := parser.Parse("SELECT " + expr)
	if err != nil || list == nil || len(list.Items) == 0 {
		return nil
	}
	sel, ok := list.Items[0].(*ast.SelectStmt)
	if !ok {
		return nil
	}
	var out []string
	seen := make(map[string]bool)
	for _, item := range sel.TargetList {
		ast.Inspect(item, func(n ast.Node) bool {
			if cr, ok := n.(*ast.ColumnRef); ok && !cr.Star && cr.Column != "" {
				k := strings.ToLower(cr.Column)
				if !seen[k] {
					seen[k] = true
					out = append(out, k)
				}
			}
			return true
		})
	}
	return out
}

func partitionInfoFromClause(pc *ast.PartitionClause) *tablePartitionInfo {
	if pc == nil {
		return &tablePartitionInfo{partitioned: false}
	}
	info := &tablePartitionInfo{
		partitioned: true,
		hasType:     true,
		partType:    pc.Type,
		cols:        make(map[string]bool),
		single:      len(pc.Columns) <= 1,
		subparts:    partitionHasSubpartitions(pc),
	}
	if len(pc.Partitions) > 0 {
		info.names = make(map[string]bool, len(pc.Partitions))
		for _, pd := range pc.Partitions {
			if pd != nil && pd.Name != "" {
				info.names[strings.ToLower(pd.Name)] = true
			}
		}
		info.count = len(pc.Partitions)
	} else if pc.NumParts > 0 {
		info.count = pc.NumParts
	}
	for _, c := range pc.Columns {
		if c != "" {
			info.cols[strings.ToLower(c)] = true
		}
	}
	for _, c := range pc.SubPartColumns {
		if c != "" {
			info.cols[strings.ToLower(c)] = true
		}
	}
	for _, c := range exprColumnNames(pc.Expr) {
		info.cols[strings.ToLower(c)] = true
	}
	for _, c := range exprColumnNames(pc.SubPartExpr) {
		info.cols[strings.ToLower(c)] = true
	}
	return info
}

// normalizeBaseType returns the lower-case base MySQL data type name from a
// metadata column type string (e.g. "int(11) unsigned" -> "int") or AST type.
func normalizeBaseType(raw string) string {
	s := strings.TrimSpace(raw)
	if s == "" {
		return ""
	}
	if i := strings.IndexAny(s, " ("); i >= 0 {
		s = s[:i]
	}
	s = strings.ToLower(s)
	switch s {
	case "integer", "int4":
		return "int"
	case "int1", "bool", "boolean":
		return "tinyint"
	case "int2":
		return "smallint"
	case "int3", "middleint":
		return "mediumint"
	case "int8":
		return "bigint"
	case "float4":
		return "float"
	case "float8", "real":
		return "double"
	case "dec", "numeric", "fixed":
		return "decimal"
	case "character", "nchar":
		return "char"
	case "nvarchar":
		return "varchar"
	case "geomcollection":
		return "geometrycollection"
	}
	return s
}

func astBaseType(dt *ast.DataType) string {
	if dt == nil {
		return ""
	}
	return normalizeBaseType(dt.Name)
}

func isIntegerBaseType(dt string) bool {
	switch dt {
	case "tinyint", "smallint", "mediumint", "int", "bigint":
		return true
	}
	return false
}

func isFulltextBaseType(dt string) bool {
	switch dt {
	case "char", "varchar", "tinytext", "text", "mediumtext", "longtext":
		return true
	}
	return false
}

func isGeometryBaseType(dt string) bool {
	switch dt {
	case "geometry", "geometrycollection", "point", "linestring", "polygon",
		"multipoint", "multilinestring", "multipolygon":
		return true
	}
	return false
}

func isPrefixableBaseType(dt string) bool {
	switch dt {
	case "char", "varchar", "binary", "varbinary",
		"tinytext", "text", "mediumtext", "longtext",
		"tinyblob", "blob", "mediumblob", "longblob":
		return true
	}
	return false
}

// spatialOnePart is the finding for a SPATIAL key of more than one part.
var spatialOnePart = "A " + review.Code("SPATIAL") + " index takes one key part."

// usingType is an index's declared type: USING (or TYPE) written after
// the key parts overrides the one written before them.
func usingType(indexType string, opts []*ast.IndexOption) string {
	for _, opt := range opts {
		if strings.EqualFold(opt.Name, "USING") {
			if s, ok := opt.Value.(*ast.StringLit); ok {
				indexType = s.Value
			}
		}
	}
	return indexType
}

// isSpatialConstraint reports whether a table constraint defines a SPATIAL
// key: SPATIAL, or USING RTREE before or after its key parts.
func isSpatialConstraint(con *ast.Constraint) bool {
	return con.Type == ast.ConstrSpatialIndex || isSpatialIndexType(usingType(con.IndexType, con.IndexOptions))
}

func isSpatialIndexType(indexType string) bool {
	switch strings.ToUpper(strings.TrimSpace(indexType)) {
	case "SPATIAL", "RTREE":
		return true
	}
	return false
}

func isInnoDBEngine(engine string) bool {
	return engine == "" || strings.EqualFold(engine, "InnoDB")
}

// exprHasBareColumnRef reports whether expr references a column in the current
// scope (outside nested subqueries and schema-qualified calls).
func exprHasBareColumnRef(expr ast.ExprNode) (string, bool) {
	if expr == nil {
		return "", false
	}
	found := ""
	ast.Inspect(expr, func(n ast.Node) bool {
		if found != "" {
			return false
		}
		switch e := n.(type) {
		case *ast.ColumnRef:
			if !e.Star && e.Column != "" {
				found = e.Column
			}
			return false
		case *ast.SubqueryExpr, *ast.ExistsExpr, *ast.SelectStmt, *ast.TableStmt, *ast.ValuesStmt:
			return false
		case *ast.FuncCallExpr:
			if e.Schema != "" {
				return false
			}
		}
		return true
	})
	return found, found != ""
}

func exprColumnNames(expr ast.ExprNode) []string {
	if expr == nil {
		return nil
	}
	var names []string
	ast.Inspect(expr, func(n ast.Node) bool {
		switch e := n.(type) {
		case *ast.ColumnRef:
			if !e.Star && e.Column != "" {
				names = append(names, e.Column)
			}
			return false
		case *ast.SubqueryExpr, *ast.ExistsExpr, *ast.SelectStmt, *ast.TableStmt, *ast.ValuesStmt:
			return false
		}
		return true
	})
	return names
}

func exprReferencesColumn(expr ast.ExprNode, colName string) bool {
	if expr == nil || colName == "" {
		return false
	}
	found := false
	ast.Inspect(expr, func(n ast.Node) bool {
		if found {
			return false
		}
		switch e := n.(type) {
		case *ast.ColumnRef:
			if !e.Star && e.Table == "" && e.Schema == "" && strings.EqualFold(e.Column, colName) {
				found = true
			}
			return false
		case *ast.SubqueryExpr, *ast.ExistsExpr, *ast.SelectStmt, *ast.TableStmt, *ast.ValuesStmt:
			return false
		}
		return true
	})
	return found
}

// selectFromlessColumn returns the first column referenced in the target list
// of a SELECT branch that has no FROM clause (or FROM DUAL).
func selectFromlessColumn(sel *ast.SelectStmt) (string, bool) {
	if sel == nil {
		return "", false
	}
	if sel.SetOp != ast.SetOpNone {
		if col, ok := selectFromlessColumn(sel.Left); ok {
			return col, true
		}
		return selectFromlessColumn(sel.Right)
	}
	if len(sel.From) == 0 {
		for _, target := range sel.TargetList {
			expr := target
			if rt, ok := target.(*ast.ResTarget); ok {
				expr = rt.Val
			}
			if col, ok := exprHasBareColumnRef(expr); ok {
				return col, true
			}
		}
	}
	return "", false
}

// isBareColumnRef reports whether expr is a bare column reference (optionally
// parenthesized), which MySQL 8.0 forbids as a functional index key part (3762).
func isBareColumnRef(expr ast.ExprNode) bool {
	for {
		p, ok := expr.(*ast.ParenExpr)
		if !ok {
			break
		}
		expr = p.Expr
	}
	cr, ok := expr.(*ast.ColumnRef)
	return ok && !cr.Star && cr.Column != ""
}

// checkExprIsBoolean reports whether a CHECK constraint's top-level expression
// is accepted as boolean by MySQL 8.0 (false -> error 3812). Unrecognized forms
// return true so no valid expression is falsely refused.
func checkExprIsBoolean(expr ast.ExprNode) bool {
	switch e := expr.(type) {
	case *ast.ParenExpr:
		return checkExprIsBoolean(e.Expr)
	case *ast.BinaryExpr:
		switch e.Op {
		case ast.BinOpEq, ast.BinOpNe, ast.BinOpLt, ast.BinOpGt, ast.BinOpLe, ast.BinOpGe,
			ast.BinOpAnd, ast.BinOpOr, ast.BinOpXor, ast.BinOpRegexp, ast.BinOpLikeEscape,
			ast.BinOpNullSafeEq, ast.BinOpSoundsLike:
			return true
		}
		return false
	case *ast.UnaryExpr:
		return e.Op == ast.UnaryNot
	case *ast.IsExpr, *ast.InExpr, *ast.BetweenExpr, *ast.LikeExpr, *ast.ExistsExpr,
		*ast.MemberOfExpr, *ast.MatchExpr:
		return true
	case *ast.BoolLit:
		return true
	case *ast.ColumnRef, *ast.IntLit, *ast.FloatLit, *ast.StringLit, *ast.NullLit,
		*ast.BitLit, *ast.CastExpr, *ast.ConvertExpr, *ast.CollateExpr,
		*ast.WeightStringExpr, *ast.IntervalExpr:
		return false
	}
	return true
}

// partitionBoundColumn returns the first column referenced in a partition
// definition's VALUES LESS THAN / VALUES IN bound (error 1054).
func partitionBoundColumn(v ast.Node) (string, bool) {
	found := ""
	var visit func(ast.Node)
	visit = func(n ast.Node) {
		if found != "" || n == nil {
			return
		}
		switch node := n.(type) {
		case *ast.List:
			for _, item := range node.Items {
				visit(item)
			}
		case ast.ExprNode:
			if col, ok := exprHasBareColumnRef(node); ok {
				found = col
			}
		}
	}
	visit(v)
	return found, found != ""
}

func partitionHasSubpartitions(pc *ast.PartitionClause) bool {
	return pc.SubPartType == ast.PartitionHash || pc.SubPartType == ast.PartitionKey ||
		pc.SubPartExpr != nil || len(pc.SubPartColumns) > 0 || pc.NumSubParts > 0
}

// createSpatialRefSys checks CREATE [OR REPLACE] SPATIAL REFERENCE SYSTEM
// required attributes (3708, 3714, 3710, 3517).
func (b *builder) createSpatialRefSys(v *ast.CreateSpatialRefSysStmt) bool {
	rng := b.anchor()
	switch {
	case v.SRID == 0:
		b.fail(rng, nil, review.UnsupportedForm, engineName+" doesn't allow modifying spatial reference system SRID 0.")
	case !v.HasName:
		b.fail(rng, nil, review.UnsupportedForm, "Spatial reference system requires a "+review.Code("NAME")+" attribute.")
	case v.Name == "" || strings.TrimSpace(v.Name) != v.Name:
		b.fail(rng, nil, review.UnsupportedForm, "Spatial reference system "+review.Code("NAME")+" cannot be empty or padded with whitespace.")
	case !v.HasDefinition:
		b.fail(rng, nil, review.UnsupportedForm, "Spatial reference system requires a "+review.Code("DEFINITION")+" attribute.")
	case v.Definition == "":
		b.fail(rng, nil, review.UnsupportedForm, "Spatial reference system "+review.Code("DEFINITION")+" cannot be empty.")
	}
	return true
}

// eventScheduleFails refuses an event schedule expression that references a
// column, which has no table scope and always fails with 1054.
func (b *builder) eventScheduleFails(sched *ast.EventSchedule) bool {
	if sched == nil {
		return false
	}
	for _, expr := range []ast.ExprNode{sched.At, sched.Every, sched.Starts, sched.Ends} {
		if col, ok := exprHasBareColumnRef(expr); ok {
			b.fail(b.anchor(), nil, review.ObjectNotFound, "Column "+review.Code(ident(col))+" doesn't exist in the event schedule.")
			return true
		}
	}
	// EVERY takes a positive interval (1542), and ENDS can't come before
	// STARTS (1543), as literals show them (verified on mysql:8.0.46).
	every := sched.Every
	if iv, ok := every.(*ast.IntervalExpr); ok {
		every = iv.Value
	}
	if every != nil {
		if n, ok := partitionIntBound(every); ok && n <= 0 {
			b.fail(b.anchor(), nil, review.UnsupportedForm, "An event's "+review.Code("EVERY")+" interval must be positive.")
			return true
		}
	}
	start, ok1 := literalDatetime(sched.Starts)
	end, ok2 := literalDatetime(sched.Ends)
	if ok1 && ok2 && end < start {
		b.fail(b.anchor(), nil, review.UnsupportedForm, "An event's "+review.Code("ENDS")+" comes before its "+review.Code("STARTS")+".")
		return true
	}
	return false
}

// literalDatetime returns a datetime a string literal writes in full,
// 'YYYY-MM-DD hh:mm:ss' or 'YYYY-MM-DD', in a form that orders as text.
func literalDatetime(e ast.ExprNode) (string, bool) {
	s, ok := e.(*ast.StringLit)
	if !ok {
		return "", false
	}
	v := s.Value
	if len(v) == 10 {
		v += " 00:00:00"
	}
	if len(v) != 19 {
		return "", false
	}
	for i := 0; i < len(v); i++ {
		switch i {
		case 4, 7:
			if v[i] != '-' {
				return "", false
			}
		case 10:
			if v[i] != ' ' {
				return "", false
			}
		case 13, 16:
			if v[i] != ':' {
				return "", false
			}
		default:
			if v[i] < '0' || v[i] > '9' {
				return "", false
			}
		}
	}
	return v, true
}

const maxDisplayWidth = 4294967295

// columnDefTextFails checks column-definition errors MySQL establishes from the
// column definition itself (1439, 1221 SRID, 1067, 1294, 3107, 3767, 3812).
func (b *builder) columnDefTextFails(col *ast.ColumnDef) bool {
	if col == nil {
		return false
	}
	rng := rangeOf(col.Loc)
	if col.TypeName != nil {
		base := astBaseType(col.TypeName)
		maxLen := maxDisplayWidth
		switch base {
		case "tinyint", "smallint", "mediumint", "int", "bigint":
			maxLen = 255
		case "float", "double":
			if col.TypeName.Scale > 0 {
				maxLen = 255
			}
		case "bit":
			maxLen = 64
		}
		// FLOAT(p) keeps p in Length without HasLength.
		if col.TypeName.Length > maxLen || col.TypeName.HasLength && col.TypeName.Length < 0 {
			b.fail(rng, nil, review.UnsupportedForm, "Column "+review.Code(ident(col.Name))+" has a display width or length past "+engineName+"'s limit.")
			return true
		}
		if col.TypeName.HasSRID && !isGeometryBaseType(base) {
			b.fail(rng, nil, review.UnsupportedForm, "Column "+review.Code(ident(col.Name))+" cannot declare "+review.Code("SRID")+" on a non-geometry type.")
			return true
		}
		if col.OnUpdate != nil && base != "datetime" && base != "timestamp" {
			b.fail(rng, nil, review.UnsupportedForm, "Column "+review.Code(ident(col.Name))+" cannot use "+review.Code("ON UPDATE")+" with type "+review.Code(base)+".")
			return true
		}
		if eff := columnDefEffectiveDefaultExpr(col); eff != nil {
			if b.columnDefaultFails(col.Name, base, eff, rng) {
				return true
			}
		}
	}
	if col.Generated != nil && col.Generated.Expr != nil {
		if exprReferencesColumn(col.Generated.Expr, col.Name) {
			b.fail(rng, nil, review.UnsupportedForm, "Generated column "+review.Code(ident(col.Name))+" cannot refer to itself.")
			return true
		}
	}
	for _, cc := range col.Constraints {
		if cc != nil && cc.Type == ast.ColConstrCheck && cc.Expr != nil {
			if !checkExprIsBoolean(cc.Expr) {
				b.fail(rng, nil, review.UnsupportedForm, "Check constraint on "+review.Code(ident(col.Name))+" must be a boolean expression.")
				return true
			}
		}
	}
	return false
}

// columnDefEffectiveDefaultExpr returns the last DEFAULT expression declared on
// col, matching MySQL's last-attribute-wins rule.
func columnDefEffectiveDefaultExpr(col *ast.ColumnDef) ast.ExprNode {
	if col == nil {
		return nil
	}
	eff := col.DefaultValue
	for _, cc := range col.Constraints {
		if cc != nil && cc.Type == ast.ColConstrDefault && cc.Expr != nil {
			eff = cc.Expr
		}
	}
	return eff
}

// columnDefaultFails checks a DEFAULT expression against its column's name and
// known base type (1067 bare NOW() / invalid integer literal, 3767 self-ref).
func (b *builder) columnDefaultFails(colName, baseType string, expr ast.ExprNode, rng review.Range) bool {
	if expr == nil {
		return false
	}
	if paren, ok := expr.(*ast.ParenExpr); ok {
		if exprReferencesColumn(paren.Expr, colName) {
			b.fail(rng, nil, review.UnsupportedForm, "Default expression of column "+review.Code(ident(colName))+" cannot refer to itself.")
			return true
		}
		return false
	}
	if fn, ok := expr.(*ast.FuncCallExpr); ok {
		switch strings.ToLower(fn.Name) {
		case "now", "current_timestamp", "localtime", "localtimestamp":
			if baseType != "" && baseType != "datetime" && baseType != "timestamp" {
				b.fail(rng, nil, review.UnsupportedForm, "Column "+review.Code(ident(colName))+" has an invalid default value.")
				return true
			}
		}
	}
	if sl, ok := expr.(*ast.StringLit); ok && isIntegerBaseType(baseType) {
		if !isStrictNumericString(sl.Value) {
			b.fail(rng, nil, review.UnsupportedForm, "Column "+review.Code(ident(colName))+" has an invalid default value.")
			return true
		}
	}
	return false
}

// isStrictNumericString reports whether raw (after trimming ASCII whitespace)
// has valid MySQL decimal or scientific-notation numeric syntax (e.g. "42",
// "-1.5", ".5", "1e2") that MySQL can coerce to an integer default.
func isStrictNumericString(raw string) bool {
	s := strings.TrimSpace(raw)
	if s == "" {
		return false
	}
	i := 0
	if s[i] == '+' || s[i] == '-' {
		i++
	}
	digits := 0
	for i < len(s) && s[i] >= '0' && s[i] <= '9' {
		i++
		digits++
	}
	if i < len(s) && s[i] == '.' {
		i++
		for i < len(s) && s[i] >= '0' && s[i] <= '9' {
			i++
			digits++
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
			i++
			expDigits++
		}
		if expDigits == 0 {
			return false
		}
	}
	return i == len(s)
}

func colHasAutoIncrement(col *ast.ColumnDef) bool {
	if col == nil {
		return false
	}
	if col.AutoIncrement {
		return true
	}
	for _, cc := range col.Constraints {
		if cc != nil && cc.Type == ast.ColConstrAutoIncrement {
			return true
		}
	}
	return false
}

func colHasInlineKey(col *ast.ColumnDef) bool {
	if col == nil {
		return false
	}
	for _, cc := range col.Constraints {
		if cc != nil && (cc.Type == ast.ColConstrPrimaryKey || cc.Type == ast.ColConstrUnique) {
			return true
		}
	}
	return false
}

// indexColPart extracts a key part's column name ("" for functional key parts)
// and prefix length.
func indexColPart(ic *ast.IndexColumn) (name string, length int, functional bool) {
	if ic == nil {
		return "", 0, false
	}
	if ic.Functional {
		return "", 0, true
	}
	if cr, ok := ic.Expr.(*ast.ColumnRef); ok && !cr.Star && cr.Column != "" {
		return cr.Column, ic.Length, false
	}
	return "", 0, true
}

type keyPartSpec struct {
	name       string
	length     int
	functional bool
	expr       ast.ExprNode
}

func constraintKeyParts(c *ast.Constraint) []keyPartSpec {
	if c == nil {
		return nil
	}
	if len(c.IndexColumns) > 0 {
		out := make([]keyPartSpec, 0, len(c.IndexColumns))
		for _, ic := range c.IndexColumns {
			name, length, fn := indexColPart(ic)
			out = append(out, keyPartSpec{name: name, length: length, functional: fn, expr: ic.Expr})
		}
		return out
	}
	out := make([]keyPartSpec, 0, len(c.Columns))
	for _, col := range c.Columns {
		out = append(out, keyPartSpec{name: col})
	}
	return out
}

func indexStmtKeyParts(cols []*ast.IndexColumn) []keyPartSpec {
	out := make([]keyPartSpec, 0, len(cols))
	for _, ic := range cols {
		name, length, fn := indexColPart(ic)
		out = append(out, keyPartSpec{name: name, length: length, functional: fn, expr: ic.Expr})
	}
	return out
}

// keyPartsCoverColumn reports whether any non-functional key part names col.
func keyPartsCoverColumn(parts []keyPartSpec, col string) bool {
	for _, p := range parts {
		if !p.functional && strings.EqualFold(p.name, col) {
			return true
		}
	}
	return false
}

// tableOptionsFails checks table options that MySQL refuses unconditionally
// (3979 START TRANSACTION, 3184 ENCRYPTION) or on InnoDB (1478 COMPRESSION and
// ROW_FORMAT FIXED, 4025 AUTOEXTEND_SIZE).
func (b *builder) tableOptionsFails(opts []*ast.TableOption, finalEngine string, engineKnown bool, rng review.Range) bool {
	var rowFormat, compression, autoextend string
	for _, o := range opts {
		if o == nil {
			continue
		}
		orng := rangeOf(o.Loc)
		if orng == (review.Range{}) {
			orng = rng
		}
		switch strings.ToUpper(strings.TrimSpace(o.Name)) {
		case "START TRANSACTION":
			b.fail(orng, nil, review.UnsupportedForm, engineName+" doesn't accept "+review.Code("START TRANSACTION")+" as a table option.")
			return true
		case "ENCRYPTION":
			val := strings.TrimSpace(strings.Trim(o.Value, "'\""))
			if val != "" && !strings.EqualFold(val, "Y") && !strings.EqualFold(val, "N") {
				b.fail(orng, nil, review.UnsupportedForm, "Invalid "+review.Code("ENCRYPTION")+" option value "+review.Code(o.Value)+".")
				return true
			}
		case "ROW_FORMAT":
			rowFormat = strings.TrimSpace(strings.Trim(o.Value, "'\""))
		case "COMPRESSION":
			compression = strings.TrimSpace(strings.Trim(o.Value, "'\""))
		case "AUTOEXTEND_SIZE":
			autoextend = strings.TrimSpace(strings.Trim(o.Value, "'\""))
		}
	}
	if !engineKnown || !isInnoDBEngine(finalEngine) {
		return false
	}
	if compression != "" {
		switch strings.ToLower(compression) {
		case "zlib", "lz4", "none":
		default:
			b.fail(rng, nil, review.UnsupportedForm, "InnoDB doesn't support "+review.Code("COMPRESSION")+" value "+review.Code(compression)+".")
			return true
		}
	}
	if strings.EqualFold(rowFormat, "FIXED") && !b.w.innodbStrictOff {
		b.fail(rng, nil, review.UnsupportedForm, "InnoDB doesn't support "+review.Code("ROW_FORMAT = FIXED")+".")
		return true
	}
	if autoextend != "" {
		if bytes, ok := parseMySQLByteSize(autoextend); !ok || (bytes != 0 && (bytes < 4<<20 || bytes > 4096<<20)) {
			b.fail(rng, nil, review.UnsupportedForm, review.Code("AUTOEXTEND_SIZE")+" value "+review.Code(autoextend)+" is out of range for InnoDB.")
			return true
		}
	}
	return false
}

func parseMySQLByteSize(raw string) (uint64, bool) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return 0, false
	}
	var mult uint64 = 1
	switch s[len(s)-1] {
	case 'k', 'K':
		mult = 1 << 10
		s = s[:len(s)-1]
	case 'm', 'M':
		mult = 1 << 20
		s = s[:len(s)-1]
	case 'g', 'G':
		mult = 1 << 30
		s = s[:len(s)-1]
	}
	n, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, false
	}
	if n > ^uint64(0)/mult {
		return 0, false
	}
	return n * mult, true
}

// checkIndexParts validates an index or foreign key's key parts:
//   - functional key part on a bare column (3762)
//   - duplicate column in the same key part list (1060)
//   - FULLTEXT on a partitioned table (1214)
//   - column existence (1072) via lookupCol
//   - FULLTEXT / SPATIAL / RTREE / prefix column type compatibility (1283, 1687, 1089)
//
// lookupCol returns (state, baseType, fromDelta) for a column name: state is
// walk.Present, walk.Absent, or walk.Unknown; baseType is "" when unknown;
// fromDelta is true when state == walk.Absent came from this statement's own
// DROP/RENAME COLUMN (so a ledger Need would be suppressed by the touched rule
// and a direct Fail must be used instead).
func (b *builder) checkIndexParts(
	r tableRef,
	parts []keyPartSpec,
	fulltext, spatial, isFK bool,
	partitioned bool,
	rng review.Range,
	lookupCol func(colName string) (walk.Fact, string, bool),
) bool {
	for _, p := range parts {
		if p.functional && isBareColumnRef(p.expr) {
			b.fail(rng, nil, review.UnsupportedForm, "Functional index key part cannot be a bare column reference.")
			return true
		}
	}
	seen := make(map[string]bool, len(parts))
	for _, p := range parts {
		if p.functional || p.name == "" {
			continue
		}
		k, ok := b.key(spaceColumn, p.name)
		if !ok {
			continue
		}
		if seen[k] {
			b.fail(rng, nil, review.DuplicateName, statesTwice("names", "column", ident(p.name)))
			return true
		}
		seen[k] = true
	}
	for _, p := range parts {
		if p.functional || p.name == "" {
			continue
		}
		state, _, fromDelta := lookupCol(p.name)
		if state == walk.Absent {
			if fromDelta {
				b.fail(rng, nil, review.ObjectNotFound, notFound("key column", ident(p.name), r.written))
				return true
			}
			if k, ok := b.key(spaceColumn, p.name); ok {
				b.need(colKey(r.db, r.name, k), walk.Present, nil, rng, notFound("key column", ident(p.name), r.written), "")
				return false // ledger check recorded; caller continues so walk.Run reports it
			}
		}
	}
	if spatial && len(parts) > 1 {
		b.fail(rng, nil, review.UnsupportedForm, spatialOnePart) // 1070
		return true
	}
	if fulltext && partitioned {
		b.fail(rng, nil, review.UnsupportedForm, "Partitioned table "+review.Code(r.written)+" doesn't support "+review.Code("FULLTEXT")+" indexes.")
		return true
	}
	for _, p := range parts {
		if p.functional || p.name == "" {
			continue
		}
		state, baseType, _ := lookupCol(p.name)
		if state != walk.Present || baseType == "" {
			continue
		}
		if fulltext {
			if !isFulltextBaseType(baseType) {
				b.fail(rng, nil, review.UnsupportedForm, "Column "+review.Code(ident(p.name))+" cannot be part of a "+review.Code("FULLTEXT")+" index.")
				return true
			}
			continue
		}
		if spatial {
			if !isGeometryBaseType(baseType) {
				b.fail(rng, nil, review.UnsupportedForm, "Column "+review.Code(ident(p.name))+" cannot be part of a "+review.Code("SPATIAL")+" index.")
				return true
			}
			continue
		}
		if !isFK && p.length > 0 && !isPrefixableBaseType(baseType) {
			b.fail(rng, nil, review.UnsupportedForm, "Column "+review.Code(ident(p.name))+" cannot use an index prefix length.")
			return true
		}
	}
	return false
}

// partitionClausePreFKFails checks PARTITION BY clause errors MySQL raises
// regardless of table foreign keys (1492, 1500, 1652, partition bound 1054).
func (b *builder) partitionClausePreFKFails(pc *ast.PartitionClause, rng review.Range) bool {
	if pc == nil {
		return false
	}
	if (pc.Type == ast.PartitionRange || pc.Type == ast.PartitionList) && len(pc.Partitions) == 0 {
		b.fail(rng, nil, review.UnsupportedForm, "Each partition must be defined for "+review.Code("RANGE")+" or "+review.Code("LIST")+" partitioning.")
		return true
	}
	if (pc.Type == ast.PartitionHash || pc.Type == ast.PartitionKey) && partitionHasSubpartitions(pc) {
		b.fail(rng, nil, review.UnsupportedForm, review.Code("SUBPARTITION BY")+" is only supported with "+review.Code("RANGE")+" or "+review.Code("LIST")+" partitioning.")
		return true
	}
	if b.dupPartitionFields(pc.Columns, rng) || b.dupPartitionFields(pc.SubPartColumns, rng) {
		return true
	}
	return b.partitionDefsBoundFails(pc.Partitions, rng)
}

func (b *builder) dupPartitionFields(cols []string, rng review.Range) bool {
	seen := make(map[string]bool, len(cols))
	for _, c := range cols {
		k := strings.ToLower(c)
		if seen[k] {
			b.fail(rng, nil, review.DuplicateName, statesTwice("lists", "partition column", ident(c)))
			return true
		}
		seen[k] = true
	}
	return false
}

func (b *builder) partitionDefsBoundFails(defs []*ast.PartitionDef, rng review.Range) bool {
	for _, pd := range defs {
		if pd == nil {
			continue
		}
		if col, ok := partitionBoundColumn(pd.Values); ok {
			b.fail(rng, nil, review.ObjectNotFound, "Column "+review.Code(ident(col))+" doesn't exist in the partition function.")
			return true
		}
	}
	return false
}

// partitionClauseColumnsAndUniqueFails validates a PARTITION BY clause against
// the table's final columns (1488 / 1054) and unique keys (1503).
func (b *builder) partitionClauseColumnsAndUniqueFails(
	r tableRef,
	pc *ast.PartitionClause,
	rng review.Range,
	lookupCol func(colName string) walk.Fact,
	uniqueIndexes map[string]*indexMeta,
	indexesComplete bool,
) bool {
	if pc == nil {
		return false
	}
	var partCols []string
	addPartCol := func(name string) {
		k := strings.ToLower(name)
		for _, existing := range partCols {
			if existing == k {
				return
			}
		}
		partCols = append(partCols, k)
	}
	checkList := func(cols []string, isKey bool) bool {
		for _, c := range cols {
			if lookupCol(c) == walk.Absent {
				if isKey {
					b.fail(rng, nil, review.ObjectNotFound, "Partition column "+review.Code(ident(c))+" doesn't exist in "+review.Code(r.written)+".")
				} else {
					b.fail(rng, nil, review.ObjectNotFound, "Column "+review.Code(ident(c))+" doesn't exist in the partition function.")
				}
				return true
			}
			addPartCol(c)
		}
		return false
	}
	if checkList(pc.Columns, pc.Type == ast.PartitionKey) {
		return true
	}
	for _, c := range exprColumnNames(pc.Expr) {
		if lookupCol(c) == walk.Absent {
			b.fail(rng, nil, review.ObjectNotFound, "Column "+review.Code(ident(c))+" doesn't exist in the partition function.")
			return true
		}
		addPartCol(c)
	}
	if checkList(pc.SubPartColumns, pc.SubPartType == ast.PartitionKey) {
		return true
	}
	for _, c := range exprColumnNames(pc.SubPartExpr) {
		if lookupCol(c) == walk.Absent {
			b.fail(rng, nil, review.ObjectNotFound, "Column "+review.Code(ident(c))+" doesn't exist in the partition function.")
			return true
		}
		addPartCol(c)
	}
	if !indexesComplete || len(partCols) == 0 {
		return false
	}
	for _, idx := range uniqueIndexes {
		if idx == nil || (!idx.unique && !idx.primary) || idx.cols == nil || idx.anyPrefix() {
			continue
		}
		for _, pcCol := range partCols {
			covered := false
			for _, ic := range idx.cols {
				if ic == pcCol {
					covered = true
					break
				}
			}
			if !covered {
				b.fail(rng, nil, review.UnsupportedForm, "Every unique key of "+review.Code(r.written)+" must include all columns in its partitioning function.")
				return true
			}
		}
	}
	return false
}

// indexCoversCols reports whether idx is a BTREE/HASH index (not FULLTEXT or
// SPATIAL) whose leading whole-column (unprefixed) key parts match fkCols in
// order. Prefix lengths on trailing columns past len(fkCols) do not disqualify
// the index.
func indexCoversCols(idx *indexMeta, fkCols []string) bool {
	if idx == nil || idx.fulltext || idx.spatial || len(fkCols) == 0 || len(idx.cols) < len(fkCols) {
		return false
	}
	for i, c := range fkCols {
		if idx.cols[i] != c {
			return false
		}
		if i < len(idx.colPrefixes) && idx.colPrefixes[i] {
			return false
		}
	}
	return true
}

// dropIndexNeededByFK reports whether dropping index dropKey (with optional
// additional dropped indexes, dropped FKs, dropped/renamed columns, and newly
// added indexes in the same ALTER TABLE) leaves a surviving FK on rel without
// any backing index (error 1553).
func (b *builder) dropIndexNeededByFK(
	r tableRef,
	dropKey string,
	droppedIndexes map[string]bool,
	droppedFKs map[string]bool,
	droppedCols map[string]bool,
	renamedFrom map[string]string,
	addedIndexes []*indexMeta,
) bool {
	rel := relKey(r.db, r.name)
	if !b.fks || !b.w.tableIndexesComplete[rel] || len(droppedCols) > 0 {
		return false
	}
	idxs := b.w.tableIndexes[rel]
	fks := b.w.tableFKs[rel]
	if len(idxs) == 0 || len(fks) == 0 {
		return false
	}
	droppedIdx := idxs[dropKey]
	if droppedIdx == nil {
		return false
	}
	for fkName, fkCols := range fks {
		if droppedFKs[fkName] || len(fkCols) == 0 {
			continue
		}
		if !indexCoversCols(droppedIdx, fkCols) {
			continue
		}
		coveredElsewhere := false
		for name, other := range idxs {
			if name == dropKey || droppedIndexes[name] {
				continue
			}
			if indexCoversCols(other, fkCols) {
				coveredElsewhere = true
				break
			}
		}
		if !coveredElsewhere {
			finalFKCols := fkCols
			if len(renamedFrom) > 0 {
				finalFKCols = make([]string, len(fkCols))
				for i, fc := range fkCols {
					if newK, ok := renamedFrom[fc]; ok {
						finalFKCols[i] = newK
					} else {
						finalFKCols[i] = fc
					}
				}
			}
			for _, added := range addedIndexes {
				if indexCoversCols(added, finalFKCols) {
					coveredElsewhere = true
					break
				}
			}
		}
		if !coveredElsewhere {
			return true
		}
	}
	return false
}

// tableHasSurvivingFK reports whether rel has any foreign key (or is referenced
// by another table's foreign key) after excluding droppedFKs.
func (b *builder) tableHasSurvivingFK(rel walk.Key, droppedFKs map[string]bool) bool {
	for fkName := range b.w.tableFKs[rel] {
		if !droppedFKs[fkName] {
			return true
		}
	}
	for ref := range b.w.referencedByFK[rel] {
		if ref.child == rel && droppedFKs[ref.fk] {
			continue
		}
		return true
	}
	return false
}
