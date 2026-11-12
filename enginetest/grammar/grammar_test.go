package grammar

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func tokens(forms map[string]string) *Tokens { return &Tokens{forms: forms} }

func rows(t *testing.T, src string, forms map[string]string) ([]Row, []string) {
	t.Helper()
	g, err := ReadBison(src, tokens(forms))
	if err != nil {
		t.Fatal(err)
	}
	return g.Rows()
}

func render(rs []Row) string {
	var b strings.Builder
	for _, r := range rs {
		b.WriteString(r.Rule + "\t" + r.Statement + "\n")
	}
	return b.String()
}

// TestRowsAreShortestStatementsAndNearMisses pins what a row is: the
// shortest statement that reaches each alternative from the start symbol,
// labelled by its symbols, and its near-misses, labelled by mutation.
func TestRowsAreShortestStatementsAndNearMisses(t *testing.T) {
	const src = `
%token <ident> identifier "identifier" ifKwd "IF"
%token	exists "EXISTS"
%start	Stmt
%%
Stmt:
	"DROP" IfExists identifier { $$ = &ast.DropStmt{Name: "}"} }
IfExists:
	/* empty */ %prec lowest
|	"IF" "EXISTS"
%%
`
	got, uncovered := rows(t, src, map[string]string{"identifier": "t"})
	// Every alternative keeps its row, even where its statement is another's
	// (IfExists/%empty is Stmt's), and only declared keywords are swapped.
	want := strings.Join([]string{
		"Stmt/DROP_IfExists_identifier\tDROP t",
		"Stmt/DROP_IfExists_identifier!drop-1\tt",
		"Stmt/DROP_IfExists_identifier!drop-2\tDROP",
		"Stmt/DROP_IfExists_identifier!insert\tDROP t zzz",
		"IfExists/%empty\tDROP t",
		"IfExists/%empty!insert\tDROP zzz t",
		"IfExists/IF_EXISTS\tDROP IF EXISTS t",
		"IfExists/IF_EXISTS!drop-1\tDROP EXISTS t",
		"IfExists/IF_EXISTS!drop-2\tDROP IF t",
		"IfExists/IF_EXISTS!truncate\tDROP IF",
		"IfExists/IF_EXISTS!insert\tDROP IF EXISTS zzz t",
		"IfExists/IF_EXISTS!swap\tDROP EXISTS EXISTS t",
	}, "\n") + "\n"
	if render(got) != want {
		t.Errorf("rows:\n%s\nwant:\n%s", render(got), want)
	}
	if len(uncovered) != 0 {
		t.Errorf("uncovered: %v", uncovered)
	}
}

// statements returns the statement of each alternative row, by label.
func statements(t *testing.T, src string, forms map[string]string) map[string]string {
	t.Helper()
	got, _ := rows(t, src, forms)
	out := map[string]string{}
	for _, r := range got {
		if !IsNearMiss(r.Rule) {
			out[r.Rule] = r.Statement
		}
	}
	return out
}

// TestEquallyShortChoicesFollowTheGrammarsOrder pins how ties are broken,
// because a grammar's rules do not come in dependency order. Among equally
// short derivations the first alternative wins: here a Value is a name, as
// the grammar lists first, though Name is defined after Value and the
// keyword could be derived a pass earlier. Among equally short contexts the
// first statement wins: an Option sits in ALTER, the first statement that
// reaches it, though DROP's rule names it directly and ALTER's does not.
func TestEquallyShortChoicesFollowTheGrammarsOrder(t *testing.T) {
	const src = `
%token NAME
%%
Stmt:
	"ALTER" Wrapped
|	"DROP" Option
|	"SHOW" Value
Value:
	Name
|	"DEFAULT"
Name:
	NAME
Wrapped:
	Option
Option:
	"FORCE"
|	Value
%%
`
	got := statements(t, src, map[string]string{"NAME": "n"})
	for label, want := range map[string]string{
		"Stmt/SHOW_Value":  "SHOW n",
		"Stmt/DROP_Option": "DROP FORCE",
		"Option/FORCE":     "ALTER FORCE",
		"Option/Value":     "ALTER n",
		"Value/DEFAULT":    "ALTER DEFAULT",
	} {
		if got[label] != want {
			t.Errorf("%s is %q, want %q", label, got[label], want)
		}
	}
}

// TestRulesThatDeriveEachOtherTerminate pins the guard on a derivation that
// would use a rule inside itself: A is first a B and B is first an A, at the
// same cost, so each must fall through to an alternative that ends.
func TestRulesThatDeriveEachOtherTerminate(t *testing.T) {
	const src = `
%%
Stmt:
	"DO" A
A:
	B
|	"X"
B:
	A
|	"Y"
%%
`
	got := statements(t, src, nil)
	for label, want := range map[string]string{
		"Stmt/DO_A": "DO Y",
		"A/B":       "DO Y",
		"A/X":       "DO X",
		"B/A":       "DO Y",
		"B/Y":       "DO Y",
	} {
		if got[label] != want {
			t.Errorf("%s is %q, want %q", label, got[label], want)
		}
	}
}

