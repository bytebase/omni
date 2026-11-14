package deparse

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// CREATE TABLE
// ---------------------------------------------------------------------------

func (w *writer) writeCreateTableStmt(n *ast.CreateTableStmt) error {
	w.buf.WriteString("CREATE")
	if n.OrReplace {
		w.buf.WriteString(" OR REPLACE")
	} else if n.OrAlter {
		w.buf.WriteString(" OR ALTER")
	}
	if n.Transient {
		w.buf.WriteString(" TRANSIENT")
	} else if n.Temporary {
		w.buf.WriteString(" TEMPORARY")
		if n.ReadOnly {
			w.buf.WriteString(" READ ONLY")
		}
	} else if n.Volatile {
		w.buf.WriteString(" VOLATILE")
	}
	if n.Hybrid {
		w.buf.WriteString(" HYBRID")
	}
	if n.Iceberg {
		w.buf.WriteString(" ICEBERG")
	}
	if n.Interactive {
		w.buf.WriteString(" INTERACTIVE")
	}
	w.buf.WriteString(" TABLE")
	if n.IfNotExists {
		w.buf.WriteString(" IF NOT EXISTS")
	}
	w.buf.WriteByte(' ')
	w.writeObjectNameNoSpace(n.Name)

	// FROM BACKUP SET / FROM ARCHIVE OF
	if n.FromBackup != nil {
		w.writeBackupSource(n.FromBackup)
		return nil
	}
	if a := n.FromArchive; a != nil {
		w.buf.WriteString(" FROM ARCHIVE OF ")
		w.writeObjectNameNoSpace(a.Table)
		if !a.Alias.IsEmpty() {
			w.buf.WriteString(" AS ")
			w.buf.WriteString(a.Alias.String())
		}
		w.buf.WriteString(" WHERE ")
		return writeExprNoLeadSpace(w, a.Where)
	}

	// [ COPY GRANTS ] USING TEMPLATE ( <query> )
	if n.UsingTemplate != nil {
		if n.CopyGrants {
			w.buf.WriteString(" COPY GRANTS")
		}
		w.buf.WriteString(" USING TEMPLATE (")
		if err := w.writeQueryExpr(n.UsingTemplate); err != nil {
			return err
		}
		w.buf.WriteByte(')')
		return nil
	}

	switch {
	case n.Like != nil:
		w.buf.WriteString(" LIKE ")
		w.writeObjectNameNoSpace(n.Like)
	case n.Clone != nil:
		w.buf.WriteString(" CLONE ")
		w.writeObjectNameNoSpace(n.Clone.Source)
		if err := w.writeCloneTimeTravelSuffix(n.Clone); err != nil {
			return err
		}
	case len(n.Columns) > 0 || len(n.Constraints) > 0 || len(n.Indexes) > 0:
		// Column definitions, table constraints and inline INDEX elements (HYBRID).
		w.buf.WriteString(" (")
		idx := 0
		for _, col := range n.Columns {
			if idx > 0 {
				w.buf.WriteString(", ")
			}
			if err := w.writeColumnDef(col); err != nil {
				return err
			}
			idx++
		}
		for _, con := range n.Constraints {
			if idx > 0 {
				w.buf.WriteString(", ")
			}
			if err := w.writeTableConstraint(con); err != nil {
				return err
			}
			idx++
		}
		for _, ti := range n.Indexes {
			if idx > 0 {
				w.buf.WriteString(", ")
			}
			w.writeTableIndex(ti)
			idx++
		}
		w.buf.WriteByte(')')
	}

	// PARTITION BY (Iceberg)
	if len(n.PartitionBy) > 0 {
		w.buf.WriteString(" PARTITION BY (")
		for i, expr := range n.PartitionBy {
			if i > 0 {
				w.buf.WriteString(", ")
			}
			if err := writeExprNoLeadSpace(w, expr); err != nil {
				return err
			}
		}
		w.buf.WriteByte(')')
	}

	// CLUSTER BY
	if err := w.writeClusterBy(n.ClusterBy, n.Linear); err != nil {
		return err
	}

	// Iceberg catalog and storage parameters, COPY = ( ... ), and an
	// interactive table's TARGET_LAG / WAREHOUSE / INITIALIZATION_WAREHOUSE
	if err := w.writeCopyOptions(n.Options); err != nil {
		return err
	}

	// COPY GRANTS / COPY TAGS / PURGE
	if n.CopyGrants {
		w.buf.WriteString(" COPY GRANTS")
	}
	if n.CopyTags {
		w.buf.WriteString(" COPY TAGS")
	}
	if n.Purge {
		w.buf.WriteString(" PURGE")
	}

	// COMMENT
	if n.Comment != nil {
		w.buf.WriteString(" COMMENT = ")
		w.buf.WriteString(quoteString(*n.Comment))
	}

	// WITH ROW ACCESS POLICY
	if n.RowPolicy != nil {
		w.buf.WriteString(" WITH ROW ACCESS POLICY ")
		w.writeObjectNameNoSpace(n.RowPolicy.PolicyName)
		if len(n.RowPolicy.Columns) > 0 {
			w.buf.WriteString(" ON (")
			for i, col := range n.RowPolicy.Columns {
				if i > 0 {
					w.buf.WriteString(", ")
				}
				w.buf.WriteString(col.String())
			}
			w.buf.WriteByte(')')
		}
	}

	// WITH AGGREGATION POLICY
	if n.AggregationPolicy != nil {
		w.buf.WriteString(" WITH AGGREGATION POLICY ")
		w.writeObjectNameNoSpace(n.AggregationPolicy)
	}

	// WITH TAG
	if len(n.Tags) > 0 {
		w.buf.WriteString(" WITH TAG (")
		w.writeTagAssignments(n.Tags)
		w.buf.WriteByte(')')
	}

	// WITH CONTACT
	if len(n.Contacts) > 0 {
		w.buf.WriteString(" WITH CONTACT (")
		for i, c := range n.Contacts {
			if i > 0 {
				w.buf.WriteString(", ")
			}
			w.buf.WriteString(c.Name)
			w.buf.WriteString(" = ")
			w.buf.WriteString(c.Words)
		}
		w.buf.WriteByte(')')
	}

	// AS SELECT, after the clauses the CTAS form writes before it
	if n.AsSelect != nil {
		w.buf.WriteString(" AS ")
		return w.writeQueryExpr(n.AsSelect)
	}
	return nil
}

// writeBackupSource emits FROM BACKUP SET <set> IDENTIFIER '<id>'.
func (w *writer) writeBackupSource(b *ast.BackupSource) {
	w.buf.WriteString(" FROM BACKUP SET ")
	w.writeObjectNameNoSpace(b.BackupSet)
	w.buf.WriteString(" IDENTIFIER ")
	w.buf.WriteString(quoteString(b.ID))
}

func (w *writer) writeColumnDef(n *ast.ColumnDef) error {
	w.buf.WriteString(n.Name.String())
	if n.DataType != nil {
		w.buf.WriteByte(' ')
		w.writeTypeName(n.DataType)
	}

	if n.VirtualExpr != nil {
		w.buf.WriteString(" AS (")
		if err := writeExprNoLeadSpace(w, n.VirtualExpr); err != nil {
			return err
		}
		w.buf.WriteByte(')')
	}

	if n.Collate != "" {
		w.buf.WriteString(" COLLATE '")
		w.buf.WriteString(n.Collate)
		w.buf.WriteByte('\'')
	}

	if n.NotNull {
		w.buf.WriteString(" NOT NULL")
	} else if n.Nullable {
		w.buf.WriteString(" NULL")
	}

	if n.Default != nil {
		w.buf.WriteString(" DEFAULT ")
		if err := writeExprNoLeadSpace(w, n.Default); err != nil {
			return err
		}
	}

	if n.Identity != nil {
		w.writeIdentitySpec(n.Identity)
	}

	if n.InlineConstraint != nil {
		if err := w.writeInlineConstraint(n.InlineConstraint); err != nil {
			return err
		}
	}

	if n.MaskingPolicy != nil {
		w.buf.WriteString(" WITH MASKING POLICY ")
		w.writeObjectNameNoSpace(n.MaskingPolicy)
	}

	if n.ProjectionPolicy != nil {
		w.buf.WriteString(" WITH PROJECTION POLICY ")
		w.writeObjectNameNoSpace(n.ProjectionPolicy)
	}

	if n.Comment != nil {
		w.buf.WriteString(" COMMENT ")
		w.buf.WriteString(quoteString(*n.Comment))
	}

	if len(n.Tags) > 0 {
		w.buf.WriteString(" WITH TAG (")
		w.writeTagAssignments(n.Tags)
		w.buf.WriteByte(')')
	}

	return nil
}

func (w *writer) writeIdentitySpec(id *ast.IdentitySpec) {
	w.buf.WriteString(" IDENTITY")
	if id.Start != nil || id.Increment != nil {
		w.buf.WriteByte('(')
		if id.Start != nil {
			w.buf.WriteString(strconv.FormatInt(*id.Start, 10))
		} else {
			w.buf.WriteByte('1')
		}
		if id.Increment != nil {
			w.buf.WriteString(", ")
			w.buf.WriteString(strconv.FormatInt(*id.Increment, 10))
		}
		w.buf.WriteByte(')')
	}
	if id.Order != nil {
		if *id.Order {
			w.buf.WriteString(" ORDER")
		} else {
			w.buf.WriteString(" NOORDER")
		}
	}
}

