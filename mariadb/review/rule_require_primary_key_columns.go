package review

import (
	"slices"
	"strconv"
	"strings"

	"github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/mariadb/parser"
	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/review"
)

// An ALTER TABLE that drops every column of a table's primary key takes
// the key with it, which only the target's synced schema tells. The rule
// judges such a drop against each target with a closed model of the
// change, read in order:
//
//   - An ALTER TABLE whose every subcommand is DROP COLUMN, or ALGORITHM
//     or LOCK left at DEFAULT, is followed exactly: the columns it drops,
//     whether the server refuses it, and whether it takes the key.
//   - A statement that cannot fail and changes no schema object passes: a
//     transaction statement, but one of a savepoint the transaction may
//     lack or one that ends the session; UNLOCK TABLES; table maintenance,
//     which reports a missing table in its result; and a SET of user
//     variables to literals, or of the session's foreign_key_checks or
//     unique_checks.
//   - Any other statement, and a drop the server refuses or may refuse,
//     ends the judging of later drops, as pg/review's scan ends at a
//     statement it cannot follow: the rule cannot establish that it runs,
//     and one that fails stops the change. DML is one, since a query of a
//     table or column the target lacks fails, and so is any other SET,
//     which may name a variable the server lacks or refuses to set. A
//     finding before it stands, unless a statement that may give its table
//     a key back, rename the table, drop it, or create a table of its name
//     withdraws the finding.
//   - A statement that may change any table, such as a CALL or an event,
//     or a drop of the target's database, ends the judging and withdraws
//     every finding.
//
// mariadb: https://mariadb.com/docs/server/reference/sql-statements/data-definition/alter/alter-table — DROP COLUMN

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
	// judging and withdraws a finding on the tables it may give a key
	// back, rename, drop, or create.
	namesTables
	// dropsDatabase is a DROP DATABASE, which withdraws every finding when
	// it drops the target's.
	dropsDatabase
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
				st.kind, st.tables = namesTables, alterSettles(v)
			}
		case *ast.CreateTableStmt:
			st.kind = namesTables
			// IF NOT EXISTS of a table that exists does nothing, and a temporary
			// table leaves the table of its name as it is.
			if !v.IfNotExists && !v.Temporary {
				st.tables = []*ast.TableRef{v.Table}
			}
		case *ast.CreateIndexStmt, *ast.DropIndexStmt:
			// CREATE INDEX makes no primary key.
			st.kind = namesTables
		case *ast.DropTableStmt:
			st.kind = namesTables
			// DROP TEMPORARY TABLE drops no permanent table.
			if !v.Temporary {
				st.tables = v.Tables
			}
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
			if setPasses(v) {
				continue
			}
			st.kind = barrier
		case *ast.BeginStmt, *ast.UnlockTablesStmt,
			*ast.CheckTableStmt, *ast.ChecksumTableStmt, *ast.OptimizeTableStmt, *ast.RepairTableStmt:
			continue
		case *ast.AnalyzeTableStmt:
			// The parser takes MySQL's histogram clauses, which MariaDB refuses.
			if v.HistogramOp == "" {
				continue
			}
			st.kind = barrier
		case *ast.CommitStmt:
			// COMMIT RELEASE ends the session.
			if !v.Release {
				continue
			}
			st.kind = barrier
		case *ast.RollbackStmt:
			// ROLLBACK TO a savepoint the transaction lacks fails.
			if !v.Release && v.Savepoint == "" {
				continue
			}
			st.kind = barrier
		case *ast.SavepointStmt:
			// RELEASE of a savepoint the transaction lacks fails.
			if !v.Release {
				continue
			}
			st.kind = barrier
		case *ast.SelectStmt, *ast.InsertStmt, *ast.UpdateStmt, *ast.DeleteStmt, *ast.LoadDataStmt, *ast.TableStmt, *ast.ValuesStmt, *ast.DoStmt,
			*ast.ExplainStmt, *ast.ShowStmt, *ast.HelpStmt, *ast.TruncateStmt, *ast.LockTablesStmt, *ast.SetTransactionStmt, *ast.UseStmt:
			// It changes no schema object, but may fail: on a table or column
			// the target lacks; for TRUNCATE, on a table a foreign key
			// references; for SET TRANSACTION, inside a transaction. Under LOCK
			// TABLES, an ALTER TABLE of a table it did not lock for writing
			// fails, and after USE, a name without a database may mean another
			// table.
			st.kind = barrier
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
		if cmd != nil && leftDefault(sql, cmd) {
			continue
		}
		if cmd == nil || dropKindOf(sql, cmd) != dropColumn {
			return nil, false
		}
		drops = append(drops, cmd)
	}
	return drops, len(drops) > 0
}

