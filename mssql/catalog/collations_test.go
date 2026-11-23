package catalog

import (
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/mssql/internal/testcontainer"
)

// TestCollations holds collations.txt, which lookupCollation reads, to the
// pinned SQL Server: the file is sys.fn_helpcollations(), and SQL Server
// accepts a column collation exactly when checkCollation does, for every
// listed name in upper and lower case and for every version-less, versioned,
// _SC, _UTF8 and reordered form of a listed name (collationFormsSQL).
//
// SQL Server's answers are on the test's tape. OMNI_RECORD=1 asks the pinned
// server again and rewrites collations.txt from its list.
func TestCollations(t *testing.T) {
	db := testcontainer.Open(t)

	var list string
	if err := db.QueryRow(listCollationsSQL).Scan(&list); err != nil {
		t.Fatal(err)
	}
	names := strings.Split(list, "\n")
	slices.Sort(names)
	want := collationsHeader + strings.Join(names, "\n") + "\n"
	mode, err := enginetest.RunMode()
	if err != nil {
		t.Fatal(err)
	}
	if mode == enginetest.ModeRecord {
		if want != collationsText {
			if err := os.WriteFile("collations.txt", []byte(want), 0o644); err != nil {
				t.Fatal(err)
			}
			t.Logf("wrote collations.txt: %d collations", len(names))
		}
	} else if want != collationsText {
		t.Fatalf("collations.txt is not the pinned server's sys.fn_helpcollations() (%d names); rewrite it with OMNI_RECORD=1 %s", len(names), enginetest.RunCommand(t))
	}
	for _, name := range names {
		for _, form := range []string{name, strings.ToLower(name), strings.ToUpper(name)} {
			if got, ok := lookupCollation(form); !ok || got != name {
				t.Errorf("lookupCollation(%q) = %q, %v; want %q", form, got, ok, name)
			}
		}
	}

	rows, err := db.Query(collationFormsSQL)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var form string
		var accepted bool
		if err := rows.Scan(&form, &accepted); err != nil {
			t.Fatal(err)
		}
		if omni := checkCollation(form, true) == nil; omni != accepted {
			t.Errorf("COLLATE %s: SQL Server accepts it %v, the catalog %v", form, accepted, omni)
		}
	}
	if !rows.NextResultSet() || !rows.Next() {
		t.Fatalf("no count of the forms tried: %v", rows.Err())
	}
	var tried int
	if err := rows.Scan(&tried); err != nil {
		t.Fatal(err)
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	t.Logf("SQL Server and the catalog agree on %d collation forms", tried)
}

const collationsHeader = `# Every collation SQL Server 2022 knows (16.0.4305.5, mssql/testdata/pins.tsv):
# sys.fn_helpcollations() on the pinned server, sorted. Written by TestCollations;
# rewrite it with OMNI_RECORD=1 go test ./mssql/catalog/ -run '^TestCollations$'.
`

// listCollationsSQL reads sys.fn_helpcollations() as one value, so the tape
// holds one answer rather than a row per collation.
const listCollationsSQL = `SELECT STRING_AGG(CAST(name AS NVARCHAR(MAX)), NCHAR(10)) FROM sys.fn_helpcollations()`

// collationFormsSQL tries each form of each listed collation as a column's
// collation and answers the forms SQL Server judges otherwise than the
// case-insensitive list does, then how many forms it tried. A form is a listed
// name in upper or lower case; its stem (the name without _SC and _UTF8) with
// and without its version, with each of _SC, _UTF8, _SC_UTF8 and _UTF8_SC; the
// name with two of its flags swapped; or the name at the next version.
// database_default and catalog_default are tried too.
const collationFormsSQL = `SET NOCOUNT ON;
DECLARE @listed TABLE (name NVARCHAR(128) COLLATE Latin1_General_CI_AS PRIMARY KEY);
INSERT @listed SELECT name FROM sys.fn_helpcollations();
DECLARE @form TABLE (name NVARCHAR(128) COLLATE Latin1_General_BIN2 PRIMARY KEY, listed BIT, accepted BIT);
WITH n AS (SELECT name COLLATE Latin1_General_BIN2 AS name FROM @listed),
stem AS (SELECT DISTINCT REPLACE(REPLACE(name, N'_UTF8', N''), N'_SC', N'') AS stem FROM n),
base AS (
	SELECT stem FROM stem
	UNION SELECT REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(stem, N'_80_', N'_'), N'_90_', N'_'), N'_100_', N'_'), N'_140_', N'_'), N'_160_', N'_') FROM stem
),
f AS (
	SELECT stem + tail AS name FROM base CROSS JOIN (VALUES (N''), (N'_SC'), (N'_UTF8'), (N'_SC_UTF8'), (N'_UTF8_SC')) x(tail)
	UNION SELECT UPPER(name) FROM n
	UNION SELECT LOWER(name) FROM n
	UNION SELECT REPLACE(name, x.f, x.t) FROM n CROSS JOIN (VALUES
		(N'_CI_AS', N'_AS_CI'), (N'_CI_AI', N'_AI_CI'), (N'_CS_AS', N'_AS_CS'), (N'_CS_AI', N'_AI_CS'),
		(N'_KS_WS', N'_WS_KS'), (N'_WS_VSS', N'_VSS_WS'), (N'_KS_VSS', N'_VSS_KS'),
		(N'_80_', N'_90_'), (N'_90_', N'_100_'), (N'_100_', N'_140_'), (N'_140_', N'_160_')) x(f, t)
	WHERE CHARINDEX(x.f, name) > 0
	UNION SELECT N'database_default'
	UNION SELECT N'CATALOG_DEFAULT'
)
INSERT @form (name, listed)
SELECT f.name, CASE WHEN l.name IS NULL THEN 0 ELSE 1 END
FROM f LEFT JOIN @listed l ON l.name = f.name COLLATE Latin1_General_CI_AS;
DECLARE @name NVARCHAR(128), @sql NVARCHAR(400);
DECLARE forms CURSOR LOCAL FAST_FORWARD FOR SELECT name FROM @form;
OPEN forms;
FETCH forms INTO @name;
WHILE @@FETCH_STATUS = 0
BEGIN
	SET @sql = N'DECLARE @t TABLE (c NVARCHAR(1) COLLATE ' + @name + N')';
	BEGIN TRY
		EXEC sp_executesql @sql;
		UPDATE @form SET accepted = 1 WHERE name = @name;
	END TRY
	BEGIN CATCH
		IF ERROR_NUMBER() <> 448 THROW;
		UPDATE @form SET accepted = 0 WHERE name = @name;
	END CATCH;
	FETCH forms INTO @name;
END;
CLOSE forms;
DEALLOCATE forms;
SELECT name, accepted FROM @form WHERE accepted <> listed ORDER BY name;
SELECT COUNT(*) FROM @form;`
