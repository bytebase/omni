// Command cup2g4 converts the grammar of Doris's legacy parser, CUP's
// sql_parser.cup with the JFlex scanner sql_scanner.flex (doris/testdata, as
// the pinned release publishes them), into an ANTLR 4 parser grammar the
// compatibility rows are generated from: doris/testdata/DorisLegacyParser.g4.
// See convert for what the conversion keeps and drops.
//
// Doris's ANTLR parser (DorisParser.g4) reads every statement first, and
// hands the statements its unsupportedStatement rule derives, and a LOAD or
// CREATE ... POLICY, to the CUP parser, which reads them again: Doris takes
// such a statement only if both parsers do. The rows of the converted
// grammar sample what the CUP parser takes.
//
// After a pin bump replaces the two files, from the repo root:
//
//	go run ./doris/parser/cmd/cup2g4
//
// TestDorisLegacyParserIsTheConversion fails while the .g4 is anything else.
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// testdata is where the grammars are, from the repo root.
const testdata = "doris/testdata"

const (
	cupFile  = "sql_parser.cup"
	flexFile = "sql_scanner.flex"
	output   = "DorisLegacyParser.g4"
)

func main() {
	text, err := convertDir(testdata)
	if err != nil {
		fmt.Fprintln(os.Stderr, "cup2g4:", err)
		os.Exit(1)
	}
	if err := os.WriteFile(filepath.Join(testdata, output), []byte(text), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, "cup2g4:", err)
		os.Exit(1)
	}
}

// convertDir converts the vendor files in dir.
func convertDir(dir string) (string, error) {
	cup, err := os.ReadFile(filepath.Join(dir, cupFile))
	if err != nil {
		return "", err
	}
	flex, err := os.ReadFile(filepath.Join(dir, flexFile))
	if err != nil {
		return "", err
	}
	return convert(header, string(cup), string(flex))
}

// header opens the .g4: the notice of the files it is derived from, and what
// was changed, as the Apache License asks of a modified file.
const header = `// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

// This file is derived from Apache Doris's fe/fe-core/src/main/cup/sql_parser.cup
// and fe/fe-core/src/main/jflex/sql_scanner.flex, which are beside it as
// published, with the license they come under (Doris.LICENSE). It was
// changed: doris/parser/cmd/cup2g4 wrote it from them, converting CUP to
// ANTLR 4 and keeping only what each rule derives. Actions, labels and
// %prec are gone, and so are the rules the start symbol does not reach. A
// keyword is written as the text sql_scanner.flex maps to it, punctuation
// as the text its scanner rule matches, and a class of texts as the
// DorisLexer.g4 token that matches the same texts, or a token of its own
// that tokens.tsv samples. Each rule is under the line of sql_parser.cup
// it came from.
//
// Do not edit: rerun the converter.

`
