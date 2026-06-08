//go:build googlesql_oracle

// Differential gate for the parser-gql node (GQL graph queries +
// CREATE PROPERTY GRAPH) against the live Cloud Spanner emulator. Run with:
//
//	SPANNER_EMULATOR_HOST=localhost:9010 \
//	  go test -tags googlesql_oracle ./googlesql/parser/ -run TestGQLOracle
//
// EMPIRICAL FINDING (probed 2026-06-05, this node) — GQL is NOT a "BigQuery-only
// blind spot" on this emulator. The analysis corpus (oracle.md) listed GQL among
// the BigQuery-only forms to triangulate; the live emulator says otherwise: it
// PARSES the GQL sub-language (Spanner natively supports Spanner Graph). A
// `GRAPH <name> MATCH … RETURN …` query parses and fails only at the resolver
// with InvalidArgument "Property graph not found: <name>" (no "Syntax error:"
// prefix), which the harness classifies ACCEPT(semantic) — exactly like
// "QUALIFY is not supported". So the positive differential is a true accept/accept
// match, and the negative differential is a true reject/reject match (the emulator
// emits a real "Syntax error:" for missing RETURN, set-op without ALL/DISTINCT,
// bad arrows, empty labels, etc.). This is a far stronger PROVE result than the
// triangulation the corpus presumed.
//
// THE ONE CLASSIFICATION CAVEAT — CREATE PROPERTY GRAPH's INNER grammar. The
// emulator runs CREATE PROPERTY GRAPH through two layers: the OUTER Spanner DDL
// grammar (its rejects carry the canonical "Error parsing Spanner DDL statement:"
// prefix → harness REJECT) and a SEPARATE ZetaSQL property-graph subparser for
// everything after the graph name, whose syntax errors come back as a plain
// "Syntax error:" WITHOUT that prefix → harness ACCEPT(semantic). So a malformed
// NODE TABLES list (missing parens / empty / no NODE TABLES at all) is a real
// inner-parser syntax error that the harness MISCLASSIFIES as accept. Those cases
// live in gqlDivergentFixtures: omni REJECTS them per the authoritative .g4 (and
// per the emulator's actual inner-parser message), the harness mislabels them
// accept, and we DEFEND omni. Divergence ledger ids 136 (IF NOT EXISTS optional +
// mutually exclusive with OR REPLACE — a TRUE outer-grammar reject the harness
// catches) and 137 (the inner-grammar misclassification). Reuses the ddlHarness
// plumbing from ddl_oracle_test.go under the shared googlesql_oracle tag.

package parser

import (
	"os"
	"testing"
)

// gqlOracleExpectation is one fixture: the SQL, the expected live-emulator
// verdict ("accept"/"reject" after harness classification), the expected omni
// accept, and a note. (Same shape as bqOracleExpectation but named distinctly to
// avoid a duplicate-symbol clash under the shared build tag.)
type gqlOracleExpectation struct {
	sql          string
	oracleAccept string // "accept" | "reject"
	omniAccept   bool
	note         string
}

