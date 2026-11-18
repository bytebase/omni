package driver

import (
	"strings"
	"testing"
)

func TestDMLStatement(t *testing.T) {
	for _, tc := range []struct {
		name      string
		statement string
		wantErr   string
	}{
		{name: "insert", statement: "INSERT INTO t VALUES (1)"},
		{name: "replace", statement: "REPLACE INTO t VALUES (1)"},
		{name: "update", statement: "UPDATE t SET c = 1 WHERE id = 1"},
		{name: "delete", statement: "DELETE FROM t WHERE id = 1"},
		{name: "select", statement: "SELECT * FROM t", wantErr: "expected an INSERT, UPDATE or DELETE"},
		{name: "script", statement: "DELETE FROM t WHERE id = 1; SELECT 1", wantErr: "expected exactly one statement"},
		{name: "empty", statement: "", wantErr: "expected exactly one statement"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			stmt, err := dmlStatement(tc.statement)
			if tc.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("dmlStatement(%q) = (%T, %v), want error containing %q", tc.statement, stmt, err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("dmlStatement(%q) = %v", tc.statement, err)
			}
			if stmt == nil {
				t.Fatalf("dmlStatement(%q) returned no statement", tc.statement)
			}
		})
	}
}

// The plans are EXPLAIN output captured from TiDB v8.5.0.
func TestGetAffectedRowsFromPlan(t *testing.T) {
	for _, tc := range []struct {
		name    string
		plan    []planRow
		want    int64
		wantErr bool
	}{
		{
			name: "join_estimates_joined_rows_not_its_inputs",
			plan: []planRow{
				{"Update_8", "N/A"},
				{"└─HashJoin_27", "800.00"},
				{"  ├─TableReader_44(Build)", "0.40"},
				{"  │ └─Selection_43", "0.40"},
				{"  │   └─TableFullScan_42", "400.00"},
				{"  └─TableReader_38(Probe)", "1000.00"},
				{"    └─Selection_37", "1000.00"},
				{"      └─TableFullScan_36", "1000.00"},
			},
			want: 800,
		},
		{
			name: "limit_is_part_of_the_plan",
			plan: []planRow{
				{"Update_6", "N/A"},
				{"└─TopN_9", "20.00"},
				{"  └─TableReader_17", "20.00"},
				{"    └─TopN_16", "20.00"},
				{"      └─Selection_15", "500.00"},
				{"        └─TableFullScan_14", "1000.00"},
			},
			want: 20,
		},
		{
			name: "foreign_key_cascade_follows_the_input",
			plan: []planRow{
				{"Delete_3", "N/A"},
				{"├─TableReader_8", "9.00"},
				{"│ └─TableRangeScan_7", "9.00"},
				{"└─Foreign_Key_Cascade_4", "0.00"},
			},
			want: 9,
		},
		{
			name: "insert_set_has_no_input",
			plan: []planRow{
				{"Insert_1", "N/A"},
			},
			want: 0,
		},
		{
			// With tidb_opt_enable_non_eval_scalar_subquery, a subquery in SET is
			// planned as a second root rather than evaluated during planning.
			name: "subquery_root_is_not_an_input",
			plan: []planRow{
				{"Insert_1", "N/A"},
				{"ScalarSubQuery_36", "N/A"},
				{"└─MaxOneRow_11", "1.00"},
				{"  └─Projection_12", "1.00"},
				{"    └─StreamAgg_14", "1.00"},
				{"      └─Limit_30", "1.00"},
				{"        └─TableReader_35", "1.00"},
				{"          └─Limit_34", "1.00"},
				{"            └─TableFullScan_24", "1.00"},
			},
			want: 0,
		},
		{
			name: "select_root_is_rejected",
			plan: []planRow{
				{"TableReader_7", "100.00"},
				{"└─Selection_6", "100.00"},
				{"  └─TableFullScan_5", "1000.00"},
			},
			wantErr: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			est, err := getAffectedRowsFromPlan(tc.plan)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("getAffectedRowsFromPlan = %v, want an error", est)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if got := roundRows(est); got != tc.want {
				t.Errorf("getAffectedRowsFromPlan = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestGetAffectedRows(t *testing.T) {
	stmt, err := dmlStatement("DELETE a, b FROM a JOIN b ON a.id = b.id")
	if err != nil {
		t.Fatal(err)
	}
	plan := []planRow{
		{"Delete_7", "N/A"},
		{"└─HashJoin_20", "0.40"},
	}
	// 0.4 joined rows change a row in each of two tables.
	got, err := getAffectedRows(plan, stmt)
	if err != nil {
		t.Fatal(err)
	}
	if got != 1 {
		t.Errorf("getAffectedRows = %d, want 1", got)
	}
}

func TestCountDMLTargets(t *testing.T) {
	for _, tc := range []struct {
		statement string
		want      int
	}{
		{statement: "UPDATE t SET a = 1, b = 2 WHERE id = 1", want: 1},
		{statement: "UPDATE a JOIN b ON a.id = b.id SET a.c = 1, a.d = 2", want: 1},
		{statement: "UPDATE a JOIN b ON a.id = b.id SET a.c = 1, b.c = 1", want: 2},
		{statement: "UPDATE a AS x, b AS y SET x.c = 1, y.c = 1 WHERE x.id = y.id", want: 2},
		// Either joined table may own each unqualified column, up to the two joined tables.
		{statement: "UPDATE a JOIN b ON a.id = b.id SET c = 1, d = 2, e = 3", want: 2},
		{statement: "DELETE a, b FROM a JOIN b ON a.id = b.id", want: 2},
		{statement: "DELETE a FROM a JOIN b ON a.id = b.id", want: 1},
		{statement: "DELETE FROM t WHERE id = 1", want: 1},
		{statement: "INSERT INTO t SELECT * FROM s", want: 1},
	} {
		t.Run(tc.statement, func(t *testing.T) {
			stmt, err := dmlStatement(tc.statement)
			if err != nil {
				t.Fatal(err)
			}
			if got := countDMLTargets(stmt); got != tc.want {
				t.Errorf("countDMLTargets = %d, want %d", got, tc.want)
			}
		})
	}
}
