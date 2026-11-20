package catalog

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
	"unicode/utf8"

	nodes "github.com/bytebase/omni/mysql/ast"
	mysqlparser "github.com/bytebase/omni/mysql/parser"
	"github.com/bytebase/omni/mysql/validate"
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
			dbColl := ""
			if c.currentDB != "" {
				if db := c.GetDatabase(c.currentDB); db != nil {
					if db.Collation != "" {
						dbColl = canonicalMySQLCollationName(db.Collation)
					}
					if db.Charset != "" {
						dbCS = normalizeCharsetName(db.Charset)
					}
				}
			}
			if dbColl == "" {
				dbColl = defaultCollationForCharset[toLower(dbCS)]
			}
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
		"default_table_encryption": makeBoolSysvarSpec(
			false,
			func(c *Catalog) bool { return c.globalDefaultTableEncrypt },
			func(c *Catalog, v bool) { c.globalDefaultTableEncrypt = v },
			func(c *Catalog) bool { return c.session.DefaultTableEncryption },
			func(c *Catalog, v bool) { c.session.DefaultTableEncryption = v },
		),
		"innodb_directories": {
			kind: sysvarGlobalOnlyNonPersist,
			assign: func(_ *Catalog, varName string, scope setVarScope, _ nodes.ExprNode) (func(), error) {
				return checkReadOnlySysvarWrite(varName, scope, false)
			},
		},
		"innodb_page_size": {
			kind: sysvarGlobalOnlyNonPersist,
			assign: func(_ *Catalog, varName string, scope setVarScope, _ nodes.ExprNode) (func(), error) {
				return checkReadOnlySysvarWrite(varName, scope, false)
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
				return userVarValue{kind: userVarString, str: expandSQLModeForRead(c.globalSQLMode)}
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
		defMode := c.globalSQLMode
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
			if c.isDefinedStoredFuncCall(expr) {
				return nil, nil
			}
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
			if tok == "DEFAULT" {
				continue
			}
		}
		seen[tok] = true
	}
	var out []string
	for _, bf := range sqlModeBitFlags {
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
		if val.isFloat || strings.ContainsAny(s, ".eE") {
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
				continue
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
	if err := c.validateSetExprStatic(expr); err != nil {
		return userVarValue{}, false, err
	}
	return c.evalSetExprNode(expr)
}

func (c *Catalog) evalSetExprNode(expr nodes.ExprNode) (userVarValue, bool, error) {
	switch e := expr.(type) {
	case *nodes.ParenExpr:
		return c.evalSetExprNode(e.Expr)
	case *nodes.NullLit:
		return userVarValue{kind: userVarNull}, true, nil
	case *nodes.StringLit:
		return userVarValue{kind: userVarString, str: e.Value}, true, nil
	case *nodes.HexLit:
		raw, ok := decodeLiteralRawBytes(cleanHexDigits(e.Value), true)
		if !ok {
			return userVarValue{}, false, nil
		}
		return userVarValue{kind: userVarString, str: string(raw), binary: true}, true, nil
	case *nodes.IntLit:
		if s, ok := exactIntLiteral(e, ""); ok {
			return userVarValue{kind: userVarNumeric, str: s}, true, nil
		}
		return userVarValue{kind: userVarNumeric, str: fmt.Sprintf("%d", e.Value)}, true, nil
	case *nodes.FloatLit:
		return userVarValue{kind: userVarNumeric, str: e.Value, isFloat: strings.ContainsAny(e.Value, "eE")}, true, nil
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
		return userVarValue{kind: userVarString, str: string(raw), binary: true}, true, nil
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
				return userVarValue{kind: userVarNumeric, str: sign + op.Value, isFloat: strings.ContainsAny(op.Value, "eE")}, true, nil
			}
			v, ok, err := c.evalSetExprNode(e.Operand)
			if err != nil || !ok {
				return userVarValue{}, false, err
			}
			if v.kind == userVarNull {
				return v, true, nil
			}
			n, ok := parseSetNumericInt64(v)
			if ok && !v.isFloat {
				if e.Op == nodes.UnaryMinus {
					if n == math.MinInt64 {
						return userVarValue{}, false, errBigintOutOfRange()
					}
					n = -n
				}
				return userVarValue{kind: userVarNumeric, str: strconv.FormatInt(n, 10)}, true, nil
			}
			r, scale, ok := parseSetDecimalRat(v)
			if !ok {
				return userVarValue{}, false, nil
			}
			if e.Op == nodes.UnaryMinus {
				r.Neg(r)
			}
			return userVarValue{kind: userVarNumeric, str: formatSetDecimalRat(r, scale), isFloat: v.isFloat}, true, nil
		}
		return userVarValue{}, false, nil
	case *nodes.BinaryExpr:
		return c.evalSetBinaryExprNode(e)
	case *nodes.CaseExpr:
		return c.evalSetCaseExprNode(e)
	case *nodes.CastExpr:
		v, ok, err := c.evalSetExprNode(e.Expr)
		if err != nil || !ok {
			return userVarValue{}, false, err
		}
		out, ok := evalSetCastValue(v, e.TypeName, "")
		return out, ok, nil
	case *nodes.ConvertExpr:
		v, ok, err := c.evalSetExprNode(e.Expr)
		if err != nil || !ok {
			return userVarValue{}, false, err
		}
		out, ok := evalSetCastValue(v, e.TypeName, e.Charset)
		return out, ok, nil
	case *nodes.VariableRef:
		if !e.System {
			if v, exists := c.userVars[toLower(e.Name)]; exists {
				return v, true, nil
			}
			return userVarValue{kind: userVarNull}, true, nil
		}
		spec, ok, err := validateSetSysvarRef(e)
		if err != nil {
			return userVarValue{}, false, err
		}
		if !ok {
			return userVarValue{}, false, nil
		}
		if spec.kind != sysvarSessionAndGlobal {
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
		if e.Schema != "" || (e.SpacedParen && !c.isIgnoreSpaceSQLMode()) {
			targetDBName := e.Schema
			if targetDBName == "" {
				targetDBName = c.currentDB
			}
			if db := c.GetDatabase(targetDBName); db != nil && db.Functions[toLower(e.Name)] != nil {
				for _, arg := range e.Args {
					if _, _, err := c.evalSetExprNode(arg); err != nil {
						return userVarValue{}, false, err
					}
				}
				return userVarValue{}, false, nil
			}
			qual := e.Name
			if targetDBName != "" {
				qual = targetDBName + "." + e.Name
			}
			return userVarValue{}, false, errNoSuchFunction(qual)
		}
		switch toLower(e.Name) {
		case "if":
			if len(e.Args) != 3 {
				return userVarValue{}, false, errWrongParamCount("IF")
			}
			cond, ok, err := c.evalSetExprNode(e.Args[0])
			if err != nil || !ok {
				return userVarValue{}, false, err
			}
			chosen := e.Args[2]
			if evalSetValueTruthy(cond) {
				chosen = e.Args[1]
			}
			res, resOK, resErr := c.evalSetExprNode(chosen)
			if resErr != nil || !resOK {
				return res, resOK, resErr
			}
			return c.aggregateConditionalResultKind(res, e.Args[1], e.Args[2]), true, nil
		case "nullif":
			if len(e.Args) != 2 {
				return userVarValue{}, false, errWrongParamCount("NULLIF")
			}
			v0, ok, err := c.evalSetExprNode(e.Args[0])
			if err != nil || !ok {
				return userVarValue{}, false, err
			}
			v1, ok, err := c.evalSetExprNode(e.Args[1])
			if err != nil || !ok {
				return userVarValue{}, false, err
			}
			if v0.kind == userVarNull {
				return userVarValue{kind: userVarNull}, true, nil
			}
			if v1.kind != userVarNull && c.compareSetValues(v0, v1) == 0 {
				return userVarValue{kind: userVarNull}, true, nil
			}
			return v0, true, nil
		case "isnull":
			if len(e.Args) != 1 {
				return userVarValue{}, false, errWrongParamCount("ISNULL")
			}
			v, ok, err := c.evalSetExprNode(e.Args[0])
			if err != nil || !ok {
				return userVarValue{}, false, err
			}
			if v.kind == userVarNull {
				return userVarValue{kind: userVarNumeric, str: "1"}, true, nil
			}
			return userVarValue{kind: userVarNumeric, str: "0"}, true, nil
		case "concat":
			if len(e.Args) == 0 {
				return userVarValue{}, false, errWrongParamCount("CONCAT")
			}
			var b strings.Builder
			isBin := false
			for _, arg := range e.Args {
				v, ok, err := c.evalSetExprNode(arg)
				if err != nil || !ok {
					return userVarValue{}, false, err
				}
				if v.kind == userVarNull {
					return userVarValue{kind: userVarNull}, true, nil
				}
				if v.binary {
					isBin = true
				}
				b.WriteString(v.str)
			}
			return userVarValue{kind: userVarString, str: b.String(), binary: isBin}, true, nil
		case "coalesce":
			if len(e.Args) == 0 {
				return userVarValue{}, false, errWrongParamCount("COALESCE")
			}
			for _, arg := range e.Args {
				v, ok, err := c.evalSetExprNode(arg)
				if err != nil || !ok {
					return userVarValue{}, false, err
				}
				if v.kind != userVarNull {
					return c.aggregateConditionalResultKind(v, e.Args...), true, nil
				}
			}
			return userVarValue{kind: userVarNull}, true, nil
		case "ifnull":
			if len(e.Args) != 2 {
				return userVarValue{}, false, errWrongParamCount("IFNULL")
			}
			v0, ok, err := c.evalSetExprNode(e.Args[0])
			if err != nil || !ok {
				return userVarValue{}, false, err
			}
			if v0.kind != userVarNull {
				return c.aggregateConditionalResultKind(v0, e.Args...), true, nil
			}
			v1, ok, err := c.evalSetExprNode(e.Args[1])
			if err != nil || !ok {
				return userVarValue{}, false, err
			}
			return c.aggregateConditionalResultKind(v1, e.Args...), true, nil
		case "trim", "ltrim", "rtrim", "upper", "lower", "ucase", "lcase", "reverse":
			if len(e.Args) != 1 {
				return userVarValue{}, false, errWrongParamCount(strings.ToUpper(e.Name))
			}
			v, ok, err := c.evalSetExprNode(e.Args[0])
			if err != nil || !ok || v.kind == userVarNull {
				return v, ok, err
			}
			s := v.str
			switch toLower(e.Name) {
			case "trim":
				s = strings.Trim(s, " ")
			case "ltrim":
				s = strings.TrimLeft(s, " ")
			case "rtrim":
				s = strings.TrimRight(s, " ")
			case "upper", "ucase":
				if !v.binary {
					s = strings.ToUpper(s)
				}
			case "lower", "lcase":
				if !v.binary {
					s = strings.ToLower(s)
				}
			case "reverse":
				if v.binary {
					bs := []byte(s)
					for i, j := 0, len(bs)-1; i < j; i, j = i+1, j-1 {
						bs[i], bs[j] = bs[j], bs[i]
					}
					s = string(bs)
				} else {
					runes := []rune(s)
					for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
						runes[i], runes[j] = runes[j], runes[i]
					}
					s = string(runes)
				}
			}
			return userVarValue{kind: userVarString, str: s, binary: v.binary}, true, nil
		case "trim_leading", "trim_trailing", "trim_both", "trim_from":
			if len(e.Args) < 1 || len(e.Args) > 2 {
				return userVarValue{}, false, errWrongParamCount("TRIM")
			}
			rem := " "
			remBin := false
			strExpr := e.Args[0]
			if len(e.Args) == 2 {
				rv, ok, err := c.evalSetExprNode(e.Args[0])
				if err != nil || !ok {
					return userVarValue{}, false, err
				}
				if rv.kind == userVarNull {
					return userVarValue{kind: userVarNull}, true, nil
				}
				rem = rv.str
				remBin = rv.binary
				strExpr = e.Args[1]
			}
			sv, ok, err := c.evalSetExprNode(strExpr)
			if err != nil || !ok || sv.kind == userVarNull {
				return sv, ok, err
			}
			return userVarValue{kind: userVarString, str: trimSetString(toLower(e.Name), sv.str, rem), binary: sv.binary || remBin}, true, nil
		case "left", "right":
			if len(e.Args) != 2 {
				return userVarValue{}, false, errWrongParamCount(strings.ToUpper(e.Name))
			}
			sv, ok, err := c.evalSetExprNode(e.Args[0])
			if err != nil || !ok {
				return userVarValue{}, false, err
			}
			nv, ok, err := c.evalSetExprNode(e.Args[1])
			if err != nil || !ok {
				return userVarValue{}, false, err
			}
			if sv.kind == userVarNull || nv.kind == userVarNull {
				return userVarValue{kind: userVarNull}, true, nil
			}
			n, nok := parseSetNumericInt64(nv)
			if !nok {
				return userVarValue{}, false, nil
			}
			if n <= 0 {
				return userVarValue{kind: userVarString, str: "", binary: sv.binary}, true, nil
			}
			if sv.binary {
				if int(n) >= len(sv.str) {
					return userVarValue{kind: userVarString, str: sv.str, binary: true}, true, nil
				}
				if toLower(e.Name) == "left" {
					return userVarValue{kind: userVarString, str: sv.str[:n], binary: true}, true, nil
				}
				return userVarValue{kind: userVarString, str: sv.str[len(sv.str)-int(n):], binary: true}, true, nil
			}
			runes := []rune(sv.str)
			if int(n) >= len(runes) {
				return userVarValue{kind: userVarString, str: sv.str}, true, nil
			}
			if toLower(e.Name) == "left" {
				return userVarValue{kind: userVarString, str: string(runes[:n])}, true, nil
			}
			return userVarValue{kind: userVarString, str: string(runes[len(runes)-int(n):])}, true, nil
		case "substr", "substring", "mid":
			if len(e.Args) < 2 || len(e.Args) > 3 {
				return userVarValue{}, false, errWrongParamCount(strings.ToUpper(e.Name))
			}
			sv, ok, err := c.evalSetExprNode(e.Args[0])
			if err != nil || !ok {
				return userVarValue{}, false, err
			}
			pv, ok, err := c.evalSetExprNode(e.Args[1])
			if err != nil || !ok {
				return userVarValue{}, false, err
			}
			var lv userVarValue
			hasLen := len(e.Args) == 3
			if hasLen {
				lv, ok, err = c.evalSetExprNode(e.Args[2])
				if err != nil || !ok {
					return userVarValue{}, false, err
				}
			}
			if sv.kind == userVarNull || pv.kind == userVarNull || (hasLen && lv.kind == userVarNull) {
				return userVarValue{kind: userVarNull}, true, nil
			}
			pos, pok := parseSetNumericInt64(pv)
			if !pok {
				return userVarValue{}, false, nil
			}
			var length int64
			if hasLen {
				var lok bool
				length, lok = parseSetNumericInt64(lv)
				if !lok {
					return userVarValue{}, false, nil
				}
			}
			if sv.binary {
				startIdx := 0
				if pos > 0 {
					startIdx = int(pos - 1)
				} else if pos < 0 {
					startIdx = len(sv.str) + int(pos)
				} else {
					return userVarValue{kind: userVarString, str: "", binary: true}, true, nil
				}
				if startIdx < 0 || startIdx >= len(sv.str) {
					return userVarValue{kind: userVarString, str: "", binary: true}, true, nil
				}
				if !hasLen {
					return userVarValue{kind: userVarString, str: sv.str[startIdx:], binary: true}, true, nil
				}
				if length <= 0 {
					return userVarValue{kind: userVarString, str: "", binary: true}, true, nil
				}
				endIdx := len(sv.str)
				if length < int64(endIdx-startIdx) {
					endIdx = startIdx + int(length)
				}
				return userVarValue{kind: userVarString, str: sv.str[startIdx:endIdx], binary: true}, true, nil
			}
			runes := []rune(sv.str)
			startIdx := 0
			if pos > 0 {
				startIdx = int(pos - 1)
			} else if pos < 0 {
				startIdx = len(runes) + int(pos)
			} else {
				return userVarValue{kind: userVarString, str: ""}, true, nil
			}
			if startIdx < 0 || startIdx >= len(runes) {
				return userVarValue{kind: userVarString, str: ""}, true, nil
			}
			if !hasLen {
				return userVarValue{kind: userVarString, str: string(runes[startIdx:])}, true, nil
			}
			if length <= 0 {
				return userVarValue{kind: userVarString, str: ""}, true, nil
			}
			endIdx := len(runes)
			if length < int64(endIdx-startIdx) {
				endIdx = startIdx + int(length)
			}
			return userVarValue{kind: userVarString, str: string(runes[startIdx:endIdx])}, true, nil
		case "char_length", "character_length", "length":
			if len(e.Args) != 1 {
				return userVarValue{}, false, errWrongParamCount(strings.ToUpper(e.Name))
			}
			v, ok, err := c.evalSetExprNode(e.Args[0])
			if err != nil || !ok || v.kind == userVarNull {
				return v, ok, err
			}
			n := len(v.str)
			if !v.binary && toLower(e.Name) != "length" {
				n = utf8.RuneCountInString(v.str)
			}
			return userVarValue{kind: userVarNumeric, str: strconv.Itoa(n)}, true, nil
		case "replace":
			if len(e.Args) != 3 {
				return userVarValue{}, false, errWrongParamCount("REPLACE")
			}
			args := make([]userVarValue, 3)
			for i := 0; i < 3; i++ {
				v, ok, err := c.evalSetExprNode(e.Args[i])
				if err != nil || !ok {
					return userVarValue{}, false, err
				}
				if v.kind == userVarNull {
					return userVarValue{kind: userVarNull}, true, nil
				}
				args[i] = v
			}
			// REPLACE takes its character set from the first argument only:
			// REPLACE('ansi', X'78', '') = 'ANSI' compares case-insensitively.
			isBin := args[0].binary
			if args[1].str == "" {
				return userVarValue{kind: userVarString, str: args[0].str, binary: isBin}, true, nil
			}
			return userVarValue{kind: userVarString, str: strings.ReplaceAll(args[0].str, args[1].str, args[2].str), binary: isBin}, true, nil
		case "elt":
			// ELT(N, str1, ...) evaluates N and then only strN; N below 1 or
			// past the last string, or NULL, answers NULL.
			if len(e.Args) < 2 {
				return userVarValue{}, false, errWrongParamCount("ELT")
			}
			nv, ok, err := c.evalSetExprNode(e.Args[0])
			if err != nil || !ok {
				return userVarValue{}, false, err
			}
			if nv.kind == userVarNull {
				return userVarValue{kind: userVarNull}, true, nil
			}
			n, nok := setValInt(nv)
			if !nok {
				return userVarValue{}, false, nil
			}
			if n < 1 || n >= int64(len(e.Args)) {
				return userVarValue{kind: userVarNull}, true, nil
			}
			v, ok, err := c.evalSetExprNode(e.Args[n])
			if err != nil || !ok || v.kind == userVarNull {
				return v, ok, err
			}
			// The result's collation aggregates the strings, not N.
			_, isBin := c.setExprResultIsStringOrBinary(e)
			return userVarValue{kind: userVarString, str: v.str, binary: isBin || v.binary}, true, nil
		case "concat_ws":
			if len(e.Args) < 2 {
				return userVarValue{}, false, errWrongParamCount("CONCAT_WS")
			}
			sep, ok, err := c.evalSetExprNode(e.Args[0])
			if err != nil || !ok {
				return userVarValue{}, false, err
			}
			if sep.kind == userVarNull {
				return userVarValue{kind: userVarNull}, true, nil
			}
			isBin := sep.binary
			var parts []string
			for _, arg := range e.Args[1:] {
				v, ok, err := c.evalSetExprNode(arg)
				if err != nil || !ok {
					return userVarValue{}, false, err
				}
				if v.kind == userVarNull {
					continue
				}
				if v.binary {
					isBin = true
				}
				parts = append(parts, v.str)
			}
			return userVarValue{kind: userVarString, str: strings.Join(parts, sep.str), binary: isBin}, true, nil
		}
	}
	return userVarValue{}, false, nil
}

