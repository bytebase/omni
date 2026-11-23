package analysis

import (
	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/googlesql/ast"
)

// StatementType is the kind of change a GoogleSQL statement makes.
type StatementType = analysis.StatementType

const (
	StatementTypeUnspecified = analysis.StatementTypeUnspecified
	Insert                   = analysis.Insert
	Update                   = analysis.Update
	Delete                   = analysis.Delete
	Merge                    = analysis.Merge
	Truncate                 = analysis.Truncate
	CreateTable              = analysis.CreateTable
	CreateView               = analysis.CreateView
	CreateIndex              = analysis.CreateIndex
	CreateSchema             = analysis.CreateSchema
	CreateDatabase           = analysis.CreateDatabase
	CreateFunction           = analysis.CreateFunction
	CreateProcedure          = analysis.CreateProcedure
	CreateSequence           = analysis.CreateSequence
	AlterTable               = analysis.AlterTable
	AlterView                = analysis.AlterView
	AlterIndex               = analysis.AlterIndex
	AlterDatabase            = analysis.AlterDatabase
	AlterSequence            = analysis.AlterSequence
	DropTable                = analysis.DropTable
	DropView                 = analysis.DropView
	DropIndex                = analysis.DropIndex
	DropSchema               = analysis.DropSchema
	DropDatabase             = analysis.DropDatabase
	DropSequence             = analysis.DropSequence
	DropFunction             = analysis.DropFunction
	DropProcedure            = analysis.DropProcedure
	Rename                   = analysis.Rename
	Execute                  = analysis.Execute
)

// StatementTypeOf returns the type of a statement.
//
// CALL and EXECUTE IMMEDIATE report Execute: what they change is decided by
// the procedure body or the run-time SQL string, which the statement does not
// show.
//
// Deferral: every GoogleSQL statement the enum has no value for reports
// StatementTypeUnspecified, so a rule selecting DML skips it. The gap includes
// statements that do write rows — LOAD DATA INTO, CLONE DATA INTO, scripting
// bodies (BEGIN…END, IF, FOR, WHILE, each of which parses as one statement, so a
// CALL or EXECUTE IMMEDIATE nested in one is not seen either). It also includes
// ALTER SCHEMA, row access policies, generic entities, and the Spanner-only DDL
// set (change streams, roles, GRANT / REVOKE, proto bundles), which do not write
// rows.
//
// A new row-writing statement needs its own enum value before it lands here.
// ClassifySQL returns Unknown for these too, so it cannot gate the list.
func StatementTypeOf(node ast.Node) StatementType {
	switch n := node.(type) {
	case *ast.SubpipelineStmt:
		return pipeStatementType(n.Pipeline)
	case *ast.PipeStatement:
		t, conflict := pipeStatementResult(n.Pipeline)
		if conflict {
			return StatementTypeUnspecified
		}
		if t != StatementTypeUnspecified {
			return t
		}
		// A reading suffix keeps the head's type: a piped CALL still executes.
		return StatementTypeOf(n.Head)
	case *ast.QueryStmt:
		// A query only makes a change through a pipe operator that writes
		// (`SELECT … |> INSERT t`); everything else is a read.
		if n.Pipe != nil {
			return pipeStatementType(n.Pipe)
		}
		return StatementTypeUnspecified
	case *ast.InsertStmt:
		return Insert
	case *ast.UpdateStmt:
		return Update
	case *ast.DeleteStmt:
		return Delete
	case *ast.MergeStmt:
		return Merge
	case *ast.TruncateStmt:
		return Truncate

	case *ast.CreateTableStmt:
		return CreateTable
	case *ast.CreateViewStmt:
		return CreateView
	case *ast.CreateMaterializedViewStmt:
		return CreateView
	case *ast.CreateIndexStmt:
		return CreateIndex
	case *ast.SearchVectorIndexStmt:
		return CreateIndex
	case *ast.CreateSchemaStmt:
		return CreateSchema
	case *ast.CreateDatabaseStmt:
		return CreateDatabase
	case *ast.CreateFunctionStmt:
		return CreateFunction
	case *ast.CreateProcedureStmt:
		return CreateProcedure
	case *ast.CreateSequenceStmt:
		return CreateSequence
	case *ast.CreateSnapshotStmt:
		return CreateTable

	case *ast.AlterStmt:
		return alterStatementType(n.Object)
	case *ast.DropStmt:
		return dropStatementType(n.Object)
	case *ast.AlterSequenceStmt:
		return AlterSequence
	case *ast.DropSequenceStmt:
		return DropSequence
	case *ast.BQAlterStmt:
		return bqAlterStatementType(n.Object)
	case *ast.BQDropStmt:
		return bqDropStatementType(n.Object)

	case *ast.RenameStmt:
		return Rename

	case *ast.CallStmt, *ast.ExecuteImmediateStmt:
		return Execute

	default:
		return StatementTypeUnspecified
	}
}

