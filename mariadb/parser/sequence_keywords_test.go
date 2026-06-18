package parser

import "testing"

// newSequenceKeywords are the 9 MariaDB sequence-related keywords added in P2
// (BYT-9135). The other 6 needed by the feature — next/value/maxvalue/cache/
// restart/start — are already inherited from the mysql fork base.
//
// Each must be wired at two structural sites that the Task 3/4 parser depends on:
//   - the lexer keywords map (so it lexes to a kw* token the parser can dispatch
//     on via p.cur.Type == kwX, instead of an opaque tokIDENT);
//   - keywordCategories as kwCatUnambiguous (MariaDB reserves none of them, and
//     they are usable as label/role/lvalue identifiers — the contexts that route
//     through isLabelKeyword/isRoleKeyword/isLvalueKeyword, all of which fail
//     closed for tokens absent from keywordCategories).
var newSequenceKeywords = []string{
	"sequence", "previous", "minvalue", "cycle", "increment",
	"nocache", "nocycle", "nominvalue", "nomaxvalue",
}

// TestSequenceKeywordsRegistered asserts each new keyword lexes to a registered
// token AND is classified kwCatUnambiguous. This is the structural guard for all
// three wiring sites (const + lexer map + category); a missing const/map entry
// trips the keywords lookup, a missing category entry trips the classification.
func TestSequenceKeywordsRegistered(t *testing.T) {
	for _, kw := range newSequenceKeywords {
		tok, ok := keywords[kw]
		if !ok {
			t.Errorf("keyword %q is not registered in the lexer keywords map", kw)
			continue
		}
		cat, ok := keywordCategories[tok]
		if !ok {
			t.Errorf("keyword %q (token %d) is absent from keywordCategories", kw, tok)
			continue
		}
		if cat != kwCatUnambiguous {
			t.Errorf("keyword %q: want category kwCatUnambiguous (%d), got %d", kw, kwCatUnambiguous, cat)
		}
	}
}

// TestSequenceKeywordsUsableAsColumnNames guards the fail-open column-name path,
// NOT the keywordCategories registration (TestSequenceKeywordsRegistered owns
// that). Column names parse via parseColumnDef -> parseIdent, which accepts any
// token >= 700 that is not reserved; this stays green whether or not the 9 are
// present in keywordCategories, and would only break if one were mis-marked
// kwCatReserved. It passed even in the pre-implementation RED state (the words
// lexed as tokIDENT), which is exactly why it cannot catch a missing keyword
// registration — that is the registration test's job.
func TestSequenceKeywordsUsableAsColumnNames(t *testing.T) {
	sql := "CREATE TABLE t (" +
		"sequence INT, previous INT, increment INT, cycle INT, minvalue INT, " +
		"nocache INT, nocycle INT, nominvalue INT, nomaxvalue INT)"
	if _, err := Parse(sql); err != nil {
		t.Fatalf("the 9 sequence keywords must remain usable as column names; got: %v", err)
	}
}
