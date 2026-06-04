package parser

import (
	"testing"
)

// This file is the expr-json node's correctness gate (correctness-protocol.md).
// It proves the SQL/JSON function parser in json.go against three layers:
//
//  1. Accept/reject corpus — every documented Trino-481 SQL/JSON form (truth1:
//     docs/migration/trino/truth1/types-expr-functions.md json-* sections) plus
//     the legacy-grammar forms (truth2: TrinoParser.g4 jsonExists/jsonValue/
//     jsonQuery/jsonObject/jsonArray and helper rules) and oracle-discovered
//     edge cases, each tagged accept or reject.
//  2. Oracle differential (TestJSON_OracleDifferential) — the authoritative
//     gate: omni's ParseExpression accept/reject MUST equal Trino 481's
//     accept/reject of `SELECT <expr>`. Skipped when no oracle is reachable.
//  3. Structural assertions (TestJSON_Structure) — for representative forms, the
//     produced AST has the expected shape (function kind, behavior spellings,
//     member separators, null-handling, RETURNING type), catching wrong
//     field-population that the accept/reject gate cannot see.
//
// Completeness checklist (every grammar rule + doc form → its test, 100%):
//
//	jsonExists                  → exists_* accept cases + missing-path/bare-ON reject
//	jsonValue                   → value_* accept cases (RETURNING, ON EMPTY/ERROR, DEFAULT)
//	jsonQuery                   → query_* accept cases (RETURNING+FORMAT, WRAPPER, QUOTES, ON EMPTY/ERROR)
//	jsonObject                  → object_* accept cases (: and VALUE and KEY members, null/uniqueness/RETURNING)
//	jsonArray                   → array_* accept cases (elements, null-handling, RETURNING+FORMAT)
//	jsonPathInvocation          → every function's path arg; PASSING single/multi
//	jsonValueExpression         → FORMAT-tagged input/value/element/arg/passing
//	jsonRepresentation          → FORMAT JSON [ENCODING UTF8|UTF16|UTF32]; FORMAT XML / bad encoding reject
//	jsonArgument                → PASSING <expr> [FORMAT] AS <ident>, single & multi
//	jsonExistsErrorBehavior     → TRUE|FALSE|UNKNOWN|ERROR ON ERROR
//	jsonValueBehavior           → ERROR|NULL|DEFAULT expr (ON EMPTY / ON ERROR)
//	jsonQueryWrapperBehavior    → WITHOUT [ARRAY] | WITH [CONDITIONAL|UNCONDITIONAL] [ARRAY]
//	jsonQueryBehavior           → ERROR|NULL|EMPTY (ARRAY|OBJECT)
//	jsonObjectMember            → KEY? expr VALUE val | expr : val  (incl. KEY-as-expr J4)
//	clause-ordering (J1)        → out-of-order clause reject cases
//	path-literal (J2)           → non-literal path reject cases
//	reserved-name (J3)          → bare JSON_* name reject (covered in expr node; re-asserted here)
//	null-vs-element (J5)        → array_null_* cases + JSON_ARRAY(NULL ON NULL) reject

