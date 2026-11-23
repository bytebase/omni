package catalog

import (
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/bytebase/omni/enginetest/simulate"
	"github.com/bytebase/omni/mssql/internal/testcontainer"
)

// TestQuerySimulation holds the catalog to testdata/queries/simulate.tsv:
// every query (SELECT, WITH, or a parenthesized query) SQL Server accepts in
// the parser's compat.tsv, wrapped as CREATE VIEW _v AS <query> on
// one_table@column, and the hand-written T-SQL query features in
// queryFeatures on queries@tables, with what SQL Server 2022 did and what
// Exec does. Wrapping a query in a view holds the catalog to SQL Server's
// binding of it: the objects and columns it names, the view's column names,
// and the rules a view's query must keep.
//
// SQL Server's side is on the test's tape, so the test starts no server.
// OMNI_RECORD=1 runs every row on the pinned SQL Server again; -update-compat
// commits the rows a catalog fix moved; and -simulate-show=<text> prints both
// schemas of the rows that hold the text.
func TestQuerySimulation(t *testing.T) {
	rows, err := simulate.FromCompat(filepath.Join("..", "parser", "testdata", "compat.tsv"), []string{"one_table@column"}, simulatedQuery)
	if err != nil {
		t.Fatal(err)
	}
	seen := make(map[string]bool, len(rows))
	wrapped := rows[:0]
	for _, r := range rows {
		// Only the grammar's rows: compat.tsv's hand rows (suite:, issue:)
		// are corpora that grow on their own, and this file should not move
		// with them.
		if handLabel(r.Rule) {
			continue
		}
		r.Statement = wrapQuery(r.Statement)
		if seen[r.Fixture+"\x00"+r.Statement] {
			continue
		}
		seen[r.Fixture+"\x00"+r.Statement] = true
		wrapped = append(wrapped, r)
	}
	for _, f := range queryFeatures {
		stmt := wrapQuery(f.query)
		if seen["queries@tables\x00"+stmt] {
			t.Fatalf("queryFeatures holds %q twice", f.query)
		}
		seen["queries@tables\x00"+stmt] = true
		wrapped = append(wrapped, simulate.Generated{Fixture: "queries@tables", Rule: "feature:" + f.feature, Statement: stmt})
	}
	simulate.Check(t, simulate.Suite{
		Path:     filepath.Join("testdata", "queries", "simulate.tsv"),
		Pins:     testcontainer.Pins,
		Fixtures: filepath.Join("testdata", "fixtures"),
		Rows:     wrapped,
		MadeUp:   testcontainer.Codec{}.MadeUp(),
		Omni:     simulationOmni,
		Engine:   simulationEngine,
	})
}

var simulatedQueryStatement = regexp.MustCompile(`(?i)^(?:SELECT\b|WITH\s|\(\s*(?:SELECT|WITH|\())`)

func simulatedQuery(statement string) bool { return simulatedQueryStatement.MatchString(statement) }

func wrapQuery(query string) string {
	if strings.HasPrefix(query, "CREATE ") || strings.HasPrefix(query, "ALTER ") {
		return query
	}
	return "CREATE VIEW _v AS " + strings.TrimSpace(strings.TrimRight(strings.TrimSpace(query), ";"))
}

// handLabel reports whether a compat.tsv rule is a label a person wrote
// (issue:BYT-1, suite:lexer) rather than a grammar rule: a word and a colon
// before any slash.
func handLabel(rule string) bool {
	i := strings.IndexAny(rule, ":/")
	return i >= 0 && rule[i] == ':'
}

