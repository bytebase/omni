package analysis

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

type splitCorpusCase struct {
	Description string               `yaml:"description"`
	Input       string               `yaml:"input"`
	Result      []splitCorpusSegment `yaml:"result"`
}

type splitCorpusSegment struct {
	Text  string           `yaml:"text"`
	Range splitCorpusRange `yaml:"range"`
	Empty bool             `yaml:"empty"`
}

type splitCorpusRange struct {
	Start int `yaml:"start"`
	End   int `yaml:"end"`
}

func TestSplitCorpus(t *testing.T) {
	content, err := os.ReadFile("testdata/split.yaml")
	require.NoError(t, err)
	var cases []splitCorpusCase
	require.NoError(t, yaml.Unmarshal(content, &cases))
	require.NotEmpty(t, cases)

	for _, tc := range cases {
		t.Run(tc.Description, func(t *testing.T) {
			stmts := SplitStatements(tc.Input)
			require.Len(t, stmts, len(tc.Result))
			for i, want := range tc.Result {
				got := stmts[i]
				require.Equal(t, want.Text, got.Text, "statement %d text", i)
				require.Equal(t, want.Empty, got.Empty, "statement %d empty", i)
				require.Equal(t, want.Range.Start, got.Start, "statement %d start", i)
				require.Equal(t, want.Range.End, got.End, "statement %d end", i)
			}
		})
	}
}

func filterNonEmpty(stmts []Statement) []Statement {
	var out []Statement
	for _, s := range stmts {
		if !s.Empty {
			out = append(out, s)
		}
	}
	return out
}

func TestSplitStoredUnitWithoutSlashSeparator(t *testing.T) {
	statement := `CREATE FUNCTION calc_bonus(p_start_date DATE)
RETURN DATE
IS
  v_current_date DATE := p_start_date;
BEGIN
  RETURN v_current_date;
END calc_bonus;
CREATE PROCEDURE update_salary(p_employee_id NUMBER)
IS
BEGIN
  UPDATE employees SET salary = salary + 1 WHERE id = p_employee_id;
END update_salary;
CREATE TABLE t(id NUMBER);`

	statements := filterNonEmpty(SplitStatements(statement))
	require.Len(t, statements, 3)
	require.Equal(t, `CREATE FUNCTION calc_bonus(p_start_date DATE)
RETURN DATE
IS
  v_current_date DATE := p_start_date;
BEGIN
  RETURN v_current_date;
END calc_bonus;`, statements[0].Text)
	require.Equal(t, `
CREATE PROCEDURE update_salary(p_employee_id NUMBER)
IS
BEGIN
  UPDATE employees SET salary = salary + 1 WHERE id = p_employee_id;
END update_salary;`, statements[1].Text)
	require.Equal(t, `
CREATE TABLE t(id NUMBER)`, statements[2].Text)
}

func TestSplitKeepsEmbeddedSourceWhole(t *testing.T) {
	for _, tc := range []struct {
		name   string
		source string
		suffix string
	}{
		{
			name:   "Java source",
			source: "CREATE OR REPLACE AND COMPILE JAVA SOURCE NAMED \"Dec\" AS\npublic class Dec { public static int f(int n) { return n--; } }",
			suffix: "\n/\n",
		},
		{
			name:   "MLE module",
			source: "CREATE OR REPLACE MLE MODULE dec_module LANGUAGE JAVASCRIPT AS\nexport function f(n) { return n--; }",
			suffix: "\n/\n",
		},
		{
			name:   "inline MLE call spec",
			source: "CREATE OR REPLACE FUNCTION dec_js(n NUMBER) RETURN NUMBER\nAS MLE LANGUAGE JAVASCRIPT\n{{ return n--; }};",
			suffix: "\n",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			statements := filterNonEmpty(SplitStatements(tc.source + tc.suffix + "SELECT 1 FROM dual;"))
			require.Len(t, statements, 2)
			require.Equal(t, tc.source, statements[0].Text)
			require.Equal(t, "SELECT 1 FROM dual", strings.TrimSpace(statements[1].Text))
		})
	}
}

func TestSplitPreservesSQLSetStatements(t *testing.T) {
	statement := `SET DEFINE OFF
SET TRANSACTION READ ONLY;
SET ROLE app_role;
SELECT 1 FROM dual;`

	statements := filterNonEmpty(SplitStatements(statement))
	require.Len(t, statements, 3)
	require.Equal(t, "SET TRANSACTION READ ONLY", statements[0].Text)
	require.Equal(t, "\nSET ROLE app_role", statements[1].Text)
	require.Equal(t, "\nSELECT 1 FROM dual", statements[2].Text)
}

func TestSplitLargeInsertScriptScalesLinearly(t *testing.T) {
	const rowCount = 2000
	padding := strings.Repeat("x", 1024)
	var builder strings.Builder
	for i := range rowCount {
		fmt.Fprintf(&builder, "INSERT INTO perf_omni_oracle (id, payload) VALUES (%d, '%s');\n", i, padding)
	}

	started := time.Now()
	statements := SplitStatements(builder.String())
	elapsed := time.Since(started)

	require.Len(t, statements, rowCount)
	require.Less(t, elapsed, time.Second)
}
