package catalog

import (
	nodes "github.com/bytebase/omni/mysql/ast"
	mysqlparser "github.com/bytebase/omni/mysql/parser"
)

type ExecOptions struct {
	ContinueOnError bool
}

type ExecResult struct {
	Index   int
	SQL     string
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

	continueOnError := false
	if opts != nil {
		continueOnError = opts.ContinueOnError
	}

	results := make([]ExecResult, 0, len(list.Items))
	for i, item := range list.Items {
		result := ExecResult{Index: i}

		if isDML(item) {
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

func isDML(stmt nodes.Node) bool {
	switch stmt.(type) {
	case *nodes.SelectStmt, *nodes.InsertStmt, *nodes.UpdateStmt, *nodes.DeleteStmt:
		return true
	default:
		return false
	}
}

func (c *Catalog) processUtility(stmt nodes.Node) error {
	switch s := stmt.(type) {
	case *nodes.CreateDatabaseStmt:
		return c.createDatabase(s)
	case *nodes.CreateTableStmt:
		return c.createTable(s)
	case *nodes.CreateIndexStmt:
		return c.createIndex(s)
	case *nodes.CreateViewStmt:
		return c.createView(s)
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
	default:
		return nil
	}
}
