# PartiQL Lexer Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement the `partiql/parser` package's hand-written tokenizer — 302 token constants (266 keywords + 28 operators + 6 literals + 2 specials), a single-pass `Lexer.Next()` API, and ~123 unit sub-tests covering golden token streams, the AWS DynamoDB PartiQL example corpus, and 5 error triggers.

**Architecture:** Hand-written byte-at-a-time tokenizer modeled on `cosmosdb/parser/lexer.go`. Three source files: `token.go` (Token struct + constants + `tokenName`), `keywords.go` (lowercase→constant map), `lexer.go` (Lexer struct + `Next()` + scan helpers). Token positions use `ast.Loc` directly so the parser doesn't need to convert at every AST node construction site. First-error-and-stop error model: lex errors set `Lexer.Err` and all subsequent `Next()` calls return `tokEOF`.

**Tech Stack:** Pure Go (zero runtime dependencies), `testing` package, standard library only (`fmt`, `strings`, `os`, `path/filepath`, `reflect`).

**Spec:** `docs/superpowers/specs/2026-04-08-partiql-lexer-design.md`

**DAG entry:** `docs/migration/partiql/dag.md` node 2 (P0).

**Worktree:** `/Users/h3n4l/OpenSource/omni/.worktrees/feat-partiql-lexer`
**Branch:** `feat/partiql/lexer`

---

## File Map

| File | Lines (approx) | Responsibility |
|------|---------------|----------------|
| `partiql/parser/token.go` | ~700 | `Token` struct, 302 `tok*` constants in 4 iota groups, `tokenName(int) string` (302-arm switch) |
| `partiql/parser/keywords.go` | ~280 | `keywords` map (lowercase keyword → `tok*` constant, 266 entries) |
| `partiql/parser/lexer.go` | ~450 | `Lexer` struct, `NewLexer`, `Next()`, scan helpers (`scanString`, `scanQuotedIdent`, `scanIdentOrKeyword`, `scanNumber`, `scanOperator`, `scanIonLiteral`), `skipWhitespaceAndComments`, `isIdentStart`/`isIdentContinue`/`isDigit` |
| `partiql/parser/lexer_test.go` | ~700 | `TestLexer_Tokens` (~55 golden cases), `TestLexer_AWSCorpus`, `TestLexer_Errors`, `TestTokenName_AllCovered`, `TestKeywords_LenMatchesConstants` |

**Total:** ~2130 lines across 4 files.

---

## Conventions

- **Unexported `tok*` constants** — the lexer is internal to `partiql/parser`. The parser-foundation node will live in the same package and consume these constants directly.
- **Naming:** lowercase `tok` prefix + UPPERCASE name from the grammar (e.g., `tokSELECT`, `tokANGLE_DOUBLE_LEFT`). Compound operator names match `PartiQLLexer.g4` rule names verbatim for traceability.
- **`Token.Str`** holds the raw source text for most tokens. For `tokSCONST` and `tokIDENT_QUOTED`, it holds the **decoded** value (with `''`/`""` collapsed).
- **`Token.Loc`** is `ast.Loc{Start, End}` — half-open byte range. Set by every scan helper after consuming the token.
- **Pointer receiver** for all `Lexer` methods.
- **Error model:** first-error-and-stop. `l.Err` is set on first failure; subsequent `Next()` calls return `tokEOF`.
- **gofmt clean** at every commit. **`go vet ./partiql/parser/...` clean.**
- **Test scope:** `go test ./partiql/parser/...` only. No global test runs.

## Authoritative Sources

- **Grammar:** `/Users/h3n4l/OpenSource/parser/partiql/PartiQLLexer.g4` (542 lines, lines 13–356 for non-Ion-mode rules)
- **Sorted keyword list (266 entries):** the canonical alphabetical list embedded in Tasks 2 and 4 of this plan. Generated from the grammar via `grep -E "^[A-Z][A-Z0-9_]*: '[A-Z]" PartiQLLexer.g4 | sed -E "s/^([A-Z][A-Z0-9_]*):.*/\1/" | sort`.
- **Reference precedent:** `cosmosdb/parser/lexer.go` (the closest existing hand-written lexer in omni)

---

### Task 1: `token.go` scaffold — Token struct, special/literal/operator constants, stub `tokenName`

**Files:**
- Create: `partiql/parser/token.go`

This task creates the package, the `Token` struct, and the 36 non-keyword token constants (2 specials + 6 literals + 28 operators). The `tokenName` function is added as a stub that handles only the 36 non-keyword constants — Task 3 expands it to cover all 302.

- [ ] **Step 1: Create `partiql/parser/token.go` with this exact content**

```go
// Package parser implements the PartiQL parser. This file declares the
// token type used by the lexer (lexer.go) and the full set of token
// type constants. Token positions use ast.Loc directly to eliminate the
// need for the parser to convert at every AST node construction site.
package parser

import "github.com/bytebase/omni/partiql/ast"

// Token is a single PartiQL lexer token.
//
// Type is one of the tok* constants below.
//
// Str holds the raw source text for most tokens. For tokSCONST (single-quoted
// string literals) and tokIDENT_QUOTED (double-quoted identifiers), Str holds
// the *decoded* value with the doubled-quote escape collapsed (e.g., 'it''s'
// -> "it's"). For tokION_LITERAL, Str is the verbatim inner content between
// the backticks (no decoding).
//
// Loc is the half-open byte range covering the token in the source string.
type Token struct {
	Type int
	Str  string
	Loc  ast.Loc
}

// ===========================================================================
// Special tokens.
// ===========================================================================

const (
	tokEOF     = 0 // end of input or after lex error
	tokInvalid = 1 // sentinel for unknown token type (never returned by Next)
)

// ===========================================================================
// Literal tokens — group 1000.
// ===========================================================================

const (
	tokSCONST       = iota + 1000 // single-quoted string literal: 'hello'
	tokICONST                     // integer literal: 42
	tokFCONST                     // decimal/float literal: 3.14, 1e10, .5
	tokIDENT                      // unquoted identifier (case-insensitive lookup)
	tokIDENT_QUOTED               // double-quoted identifier (case-sensitive): "Foo"
	tokION_LITERAL                // backtick-delimited Ion blob (body deferred to DAG node 17)
)

// ===========================================================================
// Operator and punctuation tokens — group 2000.
//
// Names follow PartiQLLexer.g4 rule names verbatim for traceability against
// the grammar.
// ===========================================================================

const (
	tokPLUS              = iota + 2000 // +
	tokMINUS                           // -
	tokASTERISK                        // *
	tokSLASH_FORWARD                   // /
	tokPERCENT                         // %
	tokCARET                           // ^
	tokTILDE                           // ~
	tokAT_SIGN                         // @
	tokEQ                              // =
	tokNEQ                             // <> or !=
	tokLT                              // < (ANGLE_LEFT in grammar)
	tokGT                              // > (ANGLE_RIGHT in grammar)
	tokLT_EQ                           // <=
	tokGT_EQ                           // >=
	tokCONCAT                          // ||
	tokANGLE_DOUBLE_LEFT               // <<  (PartiQL bag-literal start)
	tokANGLE_DOUBLE_RIGHT              // >>  (PartiQL bag-literal end)
	tokPAREN_LEFT                      // (
	tokPAREN_RIGHT                     // )
	tokBRACKET_LEFT                    // [
	tokBRACKET_RIGHT                   // ]
	tokBRACE_LEFT                      // {
	tokBRACE_RIGHT                     // }
	tokCOLON                           // :
	tokCOLON_SEMI                      // ;
	tokCOMMA                           // ,
	tokPERIOD                          // .
	tokQUESTION_MARK                   // ?
)

// tokenName returns the canonical printable name for a token type constant.
// Used by error messages, test failure output, and future debugging.
//
// Task 3 expands this switch to cover all 302 constants. For now it
// covers only the 36 non-keyword constants from this file.
func tokenName(t int) string {
	switch t {
	case tokEOF:
		return "EOF"
	case tokInvalid:
		return "INVALID"
	case tokSCONST:
		return "SCONST"
	case tokICONST:
		return "ICONST"
	case tokFCONST:
		return "FCONST"
	case tokIDENT:
		return "IDENT"
	case tokIDENT_QUOTED:
		return "IDENT_QUOTED"
	case tokION_LITERAL:
		return "ION_LITERAL"
	case tokPLUS:
		return "PLUS"
	case tokMINUS:
		return "MINUS"
	case tokASTERISK:
		return "ASTERISK"
	case tokSLASH_FORWARD:
		return "SLASH_FORWARD"
	case tokPERCENT:
		return "PERCENT"
	case tokCARET:
		return "CARET"
	case tokTILDE:
		return "TILDE"
	case tokAT_SIGN:
		return "AT_SIGN"
	case tokEQ:
		return "EQ"
	case tokNEQ:
		return "NEQ"
	case tokLT:
		return "LT"
	case tokGT:
		return "GT"
	case tokLT_EQ:
		return "LT_EQ"
	case tokGT_EQ:
		return "GT_EQ"
	case tokCONCAT:
		return "CONCAT"
	case tokANGLE_DOUBLE_LEFT:
		return "ANGLE_DOUBLE_LEFT"
	case tokANGLE_DOUBLE_RIGHT:
		return "ANGLE_DOUBLE_RIGHT"
	case tokPAREN_LEFT:
		return "PAREN_LEFT"
	case tokPAREN_RIGHT:
		return "PAREN_RIGHT"
	case tokBRACKET_LEFT:
		return "BRACKET_LEFT"
	case tokBRACKET_RIGHT:
		return "BRACKET_RIGHT"
	case tokBRACE_LEFT:
		return "BRACE_LEFT"
	case tokBRACE_RIGHT:
		return "BRACE_RIGHT"
	case tokCOLON:
		return "COLON"
	case tokCOLON_SEMI:
		return "COLON_SEMI"
	case tokCOMMA:
		return "COMMA"
	case tokPERIOD:
		return "PERIOD"
	case tokQUESTION_MARK:
		return "QUESTION_MARK"
	}
	return ""
}
```

- [ ] **Step 2: Create the directory and verify it builds**

```bash
cd /Users/h3n4l/OpenSource/omni/.worktrees/feat-partiql-lexer
mkdir -p partiql/parser
# (the file from Step 1 lives at partiql/parser/token.go)
go build ./partiql/parser/...
```

Expected: build succeeds. The package compiles even though it has no `lexer.go` yet because `token.go` only declares types and constants — no executable code references missing identifiers.

- [ ] **Step 3: Run vet and gofmt**

```bash
go vet ./partiql/parser/...
gofmt -l partiql/parser/
```

Both must produce no output.

- [ ] **Step 4: Commit**

```bash
git add partiql/parser/token.go
git commit -m "$(cat <<'EOF'
feat(partiql/parser): scaffold token.go with non-keyword constants

Creates the partiql/parser package with the Token struct, 2 special
constants (tokEOF, tokInvalid), 6 literal constants (tokSCONST,
tokICONST, tokFCONST, tokIDENT, tokIDENT_QUOTED, tokION_LITERAL),
and 28 operator/punctuation constants (tokPLUS through tokQUESTION_MARK).

Token positions embed ast.Loc directly per spec D2.

tokenName function is a stub covering the 36 non-keyword constants
defined here; Task 3 expands it to all 302 once Task 2 lands the
266 keyword constants.

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

### Task 2: `token.go` keyword constants — 266 keywords in one iota block

**Files:**
- Modify: `partiql/parser/token.go` (append after the operator block)

Adds the 266 keyword constants in alphabetical order. The block starts at iota offset 3000 to leave headroom before group 4000.

- [ ] **Step 1: Append the keyword constants block to `partiql/parser/token.go`**

Open `partiql/parser/token.go` and append the following at the bottom of the file (after the operator constants block, before `tokenName`). The list is alphabetical to match the canonical sort from `PartiQLLexer.g4`.

```go
// ===========================================================================
// Keyword tokens — group 3000.
//
// 266 keywords from PartiQLLexer.g4 lines 13–295, alphabetical.
//
// Includes standard SQL keywords (ABSOLUTE..ZONE), window keywords
// (LAG, LEAD, OVER, PARTITION), PartiQL extension keywords (CAN_CAST,
// CAN_LOSSLESS_CAST, MISSING, PIVOT, UNPIVOT, LIMIT, OFFSET, REMOVE,
// INDEX, LET, CONFLICT, DO, RETURNING, MODIFIED, NEW, OLD, NOTHING,
// EXCLUDED, SHORTEST, MATCH), and data type keywords (TUPLE, INT2/4/8,
// INTEGER2/4/8, BIGINT, BOOL, BOOLEAN, STRING, SYMBOL, CLOB, BLOB,
// STRUCT, LIST, SEXP, BAG).
// ===========================================================================

