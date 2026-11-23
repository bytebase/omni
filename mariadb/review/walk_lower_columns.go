package review

import (
	"maps"
	"slices"
	"strconv"
	"strings"

	"github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
)

// nondeterministic lists the built-in functions MariaDB refuses in a CHECK
// constraint and in a STORED generated column (1901, each verified on
// mariadb:11.8.8): their value depends on the session, the time, or
// chance. UNIX_TIMESTAMP is one only without an argument.
var nondeterministic = map[string]bool{
	"rand": true, "uuid": true, "uuid_short": true, "sys_guid": true, "now": true, "current_timestamp": true,
	"localtime": true, "localtimestamp": true, "sysdate": true, "curdate": true, "current_date": true,
	"curtime": true, "current_time": true, "utc_timestamp": true, "utc_date": true, "utc_time": true,
	"user": true, "current_user": true, "session_user": true, "system_user": true, "connection_id": true,
	"last_insert_id": true, "database": true, "schema": true, "version": true, "found_rows": true,
	"row_count": true,
}

// nondeterministicCall returns the first call of an expression whose value
// isn't deterministic, as nondeterministic lists them, "" for none.
func nondeterministicCall(e ast.ExprNode) string {
	found := ""
	ast.Inspect(e, func(n ast.Node) bool {
		if found != "" {
			return false
		}
		f, ok := n.(*ast.FuncCallExpr)
		if !ok || f.Schema != "" || f.NameQuoted {
			return true
		}
		name := strings.ToLower(f.Name)
		if nondeterministic[name] || name == "unix_timestamp" && len(f.Args) == 0 {
			found = strings.ToUpper(f.Name) + "()"
		}
		return found == ""
	})
	return found
}

// cannotCall is the text for a call of a function whose value isn't
// deterministic where MariaDB needs one that is (1901).
func cannotCall(where, call string) string {
	return where + " can't call " + review.Code(call) + ", whose value isn't deterministic."
}

// columnTextFails fails a column definition its text alone refutes (each
// verified on mariadb:11.8.8): a CHECK constraint or a STORED generated
// column calling a function whose value isn't deterministic (1901); an
// ENUM or SET that lists a value twice (1291), and a VARCHAR longer than
// its character set allows (1074), both under the default sql_mode, whose
// STRICT_TRANS_TABLES makes either an error. tableCharset is the character
// set the table's own options name, "" when the walk does not know it.
func (b *builder) columnTextFails(col *ast.ColumnDef, tableCharset string) bool {
	if col == nil {
		return false
	}
	rng := rangeOf(col.Loc)
	if g := col.Generated; g != nil && g.Stored && g.Expr != nil {
		if call := nondeterministicCall(g.Expr); call != "" {
			b.fail(rng, nil, review.UnsupportedForm, cannotCall("Generated column "+review.Code(ident(col.Name)), call))
			return true
		}
	}
	for _, cc := range col.Constraints {
		if cc != nil && cc.Type == ast.ColConstrCheck && cc.Expr != nil {
			if call := nondeterministicCall(cc.Expr); call != "" {
				b.fail(rng, nil, review.UnsupportedForm, cannotCall("A "+review.Code("CHECK")+" constraint", call))
				return true
			}
		}
	}
	if b.w.sqlModeSet || col.TypeName == nil {
		return false
	}
	dt := col.TypeName
	switch strings.ToLower(dt.Name) {
	case "enum", "set":
		seen := make(map[string]bool, len(dt.EnumValues))
		for _, v := range dt.EnumValues {
			if seen[v] {
				b.fail(rng, nil, review.UnsupportedForm, "Column "+review.Code(ident(col.Name))+" lists "+review.Code("'"+v+"'")+" twice in its "+review.Code(strings.ToUpper(dt.Name))+".")
				return true
			}
			seen[v] = true
		}
	case "varchar":
		cs := normalCharset(dt.Charset)
		if cs == "" && dt.Collate != "" {
			cs = collationCharset(dt.Collate)
		}
		if cs == "" && dt.Charset == "" && dt.Collate == "" {
			cs = tableCharset
		}
		limit := 65532
		if n := charsetWidth[cs]; n > 0 {
			limit /= n
		}
		if dt.Length > limit {
			b.fail(rng, nil, review.UnsupportedForm, "Column "+review.Code(ident(col.Name))+" is longer than "+strconv.Itoa(limit)+" characters, "+engineName+"'s limit for its type; use "+review.Code("TEXT")+" or "+review.Code("BLOB")+".")
			return true
		}
	}
	return false
}

