-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestResourceGroupOracle default#1
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| create resource group x type=user
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create resource group x type=user
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create resource group x type=user
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| create resource group x type user
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create resource group x type user
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create resource group x type user
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| create resource group x type=system
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create resource group x type=system
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create resource group x type=system
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| CREATE RESOURCE GROUP rg1 TYPE = USER
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| CREATE RESOURCE GROUP rg1 TYPE = USER
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE RESOURCE GROUP rg1 TYPE = USER
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| create resource group `default` type=user
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create resource group `default` type=user
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create resource group `default` type=user
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| create resource group x type=user vcpu=0-3
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create resource group x type=user vcpu=0-3
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create resource group x type=user vcpu=0-3
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| create resource group x type=user vcpu 0-3
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create resource group x type=user vcpu 0-3
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create resource group x type=user vcpu 0-3
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| create resource group x type=user vcpu=0-3,5
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create resource group x type=user vcpu=0-3,5
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create resource group x type=user vcpu=0-3,5
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| create resource group x type=user vcpu=0
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create resource group x type=user vcpu=0
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create resource group x type=user vcpu=0
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| create resource group x type=system vcpu=0-3 thread_priority=0 enable
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create resource group x type=system vcpu=0-3 thread_priority=0 enable
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create resource group x type=system vcpu=0-3 thread_priority=0 enable
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| create resource group x type=user vcpu=0-3 thread_priority=-10 disable
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create resource group x type=user vcpu=0-3 thread_priority=-10 disable
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create resource group x type=user vcpu=0-3 thread_priority=-10 disable
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| alter resource group x
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| alter resource group x
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| alter resource group x
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| alter resource group x vcpu=0-3
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| alter resource group x vcpu=0-3
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| alter resource group x vcpu=0-3
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| alter resource group x enable
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| alter resource group x enable
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| alter resource group x enable
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| alter resource group x disable
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| alter resource group x disable
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| alter resource group x disable
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| alter resource group x force
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| alter resource group x force
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| alter resource group x force
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| alter resource group x enable force
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| alter resource group x enable force
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| alter resource group x enable force
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| alter resource group x thread_priority=5 enable
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| alter resource group x thread_priority=5 enable
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| alter resource group x thread_priority=5 enable
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| drop resource group x
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| drop resource group x
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| drop resource group x
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| drop resource group x force
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| drop resource group x force
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| drop resource group x force
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| set resource group x
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN OR RESOURCE_GROUP_USER privilege(s) for this operation
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| set resource group x
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN OR RESOURCE_GROUP_USER privilege(s) for this operation
> exec ^
: string
| set resource group x
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| set resource group x for 1
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN OR RESOURCE_GROUP_USER privilege(s) for this operation
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| set resource group x for 1
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN OR RESOURCE_GROUP_USER privilege(s) for this operation
> exec ^
: string
| set resource group x for 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| set resource group x for 1,2
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN OR RESOURCE_GROUP_USER privilege(s) for this operation
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| set resource group x for 1,2
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN OR RESOURCE_GROUP_USER privilege(s) for this operation
> exec ^
: string
| set resource group x for 1,2
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| set resource group `default`
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN OR RESOURCE_GROUP_USER privilege(s) for this operation
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| set resource group `default`
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN OR RESOURCE_GROUP_USER privilege(s) for this operation
> exec ^
: string
| set resource group `default`
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| create resource group x
< error
| 1064 42000 …near '' at line 1
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| create resource group x vcpu=0-3
< error
| 1064 42000 …near 'vcpu=0-3' at line 1
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| create resource group x type=foo
< error
| 1064 42000 …near 'foo' at line 1
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| create resource group x type=1
< error
| 1064 42000 …near '1' at line 1
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| create resource group x type=
< error
| 1064 42000 …near '' at line 1
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| create resource group x vcpu=0-3 type=user
< error
| 1064 42000 …near 'vcpu=0-3 type=user' at line 1
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| create resource group x type=user type=user
< error
| 1064 42000 …near 'type=user' at line 1
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| create resource group x type=user vcpu=
< error
| 1064 42000 …near '' at line 1
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| create resource group x type=user vcpu
< error
| 1064 42000 …near '' at line 1
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| create resource group x type=user enable disable
< error
| 1064 42000 …near 'disable' at line 1
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| create resource group x type=user disable enable
< error
| 1064 42000 …near 'enable' at line 1
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| alter resource group x disable enable
< error
| 1064 42000 …near 'enable' at line 1
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| create resource group x type=user force
< error
| 1064 42000 …near 'force' at line 1
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| alter resource group x force enable
< error
| 1064 42000 …near 'enable' at line 1
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| alter resource group x type=user
< error
| 1064 42000 …near 'type=user' at line 1
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| create resource group x ru_per_sec=1000
< error
| 1064 42000 …near 'ru_per_sec=1000' at line 1
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| create resource group x type=user ru_per_sec=1000
< error
| 1064 42000 …near 'ru_per_sec=1000' at line 1
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| create resource group x burstable
< error
| 1064 42000 …near 'burstable' at line 1
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| create resource group x priority=low
< error
| 1064 42000 …near 'priority=low' at line 1
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| create resource group x type=user query_limit=(exec_elapsed '1s')
< error
| 1064 42000 …near 'query_limit=(exec_elapsed '1s')' at line 1
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| create resource group x type=user background=(task_types 'br')
< error
| 1064 42000 …near 'background=(task_types 'br')' at line 1
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| create resource group if not exists x type=user
< error
| 1064 42000 …near 'if not exists x type=user' at line 1
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| drop resource group if exists x
< error
| 1064 42000 …near 'if exists x' at line 1
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| set resource group default
< error
| 1064 42000 …near 'default' at line 1
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| calibrate resource
< error
| 1064 42000 …near 'calibrate resource' at line 1
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| calibrate resource workload tpcc
< error
| 1064 42000 …near 'calibrate resource workload tpcc' at line 1
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| query watch add sql digest 'a'
< error
| 1064 42000 …near 'query watch add sql digest 'a'' at line 1
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| query watch remove 1
< error
| 1064 42000 …near 'query watch remove 1' at line 1
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| create user 'u'@'h' resource group rg
< error
| 1064 42000 …near 'resource group rg' at line 1
> exec
| DROP RESOURCE GROUP x FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP rg1 FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| DROP RESOURCE GROUP `default` FORCE
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| alter user 'u'@'h' resource group rg
< error
| 1064 42000 …near 'resource group rg' at line 1
