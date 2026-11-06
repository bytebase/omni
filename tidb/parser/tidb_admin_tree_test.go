package parser

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/bytebase/omni/tidb/ast"
)

// What these statements leave in the tree, which a parse-only verdict cannot see.
// `SHOW BR JOB QUERY 1` and `SHOW BR JOB 1` read alike, `SET CONFIG 'a' b = 1` must
// not hold an instance as a component type, and `KILL TIDB QUERY 1` has three things
// to record about one number.

func parseOneTiDBAdmin(t *testing.T, sql string) ast.Node {
	t.Helper()
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	if len(list.Items) != 1 {
		t.Fatalf("Parse(%q): %d statements, want 1", sql, len(list.Items))
	}
	return list.Items[0]
}

func TestRecommendIndexInTree(t *testing.T) {
	for _, tc := range []struct {
		sql    string
		action ast.RecommendIndexAction
		hasSQL bool
		text   string
		id     int64
		opts   []string
	}{
		{sql: "RECOMMEND INDEX RUN", action: ast.RecommendIndexRun},
		{sql: "recommend index run for 'select 1'", action: ast.RecommendIndexRun, hasSQL: true, text: "select 1"},
		{sql: "RECOMMEND INDEX RUN FOR ''", action: ast.RecommendIndexRun, hasSQL: true},
		{sql: "RECOMMEND INDEX RUN WITH a = 1, b = 'x'", action: ast.RecommendIndexRun, opts: []string{"a", "b"}},
		{sql: "RECOMMEND INDEX RUN FOR 'q' WITH a = 1", action: ast.RecommendIndexRun, hasSQL: true, text: "q", opts: []string{"a"}},
		{sql: "RECOMMEND INDEX SHOW", action: ast.RecommendIndexShow},
		{sql: "RECOMMEND INDEX APPLY 7", action: ast.RecommendIndexApply, id: 7},
		{sql: "RECOMMEND INDEX IGNORE 9", action: ast.RecommendIndexIgnore, id: 9},
		{sql: "RECOMMEND INDEX SET a = 1, `b c` = 2", action: ast.RecommendIndexSet, opts: []string{"a", "b c"}},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			stmt, ok := parseOneTiDBAdmin(t, tc.sql).(*ast.RecommendIndexStmt)
			if !ok {
				t.Fatalf("got %T", parseOneTiDBAdmin(t, tc.sql))
			}
			if stmt.Action != tc.action || stmt.HasSQL != tc.hasSQL || stmt.SQL != tc.text || stmt.ID != tc.id {
				t.Errorf("got action=%d hasSQL=%v sql=%q id=%d", stmt.Action, stmt.HasSQL, stmt.SQL, stmt.ID)
			}
			var names []string
			for _, o := range stmt.Options {
				names = append(names, o.Name)
			}
			if !reflect.DeepEqual(names, tc.opts) {
				t.Errorf("options = %v, want %v", names, tc.opts)
			}
		})
	}
}

