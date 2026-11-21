package parser

import (
	"testing"

	"github.com/bytebase/omni/mongo/ast"
)

func TestDirectShellCommands(t *testing.T) {
	t.Run("use bare and hyphenated", func(t *testing.T) {
		for _, tc := range []struct {
			input  string
			wantDB string
		}{
			{"use mydb", "mydb"},
			{"use admin", "admin"},
			{"use my-app-db", "my-app-db"},
			{"use db-1", "db-1"},
			{`use "quoted-db"`, "quoted-db"},
		} {
			nodes, err := Parse(tc.input)
			if err != nil || len(nodes) != 1 {
				t.Fatalf("%q: Parse returned %d nodes, err=%v", tc.input, len(nodes), err)
			}
			call, ok := nodes[0].(*ast.NativeFunctionCall)
			if !ok || call.Name != "use" || len(call.Args) != 1 {
				t.Fatalf("%q: got %#v, want NativeFunctionCall(use, 1 arg)", tc.input, nodes[0])
			}
			lit, ok := call.Args[0].(*ast.StringLiteral)
			if !ok || lit.Value != tc.wantDB {
				t.Errorf("%q: arg = %#v, want %q", tc.input, call.Args[0], tc.wantDB)
			}
			if call.Loc.Start != 0 || call.Loc.End != len(tc.input) {
				t.Errorf("%q: Loc = %+v, want [0,%d]", tc.input, call.Loc, len(tc.input))
			}
		}
	})

	t.Run("use errors", func(t *testing.T) {
		for _, bad := range []string{"use", "use;", "use\nshow dbs", "use my-"} {
			if _, err := Parse(bad); err == nil {
				t.Errorf("%q: expected parse error, got nil", bad)
			}
		}
	})

	t.Run("exit quit it cls", func(t *testing.T) {
		for _, tc := range []struct {
			input    string
			wantName string
			wantArgs int
		}{
			{"exit", "exit", 0},
			{"exit 0", "exit", 1},
			{"exit 1", "exit", 1},
			{"quit", "quit", 0},
			{"quit 2", "quit", 1},
			{"it", "it", 0},
			{"cls", "cls", 0},
			{"exit()", "exit", 0},
			{"quit(0)", "quit", 1},
			{"it()", "it", 0},
			{"cls()", "cls", 0},
			{`show("dbs")`, "show", 1},
			{"show()", "show", 0},
		} {
			nodes, err := Parse(tc.input)
			if err != nil || len(nodes) != 1 {
				t.Fatalf("%q: Parse returned %d nodes, err=%v", tc.input, len(nodes), err)
			}
			call, ok := nodes[0].(*ast.NativeFunctionCall)
			if !ok || call.Name != tc.wantName || len(call.Args) != tc.wantArgs {
				t.Errorf("%q: got %#v, want NativeFunctionCall(%s, %d args)", tc.input, nodes[0], tc.wantName, tc.wantArgs)
			}
		}
	})

	t.Run("newline terminates optional arg", func(t *testing.T) {
		nodes, err := Parse("show log\nshow dbs")
		if err != nil || len(nodes) != 2 {
			t.Fatalf("show log\\nshow dbs: got %d nodes, err=%v", len(nodes), err)
		}
		first := nodes[0].(*ast.ShowCommand)
		second := nodes[1].(*ast.ShowCommand)
		if first.Target != "log" || first.Arg != "" || second.Target != "databases" {
			t.Errorf("got (%+v, %+v)", first, second)
		}

		nodes, err = Parse("exit\nshow dbs")
		if err != nil || len(nodes) != 2 {
			t.Fatalf("exit\\nshow dbs: got %d nodes, err=%v", len(nodes), err)
		}
	})
}

func TestShowExtendedTargets(t *testing.T) {
	for _, tc := range []struct {
		input      string
		wantTarget string
		wantArg    string
	}{
		{"show log global", "log", "global"},
		{"show log startupWarnings", "log", "startupWarnings"},
		{"show automationNotices", "automationNotices", ""},
		{"show nonGenuineMongoDBCheck", "nonGenuineMongoDBCheck", ""},
	} {
		nodes, err := Parse(tc.input)
		if err != nil || len(nodes) != 1 {
			t.Fatalf("%q: Parse returned %d nodes, err=%v", tc.input, len(nodes), err)
		}
		cmd, ok := nodes[0].(*ast.ShowCommand)
		if !ok || cmd.Target != tc.wantTarget || cmd.Arg != tc.wantArg {
			t.Errorf("%q: got %+v, want Target=%q Arg=%q", tc.input, cmd, tc.wantTarget, tc.wantArg)
		}
	}
}

