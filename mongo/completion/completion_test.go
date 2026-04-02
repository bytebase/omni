package completion

import (
	"testing"

	"github.com/bytebase/omni/mongo/catalog"
)

// detectContextFromInput is a test helper that simulates what Complete() does:
// extract prefix, tokenize, then detect context.
func detectContextFromInput(input string) completionContext {
	prefix := extractPrefix(input, len(input))
	tokens := tokenize(input, len(input)-len(prefix))
	return detectContext(tokens)
}

func TestDetectContextTopLevel(t *testing.T) {
	tests := []struct {
		input string
		name  string
	}{
		{"", "empty input"},
		{"  ", "whitespace only"},
		{"db.users.find();\n", "after semicolon and newline"},
		{"db.users.find();", "after semicolon"},
		{"var x = 1;", "after variable assignment"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectContextFromInput(tt.input)
			if got != contextTopLevel {
				t.Errorf("detectContextFromInput(%q) = %d, want contextTopLevel (%d)", tt.input, got, contextTopLevel)
			}
		})
	}
}

func TestDetectContextAfterDbDot(t *testing.T) {
	tests := []struct {
		input string
		name  string
	}{
		{"db.", "db dot"},
		{"db.u", "db dot with prefix u"},
		{"db.get", "db dot with prefix get"},
		{"db.getC", "db dot with prefix getC"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectContextFromInput(tt.input)
			if got != contextAfterDbDot {
				t.Errorf("detectContextFromInput(%q) = %d, want contextAfterDbDot (%d)", tt.input, got, contextAfterDbDot)
			}
		})
	}
}

func TestDetectContextAfterCollDot(t *testing.T) {
	tests := []struct {
		input string
		name  string
	}{
		{"db.users.", "db.coll."},
		{"db.users.f", "db.coll.prefix"},
		{`db["users"].`, "db bracket coll dot"},
		{`db["users"].f`, "db bracket coll dot prefix"},
		{`db.getCollection("users").`, "db.getCollection dot"},
		{`db.getCollection("users").f`, "db.getCollection dot prefix"},
		{"db.myCollection.", "db.myCollection dot"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectContextFromInput(tt.input)
			if got != contextAfterCollDot {
				t.Errorf("detectContextFromInput(%q) = %d, want contextAfterCollDot (%d)", tt.input, got, contextAfterCollDot)
			}
		})
	}
}

func TestDetectContextAfterBracket(t *testing.T) {
	tests := []struct {
		input string
		name  string
	}{
		{"db[", "db open bracket"},
		{`db["`, "db bracket with opening quote"},
		{`db["us`, "db bracket with partial collection name"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectContextFromInput(tt.input)
			if got != contextAfterBracket {
				t.Errorf("detectContextFromInput(%q) = %d, want contextAfterBracket (%d)", tt.input, got, contextAfterBracket)
			}
		})
	}
}

func TestDetectContextCursorChain(t *testing.T) {
	tests := []struct {
		input string
		name  string
	}{
		{"db.users.find().", "after method call dot"},
		{"db.users.find().s", "after method call dot with prefix"},
		{"db.users.find({}).sort({a:1}).", "chained cursor methods"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectContextFromInput(tt.input)
			if got != contextCursorChain {
				t.Errorf("detectContextFromInput(%q) = %d, want contextCursorChain (%d)", tt.input, got, contextCursorChain)
			}
		})
	}
}

func TestDetectContextShowTarget(t *testing.T) {
	tests := []struct {
		input string
		name  string
	}{
		{"show ", "show space"},
		{"show d", "show with prefix d"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectContextFromInput(tt.input)
			if got != contextShowTarget {
				t.Errorf("detectContextFromInput(%q) = %d, want contextShowTarget (%d)", tt.input, got, contextShowTarget)
			}
		})
	}
}

func TestDetectContextAfterRsDot(t *testing.T) {
	tests := []struct {
		input string
		name  string
	}{
		{"rs.", "rs dot"},
		{"rs.s", "rs dot with prefix s"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectContextFromInput(tt.input)
			if got != contextAfterRsDot {
				t.Errorf("detectContextFromInput(%q) = %d, want contextAfterRsDot (%d)", tt.input, got, contextAfterRsDot)
			}
		})
	}
}

