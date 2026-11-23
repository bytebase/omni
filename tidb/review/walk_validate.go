package review

import (
	"math"
	"strings"

	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/tidb/parser"
)

// tidbValidEngines is TiDB v8.5.5's storage engine allowlist
// (pkg/planner/core/preprocess.go checkTableEngine). Anything outside this
// set fails with 1286 ("Unknown storage engine").
var tidbValidEngines = map[string]bool{
	"archive":    true,
	"blackhole":  true,
	"csv":        true,
	"example":    true,
	"federated":  true,
	"innodb":     true,
	"memory":     true,
	"merge":      true,
	"mgr_myisam": true,
	"myisam":     true,
	"ndb":        true,
	"heap":       true,
}

// tidbSupportedAlterTableOptions is the set of ATTableOption names TiDB v8.5.5
// accepts on ALTER TABLE (anything outside this set, other than UNION/8232 and
// INSERT_METHOD/8233, fails with 8200).
var tidbSupportedAlterTableOptions = map[string]bool{
	"engine":                true,
	"storage engine":        true,
	"charset":               true,
	"character set":         true,
	"default charset":       true,
	"default character set": true,
	"collate":               true,
	"default collate":       true,
	"comment":               true,
	"auto_increment":        true,
	"row_format":            true,
	"shard_row_id_bits":     true,
	"auto_id_cache":         true,
	"auto_random_base":      true,
	"placement policy":      true,
	"ttl":                   true,
	"ttl_enable":            true,
	"ttl_job_interval":      true,
	"affinity":              true,
}

// tidbSupportedCharsets is the set of character sets TiDB v8.5.5 recognizes
// (including utf8mb3 as an alias for utf8). Any other charset fails with 1115.
var tidbSupportedCharsets = map[string]bool{
	"ascii":   true,
	"binary":  true,
	"gbk":     true,
	"latin1":  true,
	"utf8":    true,
	"utf8mb3": true,
	"utf8mb4": true,
}

// tidbNowDefaultFuncs are the CURRENT_TIMESTAMP / NOW family defaults that
// TiDB v8.5.5 allows only on DATETIME and TIMESTAMP columns (1067 otherwise).
var tidbNowDefaultFuncs = map[string]bool{
	"now":               true,
	"current_timestamp": true,
	"localtime":         true,
	"localtimestamp":    true,
}

// tidbAllowedDefaultExprFuncs is TiDB v8.5.5's allowlist of functions permitted
// as the top-level call in a column DEFAULT expression (3770 otherwise).
var tidbAllowedDefaultExprFuncs = map[string]bool{
	"current_timestamp": true,
	"now":               true,
	"localtime":         true,
	"localtimestamp":    true,
	"current_date":      true,
	"curdate":           true,
	"nextval":           true,
	"rand":              true,
	"uuid":              true,
	"uuid_to_bin":       true,
	"date_format":       true,
	"replace":           true,
	"upper":             true,
	"str_to_date":       true,
	"json_object":       true,
	"json_array":        true,
	"json_quote":        true,
	"vec_from_text":     true,
}

// unsupportedAlterActionClause returns the SQL clause name for an ALTER TABLE
// action that TiDB v8.5.5 unconditionally refuses with 8200, or "" when the
// action is supported.
func unsupportedAlterActionClause(t ast.AlterTableCmdType) string {
	switch t {
	case ast.ATCheckPartition:
		return "CHECK PARTITION"
	case ast.ATOptimizePartition:
		return "OPTIMIZE PARTITION"
	case ast.ATRebuildPartition:
		return "REBUILD PARTITION"
	case ast.ATRepairPartition:
		return "REPAIR PARTITION"
	case ast.ATReorganizeLastPartition:
		return "SPLIT MAXVALUE PARTITION"
	case ast.ATReorganizeFirstPartition:
		return "MERGE FIRST PARTITION"
	case ast.ATDiscardPartitionTablespace:
		return "DISCARD PARTITION TABLESPACE"
	case ast.ATImportPartitionTablespace:
		return "IMPORT PARTITION TABLESPACE"
	case ast.ATDiscardTablespace:
		return "DISCARD TABLESPACE"
	case ast.ATImportTablespace:
		return "IMPORT TABLESPACE"
	case ast.ATForce:
		return "FORCE"
	case ast.ATSecondaryLoad:
		return "SECONDARY_LOAD"
	case ast.ATSecondaryUnload:
		return "SECONDARY_UNLOAD"
	case ast.ATStatsOptions:
		return "STATS_OPTIONS"
	}
	return ""
}