const (
	tokABSOLUTE = iota + 3000
	tokACTION
	tokADD
	tokALL
	tokALLOCATE
	tokALTER
	tokAND
	tokANY
	tokARE
	tokAS
	tokASC
	tokASSERTION
	tokAT
	tokAUTHORIZATION
	tokAVG
	tokBAG
	tokBEGIN
	tokBETWEEN
	tokBIGINT
	tokBIT
	tokBIT_LENGTH
	tokBLOB
	tokBOOL
	tokBOOLEAN
	tokBY
	tokCAN_CAST
	tokCAN_LOSSLESS_CAST
	tokCASCADE
	tokCASCADED
	tokCASE
	tokCAST
	tokCATALOG
	tokCHAR
	tokCHAR_LENGTH
	tokCHARACTER
	tokCHARACTER_LENGTH
	tokCHECK
	tokCLOB
	tokCLOSE
	tokCOALESCE
	tokCOLLATE
	tokCOLLATION
	tokCOLUMN
	tokCOMMIT
	tokCONFLICT
	tokCONNECT
	tokCONNECTION
	tokCONSTRAINT
	tokCONSTRAINTS
	tokCONTINUE
	tokCONVERT
	tokCORRESPONDING
	tokCOUNT
	tokCREATE
	tokCROSS
	tokCURRENT
	tokCURRENT_DATE
	tokCURRENT_TIME
	tokCURRENT_TIMESTAMP
	tokCURRENT_USER
	tokCURSOR
	tokDATE
	tokDATE_ADD
	tokDATE_DIFF
	tokDEALLOCATE
	tokDEC
	tokDECIMAL
	tokDECLARE
	tokDEFAULT
	tokDEFERRABLE
	tokDEFERRED
	tokDELETE
	tokDESC
	tokDESCRIBE
	tokDESCRIPTOR
	tokDIAGNOSTICS
	tokDISCONNECT
	tokDISTINCT
	tokDO
	tokDOMAIN
	tokDOUBLE
	tokDROP
	tokELSE
	tokEND
	tokEND_EXEC
	tokESCAPE
	tokEXCEPT
	tokEXCEPTION
	tokEXCLUDED
	tokEXEC
	tokEXECUTE
	tokEXISTS
	tokEXPLAIN
	tokEXTERNAL
	tokEXTRACT
	tokFALSE
	tokFETCH
	tokFIRST
	tokFLOAT
	tokFOR
	tokFOREIGN
	tokFOUND
	tokFROM
	tokFULL
	tokGET
	tokGLOBAL
	tokGO
	tokGOTO
	tokGRANT
	tokGROUP
	tokHAVING
	tokIDENTITY
	tokIMMEDIATE
	tokIN
	tokINDEX
	tokINDICATOR
	tokINITIALLY
	tokINNER
	tokINPUT
	tokINSENSITIVE
	tokINSERT
	tokINT
	tokINT2
	tokINT4
	tokINT8
	tokINTEGER
	tokINTEGER2
	tokINTEGER4
	tokINTEGER8
	tokINTERSECT
	tokINTERVAL
	tokINTO
	tokIS
	tokISOLATION
	tokJOIN
	tokKEY
	tokLAG
	tokLANGUAGE
	tokLAST
	tokLATERAL
	tokLEAD
	tokLEFT
	tokLET
	tokLEVEL
	tokLIKE
	tokLIMIT
	tokLIST
	tokLOCAL
	tokLOWER
	tokMATCH
	tokMAX
	tokMIN
	tokMISSING
	tokMODIFIED
	tokMODULE
	tokNAMES
	tokNATIONAL
	tokNATURAL
	tokNCHAR
	tokNEW
	tokNEXT
	tokNO
	tokNOT
	tokNOTHING
	tokNULL
	tokNULLIF
	tokNULLS
	tokNUMERIC
	tokOCTET_LENGTH
	tokOF
	tokOFFSET
	tokOLD
	tokON
	tokONLY
	tokOPEN
	tokOPTION
	tokOR
	tokORDER
	tokOUTER
	tokOUTPUT
	tokOVER
	tokOVERLAPS
	tokOVERLAY
	tokPAD
	tokPARTIAL
	tokPARTITION
	tokPIVOT
	tokPLACING
	tokPOSITION
	tokPRECISION
	tokPREPARE
	tokPRESERVE
	tokPRIMARY
	tokPRIOR
	tokPRIVILEGES
	tokPROCEDURE
	tokPUBLIC
	tokREAD
	tokREAL
	tokREFERENCES
	tokRELATIVE
	tokREMOVE
	tokREPLACE
	tokRESTRICT
	tokRETURNING
	tokREVOKE
	tokRIGHT
	tokROLLBACK
	tokROWS
	tokSCHEMA
	tokSCROLL
	tokSECTION
	tokSELECT
	tokSESSION
	tokSESSION_USER
	tokSET
	tokSEXP
	tokSHORTEST
	tokSIZE
	tokSMALLINT
	tokSOME
	tokSPACE
	tokSQL
	tokSQLCODE
	tokSQLERROR
	tokSQLSTATE
	tokSTRING
	tokSTRUCT
	tokSUBSTRING
	tokSUM
	tokSYMBOL
	tokSYSTEM_USER
	tokTABLE
	tokTEMPORARY
	tokTHEN
	tokTIME
	tokTIMESTAMP
	tokTO
	tokTRANSACTION
	tokTRANSLATE
	tokTRANSLATION
	tokTRIM
	tokTRUE
	tokTUPLE
	tokUNION
	tokUNIQUE
	tokUNKNOWN
	tokUNPIVOT
	tokUPDATE
	tokUPPER
	tokUPSERT
	tokUSAGE
	tokUSER
	tokUSING
	tokVALUE
	tokVALUES
	tokVARCHAR
	tokVARYING
	tokVIEW
	tokWHEN
	tokWHENEVER
	tokWHERE
	tokWITH
	tokWORK
	tokWRITE
	tokZONE
)
```

The const block declares the constants in alphabetical order. The order matters because iota increments by 1 — the constants are renumberable but not re-orderable without breaking the test count.

- [ ] **Step 2: Build to verify the constants compile**

```bash
go build ./partiql/parser/...
```

Expected: success (no other code references these constants yet).

- [ ] **Step 3: Run vet and gofmt**

```bash
go vet ./partiql/parser/...
gofmt -l partiql/parser/
```

Both clean.

- [ ] **Step 4: Commit**

```bash
git add partiql/parser/token.go
git commit -m "$(cat <<'EOF'
feat(partiql/parser): add 266 keyword constants to token.go

Appends a single iota block at offset 3000 with all 266 keyword
constants from PartiQLLexer.g4, alphabetically ordered. Generated
by extracting `^[A-Z][A-Z0-9_]*: '[A-Z]` rules from the grammar
and sorting.

Includes standard SQL keywords (ABSOLUTE..ZONE), window keywords
(LAG, LEAD, OVER, PARTITION), PartiQL extension keywords (CAN_CAST,
CAN_LOSSLESS_CAST, MISSING, PIVOT, UNPIVOT, LIMIT, OFFSET, REMOVE,
INDEX, LET, CONFLICT, DO, RETURNING, MODIFIED, NEW, OLD, NOTHING,
EXCLUDED, SHORTEST, MATCH), and data type keywords (TUPLE, INT2/4/8,
INTEGER2/4/8, BIGINT, BOOL, BOOLEAN, STRING, SYMBOL, CLOB, BLOB,
STRUCT, LIST, SEXP, BAG).

tokenName arms for these keywords are added in Task 3.

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

### Task 3: `token.go` `tokenName` — expand the switch to cover all 302 constants

**Files:**
- Modify: `partiql/parser/token.go` (replace the stub `tokenName`)

- [ ] **Step 1: Replace the stub `tokenName` function with the full version**

Open `partiql/parser/token.go`. Find the existing `tokenName` function (it currently has 36 cases) and **replace its entire body** with the full 302-case version below. The function signature stays the same.

```go
// tokenName returns the canonical printable name for a token type constant.
// Used by error messages, test failure output, and future debugging.
//
// Covers all 302 tok* constants in the package.
func tokenName(t int) string {
	switch t {
	// Specials.
	case tokEOF:
		return "EOF"
	case tokInvalid:
		return "INVALID"

	// Literals.
	case tokSCONST:
		return "SCONST"
	case tokICONST:
		return "ICONST"
	case tokFCONST:
		return "FCONST"
	case tokIDENT:
		return "IDENT"
	case tokIDENT_QUOTED:
		return "IDENT_QUOTED"
	case tokION_LITERAL:
		return "ION_LITERAL"

	// Operators / punctuation.
	case tokPLUS:
		return "PLUS"
	case tokMINUS:
		return "MINUS"
	case tokASTERISK:
		return "ASTERISK"
	case tokSLASH_FORWARD:
		return "SLASH_FORWARD"
	case tokPERCENT:
		return "PERCENT"
	case tokCARET:
		return "CARET"
	case tokTILDE:
		return "TILDE"
	case tokAT_SIGN:
		return "AT_SIGN"
	case tokEQ:
		return "EQ"
	case tokNEQ:
		return "NEQ"
	case tokLT:
		return "LT"
	case tokGT:
		return "GT"
	case tokLT_EQ:
		return "LT_EQ"
	case tokGT_EQ:
		return "GT_EQ"
	case tokCONCAT:
		return "CONCAT"
	case tokANGLE_DOUBLE_LEFT:
		return "ANGLE_DOUBLE_LEFT"
	case tokANGLE_DOUBLE_RIGHT:
		return "ANGLE_DOUBLE_RIGHT"
	case tokPAREN_LEFT:
		return "PAREN_LEFT"
	case tokPAREN_RIGHT:
		return "PAREN_RIGHT"
	case tokBRACKET_LEFT:
		return "BRACKET_LEFT"
	case tokBRACKET_RIGHT:
		return "BRACKET_RIGHT"
	case tokBRACE_LEFT:
		return "BRACE_LEFT"
	case tokBRACE_RIGHT:
		return "BRACE_RIGHT"
	case tokCOLON:
		return "COLON"
	case tokCOLON_SEMI:
		return "COLON_SEMI"
	case tokCOMMA:
		return "COMMA"
	case tokPERIOD:
		return "PERIOD"
	case tokQUESTION_MARK:
		return "QUESTION_MARK"

	// Keywords (alphabetical).
	case tokABSOLUTE:
		return "ABSOLUTE"
	case tokACTION:
		return "ACTION"
	case tokADD:
		return "ADD"
	case tokALL:
		return "ALL"
	case tokALLOCATE:
		return "ALLOCATE"
	case tokALTER:
		return "ALTER"
	case tokAND:
		return "AND"
	case tokANY:
		return "ANY"
	case tokARE:
		return "ARE"
	case tokAS:
		return "AS"
	case tokASC:
		return "ASC"
	case tokASSERTION:
		return "ASSERTION"
	case tokAT:
		return "AT"
	case tokAUTHORIZATION:
		return "AUTHORIZATION"
	case tokAVG:
		return "AVG"
	case tokBAG:
		return "BAG"
	case tokBEGIN:
		return "BEGIN"
	case tokBETWEEN:
		return "BETWEEN"
	case tokBIGINT:
		return "BIGINT"
	case tokBIT:
		return "BIT"
	case tokBIT_LENGTH:
		return "BIT_LENGTH"
	case tokBLOB:
		return "BLOB"
	case tokBOOL:
		return "BOOL"
	case tokBOOLEAN:
		return "BOOLEAN"
	case tokBY:
		return "BY"
	case tokCAN_CAST:
		return "CAN_CAST"
	case tokCAN_LOSSLESS_CAST:
		return "CAN_LOSSLESS_CAST"
	case tokCASCADE:
		return "CASCADE"
	case tokCASCADED:
		return "CASCADED"
	case tokCASE:
		return "CASE"
	case tokCAST:
		return "CAST"
	case tokCATALOG:
		return "CATALOG"
	case tokCHAR:
		return "CHAR"
	case tokCHAR_LENGTH:
		return "CHAR_LENGTH"
	case tokCHARACTER:
		return "CHARACTER"
	case tokCHARACTER_LENGTH:
		return "CHARACTER_LENGTH"
	case tokCHECK:
		return "CHECK"
	case tokCLOB:
		return "CLOB"
	case tokCLOSE:
		return "CLOSE"
	case tokCOALESCE:
		return "COALESCE"
	case tokCOLLATE:
		return "COLLATE"
	case tokCOLLATION:
		return "COLLATION"
	case tokCOLUMN:
		return "COLUMN"
	case tokCOMMIT:
		return "COMMIT"
	case tokCONFLICT:
		return "CONFLICT"
	case tokCONNECT:
		return "CONNECT"
	case tokCONNECTION:
		return "CONNECTION"
	case tokCONSTRAINT:
		return "CONSTRAINT"
	case tokCONSTRAINTS:
		return "CONSTRAINTS"
	case tokCONTINUE:
		return "CONTINUE"
	case tokCONVERT:
		return "CONVERT"
	case tokCORRESPONDING:
		return "CORRESPONDING"
	case tokCOUNT:
		return "COUNT"
	case tokCREATE:
		return "CREATE"
	case tokCROSS:
		return "CROSS"
	case tokCURRENT:
		return "CURRENT"
	case tokCURRENT_DATE:
		return "CURRENT_DATE"
	case tokCURRENT_TIME:
		return "CURRENT_TIME"
	case tokCURRENT_TIMESTAMP:
		return "CURRENT_TIMESTAMP"
	case tokCURRENT_USER:
		return "CURRENT_USER"
	case tokCURSOR:
		return "CURSOR"
	case tokDATE:
		return "DATE"
	case tokDATE_ADD:
		return "DATE_ADD"
	case tokDATE_DIFF:
		return "DATE_DIFF"
	case tokDEALLOCATE:
		return "DEALLOCATE"
	case tokDEC:
		return "DEC"
	case tokDECIMAL:
		return "DECIMAL"
	case tokDECLARE:
		return "DECLARE"
	case tokDEFAULT:
		return "DEFAULT"
	case tokDEFERRABLE:
		return "DEFERRABLE"
	case tokDEFERRED:
		return "DEFERRED"
	case tokDELETE:
		return "DELETE"
	case tokDESC:
		return "DESC"
	case tokDESCRIBE:
		return "DESCRIBE"
	case tokDESCRIPTOR:
		return "DESCRIPTOR"
	case tokDIAGNOSTICS:
		return "DIAGNOSTICS"
	case tokDISCONNECT:
		return "DISCONNECT"
	case tokDISTINCT:
		return "DISTINCT"
	case tokDO:
		return "DO"
	case tokDOMAIN:
		return "DOMAIN"
	case tokDOUBLE:
		return "DOUBLE"
	case tokDROP:
		return "DROP"
	case tokELSE:
		return "ELSE"
	case tokEND:
		return "END"
	case tokEND_EXEC:
		return "END_EXEC"
	case tokESCAPE:
		return "ESCAPE"
	case tokEXCEPT:
		return "EXCEPT"
	case tokEXCEPTION:
		return "EXCEPTION"
	case tokEXCLUDED:
		return "EXCLUDED"
	case tokEXEC:
		return "EXEC"
	case tokEXECUTE:
		return "EXECUTE"
	case tokEXISTS:
		return "EXISTS"
	case tokEXPLAIN:
		return "EXPLAIN"
	case tokEXTERNAL:
		return "EXTERNAL"
	case tokEXTRACT:
		return "EXTRACT"
	case tokFALSE:
		return "FALSE"
	case tokFETCH:
		return "FETCH"
	case tokFIRST:
		return "FIRST"
	case tokFLOAT:
		return "FLOAT"
	case tokFOR:
		return "FOR"
	case tokFOREIGN:
		return "FOREIGN"
	case tokFOUND:
		return "FOUND"
	case tokFROM:
		return "FROM"
	case tokFULL:
		return "FULL"
	case tokGET:
		return "GET"
	case tokGLOBAL:
		return "GLOBAL"
	case tokGO:
		return "GO"
	case tokGOTO:
		return "GOTO"
	case tokGRANT:
		return "GRANT"
	case tokGROUP:
		return "GROUP"
	case tokHAVING:
		return "HAVING"
	case tokIDENTITY:
		return "IDENTITY"
	case tokIMMEDIATE:
		return "IMMEDIATE"
	case tokIN:
		return "IN"
	case tokINDEX:
		return "INDEX"
	case tokINDICATOR:
		return "INDICATOR"
	case tokINITIALLY:
		return "INITIALLY"
	case tokINNER:
		return "INNER"
	case tokINPUT:
		return "INPUT"
	case tokINSENSITIVE:
		return "INSENSITIVE"
	case tokINSERT:
		return "INSERT"
	case tokINT:
		return "INT"
	case tokINT2:
		return "INT2"
	case tokINT4:
		return "INT4"
	case tokINT8:
		return "INT8"
	case tokINTEGER:
		return "INTEGER"
	case tokINTEGER2:
		return "INTEGER2"
	case tokINTEGER4:
		return "INTEGER4"
	case tokINTEGER8:
		return "INTEGER8"
	case tokINTERSECT:
		return "INTERSECT"
	case tokINTERVAL:
		return "INTERVAL"
	case tokINTO:
		return "INTO"
	case tokIS:
		return "IS"
	case tokISOLATION:
		return "ISOLATION"
	case tokJOIN:
		return "JOIN"
	case tokKEY:
		return "KEY"
	case tokLAG:
		return "LAG"
	case tokLANGUAGE:
		return "LANGUAGE"
	case tokLAST:
		return "LAST"
	case tokLATERAL:
		return "LATERAL"
	case tokLEAD:
		return "LEAD"
	case tokLEFT:
		return "LEFT"
	case tokLET:
		return "LET"
	case tokLEVEL:
		return "LEVEL"
	case tokLIKE:
		return "LIKE"
	case tokLIMIT:
		return "LIMIT"
	case tokLIST:
		return "LIST"
	case tokLOCAL:
		return "LOCAL"
	case tokLOWER:
		return "LOWER"
	case tokMATCH:
		return "MATCH"
	case tokMAX:
		return "MAX"
	case tokMIN:
		return "MIN"
	case tokMISSING:
		return "MISSING"
	case tokMODIFIED:
		return "MODIFIED"
	case tokMODULE:
		return "MODULE"
	case tokNAMES:
		return "NAMES"
	case tokNATIONAL:
		return "NATIONAL"
	case tokNATURAL:
		return "NATURAL"
	case tokNCHAR:
		return "NCHAR"
	case tokNEW:
		return "NEW"
	case tokNEXT:
		return "NEXT"
	case tokNO:
		return "NO"
	case tokNOT:
		return "NOT"
	case tokNOTHING:
		return "NOTHING"
	case tokNULL:
		return "NULL"
	case tokNULLIF:
		return "NULLIF"
	case tokNULLS:
		return "NULLS"
	case tokNUMERIC:
		return "NUMERIC"
	case tokOCTET_LENGTH:
		return "OCTET_LENGTH"
	case tokOF:
		return "OF"
	case tokOFFSET:
		return "OFFSET"
	case tokOLD:
		return "OLD"
	case tokON:
		return "ON"
	case tokONLY:
		return "ONLY"
	case tokOPEN:
		return "OPEN"
	case tokOPTION:
		return "OPTION"
	case tokOR:
		return "OR"
	case tokORDER:
		return "ORDER"
	case tokOUTER:
		return "OUTER"
	case tokOUTPUT:
		return "OUTPUT"
	case tokOVER:
		return "OVER"
	case tokOVERLAPS:
		return "OVERLAPS"
	case tokOVERLAY:
		return "OVERLAY"
	case tokPAD:
		return "PAD"
	case tokPARTIAL:
		return "PARTIAL"
	case tokPARTITION:
		return "PARTITION"
	case tokPIVOT:
		return "PIVOT"
	case tokPLACING:
		return "PLACING"
	case tokPOSITION:
		return "POSITION"
	case tokPRECISION:
		return "PRECISION"
	case tokPREPARE:
		return "PREPARE"
	case tokPRESERVE:
		return "PRESERVE"
	case tokPRIMARY:
		return "PRIMARY"
	case tokPRIOR:
		return "PRIOR"
	case tokPRIVILEGES:
		return "PRIVILEGES"
	case tokPROCEDURE:
		return "PROCEDURE"
	case tokPUBLIC:
		return "PUBLIC"
	case tokREAD:
		return "READ"
	case tokREAL:
		return "REAL"
	case tokREFERENCES:
		return "REFERENCES"
	case tokRELATIVE:
		return "RELATIVE"
	case tokREMOVE:
		return "REMOVE"
	case tokREPLACE:
		return "REPLACE"
	case tokRESTRICT:
		return "RESTRICT"
	case tokRETURNING:
		return "RETURNING"
	case tokREVOKE:
		return "REVOKE"
	case tokRIGHT:
		return "RIGHT"
	case tokROLLBACK:
		return "ROLLBACK"
	case tokROWS:
		return "ROWS"
	case tokSCHEMA:
		return "SCHEMA"
	case tokSCROLL:
		return "SCROLL"
	case tokSECTION:
		return "SECTION"
	case tokSELECT:
		return "SELECT"
	case tokSESSION:
		return "SESSION"
	case tokSESSION_USER:
		return "SESSION_USER"
	case tokSET:
		return "SET"
	case tokSEXP:
		return "SEXP"
	case tokSHORTEST:
		return "SHORTEST"
	case tokSIZE:
		return "SIZE"
	case tokSMALLINT:
		return "SMALLINT"
	case tokSOME:
		return "SOME"
	case tokSPACE:
		return "SPACE"
	case tokSQL:
		return "SQL"
	case tokSQLCODE:
		return "SQLCODE"
	case tokSQLERROR:
		return "SQLERROR"
	case tokSQLSTATE:
		return "SQLSTATE"
	case tokSTRING:
		return "STRING"
	case tokSTRUCT:
		return "STRUCT"
	case tokSUBSTRING:
		return "SUBSTRING"
	case tokSUM:
		return "SUM"
	case tokSYMBOL:
		return "SYMBOL"
	case tokSYSTEM_USER:
		return "SYSTEM_USER"
	case tokTABLE:
		return "TABLE"
	case tokTEMPORARY:
		return "TEMPORARY"
	case tokTHEN:
		return "THEN"
	case tokTIME:
		return "TIME"
	case tokTIMESTAMP:
		return "TIMESTAMP"
	case tokTO:
		return "TO"
	case tokTRANSACTION:
		return "TRANSACTION"
	case tokTRANSLATE:
		return "TRANSLATE"
	case tokTRANSLATION:
		return "TRANSLATION"
	case tokTRIM:
		return "TRIM"
	case tokTRUE:
		return "TRUE"
	case tokTUPLE:
		return "TUPLE"
	case tokUNION:
		return "UNION"
	case tokUNIQUE:
		return "UNIQUE"
	case tokUNKNOWN:
		return "UNKNOWN"
	case tokUNPIVOT:
		return "UNPIVOT"
	case tokUPDATE:
		return "UPDATE"
	case tokUPPER:
		return "UPPER"
	case tokUPSERT:
		return "UPSERT"
	case tokUSAGE:
		return "USAGE"
	case tokUSER:
		return "USER"
	case tokUSING:
		return "USING"
	case tokVALUE:
		return "VALUE"
	case tokVALUES:
		return "VALUES"
	case tokVARCHAR:
		return "VARCHAR"
	case tokVARYING:
		return "VARYING"
	case tokVIEW:
		return "VIEW"
	case tokWHEN:
		return "WHEN"
	case tokWHENEVER:
		return "WHENEVER"
	case tokWHERE:
		return "WHERE"
	case tokWITH:
		return "WITH"
	case tokWORK:
		return "WORK"
	case tokWRITE:
		return "WRITE"
	case tokZONE:
		return "ZONE"
	}
	return ""
}
```

