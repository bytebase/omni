package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// Test helpers
// ---------------------------------------------------------------------------

func mustCreateStage(t *testing.T, input string) *ast.CreateStageStmt {
	t.Helper()
	node := mustParseOne(t, input)
	stmt, ok := node.(*ast.CreateStageStmt)
	if !ok {
		t.Fatalf("parse %q: got %T, want *ast.CreateStageStmt", input, node)
	}
	return stmt
}

func mustAlterStage(t *testing.T, input string) *ast.AlterStageStmt {
	t.Helper()
	node := mustParseOne(t, input)
	stmt, ok := node.(*ast.AlterStageStmt)
	if !ok {
		t.Fatalf("parse %q: got %T, want *ast.AlterStageStmt", input, node)
	}
	return stmt
}

// groupEntry returns the named entry of a parenthesized option group, or nil.
func groupEntry(opt *ast.CopyOption, name string) *ast.CopyOption {
	if opt == nil {
		return nil
	}
	return findOption(opt.Group, name)
}

// ---------------------------------------------------------------------------
// CREATE STAGE — modifiers, name, IF NOT EXISTS
// ---------------------------------------------------------------------------

func TestParseCreateStage_Modifiers(t *testing.T) {
	t.Run("minimal internal", func(t *testing.T) {
		stmt := mustCreateStage(t, "CREATE STAGE my_stage")
		if stmt.Name.String() != "my_stage" {
			t.Errorf("Name = %q, want my_stage", stmt.Name.String())
		}
		if stmt.OrReplace || stmt.Temporary || stmt.IfNotExists {
			t.Errorf("unexpected modifier set: %+v", stmt)
		}
		if len(stmt.Options) != 0 || len(stmt.Tags) != 0 {
			t.Errorf("expected no options/tags, got %+v / %+v", stmt.Options, stmt.Tags)
		}
	})

	t.Run("or replace", func(t *testing.T) {
		stmt := mustCreateStage(t, "CREATE OR REPLACE STAGE s")
		if !stmt.OrReplace {
			t.Error("OrReplace not set")
		}
	})

	t.Run("temporary", func(t *testing.T) {
		stmt := mustCreateStage(t, "CREATE TEMPORARY STAGE s")
		if !stmt.Temporary {
			t.Error("Temporary not set")
		}
	})

	t.Run("temp", func(t *testing.T) {
		stmt := mustCreateStage(t, "CREATE TEMP STAGE s")
		if !stmt.Temporary {
			t.Error("Temporary not set for TEMP")
		}
	})

	t.Run("or replace temporary", func(t *testing.T) {
		stmt := mustCreateStage(t, "CREATE OR REPLACE TEMPORARY STAGE s")
		if !stmt.OrReplace || !stmt.Temporary {
			t.Errorf("OrReplace=%v Temporary=%v, want both true", stmt.OrReplace, stmt.Temporary)
		}
	})

	t.Run("if not exists", func(t *testing.T) {
		stmt := mustCreateStage(t, "CREATE STAGE IF NOT EXISTS s")
		if !stmt.IfNotExists {
			t.Error("IfNotExists not set")
		}
	})

	t.Run("qualified name", func(t *testing.T) {
		stmt := mustCreateStage(t, "CREATE STAGE db.sch.s")
		if stmt.Name.String() != "db.sch.s" {
			t.Errorf("Name = %q, want db.sch.s", stmt.Name.String())
		}
	})

	t.Run("quoted name", func(t *testing.T) {
		stmt := mustCreateStage(t, `CREATE STAGE "My Stage"`)
		if stmt.Name.String() != `"My Stage"` {
			t.Errorf("Name = %q", stmt.Name.String())
		}
	})
}

// ---------------------------------------------------------------------------
// CREATE STAGE — internal stage params (ENCRYPTION / DIRECTORY / FILE_FORMAT /
// COPY_OPTIONS / COMMENT / WITH TAG)
// ---------------------------------------------------------------------------

