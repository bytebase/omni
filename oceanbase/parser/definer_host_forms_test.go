package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/oceanbase/internal/testcontainer"
)

// A definer's host written without quotes, as OceanBase 4.3.5 reads it, every verdict
// measured by running the statement on oceanbase-ce:4.3.5.6 on 2026-11-13.
//
// SHOW CREATE TRIGGER prints DEFINER = root@%, and the server takes that back before a
// stored program and not before a view. A function, procedure, trigger or event is read by
// the PL parser, whose lexer has `(@%)` among the spellings of a host; a view is read by
// the SQL parser, whose lexer has not. The % is the whole host and follows its @ at once.
//
// Ref: user_variable in OceanBase's src/pl/parser/pl_parser_mysql_mode.l (v4.3.5_CE_BP6,
// line 55) and src/sql/parser/sql_parser_mysql_mode.l (line 68), and opt_sp_definer in
// pl_parser_mysql_mode.y (line 1340).
var definerHostCases = inDatabase(definerHostDatabase, []parserFormCase{
	// Before each kind of stored program.
	{"CREATE DEFINER = root@% TRIGGER DB.h01 BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1", true},
	{"CREATE DEFINER = root@% PROCEDURE h02() SELECT 1", true},
	{"CREATE DEFINER = root@% FUNCTION h03() RETURNS INT DETERMINISTIC RETURN 1", true},
	{"CREATE DEFINER = root@% EVENT h04 ON SCHEDULE EVERY 1 DAY DISABLE DO SELECT 1", true},
	{"ALTER DEFINER = root@% EVENT h04 DISABLE", true},
	// The user in each of its spellings, and the quoted hosts the parser read before.
	{"CREATE DEFINER = `root`@% TRIGGER h05 BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1", true},
	{"CREATE DEFINER = 'root'@% TRIGGER h06 BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1", true},
	{"CREATE DEFINER = root @% TRIGGER h07 BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1", true},
	{"CREATE DEFINER = root@'%' TRIGGER h08 BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1", true},
	{"CREATE DEFINER = root@`%` TRIGGER h09 BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1", true},
	{"CREATE DEFINER = 'root'@'%' TRIGGER h10 BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1", true},
	{"CREATE DEFINER = root@localhost TRIGGER h11 BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1", true},
	// A view takes the quoted host.
	{"CREATE ALGORITHM = UNDEFINED DEFINER = root@'%' SQL SECURITY DEFINER VIEW h12 AS SELECT 1", true},
	{"ALTER ALGORITHM = UNDEFINED DEFINER = root@'%' VIEW h12 AS SELECT 2", true},

	// The % alone and at once.
	{"CREATE DEFINER = root@ % TRIGGER h13 BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1", false},
	{"CREATE DEFINER = root @ % TRIGGER h14 BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1", false},
	{"CREATE DEFINER = root@%example TRIGGER h15 BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1", false},
	{"CREATE DEFINER = root@%.example.com TRIGGER h16 BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1", false},
	{"CREATE DEFINER = root@%% TRIGGER h17 BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1", false},
	{"CREATE DEFINER = root@%@% TRIGGER h18 BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1", false},
	{"CREATE DEFINER = root@ % PROCEDURE h19() SELECT 1", false},
	// CURRENT_USER takes no host.
	{"CREATE DEFINER = CURRENT_USER@% TRIGGER h20 BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1", false},
	// A view does not take it.
	{"CREATE ALGORITHM = UNDEFINED DEFINER = root@% SQL SECURITY DEFINER VIEW h21 AS SELECT 1", false},
	{"CREATE ALGORITHM = UNDEFINED DEFINER = root@% VIEW h22 AS SELECT 1", false},
	{"ALTER ALGORITHM = UNDEFINED DEFINER = root@% VIEW h12 AS SELECT 2", false},
})

// definerHostDatabase is the database TestDefinerHostFormsOceanBaseOracle runs in.
var definerHostDatabase = testcontainer.DatabaseName("definer_host_forms")

// TestDefinerHostForms holds omni to the recorded verdicts without a server.
func TestDefinerHostForms(t *testing.T) {
	for _, c := range definerHostCases {
		t.Run(c.sql, func(t *testing.T) {
			got, err := parsesAsOneStatement(c.sql)
			if got != c.accept {
				t.Errorf("omni accepts=%v, want %v for %q: %v", got, c.accept, c.sql, err)
			}
		})
	}
}

// TestDefinerHostTree holds the definer an unquoted host leaves in the tree.
func TestDefinerHostTree(t *testing.T) {
	for sql, want := range map[string]string{
		"CREATE DEFINER = root@% TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET @a = 1":         "root@%",
		"CREATE DEFINER = `root`@% TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET @a = 1":       "root@%",
		"CREATE DEFINER = root@'%' TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET @a = 1":       "root@'%'",
		"CREATE DEFINER = root@`%` TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET @a = 1":       "root@%",
		"CREATE DEFINER = root@localhost TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET @a = 1": "root@localhost",
	} {
		t.Run(sql, func(t *testing.T) {
			list, err := Parse(sql)
			if err != nil {
				t.Fatalf("Parse(%q): %v", sql, err)
			}
			if got := list.Items[0].(*nodes.CreateTriggerStmt).Definer; got != want {
				t.Errorf("definer is %q, want %q", got, want)
			}
		})
	}
}
