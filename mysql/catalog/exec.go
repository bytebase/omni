package catalog

import (
	"fmt"
	"math"
	"math/big"
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
	pendingNamesCharset := ""
	activeScope := setScopeSession
	for _, asgn := range stmt.Assignments {
		if asgn.Column != nil && strings.HasPrefix(asgn.Column.Column, "@") && !strings.HasPrefix(asgn.Column.Column, "@@") {
			pendingNamesCharset = ""
			uName := toLower(strings.TrimPrefix(asgn.Column.Column, "@"))
			val, ok, err := c.evalSetExprChecked(asgn.Value)
			if err != nil {
				return err
			}
			if ok {
				updates = append(updates, func() {
					c.userVars[uName] = val
				})
			}
			continue
		}
		varName := normalizeSetVariableName(asgn.Column)
		scope := resolveSetAssignmentScope(asgn, &activeScope)
		if varName != "collate" {
			pendingNamesCharset = ""
		}
		switch varName {
		case "names":
			raw, charset, isNull, err := c.resolveSetCharsetValueChecked(asgn.Value)
			if err != nil {
				return err
			}
			if isNull {
				return errWrongValueForVar("character_set_client", "NULL")
			}
			coll, ok := defaultCollationForCharset[toLower(charset)]
			if !ok {
				return errUnknownCharacterSet(raw)
			}
			if !isPermissibleClientCharset(charset) {
				return errWrongValueForVar("character_set_client", charset)
			}
			pendingNamesCharset = charset
			updates = append(updates, func() {
				c.charsetClient = charset
				c.session.CharsetClient = charset
				c.collationConn = coll
				c.session.CollationConnection = coll
			})
		case "character set":
			raw, charset, isNull, err := c.resolveSetCharsetValueChecked(asgn.Value)
			if err != nil {
				return err
			}
			if isNull {
				return errWrongValueForVar("character_set_client", "NULL")
			}
			if _, ok := defaultCollationForCharset[toLower(charset)]; !ok {
				return errUnknownCharacterSet(raw)
			}
			if !isPermissibleClientCharset(charset) {
				return errWrongValueForVar("character_set_client", charset)
			}
			dbCS := c.defaultCharset
			if c.currentDB != "" {
				if db := c.GetDatabase(c.currentDB); db != nil && db.Charset != "" {
					dbCS = normalizeCharsetName(db.Charset)
				}
			}
			dbColl := defaultCollationForCharset[toLower(dbCS)]
			if dbColl == "" {
				dbColl = defaultCollationForCharset[toLower(c.defaultCharset)]
			}
			updates = append(updates, func() {
				c.charsetClient = charset
				c.session.CharsetClient = charset
				if dbColl != "" {
					c.collationConn = dbColl
					c.session.CollationConnection = dbColl
				}
			})
		case "collate":
			collation := nodeToSQLValue(asgn.Value)
			if collation != "" {
				if !isValidMySQLCollationName(collation) {
					return errUnknownCollation(collation)
				}
				coll := canonicalMySQLCollationName(collation)
				if pendingNamesCharset != "" && !charsetMatchesCollation(pendingNamesCharset, coll) {
					return &Error{Code: 1253, SQLState: "42000",
						Message: fmt.Sprintf("COLLATION '%s' is not valid for CHARACTER SET '%s'", coll, pendingNamesCharset)}
				}
				pendingNamesCharset = ""
				updates = append(updates, func() {
					c.collationConn = coll
					c.session.CollationConnection = coll
				})
			}
		default:
			if spec, ok := mysqlSysvarSpecs[varName]; ok {
				upd, err := spec.assign(c, varName, scope, asgn.Value)
				if err != nil {
					return err
				}
				if upd != nil {
					updates = append(updates, upd)
				}
			}
			// Silently accept all other SET variables.
		}
	}
	for _, fn := range updates {
		fn()
	}
	return nil
}

type setVarScope uint8

const (
	setScopeSession setVarScope = iota
	setScopeGlobal
	setScopePersist
	setScopePersistOnly
)

type sysvarScopeKind uint8

const (
	sysvarSessionAndGlobal sysvarScopeKind = iota
	sysvarGlobalOnlyStatic
	sysvarGlobalOnlyNonPersist
)

type sysvarSpec struct {
	kind       sysvarScopeKind
	assign     func(c *Catalog, varName string, scope setVarScope, expr nodes.ExprNode) (func(), error)
	readGlobal func(c *Catalog) userVarValue
	readSess   func(c *Catalog) userVarValue
}

var mysqlSysvarSpecs map[string]sysvarSpec

