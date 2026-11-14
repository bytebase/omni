package compat

import (
	"errors"
	"strings"
	"testing"
)

// sqlServerErrors is a stand-in for go-mssqldb's error: what SQL Server
// reported for a batch.
type sqlServerErrors []SQLServerError

func (e sqlServerErrors) Error() string {
	if len(e) == 0 {
		return "mssql: no error"
	}
	return e[0].Message
}

func readSQLServer(err error) ([]SQLServerError, bool) {
	var e sqlServerErrors
	if !errors.As(err, &e) {
		return nil, false
	}
	return e, true
}

func server(number int32, message string) error {
	return sqlServerErrors{{number, message}}
}

// TestParseOnlyVerdictReadsTheFirstError pins how a verdict reads the
// answer: every error is the parser's and a rejection, but a lookup of what
// the server holds. The answers are those of SQL Server 2022 CU27 GDR
// (16.0.4305.5) under SET PARSEONLY ON.
func TestParseOnlyVerdictReadsTheFirstError(t *testing.T) {
	never := func(string) error {
		t.Fatalf("a verdict that names no variable re-asked")
		return nil
	}
	for _, c := range []struct {
		name string
		err  error
		want string
	}{
		{"parsed", nil, Accept},
		{"a syntax error", server(102, "Incorrect syntax near 'FROM'."), Reject},
		{"near a keyword", server(156, "Incorrect syntax near the keyword 'FROM'."), Reject},
		{"refused by the parser", server(135, "Cannot use a BREAK statement outside the scope of a WHILE statement."), Reject},
		{"an option it does not know", server(195, "'a' is not a recognized SET option."), Reject},
		{"a statement that must start its batch", server(111, "'CREATE VIEW' must be the first statement in a query batch."), Reject},
		{"a database the server lacks", server(911, "Database 'a' does not exist. Make sure that the name is entered correctly."), Accept},
		{"an object the server lacks", server(15151, "Cannot find the object 'a', because it does not exist or you do not have permission."), Accept},
		{"a collation the server lacks", server(448, "Invalid collation 'a'."), Accept},
		{"the first error decides", sqlServerErrors{
			{156, "Incorrect syntax near the keyword 'FROM'."},
			{undeclaredScalar, `Must declare the scalar variable "@x".`},
		}, Reject},
	} {
		t.Run(c.name, func(t *testing.T) {
			got, err := ParseOnlyVerdict(c.err, readSQLServer, never)
			if err != nil {
				t.Fatalf("ParseOnlyVerdict: %v", err)
			}
			if got != c.want {
				t.Errorf("verdict %q, want %q", got, c.want)
			}
		})
	}
}

// TestParseOnlyVerdictDeclaresVariables pins the re-asks for a variable the
// statement does not declare: each declares the variable the last answer
// named, of the kind it needs, and the answer after them decides.
func TestParseOnlyVerdictDeclaresVariables(t *testing.T) {
	for _, c := range []struct {
		name string
		err  error
		// answers maps the declarations of each re-ask to its answer.
		answers map[string]error
		want    string
	}{
		{
			name:    "a scalar",
			err:     server(undeclaredScalar, `Must declare the scalar variable "@a".`),
			answers: map[string]error{"DECLARE @a INT;\n": nil},
			want:    Accept,
		},
		{
			name: "a scalar before a syntax error",
			err:  server(undeclaredScalar, `Must declare the scalar variable "@a".`),
			answers: map[string]error{
				"DECLARE @a INT;\n": server(102, "Incorrect syntax near 'zzz'."),
			},
			want: Reject,
		},
		{
			name:    "a table",
			err:     server(undeclaredTable, `Must declare the table variable "@a".`),
			answers: map[string]error{"DECLARE @a TABLE (a INT);\n": nil},
			want:    Accept,
		},
		{
			name: "a cursor",
			err:  server(undeclaredScalar, `Must declare the scalar variable "@a".`),
			answers: map[string]error{
				"DECLARE @a INT;\n":    server(notCursor, "The variable '@a' is not a cursor variable, but it is used in a place where a cursor variable is expected."),
				"DECLARE @a CURSOR;\n": nil,
			},
			want: Accept,
		},
		{
			name: "two variables",
			err:  server(undeclaredScalar, `Must declare the scalar variable "@a".`),
			answers: map[string]error{
				"DECLARE @a INT;\n":                            server(undeclaredTable, `Must declare the table variable "@b".`),
				"DECLARE @a INT;\nDECLARE @b TABLE (a INT);\n": nil,
			},
			want: Accept,
		},
		{
			name: "in a statement that must start its batch",
			err:  server(undeclaredScalar, `Must declare the scalar variable "@a".`),
			answers: map[string]error{
				"DECLARE @a INT;\n": server(111, "'CREATE/ALTER PROCEDURE' must be the first statement in a query batch."),
			},
			want: Reject,
		},
		{
			name: "one variable as two kinds",
			err:  server(undeclaredScalar, `Must declare the scalar variable "@a".`),
			answers: map[string]error{
				"DECLARE @a INT;\n": server(undeclaredTable, `Must declare the table variable "@a".`),
			},
			want: Reject,
		},
		{
			name: "still undeclared once declared",
			err:  server(undeclaredScalar, `Must declare the scalar variable "@a".`),
			answers: map[string]error{
				"DECLARE @a INT;\n": server(undeclaredScalar, `Must declare the scalar variable "@a".`),
			},
			want: Reject,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			declared := func(decls string) error {
				a, ok := c.answers[decls]
				if !ok {
					t.Fatalf("re-asked with %q", decls)
				}
				return a
			}
			got, err := ParseOnlyVerdict(c.err, readSQLServer, declared)
			if err != nil {
				t.Fatalf("ParseOnlyVerdict: %v", err)
			}
			if got != c.want {
				t.Errorf("verdict %q, want %q", got, c.want)
			}
		})
	}
}

// TestParseOnlyVerdictNeedsAServerAnswer fails an answer that is no
// verdict rather than read it as a rejection.
func TestParseOnlyVerdictNeedsAServerAnswer(t *testing.T) {
	never := func(string) error { return nil }
	if _, err := ParseOnlyVerdict(errors.New("connection reset"), readSQLServer, never); err == nil {
		t.Errorf("an error the server did not send gave a verdict")
	}
	if _, err := ParseOnlyVerdict(sqlServerErrors{}, readSQLServer, never); err == nil {
		t.Errorf("an error with no error in it gave a verdict")
	}
	if _, err := ParseOnlyVerdict(server(undeclaredScalar, "Must declare the scalar variable."), readSQLServer, never); err == nil {
		t.Errorf("an undeclared variable the error does not name gave a verdict")
	}
	n := 0
	endless := func(string) error {
		n++
		return server(undeclaredScalar, `Must declare the scalar variable "@v`+strings.Repeat("x", n)+`".`)
	}
	if _, err := ParseOnlyVerdict(endless(""), readSQLServer, endless); err == nil {
		t.Errorf("a server naming a new variable on every re-ask gave a verdict")
	}
}