// gqlAcceptFixtures — GQL forms BOTH the live emulator and omni ACCEPT. The
// emulator parses each and fails only at the resolver ("Property graph not
// found" / "Table not found"), classified ACCEPT(semantic). A true accept/accept
// differential across the whole §2.12 grammar surface.
var gqlAcceptFixtures = []gqlOracleExpectation{
	// ---- gql_statement: GRAPH <path> <ops> ----
	{"GRAPH my_graph MATCH (n) RETURN n", "accept", true, "basic MATCH node + RETURN"},
	{"GRAPH myproject.mydataset.my_graph MATCH (n) RETURN n", "accept", true, "dotted graph name path"},
	{"GRAPH g MATCH (a)-[e]->(b) RETURN a, b", "accept", true, "right-directed full edge"},
	{"GRAPH g MATCH (a)<-[e]-(b) RETURN a", "accept", true, "left-directed full edge"},
	{"GRAPH g MATCH (a)-[e]-(b) RETURN a", "accept", true, "undirected full edge (leading < absent)"},
	{"GRAPH g MATCH (a)-(b) RETURN a", "accept", true, "abbreviated undirected edge"},
	{"GRAPH g MATCH (a)<-(b) RETURN a", "accept", true, "abbreviated left edge"},
	{"GRAPH g MATCH (a)->(b) RETURN a", "accept", true, "abbreviated right edge"},
	{"GRAPH g OPTIONAL MATCH (a)-[e]->(b) RETURN a, b", "accept", true, "OPTIONAL MATCH"},
	// ---- label algebra ----
	{"GRAPH g MATCH (n:Person) RETURN n", "accept", true, "label via ':'"},
	{"GRAPH g MATCH (n IS Person) RETURN n", "accept", true, "label via IS"},
	{"GRAPH g MATCH (n:%) RETURN n", "accept", true, "wildcard label %"},
	{"GRAPH g MATCH (n:!Temp) RETURN n", "accept", true, "negated label"},
	{"GRAPH g MATCH (n:A|B&C) RETURN n", "accept", true, "label OR/AND precedence"},
	{"GRAPH g MATCH (n:(A|B)&C) RETURN n", "accept", true, "parenthesized label expr"},
	// ---- filler: properties, inline WHERE, empty ----
	{"GRAPH g MATCH (n:Person {name: 'Alice', age: 30}) RETURN n", "accept", true, "property specification"},
	{"GRAPH g MATCH (n:Person WHERE n.age > 18) RETURN n", "accept", true, "inline WHERE in filler"},
	{"GRAPH g MATCH () RETURN 1", "accept", true, "empty node pattern"},
	{"GRAPH g MATCH (a)-[e]->(b) WHERE a.id = b.id RETURN a", "accept", true, "pattern-level WHERE"},
	{"GRAPH g MATCH (a)-[]->(b), (b)-[]->(c) RETURN a, c", "accept", true, "multiple path patterns"},
	// ---- path prefixes / quantifier / parenthesized ----
	{"GRAPH g MATCH p = ANY SHORTEST TRAIL (a)-[e]->(b) RETURN p", "accept", true, "path var + search + mode prefix"},
	{"GRAPH g MATCH ANY (a)-[e]->(b) RETURN a", "accept", true, "bare ANY search prefix (no SHORTEST)"},
	{"GRAPH g MATCH ALL (a)-[e]->(b) RETURN a", "accept", true, "bare ALL search prefix"},
	{"GRAPH g MATCH (a)-[:KNOWS]->(b) RETURN a", "accept", true, "anonymous edge with label"},
	{"GRAPH g MATCH (a)-[e:KNOWS|LIKES]->(b) RETURN a", "accept", true, "edge with label disjunction in filler"},
	{"GRAPH g MATCH WALK PATHS (a)-[e]->(b) RETURN a", "accept", true, "path mode WALK PATHS"},
	{"GRAPH g MATCH (a) (-[e]->){1,3} (b) RETURN a", "accept", true, "quantified parenthesized path"},
	{"GRAPH g MATCH ((a)-[e]->(b) WHERE a.x > 0) RETURN a", "accept", true, "parenthesized path with WHERE"},
	// ---- F1: parenthesized path whose interior begins with a path-pattern prefix
	//      (path-var assignment / search / mode). graph_parenthesized_path_pattern
	//      wraps the FULL graph_path_pattern (prefixes included). ----
	{"GRAPH g MATCH (p = (a)-[e]->(b)) RETURN p", "accept", true, "F1: parenthesized path with path-variable assignment"},
	{"GRAPH g MATCH (ANY (a)-[e]->(b)) RETURN a", "accept", true, "F1: parenthesized path with ANY search prefix"},
	{"GRAPH g MATCH (ALL SHORTEST (a)-[e]->(b)) RETURN a", "accept", true, "F1: parenthesized path with ALL SHORTEST search prefix"},
	{"GRAPH g MATCH (WALK (a)-[e]->(b)) RETURN a", "accept", true, "F1: parenthesized path with WALK path-mode prefix"},
	// ---- F2: a node pattern whose filler starts with a hint. ----
	{"GRAPH g MATCH (@{force_index=idx} v:Person) RETURN v", "accept", true, "F2: hinted node pattern filler"},
	// ---- F6: inline WHERE in a filler with NO element identifier (opt_graph_element_identifier absent). ----
	{"GRAPH g MATCH (:Person WHERE TRUE) RETURN 1", "accept", true, "F6: labeled filler with WHERE, no identifier"},
	{"GRAPH g MATCH (WHERE foo) RETURN 1", "accept", true, "F6: bare WHERE filler, no identifier"},
	// ---- F5: an inter-factor hint FOLLOWED by a factor is valid. ----
	{"GRAPH g MATCH (a) @{h=1} (b) RETURN *", "accept", true, "F5: inter-factor hint followed by a node factor"},
	// ---- linear operators ----
	{"GRAPH g MATCH (n) LET x = n.age FILTER x > 18 RETURN x", "accept", true, "LET + bare FILTER"},
	{"GRAPH g MATCH (n) FILTER WHERE n.age > 18 RETURN n", "accept", true, "FILTER WHERE form"},
	{"GRAPH g MATCH (n) ORDER BY n.age DESC, n.name ASCENDING RETURN n", "accept", true, "ORDER BY operator w/ ASCENDING"},
	{"GRAPH g MATCH (n) ORDER BY n.name COLLATE 'und:ci' DESC RETURN n", "accept", true, "ORDER BY with COLLATE"},
	{"GRAPH g MATCH (n) WITH DISTINCT n.age AS age GROUP BY age RETURN age", "accept", true, "WITH + GROUP BY"},
	{"GRAPH g MATCH (n) FOR x IN n.items WITH OFFSET AS pos RETURN x, pos", "accept", true, "FOR ... WITH OFFSET"},
	{"GRAPH g MATCH (n) LIMIT 10 RETURN n", "accept", true, "standalone PAGE (LIMIT) operator"},
	{"GRAPH g MATCH (n) OFFSET 5 LIMIT 10 RETURN n", "accept", true, "standalone PAGE (OFFSET..LIMIT)"},
	{"GRAPH g MATCH (n) TABLESAMPLE RESERVOIR (100 ROWS) RETURN n", "accept", true, "TABLESAMPLE rows"},
	{"GRAPH g MATCH (n) TABLESAMPLE BERNOULLI (10 PERCENT) WITH WEIGHT AS w REPEATABLE (42) RETURN n", "accept", true, "TABLESAMPLE percent + WITH WEIGHT + REPEATABLE"},
	// ---- restricted numeric operands — the VALID forms (int/param/cast/float) ----
	{"GRAPH g MATCH (n) RETURN n LIMIT @p", "accept", true, "page LIMIT parameter (possibly_cast_int_literal_or_parameter)"},
	{"GRAPH g MATCH (n) RETURN n LIMIT CAST(@p AS INT64)", "accept", true, "page LIMIT cast"},
	{"GRAPH g MATCH (a) (-[e]->){@p,3} (b) RETURN a", "accept", true, "quantifier bound parameter (int_literal_or_parameter)"},
	{"GRAPH g MATCH (n) TABLESAMPLE BERNOULLI (10.5 PERCENT) RETURN n", "accept", true, "sample size float (sample_size_value)"},
	// ---- RETURN tails / composite ----
	{"GRAPH g MATCH (n) RETURN *", "accept", true, "RETURN star"},
	{"GRAPH g MATCH (n) RETURN n ORDER BY n.id OFFSET 5 LIMIT 3", "accept", true, "RETURN with ORDER BY + PAGE"},
	{"GRAPH g MATCH (a) RETURN a NEXT MATCH (b) RETURN b", "accept", true, "NEXT-separated composite blocks"},
	{"GRAPH g MATCH (a) RETURN a UNION ALL MATCH (b) RETURN b", "accept", true, "set-op UNION ALL"},
	{"GRAPH g MATCH (a) RETURN a INTERSECT DISTINCT MATCH (b) RETURN b", "accept", true, "set-op INTERSECT DISTINCT"},

	// ---- create_property_graph_statement (DDL) — OUTER grammar accepts ----
	{"CREATE PROPERTY GRAPH g NODE TABLES (Person, Account)", "accept", true, "basic NODE TABLES"},
	{"CREATE OR REPLACE PROPERTY GRAPH g NODE TABLES (T)", "accept", true, "OR REPLACE alone"},
	{"CREATE PROPERTY GRAPH IF NOT EXISTS g NODE TABLES (T)", "accept", true, "IF NOT EXISTS alone"},
	{"CREATE PROPERTY GRAPH g OPTIONS(x=true) NODE TABLES (T)", "accept", true, "OPTIONS"},
	{"CREATE PROPERTY GRAPH g NODE TABLES (Person KEY (id), Account KEY (id)) EDGE TABLES (Owns KEY (person_id, account_id) SOURCE KEY (person_id) REFERENCES Person (id) DESTINATION KEY (account_id) REFERENCES Account (id))", "accept", true, "EDGE TABLES with SOURCE/DESTINATION"},
	{"CREATE PROPERTY GRAPH g NODE TABLES (T PROPERTIES ARE ALL COLUMNS EXCEPT (secret))", "accept", true, "PROPERTIES ARE ALL COLUMNS EXCEPT"},
	{"CREATE PROPERTY GRAPH g NODE TABLES (T LABEL L NO PROPERTIES)", "accept", true, "LABEL ... NO PROPERTIES"},
	{"CREATE PROPERTY GRAPH g NODE TABLES (A, B,)", "accept", true, "trailing comma in element list"},
}

