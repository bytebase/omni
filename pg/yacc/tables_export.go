package yacc

// tables_export.go exports the goyacc-generated parse tables for use
// by the completion engine. These are thin wrappers around the unexported
// arrays in parser.go.

// Pact returns the parser action base table (indexed by state).
func Pact() []int32 { return pgPact[:] }

// Act returns the action table.
func Act() []int16 { return pgAct[:] }

// Chk returns the check table.
func Chk() []int16 { return pgChk[:] }

// Def returns the default reduction table (indexed by state).
func Def() []int16 { return pgDef[:] }

// R1 returns the LHS nonterminal for each production.
func R1() []int16 { return pgR1[:] }

// R2 returns the RHS length for each production.
func R2() []int8 { return pgR2[:] }

// Pgo returns the goto table base offsets.
func Pgo() []int16 { return pgPgo[:] }

// Exca returns the exception table.
func Exca() []int16 { return pgExca[:] }

// TokNames returns the token name table.
func TokNames() []string { return pgToknames[:] }

// Tok1 returns the single-character token mapping table.
func Tok1() []int16 { return pgTok1[:] }

// Tok2 returns the double-character token mapping table.
func Tok2() []int16 { return pgTok2[:] }

// Tok3 returns the extended token mapping table.
func Tok3() []uint16 { return pgTok3[:] }

// NTokens returns the number of tokens (length of pgToknames).
func NTokens() int { return len(pgToknames) }

// Last returns pgLast, the size of the action table.
func Last() int { return pgLast }

// Flag returns pgFlag, the sentinel value for terminal actions.
func Flag() int { return pgFlag }

// EofCode returns the EOF token code.
func EofCode() int { return pgEofCode }

// ErrCode returns the error token code.
func ErrCode() int { return pgErrCode }

// Private returns pgPrivate, the base for private token IDs.
func Private() int { return pgPrivate }
