package analysis

import (
	"testing"

	"github.com/bytebase/omni/mssql"
	"github.com/stretchr/testify/require"
)

func TestDollarTokenStatements(t *testing.T) {
	accept := []string{
		`MERGE INTO dst AS d
USING src AS s ON d.k = s.k
WHEN MATCHED THEN UPDATE SET d.v = s.v
WHEN NOT MATCHED THEN INSERT (k, v) VALUES (s.k, s.v)
OUTPUT $action, INSERTED.k, INSERTED.v;`,
		"SELECT $IDENTITY, $ROWGUID FROM t;",
		"SELECT t.$IDENTITY FROM t;",
		"SELECT $node_id FROM Person;",
		"SELECT $12.50, -$4.78, $ 4, £10;",
		"SELECT $PARTITION.pf1(10);",
		"SELECT db1.$PARTITION.pf1(10);",
		"MERGE INTO dst AS d USING src AS s ON d.k = s.k WHEN MATCHED THEN UPDATE SET d.v = s.v OUTPUT $action, INSERTED.k INTO @changes;",
		"INSERT INTO e ($from_id, $to_id) VALUES ('a', 'b');",
		"UPDATE e SET $from_id = 'x';",
		"CREATE INDEX ix ON Person ($node_id);",
		"SELECT IDENTITYCOL, ROWGUIDCOL FROM t;",
	}
	for _, sql := range accept {
		t.Run(sql, func(t *testing.T) {
			_, err := mssql.Parse(sql)
			require.NoError(t, err)
		})
	}

	for _, sql := range []string{
		"SELECT $foo FROM t;",
		"SELECT $CUID FROM t;",
	} {
		t.Run(sql, func(t *testing.T) {
			_, err := mssql.Parse(sql)
			require.Error(t, err)
		})
	}
}