func init() {
	mysqlSysvarSpecs = map[string]sysvarSpec{
		"foreign_key_checks": makeBoolSysvarSpec(
			true,
			func(c *Catalog) bool { return c.globalForeignKeyChecks },
			func(c *Catalog, v bool) { c.globalForeignKeyChecks = v },
			func(c *Catalog) bool { return c.session.ForeignKeyChecks },
			func(c *Catalog, v bool) { c.foreignKeyChecks = v; c.session.ForeignKeyChecks = v },
		),
		"sql_generate_invisible_primary_key": makeBoolSysvarSpec(
			false,
			func(c *Catalog) bool { return c.globalGenerateGIPK },
			func(c *Catalog, v bool) { c.globalGenerateGIPK = v },
			func(c *Catalog) bool { return c.session.GenerateInvisiblePrimaryKey },
			func(c *Catalog, v bool) { c.generateGIPK = v; c.session.GenerateInvisiblePrimaryKey = v },
		),
		"show_gipk_in_create_table_and_information_schema": makeBoolSysvarSpec(
			false,
			func(c *Catalog) bool { return c.globalShowGIPK },
			func(c *Catalog, v bool) { c.globalShowGIPK = v },
			func(c *Catalog) bool { return c.session.ShowGIPK },
			func(c *Catalog, v bool) { c.showGIPK = v; c.session.ShowGIPK = v },
		),
		"explicit_defaults_for_timestamp": makeBoolSysvarSpec(
			true,
			func(c *Catalog) bool { return c.globalExplicitDefaultsTS },
			func(c *Catalog, v bool) { c.globalExplicitDefaultsTS = v },
			func(c *Catalog) bool { return c.session.ExplicitDefaultsForTimestamp },
			func(c *Catalog, v bool) { c.session.ExplicitDefaultsForTimestamp = v },
		),
		"innodb_directories": {
			kind: sysvarGlobalOnlyStatic,
			assign: func(_ *Catalog, varName string, scope setVarScope, _ nodes.ExprNode) (func(), error) {
				return checkReadOnlySysvarWrite(varName, scope, true)
			},
		},
		"innodb_page_size": {
			kind: sysvarGlobalOnlyStatic,
			assign: func(_ *Catalog, varName string, scope setVarScope, _ nodes.ExprNode) (func(), error) {
				return checkReadOnlySysvarWrite(varName, scope, true)
			},
			readGlobal: func(c *Catalog) userVarValue {
				return userVarValue{kind: userVarNumeric, str: strconv.FormatUint(c.InnoDBPageSize(), 10)}
			},
		},
		"system_time_zone": {
			kind: sysvarGlobalOnlyNonPersist,
			assign: func(_ *Catalog, varName string, scope setVarScope, _ nodes.ExprNode) (func(), error) {
				return checkReadOnlySysvarWrite(varName, scope, false)
			},
			readGlobal: func(c *Catalog) userVarValue {
				return userVarValue{kind: userVarString, str: c.SystemTimeZone()}
			},
		},
		"sql_mode": {
			kind:   sysvarSessionAndGlobal,
			assign: assignSQLModeSysvar,
			readGlobal: func(c *Catalog) userVarValue {
				return userVarValue{kind: userVarString, str: expandSQLModeForRead(effectiveNonEmpty(c.globalSQLMode, "DEFAULT"))}
			},
			readSess: func(c *Catalog) userVarValue {
				return userVarValue{kind: userVarString, str: expandSQLModeForRead(c.session.SQLMode)}
			},
		},
		"time_zone": {
			kind:   sysvarSessionAndGlobal,
			assign: assignTimeZoneSysvar,
			readGlobal: func(c *Catalog) userVarValue {
				return userVarValue{kind: userVarString, str: effectiveNonEmpty(c.globalTimeZone, "SYSTEM")}
			},
			readSess: func(c *Catalog) userVarValue { return userVarValue{kind: userVarString, str: c.session.TimeZone} },
		},
		"character_set_client": {
			kind:   sysvarSessionAndGlobal,
			assign: makeCharsetSysvarAssign(true),
			readGlobal: func(c *Catalog) userVarValue {
				return userVarValue{kind: userVarString, str: c.effectiveGlobalClientCharset()}
			},
			readSess: func(c *Catalog) userVarValue {
				return userVarValue{kind: userVarString, str: c.effectiveClientCharset()}
			},
		},
		"character_set_connection": {
			kind:   sysvarSessionAndGlobal,
			assign: makeCharsetSysvarAssign(false),
			readGlobal: func(c *Catalog) userVarValue {
				return userVarValue{kind: userVarString, str: c.effectiveGlobalConnCharset()}
			},
			readSess: func(c *Catalog) userVarValue { return userVarValue{kind: userVarString, str: c.effectiveConnCharset()} },
		},
		"collation_connection": {
			kind:   sysvarSessionAndGlobal,
			assign: assignCollationConnSysvar,
			readGlobal: func(c *Catalog) userVarValue {
				return userVarValue{kind: userVarString, str: c.effectiveGlobalCollationConn()}
			},
			readSess: func(c *Catalog) userVarValue {
				return userVarValue{kind: userVarString, str: c.effectiveSessionCollationConn()}
			},
		},
	}
}

func effectiveNonEmpty(s, fallback string) string {
	if s == "" {
		return fallback
	}
	return s
}

func isPermissibleClientCharset(cs string) bool {
	switch toLower(normalizeCharsetName(cs)) {
	case "ucs2", "utf16", "utf16le", "utf32":
		return false
	default:
		return true
	}
}

func canonicalMySQLCollationName(raw string) string {
	coll := toLower(strings.TrimSpace(raw))
	if strings.HasPrefix(coll, "utf8_") {
		return "utf8mb3_" + strings.TrimPrefix(coll, "utf8_")
	}
	return coll
}

func checkReadOnlySysvarWrite(varName string, scope setVarScope, allowPersistOnly bool) (func(), error) {
	if allowPersistOnly {
		if scope != setScopePersistOnly {
			return nil, errIncorrectGlobalLocalVar(varName, "read only variable")
		}
		return nil, nil
	}
	return nil, errIncorrectGlobalLocalVar(varName, "non persistent read only variable")
}

func applyScopedUpdate(scope setVarScope, setGlobal, setSession func()) func() {
	switch scope {
	case setScopePersistOnly:
		return nil
	case setScopeGlobal, setScopePersist:
		return setGlobal
	default:
		return setSession
	}
}

