-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestResourceGroupOceanBaseOracle default#1
> exec
| create resource group x
< error
| 1064 42000 …near 'group x' at line 1
> exec
| create resource group x type=user
< error
| 1064 42000 …near 'group x type=user' at line 1
> exec
| CREATE RESOURCE GROUP rg1 TYPE = USER
< error
| 1064 42000 …near 'GROUP rg1 TYPE = USER' at line 1
> exec
| create resource group x type=system vcpu=0-3 thread_priority=0 enable
< error
| 1064 42000 …near 'group x type=system vcpu=0-3 thread_priority=0 enable' at line 1
> exec
| create resource group x type=user vcpu=0-3
< error
| 1064 42000 …near 'group x type=user vcpu=0-3' at line 1
> exec
| alter resource group x
< error
| 1064 42000 …near 'group x' at line 1
> exec
| alter resource group x vcpu=0-3
< error
| 1064 42000 …near 'group x vcpu=0-3' at line 1
> exec
| alter resource group x enable
< error
| 1064 42000 …near 'group x enable' at line 1
> exec
| alter resource group x force
< error
| 1064 42000 …near 'group x force' at line 1
> exec
| drop resource group x
< error
| 1064 42000 …near 'group x' at line 1
> exec
| drop resource group x force
< error
| 1064 42000 …near 'group x force' at line 1
> exec
| set resource group x
< error
| 1064 42000 …near 'resource group x' at line 1
> exec
| set resource group x for 1
< error
| 1064 42000 …near 'resource group x for 1' at line 1
> exec
| set resource group `default`
< error
| 1064 42000 …near 'resource group `default`' at line 1
> exec
| create resource group x ru_per_sec=1000
< error
| 1064 42000 …near 'group x ru_per_sec=1000' at line 1
> exec
| create resource group if not exists x ru_per_sec=1000
< error
| 1064 42000 …near 'group if not exists x ru_per_sec=1000' at line 1
> exec
| create resource group x burstable
< error
| 1064 42000 …near 'group x burstable' at line 1
> exec
| create resource group x priority=low
< error
| 1064 42000 …near 'group x priority=low' at line 1
> exec
| create resource group x ru_per_sec=1000 query_limit=(exec_elapsed '1s')
< error
| 1064 42000 …near 'group x ru_per_sec=1000 query_limit=(exec_elapsed '1s')' at line 1
> exec
| create resource group x ru_per_sec=1000 background=(task_types 'br')
< error
| 1064 42000 …near 'group x ru_per_sec=1000 background=(task_types 'br')' at line 1
> exec
| alter resource group x ru_per_sec=1000
< error
| 1064 42000 …near 'group x ru_per_sec=1000' at line 1
> exec
| drop resource group if exists x
< error
| 1064 42000 …near 'group if exists x' at line 1
> exec
| set resource group default
< error
| 1064 42000 …near 'resource group default' at line 1
> exec
| calibrate resource
< error
| 1064 42000 …near 'calibrate resource' at line 1
> exec
| calibrate resource workload tpcc
< error
| 1064 42000 …near 'calibrate resource workload tpcc' at line 1
> exec
| calibrate resource start_time now()
< error
| 1064 42000 …near 'calibrate resource start_time now()' at line 1
> exec
| query watch add sql digest 'a'
< error
| 1064 42000 …near 'query watch add sql digest 'a'' at line 1
> exec
| query watch add sql text similar to 'select 1'
< error
| 1064 42000 …near 'query watch add sql text similar to 'select 1'' at line 1
> exec
| query watch remove 1
< error
| 1064 42000 …near 'query watch remove 1' at line 1
> exec
| create user 'u'@'h' resource group rg
< error
| 1064 42000 …near 'resource group rg' at line 1
> exec
| alter user 'u'@'h' resource group rg
< error
| 1064 42000 …near 'resource group rg' at line 1
> exec
| show create resource group x
< error
| 1064 42000 …near 'resource group x' at line 1
