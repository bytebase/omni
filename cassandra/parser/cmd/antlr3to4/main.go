// Command antlr3to4 converts Apache Cassandra's CQL grammar, which is ANTLR 3
// (cassandra/testdata/Cql.g, Parser.g and Lexer.g, as the pinned release
// publishes them), into the ANTLR 4 grammar the compatibility rows are
// generated from and judged by, cassandra/testdata/Cql.g4. See convert for
// what the conversion keeps and drops.
//
// After a pin bump replaces the three .g files, from the repo root:
//
//	go run ./cassandra/parser/cmd/antlr3to4
//
// TestCqlG4IsTheConversion fails while Cql.g4 is anything else.
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// testdata is where the grammars are, from the repo root.
const testdata = "cassandra/testdata"

// vendor lists the ANTLR 3 files, the one with the start rule first.
var vendor = []string{"Cql.g", "Parser.g", "Lexer.g"}

const output = "Cql.g4"

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
	return convert("Cql", header, sources...)
}

// header opens Cql.g4: the notice of the files it is derived from, and what
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

// This file is derived from Apache Cassandra's src/antlr/Cql.g, Parser.g and
// Lexer.g, which are beside it as published, with the license and notice
// they come under (Cassandra.LICENSE, Cassandra.NOTICE). It was changed:
// cassandra/parser/cmd/antlr3to4 wrote it from them, converting ANTLR 3 to
// ANTLR 4 and keeping only what each rule derives. Actions, arguments,
// labels and predicates are gone; an alternative whose action refuses every
// input is gone with it; and what else was lost is noted beside its rule.
// Each rule is under the file and line it came from.
//
// Do not edit: rerun the converter.

`