func makeBoolSysvarSpec(factoryDefault bool, getG func(*Catalog) bool, setG func(*Catalog, bool), getS func(*Catalog) bool, setS func(*Catalog, bool)) sysvarSpec {
	return sysvarSpec{
		kind: sysvarSessionAndGlobal,
		assign: func(c *Catalog, varName string, scope setVarScope, expr nodes.ExprNode) (func(), error) {
			if isDefaultSetKeywordOrIdent(expr) {
				if scope == setScopeGlobal || scope == setScopePersist {
					return func() { setG(c, factoryDefault) }, nil
				}
				return applyScopedUpdate(scope, nil, func() { setS(c, getG(c)) }), nil
			}
			v, ok, err := c.evalSetBoolChecked(varName, expr)
			if err != nil {
				return nil, err
			}
			if !ok {
				return nil, errWrongValueForVar(varName, strings.TrimSpace(nodeToSQL(expr)))
			}
			return applyScopedUpdate(scope, func() { setG(c, v) }, func() { setS(c, v) }), nil
		},
		readGlobal: func(c *Catalog) userVarValue {
			if getG(c) {
				return userVarValue{kind: userVarNumeric, str: "1"}
			}
			return userVarValue{kind: userVarNumeric, str: "0"}
		},
		readSess: func(c *Catalog) userVarValue {
			if getS(c) {
				return userVarValue{kind: userVarNumeric, str: "1"}
			}
			return userVarValue{kind: userVarNumeric, str: "0"}
		},
	}
}

func assignSQLModeSysvar(c *Catalog, varName string, scope setVarScope, expr nodes.ExprNode) (func(), error) {
	if isDefaultSetKeywordOrIdent(expr) {
		if scope == setScopeGlobal || scope == setScopePersist {
			return func() { c.globalSQLMode = "DEFAULT" }, nil
		}
		defMode := effectiveNonEmpty(c.globalSQLMode, "DEFAULT")
		return applyScopedUpdate(scope, nil, func() { c.session.SQLMode = defMode }), nil
	}
	var val userVarValue
	if cr, isCol := expr.(*nodes.ColumnRef); isCol && cr.Schema == "" && cr.Table == "" {
		val = userVarValue{kind: userVarString, str: cr.Column}
	} else {
		var ok bool
		var evalErr error
		val, ok, evalErr = c.evalSetExprChecked(expr)
		if evalErr != nil {
			return nil, evalErr
		}
		if !ok {
			return nil, errWrongValueForVar(varName, strings.TrimSpace(nodeToSQL(expr)))
		}
	}
	mode, err := validateAndNormalizeSQLMode(val)
	if err != nil {
		return nil, err
	}
	return applyScopedUpdate(scope, func() { c.globalSQLMode = mode }, func() { c.session.SQLMode = mode }), nil
}

func assignTimeZoneSysvar(c *Catalog, varName string, scope setVarScope, expr nodes.ExprNode) (func(), error) {
	if isDefaultSetKeywordOrIdent(expr) {
		if scope == setScopeGlobal || scope == setScopePersist {
			return func() { c.globalTimeZone = "SYSTEM" }, nil
		}
		defTZ := effectiveNonEmpty(c.globalTimeZone, "SYSTEM")
		return applyScopedUpdate(scope, nil, func() { c.session.TimeZone = defTZ }), nil
	}
	var rawTZ string
	if cr, isCol := expr.(*nodes.ColumnRef); isCol && cr.Schema == "" && cr.Table == "" {
		rawTZ = cr.Column
	} else {
		val, ok, evalErr := c.evalSetExprChecked(expr)
		if evalErr != nil {
			return nil, evalErr
		}
		if !ok {
			return nil, errUnknownTimeZone(strings.TrimSpace(nodeToSQL(expr)))
		}
		switch val.kind {
		case userVarNull:
			return nil, errWrongValueForVar(varName, "NULL")
		case userVarNumeric:
			return nil, errWrongTypeForVar(varName)
		case userVarString:
			rawTZ = val.str
		}
	}
	tz, err := c.validateAndNormalizeSessionTimeZone(rawTZ)
	if err != nil {
		return nil, err
	}
	return applyScopedUpdate(scope, func() { c.globalTimeZone = tz }, func() { c.session.TimeZone = tz }), nil
}

func makeCharsetSysvarAssign(isClient bool) func(c *Catalog, varName string, scope setVarScope, expr nodes.ExprNode) (func(), error) {
	return func(c *Catalog, varName string, scope setVarScope, expr nodes.ExprNode) (func(), error) {
		if isDefaultSetKeywordOrIdent(expr) {
			if isClient {
				if scope == setScopeGlobal || scope == setScopePersist {
					return func() { c.globalCharsetClient = c.defaultCharset }, nil
				}
				target := c.effectiveGlobalClientCharset()
				return applyScopedUpdate(scope, nil, func() {
					c.charsetClient = target
					c.session.CharsetClient = target
				}), nil
			}
			if scope == setScopeGlobal || scope == setScopePersist {
				defColl := defaultCollationForCharset[toLower(c.defaultCharset)]
				return func() { c.globalCollationConn = defColl }, nil
			}
			targetColl := c.effectiveGlobalCollationConn()
			return applyScopedUpdate(scope, nil, func() {
				c.collationConn = targetColl
				c.session.CollationConnection = targetColl
			}), nil
		}
		raw, charset, isNull, err := c.resolveSetCharsetValueChecked(expr)
		if err != nil {
			return nil, err
		}
		if isNull {
			return nil, errWrongValueForVar(varName, "NULL")
		}
		coll, ok := defaultCollationForCharset[toLower(charset)]
		if !ok {
			return nil, errUnknownCharacterSet(raw)
		}
		if isClient && !isPermissibleClientCharset(charset) {
			return nil, errWrongValueForVar(varName, charset)
		}
		if isClient {
			return applyScopedUpdate(
				scope,
				func() { c.globalCharsetClient = charset },
				func() {
					c.charsetClient = charset
					c.session.CharsetClient = charset
				},
			), nil
		}
		return applyScopedUpdate(
			scope,
			func() { c.globalCollationConn = coll },
			func() {
				c.collationConn = coll
				c.session.CollationConnection = coll
			},
		), nil
	}
}

