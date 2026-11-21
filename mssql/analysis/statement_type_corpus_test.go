package analysis

import (
	"os"
	"slices"
	"testing"

	"github.com/bytebase/omni/mssql"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

type statementTypeCorpusCase struct {
	Statement string   `yaml:"statement"`
	Want      []string `yaml:"want"`
}

var statementTypeNames = map[StatementType]string{
	StatementTypeUnspecified: "STATEMENT_TYPE_UNSPECIFIED",
	AlterDatabase:            "ALTER_DATABASE",
	AlterIndex:               "ALTER_INDEX",
	AlterTable:               "ALTER_TABLE",
	CreateDatabase:           "CREATE_DATABASE",
	CreateIndex:              "CREATE_INDEX",
	CreateSchema:             "CREATE_SCHEMA",
	CreateTable:              "CREATE_TABLE",
	CreateView:               "CREATE_VIEW",
	DropDatabase:             "DROP_DATABASE",
	DropIndex:                "DROP_INDEX",
	DropSchema:               "DROP_SCHEMA",
	DropTable:                "DROP_TABLE",
	DropView:                 "DROP_VIEW",
	Truncate:                 "TRUNCATE",
	Delete:                   "DELETE",
	Insert:                   "INSERT",
	Update:                   "UPDATE",
	Execute:                  "EXECUTE",
}

func TestStatementTypeCorpus(t *testing.T) {
	content, err := os.ReadFile("testdata/statement-type.yaml")
	require.NoError(t, err)
	var cases []statementTypeCorpusCase
	require.NoError(t, yaml.Unmarshal(content, &cases))
	require.NotEmpty(t, cases)

	for _, tc := range cases {
		t.Run(tc.Statement, func(t *testing.T) {
			stmts, err := mssql.Parse(tc.Statement)
			require.NoError(t, err)
			require.NotEmpty(t, stmts)

			seen := make(map[string]bool)
			for _, s := range stmts {
				if s.Empty() {
					continue
				}
				name, ok := statementTypeNames[StatementTypeOf(s.AST)]
				require.True(t, ok, "statement type %d has no corpus name", StatementTypeOf(s.AST))
				seen[name] = true
			}
			got := make([]string, 0, len(seen))
			for name := range seen {
				got = append(got, name)
			}
			slices.Sort(got)
			want := slices.Clone(tc.Want)
			slices.Sort(want)
			require.Equal(t, want, got)
		})
	}
}