// checkTableOption validates a table option on CREATE TABLE or ALTER TABLE.
// When alter is true, options outside tidbSupportedAlterTableOptions are also
// refused (8200). Returns true when it records a failure.
func (b *builder) checkTableOption(r tableRef, opt *ast.TableOption, rng review.Range, alter bool) bool {
	if opt == nil {
		return false
	}
	name := strings.ToLower(opt.Name)
	switch name {
	case "union":
		b.fail(rng, nil, review.UnsupportedForm, engineName+" doesn't support table option "+review.Code("UNION")+".")
		return true
	case "insert_method":
		b.fail(rng, nil, review.UnsupportedForm, engineName+" doesn't support table option "+review.Code("INSERT_METHOD")+".")
		return true
	case "engine", "storage engine":
		if !tidbValidEngines[strings.ToLower(opt.Value)] {
			b.fail(rng, nil, review.UnsupportedForm, engineName+" doesn't recognize storage engine "+review.Code(opt.Value)+".")
			return true
		}
		return false
	case "affinity":
		switch strings.ToLower(opt.Value) {
		case "", "none", "table", "partition":
		default:
			b.fail(rng, nil, review.UnsupportedForm, engineName+" doesn't recognize "+review.Code("AFFINITY")+" value "+review.Code(opt.Value)+".")
			return true
		}
		return false
	case "charset", "character set", "default charset", "default character set":
		if opt.Value != "" && !tidbSupportedCharsets[strings.ToLower(opt.Value)] {
			b.fail(rng, nil, review.UnsupportedForm, engineName+" doesn't recognize character set "+review.Code(opt.Value)+".")
			return true
		}
		return false
	case "ttl":
		if alter && r.name != "" {
			b.checkTTLColumnExists(r, opt.Value, rng)
		}
		return false
	}
	if alter && !tidbSupportedAlterTableOptions[name] {
		b.fail(rng, nil, review.UnsupportedForm, engineName+" doesn't support table option "+review.Code(strings.ToUpper(opt.Name))+" in "+review.Code("ALTER TABLE")+".")
		return true
	}
	return false
}

// checkTTLColumnExists checks that the column named by an ALTER TABLE ... TTL
// expression exists in the table (1054).
func (b *builder) checkTTLColumnExists(r tableRef, val string, rng review.Range) {
	expr, err := parser.ParseExpr(val)
	if err != nil {
		return
	}
	var colName string
	switch e := expr.(type) {
	case *ast.BinaryExpr:
		if e.Op == ast.BinOpAdd {
			if c, ok := e.Left.(*ast.ColumnRef); ok {
				colName = c.Column
			}
		}
	case *ast.FuncCallExpr:
		fn := strings.ToUpper(e.Name)
		if (fn == "DATE_ADD" || fn == "ADDDATE") && len(e.Args) == 2 {
			if c, ok := e.Args[0].(*ast.ColumnRef); ok {
				colName = c.Column
			}
		}
	}
	if colName != "" {
		b.columnNeeds(r, colName, rng)
	}
}

// checkAttributesSpec validates TiDB's ATTRIBUTES = '...' string: non-empty
// comma-separated key=value entries with no conflicting values for the same
// key (8237). Returns true when it records a failure.
func (b *builder) checkAttributesSpec(spec *ast.AttributesSpec, rng review.Range) bool {
	if spec == nil || spec.Default || strings.TrimSpace(spec.Value) == "" {
		return false
	}
	seen := make(map[string]string)
	for raw := range strings.SplitSeq(spec.Value, ",") {
		part := strings.TrimSpace(raw)
		if part == "" {
			continue
		}
		kv := strings.Split(part, "=")
		if len(kv) != 2 || strings.TrimSpace(kv[0]) == "" || strings.TrimSpace(kv[1]) == "" {
			b.fail(rng, nil, review.UnsupportedForm, engineName+" requires "+review.Code("ATTRIBUTES")+" entries in "+review.Code("key=value")+" format without conflicts.")
			return true
		}
		k := strings.TrimSpace(kv[0])
		v := strings.TrimSpace(kv[1])
		if prev, ok := seen[k]; ok && prev != v {
			b.fail(rng, nil, review.UnsupportedForm, engineName+" requires "+review.Code("ATTRIBUTES")+" entries in "+review.Code("key=value")+" format without conflicts.")
			return true
		}
		seen[k] = v
	}
	return false
}

