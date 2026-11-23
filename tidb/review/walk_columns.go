package review

import (
	"strings"

	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
	"github.com/bytebase/omni/tidb/ast"
)

// alterShapeFails fails an ALTER TABLE that the table's shape before it
// refutes (walk_table_shape.go): TiDB judges every subcommand against the
// table as it was, so a key or a foreign key the same statement drops
// still counts (verified). Each refusal was verified on v8.5.5.
func (b *builder) alterShapeFails(v *ast.AlterTableStmt, r tableRef) bool {
	if b.idCacheFails(v, r) {
		return true
	}
	for _, c := range v.Commands {
		rng := rangeOf(c.Loc)
		switch c.Type {
		case ast.ATDropColumn:
			if b.dropColumnKeyFails(r, c.Name, rng, true) || b.dropColumnFKFails(r, c.Name, rng) || b.dropColumnKeyFails(r, c.Name, rng, false) {
				return true
			}
		case ast.ATModifyColumn, ast.ATChangeColumn:
			col := c.Name
			if c.Type == ast.ATModifyColumn && c.Column != nil {
				col = c.Column.Name
			}
			if c.Column == nil || col == "" {
				continue
			}
			if b.typeChangeFKFails(r, col, c.Column, rng) || b.rewriteFails(r, col, c.Column, rng) {
				return true
			}
		case ast.ATDropIndex:
			if dropKindOf(b.w.sql, c) == dropIndex && b.dropIndexFKFails(r, c.Name, rng) {
				return true
			}
		case ast.ATAddIndex, ast.ATAddConstraint:
			if c.Constraint == nil {
				continue
			}
			if b.addFKTypesFail(r, c.Constraint, func(col string) (columnType, bool) {
				if definesColumn(v, col) {
					return columnType{}, false
				}
				return b.typeOf(r.db, r.name, col)
			}, rng) {
				return true
			}
		case ast.ATTableOption:
			if col := ttlColumn(c.Option); col != "" && !definesColumn(v, col) {
				if t, ok := b.typeOf(r.db, r.name, col); ok && b.ttlTypeFails(col, t, rng) {
					return true
				}
			}
		case ast.ATPartitionBy:
			if len(v.Commands) == 1 && b.alterPartitionKeysFail(r, c.PartitionBy, rng) {
				return true
			}
		}
	}
	return false
}

// definesColumn reports whether an ALTER TABLE adds, modifies, or renames
// a column of the name, whose type the walk then does not read from the
// table.
func definesColumn(v *ast.AlterTableStmt, col string) bool {
	for _, c := range v.Commands {
		cols := c.Columns
		if c.Column != nil {
			cols = append([]*ast.ColumnDef{c.Column}, cols...)
		}
		for _, d := range cols {
			if d != nil && strings.EqualFold(d.Name, col) {
				return true
			}
		}
		if (c.Type == ast.ATRenameColumn || c.Type == ast.ATChangeColumn) && (strings.EqualFold(c.Name, col) || strings.EqualFold(c.NewName, col)) {
			return true
		}
	}
	return false
}

// dropColumnKeyFails fails DROP COLUMN of a column of an index of more
// than one key part, or of the primary key: TiDB drops a column with the
// one-column indexes on it, and no other (8200, "can't drop column with
// composite index covered or Primary Key covered", or "Unsupported drop
// integer primary key"; verified, the same statement dropping the index
// too). TiDB checks an index of several key parts before a foreign key
// (verified), and a one-column primary key after one, or after a
// partitioning function that reads the column (3855): multi picks which
// the call checks.
func (b *builder) dropColumnKeyFails(r tableRef, col string, rng review.Range, multi bool) bool {
	k, ok := b.key(spaceColumn, col)
	if !ok {
		return false
	}
	ix, ok := b.indexesOf(r)
	if !ok {
		return false
	}
	for _, i := range ix {
		if !i.has(k) || multi != (len(i.cols) > 1) {
			continue
		}
		switch {
		case i.primary:
			b.fail(rng, nil, review.UnsupportedForm, engineName+" can't drop column "+review.Code(ident(col))+", which is in the primary key of "+review.Code(r.written)+".")
			return true
		case multi:
			b.fail(rng, nil, review.UnsupportedForm, engineName+" can't drop column "+review.Code(ident(col))+", which index "+review.Code(ident(i.name))+" holds with other columns.")
			return true
		}
	}
	return false
}