func (c *Catalog) isDefinedStoredFuncCall(expr nodes.ExprNode) bool {
	fc, ok := expr.(*nodes.FuncCallExpr)
	if !ok || fc == nil {
		return false
	}
	targetDBName := fc.Schema
	if targetDBName == "" {
		targetDBName = c.currentDB
	}
	db := c.GetDatabase(targetDBName)
	return db != nil && db.Functions[toLower(fc.Name)] != nil
}

func evalSetValueTruthy(v userVarValue) bool {
	if v.kind == userVarNull {
		return false
	}
	r, _, _, ok := parseSetNumberPrefixRat(v.str, false)
	return ok && r.Sign() != 0
}

func validateSetSysvarRef(x *nodes.VariableRef) (sysvarSpec, bool, error) {
	if strings.EqualFold(x.Scope, "PERSIST") || strings.EqualFold(x.Scope, "PERSIST_ONLY") {
		return sysvarSpec{}, false, &Error{
			Code:     1193,
			SQLState: "HY000",
			Message:  fmt.Sprintf("Unknown system variable '%s.%s'", strings.ToUpper(x.Scope), x.Name),
		}
	}
	lowerName := toLower(x.Name)
	spec, ok := mysqlSysvarSpecs[lowerName]
	if !ok {
		if !validate.IsKnownSystemVariable(lowerName) {
			return sysvarSpec{}, false, &Error{
				Code:     1193,
				SQLState: "HY000",
				Message:  fmt.Sprintf("Unknown system variable '%s'", x.Name),
			}
		}
		return sysvarSpec{}, false, nil
	}
	if spec.kind != sysvarSessionAndGlobal {
		if strings.EqualFold(x.Scope, "SESSION") || strings.EqualFold(x.Scope, "LOCAL") {
			return sysvarSpec{}, false, errIncorrectGlobalLocalVar(lowerName, "GLOBAL variable")
		}
	}
	return spec, true, nil
}