func assignCollationConnSysvar(c *Catalog, varName string, scope setVarScope, expr nodes.ExprNode) (func(), error) {
	if isDefaultSetKeywordOrIdent(expr) {
		if scope == setScopeGlobal || scope == setScopePersist {
			defColl := effectiveNonEmpty(defaultCollationForCharset[toLower(c.defaultCharset)], "utf8mb4_0900_ai_ci")
			return func() { c.globalCollationConn = defColl }, nil
		}
		targetColl := c.effectiveGlobalCollationConn()
		return applyScopedUpdate(scope, nil, func() {
			c.collationConn = targetColl
			c.session.CollationConnection = targetColl
		}), nil
	}
	var rawColl string
	if cr, isCol := expr.(*nodes.ColumnRef); isCol && cr.Schema == "" && cr.Table == "" {
		rawColl = cr.Column
	} else {
		val, ok, evalErr := c.evalSetExprChecked(expr)
		if evalErr != nil {
			return nil, evalErr
		}
		if ok {
			if val.kind == userVarNull {
				return nil, errWrongValueForVar(varName, "NULL")
			}
			rawColl = val.str
		} else {
			rawColl = nodeToSQLValue(expr)
		}
	}
	if !isValidMySQLCollationName(rawColl) {
		return nil, errUnknownCollation(rawColl)
	}
	coll := canonicalMySQLCollationName(rawColl)
	return applyScopedUpdate(
		scope,
		func() { c.globalCollationConn = coll },
		func() {
			c.collationConn = coll
			c.session.CollationConnection = coll
		},
	), nil
}