// TestLabelsSurviveAnAddedAlternative pins that labels name symbols, not
// positions: a release that adds an alternative renames no other row.
func TestLabelsSurviveAnAddedAlternative(t *testing.T) {
	before := "%%\nStmt:\n\t\"BEGIN\"\n|\t\"COMMIT\"\n%%\n"
	after := "%%\nStmt:\n\t\"BEGIN\"\n|\t\"START\"\n|\t\"COMMIT\"\n%%\n"
	labels := func(src string) map[string]string {
		got, _ := rows(t, src, nil)
		m := map[string]string{}
		for _, r := range got {
			m[r.Rule] = r.Statement
		}
		return m
	}
	b, a := labels(before), labels(after)
	for rule, stmt := range b {
		if strings.Contains(rule, "!") {
			continue // a swap picks its neighbour from the keywords there are
		}
		if a[rule] != stmt {
			t.Errorf("%s: %q before, %q after", rule, stmt, a[rule])
		}
	}
	if a["Stmt/START"] != "START" {
		t.Errorf("the new alternative has no row: %v", a)
	}
}

func TestEqualAlternativesAreNumbered(t *testing.T) {
	got, _ := rows(t, "%%\nE:\n\t'-' E %prec neg\n|\t'-' E\n|\t'1'\n%%\n", nil)
	var labels []string
	for _, r := range got {
		if !strings.Contains(r.Rule, "!") {
			labels = append(labels, r.Rule)
		}
	}
	if got, want := strings.Join(labels, ","), "E/-_E,E/-_E#2,E/1"; got != want {
		t.Errorf("labels %s, want %s", got, want)
	}
}

// TestUnreachableTokensLeaveTheirAlternativesUncovered pins that a token the
// judge's mode never lexes produces no row: its text would lex as another
// token and exercise another production.
func TestUnreachableTokensLeaveTheirAlternativesUncovered(t *testing.T) {
	src := "%token pipes \"||\"\n%%\nE:\n\t'1'\n|\tE pipes E\n%%\n"
	got, uncovered := rows(t, src, map[string]string{"pipes": Unreachable})
	if strings.Contains(render(got), "||") {
		t.Errorf("a row uses the unreachable token:\n%s", render(got))
	}
	if strings.Join(uncovered, ",") != "E/E_pipes_E" {
		t.Errorf("uncovered %v, want E/E_pipes_E", uncovered)
	}
}

func TestIsNearMiss(t *testing.T) {
	for label, want := range map[string]bool{
		"E/!_E":           false,
		"E/E_!=_E":        false,
		"E/!_E!drop-1":    true,
		"E/E_!=_E!swap":   true,
		"E/E!truncate":    true,
		"E/E!insert":      true,
		"E/NOT_E":         false,
		"issue:BYT-1":     false,
		"E/E!drop-":       false,
		"E/E_!_insert_kw": false,
	} {
		if IsNearMiss(label) != want {
			t.Errorf("IsNearMiss(%q) = %v, want %v", label, !want, want)
		}
	}
}

func TestEmptyTokensAreWrittenAsNothing(t *testing.T) {
	got, _ := rows(t, "%%\nS:\n\t\"SELECT\" hint '1'\n%%\n", map[string]string{"hint": Empty})
	if got[0].Statement != "SELECT 1" {
		t.Errorf("got %q, want SELECT 1", got[0].Statement)
	}
}

func TestATerminalWithNoSampleFails(t *testing.T) {
	_, err := ReadBison("%%\nS:\n\t\"SELECT\" intLit\n%%\n", tokens(nil))
	if err == nil || !strings.Contains(err.Error(), "intLit") {
		t.Errorf("got %v, want an error naming intLit", err)
	}
}

func TestReadTokens(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tokens.tsv")
	if err := os.WriteFile(path, []byte("# comment\nidentifier\tident\tnote\nhint\t(empty)\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	tok, err := ReadTokens(path)
	if err != nil {
		t.Fatal(err)
	}
	if tok.forms["identifier"] != "ident" || tok.forms["hint"] != Empty {
		t.Errorf("forms %v", tok.forms)
	}
	for _, bad := range []string{"identifier\n", "identifier\t\n", "a\tb\nа\tb\na\tc\n"} {
		if err := os.WriteFile(path, []byte(bad), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, err := ReadTokens(path); err == nil {
			t.Errorf("%q: no error", bad)
		}
	}
}