// checkAddColumnUnsupportedConstraints refuses inline constraints that TiDB
// v8.5.5 rejects on ALTER TABLE ADD COLUMN ahead of duplicate-column and
// multiple-PK checks: AUTO_RANDOM (8216), AUTO_INCREMENT / SERIAL, PRIMARY KEY,
// and UNIQUE KEY (8200). Returns true when it records a failure.
func (b *builder) checkAddColumnUnsupportedConstraints(col *ast.ColumnDef, rng review.Range) bool {
	if col == nil {
		return false
	}
	if col.AutoRandom {
		b.fail(rng, nil, review.UnsupportedForm, engineName+" doesn't support adding column "+review.Code(ident(col.Name))+" with "+review.Code("AUTO_RANDOM")+".")
		return true
	}
	if col.AutoIncrement || (col.TypeName != nil && strings.EqualFold(col.TypeName.Name, "serial")) {
		b.fail(rng, nil, review.UnsupportedForm, engineName+" doesn't support adding column "+review.Code(ident(col.Name))+" with "+review.Code("AUTO_INCREMENT")+".")
		return true
	}
	for _, cc := range col.Constraints {
		switch cc.Type {
		case ast.ColConstrAutoIncrement:
			b.fail(rng, nil, review.UnsupportedForm, engineName+" doesn't support adding column "+review.Code(ident(col.Name))+" with "+review.Code("AUTO_INCREMENT")+".")
			return true
		case ast.ColConstrPrimaryKey:
			b.fail(rng, nil, review.UnsupportedForm, engineName+" doesn't support adding column "+review.Code(ident(col.Name))+" with "+review.Code("PRIMARY KEY")+".")
			return true
		case ast.ColConstrUnique:
			b.fail(rng, nil, review.UnsupportedForm, engineName+" doesn't support adding column "+review.Code(ident(col.Name))+" with "+review.Code("UNIQUE")+".")
			return true
		}
	}
	return false
}

// colDefaultExpr returns the effective default expression on a ColumnDef.
func colDefaultExpr(col *ast.ColumnDef) ast.ExprNode {
	expr := col.DefaultValue
	for _, cc := range col.Constraints {
		if cc.Type == ast.ColConstrDefault && cc.Expr != nil {
			expr = cc.Expr
		}
	}
	return expr
}

// unwrapParens strips surrounding ParenExpr wrappers.
func unwrapParens(expr ast.ExprNode) ast.ExprNode {
	for {
		p, ok := expr.(*ast.ParenExpr)
		if !ok {
			return expr
		}
		expr = p.Expr
	}
}

// checkColumnDefPreprocess checks column definition properties that TiDB
// validates during statement preparation (charset 1115, NOW()-family default
// type 1067). Returns true when it records a failure.
func (b *builder) checkColumnDefPreprocess(col *ast.ColumnDef, rng review.Range) bool {
	if col == nil {
		return false
	}
	if dt := col.TypeName; dt != nil {
		if dt.Charset != "" && !tidbSupportedCharsets[strings.ToLower(dt.Charset)] {
			b.fail(rng, nil, review.UnsupportedForm, engineName+" doesn't recognize character set "+review.Code(dt.Charset)+".")
			return true
		}
	}
	if expr := colDefaultExpr(col); expr != nil && col.TypeName != nil {
		if fn, ok := unwrapParens(expr).(*ast.FuncCallExpr); ok && tidbNowDefaultFuncs[strings.ToLower(fn.Name)] {
			switch strings.ToLower(col.TypeName.Name) {
			case "datetime", "timestamp":
				funcFsp := 0
				if len(fn.Args) > 0 {
					if lit, ok := fn.Args[0].(*ast.IntLit); ok {
						funcFsp = int(lit.Value)
					}
				}
				if funcFsp != col.TypeName.Length {
					b.fail(rng, nil, review.UnsupportedForm, review.Code(ident(col.Name))+" has an invalid default value for its type.")
					return true
				}
			default:
				b.fail(rng, nil, review.UnsupportedForm, review.Code(ident(col.Name))+" has an invalid default value for its type.")
				return true
			}
		}
	}
	return false
}