func resolveSetAssignmentScope(asgn *nodes.Assignment, activeScope *setVarScope) setVarScope {
	if asgn == nil {
		return *activeScope
	}
	if asgn.Scope != "" {
		switch strings.ToUpper(asgn.Scope) {
		case "PERSIST_ONLY":
			*activeScope = setScopePersistOnly
		case "PERSIST":
			*activeScope = setScopePersist
		case "GLOBAL":
			*activeScope = setScopeGlobal
		case "SESSION", "LOCAL":
			*activeScope = setScopeSession
		}
		return *activeScope
	}
	if asgn.Column != nil {
		col := toLower(strings.TrimPrefix(asgn.Column.Column, "@@"))
		if asgn.Column.Table != "" {
			col = toLower(strings.TrimPrefix(asgn.Column.Table, "@@")) + "." + col
		}
		switch {
		case strings.HasPrefix(col, "persist_only."):
			return setScopePersistOnly
		case strings.HasPrefix(col, "persist."):
			return setScopePersist
		case strings.HasPrefix(col, "global."):
			return setScopeGlobal
		case strings.HasPrefix(col, "session."), strings.HasPrefix(col, "local."):
			return setScopeSession
		}
	}
	return *activeScope
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

var validMySQL80Collations = map[string]bool{
	"armscii8_bin": true, "armscii8_general_ci": true, "ascii_bin": true, "ascii_general_ci": true,
	"big5_bin": true, "big5_chinese_ci": true, "binary": true, "cp1250_bin": true,
	"cp1250_croatian_ci": true, "cp1250_czech_cs": true, "cp1250_general_ci": true, "cp1250_polish_ci": true,
	"cp1251_bin": true, "cp1251_bulgarian_ci": true, "cp1251_general_ci": true, "cp1251_general_cs": true,
	"cp1251_ukrainian_ci": true, "cp1256_bin": true, "cp1256_general_ci": true, "cp1257_bin": true,
	"cp1257_general_ci": true, "cp1257_lithuanian_ci": true, "cp850_bin": true, "cp850_general_ci": true,
	"cp852_bin": true, "cp852_general_ci": true, "cp866_bin": true, "cp866_general_ci": true,
	"cp932_bin": true, "cp932_japanese_ci": true, "dec8_bin": true, "dec8_swedish_ci": true,
	"eucjpms_bin": true, "eucjpms_japanese_ci": true, "euckr_bin": true, "euckr_korean_ci": true,
	"gb18030_bin": true, "gb18030_chinese_ci": true, "gb18030_unicode_520_ci": true, "gb2312_bin": true,
	"gb2312_chinese_ci": true, "gbk_bin": true, "gbk_chinese_ci": true, "geostd8_bin": true,
	"geostd8_general_ci": true, "greek_bin": true, "greek_general_ci": true, "hebrew_bin": true,
	"hebrew_general_ci": true, "hp8_bin": true, "hp8_english_ci": true, "keybcs2_bin": true,
	"keybcs2_general_ci": true, "koi8r_bin": true, "koi8r_general_ci": true, "koi8u_bin": true,
	"koi8u_general_ci": true, "latin1_bin": true, "latin1_danish_ci": true, "latin1_general_ci": true,
	"latin1_general_cs": true, "latin1_german1_ci": true, "latin1_german2_ci": true, "latin1_spanish_ci": true,
	"latin1_swedish_ci": true, "latin2_bin": true, "latin2_croatian_ci": true, "latin2_czech_cs": true,
	"latin2_general_ci": true, "latin2_hungarian_ci": true, "latin5_bin": true, "latin5_turkish_ci": true,
	"latin7_bin": true, "latin7_estonian_cs": true, "latin7_general_ci": true, "latin7_general_cs": true,
	"macce_bin": true, "macce_general_ci": true, "macroman_bin": true, "macroman_general_ci": true,
	"sjis_bin": true, "sjis_japanese_ci": true, "swe7_bin": true, "swe7_swedish_ci": true,
	"tis620_bin": true, "tis620_thai_ci": true, "ucs2_bin": true, "ucs2_croatian_ci": true,
	"ucs2_czech_ci": true, "ucs2_danish_ci": true, "ucs2_esperanto_ci": true, "ucs2_estonian_ci": true,
	"ucs2_general_ci": true, "ucs2_general_mysql500_ci": true, "ucs2_german2_ci": true, "ucs2_hungarian_ci": true,
	"ucs2_icelandic_ci": true, "ucs2_latvian_ci": true, "ucs2_lithuanian_ci": true, "ucs2_persian_ci": true,
	"ucs2_polish_ci": true, "ucs2_romanian_ci": true, "ucs2_roman_ci": true, "ucs2_sinhala_ci": true,
	"ucs2_slovak_ci": true, "ucs2_slovenian_ci": true, "ucs2_spanish2_ci": true, "ucs2_spanish_ci": true,
	"ucs2_swedish_ci": true, "ucs2_turkish_ci": true, "ucs2_unicode_520_ci": true, "ucs2_unicode_ci": true,
	"ucs2_vietnamese_ci": true, "ujis_bin": true, "ujis_japanese_ci": true, "utf16le_bin": true,
	"utf16le_general_ci": true, "utf16_bin": true, "utf16_croatian_ci": true, "utf16_czech_ci": true,
	"utf16_danish_ci": true, "utf16_esperanto_ci": true, "utf16_estonian_ci": true, "utf16_general_ci": true,
	"utf16_german2_ci": true, "utf16_hungarian_ci": true, "utf16_icelandic_ci": true, "utf16_latvian_ci": true,
	"utf16_lithuanian_ci": true, "utf16_persian_ci": true, "utf16_polish_ci": true, "utf16_romanian_ci": true,
	"utf16_roman_ci": true, "utf16_sinhala_ci": true, "utf16_slovak_ci": true, "utf16_slovenian_ci": true,
	"utf16_spanish2_ci": true, "utf16_spanish_ci": true, "utf16_swedish_ci": true, "utf16_turkish_ci": true,
	"utf16_unicode_520_ci": true, "utf16_unicode_ci": true, "utf16_vietnamese_ci": true, "utf32_bin": true,
	"utf32_croatian_ci": true, "utf32_czech_ci": true, "utf32_danish_ci": true, "utf32_esperanto_ci": true,
	"utf32_estonian_ci": true, "utf32_general_ci": true, "utf32_german2_ci": true, "utf32_hungarian_ci": true,
	"utf32_icelandic_ci": true, "utf32_latvian_ci": true, "utf32_lithuanian_ci": true, "utf32_persian_ci": true,
	"utf32_polish_ci": true, "utf32_romanian_ci": true, "utf32_roman_ci": true, "utf32_sinhala_ci": true,
	"utf32_slovak_ci": true, "utf32_slovenian_ci": true, "utf32_spanish2_ci": true, "utf32_spanish_ci": true,
	"utf32_swedish_ci": true, "utf32_turkish_ci": true, "utf32_unicode_520_ci": true, "utf32_unicode_ci": true,
	"utf32_vietnamese_ci": true, "utf8mb3_bin": true, "utf8mb3_croatian_ci": true, "utf8mb3_czech_ci": true,
	"utf8mb3_danish_ci": true, "utf8mb3_esperanto_ci": true, "utf8mb3_estonian_ci": true, "utf8mb3_general_ci": true,
	"utf8mb3_general_mysql500_ci": true, "utf8mb3_german2_ci": true, "utf8mb3_hungarian_ci": true, "utf8mb3_icelandic_ci": true,
	"utf8mb3_latvian_ci": true, "utf8mb3_lithuanian_ci": true, "utf8mb3_persian_ci": true, "utf8mb3_polish_ci": true,
	"utf8mb3_romanian_ci": true, "utf8mb3_roman_ci": true, "utf8mb3_sinhala_ci": true, "utf8mb3_slovak_ci": true,
	"utf8mb3_slovenian_ci": true, "utf8mb3_spanish2_ci": true, "utf8mb3_spanish_ci": true, "utf8mb3_swedish_ci": true,
	"utf8mb3_tolower_ci": true, "utf8mb3_turkish_ci": true, "utf8mb3_unicode_520_ci": true, "utf8mb3_unicode_ci": true,
	"utf8mb3_vietnamese_ci": true, "utf8mb4_0900_ai_ci": true, "utf8mb4_0900_as_ci": true, "utf8mb4_0900_as_cs": true,
	"utf8mb4_0900_bin": true, "utf8mb4_bg_0900_ai_ci": true, "utf8mb4_bg_0900_as_cs": true, "utf8mb4_bin": true,
	"utf8mb4_bs_0900_ai_ci": true, "utf8mb4_bs_0900_as_cs": true, "utf8mb4_croatian_ci": true, "utf8mb4_cs_0900_ai_ci": true,
	"utf8mb4_cs_0900_as_cs": true, "utf8mb4_czech_ci": true, "utf8mb4_danish_ci": true, "utf8mb4_da_0900_ai_ci": true,
	"utf8mb4_da_0900_as_cs": true, "utf8mb4_de_pb_0900_ai_ci": true, "utf8mb4_de_pb_0900_as_cs": true, "utf8mb4_eo_0900_ai_ci": true,
	"utf8mb4_eo_0900_as_cs": true, "utf8mb4_esperanto_ci": true, "utf8mb4_estonian_ci": true, "utf8mb4_es_0900_ai_ci": true,
	"utf8mb4_es_0900_as_cs": true, "utf8mb4_es_trad_0900_ai_ci": true, "utf8mb4_es_trad_0900_as_cs": true, "utf8mb4_et_0900_ai_ci": true,
	"utf8mb4_et_0900_as_cs": true, "utf8mb4_general_ci": true, "utf8mb4_german2_ci": true, "utf8mb4_gl_0900_ai_ci": true,
	"utf8mb4_gl_0900_as_cs": true, "utf8mb4_hr_0900_ai_ci": true, "utf8mb4_hr_0900_as_cs": true, "utf8mb4_hungarian_ci": true,
	"utf8mb4_hu_0900_ai_ci": true, "utf8mb4_hu_0900_as_cs": true, "utf8mb4_icelandic_ci": true, "utf8mb4_is_0900_ai_ci": true,
	"utf8mb4_is_0900_as_cs": true, "utf8mb4_ja_0900_as_cs": true, "utf8mb4_ja_0900_as_cs_ks": true, "utf8mb4_latvian_ci": true,
	"utf8mb4_la_0900_ai_ci": true, "utf8mb4_la_0900_as_cs": true, "utf8mb4_lithuanian_ci": true, "utf8mb4_lt_0900_ai_ci": true,
	"utf8mb4_lt_0900_as_cs": true, "utf8mb4_lv_0900_ai_ci": true, "utf8mb4_lv_0900_as_cs": true, "utf8mb4_mn_cyrl_0900_ai_ci": true,
	"utf8mb4_mn_cyrl_0900_as_cs": true, "utf8mb4_nb_0900_ai_ci": true, "utf8mb4_nb_0900_as_cs": true, "utf8mb4_nn_0900_ai_ci": true,
	"utf8mb4_nn_0900_as_cs": true, "utf8mb4_persian_ci": true, "utf8mb4_pl_0900_ai_ci": true, "utf8mb4_pl_0900_as_cs": true,
	"utf8mb4_polish_ci": true, "utf8mb4_romanian_ci": true, "utf8mb4_roman_ci": true, "utf8mb4_ro_0900_ai_ci": true,
	"utf8mb4_ro_0900_as_cs": true, "utf8mb4_ru_0900_ai_ci": true, "utf8mb4_ru_0900_as_cs": true, "utf8mb4_sinhala_ci": true,
	"utf8mb4_sk_0900_ai_ci": true, "utf8mb4_sk_0900_as_cs": true, "utf8mb4_slovak_ci": true, "utf8mb4_slovenian_ci": true,
	"utf8mb4_sl_0900_ai_ci": true, "utf8mb4_sl_0900_as_cs": true, "utf8mb4_spanish2_ci": true, "utf8mb4_spanish_ci": true,
	"utf8mb4_sr_latn_0900_ai_ci": true, "utf8mb4_sr_latn_0900_as_cs": true, "utf8mb4_sv_0900_ai_ci": true, "utf8mb4_sv_0900_as_cs": true,
	"utf8mb4_swedish_ci": true, "utf8mb4_tr_0900_ai_ci": true, "utf8mb4_tr_0900_as_cs": true, "utf8mb4_turkish_ci": true,
	"utf8mb4_unicode_520_ci": true, "utf8mb4_unicode_ci": true, "utf8mb4_vietnamese_ci": true, "utf8mb4_vi_0900_ai_ci": true,
	"utf8mb4_vi_0900_as_cs": true, "utf8mb4_zh_0900_as_cs": true,
}

func isValidMySQLCollationName(coll string) bool {
	c := toLower(strings.TrimSpace(coll))
	if c == "" {
		return false
	}
	if strings.HasPrefix(c, "utf8_") {
		c = "utf8mb3_" + strings.TrimPrefix(c, "utf8_")
	}
	return validMySQL80Collations[c]
}

func (c *Catalog) evalSetExpr(expr nodes.ExprNode) (userVarValue, bool) {
	v, ok, err := c.evalSetExprChecked(expr)
	if err != nil {
		return userVarValue{}, false
	}
	return v, ok
}

func (c *Catalog) evalSetExprChecked(expr nodes.ExprNode) (userVarValue, bool, error) {
	switch e := expr.(type) {
	case *nodes.ParenExpr:
		return c.evalSetExprChecked(e.Expr)
	case *nodes.NullLit:
		return userVarValue{kind: userVarNull}, true, nil
	case *nodes.StringLit:
		return userVarValue{kind: userVarString, str: e.Value}, true, nil
	case *nodes.HexLit:
		raw, ok := decodeLiteralRawBytes(cleanHexDigits(e.Value), true)
		if !ok {
			return userVarValue{}, false, nil
		}
		return userVarValue{kind: userVarString, str: string(raw)}, true, nil
	case *nodes.IntLit:
		if s, ok := exactIntLiteral(e, ""); ok {
			return userVarValue{kind: userVarNumeric, str: s}, true, nil
		}
		return userVarValue{kind: userVarNumeric, str: fmt.Sprintf("%d", e.Value)}, true, nil
	case *nodes.FloatLit:
		return userVarValue{kind: userVarNumeric, str: e.Value}, true, nil
	case *nodes.BoolLit:
		if e.Value {
			return userVarValue{kind: userVarNumeric, str: "1"}, true, nil
		}
		return userVarValue{kind: userVarNumeric, str: "0"}, true, nil
	case *nodes.BitLit:
		raw, ok := decodeLiteralRawBytes(cleanBitDigits(e.Value), false)
		if !ok {
			return userVarValue{}, false, nil
		}
		return userVarValue{kind: userVarString, str: string(raw)}, true, nil
	case *nodes.UnaryExpr:
		if e.Op == nodes.UnaryPlus || e.Op == nodes.UnaryMinus {
			sign := ""
			if e.Op == nodes.UnaryMinus {
				sign = "-"
			}
			switch op := e.Operand.(type) {
			case *nodes.IntLit:
				if s, ok := exactIntLiteral(op, sign); ok {
					return userVarValue{kind: userVarNumeric, str: s}, true, nil
				}
				return userVarValue{kind: userVarNumeric, str: fmt.Sprintf("%s%d", sign, op.Value)}, true, nil
			case *nodes.FloatLit:
				return userVarValue{kind: userVarNumeric, str: sign + op.Value}, true, nil
			}
			v, ok, err := c.evalSetExprChecked(e.Operand)
			if err != nil || !ok {
				return userVarValue{}, false, err
			}
			if v.kind == userVarNull {
				return v, true, nil
			}
			n, ok := parseSetNumericInt64(v)
			if !ok {
				return userVarValue{}, false, nil
			}
			if e.Op == nodes.UnaryMinus {
				if n == math.MinInt64 {
					return userVarValue{}, false, errBigintOutOfRange()
				}
				n = -n
			}
			return userVarValue{kind: userVarNumeric, str: strconv.FormatInt(n, 10)}, true, nil
		}
		return userVarValue{}, false, nil
	case *nodes.BinaryExpr:
		return c.evalSetBinaryExprChecked(e)
	case *nodes.VariableRef:
		if !e.System {
			if v, exists := c.userVars[toLower(e.Name)]; exists {
				return v, true, nil
			}
			return userVarValue{kind: userVarNull}, true, nil
		}
		if strings.EqualFold(e.Scope, "PERSIST") || strings.EqualFold(e.Scope, "PERSIST_ONLY") {
			return userVarValue{}, false, &Error{
				Code:     1193,
				SQLState: "HY000",
				Message:  fmt.Sprintf("Unknown system variable '%s.%s'", strings.ToUpper(e.Scope), e.Name),
			}
		}
		spec, ok := mysqlSysvarSpecs[toLower(e.Name)]
		if !ok {
			return userVarValue{}, false, nil
		}
		if spec.kind != sysvarSessionAndGlobal {
			if strings.EqualFold(e.Scope, "SESSION") || strings.EqualFold(e.Scope, "LOCAL") {
				return userVarValue{}, false, errIncorrectGlobalLocalVar(toLower(e.Name), "GLOBAL variable")
			}
			if spec.readGlobal != nil {
				return spec.readGlobal(c), true, nil
			}
			return userVarValue{}, false, nil
		}
		if strings.EqualFold(e.Scope, "GLOBAL") {
			return spec.readGlobal(c), true, nil
		}
		return spec.readSess(c), true, nil
	case *nodes.FuncCallExpr:
		switch toLower(e.Name) {
		case "concat":
			var b strings.Builder
			for _, arg := range e.Args {
				v, ok, err := c.evalSetExprChecked(arg)
				if err != nil || !ok {
					return userVarValue{}, false, err
				}
				if v.kind == userVarNull {
					return userVarValue{kind: userVarNull}, true, nil
				}
				b.WriteString(v.str)
			}
			return userVarValue{kind: userVarString, str: b.String()}, true, nil
		case "coalesce", "ifnull":
			for _, arg := range e.Args {
				v, ok, err := c.evalSetExprChecked(arg)
				if err != nil || !ok {
					return userVarValue{}, false, err
				}
				if v.kind != userVarNull {
					return v, true, nil
				}
			}
			return userVarValue{kind: userVarNull}, true, nil
		case "trim", "ltrim", "rtrim", "upper", "lower", "ucase", "lcase":
			if len(e.Args) != 1 {
				return userVarValue{}, false, nil
			}
			v, ok, err := c.evalSetExprChecked(e.Args[0])
			if err != nil || !ok || v.kind == userVarNull {
				return v, ok, err
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
			return userVarValue{kind: userVarString, str: s}, true, nil
		}
	}
	return userVarValue{}, false, nil
}

func (c *Catalog) evalSetBinaryExprChecked(e *nodes.BinaryExpr) (userVarValue, bool, error) {
	if e == nil {
		return userVarValue{}, false, nil
	}
	left, ok, err := c.evalSetExprChecked(e.Left)
	if err != nil || !ok {
		return userVarValue{}, false, err
	}
	right, ok, err := c.evalSetExprChecked(e.Right)
	if err != nil || !ok {
		return userVarValue{}, false, err
	}
	if left.kind == userVarNull || right.kind == userVarNull {
		return userVarValue{kind: userVarNull}, true, nil
	}
	ln, lok := parseSetNumericInt64(left)
	rn, rok := parseSetNumericInt64(right)
	if !lok || !rok {
		return userVarValue{}, false, nil
	}
	var out int64
	switch e.Op {
	case nodes.BinOpAdd:
		if (rn > 0 && ln > math.MaxInt64-rn) || (rn < 0 && ln < math.MinInt64-rn) {
			return userVarValue{}, false, errBigintOutOfRange()
		}
		out = ln + rn
	case nodes.BinOpSub:
		if (rn < 0 && ln > math.MaxInt64+rn) || (rn > 0 && ln < math.MinInt64+rn) {
			return userVarValue{}, false, errBigintOutOfRange()
		}
		out = ln - rn
	case nodes.BinOpMul:
		if ln != 0 && rn != 0 {
			if (ln == -1 && rn == math.MinInt64) || (rn == -1 && ln == math.MinInt64) {
				return userVarValue{}, false, errBigintOutOfRange()
			}
			prod := ln * rn
			if prod/ln != rn {
				return userVarValue{}, false, errBigintOutOfRange()
			}
			out = prod
		}
	case nodes.BinOpDiv:
		if rn == 0 {
			return userVarValue{kind: userVarNull}, true, nil
		}
		return userVarValue{kind: userVarNumeric, str: big.NewRat(ln, rn).FloatString(4)}, true, nil
	case nodes.BinOpDivInt:
		if rn == 0 {
			return userVarValue{kind: userVarNull}, true, nil
		}
		if ln == math.MinInt64 && rn == -1 {
			return userVarValue{}, false, errBigintOutOfRange()
		}
		out = ln / rn
	case nodes.BinOpMod:
		if rn == 0 {
			return userVarValue{kind: userVarNull}, true, nil
		}
		out = ln % rn
	case nodes.BinOpBitOr:
		out = ln | rn
	case nodes.BinOpBitAnd:
		out = ln & rn
	case nodes.BinOpBitXor:
		out = ln ^ rn
	case nodes.BinOpShiftLeft:
		if rn < 0 || rn >= 64 {
			out = 0
		} else {
			out = int64(uint64(ln) << uint(rn))
		}
	case nodes.BinOpShiftRight:
		if rn < 0 || rn >= 64 {
			out = 0
		} else {
			out = int64(uint64(ln) >> uint(rn))
		}
	default:
		return userVarValue{}, false, nil
	}
	return userVarValue{kind: userVarNumeric, str: strconv.FormatInt(out, 10)}, true, nil
}

func parseSetNumericInt64(v userVarValue) (int64, bool) {
	if v.kind != userVarNumeric && v.kind != userVarString {
		return 0, false
	}
	s := strings.TrimSpace(v.str)
	if s == "" {
		return 0, false
	}
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, false
	}
	return n, true
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

func (c *Catalog) evalSetBoolChecked(varName string, expr nodes.ExprNode) (bool, bool, error) {
	if cr, isCol := expr.(*nodes.ColumnRef); isCol && cr.Schema == "" && cr.Table == "" {
		switch toLower(cr.Column) {
		case "on", "true":
			return true, true, nil
		case "off", "false":
			return false, true, nil
		default:
			return false, false, nil
		}
	}
	v, ok, err := c.evalSetExprChecked(expr)
	if err != nil {
		return false, false, err
	}
	if !ok {
		return false, false, nil
	}
	if v.kind == userVarNull {
		return false, false, errWrongValueForVar(varName, "NULL")
	}
	switch toLower(strings.TrimSpace(v.str)) {
	case "0", "off", "false":
		return false, true, nil
	case "1", "on", "true":
		return true, true, nil
	default:
		return false, false, errWrongValueForVar(varName, v.str)
	}
}

func (c *Catalog) resolveSetCharsetValueChecked(expr nodes.ExprNode) (raw, norm string, isNull bool, err error) {
	if isDefaultSetKeywordOrIdent(expr) {
		return "DEFAULT", c.defaultCharset, false, nil
	}
	if cr, isCol := expr.(*nodes.ColumnRef); isCol && cr.Schema == "" && cr.Table == "" {
		raw = cr.Column
		return raw, normalizeCharsetName(raw), false, nil
	}
	ev, ok, evalErr := c.evalSetExprChecked(expr)
	if evalErr != nil {
		return "", "", false, evalErr
	}
	if ok {
		if ev.kind == userVarNull {
			return "NULL", "", true, nil
		}
		return ev.str, normalizeCharsetName(ev.str), false, nil
	}
	raw = nodeToSQLValue(expr)
	if strings.EqualFold(raw, "DEFAULT") {
		return raw, c.defaultCharset, false, nil
	}
	return raw, normalizeCharsetName(raw), false, nil
}

func (c *Catalog) effectiveGlobalClientCharset() string {
	cs := ""
	if c != nil {
		cs = c.globalCharsetClient
		if cs == "" {
			cs = c.defaultCharset
		}
	}
	return effectiveNonEmpty(strings.ToLower(strings.TrimSpace(normalizeCharsetName(cs))), "utf8mb4")
}

func (c *Catalog) effectiveGlobalCollationConn() string {
	if c != nil && c.globalCollationConn != "" {
		return c.globalCollationConn
	}
	return effectiveNonEmpty(defaultCollationForCharset[c.effectiveGlobalClientCharset()], "utf8mb4_0900_ai_ci")
}

func (c *Catalog) effectiveGlobalConnCharset() string {
	if cs := strings.ToLower(strings.TrimSpace(normalizeCharsetName(charsetForCollation(c.effectiveGlobalCollationConn())))); cs != "" {
		return cs
	}
	return "utf8mb4"
}

func (c *Catalog) effectiveSessionCollationConn() string {
	if c != nil {
		if c.collationConn != "" {
			return c.collationConn
		}
		if c.session.CollationConnection != "" {
			return c.session.CollationConnection
		}
	}
	return effectiveNonEmpty(defaultCollationForCharset[c.effectiveConnCharset()], "utf8mb4_0900_ai_ci")
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