// jsonAcceptCorpus is every SQL/JSON expression form omni must accept, drawn from
// the truth1 docs (each json-* section's example_sql), the legacy grammar, and
// the oracle probes. The oracle differential proves these actually parse in
// Trino 481; here they are also smoke-tested oracle-free for parse success.
var jsonAcceptCorpus = []string{
	// --- JSON_EXISTS (jsonExists) ---
	"JSON_EXISTS(description, 'lax $.children[*]?(@ > 10)')",
	"JSON_EXISTS(doc, 'strict $.name' PASSING 'Alice' AS name_param)",
	"JSON_EXISTS(col, '$.x' FALSE ON ERROR)",
	"JSON_EXISTS(col, '$.x' TRUE ON ERROR)",
	"JSON_EXISTS(col, '$.x' UNKNOWN ON ERROR)",
	"JSON_EXISTS(col, '$.x' ERROR ON ERROR)",
	"JSON_EXISTS(col FORMAT JSON, '$.x')",
	"JSON_EXISTS(col FORMAT JSON ENCODING UTF8, '$.x')",
	"JSON_EXISTS(col, '$.x' PASSING 1 AS a, 2 AS b)",
	"JSON_EXISTS(c || d, '$.x')",
	"JSON_EXISTS(CAST(c AS varchar) FORMAT JSON, '$.x')",
	"JSON_EXISTS(c, '$.x' PASSING d FORMAT JSON AS doc)",
	"JSON_EXISTS(c, '$.x' PASSING d FORMAT JSON ENCODING UTF16 AS doc, 5 AS num)",

	// --- JSON_VALUE (jsonValue) ---
	"JSON_VALUE(description, 'lax $.children[0]' RETURNING tinyint)",
	"JSON_VALUE(doc, '$.name' DEFAULT 'unknown' ON EMPTY)",
	"JSON_VALUE(col, 'strict $.id' ERROR ON ERROR)",
	"JSON_VALUE(col, '$.x' NULL ON EMPTY NULL ON ERROR)",
	"JSON_VALUE(col, '$.x' DEFAULT 1 ON EMPTY DEFAULT 2 ON ERROR)",
	"JSON_VALUE(col, '$.x' RETURNING varchar DEFAULT 'a' ON EMPTY)",
	"JSON_VALUE(c, '$.x' DEFAULT 1 + 2 ON EMPTY)",
	"JSON_VALUE(c, '$.x' DEFAULT 'a' || 'b' ON ERROR)",
	"JSON_VALUE(c, '$.x' RETURNING DECIMAL(10,2))",
	"JSON_VALUE(c, '$.x' RETURNING int ERROR ON EMPTY NULL ON ERROR)",
	"JSON_VALUE(c, '$.x')", // bare (no RETURNING/behaviors)
	"JSON_VALUE(c FORMAT JSON, '$.x')",
	"JSON_VALUE(c, '$.x' PASSING a AS x, b FORMAT JSON AS y)",

	// --- JSON_QUERY (jsonQuery) ---
	"JSON_QUERY(description, 'lax $.children[last]' WITH ARRAY WRAPPER)",
	"JSON_QUERY(doc, 'strict $.tags' RETURNING varchar)",
	"JSON_QUERY(col, 'lax $.x' OMIT QUOTES ON SCALAR STRING NULL ON EMPTY)",
	"JSON_QUERY(col, '$.x' WITHOUT WRAPPER)",
	"JSON_QUERY(col, '$.x' WITHOUT ARRAY WRAPPER)",
	"JSON_QUERY(col, '$.x' WITH CONDITIONAL ARRAY WRAPPER)",
	"JSON_QUERY(col, '$.x' WITH UNCONDITIONAL WRAPPER)",
	"JSON_QUERY(col, '$.x' KEEP QUOTES)",
	"JSON_QUERY(col, '$.x' EMPTY ARRAY ON EMPTY EMPTY OBJECT ON ERROR)",
	"JSON_QUERY(col, '$.x' RETURNING varchar FORMAT JSON WITH ARRAY WRAPPER)",
	"JSON_QUERY(c, '$.x' KEEP QUOTES ON SCALAR STRING)",
	"JSON_QUERY(c, '$.x' RETURNING ARRAY(int))",
	"JSON_QUERY(c, '$.x' WITH ARRAY WRAPPER KEEP QUOTES ERROR ON EMPTY ERROR ON ERROR)",
	"JSON_QUERY(c, '$.x')",              // bare (no optional clauses)
	"JSON_QUERY(c, '$.x' WITH WRAPPER)", // WITH WRAPPER, no ARRAY/conditionality
	"JSON_QUERY(c FORMAT JSON, '$.x' NULL ON EMPTY)",

	// --- JSON_OBJECT (jsonObject) ---
	"JSON_OBJECT('x' : true, 'y' : 1.2, 'z' : 'text')",
	"JSON_OBJECT(KEY 'name' VALUE 'Alice', KEY 'age' VALUE 30)",
	"JSON_OBJECT('x' : null, 'y' : 1 ABSENT ON NULL)",
	"JSON_OBJECT('a' : 1 WITH UNIQUE KEYS)",
	"JSON_OBJECT('a' VALUE 1)",
	"JSON_OBJECT()",
	"JSON_OBJECT('a' : 1 NULL ON NULL RETURNING varchar)",
	"JSON_OBJECT('a' : 1 WITHOUT UNIQUE)",
	"JSON_OBJECT('a' : 1 WITH UNIQUE)",
	"JSON_OBJECT('k' VALUE 1 FORMAT JSON)",
	"JSON_OBJECT(KEY 'a' VALUE 1, 'b' : 2)",
	"JSON_OBJECT(KEY VALUE 1)",   // KEY as the key expression (J4)
	"JSON_OBJECT(k VALUE 1)",     // bare column key
	"JSON_OBJECT(KEY k VALUE 1)", // KEY keyword + column key
	"JSON_OBJECT(KEY 'x' VALUE 1)",
	"JSON_OBJECT(KEY KEY VALUE 1)", // KEY keyword + KEY-as-column key
	"JSON_OBJECT(value VALUE 1)",   // VALUE as the key expression
	"JSON_OBJECT('a' || 'b' VALUE 1)",
	"JSON_OBJECT(concat('a','b') VALUE 1)",
	"JSON_OBJECT('a' : 1 + 2)",
	"JSON_OBJECT('a' VALUE 1 + 2)",
	"JSON_OBJECT('k' VALUE col FORMAT JSON)",
	"JSON_OBJECT('k' : col FORMAT JSON ENCODING UTF8)",
	"JSON_OBJECT(RETURNING varchar)",
	"JSON_OBJECT('a' : 1, KEY 'b' VALUE 2, 'c' VALUE 3)",
	"JSON_OBJECT('a':1 WITHOUT UNIQUE KEYS)",
	"JSON_OBJECT(KEY 'a' VALUE 1 NULL ON NULL)",        // KEY member then null-handling
	"JSON_OBJECT('a':1 RETURNING varchar FORMAT JSON)", // RETURNING type + FORMAT (allowed here, unlike JSON_VALUE)

	// --- JSON_ARRAY (jsonArray) ---
	"JSON_ARRAY(true, 12e-1, 'text')",
	"JSON_ARRAY(true, null, 1)",
	"JSON_ARRAY(1, 2, 3 NULL ON NULL)",
	"JSON_ARRAY(x, y RETURNING varbinary FORMAT JSON ENCODING UTF8)",
	"JSON_ARRAY()",
	"JSON_ARRAY(1 ABSENT ON NULL RETURNING varchar)",
	"JSON_ARRAY(NULL)",           // single NULL element (J5)
	"JSON_ARRAY(1, NULL)",        // NULL as a later element
	"JSON_ARRAY(1 NULL ON NULL)", // element + NULL-ON-NULL (J5)
	"JSON_ARRAY(RETURNING varchar)",
	"JSON_ARRAY(absent)",           // ABSENT (non-reserved) as a column element
	"JSON_ARRAY(x ABSENT ON NULL)", // element + ABSENT-ON-NULL handling (J5)
	"JSON_ARRAY(1 RETURNING varchar FORMAT JSON)",

	// --- nesting & composition ---
	"JSON_VALUE(JSON_QUERY(c, '$.a'), '$.b')",
	"JSON_OBJECT('k' : JSON_QUERY(c, '$.a'))",
	"JSON_ARRAY(JSON_VALUE(c, '$.a'), JSON_VALUE(c, '$.b'))",
}

