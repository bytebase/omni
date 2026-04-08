package ast

// ---------------------------------------------------------------------------
// Graph Pattern Matching (GPML) nodes — implement PatternNode (or ExprNode
// for the top-level MatchExpr container).
//
// PartiQL has a graph pattern matching extension based on GPML. The full
// shape of the field set inside NodePattern/EdgePattern may be refined when
// the parser-graph DAG node (node 16) is implemented and we read the grammar
// more carefully. The marker interface, the node names, and the multi-interface
// relationship to ExprNode are stable in this initial pass.
//
// Grammar: PartiQLParser.g4 — sourced from rules:
//   gpmlPattern          (lines 315–316)
//   gpmlPatternList      (lines 318–319)
//   matchPattern         (lines 321–322)
//   graphPart            (lines 324–328)
//   matchSelector        (lines 330–334)
//   patternPathVariable  (lines 336–337)
//   patternRestrictor    (lines 339–340)
//   node                 (lines 342–343)
//   edge                 (lines 345–348)
//   pattern              (lines 350–353)
//   patternQuantifier    (lines 355–358)
//   edgeWSpec            (lines 360–368)
//   edgeSpec             (lines 370–371)
//   patternPartLabel     (lines 373–374)
//   edgeAbbrev           (lines 376–381)
//   exprGraphMatchOne    (lines 628–629)
//   exprGraphMatchMany   (lines 625–626)
// Each type below cites its specific rule#Label.
// ---------------------------------------------------------------------------

// EdgeDirection identifies the direction of an edge pattern.
type EdgeDirection int

// EdgeDirection enumerates the seven edge alternatives from PartiQLParser.g4
// `edgeWSpec` (lines 360–368) and `edgeAbbrev` (lines 376–381). Each value
// maps to the bare grammar symbol via String().
const (
	EdgeDirInvalid              EdgeDirection = iota
	EdgeDirRight                              // ->          edgeWSpec#EdgeSpecRight
	EdgeDirLeft                               // <-          edgeWSpec#EdgeSpecLeft
	EdgeDirUndirected                         // ~           edgeWSpec#EdgeSpecUndirected
	EdgeDirLeftOrRight                        // <->         edgeWSpec#EdgeSpecBidirectional
	EdgeDirLeftOrUndirected                   // <~          edgeWSpec#EdgeSpecUndirectedLeft
	EdgeDirRightOrUndirected                  // ~>          edgeWSpec#EdgeSpecUndirectedRight
	EdgeDirUndirectedBidirectional            // -           edgeWSpec#EdgeSpecUndirectedBidirectional (any of right/left/undirected)
)

func (d EdgeDirection) String() string {
	switch d {
	case EdgeDirRight:
		return "->"
	case EdgeDirLeft:
		return "<-"
	case EdgeDirUndirected:
		return "~"
	case EdgeDirLeftOrRight:
		return "<->"
	case EdgeDirLeftOrUndirected:
		return "<~"
	case EdgeDirRightOrUndirected:
		return "~>"
	case EdgeDirUndirectedBidirectional:
		return "-"
	default:
		return "INVALID"
	}
}

// PatternRestrictor identifies the optional restrictor keyword on a graph pattern.
type PatternRestrictor int

const (
	PatternRestrictorNone PatternRestrictor = iota
	PatternRestrictorTrail
	PatternRestrictorAcyclic
	PatternRestrictorSimple
)

func (r PatternRestrictor) String() string {
	switch r {
	case PatternRestrictorTrail:
		return "TRAIL"
	case PatternRestrictorAcyclic:
		return "ACYCLIC"
	case PatternRestrictorSimple:
		return "SIMPLE"
	default:
		return ""
	}
}

// SelectorKind identifies the optional selector keyword on a graph pattern.
type SelectorKind int

const (
	SelectorKindNone SelectorKind = iota
	SelectorKindAny
	SelectorKindAllShortest
	SelectorKindShortestK
)

func (s SelectorKind) String() string {
	switch s {
	case SelectorKindAny:
		return "ANY"
	case SelectorKindAllShortest:
		return "ALL_SHORTEST"
	case SelectorKindShortestK:
		return "SHORTEST_K"
	default:
		return ""
	}
}

// MatchExpr is the top-level graph-match expression: MATCH(graph_expr, pattern, …).
// Implements ExprNode because PartiQL embeds graph matching in expression
// position (exprGraphMatchOne / exprGraphMatchMany rules).
//
// Grammar: exprGraphMatchOne, exprGraphMatchMany
type MatchExpr struct {
	Expr     ExprNode        // the graph-valued expression being matched
	Patterns []*GraphPattern // one or more pattern alternatives
	Loc      Loc
}

func (*MatchExpr) nodeTag()      {}
func (n *MatchExpr) GetLoc() Loc { return n.Loc }
func (*MatchExpr) exprNode()     {}

// GraphPattern is one complete pattern: optional selector + restrictor +
// path variable + a sequence of node/edge pattern parts.
//
// Grammar: gpmlPattern
type GraphPattern struct {
	Selector   *PatternSelector
	Restrictor PatternRestrictor
	Variable   *VarRef // optional `p = ...` path variable binding
	Parts      []PatternNode
	Loc        Loc
}

func (*GraphPattern) nodeTag()      {}
func (n *GraphPattern) GetLoc() Loc { return n.Loc }
func (*GraphPattern) patternNode()  {}

// NodePattern represents `(var:Label WHERE …)`. Variable, Labels, and
// Where are all optional.
//
// Grammar: node
type NodePattern struct {
	Variable *VarRef
	Labels   []string
	Where    ExprNode
	Loc      Loc
}

func (*NodePattern) nodeTag()      {}
func (n *NodePattern) GetLoc() Loc { return n.Loc }
func (*NodePattern) patternNode()  {}

// EdgePattern represents `-[var:Label]->`, `<-[]-`, `~[]~`, etc.
// Direction is required; Variable, Labels, Where, and Quantifier are optional.
//
// Grammar: edge#EdgeWithSpec, edge#EdgeAbbreviated
//
//	(direction comes from edgeWSpec labels / edgeAbbrev; body from edgeSpec)
type EdgePattern struct {
	Direction  EdgeDirection
	Variable   *VarRef
	Labels     []string
	Where      ExprNode
	Quantifier *PatternQuantifier
	Loc        Loc
}

func (*EdgePattern) nodeTag()      {}
func (n *EdgePattern) GetLoc() Loc { return n.Loc }
func (*EdgePattern) patternNode()  {}

// PatternQuantifier represents the `+`, `*`, or `{m,n}` decorator on a
// pattern part. Min and Max use -1 to indicate "unbounded".
//
// Grammar: patternQuantifier
type PatternQuantifier struct {
	Min int
	Max int // -1 = unbounded
	Loc Loc
}

func (*PatternQuantifier) nodeTag()      {}
func (n *PatternQuantifier) GetLoc() Loc { return n.Loc }

// PatternSelector represents the optional selector keyword on a graph
// pattern: `ANY`, `ALL SHORTEST`, or `SHORTEST k`. K is non-zero only for
// SHORTEST K.
//
// Grammar: matchSelector#SelectorBasic, matchSelector#SelectorAny, matchSelector#SelectorShortest
type PatternSelector struct {
	Kind SelectorKind
	K    int
	Loc  Loc
}

func (*PatternSelector) nodeTag()      {}
func (n *PatternSelector) GetLoc() Loc { return n.Loc }