func TestParseCreateStage_InternalParams(t *testing.T) {
	t.Run("encryption snowflake_sse", func(t *testing.T) {
		stmt := mustCreateStage(t, "CREATE STAGE s ENCRYPTION = (TYPE = 'SNOWFLAKE_SSE')")
		enc := findOption(stmt.Options, "ENCRYPTION")
		if enc == nil {
			t.Fatalf("missing ENCRYPTION; options=%+v", stmt.Options)
		}
		typ := groupEntry(enc, "TYPE")
		if typ == nil || typ.Lit == nil || typ.Lit.Value != "SNOWFLAKE_SSE" {
			t.Errorf("ENCRYPTION TYPE = %+v, want 'SNOWFLAKE_SSE'", typ)
		}
	})

	t.Run("encryption snowflake_full", func(t *testing.T) {
		stmt := mustCreateStage(t, "CREATE STAGE s ENCRYPTION = (TYPE = 'SNOWFLAKE_FULL')")
		if groupEntry(findOption(stmt.Options, "ENCRYPTION"), "TYPE") == nil {
			t.Error("missing ENCRYPTION TYPE")
		}
	})

	t.Run("directory enable", func(t *testing.T) {
		stmt := mustCreateStage(t, "CREATE STAGE s DIRECTORY = (ENABLE = TRUE)")
		dir := findOption(stmt.Options, "DIRECTORY")
		en := groupEntry(dir, "ENABLE")
		if en == nil || en.Words != "TRUE" {
			t.Errorf("DIRECTORY ENABLE = %+v, want Words=TRUE", en)
		}
	})

	t.Run("directory enable auto_refresh (internal)", func(t *testing.T) {
		stmt := mustCreateStage(t, "CREATE STAGE s DIRECTORY = (ENABLE = TRUE AUTO_REFRESH = FALSE)")
		dir := findOption(stmt.Options, "DIRECTORY")
		if groupEntry(dir, "ENABLE") == nil || groupEntry(dir, "AUTO_REFRESH") == nil {
			t.Errorf("DIRECTORY group = %+v", dir)
		}
	})

	t.Run("file_format format_name", func(t *testing.T) {
		stmt := mustCreateStage(t, "CREATE STAGE s FILE_FORMAT = (FORMAT_NAME = 'my_csv')")
		ff := findOption(stmt.Options, "FILE_FORMAT")
		fn := groupEntry(ff, "FORMAT_NAME")
		if fn == nil || fn.Lit == nil || fn.Lit.Value != "my_csv" {
			t.Errorf("FORMAT_NAME = %+v, want 'my_csv'", fn)
		}
	})

	t.Run("file_format type with options", func(t *testing.T) {
		stmt := mustCreateStage(t, "CREATE STAGE s FILE_FORMAT = (TYPE = CSV SKIP_HEADER = 1 FIELD_DELIMITER = ',')")
		ff := findOption(stmt.Options, "FILE_FORMAT")
		if ff == nil || len(ff.Group) != 3 {
			t.Fatalf("FILE_FORMAT group = %+v", ff)
		}
		typ := groupEntry(ff, "TYPE")
		if typ == nil || typ.Words != "CSV" {
			t.Errorf("TYPE = %+v, want Words=CSV", typ)
		}
		if groupEntry(ff, "SKIP_HEADER") == nil || groupEntry(ff, "FIELD_DELIMITER") == nil {
			t.Errorf("missing SKIP_HEADER/FIELD_DELIMITER in %+v", ff)
		}
	})

	t.Run("copy_options", func(t *testing.T) {
		stmt := mustCreateStage(t, "CREATE STAGE s COPY_OPTIONS = (ON_ERROR = 'SKIP_FILE' PURGE = TRUE)")
		co := findOption(stmt.Options, "COPY_OPTIONS")
		if co == nil || len(co.Group) != 2 {
			t.Fatalf("COPY_OPTIONS group = %+v", co)
		}
		if groupEntry(co, "ON_ERROR") == nil || groupEntry(co, "PURGE") == nil {
			t.Errorf("COPY_OPTIONS entries = %+v", co.Group)
		}
	})

	t.Run("comment", func(t *testing.T) {
		stmt := mustCreateStage(t, "CREATE STAGE s COMMENT = 'a stage'")
		c := findOption(stmt.Options, "COMMENT")
		if c == nil || c.Lit == nil || c.Lit.Value != "a stage" {
			t.Errorf("COMMENT = %+v, want 'a stage'", c)
		}
	})

	t.Run("with tag", func(t *testing.T) {
		stmt := mustCreateStage(t, "CREATE STAGE s WITH TAG (cost_center = 'sales')")
		if len(stmt.Tags) != 1 {
			t.Fatalf("Tags = %+v, want 1", stmt.Tags)
		}
		if stmt.Tags[0].Name.String() != "cost_center" || stmt.Tags[0].Value != "sales" {
			t.Errorf("tag = %+v", stmt.Tags[0])
		}
	})

	t.Run("bare tag", func(t *testing.T) {
		stmt := mustCreateStage(t, "CREATE STAGE s TAG (a = '1', b = '2')")
		if len(stmt.Tags) != 2 {
			t.Fatalf("Tags = %+v, want 2", stmt.Tags)
		}
	})

	t.Run("comment then with tag (docs order)", func(t *testing.T) {
		stmt := mustCreateStage(t, "CREATE STAGE s COMMENT = 'c' WITH TAG (t = 'v')")
		if findOption(stmt.Options, "COMMENT") == nil {
			t.Error("missing COMMENT")
		}
		if len(stmt.Tags) != 1 {
			t.Errorf("Tags = %+v, want 1", stmt.Tags)
		}
	})

	// The legacy ANTLR grammar lists `with_tags? comment_clause?` (TAG before
	// COMMENT); accepting that order too avoids silently dropping a trailing
	// COMMENT and regresses neither order. Both spellings must keep the COMMENT.
	t.Run("with tag then comment (antlr order)", func(t *testing.T) {
		stmt := mustCreateStage(t, "CREATE STAGE s WITH TAG (t = 'v') COMMENT = 'c'")
		c := findOption(stmt.Options, "COMMENT")
		if c == nil || c.Lit == nil || c.Lit.Value != "c" {
			t.Errorf("COMMENT not captured after TAG: %+v", c)
		}
		if len(stmt.Tags) != 1 {
			t.Errorf("Tags = %+v, want 1", stmt.Tags)
		}
	})

	t.Run("bare tag then comment (antlr order)", func(t *testing.T) {
		stmt := mustCreateStage(t, "CREATE STAGE s TAG (t = 'v') COMMENT = 'c'")
		if findOption(stmt.Options, "COMMENT") == nil {
			t.Error("COMMENT dropped after bare TAG")
		}
		if len(stmt.Tags) != 1 {
			t.Errorf("Tags = %+v, want 1", stmt.Tags)
		}
	})
}

