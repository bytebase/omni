package parser

import (
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// Bulk data movement — COPY INTO <table> (load) + COPY INTO <location> (unload)
// (T5.2)
// ---------------------------------------------------------------------------
//
// Snowflake's COPY grammar carries large, version-growing option vocabularies
// (copyOptions, FILE_FORMAT type-options, transfer options). Rather than mirror
// the legacy ANTLR grammar's finite, already-stale enumerations (its
// copy_options rule lacks MAX_FILE_SIZE, SINGLE, HEADER as an option,
// INCLUDE_QUERY_ID, INCLUDE_METADATA, LOAD_MODE, CLUSTER_AT_INGEST_TIME, ... all
// of which appear in the official docs corpus), every option that follows the
// COPY source/destination is parsed as an open-ended `KEY = <value>` pair
// (ast.CopyOption). Structural keywords (INTO, FROM, PARTITION BY) anchor the
// grammar; option-name vocabulary is not enumerated. The catalog/semantic layer
// validates that an option is real and legal. This mirrors the merged
// GRANT/REVOKE (T6.1) and SHOW (T6.3) open-ended token-run approach.
//
// Two location-lexing wrinkles drive the helpers below:
//   - External cloud-storage URIs ('s3://...', 'gcs://...', 'azure://...') are
//     single-quoted, so they lex as a plain string token — easy.
//   - PUT/GET file:// paths are NOT quoted and contain '//', which omni's lexer
//     treats as a line comment. They cannot be reconstructed from tokens and are
//     scanned directly from the source text (scanRawLocation); the lexer is then
//     re-synced past the scanned span.
//   - '@'-stage references lex as a clean run of contiguous tokens
//     (@ % ~ / . ident), so they are rebuilt from tokens by source contiguity.

// parseCopyStmt parses a COPY statement and dispatches on the destination:
//
//	COPY INTO <table>    ...    (load)   → parseCopyIntoTable
//	COPY INTO <location> ...    (unload) → parseCopyIntoLocation
//
// The disambiguation is purely structural: an unload target is a stage ('@...')
// or an external URI (a single-quoted 'scheme://...' string); anything else is
// a table name and therefore a load. The COPY keyword has NOT been consumed.
func (p *Parser) parseCopyStmt() (ast.Node, error) {
	copyTok := p.advance() // consume COPY
	if _, err := p.expect(kwINTO); err != nil {
		return nil, err
	}

	if p.startsStageLocation() {
		return p.parseCopyIntoLocation(copyTok.Loc)
	}
	return p.parseCopyIntoTable(copyTok.Loc)
}

// ---------------------------------------------------------------------------
// COPY INTO <table> (load)
// ---------------------------------------------------------------------------

// parseCopyIntoTable parses the load form. INTO has been consumed; cur is the
// target table name. start is the Loc of the COPY keyword.
//
//	COPY INTO [ns.]table [ (cols) ]
//	  FROM { stage | location | ( SELECT ...$col... FROM stage ) }
//	  [ FILES = (...) ] [ PATTERN = '...' ] [ FILE_FORMAT = (...) ]
//	  [ copyOptions ] [ VALIDATION_MODE = ... ]
func (p *Parser) parseCopyIntoTable(start ast.Loc) (ast.Node, error) {
	target, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}

	stmt := &ast.CopyIntoTableStmt{
		Target: target,
		Loc:    ast.Loc{Start: start.Start},
	}

	// Optional ( col, ... ) column list (transformation form). Distinguished
	// from anything else by the '(' immediately following the table name.
	if p.cur.Type == '(' {
		p.advance() // consume '('
		cols, err := p.parseIdentList()
		if err != nil {
			return nil, err
		}
		stmt.Columns = cols
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}

	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}

	// Source: a transformation query `( SELECT ... FROM stage )` or a bare
	// stage/location. The '(' followed by SELECT (or WITH) is the transform.
	if p.cur.Type == '(' && (p.peekNext().Type == kwSELECT || p.peekNext().Type == kwWITH) {
		p.advance() // consume '('
		sel, err := p.parseCopyTransformQuery()
		if err != nil {
			return nil, err
		}
		stmt.Transform = sel
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	} else {
		from, err := p.parseStageLocation()
		if err != nil {
			return nil, err
		}
		stmt.From = from
	}

	// Trailing options: FILES / PATTERN / FILE_FORMAT / copyOptions /
	// VALIDATION_MODE — all open-ended KEY = value pairs.
	opts, err := p.parseCopyOptions()
	if err != nil {
		return nil, err
	}
	stmt.Options = opts

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseCopyTransformQuery parses the SELECT inside a COPY INTO <table> FROM (
// SELECT ... ) transformation. The transform's FROM source is a stage (e.g.
// `FROM @mystage/file.csv d`), which the general SELECT parser cannot consume,
// so this restricted reader is used in its place. The opening '(' has already
// been consumed; the closing ')' is consumed by the caller.
//
// Per the docs the transform body is:
//
//	SELECT [<alias>.]$<col>[.<element>] [ , ... ] FROM { internalStage | externalStage }
//
// but the corpus also shows `FROM TABLE(<fn>(...))` (a streaming/table-function
// source) and full expressions in the select list (e.g. `$1:num::number`). To
// stay robust and reuse the dependency's machinery, the select list and (when
// the source is not a stage) the FROM clause are delegated to the shared SELECT
// expression/table-ref parsers; only a stage source is read specially here.
func (p *Parser) parseCopyTransformQuery() (*ast.SelectStmt, error) {
	selTok, err := p.expect(kwSELECT)
	if err != nil {
		return nil, err
	}

	stmt := &ast.SelectStmt{Loc: ast.Loc{Start: selTok.Loc.Start}}

	// SELECT list (reuses the shared expression parser via parseSelectList).
	targets, err := p.parseSelectList()
	if err != nil {
		return nil, err
	}
	stmt.Targets = targets

	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}

	// FROM source. A stage ('@...') is recorded as a TableRef whose Name is a
	// synthetic single-part ObjectName carrying the raw stage text (TableRef has
	// no dedicated stage field, and the shared primary-source parser cannot
	// consume an '@' stage). Non-stage sources (a table function, a parenthesized
	// subquery, a plain name) reuse the shared primary-source parser.
	if p.cur.Type == '@' {
		stage, err := p.parseStageRef()
		if err != nil {
			return nil, err
		}
		ref := &ast.TableRef{
			Name: &ast.ObjectName{
				Name: ast.Ident{Name: stage.Raw},
				Loc:  stage.Loc,
			},
			Loc: stage.Loc,
		}
		// Optional alias after the stage (e.g. `@stage d`).
		if alias, ok := p.parseOptionalAlias(); ok {
			ref.Alias = alias
			ref.Loc.End = p.prev.Loc.End
		}
		stmt.From = []ast.Node{ref}
	} else {
		src, err := p.parsePrimarySource()
		if err != nil {
			return nil, err
		}
		stmt.From = []ast.Node{src}
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// COPY INTO <location> (unload)
// ---------------------------------------------------------------------------

// parseCopyIntoLocation parses the unload form. INTO has been consumed; cur is
// the destination stage/location. start is the Loc of the COPY keyword.
//
//	COPY INTO { stage | location }
//	  FROM { [ns.]table | ( query ) }
//	  [ PARTITION BY <expr> ] [ FILE_FORMAT = (...) ] [ copyOptions ]
//	  [ VALIDATION_MODE = RETURN_ROWS ] [ HEADER ]
func (p *Parser) parseCopyIntoLocation(start ast.Loc) (ast.Node, error) {
	into, err := p.parseStageLocation()
	if err != nil {
		return nil, err
	}

	stmt := &ast.CopyIntoLocationStmt{
		Into: into,
		Loc:  ast.Loc{Start: start.Start},
	}

	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}

	// FROM source: a parenthesized query or a table name.
	if p.cur.Type == '(' {
		p.advance() // consume '('
		var query ast.Node
		if p.cur.Type == kwWITH {
			query, err = p.parseWithQueryExpr()
		} else {
			query, err = p.parseQueryExpr()
		}
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		stmt.FromQuery = query
	} else {
		name, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		stmt.FromTable = name
	}

	// Optional PARTITION BY <expr>.
	if p.cur.Type == kwPARTITION {
		p.advance() // consume PARTITION
		if _, err := p.expect(kwBY); err != nil {
			return nil, err
		}
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		stmt.Partition = expr
	}

	// Trailing options: FILE_FORMAT / copyOptions / VALIDATION_MODE / HEADER.
	opts, err := p.parseCopyOptions()
	if err != nil {
		return nil, err
	}
	stmt.Options = opts

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// Shared: stage / location parsing
// ---------------------------------------------------------------------------

// startsStageLocation reports whether the current token begins a stage or
// external-location reference (as opposed to a table name). A stage starts with
// '@'; an external location is a single-quoted 'scheme://...' string. A bare
// identifier is a table name (load target / unload source), not a location.
func (p *Parser) startsStageLocation() bool {
	if p.cur.Type == '@' {
		return true
	}
	if p.cur.Type == tokString {
		return isExternalURI(p.cur.Str)
	}
	return false
}

// isExternalURI reports whether s looks like a cloud-storage URI, i.e. it
// contains a "scheme://" prefix (s3://, gcs://, azure://, and any future
// scheme). Used to tell an external location string apart from an ordinary
// string literal.
func isExternalURI(s string) bool {
	i := strings.Index(s, "://")
	if i <= 0 {
		return false
	}
	// Everything before "://" must be a bare scheme (letters/digits/+/-/.).
	for _, c := range s[:i] {
		if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' ||
			c == '+' || c == '-' || c == '.') {
			return false
		}
	}
	return true
}

// parseStageLocation parses a stage reference or an external-location URI in a
// COPY source/destination position:
//
//	@[ns.]name[/path] | @[ns.]%table[/path] | @~[/path]      → StageRef
//	'scheme://<...>'                                          → StageExternal
//
// (PUT/GET local file paths are parsed by parseLocalFile, not here, because
// they need raw-source scanning past the '//' comment hazard.)
func (p *Parser) parseStageLocation() (*ast.StageLocation, error) {
	if p.cur.Type == '@' {
		return p.parseStageRef()
	}
	if p.cur.Type == tokString && isExternalURI(p.cur.Str) {
		tok := p.advance()
		return &ast.StageLocation{
			Kind: ast.StageExternal,
			Raw:  tok.Str,
			Loc:  tok.Loc,
		}, nil
	}
	return nil, p.syntaxErrorAtCur()
}

// parseStageRef parses an '@'-prefixed stage reference by rebuilding it from the
// contiguous token run the lexer produced (@ then any of % ~ / . and name
// tokens, plus '$' for path segments). The run ends at the first token that is
// not source-adjacent to the previous one, or that cannot be part of a stage
// path (whitespace-separated tokens, commas, parens, clause keywords, EOF).
//
// On entry cur is '@'. Returns a StageLocation whose Raw is the verbatim source
// span "@...".
func (p *Parser) parseStageRef() (*ast.StageLocation, error) {
	atTok := p.advance() // consume '@'
	startAbs := atTok.Loc.Start
	endAbs := atTok.Loc.End

	// Greedily consume tokens that are (a) adjacent in the source to the prior
	// token (no intervening whitespace) and (b) valid stage-path constituents.
	for p.cur.Loc.Start == endAbs && p.isStagePathToken(p.cur.Type) {
		endAbs = p.cur.Loc.End
		p.advance()
	}

	return &ast.StageLocation{
		Kind: ast.StageRef,
		Raw:  p.srcSlice(startAbs, endAbs),
		Loc:  ast.Loc{Start: startAbs, End: endAbs},
	}, nil
}

// isStagePathToken reports whether a token type may appear inside a stage path
// after the leading '@'. Stage paths are composed of name parts and the
// separators % ~ / . — plus '$' which the lexer can emit when a path segment is
// scanned as a $-variable (rare but harmless to allow). Structural punctuation
// (comma, parens), operators, and EOF terminate the path.
func (p *Parser) isStagePathToken(tokType int) bool {
	switch tokType {
	case '%', '~', '/', '.', '$':
		return true
	case tokIdent, tokQuotedIdent, tokVariable:
		return true
	}
	// Numbers can appear in path segments (e.g. @stage/2020/data).
	if tokType == tokInt || tokType == tokFloat || tokType == tokReal {
		return true
	}
	// Keywords may appear as path segment names (the lexer keywordizes words
	// like DATA, RESULT, etc.). Any keyword constant is >= 700.
	return tokType >= 700
}

// ---------------------------------------------------------------------------
// Shared: local file paths (PUT source / GET target)
// ---------------------------------------------------------------------------

// parseLocalFile parses a PUT source / GET target local-file location, which is
// one of:
//
//	file://<path>            (unquoted; may contain '//' that omni's lexer would
//	                          otherwise treat as a line comment)
//	'file://<path>'          (single-quoted → a plain string token)
//	<bare-path>              (e.g. the GET target `my_target_path`)
//
// The unquoted file:// form is scanned directly from the source text and the
// lexer is re-synced past it; the other two forms are read from the token
// stream.
func (p *Parser) parseLocalFile() (*ast.StageLocation, error) {
	// Quoted form: 'file://...' (or any quoted local path).
	if p.cur.Type == tokString {
		tok := p.advance()
		return &ast.StageLocation{
			Kind: ast.StageLocalFile,
			Raw:  tok.Str,
			Loc:  tok.Loc,
		}, nil
	}

	// Unquoted file://... form: 'file' lexes as an identifier/keyword, then ':'
	// and '//...' would be a line comment. Scan the whole run from source.
	if p.curIsWord("FILE") && p.lookaheadColonSlashSlash() {
		return p.scanRawLocation()
	}

	// Bare local path (e.g. GET target `my_target_path`): a single identifier.
	if p.cur.Type == tokIdent || p.cur.Type == tokQuotedIdent {
		tok := p.advance()
		return &ast.StageLocation{
			Kind: ast.StageLocalFile,
			Raw:  p.srcSlice(tok.Loc.Start, tok.Loc.End),
			Loc:  tok.Loc,
		}, nil
	}

	return nil, p.syntaxErrorAtCur()
}

// curIsWord reports whether the current token is an identifier or keyword whose
// uppercased text equals upper.
func (p *Parser) curIsWord(upper string) bool {
	if p.cur.Str == "" {
		return false
	}
	return strings.ToUpper(p.cur.Str) == upper
}

// lookaheadColonSlashSlash reports whether the source immediately following the
// current token is "://" — i.e. the current 'file' word begins a file:// URI.
// It inspects the raw source rather than tokens because '//' is lexed as a
// comment and would not appear as a token.
func (p *Parser) lookaheadColonSlashSlash() bool {
	// Position just past the current token, in input-relative coordinates.
	i := p.cur.Loc.End - p.base
	return i >= 0 && i+3 <= len(p.input) && p.input[i:i+3] == "://"
}

// scanRawLocation scans an unquoted file:// location directly from the source,
// starting at the current token, consuming everything up to the next ASCII
// whitespace byte (or end of segment), then re-syncs the lexer to that point.
// file:// paths never contain spaces (a path with spaces is single-quoted), so
// whitespace is the correct terminator.
func (p *Parser) scanRawLocation() (*ast.StageLocation, error) {
	startAbs := p.cur.Loc.Start
	i := startAbs - p.base // input-relative scan cursor
	for i < len(p.input) {
		c := p.input[i]
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
			break
		}
		i++
	}
	endAbs := i + p.base
	raw := p.input[startAbs-p.base : i]

	// Re-sync the lexer to endAbs (input-relative) and re-prime cur. Because
	// scanRawLocation bypassed the token stream, both the buffered lookahead and
	// the current token must be discarded and re-read from the new position.
	p.lexer.pos = i
	p.hasNext = false
	p.advance()

	return &ast.StageLocation{
		Kind: ast.StageLocalFile,
		Raw:  raw,
		Loc:  ast.Loc{Start: startAbs, End: endAbs},
	}, nil
}

