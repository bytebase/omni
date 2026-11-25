package analysis

// DiagnoseError is the first syntax error Diagnose finds.
type DiagnoseError struct {
	// Offset is the error's byte offset in the script, the start of the
	// failing statement when CockroachDB's parser did not locate it.
	Offset  int
	Message string
}

// Diagnose parses a script for an editor with CockroachDB's grammar and
// returns its first syntax error, or nil.
func Diagnose(script string) *DiagnoseError {
	_, info, err := ParseRaw(script)
	if err == nil {
		return nil
	}
	return &DiagnoseError{Offset: info.Position, Message: info.Message}
}