// ---------------------------------------------------------------------------
// CREATE STAGE — external stage params (S3 / GCS / Azure)
// ---------------------------------------------------------------------------

func TestParseCreateStage_ExternalS3(t *testing.T) {
	t.Run("url + storage_integration", func(t *testing.T) {
		stmt := mustCreateStage(t, "CREATE STAGE s URL = 's3://mybucket/path/' STORAGE_INTEGRATION = my_int")
		url := findOption(stmt.Options, "URL")
		if url == nil || url.Lit == nil || url.Lit.Value != "s3://mybucket/path/" {
			t.Errorf("URL = %+v", url)
		}
		si := findOption(stmt.Options, "STORAGE_INTEGRATION")
		if si == nil || si.Words != "MY_INT" {
			t.Errorf("STORAGE_INTEGRATION = %+v, want Words=MY_INT", si)
		}
	})

	t.Run("credentials aws key/secret/token", func(t *testing.T) {
		stmt := mustCreateStage(t, "CREATE STAGE s URL = 's3://b/p' CREDENTIALS = (AWS_KEY_ID = 'k' AWS_SECRET_KEY = 's' AWS_TOKEN = 't')")
		cr := findOption(stmt.Options, "CREDENTIALS")
		if cr == nil || len(cr.Group) != 3 {
			t.Fatalf("CREDENTIALS = %+v", cr)
		}
		for _, k := range []string{"AWS_KEY_ID", "AWS_SECRET_KEY", "AWS_TOKEN"} {
			if groupEntry(cr, k) == nil {
				t.Errorf("missing %s in %+v", k, cr.Group)
			}
		}
	})

	t.Run("credentials aws_role", func(t *testing.T) {
		stmt := mustCreateStage(t, "CREATE STAGE s URL = 's3://b/p' CREDENTIALS = (AWS_ROLE = 'arn:aws:iam::123:role/r')")
		cr := findOption(stmt.Options, "CREDENTIALS")
		if groupEntry(cr, "AWS_ROLE") == nil {
			t.Errorf("missing AWS_ROLE in %+v", cr)
		}
	})

	t.Run("encryption aws_sse_kms with kms_key_id", func(t *testing.T) {
		stmt := mustCreateStage(t, "CREATE STAGE s URL = 's3://b/p' ENCRYPTION = (TYPE = 'AWS_SSE_KMS' KMS_KEY_ID = 'aws-key')")
		enc := findOption(stmt.Options, "ENCRYPTION")
		if groupEntry(enc, "TYPE") == nil || groupEntry(enc, "KMS_KEY_ID") == nil {
			t.Errorf("ENCRYPTION = %+v", enc)
		}
	})

	t.Run("encryption aws_cse with master_key", func(t *testing.T) {
		stmt := mustCreateStage(t, "CREATE STAGE s URL = 's3://b/p' ENCRYPTION = (TYPE = 'AWS_CSE' MASTER_KEY = 'm')")
		enc := findOption(stmt.Options, "ENCRYPTION")
		if groupEntry(enc, "MASTER_KEY") == nil {
			t.Errorf("missing MASTER_KEY in %+v", enc)
		}
	})

	t.Run("storage_integration and encryption combined", func(t *testing.T) {
		stmt := mustCreateStage(t, "CREATE STAGE s URL = 's3://b/p' STORAGE_INTEGRATION = i ENCRYPTION = (TYPE = 'AWS_SSE_S3')")
		if findOption(stmt.Options, "URL") == nil ||
			findOption(stmt.Options, "STORAGE_INTEGRATION") == nil ||
			findOption(stmt.Options, "ENCRYPTION") == nil {
			t.Errorf("options = %+v", stmt.Options)
		}
	})
}

