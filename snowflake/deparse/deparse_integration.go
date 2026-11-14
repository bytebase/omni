package deparse

import (
	"fmt"
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// Integrations read by their pages' properties
// ---------------------------------------------------------------------------
//
// The parser holds CATALOG, EXTERNAL ACCESS and NOTIFICATION INTEGRATION's
// CREATE, CATALOG and EXTERNAL ACCESS INTEGRATION's ALTER, and every
// integration's REFRESH to the forms their pages document, so those are the
// forms written here. The other kinds and actions keep open-ended options this
// writer cannot promise to read back the same, and stay unsupported.

// writeCreateIntegrationStmt writes
//
//	CREATE [ OR REPLACE ] <kind> [ IF NOT EXISTS ] <name> <property> [ ... ]
func (w *writer) writeCreateIntegrationStmt(n *ast.CreateIntegrationStmt) error {
	switch n.Kind {
	case ast.CatalogIntegration, ast.ExternalAccessIntegration, ast.NotificationIntegration:
	default:
		return fmt.Errorf("deparse: unsupported CREATE %s", n.Kind)
	}
	w.buf.WriteString("CREATE ")
	if n.OrReplace {
		w.buf.WriteString("OR REPLACE ")
	}
	w.buf.WriteString(n.Kind.String())
	if n.IfNotExists {
		w.buf.WriteString(" IF NOT EXISTS")
	}
	w.buf.WriteByte(' ')
	w.writeObjectNameNoSpace(n.Name)
	return w.writeIntegrationProperties(n.Options)
}

// writeAlterIntegrationStmt writes
//
//	ALTER <kind> [ IF EXISTS ] <name> SET <property> [ ... ] [ TAG <tag> = '<value>' [ , ... ] ]
//	ALTER <kind> [ IF EXISTS ] <name> SET TAG <tag> = '<value>' [ , ... ]
//	ALTER <kind> [ IF EXISTS ] <name> UNSET { <property> | TAG <tag> } [ , ... ]
//	ALTER [ SECURITY ] INTEGRATION [ IF EXISTS ] <name> REFRESH [ <what> ... ]
func (w *writer) writeAlterIntegrationStmt(n *ast.AlterIntegrationStmt) error {
	strict := n.Kind == ast.CatalogIntegration || n.Kind == ast.ExternalAccessIntegration
	if !strict && n.Action != ast.AlterIntegrationRefresh {
		return fmt.Errorf("deparse: unsupported ALTER %s", n.Kind)
	}
	w.buf.WriteString("ALTER ")
	if n.Kind == ast.StorageIntegration {
		// A bare ALTER INTEGRATION reads as a storage integration, and ALTER
		// INTEGRATION is right for either.
		w.buf.WriteString("INTEGRATION")
	} else {
		w.buf.WriteString(n.Kind.String())
	}
	if n.IfExists {
		w.buf.WriteString(" IF EXISTS")
	}
	w.buf.WriteByte(' ')
	w.writeObjectNameNoSpace(n.Name)

	switch n.Action {
	case ast.AlterIntegrationSet:
		w.buf.WriteString(" SET")
		if err := w.writeIntegrationProperties(n.Options); err != nil {
			return err
		}
		if len(n.Tags) > 0 {
			w.buf.WriteString(" TAG ")
			w.writeIntegrationTags(n.Tags)
		}
	case ast.AlterIntegrationSetTag:
		w.buf.WriteString(" SET TAG ")
		w.writeIntegrationTags(n.Tags)
	case ast.AlterIntegrationUnset, ast.AlterIntegrationUnsetTag:
		w.buf.WriteString(" UNSET ")
		w.buf.WriteString(strings.Join(n.UnsetProps, ", "))
		for i, t := range n.UnsetTags {
			if i > 0 || len(n.UnsetProps) > 0 {
				w.buf.WriteString(", ")
			}
			w.buf.WriteString("TAG ")
			w.writeObjectNameNoSpace(t)
		}
	case ast.AlterIntegrationRefresh:
		w.buf.WriteString(" REFRESH")
		for _, r := range n.Refresh {
			w.buf.WriteByte(' ')
			w.buf.WriteString(r)
		}
	default:
		return fmt.Errorf("deparse: unsupported ALTER %s action", n.Kind)
	}
	return nil
}

// writeIntegrationProperties writes each property as ` NAME = value`.
func (w *writer) writeIntegrationProperties(opts []*ast.CopyOption) error {
	for _, o := range opts {
		w.buf.WriteByte(' ')
		if err := w.writeIntegrationProperty(o); err != nil {
			return err
		}
	}
	return nil
}

// writeIntegrationProperty writes `NAME = value`, or the value alone for an
// unnamed group in a list of groups. A group of `KEY = value` entries is
// written with spaces between them, as the pages write it; a list of names,
// strings, string pairs or groups with commas.
func (w *writer) writeIntegrationProperty(o *ast.CopyOption) error {
	if o.Name != "" {
		w.buf.WriteString(o.Name)
		if o.Bare {
			return nil
		}
		w.buf.WriteString(" = ")
	}
	switch {
	case o.Lit != nil:
		return w.writeIntegrationLiteral(o.Lit)
	case o.Words != "":
		w.buf.WriteString(o.Words)
	case o.List != nil:
		w.buf.WriteByte('(')
		for i, lit := range o.List {
			if i > 0 {
				w.buf.WriteString(", ")
			}
			if err := w.writeIntegrationLiteral(lit); err != nil {
				return err
			}
		}
		w.buf.WriteByte(')')
	case o.Group != nil:
		sep := ", "
		if len(o.Group) > 0 && isIntegrationKeyEntry(o.Group[0]) {
			sep = " "
		}
		w.buf.WriteByte('(')
		for i, e := range o.Group {
			if i > 0 {
				w.buf.WriteString(sep)
			}
			if err := w.writeIntegrationProperty(e); err != nil {
				return err
			}
		}
		w.buf.WriteByte(')')
	}
	return nil
}

// isIntegrationKeyEntry reports whether e is a `KEY = value` entry of a group
// (a word key), rather than a bare name, a quoted-key pair or an unnamed group.
func isIntegrationKeyEntry(e *ast.CopyOption) bool {
	return e.Name != "" && !e.Bare && !strings.HasPrefix(e.Name, "'")
}

// writeIntegrationLiteral writes a property's literal. A string is written
// single-quoted with its backslashes and quotes escaped, so that a value read
// from a $$-quoted string (where a backslash is literal) reads back the same.
func (w *writer) writeIntegrationLiteral(lit *ast.Literal) error {
	if lit.Kind == ast.LitString {
		w.buf.WriteString(quoteIntegrationString(lit.Value))
		return nil
	}
	return w.writeCopyOptionLiteral(lit)
}

// writeIntegrationTags writes <tag> = '<value>' [ , ... ], each value quoted
// as writeIntegrationLiteral quotes a string.
func (w *writer) writeIntegrationTags(tags []*ast.TagAssignment) {
	for i, tag := range tags {
		if i > 0 {
			w.buf.WriteString(", ")
		}
		w.writeObjectNameNoSpace(tag.Name)
		w.buf.WriteString(" = ")
		w.buf.WriteString(quoteIntegrationString(tag.Value))
	}
}

// quoteIntegrationString spells s as a single-quoted string that reads back
// as s: a backslash and a quote are escaped.
func quoteIntegrationString(s string) string {
	return "'" + strings.ReplaceAll(strings.ReplaceAll(s, `\`, `\\`), "'", "''") + "'"
}