// rewriteFails fails MODIFY or CHANGE that makes TiDB rewrite a column's
// data (rewritesData) where it can't: a column of the primary key, clustered
// or not ("this column has primary key flag"), or any column of a
// partitioned table ("table is partition table"), both 8200 (verified).
func (b *builder) rewriteFails(r tableRef, col string, def *ast.ColumnDef, rng review.Range) bool {
	k, ok := b.key(spaceColumn, col)
	if !ok {
		return false
	}
	from, ok := b.typeOf(r.db, r.name, col)
	if !ok {
		return false
	}
	to, ok := parseColumnType(astColumnType(def.TypeName))
	if !ok || !rewritesData(from, to) {
		return false
	}
	if ix, ok := b.indexesOf(r); ok {
		for _, i := range ix {
			if i.primary && i.has(k) {
				b.fail(rng, nil, review.UnsupportedForm, engineName+" can't change primary key column "+review.Code(ident(col))+" from "+review.Code(from.String())+" to "+review.Code(to.String())+", which rewrites its data.")
				return true
			}
		}
	}
	if kind := b.partitioning(r); kind != "" && kind != partNone {
		b.fail(rng, nil, review.UnsupportedForm, engineName+" can't change column "+review.Code(ident(col))+" from "+review.Code(from.String())+" to "+review.Code(to.String())+", which rewrites its data, in partitioned table "+review.Code(r.written)+".")
		return true
	}
	return false
}

// ttlTypes are the types a TTL column may have (8148 otherwise, verified).
var ttlTypes = map[string]bool{"date": true, "datetime": true, "timestamp": true}

func (b *builder) ttlTypeFails(col string, t columnType, rng review.Range) bool {
	if ttlTypes[t.base] {
		return false
	}
	b.fail(rng, nil, review.UnsupportedForm, "TTL column "+review.Code(ident(col))+" is "+review.Code(t.String())+"; "+engineName+" takes only a "+review.Code("DATETIME")+", "+review.Code("DATE")+", or "+review.Code("TIMESTAMP")+" column.")
	return true
}

// idCacheFails fails AUTO_ID_CACHE that moves a table between 1 and any
// other value (1105, "Can't Alter AUTO_ID_CACHE between 1 and non-1",
// verified both ways). A statement that sets it twice is left alone.
func (b *builder) idCacheFails(v *ast.AlterTableStmt, r tableRef) bool {
	var opt *ast.AlterTableCmd
	for _, c := range v.Commands {
		if c.Type == ast.ATTableOption && c.Option != nil && strings.EqualFold(c.Option.Name, "AUTO_ID_CACHE") {
			if opt != nil {
				return false
			}
			opt = c
		}
	}
	if opt == nil {
		return false
	}
	to := idCacheKind(opt.Option.Value)
	f, from := b.l.Lookup(idCacheKey(r.db, r.name))
	if to == "" || f != walk.Present || (to == idCacheOne) == (from == idCacheOne) {
		return false
	}
	b.fail(rangeOf(opt.Loc), nil, review.UnsupportedForm, engineName+" can't move "+review.Code("AUTO_ID_CACHE")+" of "+review.Code(r.written)+" between 1 and another value.")
	return true
}

// partitionColumns returns the columns a PARTITION BY clause reads, false
// when the walk does not read them: KEY with no column list, and a clause
// the walk does not follow (clausePartitioning).
func partitionColumns(pc *ast.PartitionClause) ([]string, bool) {
	if pc == nil {
		return nil, false
	}
	if _, ok := clausePartitioning(pc); !ok {
		return nil, false
	}
	var cols []string
	if len(pc.Columns) > 0 {
		cols = append(cols, pc.Columns...)
	} else if pc.Expr != nil {
		ast.Inspect(pc.Expr, func(n ast.Node) bool {
			if c, ok := n.(*ast.ColumnRef); ok && !c.Star && c.Column != "" {
				cols = append(cols, c.Column)
			}
			return true
		})
	}
	return cols, len(cols) > 0
}

