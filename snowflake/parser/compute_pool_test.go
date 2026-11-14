package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

func mustObjectStmt(t *testing.T, input string) *ast.ObjectStmt {
	t.Helper()
	node := mustParseOne(t, input)
	stmt, ok := node.(*ast.ObjectStmt)
	if !ok {
		t.Fatalf("parse %q: got %T, want *ast.ObjectStmt", input, node)
	}
	if stmt.Loc.Start != 0 || stmt.Loc.End != len(input) {
		t.Errorf("parse %q: Loc = %+v, want {0 %d}", input, stmt.Loc, len(input))
	}
	for _, c := range stmt.Clauses {
		if c.Loc.Start < 0 || c.Loc.End > len(input) || c.Loc.Start >= c.Loc.End {
			t.Errorf("parse %q: clause %q Loc = %+v", input, c.Keywords, c.Loc)
		}
	}
	return stmt
}

func mustRejectObject(t *testing.T, input string) {
	t.Helper()
	if _, err := Parse(input); err == nil {
		t.Errorf("Parse(%q) succeeded, want an error", input)
	}
}

// clauseKeywords lists the Keywords of an ObjectStmt's clauses.
func clauseKeywords(s *ast.ObjectStmt) []string {
	var out []string
	for _, c := range s.Clauses {
		out = append(out, c.Keywords)
	}
	return out
}

func TestParseCreateComputePool(t *testing.T) {
	in := "CREATE COMPUTE POOL IF NOT EXISTS pool1 FOR APPLICATION app MIN_NODES = 1 MAX_NODES = 2 " +
		"INSTANCE_FAMILY = CPU_X64_XS AUTO_RESUME = TRUE WITH TAG (t = 'v') COMMENT = 'c' " +
		"BACKUP_INSTANCE_FAMILIES = ('a', 'b')"
	s := mustObjectStmt(t, in)
	if s.Verb != "CREATE" || s.Kind != "COMPUTE POOL" || !s.IfNotExists || s.Name.String() != "pool1" {
		t.Fatalf("got %+v", s)
	}
	if got := clauseKeywords(s); len(got) != 2 || got[0] != "FOR APPLICATION" || got[1] != "" {
		t.Fatalf("clauses = %q", got)
	}
	props := s.Clauses[1]
	if len(props.Options) != 6 || len(props.Tags) != 1 {
		t.Fatalf("options = %d, tags = %d", len(props.Options), len(props.Tags))
	}
	if w, ok := optWords(props.Options, "INSTANCE_FAMILY"); !ok || w != "CPU_X64_XS" {
		t.Errorf("INSTANCE_FAMILY = %q", w)
	}
}

func TestParseAlterComputePool(t *testing.T) {
	for in, want := range map[string]string{
		"ALTER COMPUTE POOL IF EXISTS p SUSPEND":                      "SUSPEND",
		"ALTER COMPUTE POOL p RESUME":                                 "RESUME",
		"ALTER COMPUTE POOL p STOP ALL":                               "STOP ALL",
		"ALTER COMPUTE POOL p STOP ALL OF TYPE a, b":                  "STOP ALL OF TYPE",
		"ALTER COMPUTE POOL p SET MIN_NODES = 1 TAG t = 's', u = 's'": "SET",
		"ALTER COMPUTE POOL p SET BACKUP_INSTANCE_FAMILIES = ('s')":   "SET",
		"ALTER COMPUTE POOL p UNSET AUTO_SUSPEND_SECS, COMMENT":       "UNSET",
	} {
		s := mustObjectStmt(t, in)
		if s.Verb != "ALTER" || len(s.Clauses) != 1 || s.Clauses[0].Keywords != want {
			t.Errorf("%q: clauses = %q, want [%q]", in, clauseKeywords(s), want)
		}
	}
	s := mustObjectStmt(t, "ALTER COMPUTE POOL p SET MIN_NODES = 1 TAG t = 's', u = 's'")
	if c := s.Clauses[0]; len(c.Options) != 1 || len(c.Tags) != 2 {
		t.Errorf("SET clause = %+v", c)
	}
}