// srcSlice returns the verbatim source text spanning the absolute byte range
// [startAbs, endAbs). Token Locs are absolute (input-relative + base), while
// p.input is the input-relative segment, so the base offset is subtracted.
func (p *Parser) srcSlice(startAbs, endAbs int) string {
	lo := startAbs - p.base
	hi := endAbs - p.base
	if lo < 0 {
		lo = 0
	}
	if hi > len(p.input) {
		hi = len(p.input)
	}
	if lo > hi {
		return ""
	}
	return p.input[lo:hi]
}

// ---------------------------------------------------------------------------
// Shared: open-ended copy options (KEY = value)
// ---------------------------------------------------------------------------

// parseCopyOptions parses a run of zero or more COPY options. Each option is a
// `KEY = <value>` pair (or a bare value-less keyword such as HEADER). The run
// continues while the current token begins another option, and terminates at a
// statement boundary (';' / EOF) or any token that does not begin an option.
func (p *Parser) parseCopyOptions() ([]*ast.CopyOption, error) {
	var opts []*ast.CopyOption
	for p.startsCopyOption() {
		opt, err := p.parseCopyOption()
		if err != nil {
			return nil, err
		}
		opts = append(opts, opt)
	}
	return opts, nil
}

// startsCopyOption reports whether the current token can begin a COPY option.
// An option name is an identifier or a keyword (the lexer keywordizes many
// option names: FILES, PATTERN, FILE_FORMAT, PARTITION, ON_ERROR, ...).
// Structural tokens (parens, commas, operators) and EOF do not begin an option.
func (p *Parser) startsCopyOption() bool {
	if p.cur.Type == tokIdent || p.cur.Type == tokQuotedIdent {
		return true
	}
	if p.cur.Type == tokEOF || p.cur.Type == ';' {
		return false
	}
	return p.cur.Type >= 700
}

