-- pin: default mcr.microsoft.com/mssql/server:2022-CU27-GDR1-ubuntu-22.04 sha256:cf8f0ab2ebe61feec2eed2629061a07d247f791c2a8ddda74a4ac9f06a291718 16.0.4305.5

== TestSyncDatabaseMissing default#1
> connect
< error
| 4063 1 11 1 Cannot open database "omni_no_such_database" that was requested by the login. Using the user default database "master" instead.
> connect
< error
| 4063 1 11 1 Cannot open database "omni_no_such_database" that was requested by the login. Using the user default database "master" instead.
