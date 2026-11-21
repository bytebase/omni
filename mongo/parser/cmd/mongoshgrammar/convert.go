package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"slices"
	"strings"
)

// vendorFiles lists the exact files from mongodb-js/mongosh read by the
// converter, in deterministic sorted order. The SHA-256 of their bytes
// concatenated in this order is the digest pinned in mongo/testdata/pins.tsv.
var vendorFiles = []string{
	"packages/shell-api/src/abstract-cursor.ts",
	"packages/shell-api/src/aggregate-or-find-cursor.ts",
	"packages/shell-api/src/aggregation-cursor.ts",
	"packages/shell-api/src/bulk.ts",
	"packages/shell-api/src/change-stream-cursor.ts",
	"packages/shell-api/src/collection.ts",
	"packages/shell-api/src/cursor.ts",
	"packages/shell-api/src/database.ts",
	"packages/shell-api/src/explainable-cursor.ts",
	"packages/shell-api/src/explainable.ts",
	"packages/shell-api/src/field-level-encryption.ts",
	"packages/shell-api/src/mongo.ts",
	"packages/shell-api/src/plan-cache.ts",
	"packages/shell-api/src/replica-set.ts",
	"packages/shell-api/src/run-command-cursor.ts",
	"packages/shell-api/src/session.ts",
	"packages/shell-api/src/shard.ts",
	"packages/shell-api/src/shell-api.ts",
	"packages/shell-api/src/stream-processor.ts",
	"packages/shell-api/src/streams.ts",
	"packages/shell-bson/src/shell-bson.ts",
}

// methodDecl is one public method extracted from a mongosh TypeScript class.
type methodDecl struct {
	name   string
	file   string
	line   int
	direct bool // decorated with @directShellCommand
}

// staticHelper is one BSON static method (e.g., Binary.createFromBase64).
type staticHelper struct {
	owner  string
	method string
	file   string
	line   int
}

// bsonHelper is one callable BSON constructor or helper from shell-bson.ts.
type bsonHelper struct {
	name string
	file string
	line int
}

// showTarget is one `show <target>` subcommand extracted from Mongo.show().
type showTarget struct {
	name string
	file string
	line int
}

// extracted holds every API surface element parsed from the vendor files.
type extracted struct {
	classes      map[string][]methodDecl
	showTargets  []showTarget
	bsonHelpers  []bsonHelper
	staticHelper []staticHelper
	ejsonMethods []staticHelper
}

// filesDigest returns "sha256:<hex>" over the vendor files in vendorFiles order.
func filesDigest(sources map[string][]byte) (string, error) {
	h := sha256.New()
	for _, path := range vendorFiles {
		b, ok := sources[path]
		if !ok {
			return "", fmt.Errorf("missing vendor file %s", path)
		}
		h.Write(b)
	}
	return "sha256:" + hex.EncodeToString(h.Sum(nil)), nil
}

var (
	reClassDecl  = regexp.MustCompile(`^export\s+(?:default\s+)?(?:abstract\s+)?class\s+([A-Za-z_]\w*)`)
	reMethodDecl = regexp.MustCompile(`^  (public\s+)?(private\s+|protected\s+|static\s+|get\s+|set\s+)?(async\s+)?([A-Za-z]\w*)\s*(?:<|\()`)
	reShowCase   = regexp.MustCompile(`^\s*case\s+'([A-Za-z]\w*)':`)
	reBsonProp   = regexp.MustCompile(`^  ([A-Za-z]\w*):\s`)
	rePickStart  = regexp.MustCompile(`pickWithExactKeyMatch\(\s*bson\.([A-Za-z]\w*)`)
	reSingleQuot = regexp.MustCompile(`'([A-Za-z_]\w*)'`)
)

// ignoredMethodNames are control-flow or lifecycle keywords that can appear at
// 2-space indent inside a TypeScript class body.
var ignoredMethodNames = map[string]bool{
	"constructor": true,
	"if":          true,
	"for":         true,
	"while":       true,
	"switch":      true,
	"return":      true,
	"catch":       true,
	"try":         true,
	"else":        true,
}