// TestRecommendIndexOptionValues pins what a Literal becomes, the concatenation and
// the character set above all.
func TestRecommendIndexOptionValues(t *testing.T) {
	for _, tc := range []struct {
		value string
		check func(ast.ExprNode) bool
	}{
		{"1", func(e ast.ExprNode) bool { v, ok := e.(*ast.IntLit); return ok && v.Value == 1 }},
		{"1.5", func(e ast.ExprNode) bool { v, ok := e.(*ast.FloatLit); return ok && v.Value == "1.5" }},
		{"TRUE", func(e ast.ExprNode) bool { v, ok := e.(*ast.BoolLit); return ok && v.Value }},
		{"FALSE", func(e ast.ExprNode) bool { v, ok := e.(*ast.BoolLit); return ok && !v.Value }},
		{"NULL", func(e ast.ExprNode) bool { _, ok := e.(*ast.NullLit); return ok }},
		{"0x10", func(e ast.ExprNode) bool { _, ok := e.(*ast.HexLit); return ok }},
		{"b'1'", func(e ast.ExprNode) bool { _, ok := e.(*ast.BitLit); return ok }},
		{"'x'", func(e ast.ExprNode) bool { v, ok := e.(*ast.StringLit); return ok && v.Value == "x" }},
		{"'x' 'y' 'z'", func(e ast.ExprNode) bool { v, ok := e.(*ast.StringLit); return ok && v.Value == "xyz" }},
		{"_utf8mb4 'x'", func(e ast.ExprNode) bool {
			v, ok := e.(*ast.StringLit)
			return ok && v.Value == "x" && v.Charset == "_utf8mb4"
		}},
		{"N'x'", func(e ast.ExprNode) bool {
			v, ok := e.(*ast.StringLit)
			return ok && v.Value == "x" && v.Charset == "_utf8"
		}},
		{"_utf8mb4 0x10", func(e ast.ExprNode) bool { _, ok := e.(*ast.HexLit); return ok }},
	} {
		t.Run(tc.value, func(t *testing.T) {
			stmt := parseOneTiDBAdmin(t, "RECOMMEND INDEX SET a = "+tc.value).(*ast.RecommendIndexStmt)
			if got := stmt.Options[0].Value; !tc.check(got) {
				t.Errorf("got %s", ast.NodeToString(got))
			}
		})
	}
}

func TestSetConfigInTree(t *testing.T) {
	for _, tc := range []struct {
		sql         string
		typ         string
		hasInstance bool
		instance    string
		name        string
		value       string
	}{
		{sql: "SET CONFIG TIKV LOG.LEVEL = 'info'", typ: "tikv", name: "LOG.LEVEL", value: "info"},
		{sql: "set config '127.0.0.1:3306' auto-compaction-mode := 'x'", hasInstance: true, instance: "127.0.0.1:3306", name: "auto-compaction-mode", value: "x"},
		{sql: "SET CONFIG `PD` `a`.`b`-c.1 = 'v'", typ: "pd", name: "a.b-c.1", value: "v"},
		{sql: "SET CONFIG '' a.select = 'v'", hasInstance: true, name: "a.select", value: "v"},
		{sql: "SET CONFIG tikv a = ON", typ: "tikv", name: "a", value: "ON"},
		{sql: "SET CONFIG tikv a = BINARY", typ: "tikv", name: "a", value: "BINARY"},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			stmt, ok := parseOneTiDBAdmin(t, tc.sql).(*ast.SetConfigStmt)
			if !ok {
				t.Fatalf("got %T", parseOneTiDBAdmin(t, tc.sql))
			}
			if stmt.Type != tc.typ || stmt.HasInstance != tc.hasInstance || stmt.Instance != tc.instance || stmt.Name != tc.name {
				t.Errorf("got type=%q hasInstance=%v instance=%q name=%q", stmt.Type, stmt.HasInstance, stmt.Instance, stmt.Name)
			}
			lit, ok := stmt.Value.(*ast.StringLit)
			if !ok || lit.Value != tc.value {
				t.Errorf("value = %s, want %q", ast.NodeToString(stmt.Value), tc.value)
			}
		})
	}
	// DEFAULT is not a string, and a variable called config is not this statement.
	if stmt := parseOneTiDBAdmin(t, "SET CONFIG tikv a = DEFAULT").(*ast.SetConfigStmt); stmt.Value == nil {
		t.Error("DEFAULT left no value")
	} else if _, ok := stmt.Value.(*ast.DefaultExpr); !ok {
		t.Errorf("DEFAULT = %T", stmt.Value)
	}
	if _, ok := parseOneTiDBAdmin(t, "SET config = 1").(*ast.SetStmt); !ok {
		t.Error("SET config = 1 is not an assignment")
	}
}

