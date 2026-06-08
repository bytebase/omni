package parser

import "testing"

func TestScannerBasic(t *testing.T) {
	s := NewScanner("SELECT a, b FROM t1 WHERE x = 1")
	if s.Size() == 0 {
		t.Fatal("expected tokens")
	}

	// First token should be SELECT
	if s.GetTokenType() != SELECT {
		t.Errorf("expected SELECT, got %d", s.GetTokenType())
	}

	// Forward
	s.Forward()
	if !s.IsCurrentIdentifier() {
		t.Error("expected identifier after SELECT")
	}
	if s.GetTokenText() != "a" {
		t.Errorf("expected 'a', got %q", s.GetTokenText())
	}

	// Push/Pop
	s.Push()
	s.Forward() // ','
	s.Forward() // 'b'
	if s.GetTokenText() != "b" {
		t.Errorf("expected 'b', got %q", s.GetTokenText())
	}
	s.PopAndRestore()
	if s.GetTokenText() != "a" {
		t.Errorf("after restore expected 'a', got %q", s.GetTokenText())
	}

	// Backward
	s.Backward()
	if s.GetTokenType() != SELECT {
		t.Error("expected SELECT after backward")
	}

	// GetPreviousTokenType
	s.Forward() // 'a'
	if s.GetPreviousTokenType() != SELECT {
		t.Errorf("expected previous to be SELECT, got %d", s.GetPreviousTokenType())
	}

	// IsIdentifier
	if !s.IsIdentifier(IDENT) {
		t.Error("IDENT should be identifier")
	}
	if s.IsIdentifier(SELECT) {
		t.Error("SELECT should not be identifier")
	}

	// SeekOffset
	s.SeekOffset(0)
	if s.GetTokenType() != SELECT {
		t.Error("SeekOffset(0) should go to SELECT")
	}
}

func TestScannerGetFollowingText(t *testing.T) {
	sql := "WITH cte AS (SELECT 1) SELECT 2"
	s := NewScanner(sql)

	// Find WITH position
	if s.GetTokenType() != WITH {
		t.Fatal("expected WITH at start")
	}
	text := s.GetFollowingText(sql)
	if text != sql {
		t.Errorf("expected full SQL, got %q", text)
	}

	// Seek forward to SELECT
	for s.GetTokenType() != SELECT {
		if !s.Forward() {
			break
		}
	}
	text = s.GetFollowingText(sql)
	// Should get "SELECT 1) SELECT 2" (from the inner SELECT)
	if len(text) == 0 {
		t.Error("expected non-empty following text")
	}
}