// gqlRejectFixtures — forms BOTH the live emulator and omni REJECT with a true
// syntax error. For queries the emulator emits "Syntax error:" (harness REJECT);
// for the CREATE PROPERTY GRAPH OR-REPLACE/IF-NOT-EXISTS conflict it emits the
// canonical "Error parsing Spanner DDL statement:" prefix (harness REJECT,
// ledger 136).
var gqlRejectFixtures = []gqlOracleExpectation{
	{"GRAPH g MATCH (n)", "reject", false, "missing RETURN (emulator: Expected keyword RETURN)"},
	{"GRAPH g RETURN", "reject", false, "RETURN with no items"},
	{"GRAPH g MATCH n RETURN n", "reject", false, "bare node needs parens (emulator: Expected '=' / path-var)"},
	{"GRAPH g MATCH (a)=>(b) RETURN a", "reject", false, "'=>' is not a valid edge arrow"},
	{"GRAPH g MATCH (a)<-[e]->(b) RETURN a", "reject", false, "a full edge cannot be both <- and -> (emulator: Expected '-' but got '->')"},
	{"GRAPH g MATCH (n:) RETURN n", "reject", false, "empty label after ':'"},
	{"GRAPH g MATCH (a) RETURN a UNION MATCH (b) RETURN b", "reject", false, "set-op metadata requires ALL|DISTINCT"},
	// Restricted numeric operands: a page count / quantifier bound / sample size is
	// int_literal_or_parameter (+cast/float where the grammar allows), NOT a full
	// expression — a bare identifier or arithmetic there is a true syntax error.
	{"GRAPH g MATCH (n) RETURN n LIMIT x", "reject", false, "page LIMIT must be int/param/cast, not an identifier"},
	{"GRAPH g MATCH (n) RETURN n LIMIT 1+1", "reject", false, "page LIMIT cannot be an arithmetic expression"},
	{"GRAPH g MATCH (n) LIMIT z RETURN n", "reject", false, "standalone-page LIMIT identifier"},
	{"GRAPH g MATCH (a) (-[e]->){1+1,3} (b) RETURN a", "reject", false, "quantifier bound cannot be arithmetic"},
	{"GRAPH g MATCH (n) TABLESAMPLE RESERVOIR (foo ROWS) RETURN n", "reject", false, "sample size must be numeric/param, not an identifier"},
	// F5: a hint between factors belongs to a `(hint? graph_path_factor)` group, so a
	// TRAILING hint with no factor after it is a syntax error (emulator: Expected "("
	// or "-" or "<" or -> but got the next keyword).
	{"GRAPH g MATCH (a) @{h=1} RETURN *", "reject", false, "F5: trailing path-factor hint with no following factor"},
	{"GRAPH g MATCH (a)-[e]->(b) @{h=1} RETURN *", "reject", false, "F5: trailing path-factor hint after a full path"},
	// CREATE PROPERTY GRAPH — OUTER-grammar true rejects (ledger 136).
	{"CREATE OR REPLACE PROPERTY GRAPH IF NOT EXISTS g NODE TABLES (T)", "reject", false, "OR REPLACE + IF NOT EXISTS mutually exclusive (emulator: Error parsing Spanner DDL statement: ... cannot be used with other existence modifiers)"},
	{"CREATE PROPERTY g NODE TABLES (T)", "reject", false, "PROPERTY without GRAPH (emulator outer DDL: Encountered 'PROPERTY')"},
	{"CREATE TEMP PROPERTY GRAPH g NODE TABLES (T)", "reject", false, "no create-scope before PROPERTY GRAPH (emulator outer DDL: Encountered 'TEMP')"},
}

