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
	case *nodes.DeclareVarStmt:
		v.registerDeclareVar(s)
	case *nodes.DeclareConditionStmt:
		v.registerDeclareCondition(s)
	case *nodes.DeclareCursorStmt:
		v.registerDeclareCursor(s)
	case *nodes.DeclareHandlerStmt:
		v.walkDeclareHandler(s)
	case *nodes.WhileStmt:
		v.walkWhileStmt(s)
	case *nodes.RepeatStmt:
		v.walkRepeatStmt(s)
	case *nodes.LoopStmt:
		v.walkLoopStmt(s)
	case *nodes.IfStmt:
		v.walkIfStmt(s)
	case *nodes.CaseStmtNode:
		v.walkCaseStmt(s)
	case *nodes.LeaveStmt:
		v.walkLeaveStmt(s)
	case *nodes.IterateStmt:
		v.walkIterateStmt(s)
	}
}

// handlerCondKey dedups a DECLARE HANDLER FOR condition-value list, mirroring
// the parser's local key type.
type handlerCondKey struct {
	kind  nodes.HandlerCondKind
	value string
}

func (v *validator) walkDeclareHandler(s *nodes.DeclareHandlerStmt) {
	seen := make(map[handlerCondKey]bool, len(s.Conditions))
	for _, c := range s.Conditions {
		k := handlerCondKey{kind: c.Kind, value: lower(c.Value)}
		if seen[k] {
			v.emit("duplicate_handler_condition",
				"duplicate condition value in handler declaration", s.Loc.Start)
			continue
		}
		seen[k] = true
		if c.Kind == nodes.HandlerCondName {
			if v.scope == nil || v.scope.lookupCondition(c.Value) == nil {
				v.emit("undeclared_condition", "undeclared condition: "+c.Value, s.Loc.Start)
			}
		}
	}
	// Handler body runs inside a HANDLER_SCOPE: outer vars/conditions/cursors
	// visible via parent chain, but outer labels are hidden (label barrier).
	v.push(scopeHandlerBody)
	v.walk(s.Stmt)
	v.pop()
}

func (v *validator) registerDeclareVar(s *nodes.DeclareVarStmt) {
	if v.scope == nil {
		return
	}
	for _, name := range s.Names {
		key := lower(name)
		if _, exists := v.scope.vars[key]; exists {
			v.emit("duplicate_variable", "duplicate variable declaration: "+name, s.Loc.Start)
			continue
		}
		v.scope.vars[key] = s
	}
}

func (v *validator) registerDeclareCondition(s *nodes.DeclareConditionStmt) {
	if v.scope == nil {
		return
	}
	key := lower(s.Name)
	if _, exists := v.scope.conditions[key]; exists {
		v.emit("duplicate_condition", "duplicate condition declaration: "+s.Name, s.Loc.Start)
		return
	}
	v.scope.conditions[key] = s
}

func (v *validator) registerDeclareCursor(s *nodes.DeclareCursorStmt) {
	if v.scope == nil {
		return
	}
	key := lower(s.Name)
	if _, exists := v.scope.cursors[key]; exists {
		v.emit("duplicate_cursor", "duplicate cursor declaration: "+s.Name, s.Loc.Start)
		return
	}
	v.scope.cursors[key] = s
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
	// The block's label belongs to the ENCLOSING scope, mirroring sp_pcontext:
	// BEGIN ... END's label is visible to siblings/parents, not to the body's
	// own declarations.
	v.registerLabel(b.Label, labelBlock, b.Loc.Start)
	v.push(scopeBlock)
	for _, s := range b.Stmts {
		v.walk(s)
	}
	v.pop()
}

// registerLabel installs a label into the current (enclosing) scope and emits
// duplicate_label on collision. Safe to call with an empty name.
func (v *validator) registerLabel(name string, kind labelKind, pos int) {
	if name == "" || v.scope == nil {
		return
	}
	key := lower(name)
	if _, exists := v.scope.labels[key]; exists {
		v.emit("duplicate_label", "duplicate label: "+name, pos)
		return
	}
	v.scope.labels[key] = labelInfo{kind: kind, pos: pos}
}

func (v *validator) walkWhileStmt(s *nodes.WhileStmt) {
	v.registerLabel(s.Label, labelLoop, s.Loc.Start)
	v.push(scopeBlock)
	for _, st := range s.Stmts {
		v.walk(st)
	}
	v.pop()
}

func (v *validator) walkRepeatStmt(s *nodes.RepeatStmt) {
	v.registerLabel(s.Label, labelLoop, s.Loc.Start)
	v.push(scopeBlock)
	for _, st := range s.Stmts {
		v.walk(st)
	}
	v.pop()
}

func (v *validator) walkLoopStmt(s *nodes.LoopStmt) {
	v.registerLabel(s.Label, labelLoop, s.Loc.Start)
	v.push(scopeBlock)
	for _, st := range s.Stmts {
		v.walk(st)
	}
	v.pop()
}

func (v *validator) walkIfStmt(s *nodes.IfStmt) {
	for _, st := range s.ThenList {
		v.walk(st)
	}
	for _, ei := range s.ElseIfs {
		for _, st := range ei.ThenList {
			v.walk(st)
		}
	}
	for _, st := range s.ElseList {
		v.walk(st)
	}
}

func (v *validator) walkCaseStmt(s *nodes.CaseStmtNode) {
	for _, w := range s.Whens {
		for _, st := range w.ThenList {
			v.walk(st)
		}
	}
	for _, st := range s.ElseList {
		v.walk(st)
	}
}

func (v *validator) walkLeaveStmt(s *nodes.LeaveStmt) {
	if v.scope == nil {
		return
	}
	if _, ok := v.scope.lookupLabel(s.Label, false); !ok {
		v.emit("undeclared_label", "LEAVE references undeclared label: "+s.Label, s.Loc.Start)
	}
}

func (v *validator) walkIterateStmt(s *nodes.IterateStmt) {
	if v.scope == nil {
		return
	}
	if _, ok := v.scope.lookupLabel(s.Label, true); !ok {
		v.emit("undeclared_loop_label",
			"ITERATE references undeclared loop label: "+s.Label, s.Loc.Start)
	}
}
