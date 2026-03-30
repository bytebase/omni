# CosmosDB Engine Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a `cosmosdb/` engine to Omni with a complete, pure-Go hand-written parser for Azure Cosmos DB's NoSQL SQL API.

**Architecture:** Hand-written lexer + recursive descent parser with Pratt expression parsing, producing an AST that follows existing Omni conventions (same interfaces as `oracle/`, `mssql/`). Public API matches `pg/parse.go` pattern. Test corpus from `bytebase/parser/cosmosdb/examples/`.

**Tech Stack:** Pure Go (zero runtime dependencies), `testing` package for tests.

**Spec:** `docs/superpowers/specs/2026-03-30-cosmosdb-engine-design.md`

---

## File Map

| File | Responsibility |
|------|---------------|
| `cosmosdb/ast/node.go` | Core interfaces (Node, ExprNode, TableExpr, StmtNode), Loc, literal nodes (StringLit, NumberLit, BoolLit, NullLit, UndefinedLit, InfinityLit, NanLit), List |
| `cosmosdb/ast/parsenodes.go` | All statement, clause, table-expr, and expression AST node types (~24 types) |
| `cosmosdb/ast/outfuncs.go` | `NodeToString(Node) string` — deterministic text dump of any AST node |
| `cosmosdb/parser/lexer.go` | Tokenizer: Token struct, token constants, Lexer struct, NextToken() |
| `cosmosdb/parser/parser.go` | Parser struct, Parse(), advance/peek/match/expect helpers, parseIdentifier/parsePropertyName, top-level parseSelect |
| `cosmosdb/parser/select.go` | parseSelectClause, parseTargetEntry, parseGroupBy, parseHaving, parseOrderBy, parseSortExpr, parseOffsetLimit |
| `cosmosdb/parser/from.go` | parseFromClause, parseFromSource, parseContainerExpr, parseJoinClause |
| `cosmosdb/parser/expr.go` | Pratt parser: parseExpr(precedence), all prefix and infix parse functions |
| `cosmosdb/parser/function.go` | parseFuncCall, parseUDFCall |
| `cosmosdb/parser/parser_test.go` | Golden-file tests using 37 example SQL files |
| `cosmosdb/parser/testdata/*.sql` | 37 example SQL files from bytebase/parser/cosmosdb/examples/ |
| `cosmosdb/parse.go` | Public API: Parse(sql) -> []Statement with position tracking |
| `cosmosdb/parse_test.go` | Public API tests: positions, text slicing |

---

### Task 1: AST Core Interfaces and Literal Nodes

**Files:**
- Create: `cosmosdb/ast/node.go`

- [ ] **Step 1: Create `cosmosdb/ast/node.go`**

```go
// Package ast defines parse tree node types for Azure Cosmos DB NoSQL SQL API.
package ast

// Node is the interface implemented by all parse tree nodes.
type Node interface {
	nodeTag()
}

// ExprNode is the interface for expression nodes.
type ExprNode interface {
	Node
	exprNode()
}

// TableExpr is the interface for table reference nodes in FROM clauses.
type TableExpr interface {
	Node
	tableExpr()
}

// StmtNode is the interface for statement nodes.
type StmtNode interface {
	Node
	stmtNode()
}

// Loc represents a source location range (byte offsets).
type Loc struct {
	Start int // inclusive start byte offset, -1 = unknown
	End   int // exclusive end byte offset, -1 = unknown
}

// List represents a generic list of nodes.
type List struct {
	Items []Node
}

func (l *List) nodeTag() {}

// Len returns the number of items in the list.
func (l *List) Len() int {
	if l == nil {
		return 0
	}
	return len(l.Items)
}

// ---------------------------------------------------------------------------
// Literal nodes — all implement ExprNode
// ---------------------------------------------------------------------------

// StringLit represents a string constant ('hello' or "hello").
type StringLit struct {
	Val string
	Loc Loc
}

func (*StringLit) nodeTag()  {}
func (*StringLit) exprNode() {}

// NumberLit represents a numeric constant (integer, float, hex, scientific).
// Val stores the raw text to preserve the original representation.
type NumberLit struct {
	Val string
	Loc Loc
}

func (*NumberLit) nodeTag()  {}
func (*NumberLit) exprNode() {}

// BoolLit represents true or false.
type BoolLit struct {
	Val bool
	Loc Loc
}

func (*BoolLit) nodeTag()  {}
func (*BoolLit) exprNode() {}

// NullLit represents null.
type NullLit struct {
	Loc Loc
}

func (*NullLit) nodeTag()  {}
func (*NullLit) exprNode() {}

// UndefinedLit represents the CosmosDB-specific undefined constant.
type UndefinedLit struct {
	Loc Loc
}

func (*UndefinedLit) nodeTag()  {}
func (*UndefinedLit) exprNode() {}

// InfinityLit represents the Infinity constant (case-sensitive).
type InfinityLit struct {
	Loc Loc
}

func (*InfinityLit) nodeTag()  {}
func (*InfinityLit) exprNode() {}

// NanLit represents the NaN constant (case-sensitive).
type NanLit struct {
	Loc Loc
}

func (*NanLit) nodeTag()  {}
func (*NanLit) exprNode() {}
```

- [ ] **Step 2: Verify it compiles**

Run: `cd /Users/h3n4l/OpenSource/omni && go build ./cosmosdb/ast/`
Expected: Clean build, no errors.

- [ ] **Step 3: Commit**

```bash
git add cosmosdb/ast/node.go
git commit -m "feat(cosmosdb): add AST core interfaces and literal nodes"
```

---

### Task 2: AST Parse Nodes

**Files:**
- Create: `cosmosdb/ast/parsenodes.go`

- [ ] **Step 1: Create `cosmosdb/ast/parsenodes.go`**

```go
package ast

// ---------------------------------------------------------------------------
// Wrapper / top-level
// ---------------------------------------------------------------------------

// RawStmt wraps a parsed statement with position info.
type RawStmt struct {
	Stmt         StmtNode
	StmtLocation int // start byte offset
	StmtLen      int // length in bytes; 0 means "rest of string"
}

func (*RawStmt) nodeTag() {}

// ---------------------------------------------------------------------------
// Statement nodes
// ---------------------------------------------------------------------------

// SelectStmt represents a complete SELECT query.
type SelectStmt struct {
	Top         *int               // TOP n (nil if not specified)
	Distinct    bool               // DISTINCT modifier
	Value       bool               // VALUE modifier
	Star        bool               // SELECT *
	Targets     []*TargetEntry     // select list (nil when Star is true)
	From        TableExpr          // FROM clause (nil if omitted)
	Joins       []*JoinExpr        // JOIN clauses
	Where       ExprNode           // WHERE clause (nil if omitted)
	GroupBy     []ExprNode         // GROUP BY expressions
	Having      ExprNode           // HAVING clause (nil if omitted)
	OrderBy     []*SortExpr        // ORDER BY expressions
	OffsetLimit *OffsetLimitClause // OFFSET ... LIMIT ... (nil if omitted)
	Loc         Loc
}

func (*SelectStmt) nodeTag()  {}
func (*SelectStmt) stmtNode() {}

// TargetEntry represents one item in a SELECT list: expr [AS alias].
type TargetEntry struct {
	Expr  ExprNode
	Alias *string // nil if no alias
	Loc   Loc
}

func (*TargetEntry) nodeTag() {}

// SortExpr represents one item in an ORDER BY list.
type SortExpr struct {
	Expr ExprNode
	Desc bool // true for DESC, false for ASC (default)
	Rank bool // true for RANK expr
	Loc  Loc
}

func (*SortExpr) nodeTag() {}

// OffsetLimitClause represents OFFSET n LIMIT m.
type OffsetLimitClause struct {
	Offset ExprNode
	Limit  ExprNode
	Loc    Loc
}

func (*OffsetLimitClause) nodeTag() {}

// JoinExpr represents a JOIN clause: JOIN from_source.
type JoinExpr struct {
	Source TableExpr
	Loc    Loc
}

func (*JoinExpr) nodeTag() {}

// ---------------------------------------------------------------------------
// Table expression nodes — implement TableExpr
// ---------------------------------------------------------------------------

// ContainerRef references a container by name, or ROOT.
type ContainerRef struct {
	Name string // empty when Root is true
	Root bool
	Loc  Loc
}

func (*ContainerRef) nodeTag()    {}
func (*ContainerRef) tableExpr()  {}

// AliasedTableExpr wraps a table expression with an alias: source [AS] alias.
type AliasedTableExpr struct {
	Source TableExpr
	Alias  string
	Loc    Loc
}

func (*AliasedTableExpr) nodeTag()    {}
func (*AliasedTableExpr) tableExpr()  {}

// ArrayIterationExpr represents: alias IN container_expression.
type ArrayIterationExpr struct {
	Alias  string
	Source TableExpr
	Loc    Loc
}

func (*ArrayIterationExpr) nodeTag()    {}
func (*ArrayIterationExpr) tableExpr()  {}

// SubqueryExpr represents a parenthesized SELECT in FROM position.
type SubqueryExpr struct {
	Select *SelectStmt
	Loc    Loc
}

func (*SubqueryExpr) nodeTag()    {}
func (*SubqueryExpr) tableExpr()  {}

// ---------------------------------------------------------------------------
// Expression nodes — implement ExprNode
//
// DotAccessExpr and BracketAccessExpr also implement TableExpr since they
// appear in both FROM container_expression and scalar_expression contexts.
// ---------------------------------------------------------------------------

// ColumnRef represents an identifier reference (alias, property name).
type ColumnRef struct {
	Name string
	Loc  Loc
}

func (*ColumnRef) nodeTag()  {}
func (*ColumnRef) exprNode() {}

// DotAccessExpr represents property access: expr.property.
type DotAccessExpr struct {
	Expr     ExprNode
	Property string
	Loc      Loc
}

func (*DotAccessExpr) nodeTag()   {}
func (*DotAccessExpr) exprNode()  {}
func (*DotAccessExpr) tableExpr() {}

// BracketAccessExpr represents bracket access: expr[index].
type BracketAccessExpr struct {
	Expr  ExprNode
	Index ExprNode
	Loc   Loc
}

func (*BracketAccessExpr) nodeTag()   {}
func (*BracketAccessExpr) exprNode()  {}
func (*BracketAccessExpr) tableExpr() {}

// BinaryExpr represents a binary operation: left op right.
type BinaryExpr struct {
	Op    string
	Left  ExprNode
	Right ExprNode
	Loc   Loc
}

func (*BinaryExpr) nodeTag()  {}
func (*BinaryExpr) exprNode() {}

// UnaryExpr represents a unary operation: op operand.
type UnaryExpr struct {
	Op      string
	Operand ExprNode
	Loc     Loc
}

func (*UnaryExpr) nodeTag()  {}
func (*UnaryExpr) exprNode() {}

// TernaryExpr represents: cond ? then : else.
type TernaryExpr struct {
	Cond ExprNode
	Then ExprNode
	Else ExprNode
	Loc  Loc
}

func (*TernaryExpr) nodeTag()  {}
func (*TernaryExpr) exprNode() {}

// InExpr represents: expr [NOT] IN (list).
type InExpr struct {
	Expr ExprNode
	List []ExprNode
	Not  bool
	Loc  Loc
}

func (*InExpr) nodeTag()  {}
func (*InExpr) exprNode() {}

// BetweenExpr represents: expr [NOT] BETWEEN low AND high.
type BetweenExpr struct {
	Expr ExprNode
	Low  ExprNode
	High ExprNode
	Not  bool
	Loc  Loc
}

func (*BetweenExpr) nodeTag()  {}
func (*BetweenExpr) exprNode() {}

// LikeExpr represents: expr [NOT] LIKE pattern [ESCAPE escape].
type LikeExpr struct {
	Expr    ExprNode
	Pattern ExprNode
	Escape  ExprNode // nil if no ESCAPE clause
	Not     bool
	Loc     Loc
}

func (*LikeExpr) nodeTag()  {}
func (*LikeExpr) exprNode() {}

// FuncCall represents a built-in function call: name([*|args]).
type FuncCall struct {
	Name string
	Args []ExprNode
	Star bool // true for COUNT(*)
	Loc  Loc
}

func (*FuncCall) nodeTag()  {}
func (*FuncCall) exprNode() {}

// UDFCall represents a user-defined function call: udf.name(args).
type UDFCall struct {
	Name string
	Args []ExprNode
	Loc  Loc
}

func (*UDFCall) nodeTag()  {}
func (*UDFCall) exprNode() {}

// ExistsExpr represents EXISTS(SELECT ...).
type ExistsExpr struct {
	Select *SelectStmt
	Loc    Loc
}

func (*ExistsExpr) nodeTag()  {}
func (*ExistsExpr) exprNode() {}

// ArrayExpr represents ARRAY(SELECT ...).
type ArrayExpr struct {
	Select *SelectStmt
	Loc    Loc
}

func (*ArrayExpr) nodeTag()  {}
func (*ArrayExpr) exprNode() {}

// CreateArrayExpr represents an array literal: [elem, ...].
type CreateArrayExpr struct {
	Elements []ExprNode
	Loc      Loc
}

func (*CreateArrayExpr) nodeTag()  {}
func (*CreateArrayExpr) exprNode() {}

// CreateObjectExpr represents an object literal: {key: val, ...}.
type CreateObjectExpr struct {
	Fields []*ObjectFieldPair
	Loc    Loc
}

func (*CreateObjectExpr) nodeTag()  {}
func (*CreateObjectExpr) exprNode() {}

// ObjectFieldPair represents one field in an object literal: key: value.
type ObjectFieldPair struct {
	Key   string
	Value ExprNode
	Loc   Loc
}

func (*ObjectFieldPair) nodeTag() {}

// ParamRef represents a parameter reference: @name.
type ParamRef struct {
	Name string
	Loc  Loc
}

func (*ParamRef) nodeTag()  {}
func (*ParamRef) exprNode() {}

// SubLink represents a parenthesized subquery in expression position: (SELECT ...).
type SubLink struct {
	Select *SelectStmt
	Loc    Loc
}

func (*SubLink) nodeTag()  {}
func (*SubLink) exprNode() {}
```