func TestParseCreateStage_ExternalGCS(t *testing.T) {
	t.Run("url + storage_integration + encryption gcs_sse_kms", func(t *testing.T) {
		stmt := mustCreateStage(t, "CREATE STAGE s URL = 'gcs://mybucket/path/' STORAGE_INTEGRATION = gcs_int ENCRYPTION = (TYPE = 'GCS_SSE_KMS' KMS_KEY_ID = 'k')")
		url := findOption(stmt.Options, "URL")
		if url == nil || url.Lit.Value != "gcs://mybucket/path/" {
			t.Errorf("URL = %+v", url)
		}
		enc := findOption(stmt.Options, "ENCRYPTION")
		if groupEntry(enc, "TYPE") == nil || groupEntry(enc, "KMS_KEY_ID") == nil {
			t.Errorf("ENCRYPTION = %+v", enc)
		}
	})
}

func TestParseCreateStage_ExternalAzure(t *testing.T) {
	t.Run("url + credentials azure_sas_token + encryption azure_cse", func(t *testing.T) {
		stmt := mustCreateStage(t, "CREATE STAGE s URL = 'azure://acct.blob.core.windows.net/cont/path/' CREDENTIALS = (AZURE_SAS_TOKEN = 'sas') ENCRYPTION = (TYPE = 'AZURE_CSE' MASTER_KEY = 'm')")
		if findOption(stmt.Options, "URL") == nil {
			t.Error("missing URL")
		}
		cr := findOption(stmt.Options, "CREDENTIALS")
		if groupEntry(cr, "AZURE_SAS_TOKEN") == nil {
			t.Errorf("missing AZURE_SAS_TOKEN in %+v", cr)
		}
		enc := findOption(stmt.Options, "ENCRYPTION")
		if groupEntry(enc, "MASTER_KEY") == nil {
			t.Errorf("missing MASTER_KEY in %+v", enc)
		}
	})
}

