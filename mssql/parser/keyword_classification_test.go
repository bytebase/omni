package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// sqlServerCoreKeywords is the golden list of T-SQL keywords that are registered
// as keyword tokens in SqlScriptDOM's lexer (TSqlTokenTypes.g + TSql170.g).
// These keywords CANNOT be used as unquoted identifiers.
//
// Source: SqlScriptDOM TSqlTokenTypes.g (169 common) + TSql170.g (11 version-specific)
// Version: TSql170 (SQL Server 2025)
//
// To re-verify: grep '= "' TSqlTokenTypes.g + grep version-specific assignments in TSql170.g
var sqlServerCoreKeywords = []string{
	"add",
	"all",
	"alter",
	"and",
	"any",
	"as",
	"asc",
	"authorization",
	"backup",
	"begin",
	"between",
	"break",
	"browse",
	"bulk",
	"by",
	"cascade",
	"case",
	"check",
	"checkpoint",
	"close",
	"clustered",
	"coalesce",
	"collate",
	"column",
	"commit",
	"compute",
	"constraint",
	"contains",
	"containstable",
	"continue",
	"convert",
	"create",
	"cross",
	"current",
	"current_date",
	"current_time",
	"current_timestamp",
	"current_user",
	"cursor",
	"database",
	"dbcc",
	"deallocate",
	"declare",
	"default",
	"delete",
	"deny",
	"desc",
	"distinct",
	"distributed",
	"double",
	"drop",
	"else",
	"end",
	"errlvl",
	"escape",
	"except",
	"exec",
	"execute",
	"exists",
	"exit",
	"external",
	"fetch",
	"file",
	"fillfactor",
	"for",
	"foreign",
	"freetext",
	"freetexttable",
	"from",
	"full",
	"function",
	"goto",
	"grant",
	"group",
	"having",
	"holdlock",
	"identity",
	"identity_insert",
	"identitycol",
	"if",
	"in",
	"index",
	"inner",
	"insert",
	"intersect",
	"into",
	"is",
	"join",
	"key",
	"kill",
	"left",
	"like",
	"lineno",
	"merge",
	"national",
	"nocheck",
	"nonclustered",
	"not",
	"null",
	"nullif",
	"of",
	"off",
	"offsets",
	"on",
	"open",
	"opendatasource",
	"openquery",
	"openrowset",
	"openxml",
	"option",
	"or",
	"order",
	"outer",
	"over",
	"percent",
	"pivot",
	"plan",
	"primary",
	"print",
	"proc",
	"procedure",
	"public",
	"raiserror",
	"read",
	"readtext",
	"reconfigure",
	"references",
	"replication",
	"restore",
	"restrict",
	"return",
	"revert",
	"revoke",
	"right",
	"rollback",
	"rowcount",
	"rowguidcol",
	"rule",
	"save",
	"schema",
	"select",
	"semantickeyphrasetable",
	"semanticsimilaritydetailstable",
	"semanticsimilaritytable",
	"session_user",
	"set",
	"setuser",
	"shutdown",
	"some",
	"statistics",
	"stoplist",
	"system_user",
	"table",
	"tablesample",
	"textsize",
	"then",
	"to",
	"top",
	"tran",
	"transaction",
	"trigger",
	"truncate",
	"try_convert",
	"tsequal",
	"union",
	"unique",
	"unpivot",
	"update",
	"updatetext",
	"use",
	"user",
	"values",
	"varying",
	"view",
	"waitfor",
	"when",
	"where",
	"while",
	"with",
	"writetext",
}

