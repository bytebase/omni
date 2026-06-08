package catalog

// extensionScripts contains bundled DDL for known extensions.
// Each entry is the SQL that CREATE EXTENSION <name> will execute through ProcessUtility.
//
// Scripts are derived from the official PostgreSQL extension SQL files,
// simplified to the subset relevant for DDL semantic analysis (types, casts,
// operators, functions, aggregates). Physical/index-only constructs like
// CREATE OPERATOR CLASS are included when the opclass infrastructure is present.
var extensionScripts = map[string]string{
	"btree_gin":     btreeGinSQL,
	"btree_gist":    btreeGistSQL,
	"citext":        citextSQL,
	"cube":          cubeSQL,
	"dblink":        dblinkSQL,
	"earthdistance": earthdistanceSQL,
	"fuzzystrmatch": fuzzystrmatchSQL,
	"hstore":        hstoreSQL,
	"intarray":      intarraySQL,
	"ltree":         ltreeSQL,
	"pg_trgm":       pgTrgmSQL,
	"pgcrypto":      pgcryptoSQL,
	"tablefunc":     tablefuncSQL,
	"unaccent":      unaccentSQL,
	"uuid-ossp":     uuidOSSPSQL,
	"vector":        pgvectorSQL,
}

// citextSQL is the DDL for the citext extension (case-insensitive text type).
// Derived from contrib/citext/citext--1.6.sql.
const citextSQL = `
-- Shell type
CREATE TYPE citext;

-- I/O functions (piggyback on text I/O via LANGUAGE internal)
-- pg: contrib/citext/citext--1.4.sql
CREATE FUNCTION citextin(cstring) RETURNS citext AS 'textin' LANGUAGE internal STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION citextout(citext) RETURNS cstring AS 'textout' LANGUAGE internal STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION citextrecv(internal) RETURNS citext AS 'textrecv' LANGUAGE internal STRICT STABLE PARALLEL SAFE;
CREATE FUNCTION citextsend(citext) RETURNS bytea AS 'textsend' LANGUAGE internal STRICT STABLE PARALLEL SAFE;

-- Full type definition
CREATE TYPE citext (
    INPUT          = citextin,
    OUTPUT         = citextout,
    RECEIVE        = citextrecv,
    SEND           = citextsend,
    INTERNALLENGTH = VARIABLE,
    STORAGE        = extended,
    CATEGORY       = 'S',
    PREFERRED      = false,
    COLLATABLE     = true
);

-- Casts
-- pg: contrib/citext/citext--1.4.sql (casts section)
CREATE CAST (citext AS text)              WITHOUT FUNCTION AS IMPLICIT;
CREATE CAST (citext AS character varying)  WITHOUT FUNCTION AS IMPLICIT;
CREATE CAST (citext AS character)          WITHOUT FUNCTION AS ASSIGNMENT;
CREATE CAST (text AS citext)              WITHOUT FUNCTION AS ASSIGNMENT;
CREATE CAST (character varying AS citext)  WITHOUT FUNCTION AS ASSIGNMENT;

-- Comparison functions
CREATE FUNCTION citext_eq(citext, citext) RETURNS boolean AS 'citext_eq' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION citext_ne(citext, citext) RETURNS boolean AS 'citext_ne' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION citext_lt(citext, citext) RETURNS boolean AS 'citext_lt' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION citext_le(citext, citext) RETURNS boolean AS 'citext_le' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION citext_gt(citext, citext) RETURNS boolean AS 'citext_gt' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION citext_ge(citext, citext) RETURNS boolean AS 'citext_ge' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION citext_cmp(citext, citext) RETURNS integer AS 'citext_cmp' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION citext_hash(citext) RETURNS integer AS 'citext_hash' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;

-- Operators
CREATE OPERATOR = (
    LEFTARG    = citext,
    RIGHTARG   = citext,
    COMMUTATOR = =,
    NEGATOR    = <>,
    PROCEDURE  = citext_eq,
    RESTRICT   = eqsel,
    JOIN       = eqjoinsel,
    HASHES,
    MERGES
);
CREATE OPERATOR <> (
    LEFTARG    = citext,
    RIGHTARG   = citext,
    NEGATOR    = =,
    COMMUTATOR = <>,
    PROCEDURE  = citext_ne,
    RESTRICT   = neqsel,
    JOIN       = neqjoinsel
);
CREATE OPERATOR < (
    LEFTARG    = citext,
    RIGHTARG   = citext,
    NEGATOR    = >=,
    COMMUTATOR = >,
    PROCEDURE  = citext_lt,
    RESTRICT   = scalarltsel,
    JOIN       = scalarltjoinsel
);
CREATE OPERATOR <= (
    LEFTARG    = citext,
    RIGHTARG   = citext,
    NEGATOR    = >,
    COMMUTATOR = >=,
    PROCEDURE  = citext_le,
    RESTRICT   = scalarlesel,
    JOIN       = scalarlejoinsel
);
CREATE OPERATOR >= (
    LEFTARG    = citext,
    RIGHTARG   = citext,
    NEGATOR    = <,
    COMMUTATOR = <=,
    PROCEDURE  = citext_ge,
    RESTRICT   = scalargesel,
    JOIN       = scalargejoinsel
);
CREATE OPERATOR > (
    LEFTARG    = citext,
    RIGHTARG   = citext,
    NEGATOR    = <=,
    COMMUTATOR = <,
    PROCEDURE  = citext_gt,
    RESTRICT   = scalargtsel,
    JOIN       = scalargtjoinsel
);

-- Cast functions (delegate to internal builtins)
-- pg: contrib/citext/citext--1.4.sql
CREATE FUNCTION citext(bpchar) RETURNS citext AS 'rtrim1' LANGUAGE internal IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION citext(boolean) RETURNS citext AS 'booltext' LANGUAGE internal IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION citext(inet) RETURNS citext AS 'network_show' LANGUAGE internal IMMUTABLE STRICT PARALLEL SAFE;

-- Additional casts
CREATE CAST (bpchar AS citext) WITH FUNCTION citext(bpchar) AS ASSIGNMENT;
CREATE CAST (boolean AS citext) WITH FUNCTION citext(boolean) AS ASSIGNMENT;
CREATE CAST (inet AS citext) WITH FUNCTION citext(inet) AS ASSIGNMENT;

-- Aggregates
CREATE FUNCTION citext_smaller(citext, citext) RETURNS citext AS 'citext_smaller' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION citext_larger(citext, citext) RETURNS citext AS 'citext_larger' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;

CREATE AGGREGATE min(citext) (
    SFUNC = citext_smaller,
    STYPE = citext,
    SORTOP = <,
    PARALLEL = safe
);
CREATE AGGREGATE max(citext) (
    SFUNC = citext_larger,
    STYPE = citext,
    SORTOP = >,
    PARALLEL = safe
);

-- Pattern matching (citext, citext)
-- pg: contrib/citext/citext--1.4.sql
CREATE FUNCTION texticlike(citext, citext) RETURNS bool AS 'texticlike' LANGUAGE internal IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION texticnlike(citext, citext) RETURNS bool AS 'texticnlike' LANGUAGE internal IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION texticregexeq(citext, citext) RETURNS bool AS 'texticregexeq' LANGUAGE internal IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION texticregexne(citext, citext) RETURNS bool AS 'texticregexne' LANGUAGE internal IMMUTABLE STRICT PARALLEL SAFE;

CREATE OPERATOR ~ (LEFTARG = citext, RIGHTARG = citext, PROCEDURE = texticregexeq, NEGATOR = !~, RESTRICT = icregexeqsel, JOIN = icregexeqjoinsel);
CREATE OPERATOR ~* (LEFTARG = citext, RIGHTARG = citext, PROCEDURE = texticregexeq, NEGATOR = !~*, RESTRICT = icregexeqsel, JOIN = icregexeqjoinsel);
CREATE OPERATOR !~ (LEFTARG = citext, RIGHTARG = citext, PROCEDURE = texticregexne, NEGATOR = ~, RESTRICT = icregexnesel, JOIN = icregexnejoinsel);
CREATE OPERATOR !~* (LEFTARG = citext, RIGHTARG = citext, PROCEDURE = texticregexne, NEGATOR = ~*, RESTRICT = icregexnesel, JOIN = icregexnejoinsel);
CREATE OPERATOR ~~ (LEFTARG = citext, RIGHTARG = citext, PROCEDURE = texticlike, NEGATOR = !~~, RESTRICT = iclikesel, JOIN = iclikejoinsel);
CREATE OPERATOR ~~* (LEFTARG = citext, RIGHTARG = citext, PROCEDURE = texticlike, NEGATOR = !~~*, RESTRICT = iclikesel, JOIN = iclikejoinsel);
CREATE OPERATOR !~~ (LEFTARG = citext, RIGHTARG = citext, PROCEDURE = texticnlike, NEGATOR = ~~, RESTRICT = icnlikesel, JOIN = icnlikejoinsel);
CREATE OPERATOR !~~* (LEFTARG = citext, RIGHTARG = citext, PROCEDURE = texticnlike, NEGATOR = ~~*, RESTRICT = icnlikesel, JOIN = icnlikejoinsel);

-- Pattern matching (citext, text)
-- pg: contrib/citext/citext--1.4.sql
CREATE FUNCTION texticlike(citext, text) RETURNS bool AS 'texticlike' LANGUAGE internal IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION texticnlike(citext, text) RETURNS bool AS 'texticnlike' LANGUAGE internal IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION texticregexeq(citext, text) RETURNS bool AS 'texticregexeq' LANGUAGE internal IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION texticregexne(citext, text) RETURNS bool AS 'texticregexne' LANGUAGE internal IMMUTABLE STRICT PARALLEL SAFE;

CREATE OPERATOR ~ (LEFTARG = citext, RIGHTARG = text, PROCEDURE = texticregexeq, NEGATOR = !~, RESTRICT = icregexeqsel, JOIN = icregexeqjoinsel);
CREATE OPERATOR ~* (LEFTARG = citext, RIGHTARG = text, PROCEDURE = texticregexeq, NEGATOR = !~*, RESTRICT = icregexeqsel, JOIN = icregexeqjoinsel);
CREATE OPERATOR !~ (LEFTARG = citext, RIGHTARG = text, PROCEDURE = texticregexne, NEGATOR = ~, RESTRICT = icregexnesel, JOIN = icregexnejoinsel);
CREATE OPERATOR !~* (LEFTARG = citext, RIGHTARG = text, PROCEDURE = texticregexne, NEGATOR = ~*, RESTRICT = icregexnesel, JOIN = icregexnejoinsel);
CREATE OPERATOR ~~ (LEFTARG = citext, RIGHTARG = text, PROCEDURE = texticlike, NEGATOR = !~~, RESTRICT = iclikesel, JOIN = iclikejoinsel);
CREATE OPERATOR ~~* (LEFTARG = citext, RIGHTARG = text, PROCEDURE = texticlike, NEGATOR = !~~*, RESTRICT = iclikesel, JOIN = iclikejoinsel);
CREATE OPERATOR !~~ (LEFTARG = citext, RIGHTARG = text, PROCEDURE = texticnlike, NEGATOR = ~~, RESTRICT = icnlikesel, JOIN = icnlikejoinsel);
CREATE OPERATOR !~~* (LEFTARG = citext, RIGHTARG = text, PROCEDURE = texticnlike, NEGATOR = ~~*, RESTRICT = icnlikesel, JOIN = icnlikejoinsel);

-- String matching functions (SQL wrappers)
-- pg: contrib/citext/citext--1.4.sql (regexp_match .. translate)
CREATE FUNCTION regexp_match(citext, citext) RETURNS text[] AS $$
    SELECT pg_catalog.regexp_match( $1::pg_catalog.text, $2::pg_catalog.text, 'i' );
$$ LANGUAGE SQL IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION regexp_match(citext, citext, text) RETURNS text[] AS $$
    SELECT pg_catalog.regexp_match( $1::pg_catalog.text, $2::pg_catalog.text, CASE WHEN pg_catalog.strpos($3, 'c') = 0 THEN  $3 || 'i' ELSE $3 END );
$$ LANGUAGE SQL IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION regexp_matches(citext, citext) RETURNS SETOF text[] AS $$
    SELECT pg_catalog.regexp_matches( $1::pg_catalog.text, $2::pg_catalog.text, 'i' );
$$ LANGUAGE SQL IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION regexp_matches(citext, citext, text) RETURNS SETOF text[] AS $$
    SELECT pg_catalog.regexp_matches( $1::pg_catalog.text, $2::pg_catalog.text, CASE WHEN pg_catalog.strpos($3, 'c') = 0 THEN  $3 || 'i' ELSE $3 END );
$$ LANGUAGE SQL IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION regexp_replace(citext, citext, text) RETURNS text AS $$
    SELECT pg_catalog.regexp_replace( $1::pg_catalog.text, $2::pg_catalog.text, $3, 'i');
$$ LANGUAGE SQL IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION regexp_replace(citext, citext, text, text) RETURNS text AS $$
    SELECT pg_catalog.regexp_replace( $1::pg_catalog.text, $2::pg_catalog.text, $3, CASE WHEN pg_catalog.strpos($4, 'c') = 0 THEN  $4 || 'i' ELSE $4 END);
$$ LANGUAGE SQL IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION regexp_split_to_array(citext, citext) RETURNS text[] AS $$
    SELECT pg_catalog.regexp_split_to_array( $1::pg_catalog.text, $2::pg_catalog.text, 'i' );
$$ LANGUAGE SQL IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION regexp_split_to_array(citext, citext, text) RETURNS text[] AS $$
    SELECT pg_catalog.regexp_split_to_array( $1::pg_catalog.text, $2::pg_catalog.text, CASE WHEN pg_catalog.strpos($3, 'c') = 0 THEN  $3 || 'i' ELSE $3 END );
$$ LANGUAGE SQL IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION regexp_split_to_table(citext, citext) RETURNS SETOF text AS $$
    SELECT pg_catalog.regexp_split_to_table( $1::pg_catalog.text, $2::pg_catalog.text, 'i' );
$$ LANGUAGE SQL IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION regexp_split_to_table(citext, citext, text) RETURNS SETOF text AS $$
    SELECT pg_catalog.regexp_split_to_table( $1::pg_catalog.text, $2::pg_catalog.text, CASE WHEN pg_catalog.strpos($3, 'c') = 0 THEN  $3 || 'i' ELSE $3 END );
$$ LANGUAGE SQL IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION strpos(citext, citext) RETURNS integer AS $$
    SELECT pg_catalog.strpos( pg_catalog.lower( $1::pg_catalog.text ), pg_catalog.lower( $2::pg_catalog.text ) );
$$ LANGUAGE SQL IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION replace(citext, citext, citext) RETURNS text AS $$
    SELECT pg_catalog.regexp_replace( $1::pg_catalog.text, pg_catalog.regexp_replace($2::pg_catalog.text, '([^a-zA-Z_0-9])', E'\\\\\\1', 'g'), $3::pg_catalog.text, 'gi' );
$$ LANGUAGE SQL IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION split_part(citext, citext, integer) RETURNS text AS $$
    SELECT (pg_catalog.regexp_split_to_array( $1::pg_catalog.text, pg_catalog.regexp_replace($2::pg_catalog.text, '([^a-zA-Z_0-9])', E'\\\\\\1', 'g'), 'i'))[$3];
$$ LANGUAGE SQL IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION translate(citext, citext, text) RETURNS text AS $$
    SELECT pg_catalog.translate( pg_catalog.translate( $1::pg_catalog.text, pg_catalog.lower($2::pg_catalog.text), $3), pg_catalog.upper($2::pg_catalog.text), $3);
$$ LANGUAGE SQL IMMUTABLE STRICT PARALLEL SAFE;

-- Operator classes
CREATE OPERATOR CLASS citext_ops DEFAULT FOR TYPE citext USING btree AS
    OPERATOR 1 <,
    OPERATOR 2 <=,
    OPERATOR 3 =,
    OPERATOR 4 >=,
    OPERATOR 5 >,
    FUNCTION 1 citext_cmp(citext, citext);

CREATE OPERATOR CLASS citext_ops DEFAULT FOR TYPE citext USING hash AS
    OPERATOR 1 =,
    FUNCTION 1 citext_hash(citext);

-- Upgrade functions (1.4→1.5→1.6)
-- pg: contrib/citext/citext--1.4--1.5.sql
CREATE FUNCTION citext_pattern_lt(citext, citext) RETURNS bool AS 'citext_pattern_lt' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION citext_pattern_le(citext, citext) RETURNS bool AS 'citext_pattern_le' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION citext_pattern_gt(citext, citext) RETURNS bool AS 'citext_pattern_gt' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION citext_pattern_ge(citext, citext) RETURNS bool AS 'citext_pattern_ge' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION citext_pattern_cmp(citext, citext) RETURNS integer AS 'citext_pattern_cmp' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
-- pg: contrib/citext/citext--1.5--1.6.sql
CREATE FUNCTION citext_hash_extended(citext, bigint) RETURNS bigint AS 'citext_hash_extended' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
`