func (w *writer) writeInlineConstraint(c *ast.InlineConstraint) error {
	if !c.Name.IsEmpty() {
		w.buf.WriteString(" CONSTRAINT ")
		w.buf.WriteString(c.Name.String())
	}
	switch c.Type {
	case ast.ConstrPrimaryKey:
		w.buf.WriteString(" PRIMARY KEY")
	case ast.ConstrUnique:
		w.buf.WriteString(" UNIQUE")
	case ast.ConstrForeignKey:
		w.buf.WriteString(" FOREIGN KEY REFERENCES ")
		if c.References != nil {
			w.writeObjectNameNoSpace(c.References.Table)
			if len(c.References.Columns) > 0 {
				w.buf.WriteString(" (")
				for i, col := range c.References.Columns {
					if i > 0 {
						w.buf.WriteString(", ")
					}
					w.buf.WriteString(col.String())
				}
				w.buf.WriteByte(')')
			}
			w.writeRefActions(c.References.OnDelete, c.References.OnUpdate)
		}
	case ast.ConstrCheck:
		w.buf.WriteString(" CHECK (")
		if err := writeExprNoLeadSpace(w, c.Check); err != nil {
			return err
		}
		w.buf.WriteByte(')')
	}
	return nil
}

func (w *writer) writeTableConstraint(n *ast.TableConstraint) error {
	if !n.Name.IsEmpty() {
		w.buf.WriteString("CONSTRAINT ")
		w.buf.WriteString(n.Name.String())
		w.buf.WriteByte(' ')
	}
	switch n.Type {
	case ast.ConstrPrimaryKey:
		w.buf.WriteString("PRIMARY KEY")
	case ast.ConstrUnique:
		w.buf.WriteString("UNIQUE")
	case ast.ConstrForeignKey:
		w.buf.WriteString("FOREIGN KEY")
	case ast.ConstrCheck:
		w.buf.WriteString("CHECK (")
		if err := writeExprNoLeadSpace(w, n.Check); err != nil {
			return err
		}
		w.buf.WriteByte(')')
	}
	if len(n.Columns) > 0 {
		w.buf.WriteString(" (")
		for i, col := range n.Columns {
			if i > 0 {
				w.buf.WriteString(", ")
			}
			w.buf.WriteString(col.String())
		}
		w.buf.WriteByte(')')
	}
	if n.Type == ast.ConstrForeignKey && n.References != nil {
		w.buf.WriteString(" REFERENCES ")
		w.writeObjectNameNoSpace(n.References.Table)
		if len(n.References.Columns) > 0 {
			w.buf.WriteString(" (")
			for i, col := range n.References.Columns {
				if i > 0 {
					w.buf.WriteString(", ")
				}
				w.buf.WriteString(col.String())
			}
			w.buf.WriteByte(')')
		}
		w.writeRefActions(n.References.OnDelete, n.References.OnUpdate)
	}
	if n.Comment != nil {
		w.buf.WriteString(" COMMENT ")
		w.buf.WriteString(quoteString(*n.Comment))
	}
	return nil
}

// writeTableIndex emits an inline INDEX element:
// INDEX name (col, col, ...) [INCLUDE (col, ...)].
func (w *writer) writeTableIndex(n *ast.TableIndex) {
	w.buf.WriteString("INDEX ")
	w.buf.WriteString(n.Name.String())
	w.buf.WriteString(" (")
	for i, col := range n.Columns {
		if i > 0 {
			w.buf.WriteString(", ")
		}
		w.buf.WriteString(col.String())
	}
	w.buf.WriteByte(')')
	if len(n.Include) > 0 {
		w.buf.WriteString(" INCLUDE (")
		for i, col := range n.Include {
			if i > 0 {
				w.buf.WriteString(", ")
			}
			w.buf.WriteString(col.String())
		}
		w.buf.WriteByte(')')
	}
}

func (w *writer) writeRefActions(onDelete, onUpdate ast.ReferenceAction) {
	if onDelete != ast.RefActNone {
		w.buf.WriteString(" ON DELETE ")
		w.buf.WriteString(refActionString(onDelete))
	}
	if onUpdate != ast.RefActNone {
		w.buf.WriteString(" ON UPDATE ")
		w.buf.WriteString(refActionString(onUpdate))
	}
}

func refActionString(a ast.ReferenceAction) string {
	switch a {
	case ast.RefActCascade:
		return "CASCADE"
	case ast.RefActSetNull:
		return "SET NULL"
	case ast.RefActSetDefault:
		return "SET DEFAULT"
	case ast.RefActRestrict:
		return "RESTRICT"
	case ast.RefActNoAction:
		return "NO ACTION"
	default:
		return ""
	}
}

// writeClusterBy writes CLUSTER BY [LINEAR] (exprs) if clusterBy is non-empty.
func (w *writer) writeClusterBy(clusterBy []ast.Node, linear bool) error {
	if len(clusterBy) == 0 {
		return nil
	}
	w.buf.WriteString(" CLUSTER BY")
	if linear {
		w.buf.WriteString(" LINEAR")
	}
	w.buf.WriteString(" (")
	for i, expr := range clusterBy {
		if i > 0 {
			w.buf.WriteString(", ")
		}
		if err := writeExprNoLeadSpace(w, expr); err != nil {
			return err
		}
	}
	w.buf.WriteByte(')')
	return nil
}

// writeTagAssignments writes TAG name = 'value' pairs (without the outer parens).
func (w *writer) writeTagAssignments(tags []*ast.TagAssignment) {
	for i, tag := range tags {
		if i > 0 {
			w.buf.WriteString(", ")
		}
		w.writeObjectNameNoSpace(tag.Name)
		w.buf.WriteString(" = ")
		w.buf.WriteString(quoteString(tag.Value))
	}
}

// writeCloneTimeTravelSuffix writes the AT/BEFORE time travel qualifier
// `{ AT | BEFORE } ( <kind> => <value> )`. The value is a general expression
// (a string literal, a TO_TIMESTAMP_TZ(...) call, an offset arithmetic, …).
func (w *writer) writeCloneTimeTravelSuffix(c *ast.CloneSource) error {
	if c.AtBefore == "" {
		return nil
	}
	if c.AtBefore == "BEFORE" {
		w.buf.WriteString(" BEFORE (")
	} else {
		w.buf.WriteString(" AT (")
	}
	w.buf.WriteString(c.Kind)
	w.buf.WriteString(" => ")
	if err := w.writeExpr(c.Value); err != nil {
		return err
	}
	w.buf.WriteByte(')')
	return nil
}

// ---------------------------------------------------------------------------
// DATABASE DDL
// ---------------------------------------------------------------------------

func (w *writer) writeCreateDatabaseStmt(n *ast.CreateDatabaseStmt) error {
	w.buf.WriteString("CREATE")
	if n.OrReplace {
		w.buf.WriteString(" OR REPLACE")
	} else if n.OrAlter {
		w.buf.WriteString(" OR ALTER")
	}
	if n.Transient {
		w.buf.WriteString(" TRANSIENT")
	}
	w.buf.WriteString(" DATABASE")
	if n.IfNotExists {
		w.buf.WriteString(" IF NOT EXISTS")
	}
	w.buf.WriteByte(' ')
	w.writeObjectNameNoSpace(n.Name)
	if n.Clone != nil {
		w.buf.WriteString(" CLONE ")
		w.writeObjectNameNoSpace(n.Clone.Source)
		if err := w.writeCloneTimeTravelSuffix(n.Clone); err != nil {
			return err
		}
		w.writeCloneIgnoreClauses(n.IgnoreInsufficientRetention, n.IgnoreHybridTables)
	}
	if n.FromShare != nil {
		w.buf.WriteString(" FROM SHARE ")
		w.writeObjectNameNoSpace(n.FromShare)
	}
	if n.FromListing != nil {
		w.buf.WriteString(" FROM LISTING ")
		w.buf.WriteString(quoteString(*n.FromListing))
	}
	if n.FromBackupSet != nil {
		w.writeFromBackupSet(n.FromBackupSet, n.BackupID)
	}
	if n.ReplicaOf != nil {
		w.buf.WriteString(" AS REPLICA OF ")
		w.writeObjectNameNoSpace(n.ReplicaOf)
	}
	return w.writeCreateDBSchemaTail(&n.Props, n.Tags)
}

func (w *writer) writeAlterDatabaseStmt(n *ast.AlterDatabaseStmt) error {
	w.buf.WriteString("ALTER DATABASE")
	if n.IfExists {
		w.buf.WriteString(" IF EXISTS")
	}
	w.buf.WriteByte(' ')
	w.writeObjectNameNoSpace(n.Name)
	switch n.Action {
	case ast.AlterDBRename:
		w.buf.WriteString(" RENAME TO ")
		w.writeObjectNameNoSpace(n.NewName)
	case ast.AlterDBSwap:
		w.buf.WriteString(" SWAP WITH ")
		w.writeObjectNameNoSpace(n.NewName)
	case ast.AlterDBSet:
		w.buf.WriteString(" SET")
		if err := w.writeDBSchemaSetProps(n.SetProps); err != nil {
			return err
		}
	case ast.AlterDBUnset:
		w.buf.WriteString(" UNSET ")
		w.buf.WriteString(strings.Join(n.UnsetProps, ", "))
	case ast.AlterDBSetTag:
		w.buf.WriteString(" SET TAG ")
		w.writeTagAssignments(n.Tags)
	case ast.AlterDBUnsetTag:
		w.buf.WriteString(" UNSET TAG ")
		w.writeObjectNameList(n.UnsetTags)
	case ast.AlterDBEnableReplication:
		w.buf.WriteString(" ENABLE REPLICATION TO ACCOUNTS")
	case ast.AlterDBDisableReplication:
		w.buf.WriteString(" DISABLE REPLICATION TO ACCOUNTS")
	case ast.AlterDBEnableFailover:
		w.buf.WriteString(" ENABLE FAILOVER TO ACCOUNTS")
	case ast.AlterDBDisableFailover:
		w.buf.WriteString(" DISABLE FAILOVER TO ACCOUNTS")
	case ast.AlterDBRefresh:
		w.buf.WriteString(" REFRESH")
	case ast.AlterDBPrimary:
		w.buf.WriteString(" PRIMARY")
	case ast.AlterDBSetFeaturePolicy:
		w.buf.WriteString(" SET FEATURE POLICY ")
		w.writeObjectNameNoSpace(n.FeaturePolicy)
		if n.Force {
			w.buf.WriteString(" FORCE")
		}
	case ast.AlterDBUnsetFeaturePolicy:
		w.buf.WriteString(" UNSET FEATURE POLICY")
	case ast.AlterDBSuspendDiscovery:
		w.buf.WriteString(" SUSPEND DISCOVERY")
	case ast.AlterDBResumeDiscovery:
		w.buf.WriteString(" RESUME DISCOVERY")
	case ast.AlterDBUpdateLinkedCatalog:
		if err := w.writeLinkedCatalogUpdate(n.LinkedCatalog); err != nil {
			return err
		}
	default:
		return fmt.Errorf("deparse: unsupported ALTER DATABASE action %d", n.Action)
	}
	return nil
}