// parseCopyOption parses one `KEY = <value>` (or bare `KEY`) option.
//
// The name is one or more name words joined by '_' is not needed (option names
// are single tokens). The value is one of:
//   - ( k = v [, ...] )      → nested option group (FILE_FORMAT, CREDENTIALS,
//     INCLUDE_METADATA, ENCRYPTION, ...)
//   - ( 'a' [, 'b' ...] )    → literal list (FILES = ('f1', 'f2'))
//   - ( 'a' [, 'b' ...] )    NULL_IF style string lists also fall here
//   - '<string>' | <number>  → literal value
//   - <word-run>             → verbatim uppercased token run (TRUE, FALSE,
//     CASE_INSENSITIVE, RETURN_ERRORS, a bare name, ...)
//
// A name with no following '=' is a bare option (e.g. HEADER).
func (p *Parser) parseCopyOption() (*ast.CopyOption, error) {
	nameTok := p.advance() // consume the option name
	opt := &ast.CopyOption{
		Name: strings.ToUpper(nameTok.Str),
		Loc:  ast.Loc{Start: nameTok.Loc.Start, End: nameTok.Loc.End},
	}

	// Bare option (no '='): e.g. trailing HEADER.
	if p.cur.Type != '=' {
		opt.Bare = true
		return opt, nil
	}
	p.advance() // consume '='

	// Parenthesized value: a nested group or a literal list.
	if p.cur.Type == '(' {
		if err := p.parseCopyOptionParen(opt); err != nil {
			return nil, err
		}
		opt.Loc.End = p.prev.Loc.End
		return opt, nil
	}

	// Literal value: string, number, or a double-quoted value. A double-quoted
	// token is lexically an identifier but, as an option value, is semantically a
	// string (Snowflake accepts e.g. PATTERN = "regex"); it is recorded as a
	// string literal carrying the quoted text.
	switch p.cur.Type {
	case tokString:
		tok := p.advance()
		opt.Lit = &ast.Literal{Kind: ast.LitString, Value: tok.Str, Loc: tok.Loc}
		opt.Loc.End = tok.Loc.End
		return opt, nil
	case tokQuotedIdent:
		tok := p.advance()
		opt.Lit = &ast.Literal{Kind: ast.LitString, Value: tok.Str, Loc: tok.Loc}
		opt.Loc.End = tok.Loc.End
		return opt, nil
	case tokInt:
		tok := p.advance()
		opt.Lit = &ast.Literal{Kind: ast.LitInt, Value: tok.Str, Ival: tok.Ival, Loc: tok.Loc}
		opt.Loc.End = tok.Loc.End
		return opt, nil
	case tokFloat, tokReal:
		tok := p.advance()
		opt.Lit = &ast.Literal{Kind: ast.LitFloat, Value: tok.Str, Loc: tok.Loc}
		opt.Loc.End = tok.Loc.End
		return opt, nil
	}

	// Word value: a single name word (TRUE / CASE_INSENSITIVE / a bare
	// integration name / a dotted format name), captured verbatim and uppercased.
	words, end, err := p.parseValueWord()
	if err != nil {
		return nil, err
	}
	opt.Words = words
	opt.Loc.End = end
	return opt, nil
}