func TestDetectContextAfterShDot(t *testing.T) {
	tests := []struct {
		input string
		name  string
	}{
		{"sh.", "sh dot"},
		{"sh.a", "sh dot with prefix a"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectContextFromInput(tt.input)
			if got != contextAfterShDot {
				t.Errorf("detectContextFromInput(%q) = %d, want contextAfterShDot (%d)", tt.input, got, contextAfterShDot)
			}
		})
	}
}

func TestDetectContextAggStage(t *testing.T) {
	tests := []struct {
		input string
		name  string
	}{
		{"db.users.aggregate([{$", "agg stage with dollar prefix"},
		{"db.users.aggregate([{$m", "agg stage with dollar m prefix"},
		{"db.users.aggregate([{$match: {}}, {$", "agg stage after comma"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectContextFromInput(tt.input)
			if got != contextAggStage {
				t.Errorf("detectContextFromInput(%q) = %d, want contextAggStage (%d)", tt.input, got, contextAggStage)
			}
		})
	}
}

func TestDetectContextQueryOperator(t *testing.T) {
	tests := []struct {
		input string
		name  string
	}{
		{"db.users.find({age: {$", "query operator with dollar prefix"},
		{"db.users.find({age: {$g", "query operator with dollar g prefix"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectContextFromInput(tt.input)
			if got != contextQueryOperator {
				t.Errorf("detectContextFromInput(%q) = %d, want contextQueryOperator (%d)", tt.input, got, contextQueryOperator)
			}
		})
	}
}

func TestDetectContextInsideArgs(t *testing.T) {
	tests := []struct {
		input string
		name  string
	}{
		{"db.users.find(", "find open paren"},
		{"db.users.insertOne(", "insertOne open paren"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectContextFromInput(tt.input)
			if got != contextInsideArgs {
				t.Errorf("detectContextFromInput(%q) = %d, want contextInsideArgs (%d)", tt.input, got, contextInsideArgs)
			}
		})
	}
}

func TestDetectContextDocumentKey(t *testing.T) {
	tests := []struct {
		input string
		name  string
	}{
		{"db.users.find({", "open brace after paren"},
		{"db.users.insertOne({name: 1, ", "after comma inside brace"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectContextFromInput(tt.input)
			if got != contextDocumentKey {
				t.Errorf("detectContextFromInput(%q) = %d, want contextDocumentKey (%d)", tt.input, got, contextDocumentKey)
			}
		})
	}
}

// --- Test helpers ---

func newTestCatalog(names ...string) *catalog.Catalog {
	cat := catalog.New()
	for _, name := range names {
		cat.AddCollection(name)
	}
	return cat
}

func candidateTexts(candidates []Candidate) []string {
	texts := make([]string, len(candidates))
	for i, c := range candidates {
		texts[i] = c.Text
	}
	return texts
}

func hasCandidate(candidates []Candidate, text string) bool {
	for _, c := range candidates {
		if c.Text == text {
			return true
		}
	}
	return false
}

func hasCandidateWithType(candidates []Candidate, text string, typ CandidateType) bool {
	for _, c := range candidates {
		if c.Text == text && c.Type == typ {
			return true
		}
	}
	return false
}

// --- extractPrefix tests ---

func TestExtractPrefix(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		offset int
		want   string
	}{
		{"db dot", "db.", 3, ""},
		{"db dot partial", "db.us", 5, "us"},
		{"cursor chain prefix", "db.users.find().s", 17, "s"},
		{"dollar prefix", "{age: {$g", 9, "$g"},
		{"empty input", "", 0, ""},
		{"show prefix", "show d", 6, "d"},
		{"full word", "db", 2, "db"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractPrefix(tt.input, tt.offset)
			if got != tt.want {
				t.Errorf("extractPrefix(%q, %d) = %q, want %q", tt.input, tt.offset, got, tt.want)
			}
		})
	}
}

// --- filterByPrefix tests ---