func (w *writer) writeDropDatabaseStmt(n *ast.DropDatabaseStmt) error {
	w.buf.WriteString("DROP DATABASE")
	if n.IfExists {
		w.buf.WriteString(" IF EXISTS")
	}
	w.buf.WriteByte(' ')
	w.writeObjectNameNoSpace(n.Name)
	if n.Cascade {
		w.buf.WriteString(" CASCADE")
	} else if n.Restrict {
		w.buf.WriteString(" RESTRICT")
	}
	return nil
}

func (w *writer) writeUndropDatabaseStmt(n *ast.UndropDatabaseStmt) error {
	w.buf.WriteString("UNDROP DATABASE ")
	w.writeObjectNameNoSpace(n.Name)
	return nil
}

// ---------------------------------------------------------------------------
// SCHEMA DDL
// ---------------------------------------------------------------------------

func (w *writer) writeCreateSchemaStmt(n *ast.CreateSchemaStmt) error {
	w.buf.WriteString("CREATE")
	if n.OrReplace {
		w.buf.WriteString(" OR REPLACE")
	} else if n.OrAlter {
		w.buf.WriteString(" OR ALTER")
	}
	if n.Transient {
		w.buf.WriteString(" TRANSIENT")
	}
	if n.Versioned {
		w.buf.WriteString(" VERSIONED")
	}
	w.buf.WriteString(" SCHEMA")
	if n.IfNotExists {
		w.buf.WriteString(" IF NOT EXISTS")
	}
	w.buf.WriteByte(' ')
	w.writeObjectNameNoSpace(n.Name)
	if n.Clone != nil {
		w.buf.WriteString(" CLONE ")
		w.writeObjectNameNoSpace(n.Clone.Source)
		if err := w.writeCloneTimeTravelSuffix(n.Clone); err != nil {
			return err
		}
		w.writeCloneIgnoreClauses(n.IgnoreInsufficientRetention, n.IgnoreHybridTables)
	}
	if n.FromBackupSet != nil {
		w.writeFromBackupSet(n.FromBackupSet, n.BackupID)
	}
	if n.ManagedAccess {
		w.buf.WriteString(" WITH MANAGED ACCESS")
	}
	return w.writeCreateDBSchemaTail(&n.Props, n.Tags)
}

func (w *writer) writeAlterSchemaStmt(n *ast.AlterSchemaStmt) error {
	w.buf.WriteString("ALTER SCHEMA")
	if n.IfExists {
		w.buf.WriteString(" IF EXISTS")
	}
	w.buf.WriteByte(' ')
	w.writeObjectNameNoSpace(n.Name)
	switch n.Action {
	case ast.AlterSchemaRename:
		w.buf.WriteString(" RENAME TO ")
		w.writeObjectNameNoSpace(n.NewName)
	case ast.AlterSchemaSwap:
		w.buf.WriteString(" SWAP WITH ")
		w.writeObjectNameNoSpace(n.NewName)
	case ast.AlterSchemaSet:
		w.buf.WriteString(" SET")
		if err := w.writeDBSchemaSetProps(n.SetProps); err != nil {
			return err
		}
	case ast.AlterSchemaUnset:
		w.buf.WriteString(" UNSET ")
		w.buf.WriteString(strings.Join(n.UnsetProps, ", "))
	case ast.AlterSchemaSetTag:
		w.buf.WriteString(" SET TAG ")
		w.writeTagAssignments(n.Tags)
	case ast.AlterSchemaUnsetTag:
		w.buf.WriteString(" UNSET TAG ")
		w.writeObjectNameList(n.UnsetTags)
	case ast.AlterSchemaEnableManagedAccess:
		w.buf.WriteString(" ENABLE MANAGED ACCESS")
	case ast.AlterSchemaDisableManagedAccess:
		w.buf.WriteString(" DISABLE MANAGED ACCESS")
	default:
		return fmt.Errorf("deparse: unsupported ALTER SCHEMA action %d", n.Action)
	}
	return nil
}

func (w *writer) writeDropSchemaStmt(n *ast.DropSchemaStmt) error {
	w.buf.WriteString("DROP SCHEMA")
	if n.IfExists {
		w.buf.WriteString(" IF EXISTS")
	}
	w.buf.WriteByte(' ')
	w.writeObjectNameNoSpace(n.Name)
	if n.Cascade {
		w.buf.WriteString(" CASCADE")
	} else if n.Restrict {
		w.buf.WriteString(" RESTRICT")
	}
	return nil
}

func (w *writer) writeUndropSchemaStmt(n *ast.UndropSchemaStmt) error {
	w.buf.WriteString("UNDROP SCHEMA ")
	w.writeObjectNameNoSpace(n.Name)
	return nil
}

// ---------------------------------------------------------------------------
// DROP / UNDROP (non-database/schema)
// ---------------------------------------------------------------------------

func (w *writer) writeDropStmt(n *ast.DropStmt) error {
	w.buf.WriteString("DROP ")
	if n.Iceberg {
		w.buf.WriteString("ICEBERG ")
	}
	w.buf.WriteString(n.ObjectType())
	if n.IfExists {
		w.buf.WriteString(" IF EXISTS")
	}
	w.buf.WriteByte(' ')
	w.writeObjectNameNoSpace(n.Name)
	if n.HasArgs {
		// DROP FUNCTION / PROCEDURE overload signature: f(NUMBER, VARCHAR).
		w.buf.WriteByte('(')
		for i, t := range n.ArgTypes {
			if i > 0 {
				w.buf.WriteString(", ")
			}
			w.writeTypeName(t)
		}
		w.buf.WriteByte(')')
	}
	if n.Cascade {
		w.buf.WriteString(" CASCADE")
	} else if n.Restrict {
		w.buf.WriteString(" RESTRICT")
	}
	if n.Purge {
		w.buf.WriteString(" PURGE")
	} else if n.NoPurge {
		w.buf.WriteString(" NOPURGE")
	}
	return nil
}

func (w *writer) writeUndropStmt(n *ast.UndropStmt) error {
	w.buf.WriteString("UNDROP ")
	if n.Iceberg {
		w.buf.WriteString("ICEBERG ")
	}
	w.buf.WriteString(n.Kind.String())
	w.buf.WriteByte(' ')
	w.writeObjectNameNoSpace(n.Name)
	return nil
}

// ---------------------------------------------------------------------------
// CREATE VIEW / ALTER VIEW
// ---------------------------------------------------------------------------

func (w *writer) writeCreateViewStmt(n *ast.CreateViewStmt) error {
	w.buf.WriteString("CREATE")
	if n.OrReplace {
		w.buf.WriteString(" OR REPLACE")
	} else if n.OrAlter {
		w.buf.WriteString(" OR ALTER")
	}
	if n.Secure {
		w.buf.WriteString(" SECURE")
	}
	if n.Recursive {
		w.buf.WriteString(" RECURSIVE")
	}
	w.buf.WriteString(" VIEW")
	if n.IfNotExists {
		w.buf.WriteString(" IF NOT EXISTS")
	}
	w.buf.WriteByte(' ')
	w.writeObjectNameNoSpace(n.Name)
	// Column list
	if len(n.Columns) > 0 {
		w.buf.WriteString(" (")
		for i, col := range n.Columns {
			if i > 0 {
				w.buf.WriteString(", ")
			}
			if err := w.writeViewColumn(col); err != nil {
				return err
			}
		}
		w.buf.WriteByte(')')
	}
	// ViewCols (masking/tag bindings outside parens)
	for _, vc := range n.ViewCols {
		if err := w.writeViewColumn(vc); err != nil {
			return err
		}
	}
	if n.CopyGrants {
		w.buf.WriteString(" COPY GRANTS")
	}
	if n.Comment != nil {
		w.buf.WriteString(" COMMENT = ")
		w.buf.WriteString(quoteString(*n.Comment))
	}
	if len(n.Tags) > 0 {
		w.buf.WriteString(" WITH TAG (")
		w.writeTagAssignments(n.Tags)
		w.buf.WriteByte(')')
	}
	if n.RowPolicy != nil {
		w.buf.WriteString(" WITH ROW ACCESS POLICY ")
		w.writeObjectNameNoSpace(n.RowPolicy.PolicyName)
		if len(n.RowPolicy.Columns) > 0 {
			w.buf.WriteString(" ON (")
			for i, col := range n.RowPolicy.Columns {
				if i > 0 {
					w.buf.WriteString(", ")
				}
				w.buf.WriteString(col.String())
			}
			w.buf.WriteByte(')')
		}
	}
	w.buf.WriteString(" AS ")
	return w.writeQueryExpr(n.Query)
}

