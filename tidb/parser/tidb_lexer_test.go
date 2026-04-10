package parser

import "testing"

func TestTiDBKeywords(t *testing.T) {
	tidbKeywords := []string{
		"AUTO_RANDOM", "AUTO_RANDOM_BASE", "SHARD_ROW_ID_BITS",
		"PRE_SPLIT_REGIONS", "AUTO_ID_CACHE", "CLUSTERED",
		"NONCLUSTERED", "TIFLASH", "TTL", "TTL_ENABLE",
		"TTL_JOB_INTERVAL", "PLACEMENT", "POLICY",
	}
	for _, kw := range tidbKeywords {
		t.Run(kw, func(t *testing.T) {
			lex := NewLexer(kw)
			tok := lex.NextToken()
			if tok.Type == tokIDENT {
				t.Errorf("keyword %s was lexed as plain identifier, expected keyword token", kw)
			}
			if tok.Type == tokEOF {
				t.Errorf("keyword %s produced EOF", kw)
			}
		})
	}
}

func TestTiDBKeywordsAsIdentifiers(t *testing.T) {
	// TiDB keywords should be usable as identifiers (kwCatUnambiguous)
	tidbKeywords := []string{
		"AUTO_RANDOM", "CLUSTERED", "NONCLUSTERED", "TIFLASH",
		"TTL", "PLACEMENT", "POLICY",
	}
	for _, kw := range tidbKeywords {
		t.Run(kw, func(t *testing.T) {
			// Should parse as a column name in: SELECT auto_random FROM t
			sql := "SELECT " + kw + " FROM t"
			_, err := Parse(sql)
			if err != nil {
				t.Errorf("keyword %s not usable as identifier: %v", kw, err)
			}
		})
	}
}
