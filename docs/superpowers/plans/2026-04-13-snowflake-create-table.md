# T2.2: Snowflake CREATE TABLE Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement CREATE TABLE parsing for the Snowflake engine, covering standard column/constraint definitions, CTAS, LIKE, CLONE, and all table properties required by bytebase lint rules.

**Architecture:** A new `create_table.go` file in `snowflake/parser/` houses all parse functions. The existing `parseStmt` dispatch in `parser.go` is modified to call `parseCreateStmt()` instead of `unsupported("CREATE")`. Three new AST node types (`CreateTableStmt`, `ColumnDef`, `TableConstraint`) are added to `snowflake/ast/parsenodes.go` with corresponding tags and walker support.

**Tech Stack:** Go, recursive-descent parser (hand-written), Pratt expression parser (existing), genwalker code generator (existing).

---

## File Map

| File | Action | Responsibility |
|------|--------|----------------|
| `snowflake/ast/parsenodes.go` | Modify | Add CreateTableStmt, ColumnDef, TableConstraint nodes + helper structs + enums |
| `snowflake/ast/nodetags.go` | Modify | Add T_CreateTableStmt, T_ColumnDef, T_TableConstraint tags + String cases |
| `snowflake/ast/walk_generated.go` | Regenerate | Auto-generated walker cases for 3 new node types |
| `snowflake/parser/parser.go:185-186` | Modify | Replace `unsupported("CREATE")` with `parseCreateStmt()` |
| `snowflake/parser/create_table.go` | Create | All CREATE TABLE parse functions |
| `snowflake/parser/create_table_test.go` | Create | All CREATE TABLE tests |

---

### Task 1: AST Types + Node Tags

**Files:**
- Modify: `snowflake/ast/parsenodes.go` (append after line 861, the end of SetOperationStmt)
- Modify: `snowflake/ast/nodetags.go` (add 3 tags after T_JoinExpr on line 52, add 3 String cases)

- [ ] **Step 1: Add enums and helper structs to parsenodes.go**

Append after the `var _ Node = (*SetOperationStmt)(nil)` line (line 860):

```go
// ---------------------------------------------------------------------------
// Constraint enums
// ---------------------------------------------------------------------------

// ConstraintType enumerates constraint kinds for inline and table-level constraints.
type ConstraintType int

const (
	ConstrPrimaryKey ConstraintType = iota
	ConstrForeignKey
	ConstrUnique
)

// String returns the constraint type name.
func (c ConstraintType) String() string {
	switch c {
	case ConstrPrimaryKey:
		return "PRIMARY KEY"
	case ConstrForeignKey:
		return "FOREIGN KEY"
	case ConstrUnique:
		return "UNIQUE"
	default:
		return "UNKNOWN"
	}
}

// ReferenceAction enumerates FK referential actions.
type ReferenceAction int

const (
	RefActNone       ReferenceAction = iota // not specified
	RefActCascade                           // CASCADE
	RefActSetNull                           // SET NULL
	RefActSetDefault                        // SET DEFAULT
	RefActRestrict                          // RESTRICT
	RefActNoAction                          // NO ACTION
)

// ---------------------------------------------------------------------------
// Helper structs (not Nodes)
// ---------------------------------------------------------------------------

// InlineConstraint represents a column-level constraint.
type InlineConstraint struct {
	Type       ConstraintType
	Name       Ident          // CONSTRAINT name; zero if unnamed
	References *ForeignKeyRef // for FK; nil otherwise
	Loc        Loc
}

// ForeignKeyRef holds REFERENCES clause details.
type ForeignKeyRef struct {
	Table    *ObjectName
	Columns  []Ident
	OnDelete ReferenceAction
	OnUpdate ReferenceAction
	Match    string // "FULL"/"PARTIAL"/"SIMPLE"; empty if absent
}

// IdentitySpec holds IDENTITY/AUTOINCREMENT configuration.
type IdentitySpec struct {
	Start     *int64 // START WITH value; nil if default
	Increment *int64 // INCREMENT BY value; nil if default
	Order     *bool  // true=ORDER, false=NOORDER, nil=unspecified
}

// TagAssignment is a single TAG name = 'value' pair.
type TagAssignment struct {
	Name  *ObjectName
	Value string
}

// CloneSource holds CLONE source with optional time travel.
type CloneSource struct {
	Source   *ObjectName
	AtBefore string // "AT" or "BEFORE"; empty if no time travel
	Kind     string // "TIMESTAMP"/"OFFSET"/"STATEMENT"
	Value    string // the time travel value
}
```

- [ ] **Step 2: Add Node types to parsenodes.go**

Append immediately after the helper structs:

```go
// ---------------------------------------------------------------------------
// DDL statement nodes
// ---------------------------------------------------------------------------

// CreateTableStmt represents CREATE [OR REPLACE] [TRANSIENT|TEMPORARY|VOLATILE] TABLE ...
type CreateTableStmt struct {
	OrReplace   bool
	Transient   bool
	Temporary   bool
	Volatile    bool
	IfNotExists bool
	Name        *ObjectName
	Columns     []*ColumnDef
	Constraints []*TableConstraint
	ClusterBy   []Node           // CLUSTER BY expressions; nil if absent
	Linear      bool             // CLUSTER BY LINEAR modifier
	Comment     *string          // COMMENT = 'text'; nil if absent
	CopyGrants  bool
	Tags        []*TagAssignment // WITH TAG (...); nil if absent
	AsSelect    Node             // CREATE TABLE ... AS SELECT; nil if absent
	Like        *ObjectName      // CREATE TABLE ... LIKE source; nil if absent
	Clone       *CloneSource     // CREATE TABLE ... CLONE source; nil if absent
	Loc         Loc
}

func (n *CreateTableStmt) Tag() NodeTag { return T_CreateTableStmt }

// ColumnDef represents a column definition in CREATE TABLE.
type ColumnDef struct {
	Name             Ident
	DataType         *TypeName         // nil for virtual columns without explicit type
	Default          Node              // DEFAULT expr; nil if absent
	NotNull          bool
	Nullable         bool              // explicit NULL
	Identity         *IdentitySpec     // IDENTITY/AUTOINCREMENT; nil if absent
	Collate          string            // COLLATE 'name'; empty if absent
	MaskingPolicy    *ObjectName       // WITH MASKING POLICY name; nil if absent
	InlineConstraint *InlineConstraint // inline PK/FK/UNIQUE; nil if absent
	Comment          *string           // COMMENT 'text'; nil if absent
	Tags             []*TagAssignment  // WITH TAG (...); nil if absent
	VirtualExpr      Node              // AS (expr); nil if absent
	Loc              Loc
}

func (n *ColumnDef) Tag() NodeTag { return T_ColumnDef }

// TableConstraint represents a table-level constraint (out-of-line).
type TableConstraint struct {
	Type       ConstraintType // ConstrPrimaryKey/ConstrForeignKey/ConstrUnique
	Name       Ident          // CONSTRAINT name; zero if unnamed
	Columns    []Ident        // constrained column names
	References *ForeignKeyRef // FK only; nil otherwise
	Comment    *string        // inline COMMENT 'text'; nil if absent
	Loc        Loc
}

func (n *TableConstraint) Tag() NodeTag { return T_TableConstraint }

// Compile-time assertions.
var (
	_ Node = (*CreateTableStmt)(nil)
	_ Node = (*ColumnDef)(nil)
	_ Node = (*TableConstraint)(nil)
)
```

- [ ] **Step 3: Add node tags to nodetags.go**

In `snowflake/ast/nodetags.go`, add 3 new tags after `T_JoinExpr` (line 52):

```go
	T_CreateTableStmt
	T_ColumnDef
	T_TableConstraint
```

Add 3 String cases in the `String()` method, before the `default:` line:

```go
	case T_CreateTableStmt:
		return "CreateTableStmt"
	case T_ColumnDef:
		return "ColumnDef"
	case T_TableConstraint:
		return "TableConstraint"
```

- [ ] **Step 4: Regenerate walker**

Run:
```bash
cd /Users/h3n4l/OpenSource/omni/.worktrees/create-table && go run ./snowflake/ast/cmd/genwalker
```