// hstoreSQL is the DDL for the hstore extension (key-value store type).
// Faithfully translated from contrib/hstore/hstore--1.4.sql.
const hstoreSQL = `
-- Shell type
CREATE TYPE hstore;

-- I/O functions
CREATE FUNCTION hstore_in(cstring) RETURNS hstore AS 'hstore_in' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION hstore_out(hstore) RETURNS cstring AS 'hstore_out' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION hstore_recv(internal) RETURNS hstore AS 'hstore_recv' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION hstore_send(hstore) RETURNS bytea AS 'hstore_send' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;

-- Full type definition
CREATE TYPE hstore (
    INPUT          = hstore_in,
    OUTPUT         = hstore_out,
    RECEIVE        = hstore_recv,
    SEND           = hstore_send,
    INTERNALLENGTH = VARIABLE,
    STORAGE        = extended
);

-- Version diagnostic
CREATE FUNCTION hstore_version_diag(hstore) RETURNS integer AS 'hstore_version_diag' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;

-- Core functions
CREATE FUNCTION fetchval(hstore, text) RETURNS text AS 'hstore_fetchval' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION slice_array(hstore, text[]) RETURNS text[] AS 'hstore_slice_to_array' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION slice(hstore, text[]) RETURNS hstore AS 'hstore_slice_to_hstore' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION isexists(hstore, text) RETURNS bool AS 'hstore_exists' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION exist(hstore, text) RETURNS bool AS 'hstore_exists' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION exists_any(hstore, text[]) RETURNS bool AS 'hstore_exists_any' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION exists_all(hstore, text[]) RETURNS bool AS 'hstore_exists_all' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION isdefined(hstore, text) RETURNS bool AS 'hstore_defined' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION defined(hstore, text) RETURNS bool AS 'hstore_defined' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION delete(hstore, text) RETURNS hstore AS 'hstore_delete' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION delete(hstore, text[]) RETURNS hstore AS 'hstore_delete_array' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION delete(hstore, hstore) RETURNS hstore AS 'hstore_delete_hstore' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION hs_concat(hstore, hstore) RETURNS hstore AS 'hstore_concat' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION hs_contains(hstore, hstore) RETURNS bool AS 'hstore_contains' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION hs_contained(hstore, hstore) RETURNS bool AS 'hstore_contained' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION tconvert(text, text) RETURNS hstore AS 'hstore_from_text' LANGUAGE C IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION hstore(text, text) RETURNS hstore AS 'hstore_from_text' LANGUAGE C IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION hstore(text[], text[]) RETURNS hstore AS 'hstore_from_arrays' LANGUAGE C IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION hstore(text[]) RETURNS hstore AS 'hstore_from_array' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION hstore(record) RETURNS hstore AS 'hstore_from_record' LANGUAGE C IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION hstore_to_array(hstore) RETURNS text[] AS 'hstore_to_array' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION hstore_to_matrix(hstore) RETURNS text[] AS 'hstore_to_matrix' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION akeys(hstore) RETURNS text[] AS 'hstore_akeys' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION avals(hstore) RETURNS text[] AS 'hstore_avals' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION skeys(hstore) RETURNS SETOF text AS 'hstore_skeys' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION svals(hstore) RETURNS SETOF text AS 'hstore_svals' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION each(IN hs hstore, OUT key text, OUT value text) RETURNS SETOF record AS 'hstore_each' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION populate_record(anyelement, hstore) RETURNS anyelement AS 'hstore_populate_record' LANGUAGE C IMMUTABLE PARALLEL SAFE;

-- JSON conversion functions
CREATE FUNCTION hstore_to_json(hstore) RETURNS json AS 'hstore_to_json' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION hstore_to_json_loose(hstore) RETURNS json AS 'hstore_to_json_loose' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION hstore_to_jsonb(hstore) RETURNS jsonb AS 'hstore_to_jsonb' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION hstore_to_jsonb_loose(hstore) RETURNS jsonb AS 'hstore_to_jsonb_loose' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;

-- Casts
CREATE CAST (text[] AS hstore) WITH FUNCTION hstore(text[]);
CREATE CAST (hstore AS json) WITH FUNCTION hstore_to_json(hstore);
CREATE CAST (hstore AS jsonb) WITH FUNCTION hstore_to_jsonb(hstore);

-- Operators
CREATE OPERATOR -> (LEFTARG = hstore, RIGHTARG = text, PROCEDURE = fetchval);
CREATE OPERATOR -> (LEFTARG = hstore, RIGHTARG = text[], PROCEDURE = slice_array);
CREATE OPERATOR ? (LEFTARG = hstore, RIGHTARG = text, PROCEDURE = exist, RESTRICT = contsel, JOIN = contjoinsel);
CREATE OPERATOR ?| (LEFTARG = hstore, RIGHTARG = text[], PROCEDURE = exists_any, RESTRICT = contsel, JOIN = contjoinsel);
CREATE OPERATOR ?& (LEFTARG = hstore, RIGHTARG = text[], PROCEDURE = exists_all, RESTRICT = contsel, JOIN = contjoinsel);
CREATE OPERATOR - (LEFTARG = hstore, RIGHTARG = text, PROCEDURE = delete);
CREATE OPERATOR - (LEFTARG = hstore, RIGHTARG = text[], PROCEDURE = delete);
CREATE OPERATOR - (LEFTARG = hstore, RIGHTARG = hstore, PROCEDURE = delete);
CREATE OPERATOR || (LEFTARG = hstore, RIGHTARG = hstore, PROCEDURE = hs_concat);
CREATE OPERATOR @> (LEFTARG = hstore, RIGHTARG = hstore, PROCEDURE = hs_contains, COMMUTATOR = <@, RESTRICT = contsel, JOIN = contjoinsel);
CREATE OPERATOR <@ (LEFTARG = hstore, RIGHTARG = hstore, PROCEDURE = hs_contained, COMMUTATOR = @>, RESTRICT = contsel, JOIN = contjoinsel);
CREATE OPERATOR @ (LEFTARG = hstore, RIGHTARG = hstore, PROCEDURE = hs_contains, COMMUTATOR = ~, RESTRICT = contsel, JOIN = contjoinsel);
CREATE OPERATOR ~ (LEFTARG = hstore, RIGHTARG = hstore, PROCEDURE = hs_contained, COMMUTATOR = @, RESTRICT = contsel, JOIN = contjoinsel);
CREATE OPERATOR %% (RIGHTARG = hstore, PROCEDURE = hstore_to_array);
CREATE OPERATOR %# (RIGHTARG = hstore, PROCEDURE = hstore_to_matrix);
CREATE OPERATOR #= (LEFTARG = anyelement, RIGHTARG = hstore, PROCEDURE = populate_record);

-- btree support
CREATE FUNCTION hstore_eq(hstore, hstore) RETURNS boolean AS 'hstore_eq' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION hstore_ne(hstore, hstore) RETURNS boolean AS 'hstore_ne' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION hstore_gt(hstore, hstore) RETURNS boolean AS 'hstore_gt' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION hstore_ge(hstore, hstore) RETURNS boolean AS 'hstore_ge' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION hstore_lt(hstore, hstore) RETURNS boolean AS 'hstore_lt' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION hstore_le(hstore, hstore) RETURNS boolean AS 'hstore_le' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION hstore_cmp(hstore, hstore) RETURNS integer AS 'hstore_cmp' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;

CREATE OPERATOR = (LEFTARG = hstore, RIGHTARG = hstore, PROCEDURE = hstore_eq, COMMUTATOR = =, NEGATOR = <>, RESTRICT = eqsel, JOIN = eqjoinsel, MERGES, HASHES);
CREATE OPERATOR <> (LEFTARG = hstore, RIGHTARG = hstore, PROCEDURE = hstore_ne, COMMUTATOR = <>, NEGATOR = =, RESTRICT = neqsel, JOIN = neqjoinsel);
CREATE OPERATOR #<# (LEFTARG = hstore, RIGHTARG = hstore, PROCEDURE = hstore_lt, COMMUTATOR = #>#, NEGATOR = #>=#, RESTRICT = scalarltsel, JOIN = scalarltjoinsel);
CREATE OPERATOR #<=# (LEFTARG = hstore, RIGHTARG = hstore, PROCEDURE = hstore_le, COMMUTATOR = #>=#, NEGATOR = #>#, RESTRICT = scalarltsel, JOIN = scalarltjoinsel);
CREATE OPERATOR #># (LEFTARG = hstore, RIGHTARG = hstore, PROCEDURE = hstore_gt, COMMUTATOR = #<#, NEGATOR = #<=#, RESTRICT = scalargtsel, JOIN = scalargtjoinsel);
CREATE OPERATOR #>=# (LEFTARG = hstore, RIGHTARG = hstore, PROCEDURE = hstore_ge, COMMUTATOR = #<=#, NEGATOR = #<#, RESTRICT = scalargtsel, JOIN = scalargtjoinsel);

-- hash support
CREATE FUNCTION hstore_hash(hstore) RETURNS integer AS 'hstore_hash' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;

-- Operator classes
CREATE OPERATOR CLASS btree_hstore_ops DEFAULT FOR TYPE hstore USING btree AS
    OPERATOR 1 #<#,
    OPERATOR 2 #<=#,
    OPERATOR 3 =,
    OPERATOR 4 #>=#,
    OPERATOR 5 #>#,
    FUNCTION 1 hstore_cmp(hstore, hstore);

CREATE OPERATOR CLASS hash_hstore_ops DEFAULT FOR TYPE hstore USING hash AS
    OPERATOR 1 =,
    FUNCTION 1 hstore_hash(hstore);

-- GiST support
-- pg: contrib/hstore/hstore--1.4.sql (GiST support section)

CREATE TYPE ghstore;

CREATE FUNCTION ghstore_in(cstring) RETURNS ghstore AS 'ghstore_in' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION ghstore_out(ghstore) RETURNS cstring AS 'ghstore_out' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;

CREATE TYPE ghstore (
    INTERNALLENGTH = VARIABLE,
    INPUT          = ghstore_in,
    OUTPUT         = ghstore_out
);

CREATE FUNCTION ghstore_compress(internal) RETURNS internal AS 'ghstore_compress' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION ghstore_decompress(internal) RETURNS internal AS 'ghstore_decompress' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION ghstore_penalty(internal, internal, internal) RETURNS internal AS 'ghstore_penalty' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION ghstore_picksplit(internal, internal) RETURNS internal AS 'ghstore_picksplit' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION ghstore_union(internal, internal) RETURNS ghstore AS 'ghstore_union' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION ghstore_same(ghstore, ghstore, internal) RETURNS internal AS 'ghstore_same' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION ghstore_consistent(internal, hstore, smallint, oid, internal) RETURNS bool AS 'ghstore_consistent' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;

-- GIN support
-- pg: contrib/hstore/hstore--1.4.sql (GIN support section)

CREATE FUNCTION gin_extract_hstore(hstore, internal) RETURNS internal AS 'gin_extract_hstore' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION gin_extract_hstore_query(hstore, internal, smallint, internal, internal) RETURNS internal AS 'gin_extract_hstore_query' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION gin_consistent_hstore(internal, smallint, hstore, integer, internal, internal) RETURNS bool AS 'gin_consistent_hstore' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;

-- Upgrade functions (1.5→1.6→1.7→1.8)
-- pg: contrib/hstore/hstore--1.5--1.6.sql
CREATE FUNCTION hstore_hash_extended(hstore, bigint) RETURNS bigint AS 'hstore_hash_extended' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
-- pg: contrib/hstore/hstore--1.6--1.7.sql
CREATE FUNCTION ghstore_options(internal) RETURNS void AS 'ghstore_options' LANGUAGE C IMMUTABLE PARALLEL SAFE;
-- pg: contrib/hstore/hstore--1.7--1.8.sql
CREATE FUNCTION hstore_subscript_handler(internal) RETURNS internal AS 'hstore_subscript_handler' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
`