// uncoveringPartition reports whether a key's columns miss a column of the
// partitioning.
func uncoveringPartition(keyCols, partCols []string) bool {
	for _, p := range partCols {
		found := false
		for _, k := range keyCols {
			if strings.EqualFold(k, p) {
				found = true
			}
		}
		if !found {
			return true
		}
	}
	return false
}

func (b *builder) partitionKeyFail(r tableRef, key string, rng review.Range) {
	b.fail(rng, nil, review.UnsupportedForm, "Every unique key of "+review.Code(r.written)+" must hold all columns of its partitioning, and "+key+" doesn't.")
}

// createPartitionKeysFail fails a CREATE TABLE whose partitioning reads a
// column one of its unique keys does not hold: 1503 for a clustered
// primary key, 8264 for any other (verified). A GLOBAL key, and one with a
// functional key part, is left alone, and so is a partitioning that reads
// a column the table does not define (1054 first).
func (b *builder) createPartitionKeysFail(v *ast.CreateTableStmt, r tableRef) bool {
	if v.Partitions == nil || v.Partitions.Interval != nil {
		return false
	}
	parts, ok := partitionColumns(v.Partitions)
	if !ok {
		return false
	}
	for _, p := range parts {
		found := false
		for _, c := range v.Columns {
			if strings.EqualFold(c.Name, p) {
				found = true
			}
		}
		if !found {
			return false
		}
	}
	rng := rangeOf(v.Partitions.Loc)
	for _, c := range v.Columns {
		for _, cc := range c.Constraints {
			if (cc.Type == ast.ColConstrPrimaryKey || cc.Type == ast.ColConstrUnique) && uncoveringPartition([]string{c.Name}, parts) {
				what := "the primary key"
				if cc.Type == ast.ColConstrUnique {
					what = "the unique key on " + review.Code(ident(c.Name))
				}
				b.partitionKeyFail(r, what, rng)
				return true
			}
		}
	}
	for _, c := range v.Constraints {
		if c.Type != ast.ConstrPrimaryKey && c.Type != ast.ConstrUnique {
			continue
		}
		cols, ok := plainKeyParts(c)
		if !ok || globalKey(c) {
			continue
		}
		if uncoveringPartition(cols, parts) {
			b.partitionKeyFail(r, keyWords(c), rng)
			return true
		}
	}
	return false
}

// alterPartitionKeysFail fails ALTER TABLE ... PARTITION BY, alone in its
// statement, of a table that is not partitioned, whose partitioning reads a
// column a unique key of the table does not hold (1503 or 8264, verified).
// An UPDATE INDEXES clause is left alone, as is a key with a functional key
// part.
func (b *builder) alterPartitionKeysFail(r tableRef, pc *ast.PartitionClause, rng review.Range) bool {
	// TiDB checks the unique keys after the partitions: an INTERVAL
	// clause fails first on its own (8200 or 1054, verified).
	if pc == nil || len(pc.UpdateIndexes) > 0 || pc.Interval != nil {
		return false
	}
	// A partitioned table may hold a GLOBAL key; one that is not
	// partitioned holds none.
	if b.partitioning(r) != partNone {
		return false
	}
	parts, ok := partitionColumns(pc)
	if !ok {
		return false
	}
	var keys []string
	for _, p := range parts {
		k, ok := b.key(spaceColumn, p)
		if !ok {
			return false
		}
		if f, _ := b.l.Lookup(colKey(r.db, r.name, k)); f != walk.Present {
			return false
		}
		keys = append(keys, k)
	}
	ix, ok := b.indexesOf(r)
	if !ok {
		return false
	}
	for _, i := range ix {
		if !i.unique || i.global || contains(i.cols, "") {
			continue
		}
		if uncoveringPartition(i.cols, keys) {
			what := "the primary key"
			if !i.primary {
				what = "index " + review.Code(ident(i.name))
			}
			b.partitionKeyFail(r, what, rng)
			return true
		}
	}
	return false
}