- [ ] **Step 2: Verify it compiles**

Run: `cd /Users/h3n4l/OpenSource/omni && go build ./cosmosdb/ast/`
Expected: Clean build.

- [ ] **Step 3: Commit**

```bash
git add cosmosdb/ast/parsenodes.go
git commit -m "feat(cosmosdb): add AST parse node types"
```

---

### Task 3: AST Output Functions

**Files:**
- Create: `cosmosdb/ast/outfuncs.go`

This produces deterministic text output for golden-file testing. The format follows the Oracle pattern: `{NODETYPE :field value ...}`.

- [ ] **Step 1: Create `cosmosdb/ast/outfuncs.go`**

```go
package ast

import (
	"fmt"
	"strings"
)

// NodeToString converts a Node to its string representation.
func NodeToString(node Node) string {
	if node == nil {
		return "<>"
	}
	var sb strings.Builder
	writeNode(&sb, node, 0)
	return sb.String()
}

func indent(sb *strings.Builder, level int) {
	for i := 0; i < level; i++ {
		sb.WriteString("  ")
	}
}

func writeLoc(sb *strings.Builder, loc Loc) {
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", loc.Start, loc.End))
}

func writeNode(sb *strings.Builder, node Node, level int) {
	if node == nil {
		sb.WriteString("<>")
		return
	}

	switch n := node.(type) {
	case *List:
		sb.WriteString("(")
		for i, item := range n.Items {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, item, level)
		}
		sb.WriteString(")")

	// Literals
	case *StringLit:
		sb.WriteString(fmt.Sprintf("{STRLIT :val %q", n.Val))
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *NumberLit:
		sb.WriteString(fmt.Sprintf("{NUMLIT :val %q", n.Val))
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *BoolLit:
		sb.WriteString(fmt.Sprintf("{BOOLLIT :val %t", n.Val))
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *NullLit:
		sb.WriteString("{NULL")
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *UndefinedLit:
		sb.WriteString("{UNDEFINED")
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *InfinityLit:
		sb.WriteString("{INFINITY")
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *NanLit:
		sb.WriteString("{NAN")
		writeLoc(sb, n.Loc)
		sb.WriteString("}")

	// Wrapper
	case *RawStmt:
		sb.WriteString("{RAWSTMT")
		sb.WriteString(fmt.Sprintf(" :stmt_location %d :stmt_len %d", n.StmtLocation, n.StmtLen))
		sb.WriteString("\n")
		indent(sb, level+1)
		sb.WriteString(":stmt ")
		writeNode(sb, n.Stmt, level+1)
		sb.WriteString("}")

	// Statement
	case *SelectStmt:
		writeSelectStmt(sb, n, level)

	// Clause nodes
	case *TargetEntry:
		sb.WriteString("{TARGET :expr ")
		writeNode(sb, n.Expr, level)
		if n.Alias != nil {
			sb.WriteString(fmt.Sprintf(" :alias %q", *n.Alias))
		}
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *SortExpr:
		sb.WriteString("{SORT :expr ")
		writeNode(sb, n.Expr, level)
		if n.Desc {
			sb.WriteString(" :desc true")
		}
		if n.Rank {
			sb.WriteString(" :rank true")
		}
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *OffsetLimitClause:
		sb.WriteString("{OFFSETLIMIT :offset ")
		writeNode(sb, n.Offset, level)
		sb.WriteString(" :limit ")
		writeNode(sb, n.Limit, level)
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *JoinExpr:
		sb.WriteString("{JOIN :source ")
		writeNode(sb, n.Source, level)
		writeLoc(sb, n.Loc)
		sb.WriteString("}")

	// Table expressions
	case *ContainerRef:
		if n.Root {
			sb.WriteString("{CONTAINER :root true")
		} else {
			sb.WriteString(fmt.Sprintf("{CONTAINER :name %q", n.Name))
		}
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *AliasedTableExpr:
		sb.WriteString("{ALIASED :source ")
		writeNode(sb, n.Source, level)
		sb.WriteString(fmt.Sprintf(" :alias %q", n.Alias))
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *ArrayIterationExpr:
		sb.WriteString(fmt.Sprintf("{ARRAYITER :alias %q :source ", n.Alias))
		writeNode(sb, n.Source, level)
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *SubqueryExpr:
		sb.WriteString("{SUBQUERY\n")
		indent(sb, level+1)
		sb.WriteString(":select ")
		writeNode(sb, n.Select, level+1)
		writeLoc(sb, n.Loc)
		sb.WriteString("}")

	// Expression nodes
	case *ColumnRef:
		sb.WriteString(fmt.Sprintf("{COLREF :name %q", n.Name))
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *DotAccessExpr:
		sb.WriteString("{DOT :expr ")
		writeNode(sb, n.Expr, level)
		sb.WriteString(fmt.Sprintf(" :property %q", n.Property))
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *BracketAccessExpr:
		sb.WriteString("{BRACKET :expr ")
		writeNode(sb, n.Expr, level)
		sb.WriteString(" :index ")
		writeNode(sb, n.Index, level)
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *BinaryExpr:
		sb.WriteString(fmt.Sprintf("{BINEXPR :op %q :left ", n.Op))
		writeNode(sb, n.Left, level)
		sb.WriteString(" :right ")
		writeNode(sb, n.Right, level)
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *UnaryExpr:
		sb.WriteString(fmt.Sprintf("{UNARYEXPR :op %q :operand ", n.Op))
		writeNode(sb, n.Operand, level)
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *TernaryExpr:
		sb.WriteString("{TERNARY :cond ")
		writeNode(sb, n.Cond, level)
		sb.WriteString(" :then ")
		writeNode(sb, n.Then, level)
		sb.WriteString(" :else ")
		writeNode(sb, n.Else, level)
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *InExpr:
		sb.WriteString("{IN")
		if n.Not {
			sb.WriteString(" :not true")
		}
		sb.WriteString(" :expr ")
		writeNode(sb, n.Expr, level)
		sb.WriteString(" :list (")
		for i, item := range n.List {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, item, level)
		}
		sb.WriteString(")")
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *BetweenExpr:
		sb.WriteString("{BETWEEN")
		if n.Not {
			sb.WriteString(" :not true")
		}
		sb.WriteString(" :expr ")
		writeNode(sb, n.Expr, level)
		sb.WriteString(" :low ")
		writeNode(sb, n.Low, level)
		sb.WriteString(" :high ")
		writeNode(sb, n.High, level)
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *LikeExpr:
		sb.WriteString("{LIKE")
		if n.Not {
			sb.WriteString(" :not true")
		}
		sb.WriteString(" :expr ")
		writeNode(sb, n.Expr, level)
		sb.WriteString(" :pattern ")
		writeNode(sb, n.Pattern, level)
		if n.Escape != nil {
			sb.WriteString(" :escape ")
			writeNode(sb, n.Escape, level)
		}
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *FuncCall:
		sb.WriteString(fmt.Sprintf("{FUNCCALL :name %q", n.Name))
		if n.Star {
			sb.WriteString(" :star true")
		}
		if len(n.Args) > 0 {
			sb.WriteString(" :args (")
			for i, arg := range n.Args {
				if i > 0 {
					sb.WriteString(" ")
				}
				writeNode(sb, arg, level)
			}
			sb.WriteString(")")
		}
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *UDFCall:
		sb.WriteString(fmt.Sprintf("{UDFCALL :name %q", n.Name))
		if len(n.Args) > 0 {
			sb.WriteString(" :args (")
			for i, arg := range n.Args {
				if i > 0 {
					sb.WriteString(" ")
				}
				writeNode(sb, arg, level)
			}
			sb.WriteString(")")
		}
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *ExistsExpr:
		sb.WriteString("{EXISTS\n")
		indent(sb, level+1)
		sb.WriteString(":select ")
		writeNode(sb, n.Select, level+1)
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *ArrayExpr:
		sb.WriteString("{ARRAY\n")
		indent(sb, level+1)
		sb.WriteString(":select ")
		writeNode(sb, n.Select, level+1)
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *CreateArrayExpr:
		sb.WriteString("{CREATEARRAY :elements (")
		for i, elem := range n.Elements {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, elem, level)
		}
		sb.WriteString(")")
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *CreateObjectExpr:
		sb.WriteString("{CREATEOBJECT :fields (")
		for i, f := range n.Fields {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, f, level)
		}
		sb.WriteString(")")
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *ObjectFieldPair:
		sb.WriteString(fmt.Sprintf("{FIELD :key %q :value ", n.Key))
		writeNode(sb, n.Value, level)
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *ParamRef:
		sb.WriteString(fmt.Sprintf("{PARAM :name %q", n.Name))
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *SubLink:
		sb.WriteString("{SUBLINK\n")
		indent(sb, level+1)
		sb.WriteString(":select ")
		writeNode(sb, n.Select, level+1)
		writeLoc(sb, n.Loc)
		sb.WriteString("}")

	default:
		sb.WriteString(fmt.Sprintf("{UNKNOWN %T}", node))
	}
}

func writeSelectStmt(sb *strings.Builder, n *SelectStmt, level int) {
	sb.WriteString("{SELECT")
	if n.Top != nil {
		sb.WriteString(fmt.Sprintf(" :top %d", *n.Top))
	}
	if n.Distinct {
		sb.WriteString(" :distinct true")
	}
	if n.Value {
		sb.WriteString(" :value true")
	}
	if n.Star {
		sb.WriteString(" :star true")
	}
	if len(n.Targets) > 0 {
		sb.WriteString("\n")
		indent(sb, level+1)
		sb.WriteString(":targets (")
		for i, t := range n.Targets {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, t, level+1)
		}
		sb.WriteString(")")
	}
	if n.From != nil {
		sb.WriteString("\n")
		indent(sb, level+1)
		sb.WriteString(":from ")
		writeNode(sb, n.From, level+1)
	}
	if len(n.Joins) > 0 {
		sb.WriteString("\n")
		indent(sb, level+1)
		sb.WriteString(":joins (")
		for i, j := range n.Joins {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, j, level+1)
		}
		sb.WriteString(")")
	}
	if n.Where != nil {
		sb.WriteString("\n")
		indent(sb, level+1)
		sb.WriteString(":where ")
		writeNode(sb, n.Where, level+1)
	}
	if len(n.GroupBy) > 0 {
		sb.WriteString("\n")
		indent(sb, level+1)
		sb.WriteString(":group_by (")
		for i, g := range n.GroupBy {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, g, level+1)
		}
		sb.WriteString(")")
	}
	if n.Having != nil {
		sb.WriteString("\n")
		indent(sb, level+1)
		sb.WriteString(":having ")
		writeNode(sb, n.Having, level+1)
	}
	if len(n.OrderBy) > 0 {
		sb.WriteString("\n")
		indent(sb, level+1)
		sb.WriteString(":order_by (")
		for i, o := range n.OrderBy {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, o, level+1)
		}
		sb.WriteString(")")
	}
	if n.OffsetLimit != nil {
		sb.WriteString("\n")
		indent(sb, level+1)
		sb.WriteString(":offset_limit ")
		writeNode(sb, n.OffsetLimit, level+1)
	}
	writeLoc(sb, n.Loc)
	sb.WriteString("}")
}
```