// queryFeatures are T-SQL query features the grammar's rows reach only with
// the one name a: joins over real keys, built-in functions, windows, CTEs,
// set operations, APPLY and PIVOT, on queries@tables. Each is labelled
// feature:<feature> and wrapped as CREATE VIEW _v AS <query>. Some are refused
// by SQL Server on purpose, to hold the catalog to the refusal.
var queryFeatures = []struct{ feature, query string }{
	// Select lists and naming.
	{"select-list", "SELECT * FROM a"},
	{"select-list", "SELECT a.* FROM a"},
	{"select-list", "SELECT dbo.a.* FROM dbo.a"},
	{"select-list", "SELECT b, a, d FROM a"},
	{"select-list", "SELECT b AS x, a y, z = d FROM a"},
	{"select-list", "SELECT [b] AS [x y], 'lit' AS [lit] FROM a"},
	{"select-list", "SELECT b, b FROM a"},
	{"select-list", "SELECT b AS x, a AS x FROM a"},
	{"select-list", "SELECT b + 1 FROM a"},
	{"select-list", "SELECT 1"},
	{"select-list", "SELECT 1 AS one"},
	{"select-list", "SELECT 1 AS one, 2"},
	{"select-list", "SELECT a.b FROM a AS t"},
	{"select-list", "SELECT t.b FROM a AS t"},
	{"select-list", "SELECT nope FROM a"},
	{"select-list", "SELECT x.b FROM a"},
	{"select-list", "SELECT b FROM nope"},
	{"select-list", "SELECT b FROM dbo.nope"},
	{"select-list", "SELECT * FROM s.orders"},
	{"select-list", "SELECT id, note FROM s.orders"},
	{"select-list", "SELECT * FROM v_customers"},
	{"select-list", "SELECT id, name FROM dbo.v_customers WHERE id > 1"},
	{"select-list", "SELECT DISTINCT region FROM customers"},
	{"select-list", "SELECT ALL region FROM customers"},
	{"select-list", "SELECT $IDENTITY FROM customers"},
	{"select-list", "SELECT $ROWGUID FROM customers"},
	{"select-list", "SELECT @x AS v"},
	{"select-list", "SELECT b INTO t2 FROM a"},
	{"select-list", "SELECT b FROM a ORDER BY b"},
	{"select-list", "SELECT b FROM #tmp"},
	{"select-list", "SELECT b FROM a WITH (NOLOCK)"},
	{"select-list", "SELECT b FROM a WITH (NOEXPAND)"},
	{"select-list", "SELECT b FROM a OPTION (MAXDOP 1)"},
	{"select-list", "SELECT name COLLATE Latin1_General_BIN2 AS n FROM customers"},
	{"select-list", "SELECT customers.name FROM dbo.customers"},
	{"select-list", "SELECT dbo.customers.name FROM dbo.customers"},
	{"select-list", "SELECT c.id FROM customers c WHERE c.region = 'EU'"},
	// Joins.
	{"join", "SELECT o.id, c.name FROM orders o JOIN customers c ON c.id = o.customer_id"},
	{"join", "SELECT o.id, c.name FROM orders o INNER JOIN customers c ON c.id = o.customer_id"},
	{"join", "SELECT o.id, c.name FROM orders o LEFT JOIN customers c ON c.id = o.customer_id"},
	{"join", "SELECT o.id, c.name FROM orders o LEFT OUTER JOIN customers c ON c.id = o.customer_id"},
	{"join", "SELECT o.id, c.name FROM orders o RIGHT JOIN customers c ON c.id = o.customer_id"},
	{"join", "SELECT o.id, c.name FROM orders o FULL OUTER JOIN customers c ON c.id = o.customer_id"},
	{"join", "SELECT o.id, c.name FROM orders o CROSS JOIN customers c"},
	{"join", "SELECT o.id, c.name FROM orders o, customers c WHERE c.id = o.customer_id"},
	{"join", "SELECT o.id, c.id FROM orders o JOIN customers c ON c.id = o.customer_id"},
	{"join", "SELECT id FROM orders o JOIN customers c ON c.id = o.customer_id"},
	{"join", "SELECT * FROM orders o JOIN customers c ON c.id = o.customer_id"},
	{"join", "SELECT o.*, c.name FROM orders o JOIN customers c ON c.id = o.customer_id"},
	{"join", "SELECT o.id FROM orders o JOIN customers c ON c.nope = o.customer_id"},
	{"join", "SELECT o.id FROM orders o JOIN customers o ON o.id = o.id"},
	{"join", "SELECT orders.id FROM orders JOIN s.orders ON s.orders.id = dbo.orders.id"},
	{"join", "SELECT o.id FROM orders o INNER HASH JOIN customers c ON c.id = o.customer_id"},
	{"join", "SELECT o.id FROM orders o JOIN customers c JOIN a ON a.b = c.id ON c.id = o.customer_id"},
	{"join", "SELECT o.id FROM (orders o JOIN customers c ON c.id = o.customer_id)"},
	{"join", "SELECT o.id FROM orders o JOIN customers c ON o.id = x.id"},
	// Predicates.
	{"where", "SELECT id FROM orders WHERE amount BETWEEN 1 AND 10"},
	{"where", "SELECT id FROM orders WHERE status IN ('new', 'paid')"},
	{"where", "SELECT id FROM orders WHERE status NOT IN ('new')"},
	{"where", "SELECT id FROM orders WHERE status LIKE 'p%' ESCAPE '\\'"},
	{"where", "SELECT id FROM orders WHERE status IS NULL OR amount IS NOT NULL"},
	{"where", "SELECT id FROM orders WHERE NOT (amount > 1 AND amount < 5)"},
	{"where", "SELECT id FROM orders WHERE EXISTS (SELECT 1 FROM customers c WHERE c.id = orders.customer_id)"},
	{"where", "SELECT id FROM orders WHERE customer_id IN (SELECT id FROM customers)"},
	{"where", "SELECT id FROM orders WHERE amount > ALL (SELECT amount FROM orders)"},
	{"where", "SELECT id FROM orders WHERE amount > ANY (SELECT amount FROM orders)"},
	{"where", "SELECT id FROM orders WHERE amount = SOME (SELECT amount FROM orders)"},
	{"where", "SELECT id FROM orders WHERE customer_id IN (SELECT nope FROM customers)"},
	{"where", "SELECT id FROM orders WHERE customer_id IN (SELECT id, name FROM customers)"},
	{"where", "SELECT id FROM orders WHERE amount > (SELECT MAX(amount) FROM orders)"},
	{"where", "SELECT id FROM orders WHERE amount > (SELECT amount, id FROM orders)"},
	{"where", "SELECT id FROM orders WHERE IIF(amount > 1, 1, 0) = 1"},
	{"where", "SELECT id FROM orders WHERE amount IS DISTINCT FROM 1"},
	{"where", "SELECT id FROM orders WHERE CONTAINS(status, 'x')"},
	{"where", "SELECT id FROM orders WHERE ROW_NUMBER() OVER (ORDER BY id) = 1"},
	{"where", "SELECT id FROM orders WHERE SUM(amount) > 1"},
	// Scalar expressions.
	{"case", "SELECT CASE WHEN amount > 100 THEN 'big' ELSE 'small' END AS size FROM orders"},
	{"case", "SELECT CASE status WHEN 'new' THEN 1 WHEN 'paid' THEN 2 END AS s FROM orders"},
	{"case", "SELECT CASE WHEN amount > 100 THEN 'big' END FROM orders"},
	{"case", "SELECT CASE WHEN nope > 1 THEN 1 END AS c FROM orders"},
	{"case", "SELECT IIF(amount > 1, 'y', 'n') AS c FROM orders"},
	{"case", "SELECT CHOOSE(2, 'a', 'b', 'c') AS c"},
	{"case", "SELECT COALESCE(status, 'none') AS c FROM orders"},
	{"case", "SELECT ISNULL(status, 'none') AS c FROM orders"},
	{"case", "SELECT NULLIF(status, '') AS c FROM orders"},
	{"case", "SELECT COALESCE(status) AS c FROM orders"},
	{"cast", "SELECT CAST(amount AS INT) AS c FROM orders"},
	{"cast", "SELECT CAST(amount AS VARCHAR(20)) AS c FROM orders"},
	{"cast", "SELECT CAST(amount AS NOPE) AS c FROM orders"},
	{"cast", "SELECT CONVERT(VARCHAR(10), created, 120) AS c FROM orders"},
	{"cast", "SELECT CONVERT(INT, amount) AS c FROM orders"},
	{"cast", "SELECT TRY_CAST(status AS INT) AS c FROM orders"},
	{"cast", "SELECT TRY_CONVERT(INT, status) AS c FROM orders"},
	{"cast", "SELECT PARSE('1' AS INT) AS c"},
	{"cast", "SELECT TRY_PARSE('1' AS INT USING 'en-US') AS c"},
	{"cast", "SELECT CAST(doc AS NVARCHAR(MAX)) AS c FROM orders"},
	{"cast", "SELECT created AT TIME ZONE 'UTC' AS c FROM orders"},
	{"cast", "SELECT CAST(amount AS INT) FROM orders"},
	{"cast", "SELECT id COLLATE Latin1_General_BIN2 AS c FROM customers"},
	// Built-in functions.
	{"function", "SELECT LEN(name) AS l, UPPER(name) AS u, LOWER(name) AS lo FROM customers"},
	{"function", "SELECT SUBSTRING(name, 1, 3) AS s, LEFT(name, 2) AS l, RIGHT(name, 2) AS r FROM customers"},
	{"function", "SELECT LTRIM(RTRIM(name)) AS t, TRIM(name) AS t2 FROM customers"},
	{"function", "SELECT TRIM('x' FROM name) AS t FROM customers"},
	{"function", "SELECT REPLACE(name, 'a', 'b') AS r, CHARINDEX('a', name) AS c, PATINDEX('%a%', name) AS p FROM customers"},
	{"function", "SELECT CONCAT(name, '-', region) AS c, CONCAT_WS(',', name, region) AS w FROM customers"},
	{"function", "SELECT STUFF(name, 1, 1, 'x') AS s, REVERSE(name) AS r, REPLICATE('x', 3) AS p FROM customers"},
	{"function", "SELECT FORMAT(amount, 'N2') AS f, STR(amount, 10, 2) AS s FROM orders"},
	{"function", "SELECT GETDATE() AS d, SYSDATETIME() AS s, GETUTCDATE() AS u, CURRENT_TIMESTAMP AS c"},
	{"function", "SELECT DATEADD(day, 1, created) AS d, DATEDIFF(day, created, GETDATE()) AS n FROM orders"},
	{"function", "SELECT DATEPART(year, created) AS y, DATENAME(month, created) AS m, YEAR(created) AS y2, MONTH(created) AS m2, DAY(created) AS d2 FROM orders"},
	{"function", "SELECT EOMONTH(created) AS e, DATEFROMPARTS(2020, 1, 1) AS d FROM orders"},
	{"function", "SELECT DATEADD(nope, 1, created) AS d FROM orders"},
	{"function", "SELECT ABS(amount) AS a, ROUND(amount, 1) AS r, CEILING(amount) AS c, FLOOR(amount) AS f, POWER(amount, 2) AS p, SQRT(amount) AS s FROM orders"},
	{"function", "SELECT NEWID() AS g"},
	{"function", "SELECT RAND() AS r"},
	{"function", "SELECT @@ROWCOUNT AS r"},
	{"function", "SELECT @@SPID AS s, @@VERSION AS v"},
	{"function", "SELECT DB_NAME() AS d, SCHEMA_NAME() AS s, OBJECT_ID('dbo.a') AS o, USER_NAME() AS u, SUSER_SNAME() AS n"},
	{"function", "SELECT ISJSON(status) AS j, JSON_VALUE(status, '$.a') AS v, JSON_QUERY(status, '$.a') AS q FROM orders"},
	{"function", "SELECT JSON_MODIFY(status, '$.a', 1) AS m FROM orders"},
	{"function", "SELECT doc.value('(/a)[1]', 'INT') AS v FROM orders"},
	{"function", "SELECT doc.query('/a') AS q, doc.exist('/a') AS e FROM orders"},
	{"function", "SELECT STRING_SPLIT('a,b', ',') AS s"},
	{"function", "SELECT LEN(name, 1) AS l FROM customers"},
	{"function", "SELECT LEN() AS l FROM customers"},
	{"function", "SELECT NOPE(name) AS l FROM customers"},
	{"function", "SELECT dbo.f_double(id) AS d FROM customers"},
	{"function", "SELECT f_double(id) AS d FROM customers"},
	{"function", "SELECT dbo.f_double(id, 1) AS d FROM customers"},
	{"function", "SELECT dbo.nope(id) AS d FROM customers"},
	{"function", "SELECT GREATEST(1, 2) AS g, LEAST(1, 2) AS l"},
	{"function", "SELECT DATETRUNC(month, created) AS d FROM orders"},
	{"function", "SELECT STRING_ESCAPE(name, 'json') AS e FROM customers"},
	{"function", "SELECT TRANSLATE(name, 'ab', 'cd') AS t FROM customers"},
	{"function", "SELECT SOUNDEX(name) AS s, DIFFERENCE(name, 'x') AS d FROM customers"},
	{"function", "SELECT HASHBYTES('SHA2_256', name) AS h FROM customers"},
	{"function", "SELECT CHECKSUM(*) AS c FROM customers"},
	{"function", "SELECT BINARY_CHECKSUM(name) AS c FROM customers"},
	{"function", "SELECT NEXT VALUE FOR nope AS n"},
	{"function", "SELECT CONCAT(name) AS c FROM customers"},
	{"function", "SELECT SUBSTRING(name, 1) AS s FROM customers"},
	{"function", "SELECT ROUND(amount) AS r FROM orders"},
	{"function", "SELECT DATEDIFF(nope, created, created) AS d FROM orders"},
	{"function", "SELECT TRIM(id) AS t FROM customers"},
	{"function", "SELECT UPPER(id) AS u FROM customers"},
	// Aggregates and grouping.
	{"group", "SELECT customer_id, COUNT(*) AS n, SUM(amount) AS s FROM orders GROUP BY customer_id"},
	{"group", "SELECT customer_id, COUNT(DISTINCT status) AS n FROM orders GROUP BY customer_id"},
	{"group", "SELECT customer_id, AVG(amount) AS a, MIN(amount) AS mi, MAX(amount) AS ma FROM orders GROUP BY customer_id"},
	{"group", "SELECT customer_id, COUNT_BIG(*) AS n FROM orders GROUP BY customer_id"},
	{"group", "SELECT customer_id, STDEV(amount) AS s, VAR(amount) AS v FROM orders GROUP BY customer_id"},
	{"group", "SELECT customer_id, STRING_AGG(status, ',') AS s FROM orders GROUP BY customer_id"},
	{"group", "SELECT customer_id, STRING_AGG(status, ',') WITHIN GROUP (ORDER BY status) AS s FROM orders GROUP BY customer_id"},
	{"group", "SELECT customer_id, SUM(amount) AS s FROM orders GROUP BY customer_id HAVING SUM(amount) > 10"},
	{"group", "SELECT customer_id, status, COUNT(*) AS n FROM orders GROUP BY customer_id"},
	{"group", "SELECT status FROM orders GROUP BY customer_id"},
	{"group", "SELECT customer_id FROM orders GROUP BY customer_id HAVING status = 'new'"},
	{"group", "SELECT COUNT(*) AS n FROM orders"},
	{"group", "SELECT COUNT(*) FROM orders"},
	{"group", "SELECT customer_id, COUNT(*) AS n FROM orders GROUP BY ROLLUP (customer_id)"},
	{"group", "SELECT customer_id, status, COUNT(*) AS n FROM orders GROUP BY CUBE (customer_id, status)"},
	{"group", "SELECT customer_id, status, COUNT(*) AS n FROM orders GROUP BY GROUPING SETS ((customer_id), (status), ())"},
	{"group", "SELECT customer_id, GROUPING(customer_id) AS g, COUNT(*) AS n FROM orders GROUP BY ROLLUP (customer_id)"},
	{"group", "SELECT customer_id, COUNT(*) AS n FROM orders GROUP BY customer_id WITH ROLLUP"},
	{"group", "SELECT SUM(SUM(amount)) AS s FROM orders GROUP BY customer_id"},
	{"group", "SELECT SUM(name) AS s FROM customers"},
	{"group", "SELECT customer_id FROM orders GROUP BY nope"},
	{"group", "SELECT SUM(amount) AS s FROM orders GROUP BY 1"},
	{"group", "SELECT customer_id, APPROX_COUNT_DISTINCT(status) AS n FROM orders GROUP BY customer_id"},
	{"group", "SELECT CASE WHEN amount > 10 THEN 1 ELSE 0 END AS big, COUNT(*) AS n FROM orders GROUP BY CASE WHEN amount > 10 THEN 1 ELSE 0 END"},
	{"group", "SELECT YEAR(created) AS y, COUNT(*) AS n FROM orders GROUP BY YEAR(created)"},
	{"group", "SELECT YEAR(created) AS y, COUNT(*) AS n FROM orders GROUP BY MONTH(created)"},
	{"group", "SELECT o.customer_id, c.name, SUM(o.amount) AS s FROM orders o JOIN customers c ON c.id = o.customer_id GROUP BY o.customer_id, c.name"},
	{"group", "SELECT customer_id, SUM(amount) AS s FROM orders GROUP BY customer_id HAVING COUNT(*) > 1 AND MAX(amount) > 5"},
	{"group", "SELECT TOP 3 customer_id, SUM(amount) AS s FROM orders GROUP BY customer_id ORDER BY s DESC"},
	{"group", "SELECT TOP 3 customer_id FROM orders GROUP BY customer_id ORDER BY SUM(amount) DESC"},
	{"group", "SELECT TOP 3 customer_id FROM orders GROUP BY customer_id ORDER BY amount"},
	{"group", "SELECT customer_id + 1 AS c FROM orders GROUP BY customer_id"},
	{"group", "SELECT (SELECT COUNT(*) FROM customers c WHERE c.id = o.customer_id) AS n FROM orders o GROUP BY o.customer_id"},
	{"group", "SELECT AVG(name) AS a FROM customers"},
	{"group", "SELECT SUM(born) AS s FROM customers"},
	{"group", "SELECT MAX(name) AS m, MIN(born) AS b FROM customers"},
	// Windows.
	{"window", "SELECT id, ROW_NUMBER() OVER (ORDER BY id) AS rn FROM orders"},
	{"window", "SELECT id, ROW_NUMBER() OVER (PARTITION BY customer_id ORDER BY created DESC) AS rn FROM orders"},
	{"window", "SELECT id, RANK() OVER (ORDER BY amount) AS r, DENSE_RANK() OVER (ORDER BY amount) AS d, NTILE(4) OVER (ORDER BY amount) AS n FROM orders"},
	{"window", "SELECT id, SUM(amount) OVER (PARTITION BY customer_id) AS s FROM orders"},
	{"window", "SELECT id, SUM(amount) OVER (ORDER BY id ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) AS s FROM orders"},
	{"window", "SELECT id, AVG(amount) OVER (ORDER BY id ROWS BETWEEN 2 PRECEDING AND 2 FOLLOWING) AS s FROM orders"},
	{"window", "SELECT id, SUM(amount) OVER (ORDER BY id RANGE BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) AS s FROM orders"},
	{"window", "SELECT id, LAG(amount) OVER (ORDER BY id) AS p, LEAD(amount, 1, 0) OVER (ORDER BY id) AS n FROM orders"},
	{"window", "SELECT id, FIRST_VALUE(amount) OVER (ORDER BY id) AS f, LAST_VALUE(amount) OVER (ORDER BY id) AS l FROM orders"},
	{"window", "SELECT id, PERCENT_RANK() OVER (ORDER BY amount) AS p, CUME_DIST() OVER (ORDER BY amount) AS c FROM orders"},
	{"window", "SELECT id, PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY amount) OVER (PARTITION BY customer_id) AS m FROM orders"},
	{"window", "SELECT id, ROW_NUMBER() OVER () AS rn FROM orders"},
	{"window", "SELECT id, LAG(amount) OVER (PARTITION BY customer_id) AS p FROM orders"},
	{"window", "SELECT id, ROW_NUMBER() OVER w AS rn FROM orders WINDOW w AS (ORDER BY id)"},
	{"window", "SELECT id, ROW_NUMBER() OVER (ORDER BY nope) AS rn FROM orders"},
	{"window", "SELECT id, COUNT(*) OVER () AS n FROM orders"},
	{"window", "SELECT customer_id, SUM(SUM(amount)) OVER () AS total FROM orders GROUP BY customer_id"},
	{"window", "SELECT id, ROW_NUMBER() OVER w AS rn, RANK() OVER w AS r FROM orders WINDOW w AS (PARTITION BY customer_id ORDER BY id)"},
	{"window", "SELECT id, ROW_NUMBER() OVER w AS rn FROM orders WINDOW w AS (PARTITION BY customer_id)"},
	{"window", "SELECT id, ROW_NUMBER() OVER nope AS rn FROM orders"},
	// TOP, ORDER BY, OFFSET/FETCH.
	{"top", "SELECT TOP 10 id FROM orders ORDER BY id"},
	{"top", "SELECT TOP (10) id FROM orders ORDER BY created DESC"},
	{"top", "SELECT TOP 10 PERCENT id FROM orders ORDER BY id"},
	{"top", "SELECT TOP 10 WITH TIES id FROM orders ORDER BY amount"},
	{"top", "SELECT TOP 10 WITH TIES id FROM orders"},
	{"top", "SELECT TOP 100 PERCENT id FROM orders ORDER BY id"},
	{"top", "SELECT TOP 5 id FROM orders"},
	{"top", "SELECT id FROM orders ORDER BY id OFFSET 0 ROWS"},
	{"top", "SELECT id FROM orders ORDER BY id OFFSET 10 ROWS FETCH NEXT 5 ROWS ONLY"},
	{"top", "SELECT id FROM orders ORDER BY id OFFSET 10 ROWS FETCH FIRST 5 ROW ONLY"},
	{"top", "SELECT TOP 5 id FROM orders ORDER BY id OFFSET 1 ROWS"},
	{"top", "SELECT id FROM orders ORDER BY nope"},
	{"top", "SELECT TOP 1 id FROM orders ORDER BY 1"},
	{"top", "SELECT TOP 1 id AS x FROM orders ORDER BY x"},
	{"top", "SELECT TOP 1 id FROM orders ORDER BY 2"},
	{"top", "SELECT DISTINCT TOP 1 status FROM orders ORDER BY id"},
	// Derived tables and subqueries.
	{"derived", "SELECT x.id FROM (SELECT id FROM orders) AS x"},
	{"derived", "SELECT x.n FROM (SELECT id AS n FROM orders) x"},
	{"derived", "SELECT x.n FROM (SELECT id FROM orders) x (n)"},
	{"derived", "SELECT x.id FROM (SELECT id FROM orders) x (n)"},
	{"derived", "SELECT x.id FROM (SELECT id FROM orders)"},
	{"derived", "SELECT x.c FROM (SELECT COUNT(*) FROM orders) x"},
	{"derived", "SELECT x.c FROM (SELECT COUNT(*) FROM orders) x (c)"},
	{"derived", "SELECT x.id FROM (SELECT id, id FROM orders) x"},
	{"derived", "SELECT x.a FROM (VALUES (1, 2), (3, 4)) AS x (a, b)"},
	{"derived", "SELECT x.a FROM (VALUES (1, 2), (3, 4)) AS x (a)"},
	{"derived", "SELECT x.a FROM (VALUES (1, 2), (3)) AS x (a, b)"},
	{"derived", "SELECT (SELECT MAX(amount) FROM orders) AS m"},
	{"derived", "SELECT c.id, (SELECT COUNT(*) FROM orders o WHERE o.customer_id = c.id) AS n FROM customers c"},
	{"derived", "SELECT x.id FROM (SELECT id FROM orders ORDER BY id) x"},
	{"derived", "SELECT x.id FROM (SELECT TOP 1 id FROM orders ORDER BY id) x"},
	{"derived", "SELECT f.id, f.amount FROM dbo.f_orders(1) AS f"},
	{"derived", "SELECT f.id FROM dbo.f_orders(1, 2) AS f"},
	{"derived", "SELECT * FROM dbo.f_double(1)"},
	{"derived", "SELECT f.nope FROM dbo.f_orders(1) AS f"},
	{"derived", "SELECT value FROM STRING_SPLIT('a,b', ',')"},
	{"derived", "SELECT [key], value FROM OPENJSON('{}')"},
	{"derived", "SELECT j.a FROM OPENJSON('{}') WITH (a INT '$.a') AS j"},
	{"derived", "SELECT value FROM GENERATE_SERIES(1, 10)"},
	{"derived", "SELECT d.c FROM (SELECT customer_id AS c FROM orders GROUP BY customer_id) d JOIN customers ON customers.id = d.c"},
	{"derived", "SELECT o.id FROM orders o WHERE o.amount > (SELECT AVG(i.amount) FROM orders i WHERE i.customer_id = o.customer_id)"},
	{"derived", "SELECT o.id FROM orders o JOIN (SELECT id FROM customers WHERE id = o.customer_id) c ON c.id = o.id"},
	// CTEs.
	{"cte", "WITH x AS (SELECT id FROM orders) SELECT id FROM x"},
	{"cte", "WITH x (n) AS (SELECT id FROM orders) SELECT n FROM x"},
	{"cte", "WITH x (n) AS (SELECT id FROM orders) SELECT id FROM x"},
	{"cte", "WITH x (n, m) AS (SELECT id FROM orders) SELECT n FROM x"},
	{"cte", "WITH x AS (SELECT id FROM orders), y AS (SELECT id FROM x) SELECT y.id FROM y JOIN x ON x.id = y.id"},
	{"cte", "WITH x AS (SELECT COUNT(*) FROM orders) SELECT * FROM x"},
	{"cte", "WITH x AS (SELECT id FROM orders), x AS (SELECT id FROM orders) SELECT id FROM x"},
	{"cte", "WITH r (n) AS (SELECT 1 UNION ALL SELECT n + 1 FROM r WHERE n < 10) SELECT n FROM r"},
	{"cte", "WITH r (n) AS (SELECT n + 1 FROM r) SELECT n FROM r"},
	{"cte", "WITH x AS (SELECT id FROM nope) SELECT id FROM x"},
	{"cte", "WITH x AS (SELECT id FROM orders) SELECT id FROM orders"},
	{"cte", "WITH x AS (SELECT id FROM orders ORDER BY id) SELECT id FROM x"},
	{"cte", "WITH x AS (SELECT id FROM orders) SELECT a.id FROM x a JOIN x b ON a.id = b.id"},
	{"cte", "WITH r (n, p) AS (SELECT id, CAST(NULL AS INT) FROM customers UNION ALL SELECT r.n, r.p FROM r JOIN customers c ON c.id = r.n WHERE 1 = 0) SELECT n, p FROM r"},
	{"cte", "WITH x AS (SELECT id FROM orders) SELECT id FROM x WHERE id IN (SELECT id FROM x)"},
	{"cte", "WITH x AS (SELECT 1 AS n), y AS (SELECT n FROM z), z AS (SELECT 1 AS n) SELECT n FROM y"},
	// Set operations.
	{"setop", "SELECT id FROM orders UNION SELECT id FROM customers"},
	{"setop", "SELECT id FROM orders UNION ALL SELECT id FROM customers"},
	{"setop", "SELECT id FROM orders INTERSECT SELECT id FROM customers"},
	{"setop", "SELECT id FROM orders EXCEPT SELECT id FROM customers"},
	{"setop", "SELECT id, status FROM orders UNION SELECT id FROM customers"},
	{"setop", "SELECT id AS x FROM orders UNION SELECT id AS y FROM customers"},
	{"setop", "SELECT 1 UNION SELECT 2"},
	{"setop", "(SELECT id FROM orders) UNION (SELECT id FROM customers)"},
	{"setop", "SELECT id FROM orders UNION SELECT id FROM customers ORDER BY id"},
	{"setop", "SELECT TOP 1 id FROM orders UNION SELECT id FROM customers"},
	{"setop", "SELECT id FROM orders UNION SELECT nope FROM customers"},
	{"setop", "SELECT id FROM orders UNION SELECT id FROM customers UNION ALL SELECT b FROM a"},
	{"setop", "SELECT TOP 2 id FROM orders UNION ALL SELECT TOP 2 id FROM customers"},
	{"setop", "SELECT id FROM orders EXCEPT SELECT b FROM a INTERSECT SELECT id FROM customers"},
	// APPLY.
	{"apply", "SELECT c.id, x.amount FROM customers c CROSS APPLY (SELECT TOP 1 amount FROM orders o WHERE o.customer_id = c.id ORDER BY created DESC) x"},
	{"apply", "SELECT c.id, x.amount FROM customers c OUTER APPLY (SELECT TOP 1 amount FROM orders o WHERE o.customer_id = c.id ORDER BY created DESC) x"},
	{"apply", "SELECT c.id, f.amount FROM customers c CROSS APPLY dbo.f_orders(c.id) f"},
	{"apply", "SELECT c.id, s.value FROM customers c CROSS APPLY STRING_SPLIT(c.name, ' ') s"},
	{"apply", "SELECT c.id, x.n FROM customers c CROSS APPLY (SELECT c.nope AS n) x"},
	{"apply", "SELECT c.id, x.n FROM customers c CROSS APPLY (VALUES (c.id), (c.id + 1)) x (n)"},
	{"apply", "SELECT c.id, x.n FROM customers c OUTER APPLY (SELECT COUNT(*) AS n FROM orders o WHERE o.customer_id = c.id) x"},
	{"apply", "SELECT c.id, x.n FROM customers c JOIN (SELECT COUNT(*) AS n FROM orders o WHERE o.customer_id = c.id) x ON 1 = 1"},
	// PIVOT and UNPIVOT.
	{"pivot", "SELECT p.customer_id, p.[new], p.[paid] FROM (SELECT customer_id, status, amount FROM orders) AS src PIVOT (SUM(amount) FOR status IN ([new], [paid])) AS p"},
	{"pivot", "SELECT * FROM (SELECT customer_id, status, amount FROM orders) AS src PIVOT (SUM(amount) FOR status IN ([new], [paid])) AS p"},
	{"pivot", "SELECT p.customer_id FROM (SELECT customer_id, status, amount FROM orders) AS src PIVOT (SUM(amount) FOR nope IN ([new])) AS p"},
	{"pivot", "SELECT p.amount FROM (SELECT customer_id, status, amount FROM orders) AS src PIVOT (SUM(amount) FOR status IN ([new])) AS p"},
	{"pivot", "SELECT u.id, u.k, u.v FROM (SELECT id, a, d FROM a AS t (id, a, d)) AS src UNPIVOT (v FOR k IN (a, d)) AS u"},
	{"pivot", "SELECT u.id, u.k, u.v FROM (SELECT b AS id, a, d FROM a) AS src UNPIVOT (v FOR k IN (a, d)) AS u"},
	{"pivot", "SELECT u.b, u.k, u.v FROM a UNPIVOT (v FOR k IN (a, d)) AS u"},
	{"pivot", "SELECT u.b, u.k, u.v FROM a UNPIVOT (v FOR k IN (a, nope)) AS u"},
	// FOR XML / FOR JSON.
	{"for", "SELECT id, name FROM customers FOR XML PATH('c'), ROOT('cs')"},
	{"for", "SELECT id, name FROM customers FOR XML AUTO"},
	{"for", "SELECT id, name FROM customers FOR JSON PATH"},
	{"for", "SELECT (SELECT id, name FROM customers FOR JSON PATH) AS j"},
	{"for", "SELECT (SELECT id, name FROM customers FOR XML PATH('c'), TYPE) AS x"},
	{"for", "SELECT id FROM customers ORDER BY id FOR XML AUTO"},
	// Table hints, TABLESAMPLE, system tables.
	{"source", "SELECT name FROM sys.objects"},
	{"source", "SELECT name FROM sys.nope"},
	{"source", "SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES"},
	{"source", "SELECT b FROM a TABLESAMPLE (10 PERCENT)"},
	{"source", "SELECT name FROM master.dbo.spt_values"},
	{"source", "SELECT b FROM nope.dbo.a"},
	{"source", "SELECT b FROM OPENROWSET('x', 'y', 'z')"},
	{"source", "SELECT b FROM a FOR SYSTEM_TIME AS OF '2020-01-01'"},
	// View options.
	{"view", "SELECT b FROM a WITH CHECK OPTION"},
	{"view", "CREATE VIEW _v (x, y) AS SELECT b, a FROM a"},
	{"view", "CREATE VIEW _v (x) AS SELECT b, a FROM a"},
	{"view", "CREATE VIEW _v (x, y, z) AS SELECT b, a FROM a"},
	{"view", "CREATE VIEW _v (x, x) AS SELECT b, a FROM a"},
	{"view", "CREATE VIEW _v (x, y) AS SELECT b, b FROM a"},
	{"view", "CREATE VIEW _v (x) AS SELECT COUNT(*) FROM a"},
	{"view", "CREATE VIEW _v WITH SCHEMABINDING AS SELECT b, a FROM dbo.a"},
	{"view", "CREATE VIEW _v WITH SCHEMABINDING AS SELECT b, a FROM a"},
	{"view", "CREATE VIEW _v WITH SCHEMABINDING AS SELECT * FROM dbo.a"},
	{"view", "CREATE VIEW _v WITH SCHEMABINDING AS SELECT dbo.f_double(b) AS x FROM dbo.a"},
	{"view", "CREATE VIEW _v WITH SCHEMABINDING AS SELECT id FROM dbo.v_customers"},
	{"view", "CREATE VIEW _v WITH VIEW_METADATA AS SELECT b FROM a"},
	{"view", "CREATE VIEW dbo._v AS SELECT b FROM a"},
	{"view", "CREATE VIEW s._v AS SELECT id FROM orders"},
	{"view", "CREATE VIEW nope._v AS SELECT b FROM a"},
	{"view", "CREATE VIEW a AS SELECT b FROM a"},
	{"view", "CREATE VIEW v_customers AS SELECT b FROM a"},
	{"view", "CREATE OR ALTER VIEW v_customers AS SELECT b FROM a"},
	{"view", "ALTER VIEW v_customers AS SELECT b AS id, a AS name FROM a"},
	{"view", "ALTER VIEW v_customers AS SELECT nope FROM a"},
	{"view", "CREATE VIEW _v AS SELECT b FROM a WHERE b IN (SELECT b FROM _v)"},
	{"view", "CREATE VIEW _v WITH SCHEMABINDING AS SELECT o.id, c.name FROM dbo.orders o JOIN dbo.customers c ON c.id = o.customer_id"},
	{"view", "CREATE VIEW _v WITH SCHEMABINDING AS SELECT id FROM dbo.f_orders(1)"},
	{"view", "CREATE VIEW _v WITH SCHEMABINDING AS SELECT COUNT_BIG(*) AS n, customer_id FROM dbo.orders GROUP BY customer_id"},
	{"view", "CREATE VIEW _v WITH SCHEMABINDING AS SELECT name FROM sys.objects"},
}