// charsetWidth holds the most bytes a character takes in each Unicode
// character set, which sets how long a VARCHAR of it may be; a character
// set it lacks is read as one byte, the widest limit.
var charsetWidth = map[string]int{
	"utf8mb3": 3, "utf8mb4": 4, "ucs2": 2, "utf16": 4, "utf16le": 4, "utf32": 4,
}

// tableCharsetOf returns the character set a CREATE TABLE's options name,
// by CHARSET or by COLLATE, "" for none or DEFAULT.
func tableCharsetOf(opts []*ast.TableOption) string {
	cs := ""
	for _, o := range opts {
		if o == nil || strings.EqualFold(optionValue(o), "DEFAULT") {
			continue
		}
		switch optionName(o) {
		case "CHARSET":
			cs = normalCharset(o.Value)
		case "COLLATE":
			if c := collationCharset(o.Value); c != "" {
				cs = c
			}
		}
	}
	return cs
}

// createColumnsFail fails a CREATE TABLE whose columns and CHECK
// constraints MariaDB refuses against its own text (each verified on
// mariadb:11.8.8): columnTextFails' faults; no visible column (1113); a
// generated column reading another generated column defined after it
// (4029); a CHECK constraint calling a function whose value isn't
// deterministic, or reading an AUTO_INCREMENT column (1901); two CHECK
// constraints, or a CHECK and a unique key or foreign key, of one name
// (1826); and a primary key marked IGNORED (4174).
func (b *builder) createColumnsFail(v *ast.CreateTableStmt, r tableRef) bool {
	if v.Like != nil {
		return false
	}
	cs := tableCharsetOf(v.Options)
	visible := false
	autoInc := make(map[string]bool)
	generated := make(map[string]int)
	for i, c := range v.Columns {
		if b.columnTextFails(c, cs) {
			return true
		}
		if !colDeclares(c, ast.ColConstrInvisible) {
			visible = true
		}
		k, ok := b.key(spaceColumn, c.Name)
		if !ok {
			return false
		}
		if c.AutoIncrement || colDeclares(c, ast.ColConstrAutoIncrement) {
			autoInc[k] = true
		}
		if c.Generated != nil && c.Generated.Expr != nil {
			generated[k] = i
		}
	}
	if len(v.Columns) > 0 && !visible && v.Select == nil {
		b.fail(r.rng, nil, review.UnsupportedForm, "Table "+review.Code(r.written)+" must have at least one visible column.")
		return true
	}
	for i, c := range v.Columns {
		if c.Generated == nil || c.Generated.Expr == nil {
			continue
		}
		for _, ref := range exprBareColumnRefs(c.Generated.Expr) {
			k, ok := b.key(spaceColumn, ref.Column)
			if j, gen := generated[k]; ok && gen && j > i {
				b.fail(rangeOf(c.Loc), nil, review.UnsupportedForm, "Generated column "+review.Code(ident(c.Name))+" reads generated column "+review.Code(ident(ref.Column))+", defined after it.")
				return true
			}
		}
	}
	readsAutoInc := func(e ast.ExprNode) string {
		for _, ref := range exprBareColumnRefs(e) {
			if k, ok := b.key(spaceColumn, ref.Column); ok && autoInc[k] {
				return ref.Column
			}
		}
		return ""
	}
	for _, c := range v.Columns {
		for _, cc := range c.Constraints {
			if cc != nil && cc.Type == ast.ColConstrCheck && cc.Expr != nil {
				if col := readsAutoInc(cc.Expr); col != "" {
					b.fail(rangeOf(c.Loc), nil, review.UnsupportedForm, "A "+review.Code("CHECK")+" constraint can't read "+review.Code("AUTO_INCREMENT")+" column "+review.Code(ident(col))+".")
					return true
				}
			}
		}
	}
	checks := make(map[string]bool)
	named := make(map[string]bool)
	for _, c := range v.Constraints {
		if c == nil {
			continue
		}
		rng := rangeOf(c.Loc)
		switch c.Type {
		case ast.ConstrCheck:
			if c.Expr != nil {
				if call := nondeterministicCall(c.Expr); call != "" {
					b.fail(rng, nil, review.UnsupportedForm, cannotCall("A "+review.Code("CHECK")+" constraint", call))
					return true
				}
				if col := readsAutoInc(c.Expr); col != "" {
					b.fail(rng, nil, review.UnsupportedForm, "A "+review.Code("CHECK")+" constraint can't read "+review.Code("AUTO_INCREMENT")+" column "+review.Code(ident(col))+".")
					return true
				}
			}
			if c.Name == "" {
				continue
			}
			k, ok := b.key(spaceColumn, c.Name)
			if !ok {
				continue
			}
			if checks[k] || named[k] {
				b.fail(rng, nil, review.DuplicateName, "Check constraint "+review.Code(ident(c.Name))+" shares its name with another constraint of the table.")
				return true
			}
			checks[k] = true
		case ast.ConstrUnique, ast.ConstrForeignKey:
			if c.Name == "" {
				continue
			}
			k, ok := b.key(spaceColumn, c.Name)
			if !ok {
				continue
			}
			if checks[k] {
				b.fail(rng, nil, review.DuplicateName, "Check constraint "+review.Code(ident(c.Name))+" shares its name with another constraint of the table.")
				return true
			}
			named[k] = true
		case ast.ConstrPrimaryKey:
			for _, o := range c.IndexOptions {
				if o != nil && strings.EqualFold(o.Name, "IGNORED") {
					b.fail(rng, nil, review.UnsupportedForm, "A primary key can't be "+review.Code("IGNORED")+".")
					return true
				}
			}
		}
	}
	return false
}