- [ ] **Step 2: Verify it compiles**

Run: `cd /Users/h3n4l/OpenSource/omni && go build ./cosmosdb/ast/`
Expected: Clean build.

- [ ] **Step 3: Commit**

```bash
git add cosmosdb/ast/outfuncs.go
git commit -m "feat(cosmosdb): add AST output functions for golden testing"
```

---

### Task 4: Lexer

**Files:**
- Create: `cosmosdb/parser/lexer.go`

- [ ] **Step 1: Create `cosmosdb/parser/lexer.go`**

```go
package parser

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Token type constants.
const (
	tokEOF = 0
)

// Literal and identifier tokens.
const (
	tokICONST = iota + 1000 // integer literal (DECIMAL)
	tokFCONST               // float literal (REAL, FLOAT)
	tokHCONST               // hexadecimal literal
	tokSCONST               // single-quoted string
	tokDCONST               // double-quoted string
	tokIDENT                 // identifier
	tokPARAM                 // @param
)

// Operator and punctuation tokens.
const (
	tokDOT      = iota + 2000 // .
	tokCOMMA                  // ,
	tokCOLON                  // :
	tokLPAREN                 // (
	tokRPAREN                 // )
	tokLBRACK                 // [
	tokRBRACK                 // ]
	tokLBRACE                 // {
	tokRBRACE                 // }
	tokSTAR                   // *
	tokPLUS                   // +
	tokMINUS                  // -
	tokDIV                    // /
	tokMOD                    // %
	tokEQ                     // =
	tokNE                     // !=
	tokNE2                    // <>
	tokLT                     // <
	tokLE                     // <=
	tokGT                     // >
	tokGE                     // >=
	tokBITAND                 // &
	tokBITOR                  // |
	tokBITXOR                 // ^
	tokBITNOT                 // ~
	tokLSHIFT                 // <<
	tokRSHIFT                 // >>
	tokURSHIFT                // >>>
	tokCONCAT                 // ||
	tokCOALESCE               // ??
	tokQUESTION               // ?
)

// Keyword tokens.
const (
	tokSELECT    = iota + 3000
	tokFROM
	tokWHERE
	tokAND
	tokOR
	tokNOT
	tokIN
	tokBETWEEN
	tokLIKE
	tokESCAPE
	tokAS
	tokJOIN
	tokTOP
	tokDISTINCT
	tokVALUE
	tokORDER
	tokBY
	tokGROUP
	tokHAVING
	tokOFFSET
	tokLIMIT
	tokASC
	tokDESC
	tokEXISTS
	tokTRUE
	tokFALSE
	tokNULL
	tokUNDEFINED
	tokUDF
	tokARRAY
	tokROOT
	tokRANK
	tokINFINITY // case-sensitive
	tokNAN      // case-sensitive
)

var keywords = map[string]int{
	"select":    tokSELECT,
	"from":      tokFROM,
	"where":     tokWHERE,
	"and":       tokAND,
	"or":        tokOR,
	"not":       tokNOT,
	"in":        tokIN,
	"between":   tokBETWEEN,
	"like":      tokLIKE,
	"escape":    tokESCAPE,
	"as":        tokAS,
	"join":      tokJOIN,
	"top":       tokTOP,
	"distinct":  tokDISTINCT,
	"value":     tokVALUE,
	"order":     tokORDER,
	"by":        tokBY,
	"group":     tokGROUP,
	"having":    tokHAVING,
	"offset":    tokOFFSET,
	"limit":     tokLIMIT,
	"asc":       tokASC,
	"desc":      tokDESC,
	"exists":    tokEXISTS,
	"true":      tokTRUE,
	"false":     tokFALSE,
	"null":      tokNULL,
	"undefined": tokUNDEFINED,
	"udf":       tokUDF,
	"array":     tokARRAY,
	"root":      tokROOT,
	"rank":      tokRANK,
}

// Token represents a lexical token.
type Token struct {
	Type int
	Str  string
	Loc  int // byte offset in source
}

// Lexer tokenizes CosmosDB SQL input.
type Lexer struct {
	input string
	pos   int   // current read position in input (next byte to consume)
	start int   // start position of the token currently being scanned
	Err   error
}

// NewLexer creates a new Lexer for the given input.
func NewLexer(input string) *Lexer {
	return &Lexer{input: input}
}

// NextToken returns the next token from the input.
func (l *Lexer) NextToken() Token {
	l.skipWhitespaceAndComments()
	if l.pos >= len(l.input) {
		return Token{Type: tokEOF, Loc: l.pos}
	}

	l.start = l.pos
	ch := l.input[l.pos]

	// Single-quoted string.
	if ch == '\'' {
		return l.lexString('\'', tokSCONST)
	}
	// Double-quoted string.
	if ch == '"' {
		return l.lexString('"', tokDCONST)
	}
	// Parameter reference: @ident.
	if ch == '@' {
		return l.lexParam()
	}
	// Number: starts with digit or dot-digit.
	if isDigit(ch) || (ch == '.' && l.pos+1 < len(l.input) && isDigit(l.input[l.pos+1])) {
		return l.lexNumber()
	}
	// Identifier or keyword.
	if isIdentStart(ch) {
		return l.lexIdentOrKeyword()
	}
	// Operators and punctuation.
	return l.lexOperator()
}

func (l *Lexer) skipWhitespaceAndComments() {
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n' || ch == '\f' {
			l.pos++
			continue
		}
		// Line comment: --
		if ch == '-' && l.pos+1 < len(l.input) && l.input[l.pos+1] == '-' {
			l.pos += 2
			for l.pos < len(l.input) && l.input[l.pos] != '\n' && l.input[l.pos] != '\r' {
				l.pos++
			}
			continue
		}
		break
	}
}

func (l *Lexer) lexString(quote byte, tokType int) Token {
	l.pos++ // skip opening quote
	var sb strings.Builder
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == '\\' && l.pos+1 < len(l.input) {
			l.pos++
			next := l.input[l.pos]
			switch next {
			case 'b':
				sb.WriteByte('\b')
			case 't':
				sb.WriteByte('\t')
			case 'n':
				sb.WriteByte('\n')
			case 'r':
				sb.WriteByte('\r')
			case 'f':
				sb.WriteByte('\f')
			case '"':
				sb.WriteByte('"')
			case '\'':
				sb.WriteByte('\'')
			case '\\':
				sb.WriteByte('\\')
			case '/':
				sb.WriteByte('/')
			case 'u':
				if l.pos+4 < len(l.input) {
					hex := l.input[l.pos+1 : l.pos+5]
					var r rune
					if _, err := fmt.Sscanf(hex, "%04x", &r); err == nil {
						sb.WriteRune(r)
						l.pos += 4
					} else {
						sb.WriteByte('u')
					}
				} else {
					sb.WriteByte('u')
				}
			default:
				sb.WriteByte('\\')
				sb.WriteByte(next)
			}
			l.pos++
			continue
		}
		if ch == quote {
			l.pos++ // skip closing quote
			return Token{Type: tokType, Str: sb.String(), Loc: l.start}
		}
		sb.WriteByte(ch)
		l.pos++
	}
	l.Err = fmt.Errorf("unterminated string starting at byte %d", l.start)
	return Token{Type: tokEOF, Loc: l.start}
}

func (l *Lexer) lexParam() Token {
	l.pos++ // skip @
	start := l.pos
	for l.pos < len(l.input) && isIdentPart(l.input[l.pos]) {
		l.pos++
	}
	if l.pos == start {
		l.Err = fmt.Errorf("expected identifier after @ at byte %d", l.start)
		return Token{Type: tokEOF, Loc: l.start}
	}
	return Token{Type: tokPARAM, Str: l.input[start:l.pos], Loc: l.start}
}

func (l *Lexer) lexNumber() Token {
	// Hexadecimal: 0x or 0X
	if l.input[l.pos] == '0' && l.pos+1 < len(l.input) &&
		(l.input[l.pos+1] == 'x' || l.input[l.pos+1] == 'X') {
		l.pos += 2
		for l.pos < len(l.input) && isHexDigit(l.input[l.pos]) {
			l.pos++
		}
		return Token{Type: tokHCONST, Str: l.input[l.start:l.pos], Loc: l.start}
	}

	// Integer or float part before dot.
	for l.pos < len(l.input) && isDigit(l.input[l.pos]) {
		l.pos++
	}

	isFloat := false

	// Fractional part.
	if l.pos < len(l.input) && l.input[l.pos] == '.' {
		// Only treat as float if the next char is a digit or we already consumed digits.
		if l.pos+1 < len(l.input) && isDigit(l.input[l.pos+1]) {
			isFloat = true
			l.pos++ // skip dot
			for l.pos < len(l.input) && isDigit(l.input[l.pos]) {
				l.pos++
			}
		} else if l.start == l.pos {
			// dot-only: not a number, shouldn't reach here due to caller check.
		} else {
			// Digits followed by dot but no digit after: still a float like "123."
			isFloat = true
			l.pos++ // skip dot
		}
	}

	// Exponent part: E/e [+-] digits.
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

	tok := tokICONST
	if isFloat {
		tok = tokFCONST
	}
	return Token{Type: tok, Str: l.input[l.start:l.pos], Loc: l.start}
}

func (l *Lexer) lexIdentOrKeyword() Token {
	for l.pos < len(l.input) && isIdentPart(l.input[l.pos]) {
		l.pos++
	}
	raw := l.input[l.start:l.pos]

	// Case-sensitive special constants: Infinity, NaN.
	if raw == "Infinity" {
		return Token{Type: tokINFINITY, Str: raw, Loc: l.start}
	}
	if raw == "NaN" {
		return Token{Type: tokNAN, Str: raw, Loc: l.start}
	}

	// Case-insensitive keyword lookup.
	lower := strings.ToLower(raw)
	if tok, ok := keywords[lower]; ok {
		return Token{Type: tok, Str: lower, Loc: l.start}
	}
	return Token{Type: tokIDENT, Str: raw, Loc: l.start}
}

func (l *Lexer) lexOperator() Token {
	ch := l.input[l.pos]
	l.pos++
	next := byte(0)
	if l.pos < len(l.input) {
		next = l.input[l.pos]
	}

	switch ch {
	case '.':
		return Token{Type: tokDOT, Str: ".", Loc: l.start}
	case ',':
		return Token{Type: tokCOMMA, Str: ",", Loc: l.start}
	case ':':
		return Token{Type: tokCOLON, Str: ":", Loc: l.start}
	case '(':
		return Token{Type: tokLPAREN, Str: "(", Loc: l.start}
	case ')':
		return Token{Type: tokRPAREN, Str: ")", Loc: l.start}
	case '[':
		return Token{Type: tokLBRACK, Str: "[", Loc: l.start}
	case ']':
		return Token{Type: tokRBRACK, Str: "]", Loc: l.start}
	case '{':
		return Token{Type: tokLBRACE, Str: "{", Loc: l.start}
	case '}':
		return Token{Type: tokRBRACE, Str: "}", Loc: l.start}
	case '*':
		return Token{Type: tokSTAR, Str: "*", Loc: l.start}
	case '+':
		return Token{Type: tokPLUS, Str: "+", Loc: l.start}
	case '-':
		return Token{Type: tokMINUS, Str: "-", Loc: l.start}
	case '/':
		return Token{Type: tokDIV, Str: "/", Loc: l.start}
	case '%':
		return Token{Type: tokMOD, Str: "%", Loc: l.start}
	case '=':
		return Token{Type: tokEQ, Str: "=", Loc: l.start}
	case '~':
		return Token{Type: tokBITNOT, Str: "~", Loc: l.start}
	case '^':
		return Token{Type: tokBITXOR, Str: "^", Loc: l.start}
	case '&':
		return Token{Type: tokBITAND, Str: "&", Loc: l.start}
	case '!':
		if next == '=' {
			l.pos++
			return Token{Type: tokNE, Str: "!=", Loc: l.start}
		}
		l.Err = fmt.Errorf("unexpected character '!' at byte %d", l.start)
		return Token{Type: tokEOF, Loc: l.start}
	case '<':
		if next == '=' {
			l.pos++
			return Token{Type: tokLE, Str: "<=", Loc: l.start}
		}
		if next == '>' {
			l.pos++
			return Token{Type: tokNE2, Str: "<>", Loc: l.start}
		}
		if next == '<' {
			l.pos++
			return Token{Type: tokLSHIFT, Str: "<<", Loc: l.start}
		}
		return Token{Type: tokLT, Str: "<", Loc: l.start}
	case '>':
		if next == '=' {
			l.pos++
			return Token{Type: tokGE, Str: ">=", Loc: l.start}
		}
		if next == '>' {
			l.pos++
			// Check for >>>
			if l.pos < len(l.input) && l.input[l.pos] == '>' {
				l.pos++
				return Token{Type: tokURSHIFT, Str: ">>>", Loc: l.start}
			}
			return Token{Type: tokRSHIFT, Str: ">>", Loc: l.start}
		}
		return Token{Type: tokGT, Str: ">", Loc: l.start}
	case '|':
		if next == '|' {
			l.pos++
			return Token{Type: tokCONCAT, Str: "||", Loc: l.start}
		}
		return Token{Type: tokBITOR, Str: "|", Loc: l.start}
	case '?':
		if next == '?' {
			l.pos++
			return Token{Type: tokCOALESCE, Str: "??", Loc: l.start}
		}
		return Token{Type: tokQUESTION, Str: "?", Loc: l.start}
	default:
		// Handle multi-byte UTF-8 characters gracefully.
		_, size := utf8.DecodeRuneInString(l.input[l.start:])
		if size > 1 {
			l.pos = l.start + size
		}
		l.Err = fmt.Errorf("unexpected character %q at byte %d", string(ch), l.start)
		return Token{Type: tokEOF, Loc: l.start}
	}
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isHexDigit(ch byte) bool {
	return isDigit(ch) || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
}

func isIdentStart(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

func isIdentPart(ch byte) bool {
	return isIdentStart(ch) || isDigit(ch)
}
```