// parseCopyOptionParen parses the parenthesized RHS of an option, which is
// either a key/value group ( k = v [, ...] ) or a literal list ( 'a' [, ...] ).
// It peeks past the '(' to decide: a list begins with a literal followed by ','
// or ')'; otherwise it is a key/value group. On entry cur is '('.
func (p *Parser) parseCopyOptionParen(opt *ast.CopyOption) error {
	p.advance() // consume '('

	// Empty group: ().
	if p.cur.Type == ')' {
		p.advance()
		opt.Group = []*ast.CopyOption{} // non-nil to record "()" explicitly
		return nil
	}

	// A literal list: first token is a literal and the token after it is ',' or
	// ')'. (FILES = ('f1', 'f2'); NULL_IF = ('NULL', 'null').)
	if p.curIsLiteral() {
		next := p.peekNext()
		if next.Type == ',' || next.Type == ')' {
			lits, err := p.parseLiteralList()
			if err != nil {
				return err
			}
			opt.List = lits
			if _, err := p.expect(')'); err != nil {
				return err
			}
			return nil
		}
	}

	// Otherwise a nested key/value group. Snowflake option groups separate
	// entries with WHITESPACE, not commas (e.g.
	// FILE_FORMAT = (FORMAT_NAME ='x' COMPRESSION='GZIP'),
	// CREDENTIALS = (AWS_KEY_ID='k' AWS_SECRET_KEY='s' AWS_TOKEN='t')),
	// though a comma between entries is also tolerated. The loop therefore
	// continues while another `key [= value]` entry begins, optionally
	// consuming a separating comma.
	var group []*ast.CopyOption
	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		sub, err := p.parseGroupEntry()
		if err != nil {
			return err
		}
		group = append(group, sub)
		if p.cur.Type == ',' {
			p.advance() // consume optional ',' separator
		}
	}
	if _, err := p.expect(')'); err != nil {
		return err
	}
	opt.Group = group
	return nil
}