// jsonRejectCorpus is every SQL/JSON form omni must reject, mirroring Trino 481
// (oracle-confirmed). Covers the missing-mandatory-part errors, the fixed
// clause-ordering (J1), the path-must-be-literal rule (J2), the reserved bare
// name (J3), the EMPTY-without-ARRAY/OBJECT error, the null-handling-without-a-
// member error, and the bad-FORMAT / bad-ENCODING errors.
var jsonRejectCorpus = []string{
	// JSON_EXISTS: path is mandatory; ON ERROR needs a behavior keyword.
	"JSON_EXISTS(col)",
	"JSON_EXISTS(col, '$.x' ON ERROR)",
	"JSON_EXISTS(c, somecol)",                          // path must be a string literal (J2)
	"JSON_EXISTS(c, '$.x' PASSING)",                    // PASSING needs at least one arg
	"JSON_EXISTS(c FORMAT XML, '$.x')",                 // FORMAT only allows JSON
	"JSON_EXISTS(c FORMAT JSON ENCODING UTF99, '$.x')", // bad encoding

	// JSON_VALUE: ordering is fixed (J1) — ON ERROR cannot precede ON EMPTY;
	// neither ON EMPTY nor ON ERROR may repeat.
	"JSON_VALUE(c, '$.x' NULL ON ERROR ERROR ON EMPTY)",
	"JSON_VALUE(c, '$.x' NULL ON EMPTY NULL ON EMPTY)",
	"JSON_VALUE(c, '$.x' NULL ON ERROR NULL ON ERROR)",
	"JSON_VALUE(c, '$.x' NULL ON EMPTY NULL ON ERROR NULL ON ERROR)",
	// JSON_VALUE's RETURNING takes a bare type — no FORMAT clause (unlike
	// JSON_QUERY/JSON_OBJECT/JSON_ARRAY), so a trailing FORMAT is rejected.
	"JSON_VALUE(c, '$.x' RETURNING varchar FORMAT JSON)",
	"JSON_VALUE(c, '$.x' || 'y')", // path must be a literal (J2)

	// JSON_QUERY: ordering is fixed (J1) — QUOTES cannot precede WRAPPER.
	"JSON_QUERY(c, '$.x' KEEP QUOTES WITH ARRAY WRAPPER)",
	"JSON_QUERY(col, '$.x' ON SCALAR STRING)",          // ON SCALAR STRING needs the QUOTES clause
	"JSON_QUERY(c, '$.x' KEEP QUOTES ON SCALAR)",       // ON SCALAR requires the trailing STRING
	"JSON_QUERY(c, '$.x' KEEP)",                        // KEEP requires QUOTES
	"JSON_QUERY(c, '$.x' NULL ON EMPTY NULL ON EMPTY)", // duplicate ON EMPTY

	// JSON_OBJECT: a member needs a separator; null-handling/uniqueness cannot
	// appear without a member.
	"JSON_OBJECT('k')",
	"JSON_OBJECT('k' :)",
	"JSON_OBJECT('k' VALUE)",
	"JSON_OBJECT(NULL ON NULL)",
	"JSON_OBJECT(ABSENT ON NULL)",
	"JSON_OBJECT(KEY)",                                // KEY-as-column key, but no VALUE separator
	"JSON_OBJECT('a':1, NULL ON NULL ABSENT ON NULL)", // double null-handling

	// JSON_ARRAY: a leading null-handling clause (no element) is rejected (J5);
	// the trailing token of the clause must be NULL.
	"JSON_ARRAY(NULL ON NULL)",
	"JSON_ARRAY(NULL ON NULL ABSENT ON NULL)",
	"JSON_ARRAY(ABSENT ON NULL)",    // leading null-handling without an element (J5)
	"JSON_ARRAY(1, ABSENT ON NULL)", // null-handling clause cannot follow a comma
	"JSON_ARRAY(1 NULL ON ERROR)",
	"JSON_ARRAY(1 NULL ON)",

	// Reserved bare names (J3) — cannot be column references.
	"JSON_EXISTS",
	"JSON_VALUE",
	"JSON_QUERY",
	"JSON_OBJECT",
	"JSON_ARRAY",
}

