package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/oracle/ast"
)

// This file reads the clauses of CREATE INDEX and ALTER INDEX by the
// grammar of the Oracle SQL Language Reference pages quoted in
// bnf/CREATE-INDEX.bnf and bnf/ALTER-INDEX.bnf. Where Oracle's answer
// differs from the page, the comment beside the reader names the answer
// (probed on the pinned Oracle Free 23.26.3).
//
// ixEnd reports whether the statement ends at the current token.
func (p *Parser) ixEnd() bool {
	return p.cur.Type == ';' || p.cur.Type == tokEOF
}

// ixWord consumes the unquoted word w, or fails at the current token.
func (p *Parser) ixWord(w string) error {
	if !p.isKeywordStr(w) {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	return nil
}

// isIxName reports whether the current token can be a name: a quoted
// identifier, or a word that is not reserved.
func (p *Parser) isIxName() bool {
	if p.cur.Type == tokQIDENT {
		return true
	}
	return p.isIdentLike() && !isOracleSQLReservedKeyword(p.cur)
}

// isIxNameTok reports whether tok can be a name, as isIxName does for the
// current token.
func (p *Parser) isIxNameTok(tok Token) bool {
	if tok.Type == tokQIDENT {
		return true
	}
	return (tok.Type == tokIDENT || tok.Type >= 2000) && !isOracleSQLReservedKeyword(tok)
}

// ixName reads a required name. A reserved word is not one: ALTER INDEX a
// RENAME PARTITION TO a is ORA-14006 at TO.
func (p *Parser) ixName() (string, error) {
	if !p.isIxName() {
		return "", p.syntaxErrorAtCur()
	}
	return p.advance().Str, nil
}

// ixObjectName reads [ schema. ] name.
func (p *Parser) ixObjectName() (*nodes.ObjectName, error) {
	obj := &nodes.ObjectName{Loc: nodes.Loc{Start: p.pos()}}
	name, err := p.ixName()
	if err != nil {
		return nil, err
	}
	obj.Name = name
	if p.cur.Type == '.' {
		p.advance()
		name2, err := p.ixName()
		if err != nil {
			return nil, err
		}
		obj.Schema, obj.Name = name, name2
	}
	obj.Loc.End = p.prev.End
	return obj, nil
}

// ixInteger reads a required integer.
func (p *Parser) ixInteger() (string, error) {
	if p.cur.Type != tokICONST {
		return "", p.syntaxErrorAtCur()
	}
	return p.advance().Str, nil
}

// ixSize reads size_clause: integer [ K | M | G | T | P | E ].
func (p *Parser) ixSize() (string, error) {
	v, err := p.ixInteger()
	if err != nil {
		return "", err
	}
	if p.cur.Type != tokQIDENT && p.isIdentLike() {
		switch p.cur.Str {
		case "K", "M", "G", "T", "P", "E":
			v += p.advance().Str
		}
	}
	return v, nil
}

// ixString reads a required string literal.
func (p *Parser) ixString() (string, error) {
	if p.cur.Type != tokSCONST {
		return "", p.syntaxErrorAtCur()
	}
	return p.advance().Str, nil
}

// ixParameters reads PARAMETERS ( 'string' ). With optional set the string
// may be left out: Oracle takes ALTER INDEX a REBUILD PARAMETERS () and
// MODIFY PARTITION a PARAMETERS (), but refuses ALTER INDEX a PARAMETERS ()
// and INDEXTYPE IS a PARAMETERS () (ORA-01780 string literal required).
func (p *Parser) ixParameters(optional bool) (string, error) {
	p.advance() // PARAMETERS
	if p.cur.Type != '(' {
		return "", p.syntaxErrorAtCur()
	}
	p.advance()
	var s string
	if p.cur.Type == tokSCONST || !optional {
		var err error
		if s, err = p.ixString(); err != nil {
			return "", err
		}
	}
	if p.cur.Type != ')' {
		return "", p.syntaxErrorAtCur()
	}
	p.advance()
	return s, nil
}

// isIxParallelStart reports whether parallel_clause starts here.
func (p *Parser) isIxParallelStart() bool {
	return p.cur.Type == kwPARALLEL || p.cur.Type == kwNOPARALLEL
}

// ixParallel reads parallel_clause: { NOPARALLEL | PARALLEL [ integer ] }.
// It returns the degree, "PARALLEL" without one, or "NOPARALLEL".
func (p *Parser) ixParallel() string {
	if p.cur.Type == kwNOPARALLEL {
		p.advance()
		return "NOPARALLEL"
	}
	p.advance() // PARALLEL
	if p.cur.Type == tokICONST {
		return p.advance().Str
	}
	return "PARALLEL"
}

// isIxCompressionStart reports whether index_compression starts here.
func (p *Parser) isIxCompressionStart() bool {
	return p.cur.Type == kwCOMPRESS || p.cur.Type == kwNOCOMPRESS
}

// ixCompression reads index_compression:
//
//	{ COMPRESS [ integer ] | COMPRESS ADVANCED [ LOW | HIGH ] | NOCOMPRESS }
//
// It returns "NOCOMPRESS", "COMPRESS" (no prefix length), the prefix
// length, or "ADVANCED [LOW|HIGH]".
func (p *Parser) ixCompression() string {
	if p.cur.Type == kwNOCOMPRESS {
		p.advance()
		return "NOCOMPRESS"
	}
	p.advance() // COMPRESS
	switch {
	case p.cur.Type == tokICONST:
		return p.advance().Str
	case p.isKeywordStr("ADVANCED"):
		p.advance()
		if p.isKeywordStr("LOW") || p.isKeywordStr("HIGH") {
			return "ADVANCED " + p.advance().Str
		}
		return "ADVANCED"
	}
	return "COMPRESS"
}

// isIxLoggingStart reports whether logging_clause starts here.
func (p *Parser) isIxLoggingStart() bool {
	return p.cur.Type == kwLOGGING || p.cur.Type == kwNOLOGGING
}

// ixTablespaceName reads the name after TABLESPACE. Oracle lets the clause
// go without one before the ',' or ')' that ends a partition description
// (LOCAL (PARTITION TABLESPACE) parses), but not elsewhere (ALTER INDEX a
// REBUILD TABLESPACE is ORA-02216). DEFAULT is taken where allowDefault is
// set.
func (p *Parser) ixTablespaceName(allowDefault bool) (string, error) {
	if allowDefault && p.cur.Type == kwDEFAULT {
		p.advance()
		return "DEFAULT", nil
	}
	if p.cur.Type == ',' || p.cur.Type == ')' {
		return "", nil
	}
	return p.ixName()
}

// ixStorageClause reads storage_clause and returns its text after STORAGE.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/storage_clause.html
//
//	STORAGE ( { INITIAL size_clause | NEXT size_clause | MINEXTENTS integer
//	          | MAXEXTENTS { integer | UNLIMITED } | maxsize_clause
//	          | PCTINCREASE integer | FREELISTS integer | FREELIST GROUPS integer
//	          | BUFFER_POOL { KEEP | RECYCLE | DEFAULT }
//	          | FLASH_CACHE { KEEP | NONE | DEFAULT }
//	          | CELL_FLASH_CACHE { KEEP | NONE | DEFAULT }
//	          | ENCRYPT } ... )
//
// An index takes no OPTIMAL (ORA-02143), and an existing one no INITIAL: ALTER INDEX a STORAGE (INITIAL 1) is
// ORA-02203 at INITIAL.
func (p *Parser) ixStorageClause(alter bool) (string, error) {
	p.advance() // STORAGE
	if p.cur.Type != '(' {
		return "", p.syntaxErrorAtCur()
	}
	start := p.pos()
	p.advance()
	n := 0
	for p.cur.Type != ')' {
		if err := p.ixStorageOption(alter); err != nil {
			return "", err
		}
		n++
	}
	if n == 0 {
		return "", p.syntaxErrorAtCur()
	}
	p.advance()
	return p.ixSpanText(start), nil
}

func (p *Parser) ixStorageOption(alter bool) error {
	if p.cur.Type == tokQIDENT || !p.isIdentLike() {
		return p.syntaxErrorAtCur()
	}
	word := p.cur.Str
	oneOf := func(words ...string) error {
		for _, w := range words {
			if p.isKeywordStr(w) {
				p.advance()
				return nil
			}
		}
		return p.syntaxErrorAtCur()
	}
	switch word {
	case "INITIAL":
		if alter {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		_, err := p.ixSize()
		return err
	case "NEXT":
		p.advance()
		_, err := p.ixSize()
		return err
	case "MINEXTENTS", "PCTINCREASE", "FREELISTS":
		p.advance()
		_, err := p.ixInteger()
		return err
	case "FREELIST":
		p.advance()
		if err := p.ixWord("GROUPS"); err != nil {
			return err
		}
		_, err := p.ixInteger()
		return err
	case "MAXEXTENTS", "MAXSIZE":
		p.advance()
		if p.isKeywordStr("UNLIMITED") {
			p.advance()
			return nil
		}
		if word == "MAXEXTENTS" {
			_, err := p.ixInteger()
			return err
		}
		_, err := p.ixSize()
		return err
	case "BUFFER_POOL":
		p.advance()
		return oneOf("KEEP", "RECYCLE", "DEFAULT")
	case "FLASH_CACHE", "CELL_FLASH_CACHE":
		p.advance()
		return oneOf("KEEP", "NONE", "DEFAULT")
	case "ENCRYPT":
		p.advance()
		return nil
	}
	return p.syntaxErrorAtCur()
}

// ixAnnotations reads annotations_clause and returns its text after
// ANNOTATIONS.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/annotations_clause.html
//
//	ANNOTATIONS ( { ADD [ IF NOT EXISTS ] | DROP [ IF EXISTS ] | REPLACE }
//	              annotation_name [ 'annotation_value' ] [, ...] )
//
// The operator may be left out (ANNOTATIONS (a) is ADD a). DROP takes no
// value (ANNOTATIONS (DROP a a) is ORA-11549), and a new index has nothing
// to drop: CREATE INDEX ... ANNOTATIONS (DROP a) is ORA-11555 at DROP.
func (p *Parser) ixAnnotations(create bool) (string, error) {
	p.advance() // ANNOTATIONS
	if p.cur.Type != '(' {
		return "", p.syntaxErrorAtCur()
	}
	start := p.pos()
	p.advance()
	for {
		drop := false
		switch {
		case p.cur.Type == kwADD:
			p.advance()
			if p.cur.Type == kwIF {
				p.advance()
				if p.cur.Type != kwNOT {
					return "", p.syntaxErrorAtCur()
				}
				p.advance()
				if p.cur.Type != kwEXISTS {
					return "", p.syntaxErrorAtCur()
				}
				p.advance()
			}
		case p.cur.Type == kwDROP:
			if create {
				return "", p.syntaxErrorAtCur()
			}
			drop = true
			p.advance()
			if p.cur.Type == kwIF {
				p.advance()
				if p.cur.Type != kwEXISTS {
					return "", p.syntaxErrorAtCur()
				}
				p.advance()
			}
		case p.isKeywordStr("REPLACE"):
			p.advance()
		}
		if _, err := p.ixName(); err != nil {
			return "", err
		}
		if !drop && p.cur.Type == tokSCONST {
			p.advance()
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if p.cur.Type != ')' {
		return "", p.syntaxErrorAtCur()
	}
	p.advance()
	return p.ixSpanText(start), nil
}

// ixInvalidation reads { DEFERRED | IMMEDIATE } INVALIDATION after the
// first word, which the caller has checked.
func (p *Parser) ixInvalidation() (string, error) {
	kind := p.advance().Str
	if err := p.ixWord("INVALIDATION"); err != nil {
		return "", err
	}
	return kind, nil
}

// isIxInvalidationStart reports whether { DEFERRED | IMMEDIATE }
// INVALIDATION starts here.
func (p *Parser) isIxInvalidationStart() bool {
	return p.cur.Type == kwDEFERRED || p.cur.Type == kwIMMEDIATE
}

// ixStoreIn reads STORE IN ( tablespace [, tablespace ]... ). Oracle takes
// a comma before the closing parenthesis: LOCAL STORE IN (a,) parses.
func (p *Parser) ixStoreIn() error {
	p.advance() // STORE
	if p.cur.Type != kwIN {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	if p.cur.Type != '(' {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	if _, err := p.ixName(); err != nil {
		return err
	}
	for p.cur.Type == ',' {
		p.advance()
		if p.cur.Type == ')' {
			break
		}
		if _, err := p.ixName(); err != nil {
			return err
		}
	}
	if p.cur.Type != ')' {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	return nil
}

// ixOnce records an option and fails at it when it was given already.
func (p *Parser) ixOnce(seen map[string]bool, key string) error {
	if seen[key] {
		return p.syntaxErrorAtCur()
	}
	seen[key] = true
	return nil
}

// ixDescBlame keeps where Oracle reports a word that cannot follow in an
// index partition description: at the partition's name when the word comes
// right after it (LOCAL (PARTITION a zzz) is ORA-14010 at a), and at a
// COMPRESS that opens the description (PARTITION COMPRESS zzz is ORA-14006
// at COMPRESS, which Oracle then reads as a name a reserved word cannot be).
type ixDescBlame struct {
	tok Token
	set bool
}

func (b *ixDescBlame) at(t Token) { b.tok, b.set = t, true }

func (b *ixDescBlame) clear() { b.set = false }

// ixDescEnd ends a partition description: a word other than ',' or ')'
// fails at the blamed token, when one is kept, or else where it stands.
func (p *Parser) ixDescEnd(b *ixDescBlame) error {
	if b.set && p.cur.Type != ',' && p.cur.Type != ')' && p.cur.Type != tokEOF {
		return p.syntaxErrorAtTok(b.tok)
	}
	return nil
}

// ---------------------------------------------------------------------------
// CREATE INDEX: index_attributes and partitioning
// ---------------------------------------------------------------------------

// isIxPhysicalStart reports whether a physical_attributes_clause option of
// an index starts here (PCTUSED is not one: an index has none).
func (p *Parser) isIxPhysicalStart() bool {
	return p.cur.Type == kwPCTFREE || p.isKeywordStr("INITRANS") ||
		p.isKeywordStr("MAXTRANS") || p.isKeywordStr("STORAGE")
}

// ixSegmentOption reads one option of segment_attributes_clause
// (physical_attributes_clause, TABLESPACE, logging_clause) of an index
// partition. It reports false when none starts here.
func (p *Parser) ixSegmentOption(seen map[string]bool) (bool, error) {
	switch {
	case p.cur.Type == kwPCTFREE || p.isKeywordStr("INITRANS") || p.isKeywordStr("MAXTRANS"):
		if err := p.ixOnce(seen, p.cur.Str); err != nil {
			return false, err
		}
		p.advance()
		_, err := p.ixInteger()
		return true, err
	case p.isKeywordStr("STORAGE"):
		if err := p.ixOnce(seen, "STORAGE"); err != nil {
			return false, err
		}
		_, err := p.ixStorageClause(false)
		return true, err
	case p.cur.Type == kwTABLESPACE:
		if err := p.ixOnce(seen, "TABLESPACE"); err != nil {
			return false, err
		}
		p.advance()
		_, err := p.ixTablespaceName(true)
		return true, err
	case p.isIxLoggingStart():
		if err := p.ixOnce(seen, "LOGGING"); err != nil {
			return false, err
		}
		p.advance()
		return true, nil
	}
	return false, nil
}

// isIxPartitionAttrStart reports whether an option of an index partition
// description starts here, so the word is not the partition's name.
func (p *Parser) isIxPartitionAttrStart() bool {
	return p.isIxPhysicalStart() || p.cur.Type == kwTABLESPACE || p.isIxLoggingStart() ||
		p.isIxCompressionStart() || p.isKeywordStr("USABLE") || p.isKeywordStr("UNUSABLE") ||
		p.isKeywordStr("STORE")
}

// ixLocalPartitioned reads local_partitioned_index after LOCAL:
//
//	LOCAL [ on_range_partitioned_table | on_list_partitioned_table
//	      | on_hash_partitioned_table | on_comp_partitioned_table ]
//
// The forms differ in what a partition takes; one reader takes their union:
//
//	( PARTITION [ partition ] [ segment_attributes_clause ] [ index_compression ]
//	  [ index_subpartition_clause ] [ USABLE | UNUSABLE ] [, ...] )
//	[ STORE IN ( tablespace [, tablespace ]... ) ]
func (p *Parser) ixLocalPartitioned() error {
	if p.isKeywordStr("STORE") {
		return p.ixStoreIn()
	}
	if p.cur.Type != '(' {
		return nil
	}
	p.advance()
	for {
		if err := p.ixLocalPartition(); err != nil {
			return err
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if p.cur.Type != ')' {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	if p.isKeywordStr("STORE") {
		return p.ixStoreIn()
	}
	return nil
}

func (p *Parser) ixLocalPartition() error {
	if p.cur.Type != kwPARTITION {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	var blame ixDescBlame
	first := true
	if p.isIxName() && !p.isIxPartitionAttrStart() {
		blame.at(p.advance())
		first = false
	}
	seen := map[string]bool{}
	for {
		if ok, err := p.ixSegmentOption(seen); err != nil {
			return err
		} else if ok {
			blame.clear()
			first = false
			continue
		}
		prior := blame
		blame.clear()
		switch {
		case p.isIxCompressionStart():
			if err := p.ixOnce(seen, "COMPRESS"); err != nil {
				return err
			}
			if first {
				blame.at(p.cur)
			}
			p.ixCompression()
		case p.isKeywordStr("USABLE") || p.isKeywordStr("UNUSABLE"):
			if err := p.ixOnce(seen, "USABLE"); err != nil {
				return err
			}
			p.advance()
		case p.isKeywordStr("STORE"):
			if err := p.ixOnce(seen, "SUBPARTITION"); err != nil {
				return err
			}
			if err := p.ixStoreIn(); err != nil {
				return err
			}
		case p.cur.Type == '(':
			if err := p.ixOnce(seen, "SUBPARTITION"); err != nil {
				return err
			}
			if err := p.ixSubpartitions(); err != nil {
				return err
			}
		default:
			return p.ixDescEnd(&prior)
		}
		first = false
	}
}

// ixSubpartitions reads the list form of index_subpartition_clause:
//
//	( SUBPARTITION [ subpartition ] [ TABLESPACE tablespace ]
//	  [ index_compression ] [ USABLE | UNUSABLE ] [, ...] )
//
// Oracle takes no index_compression here: SUBPARTITION COMPRESS reads
// COMPRESS as the reserved word it is and fails (ORA-14157).
func (p *Parser) ixSubpartitions() error {
	p.advance() // (
	for {
		if p.cur.Type != kwSUBPARTITION {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		var blame ixDescBlame
		if p.isIxName() && p.cur.Type != kwTABLESPACE && !p.isKeywordStr("USABLE") && !p.isKeywordStr("UNUSABLE") {
			blame.at(p.advance())
		}
		if p.cur.Type == kwTABLESPACE {
			blame.clear()
			p.advance()
			if _, err := p.ixTablespaceName(false); err != nil {
				return err
			}
		}
		if p.isKeywordStr("USABLE") || p.isKeywordStr("UNUSABLE") {
			blame.clear()
			p.advance()
		}
		if err := p.ixDescEnd(&blame); err != nil {
			return err
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if p.cur.Type != ')' {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	return nil
}

// ixGlobalPartitioned reads global_partitioned_index after GLOBAL. Oracle
// takes GLOBAL alone too, for an index that is not partitioned (CREATE
// INDEX i ON t (a) GLOBAL creates one).
//
//	GLOBAL PARTITION BY
//	{ RANGE ( column_list ) ( index_partitioning_clause [, ...] )
//	| HASH ( column_list ) { individual_hash_partitions | hash_partitions_by_quantity } }
func (p *Parser) ixGlobalPartitioned() error {
	if p.cur.Type != kwPARTITION {
		return nil
	}
	p.advance()
	if p.cur.Type != kwBY {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	switch {
	case p.isKeywordStr("RANGE"):
		p.advance()
		if err := p.ixPartitionColumnList(); err != nil {
			return err
		}
		if p.cur.Type != '(' {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		for {
			if err := p.ixRangePartition(); err != nil {
				return err
			}
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		if p.cur.Type != ')' {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		return nil
	case p.isKeywordStr("HASH"):
		p.advance()
		if err := p.ixColumnList(); err != nil {
			return err
		}
		if p.cur.Type == '(' {
			return p.ixHashPartitions()
		}
		// hash_partitions_by_quantity: PARTITIONS integer [ STORE IN ( ... ) ]
		if err := p.ixWord("PARTITIONS"); err != nil {
			return err
		}
		if _, err := p.ixInteger(); err != nil {
			return err
		}
		if p.isKeywordStr("STORE") {
			return p.ixStoreIn()
		}
		return nil
	}
	return p.syntaxErrorAtCur()
}

// ixPartitionColumnList reads ( column [, column ]... ) for GLOBAL PARTITION BY
// RANGE/HASH using parseExpr so a missing ')' before (PARTITION ...) stops inside.
func (p *Parser) ixPartitionColumnList() error {
	if p.cur.Type != '(' {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	for {
		if !p.isIxName() {
			return p.syntaxErrorAtCur()
		}
		col, err := p.parseExpr()
		if err != nil {
			return err
		}
		if _, ok := col.(*nodes.ColumnRef); !ok {
			return p.syntaxErrorAtCur()
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if p.cur.Type != ')' {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	return nil
}

// ixColumnList reads ( column [, column ]... ).
func (p *Parser) ixColumnList() error {
	if p.cur.Type != '(' {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	for {
		if _, err := p.ixName(); err != nil {
			return err
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if p.cur.Type != ')' {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	return nil
}

// ixRangePartition reads index_partitioning_clause:
//
//	PARTITION [ partition ] VALUES LESS THAN ( value_list ) [ segment_attributes_clause ]
//
// A bound is a literal or MAXVALUE: VALUES LESS THAN (a) is ORA-14019 after
// the a.
func (p *Parser) ixRangePartition() error {
	if p.cur.Type != kwPARTITION {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	if p.cur.Type != kwVALUES && p.isIxName() {
		p.advance()
	}
	if p.cur.Type != kwVALUES {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	if err := p.ixWord("LESS"); err != nil {
		return err
	}
	if err := p.ixWord("THAN"); err != nil {
		return err
	}
	if p.cur.Type != '(' {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	for {
		if err := p.ixBoundValue(); err != nil {
			return err
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if p.cur.Type != ')' {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	seen := map[string]bool{}
	for {
		ok, err := p.ixSegmentOption(seen)
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}
	}
}

// ixBoundValue reads one partition bound: a literal, an expression of
// literals such as TO_DATE('...', '...'), or MAXVALUE, but not a column
// or NULL (ORA-14019 after it).
func (p *Parser) ixBoundValue() error {
	if p.cur.Type == kwMAXVALUE {
		p.advance()
		return nil
	}
	expr, err := p.parseExpr()
	if err != nil {
		return err
	}
	if expr == nil {
		return p.syntaxErrorAtCur()
	}
	switch expr.(type) {
	case *nodes.ColumnRef, *nodes.NullLiteral:
		return p.syntaxErrorAtCur()
	}
	return nil
}

// ixHashPartitions reads individual_hash_partitions:
//
//	( PARTITION [ partition ] [ partitioning_storage_clause ] [, ...] )
//
// partitioning_storage_clause is [ TABLESPACE tablespace ] [ index_compression ].
func (p *Parser) ixHashPartitions() error {
	p.advance() // (
	for {
		if p.cur.Type != kwPARTITION {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		var blame ixDescBlame
		first := true
		if p.isIxName() && p.cur.Type != kwTABLESPACE && !p.isIxCompressionStart() {
			blame.at(p.advance())
			first = false
		}
		if p.cur.Type == kwTABLESPACE {
			blame.clear()
			first = false
			p.advance()
			if _, err := p.ixTablespaceName(false); err != nil {
				return err
			}
		}
		if p.isIxCompressionStart() {
			blame.clear()
			if first {
				blame.at(p.cur)
			}
			p.ixCompression()
		}
		if err := p.ixDescEnd(&blame); err != nil {
			return err
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if p.cur.Type != ')' {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	return nil
}

// ixDomainLocal reads local_domain_index_clause (and local_xmlindex_clause)
// after LOCAL:
//
//	LOCAL [ ( PARTITION [ partition ] [ PARAMETERS ( 'string' ) ] [, ...] ) ]
func (p *Parser) ixDomainLocal() error {
	if p.cur.Type != '(' {
		return nil
	}
	p.advance()
	for {
		if p.cur.Type != kwPARTITION {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		// Oracle reports a word that cannot follow at the name, a '(' at
		// what it opens, and a word after PARAMETERS ( 'string' ) at the
		// string: LOCAL (PARTITION PARENT ('s')) is ORA-14150 at 's'.
		var blame ixDescBlame
		if p.isIxName() && !p.isKeywordStr("PARAMETERS") {
			blame.at(p.advance())
		}
		if p.cur.Type == '(' {
			p.advance()
			return p.syntaxErrorAtCur()
		}
		if p.isKeywordStr("PARAMETERS") {
			p.advance()
			if p.cur.Type != '(' {
				return p.syntaxErrorAtCur()
			}
			p.advance()
			str := p.cur
			if _, err := p.ixString(); err != nil {
				return err
			}
			if p.cur.Type != ')' {
				if p.cur.Type == tokEOF {
					return p.syntaxErrorAtTok(str)
				}
				return p.syntaxErrorAtCur()
			}
			p.advance()
			blame.at(str)
		}
		if err := p.ixDescEnd(&blame); err != nil {
			return err
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if p.cur.Type != ')' {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	return nil
}

// parseCreateIndexAttributes reads index_properties of a table index, or
// index_attributes of a cluster index:
//
//	index_properties:
//	    [ global_partitioned_index | local_partitioned_index ]
//	    [ index_attributes ]
//	    [ domain_index_clause | XMLIndex_clause ]
//
//	index_attributes:
//	    [ physical_attributes_clause ] [ TABLESPACE { tablespace | DEFAULT } ]
//	    [ index_compression ] [ partial_index_clause ] [ parallel_clause ]
//	    [ { SORT | NOSORT } ] [ REVERSE ] [ { VISIBLE | INVISIBLE } ]
//	    [ logging_clause ] [ ONLINE ] [ annotations_clause ]
//
// Oracle takes the attributes in any order, each once. A cluster index
// takes no partitioning, domain clause, COMPRESS or NOSORT (ORA-02158,
// ORA-01716 at the word); NOCOMPRESS it takes.
func (p *Parser) parseCreateIndexAttributes(stmt *nodes.CreateIndexStmt) error {
	cluster := stmt.Cluster != nil
	seen := map[string]bool{}
	for {
		optStart := p.pos()
		switch {
		case p.cur.Type == kwPCTFREE || p.isKeywordStr("INITRANS") || p.isKeywordStr("MAXTRANS"):
			word := p.cur.Str
			if err := p.ixOnce(seen, word); err != nil {
				return err
			}
			p.advance()
			v, err := p.ixInteger()
			if err != nil {
				return err
			}
			switch word {
			case "PCTFREE":
				stmt.PctFree = v
			case "INITRANS":
				stmt.InitTrans = v
			default:
				stmt.MaxTrans = v
			}
		case p.isKeywordStr("STORAGE"):
			if err := p.ixOnce(seen, "STORAGE"); err != nil {
				return err
			}
			value, err := p.ixStorageClause(false)
			if err != nil {
				return err
			}
			stmt.StorageSpec = value
			appendCreateIndexOption(stmt, "STORAGE", value, nodes.Loc{Start: optStart, End: p.prev.End})
		case p.cur.Type == kwTABLESPACE:
			if err := p.ixOnce(seen, "TABLESPACE"); err != nil {
				return err
			}
			p.advance()
			name, err := p.ixTablespaceName(true)
			if err != nil {
				return err
			}
			stmt.Tablespace = name
		case p.cur.Type == kwCOMPRESS && cluster:
			return p.syntaxErrorAtCur()
		case p.isIxCompressionStart():
			if err := p.ixOnce(seen, "COMPRESS"); err != nil {
				return err
			}
			stmt.Compress = p.ixCompression()
		case p.isKeywordStr("INDEXING"):
			if err := p.ixOnce(seen, "INDEXING"); err != nil {
				return err
			}
			p.advance()
			switch {
			case p.isKeywordStr("FULL"):
				stmt.IndexingFull = true
				p.advance()
			case p.isKeywordStr("PARTIAL"):
				stmt.IndexingPartial = true
				p.advance()
			default:
				return p.syntaxErrorAtCur()
			}
		case p.isIxParallelStart():
			if err := p.ixOnce(seen, "PARALLEL"); err != nil {
				return err
			}
			stmt.Parallel = p.ixParallel()
		case p.isKeywordStr("SORT") || p.isKeywordStr("NOSORT"):
			if cluster && p.isKeywordStr("NOSORT") {
				return p.syntaxErrorAtCur()
			}
			if err := p.ixOnce(seen, "SORT"); err != nil {
				return err
			}
			stmt.Sort = p.cur.Str == "SORT"
			stmt.NoSort = !stmt.Sort
			p.advance()
		case p.cur.Type == kwREVERSE:
			if err := p.ixOnce(seen, "REVERSE"); err != nil {
				return err
			}
			stmt.Reverse = true
			p.advance()
		case p.isKeywordStr("VISIBLE") || p.cur.Type == kwINVISIBLE:
			if err := p.ixOnce(seen, "VISIBLE"); err != nil {
				return err
			}
			stmt.Visible = p.cur.Type != kwINVISIBLE
			stmt.Invisible = !stmt.Visible
			p.advance()
		case p.isIxLoggingStart():
			if err := p.ixOnce(seen, "LOGGING"); err != nil {
				return err
			}
			stmt.Logging = p.cur.Type == kwLOGGING
			stmt.NoLogging = !stmt.Logging
			p.advance()
		case p.cur.Type == kwONLINE:
			if err := p.ixOnce(seen, "ONLINE"); err != nil {
				return err
			}
			stmt.Online = true
			p.advance()
		case p.isKeywordStr("ANNOTATIONS"):
			if err := p.ixOnce(seen, "ANNOTATIONS"); err != nil {
				return err
			}
			value, err := p.ixAnnotations(true)
			if err != nil {
				return err
			}
			stmt.AnnotationsSpec = value
			appendCreateIndexOption(stmt, "ANNOTATIONS", value, nodes.Loc{Start: optStart, End: p.prev.End})
		case p.isKeywordStr("INMEMORY"), p.isKeywordStr("NOMAPPING"):
			p.advance()
			return p.syntaxErrorAtCur()
		case cluster:
			return nil
		case p.cur.Type == kwLOCAL && stmt.IndexType != nil:
			if err := p.ixOnce(seen, "DOMAIN LOCAL"); err != nil {
				return err
			}
			p.advance()
			valueStart := p.pos()
			if err := p.ixDomainLocal(); err != nil {
				return err
			}
			stmt.Local = true
			stmt.LocalPartition = p.ixSpanText(valueStart)
			appendCreateIndexOption(stmt, "LOCAL", stmt.LocalPartition, nodes.Loc{Start: optStart, End: p.prev.End})
		case p.cur.Type == kwLOCAL || p.cur.Type == kwGLOBAL:
			if err := p.ixOnce(seen, "PARTITIONING"); err != nil {
				return err
			}
			global := p.cur.Type == kwGLOBAL
			p.advance()
			valueStart := p.pos()
			var err error
			if global {
				err = p.ixGlobalPartitioned()
			} else {
				err = p.ixLocalPartitioned()
			}
			if err != nil {
				return err
			}
			value := p.ixSpanText(valueStart)
			if global {
				stmt.Global = true
				stmt.GlobalPartition = value
				appendCreateIndexOption(stmt, "GLOBAL", value, nodes.Loc{Start: optStart, End: p.prev.End})
			} else {
				stmt.Local = true
				stmt.LocalPartition = value
				appendCreateIndexOption(stmt, "LOCAL", value, nodes.Loc{Start: optStart, End: p.prev.End})
			}
		case p.isKeywordStr("INDEXTYPE"):
			// domain_index_clause / XMLIndex_clause:
			//	INDEXTYPE IS [ schema. ] indextype [ local_domain_index_clause ]
			//	[ parallel_clause ] [ PARAMETERS ( 'string' ) ]
			if err := p.ixOnce(seen, "INDEXTYPE"); err != nil {
				return err
			}
			p.advance()
			if p.cur.Type != kwIS {
				return p.syntaxErrorAtCur()
			}
			p.advance()
			name, err := p.ixObjectName()
			if err != nil {
				return err
			}
			stmt.IndexType = name
		case p.isKeywordStr("PARAMETERS") && stmt.IndexType != nil:
			if err := p.ixOnce(seen, "PARAMETERS"); err != nil {
				return err
			}
			s, err := p.ixParameters(false)
			if err != nil {
				return err
			}
			stmt.Parameters = s
		default:
			return nil
		}
	}
}

// ixSpanText returns the tokens from start to the last token read, joined
// by spaces as the option values of a CREATE INDEX record them, or "" when
// nothing was read.
func (p *Parser) ixSpanText(start int) string {
	if p.prev.End <= start {
		return ""
	}
	lex := NewLexerRange(p.source, start, p.prev.End)
	var words []string
	for tok := lex.NextToken(); tok.Type != tokEOF; tok = lex.NextToken() {
		words = append(words, p.ddlOptionTokenText(tok))
	}
	return strings.Join(words, " ")
}

// ---------------------------------------------------------------------------
// ALTER INDEX
// ---------------------------------------------------------------------------

// parseAlterIndexStmt parses an ALTER INDEX statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-INDEX.html
//
//	ALTER INDEX [ IF EXISTS ] [ schema. ] index_name
//	    { deallocate_unused_clause | allocate_extent_clause | shrink_clause
//	    | parallel_clause | physical_attributes_clause | logging_clause
//	    | partial_index_clause | rebuild_clause | alter_index_partitioning
//	    | PARAMETERS ( 'odci_parameters' ) | { DEFERRED | IMMEDIATE } INVALIDATION
//	    | COMPILE | ENABLE | DISABLE | UNUSABLE [ ONLINE ] | { VISIBLE | INVISIBLE }
//	    | RENAME TO new_index_name | COALESCE [ CLEANUP [ ONLY ] ] [ parallel_clause ]
//	    | MONITORING USAGE | NOMONITORING USAGE | UPDATE BLOCK REFERENCES
//	    | annotations_clause }
//
// Where Oracle answers otherwise than the page: it refuses IF EXISTS
// (ORA-11557 at IF) and USABLE (ORA-02243), and takes PCTFREE only inside
// REBUILD, MODIFY DEFAULT ATTRIBUTES and SPLIT PARTITION, as the page's
// text says (ALTER INDEX a PCTFREE 1 is ORA-02243 at PCTFREE). Outside the
// partitioning clauses it takes parallel_clause, logging_clause, INITRANS,
// MAXTRANS, storage_clause and the extent clauses together, each once
// (ALTER INDEX a LOGGING PARALLEL parses).
func (p *Parser) parseAlterIndexStmt(start int) (nodes.StmtNode, error) {
	p.advance() // consume INDEX
	stmt := &nodes.AlterIndexStmt{Loc: nodes.Loc{Start: start}}
	if p.cur.Type == kwIF && p.peekNext().Type == kwEXISTS {
		return nil, p.syntaxErrorAtCur()
	}
	name, err := p.ixObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	if err := p.parseAlterIndexAction(stmt); err != nil {
		return nil, err
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

func (p *Parser) parseAlterIndexAction(stmt *nodes.AlterIndexStmt) error {
	word := ""
	if p.cur.Type != tokQIDENT && p.isIdentLike() {
		word = p.cur.Str
	}
	switch word {
	case "REBUILD":
		return p.parseAlterIndexRebuild(stmt)
	case "RENAME":
		return p.parseAlterIndexRename(stmt)
	case "COALESCE":
		stmt.Action = "COALESCE"
		p.advance()
		if p.cur.Type == kwPARTITION {
			stmt.Action = "COALESCE_PARTITION"
			p.advance()
		} else if p.isKeywordStr("CLEANUP") {
			stmt.Cleanup = true
			p.advance()
			if p.isKeywordStr("ONLY") {
				stmt.CleanupOnly = true
				p.advance()
			}
		}
		if p.isIxParallelStart() {
			p.ixAlterParallel(stmt)
		}
		return nil
	case "MONITORING", "NOMONITORING":
		stmt.Action = word + "_USAGE"
		p.advance()
		return p.ixWord("USAGE")
	case "UNUSABLE":
		stmt.Action = "UNUSABLE"
		p.advance()
		if p.cur.Type == kwONLINE {
			stmt.Online = true
			p.advance()
		}
		if p.isIxInvalidationStart() {
			inv, err := p.ixInvalidation()
			if err != nil {
				return err
			}
			stmt.Invalidation = inv
		}
		return nil
	case "VISIBLE", "INVISIBLE", "ENABLE", "DISABLE", "COMPILE":
		stmt.Action = word
		p.advance()
		return nil
	case "SHRINK":
		stmt.Action = "SHRINK_SPACE"
		p.advance()
		if err := p.ixWord("SPACE"); err != nil {
			return err
		}
		if p.isKeywordStr("COMPACT") {
			stmt.Compact = true
			p.advance()
		}
		if p.cur.Type == kwCASCADE {
			stmt.Cascade = true
			p.advance()
		}
		return nil
	case "DEFERRED", "IMMEDIATE":
		stmt.Action = "INVALIDATION"
		inv, err := p.ixInvalidation()
		stmt.Invalidation = inv
		return err
	case "PARAMETERS":
		stmt.Action = "PARAMETERS"
		s, err := p.ixParameters(false)
		stmt.Parameters = s
		return err
	case "UPDATE":
		stmt.Action = "UPDATE_BLOCK_REFERENCES"
		p.advance()
		if err := p.ixWord("BLOCK"); err != nil {
			return err
		}
		return p.ixWord("REFERENCES")
	case "INDEXING":
		stmt.Action = "INDEXING"
		p.advance()
		return p.ixAlterIndexing(stmt)
	case "ANNOTATIONS":
		stmt.Action = "ANNOTATIONS"
		_, err := p.ixAnnotations(false)
		return err
	case "MODIFY":
		return p.parseAlterIndexModify(stmt)
	case "ADD":
		return p.parseAlterIndexAddPartition(stmt)
	case "DROP":
		stmt.Action = "DROP_PARTITION"
		p.advance()
		if p.cur.Type != kwPARTITION {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		name, err := p.ixName()
		stmt.Partition = name
		return err
	case "SPLIT":
		return p.parseAlterIndexSplit(stmt)
	case "PARTITION":
		// Oracle reads PARTITION partition and NOMAPPING here, then refuses
		// what follows, the end of the statement too (ORA-02243): ALTER
		// INDEX a PARTITION a TO a fails at TO.
		p.advance()
		if _, err := p.ixName(); err != nil {
			return err
		}
		return p.syntaxErrorAtCur()
	case "NOMAPPING", "INMEMORY":
		p.advance()
		return p.syntaxErrorAtCur()
	}
	return p.parseAlterIndexOptions(stmt)
}

// ixAlterParallel reads parallel_clause into an ALTER INDEX.
func (p *Parser) ixAlterParallel(stmt *nodes.AlterIndexStmt) {
	switch v := p.ixParallel(); v {
	case "NOPARALLEL":
		stmt.NoParallel = true
	case "PARALLEL":
	default:
		stmt.Parallel = v
	}
}

// ixAlterCompression reads index_compression into an ALTER INDEX.
func (p *Parser) ixAlterCompression(stmt *nodes.AlterIndexStmt) {
	switch v := p.ixCompression(); v {
	case "NOCOMPRESS":
		stmt.NoCompress = true
	case "COMPRESS":
		stmt.Compress = "1"
	default:
		stmt.Compress = v
	}
}

// ixAlterIndexing reads { FULL | PARTIAL } after INDEXING.
func (p *Parser) ixAlterIndexing(stmt *nodes.AlterIndexStmt) error {
	switch {
	case p.isKeywordStr("FULL"):
		stmt.IndexingFull = true
		p.advance()
	case p.isKeywordStr("PARTIAL"):
		stmt.IndexingPartial = true
		p.advance()
	default:
		return p.syntaxErrorAtCur()
	}
	return nil
}

// parseAlterIndexOptions reads the clauses ALTER INDEX takes together, each
// once: parallel_clause, logging_clause, INITRANS, MAXTRANS, storage_clause,
// deallocate_unused_clause and allocate_extent_clause.
func (p *Parser) parseAlterIndexOptions(stmt *nodes.AlterIndexStmt) error {
	seen := map[string]bool{}
	for !p.ixEnd() {
		ok, err := p.ixAlterSegmentOption(stmt, seen)
		if err != nil {
			return err
		}
		if ok {
			continue
		}
		if p.isIxParallelStart() {
			if err := p.ixOnce(seen, "PARALLEL"); err != nil {
				return err
			}
			if stmt.Action == "" {
				stmt.Action = p.cur.Str
			}
			p.ixAlterParallel(stmt)
			continue
		}
		break
	}
	if len(seen) == 0 {
		return p.syntaxErrorAtCur()
	}
	return nil
}

// ixAlterSegmentOption reads one of logging_clause, INITRANS, MAXTRANS,
// storage_clause, deallocate_unused_clause and allocate_extent_clause into
// an ALTER INDEX. It reports false when none starts here.
func (p *Parser) ixAlterSegmentOption(stmt *nodes.AlterIndexStmt, seen map[string]bool) (bool, error) {
	setAction := func(a string) {
		if stmt.Action == "" {
			stmt.Action = a
		}
	}
	switch {
	case p.isIxLoggingStart():
		if err := p.ixOnce(seen, "LOGGING"); err != nil {
			return false, err
		}
		setAction(p.cur.Str)
		stmt.Logging = p.cur.Type == kwLOGGING
		stmt.NoLogging = !stmt.Logging
		p.advance()
	case p.isKeywordStr("INITRANS") || p.isKeywordStr("MAXTRANS"):
		word := p.cur.Str
		if err := p.ixOnce(seen, word); err != nil {
			return false, err
		}
		setAction("PHYSICAL_ATTRIBUTES")
		p.advance()
		v, err := p.ixInteger()
		if err != nil {
			return false, err
		}
		if word == "INITRANS" {
			stmt.InitTrans = v
		} else {
			stmt.MaxTrans = v
		}
	case p.isKeywordStr("STORAGE"):
		if err := p.ixOnce(seen, "STORAGE"); err != nil {
			return false, err
		}
		setAction("PHYSICAL_ATTRIBUTES")
		if _, err := p.ixStorageClause(true); err != nil {
			return false, err
		}
	case p.isKeywordStr("DEALLOCATE"):
		if err := p.ixOnce(seen, "DEALLOCATE"); err != nil {
			return false, err
		}
		setAction("DEALLOCATE_UNUSED")
		keep, err := p.ixDeallocateUnused()
		if err != nil {
			return false, err
		}
		stmt.DeallocateKeep = keep
	case p.isKeywordStr("ALLOCATE"):
		if err := p.ixOnce(seen, "ALLOCATE"); err != nil {
			return false, err
		}
		setAction("ALLOCATE_EXTENT")
		if err := p.ixAllocateExtent(); err != nil {
			return false, err
		}
	default:
		return false, nil
	}
	return true, nil
}

// ixDeallocateUnused reads deallocate_unused_clause:
//
//	DEALLOCATE UNUSED [ KEEP size_clause ]
func (p *Parser) ixDeallocateUnused() (string, error) {
	p.advance() // DEALLOCATE
	if err := p.ixWord("UNUSED"); err != nil {
		return "", err
	}
	if p.cur.Type != kwKEEP {
		return "", nil
	}
	p.advance()
	return p.ixSize()
}

// ixAllocateExtent reads allocate_extent_clause:
//
//	ALLOCATE EXTENT [ ( { SIZE size_clause | DATAFILE 'filename' | INSTANCE integer } ... ) ]
//
// Each option is taken once (SIZE 1 SIZE 1 is ORA-03278 at the second).
func (p *Parser) ixAllocateExtent() error {
	p.advance() // ALLOCATE
	if err := p.ixWord("EXTENT"); err != nil {
		return err
	}
	if p.cur.Type != '(' {
		return nil
	}
	p.advance()
	seen := map[string]bool{}
	for p.cur.Type != ')' || len(seen) == 0 {
		var err error
		switch {
		case p.cur.Type == kwSIZE:
			if err = p.ixOnce(seen, "SIZE"); err == nil {
				p.advance()
				_, err = p.ixSize()
			}
		case p.isKeywordStr("DATAFILE"):
			if err = p.ixOnce(seen, "DATAFILE"); err == nil {
				p.advance()
				_, err = p.ixString()
			}
		case p.isKeywordStr("INSTANCE"):
			if err = p.ixOnce(seen, "INSTANCE"); err == nil {
				p.advance()
				_, err = p.ixInteger()
			}
		default:
			err = p.syntaxErrorAtCur()
		}
		if err != nil {
			return err
		}
	}
	p.advance()
	return nil
}

// parseAlterIndexRebuild reads rebuild_clause:
//
//	REBUILD [ PARTITION partition | SUBPARTITION subpartition ]
//	    [ ONLINE ] [ TABLESPACE tablespace ] [ PARAMETERS ( 'odci_parameters' ) ]
//	    [ { REVERSE | NOREVERSE } ] [ parallel_clause ] [ physical_attributes_clause ]
//	    [ index_compression ] [ logging_clause ] [ partial_index_clause ]
//	    [ { DEFERRED | IMMEDIATE } INVALIDATION ]
//
// Oracle takes the options after the partition in any order, each once
// (REBUILD ONLINE ONLINE is ORA-02243 at the second).
func (p *Parser) parseAlterIndexRebuild(stmt *nodes.AlterIndexStmt) error {
	stmt.Action = "REBUILD"
	p.advance() // REBUILD
	if p.cur.Type == kwPARTITION || p.cur.Type == kwSUBPARTITION {
		sub := p.cur.Type == kwSUBPARTITION
		p.advance()
		name, err := p.ixName()
		if err != nil {
			return err
		}
		if sub {
			stmt.Subpartition = name
		} else {
			stmt.Partition = name
		}
	}
	seen := map[string]bool{}
	for !p.ixEnd() {
		var err error
		switch {
		case p.cur.Type == kwONLINE:
			if err = p.ixOnce(seen, "ONLINE"); err == nil {
				stmt.Online = true
				p.advance()
			}
		case p.cur.Type == kwTABLESPACE:
			if err = p.ixOnce(seen, "TABLESPACE"); err == nil {
				p.advance()
				stmt.Tablespace, err = p.ixTablespaceName(false)
			}
		case p.isKeywordStr("PARAMETERS"):
			if err = p.ixOnce(seen, "PARAMETERS"); err == nil {
				stmt.Parameters, err = p.ixParameters(true)
			}
		case p.cur.Type == kwREVERSE || p.isKeywordStr("NOREVERSE"):
			if err = p.ixOnce(seen, "REVERSE"); err == nil {
				stmt.Reverse = p.cur.Type == kwREVERSE
				stmt.NoReverse = !stmt.Reverse
				p.advance()
			}
		case p.isIxParallelStart():
			if err = p.ixOnce(seen, "PARALLEL"); err == nil {
				p.ixAlterParallel(stmt)
			}
		case p.cur.Type == kwPCTFREE || p.isKeywordStr("INITRANS") || p.isKeywordStr("MAXTRANS"):
			word := p.cur.Str
			if err = p.ixOnce(seen, word); err == nil {
				p.advance()
				var v string
				v, err = p.ixInteger()
				switch word {
				case "PCTFREE":
					stmt.PctFree = v
				case "INITRANS":
					stmt.InitTrans = v
				default:
					stmt.MaxTrans = v
				}
			}
		case p.isKeywordStr("STORAGE"):
			if err = p.ixOnce(seen, "STORAGE"); err == nil {
				_, err = p.ixStorageClause(false)
			}
		case p.isIxCompressionStart():
			if err = p.ixOnce(seen, "COMPRESS"); err == nil {
				p.ixAlterCompression(stmt)
			}
		case p.isIxLoggingStart():
			if err = p.ixOnce(seen, "LOGGING"); err == nil {
				stmt.Logging = p.cur.Type == kwLOGGING
				stmt.NoLogging = !stmt.Logging
				p.advance()
			}
		case p.isKeywordStr("INDEXING"):
			if err = p.ixOnce(seen, "INDEXING"); err == nil {
				p.advance()
				err = p.ixAlterIndexing(stmt)
			}
		case p.isIxInvalidationStart():
			if err = p.ixOnce(seen, "INVALIDATION"); err == nil {
				stmt.Invalidation, err = p.ixInvalidation()
			}
		default:
			return p.syntaxErrorAtCur()
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// parseAlterIndexRename reads RENAME TO new_index_name and
// rename_index_partition:
//
//	RENAME { PARTITION partition | SUBPARTITION subpartition } TO new_name
func (p *Parser) parseAlterIndexRename(stmt *nodes.AlterIndexStmt) error {
	p.advance() // RENAME
	switch p.cur.Type {
	case kwTO:
		stmt.Action = "RENAME"
		p.advance()
		name, err := p.ixName()
		stmt.NewName = name
		return err
	case kwPARTITION, kwSUBPARTITION:
		sub := p.cur.Type == kwSUBPARTITION
		if sub {
			stmt.Action = "RENAME_SUBPARTITION"
		} else {
			stmt.Action = "RENAME_PARTITION"
		}
		p.advance()
		name, err := p.ixName()
		if err != nil {
			return err
		}
		if sub {
			stmt.Subpartition = name
		} else {
			stmt.Partition = name
		}
		if p.cur.Type != kwTO {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		stmt.NewName, err = p.ixName()
		return err
	}
	return p.syntaxErrorAtCur()
}

// parseAlterIndexModify reads modify_index_default_attrs,
// modify_index_partition and modify_index_subpartition.
func (p *Parser) parseAlterIndexModify(stmt *nodes.AlterIndexStmt) error {
	p.advance() // MODIFY
	switch p.cur.Type {
	case kwDEFAULT:
		return p.parseAlterIndexModifyDefault(stmt)
	case kwPARTITION:
		return p.parseAlterIndexModifyPartition(stmt)
	case kwSUBPARTITION:
		// MODIFY SUBPARTITION subpartition
		//     { UNUSABLE | allocate_extent_clause | deallocate_unused_clause }
		stmt.Action = "MODIFY_SUBPARTITION"
		p.advance()
		name, err := p.ixName()
		if err != nil {
			return err
		}
		stmt.Subpartition = name
		switch {
		case p.ixEnd():
		case p.isKeywordStr("UNUSABLE"):
			stmt.ModifyPartAction = "UNUSABLE"
			p.advance()
		case p.isKeywordStr("ALLOCATE"):
			stmt.ModifyPartAction = "ALLOCATE"
			return p.ixAllocateExtent()
		case p.isKeywordStr("DEALLOCATE"):
			stmt.ModifyPartAction = "DEALLOCATE"
			stmt.DeallocateKeep, err = p.ixDeallocateUnused()
			return err
		default:
			return p.syntaxErrorAtCur()
		}
		return nil
	}
	return p.syntaxErrorAtCur()
}

// parseAlterIndexModifyDefault reads modify_index_default_attrs:
//
//	MODIFY DEFAULT ATTRIBUTES [ FOR PARTITION partition ]
//	    [ physical_attributes_clause ] [ TABLESPACE { tablespace | DEFAULT } ]
//	    [ logging_clause ]
func (p *Parser) parseAlterIndexModifyDefault(stmt *nodes.AlterIndexStmt) error {
	stmt.Action = "MODIFY_DEFAULT_ATTRIBUTES"
	p.advance() // DEFAULT
	if err := p.ixWord("ATTRIBUTES"); err != nil {
		return err
	}
	if p.cur.Type == kwFOR {
		p.advance()
		if p.cur.Type != kwPARTITION {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		name, err := p.ixName()
		if err != nil {
			return err
		}
		stmt.ModifyDefaultFor = name
	}
	// At least one attribute: MODIFY DEFAULT ATTRIBUTES alone is ORA-02242.
	if p.ixEnd() {
		return p.syntaxErrorAtCur()
	}
	seen := map[string]bool{}
	for !p.ixEnd() {
		var err error
		switch {
		case p.cur.Type == kwPCTFREE || p.isKeywordStr("PCTUSED") || p.isKeywordStr("INITRANS") || p.isKeywordStr("MAXTRANS"):
			word := p.cur.Str
			if err = p.ixOnce(seen, word); err == nil {
				p.advance()
				var v string
				v, err = p.ixInteger()
				switch word {
				case "PCTFREE":
					stmt.PctFree = v
				case "PCTUSED":
					stmt.PctUsed = v
				case "INITRANS":
					stmt.InitTrans = v
				default:
					stmt.MaxTrans = v
				}
			}
		case p.isKeywordStr("STORAGE"):
			if err = p.ixOnce(seen, "STORAGE"); err == nil {
				_, err = p.ixStorageClause(true)
			}
		case p.cur.Type == kwTABLESPACE:
			if err = p.ixOnce(seen, "TABLESPACE"); err == nil {
				p.advance()
				stmt.Tablespace, err = p.ixTablespaceName(true)
			}
		case p.isIxLoggingStart():
			if err = p.ixOnce(seen, "LOGGING"); err == nil {
				stmt.Logging = p.cur.Type == kwLOGGING
				stmt.NoLogging = !stmt.Logging
				p.advance()
			}
		default:
			return p.syntaxErrorAtCur()
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// parseAlterIndexModifyPartition reads modify_index_partition:
//
//	MODIFY PARTITION partition
//	{ { deallocate_unused_clause | allocate_extent_clause
//	  | physical_attributes_clause | logging_clause | index_compression } ...
//	| PARAMETERS ( 'odci_parameters' ) | COALESCE [ CLEANUP [ ONLY ] ]
//	| UPDATE BLOCK REFERENCES | UNUSABLE }
//
// Oracle refuses PCTFREE and UPDATE BLOCK REFERENCES here (ORA-14050 at
// the word); the other options it takes together, each once.
func (p *Parser) parseAlterIndexModifyPartition(stmt *nodes.AlterIndexStmt) error {
	stmt.Action = "MODIFY_PARTITION"
	p.advance() // PARTITION
	name, err := p.ixName()
	if err != nil {
		return err
	}
	stmt.Partition = name
	switch {
	case p.ixEnd():
		return nil
	case p.isKeywordStr("PARAMETERS"):
		stmt.ModifyPartAction = "PARAMETERS"
		stmt.Parameters, err = p.ixParameters(true)
		return err
	case p.isKeywordStr("COALESCE"):
		stmt.ModifyPartAction = "COALESCE"
		p.advance()
		if p.isKeywordStr("CLEANUP") {
			stmt.Cleanup = true
			p.advance()
			if p.isKeywordStr("ONLY") {
				stmt.CleanupOnly = true
				p.advance()
			}
		}
		return nil
	case p.isKeywordStr("UNUSABLE"):
		stmt.ModifyPartAction = "UNUSABLE"
		p.advance()
		return nil
	}
	seen := map[string]bool{}
	for !p.ixEnd() {
		if p.isIxCompressionStart() {
			if err := p.ixOnce(seen, "COMPRESS"); err != nil {
				return err
			}
			if stmt.ModifyPartAction == "" {
				stmt.ModifyPartAction = p.cur.Str
			}
			p.ixAlterCompression(stmt)
			continue
		}
		action := p.cur.Str
		ok, err := p.ixAlterSegmentOption(stmt, seen)
		if err != nil {
			return err
		}
		if !ok {
			break
		}
		if stmt.ModifyPartAction == "" {
			stmt.ModifyPartAction = action
		}
	}
	if len(seen) == 0 {
		return p.syntaxErrorAtCur()
	}
	return nil
}

// parseAlterIndexAddPartition reads add_hash_index_partition:
//
//	ADD PARTITION [ partition ] [ TABLESPACE tablespace ]
//	    [ index_compression ] [ parallel_clause ]
func (p *Parser) parseAlterIndexAddPartition(stmt *nodes.AlterIndexStmt) error {
	stmt.Action = "ADD_PARTITION"
	p.advance() // ADD
	if p.cur.Type != kwPARTITION {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	// TABLESPACE not followed by a name is the partition's name: ADD
	// PARTITION TABLESPACE parses.
	var nameTok Token
	hasName := false
	if p.isIxName() && !p.isIxCompressionStart() && !p.isIxParallelStart() &&
		(p.cur.Type != kwTABLESPACE || !p.isIxNameTok(p.peekNext())) {
		nameTok = p.cur
		hasName = true
		stmt.AddPartitionName = p.advance().Str
	}
	afterNamePos := p.pos()
	firstAttrTok := p.cur
	if p.cur.Type == kwTABLESPACE {
		p.advance()
		name, err := p.ixTablespaceName(false)
		if err != nil {
			return err
		}
		stmt.Tablespace = name
	}
	if p.isIxCompressionStart() {
		p.ixAlterCompression(stmt)
	}
	if p.isIxParallelStart() {
		p.ixAlterParallel(stmt)
	}
	if !p.isStatementEnd() {
		if hasName && p.pos() == afterNamePos {
			return p.syntaxErrorAtTok(nameTok)
		}
		if !hasName && (firstAttrTok.Type == kwCOMPRESS || firstAttrTok.Type == kwNOCOMPRESS) {
			return p.syntaxErrorAtTok(firstAttrTok)
		}
	}
	return nil
}

// parseAlterIndexSplit reads split_index_partition:
//
//	SPLIT PARTITION partition_name_old AT ( literal [, literal ]... )
//	    [ INTO ( index_partition_description , index_partition_description ) ]
//	    [ parallel_clause ]
//
// Oracle takes one description in INTO too: INTO (PARTITION PARTITION)
// parses, the second PARTITION naming the first.
func (p *Parser) parseAlterIndexSplit(stmt *nodes.AlterIndexStmt) error {
	stmt.Action = "SPLIT_PARTITION"
	p.advance() // SPLIT
	if p.cur.Type != kwPARTITION {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	name, err := p.ixName()
	if err != nil {
		return err
	}
	stmt.SplitPartition = name
	if err := p.ixWord("AT"); err != nil {
		return err
	}
	if p.cur.Type != '(' {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	stmt.SplitValues = &nodes.List{}
	for {
		val, err := p.parseExpr()
		if err != nil {
			return err
		}
		if val == nil {
			return p.syntaxErrorAtCur()
		}
		stmt.SplitValues.Items = append(stmt.SplitValues.Items, val)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if p.cur.Type != ')' {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	if p.cur.Type == kwINTO {
		p.advance()
		if p.cur.Type != '(' {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		for {
			if err := p.ixIndexPartitionDescription(); err != nil {
				return err
			}
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		if p.cur.Type != ')' {
			return p.syntaxErrorAtCur()
		}
		p.advance()
	}
	if p.isIxParallelStart() {
		p.ixAlterParallel(stmt)
	}
	return nil
}

// ixIndexPartitionDescription reads index_partition_description:
//
//	PARTITION [ partition ] [ segment_attributes_clause ]
//	    [ index_compression ] [ USABLE | UNUSABLE ]
func (p *Parser) ixIndexPartitionDescription() error {
	if p.cur.Type != kwPARTITION {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	var blame ixDescBlame
	first := true
	if p.isIxName() && !p.isIxPartitionAttrStart() {
		blame.at(p.advance())
		first = false
	}
	seen := map[string]bool{}
	for {
		ok, err := p.ixSegmentOption(seen)
		if err != nil {
			return err
		}
		if ok {
			blame.clear()
			first = false
			continue
		}
		prior := blame
		blame.clear()
		switch {
		case p.isIxCompressionStart():
			if err := p.ixOnce(seen, "COMPRESS"); err != nil {
				return err
			}
			if first {
				blame.at(p.cur)
			}
			p.ixCompression()
		case p.isKeywordStr("USABLE") || p.isKeywordStr("UNUSABLE"):
			// Here Oracle blames an opening USABLE too (ORA-14010 at it).
			if err := p.ixOnce(seen, "USABLE"); err != nil {
				return err
			}
			if first {
				blame.at(p.cur)
			}
			p.advance()
		default:
			return p.ixDescEnd(&prior)
		}
		first = false
	}
}
