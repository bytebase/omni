-- pin: default cockroachdb/cockroach:v25.2.3 sha256:4a412731d1cfc5d6dea7431b7ce383b502966c9d589163d63db9e1ff7b51495a v25.2.3

== TestGetCurrentDatabaseOwner default#1
> query
|
| 		SELECT
| 			u.rolname
| 		FROM
| 			pg_roles AS u JOIN pg_database AS d ON (d.datdba = u.oid)
| 		WHERE
| 			d.datname = current_database();
| 		
< result
< column "rolname" "NAME"
< row
: string
| root