func TestParseDropComputePool(t *testing.T) {
	node := mustParseOne(t, "DROP COMPUTE POOL IF EXISTS p")
	d, ok := node.(*ast.DropStmt)
	if !ok || d.Kind != ast.DropObject || d.Object != "COMPUTE POOL" || !d.IfExists || d.Name.String() != "p" {
		t.Fatalf("got %#v", node)
	}
}

func TestParseComputePoolRejects(t *testing.T) {
	for _, in := range []string{
		"CREATE COMPUTE POOL p",                                                            // no properties
		"CREATE COMPUTE POOL p MIN_NODES = 1 MAX_NODES = 1",                                // INSTANCE_FAMILY is required
		"CREATE OR REPLACE COMPUTE POOL p MIN_NODES = 1 MAX_NODES = 1 INSTANCE_FAMILY = a", // not documented
		"CREATE COMPUTE POOL p MIN_NODES = 1 MAX_NODES = 1 INSTANCE_FAMILY = a SUSPEND",
		"ALTER COMPUTE POOL p",
		"ALTER COMPUTE POOL p SET",
		"ALTER COMPUTE POOL p UNSET",
		"ALTER COMPUTE POOL p RENAME TO q", // not documented
		"ALTER COMPUTE POOL p STOP",
		"ALTER COMPUTE POOL p SUSPEND RESUME",
		"DROP COMPUTE POOL p CASCADE",
		"ALTER 'COMPUTE' 'POOL' p 'SUSPEND'", // a string never stands for a keyword
		"ALTER COMPUTE \"POOL\" p SUSPEND",
		"ALTER COMPUTE POOL p STOP \"ALL\"",
		"ALTER COMPUTE POOL p UNSET 123",
		"CREATE COMPUTE POOL p MIN_NODES = () MAX_NODES = (foo) INSTANCE_FAMILY = CPU_X64_XS",
		"ALTER COMPUTE POOL p SET AUTO_RESUME = ()",
		"ALTER COMPUTE POOL p SET AUTO_RESUME = MAYBE",
		"ALTER COMPUTE POOL p SET MIN_NODES = 'one'",
		"ALTER COMPUTE POOL p SET NOT_A_PROPERTY = 1",
		"ALTER COMPUTE POOL p SET BACKUP_INSTANCE_FAMILIES = (a)",
		"ALTER COMPUTE POOL p UNSET MIN_NODES",
		"ALTER COMPUTE POOL p SET INITIALLY_SUSPENDED = TRUE", // CREATE's alone
		"ALTER COMPUTE POOL p SET WITH TAG (t = 'x')",
		"ALTER COMPUTE POOL db.p SUSPEND", // an account-level object
		"CREATE COMPUTE POOL db.p MIN_NODES = 1 MAX_NODES = 1 INSTANCE_FAMILY = a",
		"DROP COMPUTE POOL db.p",
		"CREATE COMPUTE POOL p FOR APPLICATION db.a MIN_NODES = 1 MAX_NODES = 1 INSTANCE_FAMILY = a",
		"ALTER COMPUTE POOL p STOP ALL OF TYPE db.a",
		"ALTER COMPUTE POOL p SET AUTO_RESUME = TRUE AUTO_RESUME = FALSE",
		"CREATE COMPUTE POOL p MIN_NODES = 1 MAX_NODES = 1 INSTANCE_FAMILY = a MIN_NODES = 2",
		"ALTER COMPUTE POOL p SET TAG (t = 'x')",
		"ALTER COMPUTE POOL p SET INSTANCE_FAMILY = foo.",
		"ALTER COMPUTE POOL p SET INSTANCE_FAMILY = foo..bar",
		"ALTER COMPUTE POOL p SET AUTO_RESUME = \"TRUE\"",
		"ALTER COMPUTE POOL p SET AUTO_RESUME = \"true\"",
		"CREATE COMPUTE POOL p \"MIN_NODES\" = 1 MAX_NODES = 1 INSTANCE_FAMILY = a",
		"ALTER COMPUTE POOL p SET C = (\"K\" = 1)",
		"ALTER COMPUTE POOL p UNSET COMMENT, 1.5",
		"ALTER COMPUTE POOL p UNSET 'COMMENT'",
	} {
		mustRejectObject(t, in)
	}
}

