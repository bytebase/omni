package validate

import nodes "github.com/bytebase/omni/mysql/ast"

// scopeKind classifies a lexical scope.
type scopeKind int

const (
	scopeBlock scopeKind = iota
	scopeHandlerBody
)

// valScope mirrors MySQL's sp_pcontext layout: separate namespaces per kind.
type valScope struct {
	parent     *valScope
	kind       scopeKind
	vars       map[string]*nodes.DeclareVarStmt
	conditions map[string]*nodes.DeclareConditionStmt
	cursors    map[string]*nodes.DeclareCursorStmt
	labels     map[string]labelInfo
	isFunction bool
}

type labelKind int

const (
	labelBlock labelKind = iota
	labelLoop
)

type labelInfo struct {
	kind labelKind
	pos  int
}

func newScope(parent *valScope, kind scopeKind) *valScope {
	s := &valScope{
		parent:     parent,
		kind:       kind,
		vars:       map[string]*nodes.DeclareVarStmt{},
		conditions: map[string]*nodes.DeclareConditionStmt{},
		cursors:    map[string]*nodes.DeclareCursorStmt{},
		labels:     map[string]labelInfo{},
	}
	if parent != nil {
		s.isFunction = parent.isFunction
	}
	return s
}

func (s *valScope) lookupVar(name string) *nodes.DeclareVarStmt {
	for cur := s; cur != nil; cur = cur.parent {
		if v, ok := cur.vars[lower(name)]; ok {
			return v
		}
	}
	return nil
}

func (s *valScope) lookupCondition(name string) *nodes.DeclareConditionStmt {
	for cur := s; cur != nil; cur = cur.parent {
		if v, ok := cur.conditions[lower(name)]; ok {
			return v
		}
	}
	return nil
}

func (s *valScope) lookupCursor(name string) *nodes.DeclareCursorStmt {
	for cur := s; cur != nil; cur = cur.parent {
		if v, ok := cur.cursors[lower(name)]; ok {
			return v
		}
	}
	return nil
}

// lookupLabel walks the parent chain; labels are blocked at handler-body scopes
// (MySQL's label barrier). loopOnly=true limits matches to loop labels (for ITERATE).
func (s *valScope) lookupLabel(name string, loopOnly bool) (labelInfo, bool) {
	for cur := s; cur != nil; cur = cur.parent {
		if info, ok := cur.labels[lower(name)]; ok {
			if loopOnly && info.kind != labelLoop {
				return labelInfo{}, false
			}
			return info, true
		}
		if cur.kind == scopeHandlerBody {
			return labelInfo{}, false
		}
	}
	return labelInfo{}, false
}

func lower(s string) string {
	// ASCII lowering is sufficient for MySQL identifier case-folding in the
	// contexts the validator cares about. Keep local to avoid importing strings
	// in hot paths if later moved inline.
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		b[i] = c
	}
	return string(b)
}