func TestParseCreateStage_ExternalDirectory(t *testing.T) {
	// External DIRECTORY supports REFRESH_ON_CREATE / AUTO_REFRESH /
	// NOTIFICATION_INTEGRATION beyond the internal ENABLE/AUTO_REFRESH.
	t.Run("directory full external", func(t *testing.T) {
		stmt := mustCreateStage(t, "CREATE STAGE s URL = 's3://b/p' STORAGE_INTEGRATION = i DIRECTORY = (ENABLE = TRUE REFRESH_ON_CREATE = TRUE AUTO_REFRESH = TRUE NOTIFICATION_INTEGRATION = 'ni')")
		dir := findOption(stmt.Options, "DIRECTORY")
		if dir == nil || len(dir.Group) != 4 {
			t.Fatalf("DIRECTORY group = %+v", dir)
		}
		for _, k := range []string{"ENABLE", "REFRESH_ON_CREATE", "AUTO_REFRESH", "NOTIFICATION_INTEGRATION"} {
			if groupEntry(dir, k) == nil {
				t.Errorf("missing %s in %+v", k, dir.Group)
			}
		}
	})
}

func TestParseCreateStage_ExternalFull(t *testing.T) {
	// A full external S3 stage exercising every option block + COMMENT + WITH TAG
	// in one statement, in the documented order.
	input := "CREATE OR REPLACE STAGE db.sch.my_ext_stage " +
		"URL = 's3://load/files/' " +
		"STORAGE_INTEGRATION = s3_int " +
		"ENCRYPTION = (TYPE = 'AWS_SSE_KMS' KMS_KEY_ID = 'aws-key') " +
		"DIRECTORY = (ENABLE = TRUE AUTO_REFRESH = TRUE) " +
		"FILE_FORMAT = (TYPE = JSON STRIP_OUTER_ARRAY = TRUE) " +
		"COPY_OPTIONS = (ON_ERROR = 'CONTINUE') " +
		"COMMENT = 'prod load stage' " +
		"WITH TAG (env = 'prod', team = 'data')"
	stmt := mustCreateStage(t, input)
	if !stmt.OrReplace {
		t.Error("OrReplace not set")
	}
	if stmt.Name.String() != "db.sch.my_ext_stage" {
		t.Errorf("Name = %q", stmt.Name.String())
	}
	for _, k := range []string{"URL", "STORAGE_INTEGRATION", "ENCRYPTION", "DIRECTORY", "FILE_FORMAT", "COPY_OPTIONS", "COMMENT"} {
		if findOption(stmt.Options, k) == nil {
			t.Errorf("missing option %s", k)
		}
	}
	if len(stmt.Tags) != 2 {
		t.Errorf("Tags = %+v, want 2", stmt.Tags)
	}
}

func TestParseCreateStage_DocsNewOptions(t *testing.T) {
	// Options present in the official docs but absent from the legacy ANTLR
	// grammar. The open-ended option parser must accept them with no special
	// casing. (Divergence: docs add these; antlr lacks them — docs win.)
	t.Run("aws_access_point_arn", func(t *testing.T) {
		stmt := mustCreateStage(t, "CREATE STAGE s URL = 's3://b/p' AWS_ACCESS_POINT_ARN = 'arn:aws:s3:::ap'")
		if findOption(stmt.Options, "AWS_ACCESS_POINT_ARN") == nil {
			t.Errorf("missing AWS_ACCESS_POINT_ARN; options=%+v", stmt.Options)
		}
	})

	t.Run("use_privatelink_endpoint", func(t *testing.T) {
		stmt := mustCreateStage(t, "CREATE STAGE s URL = 's3://b/p' USE_PRIVATELINK_ENDPOINT = TRUE")
		opt := findOption(stmt.Options, "USE_PRIVATELINK_ENDPOINT")
		if opt == nil || opt.Words != "TRUE" {
			t.Errorf("USE_PRIVATELINK_ENDPOINT = %+v", opt)
		}
	})

	t.Run("s3compat endpoint", func(t *testing.T) {
		stmt := mustCreateStage(t, "CREATE STAGE s URL = 's3compat://bucket/path/' ENDPOINT = 'my.endpoint.com' CREDENTIALS = (AWS_KEY_ID = 'k' AWS_SECRET_KEY = 's')")
		if findOption(stmt.Options, "ENDPOINT") == nil {
			t.Errorf("missing ENDPOINT; options=%+v", stmt.Options)
		}
	})
}

