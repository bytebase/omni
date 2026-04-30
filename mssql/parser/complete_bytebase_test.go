package parser_test

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/mssql/parser"
)

func collectMarked(t *testing.T, input string) (string, int, *parser.CandidateSet) {
	t.Helper()
	cursor := strings.Index(input, "|")
	if cursor < 0 {
		t.Fatalf("input %q has no cursor marker", input)
	}
	sql := strings.Replace(input, "|", "", 1)
	return sql, cursor, parser.Collect(sql, cursor)
}

func collectMarkedWithPrefixRetry(t *testing.T, input string) (string, int, *parser.CandidateSet) {
	t.Helper()
	sql, cursor, candidates := collectMarked(t, input)
	if len(candidates.Rules) > 0 {
		return sql, cursor, candidates
	}

	for _, tok := range parser.Tokenize(sql) {
		if tok.Loc < cursor && cursor <= tok.End && parser.IsIdentTokenType(tok.Type) {
			return sql, cursor, parser.Collect(sql, tok.Loc)
		}
		if tok.End == cursor && parser.IsIdentTokenType(tok.Type) {
			return sql, cursor, parser.Collect(sql, tok.Loc)
		}
	}
	return sql, cursor, candidates
}

func requireRule(t *testing.T, candidates *parser.CandidateSet, rule string) {
	t.Helper()
	if !candidates.HasRule(rule) {
		t.Fatalf("missing rule %q; got rules=%v tokens=%v", rule, candidates.Rules, tokenNames(candidates.Tokens))
	}
}

func requireToken(t *testing.T, candidates *parser.CandidateSet, token int) {
	t.Helper()
	if !candidates.HasToken(token) {
		t.Fatalf("missing token %q; got rules=%v tokens=%v", parser.TokenName(token), candidates.Rules, tokenNames(candidates.Tokens))
	}
}

func tokenNames(tokens []int) []string {
	result := make([]string, 0, len(tokens))
	for _, tok := range tokens {
		name := parser.TokenName(tok)
		if name == "" {
			name = string(rune(tok))
		}
		result = append(result, name)
	}
	return result
}

func TestCompletionBytebaseTableReferenceSignals(t *testing.T) {
	tests := []string{
		"select count(1) from t1 where id; SELECT * FROM |",
		"SELECT 1\nGO\nSELECT * FROM |",
		"SELECT * FROM dbo.|",
		"SELECT * FROM MySchema.|",
		"WITH MyCTE_01 AS (SELECT * FROM dbo.Employees) SELECT * FROM |",
		"INSERT INTO |",
		"UPDATE | SET Name = 'x'",
		"DELETE FROM |",
		"MERGE INTO target USING |",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			_, _, candidates := collectMarked(t, input)
			requireRule(t, candidates, "table_ref")
		})
	}
}

func TestCompletionBytebaseColumnReferenceSignals(t *testing.T) {
	tests := []string{
		"SELECT | FROM Employees",
		"SELECT tableAlias.| FROM Employees AS tableAlias",
		"WITH MyCTE_01 AS (SELECT * FROM dbo.Employees) SELECT MyCTE_01.| FROM MyCTE_01",
		"SELECT * FROM Employees e JOIN Address a ON |",
		"SELECT * FROM Employees e LEFT JOIN MySchema.SalaryLevel s ON s.|",
		"SELECT Id AS IdAlias, Name FROM Employees ORDER BY |",
		"SELECT Id, Name FROM Employees WHERE |",
		"SELECT * FROM Employees WHERE Id IN (1, |)",
		"SELECT * FROM Employees WHERE NOT |",
		"SELECT * FROM Employees WHERE Id + | > 0",
		"SELECT Id FROM Employees GROUP BY |",
		"SELECT Id FROM Employees HAVING |",
		"SELECT COALESCE(NULL, |) FROM Employees",
		"SELECT CONVERT(INT, |) FROM Employees",
		"INSERT INTO Employees SELECT | FROM Address",
		"UPDATE Employees SET |",
		"UPDATE Employees SET Name = |",
		"DELETE FROM Employees WHERE |",
		"SELECT ABS(|) FROM Employees",
		"SELECT CONCAT(Name, |) FROM Employees",
		"SELECT * FROM (SELECT Name, Amount FROM Sales) src PIVOT (SUM(Amount) FOR Name IN (|)) AS p",
		"SELECT * FROM (SELECT Name, Amount FROM Sales) src UNPIVOT (Amount FOR Name IN (|)) AS u",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			_, _, candidates := collectMarked(t, input)
			requireRule(t, candidates, "columnref")
		})
	}
}