func TestAtWords(t *testing.T) {
	p := &Parser{lexer: NewLexer("COMPUTE POOL x y"), input: "COMPUTE POOL x y"}
	p.advance()
	if !p.atWords("COMPUTE", "POOL", "X") || p.atWords("COMPUTE", "POOL", "Y") || !p.atWords("COMPUTE") {
		t.Fatal("atWords")
	}
	if p.cur.Loc.Start != 0 {
		t.Fatal("atWords consumed tokens")
	}
	if !p.acceptWords("COMPUTE", "POOL") || !p.curIsWord("X") {
		t.Fatal("acceptWords")
	}
}

// readObjectProperties reads in as a SET's property run, with no object's
// documented properties to hold it to.
func readObjectProperties(in string) (*ast.ObjectClause, error) {
	p := &Parser{lexer: NewLexer(in), input: in}
	p.advance()
	c, err := p.parseObjectProperties(&ast.ObjectStmt{}, "SET", true)
	if err == nil && p.cur.Type != tokEOF {
		err = p.syntaxErrorAtCur()
	}
	return c, err
}

func TestParseObjectValue(t *testing.T) {
	c, err := readObjectProperties("A = ('x', 'y') B = (c, d) C = (K = 's', L = (M = 1)) D = ((E = 1), (E = 2)) F = () G = -1")
	if err != nil {
		t.Fatal(err)
	}
	opts := c.Options
	if len(opts) != 6 || len(opts[0].List) != 2 || len(opts[1].Group) != 2 || !opts[1].Group[0].Bare ||
		len(opts[2].Group[1].Group) != 1 || len(opts[3].Group) != 2 || opts[4].Group == nil || opts[5].Lit.Value != "-1" {
		t.Fatalf("options = %+v", opts)
	}
	for _, in := range []string{
		"A = ('x' 'y')",
		"A = ('x' z)",
		"A = (K = )",
		"A = (",
		"A =",
		"C = (\"K\" = 1)",
		"\"A\" = 1",
	} {
		if _, err := readObjectProperties(in); err == nil {
			t.Errorf("%q: want an error", in)
		}
	}
}

func TestObjectValueMatches(t *testing.T) {
	c, err := readObjectProperties("N = 1 S = 's' W = x B = TRUE L = ('a') M = (a, b) G = (K = 1)")
	if err != nil {
		t.Fatal(err)
	}
	want := []objectValue{propNumber, propString, propName, propBool, propStrings, propNames, propGroup}
	for i, o := range c.Options {
		for j, v := range want {
			// propName takes any word, so TRUE matches it too.
			if got := v.matches(o); got != (i == j || (i == 3 && j == 2)) {
				t.Errorf("%s matches %d = %v", o.Name, j, got)
			}
		}
	}
}

func TestParseObjectAsQuery(t *testing.T) {
	p := &Parser{lexer: NewLexer("AS SELECT a FROM t"), input: "AS SELECT a FROM t"}
	p.advance()
	s := &ast.ObjectStmt{}
	if err := p.parseObjectAsQuery(s); err != nil {
		t.Fatal(err)
	}
	if c := s.Clauses[0]; c.Keywords != "AS" || c.Query == nil || c.Loc.End != len("AS SELECT a FROM t") {
		t.Fatalf("clause = %+v", c)
	}
}

func TestObjectFormIfExistsFlags(t *testing.T) {
	f := &objectForm{words: []string{"ZZ", "THING"}, noIfNotExists: true, noIfExists: true, dropNoIfExists: true, drop: true,
		create: func(*Parser, *ast.ObjectStmt) error { return nil },
		alter:  func(*Parser, *ast.ObjectStmt) error { return nil }}
	registerObject(f)
	defer delete(objectForms, "ZZ")
	for _, in := range []string{"CREATE ZZ THING t", "ALTER ZZ THING t", "DROP ZZ THING t"} {
		mustParseOne(t, in)
	}
	for _, in := range []string{"CREATE ZZ THING IF NOT EXISTS t", "ALTER ZZ THING IF EXISTS t", "DROP ZZ THING IF EXISTS t"} {
		mustRejectObject(t, in)
	}
}
