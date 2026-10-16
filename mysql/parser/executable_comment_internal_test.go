package parser

import (
	"strings"
	"testing"
)

func TestExecutableCommentDetectionStopsBeforeSplice(t *testing.T) {
	lexer := NewLexer(strings.Repeat("/*! */", 1000))
	lexer.stopAtExecutableComment = true

	lexer.NextToken()

	if !lexer.hasExecutableComment {
		t.Fatal("lexer did not detect executable comment")
	}
	if len(lexer.spliceGaps) != 0 {
		t.Fatalf("detection-only lexer performed %d splices, want 0", len(lexer.spliceGaps))
	}
}
