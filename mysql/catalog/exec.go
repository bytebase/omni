package catalog

import (
	"fmt"
	"strconv"
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

		// A `?` parameter marker binds a value only in a prepared statement, so
		// it is never valid in catalog DDL. The parser accepts it as a general
		// operand (a default, generated-column or CHECK expression, a view
		// body), so refuse it here before processUtility stores the expression.
		// Parse has already bounded the tree height, so the walk is safe.
		if containsParamMarker(item) {
			result.Error = &Error{
				Code:     0,
				SQLState: "HY000",
				Message:  "parameter marker (?) is not allowed in a DDL statement",
			}
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
	// MySQL evaluates and checks every assignment in a SET statement before
	// applying any variable update, so a failure rolls back the whole SET and an
	// RHS reading @var sees the value from before this statement.
	var updates []func()
	for _, asgn := range stmt.Assignments {
		if asgn.Column != nil && strings.HasPrefix(asgn.Column.Column, "@") && !strings.HasPrefix(asgn.Column.Column, "@@") {
			uName := toLower(strings.TrimPrefix(asgn.Column.Column, "@"))
			val, ok := c.evalSetExpr(asgn.Value)
			updates = append(updates, func() {
				if ok {
					c.userVars[uName] = val
				} else {
					delete(c.userVars, uName)
				}
			})
			continue
		}
		varName := normalizeSetVariableName(asgn.Column)
		switch varName {
		case "foreign_key_checks":
			if v, ok := parseSetBool(asgn.Value); ok {
				updates = append(updates, func() {
					c.foreignKeyChecks = v
					c.session.ForeignKeyChecks = v
				})
			}
		case "sql_generate_invisible_primary_key":
			if v, ok := parseSetBool(asgn.Value); ok {
				updates = append(updates, func() {
					c.generateGIPK = v
					c.session.GenerateInvisiblePrimaryKey = v
				})
			}
		case "show_gipk_in_create_table_and_information_schema":
			if v, ok := parseSetBool(asgn.Value); ok {
				updates = append(updates, func() {
					c.showGIPK = v
					c.session.ShowGIPK = v
				})
			}
		case "explicit_defaults_for_timestamp":
			if isDefaultSetKeywordOrIdent(asgn.Value) {
				updates = append(updates, func() {
					c.session.ExplicitDefaultsForTimestamp = true
				})
			} else if v, ok := parseSetBool(asgn.Value); ok {
				updates = append(updates, func() {
					c.session.ExplicitDefaultsForTimestamp = v
				})
			}
		case "innodb_directories":
			if v, ok := c.evalSetExpr(asgn.Value); ok && v.kind == userVarString {
				dirList := v.str
				updates = append(updates, func() {
					c.AddTablespaceDir(dirList)
				})
			}
		case "sql_mode":
			isGlobal, isPersistOnly := setAssignmentIsGlobalOrPersist(asgn)
			if isDefaultSetKeywordOrIdent(asgn.Value) {
				if isPersistOnly {
					continue
				}
				if isGlobal {
					updates = append(updates, func() {
						c.globalSQLMode = "DEFAULT"
					})
				} else {
					defMode := c.globalSQLMode
					if defMode == "" {
						defMode = "DEFAULT"
					}
					updates = append(updates, func() {
						c.session.SQLMode = defMode
					})
				}
				continue
			}
			var val userVarValue
			if cr, isCol := asgn.Value.(*nodes.ColumnRef); isCol && cr.Schema == "" && cr.Table == "" {
				val = userVarValue{kind: userVarString, str: cr.Column}
			} else {
				var ok bool
				val, ok = c.evalSetExpr(asgn.Value)
				if !ok {
					return errWrongValueForVar("sql_mode", strings.TrimSpace(nodeToSQL(asgn.Value)))
				}
			}
			mode, err := validateAndNormalizeSQLMode(val)
			if err != nil {
				return err
			}
			if isPersistOnly {
				continue
			}
			if isGlobal {
				updates = append(updates, func() {
					c.globalSQLMode = mode
				})
			} else {
				updates = append(updates, func() {
					c.session.SQLMode = mode
				})
			}
		case "time_zone":
			isGlobal, isPersistOnly := setAssignmentIsGlobalOrPersist(asgn)
			if isDefaultSetKeywordOrIdent(asgn.Value) {
				if isPersistOnly {
					continue
				}
				if isGlobal {
					updates = append(updates, func() {
						c.globalTimeZone = "SYSTEM"
					})
				} else {
					defTZ := c.globalTimeZone
					if defTZ == "" {
						defTZ = "SYSTEM"
					}
					updates = append(updates, func() {
						c.session.TimeZone = defTZ
					})
				}
				continue
			}
			var rawTZ string
			if cr, isCol := asgn.Value.(*nodes.ColumnRef); isCol && cr.Schema == "" && cr.Table == "" {
				rawTZ = cr.Column
			} else {
				val, ok := c.evalSetExpr(asgn.Value)
				if !ok {
					return errUnknownTimeZone(strings.TrimSpace(nodeToSQL(asgn.Value)))
				}
				switch val.kind {
				case userVarNull:
					return errWrongValueForVar("time_zone", "NULL")
				case userVarNumeric:
					return errWrongTypeForVar("time_zone")
				case userVarString:
					rawTZ = val.str
				}
			}
			tz, err := validateAndNormalizeSessionTimeZone(rawTZ)
			if err != nil {
				return err
			}
			if isPersistOnly {
				continue
			}
			if isGlobal {
				updates = append(updates, func() {
					c.globalTimeZone = tz
				})
			} else {
				updates = append(updates, func() {
					c.session.TimeZone = tz
				})
			}
		case "names", "character set":
			charset := normalizeCharsetName(nodeToSQLValue(asgn.Value))
			if strings.EqualFold(charset, "DEFAULT") || charset == "" {
				charset = c.defaultCharset
			}
			updates = append(updates, func() {
				c.charsetClient = charset
				c.session.CharsetClient = charset
				if coll, ok := defaultCollationForCharset[toLower(charset)]; ok {
					c.collationConn = coll
					c.session.CollationConnection = coll
				}
			})
		case "collate":
			collation := nodeToSQLValue(asgn.Value)
			if collation != "" {
				updates = append(updates, func() {
					c.collationConn = collation
					c.session.CollationConnection = collation
				})
			}
		default:
			// Silently accept all other SET variables.
		}
	}
	for _, fn := range updates {
		fn()
	}
	return nil
}

func setAssignmentIsGlobalOrPersist(asgn *nodes.Assignment) (isGlobal, isPersistOnly bool) {
	if asgn == nil {
		return false, false
	}
	switch strings.ToUpper(asgn.Scope) {
	case "PERSIST_ONLY":
		return true, true
	case "GLOBAL", "PERSIST":
		return true, false
	}
	if asgn.Column != nil {
		col := toLower(strings.TrimPrefix(asgn.Column.Column, "@@"))
		if strings.HasPrefix(col, "persist_only.") {
			return true, true
		}
		if strings.HasPrefix(col, "global.") || strings.HasPrefix(col, "persist.") {
			return true, false
		}
	}
	return false, false
}

var validMySQL80SQLModes = map[string]string{
	"REAL_AS_FLOAT":              "REAL_AS_FLOAT",
	"PIPES_AS_CONCAT":            "PIPES_AS_CONCAT",
	"ANSI_QUOTES":                "ANSI_QUOTES",
	"IGNORE_SPACE":               "IGNORE_SPACE",
	"NOT_USED":                   "NOT_USED",
	"ONLY_FULL_GROUP_BY":         "ONLY_FULL_GROUP_BY",
	"NO_UNSIGNED_SUBTRACTION":    "NO_UNSIGNED_SUBTRACTION",
	"NO_DIR_IN_CREATE":           "NO_DIR_IN_CREATE",
	"ANSI":                       "ANSI",
	"NO_AUTO_VALUE_ON_ZERO":      "NO_AUTO_VALUE_ON_ZERO",
	"NO_BACKSLASH_ESCAPES":       "NO_BACKSLASH_ESCAPES",
	"STRICT_TRANS_TABLES":        "STRICT_TRANS_TABLES",
	"STRICT_ALL_TABLES":          "STRICT_ALL_TABLES",
	"NO_ZERO_IN_DATE":            "NO_ZERO_IN_DATE",
	"NO_ZERO_DATE":               "NO_ZERO_DATE",
	"ALLOW_INVALID_DATES":        "ALLOW_INVALID_DATES",
	"ERROR_FOR_DIVISION_BY_ZERO": "ERROR_FOR_DIVISION_BY_ZERO",
	"TRADITIONAL":                "TRADITIONAL",
	"HIGH_NOT_PRECEDENCE":        "HIGH_NOT_PRECEDENCE",
	"NO_ENGINE_SUBSTITUTION":     "NO_ENGINE_SUBSTITUTION",
	"PAD_CHAR_TO_FULL_LENGTH":    "PAD_CHAR_TO_FULL_LENGTH",
	"TIME_TRUNCATE_FRACTIONAL":   "TIME_TRUNCATE_FRACTIONAL",
}

var sqlModeBitFlags = []struct {
	bit  uint64
	name string
}{
	{1 << 0, "REAL_AS_FLOAT"},
	{1 << 1, "PIPES_AS_CONCAT"},
	{1 << 2, "ANSI_QUOTES"},
	{1 << 3, "IGNORE_SPACE"},
	{1 << 4, "NOT_USED"},
	{1 << 5, "ONLY_FULL_GROUP_BY"},
	{1 << 6, "NO_UNSIGNED_SUBTRACTION"},
	{1 << 7, "NO_DIR_IN_CREATE"},
	{1 << 18, "ANSI"},
	{1 << 19, "NO_AUTO_VALUE_ON_ZERO"},
	{1 << 20, "NO_BACKSLASH_ESCAPES"},
	{1 << 21, "STRICT_TRANS_TABLES"},
	{1 << 22, "STRICT_ALL_TABLES"},
	{1 << 23, "NO_ZERO_IN_DATE"},
	{1 << 24, "NO_ZERO_DATE"},
	{1 << 25, "ALLOW_INVALID_DATES"},
	{1 << 26, "ERROR_FOR_DIVISION_BY_ZERO"},
	{1 << 27, "TRADITIONAL"},
	{1 << 29, "HIGH_NOT_PRECEDENCE"},
	{1 << 30, "NO_ENGINE_SUBSTITUTION"},
	{1 << 31, "PAD_CHAR_TO_FULL_LENGTH"},
	{1 << 32, "TIME_TRUNCATE_FRACTIONAL"},
}

const unsupportedMySQL80SQLModeBits uint64 = 0x1003FF00

func expandSQLModeForRead(mode string) string {
	trimmed := strings.TrimSpace(mode)
	if trimmed == "" {
		return ""
	}
	if strings.EqualFold(trimmed, "DEFAULT") {
		return strings.Join(compositeSQLModes["DEFAULT"], ",")
	}
	seen := make(map[string]bool)
	for _, part := range strings.Split(trimmed, ",") {
		tok := strings.ToUpper(strings.TrimSpace(part))
		if tok == "" {
			continue
		}
		if expanded, isComp := compositeSQLModes[tok]; isComp {
			for _, f := range expanded {
				seen[f] = true
			}
			continue
		}
		seen[tok] = true
	}
	var out []string
	for _, bf := range sqlModeBitFlags {
		if bf.name == "ANSI" || bf.name == "TRADITIONAL" {
			continue
		}
		if seen[bf.name] {
			out = append(out, bf.name)
		}
	}
	return strings.Join(out, ",")
}

func validateAndNormalizeSQLMode(val userVarValue) (string, error) {
	switch val.kind {
	case userVarNull:
		return "", errWrongValueForVar("sql_mode", "NULL")
	case userVarNumeric:
		s := strings.TrimSpace(val.str)
		if strings.ContainsAny(s, ".eE") {
			return "", errWrongTypeForVar("sql_mode")
		}
		if strings.HasPrefix(s, "-") {
			digits := strings.TrimLeft(s[1:], "0")
			if digits != "" {
				return "", errWrongValueForVar("sql_mode", s)
			}
			return "", nil
		}
		u, err := strconv.ParseUint(strings.TrimPrefix(s, "+"), 10, 64)
		if err != nil || u >= (uint64(1)<<33) {
			return "", errWrongValueForVar("sql_mode", s)
		}
		if bad := u & unsupportedMySQL80SQLModeBits; bad != 0 {
			return "", &Error{Code: 3899, SQLState: "HY000",
				Message: fmt.Sprintf("sql_mode=0x%x is not supported.", bad)}
		}
		var names []string
		for _, bf := range sqlModeBitFlags {
			if (u & bf.bit) != 0 {
				names = append(names, bf.name)
			}
		}
		return strings.Join(names, ","), nil
	case userVarString:
		if strings.TrimSpace(val.str) == "" {
			return "", nil
		}
		var out []string
		for _, rawPart := range strings.Split(val.str, ",") {
			tok := strings.ToUpper(strings.TrimSpace(rawPart))
			if tok == "" {
				return "", errWrongValueForVar("sql_mode", val.str)
			}
			canon, ok := validMySQL80SQLModes[tok]
			if !ok {
				return "", errWrongValueForVar("sql_mode", val.str)
			}
			out = append(out, canon)
		}
		return strings.Join(out, ","), nil
	default:
		return "", errWrongValueForVar("sql_mode", val.str)
	}
}

func isDefaultSetKeywordOrIdent(expr nodes.ExprNode) bool {
	switch e := expr.(type) {
	case *nodes.DefaultExpr:
		return true
	case *nodes.ColumnRef:
		return e.Schema == "" && e.Table == "" && strings.EqualFold(e.Column, "DEFAULT")
	default:
		return false
	}
}

func (c *Catalog) evalSetExpr(expr nodes.ExprNode) (userVarValue, bool) {
	switch e := expr.(type) {
	case *nodes.ParenExpr:
		return c.evalSetExpr(e.Expr)
	case *nodes.NullLit:
		return userVarValue{kind: userVarNull}, true
	case *nodes.StringLit:
		return userVarValue{kind: userVarString, str: e.Value}, true
	case *nodes.HexLit:
		raw, ok := decodeLiteralRawBytes(cleanHexDigits(e.Value), true)
		if !ok {
			return userVarValue{}, false
		}
		return userVarValue{kind: userVarString, str: string(raw)}, true
	case *nodes.IntLit:
		if s, ok := exactIntLiteral(e, ""); ok {
			return userVarValue{kind: userVarNumeric, str: s}, true
		}
		return userVarValue{kind: userVarNumeric, str: fmt.Sprintf("%d", e.Value)}, true
	case *nodes.FloatLit:
		return userVarValue{kind: userVarNumeric, str: e.Value}, true
	case *nodes.BoolLit:
		if e.Value {
			return userVarValue{kind: userVarNumeric, str: "1"}, true
		}
		return userVarValue{kind: userVarNumeric, str: "0"}, true
	case *nodes.BitLit:
		raw, ok := decodeLiteralRawBytes(cleanBitDigits(e.Value), false)
		if !ok {
			return userVarValue{}, false
		}
		return userVarValue{kind: userVarString, str: string(raw)}, true
	case *nodes.UnaryExpr:
		if e.Op == nodes.UnaryPlus || e.Op == nodes.UnaryMinus {
			sign := ""
			if e.Op == nodes.UnaryMinus {
				sign = "-"
			}
			switch op := e.Operand.(type) {
			case *nodes.IntLit:
				if s, ok := exactIntLiteral(op, sign); ok {
					return userVarValue{kind: userVarNumeric, str: s}, true
				}
				return userVarValue{kind: userVarNumeric, str: fmt.Sprintf("%s%d", sign, op.Value)}, true
			case *nodes.FloatLit:
				return userVarValue{kind: userVarNumeric, str: sign + op.Value}, true
			}
			return userVarValue{kind: userVarNumeric}, true
		}
		return userVarValue{}, false
	case *nodes.VariableRef:
		if !e.System {
			if v, exists := c.userVars[toLower(e.Name)]; exists {
				return v, true
			}
			return userVarValue{kind: userVarNull}, true
		}
		switch toLower(e.Name) {
		case "time_zone":
			if strings.EqualFold(e.Scope, "GLOBAL") {
				tz := c.globalTimeZone
				if tz == "" {
					tz = "SYSTEM"
				}
				return userVarValue{kind: userVarString, str: tz}, true
			}
			return userVarValue{kind: userVarString, str: c.session.TimeZone}, true
		case "system_time_zone":
			return userVarValue{kind: userVarString, str: c.SystemTimeZone()}, true
		case "sql_mode":
			if strings.EqualFold(e.Scope, "GLOBAL") {
				mode := c.globalSQLMode
				if mode == "" {
					mode = "DEFAULT"
				}
				return userVarValue{kind: userVarString, str: expandSQLModeForRead(mode)}, true
			}
			return userVarValue{kind: userVarString, str: expandSQLModeForRead(c.session.SQLMode)}, true
		case "foreign_key_checks":
			if c.session.ForeignKeyChecks {
				return userVarValue{kind: userVarNumeric, str: "1"}, true
			}
			return userVarValue{kind: userVarNumeric, str: "0"}, true
		default:
			return userVarValue{}, false
		}
	case *nodes.FuncCallExpr:
		switch toLower(e.Name) {
		case "concat":
			var b strings.Builder
			for _, arg := range e.Args {
				v, ok := c.evalSetExpr(arg)
				if !ok {
					return userVarValue{}, false
				}
				if v.kind == userVarNull {
					return userVarValue{kind: userVarNull}, true
				}
				b.WriteString(v.str)
			}
			return userVarValue{kind: userVarString, str: b.String()}, true
		case "coalesce", "ifnull":
			for _, arg := range e.Args {
				v, ok := c.evalSetExpr(arg)
				if !ok {
					return userVarValue{}, false
				}
				if v.kind != userVarNull {
					return v, true
				}
			}
			return userVarValue{kind: userVarNull}, true
		case "trim", "ltrim", "rtrim", "upper", "lower", "ucase", "lcase":
			if len(e.Args) != 1 {
				return userVarValue{}, false
			}
			v, ok := c.evalSetExpr(e.Args[0])
			if !ok || v.kind == userVarNull {
				return v, ok
			}
			s := v.str
			switch toLower(e.Name) {
			case "trim":
				s = strings.TrimSpace(s)
			case "ltrim":
				s = strings.TrimLeft(s, " ")
			case "rtrim":
				s = strings.TrimRight(s, " ")
			case "upper", "ucase":
				s = strings.ToUpper(s)
			case "lower", "lcase":
				s = strings.ToLower(s)
			}
			return userVarValue{kind: userVarString, str: s}, true
		}
	}
	return userVarValue{}, false
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
	if err := validateExpressionSemantics(stmt); err != nil {
		return err
	}
	switch s := stmt.(type) {
	case *nodes.CreateDatabaseStmt:
		return c.createDatabase(s)
	case *nodes.CreateTableStmt:
		return c.createTable(s, true, source)
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
		return c.createTrigger(s)
	case *nodes.DropTriggerStmt:
		return c.dropTrigger(s)
	case *nodes.CreateEventStmt:
		return c.createEvent(s)
	case *nodes.AlterEventStmt:
		return c.alterEvent(s)
	case *nodes.DropEventStmt:
		return c.dropEvent(s)
	case *nodes.CreateSpatialRefSysStmt:
		return validateCreateSpatialRefSys(s)
	case *nodes.CreateTablespaceStmt:
		return c.createTablespace(s)
	case *nodes.AlterTablespaceStmt:
		return c.alterTablespace(s)
	case *nodes.DropTablespaceStmt:
		return c.dropTablespace(s)
	case *nodes.SetStmt:
		return c.execSet(s)
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
	case *nodes.CreateSpatialRefSysStmt:
		return s.Loc.Start
	case *nodes.CreateTablespaceStmt:
		return s.Loc.Start
	case *nodes.AlterTablespaceStmt:
		return s.Loc.Start
	case *nodes.DropTablespaceStmt:
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
