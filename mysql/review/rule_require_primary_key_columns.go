package review

import (
	"slices"
	"strconv"
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
//     a key back, rename the table, or drop it withdraws the finding. A
//     temporary table of the name takes an ALTER TABLE or a DROP TABLE
//     instead, and CREATE TABLE of the name, or a rename onto it, fails
//     while the table exists.
//   - A statement the rule shows to fail withdraws nothing, and the change
//     stops at it: an ALTER TABLE that drops a column the table lacks, and
//     a rename onto a table the target holds, while no statement the rule
//     does not follow has named either table. A DROP TABLE of a table a
//     foreign key references is not shown to fail: the synced schema does
//     not give the database a foreign key references, which may be
//     another's.
//   - A USE of the target's database passes. A USE of another ends the
//     judging, and a name without a database after it is not the
//     target's: the USE either moves the change there or fails.
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
	// judging; tableState tells which findings it withdraws.
	namesTables
	// dropsDatabase is a DROP DATABASE, which withdraws every finding when
	// it drops the target's.
	dropsDatabase
	// uses is a USE, which moves names without a database to its database.
	uses
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
	// node is the statement a namesTables step is.
	node ast.Node
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
				st.kind, st.node = namesTables, v
			}
		case *ast.CreateTableStmt, *ast.CreateIndexStmt, *ast.DropIndexStmt, *ast.DropTableStmt, *ast.RenameTableStmt:
			st.kind, st.node = namesTables, v
		case *ast.DropDatabaseStmt:
			st.kind, st.database = dropsDatabase, v.Name
		case *ast.SetStmt:
			if setPasses(v) {
				continue
			}
			st.kind = barrier
		case *ast.BeginStmt, *ast.UnlockTablesStmt,
			*ast.AnalyzeTableStmt, *ast.CheckTableStmt, *ast.ChecksumTableStmt, *ast.OptimizeTableStmt, *ast.RepairTableStmt:
			continue
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
			*ast.ExplainStmt, *ast.ShowStmt, *ast.HelpStmt, *ast.TruncateStmt, *ast.LockTablesStmt, *ast.FlushStmt, *ast.SetTransactionStmt:
			// It changes no schema object, but may fail: on a table or column
			// the target lacks; for TRUNCATE, on a table a foreign key
			// references; for SET TRANSACTION, inside a transaction. Under LOCK
			// TABLES, an ALTER TABLE of a table it did not lock for writing
			// fails, and after FLUSH TABLES WITH READ LOCK, any does.
			st.kind = barrier
		case *ast.UseStmt:
			st.kind, st.database = uses, v.Database
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
	// The last ALGORITHM and the last LOCK of a statement hold, and either
	// at DEFAULT is the same as left out.
	last := make(map[ast.AlterTableCmdType]bool)
	for _, cmd := range v.Commands {
		if cmd != nil && (cmd.Type == ast.ATAlgorithm || cmd.Type == ast.ATLock) {
			last[cmd.Type] = leftDefault(sql, cmd)
			continue
		}
		if cmd == nil || dropKindOf(sql, cmd) != dropColumn {
			return nil, false
		}
		drops = append(drops, cmd)
	}
	for _, isDefault := range last {
		if !isDefault {
			return nil, false
		}
	}
	return drops, len(drops) > 0
}

