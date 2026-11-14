package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

func TestParseCopyFiles(t *testing.T) {
	in := "COPY FILES INTO @trg_stage/trg_path/ FROM @db.s.src_stage/src_path/ " +
		"FILES = ('file1.csv', 'file2.csv') PATTERN = '.*[.]csv' DETAILED_OUTPUT = TRUE"
	s, ok := mustParseOne(t, in).(*ast.CopyFilesStmt)
	if !ok {
		t.Fatalf("got %T", mustParseOne(t, in))
	}
	if s.Into.Raw != "@trg_stage/trg_path/" || s.From == nil || s.From.Raw != "@db.s.src_stage/src_path/" || s.FromQuery != nil {
		t.Errorf("got %+v", s)
	}
	if len(s.Options) != 3 || len(s.Options[0].List) != 2 || s.Options[2].Words != "TRUE" {
		t.Errorf("options = %+v", s.Options)
	}
	if s.Loc.Start != 0 || s.Loc.End != len(in) {
		t.Errorf("Loc = %+v", s.Loc)
	}

	q := mustParseOne(t, "COPY FILES INTO @trg FROM (SELECT src_file, trg_file FROM urls) DETAILED_OUTPUT = FALSE").(*ast.CopyFilesStmt)
	if q.FromQuery == nil || q.From != nil || len(q.Options) != 1 {
		t.Errorf("got %+v", q)
	}
	mustParseOne(t, "COPY FILES INTO @trg FROM (SELECT '@src_stage/file.txt', 'new_filename.txt')")

	for _, in := range []string{
		"COPY FILES INTO 's3://b/' FROM @b", // stages only
		"COPY FILES INTO @a FROM 's3://b/'",
		"COPY FILES INTO @a FROM @b FILES = 'x'", // a list of strings
		"COPY FILES INTO @a FROM @b FILES = ()",
		"COPY FILES INTO @a FROM @b DETAILED_OUTPUT = 'TRUE'",
		"COPY FILES INTO @a FROM @b PATTERN = 'a' PATTERN = 'b'",
		"COPY FILES INTO @a FROM @b ON_ERROR = CONTINUE",
		"COPY FILES INTO @a FROM (SELECT 'x') PATTERN = 'p'", // a query copy takes DETAILED_OUTPUT alone
		"COPY FILES INTO @a FROM (VALUES (1))",
		"COPY FILES INTO @a",
	} {
		mustRejectObject(t, in)
	}
}

// COPY's copy options and the file format options may be separated by commas.
func TestParseCopyAndFileFormatOptionCommas(t *testing.T) {
	c := mustCopyTable(t, "COPY INTO t FROM @s ON_ERROR = CONTINUE, PURGE = TRUE FORCE = FALSE")
	if len(c.Options) != 3 {
		t.Errorf("options = %+v", c.Options)
	}
	l := mustCopyLocation(t, "COPY INTO @s FROM t OVERWRITE = TRUE, SINGLE = TRUE HEADER = TRUE")
	if len(l.Options) != 3 {
		t.Errorf("options = %+v", l.Options)
	}
	f := mustCreateFileFormat(t, "CREATE FILE FORMAT f TYPE = CSV COMPRESSION = AUTO, TRIM_SPACE = TRUE")
	if len(f.Options) != 2 {
		t.Errorf("options = %+v", f.Options)
	}
	a := mustAlterFileFormat(t, "ALTER FILE FORMAT f SET COMPRESSION = AUTO, TRIM_SPACE = TRUE")
	if len(a.Options) != 2 {
		t.Errorf("options = %+v", a.Options)
	}
	for _, in := range []string{
		"COPY INTO t FROM @s ON_ERROR = CONTINUE,",
		"COPY INTO t FROM @s, ON_ERROR = CONTINUE",
		"COPY INTO t FROM @s ON_ERROR = CONTINUE,, PURGE = TRUE",
		"CREATE FILE FORMAT f TYPE = CSV COMPRESSION = AUTO,",
		"ALTER FILE FORMAT f SET COMPRESSION = AUTO,",
	} {
		mustRejectObject(t, in)
	}
}
