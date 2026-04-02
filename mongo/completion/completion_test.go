package completion

import (
	"testing"
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
