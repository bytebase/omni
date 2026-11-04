package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

// The shape of the SHOW clauses this change rewrote. A target table measured
// only by whether a statement parses cannot see a clause that landed in the
// wrong field: that parses clean and misreports for as long as it stands.

// TestShowChannelShape pins the channel of SHOW {SLAVE | REPLICA} STATUS and
// SHOW RELAYLOG EVENTS. It is a string literal, not an identifier, so the
// value that reaches Channel is the string's contents with no quotes.
func TestShowChannelShape(t *testing.T) {
	for _, tc := range []struct{ sql, typ, channel string }{
		{"SHOW SLAVE STATUS", "SLAVE STATUS", ""},
		{"SHOW SLAVE STATUS FOR CHANNEL 'ch1'", "SLAVE STATUS", "ch1"},
		{"SHOW REPLICA STATUS FOR CHANNEL 'ch1'", "REPLICA STATUS", "ch1"},
		{"SHOW REPLICA STATUS FOR CHANNEL ''", "REPLICA STATUS", ""},
		{"SHOW RELAYLOG EVENTS FOR CHANNEL 'ch1'", "RELAYLOG EVENTS", "ch1"},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			stmt := parseOneShow(t, tc.sql)
			if stmt.Type != tc.typ {
				t.Errorf("Type = %q, want %q", stmt.Type, tc.typ)
			}
			if stmt.Channel != tc.channel {
				t.Errorf("Channel = %q, want %q", stmt.Channel, tc.channel)
			}
		})
	}
}

// TestShowEngineShape pins which word of SHOW ENGINE name ... decides the
// type, and that the engine name lands on From rather than on the type.
func TestShowEngineShape(t *testing.T) {
	for _, tc := range []struct{ sql, typ, name string }{
		{"SHOW ENGINE INNODB STATUS", "ENGINE STATUS", "INNODB"},
		{"SHOW ENGINE INNODB MUTEX", "ENGINE MUTEX", "INNODB"},
		{"SHOW ENGINE INNODB LOGS", "ENGINE LOGS", "INNODB"},
		{"SHOW ENGINE PERFORMANCE_SCHEMA STATUS", "ENGINE STATUS", "PERFORMANCE_SCHEMA"},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			stmt := parseOneShow(t, tc.sql)
			if stmt.Type != tc.typ {
				t.Errorf("Type = %q, want %q", stmt.Type, tc.typ)
			}
			if stmt.From == nil || stmt.From.Name != tc.name {
				t.Errorf("From = %+v, want name %q", stmt.From, tc.name)
			}
		})
	}
}

// TestShowProfileLimitOffsetShape pins which of the two LIMIT numbers is the
// offset. SHOW PROFILE reads a whole limit clause, where `LIMIT a, b` means
// offset a and count b while `LIMIT a OFFSET b` means count a and offset b;
// reading them the same way round would have swapped the pair silently.
func TestShowProfileLimitOffsetShape(t *testing.T) {
	stmt := parseOneShow(t, "SHOW PROFILE CPU FOR QUERY 2 LIMIT 1,3")
	if stmt.Type != "PROFILE" {
		t.Fatalf("Type = %q, want PROFILE", stmt.Type)
	}
	if got := showIntLit(t, stmt.LimitOffset); got != 1 {
		t.Errorf("LimitOffset = %d, want 1", got)
	}
	if got := showIntLit(t, stmt.LimitCount); got != 3 {
		t.Errorf("LimitCount = %d, want 3", got)
	}

	stmt = parseOneShow(t, "SHOW PROFILE CPU LIMIT 1 OFFSET 3")
	if got := showIntLit(t, stmt.LimitCount); got != 1 {
		t.Errorf("LimitCount = %d, want 1", got)
	}
	if got := showIntLit(t, stmt.LimitOffset); got != 3 {
		t.Errorf("LimitOffset = %d, want 3", got)
	}
}

// TestShowLikeShape pins what the LIKE of a filter holds. It is one string
// literal, so Like is a StringLit carrying the pattern and Where stays nil.
func TestShowLikeShape(t *testing.T) {
	for _, tc := range []struct{ sql, pattern string }{
		{"SHOW DATABASES LIKE 't'", "t"},
		{"SHOW DATABASES LIKE 'a%'", "a%"},
		{"SHOW DATABASES LIKE ''", ""},
		{`SHOW TABLES LIKE "t"`, "t"},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			stmt := parseOneShow(t, tc.sql)
			lit, ok := stmt.Like.(*nodes.StringLit)
			if !ok {
				t.Fatalf("Like = %T, want *ast.StringLit", stmt.Like)
			}
			if lit.Value != tc.pattern {
				t.Errorf("Like = %q, want %q", lit.Value, tc.pattern)
			}
			if stmt.Where != nil {
				t.Errorf("Where = %v, want nil", stmt.Where)
			}
		})
	}
	// WHERE still reads a whole expression, and it lands on Where.
	stmt := parseOneShow(t, "SHOW DATABASES WHERE 1=1")
	if stmt.Where == nil {
		t.Error("Where is nil, want the comparison")
	}
	if stmt.Like != nil {
		t.Errorf("Like = %v, want nil", stmt.Like)
	}
}