// pgvectorSQL is the DDL for the pgvector extension (vector similarity search).
// Derived from pgvector's vector--0.7.0.sql — core types, operators, access methods.
const pgvectorSQL = `
-- vector type (with typmod for dimensions)
CREATE TYPE vector;
CREATE FUNCTION vector_in(cstring, oid, integer) RETURNS vector AS 'vector_in' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION vector_out(vector) RETURNS cstring AS 'vector_out' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION vector_typmod_in(cstring[]) RETURNS integer AS 'vector_typmod_in' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION vector_typmod_out(integer) RETURNS cstring AS 'vector_typmod_out' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION vector_recv(internal, oid, integer) RETURNS vector AS 'vector_recv' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION vector_send(vector) RETURNS bytea AS 'vector_send' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE TYPE vector (
    INPUT     = vector_in,
    OUTPUT    = vector_out,
    TYPMOD_IN = vector_typmod_in,
    TYPMOD_OUT = vector_typmod_out,
    RECEIVE   = vector_recv,
    SEND      = vector_send,
    INTERNALLENGTH = VARIABLE,
    STORAGE   = extended
);

-- halfvec type (half-precision vector, with typmod for dimensions)
CREATE TYPE halfvec;
CREATE FUNCTION halfvec_in(cstring, oid, integer) RETURNS halfvec AS 'halfvec_in' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION halfvec_out(halfvec) RETURNS cstring AS 'halfvec_out' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION halfvec_typmod_in(cstring[]) RETURNS integer AS 'halfvec_typmod_in' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION halfvec_typmod_out(integer) RETURNS cstring AS 'halfvec_typmod_out' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION halfvec_recv(internal, oid, integer) RETURNS halfvec AS 'halfvec_recv' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION halfvec_send(halfvec) RETURNS bytea AS 'halfvec_send' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE TYPE halfvec (
    INPUT     = halfvec_in,
    OUTPUT    = halfvec_out,
    TYPMOD_IN = halfvec_typmod_in,
    TYPMOD_OUT = halfvec_typmod_out,
    RECEIVE   = halfvec_recv,
    SEND      = halfvec_send,
    INTERNALLENGTH = VARIABLE,
    STORAGE   = extended
);

-- sparsevec type
CREATE TYPE sparsevec;
CREATE FUNCTION sparsevec_in(cstring, oid, integer) RETURNS sparsevec AS 'sparsevec_in' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION sparsevec_out(sparsevec) RETURNS cstring AS 'sparsevec_out' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION sparsevec_typmod_in(cstring[]) RETURNS integer AS 'sparsevec_typmod_in' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION sparsevec_typmod_out(integer) RETURNS cstring AS 'sparsevec_typmod_out' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION sparsevec_recv(internal, oid, integer) RETURNS sparsevec AS 'sparsevec_recv' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION sparsevec_send(sparsevec) RETURNS bytea AS 'sparsevec_send' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE TYPE sparsevec (
    INPUT     = sparsevec_in,
    OUTPUT    = sparsevec_out,
    TYPMOD_IN = sparsevec_typmod_in,
    TYPMOD_OUT = sparsevec_typmod_out,
    RECEIVE   = sparsevec_recv,
    SEND      = sparsevec_send,
    INTERNALLENGTH = VARIABLE,
    STORAGE   = extended
);

-- Casts between vector types
CREATE FUNCTION vector_to_halfvec(vector, integer, boolean) RETURNS halfvec AS 'vector_to_halfvec' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION halfvec_to_vector(halfvec, integer, boolean) RETURNS vector AS 'halfvec_to_vector' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE CAST (vector AS halfvec) WITH FUNCTION vector_to_halfvec(vector, integer, boolean) AS IMPLICIT;
CREATE CAST (halfvec AS vector) WITH FUNCTION halfvec_to_vector(halfvec, integer, boolean) AS IMPLICIT;

-- Distance functions
CREATE FUNCTION l2_distance(vector, vector) RETURNS double precision AS 'vector_l2_distance' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION cosine_distance(vector, vector) RETURNS double precision AS 'vector_cosine_distance' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION inner_product(vector, vector) RETURNS double precision AS 'vector_inner_product' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION halfvec_l2_distance(halfvec, halfvec) RETURNS double precision AS 'halfvec_l2_distance' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION halfvec_cosine_distance(halfvec, halfvec) RETURNS double precision AS 'halfvec_cosine_distance' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION halfvec_inner_product(halfvec, halfvec) RETURNS double precision AS 'halfvec_inner_product' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;

-- Distance operators
CREATE OPERATOR <-> (
    LEFTARG    = vector,
    RIGHTARG   = vector,
    PROCEDURE  = l2_distance,
    COMMUTATOR = <->
);
CREATE OPERATOR <=> (
    LEFTARG    = vector,
    RIGHTARG   = vector,
    PROCEDURE  = cosine_distance,
    COMMUTATOR = <=>
);
CREATE OPERATOR <#> (
    LEFTARG    = vector,
    RIGHTARG   = vector,
    PROCEDURE  = inner_product,
    COMMUTATOR = <#>
);
CREATE OPERATOR <-> (
    LEFTARG    = halfvec,
    RIGHTARG   = halfvec,
    PROCEDURE  = halfvec_l2_distance,
    COMMUTATOR = <->
);
CREATE OPERATOR <=> (
    LEFTARG    = halfvec,
    RIGHTARG   = halfvec,
    PROCEDURE  = halfvec_cosine_distance,
    COMMUTATOR = <=>
);
CREATE OPERATOR <#> (
    LEFTARG    = halfvec,
    RIGHTARG   = halfvec,
    PROCEDURE  = halfvec_inner_product,
    COMMUTATOR = <#>
);

-- Access method handlers (pgddl: handler functions are stubs)
CREATE FUNCTION hnswhandler(internal) RETURNS index_am_handler AS 'hnswhandler' LANGUAGE C;
CREATE FUNCTION ivfflathandler(internal) RETURNS index_am_handler AS 'ivfflathandler' LANGUAGE C;

-- Access methods
CREATE ACCESS METHOD hnsw TYPE INDEX HANDLER hnswhandler;
CREATE ACCESS METHOD ivfflat TYPE INDEX HANDLER ivfflathandler;

-- Operator classes for vector
CREATE OPERATOR CLASS vector_l2_ops DEFAULT FOR TYPE vector USING hnsw AS
    OPERATOR 1 <-> (vector, vector);
CREATE OPERATOR CLASS vector_cosine_ops FOR TYPE vector USING hnsw AS
    OPERATOR 1 <=> (vector, vector);
CREATE OPERATOR CLASS vector_ip_ops FOR TYPE vector USING hnsw AS
    OPERATOR 1 <#> (vector, vector);

-- Operator classes for halfvec
CREATE OPERATOR CLASS halfvec_l2_ops DEFAULT FOR TYPE halfvec USING hnsw AS
    OPERATOR 1 <-> (halfvec, halfvec);
CREATE OPERATOR CLASS halfvec_cosine_ops FOR TYPE halfvec USING hnsw AS
    OPERATOR 1 <=> (halfvec, halfvec);
CREATE OPERATOR CLASS halfvec_ip_ops FOR TYPE halfvec USING hnsw AS
    OPERATOR 1 <#> (halfvec, halfvec);

-- ivfflat operator classes
CREATE OPERATOR CLASS vector_l2_ops DEFAULT FOR TYPE vector USING ivfflat AS
    OPERATOR 1 <-> (vector, vector);
CREATE OPERATOR CLASS vector_cosine_ops FOR TYPE vector USING ivfflat AS
    OPERATOR 1 <=> (vector, vector);
CREATE OPERATOR CLASS halfvec_l2_ops DEFAULT FOR TYPE halfvec USING ivfflat AS
    OPERATOR 1 <-> (halfvec, halfvec);
CREATE OPERATOR CLASS halfvec_cosine_ops FOR TYPE halfvec USING ivfflat AS
    OPERATOR 1 <=> (halfvec, halfvec);
`