func trimSetString(mode, s, rem string) string {
	if rem == "" {
		return s
	}
	if mode == "trim_leading" || mode == "trim_both" || mode == "trim_from" {
		for strings.HasPrefix(s, rem) {
			s = strings.TrimPrefix(s, rem)
		}
	}
	if mode == "trim_trailing" || mode == "trim_both" || mode == "trim_from" {
		for strings.HasSuffix(s, rem) {
			s = strings.TrimSuffix(s, rem)
		}
	}
	return s
}

func (c *Catalog) validateSetExprStatic(n nodes.Node) error {
	if n == nil {
		return nil
	}
	var err error
	nodes.Inspect(n, func(cur nodes.Node) bool {
		if err != nil {
			return false
		}
		switch x := cur.(type) {
		case *nodes.VariableRef:
			if !x.System {
				return true
			}
			if _, _, err = validateSetSysvarRef(x); err != nil {
				return false
			}
		case *nodes.FuncCallExpr:
			if x.Schema != "" || (x.SpacedParen && !c.isIgnoreSpaceSQLMode()) {
				targetDBName := x.Schema
				if targetDBName == "" {
					targetDBName = c.currentDB
				}
				if db := c.GetDatabase(targetDBName); db != nil && db.Functions[toLower(x.Name)] != nil {
					return true
				}
				qual := x.Name
				if targetDBName != "" {
					qual = targetDBName + "." + x.Name
				}
				err = errNoSuchFunction(qual)
				return false
			}
			fn := toLower(x.Name)
			// A window function or an aggregate is refused in any branch:
			// mysql:8.0.46 answers 3593 and 1111 for IF(1, 'ANSI', COUNT(1)).
			if x.Over != nil || windowFuncs[fn] {
				err = &Error{Code: 3593, SQLState: "HY000",
					Message: fmt.Sprintf("You cannot use the window function '%s' in this context.'", fn)}
				return false
			}
			if functionalIndexAggregates[fn] {
				err = &Error{Code: 1111, SQLState: "HY000", Message: "Invalid use of group function"}
				return false
			}
			// Any other name is a stored function of the current database,
			// resolved before a branch is chosen: IF(1, 'ANSI', no_such()) is
			// 1305, and 1046 with no database selected.
			if !isBuiltinFunc(fn) {
				if c.currentDB == "" {
					err = errNoDatabaseSelected()
					return false
				}
				if db := c.GetDatabase(c.currentDB); db != nil && db.Functions[fn] != nil {
					return true
				}
				err = errNoSuchFunction(c.currentDB + "." + x.Name)
				return false
			}
			argc := len(x.Args)
			switch fn {
			case "if":
				if argc != 3 {
					err = errWrongParamCount("IF")
				}
			case "nullif", "ifnull", "left", "right":
				if argc != 2 {
					err = errWrongParamCount(strings.ToUpper(fn))
				}
			case "isnull", "trim", "ltrim", "rtrim", "upper", "lower", "ucase", "lcase",
				"reverse", "char_length", "character_length", "length":
				if argc != 1 {
					err = errWrongParamCount(strings.ToUpper(fn))
				}
			case "trim_leading", "trim_trailing", "trim_both":
				if argc < 1 || argc > 2 {
					err = errWrongParamCount("TRIM")
				}
			case "trim_from":
				if argc != 2 {
					err = errWrongParamCount("TRIM")
				}
			case "concat", "coalesce":
				if argc == 0 {
					err = errWrongParamCount(strings.ToUpper(fn))
				}
			case "concat_ws", "elt":
				if argc < 2 {
					err = errWrongParamCount(strings.ToUpper(fn))
				}
			case "substr", "substring", "mid":
				if argc < 2 || argc > 3 {
					err = errWrongParamCount(strings.ToUpper(fn))
				}
			case "replace":
				if argc != 3 {
					err = errWrongParamCount("REPLACE")
				}
			}
			if err != nil {
				return false
			}
		}
		return true
	})
	return err
}

