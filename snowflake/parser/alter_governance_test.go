package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// Data metric function actions on ALTER TABLE
// (https://docs.snowflake.com/en/sql-reference/sql/alter-table).
func TestAlterTable_DataMetricFunction(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t ADD DATA METRIC FUNCTION snowflake.core.null_count ON (c1, c2, TABLE r(k)) "+
		"WITHIN GROUP (g) GROUP LIMIT 10 FILTER (c1 > 0) EXPECTATION e1 (value = 0), e2 (value < 5) "+
		"EXECUTE AS ROLE r1 ANOMALY_DETECTION = TRUE SENSITIVITY = 'LOW', m2 ON ()")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableDataMetric || a.DataMetric == nil || a.DataMetric.Kind != ast.DataMetricAdd {
		t.Fatalf("got kind %v, data metric %+v", a.Kind, a.DataMetric)
	}
	fns := a.DataMetric.Functions
	if len(fns) != 2 {
		t.Fatalf("functions = %d, want 2", len(fns))
	}
	f := fns[0]
	if len(f.Columns) != 2 || f.RefTable == nil || len(f.RefColumns) != 1 || len(f.WithinGroup) != 1 ||
		f.GroupLimit == nil || *f.GroupLimit != 10 || f.Filter == nil || len(f.Expectations) != 2 ||
		f.ExecuteAs.Name != "r1" || len(f.Options) != 2 {
		t.Errorf("clause = %+v", f)
	}
	if !fns[1].HasOn || len(fns[1].Columns) != 0 {
		t.Errorf("second clause = %+v", fns[1])
	}

	for _, sql := range []string{
		"ALTER TABLE t SET DATA_METRIC_SCHEDULE = 'TRIGGER_ON_CHANGES'",
		"ALTER TABLE t UNSET DATA_METRIC_SCHEDULE",
		"ALTER TABLE t DROP DATA METRIC FUNCTION m ON (a), m2 ON (b)",
		"ALTER TABLE t MODIFY DATA METRIC FUNCTION m ON (a) SUSPEND, m ON (b) RESUME",
		"ALTER TABLE t MODIFY DATA METRIC FUNCTION m ON (a) ADD EXPECTATION e (value > 1), e2 (value < 9)",
		"ALTER TABLE t MODIFY DATA METRIC FUNCTION m ON (a) DROP EXPECTATION e, e2",
		"ALTER TABLE t MODIFY DATA METRIC FUNCTION m ON (a) FILTER (), m ON (b) FILTER (b > 1)",
		"ALTER TABLE t MODIFY DATA METRIC FUNCTION m SET ANOMALY_DETECTION = TRUE SENSITIVITY = 'HIGH'",
	} {
		requireAlterTable(t, sql)
	}

	// DATA stays a column name when METRIC does not follow.
	stmt = requireAlterTable(t, "ALTER TABLE t ALTER data SET NOT NULL")
	if a := requireAction(t, stmt, 0); a.Kind != ast.AlterTableAlterColumn {
		t.Errorf("kind = %v, want AlterTableAlterColumn", a.Kind)
	}

	for _, sql := range []string{
		"ALTER TABLE t ADD DATA METRIC FUNCTION m",                  // ON is required
		"ALTER TABLE t ADD DATA METRIC FUNCTION m ON (a) FILTER ()", // ADD's filter needs an expression
		"ALTER TABLE t MODIFY DATA METRIC FUNCTION m ON (a)",        // MODIFY needs a change
		"ALTER TABLE t MODIFY DATA METRIC FUNCTION m SET ANOMALY_DETECTION",
		"ALTER TABLE t SET DATA_METRIC_SCHEDULE = 5",
		"ALTER TABLE t ADD DATA METRIC FUNCTION m ON (a,)",           // no trailing comma
		"ALTER TABLE t MODIFY DATA METRIC FUNCTION m SET ANY = TRUE", // only the documented properties
		"ALTER TABLE t MODIFY DATA METRIC FUNCTION m SET SENSITIVITY",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) succeeded, want error", sql)
		}
	}
}

