package validate

// Severity classifies a diagnostic.
type Severity int

const (
	SeverityError Severity = iota
	SeverityWarning
)

// Diagnostic is one semantic finding produced by Validate.
type Diagnostic struct {
	Code     string   // stable machine-readable code, e.g. "undeclared_variable"
	Message  string   // human-readable message
	Severity Severity // error or warning
	Position int      // byte offset within the original SQL text
}
