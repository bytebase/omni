package review

import (
	"sort"
	"strconv"
	"strings"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
)

// disallowedFuncs are the builtins whose value depends on the session or
// the moment, which no generated column (3763), CHECK constraint (3814),
// or functional index (3758) may call (each verified on mysql:8.0.46,
// RAND with a seed too). UNIX_TIMESTAMP counts only with no argument.
var disallowedFuncs = map[string]bool{
	"rand": true, "uuid": true, "uuid_short": true, "now": true, "current_timestamp": true,
	"localtime": true, "localtimestamp": true, "sysdate": true, "curdate": true, "current_date": true,
	"curtime": true, "current_time": true, "utc_date": true, "utc_time": true, "utc_timestamp": true,
	"connection_id": true, "current_user": true, "user": true, "session_user": true, "system_user": true,
	"database": true, "schema": true, "last_insert_id": true, "found_rows": true, "row_count": true,
	"sleep": true, "get_lock": true, "unix_timestamp": true,
}

// disallowedFunc returns the first builtin of disallowedFuncs an expression
// calls, outside any nested query, "" for none.
func disallowedFunc(expr ast.ExprNode) string {
	found := ""
	if expr == nil {
		return ""
	}
	ast.Inspect(expr, func(n ast.Node) bool {
		if found != "" {
			return false
		}
		switch x := n.(type) {
		case *ast.SubqueryExpr, *ast.ExistsExpr, *ast.SelectStmt:
			return false
		case *ast.FuncCallExpr:
			name := strings.ToLower(x.Name)
			if x.Schema == "" && disallowedFuncs[name] && (name != "unix_timestamp" || len(x.Args) == 0) {
				found = name
			}
		}
		return found == ""
	})
	return found
}

// callsDisallowed is the text for an expression of a kind that can't call
// a builtin of disallowedFuncs.
func callsDisallowed(what, fn string) string {
	return what + " can't call " + review.Code(strings.ToUpper(fn)+"()") + ", whose value isn't deterministic."
}

// charsetMaxBytes are the bytes a character takes at most, for the
// character sets the walk sizes a column of.
var charsetMaxBytes = map[string]int{
	"utf8mb4": 4, "utf8mb3": 3, "utf16": 4, "utf16le": 4, "utf32": 4, "ucs2": 2, "gb18030": 4,
	"latin1": 1, "ascii": 1, "binary": 1, "latin2": 1, "latin5": 1, "latin7": 1, "cp1250": 1,
	"cp1251": 1, "cp1256": 1, "cp1257": 1, "gbk": 2, "big5": 2, "sjis": 2, "cp932": 2,
	"euckr": 2, "gb2312": 2, "ujis": 3, "eucjpms": 3,
}

// columnTypeFault returns the fault of a column type MySQL refuses under
// the default sql_mode, which is strict (each verified on mysql:8.0.46):
//
//   - a VARCHAR or VARBINARY longer than 65,535 bytes, its character set's
//     widest character counted: the column's, or charset when it has none,
//     "" when not known (1074);
//   - an ENUM or SET naming a value twice, spelled alike (1291).
func columnTypeFault(col *ast.ColumnDef, charset string) string {
	dt := col.TypeName
	if dt == nil {
		return ""
	}
	switch base := astBaseType(dt); base {
	case "varchar", "varbinary":
		if !dt.HasLength {
			return ""
		}
		width := 1
		if base == "varchar" && !dt.Binary {
			cs := charset
			switch {
			case dt.Charset != "":
				cs = normalCharset(dt.Charset)
			case dt.Collate != "":
				// A column's collation names its character set.
				cs = collationCharset(dt.Collate)
			}
			if w, ok := charsetMaxBytes[cs]; ok {
				width = w
			}
		}
		if dt.Length > 65535/width {
			return "Column " + review.Code(ident(col.Name)) + " is longer than " + strconv.Itoa(65535/width) + " characters, " + engineName + "'s limit for its type; use " + review.Code("TEXT") + " or " + review.Code("BLOB") + "."
		}
	case "enum", "set":
		seen := make(map[string]bool, len(dt.EnumValues))
		for _, v := range dt.EnumValues {
			if seen[v] {
				return "Column " + review.Code(ident(col.Name)) + " lists " + review.Code("'"+v+"'") + " twice in its " + review.Code(strings.ToUpper(base)) + "."
			}
			seen[v] = true
		}
	}
	return ""
}