// parseGroupEntry parses one `key = value` entry inside a parenthesized option
// group. The value reuses the same value-shapes as a top-level option (nested
// group, literal list, literal, or single word). Entries are separated by
// whitespace (and, tolerantly, an optional comma) by the caller's loop.
// INCLUDE_METADATA uses `key = METADATA$col` values, captured as a word value
// (the lexer emits a single METADATA$col identifier).
func (p *Parser) parseGroupEntry() (*ast.CopyOption, error) {
	if !p.startsCopyOption() {
		return nil, p.syntaxErrorAtCur()
	}
	nameTok := p.advance()
	entry := &ast.CopyOption{
		Name: strings.ToUpper(nameTok.Str),
		Loc:  ast.Loc{Start: nameTok.Loc.Start, End: nameTok.Loc.End},
	}

	// Group entries are always `key = value` (a bare key inside a group is not
	// a documented form), but tolerate a bare key defensively rather than error.
	if p.cur.Type != '=' {
		entry.Bare = true
		return entry, nil
	}
	p.advance() // consume '='

	if p.cur.Type == '(' {
		if err := p.parseCopyOptionParen(entry); err != nil {
			return nil, err
		}
		entry.Loc.End = p.prev.Loc.End
		return entry, nil
	}
	switch p.cur.Type {
	case tokString:
		tok := p.advance()
		entry.Lit = &ast.Literal{Kind: ast.LitString, Value: tok.Str, Loc: tok.Loc}
		entry.Loc.End = tok.Loc.End
		return entry, nil
	case tokInt:
		tok := p.advance()
		entry.Lit = &ast.Literal{Kind: ast.LitInt, Value: tok.Str, Ival: tok.Ival, Loc: tok.Loc}
		entry.Loc.End = tok.Loc.End
		return entry, nil
	case tokFloat, tokReal:
		tok := p.advance()
		entry.Lit = &ast.Literal{Kind: ast.LitFloat, Value: tok.Str, Loc: tok.Loc}
		entry.Loc.End = tok.Loc.End
		return entry, nil
	}
	words, end, err := p.parseValueWord()
	if err != nil {
		return nil, err
	}
	entry.Words = words
	entry.Loc.End = end
	return entry, nil
}

