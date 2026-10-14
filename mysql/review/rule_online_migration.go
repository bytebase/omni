package review

import (
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/review"
)

// Bytebase runs online migration through gh-ost. It hands gh-ost the whole
// sheet, without its directives, as the one ALTER TABLE statement gh-ost
// applies, with the flags of the sheet's gh-ost directive; gh-ost stages
// its tables in the bbdataarchive database, and refuses a table it cannot
// copy row by row.
//
// bytebase: backend/runner/taskrun/database_migrate_executor.go —
// executeGhostMigration; backend/component/ghost — ParseGhostDirective,
// GetUserFlags, NewMigrationContext, GetTableNameFromStatement.

// ghostDirective matches the sheet's gh-ost directive, as Bytebase does.
var ghostDirective = regexp.MustCompile(`(?im)^\s*--\s*gh-ost\s*=\s*(\{[^}]*\})\s*(?:/\*.*\*/)?\s*$`)

// ghostFlagKinds are the flags Bytebase passes to gh-ost, by the type it
// reads each as.
var ghostFlagKinds = map[string]string{
	"max-load":                         "load map",
	"chunk-size":                       "int",
	"dml-batch-size":                   "int",
	"default-retries":                  "int",
	"cut-over-lock-timeout-seconds":    "int",
	"exponential-backoff-max-interval": "int",
	"max-lag-millis":                   "int",
	"allow-on-master":                  "bool",
	"switch-to-rbr":                    "bool",
	"assume-rbr":                       "bool",
	"heartbeat-interval-millis":        "int",
	"nice-ratio":                       "float",
	"throttle-control-replicas":        "host list",
	"attempt-instant-ddl":              "bool",
	"assume-master-host":               "bool",
	"skip-metadata-lock-check":         "bool",
}

// checkOnlineMigration reports, from the SQL alone, what fails a change
// that asks for online migration before gh-ost starts: a directive whose
// flags Bytebase cannot read, a sheet that is not one ALTER TABLE
// statement, and an ALTER TABLE gh-ost refuses (gh-ost go/logic/migrator.go,
// validateAlterStatement): one that renames the table, and one that renames
// a column by CHANGE, which gh-ost takes only with approve-renamed-columns
// or skip-renamed-columns, flags Bytebase does not pass. gh-ost compares
// the two names of a CHANGE exactly, and splits a name at whitespace, so
// two different names are a rename to it either way. It also reports
// RENAME COLUMN, which
// gh-ost does not recognize as a rename (go/sql/parser.go,
// renameColumnRegexp reads only CHANGE): it copies no rows into the new
// name, and the cut-over loses the column's data.
func checkOnlineMigration(sql string, stmts []statement, r *reporter) {
	if msg := ghostFlagsRefused(sql); msg != "" {
		r.report(review.OnlineMigration, -1, review.Range{}, msg)
	}
	switch {
	case len(stmts) != 1:
		r.report(review.OnlineMigration, -1, review.Range{}, fmt.Sprintf("online migration runs the change as one ALTER TABLE statement, and the change has %d statements", len(stmts)))
	default:
		v, ok := stmts[0].node.(*ast.AlterTableStmt)
		if !ok {
			r.report(review.OnlineMigration, stmts[0].index, rangeOf(stmts[0].loc), "online migration runs only an ALTER TABLE statement")
			break
		}
		for _, cmd := range v.Commands {
			if cmd == nil {
				continue
			}
			switch {
			case cmd.Type == ast.ATRenameTable:
				r.report(review.OnlineMigration, stmts[0].index, rangeOf(cmd.Loc), "gh-ost does not rename the table it migrates")
			case cmd.Type == ast.ATChangeColumn && cmd.Column != nil && cmd.Name != cmd.Column.Name:
				r.report(review.OnlineMigration, stmts[0].index, rangeOf(cmd.Loc), "gh-ost renames column "+ident(cmd.Name)+" to "+ident(cmd.Column.Name)+" only with approve-renamed-columns, which Bytebase does not pass")
			case cmd.Type == ast.ATRenameColumn:
				r.report(review.OnlineMigration, stmts[0].index, rangeOf(cmd.Loc), "gh-ost does not read RENAME COLUMN as a rename, and the migrated table loses the data of column "+ident(cmd.Name))
			}
		}
	}
}

