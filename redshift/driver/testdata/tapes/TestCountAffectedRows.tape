-- pin: driver aws:redshift:1.0.539507 sha256:afc4041fe6e35ce273970308e430eb8245ea9c4bbf597c497bb8eb8bb00810d2 2026-11-17

== TestCountAffectedRows driver#1
> begin
< ok
> exec
|
| 		DROP TABLE IF EXISTS public.omni_stats_probe;
< ok 0 -
> exec
|
| 		CREATE TABLE public.omni_stats_probe (a int);
< ok 0 -
> exec
|
| 		INSERT INTO public.omni_stats_probe
| 			WITH RECURSIVE numbers(n) AS (SELECT 1 UNION ALL SELECT n + 1 FROM numbers WHERE n < 100)
| 			SELECT 1 FROM numbers;
< ok 100 -
> exec
|
| 		ANALYZE public.omni_stats_probe;
< ok 0 -
> commit
< ok
> query
| EXPLAIN DELETE FROM public.omni_stats_probe WHERE a = 1
< result
< column "QUERY PLAN" "VARCHAR" "string"
< row
: string
| XN Seq Scan on omni_stats_probe  (cost=0.00..1.25 rows=100 width=6)
< row
: string
|   Filter: (a = 1)
> exec
| DROP TABLE IF EXISTS public.omni_stats_probe
< ok 0 -