// leftDefault reports whether a subcommand is ALGORITHM = DEFAULT or LOCK =
// DEFAULT, the same as leaving it out. The last ALGORITHM and the last LOCK
// of a statement are the ones that hold, so a statement with no other
// value for them runs as without. A quoted DEFAULT is a string, which
// the server refuses, or a name, which the rule does not judge.
//
// mariadb: https://mariadb.com/docs/server/reference/sql-statements/data-definition/alter/alter-table — ALGORITHM, LOCK
func leftDefault(sql string, cmd *ast.AlterTableCmd) bool {
	if cmd.Type != ast.ATAlgorithm && cmd.Type != ast.ATLock {
		return false
	}
	if cmd.Loc.Start < 0 || cmd.Loc.End > len(sql) || cmd.Loc.End <= cmd.Loc.Start {
		return false
	}
	text := sql[cmd.Loc.Start:cmd.Loc.End]
	tokens := parser.Tokenize(text)
	return len(tokens) > 0 && strings.EqualFold(tokenText(text, tokens[len(tokens)-1]), "DEFAULT")
}

// alterSettles returns the tables an ALTER TABLE may give a primary key
// back or rename: its own, and the one it renames to.
func alterSettles(v *ast.AlterTableStmt) []*ast.TableRef {
	if v.Table == nil {
		return nil
	}
	var names []*ast.TableRef
	for _, cmd := range v.Commands {
		switch {
		case cmd == nil:
		case cmd.Type == ast.ATRenameTable:
			names = append(names, v.Table, renamedTo(cmd))
		case addsPrimaryKey(cmd):
			names = append(names, v.Table)
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

// setPasses reports whether a SET cannot fail: every assignment sets a
// user variable to a literal, or the session's foreign_key_checks or
// unique_checks, which need no privilege, to on or off. Any other may name
// a variable the server lacks, one it reads only, or one that needs a
// privilege, and a failed assignment fails the statement.
//
// mariadb: https://mariadb.com/docs/server/reference/sql-statements/administrative-sql-statements/set-commands/set
func setPasses(v *ast.SetStmt) bool {
	for i, a := range v.Assignments {
		if a == nil || a.Column == nil || a.Column.Table != "" {
			return false
		}
		name, scope := strings.ToLower(a.Column.Column), strings.ToUpper(v.Scope)
		if strings.HasPrefix(name, "@") && !strings.HasPrefix(name, "@@") {
			// A scope names the system variable after it, and the parser
			// takes SET SESSION @x = 1, which the server refuses.
			if i == 0 && v.Scope != "" || !literal(a.Value) {
				return false
			}
			continue
		}
		if strings.HasPrefix(name, "@@") {
			name = strings.TrimPrefix(name, "@@")
			if dot := strings.IndexByte(name, '.'); dot >= 0 {
				name, scope = name[dot+1:], strings.ToUpper(name[:dot])
			}
		}
		if name != "foreign_key_checks" && name != "unique_checks" || scope != "" && scope != "SESSION" && scope != "LOCAL" || !onOrOff(a.Value) {
			return false
		}
	}
	return len(v.Assignments) > 0
}

// literal reports whether an expression is a literal the server takes for
// a user variable as it is: a number, or a sign of one; a string without a
// character set introducer; a hexadecimal or bit value; TRUE, FALSE, or
// NULL. The parser takes X'4' and b'2', which the server refuses as
// malformed, and 1e400, which overflows a double.
func literal(e ast.ExprNode) bool {
	switch x := e.(type) {
	case *ast.StringLit:
		return x.Charset == ""
	case *ast.HexLit:
		digits, prefixed := strings.CutPrefix(strings.ToLower(x.Value), "0x")
		return (prefixed && digits != "" || !prefixed && len(digits)%2 == 0) && strings.Trim(digits, "0123456789abcdef") == ""
	case *ast.BitLit:
		digits, prefixed := strings.CutPrefix(strings.ToLower(x.Value), "0b")
		return (!prefixed || digits != "") && strings.Trim(digits, "01") == ""
	case *ast.BoolLit, *ast.NullLit:
		return true
	}
	return number(e)
}

// number reports whether an expression is a number the server holds, or a
// sign of one. Digits alone always fit, as a decimal or a double; an
// exponent out of a double's range does not.
func number(e ast.ExprNode) bool {
	switch x := e.(type) {
	case *ast.IntLit:
		return true
	case *ast.FloatLit:
		if !strings.ContainsAny(x.Value, "eE") {
			return true
		}
		_, err := strconv.ParseFloat(x.Value, 64)
		return err == nil
	case *ast.UnaryExpr:
		return (x.Op == ast.UnaryMinus || x.Op == ast.UnaryPlus) && number(x.Operand)
	}
	return false
}

// onOrOff reports whether an expression is a value of a boolean system
// variable: ON, OFF, TRUE, FALSE, 1, or 0.
func onOrOff(e ast.ExprNode) bool {
	switch x := e.(type) {
	case *ast.BoolLit:
		return true
	case *ast.IntLit:
		return x.Value == 0 || x.Value == 1
	case *ast.StringLit:
		return x.Charset == "" && (strings.EqualFold(x.Value, "ON") || strings.EqualFold(x.Value, "OFF"))
	}
	return false
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
	// judging is whether the rule still judges drops.
	judging := true
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
			if !runs {
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
// the table no visible column; when something depends on a column it drops;
// and when it drops only some columns of a key of several.
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

// keepsColumn reports whether something the synced schema holds may keep the
// server from dropping a column of the table, given the columns gone with
// this statement or before it: a foreign key of the table that holds it, or
// one of the database that references it; a generated column not gone that
// may read it, or whose default may name it; a check that reads it and
// another column, which MariaDB drops when every column it reads goes, but
// which may be an application-time period, synced as a check of its two
// columns, that keeps them; a unique key of it and a column not gone, which
// MariaDB refuses to narrow; or the table's partitioning or subpartitioning,
// by an expression that may read it or by KEY (), which reads the primary
// key, or a unique key of a table without one. A check that reads the column
// alone goes with it. A foreign key names the table it references without
// its database, so one that may reference a table of another database
// counts. MariaDB has no functional index.
//
// mariadb: https://mariadb.com/docs/server/reference/sql-structure/temporal-tables/application-time-periods
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
		// The synced default does not tell an expression from a literal: one
		// that may name the column counts.
		if !gone(c.GetName()) && mentions(c.GetDefault(), column) {
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
		if idx.GetUnique() && !idx.GetPrimary() && has(idx.GetExpressions()) && slices.ContainsFunc(idx.GetExpressions(), func(part string) bool { return !gone(part) }) {
			return true
		}
	}
	// KEY () is synced with no expression, and reads the primary key, or,
	// without one, a unique key: a column of none of them is not read.
	partitionedBy := func(expr string) bool {
		if expr != "" {
			return reads(expr, column)
		}
		if key := primaryKeyColumns(t); key != nil {
			return has(key)
		}
		return !slices.ContainsFunc(t.GetIndexes(), (*metadata.IndexMetadata).GetUnique) ||
			slices.ContainsFunc(t.GetIndexes(), func(idx *metadata.IndexMetadata) bool { return idx.GetUnique() && has(idx.GetExpressions()) })
	}
	for _, p := range t.GetPartitions() {
		if partitionedBy(p.GetExpression()) {
			return true
		}
		for _, sub := range p.GetSubpartitions() {
			if partitionedBy(sub.GetExpression()) {
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