// checkAlterColumnDef checks column definition rules that TiDB v8.5.5 enforces
// on ALTER TABLE ADD/MODIFY/CHANGE COLUMN: preprocess-time charset (1115) and
// NOW()-family default (1067), sequence default on ADD COLUMN (1146/1347/8230),
// and — when skipPostLookup is false (the target column lookup itself did not
// already fail) — STORED generated (3106), ON UPDATE type/precision (1294),
// CURDATE() default on narrow integer (1067), disallowed default function
// (3770), duplicate ENUM/SET members (1291), and generated column references
// (1054).
func (b *builder) checkAlterColumnDef(r tableRef, col *ast.ColumnDef, rng review.Range, isAdd, skipPostLookup bool) bool {
	if col == nil {
		return false
	}
	if b.checkColumnDefPreprocess(col, rng) {
		return true
	}
	if isAdd {
		if b.checkAddColumnSequenceDefault(r, col, rng) {
			return true
		}
	}
	if skipPostLookup {
		return false
	}
	if col.Generated != nil && col.Generated.Stored {
		b.fail(rng, nil, review.UnsupportedForm, engineName+" doesn't support a "+review.Code("STORED")+" generated column in "+review.Code("ALTER TABLE")+".")
		return true
	}
	dt := col.TypeName
	if col.OnUpdate != nil && dt != nil {
		switch strings.ToLower(dt.Name) {
		case "datetime", "timestamp":
			if fn, ok := col.OnUpdate.(*ast.FuncCallExpr); ok {
				updFsp := 0
				if len(fn.Args) > 0 {
					if lit, ok := fn.Args[0].(*ast.IntLit); ok {
						updFsp = int(lit.Value)
					}
				}
				if updFsp != dt.Length {
					b.fail(rng, nil, review.UnsupportedForm, review.Code(ident(col.Name))+" has an invalid "+review.Code("ON UPDATE")+" clause for its type.")
					return true
				}
			}
		default:
			b.fail(rng, nil, review.UnsupportedForm, review.Code(ident(col.Name))+" has an invalid "+review.Code("ON UPDATE")+" clause for its type.")
			return true
		}
	}
	if expr := colDefaultExpr(col); expr != nil {
		inner := unwrapParens(expr)
		if fn, ok := inner.(*ast.FuncCallExpr); ok {
			fnName := strings.ToLower(fn.Name)
			if !tidbAllowedDefaultExprFuncs[fnName] {
				b.fail(rng, nil, review.UnsupportedForm, engineName+" doesn't allow "+review.Code(fn.Name)+" in a column "+review.Code("DEFAULT")+" expression.")
				return true
			}
			if (fnName == "curdate" || fnName == "current_date") && dt != nil {
				switch strings.ToLower(dt.Name) {
				case "tinyint", "smallint", "mediumint", "bool", "boolean":
					b.fail(rng, nil, review.UnsupportedForm, review.Code(ident(col.Name))+" has an invalid default value for its type.")
					return true
				}
			}
		}
	}
	if dt != nil && len(dt.EnumValues) > 0 {
		kind := "ENUM"
		if strings.EqualFold(dt.Name, "set") {
			kind = "SET"
		}
		seen := make(map[string]bool, len(dt.EnumValues))
		for _, v := range dt.EnumValues {
			if seen[v] {
				b.fail(rng, nil, review.DuplicateName, review.Code(ident(col.Name))+" declares duplicate "+kind+" value "+review.Code(v)+".")
				return true
			}
			seen[v] = true
		}
	}
	if col.Generated != nil && col.Generated.Expr != nil {
		b.checkExprColumnsExist(r, col.Generated.Expr, rng)
	}
	return false
}