func TestShellConfigAndEJSONCalls(t *testing.T) {
	for _, tc := range []struct {
		input    string
		wantName string
		wantArgs int
	}{
		{`config.get("displayBatchSize")`, "config.get", 1},
		{`config.set("displayBatchSize", 50)`, "config.set", 2},
		{`config.reset("displayBatchSize")`, "config.reset", 1},
		{`EJSON.parse('{"a": 1}')`, "EJSON.parse", 1},
		{`EJSON.stringify({a: 1})`, "EJSON.stringify", 1},
		{`EJSON.serialize({a: 1})`, "EJSON.serialize", 1},
		{`EJSON.deserialize({a: 1})`, "EJSON.deserialize", 1},
	} {
		nodes, err := Parse(tc.input)
		if err != nil || len(nodes) != 1 {
			t.Fatalf("%q: Parse returned %d nodes, err=%v", tc.input, len(nodes), err)
		}
		call, ok := nodes[0].(*ast.NativeFunctionCall)
		if !ok || call.Name != tc.wantName || len(call.Args) != tc.wantArgs {
			t.Errorf("%q: got %#v, want NativeFunctionCall(%s, %d args)", tc.input, nodes[0], tc.wantName, tc.wantArgs)
		}
	}

	for _, bad := range []string{"config", "config.", "config.get", "EJSON", "EJSON.", "EJSON.parse"} {
		if _, err := Parse(bad); err == nil {
			t.Errorf("%q: expected parse error, got nil", bad)
		}
	}
}

func TestDatabaseChainingAndStandaloneGetCollection(t *testing.T) {
	t.Run("standalone getCollection", func(t *testing.T) {
		for _, input := range []string{`db.getCollection("users")`, `db.getCollection()`} {
			nodes, err := Parse(input)
			if err != nil || len(nodes) != 1 {
				t.Fatalf("%q: Parse returned %d nodes, err=%v", input, len(nodes), err)
			}
			stmt, ok := nodes[0].(*ast.DatabaseStatement)
			if !ok || stmt.Method != "getCollection" {
				t.Errorf("%q: got %#v, want DatabaseStatement(getCollection)", input, nodes[0])
			}
		}
	})

	t.Run("db method chains", func(t *testing.T) {
		for _, tc := range []struct {
			input      string
			wantMethod string
			wantChain  int
		}{
			{`db.aggregate([{$currentOp: {}}]).toArray()`, "aggregate", 1},
			{`db.watch().getResumeToken()`, "watch", 1},
			{`db.getSiblingDB("other").getMongo()`, "getSiblingDB", 1},
			{`db.getSiblingDB("other").users.find({active: true}).limit(5)`, "getSiblingDB", 3},
		} {
			nodes, err := Parse(tc.input)
			if err != nil || len(nodes) != 1 {
				t.Fatalf("%q: Parse returned %d nodes, err=%v", tc.input, len(nodes), err)
			}
			stmt, ok := nodes[0].(*ast.DatabaseStatement)
			if !ok || stmt.Method != tc.wantMethod || len(stmt.ChainedMethods) != tc.wantChain {
				t.Errorf("%q: got %#v, want DatabaseStatement(%s, %d chained)", tc.input, nodes[0], tc.wantMethod, tc.wantChain)
			}
		}
	})

	t.Run("invalid db chains still fail", func(t *testing.T) {
		for _, bad := range []string{
			`db.getSiblingDB("other").users`,
			`db.getSiblingDB("other").users.`,
			`db.aggregate().`,
			`db.getCollection(123).find()`,
		} {
			if _, err := Parse(bad); err == nil {
				t.Errorf("%q: expected parse error, got nil", bad)
			}
		}
	})
}

func TestBSONKeywordStaticHelpers(t *testing.T) {
	for _, tc := range []struct {
		input    string
		wantName string
	}{
		{`db.c.find({_id: ObjectId.createFromHexString("507f1f77bcf86cd799439011")})`, "ObjectId.createFromHexString"},
		{`db.c.find({_id: ObjectId.createFromTime(1700000000)})`, "ObjectId.createFromTime"},
		{`db.c.find({ok: ObjectId.isValid("507f1f77bcf86cd799439011")})`, "ObjectId.isValid"},
		{`db.c.find({ts: Timestamp.fromInt(1)})`, "Timestamp.fromInt"},
		{`db.c.find({ts: Timestamp.fromString("1")})`, "Timestamp.fromString"},
		{`db.c.find({now: Date.now()})`, "Date.now"},
	} {
		nodes, err := Parse(tc.input)
		if err != nil || len(nodes) != 1 {
			t.Fatalf("%q: Parse returned %d nodes, err=%v", tc.input, len(nodes), err)
		}
		doc := nodes[0].(*ast.CollectionStatement).Args[0].(*ast.Document)
		call, ok := doc.Pairs[0].Value.(*ast.HelperCall)
		if !ok || call.Name != tc.wantName {
			t.Errorf("%q: value = %#v, want HelperCall(%s)", tc.input, doc.Pairs[0].Value, tc.wantName)
		}
	}
}
