package parser

// What MariaDB 11.8 refuses of a DECLARE's variables while it parses, beyond what sp_scope.go
// holds: TYPE OF a variable no scope in reach declares (1327). A DECLARE's own names are
// declared once its type is read, so its TYPE OF does not see them, though a name it repeats,
// or one the block already has, is refused as it is read (1331). Measured on mariadb:11.8.8.
//
// Ref: sp_decl_idents, sp_decl_variable_list and sp_decl_variable_list_anchored in MariaDB
// 11.8 sql/sql_yacc.yy.

// checkSPVarFree refuses a variable a DECLARE names that the block, or the DECLARE before it,
// already names: 1331, ER_SP_DUP_VAR.
func (p *Parser) checkSPVarFree(before []string, name string, pos int) error {
	if !p.spChecks() {
		return nil
	}
	if spHasName(p.routine.sp.vars, name) || spHasName(before, name) {
		return &ParseError{Message: "Duplicate variable: " + name, Position: pos}
	}
	return nil
}

// declareSPVars adds a DECLARE's variables to the current block.
func (p *Parser) declareSPVars(names []string) {
	if s := p.routine.sp; s != nil {
		s.vars = append(s.vars, names...)
	}
}

// checkSPVarDeclared refuses TYPE OF a variable no scope in reach declares: 1327,
// ER_SP_UNDECLARED_VAR.
func (p *Parser) checkSPVarDeclared(name string, pos int) error {
	if p.spChecks() && !p.routine.sp.findVar(name) {
		return &ParseError{Message: "Undeclared variable: " + name, Position: pos}
	}
	return nil
}
