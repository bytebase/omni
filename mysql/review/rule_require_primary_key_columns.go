package review

import (
	"slices"
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/mysql/parser"
	"github.com/bytebase/omni/review"
)

// An ALTER TABLE that drops every column of a table's primary key takes
// the key with it, which only the target's synced schema tells. The rule
// judges such a drop against each target with a closed model of the
// change, read in order:
//
//   - An ALTER TABLE whose every subcommand is DROP COLUMN is followed
//     exactly: the columns it drops, whether the server refuses it, and
//     whether it takes the key.
//   - A statement that changes no schema object, such as DML, a session
//     statement, or table maintenance, passes: whether it fails depends on
//     data the review does not see, as in pg/review.
//   - Any other statement, and a drop the server refuses or may refuse,
//     ends the judging of later drops, as pg/review's scan ends at a
//     statement it cannot follow: the rule cannot establish that it runs,
//     and one that fails stops the change. A finding before it stands,
//     unless a statement on its table that may give it a key back, rename
//     it, or drop it withdraws the finding.
//   - A statement that may change any table, such as a CALL or an event,
//     or a drop of the target's database, ends the judging and withdraws
//     every finding.
//
// mysql: https://dev.mysql.com/doc/refman/8.0/en/alter-table.html#alter-table-add-drop-column

// columnDrops is the change as the targets judge its column drops: the
// top-level statements that do not pass, in order.
type columnDrops struct {
	steps []step
}

// any reports whether the change has a column drop to judge.
func (c columnDrops) any() bool {
	return slices.ContainsFunc(c.steps, func(s step) bool { return s.kind == dropsColumns })
}

type stepKind int

const (
	// dropsColumns is an ALTER TABLE whose every subcommand is DROP COLUMN.
	dropsColumns stepKind = iota
	// namesTables is any other statement on a table, which ends the
	// judging and withdraws a finding on the tables it changes.
	namesTables
	// dropsDatabase is a DROP DATABASE, which withdraws every finding when
	// it drops the target's.
	dropsDatabase
	// requiresKey is a SET of sql_require_primary_key or
	// sql_generate_invisible_primary_key for the session.
	requiresKey
	// barrier is any other statement the rule cannot follow.
	barrier
	// stops is a statement that may change any table.
	stops
)

type step struct {
	kind stepKind
	// statement is the index in Result.Statements.
	statement int
	tables    []*ast.TableRef
	drops     []*ast.AlterTableCmd
	database  string
	// settings are what a requiresKey step sets.
	settings []keySetting
}

