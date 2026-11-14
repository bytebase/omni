package parser

import (
	"testing"
)

// TestParseAnalyticObjectStmts tests parsing of CREATE/ALTER ANALYTIC VIEW,
// CREATE/ALTER ATTRIBUTE DIMENSION, and CREATE/ALTER HIERARCHY statements.
func TestParseAnalyticObjectStmts(t *testing.T) {
	tests := []string{
		// ---- CREATE ANALYTIC VIEW ----
		// Minimal
		`CREATE ANALYTIC VIEW sales_av
			USING sales
			DIMENSION BY (time_dim KEY month_id REFERENCES month HIERARCHIES (time_dim_hier))
			MEASURES (amount)`,
		// OR REPLACE
		`CREATE OR REPLACE ANALYTIC VIEW sales_av
			USING sales
			DIMENSION BY (geo KEY region_id REFERENCES region HIERARCHIES (geo_hier))
			MEASURES (total_sales)`,
		// IF NOT EXISTS
		`CREATE ANALYTIC VIEW IF NOT EXISTS sales_av
			USING sales
			DIMENSION BY (time KEY id REFERENCES day HIERARCHIES (time_hier))
			MEASURES (revenue)`,
		// SHARING
		`CREATE ANALYTIC VIEW sales_av
			SHARING = METADATA
			USING sales
			DIMENSION BY (time KEY id REFERENCES day HIERARCHIES (time_hier))
			MEASURES (revenue)`,
		`CREATE ANALYTIC VIEW sales_av
			SHARING = NONE
			USING sales
			DIMENSION BY (time KEY id REFERENCES day HIERARCHIES (time_hier))
			MEASURES (revenue)`,
		// With alias
		`CREATE ANALYTIC VIEW sales_av
			USING sales AS s
			DIMENSION BY (time KEY month_id REFERENCES month HIERARCHIES (time_hier))
			MEASURES (total FACT amount AGGREGATE BY SUM)`,
		// DEFAULT MEASURE
		`CREATE ANALYTIC VIEW sales_av
			USING sales
			DIMENSION BY (time KEY id REFERENCES day HIERARCHIES (time_hier))
			MEASURES (amount, quantity)
			DEFAULT MEASURE amount`,
		// DEFAULT AGGREGATE BY
		`CREATE ANALYTIC VIEW sales_av
			USING sales
			DIMENSION BY (time KEY id REFERENCES day HIERARCHIES (time_hier))
			MEASURES (amount)
			DEFAULT AGGREGATE BY SUM`,
		// Schema-qualified
		`CREATE ANALYTIC VIEW hr.emp_av
			USING hr.employees
			DIMENSION BY (dept KEY department_id REFERENCES department HIERARCHIES (dept_hier))
			MEASURES (salary)`,

		// ---- ALTER ANALYTIC VIEW ----
		// RENAME
		`ALTER ANALYTIC VIEW sales_av RENAME TO revenue_av`,
		// COMPILE
		`ALTER ANALYTIC VIEW sales_av COMPILE`,
		// IF EXISTS
		`ALTER ANALYTIC VIEW IF EXISTS sales_av COMPILE`,
		// Schema-qualified
		`ALTER ANALYTIC VIEW hr.emp_av COMPILE`,
		// ADD CACHE
		`ALTER ANALYTIC VIEW sales_av ADD CACHE
			MEASURE GROUP (amount)
			LEVELS (month, quarter) MATERIALIZED`,
		// DROP CACHE
		`ALTER ANALYTIC VIEW sales_av DROP CACHE
			MEASURE GROUP (amount)
			LEVELS (month)`,

		// ---- CREATE ATTRIBUTE DIMENSION ----
		// Minimal
		`CREATE ATTRIBUTE DIMENSION time_attr_dim
			USING time_table
			ATTRIBUTES (month_id, month_name, year_id)
			LEVEL month
				KEY month_id
				ORDER BY month_id`,
		// OR REPLACE
		`CREATE OR REPLACE ATTRIBUTE DIMENSION time_attr_dim
			USING time_table
			ATTRIBUTES (month_id, year_id)
			LEVEL month KEY month_id ORDER BY month_id`,
		// IF NOT EXISTS
		`CREATE ATTRIBUTE DIMENSION IF NOT EXISTS time_attr_dim
			USING time_table
			ATTRIBUTES (day_id)
			LEVEL day KEY day_id ORDER BY day_id`,
		// SHARING
		`CREATE ATTRIBUTE DIMENSION time_attr_dim
			SHARING = METADATA
			USING time_table
			ATTRIBUTES (day_id)
			LEVEL day KEY day_id ORDER BY day_id`,
		// DIMENSION TYPE TIME
		`CREATE ATTRIBUTE DIMENSION time_attr_dim
			DIMENSION TYPE TIME
			USING time_table
			ATTRIBUTES (day_id, month_id, year_id)
			LEVEL day KEY day_id ORDER BY day_id
			LEVEL month KEY month_id ORDER BY month_id
			LEVEL year KEY year_id ORDER BY year_id`,
		// DIMENSION TYPE STANDARD
		`CREATE ATTRIBUTE DIMENSION geo_attr_dim
			DIMENSION TYPE STANDARD
			USING geo_table
			ATTRIBUTES (city_id, state_id)
			LEVEL city KEY city_id ORDER BY city_id
			LEVEL state KEY state_id ORDER BY state_id`,
		// Multiple LEVEL clauses with DETERMINES
		`CREATE ATTRIBUTE DIMENSION time_attr_dim
			USING time_table
			ATTRIBUTES (day_id, day_name, month_id, month_name, year_id)
			LEVEL day
				KEY day_id
				ORDER BY day_id
				DETERMINES (day_name)
			LEVEL month
				KEY month_id
				ORDER BY month_id
				DETERMINES (month_name)
			LEVEL year
				KEY year_id
				ORDER BY year_id`,
		// Schema-qualified with multiple sources
		`CREATE ATTRIBUTE DIMENSION hr.geo_dim
			USING hr.cities, hr.states
			ATTRIBUTES (city_id, state_id, country_id)
			LEVEL city KEY city_id ORDER BY city_id
			LEVEL state KEY state_id ORDER BY state_id`,
		// ALL clause
		`CREATE ATTRIBUTE DIMENSION time_attr_dim
			USING time_table
			ATTRIBUTES (day_id)
			LEVEL day KEY day_id ORDER BY day_id
			ALL MEMBER NAME 'All Time'`,

		// ---- ALTER ATTRIBUTE DIMENSION ----
		// RENAME
		`ALTER ATTRIBUTE DIMENSION time_attr_dim RENAME TO calendar_dim`,
		// COMPILE
		`ALTER ATTRIBUTE DIMENSION time_attr_dim COMPILE`,
		// IF EXISTS
		`ALTER ATTRIBUTE DIMENSION IF EXISTS time_attr_dim COMPILE`,
		// Schema-qualified
		`ALTER ATTRIBUTE DIMENSION hr.geo_dim COMPILE`,

		// ---- CREATE HIERARCHY ----
		// Minimal
		`CREATE HIERARCHY time_hier
			USING time_attr_dim
			(day CHILD OF month CHILD OF year)`,
		// OR REPLACE
		`CREATE OR REPLACE HIERARCHY time_hier
			USING time_attr_dim
			(day CHILD OF month CHILD OF year)`,
		// IF NOT EXISTS
		`CREATE HIERARCHY IF NOT EXISTS time_hier
			USING time_attr_dim
			(day CHILD OF month CHILD OF year)`,
		// SHARING
		`CREATE HIERARCHY time_hier
			SHARING = METADATA
			USING time_attr_dim
			(day CHILD OF month CHILD OF year)`,
		// HIERARCHICAL ATTRIBUTES
		`CREATE HIERARCHY time_hier
			USING time_attr_dim
			(day CHILD OF month CHILD OF year)
			HIERARCHICAL ATTRIBUTES (
				HIER_ORDER,
				DEPTH,
				IS_LEAF
			)`,
		// Schema-qualified
		`CREATE HIERARCHY hr.org_hier
			USING hr.org_dim
			(employee CHILD OF department CHILD OF division)`,
		// Deep hierarchy
		`CREATE HIERARCHY time_hier
			USING time_attr_dim
			(day CHILD OF week CHILD OF month CHILD OF quarter CHILD OF year)`,

		// ---- ALTER HIERARCHY ----
		// RENAME
		`ALTER HIERARCHY time_hier RENAME TO calendar_hier`,
		// COMPILE
		`ALTER HIERARCHY time_hier COMPILE`,
		// IF EXISTS
		`ALTER HIERARCHY IF EXISTS time_hier COMPILE`,
		// Schema-qualified
		`ALTER HIERARCHY hr.org_hier COMPILE`,
	}
	for _, sql := range tests {
		name := sql
		if len(name) > 60 {
			name = name[:60]
		}
		t.Run(name, func(t *testing.T) {
			result := ParseAndCheck(t, sql)
			if result.Len() < 1 {
				t.Fatalf("expected at least 1 statement, got %d", result.Len())
			}
		})
	}
}