- [ ] **Step 2: Verify it compiles**

Run: `cd /Users/h3n4l/OpenSource/omni && go build ./cosmosdb/parser/`
Expected: Will fail because `parser/` has no non-test `.go` file that imports it yet. Instead run:
Run: `cd /Users/h3n4l/OpenSource/omni && go vet ./cosmosdb/parser/`
Expected: May warn about unused imports; check with `go build` by creating a minimal parser.go in next task. For now, verify syntax:
Run: `cd /Users/h3n4l/OpenSource/omni && gofmt -e cosmosdb/parser/lexer.go > /dev/null`
Expected: No errors.

- [ ] **Step 3: Commit**

```bash
git add cosmosdb/parser/lexer.go
git commit -m "feat(cosmosdb): add hand-written lexer"
```

---

### Task 5: Parser Core (Entry Point and Helpers)

**Files:**
- Create: `cosmosdb/parser/parser.go`

- [ ] **Step 1: Create `cosmosdb/parser/parser.go`**

```go
package parser

import (
	"fmt"
	"strconv"

	nodes "github.com/bytebase/omni/cosmosdb/ast"
)

// Parser is a recursive descent parser for CosmosDB SQL.
type Parser struct {
	lexer   *Lexer
	cur     Token // current token (already consumed)
	prev    Token // previous token (for error context)
	nextBuf Token // buffered lookahead token
	hasNext bool  // whether nextBuf is valid
}

// ParseError represents a syntax error.
type ParseError struct {
	Message string
	Pos     int
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("syntax error at byte %d: %s", e.Pos, e.Message)
}

// Parse parses a CosmosDB SQL query and returns a list containing one RawStmt.
func Parse(sql string) (*nodes.List, error) {
	p := &Parser{lexer: NewLexer(sql)}
	p.advance()

	if p.cur.Type == tokEOF {
		return &nodes.List{}, nil
	}

	stmtStart := p.cur.Loc
	stmt, err := p.parseSelect()
	if err != nil {
		return nil, err
	}

	if p.cur.Type != tokEOF {
		return nil, &ParseError{
			Message: fmt.Sprintf("unexpected token %q after query", p.cur.Str),
			Pos:     p.cur.Loc,
		}
	}

	raw := &nodes.RawStmt{
		Stmt:         stmt,
		StmtLocation: stmtStart,
		StmtLen:      len(sql) - stmtStart,
	}

	return &nodes.List{Items: []nodes.Node{raw}}, nil
}

// parseSelect parses a full SELECT query.
func (p *Parser) parseSelect() (*nodes.SelectStmt, error) {
	startLoc := p.cur.Loc

	selectClause, err := p.parseSelectClause()
	if err != nil {
		return nil, err
	}

	stmt := selectClause

	// FROM clause (optional).
	if p.cur.Type == tokFROM {
		from, joins, err := p.parseFromClause()
		if err != nil {
			return nil, err
		}
		stmt.From = from
		stmt.Joins = joins
	}

	// WHERE clause (optional).
	if p.cur.Type == tokWHERE {
		p.advance()
		where, err := p.parseExpr(0)
		if err != nil {
			return nil, err
		}
		stmt.Where = where
	}

	// GROUP BY clause (optional).
	if p.cur.Type == tokGROUP {
		groupBy, err := p.parseGroupBy()
		if err != nil {
			return nil, err
		}
		stmt.GroupBy = groupBy
	}

	// HAVING clause (optional).
	if p.cur.Type == tokHAVING {
		p.advance()
		having, err := p.parseExpr(0)
		if err != nil {
			return nil, err
		}
		stmt.Having = having
	}

	// ORDER BY clause (optional).
	if p.cur.Type == tokORDER {
		orderBy, err := p.parseOrderBy()
		if err != nil {
			return nil, err
		}
		stmt.OrderBy = orderBy
	}

	// OFFSET ... LIMIT ... (optional).
	if p.cur.Type == tokOFFSET {
		ol, err := p.parseOffsetLimit()
		if err != nil {
			return nil, err
		}
		stmt.OffsetLimit = ol
	}

	stmt.Loc = nodes.Loc{Start: startLoc, End: p.pos()}

	return stmt, nil
}

// -----------------------------------------------------------------------
// Helper methods
// -----------------------------------------------------------------------

// advance consumes the current token and moves to the next one.
func (p *Parser) advance() Token {
	p.prev = p.cur
	if p.hasNext {
		p.cur = p.nextBuf
		p.hasNext = false
	} else {
		p.cur = p.lexer.NextToken()
	}
	return p.prev
}

// peek returns the current token without consuming it.
func (p *Parser) peek() Token {
	return p.cur
}

// peekNext returns the next token without consuming the current one.
func (p *Parser) peekNext() Token {
	if !p.hasNext {
		p.nextBuf = p.lexer.NextToken()
		p.hasNext = true
	}
	return p.nextBuf
}

// match checks if the current token matches any of the given types.
// If it matches, it consumes the token and returns it with true.
func (p *Parser) match(types ...int) (Token, bool) {
	for _, t := range types {
		if p.cur.Type == t {
			tok := p.advance()
			return tok, true
		}
	}
	return Token{}, false
}

// expect consumes the current token if it matches the given type.
// Returns an error if it doesn't match.
func (p *Parser) expect(tokenType int) (Token, error) {
	if p.cur.Type == tokenType {
		return p.advance(), nil
	}
	return Token{}, &ParseError{
		Message: fmt.Sprintf("expected token type %d, got %q", tokenType, p.cur.Str),
		Pos:     p.cur.Loc,
	}
}

// pos returns the current byte position in the source.
func (p *Parser) pos() int {
	return p.cur.Loc
}

// identifierTokens are keyword tokens that can be used as identifiers.
var identifierTokens = map[int]bool{
	tokIN:      true,
	tokBETWEEN: true,
	tokTOP:     true,
	tokVALUE:   true,
	tokORDER:   true,
	tokBY:      true,
	tokGROUP:   true,
	tokOFFSET:  true,
	tokLIMIT:   true,
	tokASC:     true,
	tokDESC:    true,
	tokEXISTS:  true,
	tokLIKE:    true,
	tokHAVING:  true,
	tokJOIN:    true,
	tokESCAPE:  true,
	tokARRAY:   true,
	tokROOT:    true,
	tokRANK:    true,
}

// propertyNameExtraTokens are additional keywords valid as property names.
var propertyNameExtraTokens = map[int]bool{
	tokSELECT:    true,
	tokFROM:      true,
	tokWHERE:     true,
	tokNOT:       true,
	tokAND:       true,
	tokOR:        true,
	tokAS:        true,
	tokTRUE:      true,
	tokFALSE:     true,
	tokNULL:      true,
	tokUNDEFINED: true,
	tokUDF:       true,
	tokDISTINCT:  true,
}

// parseIdentifier parses an identifier, accepting keywords that are valid
// in identifier position per the grammar's `identifier` rule.
func (p *Parser) parseIdentifier() (string, int, error) {
	if p.cur.Type == tokIDENT {
		tok := p.advance()
		return tok.Str, tok.Loc, nil
	}
	if identifierTokens[p.cur.Type] {
		tok := p.advance()
		return tok.Str, tok.Loc, nil
	}
	return "", 0, &ParseError{
		Message: fmt.Sprintf("expected identifier, got %q", p.cur.Str),
		Pos:     p.cur.Loc,
	}
}

// parsePropertyName parses a property name, which allows even more keywords
// than parseIdentifier per the grammar's `property_name` rule.
func (p *Parser) parsePropertyName() (string, int, error) {
	if p.cur.Type == tokIDENT || identifierTokens[p.cur.Type] || propertyNameExtraTokens[p.cur.Type] {
		tok := p.advance()
		return tok.Str, tok.Loc, nil
	}
	return "", 0, &ParseError{
		Message: fmt.Sprintf("expected property name, got %q", p.cur.Str),
		Pos:     p.cur.Loc,
	}
}

// parseIntLiteral parses a DECIMAL token as an integer.
func (p *Parser) parseIntLiteral() (int, int, error) {
	if p.cur.Type != tokICONST {
		return 0, 0, &ParseError{
			Message: fmt.Sprintf("expected integer, got %q", p.cur.Str),
			Pos:     p.cur.Loc,
		}
	}
	tok := p.advance()
	n, err := strconv.Atoi(tok.Str)
	if err != nil {
		return 0, 0, &ParseError{Message: fmt.Sprintf("invalid integer %q", tok.Str), Pos: tok.Loc}
	}
	return n, tok.Loc, nil
}
```

- [ ] **Step 2: Verify it compiles (will need stubs for parseSelectClause, parseFromClause, parseGroupBy, parseOrderBy, parseOffsetLimit, parseExpr)**

Create temporary stubs by running the build; we'll fill them in the next tasks. For now, verify the file has no syntax errors:

Run: `cd /Users/h3n4l/OpenSource/omni && gofmt -e cosmosdb/parser/parser.go > /dev/null`
Expected: No errors.

- [ ] **Step 3: Commit**

```bash
git add cosmosdb/parser/parser.go
git commit -m "feat(cosmosdb): add parser core with entry point and helpers"
```

---

### Task 6: Parser — SELECT and Clause Parsing

**Files:**
- Create: `cosmosdb/parser/select.go`

- [ ] **Step 1: Create `cosmosdb/parser/select.go`**

```go
package parser

import (
	"fmt"

	nodes "github.com/bytebase/omni/cosmosdb/ast"
)

// parseSelectClause parses: SELECT [TOP n] [DISTINCT] [VALUE] (target_list | *)
func (p *Parser) parseSelectClause() (*nodes.SelectStmt, error) {
	if _, err := p.expect(tokSELECT); err != nil {
		return nil, err
	}

	stmt := &nodes.SelectStmt{}

	// TOP n
	if p.cur.Type == tokTOP {
		p.advance()
		n, _, err := p.parseIntLiteral()
		if err != nil {
			return nil, err
		}
		stmt.Top = &n
	}

	// DISTINCT
	if p.cur.Type == tokDISTINCT {
		p.advance()
		stmt.Distinct = true
	}

	// VALUE
	if p.cur.Type == tokVALUE {
		p.advance()
		stmt.Value = true
	}

	// * or target list
	if p.cur.Type == tokSTAR {
		p.advance()
		stmt.Star = true
	} else {
		targets, err := p.parseTargetList()
		if err != nil {
			return nil, err
		}
		stmt.Targets = targets
	}

	return stmt, nil
}

// parseTargetList parses: object_property (',' object_property)*
func (p *Parser) parseTargetList() ([]*nodes.TargetEntry, error) {
	var targets []*nodes.TargetEntry

	t, err := p.parseTargetEntry()
	if err != nil {
		return nil, err
	}
	targets = append(targets, t)

	for p.cur.Type == tokCOMMA {
		p.advance()
		t, err := p.parseTargetEntry()
		if err != nil {
			return nil, err
		}
		targets = append(targets, t)
	}

	return targets, nil
}

// parseTargetEntry parses: scalar_expression [AS? property_alias]
func (p *Parser) parseTargetEntry() (*nodes.TargetEntry, error) {
	startLoc := p.cur.Loc

	expr, err := p.parseExpr(0)
	if err != nil {
		return nil, err
	}

	entry := &nodes.TargetEntry{Expr: expr}

	// Optional alias: AS? identifier
	if p.cur.Type == tokAS {
		p.advance()
		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		entry.Alias = &name
	} else if p.cur.Type == tokIDENT || identifierTokens[p.cur.Type] {
		// Alias without AS — but only if it looks like an identifier, not a keyword
		// that starts a new clause (FROM, WHERE, etc.). We check against the
		// identifierTokens set which excludes clause-starting keywords.
		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		entry.Alias = &name
	}

	entry.Loc = nodes.Loc{Start: startLoc, End: p.pos()}

	return entry, nil
}

// parseGroupBy parses: GROUP BY expr (',' expr)*
func (p *Parser) parseGroupBy() ([]nodes.ExprNode, error) {
	p.advance() // consume GROUP
	if _, err := p.expect(tokBY); err != nil {
		return nil, err
	}

	var exprs []nodes.ExprNode
	expr, err := p.parseExpr(0)
	if err != nil {
		return nil, err
	}
	exprs = append(exprs, expr)

	for p.cur.Type == tokCOMMA {
		p.advance()
		expr, err := p.parseExpr(0)
		if err != nil {
			return nil, err
		}
		exprs = append(exprs, expr)
	}

	return exprs, nil
}

// parseOrderBy parses: ORDER BY sort_expression (',' sort_expression)*
func (p *Parser) parseOrderBy() ([]*nodes.SortExpr, error) {
	p.advance() // consume ORDER
	if _, err := p.expect(tokBY); err != nil {
		return nil, err
	}

	var sorts []*nodes.SortExpr
	s, err := p.parseSortExpr()
	if err != nil {
		return nil, err
	}
	sorts = append(sorts, s)

	for p.cur.Type == tokCOMMA {
		p.advance()
		s, err := p.parseSortExpr()
		if err != nil {
			return nil, err
		}
		sorts = append(sorts, s)
	}

	return sorts, nil
}

// parseSortExpr parses: expr [ASC|DESC] | RANK expr
func (p *Parser) parseSortExpr() (*nodes.SortExpr, error) {
	startLoc := p.cur.Loc
	sort := &nodes.SortExpr{}

	if p.cur.Type == tokRANK {
		p.advance()
		sort.Rank = true
	}

	expr, err := p.parseExpr(0)
	if err != nil {
		return nil, err
	}
	sort.Expr = expr

	if !sort.Rank {
		if p.cur.Type == tokDESC {
			p.advance()
			sort.Desc = true
		} else if p.cur.Type == tokASC {
			p.advance()
		}
	}

	sort.Loc = nodes.Loc{Start: startLoc, End: p.pos()}
	return sort, nil
}

// parseOffsetLimit parses: OFFSET n LIMIT m
func (p *Parser) parseOffsetLimit() (*nodes.OffsetLimitClause, error) {
	startLoc := p.cur.Loc
	p.advance() // consume OFFSET

	offsetVal, offsetLoc, err := p.parseIntLiteral()
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(tokLIMIT); err != nil {
		return nil, err
	}

	limitVal, limitLoc, err := p.parseIntLiteral()
	if err != nil {
		return nil, err
	}

	return &nodes.OffsetLimitClause{
		Offset: &nodes.NumberLit{Val: intToStr(offsetVal), Loc: nodes.Loc{Start: offsetLoc, End: offsetLoc + len(intToStr(offsetVal))}},
		Limit:  &nodes.NumberLit{Val: intToStr(limitVal), Loc: nodes.Loc{Start: limitLoc, End: limitLoc + len(intToStr(limitVal))}},
		Loc:    nodes.Loc{Start: startLoc, End: p.pos()},
	}, nil
}

func intToStr(n int) string {
	return fmt.Sprintf("%d", n)
}
```

- [ ] **Step 2: Verify syntax**

Run: `cd /Users/h3n4l/OpenSource/omni && gofmt -e cosmosdb/parser/select.go > /dev/null`
Expected: No errors.

- [ ] **Step 3: Commit**

```bash
git add cosmosdb/parser/select.go
git commit -m "feat(cosmosdb): add SELECT clause parsing"
```

---

### Task 7: Parser — FROM Clause

**Files:**
- Create: `cosmosdb/parser/from.go`

- [ ] **Step 1: Create `cosmosdb/parser/from.go`**

```go
package parser

import (
	"fmt"

	nodes "github.com/bytebase/omni/cosmosdb/ast"
)

// parseFromClause parses: FROM from_specification
// Returns the primary source and any JOINs.
func (p *Parser) parseFromClause() (nodes.TableExpr, []*nodes.JoinExpr, error) {
	p.advance() // consume FROM

	source, err := p.parseFromSource()
	if err != nil {
		return nil, nil, err
	}

	var joins []*nodes.JoinExpr
	for p.cur.Type == tokJOIN {
		j, err := p.parseJoinClause()
		if err != nil {
			return nil, nil, err
		}
		joins = append(joins, j)
	}

	return source, joins, nil
}

// parseJoinClause parses: JOIN from_source
func (p *Parser) parseJoinClause() (*nodes.JoinExpr, error) {
	startLoc := p.cur.Loc
	p.advance() // consume JOIN

	source, err := p.parseFromSource()
	if err != nil {
		return nil, err
	}

	return &nodes.JoinExpr{
		Source: source,
		Loc:    nodes.Loc{Start: startLoc, End: p.pos()},
	}, nil
}

// parseFromSource parses:
//
//	container_expression [AS? identifier]   (aliased form)
//	identifier IN container_expression      (iteration form)
func (p *Parser) parseFromSource() (nodes.TableExpr, error) {
	startLoc := p.cur.Loc

	// Check for iteration form: identifier IN container_expression.
	// We need lookahead: if current is an identifier-like token and next is IN.
	if (p.cur.Type == tokIDENT || identifierTokens[p.cur.Type]) && p.peekNext().Type == tokIN {
		alias := p.cur.Str
		p.advance() // consume identifier
		p.advance() // consume IN

		source, err := p.parseContainerExpr()
		if err != nil {
			return nil, err
		}

		return &nodes.ArrayIterationExpr{
			Alias:  alias,
			Source: source,
			Loc:    nodes.Loc{Start: startLoc, End: p.pos()},
		}, nil
	}

	// Aliased form: container_expression [AS? identifier].
	source, err := p.parseContainerExpr()
	if err != nil {
		return nil, err
	}

	// Optional alias.
	if p.cur.Type == tokAS {
		p.advance()
		alias, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		return &nodes.AliasedTableExpr{
			Source: source,
			Alias:  alias,
			Loc:    nodes.Loc{Start: startLoc, End: p.pos()},
		}, nil
	}
	// Implicit alias (identifier not followed by a clause keyword).
	if p.cur.Type == tokIDENT || identifierTokens[p.cur.Type] {
		// Only treat as alias if it's not a keyword that starts a clause.
		if !isClauseStart(p.cur.Type) {
			alias, _, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			return &nodes.AliasedTableExpr{
				Source: source,
				Alias:  alias,
				Loc:    nodes.Loc{Start: startLoc, End: p.pos()},
			}, nil
		}
	}

	return source, nil
}

// parseContainerExpr parses:
//
//	ROOT
//	container_name
//	container_expression '.' property_name
//	container_expression '[' (string | index | @param) ']'
//	'(' select ')'
func (p *Parser) parseContainerExpr() (nodes.TableExpr, error) {
	var expr nodes.TableExpr

	switch p.cur.Type {
	case tokROOT:
		expr = &nodes.ContainerRef{
			Root: true,
			Loc:  nodes.Loc{Start: p.cur.Loc, End: p.cur.Loc + len(p.cur.Str)},
		}
		p.advance()

	case tokLPAREN:
		startLoc := p.cur.Loc
		p.advance() // consume (
		sel, err := p.parseSelect()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(tokRPAREN); err != nil {
			return nil, err
		}
		expr = &nodes.SubqueryExpr{
			Select: sel,
			Loc:    nodes.Loc{Start: startLoc, End: p.pos()},
		}

	default:
		// container_name (identifier).
		name, loc, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		expr = &nodes.ContainerRef{
			Name: name,
			Loc:  nodes.Loc{Start: loc, End: loc + len(name)},
		}
	}

	// Postfix: .property and [index]
	for {
		if p.cur.Type == tokDOT {
			startLoc := expr.(*nodes.ContainerRef)
			_ = startLoc // use the expr's own start
			dotStart := p.cur.Loc
			_ = dotStart
			p.advance() // consume .
			prop, _, err := p.parsePropertyName()
			if err != nil {
				return nil, err
			}
			expr = &nodes.DotAccessExpr{
				Expr:     exprFromTableExpr(expr),
				Property: prop,
				Loc:      nodes.Loc{Start: tableExprStart(expr), End: p.pos()},
			}
			continue
		}
		if p.cur.Type == tokLBRACK {
			p.advance() // consume [
			idx, err := p.parseBracketIndex()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(tokRBRACK); err != nil {
				return nil, err
			}
			expr = &nodes.BracketAccessExpr{
				Expr:  exprFromTableExpr(expr),
				Index: idx,
				Loc:   nodes.Loc{Start: tableExprStart(expr), End: p.pos()},
			}
			continue
		}
		break
	}

	return expr, nil
}

// parseBracketIndex parses the index inside [...]: string literal, integer, or @param.
func (p *Parser) parseBracketIndex() (nodes.ExprNode, error) {
	switch p.cur.Type {
	case tokSCONST:
		tok := p.advance()
		return &nodes.StringLit{Val: tok.Str, Loc: nodes.Loc{Start: tok.Loc, End: p.pos()}}, nil
	case tokDCONST:
		tok := p.advance()
		return &nodes.StringLit{Val: tok.Str, Loc: nodes.Loc{Start: tok.Loc, End: p.pos()}}, nil
	case tokICONST:
		tok := p.advance()
		return &nodes.NumberLit{Val: tok.Str, Loc: nodes.Loc{Start: tok.Loc, End: p.pos()}}, nil
	case tokPARAM:
		tok := p.advance()
		return &nodes.ParamRef{Name: tok.Str, Loc: nodes.Loc{Start: tok.Loc, End: p.pos()}}, nil
	default:
		return nil, &ParseError{
			Message: fmt.Sprintf("expected string, integer, or @param inside [], got %q", p.cur.Str),
			Pos:     p.cur.Loc,
		}
	}
}

// exprFromTableExpr converts a TableExpr to an ExprNode.
// DotAccessExpr and BracketAccessExpr implement both interfaces.
// ContainerRef must be converted to a ColumnRef.
func exprFromTableExpr(te nodes.TableExpr) nodes.ExprNode {
	switch n := te.(type) {
	case *nodes.DotAccessExpr:
		return n
	case *nodes.BracketAccessExpr:
		return n
	case *nodes.ContainerRef:
		return &nodes.ColumnRef{Name: n.Name, Loc: n.Loc}
	default:
		// For SubqueryExpr and others, this shouldn't happen in valid SQL.
		return nil
	}
}

// tableExprStart returns the start byte offset of a TableExpr.
func tableExprStart(te nodes.TableExpr) int {
	switch n := te.(type) {
	case *nodes.ContainerRef:
		return n.Loc.Start
	case *nodes.DotAccessExpr:
		return n.Loc.Start
	case *nodes.BracketAccessExpr:
		return n.Loc.Start
	case *nodes.SubqueryExpr:
		return n.Loc.Start
	case *nodes.AliasedTableExpr:
		return n.Loc.Start
	case *nodes.ArrayIterationExpr:
		return n.Loc.Start
	default:
		return -1
	}
}

// isClauseStart returns true for tokens that start a main clause and should
// not be consumed as implicit aliases.
func isClauseStart(tokType int) bool {
	switch tokType {
	case tokFROM, tokWHERE, tokGROUP, tokHAVING, tokORDER, tokOFFSET, tokLIMIT, tokJOIN:
		return true
	}
	return false
}
```

