// Command antlr3to4 converts Apache Hive's HiveQL grammar, which is ANTLR 3
// (hive/testdata/vendor/*.g, as rel/release-4.0.1 publishes them), into the
// ANTLR 4 grammar the compatibility rows are generated from,
// hive/testdata/Hive.g4. See convert for what the conversion keeps and drops.
//
// After a pin bump replaces the vendor .g files, from the repo root:
//
//	go run ./hive/parser/cmd/antlr3to4
//
// TestHiveG4IsTheConversion fails while Hive.g4 is anything else.
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// testdata is where the grammars are, from the repo root.
const testdata = "hive/testdata"

// vendor lists the ANTLR 3 files in order: HiveParser.g (which holds the start
// rule statement) first, then its imported parser grammars, then HiveLexer.g
// (the default lexer subclass) before HiveLexerParent.g so subclass token
// definitions supersede parent stubs.
var vendor = []string{
	"vendor/HiveParser.g",
	"vendor/AlterClauseParser.g",
	"vendor/SelectClauseParser.g",
	"vendor/FromClauseParser.g",
	"vendor/IdentifiersParser.g",
	"vendor/ResourcePlanParser.g",
	"vendor/CreateDDLParser.g",
	"vendor/PrepareStatementParser.g",
	"vendor/HiveLexer.g",
	"vendor/HiveLexerParent.g",
}

const output = "Hive.g4"

func main() {
	text, err := convertDir(testdata)
	if err != nil {
		fmt.Fprintln(os.Stderr, "antlr3to4:", err)
		os.Exit(1)
	}
	if err := os.WriteFile(filepath.Join(testdata, output), []byte(text), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, "antlr3to4:", err)
		os.Exit(1)
	}
}

// convertDir converts the vendor files in dir.
func convertDir(dir string) (string, error) {
	var sources []source
	for _, name := range vendor {
		b, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return "", err
		}
		sources = append(sources, source{name, string(b)})
	}
	return convert("Hive", header, sources...)
}

// header opens Hive.g4: the notice of the files it is derived from, and what
// was changed, as the Apache License asks of a modified file.
const header = `/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

// This file is derived from Apache Hive's parser/src/java/org/apache/hadoop/hive/ql/parse/
// ANTLR 3 grammar files (vendor/*.g), which are beside it as published, with the
// license and notice they come under (Hive.LICENSE, Hive.NOTICE). It was
// changed: hive/parser/cmd/antlr3to4 wrote it from them, converting ANTLR 3 to
// ANTLR 4 and keeping only what each rule derives. AST rewrite rules (-> ...),
// actions, arguments, labels, tree operators, and predicates are gone; empty
// lexer stubs overridden by HiveLexer.g are dropped; and what else was lost is
// noted beside its rule. Each rule is under the file and line it came from.
//
// Do not edit: rerun the converter.

`
