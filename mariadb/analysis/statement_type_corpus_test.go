package analysis

import (
	"os"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/bytebase/omni/mariadb/parser"
)

// The corpus in testdata/statement-type.yaml is Bytebase's MariaDB statement-type
// fixture, recorded from the classifier the release check gated SDL statements
// on. It is reproduced here, never re-recorded against this implementation.

type statementTypeCorpusCase struct {
	Statement string   `yaml:"statement"`
	Want      []string `yaml:"want"`
}

// statementTypeNames is the corpus's vocabulary: storepb.StatementType names,
// which this package's StatementType values mirror.
var statementTypeNames = map[StatementType]string{
	StatementTypeUnspecified: "STATEMENT_TYPE_UNSPECIFIED",
	CreateDatabase:           "CREATE_DATABASE",
	CreateTable:              "CREATE_TABLE",
	CreateView:               "CREATE_VIEW",
	CreateIndex:              "CREATE_INDEX",
	CreateEvent:              "CREATE_EVENT",
	CreateTrigger:            "CREATE_TRIGGER",
	CreateFunction:           "CREATE_FUNCTION",
	CreateProcedure:          "CREATE_PROCEDURE",
	DropDatabase:             "DROP_DATABASE",
	DropTable:                "DROP_TABLE",
	DropView:                 "DROP_VIEW",
	DropIndex:                "DROP_INDEX",
	DropEvent:                "DROP_EVENT",
	DropTrigger:              "DROP_TRIGGER",
	DropFunction:             "DROP_FUNCTION",
	DropProcedure:            "DROP_PROCEDURE",
	AlterDatabase:            "ALTER_DATABASE",
	AlterTable:               "ALTER_TABLE",
	AlterView:                "ALTER_VIEW",
	AlterEvent:               "ALTER_EVENT",
	Truncate:                 "TRUNCATE",
	Rename:                   "RENAME",
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
		t.Run(firstLine(tc.Statement), func(t *testing.T) {
			list, err := parser.Parse(tc.Statement)
			require.NoError(t, err)
			require.Len(t, list.Items, 1)

			name, ok := statementTypeNames[StatementTypeOf(list.Items[0])]
			require.True(t, ok, "statement type %d has no corpus name", StatementTypeOf(list.Items[0]))
			want := slices.Clone(tc.Want)
			slices.Sort(want)
			require.Equal(t, want, []string{name})
		})
	}
}

func firstLine(s string) string {
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			return s[:i]
		}
	}
	return s
}
