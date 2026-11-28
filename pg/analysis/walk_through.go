package analysis

import (
	"context"

	"google.golang.org/protobuf/proto"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/pg/catalog"
)

var _ analysis.WalkThrougher = Engine

// WalkThrough runs the script on a catalog loaded with the database's schema
// and applies what the script changed to the database's metadata. The first
// statement the catalog refuses is the finding, placed by its line.
func (engine) WalkThrough(ctx context.Context, request analysis.WalkThroughRequest, _ []analysis.ParsedStatement) analysis.WalkThroughResult {
	if request.Script == "" {
		return analysis.WalkThroughResult{}
	}
	before := catalog.New()
	if request.SessionUser != "" {
		before.SetSessionUser(request.SessionUser)
	}
	if len(request.SearchPath) > 0 {
		before.SetSearchPath(request.SearchPath)
	}
	original := request.Database.Metadata
	if _, err := before.LoadMetadata(ctx, original, catalog.LoadMetadataOptions{Full: true}); err != nil {
		return walkThroughFailure(analysis.WalkThroughDDLSimulationFailed, "Failed to load schema", err, 0)
	}
	after := before.Clone()
	results, err := after.Exec(request.Script, &catalog.ExecOptions{ContinueOnError: true})
	if err != nil {
		return walkThroughFailure(analysis.WalkThroughDDLSimulationFailed, "DDL simulation failed", err, 0)
	}
	for _, r := range results {
		if r.Error != nil {
			return walkThroughFailure(walkThroughCode(r.Error), "DDL simulation failed", r.Error, r.Line)
		}
	}
	updated := catalog.ApplyDiffToMetadata(original, before, after, catalog.Diff(before, after))
	if proto.Equal(original, updated) {
		return analysis.WalkThroughResult{}
	}
	return analysis.WalkThroughResult{Database: &analysis.Database{
		Name:                request.Database.Name,
		Metadata:            updated,
		ObjectCaseSensitive: request.Database.ObjectCaseSensitive,
	}}
}

func walkThroughFailure(code analysis.WalkThroughCode, title string, err error, line int) analysis.WalkThroughResult {
	return analysis.WalkThroughResult{Finding: &analysis.WalkThroughFinding{
		Code:      code,
		Error:     true,
		Title:     title,
		Content:   err.Error(),
		Statement: -1,
		Start:     analysis.Position{Line: line},
	}}
}

// walkThroughCode maps a catalog error's SQLSTATE to a walk-through code.
func walkThroughCode(err error) analysis.WalkThroughCode {
	catalogErr, ok := err.(*catalog.Error)
	if !ok {
		return analysis.WalkThroughDDLSimulationFailed
	}
	switch catalogErr.Code {
	case catalog.CodeDuplicateTable:
		return analysis.WalkThroughTableExists
	case catalog.CodeUndefinedTable:
		return analysis.WalkThroughTableNotExists
	case catalog.CodeDuplicateColumn:
		return analysis.WalkThroughColumnExists
	case catalog.CodeUndefinedColumn, catalog.CodeInvalidColumnDefinition:
		return analysis.WalkThroughColumnNotExists
	case catalog.CodeDuplicateSchema:
		return analysis.WalkThroughRelationExists
	case catalog.CodeUndefinedSchema:
		return analysis.WalkThroughSchemaNotExists
	case catalog.CodeDuplicateObject, catalog.CodeUniqueViolation:
		return analysis.WalkThroughIndexExists
	case catalog.CodeUndefinedObject:
		return analysis.WalkThroughIndexNotExists
	case catalog.CodeInvalidFK, catalog.CodeForeignKeyViolation:
		return analysis.WalkThroughTableHasFK
	case catalog.CodeDuplicatePKey:
		return analysis.WalkThroughPrimaryKeyExists
	case catalog.CodeDatatypeMismatch:
		return analysis.WalkThroughChangeColumnType
	case catalog.CodeNotNullViolation:
		return analysis.WalkThroughColumnCannotNull
	case catalog.CodeSyntaxError:
		return analysis.WalkThroughSyntaxError
	case catalog.CodeReservedName:
		return analysis.WalkThroughNameIsKeywordIdentifier
	default:
		return analysis.WalkThroughDDLSimulationFailed
	}
}
