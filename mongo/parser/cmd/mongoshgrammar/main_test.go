package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytebase/omni/enginetest/grammar"
)

// TestExtractAndRender pins the converter's TypeScript extraction and ANTLR 4
// grammar rendering on synthetic vendor source files.
func TestExtractAndRender(t *testing.T) {
	sources := map[string][]byte{}
	for _, rel := range vendorFiles {
		sources[rel] = []byte("// empty\n")
	}
	sources["packages/shell-api/src/shell-api.ts"] = []byte(`
export class ShellConfig extends ShellApiClass {
  async set<K extends keyof ShellUserConfig>(key: K, value: any) {}
  async get<K extends keyof ShellUserConfig>(key: K) {}
}
export default class ShellApi extends ShellApiClass {
  @directShellCommand
  use(db: string) {}
  @directShellCommand
  async show(cmd: string, arg?: string) {}
  @directShellCommand
  async exit(exitCode?: number) {}
  @directShellCommand
  async it() {}
  async sleep(ms: number) {}
  async _internal() {}
  private hidden() {}
}
`)
	sources["packages/shell-api/src/mongo.ts"] = []byte(`
export default class Mongo {
  getDB<K extends StringKey<M>>(db: K) {}
  async show(cmd: string, arg?: string) {
    switch (cmd) {
      case 'databases':
      case 'dbs':
      case 'collections':
        break;
    }
  }
}
`)
	sources["packages/shell-api/src/database.ts"] = []byte(`
export class Database {
  getSiblingDB<K extends StringKey<M>>(db: K) {}
  async createCollection(name: string) {}
}
`)
	sources["packages/shell-api/src/collection.ts"] = []byte(`
export class Collection {
  find(query?: any) {}
  async insertOne(doc: any) {}
}
`)
	sources["packages/shell-api/src/abstract-cursor.ts"] = []byte(`
export abstract class BaseCursor {
  async hasNext() {}
}
export abstract class AbstractFiniteCursor {
  maxTimeMS(v: number) {}
}
`)
	sources["packages/shell-api/src/aggregate-or-find-cursor.ts"] = []byte(`
export abstract class AggregateOrFindCursor {
  sort(spec: any) {}
}
`)
	sources["packages/shell-api/src/cursor.ts"] = []byte(`
export default class Cursor {
  async hasNext() {}
  limit(n: number) {}
}
`)
	sources["packages/shell-api/src/change-stream-cursor.ts"] = []byte(`
export default class ChangeStreamCursor {
  getResumeToken() {}
}
`)
	sources["packages/shell-api/src/explainable-cursor.ts"] = []byte(`
export default class ExplainableCursor {
  finish() {}
}
`)
	sources["packages/shell-api/src/explainable.ts"] = []byte(`
export default class Explainable {
  async find(q?: any) {}
}
`)
	sources["packages/shell-api/src/bulk.ts"] = []byte(`
export class BulkFindOp {
  updateOne(u: any) {}
}
export default class Bulk {
  async execute() {}
}
`)
	sources["packages/shell-api/src/plan-cache.ts"] = []byte(`
export default class PlanCache {
  async clear() {}
}
`)
	sources["packages/shell-api/src/replica-set.ts"] = []byte(`
export default class ReplicaSet {
  async status() {}
}
`)
	sources["packages/shell-api/src/shard.ts"] = []byte(`
export default class Shard {
  async status() {}
}
`)
	sources["packages/shell-api/src/streams.ts"] = []byte(`
export class Streams {
  async listStreamProcessors(filter: any) {}
}
`)
	sources["packages/shell-api/src/stream-processor.ts"] = []byte(`
export class StreamProcessor {
  async start() {}
}
`)
	sources["packages/shell-api/src/field-level-encryption.ts"] = []byte(`
export class ClientEncryption {
  async encrypt(v: any) {}
}
export class KeyVault {
  async createKey(kms: string) {}
}
`)
	sources["packages/shell-api/src/session.ts"] = []byte(`
export default class Session {
  startTransaction() {}
}
`)
	sources["packages/shell-bson/src/shell-bson.ts"] = []byte(`
export interface ShellBsonBase {
  ObjectId: any;
  Binary: any;
  EJSON: any;
}
export function constructShellBson() {
  return {
    ObjectId: assignAll(fn, pickWithExactKeyMatch(bson.ObjectId, ['prototype', 'cacheHexString', 'createFromHexString'])),
    Binary: assignAll(fn, pickWithExactKeyMatch(bson.Binary, ['prototype', 'createFromBase64', 'SUBTYPE_UUID'])),
    EJSON: pickWithExactKeyMatch(bson.EJSON, ['parse', 'stringify']),
  };
}
`)

	digest, err := filesDigest(sources)
	if err != nil {
		t.Fatal(err)
	}
	ex, err := extract(sources)
	if err != nil {
		t.Fatal(err)
	}
	got, err := renderGrammar("v0.0.0", "deadbeef", digest, ex)
	if err != nil {
		t.Fatal(err)
	}

	for _, want := range []string{
		"grammar MongoShell;",
		"'use' IDENT // packages/shell-api/src/shell-api.ts:",
		"'show' showTarget // packages/shell-api/src/shell-api.ts:",
		"'databases' // packages/shell-api/src/mongo.ts:",
		"'getSiblingDB' // packages/shell-api/src/database.ts:",
		"'ObjectId' '.' 'createFromHexString'",
		"'Binary' '.' 'createFromBase64'",
		"'parse' // packages/shell-bson/src/shell-bson.ts:",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("rendered grammar missing %q:\n%s", want, got)
		}
	}
	for _, forbid := range []string{"_internal", "hidden", "SUBTYPE_UUID", "cacheHexString", "prototype"} {
		if strings.Contains(got, "'"+forbid+"'") {
			t.Errorf("rendered grammar should not contain %q", forbid)
		}
	}
}

// TestCommittedMongoShellG4ParsesCleanly verifies that the committed
// mongo/testdata/MongoShell.g4 parses cleanly through enginetest/grammar
// with mongo/testdata/tokens.tsv and leaves zero uncovered alternatives.
func TestCommittedMongoShellG4ParsesCleanly(t *testing.T) {
	td := filepath.Join("..", "..", "..", "testdata")
	src, err := os.ReadFile(filepath.Join(td, "MongoShell.g4"))
	if err != nil {
		t.Fatal(err)
	}
	tokens, err := grammar.ReadTokens(filepath.Join(td, "tokens.tsv"))
	if err != nil {
		t.Fatal(err)
	}
	g, err := grammar.ReadANTLR(string(src), tokens)
	if err != nil {
		t.Fatalf("ReadANTLR: %v", err)
	}
	rows, uncovered := g.Rows()
	if len(rows) == 0 {
		t.Fatal("expected generated rows, got 0")
	}
	if len(uncovered) > 0 {
		t.Errorf("uncovered alternatives (%d): %v", len(uncovered), uncovered)
	}
}