Expected: `walk_generated.go` updated with 3 new cases:
- `*CreateTableStmt`: walks Name, Columns, Constraints, ClusterBy, AsSelect, Like
- `*ColumnDef`: walks DataType, Default, MaskingPolicy, VirtualExpr
- `*TableConstraint`: no walkable Node children

- [ ] **Step 5: Verify compilation**

Run:
```bash
cd /Users/h3n4l/OpenSource/omni/.worktrees/create-table && go build ./snowflake/...
```

Expected: clean compile, no errors.

- [ ] **Step 6: Commit**

```bash
cd /Users/h3n4l/OpenSource/omni/.worktrees/create-table
git add snowflake/ast/parsenodes.go snowflake/ast/nodetags.go snowflake/ast/walk_generated.go
git commit -m "feat(snowflake): add CREATE TABLE AST types (T2.2 step 1)

Add CreateTableStmt, ColumnDef, TableConstraint node types plus
ConstraintType, ReferenceAction enums and InlineConstraint,
ForeignKeyRef, IdentitySpec, TagAssignment, CloneSource helper structs.
Regenerate walker.

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>"
```

---

### Task 2: CREATE Dispatch + Skeleton Parser

**Files:**
- Modify: `snowflake/parser/parser.go:185-186` (replace unsupported("CREATE"))
- Create: `snowflake/parser/create_table.go`

- [ ] **Step 1: Create create_table.go with parseCreateStmt and parseCreateTableStmt skeleton**

Create `snowflake/parser/create_table.go`:

```go
package parser

import "github.com/bytebase/omni/snowflake/ast"

// parseCreateStmt dispatches CREATE statements. Consumes CREATE, then
// handles OR REPLACE, table type modifiers (TRANSIENT/TEMPORARY/VOLATILE),
// and sub-dispatches on the next keyword (TABLE, or unsupported).
//
// Future DAG nodes plug in by adding cases to the inner switch.
func (p *Parser) parseCreateStmt() (ast.Node, error) {
	start := p.cur.Loc
	p.advance() // consume CREATE

	// OR REPLACE
	orReplace := false
	if p.cur.Type == kwOR && p.peekNext().Type == kwREPLACE {
		p.advance() // consume OR
		p.advance() // consume REPLACE
		orReplace = true
	}

	// Table type modifiers: [LOCAL|GLOBAL] TEMPORARY|TEMP / TRANSIENT / VOLATILE
	transient := false
	temporary := false
	volatile := false
	switch p.cur.Type {
	case kwLOCAL, kwGLOBAL:
		p.advance() // consume LOCAL/GLOBAL (no semantic difference)
		if p.cur.Type == kwTEMPORARY || p.cur.Type == kwTEMP {
			p.advance()
			temporary = true
		}
	case kwTEMPORARY, kwTEMP:
		p.advance()
		temporary = true
	case kwTRANSIENT:
		p.advance()
		transient = true
	case kwVOLATILE:
		p.advance()
		volatile = true
	}

	switch p.cur.Type {
	case kwTABLE:
		return p.parseCreateTableStmt(start, orReplace, transient, temporary, volatile)
	default:
		// All other CREATE forms remain unsupported.
		err := &ParseError{
			Loc: start,
			Msg: "CREATE statement parsing is not yet supported for this object type",
		}
		p.skipToNextStatement()
		return nil, err
	}
}

// parseCreateTableStmt parses CREATE TABLE after the table type has been consumed.
//
//	CREATE [OR REPLACE] [TRANSIENT|TEMPORARY|VOLATILE] TABLE [IF NOT EXISTS] name
//	  ( column_def | table_constraint , ... )
//	  [table_properties...]
//	  [AS query]
//
// Also handles LIKE, CLONE, and bare CTAS forms.
func (p *Parser) parseCreateTableStmt(start ast.Loc, orReplace, transient, temporary, volatile bool) (*ast.CreateTableStmt, error) {
	p.advance() // consume TABLE

	stmt := &ast.CreateTableStmt{
		OrReplace: orReplace,
		Transient: transient,
		Temporary: temporary,
		Volatile:  volatile,
		Loc:       start,
	}

	// IF NOT EXISTS
	if p.cur.Type == kwIF {
		p.advance() // consume IF
		if _, err := p.expect(kwNOT); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		stmt.IfNotExists = true
	}

	// Table name
	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Branch: LIKE / CLONE / AS / ( columns ) / ( columns ) AS
	switch p.cur.Type {
	case kwLIKE:
		p.advance() // consume LIKE
		like, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		stmt.Like = like
		// LIKE can be followed by CLUSTER BY and COPY GRANTS per legacy grammar.
		if err := p.parseTableProperties(stmt); err != nil {
			return nil, err
		}
		stmt.Loc.End = p.prev.Loc.End
		return stmt, nil

	case kwCLONE:
		clone, err := p.parseCloneSource()
		if err != nil {
			return nil, err
		}
		stmt.Clone = clone
		stmt.Loc.End = p.prev.Loc.End
		return stmt, nil

	case kwAS:
		p.advance() // consume AS
		query, err := p.parseQueryExpr()
		if err != nil {
			return nil, err
		}
		stmt.AsSelect = query
		stmt.Loc.End = p.prev.Loc.End
		return stmt, nil

	case '(':
		if err := p.parseColumnDeclItems(stmt); err != nil {
			return nil, err
		}
	}

	// Table properties (CLUSTER BY, COMMENT, COPY GRANTS, WITH TAG, etc.)
	if err := p.parseTableProperties(stmt); err != nil {
		return nil, err
	}

	// CTAS with columns: CREATE TABLE t (col INT) AS SELECT ...
	if p.cur.Type == kwAS {
		p.advance() // consume AS
		query, err := p.parseQueryExpr()
		if err != nil {
			return nil, err
		}
		stmt.AsSelect = query
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}
```

- [ ] **Step 2: Add stub parse functions**

Append to `create_table.go`:

```go
// parseColumnDeclItems parses the parenthesized list of column definitions
// and table-level constraints: ( item , item , ... )
func (p *Parser) parseColumnDeclItems(stmt *ast.CreateTableStmt) error {
	if _, err := p.expect('('); err != nil {
		return err
	}

	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		if p.isOutOfLineConstraintStart() {
			constr, err := p.parseOutOfLineConstraint()
			if err != nil {
				return err
			}
			stmt.Constraints = append(stmt.Constraints, constr)
		} else {
			col, err := p.parseColumnDef()
			if err != nil {
				return err
			}
			stmt.Columns = append(stmt.Columns, col)
		}

		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}

	if _, err := p.expect(')'); err != nil {
		return err
	}
	return nil
}

// isOutOfLineConstraintStart returns true if the current token starts
// a table-level constraint (CONSTRAINT, PRIMARY, UNIQUE, FOREIGN).
func (p *Parser) isOutOfLineConstraintStart() bool {
	switch p.cur.Type {
	case kwCONSTRAINT, kwPRIMARY, kwUNIQUE, kwFOREIGN:
		return true
	}
	return false
}

// parseColumnDef parses a single column definition:
//   col_name data_type [column_options...]
func (p *Parser) parseColumnDef() (*ast.ColumnDef, error) {
	start := p.cur.Loc
	name, err := p.parseIdent()
	if err != nil {
		return nil, err
	}

	col := &ast.ColumnDef{
		Name: name,
		Loc:  start,
	}

	// Data type — required unless this is a virtual column (AS ...) with no explicit type.
	// Peek: if next token is AS and not a data type keyword, skip type parsing.
	if p.cur.Type != kwAS {
		dt, err := p.parseDataType()
		if err != nil {
			return nil, err
		}
		col.DataType = dt
	}

	// Column options loop
	if err := p.parseColumnOptions(col); err != nil {
		return nil, err
	}

	col.Loc.End = p.prev.Loc.End
	return col, nil
}

// parseColumnOptions parses zero or more column-level options after the data type.
func (p *Parser) parseColumnOptions(col *ast.ColumnDef) error {
	for {
		switch p.cur.Type {
		case kwNOT:
			// NOT NULL
			p.advance() // consume NOT
			if _, err := p.expect(kwNULL); err != nil {
				return err
			}
			col.NotNull = true

		case kwNULL:
			p.advance()
			col.Nullable = true

		case kwDEFAULT:
			p.advance() // consume DEFAULT
			expr, err := p.parseExpr()
			if err != nil {
				return err
			}
			col.Default = expr

		case kwIDENTITY, kwAUTOINCREMENT:
			spec, err := p.parseIdentitySpec()
			if err != nil {
				return err
			}
			col.Identity = spec

		case kwCOLLATE:
			p.advance() // consume COLLATE
			tok, err := p.expect(tokString)
			if err != nil {
				return err
			}
			col.Collate = tok.Str

		case kwCONSTRAINT:
			ic, err := p.parseInlineConstraint()
			if err != nil {
				return err
			}
			col.InlineConstraint = ic

		case kwPRIMARY, kwUNIQUE:
			ic, err := p.parseInlineConstraint()
			if err != nil {
				return err
			}
			col.InlineConstraint = ic

		case kwCOMMENT:
			p.advance() // consume COMMENT
			tok, err := p.expect(tokString)
			if err != nil {
				return err
			}
			s := tok.Str
			col.Comment = &s

		case kwWITH:
			// WITH MASKING POLICY or WITH TAG
			next := p.peekNext()
			if next.Type == kwMASKING {
				p.advance() // consume WITH
				p.advance() // consume MASKING
				if _, err := p.expect(kwPOLICY); err != nil {
					return err
				}
				policyName, err := p.parseObjectName()
				if err != nil {
					return err
				}
				col.MaskingPolicy = policyName
				// Optional USING (col_list) — consume and discard
				if p.cur.Type == kwUSING {
					p.advance()
					if err := p.skipParenthesized(); err != nil {
						return err
					}
				}
			} else if next.Type == kwTAG {
				p.advance() // consume WITH
				tags, err := p.parseTagAssignments()
				if err != nil {
					return err
				}
				col.Tags = tags
			} else {
				return nil // not a column option
			}

		case kwMASKING:
			// MASKING POLICY without WITH prefix
			p.advance() // consume MASKING
			if _, err := p.expect(kwPOLICY); err != nil {
				return err
			}
			policyName, err := p.parseObjectName()
			if err != nil {
				return err
			}
			col.MaskingPolicy = policyName
			if p.cur.Type == kwUSING {
				p.advance()
				if err := p.skipParenthesized(); err != nil {
					return err
				}
			}

		case kwTAG:
			tags, err := p.parseTagAssignments()
			if err != nil {
				return err
			}
			col.Tags = tags

		case kwAS:
			// Virtual column: AS (expr) or AS expr
			p.advance() // consume AS
			if p.cur.Type == '(' {
				p.advance() // consume (
				expr, err := p.parseExpr()
				if err != nil {
					return err
				}
				if _, err := p.expect(')'); err != nil {
					return err
				}
				col.VirtualExpr = expr
			} else {
				expr, err := p.parseExpr()
				if err != nil {
					return err
				}
				col.VirtualExpr = expr
			}

		default:
			return nil // no more column options
		}
	}
}

// parseInlineConstraint parses a column-level inline constraint:
//   [CONSTRAINT name] PRIMARY KEY | UNIQUE | FOREIGN KEY REFERENCES table(cols)
func (p *Parser) parseInlineConstraint() (*ast.InlineConstraint, error) {
	start := p.cur.Loc
	ic := &ast.InlineConstraint{Loc: start}

	// Optional CONSTRAINT name
	if p.cur.Type == kwCONSTRAINT {
		p.advance() // consume CONSTRAINT
		name, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		ic.Name = name
	}

	switch p.cur.Type {
	case kwPRIMARY:
		p.advance() // consume PRIMARY
		if _, err := p.expect(kwKEY); err != nil {
			return nil, err
		}
		ic.Type = ast.ConstrPrimaryKey

	case kwUNIQUE:
		p.advance() // consume UNIQUE
		ic.Type = ast.ConstrUnique

	case kwFOREIGN:
		p.advance() // consume FOREIGN
		if _, err := p.expect(kwKEY); err != nil {
			return nil, err
		}
		ic.Type = ast.ConstrForeignKey
		ref, err := p.parseForeignKeyRef()
		if err != nil {
			return nil, err
		}
		ic.References = ref

	default:
		return nil, p.syntaxErrorAtCur()
	}

	// Consume optional constraint properties (ENFORCED, DEFERRABLE, etc.)
	p.parseConstraintProperties()

	ic.Loc.End = p.prev.Loc.End
	return ic, nil
}

// parseOutOfLineConstraint parses a table-level (out-of-line) constraint:
//   [CONSTRAINT name] PRIMARY KEY (cols) | UNIQUE (cols) | FOREIGN KEY (cols) REFERENCES ...
func (p *Parser) parseOutOfLineConstraint() (*ast.TableConstraint, error) {
	start := p.cur.Loc
	tc := &ast.TableConstraint{Loc: start}

	// Optional CONSTRAINT name
	if p.cur.Type == kwCONSTRAINT {
		p.advance() // consume CONSTRAINT
		name, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		tc.Name = name
	}

	switch p.cur.Type {
	case kwPRIMARY:
		p.advance() // consume PRIMARY
		if _, err := p.expect(kwKEY); err != nil {
			return nil, err
		}
		tc.Type = ast.ConstrPrimaryKey
		cols, err := p.parseIdentListInParens()
		if err != nil {
			return nil, err
		}
		tc.Columns = cols

	case kwUNIQUE:
		p.advance() // consume UNIQUE
		tc.Type = ast.ConstrUnique
		cols, err := p.parseIdentListInParens()
		if err != nil {
			return nil, err
		}
		tc.Columns = cols

	case kwFOREIGN:
		p.advance() // consume FOREIGN
		if _, err := p.expect(kwKEY); err != nil {
			return nil, err
		}
		tc.Type = ast.ConstrForeignKey
		cols, err := p.parseIdentListInParens()
		if err != nil {
			return nil, err
		}
		tc.Columns = cols

		if _, err := p.expect(kwREFERENCES); err != nil {
			return nil, err
		}
		ref, err := p.parseForeignKeyRefAfterReferences()
		if err != nil {
			return nil, err
		}
		tc.References = ref

	default:
		return nil, p.syntaxErrorAtCur()
	}

	// Consume optional constraint properties
	p.parseConstraintProperties()

	// Optional inline COMMENT
	if p.cur.Type == kwCOMMENT {
		p.advance()
		tok, err := p.expect(tokString)
		if err != nil {
			return nil, err
		}
		s := tok.Str
		tc.Comment = &s
	}

	tc.Loc.End = p.prev.Loc.End
	return tc, nil
}

// parseIdentListInParens parses ( ident, ident, ... ).
func (p *Parser) parseIdentListInParens() ([]ast.Ident, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	var idents []ast.Ident
	for {
		id, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		idents = append(idents, id)
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return idents, nil
}

// parseForeignKeyRef parses REFERENCES table [(cols)] [MATCH ...] [ON DELETE/UPDATE ...]
// Called from inline constraints where REFERENCES keyword has not yet been consumed.
func (p *Parser) parseForeignKeyRef() (*ast.ForeignKeyRef, error) {
	if _, err := p.expect(kwREFERENCES); err != nil {
		return nil, err
	}
	return p.parseForeignKeyRefAfterReferences()
}

// parseForeignKeyRefAfterReferences parses the part after REFERENCES:
//   table [(cols)] [MATCH ...] [ON DELETE ...] [ON UPDATE ...]
func (p *Parser) parseForeignKeyRefAfterReferences() (*ast.ForeignKeyRef, error) {
	table, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	ref := &ast.ForeignKeyRef{Table: table}

	// Optional column list
	if p.cur.Type == '(' {
		cols, err := p.parseIdentListInParens()
		if err != nil {
			return nil, err
		}
		ref.Columns = cols
	}

	// Optional MATCH FULL|PARTIAL|SIMPLE
	if p.cur.Type == kwMATCH {
		p.advance() // consume MATCH
		switch p.cur.Type {
		case kwFULL:
			ref.Match = "FULL"
			p.advance()
		case kwPARTIAL:
			ref.Match = "PARTIAL"
			p.advance()
		case kwSIMPLE:
			ref.Match = "SIMPLE"
			p.advance()
		default:
			return nil, p.syntaxErrorAtCur()
		}
	}

	// Optional ON DELETE / ON UPDATE (either order)
	for p.cur.Type == kwON {
		p.advance() // consume ON
		switch p.cur.Type {
		case kwDELETE:
			p.advance() // consume DELETE
			act, err := p.parseReferenceAction()
			if err != nil {
				return nil, err
			}
			ref.OnDelete = act
		case kwUPDATE:
			p.advance() // consume UPDATE
			act, err := p.parseReferenceAction()
			if err != nil {
				return nil, err
			}
			ref.OnUpdate = act
		default:
			return nil, p.syntaxErrorAtCur()
		}
	}

	return ref, nil
}

// parseReferenceAction parses CASCADE | SET NULL | SET DEFAULT | RESTRICT | NO ACTION.
func (p *Parser) parseReferenceAction() (ast.ReferenceAction, error) {
	switch p.cur.Type {
	case kwCASCADE:
		p.advance()
		return ast.RefActCascade, nil
	case kwSET:
		p.advance() // consume SET
		switch p.cur.Type {
		case kwNULL:
			p.advance()
			return ast.RefActSetNull, nil
		case kwDEFAULT:
			p.advance()
			return ast.RefActSetDefault, nil
		default:
			return ast.RefActNone, p.syntaxErrorAtCur()
		}
	case kwRESTRICT:
		p.advance()
		return ast.RefActRestrict, nil
	case kwNO:
		p.advance() // consume NO
		if _, err := p.expect(kwACTION); err != nil {
			return ast.RefActNone, err
		}
		return ast.RefActNoAction, nil
	default:
		return ast.RefActNone, p.syntaxErrorAtCur()
	}
}

// parseConstraintProperties consumes optional constraint property keywords
// (ENFORCED, NOT ENFORCED, DEFERRABLE, NOT DEFERRABLE, INITIALLY DEFERRED,
// INITIALLY IMMEDIATE, ENABLE, DISABLE, VALIDATE, NOVALIDATE, RELY, NORELY).
// These are parsed and discarded — not stored in the AST.
func (p *Parser) parseConstraintProperties() {
	for {
		switch p.cur.Type {
		case kwNOT:
			next := p.peekNext()
			if next.Type == kwENFORCED || next.Type == kwDEFERRABLE {
				p.advance() // consume NOT
				p.advance() // consume ENFORCED/DEFERRABLE
			} else {
				return
			}
		case kwENFORCED, kwDEFERRABLE:
			p.advance()
		case kwINITIALLY:
			p.advance() // consume INITIALLY
			// Expect DEFERRED or IMMEDIATE — consume whatever follows
			p.advance()
		case kwENABLE, kwDISABLE:
			p.advance()
			// Optional VALIDATE/NOVALIDATE
			if p.cur.Type == kwVALIDATE || p.cur.Type == kwNOVALIDATE {
				p.advance()
			}
		case kwVALIDATE, kwNOVALIDATE:
			p.advance()
		case kwRELY, kwNORELY:
			p.advance()
		default:
			return
		}
	}
}

// parseIdentitySpec parses IDENTITY or AUTOINCREMENT options:
//   IDENTITY [(start, increment)] [START [WITH] [=] n] [INCREMENT [BY] [=] n] [ORDER|NOORDER]
func (p *Parser) parseIdentitySpec() (*ast.IdentitySpec, error) {
	p.advance() // consume IDENTITY or AUTOINCREMENT
	spec := &ast.IdentitySpec{}

	// Optional (start, increment)
	if p.cur.Type == '(' {
		p.advance() // consume (
		startTok, err := p.expect(tokInt)
		if err != nil {
			return nil, err
		}
		start := startTok.Ival
		spec.Start = &start

		if _, err := p.expect(','); err != nil {
			return nil, err
		}

		incrTok, err := p.expect(tokInt)
		if err != nil {
			return nil, err
		}
		incr := incrTok.Ival
		spec.Increment = &incr

		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}

	// Optional START WITH [=] n
	if p.cur.Type == kwSTART {
		p.advance() // consume START
		if p.cur.Type == kwWITH {
			p.advance() // consume WITH
		}
		p.match('=') // optional =
		tok, err := p.expect(tokInt)
		if err != nil {
			return nil, err
		}
		start := tok.Ival
		spec.Start = &start
	}

	// Optional INCREMENT BY [=] n
	if p.cur.Type == kwINCREMENT {
		p.advance() // consume INCREMENT
		if p.cur.Type == kwBY {
			p.advance() // consume BY
		}
		p.match('=') // optional =
		tok, err := p.expect(tokInt)
		if err != nil {
			return nil, err
		}
		incr := tok.Ival
		spec.Increment = &incr
	}

	// Optional ORDER/NOORDER
	switch p.cur.Type {
	case kwORDER:
		p.advance()
		b := true
		spec.Order = &b
	case kwNOORDER:
		p.advance()
		b := false
		spec.Order = &b
	}

	return spec, nil
}

// parseTagAssignments parses TAG ( name = 'value', ... ).
// Expects current token to be kwTAG.
func (p *Parser) parseTagAssignments() ([]*ast.TagAssignment, error) {
	if _, err := p.expect(kwTAG); err != nil {
		return nil, err
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	var tags []*ast.TagAssignment
	for {
		name, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		valTok, err := p.expect(tokString)
		if err != nil {
			return nil, err
		}
		tags = append(tags, &ast.TagAssignment{Name: name, Value: valTok.Str})

		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return tags, nil
}

// parseCloneSource parses CLONE name [AT|BEFORE (TIMESTAMP => val | OFFSET => val | STATEMENT => id)].
func (p *Parser) parseCloneSource() (*ast.CloneSource, error) {
	p.advance() // consume CLONE
	source, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	cs := &ast.CloneSource{Source: source}

	// Optional AT or BEFORE with time travel
	if p.cur.Type == kwAT || p.cur.Type == kwBEFORE {
		if p.cur.Type == kwAT {
			cs.AtBefore = "AT"
		} else {
			cs.AtBefore = "BEFORE"
		}
		p.advance() // consume AT/BEFORE

		if _, err := p.expect('('); err != nil {
			return nil, err
		}

		// TIMESTAMP => val | OFFSET => val | STATEMENT => id
		switch p.cur.Type {
		case kwTIMESTAMP:
			cs.Kind = "TIMESTAMP"
			p.advance()
		case kwOFFSET:
			cs.Kind = "OFFSET"
			p.advance()
		case kwSTATEMENT:
			cs.Kind = "STATEMENT"
			p.advance()
		default:
			return nil, p.syntaxErrorAtCur()
		}

		if _, err := p.expect(tokAssoc); err != nil {
			return nil, err
		}

		// Value: string literal or identifier
		switch p.cur.Type {
		case tokString:
			cs.Value = p.cur.Str
			p.advance()
		case tokIdent:
			cs.Value = p.cur.Str
			p.advance()
		case tokInt:
			cs.Value = p.cur.Str
			p.advance()
		default:
			// Accept any token as the value string
			cs.Value = p.cur.Str
			p.advance()
		}

		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}

	return cs, nil
}

// parseTableProperties parses optional table-level properties after the column list:
// CLUSTER BY, COPY GRANTS, COMMENT, WITH TAG, DATA_RETENTION_TIME_IN_DAYS,
// CHANGE_TRACKING, DEFAULT_DDL_COLLATION, STAGE_FILE_FORMAT, STAGE_COPY_OPTIONS,
// WITH ROW ACCESS POLICY.
func (p *Parser) parseTableProperties(stmt *ast.CreateTableStmt) error {
	for {
		switch p.cur.Type {
		case kwCLUSTER:
			p.advance() // consume CLUSTER
			if _, err := p.expect(kwBY); err != nil {
				return err
			}
			// Optional LINEAR
			if p.cur.Type == kwLINEAR {
				p.advance()
				stmt.Linear = true
			}
			if _, err := p.expect('('); err != nil {
				return err
			}
			exprs, err := p.parseExprList()
			if err != nil {
				return err
			}
			stmt.ClusterBy = exprs
			if _, err := p.expect(')'); err != nil {
				return err
			}

		case kwCOPY:
			p.advance() // consume COPY
			if _, err := p.expect(kwGRANTS); err != nil {
				return err
			}
			stmt.CopyGrants = true

		case kwCOMMENT:
			p.advance() // consume COMMENT
			if _, err := p.expect('='); err != nil {
				return err
			}
			tok, err := p.expect(tokString)
			if err != nil {
				return err
			}
			s := tok.Str
			stmt.Comment = &s

		case kwWITH:
			next := p.peekNext()
			switch next.Type {
			case kwTAG:
				p.advance() // consume WITH
				tags, err := p.parseTagAssignments()
				if err != nil {
					return err
				}
				stmt.Tags = tags
			case kwROW:
				// WITH ROW ACCESS POLICY — consume and discard
				p.advance() // consume WITH
				p.advance() // consume ROW
				if _, err := p.expect(kwACCESS); err != nil {
					return err
				}
				if _, err := p.expect(kwPOLICY); err != nil {
					return err
				}
				// Policy name
				if _, err := p.parseObjectName(); err != nil {
					return err
				}
				// ON (col, col, ...)
				if _, err := p.expect(kwON); err != nil {
					return err
				}
				if err := p.skipParenthesized(); err != nil {
					return err
				}
			default:
				return nil // not a table property
			}

		case kwTAG:
			tags, err := p.parseTagAssignments()
			if err != nil {
				return err
			}
			stmt.Tags = tags

		case kwDATA_RETENTION_TIME_IN_DAYS:
			// Consume and discard: DATA_RETENTION_TIME_IN_DAYS = n
			p.advance()
			p.match('=')
			if p.cur.Type == tokInt {
				p.advance()
			}

		case kwCHANGE_TRACKING:
			// Consume and discard: CHANGE_TRACKING = TRUE|FALSE
			p.advance()
			p.match('=')
			if p.cur.Type == kwTRUE || p.cur.Type == kwFALSE {
				p.advance()
			}

		case kwDEFAULT_DDL_COLLATION:
			// Consume and discard: DEFAULT_DDL_COLLATION_ = 'string'
			p.advance()
			p.match('=')
			if p.cur.Type == tokString {
				p.advance()
			}

		case kwSTAGE_FILE_FORMAT, kwSTAGE_COPY_OPTIONS:
			// Consume and discard: STAGE_FILE_FORMAT = (...) or STAGE_COPY_OPTIONS = (...)
			p.advance()
			p.match('=')
			if p.cur.Type == '(' {
				if err := p.skipParenthesized(); err != nil {
					return err
				}
			}

		default:
			return nil // no more properties
		}
	}
}

// skipParenthesized consumes everything from ( to the matching ), including nested parens.
func (p *Parser) skipParenthesized() error {
	if _, err := p.expect('('); err != nil {
		return err
	}
	depth := 1
	for depth > 0 && p.cur.Type != tokEOF {
		switch p.cur.Type {
		case '(':
			depth++
		case ')':
			depth--
		}
		p.advance()
	}
	return nil
}
```