// ---------------------------------------------------------------------------
// ALTER STAGE
// ---------------------------------------------------------------------------

func TestParseAlterStage_Rename(t *testing.T) {
	t.Run("rename", func(t *testing.T) {
		stmt := mustAlterStage(t, "ALTER STAGE s RENAME TO s2")
		if stmt.Action != ast.AlterStageRename {
			t.Fatalf("Action = %v, want Rename", stmt.Action)
		}
		if stmt.Name.String() != "s" || stmt.NewName.String() != "s2" {
			t.Errorf("Name=%q NewName=%q", stmt.Name.String(), stmt.NewName.String())
		}
	})

	t.Run("if exists rename qualified", func(t *testing.T) {
		stmt := mustAlterStage(t, "ALTER STAGE IF EXISTS db.sch.s RENAME TO db.sch.s2")
		if !stmt.IfExists {
			t.Error("IfExists not set")
		}
		if stmt.NewName.String() != "db.sch.s2" {
			t.Errorf("NewName = %q", stmt.NewName.String())
		}
	})
}

func TestParseAlterStage_Set(t *testing.T) {
	t.Run("set file_format", func(t *testing.T) {
		stmt := mustAlterStage(t, "ALTER STAGE s SET FILE_FORMAT = (TYPE = CSV SKIP_HEADER = 1)")
		if stmt.Action != ast.AlterStageSet {
			t.Fatalf("Action = %v, want Set", stmt.Action)
		}
		if findOption(stmt.Options, "FILE_FORMAT") == nil {
			t.Errorf("missing FILE_FORMAT; options=%+v", stmt.Options)
		}
	})

	t.Run("set comment", func(t *testing.T) {
		stmt := mustAlterStage(t, "ALTER STAGE s SET COMMENT = 'updated'")
		c := findOption(stmt.Options, "COMMENT")
		if c == nil || c.Lit == nil || c.Lit.Value != "updated" {
			t.Errorf("COMMENT = %+v", c)
		}
	})

	t.Run("set external params", func(t *testing.T) {
		stmt := mustAlterStage(t, "ALTER STAGE s SET URL = 's3://b/p2' CREDENTIALS = (AWS_KEY_ID = 'k' AWS_SECRET_KEY = 's') ENCRYPTION = (TYPE = 'AWS_SSE_S3')")
		for _, k := range []string{"URL", "CREDENTIALS", "ENCRYPTION"} {
			if findOption(stmt.Options, k) == nil {
				t.Errorf("missing %s; options=%+v", k, stmt.Options)
			}
		}
	})

	t.Run("set copy_options", func(t *testing.T) {
		stmt := mustAlterStage(t, "ALTER STAGE s SET COPY_OPTIONS = (ON_ERROR = 'SKIP_FILE')")
		if findOption(stmt.Options, "COPY_OPTIONS") == nil {
			t.Errorf("missing COPY_OPTIONS; options=%+v", stmt.Options)
		}
	})

	t.Run("set directory", func(t *testing.T) {
		stmt := mustAlterStage(t, "ALTER STAGE s SET DIRECTORY = (ENABLE = TRUE)")
		if groupEntry(findOption(stmt.Options, "DIRECTORY"), "ENABLE") == nil {
			t.Errorf("missing DIRECTORY ENABLE; options=%+v", stmt.Options)
		}
	})
}

