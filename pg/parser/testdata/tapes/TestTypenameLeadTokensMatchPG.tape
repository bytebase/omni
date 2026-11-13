-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestTypenameLeadTokensMatchPG default#1
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS abort AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "abort" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS absent AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "absent" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS absolute AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "absolute" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS access AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "access" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS action AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "action" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS add AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "add" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS admin AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "admin" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS after AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "after" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS aggregate AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "aggregate" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS all AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "all"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS all(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "all"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS all(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "all"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS all precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "all"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS all character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "all"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS all character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "all"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS all varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "all"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS all varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "all"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS all character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "all"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS all character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "all"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS also AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "also" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS alter AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "alter" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS always AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "always" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS analyse AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "analyse"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS analyse(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "analyse"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS analyse(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "analyse"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS analyse precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "analyse"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS analyse character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "analyse"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS analyse character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "analyse"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS analyse varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "analyse"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS analyse varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "analyse"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS analyse character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "analyse"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS analyse character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "analyse"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS analyze AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "analyze"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS analyze(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "analyze"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS analyze(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "analyze"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS analyze precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "analyze"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS analyze character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "analyze"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS analyze character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "analyze"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS analyze varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "analyze"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS analyze varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "analyze"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS analyze character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "analyze"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS analyze character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "analyze"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS and AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "and"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS and(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "and"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS and(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "and"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS and precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "and"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS and character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "and"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS and character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "and"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS and varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "and"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS and varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "and"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS and character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "and"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS and character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "and"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS any AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "any"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS any(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "any"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS any(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "any"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS any precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "any"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS any character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "any"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS any character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "any"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS any varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "any"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS any varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "any"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS any character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "any"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS any character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "any"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS array AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "array"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS array(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "array"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS array(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "array"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS array precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "array"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS array character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "array"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS array character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "array"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS array varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "array"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS array varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "array"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS array character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "array"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS array character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "array"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS as AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "as"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS as(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "as"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS as(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "as"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS as precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "as"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS as character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "as"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS as character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "as"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS as varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "as"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS as varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "as"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS as character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "as"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS as character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "as"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS asc AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "asc"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS asc(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "asc"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS asc(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "asc"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS asc precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "asc"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS asc character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "asc"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS asc character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "asc"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS asc varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "asc"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS asc varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "asc"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS asc character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "asc"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS asc character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "asc"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS asensitive AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "asensitive" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS assertion AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "assertion" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS assignment AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "assignment" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS asymmetric AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "asymmetric"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS asymmetric(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "asymmetric"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS asymmetric(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "asymmetric"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS asymmetric precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "asymmetric"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS asymmetric character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "asymmetric"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS asymmetric character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "asymmetric"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS asymmetric varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "asymmetric"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS asymmetric varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "asymmetric"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS asymmetric character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "asymmetric"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS asymmetric character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "asymmetric"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS at AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "at" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS atomic AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "atomic" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS attach AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "attach" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS attribute AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "attribute" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS authorization AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "authorization" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS backward AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "backward" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS before AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "before" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS begin AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "begin" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS between AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "between"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS between(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "between"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS between(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "between"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS between precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "between"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS between character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "between"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS between character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "between"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS between varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "between"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS between varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "between"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS between character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "between"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS between character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "between"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS bigint AS $$ SELECT 1 $$ LANGUAGE sql
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS binary AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "binary" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS bit AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42P13 return type mismatch in function declared to return bit
| detail: Actual return type is integer.
| where: SQL function "__omni_probe"
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS boolean AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42P13 return type mismatch in function declared to return boolean
| detail: Actual return type is integer.
| where: SQL function "__omni_probe"
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS both AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "both"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS both(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "both"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS both(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "both"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS both precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "both"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS both character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "both"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS both character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "both"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS both varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "both"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS both varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "both"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS both character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "both"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS both character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "both"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS breadth AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "breadth" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS by AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "by" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS cache AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "cache" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS call AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "call" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS called AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "called" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS cascade AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "cascade" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS cascaded AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "cascaded" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS case AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "case"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS case(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "case"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS case(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "case"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS case precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "case"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS case character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "case"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS case character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "case"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS case varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "case"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS case varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "case"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS case character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "case"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS case character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "case"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS cast AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "cast"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS cast(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "cast"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS cast(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "cast"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS cast precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "cast"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS cast character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "cast"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS cast character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "cast"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS cast varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "cast"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS cast varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "cast"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS cast character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "cast"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS cast character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "cast"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS catalog AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "catalog" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS chain AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "chain" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS char AS $$ SELECT 1 $$ LANGUAGE sql
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS character AS $$ SELECT 1 $$ LANGUAGE sql
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS characteristics AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "characteristics" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS check AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "check"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS check(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "check"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS check(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "check"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS check precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "check"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS check character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "check"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS check character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "check"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS check varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "check"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS check varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "check"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS check character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "check"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS check character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "check"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS checkpoint AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "checkpoint" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS class AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "class" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS close AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "close" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS cluster AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "cluster" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS coalesce AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "coalesce"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS coalesce(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "coalesce"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS coalesce(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "coalesce"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS coalesce precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "coalesce"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS coalesce character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "coalesce"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS coalesce character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "coalesce"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS coalesce varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "coalesce"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS coalesce varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "coalesce"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS coalesce character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "coalesce"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS coalesce character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "coalesce"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS collate AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "collate"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS collate(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "collate"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS collate(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "collate"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS collate precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "collate"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS collate character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "collate"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS collate character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "collate"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS collate varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "collate"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS collate varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "collate"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS collate character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "collate"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS collate character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "collate"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS collation AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "collation" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS column AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "column"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS column(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "column"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS column(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "column"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS column precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "column"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS column character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "column"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS column character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "column"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS column varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "column"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS column varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "column"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS column character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "column"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS column character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "column"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS columns AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "columns" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS comment AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "comment" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS comments AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "comments" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS commit AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "commit" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS committed AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "committed" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS compression AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "compression" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS concurrently AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "concurrently" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS conditional AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "conditional" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS configuration AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "configuration" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS conflict AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "conflict" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS connection AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "connection" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS constraint AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "constraint"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS constraint(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "constraint"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS constraint(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "constraint"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS constraint precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "constraint"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS constraint character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "constraint"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS constraint character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "constraint"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS constraint varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "constraint"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS constraint varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "constraint"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS constraint character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "constraint"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS constraint character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "constraint"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS constraints AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "constraints" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS content AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "content" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS continue AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "continue" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS conversion AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "conversion" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS copy AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "copy" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS cost AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "cost" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS create AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "create"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS create(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "create"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS create(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "create"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS create precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "create"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS create character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "create"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS create character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "create"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS create varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "create"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS create varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "create"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS create character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "create"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS create character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "create"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS cross AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "cross" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS csv AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "csv" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS cube AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "cube" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "current" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_catalog AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_catalog"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_catalog(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_catalog"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_catalog(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_catalog"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_catalog precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_catalog"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_catalog character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_catalog"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_catalog character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_catalog"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_catalog varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_catalog"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_catalog varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_catalog"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_catalog character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_catalog"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_catalog character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_catalog"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_date AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_date"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_date(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_date"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_date(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_date"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_date precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_date"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_date character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_date"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_date character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_date"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_date varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_date"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_date varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_date"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_date character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_date"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_date character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_date"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_role AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_role"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_role(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_role"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_role(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_role"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_role precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_role"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_role character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_role"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_role character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_role"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_role varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_role"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_role varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_role"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_role character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_role"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_role character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_role"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_schema AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "current_schema" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_time AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_time"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_time(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_time"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_time(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_time"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_time precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_time"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_time character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_time"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_time character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_time"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_time varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_time"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_time varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_time"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_time character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_time"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_time character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_time"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_timestamp AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_timestamp"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_timestamp(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_timestamp"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_timestamp(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_timestamp"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_timestamp precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_timestamp"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_timestamp character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_timestamp"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_timestamp character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_timestamp"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_timestamp varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_timestamp"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_timestamp varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_timestamp"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_timestamp character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_timestamp"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_timestamp character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_timestamp"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_user AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_user"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_user(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_user"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_user(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_user"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_user precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_user"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_user character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_user"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_user character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_user"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_user varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_user"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_user varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_user"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_user character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_user"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS current_user character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "current_user"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS cursor AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "cursor" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS cycle AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "cycle" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS data AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "data" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS database AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "database" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS day AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "day" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS deallocate AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "deallocate" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS dec AS $$ SELECT 1 $$ LANGUAGE sql
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS decimal AS $$ SELECT 1 $$ LANGUAGE sql
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS declare AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "declare" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS default AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "default"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS default(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "default"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS default(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "default"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS default precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "default"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS default character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "default"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS default character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "default"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS default varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "default"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS default varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "default"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS default character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "default"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS default character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "default"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS defaults AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "defaults" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS deferrable AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "deferrable"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS deferrable(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "deferrable"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS deferrable(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "deferrable"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS deferrable precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "deferrable"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS deferrable character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "deferrable"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS deferrable character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "deferrable"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS deferrable varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "deferrable"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS deferrable varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "deferrable"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS deferrable character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "deferrable"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS deferrable character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "deferrable"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS deferred AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "deferred" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS definer AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "definer" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS delete AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "delete" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS delimiter AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "delimiter" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS delimiters AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "delimiters" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS depends AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "depends" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS depth AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "depth" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS desc AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "desc"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS desc(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "desc"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS desc(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "desc"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS desc precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "desc"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS desc character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "desc"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS desc character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "desc"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS desc varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "desc"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS desc varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "desc"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS desc character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "desc"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS desc character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "desc"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS detach AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "detach" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS dictionary AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "dictionary" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS disable AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "disable" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS discard AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "discard" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS distinct AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "distinct"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS distinct(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "distinct"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS distinct(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "distinct"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS distinct precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "distinct"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS distinct character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "distinct"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS distinct character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "distinct"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS distinct varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "distinct"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS distinct varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "distinct"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS distinct character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "distinct"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS distinct character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "distinct"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS do AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "do"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS do(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "do"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS do(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "do"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS do precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "do"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS do character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "do"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS do character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "do"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS do varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "do"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS do varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "do"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS do character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "do"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS do character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "do"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS document AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "document" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS domain AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "domain" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS double AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "double" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS drop AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "drop" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS each AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "each" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS else AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "else"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS else(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "else"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS else(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "else"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS else precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "else"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS else character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "else"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS else character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "else"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS else varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "else"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS else varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "else"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS else character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "else"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS else character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "else"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS empty AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "empty" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS enable AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "enable" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS encoding AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "encoding" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS encrypted AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "encrypted" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS end AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "end"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS end(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "end"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS end(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "end"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS end precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "end"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS end character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "end"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS end character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "end"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS end varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "end"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS end varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "end"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS end character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "end"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS end character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "end"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS enum AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "enum" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS error AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "error" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS escape AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "escape" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS event AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "event" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS except AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "except"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS except(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "except"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS except(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "except"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS except precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "except"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS except character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "except"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS except character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "except"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS except varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "except"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS except varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "except"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS except character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "except"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS except character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "except"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS exclude AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "exclude" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS excluding AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "excluding" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS exclusive AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "exclusive" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS execute AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "execute" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS exists AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "exists"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS exists(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "exists"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS exists(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "exists"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS exists precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "exists"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS exists character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "exists"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS exists character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "exists"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS exists varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "exists"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS exists varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "exists"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS exists character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "exists"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS exists character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "exists"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS explain AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "explain" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS expression AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "expression" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS extension AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "extension" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS external AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "external" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS extract AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "extract"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS extract(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "extract"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS extract(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "extract"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS extract precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "extract"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS extract character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "extract"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS extract character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "extract"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS extract varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "extract"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS extract varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "extract"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS extract character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "extract"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS extract character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "extract"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS false AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "false"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS false(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "false"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS false(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "false"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS false precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "false"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS false character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "false"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS false character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "false"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS false varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "false"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS false varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "false"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS false character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "false"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS false character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "false"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS family AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "family" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS fetch AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "fetch"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS fetch(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "fetch"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS fetch(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "fetch"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS fetch precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "fetch"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS fetch character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "fetch"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS fetch character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "fetch"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS fetch varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "fetch"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS fetch varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "fetch"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS fetch character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "fetch"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS fetch character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "fetch"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS filter AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "filter" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS finalize AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "finalize" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS first AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "first" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS float AS $$ SELECT 1 $$ LANGUAGE sql
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS following AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "following" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS for AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "for"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS for(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "for"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS for(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "for"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS for precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "for"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS for character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "for"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS for character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "for"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS for varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "for"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS for varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "for"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS for character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "for"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS for character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "for"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS force AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "force" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS foreign AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "foreign"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS foreign(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "foreign"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS foreign(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "foreign"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS foreign precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "foreign"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS foreign character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "foreign"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS foreign character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "foreign"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS foreign varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "foreign"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS foreign varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "foreign"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS foreign character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "foreign"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS foreign character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "foreign"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS format AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "format" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS forward AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "forward" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS freeze AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "freeze" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS from AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "from"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS from(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "from"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS from(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "from"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS from precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "from"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS from character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "from"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS from character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "from"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS from varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "from"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS from varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "from"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS from character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "from"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS from character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "from"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS full AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "full" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS function AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "function" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS functions AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "functions" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS generated AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "generated" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS global AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "global" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS grant AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "grant"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS grant(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "grant"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS grant(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "grant"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS grant precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "grant"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS grant character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "grant"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS grant character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "grant"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS grant varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "grant"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS grant varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "grant"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS grant character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "grant"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS grant character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "grant"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS granted AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "granted" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS greatest AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "greatest"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS greatest(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "greatest"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS greatest(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "greatest"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS greatest precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "greatest"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS greatest character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "greatest"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS greatest character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "greatest"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS greatest varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "greatest"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS greatest varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "greatest"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS greatest character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "greatest"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS greatest character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "greatest"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS group AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "group"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS group(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "group"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS group(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "group"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS group precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "group"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS group character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "group"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS group character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "group"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS group varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "group"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS group varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "group"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS group character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "group"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS group character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "group"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS grouping AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "grouping"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS grouping(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "grouping"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS grouping(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "grouping"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS grouping precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "grouping"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS grouping character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "grouping"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS grouping character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "grouping"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS grouping varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "grouping"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS grouping varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "grouping"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS grouping character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "grouping"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS grouping character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "grouping"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS groups AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "groups" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS handler AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "handler" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS having AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "having"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS having(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "having"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS having(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "having"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS having precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "having"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS having character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "having"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS having character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "having"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS having varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "having"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS having varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "having"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS having character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "having"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS having character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "having"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS header AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "header" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS hold AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "hold" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS hour AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "hour" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS identity AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "identity" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS if AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "if" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS ilike AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "ilike" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS immediate AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "immediate" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS immutable AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "immutable" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS implicit AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "implicit" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS import AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "import" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS in AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "in"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS in(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "in"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS in(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "in"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS in precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "in"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS in character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "in"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS in character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "in"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS in varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "in"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS in varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "in"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS in character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "in"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS in character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "in"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS include AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "include" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS including AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "including" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS increment AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "increment" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS indent AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "indent" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS index AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "index" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS indexes AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "indexes" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS inherit AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "inherit" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS inherits AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "inherits" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS initially AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "initially"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS initially(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "initially"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS initially(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "initially"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS initially precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "initially"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS initially character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "initially"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS initially character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "initially"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS initially varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "initially"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS initially varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "initially"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS initially character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "initially"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS initially character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "initially"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS inline AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "inline" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS inner AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "inner" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS inout AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "inout"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS inout(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "inout"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS inout(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "inout"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS inout precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "inout"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS inout character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "inout"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS inout character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "inout"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS inout varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "inout"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS inout varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "inout"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS inout character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "inout"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS inout character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "inout"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS input AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "input" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS insensitive AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "insensitive" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS insert AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "insert" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS instead AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "instead" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS integer AS $$ SELECT 1 $$ LANGUAGE sql
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS intersect AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "intersect"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS intersect(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "intersect"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS intersect(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "intersect"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS intersect precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "intersect"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS intersect character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "intersect"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS intersect character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "intersect"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS intersect varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "intersect"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS intersect varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "intersect"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS intersect character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "intersect"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS intersect character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "intersect"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS interval AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42P13 return type mismatch in function declared to return interval
| detail: Actual return type is integer.
| where: SQL function "__omni_probe"
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS into AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "into"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS into(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "into"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS into(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "into"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS into precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "into"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS into character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "into"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS into character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "into"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS into varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "into"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS into varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "into"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS into character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "into"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS into character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "into"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS invoker AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "invoker" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS is AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "is" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS isnull AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "isnull" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS isolation AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "isolation" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS join AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "join" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42P13 return type mismatch in function declared to return json
| detail: Actual return type is integer.
| where: SQL function "__omni_probe"
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_array AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_array"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_array(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_array"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_array(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_array"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_array precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_array"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_array character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_array"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_array character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_array"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_array varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_array"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_array varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_array"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_array character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_array"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_array character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_array"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_arrayagg AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_arrayagg"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_arrayagg(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_arrayagg"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_arrayagg(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_arrayagg"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_arrayagg precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_arrayagg"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_arrayagg character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_arrayagg"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_arrayagg character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_arrayagg"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_arrayagg varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_arrayagg"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_arrayagg varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_arrayagg"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_arrayagg character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_arrayagg"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_arrayagg character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_arrayagg"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_exists AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_exists"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_exists(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_exists"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_exists(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_exists"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_exists precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_exists"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_exists character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_exists"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_exists character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_exists"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_exists varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_exists"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_exists varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_exists"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_exists character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_exists"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_exists character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_exists"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_object AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_object"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_object(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_object"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_object(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_object"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_object precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_object"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_object character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_object"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_object character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_object"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_object varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_object"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_object varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_object"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_object character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_object"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_object character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_object"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_objectagg AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_objectagg"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_objectagg(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_objectagg"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_objectagg(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_objectagg"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_objectagg precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_objectagg"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_objectagg character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_objectagg"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_objectagg character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_objectagg"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_objectagg varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_objectagg"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_objectagg varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_objectagg"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_objectagg character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_objectagg"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_objectagg character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_objectagg"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_query AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_query"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_query(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_query"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_query(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_query"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_query precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_query"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_query character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_query"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_query character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_query"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_query varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_query"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_query varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_query"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_query character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_query"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_query character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_query"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_scalar AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_scalar"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_scalar(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_scalar"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_scalar(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_scalar"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_scalar precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_scalar"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_scalar character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_scalar"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_scalar character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_scalar"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_scalar varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_scalar"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_scalar varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_scalar"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_scalar character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_scalar"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_scalar character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_scalar"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_serialize AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_serialize"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_serialize(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_serialize"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_serialize(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_serialize"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_serialize precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_serialize"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_serialize character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_serialize"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_serialize character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_serialize"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_serialize varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_serialize"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_serialize varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_serialize"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_serialize character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_serialize"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_serialize character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_serialize"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_table AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_table"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_table(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_table"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_table(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_table"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_table precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_table"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_table character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_table"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_table character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_table"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_table varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_table"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_table varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_table"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_table character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_table"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_table character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_table"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_value AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_value"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_value(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_value"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_value(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_value"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_value precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_value"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_value character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_value"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_value character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_value"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_value varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_value"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_value varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_value"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_value character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_value"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS json_value character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "json_value"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS keep AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "keep" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS key AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "key" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS keys AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "keys" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS label AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "label" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS language AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "language" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS large AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "large" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS last AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "last" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS lateral AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "lateral"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS lateral(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "lateral"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS lateral(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "lateral"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS lateral precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "lateral"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS lateral character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "lateral"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS lateral character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "lateral"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS lateral varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "lateral"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS lateral varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "lateral"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS lateral character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "lateral"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS lateral character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "lateral"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS leading AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "leading"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS leading(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "leading"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS leading(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "leading"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS leading precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "leading"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS leading character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "leading"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS leading character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "leading"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS leading varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "leading"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS leading varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "leading"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS leading character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "leading"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS leading character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "leading"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS leakproof AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "leakproof" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS least AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "least"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS least(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "least"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS least(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "least"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS least precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "least"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS least character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "least"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS least character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "least"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS least varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "least"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS least varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "least"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS least character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "least"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS least character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "least"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS left AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "left" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS level AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "level" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS like AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "like" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS limit AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "limit"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS limit(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "limit"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS limit(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "limit"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS limit precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "limit"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS limit character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "limit"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS limit character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "limit"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS limit varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "limit"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS limit varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "limit"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS limit character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "limit"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS limit character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "limit"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS listen AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "listen" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS load AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "load" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS local AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "local" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS localtime AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "localtime"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS localtime(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "localtime"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS localtime(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "localtime"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS localtime precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "localtime"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS localtime character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "localtime"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS localtime character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "localtime"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS localtime varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "localtime"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS localtime varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "localtime"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS localtime character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "localtime"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS localtime character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "localtime"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS localtimestamp AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "localtimestamp"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS localtimestamp(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "localtimestamp"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS localtimestamp(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "localtimestamp"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS localtimestamp precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "localtimestamp"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS localtimestamp character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "localtimestamp"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS localtimestamp character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "localtimestamp"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS localtimestamp varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "localtimestamp"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS localtimestamp varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "localtimestamp"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS localtimestamp character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "localtimestamp"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS localtimestamp character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "localtimestamp"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS location AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "location" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS lock AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "lock" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS locked AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "locked" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS logged AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "logged" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS mapping AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "mapping" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS match AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "match" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS matched AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "matched" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS materialized AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "materialized" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS maxvalue AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "maxvalue" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS merge AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "merge" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS merge_action AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "merge_action"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS merge_action(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "merge_action"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS merge_action(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "merge_action"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS merge_action precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "merge_action"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS merge_action character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "merge_action"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS merge_action character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "merge_action"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS merge_action varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "merge_action"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS merge_action varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "merge_action"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS merge_action character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "merge_action"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS merge_action character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "merge_action"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS method AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "method" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS minute AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "minute" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS minvalue AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "minvalue" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS mode AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "mode" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS month AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "month" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS move AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "move" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS name AS $$ SELECT 1 $$ LANGUAGE sql
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS names AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "names" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS national AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "AS"
| position: 49
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS national(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "("
| position: 48
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS national(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "("
| position: 48
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS national precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "precision"
| position: 49
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS national character AS $$ SELECT 1 $$ LANGUAGE sql
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS natural AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "natural" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS nchar AS $$ SELECT 1 $$ LANGUAGE sql
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS nested AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "nested" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS new AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "new" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS next AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "next" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS nfc AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "nfc" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS nfd AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "nfd" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS nfkc AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "nfkc" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS nfkd AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "nfkd" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS no AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "no" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS none AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "none"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS none(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "none"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS none(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "none"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS none precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "none"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS none character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "none"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS none character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "none"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS none varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "none"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS none varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "none"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS none character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "none"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS none character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "none"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS normalize AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "normalize"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS normalize(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "normalize"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS normalize(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "normalize"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS normalize precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "normalize"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS normalize character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "normalize"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS normalize character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "normalize"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS normalize varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "normalize"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS normalize varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "normalize"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS normalize character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "normalize"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS normalize character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "normalize"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS normalized AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "normalized" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS not AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "not"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS not(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "not"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS not(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "not"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS not precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "not"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS not character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "not"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS not character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "not"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS not varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "not"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS not varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "not"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS not character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "not"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS not character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "not"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS nothing AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "nothing" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS notify AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "notify" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS notnull AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "notnull" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS nowait AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "nowait" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS null AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "AS"
| position: 45
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS null(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "("
| position: 44
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS null(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "("
| position: 44
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS null precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "precision"
| position: 45
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS null character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "character"
| position: 45
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS null character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "character"
| position: 45
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS null varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "varying"
| position: 45
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS null varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "varying"
| position: 45
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS null character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "character"
| position: 45
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS null character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "character"
| position: 45
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS nullif AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "nullif"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS nullif(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "nullif"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS nullif(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "nullif"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS nullif precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "nullif"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS nullif character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "nullif"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS nullif character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "nullif"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS nullif varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "nullif"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS nullif varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "nullif"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS nullif character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "nullif"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS nullif character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "nullif"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS nulls AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "nulls" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS numeric AS $$ SELECT 1 $$ LANGUAGE sql
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS object AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "object" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS of AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "of" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS off AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "off" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS offset AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "offset"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS offset(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "offset"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS offset(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "offset"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS offset precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "offset"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS offset character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "offset"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS offset character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "offset"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS offset varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "offset"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS offset varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "offset"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS offset character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "offset"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS offset character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "offset"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS oids AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "oids" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS old AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "old" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS omit AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "omit" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS on AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "on"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS on(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "on"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS on(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "on"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS on precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "on"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS on character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "on"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS on character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "on"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS on varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "on"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS on varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "on"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS on character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "on"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS on character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "on"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS only AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "only"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS only(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "only"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS only(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "only"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS only precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "only"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS only character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "only"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS only character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "only"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS only varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "only"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS only varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "only"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS only character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "only"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS only character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "only"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS operator AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "operator" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS option AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "option" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS options AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "options" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS or AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "or"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS or(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "or"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS or(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "or"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS or precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "or"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS or character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "or"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS or character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "or"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS or varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "or"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS or varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "or"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS or character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "or"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS or character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "or"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS order AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "order"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS order(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "order"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS order(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "order"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS order precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "order"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS order character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "order"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS order character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "order"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS order varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "order"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS order varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "order"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS order character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "order"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS order character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "order"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS ordinality AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "ordinality" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS others AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "others" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS out AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "out"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS out(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "out"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS out(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "out"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS out precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "out"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS out character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "out"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS out character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "out"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS out varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "out"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS out varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "out"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS out character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "out"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS out character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "out"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS outer AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "outer" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS over AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "over" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS overlaps AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "overlaps" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS overlay AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "overlay"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS overlay(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "overlay"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS overlay(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "overlay"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS overlay precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "overlay"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS overlay character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "overlay"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS overlay character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "overlay"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS overlay varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "overlay"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS overlay varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "overlay"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS overlay character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "overlay"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS overlay character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "overlay"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS overriding AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "overriding" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS owned AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "owned" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS owner AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "owner" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS parallel AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "parallel" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS parameter AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "parameter" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS parser AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "parser" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS partial AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "partial" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS partition AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "partition" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS passing AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "passing" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS password AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "password" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS path AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42P13 return type mismatch in function declared to return path
| detail: Actual return type is integer.
| where: SQL function "__omni_probe"
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS placing AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "placing"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS placing(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "placing"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS placing(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "placing"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS placing precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "placing"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS placing character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "placing"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS placing character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "placing"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS placing varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "placing"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS placing varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "placing"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS placing character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "placing"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS placing character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "placing"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS plan AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "plan" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS plans AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "plans" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS policy AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "policy" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS position AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "position"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS position(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "position"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS position(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "position"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS position precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "position"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS position character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "position"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS position character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "position"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS position varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "position"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS position varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "position"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS position character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "position"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS position character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "position"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS preceding AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "preceding" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "precision"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS precision(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "precision"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS precision(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "precision"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS precision precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "precision"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS precision character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "precision"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS precision character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "precision"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS precision varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "precision"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS precision varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "precision"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS precision character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "precision"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS precision character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "precision"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS prepare AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "prepare" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS prepared AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "prepared" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS preserve AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "preserve" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS primary AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "primary"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS primary(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "primary"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS primary(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "primary"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS primary precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "primary"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS primary character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "primary"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS primary character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "primary"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS primary varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "primary"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS primary varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "primary"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS primary character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "primary"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS primary character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "primary"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS prior AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "prior" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS privileges AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "privileges" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS procedural AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "procedural" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS procedure AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "procedure" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS procedures AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "procedures" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS program AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "program" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS publication AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "publication" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS quote AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "quote" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS quotes AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "quotes" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS range AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "range" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS read AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "read" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS real AS $$ SELECT 1 $$ LANGUAGE sql
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS reassign AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "reassign" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS recheck AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "recheck" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS recursive AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "recursive" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS ref AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "ref" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS references AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "references"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS references(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "references"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS references(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "references"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS references precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "references"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS references character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "references"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS references character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "references"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS references varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "references"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS references varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "references"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS references character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "references"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS references character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "references"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS referencing AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "referencing" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS refresh AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "refresh" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS reindex AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "reindex" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS relative AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "relative" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS release AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "release" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS rename AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "rename" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS repeatable AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "repeatable" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS replace AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "replace" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS replica AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "replica" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS reset AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "reset" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS restart AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "restart" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS restrict AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "restrict" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS return AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "return" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS returning AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "returning"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS returning(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "returning"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS returning(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "returning"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS returning precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "returning"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS returning character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "returning"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS returning character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "returning"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS returning varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "returning"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS returning varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "returning"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS returning character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "returning"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS returning character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "returning"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS returns AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "returns" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS revoke AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "revoke" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS right AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "right" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS role AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "role" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS rollback AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "rollback" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS rollup AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "rollup" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS routine AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "routine" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS routines AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "routines" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS row AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "row"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS row(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "row"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS row(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "row"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS row precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "row"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS row character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "row"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS row character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "row"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS row varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "row"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS row varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "row"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS row character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "row"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS row character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "row"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS rows AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "rows" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS rule AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "rule" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS savepoint AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "savepoint" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS scalar AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "scalar" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS schema AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "schema" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS schemas AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "schemas" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS scroll AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "scroll" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS search AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "search" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS second AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "second" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS security AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "security" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS select AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "select"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS select(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "select"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS select(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "select"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS select precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "select"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS select character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "select"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS select character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "select"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS select varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "select"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS select varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "select"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS select character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "select"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS select character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "select"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS sequence AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "sequence" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS sequences AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "sequences" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS serializable AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "serializable" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS server AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "server" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS session AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "session" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS session_user AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "session_user"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS session_user(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "session_user"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS session_user(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "session_user"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS session_user precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "session_user"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS session_user character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "session_user"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS session_user character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "session_user"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS session_user varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "session_user"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS session_user varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "session_user"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS session_user character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "session_user"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS session_user character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "session_user"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS set AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "set" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS setof int AS $$ SELECT 1 $$ LANGUAGE sql
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS sets AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "sets" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS share AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "share" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS show AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "show" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS similar AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "similar" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS simple AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "simple" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS skip AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "skip" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS smallint AS $$ SELECT 1 $$ LANGUAGE sql
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS snapshot AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "snapshot" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS some AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "some"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS some(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "some"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS some(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "some"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS some precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "some"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS some character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "some"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS some character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "some"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS some varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "some"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS some varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "some"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS some character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "some"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS some character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "some"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS source AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "source" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS sql AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "sql" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS stable AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "stable" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS standalone AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "standalone" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS start AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "start" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS statement AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "statement" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS statistics AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "statistics" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS stdin AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "stdin" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS stdout AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "stdout" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS storage AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "storage" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS stored AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "stored" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS strict AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "strict" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS string AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "string" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS strip AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "strip" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS subscription AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "subscription" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS substring AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "substring"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS substring(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "substring"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS substring(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "substring"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS substring precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "substring"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS substring character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "substring"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS substring character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "substring"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS substring varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "substring"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS substring varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "substring"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS substring character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "substring"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS substring character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "substring"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS support AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "support" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS symmetric AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "symmetric"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS symmetric(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "symmetric"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS symmetric(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "symmetric"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS symmetric precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "symmetric"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS symmetric character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "symmetric"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS symmetric character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "symmetric"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS symmetric varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "symmetric"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS symmetric varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "symmetric"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS symmetric character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "symmetric"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS symmetric character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "symmetric"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS sysid AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "sysid" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS system AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "system" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS system_user AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "system_user"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS system_user(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "system_user"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS system_user(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "system_user"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS system_user precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "system_user"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS system_user character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "system_user"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS system_user character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "system_user"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS system_user varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "system_user"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS system_user varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "system_user"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS system_user character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "system_user"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS system_user character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "system_user"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS table AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "AS"
| position: 46
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS table(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "1"
| position: 46
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS table(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "1"
| position: 46
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS table precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "precision"
| position: 46
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS table character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "character"
| position: 46
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS table character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "character"
| position: 46
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS table varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "varying"
| position: 46
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS table varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "varying"
| position: 46
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS table character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "character"
| position: 46
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS table character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "character"
| position: 46
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS tables AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "tables" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS tablesample AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "tablesample" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS tablespace AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "tablespace" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS target AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "target" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS temp AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "temp" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS template AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "template" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS temporary AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "temporary" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS text AS $$ SELECT 1 $$ LANGUAGE sql
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS then AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "then"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS then(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "then"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS then(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "then"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS then precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "then"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS then character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "then"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS then character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "then"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS then varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "then"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS then varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "then"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS then character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "then"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS then character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "then"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS ties AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "ties" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS time AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42P13 return type mismatch in function declared to return time without time zone
| detail: Actual return type is integer.
| where: SQL function "__omni_probe"
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS timestamp AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42P13 return type mismatch in function declared to return timestamp without time zone
| detail: Actual return type is integer.
| where: SQL function "__omni_probe"
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS to AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "to"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS to(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "to"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS to(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "to"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS to precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "to"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS to character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "to"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS to character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "to"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS to varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "to"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS to varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "to"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS to character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "to"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS to character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "to"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS trailing AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "trailing"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS trailing(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "trailing"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS trailing(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "trailing"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS trailing precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "trailing"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS trailing character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "trailing"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS trailing character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "trailing"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS trailing varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "trailing"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS trailing varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "trailing"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS trailing character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "trailing"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS trailing character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "trailing"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS transaction AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "transaction" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS transform AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "transform" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS treat AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "treat"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS treat(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "treat"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS treat(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "treat"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS treat precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "treat"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS treat character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "treat"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS treat character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "treat"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS treat varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "treat"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS treat varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "treat"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS treat character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "treat"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS treat character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "treat"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS trigger AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42P13 SQL functions cannot return type trigger
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS trim AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "trim"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS trim(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "trim"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS trim(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "trim"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS trim precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "trim"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS trim character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "trim"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS trim character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "trim"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS trim varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "trim"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS trim varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "trim"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS trim character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "trim"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS trim character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "trim"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS true AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "true"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS true(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "true"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS true(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "true"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS true precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "true"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS true character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "true"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS true character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "true"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS true varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "true"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS true varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "true"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS true character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "true"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS true character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "true"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS truncate AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "truncate" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS trusted AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "trusted" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS type AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "type" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS types AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "types" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS uescape AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "uescape" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS unbounded AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "unbounded" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS uncommitted AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "uncommitted" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS unconditional AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "unconditional" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS unencrypted AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "unencrypted" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS union AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "union"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS union(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "union"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS union(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "union"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS union precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "union"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS union character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "union"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS union character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "union"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS union varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "union"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS union varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "union"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS union character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "union"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS union character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "union"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS unique AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "unique"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS unique(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "unique"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS unique(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "unique"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS unique precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "unique"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS unique character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "unique"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS unique character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "unique"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS unique varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "unique"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS unique varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "unique"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS unique character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "unique"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS unique character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "unique"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS unknown AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42P13 SQL functions cannot return type unknown
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS unlisten AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "unlisten" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS unlogged AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "unlogged" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS until AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "until" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS update AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "update" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS user AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "user"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS user(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "user"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS user(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "user"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS user precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "user"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS user character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "user"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS user character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "user"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS user varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "user"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS user varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "user"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS user character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "user"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS user character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "user"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS using AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "using"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS using(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "using"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS using(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "using"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS using precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "using"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS using character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "using"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS using character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "using"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS using varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "using"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS using varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "using"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS using character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "using"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS using character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "using"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS vacuum AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "vacuum" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS valid AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "valid" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS validate AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "validate" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS validator AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "validator" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS value AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "value" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS values AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "values"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS values(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "values"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS values(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "values"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS values precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "values"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS values character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "values"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS values character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "values"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS values varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "values"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS values varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "values"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS values character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "values"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS values character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "values"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS varchar AS $$ SELECT 1 $$ LANGUAGE sql
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS variadic AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "variadic"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS variadic(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "variadic"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS variadic(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "variadic"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS variadic precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "variadic"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS variadic character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "variadic"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS variadic character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "variadic"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS variadic varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "variadic"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS variadic varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "variadic"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS variadic character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "variadic"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS variadic character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "variadic"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "varying" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS verbose AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "verbose" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS version AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "version" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS view AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "view" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS views AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "views" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS volatile AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "volatile" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS when AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "when"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS when(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "when"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS when(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "when"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS when precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "when"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS when character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "when"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS when character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "when"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS when varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "when"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS when varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "when"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS when character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "when"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS when character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "when"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS where AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "where"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS where(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "where"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS where(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "where"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS where precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "where"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS where character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "where"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS where character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "where"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS where varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "where"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS where varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "where"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS where character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "where"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS where character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "where"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS whitespace AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "whitespace" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS window AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "window"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS window(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "window"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS window(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "window"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS window precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "window"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS window character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "window"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS window character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "window"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS window varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "window"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS window varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "window"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS window character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "window"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS window character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "window"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS with AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "with"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS with(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "with"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS with(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "with"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS with precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "with"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS with character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "with"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS with character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "with"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS with varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "with"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS with varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "with"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS with character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "with"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS with character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "with"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS within AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "within" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS without AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "without" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS work AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "work" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS wrapper AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "wrapper" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS write AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "write" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xml AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42P13 return type mismatch in function declared to return xml
| detail: Actual return type is integer.
| where: SQL function "__omni_probe"
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlattributes AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlattributes"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlattributes(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlattributes"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlattributes(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlattributes"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlattributes precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlattributes"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlattributes character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlattributes"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlattributes character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlattributes"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlattributes varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlattributes"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlattributes varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlattributes"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlattributes character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlattributes"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlattributes character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlattributes"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlconcat AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlconcat"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlconcat(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlconcat"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlconcat(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlconcat"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlconcat precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlconcat"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlconcat character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlconcat"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlconcat character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlconcat"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlconcat varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlconcat"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlconcat varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlconcat"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlconcat character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlconcat"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlconcat character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlconcat"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlelement AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlelement"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlelement(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlelement"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlelement(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlelement"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlelement precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlelement"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlelement character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlelement"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlelement character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlelement"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlelement varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlelement"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlelement varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlelement"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlelement character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlelement"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlelement character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlelement"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlexists AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlexists"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlexists(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlexists"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlexists(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlexists"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlexists precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlexists"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlexists character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlexists"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlexists character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlexists"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlexists varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlexists"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlexists varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlexists"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlexists character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlexists"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlexists character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlexists"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlforest AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlforest"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlforest(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlforest"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlforest(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlforest"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlforest precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlforest"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlforest character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlforest"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlforest character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlforest"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlforest varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlforest"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlforest varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlforest"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlforest character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlforest"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlforest character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlforest"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlnamespaces AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlnamespaces"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlnamespaces(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlnamespaces"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlnamespaces(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlnamespaces"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlnamespaces precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlnamespaces"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlnamespaces character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlnamespaces"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlnamespaces character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlnamespaces"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlnamespaces varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlnamespaces"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlnamespaces varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlnamespaces"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlnamespaces character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlnamespaces"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlnamespaces character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlnamespaces"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlparse AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlparse"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlparse(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlparse"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlparse(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlparse"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlparse precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlparse"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlparse character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlparse"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlparse character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlparse"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlparse varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlparse"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlparse varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlparse"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlparse character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlparse"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlparse character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlparse"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlpi AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlpi"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlpi(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlpi"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlpi(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlpi"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlpi precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlpi"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlpi character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlpi"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlpi character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlpi"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlpi varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlpi"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlpi varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlpi"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlpi character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlpi"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlpi character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlpi"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlroot AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlroot"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlroot(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlroot"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlroot(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlroot"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlroot precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlroot"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlroot character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlroot"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlroot character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlroot"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlroot varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlroot"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlroot varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlroot"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlroot character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlroot"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlroot character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlroot"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlserialize AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlserialize"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlserialize(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlserialize"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlserialize(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlserialize"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlserialize precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlserialize"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlserialize character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlserialize"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlserialize character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlserialize"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlserialize varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlserialize"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlserialize varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlserialize"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlserialize character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlserialize"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmlserialize character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmlserialize"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmltable AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmltable"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmltable(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmltable"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmltable(1, 1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmltable"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmltable precision AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmltable"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmltable character AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmltable"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmltable character(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmltable"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmltable varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmltable"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmltable varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmltable"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmltable character varying AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmltable"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS xmltable character varying(1) AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42601 syntax error at or near "xmltable"
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS year AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "year" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS yes AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "yes" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS zone AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "zone" does not exist
> rollback
< ok
> begin
< ok
> exec
| CREATE FUNCTION __omni_probe() RETURNS omni_probe_ident AS $$ SELECT 1 $$ LANGUAGE sql
< error
| 42704 type "omni_probe_ident" does not exist
> rollback
< ok