// extract parses the vendor TypeScript sources and collects all public shell
// API methods, direct shell commands, show targets, and BSON helpers.
func extract(sources map[string][]byte) (*extracted, error) {
	out := &extracted{
		classes: map[string][]methodDecl{},
	}

	for _, path := range vendorFiles {
		b, ok := sources[path]
		if !ok {
			return nil, fmt.Errorf("missing vendor file %s", path)
		}
		text := string(b)
		if strings.HasSuffix(path, "shell-bson.ts") {
			extractShellBson(path, text, out)
			continue
		}
		extractClassMethods(path, text, out)
		if strings.HasSuffix(path, "/mongo.ts") {
			extractShowTargets(path, text, out)
		}
	}

	return out, nil
}

// extractClassMethods scans TypeScript exported classes and extracts their
// public methods along with exact file:line citations and @directShellCommand
// decorators.
func extractClassMethods(path, src string, out *extracted) {
	lines := strings.Split(src, "\n")
	var curClass string
	seen := map[string]bool{}
	pendingDirect := false

	for i, raw := range lines {
		lineNo := i + 1
		trimmed := strings.TrimSpace(raw)

		if m := reClassDecl.FindStringSubmatch(raw); m != nil {
			curClass = m[1]
			seen = map[string]bool{}
			pendingDirect = false
			continue
		}
		if curClass == "" {
			continue
		}
		if strings.HasPrefix(raw, "}") {
			curClass = ""
			pendingDirect = false
			continue
		}
		if trimmed == "@directShellCommand" {
			pendingDirect = true
			continue
		}
		if strings.HasPrefix(trimmed, "@") {
			// Another decorator on the same method: preserve pendingDirect.
			continue
		}

		if m := reMethodDecl.FindStringSubmatch(raw); m != nil {
			vis, name := m[2], m[4]
			direct := pendingDirect
			pendingDirect = false
			if vis != "" || ignoredMethodNames[name] || strings.HasPrefix(name, "_") {
				continue
			}
			if seen[name] {
				continue
			}
			seen[name] = true
			out.classes[curClass] = append(out.classes[curClass], methodDecl{
				name:   name,
				file:   path,
				line:   lineNo,
				direct: direct,
			})
			continue
		}
		if trimmed != "" && !strings.HasPrefix(trimmed, "//") && !strings.HasPrefix(trimmed, "*") && !strings.HasPrefix(trimmed, "/*") {
			pendingDirect = false
		}
	}
}

// extractShowTargets extracts the `case '<target>':` branches from Mongo.show()
// in packages/shell-api/src/mongo.ts.
func extractShowTargets(path, src string, out *extracted) {
	lines := strings.Split(src, "\n")
	inShow := false
	inSwitch := false
	seen := map[string]bool{}

	for i, raw := range lines {
		lineNo := i + 1
		if !inShow {
			if strings.Contains(raw, "async show(") {
				inShow = true
			}
			continue
		}
		if !inSwitch {
			if strings.Contains(raw, "switch (cmd)") {
				inSwitch = true
			}
			continue
		}
		if strings.HasPrefix(raw, "  }") {
			break
		}
		if m := reShowCase.FindStringSubmatch(raw); m != nil {
			name := m[1]
			if !seen[name] {
				seen[name] = true
				out.showTargets = append(out.showTargets, showTarget{
					name: name,
					file: path,
					line: lineNo,
				})
			}
		}
	}
}