func TestParseAlterStage_Tags(t *testing.T) {
	t.Run("set tag single", func(t *testing.T) {
		stmt := mustAlterStage(t, "ALTER STAGE s SET TAG cost = '42'")
		if stmt.Action != ast.AlterStageSetTag {
			t.Fatalf("Action = %v, want SetTag", stmt.Action)
		}
		if len(stmt.Tags) != 1 || stmt.Tags[0].Name.String() != "cost" || stmt.Tags[0].Value != "42" {
			t.Errorf("Tags = %+v", stmt.Tags)
		}
	})

	t.Run("set tag multiple", func(t *testing.T) {
		stmt := mustAlterStage(t, "ALTER STAGE s SET TAG a = '1', b = '2'")
		if len(stmt.Tags) != 2 {
			t.Fatalf("Tags = %+v, want 2", stmt.Tags)
		}
	})

	t.Run("set tag qualified name", func(t *testing.T) {
		stmt := mustAlterStage(t, "ALTER STAGE s SET TAG db.sch.t = 'v'")
		if stmt.Tags[0].Name.String() != "db.sch.t" {
			t.Errorf("tag name = %q", stmt.Tags[0].Name.String())
		}
	})

	t.Run("unset tag single", func(t *testing.T) {
		stmt := mustAlterStage(t, "ALTER STAGE s UNSET TAG cost")
		if stmt.Action != ast.AlterStageUnsetTag {
			t.Fatalf("Action = %v, want UnsetTag", stmt.Action)
		}
		if len(stmt.UnsetTags) != 1 || stmt.UnsetTags[0].String() != "cost" {
			t.Errorf("UnsetTags = %+v", stmt.UnsetTags)
		}
	})

	t.Run("unset tag multiple", func(t *testing.T) {
		stmt := mustAlterStage(t, "ALTER STAGE s UNSET TAG a, b, c")
		if len(stmt.UnsetTags) != 3 {
			t.Errorf("UnsetTags = %+v, want 3", stmt.UnsetTags)
		}
	})
}

func TestParseAlterStage_Unset(t *testing.T) {
	t.Run("unset comment", func(t *testing.T) {
		stmt := mustAlterStage(t, "ALTER STAGE s UNSET COMMENT")
		if stmt.Action != ast.AlterStageUnset {
			t.Fatalf("Action = %v, want Unset", stmt.Action)
		}
		if len(stmt.UnsetProps) != 1 || stmt.UnsetProps[0] != "COMMENT" {
			t.Errorf("UnsetProps = %+v, want [COMMENT]", stmt.UnsetProps)
		}
	})

	t.Run("unset dcm project (multi-word, docs)", func(t *testing.T) {
		stmt := mustAlterStage(t, "ALTER STAGE s UNSET DCM PROJECT")
		if len(stmt.UnsetProps) != 1 || stmt.UnsetProps[0] != "DCM PROJECT" {
			t.Errorf("UnsetProps = %+v, want [\"DCM PROJECT\"]", stmt.UnsetProps)
		}
	})

	t.Run("unset multiple properties", func(t *testing.T) {
		stmt := mustAlterStage(t, "ALTER STAGE s UNSET URL, CREDENTIALS")
		if len(stmt.UnsetProps) != 2 {
			t.Errorf("UnsetProps = %+v, want 2", stmt.UnsetProps)
		}
	})
}

