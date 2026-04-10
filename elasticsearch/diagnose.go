package elasticsearch

// DiagnosticSeverity represents the severity level of a diagnostic.
type DiagnosticSeverity int

const (
	// SeverityError indicates a hard syntax error.
	SeverityError DiagnosticSeverity = 1
)

// DiagnosticRange represents a character range in source text.
type DiagnosticRange struct {
	Start Position
	End   Position
}

// Diagnostic holds a single parse diagnostic (error, warning, …) produced by
// Diagnose. Positions are 0-based line and column values matching LSP
// conventions so callers can forward them without conversion.
type Diagnostic struct {
	Range    DiagnosticRange
	Severity DiagnosticSeverity
	Message  string
}

// Diagnose parses statement and returns a Diagnostic for every SyntaxError
// found. It returns nil, nil when there are no errors.
func Diagnose(statement string) ([]Diagnostic, error) {
	parseResult, _ := ParseElasticsearchREST(statement)
	if parseResult == nil {
		return nil, nil
	}

	var diagnostics []Diagnostic
	for _, err := range parseResult.Errors {
		if err == nil {
			continue
		}
		// SyntaxError positions are 0-based (line and column). The end of the
		// range is set one column past the start so callers get a non-empty span.
		start := err.Position
		end := Position{Line: start.Line, Column: start.Column + 1}
		diagnostics = append(diagnostics, Diagnostic{
			Range:    DiagnosticRange{Start: start, End: end},
			Severity: SeverityError,
			Message:  err.Message,
		})
	}
	return diagnostics, nil
}