func TestFilterByPrefix(t *testing.T) {
	candidates := []Candidate{
		{Text: "find", Type: CandidateMethod},
		{Text: "findOne", Type: CandidateMethod},
		{Text: "aggregate", Type: CandidateMethod},
		{Text: "$gt", Type: CandidateQueryOperator},
		{Text: "$gte", Type: CandidateQueryOperator},
	}

	t.Run("empty prefix returns all", func(t *testing.T) {
		got := filterByPrefix(candidates, "")
		if len(got) != len(candidates) {
			t.Errorf("filterByPrefix with empty prefix returned %d candidates, want %d", len(got), len(candidates))
		}
	})

	t.Run("case sensitive", func(t *testing.T) {
		got := filterByPrefix(candidates, "F")
		if len(got) != 0 {
			t.Errorf("filterByPrefix with prefix 'F' returned %d candidates, want 0 (case-sensitive)", len(got))
		}
	})

	t.Run("dollar prefix", func(t *testing.T) {
		got := filterByPrefix(candidates, "$g")
		if len(got) != 2 {
			t.Errorf("filterByPrefix with prefix '$g' returned %d candidates, want 2", len(got))
		}
		for _, c := range got {
			if c.Text != "$gt" && c.Text != "$gte" {
				t.Errorf("unexpected candidate %q for prefix '$g'", c.Text)
			}
		}
	})

	t.Run("prefix f", func(t *testing.T) {
		got := filterByPrefix(candidates, "f")
		if len(got) != 2 {
			t.Errorf("filterByPrefix with prefix 'f' returned %d candidates, want 2", len(got))
		}
	})
}

// --- Complete end-to-end tests ---

func TestCompleteAfterDbDot(t *testing.T) {
	cat := newTestCatalog("users", "orders")
	results := Complete("db.", 3, cat)

	// Should include collection names.
	if !hasCandidateWithType(results, "users", CandidateCollection) {
		t.Error("expected collection 'users' in results")
	}
	if !hasCandidateWithType(results, "orders", CandidateCollection) {
		t.Error("expected collection 'orders' in results")
	}
	// Should include db methods.
	if !hasCandidateWithType(results, "getName", CandidateDbMethod) {
		t.Error("expected db method 'getName' in results")
	}
	if !hasCandidateWithType(results, "runCommand", CandidateDbMethod) {
		t.Error("expected db method 'runCommand' in results")
	}
}

func TestCompleteCollectionMethodPrefix(t *testing.T) {
	results := Complete("db.users.f", 10, nil)

	expected := []string{"find", "findOne", "findOneAndDelete", "findOneAndReplace", "findOneAndUpdate"}
	for _, e := range expected {
		if !hasCandidateWithType(results, e, CandidateMethod) {
			t.Errorf("expected method %q in results", e)
		}
	}
}

func TestCompleteCursorChainPrefix(t *testing.T) {
	results := Complete("db.users.find().s", 17, nil)

	expected := []string{"sort", "skip", "size", "showRecordId"}
	for _, e := range expected {
		if !hasCandidateWithType(results, e, CandidateCursorMethod) {
			t.Errorf("expected cursor method %q in results", e)
		}
	}
}

func TestCompleteAggStage(t *testing.T) {
	results := Complete("db.users.aggregate([{$m", 23, nil)

	expected := []string{"$match", "$merge"}
	for _, e := range expected {
		if !hasCandidateWithType(results, e, CandidateAggStage) {
			t.Errorf("expected agg stage %q in results", e)
		}
	}
}

func TestCompleteQueryOperator(t *testing.T) {
	results := Complete("db.users.find({age: {$g", 23, nil)

	expected := []string{"$gt", "$gte", "$geoWithin", "$geoIntersects"}
	for _, e := range expected {
		if !hasCandidateWithType(results, e, CandidateQueryOperator) {
			t.Errorf("expected query operator %q in results", e)
		}
	}
}

