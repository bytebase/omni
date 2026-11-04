package catalog

import (
	"fmt"

	nodes "github.com/bytebase/omni/tidb/ast"
	mysqlparser "github.com/bytebase/omni/tidb/parser"
)

type ExecOptions struct {
	ContinueOnError bool
}

type ExecResult struct {
	Index   int
	SQL     string
	Line    int // 1-based start line in the original SQL
	Skipped bool
	Error   error
}

func (c *Catalog) Exec(sql string, opts *ExecOptions) ([]ExecResult, error) {
	list, err := mysqlparser.Parse(sql)
	if err != nil {
		return nil, err
	}
	if list == nil || len(list.Items) == 0 {
		return nil, nil
	}

	lineIndex := buildLineIndex(sql)

	continueOnError := false
	if opts != nil {
		continueOnError = opts.ContinueOnError
	}

	results := make([]ExecResult, 0, len(list.Items))
	for i, item := range list.Items {
		locStart := stmtLocStart(item)
		result := ExecResult{
			Index: i,
			Line:  offsetToLine(lineIndex, locStart),
		}

		if isDML(item) {
			result.Skipped = true
			results = append(results, result)
			continue
		}

		execErr := c.processUtility(item, sql)
		result.Error = execErr
		results = append(results, result)

		if execErr != nil && !continueOnError {
			break
		}
	}
	return results, nil
}

func LoadSQL(sql string) (*Catalog, error) {
	c := New()
	results, err := c.Exec(sql, nil)
	if err != nil {
		return nil, err
	}
	for _, r := range results {
		if r.Error != nil {
			return c, r.Error
		}
	}
	return c, nil
}

// execSet handles SET statements that affect catalog behavior.
// Most SET variables are silently accepted (session-level settings like NAMES,
// CHARACTER SET, sql_mode). Variables that affect DDL behavior (foreign_key_checks)
// update the catalog state.
func (c *Catalog) execSet(stmt *nodes.SetStmt) error {
	for _, asgn := range stmt.Assignments {
		varName := toLower(asgn.Column.Column)
		switch varName {
		case "foreign_key_checks":
			// Extract the value.
			val := nodeToSQLValue(asgn.Value)
			switch toLower(val) {
			case "0", "off", "false":
				c.foreignKeyChecks = false
			case "1", "on", "true":
				c.foreignKeyChecks = true
			}
		case "names", "character set":
			// Silently accept — these affect character encoding but
			// the in-memory catalog doesn't need to change behavior.
		default:
			// Silently accept all other SET variables (sql_mode, etc.).
		}
	}
	return nil
}

// nodeToSQLValue extracts a simple string value from an expression node.
func nodeToSQLValue(expr nodes.ExprNode) string {
	switch e := expr.(type) {
	case *nodes.StringLit:
		return e.Value
	case *nodes.IntLit:
		return fmt.Sprintf("%d", e.Value)
	case *nodes.FloatLit:
		return e.Value
	case *nodes.BoolLit:
		if e.Value {
			return "1"
		}
		return "0"
	case *nodes.ColumnRef:
		return e.Column
	default:
		return ""
	}
}

func isDML(stmt nodes.Node) bool {
	switch stmt.(type) {
	case *nodes.SelectStmt, *nodes.InsertStmt, *nodes.UpdateStmt, *nodes.DeleteStmt:
		return true
	default:
		return false
	}
}