// leftDefault reports whether a subcommand is ALGORITHM = DEFAULT or LOCK =
// DEFAULT, the same as leaving it out. A quoted DEFAULT is a string, which
// the server refuses, or a name, which it does not know.
//
// mysql: https://dev.mysql.com/doc/refman/8.0/en/alter-table.html#alter-table-performance
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
// mysql: https://dev.mysql.com/doc/refman/8.0/en/set-variable.html
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
// NULL. The parser takes X'4', b'2', and 0X41, which the server refuses
// as malformed, and 1e400, which overflows a double.
func literal(e ast.ExprNode) bool {
	switch x := e.(type) {
	case *ast.StringLit:
		return x.Charset == ""
	case *ast.HexLit:
		// The 0x and 0b prefixes are lowercase: 0X41 is a name.
		digits, prefixed := strings.CutPrefix(x.Value, "0x")
		return (prefixed && digits != "" || !prefixed && len(digits)%2 == 0) && strings.Trim(digits, "0123456789abcdefABCDEF") == ""
	case *ast.BitLit:
		digits, prefixed := strings.CutPrefix(x.Value, "0b")
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
// variable: ON, OFF, TRUE, FALSE, 1, 0, or DEFAULT, which sets the
// session's value to the global one.
func onOrOff(e ast.ExprNode) bool {
	switch x := e.(type) {
	case *ast.BoolLit, *ast.DefaultExpr:
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
	w := &tableState{x: x, touched: make(map[*metadata.TableMetadata]bool), shadowed: make(map[*metadata.TableMetadata]bool), gone: make(map[*metadata.TableMetadata][]string)}
	// judging is whether the rule still judges drops.
	judging := true
read:
	for _, s := range c.steps {
		switch s.kind {
		case stops:
			return
		case dropsDatabase:
			if x.same(s.database, x.schema.GetName()) {
				return
			}
			judging = false
		case uses:
			if w.moved = !x.same(s.database, x.schema.GetName()); w.moved {
				judging = false
			}
		case barrier:
			judging = false
		case namesTables:
			judging = false
			// The change stops at a statement that fails, and what stands
			// before it is the change's.
			if w.fails(s.node) {
				break read
			}
			for _, t := range w.settles(s.node) {
				findings = slices.DeleteFunc(findings, func(f found) bool { return f.table == t })
			}
			w.touch(s.node)
		case dropsColumns:
			ref := s.tables[0]
			t := w.table(ref)
			// So does one that drops a column a known table lacks.
			if _, ok := w.known(ref); ok && w.dropsLacking(t, s.drops) {
				break read
			}
			if judging {
				var columns []string
				var lost *ast.AlterTableCmd
				runs := false
				// A table the synced schema does not hold may be missing, and
				// the statement then fails.
				if t != nil {
					columns, lost, runs = j.dropColumns(t, s.drops, dropped[t])
				}
				if judging = runs; runs {
					dropped[t] = append(dropped[t], columns...)
					if lost != nil {
						findings = append(findings, found{t, s.statement, rangeOf(lost.Loc), "removes the primary key of " + table(ref) + ", and the change adds none back"})
					}
				}
			}
			// The table lacks the columns the statement names after it, or
			// the change stops at it.
			if t != nil {
				for _, cmd := range s.drops {
					w.gone[t] = append(w.gone[t], cmd.Name)
				}
			}
		}
	}
	for _, f := range findings {
		r.report(review.RequirePrimaryKey, f.statement, f.rng, f.message)
	}
}

// tableState is what the rule knows of the target's tables as it reads the
// change past the statements it follows. Only a statement it does not
// follow that names a synced table leaves the table unknown: the
// statements it passes and its barriers change no table, and a DROP COLUMN
// only takes columns away, or fails and stops the change.
type tableState struct {
	x schemaIndex
	// moved is whether a USE has moved names without a database out of the
	// target's.
	moved bool
	// touched holds the synced tables a statement the rule does not follow
	// has named, and shadowed those a temporary table of their name hides.
	touched, shadowed map[*metadata.TableMetadata]bool
	// gone holds the columns DROP COLUMN statements have named, which the
	// table lacks after them.
	gone map[*metadata.TableMetadata][]string
}

func (w *tableState) table(ref *ast.TableRef) *metadata.TableMetadata {
	return w.x.table(ref, w.moved)
}

// known returns the synced table a reference means, and whether the rule
// knows it as the synced schema holds it, but for the columns DROP COLUMN
// took away.
func (w *tableState) known(ref *ast.TableRef) (*metadata.TableMetadata, bool) {
	t := w.table(ref)
	return t, t != nil && !w.touched[t]
}

// nameKey returns a reference's name as the target's database compares
// it, and whether the reference names that database.
func (w *tableState) nameKey(ref *ast.TableRef) (string, bool) {
	if ref == nil || ref.Schema == "" && w.moved || ref.Schema != "" && !w.x.same(ref.Schema, w.x.schema.GetName()) {
		return "", false
	}
	if w.x.caseSensitive {
		return ref.Name, true
	}
	return strings.ToLower(ref.Name), true
}

// lacks reports whether a known table lacks a column.
func (w *tableState) lacks(t *metadata.TableMetadata, column string) bool {
	return !slices.ContainsFunc(t.GetColumns(), func(c *metadata.ColumnMetadata) bool { return sameName(c.GetName(), column) }) ||
		slices.ContainsFunc(w.gone[t], func(c string) bool { return sameName(c, column) })
}

// dropsLacking reports whether subcommands of an ALTER TABLE of a known
// table drop a column it lacks, without IF EXISTS, which fails the
// statement whatever else it does.
func (w *tableState) dropsLacking(t *metadata.TableMetadata, cmds []*ast.AlterTableCmd) bool {
	return slices.ContainsFunc(cmds, func(cmd *ast.AlterTableCmd) bool {
		return cmd != nil && cmd.Type == ast.ATDropColumn && !cmd.IfExists && w.lacks(t, cmd.Name)
	})
}

// fails reports whether a statement the rule does not follow fails for
// certain: an ALTER TABLE of a known table that drops a column the table
// lacks, whatever else it does, or that renames the table onto another known
// table; and a RENAME TABLE with a pair onto a name a known table holds, or
// an earlier pair renamed a table onto, that no pair between renames away,
// the table itself among them. A temporary table of either name, which the
// statement may mean, leaves the statement unknown.
//
// mysql: https://dev.mysql.com/doc/refman/8.0/en/rename-table.html
func (w *tableState) fails(n ast.Node) bool {
	switch v := n.(type) {
	case *ast.AlterTableStmt:
		t, ok := w.known(v.Table)
		if !ok {
			return false
		}
		if w.dropsLacking(t, v.Commands) {
			return true
		}
		for _, cmd := range v.Commands {
			if cmd != nil && cmd.Type == ast.ATRenameTable {
				if d, ok := w.known(renamedTo(cmd)); ok && d != t {
					return true
				}
			}
		}
	case *ast.RenameTableStmt:
		// The pairs run in order: a name a known table holds, or one a pair
		// renamed a table onto, is taken until a pair renames it away.
		taken := make(map[string]bool)
		for _, pair := range v.Pairs {
			if pair == nil {
				continue
			}
			// MariaDB renames a temporary table of the name, which may take a
			// name a table holds.
			if o := w.table(pair.Old); o != nil && w.shadowed[o] {
				return false
			}
			d, inTarget := w.nameKey(pair.New)
			if inTarget {
				held, seen := taken[d]
				if !seen {
					_, held = w.known(pair.New)
				}
				if held {
					return true
				}
			}
			if o, ok := w.nameKey(pair.Old); ok {
				taken[o] = false
			}
			if inTarget {
				taken[d] = true
			}
		}
	}
	return false
}

// settles returns the synced tables a statement the rule does not follow may
// give a primary key back, rename away, or drop, whose findings it
// withdraws. A temporary table of the name takes an ALTER TABLE, and a DROP
// TABLE drops it before the synced table; RENAME TABLE renames the synced
// table, though a temporary table of its name exists.
func (w *tableState) settles(n ast.Node) []*metadata.TableMetadata {
	var out []*metadata.TableMetadata
	switch v := n.(type) {
	case *ast.AlterTableStmt:
		t := w.table(v.Table)
		if t != nil && !w.shadowed[t] && slices.ContainsFunc(v.Commands, func(cmd *ast.AlterTableCmd) bool {
			return cmd != nil && (cmd.Type == ast.ATRenameTable || addsPrimaryKey(cmd))
		}) {
			out = append(out, t)
		}
	case *ast.DropTableStmt:
		// DROP TEMPORARY TABLE drops no synced table.
		if v.Temporary {
			break
		}
		for _, ref := range v.Tables {
			if t := w.table(ref); t != nil && !w.shadowed[t] {
				out = append(out, t)
			}
		}
	case *ast.RenameTableStmt:
		for _, pair := range v.Pairs {
			if pair == nil {
				continue
			}
			if t := w.table(pair.Old); t != nil {
				out = append(out, t)
			}
		}
	}
	return out
}

// touch records the tables a statement the rule does not follow names,
// which it no longer knows, and the temporary tables of their names it
// makes, moves, or drops. A rename onto a known table fails, and one onto
// a table an earlier pair moved away names it as a source.
func (w *tableState) touch(n ast.Node) {
	mark := func(ref *ast.TableRef) *metadata.TableMetadata {
		t := w.table(ref)
		if t != nil {
			w.touched[t] = true
		}
		return t
	}
	switch v := n.(type) {
	case *ast.AlterTableStmt:
		// A rename moves a temporary table of the name away.
		if t := mark(v.Table); t != nil && slices.ContainsFunc(v.Commands, func(cmd *ast.AlterTableCmd) bool { return cmd != nil && cmd.Type == ast.ATRenameTable }) {
			w.shadowed[t] = false
		}
	case *ast.CreateTableStmt:
		if t := mark(v.Table); t != nil && v.Temporary {
			w.shadowed[t] = true
		}
	case *ast.CreateIndexStmt:
		mark(v.Table)
	case *ast.DropIndexStmt:
		mark(v.Table)
	case *ast.DropTableStmt:
		for _, ref := range v.Tables {
			if t := mark(ref); t != nil {
				w.shadowed[t] = false
			}
		}
	case *ast.RenameTableStmt:
		for _, pair := range v.Pairs {
			if pair != nil {
				mark(pair.Old)
			}
		}
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
// known, unless a unique key keeps them apart.
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
		if !keptApart(t, t.GetIndexes()[slices.IndexFunc(t.GetIndexes(), (*metadata.IndexMetadata).GetPrimary)], gone) {
			return nil, nil, false
		}
	default:
		lost = effective[slices.IndexFunc(effective, func(cmd *ast.AlterTableCmd) bool { return in(key, cmd.Name) })]
	}
	return columns, lost, true
}

// keptApart reports whether another unique key of the table keeps apart
// the values of a unique key, the primary key among them, once it narrows
// to its parts not gone: one whose every part is a column the narrowed key
// holds whole, which the narrowed key, losing a part, is not. Values a unique key keeps apart, even by prefixes, the
// whole columns keep apart too; a part whose length the synced schema does
// not give may be a prefix, which need not.
func keptApart(t *metadata.TableMetadata, narrowed *metadata.IndexMetadata, gone func(string) bool) bool {
	lengths := narrowed.GetKeyLength()
	holds := func(column string) bool {
		for n, part := range narrowed.GetExpressions() {
			if sameName(part, column) {
				return !gone(part) && n < len(lengths) && lengths[n] <= 0
			}
		}
		return false
	}
	return slices.ContainsFunc(t.GetIndexes(), func(idx *metadata.IndexMetadata) bool {
		return idx.GetUnique() && len(idx.GetExpressions()) > 0 &&
			!slices.ContainsFunc(idx.GetExpressions(), func(part string) bool { return !holds(part) })
	})
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
// may read it, or whose expression default may; a check that reads it and
// another column; a unique key of it and a part not gone, which narrows to
// the rest and fails when their values repeat, unless another unique key
// keeps them apart; a functional index part that may read it; or the table's
// partitioning or subpartitioning, by an expression that may read it or by
// KEY (), which reads the primary key, or a unique key of a table without
// one. A check that reads the column alone goes with it. A foreign key names
// the table it references without its database, so one that may reference a
// table of another database counts.
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
		if idx.GetUnique() && !idx.GetPrimary() && has(idx.GetExpressions()) && slices.ContainsFunc(idx.GetExpressions(), func(part string) bool { return !gone(part) }) && !keptApart(t, idx, gone) {
			return true
		}
		for _, part := range idx.GetExpressions() {
			if !sameName(part, column) && reads(part, column) {
				return true
			}
		}
	}
	// KEY () is synced with no expression, and reads the primary key, or,
	// without one, a unique key: a column of none of them is not read.
	partitionedBy := func(expr string) bool {
		if expr != "" {
			return reads(expr, column)
		}
		if primaryKeyColumns(t) != nil {
			// A part of the key that is a prefix of its column is left out.
			//
			// mysql: https://dev.mysql.com/doc/refman/8.0/en/partitioning-limitations.html
			pk := t.GetIndexes()[slices.IndexFunc(t.GetIndexes(), (*metadata.IndexMetadata).GetPrimary)]
			lengths := pk.GetKeyLength()
			for n, part := range pk.GetExpressions() {
				if sameName(part, column) && (n >= len(lengths) || lengths[n] <= 0) {
					return true
				}
			}
			return false
		}
		// The partitioning columns are in every unique key.
		//
		// mysql: https://dev.mysql.com/doc/refman/8.0/en/partitioning-limitations-partitioning-keys-unique-keys.html
		return !slices.ContainsFunc(t.GetIndexes(), func(idx *metadata.IndexMetadata) bool { return idx.GetUnique() && !has(idx.GetExpressions()) })
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