// checkAddColumnSequenceDefault resolves any NEXT VALUE FOR / NEXTVAL sequence
// reference in an ALTER TABLE ADD COLUMN default: missing (1146), wrong kind
// (1347), or existing sequence (8230, since ADD COLUMN does not support
// sequence defaults).
func (b *builder) checkAddColumnSequenceDefault(r tableRef, col *ast.ColumnDef, rng review.Range) bool {
	expr := colDefaultExpr(col)
	if expr == nil {
		return false
	}
	var seqRef *ast.ColumnRef
	ast.Inspect(expr, func(n ast.Node) bool {
		if seqRef != nil {
			return false
		}
		if fn, ok := n.(*ast.FuncCallExpr); ok && strings.EqualFold(fn.Name, "nextval") && len(fn.Args) == 1 {
			if arg, ok := fn.Args[0].(*ast.ColumnRef); ok && arg.Column != "" {
				seqRef = arg
				return false
			}
		}
		return true
	})
	if seqRef == nil {
		return false
	}
	tr := &ast.TableRef{Schema: seqRef.Table, Name: seqRef.Column, Loc: seqRef.Loc}
	if tr.Schema == "" {
		tr.Schema = r.db
		if strings.HasPrefix(r.db, nonASCII) {
			return false
		}
	}
	sr, st := b.ref(tr)
	if st != refOK {
		return false
	}
	sr.written = ident(seqRef.Column)
	if seqRef.Table != "" {
		sr.written = ident(seqRef.Table) + "." + ident(seqRef.Column)
	}
	k := relKey(sr.db, sr.name)
	f, kind := b.l.Lookup(k)
	switch {
	case f == walk.Absent, f == walk.Present && kind != "" && kind != relSequence:
		b.needKind(k, []string{relSequence}, rng, notFound("sequence", sr.written, ""), sr.written, "sequence")
		return true
	case f == walk.Present && kind == relSequence:
		b.fail(rng, []walk.Key{k}, review.UnsupportedForm, engineName+" doesn't support a sequence default in "+review.Code("ALTER TABLE ADD COLUMN")+".")
		return true
	}
	return false
}

// checkExprColumnsExist adds a Present check on every column referenced in
// expr against table r (1054), and reports whether any referenced column is
// already known Absent in the ledger.
func (b *builder) checkExprColumnsExist(r tableRef, expr ast.ExprNode, rng review.Range) bool {
	if expr == nil {
		return false
	}
	missing := false
	ast.Inspect(expr, func(n ast.Node) bool {
		if ref, ok := n.(*ast.ColumnRef); ok && ref.Column != "" {
			if b.columnNeeds(r, ref.Column, rng) {
				if k, ok := b.key(spaceColumn, ref.Column); ok {
					if f, _ := b.l.Lookup(colKey(r.db, r.name, k)); f == walk.Absent {
						missing = true
					}
				}
			}
		}
		return true
	})
	return missing
}

// checkIndexColumnsExist checks that every non-functional column key part in
// cols exists in table r (1072).
func (b *builder) checkIndexColumnsExist(r tableRef, cols []*ast.IndexColumn, names []string, rng review.Range) {
	if len(cols) > 0 {
		for _, ic := range cols {
			if ic == nil || ic.Functional {
				continue
			}
			if cr, ok := ic.Expr.(*ast.ColumnRef); ok && cr.Column != "" {
				b.columnNeeds(r, cr.Column, rng)
			}
		}
		return
	}
	for _, name := range names {
		if name != "" {
			b.columnNeeds(r, name, rng)
		}
	}
}

// checkConstraintForeignKeyShape refuses a table-level FOREIGN KEY that omits
// its referenced columns (1239) or references a functional expression (1105).
func (b *builder) checkConstraintForeignKeyShape(c *ast.Constraint, rng review.Range) bool {
	if c == nil || c.Type != ast.ConstrForeignKey {
		return false
	}
	for _, kp := range c.RefKeyParts {
		if kp != nil && kp.Functional {
			b.fail(rng, nil, review.UnsupportedForm, review.Code("FOREIGN KEY")+" can't reference a functional expression.")
			return true
		}
	}
	if len(c.RefColumns) == 0 {
		b.fail(rng, nil, review.UnsupportedForm, review.Code("FOREIGN KEY")+" must name the referenced columns.")
		return true
	}
	return false
}