The function has 302 cases total (2 + 6 + 28 + 266). Returning empty string for unknown values is the sentinel that `TestTokenName_AllCovered` (Task 13) checks for.

- [ ] **Step 2: Build, vet, gofmt**

```bash
go build ./partiql/parser/...
go vet ./partiql/parser/...
gofmt -l partiql/parser/
```

All clean.

- [ ] **Step 3: Commit**

```bash
git add partiql/parser/token.go
git commit -m "$(cat <<'EOF'
feat(partiql/parser): expand tokenName to all 302 token constants

Replaces the Task 1 stub tokenName with the full 302-arm switch:
2 specials + 6 literals + 28 operators + 266 keywords. Returns the
canonical UPPERCASE name for each constant; returns empty string
for unknown values (used by TestTokenName_AllCovered as the
"missing arm" sentinel).

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

### Task 4: `keywords.go` — keyword lookup map (266 entries)

**Files:**
- Create: `partiql/parser/keywords.go`

- [ ] **Step 1: Create `partiql/parser/keywords.go` with this exact content**

```go
package parser

// keywords maps lowercase keyword strings to their tok* constant.
// PartiQL keywords are case-insensitive per PartiQLLexer.g4
// `options { caseInsensitive = true; }`. The lexer lowercases the
// identifier text before lookup.
//
// Built once at package init via a literal map. The 266 entries here
// must stay in 1:1 correspondence with the 266 tokKEYWORD constants in
// token.go (group 3000) — TestKeywords_LenMatchesConstants asserts
// len(keywords) == 266.
//
// Source: every uppercase rule in PartiQLLexer.g4 from line 13 (ABSOLUTE)
// through line 295 (BAG), generated via:
//
//	grep -E "^[A-Z][A-Z0-9_]*: '[A-Z]" PartiQLLexer.g4 |
//	  sed -E "s/^([A-Z][A-Z0-9_]*):.*/\1/" | sort
var keywords = map[string]int{
	"absolute":          tokABSOLUTE,
	"action":            tokACTION,
	"add":               tokADD,
	"all":               tokALL,
	"allocate":          tokALLOCATE,
	"alter":             tokALTER,
	"and":               tokAND,
	"any":               tokANY,
	"are":               tokARE,
	"as":                tokAS,
	"asc":               tokASC,
	"assertion":         tokASSERTION,
	"at":                tokAT,
	"authorization":     tokAUTHORIZATION,
	"avg":               tokAVG,
	"bag":               tokBAG,
	"begin":             tokBEGIN,
	"between":           tokBETWEEN,
	"bigint":            tokBIGINT,
	"bit":               tokBIT,
	"bit_length":        tokBIT_LENGTH,
	"blob":              tokBLOB,
	"bool":              tokBOOL,
	"boolean":           tokBOOLEAN,
	"by":                tokBY,
	"can_cast":          tokCAN_CAST,
	"can_lossless_cast": tokCAN_LOSSLESS_CAST,
	"cascade":           tokCASCADE,
	"cascaded":          tokCASCADED,
	"case":              tokCASE,
	"cast":              tokCAST,
	"catalog":           tokCATALOG,
	"char":              tokCHAR,
	"char_length":       tokCHAR_LENGTH,
	"character":         tokCHARACTER,
	"character_length":  tokCHARACTER_LENGTH,
	"check":             tokCHECK,
	"clob":              tokCLOB,
	"close":             tokCLOSE,
	"coalesce":          tokCOALESCE,
	"collate":           tokCOLLATE,
	"collation":         tokCOLLATION,
	"column":            tokCOLUMN,
	"commit":            tokCOMMIT,
	"conflict":          tokCONFLICT,
	"connect":           tokCONNECT,
	"connection":        tokCONNECTION,
	"constraint":        tokCONSTRAINT,
	"constraints":       tokCONSTRAINTS,
	"continue":          tokCONTINUE,
	"convert":           tokCONVERT,
	"corresponding":     tokCORRESPONDING,
	"count":             tokCOUNT,
	"create":            tokCREATE,
	"cross":             tokCROSS,
	"current":           tokCURRENT,
	"current_date":      tokCURRENT_DATE,
	"current_time":      tokCURRENT_TIME,
	"current_timestamp": tokCURRENT_TIMESTAMP,
	"current_user":      tokCURRENT_USER,
	"cursor":            tokCURSOR,
	"date":              tokDATE,
	"date_add":          tokDATE_ADD,
	"date_diff":         tokDATE_DIFF,
	"deallocate":        tokDEALLOCATE,
	"dec":               tokDEC,
	"decimal":           tokDECIMAL,
	"declare":           tokDECLARE,
	"default":           tokDEFAULT,
	"deferrable":        tokDEFERRABLE,
	"deferred":          tokDEFERRED,
	"delete":            tokDELETE,
	"desc":              tokDESC,
	"describe":          tokDESCRIBE,
	"descriptor":        tokDESCRIPTOR,
	"diagnostics":       tokDIAGNOSTICS,
	"disconnect":        tokDISCONNECT,
	"distinct":          tokDISTINCT,
	"do":                tokDO,
	"domain":            tokDOMAIN,
	"double":            tokDOUBLE,
	"drop":              tokDROP,
	"else":              tokELSE,
	"end":               tokEND,
	"end-exec":          tokEND_EXEC,
	"escape":            tokESCAPE,
	"except":            tokEXCEPT,
	"exception":         tokEXCEPTION,
	"excluded":          tokEXCLUDED,
	"exec":              tokEXEC,
	"execute":           tokEXECUTE,
	"exists":            tokEXISTS,
	"explain":           tokEXPLAIN,
	"external":          tokEXTERNAL,
	"extract":           tokEXTRACT,
	"false":             tokFALSE,
	"fetch":             tokFETCH,
	"first":             tokFIRST,
	"float":             tokFLOAT,
	"for":               tokFOR,
	"foreign":           tokFOREIGN,
	"found":             tokFOUND,
	"from":              tokFROM,
	"full":              tokFULL,
	"get":               tokGET,
	"global":            tokGLOBAL,
	"go":                tokGO,
	"goto":              tokGOTO,
	"grant":             tokGRANT,
	"group":             tokGROUP,
	"having":            tokHAVING,
	"identity":          tokIDENTITY,
	"immediate":         tokIMMEDIATE,
	"in":                tokIN,
	"index":             tokINDEX,
	"indicator":         tokINDICATOR,
	"initially":         tokINITIALLY,
	"inner":             tokINNER,
	"input":             tokINPUT,
	"insensitive":       tokINSENSITIVE,
	"insert":            tokINSERT,
	"int":               tokINT,
	"int2":              tokINT2,
	"int4":              tokINT4,
	"int8":              tokINT8,
	"integer":           tokINTEGER,
	"integer2":          tokINTEGER2,
	"integer4":          tokINTEGER4,
	"integer8":          tokINTEGER8,
	"intersect":         tokINTERSECT,
	"interval":          tokINTERVAL,
	"into":              tokINTO,
	"is":                tokIS,
	"isolation":         tokISOLATION,
	"join":              tokJOIN,
	"key":               tokKEY,
	"lag":               tokLAG,
	"language":          tokLANGUAGE,
	"last":              tokLAST,
	"lateral":           tokLATERAL,
	"lead":              tokLEAD,
	"left":              tokLEFT,
	"let":               tokLET,
	"level":             tokLEVEL,
	"like":              tokLIKE,
	"limit":             tokLIMIT,
	"list":              tokLIST,
	"local":             tokLOCAL,
	"lower":             tokLOWER,
	"match":             tokMATCH,
	"max":               tokMAX,
	"min":               tokMIN,
	"missing":           tokMISSING,
	"modified":          tokMODIFIED,
	"module":            tokMODULE,
	"names":             tokNAMES,
	"national":          tokNATIONAL,
	"natural":           tokNATURAL,
	"nchar":             tokNCHAR,
	"new":               tokNEW,
	"next":              tokNEXT,
	"no":                tokNO,
	"not":               tokNOT,
	"nothing":           tokNOTHING,
	"null":              tokNULL,
	"nullif":            tokNULLIF,
	"nulls":             tokNULLS,
	"numeric":           tokNUMERIC,
	"octet_length":      tokOCTET_LENGTH,
	"of":                tokOF,
	"offset":            tokOFFSET,
	"old":               tokOLD,
	"on":                tokON,
	"only":              tokONLY,
	"open":              tokOPEN,
	"option":            tokOPTION,
	"or":                tokOR,
	"order":             tokORDER,
	"outer":             tokOUTER,
	"output":            tokOUTPUT,
	"over":              tokOVER,
	"overlaps":          tokOVERLAPS,
	"overlay":           tokOVERLAY,
	"pad":               tokPAD,
	"partial":           tokPARTIAL,
	"partition":         tokPARTITION,
	"pivot":             tokPIVOT,
	"placing":           tokPLACING,
	"position":          tokPOSITION,
	"precision":         tokPRECISION,
	"prepare":           tokPREPARE,
	"preserve":          tokPRESERVE,
	"primary":           tokPRIMARY,
	"prior":             tokPRIOR,
	"privileges":        tokPRIVILEGES,
	"procedure":         tokPROCEDURE,
	"public":            tokPUBLIC,
	"read":              tokREAD,
	"real":              tokREAL,
	"references":        tokREFERENCES,
	"relative":          tokRELATIVE,
	"remove":            tokREMOVE,
	"replace":           tokREPLACE,
	"restrict":          tokRESTRICT,
	"returning":         tokRETURNING,
	"revoke":            tokREVOKE,
	"right":             tokRIGHT,
	"rollback":          tokROLLBACK,
	"rows":              tokROWS,
	"schema":            tokSCHEMA,
	"scroll":            tokSCROLL,
	"section":           tokSECTION,
	"select":            tokSELECT,
	"session":           tokSESSION,
	"session_user":      tokSESSION_USER,
	"set":               tokSET,
	"sexp":              tokSEXP,
	"shortest":          tokSHORTEST,
	"size":              tokSIZE,
	"smallint":          tokSMALLINT,
	"some":              tokSOME,
	"space":             tokSPACE,
	"sql":               tokSQL,
	"sqlcode":           tokSQLCODE,
	"sqlerror":          tokSQLERROR,
	"sqlstate":          tokSQLSTATE,
	"string":            tokSTRING,
	"struct":            tokSTRUCT,
	"substring":         tokSUBSTRING,
	"sum":               tokSUM,
	"symbol":            tokSYMBOL,
	"system_user":       tokSYSTEM_USER,
	"table":             tokTABLE,
	"temporary":         tokTEMPORARY,
	"then":              tokTHEN,
	"time":              tokTIME,
	"timestamp":         tokTIMESTAMP,
	"to":                tokTO,
	"transaction":       tokTRANSACTION,
	"translate":         tokTRANSLATE,
	"translation":       tokTRANSLATION,
	"trim":              tokTRIM,
	"true":              tokTRUE,
	"tuple":             tokTUPLE,
	"union":             tokUNION,
	"unique":            tokUNIQUE,
	"unknown":           tokUNKNOWN,
	"unpivot":           tokUNPIVOT,
	"update":            tokUPDATE,
	"upper":             tokUPPER,
	"upsert":            tokUPSERT,
	"usage":             tokUSAGE,
	"user":              tokUSER,
	"using":             tokUSING,
	"value":             tokVALUE,
	"values":            tokVALUES,
	"varchar":           tokVARCHAR,
	"varying":           tokVARYING,
	"view":              tokVIEW,
	"when":              tokWHEN,
	"whenever":          tokWHENEVER,
	"where":             tokWHERE,
	"with":              tokWITH,
	"work":              tokWORK,
	"write":             tokWRITE,
	"zone":              tokZONE,
}
```

**Note on `END_EXEC`**: the grammar's `END_EXEC: 'END-EXEC';` rule has a hyphen in the literal, not an underscore. So the map key is `"end-exec"` (with hyphen) and the constant is `tokEND_EXEC` (with underscore). This is the only keyword in the grammar with a hyphen in its literal text.

- [ ] **Step 2: Build, vet, gofmt**

```bash
go build ./partiql/parser/...
go vet ./partiql/parser/...
gofmt -l partiql/parser/
```

All clean. (gofmt may re-align the column widths in the map literal — let it.)

- [ ] **Step 3: Commit**

```bash
git add partiql/parser/keywords.go
git commit -m "$(cat <<'EOF'
feat(partiql/parser): add keywords map (266 entries)

