# pg test data

What `pg`'s committed engine answers are judged by and generated from. The
design is `docs/specs/2026-11-11-engine-compatibility-without-servers-design.md`.

- **`pins.tsv`** names each PostgreSQL server by image digest and the source
  tag it was built from. `default` (17.12) judges `../parser/testdata/compat.tsv`;
  `pg18` (18.7) serves the tests that ask the next major.
- **`gram.y`** is `src/backend/parser/gram.y` of
  [postgres/postgres](https://github.com/postgres/postgres) at tag `REL_17_12`
  (commit `702b3dc86d35bcc08c8e23cffa3c4d9fe9d32c85`), unchanged: SHA-256
  `ba8ad5a6e1df5768dee7dbca5f2594400a741806b4f671579ecf806969b9e3bf`. The
  generated rows of `compat.tsv` are its alternatives and their near-misses.
  A pin bump replaces it with the new tag's.
- **`tokens.tsv`** gives the text each terminal of `gram.y` is written as.
  PostgreSQL's lexer is flex (`src/backend/parser/scan.l`), not a grammar the
  generator reads, so the file is written by hand and cites `scan.l`, the
  filter between the lexer and the grammar (`base_yylex` in
  `src/backend/parser/parser.c`) and the keyword table
  (`src/include/parser/kwlist.h`), all at `REL_17_12`. A new release's
  keywords need rows here.

## License

`gram.y` is distributed under the PostgreSQL License, which asks that this
notice appear in all copies:

> PostgreSQL Database Management System
> (also known as Postgres, formerly known as Postgres95)
>
> Portions Copyright (c) 1996-2026, PostgreSQL Global Development Group
>
> Portions Copyright (c) 1994, The Regents of the University of California
>
> Permission to use, copy, modify, and distribute this software and its
> documentation for any purpose, without fee, and without a written agreement
> is hereby granted, provided that the above copyright notice and this
> paragraph and the following two paragraphs appear in all copies.
>
> IN NO EVENT SHALL THE UNIVERSITY OF CALIFORNIA BE LIABLE TO ANY PARTY FOR
> DIRECT, INDIRECT, SPECIAL, INCIDENTAL, OR CONSEQUENTIAL DAMAGES, INCLUDING
> LOST PROFITS, ARISING OUT OF THE USE OF THIS SOFTWARE AND ITS
> DOCUMENTATION, EVEN IF THE UNIVERSITY OF CALIFORNIA HAS BEEN ADVISED OF THE
> POSSIBILITY OF SUCH DAMAGE.
>
> THE UNIVERSITY OF CALIFORNIA SPECIFICALLY DISCLAIMS ANY WARRANTIES,
> INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
> AND FITNESS FOR A PARTICULAR PURPOSE. THE SOFTWARE PROVIDED HEREUNDER IS
> ON AN "AS IS" BASIS, AND THE UNIVERSITY OF CALIFORNIA HAS NO OBLIGATIONS TO
> PROVIDE MAINTENANCE, SUPPORT, UPDATES, ENHANCEMENTS, OR MODIFICATIONS.

The file is a test fixture read at test time and is not compiled into omni.