// alterColumnsFail fails an ALTER TABLE whose column definitions or CHECK
// constraints MariaDB refuses (each verified on mariadb:11.8.8):
// columnTextFails' faults in what it adds or redefines, a CHECK constraint
// it adds calling a function whose value isn't deterministic (1901), two
// CHECK constraints it adds of one name (1826), a generated column it adds
// reading one it adds after it (4029), and ALTER INDEX of the primary key's
// IGNORED (4174).
func (b *builder) alterColumnsFail(v *ast.AlterTableStmt) bool {
	checks := make(map[string]bool)
	added := make(map[string]int)
	var adds []*ast.ColumnDef
	for _, c := range v.Commands {
		if c == nil {
			continue
		}
		rng := rangeOf(c.Loc)
		switch c.Type {
		case ast.ATAddColumn, ast.ATModifyColumn, ast.ATChangeColumn:
			if c.IfNotExists || c.IfExists {
				continue
			}
			for _, col := range columnsOf(c) {
				if b.columnTextFails(col, "") {
					return true
				}
				if c.Type == ast.ATAddColumn && col != nil {
					if k, ok := b.key(spaceColumn, col.Name); ok {
						added[k] = len(adds)
						adds = append(adds, col)
					}
				}
			}
		case ast.ATAddConstraint:
			con := c.Constraint
			if con == nil || con.Type != ast.ConstrCheck {
				continue
			}
			if con.Expr != nil && !con.IfNotExists {
				if call := nondeterministicCall(con.Expr); call != "" {
					b.fail(rng, nil, review.UnsupportedForm, cannotCall("A "+review.Code("CHECK")+" constraint", call))
					return true
				}
			}
			if con.Name == "" || con.IfNotExists {
				continue
			}
			if k, ok := b.key(spaceColumn, con.Name); ok {
				if checks[k] {
					b.fail(rng, nil, review.DuplicateName, statesTwice("adds", "check constraint", ident(con.Name)))
					return true
				}
				checks[k] = true
			}
		case ast.ATAlterIndexVisible, ast.ATAlterIndexInvisible:
			if k, ok := b.key(spaceIndex, c.Name); ok && k == primaryIndex {
				b.fail(rng, nil, review.UnsupportedForm, "A primary key can't be "+review.Code("IGNORED")+" or "+review.Code("NOT IGNORED")+".")
				return true
			}
		}
	}
	for i, col := range adds {
		if col.Generated == nil || col.Generated.Expr == nil {
			continue
		}
		for _, ref := range exprBareColumnRefs(col.Generated.Expr) {
			k, ok := b.key(spaceColumn, ref.Column)
			if j, isAdded := added[k]; ok && isAdded && j > i && adds[j].Generated != nil && adds[j].Generated.Expr != nil {
				b.fail(rangeOf(col.Loc), nil, review.UnsupportedForm, "Generated column "+review.Code(ident(col.Name))+" reads generated column "+review.Code(ident(ref.Column))+", defined after it.")
				return true
			}
		}
	}
	return false
}