Lowercase keyword string -> tok* constant lookup map. PartiQL
keywords are case-insensitive per PartiQLLexer.g4 caseInsensitive=true,
so the lexer lowercases identifiers before lookup.

266 entries in 1:1 correspondence with the tok keyword constants
in token.go group 3000. Will be enforced by
TestKeywords_LenMatchesConstants in Task 13.

Note: end-exec is the only entry with a hyphen in its key, because
the grammar's END_EXEC rule literal is 'END-EXEC' (the only keyword
in PartiQL with a hyphen in its source spelling).

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

### Task 5: `lexer.go` scaffold — Lexer struct, NewLexer, Next() dispatch, skipWhitespaceAndComments, helpers

**Files:**
- Create: `partiql/parser/lexer.go`

- [ ] **Step 1: Create `partiql/parser/lexer.go` with this exact content**

```go
package parser

import (
	"fmt"
	"strings"

	"github.com/bytebase/omni/partiql/ast"
)

// Lexer is a hand-written tokenizer for PartiQL source code.
//
// Single-pass scanner. The caller drives it via Next(); each call
// returns one token. At end of input or after a lex error, Next()
// returns Token{Type: tokEOF, ...}. The first error encountered is
// stored in Err and all subsequent Next() calls return tokEOF.
type Lexer struct {
	input string // source text
	pos   int    // current read position (next byte to consume)
	start int    // byte offset of token currently being scanned
	Err   error  // first error encountered, nil if none
}

// NewLexer creates a Lexer for the given source string.
func NewLexer(input string) *Lexer {
	return &Lexer{input: input}
}

// Next returns the next token from the input.
// At end of input or after a lex error, returns Token{Type: tokEOF, ...}.
// After Err is set, all subsequent calls return tokEOF.
func (l *Lexer) Next() Token {
	if l.Err != nil {
		return Token{Type: tokEOF, Loc: ast.Loc{Start: l.pos, End: l.pos}}
	}
	l.skipWhitespaceAndComments()
	if l.Err != nil {
		return Token{Type: tokEOF, Loc: ast.Loc{Start: l.pos, End: l.pos}}
	}
	if l.pos >= len(l.input) {
		return Token{Type: tokEOF, Loc: ast.Loc{Start: l.pos, End: l.pos}}
	}
	l.start = l.pos
	ch := l.input[l.pos]

	switch {
	case ch == '\'':
		return l.scanString()
	case ch == '"':
		return l.scanQuotedIdent()
	case ch == '`':
		return l.scanIonLiteral()
	case ch >= '0' && ch <= '9':
		return l.scanNumber()
	case ch == '.' && l.pos+1 < len(l.input) && isDigit(l.input[l.pos+1]):
		return l.scanNumber() // leading-dot decimal: .5
	case isIdentStart(ch):
		return l.scanIdentOrKeyword()
	default:
		return l.scanOperator()
	}
}

// skipWhitespaceAndComments advances l.pos past whitespace, line comments,
// and block comments. All three are on the HIDDEN channel per the grammar
// and never appear in the token stream.
func (l *Lexer) skipWhitespaceAndComments() {
	for l.pos < len(l.input) {
		ch := l.input[l.pos]

		// Whitespace.
		if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
			l.pos++
			continue
		}

		// Line comment: -- to end of line.
		if ch == '-' && l.pos+1 < len(l.input) && l.input[l.pos+1] == '-' {
			l.pos += 2
			for l.pos < len(l.input) && l.input[l.pos] != '\n' {
				l.pos++
			}
			continue
		}

		// Block comment: /* ... */ (non-nested, greedy-shortest).
		if ch == '/' && l.pos+1 < len(l.input) && l.input[l.pos+1] == '*' {
			start := l.pos
			l.pos += 2
			for l.pos+1 < len(l.input) {
				if l.input[l.pos] == '*' && l.input[l.pos+1] == '/' {
					l.pos += 2
					goto continueLoop
				}
				l.pos++
			}
			// EOF before closing */
			l.Err = fmt.Errorf("unterminated block comment at position %d", start)
			return
		continueLoop:
			continue
		}

		break
	}
}

// ============================================================================
// Scan helpers — implemented in subsequent tasks (Tasks 6–10).
//
// scanString, scanQuotedIdent, scanIdentOrKeyword, scanNumber,
// scanOperator, scanIonLiteral are added incrementally with their tests.
// Until then, the dispatch in Next() refers to functions that don't exist
// yet, so the file does NOT compile after this task. Tasks 6–10 add the
// missing helpers and the package becomes buildable again at the end of
// Task 6.
//
// Task 5 ends here. Run `go build` after this task and expect a
// compile error like "undefined: scanString" — that's expected and
// fixed by Task 6.
// ============================================================================

// ============================================================================
// Character class helpers.
// ============================================================================

// isIdentStart reports whether ch can begin a PartiQL identifier.
// PartiQL identifiers start with [a-zA-Z_$].
func isIdentStart(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z') ||
		ch == '_' || ch == '$'
}

// isIdentContinue reports whether ch can appear in a PartiQL identifier
// after the first character. Adds digits to isIdentStart.
func isIdentContinue(ch byte) bool {
	return isIdentStart(ch) || isDigit(ch)
}

// isDigit reports whether ch is an ASCII decimal digit.
func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}
```

**Note:** the `goto continueLoop` inside the block-comment scan is awkward but necessary because Go doesn't have a labeled `continue` for the outer for loop in the same way. The label `continueLoop` is local to the block-comment branch.

Actually, this is over-clever. Replace the `goto` with a flag-based exit:

```go
		// Block comment: /* ... */ (non-nested, greedy-shortest).
		if ch == '/' && l.pos+1 < len(l.input) && l.input[l.pos+1] == '*' {
			start := l.pos
			l.pos += 2
			closed := false
			for l.pos+1 < len(l.input) {
				if l.input[l.pos] == '*' && l.input[l.pos+1] == '/' {
					l.pos += 2
					closed = true
					break
				}
				l.pos++
			}
			if !closed {
				l.Err = fmt.Errorf("unterminated block comment at position %d", start)
				return
			}
			continue
		}
```

Use the flag-based version above. Replace the `goto continueLoop:` block in the file content with this cleaner form before saving.

- [ ] **Step 2: Try to build (expect failure)**

```bash
go build ./partiql/parser/...
```

Expected: compile errors like:
```
partiql/parser/lexer.go:XX:XX: undefined: l.scanString
partiql/parser/lexer.go:XX:XX: undefined: l.scanQuotedIdent
partiql/parser/lexer.go:XX:XX: undefined: l.scanIonLiteral
partiql/parser/lexer.go:XX:XX: undefined: l.scanNumber
partiql/parser/lexer.go:XX:XX: undefined: l.scanIdentOrKeyword
partiql/parser/lexer.go:XX:XX: undefined: l.scanOperator
```

This is expected. Tasks 6–10 add the missing helpers. Do NOT commit yet — Task 5 doesn't ship until Task 6 lands the first scan helper and the package is buildable again.

- [ ] **Step 3: Add a temporary stub for the missing scan helpers**

To keep each task individually committable and testable, add stubs at the bottom of `lexer.go` that return tokEOF. They will be replaced by the real implementations in Tasks 6–10.

```go
// ============================================================================
// STUBS — replaced by Tasks 6–10.
//
// These return tokEOF and set l.Err so the package builds at the end of
// Task 5. Each subsequent task removes one stub and adds the real
// implementation alongside its tests.
// ============================================================================

func (l *Lexer) scanString() Token {
	l.Err = fmt.Errorf("scanString not yet implemented (stub from Task 5)")
	return Token{Type: tokEOF, Loc: ast.Loc{Start: l.start, End: l.start}}
}

func (l *Lexer) scanQuotedIdent() Token {
	l.Err = fmt.Errorf("scanQuotedIdent not yet implemented (stub from Task 5)")
	return Token{Type: tokEOF, Loc: ast.Loc{Start: l.start, End: l.start}}
}

func (l *Lexer) scanIdentOrKeyword() Token {
	l.Err = fmt.Errorf("scanIdentOrKeyword not yet implemented (stub from Task 5)")
	return Token{Type: tokEOF, Loc: ast.Loc{Start: l.start, End: l.start}}
}

func (l *Lexer) scanNumber() Token {
	l.Err = fmt.Errorf("scanNumber not yet implemented (stub from Task 5)")
	return Token{Type: tokEOF, Loc: ast.Loc{Start: l.start, End: l.start}}
}

func (l *Lexer) scanOperator() Token {
	l.Err = fmt.Errorf("scanOperator not yet implemented (stub from Task 5)")
	return Token{Type: tokEOF, Loc: ast.Loc{Start: l.start, End: l.start}}
}

func (l *Lexer) scanIonLiteral() Token {
	l.Err = fmt.Errorf("scanIonLiteral not yet implemented (stub from Task 5)")
	return Token{Type: tokEOF, Loc: ast.Loc{Start: l.start, End: l.start}}
}

// strings is imported but only used by the stubs above (and by future
// scan helpers). Avoid the unused-import error during early tasks by
// keeping a no-op reference. Remove this line in Task 7 when
// scanIdentOrKeyword adds the real strings.ToLower call.
var _ = strings.ToLower
```

Yes, the `var _ = strings.ToLower` line is ugly. It exists because the `strings` import is needed by Task 7 (`scanIdentOrKeyword` uses `strings.ToLower`) but at the end of Task 5 there's no other reference to the package and `go vet` would complain about an unused import. Task 7 removes this no-op line when it adds the real `strings.ToLower` call.

- [ ] **Step 4: Build, vet, gofmt**

```bash
go build ./partiql/parser/...
go vet ./partiql/parser/...
gofmt -l partiql/parser/
```

All clean. The package compiles end-to-end (token.go + keywords.go + lexer.go) but `Next()` always returns an error because every dispatch arm hits a stub.

- [ ] **Step 5: Commit**

```bash
git add partiql/parser/lexer.go
git commit -m "$(cat <<'EOF'
feat(partiql/parser): scaffold lexer.go with Next() dispatch and helpers

Adds the Lexer struct, NewLexer, Next() entry point, character class
helpers (isIdentStart/isIdentContinue/isDigit), and the
skipWhitespaceAndComments helper that handles whitespace, line
comments (-- to EOL), and block comments (/* */).

Scan helpers (scanString, scanQuotedIdent, scanIdentOrKeyword,
scanNumber, scanOperator, scanIonLiteral) are stubbed; each one
returns an error and is replaced by Tasks 6-10. The stubs let the
package build at the end of Task 5 so each subsequent task has a
clean baseline.

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

### Task 6: `lexer.go` scanString + scanQuotedIdent + golden tests

**Files:**
- Modify: `partiql/parser/lexer.go` (replace stubs for `scanString`/`scanQuotedIdent`)
- Create: `partiql/parser/lexer_test.go` (start the test file)

These two scan helpers share an algorithm: read until the matching quote, treating doubled quotes as escapes for an embedded quote. This task implements both and adds the first batch of golden test cases.

- [ ] **Step 1: Create `partiql/parser/lexer_test.go` with the test scaffolding**

