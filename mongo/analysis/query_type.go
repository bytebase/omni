package analysis

import (
	"errors"
	"strings"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/mongo/ast"
	"github.com/bytebase/omni/mongo/parser"
)

// QueryType is the permission class of a statement. The values mirror
// Bytebase's base.QueryType.
type QueryType = analysis.QueryType

const (
	// QueryTypeUnknown is a statement not deliberately mapped.
	QueryTypeUnknown = analysis.QueryTypeUnknown
	// Select is a read-only user query.
	Select = analysis.Select
	// Explain is a query plan statement.
	Explain = analysis.Explain
	// SelectInfoSchema reads system or metadata information.
	SelectInfoSchema = analysis.SelectInfoSchema
	// DDL changes schema or administrative state.
	DDL = analysis.DDL
	// DML changes data.
	DML = analysis.DML
)

// safeBSONHelpers lists pure BSON value constructors and static factory methods
// that do not perform database writes or arbitrary shell side effects.
var safeBSONHelpers = map[string]bool{
	"ObjectId":                          true,
	"ObjectId.createFromHexString":      true,
	"ObjectId.createFromBase64":         true,
	"ObjectId.createFromTime":           true,
	"ObjectId.isValid":                  true,
	"ObjectId.generate":                 true,
	"ISODate":                           true,
	"Date":                              true,
	"Date.now":                          true,
	"Date.UTC":                          true,
	"Date.parse":                        true,
	"UUID":                              true,
	"NumberLong":                        true,
	"NumberInt":                         true,
	"NumberDecimal":                     true,
	"Long":                              true,
	"Long.fromString":                   true,
	"Long.fromNumber":                   true,
	"Long.fromInt":                      true,
	"Long.fromBits":                     true,
	"Int32":                             true,
	"Int32.fromString":                  true,
	"Double":                            true,
	"Double.fromString":                 true,
	"Decimal128":                        true,
	"Decimal128.fromString":             true,
	"Decimal128.fromStringWithRounding": true,
	"Timestamp":                         true,
	"Timestamp.fromBits":                true,
	"Timestamp.fromInt":                 true,
	"Timestamp.fromNumber":              true,
	"Timestamp.fromString":              true,
	"RegExp":                            true,
	"BSONRegExp":                        true,
	"BSONRegExp.parseOptions":           true,
	"BinData":                           true,
	"Binary":                            true,
	"Binary.createFromBase64":           true,
	"Binary.createFromHexString":        true,
	"Binary.fromBits":                   true,
	"Binary.fromFloat32Array":           true,
	"Binary.fromInt8Array":              true,
	"Binary.fromPackedBits":             true,
	"HexData":                           true,
	"MinKey":                            true,
	"MaxKey":                            true,
	"Code":                              true,
	"DBRef":                             true,
	"Symbol":                            true,
	"BSONSymbol":                        true,
	"MD5":                               true,
	"EJSON.parse":                       true,
	"EJSON.stringify":                   true,
	"EJSON.serialize":                   true,
	"EJSON.deserialize":                 true,
}

func hasWriteOrUnknownStage(a *StatementAnalysis) bool {
	if a.UnsupportedStage == "unknown" {
		return true
	}
	for _, stage := range a.PipelineStages {
		if stage == "$out" || stage == "$merge" || stage == "unknown" {
			return true
		}
	}
	return false
}

// hasUnsafeStatementArgs reports whether a statement's argument expressions (or
// chained database calls) invoke anything outside safeBSONHelpers.
func hasUnsafeStatementArgs(node ast.Node, source string) bool {
	switch n := node.(type) {
	case *ast.CollectionStatement:
		if hasUnsafeNodes(n.Args, source) || hasUnsafeNodes(n.ExplainArgs, source) {
			return true
		}
		for _, cm := range n.CursorMethods {
			if hasUnsafeNodes(cm.Args, source) {
				return true
			}
		}
		return false
	case *ast.DatabaseStatement:
		if len(n.ChainedMethods) > 0 {
			return true
		}
		return hasUnsafeNodes(n.Args, source)
	case *ast.RsStatement:
		return hasUnsafeNodes(n.Args, source)
	case *ast.ShStatement:
		return hasUnsafeNodes(n.Args, source)
	default:
		return false
	}
}

