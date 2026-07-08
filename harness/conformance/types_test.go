package main

import "testing"

func TestClassifyFamily(t *testing.T) {
	cases := []struct{ sql, want string }{
		{"SELECT 1", "SELECT"},
		{"  select * from t", "SELECT"},
		{"CREATE TABLE t (a int)", "CREATE TABLE"},
		{"CREATE UNIQUE INDEX i ON t(a)", "CREATE INDEX"},
		{"ALTER TABLE t ADD COLUMN b int", "ALTER TABLE"},
		{"INSERT INTO t VALUES (1)", "INSERT"},
		{"WITH c AS (SELECT 1) SELECT * FROM c", "SELECT"},
		{"(SELECT 1) UNION (SELECT 2)", "SELECT"},
		{"ADMIN CHECK TABLE t", "ADMIN"},
		{"/* hint */ SELECT 1", "SELECT"},
		{"", "UNKNOWN"},

		// Prefix-table order: the specific CREATE arms must win over the bare
		// CREATE catch-all, which only works because CREATE sorts last.
		{"CREATE PROCEDURE p() SELECT 1", "CREATE OTHER"},
		{"CREATE VIEW v AS SELECT 1", "CREATE VIEW"},
		{"CREATE USER u IDENTIFIED BY 'x'", "DCL"},

		{"DESC t", "EXPLAIN"},
		{"DESCRIBE t", "EXPLAIN"},
		{"SET NAMES utf8mb4", "SET"},
		{"USE test", "SET"},
		{"TABLE t", "SELECT"},
		{"VALUES ROW(1)", "SELECT"},
		{"BEGIN", "TXN"},
		{"GRANT SELECT ON *.* TO u", "DCL"},
		{"UPDATE t SET a = 1", "UPDATE"},
		{"XA COMMIT 'x'", "OTHER"},

		{"-- note\nSELECT 1", "SELECT"},
		{"# note\nSELECT 1", "SELECT"},
		{"/*c1*/ /*c2*/ SELECT 1", "SELECT"},
		{"/*+ SET_VAR(sql_mode='') */ INSERT INTO t VALUES (1)", "INSERT"},
	}
	for _, c := range cases {
		if got := classifyFamily(c.sql); got != c.want {
			t.Errorf("classifyFamily(%q) = %q, want %q", c.sql, got, c.want)
		}
	}
}

func TestClusterKey(t *testing.T) {
	// Positions, numbers, and quoted identifiers must normalize away so one
	// grammar divergence = one cluster.
	a := clusterKey(`syntax error at line 1 column 27 near "FOO"`)
	b := clusterKey(`syntax error at line 3 column 9 near "BAR"`)
	if a != b {
		t.Errorf("cluster keys differ: %q vs %q", a, b)
	}
	if want := "syntax error at line N column N near ?"; a != want {
		t.Errorf("clusterKey normal form = %q, want %q", a, want)
	}
	if got, want := clusterKey("unknown column `x1` in `t23`"), "unknown column ? in ?"; got != want {
		t.Errorf("clusterKey backtick form = %q, want %q", got, want)
	}
	if got, want := clusterKey("  a  b\t c  "), "a b c"; got != want {
		t.Errorf("clusterKey whitespace collapse = %q, want %q", got, want)
	}

	// omni's real Parse error format appends "\nrelated text: <raw source
	// line>" (tidb/parser/parser.go); the related text never normalizes, so
	// only the first line may contribute to the key.
	a = clusterKey("syntax error at or near \"LIKEX\" (line 1, column 25)\nrelated text: SELECT 1 FROM t WHERE a LIKEX 'p'")
	b = clusterKey("syntax error at or near \"LIKEX\" (line 1, column 26)\nrelated text: SELECT 2 FROM u WHERE zz LIKEX 'q'")
	if a != b {
		t.Errorf("multi-line cluster keys differ: %q vs %q", a, b)
	}
	if want := "syntax error at or near ? (line N, column N)"; a != want {
		t.Errorf("multi-line clusterKey = %q, want %q", a, want)
	}
}

func TestStmtHash_Stable(t *testing.T) {
	if stmtHash("SELECT 1") != stmtHash("SELECT 1") {
		t.Error("hash not deterministic")
	}
	if stmtHash("SELECT 1") == stmtHash("SELECT 2") {
		t.Error("hash collision on different SQL")
	}
}