// plainKeyParts returns a key's columns, false when one is a functional
// key part.
func plainKeyParts(c *ast.Constraint) ([]string, bool) {
	if len(c.IndexColumns) == 0 {
		return c.Columns, len(c.Columns) > 0
	}
	var cols []string
	for _, ic := range c.IndexColumns {
		cr, ok := ic.Expr.(*ast.ColumnRef)
		if ic.Functional || !ok || cr.Star || cr.Column == "" {
			return nil, false
		}
		cols = append(cols, cr.Column)
	}
	return cols, true
}

// globalKey reports whether a key's options make it GLOBAL: the last of
// GLOBAL and LOCAL wins.
func globalKey(c *ast.Constraint) bool {
	global := false
	for _, o := range c.IndexOptions {
		switch strings.ToUpper(o.Name) {
		case "GLOBAL":
			global = true
		case "LOCAL":
			global = false
		}
	}
	return global
}

// keyWords names a table constraint in a finding.
func keyWords(c *ast.Constraint) string {
	if c.Type == ast.ConstrPrimaryKey {
		return "the primary key"
	}
	if c.Name != "" {
		return "unique key " + review.Code(ident(c.Name))
	}
	return "a unique key"
}

// algorithmFails fails ALTER TABLE ... ALGORITHM = INSTANT beside a
// subcommand that adds an index or a constraint of any kind, which TiDB
// runs only in place (1846, verified for INDEX, UNIQUE, PRIMARY KEY,
// FOREIGN KEY, CHECK, and FULLTEXT). CREATE INDEX ... ALGORITHM = INSTANT
// passes (verified), and so does every other subcommand.
func (b *builder) algorithmFails(v *ast.AlterTableStmt) bool {
	instant := false
	algorithms := 0
	var add *ast.AlterTableCmd
	for _, c := range v.Commands {
		switch c.Type {
		case ast.ATAlgorithm:
			algorithms++
			instant = strings.EqualFold(c.Name, "INSTANT")
		case ast.ATAddIndex, ast.ATAddConstraint:
			if add == nil {
				add = c
			}
		}
	}
	if !instant || algorithms != 1 || add == nil {
		return false
	}
	b.fail(rangeOf(add.Loc), nil, review.UnsupportedForm, engineName+" can't add an index or a constraint with "+review.Code("ALGORITHM=INSTANT")+"; it takes "+review.Code("INPLACE")+".")
	return true
}

// collationCharsets maps each collation TiDB v8.5.5 lists (SHOW
// COLLATION) to its character set; utf8mb3's names are utf8's.
var collationCharsets = map[string]string{
	"ascii_bin": "ascii", "binary": "binary", "gbk_bin": "gbk", "gbk_chinese_ci": "gbk", "latin1_bin": "latin1",
	"utf8_bin": "utf8", "utf8_general_ci": "utf8", "utf8_unicode_ci": "utf8",
	"utf8mb3_bin": "utf8", "utf8mb3_general_ci": "utf8", "utf8mb3_unicode_ci": "utf8",
	"utf8mb4_0900_ai_ci": "utf8mb4", "utf8mb4_0900_bin": "utf8mb4", "utf8mb4_bin": "utf8mb4",
	"utf8mb4_general_ci": "utf8mb4", "utf8mb4_unicode_ci": "utf8mb4",
}

// collationMismatch returns the text of a finding for a collation that is
// no collation of a character set, "" when it is or the walk does not
// know either (1253, verified).
func collationMismatch(charset, collation string) string {
	cs := normalCharset(charset)
	if cs == "" || !tidbSupportedCharsets[cs] && cs != "utf8" {
		return ""
	}
	of, ok := collationCharsets[strings.ToLower(strings.TrimSpace(collation))]
	if !ok || of == cs {
		return ""
	}
	return "Collation " + review.Code(collation) + " isn't one of character set " + review.Code(charset) + "."
}

