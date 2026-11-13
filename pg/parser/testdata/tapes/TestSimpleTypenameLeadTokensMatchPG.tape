-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestSimpleTypenameLeadTokensMatchPG default#1
> begin
< ok
> exec
| SELECT CAST(NULL AS abort)
< error
| 42704 type "abort" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS absent)
< error
| 42704 type "absent" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS absolute)
< error
| 42704 type "absolute" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS access)
< error
| 42704 type "access" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS action)
< error
| 42704 type "action" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS add)
< error
| 42704 type "add" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS admin)
< error
| 42704 type "admin" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS after)
< error
| 42704 type "after" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS aggregate)
< error
| 42704 type "aggregate" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS all)
< error
| 42601 syntax error at or near "all"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS all(1))
< error
| 42601 syntax error at or near "all"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS all(1, 1))
< error
| 42601 syntax error at or near "all"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS all precision)
< error
| 42601 syntax error at or near "all"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS all character)
< error
| 42601 syntax error at or near "all"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS all character(1))
< error
| 42601 syntax error at or near "all"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS all varying)
< error
| 42601 syntax error at or near "all"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS all varying(1))
< error
| 42601 syntax error at or near "all"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS all character varying)
< error
| 42601 syntax error at or near "all"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS all character varying(1))
< error
| 42601 syntax error at or near "all"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS also)
< error
| 42704 type "also" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS alter)
< error
| 42704 type "alter" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS always)
< error
| 42704 type "always" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS analyse)
< error
| 42601 syntax error at or near "analyse"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS analyse(1))
< error
| 42601 syntax error at or near "analyse"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS analyse(1, 1))
< error
| 42601 syntax error at or near "analyse"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS analyse precision)
< error
| 42601 syntax error at or near "analyse"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS analyse character)
< error
| 42601 syntax error at or near "analyse"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS analyse character(1))
< error
| 42601 syntax error at or near "analyse"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS analyse varying)
< error
| 42601 syntax error at or near "analyse"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS analyse varying(1))
< error
| 42601 syntax error at or near "analyse"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS analyse character varying)
< error
| 42601 syntax error at or near "analyse"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS analyse character varying(1))
< error
| 42601 syntax error at or near "analyse"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS analyze)
< error
| 42601 syntax error at or near "analyze"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS analyze(1))
< error
| 42601 syntax error at or near "analyze"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS analyze(1, 1))
< error
| 42601 syntax error at or near "analyze"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS analyze precision)
< error
| 42601 syntax error at or near "analyze"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS analyze character)
< error
| 42601 syntax error at or near "analyze"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS analyze character(1))
< error
| 42601 syntax error at or near "analyze"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS analyze varying)
< error
| 42601 syntax error at or near "analyze"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS analyze varying(1))
< error
| 42601 syntax error at or near "analyze"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS analyze character varying)
< error
| 42601 syntax error at or near "analyze"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS analyze character varying(1))
< error
| 42601 syntax error at or near "analyze"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS and)
< error
| 42601 syntax error at or near "and"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS and(1))
< error
| 42601 syntax error at or near "and"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS and(1, 1))
< error
| 42601 syntax error at or near "and"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS and precision)
< error
| 42601 syntax error at or near "and"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS and character)
< error
| 42601 syntax error at or near "and"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS and character(1))
< error
| 42601 syntax error at or near "and"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS and varying)
< error
| 42601 syntax error at or near "and"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS and varying(1))
< error
| 42601 syntax error at or near "and"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS and character varying)
< error
| 42601 syntax error at or near "and"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS and character varying(1))
< error
| 42601 syntax error at or near "and"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS any)
< error
| 42601 syntax error at or near "any"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS any(1))
< error
| 42601 syntax error at or near "any"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS any(1, 1))
< error
| 42601 syntax error at or near "any"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS any precision)
< error
| 42601 syntax error at or near "any"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS any character)
< error
| 42601 syntax error at or near "any"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS any character(1))
< error
| 42601 syntax error at or near "any"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS any varying)
< error
| 42601 syntax error at or near "any"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS any varying(1))
< error
| 42601 syntax error at or near "any"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS any character varying)
< error
| 42601 syntax error at or near "any"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS any character varying(1))
< error
| 42601 syntax error at or near "any"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS array)
< error
| 42601 syntax error at or near "array"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS array(1))
< error
| 42601 syntax error at or near "array"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS array(1, 1))
< error
| 42601 syntax error at or near "array"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS array precision)
< error
| 42601 syntax error at or near "array"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS array character)
< error
| 42601 syntax error at or near "array"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS array character(1))
< error
| 42601 syntax error at or near "array"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS array varying)
< error
| 42601 syntax error at or near "array"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS array varying(1))
< error
| 42601 syntax error at or near "array"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS array character varying)
< error
| 42601 syntax error at or near "array"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS array character varying(1))
< error
| 42601 syntax error at or near "array"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS as)
< error
| 42601 syntax error at or near "as"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS as(1))
< error
| 42601 syntax error at or near "as"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS as(1, 1))
< error
| 42601 syntax error at or near "as"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS as precision)
< error
| 42601 syntax error at or near "as"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS as character)
< error
| 42601 syntax error at or near "as"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS as character(1))
< error
| 42601 syntax error at or near "as"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS as varying)
< error
| 42601 syntax error at or near "as"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS as varying(1))
< error
| 42601 syntax error at or near "as"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS as character varying)
< error
| 42601 syntax error at or near "as"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS as character varying(1))
< error
| 42601 syntax error at or near "as"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS asc)
< error
| 42601 syntax error at or near "asc"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS asc(1))
< error
| 42601 syntax error at or near "asc"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS asc(1, 1))
< error
| 42601 syntax error at or near "asc"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS asc precision)
< error
| 42601 syntax error at or near "asc"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS asc character)
< error
| 42601 syntax error at or near "asc"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS asc character(1))
< error
| 42601 syntax error at or near "asc"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS asc varying)
< error
| 42601 syntax error at or near "asc"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS asc varying(1))
< error
| 42601 syntax error at or near "asc"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS asc character varying)
< error
| 42601 syntax error at or near "asc"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS asc character varying(1))
< error
| 42601 syntax error at or near "asc"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS asensitive)
< error
| 42704 type "asensitive" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS assertion)
< error
| 42704 type "assertion" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS assignment)
< error
| 42704 type "assignment" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS asymmetric)
< error
| 42601 syntax error at or near "asymmetric"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS asymmetric(1))
< error
| 42601 syntax error at or near "asymmetric"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS asymmetric(1, 1))
< error
| 42601 syntax error at or near "asymmetric"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS asymmetric precision)
< error
| 42601 syntax error at or near "asymmetric"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS asymmetric character)
< error
| 42601 syntax error at or near "asymmetric"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS asymmetric character(1))
< error
| 42601 syntax error at or near "asymmetric"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS asymmetric varying)
< error
| 42601 syntax error at or near "asymmetric"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS asymmetric varying(1))
< error
| 42601 syntax error at or near "asymmetric"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS asymmetric character varying)
< error
| 42601 syntax error at or near "asymmetric"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS asymmetric character varying(1))
< error
| 42601 syntax error at or near "asymmetric"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS at)
< error
| 42704 type "at" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS atomic)
< error
| 42704 type "atomic" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS attach)
< error
| 42704 type "attach" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS attribute)
< error
| 42704 type "attribute" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS authorization)
< error
| 42704 type "authorization" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS backward)
< error
| 42704 type "backward" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS before)
< error
| 42704 type "before" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS begin)
< error
| 42704 type "begin" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS between)
< error
| 42601 syntax error at or near "between"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS between(1))
< error
| 42601 syntax error at or near "between"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS between(1, 1))
< error
| 42601 syntax error at or near "between"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS between precision)
< error
| 42601 syntax error at or near "between"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS between character)
< error
| 42601 syntax error at or near "between"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS between character(1))
< error
| 42601 syntax error at or near "between"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS between varying)
< error
| 42601 syntax error at or near "between"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS between varying(1))
< error
| 42601 syntax error at or near "between"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS between character varying)
< error
| 42601 syntax error at or near "between"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS between character varying(1))
< error
| 42601 syntax error at or near "between"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS bigint)
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS binary)
< error
| 42704 type "binary" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS bit)
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS boolean)
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS both)
< error
| 42601 syntax error at or near "both"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS both(1))
< error
| 42601 syntax error at or near "both"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS both(1, 1))
< error
| 42601 syntax error at or near "both"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS both precision)
< error
| 42601 syntax error at or near "both"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS both character)
< error
| 42601 syntax error at or near "both"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS both character(1))
< error
| 42601 syntax error at or near "both"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS both varying)
< error
| 42601 syntax error at or near "both"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS both varying(1))
< error
| 42601 syntax error at or near "both"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS both character varying)
< error
| 42601 syntax error at or near "both"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS both character varying(1))
< error
| 42601 syntax error at or near "both"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS breadth)
< error
| 42704 type "breadth" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS by)
< error
| 42704 type "by" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS cache)
< error
| 42704 type "cache" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS call)
< error
| 42704 type "call" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS called)
< error
| 42704 type "called" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS cascade)
< error
| 42704 type "cascade" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS cascaded)
< error
| 42704 type "cascaded" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS case)
< error
| 42601 syntax error at or near "case"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS case(1))
< error
| 42601 syntax error at or near "case"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS case(1, 1))
< error
| 42601 syntax error at or near "case"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS case precision)
< error
| 42601 syntax error at or near "case"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS case character)
< error
| 42601 syntax error at or near "case"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS case character(1))
< error
| 42601 syntax error at or near "case"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS case varying)
< error
| 42601 syntax error at or near "case"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS case varying(1))
< error
| 42601 syntax error at or near "case"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS case character varying)
< error
| 42601 syntax error at or near "case"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS case character varying(1))
< error
| 42601 syntax error at or near "case"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS cast)
< error
| 42601 syntax error at or near "cast"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS cast(1))
< error
| 42601 syntax error at or near "cast"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS cast(1, 1))
< error
| 42601 syntax error at or near "cast"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS cast precision)
< error
| 42601 syntax error at or near "cast"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS cast character)
< error
| 42601 syntax error at or near "cast"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS cast character(1))
< error
| 42601 syntax error at or near "cast"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS cast varying)
< error
| 42601 syntax error at or near "cast"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS cast varying(1))
< error
| 42601 syntax error at or near "cast"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS cast character varying)
< error
| 42601 syntax error at or near "cast"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS cast character varying(1))
< error
| 42601 syntax error at or near "cast"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS catalog)
< error
| 42704 type "catalog" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS chain)
< error
| 42704 type "chain" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS char)
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS character)
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS characteristics)
< error
| 42704 type "characteristics" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS check)
< error
| 42601 syntax error at or near "check"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS check(1))
< error
| 42601 syntax error at or near "check"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS check(1, 1))
< error
| 42601 syntax error at or near "check"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS check precision)
< error
| 42601 syntax error at or near "check"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS check character)
< error
| 42601 syntax error at or near "check"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS check character(1))
< error
| 42601 syntax error at or near "check"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS check varying)
< error
| 42601 syntax error at or near "check"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS check varying(1))
< error
| 42601 syntax error at or near "check"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS check character varying)
< error
| 42601 syntax error at or near "check"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS check character varying(1))
< error
| 42601 syntax error at or near "check"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS checkpoint)
< error
| 42704 type "checkpoint" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS class)
< error
| 42704 type "class" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS close)
< error
| 42704 type "close" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS cluster)
< error
| 42704 type "cluster" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS coalesce)
< error
| 42601 syntax error at or near "coalesce"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS coalesce(1))
< error
| 42601 syntax error at or near "coalesce"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS coalesce(1, 1))
< error
| 42601 syntax error at or near "coalesce"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS coalesce precision)
< error
| 42601 syntax error at or near "coalesce"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS coalesce character)
< error
| 42601 syntax error at or near "coalesce"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS coalesce character(1))
< error
| 42601 syntax error at or near "coalesce"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS coalesce varying)
< error
| 42601 syntax error at or near "coalesce"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS coalesce varying(1))
< error
| 42601 syntax error at or near "coalesce"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS coalesce character varying)
< error
| 42601 syntax error at or near "coalesce"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS coalesce character varying(1))
< error
| 42601 syntax error at or near "coalesce"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS collate)
< error
| 42601 syntax error at or near "collate"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS collate(1))
< error
| 42601 syntax error at or near "collate"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS collate(1, 1))
< error
| 42601 syntax error at or near "collate"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS collate precision)
< error
| 42601 syntax error at or near "collate"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS collate character)
< error
| 42601 syntax error at or near "collate"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS collate character(1))
< error
| 42601 syntax error at or near "collate"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS collate varying)
< error
| 42601 syntax error at or near "collate"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS collate varying(1))
< error
| 42601 syntax error at or near "collate"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS collate character varying)
< error
| 42601 syntax error at or near "collate"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS collate character varying(1))
< error
| 42601 syntax error at or near "collate"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS collation)
< error
| 42704 type "collation" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS column)
< error
| 42601 syntax error at or near "column"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS column(1))
< error
| 42601 syntax error at or near "column"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS column(1, 1))
< error
| 42601 syntax error at or near "column"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS column precision)
< error
| 42601 syntax error at or near "column"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS column character)
< error
| 42601 syntax error at or near "column"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS column character(1))
< error
| 42601 syntax error at or near "column"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS column varying)
< error
| 42601 syntax error at or near "column"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS column varying(1))
< error
| 42601 syntax error at or near "column"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS column character varying)
< error
| 42601 syntax error at or near "column"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS column character varying(1))
< error
| 42601 syntax error at or near "column"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS columns)
< error
| 42704 type "columns" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS comment)
< error
| 42704 type "comment" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS comments)
< error
| 42704 type "comments" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS commit)
< error
| 42704 type "commit" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS committed)
< error
| 42704 type "committed" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS compression)
< error
| 42704 type "compression" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS concurrently)
< error
| 42704 type "concurrently" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS conditional)
< error
| 42704 type "conditional" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS configuration)
< error
| 42704 type "configuration" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS conflict)
< error
| 42704 type "conflict" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS connection)
< error
| 42704 type "connection" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS constraint)
< error
| 42601 syntax error at or near "constraint"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS constraint(1))
< error
| 42601 syntax error at or near "constraint"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS constraint(1, 1))
< error
| 42601 syntax error at or near "constraint"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS constraint precision)
< error
| 42601 syntax error at or near "constraint"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS constraint character)
< error
| 42601 syntax error at or near "constraint"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS constraint character(1))
< error
| 42601 syntax error at or near "constraint"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS constraint varying)
< error
| 42601 syntax error at or near "constraint"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS constraint varying(1))
< error
| 42601 syntax error at or near "constraint"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS constraint character varying)
< error
| 42601 syntax error at or near "constraint"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS constraint character varying(1))
< error
| 42601 syntax error at or near "constraint"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS constraints)
< error
| 42704 type "constraints" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS content)
< error
| 42704 type "content" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS continue)
< error
| 42704 type "continue" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS conversion)
< error
| 42704 type "conversion" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS copy)
< error
| 42704 type "copy" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS cost)
< error
| 42704 type "cost" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS create)
< error
| 42601 syntax error at or near "create"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS create(1))
< error
| 42601 syntax error at or near "create"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS create(1, 1))
< error
| 42601 syntax error at or near "create"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS create precision)
< error
| 42601 syntax error at or near "create"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS create character)
< error
| 42601 syntax error at or near "create"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS create character(1))
< error
| 42601 syntax error at or near "create"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS create varying)
< error
| 42601 syntax error at or near "create"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS create varying(1))
< error
| 42601 syntax error at or near "create"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS create character varying)
< error
| 42601 syntax error at or near "create"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS create character varying(1))
< error
| 42601 syntax error at or near "create"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS cross)
< error
| 42704 type "cross" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS csv)
< error
| 42704 type "csv" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS cube)
< error
| 42704 type "cube" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current)
< error
| 42704 type "current" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_catalog)
< error
| 42601 syntax error at or near "current_catalog"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_catalog(1))
< error
| 42601 syntax error at or near "current_catalog"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_catalog(1, 1))
< error
| 42601 syntax error at or near "current_catalog"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_catalog precision)
< error
| 42601 syntax error at or near "current_catalog"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_catalog character)
< error
| 42601 syntax error at or near "current_catalog"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_catalog character(1))
< error
| 42601 syntax error at or near "current_catalog"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_catalog varying)
< error
| 42601 syntax error at or near "current_catalog"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_catalog varying(1))
< error
| 42601 syntax error at or near "current_catalog"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_catalog character varying)
< error
| 42601 syntax error at or near "current_catalog"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_catalog character varying(1))
< error
| 42601 syntax error at or near "current_catalog"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_date)
< error
| 42601 syntax error at or near "current_date"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_date(1))
< error
| 42601 syntax error at or near "current_date"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_date(1, 1))
< error
| 42601 syntax error at or near "current_date"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_date precision)
< error
| 42601 syntax error at or near "current_date"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_date character)
< error
| 42601 syntax error at or near "current_date"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_date character(1))
< error
| 42601 syntax error at or near "current_date"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_date varying)
< error
| 42601 syntax error at or near "current_date"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_date varying(1))
< error
| 42601 syntax error at or near "current_date"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_date character varying)
< error
| 42601 syntax error at or near "current_date"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_date character varying(1))
< error
| 42601 syntax error at or near "current_date"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_role)
< error
| 42601 syntax error at or near "current_role"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_role(1))
< error
| 42601 syntax error at or near "current_role"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_role(1, 1))
< error
| 42601 syntax error at or near "current_role"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_role precision)
< error
| 42601 syntax error at or near "current_role"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_role character)
< error
| 42601 syntax error at or near "current_role"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_role character(1))
< error
| 42601 syntax error at or near "current_role"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_role varying)
< error
| 42601 syntax error at or near "current_role"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_role varying(1))
< error
| 42601 syntax error at or near "current_role"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_role character varying)
< error
| 42601 syntax error at or near "current_role"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_role character varying(1))
< error
| 42601 syntax error at or near "current_role"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_schema)
< error
| 42704 type "current_schema" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_time)
< error
| 42601 syntax error at or near "current_time"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_time(1))
< error
| 42601 syntax error at or near "current_time"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_time(1, 1))
< error
| 42601 syntax error at or near "current_time"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_time precision)
< error
| 42601 syntax error at or near "current_time"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_time character)
< error
| 42601 syntax error at or near "current_time"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_time character(1))
< error
| 42601 syntax error at or near "current_time"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_time varying)
< error
| 42601 syntax error at or near "current_time"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_time varying(1))
< error
| 42601 syntax error at or near "current_time"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_time character varying)
< error
| 42601 syntax error at or near "current_time"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_time character varying(1))
< error
| 42601 syntax error at or near "current_time"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_timestamp)
< error
| 42601 syntax error at or near "current_timestamp"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_timestamp(1))
< error
| 42601 syntax error at or near "current_timestamp"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_timestamp(1, 1))
< error
| 42601 syntax error at or near "current_timestamp"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_timestamp precision)
< error
| 42601 syntax error at or near "current_timestamp"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_timestamp character)
< error
| 42601 syntax error at or near "current_timestamp"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_timestamp character(1))
< error
| 42601 syntax error at or near "current_timestamp"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_timestamp varying)
< error
| 42601 syntax error at or near "current_timestamp"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_timestamp varying(1))
< error
| 42601 syntax error at or near "current_timestamp"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_timestamp character varying)
< error
| 42601 syntax error at or near "current_timestamp"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_timestamp character varying(1))
< error
| 42601 syntax error at or near "current_timestamp"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_user)
< error
| 42601 syntax error at or near "current_user"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_user(1))
< error
| 42601 syntax error at or near "current_user"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_user(1, 1))
< error
| 42601 syntax error at or near "current_user"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_user precision)
< error
| 42601 syntax error at or near "current_user"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_user character)
< error
| 42601 syntax error at or near "current_user"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_user character(1))
< error
| 42601 syntax error at or near "current_user"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_user varying)
< error
| 42601 syntax error at or near "current_user"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_user varying(1))
< error
| 42601 syntax error at or near "current_user"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_user character varying)
< error
| 42601 syntax error at or near "current_user"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS current_user character varying(1))
< error
| 42601 syntax error at or near "current_user"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS cursor)
< error
| 42704 type "cursor" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS cycle)
< error
| 42704 type "cycle" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS data)
< error
| 42704 type "data" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS database)
< error
| 42704 type "database" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS day)
< error
| 42704 type "day" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS deallocate)
< error
| 42704 type "deallocate" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS dec)
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS decimal)
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS declare)
< error
| 42704 type "declare" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS default)
< error
| 42601 syntax error at or near "default"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS default(1))
< error
| 42601 syntax error at or near "default"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS default(1, 1))
< error
| 42601 syntax error at or near "default"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS default precision)
< error
| 42601 syntax error at or near "default"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS default character)
< error
| 42601 syntax error at or near "default"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS default character(1))
< error
| 42601 syntax error at or near "default"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS default varying)
< error
| 42601 syntax error at or near "default"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS default varying(1))
< error
| 42601 syntax error at or near "default"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS default character varying)
< error
| 42601 syntax error at or near "default"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS default character varying(1))
< error
| 42601 syntax error at or near "default"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS defaults)
< error
| 42704 type "defaults" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS deferrable)
< error
| 42601 syntax error at or near "deferrable"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS deferrable(1))
< error
| 42601 syntax error at or near "deferrable"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS deferrable(1, 1))
< error
| 42601 syntax error at or near "deferrable"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS deferrable precision)
< error
| 42601 syntax error at or near "deferrable"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS deferrable character)
< error
| 42601 syntax error at or near "deferrable"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS deferrable character(1))
< error
| 42601 syntax error at or near "deferrable"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS deferrable varying)
< error
| 42601 syntax error at or near "deferrable"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS deferrable varying(1))
< error
| 42601 syntax error at or near "deferrable"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS deferrable character varying)
< error
| 42601 syntax error at or near "deferrable"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS deferrable character varying(1))
< error
| 42601 syntax error at or near "deferrable"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS deferred)
< error
| 42704 type "deferred" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS definer)
< error
| 42704 type "definer" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS delete)
< error
| 42704 type "delete" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS delimiter)
< error
| 42704 type "delimiter" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS delimiters)
< error
| 42704 type "delimiters" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS depends)
< error
| 42704 type "depends" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS depth)
< error
| 42704 type "depth" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS desc)
< error
| 42601 syntax error at or near "desc"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS desc(1))
< error
| 42601 syntax error at or near "desc"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS desc(1, 1))
< error
| 42601 syntax error at or near "desc"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS desc precision)
< error
| 42601 syntax error at or near "desc"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS desc character)
< error
| 42601 syntax error at or near "desc"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS desc character(1))
< error
| 42601 syntax error at or near "desc"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS desc varying)
< error
| 42601 syntax error at or near "desc"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS desc varying(1))
< error
| 42601 syntax error at or near "desc"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS desc character varying)
< error
| 42601 syntax error at or near "desc"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS desc character varying(1))
< error
| 42601 syntax error at or near "desc"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS detach)
< error
| 42704 type "detach" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS dictionary)
< error
| 42704 type "dictionary" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS disable)
< error
| 42704 type "disable" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS discard)
< error
| 42704 type "discard" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS distinct)
< error
| 42601 syntax error at or near "distinct"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS distinct(1))
< error
| 42601 syntax error at or near "distinct"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS distinct(1, 1))
< error
| 42601 syntax error at or near "distinct"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS distinct precision)
< error
| 42601 syntax error at or near "distinct"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS distinct character)
< error
| 42601 syntax error at or near "distinct"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS distinct character(1))
< error
| 42601 syntax error at or near "distinct"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS distinct varying)
< error
| 42601 syntax error at or near "distinct"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS distinct varying(1))
< error
| 42601 syntax error at or near "distinct"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS distinct character varying)
< error
| 42601 syntax error at or near "distinct"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS distinct character varying(1))
< error
| 42601 syntax error at or near "distinct"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS do)
< error
| 42601 syntax error at or near "do"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS do(1))
< error
| 42601 syntax error at or near "do"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS do(1, 1))
< error
| 42601 syntax error at or near "do"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS do precision)
< error
| 42601 syntax error at or near "do"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS do character)
< error
| 42601 syntax error at or near "do"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS do character(1))
< error
| 42601 syntax error at or near "do"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS do varying)
< error
| 42601 syntax error at or near "do"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS do varying(1))
< error
| 42601 syntax error at or near "do"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS do character varying)
< error
| 42601 syntax error at or near "do"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS do character varying(1))
< error
| 42601 syntax error at or near "do"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS document)
< error
| 42704 type "document" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS domain)
< error
| 42704 type "domain" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS double)
< error
| 42704 type "double" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS drop)
< error
| 42704 type "drop" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS each)
< error
| 42704 type "each" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS else)
< error
| 42601 syntax error at or near "else"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS else(1))
< error
| 42601 syntax error at or near "else"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS else(1, 1))
< error
| 42601 syntax error at or near "else"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS else precision)
< error
| 42601 syntax error at or near "else"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS else character)
< error
| 42601 syntax error at or near "else"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS else character(1))
< error
| 42601 syntax error at or near "else"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS else varying)
< error
| 42601 syntax error at or near "else"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS else varying(1))
< error
| 42601 syntax error at or near "else"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS else character varying)
< error
| 42601 syntax error at or near "else"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS else character varying(1))
< error
| 42601 syntax error at or near "else"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS empty)
< error
| 42704 type "empty" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS enable)
< error
| 42704 type "enable" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS encoding)
< error
| 42704 type "encoding" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS encrypted)
< error
| 42704 type "encrypted" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS end)
< error
| 42601 syntax error at or near "end"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS end(1))
< error
| 42601 syntax error at or near "end"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS end(1, 1))
< error
| 42601 syntax error at or near "end"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS end precision)
< error
| 42601 syntax error at or near "end"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS end character)
< error
| 42601 syntax error at or near "end"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS end character(1))
< error
| 42601 syntax error at or near "end"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS end varying)
< error
| 42601 syntax error at or near "end"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS end varying(1))
< error
| 42601 syntax error at or near "end"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS end character varying)
< error
| 42601 syntax error at or near "end"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS end character varying(1))
< error
| 42601 syntax error at or near "end"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS enum)
< error
| 42704 type "enum" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS error)
< error
| 42704 type "error" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS escape)
< error
| 42704 type "escape" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS event)
< error
| 42704 type "event" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS except)
< error
| 42601 syntax error at or near "except"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS except(1))
< error
| 42601 syntax error at or near "except"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS except(1, 1))
< error
| 42601 syntax error at or near "except"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS except precision)
< error
| 42601 syntax error at or near "except"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS except character)
< error
| 42601 syntax error at or near "except"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS except character(1))
< error
| 42601 syntax error at or near "except"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS except varying)
< error
| 42601 syntax error at or near "except"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS except varying(1))
< error
| 42601 syntax error at or near "except"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS except character varying)
< error
| 42601 syntax error at or near "except"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS except character varying(1))
< error
| 42601 syntax error at or near "except"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS exclude)
< error
| 42704 type "exclude" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS excluding)
< error
| 42704 type "excluding" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS exclusive)
< error
| 42704 type "exclusive" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS execute)
< error
| 42704 type "execute" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS exists)
< error
| 42601 syntax error at or near "exists"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS exists(1))
< error
| 42601 syntax error at or near "exists"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS exists(1, 1))
< error
| 42601 syntax error at or near "exists"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS exists precision)
< error
| 42601 syntax error at or near "exists"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS exists character)
< error
| 42601 syntax error at or near "exists"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS exists character(1))
< error
| 42601 syntax error at or near "exists"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS exists varying)
< error
| 42601 syntax error at or near "exists"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS exists varying(1))
< error
| 42601 syntax error at or near "exists"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS exists character varying)
< error
| 42601 syntax error at or near "exists"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS exists character varying(1))
< error
| 42601 syntax error at or near "exists"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS explain)
< error
| 42704 type "explain" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS expression)
< error
| 42704 type "expression" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS extension)
< error
| 42704 type "extension" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS external)
< error
| 42704 type "external" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS extract)
< error
| 42601 syntax error at or near "extract"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS extract(1))
< error
| 42601 syntax error at or near "extract"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS extract(1, 1))
< error
| 42601 syntax error at or near "extract"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS extract precision)
< error
| 42601 syntax error at or near "extract"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS extract character)
< error
| 42601 syntax error at or near "extract"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS extract character(1))
< error
| 42601 syntax error at or near "extract"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS extract varying)
< error
| 42601 syntax error at or near "extract"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS extract varying(1))
< error
| 42601 syntax error at or near "extract"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS extract character varying)
< error
| 42601 syntax error at or near "extract"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS extract character varying(1))
< error
| 42601 syntax error at or near "extract"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS false)
< error
| 42601 syntax error at or near "false"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS false(1))
< error
| 42601 syntax error at or near "false"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS false(1, 1))
< error
| 42601 syntax error at or near "false"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS false precision)
< error
| 42601 syntax error at or near "false"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS false character)
< error
| 42601 syntax error at or near "false"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS false character(1))
< error
| 42601 syntax error at or near "false"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS false varying)
< error
| 42601 syntax error at or near "false"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS false varying(1))
< error
| 42601 syntax error at or near "false"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS false character varying)
< error
| 42601 syntax error at or near "false"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS false character varying(1))
< error
| 42601 syntax error at or near "false"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS family)
< error
| 42704 type "family" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS fetch)
< error
| 42601 syntax error at or near "fetch"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS fetch(1))
< error
| 42601 syntax error at or near "fetch"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS fetch(1, 1))
< error
| 42601 syntax error at or near "fetch"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS fetch precision)
< error
| 42601 syntax error at or near "fetch"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS fetch character)
< error
| 42601 syntax error at or near "fetch"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS fetch character(1))
< error
| 42601 syntax error at or near "fetch"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS fetch varying)
< error
| 42601 syntax error at or near "fetch"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS fetch varying(1))
< error
| 42601 syntax error at or near "fetch"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS fetch character varying)
< error
| 42601 syntax error at or near "fetch"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS fetch character varying(1))
< error
| 42601 syntax error at or near "fetch"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS filter)
< error
| 42704 type "filter" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS finalize)
< error
| 42704 type "finalize" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS first)
< error
| 42704 type "first" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS float)
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS following)
< error
| 42704 type "following" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS for)
< error
| 42601 syntax error at or near "for"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS for(1))
< error
| 42601 syntax error at or near "for"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS for(1, 1))
< error
| 42601 syntax error at or near "for"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS for precision)
< error
| 42601 syntax error at or near "for"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS for character)
< error
| 42601 syntax error at or near "for"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS for character(1))
< error
| 42601 syntax error at or near "for"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS for varying)
< error
| 42601 syntax error at or near "for"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS for varying(1))
< error
| 42601 syntax error at or near "for"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS for character varying)
< error
| 42601 syntax error at or near "for"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS for character varying(1))
< error
| 42601 syntax error at or near "for"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS force)
< error
| 42704 type "force" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS foreign)
< error
| 42601 syntax error at or near "foreign"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS foreign(1))
< error
| 42601 syntax error at or near "foreign"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS foreign(1, 1))
< error
| 42601 syntax error at or near "foreign"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS foreign precision)
< error
| 42601 syntax error at or near "foreign"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS foreign character)
< error
| 42601 syntax error at or near "foreign"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS foreign character(1))
< error
| 42601 syntax error at or near "foreign"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS foreign varying)
< error
| 42601 syntax error at or near "foreign"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS foreign varying(1))
< error
| 42601 syntax error at or near "foreign"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS foreign character varying)
< error
| 42601 syntax error at or near "foreign"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS foreign character varying(1))
< error
| 42601 syntax error at or near "foreign"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS format)
< error
| 42704 type "format" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS forward)
< error
| 42704 type "forward" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS freeze)
< error
| 42704 type "freeze" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS from)
< error
| 42601 syntax error at or near "from"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS from(1))
< error
| 42601 syntax error at or near "from"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS from(1, 1))
< error
| 42601 syntax error at or near "from"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS from precision)
< error
| 42601 syntax error at or near "from"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS from character)
< error
| 42601 syntax error at or near "from"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS from character(1))
< error
| 42601 syntax error at or near "from"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS from varying)
< error
| 42601 syntax error at or near "from"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS from varying(1))
< error
| 42601 syntax error at or near "from"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS from character varying)
< error
| 42601 syntax error at or near "from"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS from character varying(1))
< error
| 42601 syntax error at or near "from"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS full)
< error
| 42704 type "full" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS function)
< error
| 42704 type "function" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS functions)
< error
| 42704 type "functions" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS generated)
< error
| 42704 type "generated" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS global)
< error
| 42704 type "global" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS grant)
< error
| 42601 syntax error at or near "grant"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS grant(1))
< error
| 42601 syntax error at or near "grant"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS grant(1, 1))
< error
| 42601 syntax error at or near "grant"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS grant precision)
< error
| 42601 syntax error at or near "grant"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS grant character)
< error
| 42601 syntax error at or near "grant"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS grant character(1))
< error
| 42601 syntax error at or near "grant"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS grant varying)
< error
| 42601 syntax error at or near "grant"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS grant varying(1))
< error
| 42601 syntax error at or near "grant"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS grant character varying)
< error
| 42601 syntax error at or near "grant"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS grant character varying(1))
< error
| 42601 syntax error at or near "grant"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS granted)
< error
| 42704 type "granted" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS greatest)
< error
| 42601 syntax error at or near "greatest"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS greatest(1))
< error
| 42601 syntax error at or near "greatest"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS greatest(1, 1))
< error
| 42601 syntax error at or near "greatest"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS greatest precision)
< error
| 42601 syntax error at or near "greatest"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS greatest character)
< error
| 42601 syntax error at or near "greatest"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS greatest character(1))
< error
| 42601 syntax error at or near "greatest"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS greatest varying)
< error
| 42601 syntax error at or near "greatest"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS greatest varying(1))
< error
| 42601 syntax error at or near "greatest"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS greatest character varying)
< error
| 42601 syntax error at or near "greatest"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS greatest character varying(1))
< error
| 42601 syntax error at or near "greatest"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS group)
< error
| 42601 syntax error at or near "group"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS group(1))
< error
| 42601 syntax error at or near "group"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS group(1, 1))
< error
| 42601 syntax error at or near "group"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS group precision)
< error
| 42601 syntax error at or near "group"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS group character)
< error
| 42601 syntax error at or near "group"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS group character(1))
< error
| 42601 syntax error at or near "group"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS group varying)
< error
| 42601 syntax error at or near "group"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS group varying(1))
< error
| 42601 syntax error at or near "group"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS group character varying)
< error
| 42601 syntax error at or near "group"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS group character varying(1))
< error
| 42601 syntax error at or near "group"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS grouping)
< error
| 42601 syntax error at or near "grouping"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS grouping(1))
< error
| 42601 syntax error at or near "grouping"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS grouping(1, 1))
< error
| 42601 syntax error at or near "grouping"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS grouping precision)
< error
| 42601 syntax error at or near "grouping"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS grouping character)
< error
| 42601 syntax error at or near "grouping"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS grouping character(1))
< error
| 42601 syntax error at or near "grouping"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS grouping varying)
< error
| 42601 syntax error at or near "grouping"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS grouping varying(1))
< error
| 42601 syntax error at or near "grouping"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS grouping character varying)
< error
| 42601 syntax error at or near "grouping"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS grouping character varying(1))
< error
| 42601 syntax error at or near "grouping"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS groups)
< error
| 42704 type "groups" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS handler)
< error
| 42704 type "handler" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS having)
< error
| 42601 syntax error at or near "having"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS having(1))
< error
| 42601 syntax error at or near "having"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS having(1, 1))
< error
| 42601 syntax error at or near "having"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS having precision)
< error
| 42601 syntax error at or near "having"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS having character)
< error
| 42601 syntax error at or near "having"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS having character(1))
< error
| 42601 syntax error at or near "having"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS having varying)
< error
| 42601 syntax error at or near "having"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS having varying(1))
< error
| 42601 syntax error at or near "having"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS having character varying)
< error
| 42601 syntax error at or near "having"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS having character varying(1))
< error
| 42601 syntax error at or near "having"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS header)
< error
| 42704 type "header" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS hold)
< error
| 42704 type "hold" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS hour)
< error
| 42704 type "hour" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS identity)
< error
| 42704 type "identity" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS if)
< error
| 42704 type "if" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS ilike)
< error
| 42704 type "ilike" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS immediate)
< error
| 42704 type "immediate" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS immutable)
< error
| 42704 type "immutable" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS implicit)
< error
| 42704 type "implicit" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS import)
< error
| 42704 type "import" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS in)
< error
| 42601 syntax error at or near "in"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS in(1))
< error
| 42601 syntax error at or near "in"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS in(1, 1))
< error
| 42601 syntax error at or near "in"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS in precision)
< error
| 42601 syntax error at or near "in"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS in character)
< error
| 42601 syntax error at or near "in"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS in character(1))
< error
| 42601 syntax error at or near "in"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS in varying)
< error
| 42601 syntax error at or near "in"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS in varying(1))
< error
| 42601 syntax error at or near "in"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS in character varying)
< error
| 42601 syntax error at or near "in"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS in character varying(1))
< error
| 42601 syntax error at or near "in"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS include)
< error
| 42704 type "include" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS including)
< error
| 42704 type "including" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS increment)
< error
| 42704 type "increment" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS indent)
< error
| 42704 type "indent" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS index)
< error
| 42704 type "index" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS indexes)
< error
| 42704 type "indexes" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS inherit)
< error
| 42704 type "inherit" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS inherits)
< error
| 42704 type "inherits" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS initially)
< error
| 42601 syntax error at or near "initially"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS initially(1))
< error
| 42601 syntax error at or near "initially"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS initially(1, 1))
< error
| 42601 syntax error at or near "initially"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS initially precision)
< error
| 42601 syntax error at or near "initially"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS initially character)
< error
| 42601 syntax error at or near "initially"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS initially character(1))
< error
| 42601 syntax error at or near "initially"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS initially varying)
< error
| 42601 syntax error at or near "initially"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS initially varying(1))
< error
| 42601 syntax error at or near "initially"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS initially character varying)
< error
| 42601 syntax error at or near "initially"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS initially character varying(1))
< error
| 42601 syntax error at or near "initially"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS inline)
< error
| 42704 type "inline" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS inner)
< error
| 42704 type "inner" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS inout)
< error
| 42601 syntax error at or near "inout"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS inout(1))
< error
| 42601 syntax error at or near "inout"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS inout(1, 1))
< error
| 42601 syntax error at or near "inout"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS inout precision)
< error
| 42601 syntax error at or near "inout"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS inout character)
< error
| 42601 syntax error at or near "inout"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS inout character(1))
< error
| 42601 syntax error at or near "inout"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS inout varying)
< error
| 42601 syntax error at or near "inout"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS inout varying(1))
< error
| 42601 syntax error at or near "inout"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS inout character varying)
< error
| 42601 syntax error at or near "inout"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS inout character varying(1))
< error
| 42601 syntax error at or near "inout"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS input)
< error
| 42704 type "input" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS insensitive)
< error
| 42704 type "insensitive" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS insert)
< error
| 42704 type "insert" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS instead)
< error
| 42704 type "instead" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS int)
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS integer)
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS intersect)
< error
| 42601 syntax error at or near "intersect"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS intersect(1))
< error
| 42601 syntax error at or near "intersect"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS intersect(1, 1))
< error
| 42601 syntax error at or near "intersect"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS intersect precision)
< error
| 42601 syntax error at or near "intersect"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS intersect character)
< error
| 42601 syntax error at or near "intersect"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS intersect character(1))
< error
| 42601 syntax error at or near "intersect"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS intersect varying)
< error
| 42601 syntax error at or near "intersect"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS intersect varying(1))
< error
| 42601 syntax error at or near "intersect"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS intersect character varying)
< error
| 42601 syntax error at or near "intersect"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS intersect character varying(1))
< error
| 42601 syntax error at or near "intersect"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS interval)
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS into)
< error
| 42601 syntax error at or near "into"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS into(1))
< error
| 42601 syntax error at or near "into"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS into(1, 1))
< error
| 42601 syntax error at or near "into"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS into precision)
< error
| 42601 syntax error at or near "into"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS into character)
< error
| 42601 syntax error at or near "into"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS into character(1))
< error
| 42601 syntax error at or near "into"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS into varying)
< error
| 42601 syntax error at or near "into"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS into varying(1))
< error
| 42601 syntax error at or near "into"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS into character varying)
< error
| 42601 syntax error at or near "into"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS into character varying(1))
< error
| 42601 syntax error at or near "into"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS invoker)
< error
| 42704 type "invoker" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS is)
< error
| 42704 type "is" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS isnull)
< error
| 42704 type "isnull" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS isolation)
< error
| 42704 type "isolation" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS join)
< error
| 42704 type "join" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json)
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_array)
< error
| 42601 syntax error at or near "json_array"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_array(1))
< error
| 42601 syntax error at or near "json_array"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_array(1, 1))
< error
| 42601 syntax error at or near "json_array"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_array precision)
< error
| 42601 syntax error at or near "json_array"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_array character)
< error
| 42601 syntax error at or near "json_array"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_array character(1))
< error
| 42601 syntax error at or near "json_array"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_array varying)
< error
| 42601 syntax error at or near "json_array"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_array varying(1))
< error
| 42601 syntax error at or near "json_array"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_array character varying)
< error
| 42601 syntax error at or near "json_array"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_array character varying(1))
< error
| 42601 syntax error at or near "json_array"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_arrayagg)
< error
| 42601 syntax error at or near "json_arrayagg"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_arrayagg(1))
< error
| 42601 syntax error at or near "json_arrayagg"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_arrayagg(1, 1))
< error
| 42601 syntax error at or near "json_arrayagg"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_arrayagg precision)
< error
| 42601 syntax error at or near "json_arrayagg"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_arrayagg character)
< error
| 42601 syntax error at or near "json_arrayagg"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_arrayagg character(1))
< error
| 42601 syntax error at or near "json_arrayagg"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_arrayagg varying)
< error
| 42601 syntax error at or near "json_arrayagg"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_arrayagg varying(1))
< error
| 42601 syntax error at or near "json_arrayagg"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_arrayagg character varying)
< error
| 42601 syntax error at or near "json_arrayagg"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_arrayagg character varying(1))
< error
| 42601 syntax error at or near "json_arrayagg"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_exists)
< error
| 42601 syntax error at or near "json_exists"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_exists(1))
< error
| 42601 syntax error at or near "json_exists"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_exists(1, 1))
< error
| 42601 syntax error at or near "json_exists"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_exists precision)
< error
| 42601 syntax error at or near "json_exists"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_exists character)
< error
| 42601 syntax error at or near "json_exists"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_exists character(1))
< error
| 42601 syntax error at or near "json_exists"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_exists varying)
< error
| 42601 syntax error at or near "json_exists"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_exists varying(1))
< error
| 42601 syntax error at or near "json_exists"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_exists character varying)
< error
| 42601 syntax error at or near "json_exists"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_exists character varying(1))
< error
| 42601 syntax error at or near "json_exists"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_object)
< error
| 42601 syntax error at or near "json_object"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_object(1))
< error
| 42601 syntax error at or near "json_object"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_object(1, 1))
< error
| 42601 syntax error at or near "json_object"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_object precision)
< error
| 42601 syntax error at or near "json_object"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_object character)
< error
| 42601 syntax error at or near "json_object"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_object character(1))
< error
| 42601 syntax error at or near "json_object"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_object varying)
< error
| 42601 syntax error at or near "json_object"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_object varying(1))
< error
| 42601 syntax error at or near "json_object"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_object character varying)
< error
| 42601 syntax error at or near "json_object"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_object character varying(1))
< error
| 42601 syntax error at or near "json_object"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_objectagg)
< error
| 42601 syntax error at or near "json_objectagg"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_objectagg(1))
< error
| 42601 syntax error at or near "json_objectagg"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_objectagg(1, 1))
< error
| 42601 syntax error at or near "json_objectagg"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_objectagg precision)
< error
| 42601 syntax error at or near "json_objectagg"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_objectagg character)
< error
| 42601 syntax error at or near "json_objectagg"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_objectagg character(1))
< error
| 42601 syntax error at or near "json_objectagg"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_objectagg varying)
< error
| 42601 syntax error at or near "json_objectagg"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_objectagg varying(1))
< error
| 42601 syntax error at or near "json_objectagg"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_objectagg character varying)
< error
| 42601 syntax error at or near "json_objectagg"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_objectagg character varying(1))
< error
| 42601 syntax error at or near "json_objectagg"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_query)
< error
| 42601 syntax error at or near "json_query"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_query(1))
< error
| 42601 syntax error at or near "json_query"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_query(1, 1))
< error
| 42601 syntax error at or near "json_query"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_query precision)
< error
| 42601 syntax error at or near "json_query"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_query character)
< error
| 42601 syntax error at or near "json_query"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_query character(1))
< error
| 42601 syntax error at or near "json_query"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_query varying)
< error
| 42601 syntax error at or near "json_query"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_query varying(1))
< error
| 42601 syntax error at or near "json_query"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_query character varying)
< error
| 42601 syntax error at or near "json_query"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_query character varying(1))
< error
| 42601 syntax error at or near "json_query"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_scalar)
< error
| 42601 syntax error at or near "json_scalar"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_scalar(1))
< error
| 42601 syntax error at or near "json_scalar"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_scalar(1, 1))
< error
| 42601 syntax error at or near "json_scalar"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_scalar precision)
< error
| 42601 syntax error at or near "json_scalar"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_scalar character)
< error
| 42601 syntax error at or near "json_scalar"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_scalar character(1))
< error
| 42601 syntax error at or near "json_scalar"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_scalar varying)
< error
| 42601 syntax error at or near "json_scalar"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_scalar varying(1))
< error
| 42601 syntax error at or near "json_scalar"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_scalar character varying)
< error
| 42601 syntax error at or near "json_scalar"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_scalar character varying(1))
< error
| 42601 syntax error at or near "json_scalar"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_serialize)
< error
| 42601 syntax error at or near "json_serialize"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_serialize(1))
< error
| 42601 syntax error at or near "json_serialize"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_serialize(1, 1))
< error
| 42601 syntax error at or near "json_serialize"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_serialize precision)
< error
| 42601 syntax error at or near "json_serialize"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_serialize character)
< error
| 42601 syntax error at or near "json_serialize"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_serialize character(1))
< error
| 42601 syntax error at or near "json_serialize"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_serialize varying)
< error
| 42601 syntax error at or near "json_serialize"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_serialize varying(1))
< error
| 42601 syntax error at or near "json_serialize"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_serialize character varying)
< error
| 42601 syntax error at or near "json_serialize"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_serialize character varying(1))
< error
| 42601 syntax error at or near "json_serialize"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_table)
< error
| 42601 syntax error at or near "json_table"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_table(1))
< error
| 42601 syntax error at or near "json_table"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_table(1, 1))
< error
| 42601 syntax error at or near "json_table"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_table precision)
< error
| 42601 syntax error at or near "json_table"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_table character)
< error
| 42601 syntax error at or near "json_table"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_table character(1))
< error
| 42601 syntax error at or near "json_table"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_table varying)
< error
| 42601 syntax error at or near "json_table"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_table varying(1))
< error
| 42601 syntax error at or near "json_table"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_table character varying)
< error
| 42601 syntax error at or near "json_table"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_table character varying(1))
< error
| 42601 syntax error at or near "json_table"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_value)
< error
| 42601 syntax error at or near "json_value"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_value(1))
< error
| 42601 syntax error at or near "json_value"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_value(1, 1))
< error
| 42601 syntax error at or near "json_value"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_value precision)
< error
| 42601 syntax error at or near "json_value"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_value character)
< error
| 42601 syntax error at or near "json_value"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_value character(1))
< error
| 42601 syntax error at or near "json_value"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_value varying)
< error
| 42601 syntax error at or near "json_value"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_value varying(1))
< error
| 42601 syntax error at or near "json_value"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_value character varying)
< error
| 42601 syntax error at or near "json_value"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS json_value character varying(1))
< error
| 42601 syntax error at or near "json_value"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS keep)
< error
| 42704 type "keep" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS key)
< error
| 42704 type "key" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS keys)
< error
| 42704 type "keys" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS label)
< error
| 42704 type "label" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS language)
< error
| 42704 type "language" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS large)
< error
| 42704 type "large" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS last)
< error
| 42704 type "last" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS lateral)
< error
| 42601 syntax error at or near "lateral"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS lateral(1))
< error
| 42601 syntax error at or near "lateral"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS lateral(1, 1))
< error
| 42601 syntax error at or near "lateral"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS lateral precision)
< error
| 42601 syntax error at or near "lateral"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS lateral character)
< error
| 42601 syntax error at or near "lateral"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS lateral character(1))
< error
| 42601 syntax error at or near "lateral"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS lateral varying)
< error
| 42601 syntax error at or near "lateral"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS lateral varying(1))
< error
| 42601 syntax error at or near "lateral"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS lateral character varying)
< error
| 42601 syntax error at or near "lateral"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS lateral character varying(1))
< error
| 42601 syntax error at or near "lateral"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS leading)
< error
| 42601 syntax error at or near "leading"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS leading(1))
< error
| 42601 syntax error at or near "leading"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS leading(1, 1))
< error
| 42601 syntax error at or near "leading"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS leading precision)
< error
| 42601 syntax error at or near "leading"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS leading character)
< error
| 42601 syntax error at or near "leading"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS leading character(1))
< error
| 42601 syntax error at or near "leading"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS leading varying)
< error
| 42601 syntax error at or near "leading"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS leading varying(1))
< error
| 42601 syntax error at or near "leading"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS leading character varying)
< error
| 42601 syntax error at or near "leading"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS leading character varying(1))
< error
| 42601 syntax error at or near "leading"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS leakproof)
< error
| 42704 type "leakproof" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS least)
< error
| 42601 syntax error at or near "least"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS least(1))
< error
| 42601 syntax error at or near "least"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS least(1, 1))
< error
| 42601 syntax error at or near "least"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS least precision)
< error
| 42601 syntax error at or near "least"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS least character)
< error
| 42601 syntax error at or near "least"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS least character(1))
< error
| 42601 syntax error at or near "least"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS least varying)
< error
| 42601 syntax error at or near "least"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS least varying(1))
< error
| 42601 syntax error at or near "least"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS least character varying)
< error
| 42601 syntax error at or near "least"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS least character varying(1))
< error
| 42601 syntax error at or near "least"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS left)
< error
| 42704 type "left" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS level)
< error
| 42704 type "level" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS like)
< error
| 42704 type "like" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS limit)
< error
| 42601 syntax error at or near "limit"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS limit(1))
< error
| 42601 syntax error at or near "limit"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS limit(1, 1))
< error
| 42601 syntax error at or near "limit"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS limit precision)
< error
| 42601 syntax error at or near "limit"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS limit character)
< error
| 42601 syntax error at or near "limit"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS limit character(1))
< error
| 42601 syntax error at or near "limit"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS limit varying)
< error
| 42601 syntax error at or near "limit"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS limit varying(1))
< error
| 42601 syntax error at or near "limit"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS limit character varying)
< error
| 42601 syntax error at or near "limit"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS limit character varying(1))
< error
| 42601 syntax error at or near "limit"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS listen)
< error
| 42704 type "listen" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS load)
< error
| 42704 type "load" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS local)
< error
| 42704 type "local" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS localtime)
< error
| 42601 syntax error at or near "localtime"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS localtime(1))
< error
| 42601 syntax error at or near "localtime"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS localtime(1, 1))
< error
| 42601 syntax error at or near "localtime"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS localtime precision)
< error
| 42601 syntax error at or near "localtime"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS localtime character)
< error
| 42601 syntax error at or near "localtime"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS localtime character(1))
< error
| 42601 syntax error at or near "localtime"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS localtime varying)
< error
| 42601 syntax error at or near "localtime"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS localtime varying(1))
< error
| 42601 syntax error at or near "localtime"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS localtime character varying)
< error
| 42601 syntax error at or near "localtime"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS localtime character varying(1))
< error
| 42601 syntax error at or near "localtime"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS localtimestamp)
< error
| 42601 syntax error at or near "localtimestamp"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS localtimestamp(1))
< error
| 42601 syntax error at or near "localtimestamp"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS localtimestamp(1, 1))
< error
| 42601 syntax error at or near "localtimestamp"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS localtimestamp precision)
< error
| 42601 syntax error at or near "localtimestamp"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS localtimestamp character)
< error
| 42601 syntax error at or near "localtimestamp"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS localtimestamp character(1))
< error
| 42601 syntax error at or near "localtimestamp"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS localtimestamp varying)
< error
| 42601 syntax error at or near "localtimestamp"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS localtimestamp varying(1))
< error
| 42601 syntax error at or near "localtimestamp"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS localtimestamp character varying)
< error
| 42601 syntax error at or near "localtimestamp"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS localtimestamp character varying(1))
< error
| 42601 syntax error at or near "localtimestamp"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS location)
< error
| 42704 type "location" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS lock)
< error
| 42704 type "lock" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS locked)
< error
| 42704 type "locked" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS logged)
< error
| 42704 type "logged" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS mapping)
< error
| 42704 type "mapping" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS match)
< error
| 42704 type "match" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS matched)
< error
| 42704 type "matched" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS materialized)
< error
| 42704 type "materialized" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS maxvalue)
< error
| 42704 type "maxvalue" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS merge)
< error
| 42704 type "merge" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS merge_action)
< error
| 42601 syntax error at or near "merge_action"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS merge_action(1))
< error
| 42601 syntax error at or near "merge_action"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS merge_action(1, 1))
< error
| 42601 syntax error at or near "merge_action"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS merge_action precision)
< error
| 42601 syntax error at or near "merge_action"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS merge_action character)
< error
| 42601 syntax error at or near "merge_action"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS merge_action character(1))
< error
| 42601 syntax error at or near "merge_action"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS merge_action varying)
< error
| 42601 syntax error at or near "merge_action"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS merge_action varying(1))
< error
| 42601 syntax error at or near "merge_action"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS merge_action character varying)
< error
| 42601 syntax error at or near "merge_action"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS merge_action character varying(1))
< error
| 42601 syntax error at or near "merge_action"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS method)
< error
| 42704 type "method" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS minute)
< error
| 42704 type "minute" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS minvalue)
< error
| 42704 type "minvalue" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS mode)
< error
| 42704 type "mode" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS month)
< error
| 42704 type "month" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS move)
< error
| 42704 type "move" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS name)
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS names)
< error
| 42704 type "names" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS national)
< error
| 42601 syntax error at or near ")"
| position: 29
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS national(1))
< error
| 42601 syntax error at or near "("
| position: 29
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS national(1, 1))
< error
| 42601 syntax error at or near "("
| position: 29
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS national precision)
< error
| 42601 syntax error at or near "precision"
| position: 30
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS national character)
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS natural)
< error
| 42704 type "natural" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS nchar)
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS nested)
< error
| 42704 type "nested" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS new)
< error
| 42704 type "new" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS next)
< error
| 42704 type "next" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS nfc)
< error
| 42704 type "nfc" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS nfd)
< error
| 42704 type "nfd" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS nfkc)
< error
| 42704 type "nfkc" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS nfkd)
< error
| 42704 type "nfkd" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS no)
< error
| 42704 type "no" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS none)
< error
| 42601 syntax error at or near "none"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS none(1))
< error
| 42601 syntax error at or near "none"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS none(1, 1))
< error
| 42601 syntax error at or near "none"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS none precision)
< error
| 42601 syntax error at or near "none"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS none character)
< error
| 42601 syntax error at or near "none"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS none character(1))
< error
| 42601 syntax error at or near "none"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS none varying)
< error
| 42601 syntax error at or near "none"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS none varying(1))
< error
| 42601 syntax error at or near "none"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS none character varying)
< error
| 42601 syntax error at or near "none"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS none character varying(1))
< error
| 42601 syntax error at or near "none"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS normalize)
< error
| 42601 syntax error at or near "normalize"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS normalize(1))
< error
| 42601 syntax error at or near "normalize"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS normalize(1, 1))
< error
| 42601 syntax error at or near "normalize"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS normalize precision)
< error
| 42601 syntax error at or near "normalize"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS normalize character)
< error
| 42601 syntax error at or near "normalize"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS normalize character(1))
< error
| 42601 syntax error at or near "normalize"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS normalize varying)
< error
| 42601 syntax error at or near "normalize"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS normalize varying(1))
< error
| 42601 syntax error at or near "normalize"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS normalize character varying)
< error
| 42601 syntax error at or near "normalize"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS normalize character varying(1))
< error
| 42601 syntax error at or near "normalize"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS normalized)
< error
| 42704 type "normalized" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS not)
< error
| 42601 syntax error at or near "not"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS not(1))
< error
| 42601 syntax error at or near "not"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS not(1, 1))
< error
| 42601 syntax error at or near "not"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS not precision)
< error
| 42601 syntax error at or near "not"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS not character)
< error
| 42601 syntax error at or near "not"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS not character(1))
< error
| 42601 syntax error at or near "not"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS not varying)
< error
| 42601 syntax error at or near "not"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS not varying(1))
< error
| 42601 syntax error at or near "not"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS not character varying)
< error
| 42601 syntax error at or near "not"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS not character varying(1))
< error
| 42601 syntax error at or near "not"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS nothing)
< error
| 42704 type "nothing" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS notify)
< error
| 42704 type "notify" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS notnull)
< error
| 42704 type "notnull" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS nowait)
< error
| 42704 type "nowait" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS null)
< error
| 42601 syntax error at or near "null"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS null(1))
< error
| 42601 syntax error at or near "null"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS null(1, 1))
< error
| 42601 syntax error at or near "null"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS null precision)
< error
| 42601 syntax error at or near "null"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS null character)
< error
| 42601 syntax error at or near "null"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS null character(1))
< error
| 42601 syntax error at or near "null"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS null varying)
< error
| 42601 syntax error at or near "null"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS null varying(1))
< error
| 42601 syntax error at or near "null"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS null character varying)
< error
| 42601 syntax error at or near "null"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS null character varying(1))
< error
| 42601 syntax error at or near "null"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS nullif)
< error
| 42601 syntax error at or near "nullif"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS nullif(1))
< error
| 42601 syntax error at or near "nullif"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS nullif(1, 1))
< error
| 42601 syntax error at or near "nullif"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS nullif precision)
< error
| 42601 syntax error at or near "nullif"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS nullif character)
< error
| 42601 syntax error at or near "nullif"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS nullif character(1))
< error
| 42601 syntax error at or near "nullif"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS nullif varying)
< error
| 42601 syntax error at or near "nullif"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS nullif varying(1))
< error
| 42601 syntax error at or near "nullif"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS nullif character varying)
< error
| 42601 syntax error at or near "nullif"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS nullif character varying(1))
< error
| 42601 syntax error at or near "nullif"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS nulls)
< error
| 42704 type "nulls" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS numeric)
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS object)
< error
| 42704 type "object" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS of)
< error
| 42704 type "of" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS off)
< error
| 42704 type "off" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS offset)
< error
| 42601 syntax error at or near "offset"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS offset(1))
< error
| 42601 syntax error at or near "offset"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS offset(1, 1))
< error
| 42601 syntax error at or near "offset"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS offset precision)
< error
| 42601 syntax error at or near "offset"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS offset character)
< error
| 42601 syntax error at or near "offset"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS offset character(1))
< error
| 42601 syntax error at or near "offset"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS offset varying)
< error
| 42601 syntax error at or near "offset"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS offset varying(1))
< error
| 42601 syntax error at or near "offset"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS offset character varying)
< error
| 42601 syntax error at or near "offset"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS offset character varying(1))
< error
| 42601 syntax error at or near "offset"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS oids)
< error
| 42704 type "oids" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS old)
< error
| 42704 type "old" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS omit)
< error
| 42704 type "omit" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS on)
< error
| 42601 syntax error at or near "on"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS on(1))
< error
| 42601 syntax error at or near "on"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS on(1, 1))
< error
| 42601 syntax error at or near "on"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS on precision)
< error
| 42601 syntax error at or near "on"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS on character)
< error
| 42601 syntax error at or near "on"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS on character(1))
< error
| 42601 syntax error at or near "on"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS on varying)
< error
| 42601 syntax error at or near "on"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS on varying(1))
< error
| 42601 syntax error at or near "on"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS on character varying)
< error
| 42601 syntax error at or near "on"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS on character varying(1))
< error
| 42601 syntax error at or near "on"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS only)
< error
| 42601 syntax error at or near "only"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS only(1))
< error
| 42601 syntax error at or near "only"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS only(1, 1))
< error
| 42601 syntax error at or near "only"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS only precision)
< error
| 42601 syntax error at or near "only"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS only character)
< error
| 42601 syntax error at or near "only"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS only character(1))
< error
| 42601 syntax error at or near "only"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS only varying)
< error
| 42601 syntax error at or near "only"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS only varying(1))
< error
| 42601 syntax error at or near "only"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS only character varying)
< error
| 42601 syntax error at or near "only"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS only character varying(1))
< error
| 42601 syntax error at or near "only"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS operator)
< error
| 42704 type "operator" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS option)
< error
| 42704 type "option" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS options)
< error
| 42704 type "options" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS or)
< error
| 42601 syntax error at or near "or"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS or(1))
< error
| 42601 syntax error at or near "or"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS or(1, 1))
< error
| 42601 syntax error at or near "or"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS or precision)
< error
| 42601 syntax error at or near "or"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS or character)
< error
| 42601 syntax error at or near "or"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS or character(1))
< error
| 42601 syntax error at or near "or"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS or varying)
< error
| 42601 syntax error at or near "or"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS or varying(1))
< error
| 42601 syntax error at or near "or"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS or character varying)
< error
| 42601 syntax error at or near "or"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS or character varying(1))
< error
| 42601 syntax error at or near "or"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS order)
< error
| 42601 syntax error at or near "order"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS order(1))
< error
| 42601 syntax error at or near "order"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS order(1, 1))
< error
| 42601 syntax error at or near "order"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS order precision)
< error
| 42601 syntax error at or near "order"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS order character)
< error
| 42601 syntax error at or near "order"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS order character(1))
< error
| 42601 syntax error at or near "order"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS order varying)
< error
| 42601 syntax error at or near "order"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS order varying(1))
< error
| 42601 syntax error at or near "order"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS order character varying)
< error
| 42601 syntax error at or near "order"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS order character varying(1))
< error
| 42601 syntax error at or near "order"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS ordinality)
< error
| 42704 type "ordinality" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS others)
< error
| 42704 type "others" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS out)
< error
| 42601 syntax error at or near "out"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS out(1))
< error
| 42601 syntax error at or near "out"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS out(1, 1))
< error
| 42601 syntax error at or near "out"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS out precision)
< error
| 42601 syntax error at or near "out"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS out character)
< error
| 42601 syntax error at or near "out"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS out character(1))
< error
| 42601 syntax error at or near "out"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS out varying)
< error
| 42601 syntax error at or near "out"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS out varying(1))
< error
| 42601 syntax error at or near "out"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS out character varying)
< error
| 42601 syntax error at or near "out"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS out character varying(1))
< error
| 42601 syntax error at or near "out"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS outer)
< error
| 42704 type "outer" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS over)
< error
| 42704 type "over" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS overlaps)
< error
| 42704 type "overlaps" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS overlay)
< error
| 42601 syntax error at or near "overlay"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS overlay(1))
< error
| 42601 syntax error at or near "overlay"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS overlay(1, 1))
< error
| 42601 syntax error at or near "overlay"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS overlay precision)
< error
| 42601 syntax error at or near "overlay"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS overlay character)
< error
| 42601 syntax error at or near "overlay"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS overlay character(1))
< error
| 42601 syntax error at or near "overlay"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS overlay varying)
< error
| 42601 syntax error at or near "overlay"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS overlay varying(1))
< error
| 42601 syntax error at or near "overlay"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS overlay character varying)
< error
| 42601 syntax error at or near "overlay"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS overlay character varying(1))
< error
| 42601 syntax error at or near "overlay"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS overriding)
< error
| 42704 type "overriding" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS owned)
< error
| 42704 type "owned" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS owner)
< error
| 42704 type "owner" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS parallel)
< error
| 42704 type "parallel" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS parameter)
< error
| 42704 type "parameter" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS parser)
< error
| 42704 type "parser" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS partial)
< error
| 42704 type "partial" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS partition)
< error
| 42704 type "partition" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS passing)
< error
| 42704 type "passing" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS password)
< error
| 42704 type "password" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS path)
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS placing)
< error
| 42601 syntax error at or near "placing"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS placing(1))
< error
| 42601 syntax error at or near "placing"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS placing(1, 1))
< error
| 42601 syntax error at or near "placing"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS placing precision)
< error
| 42601 syntax error at or near "placing"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS placing character)
< error
| 42601 syntax error at or near "placing"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS placing character(1))
< error
| 42601 syntax error at or near "placing"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS placing varying)
< error
| 42601 syntax error at or near "placing"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS placing varying(1))
< error
| 42601 syntax error at or near "placing"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS placing character varying)
< error
| 42601 syntax error at or near "placing"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS placing character varying(1))
< error
| 42601 syntax error at or near "placing"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS plan)
< error
| 42704 type "plan" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS plans)
< error
| 42704 type "plans" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS policy)
< error
| 42704 type "policy" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS position)
< error
| 42601 syntax error at or near "position"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS position(1))
< error
| 42601 syntax error at or near "position"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS position(1, 1))
< error
| 42601 syntax error at or near "position"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS position precision)
< error
| 42601 syntax error at or near "position"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS position character)
< error
| 42601 syntax error at or near "position"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS position character(1))
< error
| 42601 syntax error at or near "position"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS position varying)
< error
| 42601 syntax error at or near "position"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS position varying(1))
< error
| 42601 syntax error at or near "position"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS position character varying)
< error
| 42601 syntax error at or near "position"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS position character varying(1))
< error
| 42601 syntax error at or near "position"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS preceding)
< error
| 42704 type "preceding" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS precision)
< error
| 42601 syntax error at or near "precision"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS precision(1))
< error
| 42601 syntax error at or near "precision"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS precision(1, 1))
< error
| 42601 syntax error at or near "precision"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS precision precision)
< error
| 42601 syntax error at or near "precision"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS precision character)
< error
| 42601 syntax error at or near "precision"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS precision character(1))
< error
| 42601 syntax error at or near "precision"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS precision varying)
< error
| 42601 syntax error at or near "precision"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS precision varying(1))
< error
| 42601 syntax error at or near "precision"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS precision character varying)
< error
| 42601 syntax error at or near "precision"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS precision character varying(1))
< error
| 42601 syntax error at or near "precision"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS prepare)
< error
| 42704 type "prepare" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS prepared)
< error
| 42704 type "prepared" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS preserve)
< error
| 42704 type "preserve" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS primary)
< error
| 42601 syntax error at or near "primary"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS primary(1))
< error
| 42601 syntax error at or near "primary"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS primary(1, 1))
< error
| 42601 syntax error at or near "primary"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS primary precision)
< error
| 42601 syntax error at or near "primary"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS primary character)
< error
| 42601 syntax error at or near "primary"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS primary character(1))
< error
| 42601 syntax error at or near "primary"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS primary varying)
< error
| 42601 syntax error at or near "primary"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS primary varying(1))
< error
| 42601 syntax error at or near "primary"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS primary character varying)
< error
| 42601 syntax error at or near "primary"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS primary character varying(1))
< error
| 42601 syntax error at or near "primary"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS prior)
< error
| 42704 type "prior" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS privileges)
< error
| 42704 type "privileges" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS procedural)
< error
| 42704 type "procedural" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS procedure)
< error
| 42704 type "procedure" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS procedures)
< error
| 42704 type "procedures" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS program)
< error
| 42704 type "program" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS publication)
< error
| 42704 type "publication" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS quote)
< error
| 42704 type "quote" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS quotes)
< error
| 42704 type "quotes" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS range)
< error
| 42704 type "range" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS read)
< error
| 42704 type "read" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS real)
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS reassign)
< error
| 42704 type "reassign" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS recheck)
< error
| 42704 type "recheck" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS recursive)
< error
| 42704 type "recursive" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS ref)
< error
| 42704 type "ref" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS references)
< error
| 42601 syntax error at or near "references"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS references(1))
< error
| 42601 syntax error at or near "references"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS references(1, 1))
< error
| 42601 syntax error at or near "references"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS references precision)
< error
| 42601 syntax error at or near "references"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS references character)
< error
| 42601 syntax error at or near "references"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS references character(1))
< error
| 42601 syntax error at or near "references"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS references varying)
< error
| 42601 syntax error at or near "references"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS references varying(1))
< error
| 42601 syntax error at or near "references"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS references character varying)
< error
| 42601 syntax error at or near "references"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS references character varying(1))
< error
| 42601 syntax error at or near "references"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS referencing)
< error
| 42704 type "referencing" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS refresh)
< error
| 42704 type "refresh" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS reindex)
< error
| 42704 type "reindex" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS relative)
< error
| 42704 type "relative" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS release)
< error
| 42704 type "release" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS rename)
< error
| 42704 type "rename" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS repeatable)
< error
| 42704 type "repeatable" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS replace)
< error
| 42704 type "replace" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS replica)
< error
| 42704 type "replica" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS reset)
< error
| 42704 type "reset" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS restart)
< error
| 42704 type "restart" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS restrict)
< error
| 42704 type "restrict" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS return)
< error
| 42704 type "return" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS returning)
< error
| 42601 syntax error at or near "returning"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS returning(1))
< error
| 42601 syntax error at or near "returning"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS returning(1, 1))
< error
| 42601 syntax error at or near "returning"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS returning precision)
< error
| 42601 syntax error at or near "returning"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS returning character)
< error
| 42601 syntax error at or near "returning"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS returning character(1))
< error
| 42601 syntax error at or near "returning"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS returning varying)
< error
| 42601 syntax error at or near "returning"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS returning varying(1))
< error
| 42601 syntax error at or near "returning"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS returning character varying)
< error
| 42601 syntax error at or near "returning"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS returning character varying(1))
< error
| 42601 syntax error at or near "returning"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS returns)
< error
| 42704 type "returns" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS revoke)
< error
| 42704 type "revoke" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS right)
< error
| 42704 type "right" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS role)
< error
| 42704 type "role" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS rollback)
< error
| 42704 type "rollback" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS rollup)
< error
| 42704 type "rollup" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS routine)
< error
| 42704 type "routine" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS routines)
< error
| 42704 type "routines" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS row)
< error
| 42601 syntax error at or near "row"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS row(1))
< error
| 42601 syntax error at or near "row"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS row(1, 1))
< error
| 42601 syntax error at or near "row"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS row precision)
< error
| 42601 syntax error at or near "row"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS row character)
< error
| 42601 syntax error at or near "row"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS row character(1))
< error
| 42601 syntax error at or near "row"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS row varying)
< error
| 42601 syntax error at or near "row"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS row varying(1))
< error
| 42601 syntax error at or near "row"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS row character varying)
< error
| 42601 syntax error at or near "row"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS row character varying(1))
< error
| 42601 syntax error at or near "row"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS rows)
< error
| 42704 type "rows" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS rule)
< error
| 42704 type "rule" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS savepoint)
< error
| 42704 type "savepoint" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS scalar)
< error
| 42704 type "scalar" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS schema)
< error
| 42704 type "schema" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS schemas)
< error
| 42704 type "schemas" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS scroll)
< error
| 42704 type "scroll" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS search)
< error
| 42704 type "search" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS second)
< error
| 42704 type "second" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS security)
< error
| 42704 type "security" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS select)
< error
| 42601 syntax error at or near "select"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS select(1))
< error
| 42601 syntax error at or near "select"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS select(1, 1))
< error
| 42601 syntax error at or near "select"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS select precision)
< error
| 42601 syntax error at or near "select"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS select character)
< error
| 42601 syntax error at or near "select"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS select character(1))
< error
| 42601 syntax error at or near "select"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS select varying)
< error
| 42601 syntax error at or near "select"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS select varying(1))
< error
| 42601 syntax error at or near "select"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS select character varying)
< error
| 42601 syntax error at or near "select"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS select character varying(1))
< error
| 42601 syntax error at or near "select"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS sequence)
< error
| 42704 type "sequence" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS sequences)
< error
| 42704 type "sequences" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS serializable)
< error
| 42704 type "serializable" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS server)
< error
| 42704 type "server" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS session)
< error
| 42704 type "session" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS session_user)
< error
| 42601 syntax error at or near "session_user"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS session_user(1))
< error
| 42601 syntax error at or near "session_user"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS session_user(1, 1))
< error
| 42601 syntax error at or near "session_user"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS session_user precision)
< error
| 42601 syntax error at or near "session_user"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS session_user character)
< error
| 42601 syntax error at or near "session_user"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS session_user character(1))
< error
| 42601 syntax error at or near "session_user"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS session_user varying)
< error
| 42601 syntax error at or near "session_user"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS session_user varying(1))
< error
| 42601 syntax error at or near "session_user"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS session_user character varying)
< error
| 42601 syntax error at or near "session_user"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS session_user character varying(1))
< error
| 42601 syntax error at or near "session_user"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS set)
< error
| 42704 type "set" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS setof int)
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS sets)
< error
| 42704 type "sets" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS share)
< error
| 42704 type "share" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS show)
< error
| 42704 type "show" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS similar)
< error
| 42704 type "similar" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS simple)
< error
| 42704 type "simple" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS skip)
< error
| 42704 type "skip" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS smallint)
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS snapshot)
< error
| 42704 type "snapshot" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS some)
< error
| 42601 syntax error at or near "some"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS some(1))
< error
| 42601 syntax error at or near "some"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS some(1, 1))
< error
| 42601 syntax error at or near "some"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS some precision)
< error
| 42601 syntax error at or near "some"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS some character)
< error
| 42601 syntax error at or near "some"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS some character(1))
< error
| 42601 syntax error at or near "some"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS some varying)
< error
| 42601 syntax error at or near "some"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS some varying(1))
< error
| 42601 syntax error at or near "some"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS some character varying)
< error
| 42601 syntax error at or near "some"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS some character varying(1))
< error
| 42601 syntax error at or near "some"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS source)
< error
| 42704 type "source" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS sql)
< error
| 42704 type "sql" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS stable)
< error
| 42704 type "stable" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS standalone)
< error
| 42704 type "standalone" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS start)
< error
| 42704 type "start" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS statement)
< error
| 42704 type "statement" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS statistics)
< error
| 42704 type "statistics" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS stdin)
< error
| 42704 type "stdin" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS stdout)
< error
| 42704 type "stdout" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS storage)
< error
| 42704 type "storage" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS stored)
< error
| 42704 type "stored" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS strict)
< error
| 42704 type "strict" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS string)
< error
| 42704 type "string" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS strip)
< error
| 42704 type "strip" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS subscription)
< error
| 42704 type "subscription" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS substring)
< error
| 42601 syntax error at or near "substring"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS substring(1))
< error
| 42601 syntax error at or near "substring"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS substring(1, 1))
< error
| 42601 syntax error at or near "substring"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS substring precision)
< error
| 42601 syntax error at or near "substring"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS substring character)
< error
| 42601 syntax error at or near "substring"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS substring character(1))
< error
| 42601 syntax error at or near "substring"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS substring varying)
< error
| 42601 syntax error at or near "substring"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS substring varying(1))
< error
| 42601 syntax error at or near "substring"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS substring character varying)
< error
| 42601 syntax error at or near "substring"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS substring character varying(1))
< error
| 42601 syntax error at or near "substring"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS support)
< error
| 42704 type "support" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS symmetric)
< error
| 42601 syntax error at or near "symmetric"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS symmetric(1))
< error
| 42601 syntax error at or near "symmetric"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS symmetric(1, 1))
< error
| 42601 syntax error at or near "symmetric"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS symmetric precision)
< error
| 42601 syntax error at or near "symmetric"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS symmetric character)
< error
| 42601 syntax error at or near "symmetric"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS symmetric character(1))
< error
| 42601 syntax error at or near "symmetric"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS symmetric varying)
< error
| 42601 syntax error at or near "symmetric"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS symmetric varying(1))
< error
| 42601 syntax error at or near "symmetric"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS symmetric character varying)
< error
| 42601 syntax error at or near "symmetric"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS symmetric character varying(1))
< error
| 42601 syntax error at or near "symmetric"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS sysid)
< error
| 42704 type "sysid" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS system)
< error
| 42704 type "system" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS system_user)
< error
| 42601 syntax error at or near "system_user"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS system_user(1))
< error
| 42601 syntax error at or near "system_user"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS system_user(1, 1))
< error
| 42601 syntax error at or near "system_user"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS system_user precision)
< error
| 42601 syntax error at or near "system_user"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS system_user character)
< error
| 42601 syntax error at or near "system_user"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS system_user character(1))
< error
| 42601 syntax error at or near "system_user"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS system_user varying)
< error
| 42601 syntax error at or near "system_user"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS system_user varying(1))
< error
| 42601 syntax error at or near "system_user"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS system_user character varying)
< error
| 42601 syntax error at or near "system_user"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS system_user character varying(1))
< error
| 42601 syntax error at or near "system_user"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS table)
< error
| 42601 syntax error at or near "table"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS table(1))
< error
| 42601 syntax error at or near "table"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS table(1, 1))
< error
| 42601 syntax error at or near "table"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS table precision)
< error
| 42601 syntax error at or near "table"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS table character)
< error
| 42601 syntax error at or near "table"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS table character(1))
< error
| 42601 syntax error at or near "table"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS table varying)
< error
| 42601 syntax error at or near "table"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS table varying(1))
< error
| 42601 syntax error at or near "table"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS table character varying)
< error
| 42601 syntax error at or near "table"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS table character varying(1))
< error
| 42601 syntax error at or near "table"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS tables)
< error
| 42704 type "tables" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS tablesample)
< error
| 42704 type "tablesample" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS tablespace)
< error
| 42704 type "tablespace" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS target)
< error
| 42704 type "target" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS temp)
< error
| 42704 type "temp" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS template)
< error
| 42704 type "template" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS temporary)
< error
| 42704 type "temporary" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS text)
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS then)
< error
| 42601 syntax error at or near "then"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS then(1))
< error
| 42601 syntax error at or near "then"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS then(1, 1))
< error
| 42601 syntax error at or near "then"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS then precision)
< error
| 42601 syntax error at or near "then"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS then character)
< error
| 42601 syntax error at or near "then"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS then character(1))
< error
| 42601 syntax error at or near "then"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS then varying)
< error
| 42601 syntax error at or near "then"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS then varying(1))
< error
| 42601 syntax error at or near "then"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS then character varying)
< error
| 42601 syntax error at or near "then"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS then character varying(1))
< error
| 42601 syntax error at or near "then"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS ties)
< error
| 42704 type "ties" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS time)
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS timestamp)
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS to)
< error
| 42601 syntax error at or near "to"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS to(1))
< error
| 42601 syntax error at or near "to"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS to(1, 1))
< error
| 42601 syntax error at or near "to"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS to precision)
< error
| 42601 syntax error at or near "to"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS to character)
< error
| 42601 syntax error at or near "to"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS to character(1))
< error
| 42601 syntax error at or near "to"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS to varying)
< error
| 42601 syntax error at or near "to"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS to varying(1))
< error
| 42601 syntax error at or near "to"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS to character varying)
< error
| 42601 syntax error at or near "to"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS to character varying(1))
< error
| 42601 syntax error at or near "to"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS trailing)
< error
| 42601 syntax error at or near "trailing"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS trailing(1))
< error
| 42601 syntax error at or near "trailing"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS trailing(1, 1))
< error
| 42601 syntax error at or near "trailing"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS trailing precision)
< error
| 42601 syntax error at or near "trailing"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS trailing character)
< error
| 42601 syntax error at or near "trailing"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS trailing character(1))
< error
| 42601 syntax error at or near "trailing"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS trailing varying)
< error
| 42601 syntax error at or near "trailing"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS trailing varying(1))
< error
| 42601 syntax error at or near "trailing"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS trailing character varying)
< error
| 42601 syntax error at or near "trailing"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS trailing character varying(1))
< error
| 42601 syntax error at or near "trailing"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS transaction)
< error
| 42704 type "transaction" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS transform)
< error
| 42704 type "transform" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS treat)
< error
| 42601 syntax error at or near "treat"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS treat(1))
< error
| 42601 syntax error at or near "treat"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS treat(1, 1))
< error
| 42601 syntax error at or near "treat"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS treat precision)
< error
| 42601 syntax error at or near "treat"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS treat character)
< error
| 42601 syntax error at or near "treat"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS treat character(1))
< error
| 42601 syntax error at or near "treat"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS treat varying)
< error
| 42601 syntax error at or near "treat"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS treat varying(1))
< error
| 42601 syntax error at or near "treat"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS treat character varying)
< error
| 42601 syntax error at or near "treat"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS treat character varying(1))
< error
| 42601 syntax error at or near "treat"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS trigger)
< error
| 0A000 cannot accept a value of type trigger
| position: 13
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS trim)
< error
| 42601 syntax error at or near "trim"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS trim(1))
< error
| 42601 syntax error at or near "trim"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS trim(1, 1))
< error
| 42601 syntax error at or near "trim"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS trim precision)
< error
| 42601 syntax error at or near "trim"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS trim character)
< error
| 42601 syntax error at or near "trim"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS trim character(1))
< error
| 42601 syntax error at or near "trim"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS trim varying)
< error
| 42601 syntax error at or near "trim"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS trim varying(1))
< error
| 42601 syntax error at or near "trim"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS trim character varying)
< error
| 42601 syntax error at or near "trim"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS trim character varying(1))
< error
| 42601 syntax error at or near "trim"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS true)
< error
| 42601 syntax error at or near "true"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS true(1))
< error
| 42601 syntax error at or near "true"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS true(1, 1))
< error
| 42601 syntax error at or near "true"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS true precision)
< error
| 42601 syntax error at or near "true"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS true character)
< error
| 42601 syntax error at or near "true"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS true character(1))
< error
| 42601 syntax error at or near "true"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS true varying)
< error
| 42601 syntax error at or near "true"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS true varying(1))
< error
| 42601 syntax error at or near "true"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS true character varying)
< error
| 42601 syntax error at or near "true"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS true character varying(1))
< error
| 42601 syntax error at or near "true"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS truncate)
< error
| 42704 type "truncate" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS trusted)
< error
| 42704 type "trusted" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS type)
< error
| 42704 type "type" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS types)
< error
| 42704 type "types" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS uescape)
< error
| 42704 type "uescape" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS unbounded)
< error
| 42704 type "unbounded" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS uncommitted)
< error
| 42704 type "uncommitted" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS unconditional)
< error
| 42704 type "unconditional" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS unencrypted)
< error
| 42704 type "unencrypted" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS union)
< error
| 42601 syntax error at or near "union"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS union(1))
< error
| 42601 syntax error at or near "union"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS union(1, 1))
< error
| 42601 syntax error at or near "union"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS union precision)
< error
| 42601 syntax error at or near "union"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS union character)
< error
| 42601 syntax error at or near "union"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS union character(1))
< error
| 42601 syntax error at or near "union"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS union varying)
< error
| 42601 syntax error at or near "union"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS union varying(1))
< error
| 42601 syntax error at or near "union"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS union character varying)
< error
| 42601 syntax error at or near "union"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS union character varying(1))
< error
| 42601 syntax error at or near "union"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS unique)
< error
| 42601 syntax error at or near "unique"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS unique(1))
< error
| 42601 syntax error at or near "unique"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS unique(1, 1))
< error
| 42601 syntax error at or near "unique"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS unique precision)
< error
| 42601 syntax error at or near "unique"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS unique character)
< error
| 42601 syntax error at or near "unique"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS unique character(1))
< error
| 42601 syntax error at or near "unique"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS unique varying)
< error
| 42601 syntax error at or near "unique"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS unique varying(1))
< error
| 42601 syntax error at or near "unique"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS unique character varying)
< error
| 42601 syntax error at or near "unique"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS unique character varying(1))
< error
| 42601 syntax error at or near "unique"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS unknown)
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS unlisten)
< error
| 42704 type "unlisten" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS unlogged)
< error
| 42704 type "unlogged" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS until)
< error
| 42704 type "until" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS update)
< error
| 42704 type "update" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS user)
< error
| 42601 syntax error at or near "user"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS user(1))
< error
| 42601 syntax error at or near "user"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS user(1, 1))
< error
| 42601 syntax error at or near "user"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS user precision)
< error
| 42601 syntax error at or near "user"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS user character)
< error
| 42601 syntax error at or near "user"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS user character(1))
< error
| 42601 syntax error at or near "user"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS user varying)
< error
| 42601 syntax error at or near "user"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS user varying(1))
< error
| 42601 syntax error at or near "user"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS user character varying)
< error
| 42601 syntax error at or near "user"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS user character varying(1))
< error
| 42601 syntax error at or near "user"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS using)
< error
| 42601 syntax error at or near "using"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS using(1))
< error
| 42601 syntax error at or near "using"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS using(1, 1))
< error
| 42601 syntax error at or near "using"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS using precision)
< error
| 42601 syntax error at or near "using"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS using character)
< error
| 42601 syntax error at or near "using"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS using character(1))
< error
| 42601 syntax error at or near "using"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS using varying)
< error
| 42601 syntax error at or near "using"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS using varying(1))
< error
| 42601 syntax error at or near "using"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS using character varying)
< error
| 42601 syntax error at or near "using"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS using character varying(1))
< error
| 42601 syntax error at or near "using"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS vacuum)
< error
| 42704 type "vacuum" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS valid)
< error
| 42704 type "valid" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS validate)
< error
| 42704 type "validate" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS validator)
< error
| 42704 type "validator" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS value)
< error
| 42704 type "value" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS values)
< error
| 42601 syntax error at or near "values"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS values(1))
< error
| 42601 syntax error at or near "values"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS values(1, 1))
< error
| 42601 syntax error at or near "values"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS values precision)
< error
| 42601 syntax error at or near "values"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS values character)
< error
| 42601 syntax error at or near "values"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS values character(1))
< error
| 42601 syntax error at or near "values"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS values varying)
< error
| 42601 syntax error at or near "values"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS values varying(1))
< error
| 42601 syntax error at or near "values"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS values character varying)
< error
| 42601 syntax error at or near "values"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS values character varying(1))
< error
| 42601 syntax error at or near "values"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS varchar)
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS variadic)
< error
| 42601 syntax error at or near "variadic"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS variadic(1))
< error
| 42601 syntax error at or near "variadic"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS variadic(1, 1))
< error
| 42601 syntax error at or near "variadic"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS variadic precision)
< error
| 42601 syntax error at or near "variadic"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS variadic character)
< error
| 42601 syntax error at or near "variadic"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS variadic character(1))
< error
| 42601 syntax error at or near "variadic"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS variadic varying)
< error
| 42601 syntax error at or near "variadic"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS variadic varying(1))
< error
| 42601 syntax error at or near "variadic"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS variadic character varying)
< error
| 42601 syntax error at or near "variadic"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS variadic character varying(1))
< error
| 42601 syntax error at or near "variadic"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS varying)
< error
| 42704 type "varying" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS verbose)
< error
| 42704 type "verbose" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS version)
< error
| 42704 type "version" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS view)
< error
| 42704 type "view" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS views)
< error
| 42704 type "views" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS volatile)
< error
| 42704 type "volatile" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS when)
< error
| 42601 syntax error at or near "when"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS when(1))
< error
| 42601 syntax error at or near "when"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS when(1, 1))
< error
| 42601 syntax error at or near "when"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS when precision)
< error
| 42601 syntax error at or near "when"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS when character)
< error
| 42601 syntax error at or near "when"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS when character(1))
< error
| 42601 syntax error at or near "when"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS when varying)
< error
| 42601 syntax error at or near "when"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS when varying(1))
< error
| 42601 syntax error at or near "when"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS when character varying)
< error
| 42601 syntax error at or near "when"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS when character varying(1))
< error
| 42601 syntax error at or near "when"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS where)
< error
| 42601 syntax error at or near "where"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS where(1))
< error
| 42601 syntax error at or near "where"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS where(1, 1))
< error
| 42601 syntax error at or near "where"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS where precision)
< error
| 42601 syntax error at or near "where"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS where character)
< error
| 42601 syntax error at or near "where"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS where character(1))
< error
| 42601 syntax error at or near "where"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS where varying)
< error
| 42601 syntax error at or near "where"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS where varying(1))
< error
| 42601 syntax error at or near "where"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS where character varying)
< error
| 42601 syntax error at or near "where"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS where character varying(1))
< error
| 42601 syntax error at or near "where"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS whitespace)
< error
| 42704 type "whitespace" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS window)
< error
| 42601 syntax error at or near "window"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS window(1))
< error
| 42601 syntax error at or near "window"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS window(1, 1))
< error
| 42601 syntax error at or near "window"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS window precision)
< error
| 42601 syntax error at or near "window"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS window character)
< error
| 42601 syntax error at or near "window"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS window character(1))
< error
| 42601 syntax error at or near "window"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS window varying)
< error
| 42601 syntax error at or near "window"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS window varying(1))
< error
| 42601 syntax error at or near "window"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS window character varying)
< error
| 42601 syntax error at or near "window"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS window character varying(1))
< error
| 42601 syntax error at or near "window"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS with)
< error
| 42601 syntax error at or near "with"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS with(1))
< error
| 42601 syntax error at or near "with"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS with(1, 1))
< error
| 42601 syntax error at or near "with"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS with precision)
< error
| 42601 syntax error at or near "with"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS with character)
< error
| 42601 syntax error at or near "with"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS with character(1))
< error
| 42601 syntax error at or near "with"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS with varying)
< error
| 42601 syntax error at or near "with"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS with varying(1))
< error
| 42601 syntax error at or near "with"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS with character varying)
< error
| 42601 syntax error at or near "with"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS with character varying(1))
< error
| 42601 syntax error at or near "with"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS within)
< error
| 42704 type "within" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS without)
< error
| 42704 type "without" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS work)
< error
| 42704 type "work" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS wrapper)
< error
| 42704 type "wrapper" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS write)
< error
| 42704 type "write" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xml)
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlattributes)
< error
| 42601 syntax error at or near "xmlattributes"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlattributes(1))
< error
| 42601 syntax error at or near "xmlattributes"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlattributes(1, 1))
< error
| 42601 syntax error at or near "xmlattributes"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlattributes precision)
< error
| 42601 syntax error at or near "xmlattributes"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlattributes character)
< error
| 42601 syntax error at or near "xmlattributes"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlattributes character(1))
< error
| 42601 syntax error at or near "xmlattributes"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlattributes varying)
< error
| 42601 syntax error at or near "xmlattributes"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlattributes varying(1))
< error
| 42601 syntax error at or near "xmlattributes"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlattributes character varying)
< error
| 42601 syntax error at or near "xmlattributes"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlattributes character varying(1))
< error
| 42601 syntax error at or near "xmlattributes"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlconcat)
< error
| 42601 syntax error at or near "xmlconcat"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlconcat(1))
< error
| 42601 syntax error at or near "xmlconcat"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlconcat(1, 1))
< error
| 42601 syntax error at or near "xmlconcat"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlconcat precision)
< error
| 42601 syntax error at or near "xmlconcat"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlconcat character)
< error
| 42601 syntax error at or near "xmlconcat"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlconcat character(1))
< error
| 42601 syntax error at or near "xmlconcat"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlconcat varying)
< error
| 42601 syntax error at or near "xmlconcat"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlconcat varying(1))
< error
| 42601 syntax error at or near "xmlconcat"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlconcat character varying)
< error
| 42601 syntax error at or near "xmlconcat"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlconcat character varying(1))
< error
| 42601 syntax error at or near "xmlconcat"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlelement)
< error
| 42601 syntax error at or near "xmlelement"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlelement(1))
< error
| 42601 syntax error at or near "xmlelement"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlelement(1, 1))
< error
| 42601 syntax error at or near "xmlelement"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlelement precision)
< error
| 42601 syntax error at or near "xmlelement"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlelement character)
< error
| 42601 syntax error at or near "xmlelement"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlelement character(1))
< error
| 42601 syntax error at or near "xmlelement"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlelement varying)
< error
| 42601 syntax error at or near "xmlelement"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlelement varying(1))
< error
| 42601 syntax error at or near "xmlelement"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlelement character varying)
< error
| 42601 syntax error at or near "xmlelement"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlelement character varying(1))
< error
| 42601 syntax error at or near "xmlelement"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlexists)
< error
| 42601 syntax error at or near "xmlexists"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlexists(1))
< error
| 42601 syntax error at or near "xmlexists"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlexists(1, 1))
< error
| 42601 syntax error at or near "xmlexists"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlexists precision)
< error
| 42601 syntax error at or near "xmlexists"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlexists character)
< error
| 42601 syntax error at or near "xmlexists"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlexists character(1))
< error
| 42601 syntax error at or near "xmlexists"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlexists varying)
< error
| 42601 syntax error at or near "xmlexists"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlexists varying(1))
< error
| 42601 syntax error at or near "xmlexists"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlexists character varying)
< error
| 42601 syntax error at or near "xmlexists"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlexists character varying(1))
< error
| 42601 syntax error at or near "xmlexists"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlforest)
< error
| 42601 syntax error at or near "xmlforest"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlforest(1))
< error
| 42601 syntax error at or near "xmlforest"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlforest(1, 1))
< error
| 42601 syntax error at or near "xmlforest"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlforest precision)
< error
| 42601 syntax error at or near "xmlforest"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlforest character)
< error
| 42601 syntax error at or near "xmlforest"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlforest character(1))
< error
| 42601 syntax error at or near "xmlforest"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlforest varying)
< error
| 42601 syntax error at or near "xmlforest"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlforest varying(1))
< error
| 42601 syntax error at or near "xmlforest"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlforest character varying)
< error
| 42601 syntax error at or near "xmlforest"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlforest character varying(1))
< error
| 42601 syntax error at or near "xmlforest"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlnamespaces)
< error
| 42601 syntax error at or near "xmlnamespaces"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlnamespaces(1))
< error
| 42601 syntax error at or near "xmlnamespaces"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlnamespaces(1, 1))
< error
| 42601 syntax error at or near "xmlnamespaces"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlnamespaces precision)
< error
| 42601 syntax error at or near "xmlnamespaces"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlnamespaces character)
< error
| 42601 syntax error at or near "xmlnamespaces"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlnamespaces character(1))
< error
| 42601 syntax error at or near "xmlnamespaces"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlnamespaces varying)
< error
| 42601 syntax error at or near "xmlnamespaces"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlnamespaces varying(1))
< error
| 42601 syntax error at or near "xmlnamespaces"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlnamespaces character varying)
< error
| 42601 syntax error at or near "xmlnamespaces"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlnamespaces character varying(1))
< error
| 42601 syntax error at or near "xmlnamespaces"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlparse)
< error
| 42601 syntax error at or near "xmlparse"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlparse(1))
< error
| 42601 syntax error at or near "xmlparse"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlparse(1, 1))
< error
| 42601 syntax error at or near "xmlparse"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlparse precision)
< error
| 42601 syntax error at or near "xmlparse"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlparse character)
< error
| 42601 syntax error at or near "xmlparse"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlparse character(1))
< error
| 42601 syntax error at or near "xmlparse"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlparse varying)
< error
| 42601 syntax error at or near "xmlparse"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlparse varying(1))
< error
| 42601 syntax error at or near "xmlparse"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlparse character varying)
< error
| 42601 syntax error at or near "xmlparse"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlparse character varying(1))
< error
| 42601 syntax error at or near "xmlparse"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlpi)
< error
| 42601 syntax error at or near "xmlpi"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlpi(1))
< error
| 42601 syntax error at or near "xmlpi"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlpi(1, 1))
< error
| 42601 syntax error at or near "xmlpi"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlpi precision)
< error
| 42601 syntax error at or near "xmlpi"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlpi character)
< error
| 42601 syntax error at or near "xmlpi"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlpi character(1))
< error
| 42601 syntax error at or near "xmlpi"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlpi varying)
< error
| 42601 syntax error at or near "xmlpi"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlpi varying(1))
< error
| 42601 syntax error at or near "xmlpi"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlpi character varying)
< error
| 42601 syntax error at or near "xmlpi"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlpi character varying(1))
< error
| 42601 syntax error at or near "xmlpi"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlroot)
< error
| 42601 syntax error at or near "xmlroot"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlroot(1))
< error
| 42601 syntax error at or near "xmlroot"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlroot(1, 1))
< error
| 42601 syntax error at or near "xmlroot"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlroot precision)
< error
| 42601 syntax error at or near "xmlroot"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlroot character)
< error
| 42601 syntax error at or near "xmlroot"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlroot character(1))
< error
| 42601 syntax error at or near "xmlroot"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlroot varying)
< error
| 42601 syntax error at or near "xmlroot"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlroot varying(1))
< error
| 42601 syntax error at or near "xmlroot"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlroot character varying)
< error
| 42601 syntax error at or near "xmlroot"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlroot character varying(1))
< error
| 42601 syntax error at or near "xmlroot"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlserialize)
< error
| 42601 syntax error at or near "xmlserialize"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlserialize(1))
< error
| 42601 syntax error at or near "xmlserialize"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlserialize(1, 1))
< error
| 42601 syntax error at or near "xmlserialize"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlserialize precision)
< error
| 42601 syntax error at or near "xmlserialize"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlserialize character)
< error
| 42601 syntax error at or near "xmlserialize"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlserialize character(1))
< error
| 42601 syntax error at or near "xmlserialize"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlserialize varying)
< error
| 42601 syntax error at or near "xmlserialize"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlserialize varying(1))
< error
| 42601 syntax error at or near "xmlserialize"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlserialize character varying)
< error
| 42601 syntax error at or near "xmlserialize"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmlserialize character varying(1))
< error
| 42601 syntax error at or near "xmlserialize"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmltable)
< error
| 42601 syntax error at or near "xmltable"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmltable(1))
< error
| 42601 syntax error at or near "xmltable"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmltable(1, 1))
< error
| 42601 syntax error at or near "xmltable"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmltable precision)
< error
| 42601 syntax error at or near "xmltable"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmltable character)
< error
| 42601 syntax error at or near "xmltable"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmltable character(1))
< error
| 42601 syntax error at or near "xmltable"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmltable varying)
< error
| 42601 syntax error at or near "xmltable"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmltable varying(1))
< error
| 42601 syntax error at or near "xmltable"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmltable character varying)
< error
| 42601 syntax error at or near "xmltable"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS xmltable character varying(1))
< error
| 42601 syntax error at or near "xmltable"
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS year)
< error
| 42704 type "year" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS yes)
< error
| 42704 type "yes" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS zone)
< error
| 42704 type "zone" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS omni_probe_ident)
< error
| 42704 type "omni_probe_ident" does not exist
| position: 21
> rollback
< ok
