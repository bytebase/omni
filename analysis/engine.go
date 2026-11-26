package analysis

// An engine analysis package exports one value, Engine, implementing each
// capability interface below that the engine supports. Bytebase registers that
// value per engine and converts these types once, so every engine-specific
// decision behind a capability (where a statement's text ends, which segments a
// parse keeps, how a syntax error reads, the fallback a result limit takes) is
// the engine package's.

// Position is a 1-based line and column in a script, the column counted in
// code points.
type Position struct {
	Line   int
	Column int
}

// SyntaxError is a statement the engine rejects, positioned in the script the
// caller passed.
type SyntaxError struct {
	Position Position
	// Message is the error as a user reads it.
	Message string
	// RawMessage is the parser's own message, without position or context.
	// Empty means Message.
	RawMessage string
}

func (e *SyntaxError) Error() string { return e.Message }

// ParsedStatement is one statement of a parsed script.
type ParsedStatement struct {
	Statement
	// Node is the engine's AST node for the statement. It is nil for an empty
	// segment, and an engine whose parse tolerates a statement it cannot parse
	// says so in its Parse documentation.
	Node any
}

// DiagnosticSeverity mirrors the LSP severities.
type DiagnosticSeverity int

const (
	SeverityError       DiagnosticSeverity = 1
	SeverityWarning     DiagnosticSeverity = 2
	SeverityInformation DiagnosticSeverity = 3
	SeverityHint        DiagnosticSeverity = 4
)

// Diagnostic is one problem the editor marks in a script.
type Diagnostic struct {
	Start Position
	// End is exclusive. A zero End marks the one character at Start.
	End      Position
	Severity DiagnosticSeverity
	Message  string
}

// LSPPosition is a 0-based line and a 0-based character offset counted in
// UTF-16 code units, the position the editor protocol uses.
type LSPPosition struct {
	Line      uint32
	Character uint32
}

// LSPRange is a half-open range of LSP positions.
type LSPRange struct {
	Start LSPPosition
	End   LSPPosition
}

// ExplainFormat is the plan format a caller asks for.
type ExplainFormat int

const (
	// ExplainFormatDefault leaves the format to the engine.
	ExplainFormatDefault ExplainFormat = iota
	ExplainFormatText
	ExplainFormatJSON
	ExplainFormatXML
	ExplainFormatYAML
)

// Splitter cuts a script into statements, valid SQL or not. It never fails.
type Splitter interface {
	// Split returns the script's segments in order. Each segment's text is
	// script[Start:End].
	Split(script string) []Statement
}

// Parser parses a script into statements.
type Parser interface {
	// Parse returns the script's statements, or a *SyntaxError for the first
	// statement the engine rejects.
	Parse(script string) ([]ParsedStatement, error)
}

// StatementClassifier names what a parsed statement does.
type StatementClassifier interface {
	// StatementTypes returns the types of one statement Parse returned. An
	// empty result records nothing for the statement.
	StatementTypes(statement ParsedStatement) []StatementType
}

// StatementRanger finds the statements an editor highlights.
type StatementRanger interface {
	StatementRanges(script string) []LSPRange
}

// Diagnoser reports the problems an editor marks in a script.
type Diagnoser interface {
	Diagnose(script string) []Diagnostic
}

// QueryValidator decides whether a script may run in a read-only editor.
type QueryValidator interface {
	// ValidateQuery reports whether every statement of the script is read-only
	// and whether every statement returns rows. A script the engine rejects
	// returns a *SyntaxError.
	ValidateQuery(script string) (readOnly bool, returnsData bool, err error)
}

// Explainer writes the EXPLAIN statement that plans a statement.
type Explainer interface {
	// Explain returns the statement that plans statement in format, falling
	// back to the engine's own format for one it does not offer. It returns
	// ErrExplainAnalyzeExecutes for a statement that would execute what it
	// plans, and a *SyntaxError for one the engine rejects.
	Explain(statement string, format ExplainFormat) (string, error)
}

// ResultLimiter caps the rows a query returns.
type ResultLimiter interface {
	// ResultLimit returns statement rewritten to return at most limit rows.
	// engineVersion is the server version the caller last recorded, or empty.
	ResultLimit(statement string, limit int, engineVersion string) string
}
