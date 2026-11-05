package parser

import "strings"

// keywordsNotFunctionNames are the keywords MariaDB 11.8.8 refuses as an unquoted function name,
// in every argument shape: `SELECT comment(3)` is 1064 where `SELECT foo(3)` and
// the quoted `comment`(3) parse. They are no other word's business: all of them still
// name a column or a table where their category lets them, and a qualified `d.comment(3)`
// parses.
//
// The list is measured and not derived, because this parser's keyword table is not the
// engine's: it holds words the engine lexes as plain identifiers, which do take the call.
// TestKeywordFunctionNamesOracle holds it to mariadb:11.8.8 over every keyword the lexer has.
//
// MariaDB takes far more keywords here than MySQL does, so its list is short. Nothing in
// the upstream grammar was read for this; every word is what 11.8.8 answered.
var keywordsNotFunctionNames = wordSet(`
any backup begin binlog bit bool boolean byte cache checksum close code comment commit
compressed datetime deallocate do end enum exclude execute fixed flush following follows
global handler help host install json language local medium names national nchar no
number nvarchar offset open options others owner parser port precedes preceding prepare
remove repair replica replicas reset restore role rollback savepoint security serial
server session shutdown signed slave socket some soname sounds sql_after_gtids
sql_before_gtids start stats_auto_recalc stats_persistent stats_sample_pages stop text
ties unbounded unicode uninstall upgrade vector wrapper xa
`)

func wordSet(words string) map[string]bool {
	set := make(map[string]bool)
	for _, word := range strings.Fields(words) {
		set[word] = true
	}
	return set
}

// refusesFunctionName reports whether the current token is a keyword the engine does not
// take as the name of the call that follows it. A quoted name lexes as an identifier
// whatever it spells, so it never comes here.
func (p *Parser) refusesFunctionName() bool {
	return p.cur.Type != tokIDENT && p.peekNext().Type == '(' && keywordsNotFunctionNames[strings.ToLower(p.cur.Str)]
}