func parseSetCoercedNumberRat(v userVarValue) *big.Rat {
	r, _, _, ok := parseSetNumberPrefixRat(v.str, false)
	if !ok {
		return new(big.Rat)
	}
	return r
}

func (c *Catalog) compareSetValues(left, right userVarValue) int {
	if left.kind == userVarNumeric || right.kind == userVarNumeric {
		lr := parseSetCoercedNumberRat(left)
		rr := parseSetCoercedNumberRat(right)
		return lr.Cmp(rr)
	}
	if left.binary || right.binary {
		return strings.Compare(left.str, right.str)
	}
	coll := c.effectiveSessionCollationConn()
	cs := charsetForCollation(coll)
	lk := collationComparisonKey(left.str, coll, cs)
	rk := collationComparisonKey(right.str, coll, cs)
	return strings.Compare(lk, rk)
}

func boolSetVar(b bool) userVarValue {
	if b {
		return userVarValue{kind: userVarNumeric, str: "1"}
	}
	return userVarValue{kind: userVarNumeric, str: "0"}
}

func evalSetCastValue(v userVarValue, dt *nodes.DataType, usingCharset string) (userVarValue, bool) {
	if v.kind == userVarNull {
		return userVarValue{kind: userVarNull}, true
	}
	if usingCharset != "" {
		return userVarValue{
			kind:   userVarString,
			str:    v.str,
			binary: strings.EqualFold(strings.TrimSpace(usingCharset), "binary"),
		}, true
	}
	if dt == nil {
		return userVarValue{}, false
	}
	switch strings.ToLower(strings.TrimSpace(dt.Name)) {
	case "char", "nchar", "varchar":
		isBin := strings.EqualFold(dt.Charset, "binary")
		s := v.str
		if dt.HasLength && dt.Length >= 0 {
			if isBin {
				if len(s) > dt.Length {
					s = s[:dt.Length]
				}
			} else {
				runes := []rune(s)
				if len(runes) > dt.Length {
					s = string(runes[:dt.Length])
				}
			}
		}
		return userVarValue{kind: userVarString, str: s, binary: isBin}, true
	case "binary", "varbinary":
		s := v.str
		if dt.HasLength && dt.Length >= 0 {
			if len(s) > dt.Length {
				s = s[:dt.Length]
			} else if strings.EqualFold(dt.Name, "binary") && len(s) < dt.Length {
				s += strings.Repeat("\x00", dt.Length-len(s))
			}
		}
		return userVarValue{kind: userVarString, str: s, binary: true}, true
	case "signed", "unsigned":
		s, ok := setCastInt(v, strings.EqualFold(strings.TrimSpace(dt.Name), "unsigned"))
		if !ok {
			return userVarValue{}, false
		}
		return userVarValue{kind: userVarNumeric, str: s}, true
	case "decimal", "numeric", "float", "double", "real":
		r := parseSetCoercedNumberRat(v)
		scale := dt.Scale
		if scale < 0 {
			scale = 0
		}
		return userVarValue{kind: userVarNumeric, str: formatSetDecimalRat(r, scale)}, true
	default:
		return userVarValue{}, false
	}
}

