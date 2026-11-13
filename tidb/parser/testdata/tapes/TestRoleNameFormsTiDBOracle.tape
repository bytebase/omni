-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestRoleNameFormsTiDBOracle default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| grant proxy@h to u
< ok 0 0
> exec ^
: string
| grant proxy@h, a to u
< ok 0 0
> exec ^
: string
| grant a, proxy@h to u
< ok 0 0
> exec ^
: string
| grant proxy@'h' to u
< ok 0 0
> exec ^
: string
| revoke proxy@h from u
< ok 0 0
> exec ^
: string
| create role proxy@h
< ok 0 0
> exec ^
: string
| grant proxy to u
< error
| 1064 42000 …line 1 column 14 near "to u" 
> exec ^
: string
| grant proxy, a to u
< error
| 1064 42000 …line 1 column 12 near ", a to u" 
> exec ^
: string
| grant proxy on *.* to u
< error
| 1064 42000 …line 1 column 16 near "*.* to u" 
> exec ^
: string
| grant proxy on a to u
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| grant proxy on a to u
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| grant proxy on a@h to u
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| grant proxy on a@h to u
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| grant `proxy` to u
< ok 0 0
> exec ^
: string
| create role 'transaction'
< ok 0 0
> exec ^
: string
| create role LOOP
< ok 0 0
> exec ^
: string
| create role Loop, GET
< ok 0 0
> exec ^
: string
| create role TRANSACTION
< error
| 1064 42000 …line 1 column 23 near "" 
> exec ^
: string
| create role `loop`
< ok 0 0
> exec ^
: string
| create role `select`
< ok 0 0
> exec ^
: string
| create role `transaction`
< ok 0 0
> exec ^
: string
| create role a, loop
< ok 0 0
> exec ^
: string
| create role a, transaction
< error
| 1064 42000 …line 1 column 26 near "" 
> exec ^
: string
| create role a@h, loop
< ok 0 0
> exec ^
: string
| create role count
< ok 0 0
> exec ^
: string
| create role count@h
< ok 0 0
> exec ^
: string
| create role if
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| create role if exists a
< error
| 1064 42000 …line 1 column 21 near "exists a" 
> exec ^
: string
| create role if not a
< error
| 1064 42000 …line 1 column 20 near "a" 
> exec ^
: string
| create role if not exists a
< ok 0 0
> exec ^
: string
| create role if not exists loop
< ok 0 0
> exec ^
: string
| create role if@h
< error
| 1064 42000 …line 1 column 16 near "@h" 
> exec ^
: string
| create role loop
< ok 0 0
> exec ^
: string
| create role loop, a
< ok 0 0
> exec ^
: string
| create role loop@'h'
< ok 0 0
> exec ^
: string
| create role loop@h
< ok 0 0
> exec ^
: string
| create role select
< error
| 1064 42000 …line 1 column 18 near "select" 
> exec ^
: string
| create role select@h
< error
| 1064 42000 …line 1 column 18 near "select@h" 
> exec ^
: string
| create role transaction@h
< ok 0 0
> exec ^
: string
| drop role LOOP
< ok 0 0
> exec ^
: string
| drop role a, loop
< ok 0 0
> exec ^
: string
| drop role if
< error
| 1064 42000 …line 1 column 12 near "" 
> exec ^
: string
| drop role if a
< error
| 1064 42000 …line 1 column 14 near "a" 
> exec ^
: string
| drop role if exists a
< ok 0 0
> exec ^
: string
| drop role if exists loop
< ok 0 0
> exec ^
: string
| drop role if not exists a
< error
| 1064 42000 …line 1 column 16 near "not exists a" 
> exec ^
: string
| drop role if@h
< error
| 1064 42000 …line 1 column 14 near "@h" 
> exec ^
: string
| drop role loop
< ok 0 0
> exec ^
: string
| drop role transaction
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| drop role transaction@h
< ok 0 0
> exec ^
: string
| grant @v to current_user
< error
| 1064 42000 …line 1 column 8 near "@v to current_user" 
> exec ^
: string
| revoke @v from u
< error
| 1064 42000 …line 1 column 9 near "@v from u" 
> exec ^
: string
| grant 'loop' to u
< ok 0 0
> exec ^
: string
| grant Create View on *.* to u
< ok 0 0
> exec ^
: string
| grant LOOP to u
< ok 0 0
> exec ^
: string
| grant Loop Get on *.* to u
< ok 0 0
> exec ^
: string
| grant PROXY on a to u
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| grant PROXY on a to u
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| grant SELECT on *.* to u
< ok 0 0
> exec ^
: string
| grant TRANSACTION to u
< error
| 1064 42000 …line 1 column 20 near "to u" 
> exec ^
: string
| grant `admin` on *.* to u
< ok 0 0
> exec ^
: string
| grant `loop` to u
< ok 0 0
> exec ^
: string
| grant `select` on *.* to u
< ok 0 0
> exec ^
: string
| grant `select`(a) on t to u
< error
| 1064 42000 …line 1 column 15 near "(a) on t to u" 
> exec ^
: string
| grant `transaction` to u
< ok 0 0
> exec ^
: string
| grant a loop on *.* to u
< ok 0 0
> exec ^
: string
| grant a loop to u
< ok 0 0
> exec ^
: string
| grant a(b) on t to u
< error
| 1064 42000 …line 1 column 8 near "(b) on t to u" 
> exec ^
: string
| grant a, all to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant a, loop to u
< ok 0 0
> exec ^
: string
| grant add on *.* to u
< error
| 1064 42000 …line 1 column 9 near "add on *.* to u" 
> exec ^
: string
| grant add(a) on t to u
< error
| 1064 42000 …line 1 column 9 near "add(a) on t to u" 
> exec ^
: string
| grant admin on *.* to u
< error
| 1064 42000 …line 1 column 14 near "on *.* to u" 
> exec ^
: string
| grant all on *.* to u
< ok 0 0
> exec ^
: string
| grant all on *.* to u with grant option
< ok 0 0
> exec ^
: string
| grant all on t to u
< ok 0 0
> exec ^
: string
| grant all privileges on *.* to u
< ok 0 0
> exec ^
: string
| grant all privileges on t to u
< ok 0 0
> exec ^
: string
| grant all privileges to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant all to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant all, a to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant alter on *.* to u
< ok 0 0
> exec ^
: string
| grant alter routine on *.* to u
< ok 0 0
> exec ^
: string
| grant alter routine to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant alter to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant backup_admin on *.* to u
< ok 0 0
> exec ^
: string
| grant comment on *.* to u1
< error
| 1064 42000 …line 1 column 16 near "on *.* to u1" 
> exec ^
: string
| grant comment, select on *.* to u1
< error
| 1064 42000 …line 1 column 14 near ", select on *.* to u1" 
> exec ^
: string
| grant concat to u
< ok 0 0
> exec ^
: string
| grant config on *.* to u
< ok 0 0
> exec ^
: string
| grant config to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant count on *.* to u
< ok 0 0
> exec ^
: string
| grant count to u
< ok 0 0
> exec ^
: string
| grant count(a) on t to u
< error
| 1064 42000 …line 1 column 11 near "count(a) on t to u" 
> exec ^
: string
| grant create on *.* to u
< ok 0 0
> exec ^
: string
| grant create role on *.* to u
< ok 0 0
> exec ^
: string
| grant create role to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant create routine on *.* to u
< ok 0 0
> exec ^
: string
| grant create routine to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant create tablespace on *.* to u
< ok 0 0
> exec ^
: string
| grant create tablespace to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant create temporary on *.* to u
< error
| 1064 42000 …line 1 column 25 near "on *.* to u" 
> exec ^
: string
| grant create temporary tables on *.* to u
< ok 0 0
> exec ^
: string
| grant create temporary tables to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant create to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant create user on *.* to u
< ok 0 0
> exec ^
: string
| grant create user to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant create view on *.* to u
< ok 0 0
> exec ^
: string
| grant create view to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant delete on *.* to u
< ok 0 0
> exec ^
: string
| grant delete to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant drop on *.* to u
< ok 0 0
> exec ^
: string
| grant drop role on *.* to u
< ok 0 0
> exec ^
: string
| grant drop role to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant drop to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant event on *.* to u
< ok 0 0
> exec ^
: string
| grant event to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant execute on *.* to u
< ok 0 0
> exec ^
: string
| grant execute to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant file on *.* to u
< ok 0 0
> exec ^
: string
| grant file to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant grant on *.* to u
< error
| 1064 42000 …line 1 column 14 near "on *.* to u" 
> exec ^
: string
| grant grant option on *.* to u
< ok 0 0
> exec ^
: string
| grant grant option to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant grant to u
< error
| 1064 42000 …line 1 column 14 near "to u" 
> exec ^
: string
| grant grant, select to u
< error
| 1064 42000 …line 1 column 12 near ", select to u" 
> exec ^
: string
| grant index on *.* to u
< ok 0 0
> exec ^
: string
| grant index to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant insert on *.* to u
< ok 0 0
> exec ^
: string
| grant insert to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant insert(a), update(b) on t to u
< ok 0 0
> exec ^
: string
| grant invoke `lambda` on *.* to u
< ok 0 0
> exec ^
: string
| grant invoke lambda on *.* to u
< ok 0 0
> exec ^
: string
| grant invoke loop on *.* to u
< ok 0 0
> exec ^
: string
| grant load from s3 on *.* to u
< ok 0 0
> exec ^
: string
| grant load on *.* to u
< error
| 1064 42000 …line 1 column 13 near "on *.* to u" 
> exec ^
: string
| grant lock on *.* to u
< error
| 1064 42000 …line 1 column 13 near "on *.* to u" 
> exec ^
: string
| grant lock tables on *.* to u
< ok 0 0
> exec ^
: string
| grant lock tables to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant loop a on *.* to u
< ok 0 0
> exec ^
: string
| grant loop a to u
< ok 0 0
> exec ^
: string
| grant loop get on *.* to u
< ok 0 0
> exec ^
: string
| grant loop on *.* to u
< ok 0 0
> exec ^
: string
| grant loop to loop
< ok 0 0
> exec ^
: string
| grant loop to transaction
< ok 0 0
> exec ^
: string
| grant loop to u
< ok 0 0
> exec ^
: string
| grant loop transaction on *.* to u
< error
| 1064 42000 …line 1 column 22 near "transaction on *.* to u" 
> exec ^
: string
| grant loop(a) on t to u
< error
| 1064 42000 …line 1 column 11 near "(a) on t to u" 
> exec ^
: string
| grant loop, a to u
< ok 0 0
> exec ^
: string
| grant loop, select on *.* to u
< ok 0 0
> exec ^
: string
| grant loop, transaction to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| grant loop@'h' to u
< ok 0 0
> exec ^
: string
| grant loop@h to u
< ok 0 0
> exec ^
: string
| grant process on *.* to u
< ok 0 0
> exec ^
: string
| grant process to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant proxy on 'a'@'b' to u
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| grant proxy on 'a'@'b' to u
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| grant proxy on *.* to u
< error
| 1064 42000 …line 1 column 16 near "*.* to u" 
> exec ^
: string
| grant proxy on a to u
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| grant proxy on a to u
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| grant proxy to u
< error
| 1064 42000 …line 1 column 14 near "to u" 
> exec ^
: string
| grant proxy(a) on t to u
< error
| 1064 42000 …line 1 column 12 near "(a) on t to u" 
> exec ^
: string
| grant proxy, select on *.* to u
< error
| 1064 42000 …line 1 column 12 near ", select on *.* to u" 
> exec ^
: string
| grant references on *.* to u
< ok 0 0
> exec ^
: string
| grant references to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant reload on *.* to u
< ok 0 0
> exec ^
: string
| grant reload to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant replication client on *.* to u
< ok 0 0
> exec ^
: string
| grant replication client to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant replication client, replication slave on *.* to u
< ok 0 0
> exec ^
: string
| grant replication on *.* to u
< error
| 1064 42000 …line 1 column 20 near "on *.* to u" 
> exec ^
: string
| grant replication slave on *.* to u
< ok 0 0
> exec ^
: string
| grant replication slave to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant restricted_tables_admin on *.* to u
< ok 0 0
> exec ^
: string
| grant s3 on *.* to u1
< error
| 1064 42000 …line 1 column 11 near "on *.* to u1" 
> exec ^
: string
| grant select (a, b), insert(c) on t to u
< ok 0 0
> exec ^
: string
| grant select into s3 on *.* to u
< ok 0 0
> exec ^
: string
| grant select on *.* to u
< ok 0 0
> exec ^
: string
| grant select to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant select(a) on t to u
< ok 0 0
> exec ^
: string
| grant select(a) to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant select, grant on *.* to u
< error
| 1064 42000 …line 1 column 22 near "on *.* to u" 
> exec ^
: string
| grant select, loop on *.* to u
< ok 0 0
> exec ^
: string
| grant select, loop to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant select, proxy on *.* to u
< error
| 1064 42000 …line 1 column 22 near "on *.* to u" 
> exec ^
: string
| grant show databases on *.* to u
< ok 0 0
> exec ^
: string
| grant show databases to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant show on *.* to u
< error
| 1064 42000 …line 1 column 13 near "on *.* to u" 
> exec ^
: string
| grant show view on *.* to u
< ok 0 0
> exec ^
: string
| grant show view to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant shutdown on *.* to u
< ok 0 0
> exec ^
: string
| grant shutdown to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant super on *.* to u
< ok 0 0
> exec ^
: string
| grant super to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant system_variables_admin on *.* to u
< ok 0 0
> exec ^
: string
| grant transaction loop on *.* to u
< error
| 1064 42000 …line 1 column 22 near "loop on *.* to u" 
> exec ^
: string
| grant transaction, loop to u
< error
| 1064 42000 …line 1 column 18 near ", loop to u" 
> exec ^
: string
| grant transaction@h to u
< ok 0 0
> exec ^
: string
| grant trigger on *.* to u
< ok 0 0
> exec ^
: string
| grant trigger to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant update on *.* to u
< ok 0 0
> exec ^
: string
| grant update to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant update(a) on t to u
< ok 0 0
> exec ^
: string
| grant usage on *.* to u
< ok 0 0
> exec ^
: string
| grant usage to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| revoke LOOP from u
< ok 0 0
> exec ^
: string
| revoke a, all from u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| revoke a, transaction from u
< error
| 1064 42000 …line 1 column 26 near "from u" 
> exec ^
: string
| revoke all from u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| revoke all on *.* from u
< ok 0 0
> exec ^
: string
| revoke all privileges from u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| revoke all privileges on *.* from u
< ok 0 0
> exec ^
: string
| revoke all privileges, grant option from u
< ok 0 0
> exec ^
: string
| revoke all, grant option from u
< ok 0 0
> exec ^
: string
| revoke all, grant option from u, v
< ok 0 0
> exec ^
: string
| revoke loop from u
< ok 0 0
> exec ^
: string
| revoke loop on *.* from u
< ok 0 0
> exec ^
: string
| revoke loop, a from u
< ok 0 0
> exec ^
: string
| revoke loop@h from u
< ok 0 0
> exec ^
: string
| revoke proxy from u
< error
| 1064 42000 …line 1 column 17 near "from u" 
> exec ^
: string
| revoke proxy on 'a'@'b' from u
< error
| 1064 42000 …line 1 column 15 near "on 'a'@'b' from u" 
> exec ^
: string
| revoke proxy on *.* from u
< error
| 1064 42000 …line 1 column 15 near "on *.* from u" 
> exec ^
: string
| set default role a, loop to u
< ok 0 0
> exec ^
: string
| set default role all to u
< ok 0 0
> exec ^
: string
| set default role loop to loop
< ok 0 0
> exec ^
: string
| set default role loop to u
< ok 0 0
> exec ^
: string
| set default role none to u
< ok 0 0
> exec ^
: string
| set default role transaction to u
< error
| 1064 42000 …line 1 column 31 near "to u" 
> exec ^
: string
| set default role transaction@h to u
< ok 0 0
> exec ^
: string
| set role LOOP
< ok 0 0
> exec ^
: string
| set role a, loop
< ok 0 0
> exec ^
: string
| set role all
< ok 0 0
> exec ^
: string
| set role all except a, loop
< ok 0 0
> exec ^
: string
| set role all except loop
< ok 0 0
> exec ^
: string
| set role all except transaction
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| set role default
< ok 0 0
> exec ^
: string
| set role defined
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role cast
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| set role loop
< ok 0 0
> exec ^
: string
| set role loop, a
< ok 0 0
> exec ^
: string
| set role none
< ok 0 0
> exec ^
: string
| set role transaction@h
< ok 0 0
> exec ^
: string
| show grants
< ok 0 0
> exec ^
: string
| show grants for u
< ok 0 0
> exec ^
: string
| show grants for u using 'a'
< ok 0 0
> exec ^
: string
| show grants for u using a
< ok 0 0
> exec ^
: string
| show grants for u using a, loop
< ok 0 0
> exec ^
: string
| show grants for u using a@h
< ok 0 0
> exec ^
: string
| show grants for u using loop
< ok 0 0
> exec ^
: string
| show grants for u using select
< error
| 1064 42000 …line 1 column 30 near "select" 
> exec ^
: string
| show grants for u using transaction
< error
| 1064 42000 …line 1 column 35 near "" 
> exec ^
: string
| show grants for u using transaction@h
< ok 0 0
