package parser

import "testing"

func TestSQLJSONTypedJsonLiterals(t *testing.T) {
	tests := []string{
		`SELECT json '{"a": 1}'`,
		`SELECT json '{"a": 1}' -> 'a'`,
		`SELECT JSON_ARRAY('aaa', 111, true, array[1,2,3], NULL, json '{"a": [1]}', jsonb '["a",3]')`,
		`SELECT *
		   FROM JSON_TABLE(
		     jsonb 'null', 'lax $[*]' PASSING 1 + 2 AS a, json '"foo"' AS "b c"
		     COLUMNS (v text PATH '$'))`,
		`SELECT JSON_EXISTS(json '{"a": 123}', '$' || '.' || 'a')`,
	}

	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			parseOK(t, sql)
		})
	}
}

func TestSQLJSONQueryFunctionPassingArguments(t *testing.T) {
	tests := []string{
		`SELECT JSON_QUERY(js, '$.date() < $x' PASSING '12:34'::timetz AS x) FROM t`,
		`SELECT JSON_QUERY(js, '$.date() < $x' PASSING '1234'::int AS x) FROM t`,
		`SELECT JSON_QUERY(js, '$.date() < $x' PASSING time '12:34' AS x) FROM t`,
		`CREATE INDEX ON t (JSON_QUERY(js, '$.date() < $x' PASSING '12:34'::timetz AS x))`,
	}

	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			parseOK(t, sql)
		})
	}
}

func TestSQLJSONTableExistsColumnBehaviorClauses(t *testing.T) {
	tests := []string{
		`SELECT * FROM JSON_TABLE(jsonb '1', '$' COLUMNS (a int EXISTS EMPTY OBJECT ON EMPTY))`,
		`SELECT * FROM JSON_TABLE(jsonb '1', '$' COLUMNS (a int EXISTS EMPTY OBJECT ON ERROR))`,
	}

	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			parseOK(t, sql)
		})
	}
}
