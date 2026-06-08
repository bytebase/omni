package parser

import (
	"strings"

	"github.com/bytebase/omni/redshift/plpgsql/ast"
)

// parseOptions parses compiler option directives before the block.
// Directives start with # and must appear before the block (before any label or DECLARE/BEGIN).
//
// Supported directives:
//   - #option dump
//   - #print_strict_params on|off
//   - #variable_conflict error|use_variable|use_column
//
// The # comes through as an Op token with Str="#" from the SQL lexer.
func (p *Parser) parseOptions() ([]ast.PLOption, error) {
	var options []ast.PLOption

	for p.isOp("#") {
		startPos := p.pos()
		p.advance() // consume #

		// Read the directive name
		if p.isEOF() {
			return nil, p.errorf("syntax error at or near end of input, expected option directive")
		}

		name := strings.ToLower(p.identText())
		p.advance() // consume directive name

		switch name {
		case "option":
			// #option <name>
			if p.isEOF() {
				return nil, p.errorf("syntax error at or near end of input, expected option name after #option")
			}
			optName := strings.ToLower(p.identText())
			p.advance()
			switch optName {
			case "dump":
				options = append(options, ast.PLOption{
					Name:  "dump",
					Value: "",
					Loc:   ast.Loc{Start: startPos, End: p.prev.End},
				})
			default:
				return nil, &ParseError{
					Message:  "unrecognized #option: " + optName,
					Position: p.prev.Loc,
				}
			}

		case "print_strict_params":
			// #print_strict_params on|off
			if p.isEOF() {
				return nil, p.errorf("syntax error at or near end of input, expected on/off after #print_strict_params")
			}
			val := strings.ToLower(p.identText())
			p.advance()
			if val != "on" && val != "off" {
				return nil, &ParseError{
					Message:  "#print_strict_params expects on or off, got " + val,
					Position: p.prev.Loc,
				}
			}
			options = append(options, ast.PLOption{
				Name:  "print_strict_params",
				Value: val,
				Loc:   ast.Loc{Start: startPos, End: p.prev.End},
			})

		case "variable_conflict":
			// #variable_conflict error|use_variable|use_column
			if p.isEOF() {
				return nil, p.errorf("syntax error at or near end of input, expected value after #variable_conflict")
			}
			val := strings.ToLower(p.identText())
			p.advance()
			switch val {
			case "error", "use_variable", "use_column":
				// valid
			default:
				return nil, &ParseError{
					Message:  "#variable_conflict expects error, use_variable, or use_column, got " + val,
					Position: p.prev.Loc,
				}
			}
			options = append(options, ast.PLOption{
				Name:  "variable_conflict",
				Value: val,
				Loc:   ast.Loc{Start: startPos, End: p.prev.End},
			})

		default:
			return nil, &ParseError{
				Message:  "unrecognized compiler directive: #" + name,
				Position: startPos,
			}
		}
	}

	return options, nil
}