// TestJSON_AcceptCorpusParses verifies omni accepts every form in
// jsonAcceptCorpus (oracle-free completeness smoke test). The oracle differential
// proves the verdicts actually match Trino.
func TestJSON_AcceptCorpusParses(t *testing.T) {
	for _, expr := range jsonAcceptCorpus {
		t.Run(truncateName(expr), func(t *testing.T) {
			node, errs := ParseExpression(expr)
			if len(errs) != 0 {
				t.Errorf("ParseExpression(%q) should accept, got errors: %v", expr, errs)
			}
			if node == nil {
				t.Errorf("ParseExpression(%q) returned nil node", expr)
			}
		})
	}
}

// TestJSON_RejectCorpusRejected verifies omni rejects every form in
// jsonRejectCorpus (required negative coverage per the correctness protocol).
func TestJSON_RejectCorpusRejected(t *testing.T) {
	for _, expr := range jsonRejectCorpus {
		t.Run(truncateName(expr), func(t *testing.T) {
			_, errs := ParseExpression(expr)
			if len(errs) == 0 {
				t.Errorf("ParseExpression(%q) should reject, but accepted", expr)
			}
		})
	}
}

// TestJSON_OracleDifferential is the authoritative accept/reject gate: for every
// form in both corpora, omni's ParseExpression accept/reject must equal Trino
// 481's accept/reject of `SELECT <expr>`. Skipped when no oracle is reachable.
func TestJSON_OracleDifferential(t *testing.T) {
	if testing.Short() {
		t.Skip("trino oracle: skipped in -short mode")
	}
	o := connectOracle(t)

	check := func(t *testing.T, expr string) {
		_, errs := ParseExpression(expr)
		omniAccepts := len(errs) == 0

		trinoAccepts, ok := oracleAccepts(t, o, wrapExpr(expr))
		if !ok {
			t.Skip("oracle unreachable for this case")
		}
		if omniAccepts != trinoAccepts {
			t.Errorf("MISMATCH expr=%q: omni accepts=%v (errs=%v), Trino accepts=%v",
				expr, omniAccepts, errs, trinoAccepts)
		}
	}

	t.Run("accept", func(t *testing.T) {
		for _, expr := range jsonAcceptCorpus {
			expr := expr
			t.Run(truncateName(expr), func(t *testing.T) { check(t, expr) })
		}
	})
	t.Run("reject", func(t *testing.T) {
		for _, expr := range jsonRejectCorpus {
			expr := expr
			t.Run(truncateName(expr), func(t *testing.T) { check(t, expr) })
		}
	})
}