// extractShellBson extracts callable BSON constructors/helpers from
// ShellBsonBase and static methods from pickWithExactKeyMatch(bson.<Class>, [...])
// in packages/shell-bson/src/shell-bson.ts.
func extractShellBson(path, src string, out *extracted) {
	lines := strings.Split(src, "\n")
	inInterface := false
	seenHelper := map[string]bool{}

	for i, raw := range lines {
		lineNo := i + 1
		if !inInterface {
			if strings.HasPrefix(raw, "export interface ShellBsonBase") {
				inInterface = true
			}
			continue
		}
		if strings.HasPrefix(raw, "}") {
			break
		}
		if m := reBsonProp.FindStringSubmatch(raw); m != nil {
			name := m[1]
			if name == "EJSON" || seenHelper[name] {
				// EJSON is a namespace object (EJSON.parse, etc.), not a direct constructor.
				continue
			}
			seenHelper[name] = true
			out.bsonHelpers = append(out.bsonHelpers, bsonHelper{
				name: name,
				file: path,
				line: lineNo,
			})
		}
	}

	// Extract static methods from pickWithExactKeyMatch(bson.<Owner>, [ ... ]).
	inPick := false
	curOwner := ""
	curLine := 0
	seenStatic := map[string]bool{}
	for i, raw := range lines {
		lineNo := i + 1
		if !inPick {
			if m := rePickStart.FindStringSubmatch(raw); m != nil {
				curOwner = m[1]
				curLine = lineNo
				inPick = true
				// The list may be on the same line.
				collectStaticPicks(raw[strings.Index(raw, "["):], curOwner, path, curLine, seenStatic, out)
				if strings.Contains(raw, "]") {
					inPick = false
				}
			}
			continue
		}
		collectStaticPicks(raw, curOwner, path, lineNo, seenStatic, out)
		if strings.Contains(raw, "]") {
			inPick = false
		}
	}
}

func collectStaticPicks(segment, owner, path string, lineNo int, seen map[string]bool, out *extracted) {
	// Skip non-callable properties (prototype, cacheHexString, ALL_CAPS constants).
	for _, m := range reSingleQuot.FindAllStringSubmatch(segment, -1) {
		name := m[1]
		if name == "prototype" || name == "cacheHexString" || isAllCaps(name) {
			continue
		}
		key := owner + "." + name
		if seen[key] {
			continue
		}
		seen[key] = true
		item := staticHelper{owner: owner, method: name, file: path, line: lineNo}
		if owner == "EJSON" {
			out.ejsonMethods = append(out.ejsonMethods, item)
		} else {
			out.staticHelper = append(out.staticHelper, item)
		}
	}
}

func isAllCaps(s string) bool {
	hasLetter := false
	for _, r := range s {
		if r >= 'a' && r <= 'z' {
			return false
		}
		if r >= 'A' && r <= 'Z' {
			hasLetter = true
		}
	}
	return hasLetter
}