// gqlDivergentFixtures — DEFENDED divergences (ledger 137). omni REJECTS these
// per the authoritative .g4 (NODE TABLES requires a parenthesized non-empty
// element_table_list; LABEL needs a name; KEY needs parens). The live emulator's
// INNER property-graph subparser ALSO rejects them (its messages are real
// "Syntax error: Expected ( …" / "Unexpected )" complaints) — but those errors
// lack the "Error parsing Spanner DDL statement:" prefix, so the harness
// MISCLASSIFIES them as accept(semantic). The oracle therefore did NOT truly
// decide; omni follows the .g4 and the emulator's actual inner-parser behavior.
// Recorded so any future flip (e.g. the emulator starting to emit the DDL prefix
// for these) is caught.
var gqlDivergentFixtures = []gqlOracleExpectation{
	{"CREATE PROPERTY GRAPH g", "accept", false, "no NODE TABLES — .g4 requires it; emulator inner-parser rejects but harness misclassifies accept"},
	{"CREATE PROPERTY GRAPH g NODE TABLES ()", "accept", false, "empty element list — .g4 requires >=1; misclassified accept"},
	{"CREATE PROPERTY GRAPH g NODE TABLES T", "accept", false, "element list needs parens; misclassified accept"},
	{"CREATE PROPERTY GRAPH g NODE TABLES (T LABEL)", "accept", false, "LABEL needs a name; misclassified accept"},
	{"CREATE PROPERTY GRAPH g NODE TABLES (T KEY id)", "accept", false, "KEY needs parens; misclassified accept"},
	{"CREATE PROPERTY GRAPH g EDGE TABLES (E)", "accept", false, "NODE TABLES is required first; misclassified accept"},
}