// TestJSON_Structure asserts the parsed AST shape for representative forms: the
// accept/reject gate cannot see whether fields are populated correctly, so this
// checks the function kind, behavior spellings, member separators, null-handling
// and RETURNING type are set as expected. (Structural gate, correctness-protocol.md.)
func TestJSON_Structure(t *testing.T) {
	t.Run("exists_on_error", func(t *testing.T) {
		e := mustParseJSON(t, "JSON_EXISTS(c, '$.x' FALSE ON ERROR)")
		je, ok := e.(*JSONExistsExpr)
		if !ok {
			t.Fatalf("got %T, want *JSONExistsExpr", e)
		}
		if je.OnError != "FALSE" {
			t.Errorf("OnError = %q, want FALSE", je.OnError)
		}
		if je.Path == nil || je.Path.Path != "$.x" {
			t.Errorf("path = %+v, want path string %q", je.Path, "$.x")
		}
	})

	t.Run("exists_passing", func(t *testing.T) {
		e := mustParseJSON(t, "JSON_EXISTS(c, '$.x' PASSING 1 AS a, 2 AS b)")
		je := e.(*JSONExistsExpr)
		if got := len(je.Path.Passing); got != 2 {
			t.Fatalf("PASSING count = %d, want 2", got)
		}
		if je.Path.Passing[0].Name.Value != "a" || je.Path.Passing[1].Name.Value != "b" {
			t.Errorf("PASSING names = %q, %q, want a, b",
				je.Path.Passing[0].Name.Value, je.Path.Passing[1].Name.Value)
		}
	})

	t.Run("value_returning_default", func(t *testing.T) {
		e := mustParseJSON(t, "JSON_VALUE(c, '$.x' RETURNING varchar DEFAULT 'a' ON EMPTY)")
		jv, ok := e.(*JSONValueFunc)
		if !ok {
			t.Fatalf("got %T, want *JSONValueFunc", e)
		}
		if jv.Returning == nil {
			t.Errorf("Returning is nil, want a varchar type")
		}
		if jv.OnEmpty == nil || jv.OnEmpty.Kind != "DEFAULT" || jv.OnEmpty.Default == nil {
			t.Errorf("OnEmpty = %+v, want DEFAULT with expr", jv.OnEmpty)
		}
		if jv.OnError != nil {
			t.Errorf("OnError = %+v, want nil", jv.OnError)
		}
	})

	t.Run("query_wrapper_quotes", func(t *testing.T) {
		e := mustParseJSON(t, "JSON_QUERY(c, '$.x' WITH CONDITIONAL ARRAY WRAPPER OMIT QUOTES ON SCALAR STRING NULL ON EMPTY)")
		jq, ok := e.(*JSONQueryFunc)
		if !ok {
			t.Fatalf("got %T, want *JSONQueryFunc", e)
		}
		if jq.Wrapper != "WITH CONDITIONAL" || !jq.WrapperArray {
			t.Errorf("wrapper = %q array=%v, want WITH CONDITIONAL array=true", jq.Wrapper, jq.WrapperArray)
		}
		if jq.Quotes != "OMIT" || !jq.QuotesOnScalar {
			t.Errorf("quotes = %q onScalar=%v, want OMIT onScalar=true", jq.Quotes, jq.QuotesOnScalar)
		}
		if jq.OnEmpty == nil || jq.OnEmpty.Kind != "NULL" {
			t.Errorf("OnEmpty = %+v, want NULL", jq.OnEmpty)
		}
	})

	t.Run("query_empty_behavior", func(t *testing.T) {
		e := mustParseJSON(t, "JSON_QUERY(c, '$.x' EMPTY ARRAY ON EMPTY EMPTY OBJECT ON ERROR)")
		jq := e.(*JSONQueryFunc)
		if jq.OnEmpty == nil || jq.OnEmpty.Kind != "EMPTY ARRAY" {
			t.Errorf("OnEmpty = %+v, want EMPTY ARRAY", jq.OnEmpty)
		}
		if jq.OnError == nil || jq.OnError.Kind != "EMPTY OBJECT" {
			t.Errorf("OnError = %+v, want EMPTY OBJECT", jq.OnError)
		}
	})

	t.Run("object_members_mixed_separators", func(t *testing.T) {
		e := mustParseJSON(t, "JSON_OBJECT('a' : 1, KEY 'b' VALUE 2, 'c' VALUE 3)")
		jo, ok := e.(*JSONObjectExpr)
		if !ok {
			t.Fatalf("got %T, want *JSONObjectExpr", e)
		}
		if len(jo.Members) != 3 {
			t.Fatalf("members = %d, want 3", len(jo.Members))
		}
		if jo.Members[0].Separator != ":" || jo.Members[0].KeyKeyword {
			t.Errorf("member 0 = %+v, want ':' no-KEY", jo.Members[0])
		}
		if jo.Members[1].Separator != "VALUE" || !jo.Members[1].KeyKeyword {
			t.Errorf("member 1 = %+v, want VALUE with KEY", jo.Members[1])
		}
		if jo.Members[2].Separator != "VALUE" || jo.Members[2].KeyKeyword {
			t.Errorf("member 2 = %+v, want VALUE no-KEY", jo.Members[2])
		}
	})

	t.Run("object_key_as_expression", func(t *testing.T) {
		// JSON_OBJECT(KEY VALUE 1): KEY is the key column, not the keyword (J4).
		e := mustParseJSON(t, "JSON_OBJECT(KEY VALUE 1)")
		jo := e.(*JSONObjectExpr)
		if len(jo.Members) != 1 {
			t.Fatalf("members = %d, want 1", len(jo.Members))
		}
		m := jo.Members[0]
		if m.KeyKeyword {
			t.Errorf("KeyKeyword = true, want false (KEY is the key expression here)")
		}
		// KEY is parsed as a non-reserved-keyword identifier; its Value preserves
		// the source case ("KEY") and Normalize folds it to "key".
		if cr, ok := m.Key.(*ColumnRef); !ok || cr.Name.Normalize() != "key" {
			t.Errorf("key = %+v, want ColumnRef normalizing to key", m.Key)
		}
	})

	t.Run("object_null_handling_uniqueness", func(t *testing.T) {
		e := mustParseJSON(t, "JSON_OBJECT('a' : 1 ABSENT ON NULL WITH UNIQUE KEYS RETURNING varchar)")
		jo := e.(*JSONObjectExpr)
		if jo.NullHandling != "ABSENT ON NULL" {
			t.Errorf("NullHandling = %q, want ABSENT ON NULL", jo.NullHandling)
		}
		if jo.Uniqueness != "WITH UNIQUE" || !jo.UniqueKeys {
			t.Errorf("uniqueness = %q keys=%v, want WITH UNIQUE keys=true", jo.Uniqueness, jo.UniqueKeys)
		}
		if jo.Returning == nil {
			t.Errorf("Returning is nil, want varchar")
		}
	})

	t.Run("array_elements_null_handling", func(t *testing.T) {
		e := mustParseJSON(t, "JSON_ARRAY(1, 2, 3 NULL ON NULL)")
		ja, ok := e.(*JSONArrayExpr)
		if !ok {
			t.Fatalf("got %T, want *JSONArrayExpr", e)
		}
		if len(ja.Elements) != 3 {
			t.Fatalf("elements = %d, want 3", len(ja.Elements))
		}
		if ja.NullHandling != "NULL ON NULL" {
			t.Errorf("NullHandling = %q, want NULL ON NULL", ja.NullHandling)
		}
	})

	t.Run("array_single_null_element", func(t *testing.T) {
		// JSON_ARRAY(NULL): one NULL element, NO null-handling clause (J5).
		e := mustParseJSON(t, "JSON_ARRAY(NULL)")
		ja := e.(*JSONArrayExpr)
		if len(ja.Elements) != 1 {
			t.Fatalf("elements = %d, want 1", len(ja.Elements))
		}
		if ja.NullHandling != "" {
			t.Errorf("NullHandling = %q, want empty (NULL is the element)", ja.NullHandling)
		}
	})

	t.Run("value_input_format", func(t *testing.T) {
		e := mustParseJSON(t, "JSON_EXISTS(c FORMAT JSON ENCODING UTF16, '$.x')")
		je := e.(*JSONExistsExpr)
		if je.Path.Input.Format == nil || je.Path.Input.Format.Encoding != "UTF16" {
			t.Errorf("input format = %+v, want JSON ENCODING UTF16", je.Path.Input.Format)
		}
	})
}

// mustParseJSON parses expr and fails the test if it does not parse cleanly to a
// single Expr. Returns the parsed Expr.
func mustParseJSON(t *testing.T, expr string) Expr {
	t.Helper()
	e, errs := ParseExpression(expr)
	if len(errs) != 0 {
		t.Fatalf("ParseExpression(%q): unexpected errors: %v", expr, errs)
	}
	if e == nil {
		t.Fatalf("ParseExpression(%q): nil expr", expr)
	}
	return e
}