- [ ] **Step 3: Update parser.go dispatch**

In `snowflake/parser/parser.go`, replace line 186:

Change:
```go
	case kwCREATE:
		return p.unsupported("CREATE")
```

To:
```go
	case kwCREATE:
		return p.parseCreateStmt()
```

- [ ] **Step 4: Check for kwTRUE/kwFALSE/kwUSING/kwBY/kwUPDATE existence**

Run:
```bash
cd /Users/h3n4l/OpenSource/omni/.worktrees/create-table && grep -E 'kwTRUE\b|kwFALSE\b|kwUSING\b|kwBY\b|kwUPDATE\b' snowflake/parser/tokens.go | head -10
```

Expected: all five should exist. If any are missing, add them.

- [ ] **Step 5: Verify compilation**

Run:
```bash
cd /Users/h3n4l/OpenSource/omni/.worktrees/create-table && go build ./snowflake/...
```

Expected: clean compile. Fix any undefined constants or signature mismatches.

- [ ] **Step 6: Commit**

```bash
cd /Users/h3n4l/OpenSource/omni/.worktrees/create-table
git add snowflake/parser/create_table.go snowflake/parser/parser.go
git commit -m "feat(snowflake): add CREATE TABLE parser (T2.2 step 2)

Implement parseCreateStmt dispatch and full parseCreateTableStmt with
column definitions, inline/out-of-line constraints, FK references,
IDENTITY/AUTOINCREMENT, MASKING POLICY, COLLATE, virtual columns,
CLUSTER BY, COPY GRANTS, COMMENT, WITH TAG, CTAS, LIKE, CLONE.

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>"
```

