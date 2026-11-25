package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/mariadb/ast"
	mysqlparser "github.com/bytebase/omni/mariadb/parser"
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

		if err := validateDDL(item); err != nil {
			result.Error = err
			results = append(results, result)
			if !continueOnError {
				break
			}
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

// validateDDL refuses a DDL statement the server would not take that the parser
// does. A `?` parameter marker binds a value only in a prepared statement, so
// it is never valid in catalog DDL, though the parser accepts it as a general
// operand (a default, generated-column or CHECK expression, a view body). Parse
// has already bounded the tree height, so the walk is safe.
func validateDDL(stmt nodes.Node) error {
	if containsParamMarker(stmt) {
		return &Error{
			Code:     1064,
			SQLState: "42000",
			Message:  "parameter marker (?) is not allowed in a DDL statement",
		}
	}
	return validateParseTimeLookups(stmt)
}

// validateParseTimeLookups refuses the checks MariaDB makes while it parses that
// answer with the code of the name they looked up, before anything of the statement
// runs: a column whose type writes a character set and a collation of another one is
// 1253, and a FULLTEXT key's WITH PARSER names a parser plugin the pinned server does
// not ship, which is 1128. The server stops at the first of them, so the one earliest
// in the text is the one reported: `ADD a CHAR ASCII COLLATE utf8mb4_bin, ADD FULLTEXT
// (a) WITH PARSER p` is 1253.
//
// Measured on mariadb:11.8.8: both answer before the table is looked up (1146) and
// before a name collides (1050) or a column repeats (1060), so this runs before the
// statement does. The parser records the same two checks for a statement that does
// not parse (parser.ParseError.Code), which never reaches here.
func validateParseTimeLookups(stmt nodes.Node) error {
	var err error
	first := -1
	record := func(pos int, e error) {
		if first < 0 || pos < first {
			first, err = pos, e
		}
	}
	nodes.Inspect(stmt, func(n nodes.Node) bool {
		switch x := n.(type) {
		case *nodes.ColumnDef:
			dt := x.TypeName
			collate := declaredColumnCollate(x)
			if dt == nil || dt.Charset == "" || collate == "" || charsetMatchesCollation(dt.Charset, collate) {
				return true
			}
			charset := normalizeCharsetName(dt.Charset)
			record(dt.Loc.Start, &Error{
				Code:     1253,
				SQLState: sqlState(1253),
				Message:  fmt.Sprintf("COLLATION '%s' is not valid for CHARACTER SET '%s'", strings.ToLower(collate), charset),
			})
		case *nodes.CreateIndexStmt:
			if x.Fulltext {
				recordFulltextParser(record, x.Options)
			}
		case *nodes.Constraint:
			if x.Type == nodes.ConstrFulltextIndex {
				recordFulltextParser(record, x.IndexOptions)
			}
		}
		return true
	})
	return err
}

// recordFulltextParser records the 1128 a FULLTEXT key's WITH PARSER answers for its
// parser plugin, which is the only option whose name the server looks up.
func recordFulltextParser(record func(int, error), opts []*nodes.IndexOption) {
	for _, opt := range opts {
		if opt == nil || !strings.EqualFold(opt.Name, "PARSER") {
			continue
		}
		name := ""
		if s, ok := opt.Value.(*nodes.StringLit); ok {
			name = s.Value
		}
		record(opt.Loc.Start, errFunctionNotDefined(name))
		return
	}
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
		varName := normalizeSetVariableName(asgn.Column)
		switch varName {
		case "foreign_key_checks":
			if v, ok := parseSetBool(asgn.Value); ok {
				c.foreignKeyChecks = v
				c.session.ForeignKeyChecks = v
			}
		case "explicit_defaults_for_timestamp":
			if isDefaultSetValue(asgn.Value) {
				c.session.ExplicitDefaultsForTimestamp = true
			} else if v, ok := parseSetBool(asgn.Value); ok {
				c.session.ExplicitDefaultsForTimestamp = v
			}
		case "sql_mode":
			c.session.SQLMode = nodeToSQLValue(asgn.Value)
			if c.session.SQLMode == "" && isDefaultSetValue(asgn.Value) {
				c.session.SQLMode = "DEFAULT"
			}
		case "names", "character set":
			charset := normalizeCharsetName(nodeToSQLValue(asgn.Value))
			if strings.EqualFold(charset, "DEFAULT") || charset == "" {
				charset = c.defaultCharset
			}
			c.charsetClient = charset
			c.session.CharsetClient = charset
			if coll := c.charsetCollation(charset); coll != "" {
				c.collationConn = coll
				c.session.CollationConnection = coll
			}
		case "collate":
			collation := nodeToSQLValue(asgn.Value)
			if collation != "" {
				c.collationConn = collation
				c.session.CollationConnection = collation
			}
		default:
			// Silently accept all other SET variables (sql_mode, etc.).
		}
	}
	return nil
}