func (w *writer) writeViewColumn(vc *ast.ViewColumn) error {
	w.buf.WriteString(vc.Name.String())
	if vc.MaskingPolicy != nil {
		w.buf.WriteString(" WITH MASKING POLICY ")
		w.writeObjectNameNoSpace(vc.MaskingPolicy)
		if len(vc.MaskingUsing) > 0 {
			w.buf.WriteString(" USING (")
			for i, col := range vc.MaskingUsing {
				if i > 0 {
					w.buf.WriteString(", ")
				}
				w.buf.WriteString(col.String())
			}
			w.buf.WriteByte(')')
		}
	}
	if vc.ProjectionPolicy != nil {
		w.buf.WriteString(" WITH PROJECTION POLICY ")
		w.writeObjectNameNoSpace(vc.ProjectionPolicy)
	}
	if len(vc.Tags) > 0 {
		w.buf.WriteString(" WITH TAG (")
		w.writeTagAssignments(vc.Tags)
		w.buf.WriteByte(')')
	}
	if vc.Comment != nil {
		w.buf.WriteString(" COMMENT ")
		w.buf.WriteString(quoteString(*vc.Comment))
	}
	return nil
}

func (w *writer) writeAlterViewStmt(n *ast.AlterViewStmt) error {
	w.buf.WriteString("ALTER VIEW")
	if n.IfExists {
		w.buf.WriteString(" IF EXISTS")
	}
	w.buf.WriteByte(' ')
	w.writeObjectNameNoSpace(n.Name)
	// Column action list: MODIFY COLUMN a SET MASKING POLICY p, COLUMN b
	// UNSET MASKING POLICY. (A statement built without ColumnActions keeps
	// the legacy single-action rendering below, driven by the mirrored fields.)
	if len(n.ColumnActions) > 0 {
		w.buf.WriteString(" MODIFY ")
		for i, a := range n.ColumnActions {
			if i > 0 {
				w.buf.WriteString(", ")
			}
			w.writeAlterViewColumnAction(a)
		}
		return nil
	}
	switch n.Action {
	case ast.AlterViewRename:
		w.buf.WriteString(" RENAME TO ")
		w.writeObjectNameNoSpace(n.NewName)
	case ast.AlterViewSetComment:
		w.buf.WriteString(" SET COMMENT = ")
		if n.Comment != nil {
			w.buf.WriteString(quoteString(*n.Comment))
		} else {
			w.buf.WriteString("NULL")
		}
	case ast.AlterViewUnsetComment:
		w.buf.WriteString(" UNSET COMMENT")
	case ast.AlterViewSetSecure:
		w.buf.WriteString(" SET SECURE")
	case ast.AlterViewUnsetSecure:
		w.buf.WriteString(" UNSET SECURE")
	case ast.AlterViewSetTag:
		w.buf.WriteString(" SET TAG (")
		w.writeTagAssignments(n.Tags)
		w.buf.WriteByte(')')
	case ast.AlterViewUnsetTag:
		w.buf.WriteString(" UNSET TAG (")
		w.writeObjectNameList(n.UnsetTags)
		w.buf.WriteByte(')')
	case ast.AlterViewAddRowAccessPolicy:
		w.buf.WriteString(" ADD ROW ACCESS POLICY ")
		w.writeObjectNameNoSpace(n.PolicyName)
		if len(n.PolicyCols) > 0 {
			w.buf.WriteString(" ON (")
			for i, col := range n.PolicyCols {
				if i > 0 {
					w.buf.WriteString(", ")
				}
				w.buf.WriteString(col.String())
			}
			w.buf.WriteByte(')')
		}
	case ast.AlterViewDropRowAccessPolicy:
		w.buf.WriteString(" DROP ROW ACCESS POLICY ")
		w.writeObjectNameNoSpace(n.PolicyName)
		if n.AddPolicyName != nil {
			// Drop-and-add combination form.
			w.buf.WriteString(", ADD ROW ACCESS POLICY ")
			w.writeObjectNameNoSpace(n.AddPolicyName)
			w.buf.WriteString(" ON (")
			for i, col := range n.AddPolicyCols {
				if i > 0 {
					w.buf.WriteString(", ")
				}
				w.buf.WriteString(col.String())
			}
			w.buf.WriteByte(')')
		}
	case ast.AlterViewDropAllRowAccessPolicies:
		w.buf.WriteString(" DROP ALL ROW ACCESS POLICIES")
	case ast.AlterViewColumnSetMaskingPolicy:
		w.buf.WriteString(" ALTER COLUMN ")
		w.buf.WriteString(n.Column.String())
		w.buf.WriteString(" SET MASKING POLICY ")
		w.writeObjectNameNoSpace(n.MaskingPolicy)
		if len(n.MaskingUsing) > 0 {
			w.buf.WriteString(" USING (")
			for i, col := range n.MaskingUsing {
				if i > 0 {
					w.buf.WriteString(", ")
				}
				w.buf.WriteString(col.String())
			}
			w.buf.WriteByte(')')
		}
	case ast.AlterViewColumnUnsetMaskingPolicy:
		w.buf.WriteString(" ALTER COLUMN ")
		w.buf.WriteString(n.Column.String())
		w.buf.WriteString(" UNSET MASKING POLICY")
	case ast.AlterViewColumnSetTag:
		w.buf.WriteString(" ALTER COLUMN ")
		w.buf.WriteString(n.Column.String())
		w.buf.WriteString(" SET TAG (")
		w.writeTagAssignments(n.Tags)
		w.buf.WriteByte(')')
	case ast.AlterViewColumnUnsetTag:
		w.buf.WriteString(" ALTER COLUMN ")
		w.buf.WriteString(n.Column.String())
		w.buf.WriteString(" UNSET TAG (")
		w.writeObjectNameList(n.UnsetTags)
		w.buf.WriteByte(')')
	case ast.AlterViewSetAggregationPolicy:
		w.buf.WriteString(" SET AGGREGATION POLICY ")
		w.writeObjectNameNoSpace(n.PolicyName)
		if len(n.PolicyKeyCols) > 0 {
			w.buf.WriteString(" ENTITY KEY (")
			for i, col := range n.PolicyKeyCols {
				if i > 0 {
					w.buf.WriteString(", ")
				}
				w.buf.WriteString(col.String())
			}
			w.buf.WriteByte(')')
		}
		if n.PolicyForce {
			w.buf.WriteString(" FORCE")
		}
	case ast.AlterViewUnsetAggregationPolicy:
		w.buf.WriteString(" UNSET AGGREGATION POLICY")
	case ast.AlterViewSetJoinPolicy:
		w.buf.WriteString(" SET JOIN POLICY ")
		w.writeObjectNameNoSpace(n.PolicyName)
		if n.PolicyForce {
			w.buf.WriteString(" FORCE")
		}
		if len(n.PolicyKeyCols) > 0 {
			w.buf.WriteByte(' ')
			w.buf.WriteString(n.PolicyKeyKind)
			w.buf.WriteString(" JOIN KEYS (")
			for i, col := range n.PolicyKeyCols {
				if i > 0 {
					w.buf.WriteString(", ")
				}
				w.buf.WriteString(col.String())
			}
			w.buf.WriteByte(')')
		}
	case ast.AlterViewUnsetJoinPolicy:
		w.buf.WriteString(" UNSET JOIN POLICY")
	case ast.AlterViewSetChangeTracking:
		w.buf.WriteString(" SET CHANGE_TRACKING = ")
		if n.ChangeTracking {
			w.buf.WriteString("TRUE")
		} else {
			w.buf.WriteString("FALSE")
		}
	case ast.AlterViewSetContact:
		w.buf.WriteString(" SET CONTACT ")
		w.buf.WriteString(n.Contact)
	case ast.AlterViewUnsetContact:
		w.buf.WriteString(" UNSET CONTACT ")
		w.buf.WriteString(n.Contact)
	case ast.AlterViewDataMetric:
		w.buf.WriteByte(' ')
		return w.writeDataMetricAction(n.DataMetric)
	case ast.AlterViewAddColumn:
		w.buf.WriteString(" ADD COLUMN")
		if n.IfNotExists {
			w.buf.WriteString(" IF NOT EXISTS")
		}
		w.buf.WriteByte(' ')
		w.buf.WriteString(n.Column.String())
		w.buf.WriteByte(' ')
		w.writeTypeName(n.ColumnType)
		if n.MaskingPolicy != nil {
			w.buf.WriteString(" WITH MASKING POLICY ")
			w.writeObjectNameNoSpace(n.MaskingPolicy)
			if len(n.MaskingUsing) > 0 {
				w.buf.WriteString(" USING")
				w.writeParenIdentList(n.MaskingUsing)
			}
		}
		if n.ProjectionPolicy != nil {
			w.buf.WriteString(" WITH PROJECTION POLICY ")
			w.writeObjectNameNoSpace(n.ProjectionPolicy)
		}
		if len(n.Tags) > 0 {
			w.buf.WriteString(" WITH TAG (")
			w.writeTagAssignments(n.Tags)
			w.buf.WriteByte(')')
		}
	default:
		return fmt.Errorf("deparse: unsupported ALTER VIEW action %d", n.Action)
	}
	return nil
}

// writeAlterViewColumnAction emits one element of a multi-element ALTER VIEW
// MODIFY column action list: COLUMN <col> SET MASKING POLICY <p> [USING (...)]
// [FORCE] | UNSET MASKING POLICY | SET TAG (...) | UNSET TAG (...).
func (w *writer) writeAlterViewColumnAction(a *ast.AlterViewColumnAction) {
	w.buf.WriteString("COLUMN ")
	w.buf.WriteString(a.Column.String())
	switch a.Action {
	case ast.AlterViewColumnSetMaskingPolicy:
		w.buf.WriteString(" SET MASKING POLICY ")
		w.writeObjectNameNoSpace(a.MaskingPolicy)
		if len(a.MaskingUsing) > 0 {
			w.buf.WriteString(" USING (")
			for i, col := range a.MaskingUsing {
				if i > 0 {
					w.buf.WriteString(", ")
				}
				w.buf.WriteString(col.String())
			}
			w.buf.WriteByte(')')
		}
		if a.Force {
			w.buf.WriteString(" FORCE")
		}
	case ast.AlterViewColumnUnsetMaskingPolicy:
		w.buf.WriteString(" UNSET MASKING POLICY")
	case ast.AlterViewColumnSetTag:
		w.buf.WriteString(" SET TAG (")
		w.writeTagAssignments(a.Tags)
		w.buf.WriteByte(')')
	case ast.AlterViewColumnUnsetTag:
		w.buf.WriteString(" UNSET TAG (")
		w.writeObjectNameList(a.UnsetTags)
		w.buf.WriteByte(')')
	case ast.AlterViewColumnSetProjectionPolicy:
		w.buf.WriteString(" SET PROJECTION POLICY ")
		w.writeObjectNameNoSpace(a.ProjectionPolicy)
		if a.Force {
			w.buf.WriteString(" FORCE")
		}
	case ast.AlterViewColumnUnsetProjectionPolicy:
		w.buf.WriteString(" UNSET PROJECTION POLICY")
	}
}