// ltreeSQL is the DDL subset for the ltree extension.
// Derived from contrib/ltree/ltree--1.1.sql through ltree--1.2--1.3.sql.
const ltreeSQL = `
CREATE FUNCTION ltree_in(cstring) RETURNS ltree AS 'ltree_in' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION ltree_out(ltree) RETURNS cstring AS 'ltree_out' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION ltree_recv(internal) RETURNS ltree AS 'ltree_recv' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION ltree_send(ltree) RETURNS bytea AS 'ltree_send' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE TYPE ltree (INPUT = ltree_in, OUTPUT = ltree_out, RECEIVE = ltree_recv, SEND = ltree_send, INTERNALLENGTH = VARIABLE, STORAGE = extended);

CREATE FUNCTION lquery_in(cstring) RETURNS lquery AS 'lquery_in' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION lquery_out(lquery) RETURNS cstring AS 'lquery_out' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION lquery_recv(internal) RETURNS lquery AS 'lquery_recv' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION lquery_send(lquery) RETURNS bytea AS 'lquery_send' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE TYPE lquery (INPUT = lquery_in, OUTPUT = lquery_out, RECEIVE = lquery_recv, SEND = lquery_send, INTERNALLENGTH = VARIABLE, STORAGE = extended);

CREATE FUNCTION ltxtq_in(cstring) RETURNS ltxtquery AS 'ltxtq_in' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION ltxtq_out(ltxtquery) RETURNS cstring AS 'ltxtq_out' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION ltxtq_recv(internal) RETURNS ltxtquery AS 'ltxtq_recv' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION ltxtq_send(ltxtquery) RETURNS bytea AS 'ltxtq_send' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE TYPE ltxtquery (INPUT = ltxtq_in, OUTPUT = ltxtq_out, RECEIVE = ltxtq_recv, SEND = ltxtq_send, INTERNALLENGTH = VARIABLE, STORAGE = extended);

CREATE FUNCTION ltree_gist_in(cstring) RETURNS ltree_gist AS 'ltree_gist_in' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION ltree_gist_out(ltree_gist) RETURNS cstring AS 'ltree_gist_out' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE TYPE ltree_gist (INPUT = ltree_gist_in, OUTPUT = ltree_gist_out, INTERNALLENGTH = VARIABLE, STORAGE = plain);

CREATE FUNCTION ltree_cmp(ltree, ltree) RETURNS integer AS 'ltree_cmp' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION ltree_lt(ltree, ltree) RETURNS boolean AS 'ltree_lt' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION ltree_le(ltree, ltree) RETURNS boolean AS 'ltree_le' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION ltree_eq(ltree, ltree) RETURNS boolean AS 'ltree_eq' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION ltree_ge(ltree, ltree) RETURNS boolean AS 'ltree_ge' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION ltree_gt(ltree, ltree) RETURNS boolean AS 'ltree_gt' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION ltree_ne(ltree, ltree) RETURNS boolean AS 'ltree_ne' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;

CREATE OPERATOR < (LEFTARG = ltree, RIGHTARG = ltree, PROCEDURE = ltree_lt, COMMUTATOR = >, NEGATOR = >=);
CREATE OPERATOR <= (LEFTARG = ltree, RIGHTARG = ltree, PROCEDURE = ltree_le, COMMUTATOR = >=, NEGATOR = >);
CREATE OPERATOR = (LEFTARG = ltree, RIGHTARG = ltree, PROCEDURE = ltree_eq, COMMUTATOR = =, NEGATOR = <>, HASHES, MERGES);
CREATE OPERATOR <> (LEFTARG = ltree, RIGHTARG = ltree, PROCEDURE = ltree_ne, COMMUTATOR = <>, NEGATOR = =);
CREATE OPERATOR >= (LEFTARG = ltree, RIGHTARG = ltree, PROCEDURE = ltree_ge, COMMUTATOR = <=, NEGATOR = <);
CREATE OPERATOR > (LEFTARG = ltree, RIGHTARG = ltree, PROCEDURE = ltree_gt, COMMUTATOR = <, NEGATOR = <=);

CREATE FUNCTION subltree(ltree, integer, integer) RETURNS ltree AS 'subltree' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION subpath(ltree, integer, integer) RETURNS ltree AS 'subpath' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION subpath(ltree, integer) RETURNS ltree AS 'subpath' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION nlevel(ltree) RETURNS integer AS 'nlevel' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION ltree2text(ltree) RETURNS text AS 'ltree2text' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION text2ltree(text) RETURNS ltree AS 'text2ltree' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION ltree_isparent(ltree, ltree) RETURNS boolean AS 'ltree_isparent' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION ltree_risparent(ltree, ltree) RETURNS boolean AS 'ltree_risparent' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION ltree_addltree(ltree, ltree) RETURNS ltree AS 'ltree_addltree' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION ltree_addtext(ltree, text) RETURNS ltree AS 'ltree_addtext' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION ltree_textadd(text, ltree) RETURNS ltree AS 'ltree_textadd' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;

CREATE OPERATOR @> (LEFTARG = ltree, RIGHTARG = ltree, PROCEDURE = ltree_isparent, COMMUTATOR = <@);
CREATE OPERATOR <@ (LEFTARG = ltree, RIGHTARG = ltree, PROCEDURE = ltree_risparent, COMMUTATOR = @>);
CREATE OPERATOR || (LEFTARG = ltree, RIGHTARG = ltree, PROCEDURE = ltree_addltree);
CREATE OPERATOR || (LEFTARG = ltree, RIGHTARG = text, PROCEDURE = ltree_addtext);
CREATE OPERATOR || (LEFTARG = text, RIGHTARG = ltree, PROCEDURE = ltree_textadd);

CREATE FUNCTION ltq_regex(ltree, lquery) RETURNS boolean AS 'ltq_regex' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION ltq_rregex(lquery, ltree) RETURNS boolean AS 'ltq_rregex' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION lt_q_regex(ltree, lquery[]) RETURNS boolean AS 'lt_q_regex' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION lt_q_rregex(lquery[], ltree) RETURNS boolean AS 'lt_q_rregex' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION ltxtq_exec(ltree, ltxtquery) RETURNS boolean AS 'ltxtq_exec' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION ltxtq_rexec(ltxtquery, ltree) RETURNS boolean AS 'ltxtq_rexec' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;

CREATE OPERATOR ~ (LEFTARG = ltree, RIGHTARG = lquery, PROCEDURE = ltq_regex, COMMUTATOR = ~);
CREATE OPERATOR ~ (LEFTARG = lquery, RIGHTARG = ltree, PROCEDURE = ltq_rregex, COMMUTATOR = ~);
CREATE OPERATOR ? (LEFTARG = ltree, RIGHTARG = lquery[], PROCEDURE = lt_q_regex, COMMUTATOR = ?);
CREATE OPERATOR ? (LEFTARG = lquery[], RIGHTARG = ltree, PROCEDURE = lt_q_rregex, COMMUTATOR = ?);
CREATE OPERATOR @ (LEFTARG = ltree, RIGHTARG = ltxtquery, PROCEDURE = ltxtq_exec, COMMUTATOR = @);
CREATE OPERATOR @ (LEFTARG = ltxtquery, RIGHTARG = ltree, PROCEDURE = ltxtq_rexec, COMMUTATOR = @);

CREATE FUNCTION ltree_consistent(internal, ltree, smallint, oid, internal) RETURNS boolean AS 'ltree_consistent' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION ltree_union(internal, internal) RETURNS ltree_gist AS 'ltree_union' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION ltree_compress(internal) RETURNS internal AS 'ltree_compress' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION ltree_decompress(internal) RETURNS internal AS 'ltree_decompress' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION ltree_penalty(internal, internal, internal) RETURNS internal AS 'ltree_penalty' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION ltree_picksplit(internal, internal) RETURNS internal AS 'ltree_picksplit' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION ltree_same(ltree_gist, ltree_gist, internal) RETURNS internal AS 'ltree_same' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION hash_ltree(ltree) RETURNS integer AS 'hash_ltree' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION hash_ltree_extended(ltree, bigint) RETURNS bigint AS 'hash_ltree_extended' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;

CREATE OPERATOR CLASS ltree_ops DEFAULT FOR TYPE ltree USING btree AS OPERATOR 1 <, OPERATOR 2 <=, OPERATOR 3 =, OPERATOR 4 >=, OPERATOR 5 >, FUNCTION 1 ltree_cmp(ltree, ltree);
CREATE OPERATOR CLASS gist_ltree_ops DEFAULT FOR TYPE ltree USING gist AS OPERATOR 10 @>, OPERATOR 11 <@, OPERATOR 12 ~ (ltree, lquery), OPERATOR 14 @ (ltree, ltxtquery), FUNCTION 1 ltree_consistent(internal, ltree, smallint, oid, internal), STORAGE ltree_gist;
CREATE OPERATOR CLASS hash_ltree_ops DEFAULT FOR TYPE ltree USING hash AS OPERATOR 1 =, FUNCTION 1 hash_ltree(ltree);
`

