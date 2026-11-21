package completion

import (
	"context"
	"errors"
	"slices"
	"testing"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mongo"
)

func testSessionConfig(collections ...string) SessionConfig {
	tables := make([]*metadatapb.TableMetadata, 0, len(collections))
	for _, name := range collections {
		tables = append(tables, &metadatapb.TableMetadata{Name: name})
	}
	return SessionConfig{
		DefaultDatabase: "test",
		LoadDatabase: func(_ context.Context, name string) (*metadatapb.DatabaseSchemaMetadata, error) {
			if name != "test" {
				return nil, nil
			}
			return &metadatapb.DatabaseSchemaMetadata{
				Name: name,
				Schemas: []*metadatapb.SchemaMetadata{
					{Tables: tables},
				},
			}, nil
		},
	}
}

func TestLineOffsetToBytePos(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		line   int
		offset int
		want   int
	}{
		{"start", "db.", 1, 0, 0},
		{"end of line", "db.", 1, 3, 3},
		{"second line", "db.users.find()\ndb.", 2, 3, len("db.users.find()\ndb.")},
		{"offset past newline clamps to line end", "db.\ndb.users.", 1, 99, 3},
		{"utf16 surrogate pair", "a😀b", 1, 3, len("a😀")},
		{"invalid line falls back to end", "db.", -1, -1, 3},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := LineOffsetToBytePos(tc.input, tc.line, tc.offset); got != tc.want {
				t.Errorf("LineOffsetToBytePos(%q, %d, %d) = %d, want %d", tc.input, tc.line, tc.offset, got, tc.want)
			}
		})
	}
}

func TestCompleteSessionDotAccess(t *testing.T) {
	cfg := testSessionConfig("users", "orders", "my-collection", `a"b\c`, "$valid_1")
	results := CompleteSession(context.Background(), cfg, "db.", 1, 3)

	if !hasCandidateWithType(results, "users", CandidateCollection) {
		t.Error("expected bare collection 'users'")
	}
	if !hasCandidateWithType(results, "orders", CandidateCollection) {
		t.Error("expected bare collection 'orders'")
	}
	if !hasCandidateWithType(results, "$valid_1", CandidateCollection) {
		t.Error("expected bare collection '$valid_1'")
	}
	if !hasCandidateWithType(results, `getCollection("my-collection")`, CandidateCollection) {
		t.Errorf("expected getCollection(...) for 'my-collection', got %v", candidateTexts(results))
	}
	if !hasCandidateWithType(results, `getCollection("a\"b\\c")`, CandidateCollection) {
		t.Errorf("expected escaped getCollection(...) collection, got %v", candidateTexts(results))
	}
	for _, c := range results {
		if c.Type != CandidateCollection {
			continue
		}
		if _, err := mongo.Parse("db." + c.Text + ".find()"); err != nil {
			t.Errorf("inserting dot-access candidate %q produced unparseable statement: %v", c.Text, err)
		}
	}
	if !hasCandidateWithType(results, "getName", CandidateDbMethod) {
		t.Error("expected db method 'getName'")
	}
	if !slices.IsSortedFunc(results, compareSessionCandidates) {
		t.Error("expected deterministically sorted candidates")
	}
}

func TestCompleteSessionBracketAccess(t *testing.T) {
	cfg := testSessionConfig("users", "orders", "my-collection", `a"b\c`, `a'b\c`)

	t.Run("after open bracket without quote", func(t *testing.T) {
		results := CompleteSession(context.Background(), cfg, "db[", 1, 3)
		for _, want := range []string{`"users"]`, `"orders"]`, `"my-collection"]`, `"a\"b\\c"]`, `"a'b\\c"]`} {
			if !hasCandidateWithType(results, want, CandidateCollection) {
				t.Errorf("expected %q in results, got %v", want, candidateTexts(results))
			}
		}
	})

	t.Run("inside double-quoted bracket", func(t *testing.T) {
		results := CompleteSession(context.Background(), cfg, `db["`, 1, 4)
		for _, want := range []string{"users", "orders", "my-collection", `a\"b\\c`} {
			if !hasCandidateWithType(results, want, CandidateCollection) {
				t.Errorf("expected %q in results, got %v", want, candidateTexts(results))
			}
		}
	})

	t.Run("inside double-quoted bracket with prefix", func(t *testing.T) {
		results := CompleteSession(context.Background(), cfg, `db["us`, 1, 6)
		if len(results) != 1 || results[0].Text != "users" || results[0].Type != CandidateCollection {
			t.Errorf("got %+v, want single 'users' collection candidate", results)
		}
	})

	t.Run("inside single-quoted bracket", func(t *testing.T) {
		results := CompleteSession(context.Background(), cfg, `db['`, 1, 4)
		if !hasCandidateWithType(results, `a\'b\\c`, CandidateCollection) {
			t.Errorf("expected single-quote-escaped collection, got %v", candidateTexts(results))
		}
	})

	t.Run("dot access inside array literal is not treated as bracket access", func(t *testing.T) {
		stmt := "db.users.aggregate([{$match: {x: db."
		results := CompleteSession(context.Background(), cfg, stmt, 1, len(stmt))
		if !hasCandidateWithType(results, "users", CandidateCollection) {
			t.Errorf("expected bare 'users', got %v", candidateTexts(results))
		}
		if !hasCandidateWithType(results, `getCollection("my-collection")`, CandidateCollection) {
			t.Errorf("expected getCollection(...) for 'my-collection', got %v", candidateTexts(results))
		}
	})
}

