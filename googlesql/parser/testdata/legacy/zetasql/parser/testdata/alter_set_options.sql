ALTER MATERIALIZED VIEW adsdw.foo SET OPTIONS(), SET OPTIONS(), SET OPTIONS();
ALTER APPROX VIEW IF EXISTS adsdw.foo SET OPTIONS();
ALTER APPROX VIEW adsdw.foo SET OPTIONS(quota_accounting_owner = 'adsdw-etl@google.com');
ALTER APPROX VIEW adsdw.foo SET OPTIONS(ttl_seconds = 3600);
ALTER APPROX VIEW adsdw.foo SET OPTIONS(ttl_seconds = NULL);
ALTER APPROX VIEW foo SET OPTIONS(quota_accounting_owner = 'adsdw-etl@google.com', ttl_seconds = 3600);
ALTER MATERIALIZED VIEW foo SET OPTIONS(quota_accounting_owner = 'adsdw-etl@google.com', quota_accounting_owner =
    'adsdw-etl@google.com');
ALTER MODEL foo SET OPTIONS(a = 'a', b = 1), SET OPTIONS(c = 'c'), SET OPTIONS(a = 'a', b = 1);