// dblinkSQL is the DDL subset for the dblink extension.
// Derived from contrib/dblink/dblink--1.2.sql.
const dblinkSQL = `
CREATE FUNCTION dblink_connect(text) RETURNS text AS 'dblink_connect' LANGUAGE C STRICT PARALLEL RESTRICTED;
CREATE FUNCTION dblink_connect(text, text) RETURNS text AS 'dblink_connect' LANGUAGE C STRICT PARALLEL RESTRICTED;
CREATE FUNCTION dblink_disconnect() RETURNS text AS 'dblink_disconnect' LANGUAGE C STRICT PARALLEL RESTRICTED;
CREATE FUNCTION dblink_disconnect(text) RETURNS text AS 'dblink_disconnect' LANGUAGE C STRICT PARALLEL RESTRICTED;
CREATE FUNCTION dblink_open(text, text) RETURNS text AS 'dblink_open' LANGUAGE C STRICT PARALLEL RESTRICTED;
CREATE FUNCTION dblink_open(text, text, boolean) RETURNS text AS 'dblink_open' LANGUAGE C STRICT PARALLEL RESTRICTED;
CREATE FUNCTION dblink_open(text, text, text) RETURNS text AS 'dblink_open' LANGUAGE C STRICT PARALLEL RESTRICTED;
CREATE FUNCTION dblink_open(text, text, text, boolean) RETURNS text AS 'dblink_open' LANGUAGE C STRICT PARALLEL RESTRICTED;
CREATE FUNCTION dblink_fetch(text, integer) RETURNS SETOF record AS 'dblink_fetch' LANGUAGE C STRICT PARALLEL RESTRICTED;
CREATE FUNCTION dblink_fetch(text, integer, boolean) RETURNS SETOF record AS 'dblink_fetch' LANGUAGE C STRICT PARALLEL RESTRICTED;
CREATE FUNCTION dblink_fetch(text, text, integer) RETURNS SETOF record AS 'dblink_fetch' LANGUAGE C STRICT PARALLEL RESTRICTED;
CREATE FUNCTION dblink_fetch(text, text, integer, boolean) RETURNS SETOF record AS 'dblink_fetch' LANGUAGE C STRICT PARALLEL RESTRICTED;
CREATE FUNCTION dblink_close(text) RETURNS text AS 'dblink_close' LANGUAGE C STRICT PARALLEL RESTRICTED;
CREATE FUNCTION dblink_close(text, boolean) RETURNS text AS 'dblink_close' LANGUAGE C STRICT PARALLEL RESTRICTED;
CREATE FUNCTION dblink_close(text, text) RETURNS text AS 'dblink_close' LANGUAGE C STRICT PARALLEL RESTRICTED;
CREATE FUNCTION dblink_close(text, text, boolean) RETURNS text AS 'dblink_close' LANGUAGE C STRICT PARALLEL RESTRICTED;
CREATE FUNCTION dblink(text, text) RETURNS SETOF record AS 'dblink_record' LANGUAGE C STRICT PARALLEL RESTRICTED;
CREATE FUNCTION dblink(text, text, boolean) RETURNS SETOF record AS 'dblink_record' LANGUAGE C STRICT PARALLEL RESTRICTED;
CREATE FUNCTION dblink(text) RETURNS SETOF record AS 'dblink_record' LANGUAGE C STRICT PARALLEL RESTRICTED;
CREATE FUNCTION dblink(text, boolean) RETURNS SETOF record AS 'dblink_record' LANGUAGE C STRICT PARALLEL RESTRICTED;
CREATE FUNCTION dblink_exec(text, text) RETURNS text AS 'dblink_exec' LANGUAGE C STRICT PARALLEL RESTRICTED;
CREATE FUNCTION dblink_exec(text, text, boolean) RETURNS text AS 'dblink_exec' LANGUAGE C STRICT PARALLEL RESTRICTED;
CREATE FUNCTION dblink_exec(text) RETURNS text AS 'dblink_exec' LANGUAGE C STRICT PARALLEL RESTRICTED;
CREATE FUNCTION dblink_exec(text, boolean) RETURNS text AS 'dblink_exec' LANGUAGE C STRICT PARALLEL RESTRICTED;
CREATE TYPE dblink_pkey_results AS (position integer, colname text);
CREATE FUNCTION dblink_get_pkey(text) RETURNS SETOF dblink_pkey_results AS 'dblink_get_pkey' LANGUAGE C STRICT PARALLEL RESTRICTED;
CREATE FUNCTION dblink_build_sql_insert(text, int2vector, integer, text[], text[]) RETURNS text AS 'dblink_build_sql_insert' LANGUAGE C STRICT PARALLEL RESTRICTED;
CREATE FUNCTION dblink_build_sql_delete(text, int2vector, integer, text[]) RETURNS text AS 'dblink_build_sql_delete' LANGUAGE C STRICT PARALLEL RESTRICTED;
CREATE FUNCTION dblink_build_sql_update(text, int2vector, integer, text[], text[]) RETURNS text AS 'dblink_build_sql_update' LANGUAGE C STRICT PARALLEL RESTRICTED;
CREATE FUNCTION dblink_current_query() RETURNS text AS 'dblink_current_query' LANGUAGE C PARALLEL RESTRICTED;
CREATE FUNCTION dblink_send_query(text, text) RETURNS integer AS 'dblink_send_query' LANGUAGE C STRICT PARALLEL RESTRICTED;
CREATE FUNCTION dblink_is_busy(text) RETURNS integer AS 'dblink_is_busy' LANGUAGE C STRICT PARALLEL RESTRICTED;
CREATE FUNCTION dblink_get_result(text) RETURNS SETOF record AS 'dblink_get_result' LANGUAGE C STRICT PARALLEL RESTRICTED;
CREATE FUNCTION dblink_get_result(text, boolean) RETURNS SETOF record AS 'dblink_get_result' LANGUAGE C STRICT PARALLEL RESTRICTED;
CREATE FUNCTION dblink_get_connections() RETURNS text[] AS 'dblink_get_connections' LANGUAGE C PARALLEL RESTRICTED;
CREATE FUNCTION dblink_cancel_query(text) RETURNS text AS 'dblink_cancel_query' LANGUAGE C STRICT PARALLEL RESTRICTED;
CREATE FUNCTION dblink_error_message(text) RETURNS text AS 'dblink_error_message' LANGUAGE C STRICT PARALLEL RESTRICTED;
CREATE FUNCTION dblink_get_notify() RETURNS SETOF record AS 'dblink_get_notify' LANGUAGE C STRICT PARALLEL RESTRICTED;
CREATE FUNCTION dblink_get_notify(text) RETURNS SETOF record AS 'dblink_get_notify' LANGUAGE C STRICT PARALLEL RESTRICTED;
CREATE FUNCTION dblink_fdw_validator(text[], oid) RETURNS void AS 'dblink_fdw_validator' LANGUAGE C STRICT PARALLEL SAFE;
CREATE FOREIGN DATA WRAPPER dblink_fdw VALIDATOR dblink_fdw_validator;
`

