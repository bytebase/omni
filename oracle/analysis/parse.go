package analysis

import (
	"errors"
	"reflect"

	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/oracle/parser"
)

// Parse parses a script statement by statement, as SplitStatements splits
// it, and returns every statement in one list. Locations and a parse error's
// position are byte offsets in the whole script.
func Parse(sql string) (*ast.List, error) {
	list := &ast.List{}
	for _, statement := range SplitStatements(sql) {
		if statement.Empty {
			continue
		}
		parsed, err := parser.Parse(statement.Text)
		if err != nil {
			return nil, offsetParseError(err, statement.Start)
		}
		if parsed == nil {
			continue
		}
		offsetLocs(parsed, statement.Start)
		list.Items = append(list.Items, parsed.Items...)
	}
	return list, nil
}

func offsetParseError(err error, offset int) error {
	var parseErr *parser.ParseError
	if !errors.As(err, &parseErr) {
		return err
	}
	adjusted := *parseErr
	adjusted.Position += offset
	return &adjusted
}

type locOffsetter int

// offsetLocs reaches only the nodes ast.Walk visits, as the Bytebase code this
// was moved from did. Parsing each statement in place with parser.ParseRange
// would place every node, which changes analysis results, so it is a separate
// change.
func offsetLocs(node ast.Node, offset int) {
	if offset == 0 {
		return
	}
	ast.Walk(locOffsetter(offset), node)
}

func (o locOffsetter) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}
	value := reflect.ValueOf(node)
	if value.Kind() != reflect.Pointer || value.IsNil() {
		return o
	}
	elem := value.Elem()
	if elem.Kind() != reflect.Struct {
		return o
	}
	locField := elem.FieldByName("Loc")
	if !locField.IsValid() || !locField.CanSet() || locField.Type() != reflect.TypeOf(ast.Loc{}) {
		return o
	}

	loc, ok := locField.Interface().(ast.Loc)
	if !ok {
		return o
	}
	offset := int(o)
	if loc.Start >= 0 {
		loc.Start += offset
	}
	if loc.End >= 0 {
		loc.End += offset
	}
	locField.Set(reflect.ValueOf(loc))
	return o
}
