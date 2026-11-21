package parser

import "regexp"

// grammarRefusals match the messages GoogleSQL's grammar writes when one of
// its parser actions refuses input the rules derive, without the "Syntax error:"
// prefix (mirrors googlesql/spanner-oracle/grammar_refusals.go, generated from
// googlesql/testdata/googlesql.tm, plus the lookahead-transformer refusal on
// bare EXCEPT and the MergeDynamicLabelProperties action refusal on duplicate
// DYNAMIC LABEL/PROPERTIES). Each pattern accepts both "[at 1:N]" and
// BigQuery's "at [1:N]" position suffix.
var grammarRefusals = []*regexp.Regexp{
	regexp.MustCompile(`^EXCEPT must be followed by ALL, DISTINCT, or "\("(?:$|\s+\[at |\s+at \[)`),                                    // lookahead_transformer
	regexp.MustCompile(`^DYNAMIC (?:LABEL|PROPERTIES) cannot be used more than once in a single element table(?:$|\s+\[at |\s+at \[)`), // googlesql.tm:2862 (MergeDynamicLabelProperties)
	regexp.MustCompile(`^Hints are not allowed on DEFINE MACRO statements\.(?:$|\s+\[at |\s+at \[)`),                                   // googlesql.tm:1195
	regexp.MustCompile(`^Hints are not allowed on DEFINE MACRO statements(?:$|\s+\[at |\s+at \[)`),                                     // googlesql.tm:1206
	regexp.MustCompile(`^Mismatched end label; expected .*, got .*`),                                                                   // googlesql.tm:1244
	regexp.MustCompile(`^Pipe operators are not supported on this statement(?:$|\s+\[at |\s+at \[)`),                                   // googlesql.tm:1269
	regexp.MustCompile(`^RUN statement is not supported\.(?:$|\s+\[at |\s+at \[)`),                                                     // googlesql.tm:1375
	regexp.MustCompile(`^ADD ROW DELETION POLICY clause is not supported\.(?:$|\s+\[at |\s+at \[)`),                                    // googlesql.tm:1587
	regexp.MustCompile(`^REPLACE ROW DELETION POLICY clause is not supported\.(?:$|\s+\[at |\s+at \[)`),                                // googlesql.tm:1597
	regexp.MustCompile(`^DROP ROW DELETION POLICY clause is not supported\.(?:$|\s+\[at |\s+at \[)`),                                   // googlesql.tm:1607
	regexp.MustCompile(`^EXTERNAL TABLE FUNCTION is not supported(?:$|\s+\[at |\s+at \[)`),                                             // googlesql.tm:1800
	regexp.MustCompile(`^ALTER TABLE FUNCTION is not supported(?:$|\s+\[at |\s+at \[)`),                                                // googlesql.tm:1852
	regexp.MustCompile(`^ALTER .* is not supported(?:$|\s+\[at |\s+at \[)`),                                                            // googlesql.tm:1883
	regexp.MustCompile(`^Parenthesized SET statement requires a variable list(?:$|\s+\[at |\s+at \[)`),                                 // googlesql.tm:2041
	regexp.MustCompile(`^Using SET with multiple variables requires parentheses around the variable list(?:$|\s+\[at |\s+at \[)`),      // googlesql.tm:2048
	regexp.MustCompile(`^Keyword REMOTE is not supported(?:$|\s+\[at |\s+at \[)`),                                                      // googlesql.tm:2313
	regexp.MustCompile(`^WITH CONNECTION clause is not supported(?:$|\s+\[at |\s+at \[)`),                                              // googlesql.tm:2322
	regexp.MustCompile(`^The body of each CREATE FUNCTION statement is an expression, not a query; to use a query as an expression, the query must be wrapped with additional parentheses to make it a scalar subquery expression(?:$|\s+\[at |\s+at \[)`), // googlesql.tm:2420
	regexp.MustCompile(`^CREATE LOCALITY GROUP is not supported(?:$|\s+\[at |\s+at \[)`),                                             // googlesql.tm:2905
	regexp.MustCompile(`^UNDROP .* is not supported(?:$|\s+\[at |\s+at \[)`),                                                         // googlesql.tm:2939
	regexp.MustCompile(`^CREATE SNAPSHOT .* is not supported(?:$|\s+\[at |\s+at \[)`),                                                // googlesql.tm:2967
	regexp.MustCompile(`^LOAD DATA statement with PARTITIONS is not supported(?:$|\s+\[at |\s+at \[)`),                               // googlesql.tm:3104
	regexp.MustCompile(`^LOAD DATA statement with TEMP TABLE is not supported(?:$|\s+\[at |\s+at \[)`),                               // googlesql.tm:3156
	regexp.MustCompile(`^.* is not a supported object type(?:$|\s+\[at |\s+at \[)`),                                                  // googlesql.tm:3184
	regexp.MustCompile(`^.* is not a supported nested object type(?:$|\s+\[at |\s+at \[)`),                                           // googlesql.tm:3213
	regexp.MustCompile(`^A table must define at least one column\.(?:$|\s+\[at |\s+at \[)`),                                          // googlesql.tm:3295
	regexp.MustCompile(`^Column DEFAULT value is not supported\.(?:$|\s+\[at |\s+at \[)`),                                            // googlesql.tm:3643
	regexp.MustCompile(`^Internal Error: Expected foreign key reference(?:$|\s+\[at |\s+at \[)`),                                     // googlesql.tm:3768
	regexp.MustCompile(`^Ordering for primary keys is not supported(?:$|\s+\[at |\s+at \[)`),                                         // googlesql.tm:3856
	regexp.MustCompile(`^EXPORT TABLE FUNCTION METADATA is not supported(?:$|\s+\[at |\s+at \[)`),                                    // googlesql.tm:4302
	regexp.MustCompile(`^PIVOT clause requires at least one pivot expression(?:$|\s+\[at |\s+at \[)`),                                // googlesql.tm:6246
	regexp.MustCompile(`^QUALIFY clause must be used in conjunction with WHERE or GROUP BY or HAVING clause(?:$|\s+\[at |\s+at \[)`), // googlesql.tm:6378
	regexp.MustCompile(`^Hint cannot be used in between two .*s(?:$|\s+\[at |\s+at \[)`),                                             // googlesql.tm:7468
	regexp.MustCompile(`^Unable to parse hints in the right order(?:$|\s+\[at |\s+at \[)`),                                           // googlesql.tm:7495
	regexp.MustCompile(`^Hint cannot be used at beginning of path pattern(?:$|\s+\[at |\s+at \[)`),                                   // googlesql.tm:7532
	regexp.MustCompile(`^Quantifier cannot be used on a node pattern(?:$|\s+\[at |\s+at \[)`),                                        // googlesql.tm:7579
	regexp.MustCompile(`^WHERE clause cannot be used together with property specification(?:$|\s+\[at |\s+at \[)`),                   // googlesql.tm:7635
	regexp.MustCompile(`^Query parameters cannot be used in place of table names(?:$|\s+\[at |\s+at \[)`),                            // googlesql.tm:8092
	regexp.MustCompile(`^System variables cannot be used in place of table names(?:$|\s+\[at |\s+at \[)`),                            // googlesql.tm:8102
	regexp.MustCompile(`^QUALIFY is not supported(?:$|\s+\[at |\s+at \[)`),                                                           // googlesql.tm:8461
	regexp.MustCompile(`^FOR UPDATE is not supported(?:$|\s+\[at |\s+at \[)`),                                                        // googlesql.tm:8523
	regexp.MustCompile(`^GROUP ROWS is not supported\.(?:$|\s+\[at |\s+at \[)`),                                                      // googlesql.tm:8685
	regexp.MustCompile(`^The argument to UNNEST is an expression, not a query; to use a query as an expression, the query must be wrapped with additional parentheses to make it a scalar subquery expression(?:$|\s+\[at |\s+at \[)`), // googlesql.tm:8941
	regexp.MustCompile(`^IS DISTINCT FROM is not supported(?:$|\s+\[at |\s+at \[)`),     // googlesql.tm:9472
	regexp.MustCompile(`^Braced constructors are not supported(?:$|\s+\[at |\s+at \[)`), // googlesql.tm:10576
	regexp.MustCompile(`^The argument to CAST is an expression, not a query; to use a query as an expression, the query must be wrapped with additional parentheses to make it a scalar subquery expression(?:$|\s+\[at |\s+at \[)`),      // googlesql.tm:10668
	regexp.MustCompile(`^The argument to SAFE_CAST is an expression, not a query; to use a query as an expression, the query must be wrapped with additional parentheses to make it a scalar subquery expression(?:$|\s+\[at |\s+at \[)`), // googlesql.tm:10676
	regexp.MustCompile(`^Each function argument is an expression, not a query; to use a query as an expression, the query must be wrapped with additional parentheses to make it a scalar subquery expression(?:$|\s+\[at |\s+at \[)`),    // googlesql.tm:10892
	regexp.MustCompile(`^ROW DELETION POLICY clause is not supported\.(?:$|\s+\[at |\s+at \[)`),                                                                                                                                           // googlesql.tm:11103
	regexp.MustCompile(`^DROP AGGREGATE FUNCTION is not supported, use DROP FUNCTION(?:$|\s+\[at |\s+at \[)`),                                                                                                                             // googlesql.tm:12573
	regexp.MustCompile(`^DEFINE MACRO statements cannot be nested under other statements or blocks\.(?:$|\s+\[at |\s+at \[)`),                                                                                                             // googlesql.tm:12652
	regexp.MustCompile(`^Statement CASE\.\.\.WHEN is not supported(?:$|\s+\[at |\s+at \[)`),                                                                                                                                               // googlesql.tm:12813
	regexp.MustCompile(`^REPEAT is not supported(?:$|\s+\[at |\s+at \[)`),                                                                                                                                                                 // googlesql.tm:12903
	regexp.MustCompile(`^FOR\.\.\.IN is not supported(?:$|\s+\[at |\s+at \[)`),                                                                                                                                                            // googlesql.tm:12917
	regexp.MustCompile(`^PRIMARY KEY must be defined in the table element list as column attribute or constraint\.(?:$|\s+\[at |\s+at \[)`),                                                                                               // googlesql.tm:13539
	regexp.MustCompile(`^null_filtered is not a supported object type(?:$|\s+\[at |\s+at \[)`),                                                                                                                                            // googlesql.tm:13603
	regexp.MustCompile(`^Expected keyword DROP or keyword SET but got identifier(?:$|\s+\[at |\s+at \[)`),                                                                                                                                 // googlesql.tm:13645
	regexp.MustCompile(`^Script labels are not supported(?:$|\s+\[at |\s+at \[)`),                                                                                                                                                         // googlesql.tm:13770
}
