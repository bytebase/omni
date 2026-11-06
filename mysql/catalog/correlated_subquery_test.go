package catalog

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"

	"github.com/bytebase/omni/mysql/internal/testcontainer"
)

// correlatedDDL and correlatedRows are the database every correlated-subquery
// test reads. Department 3 has no employee and employee 4 has no department, so
// EXISTS, IN, ANY and ALL each answer differently from row to row.
const correlatedDDL = `CREATE TABLE employees (id INT PRIMARY KEY, name VARCHAR(100), salary DECIMAL(10,2), department_id INT);
CREATE TABLE departments (id INT PRIMARY KEY, name VARCHAR(100))`

const correlatedRows = `INSERT INTO departments VALUES (1, 'eng'), (2, 'ops'), (3, 'empty');
INSERT INTO employees VALUES (1, 'ann', 100, 1), (2, 'bob', 200, 1), (3, 'cy', 150, 2), (4, 'eng', 50, NULL)`

// correlatedCatalog returns a catalog holding correlatedDDL's tables in db.
func correlatedCatalog(t *testing.T, db string) *Catalog {
	t.Helper()
	c := New()
	results, err := c.Exec(fmt.Sprintf("CREATE DATABASE %s; USE %s; %s;", db, db, correlatedDDL), nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range results {
		if r.Error != nil {
			t.Fatal(r.Error)
		}
	}
	return c
}

// correlatedCases are subqueries that read a column of an enclosing query, with
// the text DeparseQuery writes for them. A column LevelsUp queries out is
// written with the alias of the query it belongs to, not resolved against the
// subquery's own tables.
var correlatedCases = []struct {
	name, sql, want string
}{
	{
		"in",
		"SELECT e.id, e.salary IN (SELECT d.id FROM departments d WHERE d.id = e.department_id) AS x FROM employees e",
		"select `e`.`id`, `e`.`salary` in (select `d`.`id` from `corr`.`departments` as `d` where (`d`.`id` = `e`.`department_id`)) as `x` from `corr`.`employees` as `e`",
	},
	{
		"exists",
		"SELECT e.id FROM employees e WHERE EXISTS (SELECT 1 FROM departments d WHERE d.id = e.department_id)",
		"select `e`.`id` from `corr`.`employees` as `e` where exists (select 1 as `1` from `corr`.`departments` as `d` where (`d`.`id` = `e`.`department_id`))",
	},
	{
		"scalar",
		"SELECT e.id, (SELECT d.name FROM departments d WHERE d.id = e.department_id) AS dept FROM employees e",
		"select `e`.`id`, (select `d`.`name` from `corr`.`departments` as `d` where (`d`.`id` = `e`.`department_id`)) as `dept` from `corr`.`employees` as `e`",
	},
	{
		"any",
		"SELECT e.id FROM employees e WHERE e.salary > ANY (SELECT o.salary FROM employees o WHERE o.department_id = e.department_id)",
		"select `e`.`id` from `corr`.`employees` as `e` where (`e`.`salary` > any (select `o`.`salary` from `corr`.`employees` as `o` where (`o`.`department_id` = `e`.`department_id`)))",
	},
	{
		"all",
		"SELECT e.id FROM employees e WHERE e.salary >= ALL (SELECT o.salary FROM employees o WHERE o.department_id = e.department_id)",
		"select `e`.`id` from `corr`.`employees` as `e` where (`e`.`salary` >= all (select `o`.`salary` from `corr`.`employees` as `o` where (`o`.`department_id` = `e`.`department_id`)))",
	},
	{
		// The innermost subquery reads both enclosing queries: d one level
		// out and e two.
		"two levels",
		"SELECT e.id FROM employees e WHERE EXISTS (SELECT 1 FROM departments d WHERE d.id = e.department_id AND EXISTS (SELECT 1 FROM employees o WHERE o.department_id = d.id AND o.salary > e.salary))",
		"select `e`.`id` from `corr`.`employees` as `e` where exists (select 1 as `1` from `corr`.`departments` as `d` where ((`d`.`id` = `e`.`department_id`) and exists (select 1 as `1` from `corr`.`employees` as `o` where ((`o`.`department_id` = `d`.`id`) and (`o`.`salary` > `e`.`salary`)))))",
	},
	{
		// An outer column as the subquery's select item keeps its own name,
		// so it needs no alias.
		"outer column as select item",
		"SELECT (SELECT e.name FROM departments d WHERE d.id = e.department_id) AS n FROM employees e",
		"select (select `e`.`name` from `corr`.`departments` as `d` where (`d`.`id` = `e`.`department_id`)) as `n` from `corr`.`employees` as `e`",
	},
	{
		// The arms of a set operation see what the set operation sees.
		"set operation arm",
		"SELECT e.id FROM employees e WHERE e.department_id IN (SELECT d.id FROM departments d WHERE d.name = e.name UNION ALL SELECT o.department_id FROM employees o WHERE o.id = e.id)",
		"select `e`.`id` from `corr`.`employees` as `e` where `e`.`department_id` in (select `d`.`id` from `corr`.`departments` as `d` where (`d`.`name` = `e`.`name`) union all select `o`.`department_id` from `corr`.`employees` as `o` where (`o`.`id` = `e`.`id`))",
	},
	{
		"cte in subquery",
		"SELECT e.id FROM employees e WHERE EXISTS (WITH c AS (SELECT d.id FROM departments d) SELECT 1 FROM c WHERE c.id = e.department_id)",
		"select `e`.`id` from `corr`.`employees` as `e` where exists (with `c`(`id`) as (select `d`.`id` from `corr`.`departments` as `d`) select 1 as `1` from `c` where (`c`.`id` = `e`.`department_id`))",
	},
	{
		// A CTE body sees what the query holding its WITH clause sees.
		"cte body reads outer query",
		"SELECT e.id FROM employees e WHERE EXISTS (WITH c AS (SELECT d.id FROM departments d WHERE d.id = e.department_id) SELECT 1 FROM c)",
		"select `e`.`id` from `corr`.`employees` as `e` where exists (with `c`(`id`) as (select `d`.`id` from `corr`.`departments` as `d` where (`d`.`id` = `e`.`department_id`)) select 1 as `1` from `c`)",
	},
	{
		// A LATERAL derived table sees the query whose FROM holds it.
		"lateral",
		"SELECT e.id, l.n FROM employees e, LATERAL (SELECT d.name AS n FROM departments d WHERE d.id = e.department_id) AS l",
		"select `e`.`id`, `l`.`n` from `corr`.`employees` as `e`, lateral (select `d`.`name` as `n` from `corr`.`departments` as `d` where (`d`.`id` = `e`.`department_id`)) as `l`",
	},
}

// TestDeparseCorrelatedSubquery checks DeparseQuery resolves a column of an
// enclosing query against that query, and that the text it writes analyzes to
// a query it writes the same way again.
func TestDeparseCorrelatedSubquery(t *testing.T) {
	c := correlatedCatalog(t, "corr")
	for _, tc := range correlatedCases {
		t.Run(tc.name, func(t *testing.T) {
			got := DeparseQuery(analyzeQuantified(t, c, tc.sql))
			if got != tc.want {
				t.Errorf("DeparseQuery:\n got: %s\nwant: %s", got, tc.want)
			}
			if again := DeparseQuery(analyzeQuantified(t, c, got)); again != got {
				t.Errorf("deparse is not stable:\n first: %s\nsecond: %s", got, again)
			}
		})
	}

	t.Run("view", func(t *testing.T) {
		vc := correlatedCatalog(t, "corr")
		results, err := vc.Exec("CREATE VIEW v AS SELECT e.id FROM employees e WHERE EXISTS (SELECT 1 FROM departments d WHERE d.id = e.department_id)", nil)
		if err != nil {
			t.Fatal(err)
		}
		for _, r := range results {
			if r.Error != nil {
				t.Fatal(r.Error)
			}
		}
		v := vc.GetDatabase("corr").Views["v"]
		if v == nil || v.AnalyzedQuery == nil {
			t.Fatalf("view v has no analyzed query: %#v", v)
		}
		if got, want := DeparseQuery(v.AnalyzedQuery), "where (`d`.`id` = `e`.`department_id`)"; !strings.Contains(got, want) {
			t.Errorf("DeparseQuery = %s, want it to hold %s", got, want)
		}
	})
}

// checkCorrelatedOracle runs each correlated case and its DeparseQuery text on
// db, which holds correlatedDDL's tables and correlatedRows, and requires the
// same rows. c holds the same tables.
func checkCorrelatedOracle(t *testing.T, db *sql.DB, c *Catalog) {
	t.Helper()
	for _, cc := range correlatedCases {
		t.Run(cc.name, func(t *testing.T) {
			query := cc.sql + " ORDER BY 1"
			deparsed := DeparseQuery(analyzeQuantified(t, c, query))
			want, err := queryRowsText(db, query)
			if err != nil {
				t.Fatalf("original query failed: %v", err)
			}
			if want == "" {
				t.Fatalf("original query returned no rows, so it tells nothing apart")
			}
			got, err := queryRowsText(db, deparsed)
			if err != nil {
				t.Fatalf("deparsed: %s\nfailed: %v", deparsed, err)
			}
			if got != want {
				t.Errorf("deparsed: %s\noriginal rows: %s\ndeparsed rows: %s", deparsed, want, got)
			}
		})
	}
}

// execCorrelatedRows creates correlatedDDL's tables on db and fills them.
func execCorrelatedRows(t *testing.T, db *sql.DB) {
	t.Helper()
	for _, stmt := range strings.Split(correlatedDDL+";\n"+correlatedRows, ";\n") {
		if _, err := db.Exec(stmt); err != nil {
			t.Fatalf("%s: %v", stmt, err)
		}
	}
}

// TestDeparseCorrelatedSubqueryMySQL runs each correlated case and its
// DeparseQuery text on MySQL 8.0 and requires the same rows.
func TestDeparseCorrelatedSubqueryMySQL(t *testing.T) {
	ctx := context.Background()
	d, err := testcontainer.MySQL80.Database(ctx, "mysql_catalog_correlated")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { d.Close(ctx) })
	execCorrelatedRows(t, d.DB)
	checkCorrelatedOracle(t, d.DB, correlatedCatalog(t, d.Name))
}