// optionsCollationMismatch reads a table's CHARSET and COLLATE options, in
// that order and one of each, as TiDB pairs them; a COLLATE first is 1302
// (verified), which the walk leaves alone.
func optionsCollationMismatch(opts []*ast.TableOption) string {
	charset, collation := "", ""
	charsets, collations := 0, 0
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		switch strings.ToLower(opt.Name) {
		case "charset", "character set", "default charset", "default character set":
			charsets++
			charset = opt.Value
		case "collate", "default collate":
			collations++
			if charsets == 0 {
				return ""
			}
			collation = opt.Value
		}
	}
	if charsets != 1 || collations != 1 {
		return ""
	}
	return collationMismatch(charset, collation)
}

// createCollationFails fails a CREATE TABLE whose column or table names a
// collation of another character set than it names (1253).
func (b *builder) createCollationFails(v *ast.CreateTableStmt) bool {
	for _, c := range v.Columns {
		if c.TypeName == nil || c.TypeName.Charset == "" || c.TypeName.Collate == "" {
			continue
		}
		if msg := collationMismatch(c.TypeName.Charset, c.TypeName.Collate); msg != "" {
			b.fail(rangeOf(c.Loc), nil, review.UnsupportedForm, msg)
			return true
		}
	}
	if msg := optionsCollationMismatch(v.Options); msg != "" {
		b.fail(b.anchor(), nil, review.UnsupportedForm, msg)
		return true
	}
	return false
}

// alterCollationFails fails an ALTER TABLE whose CONVERT TO, or whose
// CHARSET and COLLATE options, name a collation of another character set
// (1253). A column definition's pair is judged too.
func (b *builder) alterCollationFails(v *ast.AlterTableStmt) bool {
	var opts []*ast.TableOption
	for _, c := range v.Commands {
		switch c.Type {
		case ast.ATConvertCharset:
			if c.Name != "" && c.NewName != "" {
				if msg := collationMismatch(c.Name, c.NewName); msg != "" {
					b.fail(rangeOf(c.Loc), nil, review.UnsupportedForm, msg)
					return true
				}
			}
		case ast.ATTableOption:
			opts = append(opts, c.Option)
		}
		cols := c.Columns
		if c.Column != nil {
			cols = append([]*ast.ColumnDef{c.Column}, cols...)
		}
		for _, d := range cols {
			if d == nil || d.TypeName == nil || d.TypeName.Charset == "" || d.TypeName.Collate == "" {
				continue
			}
			if msg := collationMismatch(d.TypeName.Charset, d.TypeName.Collate); msg != "" {
				b.fail(rangeOf(d.Loc), nil, review.UnsupportedForm, msg)
				return true
			}
		}
	}
	if msg := optionsCollationMismatch(opts); msg != "" {
		b.fail(b.anchor(), nil, review.UnsupportedForm, msg)
		return true
	}
	return false
}

// nondeterministic are functions TiDB refuses in a generated column (3102)
// and in an expression index (3758), from v8.5.5's
// IllegalFunctions4GeneratedColumns; each was verified in one of them.
var nondeterministic = map[string]bool{
	"rand": true, "uuid": true, "uuid_short": true, "now": true, "current_timestamp": true,
	"sysdate": true, "connection_id": true, "unix_timestamp": true, "curdate": true,
	"current_date": true, "curtime": true, "current_time": true, "localtime": true,
	"localtimestamp": true, "utc_date": true, "utc_time": true, "utc_timestamp": true,
	"last_insert_id": true, "found_rows": true, "row_count": true, "sleep": true,
	"current_user": true, "user": true, "session_user": true, "system_user": true,
	"database": true, "schema": true, "version": true,
}

// nondeterministicCall returns the first function of nondeterministic an
// expression calls, "" for none. A name written quoted is a column, or a
// call the walk doesn't judge.
func nondeterministicCall(sql string, e ast.Node) string {
	found := ""
	if e == nil {
		return ""
	}
	ast.Inspect(e, func(n ast.Node) bool {
		if found != "" {
			return false
		}
		switch x := n.(type) {
		case *ast.SelectStmt:
			return false
		case *ast.FuncCallExpr:
			if x.Schema == "" && nondeterministic[strings.ToLower(x.Name)] && !quotedAt(sql, x.Loc.Start) {
				found = strings.ToUpper(x.Name)
			}
		case *ast.ColumnRef:
			if x.Table == "" && x.Schema == "" && niladicNames[strings.ToLower(x.Column)] && nondeterministic[strings.ToLower(x.Column)] && !quotedAt(sql, x.Loc.Start) {
				found = strings.ToUpper(x.Column)
			}
		}
		return found == ""
	})
	return found
}