// TestBRIELogTasksAndJobsInTree asserts the kind, the storage, the options and the job
// id of each log task and job statement, all of which land in one node.
func TestBRIELogTasksAndJobsInTree(t *testing.T) {
	for _, tc := range []struct {
		sql     string
		kind    ast.BRIEKind
		storage string
		options []string
		job     int64
	}{
		{sql: "STOP BACKUP LOGS", kind: ast.BRIEStopLogs},
		{sql: "PAUSE BACKUP LOGS", kind: ast.BRIEPauseLogs},
		{sql: "PAUSE BACKUP LOGS gc_ttl = '20220304' concurrency 4", kind: ast.BRIEPauseLogs, options: []string{"GC_TTL", "CONCURRENCY"}},
		{sql: "RESUME BACKUP LOGS", kind: ast.BRIEResumeLogs},
		{sql: "PURGE BACKUP LOGS FROM 'noop://'", kind: ast.BRIEPurgeLogs, storage: "noop://"},
		{sql: "PURGE BACKUP LOGS FROM 's3://b/' until_ts = '2012122304'", kind: ast.BRIEPurgeLogs, storage: "s3://b/", options: []string{"UNTIL_TS"}},
		{sql: "SHOW BACKUP LOGS STATUS", kind: ast.BRIEShowLogsStatus},
		{sql: "SHOW BACKUP LOGS METADATA FROM 'noop://'", kind: ast.BRIEShowLogsMetadata, storage: "noop://"},
		{sql: "SHOW BACKUP METADATA FROM 'noop://x'", kind: ast.BRIEShowBackupMetadata, storage: "noop://x"},
		{sql: "SHOW BR JOB 1234", kind: ast.BRIEShowJob, job: 1234},
		{sql: "SHOW BR JOB QUERY 1234", kind: ast.BRIEShowJobQuery, job: 1234},
		{sql: "SHOW BR JOB 9223372036854775807", kind: ast.BRIEShowJob, job: 9223372036854775807},
		{sql: "CANCEL BR JOB 77", kind: ast.BRIECancelJob, job: 77},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			stmt, ok := parseOneTiDBAdmin(t, tc.sql).(*ast.BRIEStmt)
			if !ok {
				t.Fatalf("got %T", parseOneTiDBAdmin(t, tc.sql))
			}
			if stmt.Kind != tc.kind {
				t.Errorf("kind = %d, want %d", stmt.Kind, tc.kind)
			}
			if !stmt.Kind.IsLogTaskOrJob() {
				t.Error("IsLogTaskOrJob = false")
			}
			if stmt.Target != ast.BRIETargetNone {
				t.Errorf("target = %d, want none", stmt.Target)
			}
			if stmt.Storage != tc.storage {
				t.Errorf("storage = %q, want %q", stmt.Storage, tc.storage)
			}
			if stmt.JobID != tc.job {
				t.Errorf("job = %d, want %d", stmt.JobID, tc.job)
			}
			var names []string
			for _, o := range stmt.Options {
				names = append(names, o.Name)
			}
			if !reflect.DeepEqual(names, tc.options) {
				t.Errorf("options = %v, want %v", names, tc.options)
			}
		})
	}
	// The statements that take or restore a backup stay outside IsLogTaskOrJob: the
	// walk ends at them, and it follows the kinds above.
	for _, sql := range []string{"BACKUP DATABASE * TO 'noop://'", "BACKUP LOGS TO 'noop://'", "RESTORE DATABASE * FROM 'noop://'", "RESTORE POINT FROM 'noop://'"} {
		if parseOneTiDBAdmin(t, sql).(*ast.BRIEStmt).Kind.IsLogTaskOrJob() {
			t.Errorf("%s: IsLogTaskOrJob = true", sql)
		}
	}
}

func TestCancelDistributionJobInTree(t *testing.T) {
	stmt, ok := parseOneTiDBAdmin(t, "CANCEL DISTRIBUTION JOB 42").(*ast.CancelDistributionJobStmt)
	if !ok || stmt.JobID != 42 {
		t.Fatalf("got %+v", stmt)
	}
}