- [ ] **Step 2: Verify syntax**

Run: `cd /Users/h3n4l/OpenSource/omni && gofmt -e cosmosdb/parser/from.go > /dev/null`
Expected: No errors.

- [ ] **Step 3: Commit**

```bash
git add cosmosdb/parser/from.go
git commit -m "feat(cosmosdb): add FROM clause parsing"
```

---

### Task 8: Parser — Expression Parsing (Pratt Parser)

**Files:**
- Create: `cosmosdb/parser/expr.go`

- [ ] **Step 1: Create `cosmosdb/parser/expr.go`**

```go
package parser

import (
	"fmt"

	nodes "github.com/bytebase/omni/cosmosdb/ast"
)

// Precedence levels for the Pratt parser.
const (
	precNone       = 0
	precTernary    = 1  // ? :
	precCoalesce   = 2  // ??
	precOr         = 3  // OR
	precAnd        = 4  // AND
	precInBetween  = 5  // IN, BETWEEN, LIKE
	precComparison = 6  // = != <> < <= > >=
	precConcat     = 7  // ||
	precBitOr      = 8  // |
	precBitXor     = 9  // ^
	precBitAnd     = 10 // &
	precShift      = 11 // << >> >>>
	precAdditive   = 12 // + -
	precMult       = 13 // * / %
	precUnary      = 14 // NOT ~ + -
	precPostfix    = 15 // . []
)

// parseExpr parses a scalar expression using Pratt parsing.
func (p *Parser) parseExpr(minPrec int) (nodes.ExprNode, error) {
	left, err := p.parsePrefixExpr()
	if err != nil {
		return nil, err
	}

	for {
		prec, ok := p.infixPrecedence()
		if !ok || prec < minPrec {
			break
		}

		left, err = p.parseInfixExpr(left, prec)
		if err != nil {
			return nil, err
		}
	}

	return left, nil
}

// parsePrefixExpr parses a prefix expression (constants, identifiers, unary ops, etc.).
func (p *Parser) parsePrefixExpr() (nodes.ExprNode, error) {
	switch p.cur.Type {
	// Constants.
	case tokTRUE:
		tok := p.advance()
		return &nodes.BoolLit{Val: true, Loc: nodes.Loc{Start: tok.Loc, End: tok.Loc + len(tok.Str)}}, nil
	case tokFALSE:
		tok := p.advance()
		return &nodes.BoolLit{Val: false, Loc: nodes.Loc{Start: tok.Loc, End: tok.Loc + len(tok.Str)}}, nil
	case tokNULL:
		tok := p.advance()
		return &nodes.NullLit{Loc: nodes.Loc{Start: tok.Loc, End: tok.Loc + len(tok.Str)}}, nil
	case tokUNDEFINED:
		tok := p.advance()
		return &nodes.UndefinedLit{Loc: nodes.Loc{Start: tok.Loc, End: tok.Loc + len(tok.Str)}}, nil
	case tokINFINITY:
		tok := p.advance()
		return &nodes.InfinityLit{Loc: nodes.Loc{Start: tok.Loc, End: tok.Loc + len(tok.Str)}}, nil
	case tokNAN:
		tok := p.advance()
		return &nodes.NanLit{Loc: nodes.Loc{Start: tok.Loc, End: tok.Loc + len(tok.Str)}}, nil

	// Numbers.
	case tokICONST, tokFCONST, tokHCONST:
		tok := p.advance()
		return &nodes.NumberLit{Val: tok.Str, Loc: nodes.Loc{Start: tok.Loc, End: tok.Loc + len(tok.Str)}}, nil

	// Strings.
	case tokSCONST:
		tok := p.advance()
		return &nodes.StringLit{Val: tok.Str, Loc: nodes.Loc{Start: tok.Loc, End: p.pos()}}, nil
	case tokDCONST:
		tok := p.advance()
		return &nodes.StringLit{Val: tok.Str, Loc: nodes.Loc{Start: tok.Loc, End: p.pos()}}, nil

	// Parameter reference.
	case tokPARAM:
		tok := p.advance()
		return &nodes.ParamRef{Name: tok.Str, Loc: nodes.Loc{Start: tok.Loc, End: p.pos()}}, nil

	// Unary operators.
	case tokNOT:
		tok := p.advance()
		operand, err := p.parseExpr(precUnary)
		if err != nil {
			return nil, err
		}
		return &nodes.UnaryExpr{Op: "NOT", Operand: operand, Loc: nodes.Loc{Start: tok.Loc, End: p.pos()}}, nil
	case tokBITNOT:
		tok := p.advance()
		operand, err := p.parseExpr(precUnary)
		if err != nil {
			return nil, err
		}
		return &nodes.UnaryExpr{Op: "~", Operand: operand, Loc: nodes.Loc{Start: tok.Loc, End: p.pos()}}, nil
	case tokPLUS:
		tok := p.advance()
		operand, err := p.parseExpr(precUnary)
		if err != nil {
			return nil, err
		}
		return &nodes.UnaryExpr{Op: "+", Operand: operand, Loc: nodes.Loc{Start: tok.Loc, End: p.pos()}}, nil
	case tokMINUS:
		tok := p.advance()
		operand, err := p.parseExpr(precUnary)
		if err != nil {
			return nil, err
		}
		return &nodes.UnaryExpr{Op: "-", Operand: operand, Loc: nodes.Loc{Start: tok.Loc, End: p.pos()}}, nil

	// EXISTS(SELECT ...)
	case tokEXISTS:
		return p.parseExistsExpr()

	// ARRAY(SELECT ...)
	case tokARRAY:
		// Could be ARRAY(SELECT ...) or ARRAY used as an identifier.
		if p.peekNext().Type == tokLPAREN {
			return p.parseArrayExpr()
		}
		// Fall through to identifier.
		return p.parseIdentExprOrFuncCall()

	// UDF.name(...)
	case tokUDF:
		return p.parseUDFCall()

	// Parenthesized expression or subquery.
	case tokLPAREN:
		return p.parseParenExpr()

	// Array literal: [...]
	case tokLBRACK:
		return p.parseCreateArrayExpr()

	// Object literal: {...}
	case tokLBRACE:
		return p.parseCreateObjectExpr()

	// Identifier or function call.
	case tokIDENT:
		return p.parseIdentExprOrFuncCall()

	// Keywords usable as identifiers in expression position.
	default:
		if identifierTokens[p.cur.Type] {
			return p.parseIdentExprOrFuncCall()
		}
		return nil, &ParseError{
			Message: fmt.Sprintf("unexpected token %q in expression", p.cur.Str),
			Pos:     p.cur.Loc,
		}
	}
}

// infixPrecedence returns the precedence of the current token as an infix operator.
func (p *Parser) infixPrecedence() (int, bool) {
	switch p.cur.Type {
	case tokQUESTION:
		return precTernary, true
	case tokCOALESCE:
		return precCoalesce, true
	case tokOR:
		return precOr, true
	case tokAND:
		return precAnd, true
	case tokNOT:
		// NOT IN, NOT BETWEEN, NOT LIKE — only valid as infix when preceding IN/BETWEEN/LIKE.
		next := p.peekNext()
		if next.Type == tokIN || next.Type == tokBETWEEN || next.Type == tokLIKE {
			return precInBetween, true
		}
		return 0, false
	case tokIN:
		return precInBetween, true
	case tokBETWEEN:
		return precInBetween, true
	case tokLIKE:
		return precInBetween, true
	case tokEQ, tokNE, tokNE2, tokLT, tokLE, tokGT, tokGE:
		return precComparison, true
	case tokCONCAT:
		return precConcat, true
	case tokBITOR:
		return precBitOr, true
	case tokBITXOR:
		return precBitXor, true
	case tokBITAND:
		return precBitAnd, true
	case tokLSHIFT, tokRSHIFT, tokURSHIFT:
		return precShift, true
	case tokPLUS, tokMINUS:
		return precAdditive, true
	case tokSTAR, tokDIV, tokMOD:
		return precMult, true
	case tokDOT:
		return precPostfix, true
	case tokLBRACK:
		return precPostfix, true
	default:
		return 0, false
	}
}

// parseInfixExpr parses an infix expression given the left operand.
func (p *Parser) parseInfixExpr(left nodes.ExprNode, prec int) (nodes.ExprNode, error) {
	startLoc := locStart(left)

	switch p.cur.Type {
	// Ternary: left ? then : else
	case tokQUESTION:
		p.advance()
		then, err := p.parseExpr(0) // ternary is right-associative, parse fully
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(tokCOLON); err != nil {
			return nil, err
		}
		elseExpr, err := p.parseExpr(precTernary) // right-assoc
		if err != nil {
			return nil, err
		}
		return &nodes.TernaryExpr{Cond: left, Then: then, Else: elseExpr, Loc: nodes.Loc{Start: startLoc, End: p.pos()}}, nil

	// NOT IN / NOT BETWEEN / NOT LIKE
	case tokNOT:
		p.advance() // consume NOT
		switch p.cur.Type {
		case tokIN:
			return p.parseInExpr(left, true, startLoc)
		case tokBETWEEN:
			return p.parseBetweenExpr(left, true, startLoc)
		case tokLIKE:
			return p.parseLikeExpr(left, true, startLoc)
		default:
			return nil, &ParseError{
				Message: fmt.Sprintf("expected IN, BETWEEN, or LIKE after NOT, got %q", p.cur.Str),
				Pos:     p.cur.Loc,
			}
		}

	case tokIN:
		return p.parseInExpr(left, false, startLoc)
	case tokBETWEEN:
		return p.parseBetweenExpr(left, false, startLoc)
	case tokLIKE:
		return p.parseLikeExpr(left, false, startLoc)

	// Property access: left.property
	case tokDOT:
		p.advance()
		prop, _, err := p.parsePropertyName()
		if err != nil {
			return nil, err
		}
		return &nodes.DotAccessExpr{Expr: left, Property: prop, Loc: nodes.Loc{Start: startLoc, End: p.pos()}}, nil

	// Bracket access: left[index]
	case tokLBRACK:
		p.advance()
		idx, err := p.parseBracketIndex()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(tokRBRACK); err != nil {
			return nil, err
		}
		return &nodes.BracketAccessExpr{Expr: left, Index: idx, Loc: nodes.Loc{Start: startLoc, End: p.pos()}}, nil

	// Binary operators.
	default:
		op := p.cur.Str
		// Normalize keyword ops.
		switch p.cur.Type {
		case tokAND:
			op = "AND"
		case tokOR:
			op = "OR"
		}
		p.advance()
		right, err := p.parseExpr(prec + 1) // left-associative
		if err != nil {
			return nil, err
		}
		return &nodes.BinaryExpr{Op: op, Left: left, Right: right, Loc: nodes.Loc{Start: startLoc, End: p.pos()}}, nil
	}
}

// parseInExpr parses: [NOT] IN (expr, ...)
func (p *Parser) parseInExpr(left nodes.ExprNode, not bool, startLoc int) (nodes.ExprNode, error) {
	p.advance() // consume IN
	if _, err := p.expect(tokLPAREN); err != nil {
		return nil, err
	}

	var list []nodes.ExprNode
	if p.cur.Type != tokRPAREN {
		expr, err := p.parseExpr(0)
		if err != nil {
			return nil, err
		}
		list = append(list, expr)
		for p.cur.Type == tokCOMMA {
			p.advance()
			expr, err := p.parseExpr(0)
			if err != nil {
				return nil, err
			}
			list = append(list, expr)
		}
	}

	if _, err := p.expect(tokRPAREN); err != nil {
		return nil, err
	}

	return &nodes.InExpr{Expr: left, List: list, Not: not, Loc: nodes.Loc{Start: startLoc, End: p.pos()}}, nil
}

// parseBetweenExpr parses: [NOT] BETWEEN low AND high
func (p *Parser) parseBetweenExpr(left nodes.ExprNode, not bool, startLoc int) (nodes.ExprNode, error) {
	p.advance() // consume BETWEEN
	low, err := p.parseExpr(precComparison) // parse up to comparison level to avoid consuming AND
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(tokAND); err != nil {
		return nil, err
	}
	high, err := p.parseExpr(precComparison)
	if err != nil {
		return nil, err
	}
	return &nodes.BetweenExpr{Expr: left, Low: low, High: high, Not: not, Loc: nodes.Loc{Start: startLoc, End: p.pos()}}, nil
}

// parseLikeExpr parses: [NOT] LIKE pattern [ESCAPE escape]
func (p *Parser) parseLikeExpr(left nodes.ExprNode, not bool, startLoc int) (nodes.ExprNode, error) {
	p.advance() // consume LIKE
	pattern, err := p.parseExpr(precComparison)
	if err != nil {
		return nil, err
	}

	var escape nodes.ExprNode
	if p.cur.Type == tokESCAPE {
		p.advance()
		escape, err = p.parseExpr(precComparison)
		if err != nil {
			return nil, err
		}
	}

	return &nodes.LikeExpr{Expr: left, Pattern: pattern, Escape: escape, Not: not, Loc: nodes.Loc{Start: startLoc, End: p.pos()}}, nil
}

// parseExistsExpr parses: EXISTS(SELECT ...)
func (p *Parser) parseExistsExpr() (nodes.ExprNode, error) {
	startLoc := p.cur.Loc
	p.advance() // consume EXISTS
	if _, err := p.expect(tokLPAREN); err != nil {
		return nil, err
	}
	sel, err := p.parseSelect()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(tokRPAREN); err != nil {
		return nil, err
	}
	return &nodes.ExistsExpr{Select: sel, Loc: nodes.Loc{Start: startLoc, End: p.pos()}}, nil
}

// parseArrayExpr parses: ARRAY(SELECT ...)
func (p *Parser) parseArrayExpr() (nodes.ExprNode, error) {
	startLoc := p.cur.Loc
	p.advance() // consume ARRAY
	if _, err := p.expect(tokLPAREN); err != nil {
		return nil, err
	}
	sel, err := p.parseSelect()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(tokRPAREN); err != nil {
		return nil, err
	}
	return &nodes.ArrayExpr{Select: sel, Loc: nodes.Loc{Start: startLoc, End: p.pos()}}, nil
}

// parseParenExpr parses: '(' expr ')' or '(' SELECT ... ')'
func (p *Parser) parseParenExpr() (nodes.ExprNode, error) {
	startLoc := p.cur.Loc
	p.advance() // consume (

	if p.cur.Type == tokSELECT {
		sel, err := p.parseSelect()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(tokRPAREN); err != nil {
			return nil, err
		}
		return &nodes.SubLink{Select: sel, Loc: nodes.Loc{Start: startLoc, End: p.pos()}}, nil
	}

	expr, err := p.parseExpr(0)
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(tokRPAREN); err != nil {
		return nil, err
	}
	return expr, nil
}

// parseCreateArrayExpr parses: '[' [expr (',' expr)*] ']'
func (p *Parser) parseCreateArrayExpr() (nodes.ExprNode, error) {
	startLoc := p.cur.Loc
	p.advance() // consume [

	var elements []nodes.ExprNode
	if p.cur.Type != tokRBRACK {
		elem, err := p.parseExpr(0)
		if err != nil {
			return nil, err
		}
		elements = append(elements, elem)
		for p.cur.Type == tokCOMMA {
			p.advance()
			elem, err := p.parseExpr(0)
			if err != nil {
				return nil, err
			}
			elements = append(elements, elem)
		}
	}

	if _, err := p.expect(tokRBRACK); err != nil {
		return nil, err
	}

	return &nodes.CreateArrayExpr{Elements: elements, Loc: nodes.Loc{Start: startLoc, End: p.pos()}}, nil
}

// parseCreateObjectExpr parses: '{' [field_pair (',' field_pair)*] '}'
func (p *Parser) parseCreateObjectExpr() (nodes.ExprNode, error) {
	startLoc := p.cur.Loc
	p.advance() // consume {

	var fields []*nodes.ObjectFieldPair
	if p.cur.Type != tokRBRACE {
		f, err := p.parseObjectFieldPair()
		if err != nil {
			return nil, err
		}
		fields = append(fields, f)
		for p.cur.Type == tokCOMMA {
			p.advance()
			f, err := p.parseObjectFieldPair()
			if err != nil {
				return nil, err
			}
			fields = append(fields, f)
		}
	}

	if _, err := p.expect(tokRBRACE); err != nil {
		return nil, err
	}

	return &nodes.CreateObjectExpr{Fields: fields, Loc: nodes.Loc{Start: startLoc, End: p.pos()}}, nil
}

// parseObjectFieldPair parses: (string_literal | property_name) ':' scalar_expression
func (p *Parser) parseObjectFieldPair() (*nodes.ObjectFieldPair, error) {
	startLoc := p.cur.Loc
	var key string

	if p.cur.Type == tokSCONST || p.cur.Type == tokDCONST {
		tok := p.advance()
		key = tok.Str
	} else {
		name, _, err := p.parsePropertyName()
		if err != nil {
			return nil, err
		}
		key = name
	}

	if _, err := p.expect(tokCOLON); err != nil {
		return nil, err
	}

	value, err := p.parseExpr(0)
	if err != nil {
		return nil, err
	}

	return &nodes.ObjectFieldPair{Key: key, Value: value, Loc: nodes.Loc{Start: startLoc, End: p.pos()}}, nil
}

// parseIdentExprOrFuncCall parses an identifier that could be:
//   - a simple column reference
//   - a function call: name(args)
func (p *Parser) parseIdentExprOrFuncCall() (nodes.ExprNode, error) {
	name, loc, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	// Function call: name(...)
	if p.cur.Type == tokLPAREN {
		return p.parseFuncCall(name, loc)
	}

	return &nodes.ColumnRef{Name: name, Loc: nodes.Loc{Start: loc, End: loc + len(name)}}, nil
}

// locStart returns the start byte offset of an ExprNode.
func locStart(expr nodes.ExprNode) int {
	switch n := expr.(type) {
	case *nodes.ColumnRef:
		return n.Loc.Start
	case *nodes.DotAccessExpr:
		return n.Loc.Start
	case *nodes.BracketAccessExpr:
		return n.Loc.Start
	case *nodes.BinaryExpr:
		return n.Loc.Start
	case *nodes.UnaryExpr:
		return n.Loc.Start
	case *nodes.TernaryExpr:
		return n.Loc.Start
	case *nodes.InExpr:
		return n.Loc.Start
	case *nodes.BetweenExpr:
		return n.Loc.Start
	case *nodes.LikeExpr:
		return n.Loc.Start
	case *nodes.FuncCall:
		return n.Loc.Start
	case *nodes.UDFCall:
		return n.Loc.Start
	case *nodes.ExistsExpr:
		return n.Loc.Start
	case *nodes.ArrayExpr:
		return n.Loc.Start
	case *nodes.CreateArrayExpr:
		return n.Loc.Start
	case *nodes.CreateObjectExpr:
		return n.Loc.Start
	case *nodes.ParamRef:
		return n.Loc.Start
	case *nodes.SubLink:
		return n.Loc.Start
	case *nodes.StringLit:
		return n.Loc.Start
	case *nodes.NumberLit:
		return n.Loc.Start
	case *nodes.BoolLit:
		return n.Loc.Start
	case *nodes.NullLit:
		return n.Loc.Start
	case *nodes.UndefinedLit:
		return n.Loc.Start
	case *nodes.InfinityLit:
		return n.Loc.Start
	case *nodes.NanLit:
		return n.Loc.Start
	default:
		return -1
	}
}
```