func (c *Catalog) setExprResultIsStringOrBinary(expr nodes.ExprNode) (isStr, isBin bool) {
	switch e := expr.(type) {
	case nil, *nodes.IntLit, *nodes.FloatLit, *nodes.BoolLit, *nodes.NullLit:
		return false, false
	case *nodes.ParenExpr:
		return c.setExprResultIsStringOrBinary(e.Expr)
	case *nodes.StringLit:
		return true, false
	case *nodes.HexLit, *nodes.BitLit:
		return true, true
	case *nodes.CastExpr:
		if e.TypeName == nil {
			return false, false
		}
		switch strings.ToLower(strings.TrimSpace(e.TypeName.Name)) {
		case "binary", "varbinary":
			return true, true
		case "char", "nchar", "varchar":
			return true, strings.EqualFold(e.TypeName.Charset, "binary")
		}
		return false, false
	case *nodes.ConvertExpr:
		if e.Charset != "" {
			return true, strings.EqualFold(strings.TrimSpace(e.Charset), "binary")
		}
		if e.TypeName != nil {
			return c.setExprResultIsStringOrBinary(&nodes.CastExpr{TypeName: e.TypeName})
		}
		return false, false
	case *nodes.VariableRef:
		if !e.System {
			if v, exists := c.userVars[toLower(e.Name)]; exists {
				return v.kind == userVarString, v.binary
			}
			return false, false
		}
		switch toLower(e.Name) {
		case "sql_mode", "time_zone", "system_time_zone", "character_set_client", "collation_connection":
			return true, false
		}
		return false, false
	case *nodes.CaseExpr:
		for _, w := range e.Whens {
			if w != nil {
				s, b := c.setExprResultIsStringOrBinary(w.Result)
				isStr = isStr || s
				isBin = isBin || b
			}
		}
		if e.Default != nil {
			s, b := c.setExprResultIsStringOrBinary(e.Default)
			isStr = isStr || s
			isBin = isBin || b
		}
		return isStr, isBin
	case *nodes.FuncCallExpr:
		switch toLower(e.Name) {
		case "concat", "concat_ws", "trim", "ltrim", "rtrim", "upper", "lower",
			"ucase", "lcase", "reverse", "trim_leading", "trim_trailing", "trim_both",
			"trim_from", "left", "right", "substr", "substring", "mid", "replace":
			for _, arg := range e.Args {
				_, b := c.setExprResultIsStringOrBinary(arg)
				isBin = isBin || b
			}
			return true, isBin
		case "elt":
			for i, arg := range e.Args {
				if i > 0 {
					_, b := c.setExprResultIsStringOrBinary(arg)
					isBin = isBin || b
				}
			}
			return true, isBin
		case "if":
			if len(e.Args) == 3 {
				s1, b1 := c.setExprResultIsStringOrBinary(e.Args[1])
				s2, b2 := c.setExprResultIsStringOrBinary(e.Args[2])
				return s1 || s2, b1 || b2
			}
		case "ifnull", "coalesce":
			for _, arg := range e.Args {
				s, b := c.setExprResultIsStringOrBinary(arg)
				isStr = isStr || s
				isBin = isBin || b
			}
			return isStr, isBin
		case "nullif":
			if len(e.Args) >= 1 {
				return c.setExprResultIsStringOrBinary(e.Args[0])
			}
		}
	}
	return false, false
}

