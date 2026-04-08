package ast

// ---------------------------------------------------------------------------
// Literal nodes — all implement ExprNode.
//
// Grammar: PartiQLParser.g4 lines 661–672 (the `literal` rule). Each type
// here corresponds to one labeled alternative of that rule. There is no
// TIMESTAMP literal — TIMESTAMP appears only as a type-name keyword in the
// `type` rule (line 677), used by CAST and DDL.
// ---------------------------------------------------------------------------

// StringLit represents a single-quoted string literal: 'hello'.
//
// Grammar: literal#LiteralString (LITERAL_STRING)
type StringLit struct {
	Val string // SQL escapes already decoded ('' → ')
	Loc Loc
}

func (*StringLit) nodeTag()      {}
func (n *StringLit) GetLoc() Loc { return n.Loc }
func (*StringLit) exprNode()     {}

// NumberLit represents a numeric literal. Val stores the raw text to
// preserve the original representation (integer vs decimal vs scientific).
// Consumers needing a typed value should call strconv.ParseFloat or
// shopspring/decimal at the call site — this AST does not normalize.
//
// Grammar: literal#LiteralInteger / literal#LiteralDecimal
type NumberLit struct {
	Val string // raw text as it appears in source
	Loc Loc
}

func (*NumberLit) nodeTag()      {}
func (n *NumberLit) GetLoc() Loc { return n.Loc }
func (*NumberLit) exprNode()     {}

// BoolLit represents TRUE or FALSE.
//
// Grammar: literal#LiteralTrue / literal#LiteralFalse
type BoolLit struct {
	Val bool
	Loc Loc
}

func (*BoolLit) nodeTag()      {}
func (n *BoolLit) GetLoc() Loc { return n.Loc }
func (*BoolLit) exprNode()     {}

// NullLit represents NULL.
//
// Grammar: literal#LiteralNull
type NullLit struct {
	Loc Loc
}

func (*NullLit) nodeTag()      {}
func (n *NullLit) GetLoc() Loc { return n.Loc }
func (*NullLit) exprNode()     {}

// MissingLit represents the PartiQL-distinct MISSING value.
//
// Grammar: literal#LiteralMissing
type MissingLit struct {
	Loc Loc
}

func (*MissingLit) nodeTag()      {}
func (n *MissingLit) GetLoc() Loc { return n.Loc }
func (*MissingLit) exprNode()     {}

// DateLit represents `DATE 'YYYY-MM-DD'`.
//
// Grammar: literal#LiteralDate (DATE LITERAL_STRING)
type DateLit struct {
	Val string // YYYY-MM-DD body
	Loc Loc
}

func (*DateLit) nodeTag()      {}
func (n *DateLit) GetLoc() Loc { return n.Loc }
func (*DateLit) exprNode()     {}

// TimeLit represents `TIME [(p)] [WITH TIME ZONE] 'HH:MM:SS[.frac]'`.
//
// Grammar: literal#LiteralTime
//
//	TIME (PAREN_LEFT LITERAL_INTEGER PAREN_RIGHT)? (WITH TIME ZONE)? LITERAL_STRING
type TimeLit struct {
	Val          string // HH:MM:SS[.frac] body
	Precision    *int   // optional fractional-seconds precision
	WithTimeZone bool   // WITH TIME ZONE clause present
	Loc          Loc
}

func (*TimeLit) nodeTag()      {}
func (n *TimeLit) GetLoc() Loc { return n.Loc }
func (*TimeLit) exprNode()     {}

// IonLit represents a backtick-delimited inline Ion value: `…`.
// Text holds the verbatim contents between the backticks (no parsing,
// no normalization). PartiQL-unique.
//
// Grammar: literal#LiteralIon (ION_CLOSURE)
type IonLit struct {
	Text string
	Loc  Loc
}

func (*IonLit) nodeTag()      {}
func (n *IonLit) GetLoc() Loc { return n.Loc }
func (*IonLit) exprNode()     {}