// ghostFlagsRefused returns why Bytebase cannot read the flags of the
// sheet's gh-ost directive, "" when it can: a flag object that is not
// JSON of strings, a flag it does not pass, a value it cannot read as its
// type, or switch-to-rbr with assume-rbr. max-load and
// throttle-control-replicas are read as gh-ost reads them (go/base,
// ParseLoadMap; go/mysql, ReadCommaDelimitedList).
func ghostFlagsRefused(sql string) string {
	match := ghostDirective.FindStringSubmatch(sql)
	if len(match) < 2 {
		return ""
	}
	var flags map[string]string
	if err := json.Unmarshal([]byte(match[1]), &flags); err != nil {
		return "online migration cannot read the gh-ost directive: " + err.Error()
	}
	keys := make([]string, 0, len(flags))
	for k := range flags {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	for _, k := range keys {
		kind, ok := ghostFlagKinds[k]
		if !ok {
			return "online migration does not support gh-ost flag " + strconv.Quote(k)
		}
		v := flags[k]
		var err error
		switch kind {
		case "int":
			_, err = strconv.ParseInt(v, 10, 64)
		case "bool":
			_, err = strconv.ParseBool(v)
		case "float":
			_, err = strconv.ParseFloat(v, 64)
		case "load map":
			err = parseLoadMap(v)
		case "host list":
			err = parseHostList(v)
		}
		if err != nil {
			return fmt.Sprintf("online migration cannot read gh-ost flag %q as %s: %q", k, kind, v)
		}
		if lo, hi, ok := ghostFlagRange(k); ok {
			if n, _ := strconv.ParseInt(v, 10, 64); n < lo || n > hi {
				bound := fmt.Sprintf("below %d", lo)
				if hi < math.MaxInt64 {
					bound = fmt.Sprintf("outside %d to %d", lo, hi)
				}
				return fmt.Sprintf("online migration cannot take gh-ost flag %q = %d, %s", k, n, bound)
			}
		}
	}
	if isTrue(flags["switch-to-rbr"]) && isTrue(flags["assume-rbr"]) {
		return "online migration cannot take gh-ost flags switch-to-rbr and assume-rbr together"
	}
	return ""
}

// ghostFlagRange returns the bounds gh-ost's setter of an integer flag
// holds a value to, failing Bytebase's migration context outside them.
//
// gh-ost (bytebase/gh-ost2): go/base/context.go —
// SetCutOverLockTimeoutSeconds, SetExponentialBackoffMaxInterval.
func ghostFlagRange(flag string) (lo, hi int64, ok bool) {
	switch flag {
	case "cut-over-lock-timeout-seconds":
		return 1, 10, true
	case "exponential-backoff-max-interval":
		return 2, math.MaxInt64, true
	}
	return 0, 0, false
}

// parseLoadMap reads a max-load value as gh-ost's ParseLoadMap does:
// comma-separated status=number pairs.
func parseLoadMap(v string) error {
	if v == "" {
		return nil
	}
	for _, cond := range strings.Split(v, ",") {
		kv := strings.Split(cond, "=")
		if len(kv) != 2 || kv[0] == "" {
			return fmt.Errorf("bad condition %q", cond)
		}
		if _, err := strconv.ParseInt(kv[1], 10, 0); err != nil {
			return err
		}
	}
	return nil
}

// ghostHost matches one host of a throttle-control-replicas value as
// gh-ost's NewRawInstanceKey does: host, host:port, [ipv6]:port, or ipv6.
var ghostHost = regexp.MustCompile(`^(?:[^:]+(?::([0-9]+))?|\[[:0-9a-fA-F]+\]:([0-9]+)|[:0-9a-fA-F]+)$`)

// parseHostList reads a throttle-control-replicas value as gh-ost's
// ReadCommaDelimitedList does.
func parseHostList(v string) error {
	if v == "" {
		return nil
	}
	for _, host := range strings.Split(v, ",") {
		m := ghostHost.FindStringSubmatch(host)
		if m == nil {
			return fmt.Errorf("cannot parse address %q", host)
		}
		for _, port := range m[1:] {
			if port == "" {
				continue
			}
			if _, err := strconv.Atoi(port); err != nil {
				return err
			}
		}
	}
	return nil
}

func isTrue(v string) bool {
	b, err := strconv.ParseBool(v)
	return err == nil && b
}

// checkOnlineMigrationTarget reports what fails online migration on a
// target: a missing bbdataarchive database, and an ALTER TABLE of a synced
// table gh-ost refuses to copy, one without a primary or unique key, with
// a foreign key of its own or one another table of the database holds to
// it, or with a trigger.
//
// gh-ost (bytebase/gh-ost2): go/logic/inspect.go —
// InspectTableColumnsAndUniqueKeys, validateTableForeignKeys,
// validateTableTriggers.
//
// gh-ost also refuses a view, and a table whose unique keys are all of
// expressions, which its key query cannot see, or all hold a nullable
// column, or a FLOAT or JSON one, which it skips: it takes the primary key, else a unique key without one, first,
// and bails out on a nullable key without allow-nullable-unique-key, a
// flag Bytebase does not pass.
func checkOnlineMigrationTarget(stmts []statement, t review.Target, x schemaIndex, r *reporter) {
	if !t.BackupDatabaseExists {
		r.report(review.OnlineMigration, -1, review.Range{}, "online migration needs database "+backupDatabase+", which does not exist")
	}
	if len(stmts) != 1 {
		return
	}
	v, ok := stmts[0].node.(*ast.AlterTableStmt)
	if !ok {
		return
	}
	// gh-ost migrates the table of the name in the task's database, which
	// Bytebase passes as gh-ost's DatabaseName whatever database the ALTER
	// names (backend/component/ghost, NewMigrationContext).
	tbl := x.table(&ast.TableRef{Name: v.Table.Name}, false)
	if tbl == nil {
		// gh-ost's validateTable refuses a view.
		if x.isView(v.Table.Name) {
			r.report(review.OnlineMigration, stmts[0].index, rangeOf(stmts[0].loc), "online migration needs a table, and "+table(v.Table)+" is a view")
		}
		return
	}
	at := func(msg string) {
		r.report(review.OnlineMigration, stmts[0].index, rangeOf(stmts[0].loc), msg)
	}
	name := table(v.Table)
	switch {
	case !slices.ContainsFunc(tbl.GetIndexes(), func(i *metadata.IndexMetadata) bool { return i.GetPrimary() || i.GetUnique() && !functionalKey(tbl, i) }):
		at("online migration needs a primary or unique key on " + name + ", which has none")
	// One key must be all of these at once.
	case !slices.ContainsFunc(tbl.GetIndexes(), func(i *metadata.IndexMetadata) bool {
		return (i.GetPrimary() || i.GetUnique() && !nullableKey(tbl, i)) && !functionalKey(tbl, i) && !floatOrJSONKey(tbl, i)
	}):
		at("online migration needs a primary key or a unique key of NOT NULL columns, none of them FLOAT or JSON, on " + name + ", which has none")
	}
	if len(tbl.GetForeignKeys()) > 0 {
		at("online migration does not support " + name + ", which has a foreign key")
	} else if x.referenced(tbl.GetName()) {
		at("online migration does not support " + name + ", which a foreign key references")
	}
	if len(tbl.GetTriggers()) > 0 {
		at("online migration does not support " + name + ", which has a trigger")
	}
}

// functionalKey reports whether every part of an index is an expression,
// one naming no column of the table: gh-ost finds its keys by their
// columns in information_schema.STATISTICS, where such a part has none.
func functionalKey(t *metadata.TableMetadata, idx *metadata.IndexMetadata) bool {
	parts := idx.GetExpressions()
	if len(parts) == 0 {
		return false
	}
	for _, part := range parts {
		if slices.ContainsFunc(t.GetColumns(), func(c *metadata.ColumnMetadata) bool { return strings.EqualFold(c.GetName(), part) }) {
			return false
		}
	}
	return true
}

// floatOrJSONKey reports whether an index holds a column whose type gh-ost
// will not copy by, FLOAT or JSON: it skips such a key.
//
// gh-ost (bytebase/gh-ost2): go/logic/inspect.go — applyColumnTypes,
// InspectOriginalAndGhostTables.
func floatOrJSONKey(t *metadata.TableMetadata, idx *metadata.IndexMetadata) bool {
	for _, part := range idx.GetExpressions() {
		for _, c := range t.GetColumns() {
			// gh-ost reads ENUM and BINARY types after FLOAT and JSON, so
			// ENUM('float') is an ENUM.
			typ := strings.ToLower(c.GetType())
			typed := strings.Contains(typ, "float") || strings.Contains(typ, "json")
			if strings.HasPrefix(typ, "enum") || strings.HasPrefix(typ, "binary") {
				typed = false
			}
			if strings.EqualFold(c.GetName(), part) && typed {
				return true
			}
		}
	}
	return false
}

// nullableKey reports whether an index certainly holds a nullable column:
// one of its parts names a column of the table that is nullable. A part
// that names no column, an expression, is not judged.
func nullableKey(t *metadata.TableMetadata, idx *metadata.IndexMetadata) bool {
	for _, part := range idx.GetExpressions() {
		for _, c := range t.GetColumns() {
			if strings.EqualFold(c.GetName(), part) && c.GetNullable() {
				return true
			}
		}
	}
	return false
}
