package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mysql/ast"
)

// procScope is the static-validation scope MySQL maintains via sp_pcontext
// while parsing a stored-program body. omni mirrors the model: each
// BEGIN...END (and each DECLARE HANDLER body) opens a scope; declarations
// register into kind-specific namespaces; references walk the parent chain
// with kind-specific rules.
//
// Namespaces are kept separate (matching MySQL's sp_pcontext m_vars /
// m_conditions / m_cursors / m_labels). Lookup is case-insensitive: MySQL
// procedure-local names are case-insensitive (container-verified
// 2026-04-21 round; consistent with begin/end label matching from PR series
// commit 4c75a83).
//
// Handler body scopes are tagged scopeHandlerBody. Label lookup stops at a
// handler-scope barrier: a LEAVE/ITERATE inside a handler body cannot
// target a label declared in the enclosing routine (MySQL's HANDLER_SCOPE
// flag in sp_pcontext).
type procScope struct {
	parent     *procScope
	kind       scopeKind
	vars       map[string]*nodes.DeclareVarStmt       // includes routine parameters seeded at top
	conditions map[string]*nodes.DeclareConditionStmt // user-declared conditions
	cursors    map[string]*nodes.DeclareCursorStmt
	labels     map[string]labelInfo
	// isFunction is set on the outermost scope of a CREATE FUNCTION body and
	// inherited by lookup helpers via the parent chain. Used to gate
	// RETURN-allowed and result-set-restriction checks.
	isFunction bool
}

type scopeKind int

const (
	scopeBlock scopeKind = iota
	scopeHandlerBody
)

// labelKind distinguishes BEGIN-style labels (LEAVE-only target) from
// loop-style labels (LEAVE + ITERATE target), per MySQL semantics.
type labelKind int

const (
	labelBegin labelKind = iota
	labelLoop
)

type labelInfo struct {
	kind labelKind
	node nodes.Node
}

// pushScope creates a child scope with the given kind and makes it current.
// Returns the new scope so callers can stash a defer popScope() pair if
// preferred (the standard pattern in this file).
func (p *Parser) pushScope(kind scopeKind) *procScope {
	s := &procScope{
		parent:     p.procScope,
		kind:       kind,
		vars:       make(map[string]*nodes.DeclareVarStmt),
		conditions: make(map[string]*nodes.DeclareConditionStmt),
		cursors:    make(map[string]*nodes.DeclareCursorStmt),
		labels:     make(map[string]labelInfo),
	}
	if p.procScope != nil {
		s.isFunction = p.procScope.isFunction
	}
	p.procScope = s
	return s
}

// popScope restores the parent scope.
func (p *Parser) popScope() {
	if p.procScope != nil {
		p.procScope = p.procScope.parent
	}
}

// lookupVar walks the parent chain for a variable declaration matching name
// (case-insensitive). Returns nil if not found.
func (p *Parser) lookupVar(name string) *nodes.DeclareVarStmt {
	key := strings.ToLower(name)
	for cur := p.procScope; cur != nil; cur = cur.parent {
		if v, ok := cur.vars[key]; ok {
			return v
		}
	}
	return nil
}

// lookupCondition walks the parent chain for a user-declared condition.
func (p *Parser) lookupCondition(name string) *nodes.DeclareConditionStmt {
	key := strings.ToLower(name)
	for cur := p.procScope; cur != nil; cur = cur.parent {
		if c, ok := cur.conditions[key]; ok {
			return c
		}
	}
	return nil
}

// lookupCursor walks the parent chain for a cursor declaration.
func (p *Parser) lookupCursor(name string) *nodes.DeclareCursorStmt {
	key := strings.ToLower(name)
	for cur := p.procScope; cur != nil; cur = cur.parent {
		if c, ok := cur.cursors[key]; ok {
			return c
		}
	}
	return nil
}

// lookupLabel walks the parent chain for a label, applying the
// handler-scope barrier (labels declared outside a handler are not visible
// from within its body). When mustBeLoop is true, only loop-kind labels
// match (used by ITERATE).
func (p *Parser) lookupLabel(name string, mustBeLoop bool) (labelInfo, bool) {
	key := strings.ToLower(name)
	for cur := p.procScope; cur != nil; cur = cur.parent {
		if l, ok := cur.labels[key]; ok {
			if mustBeLoop && l.kind != labelLoop {
				return labelInfo{}, false
			}
			return l, true
		}
		if cur.kind == scopeHandlerBody {
			// label barrier — handler bodies don't see outer labels
			return labelInfo{}, false
		}
	}
	return labelInfo{}, false
}

// declareVar inserts a variable into the current scope. Returns an error
// if a same-scope same-name variable, condition, or cursor already exists
// (MySQL declares vars/conditions/cursors share insertion-time conflict
// detection; ER_SP_DUP_VAR / ER_SP_DUP_COND / ER_SP_DUP_CURS).
func (p *Parser) declareVar(name string, stmt *nodes.DeclareVarStmt, pos int) error {
	if p.procScope == nil {
		return nil
	}
	key := strings.ToLower(name)
	if _, exists := p.procScope.vars[key]; exists {
		return &ParseError{Message: "duplicate variable declaration: " + name, Position: pos}
	}
	p.procScope.vars[key] = stmt
	return nil
}

// declareCondition inserts a user-declared condition into the current scope.
func (p *Parser) declareCondition(name string, stmt *nodes.DeclareConditionStmt, pos int) error {
	if p.procScope == nil {
		return nil
	}
	key := strings.ToLower(name)
	if _, exists := p.procScope.conditions[key]; exists {
		return &ParseError{Message: "duplicate condition declaration: " + name, Position: pos}
	}
	p.procScope.conditions[key] = stmt
	return nil
}

// declareCursor inserts a cursor into the current scope.
func (p *Parser) declareCursor(name string, stmt *nodes.DeclareCursorStmt, pos int) error {
	if p.procScope == nil {
		return nil
	}
	key := strings.ToLower(name)
	if _, exists := p.procScope.cursors[key]; exists {
		return &ParseError{Message: "duplicate cursor declaration: " + name, Position: pos}
	}
	p.procScope.cursors[key] = stmt
	return nil
}

// declareLabel registers a label in the current scope. The label belongs
// to the scope active when the labeled compound is encountered (so a
// LEAVE inside the body can reach it via the parent chain). MySQL forbids
// label reuse anywhere in the visibility chain, so duplicate detection
// walks parents up to (but not across) a handler-scope barrier.
func (p *Parser) declareLabel(name string, kind labelKind, node nodes.Node, pos int) error {
	if p.procScope == nil {
		return nil
	}
	key := strings.ToLower(name)
	for cur := p.procScope; cur != nil; cur = cur.parent {
		if _, exists := cur.labels[key]; exists {
			return &ParseError{Message: "duplicate label: " + name, Position: pos}
		}
		if cur.kind == scopeHandlerBody {
			break
		}
	}
	p.procScope.labels[key] = labelInfo{kind: kind, node: node}
	return nil
}
