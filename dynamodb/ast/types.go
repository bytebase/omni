package ast

// ---------------------------------------------------------------------------
// PartiQL type references — implement TypeName.
//
// One flat TypeRef covers all PartiQL types. Discrimination is by Name
// (uppercase canonical: INT, DECIMAL, TIME, BAG, ANY, ...) rather than a
// 25-arm Go enum. The grammar treats types as a name token with optional
// modifiers (Args for parameterized types, WithTimeZone for TIME/TIMESTAMP),
// so the AST mirrors that shape.
//
// Grammar: PartiQLParser.g4 — sourced from rules:
//   type#TypeAtomic    (lines 675–680)
//   type#TypeArgSingle (line 681)
//   type#TypeVarChar   (line 682)
//   type#TypeArgDouble (line 683)
//   type#TypeTimeZone  (line 684)
//   type#TypeCustom    (line 685)
// ---------------------------------------------------------------------------

// TypeRef represents any PartiQL type reference, used by CAST and DDL.
//
// Args carries optional precision/scale/length, with positional indexing:
//   - Args[0] = precision / length
//   - Args[1] = scale (DECIMAL/NUMERIC only)
//
// Examples:
//   - DECIMAL(10,2)         -> Args=[10, 2]
//   - VARCHAR(255)          -> Args=[255]
//   - FLOAT(53)             -> Args=[53]
//   - TIME(6) WITH TIME ZONE -> Args=[6], WithTimeZone=true
//   - INT                   -> Args=nil
//
// WithTimeZone is set for `TIME WITH TIME ZONE`. (Note: the legacy
// PartiQLParser.g4 grammar only carries the WITH TIME ZONE form on the
// `TypeTimeZone` alternative for TIME, not TIMESTAMP — see line 684.)
//
// Names covered (canonical uppercase, drawn from PartiQLParser.g4 lines
// 676–685):
//   - Numeric: INT, INT2, INT4, INT8, INTEGER, INTEGER2, INTEGER4, INTEGER8,
//     SMALLINT, BIGINT, REAL, DOUBLE PRECISION, FLOAT, DECIMAL, DEC, NUMERIC
//   - Boolean: BOOL, BOOLEAN
//   - Null / missing: NULL, MISSING
//   - String: STRING, SYMBOL, VARCHAR, CHAR, CHARACTER, "CHARACTER VARYING"
//   - Large objects: BLOB, CLOB
//   - Temporal: DATE, TIME, TIMESTAMP
//   - Collection: STRUCT, TUPLE, LIST, BAG, SEXP
//   - Any: ANY
//   - User-defined: any symbolPrimitive identifier (the `TypeCustom`
//     alternative on line 685)
type TypeRef struct {
	Name         string
	Args         []int
	WithTimeZone bool
	Loc          Loc
}

func (*TypeRef) nodeTag()      {}
func (n *TypeRef) GetLoc() Loc { return n.Loc }
func (*TypeRef) typeName()     {}
