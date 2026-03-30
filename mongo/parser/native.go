package parser

import "github.com/bytebase/omni/mongo/ast"

// isNativeFunction returns true if the current token is a native mongosh function name.
func (p *Parser) isNativeFunction() bool {
	switch p.cur.Type {
	case kwSleep, kwLoad, kwPrint, kwPrintjson, kwQuit, kwExit,
		kwVersion, kwCls, kwHelp, kwIt,
		kwIsNaN, kwIsFinite, kwParseInt, kwParseFloat,
		kwEncodeURI, kwEncodeURIComponent, kwDecodeURI, kwDecodeURIComponent:
		return true
	}
	return false
}

// parseNativeFunctionCall parses a top-level native function call like sleep(1000).
func (p *Parser) parseNativeFunctionCall() (ast.Node, error) {
	return nil, p.syntaxErrorAtCur()
}

// parseIdentStatement parses statements starting with an unrecognized identifier.
func (p *Parser) parseIdentStatement() (ast.Node, error) {
	return nil, p.syntaxErrorAtCur()
}
