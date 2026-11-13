-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestResourceGroupOracle default#1
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| create resource group x
< error
| 1064 42000 …near 'resource group x' at line 1
> exec ^
: string
| create resource group x type=user
< error
| 1064 42000 …near 'resource group x type=user' at line 1
> exec ^
: string
| CREATE RESOURCE GROUP rg1 TYPE = USER
< error
| 1064 42000 …near 'RESOURCE GROUP rg1 TYPE = USER' at line 1
> exec ^
: string
| create resource group x type=system vcpu=0-3 thread_priority=0 enable
< error
| 1064 42000 …near 'resource group x type=system vcpu=0-3 thread_priority=0 enable' at line 1
> exec ^
: string
| create resource group x type=user vcpu=0-3
< error
| 1064 42000 …near 'resource group x type=user vcpu=0-3' at line 1
> exec ^
: string
| alter resource group x
< error
| 1064 42000 …near 'resource group x' at line 1
> exec ^
: string
| alter resource group x vcpu=0-3
< error
| 1064 42000 …near 'resource group x vcpu=0-3' at line 1
> exec ^
: string
| alter resource group x enable
< error
| 1064 42000 …near 'resource group x enable' at line 1
> exec ^
: string
| alter resource group x force
< error
| 1064 42000 …near 'resource group x force' at line 1
> exec ^
: string
| drop resource group x
< error
| 1064 42000 …near 'resource group x' at line 1
> exec ^
: string
| drop resource group x force
< error
| 1064 42000 …near 'resource group x force' at line 1
> exec ^
: string
| set resource group x
< error
| 1064 42000 …near 'group x' at line 1
> exec ^
: string
| set resource group x for 1
< error
| 1064 42000 …near 'group x for 1' at line 1
> exec ^
: string
| set resource group `default`
< error
| 1064 42000 …near 'group `default`' at line 1
> exec ^
: string
| create resource group x ru_per_sec=1000
< error
| 1064 42000 …near 'resource group x ru_per_sec=1000' at line 1
> exec ^
: string
| create resource group if not exists x ru_per_sec=1000
< error
| 1064 42000 …near 'resource group if not exists x ru_per_sec=1000' at line 1
> exec ^
: string
| create resource group x burstable
< error
| 1064 42000 …near 'resource group x burstable' at line 1
> exec ^
: string
| create resource group x priority=low
< error
| 1064 42000 …near 'resource group x priority=low' at line 1
> exec ^
: string
| create resource group x ru_per_sec=1000 query_limit=(exec_elapsed '1s')
< error
| 1064 42000 …near 'resource group x ru_per_sec=1000 query_limit=(exec_elapsed '1s')' at line 1
> exec ^
: string
| create resource group x ru_per_sec=1000 background=(task_types 'br')
< error
| 1064 42000 …near 'resource group x ru_per_sec=1000 background=(task_types 'br')' at line 1
> exec ^
: string
| alter resource group x ru_per_sec=1000
< error
| 1064 42000 …near 'resource group x ru_per_sec=1000' at line 1
> exec ^
: string
| drop resource group if exists x
< error
| 1064 42000 …near 'resource group if exists x' at line 1
> exec ^
: string
| set resource group default
< error
| 1064 42000 …near 'group default' at line 1
> exec ^
: string
| calibrate resource
< error
| 1064 42000 …near 'calibrate resource' at line 1
> exec ^
: string
| calibrate resource workload tpcc
< error
| 1064 42000 …near 'calibrate resource workload tpcc' at line 1
> exec ^
: string
| calibrate resource start_time now()
< error
| 1064 42000 …near 'calibrate resource start_time now()' at line 1
> exec ^
: string
| query watch add sql digest 'a'
< error
| 1064 42000 …near 'query watch add sql digest 'a'' at line 1
> exec ^
: string
| query watch add sql text similar to 'select 1'
< error
| 1064 42000 …near 'query watch add sql text similar to 'select 1'' at line 1
> exec ^
: string
| query watch remove 1
< error
| 1064 42000 …near 'query watch remove 1' at line 1
> exec ^
: string
| create user 'u'@'h' resource group rg
< error
| 1064 42000 …near 'resource group rg' at line 1
> exec ^
: string
| alter user 'u'@'h' resource group rg
< error
| 1064 42000 …near 'resource group rg' at line 1
> exec ^
: string
| show create resource group x
< error
| 1064 42000 …near 'resource group x' at line 1