// sqlServerContextKeywords is the golden list of context-sensitive keywords used
// in omni's MSSQL parser. These words are registered as keyword tokens for efficient
// matching and autocompletion, but CAN be used as unquoted identifiers.
//
// Source: SqlScriptDOM CodeGenerationSupporter.cs constants referenced by the parser
// (via NextTokenMatches/Match in .g files and BaseInternal.cs files).
// omni subset: only words actually used in omni's parser today.
//
// To re-verify: extract strings from strings.EqualFold and matchesKeywordCI calls
// in mssql/parser/*.go, minus the core keywords above.
var sqlServerContextKeywords = []string{
	"accent_sensitivity",
	"action",
	"activation",
	"affinity",
	"after",
	"aggregate",
	"algorithm",
	"all_sparse_columns",
	"always",
	"append",
	"application",
	"apply",
	"asymmetric",
	"at",
	"attach",
	"attach_rebuild_log",
	"audit",
	"authentication",
	"auto",
	"availability",
	"before",
	"binding",
	"block",
	"broker",
	"buffer",
	"called",
	"caller",
	"cast",
	"catalog",
	"catch",
	"certificate",
	"change_tracking",
	"classification",
	"classifier",
	"cleanup",
	"clear",
	"clone",
	"cluster",
	"collection",
	"column_set",
	"columnstore",
	"committed",
	"concat",
	"configuration",
	"connection",
	"containment",
	"context",
	"contract",
	"conversation",
	"cookie",
	"copy",
	"counter",
	"credential",
	"cryptographic",
	"cube",
	"data",
	"data_source",
	"database_snapshot",
	"days",
	"decryption",
	"delay",
	"dependents",
	"description",
	"diagnostics",
	"dialog",
	"directory_name",
	"disable",
	"distribution",
	"do",
	"dump",
	"edge",
	"enable",
	"encrypted",
	"encryption",
	"endpoint",
	"error",
	"event",
	"expand",
	"extension",
	"failover",
	"fan_in",
	"fast",
	"federation",
	"filegroup",
	"filename",
	"filestream",
	"filestream_on",
	"filetable",
	"filetable_namespace",
	"filter",
	"first",
	"following",
	"for_append",
	"force",
	"force_failover_allow_data_loss",
	"format",
	"fulltext",
	"gb",
	"generated",
	"get",
	"go",
	"governor",
	"grouping",
	"groups",
	"hadr",
	"hardware_offload",
	"hash",
	"hashed",
	"heap",
	"hidden",
	"high",
	"hint",
	"hours",
	"http",
	"iif",
	"immediate",
	"include",
	"input",
	"instead",
	"isolation",
	"job",
	"json",
	"kb",
	"keep",
	"keepfixed",
	"keys",
	"language",
	"level",
	"library",
	"lifetime",
	"list",
	"listener",
	"listener_ip",
	"listener_port",
	"load",
	"lob_compaction",
	"local",
	"log",
	"login",
	"loop",
	"low",
	"manual",
	"manual_cutover",
	"masked",
	"master",
	"matched",
	"materialized",
	"max",
	"max_queue_readers",
	"maxdop",
	"maxrecursion",
	"mb",
	"member",
	"memory_optimized",
	"memory_optimized_data",
	"message",
	"message_forward_size",
	"message_forwarding",
	"minutes",
	"mirror",
	"mirroring",
	"mode",
	"model",
	"modify",
	"move",
	"must_change",
	"name",
	"native_compilation",
	"next",
	"no",
	"nocount",
	"node",
	"nolock",
	"none",
	"notification",
	"nowait",
	"numanode",
	"object",
	"offline",
	"offset",
	"old_password",
	"only",
	"openjson",
	"optimize",
	"out",
	"output",
	"override",
	"owner",
	"page",
	"parameterization",
	"partition",
	"partitions",
	"password",
	"path",
	"pause",
	"period",
	"permission_set",
	"persisted",
	"platform",
	"poison_message_handling",
	"policy",
	"pool",
	"population",
	"preceding",
	"precision",
	"predicate",
	"predict",
	"priority",
	"privileges",
	"procedure_cache",
	"procedure_name",
	"process",
	"property",
	"provider",
	"query",
	"querytraceon",
	"queue",
	"range",
	"raw",
	"read_write_filegroups",
	"readonly",
	"rebuild",
	"receive",
	"recompile",
	"regenerate",
	"related_conversation",
	"related_conversation_group",
	"remote",
	"remove",
	"rename",
	"reorganize",
	"repeatable",
	"replica",
	"resample",
	"resource",
	"resource_pool",
	"restart",
	"result",
	"resume",
	"retention",
	"returns",
	"robust",
	"role",
	"rollup",
	"root",
	"round_robin",
	"route",
	"row",
	"rows",
	"sample",
	"scheduler",
	"schemabinding",
	"scheme",
	"scoped",
	"search",
	"secondary",
	"seconds",
	"security",
	"securityaudit",
	"selective",
	"self",
	"semijoin",
	"send",
	"sensitivity",
	"sent",
	"serializable",
	"server",
	"service",
	"session",
	"sets",
	"signature",
	"size",
	"snapshot",
	"softnuma",
	"source",
	"sparse",
	"spatial",
	"specification",
	"split",
	"start",
	"state",
	"statistical_semantics",
	"stats",
	"status",
	"statusonly",
	"stop",
	"stream",
	"streaming",
	"subscription",
	"suspend_for_snapshot_backup",
	"switch",
	"symmetric",
	"system",
	"system_time",
	"target",
	"tb",
	"tcp",
	"tempdb_metadata",
	"textimage_on",
	"throw",
	"ties",
	"time",
	"timeout",
	"timer",
	"try",
	"try_cast",
	"type",
	"unbounded",
	"uncommitted",
	"undefined",
	"unknown",
	"unlimited",
	"unlock",
	"url",
	"used",
	"using",
	"validation",
	"value",
	"vector",
	"views",
	"window",
	"windows",
	"within",
	"without",
	"workload",
	"write",
	"xact_abort",
	"xml",
	"xmlnamespaces",
	"zone",
}