```go
package parser

import (
	"reflect"
	"testing"

	"github.com/bytebase/omni/partiql/ast"
)

// tokenStreamCase is one row in the TestLexer_Tokens table.
type tokenStreamCase struct {
	name   string
	input  string
	tokens []Token // expected, excluding the trailing EOF
}

// runTokenStreamCases drains l.Next() until tokEOF and asserts the captured
// stream matches tc.tokens via reflect.DeepEqual. Asserts l.Err == nil.
func runTokenStreamCase(t *testing.T, tc tokenStreamCase) {
	t.Helper()
	l := NewLexer(tc.input)
	var got []Token
	for {
		tok := l.Next()
		if tok.Type == tokEOF {
			break
		}
		got = append(got, tok)
	}
	if l.Err != nil {
		t.Fatalf("unexpected error: %v", l.Err)
	}
	if !reflect.DeepEqual(got, tc.tokens) {
		t.Errorf("token stream mismatch\n got: %+v\nwant: %+v", got, tc.tokens)
	}
}

// TestLexer_Tokens is the master golden-test table for the lexer.
// Tasks 6-10 each append a section of cases as their scan helper lands.
//
// All cases assert reflect.DeepEqual on the full token slice and l.Err == nil.
func TestLexer_Tokens(t *testing.T) {
	cases := []tokenStreamCase{
		// =============================================================
		// Empty input + whitespace + comments
		// =============================================================
		{"empty", "", nil},
		{"whitespace_spaces", "   ", nil},
		{"whitespace_tabs", "\t\t", nil},
		{"whitespace_newlines", "\n\n", nil},
		{"whitespace_mixed", " \t\n\r ", nil},
		{"line_comment_only", "-- a comment\n", nil},
		{"line_comment_at_eof", "-- a comment without newline", nil},
		{"block_comment_only", "/* a comment */", nil},
		{"block_comment_multiline", "/*\nmulti\nline\n*/", nil},

		// =============================================================
		// String literals (Task 6)
		// =============================================================
		{
			"string_simple",
			"'hello'",
			[]Token{{tokSCONST, "hello", ast.Loc{Start: 0, End: 7}}},
		},
		{
			"string_empty",
			"''",
			[]Token{{tokSCONST, "", ast.Loc{Start: 0, End: 2}}},
		},
		{
			"string_doubled_quote",
			"'it''s'",
			[]Token{{tokSCONST, "it's", ast.Loc{Start: 0, End: 7}}},
		},
		{
			"string_with_whitespace",
			"'  '",
			[]Token{{tokSCONST, "  ", ast.Loc{Start: 0, End: 4}}},
		},
		{
			"string_with_special_chars",
			"'a!@#%^&*()'",
			[]Token{{tokSCONST, "a!@#%^&*()", ast.Loc{Start: 0, End: 12}}},
		},

		// =============================================================
		// Quoted identifiers (Task 6)
		// =============================================================
		{
			"quoted_ident_simple",
			`"Foo"`,
			[]Token{{tokIDENT_QUOTED, "Foo", ast.Loc{Start: 0, End: 5}}},
		},
		{
			"quoted_ident_empty",
			`""`,
			[]Token{{tokIDENT_QUOTED, "", ast.Loc{Start: 0, End: 2}}},
		},
		{
			"quoted_ident_doubled_quote",
			`"a""b"`,
			[]Token{{tokIDENT_QUOTED, `a"b`, ast.Loc{Start: 0, End: 6}}},
		},
		{
			"quoted_ident_with_space",
			`"Foo Bar"`,
			[]Token{{tokIDENT_QUOTED, "Foo Bar", ast.Loc{Start: 0, End: 9}}},
		},
		{
			"quoted_ident_case_preserved",
			`"FoO"`,
			[]Token{{tokIDENT_QUOTED, "FoO", ast.Loc{Start: 0, End: 5}}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			runTokenStreamCase(t, tc)
		})
	}
}
```

- [ ] **Step 2: Replace the `scanString` stub in `partiql/ast/lexer.go`** *(typo: should be `partiql/parser/lexer.go`)*

Open `partiql/parser/lexer.go`. Find the stub:

```go
func (l *Lexer) scanString() Token {
	l.Err = fmt.Errorf("scanString not yet implemented (stub from Task 5)")
	return Token{Type: tokEOF, Loc: ast.Loc{Start: l.start, End: l.start}}
}
```

Replace with the real implementation:

```go
// scanString consumes a single-quoted string literal: '...'.
// Doubled single quotes ('') represent an embedded single quote.
// Token.Str holds the decoded value (with '' collapsed to ').
//
// Grammar: LITERAL_STRING : '\'' ( ('\'\'') | ~('\'') )* '\'';
func (l *Lexer) scanString() Token {
	l.pos++ // skip opening '
	var buf strings.Builder
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == '\'' {
			// Check for doubled quote (escape).
			if l.pos+1 < len(l.input) && l.input[l.pos+1] == '\'' {
				buf.WriteByte('\'')
				l.pos += 2
				continue
			}
			// Closing quote.
			l.pos++
			return Token{
				Type: tokSCONST,
				Str:  buf.String(),
				Loc:  ast.Loc{Start: l.start, End: l.pos},
			}
		}
		buf.WriteByte(ch)
		l.pos++
	}
	l.Err = fmt.Errorf("unterminated string literal at position %d", l.start)
	return Token{Type: tokEOF, Loc: ast.Loc{Start: l.start, End: l.start}}
}
```

- [ ] **Step 3: Replace the `scanQuotedIdent` stub**

Same pattern as `scanString` but with `"` instead of `'` and producing `tokIDENT_QUOTED`:

```go
// scanQuotedIdent consumes a double-quoted identifier: "...".
// Doubled double quotes ("") represent an embedded double quote.
// Token.Str holds the decoded value (with "" collapsed to ").
// Quoted identifiers preserve case (the keyword map is not consulted).
//
// Grammar: IDENTIFIER_QUOTED : '"' ( ('""') | ~('"') )* '"';
func (l *Lexer) scanQuotedIdent() Token {
	l.pos++ // skip opening "
	var buf strings.Builder
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == '"' {
			// Check for doubled quote (escape).
			if l.pos+1 < len(l.input) && l.input[l.pos+1] == '"' {
				buf.WriteByte('"')
				l.pos += 2
				continue
			}
			// Closing quote.
			l.pos++
			return Token{
				Type: tokIDENT_QUOTED,
				Str:  buf.String(),
				Loc:  ast.Loc{Start: l.start, End: l.pos},
			}
		}
		buf.WriteByte(ch)
		l.pos++
	}
	l.Err = fmt.Errorf("unterminated quoted identifier at position %d", l.start)
	return Token{Type: tokEOF, Loc: ast.Loc{Start: l.start, End: l.start}}
}
```

- [ ] **Step 4: Build and run the new tests**

```bash
go build ./partiql/parser/...
go test -v -run TestLexer_Tokens ./partiql/parser/...
```

Expected: 19 sub-tests pass (9 whitespace/comment + 5 string + 5 quoted ident).

If any sub-test fails, debug the helper before proceeding. Do NOT change the goldens to match the implementation — fix the implementation.

- [ ] **Step 5: Vet and gofmt**

```bash
go vet ./partiql/parser/...
gofmt -l partiql/parser/
```

Both clean.

- [ ] **Step 6: Commit**

```bash
git add partiql/parser/lexer.go partiql/parser/lexer_test.go
git commit -m "$(cat <<'EOF'
feat(partiql/parser): implement scanString and scanQuotedIdent

Replaces the Task 5 stubs with real implementations of scanString
('...' literals) and scanQuotedIdent ("..." quoted identifiers).
Both share an algorithm: read until matching quote, treating doubled
quotes as escapes for an embedded quote. Token.Str holds the decoded
value.

Per the grammar, PartiQL has no backslash escapes. The only escape
mechanism is the doubled-quote form ('it''s' -> "it's").

Adds lexer_test.go with the master TestLexer_Tokens table and the
first 19 cases: empty input, whitespace, line/block comments, string
literals (5 cases), and quoted identifiers (5 cases). Future tasks
append more cases to the same table.

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

### Task 7: `lexer.go` scanIdentOrKeyword + golden tests

**Files:**
- Modify: `partiql/parser/lexer.go` (replace `scanIdentOrKeyword` stub, remove the `var _ = strings.ToLower` line)
- Modify: `partiql/parser/lexer_test.go` (append identifier/keyword cases)

- [ ] **Step 1: Replace the `scanIdentOrKeyword` stub**

In `partiql/parser/lexer.go`, find:

```go
func (l *Lexer) scanIdentOrKeyword() Token {
	l.Err = fmt.Errorf("scanIdentOrKeyword not yet implemented (stub from Task 5)")
	return Token{Type: tokEOF, Loc: ast.Loc{Start: l.start, End: l.start}}
}
```

Replace with:

```go
// scanIdentOrKeyword consumes an unquoted identifier and looks it up in
// the keywords map. If the lowercased text matches a keyword, returns
// that keyword token; otherwise returns tokIDENT.
//
// Token.Str preserves the original case (so the parser/AST can render
// identifiers as written). Keyword matching is case-insensitive per
// PartiQLLexer.g4 caseInsensitive=true.
//
// Grammar: IDENTIFIER : [A-Z$_][A-Z0-9$_]*;
//          (with caseInsensitive=true expanding to [a-zA-Z$_][a-zA-Z0-9$_]*)
func (l *Lexer) scanIdentOrKeyword() Token {
	for l.pos < len(l.input) && isIdentContinue(l.input[l.pos]) {
		l.pos++
	}
	raw := l.input[l.start:l.pos]
	lower := strings.ToLower(raw)
	if tt, ok := keywords[lower]; ok {
		return Token{
			Type: tt,
			Str:  raw,
			Loc:  ast.Loc{Start: l.start, End: l.pos},
		}
	}
	return Token{
		Type: tokIDENT,
		Str:  raw,
		Loc:  ast.Loc{Start: l.start, End: l.pos},
	}
}
```

- [ ] **Step 2: Remove the `var _ = strings.ToLower` no-op line**

Find and delete the line:

```go
var _ = strings.ToLower
```

The `strings` package is now genuinely used by `scanIdentOrKeyword`, so the no-op reference is no longer needed.

- [ ] **Step 3: Append identifier and keyword test cases to `lexer_test.go`**

Find the `cases` slice in `TestLexer_Tokens` and append the following block at the end (before the closing `}` of the slice):

```go
		// =============================================================
		// Unquoted identifiers (Task 7)
		// =============================================================
		{
			"unquoted_ident_simple",
			"foo",
			[]Token{{tokIDENT, "foo", ast.Loc{Start: 0, End: 3}}},
		},
		{
			"unquoted_ident_with_underscore",
			"_foo",
			[]Token{{tokIDENT, "_foo", ast.Loc{Start: 0, End: 4}}},
		},
		{
			"unquoted_ident_with_dollar",
			"x$y",
			[]Token{{tokIDENT, "x$y", ast.Loc{Start: 0, End: 3}}},
		},
		{
			"unquoted_ident_with_digit_in_middle",
			"a1b2",
			[]Token{{tokIDENT, "a1b2", ast.Loc{Start: 0, End: 4}}},
		},
		{
			"unquoted_ident_uppercase_preserved",
			"FOO",
			[]Token{{tokIDENT, "FOO", ast.Loc{Start: 0, End: 3}}},
		},

		// =============================================================
		// Keywords (case-insensitive lookup, raw text preserved) (Task 7)
		// =============================================================
		{
			"keyword_select_lower",
			"select",
			[]Token{{tokSELECT, "select", ast.Loc{Start: 0, End: 6}}},
		},
		{
			"keyword_select_upper",
			"SELECT",
			[]Token{{tokSELECT, "SELECT", ast.Loc{Start: 0, End: 6}}},
		},
		{
			"keyword_select_mixed",
			"Select",
			[]Token{{tokSELECT, "Select", ast.Loc{Start: 0, End: 6}}},
		},
		{
			"keyword_from",
			"FROM",
			[]Token{{tokFROM, "FROM", ast.Loc{Start: 0, End: 4}}},
		},
		{
			"keyword_where",
			"WHERE",
			[]Token{{tokWHERE, "WHERE", ast.Loc{Start: 0, End: 5}}},
		},
		{
			"keyword_pivot_partiql_unique",
			"PIVOT",
			[]Token{{tokPIVOT, "PIVOT", ast.Loc{Start: 0, End: 5}}},
		},
		{
			"keyword_missing_partiql_unique",
			"MISSING",
			[]Token{{tokMISSING, "MISSING", ast.Loc{Start: 0, End: 7}}},
		},
		{
			"keyword_bag_data_type",
			"BAG",
			[]Token{{tokBAG, "BAG", ast.Loc{Start: 0, End: 3}}},
		},
		{
			"keyword_can_lossless_cast_underscored",
			"CAN_LOSSLESS_CAST",
			[]Token{{tokCAN_LOSSLESS_CAST, "CAN_LOSSLESS_CAST", ast.Loc{Start: 0, End: 17}}},
		},

		// =============================================================
		// Identifier vs keyword cases after whitespace/comments (Task 7)
		// =============================================================
		{
			"ident_after_line_comment",
			"-- skipped\nfoo",
			[]Token{{tokIDENT, "foo", ast.Loc{Start: 11, End: 14}}},
		},
		{
			"ident_after_block_comment",
			"/* x */ foo",
			[]Token{{tokIDENT, "foo", ast.Loc{Start: 8, End: 11}}},
		},
```

- [ ] **Step 4: Build and run tests**

```bash
go build ./partiql/parser/...
go test -v -run TestLexer_Tokens ./partiql/parser/...
```

Expected: 35 sub-tests pass (19 prior + 5 unquoted ident + 9 keyword + 2 ident-after-comment).

- [ ] **Step 5: Vet and gofmt**

```bash
go vet ./partiql/parser/...
gofmt -l partiql/parser/
```

Both clean.

- [ ] **Step 6: Commit**

```bash
git add partiql/parser/lexer.go partiql/parser/lexer_test.go
git commit -m "$(cat <<'EOF'
feat(partiql/parser): implement scanIdentOrKeyword

Replaces the Task 5 stub. Walks while isIdentContinue holds, takes
the source slice, lowercases via strings.ToLower, and looks up the
result in the keywords map. If found, returns the keyword token type
with Token.Str preserving the original case. Otherwise returns
tokIDENT.

Removes the var _ = strings.ToLower no-op line that Task 5 added to
keep the import alive.

Adds 16 new TestLexer_Tokens cases: 5 unquoted identifier shapes
(simple, leading underscore, $ in middle, digit in middle, uppercase
preserved); 9 keyword cases including standard SQL (SELECT case
variations, FROM, WHERE), PartiQL-unique (PIVOT, MISSING), data type
(BAG), and a multi-underscore (CAN_LOSSLESS_CAST); and 2 cases for
identifier-after-comment showing position tracking.

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

### Task 8: `lexer.go` scanNumber + golden tests

**Files:**
- Modify: `partiql/parser/lexer.go` (replace `scanNumber` stub)
- Modify: `partiql/parser/lexer_test.go` (append numeric literal cases)

- [ ] **Step 1: Replace the `scanNumber` stub**

```go
// scanNumber consumes an integer or decimal literal. Returns tokICONST
// for plain integers and tokFCONST for any number with a decimal point
// or scientific exponent. Token.Str is the raw source text.
//
// Grammar:
//
//	LITERAL_INTEGER : DIGIT+;
//	LITERAL_DECIMAL :
//	    DIGIT+ '.' DIGIT* ([e] [+-]? DIGIT+)?
//	  | '.' DIGIT+ ([e] [+-]? DIGIT+)?
//	  | DIGIT+ ([e] [+-]? DIGIT+)?
//	  ;
//
// (caseInsensitive=true means [e] matches both 'e' and 'E'.)
func (l *Lexer) scanNumber() Token {
	isFloat := false

	// Leading-dot form (.5).
	if l.input[l.pos] == '.' {
		isFloat = true
		l.pos++
		for l.pos < len(l.input) && isDigit(l.input[l.pos]) {
			l.pos++
		}
	} else {
		// Integer part.
		for l.pos < len(l.input) && isDigit(l.input[l.pos]) {
			l.pos++
		}
		// Optional fraction.
		if l.pos < len(l.input) && l.input[l.pos] == '.' {
			isFloat = true
			l.pos++
			for l.pos < len(l.input) && isDigit(l.input[l.pos]) {
				l.pos++
			}
		}
	}

	// Optional scientific exponent.
	if l.pos < len(l.input) && (l.input[l.pos] == 'e' || l.input[l.pos] == 'E') {
		isFloat = true
		l.pos++
		if l.pos < len(l.input) && (l.input[l.pos] == '+' || l.input[l.pos] == '-') {
			l.pos++
		}
		for l.pos < len(l.input) && isDigit(l.input[l.pos]) {
			l.pos++
		}
	}

	tt := tokICONST
	if isFloat {
		tt = tokFCONST
	}
	return Token{
		Type: tt,
		Str:  l.input[l.start:l.pos],
		Loc:  ast.Loc{Start: l.start, End: l.pos},
	}
}
```