// parseOneShow parses sql through the whole entry point and returns the one
// SHOW statement it holds. Going through Parse rather than parseShowStmt is
// deliberate: a SHOW that leaves tokens behind is a second statement in the
// list, which is the shape the empty-target bug took, and asserting on one
// item catches it.
func parseOneShow(t *testing.T, sql string) *nodes.ShowStmt {
	t.Helper()
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	if len(list.Items) != 1 {
		t.Fatalf("Parse(%q) returned %d statements, want 1", sql, len(list.Items))
	}
	stmt, ok := list.Items[0].(*nodes.ShowStmt)
	if !ok {
		t.Fatalf("Parse(%q) returned %T, want *ast.ShowStmt", sql, list.Items[0])
	}
	return stmt
}

func showIntLit(t *testing.T, e nodes.ExprNode) int64 {
	t.Helper()
	lit, ok := e.(*nodes.IntLit)
	if !ok {
		t.Fatalf("expression is %T, want *ast.IntLit", e)
	}
	return lit.Value
}

// TestShowMariaDBOnlyShapes pins the four SHOW forms MariaDB has and MySQL
// 8.0 does not, each measured against mariadb:11.8.8 in showTargetCases.
func TestShowMariaDBOnlyShapes(t *testing.T) {
	// SHOW CREATE SEQUENCE: MariaDB has SEQUENCE objects (10.3+), so the
	// sequence lands on From like any other named object.
	//
	// Ref: https://mariadb.com/kb/en/show-create-sequence/
	t.Run("SHOW CREATE SEQUENCE seq", func(t *testing.T) {
		stmt := parseOneShow(t, "SHOW CREATE SEQUENCE seq")
		if stmt.Type != "CREATE SEQUENCE" {
			t.Errorf("Type = %q, want CREATE SEQUENCE", stmt.Type)
		}
		if stmt.From == nil || stmt.From.Name != "seq" || stmt.From.Schema != "" {
			t.Errorf("From = %+v, want name seq", stmt.From)
		}
	})
	t.Run("SHOW CREATE SEQUENCE test.seq", func(t *testing.T) {
		stmt := parseOneShow(t, "SHOW CREATE SEQUENCE test.seq")
		if stmt.From == nil || stmt.From.Schema != "test" || stmt.From.Name != "seq" {
			t.Errorf("From = %+v, want test.seq", stmt.From)
		}
	})

	// SHOW BINLOG STATUS, MariaDB 10.5.2's name for SHOW MASTER STATUS. It
	// shares the BINLOG keyword with SHOW BINLOG EVENTS and must not be read
	// as one.
	//
	// Ref: https://mariadb.com/kb/en/show-binlog-status/
	t.Run("SHOW BINLOG STATUS", func(t *testing.T) {
		stmt := parseOneShow(t, "SHOW BINLOG STATUS")
		if stmt.Type != "BINLOG STATUS" {
			t.Errorf("Type = %q, want BINLOG STATUS", stmt.Type)
		}
		if stmt.Like != nil || stmt.Where != nil || stmt.LimitCount != nil {
			t.Errorf("filter or limit set on %+v, want none", stmt)
		}
	})

	// SHOW REPLICA HOSTS, MariaDB 10.5.1's alias for SHOW SLAVE HOSTS. Both
	// read as the same statement; MySQL's SHOW REPLICAS is a 1064 here.
	//
	// Ref: https://mariadb.com/kb/en/show-replica-hosts/
	for _, sql := range []string{"SHOW REPLICA HOSTS", "SHOW SLAVE HOSTS"} {
		t.Run(sql, func(t *testing.T) {
			if got := parseOneShow(t, sql).Type; got != "REPLICAS" {
				t.Errorf("Type = %q, want REPLICAS", got)
			}
		})
	}

	// SHOW CREATE USER takes no user on MariaDB, which then reports the
	// current one, so ForUser stays nil rather than being filled in with a
	// name the text does not carry.
	//
	// Ref: https://mariadb.com/kb/en/show-create-user/
	t.Run("SHOW CREATE USER", func(t *testing.T) {
		stmt := parseOneShow(t, "SHOW CREATE USER")
		if stmt.Type != "CREATE USER" {
			t.Errorf("Type = %q, want CREATE USER", stmt.Type)
		}
		if stmt.ForUser != nil {
			t.Errorf("ForUser = %+v, want nil", stmt.ForUser)
		}
	})
	t.Run("SHOW CREATE USER u", func(t *testing.T) {
		stmt := parseOneShow(t, "SHOW CREATE USER u")
		if stmt.ForUser == nil || stmt.ForUser.Name != "u" {
			t.Errorf("ForUser = %+v, want name u", stmt.ForUser)
		}
	})
}