// generatedFuncFails fails a generated column that calls a function
// nondeterministic lists (3102).
func (b *builder) generatedFuncFails(cols []*ast.ColumnDef) bool {
	for _, c := range cols {
		if c == nil || c.Generated == nil {
			continue
		}
		if fn := nondeterministicCall(b.w.sql, c.Generated.Expr); fn != "" {
			b.fail(rangeOf(c.Loc), nil, review.UnsupportedForm, "Generated column "+review.Code(ident(c.Name))+" can't call "+review.Code(fn)+", whose value isn't deterministic.")
			return true
		}
	}
	return false
}

// indexFuncFails fails an expression index whose key part calls a
// function nondeterministic lists (3758).
func (b *builder) indexFuncFails(parts []*ast.IndexColumn, rng review.Range) bool {
	for _, ic := range parts {
		if ic == nil || !ic.Functional {
			continue
		}
		if fn := nondeterministicCall(b.w.sql, ic.Expr); fn != "" {
			b.fail(rng, nil, review.UnsupportedForm, "An expression index can't call "+review.Code(fn)+", whose value isn't deterministic.")
			return true
		}
	}
	return false
}

// alterTextShapeFails fails an ALTER TABLE its text refutes by the rules
// of this file: ALGORITHM = INSTANT beside an added key, a collation of
// another character set, and a nondeterministic generated column or
// expression index.
func (b *builder) alterTextShapeFails(v *ast.AlterTableStmt) bool {
	if b.algorithmFails(v) || b.alterCollationFails(v) {
		return true
	}
	for _, c := range v.Commands {
		cols := c.Columns
		if c.Column != nil {
			cols = append([]*ast.ColumnDef{c.Column}, cols...)
		}
		if b.generatedFuncFails(cols) {
			return true
		}
		if (c.Type == ast.ATAddIndex || c.Type == ast.ATAddConstraint) && c.Constraint != nil && b.indexFuncFails(c.Constraint.IndexColumns, rangeOf(c.Loc)) {
			return true
		}
	}
	return false
}

// createShapeFails fails a CREATE TABLE with a column list by the rules of
// this file and the foreign key checks: a collation of another character
// set, a nondeterministic generated column or expression index, a unique
// key the partitioning reads past, a TTL column of a type TiDB doesn't
// take, and a foreign key whose types don't match.
func (b *builder) createShapeFails(v *ast.CreateTableStmt, r tableRef) bool {
	if v.Like != nil || v.Select != nil {
		return false
	}
	if b.createCollationFails(v) || b.generatedFuncFails(v.Columns) {
		return true
	}
	for _, c := range v.Constraints {
		if b.indexFuncFails(c.IndexColumns, rangeOf(c.Loc)) {
			return true
		}
	}
	if b.createPartitionKeysFail(v, r) {
		return true
	}
	types := make(map[string]columnType)
	for _, c := range v.Columns {
		if k, ok := b.key(spaceColumn, c.Name); ok {
			if t, ok := parseColumnType(astColumnType(c.TypeName)); ok {
				types[k] = t
			}
		}
	}
	for _, opt := range v.Options {
		if col := ttlColumn(opt); col != "" {
			if k, ok := b.key(spaceColumn, col); ok {
				if t, ok := types[k]; ok && b.ttlTypeFails(col, t, rangeOf(opt.Loc)) {
					return true
				}
			}
		}
	}
	for _, c := range v.Constraints {
		if b.addFKTypesFail(r, c, func(col string) (columnType, bool) {
			k, ok := b.key(spaceColumn, col)
			if !ok {
				return columnType{}, false
			}
			t, ok := types[k]
			return t, ok
		}, rangeOf(c.Loc)) {
			return true
		}
	}
	return false
}