func (c *Catalog) aggregateConditionalResultKind(v userVarValue, exprs ...nodes.ExprNode) userVarValue {
	if v.kind == userVarNull {
		return v
	}
	isStr, isBin := false, false
	for _, ex := range exprs {
		s, b := c.setExprResultIsStringOrBinary(ex)
		isStr = isStr || s
		isBin = isBin || b
	}
	if isStr || v.kind == userVarString {
		v.kind = userVarString
		if isBin {
			v.binary = true
		}
	}
	return v
}

func (c *Catalog) evalSetCaseExprNode(e *nodes.CaseExpr) (userVarValue, bool, error) {
	if e == nil {
		return userVarValue{}, false, nil
	}
	var hasOperand bool
	var opVal userVarValue
	if e.Operand != nil {
		var ok bool
		var err error
		opVal, ok, err = c.evalSetExprNode(e.Operand)
		if err != nil || !ok {
			return userVarValue{}, false, err
		}
		hasOperand = true
	}
	allResults := make([]nodes.ExprNode, 0, len(e.Whens)+1)
	for _, w := range e.Whens {
		if w != nil {
			allResults = append(allResults, w.Result)
		}
	}
	if e.Default != nil {
		allResults = append(allResults, e.Default)
	}
	for _, w := range e.Whens {
		if w == nil {
			continue
		}
		cond, ok, err := c.evalSetExprNode(w.Cond)
		if err != nil || !ok {
			return userVarValue{}, false, err
		}
		matched := false
		if hasOperand {
			matched = opVal.kind != userVarNull && cond.kind != userVarNull && c.compareSetValues(opVal, cond) == 0
		} else {
			matched = evalSetValueTruthy(cond)
		}
		if matched {
			res, resOK, resErr := c.evalSetExprNode(w.Result)
			if resErr != nil || !resOK {
				return res, resOK, resErr
			}
			return c.aggregateConditionalResultKind(res, allResults...), true, nil
		}
	}
	if e.Default != nil {
		res, resOK, resErr := c.evalSetExprNode(e.Default)
		if resErr != nil || !resOK {
			return res, resOK, resErr
		}
		return c.aggregateConditionalResultKind(res, allResults...), true, nil
	}
	return userVarValue{kind: userVarNull}, true, nil
}