func TestAlterTable_UnparenthesizedTags(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t SET TAG a.b.t1 = 'x', t2 = 'y'")
	if a := requireAction(t, stmt, 0); a.Kind != ast.AlterTableSetTag || len(a.Tags) != 2 {
		t.Errorf("action = %+v", a)
	}
	stmt = requireAlterTable(t, "ALTER TABLE t UNSET TAG t1, t2")
	if a := requireAction(t, stmt, 0); a.Kind != ast.AlterTableUnsetTag || len(a.UnsetTags) != 2 {
		t.Errorf("action = %+v", a)
	}

	// A column list continues without repeating ALTER: each column is an action.
	stmt = requireAlterTable(t, "ALTER TABLE t ALTER c1 SET TAG t1 = 'x', t2 = 'y', COLUMN c2 SET TAG t1 = 'z'")
	if len(stmt.Actions) != 2 || len(stmt.Actions[0].Tags) != 2 || stmt.Actions[1].TagColumn.Name != "c2" {
		t.Errorf("actions = %+v", stmt.Actions)
	}
	stmt = requireAlterTable(t, "ALTER TABLE t MODIFY c1 UNSET TAG t1, COLUMN c2 UNSET TAG t2")
	if len(stmt.Actions) != 2 || len(stmt.Actions[0].UnsetTags) != 1 || stmt.Actions[1].TagColumn.Name != "c2" {
		t.Errorf("actions = %+v", stmt.Actions)
	}
	stmt = requireAlterTable(t, "ALTER TABLE t MODIFY c1 SET TAG t1 = 'x', COLUMN c2 SET TAG t2 = 'y'")
	if len(stmt.Actions) != 2 || stmt.Actions[1].TagColumn.Name != "c2" {
		t.Errorf("actions = %+v", stmt.Actions)
	}
	stmt = requireAlterTable(t, `ALTER TABLE t MODIFY c1 SET TAG t1 = 'v', "c2" SET TAG t2 = 'w'`)
	if len(stmt.Actions) != 2 || stmt.Actions[1].TagColumn.Name != "c2" {
		t.Errorf("actions = %+v", stmt.Actions)
	}
	stmt = requireAlterTable(t, `ALTER TABLE t MODIFY c1 UNSET TAG t1, "c2" UNSET TAG t2`)
	if len(stmt.Actions) != 2 || stmt.Actions[1].TagColumn.Name != "c2" {
		t.Errorf("actions = %+v", stmt.Actions)
	}
	stmt = requireAlterTable(t, "ALTER TABLE t MODIFY c1 UNSET TAG t1, t2, c2 UNSET TAG t3")
	if len(stmt.Actions) != 2 || len(stmt.Actions[0].UnsetTags) != 2 || stmt.Actions[1].TagColumn.Name != "c2" {
		t.Errorf("actions = %+v", stmt.Actions)
	}
	// A plain column change followed by a tag change on another column.
	stmt = requireAlterTable(t, "ALTER TABLE t ALTER c1 SET NOT NULL, c2 SET TAG t1 = 'z'")
	if len(stmt.Actions) != 2 || stmt.Actions[1].Kind != ast.AlterTableSetColumnTag {
		t.Errorf("actions = %+v", stmt.Actions)
	}
}

func TestAlterTable_Policies(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t SET AGGREGATION POLICY p ENTITY KEY (a, b) FORCE")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableSetAggregationPolicy || len(a.PolicyKeyCols) != 2 || !a.PolicyForce {
		t.Errorf("action = %+v", a)
	}
	stmt = requireAlterTable(t, "ALTER TABLE t ALTER COLUMN c SET PROJECTION POLICY p FORCE")
	a = requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableSetProjectionPolicy || a.ProjectionColumn.Name != "c" || !a.PolicyForce {
		t.Errorf("action = %+v", a)
	}
	for sql, kind := range map[string]ast.AlterTableActionKind{
		"ALTER TABLE t UNSET AGGREGATION POLICY":                 ast.AlterTableUnsetAggregationPolicy,
		"ALTER TABLE t SET JOIN POLICY p FORCE":                  ast.AlterTableSetJoinPolicy,
		"ALTER TABLE t UNSET JOIN POLICY":                        ast.AlterTableUnsetJoinPolicy,
		"ALTER TABLE t MODIFY c UNSET PROJECTION POLICY":         ast.AlterTableUnsetProjectionPolicy,
		"ALTER TABLE t ADD STORAGE LIFECYCLE POLICY p ON (a, b)": ast.AlterTableAddStorageLifecyclePolicy,
		"ALTER TABLE t DROP STORAGE LIFECYCLE POLICY":            ast.AlterTableDropStorageLifecyclePolicy,
	} {
		if a := requireAction(t, requireAlterTable(t, sql), 0); a.Kind != kind {
			t.Errorf("%q: kind = %v, want %v", sql, a.Kind, kind)
		}
	}
}

