// Package analysis reads one T-SQL statement: its query span (the columns a
// result traces to, the columns its conditions read, the tables it accesses),
// its permission class, the statement type, the resources it changes, and the
// splitter that finds a script's batches.
//
// # What a statement reads and this package cannot trace
//
// A masking caller masks a result column by matching the columns it traces to
// against column policies. Three T-SQL shapes defeat that trace, because the
// value comes from a place the connected instance's synced metadata does not
// describe or this call does not see. For each, this package answers with what
// it knows instead of a plausible lineage, and the answer says the rows must be
// withheld rather than masked. The engine behavior behind each answer is
// recorded from the pinned SQL Server (testdata/pins.tsv, 2022 CU27 GDR,
// 16.0.4305.5) by the tests in this package that read testdata/tapes:
//
//   - A linked server read. A four-part name (srv.db.schema.table), OPENQUERY,
//     OPENROWSET and OPENDATASOURCE are resolved by the server through
//     sys.servers, against another instance's metadata: no object here
//     describes the relation, and a loopback linked server returns exactly the
//     rows a local read of the same table returns
//     (TestOracleLinkedServerRead). ResolveQuerySpan returns a
//     *LinkedServerError, so the statement is refused. Naming the relation on
//     the span with no columns is not an option: a masking caller reads that as
//     a statement that touches nothing sensitive and returns the rows. A
//     statement that only writes returns no rows, so it keeps its span: a
//     four-part name appears as the table it writes, with ColumnResource.Server
//     set, while a write that reads through OPENQUERY or OPENROWSET names only
//     its local target (TestResolveQuerySpanLinkedServerWrite). The rows of
//     EXEC ... AT, a pass-through call, are withheld by the driver, which
//     treats a procedure call as untraced.
//   - A #temp table, or a @table variable, whose rows the batch's own
//     statements wrote (SELECT ... INTO, INSERT ... VALUES/SELECT/EXEC,
//     UPDATE ... FROM, MERGE ... USING; TestOracleTempTableRows). The caller
//     carries those tables from one statement to the next as names and columns
//     (ResolveOptions.TempTables), not as lineage, and a temp table outlives
//     the GO that ends a batch. So a read of one is reported in
//     ResolvedSpan.Untraced: the result columns still point at the temp table,
//     and the field says the chain stops there. Tracing it would need every
//     write shape modelled and each column's lineage carried across calls; a
//     shape left out would claim a trace that is wrong.
//   - A local variable, which a statement of the same batch assigned from a
//     column (DECLARE @v = (SELECT ...), SET @v = (SELECT ...), SELECT @v =
//     column; TestOracleVariableRows). Variable scope is the batch, so no
//     caller-side state spans it; a statement that reads a local variable is
//     reported in ResolvedSpan.Untraced. A @@name is a system function, whose
//     value is the server's state, so it is not a variable here.
//
// A masking caller must therefore withhold the rows of a span whose Untraced
// field is set, and must treat a *LinkedServerError as a refusal, as it treats
// a statement that does not parse. A caller that masks by column position and
// ignores Untraced returns rows masking could not evaluate.
package analysis