func TestCompleteBracketWithCatalog(t *testing.T) {
	cat := newTestCatalog("system.profile", "users", "system.views")
	results := Complete(`db["sys`, 7, cat)

	if !hasCandidateWithType(results, "system.profile", CandidateCollection) {
		t.Error("expected 'system.profile' in results")
	}
	if !hasCandidateWithType(results, "system.views", CandidateCollection) {
		t.Error("expected 'system.views' in results")
	}
	if hasCandidate(results, "users") {
		t.Error("should NOT include 'users' for prefix 'sys'")
	}
}

func TestCompleteShowTarget(t *testing.T) {
	results := Complete("show d", 6, nil)

	expected := []string{"dbs", "databases"}
	for _, e := range expected {
		if !hasCandidateWithType(results, e, CandidateShowTarget) {
			t.Errorf("expected show target %q in results", e)
		}
	}
}

func TestCompleteRsMethods(t *testing.T) {
	results := Complete("rs.", 3, nil)

	expected := []string{"status", "conf", "initiate"}
	for _, e := range expected {
		if !hasCandidateWithType(results, e, CandidateRsMethod) {
			t.Errorf("expected rs method %q in results", e)
		}
	}
}

func TestCompleteShMethods(t *testing.T) {
	results := Complete("sh.", 3, nil)

	expected := []string{"addShard", "enableSharding", "status"}
	for _, e := range expected {
		if !hasCandidateWithType(results, e, CandidateShMethod) {
			t.Errorf("expected sh method %q in results", e)
		}
	}
}

func TestCompleteTopLevelEmpty(t *testing.T) {
	results := Complete("", 0, nil)

	expected := []string{"db", "rs", "sh", "show"}
	for _, e := range expected {
		if !hasCandidateWithType(results, e, CandidateKeyword) {
			t.Errorf("expected keyword %q in results", e)
		}
	}
}

func TestCompleteGetCollection(t *testing.T) {
	results := Complete(`db.getCollection("users").f`, 27, nil)

	expected := []string{"find", "findOne"}
	for _, e := range expected {
		if !hasCandidateWithType(results, e, CandidateMethod) {
			t.Errorf("expected method %q in results", e)
		}
	}
}

// --- Edge case tests ---

func TestCompleteNilCatalog(t *testing.T) {
	results := Complete("db.", 3, nil)

	// Should still have db methods.
	if !hasCandidateWithType(results, "getName", CandidateDbMethod) {
		t.Error("expected db method 'getName' with nil catalog")
	}
	// Should NOT have any collection candidates.
	for _, c := range results {
		if c.Type == CandidateCollection {
			t.Errorf("unexpected collection candidate %q with nil catalog", c.Text)
		}
	}
}

func TestCompleteCursorOvershoot(t *testing.T) {
	// Should not panic when cursor offset exceeds input length.
	results := Complete("db.", 100, nil)
	if !hasCandidateWithType(results, "getName", CandidateDbMethod) {
		t.Error("expected db method 'getName' with overshooting cursor")
	}
}

func TestCompleteBracketWithCatalogQuote(t *testing.T) {
	cat := newTestCatalog("users", "orders")
	results := Complete(`db["`, 4, cat)

	if !hasCandidateWithType(results, "users", CandidateCollection) {
		t.Error("expected 'users' in bracket completion")
	}
	if !hasCandidateWithType(results, "orders", CandidateCollection) {
		t.Error("expected 'orders' in bracket completion")
	}
}

// --- Negative tests ---

func TestCompleteNegativeCases(t *testing.T) {
	t.Run("collection method prefix should not include unrelated", func(t *testing.T) {
		results := Complete("db.users.f", 10, nil)

		if hasCandidate(results, "aggregate") {
			t.Error("should NOT include 'aggregate' for prefix 'f'")
		}
		if hasCandidate(results, "sort") {
			t.Error("should NOT include cursor method 'sort' in collection methods")
		}
	})

	t.Run("case sensitivity", func(t *testing.T) {
		results := Complete("db.users.F", 10, nil)

		if hasCandidate(results, "find") {
			t.Error("should NOT match 'find' for prefix 'F' (case-sensitive)")
		}
		if hasCandidate(results, "findOne") {
			t.Error("should NOT match 'findOne' for prefix 'F' (case-sensitive)")
		}
	})
}