---

### Task 3: Tests — Basic Table Forms

**Files:**
- Create: `snowflake/parser/create_table_test.go`

- [ ] **Step 1: Create test file with helper and basic tests**

Create `snowflake/parser/create_table_test.go`:

```go
package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// testParseCreateTable parses input and returns the first statement as
// *ast.CreateTableStmt plus any errors.
func testParseCreateTable(input string) (*ast.CreateTableStmt, []ParseError) {
	result := ParseBestEffort(input)
	if len(result.File.Stmts) == 0 {
		return nil, result.Errors
	}
	stmt, ok := result.File.Stmts[0].(*ast.CreateTableStmt)
	if !ok {
		return nil, append(result.Errors, ParseError{Msg: "not a CreateTableStmt"})
	}
	return stmt, result.Errors
}

// ---------------------------------------------------------------------------
// 1. Simplest CREATE TABLE
// ---------------------------------------------------------------------------

func TestCreateTable_Basic(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (id INT)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Name.Name.Name != "t" {
		t.Errorf("table name = %q, want %q", stmt.Name.Name.Name, "t")
	}
	if len(stmt.Columns) != 1 {
		t.Fatalf("columns = %d, want 1", len(stmt.Columns))
	}
	col := stmt.Columns[0]
	if col.Name.Name != "id" {
		t.Errorf("col name = %q, want %q", col.Name.Name, "id")
	}
	if col.DataType == nil || col.DataType.Kind != ast.TypeInt {
		t.Errorf("col type = %v, want TypeInt", col.DataType)
	}
}

// ---------------------------------------------------------------------------
// 2. Multiple columns
// ---------------------------------------------------------------------------

func TestCreateTable_MultipleColumns(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (id INT, name VARCHAR(100), active BOOLEAN)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(stmt.Columns) != 3 {
		t.Fatalf("columns = %d, want 3", len(stmt.Columns))
	}
	if stmt.Columns[0].Name.Name != "id" {
		t.Errorf("col[0] = %q, want %q", stmt.Columns[0].Name.Name, "id")
	}
	if stmt.Columns[1].Name.Name != "name" {
		t.Errorf("col[1] = %q, want %q", stmt.Columns[1].Name.Name, "name")
	}
	if stmt.Columns[1].DataType.Kind != ast.TypeVarchar {
		t.Errorf("col[1] type = %v, want TypeVarchar", stmt.Columns[1].DataType.Kind)
	}
	if len(stmt.Columns[1].DataType.Params) != 1 || stmt.Columns[1].DataType.Params[0] != 100 {
		t.Errorf("col[1] params = %v, want [100]", stmt.Columns[1].DataType.Params)
	}
	if stmt.Columns[2].DataType.Kind != ast.TypeBoolean {
		t.Errorf("col[2] type = %v, want TypeBoolean", stmt.Columns[2].DataType.Kind)
	}
}

// ---------------------------------------------------------------------------
// 3. Modifiers: OR REPLACE, TRANSIENT, IF NOT EXISTS
// ---------------------------------------------------------------------------

func TestCreateTable_OrReplace(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE OR REPLACE TABLE t (id INT)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.OrReplace {
		t.Error("OrReplace should be true")
	}
}

func TestCreateTable_Transient(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TRANSIENT TABLE t (id INT)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.Transient {
		t.Error("Transient should be true")
	}
}

func TestCreateTable_Temporary(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TEMPORARY TABLE t (id INT)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.Temporary {
		t.Error("Temporary should be true")
	}
}

func TestCreateTable_Volatile(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE VOLATILE TABLE t (id INT)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.Volatile {
		t.Error("Volatile should be true")
	}
}

func TestCreateTable_LocalTemporary(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE LOCAL TEMPORARY TABLE t (id INT)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.Temporary {
		t.Error("Temporary should be true")
	}
}

func TestCreateTable_IfNotExists(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE IF NOT EXISTS t (id INT)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.IfNotExists {
		t.Error("IfNotExists should be true")
	}
}

func TestCreateTable_AllModifiers(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE OR REPLACE TRANSIENT TABLE IF NOT EXISTS mydb.myschema.t (id INT)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.OrReplace {
		t.Error("OrReplace should be true")
	}
	if !stmt.Transient {
		t.Error("Transient should be true")
	}
	if !stmt.IfNotExists {
		t.Error("IfNotExists should be true")
	}
	if stmt.Name.Normalize() != "MYDB.MYSCHEMA.T" {
		t.Errorf("name = %q, want MYDB.MYSCHEMA.T", stmt.Name.Normalize())
	}
}

// ---------------------------------------------------------------------------
// 4. LIKE
// ---------------------------------------------------------------------------

func TestCreateTable_Like(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t2 LIKE db.schema.t1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Like == nil {
		t.Fatal("Like should not be nil")
	}
	if stmt.Like.Normalize() != "DB.SCHEMA.T1" {
		t.Errorf("like = %q, want DB.SCHEMA.T1", stmt.Like.Normalize())
	}
	if len(stmt.Columns) != 0 {
		t.Errorf("columns = %d, want 0", len(stmt.Columns))
	}
}

// ---------------------------------------------------------------------------
// 5. CLONE
// ---------------------------------------------------------------------------

func TestCreateTable_Clone(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t2 CLONE t1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Clone == nil {
		t.Fatal("Clone should not be nil")
	}
	if stmt.Clone.Source.Name.Name != "t1" {
		t.Errorf("clone source = %q, want t1", stmt.Clone.Source.Name.Name)
	}
	if stmt.Clone.AtBefore != "" {
		t.Errorf("clone AtBefore = %q, want empty", stmt.Clone.AtBefore)
	}
}

func TestCreateTable_CloneAtTimestamp(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t2 CLONE t1 AT (TIMESTAMP => '2024-01-01')")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Clone.AtBefore != "AT" {
		t.Errorf("AtBefore = %q, want AT", stmt.Clone.AtBefore)
	}
	if stmt.Clone.Kind != "TIMESTAMP" {
		t.Errorf("Kind = %q, want TIMESTAMP", stmt.Clone.Kind)
	}
	if stmt.Clone.Value != "2024-01-01" {
		t.Errorf("Value = %q, want 2024-01-01", stmt.Clone.Value)
	}
}

func TestCreateTable_CloneBeforeStatement(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t2 CLONE t1 BEFORE (STATEMENT => '8e5d0ca9-005e-44e6-b858-a8f5b37c5726')")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Clone.AtBefore != "BEFORE" {
		t.Errorf("AtBefore = %q, want BEFORE", stmt.Clone.AtBefore)
	}
	if stmt.Clone.Kind != "STATEMENT" {
		t.Errorf("Kind = %q, want STATEMENT", stmt.Clone.Kind)
	}
}

// ---------------------------------------------------------------------------
// 6. CTAS
// ---------------------------------------------------------------------------

func TestCreateTable_CTAS(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t AS SELECT 1 AS id, 'hello' AS name")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.AsSelect == nil {
		t.Fatal("AsSelect should not be nil")
	}
	sel, ok := stmt.AsSelect.(*ast.SelectStmt)
	if !ok {
		t.Fatalf("AsSelect type = %T, want *ast.SelectStmt", stmt.AsSelect)
	}
	if len(sel.Targets) != 2 {
		t.Errorf("select targets = %d, want 2", len(sel.Targets))
	}
}

func TestCreateTable_CTASWithColumns(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (id INT, name VARCHAR) AS SELECT 1, 'hello'")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(stmt.Columns) != 2 {
		t.Fatalf("columns = %d, want 2", len(stmt.Columns))
	}
	if stmt.AsSelect == nil {
		t.Fatal("AsSelect should not be nil")
	}
}
```