// columnDropSteps reads the change's top-level statements into steps.
func columnDropSteps(sql string, stmts []statement) columnDrops {
	var c columnDrops
	for i := range stmts {
		s := &stmts[i]
		st := step{statement: s.index}
		switch v := s.node.(type) {
		case *ast.AlterTableStmt:
			if drops, ok := onlyColumnDrops(sql, v); ok {
				st.kind, st.tables, st.drops = dropsColumns, []*ast.TableRef{v.Table}, drops
			} else {
				st.kind, st.tables = namesTables, alterNames(v)
			}
		case *ast.CreateTableStmt:
			st.kind, st.tables = namesTables, []*ast.TableRef{v.Table}
		case *ast.CreateIndexStmt:
			st.kind, st.tables = namesTables, []*ast.TableRef{v.Table}
		case *ast.DropIndexStmt:
			st.kind, st.tables = namesTables, []*ast.TableRef{v.Table}
		case *ast.DropTableStmt:
			st.kind, st.tables = namesTables, v.Tables
		case *ast.RenameTableStmt:
			st.kind = namesTables
			for _, pair := range v.Pairs {
				if pair != nil {
					st.tables = append(st.tables, pair.Old, pair.New)
				}
			}
		case *ast.DropDatabaseStmt:
			st.kind, st.database = dropsDatabase, v.Name
		case *ast.SetStmt:
			settings := keySettings(v)
			if len(settings) == 0 {
				continue
			}
			st.kind, st.settings = requiresKey, settings
		case *ast.UseStmt:
			st.kind = barrier
		case *ast.SelectStmt, *ast.InsertStmt, *ast.UpdateStmt, *ast.DeleteStmt, *ast.LoadDataStmt, *ast.TableStmt, *ast.ValuesStmt, *ast.DoStmt,
			*ast.ExplainStmt, *ast.ShowStmt, *ast.HelpStmt,
			*ast.BeginStmt, *ast.CommitStmt, *ast.RollbackStmt, *ast.SavepointStmt, *ast.SetTransactionStmt,
			*ast.AnalyzeTableStmt, *ast.CheckTableStmt, *ast.ChecksumTableStmt, *ast.OptimizeTableStmt, *ast.RepairTableStmt, *ast.TruncateStmt:
			// It changes no schema object.
			continue
		case *ast.CreateViewStmt, *ast.AlterViewStmt, *ast.DropViewStmt, *ast.CreateTriggerStmt, *ast.DropTriggerStmt,
			*ast.CreateFunctionStmt, *ast.AlterRoutineStmt, *ast.DropRoutineStmt, *ast.DropEventStmt,
			*ast.CreateDatabaseStmt, *ast.AlterDatabaseStmt,
			*ast.CreateUserStmt, *ast.AlterUserStmt, *ast.DropUserStmt, *ast.RenameUserStmt, *ast.SetPasswordStmt, *ast.GrantStmt, *ast.RevokeStmt,
			*ast.CreateRoleStmt, *ast.DropRoleStmt, *ast.GrantRoleStmt, *ast.RevokeRoleStmt, *ast.SetRoleStmt, *ast.SetDefaultRoleStmt:
			// It changes no table, but may fail. A trigger or a function cannot
			// run DDL, and a procedure runs only when called.
			st.kind = barrier
		default:
			// CALL, PREPARE and EXECUTE, an event, which runs on its own
			// schedule, and whatever else may change any table.
			st.kind = stops
		}
		c.steps = append(c.steps, st)
	}
	return c
}

// onlyColumnDrops returns the subcommands of an ALTER TABLE whose every
// subcommand is DROP COLUMN.
func onlyColumnDrops(sql string, v *ast.AlterTableStmt) ([]*ast.AlterTableCmd, bool) {
	if v.Table == nil {
		return nil, false
	}
	var drops []*ast.AlterTableCmd
	for _, cmd := range v.Commands {
		if cmd == nil || dropKindOf(sql, cmd) != dropColumn {
			return nil, false
		}
		drops = append(drops, cmd)
	}
	return drops, len(drops) > 0
}

// alterNames returns the tables an ALTER TABLE changes: its own, and the
// one it renames to.
func alterNames(v *ast.AlterTableStmt) []*ast.TableRef {
	if v.Table == nil {
		return nil
	}
	names := []*ast.TableRef{v.Table}
	for _, cmd := range v.Commands {
		if cmd != nil && cmd.Type == ast.ATRenameTable {
			names = append(names, renamedTo(cmd))
		}
	}
	return names
}

// renamedTo is the table an ALTER TABLE ... RENAME subcommand renames to.
// A new name without a database is in the current one, whatever database
// the table is in.
//
// mysql: https://dev.mysql.com/doc/refman/8.0/en/identifier-qualifiers.html
func renamedTo(cmd *ast.AlterTableCmd) *ast.TableRef {
	if cmd.NewTable != nil {
		return cmd.NewTable
	}
	return &ast.TableRef{Name: cmd.NewName}
}

// keyRequired reports whether a session setting refuses, or may refuse, a
// statement that leaves a table without a primary key.
func keyRequired(settings map[string]bool) bool {
	for _, on := range settings {
		if on {
			return true
		}
	}
	return false
}

// keySetting is what a SET does to one of the session's settings under
// which the server refuses a statement that leaves a table without a
// primary key: on or off, when known.
type keySetting struct {
	name      string
	on, known bool
}

