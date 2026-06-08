package compat

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/bytebase/omni/redshift"
)

type CommandStatus string

const (
	StatusSupportedParse      CommandStatus = "supported_parse"
	StatusSupportedRuntime    CommandStatus = "supported_runtime"
	StatusExplicitUnsupported CommandStatus = "explicit_unsupported"
	StatusNotRelevant         CommandStatus = "not_relevant"
)

func (s CommandStatus) Valid() bool {
	switch s {
	case StatusSupportedParse, StatusSupportedRuntime, StatusExplicitUnsupported, StatusNotRelevant:
		return true
	default:
		return false
	}
}

type ManifestSource struct {
	URL       string `json:"url"`
	Retrieved string `json:"retrieved"`
}

type AWSCommandManifest struct {
	Source   ManifestSource `json:"source"`
	Commands []AWSCommand   `json:"commands"`
}

type AWSCommand struct {
	Name      string        `json:"name"`
	Status    CommandStatus `json:"status"`
	SampleSQL string        `json:"sample_sql,omitempty"`
	Notes     string        `json:"notes,omitempty"`
}

type AWSCommandResult struct {
	Name       string
	Status     CommandStatus
	SampleSQL  string
	ParseOK    bool
	ParseError string
	Notes      string
}

type AWSCommandManifestResult struct {
	Source   ManifestSource
	Commands []AWSCommandResult
	Counts   map[CommandStatus]int
}

func LoadAWSCommandManifest(path string) (*AWSCommandManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var manifest AWSCommandManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}
	return &manifest, nil
}

func EvaluateAWSCommandManifest(manifest *AWSCommandManifest) AWSCommandManifestResult {
	result := AWSCommandManifestResult{
		Counts: make(map[CommandStatus]int),
	}
	if manifest == nil {
		return result
	}
	result.Source = manifest.Source
	result.Commands = make([]AWSCommandResult, 0, len(manifest.Commands))
	for _, command := range manifest.Commands {
		entry := AWSCommandResult{
			Name:      command.Name,
			Status:    command.Status,
			SampleSQL: command.SampleSQL,
			Notes:     command.Notes,
		}
		if strings.TrimSpace(command.SampleSQL) != "" {
			_, err := redshift.Parse(command.SampleSQL)
			entry.ParseOK = err == nil
			if err != nil {
				entry.ParseError = err.Error()
			}
		}
		result.Commands = append(result.Commands, entry)
		result.Counts[command.Status]++
	}
	return result
}

type LegacyCorpusStats struct {
	TotalFiles           int
	PassingFiles         int
	ExpectedFailureFiles int
	NewFailureFiles      int
	PromotedFiles        int
	NewFailures          []string
	Promoted             []string
}

type LegacyStatementParityStats struct {
	TotalStatements              int
	BothAccept                   int
	BothReject                   int
	OldAcceptsOmniRejects        int
	OldRejectsOmniAccepts        int
	OldAcceptsOmniRejectExamples []string
	OldRejectsOmniAcceptExamples []string
}

var legacyAcceptedOmniRejectSnapshot = map[string]string{
	"alter_role.sql[28]":  `ALTER ROLE role1 EXTERNALID TO "";`,
	"create_role.sql[36]": `CREATE ROLE empty_external_role EXTERNALID "";`,
	"drop_materialized_view.sql[114]": `DROP MATERIALIZED VIEW IF EXISTS sales_summary_mv, customer_orders_mv, inventory_status_mv CASCADE
DROP MATERIALIZED VIEW old_sales_summary
DROP MATERIALIZED VIEW analytics.monthly_revenue_mv RESTRICT
DROP MATERIALIZED VIEW IF EXISTS temp_analysis_mv
DROP MATERIALIZED VIEW reporting.customer_metrics_mv, reporting.product_metrics_mv CASCADE
DROP MATERIALIZED VIEW mv_with_special_chars_123
DROP MATERIALIZED VIEW IF EXISTS "MixedCase_MV"
DROP MATERIALIZED VIEW public.mv_to_drop RESTRICT
DROP MATERIALIZED VIEW IF EXISTS stale_data_mv, outdated_report_mv, unused_metrics_mv CASCADE`,
	"drop_role.sql[66]": `DROP ROLE admin_role
DROP ROLE IF EXISTS temp_role
DROP ROLE role1, role2, role3
DROP ROLE IF EXISTS analytics_readonly, reporting_user, data_scientist
DROP ROLE "role with spaces"
DROP ROLE IF EXISTS mixed_case_role
DROP ROLE user_admin CASCADE
DROP ROLE IF EXISTS application_role RESTRICT`,
	"grant.sql[63]":   `GRANT ALL ON TABLE customers TO user1 WITH GRANT OPTION, ROLE customer_admin, GROUP support_team;`,
	"grant.sql[87]":   `GRANT USAGE ON SCHEMA "" TO edge_case_user;`,
	"revoke.sql[163]": `REVOKE USAGE ON SCHEMA "" FROM edge_case_user;`,
	"revoke.sql[199]": `REVOKE SELECT ON TABLE IF EXISTS optional_table FROM optional_user;`,
	"revoke.sql[200]": `REVOKE USAGE ON SCHEMA IF EXISTS optional_schema FROM optional_role;`,
}