- [ ] **Step 2: Run tests**

Run:
```bash
cd /Users/h3n4l/OpenSource/omni/.worktrees/create-table && go test ./snowflake/parser/ -run TestCreateTable -v -count=1
```

Expected: all tests PASS.

- [ ] **Step 3: Commit**

```bash
cd /Users/h3n4l/OpenSource/omni/.worktrees/create-table
git add snowflake/parser/create_table_test.go
git commit -m "test(snowflake): CREATE TABLE basic forms + modifiers + LIKE/CLONE/CTAS (T2.2 step 3)

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>"
```

---

### Task 4: Tests — Column Features + Constraints

**Files:**
- Modify: `snowflake/parser/create_table_test.go`

- [ ] **Step 1: Add column option tests**

Append to `create_table_test.go`:

```go
// ---------------------------------------------------------------------------
// 7. Column NOT NULL / NULL
// ---------------------------------------------------------------------------

func TestCreateTable_ColumnNotNull(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (id INT NOT NULL, name VARCHAR NULL)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.Columns[0].NotNull {
		t.Error("col[0] NotNull should be true")
	}
	if !stmt.Columns[1].Nullable {
		t.Error("col[1] Nullable should be true")
	}
}

// ---------------------------------------------------------------------------
// 8. Column DEFAULT
// ---------------------------------------------------------------------------

func TestCreateTable_ColumnDefault(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (id INT DEFAULT 0, name VARCHAR DEFAULT 'unknown')")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Columns[0].Default == nil {
		t.Fatal("col[0] Default should not be nil")
	}
	lit, ok := stmt.Columns[0].Default.(*ast.Literal)
	if !ok {
		t.Fatalf("col[0] Default type = %T, want *ast.Literal", stmt.Columns[0].Default)
	}
	if lit.Kind != ast.LitInt || lit.Ival != 0 {
		t.Errorf("default = %v/%d, want LitInt/0", lit.Kind, lit.Ival)
	}
}

// ---------------------------------------------------------------------------
// 9. IDENTITY / AUTOINCREMENT
// ---------------------------------------------------------------------------

func TestCreateTable_Identity(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (id INT IDENTITY(1, 1))")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	spec := stmt.Columns[0].Identity
	if spec == nil {
		t.Fatal("Identity should not be nil")
	}
	if spec.Start == nil || *spec.Start != 1 {
		t.Errorf("Start = %v, want 1", spec.Start)
	}
	if spec.Increment == nil || *spec.Increment != 1 {
		t.Errorf("Increment = %v, want 1", spec.Increment)
	}
}

func TestCreateTable_AutoincrementStartIncrement(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (id INT AUTOINCREMENT START 100 INCREMENT 10)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	spec := stmt.Columns[0].Identity
	if spec == nil {
		t.Fatal("Identity should not be nil")
	}
	if *spec.Start != 100 {
		t.Errorf("Start = %d, want 100", *spec.Start)
	}
	if *spec.Increment != 10 {
		t.Errorf("Increment = %d, want 10", *spec.Increment)
	}
}

func TestCreateTable_IdentityOrder(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (id INT IDENTITY ORDER)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	spec := stmt.Columns[0].Identity
	if spec == nil {
		t.Fatal("Identity should not be nil")
	}
	if spec.Order == nil || !*spec.Order {
		t.Error("Order should be true")
	}
}

// ---------------------------------------------------------------------------
// 10. COLLATE
// ---------------------------------------------------------------------------

func TestCreateTable_Collate(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (name VARCHAR COLLATE 'en-ci')")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Columns[0].Collate != "en-ci" {
		t.Errorf("Collate = %q, want %q", stmt.Columns[0].Collate, "en-ci")
	}
}

// ---------------------------------------------------------------------------
// 11. COMMENT on column
// ---------------------------------------------------------------------------

func TestCreateTable_ColumnComment(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (id INT COMMENT 'primary identifier')")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Columns[0].Comment == nil || *stmt.Columns[0].Comment != "primary identifier" {
		t.Errorf("Comment = %v, want 'primary identifier'", stmt.Columns[0].Comment)
	}
}

// ---------------------------------------------------------------------------
// 12. Inline constraints
// ---------------------------------------------------------------------------

func TestCreateTable_InlinePrimaryKey(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (id INT PRIMARY KEY)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	ic := stmt.Columns[0].InlineConstraint
	if ic == nil {
		t.Fatal("InlineConstraint should not be nil")
	}
	if ic.Type != ast.ConstrPrimaryKey {
		t.Errorf("type = %v, want ConstrPrimaryKey", ic.Type)
	}
}

func TestCreateTable_InlineUnique(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (email VARCHAR UNIQUE)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	ic := stmt.Columns[0].InlineConstraint
	if ic == nil {
		t.Fatal("InlineConstraint should not be nil")
	}
	if ic.Type != ast.ConstrUnique {
		t.Errorf("type = %v, want ConstrUnique", ic.Type)
	}
}

func TestCreateTable_InlineForeignKey(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE orders (customer_id INT FOREIGN KEY REFERENCES customers (id))")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	ic := stmt.Columns[0].InlineConstraint
	if ic == nil {
		t.Fatal("InlineConstraint should not be nil")
	}
	if ic.Type != ast.ConstrForeignKey {
		t.Errorf("type = %v, want ConstrForeignKey", ic.Type)
	}
	if ic.References == nil {
		t.Fatal("References should not be nil")
	}
	if ic.References.Table.Name.Name != "customers" {
		t.Errorf("ref table = %q, want customers", ic.References.Table.Name.Name)
	}
	if len(ic.References.Columns) != 1 || ic.References.Columns[0].Name != "id" {
		t.Errorf("ref cols = %v, want [id]", ic.References.Columns)
	}
}

func TestCreateTable_NamedInlineConstraint(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (id INT CONSTRAINT pk_t PRIMARY KEY)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	ic := stmt.Columns[0].InlineConstraint
	if ic.Name.Name != "pk_t" {
		t.Errorf("constraint name = %q, want pk_t", ic.Name.Name)
	}
}

// ---------------------------------------------------------------------------
// 13. Out-of-line (table-level) constraints
// ---------------------------------------------------------------------------

func TestCreateTable_OutOfLinePK(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (id INT, name VARCHAR, PRIMARY KEY (id))")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(stmt.Constraints) != 1 {
		t.Fatalf("constraints = %d, want 1", len(stmt.Constraints))
	}
	c := stmt.Constraints[0]
	if c.Type != ast.ConstrPrimaryKey {
		t.Errorf("type = %v, want ConstrPrimaryKey", c.Type)
	}
	if len(c.Columns) != 1 || c.Columns[0].Name != "id" {
		t.Errorf("columns = %v, want [id]", c.Columns)
	}
}

func TestCreateTable_OutOfLineCompositePK(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (a INT, b INT, PRIMARY KEY (a, b))")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	c := stmt.Constraints[0]
	if len(c.Columns) != 2 {
		t.Fatalf("columns = %d, want 2", len(c.Columns))
	}
	if c.Columns[0].Name != "a" || c.Columns[1].Name != "b" {
		t.Errorf("columns = [%q, %q], want [a, b]", c.Columns[0].Name, c.Columns[1].Name)
	}
}

func TestCreateTable_OutOfLineFK(t *testing.T) {
	stmt, errs := testParseCreateTable(`CREATE TABLE orders (
		id INT,
		customer_id INT,
		FOREIGN KEY (customer_id) REFERENCES customers (id) ON DELETE CASCADE ON UPDATE SET NULL
	)`)
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(stmt.Constraints) != 1 {
		t.Fatalf("constraints = %d, want 1", len(stmt.Constraints))
	}
	c := stmt.Constraints[0]
	if c.Type != ast.ConstrForeignKey {
		t.Errorf("type = %v, want ConstrForeignKey", c.Type)
	}
	if c.References.OnDelete != ast.RefActCascade {
		t.Errorf("OnDelete = %v, want RefActCascade", c.References.OnDelete)
	}
	if c.References.OnUpdate != ast.RefActSetNull {
		t.Errorf("OnUpdate = %v, want RefActSetNull", c.References.OnUpdate)
	}
}

func TestCreateTable_OutOfLineUnique(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (a INT, b INT, UNIQUE (a, b))")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(stmt.Constraints) != 1 {
		t.Fatalf("constraints = %d, want 1", len(stmt.Constraints))
	}
	if stmt.Constraints[0].Type != ast.ConstrUnique {
		t.Errorf("type = %v, want ConstrUnique", stmt.Constraints[0].Type)
	}
}

func TestCreateTable_NamedConstraint(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (id INT, CONSTRAINT pk_t PRIMARY KEY (id))")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Constraints[0].Name.Name != "pk_t" {
		t.Errorf("constraint name = %q, want pk_t", stmt.Constraints[0].Name.Name)
	}
}

func TestCreateTable_MixedColumnsAndConstraints(t *testing.T) {
	stmt, errs := testParseCreateTable(`CREATE TABLE t (
		id INT NOT NULL,
		name VARCHAR(100),
		email VARCHAR(255),
		PRIMARY KEY (id),
		UNIQUE (email),
		FOREIGN KEY (name) REFERENCES other (name)
	)`)
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(stmt.Columns) != 3 {
		t.Errorf("columns = %d, want 3", len(stmt.Columns))
	}
	if len(stmt.Constraints) != 3 {
		t.Errorf("constraints = %d, want 3", len(stmt.Constraints))
	}
}
```