func TestGQLOracleDifferential(t *testing.T) {
	if os.Getenv("SPANNER_EMULATOR_HOST") == "" {
		t.Skip("SPANNER_EMULATOR_HOST not set; skipping live GQL oracle differential")
	}
	h := newDDLHarness(t)
	defer h.close()

	// Positive + negative: omni's accept/reject MUST equal the (correctly
	// classified) emulator verdict.
	run := func(name string, fixtures []gqlOracleExpectation, enforce bool) {
		t.Run(name, func(t *testing.T) {
			for _, fx := range fixtures {
				v := h.verdict(t, fx.sql)
				if v.Verdict == "error" {
					t.Fatalf("oracle returned ERROR (no verdict) for %q: %s %s — fix the emulator/fixture, do not proceed", fx.sql, v.Reason, v.Message)
				}
				oracleAccepts := v.Verdict == "accept"
				if oracleAccepts != (fx.oracleAccept == "accept") {
					t.Errorf("oracle verdict drift on %q: got %s (%s: %s), fixture expected %s", fx.sql, v.Verdict, v.Reason, v.Message, fx.oracleAccept)
				}
				_, errs := Parse(fx.sql)
				omniAccepts := len(errs) == 0
				if omniAccepts != fx.omniAccept {
					t.Errorf("omni accept drift on %q: omni accepts=%v, fixture expected %v; errs=%v", fx.sql, omniAccepts, fx.omniAccept, errs)
				}
				if enforce && omniAccepts != oracleAccepts {
					t.Errorf("DIVERGENCE on %q: omni accepts=%v, oracle accepts=%v (%s: %s)", fx.sql, omniAccepts, oracleAccepts, v.Reason, v.Message)
				}
			}
		})
	}

	run("Accept", gqlAcceptFixtures, true)
	run("Reject", gqlRejectFixtures, true)
	// Divergent: assert omni REJECTS and the emulator MISCLASSIFIES as accept (the
	// documented inner-grammar classification gap). enforce=false: the
	// accept/reject mismatch is EXPECTED and DEFENDED (ledger 137); we still pin
	// both sides so a future emulator change (it starting to surface the DDL
	// prefix for these → harness reject → match) trips the per-side checks.
	run("DivergentDefended", gqlDivergentFixtures, false)
}
