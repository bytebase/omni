package main

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytebase/omni/enginetest/grammar/docs"
)

// TestRewriteReadsTheBNF pins how a quoted BNF line is read into the
// notation docs converts: a defined name stays a fragment in lower case, a
// name no page defines is a <value>, a quoted value is a STRING, a keyword
// the docs write in mixed case is a keyword, and @ stands apart.
func TestRewriteReadsTheBNF(t *testing.T) {
	defs := map[string]bool{"quota_clause": true, "xmlindex_clause": true}
	for _, c := range []struct{ in, want string }{
		{"ALTER USER user [ quota_clause ]", "ALTER USER <user> [ quota_clause ]"},
		{"[ schema. ] table @dblink", "[ <schema>. ] <table> @ <dblink>"},
		{"AS 'certificate_DN' | AS '[ directory_DN ]'", "AS '<certificate_DN>' | AS '<directory_DN>'"},
		{"OF XMLType [ XMLIndex_clause ]", "OF XMLTYPE [ xmlindex_clause ]"},
		{"{ IS | AS } plsql_function_source", "{ IS | AS } <plsql_function_source>"},
		{"ON subquery", "ON <subquery>"},
	} {
		leftOut = map[string]bool{}
		if got := rewrite(c.in, defs); got != c.want {
			t.Errorf("%q\n got %q\nwant %q", c.in, got, c.want)
		}
	}
}

// TestUndefinedClausesAreLeftOut pins where a clause no page defines goes:
// its alternative where its group has others, its [ ] with a ... after it,
// a level out from a { } of one alternative, and nowhere optional, a name
// the converter refuses.
func TestUndefinedClausesAreLeftOut(t *testing.T) {
	for _, c := range []struct{ in, want string }{
		{"A [ storage_clause ] B", "A  B"},
		{"{ PCTFREE n | storage_clause | INITRANS n }", "{ PCTFREE <n> | INITRANS <n> }"},
		{"{ PCTFREE n | storage_clause }", "{ PCTFREE <n> }"},
		{"[ { storage_clause } ]... B", "B"},
		{"( TABLESPACE t | storage_clause )", "( TABLESPACE <t> )"},
		{"USING storage_clause", "USING <storageClause>"},
	} {
		leftOut = map[string]bool{}
		if got := rewrite(c.in, map[string]bool{}); got != c.want && !(c.want == "B" && got == " B") {
			t.Errorf("%q\n got %q\nwant %q", c.in, got, c.want)
		}
		if !leftOut["storage_clause"] && c.want != "USING <storageClause>" {
			t.Errorf("%q: storage_clause not noted as left out", c.in)
		}
	}
}

func TestPlaceholderClasses(t *testing.T) {
	for name, want := range map[string]string{
		"table": docs.ClassName, "schema": docs.ClassName, "c_alias": docs.ClassName, "password": docs.ClassName,
		"integer": docs.ClassNumber, "scn_value": docs.ClassNumber, "n": docs.ClassNumber,
		"expr": docs.ClassExpr, "default_expression": docs.ClassExpr, "literal": docs.ClassExpr,
		"condition": "CONDITION", "join_condition": "CONDITION",
		"filename": docs.ClassString, "keystore_location": docs.ClassString, "XMLSchema_URL": docs.ClassString,
		"datatype": docs.ClassDataType, "subquery": docs.ClassQuery, "hint": "HINT",
		"plsql_trigger_source": "TRIGGER_SOURCE", "redo_log_file_spec": docs.ClassString,
	} {
		if got := placeholderClass(name); got != want {
			t.Errorf("<%s> is %s, want %s", name, got, want)
		}
	}
}

// TestUnjudgedStatementsAreNotAsked pins what the converter leaves out
// because the pinned Oracle refuses it before it parses it: the pages of
// unjudgedPages are not in the index, and the clauses of unjudgedClauses
// are no definition of their page's syntax.
func TestUnjudgedStatementsAreNotAsked(t *testing.T) {
	t.Chdir(filepath.Join("..", "..", "..", ".."))
	if err := load(); err != nil {
		t.Fatal(err)
	}
	listed, err := read(index)
	if err != nil {
		t.Fatal(err)
	}
	bySlug := map[string]string{}
	for _, p := range pages {
		bySlug[p.slug] = p.path
	}
	for slug := range unjudgedPages {
		if strings.Contains("\n"+listed+"\n", "\n"+bySlug[slug]+"\n") {
			t.Errorf("%s is in the index", slug)
		}
	}
	for slug, clauses := range unjudgedClauses {
		for _, b := range blocks(bySlug[slug]) {
			for n := range clauses {
				if strings.HasPrefix(b, n+" ::=") {
					t.Errorf("%s still defines %s", slug, n)
				}
			}
		}
		if len(askedNot[slug]) != len(clauses) {
			t.Errorf("%s: noted %v as left out, want %d clauses", slug, askedNot[slug], len(clauses))
		}
	}
}