// renderGrammar emits mongo/testdata/MongoShell.g4 from the extracted mongosh
// API definitions. Every rule cites the exact vendor file and line it came from.
func renderGrammar(tag, commit, digest string, ex *extracted) (string, error) {
	// Validate required classes were extracted so a changed upstream layout
	// fails loudly instead of silently emitting an incomplete grammar.
	for _, cls := range []string{
		"ShellApi", "ShellConfig", "Database", "Collection",
		"BaseCursor", "AbstractFiniteCursor", "AggregateOrFindCursor", "Cursor",
		"ChangeStreamCursor", "ExplainableCursor",
		"Explainable", "Bulk", "BulkFindOp", "PlanCache",
		"ReplicaSet", "Shard", "Streams", "StreamProcessor",
		"Mongo", "ClientEncryption", "KeyVault", "Session",
	} {
		if len(ex.classes[cls]) == 0 {
			return "", fmt.Errorf("no methods extracted for class %s", cls)
		}
	}
	if len(ex.showTargets) == 0 {
		return "", fmt.Errorf("no show targets extracted from Mongo.show")
	}
	if len(ex.bsonHelpers) == 0 {
		return "", fmt.Errorf("no BSON helpers extracted from ShellBsonBase")
	}

	var b strings.Builder
	fmt.Fprintf(&b, `/*
 * Copyright MongoDB, Inc. (mongodb-js/mongosh, Apache-2.0)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// MongoDB Shell (mongosh) call-shape grammar as ANTLR 4, converted from the
// official mongodb-js/mongosh TypeScript shell-api and shell-bson definitions
// by mongo/parser/cmd/mongoshgrammar. Do not edit by hand: rerun the converter.
//
// Scope (§2.2 of docs/specs/2026-11-11-engine-compatibility-without-servers-design.md):
// covers mongosh call shapes (direct shell commands, ShellApi/ShellConfig calls,
// db.*, db.<coll>.*, cursor/explain/bulk/planCache chains, rs.*, sh.*, sp.*,
// Mongo/connect/session/encryption chains, BSON constructors & static helpers,
// and representative document/array/literal argument shapes) rather than
// arbitrary JavaScript syntax.
//
// pin: mongosh github.com/mongodb-js/mongosh@%s %s %s

grammar MongoShell;

statement
    : directShellCommand
    | shellApiCall
    | shellConfigCall
    | dbStatement
    | collectionStatement
    | explainableStatement
    | bulkStatement
    | planCacheStatement
    | rsStatement
    | shStatement
    | spStatement
    | mongoConnectionStatement
    | encryptionStatement
    | sessionStatement
    | ejsonStatement
    ;

`, commit, digest, tag)

	// 1. Direct shell commands (@directShellCommand in ShellApi + Mongo.show).
	b.WriteString("// packages/shell-api/src/shell-api.ts (@directShellCommand) and packages/shell-api/src/mongo.ts (Mongo.show)\n")
	b.WriteString("directShellCommand\n")
	directAlts := []string{}
	for _, m := range ex.classes["ShellApi"] {
		if !m.direct {
			continue
		}
		switch m.name {
		case "use":
			directAlts = append(directAlts, fmt.Sprintf("'use' IDENT // %s:%d", m.file, m.line))
		case "show":
			directAlts = append(directAlts, fmt.Sprintf("'show' showTarget // %s:%d", m.file, m.line))
			directAlts = append(directAlts, fmt.Sprintf("'show' 'log' IDENT // %s:%d (show(cmd, arg?))", m.file, m.line))
		case "exit", "quit":
			directAlts = append(directAlts, fmt.Sprintf("'%s' // %s:%d", m.name, m.file, m.line))
			directAlts = append(directAlts, fmt.Sprintf("'%s' NUMBER // %s:%d (exitCode?: number)", m.name, m.file, m.line))
		default:
			directAlts = append(directAlts, fmt.Sprintf("'%s' // %s:%d", m.name, m.file, m.line))
		}
	}
	writeAlts(&b, directAlts)

	// showTarget rule.
	b.WriteString("// packages/shell-api/src/mongo.ts (Mongo.show)\n")
	b.WriteString("showTarget\n")
	showAlts := make([]string, 0, len(ex.showTargets))
	for _, st := range ex.showTargets {
		showAlts = append(showAlts, fmt.Sprintf("'%s' // %s:%d", st.name, st.file, st.line))
	}
	writeAlts(&b, showAlts)

	// 2. ShellApi global function calls (excluding Mongo/connect which have dedicated connection chains).
	b.WriteString("// packages/shell-api/src/shell-api.ts (ShellApi)\n")
	b.WriteString("shellApiCall\n    : shellApiMethod '(' argList? ')'\n    ;\n\n")
	var shellApiAlts []string
	for _, m := range ex.classes["ShellApi"] {
		if m.name == "Mongo" || m.name == "connect" {
			continue
		}
		shellApiAlts = append(shellApiAlts, fmt.Sprintf("'%s' // %s:%d", m.name, m.file, m.line))
	}
	b.WriteString("shellApiMethod\n")
	writeAlts(&b, shellApiAlts)

	// 3. ShellConfig calls: config.set / config.get / config.reset.
	b.WriteString("// packages/shell-api/src/shell-api.ts (ShellConfig)\n")
	b.WriteString("shellConfigCall\n    : 'config' '.' shellConfigMethod '(' argList? ')'\n    ;\n\n")
	writeMethodRule(&b, "shellConfigMethod", ex.classes["ShellConfig"])

	// 4. Database statements: db.<method>(), db.getSiblingDB('s').<method>(),
	// db.getSiblingDB('s').<coll>.<method>(), db.aggregate().<cursor>(),
	// db.watch().<changeStream>(), db.getMongo().<mongoMethod>().
	b.WriteString("// packages/shell-api/src/database.ts (Database)\n")
	b.WriteString(`dbStatement
    : 'db' '.' databaseMethod '(' argList? ')'
    | 'db' '.' 'getSiblingDB' '(' STRING ')' '.' databaseMethod '(' argList? ')'
    | 'db' '.' 'getSiblingDB' '(' STRING ')' '.' IDENT '.' collectionMethod '(' argList? ')'
    | 'db' '.' 'aggregate' '(' argList? ')' '.' cursorMethod '(' argList? ')'
    | 'db' '.' 'watch' '(' argList? ')' '.' changeStreamMethod '(' argList? ')'
    | 'db' '.' 'getMongo' '(' ')' '.' mongoMethod '(' argList? ')'
    ;

`)
	writeMethodRule(&b, "databaseMethod", ex.classes["Database"])

	// 5. Collection statements + cursor chaining.
	b.WriteString("// packages/shell-api/src/collection.ts (Collection) and cursor classes\n")
	b.WriteString(`collectionRef
    : 'db' '.' IDENT
    | 'db' '[' STRING ']'
    | 'db' '.' 'getCollection' '(' STRING ')'
    ;

collectionStatement
    : collectionRef '.' collectionMethod '(' argList? ')' cursorChain?
    ;

cursorChain
    : '.' cursorMethod '(' argList? ')'
    ;

`)
	writeMethodRule(&b, "collectionMethod", ex.classes["Collection"])

	// Deduplicated union of cursor methods across BaseCursor, AbstractFiniteCursor,
	// AggregateOrFindCursor, Cursor, ChangeStreamCursor, ExplainableCursor.
	cursorMethods := mergeMethods(
		ex.classes["BaseCursor"],
		ex.classes["AbstractFiniteCursor"],
		ex.classes["AggregateOrFindCursor"],
		ex.classes["Cursor"],
		ex.classes["ChangeStreamCursor"],
		ex.classes["ExplainableCursor"],
	)
	b.WriteString("// packages/shell-api/src/{abstract-cursor,aggregate-or-find-cursor,cursor,change-stream-cursor,explainable-cursor}.ts\n")
	writeMethodRule(&b, "cursorMethod", cursorMethods)

	// ChangeStreamCursor-specific methods (also included in cursorMethod, plus BaseCursor methods).
	b.WriteString("// packages/shell-api/src/change-stream-cursor.ts (ChangeStreamCursor)\n")
	writeMethodRule(&b, "changeStreamMethod", ex.classes["ChangeStreamCursor"])

	// 6. Explainable statements.
	b.WriteString("// packages/shell-api/src/explainable.ts (Explainable)\n")
	b.WriteString(`explainableStatement
    : collectionRef '.' 'explain' '(' argList? ')' ( '.' explainableMethod '(' argList? ')' )?
    ;

`)
	writeMethodRule(&b, "explainableMethod", ex.classes["Explainable"])

	// 7. Bulk statements.
	b.WriteString("// packages/shell-api/src/bulk.ts (Bulk, BulkFindOp)\n")
	b.WriteString(`bulkStatement
    : collectionRef '.' ( 'initializeOrderedBulkOp' | 'initializeUnorderedBulkOp' ) '(' ')' '.' bulkMethod '(' argList? ')'
    ;

`)
	writeMethodRule(&b, "bulkMethod", mergeMethods(ex.classes["Bulk"], ex.classes["BulkFindOp"]))

	// 8. PlanCache statements.
	b.WriteString("// packages/shell-api/src/plan-cache.ts (PlanCache)\n")
	b.WriteString(`planCacheStatement
    : collectionRef '.' 'getPlanCache' '(' ')' ( '.' planCacheMethod '(' argList? ')' )?
    ;

`)
	writeMethodRule(&b, "planCacheMethod", ex.classes["PlanCache"])

	// 9. ReplicaSet statements.
	b.WriteString("// packages/shell-api/src/replica-set.ts (ReplicaSet)\n")
	b.WriteString("rsStatement\n    : 'rs' '.' rsMethod '(' argList? ')'\n    ;\n\n")
	writeMethodRule(&b, "rsMethod", ex.classes["ReplicaSet"])

	// 10. Shard statements.
	b.WriteString("// packages/shell-api/src/shard.ts (Shard)\n")
	b.WriteString("shStatement\n    : 'sh' '.' shMethod '(' argList? ')'\n    ;\n\n")
	writeMethodRule(&b, "shMethod", ex.classes["Shard"])

	// 11. Atlas Stream Processing statements (Streams & StreamProcessor).
	b.WriteString("// packages/shell-api/src/streams.ts (Streams) and packages/shell-api/src/stream-processor.ts (StreamProcessor)\n")
	b.WriteString(`spStatement
    : 'sp' '.' streamsMethod '(' argList? ')'
    | 'sp' '.' IDENT '.' streamProcessorMethod '(' argList? ')'
    ;

`)
	writeMethodRule(&b, "streamsMethod", ex.classes["Streams"])
	writeMethodRule(&b, "streamProcessorMethod", ex.classes["StreamProcessor"])

	// 12. Connection statements: Mongo(...) and connect(...).
	b.WriteString("// packages/shell-api/src/shell-api.ts (Mongo, connect) and packages/shell-api/src/mongo.ts (Mongo)\n")
	b.WriteString(`mongoConnectionStatement
    : 'Mongo' '(' argList? ')' ( '.' mongoMethod '(' argList? ')' )?
    | 'connect' '(' argList? ')' ( '.' databaseMethod '(' argList? ')' )?
    ;

`)
	writeMethodRule(&b, "mongoMethod", ex.classes["Mongo"])

	// 13. Client-side field-level encryption (KeyVault, ClientEncryption).
	b.WriteString("// packages/shell-api/src/field-level-encryption.ts (KeyVault, ClientEncryption)\n")
	b.WriteString(`encryptionStatement
    : 'db' '.' 'getMongo' '(' ')' '.' 'getKeyVault' '(' ')' ( '.' keyVaultMethod '(' argList? ')' )?
    | 'db' '.' 'getMongo' '(' ')' '.' 'getClientEncryption' '(' ')' ( '.' clientEncryptionMethod '(' argList? ')' )?
    ;

`)
	writeMethodRule(&b, "keyVaultMethod", ex.classes["KeyVault"])
	writeMethodRule(&b, "clientEncryptionMethod", ex.classes["ClientEncryption"])

	// 14. Session statements.
	b.WriteString("// packages/shell-api/src/session.ts (Session)\n")
	b.WriteString(`sessionStatement
    : 'db' '.' 'getMongo' '(' ')' '.' 'startSession' '(' argList? ')' ( '.' sessionMethod '(' argList? ')' )?
    ;

`)
	writeMethodRule(&b, "sessionMethod", ex.classes["Session"])

	// 15. EJSON top-level helper calls.
	b.WriteString("// packages/shell-bson/src/shell-bson.ts (EJSON)\n")
	b.WriteString("ejsonStatement\n    : 'EJSON' '.' ejsonMethod '(' argList? ')'\n    ;\n\n")
	ejsonAlts := make([]string, 0, len(ex.ejsonMethods))
	for _, em := range ex.ejsonMethods {
		ejsonAlts = append(ejsonAlts, fmt.Sprintf("'%s' // %s:%d", em.method, em.file, em.line))
	}
	b.WriteString("ejsonMethod\n")
	writeAlts(&b, ejsonAlts)

	// 16. Representative argument, document, array, literal, and BSON helper shapes.
	b.WriteString(`// Representative argument, document, array, and literal value shapes
argList
    : value ( ',' value )* ','?
    ;

value
    : document
    | array
    | STRING
    | NUMBER
    | NEG_NUMBER
    | 'true'
    | 'false'
    | 'null'
    | REGEX
    | bsonHelperCall
    | bsonStaticCall
    | IDENT
    ;

document
    : '{' '}'
    | '{' pair ( ',' pair )* ','? '}'
    ;

pair
    : docKey ':' value
    ;

docKey
    : IDENT
    | STRING
    | FIELD_OP
    ;

array
    : '[' ']'
    | '[' value ( ',' value )* ','? ']'
    ;

`)

	// BSON helpers from ShellBsonBase.
	b.WriteString("// packages/shell-bson/src/shell-bson.ts (ShellBsonBase)\n")
	b.WriteString("bsonHelperCall\n    : bsonHelperName '(' argList? ')'\n    ;\n\n")
	bsonAlts := make([]string, 0, len(ex.bsonHelpers))
	for _, h := range ex.bsonHelpers {
		bsonAlts = append(bsonAlts, fmt.Sprintf("'%s' // %s:%d", h.name, h.file, h.line))
	}
	b.WriteString("bsonHelperName\n")
	writeAlts(&b, bsonAlts)

	// BSON static helper calls (e.g., Binary.createFromBase64(...), ObjectId.createFromHexString(...)).
	b.WriteString("// packages/shell-bson/src/shell-bson.ts (static BSON helpers)\n")
	b.WriteString("bsonStaticCall\n    : bsonStaticTarget '(' argList? ')'\n    ;\n\n")
	staticAlts := make([]string, 0, len(ex.staticHelper)+len(ex.ejsonMethods))
	for _, sh := range ex.staticHelper {
		staticAlts = append(staticAlts, fmt.Sprintf("'%s' '.' '%s' // %s:%d", sh.owner, sh.method, sh.file, sh.line))
	}
	for _, em := range ex.ejsonMethods {
		staticAlts = append(staticAlts, fmt.Sprintf("'%s' '.' '%s' // %s:%d", em.owner, em.method, em.file, em.line))
	}
	b.WriteString("bsonStaticTarget\n")
	writeAlts(&b, staticAlts)

	return b.String(), nil
}