// columnExprFault returns the fault of a column definition's generated
// expression or CHECK constraint calling a builtin it can't (3763, 3814).
func columnExprFault(col *ast.ColumnDef) string {
	if col.Generated != nil {
		if fn := disallowedFunc(col.Generated.Expr); fn != "" {
			return callsDisallowed("Generated column "+review.Code(ident(col.Name)), fn)
		}
	}
	for _, cc := range col.Constraints {
		if cc != nil && cc.Type == ast.ColConstrCheck {
			if fn := disallowedFunc(cc.Expr); fn != "" {
				return callsDisallowed("A "+review.Code("CHECK")+" constraint", fn)
			}
		}
	}
	return ""
}

// keyTextFault returns the fault of a key or index a statement defines,
// from its text: a functional key part calling a builtin it can't (3758),
// a primary key with a functional part (3756), or an invisible primary key
// (3522).
func keyTextFault(c *ast.Constraint) string {
	if c == nil {
		return ""
	}
	if c.Type == ast.ConstrCheck {
		if fn := disallowedFunc(c.Expr); fn != "" {
			return callsDisallowed("A "+review.Code("CHECK")+" constraint", fn)
		}
		return ""
	}
	for _, ic := range c.IndexColumns {
		if ic == nil {
			continue
		}
		if _, _, functional := indexColPart(ic); !functional {
			continue
		}
		if fn := disallowedFunc(ic.Expr); fn != "" {
			return callsDisallowed("A functional key part", fn)
		}
		if c.Type == ast.ConstrPrimaryKey {
			return "A primary key can't have a functional key part."
		}
	}
	if c.Type == ast.ConstrPrimaryKey && isInvisibleKey(c.IndexOptions) {
		return "A primary key can't be invisible."
	}
	return ""
}

// implicitPrimaryKey returns the UNIQUE key of a CREATE TABLE with no
// primary key that MySQL makes its primary key: the first one declared
// whose parts are all NOT NULL columns (AUTO_INCREMENT is NOT NULL too),
// as mysql/catalog's checkImplicitPrimaryKeyVisible verifies. A column's own UNIQUE comes first
// where the column stands and is visible, so it returns nil, as does a
// key part the walk can't judge.
func implicitPrimaryKey(v *ast.CreateTableStmt) *ast.Constraint {
	notNull := make(map[string]bool, len(v.Columns))
	type key struct {
		pos int
		con *ast.Constraint // nil for a column's own UNIQUE
	}
	var keys []key
	for _, c := range v.Columns {
		if c == nil {
			return nil
		}
		nn, null, autoInc, unique := false, false, false, false
		for _, cc := range c.Constraints {
			switch {
			case cc == nil:
			case cc.Type == ast.ColConstrNotNull:
				nn, null = true, false
			case cc.Type == ast.ColConstrNull:
				nn, null = false, true
			case cc.Type == ast.ColConstrAutoIncrement:
				autoInc = true
			case cc.Type == ast.ColConstrPrimaryKey:
				return nil
			case cc.Type == ast.ColConstrUnique:
				unique = true
			}
		}
		if autoInc {
			if null {
				return nil // NULL AUTO_INCREMENT: not judged
			}
			nn = true
		}
		if c.TypeName != nil && strings.EqualFold(c.TypeName.Name, "SERIAL") {
			nn, unique = true, true
		}
		notNull[strings.ToLower(c.Name)] = nn
		if unique && nn {
			keys = append(keys, key{pos: c.Loc.Start})
		}
	}
	for _, c := range v.Constraints {
		switch {
		case c == nil:
		case c.Type == ast.ConstrPrimaryKey:
			return nil
		case c.Type == ast.ConstrUnique:
			keys = append(keys, key{pos: c.Loc.Start, con: c})
		}
	}
	sort.SliceStable(keys, func(i, j int) bool { return keys[i].pos < keys[j].pos })
	for _, k := range keys {
		if k.con == nil || k.pos < 0 {
			return nil
		}
		parts := constraintKeyParts(k.con)
		if len(parts) == 0 {
			return nil
		}
		nullable := false
		for _, p := range parts {
			if p.functional || p.name == "" {
				return nil
			}
			nn, ok := notNull[strings.ToLower(p.name)]
			if !ok {
				return nil
			}
			nullable = nullable || !nn
		}
		if !nullable {
			return k.con
		}
	}
	return nil
}