// ---------------------------------------------------------------------------
// CREATE MATERIALIZED VIEW / ALTER MATERIALIZED VIEW
// ---------------------------------------------------------------------------

func (w *writer) writeCreateMaterializedViewStmt(n *ast.CreateMaterializedViewStmt) error {
	w.buf.WriteString("CREATE")
	if n.OrReplace {
		w.buf.WriteString(" OR REPLACE")
	} else if n.OrAlter {
		w.buf.WriteString(" OR ALTER")
	}
	if n.Secure {
		w.buf.WriteString(" SECURE")
	}
	if n.Interactive {
		w.buf.WriteString(" INTERACTIVE")
	}
	w.buf.WriteString(" MATERIALIZED VIEW")
	if n.IfNotExists {
		w.buf.WriteString(" IF NOT EXISTS")
	}
	w.buf.WriteByte(' ')
	w.writeObjectNameNoSpace(n.Name)
	if len(n.Columns) > 0 {
		w.buf.WriteString(" (")
		for i, col := range n.Columns {
			if i > 0 {
				w.buf.WriteString(", ")
			}
			if err := w.writeViewColumn(col); err != nil {
				return err
			}
		}
		w.buf.WriteByte(')')
	}
	for _, vc := range n.ViewCols {
		if err := w.writeViewColumn(vc); err != nil {
			return err
		}
	}
	if n.CopyGrants {
		w.buf.WriteString(" COPY GRANTS")
	}
	if n.Comment != nil {
		w.buf.WriteString(" COMMENT = ")
		w.buf.WriteString(quoteString(*n.Comment))
	}
	if err := w.writeClusterBy(n.ClusterBy, n.Linear); err != nil {
		return err
	}
	if len(n.Tags) > 0 {
		w.buf.WriteString(" WITH TAG (")
		w.writeTagAssignments(n.Tags)
		w.buf.WriteByte(')')
	}
	if n.RowPolicy != nil {
		w.buf.WriteString(" WITH ROW ACCESS POLICY ")
		w.writeObjectNameNoSpace(n.RowPolicy.PolicyName)
		if len(n.RowPolicy.Columns) > 0 {
			w.buf.WriteString(" ON (")
			for i, col := range n.RowPolicy.Columns {
				if i > 0 {
					w.buf.WriteString(", ")
				}
				w.buf.WriteString(col.String())
			}
			w.buf.WriteByte(')')
		}
	}
	w.buf.WriteString(" AS ")
	return w.writeQueryExpr(n.Query)
}

func (w *writer) writeAlterMaterializedViewStmt(n *ast.AlterMaterializedViewStmt) error {
	w.buf.WriteString("ALTER MATERIALIZED VIEW ")
	w.writeObjectNameNoSpace(n.Name)
	switch n.Action {
	case ast.AlterMVRename:
		w.buf.WriteString(" RENAME TO ")
		w.writeObjectNameNoSpace(n.NewName)
	case ast.AlterMVClusterBy:
		if err := w.writeClusterBy(n.ClusterBy, n.Linear); err != nil {
			return err
		}
	case ast.AlterMVDropClusteringKey:
		w.buf.WriteString(" DROP CLUSTERING KEY")
	case ast.AlterMVSuspend:
		w.buf.WriteString(" SUSPEND")
	case ast.AlterMVResume:
		w.buf.WriteString(" RESUME")
	case ast.AlterMVSuspendRecluster:
		w.buf.WriteString(" SUSPEND RECLUSTER")
	case ast.AlterMVResumeRecluster:
		w.buf.WriteString(" RESUME RECLUSTER")
	case ast.AlterMVSetSecure, ast.AlterMVSetContact:
		w.buf.WriteString(" SET")
		if n.Action == ast.AlterMVSetSecure {
			w.buf.WriteString(" SECURE")
		}
		if n.Contact != "" {
			w.buf.WriteString(" CONTACT ")
			w.buf.WriteString(n.Contact)
		}
		if n.Comment != nil {
			w.buf.WriteString(" COMMENT = ")
			w.buf.WriteString(quoteString(*n.Comment))
		}
	case ast.AlterMVUnsetContact:
		w.buf.WriteString(" UNSET CONTACT ")
		w.buf.WriteString(n.Contact)
	case ast.AlterMVDataMetric:
		w.buf.WriteByte(' ')
		return w.writeDataMetricAction(n.DataMetric)
	case ast.AlterMVUnsetSecure:
		w.buf.WriteString(" UNSET SECURE")
	case ast.AlterMVSetComment:
		w.buf.WriteString(" SET COMMENT = ")
		if n.Comment != nil {
			w.buf.WriteString(quoteString(*n.Comment))
		} else {
			w.buf.WriteString("NULL")
		}
	case ast.AlterMVUnsetComment:
		w.buf.WriteString(" UNSET COMMENT")
	default:
		return fmt.Errorf("deparse: unsupported ALTER MATERIALIZED VIEW action %d", n.Action)
	}
	return nil
}

// ---------------------------------------------------------------------------
// ALTER TABLE
// ---------------------------------------------------------------------------

func (w *writer) writeAlterTableStmt(n *ast.AlterTableStmt) error {
	if n.Iceberg {
		w.buf.WriteString("ALTER ICEBERG TABLE")
	} else {
		w.buf.WriteString("ALTER TABLE")
	}
	if n.IfExists {
		w.buf.WriteString(" IF EXISTS")
	}
	w.buf.WriteByte(' ')
	w.writeObjectNameNoSpace(n.Name)
	for i, action := range n.Actions {
		if i > 0 {
			w.buf.WriteByte(',')
		}
		if err := w.writeAlterTableAction(action); err != nil {
			return err
		}
	}
	return nil
}

