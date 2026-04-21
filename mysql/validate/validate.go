// Package validate performs static semantic checks on a MySQL AST produced by
// mysql/parser. It is the omni analogue of MySQL's sp_head::parse /
// sp_pcontext phase: grammar errors come from the parser, semantic errors
// (undeclared var/cursor/label, duplicate DECLARE, missing RETURN, etc.) come
// from here.
package validate

import nodes "github.com/bytebase/omni/mysql/ast"

// Options tunes which validators run. Reserved for future strictness toggles.
type Options struct{}

// Validate walks a parsed AST and returns all semantic diagnostics. An empty
// slice means "no issues"; a nil AST returns nil.
func Validate(list *nodes.List, _ Options) []Diagnostic {
	if list == nil {
		return nil
	}
	v := &validator{}
	for _, stmt := range list.Items {
		v.walk(stmt)
	}
	return v.diagnostics
}

type validator struct {
	scope       *valScope
	diagnostics []Diagnostic
}

func (v *validator) push(kind scopeKind) *valScope {
	s := newScope(v.scope, kind)
	v.scope = s
	return s
}

func (v *validator) pop() {
	if v.scope != nil {
		v.scope = v.scope.parent
	}
}

func (v *validator) emit(code, msg string, pos int) {
	v.diagnostics = append(v.diagnostics, Diagnostic{
		Code:     code,
		Message:  msg,
		Severity: SeverityError,
		Position: pos,
	})
}

// walk dispatches on the node type. In this scaffolding commit it only
// descends into routine bodies and compound blocks; each semantic check is
// wired up by a later task.
func (v *validator) walk(n nodes.Node) {
	switch s := n.(type) {
	case nil:
		return
	case *nodes.CreateFunctionStmt:
		v.walkRoutine(s.Body, !s.IsProcedure, s.Params)
	case *nodes.CreateTriggerStmt:
		v.walkRoutine(s.Body, false, nil)
	case *nodes.CreateEventStmt:
		// Event body, if present, is a compound statement. The field is not
		// directly named here; follow the AST definition.
		v.walkEventBody(s)
	case *nodes.BeginEndBlock:
		v.walkBeginEnd(s)
	}
}

func (v *validator) walkRoutine(body nodes.Node, isFunction bool, params []*nodes.FuncParam) {
	if body == nil {
		return
	}
	scope := v.push(scopeBlock)
	scope.isFunction = isFunction
	for _, p := range params {
		scope.vars[lower(p.Name)] = &nodes.DeclareVarStmt{
			Loc:      p.Loc,
			Names:    []string{p.Name},
			TypeName: p.TypeName,
		}
	}
	v.walk(body)
	v.pop()
}

func (v *validator) walkEventBody(_ *nodes.CreateEventStmt) {
	// Fill in once we wire up event-body walking; skeletal for now.
}

func (v *validator) walkBeginEnd(b *nodes.BeginEndBlock) {
	v.push(scopeBlock)
	for _, s := range b.Stmts {
		v.walk(s)
	}
	v.pop()
}