// constraintNames returns the keys of the constraints and keys an ALTER
// TABLE drops by name, IF EXISTS or not, and of the names it adds: those
// of its constraints and keys, and of its columns, whose own CHECK takes
// the column's name.
func (b *builder) constraintNames(v *ast.AlterTableStmt) (dropped, added map[string]bool) {
	dropped, added = make(map[string]bool), make(map[string]bool)
	for _, c := range v.Commands {
		if c == nil {
			continue
		}
		switch c.Type {
		case ast.ATDropConstraint, ast.ATDropIndex:
			if k, ok := b.key(spaceColumn, c.Name); ok {
				dropped[k] = true
			}
		case ast.ATAddConstraint, ast.ATAddIndex:
			if c.Constraint != nil {
				for _, n := range []string{c.Constraint.Name, c.Constraint.IndexName} {
					if k, ok := b.key(spaceColumn, n); ok {
						added[k] = true
					}
				}
			}
		case ast.ATAddColumn:
			// A column's CHECK takes the column's name.
			for _, col := range columnsOf(c) {
				if col != nil {
					if k, ok := b.key(spaceColumn, col.Name); ok {
						added[k] = true
					}
				}
			}
		}
	}
	return dropped, added
}

// columnCmds returns the subcommands of an ALTER TABLE that name a column,
// as namedCount counts them.
func (b *builder) columnCmds(v *ast.AlterTableStmt, k string) []*ast.AlterTableCmd {
	var out []*ast.AlterTableCmd
	for _, c := range v.Commands {
		if c != nil && b.namedCount(&ast.AlterTableStmt{Commands: []*ast.AlterTableCmd{c}}, k) > 0 {
			out = append(out, c)
		}
	}
	return out
}

// droppedOnly reports whether an ALTER TABLE drops a column and names it
// nowhere else but in an ADD COLUMN after the drop, which a CHECK
// constraint that read the column does not read (verified on
// mariadb:11.8.8).
func (b *builder) droppedOnly(v *ast.AlterTableStmt, k string) bool {
	cmds := b.columnCmds(v, k)
	switch {
	case len(cmds) == 1:
		return cmds[0].Type == ast.ATDropColumn
	case len(cmds) == 2:
		return cmds[0].Type == ast.ATDropColumn && cmds[1].Type == ast.ATAddColumn && !cmds[1].IfNotExists
	}
	return false
}

