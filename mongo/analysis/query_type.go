package analysis

import (
	"errors"

	"github.com/bytebase/omni/analysis"
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
	return QueryTypeOf(Analyze(nodes[0])), nil
}