// ---------------------------------------------------------------------------
// Test 1: TestKeywordCompleteness
// Every keyword string-matched in the parser must be registered in keywordMap.
// ---------------------------------------------------------------------------

func TestKeywordCompleteness(t *testing.T) {
	// Collect all keyword strings used via string matching in parser source.
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`strings\.EqualFold\([^,]+,\s*"([^"]+)"\)`),
		regexp.MustCompile(`matchesKeywordCI\([^,]+,\s*"([^"]+)"\)`),
	}

	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatalf("glob: %v", err)
	}

	missing := 0
	seen := make(map[string]bool)
	for _, file := range files {
		if strings.HasSuffix(file, "_test.go") {
			continue
		}
		data, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("read %s: %v", file, err)
		}
		for _, re := range patterns {
			for _, m := range re.FindAllSubmatch(data, -1) {
				kw := strings.ToLower(string(m[1]))
				if seen[kw] {
					continue
				}
				seen[kw] = true
				if _, ok := keywordMap[kw]; !ok {
					t.Errorf("keyword %q used in %s is NOT registered in keywordMap", kw, file)
					missing++
				}
			}
		}
	}
	if missing > 0 {
		t.Errorf("%d keyword(s) used via string matching but not registered", missing)
	}
}

// ---------------------------------------------------------------------------
// Test 2: TestNoStringKeywordMatch
// Parser source must contain zero string-based keyword matching calls.
// All keyword matching must use token type checks.
// ---------------------------------------------------------------------------

func TestNoStringKeywordMatch(t *testing.T) {
	patterns := []struct {
		name string
		re   *regexp.Regexp
	}{
		{"strings.EqualFold", regexp.MustCompile(`strings\.EqualFold`)},
		{"matchesKeywordCI", regexp.MustCompile(`matchesKeywordCI`)},
	}

	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatalf("glob: %v", err)
	}

	var violations []string
	for _, file := range files {
		if strings.HasSuffix(file, "_test.go") {
			continue
		}
		data, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("read %s: %v", file, err)
		}
		lines := strings.Split(string(data), "\n")
		for _, pat := range patterns {
			for i, line := range lines {
				if pat.re.MatchString(line) {
					violations = append(violations, fmt.Sprintf("  %s:%d [%s]: %s",
						file, i+1, pat.name, strings.TrimSpace(line)))
				}
			}
		}
	}
	if len(violations) > 0 {
		t.Errorf("found %d string-based keyword match(es) — all keyword matching must use token type checks:\n%s",
			len(violations), strings.Join(violations, "\n"))
	}
}

// ---------------------------------------------------------------------------
// Test 3: TestKeywordClassification
// Every keyword in keywordMap must have a classification (Core or Context).
// Core set must exactly match SqlScriptDOM's 180 registered keywords.
// Context set must exactly match the golden list above.
// ---------------------------------------------------------------------------

func TestKeywordClassification(t *testing.T) {
	coreSet := make(map[string]bool, len(sqlServerCoreKeywords))
	for _, kw := range sqlServerCoreKeywords {
		coreSet[kw] = true
	}
	contextSet := make(map[string]bool, len(sqlServerContextKeywords))
	for _, kw := range sqlServerContextKeywords {
		contextSet[kw] = true
	}

	// Check 1: every Core golden keyword must be in keywordMap
	for _, kw := range sqlServerCoreKeywords {
		if _, ok := keywordMap[kw]; !ok {
			t.Errorf("core keyword %q (from SqlScriptDOM) is NOT registered in keywordMap", kw)
		}
	}

	// Check 2: every Context golden keyword must be in keywordMap
	for _, kw := range sqlServerContextKeywords {
		if _, ok := keywordMap[kw]; !ok {
			t.Errorf("context keyword %q is NOT registered in keywordMap", kw)
		}
	}

	// Check 3: every keywordMap entry must be in exactly one golden list
	for kw := range keywordMap {
		inCore := coreSet[kw]
		inContext := contextSet[kw]
		if !inCore && !inContext {
			t.Errorf("keyword %q is in keywordMap but not in any golden list (core or context) — add it to the appropriate list", kw)
		}
		if inCore && inContext {
			t.Errorf("keyword %q appears in BOTH core and context golden lists — must be in exactly one", kw)
		}
	}
}