// isInvisibleKey reports whether a key's options leave it invisible: the
// last of VISIBLE and INVISIBLE wins.
func isInvisibleKey(opts []*ast.IndexOption) bool {
	invisible := false
	for _, o := range opts {
		switch {
		case o == nil:
		case strings.EqualFold(o.Name, "INVISIBLE"):
			invisible = true
		case strings.EqualFold(o.Name, "VISIBLE"):
			invisible = false
		}
	}
	return invisible
}

// createColumnsFail fails a CREATE TABLE whose column list MySQL refuses
// for its text alone (each verified on mysql:8.0.46): no visible column
// (4028); a generated column reading an AUTO_INCREMENT column (3109), or a
// generated column defined after it (3107); a CHECK constraint reading an
// AUTO_INCREMENT column (3818); and the faults of columnTypeFault, with
// the table's CHARSET, and keyTextFault.
func (b *builder) createColumnsFail(v *ast.CreateTableStmt, r tableRef) bool {
	if len(v.Columns) == 0 || v.Like != nil || v.Select != nil {
		return false
	}
	fail := func(rng review.Range, msg string) bool {
		b.fail(rng, nil, review.UnsupportedForm, msg)
		return true
	}
	charset := ""
	for _, o := range v.Options {
		if o != nil && optionName(o) == "CHARSET" {
			charset = normalCharset(o.Value)
		}
		if o != nil && optionName(o) == "COLLATE" && charset == "" {
			charset = collationCharset(o.Value)
		}
	}
	for _, c := range v.Columns {
		if msg := columnExprFault(c); msg != "" {
			return fail(rangeOf(c.Loc), msg)
		}
		if b.w.sqlModeSet {
			continue
		}
		if msg := columnTypeFault(c, charset); msg != "" {
			return fail(rangeOf(c.Loc), msg)
		}
	}
	visible := false
	autoInc := make(map[string]bool)
	generatedAt := make(map[string]int)
	for i, c := range v.Columns {
		invisible := false
		for _, cc := range c.Constraints {
			if cc != nil && cc.Type == ast.ColConstrInvisible {
				invisible = true
			}
		}
		visible = visible || !invisible
		if colHasAutoIncrement(c) {
			autoInc[strings.ToLower(c.Name)] = true
		}
		if c.Generated != nil {
			generatedAt[strings.ToLower(c.Name)] = i
		}
	}
	if !visible {
		return fail(r.rng, "Table "+review.Code(r.written)+" must have at least one visible column.")
	}
	if c := implicitPrimaryKey(v); c != nil && isInvisibleKey(c.IndexOptions) {
		return fail(rangeOf(c.Loc), "A primary key can't be invisible.") // 3522
	}
	if !v.Temporary {
		names, rngs := checkNames(v)
		for i, n := range names {
			if b.checkNameTakenFails(r.db, n, walk.Key{}, rngs[i]) {
				return true
			}
		}
	}
	for i, c := range v.Columns {
		if c.Generated == nil {
			continue
		}
		for _, name := range exprColumnNames(c.Generated.Expr) {
			k := strings.ToLower(name)
			if autoInc[k] {
				return fail(rangeOf(c.Loc), "Generated column "+review.Code(ident(c.Name))+" can't read "+review.Code("AUTO_INCREMENT")+" column "+review.Code(ident(name))+".")
			}
			if at, ok := generatedAt[k]; ok && at > i {
				return fail(rangeOf(c.Loc), "Generated column "+review.Code(ident(c.Name))+" can't read generated column "+review.Code(ident(name))+", defined after it.")
			}
		}
	}
	checkReadsAutoInc := func(expr ast.ExprNode) string {
		for _, name := range exprColumnNames(expr) {
			if autoInc[strings.ToLower(name)] {
				return name
			}
		}
		return ""
	}
	for _, c := range v.Columns {
		for _, cc := range c.Constraints {
			if cc != nil && cc.Type == ast.ColConstrCheck {
				if name := checkReadsAutoInc(cc.Expr); name != "" {
					return fail(rangeOf(c.Loc), "A "+review.Code("CHECK")+" constraint can't read "+review.Code("AUTO_INCREMENT")+" column "+review.Code(ident(name))+".")
				}
			}
		}
	}
	for _, c := range v.Constraints {
		if c == nil {
			continue
		}
		if c.Type == ast.ConstrCheck {
			if name := checkReadsAutoInc(c.Expr); name != "" {
				return fail(rangeOf(c.Loc), "A "+review.Code("CHECK")+" constraint can't read "+review.Code("AUTO_INCREMENT")+" column "+review.Code(ident(name))+".")
			}
		}
		if msg := keyTextFault(c); msg != "" {
			return fail(rangeOf(c.Loc), msg)
		}
	}
	return false
}