- [ ] **Step 2: Verify syntax**

Run: `cd /Users/h3n4l/OpenSource/omni && gofmt -e cosmosdb/parser/expr.go > /dev/null`
Expected: No errors.

- [ ] **Step 3: Commit**

```bash
git add cosmosdb/parser/expr.go
git commit -m "feat(cosmosdb): add Pratt expression parser"
```

---

### Task 9: Parser — Function Calls

**Files:**
- Create: `cosmosdb/parser/function.go`

- [ ] **Step 1: Create `cosmosdb/parser/function.go`**

```go
package parser

import (
	nodes "github.com/bytebase/omni/cosmosdb/ast"
)

// parseFuncCall parses a built-in function call: name([*|args]).
// name and loc have already been consumed by the caller.
func (p *Parser) parseFuncCall(name string, startLoc int) (nodes.ExprNode, error) {
	p.advance() // consume (

	fc := &nodes.FuncCall{Name: name}

	if p.cur.Type == tokSTAR {
		p.advance()
		fc.Star = true
	} else if p.cur.Type != tokRPAREN {
		arg, err := p.parseExpr(0)
		if err != nil {
			return nil, err
		}
		fc.Args = append(fc.Args, arg)
		for p.cur.Type == tokCOMMA {
			p.advance()
			arg, err := p.parseExpr(0)
			if err != nil {
				return nil, err
			}
			fc.Args = append(fc.Args, arg)
		}
	}

	if _, err := p.expect(tokRPAREN); err != nil {
		return nil, err
	}

	fc.Loc = nodes.Loc{Start: startLoc, End: p.pos()}
	return fc, nil
}

// parseUDFCall parses: UDF.name(args)
func (p *Parser) parseUDFCall() (nodes.ExprNode, error) {
	startLoc := p.cur.Loc
	p.advance() // consume UDF

	if _, err := p.expect(tokDOT); err != nil {
		return nil, err
	}

	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(tokLPAREN); err != nil {
		return nil, err
	}

	var args []nodes.ExprNode
	if p.cur.Type != tokRPAREN {
		arg, err := p.parseExpr(0)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
		for p.cur.Type == tokCOMMA {
			p.advance()
			arg, err := p.parseExpr(0)
			if err != nil {
				return nil, err
			}
			args = append(args, arg)
		}
	}

	if _, err := p.expect(tokRPAREN); err != nil {
		return nil, err
	}

	return &nodes.UDFCall{Name: name, Args: args, Loc: nodes.Loc{Start: startLoc, End: p.pos()}}, nil
}
```