- [ ] **Step 2: Run tests**

Run:
```bash
cd /Users/h3n4l/OpenSource/omni/.worktrees/create-table && go test ./snowflake/parser/ -run TestCreateTable -v -count=1
```

Expected: all tests PASS.

- [ ] **Step 3: Commit**

```bash
cd /Users/h3n4l/OpenSource/omni/.worktrees/create-table
git add snowflake/parser/create_table_test.go
git commit -m "test(snowflake): CREATE TABLE column options + constraints (T2.2 step 4)

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>"
```

---

### Task 5: Tests — Table Properties + Acceptance Sweep

**Files:**
- Modify: `snowflake/parser/create_table_test.go`

- [ ] **Step 1: Add table property tests**

Append to `create_table_test.go`:

```go
// ---------------------------------------------------------------------------
// 14. CLUSTER BY
// ---------------------------------------------------------------------------

func TestCreateTable_ClusterBy(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (id INT) CLUSTER BY (id)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(stmt.ClusterBy) != 1 {
		t.Fatalf("ClusterBy = %d, want 1", len(stmt.ClusterBy))
	}
}

func TestCreateTable_ClusterByLinear(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (a INT, b INT) CLUSTER BY LINEAR (a, b)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.Linear {
		t.Error("Linear should be true")
	}
	if len(stmt.ClusterBy) != 2 {
		t.Errorf("ClusterBy = %d, want 2", len(stmt.ClusterBy))
	}
}

// ---------------------------------------------------------------------------
// 15. Table COMMENT
// ---------------------------------------------------------------------------

func TestCreateTable_TableComment(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (id INT) COMMENT = 'my table'")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Comment == nil || *stmt.Comment != "my table" {
		t.Errorf("Comment = %v, want 'my table'", stmt.Comment)
	}
}

// ---------------------------------------------------------------------------
// 16. COPY GRANTS
// ---------------------------------------------------------------------------

func TestCreateTable_CopyGrants(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (id INT) COPY GRANTS")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.CopyGrants {
		t.Error("CopyGrants should be true")
	}
}

// ---------------------------------------------------------------------------
// 17. WITH TAG
// ---------------------------------------------------------------------------

func TestCreateTable_WithTag(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (id INT) WITH TAG (env = 'prod', team = 'data')")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(stmt.Tags) != 2 {
		t.Fatalf("Tags = %d, want 2", len(stmt.Tags))
	}
	if stmt.Tags[0].Name.Name.Name != "env" || stmt.Tags[0].Value != "prod" {
		t.Errorf("tag[0] = %v/%v, want env/prod", stmt.Tags[0].Name, stmt.Tags[0].Value)
	}
	if stmt.Tags[1].Name.Name.Name != "team" || stmt.Tags[1].Value != "data" {
		t.Errorf("tag[1] = %v/%v, want team/data", stmt.Tags[1].Name, stmt.Tags[1].Value)
	}
}

// ---------------------------------------------------------------------------
// 18. Data retention / change tracking (consumed without error)
// ---------------------------------------------------------------------------

func TestCreateTable_DataRetention(t *testing.T) {
	_, errs := testParseCreateTable("CREATE TABLE t (id INT) DATA_RETENTION_TIME_IN_DAYS = 90")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
}

func TestCreateTable_ChangeTracking(t *testing.T) {
	_, errs := testParseCreateTable("CREATE TABLE t (id INT) CHANGE_TRACKING = TRUE")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
}

// ---------------------------------------------------------------------------
// 19. WITH MASKING POLICY on column
// ---------------------------------------------------------------------------

func TestCreateTable_MaskingPolicy(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (ssn VARCHAR WITH MASKING POLICY ssn_mask)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Columns[0].MaskingPolicy == nil {
		t.Fatal("MaskingPolicy should not be nil")
	}
	if stmt.Columns[0].MaskingPolicy.Name.Name != "ssn_mask" {
		t.Errorf("policy = %q, want ssn_mask", stmt.Columns[0].MaskingPolicy.Name.Name)
	}
}

// ---------------------------------------------------------------------------
// 20. Virtual column
// ---------------------------------------------------------------------------

func TestCreateTable_VirtualColumn(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (a INT, b INT, c INT AS (a + b))")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(stmt.Columns) != 3 {
		t.Fatalf("columns = %d, want 3", len(stmt.Columns))
	}
	if stmt.Columns[2].VirtualExpr == nil {
		t.Fatal("col[2] VirtualExpr should not be nil")
	}
	_, ok := stmt.Columns[2].VirtualExpr.(*ast.BinaryExpr)
	if !ok {
		t.Errorf("VirtualExpr type = %T, want *ast.BinaryExpr", stmt.Columns[2].VirtualExpr)
	}
}

// ---------------------------------------------------------------------------
// 21. Complex real-world CREATE TABLE
// ---------------------------------------------------------------------------

func TestCreateTable_Complex(t *testing.T) {
	stmt, errs := testParseCreateTable(`
		CREATE OR REPLACE TABLE mydb.myschema.orders (
			id INT IDENTITY(1, 1) NOT NULL,
			customer_id INT NOT NULL,
			amount NUMBER(10, 2) DEFAULT 0,
			status VARCHAR(50) COLLATE 'en-ci',
			created_at TIMESTAMP_NTZ DEFAULT CURRENT_TIMESTAMP(),
			CONSTRAINT pk_orders PRIMARY KEY (id),
			CONSTRAINT fk_customer FOREIGN KEY (customer_id) REFERENCES customers (id) ON DELETE CASCADE,
			UNIQUE (customer_id, created_at)
		)
		CLUSTER BY (customer_id)
		COMMENT = 'Order records'
		COPY GRANTS
	`)
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Name.Normalize() != "MYDB.MYSCHEMA.ORDERS" {
		t.Errorf("name = %q", stmt.Name.Normalize())
	}
	if len(stmt.Columns) != 5 {
		t.Errorf("columns = %d, want 5", len(stmt.Columns))
	}
	if len(stmt.Constraints) != 3 {
		t.Errorf("constraints = %d, want 3", len(stmt.Constraints))
	}
	if !stmt.OrReplace {
		t.Error("OrReplace should be true")
	}
	if !stmt.CopyGrants {
		t.Error("CopyGrants should be true")
	}
	if stmt.Comment == nil || *stmt.Comment != "Order records" {
		t.Errorf("Comment = %v", stmt.Comment)
	}
	if len(stmt.ClusterBy) != 1 {
		t.Errorf("ClusterBy = %d, want 1", len(stmt.ClusterBy))
	}
	// Identity on first column
	if stmt.Columns[0].Identity == nil {
		t.Error("col[0] Identity should not be nil")
	}
	// NOT NULL on first two columns
	if !stmt.Columns[0].NotNull || !stmt.Columns[1].NotNull {
		t.Error("first two columns should be NOT NULL")
	}
	// FK constraint
	fk := stmt.Constraints[1]
	if fk.Type != ast.ConstrForeignKey {
		t.Errorf("constraint[1] type = %v, want FK", fk.Type)
	}
	if fk.Name.Name != "fk_customer" {
		t.Errorf("constraint[1] name = %q, want fk_customer", fk.Name.Name)
	}
}
```