// alterKeysFail fails an ALTER TABLE adding a key or CHECK constraint
// keyTextFault refuses, or a CHECK constraint of a name the database has
// (3822), or making the primary key invisible (3522).
func (b *builder) alterKeysFail(r tableRef, v *ast.AlterTableStmt) bool {
	// A CHECK the statement drops frees its name for one it adds.
	dropped := make(map[string]bool)
	for _, c := range v.Commands {
		if c != nil && c.Type == ast.ATDropConstraint {
			if k, ok := b.key(spaceCheckConstraint, c.Name); ok {
				dropped[k] = true
			}
		}
	}
	for _, c := range v.Commands {
		if c == nil {
			continue
		}
		switch c.Type {
		case ast.ATAddConstraint, ast.ATAddIndex:
			if msg := keyTextFault(c.Constraint); msg != "" {
				b.fail(rangeOf(c.Loc), nil, review.UnsupportedForm, msg)
				return true
			}
			if c.Constraint != nil && c.Constraint.Type == ast.ConstrCheck {
				var freed walk.Key
				if k, ok := b.key(spaceCheckConstraint, c.Constraint.Name); ok && dropped[k] {
					freed = chkKey(r.db, r.name, k)
				}
				if b.checkNameTakenFails(r.db, c.Constraint.Name, freed, rangeOf(c.Loc)) {
					return true
				}
			}
		case ast.ATAlterIndexInvisible:
			if strings.EqualFold(c.Name, "PRIMARY") {
				b.fail(rangeOf(c.Loc), nil, review.UnsupportedForm, "A primary key can't be invisible.")
				return true
			}
		}
	}
	return false
}

// exprColumnsFail fails an ALTER TABLE whose generated column, CHECK
// constraint, or functional key part reads a column the table won't have
// (1054).
func (b *builder) exprColumnsFail(r tableRef, v *ast.AlterTableStmt, d *alterDelta) bool {
	check := func(expr ast.ExprNode, rng review.Range) bool {
		for _, name := range exprColumnNames(expr) {
			if outsideASCII(name) {
				continue
			}
			if f, _, _ := b.alterFinalCol(r, d, name); f == walk.Absent {
				b.fail(rng, nil, review.ObjectNotFound, notFound("column", ident(name), r.written))
				return true
			}
		}
		return false
	}
	for _, c := range v.Commands {
		if c == nil {
			continue
		}
		rng := rangeOf(c.Loc)
		switch c.Type {
		case ast.ATAddColumn, ast.ATModifyColumn, ast.ATChangeColumn:
			cols := c.Columns
			if c.Column != nil {
				cols = append([]*ast.ColumnDef{c.Column}, cols...)
			}
			for _, col := range cols {
				if col.Generated != nil && check(col.Generated.Expr, rng) {
					return true
				}
				// The column's own faults, where MySQL meets them: an
				// added column that isn't there, or one changed that is.
				k, ok := b.key(spaceColumn, col.Name)
				if c.Type != ast.ATAddColumn {
					k, ok = b.key(spaceColumn, c.Name)
				}
				if !ok {
					continue
				}
				f, _ := b.l.Lookup(colKey(r.db, r.name, k))
				if (c.Type == ast.ATAddColumn) != (f == walk.Absent) || f == walk.Unknown {
					continue
				}
				msg := columnExprFault(col)
				if msg == "" && !b.w.sqlModeSet {
					msg = columnTypeFault(col, "")
				}
				if msg != "" {
					b.fail(rangeOf(col.Loc), nil, review.UnsupportedForm, msg)
					return true
				}
			}
		case ast.ATAddConstraint, ast.ATAddIndex:
			if c.Constraint == nil {
				continue
			}
			if c.Constraint.Type == ast.ConstrCheck {
				if check(c.Constraint.Expr, rng) {
					return true
				}
				continue
			}
			for _, ic := range c.Constraint.IndexColumns {
				// A bare column in parentheses is 3762, which MySQL
				// meets first.
				if _, _, functional := indexColPart(ic); functional && !isBareColumnRef(ic.Expr) && check(ic.Expr, rng) {
					return true
				}
			}
		}
	}
	return false
}