// checkPartitionBy checks column existence (1054/1488) and INTERVAL bounds
// (8200/1054) on an ALTER TABLE ... PARTITION BY clause. Column existence runs
// first in TiDB v8.5.5, so if a referenced column is already known Absent, the
// column check is left to report 1054/1488 rather than being overridden by an
// INTERVAL failure.
func (b *builder) checkPartitionBy(r tableRef, pc *ast.PartitionClause, rng review.Range) bool {
	if pc == nil || pc.Type == ast.PartitionSystemTime {
		return false
	}
	missingCol := false
	if len(pc.Columns) > 0 {
		for _, name := range pc.Columns {
			if b.columnNeeds(r, name, rng) {
				if k, ok := b.key(spaceColumn, name); ok {
					if f, _ := b.l.Lookup(colKey(r.db, r.name, k)); f == walk.Absent {
						missingCol = true
					}
				}
			}
		}
	} else if pc.Expr != nil {
		missingCol = b.checkExprColumnsExist(r, pc.Expr, rng)
	}
	if missingCol {
		return false
	}
	if pc.Interval != nil {
		if pc.Interval.First == nil || pc.Interval.Last == nil {
			b.fail(rng, nil, review.UnsupportedForm, engineName+" requires "+review.Code("FIRST")+" and "+review.Code("LAST")+" partitions in "+review.Code("INTERVAL")+" partitioning.")
			return true
		}
		for _, bound := range []ast.ExprNode{pc.Interval.First, pc.Interval.Last} {
			if name := firstExprColumnRef(bound); name != "" {
				b.fail(rng, nil, review.UnsupportedForm, "Partition bound expression can't reference column "+review.Code(ident(name))+".")
				return true
			}
		}
	}
	return false
}

// firstExprColumnRef returns the first column reference in expr, or "".
func firstExprColumnRef(expr ast.ExprNode) string {
	if expr == nil {
		return ""
	}
	var name string
	ast.Inspect(expr, func(n ast.Node) bool {
		if ref, ok := n.(*ast.ColumnRef); ok && ref.Column != "" {
			name = ref.Column
			return false
		}
		return true
	})
	return name
}

// firstPartitionValueColumnRef returns the first column reference inside a
// partition bound value node (single expression or tuple/list), or "".
func firstPartitionValueColumnRef(v ast.Node) string {
	switch n := v.(type) {
	case *ast.List:
		for _, item := range n.Items {
			if name := firstPartitionValueColumnRef(item); name != "" {
				return name
			}
		}
	case ast.ExprNode:
		return firstExprColumnRef(n)
	}
	return ""
}

// checkPartitionDefs validates partition definitions in ADD PARTITION /
// REORGANIZE PARTITION: TiDB v8.5.5 does not support SYSTEM_TIME HISTORY or
// CURRENT partitions (4113/1480/1505), rejects duplicate partition or
// subpartition names in one statement (1517/1479/1505), and rejects column
// references in VALUES LESS THAN bounds (1054/1480/1505).
func (b *builder) checkPartitionDefs(defs []*ast.PartitionDef, rng review.Range) bool {
	seenPart := make(map[string]bool, len(defs))
	seenSub := make(map[string]bool)
	for _, pd := range defs {
		if pd == nil {
			continue
		}
		if s, ok := pd.Values.(*ast.String); ok && (s.Str == "HISTORY" || s.Str == "CURRENT") {
			b.fail(rng, nil, review.UnsupportedForm, engineName+" doesn't support "+review.Code(s.Str)+" partitions.")
			return true
		}
		if pd.Name != "" {
			k := strings.ToLower(pd.Name)
			if seenPart[k] {
				b.fail(rng, nil, review.DuplicateName, statesTwice("declares", "partition", ident(pd.Name)))
				return true
			}
			seenPart[k] = true
		}
		for _, sp := range pd.SubPartitions {
			if sp == nil || sp.Name == "" {
				continue
			}
			sk := strings.ToLower(sp.Name)
			if seenSub[sk] {
				b.fail(rng, nil, review.DuplicateName, statesTwice("declares", "subpartition", ident(sp.Name)))
				return true
			}
			seenSub[sk] = true
		}
		if pd.BoundKind == ast.PartBoundLessThan {
			if name := firstPartitionValueColumnRef(pd.Values); name != "" {
				b.fail(rng, nil, review.UnsupportedForm, "Partition bound expression can't reference column "+review.Code(ident(name))+".")
				return true
			}
		}
	}
	return false
}

