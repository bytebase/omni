# MySQL 8.0 system variable table — provenance

`sysvars.go` embeds the MySQL 8.0 session/global system variable name set used
by the SET-assignment validator to implement the `sp_head::find_variable`
fallback (see Bug 2 / Phase 6 of `docs/plans/2026-04-21-mysql-parser-validator-split.md`).

## Source: Strategy B (running container — performance_schema dump)

**Date**: 2026-04-21
**Image**: `mysql:8.0` (official Docker Hub image, pulled and started fresh for
the dump).
**Entry count**: 655 names.

### Procedure

```bash
# 1. Start a throwaway MySQL 8.0 container.
docker run -d --name mysql-sysvars-dump -e MYSQL_ROOT_PASSWORD=rootpw mysql:8.0

# 2. Wait until the server is ready (mysqladmin ping).
for i in $(seq 1 12); do
  docker exec mysql-sysvars-dump mysqladmin -uroot -prootpw ping 2>/dev/null \
    | grep -q alive && break
  sleep 5
done

# 3. Dump the union of global + session system variables, lowercased.
docker exec mysql-sysvars-dump mysql -uroot -prootpw -NBe "
  SELECT LOWER(VARIABLE_NAME) FROM performance_schema.global_variables
  UNION
  SELECT LOWER(VARIABLE_NAME) FROM performance_schema.session_variables
  ORDER BY 1;" > /tmp/sysvars_psched.txt

# 4. (Cleanup)
docker rm -f mysql-sysvars-dump
```

The resulting 655 names were dropped into `sysvars.go`'s map body verbatim
(alphabetized, lowercased, one entry per line for hand-readability).

## Why not Strategy A (`mysqld --verbose --help`)?

Strategy A (`mysqld --verbose --help`, Variables section) was attempted first
and produced 719 entries, but it has two problems:

1. **False positives**: it includes CLI-only startup options that are NOT
   session/global system variables — e.g. `console`, `daemonize`, `initialize`,
   `gdb`, `help`, `blackhole`, `archive`, NDB cluster variables that are absent
   in a non-NDB build. Treating these as sysvars would silently accept
   misspelled targets and bury real bugs.
2. **False negatives**: it misses session-only variables that are the most
   common routine-body SET targets — `sql_safe_updates` (the original Bug 2
   repro), `foreign_key_checks`, `unique_checks`, `sql_log_bin`, `time_zone`,
   the full `character_set_*` family, etc. These are exactly the variables the
   fallback needs to recognize.

`performance_schema.global_variables ∪ performance_schema.session_variables`
is the authoritative runtime set the server itself reports, and matches what
`SHOW VARIABLES` exposes.

## Refresh procedure

To refresh the list for a newer MySQL server version:

1. Rerun the Strategy B dump (above) against the target version's container.
2. Replace the map body of `sysvars.go` with the new entries (keep them
   alphabetized + lowercased).
3. Bump the version reference in the comment at the top of `sysvars.go`.

If Docker is unavailable, Strategy A (`mysqld --verbose --help` / Variables
section) can be used as a last-resort fallback but must be manually merged with
a curated list of session-only variables to avoid the false negatives noted
above.