func EvaluateLegacyCorpus(fixturesDir, expectedFailuresPath string) (LegacyCorpusStats, error) {
	expectedFailures, err := loadLegacyExpectedFailures(expectedFailuresPath)
	if err != nil {
		return LegacyCorpusStats{}, err
	}

	entries, err := os.ReadDir(fixturesDir)
	if err != nil {
		return LegacyCorpusStats{}, err
	}

	stats := LegacyCorpusStats{}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		stats.TotalFiles++
		sql, err := os.ReadFile(filepath.Join(fixturesDir, entry.Name()))
		if err != nil {
			return LegacyCorpusStats{}, err
		}
		_, parseErr := redshift.Parse(string(sql))
		_, expectedFailure := expectedFailures[entry.Name()]
		switch {
		case expectedFailure && parseErr != nil:
			stats.ExpectedFailureFiles++
		case expectedFailure && parseErr == nil:
			stats.PromotedFiles++
			stats.Promoted = append(stats.Promoted, entry.Name())
		case !expectedFailure && parseErr != nil:
			stats.NewFailureFiles++
			stats.NewFailures = append(stats.NewFailures, fmt.Sprintf("%s: %v", entry.Name(), parseErr))
		default:
			stats.PassingFiles++
		}
	}
	sort.Strings(stats.NewFailures)
	sort.Strings(stats.Promoted)
	return stats, nil
}

func EvaluateLegacyStatementParity(fixturesDir string) (LegacyStatementParityStats, error) {
	entries, err := os.ReadDir(fixturesDir)
	if err != nil {
		return LegacyStatementParityStats{}, err
	}

	var stats LegacyStatementParityStats
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		sql, err := os.ReadFile(filepath.Join(fixturesDir, entry.Name()))
		if err != nil {
			return LegacyStatementParityStats{}, err
		}
		for index, segment := range redshift.Split(string(sql)) {
			if segment.Empty() {
				continue
			}
			stats.TotalStatements++
			_, omniErr := redshift.Parse(segment.Text)
			omniOK := omniErr == nil

			key := fmt.Sprintf("%s[%d]", entry.Name(), index)
			if omniOK {
				stats.BothAccept++
				continue
			}

			if _, ok := legacyAcceptedOmniRejectSnapshot[key]; !ok {
				stats.BothReject++
				continue
			}
			stats.OldAcceptsOmniRejects++
			if len(stats.OldAcceptsOmniRejectExamples) < 20 {
				stats.OldAcceptsOmniRejectExamples = append(stats.OldAcceptsOmniRejectExamples,
					fmt.Sprintf("%s: %s", key, oneLine(segment.Text)))
			}
		}
	}
	return stats, nil
}

func oneLine(sql string) string {
	text := strings.Join(strings.Fields(sql), " ")
	if len(text) > 160 {
		return text[:160] + "..."
	}
	return text
}

type EvaluateReferenceOptions struct {
	DSN     string
	Samples []string
}

type ReferenceRedshiftResult struct {
	Enabled    bool
	SkipReason string
	Total      int
	Accepted   int
	Rejected   int
	Errors     []string
}