// keySettings reads what a SET does to the session's
// sql_require_primary_key and sql_generate_invisible_primary_key; under
// either, on, the server refuses a statement that leaves a table without a
// primary key. A GLOBAL or PERSIST assignment leaves the session as it is,
// and a user variable of the name is another variable. DEFAULT, or an
// expression, sets a value the change does not tell.
//
// mysql: https://dev.mysql.com/doc/refman/8.0/en/server-system-variables.html#sysvar_sql_require_primary_key
// mysql: https://dev.mysql.com/doc/refman/8.0/en/create-table-gipks.html
func keySettings(v *ast.SetStmt) []keySetting {
	var out []keySetting
	for _, a := range v.Assignments {
		if a == nil || a.Column == nil {
			continue
		}
		name, scope := strings.ToLower(a.Column.Column), strings.ToUpper(v.Scope)
		if strings.HasPrefix(name, "@@") {
			name = strings.TrimPrefix(name, "@@")
			for _, s := range []string{"global.", "persist.", "persist_only.", "session.", "local."} {
				if strings.HasPrefix(name, s) {
					name, scope = strings.TrimPrefix(name, s), strings.ToUpper(strings.TrimSuffix(s, "."))
				}
			}
		}
		if name != "sql_require_primary_key" && name != "sql_generate_invisible_primary_key" || scope == "GLOBAL" || scope == "PERSIST" || scope == "PERSIST_ONLY" {
			continue
		}
		k := keySetting{name: name}
		switch x := a.Value.(type) {
		case *ast.BoolLit:
			k.on, k.known = x.Value, true
		case *ast.IntLit:
			k.on, k.known = x.Value != 0, x.Value == 0 || x.Value == 1
		case *ast.StringLit:
			k.on, k.known = strings.EqualFold(x.Value, "ON"), strings.EqualFold(x.Value, "ON") || strings.EqualFold(x.Value, "OFF")
		}
		out = append(out, k)
	}
	return out
}

// checkColumnDrops judges the change's column drops against one target.
func checkColumnDrops(c columnDrops, x schemaIndex, r *reporter) {
	if !c.any() {
		return
	}
	type found struct {
		table     *metadata.TableMetadata
		statement int
		rng       review.Range
		message   string
	}
	var findings []found
	dropped := make(map[*metadata.TableMetadata][]string)
	j := &columnJudge{x: x}
	// judging is whether the rule still judges drops, and requireKey holds
	// the session settings under which the server refuses one that takes a
	// table's key, or may, when the value is not known.
	judging, requireKey := true, make(map[string]bool)
	for _, s := range c.steps {
		switch s.kind {
		case stops:
			return
		case dropsDatabase:
			if x.same(s.database, x.schema.GetName()) {
				return
			}
			judging = false
		case barrier:
			judging = false
		case requiresKey:
			for _, k := range s.settings {
				requireKey[k.name] = k.on || !k.known
			}
		case namesTables:
			judging = false
			for _, ref := range s.tables {
				if t := x.table(ref, false); t != nil {
					findings = slices.DeleteFunc(findings, func(f found) bool { return f.table == t })
				}
			}
		case dropsColumns:
			if !judging {
				continue
			}
			ref := s.tables[0]
			t := x.table(ref, false)
			// A table the synced schema does not hold may be missing, and the
			// statement then fails.
			if t == nil {
				judging = false
				continue
			}
			columns, lost, runs := j.dropColumns(t, s.drops, dropped[t])
			if !runs || lost != nil && keyRequired(requireKey) {
				judging = false
				continue
			}
			dropped[t] = append(dropped[t], columns...)
			if lost != nil {
				findings = append(findings, found{t, s.statement, rangeOf(lost.Loc), "removes the primary key of " + table(ref) + ", and the change adds none back"})
			}
		}
	}
	for _, f := range findings {
		r.report(review.RequirePrimaryKey, f.statement, f.rng, f.message)
	}
}

// columnJudge judges the column drops of one target.
type columnJudge struct {
	x schemaIndex
	// referencing holds the foreign keys of the synced schema by the table
	// they reference, built when first needed.
	referencing map[string][]referencingKey
}

