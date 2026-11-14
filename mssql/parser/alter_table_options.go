// Package parser - alter_table_options.go reads the option lists of ALTER
// TABLE (SET, REBUILD, SWITCH, DROP ... WITH, ALTER COLUMN ... WITH) by
// their grammar: each option name comes from a closed set and each value
// from what that option takes, where SQL Server refuses anything else.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-table-transact-sql
package parser

import (
	"strconv"
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// tableWordIs reports whether the current token is word, an upper-case
// option word: its keyword token when the lexer registers one, else an
// undelimited identifier spelling it. SQL Server reads option words by
// their text, so [TABLE] is not the word TABLE.
func (p *Parser) tableWordIs(word string) bool {
	if kw := lookupKeyword(word); kw != tokIDENT {
		return p.cur.Type == kw
	}
	return p.cur.Type == tokIDENT && !p.isDelimited(p.cur) && strings.ToUpper(p.cur.Str) == word
}

// tableWordIn returns the current token's word, upper-cased, when it is
// one of words, and "" otherwise.
func (p *Parser) tableWordIn(words ...string) string {
	for _, w := range words {
		if p.tableWordIs(w) {
			return w
		}
	}
	return ""
}

// tableExpectWord consumes the current token when it is word.
func (p *Parser) tableExpectWord(word string) error {
	if !p.tableWordIs(word) {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	return nil
}

// tableExpectWordIn consumes the current token when it is one of words
// and returns it upper-cased.
func (p *Parser) tableExpectWordIn(words ...string) (string, error) {
	w := p.tableWordIn(words...)
	if w == "" {
		return "", p.syntaxErrorAtCur()
	}
	p.advance()
	return w, nil
}

// tableExpectInteger consumes an unsigned integer literal and returns
// its text.
func (p *Parser) tableExpectInteger() (string, error) {
	if p.cur.Type != tokICONST {
		return "", p.syntaxErrorAtCur()
	}
	s := p.cur.Str
	if s == "" {
		s = strconv.FormatInt(p.cur.Ival, 10)
	}
	p.advance()
	return s, nil
}

// tableExpectOnOff consumes ON or OFF and returns it.
func (p *Parser) tableExpectOnOff() (string, error) {
	switch p.cur.Type {
	case kwON:
		p.advance()
		return "ON", nil
	case kwOFF:
		p.advance()
		return "OFF", nil
	}
	return "", p.syntaxErrorAtCur()
}

// tableOptionName consumes an option name of names followed by =, and
// returns it upper-cased. A name seen before is refused when seen is not
// nil: SQL Server takes each of these options once.
func (p *Parser) tableOptionName(seen map[string]bool, names ...string) (string, error) {
	name, err := p.tableExpectWordIn(names...)
	if err != nil {
		return "", err
	}
	if seen != nil {
		if seen[name] {
			return "", &ParseError{Message: "option " + name + " is specified more than once", Position: p.prev.Loc}
		}
		seen[name] = true
	}
	if _, err := p.expect('='); err != nil {
		return "", err
	}
	return name, nil
}

// parseAlterTableSetOptions parses the ( table_option [ ,...n ] ) of ALTER
// TABLE ... SET, the current token being the open parenthesis. It returns the
// names and values as alternating String nodes, a list of suboptions after
// the value that takes them.
//
//	alterTableSetStatement (TSql160.g:28051), tableOption (TSql160.g:28070):
//	  LOCK_ESCALATION = { AUTO | TABLE | DISABLE }
//	| FILESTREAM_ON = { filegroup | partition_scheme | 'string' | "default" }
//	| FILETABLE_DIRECTORY = 'directory_name'
//	| SYSTEM_VERSIONING = { OFF | ON [ ( system_versioning_option [ ,...n ] ) ] }
//	| REMOTE_DATA_ARCHIVE = { ON | OFF | OFF_WITHOUT_DATA_RECOVERY }
//	      ( [ FILTER_PREDICATE = { NULL | function_call } , ] MIGRATION_STATE = { OUTBOUND | INBOUND | PAUSED } )
//
// DATA_DELETION is Azure SQL Edge's: SQL Server 2022 refuses it (Msg 102).
func (p *Parser) parseAlterTableSetOptions() (*nodes.List, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	var items []nodes.Node
	for {
		name, err := p.tableOptionName(nil, "LOCK_ESCALATION", "FILESTREAM_ON", "FILETABLE_DIRECTORY", "SYSTEM_VERSIONING", "REMOTE_DATA_ARCHIVE")
		if err != nil {
			return nil, err
		}
		items = append(items, &nodes.String{Str: name})
		switch name {
		case "LOCK_ESCALATION":
			// lockEscalationTableOption, TSql160.g:28091.
			var v string
			if p.cur.Type == kwTABLE {
				v = "TABLE"
				p.advance()
			} else if v, err = p.tableExpectWordIn("AUTO", "DISABLE"); err != nil {
				return nil, err
			}
			items = append(items, &nodes.String{Str: v})
		case "FILESTREAM_ON":
			// fileStreamOnTableOption, TSql160.g:28113: a string or an
			// identifier; "default" is a quoted identifier.
			if p.cur.Type != tokSCONST && p.cur.Type != tokNSCONST && !p.isIdentLike() {
				return nil, p.syntaxErrorAtCur()
			}
			items = append(items, &nodes.String{Str: p.cur.Str})
			p.advance()
		case "FILETABLE_DIRECTORY":
			// SQL Server refuses NULL here though the grammar's
			// stringLiteralOrNull takes it.
			if p.cur.Type != tokSCONST && p.cur.Type != tokNSCONST {
				return nil, p.syntaxErrorAtCur()
			}
			items = append(items, &nodes.String{Str: p.cur.Str})
			p.advance()
		case "SYSTEM_VERSIONING":
			v, err := p.tableExpectOnOff()
			if err != nil {
				return nil, err
			}
			items = append(items, &nodes.String{Str: v})
			if v == "ON" && p.cur.Type == '(' {
				sub, err := p.parseAlterTableSystemVersioningSuboptions()
				if err != nil {
					return nil, err
				}
				items = append(items, sub)
			}
		case "REMOTE_DATA_ARCHIVE":
			var v string
			switch {
			case p.cur.Type == kwON || p.cur.Type == kwOFF:
				v, _ = p.tableExpectOnOff()
			default:
				if v, err = p.tableExpectWordIn("OFF_WITHOUT_DATA_RECOVERY"); err != nil {
					return nil, err
				}
			}
			items = append(items, &nodes.String{Str: v})
			sub, err := p.parseAlterTableRemoteDataArchiveSettings()
			if err != nil {
				return nil, err
			}
			items = append(items, sub)
		}
		if _, ok := p.match(','); !ok {
			break
		}
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.List{Items: items}, nil
}

// parseAlterTableSystemVersioningSuboptions parses the parenthesized options
// of SYSTEM_VERSIONING = ON, each at most once (systemVersioningTableSuboption,
// TSql160.g:27218):
//
//	HISTORY_TABLE = schema_name.history_table_name
//	| DATA_CONSISTENCY_CHECK = { ON | OFF }
//	| HISTORY_RETENTION_PERIOD = { INFINITE | number { DAY | DAYS | WEEK | WEEKS | MONTH | MONTHS | YEAR | YEARS } }
//
// SQL Server's parser wants the history table in two parts and a positive
// retention period.
func (p *Parser) parseAlterTableSystemVersioningSuboptions() (*nodes.List, error) {
	p.advance() // (
	seen := map[string]bool{}
	var items []nodes.Node
	for {
		name, err := p.tableOptionName(seen, "HISTORY_TABLE", "DATA_CONSISTENCY_CHECK", "HISTORY_RETENTION_PERIOD")
		if err != nil {
			return nil, err
		}
		items = append(items, &nodes.String{Str: name})
		switch name {
		case "HISTORY_TABLE":
			start := p.cur
			ref, err := p.parseTableRef()
			if err != nil {
				return nil, err
			}
			if ref == nil || ref.Schema == "" || ref.Database != "" || ref.Server != "" {
				return nil, &ParseError{Message: "the history table must be named in two parts", Position: start.Loc}
			}
			items = append(items, &nodes.String{Str: ref.Schema + "." + ref.Object})
		case "DATA_CONSISTENCY_CHECK":
			v, err := p.tableExpectOnOff()
			if err != nil {
				return nil, err
			}
			items = append(items, &nodes.String{Str: v})
		case "HISTORY_RETENTION_PERIOD":
			v, err := p.parseAlterTableRetentionPeriod()
			if err != nil {
				return nil, err
			}
			items = append(items, &nodes.String{Str: v})
		}
		if _, ok := p.match(','); !ok {
			break
		}
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.List{Items: items}, nil
}

// parseAlterTableRetentionPeriod parses INFINITE or a positive number and a
// unit, and returns them as written.
func (p *Parser) parseAlterTableRetentionPeriod() (string, error) {
	if p.tableWordIs("INFINITE") {
		p.advance()
		return "INFINITE", nil
	}
	if p.cur.Type != tokICONST || p.cur.Ival <= 0 {
		return "", p.syntaxErrorAtCur()
	}
	n, _ := p.tableExpectInteger()
	unit, err := p.tableExpectWordIn("DAY", "DAYS", "WEEK", "WEEKS", "MONTH", "MONTHS", "YEAR", "YEARS")
	if err != nil {
		return "", err
	}
	return n + " " + unit, nil
}

// parseAlterTableRemoteDataArchiveSettings parses the required
// ( [ FILTER_PREDICATE = { NULL | function_call } , ] MIGRATION_STATE = state )
// after REMOTE_DATA_ARCHIVE = value, each setting once and MIGRATION_STATE
// required (remoteDataArchiveTableSettings, TSql160.g:27100; Msg 201 when it
// is missing).
func (p *Parser) parseAlterTableRemoteDataArchiveSettings() (*nodes.List, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	var items []nodes.Node
	for {
		name, err := p.tableOptionName(seen, "FILTER_PREDICATE", "MIGRATION_STATE")
		if err != nil {
			return nil, err
		}
		items = append(items, &nodes.String{Str: name})
		if name == "MIGRATION_STATE" {
			v, err := p.tableExpectWordIn("OUTBOUND", "INBOUND", "PAUSED")
			if err != nil {
				return nil, err
			}
			items = append(items, &nodes.String{Str: v})
		} else if p.cur.Type == kwNULL {
			p.advance()
			items = append(items, &nodes.String{Str: "NULL"})
		} else {
			start := p.pos()
			e, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			if _, ok := e.(*nodes.FuncCallExpr); !ok {
				return nil, &ParseError{Message: "FILTER_PREDICATE takes NULL or a function call", Position: start}
			}
			items = append(items, &nodes.String{Str: p.source[start:p.prevEnd()]})
		}
		if _, ok := p.match(','); !ok {
			break
		}
	}
	if !seen["MIGRATION_STATE"] {
		return nil, p.syntaxErrorAtCur()
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.List{Items: items}, nil
}

// parseAlterTableLowPriorityLockWait parses
//
//	WAIT_AT_LOW_PRIORITY ( MAX_DURATION = integer [ MINUTES ] , ABORT_AFTER_WAIT = { NONE | SELF | BLOCKERS } )
//
// (lowPriorityLockWaitOption, TSql160.g:17465) with WAIT_AT_LOW_PRIORITY the
// current token, and returns the word and its settings as String nodes and a
// List.
func (p *Parser) parseAlterTableLowPriorityLockWait() ([]nodes.Node, error) {
	if err := p.tableExpectWord("WAIT_AT_LOW_PRIORITY"); err != nil {
		return nil, err
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	if _, err := p.tableOptionName(nil, "MAX_DURATION"); err != nil {
		return nil, err
	}
	d, err := p.tableExpectInteger()
	if err != nil {
		return nil, err
	}
	if p.tableWordIs("MINUTES") {
		p.advance()
		d += " MINUTES"
	}
	if _, err := p.expect(','); err != nil {
		return nil, err
	}
	if _, err := p.tableOptionName(nil, "ABORT_AFTER_WAIT"); err != nil {
		return nil, err
	}
	abort, err := p.tableExpectWordIn("NONE", "SELF", "BLOCKERS")
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return []nodes.Node{
		&nodes.String{Str: "WAIT_AT_LOW_PRIORITY"},
		&nodes.List{Items: []nodes.Node{
			&nodes.String{Str: "MAX_DURATION"}, &nodes.String{Str: d},
			&nodes.String{Str: "ABORT_AFTER_WAIT"}, &nodes.String{Str: abort},
		}},
	}, nil
}

// parseAlterTableSwitchOptions parses the ( <low_priority_lock_wait> ) after
// SWITCH ... WITH, the current token being the open parenthesis. TRUNCATE_TARGET
// is Azure SQL's: SQL Server 2022 refuses it (Msg 102), and it takes one
// option (tableSwitchOptionList, TSql160.g:28559).
func (p *Parser) parseAlterTableSwitchOptions() (*nodes.List, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	items, err := p.parseAlterTableLowPriorityLockWait()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.List{Items: items}, nil
}

// parseAlterTableDropOptions parses the ( drop_clustered_constraint_option
// [ ,...n ] ) after a dropped constraint, the current token being the open
// parenthesis, each option once (Msg 153 for a second one):
//
//	dropClusteredConstraintOption (TSql160.g:28476):
//	  MAXDOP = max_degree_of_parallelism
//	| ONLINE = { ON | OFF }
//	| MOVE TO { partition_scheme_name ( column_name ) | filegroup | "default" }
//	| WAIT_AT_LOW_PRIORITY ( ... )
func (p *Parser) parseAlterTableDropOptions() (*nodes.List, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	var items []nodes.Node
	for {
		word := p.tableWordIn("MAXDOP", "ONLINE", "MOVE", "WAIT_AT_LOW_PRIORITY")
		if word == "" {
			return nil, p.syntaxErrorAtCur()
		}
		if seen[word] {
			return nil, &ParseError{Message: "option " + word + " is specified more than once", Position: p.cur.Loc}
		}
		seen[word] = true
		switch word {
		case "MAXDOP":
			p.advance()
			if _, err := p.expect('='); err != nil {
				return nil, err
			}
			n, err := p.tableExpectInteger()
			if err != nil {
				return nil, err
			}
			items = append(items, &nodes.String{Str: word}, &nodes.String{Str: n})
		case "ONLINE":
			p.advance()
			if _, err := p.expect('='); err != nil {
				return nil, err
			}
			v, err := p.tableExpectOnOff()
			if err != nil {
				return nil, err
			}
			items = append(items, &nodes.String{Str: word}, &nodes.String{Str: v})
		case "MOVE":
			p.advance()
			if _, err := p.expect(kwTO); err != nil {
				return nil, err
			}
			target, err := p.parseAlterTableFilegroupOrPartitionScheme()
			if err != nil {
				return nil, err
			}
			items = append(items, &nodes.String{Str: word}, &nodes.String{Str: target})
		case "WAIT_AT_LOW_PRIORITY":
			wait, err := p.parseAlterTableLowPriorityLockWait()
			if err != nil {
				return nil, err
			}
			items = append(items, wait...)
		}
		if _, ok := p.match(','); !ok {
			break
		}
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.List{Items: items}, nil
}

// parseAlterTableFilegroupOrPartitionScheme parses a filegroup, a string, or
// a partition scheme and its column (filegroupOrPartitionScheme,
// TSql160.g:17224), and returns it as written.
func (p *Parser) parseAlterTableFilegroupOrPartitionScheme() (string, error) {
	start := p.pos()
	if p.cur.Type != tokSCONST && p.cur.Type != tokNSCONST && !p.isIdentLike() {
		return "", p.syntaxErrorAtCur()
	}
	p.advance()
	if p.cur.Type == '(' {
		p.advance()
		if _, ok := p.parseIdentifier(); !ok {
			return "", p.syntaxErrorAtCur()
		}
		if _, err := p.expect(')'); err != nil {
			return "", err
		}
	}
	return p.source[start:p.prevEnd()], nil
}

// parseAlterTableRebuildOptions parses the ( rebuild_option [ ,...n ] ) of
// ALTER TABLE ... REBUILD WITH, the current token being the open parenthesis.
// Each option is a name = value BinaryExpr, as parseOptionList gives.
//
//	<rebuild_option> (ALTER TABLE; SQL Server refuses the other index options,
//	Msg 155):
//	  PAD_INDEX | SORT_IN_TEMPDB | STATISTICS_NORECOMPUTE | STATISTICS_INCREMENTAL
//	  | ALLOW_ROW_LOCKS | ALLOW_PAGE_LOCKS = { ON | OFF }
//	| FILLFACTOR | MAXDOP = integer
//	| IGNORE_DUP_KEY = { ON [ ( SUPPRESS_MESSAGES = { ON | OFF } ) ] | OFF }
//	| DATA_COMPRESSION = { NONE | ROW | PAGE | COLUMNSTORE | COLUMNSTORE_ARCHIVE } [ ON PARTITIONS ( range [ ,...n ] ) ]
//	| XML_COMPRESSION = { ON | OFF } [ ON PARTITIONS ( range [ ,...n ] ) ]
//	| ONLINE = { ON [ ( <low_priority_lock_wait> ) ] | OFF }
func (p *Parser) parseAlterTableRebuildOptions() (*nodes.List, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	var items []nodes.Node
	for {
		loc := p.pos()
		nameEnd := p.cur.End
		name, err := p.tableOptionName(nil,
			"PAD_INDEX", "SORT_IN_TEMPDB", "STATISTICS_NORECOMPUTE", "STATISTICS_INCREMENTAL",
			"ALLOW_ROW_LOCKS", "ALLOW_PAGE_LOCKS", "FILLFACTOR", "MAXDOP", "IGNORE_DUP_KEY",
			"DATA_COMPRESSION", "XML_COMPRESSION", "ONLINE")
		if err != nil {
			return nil, err
		}
		valLoc := p.pos()
		switch name {
		case "FILLFACTOR", "MAXDOP":
			_, err = p.tableExpectInteger()
		case "DATA_COMPRESSION":
			_, err = p.tableExpectWordIn("NONE", "ROW", "PAGE", "COLUMNSTORE", "COLUMNSTORE_ARCHIVE")
			if err == nil {
				err = p.parseAlterTableOnPartitions()
			}
		case "XML_COMPRESSION":
			_, err = p.tableExpectOnOff()
			if err == nil {
				err = p.parseAlterTableOnPartitions()
			}
		case "IGNORE_DUP_KEY", "ONLINE":
			var value string
			value, err = p.tableExpectOnOff()
			if err != nil || value != "ON" || p.cur.Type != '(' {
				break
			}
			p.advance()
			if name == "ONLINE" {
				if _, err = p.parseAlterTableLowPriorityLockWait(); err == nil {
					_, err = p.expect(')')
				}
				break
			}
			if _, err = p.tableOptionName(nil, "SUPPRESS_MESSAGES"); err == nil {
				if _, err = p.tableExpectOnOff(); err == nil {
					_, err = p.expect(')')
				}
			}
		default:
			_, err = p.tableExpectOnOff()
		}
		if err != nil {
			return nil, err
		}
		items = append(items, &nodes.BinaryExpr{
			Op:    nodes.BinOpEq,
			Left:  &nodes.ColumnRef{Column: name, Loc: nodes.Loc{Start: loc, End: nameEnd}},
			Right: &nodes.ColumnRef{Column: p.source[valLoc:p.prevEnd()], Loc: nodes.Loc{Start: valLoc, End: p.prevEnd()}},
			Loc:   nodes.Loc{Start: loc, End: p.prevEnd()},
		})
		if _, ok := p.match(','); !ok {
			break
		}
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.List{Items: items}, nil
}

// parseAlterTableOnPartitions parses an optional
// ON PARTITIONS ( range [ ,...n ] ), range being expression [ TO expression ]
// (compressionPartitionRange, TSql160.g:27623).
func (p *Parser) parseAlterTableOnPartitions() error {
	if p.cur.Type != kwON {
		return nil
	}
	p.advance()
	if err := p.tableExpectWord("PARTITIONS"); err != nil {
		return err
	}
	if _, err := p.expect('('); err != nil {
		return err
	}
	for {
		for i := 0; i < 2; i++ {
			e, err := p.parseExpr()
			if err != nil {
				return err
			}
			if e == nil {
				return p.syntaxErrorAtCur()
			}
			if i == 1 || p.cur.Type != kwTO {
				break
			}
			p.advance() // TO
		}
		if _, ok := p.match(','); !ok {
			break
		}
	}
	_, err := p.expect(')')
	return err
}

// parseAlterTableAlterColumnOptions parses the ( ONLINE = { ON | OFF } ) after
// ALTER COLUMN ... WITH, the current token being the open parenthesis: SQL
// Server refuses the other index options there (Msg 156).
func (p *Parser) parseAlterTableAlterColumnOptions() (*nodes.List, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	if _, err := p.tableOptionName(nil, "ONLINE"); err != nil {
		return nil, err
	}
	v, err := p.tableExpectOnOff()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.List{Items: []nodes.Node{&nodes.String{Str: "ONLINE"}, &nodes.String{Str: v}}}, nil
}