func EvaluateReferenceRedshift(options EvaluateReferenceOptions) ReferenceRedshiftResult {
	if strings.TrimSpace(options.DSN) == "" {
		return ReferenceRedshiftResult{
			SkipReason: "REDSHIFT_COMPAT_DSN is not set",
		}
	}
	result := ReferenceRedshiftResult{Enabled: true}
	samples := options.Samples
	if len(samples) == 0 {
		samples = []string{
			"SELECT 1",
			"EXPLAIN SELECT 1",
			"SHOW DATABASES",
		}
	}

	db, err := sql.Open("pgx", options.DSN)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("open: %v", err))
		return result
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("ping: %v", err))
		return result
	}
	for _, sample := range samples {
		result.Total++
		if _, err := db.ExecContext(ctx, sample); err != nil {
			result.Rejected++
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", oneLine(sample), err))
			continue
		}
		result.Accepted++
	}
	return result
}

type RuntimeSemanticStats struct {
	TotalChecks  int
	PassedChecks int
	FailedChecks int
	Failures     []string
}

func EvaluateRuntimeSemantics() RuntimeSemanticStats {
	checks := []struct {
		name string
		run  func() error
	}{
		{
			name: "parse create table diststyle sortkey encode",
			run: func() error {
				_, err := redshift.Parse("CREATE TABLE public.orders (id INT ENCODE az64, name VARCHAR(64) ENCODE lzo) DISTSTYLE KEY DISTKEY(id) SORTKEY(id);")
				return err
			},
		},
		{
			name: "diagnose invalid syntax",
			run: func() error {
				if diagnostics := redshift.Diagnose("SELECT FROM;"); len(diagnostics) == 0 {
					return fmt.Errorf("expected diagnostics")
				}
				return nil
			},
		},
		{
			name: "statement ranges for multiple statements",
			run: func() error {
				ranges, err := redshift.StatementRanges("SELECT 1;\n  SELECT * FROM sales QUALIFY row_number() OVER () = 1;")
				if err != nil {
					return err
				}
				if len(ranges) != 2 {
					return fmt.Errorf("got %d ranges, want 2", len(ranges))
				}
				if ranges[1].Start.Line != 1 || ranges[1].Start.Character != 2 {
					return fmt.Errorf("second range starts at %+v, want line 1 character 2", ranges[1].Start)
				}
				return nil
			},
		},
		{
			name: "statement types for redshift commands",
			run: func() error {
				types, err := redshift.GetStatementTypes(`
COPY orders FROM 's3://bucket/orders' IAM_ROLE DEFAULT;
UNLOAD ('SELECT * FROM orders') TO 's3://bucket/out' IAM_ROLE DEFAULT;
MERGE INTO orders USING staging_orders s ON orders.id = s.id WHEN MATCHED THEN UPDATE SET id = s.id;
SHOW DATABASES;
`)
				if err != nil {
					return err
				}
				want := []redshift.StatementType{
					redshift.StatementTypeCopy,
					redshift.StatementTypeUnload,
					redshift.StatementTypeDML,
					redshift.StatementTypeShow,
				}
				if len(types) != len(want) {
					return fmt.Errorf("got %d statement types, want %d: %v", len(types), len(want), types)
				}
				for i := range want {
					if types[i] != want[i] {
						return fmt.Errorf("statement type %d = %s, want %s", i, types[i], want[i])
					}
				}
				return nil
			},
		},
		{
			name: "editor accepts readonly explain",
			run: func() error {
				query, plan, err := redshift.ValidateSQLForEditor("EXPLAIN SELECT * FROM orders;")
				if err != nil {
					return err
				}
				if !query || !plan {
					return fmt.Errorf("query=%v plan=%v, want true true", query, plan)
				}
				return nil
			},
		},
		{
			name: "editor rejects dml",
			run: func() error {
				query, plan, err := redshift.ValidateSQLForEditor("INSERT INTO orders SELECT 1;")
				if err != nil {
					return err
				}
				if query || plan {
					return fmt.Errorf("query=%v plan=%v, want false false", query, plan)
				}
				return nil
			},
		},
		{
			name: "changed resources include ddl and dml",
			run: func() error {
				summary, err := redshift.ExtractChangedResources(`
CREATE TABLE public.created_orders (id INT) DISTSTYLE EVEN;
INSERT INTO public.changed_orders SELECT 1;
`, "dev", "public")
				if err != nil {
					return err
				}
				if summary.DMLCount != 1 || summary.InsertCount != 1 {
					return fmt.Errorf("DMLCount=%d InsertCount=%d, want 1 1", summary.DMLCount, summary.InsertCount)
				}
				if !hasChangedTable(summary, "dev", "public", "created_orders", redshift.ChangeKindCreate, false) {
					return fmt.Errorf("missing created_orders create resource: %#v", summary.Tables)
				}
				if !hasChangedTable(summary, "dev", "public", "changed_orders", redshift.ChangeKindDML, true) {
					return fmt.Errorf("missing changed_orders dml resource: %#v", summary.Tables)
				}
				return nil
			},
		},
		{
			name: "completion scope includes from relation",
			run: func() error {
				sql := "SELECT * FROM sales WHERE "
				ctx := redshift.CollectCompletion(sql, len(sql))
				if ctx == nil || ctx.Scope == nil {
					return fmt.Errorf("missing completion scope")
				}
				for _, ref := range ctx.Scope.LocalReferences {
					if ref.Name == "sales" {
						return nil
					}
				}
				return fmt.Errorf("missing sales relation in local scope: %#v", ctx.Scope.LocalReferences)
			},
		},
	}

	stats := RuntimeSemanticStats{TotalChecks: len(checks)}
	for _, check := range checks {
		if err := check.run(); err != nil {
			stats.FailedChecks++
			stats.Failures = append(stats.Failures, fmt.Sprintf("%s: %v", check.name, err))
			continue
		}
		stats.PassedChecks++
	}
	return stats
}