func (w *writer) writeAlterTableAction(a *ast.AlterTableAction) error {
	switch a.Kind {
	case ast.AlterTableRename:
		w.buf.WriteString(" RENAME TO ")
		w.writeObjectNameNoSpace(a.NewName)
	case ast.AlterTableSwapWith:
		w.buf.WriteString(" SWAP WITH ")
		w.writeObjectNameNoSpace(a.NewName)
	case ast.AlterTableAddColumn:
		w.buf.WriteString(" ADD COLUMN")
		if a.IfNotExists {
			w.buf.WriteString(" IF NOT EXISTS")
		}
		for i, col := range a.Columns {
			if i > 0 {
				w.buf.WriteByte(',')
			}
			w.buf.WriteByte(' ')
			if err := w.writeColumnDef(col); err != nil {
				return err
			}
		}
	case ast.AlterTableDropColumn:
		w.buf.WriteString(" DROP COLUMN")
		if a.IfExists {
			w.buf.WriteString(" IF EXISTS")
		}
		for i, col := range a.DropColumnNames {
			if i > 0 {
				w.buf.WriteByte(',')
			}
			w.buf.WriteByte(' ')
			w.buf.WriteString(col.String())
		}
	case ast.AlterTableRenameColumn:
		w.buf.WriteString(" RENAME COLUMN ")
		w.buf.WriteString(a.OldName.String())
		w.buf.WriteString(" TO ")
		w.buf.WriteString(a.NewColName.String())
	case ast.AlterTableAlterColumn:
		for i, ca := range a.ColumnAlters {
			if i > 0 {
				w.buf.WriteByte(',')
			}
			// The column alters of one action are one list: ALTER COLUMN a
			// ..., COLUMN b ....
			if err := w.writeColumnAlter(ca, a.Continuation || i > 0); err != nil {
				return err
			}
		}
	case ast.AlterTableAddConstraint:
		w.buf.WriteString(" ADD")
		if a.Constraint != nil {
			w.buf.WriteByte(' ')
			if err := w.writeTableConstraint(a.Constraint); err != nil {
				return err
			}
		}
	case ast.AlterTableDropConstraint:
		if a.IsPrimaryKey {
			w.buf.WriteString(" DROP PRIMARY KEY")
		} else if a.DropUnique {
			w.buf.WriteString(" DROP UNIQUE (")
			// ConstraintName holds the unique key name here
			w.buf.WriteString(a.ConstraintName.String())
			w.buf.WriteByte(')')
		} else {
			w.buf.WriteString(" DROP CONSTRAINT ")
			w.buf.WriteString(a.ConstraintName.String())
		}
		if a.Cascade {
			w.buf.WriteString(" CASCADE")
		} else if a.Restrict {
			w.buf.WriteString(" RESTRICT")
		}
	case ast.AlterTableRenameConstraint:
		w.buf.WriteString(" RENAME CONSTRAINT ")
		w.buf.WriteString(a.ConstraintName.String())
		w.buf.WriteString(" TO ")
		w.buf.WriteString(a.NewConstraintName.String())
	case ast.AlterTableClusterBy:
		if err := w.writeClusterBy(a.ClusterBy, a.Linear); err != nil {
			return err
		}
	case ast.AlterTableDropClusterKey:
		w.buf.WriteString(" DROP CLUSTERING KEY")
	case ast.AlterTableRecluster:
		w.buf.WriteString(" RECLUSTER")
		if a.ReclusterMaxSize != nil {
			w.buf.WriteString(" MAX_SIZE = ")
			w.buf.WriteString(strconv.FormatInt(*a.ReclusterMaxSize, 10))
		}
		if a.ReclusterWhere != nil {
			w.buf.WriteString(" WHERE ")
			if err := writeExprNoLeadSpace(w, a.ReclusterWhere); err != nil {
				return err
			}
		}
	case ast.AlterTableSuspendRecluster:
		w.buf.WriteString(" SUSPEND RECLUSTER")
	case ast.AlterTableResumeRecluster:
		w.buf.WriteString(" RESUME RECLUSTER")
	case ast.AlterTableSet:
		w.buf.WriteString(" SET")
		for _, prop := range a.Props {
			w.buf.WriteByte(' ')
			w.buf.WriteString(prop.Name)
			if prop.Name == "CONTACT" {
				// Value is the verbatim `purpose = contact, ...` or `( ... )`.
				w.buf.WriteByte(' ')
			} else {
				w.buf.WriteString(" = ")
			}
			w.buf.WriteString(prop.Value)
		}
	case ast.AlterTableUnset:
		w.buf.WriteString(" UNSET ")
		w.buf.WriteString(strings.Join(a.UnsetProps, ", "))
	case ast.AlterTableSetTag:
		w.buf.WriteString(" SET TAG (")
		w.writeTagAssignments(a.Tags)
		w.buf.WriteByte(')')
	case ast.AlterTableUnsetTag:
		w.buf.WriteString(" UNSET TAG (")
		w.writeObjectNameList(a.UnsetTags)
		w.buf.WriteByte(')')
	case ast.AlterTableAddRowAccessPolicy:
		w.buf.WriteString(" ADD ROW ACCESS POLICY ")
		w.writeObjectNameNoSpace(a.PolicyName)
		if len(a.PolicyCols) > 0 {
			w.buf.WriteString(" ON (")
			for i, col := range a.PolicyCols {
				if i > 0 {
					w.buf.WriteString(", ")
				}
				w.buf.WriteString(col.String())
			}
			w.buf.WriteByte(')')
		}
	case ast.AlterTableDropRowAccessPolicy:
		w.buf.WriteString(" DROP ROW ACCESS POLICY ")
		w.writeObjectNameNoSpace(a.PolicyName)
	case ast.AlterTableDropAllRowAccessPolicies:
		w.buf.WriteString(" DROP ALL ROW ACCESS POLICIES")
	case ast.AlterTableAddSearchOpt:
		w.buf.WriteString(" ADD SEARCH OPTIMIZATION")
		if len(a.SearchOptOn) > 0 {
			w.buf.WriteString(" ON ")
			w.buf.WriteString(strings.Join(a.SearchOptOn, ", "))
		}
	case ast.AlterTableDropSearchOpt:
		w.buf.WriteString(" DROP SEARCH OPTIMIZATION")
		if len(a.SearchOptOn) > 0 {
			w.buf.WriteString(" ON ")
			w.buf.WriteString(strings.Join(a.SearchOptOn, ", "))
		}
	case ast.AlterTableSetMaskingPolicy:
		w.writeAlterColumnPrefix(a.Continuation)
		w.buf.WriteString(a.MaskColumn.String())
		w.buf.WriteString(" SET MASKING POLICY ")
		w.writeObjectNameNoSpace(a.MaskingPolicy)
	case ast.AlterTableUnsetMaskingPolicy:
		w.writeAlterColumnPrefix(a.Continuation)
		w.buf.WriteString(a.MaskColumn.String())
		w.buf.WriteString(" UNSET MASKING POLICY")
	case ast.AlterTableSetColumnTag:
		w.writeAlterColumnPrefix(a.Continuation)
		w.buf.WriteString(a.TagColumn.String())
		w.buf.WriteString(" SET TAG (")
		w.writeTagAssignments(a.Tags)
		w.buf.WriteByte(')')
	case ast.AlterTableUnsetColumnTag:
		w.writeAlterColumnPrefix(a.Continuation)
		w.buf.WriteString(a.TagColumn.String())
		w.buf.WriteString(" UNSET TAG (")
		w.writeObjectNameList(a.UnsetTags)
		w.buf.WriteByte(')')
	case ast.AlterTableAlterConstraint:
		w.buf.WriteString(" ALTER ")
		if a.ConstraintTarget != "" {
			w.buf.WriteString(a.ConstraintTarget)
		} else {
			w.buf.WriteString("CONSTRAINT ")
			w.buf.WriteString(a.ConstraintName.String())
		}
		w.writeParenIdentList(a.ConstraintCols)
		for _, prop := range a.ConstraintProps {
			w.buf.WriteByte(' ')
			w.buf.WriteString(prop)
		}
	case ast.AlterTableSetAggregationPolicy:
		w.buf.WriteString(" SET AGGREGATION POLICY ")
		w.writeObjectNameNoSpace(a.PolicyName)
		if len(a.PolicyKeyCols) > 0 {
			w.buf.WriteString(" ENTITY KEY")
			w.writeParenIdentList(a.PolicyKeyCols)
		}
		if a.PolicyForce {
			w.buf.WriteString(" FORCE")
		}
	case ast.AlterTableUnsetAggregationPolicy:
		w.buf.WriteString(" UNSET AGGREGATION POLICY")
	case ast.AlterTableSetJoinPolicy:
		w.buf.WriteString(" SET JOIN POLICY ")
		w.writeObjectNameNoSpace(a.PolicyName)
		if a.PolicyForce {
			w.buf.WriteString(" FORCE")
		}
	case ast.AlterTableUnsetJoinPolicy:
		w.buf.WriteString(" UNSET JOIN POLICY")
	case ast.AlterTableSetProjectionPolicy:
		w.writeAlterColumnPrefix(a.Continuation)
		w.buf.WriteString(a.ProjectionColumn.String())
		w.buf.WriteString(" SET PROJECTION POLICY ")
		w.writeObjectNameNoSpace(a.ProjectionPolicy)
		if a.PolicyForce {
			w.buf.WriteString(" FORCE")
		}
	case ast.AlterTableUnsetProjectionPolicy:
		w.writeAlterColumnPrefix(a.Continuation)
		w.buf.WriteString(a.ProjectionColumn.String())
		w.buf.WriteString(" UNSET PROJECTION POLICY")
	case ast.AlterTableAddStorageLifecyclePolicy:
		w.buf.WriteString(" ADD STORAGE LIFECYCLE POLICY ")
		w.writeObjectNameNoSpace(a.PolicyName)
		w.buf.WriteString(" ON")
		w.writeParenIdentList(a.PolicyCols)
	case ast.AlterTableDropStorageLifecyclePolicy:
		w.buf.WriteString(" DROP STORAGE LIFECYCLE POLICY")
	case ast.AlterTableDataMetric:
		w.buf.WriteByte(' ')
		return w.writeDataMetricAction(a.DataMetric)
	case ast.AlterTableRefresh:
		w.buf.WriteString(" REFRESH")
		if a.RefreshPath != nil {
			w.buf.WriteByte(' ')
			w.buf.WriteString(quoteString(*a.RefreshPath))
		}
	case ast.AlterTableConvertToManaged:
		w.buf.WriteString(" CONVERT TO MANAGED")
		return w.writeCopyOptions(a.Options)
	case ast.AlterTableOptimize:
		w.buf.WriteString(" OPTIMIZE (")
		w.buf.WriteString(strings.Join(a.OptimizeOps, ", "))
		w.buf.WriteByte(')')
	case ast.AlterTableAddPartitionBy, ast.AlterTableDropPartitionBy, ast.AlterTableReplacePartitionBy:
		switch a.Kind {
		case ast.AlterTableAddPartitionBy:
			w.buf.WriteString(" ADD")
		case ast.AlterTableDropPartitionBy:
			w.buf.WriteString(" DROP")
		default:
			w.buf.WriteString(" REPLACE")
		}
		w.buf.WriteString(" PARTITION BY (")
		if err := writeExprNoLeadSpace(w, a.Partition); err != nil {
			return err
		}
		w.buf.WriteByte(')')
		if a.Kind == ast.AlterTableReplacePartitionBy {
			w.buf.WriteString(" WITH (")
			if err := writeExprNoLeadSpace(w, a.NewPartition); err != nil {
				return err
			}
			w.buf.WriteByte(')')
		}
	case ast.AlterTableCreateVersionTag, ast.AlterTableAlterVersionTag, ast.AlterTableReplaceVersionTag:
		switch a.Kind {
		case ast.AlterTableCreateVersionTag:
			w.buf.WriteString(" CREATE")
			if a.OrReplace {
				w.buf.WriteString(" OR REPLACE")
			}
			w.buf.WriteString(" VERSION_TAG")
			if a.IfNotExists {
				w.buf.WriteString(" IF NOT EXISTS")
			}
		case ast.AlterTableAlterVersionTag:
			w.buf.WriteString(" ALTER VERSION_TAG")
		default:
			w.buf.WriteString(" REPLACE VERSION_TAG")
		}
		w.buf.WriteByte(' ')
		w.writeObjectNameNoSpace(a.VersionTag)
		if a.AsOfKind != "" {
			w.buf.WriteString(" AS OF ")
			w.buf.WriteString(a.AsOfKind)
			w.buf.WriteByte(' ')
			w.buf.WriteString(a.AsOfValue)
		}
		if a.RetainCount != nil {
			w.buf.WriteString(" RETAIN ")
			w.buf.WriteString(strconv.FormatInt(*a.RetainCount, 10))
			w.buf.WriteByte(' ')
			w.buf.WriteString(a.RetainUnit)
		}
	case ast.AlterTableDropVersionTag:
		w.buf.WriteString(" DROP VERSION_TAG")
		if a.IfExists {
			w.buf.WriteString(" IF EXISTS")
		}
		w.buf.WriteByte(' ')
		w.writeObjectNameNoSpace(a.VersionTag)
	case ast.AlterTableUpgradeIcebergVersion:
		w.buf.WriteString(" UPGRADE TO ICEBERG VERSION ")
		w.buf.WriteString(a.IcebergVersion)
	default:
		return fmt.Errorf("deparse: unsupported ALTER TABLE action kind %d", a.Kind)
	}
	return nil
}