// checkSequenceOptionsValid checks self-contradictory sequence options that
// TiDB v8.5.5 always refuses with 4136, regardless of any existing state:
// CACHE <= 0, INCREMENT 0, MINVALUE == math.MinInt64, MAXVALUE == math.MaxInt64.
func (b *builder) checkSequenceOptionsValid(written string, opts []*ast.SequenceOption, rng review.Range) bool {
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		switch opt.Type {
		case ast.SequenceOptionCache:
			if opt.Value <= 0 {
				b.fail(rng, nil, review.UnsupportedForm, "Sequence "+review.Code(written)+" has conflicting options.")
				return true
			}
		case ast.SequenceOptionIncrementBy:
			if opt.Value == 0 {
				b.fail(rng, nil, review.UnsupportedForm, "Sequence "+review.Code(written)+" has conflicting options.")
				return true
			}
		case ast.SequenceOptionMinValue:
			if opt.Value == math.MinInt64 {
				b.fail(rng, nil, review.UnsupportedForm, "Sequence "+review.Code(written)+" has conflicting options.")
				return true
			}
		case ast.SequenceOptionMaxValue:
			if opt.Value == math.MaxInt64 {
				b.fail(rng, nil, review.UnsupportedForm, "Sequence "+review.Code(written)+" has conflicting options.")
				return true
			}
		}
	}
	return false
}

// viewBodyHasNoTableStar reports whether a CREATE VIEW's SELECT body contains
// a VALUES query primary or an unexpandable `*` / `t.*` with no matching table
// in its FROM clause, which TiDB v8.5.5 always refuses with 1051 before checking
// whether the view name is already taken.
func viewBodyHasNoTableStar(sel *ast.SelectStmt, source string) bool {
	if sel == nil {
		return false
	}
	bad := false
	ast.Inspect(sel, func(n ast.Node) bool {
		if bad {
			return false
		}
		block, ok := n.(*ast.SelectStmt)
		if !ok {
			return true
		}
		if viewBlockHasNoTableStar(block, source) {
			bad = true
			return false
		}
		return true
	})
	return bad
}

func viewBlockHasNoTableStar(sel *ast.SelectStmt, source string) bool {
	if sel.ValuesSource != nil {
		return true
	}
	if viewBlockReadsOnlyNoTable(sel, source) {
		for _, field := range sel.TargetList {
			expr := field
			if rt, ok := field.(*ast.ResTarget); ok {
				expr = rt.Val
			}
			switch e := expr.(type) {
			case *ast.StarExpr:
				return true
			case *ast.ColumnRef:
				if e.Star {
					return true
				}
			}
		}
		return false
	}
	relations := viewBlockRelationNames(sel)
	for _, field := range sel.TargetList {
		expr := field
		if rt, ok := field.(*ast.ResTarget); ok {
			expr = rt.Val
		}
		if e, ok := expr.(*ast.ColumnRef); ok && e.Star {
			q := strings.ToLower(e.Table)
			if e.Schema != "" {
				q = strings.ToLower(e.Schema) + "." + q
			}
			if !relations[q] {
				return true
			}
		}
	}
	return false
}

func viewBlockReadsOnlyNoTable(sel *ast.SelectStmt, source string) bool {
	if len(sel.From) > 1 {
		return false
	}
	if len(sel.From) == 0 {
		return true
	}
	ref, ok := sel.From[0].(*ast.TableRef)
	if !ok || ref.Schema != "" || !strings.EqualFold(ref.Name, "dual") {
		return false
	}
	if ref.Loc.Start < 0 || ref.Loc.Start >= len(source) {
		return false
	}
	return source[ref.Loc.Start] != '`'
}

func viewBlockRelationNames(sel *ast.SelectStmt) map[string]bool {
	names := make(map[string]bool, len(sel.From))
	var add func(ast.TableExpr)
	add = func(te ast.TableExpr) {
		switch e := te.(type) {
		case *ast.TableRef:
			if e.Alias != "" {
				names[strings.ToLower(e.Alias)] = true
				return
			}
			if e.Name != "" {
				names[strings.ToLower(e.Name)] = true
			}
			if e.Schema != "" && e.Name != "" {
				names[strings.ToLower(e.Schema)+"."+strings.ToLower(e.Name)] = true
			}
		case *ast.SubqueryExpr:
			if e.Alias != "" {
				names[strings.ToLower(e.Alias)] = true
			}
		case *ast.JoinClause:
			add(e.Left)
			add(e.Right)
		}
	}
	for _, te := range sel.From {
		add(te)
	}
	return names
}
