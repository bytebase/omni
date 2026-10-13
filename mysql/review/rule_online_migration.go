package review

import (
	"encoding/json"
	"fmt"
	"regexp"
	"slices"
	"strconv"

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
	"max-load":                         "string",
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
	"throttle-control-replicas":        "string",
	"attempt-instant-ddl":              "bool",
	"assume-master-host":               "bool",
	"skip-metadata-lock-check":         "bool",
}

// checkOnlineMigration reports, from the SQL alone, what fails a change
// that asks for online migration before gh-ost starts: a directive whose
// flags Bytebase cannot read, and a sheet that is not one ALTER TABLE
// statement.
func checkOnlineMigration(sql string, stmts []statement, r *reporter) {
	if msg := ghostFlagsRefused(sql); msg != "" {
		r.report(review.OnlineMigration, -1, review.Range{}, msg)
	}
	switch {
	case len(stmts) != 1:
		r.report(review.OnlineMigration, -1, review.Range{}, fmt.Sprintf("online migration runs the change as one ALTER TABLE statement, and the change has %d statements", len(stmts)))
	default:
		if _, ok := stmts[0].node.(*ast.AlterTableStmt); !ok {
			r.report(review.OnlineMigration, stmts[0].index, rangeOf(stmts[0].loc), "online migration runs only an ALTER TABLE statement")
		}
	}
}

// ghostFlagsRefused returns why Bytebase cannot read the flags of the
// sheet's gh-ost directive, "" when it can: a flag object that is not
// JSON of strings, a flag it does not pass, a value it cannot read as its
// type, or switch-to-rbr with assume-rbr. The value of max-load is checked
// by gh-ost, and is not judged.
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
		}
		if err != nil {
			return fmt.Sprintf("online migration cannot read gh-ost flag %q as %s: %q", k, kind, v)
		}
	}
	if isTrue(flags["switch-to-rbr"]) && isTrue(flags["assume-rbr"]) {
		return "online migration cannot take gh-ost flags switch-to-rbr and assume-rbr together"
	}
	return ""
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
	// gh-ost migrates the table of the task's database.
	tbl := x.table(v.Table, false)
	if tbl == nil {
		return
	}
	at := func(msg string) {
		r.report(review.OnlineMigration, stmts[0].index, rangeOf(stmts[0].loc), msg)
	}
	name := table(v.Table)
	if !slices.ContainsFunc(tbl.GetIndexes(), func(i *metadata.IndexMetadata) bool { return i.GetPrimary() || i.GetUnique() }) {
		at("online migration needs a primary or unique key on " + name + ", which has none")
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
