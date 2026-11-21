package completion

import (
	"context"
	"strings"
	"unicode/utf8"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/trino/catalog"
	"github.com/bytebase/omni/trino/parser"
)

// SessionConfig configures multi-catalog, session-aware Trino SQL completion.
type SessionConfig struct {
	// DefaultCatalog is the active Trino catalog (Bytebase database), if any.
	DefaultCatalog string
	// DefaultSchema is the active schema within DefaultCatalog, if any.
	DefaultSchema string
	// LoadCatalog fetches metadata for a catalog name on demand.
	LoadCatalog func(ctx context.Context, name string) (*metadatapb.DatabaseSchemaMetadata, error)
	// ListCatalogNames enumerates known catalog names on the instance.
	ListCatalogNames func(ctx context.Context) ([]string, error)
}

// CompleteSession runs session-aware Trino auto-completion over statement at
// the 1-based line and 0-based UTF-16 code-unit caretOffset.
//
// All known catalog names are registered cheaply so catalog-level completion
// lists them, while full metadata is loaded only for DefaultCatalog and any
// catalog referenced in the statement at the caret.
func CompleteSession(ctx context.Context, cfg SessionConfig, statement string, caretLine, caretOffset int) []Candidate {
	bytePos := LineOffsetToBytePos(statement, caretLine, caretOffset)
	stmtScope := statementScopeAtCaret(statement, bytePos)
	cat := buildSessionCatalog(ctx, cfg, stmtScope)
	return Complete(statement, bytePos, cat)
}

// CatalogNeeded reports whether the current completion statement needs this
// catalog's full schema/table/column metadata loaded: true for the session
// default catalog and for any catalog whose name appears in lowerStmt.
func CatalogNeeded(normName, defaultCatalog, lowerStmt string) bool {
	if normName == "" {
		return false
	}
	if defaultCatalog != "" && normName == defaultCatalog {
		return true
	}
	// A catalog name containing a double quote is escaped as "" inside a quoted
	// reference, so its normalized form may not appear verbatim in lowerStmt;
	// load it unconditionally rather than dropping completion.
	if strings.Contains(normName, `"`) {
		return true
	}
	return strings.Contains(lowerStmt, strings.ToLower(normName))
}

// LineOffsetToBytePos converts a 1-based line and 0-based UTF-16 code-unit
// column offset into a byte offset into s, matching LSP/Monaco positions.
func LineOffsetToBytePos(s string, line, offset int) int {
	currentLine := 1
	col := 0
	for i := 0; i < len(s); {
		if currentLine == line && col >= offset {
			return i
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == '\n' {
			if currentLine == line {
				return i
			}
			currentLine++
			col = 0
		} else if r <= 0xFFFF {
			col++
		} else {
			col += 2
		}
		i += size
	}
	return len(s)
}

func statementScopeAtCaret(statement string, bytePos int) string {
	bytePos = max(0, min(bytePos, len(statement)))
	segs := parser.Split(statement)
	if len(segs) == 0 {
		return statement
	}
	for _, s := range segs {
		if bytePos >= s.ByteStart && bytePos <= s.ByteEnd {
			return s.Text
		}
	}
	// Caret sits outside all non-empty statements (e.g., after a trailing ';'):
	// avoid pulling catalog references from earlier statements.
	return ""
}

func buildSessionCatalog(ctx context.Context, cfg SessionConfig, statement string) *catalog.Catalog {
	if cfg.LoadCatalog == nil && cfg.ListCatalogNames == nil {
		return nil
	}

	var names []string
	if cfg.ListCatalogNames != nil {
		if listed, err := cfg.ListCatalogNames(ctx); err == nil {
			names = listed
		}
	}
	if len(names) == 0 && (cfg.DefaultCatalog == "" || cfg.LoadCatalog == nil) {
		return nil
	}

	lowerStmt := strings.ToLower(statement)
	defaultCat := catalog.Normalize(cfg.DefaultCatalog)

	cat := catalog.New()
	attemptedDefault := false
	for _, name := range names {
		norm := catalog.Normalize(name)
		cat.EnsureCatalog(norm)
		if !CatalogNeeded(norm, defaultCat, lowerStmt) {
			continue
		}
		if norm == defaultCat {
			attemptedDefault = true
		}
		if cfg.LoadCatalog == nil {
			continue
		}
		meta, err := cfg.LoadCatalog(ctx, name)
		if err != nil || meta == nil {
			continue
		}
		cat.LoadMetadata(meta)
	}

	if !attemptedDefault && cfg.DefaultCatalog != "" {
		cat.EnsureCatalog(defaultCat)
		if cfg.LoadCatalog != nil {
			if meta, err := cfg.LoadCatalog(ctx, cfg.DefaultCatalog); err == nil && meta != nil {
				cat.LoadMetadata(meta)
			}
		}
	}

	if cfg.DefaultCatalog != "" {
		cat.SetCurrentCatalog(defaultCat)
	}
	if cfg.DefaultSchema != "" {
		cat.SetCurrentSchema(catalog.Normalize(cfg.DefaultSchema))
	}
	return cat
}
