-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestFuncTypeLeadTokensMatchPG default#1
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name abort) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type abort does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name absent) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type absent does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name absolute) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type absolute does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name access) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type access does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name action) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type action does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name add) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type add does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name admin) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type admin does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name after) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type after does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name aggregate) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type aggregate does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name all) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "all"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name all(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "all"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name all(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "all"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name all precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "all"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name all character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "all"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name all character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "all"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name all varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "all"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name all varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "all"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name all character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "all"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name all character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "all"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name also) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type also does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name alter) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type alter does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name always) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type always does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name analyse) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "analyse"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name analyse(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "analyse"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name analyse(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "analyse"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name analyse precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "analyse"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name analyse character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "analyse"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name analyse character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "analyse"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name analyse varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "analyse"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name analyse varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "analyse"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name analyse character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "analyse"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name analyse character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "analyse"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name analyze) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "analyze"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name analyze(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "analyze"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name analyze(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "analyze"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name analyze precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "analyze"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name analyze character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "analyze"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name analyze character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "analyze"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name analyze varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "analyze"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name analyze varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "analyze"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name analyze character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "analyze"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name analyze character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "analyze"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name and) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "and"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name and(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "and"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name and(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "and"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name and precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "and"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name and character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "and"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name and character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "and"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name and varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "and"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name and varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "and"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name and character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "and"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name and character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "and"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name any) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "any"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name any(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "any"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name any(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "any"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name any precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "any"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name any character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "any"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name any character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "any"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name any varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "any"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name any varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "any"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name any character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "any"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name any character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "any"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name array) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type dummy_name[] does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name as) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "as"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name as(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "as"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name as(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "as"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name as precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "as"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name as character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "as"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name as character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "as"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name as varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "as"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name as varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "as"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name as character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "as"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name as character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "as"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name asc) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "asc"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name asc(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "asc"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name asc(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "asc"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name asc precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "asc"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name asc character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "asc"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name asc character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "asc"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name asc varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "asc"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name asc varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "asc"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name asc character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "asc"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name asc character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "asc"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name asensitive) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type asensitive does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name assertion) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type assertion does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name assignment) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type assignment does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name asymmetric) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "asymmetric"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name asymmetric(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "asymmetric"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name asymmetric(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "asymmetric"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name asymmetric precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "asymmetric"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name asymmetric character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "asymmetric"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name asymmetric character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "asymmetric"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name asymmetric varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "asymmetric"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name asymmetric varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "asymmetric"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name asymmetric character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "asymmetric"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name asymmetric character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "asymmetric"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name at) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type at does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name atomic) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type atomic does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name attach) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type attach does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name attribute) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type attribute does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name authorization) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type authorization does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name backward) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type backward does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name before) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type before does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name begin) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type begin does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name between) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "between"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name between(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "between"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name between(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "between"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name between precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "between"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name between character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "between"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name between character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "between"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name between varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "between"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name between varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "between"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name between character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "between"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name between character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "between"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name bigint) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name binary) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type binary does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name bit) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name boolean) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name both) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "both"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name both(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "both"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name both(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "both"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name both precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "both"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name both character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "both"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name both character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "both"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name both varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "both"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name both varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "both"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name both character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "both"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name both character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "both"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name breadth) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type breadth does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name by) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type by does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name cache) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type cache does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name call) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type call does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name called) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type called does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name cascade) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type cascade does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name cascaded) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type cascaded does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name case) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "case"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name case(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "case"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name case(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "case"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name case precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "case"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name case character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "case"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name case character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "case"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name case varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "case"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name case varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "case"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name case character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "case"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name case character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "case"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name cast) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "cast"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name cast(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "cast"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name cast(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "cast"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name cast precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "cast"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name cast character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "cast"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name cast character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "cast"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name cast varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "cast"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name cast varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "cast"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name cast character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "cast"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name cast character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "cast"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name catalog) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type catalog does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name chain) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type chain does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name char) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name characteristics) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type characteristics does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name check) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "check"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name check(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "check"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name check(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "check"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name check precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "check"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name check character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "check"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name check character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "check"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name check varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "check"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name check varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "check"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name check character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "check"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name check character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "check"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name checkpoint) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type checkpoint does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name class) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type class does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name close) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type close does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name cluster) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type cluster does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name coalesce) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "coalesce"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name coalesce(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "coalesce"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name coalesce(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "coalesce"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name coalesce precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "coalesce"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name coalesce character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "coalesce"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name coalesce character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "coalesce"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name coalesce varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "coalesce"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name coalesce varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "coalesce"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name coalesce character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "coalesce"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name coalesce character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "coalesce"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name collate) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "collate"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name collate(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "collate"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name collate(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "collate"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name collate precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "collate"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name collate character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "collate"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name collate character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "collate"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name collate varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "collate"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name collate varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "collate"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name collate character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "collate"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name collate character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "collate"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name collation) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type collation does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name column) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "column"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name column(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "column"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name column(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "column"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name column precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "column"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name column character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "column"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name column character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "column"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name column varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "column"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name column varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "column"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name column character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "column"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name column character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "column"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name columns) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type columns does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name comment) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type comment does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name comments) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type comments does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name commit) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type commit does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name committed) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type committed does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name compression) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type compression does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name concurrently) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type concurrently does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name conditional) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type conditional does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name configuration) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type configuration does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name conflict) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type conflict does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name connection) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type connection does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name constraint) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "constraint"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name constraint(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "constraint"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name constraint(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "constraint"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name constraint precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "constraint"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name constraint character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "constraint"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name constraint character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "constraint"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name constraint varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "constraint"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name constraint varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "constraint"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name constraint character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "constraint"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name constraint character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "constraint"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name constraints) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type constraints does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name content) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type content does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name continue) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type continue does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name conversion) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type conversion does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name copy) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type copy does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name cost) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type cost does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name create) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "create"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name create(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "create"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name create(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "create"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name create precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "create"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name create character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "create"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name create character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "create"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name create varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "create"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name create varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "create"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name create character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "create"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name create character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "create"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name cross) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type cross does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name csv) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type csv does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name cube) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type cube does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type current does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_catalog) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_catalog"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_catalog(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_catalog"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_catalog(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_catalog"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_catalog precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_catalog"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_catalog character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_catalog"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_catalog character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_catalog"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_catalog varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_catalog"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_catalog varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_catalog"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_catalog character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_catalog"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_catalog character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_catalog"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_date) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_date"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_date(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_date"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_date(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_date"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_date precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_date"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_date character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_date"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_date character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_date"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_date varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_date"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_date varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_date"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_date character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_date"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_date character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_date"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_role) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_role"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_role(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_role"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_role(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_role"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_role precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_role"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_role character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_role"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_role character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_role"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_role varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_role"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_role varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_role"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_role character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_role"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_role character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_role"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_schema) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type current_schema does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_time) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_time"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_time(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_time"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_time(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_time"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_time precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_time"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_time character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_time"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_time character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_time"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_time varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_time"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_time varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_time"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_time character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_time"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_time character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_time"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_timestamp) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_timestamp"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_timestamp(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_timestamp"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_timestamp(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_timestamp"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_timestamp precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_timestamp"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_timestamp character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_timestamp"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_timestamp character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_timestamp"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_timestamp varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_timestamp"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_timestamp varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_timestamp"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_timestamp character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_timestamp"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_timestamp character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_timestamp"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_user) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_user"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_user(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_user"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_user(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_user"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_user precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_user"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_user character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_user"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_user character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_user"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_user varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_user"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_user varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_user"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_user character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_user"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name current_user character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_user"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name cursor) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type cursor does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name cycle) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type cycle does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name data) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type data does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name database) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type database does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name day) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type day does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name deallocate) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type deallocate does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name dec) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name decimal) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name declare) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type declare does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name default) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near ")"
| position: 48
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name default(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type dummy_name does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name defaults) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type defaults does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name deferrable) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "deferrable"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name deferrable(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "deferrable"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name deferrable(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "deferrable"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name deferrable precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "deferrable"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name deferrable character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "deferrable"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name deferrable character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "deferrable"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name deferrable varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "deferrable"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name deferrable varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "deferrable"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name deferrable character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "deferrable"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name deferrable character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "deferrable"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name deferred) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type deferred does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name definer) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type definer does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name delete) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type delete does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name delimiter) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type delimiter does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name delimiters) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type delimiters does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name depends) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type depends does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name depth) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type depth does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name desc) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "desc"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name desc(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "desc"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name desc(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "desc"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name desc precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "desc"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name desc character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "desc"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name desc character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "desc"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name desc varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "desc"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name desc varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "desc"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name desc character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "desc"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name desc character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "desc"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name detach) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type detach does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name dictionary) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type dictionary does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name disable) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type disable does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name discard) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type discard does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name distinct) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "distinct"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name distinct(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "distinct"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name distinct(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "distinct"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name distinct precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "distinct"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name distinct character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "distinct"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name distinct character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "distinct"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name distinct varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "distinct"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name distinct varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "distinct"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name distinct character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "distinct"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name distinct character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "distinct"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name do) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "do"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name do(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "do"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name do(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "do"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name do precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "do"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name do character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "do"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name do character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "do"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name do varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "do"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name do varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "do"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name do character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "do"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name do character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "do"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name document) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type document does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name domain) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type domain does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name double) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type double does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name drop) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type drop does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name each) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type each does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name else) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "else"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name else(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "else"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name else(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "else"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name else precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "else"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name else character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "else"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name else character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "else"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name else varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "else"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name else varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "else"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name else character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "else"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name else character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "else"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name empty) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type empty does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name enable) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type enable does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name encoding) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type encoding does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name encrypted) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type encrypted does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name end) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "end"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name end(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "end"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name end(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "end"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name end precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "end"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name end character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "end"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name end character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "end"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name end varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "end"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name end varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "end"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name end character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "end"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name end character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "end"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name enum) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type enum does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name error) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type error does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name escape) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type escape does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name event) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type event does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name except) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "except"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name except(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "except"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name except(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "except"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name except precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "except"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name except character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "except"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name except character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "except"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name except varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "except"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name except varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "except"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name except character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "except"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name except character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "except"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name exclude) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type exclude does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name excluding) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type excluding does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name exclusive) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type exclusive does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name execute) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type execute does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name exists) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "exists"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name exists(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "exists"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name exists(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "exists"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name exists precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "exists"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name exists character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "exists"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name exists character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "exists"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name exists varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "exists"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name exists varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "exists"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name exists character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "exists"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name exists character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "exists"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name explain) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type explain does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name expression) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type expression does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name extension) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type extension does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name external) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type external does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name extract) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "extract"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name extract(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "extract"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name extract(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "extract"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name extract precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "extract"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name extract character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "extract"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name extract character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "extract"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name extract varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "extract"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name extract varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "extract"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name extract character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "extract"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name extract character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "extract"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name false) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "false"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name false(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "false"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name false(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "false"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name false precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "false"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name false character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "false"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name false character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "false"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name false varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "false"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name false varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "false"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name false character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "false"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name false character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "false"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name family) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type family does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name fetch) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "fetch"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name fetch(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "fetch"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name fetch(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "fetch"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name fetch precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "fetch"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name fetch character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "fetch"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name fetch character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "fetch"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name fetch varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "fetch"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name fetch varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "fetch"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name fetch character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "fetch"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name fetch character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "fetch"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name filter) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type filter does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name finalize) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type finalize does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name first) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type first does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name float) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name following) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type following does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name for) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "for"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name for(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "for"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name for(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "for"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name for precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "for"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name for character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "for"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name for character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "for"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name for varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "for"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name for varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "for"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name for character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "for"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name for character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "for"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name force) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type force does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name foreign) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "foreign"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name foreign(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "foreign"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name foreign(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "foreign"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name foreign precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "foreign"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name foreign character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "foreign"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name foreign character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "foreign"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name foreign varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "foreign"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name foreign varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "foreign"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name foreign character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "foreign"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name foreign character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "foreign"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name format) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type format does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name forward) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type forward does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name freeze) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type freeze does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name from) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "from"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name from(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "from"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name from(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "from"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name from precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "from"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name from character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "from"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name from character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "from"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name from varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "from"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name from varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "from"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name from character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "from"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name from character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "from"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name full) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type full does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name function) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type function does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name functions) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type functions does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name generated) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type generated does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name global) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type global does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name grant) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "grant"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name grant(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "grant"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name grant(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "grant"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name grant precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "grant"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name grant character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "grant"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name grant character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "grant"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name grant varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "grant"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name grant varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "grant"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name grant character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "grant"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name grant character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "grant"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name granted) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type granted does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name greatest) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "greatest"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name greatest(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "greatest"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name greatest(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "greatest"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name greatest precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "greatest"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name greatest character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "greatest"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name greatest character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "greatest"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name greatest varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "greatest"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name greatest varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "greatest"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name greatest character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "greatest"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name greatest character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "greatest"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name group) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "group"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name group(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "group"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name group(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "group"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name group precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "group"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name group character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "group"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name group character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "group"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name group varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "group"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name group varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "group"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name group character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "group"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name group character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "group"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name grouping) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "grouping"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name grouping(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "grouping"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name grouping(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "grouping"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name grouping precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "grouping"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name grouping character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "grouping"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name grouping character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "grouping"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name grouping varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "grouping"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name grouping varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "grouping"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name grouping character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "grouping"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name grouping character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "grouping"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name groups) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type groups does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name handler) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type handler does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name having) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "having"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name having(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "having"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name having(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "having"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name having precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "having"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name having character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "having"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name having character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "having"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name having varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "having"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name having varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "having"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name having character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "having"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name having character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "having"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name header) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type header does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name hold) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type hold does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name hour) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type hour does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name identity) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type identity does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name if) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type if does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name ilike) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type ilike does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name immediate) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type immediate does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name immutable) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type immutable does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name implicit) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type implicit does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name import) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type import does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name in) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near ")"
| position: 43
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name in(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "("
| position: 43
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name in(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "("
| position: 43
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name in precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "precision"
| position: 44
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name in character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name include) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type include does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name including) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type including does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name increment) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type increment does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name indent) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type indent does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name index) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type index does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name indexes) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type indexes does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name inherit) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type inherit does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name inherits) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type inherits does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name initially) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "initially"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name initially(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "initially"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name initially(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "initially"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name initially precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "initially"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name initially character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "initially"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name initially character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "initially"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name initially varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "initially"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name initially varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "initially"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name initially character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "initially"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name initially character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "initially"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name inline) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type inline does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name inner) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type inner does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name inout) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near ")"
| position: 46
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name inout(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "("
| position: 46
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name inout(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "("
| position: 46
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name inout precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "precision"
| position: 47
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name inout character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42P13 function result type must be character because of OUT parameters
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name input) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type input does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name insensitive) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type insensitive does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name insert) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type insert does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name instead) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type instead does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name int) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name integer) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name intersect) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "intersect"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name intersect(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "intersect"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name intersect(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "intersect"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name intersect precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "intersect"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name intersect character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "intersect"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name intersect character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "intersect"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name intersect varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "intersect"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name intersect varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "intersect"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name intersect character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "intersect"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name intersect character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "intersect"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name interval) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name into) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "into"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name into(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "into"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name into(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "into"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name into precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "into"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name into character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "into"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name into character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "into"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name into varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "into"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name into varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "into"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name into character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "into"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name into character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "into"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name invoker) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type invoker does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name is) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type is does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name isnull) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type isnull does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name isolation) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type isolation does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name join) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type join does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_array) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_array"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_array(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_array"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_array(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_array"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_array precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_array"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_array character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_array"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_array character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_array"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_array varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_array"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_array varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_array"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_array character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_array"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_array character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_array"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_arrayagg) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_arrayagg"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_arrayagg(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_arrayagg"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_arrayagg(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_arrayagg"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_arrayagg precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_arrayagg"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_arrayagg character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_arrayagg"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_arrayagg character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_arrayagg"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_arrayagg varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_arrayagg"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_arrayagg varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_arrayagg"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_arrayagg character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_arrayagg"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_arrayagg character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_arrayagg"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_exists) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_exists"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_exists(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_exists"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_exists(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_exists"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_exists precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_exists"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_exists character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_exists"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_exists character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_exists"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_exists varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_exists"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_exists varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_exists"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_exists character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_exists"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_exists character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_exists"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_object) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_object"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_object(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_object"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_object(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_object"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_object precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_object"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_object character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_object"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_object character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_object"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_object varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_object"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_object varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_object"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_object character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_object"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_object character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_object"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_objectagg) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_objectagg"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_objectagg(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_objectagg"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_objectagg(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_objectagg"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_objectagg precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_objectagg"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_objectagg character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_objectagg"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_objectagg character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_objectagg"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_objectagg varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_objectagg"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_objectagg varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_objectagg"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_objectagg character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_objectagg"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_objectagg character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_objectagg"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_query) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_query"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_query(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_query"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_query(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_query"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_query precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_query"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_query character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_query"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_query character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_query"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_query varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_query"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_query varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_query"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_query character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_query"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_query character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_query"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_scalar) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_scalar"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_scalar(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_scalar"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_scalar(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_scalar"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_scalar precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_scalar"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_scalar character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_scalar"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_scalar character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_scalar"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_scalar varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_scalar"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_scalar varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_scalar"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_scalar character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_scalar"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_scalar character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_scalar"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_serialize) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_serialize"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_serialize(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_serialize"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_serialize(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_serialize"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_serialize precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_serialize"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_serialize character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_serialize"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_serialize character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_serialize"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_serialize varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_serialize"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_serialize varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_serialize"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_serialize character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_serialize"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_serialize character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_serialize"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_table) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_table"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_table(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_table"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_table(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_table"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_table precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_table"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_table character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_table"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_table character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_table"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_table varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_table"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_table varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_table"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_table character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_table"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_table character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_table"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_value) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_value"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_value(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_value"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_value(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_value"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_value precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_value"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_value character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_value"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_value character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_value"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_value varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_value"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_value varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_value"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_value character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_value"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name json_value character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_value"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name keep) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type keep does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name key) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type key does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name keys) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type keys does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name label) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type label does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name language) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type language does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name large) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type large does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name last) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type last does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name lateral) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "lateral"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name lateral(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "lateral"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name lateral(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "lateral"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name lateral precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "lateral"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name lateral character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "lateral"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name lateral character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "lateral"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name lateral varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "lateral"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name lateral varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "lateral"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name lateral character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "lateral"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name lateral character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "lateral"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name leading) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "leading"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name leading(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "leading"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name leading(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "leading"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name leading precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "leading"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name leading character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "leading"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name leading character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "leading"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name leading varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "leading"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name leading varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "leading"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name leading character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "leading"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name leading character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "leading"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name leakproof) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type leakproof does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name least) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "least"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name least(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "least"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name least(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "least"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name least precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "least"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name least character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "least"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name least character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "least"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name least varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "least"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name least varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "least"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name least character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "least"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name least character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "least"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name left) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type left does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name level) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type level does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name like) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type like does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name limit) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "limit"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name limit(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "limit"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name limit(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "limit"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name limit precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "limit"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name limit character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "limit"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name limit character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "limit"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name limit varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "limit"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name limit varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "limit"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name limit character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "limit"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name limit character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "limit"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name listen) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type listen does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name load) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type load does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name local) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type local does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name localtime) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "localtime"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name localtime(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "localtime"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name localtime(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "localtime"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name localtime precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "localtime"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name localtime character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "localtime"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name localtime character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "localtime"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name localtime varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "localtime"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name localtime varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "localtime"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name localtime character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "localtime"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name localtime character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "localtime"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name localtimestamp) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "localtimestamp"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name localtimestamp(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "localtimestamp"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name localtimestamp(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "localtimestamp"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name localtimestamp precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "localtimestamp"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name localtimestamp character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "localtimestamp"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name localtimestamp character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "localtimestamp"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name localtimestamp varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "localtimestamp"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name localtimestamp varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "localtimestamp"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name localtimestamp character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "localtimestamp"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name localtimestamp character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "localtimestamp"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name location) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type location does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name lock) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type lock does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name locked) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type locked does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name logged) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type logged does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name mapping) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type mapping does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name match) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type match does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name matched) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type matched does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name materialized) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type materialized does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name maxvalue) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type maxvalue does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name merge) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type merge does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name merge_action) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "merge_action"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name merge_action(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "merge_action"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name merge_action(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "merge_action"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name merge_action precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "merge_action"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name merge_action character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "merge_action"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name merge_action character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "merge_action"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name merge_action varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "merge_action"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name merge_action varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "merge_action"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name merge_action character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "merge_action"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name merge_action character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "merge_action"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name method) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type method does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name minute) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type minute does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name minvalue) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type minvalue does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name mode) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type mode does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name month) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type month does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name move) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type move does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name name) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name names) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type names does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name national) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near ")"
| position: 49
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name national(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "("
| position: 49
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name national(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "("
| position: 49
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name national precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "precision"
| position: 50
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name national character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name natural) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type natural does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name nchar) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name nested) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type nested does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name new) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type new does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name next) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type next does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name nfc) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type nfc does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name nfd) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type nfd does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name nfkc) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type nfkc does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name nfkd) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type nfkd does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name no) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type no does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name none) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "none"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name none(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "none"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name none(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "none"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name none precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "none"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name none character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "none"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name none character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "none"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name none varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "none"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name none varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "none"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name none character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "none"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name none character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "none"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name normalize) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "normalize"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name normalize(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "normalize"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name normalize(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "normalize"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name normalize precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "normalize"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name normalize character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "normalize"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name normalize character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "normalize"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name normalize varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "normalize"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name normalize varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "normalize"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name normalize character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "normalize"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name normalize character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "normalize"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name normalized) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type normalized does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name not) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "not"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name not(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "not"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name not(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "not"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name not precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "not"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name not character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "not"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name not character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "not"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name not varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "not"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name not varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "not"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name not character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "not"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name not character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "not"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name nothing) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type nothing does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name notify) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type notify does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name notnull) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type notnull does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name nowait) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type nowait does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name null) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "null"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name null(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "null"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name null(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "null"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name null precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "null"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name null character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "null"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name null character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "null"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name null varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "null"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name null varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "null"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name null character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "null"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name null character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "null"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name nullif) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "nullif"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name nullif(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "nullif"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name nullif(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "nullif"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name nullif precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "nullif"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name nullif character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "nullif"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name nullif character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "nullif"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name nullif varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "nullif"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name nullif varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "nullif"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name nullif character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "nullif"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name nullif character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "nullif"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name nulls) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type nulls does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name numeric) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name object) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type object does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name of) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type of does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name off) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type off does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name offset) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "offset"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name offset(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "offset"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name offset(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "offset"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name offset precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "offset"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name offset character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "offset"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name offset character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "offset"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name offset varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "offset"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name offset varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "offset"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name offset character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "offset"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name offset character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "offset"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name oids) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type oids does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name old) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type old does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name omit) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type omit does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name on) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "on"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name on(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "on"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name on(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "on"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name on precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "on"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name on character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "on"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name on character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "on"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name on varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "on"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name on varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "on"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name on character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "on"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name on character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "on"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name only) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "only"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name only(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "only"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name only(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "only"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name only precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "only"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name only character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "only"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name only character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "only"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name only varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "only"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name only varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "only"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name only character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "only"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name only character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "only"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name operator) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type operator does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name option) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type option does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name options) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type options does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name or) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "or"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name or(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "or"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name or(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "or"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name or precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "or"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name or character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "or"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name or character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "or"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name or varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "or"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name or varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "or"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name or character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "or"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name or character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "or"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name order) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "order"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name order(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "order"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name order(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "order"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name order precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "order"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name order character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "order"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name order character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "order"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name order varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "order"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name order varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "order"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name order character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "order"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name order character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "order"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name ordinality) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type ordinality does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name others) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type others does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name out) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near ")"
| position: 44
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name out(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "("
| position: 44
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name out(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "("
| position: 44
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name out precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "precision"
| position: 45
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name out character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42P13 function result type must be character because of OUT parameters
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name outer) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type outer does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name over) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type over does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name overlaps) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type overlaps does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name overlay) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "overlay"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name overlay(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "overlay"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name overlay(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "overlay"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name overlay precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "overlay"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name overlay character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "overlay"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name overlay character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "overlay"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name overlay varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "overlay"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name overlay varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "overlay"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name overlay character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "overlay"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name overlay character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "overlay"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name overriding) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type overriding does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name owned) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type owned does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name owner) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type owner does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name parallel) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type parallel does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name parameter) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type parameter does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name parser) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type parser does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name partial) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type partial does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name partition) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type partition does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name passing) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type passing does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name password) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type password does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name path) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name placing) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "placing"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name placing(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "placing"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name placing(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "placing"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name placing precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "placing"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name placing character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "placing"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name placing character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "placing"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name placing varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "placing"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name placing varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "placing"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name placing character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "placing"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name placing character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "placing"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name plan) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type plan does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name plans) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type plans does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name policy) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type policy does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name position) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "position"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name position(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "position"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name position(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "position"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name position precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "position"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name position character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "position"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name position character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "position"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name position varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "position"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name position varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "position"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name position character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "position"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name position character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "position"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name preceding) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type preceding does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "precision"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name precision(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "precision"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name precision(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "precision"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name precision precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "precision"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name precision character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "precision"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name precision character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "precision"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name precision varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "precision"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name precision varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "precision"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name precision character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "precision"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name precision character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "precision"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name prepare) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type prepare does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name prepared) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type prepared does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name preserve) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type preserve does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name primary) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "primary"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name primary(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "primary"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name primary(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "primary"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name primary precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "primary"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name primary character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "primary"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name primary character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "primary"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name primary varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "primary"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name primary varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "primary"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name primary character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "primary"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name primary character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "primary"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name prior) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type prior does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name privileges) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type privileges does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name procedural) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type procedural does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name procedure) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type procedure does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name procedures) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type procedures does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name program) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type program does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name publication) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type publication does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name quote) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type quote does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name quotes) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type quotes does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name range) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type range does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name read) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type read does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name real) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name reassign) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type reassign does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name recheck) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type recheck does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name recursive) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type recursive does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name ref) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type ref does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name references) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "references"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name references(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "references"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name references(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "references"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name references precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "references"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name references character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "references"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name references character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "references"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name references varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "references"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name references varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "references"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name references character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "references"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name references character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "references"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name referencing) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type referencing does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name refresh) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type refresh does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name reindex) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type reindex does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name relative) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type relative does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name release) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type release does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name rename) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type rename does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name repeatable) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type repeatable does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name replace) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type replace does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name replica) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type replica does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name reset) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type reset does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name restart) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type restart does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name restrict) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type restrict does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name return) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type return does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name returning) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "returning"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name returning(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "returning"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name returning(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "returning"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name returning precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "returning"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name returning character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "returning"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name returning character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "returning"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name returning varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "returning"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name returning varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "returning"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name returning character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "returning"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name returning character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "returning"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name returns) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type returns does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name revoke) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type revoke does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name right) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type right does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name role) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type role does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name rollback) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type rollback does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name rollup) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type rollup does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name routine) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type routine does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name routines) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type routines does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name row) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "row"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name row(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "row"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name row(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "row"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name row precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "row"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name row character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "row"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name row character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "row"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name row varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "row"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name row varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "row"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name row character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "row"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name row character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "row"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name rows) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type rows does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name rule) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type rule does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name savepoint) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type savepoint does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name scalar) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type scalar does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name schema) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type schema does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name schemas) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type schemas does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name scroll) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type scroll does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name search) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type search does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name second) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type second does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name security) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type security does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name select) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "select"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name select(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "select"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name select(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "select"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name select precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "select"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name select character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "select"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name select character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "select"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name select varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "select"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name select varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "select"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name select character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "select"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name select character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "select"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name sequence) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type sequence does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name sequences) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type sequences does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name serializable) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type serializable does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name server) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type server does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name session) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type session does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name session_user) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "session_user"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name session_user(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "session_user"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name session_user(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "session_user"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name session_user precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "session_user"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name session_user character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "session_user"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name session_user character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "session_user"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name session_user varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "session_user"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name session_user varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "session_user"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name session_user character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "session_user"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name session_user character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "session_user"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name set) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type set does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name setof int) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42P13 functions cannot accept set arguments
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name sets) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type sets does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name share) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type share does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name show) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type show does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name similar) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type similar does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name simple) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type simple does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name skip) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type skip does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name smallint) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name snapshot) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type snapshot does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name some) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "some"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name some(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "some"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name some(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "some"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name some precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "some"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name some character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "some"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name some character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "some"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name some varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "some"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name some varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "some"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name some character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "some"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name some character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "some"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name source) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type source does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name sql) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type sql does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name stable) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type stable does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name standalone) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type standalone does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name start) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type start does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name statement) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type statement does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name statistics) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type statistics does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name stdin) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type stdin does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name stdout) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type stdout does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name storage) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type storage does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name stored) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type stored does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name strict) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type strict does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name string) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type string does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name strip) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type strip does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name subscription) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type subscription does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name substring) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "substring"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name substring(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "substring"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name substring(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "substring"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name substring precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "substring"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name substring character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "substring"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name substring character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "substring"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name substring varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "substring"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name substring varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "substring"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name substring character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "substring"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name substring character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "substring"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name support) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type support does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name symmetric) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "symmetric"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name symmetric(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "symmetric"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name symmetric(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "symmetric"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name symmetric precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "symmetric"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name symmetric character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "symmetric"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name symmetric character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "symmetric"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name symmetric varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "symmetric"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name symmetric varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "symmetric"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name symmetric character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "symmetric"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name symmetric character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "symmetric"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name sysid) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type sysid does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name system) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type system does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name system_user) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "system_user"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name system_user(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "system_user"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name system_user(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "system_user"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name system_user precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "system_user"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name system_user character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "system_user"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name system_user character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "system_user"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name system_user varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "system_user"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name system_user varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "system_user"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name system_user character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "system_user"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name system_user character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "system_user"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name table) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "table"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name table(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "table"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name table(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "table"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name table precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "table"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name table character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "table"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name table character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "table"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name table varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "table"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name table varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "table"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name table character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "table"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name table character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "table"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name tables) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type tables does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name tablesample) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type tablesample does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name tablespace) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type tablespace does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name target) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type target does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name temp) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type temp does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name template) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type template does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name temporary) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type temporary does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name text) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name then) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "then"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name then(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "then"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name then(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "then"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name then precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "then"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name then character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "then"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name then character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "then"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name then varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "then"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name then varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "then"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name then character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "then"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name then character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "then"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name ties) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type ties does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name time) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name timestamp) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name to) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "to"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name to(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "to"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name to(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "to"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name to precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "to"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name to character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "to"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name to character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "to"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name to varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "to"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name to varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "to"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name to character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "to"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name to character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "to"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name trailing) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "trailing"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name trailing(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "trailing"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name trailing(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "trailing"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name trailing precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "trailing"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name trailing character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "trailing"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name trailing character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "trailing"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name trailing varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "trailing"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name trailing varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "trailing"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name trailing character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "trailing"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name trailing character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "trailing"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name transaction) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type transaction does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name transform) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type transform does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name treat) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "treat"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name treat(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "treat"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name treat(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "treat"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name treat precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "treat"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name treat character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "treat"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name treat character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "treat"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name treat varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "treat"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name treat varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "treat"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name treat character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "treat"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name treat character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "treat"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name trigger) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42P13 SQL functions cannot have arguments of type trigger
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name trim) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "trim"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name trim(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "trim"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name trim(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "trim"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name trim precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "trim"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name trim character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "trim"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name trim character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "trim"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name trim varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "trim"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name trim varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "trim"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name trim character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "trim"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name trim character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "trim"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name true) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "true"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name true(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "true"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name true(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "true"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name true precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "true"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name true character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "true"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name true character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "true"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name true varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "true"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name true varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "true"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name true character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "true"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name true character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "true"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name truncate) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type truncate does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name trusted) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type trusted does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name type) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type type does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name types) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type types does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name uescape) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type uescape does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name unbounded) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type unbounded does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name uncommitted) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type uncommitted does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name unconditional) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type unconditional does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name unencrypted) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type unencrypted does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name union) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "union"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name union(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "union"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name union(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "union"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name union precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "union"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name union character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "union"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name union character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "union"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name union varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "union"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name union varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "union"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name union character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "union"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name union character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "union"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name unique) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "unique"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name unique(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "unique"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name unique(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "unique"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name unique precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "unique"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name unique character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "unique"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name unique character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "unique"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name unique varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "unique"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name unique varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "unique"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name unique character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "unique"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name unique character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "unique"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name unknown) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42P13 SQL functions cannot have arguments of type unknown
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name unlisten) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type unlisten does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name unlogged) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type unlogged does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name until) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type until does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name update) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type update does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name user) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "user"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name user(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "user"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name user(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "user"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name user precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "user"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name user character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "user"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name user character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "user"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name user varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "user"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name user varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "user"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name user character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "user"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name user character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "user"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name using) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "using"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name using(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "using"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name using(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "using"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name using precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "using"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name using character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "using"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name using character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "using"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name using varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "using"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name using varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "using"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name using character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "using"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name using character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "using"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name vacuum) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type vacuum does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name valid) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type valid does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name validate) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type validate does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name validator) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type validator does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name value) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type value does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name values) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "values"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name values(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "values"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name values(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "values"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name values precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "values"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name values character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "values"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name values character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "values"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name values varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "values"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name values varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "values"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name values character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "values"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name values character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "values"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name varchar) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name variadic) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near ")"
| position: 49
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name variadic(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "("
| position: 49
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name variadic(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "("
| position: 49
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name variadic precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "precision"
| position: 50
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name variadic character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42P13 VARIADIC parameter must be an array
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type varying does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name verbose) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type verbose does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name version) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type version does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name view) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type view does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name views) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type views does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name volatile) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type volatile does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name when) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "when"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name when(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "when"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name when(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "when"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name when precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "when"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name when character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "when"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name when character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "when"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name when varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "when"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name when varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "when"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name when character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "when"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name when character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "when"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name where) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "where"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name where(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "where"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name where(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "where"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name where precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "where"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name where character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "where"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name where character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "where"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name where varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "where"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name where varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "where"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name where character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "where"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name where character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "where"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name whitespace) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type whitespace does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name window) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "window"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name window(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "window"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name window(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "window"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name window precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "window"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name window character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "window"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name window character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "window"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name window varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "window"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name window varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "window"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name window character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "window"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name window character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "window"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name with) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "with"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name with(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "with"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name with(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "with"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name with precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "with"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name with character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "with"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name with character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "with"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name with varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "with"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name with varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "with"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name with character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "with"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name with character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "with"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name within) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type within does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name without) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type without does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name work) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type work does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name wrapper) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type wrapper does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name write) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type write does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xml) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlattributes) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlattributes"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlattributes(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlattributes"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlattributes(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlattributes"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlattributes precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlattributes"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlattributes character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlattributes"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlattributes character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlattributes"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlattributes varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlattributes"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlattributes varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlattributes"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlattributes character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlattributes"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlattributes character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlattributes"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlconcat) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlconcat"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlconcat(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlconcat"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlconcat(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlconcat"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlconcat precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlconcat"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlconcat character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlconcat"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlconcat character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlconcat"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlconcat varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlconcat"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlconcat varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlconcat"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlconcat character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlconcat"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlconcat character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlconcat"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlelement) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlelement"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlelement(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlelement"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlelement(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlelement"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlelement precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlelement"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlelement character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlelement"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlelement character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlelement"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlelement varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlelement"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlelement varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlelement"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlelement character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlelement"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlelement character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlelement"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlexists) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlexists"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlexists(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlexists"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlexists(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlexists"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlexists precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlexists"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlexists character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlexists"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlexists character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlexists"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlexists varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlexists"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlexists varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlexists"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlexists character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlexists"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlexists character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlexists"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlforest) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlforest"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlforest(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlforest"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlforest(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlforest"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlforest precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlforest"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlforest character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlforest"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlforest character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlforest"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlforest varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlforest"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlforest varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlforest"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlforest character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlforest"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlforest character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlforest"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlnamespaces) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlnamespaces"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlnamespaces(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlnamespaces"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlnamespaces(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlnamespaces"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlnamespaces precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlnamespaces"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlnamespaces character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlnamespaces"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlnamespaces character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlnamespaces"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlnamespaces varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlnamespaces"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlnamespaces varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlnamespaces"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlnamespaces character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlnamespaces"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlnamespaces character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlnamespaces"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlparse) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlparse"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlparse(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlparse"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlparse(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlparse"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlparse precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlparse"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlparse character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlparse"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlparse character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlparse"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlparse varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlparse"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlparse varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlparse"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlparse character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlparse"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlparse character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlparse"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlpi) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlpi"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlpi(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlpi"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlpi(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlpi"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlpi precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlpi"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlpi character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlpi"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlpi character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlpi"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlpi varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlpi"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlpi varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlpi"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlpi character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlpi"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlpi character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlpi"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlroot) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlroot"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlroot(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlroot"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlroot(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlroot"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlroot precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlroot"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlroot character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlroot"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlroot character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlroot"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlroot varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlroot"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlroot varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlroot"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlroot character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlroot"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlroot character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlroot"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlserialize) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlserialize"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlserialize(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlserialize"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlserialize(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlserialize"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlserialize precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlserialize"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlserialize character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlserialize"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlserialize character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlserialize"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlserialize varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlserialize"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlserialize varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlserialize"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlserialize character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlserialize"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmlserialize character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlserialize"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmltable) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmltable"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmltable(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmltable"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmltable(1, 1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmltable"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmltable precision) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmltable"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmltable character) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmltable"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmltable character(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmltable"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmltable varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmltable"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmltable varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmltable"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmltable character varying) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmltable"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name xmltable character varying(1)) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmltable"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name year) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type year does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name yes) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type yes does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name zone) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type zone does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe(dummy_name omni_probe_ident) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type omni_probe_ident does not exist
> rollback
< ok