// processUtility applies one statement. source is the text it was parsed
// from, which its Loc offsets index into, or "" when there is none.
func (c *Catalog) processUtility(stmt nodes.Node, source string) error {
	switch s := stmt.(type) {
	case *nodes.CreateDatabaseStmt:
		return c.createDatabase(s)
	case *nodes.CreateTableStmt:
		return c.createTable(s, true)
	case *nodes.CreateIndexStmt:
		return c.createIndex(s)
	case *nodes.CreateViewStmt:
		return c.createView(s, source)
	case *nodes.AlterViewStmt:
		return c.alterView(s, source)
	case *nodes.AlterTableStmt:
		return c.alterTable(s)
	case *nodes.AlterDatabaseStmt:
		return c.alterDatabase(s)
	case *nodes.CreatePlacementPolicyStmt:
		return c.createPlacementPolicy(s)
	case *nodes.AlterPlacementPolicyStmt:
		return c.alterPlacementPolicy(s)
	case *nodes.DropPlacementPolicyStmt:
		return c.dropPlacementPolicy(s)
	case *nodes.DropTableStmt:
		return c.dropTable(s)
	case *nodes.DropDatabaseStmt:
		return c.dropDatabase(s)
	case *nodes.DropIndexStmt:
		return c.dropIndex(s)
	case *nodes.DropViewStmt:
		return c.dropView(s)
	case *nodes.RenameTableStmt:
		return c.renameTable(s)
	case *nodes.TruncateStmt:
		return c.truncateTable(s)
	case *nodes.UseStmt:
		return c.useDatabase(s)
	case *nodes.CreateFunctionStmt:
		return c.createRoutine(s)
	case *nodes.DropRoutineStmt:
		return c.dropRoutine(s)
	case *nodes.AlterRoutineStmt:
		return c.alterRoutine(s)
	case *nodes.CreateTriggerStmt:
		return c.createTrigger(s)
	case *nodes.DropTriggerStmt:
		return c.dropTrigger(s)
	case *nodes.CreateEventStmt:
		return c.createEvent(s)
	case *nodes.AlterEventStmt:
		return c.alterEvent(s)
	case *nodes.DropEventStmt:
		return c.dropEvent(s)
	case *nodes.SetStmt:
		return c.execSet(s)
	case *nodes.CreateBindingStmt, *nodes.DropBindingStmt, *nodes.SetBindingStmt:
		// A SQL plan binding records a plan for later statements and
		// changes no schema object. The statements it carries are not run
		// either, so the catalog reads neither of them.
		return nil
	case *nodes.DropStatsStmt, *nodes.LoadStatsStmt, *nodes.CreateStatisticsStmt, *nodes.DropStatisticsStmt:
		// TiDB's standalone statistics statements. DROP STATS discards a
		// table's statistics, LOAD STATS reads a dump back, and CREATE and
		// DROP STATISTICS make and drop an extended statistic. None
		// changes a schema the catalog holds -- no table, column, index,
		// or option -- so each is applied as a no-op rather than left to
		// the default arm, where an unmodelled statement and a deliberate
		// no-op read the same. ANALYZE TABLE reaches the default arm for
		// the same reason and has always been a no-op there.
		return nil
	case *nodes.LoadDataStmt, *nodes.ImportIntoStmt:
		// The bulk-load statements. Both write rows into a table that
		// exists already and neither changes a schema the catalog holds --
		// no table, column, index, or option -- so each is applied as a
		// no-op rather than left to the default arm, where an unmodelled
		// statement and a deliberate no-op read the same. A change script
		// may hold either next to its DDL.
		return nil
	case *nodes.CreateUserStmt, *nodes.AlterUserStmt, *nodes.DropUserStmt,
		*nodes.RenameUserStmt, *nodes.CreateRoleStmt, *nodes.DropRoleStmt,
		*nodes.SetPasswordStmt, *nodes.GrantStmt, *nodes.RevokeStmt,
		*nodes.GrantRoleStmt, *nodes.RevokeRoleStmt, *nodes.SetRoleStmt,
		*nodes.SetDefaultRoleStmt:
		// The user and privilege statements. A change script holds them
		// next to its DDL -- a GRANT on the table a migration just made,
		// a CREATE USER for the service that will read it -- and none of
		// them changes a schema the catalog holds: it models no user, no
		// role and no grant. Each is applied as a no-op rather than left
		// to the default arm, where an unmodelled statement and a
		// deliberate no-op read the same.
		return nil
	case *nodes.SplitRegionStmt, *nodes.DistributeTableStmt:
		// SPLIT ... TABLE pre-splits regions and DISTRIBUTE TABLE
		// schedules replicas. Neither changes a schema the catalog
		// holds -- no table, column, index, or option -- so both are
		// applied as a no-op rather than left to the default arm, where
		// an unmodelled statement and a deliberate no-op read the same.
		return nil
	case *nodes.BRIEStmt, *nodes.RecoverTableStmt, *nodes.FlashbackTableStmt,
		*nodes.FlashbackDatabaseStmt, *nodes.FlashbackToTimestampStmt:
		// TiDB's backup, restore, recover and flashback statements. A
		// BACKUP reads and changes nothing. The others do change a
		// schema, but none of them says how: RESTORE, RECOVER TABLE and
		// the two bring-back FLASHBACK forms each put back an object
		// whose definition is in a backup or in a dropped table the
		// catalog keeps no record of, and FLASHBACK ... TO TIMESTAMP
		// rolls data back to a point the catalog never saw. So the
		// catalog has nothing to apply and applies nothing, rather than
		// guess at a definition it does not hold.
		//
		// In particular FLASHBACK TABLE t TO t1 is not a rename, and is
		// deliberately not read as one: t has to be gone already for the
		// statement to succeed, and t1 is created. Its sibling makes
		// that plain -- v8.5.5 answers 1007, "Can't create database 'd';
		// database exists", for FLASHBACK DATABASE d when d is still
		// there.
		//
		// Each is dispatched here rather than left to the default arm,
		// where an unmodelled statement and a deliberate no-op read the
		// same.
		return nil
	case *nodes.ExplainStmt:
		// EXPLAIN ANALYZE runs the statement it carries, and on TiDB that
		// reaches ALTER TABLE: EXPLAIN ANALYZE ALTER TABLE t ADD COLUMN c
		// adds the column (tidb/parser/explain_forms_test.go measures it
		// against v8.5.5), so the catalog applies what it carries. A plain
		// EXPLAIN reports a plan and changes nothing.
		if s.Analyze && s.Stmt != nil {
			return c.processUtility(s.Stmt, source)
		}
		return nil
	default:
		return nil
	}
}

// buildLineIndex returns the byte offset of each line start.
func buildLineIndex(sql string) []int {
	index := []int{0}
	for i := 0; i < len(sql); i++ {
		if sql[i] == '\n' {
			index = append(index, i+1)
		}
	}
	return index
}