// writeParenIdentList writes " (a, b, ...)".
func (w *writer) writeParenIdentList(ids []ast.Ident) {
	w.buf.WriteString(" (")
	for i, id := range ids {
		if i > 0 {
			w.buf.WriteString(", ")
		}
		w.buf.WriteString(id.String())
	}
	w.buf.WriteByte(')')
}

// writeDataMetricAction writes a data metric function action (shared by the
// ALTER TABLE / VIEW / DYNAMIC TABLE / MATERIALIZED VIEW writers) without a
// leading space.
func (w *writer) writeDataMetricAction(dm *ast.DataMetricAction) error {
	if dm == nil {
		return fmt.Errorf("deparse: missing data metric action")
	}
	switch dm.Kind {
	case ast.DataMetricSetSchedule:
		w.buf.WriteString("SET DATA_METRIC_SCHEDULE = ")
		w.buf.WriteString(quoteString(dm.Schedule))
		return nil
	case ast.DataMetricUnsetSchedule:
		w.buf.WriteString("UNSET DATA_METRIC_SCHEDULE")
		return nil
	case ast.DataMetricAdd:
		w.buf.WriteString("ADD")
	case ast.DataMetricDrop:
		w.buf.WriteString("DROP")
	case ast.DataMetricModify:
		w.buf.WriteString("MODIFY")
	default:
		return fmt.Errorf("deparse: unsupported data metric action kind %d", dm.Kind)
	}
	w.buf.WriteString(" DATA METRIC FUNCTION ")
	for i, fn := range dm.Functions {
		if i > 0 {
			w.buf.WriteString(", ")
		}
		if err := w.writeDataMetricFunctionClause(fn); err != nil {
			return err
		}
	}
	return nil
}

func (w *writer) writeDataMetricFunctionClause(fn *ast.DataMetricFunctionClause) error {
	w.writeObjectNameNoSpace(fn.Name)
	if fn.Modify == ast.DataMetricModifySet {
		w.buf.WriteString(" SET")
		return w.writeCopyOptions(fn.Options)
	}
	if fn.HasOn {
		w.buf.WriteString(" ON (")
		for i, col := range fn.Columns {
			if i > 0 {
				w.buf.WriteString(", ")
			}
			w.buf.WriteString(col.String())
		}
		if fn.RefTable != nil {
			if len(fn.Columns) > 0 {
				w.buf.WriteString(", ")
			}
			w.buf.WriteString("TABLE ")
			w.writeObjectNameNoSpace(fn.RefTable)
			w.writeParenIdentList(fn.RefColumns)
		}
		w.buf.WriteByte(')')
	}
	if len(fn.WithinGroup) > 0 {
		w.buf.WriteString(" WITHIN GROUP")
		w.writeParenIdentList(fn.WithinGroup)
		if fn.GroupLimit != nil {
			w.buf.WriteString(" GROUP LIMIT ")
			w.buf.WriteString(strconv.FormatInt(*fn.GroupLimit, 10))
		}
	}
	switch fn.Modify {
	case ast.DataMetricModifySuspend:
		w.buf.WriteString(" SUSPEND")
	case ast.DataMetricModifyResume:
		w.buf.WriteString(" RESUME")
	case ast.DataMetricModifyAddExpectation:
		w.buf.WriteString(" ADD")
	case ast.DataMetricModifyModifyExpectation:
		w.buf.WriteString(" MODIFY")
	case ast.DataMetricModifyDropExpectation:
		w.buf.WriteString(" DROP")
	}
	if fn.HasFilter {
		w.buf.WriteString(" FILTER (")
		if fn.Filter != nil {
			if err := writeExprNoLeadSpace(w, fn.Filter); err != nil {
				return err
			}
		}
		w.buf.WriteByte(')')
	}
	if len(fn.Expectations) > 0 {
		w.buf.WriteString(" EXPECTATION ")
		for i, exp := range fn.Expectations {
			if i > 0 {
				w.buf.WriteString(", ")
			}
			w.buf.WriteString(exp.Name.String())
			if exp.Expr != nil {
				w.buf.WriteString(" (")
				if err := writeExprNoLeadSpace(w, exp.Expr); err != nil {
					return err
				}
				w.buf.WriteByte(')')
			}
		}
	}
	if fn.ExecuteAs.Name != "" {
		w.buf.WriteString(" EXECUTE AS ROLE ")
		w.buf.WriteString(fn.ExecuteAs.String())
	}
	return w.writeCopyOptions(fn.Options)
}

// writeAlterColumnPrefix writes " ALTER COLUMN ", or " COLUMN " for an action
// that continues the previous ALTER / MODIFY column list.
func (w *writer) writeAlterColumnPrefix(continuation bool) {
	if continuation {
		w.buf.WriteString(" COLUMN ")
		return
	}
	w.buf.WriteString(" ALTER COLUMN ")
}

