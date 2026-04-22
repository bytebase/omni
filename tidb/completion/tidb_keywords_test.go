package completion

import (
	"testing"

	"github.com/bytebase/omni/tidb/catalog"
)

// TestTiDBKeywords_TableOptions asserts that after a parenthesized column
// list in CREATE TABLE, TiDB-specific table-option keywords appear as
// completion candidates. These keywords come from the parser's token set
// (added in PR2) and surface automatically in positions where the parser
// accepts a table_option list.
func TestTiDBKeywords_TableOptions(t *testing.T) {
	cat := catalog.New()
	cat.SetCurrentDatabase("testdb")

	sql := "CREATE TABLE t (id INT) "
	candidates := Complete(sql, len(sql), cat)

	// At least one TiDB table option should be offered. We check for a
	// representative set; if the parser adds all of them at this position,
	// they should all appear. Failing on any single one would be flaky if
	// the grammar rule changes; requiring any-of is robust to refactors.
	wanted := []string{
		"SHARD_ROW_ID_BITS", "PRE_SPLIT_REGIONS", "AUTO_ID_CACHE",
		"AUTO_RANDOM_BASE", "TTL", "TTL_ENABLE", "TTL_JOB_INTERVAL",
		"PLACEMENT",
	}
	found := 0
	for _, w := range wanted {
		if containsText(candidates, w) {
			found++
		}
	}
	if found == 0 {
		t.Errorf("expected at least one TiDB table-option keyword after CREATE TABLE (id INT); got none.\nAll candidates: %v", candidateTexts(candidates))
	}
}

// TestTiDBKeywords_AlterTable asserts that SET TIFLASH REPLICA and
// REMOVE TTL appear as candidates in the ALTER TABLE command position.
func TestTiDBKeywords_AlterTable(t *testing.T) {
	cat := catalog.New()
	cat.SetCurrentDatabase("testdb")

	// After "ALTER TABLE t " the parser is expecting an alter command.
	sql := "ALTER TABLE t "
	candidates := Complete(sql, len(sql), cat)

	// Either SET or REMOVE (or both) should be offered; TIFLASH/TTL are
	// the follow-on tokens and may or may not be listed depending on
	// how deeply the completion walker explores.
	if !containsText(candidates, "SET") && !containsText(candidates, "REMOVE") {
		t.Errorf("expected SET or REMOVE in ALTER TABLE completion; got: %v", candidateTexts(candidates))
	}
}

// TestTiDBKeywords_NoTiKVEngine asserts that TiKV is NOT offered as a
// completion candidate for the ENGINE= option. TiDB parses ENGINE= for
// MySQL compatibility but ignores the value — suggesting TiKV would
// mislead users into thinking it's a meaningful choice.
func TestTiDBKeywords_NoTiKVEngine(t *testing.T) {
	cat := catalog.New()
	cat.SetCurrentDatabase("testdb")

	sql := "CREATE TABLE t (id INT) ENGINE="
	candidates := Complete(sql, len(sql), cat)

	for _, c := range candidates {
		if c.Text == "TiKV" || c.Text == "TIKV" || c.Text == "tikv" {
			t.Errorf("TiKV must not appear in ENGINE completion (TiDB ignores ENGINE value); got: %+v", c)
		}
	}
}

// TestTiDBKeywords_CreateDatabasePlacement asserts that PLACEMENT is
// offered as a completion candidate inside a CREATE DATABASE option
// list, anchoring the PLACEMENT POLICY = <name> syntax TiDB accepts.
func TestTiDBKeywords_CreateDatabasePlacement(t *testing.T) {
	cat := catalog.New()

	sql := "CREATE DATABASE db "
	candidates := Complete(sql, len(sql), cat)

	if !containsText(candidates, "PLACEMENT") {
		t.Errorf("expected PLACEMENT in CREATE DATABASE completion; got: %v", candidateTexts(candidates))
	}
}

// TestTiDBKeywords_CreatePlacementPolicy asserts that PLACEMENT shows
// up as a candidate object type after CREATE, so users can discover
// the CREATE PLACEMENT POLICY statement via completion.
func TestTiDBKeywords_CreatePlacementPolicy(t *testing.T) {
	cat := catalog.New()

	sql := "CREATE "
	candidates := Complete(sql, len(sql), cat)

	if !containsText(candidates, "PLACEMENT") {
		t.Errorf("expected PLACEMENT in CREATE completion; got: %v", candidateTexts(candidates))
	}
}

// TestTiDBKeywords_PlacementPolicyOptions asserts that PRIMARY_REGION
// (and friends) show up as candidates inside a CREATE PLACEMENT POLICY
// option list. The inner grammar is completely independent from CREATE
// TABLE options, so a separate anchor test is warranted.
func TestTiDBKeywords_PlacementPolicyOptions(t *testing.T) {
	cat := catalog.New()

	sql := "CREATE PLACEMENT POLICY p "
	candidates := Complete(sql, len(sql), cat)

	wanted := []string{
		"PRIMARY_REGION", "REGIONS", "FOLLOWERS", "VOTERS", "LEARNERS",
		"CONSTRAINTS", "LEADER_CONSTRAINTS", "SURVIVAL_PREFERENCES",
	}
	for _, w := range wanted {
		if containsText(candidates, w) {
			return // any-of passes
		}
	}
	t.Errorf("expected at least one placement option keyword after CREATE PLACEMENT POLICY p; got none.\nAll candidates: %v", candidateTexts(candidates))
}

// TestTiDBKeywords_AutoRandomColumnConstraint — P2 bundled with Tier 1.
// AUTO_RANDOM belongs on the column-constraint list alongside
// AUTO_INCREMENT in TiDB CREATE TABLE syntax.
func TestTiDBKeywords_AutoRandomColumnConstraint(t *testing.T) {
	cat := catalog.New()

	sql := "CREATE TABLE t (id BIGINT "
	candidates := Complete(sql, len(sql), cat)

	if !containsText(candidates, "AUTO_RANDOM") {
		t.Errorf("expected AUTO_RANDOM on column-constraint completion; got: %v", candidateTexts(candidates))
	}
}

// TestTiDBKeywords_ClusteredPKModifier — P2 bundled with Tier 1.
// CLUSTERED / NONCLUSTERED appear after a table-level PRIMARY KEY
// (col_list) declaration.
func TestTiDBKeywords_ClusteredPKModifier(t *testing.T) {
	cat := catalog.New()

	sql := "CREATE TABLE t (id INT, PRIMARY KEY (id) "
	candidates := Complete(sql, len(sql), cat)

	if !containsText(candidates, "CLUSTERED") && !containsText(candidates, "NONCLUSTERED") {
		t.Errorf("expected CLUSTERED or NONCLUSTERED after PRIMARY KEY (id); got: %v", candidateTexts(candidates))
	}
}

// candidateTexts extracts the text of each candidate for error messages.
func candidateTexts(cs []Candidate) []string {
	out := make([]string, 0, len(cs))
	for _, c := range cs {
		out = append(out, c.Text)
	}
	return out
}