func TestParseAlterStage_Refresh(t *testing.T) {
	t.Run("refresh no subpath", func(t *testing.T) {
		stmt := mustAlterStage(t, "ALTER STAGE s REFRESH")
		if stmt.Action != ast.AlterStageRefresh {
			t.Fatalf("Action = %v, want Refresh", stmt.Action)
		}
		if stmt.Subpath != nil {
			t.Errorf("Subpath = %v, want nil", *stmt.Subpath)
		}
	})

	t.Run("refresh subpath", func(t *testing.T) {
		stmt := mustAlterStage(t, "ALTER STAGE s REFRESH SUBPATH = '/sub/dir'")
		if stmt.Subpath == nil || *stmt.Subpath != "/sub/dir" {
			t.Errorf("Subpath = %v, want /sub/dir", stmt.Subpath)
		}
	})

	t.Run("if exists refresh subpath", func(t *testing.T) {
		stmt := mustAlterStage(t, "ALTER STAGE IF EXISTS s REFRESH SUBPATH = 'p'")
		if !stmt.IfExists || stmt.Subpath == nil {
			t.Errorf("IfExists=%v Subpath=%v", stmt.IfExists, stmt.Subpath)
		}
	})
}

// ---------------------------------------------------------------------------
// Negative tests — malformed CREATE / ALTER STAGE must be rejected
// ---------------------------------------------------------------------------

func TestParseStage_Negative(t *testing.T) {
	bad := []string{
		"CREATE STAGE",                             // missing name
		"CREATE STAGE s URL =",                     // option missing value
		"CREATE STAGE s FILE_FORMAT = (TYPE = CSV", // unterminated group
		"CREATE STAGE s ENCRYPTION = (TYPE = )",    // group entry missing value
		"CREATE STAGE IF NOT s",                    // malformed IF NOT EXISTS
		"ALTER STAGE",                              // missing name
		"ALTER STAGE s",                            // missing action
		"ALTER STAGE s RENAME",                     // missing TO
		"ALTER STAGE s RENAME TO",                  // missing new name
		"ALTER STAGE s SET",                        // nothing to set
		"ALTER STAGE s SET TAG t =",                // tag missing value
		"ALTER STAGE s SET TAG",                    // missing tag assignment
		"ALTER STAGE s UNSET",                      // missing property
		"ALTER STAGE s UNSET TAG",                  // missing tag name
		"ALTER STAGE s REFRESH SUBPATH =",          // subpath missing value
		"ALTER STAGE s REFRESH SUBPATH = bad",      // subpath not a string
		"ALTER STAGE s FROBNICATE",                 // unknown action
	}
	for _, in := range bad {
		t.Run(in, func(t *testing.T) {
			result := ParseBestEffort(in)
			if len(result.Errors) == 0 {
				t.Errorf("expected parse error for %q, got none (stmts=%d)", in, len(result.File.Stmts))
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Loc accuracy
// ---------------------------------------------------------------------------

func TestParseStage_Loc(t *testing.T) {
	input := "CREATE STAGE my_stage"
	stmt := mustCreateStage(t, input)
	if stmt.Loc.Start != 0 || stmt.Loc.End != len(input) {
		t.Errorf("CREATE Loc = %+v, want {0, %d}", stmt.Loc, len(input))
	}

	// ALTER sub-parsers set Loc.Start at the object-type keyword (STAGE), not the
	// ALTER keyword, matching the established ALTER TABLE/DATABASE/SCHEMA
	// convention (the dispatcher consumes ALTER before the sub-parser runs).
	ainput := "ALTER STAGE s RENAME TO s2"
	const stageKwOff = len("ALTER ")
	astmt := mustAlterStage(t, ainput)
	if astmt.Loc.Start != stageKwOff || astmt.Loc.End != len(ainput) {
		t.Errorf("ALTER Loc = %+v, want {%d, %d}", astmt.Loc, stageKwOff, len(ainput))
	}

	// Second statement in a multi-statement input gets a non-zero base offset;
	// the option literal must still be correct (exercises base handling).
	multi := "SELECT 1;\nCREATE STAGE s URL = 's3://b/p'"
	res := ParseBestEffort(multi)
	if len(res.Errors) != 0 {
		t.Fatalf("multi parse errors: %v", res.Errors)
	}
	cs, ok := res.File.Stmts[1].(*ast.CreateStageStmt)
	if !ok {
		t.Fatalf("stmt[1] = %T", res.File.Stmts[1])
	}
	if findOption(cs.Options, "URL").Lit.Value != "s3://b/p" {
		t.Errorf("URL value wrong in multi-statement input")
	}
}