// referencingKey is a foreign key of a table of the synced schema.
type referencingKey struct {
	table *metadata.TableMetadata
	fk    *metadata.ForeignKeyMetadata
}

// dropColumns judges an ALTER TABLE of DROP COLUMN subcommands on a synced
// table, given the columns statements before it dropped: the columns it
// drops, the subcommand that takes the primary key, nil when it leaves the
// key or the table has none, and whether it is known to run.
//
// The server refuses the statement when it drops a column the table lacks,
// or one twice, but under IF EXISTS, which passes over it; when it leaves
// the table no visible column; and when something depends on a column it
// drops. A drop of only some columns of a key of several shrinks the key
// to the rest, which fails when its values repeat: what follows is not
// known.
func (j *columnJudge) dropColumns(t *metadata.TableMetadata, drops []*ast.AlterTableCmd, before []string) (columns []string, lost *ast.AlterTableCmd, runs bool) {
	in := func(list []string, column string) bool {
		return slices.ContainsFunc(list, func(c string) bool { return sameName(c, column) })
	}
	has := func(column string) bool {
		return !in(before, column) && slices.ContainsFunc(t.GetColumns(), func(c *metadata.ColumnMetadata) bool { return sameName(c.GetName(), column) })
	}
	var effective []*ast.AlterTableCmd
	for _, cmd := range drops {
		if in(columns, cmd.Name) || !has(cmd.Name) {
			if !cmd.IfExists {
				return nil, nil, false
			}
			continue
		}
		columns = append(columns, cmd.Name)
		effective = append(effective, cmd)
	}
	gone := func(column string) bool { return in(before, column) || in(columns, column) }
	if !slices.ContainsFunc(t.GetColumns(), func(c *metadata.ColumnMetadata) bool { return !gone(c.GetName()) && !c.GetIsInvisible() }) {
		return nil, nil, false
	}
	for _, column := range columns {
		if j.keepsColumn(t, column, gone) {
			return nil, nil, false
		}
	}
	var key []string
	for _, column := range primaryKeyColumns(t) {
		if !in(before, column) {
			key = append(key, column)
		}
	}
	taken := 0
	for _, column := range key {
		if in(columns, column) {
			taken++
		}
	}
	switch {
	case taken == 0:
	case taken < len(key):
		return nil, nil, false
	default:
		lost = effective[slices.IndexFunc(effective, func(cmd *ast.AlterTableCmd) bool { return in(key, cmd.Name) })]
	}
	return columns, lost, true
}

// primaryKeyColumns returns the columns of the synced table's primary
// key, nil when it has none, or a part of it names no column of the table.
func primaryKeyColumns(t *metadata.TableMetadata) []string {
	i := slices.IndexFunc(t.GetIndexes(), (*metadata.IndexMetadata).GetPrimary)
	if i < 0 {
		return nil
	}
	key := t.GetIndexes()[i].GetExpressions()
	for _, part := range key {
		if !slices.ContainsFunc(t.GetColumns(), func(c *metadata.ColumnMetadata) bool { return sameName(c.GetName(), part) }) {
			return nil
		}
	}
	if len(key) == 0 {
		return nil
	}
	return key
}