func TestCompletionBytebaseSequenceReferenceSignals(t *testing.T) {
	tests := []string{
		"SELECT NEXT VALUE FOR |",
		"INSERT INTO Orders(Id) VALUES (NEXT VALUE FOR |)",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			_, _, candidates := collectMarked(t, input)
			requireRule(t, candidates, "sequence_ref")
		})
	}
}

func TestCompletionBytebaseProcedureReferenceSignals(t *testing.T) {
	tests := []string{
		"EXEC |",
		"EXECUTE |",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			_, _, candidates := collectMarked(t, input)
			requireRule(t, candidates, "proc_ref")
		})
	}
}

func TestCompletionBytebaseOpenJSONWithTypeSignals(t *testing.T) {
	tests := []string{
		"SELECT * FROM OPENJSON(@json) WITH (Name |)",
		"SELECT * FROM OPENJSON(@json) WITH (Name nvarchar(100), Age |)",
		"CREATE PROCEDURE p @Name | AS SELECT 1",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			_, _, candidates := collectMarked(t, input)
			requireRule(t, candidates, "type_name")
		})
	}
}

func TestCompletionBytebaseIncompleteBracketIdentifierSignals(t *testing.T) {
	_, _, candidates := collectMarked(t, "SELECT * FROM [|")
	requireRule(t, candidates, "table_ref")
}

func TestCompletionBytebaseAsteriskQualifierSignals(t *testing.T) {
	tests := []string{
		"SELECT |.* FROM Employees AS e",
		"WITH MyCTE_01 AS (SELECT * FROM dbo.Employees) SELECT |.* FROM MyCTE_01 JOIN dbo.Address ON MyCTE_01.EmployeeID = dbo.Address.EmployeeID",
		"SELECT e.|* FROM Employees AS e",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			_, _, candidates := collectMarked(t, input)
			requireRule(t, candidates, "columnref")
		})
	}
}

func TestCompletionBytebasePrefixRetrySignals(t *testing.T) {
	tests := []struct {
		input string
		rule  string
	}{
		{input: "SELECT * FROM dbo.Us|", rule: "table_ref"},
		{input: "SELECT * FROM Emp|", rule: "table_ref"},
		{input: "SELECT tableAlias.Na| FROM Employees AS tableAlias", rule: "columnref"},
		{input: "SELECT Na| FROM Employees", rule: "columnref"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			_, _, candidates := collectMarkedWithPrefixRetry(t, tt.input)
			requireRule(t, candidates, tt.rule)
		})
	}
}

func TestCompletionIncompleteSQLSignals(t *testing.T) {
	tests := []struct {
		input string
		rule  string
		token int
	}{
		{input: "SELECT * FROM|", rule: "table_ref"},
		{input: "SELECT * FROM |", rule: "table_ref"},
		{input: "SELECT * FROM dbo.|", rule: "table_ref"},
		{input: "SELECT * FROM Employees JOIN|", rule: "table_ref"},
		{input: "SELECT Id,|", rule: "columnref"},
		{input: "SELECT Id, |", rule: "columnref"},
		{input: "SELECT * FROM Employees WHERE Id =|", rule: "columnref"},
		{input: "SELECT * FROM Employees WHERE Id = |", rule: "columnref"},
		{input: "SELECT Id FROM Employees ORDER BY|", rule: "columnref"},
		{input: "SELECT Id FROM Employees ORDER BY |", rule: "columnref"},
		{input: "WITH cte AS (|", token: parser.SELECT},
		{input: "WITH cte AS (SELECT Id FROM dbo.Employees) SELECT * FROM|", rule: "table_ref"},
		{input: "ALTER TABLE Employees ADD CONSTRAINT fk FOREIGN KEY (Id) REFERENCES |", rule: "table_ref"},
		{input: "DROP VIEW MySchema.|", rule: "table_ref"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			_, _, candidates := collectMarked(t, tt.input)
			if tt.rule != "" {
				requireRule(t, candidates, tt.rule)
			}
			if tt.token != 0 {
				requireToken(t, candidates, tt.token)
			}
		})
	}
}

func TestCompletionIncompleteSQLDoesNotFallBackToOnlyTopLevel(t *testing.T) {
	tests := []string{
		"SELECT * FROM|",
		"SELECT Id,|",
		"SELECT * FROM Employees WHERE Id =|",
		"SELECT Id FROM Employees ORDER BY|",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			_, _, candidates := collectMarked(t, input)
			if len(candidates.Rules) == 0 {
				t.Fatalf("expected grammar rule candidates, got only tokens=%v", tokenNames(candidates.Tokens))
			}
		})
	}
}