func TestAlterTable_AlterConstraint(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t ALTER CONSTRAINT c1 (a, b) NOT ENFORCED NOVALIDATE RELY")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableAlterConstraint || a.ConstraintName.Name != "c1" || len(a.ConstraintCols) != 2 ||
		len(a.ConstraintProps) != 3 || a.ConstraintProps[0] != "NOT ENFORCED" {
		t.Errorf("action = %+v", a)
	}
	for _, sql := range []string{
		"ALTER TABLE t MODIFY UNIQUE (a) ENFORCED",
		"ALTER TABLE t ALTER PRIMARY KEY (a) NORELY",
		"ALTER TABLE t ALTER FOREIGN KEY (a, b)",
	} {
		if a := requireAction(t, requireAlterTable(t, sql), 0); a.Kind != ast.AlterTableAlterConstraint {
			t.Errorf("%q: kind = %v", sql, a.Kind)
		}
	}
	// PRIMARY and FOREIGN are not reserved: without KEY they are columns.
	if a := requireAction(t, requireAlterTable(t, "ALTER TABLE t ALTER primary SET NOT NULL"), 0); a.Kind != ast.AlterTableAlterColumn {
		t.Errorf("kind = %v, want AlterTableAlterColumn", a.Kind)
	}
	if _, err := Parse("ALTER TABLE t ALTER UNIQUE (a) RELY ENFORCED"); err == nil {
		t.Error("properties out of the documented order should not parse")
	}
}

func TestAlterTable_NewProperties(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t SET DATA_RETENTION_TIME_IN_DAYS = 1 ICEBERG_DEFAULT_DDL_COLLATION = 'en' "+
		"ENABLE_SCHEMA_EVOLUTION = TRUE ERROR_LOGGING = FALSE CONTACT steward = alice, support = bob COMMENT = 'c' ROW_TIMESTAMP = TRUE")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableSet || len(a.Props) != 7 {
		t.Fatalf("action = %+v", a)
	}
	if a.Props[4].Name != "CONTACT" || a.Props[4].Value != "steward = alice, support = bob" {
		t.Errorf("contact = %+v", a.Props[4])
	}
	if a.Props[1].Value != "'en'" {
		t.Errorf("collation value = %q", a.Props[1].Value)
	}
	stmt = requireAlterTable(t, "ALTER TABLE t UNSET CONTACT steward, ERROR_LOGGING, ROW_TIMESTAMP")
	if a := requireAction(t, stmt, 0); len(a.UnsetProps) != 3 || a.UnsetProps[0] != "CONTACT steward" {
		t.Errorf("unset = %v", a.UnsetProps)
	}
	for _, sql := range []string{
		"ALTER TABLE t SET",            // nothing to set
		"ALTER TABLE t SET ANY = TRUE", // a reserved word is no property name
		"ALTER TABLE t SET NOT_A_PROPERTY = 12",
		"ALTER TABLE t UNSET NOT_A_PROPERTY",
		"ALTER TABLE t UNSET COMMENT, PATH_LAYOUT", // an ALTER ICEBERG TABLE property
		"ALTER TABLE t SET PATH_LAYOUT = FLAT",     // an ALTER ICEBERG TABLE property
		"ALTER TABLE t SET ERROR_LOGGING",          // a property needs = value
		"ALTER TABLE t SET ERROR_LOGGING = 1 2",    // one value per property
		// Documented properties take their documented values.
		"ALTER TABLE t SET ENABLE_SCHEMA_EVOLUTION = 'garbage'",
		"ALTER TABLE t SET ERROR_LOGGING = 1",
		"ALTER TABLE t SET ROW_TIMESTAMP = (a = b)",
		"ALTER TABLE t SET ICEBERG_DEFAULT_DDL_COLLATION = en",
		"ALTER TABLE t MODIFY DATA METRIC FUNCTION m SET SENSITIVITY = 'EXTREME'",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) succeeded, want error", sql)
		}
	}
}