// curIsLiteral reports whether the current token is a string or numeric literal.
func (p *Parser) curIsLiteral() bool {
	switch p.cur.Type {
	case tokString, tokInt, tokFloat, tokReal:
		return true
	}
	return false
}

// parseLiteralList parses a comma-separated list of literals: lit [, lit ...].
// The opening '(' has been consumed; the closing ')' is consumed by the caller.
func (p *Parser) parseLiteralList() ([]*ast.Literal, error) {
	var lits []*ast.Literal
	for {
		lit, err := p.parseCopyLiteral()
		if err != nil {
			return nil, err
		}
		lits = append(lits, lit)
		if p.cur.Type == ',' {
			p.advance() // consume ','
			continue
		}
		break
	}
	return lits, nil
}

// parseCopyLiteral parses one string/number literal in a COPY option list.
func (p *Parser) parseCopyLiteral() (*ast.Literal, error) {
	switch p.cur.Type {
	case tokString:
		tok := p.advance()
		return &ast.Literal{Kind: ast.LitString, Value: tok.Str, Loc: tok.Loc}, nil
	case tokInt:
		tok := p.advance()
		return &ast.Literal{Kind: ast.LitInt, Value: tok.Str, Ival: tok.Ival, Loc: tok.Loc}, nil
	case tokFloat, tokReal:
		tok := p.advance()
		return &ast.Literal{Kind: ast.LitFloat, Value: tok.Str, Loc: tok.Loc}, nil
	}
	return nil, p.syntaxErrorAtCur()
}

