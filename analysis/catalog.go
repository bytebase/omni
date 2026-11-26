package analysis

import (
	"context"

	"github.com/bytebase/omni/metadata"
)

// Catalog is what the caller knows of an instance's synced databases. A nil
// function means the caller knows nothing of that kind.
type Catalog struct {
	// ListDatabaseNames lists the instance's databases.
	ListDatabaseNames func(ctx context.Context) ([]string, error)
	// GetDatabase returns a database's synced metadata, or nil for a database
	// the caller has not synced.
	GetDatabase func(ctx context.Context, name string) (*Database, error)
}

// Database is one synced database.
type Database struct {
	// Name is the database's name as synced, which a lookup that ignores case
	// can spell differently from the name asked for.
	Name     string
	Metadata *metadata.DatabaseSchemaMetadata
	// ObjectCaseSensitive reports whether the database compares schema and
	// table names as written.
	ObjectCaseSensitive bool
}

// MetadataLoader returns a function that loads a database's synced
// metadata, nil for a database the catalog does not have, or nil when the
// catalog cannot look databases up.
func (c Catalog) MetadataLoader() func(ctx context.Context, name string) (*metadata.DatabaseSchemaMetadata, error) {
	if c.GetDatabase == nil {
		return nil
	}
	return func(ctx context.Context, name string) (*metadata.DatabaseSchemaMetadata, error) {
		database, err := c.GetDatabase(ctx, name)
		if err != nil || database == nil {
			return nil, err
		}
		return database.Metadata, nil
	}
}

// CompletionRequest is the session a completion runs in.
type CompletionRequest struct {
	// ReadOnly asks for the candidates of a read-only editor, which leave out
	// what only a change could use.
	ReadOnly        bool
	DefaultDatabase string
	DefaultSchema   string
	Catalog         Catalog
}

// Candidate is one completion the editor offers.
type Candidate struct {
	Text       string
	Type       CandidateType
	Definition string
	Comment    string
}

// Completer offers the completions at a caret. An engine's completion package
// exports its own Engine value implementing it, apart from the analysis
// package's.
type Completer interface {
	// Complete returns the candidates at the caret, given as a 1-based line and
	// a 0-based offset in UTF-16 code units. Each candidate's Type is one the
	// caller can name; the caller orders them.
	Complete(ctx context.Context, request CompletionRequest, script string, caretLine, caretOffset int) []Candidate
}
