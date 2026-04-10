# Official Snowflake SQL Reference Examples

These SQL files are scraped from the official Snowflake SQL reference at https://docs.snowflake.com/en/sql-reference-commands and its subpages. They form **Corpus B** of the omni snowflake parser test suite.

Each subdirectory corresponds to one command (kebab-case slug from the docs URL). Each `example_NN.sql` is a single SQL example extracted verbatim from the corresponding command's "Examples" section.

The official docs are the **authoritative source of truth** for Snowflake syntax. When the legacy ANTLR4 grammar disagrees with the docs, the docs win.

Every file in this directory MUST eventually parse cleanly under `omni/snowflake/parser`. Files that fail to parse indicate either a parser bug or a docs corner case worth investigating.

## Resync

To refresh this corpus, re-scrape `https://docs.snowflake.com/en/sql-reference-commands` and overwrite the contents of this directory.

## File count

630 SQL files across 78 command directories.
Scrape date: 2026-04-07.
