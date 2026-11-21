package analysis

import (
	"os"
	"slices"
	"testing"

	"github.com/bytebase/omni/oracle/ast"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

type statementTypeCorpusCase struct {
	Statement string   `yaml:"statement"`
	Want      []string `yaml:"want"`
}

var statementTypeNames = map[StatementType]string{
	StatementTypeUnspecified: "STATEMENT_TYPE_UNSPECIFIED",
	CreateDatabase:           "CREATE_DATABASE",
	CreateTable:              "CREATE_TABLE",
	CreateIndex:              "CREATE_INDEX",
	CreateView:               "CREATE_VIEW",
	CreateSequence:           "CREATE_SEQUENCE",
	CreateSchema:             "CREATE_SCHEMA",
	CreateFunction:           "CREATE_FUNCTION",
	CreateTrigger:            "CREATE_TRIGGER",
	CreateProcedure:          "CREATE_PROCEDURE",
	CreateType:               "CREATE_TYPE",
	DropDatabase:             "DROP_DATABASE",
	DropTable:                "DROP_TABLE",
	DropView:                 "DROP_VIEW",
	DropIndex:                "DROP_INDEX",
	DropSequence:             "DROP_SEQUENCE",
	DropFunction:             "DROP_FUNCTION",
	DropTrigger:              "DROP_TRIGGER",
	DropProcedure:            "DROP_PROCEDURE",
	DropType:                 "DROP_TYPE",
	AlterDatabase:            "ALTER_DATABASE",
	AlterTable:               "ALTER_TABLE",
	AlterIndex:               "ALTER_INDEX",
	AlterView:                "ALTER_VIEW",
	AlterSequence:            "ALTER_SEQUENCE",
	AlterType:                "ALTER_TYPE",
	Truncate:                 "TRUNCATE",
	Rename:                   "RENAME",
	Comment:                  "COMMENT",
	Insert:                   "INSERT",
	Update:                   "UPDATE",
	Delete:                   "DELETE",
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
			list, err := Parse(tc.Statement)
			require.NoError(t, err)
			require.NotEmpty(t, list.Items)

			seen := make(map[string]bool)
			for _, item := range list.Items {
				raw, ok := item.(*ast.RawStmt)
				require.True(t, ok)
				name, ok := statementTypeNames[StatementTypeOf(raw.Stmt)]
				require.True(t, ok, "statement type %d has no corpus name", StatementTypeOf(raw.Stmt))
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