// execSetStatement applies SET STATEMENT: the statement runs with the options set, and each
// option's variable then goes back to what it was, even if the statement set it too. A variable
// the statement sets that no option names keeps its new value. Measured on mariadb:11.8.8.
func (c *Catalog) execSetStatement(s *nodes.SetStatementStmt, source string) error {
	fkChecks, session := c.foreignKeyChecks, c.session
	if err := c.execSet(&nodes.SetStmt{Assignments: s.Options}); err != nil {
		return err
	}
	var err error
	if isDML(s.Stmt) {
		err = validateExpressionSemantics(s.Stmt)
	} else {
		err = c.processUtility(s.Stmt, source)
	}
	for _, opt := range s.Options {
		switch normalizeSetVariableName(opt.Column) {
		case "foreign_key_checks":
			c.foreignKeyChecks, c.session.ForeignKeyChecks = fkChecks, session.ForeignKeyChecks
		case "explicit_defaults_for_timestamp":
			c.session.ExplicitDefaultsForTimestamp = session.ExplicitDefaultsForTimestamp
		case "sql_mode":
			c.session.SQLMode = session.SQLMode
		}
	}
	return err
}

func normalizeSetVariableName(col *nodes.ColumnRef) string {
	if col == nil {
		return ""
	}
	name := col.Column
	if col.Table != "" {
		name = col.Table + "." + name
	}
	name = strings.TrimPrefix(name, "@@")
	name = toLower(name)
	for _, prefix := range []string{"session.", "global.", "local.", "persist.", "persist_only."} {
		name = strings.TrimPrefix(name, prefix)
	}
	return name
}

func parseSetBool(expr nodes.ExprNode) (bool, bool) {
	switch toLower(nodeToSQLValue(expr)) {
	case "0", "off", "false":
		return false, true
	case "1", "on", "true":
		return true, true
	default:
		return false, false
	}
}

func isDefaultSetValue(expr nodes.ExprNode) bool {
	return strings.EqualFold(nodeToSQLValue(expr), "DEFAULT")
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
	case *nodes.DefaultExpr:
		return "DEFAULT"
	default:
		return ""
	}
}

func isDML(stmt nodes.Node) bool {
	switch s := stmt.(type) {
	case *nodes.SelectStmt, *nodes.InsertStmt, *nodes.UpdateStmt, *nodes.DeleteStmt:
		return true
	case *nodes.SetStatementStmt:
		return isDML(s.Stmt)
	default:
		return false
	}
}

// processUtility applies one statement. source is the text it was parsed
// from, which its Loc offsets index into, or "" when there is none.
func (c *Catalog) processUtility(stmt nodes.Node, source string) error {
	if !c.loadingSDL {
		if err := checkProgramSystemVariables(stmt); err != nil {
			return err
		}
	}
	if err := validateExpressionSemantics(stmt); err != nil {
		return err
	}
	switch s := stmt.(type) {
	case *nodes.CreateDatabaseStmt:
		return c.createDatabase(s)
	case *nodes.CreateTableStmt:
		return c.createTable(s, source)
	case *nodes.CreateIndexStmt:
		return c.createIndex(s)
	case *nodes.CreateViewStmt:
		return c.createView(s, source, true)
	case *nodes.AlterViewStmt:
		return c.alterView(s, source)
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
		return c.createTrigger(s, source)
	case *nodes.DropTriggerStmt:
		return c.dropTrigger(s)
	case *nodes.CreateEventStmt:
		return c.createEvent(s, true)
	case *nodes.AlterEventStmt:
		return c.alterEvent(s)
	case *nodes.DropEventStmt:
		return c.dropEvent(s)
	case *nodes.CreateSequenceStmt:
		return c.createSequence(s)
	case *nodes.AlterSequenceStmt:
		return c.alterSequence(s)
	case *nodes.DropSequenceStmt:
		return c.dropSequence(s)
	case *nodes.SetStmt:
		return c.execSet(s)
	case *nodes.SetStatementStmt:
		return c.execSetStatement(s, source)
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
	case *nodes.CreateSequenceStmt:
		return s.Loc.Start
	case *nodes.AlterSequenceStmt:
		return s.Loc.Start
	case *nodes.DropSequenceStmt:
		return s.Loc.Start
	case *nodes.SetStmt:
		return s.Loc.Start
	case *nodes.SetStatementStmt:
		return s.Loc.Start
	case *nodes.SelectStmt:
		return s.Loc.Start
	case *nodes.InsertStmt:
		return s.Loc.Start
	case *nodes.UpdateStmt:
		return s.Loc.Start
	case *nodes.DeleteStmt:
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
	default:
		return 0
	}
}