// offsetToLine converts a byte offset to a 1-based line number.
func offsetToLine(lineIndex []int, offset int) int {
	lo, hi := 0, len(lineIndex)-1
	for lo < hi {
		mid := (lo + hi + 1) / 2
		if lineIndex[mid] <= offset {
			lo = mid
		} else {
			hi = mid - 1
		}
	}
	return lo + 1
}

// stmtLocStart extracts Loc.Start from a statement node.
// Uses a type switch over the statement types handled by processUtility,
// plus common DML types. All have a Loc field set by the parser.
func stmtLocStart(node nodes.Node) int {
	switch s := node.(type) {
	case *nodes.CreateDatabaseStmt:
		return s.Loc.Start
	case *nodes.CreateTableStmt:
		return s.Loc.Start
	case *nodes.CreateIndexStmt:
		return s.Loc.Start
	case *nodes.CreateViewStmt:
		return s.Loc.Start
	case *nodes.AlterViewStmt:
		return s.Loc.Start
	case *nodes.AlterTableStmt:
		return s.Loc.Start
	case *nodes.AlterDatabaseStmt:
		return s.Loc.Start
	case *nodes.CreatePlacementPolicyStmt:
		return s.Loc.Start
	case *nodes.AlterPlacementPolicyStmt:
		return s.Loc.Start
	case *nodes.DropPlacementPolicyStmt:
		return s.Loc.Start
	case *nodes.CreateSequenceStmt:
		return s.Loc.Start
	case *nodes.AlterSequenceStmt:
		return s.Loc.Start
	case *nodes.DropSequenceStmt:
		return s.Loc.Start
	case *nodes.DropTableStmt:
		return s.Loc.Start
	case *nodes.DropDatabaseStmt:
		return s.Loc.Start
	case *nodes.DropIndexStmt:
		return s.Loc.Start
	case *nodes.DropViewStmt:
		return s.Loc.Start
	case *nodes.RenameTableStmt:
		return s.Loc.Start
	case *nodes.TruncateStmt:
		return s.Loc.Start
	case *nodes.UseStmt:
		return s.Loc.Start
	case *nodes.CreateFunctionStmt:
		return s.Loc.Start
	case *nodes.DropRoutineStmt:
		return s.Loc.Start
	case *nodes.AlterRoutineStmt:
		return s.Loc.Start
	case *nodes.CreateTriggerStmt:
		return s.Loc.Start
	case *nodes.DropTriggerStmt:
		return s.Loc.Start
	case *nodes.CreateEventStmt:
		return s.Loc.Start
	case *nodes.AlterEventStmt:
		return s.Loc.Start
	case *nodes.DropEventStmt:
		return s.Loc.Start
	case *nodes.SetStmt:
		return s.Loc.Start
	case *nodes.CreateBindingStmt:
		return s.Loc.Start
	case *nodes.DropBindingStmt:
		return s.Loc.Start
	case *nodes.SetBindingStmt:
		return s.Loc.Start
	case *nodes.CreateResourceGroupStmt:
		return s.Loc.Start
	case *nodes.AlterResourceGroupStmt:
		return s.Loc.Start
	case *nodes.DropResourceGroupStmt:
		return s.Loc.Start
	case *nodes.SetResourceGroupStmt:
		return s.Loc.Start
	case *nodes.CalibrateResourceStmt:
		return s.Loc.Start
	case *nodes.AddQueryWatchStmt:
		return s.Loc.Start
	case *nodes.DropQueryWatchStmt:
		return s.Loc.Start
	case *nodes.CreateUserStmt:
		return s.Loc.Start
	case *nodes.AlterUserStmt:
		return s.Loc.Start
	case *nodes.DropUserStmt:
		return s.Loc.Start
	case *nodes.RenameUserStmt:
		return s.Loc.Start
	case *nodes.CreateRoleStmt:
		return s.Loc.Start
	case *nodes.DropRoleStmt:
		return s.Loc.Start
	case *nodes.SetPasswordStmt:
		return s.Loc.Start
	case *nodes.GrantStmt:
		return s.Loc.Start
	case *nodes.RevokeStmt:
		return s.Loc.Start
	case *nodes.GrantRoleStmt:
		return s.Loc.Start
	case *nodes.RevokeRoleStmt:
		return s.Loc.Start
	case *nodes.SetRoleStmt:
		return s.Loc.Start
	case *nodes.SetDefaultRoleStmt:
		return s.Loc.Start
	case *nodes.ExplainStmt:
		// The whole EXPLAIN, not the statement it carries: an ExecResult's
		// line is the line the caller wrote.
		return s.Loc.Start
	case *nodes.LoadDataStmt:
		return s.Loc.Start
	case *nodes.ImportIntoStmt:
		return s.Loc.Start
	case *nodes.SelectStmt:
		return s.Loc.Start
	case *nodes.InsertStmt:
		return s.Loc.Start
	case *nodes.UpdateStmt:
		return s.Loc.Start
	case *nodes.DeleteStmt:
		return s.Loc.Start
	default:
		return 0
	}
}
