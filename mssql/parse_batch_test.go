package mssql

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bytebase/omni/mssql/ast"
)

// procBatchSQL builds the customer-shaped fixture from BYT-10207: an optional
// preceding statement (with or without a GO batch separator), Spanish
// accented comments, and a CREATE OR ALTER PROCEDURE whose body contains ';'
// terminators. The prefix is padded so the procedure body's "SET NOCOUNT ON;"
// lands on line 98, matching the customer's report.
func procBatchSQL(priorStmt, withGO bool) string {
	var b strings.Builder
	if priorStmt {
		b.WriteString("ALTER TABLE dbo.STG_Ventas ADD STS_STG INT NULL;\n")
		if withGO {
			b.WriteString("GO\n")
		}
	}
	for strings.Count(b.String(), "\n") < 91 {
		b.WriteString("-- relleno de línea con acentos: configuración, número, año\n")
	}
	b.WriteString(`-- Filas sin match activo (STS_STG = 3). Para alertas en ADF.
CREATE OR ALTER PROCEDURE dbo.usp_Procesar_Ventas
    @Fecha DATE,
    @Región NVARCHAR(50) = NULL
AS
BEGIN
    SET NOCOUNT ON;
    DECLARE @cnt INT = 0;
    -- Actualización de filas sin coincidencia
    UPDATE s SET s.STS_STG = 3
    FROM dbo.STG_Ventas s
    LEFT JOIN dbo.DIM_Tienda t ON t.Id = s.TiendaId
    WHERE t.Id IS NULL AND s.Fecha = @Fecha;
    SELECT @cnt = @@ROWCOUNT;
    IF @cnt > 0
    BEGIN
        INSERT INTO dbo.LOG_Proceso (Fecha, Mensaje) VALUES (@Fecha, N'Filas sin match: ' + CAST(@cnt AS NVARCHAR(10)));
    END;
    RETURN 0;
END;
`)
	return b.String()
}

// TestParseProcedureAfterStatementWithoutGO pins the batch-boundary behavior
// behind BYT-10207. The legacy ANTLR splitter treated the ';' terminators
// inside the procedure body as statement boundaries whenever a preceding
// statement was not separated by GO, cutting the procedure into fragments
// that then failed to parse standalone. Parse must keep the procedure whole
// regardless of whether the preceding statement is followed by GO.
func TestParseProcedureAfterStatementWithoutGO(t *testing.T) {
	tests := []struct {
		name      string
		priorStmt bool
		withGO    bool
		wantASTs  []ast.Node
	}{
		{
			name:     "procedure alone",
			wantASTs: []ast.Node{&ast.CreateProcedureStmt{}},
		},
		{
			name:      "prior statement without GO",
			priorStmt: true,
			wantASTs:  []ast.Node{&ast.AlterTableStmt{}, &ast.CreateProcedureStmt{}},
		},
		{
			name:      "prior statement with GO",
			priorStmt: true,
			withGO:    true,
			wantASTs:  []ast.Node{&ast.AlterTableStmt{}, &ast.GoStmt{}, &ast.CreateProcedureStmt{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql := procBatchSQL(tt.priorStmt, tt.withGO)

			// Fixture invariant: the customer's error pointed at line 98.
			lines := strings.Split(sql, "\n")
			if got := strings.TrimSpace(lines[97]); got != "SET NOCOUNT ON;" {
				t.Fatalf("fixture line 98 = %q, want SET NOCOUNT ON;", got)
			}

			stmts, err := Parse(sql)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}
			if len(stmts) != len(tt.wantASTs) {
				for i, s := range stmts {
					t.Logf("stmt[%d] %T %d:%d-%d:%d", i, s.AST, s.Start.Line, s.Start.Column, s.End.Line, s.End.Column)
				}
				t.Fatalf("got %d statements, want %d", len(stmts), len(tt.wantASTs))
			}
			for i, want := range tt.wantASTs {
				if gotT, wantT := typeName(stmts[i].AST), typeName(want); gotT != wantT {
					t.Errorf("stmt[%d] AST = %s, want %s", i, gotT, wantT)
				}
			}

			// The procedure must be the last statement, start at its CREATE
			// keyword on line 93, and carry its whole body through END;.
			proc := stmts[len(stmts)-1]
			if proc.Start.Line != 93 || proc.Start.Column != 1 {
				t.Errorf("procedure Start = %d:%d, want 93:1", proc.Start.Line, proc.Start.Column)
			}
			if !strings.HasPrefix(strings.TrimSpace(stripLeadingComments(proc.Text)), "CREATE OR ALTER PROCEDURE") {
				t.Errorf("procedure Text does not begin at CREATE OR ALTER PROCEDURE: %q", head(proc.Text))
			}
			if !strings.HasSuffix(strings.TrimSpace(proc.Text), "END;") {
				t.Errorf("procedure Text does not end at END;: %q", tail(proc.Text))
			}
			if proc.ByteEnd != len(sql)-1 {
				t.Errorf("procedure ByteEnd = %d, want %d (end of input minus trailing newline)", proc.ByteEnd, len(sql)-1)
			}
		})
	}
}

func typeName(n ast.Node) string {
	return fmt.Sprintf("%T", n)
}

// stripLeadingComments drops leading blank lines and `--` comment lines.
func stripLeadingComments(text string) string {
	lines := strings.Split(text, "\n")
	for len(lines) > 0 {
		l := strings.TrimSpace(lines[0])
		if l == "" || strings.HasPrefix(l, "--") {
			lines = lines[1:]
			continue
		}
		break
	}
	return strings.Join(lines, "\n")
}

func head(s string) string {
	if len(s) > 80 {
		return s[:80]
	}
	return s
}

func tail(s string) string {
	if len(s) > 80 {
		return s[len(s)-80:]
	}
	return s
}