func TestKillInTree(t *testing.T) {
	for _, tc := range []struct {
		sql         string
		tidb, query bool
		id          uint64
		fn          string
		args        int
	}{
		{sql: "KILL 23123", id: 23123},
		{sql: "KILL CONNECTION 23123", id: 23123},
		{sql: "KILL QUERY 23123", query: true, id: 23123},
		{sql: "KILL TIDB 23123", tidb: true, id: 23123},
		{sql: "kill tidb connection 23123", tidb: true, id: 23123},
		{sql: "KILL TIDB QUERY 23123", tidb: true, query: true, id: 23123},
		{sql: "KILL 18446744073709551615", id: 18446744073709551615},
		{sql: "KILL CONNECTION_ID()", fn: "CONNECTION_ID"},
		{sql: "KILL TIDB CONNECTION_ID()", tidb: true, fn: "CONNECTION_ID"},
		{sql: "KILL ((foo(1, 2)))", fn: "foo", args: 2},
		{sql: "KILL REPLACE('a', 'b', 'c')", fn: "REPLACE", args: 3},
		{sql: "KILL `tidb`(1)", fn: "tidb", args: 1},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			stmt, ok := parseOneTiDBAdmin(t, tc.sql).(*ast.KillStmt)
			if !ok {
				t.Fatalf("got %T", parseOneTiDBAdmin(t, tc.sql))
			}
			if stmt.TiDB != tc.tidb || stmt.Query != tc.query || stmt.ID != tc.id {
				t.Errorf("got tidb=%v query=%v id=%d", stmt.TiDB, stmt.Query, stmt.ID)
			}
			if tc.fn == "" {
				if _, ok := stmt.ConnectionID.(*ast.IntLit); !ok {
					t.Errorf("ConnectionID = %T, want an IntLit", stmt.ConnectionID)
				}
				return
			}
			call, ok := stmt.ConnectionID.(*ast.FuncCallExpr)
			if !ok {
				t.Fatalf("ConnectionID = %T, want a call", stmt.ConnectionID)
			}
			if call.Name != tc.fn || len(call.Args) != tc.args || !call.HasParens {
				t.Errorf("call = %s", ast.NodeToString(call))
			}
		})
	}
}

func TestFlushInTree(t *testing.T) {
	for _, tc := range []struct {
		sql     string
		local   bool
		option  string
		plugins []string
		tables  int
		lock    bool
	}{
		{sql: "FLUSH TIDB PLUGINS a", option: "TIDB PLUGINS", plugins: []string{"a"}},
		{sql: "flush tidb plugins a, `b c`, status", option: "TIDB PLUGINS", plugins: []string{"a", "b c", "status"}},
		{sql: "FLUSH LOCAL TIDB PLUGINS a", local: true, option: "TIDB PLUGINS", plugins: []string{"a"}},
		{sql: "FLUSH NO_WRITE_TO_BINLOG TIDB PLUGINS a", local: true, option: "TIDB PLUGINS", plugins: []string{"a"}},
		{sql: "FLUSH PRIVILEGES", option: "PRIVILEGES"},
		{sql: "FLUSH STATUS", option: "STATUS"},
		{sql: "FLUSH HOSTS", option: "HOSTS"},
		{sql: "FLUSH CLIENT_ERRORS_SUMMARY", option: "CLIENT_ERRORS_SUMMARY"},
		{sql: "FLUSH LOGS", option: "LOGS"},
		{sql: "FLUSH BINARY LOGS", option: "BINARY LOGS"},
		{sql: "FLUSH ENGINE LOGS", option: "ENGINE LOGS"},
		{sql: "FLUSH ERROR LOGS", option: "ERROR LOGS"},
		{sql: "FLUSH GENERAL LOGS", option: "GENERAL LOGS"},
		{sql: "FLUSH SLOW LOGS", option: "SLOW LOGS"},
		{sql: "FLUSH TABLES", option: "TABLES"},
		{sql: "FLUSH TABLE t, u WITH READ LOCK", option: "TABLES", tables: 2, lock: true},
		{sql: "FLUSH LOCAL TABLES WITH READ LOCK", local: true, option: "TABLES", lock: true},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			stmt, ok := parseOneTiDBAdmin(t, tc.sql).(*ast.FlushStmt)
			if !ok {
				t.Fatalf("got %T", parseOneTiDBAdmin(t, tc.sql))
			}
			if stmt.NoWriteToBinlog != tc.local || len(stmt.Tables) != tc.tables || stmt.WithReadLock != tc.lock {
				t.Errorf("got local=%v tables=%d lock=%v", stmt.NoWriteToBinlog, len(stmt.Tables), stmt.WithReadLock)
			}
			if len(stmt.Options) != 1 || stmt.Options[0] != tc.option {
				t.Errorf("options = %v, want [%s]", stmt.Options, tc.option)
			}
			if !reflect.DeepEqual(stmt.Plugins, tc.plugins) {
				t.Errorf("plugins = %v, want %v", stmt.Plugins, tc.plugins)
			}
		})
	}
}