// intarraySQL is the DDL subset for the intarray extension.
// Derived from contrib/intarray/intarray--1.5.sql.
const intarraySQL = `
CREATE FUNCTION sort(integer[]) RETURNS integer[] AS 'sort' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION sort(integer[], text) RETURNS integer[] AS 'sort' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION uniq(integer[]) RETURNS integer[] AS 'uniq' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION idx(integer[], integer) RETURNS integer AS 'idx' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION subarray(integer[], integer, integer) RETURNS integer[] AS 'subarray' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION icount(integer[]) RETURNS integer AS 'icount' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION intarray_push_elem(integer[], integer) RETURNS integer[] AS 'intarray_push_elem' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION intarray_push_array(integer[], integer[]) RETURNS integer[] AS 'intarray_push_array' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION intarray_del_elem(integer[], integer) RETURNS integer[] AS 'intarray_del_elem' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION intset(integer) RETURNS integer[] AS 'intset' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION intarray_union(integer[], integer[]) RETURNS integer[] AS 'intarray_union' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION intarray_inter(integer[], integer[]) RETURNS integer[] AS 'intarray_inter' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION intarray_contained(integer[], integer[]) RETURNS boolean AS 'intarray_contained' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION intarray_contains(integer[], integer[]) RETURNS boolean AS 'intarray_contains' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION intarray_overlap(integer[], integer[]) RETURNS boolean AS 'intarray_overlap' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;

CREATE OPERATOR @> (LEFTARG = integer[], RIGHTARG = integer[], PROCEDURE = intarray_contains, COMMUTATOR = <@);
CREATE OPERATOR <@ (LEFTARG = integer[], RIGHTARG = integer[], PROCEDURE = intarray_contained, COMMUTATOR = @>);
CREATE OPERATOR && (LEFTARG = integer[], RIGHTARG = integer[], PROCEDURE = intarray_overlap, COMMUTATOR = &&);
CREATE OPERATOR + (LEFTARG = integer[], RIGHTARG = integer, PROCEDURE = intarray_push_elem);
CREATE OPERATOR + (LEFTARG = integer[], RIGHTARG = integer[], PROCEDURE = intarray_push_array);
CREATE OPERATOR - (LEFTARG = integer[], RIGHTARG = integer, PROCEDURE = intarray_del_elem);
CREATE OPERATOR | (LEFTARG = integer[], RIGHTARG = integer[], PROCEDURE = intarray_union, COMMUTATOR = |);
CREATE OPERATOR & (LEFTARG = integer[], RIGHTARG = integer[], PROCEDURE = intarray_inter, COMMUTATOR = &);

CREATE FUNCTION ginint4_queryextract(integer[], internal, smallint, internal, internal, internal, internal) RETURNS internal AS 'ginint4_queryextract' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION ginint4_consistent(internal, smallint, integer[], integer, internal, internal, internal, internal) RETURNS boolean AS 'ginint4_consistent' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION g_int_consistent(internal, integer[], smallint, oid, internal) RETURNS boolean AS 'g_int_consistent' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION g_int_union(internal, internal) RETURNS integer[] AS 'g_int_union' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION g_int_compress(internal) RETURNS internal AS 'g_int_compress' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION g_int_decompress(internal) RETURNS internal AS 'g_int_decompress' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION g_int_penalty(internal, internal, internal) RETURNS internal AS 'g_int_penalty' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION g_int_picksplit(internal, internal) RETURNS internal AS 'g_int_picksplit' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION g_int_same(integer[], integer[], internal) RETURNS internal AS 'g_int_same' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;

CREATE OPERATOR CLASS gist__int_ops DEFAULT FOR TYPE integer[] USING gist AS OPERATOR 3 &&, OPERATOR 6 =, OPERATOR 7 @>, OPERATOR 8 <@, FUNCTION 1 g_int_consistent(internal, integer[], smallint, oid, internal);
CREATE OPERATOR CLASS gin__int_ops DEFAULT FOR TYPE integer[] USING gin AS OPERATOR 3 &&, OPERATOR 6 =, OPERATOR 7 @>, OPERATOR 8 <@, FUNCTION 1 ginint4_queryextract(integer[], internal, smallint, internal, internal, internal, internal);
`

// pgTrgmSQL is the DDL subset for the pg_trgm extension.
// Derived from contrib/pg_trgm/pg_trgm--1.6.sql.
const pgTrgmSQL = `
CREATE FUNCTION set_limit(real) RETURNS real AS 'set_limit' LANGUAGE C STRICT VOLATILE PARALLEL SAFE;
CREATE FUNCTION show_limit() RETURNS real AS 'show_limit' LANGUAGE C STRICT STABLE PARALLEL SAFE;
CREATE FUNCTION show_trgm(text) RETURNS text[] AS 'show_trgm' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION similarity(text, text) RETURNS real AS 'similarity' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION similarity_op(text, text) RETURNS boolean AS 'similarity_op' LANGUAGE C STRICT STABLE PARALLEL SAFE;
CREATE FUNCTION word_similarity(text, text) RETURNS real AS 'word_similarity' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION word_similarity_op(text, text) RETURNS boolean AS 'word_similarity_op' LANGUAGE C STRICT STABLE PARALLEL SAFE;
CREATE FUNCTION strict_word_similarity(text, text) RETURNS real AS 'strict_word_similarity' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION strict_word_similarity_op(text, text) RETURNS boolean AS 'strict_word_similarity_op' LANGUAGE C STRICT STABLE PARALLEL SAFE;
CREATE FUNCTION similarity_dist(text, text) RETURNS real AS 'similarity_dist' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION word_similarity_dist_op(text, text) RETURNS real AS 'word_similarity_dist_op' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION strict_word_similarity_dist_op(text, text) RETURNS real AS 'strict_word_similarity_dist_op' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;

CREATE OPERATOR % (LEFTARG = text, RIGHTARG = text, PROCEDURE = similarity_op, COMMUTATOR = %);
CREATE OPERATOR <% (LEFTARG = text, RIGHTARG = text, PROCEDURE = word_similarity_op);
CREATE OPERATOR %> (LEFTARG = text, RIGHTARG = text, PROCEDURE = word_similarity_op);
CREATE OPERATOR <<% (LEFTARG = text, RIGHTARG = text, PROCEDURE = strict_word_similarity_op);
CREATE OPERATOR %>> (LEFTARG = text, RIGHTARG = text, PROCEDURE = strict_word_similarity_op);
CREATE OPERATOR <-> (LEFTARG = text, RIGHTARG = text, PROCEDURE = similarity_dist, COMMUTATOR = <->);
CREATE OPERATOR <<-> (LEFTARG = text, RIGHTARG = text, PROCEDURE = word_similarity_dist_op);
CREATE OPERATOR <->> (LEFTARG = text, RIGHTARG = text, PROCEDURE = word_similarity_dist_op);
CREATE OPERATOR <<<-> (LEFTARG = text, RIGHTARG = text, PROCEDURE = strict_word_similarity_dist_op);
CREATE OPERATOR <->>> (LEFTARG = text, RIGHTARG = text, PROCEDURE = strict_word_similarity_dist_op);

CREATE FUNCTION gtrgm_in(cstring) RETURNS gtrgm AS 'gtrgm_in' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION gtrgm_out(gtrgm) RETURNS cstring AS 'gtrgm_out' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE TYPE gtrgm (INPUT = gtrgm_in, OUTPUT = gtrgm_out, INTERNALLENGTH = VARIABLE, STORAGE = plain);
CREATE FUNCTION gtrgm_consistent(internal, text, smallint, oid, internal) RETURNS boolean AS 'gtrgm_consistent' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION gtrgm_union(internal, internal) RETURNS gtrgm AS 'gtrgm_union' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION gtrgm_compress(internal) RETURNS internal AS 'gtrgm_compress' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION gtrgm_decompress(internal) RETURNS internal AS 'gtrgm_decompress' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION gtrgm_penalty(internal, internal, internal) RETURNS internal AS 'gtrgm_penalty' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION gtrgm_picksplit(internal, internal) RETURNS internal AS 'gtrgm_picksplit' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION gtrgm_same(gtrgm, gtrgm, internal) RETURNS internal AS 'gtrgm_same' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION gin_extract_value_trgm(text, internal) RETURNS internal AS 'gin_extract_value_trgm' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION gin_extract_query_trgm(text, internal, smallint, internal, internal, internal, internal) RETURNS internal AS 'gin_extract_query_trgm' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION gin_trgm_consistent(internal, smallint, text, integer, internal, internal, internal, internal) RETURNS boolean AS 'gin_trgm_consistent' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;

CREATE OPERATOR CLASS gist_trgm_ops DEFAULT FOR TYPE text USING gist AS OPERATOR 1 %, FUNCTION 1 gtrgm_consistent(internal, text, smallint, oid, internal), STORAGE gtrgm;
CREATE OPERATOR CLASS gin_trgm_ops FOR TYPE text USING gin AS OPERATOR 1 %, FUNCTION 1 gin_extract_value_trgm(text, internal);
`

