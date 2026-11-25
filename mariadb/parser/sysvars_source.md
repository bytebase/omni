# MariaDB system variables: provenance

`sysvars.go` lists the system variables of every release in `mariadb/testdata/pins.tsv` (10.7.8, 10.11.19, 11.4.13, 11.8.8), each with the plugins MariaDB packages for it installed: 1,145 names. A stored program's SET of a bare name that is no variable of the program reads as a system variable while MariaDB parses the program, and one the server lacks is 1193. The catalog and the review walk refuse a name no listed release or plugin has. A plugin outside the list can still own a name the list lacks.

Dumped 2026-11-25. For each pin, in a throwaway container started from `image@digest` with `--network host --skip-networking`:

1. `apt-get update` with the image's own MariaDB repository, then `apt-get install --no-install-recommends mariadb-plugin-<p>` for `connect cracklib-password-check gssapi-server hashicorp-key-management mroonga oqgraph rocksdb s3 spider columnstore`.
2. `INSTALL SONAME '<p>'` for every `.so` in `@@plugin_dir`, then `SELECT LOWER(VARIABLE_NAME) FROM information_schema.SYSTEM_VARIABLES`.
3. Again from a fresh container started with `--plugin-maturity=unknown --loose-file-key-management-filename=<a key file>`, for `file_key_management` and `handlersocket`, which step 2 cannot load.

The list is the union of the dumps. The client-side and MySQL-only libraries in `@@plugin_dir` load into no server. Two plugins are missing: `hashicorp_key_management` fails to start without a reachable Vault, and `auth_gssapi` without a Kerberos keytab.