func TestAlterRangeInTree(t *testing.T) {
	for _, tc := range []struct {
		sql, name, policy string
	}{
		{"ALTER RANGE global PLACEMENT POLICY mypolicy", "global", "mypolicy"},
		{"alter range `META` placement policy = `p`", "META", "p"},
		{"ALTER RANGE global PLACEMENT POLICY 'p q'", "global", "p q"},
		{"ALTER RANGE global PLACEMENT POLICY = \"p\"", "global", "p"},
		{"ALTER RANGE global PLACEMENT POLICY default", "global", "DEFAULT"},
		{"ALTER RANGE global PLACEMENT POLICY = DEFAULT", "global", "DEFAULT"},
		{"ALTER RANGE global PLACEMENT POLICY SET DEFAULT", "global", "DEFAULT"},
		{"ALTER RANGE global PLACEMENT POLICY `default`", "global", "default"},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			stmt, ok := parseOneTiDBAdmin(t, tc.sql).(*ast.AlterRangeStmt)
			if !ok {
				t.Fatalf("got %T", parseOneTiDBAdmin(t, tc.sql))
			}
			if stmt.Name != tc.name || stmt.Policy != tc.policy {
				t.Errorf("got name=%q policy=%q", stmt.Name, stmt.Policy)
			}
		})
	}
}

// TestNameAfterDotKeepsItsDigits pins what the lexer's identDot leaves in the tree:
// the digits are the name's, not a float's.
func TestNameAfterDotKeepsItsDigits(t *testing.T) {
	sel := parseOneTiDBAdmin(t, "SELECT a.1, a. 1e5 FROM t").(*ast.SelectStmt)
	for i, want := range []string{"1", "1e5"} {
		col := findColumnRef(t, sel.TargetList[i])
		if col.Table != "a" || col.Column != want {
			t.Errorf("column %d = %s.%s, want a.%s", i, col.Table, col.Column, want)
		}
	}
	ref := parseOneTiDBAdmin(t, "DROP TABLE a.1b").(*ast.DropTableStmt).Tables[0]
	if ref.Schema != "a" || ref.Name != "1b" {
		t.Errorf("table = %s.%s, want a.1b", ref.Schema, ref.Name)
	}
	// A space before the dot leaves the number a number: `a .1` is a 1064, and
	// nothing else about numbers moved.
	if _, err := Parse("SELECT a .1 FROM t"); err == nil {
		t.Error("SELECT a .1 FROM t parsed")
	}
	sel = parseOneTiDBAdmin(t, "SELECT .5 FROM t").(*ast.SelectStmt)
	if !strings.Contains(ast.NodeToString(sel.TargetList[0]), ".5") {
		t.Errorf("SELECT .5 reads %s", ast.NodeToString(sel.TargetList[0]))
	}
}