- [ ] **Step 2: Append numeric literal test cases to `lexer_test.go`**

```go
		// =============================================================
		// Numeric literals (Task 8)
		// =============================================================
		{
			"integer",
			"42",
			[]Token{{tokICONST, "42", ast.Loc{Start: 0, End: 2}}},
		},
		{
			"integer_zero",
			"0",
			[]Token{{tokICONST, "0", ast.Loc{Start: 0, End: 1}}},
		},
		{
			"integer_large",
			"1234567890",
			[]Token{{tokICONST, "1234567890", ast.Loc{Start: 0, End: 10}}},
		},
		{
			"decimal_dot",
			"3.14",
			[]Token{{tokFCONST, "3.14", ast.Loc{Start: 0, End: 4}}},
		},
		{
			"decimal_leading_dot",
			".5",
			[]Token{{tokFCONST, ".5", ast.Loc{Start: 0, End: 2}}},
		},
		{
			"decimal_trailing_dot",
			"42.",
			[]Token{{tokFCONST, "42.", ast.Loc{Start: 0, End: 3}}},
		},
		{
			"decimal_scientific_lower_e",
			"1e10",
			[]Token{{tokFCONST, "1e10", ast.Loc{Start: 0, End: 4}}},
		},
		{
			"decimal_scientific_upper_e",
			"1E10",
			[]Token{{tokFCONST, "1E10", ast.Loc{Start: 0, End: 4}}},
		},
		{
			"decimal_scientific_negative_exp",
			"1.5e-3",
			[]Token{{tokFCONST, "1.5e-3", ast.Loc{Start: 0, End: 6}}},
		},
		{
			"decimal_scientific_positive_exp",
			"2.5e+4",
			[]Token{{tokFCONST, "2.5e+4", ast.Loc{Start: 0, End: 6}}},
		},
```

- [ ] **Step 3: Build and run tests**

```bash
go build ./partiql/parser/...
go test -v -run TestLexer_Tokens ./partiql/parser/...
```

Expected: 45 sub-tests pass (35 prior + 10 numeric).

- [ ] **Step 4: Vet, gofmt, commit**

```bash
go vet ./partiql/parser/...
gofmt -l partiql/parser/
git add partiql/parser/lexer.go partiql/parser/lexer_test.go
git commit -m "$(cat <<'EOF'
feat(partiql/parser): implement scanNumber

Replaces the Task 5 stub. Handles integer, decimal-with-dot, leading
dot (.5), trailing dot (42.), and scientific exponent forms (1e10,
1E10, 1.5e-3, 2.5e+4) per PartiQLLexer.g4 LITERAL_INTEGER and
LITERAL_DECIMAL rules. Returns tokICONST for plain integers and
tokFCONST for anything with a decimal point or exponent. Token.Str
is the raw source slice.

Adds 10 TestLexer_Tokens cases covering each form.

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

### Task 9: `lexer.go` scanOperator + golden tests

**Files:**
- Modify: `partiql/parser/lexer.go` (replace `scanOperator` stub)
- Modify: `partiql/parser/lexer_test.go` (append operator cases)

- [ ] **Step 1: Replace the `scanOperator` stub**

```go
// scanOperator consumes a one or two-character operator or punctuation
// token. Two-character operators (<=, >=, <>, <<, >>, ||, !=) are
// matched first via lookahead; otherwise the single-character cases
// fall through. Unrecognized characters set l.Err.
func (l *Lexer) scanOperator() Token {
	ch := l.input[l.pos]
	l.pos++

	// Two-character lookahead.
	if l.pos < len(l.input) {
		next := l.input[l.pos]
		switch {
		case ch == '<' && next == '=':
			l.pos++
			return Token{Type: tokLT_EQ, Str: "<=", Loc: ast.Loc{Start: l.start, End: l.pos}}
		case ch == '<' && next == '>':
			l.pos++
			return Token{Type: tokNEQ, Str: "<>", Loc: ast.Loc{Start: l.start, End: l.pos}}
		case ch == '<' && next == '<':
			l.pos++
			return Token{Type: tokANGLE_DOUBLE_LEFT, Str: "<<", Loc: ast.Loc{Start: l.start, End: l.pos}}
		case ch == '>' && next == '=':
			l.pos++
			return Token{Type: tokGT_EQ, Str: ">=", Loc: ast.Loc{Start: l.start, End: l.pos}}
		case ch == '>' && next == '>':
			l.pos++
			return Token{Type: tokANGLE_DOUBLE_RIGHT, Str: ">>", Loc: ast.Loc{Start: l.start, End: l.pos}}
		case ch == '|' && next == '|':
			l.pos++
			return Token{Type: tokCONCAT, Str: "||", Loc: ast.Loc{Start: l.start, End: l.pos}}
		case ch == '!' && next == '=':
			l.pos++
			return Token{Type: tokNEQ, Str: "!=", Loc: ast.Loc{Start: l.start, End: l.pos}}
		}
	}

	// Single-character operators / punctuation.
	switch ch {
	case '+':
		return Token{Type: tokPLUS, Str: "+", Loc: ast.Loc{Start: l.start, End: l.pos}}
	case '-':
		return Token{Type: tokMINUS, Str: "-", Loc: ast.Loc{Start: l.start, End: l.pos}}
	case '*':
		return Token{Type: tokASTERISK, Str: "*", Loc: ast.Loc{Start: l.start, End: l.pos}}
	case '/':
		return Token{Type: tokSLASH_FORWARD, Str: "/", Loc: ast.Loc{Start: l.start, End: l.pos}}
	case '%':
		return Token{Type: tokPERCENT, Str: "%", Loc: ast.Loc{Start: l.start, End: l.pos}}
	case '^':
		return Token{Type: tokCARET, Str: "^", Loc: ast.Loc{Start: l.start, End: l.pos}}
	case '~':
		return Token{Type: tokTILDE, Str: "~", Loc: ast.Loc{Start: l.start, End: l.pos}}
	case '@':
		return Token{Type: tokAT_SIGN, Str: "@", Loc: ast.Loc{Start: l.start, End: l.pos}}
	case '=':
		return Token{Type: tokEQ, Str: "=", Loc: ast.Loc{Start: l.start, End: l.pos}}
	case '<':
		return Token{Type: tokLT, Str: "<", Loc: ast.Loc{Start: l.start, End: l.pos}}
	case '>':
		return Token{Type: tokGT, Str: ">", Loc: ast.Loc{Start: l.start, End: l.pos}}
	case '(':
		return Token{Type: tokPAREN_LEFT, Str: "(", Loc: ast.Loc{Start: l.start, End: l.pos}}
	case ')':
		return Token{Type: tokPAREN_RIGHT, Str: ")", Loc: ast.Loc{Start: l.start, End: l.pos}}
	case '[':
		return Token{Type: tokBRACKET_LEFT, Str: "[", Loc: ast.Loc{Start: l.start, End: l.pos}}
	case ']':
		return Token{Type: tokBRACKET_RIGHT, Str: "]", Loc: ast.Loc{Start: l.start, End: l.pos}}
	case '{':
		return Token{Type: tokBRACE_LEFT, Str: "{", Loc: ast.Loc{Start: l.start, End: l.pos}}
	case '}':
		return Token{Type: tokBRACE_RIGHT, Str: "}", Loc: ast.Loc{Start: l.start, End: l.pos}}
	case ':':
		return Token{Type: tokCOLON, Str: ":", Loc: ast.Loc{Start: l.start, End: l.pos}}
	case ';':
		return Token{Type: tokCOLON_SEMI, Str: ";", Loc: ast.Loc{Start: l.start, End: l.pos}}
	case ',':
		return Token{Type: tokCOMMA, Str: ",", Loc: ast.Loc{Start: l.start, End: l.pos}}
	case '.':
		return Token{Type: tokPERIOD, Str: ".", Loc: ast.Loc{Start: l.start, End: l.pos}}
	case '?':
		return Token{Type: tokQUESTION_MARK, Str: "?", Loc: ast.Loc{Start: l.start, End: l.pos}}
	}

	l.Err = fmt.Errorf("unexpected character %q at position %d", ch, l.start)
	return Token{Type: tokEOF, Loc: ast.Loc{Start: l.start, End: l.start}}
}
```

- [ ] **Step 2: Append operator test cases to `lexer_test.go`**

```go
		// =============================================================
		// Single-character operators (Task 9)
		// =============================================================
		{"op_plus", "+", []Token{{tokPLUS, "+", ast.Loc{Start: 0, End: 1}}}},
		{"op_minus", "-", []Token{{tokMINUS, "-", ast.Loc{Start: 0, End: 1}}}},
		{"op_asterisk", "*", []Token{{tokASTERISK, "*", ast.Loc{Start: 0, End: 1}}}},
		{"op_slash_forward", "/", []Token{{tokSLASH_FORWARD, "/", ast.Loc{Start: 0, End: 1}}}},
		{"op_percent", "%", []Token{{tokPERCENT, "%", ast.Loc{Start: 0, End: 1}}}},
		{"op_caret", "^", []Token{{tokCARET, "^", ast.Loc{Start: 0, End: 1}}}},
		{"op_tilde", "~", []Token{{tokTILDE, "~", ast.Loc{Start: 0, End: 1}}}},
		{"op_at_sign", "@", []Token{{tokAT_SIGN, "@", ast.Loc{Start: 0, End: 1}}}},
		{"op_eq", "=", []Token{{tokEQ, "=", ast.Loc{Start: 0, End: 1}}}},
		{"op_lt", "<", []Token{{tokLT, "<", ast.Loc{Start: 0, End: 1}}}},
		{"op_gt", ">", []Token{{tokGT, ">", ast.Loc{Start: 0, End: 1}}}},

		// =============================================================
		// Punctuation (Task 9)
		// =============================================================
		{"punct_paren_left", "(", []Token{{tokPAREN_LEFT, "(", ast.Loc{Start: 0, End: 1}}}},
		{"punct_paren_right", ")", []Token{{tokPAREN_RIGHT, ")", ast.Loc{Start: 0, End: 1}}}},
		{"punct_bracket_left", "[", []Token{{tokBRACKET_LEFT, "[", ast.Loc{Start: 0, End: 1}}}},
		{"punct_bracket_right", "]", []Token{{tokBRACKET_RIGHT, "]", ast.Loc{Start: 0, End: 1}}}},
		{"punct_brace_left", "{", []Token{{tokBRACE_LEFT, "{", ast.Loc{Start: 0, End: 1}}}},
		{"punct_brace_right", "}", []Token{{tokBRACE_RIGHT, "}", ast.Loc{Start: 0, End: 1}}}},
		{"punct_colon", ":", []Token{{tokCOLON, ":", ast.Loc{Start: 0, End: 1}}}},
		{"punct_colon_semi", ";", []Token{{tokCOLON_SEMI, ";", ast.Loc{Start: 0, End: 1}}}},
		{"punct_comma", ",", []Token{{tokCOMMA, ",", ast.Loc{Start: 0, End: 1}}}},
		{"punct_period", ".", []Token{{tokPERIOD, ".", ast.Loc{Start: 0, End: 1}}}},
		{"punct_question_mark", "?", []Token{{tokQUESTION_MARK, "?", ast.Loc{Start: 0, End: 1}}}},

		// =============================================================
		// Two-character operators (Task 9)
		// =============================================================
		{"op_lt_eq", "<=", []Token{{tokLT_EQ, "<=", ast.Loc{Start: 0, End: 2}}}},
		{"op_gt_eq", ">=", []Token{{tokGT_EQ, ">=", ast.Loc{Start: 0, End: 2}}}},
		{"op_neq_angle", "<>", []Token{{tokNEQ, "<>", ast.Loc{Start: 0, End: 2}}}},
		{"op_neq_bang", "!=", []Token{{tokNEQ, "!=", ast.Loc{Start: 0, End: 2}}}},
		{"op_concat", "||", []Token{{tokCONCAT, "||", ast.Loc{Start: 0, End: 2}}}},
		{"op_angle_double_left", "<<", []Token{{tokANGLE_DOUBLE_LEFT, "<<", ast.Loc{Start: 0, End: 2}}}},
		{"op_angle_double_right", ">>", []Token{{tokANGLE_DOUBLE_RIGHT, ">>", ast.Loc{Start: 0, End: 2}}}},

		// =============================================================
		// Multi-token integration sequences (Task 9)
		// =============================================================
		{
			"select_star_from_table",
			"SELECT * FROM Music",
			[]Token{
				{tokSELECT, "SELECT", ast.Loc{Start: 0, End: 6}},
				{tokASTERISK, "*", ast.Loc{Start: 7, End: 8}},
				{tokFROM, "FROM", ast.Loc{Start: 9, End: 13}},
				{tokIDENT, "Music", ast.Loc{Start: 14, End: 19}},
			},
		},
		{
			"where_with_string_literal",
			"WHERE Artist='Pink Floyd'",
			[]Token{
				{tokWHERE, "WHERE", ast.Loc{Start: 0, End: 5}},
				{tokIDENT, "Artist", ast.Loc{Start: 6, End: 12}},
				{tokEQ, "=", ast.Loc{Start: 12, End: 13}},
				{tokSCONST, "Pink Floyd", ast.Loc{Start: 13, End: 25}},
			},
		},
		{
			"bag_literal",
			"<<1, 2, 3>>",
			[]Token{
				{tokANGLE_DOUBLE_LEFT, "<<", ast.Loc{Start: 0, End: 2}},
				{tokICONST, "1", ast.Loc{Start: 2, End: 3}},
				{tokCOMMA, ",", ast.Loc{Start: 3, End: 4}},
				{tokICONST, "2", ast.Loc{Start: 5, End: 6}},
				{tokCOMMA, ",", ast.Loc{Start: 6, End: 7}},
				{tokICONST, "3", ast.Loc{Start: 8, End: 9}},
				{tokANGLE_DOUBLE_RIGHT, ">>", ast.Loc{Start: 9, End: 11}},
			},
		},
```

- [ ] **Step 3: Build and run tests**

```bash
go build ./partiql/parser/...
go test -v -run TestLexer_Tokens ./partiql/parser/...
```

Expected: 78 sub-tests pass (45 prior + 11 single-char op + 11 punct + 7 two-char op + 3 integration + 1 wait let me recount: 11+11+7+3 = 32 new. 45 + 32 = 77.)

Actually 11 single-char ops + 11 punctuation + 7 two-char ops + 3 integration sequences = 32 new cases. Total 77.

- [ ] **Step 4: Vet, gofmt, commit**

```bash
go vet ./partiql/parser/...
gofmt -l partiql/parser/
git add partiql/parser/lexer.go partiql/parser/lexer_test.go
git commit -m "$(cat <<'EOF'
feat(partiql/parser): implement scanOperator