func (w *writer) writeColumnAlter(ca *ast.ColumnAlter, continuation bool) error {
	w.writeAlterColumnPrefix(continuation)
	w.buf.WriteString(ca.Column.String())
	switch ca.Kind {
	case ast.ColumnAlterSetDataType:
		w.buf.WriteString(" SET DATA TYPE ")
		w.writeTypeName(ca.DataType)
		if ca.RenameFields {
			w.buf.WriteString(" RENAME FIELDS")
		}
	case ast.ColumnAlterSetDefault:
		w.buf.WriteString(" SET DEFAULT ")
		if err := writeExprNoLeadSpace(w, ca.DefaultExpr); err != nil {
			return err
		}
	case ast.ColumnAlterDropDefault:
		w.buf.WriteString(" DROP DEFAULT")
	case ast.ColumnAlterSetNotNull:
		w.buf.WriteString(" SET NOT NULL")
	case ast.ColumnAlterDropNotNull:
		w.buf.WriteString(" DROP NOT NULL")
	case ast.ColumnAlterSetComment:
		w.buf.WriteString(" COMMENT ")
		if ca.Comment != nil {
			w.buf.WriteString(quoteString(*ca.Comment))
		}
	case ast.ColumnAlterUnsetComment:
		w.buf.WriteString(" UNSET COMMENT")
	case ast.ColumnAlterSetWriteDefault:
		w.buf.WriteString(" SET WRITE DEFAULT ")
		if err := writeExprNoLeadSpace(w, ca.DefaultExpr); err != nil {
			return err
		}
	case ast.ColumnAlterDropWriteDefault:
		w.buf.WriteString(" DROP WRITE DEFAULT")
	default:
		return fmt.Errorf("deparse: unsupported ColumnAlterKind %d", ca.Kind)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// writeDBSchemaProps writes optional database/schema properties: the typed
// ones, then the open-ended KEY = value options. Contacts are written by the
// caller, whose syntax differs between CREATE and ALTER ... SET.
func (w *writer) writeDBSchemaProps(p *ast.DBSchemaProps) error {
	if p == nil {
		return nil
	}
	if p.DataRetention != nil {
		w.buf.WriteString(" DATA_RETENTION_TIME_IN_DAYS = ")
		w.buf.WriteString(strconv.FormatInt(*p.DataRetention, 10))
	}
	if p.MaxDataExt != nil {
		w.buf.WriteString(" MAX_DATA_EXTENSION_TIME_IN_DAYS = ")
		w.buf.WriteString(strconv.FormatInt(*p.MaxDataExt, 10))
	}
	if p.DefaultDDLCol != nil {
		w.buf.WriteString(" DEFAULT_DDL_COLLATION = ")
		w.buf.WriteString(quoteString(*p.DefaultDDLCol))
	}
	if p.Comment != nil {
		w.buf.WriteString(" COMMENT = ")
		w.buf.WriteString(quoteString(*p.Comment))
	}
	return w.writeCopyOptions(p.Options)
}

// writeDBSchemaSetProps writes the property list of ALTER { DATABASE | SCHEMA }
// ... SET, with contacts in the inline CONTACT <purpose> = <contact> form.
func (w *writer) writeDBSchemaSetProps(p *ast.DBSchemaProps) error {
	if err := w.writeDBSchemaProps(p); err != nil {
		return err
	}
	if p != nil && len(p.Contacts) > 0 {
		w.buf.WriteString(" CONTACT ")
		w.writeContacts(p.Contacts)
	}
	return nil
}

// writeCreateDBSchemaTail writes the trailing clauses of CREATE DATABASE /
// CREATE SCHEMA: properties, WITH TAG ( ... ) and WITH CONTACT ( ... ).
func (w *writer) writeCreateDBSchemaTail(p *ast.DBSchemaProps, tags []*ast.TagAssignment) error {
	if err := w.writeDBSchemaProps(p); err != nil {
		return err
	}
	if len(tags) > 0 {
		w.buf.WriteString(" WITH TAG (")
		w.writeTagAssignments(tags)
		w.buf.WriteByte(')')
	}
	if len(p.Contacts) > 0 {
		w.buf.WriteString(" WITH CONTACT (")
		w.writeContacts(p.Contacts)
		w.buf.WriteByte(')')
	}
	return nil
}

// writeContacts writes <purpose> = <contact> [ , ... ].
func (w *writer) writeContacts(contacts []*ast.CopyOption) {
	for i, c := range contacts {
		if i > 0 {
			w.buf.WriteString(", ")
		}
		w.buf.WriteString(c.Name)
		w.buf.WriteString(" = ")
		w.buf.WriteString(c.Words)
	}
}

// writeCloneIgnoreClauses writes the IGNORE clauses of a database/schema CLONE.
func (w *writer) writeCloneIgnoreClauses(insufficient, hybrid bool) {
	if insufficient {
		w.buf.WriteString(" IGNORE TABLES WITH INSUFFICIENT DATA RETENTION")
	}
	if hybrid {
		w.buf.WriteString(" IGNORE HYBRID TABLES")
	}
}

// writeFromBackupSet writes FROM BACKUP SET <set> IDENTIFIER '<id>'.
func (w *writer) writeFromBackupSet(set *ast.ObjectName, id string) {
	w.buf.WriteString(" FROM BACKUP SET ")
	w.writeObjectNameNoSpace(set)
	w.buf.WriteString(" IDENTIFIER ")
	w.buf.WriteString(quoteString(id))
}

// writeLinkedCatalogUpdate writes UPDATE LINKED_CATALOG ... of ALTER DATABASE.
func (w *writer) writeLinkedCatalogUpdate(u *ast.LinkedCatalogUpdate) error {
	if u == nil {
		return fmt.Errorf("deparse: ALTER DATABASE UPDATE LINKED_CATALOG without a body")
	}
	w.buf.WriteString(" UPDATE LINKED_CATALOG ")
	switch u.Op {
	case "ADD", "REMOVE":
		w.buf.WriteString(u.Op)
		w.buf.WriteString(" (")
		for i, ns := range u.Namespaces {
			if i > 0 {
				w.buf.WriteString(", ")
			}
			w.buf.WriteString(quoteString(ns))
		}
		if u.Op == "ADD" {
			w.buf.WriteString(") TO ")
		} else {
			w.buf.WriteString(") FROM ")
		}
		w.buf.WriteString(u.Target)
	case "UNSET":
		w.buf.WriteString("UNSET ")
		w.buf.WriteString(u.Target)
	case "SET":
		w.buf.WriteString("SET")
		return w.writeCopyOptions(u.Options)
	default:
		return fmt.Errorf("deparse: unsupported UPDATE LINKED_CATALOG operation %q", u.Op)
	}
	return nil
}

// writeObjectNameList writes a comma-separated list of object names.
func (w *writer) writeObjectNameList(names []*ast.ObjectName) {
	for i, n := range names {
		if i > 0 {
			w.buf.WriteString(", ")
		}
		w.writeObjectNameNoSpace(n)
	}
}

// ---------------------------------------------------------------------------
// WAREHOUSE DDL (gap-warehouse)
// ---------------------------------------------------------------------------

// writeCopyOption emits one open-ended `KEY = value` option (or a bare KEY),
// covering every CopyOption value shape: a literal, a verbatim word run, a
// parenthesized key/value group, or a parenthesized literal list. It is shared
// by the warehouse writers; option names are emitted verbatim (already
// uppercased by the parser).
func (w *writer) writeCopyOption(opt *ast.CopyOption) error {
	w.buf.WriteString(opt.Name)
	if opt.Bare {
		return nil
	}
	w.buf.WriteString(" = ")
	switch {
	case opt.Lit != nil:
		return w.writeCopyOptionLiteral(opt.Lit)
	case opt.Words != "":
		w.buf.WriteString(opt.Words)
	case opt.Group != nil:
		w.buf.WriteByte('(')
		for i, sub := range opt.Group {
			if i > 0 {
				w.buf.WriteString(", ")
			}
			if err := w.writeCopyOption(sub); err != nil {
				return err
			}
		}
		w.buf.WriteByte(')')
	case opt.List != nil:
		w.buf.WriteByte('(')
		for i, lit := range opt.List {
			if i > 0 {
				w.buf.WriteString(", ")
			}
			if err := w.writeCopyOptionLiteral(lit); err != nil {
				return err
			}
		}
		w.buf.WriteByte(')')
	default:
		return fmt.Errorf("deparse: empty option value for %q", opt.Name)
	}
	return nil
}

// writeCopyOptionLiteral emits a CopyOption literal value without a leading
// space (the caller has already written " = "). String literals are single
// quoted; numeric literals use their raw text / integer value.
func (w *writer) writeCopyOptionLiteral(lit *ast.Literal) error {
	switch lit.Kind {
	case ast.LitString:
		w.buf.WriteString(quoteString(lit.Value))
	case ast.LitInt:
		w.buf.WriteString(strconv.FormatInt(lit.Ival, 10))
	case ast.LitFloat:
		w.buf.WriteString(lit.Value)
	case ast.LitBool:
		if lit.Bval {
			w.buf.WriteString("TRUE")
		} else {
			w.buf.WriteString("FALSE")
		}
	case ast.LitNull:
		w.buf.WriteString("NULL")
	default:
		w.buf.WriteString(lit.Value)
	}
	return nil
}

// writeCopyOptions emits a space-separated run of options.
func (w *writer) writeCopyOptions(opts []*ast.CopyOption) error {
	for _, opt := range opts {
		w.buf.WriteByte(' ')
		if err := w.writeCopyOption(opt); err != nil {
			return err
		}
	}
	return nil
}

func (w *writer) writeCreateWarehouseStmt(n *ast.CreateWarehouseStmt) error {
	w.buf.WriteString("CREATE")
	if n.OrReplace {
		w.buf.WriteString(" OR REPLACE")
	} else if n.OrAlter {
		w.buf.WriteString(" OR ALTER")
	}
	w.buf.WriteString(" WAREHOUSE")
	if n.IfNotExists {
		w.buf.WriteString(" IF NOT EXISTS")
	}
	w.buf.WriteByte(' ')
	w.writeObjectNameNoSpace(n.Name)
	if err := w.writeCopyOptions(n.Options); err != nil {
		return err
	}
	if len(n.Tags) > 0 {
		w.buf.WriteString(" WITH TAG (")
		w.writeTagAssignments(n.Tags)
		w.buf.WriteByte(')')
	}
	return nil
}

func (w *writer) writeAlterWarehouseStmt(n *ast.AlterWarehouseStmt) error {
	w.buf.WriteString("ALTER WAREHOUSE")
	if n.IfExists {
		w.buf.WriteString(" IF EXISTS")
	}
	w.buf.WriteByte(' ')
	w.writeObjectNameNoSpace(n.Name)
	switch n.Action {
	case ast.AlterWarehouseSuspend:
		w.buf.WriteString(" SUSPEND")
	case ast.AlterWarehouseResume:
		w.buf.WriteString(" RESUME")
		if n.ResumeIfSuspended {
			w.buf.WriteString(" IF SUSPENDED")
		}
	case ast.AlterWarehouseAbort:
		w.buf.WriteString(" ABORT ALL QUERIES")
	case ast.AlterWarehouseRename:
		w.buf.WriteString(" RENAME TO ")
		w.writeObjectNameNoSpace(n.NewName)
	case ast.AlterWarehouseSet:
		w.buf.WriteString(" SET")
		if err := w.writeCopyOptions(n.Options); err != nil {
			return err
		}
	case ast.AlterWarehouseUnset:
		w.buf.WriteString(" UNSET ")
		w.buf.WriteString(strings.Join(n.UnsetKeys, ", "))
	case ast.AlterWarehouseSetTag:
		w.buf.WriteString(" SET TAG ")
		w.writeTagAssignments(n.Tags)
	case ast.AlterWarehouseUnsetTag:
		w.buf.WriteString(" UNSET TAG ")
		w.writeObjectNameList(n.UnsetTags)
	case ast.AlterWarehouseAddTables:
		w.buf.WriteString(" ADD TABLES (")
		w.writeObjectNameList(n.Tables)
		w.buf.WriteByte(')')
	case ast.AlterWarehouseRemoveTables:
		w.buf.WriteString(" REMOVE TABLES (")
		w.writeObjectNameList(n.Tables)
		w.buf.WriteByte(')')
	default:
		return fmt.Errorf("deparse: unsupported ALTER WAREHOUSE action %d", n.Action)
	}
	return nil
}

// ---------------------------------------------------------------------------
// CREATE NETWORK RULE / ALTER NETWORK RULE (gap-network-rule)
// ---------------------------------------------------------------------------

func (w *writer) writeCreateNetworkRuleStmt(n *ast.CreateNetworkRuleStmt) error {
	w.buf.WriteString("CREATE")
	if n.OrReplace {
		w.buf.WriteString(" OR REPLACE")
	}
	w.buf.WriteString(" NETWORK RULE")
	if n.IfNotExists {
		w.buf.WriteString(" IF NOT EXISTS")
	}
	w.buf.WriteByte(' ')
	w.writeObjectNameNoSpace(n.Name)
	if err := w.writeCopyOptions(n.Options); err != nil {
		return err
	}
	return nil
}

func (w *writer) writeAlterNetworkRuleStmt(n *ast.AlterNetworkRuleStmt) error {
	w.buf.WriteString("ALTER NETWORK RULE")
	if n.IfExists {
		w.buf.WriteString(" IF EXISTS")
	}
	w.buf.WriteByte(' ')
	w.writeObjectNameNoSpace(n.Name)
	switch n.Action {
	case ast.AlterNetworkRuleSet:
		w.buf.WriteString(" SET")
		if err := w.writeCopyOptions(n.Options); err != nil {
			return err
		}
	case ast.AlterNetworkRuleUnset:
		w.buf.WriteString(" UNSET ")
		w.buf.WriteString(strings.Join(n.UnsetKeys, ", "))
	default:
		return fmt.Errorf("deparse: unsupported ALTER NETWORK RULE action %d", n.Action)
	}
	return nil
}

// writeAlterSessionStmt emits ALTER SESSION { SET | UNSET } <session_parameter>
// ... (session parameters, gap-alter-family).
func (w *writer) writeAlterSessionStmt(n *ast.AlterSessionStmt) error {
	w.buf.WriteString("ALTER SESSION")
	switch n.Action {
	case ast.AlterSessionSet:
		w.buf.WriteString(" SET")
		if err := w.writeCopyOptions(n.Options); err != nil {
			return err
		}
	case ast.AlterSessionUnset:
		w.buf.WriteString(" UNSET ")
		w.buf.WriteString(strings.Join(n.UnsetKeys, ", "))
	default:
		return fmt.Errorf("deparse: unsupported ALTER SESSION action %d", n.Action)
	}
	return nil
}