- [ ] **Step 2: Run full test suite**

Run:
```bash
cd /Users/h3n4l/OpenSource/omni/.worktrees/create-table && go test ./snowflake/... -count=1
```

Expected: all tests PASS (both new CREATE TABLE tests and existing SELECT/expression tests).

- [ ] **Step 3: Run gofmt**

Run:
```bash
cd /Users/h3n4l/OpenSource/omni/.worktrees/create-table && gofmt -l ./snowflake/
```

Expected: no files listed (all formatted). If files are listed, run `gofmt -w` on them.

- [ ] **Step 4: Check legacy corpus for CREATE TABLE**

Run:
```bash
cd /Users/h3n4l/OpenSource/omni/.worktrees/create-table && grep -ril 'CREATE TABLE' snowflake/parser/testdata/ 2>/dev/null || echo "no corpus CREATE TABLE files found"
```

If any legacy corpus files contain CREATE TABLE, verify they parse without errors. If not found, that's fine — the unit tests provide coverage.

- [ ] **Step 5: Commit**

```bash
cd /Users/h3n4l/OpenSource/omni/.worktrees/create-table
git add snowflake/parser/create_table_test.go
git commit -m "test(snowflake): CREATE TABLE properties + complex real-world + acceptance (T2.2 step 5)

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>"
```

---

## Self-Review

**1. Spec coverage:**
- Standard CREATE TABLE with columns: Task 2 (parser) + Task 3 (test 1-2)
- Modifiers (OR REPLACE, TRANSIENT, TEMP, VOLATILE, IF NOT EXISTS): Task 2 + Task 3 (tests 3)
- LIKE: Task 2 + Task 3 (test 4)
- CLONE with time travel: Task 2 + Task 3 (test 5)
- CTAS / CTAS with columns: Task 2 + Task 3 (test 6)
- NOT NULL / NULL: Task 2 + Task 4 (test 7)
- DEFAULT expr: Task 2 + Task 4 (test 8)
- IDENTITY/AUTOINCREMENT: Task 2 + Task 4 (test 9)
- COLLATE: Task 2 + Task 4 (test 10)
- Column COMMENT: Task 2 + Task 4 (test 11)
- Inline PK/UNIQUE/FK: Task 2 + Task 4 (test 12)
- Out-of-line PK/UNIQUE/FK: Task 2 + Task 4 (test 13)
- FK ON DELETE/UPDATE: Task 2 + Task 4 (test 13)
- Named constraints: Task 2 + Task 4 (test 13)
- CLUSTER BY / LINEAR: Task 2 + Task 5 (test 14)
- Table COMMENT: Task 2 + Task 5 (test 15)
- COPY GRANTS: Task 2 + Task 5 (test 16)
- WITH TAG: Task 2 + Task 5 (test 17)
- DATA_RETENTION / CHANGE_TRACKING consumed: Task 2 + Task 5 (test 18)
- MASKING POLICY: Task 2 + Task 5 (test 19)
- Virtual columns: Task 2 + Task 5 (test 20)
- Complex real-world: Task 5 (test 21)
- Node tags: Task 1
- Walker regeneration: Task 1
- Dispatch change: Task 2

**2. Placeholder scan:** No TBD/TODO/placeholder text found.

**3. Type consistency:** All types (CreateTableStmt, ColumnDef, TableConstraint, InlineConstraint, ForeignKeyRef, IdentitySpec, TagAssignment, CloneSource, ConstraintType, ReferenceAction) are defined in Task 1 and used consistently in Task 2-5. Field names match between definition and test assertions.