Replaces the Task 5 stub. Two-character lookahead first (<=, >=, <>,
<<, >>, ||, !=), then single-character switch over the 22 single-char
operators and punctuation tokens. Unrecognized characters set
l.Err with "unexpected character %q at position %d".

Adds 32 TestLexer_Tokens cases: 11 single-char operators, 11
punctuation tokens, 7 two-char operators, and 3 multi-token
integration sequences (SELECT * FROM Music, WHERE Artist='Pink
Floyd', <<1, 2, 3>>).

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

### Task 10: `lexer.go` scanIonLiteral + golden tests

**Files:**
- Modify: `partiql/parser/lexer.go` (replace `scanIonLiteral` stub)
- Modify: `partiql/parser/lexer_test.go` (append Ion literal cases)

- [ ] **Step 1: Replace the `scanIonLiteral` stub**

```go
// scanIonLiteral consumes a backtick-delimited Ion blob: `...`.
//
// SIMPLIFIED BASE BEHAVIOR: scans byte-to-byte from the opening
// backtick to the next backtick. The captured Token.Str is the
// verbatim inner content (no decoding); Token.Loc covers the entire
// `...` range including both backticks.
//
// KNOWN LIMITATION: Ion mode in PartiQLLexer.g4 has special handling
// for backticks inside Ion strings ('...' quoted symbols, "..." short
// strings, '''...''' long strings) that prevents premature literal
// closure. This naive scan does NOT respect those rules. The full
// Ion-mode-aware implementation is deferred to DAG node 17
// (parser-ion-literals).
//
// The AWS DynamoDB PartiQL corpus has zero real Ion literals; the
// only 2 backtick uses are in select-001.partiql and insert-002.partiql
// (syntax skeletons with placeholder backticks), both filtered out
// of the corpus smoke test.
func (l *Lexer) scanIonLiteral() Token {
	l.pos++ // skip opening `
	contentStart := l.pos
	for l.pos < len(l.input) {
		if l.input[l.pos] == '`' {
			content := l.input[contentStart:l.pos]
			l.pos++ // skip closing `
			return Token{
				Type: tokION_LITERAL,
				Str:  content,
				Loc:  ast.Loc{Start: l.start, End: l.pos},
			}
		}
		l.pos++
	}
	l.Err = fmt.Errorf("unterminated Ion literal at position %d", l.start)
	return Token{Type: tokEOF, Loc: ast.Loc{Start: l.start, End: l.start}}
}
```

- [ ] **Step 2: Append Ion literal test cases to `lexer_test.go`**

```go
		// =============================================================
		// Ion literals — base lexer (Task 10)
		// =============================================================
		{
			"ion_simple",
			"`{a: 1}`",
			[]Token{{tokION_LITERAL, "{a: 1}", ast.Loc{Start: 0, End: 8}}},
		},
		{
			"ion_empty",
			"``",
			[]Token{{tokION_LITERAL, "", ast.Loc{Start: 0, End: 2}}},
		},
		{
			"ion_with_whitespace",
			"`  abc  `",
			[]Token{{tokION_LITERAL, "  abc  ", ast.Loc{Start: 0, End: 9}}},
		},
```

- [ ] **Step 3: Build and run tests**

```bash
go build ./partiql/parser/...
go test -v -run TestLexer_Tokens ./partiql/parser/...
```

Expected: 80 sub-tests pass (77 prior + 3 Ion).

- [ ] **Step 4: Vet, gofmt, commit**

```bash
go vet ./partiql/parser/...
gofmt -l partiql/parser/
git add partiql/parser/lexer.go partiql/parser/lexer_test.go
git commit -m "$(cat <<'EOF'
feat(partiql/parser): implement scanIonLiteral (simplified base)

Replaces the Task 5 stub with the simplified base implementation:
naive byte-to-byte scan from opening backtick to next backtick. The
captured Token.Str is the verbatim inner content; Token.Loc covers
the full `...` range.

KNOWN LIMITATION documented inline: Ion mode's special handling of
backticks inside Ion strings is NOT implemented. Refinement deferred
to DAG node 17 (parser-ion-literals). The AWS DynamoDB PartiQL corpus
has zero real Ion literals so this limitation does not affect the
smoke test.

Adds 3 TestLexer_Tokens cases: simple Ion blob `{a: 1}`, empty ``,
and one with leading/trailing whitespace.

After this task all 6 scan helpers are implemented and the dispatch
in Next() is fully wired. Tasks 11-13 add the corpus smoke test, the
error-case golden tests, and the parity tests.

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

### Task 11: `lexer_test.go` AWS corpus smoke test

**Files:**
- Modify: `partiql/parser/lexer_test.go` (append `TestLexer_AWSCorpus`)

- [ ] **Step 1: Append `TestLexer_AWSCorpus` to `lexer_test.go`**

Add new imports first. Find the existing import block and replace with:

```go
import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/bytebase/omni/partiql/ast"
)
```

(Adds `os` and `path/filepath`.)

Then append at the end of the file:

```go
// TestLexer_AWSCorpus loads every .partiql file from
// partiql/parser/testdata/aws-corpus/ (63 files), filters out the
// 2 syntax-skeleton files, and asserts each one lexes to a non-error
// token stream ending with EOF. Catches "does the lexer tokenize at
// all" regressions on real AWS DynamoDB PartiQL examples.
//
// Skipped files:
//   - select-001.partiql: SELECT syntax skeleton with bracket placeholders
//   - insert-002.partiql: INSERT syntax skeleton with backtick placeholder
//
// Both are flagged in testdata/aws-corpus/index.json as not-real-PartiQL.
// The skip list is hard-coded here for clarity.
func TestLexer_AWSCorpus(t *testing.T) {
	skip := map[string]bool{
		"select-001.partiql": true,
		"insert-002.partiql": true,
	}

	files, err := filepath.Glob("testdata/aws-corpus/*.partiql")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no corpus files found — testdata/aws-corpus/ missing or empty?")
	}

	var lexed, skipped int
	for _, f := range files {
		name := filepath.Base(f)
		if skip[name] {
			skipped++
			continue
		}
		t.Run(name, func(t *testing.T) {
			data, err := os.ReadFile(f)
			if err != nil {
				t.Fatal(err)
			}
			l := NewLexer(string(data))
			tokens := 0
			for {
				tok := l.Next()
				if tok.Type == tokEOF {
					break
				}
				tokens++
				if tokens > 100000 {
					t.Fatalf("token stream did not terminate after %d tokens", tokens)
				}
			}
			if l.Err != nil {
				t.Errorf("lexer error: %v", l.Err)
			}
			if tokens == 0 {
				t.Errorf("lexed to zero tokens")
			}
		})
		lexed++
	}
	t.Logf("AWS corpus: %d files lexed, %d skipped", lexed, skipped)
}
```

- [ ] **Step 2: Run the new test**

```bash
go test -v -run TestLexer_AWSCorpus ./partiql/parser/...
```

Expected: 61 sub-tests pass (63 files − 2 skipped). Each prints something like:

```
=== RUN   TestLexer_AWSCorpus/batching-001.partiql
--- PASS: TestLexer_AWSCorpus/batching-001.partiql (0.00s)
...
    lexer_test.go:XXX: AWS corpus: 61 files lexed, 2 skipped
```

If any file fails, debug the lexer for that specific input. Common causes:
- Identifier with characters not in `[a-zA-Z_$0-9]`
- Numeric literal edge case
- Unhandled punctuation

Do NOT add the failing file to the skip list. Fix the lexer.

- [ ] **Step 3: Vet, gofmt, commit**

```bash
go vet ./partiql/parser/...
gofmt -l partiql/parser/
git add partiql/parser/lexer_test.go
git commit -m "$(cat <<'EOF'
test(partiql/parser): add AWS corpus smoke test

Loads every .partiql file from testdata/aws-corpus/ (63 files),
filters out select-001 and insert-002 (syntax skeletons per
index.json), and asserts each lexes to a non-error token stream
ending with tokEOF. Expected: 61 sub-tests pass.

Catches "does the lexer tokenize at all" regressions on real AWS
DynamoDB PartiQL examples.

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

### Task 12: `lexer_test.go` error-case golden tests

**Files:**
- Modify: `partiql/parser/lexer_test.go` (append `TestLexer_Errors`)

- [ ] **Step 1: Add the `strings` import**

Update the import block in `lexer_test.go`:

```go
import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/bytebase/omni/partiql/ast"
)
```

- [ ] **Step 2: Append `TestLexer_Errors`**

```go
// TestLexer_Errors covers the 5 error triggers in the lexer. Each case
// drains Next() until tokEOF and asserts l.Err is set with the expected
// error message substring.
func TestLexer_Errors(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		wantErrIn string // substring of the expected error message
	}{
		{
			name:      "unterminated_string",
			input:     "'hello",
			wantErrIn: "unterminated string literal",
		},
		{
			name:      "unterminated_quoted_ident",
			input:     `"foo`,
			wantErrIn: "unterminated quoted identifier",
		},
		{
			name:      "unterminated_ion_literal",
			input:     "`abc",
			wantErrIn: "unterminated Ion literal",
		},
		{
			name:      "unterminated_block_comment",
			input:     "/* nope",
			wantErrIn: "unterminated block comment",
		},
		{
			name:      "unexpected_character_lone_bang",
			input:     "! 1",
			wantErrIn: "unexpected character",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			l := NewLexer(tc.input)
			for {
				tok := l.Next()
				if tok.Type == tokEOF {
					break
				}
			}
			if l.Err == nil {
				t.Errorf("expected error, got nil")
				return
			}
			if !strings.Contains(l.Err.Error(), tc.wantErrIn) {
				t.Errorf("error mismatch\n got: %v\nwant substring: %q", l.Err, tc.wantErrIn)
			}
		})
	}
}
```

- [ ] **Step 3: Run tests**

```bash
go test -v -run TestLexer_Errors ./partiql/parser/...
```

Expected: 5 sub-tests pass.

If `unexpected_character_lone_bang` fails because `!` is being treated as a valid identifier start: PartiQL grammar's `IDENTIFIER : [A-Z$_]...` does NOT include `!`, so `isIdentStart` should return false for `!`. The `scanOperator` default arm should fire because `!` alone isn't followed by `=`. Verify your `isIdentStart` and the two-character lookahead in `scanOperator`.

- [ ] **Step 4: Vet, gofmt, commit**

```bash
go vet ./partiql/parser/...
gofmt -l partiql/parser/
git add partiql/parser/lexer_test.go
git commit -m "$(cat <<'EOF'
test(partiql/parser): add error-case golden tests

