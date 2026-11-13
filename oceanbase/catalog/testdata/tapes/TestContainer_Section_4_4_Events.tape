-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestContainer_Section_4_4_Events default#1
> exec
| CREATE EVENT ev_rename_old ON SCHEDULE EVERY 1 HOUR DO SELECT 1
< ok 0 0
> exec
| ALTER EVENT ev_rename_old RENAME TO ev_rename_new
< error
| 1235 0A000 modify EVENT NAME that has already started not supported
> exec
| DROP EVENT IF EXISTS ev_drop_test
< ok 0 0
> exec
| CREATE EVENT ev_drop_test ON SCHEDULE EVERY 1 HOUR DO SELECT 1
< ok 0 0
> exec
| DROP EVENT ev_drop_test
< ok 0 0
> exec
| DROP EVENT IF EXISTS ev_nonexistent
< ok 0 0
> exec
| DROP EVENT IF EXISTS ev_no_exist
< ok 0 0
> exec
| DROP EVENT ev_no_exist
< error
| 1539 HY000 Unknown event
> exec
| CREATE EVENT ev_show_test ON SCHEDULE EVERY 1 HOUR COMMENT 'test event' DO SELECT 1
< ok 0 0
> exec
| SHOW CREATE EVENT ev_show_test
< error
| 1064 42000 …near 'EVENT ev_show_test' at line 1