// keptByName reports whether an ALTER TABLE keeps a column a CHECK
// constraint of cols reads: it names it nowhere, or only to rename it, by
// RENAME COLUMN or CHANGE, to a name it names nowhere else and the CHECK
// does not read, which the CHECK follows (verified on mariadb:11.8.8).
func (b *builder) keptByName(v *ast.AlterTableStmt, k string, cols []string) bool {
	cmds := b.columnCmds(v, k)
	if len(cmds) == 0 {
		return true
	}
	if len(cmds) != 1 {
		return false
	}
	c := cmds[0]
	var to string
	switch {
	case c.Type == ast.ATRenameColumn:
		to = c.NewName
	case c.Type == ast.ATChangeColumn && c.Column != nil:
		to = c.Column.Name
	default:
		return false
	}
	from, ok1 := b.key(spaceColumn, c.Name)
	nk, ok2 := b.key(spaceColumn, to)
	if !ok1 || !ok2 || from != k {
		return false
	}
	if nk == k {
		return true
	}
	return len(b.columnCmds(v, nk)) == 1 && !slices.Contains(cols, nk)
}

// alterChecksFail fails an ALTER TABLE of a synced table the change has
// not touched against its CHECK constraints and generated columns (each
// verified on mariadb:11.8.8):
//
//   - dropping some but not all of the columns a CHECK constraint reads
//     (1054), whether it renames the others or adds a dropped one back:
//     MariaDB drops a constraint all of whose columns go, and a column's
//     own CHECK with its column;
//   - adding a CHECK constraint of a name the table's CHECK constraints
//     hold (1826), or of its foreign key's, IF NOT EXISTS or not, and a
//     unique key or foreign key of a CHECK constraint's name (1826);
//   - DROP CONSTRAINT of a name that is no foreign key, key, or CHECK
//     constraint of the table (1091);
//   - MODIFY or CHANGE of a column to or from a VIRTUAL generated column,
//     but to another VIRTUAL one (1907).
func (b *builder) alterChecksFail(v *ast.AlterTableStmt, r tableRef) bool {
	st := b.syncedOf(r)
	if st == nil {
		return false
	}
	rel := relKey(r.db, r.name)
	dropped, added := b.constraintNames(v)
	dropCols := make(map[string]bool)
	for _, c := range v.Commands {
		if c != nil && c.Type == ast.ATDropColumn {
			if k, ok := b.key(spaceColumn, c.Name); ok {
				dropCols[k] = true
			}
		}
	}
	if len(dropCols) > 0 {
		names := slices.Sorted(maps.Keys(st.checks))
		for _, name := range names {
			cols := st.checks[name]
			if cols == nil || dropped[name] {
				continue
			}
			gone, kept := "", ""
			clear := true
			for _, k := range cols {
				if dropCols[k] {
					gone = k
					clear = clear && b.droppedOnly(v, k)
				} else {
					kept = k
					clear = clear && b.keptByName(v, k, cols)
				}
			}
			// A column's own CHECK, named for it, goes with the column.
			if !clear || gone == "" || kept == "" || dropCols[name] {
				continue
			}
			b.fail(b.anchor(), []walk.Key{rel, colKey(r.db, r.name, kept)}, review.UnsupportedForm, "Check constraint "+review.Code(ident(name))+" reads column "+review.Code(ident(gone))+", and keeps column "+review.Code(ident(kept))+".")
			return true
		}
	}
	for _, c := range v.Commands {
		if c == nil {
			continue
		}
		rng := rangeOf(c.Loc)
		switch c.Type {
		case ast.ATAddConstraint, ast.ATAddIndex:
			con := c.Constraint
			if con == nil {
				continue
			}
			k, ok := b.key(spaceColumn, con.Name)
			if !ok || dropped[k] {
				continue
			}
			switch con.Type {
			case ast.ConstrCheck:
				if _, has := st.checks[k]; has && !con.IfNotExists {
					b.fail(rng, []walk.Key{rel}, review.NameTaken, "Check constraint "+review.Code(ident(con.Name))+" already exists in "+review.Code(r.written)+".")
					return true
				}
				own := walk.Key{Space: spaceTableForeignKey, Parent: tablePath(r.db, r.name), Name: k}
				if f, _ := b.l.Lookup(own); f == walk.Present {
					b.fail(rng, []walk.Key{own}, review.NameTaken, "Check constraint "+review.Code(ident(con.Name))+" takes the name of a foreign key of "+review.Code(r.written)+".")
					return true
				}
			case ast.ConstrUnique, ast.ConstrForeignKey:
				if _, has := st.checks[k]; has && !con.IfNotExists {
					b.fail(rng, []walk.Key{rel}, review.NameTaken, "Check constraint "+review.Code(ident(con.Name))+" already exists in "+review.Code(r.written)+".")
					return true
				}
			}
		case ast.ATDropConstraint:
			if c.IfExists || dropKindOf(b.w.sql, c) != dropConstraint {
				continue
			}
			k, ok := b.key(spaceColumn, c.Name)
			if !ok || strings.HasPrefix(k, nonASCII) || added[k] {
				continue
			}
			if _, has := st.checks[k]; has {
				continue
			}
			own := walk.Key{Space: spaceTableForeignKey, Parent: tablePath(r.db, r.name), Name: k}
			idx := idxKey(r.db, r.name, k)
			f1, _ := b.l.Lookup(own)
			f2, _ := b.l.Lookup(idx)
			if f1 == walk.Absent && f2 == walk.Absent {
				b.fail(rng, []walk.Key{own, idx}, review.ObjectNotFound, notFound("constraint", ident(c.Name), r.written))
				return true
			}
		case ast.ATModifyColumn, ast.ATChangeColumn:
			if c.Column == nil || c.IfExists {
				continue
			}
			from := c.Column.Name
			if c.Type == ast.ATChangeColumn {
				from = c.Name
			}
			k, ok := b.key(spaceColumn, from)
			if !ok || !b.namedOnce(v, k) {
				continue
			}
			ck := colKey(r.db, r.name, k)
			if f, _ := b.l.Lookup(ck); f != walk.Present {
				continue
			}
			stored, gen := st.generated[k]
			wasVirtual := gen && !stored
			g := c.Column.Generated
			if g != nil && g.Expr == nil {
				continue
			}
			isVirtual := g != nil && !g.Stored
			if wasVirtual != isVirtual {
				b.fail(rng, []walk.Key{ck}, review.UnsupportedForm, "Column "+review.Code(ident(from))+" can't change between a "+review.Code("VIRTUAL")+" generated column and a stored or plain one.")
				return true
			}
		}
	}
	return false
}

// eventScheduleFails fails an event's schedule MariaDB refuses as it
// reads it, before it looks the event up (verified on mariadb:11.8.8): an
// EVERY interval that isn't positive (1542), and ENDS at or before STARTS
// (1543), as literals show them.
func (b *builder) eventScheduleFails(sched *ast.EventSchedule) bool {
	if sched == nil {
		return false
	}
	every := sched.Every
	if iv, ok := every.(*ast.IntervalExpr); ok {
		every = iv.Value
	}
	if every != nil {
		n, ok := partitionInt(every)
		if s, isStr := every.(*ast.StringLit); isStr {
			v, err := strconv.ParseInt(strings.TrimSpace(s.Value), 10, 64)
			n, ok = v, err == nil
		}
		if ok && n <= 0 {
			b.fail(b.anchor(), nil, review.UnsupportedForm, "An event's "+review.Code("EVERY")+" interval must be positive.")
			return true
		}
	}
	start, ok1 := literalDatetime(sched.Starts)
	end, ok2 := literalDatetime(sched.Ends)
	if ok1 && ok2 && end <= start {
		b.fail(b.anchor(), nil, review.UnsupportedForm, "An event's "+review.Code("ENDS")+" comes at or before its "+review.Code("STARTS")+".")
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
