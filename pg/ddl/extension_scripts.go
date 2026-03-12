package ddl

// extensionScripts contains bundled DDL for known extensions.
// Each entry is the SQL that CREATE EXTENSION <name> will execute through ProcessUtility.
//
// Scripts are derived from the official PostgreSQL extension SQL files,
// simplified to the subset relevant for DDL semantic analysis (types, casts,
// operators, functions, aggregates). Physical/index-only constructs like
// CREATE OPERATOR CLASS are included when the opclass infrastructure is present.
var extensionScripts = map[string]string{
	"citext": citextSQL,
	"hstore": hstoreSQL,
	"vector": pgvectorSQL,
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