func (c *Catalog) evalSetBinaryExprNode(e *nodes.BinaryExpr) (userVarValue, bool, error) {
	if e == nil {
		return userVarValue{}, false, nil
	}
	left, ok, err := c.evalSetExprNode(e.Left)
	if err != nil || !ok {
		return userVarValue{}, false, err
	}
	right, ok, err := c.evalSetExprNode(e.Right)
	if err != nil || !ok {
		return userVarValue{}, false, err
	}
	switch e.Op {
	case nodes.BinOpNullSafeEq:
		if left.kind == userVarNull && right.kind == userVarNull {
			return boolSetVar(true), true, nil
		}
		if left.kind == userVarNull || right.kind == userVarNull {
			return boolSetVar(false), true, nil
		}
		return boolSetVar(c.compareSetValues(left, right) == 0), true, nil
	case nodes.BinOpAnd:
		if (left.kind != userVarNull && !evalSetValueTruthy(left)) || (right.kind != userVarNull && !evalSetValueTruthy(right)) {
			return boolSetVar(false), true, nil
		}
		if left.kind == userVarNull || right.kind == userVarNull {
			return userVarValue{kind: userVarNull}, true, nil
		}
		return boolSetVar(true), true, nil
	case nodes.BinOpOr:
		if (left.kind != userVarNull && evalSetValueTruthy(left)) || (right.kind != userVarNull && evalSetValueTruthy(right)) {
			return boolSetVar(true), true, nil
		}
		if left.kind == userVarNull || right.kind == userVarNull {
			return userVarValue{kind: userVarNull}, true, nil
		}
		return boolSetVar(false), true, nil
	case nodes.BinOpXor:
		if left.kind == userVarNull || right.kind == userVarNull {
			return userVarValue{kind: userVarNull}, true, nil
		}
		return boolSetVar(evalSetValueTruthy(left) != evalSetValueTruthy(right)), true, nil
	}
	if left.kind == userVarNull || right.kind == userVarNull {
		return userVarValue{kind: userVarNull}, true, nil
	}
	switch e.Op {
	case nodes.BinOpEq:
		return boolSetVar(c.compareSetValues(left, right) == 0), true, nil
	case nodes.BinOpNe:
		return boolSetVar(c.compareSetValues(left, right) != 0), true, nil
	case nodes.BinOpLt:
		return boolSetVar(c.compareSetValues(left, right) < 0), true, nil
	case nodes.BinOpLe:
		return boolSetVar(c.compareSetValues(left, right) <= 0), true, nil
	case nodes.BinOpGt:
		return boolSetVar(c.compareSetValues(left, right) > 0), true, nil
	case nodes.BinOpGe:
		return boolSetVar(c.compareSetValues(left, right) >= 0), true, nil
	}
	var ln, rn int64
	var lok, rok bool
	if !left.isFloat && !right.isFloat {
		ln, lok = parseSetNumericInt64(left)
		rn, rok = parseSetNumericInt64(right)
	}
	if !lok || !rok {
		return evalSetDecimalBinaryExpr(e.Op, left, right)
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

var (
	setInt64Min  = big.NewInt(math.MinInt64)
	setInt64Max  = big.NewInt(math.MaxInt64)
	setUint64Max = new(big.Int).SetUint64(math.MaxUint64)
	setTwoTo64   = new(big.Int).Lsh(big.NewInt(1), 64)
)

// setCastInt is CAST(v AS SIGNED / UNSIGNED) as mysql:8.0.46 answers it. A
// string is read by its integer prefix (CAST('3.6' AS SIGNED) is 3, '1e2' is
// 1) and saturates at 2^64-1 or -2^63 before wrapping into the target. An
// exact number rounds half away from zero (CAST(3.6 AS SIGNED) is 4,
// CAST(-1.6 AS UNSIGNED) is 18446744073709551614); an integer up to 2^64-1
// wraps (CAST(18446744073709551615 AS SIGNED) is -1), a decimal saturates
// (CAST(9223372036854775807.5 AS SIGNED) is 9223372036854775807). A float
// rounds half to even and saturates at the signed bounds, for UNSIGNED too
// (CAST(1.5e19 AS UNSIGNED) is 9223372036854775807).
func setCastInt(v userVarValue, unsigned bool) (string, bool) {
	var n *big.Int
	wrapPositive := false
	switch v.kind {
	case userVarString:
		t := strings.TrimLeft(v.str, " \t")
		i := 0
		if i < len(t) && (t[i] == '+' || t[i] == '-') {
			i++
		}
		j := i
		for j < len(t) && t[j] >= '0' && t[j] <= '9' {
			j++
		}
		n = new(big.Int)
		if j > i {
			n.SetString(t[i:j], 10)
			if t[0] == '-' {
				n.Neg(n)
			}
		}
		if n.Cmp(setUint64Max) > 0 {
			n.Set(setUint64Max)
		}
		wrapPositive = true
	case userVarNumeric:
		r, _, _, ok := parseSetNumberPrefixRat(v.str, true)
		if !ok {
			return "", false
		}
		if v.isFloat {
			f, _ := r.Float64()
			f = math.RoundToEven(f)
			switch {
			case f >= math.MaxInt64:
				n = new(big.Int).Set(setInt64Max)
			case f <= math.MinInt64:
				n = new(big.Int).Set(setInt64Min)
			default:
				n = big.NewInt(int64(f))
			}
			break
		}
		num := new(big.Int).Abs(r.Num())
		num.Lsh(num, 1).Add(num, r.Denom())
		n = num.Quo(num, new(big.Int).Lsh(r.Denom(), 1))
		if r.Sign() < 0 {
			n.Neg(n)
		}
		integral := r.IsInt() && !strings.ContainsAny(v.str, ".eE")
		switch {
		case unsigned && n.Cmp(setUint64Max) > 0:
			n.Set(setUint64Max)
		case !unsigned && n.Cmp(setInt64Max) > 0 && !(integral && n.Cmp(setUint64Max) <= 0):
			n.Set(setInt64Max)
		}
		wrapPositive = integral
	default:
		return "", false
	}
	if n.Cmp(setInt64Min) < 0 {
		n.Set(setInt64Min)
	}
	if unsigned {
		if n.Sign() < 0 {
			n.Add(n, setTwoTo64)
		}
		return n.String(), true
	}
	if wrapPositive && n.Cmp(setInt64Max) > 0 {
		n.Sub(n, setTwoTo64)
	}
	return n.String(), true
}

// setValInt reads a value as MySQL's val_int does: a string by its integer
// prefix (' 2' is 2, '1.6' and '1e1' are 1, 'x' is 0); an exact number
// rounded half away from zero (1.5 is 2, 2.5 is 3); an approximate one
// rounded half to even, as rint does (2.5e0 is 2).
func setValInt(v userVarValue) (int64, bool) {
	switch v.kind {
	case userVarString:
		s := strings.TrimLeft(v.str, " \t")
		i := 0
		if i < len(s) && (s[i] == '+' || s[i] == '-') {
			i++
		}
		j := i
		for j < len(s) && s[j] >= '0' && s[j] <= '9' {
			j++
		}
		if j == i {
			return 0, true
		}
		n, err := strconv.ParseInt(s[:j], 10, 64)
		if err != nil {
			return 0, false
		}
		return n, true
	case userVarNumeric:
		if n, err := strconv.ParseInt(strings.TrimSpace(v.str), 10, 64); err == nil {
			return n, true
		}
		r, _, _, ok := parseSetNumberPrefixRat(v.str, true)
		if !ok {
			return 0, false
		}
		if v.isFloat {
			f, _ := r.Float64()
			f = math.RoundToEven(f)
			if f < math.MinInt64 || f >= math.MaxInt64 {
				return 0, false
			}
			return int64(f), true
		}
		// Half away from zero, exactly: (2|num| + den) / 2den.
		num := new(big.Int).Abs(r.Num())
		num.Lsh(num, 1).Add(num, r.Denom())
		q := num.Quo(num, new(big.Int).Lsh(r.Denom(), 1))
		if r.Sign() < 0 {
			q.Neg(q)
		}
		if !q.IsInt64() {
			return 0, false
		}
		return q.Int64(), true
	}
	return 0, false
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
		// An integer past the int64 range saturates, as val_int does:
		// SUBSTRING('ANSI', 2, 18446744073709551615) is 'NSI'.
		if ne, ok := err.(*strconv.NumError); ok && ne.Err == strconv.ErrRange {
			return n, true
		}
		return 0, false
	}
	return n, true
}

// parseSetNumberPrefixRat parses a leading or full MySQL decimal/scientific
// numeric literal (with exponent backtracking so '123euro' parses as 123).
func parseSetNumberPrefixRat(raw string, strictFull bool) (*big.Rat, int, bool, bool) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return nil, 0, false, false
	}
	i := 0
	neg := false
	if s[i] == '+' || s[i] == '-' {
		neg = s[i] == '-'
		i++
	}
	var digits strings.Builder
	sawDot := false
	fracDigits := 0
	for i < len(s) {
		ch := s[i]
		if ch >= '0' && ch <= '9' {
			digits.WriteByte(ch)
			if sawDot {
				fracDigits++
			}
			i++
			continue
		}
		if ch == '.' && !sawDot {
			sawDot = true
			i++
			continue
		}
		break
	}
	if digits.Len() == 0 {
		return nil, 0, false, false
	}
	hasExp := false
	expVal := 0
	if i < len(s) && (s[i] == 'e' || s[i] == 'E') {
		expPos := i + 1
		expNeg := false
		if expPos < len(s) && (s[expPos] == '+' || s[expPos] == '-') {
			expNeg = s[expPos] == '-'
			expPos++
		}
		expDigitStart := expPos
		for expPos < len(s) && s[expPos] >= '0' && s[expPos] <= '9' {
			if expVal <= 1000 {
				expVal = expVal*10 + int(s[expPos]-'0')
			}
			expPos++
		}
		if expPos > expDigitStart {
			hasExp = true
			if expNeg {
				expVal = -expVal
			}
			i = expPos
		}
	}
	if strictFull && i != len(s) {
		return nil, 0, false, false
	}
	num, ok := new(big.Int).SetString(digits.String(), 10)
	if !ok {
		return nil, 0, false, false
	}
	if num.Sign() == 0 {
		scale := fracDigits
		if hasExp {
			scale = 0
		}
		return new(big.Rat), scale, hasExp, true
	}
	totalExp := expVal - fracDigits
	if hasExp {
		if totalExp < -350 {
			return new(big.Rat), 0, true, true
		}
		if totalExp > 350 {
			totalExp = 350
		}
	}
	ten := big.NewInt(10)
	r := new(big.Rat)
	if totalExp >= 0 {
		pow := new(big.Int).Exp(ten, big.NewInt(int64(totalExp)), nil)
		r.SetInt(new(big.Int).Mul(num, pow))
	} else {
		pow := new(big.Int).Exp(ten, big.NewInt(int64(-totalExp)), nil)
		r.SetFrac(num, pow)
	}
	if neg {
		r.Neg(r)
	}
	scale := fracDigits
	if hasExp {
		scale = 0
	}
	return r, scale, hasExp, true
}