func TestClassify(t *testing.T) {
	// wantReason and wantKey are asserted exactly: an empty want means the
	// classifier must leave the field empty.
	cases := []struct {
		name       string
		row        Row
		wantClass  Class
		wantReason string
		wantKey    string
	}{
		{
			name:      "label only, both accept",
			row:       Row{SQL: "SELECT 1", Expected: VerdictAccept, OmniVerdict: VerdictAccept},
			wantClass: ClassAgreeAccept,
		},
		{
			name:      "label only, both reject",
			row:       Row{SQL: "SELECT SELECT", Expected: VerdictReject, OmniVerdict: VerdictReject},
			wantClass: ClassAgreeReject,
		},
		{
			name: "label accept, omni reject is GAP keyed on the omni error",
			row: Row{
				SQL:         "SELECT 1",
				Expected:    VerdictAccept,
				OmniVerdict: VerdictReject,
				OmniError:   `syntax error at line 1 column 8 near "1"`,
			},
			wantClass: ClassGap,
			wantKey:   "syntax error at line N column N near ?",
		},
		{
			name: "label reject, omni accept is OVER keyed on leading tokens before adjudication",
			row: Row{
				SQL:         "ALTER TABLE t BOGUS CLAUSE here",
				Expected:    VerdictReject,
				OmniVerdict: VerdictAccept,
			},
			wantClass: ClassOver,
			wantKey:   "ALTER TABLE T BOGUS",
		},
		{
			name: "OVER leading-token key keeps short statements whole",
			row: Row{
				SQL:         "BOGUS TOKEN",
				Expected:    VerdictReject,
				OmniVerdict: VerdictAccept,
			},
			wantClass: ClassOver,
			wantKey:   "BOGUS TOKEN",
		},
		{
			name: "OVER leading-token key normalizes numbered identifiers",
			row: Row{
				SQL:         "ALTER TABLE t42 BOGUS x",
				Expected:    VerdictReject,
				OmniVerdict: VerdictAccept,
			},
			wantClass: ClassOver,
			wantKey:   "ALTER TABLE TN BOGUS",
		},
		{
			name: "adjudicated OVER upgrades the key to the engine message",
			row: Row{
				SQL:             "ALTER TABLE t BOGUS",
				EngineVerdict:   VerdictReject,
				OmniVerdict:     VerdictAccept,
				RawErrorMessage: "You have an error in your SQL syntax near 'BOGUS' at line 1",
			},
			wantClass: ClassOver,
			wantKey:   "You have an error in your SQL syntax near ? at line N",
		},
		{
			name: "container supplies ground truth when the label is absent",
			row: Row{
				SQL:           "SELECT 1",
				EngineVerdict: VerdictAccept,
				OmniVerdict:   VerdictReject,
				OmniError:     "unexpected token FOO",
			},
			wantClass: ClassGap,
			wantKey:   "unexpected token FOO",
		},
		{
			name: "container agreeing with the label classifies normally",
			row: Row{
				SQL:           "SELECT 1",
				Expected:      VerdictAccept,
				EngineVerdict: VerdictAccept,
				OmniVerdict:   VerdictReject,
				OmniError:     "unexpected token FOO",
			},
			wantClass: ClassGap,
			wantKey:   "unexpected token FOO",
		},
		{
			name: "label vs container disagreement is INDETERMINATE, never silently trusted",
			row: Row{
				SQL:           "SELECT 1",
				Expected:      VerdictReject,
				EngineVerdict: VerdictAccept,
				OmniVerdict:   VerdictReject,
			},
			wantClass:  ClassIndeterminate,
			wantReason: "label_container_disagree",
		},
		{
			name:       "no ground truth at all is INDETERMINATE",
			row:        Row{SQL: "SELECT 1", OmniVerdict: VerdictAccept},
			wantClass:  ClassIndeterminate,
			wantReason: "no_ground_truth",
		},
		{
			name:       "missing omni verdict is INDETERMINATE, not OVER",
			row:        Row{SQL: "SELECT 1", Expected: VerdictReject},
			wantClass:  ClassIndeterminate,
			wantReason: "no_omni_verdict",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := c.row
			classify(&r)
			if r.Class != c.wantClass {
				t.Fatalf("class = %q, want %q", r.Class, c.wantClass)
			}
			if r.ClassifierReason != c.wantReason {
				t.Errorf("classifier reason = %q, want %q", r.ClassifierReason, c.wantReason)
			}
			if r.DivergenceKey != c.wantKey {
				t.Errorf("divergence key = %q, want %q", r.DivergenceKey, c.wantKey)
			}
		})
	}
}

func TestClassify_Reclassification(t *testing.T) {
	// Rows are re-classified after container adjudication; fields from the
	// first pass must not leak into the second.
	r := Row{SQL: "SELECT 1", OmniVerdict: VerdictReject, OmniError: "unexpected token FOO"}
	classify(&r)
	if r.Class != ClassIndeterminate || r.ClassifierReason != "no_ground_truth" {
		t.Fatalf("first pass = %q/%q, want INDETERMINATE/no_ground_truth", r.Class, r.ClassifierReason)
	}
	r.EngineVerdict = VerdictAccept
	classify(&r)
	if r.Class != ClassGap {
		t.Fatalf("second pass class = %q, want %q", r.Class, ClassGap)
	}
	if r.ClassifierReason != "" {
		t.Errorf("stale classifier reason survived re-classification: %q", r.ClassifierReason)
	}
	if r.DivergenceKey != "unexpected token FOO" {
		t.Errorf("divergence key = %q, want %q", r.DivergenceKey, "unexpected token FOO")
	}

	// The reverse leak: an OVER divergence key must clear when adjudication
	// turns the row INDETERMINATE.
	r2 := Row{SQL: "ALTER TABLE t BOGUS", Expected: VerdictReject, OmniVerdict: VerdictAccept}
	classify(&r2)
	if r2.Class != ClassOver || r2.DivergenceKey == "" {
		t.Fatalf("first pass = %q key %q, want OVER with a key", r2.Class, r2.DivergenceKey)
	}
	r2.EngineVerdict = VerdictAccept
	classify(&r2)
	if r2.Class != ClassIndeterminate || r2.ClassifierReason != "label_container_disagree" {
		t.Fatalf("second pass = %q/%q, want INDETERMINATE/label_container_disagree", r2.Class, r2.ClassifierReason)
	}
	if r2.DivergenceKey != "" {
		t.Errorf("stale divergence key survived re-classification: %q", r2.DivergenceKey)
	}
}
