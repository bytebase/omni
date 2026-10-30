package completion

import (
	"testing"
	"time"

	"github.com/bytebase/omni/mssql/parser"
)

// tableSourceForms are the FROM-clause and rowset forms whose argument and
// column lists completion parses with loops of their own.
var tableSourceForms = []string{
	"SELECT * FROM PREDICT (MODEL = @model, DATA = dbo.mytable AS d) WITH (Score FLOAT, Label NVARCHAR(10)) AS p",
	"SELECT * FROM PREDICT (MODEL = @model, DATA = dbo.mytable AS d, RUNTIME = ONNX) WITH (Score DECIMAL(10, 2) COLLATE Latin1_General_CI_AS NOT NULL)",
	"SELECT d.*, p.Score FROM PREDICT (MODEL = (SELECT model FROM models WHERE id = 1), DATA = dbo.mytable AS d) WITH (Score FLOAT, Label NVARCHAR(10)) AS p",
	"SELECT * FROM OPENJSON(@json) WITH (id INT '$.id', name NVARCHAR(100) '$.name', tags NVARCHAR(MAX) '$.tags' AS JSON) AS j",
	"SELECT * FROM OPENJSON(@json, '$.items') AS j CROSS APPLY OPENJSON(j.value) WITH (k INT) AS x",
	"SELECT * FROM OPENROWSET(BULK 'data.csv', FORMATFILE = 'fmt.xml', FIRSTROW = 2) AS r",
	"SELECT * FROM OPENROWSET(BULK 'data.json', SINGLE_CLOB) AS j",
	"SELECT * FROM OPENROWSET('SQLNCLI', 'Server=s;Trusted_Connection=yes;', 'SELECT a FROM t') AS r",
	"SELECT * FROM OPENROWSET(BULK 'https://x/y.parquet', FORMAT = 'PARQUET') WITH (a INT, b VARCHAR(10) COLLATE Latin1_General_100_BIN2_UTF8) AS r",
	"SELECT * FROM OPENQUERY(lnk, 'SELECT 1') AS q",
	"SELECT * FROM OPENDATASOURCE('SQLNCLI', 'Data Source=s;Integrated Security=SSPI').db.dbo.t AS x",
	"SELECT * FROM OPENXML(@idoc, '/ROOT/Customer', 1) WITH (CustomerID VARCHAR(10), ContactName VARCHAR(20)) AS x",
	"SELECT * FROM CHANGETABLE(CHANGES dbo.t, @last) AS ct",
	"SELECT * FROM CONTAINSTABLE(dbo.t, (a, b), 'term', LANGUAGE 1033, 10) AS k",
	"SELECT * FROM FREETEXTTABLE(dbo.t, *, 'term') AS k",
	"SELECT * FROM SEMANTICKEYPHRASETABLE(dbo.t, *) AS k",
	"SELECT value FROM STRING_SPLIT('a,b', ',', 1) AS s",
	"SELECT value FROM GENERATE_SERIES(1, 10, 2) AS g",
	"SELECT * FROM dbo.fn_tvf(1, 'a') AS f CROSS APPLY dbo.fn_other(f.id) AS o",
	"SELECT * FROM (VALUES (1, 'a'), (2, 'b')) AS v(id, name)",
	"SELECT * FROM t PIVOT (SUM(v) FOR k IN ([a], [b])) AS p UNPIVOT (v FOR k IN (a, b)) AS u",
}

// TestCompletionFinishesAtEveryOffset runs completion at every cursor offset
// of tableSourceForms, with the text after the cursor and without it, so a
// loop that stops advancing at some cursor fails here instead of pinning a
// worker. PREDICT's WITH list did, at 76 of the first statement's offsets.
func TestCompletionFinishesAtEveryOffset(t *testing.T) {
	for _, sql := range tableSourceForms {
		for off := 0; off <= len(sql); off++ {
			done := make(chan struct{})
			go func() {
				defer close(done)
				Complete(sql, off, nil)
				Complete(sql[:off], off, nil)
				parser.CollectCompletion(sql, off)
			}()
			select {
			case <-done:
			case <-time.After(10 * time.Second):
				t.Fatalf("completion at %q|%q: no result within 10s", sql[:off], sql[off:])
			}
		}
	}
}
