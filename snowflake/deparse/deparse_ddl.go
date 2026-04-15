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
	}
	if n.Transient {
		w.buf.WriteString(" TRANSIENT")
	} else if n.Temporary {
		w.buf.WriteString(" TEMPORARY")
	} else if n.Volatile {
		w.buf.WriteString(" VOLATILE")
	}
	w.buf.WriteString(" TABLE")
	if n.IfNotExists {
		w.buf.WriteString(" IF NOT EXISTS")
	}
	w.buf.WriteByte(' ')
	w.writeObjectNameNoSpace(n.Name)

	// LIKE
	if n.Like != nil {
		w.buf.WriteString(" LIKE ")
		w.writeObjectNameNoSpace(n.Like)
		return nil
	}

	// CLONE
	if n.Clone != nil {
		w.buf.WriteString(" CLONE ")
		w.writeObjectNameNoSpace(n.Clone.Source)
		w.writeCloneTimeTravelSuffix(n.Clone)
		return nil
	}

	// Column definitions and table constraints
	if len(n.Columns) > 0 || len(n.Constraints) > 0 {
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
		w.buf.WriteByte(')')
	}

	// AS SELECT
	if n.AsSelect != nil {
		w.buf.WriteString(" AS ")
		return w.writeQueryExpr(n.AsSelect)
	}

	// CLUSTER BY
	if err := w.writeClusterBy(n.ClusterBy, n.Linear); err != nil {
		return err
	}

	// COPY GRANTS
	if n.CopyGrants {
		w.buf.WriteString(" COPY GRANTS")
	}

	// COMMENT
	if n.Comment != nil {
		w.buf.WriteString(" COMMENT = ")
		w.buf.WriteString(quoteString(*n.Comment))
	}

	// WITH TAG
	if len(n.Tags) > 0 {
		w.buf.WriteString(" WITH TAG (")
		w.writeTagAssignments(n.Tags)
		w.buf.WriteByte(')')
	}

	return nil
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

// writeCloneTimeTravelSuffix writes the AT/BEFORE time travel qualifier.
// The value is always a string literal (single-quoted in SQL).
func (w *writer) writeCloneTimeTravelSuffix(c *ast.CloneSource) {
	if c.AtBefore == "" {
		return
	}
	if c.AtBefore == "BEFORE" {
		w.buf.WriteString(" BEFORE (")
	} else {
		w.buf.WriteString(" AT (")
	}
	w.buf.WriteString(c.Kind)
	w.buf.WriteString(" => ")
	w.buf.WriteString(quoteString(c.Value))
	w.buf.WriteByte(')')
}

// ---------------------------------------------------------------------------
// DATABASE DDL
// ---------------------------------------------------------------------------

func (w *writer) writeCreateDatabaseStmt(n *ast.CreateDatabaseStmt) error {
	w.buf.WriteString("CREATE")
	if n.OrReplace {
		w.buf.WriteString(" OR REPLACE")
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
		w.writeCloneTimeTravelSuffix(n.Clone)
	}
	w.writeDBSchemaProps(&n.Props)
	if len(n.Tags) > 0 {
		w.buf.WriteString(" WITH TAG (")
		w.writeTagAssignments(n.Tags)
		w.buf.WriteByte(')')
	}
	return nil
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
		if n.SetProps != nil {
			w.writeDBSchemaProps(n.SetProps)
		}
	case ast.AlterDBUnset:
		w.buf.WriteString(" UNSET ")
		w.buf.WriteString(strings.Join(n.UnsetProps, ", "))
	case ast.AlterDBSetTag:
		w.buf.WriteString(" SET TAG (")
		w.writeTagAssignments(n.Tags)
		w.buf.WriteByte(')')
	case ast.AlterDBUnsetTag:
		w.buf.WriteString(" UNSET TAG (")
		w.writeObjectNameList(n.UnsetTags)
		w.buf.WriteByte(')')
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
	}
	if n.Transient {
		w.buf.WriteString(" TRANSIENT")
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
		w.writeCloneTimeTravelSuffix(n.Clone)
	}
	if n.ManagedAccess {
		w.buf.WriteString(" WITH MANAGED ACCESS")
	}
	w.writeDBSchemaProps(&n.Props)
	if len(n.Tags) > 0 {
		w.buf.WriteString(" WITH TAG (")
		w.writeTagAssignments(n.Tags)
		w.buf.WriteByte(')')
	}
	return nil
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
		if n.SetProps != nil {
			w.writeDBSchemaProps(n.SetProps)
		}
	case ast.AlterSchemaUnset:
		w.buf.WriteString(" UNSET ")
		w.buf.WriteString(strings.Join(n.UnsetProps, ", "))
	case ast.AlterSchemaSetTag:
		w.buf.WriteString(" SET TAG (")
		w.writeTagAssignments(n.Tags)
		w.buf.WriteByte(')')
	case ast.AlterSchemaUnsetTag:
		w.buf.WriteString(" UNSET TAG (")
		w.writeObjectNameList(n.UnsetTags)
		w.buf.WriteByte(')')
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
	w.buf.WriteString(n.Kind.String())
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