// parseValueWord reads a single word-run option value: one name word plus any
// adjacency-joined '.' / '$' continuations (no intervening whitespace), captured
// verbatim and uppercased. Returns the joined text and the end offset.
//
// COPY/PUT/GET option values are always a single token (TRUE, CASE_INSENSITIVE,
// CONTINUE, RETURN_ERRORS, a bare integration name), optionally extended by
// source-adjacent dotted or '$' segments for dotted format names
// (my_db.my_format) and metadata pseudo-columns (METADATA$FILENAME). Crucially,
// a following *space-separated* word is NOT absorbed: it is either the next
// `KEY = value` option/entry or a trailing bare option (e.g. a final HEADER).
// This is the same single-token rule used by both top-level options and group
// entries, so both call this one helper.
func (p *Parser) parseValueWord() (string, int, error) {
	if !p.isOptionWord(p.cur.Type) {
		return "", 0, p.syntaxErrorAtCur()
	}
	var b strings.Builder
	first := p.advance()
	b.WriteString(strings.ToUpper(first.Str))
	end := first.Loc.End

	// Adjacency-joined dotted / $ continuations only (no spaces).
	for (p.cur.Type == '.' || p.cur.Type == '$') && p.cur.Loc.Start == end {
		sep := p.advance()
		b.WriteString(p.srcSlice(sep.Loc.Start, sep.Loc.End))
		end = sep.Loc.End
		if p.isOptionWord(p.cur.Type) && p.cur.Loc.Start == end {
			part := p.advance()
			b.WriteString(strings.ToUpper(part.Str))
			end = part.Loc.End
		}
	}
	return b.String(), end, nil
}

// isOptionWord reports whether a token type may appear inside an option value
// word run. Identifiers, quoted identifiers, keywords, and bare numbers all
// qualify (a value like RETURN_10_ROWS is one keyword token; CASE_INSENSITIVE,
// TRUE, a bare integration name are identifiers/keywords). Structural
// punctuation, operators, and EOF do not.
func (p *Parser) isOptionWord(tokType int) bool {
	switch tokType {
	case tokIdent, tokQuotedIdent:
		return true
	case tokInt, tokFloat, tokReal:
		return true
	}
	return tokType >= 700
}