// pipeStatementType reports the change a pipe suffix makes. A pipe INSERT
// writes rows and a pipe CREATE TABLE creates one; a nested branch that does
// either makes the whole statement that kind. Every other operator only reads
// or transforms the pipeline, so it reports Unspecified like a query does.
type pipeStatementTypes struct {
	t        StatementType
	conflict bool
}

// add records one path's type.
func (a *pipeStatementTypes) add(t StatementType) {
	if t == StatementTypeUnspecified {
		return
	}
	switch {
	case a.t == StatementTypeUnspecified:
		a.t = t
	case a.t != t:
		a.conflict = true
	}
}

// addAll records the types of a nested pipeline; its own conflict carries over.
func (a *pipeStatementTypes) addAll(o pipeStatementTypes) {
	if o.conflict {
		a.conflict = true
	}
	a.add(o.t)
}

func (a pipeStatementTypes) result() StatementType {
	if a.conflict {
		return StatementTypeUnspecified
	}
	return a.t
}

// pipeStatementType reports the change a pipe suffix makes. A pipe INSERT
// writes rows and a pipe CREATE TABLE creates one; a nested branch that does
// either makes the whole statement that kind. Every other operator only reads
// or transforms the pipeline, so it reports Unspecified like a query does. Two
// different writes on different runtime paths conflict (no single type covers
// both) and report Unspecified.
func pipeStatementType(p *ast.Subpipeline) StatementType {
	t, _ := pipeStatementResult(p)
	return t
}

// pipeStatementResult reports the pipe suffix's type and whether two runtime
// paths make different changes, which no single type describes.
func pipeStatementResult(p *ast.Subpipeline) (StatementType, bool) {
	acc := pipeStatementTypesOf(p)
	return acc.result(), acc.conflict
}

func pipeStatementTypesOf(p *ast.Subpipeline) pipeStatementTypes {
	var acc pipeStatementTypes
	if p == nil {
		return acc
	}
	for _, op := range p.Ops {
		switch n := op.(type) {
		case *ast.PipeInsert:
			acc.add(Insert)
		case *ast.PipeCreateTable:
			acc.add(CreateTable)
		case *ast.PipeLog:
			acc.addAll(pipeStatementTypesOf(n.Pipeline))
		case *ast.PipeIf:
			for _, c := range n.Cases {
				acc.addAll(pipeStatementTypesOf(c.Then))
			}
			acc.addAll(pipeStatementTypesOf(n.Else))
		case *ast.PipeFork:
			for _, branch := range n.Pipelines {
				acc.addAll(pipeStatementTypesOf(branch))
			}
		case *ast.PipeTee:
			for _, branch := range n.Pipelines {
				acc.addAll(pipeStatementTypesOf(branch))
			}
		}
	}
	return acc
}

func alterStatementType(kind ast.AlterObjectKind) StatementType {
	switch kind {
	case ast.AlterTable:
		return AlterTable
	case ast.AlterView:
		return AlterView
	case ast.AlterIndex, ast.AlterSearchIndex:
		return AlterIndex
	case ast.AlterDatabase:
		return AlterDatabase
	default:
		// ALTER SCHEMA has no enum value; see the deferral on StatementTypeOf.
		return StatementTypeUnspecified
	}
}

func dropStatementType(kind ast.DropObjectKind) StatementType {
	switch kind {
	case ast.DropTable:
		return DropTable
	case ast.DropView:
		return DropView
	case ast.DropIndex:
		return DropIndex
	case ast.DropSchema:
		return DropSchema
	case ast.DropDatabase:
		return DropDatabase
	case ast.DropTableFunction:
		return DropFunction
	default:
		return StatementTypeUnspecified
	}
}

func bqAlterStatementType(kind ast.BQAlterObjectKind) StatementType {
	switch kind {
	case ast.BQAlterMaterializedView, ast.BQAlterApproxView:
		return AlterView
	case ast.BQAlterSearchIndex, ast.BQAlterVectorIndex:
		return AlterIndex
	default:
		return StatementTypeUnspecified
	}
}

func bqDropStatementType(kind ast.BQDropObjectKind) StatementType {
	switch kind {
	case ast.BQDropFunction, ast.BQDropTableFunction:
		return DropFunction
	case ast.BQDropProcedure:
		return DropProcedure
	case ast.BQDropMaterializedView, ast.BQDropApproxView:
		return DropView
	case ast.BQDropSnapshotTable:
		return DropTable
	case ast.BQDropSearchIndex, ast.BQDropVectorIndex:
		return DropIndex
	default:
		return StatementTypeUnspecified
	}
}