func (w *writer) writeUndropStmt(n *ast.UndropStmt) error {
	w.buf.WriteString("UNDROP ")
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
	default:
		return fmt.Errorf("deparse: unsupported ALTER VIEW action %d", n.Action)
	}
	return nil
}

// ---------------------------------------------------------------------------
// CREATE MATERIALIZED VIEW / ALTER MATERIALIZED VIEW
// ---------------------------------------------------------------------------

func (w *writer) writeCreateMaterializedViewStmt(n *ast.CreateMaterializedViewStmt) error {
	w.buf.WriteString("CREATE")
	if n.OrReplace {
		w.buf.WriteString(" OR REPLACE")
	}
	if n.Secure {
		w.buf.WriteString(" SECURE")
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
	case ast.AlterMVSetSecure:
		w.buf.WriteString(" SET SECURE")
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
	w.buf.WriteString("ALTER TABLE")
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
			if err := w.writeColumnAlter(ca); err != nil {
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
			w.buf.WriteString(" = ")
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
		w.buf.WriteString(" ALTER COLUMN ")
		w.buf.WriteString(a.MaskColumn.String())
		w.buf.WriteString(" SET MASKING POLICY ")
		w.writeObjectNameNoSpace(a.MaskingPolicy)
	case ast.AlterTableUnsetMaskingPolicy:
		w.buf.WriteString(" ALTER COLUMN ")
		w.buf.WriteString(a.MaskColumn.String())
		w.buf.WriteString(" UNSET MASKING POLICY")
	case ast.AlterTableSetColumnTag:
		w.buf.WriteString(" ALTER COLUMN ")
		w.buf.WriteString(a.TagColumn.String())
		w.buf.WriteString(" SET TAG (")
		w.writeTagAssignments(a.Tags)
		w.buf.WriteByte(')')
	case ast.AlterTableUnsetColumnTag:
		w.buf.WriteString(" ALTER COLUMN ")
		w.buf.WriteString(a.TagColumn.String())
		w.buf.WriteString(" UNSET TAG (")
		w.writeObjectNameList(a.UnsetTags)
		w.buf.WriteByte(')')
	default:
		return fmt.Errorf("deparse: unsupported ALTER TABLE action kind %d", a.Kind)
	}
	return nil
}

func (w *writer) writeColumnAlter(ca *ast.ColumnAlter) error {
	w.buf.WriteString(" ALTER COLUMN ")
	w.buf.WriteString(ca.Column.String())
	switch ca.Kind {
	case ast.ColumnAlterSetDataType:
		w.buf.WriteString(" SET DATA TYPE ")
		w.writeTypeName(ca.DataType)
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
	default:
		return fmt.Errorf("deparse: unsupported ColumnAlterKind %d", ca.Kind)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// writeDBSchemaProps writes optional database/schema properties.
func (w *writer) writeDBSchemaProps(p *ast.DBSchemaProps) {
	if p == nil {
		return
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
