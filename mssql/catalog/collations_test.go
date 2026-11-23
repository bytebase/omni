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
// pinned SQL Server: the file is sys.fn_helpcollations(), and SQL Server 2022
// accepts a column collation exactly when the list holds it, for every listed
// name in upper and lower case and for every version-less, versioned, _SC,
// _UTF8 and reordered form of a listed name (collationFormsSQL). The catalog
// accepts more than SQL Server 2022 does, the names another version may know
// (TestUnknownCollation).
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
		_, listed := lookupCollation(form)
		if listed = listed || isDefaultCollation(form); listed != accepted {
			t.Errorf("COLLATE %s: SQL Server accepts it %v, the list %v", form, accepted, listed)
		}
		if accepted && checkCollation(form, true) != nil {
			t.Errorf("COLLATE %s: SQL Server accepts it, the catalog refuses it", form)
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

// TestUnknownCollation holds UnknownCollation, which the catalog and
// mssql/review's WALK_THROUGH refuse a collation by, to what SQL Server
// 2017, 2019, 2022, and 2025 accept and refuse: every name SQL Server 2022
// lists, the 32 SQL Server 2025 adds, and Compatibility_52_409_30003, which
// each version takes and none lists, are known; names each refuses are not.
// The versions' lists are those measured for bytebase/omni#1705 (2017
// 14.0.3570.2 and 2019 15.0.4510.2 list subsets of 2022's names, 2025
// 17.0.5065.9 lists 2022's and these 32); only SQL Server 2022 runs here, and
// it refuses the 2025 names.
func TestUnknownCollation(t *testing.T) {
	known := make([]string, 0, len(collationNames())+33)
	for _, name := range collationNames() {
		known = append(known, name)
	}
	if len(known) < 5000 {
		t.Fatalf("read %d collations from collations.txt", len(known))
	}
	// The 32 collations SQL Server 2025 (17.0.5065.9) adds.
	var added []string
	for _, designator := range []string{"Chinese_Simplified_Pinyin_160", "Chinese_Simplified_Stroke_Order_160"} {
		for _, cs := range []string{"_CI", "_CS"} {
			for _, as := range []string{"_AI", "_AS"} {
				for _, ks := range []string{"", "_KS"} {
					for _, ws := range []string{"", "_WS"} {
						added = append(added, designator+cs+as+ks+ws+"_SC_UTF8")
					}
				}
			}
		}
	}
	for _, name := range added {
		if _, listed := lookupCollation(name); listed {
			t.Errorf("collations.txt lists %s, which SQL Server 2022 refuses", name)
		}
	}
	known = append(known, added...)
	known = append(known, "Compatibility_52_409_30003")
	for _, name := range known {
		for _, form := range []string{name, strings.ToLower(name), strings.ToUpper(name)} {
			for _, defaultAllowed := range []bool{true, false} {
				if UnknownCollation(form, defaultAllowed) {
					t.Errorf("UnknownCollation(%q, %v) = true, want false", form, defaultAllowed)
				}
			}
			if err := checkCollation(form, false); err != nil {
				t.Errorf("checkCollation(%q) = %v, want nil", form, err)
			}
		}
	}
	// Refused by SQL Server 2017, 2019, 2022, and 2025.
	for _, name := range []string{
		"a", "nope", "Latin1_General", "Latin1_General_CS",
		"Latin1_General_CI_AS_WS_KS", "Latin1_General_CI_AS_VSS_WS", "Latin1_General_CI_AS_KS_KS",
	} {
		if !UnknownCollation(name, true) {
			t.Errorf("UnknownCollation(%q) = false, want true", name)
		}
		if err := checkCollation(name, true); err == nil {
			t.Errorf("checkCollation(%q) = nil, want Msg 448", name)
		}
	}
	// database_default and catalog_default name a column's collation, not a
	// database's.
	for _, name := range []string{"database_default", "CATALOG_DEFAULT"} {
		if UnknownCollation(name, true) || !UnknownCollation(name, false) {
			t.Errorf("UnknownCollation(%q, true/false) = %v/%v, want false/true", name, UnknownCollation(name, true), UnknownCollation(name, false))
		}
	}
	// SQL Server stores a compatibility collation's prefix in its own case.
	if got := storedCollation("compatibility_52_409_30003"); got != "Compatibility_52_409_30003" {
		t.Errorf("storedCollation(compatibility_52_409_30003) = %q", got)
	}
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