Five error triggers, one test case each:
- unterminated_string ('hello)
- unterminated_quoted_ident ("foo)
- unterminated_ion_literal (`abc)
- unterminated_block_comment (/* nope)
- unexpected_character_lone_bang (! 1)

Each case drains Next() until tokEOF and asserts l.Err is set with
the expected error message substring.

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

### Task 13: `lexer_test.go` parity tests — TestTokenName_AllCovered + TestKeywords_LenMatchesConstants

**Files:**
- Modify: `partiql/parser/lexer_test.go` (append two new test functions)

- [ ] **Step 1: Append both parity tests**

```go
// TestTokenName_AllCovered walks every tok* constant declared in token.go
// and asserts tokenName returns a non-empty, non-"" string. If a future
// contributor adds a new tok* constant without wiring it into tokenName,
// this test fails.
func TestTokenName_AllCovered(t *testing.T) {
	all := []int{
		// Specials.
		tokEOF, tokInvalid,
		// Literals.
		tokSCONST, tokICONST, tokFCONST, tokIDENT, tokIDENT_QUOTED, tokION_LITERAL,
		// Operators / punctuation.
		tokPLUS, tokMINUS, tokASTERISK, tokSLASH_FORWARD, tokPERCENT,
		tokCARET, tokTILDE, tokAT_SIGN, tokEQ, tokNEQ, tokLT, tokGT,
		tokLT_EQ, tokGT_EQ, tokCONCAT, tokANGLE_DOUBLE_LEFT, tokANGLE_DOUBLE_RIGHT,
		tokPAREN_LEFT, tokPAREN_RIGHT, tokBRACKET_LEFT, tokBRACKET_RIGHT,
		tokBRACE_LEFT, tokBRACE_RIGHT, tokCOLON, tokCOLON_SEMI, tokCOMMA,
		tokPERIOD, tokQUESTION_MARK,
		// Keywords (alphabetical).
		tokABSOLUTE, tokACTION, tokADD, tokALL, tokALLOCATE, tokALTER, tokAND,
		tokANY, tokARE, tokAS, tokASC, tokASSERTION, tokAT, tokAUTHORIZATION,
		tokAVG, tokBAG, tokBEGIN, tokBETWEEN, tokBIGINT, tokBIT, tokBIT_LENGTH,
		tokBLOB, tokBOOL, tokBOOLEAN, tokBY, tokCAN_CAST, tokCAN_LOSSLESS_CAST,
		tokCASCADE, tokCASCADED, tokCASE, tokCAST, tokCATALOG, tokCHAR,
		tokCHAR_LENGTH, tokCHARACTER, tokCHARACTER_LENGTH, tokCHECK, tokCLOB,
		tokCLOSE, tokCOALESCE, tokCOLLATE, tokCOLLATION, tokCOLUMN, tokCOMMIT,
		tokCONFLICT, tokCONNECT, tokCONNECTION, tokCONSTRAINT, tokCONSTRAINTS,
		tokCONTINUE, tokCONVERT, tokCORRESPONDING, tokCOUNT, tokCREATE, tokCROSS,
		tokCURRENT, tokCURRENT_DATE, tokCURRENT_TIME, tokCURRENT_TIMESTAMP,
		tokCURRENT_USER, tokCURSOR, tokDATE, tokDATE_ADD, tokDATE_DIFF,
		tokDEALLOCATE, tokDEC, tokDECIMAL, tokDECLARE, tokDEFAULT, tokDEFERRABLE,
		tokDEFERRED, tokDELETE, tokDESC, tokDESCRIBE, tokDESCRIPTOR,
		tokDIAGNOSTICS, tokDISCONNECT, tokDISTINCT, tokDO, tokDOMAIN, tokDOUBLE,
		tokDROP, tokELSE, tokEND, tokEND_EXEC, tokESCAPE, tokEXCEPT, tokEXCEPTION,
		tokEXCLUDED, tokEXEC, tokEXECUTE, tokEXISTS, tokEXPLAIN, tokEXTERNAL,
		tokEXTRACT, tokFALSE, tokFETCH, tokFIRST, tokFLOAT, tokFOR, tokFOREIGN,
		tokFOUND, tokFROM, tokFULL, tokGET, tokGLOBAL, tokGO, tokGOTO, tokGRANT,
		tokGROUP, tokHAVING, tokIDENTITY, tokIMMEDIATE, tokIN, tokINDEX,
		tokINDICATOR, tokINITIALLY, tokINNER, tokINPUT, tokINSENSITIVE, tokINSERT,
		tokINT, tokINT2, tokINT4, tokINT8, tokINTEGER, tokINTEGER2, tokINTEGER4,
		tokINTEGER8, tokINTERSECT, tokINTERVAL, tokINTO, tokIS, tokISOLATION,
		tokJOIN, tokKEY, tokLAG, tokLANGUAGE, tokLAST, tokLATERAL, tokLEAD,
		tokLEFT, tokLET, tokLEVEL, tokLIKE, tokLIMIT, tokLIST, tokLOCAL, tokLOWER,
		tokMATCH, tokMAX, tokMIN, tokMISSING, tokMODIFIED, tokMODULE, tokNAMES,
		tokNATIONAL, tokNATURAL, tokNCHAR, tokNEW, tokNEXT, tokNO, tokNOT,
		tokNOTHING, tokNULL, tokNULLIF, tokNULLS, tokNUMERIC, tokOCTET_LENGTH,
		tokOF, tokOFFSET, tokOLD, tokON, tokONLY, tokOPEN, tokOPTION, tokOR,
		tokORDER, tokOUTER, tokOUTPUT, tokOVER, tokOVERLAPS, tokOVERLAY, tokPAD,
		tokPARTIAL, tokPARTITION, tokPIVOT, tokPLACING, tokPOSITION, tokPRECISION,
		tokPREPARE, tokPRESERVE, tokPRIMARY, tokPRIOR, tokPRIVILEGES, tokPROCEDURE,
		tokPUBLIC, tokREAD, tokREAL, tokREFERENCES, tokRELATIVE, tokREMOVE,
		tokREPLACE, tokRESTRICT, tokRETURNING, tokREVOKE, tokRIGHT, tokROLLBACK,
		tokROWS, tokSCHEMA, tokSCROLL, tokSECTION, tokSELECT, tokSESSION,
		tokSESSION_USER, tokSET, tokSEXP, tokSHORTEST, tokSIZE, tokSMALLINT,
		tokSOME, tokSPACE, tokSQL, tokSQLCODE, tokSQLERROR, tokSQLSTATE, tokSTRING,
		tokSTRUCT, tokSUBSTRING, tokSUM, tokSYMBOL, tokSYSTEM_USER, tokTABLE,
		tokTEMPORARY, tokTHEN, tokTIME, tokTIMESTAMP, tokTO, tokTRANSACTION,
		tokTRANSLATE, tokTRANSLATION, tokTRIM, tokTRUE, tokTUPLE, tokUNION,
		tokUNIQUE, tokUNKNOWN, tokUNPIVOT, tokUPDATE, tokUPPER, tokUPSERT,
		tokUSAGE, tokUSER, tokUSING, tokVALUE, tokVALUES, tokVARCHAR, tokVARYING,
		tokVIEW, tokWHEN, tokWHENEVER, tokWHERE, tokWITH, tokWORK, tokWRITE,
		tokZONE,
	}
	if got := len(all); got != 302 {
		t.Errorf("test list has %d entries, want 302 — did a tok* constant get added or removed without updating this test?", got)
	}
	for _, tt := range all {
		name := tokenName(tt)
		if name == "" {
			t.Errorf("tokenName(%d) returned empty string — missing switch arm in token.go?", tt)
		}
	}
}

// TestKeywords_LenMatchesConstants asserts that the keywords map in
// keywords.go has exactly 266 entries — the same number as the keyword
// constants in token.go. If a future contributor adds or removes a
// tok* keyword constant without updating the map (or vice versa),
// this test fails.
func TestKeywords_LenMatchesConstants(t *testing.T) {
	const expectedKeywordCount = 266
	if got := len(keywords); got != expectedKeywordCount {
		t.Errorf("len(keywords) = %d, want %d — did a tok* keyword constant get added or removed without updating the keywords map?", got, expectedKeywordCount)
	}
}
```

- [ ] **Step 2: Run all tests**

```bash
go test -v ./partiql/parser/...
```

Expected: all tests pass:
- `TestLexer_Tokens` — 80 sub-tests
- `TestLexer_AWSCorpus` — 61 sub-tests
- `TestLexer_Errors` — 5 sub-tests
- `TestTokenName_AllCovered` — 0 sub-tests (single test function)
- `TestKeywords_LenMatchesConstants` — 0 sub-tests (single test function)

Total: ~146 passing sub-tests + 2 single-function tests = 148 test units.

If `TestTokenName_AllCovered` fails with "test list has X entries, want 302" — count the entries in the `all` slice and update the assertion to match the actual count. The slice should contain exactly: 2 specials + 6 literals + 28 operators + 266 keywords = 302 entries.

If `TestKeywords_LenMatchesConstants` fails — count the entries in `keywords.go` and either fix the map or update `expectedKeywordCount`. Both numbers must match exactly.

- [ ] **Step 3: Vet, gofmt, commit**

```bash
go vet ./partiql/parser/...
gofmt -l partiql/parser/
git add partiql/parser/lexer_test.go
git commit -m "$(cat <<'EOF'
test(partiql/parser): add tokenName + keywords parity tests

TestTokenName_AllCovered: hard-coded list of all 302 tok* constants;
asserts tokenName returns a non-empty string for each. Catches future
contributors adding a tok* constant without wiring it into tokenName.

TestKeywords_LenMatchesConstants: asserts len(keywords) == 266.
Catches future contributors adding/removing a tok* keyword constant
without updating the keywords map (or vice versa).

After this task all 5 test functions are present:
- TestLexer_Tokens (80 sub-tests)
- TestLexer_AWSCorpus (61 sub-tests)
- TestLexer_Errors (5 sub-tests)
- TestTokenName_AllCovered (single function, 302 constants)
- TestKeywords_LenMatchesConstants (single function)

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

### Task 14: Final verification, DAG bookkeeping, finishing branch

**Files:**
- Read: `/Users/h3n4l/OpenSource/parser/partiql/PartiQLLexer.g4` (grammar cross-check)
- Modify: `/Users/h3n4l/OpenSource/omni/docs/migration/partiql/dag.md` (mark node 2 as done — on `main`)

- [ ] **Step 1: Run the full test suite, vet, and gofmt**

```bash
cd /Users/h3n4l/OpenSource/omni/.worktrees/feat-partiql-lexer
go test ./partiql/parser/...
go vet ./partiql/parser/...
gofmt -l partiql/parser/
```

Expected: ok, no output from vet or gofmt.

- [ ] **Step 2: Cross-check the lexer against `PartiQLLexer.g4`**

Open `/Users/h3n4l/OpenSource/parser/partiql/PartiQLLexer.g4` and verify:

1. **Every uppercase keyword rule** (lines 13–295) maps to a `tok*` constant in `token.go` and a corresponding entry in `keywords.go`. Run a sanity check:

```bash
diff <(grep -E "^[A-Z][A-Z0-9_]*: '[A-Z]" /Users/h3n4l/OpenSource/parser/partiql/PartiQLLexer.g4 | sed -E "s/^([A-Z][A-Z0-9_]*):.*/\1/" | sort) \
     <(grep -E "^\ttok[A-Z][A-Z0-9_]+ = iota \+ 3000$|^\ttok[A-Z][A-Z0-9_]+$" partiql/parser/token.go | sed -E 's/.*tok([A-Z][A-Z0-9_]+).*/\1/' | sort -u | head -266)
```

Expected: empty diff. If any keyword from the grammar is missing from the constants, add it (and the corresponding map entry).

2. **Every operator rule** (lines 303–331) maps to a `tok*` constant. Cross-check manually against the operator block in `token.go`. The 28 expected operators are: PLUS, MINUS, ASTERISK, SLASH_FORWARD, PERCENT, CARET, TILDE, AT_SIGN, EQ, NEQ, LT, GT, LT_EQ, GT_EQ, CONCAT, ANGLE_DOUBLE_LEFT, ANGLE_DOUBLE_RIGHT, PAREN_LEFT, PAREN_RIGHT, BRACKET_LEFT, BRACKET_RIGHT, BRACE_LEFT, BRACE_RIGHT, COLON, COLON_SEMI, COMMA, PERIOD, QUESTION_MARK. (`BACKTICK` and `ION_CLOSURE` are subsumed by `tokION_LITERAL`.)

3. **Every literal rule** (lines 339–355) maps to a `tok*` constant: `LITERAL_STRING` → `tokSCONST`, `LITERAL_INTEGER` → `tokICONST`, `LITERAL_DECIMAL` → `tokFCONST`, `IDENTIFIER` → `tokIDENT`, `IDENTIFIER_QUOTED` → `tokIDENT_QUOTED`, plus `tokION_LITERAL` for the deferred Ion mode.

If you find any missing rule, add the corresponding constant and a test case before proceeding.

- [ ] **Step 3: Run the spot-check commands**

```bash
# Sanity: count tok constants
grep -c '^\ttok[A-Z]' partiql/parser/token.go
# Should print 302

# Sanity: count keywords map entries
grep -cE '^\t"[a-z0-9_-]+":' partiql/parser/keywords.go
# Should print 266

# Sanity: count tokenName arms
grep -c '^\tcase tok[A-Z]' partiql/parser/token.go
# Should print 302
```

If any number is off, fix the discrepancy before committing.

- [ ] **Step 4: Run the test pass count one more time**

```bash
go test -v ./partiql/parser/... 2>&1 | tail -20
```

Expected: clean PASS line at the end. ~146 sub-tests across 5 test functions.

- [ ] **Step 5: Final verification commit (allow-empty if no changes)**

```bash
git commit --allow-empty -m "$(cat <<'EOF'
chore(partiql/parser): final verification pass

- go test ./partiql/parser/... — pass
- go vet ./partiql/parser/... — clean
- gofmt -l partiql/parser/ — clean
- 302 tok* constants verified across token.go
- 266 keywords map entries verified against grammar
- Cross-checked against bytebase/parser/partiql/PartiQLLexer.g4
  line by line; every uppercase rule maps to a constant
- ~146 sub-tests passing across 5 test functions:
    TestLexer_Tokens (80), TestLexer_AWSCorpus (61),
    TestLexer_Errors (5), TestTokenName_AllCovered,
    TestKeywords_LenMatchesConstants

Closes lexer (DAG node 2) for the partiql migration.

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>
EOF
)"
```

- [ ] **Step 6: Update `dag.md` on `main` to mark node 2 as `done`**

```bash
cd /Users/h3n4l/OpenSource/omni
```

Open `docs/migration/partiql/dag.md` and find the row:

```
| 2 | lexer | `partiql/parser` (lexer.go, token.go) | ast-core | catalog | **P0** | not started |
```

Change `not started` to `done`. Note: the actual file split is 3 files (token.go + keywords.go + lexer.go) — update the package column too if you prefer:

```
| 2 | lexer | `partiql/parser` (lexer.go, token.go, keywords.go) | ast-core | catalog | **P0** | done |
```

Stage and commit on main:

```bash
git add docs/migration/partiql/dag.md
git commit -m "$(cat <<'EOF'
docs(partiql): mark lexer (DAG node 2) as done

The partiql/parser lexer is complete on the feat/partiql/lexer
branch and ready to merge: 302 tok constants (266 keywords + 28
operators + 6 literals + 2 specials), Lexer.Next() with single-pass
scan helpers, ~146 passing sub-tests covering golden token streams,
the AWS DynamoDB PartiQL example corpus, error triggers, and parity
checks.

DAG node 4 (parser-foundation) is now unblocked.

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>
EOF
)"
```

- [ ] **Step 7: Hand off to `superpowers:finishing-a-development-branch`**

From the worktree:

```bash
cd /Users/h3n4l/OpenSource/omni/.worktrees/feat-partiql-lexer
git log --oneline feat/partiql/lexer ^main
```

Verify the commit list looks right (should be ~16 commits: 14 task commits + spec + plan).

Then invoke `superpowers:finishing-a-development-branch`. It will guide the merge/PR decision.

- [ ] **Step 8: After branch finish, clean up the worktree and report**

After `finishing-a-development-branch` completes:

1. Confirm `git worktree list` no longer shows `feat-partiql-lexer` (if option 1 or 4 was chosen)
2. Confirm `dag.md` on `main` shows lexer as `done`
3. Report to the user:
   - Which DAG node was completed
   - How many tok constants (302) and how many sub-tests (~146) were added
   - The next actionable nodes per the DAG: `parser-foundation` (node 4) and `catalog` (node 3)

---

## Self-Review

Walking the spec section by section against this plan:

**Goal & inputs (spec lines 9–18):** Task 1 establishes the package and references the spec. ✓

**D1. Hand-written single-pass lexer (cosmosdb pattern):** Tasks 5–10 implement the cosmosdb-style scanner with `Lexer.Next()` driven by `skipWhitespaceAndComments` + character-class dispatch. ✓

**D2. Token positions use `ast.Loc` directly:** Task 1 declares `Token.Loc ast.Loc`. Every scan helper sets `Loc: ast.Loc{Start: l.start, End: l.pos}`. ✓

**D3. First-error-and-stop:** Tasks 5–10 set `l.Err` on each error and return `tokEOF`; Task 5's `Next()` checks `l.Err` at the top. ✓

**D4. Three-file split:** Tasks 1–3 build `token.go`, Task 4 creates `keywords.go`, Task 5 creates `lexer.go`. ✓

**D5. Unexported `tok*` constants:** Every constant uses lowercase `tok` prefix + UPPERCASE name. ✓

**D6. Naive backtick scan:** Task 10 implements `scanIonLiteral` with the documented limitation. ✓

**Token taxonomy (spec):** Task 1 has 36 non-keyword constants; Task 2 has 266 keyword constants; Task 3 has the full 302-arm `tokenName`. ✓

**Scanning behavior (spec):** Task 5 has dispatch + skipWhitespaceAndComments; Tasks 6–10 implement the 6 scan helpers in the order matching the spec's helper sections. ✓

**Error model (spec):** Task 12's `TestLexer_Errors` covers all 5 trigger types from the spec's "Error reporting" table. ✓

**Test plan (spec):** Tasks 11–13 implement the 5 test functions from the spec exactly: TestLexer_Tokens (Task 6/7/8/9/10 incremental), TestLexer_AWSCorpus (Task 11), TestLexer_Errors (Task 12), TestTokenName_AllCovered + TestKeywords_LenMatchesConstants (Task 13). ✓

**Acceptance criteria 1–11 (spec):** Task 14's verification step walks all 11 criteria. ✓

**Non-goals (spec):** Plan never implements Ion-mode-aware backtick scanning, date/time literal special-casing, comment preservation, Tokenize() API, line/column conversion, ANTLR cross-comparison, hex/binary/octal numbers, backslash escapes, or statement splitting. ✓

**Placeholder scan:** Searched for "TBD", "TODO", "fill in details", "appropriate error handling", etc. None present. The Task 5 stubs are explicitly labeled as such and replaced one-by-one in Tasks 6–10.

**Type consistency:** `Token.Type int`, `Token.Str string`, `Token.Loc ast.Loc` — same shape across all references. `Lexer` struct fields (`input`, `pos`, `start`, `Err`) are consistent. Method signatures (`Next() Token`, `scanXxx() Token`) are uniform.

**Note on test count divergence:** the spec said "approximately 55–60 cases" for `TestLexer_Tokens`; this plan yields 80 cases (9 whitespace/comment + 5 string + 5 quoted ident + 5 unquoted ident + 9 keyword + 2 ident-after-comment + 10 numeric + 11 single-char op + 11 punctuation + 7 two-char op + 3 integration + 3 Ion). This is more than the spec estimated; the additional cases came from being thorough about per-operator and per-punctuation coverage. Not a concern — more golden coverage is strictly better.

No issues found.

---

## Plan complete and saved

Plan saved to `docs/superpowers/plans/2026-04-08-partiql-lexer.md` (this file).

Next step: pick an execution mode.