func mergeMethods(lists ...[]methodDecl) []methodDecl {
	seen := map[string]bool{}
	var out []methodDecl
	for _, list := range lists {
		for _, m := range list {
			if seen[m.name] {
				continue
			}
			seen[m.name] = true
			out = append(out, m)
		}
	}
	return out
}

func writeMethodRule(b *strings.Builder, ruleName string, methods []methodDecl) {
	b.WriteString(ruleName + "\n")
	alts := make([]string, 0, len(methods))
	for _, m := range methods {
		alts = append(alts, fmt.Sprintf("'%s' // %s:%d", m.name, m.file, m.line))
	}
	writeAlts(b, alts)
}

func writeAlts(b *strings.Builder, alts []string) {
	for i, alt := range alts {
		if i == 0 {
			b.WriteString("    : " + alt + "\n")
		} else {
			b.WriteString("    | " + alt + "\n")
		}
	}
	b.WriteString("    ;\n\n")
}

// updatePinRow replaces or appends the "mongosh" role row in a pins.tsv file.
func updatePinRow(pinsText, commit, digest, tag string) string {
	row := fmt.Sprintf("mongosh\tgithub.com/mongodb-js/mongosh@%s\t%s\t%s", commit, digest, tag)
	lines := strings.Split(strings.TrimRight(pinsText, "\n"), "\n")
	replaced := false
	for i, line := range lines {
		if strings.HasPrefix(line, "mongosh\t") {
			lines[i] = row
			replaced = true
			break
		}
	}
	if !replaced {
		lines = append(lines, row)
	}
	// Ensure trailing newline.
	return strings.Join(slices.Clone(lines), "\n") + "\n"
}