// TestTokenizeLeavesDigitsAfterADot pins that identDot is Parse's alone. Tokenize reads
// text at a caret for completion and limit rewriting, and `a.1` there stays the two
// tokens it always was.
func TestTokenizeLeavesDigitsAfterADot(t *testing.T) {
	toks := Tokenize("a.1")
	if len(toks) < 2 || toks[1].Type != tokFCONST || toks[1].Str != ".1" {
		t.Errorf("Tokenize(a.1) = %+v, want a float after the name", toks)
	}
}

// TestSnapshotRestoresIdentDot pins that a rewind puts identDot back. Every scan in the
// parser that rewinds starts somewhere no digit can follow a dotted word in a statement
// that parses, so no statement can show this, and the state is asserted directly: a word
// that ran into a dot, a scan past it, a rewind, and the digits must read as a name again.
func TestSnapshotRestoresIdentDot(t *testing.T) {
	p := &Parser{lexer: NewLexerWithOffset("a.1 b", 0)}
	p.lexer.dotNamesWord = true
	p.advance() // a
	snap := p.snapshotTokenStream()
	p.advance() // .
	p.advance() // 1
	p.advance() // b, after which the flag is off
	p.restoreTokenStream(snap)
	p.advance() // .
	if p.cur.Type != '.' {
		t.Fatalf("after the rewind the next token is %d %q, want the dot", p.cur.Type, p.cur.Str)
	}
	p.advance() // 1
	if p.cur.Type != tokIDENT || p.cur.Str != "1" {
		t.Errorf("after the rewind the digits read %d %q, want a name", p.cur.Type, p.cur.Str)
	}
}

func findColumnRef(t *testing.T, n ast.ExprNode) *ast.ColumnRef {
	t.Helper()
	switch v := n.(type) {
	case *ast.ColumnRef:
		return v
	case *ast.ResTarget:
		return findColumnRef(t, v.Val)
	}
	t.Fatalf("not a column: %T", n)
	return nil
}

// TestAdminStatementsSurviveDeepInput holds these readers to the contract that no
// input overflows the stack. KILL's parentheses and a config item's parts are read in
// loops; the arguments of a KILL call go through parseExpr, which has its own limit.
func TestAdminStatementsSurviveDeepInput(t *testing.T) {
	const n = 200000
	if _, err := Parse("KILL " + strings.Repeat("(", n) + "foo()" + strings.Repeat(")", n)); err != nil {
		t.Errorf("KILL in %d parentheses: %v", n, err)
	}
	if _, err := Parse("KILL " + strings.Repeat("(", n)); err == nil {
		t.Error("KILL and a run of open parentheses parsed")
	}
	if _, err := Parse("KILL foo(" + strings.Repeat("(", n) + "1" + strings.Repeat(")", n) + ")"); err == nil {
		t.Error("KILL with an argument past the nesting limit parsed")
	} else {
		var pe *ParseError
		if !errors.As(err, &pe) {
			t.Errorf("nesting error is %T", err)
		}
	}
	if _, err := Parse("SET CONFIG tikv a" + strings.Repeat(".a", n) + " = 1"); err != nil {
		t.Errorf("a config item of %d parts: %v", n, err)
	}
	if _, err := Parse("FLUSH TIDB PLUGINS a" + strings.Repeat(", a", n)); err != nil {
		t.Errorf("%d plugins: %v", n, err)
	}
	// Adjacent literals concatenate in linear time: a quadratic join of this run
	// allocates tens of gigabytes.
	sql := "RECOMMEND INDEX RUN WITH a = " + strings.Repeat("'x' ", n)
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("%d adjacent literals: %v", n, err)
	}
	if got := list.Items[0].(*ast.RecommendIndexStmt).Options[0].Value.(*ast.StringLit).Value; len(got) != n {
		t.Errorf("concatenated %d bytes, want %d", len(got), n)
	}
}