// uuidOSSPSQL is the DDL subset for the uuid-ossp extension.
// Derived from contrib/uuid-ossp/uuid-ossp--1.1.sql.
const uuidOSSPSQL = `
CREATE FUNCTION uuid_nil() RETURNS uuid AS 'uuid_nil' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION uuid_ns_dns() RETURNS uuid AS 'uuid_ns_dns' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION uuid_ns_url() RETURNS uuid AS 'uuid_ns_url' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION uuid_ns_oid() RETURNS uuid AS 'uuid_ns_oid' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION uuid_ns_x500() RETURNS uuid AS 'uuid_ns_x500' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION uuid_generate_v1() RETURNS uuid AS 'uuid_generate_v1' LANGUAGE C STRICT VOLATILE PARALLEL SAFE;
CREATE FUNCTION uuid_generate_v1mc() RETURNS uuid AS 'uuid_generate_v1mc' LANGUAGE C STRICT VOLATILE PARALLEL SAFE;
CREATE FUNCTION uuid_generate_v3(uuid, text) RETURNS uuid AS 'uuid_generate_v3' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION uuid_generate_v4() RETURNS uuid AS 'uuid_generate_v4' LANGUAGE C STRICT VOLATILE PARALLEL SAFE;
CREATE FUNCTION uuid_generate_v5(uuid, text) RETURNS uuid AS 'uuid_generate_v5' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
`

// pgcryptoSQL is the DDL subset for the pgcrypto extension.
// Derived from contrib/pgcrypto/pgcrypto--1.3.sql.
const pgcryptoSQL = `
CREATE FUNCTION digest(bytea, text) RETURNS bytea AS 'pg_digest' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION digest(text, text) RETURNS bytea AS 'pg_digest' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION hmac(bytea, bytea, text) RETURNS bytea AS 'pg_hmac' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION hmac(text, text, text) RETURNS bytea AS 'pg_hmac' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION crypt(text, text) RETURNS text AS 'pg_crypt' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION gen_salt(text) RETURNS text AS 'pg_gen_salt' LANGUAGE C STRICT VOLATILE PARALLEL SAFE;
CREATE FUNCTION gen_salt(text, integer) RETURNS text AS 'pg_gen_salt_rounds' LANGUAGE C STRICT VOLATILE PARALLEL SAFE;
CREATE FUNCTION encrypt(bytea, bytea, text) RETURNS bytea AS 'pg_encrypt' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION decrypt(bytea, bytea, text) RETURNS bytea AS 'pg_decrypt' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION encrypt_iv(bytea, bytea, bytea, text) RETURNS bytea AS 'pg_encrypt_iv' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION decrypt_iv(bytea, bytea, bytea, text) RETURNS bytea AS 'pg_decrypt_iv' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION gen_random_bytes(integer) RETURNS bytea AS 'pg_random_bytes' LANGUAGE C STRICT VOLATILE PARALLEL SAFE;
CREATE FUNCTION gen_random_uuid() RETURNS uuid AS 'pg_random_uuid' LANGUAGE C STRICT VOLATILE PARALLEL SAFE;
CREATE FUNCTION armor(bytea) RETURNS text AS 'pg_armor' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION dearmor(text) RETURNS bytea AS 'pg_dearmor' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION pgp_sym_encrypt(text, text) RETURNS bytea AS 'pgp_sym_encrypt_text' LANGUAGE C STRICT VOLATILE PARALLEL SAFE;
CREATE FUNCTION pgp_sym_encrypt_bytea(bytea, text) RETURNS bytea AS 'pgp_sym_encrypt_bytea' LANGUAGE C STRICT VOLATILE PARALLEL SAFE;
CREATE FUNCTION pgp_sym_decrypt(bytea, text) RETURNS text AS 'pgp_sym_decrypt_text' LANGUAGE C STRICT VOLATILE PARALLEL SAFE;
CREATE FUNCTION pgp_sym_decrypt_bytea(bytea, text) RETURNS bytea AS 'pgp_sym_decrypt_bytea' LANGUAGE C STRICT VOLATILE PARALLEL SAFE;
CREATE FUNCTION pgp_pub_encrypt(text, bytea) RETURNS bytea AS 'pgp_pub_encrypt_text' LANGUAGE C STRICT VOLATILE PARALLEL SAFE;
CREATE FUNCTION pgp_pub_encrypt_bytea(bytea, bytea) RETURNS bytea AS 'pgp_pub_encrypt_bytea' LANGUAGE C STRICT VOLATILE PARALLEL SAFE;
CREATE FUNCTION pgp_pub_decrypt(bytea, bytea) RETURNS text AS 'pgp_pub_decrypt_text' LANGUAGE C STRICT VOLATILE PARALLEL SAFE;
CREATE FUNCTION pgp_pub_decrypt_bytea(bytea, bytea) RETURNS bytea AS 'pgp_pub_decrypt_bytea' LANGUAGE C STRICT VOLATILE PARALLEL SAFE;
`

// btreeGistSQL is the DDL subset for the btree_gist extension.
// Derived from contrib/btree_gist/btree_gist--1.2.sql through btree_gist--1.7.sql.
const btreeGistSQL = `
CREATE FUNCTION gbt_int4_consistent(internal, integer, smallint, oid, internal) RETURNS boolean AS 'gbt_int4_consistent' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION gbt_text_consistent(internal, text, smallint, oid, internal) RETURNS boolean AS 'gbt_text_consistent' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION gbt_ts_consistent(internal, timestamp, smallint, oid, internal) RETURNS boolean AS 'gbt_ts_consistent' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;

CREATE OPERATOR CLASS gist_int4_ops DEFAULT FOR TYPE integer USING gist AS OPERATOR 1 <, OPERATOR 2 <=, OPERATOR 3 =, OPERATOR 4 >=, OPERATOR 5 >, FUNCTION 1 gbt_int4_consistent(internal, integer, smallint, oid, internal);
CREATE OPERATOR CLASS gist_text_ops DEFAULT FOR TYPE text USING gist AS OPERATOR 1 <, OPERATOR 2 <=, OPERATOR 3 =, OPERATOR 4 >=, OPERATOR 5 >, FUNCTION 1 gbt_text_consistent(internal, text, smallint, oid, internal);
CREATE OPERATOR CLASS gist_timestamp_ops DEFAULT FOR TYPE timestamp USING gist AS OPERATOR 1 <, OPERATOR 2 <=, OPERATOR 3 =, OPERATOR 4 >=, OPERATOR 5 >, FUNCTION 1 gbt_ts_consistent(internal, timestamp, smallint, oid, internal);
`

// btreeGinSQL is the DDL subset for the btree_gin extension.
// Derived from contrib/btree_gin/btree_gin--1.0.sql through btree_gin--1.3.sql.
const btreeGinSQL = `
CREATE FUNCTION gin_btree_consistent(internal, smallint, anyelement, integer, internal, internal) RETURNS boolean AS 'gin_btree_consistent' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION gin_extract_value_int4(integer, internal) RETURNS internal AS 'gin_extract_value_int4' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION gin_extract_query_int4(integer, internal, smallint, internal, internal) RETURNS internal AS 'gin_extract_query_int4' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION gin_compare_prefix_int4(integer, integer, smallint, internal) RETURNS integer AS 'gin_compare_prefix_int4' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION gin_extract_value_text(text, internal) RETURNS internal AS 'gin_extract_value_text' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION gin_extract_query_text(text, internal, smallint, internal, internal) RETURNS internal AS 'gin_extract_query_text' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION gin_compare_prefix_text(text, text, smallint, internal) RETURNS integer AS 'gin_compare_prefix_text' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION gin_extract_value_numeric(numeric, internal) RETURNS internal AS 'gin_extract_value_numeric' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION gin_extract_query_numeric(numeric, internal, smallint, internal, internal) RETURNS internal AS 'gin_extract_query_numeric' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION gin_compare_prefix_numeric(numeric, numeric, smallint, internal) RETURNS integer AS 'gin_compare_prefix_numeric' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION gin_numeric_cmp(numeric, numeric) RETURNS integer AS 'gin_numeric_cmp' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;

CREATE OPERATOR CLASS int4_ops DEFAULT FOR TYPE integer USING gin AS OPERATOR 1 <, OPERATOR 2 <=, OPERATOR 3 =, OPERATOR 4 >=, OPERATOR 5 >, FUNCTION 2 gin_extract_value_int4(integer, internal);
CREATE OPERATOR CLASS text_ops DEFAULT FOR TYPE text USING gin AS OPERATOR 1 <, OPERATOR 2 <=, OPERATOR 3 =, OPERATOR 4 >=, OPERATOR 5 >, FUNCTION 2 gin_extract_value_text(text, internal);
CREATE OPERATOR CLASS numeric_ops DEFAULT FOR TYPE numeric USING gin AS OPERATOR 1 <, OPERATOR 2 <=, OPERATOR 3 =, OPERATOR 4 >=, OPERATOR 5 >, FUNCTION 2 gin_extract_value_numeric(numeric, internal);
`