// Generated column kinds, as the walk records a synced column: "" for a
// column that isn't generated.
const (
	genVirtual = "virtual"
	genStored  = "stored"
)

// generatedChangeFails fails an ALTER TABLE, the table's first in the
// change, that turns a virtual generated column into a stored or plain
// one, or a plain column into a virtual one (3106), or sets a generated
// column's default (1221), each verified on mysql:8.0.46.
func (b *builder) generatedChangeFails(r tableRef, v *ast.AlterTableStmt) bool {
	kind := func(name string) (string, walk.Key, bool) {
		k, ok := b.key(spaceColumn, name)
		if !ok {
			return "", walk.Key{}, false
		}
		ck := colKey(r.db, r.name, k)
		g, known := b.w.colGen[ck]
		if !known {
			return "", ck, false
		}
		if f, _ := b.l.Lookup(ck); f != walk.Present {
			return "", ck, false
		}
		return g, ck, true
	}
	for _, c := range v.Commands {
		if c == nil {
			continue
		}
		switch c.Type {
		case ast.ATModifyColumn, ast.ATChangeColumn:
			if c.Column == nil {
				continue
			}
			from := c.Column.Name
			if c.Type == ast.ATChangeColumn {
				from = c.Name
			}
			old, ck, ok := kind(from)
			if !ok {
				continue
			}
			now := ""
			if c.Column.Generated != nil {
				now = genVirtual
				if c.Column.Generated.Stored {
					now = genStored
				}
			}
			if old == genVirtual && now != genVirtual || old == "" && now == genVirtual {
				b.failReferenced(rangeOf(c.Loc), []walk.Key{ck}, "Column "+review.Code(ident(from))+" can't change between a virtual generated column and a stored or plain one.")
				return true
			}
		case ast.ATAlterColumnDefault:
			if c.DefaultExpr == nil {
				continue
			}
			if g, ck, ok := kind(c.Name); ok && g != "" {
				b.failReferenced(rangeOf(c.Loc), []walk.Key{ck}, "Generated column "+review.Code(ident(c.Name))+" can't take a default.")
				return true
			}
		}
	}
	return false
}

// forgetGenerated records that a statement may change what the walk knows
// of a table's generated columns.
func (b *builder) forgetGenerated(r tableRef) {
	parent := tablePath(r.db, r.name)
	b.Session(func() {
		for ck := range b.w.colGen {
			if ck.Parent == parent {
				delete(b.w.colGen, ck)
			}
		}
	})
}

// checkNameTakenFails fails a statement adding a CHECK constraint whose
// name a CHECK constraint the sync lists in the database already has, on
// any table: the names are the database's (3822, verified on mysql:8.0.46).
func (b *builder) checkNameTakenFails(db, name string, freed walk.Key, rng review.Range) bool {
	if name == "" || db != b.w.snap.db {
		return false
	}
	k, ok := b.key(spaceCheckConstraint, name)
	if !ok {
		return false
	}
	tables := make([]string, 0, len(b.w.snap.checkConstraints))
	for t, names := range b.w.snap.checkConstraints {
		if names[k] {
			tables = append(tables, t)
		}
	}
	sort.Strings(tables)
	for _, t := range tables {
		ck := chkKey(db, t, k)
		if ck == freed {
			continue
		}
		if f, _ := b.l.Lookup(ck); f == walk.Present {
			b.failReferenced(rng, []walk.Key{ck}, "Check constraint "+review.Code(ident(name))+" already exists in "+review.Code(t)+".")
			return true
		}
	}
	return false
}

// checkNames returns the CHECK constraints a CREATE TABLE names, with
// where each is written.
func checkNames(v *ast.CreateTableStmt) (names []string, rngs []review.Range) {
	for _, c := range v.Columns {
		for _, cc := range c.Constraints {
			if cc != nil && cc.Type == ast.ColConstrCheck && cc.Name != "" {
				names, rngs = append(names, cc.Name), append(rngs, rangeOf(cc.Loc))
			}
		}
	}
	for _, c := range v.Constraints {
		if c != nil && c.Type == ast.ConstrCheck && c.Name != "" {
			names, rngs = append(names, c.Name), append(rngs, rangeOf(c.Loc))
		}
	}
	return names, rngs
}
