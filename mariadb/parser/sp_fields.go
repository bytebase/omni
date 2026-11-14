package parser

// What MariaDB 11.8 refuses of a name in an expression that no query is around, while it
// parses. Such an expression reads no table, so a bare name in it can only be a stored
// program's variable or parameter: inside a stored program LEX::create_item_ident_sp answers
// 1327, ER_SP_UNDECLARED_VAR, for any other name when fields_are_impossible() (no query is
// open, or only the service query of a DO or a SET of a variable); outside one,
// LEX::check_expr_allows_fields_or_error answers 1054, ER_BAD_FIELD_ERROR, for any name of a
// column when no query is open at all.
//
// The expressions with no query around them: CALL's arguments, EXECUTE ... USING's
// parameters, the values of SIGNAL's and RESIGNAL's SET, KILL's, and GET DIAGNOSTICS
// CONDITION's number; and in a stored program also RETURN's, IF's, ELSEIF's, WHILE's,
// UNTIL's, CASE's, DECLARE ... DEFAULT's, DO's and a SET of a variable. mariadb:11.8.8
// takes a qualified name in the last group (`RETURN a.b`), which a ROW variable may answer
// to, and SET of a user or system variable outside the group.
//
// Ref: LEX::create_item_ident_sp, LEX::check_expr_allows_fields_or_error and
// LEX::fields_are_impossible in MariaDB 11.8 sql/sql_lex.cc and sql/sql_lex.h.

// noFieldsScope is the expression with no query around it that the parser is inside.
type noFieldsScope struct {
	where   string // the statement, as the 1054 names it; "" outside such an expression
	outside bool   // refused outside a stored program too
}

// enterNoFields marks the expressions parsed until the returned function runs as having no
// query around them.
func (p *Parser) enterNoFields(where string, outside bool) func() {
	outer := p.noFields
	p.noFields = noFieldsScope{where: where, outside: outside}
	return func() { p.noFields = outer }
}

// enterQuery marks a query's own clauses, where a name may be a column.
func (p *Parser) enterQuery() func() {
	if p.noFields.where == "" {
		return func() {}
	}
	return p.enterNoFields("", false)
}

// refuseNoFieldsName is the answer to a name read as a column in such an expression.
func (p *Parser) refuseNoFieldsName(column string, qualified bool, pos int) error {
	if p.completing || p.noFields.where == "" {
		return nil
	}
	if s := p.routine.sp; s != nil {
		if qualified || s.findVar(column) {
			return nil
		}
		return &ParseError{Message: "Undeclared variable: " + column, Position: pos}
	}
	if !p.noFields.outside {
		return nil
	}
	return &ParseError{Message: "Unknown column '" + column + "' in '" + p.noFields.where + "'", Position: pos}
}

// hasVar reports a variable or parameter of exactly the name in reach, the test for a SET of
// a variable, whose value then has no query around it.
func (s *spScope) hasVar(name string) bool {
	for ; s != nil; s = s.parent {
		if spHasName(s.vars, name) {
			return true
		}
	}
	return false
}