func TestAlterTable_SearchOptimizationTargets(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t ADD SEARCH OPTIMIZATION ON FULL_TEXT(a, b, ANALYZER => 'UNICODE_ANALYZER')")
	if a := requireAction(t, stmt, 0); len(a.SearchOptOn) != 1 || a.SearchOptOn[0] != "FULL_TEXT(a, b, ANALYZER => 'UNICODE_ANALYZER')" {
		t.Errorf("on = %v", a.SearchOptOn)
	}
	for _, sql := range []string{
		"ALTER TABLE t ADD SEARCH OPTIMIZATION ON FULL_TEXT(*, ANALYZER => 'UNICODE_ANALYZER')",
		"ALTER TABLE t DROP SEARCH OPTIMIZATION ON FULL_TEXT(*, ANALYZER => 'UNICODE_ANALYZER')",
	} {
		if a := requireAction(t, requireAlterTable(t, sql), 0); len(a.SearchOptOn) != 1 || a.SearchOptOn[0] != "FULL_TEXT(*, ANALYZER => 'UNICODE_ANALYZER')" {
			t.Errorf("%q: on = %v", sql, a.SearchOptOn)
		}
	}
	if _, err := Parse("ALTER TABLE t ADD SEARCH OPTIMIZATION ON FULL_TEXT(*, a)"); err == nil {
		t.Error("only ANALYZER may follow *")
	}
	stmt = requireAlterTable(t, "ALTER TABLE t DROP SEARCH OPTIMIZATION ON c1, 1, EQUALITY(c2)")
	if a := requireAction(t, stmt, 0); len(a.SearchOptOn) != 3 || a.SearchOptOn[1] != "1" {
		t.Errorf("on = %v", a.SearchOptOn)
	}
	// Columns named CLUSTERING and SEARCH.
	for _, sql := range []string{"ALTER TABLE t DROP clustering", "ALTER TABLE t DROP search"} {
		if a := requireAction(t, requireAlterTable(t, sql), 0); a.Kind != ast.AlterTableDropColumn {
			t.Errorf("%q: kind = %v", sql, a.Kind)
		}
	}
	if _, err := Parse("ALTER TABLE t ADD SEARCH OPTIMIZATION ON c1"); err == nil {
		t.Error("ADD SEARCH OPTIMIZATION takes a method, not a bare column")
	}
}

func TestLookaheadHelpers(t *testing.T) {
	cases := []struct {
		sql  string
		f    func(p *Parser) bool
		want bool
	}{
		{", t = 'v'", (*Parser).commaStartsTagAssignment, true},
		{", c SET TAG", (*Parser).commaStartsTagAssignment, false},
		{", t2", (*Parser).commaStartsUnsetTagName, true},
		{", c UNSET TAG", (*Parser).commaStartsUnsetTagName, false},
		{", e (1)", (*Parser).commaStartsExpectation, true},
		{", m ON (a)", (*Parser).commaStartsExpectation, false},
		{", m ON (a)", (*Parser).commaStartsDataMetricClause, true},
		{", e (1)", (*Parser).commaStartsDataMetricClause, false},
		{", b = c", (*Parser).commaStartsContactEntry, true},
		{", c2 INT", (*Parser).commaStartsPlainColumnAlter, true},
		{", c2 SET MASKING POLICY p", (*Parser).commaStartsPlainColumnAlter, false},
		{"DATA METRIC", (*Parser).startsDataMetricFunction, true},
		{"DATA SET", (*Parser).startsDataMetricFunction, false},
		{"ICEBERG TABLE", (*Parser).startsAlterIcebergTable, true},
		{"UNIQUE (a)", (*Parser).startsAlterConstraint, true},
		{"PRIMARY SET", (*Parser).startsAlterConstraint, false},
	}
	for _, c := range cases {
		p := &Parser{lexer: NewLexer(c.sql), input: c.sql}
		p.advance()
		before := p.cur
		if got := c.f(p); got != c.want {
			t.Errorf("%q: got %v, want %v", c.sql, got, c.want)
		}
		if p.cur != before {
			t.Errorf("%q: lookahead moved the parser", c.sql)
		}
	}
}