- [ ] **Step 2: Verify the entire parser package compiles**

Run: `cd /Users/h3n4l/OpenSource/omni && go build ./cosmosdb/parser/`
Expected: Clean build. All files compile together now.

- [ ] **Step 3: Commit**

```bash
git add cosmosdb/parser/function.go
git commit -m "feat(cosmosdb): add function call parsing"
```

---

### Task 10: Public API

**Files:**
- Create: `cosmosdb/parse.go`

- [ ] **Step 1: Create `cosmosdb/parse.go`**

```go
// Package cosmosdb provides a parser for Azure Cosmos DB NoSQL SQL API queries.
package cosmosdb

import (
	"sort"

	"github.com/bytebase/omni/cosmosdb/ast"
	"github.com/bytebase/omni/cosmosdb/parser"
)

// Statement represents a single parsed SQL statement with its source location.
type Statement struct {
	Text      string   // original SQL text
	AST       ast.Node // inner statement node (*ast.SelectStmt), nil for empty
	ByteStart int      // inclusive start byte offset
	ByteEnd   int      // exclusive end byte offset
	Start     Position // start position (line:column, 1-based)
	End       Position // end position (line:column, 1-based)
}

// Position represents a line:column position in source text.
type Position struct {
	Line   int // 1-based line number
	Column int // 1-based column in bytes
}

// Empty returns true if this statement has no AST (e.g., whitespace-only).
func (s Statement) Empty() bool {
	return s.AST == nil
}

// Parse parses a CosmosDB SQL query and returns a list of statements.
// CosmosDB SQL typically contains a single SELECT statement.
func Parse(sql string) ([]Statement, error) {
	list, err := parser.Parse(sql)
	if err != nil {
		return nil, err
	}

	if list.Len() == 0 {
		return nil, nil
	}

	idx := buildLineIndex(sql)

	var stmts []Statement
	for _, item := range list.Items {
		raw, ok := item.(*ast.RawStmt)
		if !ok {
			continue
		}

		byteStart := raw.StmtLocation
		byteEnd := byteStart + raw.StmtLen

		// Trim trailing whitespace.
		for byteEnd > byteStart && isSpace(sql[byteEnd-1]) {
			byteEnd--
		}

		stmts = append(stmts, Statement{
			Text:      sql[byteStart:byteEnd],
			AST:       raw.Stmt,
			ByteStart: byteStart,
			ByteEnd:   byteEnd,
			Start:     offsetToPosition(idx, byteStart),
			End:       offsetToPosition(idx, byteEnd),
		})
	}

	return stmts, nil
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\f'
}

// lineIndex stores byte offsets for the start of each line.
type lineIndex []int

func buildLineIndex(s string) lineIndex {
	idx := lineIndex{0}
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			idx = append(idx, i+1)
		}
	}
	return idx
}

func offsetToPosition(idx lineIndex, offset int) Position {
	line := sort.SearchInts(idx, offset+1)
	col := offset - idx[line-1] + 1
	return Position{Line: line, Column: col}
}
```

- [ ] **Step 2: Verify it compiles**

Run: `cd /Users/h3n4l/OpenSource/omni && go build ./cosmosdb/`
Expected: Clean build.

- [ ] **Step 3: Commit**

```bash
git add cosmosdb/parse.go
git commit -m "feat(cosmosdb): add public Parse API"
```

---

### Task 11: Add Test Data Files

**Files:**
- Create: `cosmosdb/parser/testdata/*.sql` (37 files from bytebase/parser/cosmosdb/examples/)

- [ ] **Step 1: Create testdata directory and download all 37 example SQL files**

Create each file in `cosmosdb/parser/testdata/`. The files are fetched from `https://raw.githubusercontent.com/bytebase/parser/main/cosmosdb/examples/`. Create all 37 files:

```bash
mkdir -p /Users/h3n4l/OpenSource/omni/cosmosdb/parser/testdata
cd /Users/h3n4l/OpenSource/omni/cosmosdb/parser/testdata
# Download all 37 example files.
for f in \
  simple_select.sql select_top.sql select_functions.sql \
  from_root.sql from_subroot.sql from_subquery.sql from_in_iteration.sql \
  where_expression.sql where_operators.sql between_expression.sql \
  in_expression.sql not_equal.sql not_equal_diamond.sql \
  math_functions.sql string_functions.sql type_check_functions.sql \
  aggregation.sql group_by.sql order_by.sql offset_limit.sql \
  distinct_value.sql join_subquery.sql array_subquery.sql \
  fulltext_search.sql geospatial.sql like_escape.sql \
  coalesce_operator.sql bracket_parameter.sql \
  property_field_name_project.sql property_array_number_project.sql \
  property_field_name_ls_bracket_project.sql line_comment.sql \
  keyword_property_names.sql underscore_fields.sql infinity_nan.sql \
  value_count.sql value_keyword.sql; do
  curl -sL "https://raw.githubusercontent.com/bytebase/parser/main/cosmosdb/examples/$f" -o "$f"
done
```

- [ ] **Step 2: Verify files were downloaded**

Run: `ls /Users/h3n4l/OpenSource/omni/cosmosdb/parser/testdata/*.sql | wc -l`
Expected: `37`

- [ ] **Step 3: Commit**

```bash
git add cosmosdb/parser/testdata/
git commit -m "feat(cosmosdb): add 37 example SQL test files"
```

---

### Task 12: Parser Tests

**Files:**
- Create: `cosmosdb/parser/parser_test.go`

- [ ] **Step 1: Create `cosmosdb/parser/parser_test.go`**

```go
package parser_test

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytebase/omni/cosmosdb/ast"
	"github.com/bytebase/omni/cosmosdb/parser"
)

var update = flag.Bool("update", false, "update golden files")

func TestParseExamples(t *testing.T) {
	files, err := filepath.Glob(filepath.Join("testdata", "*.sql"))
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no test SQL files found in testdata/")
	}

	goldenDir := filepath.Join("testdata", "golden")
	if *update {
		if err := os.MkdirAll(goldenDir, 0o755); err != nil {
			t.Fatal(err)
		}
	}

	for _, file := range files {
		name := filepath.Base(file)
		t.Run(name, func(t *testing.T) {
			data, err := os.ReadFile(file)
			if err != nil {
				t.Fatal(err)
			}
			sql := string(data)

			list, err := parser.Parse(sql)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			got := ast.NodeToString(list)

			goldenFile := filepath.Join(goldenDir, strings.TrimSuffix(name, ".sql")+".txt")

			if *update {
				if err := os.WriteFile(goldenFile, []byte(got+"\n"), 0o644); err != nil {
					t.Fatal(err)
				}
				return
			}

			want, err := os.ReadFile(goldenFile)
			if err != nil {
				t.Fatalf("Golden file not found: %s (run with -update to create)", goldenFile)
			}

			if got+"\n" != string(want) {
				t.Errorf("AST output mismatch for %s.\nGot:\n%s\nWant:\n%s", name, got, string(want))
			}
		})
	}
}
```

- [ ] **Step 2: Generate golden files**

Run: `cd /Users/h3n4l/OpenSource/omni && go test ./cosmosdb/parser/ -run TestParseExamples -update -v`
Expected: All 37 tests pass. Golden files created in `testdata/golden/`.

- [ ] **Step 3: Run tests without -update to verify golden file matching**

Run: `cd /Users/h3n4l/OpenSource/omni && go test ./cosmosdb/parser/ -run TestParseExamples -v`
Expected: All 37 tests pass.

- [ ] **Step 4: Commit**

```bash
git add cosmosdb/parser/parser_test.go cosmosdb/parser/testdata/golden/
git commit -m "test(cosmosdb): add parser golden-file tests for all 37 examples"
```

---

### Task 13: Public API Tests

**Files:**
- Create: `cosmosdb/parse_test.go`

- [ ] **Step 1: Create `cosmosdb/parse_test.go`**

```go
package cosmosdb_test

import (
	"testing"

	"github.com/bytebase/omni/cosmosdb"
	"github.com/bytebase/omni/cosmosdb/ast"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		wantCount int
		check     func(t *testing.T, stmts []cosmosdb.Statement)
	}{
		{
			name:      "simple select star",
			sql:       "SELECT * FROM c",
			wantCount: 1,
			check: func(t *testing.T, stmts []cosmosdb.Statement) {
				s := stmts[0]
				if s.Text != "SELECT * FROM c" {
					t.Errorf("Text = %q, want %q", s.Text, "SELECT * FROM c")
				}
				if _, ok := s.AST.(*ast.SelectStmt); !ok {
					t.Errorf("AST type = %T, want *ast.SelectStmt", s.AST)
				}
				if s.ByteStart != 0 {
					t.Errorf("ByteStart = %d, want 0", s.ByteStart)
				}
				if s.ByteEnd != 15 {
					t.Errorf("ByteEnd = %d, want 15", s.ByteEnd)
				}
				if s.Start.Line != 1 || s.Start.Column != 1 {
					t.Errorf("Start = %v, want {1 1}", s.Start)
				}
			},
		},
		{
			name:      "multiline",
			sql:       "SELECT\n  *\nFROM\n  c",
			wantCount: 1,
			check: func(t *testing.T, stmts []cosmosdb.Statement) {
				s := stmts[0]
				if s.Start.Line != 1 {
					t.Errorf("Start.Line = %d, want 1", s.Start.Line)
				}
				if s.End.Line != 4 {
					t.Errorf("End.Line = %d, want 4", s.End.Line)
				}
			},
		},
		{
			name:      "with trailing whitespace",
			sql:       "SELECT * FROM c   ",
			wantCount: 1,
			check: func(t *testing.T, stmts []cosmosdb.Statement) {
				s := stmts[0]
				if s.Text != "SELECT * FROM c" {
					t.Errorf("Text = %q, want %q", s.Text, "SELECT * FROM c")
				}
				if s.ByteEnd != 15 {
					t.Errorf("ByteEnd = %d, want 15", s.ByteEnd)
				}
			},
		},
		{
			name:      "empty input",
			sql:       "",
			wantCount: 0,
		},
		{
			name:      "whitespace only",
			sql:       "   \n\t  ",
			wantCount: 0,
		},
		{
			name:      "select with value",
			sql:       "SELECT VALUE c.name FROM c",
			wantCount: 1,
			check: func(t *testing.T, stmts []cosmosdb.Statement) {
				sel := stmts[0].AST.(*ast.SelectStmt)
				if !sel.Value {
					t.Error("Value should be true")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmts, err := cosmosdb.Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}
			if len(stmts) != tt.wantCount {
				t.Fatalf("got %d statements, want %d", len(stmts), tt.wantCount)
			}
			if tt.check != nil {
				tt.check(t, stmts)
			}
		})
	}
}
```

- [ ] **Step 2: Run the tests**

Run: `cd /Users/h3n4l/OpenSource/omni && go test ./cosmosdb/ -run TestParse -v`
Expected: All tests pass.

- [ ] **Step 3: Commit**

```bash
git add cosmosdb/parse_test.go
git commit -m "test(cosmosdb): add public API tests"
```

---

### Task 14: Integration Verification

- [ ] **Step 1: Run all CosmosDB tests together**

Run: `cd /Users/h3n4l/OpenSource/omni && go test ./cosmosdb/... -v`
Expected: All tests pass (parser golden tests + public API tests).

- [ ] **Step 2: Run go vet on the new package**

Run: `cd /Users/h3n4l/OpenSource/omni && go vet ./cosmosdb/...`
Expected: No issues.

- [ ] **Step 3: Verify no impact on existing engines**

Run: `cd /Users/h3n4l/OpenSource/omni && go build ./...`
Expected: Clean build for all packages.

- [ ] **Step 4: Fix any issues found, then commit fixes if needed**

If all clean, no commit needed for this task.
