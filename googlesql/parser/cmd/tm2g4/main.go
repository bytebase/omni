// Command tm2g4 converts GoogleSQL's grammar, the Textmapper file
// googlesql.tm (googlesql/testdata, as the release the pinned Spanner
// emulator builds against publishes it), into the ANTLR 4 grammar the
// compatibility rows are generated from: googlesql/testdata/GoogleSQL.g4.
// See convert for what the conversion keeps and drops.
//
// After a pin bump replaces googlesql.tm, from the repo root:
//
//	go run ./googlesql/parser/cmd/tm2g4
//
// TestGoogleSQLG4IsTheConversion fails while the .g4 is anything else.
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// testdata is where the grammars are, from the repo root.
const testdata = "googlesql/testdata"

const (
	tmFile = "googlesql.tm"
	output = "GoogleSQL.g4"
	// start is the rule the rows are generated from: one statement of a
	// script, a SQL statement or a scripting one, without the ; or end of
	// input next_script_statement reads after it.
	start = "unterminated_statement"
)

func main() {
	text, err := convertDir(testdata)
	if err != nil {
		fmt.Fprintln(os.Stderr, "tm2g4:", err)
		os.Exit(1)
	}
	if err := os.WriteFile(filepath.Join(testdata, output), []byte(text), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, "tm2g4:", err)
		os.Exit(1)
	}
}

// convertDir converts the vendor file in dir.
func convertDir(dir string) (string, error) {
	tm, err := os.ReadFile(filepath.Join(dir, tmFile))
	if err != nil {
		return "", err
	}
	text, _, err := convert(header, string(tm), start)
	return text, err
}

// header opens the .g4: the notice of the file it is derived from, and what
// was changed, as the Apache License asks of a modified file.
const header = `//
// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// This file is derived from GoogleSQL's googlesql/parser/googlesql.tm,
// which is beside it as published, with the license it comes under
// (GoogleSQL.LICENSE). It was changed: googlesql/parser/cmd/tm2g4 wrote it
// from it, converting Textmapper to ANTLR 4 and keeping only what each rule
// derives. Actions, rule types, aliases and %prec are gone, and so are the
// rules unterminated_statement, the first rule, does not reach; a rule
// written name_opt is name?, and (x separator ",")+ is x (',' x)*. A token
// is written as the text it matches, and a token that matches a class of
// texts, or that the tokenizer makes in code, by its id, which tokens.tsv
// samples. Each rule is under the line of googlesql.tm it came from.
//
// Do not edit: rerun the converter.

`
