package parser

import "testing"

func TestViewBodyFromSubqueryJoinShapes(t *testing.T) {
	tests := []string{
		`CREATE VIEW v AS
		 SELECT a.id, b.id
		   FROM (SELECT 1 AS id) a
		   JOIN (SELECT 1 AS id) b ON (a.id = b.id)`,
		`CREATE VIEW v AS
		 SELECT a.id, b.id
		   FROM (SELECT 1 AS id) a
		   FULL JOIN (SELECT 1 AS id) b ON (a.id = b.id)`,
		`CREATE VIEW v AS
		 SELECT a.id, b.id
		   FROM (SELECT 1 AS id) AS a
		   LEFT JOIN ((SELECT 1 AS id) b JOIN (SELECT 1 AS id) c ON (b.id = c.id))
		     ON (a.id = b.id)`,
		`CREATE VIEW v AS
		 SELECT r.flag
		   FROM (SELECT true AS flag) r
		  WHERE r.flag`,
	}

	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			parseOK(t, sql)
		})
	}
}

func TestViewBodyNestedParenthesizedJoinTree(t *testing.T) {
	tests := []string{
		`CREATE VIEW v AS
		 SELECT *
		   FROM (((a JOIN b ON (a.id = b.id))
		   LEFT JOIN c ON (b.id = c.id))
		   LEFT JOIN d ON (c.id = d.id))`,
		`CREATE VIEW v AS
		 SELECT *
		   FROM ((a JOIN b ON (a.id = b.id)) AS ab
		   LEFT JOIN c ON (ab.id = c.id))`,
		`CREATE VIEW v AS
		 SELECT *
		   FROM (a JOIN b ON (a.id = b.id)) ab
		   FULL JOIN (SELECT 1 AS id) s ON (ab.id = s.id)`,
	}

	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			parseOK(t, sql)
		})
	}
}

func TestViewBodyUnparenthesizedNestedJoinTree(t *testing.T) {
	tests := []string{
		`CREATE VIEW v AS
		 SELECT *
		   FROM a
		   LEFT JOIN b
		     JOIN c ON (b.id = c.id)
		   ON (a.id = b.id)`,
		`CREATE VIEW v AS
		 SELECT *
		   FROM a
		   INNER JOIN b
		     LEFT JOIN c ON (b.id = c.id)
		   ON (a.id = b.id)`,
		`CREATE VIEW v AS
		 SELECT *
		   FROM a
		   LEFT JOIN b
		     LEFT JOIN c
		       FULL JOIN d ON false
		     ON false
		   ON true`,
	}

	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			parseOK(t, sql)
		})
	}
}

func TestViewBodyTypeFuncKeywordFunctionNames(t *testing.T) {
	tests := []string{
		`CREATE VIEW v AS
		 SELECT left(a, 10), right(a, 2)
		   FROM t`,
		`CREATE VIEW v AS
		 SELECT full(a), inner(a), cross(a), natural(a)
		   FROM t`,
	}

	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			parseOK(t, sql)
		})
	}
}

func TestViewBodySpecialFunctionSyntaxes(t *testing.T) {
	tests := []string{
		`CREATE VIEW v AS
		 SELECT substring('foo' similar 'f' escape '#') AS ss`,
		`CREATE VIEW v AS
		 SELECT substring('foo' from 2 for 3) AS s,
		        overlay('foo' placing 'bar' from 2 for 3) AS ovl,
		        position('foo' in 'foobar') AS p`,
	}

	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			parseOK(t, sql)
		})
	}
}

func TestViewBodyParenthesizedSetOpSubquery(t *testing.T) {
	tests := []string{
		`CREATE VIEW v AS
		 SELECT ((SELECT 2) UNION SELECT 2)`,
		`CREATE VIEW v AS
		 SELECT (((SELECT 2)) UNION SELECT 2)`,
		`CREATE VIEW v AS
		 SELECT ((SELECT max(x) FROM t) - '1 day'::interval) AS cutoff`,
		`CREATE VIEW v AS
		 SELECT ((SELECT max(x) FROM t) + '2 days'::interval) AS cutoff`,
	}

	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			parseOK(t, sql)
		})
	}
}