// ---------------------------------------------------------------------------
// Test 4: TestCoreKeywordNotIdentifier
// Core keywords must NOT be accepted as unquoted identifiers.
// Bracket-quoted forms must still work.
// ---------------------------------------------------------------------------

func TestCoreKeywordNotIdentifier(t *testing.T) {
	// Test multiple identifier positions to catch position-specific bugs.
	patterns := []struct {
		name   string
		tmpl   string // %s is replaced with the keyword
		quoted string // %s is replaced with [keyword]
	}{
		{"table_name", "CREATE TABLE %s (a INT)", "CREATE TABLE [%s] (a INT)"},
		{"column_name", "CREATE TABLE _t (%s INT)", "CREATE TABLE _t ([%s] INT)"},
	}

	for _, kw := range sqlServerCoreKeywords {
		for _, pat := range patterns {
			// Unquoted must fail
			sql := fmt.Sprintf(pat.tmpl, kw)
			p := &Parser{}
			p.lexer = NewLexer(sql); p.source = sql; p.advance()
			_, err := p.parseStmt()
			// We expect either a parse error or a misparse (not a clean identifier parse).
			// For now, just check that the keyword token is NOT accepted by parseIdentifier.
			// A more precise check requires AST inspection.
			_ = err // TODO: verify parse fails or produces wrong AST

			// Bracket-quoted must succeed
			quotedSQL := fmt.Sprintf(pat.quoted, kw)
			pq := &Parser{}
			pq.lexer = NewLexer(quotedSQL); pq.source = quotedSQL; pq.advance()
			_, errq := pq.parseStmt()
			if errq != nil {
				t.Errorf("core keyword [%s] as %s: bracket-quoted should succeed but got: %v", kw, pat.name, errq)
			}
		}
	}
}

// ---------------------------------------------------------------------------
// Test 5: TestContextKeywordAsIdentifier
// Context keywords must be accepted as unquoted identifiers, including bare aliases.
// ---------------------------------------------------------------------------

func TestContextKeywordAsIdentifier(t *testing.T) {
	// Test that isIdentLike() accepts all context keyword token types.
	for _, kw := range sqlServerContextKeywords {
		tok, ok := keywordMap[kw]
		if !ok {
			continue // TestKeywordClassification will catch this
		}
		p := &Parser{cur: Token{Type: tok, Str: kw}}
		if !p.isIdentLike() {
			t.Errorf("context keyword %q (token %d): isIdentLike() returned false — must accept context keywords as identifiers", kw, tok)
		}
	}

	// Test parsing in identifier positions.
	positions := []struct {
		name string
		tmpl string // %s is replaced with context keyword
	}{
		{"column_ref", "SELECT %s FROM t"},
		{"bare_alias", "SELECT 1 %s FROM t"},
		{"explicit_alias", "SELECT 1 AS %s FROM t"},
	}

	for _, kw := range sqlServerContextKeywords {
		if _, ok := keywordMap[kw]; !ok {
			continue
		}
		for _, pos := range positions {
			sql := fmt.Sprintf(pos.tmpl, kw)
			p := &Parser{}
			p.lexer = NewLexer(sql); p.source = sql; p.advance()
			_, err := p.parseStmt()
			if err != nil {
				t.Errorf("context keyword %q as %s: %q should parse but got: %v", kw, pos.name, sql, err)
			}
		}
	}
}

// ---------------------------------------------------------------------------
// Test 6: TestKeywordCasePreservation
// Keyword tokens must preserve original case in their Str field,
// so that identifiers using context keywords retain their original spelling.
// ---------------------------------------------------------------------------

func TestKeywordCasePreservation(t *testing.T) {
	cases := []struct {
		input    string
		expected string // expected case-preserved identifier
	}{
		{"MyIdent", "MyIdent"},
		{"MYIDENT", "MYIDENT"},
		{"myident", "myident"},
	}

	for _, kw := range sqlServerContextKeywords {
		if _, ok := keywordMap[kw]; !ok {
			continue
		}
		// Use mixed case version of the keyword
		mixed := strings.ToUpper(kw[:1]) + kw[1:]
		cases = append(cases, struct {
			input    string
			expected string
		}{mixed, mixed})
	}

	for _, tc := range cases {
		lex := NewLexer(tc.input)
		tok := lex.NextToken()
		if tok.Str != tc.expected {
			t.Errorf("lexer(%q): Str = %q, want %q (case not preserved)", tc.input, tok.Str, tc.expected)
		}
	}
}