// cubeSQL is the DDL subset for the cube extension.
// Derived from contrib/cube/cube--1.2.sql through cube--1.5.sql.
const cubeSQL = `
CREATE FUNCTION cube_in(cstring) RETURNS cube AS 'cube_in' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION cube_out(cube) RETURNS cstring AS 'cube_out' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION cube_recv(internal) RETURNS cube AS 'cube_recv' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION cube_send(cube) RETURNS bytea AS 'cube_send' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE TYPE cube (INPUT = cube_in, OUTPUT = cube_out, RECEIVE = cube_recv, SEND = cube_send, INTERNALLENGTH = VARIABLE, STORAGE = plain);

CREATE FUNCTION cube(double precision[], double precision[]) RETURNS cube AS 'cube_a_f8_f8' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION cube(double precision[]) RETURNS cube AS 'cube_a_f8' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION cube(double precision) RETURNS cube AS 'cube_f8' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION cube(double precision, double precision) RETURNS cube AS 'cube_f8_f8' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION cube(cube, double precision) RETURNS cube AS 'cube_c_f8' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION cube(cube, double precision, double precision) RETURNS cube AS 'cube_c_f8_f8' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;

CREATE FUNCTION cube_eq(cube, cube) RETURNS boolean AS 'cube_eq' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION cube_ne(cube, cube) RETURNS boolean AS 'cube_ne' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION cube_lt(cube, cube) RETURNS boolean AS 'cube_lt' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION cube_gt(cube, cube) RETURNS boolean AS 'cube_gt' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION cube_le(cube, cube) RETURNS boolean AS 'cube_le' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION cube_ge(cube, cube) RETURNS boolean AS 'cube_ge' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION cube_cmp(cube, cube) RETURNS integer AS 'cube_cmp' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION cube_contains(cube, cube) RETURNS boolean AS 'cube_contains' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION cube_contained(cube, cube) RETURNS boolean AS 'cube_contained' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION cube_overlap(cube, cube) RETURNS boolean AS 'cube_overlap' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION cube_union(cube, cube) RETURNS cube AS 'cube_union' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION cube_inter(cube, cube) RETURNS cube AS 'cube_inter' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION cube_size(cube) RETURNS double precision AS 'cube_size' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION cube_distance(cube, cube) RETURNS double precision AS 'cube_distance' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION distance_chebyshev(cube, cube) RETURNS double precision AS 'distance_chebyshev' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION distance_taxicab(cube, cube) RETURNS double precision AS 'distance_taxicab' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION cube_dim(cube) RETURNS integer AS 'cube_dim' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION cube_ll_coord(cube, integer) RETURNS double precision AS 'cube_ll_coord' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION cube_ur_coord(cube, integer) RETURNS double precision AS 'cube_ur_coord' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION cube_coord(cube, integer) RETURNS double precision AS 'cube_coord' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION cube_is_point(cube) RETURNS boolean AS 'cube_is_point' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION cube_enlarge(cube, double precision, integer) RETURNS cube AS 'cube_enlarge' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;

CREATE OPERATOR < (LEFTARG = cube, RIGHTARG = cube, PROCEDURE = cube_lt, COMMUTATOR = >, NEGATOR = >=);
CREATE OPERATOR > (LEFTARG = cube, RIGHTARG = cube, PROCEDURE = cube_gt, COMMUTATOR = <, NEGATOR = <=);
CREATE OPERATOR <= (LEFTARG = cube, RIGHTARG = cube, PROCEDURE = cube_le, COMMUTATOR = >=, NEGATOR = >);
CREATE OPERATOR >= (LEFTARG = cube, RIGHTARG = cube, PROCEDURE = cube_ge, COMMUTATOR = <=, NEGATOR = <);
CREATE OPERATOR = (LEFTARG = cube, RIGHTARG = cube, PROCEDURE = cube_eq, COMMUTATOR = =, NEGATOR = <>, HASHES, MERGES);
CREATE OPERATOR <> (LEFTARG = cube, RIGHTARG = cube, PROCEDURE = cube_ne, COMMUTATOR = <>, NEGATOR = =);
CREATE OPERATOR && (LEFTARG = cube, RIGHTARG = cube, PROCEDURE = cube_overlap, COMMUTATOR = &&);
CREATE OPERATOR @> (LEFTARG = cube, RIGHTARG = cube, PROCEDURE = cube_contains, COMMUTATOR = <@);
CREATE OPERATOR <@ (LEFTARG = cube, RIGHTARG = cube, PROCEDURE = cube_contained, COMMUTATOR = @>);
CREATE OPERATOR <#> (LEFTARG = cube, RIGHTARG = cube, PROCEDURE = distance_taxicab, COMMUTATOR = <#>);
CREATE OPERATOR <-> (LEFTARG = cube, RIGHTARG = cube, PROCEDURE = cube_distance, COMMUTATOR = <->);
CREATE OPERATOR <=> (LEFTARG = cube, RIGHTARG = cube, PROCEDURE = distance_chebyshev, COMMUTATOR = <=>);

CREATE FUNCTION g_cube_consistent(internal, cube, smallint, oid, internal) RETURNS boolean AS 'g_cube_consistent' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION g_cube_distance(internal, cube, smallint, oid, internal) RETURNS double precision AS 'g_cube_distance' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;

CREATE OPERATOR CLASS cube_ops DEFAULT FOR TYPE cube USING btree AS OPERATOR 1 <, OPERATOR 2 <=, OPERATOR 3 =, OPERATOR 4 >=, OPERATOR 5 >, FUNCTION 1 cube_cmp(cube, cube);
CREATE OPERATOR CLASS gist_cube_ops DEFAULT FOR TYPE cube USING gist AS OPERATOR 3 &&, OPERATOR 6 =, OPERATOR 7 @>, OPERATOR 8 <@, OPERATOR 17 <-> (cube, cube), FUNCTION 1 g_cube_consistent(internal, cube, smallint, oid, internal);
`

// earthdistanceSQL is the DDL subset for the earthdistance extension.
// Derived from contrib/earthdistance/earthdistance--1.1.sql through earthdistance--1.2.sql.
const earthdistanceSQL = `
CREATE DOMAIN earth AS cube;
CREATE FUNCTION earth() RETURNS double precision AS 'earth' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION sec_to_gc(double precision) RETURNS double precision AS 'sec_to_gc' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION gc_to_sec(double precision) RETURNS double precision AS 'gc_to_sec' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION ll_to_earth(double precision, double precision) RETURNS earth AS 'll_to_earth' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION latitude(earth) RETURNS double precision AS 'latitude' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION longitude(earth) RETURNS double precision AS 'longitude' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION earth_distance(earth, earth) RETURNS double precision AS 'earth_distance' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION earth_box(earth, double precision) RETURNS cube AS 'earth_box' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION geo_distance(point, point) RETURNS double precision AS 'geo_distance' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE OPERATOR <@> (LEFTARG = point, RIGHTARG = point, PROCEDURE = geo_distance, COMMUTATOR = <@>);
`

// tablefuncSQL is the DDL subset for the tablefunc extension.
// Derived from contrib/tablefunc/tablefunc--1.0.sql.
const tablefuncSQL = `
CREATE FUNCTION normal_rand(integer, double precision, double precision) RETURNS SETOF double precision AS 'normal_rand' LANGUAGE C STRICT VOLATILE PARALLEL SAFE;
CREATE FUNCTION crosstab(text) RETURNS SETOF record AS 'crosstab' LANGUAGE C STRICT STABLE PARALLEL SAFE;
CREATE TYPE tablefunc_crosstab_2 AS (row_name text, category_1 text, category_2 text);
CREATE TYPE tablefunc_crosstab_3 AS (row_name text, category_1 text, category_2 text, category_3 text);
CREATE TYPE tablefunc_crosstab_4 AS (row_name text, category_1 text, category_2 text, category_3 text, category_4 text);
CREATE FUNCTION crosstab2(text) RETURNS SETOF tablefunc_crosstab_2 AS 'crosstab' LANGUAGE C STRICT STABLE PARALLEL SAFE;
CREATE FUNCTION crosstab3(text) RETURNS SETOF tablefunc_crosstab_3 AS 'crosstab' LANGUAGE C STRICT STABLE PARALLEL SAFE;
CREATE FUNCTION crosstab4(text) RETURNS SETOF tablefunc_crosstab_4 AS 'crosstab' LANGUAGE C STRICT STABLE PARALLEL SAFE;
CREATE FUNCTION crosstab(text, integer) RETURNS SETOF record AS 'crosstab' LANGUAGE C STRICT STABLE PARALLEL SAFE;
CREATE FUNCTION crosstab(text, text) RETURNS SETOF record AS 'crosstab_hash' LANGUAGE C STRICT STABLE PARALLEL SAFE;
CREATE FUNCTION connectby(text, text, text, text, integer, text) RETURNS SETOF record AS 'connectby_text' LANGUAGE C STRICT STABLE PARALLEL SAFE;
CREATE FUNCTION connectby(text, text, text, text, integer) RETURNS SETOF record AS 'connectby_text' LANGUAGE C STRICT STABLE PARALLEL SAFE;
CREATE FUNCTION connectby(text, text, text, text, text, integer, text) RETURNS SETOF record AS 'connectby_text_serial' LANGUAGE C STRICT STABLE PARALLEL SAFE;
CREATE FUNCTION connectby(text, text, text, text, text, integer) RETURNS SETOF record AS 'connectby_text_serial' LANGUAGE C STRICT STABLE PARALLEL SAFE;
`

// unaccentSQL is the DDL subset for the unaccent extension.
// Derived from contrib/unaccent/unaccent--1.1.sql.
const unaccentSQL = `
CREATE FUNCTION unaccent(regdictionary, text) RETURNS text AS 'unaccent_dict' LANGUAGE C STRICT STABLE PARALLEL SAFE;
CREATE FUNCTION unaccent(text) RETURNS text AS 'unaccent_dict' LANGUAGE C STRICT STABLE PARALLEL SAFE;
CREATE FUNCTION unaccent_init(internal) RETURNS internal AS 'unaccent_init' LANGUAGE C PARALLEL SAFE;
CREATE FUNCTION unaccent_lexize(internal, internal, internal, internal) RETURNS internal AS 'unaccent_lexize' LANGUAGE C PARALLEL SAFE;
`

// fuzzystrmatchSQL is the DDL subset for the fuzzystrmatch extension.
// Derived from contrib/fuzzystrmatch/fuzzystrmatch--1.1.sql through fuzzystrmatch--1.2.sql.
const fuzzystrmatchSQL = `
CREATE FUNCTION levenshtein(text, text) RETURNS integer AS 'levenshtein' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION levenshtein(text, text, integer, integer, integer) RETURNS integer AS 'levenshtein_with_costs' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION levenshtein_less_equal(text, text, integer) RETURNS integer AS 'levenshtein_less_equal' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION levenshtein_less_equal(text, text, integer, integer, integer, integer) RETURNS integer AS 'levenshtein_less_equal_with_costs' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION metaphone(text, integer) RETURNS text AS 'metaphone' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION soundex(text) RETURNS text AS 'soundex' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION text_soundex(text) RETURNS text AS 'soundex' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION difference(text, text) RETURNS integer AS 'difference' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION dmetaphone(text) RETURNS text AS 'dmetaphone' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION dmetaphone_alt(text) RETURNS text AS 'dmetaphone_alt' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
`