func hasUnsafeNodes(nodes []ast.Node, source string) bool {
	stack := make([]ast.Node, 0, len(nodes))
	stack = append(stack, nodes...)
	for len(stack) > 0 {
		n := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if n == nil {
			continue
		}
		switch v := n.(type) {
		case *ast.StringLiteral, *ast.NumberLiteral, *ast.BoolLiteral,
			*ast.NullLiteral, *ast.RegexLiteral, *ast.Identifier:
			continue
		case *ast.Document:
			for _, pair := range v.Pairs {
				stack = append(stack, pair.Value)
			}
		case *ast.Array:
			stack = append(stack, v.Elements...)
		case *ast.HelperCall:
			if !safeBSONHelpers[v.Name] {
				return true
			}
			if source != "" && strings.Contains(v.Name, ".") && !isPlainStaticCall(source, v) {
				return true
			}
			stack = append(stack, v.Args...)
		default:
			return true
		}
	}
	return false
}

// isPlainStaticCall verifies that a dotted HelperCall (e.g.
// ObjectId.createFromHexString("...")) does not hide an earlier parenthesized
// call in its receiver (e.g. ObjectId(mut()).createFromHexString("...")), which
// mongo/parser collapses into a single HelperCall keeping only the last args.
func isPlainStaticCall(source string, call *ast.HelperCall) bool {
	if call.Loc.Start < 0 || call.Loc.End > len(source) || call.Loc.Start >= call.Loc.End {
		return false
	}
	text := source[call.Loc.Start:call.Loc.End]
	openIdx := strings.IndexByte(text, '(')
	if openIdx < 0 {
		return false
	}
	var b strings.Builder
	b.Grow(openIdx)
	for i := range openIdx {
		switch text[i] {
		case ' ', '\t', '\n', '\r':
		default:
			b.WriteByte(text[i])
		}
	}
	return b.String() == call.Name
}

// QueryTypeOf maps a StatementAnalysis to its QueryType.
// Nil or unrecognized analyses default to DML for safety.
func QueryTypeOf(a *StatementAnalysis) QueryType {
	if a == nil {
		return DML
	}
	switch a.Operation {
	case OpFind, OpFindOne, OpCount, OpDistinct:
		return Select
	case OpAggregate:
		if hasWriteOrUnknownStage(a) {
			return DML
		}
		return Select
	case OpRead:
		// OpRead covers runCommand({find:...}) and collection metadata reads
		// (getIndexes, stats, etc.). Treat runCommand/adminCommand as Select
		// (data access) unless a pipeline write/unknown stage is present, and
		// collection metadata methods as SelectInfoSchema.
		if a.MethodName == "runCommand" || a.MethodName == "adminCommand" {
			if hasWriteOrUnknownStage(a) {
				return DML
			}
			return Select
		}
		return SelectInfoSchema
	case OpInfo:
		// OpInfo covers show commands, db/rs/sh info methods, and info
		// runCommands, plus the fallback for unknown db methods in Analyze.
		// Classify only known info methods as SelectInfoSchema and default
		// unknowns to DML for safety.
		switch a.MethodName {
		case "show", "status",
			"getCollectionNames", "getCollectionInfos",
			"serverStatus", "serverBuildInfo", "version",
			"hostInfo", "getName", "listCommands", "stats",
			"conf", "config", "printReplicationInfo", "printSecondaryReplicationInfo",
			"getBalancerState", "isBalancerRunning":
			return SelectInfoSchema
		case "runCommand", "adminCommand":
			return SelectInfoSchema
		default:
			return DML
		}
	case OpWrite:
		return DML
	case OpAdmin:
		return DDL
	case OpExplain:
		return Explain
	default:
		return DML
	}
}

// QueryType returns the QueryType for this statement analysis.
// A nil receiver returns DML.
func (a *StatementAnalysis) QueryType() QueryType {
	return QueryTypeOf(a)
}

// ClassifyStatement parses a single MongoDB shell statement and returns its
// QueryType. Empty input returns DML with a nil error. Syntax errors return DML
// alongside the *parser.ParseError from mongo/parser, and multi-statement input
// returns DML with an error.
func ClassifyStatement(statement string) (QueryType, error) {
	nodes, err := parser.Parse(statement)
	if err != nil {
		return DML, err
	}
	if len(nodes) == 0 {
		return DML, nil
	}
	if len(nodes) > 1 {
		return DML, errors.New("expected a single MongoDB statement")
	}
	qt := QueryTypeOf(Analyze(nodes[0]))
	if (qt == Select || qt == SelectInfoSchema || qt == Explain) && hasUnsafeStatementArgs(nodes[0], statement) {
		return DML, nil
	}
	return qt, nil
}