func TestCompleteSessionMultiline(t *testing.T) {
	cfg := testSessionConfig("users", "orders", "my-collection")

	t.Run("db dot on second line", func(t *testing.T) {
		stmt := "db.users.find().limit(20)\ndb."
		results := CompleteSession(context.Background(), cfg, stmt, 2, 3)

		if !hasCandidateWithType(results, "users", CandidateCollection) {
			t.Error("expected collection 'users' on second line")
		}
		if !hasCandidateWithType(results, `getCollection("my-collection")`, CandidateCollection) {
			t.Error("expected getCollection(...) for 'my-collection' on second line")
		}
		if !hasCandidateWithType(results, "runCommand", CandidateDbMethod) {
			t.Error("expected db method 'runCommand' on second line")
		}
	})

	t.Run("collection dot on second line", func(t *testing.T) {
		stmt := "db.users.find().limit(20)\ndb.users."
		results := CompleteSession(context.Background(), cfg, stmt, 2, len("db.users."))
		for _, want := range []string{"aggregate", "bulkWrite", "find", "findOne"} {
			if !hasCandidateWithType(results, want, CandidateMethod) {
				t.Errorf("expected collection method %q on second line", want)
			}
		}
	})

	t.Run("cursor chain on second line", func(t *testing.T) {
		stmt := "db.users.find().limit(20)\ndb.users.find()."
		results := CompleteSession(context.Background(), cfg, stmt, 2, len("db.users.find()."))
		for _, want := range []string{"allowDiskUse", "batchSize", "limit", "sort"} {
			if !hasCandidateWithType(results, want, CandidateCursorMethod) {
				t.Errorf("expected cursor method %q on second line", want)
			}
		}
	})
}

func TestCompleteSessionContexts(t *testing.T) {
	cfg := testSessionConfig("users", "orders", "my-collection")

	t.Run("db dot with prefix filters collections and db methods", func(t *testing.T) {
		results := CompleteSession(context.Background(), cfg, "db.u", 1, 4)
		if !hasCandidateWithType(results, "users", CandidateCollection) {
			t.Errorf("expected collection 'users', got %v", candidateTexts(results))
		}
		if hasCandidate(results, "orders") {
			t.Error("did not expect 'orders' for prefix 'u'")
		}
		for _, want := range []string{"updateRole", "updateUser"} {
			if !hasCandidateWithType(results, want, CandidateDbMethod) {
				t.Errorf("expected db method %q", want)
			}
		}
	})

	t.Run("collection methods via dot and bracket access", func(t *testing.T) {
		for _, stmt := range []string{"db.users.", `db["users"].`} {
			results := CompleteSession(context.Background(), cfg, stmt, 1, len(stmt))
			for _, want := range []string{"aggregate", "bulkWrite", "count", "find", "findOne"} {
				if !hasCandidateWithType(results, want, CandidateMethod) {
					t.Errorf("%q: expected collection method %q", stmt, want)
				}
			}
		}
	})

	t.Run("cursor methods after find, chained sort, and bracket find", func(t *testing.T) {
		for _, stmt := range []string{
			"db.users.find().",
			"db.users.find().sort().",
			`db["users"].find().`,
		} {
			results := CompleteSession(context.Background(), cfg, stmt, 1, len(stmt))
			for _, want := range []string{"allowDiskUse", "batchSize", "limit", "sort"} {
				if !hasCandidateWithType(results, want, CandidateCursorMethod) {
					t.Errorf("%q: expected cursor method %q", stmt, want)
				}
			}
		}
	})
}

func TestCompleteSessionNoMetadata(t *testing.T) {
	configs := []struct {
		name string
		cfg  SessionConfig
	}{
		{"nil loader", SessionConfig{DefaultDatabase: "test"}},
		{"empty default db", SessionConfig{LoadDatabase: testSessionConfig("users").LoadDatabase}},
		{"loader error", SessionConfig{
			DefaultDatabase: "test",
			LoadDatabase: func(context.Context, string) (*metadatapb.DatabaseSchemaMetadata, error) {
				return nil, errors.New("boom")
			},
		}},
		{"nil metadata", SessionConfig{
			DefaultDatabase: "other",
			LoadDatabase:    testSessionConfig("users").LoadDatabase,
		}},
	}
	for _, tc := range configs {
		t.Run(tc.name, func(t *testing.T) {
			results := CompleteSession(context.Background(), tc.cfg, "db.", 1, 3)
			if !hasCandidateWithType(results, "getName", CandidateDbMethod) {
				t.Error("expected db methods even without metadata")
			}
			for _, c := range results {
				if c.Type == CandidateCollection {
					t.Errorf("unexpected collection candidate %q without metadata", c.Text)
				}
			}
		})
	}
}
