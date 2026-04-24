package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/mysql/ast"
	mysqlparser "github.com/bytebase/omni/mysql/parser"
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
			if err := validateExpressionSemantics(item); err != nil {
				result.Error = err
				results = append(results, result)
				if !continueOnError {
					break
				}
				continue
			}
			result.Skipped = true
			results = append(results, result)
			continue
		}

		execErr := c.processUtility(item)
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
			charset := normalizeCharsetName(nodeToSQLValue(asgn.Value))
			if strings.EqualFold(charset, "DEFAULT") || charset == "" {
				charset = c.defaultCharset
			}
			c.charsetClient = charset
			if coll, ok := defaultCollationForCharset[toLower(charset)]; ok {
				c.collationConn = coll
			}
		case "collate":
			collation := nodeToSQLValue(asgn.Value)
			if collation != "" {
				c.collationConn = collation
			}
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

func (c *Catalog) processUtility(stmt nodes.Node) error {
	if err := validateExpressionSemantics(stmt); err != nil {
		return err
	}
	switch s := stmt.(type) {
	case *nodes.CreateDatabaseStmt:
		return c.createDatabase(s)
	case *nodes.CreateTableStmt:
		return c.createTable(s)
	case *nodes.CreateIndexStmt:
		return c.createIndex(s)
	case *nodes.CreateViewStmt:
		return c.createView(s)
	case *nodes.AlterViewStmt:
		return c.alterView(s)
	case *nodes.AlterTableStmt:
		return c.alterTable(s)
	case *nodes.AlterDatabaseStmt:
		return c.alterDatabase(s)
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
