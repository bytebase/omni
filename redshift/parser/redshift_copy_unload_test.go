package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/redshift/ast"
)

func TestRedshiftCopyOptionsParse(t *testing.T) {
	tests := []struct {
		name        string
		sql         string
		optionName  string
		optionValue string
	}{
		{
			name:        "iam role default",
			sql:         "COPY customer FROM 's3://bucket/customer.tbl' IAM_ROLE DEFAULT;",
			optionName:  "iam_role",
			optionValue: "default",
		},
		{
			name:        "iam role arn with csv",
			sql:         "COPY customer FROM 's3://bucket/customer.tbl' IAM_ROLE 'arn:aws:iam::123456789012:role/r' FORMAT AS CSV;",
			optionName:  "format",
			optionValue: "csv",
		},
		{
			name:        "credentials",
			sql:         "COPY customer FROM 's3://bucket/customer.tbl' CREDENTIALS 'aws_access_key_id=abc';",
			optionName:  "credentials",
			optionValue: "aws_access_key_id=abc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := firstCopyStmt(t, tt.sql)
			assertDefElemString(t, stmt.Options, tt.optionName, tt.optionValue)
		})
	}
}

func TestRedshiftCopyDocumentedOptionsParse(t *testing.T) {
	tests := []struct {
		name        string
		sql         string
		optionName  string
		optionValue string
		optionInt   int64
	}{
		{
			name: "s3 load conversion and load options",
			sql: `COPY listing FROM 's3://bucket/listing/'
IAM_ROLE 'arn:aws:iam::123456789012:role/r'
DELIMITER '|'
IGNOREHEADER 1
DATEFORMAT 'auto'
TIMEFORMAT 'auto'
MAXERROR 10
ACCEPTINVCHARS
GZIP;`,
			optionName: "ignoreheader",
			optionInt:  1,
		},
		{
			name: "manifest region and bzip2",
			sql: `COPY sales FROM 's3://bucket/sales/'
IAM_ROLE 'arn:aws:iam::123456789012:role/r'
MANIFEST
REGION 'eu-west-1'
DELIMITER '\t'
BZIP2;`,
			optionName:  "region",
			optionValue: "eu-west-1",
		},
		{
			name: "dynamodb readratio and stats",
			sql: `COPY myredshifttable (col1, col2)
FROM 'dynamodb://ProductCatalog'
IAM_ROLE 'arn:aws:iam::123456789012:role/r'
READRATIO 25
COMPUPDATE OFF
STATUPDATE ON;`,
			optionName: "readratio",
			optionInt:  25,
		},
		{
			name: "json avro fixedwidth and boolean flags",
			sql: `COPY event FROM 's3://bucket/event/'
IAM_ROLE 'arn:aws:iam::123456789012:role/r'
JSON 'auto'
FIXEDWIDTH 'id:3,name:25'
ACCEPTINVCHARS '?'
NULL AS '\000'
ESCAPE
ROUNDEC
TRIMBLANKS
TRUNCATECOLUMNS
IGNOREBLANKLINES
ACCEPTANYDATE
EMPTYASNULL
BLANKSASNULL
FILLRECORD
REMOVEQUOTES;`,
			optionName:  "json",
			optionValue: "auto",
		},
		{
			name: "format as csv with optional as values",
			sql: `COPY venue FROM 's3://bucket/venue.csv'
IAM_ROLE 'arn:aws:iam::123456789012:role/r'
FORMAT AS CSV
QUOTE AS '"'
IGNOREHEADER AS 1
DATEFORMAT AS 'MM/DD/YYYY'
REGION AS 'us-west-2';`,
			optionName:  "quote",
			optionValue: `"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := firstCopyStmt(t, tt.sql)
			if tt.optionValue != "" {
				assertDefElemString(t, stmt.Options, tt.optionName, tt.optionValue)
			} else {
				assertDefElemInt(t, stmt.Options, tt.optionName, tt.optionInt)
			}
		})
	}
}

func TestRedshiftUnloadOptionsParse(t *testing.T) {
	tests := []struct {
		name        string
		sql         string
		optionName  string
		optionValue string
	}{
		{
			name:        "iam role default",
			sql:         "UNLOAD ('SELECT * FROM sales') TO 's3://bucket/out' IAM_ROLE DEFAULT;",
			optionName:  "iam_role",
			optionValue: "default",
		},
		{
			name:        "csv header manifest",
			sql:         "UNLOAD ('SELECT * FROM sales') TO 's3://bucket/out' IAM_ROLE DEFAULT FORMAT AS CSV HEADER MANIFEST;",
			optionName:  "format",
			optionValue: "csv",
		},
		{
			name:        "delimiter encrypted gzip",
			sql:         "UNLOAD ('SELECT * FROM sales') TO 's3://bucket/out' IAM_ROLE DEFAULT DELIMITER AS '|' ENCRYPTED GZIP;",
			optionName:  "delimiter",
			optionValue: "|",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := firstUnloadStmt(t, tt.sql)
			assertDefElemString(t, stmt.Options, tt.optionName, tt.optionValue)
		})
	}
}

func TestRedshiftUnloadDocumentedOptionsParse(t *testing.T) {
	tests := []struct {
		name        string
		sql         string
		optionName  string
		optionValue string
	}{
		{
			name: "partition include and parquet sizes",
			sql: `UNLOAD ('SELECT * FROM fact_sales')
TO 's3://bucket/fact_sales/'
IAM_ROLE 'arn:aws:iam::123456789012:role/r'
FORMAT PARQUET
PARTITION BY (year, month, day) INCLUDE
ROWGROUPSIZE 512 MB
MAXFILESIZE 1 GB
MANIFEST VERBOSE
CLEANPATH;`,
			optionName:  "manifest",
			optionValue: "verbose",
		},
		{
			name: "encryption kms and overwrite",
			sql: `UNLOAD ('SELECT * FROM sensitive_data')
TO 's3://bucket/encrypted/'
IAM_ROLE 'arn:aws:iam::123456789012:role/r'
ENCRYPTED AUTO
KMS_KEY_ID '1234abcd-12ab-34cd-56ef-1234567890ab'
ALLOWOVERWRITE;`,
			optionName:  "encrypted",
			optionValue: "auto",
		},
		{
			name: "parallel extension addquotes null escape",
			sql: `UNLOAD ('SELECT * FROM text_data')
TO 's3://bucket/text/'
IAM_ROLE 'arn:aws:iam::123456789012:role/r'
FORMAT JSON
ADDQUOTES
NULL AS ''
ESCAPE
PARALLEL OFF
MAXFILESIZE 250 MB
ZSTD
REGION 'ap-southeast-1'
EXTENSION 'json.zst';`,
			optionName:  "parallel",
			optionValue: "off",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := firstUnloadStmt(t, tt.sql)
			assertDefElemString(t, stmt.Options, tt.optionName, tt.optionValue)
		})
	}
}

func firstCopyStmt(t *testing.T, sql string) *nodes.CopyStmt {
	t.Helper()
	tree, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if len(tree.Items) != 1 {
		t.Fatalf("expected one statement, got %d", len(tree.Items))
	}
	raw, ok := tree.Items[0].(*nodes.RawStmt)
	if !ok {
		t.Fatalf("expected RawStmt, got %T", tree.Items[0])
	}
	stmt, ok := raw.Stmt.(*nodes.CopyStmt)
	if !ok {
		t.Fatalf("expected CopyStmt, got %T", raw.Stmt)
	}
	return stmt
}

func firstUnloadStmt(t *testing.T, sql string) *nodes.UnloadStmt {
	t.Helper()
	tree, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if len(tree.Items) != 1 {
		t.Fatalf("expected one statement, got %d", len(tree.Items))
	}
	raw, ok := tree.Items[0].(*nodes.RawStmt)
	if !ok {
		t.Fatalf("expected RawStmt, got %T", tree.Items[0])
	}
	stmt, ok := raw.Stmt.(*nodes.UnloadStmt)
	if !ok {
		t.Fatalf("expected UnloadStmt, got %T", raw.Stmt)
	}
	return stmt
}