// keepsColumn reports whether something the synced schema holds may keep
// the server from dropping a column of the table, given the columns gone
// with this statement or before it: a foreign key of the table that holds
// it, or one of the database that references it; a generated column not
// gone that may read it, or whose expression default may; a check that
// reads it and another column; a
// functional index part that may read it; or the table's partitioning or
// subpartitioning, by an expression that may read it or by KEY (), which
// reads the primary key. A check that reads the column alone goes with
// it. A foreign key names the table it references without its database,
// so one that may reference a table of another database counts.
func (j *columnJudge) keepsColumn(t *metadata.TableMetadata, column string, gone func(string) bool) bool {
	has := func(list []string) bool {
		return slices.ContainsFunc(list, func(c string) bool { return sameName(c, column) })
	}
	if slices.ContainsFunc(t.GetForeignKeys(), func(fk *metadata.ForeignKeyMetadata) bool { return has(fk.GetColumns()) }) {
		return true
	}
	for _, ref := range j.referencingKeys(t) {
		if has(ref.fk.GetReferencedColumns()) {
			return true
		}
	}
	for _, c := range t.GetColumns() {
		if g := c.GetGeneration(); g != nil && !gone(c.GetName()) && reads(g.GetExpression(), column) {
			return true
		}
		// An expression default is synced in parentheses, a literal one
		// quoted.
		if d := c.GetDefault(); strings.HasPrefix(d, "(") && !gone(c.GetName()) && reads(d, column) {
			return true
		}
	}
	for _, c := range t.GetCheckConstraints() {
		reads, ok := expressionColumns(c.GetExpression())
		if !ok && mentions(c.GetExpression(), column) || has(reads) && slices.ContainsFunc(reads, func(r string) bool { return !sameName(r, column) }) {
			return true
		}
	}
	for _, idx := range t.GetIndexes() {
		for _, part := range idx.GetExpressions() {
			if !sameName(part, column) && reads(part, column) {
				return true
			}
		}
	}
	for _, p := range t.GetPartitions() {
		if p.GetExpression() == "" || reads(p.GetExpression(), column) {
			return true
		}
		for _, sub := range p.GetSubpartitions() {
			if sub.GetExpression() == "" || reads(sub.GetExpression(), column) {
				return true
			}
		}
	}
	return false
}

// referencingKeys returns the foreign keys of the synced schema that
// reference the table, in its database or one they do not name.
func (j *columnJudge) referencingKeys(t *metadata.TableMetadata) []referencingKey {
	if j.referencing == nil {
		j.referencing = make(map[string][]referencingKey)
		for _, s := range j.x.schema.GetSchemas() {
			for _, table := range s.GetTables() {
				for _, fk := range table.GetForeignKeys() {
					if db := fk.GetReferencedSchema(); db == "" || j.x.same(db, j.x.schema.GetName()) {
						k := strings.ToLower(fk.GetReferencedTable())
						j.referencing[k] = append(j.referencing[k], referencingKey{table, fk})
					}
				}
			}
		}
	}
	var out []referencingKey
	for _, ref := range j.referencing[strings.ToLower(t.GetName())] {
		if j.x.same(ref.fk.GetReferencedTable(), t.GetName()) {
			out = append(out, ref)
		}
	}
	return out
}

// expressionColumns parses an expression of the synced schema and returns
// the columns it reads. ok is false when it does not parse.
func expressionColumns(expr string) (columns []string, ok bool) {
	list, err := parser.Parse("SELECT " + expr)
	if err != nil || list == nil || len(list.Items) != 1 {
		return nil, false
	}
	ast.Inspect(list.Items[0], func(n ast.Node) bool {
		if c, ok := n.(*ast.ColumnRef); ok {
			columns = append(columns, c.Column)
		}
		return true
	})
	return columns, true
}

// reads reports whether an expression of the synced schema reads a column:
// one it names, when it parses, which a literal of the text does not; when
// it does not, one it may name.
func reads(expr, column string) bool {
	if columns, ok := expressionColumns(expr); ok {
		return slices.ContainsFunc(columns, func(c string) bool { return sameName(c, column) })
	}
	return mentions(expr, column)
}

// mentions reports whether an expression may name a column: the name,
// compared without case, as a whole word or quoted in backticks.
func mentions(expr, column string) bool {
	expr, column = strings.ToLower(expr), strings.ToLower(column)
	if column == "" {
		return false
	}
	if strings.Contains(expr, "`"+strings.ReplaceAll(column, "`", "``")+"`") {
		return true
	}
	word := func(b byte) bool { return isWordByte(b) || b >= 0x80 }
	for i := 0; ; {
		j := strings.Index(expr[i:], column)
		if j < 0 {
			return false
		}
		j += i
		end := j + len(column)
		if (j == 0 || !word(expr[j-1])) && (end == len(expr) || !word(expr[end])) {
			return true
		}
		i = j + 1
	}
}
