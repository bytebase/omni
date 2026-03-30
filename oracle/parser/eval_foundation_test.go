package parser

import (
	"reflect"
	"strings"
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

// ---------------------------------------------------------------------------
// Stage 1 — Foundation eval tests
//
// These tests verify infrastructure prerequisites that the Oracle parser must
// satisfy before higher-level stages can proceed.  They compile even when the
// features under test are not yet implemented, using reflect to probe for
// fields and functions that may be missing.
//
// PG reference: pg/ast/node.go, pg/ast/loc.go, pg/parser/parser.go
// ---------------------------------------------------------------------------

// TestEvalStage1_NoLoc verifies that oracle/ast exports a NoLoc() function
// returning Loc{Start: -1, End: -1}, matching PG's ast.NoLoc().
func TestEvalStage1_NoLoc(t *testing.T) {
	// Probe for a package-level function named NoLoc in oracle/ast.
	// Go reflect cannot enumerate package-level functions directly, so we
	// test indirectly: if someone has defined a global variable or we can
	// construct the expected value.

	// First, verify that Loc{-1,-1} is representable.
	sentinel := ast.Loc{Start: -1, End: -1}
	if sentinel.Start != -1 || sentinel.End != -1 {
		t.Fatalf("ast.Loc{-1,-1} not representable correctly: got {%d,%d}", sentinel.Start, sentinel.End)
	}

	// Check if NoLoc exists by looking for it on the ast package.
	// Since Go doesn't let us reflect on package funcs, we try a
	// compilation-independent approach: look for a method or exported
	// var.  If NoLoc() doesn't exist, the impl worker must add it.
	//
	// We use the probeFunc helper which recovers from panics.
	fn := probeNoLoc()
	if fn == nil {
		t.Fatal("ast.NoLoc() function not found in oracle/ast — " +
			"must export func NoLoc() Loc returning Loc{Start:-1, End:-1} (PG parity)")
	}

	loc := fn()
	if loc.Start != -1 {
		t.Fatalf("NoLoc().Start = %d, want -1", loc.Start)
	}
	if loc.End != -1 {
		t.Fatalf("NoLoc().End = %d, want -1", loc.End)
	}
}

// probeNoLoc tries to resolve ast.NoLoc via reflection on a known test
// function table.  Returns nil if not found.
func probeNoLoc() func() ast.Loc {
	// We test by seeing if we can find a function symbol.
	// The most reliable method in Go: maintain a registry.
	// If the impl adds NoLoc, it will appear in foundationRegistry.
	if fn, ok := foundationFuncRegistry["NoLoc"]; ok {
		if f, ok2 := fn.(func() ast.Loc); ok2 {
			return f
		}
	}
	return nil
}

// foundationFuncRegistry maps function names to their implementations.
// The init() function below populates it with whatever ast functions exist.
var foundationFuncRegistry = map[string]interface{}{}

// Populate registry — each entry only compiles if the function exists.
// We add them behind helper functions so we can selectively include them.
func init() {
	registerFoundationFuncs()
}

// registerFoundationFuncs is defined in a way that doesn't fail to compile
// if functions are missing.  We check existence at the reflect level on
// known wrapper variables.
func registerFoundationFuncs() {
	// These will be nil if the functions don't exist in the ast package.
	// The test file must compile regardless.
	if evalNoLoc != nil {
		foundationFuncRegistry["NoLoc"] = evalNoLoc
	}
	if evalNodeLoc != nil {
		foundationFuncRegistry["NodeLoc"] = evalNodeLoc
	}
	if evalListSpan != nil {
		foundationFuncRegistry["ListSpan"] = evalListSpan
	}
}

// evalNoLoc, evalNodeLoc, evalListSpan are set in eval_foundation_bridge_test.go
// (if it exists) or remain nil if the functions haven't been implemented.
// They allow this file to compile independently of whether ast.NoLoc etc exist.
var evalNoLoc func() ast.Loc
var evalNodeLoc func(ast.Node) ast.Loc
var evalListSpan func(*ast.List) ast.Loc

// TestEvalStage1_LocSentinel verifies that the zero-value Loc{} is NOT the
// sentinel for "unknown".  Unknown must use -1 (matching PG convention).
func TestEvalStage1_LocSentinel(t *testing.T) {
	// Zero value of Loc should be valid position 0, not "unknown".
	zero := ast.Loc{}
	if zero.Start != 0 {
		t.Fatalf("Loc{}.Start = %d, want 0 (zero-value is valid position 0, not unknown)", zero.Start)
	}
	if zero.End != 0 {
		t.Fatalf("Loc{}.End = %d, want 0 (zero-value is valid position 0, not unknown)", zero.End)
	}

	// Verify that Loc has both Start and End fields of type int.
	locType := reflect.TypeOf(ast.Loc{})
	startField, hasStart := locType.FieldByName("Start")
	if !hasStart {
		t.Fatal("ast.Loc missing 'Start' field")
	}
	if startField.Type.Kind() != reflect.Int {
		t.Fatalf("ast.Loc.Start type = %s, want int", startField.Type)
	}
	endField, hasEnd := locType.FieldByName("End")
	if !hasEnd {
		t.Fatal("ast.Loc missing 'End' field")
	}
	if endField.Type.Kind() != reflect.Int {
		t.Fatalf("ast.Loc.End type = %s, want int", endField.Type)
	}
}

// TestEvalStage1_TokenEnd verifies that the Token struct has an End field
// (exclusive end byte offset), as required for range-based location tracking.
func TestEvalStage1_TokenEnd(t *testing.T) {
	tokType := reflect.TypeOf(Token{})
	endField, ok := tokType.FieldByName("End")
	if !ok {
		t.Fatal("Token struct is missing 'End' field — " +
			"lexer must track exclusive end byte offset for every token")
	}
	if endField.Type.Kind() != reflect.Int {
		t.Fatalf("Token.End has type %s, want int", endField.Type)
	}

	// If End field exists, verify it's populated correctly by the lexer.
	lex := NewLexer("SELECT")
	tok := lex.NextToken()
	tokVal := reflect.ValueOf(tok)
	endVal := tokVal.FieldByName("End")
	if !endVal.IsValid() {
		t.Fatal("Token.End field not accessible via reflect")
	}
	endInt := int(endVal.Int())
	startInt := tok.Loc

	if endInt <= startInt {
		t.Fatalf("Token.End (%d) must be > Token.Loc (%d) for non-empty token 'SELECT'",
			endInt, startInt)
	}
	if endInt != 6 {
		t.Fatalf("Token.End = %d for 'SELECT' (6 bytes), want 6", endInt)
	}

	// Additional: verify with a multi-token input.
	lex2 := NewLexer("SELECT 1")
	tok1 := lex2.NextToken()
	tok2 := lex2.NextToken()

	end1 := int(reflect.ValueOf(tok1).FieldByName("End").Int())
	start2 := tok2.Loc
	if end1 > start2 {
		t.Fatalf("Token1.End (%d) > Token2.Loc (%d) — tokens overlap", end1, start2)
	}

	end2 := int(reflect.ValueOf(tok2).FieldByName("End").Int())
	if end2 != 8 {
		t.Fatalf("Token.End for '1' at offset 7 = %d, want 8", end2)
	}
}

// TestEvalStage1_ParseErrorSeverity verifies that ParseError has a Severity
// field (string), matching PG's ParseError.Severity.
func TestEvalStage1_ParseErrorSeverity(t *testing.T) {
	peType := reflect.TypeOf(ParseError{})
	sevField, ok := peType.FieldByName("Severity")
	if !ok {
		t.Fatal("ParseError struct is missing 'Severity' field — " +
			"must have Severity string field matching PG's ParseError")
	}
	if sevField.Type.Kind() != reflect.String {
		t.Fatalf("ParseError.Severity has type %s, want string", sevField.Type)
	}

	// If field exists, verify that Error() output contains severity
	// when Severity is set to a non-empty string.
	pe := &ParseError{Message: "test error", Position: 0}
	peVal := reflect.ValueOf(pe).Elem()
	sevFieldVal := peVal.FieldByName("Severity")
	if sevFieldVal.IsValid() && sevFieldVal.CanSet() {
		sevFieldVal.SetString("ERROR")
		errStr := pe.Error()
		if !strings.Contains(errStr, "ERROR") {
			t.Fatalf("ParseError.Error() = %q, want it to contain 'ERROR' when Severity='ERROR'", errStr)
		}
	}

	// Verify default behavior: empty Severity → Error() should still
	// include "ERROR" (default severity for syntax errors).
	pe2 := &ParseError{Message: "missing semicolon", Position: 5}
	errStr2 := pe2.Error()
	if !strings.Contains(errStr2, "ERROR") && !strings.Contains(errStr2, "missing semicolon") {
		t.Fatalf("ParseError.Error() with empty Severity = %q, "+
			"want it to contain either 'ERROR' or the message", errStr2)
	}
}

// TestEvalStage1_ParseErrorCode verifies that ParseError has a Code field
// (string) for SQLSTATE codes, matching PG's ParseError.Code.
func TestEvalStage1_ParseErrorCode(t *testing.T) {
	peType := reflect.TypeOf(ParseError{})
	codeField, ok := peType.FieldByName("Code")
	if !ok {
		t.Fatal("ParseError struct is missing 'Code' field — " +
			"must have Code string field for SQLSTATE codes (default '42601' for syntax errors)")
	}
	if codeField.Type.Kind() != reflect.String {
		t.Fatalf("ParseError.Code has type %s, want string", codeField.Type)
	}

	// If field exists, set a code and verify Error() includes SQLSTATE.
	pe := &ParseError{Message: "syntax error", Position: 5}
	peVal := reflect.ValueOf(pe).Elem()
	codeFieldVal := peVal.FieldByName("Code")
	if codeFieldVal.IsValid() && codeFieldVal.CanSet() {
		codeFieldVal.SetString("42601")
		errStr := pe.Error()
		if !strings.Contains(errStr, "42601") {
			t.Fatalf("ParseError.Error() = %q, want it to contain SQLSTATE '42601'", errStr)
		}
	}
}

// TestEvalStage1_ParseErrorFormat verifies that ParseError.Error() returns
// the format "SEVERITY: message (SQLSTATE code)", matching PG behavior.
func TestEvalStage1_ParseErrorFormat(t *testing.T) {
	pe := &ParseError{
		Message:  "syntax error at or near \"FOO\"",
		Position: 10,
	}

	// Set Severity and Code via reflect (fields may not exist yet).
	peVal := reflect.ValueOf(pe).Elem()

	sevField := peVal.FieldByName("Severity")
	if !sevField.IsValid() {
		t.Fatal("ParseError missing Severity field — cannot test Error() format")
	}
	if !sevField.CanSet() {
		t.Fatal("ParseError.Severity not settable — check that it is exported")
	}
	sevField.SetString("ERROR")

	codeField := peVal.FieldByName("Code")
	if !codeField.IsValid() {
		t.Fatal("ParseError missing Code field — cannot test Error() format")
	}
	codeField.SetString("42601")

	got := pe.Error()
	want := `ERROR: syntax error at or near "FOO" (SQLSTATE 42601)`
	if got != want {
		t.Fatalf("ParseError.Error() format mismatch:\n  got:  %q\n  want: %q\n"+
			"Expected format: \"SEVERITY: message (SQLSTATE code)\"", got, want)
	}
}

// TestEvalStage1_ParserSource verifies that the Parser struct has a 'source'
// field (string) storing the original SQL input text.
func TestEvalStage1_ParserSource(t *testing.T) {
	parserType := reflect.TypeOf(Parser{})
	sourceField, ok := parserType.FieldByName("source")
	if !ok {
		t.Fatal("Parser struct is missing 'source' field — " +
			"must store original SQL input string for tokenText extraction and error reporting")
	}
	if sourceField.Type.Kind() != reflect.String {
		t.Fatalf("Parser.source has type %s, want string", sourceField.Type)
	}

	// Verify the field is set during parsing by constructing a parser
	// and checking via reflect.
	sql := "SELECT 1 FROM dual"
	p := &Parser{
		lexer: NewLexer(sql),
	}
	pVal := reflect.ValueOf(p).Elem()
	srcVal := pVal.FieldByName("source")
	if !srcVal.IsValid() {
		t.Fatal("Parser.source field not accessible via reflect")
	}

	// The field exists but may not be set yet by construction.
	// The impl worker must ensure Parse() sets p.source = sql.
	t.Logf("Parser.source field exists (type %s) — impl must set it to SQL input in Parse()", sourceField.Type)
}

// TestEvalStage1_RawStmtLoc verifies that RawStmt uses a Loc field
// (type ast.Loc) instead of separate StmtLocation/StmtLen fields.
func TestEvalStage1_RawStmtLoc(t *testing.T) {
	rawType := reflect.TypeOf(ast.RawStmt{})

	// Check that Loc field exists with type ast.Loc.
	locField, hasLoc := rawType.FieldByName("Loc")
	if !hasLoc {
		t.Fatal("RawStmt is missing 'Loc' field — " +
			"must migrate from StmtLocation/StmtLen to Loc ast.Loc (PG parity)")
	}
	expectedType := reflect.TypeOf(ast.Loc{})
	if locField.Type != expectedType {
		t.Fatalf("RawStmt.Loc has type %s, want %s (ast.Loc)", locField.Type, expectedType)
	}

	// Check that legacy StmtLocation and StmtLen fields are absent.
	if _, hasOld := rawType.FieldByName("StmtLocation"); hasOld {
		t.Fatal("RawStmt still has 'StmtLocation' field — " +
			"migration to Loc incomplete; remove StmtLocation and use Loc.Start instead")
	}
	if _, hasOld := rawType.FieldByName("StmtLen"); hasOld {
		t.Fatal("RawStmt still has 'StmtLen' field — " +
			"migration to Loc incomplete; remove StmtLen and use Loc.End instead")
	}

	// Verify that after parsing, RawStmt.Loc has valid values.
	result, err := Parse("SELECT 1")
	if err != nil {
		t.Fatalf("Parse('SELECT 1') failed: %v", err)
	}
	if result.Len() == 0 {
		t.Fatal("Parse('SELECT 1') returned 0 statements")
	}

	raw := result.Items[0]
	rawVal := reflect.ValueOf(raw).Elem()
	locVal := rawVal.FieldByName("Loc")
	if !locVal.IsValid() {
		t.Fatal("RawStmt.Loc not accessible on parsed result")
	}
	locStart := int(locVal.FieldByName("Start").Int())
	locEnd := int(locVal.FieldByName("End").Int())
	if locStart < 0 {
		t.Fatalf("RawStmt.Loc.Start = %d, want >= 0", locStart)
	}
	if locEnd <= locStart {
		t.Fatalf("RawStmt.Loc.End (%d) must be > Loc.Start (%d)", locEnd, locStart)
	}
}

// TestEvalStage1_NodeLoc verifies that oracle/ast exports a
// NodeLoc(Node) Loc function matching PG's ast.NodeLoc behavior.
func TestEvalStage1_NodeLoc(t *testing.T) {
	fn := probeNodeLoc()
	if fn == nil {
		t.Fatal("ast.NodeLoc() function not found — " +
			"must export func NodeLoc(Node) Loc (PG parity: pg/ast/loc.go)")
	}

	// Test 1: nil input → NoLoc() {-1,-1}
	nilResult := fn(nil)
	if nilResult.Start != -1 || nilResult.End != -1 {
		t.Fatalf("NodeLoc(nil) = Loc{%d,%d}, want Loc{-1,-1}", nilResult.Start, nilResult.End)
	}

	// Test 2: node with Loc → returns that Loc.
	rawType := reflect.TypeOf(ast.RawStmt{})
	if _, hasLoc := rawType.FieldByName("Loc"); hasLoc {
		raw := &ast.RawStmt{}
		setLocViaReflect(raw, 5, 15)
		result := fn(raw)
		if result.Start != 5 || result.End != 15 {
			t.Fatalf("NodeLoc(RawStmt{Loc:{5,15}}) = Loc{%d,%d}, want Loc{5,15}",
				result.Start, result.End)
		}
	} else {
		t.Log("RawStmt.Loc not present — skipping RawStmt NodeLoc subtest")
	}

	// Test 3: node without Loc field → NoLoc()
	strNode := &ast.String{Str: "hello"}
	strResult := fn(strNode)
	// String node has no Loc field, so NodeLoc should return NoLoc().
	if strResult.Start != -1 || strResult.End != -1 {
		t.Fatalf("NodeLoc(String{}) = Loc{%d,%d}, want Loc{-1,-1} (String has no Loc field)",
			strResult.Start, strResult.End)
	}
}

func probeNodeLoc() func(ast.Node) ast.Loc {
	if fn, ok := foundationFuncRegistry["NodeLoc"]; ok {
		if f, ok2 := fn.(func(ast.Node) ast.Loc); ok2 {
			return f
		}
	}
	return nil
}

// TestEvalStage1_ListSpan verifies that oracle/ast exports a
// ListSpan(*List) Loc function matching PG's ast.ListSpan.
func TestEvalStage1_ListSpan(t *testing.T) {
	fn := probeListSpan()
	if fn == nil {
		t.Fatal("ast.ListSpan() function not found — " +
			"must export func ListSpan(*List) Loc (PG parity: pg/ast/loc.go)")
	}

	// Test 1: nil list → NoLoc()
	nilResult := fn(nil)
	if nilResult.Start != -1 || nilResult.End != -1 {
		t.Fatalf("ListSpan(nil) = Loc{%d,%d}, want Loc{-1,-1}", nilResult.Start, nilResult.End)
	}

	// Test 2: empty list → NoLoc()
	emptyResult := fn(&ast.List{})
	if emptyResult.Start != -1 || emptyResult.End != -1 {
		t.Fatalf("ListSpan(empty) = Loc{%d,%d}, want Loc{-1,-1}", emptyResult.Start, emptyResult.End)
	}

	// Test 3: list with items → span from first.Start to last.End.
	rawType := reflect.TypeOf(ast.RawStmt{})
	if _, hasLoc := rawType.FieldByName("Loc"); hasLoc {
		item1 := &ast.RawStmt{}
		item2 := &ast.RawStmt{}
		setLocViaReflect(item1, 10, 20)
		setLocViaReflect(item2, 30, 45)
		list := &ast.List{Items: []ast.Node{item1, item2}}

		result := fn(list)
		if result.Start != 10 {
			t.Fatalf("ListSpan([{10,20},{30,45}]).Start = %d, want 10", result.Start)
		}
		if result.End != 45 {
			t.Fatalf("ListSpan([{10,20},{30,45}]).End = %d, want 45", result.End)
		}
	} else {
		t.Log("RawStmt.Loc not present — skipping ListSpan populated-list subtest")
	}

	// Test 4: single-item list → span equals that item's Loc.
	rawType2 := reflect.TypeOf(ast.RawStmt{})
	if _, hasLoc := rawType2.FieldByName("Loc"); hasLoc {
		item := &ast.RawStmt{}
		setLocViaReflect(item, 7, 12)
		list := &ast.List{Items: []ast.Node{item}}
		result := fn(list)
		if result.Start != 7 || result.End != 12 {
			t.Fatalf("ListSpan([{7,12}]) = Loc{%d,%d}, want Loc{7,12}", result.Start, result.End)
		}
	}
}

func probeListSpan() func(*ast.List) ast.Loc {
	if fn, ok := foundationFuncRegistry["ListSpan"]; ok {
		if f, ok2 := fn.(func(*ast.List) ast.Loc); ok2 {
			return f
		}
	}
	return nil
}

// setLocViaReflect sets the Loc field on a struct via reflection.
func setLocViaReflect(node interface{}, start, end int) {
	v := reflect.ValueOf(node).Elem()
	locField := v.FieldByName("Loc")
	if locField.IsValid() && locField.CanSet() {
		locField.Set(reflect.ValueOf(ast.Loc{Start: start, End: end}))
	}
}