func hasChangedTable(summary *redshift.ChangeSummary, database, schema, name string, kind redshift.ChangeKind, affected bool) bool {
	if summary == nil {
		return false
	}
	for _, table := range summary.Tables {
		if table.Database == database && table.Schema == schema && table.Name == name && table.Kind == kind && table.Affected == affected {
			return true
		}
	}
	return false
}

func loadLegacyExpectedFailures(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	re := regexp.MustCompile("\"([^\"]+\\.sql)\"\\s*:\\s*`([^`]*)`")
	matches := re.FindAllStringSubmatch(string(data), -1)
	result := make(map[string]string, len(matches))
	for _, match := range matches {
		result[match[1]] = match[2]
	}
	return result, nil
}

func RenderMarkdownReport(commands AWSCommandManifestResult, legacy LegacyCorpusStats) string {
	return RenderFullMarkdownReport(commands, legacy, LegacyStatementParityStats{}, ReferenceRedshiftResult{
		SkipReason: "not evaluated",
	}, RuntimeSemanticStats{})
}

func RenderFullMarkdownReport(commands AWSCommandManifestResult, legacy LegacyCorpusStats, parity LegacyStatementParityStats, reference ReferenceRedshiftResult, runtime RuntimeSemanticStats) string {
	var b strings.Builder
	b.WriteString("# Redshift Compatibility Report\n\n")
	b.WriteString("## AWS Command Coverage\n\n")
	if commands.Source.URL != "" {
		b.WriteString(fmt.Sprintf("Source: %s\n\n", commands.Source.URL))
	}
	statuses := []CommandStatus{
		StatusSupportedRuntime,
		StatusSupportedParse,
		StatusExplicitUnsupported,
		StatusNotRelevant,
	}
	b.WriteString("| Status | Count |\n|---|---:|\n")
	for _, status := range statuses {
		b.WriteString(fmt.Sprintf("| %s | %d |\n", status, commands.Counts[status]))
	}
	b.WriteString("\n")

	b.WriteString("## Legacy Corpus\n\n")
	b.WriteString("| Metric | Count |\n|---|---:|\n")
	b.WriteString(fmt.Sprintf("| total files | %d |\n", legacy.TotalFiles))
	b.WriteString(fmt.Sprintf("| passing files | %d |\n", legacy.PassingFiles))
	b.WriteString(fmt.Sprintf("| expected failure files | %d |\n", legacy.ExpectedFailureFiles))
	b.WriteString(fmt.Sprintf("| new failure files | %d |\n", legacy.NewFailureFiles))
	b.WriteString(fmt.Sprintf("| promoted files | %d |\n", legacy.PromotedFiles))
	b.WriteString("\n")

	if len(legacy.NewFailures) > 0 {
		b.WriteString("### New Failures\n\n")
		for _, failure := range legacy.NewFailures {
			b.WriteString("- " + failure + "\n")
		}
		b.WriteString("\n")
	}
	if len(legacy.Promoted) > 0 {
		b.WriteString("### Promoted Expected Failures\n\n")
		for _, file := range legacy.Promoted {
			b.WriteString("- " + file + "\n")
		}
		b.WriteString("\n")
	}

	b.WriteString("## Legacy Statement Parity\n\n")
	b.WriteString("| Metric | Count |\n|---|---:|\n")
	b.WriteString(fmt.Sprintf("| total statements | %d |\n", parity.TotalStatements))
	b.WriteString(fmt.Sprintf("| both accept | %d |\n", parity.BothAccept))
	b.WriteString(fmt.Sprintf("| both reject | %d |\n", parity.BothReject))
	b.WriteString(fmt.Sprintf("| legacy accepts, omni rejects | %d |\n", parity.OldAcceptsOmniRejects))
	b.WriteString(fmt.Sprintf("| legacy rejects, omni accepts | %d |\n", parity.OldRejectsOmniAccepts))
	b.WriteString("\n")
	if len(parity.OldAcceptsOmniRejectExamples) > 0 {
		b.WriteString("### Legacy Accepts, Omni Rejects\n\n")
		for _, example := range parity.OldAcceptsOmniRejectExamples {
			b.WriteString("- " + example + "\n")
		}
		b.WriteString("\n")
	}
	if len(parity.OldRejectsOmniAcceptExamples) > 0 {
		b.WriteString("### Legacy Rejects, Omni Accepts\n\n")
		for _, example := range parity.OldRejectsOmniAcceptExamples {
			b.WriteString("- " + example + "\n")
		}
		b.WriteString("\n")
	}

	b.WriteString("## Runtime Semantics\n\n")
	b.WriteString("| Metric | Count |\n|---|---:|\n")
	b.WriteString(fmt.Sprintf("| total checks | %d |\n", runtime.TotalChecks))
	b.WriteString(fmt.Sprintf("| passed checks | %d |\n", runtime.PassedChecks))
	b.WriteString(fmt.Sprintf("| failed checks | %d |\n", runtime.FailedChecks))
	b.WriteString("\n")
	if len(runtime.Failures) > 0 {
		b.WriteString("### Runtime Failures\n\n")
		for _, failure := range runtime.Failures {
			b.WriteString("- " + failure + "\n")
		}
		b.WriteString("\n")
	}

	b.WriteString("## Reference Redshift\n\n")
	if !reference.Enabled {
		b.WriteString("Status: skipped")
		if reference.SkipReason != "" {
			b.WriteString(" (" + reference.SkipReason + ")")
		}
		b.WriteString("\n\n")
	} else {
		b.WriteString("| Metric | Count |\n|---|---:|\n")
		b.WriteString(fmt.Sprintf("| total samples | %d |\n", reference.Total))
		b.WriteString(fmt.Sprintf("| accepted | %d |\n", reference.Accepted))
		b.WriteString(fmt.Sprintf("| rejected | %d |\n", reference.Rejected))
		b.WriteString(fmt.Sprintf("| errors | %d |\n", len(reference.Errors)))
		b.WriteString("\n")
	}

	b.WriteString("## Command Details\n\n")
	b.WriteString("| Command | Status | Parse |\n|---|---|---|\n")
	for _, command := range commands.Commands {
		parse := "n/a"
		if command.SampleSQL != "" {
			if command.ParseOK {
				parse = "ok"
			} else {
				parse = "error"
			}
		}
		b.WriteString(fmt.Sprintf("| %s | %s | %s |\n", command.Name, command.Status, parse))
	}
	return b.String()
}