func parseSetDecimalRat(v userVarValue) (*big.Rat, int, bool) {
	if v.kind != userVarNumeric && v.kind != userVarString {
		return nil, 0, false
	}
	r, scale, _, ok := parseSetNumberPrefixRat(v.str, true)
	return r, scale, ok
}

func formatSetDecimalRat(r *big.Rat, scale int) string {
	if r == nil {
		return "0"
	}
	if scale < 0 {
		scale = 0
	}
	return r.FloatString(scale)
}

func evalSetDecimalBinaryExpr(op nodes.BinaryOp, left, right userVarValue) (userVarValue, bool, error) {
	if (left.kind != userVarNumeric && left.kind != userVarString) ||
		(right.kind != userVarNumeric && right.kind != userVarString) {
		return userVarValue{}, false, nil
	}
	lr, lScale, lFloat, lok := parseSetNumberPrefixRat(left.str, true)
	rr, rScale, rFloat, rok := parseSetNumberPrefixRat(right.str, true)
	if !lok || !rok {
		return userVarValue{}, false, nil
	}
	isFloat := left.isFloat || right.isFloat || lFloat || rFloat
	maxScale := lScale
	if rScale > maxScale {
		maxScale = rScale
	}
	out := new(big.Rat)
	formatFloatOut := func(q *big.Rat) (userVarValue, bool, error) {
		f, _ := q.Float64()
		return userVarValue{kind: userVarNumeric, str: strconv.FormatFloat(f, 'g', -1, 64), isFloat: true}, true, nil
	}
	switch op {
	case nodes.BinOpAdd:
		out.Add(lr, rr)
		if isFloat {
			return formatFloatOut(out)
		}
		return userVarValue{kind: userVarNumeric, str: formatSetDecimalRat(out, maxScale)}, true, nil
	case nodes.BinOpSub:
		out.Sub(lr, rr)
		if isFloat {
			return formatFloatOut(out)
		}
		return userVarValue{kind: userVarNumeric, str: formatSetDecimalRat(out, maxScale)}, true, nil
	case nodes.BinOpMul:
		out.Mul(lr, rr)
		if isFloat {
			return formatFloatOut(out)
		}
		return userVarValue{kind: userVarNumeric, str: formatSetDecimalRat(out, lScale+rScale)}, true, nil
	case nodes.BinOpDiv:
		if rr.Sign() == 0 {
			return userVarValue{kind: userVarNull}, true, nil
		}
		out.Quo(lr, rr)
		if isFloat {
			return formatFloatOut(out)
		}
		divScale := lScale + 4
		if divScale < 4 {
			divScale = 4
		}
		return userVarValue{kind: userVarNumeric, str: formatSetDecimalRat(out, divScale)}, true, nil
	case nodes.BinOpDivInt:
		if rr.Sign() == 0 {
			return userVarValue{kind: userVarNull}, true, nil
		}
		out.Quo(lr, rr)
		q := new(big.Int).Quo(out.Num(), out.Denom())
		return userVarValue{kind: userVarNumeric, str: q.String()}, true, nil
	case nodes.BinOpMod:
		if rr.Sign() == 0 {
			return userVarValue{kind: userVarNull}, true, nil
		}
		out.Quo(lr, rr)
		q := new(big.Rat).SetInt(new(big.Int).Quo(out.Num(), out.Denom()))
		q.Mul(q, rr)
		out.Sub(lr, q)
		if isFloat {
			return formatFloatOut(out)
		}
		return userVarValue{kind: userVarNumeric, str: formatSetDecimalRat(out, maxScale)}, true, nil
	default:
		return userVarValue{}, false, nil
	}
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
