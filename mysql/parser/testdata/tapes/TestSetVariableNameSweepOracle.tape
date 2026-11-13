-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestSetVariableNameSweepOracle default#1
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| set absent = 1
< error
| 1193 HY000 Unknown system variable 'absent'
> exec ^
: string
| set accessible = 1
< error
| 1064 42000 …near 'accessible = 1' at line 1
> exec ^
: string
| set account = 1
< error
| 1193 HY000 Unknown system variable 'account'
> exec ^
: string
| set action = 1
< error
| 1193 HY000 Unknown system variable 'action'
> exec ^
: string
| set active = 1
< error
| 1193 HY000 Unknown system variable 'active'
> exec ^
: string
| set add = 1
< error
| 1064 42000 …near 'add = 1' at line 1
> exec ^
: string
| set adddate = 1
< error
| 1193 HY000 Unknown system variable 'adddate'
> exec ^
: string
| set admin = 1
< error
| 1193 HY000 Unknown system variable 'admin'
> exec ^
: string
| set after = 1
< error
| 1193 HY000 Unknown system variable 'after'
> exec ^
: string
| set against = 1
< error
| 1193 HY000 Unknown system variable 'against'
> exec ^
: string
| set aggregate = 1
< error
| 1193 HY000 Unknown system variable 'aggregate'
> exec ^
: string
| set algorithm = 1
< error
| 1193 HY000 Unknown system variable 'algorithm'
> exec ^
: string
| set all = 1
< error
| 1064 42000 …near 'all = 1' at line 1
> exec ^
: string
| set allow_missing_files = 1
< error
| 1193 HY000 Unknown system variable 'allow_missing_files'
> exec ^
: string
| set alter = 1
< error
| 1064 42000 …near 'alter = 1' at line 1
> exec ^
: string
| set always = 1
< error
| 1193 HY000 Unknown system variable 'always'
> exec ^
: string
| set analyze = 1
< error
| 1064 42000 …near 'analyze = 1' at line 1
> exec ^
: string
| set and = 1
< error
| 1064 42000 …near 'and = 1' at line 1
> exec ^
: string
| set any = 1
< error
| 1193 HY000 Unknown system variable 'any'
> exec ^
: string
| set array = 1
< error
| 1193 HY000 Unknown system variable 'array'
> exec ^
: string
| set as = 1
< error
| 1064 42000 …near 'as = 1' at line 1
> exec ^
: string
| set asc = 1
< error
| 1064 42000 …near 'asc = 1' at line 1
> exec ^
: string
| set ascii = 1
< error
| 1193 HY000 Unknown system variable 'ascii'
> exec ^
: string
| set asensitive = 1
< error
| 1064 42000 …near 'asensitive = 1' at line 1
> exec ^
: string
| set assign_gtids_to_anonymous_transactions = 1
< error
| 1193 HY000 Unknown system variable 'assign_gtids_to_anonymous_transactions'
> exec ^
: string
| set at = 1
< error
| 1193 HY000 Unknown system variable 'at'
> exec ^
: string
| set attribute = 1
< error
| 1193 HY000 Unknown system variable 'attribute'
> exec ^
: string
| set authentication = 1
< error
| 1193 HY000 Unknown system variable 'authentication'
> exec ^
: string
| set auto = 1
< error
| 1193 HY000 Unknown system variable 'auto'
> exec ^
: string
| set auto_increment = 1
< error
| 1193 HY000 Unknown system variable 'auto_increment'
> exec ^
: string
| set auto_refresh = 1
< error
| 1193 HY000 Unknown system variable 'auto_refresh'
> exec ^
: string
| set auto_refresh_source = 1
< error
| 1193 HY000 Unknown system variable 'auto_refresh_source'
> exec ^
: string
| set autoextend_size = 1
< error
| 1193 HY000 Unknown system variable 'autoextend_size'
> exec ^
: string
| set avg = 1
< error
| 1193 HY000 Unknown system variable 'avg'
> exec ^
: string
| set avg_row_length = 1
< error
| 1193 HY000 Unknown system variable 'avg_row_length'
> exec ^
: string
| set backup = 1
< error
| 1193 HY000 Unknown system variable 'backup'
> exec ^
: string
| set before = 1
< error
| 1064 42000 …near 'before = 1' at line 1
> exec ^
: string
| set begin = 1
< error
| 1193 HY000 Unknown system variable 'begin'
> exec ^
: string
| set bernoulli = 1
< error
| 1193 HY000 Unknown system variable 'bernoulli'
> exec ^
: string
| set between = 1
< error
| 1064 42000 …near 'between = 1' at line 1
> exec ^
: string
| set bigint = 1
< error
| 1064 42000 …near 'bigint = 1' at line 1
> exec ^
: string
| set binary = 1
< error
| 1064 42000 …near 'binary = 1' at line 1
> exec ^
: string
| set binlog = 1
< error
| 1193 HY000 Unknown system variable 'binlog'
> exec ^
: string
| set bit = 1
< error
| 1193 HY000 Unknown system variable 'bit'
> exec ^
: string
| set bit_and = 1
< error
| 1193 HY000 Unknown system variable 'bit_and'
> exec ^
: string
| set bit_or = 1
< error
| 1193 HY000 Unknown system variable 'bit_or'
> exec ^
: string
| set bit_xor = 1
< error
| 1193 HY000 Unknown system variable 'bit_xor'
> exec ^
: string
| set blob = 1
< error
| 1064 42000 …near 'blob = 1' at line 1
> exec ^
: string
| set block = 1
< error
| 1193 HY000 Unknown system variable 'block'
> exec ^
: string
| set bool = 1
< error
| 1193 HY000 Unknown system variable 'bool'
> exec ^
: string
| set boolean = 1
< error
| 1193 HY000 Unknown system variable 'boolean'
> exec ^
: string
| set both = 1
< error
| 1064 42000 …near 'both = 1' at line 1
> exec ^
: string
| set btree = 1
< error
| 1193 HY000 Unknown system variable 'btree'
> exec ^
: string
| set buckets = 1
< error
| 1193 HY000 Unknown system variable 'buckets'
> exec ^
: string
| set bulk = 1
< error
| 1193 HY000 Unknown system variable 'bulk'
> exec ^
: string
| set by = 1
< error
| 1064 42000 …near 'by = 1' at line 1
> exec ^
: string
| set byte = 1
< error
| 1193 HY000 Unknown system variable 'byte'
> exec ^
: string
| set cache = 1
< error
| 1193 HY000 Unknown system variable 'cache'
> exec ^
: string
| set call = 1
< error
| 1064 42000 …near 'call = 1' at line 1
> exec ^
: string
| set cascade = 1
< error
| 1064 42000 …near 'cascade = 1' at line 1
> exec ^
: string
| set cascaded = 1
< error
| 1193 HY000 Unknown system variable 'cascaded'
> exec ^
: string
| set case = 1
< error
| 1064 42000 …near 'case = 1' at line 1
> exec ^
: string
| set cast = 1
< error
| 1193 HY000 Unknown system variable 'cast'
> exec ^
: string
| set catalog_name = 1
< error
| 1193 HY000 Unknown system variable 'catalog_name'
> exec ^
: string
| set chain = 1
< error
| 1193 HY000 Unknown system variable 'chain'
> exec ^
: string
| set challenge_response = 1
< error
| 1193 HY000 Unknown system variable 'challenge_response'
> exec ^
: string
| set change = 1
< error
| 1064 42000 …near 'change = 1' at line 1
> exec ^
: string
| set changed = 1
< error
| 1193 HY000 Unknown system variable 'changed'
> exec ^
: string
| set channel = 1
< error
| 1193 HY000 Unknown system variable 'channel'
> exec ^
: string
| set char = 1
< error
| 1064 42000 …near '= 1' at line 1
> exec ^
: string
| set character = 1
< error
| 1064 42000 …near '= 1' at line 1
> exec ^
: string
| set charset = 1
< error
| 1193 HY000 Unknown system variable 'charset'
> exec ^
: string
| set check = 1
< error
| 1064 42000 …near 'check = 1' at line 1
> exec ^
: string
| set checksum = 1
< error
| 1193 HY000 Unknown system variable 'checksum'
> exec ^
: string
| set cipher = 1
< error
| 1193 HY000 Unknown system variable 'cipher'
> exec ^
: string
| set class_origin = 1
< error
| 1193 HY000 Unknown system variable 'class_origin'
> exec ^
: string
| set client = 1
< error
| 1193 HY000 Unknown system variable 'client'
> exec ^
: string
| set clone = 1
< error
| 1193 HY000 Unknown system variable 'clone'
> exec ^
: string
| set close = 1
< error
| 1193 HY000 Unknown system variable 'close'
> exec ^
: string
| set coalesce = 1
< error
| 1193 HY000 Unknown system variable 'coalesce'
> exec ^
: string
| set code = 1
< error
| 1193 HY000 Unknown system variable 'code'
> exec ^
: string
| set collate = 1
< error
| 1064 42000 …near 'collate = 1' at line 1
> exec ^
: string
| set collation = 1
< error
| 1193 HY000 Unknown system variable 'collation'
> exec ^
: string
| set column = 1
< error
| 1064 42000 …near 'column = 1' at line 1
> exec ^
: string
| set column_format = 1
< error
| 1193 HY000 Unknown system variable 'column_format'
> exec ^
: string
| set column_name = 1
< error
| 1193 HY000 Unknown system variable 'column_name'
> exec ^
: string
| set columns = 1
< error
| 1193 HY000 Unknown system variable 'columns'
> exec ^
: string
| set comment = 1
< error
| 1193 HY000 Unknown system variable 'comment'
> exec ^
: string
| set commit = 1
< error
| 1193 HY000 Unknown system variable 'commit'
> exec ^
: string
| set committed = 1
< error
| 1193 HY000 Unknown system variable 'committed'
> exec ^
: string
| set compact = 1
< error
| 1193 HY000 Unknown system variable 'compact'
> exec ^
: string
| set completion = 1
< error
| 1193 HY000 Unknown system variable 'completion'
> exec ^
: string
| set component = 1
< error
| 1193 HY000 Unknown system variable 'component'
> exec ^
: string
| set compressed = 1
< error
| 1193 HY000 Unknown system variable 'compressed'
> exec ^
: string
| set compression = 1
< error
| 1193 HY000 Unknown system variable 'compression'
> exec ^
: string
| set concat = 1
< error
| 1193 HY000 Unknown system variable 'concat'
> exec ^
: string
| set concurrent = 1
< error
| 1193 HY000 Unknown system variable 'concurrent'
> exec ^
: string
| set condition = 1
< error
| 1064 42000 …near 'condition = 1' at line 1
> exec ^
: string
| set connection = 1
< error
| 1193 HY000 Unknown system variable 'connection'
> exec ^
: string
| set consistent = 1
< error
| 1193 HY000 Unknown system variable 'consistent'
> exec ^
: string
| set constraint = 1
< error
| 1064 42000 …near 'constraint = 1' at line 1
> exec ^
: string
| set constraint_catalog = 1
< error
| 1193 HY000 Unknown system variable 'constraint_catalog'
> exec ^
: string
| set constraint_name = 1
< error
| 1193 HY000 Unknown system variable 'constraint_name'
> exec ^
: string
| set constraint_schema = 1
< error
| 1193 HY000 Unknown system variable 'constraint_schema'
> exec ^
: string
| set contains = 1
< error
| 1193 HY000 Unknown system variable 'contains'
> exec ^
: string
| set context = 1
< error
| 1193 HY000 Unknown system variable 'context'
> exec ^
: string
| set continue = 1
< error
| 1064 42000 …near 'continue = 1' at line 1
> exec ^
: string
| set convert = 1
< error
| 1064 42000 …near 'convert = 1' at line 1
> exec ^
: string
| set count = 1
< error
| 1193 HY000 Unknown system variable 'count'
> exec ^
: string
| set cpu = 1
< error
| 1193 HY000 Unknown system variable 'cpu'
> exec ^
: string
| set create = 1
< error
| 1064 42000 …near 'create = 1' at line 1
> exec ^
: string
| set cross = 1
< error
| 1064 42000 …near 'cross = 1' at line 1
> exec ^
: string
| set csv = 1
< error
| 1193 HY000 Unknown system variable 'csv'
> exec ^
: string
| set cube = 1
< error
| 1064 42000 …near 'cube = 1' at line 1
> exec ^
: string
| set cume_dist = 1
< error
| 1064 42000 …near 'cume_dist = 1' at line 1
> exec ^
: string
| set curdate = 1
< error
| 1193 HY000 Unknown system variable 'curdate'
> exec ^
: string
| set current = 1
< error
| 1193 HY000 Unknown system variable 'current'
> exec ^
: string
| set current_date = 1
< error
| 1064 42000 …near 'current_date = 1' at line 1
> exec ^
: string
| set current_time = 1
< error
| 1064 42000 …near 'current_time = 1' at line 1
> exec ^
: string
| set current_timestamp = 1
< error
| 1064 42000 …near 'current_timestamp = 1' at line 1
> exec ^
: string
| set current_user = 1
< error
| 1064 42000 …near 'current_user = 1' at line 1
> exec ^
: string
| set cursor = 1
< error
| 1064 42000 …near 'cursor = 1' at line 1
> exec ^
: string
| set cursor_name = 1
< error
| 1193 HY000 Unknown system variable 'cursor_name'
> exec ^
: string
| set curtime = 1
< error
| 1193 HY000 Unknown system variable 'curtime'
> exec ^
: string
| set data = 1
< error
| 1193 HY000 Unknown system variable 'data'
> exec ^
: string
| set database = 1
< error
| 1064 42000 …near 'database = 1' at line 1
> exec ^
: string
| set databases = 1
< error
| 1064 42000 …near 'databases = 1' at line 1
> exec ^
: string
| set datafile = 1
< error
| 1193 HY000 Unknown system variable 'datafile'
> exec ^
: string
| set date = 1
< error
| 1193 HY000 Unknown system variable 'date'
> exec ^
: string
| set date_add = 1
< error
| 1193 HY000 Unknown system variable 'date_add'
> exec ^
: string
| set date_sub = 1
< error
| 1193 HY000 Unknown system variable 'date_sub'
> exec ^
: string
| set datetime = 1
< error
| 1193 HY000 Unknown system variable 'datetime'
> exec ^
: string
| set day = 1
< error
| 1193 HY000 Unknown system variable 'day'
> exec ^
: string
| set day_hour = 1
< error
| 1064 42000 …near 'day_hour = 1' at line 1
> exec ^
: string
| set day_microsecond = 1
< error
| 1064 42000 …near 'day_microsecond = 1' at line 1
> exec ^
: string
| set day_minute = 1
< error
| 1064 42000 …near 'day_minute = 1' at line 1
> exec ^
: string
| set day_second = 1
< error
| 1064 42000 …near 'day_second = 1' at line 1
> exec ^
: string
| set deallocate = 1
< error
| 1193 HY000 Unknown system variable 'deallocate'
> exec ^
: string
| set dec = 1
< error
| 1064 42000 …near 'dec = 1' at line 1
> exec ^
: string
| set decimal = 1
< error
| 1064 42000 …near 'decimal = 1' at line 1
> exec ^
: string
| set declare = 1
< error
| 1064 42000 …near 'declare = 1' at line 1
> exec ^
: string
| set default = 1
< error
| 1064 42000 …near '= 1' at line 1
> exec ^
: string
| set default_auth = 1
< error
| 1193 HY000 Unknown system variable 'default_auth'
> exec ^
: string
| set definer = 1
< error
| 1193 HY000 Unknown system variable 'definer'
> exec ^
: string
| set definition = 1
< error
| 1193 HY000 Unknown system variable 'definition'
> exec ^
: string
| set delay_key_write = 1
< error
| 1229 HY000 Variable 'delay_key_write' is a GLOBAL variable and should be set with SET GLOBAL
> exec ^
: string
| set delay_key_write = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set delayed = 1
< error
| 1064 42000 …near 'delayed = 1' at line 1
> exec ^
: string
| set delete = 1
< error
| 1064 42000 …near 'delete = 1' at line 1
> exec ^
: string
| set dense_rank = 1
< error
| 1064 42000 …near 'dense_rank = 1' at line 1
> exec ^
: string
| set desc = 1
< error
| 1064 42000 …near 'desc = 1' at line 1
> exec ^
: string
| set describe = 1
< error
| 1064 42000 …near 'describe = 1' at line 1
> exec ^
: string
| set description = 1
< error
| 1193 HY000 Unknown system variable 'description'
> exec ^
: string
| set deterministic = 1
< error
| 1064 42000 …near 'deterministic = 1' at line 1
> exec ^
: string
| set diagnostics = 1
< error
| 1193 HY000 Unknown system variable 'diagnostics'
> exec ^
: string
| set directory = 1
< error
| 1193 HY000 Unknown system variable 'directory'
> exec ^
: string
| set disable = 1
< error
| 1193 HY000 Unknown system variable 'disable'
> exec ^
: string
| set discard = 1
< error
| 1193 HY000 Unknown system variable 'discard'
> exec ^
: string
| set disk = 1
< error
| 1193 HY000 Unknown system variable 'disk'
> exec ^
: string
| set distinct = 1
< error
| 1064 42000 …near 'distinct = 1' at line 1
> exec ^
: string
| set distinctrow = 1
< error
| 1064 42000 …near 'distinctrow = 1' at line 1
> exec ^
: string
| set div = 1
< error
| 1064 42000 …near 'div = 1' at line 1
> exec ^
: string
| set do = 1
< error
| 1193 HY000 Unknown system variable 'do'
> exec ^
: string
| set double = 1
< error
| 1064 42000 …near 'double = 1' at line 1
> exec ^
: string
| set drop = 1
< error
| 1064 42000 …near 'drop = 1' at line 1
> exec ^
: string
| set dual = 1
< error
| 1064 42000 …near 'dual = 1' at line 1
> exec ^
: string
| set duality = 1
< error
| 1193 HY000 Unknown system variable 'duality'
> exec ^
: string
| set dumpfile = 1
< error
| 1193 HY000 Unknown system variable 'dumpfile'
> exec ^
: string
| set duplicate = 1
< error
| 1193 HY000 Unknown system variable 'duplicate'
> exec ^
: string
| set dynamic = 1
< error
| 1193 HY000 Unknown system variable 'dynamic'
> exec ^
: string
| set each = 1
< error
| 1064 42000 …near 'each = 1' at line 1
> exec ^
: string
| set else = 1
< error
| 1064 42000 …near 'else = 1' at line 1
> exec ^
: string
| set elseif = 1
< error
| 1064 42000 …near 'elseif = 1' at line 1
> exec ^
: string
| set empty = 1
< error
| 1064 42000 …near 'empty = 1' at line 1
> exec ^
: string
| set enable = 1
< error
| 1193 HY000 Unknown system variable 'enable'
> exec ^
: string
| set enclosed = 1
< error
| 1064 42000 …near 'enclosed = 1' at line 1
> exec ^
: string
| set encryption = 1
< error
| 1193 HY000 Unknown system variable 'encryption'
> exec ^
: string
| set end = 1
< error
| 1193 HY000 Unknown system variable 'end'
> exec ^
: string
| set ends = 1
< error
| 1193 HY000 Unknown system variable 'ends'
> exec ^
: string
| set enforced = 1
< error
| 1193 HY000 Unknown system variable 'enforced'
> exec ^
: string
| set engine = 1
< error
| 1193 HY000 Unknown system variable 'engine'
> exec ^
: string
| set engine_attribute = 1
< error
| 1193 HY000 Unknown system variable 'engine_attribute'
> exec ^
: string
| set engines = 1
< error
| 1193 HY000 Unknown system variable 'engines'
> exec ^
: string
| set enum = 1
< error
| 1193 HY000 Unknown system variable 'enum'
> exec ^
: string
| set error = 1
< error
| 1193 HY000 Unknown system variable 'error'
> exec ^
: string
| set errors = 1
< error
| 1193 HY000 Unknown system variable 'errors'
> exec ^
: string
| set escape = 1
< error
| 1193 HY000 Unknown system variable 'escape'
> exec ^
: string
| set escaped = 1
< error
| 1064 42000 …near 'escaped = 1' at line 1
> exec ^
: string
| set event = 1
< error
| 1193 HY000 Unknown system variable 'event'
> exec ^
: string
| set events = 1
< error
| 1193 HY000 Unknown system variable 'events'
> exec ^
: string
| set every = 1
< error
| 1193 HY000 Unknown system variable 'every'
> exec ^
: string
| set except = 1
< error
| 1064 42000 …near 'except = 1' at line 1
> exec ^
: string
| set exchange = 1
< error
| 1193 HY000 Unknown system variable 'exchange'
> exec ^
: string
| set exclude = 1
< error
| 1193 HY000 Unknown system variable 'exclude'
> exec ^
: string
| set execute = 1
< error
| 1193 HY000 Unknown system variable 'execute'
> exec ^
: string
| set exists = 1
< error
| 1064 42000 …near 'exists = 1' at line 1
> exec ^
: string
| set exit = 1
< error
| 1064 42000 …near 'exit = 1' at line 1
> exec ^
: string
| set expansion = 1
< error
| 1193 HY000 Unknown system variable 'expansion'
> exec ^
: string
| set expire = 1
< error
| 1193 HY000 Unknown system variable 'expire'
> exec ^
: string
| set explain = 1
< error
| 1064 42000 …near 'explain = 1' at line 1
> exec ^
: string
| set export = 1
< error
| 1193 HY000 Unknown system variable 'export'
> exec ^
: string
| set extended = 1
< error
| 1193 HY000 Unknown system variable 'extended'
> exec ^
: string
| set extent_size = 1
< error
| 1193 HY000 Unknown system variable 'extent_size'
> exec ^
: string
| set external = 1
< error
| 1193 HY000 Unknown system variable 'external'
> exec ^
: string
| set external_format = 1
< error
| 1193 HY000 Unknown system variable 'external_format'
> exec ^
: string
| set extract = 1
< error
| 1193 HY000 Unknown system variable 'extract'
> exec ^
: string
| set factor = 1
< error
| 1193 HY000 Unknown system variable 'factor'
> exec ^
: string
| set failed_login_attempts = 1
< error
| 1193 HY000 Unknown system variable 'failed_login_attempts'
> exec ^
: string
| set false = 1
< error
| 1064 42000 …near 'false = 1' at line 1
> exec ^
: string
| set fast = 1
< error
| 1193 HY000 Unknown system variable 'fast'
> exec ^
: string
| set faults = 1
< error
| 1193 HY000 Unknown system variable 'faults'
> exec ^
: string
| set fetch = 1
< error
| 1064 42000 …near 'fetch = 1' at line 1
> exec ^
: string
| set fields = 1
< error
| 1193 HY000 Unknown system variable 'fields'
> exec ^
: string
| set file = 1
< error
| 1193 HY000 Unknown system variable 'file'
> exec ^
: string
| set file_block_size = 1
< error
| 1193 HY000 Unknown system variable 'file_block_size'
> exec ^
: string
| set file_format = 1
< error
| 1193 HY000 Unknown system variable 'file_format'
> exec ^
: string
| set file_name = 1
< error
| 1193 HY000 Unknown system variable 'file_name'
> exec ^
: string
| set file_pattern = 1
< error
| 1193 HY000 Unknown system variable 'file_pattern'
> exec ^
: string
| set file_prefix = 1
< error
| 1193 HY000 Unknown system variable 'file_prefix'
> exec ^
: string
| set files = 1
< error
| 1193 HY000 Unknown system variable 'files'
> exec ^
: string
| set filter = 1
< error
| 1193 HY000 Unknown system variable 'filter'
> exec ^
: string
| set finish = 1
< error
| 1193 HY000 Unknown system variable 'finish'
> exec ^
: string
| set first = 1
< error
| 1193 HY000 Unknown system variable 'first'
> exec ^
: string
| set first_value = 1
< error
| 1064 42000 …near 'first_value = 1' at line 1
> exec ^
: string
| set fixed = 1
< error
| 1193 HY000 Unknown system variable 'fixed'
> exec ^
: string
| set float = 1
< error
| 1064 42000 …near 'float = 1' at line 1
> exec ^
: string
| set float4 = 1
< error
| 1064 42000 …near 'float4 = 1' at line 1
> exec ^
: string
| set float8 = 1
< error
| 1064 42000 …near 'float8 = 1' at line 1
> exec ^
: string
| set flush = 1
< error
| 1229 HY000 Variable 'flush' is a GLOBAL variable and should be set with SET GLOBAL
> exec ^
: string
| set flush = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set following = 1
< error
| 1193 HY000 Unknown system variable 'following'
> exec ^
: string
| set follows = 1
< error
| 1193 HY000 Unknown system variable 'follows'
> exec ^
: string
| set for = 1
< error
| 1064 42000 …near 'for = 1' at line 1
> exec ^
: string
| set force = 1
< error
| 1064 42000 …near 'force = 1' at line 1
> exec ^
: string
| set foreign = 1
< error
| 1064 42000 …near 'foreign = 1' at line 1
> exec ^
: string
| set format = 1
< error
| 1193 HY000 Unknown system variable 'format'
> exec ^
: string
| set found = 1
< error
| 1193 HY000 Unknown system variable 'found'
> exec ^
: string
| set from = 1
< error
| 1064 42000 …near 'from = 1' at line 1
> exec ^
: string
| set full = 1
< error
| 1193 HY000 Unknown system variable 'full'
> exec ^
: string
| set fulltext = 1
< error
| 1064 42000 …near 'fulltext = 1' at line 1
> exec ^
: string
| set function = 1
< error
| 1064 42000 …near 'function = 1' at line 1
> exec ^
: string
| set general = 1
< error
| 1193 HY000 Unknown system variable 'general'
> exec ^
: string
| set generate = 1
< error
| 1193 HY000 Unknown system variable 'generate'
> exec ^
: string
| set generated = 1
< error
| 1064 42000 …near 'generated = 1' at line 1
> exec ^
: string
| set geomcollection = 1
< error
| 1193 HY000 Unknown system variable 'geomcollection'
> exec ^
: string
| set geometry = 1
< error
| 1193 HY000 Unknown system variable 'geometry'
> exec ^
: string
| set geometrycollection = 1
< error
| 1193 HY000 Unknown system variable 'geometrycollection'
> exec ^
: string
| set get = 1
< error
| 1064 42000 …near 'get = 1' at line 1
> exec ^
: string
| set get_format = 1
< error
| 1193 HY000 Unknown system variable 'get_format'
> exec ^
: string
| set get_master_public_key = 1
< error
| 1193 HY000 Unknown system variable 'get_master_public_key'
> exec ^
: string
| set get_source_public_key = 1
< error
| 1193 HY000 Unknown system variable 'get_source_public_key'
> exec ^
: string
| set global = 1
< error
| 1064 42000 …near '= 1' at line 1
> exec ^
: string
| set grant = 1
< error
| 1064 42000 …near 'grant = 1' at line 1
> exec ^
: string
| set grants = 1
< error
| 1193 HY000 Unknown system variable 'grants'
> exec ^
: string
| set greatest = 1
< error
| 1193 HY000 Unknown system variable 'greatest'
> exec ^
: string
| set group = 1
< error
| 1064 42000 …near 'group = 1' at line 1
> exec ^
: string
| set group_concat = 1
< error
| 1193 HY000 Unknown system variable 'group_concat'
> exec ^
: string
| set group_replication = 1
< error
| 1193 HY000 Unknown system variable 'group_replication'
> exec ^
: string
| set grouping = 1
< error
| 1064 42000 …near 'grouping = 1' at line 1
> exec ^
: string
| set groups = 1
< error
| 1064 42000 …near 'groups = 1' at line 1
> exec ^
: string
| set gtid_only = 1
< error
| 1193 HY000 Unknown system variable 'gtid_only'
> exec ^
: string
| set gtids = 1
< error
| 1193 HY000 Unknown system variable 'gtids'
> exec ^
: string
| set guided = 1
< error
| 1193 HY000 Unknown system variable 'guided'
> exec ^
: string
| set handler = 1
< error
| 1193 HY000 Unknown system variable 'handler'
> exec ^
: string
| set hash = 1
< error
| 1193 HY000 Unknown system variable 'hash'
> exec ^
: string
| set having = 1
< error
| 1064 42000 …near 'having = 1' at line 1
> exec ^
: string
| set header = 1
< error
| 1193 HY000 Unknown system variable 'header'
> exec ^
: string
| set help = 1
< error
| 1193 HY000 Unknown system variable 'help'
> exec ^
: string
| set high_priority = 1
< error
| 1064 42000 …near 'high_priority = 1' at line 1
> exec ^
: string
| set histogram = 1
< error
| 1193 HY000 Unknown system variable 'histogram'
> exec ^
: string
| set history = 1
< error
| 1193 HY000 Unknown system variable 'history'
> exec ^
: string
| set host = 1
< error
| 1193 HY000 Unknown system variable 'host'
> exec ^
: string
| set hosts = 1
< error
| 1193 HY000 Unknown system variable 'hosts'
> exec ^
: string
| set hour = 1
< error
| 1193 HY000 Unknown system variable 'hour'
> exec ^
: string
| set hour_microsecond = 1
< error
| 1064 42000 …near 'hour_microsecond = 1' at line 1
> exec ^
: string
| set hour_minute = 1
< error
| 1064 42000 …near 'hour_minute = 1' at line 1
> exec ^
: string
| set hour_second = 1
< error
| 1064 42000 …near 'hour_second = 1' at line 1
> exec ^
: string
| set identified = 1
< error
| 1193 HY000 Unknown system variable 'identified'
> exec ^
: string
| set if = 1
< error
| 1064 42000 …near 'if = 1' at line 1
> exec ^
: string
| set ignore = 1
< error
| 1064 42000 …near 'ignore = 1' at line 1
> exec ^
: string
| set ignore_server_ids = 1
< error
| 1193 HY000 Unknown system variable 'ignore_server_ids'
> exec ^
: string
| set import = 1
< error
| 1193 HY000 Unknown system variable 'import'
> exec ^
: string
| set in = 1
< error
| 1064 42000 …near 'in = 1' at line 1
> exec ^
: string
| set inactive = 1
< error
| 1193 HY000 Unknown system variable 'inactive'
> exec ^
: string
| set index = 1
< error
| 1064 42000 …near 'index = 1' at line 1
> exec ^
: string
| set indexes = 1
< error
| 1193 HY000 Unknown system variable 'indexes'
> exec ^
: string
| set infile = 1
< error
| 1064 42000 …near 'infile = 1' at line 1
> exec ^
: string
| set initial = 1
< error
| 1193 HY000 Unknown system variable 'initial'
> exec ^
: string
| set initial_size = 1
< error
| 1193 HY000 Unknown system variable 'initial_size'
> exec ^
: string
| set initiate = 1
< error
| 1193 HY000 Unknown system variable 'initiate'
> exec ^
: string
| set inner = 1
< error
| 1064 42000 …near 'inner = 1' at line 1
> exec ^
: string
| set inout = 1
< error
| 1064 42000 …near 'inout = 1' at line 1
> exec ^
: string
| set insensitive = 1
< error
| 1064 42000 …near 'insensitive = 1' at line 1
> exec ^
: string
| set insert = 1
< error
| 1064 42000 …near 'insert = 1' at line 1
> exec ^
: string
| set insert_method = 1
< error
| 1193 HY000 Unknown system variable 'insert_method'
> exec ^
: string
| set install = 1
< error
| 1193 HY000 Unknown system variable 'install'
> exec ^
: string
| set instance = 1
< error
| 1193 HY000 Unknown system variable 'instance'
> exec ^
: string
| set int = 1
< error
| 1064 42000 …near 'int = 1' at line 1
> exec ^
: string
| set int1 = 1
< error
| 1064 42000 …near 'int1 = 1' at line 1
> exec ^
: string
| set int2 = 1
< error
| 1064 42000 …near 'int2 = 1' at line 1
> exec ^
: string
| set int3 = 1
< error
| 1064 42000 …near 'int3 = 1' at line 1
> exec ^
: string
| set int4 = 1
< error
| 1064 42000 …near 'int4 = 1' at line 1
> exec ^
: string
| set int8 = 1
< error
| 1064 42000 …near 'int8 = 1' at line 1
> exec ^
: string
| set integer = 1
< error
| 1064 42000 …near 'integer = 1' at line 1
> exec ^
: string
| set intersect = 1
< error
| 1064 42000 …near 'intersect = 1' at line 1
> exec ^
: string
| set interval = 1
< error
| 1064 42000 …near 'interval = 1' at line 1
> exec ^
: string
| set into = 1
< error
| 1064 42000 …near 'into = 1' at line 1
> exec ^
: string
| set invisible = 1
< error
| 1193 HY000 Unknown system variable 'invisible'
> exec ^
: string
| set invoker = 1
< error
| 1193 HY000 Unknown system variable 'invoker'
> exec ^
: string
| set io = 1
< error
| 1193 HY000 Unknown system variable 'io'
> exec ^
: string
| set io_after_gtids = 1
< error
| 1064 42000 …near 'io_after_gtids = 1' at line 1
> exec ^
: string
| set io_before_gtids = 1
< error
| 1064 42000 …near 'io_before_gtids = 1' at line 1
> exec ^
: string
| set io_thread = 1
< error
| 1193 HY000 Unknown system variable 'io_thread'
> exec ^
: string
| set ipc = 1
< error
| 1193 HY000 Unknown system variable 'ipc'
> exec ^
: string
| set is = 1
< error
| 1064 42000 …near 'is = 1' at line 1
> exec ^
: string
| set isolation = 1
< error
| 1193 HY000 Unknown system variable 'isolation'
> exec ^
: string
| set issuer = 1
< error
| 1193 HY000 Unknown system variable 'issuer'
> exec ^
: string
| set iterate = 1
< error
| 1064 42000 …near 'iterate = 1' at line 1
> exec ^
: string
| set join = 1
< error
| 1064 42000 …near 'join = 1' at line 1
> exec ^
: string
| set json = 1
< error
| 1193 HY000 Unknown system variable 'json'
> exec ^
: string
| set json_arrayagg = 1
< error
| 1193 HY000 Unknown system variable 'json_arrayagg'
> exec ^
: string
| set json_duality_object = 1
< error
| 1193 HY000 Unknown system variable 'json_duality_object'
> exec ^
: string
| set json_objectagg = 1
< error
| 1193 HY000 Unknown system variable 'json_objectagg'
> exec ^
: string
| set json_table = 1
< error
| 1064 42000 …near 'json_table = 1' at line 1
> exec ^
: string
| set json_value = 1
< error
| 1193 HY000 Unknown system variable 'json_value'
> exec ^
: string
| set key = 1
< error
| 1064 42000 …near 'key = 1' at line 1
> exec ^
: string
| set key_block_size = 1
< error
| 1193 HY000 Unknown system variable 'key_block_size'
> exec ^
: string
| set keyring = 1
< error
| 1193 HY000 Unknown system variable 'keyring'
> exec ^
: string
| set keys = 1
< error
| 1064 42000 …near 'keys = 1' at line 1
> exec ^
: string
| set kill = 1
< error
| 1064 42000 …near 'kill = 1' at line 1
> exec ^
: string
| set lag = 1
< error
| 1064 42000 …near 'lag = 1' at line 1
> exec ^
: string
| set language = 1
< error
| 1193 HY000 Unknown system variable 'language'
> exec ^
: string
| set last = 1
< error
| 1193 HY000 Unknown system variable 'last'
> exec ^
: string
| set last_value = 1
< error
| 1064 42000 …near 'last_value = 1' at line 1
> exec ^
: string
| set lateral = 1
< error
| 1064 42000 …near 'lateral = 1' at line 1
> exec ^
: string
| set lead = 1
< error
| 1064 42000 …near 'lead = 1' at line 1
> exec ^
: string
| set leading = 1
< error
| 1064 42000 …near 'leading = 1' at line 1
> exec ^
: string
| set least = 1
< error
| 1193 HY000 Unknown system variable 'least'
> exec ^
: string
| set leave = 1
< error
| 1064 42000 …near 'leave = 1' at line 1
> exec ^
: string
| set leaves = 1
< error
| 1193 HY000 Unknown system variable 'leaves'
> exec ^
: string
| set left = 1
< error
| 1064 42000 …near 'left = 1' at line 1
> exec ^
: string
| set less = 1
< error
| 1193 HY000 Unknown system variable 'less'
> exec ^
: string
| set level = 1
< error
| 1193 HY000 Unknown system variable 'level'
> exec ^
: string
| set library = 1
< error
| 1193 HY000 Unknown system variable 'library'
> exec ^
: string
| set like = 1
< error
| 1064 42000 …near 'like = 1' at line 1
> exec ^
: string
| set limit = 1
< error
| 1064 42000 …near 'limit = 1' at line 1
> exec ^
: string
| set linear = 1
< error
| 1064 42000 …near 'linear = 1' at line 1
> exec ^
: string
| set lines = 1
< error
| 1064 42000 …near 'lines = 1' at line 1
> exec ^
: string
| set linestring = 1
< error
| 1193 HY000 Unknown system variable 'linestring'
> exec ^
: string
| set list = 1
< error
| 1193 HY000 Unknown system variable 'list'
> exec ^
: string
| set load = 1
< error
| 1064 42000 …near 'load = 1' at line 1
> exec ^
: string
| set local = 1
< error
| 1064 42000 …near '= 1' at line 1
> exec ^
: string
| set localtime = 1
< error
| 1064 42000 …near 'localtime = 1' at line 1
> exec ^
: string
| set localtimestamp = 1
< error
| 1064 42000 …near 'localtimestamp = 1' at line 1
> exec ^
: string
| set lock = 1
< error
| 1064 42000 …near 'lock = 1' at line 1
> exec ^
: string
| set locked = 1
< error
| 1193 HY000 Unknown system variable 'locked'
> exec ^
: string
| set locks = 1
< error
| 1193 HY000 Unknown system variable 'locks'
> exec ^
: string
| set log = 1
< error
| 1193 HY000 Unknown system variable 'log'
> exec ^
: string
| set logfile = 1
< error
| 1193 HY000 Unknown system variable 'logfile'
> exec ^
: string
| set logs = 1
< error
| 1193 HY000 Unknown system variable 'logs'
> exec ^
: string
| set long = 1
< error
| 1064 42000 …near 'long = 1' at line 1
> exec ^
: string
| set longblob = 1
< error
| 1064 42000 …near 'longblob = 1' at line 1
> exec ^
: string
| set longtext = 1
< error
| 1064 42000 …near 'longtext = 1' at line 1
> exec ^
: string
| set loop = 1
< error
| 1064 42000 …near 'loop = 1' at line 1
> exec ^
: string
| set low_priority = 1
< error
| 1064 42000 …near 'low_priority = 1' at line 1
> exec ^
: string
| set manifest = 1
< error
| 1193 HY000 Unknown system variable 'manifest'
> exec ^
: string
| set manual = 1
< error
| 1193 HY000 Unknown system variable 'manual'
> exec ^
: string
| set master = 1
< error
| 1193 HY000 Unknown system variable 'master'
> exec ^
: string
| set master_auto_position = 1
< error
| 1193 HY000 Unknown system variable 'master_auto_position'
> exec ^
: string
| set master_bind = 1
< error
| 1064 42000 …near 'master_bind = 1' at line 1
> exec ^
: string
| set master_compression_algorithms = 1
< error
| 1193 HY000 Unknown system variable 'master_compression_algorithms'
> exec ^
: string
| set master_connect_retry = 1
< error
| 1193 HY000 Unknown system variable 'master_connect_retry'
> exec ^
: string
| set master_delay = 1
< error
| 1193 HY000 Unknown system variable 'master_delay'
> exec ^
: string
| set master_heartbeat_period = 1
< error
| 1193 HY000 Unknown system variable 'master_heartbeat_period'
> exec ^
: string
| set master_host = 1
< error
| 1193 HY000 Unknown system variable 'master_host'
> exec ^
: string
| set master_log_file = 1
< error
| 1193 HY000 Unknown system variable 'master_log_file'
> exec ^
: string
| set master_log_pos = 1
< error
| 1193 HY000 Unknown system variable 'master_log_pos'
> exec ^
: string
| set master_password = 1
< error
| 1193 HY000 Unknown system variable 'master_password'
> exec ^
: string
| set master_port = 1
< error
| 1193 HY000 Unknown system variable 'master_port'
> exec ^
: string
| set master_public_key_path = 1
< error
| 1193 HY000 Unknown system variable 'master_public_key_path'
> exec ^
: string
| set master_retry_count = 1
< error
| 1193 HY000 Unknown system variable 'master_retry_count'
> exec ^
: string
| set master_ssl = 1
< error
| 1193 HY000 Unknown system variable 'master_ssl'
> exec ^
: string
| set master_ssl_ca = 1
< error
| 1193 HY000 Unknown system variable 'master_ssl_ca'
> exec ^
: string
| set master_ssl_capath = 1
< error
| 1193 HY000 Unknown system variable 'master_ssl_capath'
> exec ^
: string
| set master_ssl_cert = 1
< error
| 1193 HY000 Unknown system variable 'master_ssl_cert'
> exec ^
: string
| set master_ssl_cipher = 1
< error
| 1193 HY000 Unknown system variable 'master_ssl_cipher'
> exec ^
: string
| set master_ssl_crl = 1
< error
| 1193 HY000 Unknown system variable 'master_ssl_crl'
> exec ^
: string
| set master_ssl_crlpath = 1
< error
| 1193 HY000 Unknown system variable 'master_ssl_crlpath'
> exec ^
: string
| set master_ssl_key = 1
< error
| 1193 HY000 Unknown system variable 'master_ssl_key'
> exec ^
: string
| set master_ssl_verify_server_cert = 1
< error
| 1064 42000 …near 'master_ssl_verify_server_cert = 1' at line 1
> exec ^
: string
| set master_tls_ciphersuites = 1
< error
| 1193 HY000 Unknown system variable 'master_tls_ciphersuites'
> exec ^
: string
| set master_tls_version = 1
< error
| 1193 HY000 Unknown system variable 'master_tls_version'
> exec ^
: string
| set master_user = 1
< error
| 1193 HY000 Unknown system variable 'master_user'
> exec ^
: string
| set master_zstd_compression_level = 1
< error
| 1193 HY000 Unknown system variable 'master_zstd_compression_level'
> exec ^
: string
| set match = 1
< error
| 1064 42000 …near 'match = 1' at line 1
> exec ^
: string
| set materialized = 1
< error
| 1193 HY000 Unknown system variable 'materialized'
> exec ^
: string
| set max = 1
< error
| 1193 HY000 Unknown system variable 'max'
> exec ^
: string
| set max_connections_per_hour = 1
< error
| 1193 HY000 Unknown system variable 'max_connections_per_hour'
> exec ^
: string
| set max_queries_per_hour = 1
< error
| 1193 HY000 Unknown system variable 'max_queries_per_hour'
> exec ^
: string
| set max_rows = 1
< error
| 1193 HY000 Unknown system variable 'max_rows'
> exec ^
: string
| set max_size = 1
< error
| 1193 HY000 Unknown system variable 'max_size'
> exec ^
: string
| set max_updates_per_hour = 1
< error
| 1193 HY000 Unknown system variable 'max_updates_per_hour'
> exec ^
: string
| set max_user_connections = 1
< ok 0 0
> exec ^
: string
| set maxvalue = 1
< error
| 1064 42000 …near 'maxvalue = 1' at line 1
> exec ^
: string
| set medium = 1
< error
| 1193 HY000 Unknown system variable 'medium'
> exec ^
: string
| set mediumblob = 1
< error
| 1064 42000 …near 'mediumblob = 1' at line 1
> exec ^
: string
| set mediumint = 1
< error
| 1064 42000 …near 'mediumint = 1' at line 1
> exec ^
: string
| set mediumtext = 1
< error
| 1064 42000 …near 'mediumtext = 1' at line 1
> exec ^
: string
| set member = 1
< error
| 1193 HY000 Unknown system variable 'member'
> exec ^
: string
| set memory = 1
< error
| 1193 HY000 Unknown system variable 'memory'
> exec ^
: string
| set merge = 1
< error
| 1193 HY000 Unknown system variable 'merge'
> exec ^
: string
| set message_text = 1
< error
| 1193 HY000 Unknown system variable 'message_text'
> exec ^
: string
| set microsecond = 1
< error
| 1193 HY000 Unknown system variable 'microsecond'
> exec ^
: string
| set mid = 1
< error
| 1193 HY000 Unknown system variable 'mid'
> exec ^
: string
| set middleint = 1
< error
| 1064 42000 …near 'middleint = 1' at line 1
> exec ^
: string
| set migrate = 1
< error
| 1193 HY000 Unknown system variable 'migrate'
> exec ^
: string
| set min = 1
< error
| 1193 HY000 Unknown system variable 'min'
> exec ^
: string
| set min_rows = 1
< error
| 1193 HY000 Unknown system variable 'min_rows'
> exec ^
: string
| set minute = 1
< error
| 1193 HY000 Unknown system variable 'minute'
> exec ^
: string
| set minute_microsecond = 1
< error
| 1064 42000 …near 'minute_microsecond = 1' at line 1
> exec ^
: string
| set minute_second = 1
< error
| 1064 42000 …near 'minute_second = 1' at line 1
> exec ^
: string
| set mod = 1
< error
| 1064 42000 …near 'mod = 1' at line 1
> exec ^
: string
| set mode = 1
< error
| 1193 HY000 Unknown system variable 'mode'
> exec ^
: string
| set modifies = 1
< error
| 1064 42000 …near 'modifies = 1' at line 1
> exec ^
: string
| set modify = 1
< error
| 1193 HY000 Unknown system variable 'modify'
> exec ^
: string
| set month = 1
< error
| 1193 HY000 Unknown system variable 'month'
> exec ^
: string
| set multilinestring = 1
< error
| 1193 HY000 Unknown system variable 'multilinestring'
> exec ^
: string
| set multipoint = 1
< error
| 1193 HY000 Unknown system variable 'multipoint'
> exec ^
: string
| set multipolygon = 1
< error
| 1193 HY000 Unknown system variable 'multipolygon'
> exec ^
: string
| set mutex = 1
< error
| 1193 HY000 Unknown system variable 'mutex'
> exec ^
: string
| set mysql_errno = 1
< error
| 1193 HY000 Unknown system variable 'mysql_errno'
> exec ^
: string
| set name = 1
< error
| 1193 HY000 Unknown system variable 'name'
> exec ^
: string
| set names = 1
< error
| 1064 42000 …near '= 1' at line 1
> exec ^
: string
| set national = 1
< error
| 1193 HY000 Unknown system variable 'national'
> exec ^
: string
| set natural = 1
< error
| 1064 42000 …near 'natural = 1' at line 1
> exec ^
: string
| set nchar = 1
< error
| 1193 HY000 Unknown system variable 'nchar'
> exec ^
: string
| set ndb = 1
< error
| 1193 HY000 Unknown system variable 'ndb'
> exec ^
: string
| set ndbcluster = 1
< error
| 1193 HY000 Unknown system variable 'ndbcluster'
> exec ^
: string
| set nested = 1
< error
| 1193 HY000 Unknown system variable 'nested'
> exec ^
: string
| set network_namespace = 1
< error
| 1193 HY000 Unknown system variable 'network_namespace'
> exec ^
: string
| set never = 1
< error
| 1193 HY000 Unknown system variable 'never'
> exec ^
: string
| set new = 1
< ok 0 0
> exec ^
: string
| set next = 1
< error
| 1193 HY000 Unknown system variable 'next'
> exec ^
: string
| set no = 1
< error
| 1193 HY000 Unknown system variable 'no'
> exec ^
: string
| set no_wait = 1
< error
| 1193 HY000 Unknown system variable 'no_wait'
> exec ^
: string
| set no_write_to_binlog = 1
< error
| 1064 42000 …near 'no_write_to_binlog = 1' at line 1
> exec ^
: string
| set nodegroup = 1
< error
| 1193 HY000 Unknown system variable 'nodegroup'
> exec ^
: string
| set none = 1
< error
| 1193 HY000 Unknown system variable 'none'
> exec ^
: string
| set not = 1
< error
| 1064 42000 …near 'not = 1' at line 1
> exec ^
: string
| set now = 1
< error
| 1193 HY000 Unknown system variable 'now'
> exec ^
: string
| set nowait = 1
< error
| 1193 HY000 Unknown system variable 'nowait'
> exec ^
: string
| set nth_value = 1
< error
| 1064 42000 …near 'nth_value = 1' at line 1
> exec ^
: string
| set ntile = 1
< error
| 1064 42000 …near 'ntile = 1' at line 1
> exec ^
: string
| set null = 1
< error
| 1064 42000 …near 'null = 1' at line 1
> exec ^
: string
| set nullif = 1
< error
| 1193 HY000 Unknown system variable 'nullif'
> exec ^
: string
| set nulls = 1
< error
| 1193 HY000 Unknown system variable 'nulls'
> exec ^
: string
| set number = 1
< error
| 1193 HY000 Unknown system variable 'number'
> exec ^
: string
| set numeric = 1
< error
| 1064 42000 …near 'numeric = 1' at line 1
> exec ^
: string
| set nvarchar = 1
< error
| 1193 HY000 Unknown system variable 'nvarchar'
> exec ^
: string
| set of = 1
< error
| 1064 42000 …near 'of = 1' at line 1
> exec ^
: string
| set off = 1
< error
| 1193 HY000 Unknown system variable 'off'
> exec ^
: string
| set offset = 1
< error
| 1193 HY000 Unknown system variable 'offset'
> exec ^
: string
| set oj = 1
< error
| 1193 HY000 Unknown system variable 'oj'
> exec ^
: string
| set old = 1
< error
| 1229 HY000 Variable 'old' is a GLOBAL variable and should be set with SET GLOBAL
> exec ^
: string
| set old = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set on = 1
< error
| 1064 42000 …near 'on = 1' at line 1
> exec ^
: string
| set one = 1
< error
| 1193 HY000 Unknown system variable 'one'
> exec ^
: string
| set only = 1
< error
| 1193 HY000 Unknown system variable 'only'
> exec ^
: string
| set open = 1
< error
| 1193 HY000 Unknown system variable 'open'
> exec ^
: string
| set optimize = 1
< error
| 1064 42000 …near 'optimize = 1' at line 1
> exec ^
: string
| set optimizer_costs = 1
< error
| 1064 42000 …near 'optimizer_costs = 1' at line 1
> exec ^
: string
| set option = 1
< error
| 1064 42000 …near 'option = 1' at line 1
> exec ^
: string
| set optional = 1
< error
| 1193 HY000 Unknown system variable 'optional'
> exec ^
: string
| set optionally = 1
< error
| 1064 42000 …near 'optionally = 1' at line 1
> exec ^
: string
| set options = 1
< error
| 1193 HY000 Unknown system variable 'options'
> exec ^
: string
| set or = 1
< error
| 1064 42000 …near 'or = 1' at line 1
> exec ^
: string
| set order = 1
< error
| 1064 42000 …near 'order = 1' at line 1
> exec ^
: string
| set ordinality = 1
< error
| 1193 HY000 Unknown system variable 'ordinality'
> exec ^
: string
| set organization = 1
< error
| 1193 HY000 Unknown system variable 'organization'
> exec ^
: string
| set others = 1
< error
| 1193 HY000 Unknown system variable 'others'
> exec ^
: string
| set out = 1
< error
| 1064 42000 …near 'out = 1' at line 1
> exec ^
: string
| set outer = 1
< error
| 1064 42000 …near 'outer = 1' at line 1
> exec ^
: string
| set outfile = 1
< error
| 1064 42000 …near 'outfile = 1' at line 1
> exec ^
: string
| set over = 1
< error
| 1064 42000 …near 'over = 1' at line 1
> exec ^
: string
| set overlay = 1
< error
| 1193 HY000 Unknown system variable 'overlay'
> exec ^
: string
| set overwrite = 1
< error
| 1193 HY000 Unknown system variable 'overwrite'
> exec ^
: string
| set owner = 1
< error
| 1193 HY000 Unknown system variable 'owner'
> exec ^
: string
| set pack_keys = 1
< error
| 1193 HY000 Unknown system variable 'pack_keys'
> exec ^
: string
| set page = 1
< error
| 1193 HY000 Unknown system variable 'page'
> exec ^
: string
| set parallel = 1
< error
| 1193 HY000 Unknown system variable 'parallel'
> exec ^
: string
| set parameters = 1
< error
| 1193 HY000 Unknown system variable 'parameters'
> exec ^
: string
| set parse_tree = 1
< error
| 1193 HY000 Unknown system variable 'parse_tree'
> exec ^
: string
| set parser = 1
< error
| 1193 HY000 Unknown system variable 'parser'
> exec ^
: string
| set partial = 1
< error
| 1193 HY000 Unknown system variable 'partial'
> exec ^
: string
| set partition = 1
< error
| 1064 42000 …near 'partition = 1' at line 1
> exec ^
: string
| set partitioning = 1
< error
| 1193 HY000 Unknown system variable 'partitioning'
> exec ^
: string
| set partitions = 1
< error
| 1193 HY000 Unknown system variable 'partitions'
> exec ^
: string
| set password = 1
< error
| 1064 42000 …near '1' at line 1
> exec ^
: string
| set password_lock_time = 1
< error
| 1193 HY000 Unknown system variable 'password_lock_time'
> exec ^
: string
| set path = 1
< error
| 1193 HY000 Unknown system variable 'path'
> exec ^
: string
| set percent_rank = 1
< error
| 1064 42000 …near 'percent_rank = 1' at line 1
> exec ^
: string
| set persist = 1
< error
| 1064 42000 …near '= 1' at line 1
> exec ^
: string
| set persist_only = 1
< error
| 1064 42000 …near '= 1' at line 1
> exec ^
: string
| set phase = 1
< error
| 1193 HY000 Unknown system variable 'phase'
> exec ^
: string
| set placing = 1
< error
| 1193 HY000 Unknown system variable 'placing'
> exec ^
: string
| set plugin = 1
< error
| 1193 HY000 Unknown system variable 'plugin'
> exec ^
: string
| set plugin_dir = 1
< error
| 1229 HY000 Variable 'plugin_dir' is a GLOBAL variable and should be set with SET GLOBAL
> exec ^
: string
| set plugin_dir = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set plugins = 1
< error
| 1193 HY000 Unknown system variable 'plugins'
> exec ^
: string
| set point = 1
< error
| 1193 HY000 Unknown system variable 'point'
> exec ^
: string
| set polygon = 1
< error
| 1193 HY000 Unknown system variable 'polygon'
> exec ^
: string
| set port = 1
< error
| 1229 HY000 Variable 'port' is a GLOBAL variable and should be set with SET GLOBAL
> exec ^
: string
| set port = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set position = 1
< error
| 1193 HY000 Unknown system variable 'position'
> exec ^
: string
| set precedes = 1
< error
| 1193 HY000 Unknown system variable 'precedes'
> exec ^
: string
| set preceding = 1
< error
| 1193 HY000 Unknown system variable 'preceding'
> exec ^
: string
| set precision = 1
< error
| 1064 42000 …near 'precision = 1' at line 1
> exec ^
: string
| set prefix = 1
< error
| 1193 HY000 Unknown system variable 'prefix'
> exec ^
: string
| set prepare = 1
< error
| 1193 HY000 Unknown system variable 'prepare'
> exec ^
: string
| set preserve = 1
< error
| 1193 HY000 Unknown system variable 'preserve'
> exec ^
: string
| set prev = 1
< error
| 1193 HY000 Unknown system variable 'prev'
> exec ^
: string
| set primary = 1
< error
| 1064 42000 …near 'primary = 1' at line 1
> exec ^
: string
| set privilege_checks_user = 1
< error
| 1193 HY000 Unknown system variable 'privilege_checks_user'
> exec ^
: string
| set privileges = 1
< error
| 1193 HY000 Unknown system variable 'privileges'
> exec ^
: string
| set procedure = 1
< error
| 1064 42000 …near 'procedure = 1' at line 1
> exec ^
: string
| set process = 1
< error
| 1193 HY000 Unknown system variable 'process'
> exec ^
: string
| set processlist = 1
< error
| 1193 HY000 Unknown system variable 'processlist'
> exec ^
: string
| set profile = 1
< error
| 1193 HY000 Unknown system variable 'profile'
> exec ^
: string
| set profiles = 1
< error
| 1193 HY000 Unknown system variable 'profiles'
> exec ^
: string
| set proxy = 1
< error
| 1193 HY000 Unknown system variable 'proxy'
> exec ^
: string
| set purge = 1
< error
| 1064 42000 …near 'purge = 1' at line 1
> exec ^
: string
| set qualify = 1
< error
| 1193 HY000 Unknown system variable 'qualify'
> exec ^
: string
| set quarter = 1
< error
| 1193 HY000 Unknown system variable 'quarter'
> exec ^
: string
| set query = 1
< error
| 1193 HY000 Unknown system variable 'query'
> exec ^
: string
| set quick = 1
< error
| 1193 HY000 Unknown system variable 'quick'
> exec ^
: string
| set random = 1
< error
| 1193 HY000 Unknown system variable 'random'
> exec ^
: string
| set range = 1
< error
| 1064 42000 …near 'range = 1' at line 1
> exec ^
: string
| set rank = 1
< error
| 1064 42000 …near 'rank = 1' at line 1
> exec ^
: string
| set read = 1
< error
| 1064 42000 …near 'read = 1' at line 1
> exec ^
: string
| set read_only = 1
< error
| 1229 HY000 Variable 'read_only' is a GLOBAL variable and should be set with SET GLOBAL
> exec ^
: string
| set read_only = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set read_write = 1
< error
| 1064 42000 …near 'read_write = 1' at line 1
> exec ^
: string
| set reads = 1
< error
| 1064 42000 …near 'reads = 1' at line 1
> exec ^
: string
| set real = 1
< error
| 1064 42000 …near 'real = 1' at line 1
> exec ^
: string
| set rebuild = 1
< error
| 1193 HY000 Unknown system variable 'rebuild'
> exec ^
: string
| set recover = 1
< error
| 1193 HY000 Unknown system variable 'recover'
> exec ^
: string
| set recursive = 1
< error
| 1064 42000 …near 'recursive = 1' at line 1
> exec ^
: string
| set redo_buffer_size = 1
< error
| 1193 HY000 Unknown system variable 'redo_buffer_size'
> exec ^
: string
| set redundant = 1
< error
| 1193 HY000 Unknown system variable 'redundant'
> exec ^
: string
| set reference = 1
< error
| 1193 HY000 Unknown system variable 'reference'
> exec ^
: string
| set references = 1
< error
| 1064 42000 …near 'references = 1' at line 1
> exec ^
: string
| set regexp = 1
< error
| 1064 42000 …near 'regexp = 1' at line 1
> exec ^
: string
| set registration = 1
< error
| 1193 HY000 Unknown system variable 'registration'
> exec ^
: string
| set relational = 1
< error
| 1193 HY000 Unknown system variable 'relational'
> exec ^
: string
| set relay = 1
< error
| 1193 HY000 Unknown system variable 'relay'
> exec ^
: string
| set relay_log_file = 1
< error
| 1193 HY000 Unknown system variable 'relay_log_file'
> exec ^
: string
| set relay_log_pos = 1
< error
| 1193 HY000 Unknown system variable 'relay_log_pos'
> exec ^
: string
| set relay_thread = 1
< error
| 1193 HY000 Unknown system variable 'relay_thread'
> exec ^
: string
| set relaylog = 1
< error
| 1193 HY000 Unknown system variable 'relaylog'
> exec ^
: string
| set release = 1
< error
| 1064 42000 …near 'release = 1' at line 1
> exec ^
: string
| set reload = 1
< error
| 1193 HY000 Unknown system variable 'reload'
> exec ^
: string
| set remove = 1
< error
| 1193 HY000 Unknown system variable 'remove'
> exec ^
: string
| set rename = 1
< error
| 1064 42000 …near 'rename = 1' at line 1
> exec ^
: string
| set reorganize = 1
< error
| 1193 HY000 Unknown system variable 'reorganize'
> exec ^
: string
| set repair = 1
< error
| 1193 HY000 Unknown system variable 'repair'
> exec ^
: string
| set repeat = 1
< error
| 1064 42000 …near 'repeat = 1' at line 1
> exec ^
: string
| set repeatable = 1
< error
| 1193 HY000 Unknown system variable 'repeatable'
> exec ^
: string
| set replace = 1
< error
| 1064 42000 …near 'replace = 1' at line 1
> exec ^
: string
| set replica = 1
< error
| 1193 HY000 Unknown system variable 'replica'
> exec ^
: string
| set replicas = 1
< error
| 1193 HY000 Unknown system variable 'replicas'
> exec ^
: string
| set replicate_do_db = 1
< error
| 1193 HY000 Unknown system variable 'replicate_do_db'
> exec ^
: string
| set replicate_do_table = 1
< error
| 1193 HY000 Unknown system variable 'replicate_do_table'
> exec ^
: string
| set replicate_ignore_db = 1
< error
| 1193 HY000 Unknown system variable 'replicate_ignore_db'
> exec ^
: string
| set replicate_ignore_table = 1
< error
| 1193 HY000 Unknown system variable 'replicate_ignore_table'
> exec ^
: string
| set replicate_rewrite_db = 1
< error
| 1193 HY000 Unknown system variable 'replicate_rewrite_db'
> exec ^
: string
| set replicate_wild_do_table = 1
< error
| 1193 HY000 Unknown system variable 'replicate_wild_do_table'
> exec ^
: string
| set replicate_wild_ignore_table = 1
< error
| 1193 HY000 Unknown system variable 'replicate_wild_ignore_table'
> exec ^
: string
| set replication = 1
< error
| 1193 HY000 Unknown system variable 'replication'
> exec ^
: string
| set require = 1
< error
| 1064 42000 …near 'require = 1' at line 1
> exec ^
: string
| set require_row_format = 1
< ok 0 0
> exec ^
: string
| set require_table_primary_key_check = 1
< error
| 1193 HY000 Unknown system variable 'require_table_primary_key_check'
> exec ^
: string
| set reset = 1
< error
| 1193 HY000 Unknown system variable 'reset'
> exec ^
: string
| set resignal = 1
< error
| 1064 42000 …near 'resignal = 1' at line 1
> exec ^
: string
| set resource = 1
< error
| 1193 HY000 Unknown system variable 'resource'
> exec ^
: string
| set respect = 1
< error
| 1193 HY000 Unknown system variable 'respect'
> exec ^
: string
| set restart = 1
< error
| 1193 HY000 Unknown system variable 'restart'
> exec ^
: string
| set restore = 1
< error
| 1193 HY000 Unknown system variable 'restore'
> exec ^
: string
| set restrict = 1
< error
| 1064 42000 …near 'restrict = 1' at line 1
> exec ^
: string
| set resume = 1
< error
| 1193 HY000 Unknown system variable 'resume'
> exec ^
: string
| set retain = 1
< error
| 1193 HY000 Unknown system variable 'retain'
> exec ^
: string
| set return = 1
< error
| 1064 42000 …near 'return = 1' at line 1
> exec ^
: string
| set returned_sqlstate = 1
< error
| 1193 HY000 Unknown system variable 'returned_sqlstate'
> exec ^
: string
| set returning = 1
< error
| 1193 HY000 Unknown system variable 'returning'
> exec ^
: string
| set returns = 1
< error
| 1193 HY000 Unknown system variable 'returns'
> exec ^
: string
| set reuse = 1
< error
| 1193 HY000 Unknown system variable 'reuse'
> exec ^
: string
| set reverse = 1
< error
| 1193 HY000 Unknown system variable 'reverse'
> exec ^
: string
| set revoke = 1
< error
| 1064 42000 …near 'revoke = 1' at line 1
> exec ^
: string
| set right = 1
< error
| 1064 42000 …near 'right = 1' at line 1
> exec ^
: string
| set rlike = 1
< error
| 1064 42000 …near 'rlike = 1' at line 1
> exec ^
: string
| set role = 1
< error
| 1193 HY000 Unknown system variable 'role'
> exec ^
: string
| set rollback = 1
< error
| 1193 HY000 Unknown system variable 'rollback'
> exec ^
: string
| set rollup = 1
< error
| 1193 HY000 Unknown system variable 'rollup'
> exec ^
: string
| set rotate = 1
< error
| 1193 HY000 Unknown system variable 'rotate'
> exec ^
: string
| set routine = 1
< error
| 1193 HY000 Unknown system variable 'routine'
> exec ^
: string
| set row = 1
< error
| 1064 42000 …near 'row = 1' at line 1
> exec ^
: string
| set row_count = 1
< error
| 1193 HY000 Unknown system variable 'row_count'
> exec ^
: string
| set row_format = 1
< error
| 1193 HY000 Unknown system variable 'row_format'
> exec ^
: string
| set row_number = 1
< error
| 1064 42000 …near 'row_number = 1' at line 1
> exec ^
: string
| set rows = 1
< error
| 1064 42000 …near 'rows = 1' at line 1
> exec ^
: string
| set rtree = 1
< error
| 1193 HY000 Unknown system variable 'rtree'
> exec ^
: string
| set s3 = 1
< error
| 1193 HY000 Unknown system variable 's3'
> exec ^
: string
| set savepoint = 1
< error
| 1193 HY000 Unknown system variable 'savepoint'
> exec ^
: string
| set schedule = 1
< error
| 1193 HY000 Unknown system variable 'schedule'
> exec ^
: string
| set schema = 1
< error
| 1064 42000 …near 'schema = 1' at line 1
> exec ^
: string
| set schema_name = 1
< error
| 1193 HY000 Unknown system variable 'schema_name'
> exec ^
: string
| set schemas = 1
< error
| 1064 42000 …near 'schemas = 1' at line 1
> exec ^
: string
| set second = 1
< error
| 1193 HY000 Unknown system variable 'second'
> exec ^
: string
| set second_microsecond = 1
< error
| 1064 42000 …near 'second_microsecond = 1' at line 1
> exec ^
: string
| set secondary = 1
< error
| 1193 HY000 Unknown system variable 'secondary'
> exec ^
: string
| set secondary_engine = 1
< error
| 1193 HY000 Unknown system variable 'secondary_engine'
> exec ^
: string
| set secondary_engine_attribute = 1
< error
| 1193 HY000 Unknown system variable 'secondary_engine_attribute'
> exec ^
: string
| set secondary_load = 1
< error
| 1193 HY000 Unknown system variable 'secondary_load'
> exec ^
: string
| set secondary_unload = 1
< error
| 1193 HY000 Unknown system variable 'secondary_unload'
> exec ^
: string
| set security = 1
< error
| 1193 HY000 Unknown system variable 'security'
> exec ^
: string
| set select = 1
< error
| 1064 42000 …near 'select = 1' at line 1
> exec ^
: string
| set sensitive = 1
< error
| 1064 42000 …near 'sensitive = 1' at line 1
> exec ^
: string
| set separator = 1
< error
| 1064 42000 …near 'separator = 1' at line 1
> exec ^
: string
| set serial = 1
< error
| 1193 HY000 Unknown system variable 'serial'
> exec ^
: string
| set serializable = 1
< error
| 1193 HY000 Unknown system variable 'serializable'
> exec ^
: string
| set server = 1
< error
| 1193 HY000 Unknown system variable 'server'
> exec ^
: string
| set session = 1
< error
| 1064 42000 …near '= 1' at line 1
> exec ^
: string
| set session_user = 1
< error
| 1193 HY000 Unknown system variable 'session_user'
> exec ^
: string
| set set = 1
< error
| 1064 42000 …near 'set = 1' at line 1
> exec ^
: string
| set sets = 1
< error
| 1193 HY000 Unknown system variable 'sets'
> exec ^
: string
| set share = 1
< error
| 1193 HY000 Unknown system variable 'share'
> exec ^
: string
| set show = 1
< error
| 1064 42000 …near 'show = 1' at line 1
> exec ^
: string
| set shutdown = 1
< error
| 1193 HY000 Unknown system variable 'shutdown'
> exec ^
: string
| set signal = 1
< error
| 1064 42000 …near 'signal = 1' at line 1
> exec ^
: string
| set signed = 1
< error
| 1193 HY000 Unknown system variable 'signed'
> exec ^
: string
| set simple = 1
< error
| 1193 HY000 Unknown system variable 'simple'
> exec ^
: string
| set skip = 1
< error
| 1193 HY000 Unknown system variable 'skip'
> exec ^
: string
| set slave = 1
< error
| 1193 HY000 Unknown system variable 'slave'
> exec ^
: string
| set slow = 1
< error
| 1193 HY000 Unknown system variable 'slow'
> exec ^
: string
| set smallint = 1
< error
| 1064 42000 …near 'smallint = 1' at line 1
> exec ^
: string
| set snapshot = 1
< error
| 1193 HY000 Unknown system variable 'snapshot'
> exec ^
: string
| set socket = 1
< error
| 1229 HY000 Variable 'socket' is a GLOBAL variable and should be set with SET GLOBAL
> exec ^
: string
| set socket = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set some = 1
< error
| 1193 HY000 Unknown system variable 'some'
> exec ^
: string
| set soname = 1
< error
| 1193 HY000 Unknown system variable 'soname'
> exec ^
: string
| set sounds = 1
< error
| 1193 HY000 Unknown system variable 'sounds'
> exec ^
: string
| set source = 1
< error
| 1193 HY000 Unknown system variable 'source'
> exec ^
: string
| set source_auto_position = 1
< error
| 1193 HY000 Unknown system variable 'source_auto_position'
> exec ^
: string
| set source_bind = 1
< error
| 1193 HY000 Unknown system variable 'source_bind'
> exec ^
: string
| set source_compression_algorithms = 1
< error
| 1193 HY000 Unknown system variable 'source_compression_algorithms'
> exec ^
: string
| set source_connect_retry = 1
< error
| 1193 HY000 Unknown system variable 'source_connect_retry'
> exec ^
: string
| set source_connection_auto_failover = 1
< error
| 1193 HY000 Unknown system variable 'source_connection_auto_failover'
> exec ^
: string
| set source_delay = 1
< error
| 1193 HY000 Unknown system variable 'source_delay'
> exec ^
: string
| set source_heartbeat_period = 1
< error
| 1193 HY000 Unknown system variable 'source_heartbeat_period'
> exec ^
: string
| set source_host = 1
< error
| 1193 HY000 Unknown system variable 'source_host'
> exec ^
: string
| set source_log_file = 1
< error
| 1193 HY000 Unknown system variable 'source_log_file'
> exec ^
: string
| set source_log_pos = 1
< error
| 1193 HY000 Unknown system variable 'source_log_pos'
> exec ^
: string
| set source_password = 1
< error
| 1193 HY000 Unknown system variable 'source_password'
> exec ^
: string
| set source_port = 1
< error
| 1193 HY000 Unknown system variable 'source_port'
> exec ^
: string
| set source_public_key_path = 1
< error
| 1193 HY000 Unknown system variable 'source_public_key_path'
> exec ^
: string
| set source_retry_count = 1
< error
| 1193 HY000 Unknown system variable 'source_retry_count'
> exec ^
: string
| set source_ssl = 1
< error
| 1193 HY000 Unknown system variable 'source_ssl'
> exec ^
: string
| set source_ssl_ca = 1
< error
| 1193 HY000 Unknown system variable 'source_ssl_ca'
> exec ^
: string
| set source_ssl_capath = 1
< error
| 1193 HY000 Unknown system variable 'source_ssl_capath'
> exec ^
: string
| set source_ssl_cert = 1
< error
| 1193 HY000 Unknown system variable 'source_ssl_cert'
> exec ^
: string
| set source_ssl_cipher = 1
< error
| 1193 HY000 Unknown system variable 'source_ssl_cipher'
> exec ^
: string
| set source_ssl_crl = 1
< error
| 1193 HY000 Unknown system variable 'source_ssl_crl'
> exec ^
: string
| set source_ssl_crlpath = 1
< error
| 1193 HY000 Unknown system variable 'source_ssl_crlpath'
> exec ^
: string
| set source_ssl_key = 1
< error
| 1193 HY000 Unknown system variable 'source_ssl_key'
> exec ^
: string
| set source_ssl_verify_server_cert = 1
< error
| 1193 HY000 Unknown system variable 'source_ssl_verify_server_cert'
> exec ^
: string
| set source_tls_ciphersuites = 1
< error
| 1193 HY000 Unknown system variable 'source_tls_ciphersuites'
> exec ^
: string
| set source_tls_version = 1
< error
| 1193 HY000 Unknown system variable 'source_tls_version'
> exec ^
: string
| set source_user = 1
< error
| 1193 HY000 Unknown system variable 'source_user'
> exec ^
: string
| set source_zstd_compression_level = 1
< error
| 1193 HY000 Unknown system variable 'source_zstd_compression_level'
> exec ^
: string
| set spatial = 1
< error
| 1064 42000 …near 'spatial = 1' at line 1
> exec ^
: string
| set specific = 1
< error
| 1064 42000 …near 'specific = 1' at line 1
> exec ^
: string
| set sql = 1
< error
| 1064 42000 …near 'sql = 1' at line 1
> exec ^
: string
| set sql_after_gtids = 1
< error
| 1193 HY000 Unknown system variable 'sql_after_gtids'
> exec ^
: string
| set sql_after_mts_gaps = 1
< error
| 1193 HY000 Unknown system variable 'sql_after_mts_gaps'
> exec ^
: string
| set sql_before_gtids = 1
< error
| 1193 HY000 Unknown system variable 'sql_before_gtids'
> exec ^
: string
| set sql_big_result = 1
< error
| 1064 42000 …near 'sql_big_result = 1' at line 1
> exec ^
: string
| set sql_buffer_result = 1
< ok 0 0
> exec ^
: string
| set sql_calc_found_rows = 1
< error
| 1064 42000 …near 'sql_calc_found_rows = 1' at line 1
> exec ^
: string
| set sql_no_cache = 1
< error
| 1193 HY000 Unknown system variable 'sql_no_cache'
> exec ^
: string
| set sql_small_result = 1
< error
| 1064 42000 …near 'sql_small_result = 1' at line 1
> exec ^
: string
| set sql_thread = 1
< error
| 1193 HY000 Unknown system variable 'sql_thread'
> exec ^
: string
| set sql_tsi_day = 1
< error
| 1193 HY000 Unknown system variable 'sql_tsi_day'
> exec ^
: string
| set sql_tsi_hour = 1
< error
| 1193 HY000 Unknown system variable 'sql_tsi_hour'
> exec ^
: string
| set sql_tsi_minute = 1
< error
| 1193 HY000 Unknown system variable 'sql_tsi_minute'
> exec ^
: string
| set sql_tsi_month = 1
< error
| 1193 HY000 Unknown system variable 'sql_tsi_month'
> exec ^
: string
| set sql_tsi_quarter = 1
< error
| 1193 HY000 Unknown system variable 'sql_tsi_quarter'
> exec ^
: string
| set sql_tsi_second = 1
< error
| 1193 HY000 Unknown system variable 'sql_tsi_second'
> exec ^
: string
| set sql_tsi_week = 1
< error
| 1193 HY000 Unknown system variable 'sql_tsi_week'
> exec ^
: string
| set sql_tsi_year = 1
< error
| 1193 HY000 Unknown system variable 'sql_tsi_year'
> exec ^
: string
| set sqlexception = 1
< error
| 1064 42000 …near 'sqlexception = 1' at line 1
> exec ^
: string
| set sqlstate = 1
< error
| 1064 42000 …near 'sqlstate = 1' at line 1
> exec ^
: string
| set sqlwarning = 1
< error
| 1064 42000 …near 'sqlwarning = 1' at line 1
> exec ^
: string
| set srid = 1
< error
| 1193 HY000 Unknown system variable 'srid'
> exec ^
: string
| set sse_kms = 1
< error
| 1193 HY000 Unknown system variable 'sse_kms'
> exec ^
: string
| set sse_s3 = 1
< error
| 1193 HY000 Unknown system variable 'sse_s3'
> exec ^
: string
| set ssl = 1
< error
| 1064 42000 …near 'ssl = 1' at line 1
> exec ^
: string
| set st_collect = 1
< error
| 1193 HY000 Unknown system variable 'st_collect'
> exec ^
: string
| set stacked = 1
< error
| 1193 HY000 Unknown system variable 'stacked'
> exec ^
: string
| set start = 1
< error
| 1193 HY000 Unknown system variable 'start'
> exec ^
: string
| set starting = 1
< error
| 1064 42000 …near 'starting = 1' at line 1
> exec ^
: string
| set starts = 1
< error
| 1193 HY000 Unknown system variable 'starts'
> exec ^
: string
| set stats_auto_recalc = 1
< error
| 1193 HY000 Unknown system variable 'stats_auto_recalc'
> exec ^
: string
| set stats_persistent = 1
< error
| 1193 HY000 Unknown system variable 'stats_persistent'
> exec ^
: string
| set stats_sample_pages = 1
< error
| 1193 HY000 Unknown system variable 'stats_sample_pages'
> exec ^
: string
| set status = 1
< error
| 1193 HY000 Unknown system variable 'status'
> exec ^
: string
| set std = 1
< error
| 1193 HY000 Unknown system variable 'std'
> exec ^
: string
| set stddev = 1
< error
| 1193 HY000 Unknown system variable 'stddev'
> exec ^
: string
| set stddev_pop = 1
< error
| 1193 HY000 Unknown system variable 'stddev_pop'
> exec ^
: string
| set stddev_samp = 1
< error
| 1193 HY000 Unknown system variable 'stddev_samp'
> exec ^
: string
| set stop = 1
< error
| 1193 HY000 Unknown system variable 'stop'
> exec ^
: string
| set storage = 1
< error
| 1193 HY000 Unknown system variable 'storage'
> exec ^
: string
| set stored = 1
< error
| 1064 42000 …near 'stored = 1' at line 1
> exec ^
: string
| set straight_join = 1
< error
| 1064 42000 …near 'straight_join = 1' at line 1
> exec ^
: string
| set stream = 1
< error
| 1193 HY000 Unknown system variable 'stream'
> exec ^
: string
| set strict_load = 1
< error
| 1193 HY000 Unknown system variable 'strict_load'
> exec ^
: string
| set string = 1
< error
| 1193 HY000 Unknown system variable 'string'
> exec ^
: string
| set subclass_origin = 1
< error
| 1193 HY000 Unknown system variable 'subclass_origin'
> exec ^
: string
| set subdate = 1
< error
| 1193 HY000 Unknown system variable 'subdate'
> exec ^
: string
| set subject = 1
< error
| 1193 HY000 Unknown system variable 'subject'
> exec ^
: string
| set subpartition = 1
< error
| 1193 HY000 Unknown system variable 'subpartition'
> exec ^
: string
| set subpartitions = 1
< error
| 1193 HY000 Unknown system variable 'subpartitions'
> exec ^
: string
| set substr = 1
< error
| 1193 HY000 Unknown system variable 'substr'
> exec ^
: string
| set substring = 1
< error
| 1193 HY000 Unknown system variable 'substring'
> exec ^
: string
| set sum = 1
< error
| 1193 HY000 Unknown system variable 'sum'
> exec ^
: string
| set super = 1
< error
| 1193 HY000 Unknown system variable 'super'
> exec ^
: string
| set suspend = 1
< error
| 1193 HY000 Unknown system variable 'suspend'
> exec ^
: string
| set swaps = 1
< error
| 1193 HY000 Unknown system variable 'swaps'
> exec ^
: string
| set switches = 1
< error
| 1193 HY000 Unknown system variable 'switches'
> exec ^
: string
| set sysdate = 1
< error
| 1193 HY000 Unknown system variable 'sysdate'
> exec ^
: string
| set system = 1
< error
| 1064 42000 …near 'system = 1' at line 1
> exec ^
: string
| set system_user = 1
< error
| 1193 HY000 Unknown system variable 'system_user'
> exec ^
: string
| set table = 1
< error
| 1064 42000 …near 'table = 1' at line 1
> exec ^
: string
| set table_checksum = 1
< error
| 1193 HY000 Unknown system variable 'table_checksum'
> exec ^
: string
| set table_name = 1
< error
| 1193 HY000 Unknown system variable 'table_name'
> exec ^
: string
| set tables = 1
< error
| 1193 HY000 Unknown system variable 'tables'
> exec ^
: string
| set tablesample = 1
< error
| 1193 HY000 Unknown system variable 'tablesample'
> exec ^
: string
| set tablespace = 1
< error
| 1193 HY000 Unknown system variable 'tablespace'
> exec ^
: string
| set temporary = 1
< error
| 1193 HY000 Unknown system variable 'temporary'
> exec ^
: string
| set temptable = 1
< error
| 1193 HY000 Unknown system variable 'temptable'
> exec ^
: string
| set terminated = 1
< error
| 1064 42000 …near 'terminated = 1' at line 1
> exec ^
: string
| set text = 1
< error
| 1193 HY000 Unknown system variable 'text'
> exec ^
: string
| set than = 1
< error
| 1193 HY000 Unknown system variable 'than'
> exec ^
: string
| set then = 1
< error
| 1064 42000 …near 'then = 1' at line 1
> exec ^
: string
| set thread_priority = 1
< error
| 1193 HY000 Unknown system variable 'thread_priority'
> exec ^
: string
| set ties = 1
< error
| 1193 HY000 Unknown system variable 'ties'
> exec ^
: string
| set time = 1
< error
| 1193 HY000 Unknown system variable 'time'
> exec ^
: string
| set timestamp = 1
< ok 0 0
> exec ^
: string
| set timestampadd = 1
< error
| 1193 HY000 Unknown system variable 'timestampadd'
> exec ^
: string
| set timestampdiff = 1
< error
| 1193 HY000 Unknown system variable 'timestampdiff'
> exec ^
: string
| set tinyblob = 1
< error
| 1064 42000 …near 'tinyblob = 1' at line 1
> exec ^
: string
| set tinyint = 1
< error
| 1064 42000 …near 'tinyint = 1' at line 1
> exec ^
: string
| set tinytext = 1
< error
| 1064 42000 …near 'tinytext = 1' at line 1
> exec ^
: string
| set tls = 1
< error
| 1193 HY000 Unknown system variable 'tls'
> exec ^
: string
| set to = 1
< error
| 1064 42000 …near 'to = 1' at line 1
> exec ^
: string
| set trailing = 1
< error
| 1064 42000 …near 'trailing = 1' at line 1
> exec ^
: string
| set transaction = 1
< error
| 1193 HY000 Unknown system variable 'transaction'
> exec ^
: string
| set trigger = 1
< error
| 1064 42000 …near 'trigger = 1' at line 1
> exec ^
: string
| set triggers = 1
< error
| 1193 HY000 Unknown system variable 'triggers'
> exec ^
: string
| set trim = 1
< error
| 1193 HY000 Unknown system variable 'trim'
> exec ^
: string
| set true = 1
< error
| 1064 42000 …near 'true = 1' at line 1
> exec ^
: string
| set truncate = 1
< error
| 1193 HY000 Unknown system variable 'truncate'
> exec ^
: string
| set type = 1
< error
| 1193 HY000 Unknown system variable 'type'
> exec ^
: string
| set types = 1
< error
| 1193 HY000 Unknown system variable 'types'
> exec ^
: string
| set unbounded = 1
< error
| 1193 HY000 Unknown system variable 'unbounded'
> exec ^
: string
| set uncommitted = 1
< error
| 1193 HY000 Unknown system variable 'uncommitted'
> exec ^
: string
| set undefined = 1
< error
| 1193 HY000 Unknown system variable 'undefined'
> exec ^
: string
| set undo = 1
< error
| 1064 42000 …near 'undo = 1' at line 1
> exec ^
: string
| set undo_buffer_size = 1
< error
| 1193 HY000 Unknown system variable 'undo_buffer_size'
> exec ^
: string
| set undofile = 1
< error
| 1193 HY000 Unknown system variable 'undofile'
> exec ^
: string
| set unicode = 1
< error
| 1193 HY000 Unknown system variable 'unicode'
> exec ^
: string
| set uninstall = 1
< error
| 1193 HY000 Unknown system variable 'uninstall'
> exec ^
: string
| set union = 1
< error
| 1064 42000 …near 'union = 1' at line 1
> exec ^
: string
| set unique = 1
< error
| 1064 42000 …near 'unique = 1' at line 1
> exec ^
: string
| set unknown = 1
< error
| 1193 HY000 Unknown system variable 'unknown'
> exec ^
: string
| set unlock = 1
< error
| 1064 42000 …near 'unlock = 1' at line 1
> exec ^
: string
| set unregister = 1
< error
| 1193 HY000 Unknown system variable 'unregister'
> exec ^
: string
| set unsigned = 1
< error
| 1064 42000 …near 'unsigned = 1' at line 1
> exec ^
: string
| set until = 1
< error
| 1193 HY000 Unknown system variable 'until'
> exec ^
: string
| set update = 1
< error
| 1064 42000 …near 'update = 1' at line 1
> exec ^
: string
| set upgrade = 1
< error
| 1193 HY000 Unknown system variable 'upgrade'
> exec ^
: string
| set uri = 1
< error
| 1193 HY000 Unknown system variable 'uri'
> exec ^
: string
| set url = 1
< error
| 1193 HY000 Unknown system variable 'url'
> exec ^
: string
| set usage = 1
< error
| 1064 42000 …near 'usage = 1' at line 1
> exec ^
: string
| set use = 1
< error
| 1064 42000 …near 'use = 1' at line 1
> exec ^
: string
| set use_frm = 1
< error
| 1193 HY000 Unknown system variable 'use_frm'
> exec ^
: string
| set user = 1
< error
| 1193 HY000 Unknown system variable 'user'
> exec ^
: string
| set user_resources = 1
< error
| 1193 HY000 Unknown system variable 'user_resources'
> exec ^
: string
| set using = 1
< error
| 1064 42000 …near 'using = 1' at line 1
> exec ^
: string
| set utc_date = 1
< error
| 1064 42000 …near 'utc_date = 1' at line 1
> exec ^
: string
| set utc_time = 1
< error
| 1064 42000 …near 'utc_time = 1' at line 1
> exec ^
: string
| set utc_timestamp = 1
< error
| 1064 42000 …near 'utc_timestamp = 1' at line 1
> exec ^
: string
| set validate = 1
< error
| 1193 HY000 Unknown system variable 'validate'
> exec ^
: string
| set validation = 1
< error
| 1193 HY000 Unknown system variable 'validation'
> exec ^
: string
| set value = 1
< error
| 1193 HY000 Unknown system variable 'value'
> exec ^
: string
| set values = 1
< error
| 1064 42000 …near 'values = 1' at line 1
> exec ^
: string
| set var_pop = 1
< error
| 1193 HY000 Unknown system variable 'var_pop'
> exec ^
: string
| set var_samp = 1
< error
| 1193 HY000 Unknown system variable 'var_samp'
> exec ^
: string
| set varbinary = 1
< error
| 1064 42000 …near 'varbinary = 1' at line 1
> exec ^
: string
| set varchar = 1
< error
| 1064 42000 …near 'varchar = 1' at line 1
> exec ^
: string
| set varcharacter = 1
< error
| 1064 42000 …near 'varcharacter = 1' at line 1
> exec ^
: string
| set variables = 1
< error
| 1193 HY000 Unknown system variable 'variables'
> exec ^
: string
| set variance = 1
< error
| 1193 HY000 Unknown system variable 'variance'
> exec ^
: string
| set varying = 1
< error
| 1064 42000 …near 'varying = 1' at line 1
> exec ^
: string
| set vcpu = 1
< error
| 1193 HY000 Unknown system variable 'vcpu'
> exec ^
: string
| set vector = 1
< error
| 1193 HY000 Unknown system variable 'vector'
> exec ^
: string
| set verify_key_constraints = 1
< error
| 1193 HY000 Unknown system variable 'verify_key_constraints'
> exec ^
: string
| set view = 1
< error
| 1193 HY000 Unknown system variable 'view'
> exec ^
: string
| set virtual = 1
< error
| 1064 42000 …near 'virtual = 1' at line 1
> exec ^
: string
| set visible = 1
< error
| 1193 HY000 Unknown system variable 'visible'
> exec ^
: string
| set wait = 1
< error
| 1193 HY000 Unknown system variable 'wait'
> exec ^
: string
| set warnings = 1
< error
| 1193 HY000 Unknown system variable 'warnings'
> exec ^
: string
| set week = 1
< error
| 1193 HY000 Unknown system variable 'week'
> exec ^
: string
| set weight_string = 1
< error
| 1193 HY000 Unknown system variable 'weight_string'
> exec ^
: string
| set when = 1
< error
| 1064 42000 …near 'when = 1' at line 1
> exec ^
: string
| set where = 1
< error
| 1064 42000 …near 'where = 1' at line 1
> exec ^
: string
| set while = 1
< error
| 1064 42000 …near 'while = 1' at line 1
> exec ^
: string
| set window = 1
< error
| 1064 42000 …near 'window = 1' at line 1
> exec ^
: string
| set with = 1
< error
| 1064 42000 …near 'with = 1' at line 1
> exec ^
: string
| set without = 1
< error
| 1193 HY000 Unknown system variable 'without'
> exec ^
: string
| set work = 1
< error
| 1193 HY000 Unknown system variable 'work'
> exec ^
: string
| set wrapper = 1
< error
| 1193 HY000 Unknown system variable 'wrapper'
> exec ^
: string
| set write = 1
< error
| 1064 42000 …near 'write = 1' at line 1
> exec ^
: string
| set x509 = 1
< error
| 1193 HY000 Unknown system variable 'x509'
> exec ^
: string
| set xa = 1
< error
| 1193 HY000 Unknown system variable 'xa'
> exec ^
: string
| set xid = 1
< error
| 1193 HY000 Unknown system variable 'xid'
> exec ^
: string
| set xml = 1
< error
| 1193 HY000 Unknown system variable 'xml'
> exec ^
: string
| set xor = 1
< error
| 1064 42000 …near 'xor = 1' at line 1
> exec ^
: string
| set year = 1
< error
| 1193 HY000 Unknown system variable 'year'
> exec ^
: string
| set year_month = 1
< error
| 1064 42000 …near 'year_month = 1' at line 1
> exec ^
: string
| set zerofill = 1
< error
| 1064 42000 …near 'zerofill = 1' at line 1
> exec ^
: string
| set zone = 1
< error
| 1193 HY000 Unknown system variable 'zone'
> exec ^
: string
| set global absent = 1
< error
| 1193 HY000 Unknown system variable 'absent'
> exec ^
: string
| set global accessible = 1
< error
| 1064 42000 …near 'accessible = 1' at line 1
> exec ^
: string
| set global account = 1
< error
| 1193 HY000 Unknown system variable 'account'
> exec ^
: string
| set global action = 1
< error
| 1193 HY000 Unknown system variable 'action'
> exec ^
: string
| set global active = 1
< error
| 1193 HY000 Unknown system variable 'active'
> exec ^
: string
| set global add = 1
< error
| 1064 42000 …near 'add = 1' at line 1
> exec ^
: string
| set global adddate = 1
< error
| 1193 HY000 Unknown system variable 'adddate'
> exec ^
: string
| set global admin = 1
< error
| 1193 HY000 Unknown system variable 'admin'
> exec ^
: string
| set global after = 1
< error
| 1193 HY000 Unknown system variable 'after'
> exec ^
: string
| set global against = 1
< error
| 1193 HY000 Unknown system variable 'against'
> exec ^
: string
| set global aggregate = 1
< error
| 1193 HY000 Unknown system variable 'aggregate'
> exec ^
: string
| set global algorithm = 1
< error
| 1193 HY000 Unknown system variable 'algorithm'
> exec ^
: string
| set global all = 1
< error
| 1064 42000 …near 'all = 1' at line 1
> exec ^
: string
| set global allow_missing_files = 1
< error
| 1193 HY000 Unknown system variable 'allow_missing_files'
> exec ^
: string
| set global alter = 1
< error
| 1064 42000 …near 'alter = 1' at line 1
> exec ^
: string
| set global always = 1
< error
| 1193 HY000 Unknown system variable 'always'
> exec ^
: string
| set global analyze = 1
< error
| 1064 42000 …near 'analyze = 1' at line 1
> exec ^
: string
| set global and = 1
< error
| 1064 42000 …near 'and = 1' at line 1
> exec ^
: string
| set global any = 1
< error
| 1193 HY000 Unknown system variable 'any'
> exec ^
: string
| set global array = 1
< error
| 1193 HY000 Unknown system variable 'array'
> exec ^
: string
| set global as = 1
< error
| 1064 42000 …near 'as = 1' at line 1
> exec ^
: string
| set global asc = 1
< error
| 1064 42000 …near 'asc = 1' at line 1
> exec ^
: string
| set global ascii = 1
< error
| 1193 HY000 Unknown system variable 'ascii'
> exec ^
: string
| set global asensitive = 1
< error
| 1064 42000 …near 'asensitive = 1' at line 1
> exec ^
: string
| set global assign_gtids_to_anonymous_transactions = 1
< error
| 1193 HY000 Unknown system variable 'assign_gtids_to_anonymous_transactions'
> exec ^
: string
| set global at = 1
< error
| 1193 HY000 Unknown system variable 'at'
> exec ^
: string
| set global attribute = 1
< error
| 1193 HY000 Unknown system variable 'attribute'
> exec ^
: string
| set global authentication = 1
< error
| 1193 HY000 Unknown system variable 'authentication'
> exec ^
: string
| set global auto = 1
< error
| 1193 HY000 Unknown system variable 'auto'
> exec ^
: string
| set global auto_increment = 1
< error
| 1193 HY000 Unknown system variable 'auto_increment'
> exec ^
: string
| set global auto_refresh = 1
< error
| 1193 HY000 Unknown system variable 'auto_refresh'
> exec ^
: string
| set global auto_refresh_source = 1
< error
| 1193 HY000 Unknown system variable 'auto_refresh_source'
> exec ^
: string
| set global autoextend_size = 1
< error
| 1193 HY000 Unknown system variable 'autoextend_size'
> exec ^
: string
| set global avg = 1
< error
| 1193 HY000 Unknown system variable 'avg'
> exec ^
: string
| set global avg_row_length = 1
< error
| 1193 HY000 Unknown system variable 'avg_row_length'
> exec ^
: string
| set global backup = 1
< error
| 1193 HY000 Unknown system variable 'backup'
> exec ^
: string
| set global before = 1
< error
| 1064 42000 …near 'before = 1' at line 1
> exec ^
: string
| set global begin = 1
< error
| 1193 HY000 Unknown system variable 'begin'
> exec ^
: string
| set global bernoulli = 1
< error
| 1193 HY000 Unknown system variable 'bernoulli'
> exec ^
: string
| set global between = 1
< error
| 1064 42000 …near 'between = 1' at line 1
> exec ^
: string
| set global bigint = 1
< error
| 1064 42000 …near 'bigint = 1' at line 1
> exec ^
: string
| set global binary = 1
< error
| 1064 42000 …near 'binary = 1' at line 1
> exec ^
: string
| set global binlog = 1
< error
| 1193 HY000 Unknown system variable 'binlog'
> exec ^
: string
| set global bit = 1
< error
| 1193 HY000 Unknown system variable 'bit'
> exec ^
: string
| set global bit_and = 1
< error
| 1193 HY000 Unknown system variable 'bit_and'
> exec ^
: string
| set global bit_or = 1
< error
| 1193 HY000 Unknown system variable 'bit_or'
> exec ^
: string
| set global bit_xor = 1
< error
| 1193 HY000 Unknown system variable 'bit_xor'
> exec ^
: string
| set global blob = 1
< error
| 1064 42000 …near 'blob = 1' at line 1
> exec ^
: string
| set global block = 1
< error
| 1193 HY000 Unknown system variable 'block'
> exec ^
: string
| set global bool = 1
< error
| 1193 HY000 Unknown system variable 'bool'
> exec ^
: string
| set global boolean = 1
< error
| 1193 HY000 Unknown system variable 'boolean'
> exec ^
: string
| set global both = 1
< error
| 1064 42000 …near 'both = 1' at line 1
> exec ^
: string
| set global btree = 1
< error
| 1193 HY000 Unknown system variable 'btree'
> exec ^
: string
| set global buckets = 1
< error
| 1193 HY000 Unknown system variable 'buckets'
> exec ^
: string
| set global bulk = 1
< error
| 1193 HY000 Unknown system variable 'bulk'
> exec ^
: string
| set global by = 1
< error
| 1064 42000 …near 'by = 1' at line 1
> exec ^
: string
| set global byte = 1
< error
| 1193 HY000 Unknown system variable 'byte'
> exec ^
: string
| set global cache = 1
< error
| 1193 HY000 Unknown system variable 'cache'
> exec ^
: string
| set global call = 1
< error
| 1064 42000 …near 'call = 1' at line 1
> exec ^
: string
| set global cascade = 1
< error
| 1064 42000 …near 'cascade = 1' at line 1
> exec ^
: string
| set global cascaded = 1
< error
| 1193 HY000 Unknown system variable 'cascaded'
> exec ^
: string
| set global case = 1
< error
| 1064 42000 …near 'case = 1' at line 1
> exec ^
: string
| set global cast = 1
< error
| 1193 HY000 Unknown system variable 'cast'
> exec ^
: string
| set global catalog_name = 1
< error
| 1193 HY000 Unknown system variable 'catalog_name'
> exec ^
: string
| set global chain = 1
< error
| 1193 HY000 Unknown system variable 'chain'
> exec ^
: string
| set global challenge_response = 1
< error
| 1193 HY000 Unknown system variable 'challenge_response'
> exec ^
: string
| set global change = 1
< error
| 1064 42000 …near 'change = 1' at line 1
> exec ^
: string
| set global changed = 1
< error
| 1193 HY000 Unknown system variable 'changed'
> exec ^
: string
| set global channel = 1
< error
| 1193 HY000 Unknown system variable 'channel'
> exec ^
: string
| set global char = 1
< error
| 1064 42000 …near 'char = 1' at line 1
> exec ^
: string
| set global character = 1
< error
| 1064 42000 …near 'character = 1' at line 1
> exec ^
: string
| set global charset = 1
< error
| 1193 HY000 Unknown system variable 'charset'
> exec ^
: string
| set global check = 1
< error
| 1064 42000 …near 'check = 1' at line 1
> exec ^
: string
| set global checksum = 1
< error
| 1193 HY000 Unknown system variable 'checksum'
> exec ^
: string
| set global cipher = 1
< error
| 1193 HY000 Unknown system variable 'cipher'
> exec ^
: string
| set global class_origin = 1
< error
| 1193 HY000 Unknown system variable 'class_origin'
> exec ^
: string
| set global client = 1
< error
| 1193 HY000 Unknown system variable 'client'
> exec ^
: string
| set global clone = 1
< error
| 1193 HY000 Unknown system variable 'clone'
> exec ^
: string
| set global close = 1
< error
| 1193 HY000 Unknown system variable 'close'
> exec ^
: string
| set global coalesce = 1
< error
| 1193 HY000 Unknown system variable 'coalesce'
> exec ^
: string
| set global code = 1
< error
| 1193 HY000 Unknown system variable 'code'
> exec ^
: string
| set global collate = 1
< error
| 1064 42000 …near 'collate = 1' at line 1
> exec ^
: string
| set global collation = 1
< error
| 1193 HY000 Unknown system variable 'collation'
> exec ^
: string
| set global column = 1
< error
| 1064 42000 …near 'column = 1' at line 1
> exec ^
: string
| set global column_format = 1
< error
| 1193 HY000 Unknown system variable 'column_format'
> exec ^
: string
| set global column_name = 1
< error
| 1193 HY000 Unknown system variable 'column_name'
> exec ^
: string
| set global columns = 1
< error
| 1193 HY000 Unknown system variable 'columns'
> exec ^
: string
| set global comment = 1
< error
| 1193 HY000 Unknown system variable 'comment'
> exec ^
: string
| set global commit = 1
< error
| 1193 HY000 Unknown system variable 'commit'
> exec ^
: string
| set global committed = 1
< error
| 1193 HY000 Unknown system variable 'committed'
> exec ^
: string
| set global compact = 1
< error
| 1193 HY000 Unknown system variable 'compact'
> exec ^
: string
| set global completion = 1
< error
| 1193 HY000 Unknown system variable 'completion'
> exec ^
: string
| set global component = 1
< error
| 1193 HY000 Unknown system variable 'component'
> exec ^
: string
| set global compressed = 1
< error
| 1193 HY000 Unknown system variable 'compressed'
> exec ^
: string
| set global compression = 1
< error
| 1193 HY000 Unknown system variable 'compression'
> exec ^
: string
| set global concat = 1
< error
| 1193 HY000 Unknown system variable 'concat'
> exec ^
: string
| set global concurrent = 1
< error
| 1193 HY000 Unknown system variable 'concurrent'
> exec ^
: string
| set global condition = 1
< error
| 1064 42000 …near 'condition = 1' at line 1
> exec ^
: string
| set global connection = 1
< error
| 1193 HY000 Unknown system variable 'connection'
> exec ^
: string
| set global consistent = 1
< error
| 1193 HY000 Unknown system variable 'consistent'
> exec ^
: string
| set global constraint = 1
< error
| 1064 42000 …near 'constraint = 1' at line 1
> exec ^
: string
| set global constraint_catalog = 1
< error
| 1193 HY000 Unknown system variable 'constraint_catalog'
> exec ^
: string
| set global constraint_name = 1
< error
| 1193 HY000 Unknown system variable 'constraint_name'
> exec ^
: string
| set global constraint_schema = 1
< error
| 1193 HY000 Unknown system variable 'constraint_schema'
> exec ^
: string
| set global contains = 1
< error
| 1193 HY000 Unknown system variable 'contains'
> exec ^
: string
| set global context = 1
< error
| 1193 HY000 Unknown system variable 'context'
> exec ^
: string
| set global continue = 1
< error
| 1064 42000 …near 'continue = 1' at line 1
> exec ^
: string
| set global convert = 1
< error
| 1064 42000 …near 'convert = 1' at line 1
> exec ^
: string
| set global count = 1
< error
| 1193 HY000 Unknown system variable 'count'
> exec ^
: string
| set global cpu = 1
< error
| 1193 HY000 Unknown system variable 'cpu'
> exec ^
: string
| set global create = 1
< error
| 1064 42000 …near 'create = 1' at line 1
> exec ^
: string
| set global cross = 1
< error
| 1064 42000 …near 'cross = 1' at line 1
> exec ^
: string
| set global csv = 1
< error
| 1193 HY000 Unknown system variable 'csv'
> exec ^
: string
| set global cube = 1
< error
| 1064 42000 …near 'cube = 1' at line 1
> exec ^
: string
| set global cume_dist = 1
< error
| 1064 42000 …near 'cume_dist = 1' at line 1
> exec ^
: string
| set global curdate = 1
< error
| 1193 HY000 Unknown system variable 'curdate'
> exec ^
: string
| set global current = 1
< error
| 1193 HY000 Unknown system variable 'current'
> exec ^
: string
| set global current_date = 1
< error
| 1064 42000 …near 'current_date = 1' at line 1
> exec ^
: string
| set global current_time = 1
< error
| 1064 42000 …near 'current_time = 1' at line 1
> exec ^
: string
| set global current_timestamp = 1
< error
| 1064 42000 …near 'current_timestamp = 1' at line 1
> exec ^
: string
| set global current_user = 1
< error
| 1064 42000 …near 'current_user = 1' at line 1
> exec ^
: string
| set global cursor = 1
< error
| 1064 42000 …near 'cursor = 1' at line 1
> exec ^
: string
| set global cursor_name = 1
< error
| 1193 HY000 Unknown system variable 'cursor_name'
> exec ^
: string
| set global curtime = 1
< error
| 1193 HY000 Unknown system variable 'curtime'
> exec ^
: string
| set global data = 1
< error
| 1193 HY000 Unknown system variable 'data'
> exec ^
: string
| set global database = 1
< error
| 1064 42000 …near 'database = 1' at line 1
> exec ^
: string
| set global databases = 1
< error
| 1064 42000 …near 'databases = 1' at line 1
> exec ^
: string
| set global datafile = 1
< error
| 1193 HY000 Unknown system variable 'datafile'
> exec ^
: string
| set global date = 1
< error
| 1193 HY000 Unknown system variable 'date'
> exec ^
: string
| set global date_add = 1
< error
| 1193 HY000 Unknown system variable 'date_add'
> exec ^
: string
| set global date_sub = 1
< error
| 1193 HY000 Unknown system variable 'date_sub'
> exec ^
: string
| set global datetime = 1
< error
| 1193 HY000 Unknown system variable 'datetime'
> exec ^
: string
| set global day = 1
< error
| 1193 HY000 Unknown system variable 'day'
> exec ^
: string
| set global day_hour = 1
< error
| 1064 42000 …near 'day_hour = 1' at line 1
> exec ^
: string
| set global day_microsecond = 1
< error
| 1064 42000 …near 'day_microsecond = 1' at line 1
> exec ^
: string
| set global day_minute = 1
< error
| 1064 42000 …near 'day_minute = 1' at line 1
> exec ^
: string
| set global day_second = 1
< error
| 1064 42000 …near 'day_second = 1' at line 1
> exec ^
: string
| set global deallocate = 1
< error
| 1193 HY000 Unknown system variable 'deallocate'
> exec ^
: string
| set global dec = 1
< error
| 1064 42000 …near 'dec = 1' at line 1
> exec ^
: string
| set global decimal = 1
< error
| 1064 42000 …near 'decimal = 1' at line 1
> exec ^
: string
| set global declare = 1
< error
| 1064 42000 …near 'declare = 1' at line 1
> exec ^
: string
| set global default = 1
< error
| 1064 42000 …near '= 1' at line 1
> exec ^
: string
| set global default_auth = 1
< error
| 1193 HY000 Unknown system variable 'default_auth'
> exec ^
: string
| set global definer = 1
< error
| 1193 HY000 Unknown system variable 'definer'
> exec ^
: string
| set global definition = 1
< error
| 1193 HY000 Unknown system variable 'definition'
> exec ^
: string
| set global delay_key_write = 1
< ok 0 0
> exec ^
: string
| set global delayed = 1
< error
| 1064 42000 …near 'delayed = 1' at line 1
> exec ^
: string
| set global delete = 1
< error
| 1064 42000 …near 'delete = 1' at line 1
> exec ^
: string
| set global dense_rank = 1
< error
| 1064 42000 …near 'dense_rank = 1' at line 1
> exec ^
: string
| set global desc = 1
< error
| 1064 42000 …near 'desc = 1' at line 1
> exec ^
: string
| set global describe = 1
< error
| 1064 42000 …near 'describe = 1' at line 1
> exec ^
: string
| set global description = 1
< error
| 1193 HY000 Unknown system variable 'description'
> exec ^
: string
| set global deterministic = 1
< error
| 1064 42000 …near 'deterministic = 1' at line 1
> exec ^
: string
| set global diagnostics = 1
< error
| 1193 HY000 Unknown system variable 'diagnostics'
> exec ^
: string
| set global directory = 1
< error
| 1193 HY000 Unknown system variable 'directory'
> exec ^
: string
| set global disable = 1
< error
| 1193 HY000 Unknown system variable 'disable'
> exec ^
: string
| set global discard = 1
< error
| 1193 HY000 Unknown system variable 'discard'
> exec ^
: string
| set global disk = 1
< error
| 1193 HY000 Unknown system variable 'disk'
> exec ^
: string
| set global distinct = 1
< error
| 1064 42000 …near 'distinct = 1' at line 1
> exec ^
: string
| set global distinctrow = 1
< error
| 1064 42000 …near 'distinctrow = 1' at line 1
> exec ^
: string
| set global div = 1
< error
| 1064 42000 …near 'div = 1' at line 1
> exec ^
: string
| set global do = 1
< error
| 1193 HY000 Unknown system variable 'do'
> exec ^
: string
| set global double = 1
< error
| 1064 42000 …near 'double = 1' at line 1
> exec ^
: string
| set global drop = 1
< error
| 1064 42000 …near 'drop = 1' at line 1
> exec ^
: string
| set global dual = 1
< error
| 1064 42000 …near 'dual = 1' at line 1
> exec ^
: string
| set global duality = 1
< error
| 1193 HY000 Unknown system variable 'duality'
> exec ^
: string
| set global dumpfile = 1
< error
| 1193 HY000 Unknown system variable 'dumpfile'
> exec ^
: string
| set global duplicate = 1
< error
| 1193 HY000 Unknown system variable 'duplicate'
> exec ^
: string
| set global dynamic = 1
< error
| 1193 HY000 Unknown system variable 'dynamic'
> exec ^
: string
| set global each = 1
< error
| 1064 42000 …near 'each = 1' at line 1
> exec ^
: string
| set global else = 1
< error
| 1064 42000 …near 'else = 1' at line 1
> exec ^
: string
| set global elseif = 1
< error
| 1064 42000 …near 'elseif = 1' at line 1
> exec ^
: string
| set global empty = 1
< error
| 1064 42000 …near 'empty = 1' at line 1
> exec ^
: string
| set global enable = 1
< error
| 1193 HY000 Unknown system variable 'enable'
> exec ^
: string
| set global enclosed = 1
< error
| 1064 42000 …near 'enclosed = 1' at line 1
> exec ^
: string
| set global encryption = 1
< error
| 1193 HY000 Unknown system variable 'encryption'
> exec ^
: string
| set global end = 1
< error
| 1193 HY000 Unknown system variable 'end'
> exec ^
: string
| set global ends = 1
< error
| 1193 HY000 Unknown system variable 'ends'
> exec ^
: string
| set global enforced = 1
< error
| 1193 HY000 Unknown system variable 'enforced'
> exec ^
: string
| set global engine = 1
< error
| 1193 HY000 Unknown system variable 'engine'
> exec ^
: string
| set global engine_attribute = 1
< error
| 1193 HY000 Unknown system variable 'engine_attribute'
> exec ^
: string
| set global engines = 1
< error
| 1193 HY000 Unknown system variable 'engines'
> exec ^
: string
| set global enum = 1
< error
| 1193 HY000 Unknown system variable 'enum'
> exec ^
: string
| set global error = 1
< error
| 1193 HY000 Unknown system variable 'error'
> exec ^
: string
| set global errors = 1
< error
| 1193 HY000 Unknown system variable 'errors'
> exec ^
: string
| set global escape = 1
< error
| 1193 HY000 Unknown system variable 'escape'
> exec ^
: string
| set global escaped = 1
< error
| 1064 42000 …near 'escaped = 1' at line 1
> exec ^
: string
| set global event = 1
< error
| 1193 HY000 Unknown system variable 'event'
> exec ^
: string
| set global events = 1
< error
| 1193 HY000 Unknown system variable 'events'
> exec ^
: string
| set global every = 1
< error
| 1193 HY000 Unknown system variable 'every'
> exec ^
: string
| set global except = 1
< error
| 1064 42000 …near 'except = 1' at line 1
> exec ^
: string
| set global exchange = 1
< error
| 1193 HY000 Unknown system variable 'exchange'
> exec ^
: string
| set global exclude = 1
< error
| 1193 HY000 Unknown system variable 'exclude'
> exec ^
: string
| set global execute = 1
< error
| 1193 HY000 Unknown system variable 'execute'
> exec ^
: string
| set global exists = 1
< error
| 1064 42000 …near 'exists = 1' at line 1
> exec ^
: string
| set global exit = 1
< error
| 1064 42000 …near 'exit = 1' at line 1
> exec ^
: string
| set global expansion = 1
< error
| 1193 HY000 Unknown system variable 'expansion'
> exec ^
: string
| set global expire = 1
< error
| 1193 HY000 Unknown system variable 'expire'
> exec ^
: string
| set global explain = 1
< error
| 1064 42000 …near 'explain = 1' at line 1
> exec ^
: string
| set global export = 1
< error
| 1193 HY000 Unknown system variable 'export'
> exec ^
: string
| set global extended = 1
< error
| 1193 HY000 Unknown system variable 'extended'
> exec ^
: string
| set global extent_size = 1
< error
| 1193 HY000 Unknown system variable 'extent_size'
> exec ^
: string
| set global external = 1
< error
| 1193 HY000 Unknown system variable 'external'
> exec ^
: string
| set global external_format = 1
< error
| 1193 HY000 Unknown system variable 'external_format'
> exec ^
: string
| set global extract = 1
< error
| 1193 HY000 Unknown system variable 'extract'
> exec ^
: string
| set global factor = 1
< error
| 1193 HY000 Unknown system variable 'factor'
> exec ^
: string
| set global failed_login_attempts = 1
< error
| 1193 HY000 Unknown system variable 'failed_login_attempts'
> exec ^
: string
| set global false = 1
< error
| 1064 42000 …near 'false = 1' at line 1
> exec ^
: string
| set global fast = 1
< error
| 1193 HY000 Unknown system variable 'fast'
> exec ^
: string
| set global faults = 1
< error
| 1193 HY000 Unknown system variable 'faults'
> exec ^
: string
| set global fetch = 1
< error
| 1064 42000 …near 'fetch = 1' at line 1
> exec ^
: string
| set global fields = 1
< error
| 1193 HY000 Unknown system variable 'fields'
> exec ^
: string
| set global file = 1
< error
| 1193 HY000 Unknown system variable 'file'
> exec ^
: string
| set global file_block_size = 1
< error
| 1193 HY000 Unknown system variable 'file_block_size'
> exec ^
: string
| set global file_format = 1
< error
| 1193 HY000 Unknown system variable 'file_format'
> exec ^
: string
| set global file_name = 1
< error
| 1193 HY000 Unknown system variable 'file_name'
> exec ^
: string
| set global file_pattern = 1
< error
| 1193 HY000 Unknown system variable 'file_pattern'
> exec ^
: string
| set global file_prefix = 1
< error
| 1193 HY000 Unknown system variable 'file_prefix'
> exec ^
: string
| set global files = 1
< error
| 1193 HY000 Unknown system variable 'files'
> exec ^
: string
| set global filter = 1
< error
| 1193 HY000 Unknown system variable 'filter'
> exec ^
: string
| set global finish = 1
< error
| 1193 HY000 Unknown system variable 'finish'
> exec ^
: string
| set global first = 1
< error
| 1193 HY000 Unknown system variable 'first'
> exec ^
: string
| set global first_value = 1
< error
| 1064 42000 …near 'first_value = 1' at line 1
> exec ^
: string
| set global fixed = 1
< error
| 1193 HY000 Unknown system variable 'fixed'
> exec ^
: string
| set global float = 1
< error
| 1064 42000 …near 'float = 1' at line 1
> exec ^
: string
| set global float4 = 1
< error
| 1064 42000 …near 'float4 = 1' at line 1
> exec ^
: string
| set global float8 = 1
< error
| 1064 42000 …near 'float8 = 1' at line 1
> exec ^
: string
| set global flush = 1
< ok 0 0
> exec ^
: string
| set global following = 1
< error
| 1193 HY000 Unknown system variable 'following'
> exec ^
: string
| set global follows = 1
< error
| 1193 HY000 Unknown system variable 'follows'
> exec ^
: string
| set global for = 1
< error
| 1064 42000 …near 'for = 1' at line 1
> exec ^
: string
| set global force = 1
< error
| 1064 42000 …near 'force = 1' at line 1
> exec ^
: string
| set global foreign = 1
< error
| 1064 42000 …near 'foreign = 1' at line 1
> exec ^
: string
| set global format = 1
< error
| 1193 HY000 Unknown system variable 'format'
> exec ^
: string
| set global found = 1
< error
| 1193 HY000 Unknown system variable 'found'
> exec ^
: string
| set global from = 1
< error
| 1064 42000 …near 'from = 1' at line 1
> exec ^
: string
| set global full = 1
< error
| 1193 HY000 Unknown system variable 'full'
> exec ^
: string
| set global fulltext = 1
< error
| 1064 42000 …near 'fulltext = 1' at line 1
> exec ^
: string
| set global function = 1
< error
| 1064 42000 …near 'function = 1' at line 1
> exec ^
: string
| set global general = 1
< error
| 1193 HY000 Unknown system variable 'general'
> exec ^
: string
| set global generate = 1
< error
| 1193 HY000 Unknown system variable 'generate'
> exec ^
: string
| set global generated = 1
< error
| 1064 42000 …near 'generated = 1' at line 1
> exec ^
: string
| set global geomcollection = 1
< error
| 1193 HY000 Unknown system variable 'geomcollection'
> exec ^
: string
| set global geometry = 1
< error
| 1193 HY000 Unknown system variable 'geometry'
> exec ^
: string
| set global geometrycollection = 1
< error
| 1193 HY000 Unknown system variable 'geometrycollection'
> exec ^
: string
| set global get = 1
< error
| 1064 42000 …near 'get = 1' at line 1
> exec ^
: string
| set global get_format = 1
< error
| 1193 HY000 Unknown system variable 'get_format'
> exec ^
: string
| set global get_master_public_key = 1
< error
| 1193 HY000 Unknown system variable 'get_master_public_key'
> exec ^
: string
| set global get_source_public_key = 1
< error
| 1193 HY000 Unknown system variable 'get_source_public_key'
> exec ^
: string
| set global global = 1
< error
| 1064 42000 …near 'global = 1' at line 1
> exec ^
: string
| set global grant = 1
< error
| 1064 42000 …near 'grant = 1' at line 1
> exec ^
: string
| set global grants = 1
< error
| 1193 HY000 Unknown system variable 'grants'
> exec ^
: string
| set global greatest = 1
< error
| 1193 HY000 Unknown system variable 'greatest'
> exec ^
: string
| set global group = 1
< error
| 1064 42000 …near 'group = 1' at line 1
> exec ^
: string
| set global group_concat = 1
< error
| 1193 HY000 Unknown system variable 'group_concat'
> exec ^
: string
| set global group_replication = 1
< error
| 1193 HY000 Unknown system variable 'group_replication'
> exec ^
: string
| set global grouping = 1
< error
| 1064 42000 …near 'grouping = 1' at line 1
> exec ^
: string
| set global groups = 1
< error
| 1064 42000 …near 'groups = 1' at line 1
> exec ^
: string
| set global gtid_only = 1
< error
| 1193 HY000 Unknown system variable 'gtid_only'
> exec ^
: string
| set global gtids = 1
< error
| 1193 HY000 Unknown system variable 'gtids'
> exec ^
: string
| set global guided = 1
< error
| 1193 HY000 Unknown system variable 'guided'
> exec ^
: string
| set global handler = 1
< error
| 1193 HY000 Unknown system variable 'handler'
> exec ^
: string
| set global hash = 1
< error
| 1193 HY000 Unknown system variable 'hash'
> exec ^
: string
| set global having = 1
< error
| 1064 42000 …near 'having = 1' at line 1
> exec ^
: string
| set global header = 1
< error
| 1193 HY000 Unknown system variable 'header'
> exec ^
: string
| set global help = 1
< error
| 1193 HY000 Unknown system variable 'help'
> exec ^
: string
| set global high_priority = 1
< error
| 1064 42000 …near 'high_priority = 1' at line 1
> exec ^
: string
| set global histogram = 1
< error
| 1193 HY000 Unknown system variable 'histogram'
> exec ^
: string
| set global history = 1
< error
| 1193 HY000 Unknown system variable 'history'
> exec ^
: string
| set global host = 1
< error
| 1193 HY000 Unknown system variable 'host'
> exec ^
: string
| set global hosts = 1
< error
| 1193 HY000 Unknown system variable 'hosts'
> exec ^
: string
| set global hour = 1
< error
| 1193 HY000 Unknown system variable 'hour'
> exec ^
: string
| set global hour_microsecond = 1
< error
| 1064 42000 …near 'hour_microsecond = 1' at line 1
> exec ^
: string
| set global hour_minute = 1
< error
| 1064 42000 …near 'hour_minute = 1' at line 1
> exec ^
: string
| set global hour_second = 1
< error
| 1064 42000 …near 'hour_second = 1' at line 1
> exec ^
: string
| set global identified = 1
< error
| 1193 HY000 Unknown system variable 'identified'
> exec ^
: string
| set global if = 1
< error
| 1064 42000 …near 'if = 1' at line 1
> exec ^
: string
| set global ignore = 1
< error
| 1064 42000 …near 'ignore = 1' at line 1
> exec ^
: string
| set global ignore_server_ids = 1
< error
| 1193 HY000 Unknown system variable 'ignore_server_ids'
> exec ^
: string
| set global import = 1
< error
| 1193 HY000 Unknown system variable 'import'
> exec ^
: string
| set global in = 1
< error
| 1064 42000 …near 'in = 1' at line 1
> exec ^
: string
| set global inactive = 1
< error
| 1193 HY000 Unknown system variable 'inactive'
> exec ^
: string
| set global index = 1
< error
| 1064 42000 …near 'index = 1' at line 1
> exec ^
: string
| set global indexes = 1
< error
| 1193 HY000 Unknown system variable 'indexes'
> exec ^
: string
| set global infile = 1
< error
| 1064 42000 …near 'infile = 1' at line 1
> exec ^
: string
| set global initial = 1
< error
| 1193 HY000 Unknown system variable 'initial'
> exec ^
: string
| set global initial_size = 1
< error
| 1193 HY000 Unknown system variable 'initial_size'
> exec ^
: string
| set global initiate = 1
< error
| 1193 HY000 Unknown system variable 'initiate'
> exec ^
: string
| set global inner = 1
< error
| 1064 42000 …near 'inner = 1' at line 1
> exec ^
: string
| set global inout = 1
< error
| 1064 42000 …near 'inout = 1' at line 1
> exec ^
: string
| set global insensitive = 1
< error
| 1064 42000 …near 'insensitive = 1' at line 1
> exec ^
: string
| set global insert = 1
< error
| 1064 42000 …near 'insert = 1' at line 1
> exec ^
: string
| set global insert_method = 1
< error
| 1193 HY000 Unknown system variable 'insert_method'
> exec ^
: string
| set global install = 1
< error
| 1193 HY000 Unknown system variable 'install'
> exec ^
: string
| set global instance = 1
< error
| 1193 HY000 Unknown system variable 'instance'
> exec ^
: string
| set global int = 1
< error
| 1064 42000 …near 'int = 1' at line 1
> exec ^
: string
| set global int1 = 1
< error
| 1064 42000 …near 'int1 = 1' at line 1
> exec ^
: string
| set global int2 = 1
< error
| 1064 42000 …near 'int2 = 1' at line 1
> exec ^
: string
| set global int3 = 1
< error
| 1064 42000 …near 'int3 = 1' at line 1
> exec ^
: string
| set global int4 = 1
< error
| 1064 42000 …near 'int4 = 1' at line 1
> exec ^
: string
| set global int8 = 1
< error
| 1064 42000 …near 'int8 = 1' at line 1
> exec ^
: string
| set global integer = 1
< error
| 1064 42000 …near 'integer = 1' at line 1
> exec ^
: string
| set global intersect = 1
< error
| 1064 42000 …near 'intersect = 1' at line 1
> exec ^
: string
| set global interval = 1
< error
| 1064 42000 …near 'interval = 1' at line 1
> exec ^
: string
| set global into = 1
< error
| 1064 42000 …near 'into = 1' at line 1
> exec ^
: string
| set global invisible = 1
< error
| 1193 HY000 Unknown system variable 'invisible'
> exec ^
: string
| set global invoker = 1
< error
| 1193 HY000 Unknown system variable 'invoker'
> exec ^
: string
| set global io = 1
< error
| 1193 HY000 Unknown system variable 'io'
> exec ^
: string
| set global io_after_gtids = 1
< error
| 1064 42000 …near 'io_after_gtids = 1' at line 1
> exec ^
: string
| set global io_before_gtids = 1
< error
| 1064 42000 …near 'io_before_gtids = 1' at line 1
> exec ^
: string
| set global io_thread = 1
< error
| 1193 HY000 Unknown system variable 'io_thread'
> exec ^
: string
| set global ipc = 1
< error
| 1193 HY000 Unknown system variable 'ipc'
> exec ^
: string
| set global is = 1
< error
| 1064 42000 …near 'is = 1' at line 1
> exec ^
: string
| set global isolation = 1
< error
| 1193 HY000 Unknown system variable 'isolation'
> exec ^
: string
| set global issuer = 1
< error
| 1193 HY000 Unknown system variable 'issuer'
> exec ^
: string
| set global iterate = 1
< error
| 1064 42000 …near 'iterate = 1' at line 1
> exec ^
: string
| set global join = 1
< error
| 1064 42000 …near 'join = 1' at line 1
> exec ^
: string
| set global json = 1
< error
| 1193 HY000 Unknown system variable 'json'
> exec ^
: string
| set global json_arrayagg = 1
< error
| 1193 HY000 Unknown system variable 'json_arrayagg'
> exec ^
: string
| set global json_duality_object = 1
< error
| 1193 HY000 Unknown system variable 'json_duality_object'
> exec ^
: string
| set global json_objectagg = 1
< error
| 1193 HY000 Unknown system variable 'json_objectagg'
> exec ^
: string
| set global json_table = 1
< error
| 1064 42000 …near 'json_table = 1' at line 1
> exec ^
: string
| set global json_value = 1
< error
| 1193 HY000 Unknown system variable 'json_value'
> exec ^
: string
| set global key = 1
< error
| 1064 42000 …near 'key = 1' at line 1
> exec ^
: string
| set global key_block_size = 1
< error
| 1193 HY000 Unknown system variable 'key_block_size'
> exec ^
: string
| set global keyring = 1
< error
| 1193 HY000 Unknown system variable 'keyring'
> exec ^
: string
| set global keys = 1
< error
| 1064 42000 …near 'keys = 1' at line 1
> exec ^
: string
| set global kill = 1
< error
| 1064 42000 …near 'kill = 1' at line 1
> exec ^
: string
| set global lag = 1
< error
| 1064 42000 …near 'lag = 1' at line 1
> exec ^
: string
| set global language = 1
< error
| 1193 HY000 Unknown system variable 'language'
> exec ^
: string
| set global last = 1
< error
| 1193 HY000 Unknown system variable 'last'
> exec ^
: string
| set global last_value = 1
< error
| 1064 42000 …near 'last_value = 1' at line 1
> exec ^
: string
| set global lateral = 1
< error
| 1064 42000 …near 'lateral = 1' at line 1
> exec ^
: string
| set global lead = 1
< error
| 1064 42000 …near 'lead = 1' at line 1
> exec ^
: string
| set global leading = 1
< error
| 1064 42000 …near 'leading = 1' at line 1
> exec ^
: string
| set global least = 1
< error
| 1193 HY000 Unknown system variable 'least'
> exec ^
: string
| set global leave = 1
< error
| 1064 42000 …near 'leave = 1' at line 1
> exec ^
: string
| set global leaves = 1
< error
| 1193 HY000 Unknown system variable 'leaves'
> exec ^
: string
| set global left = 1
< error
| 1064 42000 …near 'left = 1' at line 1
> exec ^
: string
| set global less = 1
< error
| 1193 HY000 Unknown system variable 'less'
> exec ^
: string
| set global level = 1
< error
| 1193 HY000 Unknown system variable 'level'
> exec ^
: string
| set global library = 1
< error
| 1193 HY000 Unknown system variable 'library'
> exec ^
: string
| set global like = 1
< error
| 1064 42000 …near 'like = 1' at line 1
> exec ^
: string
| set global limit = 1
< error
| 1064 42000 …near 'limit = 1' at line 1
> exec ^
: string
| set global linear = 1
< error
| 1064 42000 …near 'linear = 1' at line 1
> exec ^
: string
| set global lines = 1
< error
| 1064 42000 …near 'lines = 1' at line 1
> exec ^
: string
| set global linestring = 1
< error
| 1193 HY000 Unknown system variable 'linestring'
> exec ^
: string
| set global list = 1
< error
| 1193 HY000 Unknown system variable 'list'
> exec ^
: string
| set global load = 1
< error
| 1064 42000 …near 'load = 1' at line 1
> exec ^
: string
| set global local = 1
< error
| 1064 42000 …near 'local = 1' at line 1
> exec ^
: string
| set global localtime = 1
< error
| 1064 42000 …near 'localtime = 1' at line 1
> exec ^
: string
| set global localtimestamp = 1
< error
| 1064 42000 …near 'localtimestamp = 1' at line 1
> exec ^
: string
| set global lock = 1
< error
| 1064 42000 …near 'lock = 1' at line 1
> exec ^
: string
| set global locked = 1
< error
| 1193 HY000 Unknown system variable 'locked'
> exec ^
: string
| set global locks = 1
< error
| 1193 HY000 Unknown system variable 'locks'
> exec ^
: string
| set global log = 1
< error
| 1193 HY000 Unknown system variable 'log'
> exec ^
: string
| set global logfile = 1
< error
| 1193 HY000 Unknown system variable 'logfile'
> exec ^
: string
| set global logs = 1
< error
| 1193 HY000 Unknown system variable 'logs'
> exec ^
: string
| set global long = 1
< error
| 1064 42000 …near 'long = 1' at line 1
> exec ^
: string
| set global longblob = 1
< error
| 1064 42000 …near 'longblob = 1' at line 1
> exec ^
: string
| set global longtext = 1
< error
| 1064 42000 …near 'longtext = 1' at line 1
> exec ^
: string
| set global loop = 1
< error
| 1064 42000 …near 'loop = 1' at line 1
> exec ^
: string
| set global low_priority = 1
< error
| 1064 42000 …near 'low_priority = 1' at line 1
> exec ^
: string
| set global manifest = 1
< error
| 1193 HY000 Unknown system variable 'manifest'
> exec ^
: string
| set global manual = 1
< error
| 1193 HY000 Unknown system variable 'manual'
> exec ^
: string
| set global master = 1
< error
| 1193 HY000 Unknown system variable 'master'
> exec ^
: string
| set global master_auto_position = 1
< error
| 1193 HY000 Unknown system variable 'master_auto_position'
> exec ^
: string
| set global master_bind = 1
< error
| 1064 42000 …near 'master_bind = 1' at line 1
> exec ^
: string
| set global master_compression_algorithms = 1
< error
| 1193 HY000 Unknown system variable 'master_compression_algorithms'
> exec ^
: string
| set global master_connect_retry = 1
< error
| 1193 HY000 Unknown system variable 'master_connect_retry'
> exec ^
: string
| set global master_delay = 1
< error
| 1193 HY000 Unknown system variable 'master_delay'
> exec ^
: string
| set global master_heartbeat_period = 1
< error
| 1193 HY000 Unknown system variable 'master_heartbeat_period'
> exec ^
: string
| set global master_host = 1
< error
| 1193 HY000 Unknown system variable 'master_host'
> exec ^
: string
| set global master_log_file = 1
< error
| 1193 HY000 Unknown system variable 'master_log_file'
> exec ^
: string
| set global master_log_pos = 1
< error
| 1193 HY000 Unknown system variable 'master_log_pos'
> exec ^
: string
| set global master_password = 1
< error
| 1193 HY000 Unknown system variable 'master_password'
> exec ^
: string
| set global master_port = 1
< error
| 1193 HY000 Unknown system variable 'master_port'
> exec ^
: string
| set global master_public_key_path = 1
< error
| 1193 HY000 Unknown system variable 'master_public_key_path'
> exec ^
: string
| set global master_retry_count = 1
< error
| 1193 HY000 Unknown system variable 'master_retry_count'
> exec ^
: string
| set global master_ssl = 1
< error
| 1193 HY000 Unknown system variable 'master_ssl'
> exec ^
: string
| set global master_ssl_ca = 1
< error
| 1193 HY000 Unknown system variable 'master_ssl_ca'
> exec ^
: string
| set global master_ssl_capath = 1
< error
| 1193 HY000 Unknown system variable 'master_ssl_capath'
> exec ^
: string
| set global master_ssl_cert = 1
< error
| 1193 HY000 Unknown system variable 'master_ssl_cert'
> exec ^
: string
| set global master_ssl_cipher = 1
< error
| 1193 HY000 Unknown system variable 'master_ssl_cipher'
> exec ^
: string
| set global master_ssl_crl = 1
< error
| 1193 HY000 Unknown system variable 'master_ssl_crl'
> exec ^
: string
| set global master_ssl_crlpath = 1
< error
| 1193 HY000 Unknown system variable 'master_ssl_crlpath'
> exec ^
: string
| set global master_ssl_key = 1
< error
| 1193 HY000 Unknown system variable 'master_ssl_key'
> exec ^
: string
| set global master_ssl_verify_server_cert = 1
< error
| 1064 42000 …near 'master_ssl_verify_server_cert = 1' at line 1
> exec ^
: string
| set global master_tls_ciphersuites = 1
< error
| 1193 HY000 Unknown system variable 'master_tls_ciphersuites'
> exec ^
: string
| set global master_tls_version = 1
< error
| 1193 HY000 Unknown system variable 'master_tls_version'
> exec ^
: string
| set global master_user = 1
< error
| 1193 HY000 Unknown system variable 'master_user'
> exec ^
: string
| set global master_zstd_compression_level = 1
< error
| 1193 HY000 Unknown system variable 'master_zstd_compression_level'
> exec ^
: string
| set global match = 1
< error
| 1064 42000 …near 'match = 1' at line 1
> exec ^
: string
| set global materialized = 1
< error
| 1193 HY000 Unknown system variable 'materialized'
> exec ^
: string
| set global max = 1
< error
| 1193 HY000 Unknown system variable 'max'
> exec ^
: string
| set global max_connections_per_hour = 1
< error
| 1193 HY000 Unknown system variable 'max_connections_per_hour'
> exec ^
: string
| set global max_queries_per_hour = 1
< error
| 1193 HY000 Unknown system variable 'max_queries_per_hour'
> exec ^
: string
| set global max_rows = 1
< error
| 1193 HY000 Unknown system variable 'max_rows'
> exec ^
: string
| set global max_size = 1
< error
| 1193 HY000 Unknown system variable 'max_size'
> exec ^
: string
| set global max_updates_per_hour = 1
< error
| 1193 HY000 Unknown system variable 'max_updates_per_hour'
> exec ^
: string
| set global max_user_connections = 1
< ok 0 0
> exec ^
: string
| set global maxvalue = 1
< error
| 1064 42000 …near 'maxvalue = 1' at line 1
> exec ^
: string
| set global medium = 1
< error
| 1193 HY000 Unknown system variable 'medium'
> exec ^
: string
| set global mediumblob = 1
< error
| 1064 42000 …near 'mediumblob = 1' at line 1
> exec ^
: string
| set global mediumint = 1
< error
| 1064 42000 …near 'mediumint = 1' at line 1
> exec ^
: string
| set global mediumtext = 1
< error
| 1064 42000 …near 'mediumtext = 1' at line 1
> exec ^
: string
| set global member = 1
< error
| 1193 HY000 Unknown system variable 'member'
> exec ^
: string
| set global memory = 1
< error
| 1193 HY000 Unknown system variable 'memory'
> exec ^
: string
| set global merge = 1
< error
| 1193 HY000 Unknown system variable 'merge'
> exec ^
: string
| set global message_text = 1
< error
| 1193 HY000 Unknown system variable 'message_text'
> exec ^
: string
| set global microsecond = 1
< error
| 1193 HY000 Unknown system variable 'microsecond'
> exec ^
: string
| set global mid = 1
< error
| 1193 HY000 Unknown system variable 'mid'
> exec ^
: string
| set global middleint = 1
< error
| 1064 42000 …near 'middleint = 1' at line 1
> exec ^
: string
| set global migrate = 1
< error
| 1193 HY000 Unknown system variable 'migrate'
> exec ^
: string
| set global min = 1
< error
| 1193 HY000 Unknown system variable 'min'
> exec ^
: string
| set global min_rows = 1
< error
| 1193 HY000 Unknown system variable 'min_rows'
> exec ^
: string
| set global minute = 1
< error
| 1193 HY000 Unknown system variable 'minute'
> exec ^
: string
| set global minute_microsecond = 1
< error
| 1064 42000 …near 'minute_microsecond = 1' at line 1
> exec ^
: string
| set global minute_second = 1
< error
| 1064 42000 …near 'minute_second = 1' at line 1
> exec ^
: string
| set global mod = 1
< error
| 1064 42000 …near 'mod = 1' at line 1
> exec ^
: string
| set global mode = 1
< error
| 1193 HY000 Unknown system variable 'mode'
> exec ^
: string
| set global modifies = 1
< error
| 1064 42000 …near 'modifies = 1' at line 1
> exec ^
: string
| set global modify = 1
< error
| 1193 HY000 Unknown system variable 'modify'
> exec ^
: string
| set global month = 1
< error
| 1193 HY000 Unknown system variable 'month'
> exec ^
: string
| set global multilinestring = 1
< error
| 1193 HY000 Unknown system variable 'multilinestring'
> exec ^
: string
| set global multipoint = 1
< error
| 1193 HY000 Unknown system variable 'multipoint'
> exec ^
: string
| set global multipolygon = 1
< error
| 1193 HY000 Unknown system variable 'multipolygon'
> exec ^
: string
| set global mutex = 1
< error
| 1193 HY000 Unknown system variable 'mutex'
> exec ^
: string
| set global mysql_errno = 1
< error
| 1193 HY000 Unknown system variable 'mysql_errno'
> exec ^
: string
| set global name = 1
< error
| 1193 HY000 Unknown system variable 'name'
> exec ^
: string
| set global names = 1
< error
| 1193 HY000 Unknown system variable 'names'
> exec ^
: string
| set global national = 1
< error
| 1193 HY000 Unknown system variable 'national'
> exec ^
: string
| set global natural = 1
< error
| 1064 42000 …near 'natural = 1' at line 1
> exec ^
: string
| set global nchar = 1
< error
| 1193 HY000 Unknown system variable 'nchar'
> exec ^
: string
| set global ndb = 1
< error
| 1193 HY000 Unknown system variable 'ndb'
> exec ^
: string
| set global ndbcluster = 1
< error
| 1193 HY000 Unknown system variable 'ndbcluster'
> exec ^
: string
| set global nested = 1
< error
| 1193 HY000 Unknown system variable 'nested'
> exec ^
: string
| set global network_namespace = 1
< error
| 1193 HY000 Unknown system variable 'network_namespace'
> exec ^
: string
| set global never = 1
< error
| 1193 HY000 Unknown system variable 'never'
> exec ^
: string
| set global new = 1
< ok 0 0
> exec ^
: string
| set global next = 1
< error
| 1193 HY000 Unknown system variable 'next'
> exec ^
: string
| set global no = 1
< error
| 1193 HY000 Unknown system variable 'no'
> exec ^
: string
| set global no_wait = 1
< error
| 1193 HY000 Unknown system variable 'no_wait'
> exec ^
: string
| set global no_write_to_binlog = 1
< error
| 1064 42000 …near 'no_write_to_binlog = 1' at line 1
> exec ^
: string
| set global nodegroup = 1
< error
| 1193 HY000 Unknown system variable 'nodegroup'
> exec ^
: string
| set global none = 1
< error
| 1193 HY000 Unknown system variable 'none'
> exec ^
: string
| set global not = 1
< error
| 1064 42000 …near 'not = 1' at line 1
> exec ^
: string
| set global now = 1
< error
| 1193 HY000 Unknown system variable 'now'
> exec ^
: string
| set global nowait = 1
< error
| 1193 HY000 Unknown system variable 'nowait'
> exec ^
: string
| set global nth_value = 1
< error
| 1064 42000 …near 'nth_value = 1' at line 1
> exec ^
: string
| set global ntile = 1
< error
| 1064 42000 …near 'ntile = 1' at line 1
> exec ^
: string
| set global null = 1
< error
| 1064 42000 …near 'null = 1' at line 1
> exec ^
: string
| set global nullif = 1
< error
| 1193 HY000 Unknown system variable 'nullif'
> exec ^
: string
| set global nulls = 1
< error
| 1193 HY000 Unknown system variable 'nulls'
> exec ^
: string
| set global number = 1
< error
| 1193 HY000 Unknown system variable 'number'
> exec ^
: string
| set global numeric = 1
< error
| 1064 42000 …near 'numeric = 1' at line 1
> exec ^
: string
| set global nvarchar = 1
< error
| 1193 HY000 Unknown system variable 'nvarchar'
> exec ^
: string
| set global of = 1
< error
| 1064 42000 …near 'of = 1' at line 1
> exec ^
: string
| set global off = 1
< error
| 1193 HY000 Unknown system variable 'off'
> exec ^
: string
| set global offset = 1
< error
| 1193 HY000 Unknown system variable 'offset'
> exec ^
: string
| set global oj = 1
< error
| 1193 HY000 Unknown system variable 'oj'
> exec ^
: string
| set global old = 1
< ok 0 0
> exec ^
: string
| set global on = 1
< error
| 1064 42000 …near 'on = 1' at line 1
> exec ^
: string
| set global one = 1
< error
| 1193 HY000 Unknown system variable 'one'
> exec ^
: string
| set global only = 1
< error
| 1193 HY000 Unknown system variable 'only'
> exec ^
: string
| set global open = 1
< error
| 1193 HY000 Unknown system variable 'open'
> exec ^
: string
| set global optimize = 1
< error
| 1064 42000 …near 'optimize = 1' at line 1
> exec ^
: string
| set global optimizer_costs = 1
< error
| 1064 42000 …near 'optimizer_costs = 1' at line 1
> exec ^
: string
| set global option = 1
< error
| 1064 42000 …near 'option = 1' at line 1
> exec ^
: string
| set global optional = 1
< error
| 1193 HY000 Unknown system variable 'optional'
> exec ^
: string
| set global optionally = 1
< error
| 1064 42000 …near 'optionally = 1' at line 1
> exec ^
: string
| set global options = 1
< error
| 1193 HY000 Unknown system variable 'options'
> exec ^
: string
| set global or = 1
< error
| 1064 42000 …near 'or = 1' at line 1
> exec ^
: string
| set global order = 1
< error
| 1064 42000 …near 'order = 1' at line 1
> exec ^
: string
| set global ordinality = 1
< error
| 1193 HY000 Unknown system variable 'ordinality'
> exec ^
: string
| set global organization = 1
< error
| 1193 HY000 Unknown system variable 'organization'
> exec ^
: string
| set global others = 1
< error
| 1193 HY000 Unknown system variable 'others'
> exec ^
: string
| set global out = 1
< error
| 1064 42000 …near 'out = 1' at line 1
> exec ^
: string
| set global outer = 1
< error
| 1064 42000 …near 'outer = 1' at line 1
> exec ^
: string
| set global outfile = 1
< error
| 1064 42000 …near 'outfile = 1' at line 1
> exec ^
: string
| set global over = 1
< error
| 1064 42000 …near 'over = 1' at line 1
> exec ^
: string
| set global overlay = 1
< error
| 1193 HY000 Unknown system variable 'overlay'
> exec ^
: string
| set global overwrite = 1
< error
| 1193 HY000 Unknown system variable 'overwrite'
> exec ^
: string
| set global owner = 1
< error
| 1193 HY000 Unknown system variable 'owner'
> exec ^
: string
| set global pack_keys = 1
< error
| 1193 HY000 Unknown system variable 'pack_keys'
> exec ^
: string
| set global page = 1
< error
| 1193 HY000 Unknown system variable 'page'
> exec ^
: string
| set global parallel = 1
< error
| 1193 HY000 Unknown system variable 'parallel'
> exec ^
: string
| set global parameters = 1
< error
| 1193 HY000 Unknown system variable 'parameters'
> exec ^
: string
| set global parse_tree = 1
< error
| 1193 HY000 Unknown system variable 'parse_tree'
> exec ^
: string
| set global parser = 1
< error
| 1193 HY000 Unknown system variable 'parser'
> exec ^
: string
| set global partial = 1
< error
| 1193 HY000 Unknown system variable 'partial'
> exec ^
: string
| set global partition = 1
< error
| 1064 42000 …near 'partition = 1' at line 1
> exec ^
: string
| set global partitioning = 1
< error
| 1193 HY000 Unknown system variable 'partitioning'
> exec ^
: string
| set global partitions = 1
< error
| 1193 HY000 Unknown system variable 'partitions'
> exec ^
: string
| set global password = 1
< error
| 1193 HY000 Unknown system variable 'password'
> exec ^
: string
| set global password_lock_time = 1
< error
| 1193 HY000 Unknown system variable 'password_lock_time'
> exec ^
: string
| set global path = 1
< error
| 1193 HY000 Unknown system variable 'path'
> exec ^
: string
| set global percent_rank = 1
< error
| 1064 42000 …near 'percent_rank = 1' at line 1
> exec ^
: string
| set global persist = 1
< error
| 1064 42000 …near 'persist = 1' at line 1
> exec ^
: string
| set global persist_only = 1
< error
| 1064 42000 …near 'persist_only = 1' at line 1
> exec ^
: string
| set global phase = 1
< error
| 1193 HY000 Unknown system variable 'phase'
> exec ^
: string
| set global placing = 1
< error
| 1193 HY000 Unknown system variable 'placing'
> exec ^
: string
| set global plugin = 1
< error
| 1193 HY000 Unknown system variable 'plugin'
> exec ^
: string
| set global plugin_dir = 1
< ok 0 0
> exec ^
: string
| set global plugins = 1
< error
| 1193 HY000 Unknown system variable 'plugins'
> exec ^
: string
| set global point = 1
< error
| 1193 HY000 Unknown system variable 'point'
> exec ^
: string
| set global polygon = 1
< error
| 1193 HY000 Unknown system variable 'polygon'
> exec ^
: string
| set global port = 1
< ok 0 0
> exec ^
: string
| set global position = 1
< error
| 1193 HY000 Unknown system variable 'position'
> exec ^
: string
| set global precedes = 1
< error
| 1193 HY000 Unknown system variable 'precedes'
> exec ^
: string
| set global preceding = 1
< error
| 1193 HY000 Unknown system variable 'preceding'
> exec ^
: string
| set global precision = 1
< error
| 1064 42000 …near 'precision = 1' at line 1
> exec ^
: string
| set global prefix = 1
< error
| 1193 HY000 Unknown system variable 'prefix'
> exec ^
: string
| set global prepare = 1
< error
| 1193 HY000 Unknown system variable 'prepare'
> exec ^
: string
| set global preserve = 1
< error
| 1193 HY000 Unknown system variable 'preserve'
> exec ^
: string
| set global prev = 1
< error
| 1193 HY000 Unknown system variable 'prev'
> exec ^
: string
| set global primary = 1
< error
| 1064 42000 …near 'primary = 1' at line 1
> exec ^
: string
| set global privilege_checks_user = 1
< error
| 1193 HY000 Unknown system variable 'privilege_checks_user'
> exec ^
: string
| set global privileges = 1
< error
| 1193 HY000 Unknown system variable 'privileges'
> exec ^
: string
| set global procedure = 1
< error
| 1064 42000 …near 'procedure = 1' at line 1
> exec ^
: string
| set global process = 1
< error
| 1193 HY000 Unknown system variable 'process'
> exec ^
: string
| set global processlist = 1
< error
| 1193 HY000 Unknown system variable 'processlist'
> exec ^
: string
| set global profile = 1
< error
| 1193 HY000 Unknown system variable 'profile'
> exec ^
: string
| set global profiles = 1
< error
| 1193 HY000 Unknown system variable 'profiles'
> exec ^
: string
| set global proxy = 1
< error
| 1193 HY000 Unknown system variable 'proxy'
> exec ^
: string
| set global purge = 1
< error
| 1064 42000 …near 'purge = 1' at line 1
> exec ^
: string
| set global qualify = 1
< error
| 1193 HY000 Unknown system variable 'qualify'
> exec ^
: string
| set global quarter = 1
< error
| 1193 HY000 Unknown system variable 'quarter'
> exec ^
: string
| set global query = 1
< error
| 1193 HY000 Unknown system variable 'query'
> exec ^
: string
| set global quick = 1
< error
| 1193 HY000 Unknown system variable 'quick'
> exec ^
: string
| set global random = 1
< error
| 1193 HY000 Unknown system variable 'random'
> exec ^
: string
| set global range = 1
< error
| 1064 42000 …near 'range = 1' at line 1
> exec ^
: string
| set global rank = 1
< error
| 1064 42000 …near 'rank = 1' at line 1
> exec ^
: string
| set global read = 1
< error
| 1064 42000 …near 'read = 1' at line 1
> exec ^
: string
| set global read_only = 1
< ok 0 0
> exec ^
: string
| set global read_write = 1
< error
| 1064 42000 …near 'read_write = 1' at line 1
> exec ^
: string
| set global reads = 1
< error
| 1064 42000 …near 'reads = 1' at line 1
> exec ^
: string
| set global real = 1
< error
| 1064 42000 …near 'real = 1' at line 1
> exec ^
: string
| set global rebuild = 1
< error
| 1193 HY000 Unknown system variable 'rebuild'
> exec ^
: string
| set global recover = 1
< error
| 1193 HY000 Unknown system variable 'recover'
> exec ^
: string
| set global recursive = 1
< error
| 1064 42000 …near 'recursive = 1' at line 1
> exec ^
: string
| set global redo_buffer_size = 1
< error
| 1193 HY000 Unknown system variable 'redo_buffer_size'
> exec ^
: string
| set global redundant = 1
< error
| 1193 HY000 Unknown system variable 'redundant'
> exec ^
: string
| set global reference = 1
< error
| 1193 HY000 Unknown system variable 'reference'
> exec ^
: string
| set global references = 1
< error
| 1064 42000 …near 'references = 1' at line 1
> exec ^
: string
| set global regexp = 1
< error
| 1064 42000 …near 'regexp = 1' at line 1
> exec ^
: string
| set global registration = 1
< error
| 1193 HY000 Unknown system variable 'registration'
> exec ^
: string
| set global relational = 1
< error
| 1193 HY000 Unknown system variable 'relational'
> exec ^
: string
| set global relay = 1
< error
| 1193 HY000 Unknown system variable 'relay'
> exec ^
: string
| set global relay_log_file = 1
< error
| 1193 HY000 Unknown system variable 'relay_log_file'
> exec ^
: string
| set global relay_log_pos = 1
< error
| 1193 HY000 Unknown system variable 'relay_log_pos'
> exec ^
: string
| set global relay_thread = 1
< error
| 1193 HY000 Unknown system variable 'relay_thread'
> exec ^
: string
| set global relaylog = 1
< error
| 1193 HY000 Unknown system variable 'relaylog'
> exec ^
: string
| set global release = 1
< error
| 1064 42000 …near 'release = 1' at line 1
> exec ^
: string
| set global reload = 1
< error
| 1193 HY000 Unknown system variable 'reload'
> exec ^
: string
| set global remove = 1
< error
| 1193 HY000 Unknown system variable 'remove'
> exec ^
: string
| set global rename = 1
< error
| 1064 42000 …near 'rename = 1' at line 1
> exec ^
: string
| set global reorganize = 1
< error
| 1193 HY000 Unknown system variable 'reorganize'
> exec ^
: string
| set global repair = 1
< error
| 1193 HY000 Unknown system variable 'repair'
> exec ^
: string
| set global repeat = 1
< error
| 1064 42000 …near 'repeat = 1' at line 1
> exec ^
: string
| set global repeatable = 1
< error
| 1193 HY000 Unknown system variable 'repeatable'
> exec ^
: string
| set global replace = 1
< error
| 1064 42000 …near 'replace = 1' at line 1
> exec ^
: string
| set global replica = 1
< error
| 1193 HY000 Unknown system variable 'replica'
> exec ^
: string
| set global replicas = 1
< error
| 1193 HY000 Unknown system variable 'replicas'
> exec ^
: string
| set global replicate_do_db = 1
< error
| 1193 HY000 Unknown system variable 'replicate_do_db'
> exec ^
: string
| set global replicate_do_table = 1
< error
| 1193 HY000 Unknown system variable 'replicate_do_table'
> exec ^
: string
| set global replicate_ignore_db = 1
< error
| 1193 HY000 Unknown system variable 'replicate_ignore_db'
> exec ^
: string
| set global replicate_ignore_table = 1
< error
| 1193 HY000 Unknown system variable 'replicate_ignore_table'
> exec ^
: string
| set global replicate_rewrite_db = 1
< error
| 1193 HY000 Unknown system variable 'replicate_rewrite_db'
> exec ^
: string
| set global replicate_wild_do_table = 1
< error
| 1193 HY000 Unknown system variable 'replicate_wild_do_table'
> exec ^
: string
| set global replicate_wild_ignore_table = 1
< error
| 1193 HY000 Unknown system variable 'replicate_wild_ignore_table'
> exec ^
: string
| set global replication = 1
< error
| 1193 HY000 Unknown system variable 'replication'
> exec ^
: string
| set global require = 1
< error
| 1064 42000 …near 'require = 1' at line 1
> exec ^
: string
| set global require_row_format = 1
< error
| 1228 HY000 Variable 'require_row_format' is a SESSION variable and can't be used with SET GLOBAL
> exec ^
: string
| set global require_row_format = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set global require_table_primary_key_check = 1
< error
| 1193 HY000 Unknown system variable 'require_table_primary_key_check'
> exec ^
: string
| set global reset = 1
< error
| 1193 HY000 Unknown system variable 'reset'
> exec ^
: string
| set global resignal = 1
< error
| 1064 42000 …near 'resignal = 1' at line 1
> exec ^
: string
| set global resource = 1
< error
| 1193 HY000 Unknown system variable 'resource'
> exec ^
: string
| set global respect = 1
< error
| 1193 HY000 Unknown system variable 'respect'
> exec ^
: string
| set global restart = 1
< error
| 1193 HY000 Unknown system variable 'restart'
> exec ^
: string
| set global restore = 1
< error
| 1193 HY000 Unknown system variable 'restore'
> exec ^
: string
| set global restrict = 1
< error
| 1064 42000 …near 'restrict = 1' at line 1
> exec ^
: string
| set global resume = 1
< error
| 1193 HY000 Unknown system variable 'resume'
> exec ^
: string
| set global retain = 1
< error
| 1193 HY000 Unknown system variable 'retain'
> exec ^
: string
| set global return = 1
< error
| 1064 42000 …near 'return = 1' at line 1
> exec ^
: string
| set global returned_sqlstate = 1
< error
| 1193 HY000 Unknown system variable 'returned_sqlstate'
> exec ^
: string
| set global returning = 1
< error
| 1193 HY000 Unknown system variable 'returning'
> exec ^
: string
| set global returns = 1
< error
| 1193 HY000 Unknown system variable 'returns'
> exec ^
: string
| set global reuse = 1
< error
| 1193 HY000 Unknown system variable 'reuse'
> exec ^
: string
| set global reverse = 1
< error
| 1193 HY000 Unknown system variable 'reverse'
> exec ^
: string
| set global revoke = 1
< error
| 1064 42000 …near 'revoke = 1' at line 1
> exec ^
: string
| set global right = 1
< error
| 1064 42000 …near 'right = 1' at line 1
> exec ^
: string
| set global rlike = 1
< error
| 1064 42000 …near 'rlike = 1' at line 1
> exec ^
: string
| set global role = 1
< error
| 1193 HY000 Unknown system variable 'role'
> exec ^
: string
| set global rollback = 1
< error
| 1193 HY000 Unknown system variable 'rollback'
> exec ^
: string
| set global rollup = 1
< error
| 1193 HY000 Unknown system variable 'rollup'
> exec ^
: string
| set global rotate = 1
< error
| 1193 HY000 Unknown system variable 'rotate'
> exec ^
: string
| set global routine = 1
< error
| 1193 HY000 Unknown system variable 'routine'
> exec ^
: string
| set global row = 1
< error
| 1064 42000 …near 'row = 1' at line 1
> exec ^
: string
| set global row_count = 1
< error
| 1193 HY000 Unknown system variable 'row_count'
> exec ^
: string
| set global row_format = 1
< error
| 1193 HY000 Unknown system variable 'row_format'
> exec ^
: string
| set global row_number = 1
< error
| 1064 42000 …near 'row_number = 1' at line 1
> exec ^
: string
| set global rows = 1
< error
| 1064 42000 …near 'rows = 1' at line 1
> exec ^
: string
| set global rtree = 1
< error
| 1193 HY000 Unknown system variable 'rtree'
> exec ^
: string
| set global s3 = 1
< error
| 1193 HY000 Unknown system variable 's3'
> exec ^
: string
| set global savepoint = 1
< error
| 1193 HY000 Unknown system variable 'savepoint'
> exec ^
: string
| set global schedule = 1
< error
| 1193 HY000 Unknown system variable 'schedule'
> exec ^
: string
| set global schema = 1
< error
| 1064 42000 …near 'schema = 1' at line 1
> exec ^
: string
| set global schema_name = 1
< error
| 1193 HY000 Unknown system variable 'schema_name'
> exec ^
: string
| set global schemas = 1
< error
| 1064 42000 …near 'schemas = 1' at line 1
> exec ^
: string
| set global second = 1
< error
| 1193 HY000 Unknown system variable 'second'
> exec ^
: string
| set global second_microsecond = 1
< error
| 1064 42000 …near 'second_microsecond = 1' at line 1
> exec ^
: string
| set global secondary = 1
< error
| 1193 HY000 Unknown system variable 'secondary'
> exec ^
: string
| set global secondary_engine = 1
< error
| 1193 HY000 Unknown system variable 'secondary_engine'
> exec ^
: string
| set global secondary_engine_attribute = 1
< error
| 1193 HY000 Unknown system variable 'secondary_engine_attribute'
> exec ^
: string
| set global secondary_load = 1
< error
| 1193 HY000 Unknown system variable 'secondary_load'
> exec ^
: string
| set global secondary_unload = 1
< error
| 1193 HY000 Unknown system variable 'secondary_unload'
> exec ^
: string
| set global security = 1
< error
| 1193 HY000 Unknown system variable 'security'
> exec ^
: string
| set global select = 1
< error
| 1064 42000 …near 'select = 1' at line 1
> exec ^
: string
| set global sensitive = 1
< error
| 1064 42000 …near 'sensitive = 1' at line 1
> exec ^
: string
| set global separator = 1
< error
| 1064 42000 …near 'separator = 1' at line 1
> exec ^
: string
| set global serial = 1
< error
| 1193 HY000 Unknown system variable 'serial'
> exec ^
: string
| set global serializable = 1
< error
| 1193 HY000 Unknown system variable 'serializable'
> exec ^
: string
| set global server = 1
< error
| 1193 HY000 Unknown system variable 'server'
> exec ^
: string
| set global session = 1
< error
| 1064 42000 …near 'session = 1' at line 1
> exec ^
: string
| set global session_user = 1
< error
| 1193 HY000 Unknown system variable 'session_user'
> exec ^
: string
| set global set = 1
< error
| 1064 42000 …near 'set = 1' at line 1
> exec ^
: string
| set global sets = 1
< error
| 1193 HY000 Unknown system variable 'sets'
> exec ^
: string
| set global share = 1
< error
| 1193 HY000 Unknown system variable 'share'
> exec ^
: string
| set global show = 1
< error
| 1064 42000 …near 'show = 1' at line 1
> exec ^
: string
| set global shutdown = 1
< error
| 1193 HY000 Unknown system variable 'shutdown'
> exec ^
: string
| set global signal = 1
< error
| 1064 42000 …near 'signal = 1' at line 1
> exec ^
: string
| set global signed = 1
< error
| 1193 HY000 Unknown system variable 'signed'
> exec ^
: string
| set global simple = 1
< error
| 1193 HY000 Unknown system variable 'simple'
> exec ^
: string
| set global skip = 1
< error
| 1193 HY000 Unknown system variable 'skip'
> exec ^
: string
| set global slave = 1
< error
| 1193 HY000 Unknown system variable 'slave'
> exec ^
: string
| set global slow = 1
< error
| 1193 HY000 Unknown system variable 'slow'
> exec ^
: string
| set global smallint = 1
< error
| 1064 42000 …near 'smallint = 1' at line 1
> exec ^
: string
| set global snapshot = 1
< error
| 1193 HY000 Unknown system variable 'snapshot'
> exec ^
: string
| set global socket = 1
< ok 0 0
> exec ^
: string
| set global some = 1
< error
| 1193 HY000 Unknown system variable 'some'
> exec ^
: string
| set global soname = 1
< error
| 1193 HY000 Unknown system variable 'soname'
> exec ^
: string
| set global sounds = 1
< error
| 1193 HY000 Unknown system variable 'sounds'
> exec ^
: string
| set global source = 1
< error
| 1193 HY000 Unknown system variable 'source'
> exec ^
: string
| set global source_auto_position = 1
< error
| 1193 HY000 Unknown system variable 'source_auto_position'
> exec ^
: string
| set global source_bind = 1
< error
| 1193 HY000 Unknown system variable 'source_bind'
> exec ^
: string
| set global source_compression_algorithms = 1
< error
| 1193 HY000 Unknown system variable 'source_compression_algorithms'
> exec ^
: string
| set global source_connect_retry = 1
< error
| 1193 HY000 Unknown system variable 'source_connect_retry'
> exec ^
: string
| set global source_connection_auto_failover = 1
< error
| 1193 HY000 Unknown system variable 'source_connection_auto_failover'
> exec ^
: string
| set global source_delay = 1
< error
| 1193 HY000 Unknown system variable 'source_delay'
> exec ^
: string
| set global source_heartbeat_period = 1
< error
| 1193 HY000 Unknown system variable 'source_heartbeat_period'
> exec ^
: string
| set global source_host = 1
< error
| 1193 HY000 Unknown system variable 'source_host'
> exec ^
: string
| set global source_log_file = 1
< error
| 1193 HY000 Unknown system variable 'source_log_file'
> exec ^
: string
| set global source_log_pos = 1
< error
| 1193 HY000 Unknown system variable 'source_log_pos'
> exec ^
: string
| set global source_password = 1
< error
| 1193 HY000 Unknown system variable 'source_password'
> exec ^
: string
| set global source_port = 1
< error
| 1193 HY000 Unknown system variable 'source_port'
> exec ^
: string
| set global source_public_key_path = 1
< error
| 1193 HY000 Unknown system variable 'source_public_key_path'
> exec ^
: string
| set global source_retry_count = 1
< error
| 1193 HY000 Unknown system variable 'source_retry_count'
> exec ^
: string
| set global source_ssl = 1
< error
| 1193 HY000 Unknown system variable 'source_ssl'
> exec ^
: string
| set global source_ssl_ca = 1
< error
| 1193 HY000 Unknown system variable 'source_ssl_ca'
> exec ^
: string
| set global source_ssl_capath = 1
< error
| 1193 HY000 Unknown system variable 'source_ssl_capath'
> exec ^
: string
| set global source_ssl_cert = 1
< error
| 1193 HY000 Unknown system variable 'source_ssl_cert'
> exec ^
: string
| set global source_ssl_cipher = 1
< error
| 1193 HY000 Unknown system variable 'source_ssl_cipher'
> exec ^
: string
| set global source_ssl_crl = 1
< error
| 1193 HY000 Unknown system variable 'source_ssl_crl'
> exec ^
: string
| set global source_ssl_crlpath = 1
< error
| 1193 HY000 Unknown system variable 'source_ssl_crlpath'
> exec ^
: string
| set global source_ssl_key = 1
< error
| 1193 HY000 Unknown system variable 'source_ssl_key'
> exec ^
: string
| set global source_ssl_verify_server_cert = 1
< error
| 1193 HY000 Unknown system variable 'source_ssl_verify_server_cert'
> exec ^
: string
| set global source_tls_ciphersuites = 1
< error
| 1193 HY000 Unknown system variable 'source_tls_ciphersuites'
> exec ^
: string
| set global source_tls_version = 1
< error
| 1193 HY000 Unknown system variable 'source_tls_version'
> exec ^
: string
| set global source_user = 1
< error
| 1193 HY000 Unknown system variable 'source_user'
> exec ^
: string
| set global source_zstd_compression_level = 1
< error
| 1193 HY000 Unknown system variable 'source_zstd_compression_level'
> exec ^
: string
| set global spatial = 1
< error
| 1064 42000 …near 'spatial = 1' at line 1
> exec ^
: string
| set global specific = 1
< error
| 1064 42000 …near 'specific = 1' at line 1
> exec ^
: string
| set global sql = 1
< error
| 1064 42000 …near 'sql = 1' at line 1
> exec ^
: string
| set global sql_after_gtids = 1
< error
| 1193 HY000 Unknown system variable 'sql_after_gtids'
> exec ^
: string
| set global sql_after_mts_gaps = 1
< error
| 1193 HY000 Unknown system variable 'sql_after_mts_gaps'
> exec ^
: string
| set global sql_before_gtids = 1
< error
| 1193 HY000 Unknown system variable 'sql_before_gtids'
> exec ^
: string
| set global sql_big_result = 1
< error
| 1064 42000 …near 'sql_big_result = 1' at line 1
> exec ^
: string
| set global sql_buffer_result = 1
< ok 0 0
> exec ^
: string
| set global sql_calc_found_rows = 1
< error
| 1064 42000 …near 'sql_calc_found_rows = 1' at line 1
> exec ^
: string
| set global sql_no_cache = 1
< error
| 1193 HY000 Unknown system variable 'sql_no_cache'
> exec ^
: string
| set global sql_small_result = 1
< error
| 1064 42000 …near 'sql_small_result = 1' at line 1
> exec ^
: string
| set global sql_thread = 1
< error
| 1193 HY000 Unknown system variable 'sql_thread'
> exec ^
: string
| set global sql_tsi_day = 1
< error
| 1193 HY000 Unknown system variable 'sql_tsi_day'
> exec ^
: string
| set global sql_tsi_hour = 1
< error
| 1193 HY000 Unknown system variable 'sql_tsi_hour'
> exec ^
: string
| set global sql_tsi_minute = 1
< error
| 1193 HY000 Unknown system variable 'sql_tsi_minute'
> exec ^
: string
| set global sql_tsi_month = 1
< error
| 1193 HY000 Unknown system variable 'sql_tsi_month'
> exec ^
: string
| set global sql_tsi_quarter = 1
< error
| 1193 HY000 Unknown system variable 'sql_tsi_quarter'
> exec ^
: string
| set global sql_tsi_second = 1
< error
| 1193 HY000 Unknown system variable 'sql_tsi_second'
> exec ^
: string
| set global sql_tsi_week = 1
< error
| 1193 HY000 Unknown system variable 'sql_tsi_week'
> exec ^
: string
| set global sql_tsi_year = 1
< error
| 1193 HY000 Unknown system variable 'sql_tsi_year'
> exec ^
: string
| set global sqlexception = 1
< error
| 1064 42000 …near 'sqlexception = 1' at line 1
> exec ^
: string
| set global sqlstate = 1
< error
| 1064 42000 …near 'sqlstate = 1' at line 1
> exec ^
: string
| set global sqlwarning = 1
< error
| 1064 42000 …near 'sqlwarning = 1' at line 1
> exec ^
: string
| set global srid = 1
< error
| 1193 HY000 Unknown system variable 'srid'
> exec ^
: string
| set global sse_kms = 1
< error
| 1193 HY000 Unknown system variable 'sse_kms'
> exec ^
: string
| set global sse_s3 = 1
< error
| 1193 HY000 Unknown system variable 'sse_s3'
> exec ^
: string
| set global ssl = 1
< error
| 1064 42000 …near 'ssl = 1' at line 1
> exec ^
: string
| set global st_collect = 1
< error
| 1193 HY000 Unknown system variable 'st_collect'
> exec ^
: string
| set global stacked = 1
< error
| 1193 HY000 Unknown system variable 'stacked'
> exec ^
: string
| set global start = 1
< error
| 1193 HY000 Unknown system variable 'start'
> exec ^
: string
| set global starting = 1
< error
| 1064 42000 …near 'starting = 1' at line 1
> exec ^
: string
| set global starts = 1
< error
| 1193 HY000 Unknown system variable 'starts'
> exec ^
: string
| set global stats_auto_recalc = 1
< error
| 1193 HY000 Unknown system variable 'stats_auto_recalc'
> exec ^
: string
| set global stats_persistent = 1
< error
| 1193 HY000 Unknown system variable 'stats_persistent'
> exec ^
: string
| set global stats_sample_pages = 1
< error
| 1193 HY000 Unknown system variable 'stats_sample_pages'
> exec ^
: string
| set global status = 1
< error
| 1193 HY000 Unknown system variable 'status'
> exec ^
: string
| set global std = 1
< error
| 1193 HY000 Unknown system variable 'std'
> exec ^
: string
| set global stddev = 1
< error
| 1193 HY000 Unknown system variable 'stddev'
> exec ^
: string
| set global stddev_pop = 1
< error
| 1193 HY000 Unknown system variable 'stddev_pop'
> exec ^
: string
| set global stddev_samp = 1
< error
| 1193 HY000 Unknown system variable 'stddev_samp'
> exec ^
: string
| set global stop = 1
< error
| 1193 HY000 Unknown system variable 'stop'
> exec ^
: string
| set global storage = 1
< error
| 1193 HY000 Unknown system variable 'storage'
> exec ^
: string
| set global stored = 1
< error
| 1064 42000 …near 'stored = 1' at line 1
> exec ^
: string
| set global straight_join = 1
< error
| 1064 42000 …near 'straight_join = 1' at line 1
> exec ^
: string
| set global stream = 1
< error
| 1193 HY000 Unknown system variable 'stream'
> exec ^
: string
| set global strict_load = 1
< error
| 1193 HY000 Unknown system variable 'strict_load'
> exec ^
: string
| set global string = 1
< error
| 1193 HY000 Unknown system variable 'string'
> exec ^
: string
| set global subclass_origin = 1
< error
| 1193 HY000 Unknown system variable 'subclass_origin'
> exec ^
: string
| set global subdate = 1
< error
| 1193 HY000 Unknown system variable 'subdate'
> exec ^
: string
| set global subject = 1
< error
| 1193 HY000 Unknown system variable 'subject'
> exec ^
: string
| set global subpartition = 1
< error
| 1193 HY000 Unknown system variable 'subpartition'
> exec ^
: string
| set global subpartitions = 1
< error
| 1193 HY000 Unknown system variable 'subpartitions'
> exec ^
: string
| set global substr = 1
< error
| 1193 HY000 Unknown system variable 'substr'
> exec ^
: string
| set global substring = 1
< error
| 1193 HY000 Unknown system variable 'substring'
> exec ^
: string
| set global sum = 1
< error
| 1193 HY000 Unknown system variable 'sum'
> exec ^
: string
| set global super = 1
< error
| 1193 HY000 Unknown system variable 'super'
> exec ^
: string
| set global suspend = 1
< error
| 1193 HY000 Unknown system variable 'suspend'
> exec ^
: string
| set global swaps = 1
< error
| 1193 HY000 Unknown system variable 'swaps'
> exec ^
: string
| set global switches = 1
< error
| 1193 HY000 Unknown system variable 'switches'
> exec ^
: string
| set global sysdate = 1
< error
| 1193 HY000 Unknown system variable 'sysdate'
> exec ^
: string
| set global system = 1
< error
| 1064 42000 …near 'system = 1' at line 1
> exec ^
: string
| set global system_user = 1
< error
| 1193 HY000 Unknown system variable 'system_user'
> exec ^
: string
| set global table = 1
< error
| 1064 42000 …near 'table = 1' at line 1
> exec ^
: string
| set global table_checksum = 1
< error
| 1193 HY000 Unknown system variable 'table_checksum'
> exec ^
: string
| set global table_name = 1
< error
| 1193 HY000 Unknown system variable 'table_name'
> exec ^
: string
| set global tables = 1
< error
| 1193 HY000 Unknown system variable 'tables'
> exec ^
: string
| set global tablesample = 1
< error
| 1193 HY000 Unknown system variable 'tablesample'
> exec ^
: string
| set global tablespace = 1
< error
| 1193 HY000 Unknown system variable 'tablespace'
> exec ^
: string
| set global temporary = 1
< error
| 1193 HY000 Unknown system variable 'temporary'
> exec ^
: string
| set global temptable = 1
< error
| 1193 HY000 Unknown system variable 'temptable'
> exec ^
: string
| set global terminated = 1
< error
| 1064 42000 …near 'terminated = 1' at line 1
> exec ^
: string
| set global text = 1
< error
| 1193 HY000 Unknown system variable 'text'
> exec ^
: string
| set global than = 1
< error
| 1193 HY000 Unknown system variable 'than'
> exec ^
: string
| set global then = 1
< error
| 1064 42000 …near 'then = 1' at line 1
> exec ^
: string
| set global thread_priority = 1
< error
| 1193 HY000 Unknown system variable 'thread_priority'
> exec ^
: string
| set global ties = 1
< error
| 1193 HY000 Unknown system variable 'ties'
> exec ^
: string
| set global time = 1
< error
| 1193 HY000 Unknown system variable 'time'
> exec ^
: string
| set global timestamp = 1
< error
| 1228 HY000 Variable 'timestamp' is a SESSION variable and can't be used with SET GLOBAL
> exec ^
: string
| set global timestamp = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set global timestampadd = 1
< error
| 1193 HY000 Unknown system variable 'timestampadd'
> exec ^
: string
| set global timestampdiff = 1
< error
| 1193 HY000 Unknown system variable 'timestampdiff'
> exec ^
: string
| set global tinyblob = 1
< error
| 1064 42000 …near 'tinyblob = 1' at line 1
> exec ^
: string
| set global tinyint = 1
< error
| 1064 42000 …near 'tinyint = 1' at line 1
> exec ^
: string
| set global tinytext = 1
< error
| 1064 42000 …near 'tinytext = 1' at line 1
> exec ^
: string
| set global tls = 1
< error
| 1193 HY000 Unknown system variable 'tls'
> exec ^
: string
| set global to = 1
< error
| 1064 42000 …near 'to = 1' at line 1
> exec ^
: string
| set global trailing = 1
< error
| 1064 42000 …near 'trailing = 1' at line 1
> exec ^
: string
| set global transaction = 1
< error
| 1193 HY000 Unknown system variable 'transaction'
> exec ^
: string
| set global trigger = 1
< error
| 1064 42000 …near 'trigger = 1' at line 1
> exec ^
: string
| set global triggers = 1
< error
| 1193 HY000 Unknown system variable 'triggers'
> exec ^
: string
| set global trim = 1
< error
| 1193 HY000 Unknown system variable 'trim'
> exec ^
: string
| set global true = 1
< error
| 1064 42000 …near 'true = 1' at line 1
> exec ^
: string
| set global truncate = 1
< error
| 1193 HY000 Unknown system variable 'truncate'
> exec ^
: string
| set global type = 1
< error
| 1193 HY000 Unknown system variable 'type'
> exec ^
: string
| set global types = 1
< error
| 1193 HY000 Unknown system variable 'types'
> exec ^
: string
| set global unbounded = 1
< error
| 1193 HY000 Unknown system variable 'unbounded'
> exec ^
: string
| set global uncommitted = 1
< error
| 1193 HY000 Unknown system variable 'uncommitted'
> exec ^
: string
| set global undefined = 1
< error
| 1193 HY000 Unknown system variable 'undefined'
> exec ^
: string
| set global undo = 1
< error
| 1064 42000 …near 'undo = 1' at line 1
> exec ^
: string
| set global undo_buffer_size = 1
< error
| 1193 HY000 Unknown system variable 'undo_buffer_size'
> exec ^
: string
| set global undofile = 1
< error
| 1193 HY000 Unknown system variable 'undofile'
> exec ^
: string
| set global unicode = 1
< error
| 1193 HY000 Unknown system variable 'unicode'
> exec ^
: string
| set global uninstall = 1
< error
| 1193 HY000 Unknown system variable 'uninstall'
> exec ^
: string
| set global union = 1
< error
| 1064 42000 …near 'union = 1' at line 1
> exec ^
: string
| set global unique = 1
< error
| 1064 42000 …near 'unique = 1' at line 1
> exec ^
: string
| set global unknown = 1
< error
| 1193 HY000 Unknown system variable 'unknown'
> exec ^
: string
| set global unlock = 1
< error
| 1064 42000 …near 'unlock = 1' at line 1
> exec ^
: string
| set global unregister = 1
< error
| 1193 HY000 Unknown system variable 'unregister'
> exec ^
: string
| set global unsigned = 1
< error
| 1064 42000 …near 'unsigned = 1' at line 1
> exec ^
: string
| set global until = 1
< error
| 1193 HY000 Unknown system variable 'until'
> exec ^
: string
| set global update = 1
< error
| 1064 42000 …near 'update = 1' at line 1
> exec ^
: string
| set global upgrade = 1
< error
| 1193 HY000 Unknown system variable 'upgrade'
> exec ^
: string
| set global uri = 1
< error
| 1193 HY000 Unknown system variable 'uri'
> exec ^
: string
| set global url = 1
< error
| 1193 HY000 Unknown system variable 'url'
> exec ^
: string
| set global usage = 1
< error
| 1064 42000 …near 'usage = 1' at line 1
> exec ^
: string
| set global use = 1
< error
| 1064 42000 …near 'use = 1' at line 1
> exec ^
: string
| set global use_frm = 1
< error
| 1193 HY000 Unknown system variable 'use_frm'
> exec ^
: string
| set global user = 1
< error
| 1193 HY000 Unknown system variable 'user'
> exec ^
: string
| set global user_resources = 1
< error
| 1193 HY000 Unknown system variable 'user_resources'
> exec ^
: string
| set global using = 1
< error
| 1064 42000 …near 'using = 1' at line 1
> exec ^
: string
| set global utc_date = 1
< error
| 1064 42000 …near 'utc_date = 1' at line 1
> exec ^
: string
| set global utc_time = 1
< error
| 1064 42000 …near 'utc_time = 1' at line 1
> exec ^
: string
| set global utc_timestamp = 1
< error
| 1064 42000 …near 'utc_timestamp = 1' at line 1
> exec ^
: string
| set global validate = 1
< error
| 1193 HY000 Unknown system variable 'validate'
> exec ^
: string
| set global validation = 1
< error
| 1193 HY000 Unknown system variable 'validation'
> exec ^
: string
| set global value = 1
< error
| 1193 HY000 Unknown system variable 'value'
> exec ^
: string
| set global values = 1
< error
| 1064 42000 …near 'values = 1' at line 1
> exec ^
: string
| set global var_pop = 1
< error
| 1193 HY000 Unknown system variable 'var_pop'
> exec ^
: string
| set global var_samp = 1
< error
| 1193 HY000 Unknown system variable 'var_samp'
> exec ^
: string
| set global varbinary = 1
< error
| 1064 42000 …near 'varbinary = 1' at line 1
> exec ^
: string
| set global varchar = 1
< error
| 1064 42000 …near 'varchar = 1' at line 1
> exec ^
: string
| set global varcharacter = 1
< error
| 1064 42000 …near 'varcharacter = 1' at line 1
> exec ^
: string
| set global variables = 1
< error
| 1193 HY000 Unknown system variable 'variables'
> exec ^
: string
| set global variance = 1
< error
| 1193 HY000 Unknown system variable 'variance'
> exec ^
: string
| set global varying = 1
< error
| 1064 42000 …near 'varying = 1' at line 1
> exec ^
: string
| set global vcpu = 1
< error
| 1193 HY000 Unknown system variable 'vcpu'
> exec ^
: string
| set global vector = 1
< error
| 1193 HY000 Unknown system variable 'vector'
> exec ^
: string
| set global verify_key_constraints = 1
< error
| 1193 HY000 Unknown system variable 'verify_key_constraints'
> exec ^
: string
| set global view = 1
< error
| 1193 HY000 Unknown system variable 'view'
> exec ^
: string
| set global virtual = 1
< error
| 1064 42000 …near 'virtual = 1' at line 1
> exec ^
: string
| set global visible = 1
< error
| 1193 HY000 Unknown system variable 'visible'
> exec ^
: string
| set global wait = 1
< error
| 1193 HY000 Unknown system variable 'wait'
> exec ^
: string
| set global warnings = 1
< error
| 1193 HY000 Unknown system variable 'warnings'
> exec ^
: string
| set global week = 1
< error
| 1193 HY000 Unknown system variable 'week'
> exec ^
: string
| set global weight_string = 1
< error
| 1193 HY000 Unknown system variable 'weight_string'
> exec ^
: string
| set global when = 1
< error
| 1064 42000 …near 'when = 1' at line 1
> exec ^
: string
| set global where = 1
< error
| 1064 42000 …near 'where = 1' at line 1
> exec ^
: string
| set global while = 1
< error
| 1064 42000 …near 'while = 1' at line 1
> exec ^
: string
| set global window = 1
< error
| 1064 42000 …near 'window = 1' at line 1
> exec ^
: string
| set global with = 1
< error
| 1064 42000 …near 'with = 1' at line 1
> exec ^
: string
| set global without = 1
< error
| 1193 HY000 Unknown system variable 'without'
> exec ^
: string
| set global work = 1
< error
| 1193 HY000 Unknown system variable 'work'
> exec ^
: string
| set global wrapper = 1
< error
| 1193 HY000 Unknown system variable 'wrapper'
> exec ^
: string
| set global write = 1
< error
| 1064 42000 …near 'write = 1' at line 1
> exec ^
: string
| set global x509 = 1
< error
| 1193 HY000 Unknown system variable 'x509'
> exec ^
: string
| set global xa = 1
< error
| 1193 HY000 Unknown system variable 'xa'
> exec ^
: string
| set global xid = 1
< error
| 1193 HY000 Unknown system variable 'xid'
> exec ^
: string
| set global xml = 1
< error
| 1193 HY000 Unknown system variable 'xml'
> exec ^
: string
| set global xor = 1
< error
| 1064 42000 …near 'xor = 1' at line 1
> exec ^
: string
| set global year = 1
< error
| 1193 HY000 Unknown system variable 'year'
> exec ^
: string
| set global year_month = 1
< error
| 1064 42000 …near 'year_month = 1' at line 1
> exec ^
: string
| set global zerofill = 1
< error
| 1064 42000 …near 'zerofill = 1' at line 1
> exec ^
: string
| set global zone = 1
< error
| 1193 HY000 Unknown system variable 'zone'
> exec ^
: string
| set absent.a = 1
< error
| 1193 HY000 Unknown system variable 'absent.a'
> exec ^
: string
| set accessible.a = 1
< error
| 1193 HY000 Unknown system variable 'accessible.a'
> exec ^
: string
| set account.a = 1
< error
| 1193 HY000 Unknown system variable 'account.a'
> exec ^
: string
| set action.a = 1
< error
| 1193 HY000 Unknown system variable 'action.a'
> exec ^
: string
| set active.a = 1
< error
| 1193 HY000 Unknown system variable 'active.a'
> exec ^
: string
| set add.a = 1
< error
| 1193 HY000 Unknown system variable 'add.a'
> exec ^
: string
| set adddate.a = 1
< error
| 1193 HY000 Unknown system variable 'adddate.a'
> exec ^
: string
| set admin.a = 1
< error
| 1193 HY000 Unknown system variable 'admin.a'
> exec ^
: string
| set after.a = 1
< error
| 1193 HY000 Unknown system variable 'after.a'
> exec ^
: string
| set against.a = 1
< error
| 1193 HY000 Unknown system variable 'against.a'
> exec ^
: string
| set aggregate.a = 1
< error
| 1193 HY000 Unknown system variable 'aggregate.a'
> exec ^
: string
| set algorithm.a = 1
< error
| 1193 HY000 Unknown system variable 'algorithm.a'
> exec ^
: string
| set all.a = 1
< error
| 1193 HY000 Unknown system variable 'all.a'
> exec ^
: string
| set allow_missing_files.a = 1
< error
| 1193 HY000 Unknown system variable 'allow_missing_files.a'
> exec ^
: string
| set alter.a = 1
< error
| 1193 HY000 Unknown system variable 'alter.a'
> exec ^
: string
| set always.a = 1
< error
| 1193 HY000 Unknown system variable 'always.a'
> exec ^
: string
| set analyze.a = 1
< error
| 1193 HY000 Unknown system variable 'analyze.a'
> exec ^
: string
| set and.a = 1
< error
| 1193 HY000 Unknown system variable 'and.a'
> exec ^
: string
| set any.a = 1
< error
| 1193 HY000 Unknown system variable 'any.a'
> exec ^
: string
| set array.a = 1
< error
| 1193 HY000 Unknown system variable 'array.a'
> exec ^
: string
| set as.a = 1
< error
| 1193 HY000 Unknown system variable 'as.a'
> exec ^
: string
| set asc.a = 1
< error
| 1193 HY000 Unknown system variable 'asc.a'
> exec ^
: string
| set ascii.a = 1
< error
| 1193 HY000 Unknown system variable 'ascii.a'
> exec ^
: string
| set asensitive.a = 1
< error
| 1193 HY000 Unknown system variable 'asensitive.a'
> exec ^
: string
| set assign_gtids_to_anonymous_transactions.a = 1
< error
| 1193 HY000 Unknown system variable 'assign_gtids_to_anonymous_transactions.a'
> exec ^
: string
| set at.a = 1
< error
| 1193 HY000 Unknown system variable 'at.a'
> exec ^
: string
| set attribute.a = 1
< error
| 1193 HY000 Unknown system variable 'attribute.a'
> exec ^
: string
| set authentication.a = 1
< error
| 1193 HY000 Unknown system variable 'authentication.a'
> exec ^
: string
| set auto.a = 1
< error
| 1193 HY000 Unknown system variable 'auto.a'
> exec ^
: string
| set auto_increment.a = 1
< error
| 1193 HY000 Unknown system variable 'auto_increment.a'
> exec ^
: string
| set auto_refresh.a = 1
< error
| 1193 HY000 Unknown system variable 'auto_refresh.a'
> exec ^
: string
| set auto_refresh_source.a = 1
< error
| 1193 HY000 Unknown system variable 'auto_refresh_source.a'
> exec ^
: string
| set autoextend_size.a = 1
< error
| 1193 HY000 Unknown system variable 'autoextend_size.a'
> exec ^
: string
| set avg.a = 1
< error
| 1193 HY000 Unknown system variable 'avg.a'
> exec ^
: string
| set avg_row_length.a = 1
< error
| 1193 HY000 Unknown system variable 'avg_row_length.a'
> exec ^
: string
| set backup.a = 1
< error
| 1193 HY000 Unknown system variable 'backup.a'
> exec ^
: string
| set before.a = 1
< error
| 1193 HY000 Unknown system variable 'before.a'
> exec ^
: string
| set begin.a = 1
< error
| 1193 HY000 Unknown system variable 'begin.a'
> exec ^
: string
| set bernoulli.a = 1
< error
| 1193 HY000 Unknown system variable 'bernoulli.a'
> exec ^
: string
| set between.a = 1
< error
| 1193 HY000 Unknown system variable 'between.a'
> exec ^
: string
| set bigint.a = 1
< error
| 1193 HY000 Unknown system variable 'bigint.a'
> exec ^
: string
| set binary.a = 1
< error
| 1193 HY000 Unknown system variable 'binary.a'
> exec ^
: string
| set binlog.a = 1
< error
| 1193 HY000 Unknown system variable 'binlog.a'
> exec ^
: string
| set bit.a = 1
< error
| 1193 HY000 Unknown system variable 'bit.a'
> exec ^
: string
| set bit_and.a = 1
< error
| 1193 HY000 Unknown system variable 'bit_and.a'
> exec ^
: string
| set bit_or.a = 1
< error
| 1193 HY000 Unknown system variable 'bit_or.a'
> exec ^
: string
| set bit_xor.a = 1
< error
| 1193 HY000 Unknown system variable 'bit_xor.a'
> exec ^
: string
| set blob.a = 1
< error
| 1193 HY000 Unknown system variable 'blob.a'
> exec ^
: string
| set block.a = 1
< error
| 1193 HY000 Unknown system variable 'block.a'
> exec ^
: string
| set bool.a = 1
< error
| 1193 HY000 Unknown system variable 'bool.a'
> exec ^
: string
| set boolean.a = 1
< error
| 1193 HY000 Unknown system variable 'boolean.a'
> exec ^
: string
| set both.a = 1
< error
| 1193 HY000 Unknown system variable 'both.a'
> exec ^
: string
| set btree.a = 1
< error
| 1193 HY000 Unknown system variable 'btree.a'
> exec ^
: string
| set buckets.a = 1
< error
| 1193 HY000 Unknown system variable 'buckets.a'
> exec ^
: string
| set bulk.a = 1
< error
| 1193 HY000 Unknown system variable 'bulk.a'
> exec ^
: string
| set by.a = 1
< error
| 1193 HY000 Unknown system variable 'by.a'
> exec ^
: string
| set byte.a = 1
< error
| 1193 HY000 Unknown system variable 'byte.a'
> exec ^
: string
| set cache.a = 1
< error
| 1193 HY000 Unknown system variable 'cache.a'
> exec ^
: string
| set call.a = 1
< error
| 1193 HY000 Unknown system variable 'call.a'
> exec ^
: string
| set cascade.a = 1
< error
| 1193 HY000 Unknown system variable 'cascade.a'
> exec ^
: string
| set cascaded.a = 1
< error
| 1193 HY000 Unknown system variable 'cascaded.a'
> exec ^
: string
| set case.a = 1
< error
| 1193 HY000 Unknown system variable 'case.a'
> exec ^
: string
| set cast.a = 1
< error
| 1193 HY000 Unknown system variable 'cast.a'
> exec ^
: string
| set catalog_name.a = 1
< error
| 1193 HY000 Unknown system variable 'catalog_name.a'
> exec ^
: string
| set chain.a = 1
< error
| 1193 HY000 Unknown system variable 'chain.a'
> exec ^
: string
| set challenge_response.a = 1
< error
| 1193 HY000 Unknown system variable 'challenge_response.a'
> exec ^
: string
| set change.a = 1
< error
| 1193 HY000 Unknown system variable 'change.a'
> exec ^
: string
| set changed.a = 1
< error
| 1193 HY000 Unknown system variable 'changed.a'
> exec ^
: string
| set channel.a = 1
< error
| 1193 HY000 Unknown system variable 'channel.a'
> exec ^
: string
| set char.a = 1
< error
| 1193 HY000 Unknown system variable 'char.a'
> exec ^
: string
| set character.a = 1
< error
| 1193 HY000 Unknown system variable 'character.a'
> exec ^
: string
| set charset.a = 1
< error
| 1193 HY000 Unknown system variable 'charset.a'
> exec ^
: string
| set check.a = 1
< error
| 1193 HY000 Unknown system variable 'check.a'
> exec ^
: string
| set checksum.a = 1
< error
| 1193 HY000 Unknown system variable 'checksum.a'
> exec ^
: string
| set cipher.a = 1
< error
| 1193 HY000 Unknown system variable 'cipher.a'
> exec ^
: string
| set class_origin.a = 1
< error
| 1193 HY000 Unknown system variable 'class_origin.a'
> exec ^
: string
| set client.a = 1
< error
| 1193 HY000 Unknown system variable 'client.a'
> exec ^
: string
| set clone.a = 1
< error
| 1193 HY000 Unknown system variable 'clone.a'
> exec ^
: string
| set close.a = 1
< error
| 1193 HY000 Unknown system variable 'close.a'
> exec ^
: string
| set coalesce.a = 1
< error
| 1193 HY000 Unknown system variable 'coalesce.a'
> exec ^
: string
| set code.a = 1
< error
| 1193 HY000 Unknown system variable 'code.a'
> exec ^
: string
| set collate.a = 1
< error
| 1193 HY000 Unknown system variable 'collate.a'
> exec ^
: string
| set collation.a = 1
< error
| 1193 HY000 Unknown system variable 'collation.a'
> exec ^
: string
| set column.a = 1
< error
| 1193 HY000 Unknown system variable 'column.a'
> exec ^
: string
| set column_format.a = 1
< error
| 1193 HY000 Unknown system variable 'column_format.a'
> exec ^
: string
| set column_name.a = 1
< error
| 1193 HY000 Unknown system variable 'column_name.a'
> exec ^
: string
| set columns.a = 1
< error
| 1193 HY000 Unknown system variable 'columns.a'
> exec ^
: string
| set comment.a = 1
< error
| 1193 HY000 Unknown system variable 'comment.a'
> exec ^
: string
| set commit.a = 1
< error
| 1193 HY000 Unknown system variable 'commit.a'
> exec ^
: string
| set committed.a = 1
< error
| 1193 HY000 Unknown system variable 'committed.a'
> exec ^
: string
| set compact.a = 1
< error
| 1193 HY000 Unknown system variable 'compact.a'
> exec ^
: string
| set completion.a = 1
< error
| 1193 HY000 Unknown system variable 'completion.a'
> exec ^
: string
| set component.a = 1
< error
| 1193 HY000 Unknown system variable 'component.a'
> exec ^
: string
| set compressed.a = 1
< error
| 1193 HY000 Unknown system variable 'compressed.a'
> exec ^
: string
| set compression.a = 1
< error
| 1193 HY000 Unknown system variable 'compression.a'
> exec ^
: string
| set concat.a = 1
< error
| 1193 HY000 Unknown system variable 'concat.a'
> exec ^
: string
| set concurrent.a = 1
< error
| 1193 HY000 Unknown system variable 'concurrent.a'
> exec ^
: string
| set condition.a = 1
< error
| 1193 HY000 Unknown system variable 'condition.a'
> exec ^
: string
| set connection.a = 1
< error
| 1193 HY000 Unknown system variable 'connection.a'
> exec ^
: string
| set consistent.a = 1
< error
| 1193 HY000 Unknown system variable 'consistent.a'
> exec ^
: string
| set constraint.a = 1
< error
| 1193 HY000 Unknown system variable 'constraint.a'
> exec ^
: string
| set constraint_catalog.a = 1
< error
| 1193 HY000 Unknown system variable 'constraint_catalog.a'
> exec ^
: string
| set constraint_name.a = 1
< error
| 1193 HY000 Unknown system variable 'constraint_name.a'
> exec ^
: string
| set constraint_schema.a = 1
< error
| 1193 HY000 Unknown system variable 'constraint_schema.a'
> exec ^
: string
| set contains.a = 1
< error
| 1193 HY000 Unknown system variable 'contains.a'
> exec ^
: string
| set context.a = 1
< error
| 1193 HY000 Unknown system variable 'context.a'
> exec ^
: string
| set continue.a = 1
< error
| 1193 HY000 Unknown system variable 'continue.a'
> exec ^
: string
| set convert.a = 1
< error
| 1193 HY000 Unknown system variable 'convert.a'
> exec ^
: string
| set count.a = 1
< error
| 1193 HY000 Unknown system variable 'count.a'
> exec ^
: string
| set cpu.a = 1
< error
| 1193 HY000 Unknown system variable 'cpu.a'
> exec ^
: string
| set create.a = 1
< error
| 1193 HY000 Unknown system variable 'create.a'
> exec ^
: string
| set cross.a = 1
< error
| 1193 HY000 Unknown system variable 'cross.a'
> exec ^
: string
| set csv.a = 1
< error
| 1193 HY000 Unknown system variable 'csv.a'
> exec ^
: string
| set cube.a = 1
< error
| 1193 HY000 Unknown system variable 'cube.a'
> exec ^
: string
| set cume_dist.a = 1
< error
| 1193 HY000 Unknown system variable 'cume_dist.a'
> exec ^
: string
| set curdate.a = 1
< error
| 1193 HY000 Unknown system variable 'curdate.a'
> exec ^
: string
| set current.a = 1
< error
| 1193 HY000 Unknown system variable 'current.a'
> exec ^
: string
| set current_date.a = 1
< error
| 1193 HY000 Unknown system variable 'current_date.a'
> exec ^
: string
| set current_time.a = 1
< error
| 1193 HY000 Unknown system variable 'current_time.a'
> exec ^
: string
| set current_timestamp.a = 1
< error
| 1193 HY000 Unknown system variable 'current_timestamp.a'
> exec ^
: string
| set current_user.a = 1
< error
| 1193 HY000 Unknown system variable 'current_user.a'
> exec ^
: string
| set cursor.a = 1
< error
| 1193 HY000 Unknown system variable 'cursor.a'
> exec ^
: string
| set cursor_name.a = 1
< error
| 1193 HY000 Unknown system variable 'cursor_name.a'
> exec ^
: string
| set curtime.a = 1
< error
| 1193 HY000 Unknown system variable 'curtime.a'
> exec ^
: string
| set data.a = 1
< error
| 1193 HY000 Unknown system variable 'data.a'
> exec ^
: string
| set database.a = 1
< error
| 1193 HY000 Unknown system variable 'database.a'
> exec ^
: string
| set databases.a = 1
< error
| 1193 HY000 Unknown system variable 'databases.a'
> exec ^
: string
| set datafile.a = 1
< error
| 1193 HY000 Unknown system variable 'datafile.a'
> exec ^
: string
| set date.a = 1
< error
| 1193 HY000 Unknown system variable 'date.a'
> exec ^
: string
| set date_add.a = 1
< error
| 1193 HY000 Unknown system variable 'date_add.a'
> exec ^
: string
| set date_sub.a = 1
< error
| 1193 HY000 Unknown system variable 'date_sub.a'
> exec ^
: string
| set datetime.a = 1
< error
| 1193 HY000 Unknown system variable 'datetime.a'
> exec ^
: string
| set day.a = 1
< error
| 1193 HY000 Unknown system variable 'day.a'
> exec ^
: string
| set day_hour.a = 1
< error
| 1193 HY000 Unknown system variable 'day_hour.a'
> exec ^
: string
| set day_microsecond.a = 1
< error
| 1193 HY000 Unknown system variable 'day_microsecond.a'
> exec ^
: string
| set day_minute.a = 1
< error
| 1193 HY000 Unknown system variable 'day_minute.a'
> exec ^
: string
| set day_second.a = 1
< error
| 1193 HY000 Unknown system variable 'day_second.a'
> exec ^
: string
| set deallocate.a = 1
< error
| 1193 HY000 Unknown system variable 'deallocate.a'
> exec ^
: string
| set dec.a = 1
< error
| 1193 HY000 Unknown system variable 'dec.a'
> exec ^
: string
| set decimal.a = 1
< error
| 1193 HY000 Unknown system variable 'decimal.a'
> exec ^
: string
| set declare.a = 1
< error
| 1193 HY000 Unknown system variable 'declare.a'
> exec ^
: string
| set default.a = 1
< error
| 1193 HY000 Unknown system variable 'default.a'
> exec ^
: string
| set default_auth.a = 1
< error
| 1193 HY000 Unknown system variable 'default_auth.a'
> exec ^
: string
| set definer.a = 1
< error
| 1193 HY000 Unknown system variable 'definer.a'
> exec ^
: string
| set definition.a = 1
< error
| 1193 HY000 Unknown system variable 'definition.a'
> exec ^
: string
| set delay_key_write.a = 1
< error
| 1193 HY000 Unknown system variable 'delay_key_write.a'
> exec ^
: string
| set delayed.a = 1
< error
| 1193 HY000 Unknown system variable 'delayed.a'
> exec ^
: string
| set delete.a = 1
< error
| 1193 HY000 Unknown system variable 'delete.a'
> exec ^
: string
| set dense_rank.a = 1
< error
| 1193 HY000 Unknown system variable 'dense_rank.a'
> exec ^
: string
| set desc.a = 1
< error
| 1193 HY000 Unknown system variable 'desc.a'
> exec ^
: string
| set describe.a = 1
< error
| 1193 HY000 Unknown system variable 'describe.a'
> exec ^
: string
| set description.a = 1
< error
| 1193 HY000 Unknown system variable 'description.a'
> exec ^
: string
| set deterministic.a = 1
< error
| 1193 HY000 Unknown system variable 'deterministic.a'
> exec ^
: string
| set diagnostics.a = 1
< error
| 1193 HY000 Unknown system variable 'diagnostics.a'
> exec ^
: string
| set directory.a = 1
< error
| 1193 HY000 Unknown system variable 'directory.a'
> exec ^
: string
| set disable.a = 1
< error
| 1193 HY000 Unknown system variable 'disable.a'
> exec ^
: string
| set discard.a = 1
< error
| 1193 HY000 Unknown system variable 'discard.a'
> exec ^
: string
| set disk.a = 1
< error
| 1193 HY000 Unknown system variable 'disk.a'
> exec ^
: string
| set distinct.a = 1
< error
| 1193 HY000 Unknown system variable 'distinct.a'
> exec ^
: string
| set distinctrow.a = 1
< error
| 1193 HY000 Unknown system variable 'distinctrow.a'
> exec ^
: string
| set div.a = 1
< error
| 1193 HY000 Unknown system variable 'div.a'
> exec ^
: string
| set do.a = 1
< error
| 1193 HY000 Unknown system variable 'do.a'
> exec ^
: string
| set double.a = 1
< error
| 1193 HY000 Unknown system variable 'double.a'
> exec ^
: string
| set drop.a = 1
< error
| 1193 HY000 Unknown system variable 'drop.a'
> exec ^
: string
| set dual.a = 1
< error
| 1193 HY000 Unknown system variable 'dual.a'
> exec ^
: string
| set duality.a = 1
< error
| 1193 HY000 Unknown system variable 'duality.a'
> exec ^
: string
| set dumpfile.a = 1
< error
| 1193 HY000 Unknown system variable 'dumpfile.a'
> exec ^
: string
| set duplicate.a = 1
< error
| 1193 HY000 Unknown system variable 'duplicate.a'
> exec ^
: string
| set dynamic.a = 1
< error
| 1193 HY000 Unknown system variable 'dynamic.a'
> exec ^
: string
| set each.a = 1
< error
| 1193 HY000 Unknown system variable 'each.a'
> exec ^
: string
| set else.a = 1
< error
| 1193 HY000 Unknown system variable 'else.a'
> exec ^
: string
| set elseif.a = 1
< error
| 1193 HY000 Unknown system variable 'elseif.a'
> exec ^
: string
| set empty.a = 1
< error
| 1193 HY000 Unknown system variable 'empty.a'
> exec ^
: string
| set enable.a = 1
< error
| 1193 HY000 Unknown system variable 'enable.a'
> exec ^
: string
| set enclosed.a = 1
< error
| 1193 HY000 Unknown system variable 'enclosed.a'
> exec ^
: string
| set encryption.a = 1
< error
| 1193 HY000 Unknown system variable 'encryption.a'
> exec ^
: string
| set end.a = 1
< error
| 1193 HY000 Unknown system variable 'end.a'
> exec ^
: string
| set ends.a = 1
< error
| 1193 HY000 Unknown system variable 'ends.a'
> exec ^
: string
| set enforced.a = 1
< error
| 1193 HY000 Unknown system variable 'enforced.a'
> exec ^
: string
| set engine.a = 1
< error
| 1193 HY000 Unknown system variable 'engine.a'
> exec ^
: string
| set engine_attribute.a = 1
< error
| 1193 HY000 Unknown system variable 'engine_attribute.a'
> exec ^
: string
| set engines.a = 1
< error
| 1193 HY000 Unknown system variable 'engines.a'
> exec ^
: string
| set enum.a = 1
< error
| 1193 HY000 Unknown system variable 'enum.a'
> exec ^
: string
| set error.a = 1
< error
| 1193 HY000 Unknown system variable 'error.a'
> exec ^
: string
| set errors.a = 1
< error
| 1193 HY000 Unknown system variable 'errors.a'
> exec ^
: string
| set escape.a = 1
< error
| 1193 HY000 Unknown system variable 'escape.a'
> exec ^
: string
| set escaped.a = 1
< error
| 1193 HY000 Unknown system variable 'escaped.a'
> exec ^
: string
| set event.a = 1
< error
| 1193 HY000 Unknown system variable 'event.a'
> exec ^
: string
| set events.a = 1
< error
| 1193 HY000 Unknown system variable 'events.a'
> exec ^
: string
| set every.a = 1
< error
| 1193 HY000 Unknown system variable 'every.a'
> exec ^
: string
| set except.a = 1
< error
| 1193 HY000 Unknown system variable 'except.a'
> exec ^
: string
| set exchange.a = 1
< error
| 1193 HY000 Unknown system variable 'exchange.a'
> exec ^
: string
| set exclude.a = 1
< error
| 1193 HY000 Unknown system variable 'exclude.a'
> exec ^
: string
| set execute.a = 1
< error
| 1193 HY000 Unknown system variable 'execute.a'
> exec ^
: string
| set exists.a = 1
< error
| 1193 HY000 Unknown system variable 'exists.a'
> exec ^
: string
| set exit.a = 1
< error
| 1193 HY000 Unknown system variable 'exit.a'
> exec ^
: string
| set expansion.a = 1
< error
| 1193 HY000 Unknown system variable 'expansion.a'
> exec ^
: string
| set expire.a = 1
< error
| 1193 HY000 Unknown system variable 'expire.a'
> exec ^
: string
| set explain.a = 1
< error
| 1193 HY000 Unknown system variable 'explain.a'
> exec ^
: string
| set export.a = 1
< error
| 1193 HY000 Unknown system variable 'export.a'
> exec ^
: string
| set extended.a = 1
< error
| 1193 HY000 Unknown system variable 'extended.a'
> exec ^
: string
| set extent_size.a = 1
< error
| 1193 HY000 Unknown system variable 'extent_size.a'
> exec ^
: string
| set external.a = 1
< error
| 1193 HY000 Unknown system variable 'external.a'
> exec ^
: string
| set external_format.a = 1
< error
| 1193 HY000 Unknown system variable 'external_format.a'
> exec ^
: string
| set extract.a = 1
< error
| 1193 HY000 Unknown system variable 'extract.a'
> exec ^
: string
| set factor.a = 1
< error
| 1193 HY000 Unknown system variable 'factor.a'
> exec ^
: string
| set failed_login_attempts.a = 1
< error
| 1193 HY000 Unknown system variable 'failed_login_attempts.a'
> exec ^
: string
| set false.a = 1
< error
| 1193 HY000 Unknown system variable 'false.a'
> exec ^
: string
| set fast.a = 1
< error
| 1193 HY000 Unknown system variable 'fast.a'
> exec ^
: string
| set faults.a = 1
< error
| 1193 HY000 Unknown system variable 'faults.a'
> exec ^
: string
| set fetch.a = 1
< error
| 1193 HY000 Unknown system variable 'fetch.a'
> exec ^
: string
| set fields.a = 1
< error
| 1193 HY000 Unknown system variable 'fields.a'
> exec ^
: string
| set file.a = 1
< error
| 1193 HY000 Unknown system variable 'file.a'
> exec ^
: string
| set file_block_size.a = 1
< error
| 1193 HY000 Unknown system variable 'file_block_size.a'
> exec ^
: string
| set file_format.a = 1
< error
| 1193 HY000 Unknown system variable 'file_format.a'
> exec ^
: string
| set file_name.a = 1
< error
| 1193 HY000 Unknown system variable 'file_name.a'
> exec ^
: string
| set file_pattern.a = 1
< error
| 1193 HY000 Unknown system variable 'file_pattern.a'
> exec ^
: string
| set file_prefix.a = 1
< error
| 1193 HY000 Unknown system variable 'file_prefix.a'
> exec ^
: string
| set files.a = 1
< error
| 1193 HY000 Unknown system variable 'files.a'
> exec ^
: string
| set filter.a = 1
< error
| 1193 HY000 Unknown system variable 'filter.a'
> exec ^
: string
| set finish.a = 1
< error
| 1193 HY000 Unknown system variable 'finish.a'
> exec ^
: string
| set first.a = 1
< error
| 1193 HY000 Unknown system variable 'first.a'
> exec ^
: string
| set first_value.a = 1
< error
| 1193 HY000 Unknown system variable 'first_value.a'
> exec ^
: string
| set fixed.a = 1
< error
| 1193 HY000 Unknown system variable 'fixed.a'
> exec ^
: string
| set float.a = 1
< error
| 1193 HY000 Unknown system variable 'float.a'
> exec ^
: string
| set float4.a = 1
< error
| 1193 HY000 Unknown system variable 'float4.a'
> exec ^
: string
| set float8.a = 1
< error
| 1193 HY000 Unknown system variable 'float8.a'
> exec ^
: string
| set flush.a = 1
< error
| 1193 HY000 Unknown system variable 'flush.a'
> exec ^
: string
| set following.a = 1
< error
| 1193 HY000 Unknown system variable 'following.a'
> exec ^
: string
| set follows.a = 1
< error
| 1193 HY000 Unknown system variable 'follows.a'
> exec ^
: string
| set for.a = 1
< error
| 1193 HY000 Unknown system variable 'for.a'
> exec ^
: string
| set force.a = 1
< error
| 1193 HY000 Unknown system variable 'force.a'
> exec ^
: string
| set foreign.a = 1
< error
| 1193 HY000 Unknown system variable 'foreign.a'
> exec ^
: string
| set format.a = 1
< error
| 1193 HY000 Unknown system variable 'format.a'
> exec ^
: string
| set found.a = 1
< error
| 1193 HY000 Unknown system variable 'found.a'
> exec ^
: string
| set from.a = 1
< error
| 1193 HY000 Unknown system variable 'from.a'
> exec ^
: string
| set full.a = 1
< error
| 1193 HY000 Unknown system variable 'full.a'
> exec ^
: string
| set fulltext.a = 1
< error
| 1193 HY000 Unknown system variable 'fulltext.a'
> exec ^
: string
| set function.a = 1
< error
| 1193 HY000 Unknown system variable 'function.a'
> exec ^
: string
| set general.a = 1
< error
| 1193 HY000 Unknown system variable 'general.a'
> exec ^
: string
| set generate.a = 1
< error
| 1193 HY000 Unknown system variable 'generate.a'
> exec ^
: string
| set generated.a = 1
< error
| 1193 HY000 Unknown system variable 'generated.a'
> exec ^
: string
| set geomcollection.a = 1
< error
| 1193 HY000 Unknown system variable 'geomcollection.a'
> exec ^
: string
| set geometry.a = 1
< error
| 1193 HY000 Unknown system variable 'geometry.a'
> exec ^
: string
| set geometrycollection.a = 1
< error
| 1193 HY000 Unknown system variable 'geometrycollection.a'
> exec ^
: string
| set get.a = 1
< error
| 1193 HY000 Unknown system variable 'get.a'
> exec ^
: string
| set get_format.a = 1
< error
| 1193 HY000 Unknown system variable 'get_format.a'
> exec ^
: string
| set get_master_public_key.a = 1
< error
| 1193 HY000 Unknown system variable 'get_master_public_key.a'
> exec ^
: string
| set get_source_public_key.a = 1
< error
| 1193 HY000 Unknown system variable 'get_source_public_key.a'
> exec ^
: string
| set global.a = 1
< error
| 1064 42000 …near 'global.a = 1' at line 1
> exec ^
: string
| set grant.a = 1
< error
| 1193 HY000 Unknown system variable 'grant.a'
> exec ^
: string
| set grants.a = 1
< error
| 1193 HY000 Unknown system variable 'grants.a'
> exec ^
: string
| set greatest.a = 1
< error
| 1193 HY000 Unknown system variable 'greatest.a'
> exec ^
: string
| set group.a = 1
< error
| 1193 HY000 Unknown system variable 'group.a'
> exec ^
: string
| set group_concat.a = 1
< error
| 1193 HY000 Unknown system variable 'group_concat.a'
> exec ^
: string
| set group_replication.a = 1
< error
| 1193 HY000 Unknown system variable 'group_replication.a'
> exec ^
: string
| set grouping.a = 1
< error
| 1193 HY000 Unknown system variable 'grouping.a'
> exec ^
: string
| set groups.a = 1
< error
| 1193 HY000 Unknown system variable 'groups.a'
> exec ^
: string
| set gtid_only.a = 1
< error
| 1193 HY000 Unknown system variable 'gtid_only.a'
> exec ^
: string
| set gtids.a = 1
< error
| 1193 HY000 Unknown system variable 'gtids.a'
> exec ^
: string
| set guided.a = 1
< error
| 1193 HY000 Unknown system variable 'guided.a'
> exec ^
: string
| set handler.a = 1
< error
| 1193 HY000 Unknown system variable 'handler.a'
> exec ^
: string
| set hash.a = 1
< error
| 1193 HY000 Unknown system variable 'hash.a'
> exec ^
: string
| set having.a = 1
< error
| 1193 HY000 Unknown system variable 'having.a'
> exec ^
: string
| set header.a = 1
< error
| 1193 HY000 Unknown system variable 'header.a'
> exec ^
: string
| set help.a = 1
< error
| 1193 HY000 Unknown system variable 'help.a'
> exec ^
: string
| set high_priority.a = 1
< error
| 1193 HY000 Unknown system variable 'high_priority.a'
> exec ^
: string
| set histogram.a = 1
< error
| 1193 HY000 Unknown system variable 'histogram.a'
> exec ^
: string
| set history.a = 1
< error
| 1193 HY000 Unknown system variable 'history.a'
> exec ^
: string
| set host.a = 1
< error
| 1193 HY000 Unknown system variable 'host.a'
> exec ^
: string
| set hosts.a = 1
< error
| 1193 HY000 Unknown system variable 'hosts.a'
> exec ^
: string
| set hour.a = 1
< error
| 1193 HY000 Unknown system variable 'hour.a'
> exec ^
: string
| set hour_microsecond.a = 1
< error
| 1193 HY000 Unknown system variable 'hour_microsecond.a'
> exec ^
: string
| set hour_minute.a = 1
< error
| 1193 HY000 Unknown system variable 'hour_minute.a'
> exec ^
: string
| set hour_second.a = 1
< error
| 1193 HY000 Unknown system variable 'hour_second.a'
> exec ^
: string
| set identified.a = 1
< error
| 1193 HY000 Unknown system variable 'identified.a'
> exec ^
: string
| set if.a = 1
< error
| 1193 HY000 Unknown system variable 'if.a'
> exec ^
: string
| set ignore.a = 1
< error
| 1193 HY000 Unknown system variable 'ignore.a'
> exec ^
: string
| set ignore_server_ids.a = 1
< error
| 1193 HY000 Unknown system variable 'ignore_server_ids.a'
> exec ^
: string
| set import.a = 1
< error
| 1193 HY000 Unknown system variable 'import.a'
> exec ^
: string
| set in.a = 1
< error
| 1193 HY000 Unknown system variable 'in.a'
> exec ^
: string
| set inactive.a = 1
< error
| 1193 HY000 Unknown system variable 'inactive.a'
> exec ^
: string
| set index.a = 1
< error
| 1193 HY000 Unknown system variable 'index.a'
> exec ^
: string
| set indexes.a = 1
< error
| 1193 HY000 Unknown system variable 'indexes.a'
> exec ^
: string
| set infile.a = 1
< error
| 1193 HY000 Unknown system variable 'infile.a'
> exec ^
: string
| set initial.a = 1
< error
| 1193 HY000 Unknown system variable 'initial.a'
> exec ^
: string
| set initial_size.a = 1
< error
| 1193 HY000 Unknown system variable 'initial_size.a'
> exec ^
: string
| set initiate.a = 1
< error
| 1193 HY000 Unknown system variable 'initiate.a'
> exec ^
: string
| set inner.a = 1
< error
| 1193 HY000 Unknown system variable 'inner.a'
> exec ^
: string
| set inout.a = 1
< error
| 1193 HY000 Unknown system variable 'inout.a'
> exec ^
: string
| set insensitive.a = 1
< error
| 1193 HY000 Unknown system variable 'insensitive.a'
> exec ^
: string
| set insert.a = 1
< error
| 1193 HY000 Unknown system variable 'insert.a'
> exec ^
: string
| set insert_method.a = 1
< error
| 1193 HY000 Unknown system variable 'insert_method.a'
> exec ^
: string
| set install.a = 1
< error
| 1193 HY000 Unknown system variable 'install.a'
> exec ^
: string
| set instance.a = 1
< error
| 1193 HY000 Unknown system variable 'instance.a'
> exec ^
: string
| set int.a = 1
< error
| 1193 HY000 Unknown system variable 'int.a'
> exec ^
: string
| set int1.a = 1
< error
| 1193 HY000 Unknown system variable 'int1.a'
> exec ^
: string
| set int2.a = 1
< error
| 1193 HY000 Unknown system variable 'int2.a'
> exec ^
: string
| set int3.a = 1
< error
| 1193 HY000 Unknown system variable 'int3.a'
> exec ^
: string
| set int4.a = 1
< error
| 1193 HY000 Unknown system variable 'int4.a'
> exec ^
: string
| set int8.a = 1
< error
| 1193 HY000 Unknown system variable 'int8.a'
> exec ^
: string
| set integer.a = 1
< error
| 1193 HY000 Unknown system variable 'integer.a'
> exec ^
: string
| set intersect.a = 1
< error
| 1193 HY000 Unknown system variable 'intersect.a'
> exec ^
: string
| set interval.a = 1
< error
| 1193 HY000 Unknown system variable 'interval.a'
> exec ^
: string
| set into.a = 1
< error
| 1193 HY000 Unknown system variable 'into.a'
> exec ^
: string
| set invisible.a = 1
< error
| 1193 HY000 Unknown system variable 'invisible.a'
> exec ^
: string
| set invoker.a = 1
< error
| 1193 HY000 Unknown system variable 'invoker.a'
> exec ^
: string
| set io.a = 1
< error
| 1193 HY000 Unknown system variable 'io.a'
> exec ^
: string
| set io_after_gtids.a = 1
< error
| 1193 HY000 Unknown system variable 'io_after_gtids.a'
> exec ^
: string
| set io_before_gtids.a = 1
< error
| 1193 HY000 Unknown system variable 'io_before_gtids.a'
> exec ^
: string
| set io_thread.a = 1
< error
| 1193 HY000 Unknown system variable 'io_thread.a'
> exec ^
: string
| set ipc.a = 1
< error
| 1193 HY000 Unknown system variable 'ipc.a'
> exec ^
: string
| set is.a = 1
< error
| 1193 HY000 Unknown system variable 'is.a'
> exec ^
: string
| set isolation.a = 1
< error
| 1193 HY000 Unknown system variable 'isolation.a'
> exec ^
: string
| set issuer.a = 1
< error
| 1193 HY000 Unknown system variable 'issuer.a'
> exec ^
: string
| set iterate.a = 1
< error
| 1193 HY000 Unknown system variable 'iterate.a'
> exec ^
: string
| set join.a = 1
< error
| 1193 HY000 Unknown system variable 'join.a'
> exec ^
: string
| set json.a = 1
< error
| 1193 HY000 Unknown system variable 'json.a'
> exec ^
: string
| set json_arrayagg.a = 1
< error
| 1193 HY000 Unknown system variable 'json_arrayagg.a'
> exec ^
: string
| set json_duality_object.a = 1
< error
| 1193 HY000 Unknown system variable 'json_duality_object.a'
> exec ^
: string
| set json_objectagg.a = 1
< error
| 1193 HY000 Unknown system variable 'json_objectagg.a'
> exec ^
: string
| set json_table.a = 1
< error
| 1193 HY000 Unknown system variable 'json_table.a'
> exec ^
: string
| set json_value.a = 1
< error
| 1193 HY000 Unknown system variable 'json_value.a'
> exec ^
: string
| set key.a = 1
< error
| 1193 HY000 Unknown system variable 'key.a'
> exec ^
: string
| set key_block_size.a = 1
< error
| 1193 HY000 Unknown system variable 'key_block_size.a'
> exec ^
: string
| set keyring.a = 1
< error
| 1193 HY000 Unknown system variable 'keyring.a'
> exec ^
: string
| set keys.a = 1
< error
| 1193 HY000 Unknown system variable 'keys.a'
> exec ^
: string
| set kill.a = 1
< error
| 1193 HY000 Unknown system variable 'kill.a'
> exec ^
: string
| set lag.a = 1
< error
| 1193 HY000 Unknown system variable 'lag.a'
> exec ^
: string
| set language.a = 1
< error
| 1193 HY000 Unknown system variable 'language.a'
> exec ^
: string
| set last.a = 1
< error
| 1193 HY000 Unknown system variable 'last.a'
> exec ^
: string
| set last_value.a = 1
< error
| 1193 HY000 Unknown system variable 'last_value.a'
> exec ^
: string
| set lateral.a = 1
< error
| 1193 HY000 Unknown system variable 'lateral.a'
> exec ^
: string
| set lead.a = 1
< error
| 1193 HY000 Unknown system variable 'lead.a'
> exec ^
: string
| set leading.a = 1
< error
| 1193 HY000 Unknown system variable 'leading.a'
> exec ^
: string
| set least.a = 1
< error
| 1193 HY000 Unknown system variable 'least.a'
> exec ^
: string
| set leave.a = 1
< error
| 1193 HY000 Unknown system variable 'leave.a'
> exec ^
: string
| set leaves.a = 1
< error
| 1193 HY000 Unknown system variable 'leaves.a'
> exec ^
: string
| set left.a = 1
< error
| 1193 HY000 Unknown system variable 'left.a'
> exec ^
: string
| set less.a = 1
< error
| 1193 HY000 Unknown system variable 'less.a'
> exec ^
: string
| set level.a = 1
< error
| 1193 HY000 Unknown system variable 'level.a'
> exec ^
: string
| set library.a = 1
< error
| 1193 HY000 Unknown system variable 'library.a'
> exec ^
: string
| set like.a = 1
< error
| 1193 HY000 Unknown system variable 'like.a'
> exec ^
: string
| set limit.a = 1
< error
| 1193 HY000 Unknown system variable 'limit.a'
> exec ^
: string
| set linear.a = 1
< error
| 1193 HY000 Unknown system variable 'linear.a'
> exec ^
: string
| set lines.a = 1
< error
| 1193 HY000 Unknown system variable 'lines.a'
> exec ^
: string
| set linestring.a = 1
< error
| 1193 HY000 Unknown system variable 'linestring.a'
> exec ^
: string
| set list.a = 1
< error
| 1193 HY000 Unknown system variable 'list.a'
> exec ^
: string
| set load.a = 1
< error
| 1193 HY000 Unknown system variable 'load.a'
> exec ^
: string
| set local.a = 1
< error
| 1064 42000 …near 'local.a = 1' at line 1
> exec ^
: string
| set localtime.a = 1
< error
| 1193 HY000 Unknown system variable 'localtime.a'
> exec ^
: string
| set localtimestamp.a = 1
< error
| 1193 HY000 Unknown system variable 'localtimestamp.a'
> exec ^
: string
| set lock.a = 1
< error
| 1193 HY000 Unknown system variable 'lock.a'
> exec ^
: string
| set locked.a = 1
< error
| 1193 HY000 Unknown system variable 'locked.a'
> exec ^
: string
| set locks.a = 1
< error
| 1193 HY000 Unknown system variable 'locks.a'
> exec ^
: string
| set log.a = 1
< error
| 1193 HY000 Unknown system variable 'log.a'
> exec ^
: string
| set logfile.a = 1
< error
| 1193 HY000 Unknown system variable 'logfile.a'
> exec ^
: string
| set logs.a = 1
< error
| 1193 HY000 Unknown system variable 'logs.a'
> exec ^
: string
| set long.a = 1
< error
| 1193 HY000 Unknown system variable 'long.a'
> exec ^
: string
| set longblob.a = 1
< error
| 1193 HY000 Unknown system variable 'longblob.a'
> exec ^
: string
| set longtext.a = 1
< error
| 1193 HY000 Unknown system variable 'longtext.a'
> exec ^
: string
| set loop.a = 1
< error
| 1193 HY000 Unknown system variable 'loop.a'
> exec ^
: string
| set low_priority.a = 1
< error
| 1193 HY000 Unknown system variable 'low_priority.a'
> exec ^
: string
| set manifest.a = 1
< error
| 1193 HY000 Unknown system variable 'manifest.a'
> exec ^
: string
| set manual.a = 1
< error
| 1193 HY000 Unknown system variable 'manual.a'
> exec ^
: string
| set master.a = 1
< error
| 1193 HY000 Unknown system variable 'master.a'
> exec ^
: string
| set master_auto_position.a = 1
< error
| 1193 HY000 Unknown system variable 'master_auto_position.a'
> exec ^
: string
| set master_bind.a = 1
< error
| 1193 HY000 Unknown system variable 'master_bind.a'
> exec ^
: string
| set master_compression_algorithms.a = 1
< error
| 1193 HY000 Unknown system variable 'master_compression_algorithms.a'
> exec ^
: string
| set master_connect_retry.a = 1
< error
| 1193 HY000 Unknown system variable 'master_connect_retry.a'
> exec ^
: string
| set master_delay.a = 1
< error
| 1193 HY000 Unknown system variable 'master_delay.a'
> exec ^
: string
| set master_heartbeat_period.a = 1
< error
| 1193 HY000 Unknown system variable 'master_heartbeat_period.a'
> exec ^
: string
| set master_host.a = 1
< error
| 1193 HY000 Unknown system variable 'master_host.a'
> exec ^
: string
| set master_log_file.a = 1
< error
| 1193 HY000 Unknown system variable 'master_log_file.a'
> exec ^
: string
| set master_log_pos.a = 1
< error
| 1193 HY000 Unknown system variable 'master_log_pos.a'
> exec ^
: string
| set master_password.a = 1
< error
| 1193 HY000 Unknown system variable 'master_password.a'
> exec ^
: string
| set master_port.a = 1
< error
| 1193 HY000 Unknown system variable 'master_port.a'
> exec ^
: string
| set master_public_key_path.a = 1
< error
| 1193 HY000 Unknown system variable 'master_public_key_path.a'
> exec ^
: string
| set master_retry_count.a = 1
< error
| 1193 HY000 Unknown system variable 'master_retry_count.a'
> exec ^
: string
| set master_ssl.a = 1
< error
| 1193 HY000 Unknown system variable 'master_ssl.a'
> exec ^
: string
| set master_ssl_ca.a = 1
< error
| 1193 HY000 Unknown system variable 'master_ssl_ca.a'
> exec ^
: string
| set master_ssl_capath.a = 1
< error
| 1193 HY000 Unknown system variable 'master_ssl_capath.a'
> exec ^
: string
| set master_ssl_cert.a = 1
< error
| 1193 HY000 Unknown system variable 'master_ssl_cert.a'
> exec ^
: string
| set master_ssl_cipher.a = 1
< error
| 1193 HY000 Unknown system variable 'master_ssl_cipher.a'
> exec ^
: string
| set master_ssl_crl.a = 1
< error
| 1193 HY000 Unknown system variable 'master_ssl_crl.a'
> exec ^
: string
| set master_ssl_crlpath.a = 1
< error
| 1193 HY000 Unknown system variable 'master_ssl_crlpath.a'
> exec ^
: string
| set master_ssl_key.a = 1
< error
| 1193 HY000 Unknown system variable 'master_ssl_key.a'
> exec ^
: string
| set master_ssl_verify_server_cert.a = 1
< error
| 1193 HY000 Unknown system variable 'master_ssl_verify_server_cert.a'
> exec ^
: string
| set master_tls_ciphersuites.a = 1
< error
| 1193 HY000 Unknown system variable 'master_tls_ciphersuites.a'
> exec ^
: string
| set master_tls_version.a = 1
< error
| 1193 HY000 Unknown system variable 'master_tls_version.a'
> exec ^
: string
| set master_user.a = 1
< error
| 1193 HY000 Unknown system variable 'master_user.a'
> exec ^
: string
| set master_zstd_compression_level.a = 1
< error
| 1193 HY000 Unknown system variable 'master_zstd_compression_level.a'
> exec ^
: string
| set match.a = 1
< error
| 1193 HY000 Unknown system variable 'match.a'
> exec ^
: string
| set materialized.a = 1
< error
| 1193 HY000 Unknown system variable 'materialized.a'
> exec ^
: string
| set max.a = 1
< error
| 1193 HY000 Unknown system variable 'max.a'
> exec ^
: string
| set max_connections_per_hour.a = 1
< error
| 1193 HY000 Unknown system variable 'max_connections_per_hour.a'
> exec ^
: string
| set max_queries_per_hour.a = 1
< error
| 1193 HY000 Unknown system variable 'max_queries_per_hour.a'
> exec ^
: string
| set max_rows.a = 1
< error
| 1193 HY000 Unknown system variable 'max_rows.a'
> exec ^
: string
| set max_size.a = 1
< error
| 1193 HY000 Unknown system variable 'max_size.a'
> exec ^
: string
| set max_updates_per_hour.a = 1
< error
| 1193 HY000 Unknown system variable 'max_updates_per_hour.a'
> exec ^
: string
| set max_user_connections.a = 1
< error
| 1193 HY000 Unknown system variable 'max_user_connections.a'
> exec ^
: string
| set maxvalue.a = 1
< error
| 1193 HY000 Unknown system variable 'maxvalue.a'
> exec ^
: string
| set medium.a = 1
< error
| 1193 HY000 Unknown system variable 'medium.a'
> exec ^
: string
| set mediumblob.a = 1
< error
| 1193 HY000 Unknown system variable 'mediumblob.a'
> exec ^
: string
| set mediumint.a = 1
< error
| 1193 HY000 Unknown system variable 'mediumint.a'
> exec ^
: string
| set mediumtext.a = 1
< error
| 1193 HY000 Unknown system variable 'mediumtext.a'
> exec ^
: string
| set member.a = 1
< error
| 1193 HY000 Unknown system variable 'member.a'
> exec ^
: string
| set memory.a = 1
< error
| 1193 HY000 Unknown system variable 'memory.a'
> exec ^
: string
| set merge.a = 1
< error
| 1193 HY000 Unknown system variable 'merge.a'
> exec ^
: string
| set message_text.a = 1
< error
| 1193 HY000 Unknown system variable 'message_text.a'
> exec ^
: string
| set microsecond.a = 1
< error
| 1193 HY000 Unknown system variable 'microsecond.a'
> exec ^
: string
| set mid.a = 1
< error
| 1193 HY000 Unknown system variable 'mid.a'
> exec ^
: string
| set middleint.a = 1
< error
| 1193 HY000 Unknown system variable 'middleint.a'
> exec ^
: string
| set migrate.a = 1
< error
| 1193 HY000 Unknown system variable 'migrate.a'
> exec ^
: string
| set min.a = 1
< error
| 1193 HY000 Unknown system variable 'min.a'
> exec ^
: string
| set min_rows.a = 1
< error
| 1193 HY000 Unknown system variable 'min_rows.a'
> exec ^
: string
| set minute.a = 1
< error
| 1193 HY000 Unknown system variable 'minute.a'
> exec ^
: string
| set minute_microsecond.a = 1
< error
| 1193 HY000 Unknown system variable 'minute_microsecond.a'
> exec ^
: string
| set minute_second.a = 1
< error
| 1193 HY000 Unknown system variable 'minute_second.a'
> exec ^
: string
| set mod.a = 1
< error
| 1193 HY000 Unknown system variable 'mod.a'
> exec ^
: string
| set mode.a = 1
< error
| 1193 HY000 Unknown system variable 'mode.a'
> exec ^
: string
| set modifies.a = 1
< error
| 1193 HY000 Unknown system variable 'modifies.a'
> exec ^
: string
| set modify.a = 1
< error
| 1193 HY000 Unknown system variable 'modify.a'
> exec ^
: string
| set month.a = 1
< error
| 1193 HY000 Unknown system variable 'month.a'
> exec ^
: string
| set multilinestring.a = 1
< error
| 1193 HY000 Unknown system variable 'multilinestring.a'
> exec ^
: string
| set multipoint.a = 1
< error
| 1193 HY000 Unknown system variable 'multipoint.a'
> exec ^
: string
| set multipolygon.a = 1
< error
| 1193 HY000 Unknown system variable 'multipolygon.a'
> exec ^
: string
| set mutex.a = 1
< error
| 1193 HY000 Unknown system variable 'mutex.a'
> exec ^
: string
| set mysql_errno.a = 1
< error
| 1193 HY000 Unknown system variable 'mysql_errno.a'
> exec ^
: string
| set name.a = 1
< error
| 1193 HY000 Unknown system variable 'name.a'
> exec ^
: string
| set names.a = 1
< error
| 1193 HY000 Unknown system variable 'names.a'
> exec ^
: string
| set national.a = 1
< error
| 1193 HY000 Unknown system variable 'national.a'
> exec ^
: string
| set natural.a = 1
< error
| 1193 HY000 Unknown system variable 'natural.a'
> exec ^
: string
| set nchar.a = 1
< error
| 1193 HY000 Unknown system variable 'nchar.a'
> exec ^
: string
| set ndb.a = 1
< error
| 1193 HY000 Unknown system variable 'ndb.a'
> exec ^
: string
| set ndbcluster.a = 1
< error
| 1193 HY000 Unknown system variable 'ndbcluster.a'
> exec ^
: string
| set nested.a = 1
< error
| 1193 HY000 Unknown system variable 'nested.a'
> exec ^
: string
| set network_namespace.a = 1
< error
| 1193 HY000 Unknown system variable 'network_namespace.a'
> exec ^
: string
| set never.a = 1
< error
| 1193 HY000 Unknown system variable 'never.a'
> exec ^
: string
| set new.a = 1
< error
| 1193 HY000 Unknown system variable 'new.a'
> exec ^
: string
| set next.a = 1
< error
| 1193 HY000 Unknown system variable 'next.a'
> exec ^
: string
| set no.a = 1
< error
| 1193 HY000 Unknown system variable 'no.a'
> exec ^
: string
| set no_wait.a = 1
< error
| 1193 HY000 Unknown system variable 'no_wait.a'
> exec ^
: string
| set no_write_to_binlog.a = 1
< error
| 1193 HY000 Unknown system variable 'no_write_to_binlog.a'
> exec ^
: string
| set nodegroup.a = 1
< error
| 1193 HY000 Unknown system variable 'nodegroup.a'
> exec ^
: string
| set none.a = 1
< error
| 1193 HY000 Unknown system variable 'none.a'
> exec ^
: string
| set not.a = 1
< error
| 1193 HY000 Unknown system variable 'not.a'
> exec ^
: string
| set now.a = 1
< error
| 1193 HY000 Unknown system variable 'now.a'
> exec ^
: string
| set nowait.a = 1
< error
| 1193 HY000 Unknown system variable 'nowait.a'
> exec ^
: string
| set nth_value.a = 1
< error
| 1193 HY000 Unknown system variable 'nth_value.a'
> exec ^
: string
| set ntile.a = 1
< error
| 1193 HY000 Unknown system variable 'ntile.a'
> exec ^
: string
| set null.a = 1
< error
| 1193 HY000 Unknown system variable 'null.a'
> exec ^
: string
| set nullif.a = 1
< error
| 1193 HY000 Unknown system variable 'nullif.a'
> exec ^
: string
| set nulls.a = 1
< error
| 1193 HY000 Unknown system variable 'nulls.a'
> exec ^
: string
| set number.a = 1
< error
| 1193 HY000 Unknown system variable 'number.a'
> exec ^
: string
| set numeric.a = 1
< error
| 1193 HY000 Unknown system variable 'numeric.a'
> exec ^
: string
| set nvarchar.a = 1
< error
| 1193 HY000 Unknown system variable 'nvarchar.a'
> exec ^
: string
| set of.a = 1
< error
| 1193 HY000 Unknown system variable 'of.a'
> exec ^
: string
| set off.a = 1
< error
| 1193 HY000 Unknown system variable 'off.a'
> exec ^
: string
| set offset.a = 1
< error
| 1193 HY000 Unknown system variable 'offset.a'
> exec ^
: string
| set oj.a = 1
< error
| 1193 HY000 Unknown system variable 'oj.a'
> exec ^
: string
| set old.a = 1
< error
| 1193 HY000 Unknown system variable 'old.a'
> exec ^
: string
| set on.a = 1
< error
| 1193 HY000 Unknown system variable 'on.a'
> exec ^
: string
| set one.a = 1
< error
| 1193 HY000 Unknown system variable 'one.a'
> exec ^
: string
| set only.a = 1
< error
| 1193 HY000 Unknown system variable 'only.a'
> exec ^
: string
| set open.a = 1
< error
| 1193 HY000 Unknown system variable 'open.a'
> exec ^
: string
| set optimize.a = 1
< error
| 1193 HY000 Unknown system variable 'optimize.a'
> exec ^
: string
| set optimizer_costs.a = 1
< error
| 1193 HY000 Unknown system variable 'optimizer_costs.a'
> exec ^
: string
| set option.a = 1
< error
| 1193 HY000 Unknown system variable 'option.a'
> exec ^
: string
| set optional.a = 1
< error
| 1193 HY000 Unknown system variable 'optional.a'
> exec ^
: string
| set optionally.a = 1
< error
| 1193 HY000 Unknown system variable 'optionally.a'
> exec ^
: string
| set options.a = 1
< error
| 1193 HY000 Unknown system variable 'options.a'
> exec ^
: string
| set or.a = 1
< error
| 1193 HY000 Unknown system variable 'or.a'
> exec ^
: string
| set order.a = 1
< error
| 1193 HY000 Unknown system variable 'order.a'
> exec ^
: string
| set ordinality.a = 1
< error
| 1193 HY000 Unknown system variable 'ordinality.a'
> exec ^
: string
| set organization.a = 1
< error
| 1193 HY000 Unknown system variable 'organization.a'
> exec ^
: string
| set others.a = 1
< error
| 1193 HY000 Unknown system variable 'others.a'
> exec ^
: string
| set out.a = 1
< error
| 1193 HY000 Unknown system variable 'out.a'
> exec ^
: string
| set outer.a = 1
< error
| 1193 HY000 Unknown system variable 'outer.a'
> exec ^
: string
| set outfile.a = 1
< error
| 1193 HY000 Unknown system variable 'outfile.a'
> exec ^
: string
| set over.a = 1
< error
| 1193 HY000 Unknown system variable 'over.a'
> exec ^
: string
| set overlay.a = 1
< error
| 1193 HY000 Unknown system variable 'overlay.a'
> exec ^
: string
| set overwrite.a = 1
< error
| 1193 HY000 Unknown system variable 'overwrite.a'
> exec ^
: string
| set owner.a = 1
< error
| 1193 HY000 Unknown system variable 'owner.a'
> exec ^
: string
| set pack_keys.a = 1
< error
| 1193 HY000 Unknown system variable 'pack_keys.a'
> exec ^
: string
| set page.a = 1
< error
| 1193 HY000 Unknown system variable 'page.a'
> exec ^
: string
| set parallel.a = 1
< error
| 1193 HY000 Unknown system variable 'parallel.a'
> exec ^
: string
| set parameters.a = 1
< error
| 1193 HY000 Unknown system variable 'parameters.a'
> exec ^
: string
| set parse_tree.a = 1
< error
| 1193 HY000 Unknown system variable 'parse_tree.a'
> exec ^
: string
| set parser.a = 1
< error
| 1193 HY000 Unknown system variable 'parser.a'
> exec ^
: string
| set partial.a = 1
< error
| 1193 HY000 Unknown system variable 'partial.a'
> exec ^
: string
| set partition.a = 1
< error
| 1193 HY000 Unknown system variable 'partition.a'
> exec ^
: string
| set partitioning.a = 1
< error
| 1193 HY000 Unknown system variable 'partitioning.a'
> exec ^
: string
| set partitions.a = 1
< error
| 1193 HY000 Unknown system variable 'partitions.a'
> exec ^
: string
| set password.a = 1
< error
| 1193 HY000 Unknown system variable 'password.a'
> exec ^
: string
| set password_lock_time.a = 1
< error
| 1193 HY000 Unknown system variable 'password_lock_time.a'
> exec ^
: string
| set path.a = 1
< error
| 1193 HY000 Unknown system variable 'path.a'
> exec ^
: string
| set percent_rank.a = 1
< error
| 1193 HY000 Unknown system variable 'percent_rank.a'
> exec ^
: string
| set persist.a = 1
< error
| 1193 HY000 Unknown system variable 'persist.a'
> exec ^
: string
| set persist_only.a = 1
< error
| 1193 HY000 Unknown system variable 'persist_only.a'
> exec ^
: string
| set phase.a = 1
< error
| 1193 HY000 Unknown system variable 'phase.a'
> exec ^
: string
| set placing.a = 1
< error
| 1193 HY000 Unknown system variable 'placing.a'
> exec ^
: string
| set plugin.a = 1
< error
| 1193 HY000 Unknown system variable 'plugin.a'
> exec ^
: string
| set plugin_dir.a = 1
< error
| 1193 HY000 Unknown system variable 'plugin_dir.a'
> exec ^
: string
| set plugins.a = 1
< error
| 1193 HY000 Unknown system variable 'plugins.a'
> exec ^
: string
| set point.a = 1
< error
| 1193 HY000 Unknown system variable 'point.a'
> exec ^
: string
| set polygon.a = 1
< error
| 1193 HY000 Unknown system variable 'polygon.a'
> exec ^
: string
| set port.a = 1
< error
| 1193 HY000 Unknown system variable 'port.a'
> exec ^
: string
| set position.a = 1
< error
| 1193 HY000 Unknown system variable 'position.a'
> exec ^
: string
| set precedes.a = 1
< error
| 1193 HY000 Unknown system variable 'precedes.a'
> exec ^
: string
| set preceding.a = 1
< error
| 1193 HY000 Unknown system variable 'preceding.a'
> exec ^
: string
| set precision.a = 1
< error
| 1193 HY000 Unknown system variable 'precision.a'
> exec ^
: string
| set prefix.a = 1
< error
| 1193 HY000 Unknown system variable 'prefix.a'
> exec ^
: string
| set prepare.a = 1
< error
| 1193 HY000 Unknown system variable 'prepare.a'
> exec ^
: string
| set preserve.a = 1
< error
| 1193 HY000 Unknown system variable 'preserve.a'
> exec ^
: string
| set prev.a = 1
< error
| 1193 HY000 Unknown system variable 'prev.a'
> exec ^
: string
| set primary.a = 1
< error
| 1193 HY000 Unknown system variable 'primary.a'
> exec ^
: string
| set privilege_checks_user.a = 1
< error
| 1193 HY000 Unknown system variable 'privilege_checks_user.a'
> exec ^
: string
| set privileges.a = 1
< error
| 1193 HY000 Unknown system variable 'privileges.a'
> exec ^
: string
| set procedure.a = 1
< error
| 1193 HY000 Unknown system variable 'procedure.a'
> exec ^
: string
| set process.a = 1
< error
| 1193 HY000 Unknown system variable 'process.a'
> exec ^
: string
| set processlist.a = 1
< error
| 1193 HY000 Unknown system variable 'processlist.a'
> exec ^
: string
| set profile.a = 1
< error
| 1193 HY000 Unknown system variable 'profile.a'
> exec ^
: string
| set profiles.a = 1
< error
| 1193 HY000 Unknown system variable 'profiles.a'
> exec ^
: string
| set proxy.a = 1
< error
| 1193 HY000 Unknown system variable 'proxy.a'
> exec ^
: string
| set purge.a = 1
< error
| 1193 HY000 Unknown system variable 'purge.a'
> exec ^
: string
| set qualify.a = 1
< error
| 1193 HY000 Unknown system variable 'qualify.a'
> exec ^
: string
| set quarter.a = 1
< error
| 1193 HY000 Unknown system variable 'quarter.a'
> exec ^
: string
| set query.a = 1
< error
| 1193 HY000 Unknown system variable 'query.a'
> exec ^
: string
| set quick.a = 1
< error
| 1193 HY000 Unknown system variable 'quick.a'
> exec ^
: string
| set random.a = 1
< error
| 1193 HY000 Unknown system variable 'random.a'
> exec ^
: string
| set range.a = 1
< error
| 1193 HY000 Unknown system variable 'range.a'
> exec ^
: string
| set rank.a = 1
< error
| 1193 HY000 Unknown system variable 'rank.a'
> exec ^
: string
| set read.a = 1
< error
| 1193 HY000 Unknown system variable 'read.a'
> exec ^
: string
| set read_only.a = 1
< error
| 1193 HY000 Unknown system variable 'read_only.a'
> exec ^
: string
| set read_write.a = 1
< error
| 1193 HY000 Unknown system variable 'read_write.a'
> exec ^
: string
| set reads.a = 1
< error
| 1193 HY000 Unknown system variable 'reads.a'
> exec ^
: string
| set real.a = 1
< error
| 1193 HY000 Unknown system variable 'real.a'
> exec ^
: string
| set rebuild.a = 1
< error
| 1193 HY000 Unknown system variable 'rebuild.a'
> exec ^
: string
| set recover.a = 1
< error
| 1193 HY000 Unknown system variable 'recover.a'
> exec ^
: string
| set recursive.a = 1
< error
| 1193 HY000 Unknown system variable 'recursive.a'
> exec ^
: string
| set redo_buffer_size.a = 1
< error
| 1193 HY000 Unknown system variable 'redo_buffer_size.a'
> exec ^
: string
| set redundant.a = 1
< error
| 1193 HY000 Unknown system variable 'redundant.a'
> exec ^
: string
| set reference.a = 1
< error
| 1193 HY000 Unknown system variable 'reference.a'
> exec ^
: string
| set references.a = 1
< error
| 1193 HY000 Unknown system variable 'references.a'
> exec ^
: string
| set regexp.a = 1
< error
| 1193 HY000 Unknown system variable 'regexp.a'
> exec ^
: string
| set registration.a = 1
< error
| 1193 HY000 Unknown system variable 'registration.a'
> exec ^
: string
| set relational.a = 1
< error
| 1193 HY000 Unknown system variable 'relational.a'
> exec ^
: string
| set relay.a = 1
< error
| 1193 HY000 Unknown system variable 'relay.a'
> exec ^
: string
| set relay_log_file.a = 1
< error
| 1193 HY000 Unknown system variable 'relay_log_file.a'
> exec ^
: string
| set relay_log_pos.a = 1
< error
| 1193 HY000 Unknown system variable 'relay_log_pos.a'
> exec ^
: string
| set relay_thread.a = 1
< error
| 1193 HY000 Unknown system variable 'relay_thread.a'
> exec ^
: string
| set relaylog.a = 1
< error
| 1193 HY000 Unknown system variable 'relaylog.a'
> exec ^
: string
| set release.a = 1
< error
| 1193 HY000 Unknown system variable 'release.a'
> exec ^
: string
| set reload.a = 1
< error
| 1193 HY000 Unknown system variable 'reload.a'
> exec ^
: string
| set remove.a = 1
< error
| 1193 HY000 Unknown system variable 'remove.a'
> exec ^
: string
| set rename.a = 1
< error
| 1193 HY000 Unknown system variable 'rename.a'
> exec ^
: string
| set reorganize.a = 1
< error
| 1193 HY000 Unknown system variable 'reorganize.a'
> exec ^
: string
| set repair.a = 1
< error
| 1193 HY000 Unknown system variable 'repair.a'
> exec ^
: string
| set repeat.a = 1
< error
| 1193 HY000 Unknown system variable 'repeat.a'
> exec ^
: string
| set repeatable.a = 1
< error
| 1193 HY000 Unknown system variable 'repeatable.a'
> exec ^
: string
| set replace.a = 1
< error
| 1193 HY000 Unknown system variable 'replace.a'
> exec ^
: string
| set replica.a = 1
< error
| 1193 HY000 Unknown system variable 'replica.a'
> exec ^
: string
| set replicas.a = 1
< error
| 1193 HY000 Unknown system variable 'replicas.a'
> exec ^
: string
| set replicate_do_db.a = 1
< error
| 1193 HY000 Unknown system variable 'replicate_do_db.a'
> exec ^
: string
| set replicate_do_table.a = 1
< error
| 1193 HY000 Unknown system variable 'replicate_do_table.a'
> exec ^
: string
| set replicate_ignore_db.a = 1
< error
| 1193 HY000 Unknown system variable 'replicate_ignore_db.a'
> exec ^
: string
| set replicate_ignore_table.a = 1
< error
| 1193 HY000 Unknown system variable 'replicate_ignore_table.a'
> exec ^
: string
| set replicate_rewrite_db.a = 1
< error
| 1193 HY000 Unknown system variable 'replicate_rewrite_db.a'
> exec ^
: string
| set replicate_wild_do_table.a = 1
< error
| 1193 HY000 Unknown system variable 'replicate_wild_do_table.a'
> exec ^
: string
| set replicate_wild_ignore_table.a = 1
< error
| 1193 HY000 Unknown system variable 'replicate_wild_ignore_table.a'
> exec ^
: string
| set replication.a = 1
< error
| 1193 HY000 Unknown system variable 'replication.a'
> exec ^
: string
| set require.a = 1
< error
| 1193 HY000 Unknown system variable 'require.a'
> exec ^
: string
| set require_row_format.a = 1
< error
| 1193 HY000 Unknown system variable 'require_row_format.a'
> exec ^
: string
| set require_table_primary_key_check.a = 1
< error
| 1193 HY000 Unknown system variable 'require_table_primary_key_check.a'
> exec ^
: string
| set reset.a = 1
< error
| 1193 HY000 Unknown system variable 'reset.a'
> exec ^
: string
| set resignal.a = 1
< error
| 1193 HY000 Unknown system variable 'resignal.a'
> exec ^
: string
| set resource.a = 1
< error
| 1193 HY000 Unknown system variable 'resource.a'
> exec ^
: string
| set respect.a = 1
< error
| 1193 HY000 Unknown system variable 'respect.a'
> exec ^
: string
| set restart.a = 1
< error
| 1193 HY000 Unknown system variable 'restart.a'
> exec ^
: string
| set restore.a = 1
< error
| 1193 HY000 Unknown system variable 'restore.a'
> exec ^
: string
| set restrict.a = 1
< error
| 1193 HY000 Unknown system variable 'restrict.a'
> exec ^
: string
| set resume.a = 1
< error
| 1193 HY000 Unknown system variable 'resume.a'
> exec ^
: string
| set retain.a = 1
< error
| 1193 HY000 Unknown system variable 'retain.a'
> exec ^
: string
| set return.a = 1
< error
| 1193 HY000 Unknown system variable 'return.a'
> exec ^
: string
| set returned_sqlstate.a = 1
< error
| 1193 HY000 Unknown system variable 'returned_sqlstate.a'
> exec ^
: string
| set returning.a = 1
< error
| 1193 HY000 Unknown system variable 'returning.a'
> exec ^
: string
| set returns.a = 1
< error
| 1193 HY000 Unknown system variable 'returns.a'
> exec ^
: string
| set reuse.a = 1
< error
| 1193 HY000 Unknown system variable 'reuse.a'
> exec ^
: string
| set reverse.a = 1
< error
| 1193 HY000 Unknown system variable 'reverse.a'
> exec ^
: string
| set revoke.a = 1
< error
| 1193 HY000 Unknown system variable 'revoke.a'
> exec ^
: string
| set right.a = 1
< error
| 1193 HY000 Unknown system variable 'right.a'
> exec ^
: string
| set rlike.a = 1
< error
| 1193 HY000 Unknown system variable 'rlike.a'
> exec ^
: string
| set role.a = 1
< error
| 1193 HY000 Unknown system variable 'role.a'
> exec ^
: string
| set rollback.a = 1
< error
| 1193 HY000 Unknown system variable 'rollback.a'
> exec ^
: string
| set rollup.a = 1
< error
| 1193 HY000 Unknown system variable 'rollup.a'
> exec ^
: string
| set rotate.a = 1
< error
| 1193 HY000 Unknown system variable 'rotate.a'
> exec ^
: string
| set routine.a = 1
< error
| 1193 HY000 Unknown system variable 'routine.a'
> exec ^
: string
| set row.a = 1
< error
| 1193 HY000 Unknown system variable 'row.a'
> exec ^
: string
| set row_count.a = 1
< error
| 1193 HY000 Unknown system variable 'row_count.a'
> exec ^
: string
| set row_format.a = 1
< error
| 1193 HY000 Unknown system variable 'row_format.a'
> exec ^
: string
| set row_number.a = 1
< error
| 1193 HY000 Unknown system variable 'row_number.a'
> exec ^
: string
| set rows.a = 1
< error
| 1193 HY000 Unknown system variable 'rows.a'
> exec ^
: string
| set rtree.a = 1
< error
| 1193 HY000 Unknown system variable 'rtree.a'
> exec ^
: string
| set s3.a = 1
< error
| 1193 HY000 Unknown system variable 's3.a'
> exec ^
: string
| set savepoint.a = 1
< error
| 1193 HY000 Unknown system variable 'savepoint.a'
> exec ^
: string
| set schedule.a = 1
< error
| 1193 HY000 Unknown system variable 'schedule.a'
> exec ^
: string
| set schema.a = 1
< error
| 1193 HY000 Unknown system variable 'schema.a'
> exec ^
: string
| set schema_name.a = 1
< error
| 1193 HY000 Unknown system variable 'schema_name.a'
> exec ^
: string
| set schemas.a = 1
< error
| 1193 HY000 Unknown system variable 'schemas.a'
> exec ^
: string
| set second.a = 1
< error
| 1193 HY000 Unknown system variable 'second.a'
> exec ^
: string
| set second_microsecond.a = 1
< error
| 1193 HY000 Unknown system variable 'second_microsecond.a'
> exec ^
: string
| set secondary.a = 1
< error
| 1193 HY000 Unknown system variable 'secondary.a'
> exec ^
: string
| set secondary_engine.a = 1
< error
| 1193 HY000 Unknown system variable 'secondary_engine.a'
> exec ^
: string
| set secondary_engine_attribute.a = 1
< error
| 1193 HY000 Unknown system variable 'secondary_engine_attribute.a'
> exec ^
: string
| set secondary_load.a = 1
< error
| 1193 HY000 Unknown system variable 'secondary_load.a'
> exec ^
: string
| set secondary_unload.a = 1
< error
| 1193 HY000 Unknown system variable 'secondary_unload.a'
> exec ^
: string
| set security.a = 1
< error
| 1193 HY000 Unknown system variable 'security.a'
> exec ^
: string
| set select.a = 1
< error
| 1193 HY000 Unknown system variable 'select.a'
> exec ^
: string
| set sensitive.a = 1
< error
| 1193 HY000 Unknown system variable 'sensitive.a'
> exec ^
: string
| set separator.a = 1
< error
| 1193 HY000 Unknown system variable 'separator.a'
> exec ^
: string
| set serial.a = 1
< error
| 1193 HY000 Unknown system variable 'serial.a'
> exec ^
: string
| set serializable.a = 1
< error
| 1193 HY000 Unknown system variable 'serializable.a'
> exec ^
: string
| set server.a = 1
< error
| 1193 HY000 Unknown system variable 'server.a'
> exec ^
: string
| set session.a = 1
< error
| 1064 42000 …near 'session.a = 1' at line 1
> exec ^
: string
| set session_user.a = 1
< error
| 1193 HY000 Unknown system variable 'session_user.a'
> exec ^
: string
| set set.a = 1
< error
| 1193 HY000 Unknown system variable 'set.a'
> exec ^
: string
| set sets.a = 1
< error
| 1193 HY000 Unknown system variable 'sets.a'
> exec ^
: string
| set share.a = 1
< error
| 1193 HY000 Unknown system variable 'share.a'
> exec ^
: string
| set show.a = 1
< error
| 1193 HY000 Unknown system variable 'show.a'
> exec ^
: string
| set shutdown.a = 1
< error
| 1193 HY000 Unknown system variable 'shutdown.a'
> exec ^
: string
| set signal.a = 1
< error
| 1193 HY000 Unknown system variable 'signal.a'
> exec ^
: string
| set signed.a = 1
< error
| 1193 HY000 Unknown system variable 'signed.a'
> exec ^
: string
| set simple.a = 1
< error
| 1193 HY000 Unknown system variable 'simple.a'
> exec ^
: string
| set skip.a = 1
< error
| 1193 HY000 Unknown system variable 'skip.a'
> exec ^
: string
| set slave.a = 1
< error
| 1193 HY000 Unknown system variable 'slave.a'
> exec ^
: string
| set slow.a = 1
< error
| 1193 HY000 Unknown system variable 'slow.a'
> exec ^
: string
| set smallint.a = 1
< error
| 1193 HY000 Unknown system variable 'smallint.a'
> exec ^
: string
| set snapshot.a = 1
< error
| 1193 HY000 Unknown system variable 'snapshot.a'
> exec ^
: string
| set socket.a = 1
< error
| 1193 HY000 Unknown system variable 'socket.a'
> exec ^
: string
| set some.a = 1
< error
| 1193 HY000 Unknown system variable 'some.a'
> exec ^
: string
| set soname.a = 1
< error
| 1193 HY000 Unknown system variable 'soname.a'
> exec ^
: string
| set sounds.a = 1
< error
| 1193 HY000 Unknown system variable 'sounds.a'
> exec ^
: string
| set source.a = 1
< error
| 1193 HY000 Unknown system variable 'source.a'
> exec ^
: string
| set source_auto_position.a = 1
< error
| 1193 HY000 Unknown system variable 'source_auto_position.a'
> exec ^
: string
| set source_bind.a = 1
< error
| 1193 HY000 Unknown system variable 'source_bind.a'
> exec ^
: string
| set source_compression_algorithms.a = 1
< error
| 1193 HY000 Unknown system variable 'source_compression_algorithms.a'
> exec ^
: string
| set source_connect_retry.a = 1
< error
| 1193 HY000 Unknown system variable 'source_connect_retry.a'
> exec ^
: string
| set source_connection_auto_failover.a = 1
< error
| 1193 HY000 Unknown system variable 'source_connection_auto_failover.a'
> exec ^
: string
| set source_delay.a = 1
< error
| 1193 HY000 Unknown system variable 'source_delay.a'
> exec ^
: string
| set source_heartbeat_period.a = 1
< error
| 1193 HY000 Unknown system variable 'source_heartbeat_period.a'
> exec ^
: string
| set source_host.a = 1
< error
| 1193 HY000 Unknown system variable 'source_host.a'
> exec ^
: string
| set source_log_file.a = 1
< error
| 1193 HY000 Unknown system variable 'source_log_file.a'
> exec ^
: string
| set source_log_pos.a = 1
< error
| 1193 HY000 Unknown system variable 'source_log_pos.a'
> exec ^
: string
| set source_password.a = 1
< error
| 1193 HY000 Unknown system variable 'source_password.a'
> exec ^
: string
| set source_port.a = 1
< error
| 1193 HY000 Unknown system variable 'source_port.a'
> exec ^
: string
| set source_public_key_path.a = 1
< error
| 1193 HY000 Unknown system variable 'source_public_key_path.a'
> exec ^
: string
| set source_retry_count.a = 1
< error
| 1193 HY000 Unknown system variable 'source_retry_count.a'
> exec ^
: string
| set source_ssl.a = 1
< error
| 1193 HY000 Unknown system variable 'source_ssl.a'
> exec ^
: string
| set source_ssl_ca.a = 1
< error
| 1193 HY000 Unknown system variable 'source_ssl_ca.a'
> exec ^
: string
| set source_ssl_capath.a = 1
< error
| 1193 HY000 Unknown system variable 'source_ssl_capath.a'
> exec ^
: string
| set source_ssl_cert.a = 1
< error
| 1193 HY000 Unknown system variable 'source_ssl_cert.a'
> exec ^
: string
| set source_ssl_cipher.a = 1
< error
| 1193 HY000 Unknown system variable 'source_ssl_cipher.a'
> exec ^
: string
| set source_ssl_crl.a = 1
< error
| 1193 HY000 Unknown system variable 'source_ssl_crl.a'
> exec ^
: string
| set source_ssl_crlpath.a = 1
< error
| 1193 HY000 Unknown system variable 'source_ssl_crlpath.a'
> exec ^
: string
| set source_ssl_key.a = 1
< error
| 1193 HY000 Unknown system variable 'source_ssl_key.a'
> exec ^
: string
| set source_ssl_verify_server_cert.a = 1
< error
| 1193 HY000 Unknown system variable 'source_ssl_verify_server_cert.a'
> exec ^
: string
| set source_tls_ciphersuites.a = 1
< error
| 1193 HY000 Unknown system variable 'source_tls_ciphersuites.a'
> exec ^
: string
| set source_tls_version.a = 1
< error
| 1193 HY000 Unknown system variable 'source_tls_version.a'
> exec ^
: string
| set source_user.a = 1
< error
| 1193 HY000 Unknown system variable 'source_user.a'
> exec ^
: string
| set source_zstd_compression_level.a = 1
< error
| 1193 HY000 Unknown system variable 'source_zstd_compression_level.a'
> exec ^
: string
| set spatial.a = 1
< error
| 1193 HY000 Unknown system variable 'spatial.a'
> exec ^
: string
| set specific.a = 1
< error
| 1193 HY000 Unknown system variable 'specific.a'
> exec ^
: string
| set sql.a = 1
< error
| 1193 HY000 Unknown system variable 'sql.a'
> exec ^
: string
| set sql_after_gtids.a = 1
< error
| 1193 HY000 Unknown system variable 'sql_after_gtids.a'
> exec ^
: string
| set sql_after_mts_gaps.a = 1
< error
| 1193 HY000 Unknown system variable 'sql_after_mts_gaps.a'
> exec ^
: string
| set sql_before_gtids.a = 1
< error
| 1193 HY000 Unknown system variable 'sql_before_gtids.a'
> exec ^
: string
| set sql_big_result.a = 1
< error
| 1193 HY000 Unknown system variable 'sql_big_result.a'
> exec ^
: string
| set sql_buffer_result.a = 1
< error
| 1193 HY000 Unknown system variable 'sql_buffer_result.a'
> exec ^
: string
| set sql_calc_found_rows.a = 1
< error
| 1193 HY000 Unknown system variable 'sql_calc_found_rows.a'
> exec ^
: string
| set sql_no_cache.a = 1
< error
| 1193 HY000 Unknown system variable 'sql_no_cache.a'
> exec ^
: string
| set sql_small_result.a = 1
< error
| 1193 HY000 Unknown system variable 'sql_small_result.a'
> exec ^
: string
| set sql_thread.a = 1
< error
| 1193 HY000 Unknown system variable 'sql_thread.a'
> exec ^
: string
| set sql_tsi_day.a = 1
< error
| 1193 HY000 Unknown system variable 'sql_tsi_day.a'
> exec ^
: string
| set sql_tsi_hour.a = 1
< error
| 1193 HY000 Unknown system variable 'sql_tsi_hour.a'
> exec ^
: string
| set sql_tsi_minute.a = 1
< error
| 1193 HY000 Unknown system variable 'sql_tsi_minute.a'
> exec ^
: string
| set sql_tsi_month.a = 1
< error
| 1193 HY000 Unknown system variable 'sql_tsi_month.a'
> exec ^
: string
| set sql_tsi_quarter.a = 1
< error
| 1193 HY000 Unknown system variable 'sql_tsi_quarter.a'
> exec ^
: string
| set sql_tsi_second.a = 1
< error
| 1193 HY000 Unknown system variable 'sql_tsi_second.a'
> exec ^
: string
| set sql_tsi_week.a = 1
< error
| 1193 HY000 Unknown system variable 'sql_tsi_week.a'
> exec ^
: string
| set sql_tsi_year.a = 1
< error
| 1193 HY000 Unknown system variable 'sql_tsi_year.a'
> exec ^
: string
| set sqlexception.a = 1
< error
| 1193 HY000 Unknown system variable 'sqlexception.a'
> exec ^
: string
| set sqlstate.a = 1
< error
| 1193 HY000 Unknown system variable 'sqlstate.a'
> exec ^
: string
| set sqlwarning.a = 1
< error
| 1193 HY000 Unknown system variable 'sqlwarning.a'
> exec ^
: string
| set srid.a = 1
< error
| 1193 HY000 Unknown system variable 'srid.a'
> exec ^
: string
| set sse_kms.a = 1
< error
| 1193 HY000 Unknown system variable 'sse_kms.a'
> exec ^
: string
| set sse_s3.a = 1
< error
| 1193 HY000 Unknown system variable 'sse_s3.a'
> exec ^
: string
| set ssl.a = 1
< error
| 1193 HY000 Unknown system variable 'ssl.a'
> exec ^
: string
| set st_collect.a = 1
< error
| 1193 HY000 Unknown system variable 'st_collect.a'
> exec ^
: string
| set stacked.a = 1
< error
| 1193 HY000 Unknown system variable 'stacked.a'
> exec ^
: string
| set start.a = 1
< error
| 1193 HY000 Unknown system variable 'start.a'
> exec ^
: string
| set starting.a = 1
< error
| 1193 HY000 Unknown system variable 'starting.a'
> exec ^
: string
| set starts.a = 1
< error
| 1193 HY000 Unknown system variable 'starts.a'
> exec ^
: string
| set stats_auto_recalc.a = 1
< error
| 1193 HY000 Unknown system variable 'stats_auto_recalc.a'
> exec ^
: string
| set stats_persistent.a = 1
< error
| 1193 HY000 Unknown system variable 'stats_persistent.a'
> exec ^
: string
| set stats_sample_pages.a = 1
< error
| 1193 HY000 Unknown system variable 'stats_sample_pages.a'
> exec ^
: string
| set status.a = 1
< error
| 1193 HY000 Unknown system variable 'status.a'
> exec ^
: string
| set std.a = 1
< error
| 1193 HY000 Unknown system variable 'std.a'
> exec ^
: string
| set stddev.a = 1
< error
| 1193 HY000 Unknown system variable 'stddev.a'
> exec ^
: string
| set stddev_pop.a = 1
< error
| 1193 HY000 Unknown system variable 'stddev_pop.a'
> exec ^
: string
| set stddev_samp.a = 1
< error
| 1193 HY000 Unknown system variable 'stddev_samp.a'
> exec ^
: string
| set stop.a = 1
< error
| 1193 HY000 Unknown system variable 'stop.a'
> exec ^
: string
| set storage.a = 1
< error
| 1193 HY000 Unknown system variable 'storage.a'
> exec ^
: string
| set stored.a = 1
< error
| 1193 HY000 Unknown system variable 'stored.a'
> exec ^
: string
| set straight_join.a = 1
< error
| 1193 HY000 Unknown system variable 'straight_join.a'
> exec ^
: string
| set stream.a = 1
< error
| 1193 HY000 Unknown system variable 'stream.a'
> exec ^
: string
| set strict_load.a = 1
< error
| 1193 HY000 Unknown system variable 'strict_load.a'
> exec ^
: string
| set string.a = 1
< error
| 1193 HY000 Unknown system variable 'string.a'
> exec ^
: string
| set subclass_origin.a = 1
< error
| 1193 HY000 Unknown system variable 'subclass_origin.a'
> exec ^
: string
| set subdate.a = 1
< error
| 1193 HY000 Unknown system variable 'subdate.a'
> exec ^
: string
| set subject.a = 1
< error
| 1193 HY000 Unknown system variable 'subject.a'
> exec ^
: string
| set subpartition.a = 1
< error
| 1193 HY000 Unknown system variable 'subpartition.a'
> exec ^
: string
| set subpartitions.a = 1
< error
| 1193 HY000 Unknown system variable 'subpartitions.a'
> exec ^
: string
| set substr.a = 1
< error
| 1193 HY000 Unknown system variable 'substr.a'
> exec ^
: string
| set substring.a = 1
< error
| 1193 HY000 Unknown system variable 'substring.a'
> exec ^
: string
| set sum.a = 1
< error
| 1193 HY000 Unknown system variable 'sum.a'
> exec ^
: string
| set super.a = 1
< error
| 1193 HY000 Unknown system variable 'super.a'
> exec ^
: string
| set suspend.a = 1
< error
| 1193 HY000 Unknown system variable 'suspend.a'
> exec ^
: string
| set swaps.a = 1
< error
| 1193 HY000 Unknown system variable 'swaps.a'
> exec ^
: string
| set switches.a = 1
< error
| 1193 HY000 Unknown system variable 'switches.a'
> exec ^
: string
| set sysdate.a = 1
< error
| 1193 HY000 Unknown system variable 'sysdate.a'
> exec ^
: string
| set system.a = 1
< error
| 1193 HY000 Unknown system variable 'system.a'
> exec ^
: string
| set system_user.a = 1
< error
| 1193 HY000 Unknown system variable 'system_user.a'
> exec ^
: string
| set table.a = 1
< error
| 1193 HY000 Unknown system variable 'table.a'
> exec ^
: string
| set table_checksum.a = 1
< error
| 1193 HY000 Unknown system variable 'table_checksum.a'
> exec ^
: string
| set table_name.a = 1
< error
| 1193 HY000 Unknown system variable 'table_name.a'
> exec ^
: string
| set tables.a = 1
< error
| 1193 HY000 Unknown system variable 'tables.a'
> exec ^
: string
| set tablesample.a = 1
< error
| 1193 HY000 Unknown system variable 'tablesample.a'
> exec ^
: string
| set tablespace.a = 1
< error
| 1193 HY000 Unknown system variable 'tablespace.a'
> exec ^
: string
| set temporary.a = 1
< error
| 1193 HY000 Unknown system variable 'temporary.a'
> exec ^
: string
| set temptable.a = 1
< error
| 1193 HY000 Unknown system variable 'temptable.a'
> exec ^
: string
| set terminated.a = 1
< error
| 1193 HY000 Unknown system variable 'terminated.a'
> exec ^
: string
| set text.a = 1
< error
| 1193 HY000 Unknown system variable 'text.a'
> exec ^
: string
| set than.a = 1
< error
| 1193 HY000 Unknown system variable 'than.a'
> exec ^
: string
| set then.a = 1
< error
| 1193 HY000 Unknown system variable 'then.a'
> exec ^
: string
| set thread_priority.a = 1
< error
| 1193 HY000 Unknown system variable 'thread_priority.a'
> exec ^
: string
| set ties.a = 1
< error
| 1193 HY000 Unknown system variable 'ties.a'
> exec ^
: string
| set time.a = 1
< error
| 1193 HY000 Unknown system variable 'time.a'
> exec ^
: string
| set timestamp.a = 1
< error
| 1193 HY000 Unknown system variable 'timestamp.a'
> exec ^
: string
| set timestampadd.a = 1
< error
| 1193 HY000 Unknown system variable 'timestampadd.a'
> exec ^
: string
| set timestampdiff.a = 1
< error
| 1193 HY000 Unknown system variable 'timestampdiff.a'
> exec ^
: string
| set tinyblob.a = 1
< error
| 1193 HY000 Unknown system variable 'tinyblob.a'
> exec ^
: string
| set tinyint.a = 1
< error
| 1193 HY000 Unknown system variable 'tinyint.a'
> exec ^
: string
| set tinytext.a = 1
< error
| 1193 HY000 Unknown system variable 'tinytext.a'
> exec ^
: string
| set tls.a = 1
< error
| 1193 HY000 Unknown system variable 'tls.a'
> exec ^
: string
| set to.a = 1
< error
| 1193 HY000 Unknown system variable 'to.a'
> exec ^
: string
| set trailing.a = 1
< error
| 1193 HY000 Unknown system variable 'trailing.a'
> exec ^
: string
| set transaction.a = 1
< error
| 1193 HY000 Unknown system variable 'transaction.a'
> exec ^
: string
| set trigger.a = 1
< error
| 1193 HY000 Unknown system variable 'trigger.a'
> exec ^
: string
| set triggers.a = 1
< error
| 1193 HY000 Unknown system variable 'triggers.a'
> exec ^
: string
| set trim.a = 1
< error
| 1193 HY000 Unknown system variable 'trim.a'
> exec ^
: string
| set true.a = 1
< error
| 1193 HY000 Unknown system variable 'true.a'
> exec ^
: string
| set truncate.a = 1
< error
| 1193 HY000 Unknown system variable 'truncate.a'
> exec ^
: string
| set type.a = 1
< error
| 1193 HY000 Unknown system variable 'type.a'
> exec ^
: string
| set types.a = 1
< error
| 1193 HY000 Unknown system variable 'types.a'
> exec ^
: string
| set unbounded.a = 1
< error
| 1193 HY000 Unknown system variable 'unbounded.a'
> exec ^
: string
| set uncommitted.a = 1
< error
| 1193 HY000 Unknown system variable 'uncommitted.a'
> exec ^
: string
| set undefined.a = 1
< error
| 1193 HY000 Unknown system variable 'undefined.a'
> exec ^
: string
| set undo.a = 1
< error
| 1193 HY000 Unknown system variable 'undo.a'
> exec ^
: string
| set undo_buffer_size.a = 1
< error
| 1193 HY000 Unknown system variable 'undo_buffer_size.a'
> exec ^
: string
| set undofile.a = 1
< error
| 1193 HY000 Unknown system variable 'undofile.a'
> exec ^
: string
| set unicode.a = 1
< error
| 1193 HY000 Unknown system variable 'unicode.a'
> exec ^
: string
| set uninstall.a = 1
< error
| 1193 HY000 Unknown system variable 'uninstall.a'
> exec ^
: string
| set union.a = 1
< error
| 1193 HY000 Unknown system variable 'union.a'
> exec ^
: string
| set unique.a = 1
< error
| 1193 HY000 Unknown system variable 'unique.a'
> exec ^
: string
| set unknown.a = 1
< error
| 1193 HY000 Unknown system variable 'unknown.a'
> exec ^
: string
| set unlock.a = 1
< error
| 1193 HY000 Unknown system variable 'unlock.a'
> exec ^
: string
| set unregister.a = 1
< error
| 1193 HY000 Unknown system variable 'unregister.a'
> exec ^
: string
| set unsigned.a = 1
< error
| 1193 HY000 Unknown system variable 'unsigned.a'
> exec ^
: string
| set until.a = 1
< error
| 1193 HY000 Unknown system variable 'until.a'
> exec ^
: string
| set update.a = 1
< error
| 1193 HY000 Unknown system variable 'update.a'
> exec ^
: string
| set upgrade.a = 1
< error
| 1193 HY000 Unknown system variable 'upgrade.a'
> exec ^
: string
| set uri.a = 1
< error
| 1193 HY000 Unknown system variable 'uri.a'
> exec ^
: string
| set url.a = 1
< error
| 1193 HY000 Unknown system variable 'url.a'
> exec ^
: string
| set usage.a = 1
< error
| 1193 HY000 Unknown system variable 'usage.a'
> exec ^
: string
| set use.a = 1
< error
| 1193 HY000 Unknown system variable 'use.a'
> exec ^
: string
| set use_frm.a = 1
< error
| 1193 HY000 Unknown system variable 'use_frm.a'
> exec ^
: string
| set user.a = 1
< error
| 1193 HY000 Unknown system variable 'user.a'
> exec ^
: string
| set user_resources.a = 1
< error
| 1193 HY000 Unknown system variable 'user_resources.a'
> exec ^
: string
| set using.a = 1
< error
| 1193 HY000 Unknown system variable 'using.a'
> exec ^
: string
| set utc_date.a = 1
< error
| 1193 HY000 Unknown system variable 'utc_date.a'
> exec ^
: string
| set utc_time.a = 1
< error
| 1193 HY000 Unknown system variable 'utc_time.a'
> exec ^
: string
| set utc_timestamp.a = 1
< error
| 1193 HY000 Unknown system variable 'utc_timestamp.a'
> exec ^
: string
| set validate.a = 1
< error
| 1193 HY000 Unknown system variable 'validate.a'
> exec ^
: string
| set validation.a = 1
< error
| 1193 HY000 Unknown system variable 'validation.a'
> exec ^
: string
| set value.a = 1
< error
| 1193 HY000 Unknown system variable 'value.a'
> exec ^
: string
| set values.a = 1
< error
| 1193 HY000 Unknown system variable 'values.a'
> exec ^
: string
| set var_pop.a = 1
< error
| 1193 HY000 Unknown system variable 'var_pop.a'
> exec ^
: string
| set var_samp.a = 1
< error
| 1193 HY000 Unknown system variable 'var_samp.a'
> exec ^
: string
| set varbinary.a = 1
< error
| 1193 HY000 Unknown system variable 'varbinary.a'
> exec ^
: string
| set varchar.a = 1
< error
| 1193 HY000 Unknown system variable 'varchar.a'
> exec ^
: string
| set varcharacter.a = 1
< error
| 1193 HY000 Unknown system variable 'varcharacter.a'
> exec ^
: string
| set variables.a = 1
< error
| 1193 HY000 Unknown system variable 'variables.a'
> exec ^
: string
| set variance.a = 1
< error
| 1193 HY000 Unknown system variable 'variance.a'
> exec ^
: string
| set varying.a = 1
< error
| 1193 HY000 Unknown system variable 'varying.a'
> exec ^
: string
| set vcpu.a = 1
< error
| 1193 HY000 Unknown system variable 'vcpu.a'
> exec ^
: string
| set vector.a = 1
< error
| 1193 HY000 Unknown system variable 'vector.a'
> exec ^
: string
| set verify_key_constraints.a = 1
< error
| 1193 HY000 Unknown system variable 'verify_key_constraints.a'
> exec ^
: string
| set view.a = 1
< error
| 1193 HY000 Unknown system variable 'view.a'
> exec ^
: string
| set virtual.a = 1
< error
| 1193 HY000 Unknown system variable 'virtual.a'
> exec ^
: string
| set visible.a = 1
< error
| 1193 HY000 Unknown system variable 'visible.a'
> exec ^
: string
| set wait.a = 1
< error
| 1193 HY000 Unknown system variable 'wait.a'
> exec ^
: string
| set warnings.a = 1
< error
| 1193 HY000 Unknown system variable 'warnings.a'
> exec ^
: string
| set week.a = 1
< error
| 1193 HY000 Unknown system variable 'week.a'
> exec ^
: string
| set weight_string.a = 1
< error
| 1193 HY000 Unknown system variable 'weight_string.a'
> exec ^
: string
| set when.a = 1
< error
| 1193 HY000 Unknown system variable 'when.a'
> exec ^
: string
| set where.a = 1
< error
| 1193 HY000 Unknown system variable 'where.a'
> exec ^
: string
| set while.a = 1
< error
| 1193 HY000 Unknown system variable 'while.a'
> exec ^
: string
| set window.a = 1
< error
| 1193 HY000 Unknown system variable 'window.a'
> exec ^
: string
| set with.a = 1
< error
| 1193 HY000 Unknown system variable 'with.a'
> exec ^
: string
| set without.a = 1
< error
| 1193 HY000 Unknown system variable 'without.a'
> exec ^
: string
| set work.a = 1
< error
| 1193 HY000 Unknown system variable 'work.a'
> exec ^
: string
| set wrapper.a = 1
< error
| 1193 HY000 Unknown system variable 'wrapper.a'
> exec ^
: string
| set write.a = 1
< error
| 1193 HY000 Unknown system variable 'write.a'
> exec ^
: string
| set x509.a = 1
< error
| 1193 HY000 Unknown system variable 'x509.a'
> exec ^
: string
| set xa.a = 1
< error
| 1193 HY000 Unknown system variable 'xa.a'
> exec ^
: string
| set xid.a = 1
< error
| 1193 HY000 Unknown system variable 'xid.a'
> exec ^
: string
| set xml.a = 1
< error
| 1193 HY000 Unknown system variable 'xml.a'
> exec ^
: string
| set xor.a = 1
< error
| 1193 HY000 Unknown system variable 'xor.a'
> exec ^
: string
| set year.a = 1
< error
| 1193 HY000 Unknown system variable 'year.a'
> exec ^
: string
| set year_month.a = 1
< error
| 1193 HY000 Unknown system variable 'year_month.a'
> exec ^
: string
| set zerofill.a = 1
< error
| 1193 HY000 Unknown system variable 'zerofill.a'
> exec ^
: string
| set zone.a = 1
< error
| 1193 HY000 Unknown system variable 'zone.a'
> exec ^
: string
| set a.absent = 1
< error
| 1193 HY000 Unknown system variable 'a.absent'
> exec ^
: string
| set a.accessible = 1
< error
| 1193 HY000 Unknown system variable 'a.accessible'
> exec ^
: string
| set a.account = 1
< error
| 1193 HY000 Unknown system variable 'a.account'
> exec ^
: string
| set a.action = 1
< error
| 1193 HY000 Unknown system variable 'a.action'
> exec ^
: string
| set a.active = 1
< error
| 1193 HY000 Unknown system variable 'a.active'
> exec ^
: string
| set a.add = 1
< error
| 1193 HY000 Unknown system variable 'a.add'
> exec ^
: string
| set a.adddate = 1
< error
| 1193 HY000 Unknown system variable 'a.adddate'
> exec ^
: string
| set a.admin = 1
< error
| 1193 HY000 Unknown system variable 'a.admin'
> exec ^
: string
| set a.after = 1
< error
| 1193 HY000 Unknown system variable 'a.after'
> exec ^
: string
| set a.against = 1
< error
| 1193 HY000 Unknown system variable 'a.against'
> exec ^
: string
| set a.aggregate = 1
< error
| 1193 HY000 Unknown system variable 'a.aggregate'
> exec ^
: string
| set a.algorithm = 1
< error
| 1193 HY000 Unknown system variable 'a.algorithm'
> exec ^
: string
| set a.all = 1
< error
| 1193 HY000 Unknown system variable 'a.all'
> exec ^
: string
| set a.allow_missing_files = 1
< error
| 1193 HY000 Unknown system variable 'a.allow_missing_files'
> exec ^
: string
| set a.alter = 1
< error
| 1193 HY000 Unknown system variable 'a.alter'
> exec ^
: string
| set a.always = 1
< error
| 1193 HY000 Unknown system variable 'a.always'
> exec ^
: string
| set a.analyze = 1
< error
| 1193 HY000 Unknown system variable 'a.analyze'
> exec ^
: string
| set a.and = 1
< error
| 1193 HY000 Unknown system variable 'a.and'
> exec ^
: string
| set a.any = 1
< error
| 1193 HY000 Unknown system variable 'a.any'
> exec ^
: string
| set a.array = 1
< error
| 1193 HY000 Unknown system variable 'a.array'
> exec ^
: string
| set a.as = 1
< error
| 1193 HY000 Unknown system variable 'a.as'
> exec ^
: string
| set a.asc = 1
< error
| 1193 HY000 Unknown system variable 'a.asc'
> exec ^
: string
| set a.ascii = 1
< error
| 1193 HY000 Unknown system variable 'a.ascii'
> exec ^
: string
| set a.asensitive = 1
< error
| 1193 HY000 Unknown system variable 'a.asensitive'
> exec ^
: string
| set a.assign_gtids_to_anonymous_transactions = 1
< error
| 1193 HY000 Unknown system variable 'a.assign_gtids_to_anonymous_transactions'
> exec ^
: string
| set a.at = 1
< error
| 1193 HY000 Unknown system variable 'a.at'
> exec ^
: string
| set a.attribute = 1
< error
| 1193 HY000 Unknown system variable 'a.attribute'
> exec ^
: string
| set a.authentication = 1
< error
| 1193 HY000 Unknown system variable 'a.authentication'
> exec ^
: string
| set a.auto = 1
< error
| 1193 HY000 Unknown system variable 'a.auto'
> exec ^
: string
| set a.auto_increment = 1
< error
| 1193 HY000 Unknown system variable 'a.auto_increment'
> exec ^
: string
| set a.auto_refresh = 1
< error
| 1193 HY000 Unknown system variable 'a.auto_refresh'
> exec ^
: string
| set a.auto_refresh_source = 1
< error
| 1193 HY000 Unknown system variable 'a.auto_refresh_source'
> exec ^
: string
| set a.autoextend_size = 1
< error
| 1193 HY000 Unknown system variable 'a.autoextend_size'
> exec ^
: string
| set a.avg = 1
< error
| 1193 HY000 Unknown system variable 'a.avg'
> exec ^
: string
| set a.avg_row_length = 1
< error
| 1193 HY000 Unknown system variable 'a.avg_row_length'
> exec ^
: string
| set a.backup = 1
< error
| 1193 HY000 Unknown system variable 'a.backup'
> exec ^
: string
| set a.before = 1
< error
| 1193 HY000 Unknown system variable 'a.before'
> exec ^
: string
| set a.begin = 1
< error
| 1193 HY000 Unknown system variable 'a.begin'
> exec ^
: string
| set a.bernoulli = 1
< error
| 1193 HY000 Unknown system variable 'a.bernoulli'
> exec ^
: string
| set a.between = 1
< error
| 1193 HY000 Unknown system variable 'a.between'
> exec ^
: string
| set a.bigint = 1
< error
| 1193 HY000 Unknown system variable 'a.bigint'
> exec ^
: string
| set a.binary = 1
< error
| 1193 HY000 Unknown system variable 'a.binary'
> exec ^
: string
| set a.binlog = 1
< error
| 1193 HY000 Unknown system variable 'a.binlog'
> exec ^
: string
| set a.bit = 1
< error
| 1193 HY000 Unknown system variable 'a.bit'
> exec ^
: string
| set a.bit_and = 1
< error
| 1193 HY000 Unknown system variable 'a.bit_and'
> exec ^
: string
| set a.bit_or = 1
< error
| 1193 HY000 Unknown system variable 'a.bit_or'
> exec ^
: string
| set a.bit_xor = 1
< error
| 1193 HY000 Unknown system variable 'a.bit_xor'
> exec ^
: string
| set a.blob = 1
< error
| 1193 HY000 Unknown system variable 'a.blob'
> exec ^
: string
| set a.block = 1
< error
| 1193 HY000 Unknown system variable 'a.block'
> exec ^
: string
| set a.bool = 1
< error
| 1193 HY000 Unknown system variable 'a.bool'
> exec ^
: string
| set a.boolean = 1
< error
| 1193 HY000 Unknown system variable 'a.boolean'
> exec ^
: string
| set a.both = 1
< error
| 1193 HY000 Unknown system variable 'a.both'
> exec ^
: string
| set a.btree = 1
< error
| 1193 HY000 Unknown system variable 'a.btree'
> exec ^
: string
| set a.buckets = 1
< error
| 1193 HY000 Unknown system variable 'a.buckets'
> exec ^
: string
| set a.bulk = 1
< error
| 1193 HY000 Unknown system variable 'a.bulk'
> exec ^
: string
| set a.by = 1
< error
| 1193 HY000 Unknown system variable 'a.by'
> exec ^
: string
| set a.byte = 1
< error
| 1193 HY000 Unknown system variable 'a.byte'
> exec ^
: string
| set a.cache = 1
< error
| 1193 HY000 Unknown system variable 'a.cache'
> exec ^
: string
| set a.call = 1
< error
| 1193 HY000 Unknown system variable 'a.call'
> exec ^
: string
| set a.cascade = 1
< error
| 1193 HY000 Unknown system variable 'a.cascade'
> exec ^
: string
| set a.cascaded = 1
< error
| 1193 HY000 Unknown system variable 'a.cascaded'
> exec ^
: string
| set a.case = 1
< error
| 1193 HY000 Unknown system variable 'a.case'
> exec ^
: string
| set a.cast = 1
< error
| 1193 HY000 Unknown system variable 'a.cast'
> exec ^
: string
| set a.catalog_name = 1
< error
| 1193 HY000 Unknown system variable 'a.catalog_name'
> exec ^
: string
| set a.chain = 1
< error
| 1193 HY000 Unknown system variable 'a.chain'
> exec ^
: string
| set a.challenge_response = 1
< error
| 1193 HY000 Unknown system variable 'a.challenge_response'
> exec ^
: string
| set a.change = 1
< error
| 1193 HY000 Unknown system variable 'a.change'
> exec ^
: string
| set a.changed = 1
< error
| 1193 HY000 Unknown system variable 'a.changed'
> exec ^
: string
| set a.channel = 1
< error
| 1193 HY000 Unknown system variable 'a.channel'
> exec ^
: string
| set a.char = 1
< error
| 1193 HY000 Unknown system variable 'a.char'
> exec ^
: string
| set a.character = 1
< error
| 1193 HY000 Unknown system variable 'a.character'
> exec ^
: string
| set a.charset = 1
< error
| 1193 HY000 Unknown system variable 'a.charset'
> exec ^
: string
| set a.check = 1
< error
| 1193 HY000 Unknown system variable 'a.check'
> exec ^
: string
| set a.checksum = 1
< error
| 1193 HY000 Unknown system variable 'a.checksum'
> exec ^
: string
| set a.cipher = 1
< error
| 1193 HY000 Unknown system variable 'a.cipher'
> exec ^
: string
| set a.class_origin = 1
< error
| 1193 HY000 Unknown system variable 'a.class_origin'
> exec ^
: string
| set a.client = 1
< error
| 1193 HY000 Unknown system variable 'a.client'
> exec ^
: string
| set a.clone = 1
< error
| 1193 HY000 Unknown system variable 'a.clone'
> exec ^
: string
| set a.close = 1
< error
| 1193 HY000 Unknown system variable 'a.close'
> exec ^
: string
| set a.coalesce = 1
< error
| 1193 HY000 Unknown system variable 'a.coalesce'
> exec ^
: string
| set a.code = 1
< error
| 1193 HY000 Unknown system variable 'a.code'
> exec ^
: string
| set a.collate = 1
< error
| 1193 HY000 Unknown system variable 'a.collate'
> exec ^
: string
| set a.collation = 1
< error
| 1193 HY000 Unknown system variable 'a.collation'
> exec ^
: string
| set a.column = 1
< error
| 1193 HY000 Unknown system variable 'a.column'
> exec ^
: string
| set a.column_format = 1
< error
| 1193 HY000 Unknown system variable 'a.column_format'
> exec ^
: string
| set a.column_name = 1
< error
| 1193 HY000 Unknown system variable 'a.column_name'
> exec ^
: string
| set a.columns = 1
< error
| 1193 HY000 Unknown system variable 'a.columns'
> exec ^
: string
| set a.comment = 1
< error
| 1193 HY000 Unknown system variable 'a.comment'
> exec ^
: string
| set a.commit = 1
< error
| 1193 HY000 Unknown system variable 'a.commit'
> exec ^
: string
| set a.committed = 1
< error
| 1193 HY000 Unknown system variable 'a.committed'
> exec ^
: string
| set a.compact = 1
< error
| 1193 HY000 Unknown system variable 'a.compact'
> exec ^
: string
| set a.completion = 1
< error
| 1193 HY000 Unknown system variable 'a.completion'
> exec ^
: string
| set a.component = 1
< error
| 1193 HY000 Unknown system variable 'a.component'
> exec ^
: string
| set a.compressed = 1
< error
| 1193 HY000 Unknown system variable 'a.compressed'
> exec ^
: string
| set a.compression = 1
< error
| 1193 HY000 Unknown system variable 'a.compression'
> exec ^
: string
| set a.concat = 1
< error
| 1193 HY000 Unknown system variable 'a.concat'
> exec ^
: string
| set a.concurrent = 1
< error
| 1193 HY000 Unknown system variable 'a.concurrent'
> exec ^
: string
| set a.condition = 1
< error
| 1193 HY000 Unknown system variable 'a.condition'
> exec ^
: string
| set a.connection = 1
< error
| 1193 HY000 Unknown system variable 'a.connection'
> exec ^
: string
| set a.consistent = 1
< error
| 1193 HY000 Unknown system variable 'a.consistent'
> exec ^
: string
| set a.constraint = 1
< error
| 1193 HY000 Unknown system variable 'a.constraint'
> exec ^
: string
| set a.constraint_catalog = 1
< error
| 1193 HY000 Unknown system variable 'a.constraint_catalog'
> exec ^
: string
| set a.constraint_name = 1
< error
| 1193 HY000 Unknown system variable 'a.constraint_name'
> exec ^
: string
| set a.constraint_schema = 1
< error
| 1193 HY000 Unknown system variable 'a.constraint_schema'
> exec ^
: string
| set a.contains = 1
< error
| 1193 HY000 Unknown system variable 'a.contains'
> exec ^
: string
| set a.context = 1
< error
| 1193 HY000 Unknown system variable 'a.context'
> exec ^
: string
| set a.continue = 1
< error
| 1193 HY000 Unknown system variable 'a.continue'
> exec ^
: string
| set a.convert = 1
< error
| 1193 HY000 Unknown system variable 'a.convert'
> exec ^
: string
| set a.count = 1
< error
| 1193 HY000 Unknown system variable 'a.count'
> exec ^
: string
| set a.cpu = 1
< error
| 1193 HY000 Unknown system variable 'a.cpu'
> exec ^
: string
| set a.create = 1
< error
| 1193 HY000 Unknown system variable 'a.create'
> exec ^
: string
| set a.cross = 1
< error
| 1193 HY000 Unknown system variable 'a.cross'
> exec ^
: string
| set a.csv = 1
< error
| 1193 HY000 Unknown system variable 'a.csv'
> exec ^
: string
| set a.cube = 1
< error
| 1193 HY000 Unknown system variable 'a.cube'
> exec ^
: string
| set a.cume_dist = 1
< error
| 1193 HY000 Unknown system variable 'a.cume_dist'
> exec ^
: string
| set a.curdate = 1
< error
| 1193 HY000 Unknown system variable 'a.curdate'
> exec ^
: string
| set a.current = 1
< error
| 1193 HY000 Unknown system variable 'a.current'
> exec ^
: string
| set a.current_date = 1
< error
| 1193 HY000 Unknown system variable 'a.current_date'
> exec ^
: string
| set a.current_time = 1
< error
| 1193 HY000 Unknown system variable 'a.current_time'
> exec ^
: string
| set a.current_timestamp = 1
< error
| 1193 HY000 Unknown system variable 'a.current_timestamp'
> exec ^
: string
| set a.current_user = 1
< error
| 1193 HY000 Unknown system variable 'a.current_user'
> exec ^
: string
| set a.cursor = 1
< error
| 1193 HY000 Unknown system variable 'a.cursor'
> exec ^
: string
| set a.cursor_name = 1
< error
| 1193 HY000 Unknown system variable 'a.cursor_name'
> exec ^
: string
| set a.curtime = 1
< error
| 1193 HY000 Unknown system variable 'a.curtime'
> exec ^
: string
| set a.data = 1
< error
| 1193 HY000 Unknown system variable 'a.data'
> exec ^
: string
| set a.database = 1
< error
| 1193 HY000 Unknown system variable 'a.database'
> exec ^
: string
| set a.databases = 1
< error
| 1193 HY000 Unknown system variable 'a.databases'
> exec ^
: string
| set a.datafile = 1
< error
| 1193 HY000 Unknown system variable 'a.datafile'
> exec ^
: string
| set a.date = 1
< error
| 1193 HY000 Unknown system variable 'a.date'
> exec ^
: string
| set a.date_add = 1
< error
| 1193 HY000 Unknown system variable 'a.date_add'
> exec ^
: string
| set a.date_sub = 1
< error
| 1193 HY000 Unknown system variable 'a.date_sub'
> exec ^
: string
| set a.datetime = 1
< error
| 1193 HY000 Unknown system variable 'a.datetime'
> exec ^
: string
| set a.day = 1
< error
| 1193 HY000 Unknown system variable 'a.day'
> exec ^
: string
| set a.day_hour = 1
< error
| 1193 HY000 Unknown system variable 'a.day_hour'
> exec ^
: string
| set a.day_microsecond = 1
< error
| 1193 HY000 Unknown system variable 'a.day_microsecond'
> exec ^
: string
| set a.day_minute = 1
< error
| 1193 HY000 Unknown system variable 'a.day_minute'
> exec ^
: string
| set a.day_second = 1
< error
| 1193 HY000 Unknown system variable 'a.day_second'
> exec ^
: string
| set a.deallocate = 1
< error
| 1193 HY000 Unknown system variable 'a.deallocate'
> exec ^
: string
| set a.dec = 1
< error
| 1193 HY000 Unknown system variable 'a.dec'
> exec ^
: string
| set a.decimal = 1
< error
| 1193 HY000 Unknown system variable 'a.decimal'
> exec ^
: string
| set a.declare = 1
< error
| 1193 HY000 Unknown system variable 'a.declare'
> exec ^
: string
| set a.default = 1
< error
| 1193 HY000 Unknown system variable 'a.default'
> exec ^
: string
| set a.default_auth = 1
< error
| 1193 HY000 Unknown system variable 'a.default_auth'
> exec ^
: string
| set a.definer = 1
< error
| 1193 HY000 Unknown system variable 'a.definer'
> exec ^
: string
| set a.definition = 1
< error
| 1193 HY000 Unknown system variable 'a.definition'
> exec ^
: string
| set a.delay_key_write = 1
< error
| 1193 HY000 Unknown system variable 'a.delay_key_write'
> exec ^
: string
| set a.delayed = 1
< error
| 1193 HY000 Unknown system variable 'a.delayed'
> exec ^
: string
| set a.delete = 1
< error
| 1193 HY000 Unknown system variable 'a.delete'
> exec ^
: string
| set a.dense_rank = 1
< error
| 1193 HY000 Unknown system variable 'a.dense_rank'
> exec ^
: string
| set a.desc = 1
< error
| 1193 HY000 Unknown system variable 'a.desc'
> exec ^
: string
| set a.describe = 1
< error
| 1193 HY000 Unknown system variable 'a.describe'
> exec ^
: string
| set a.description = 1
< error
| 1193 HY000 Unknown system variable 'a.description'
> exec ^
: string
| set a.deterministic = 1
< error
| 1193 HY000 Unknown system variable 'a.deterministic'
> exec ^
: string
| set a.diagnostics = 1
< error
| 1193 HY000 Unknown system variable 'a.diagnostics'
> exec ^
: string
| set a.directory = 1
< error
| 1193 HY000 Unknown system variable 'a.directory'
> exec ^
: string
| set a.disable = 1
< error
| 1193 HY000 Unknown system variable 'a.disable'
> exec ^
: string
| set a.discard = 1
< error
| 1193 HY000 Unknown system variable 'a.discard'
> exec ^
: string
| set a.disk = 1
< error
| 1193 HY000 Unknown system variable 'a.disk'
> exec ^
: string
| set a.distinct = 1
< error
| 1193 HY000 Unknown system variable 'a.distinct'
> exec ^
: string
| set a.distinctrow = 1
< error
| 1193 HY000 Unknown system variable 'a.distinctrow'
> exec ^
: string
| set a.div = 1
< error
| 1193 HY000 Unknown system variable 'a.div'
> exec ^
: string
| set a.do = 1
< error
| 1193 HY000 Unknown system variable 'a.do'
> exec ^
: string
| set a.double = 1
< error
| 1193 HY000 Unknown system variable 'a.double'
> exec ^
: string
| set a.drop = 1
< error
| 1193 HY000 Unknown system variable 'a.drop'
> exec ^
: string
| set a.dual = 1
< error
| 1193 HY000 Unknown system variable 'a.dual'
> exec ^
: string
| set a.duality = 1
< error
| 1193 HY000 Unknown system variable 'a.duality'
> exec ^
: string
| set a.dumpfile = 1
< error
| 1193 HY000 Unknown system variable 'a.dumpfile'
> exec ^
: string
| set a.duplicate = 1
< error
| 1193 HY000 Unknown system variable 'a.duplicate'
> exec ^
: string
| set a.dynamic = 1
< error
| 1193 HY000 Unknown system variable 'a.dynamic'
> exec ^
: string
| set a.each = 1
< error
| 1193 HY000 Unknown system variable 'a.each'
> exec ^
: string
| set a.else = 1
< error
| 1193 HY000 Unknown system variable 'a.else'
> exec ^
: string
| set a.elseif = 1
< error
| 1193 HY000 Unknown system variable 'a.elseif'
> exec ^
: string
| set a.empty = 1
< error
| 1193 HY000 Unknown system variable 'a.empty'
> exec ^
: string
| set a.enable = 1
< error
| 1193 HY000 Unknown system variable 'a.enable'
> exec ^
: string
| set a.enclosed = 1
< error
| 1193 HY000 Unknown system variable 'a.enclosed'
> exec ^
: string
| set a.encryption = 1
< error
| 1193 HY000 Unknown system variable 'a.encryption'
> exec ^
: string
| set a.end = 1
< error
| 1193 HY000 Unknown system variable 'a.end'
> exec ^
: string
| set a.ends = 1
< error
| 1193 HY000 Unknown system variable 'a.ends'
> exec ^
: string
| set a.enforced = 1
< error
| 1193 HY000 Unknown system variable 'a.enforced'
> exec ^
: string
| set a.engine = 1
< error
| 1193 HY000 Unknown system variable 'a.engine'
> exec ^
: string
| set a.engine_attribute = 1
< error
| 1193 HY000 Unknown system variable 'a.engine_attribute'
> exec ^
: string
| set a.engines = 1
< error
| 1193 HY000 Unknown system variable 'a.engines'
> exec ^
: string
| set a.enum = 1
< error
| 1193 HY000 Unknown system variable 'a.enum'
> exec ^
: string
| set a.error = 1
< error
| 1193 HY000 Unknown system variable 'a.error'
> exec ^
: string
| set a.errors = 1
< error
| 1193 HY000 Unknown system variable 'a.errors'
> exec ^
: string
| set a.escape = 1
< error
| 1193 HY000 Unknown system variable 'a.escape'
> exec ^
: string
| set a.escaped = 1
< error
| 1193 HY000 Unknown system variable 'a.escaped'
> exec ^
: string
| set a.event = 1
< error
| 1193 HY000 Unknown system variable 'a.event'
> exec ^
: string
| set a.events = 1
< error
| 1193 HY000 Unknown system variable 'a.events'
> exec ^
: string
| set a.every = 1
< error
| 1193 HY000 Unknown system variable 'a.every'
> exec ^
: string
| set a.except = 1
< error
| 1193 HY000 Unknown system variable 'a.except'
> exec ^
: string
| set a.exchange = 1
< error
| 1193 HY000 Unknown system variable 'a.exchange'
> exec ^
: string
| set a.exclude = 1
< error
| 1193 HY000 Unknown system variable 'a.exclude'
> exec ^
: string
| set a.execute = 1
< error
| 1193 HY000 Unknown system variable 'a.execute'
> exec ^
: string
| set a.exists = 1
< error
| 1193 HY000 Unknown system variable 'a.exists'
> exec ^
: string
| set a.exit = 1
< error
| 1193 HY000 Unknown system variable 'a.exit'
> exec ^
: string
| set a.expansion = 1
< error
| 1193 HY000 Unknown system variable 'a.expansion'
> exec ^
: string
| set a.expire = 1
< error
| 1193 HY000 Unknown system variable 'a.expire'
> exec ^
: string
| set a.explain = 1
< error
| 1193 HY000 Unknown system variable 'a.explain'
> exec ^
: string
| set a.export = 1
< error
| 1193 HY000 Unknown system variable 'a.export'
> exec ^
: string
| set a.extended = 1
< error
| 1193 HY000 Unknown system variable 'a.extended'
> exec ^
: string
| set a.extent_size = 1
< error
| 1193 HY000 Unknown system variable 'a.extent_size'
> exec ^
: string
| set a.external = 1
< error
| 1193 HY000 Unknown system variable 'a.external'
> exec ^
: string
| set a.external_format = 1
< error
| 1193 HY000 Unknown system variable 'a.external_format'
> exec ^
: string
| set a.extract = 1
< error
| 1193 HY000 Unknown system variable 'a.extract'
> exec ^
: string
| set a.factor = 1
< error
| 1193 HY000 Unknown system variable 'a.factor'
> exec ^
: string
| set a.failed_login_attempts = 1
< error
| 1193 HY000 Unknown system variable 'a.failed_login_attempts'
> exec ^
: string
| set a.false = 1
< error
| 1193 HY000 Unknown system variable 'a.false'
> exec ^
: string
| set a.fast = 1
< error
| 1193 HY000 Unknown system variable 'a.fast'
> exec ^
: string
| set a.faults = 1
< error
| 1193 HY000 Unknown system variable 'a.faults'
> exec ^
: string
| set a.fetch = 1
< error
| 1193 HY000 Unknown system variable 'a.fetch'
> exec ^
: string
| set a.fields = 1
< error
| 1193 HY000 Unknown system variable 'a.fields'
> exec ^
: string
| set a.file = 1
< error
| 1193 HY000 Unknown system variable 'a.file'
> exec ^
: string
| set a.file_block_size = 1
< error
| 1193 HY000 Unknown system variable 'a.file_block_size'
> exec ^
: string
| set a.file_format = 1
< error
| 1193 HY000 Unknown system variable 'a.file_format'
> exec ^
: string
| set a.file_name = 1
< error
| 1193 HY000 Unknown system variable 'a.file_name'
> exec ^
: string
| set a.file_pattern = 1
< error
| 1193 HY000 Unknown system variable 'a.file_pattern'
> exec ^
: string
| set a.file_prefix = 1
< error
| 1193 HY000 Unknown system variable 'a.file_prefix'
> exec ^
: string
| set a.files = 1
< error
| 1193 HY000 Unknown system variable 'a.files'
> exec ^
: string
| set a.filter = 1
< error
| 1193 HY000 Unknown system variable 'a.filter'
> exec ^
: string
| set a.finish = 1
< error
| 1193 HY000 Unknown system variable 'a.finish'
> exec ^
: string
| set a.first = 1
< error
| 1193 HY000 Unknown system variable 'a.first'
> exec ^
: string
| set a.first_value = 1
< error
| 1193 HY000 Unknown system variable 'a.first_value'
> exec ^
: string
| set a.fixed = 1
< error
| 1193 HY000 Unknown system variable 'a.fixed'
> exec ^
: string
| set a.float = 1
< error
| 1193 HY000 Unknown system variable 'a.float'
> exec ^
: string
| set a.float4 = 1
< error
| 1193 HY000 Unknown system variable 'a.float4'
> exec ^
: string
| set a.float8 = 1
< error
| 1193 HY000 Unknown system variable 'a.float8'
> exec ^
: string
| set a.flush = 1
< error
| 1193 HY000 Unknown system variable 'a.flush'
> exec ^
: string
| set a.following = 1
< error
| 1193 HY000 Unknown system variable 'a.following'
> exec ^
: string
| set a.follows = 1
< error
| 1193 HY000 Unknown system variable 'a.follows'
> exec ^
: string
| set a.for = 1
< error
| 1193 HY000 Unknown system variable 'a.for'
> exec ^
: string
| set a.force = 1
< error
| 1193 HY000 Unknown system variable 'a.force'
> exec ^
: string
| set a.foreign = 1
< error
| 1193 HY000 Unknown system variable 'a.foreign'
> exec ^
: string
| set a.format = 1
< error
| 1193 HY000 Unknown system variable 'a.format'
> exec ^
: string
| set a.found = 1
< error
| 1193 HY000 Unknown system variable 'a.found'
> exec ^
: string
| set a.from = 1
< error
| 1193 HY000 Unknown system variable 'a.from'
> exec ^
: string
| set a.full = 1
< error
| 1193 HY000 Unknown system variable 'a.full'
> exec ^
: string
| set a.fulltext = 1
< error
| 1193 HY000 Unknown system variable 'a.fulltext'
> exec ^
: string
| set a.function = 1
< error
| 1193 HY000 Unknown system variable 'a.function'
> exec ^
: string
| set a.general = 1
< error
| 1193 HY000 Unknown system variable 'a.general'
> exec ^
: string
| set a.generate = 1
< error
| 1193 HY000 Unknown system variable 'a.generate'
> exec ^
: string
| set a.generated = 1
< error
| 1193 HY000 Unknown system variable 'a.generated'
> exec ^
: string
| set a.geomcollection = 1
< error
| 1193 HY000 Unknown system variable 'a.geomcollection'
> exec ^
: string
| set a.geometry = 1
< error
| 1193 HY000 Unknown system variable 'a.geometry'
> exec ^
: string
| set a.geometrycollection = 1
< error
| 1193 HY000 Unknown system variable 'a.geometrycollection'
> exec ^
: string
| set a.get = 1
< error
| 1193 HY000 Unknown system variable 'a.get'
> exec ^
: string
| set a.get_format = 1
< error
| 1193 HY000 Unknown system variable 'a.get_format'
> exec ^
: string
| set a.get_master_public_key = 1
< error
| 1193 HY000 Unknown system variable 'a.get_master_public_key'
> exec ^
: string
| set a.get_source_public_key = 1
< error
| 1193 HY000 Unknown system variable 'a.get_source_public_key'
> exec ^
: string
| set a.global = 1
< error
| 1193 HY000 Unknown system variable 'a.global'
> exec ^
: string
| set a.grant = 1
< error
| 1193 HY000 Unknown system variable 'a.grant'
> exec ^
: string
| set a.grants = 1
< error
| 1193 HY000 Unknown system variable 'a.grants'
> exec ^
: string
| set a.greatest = 1
< error
| 1193 HY000 Unknown system variable 'a.greatest'
> exec ^
: string
| set a.group = 1
< error
| 1193 HY000 Unknown system variable 'a.group'
> exec ^
: string
| set a.group_concat = 1
< error
| 1193 HY000 Unknown system variable 'a.group_concat'
> exec ^
: string
| set a.group_replication = 1
< error
| 1193 HY000 Unknown system variable 'a.group_replication'
> exec ^
: string
| set a.grouping = 1
< error
| 1193 HY000 Unknown system variable 'a.grouping'
> exec ^
: string
| set a.groups = 1
< error
| 1193 HY000 Unknown system variable 'a.groups'
> exec ^
: string
| set a.gtid_only = 1
< error
| 1193 HY000 Unknown system variable 'a.gtid_only'
> exec ^
: string
| set a.gtids = 1
< error
| 1193 HY000 Unknown system variable 'a.gtids'
> exec ^
: string
| set a.guided = 1
< error
| 1193 HY000 Unknown system variable 'a.guided'
> exec ^
: string
| set a.handler = 1
< error
| 1193 HY000 Unknown system variable 'a.handler'
> exec ^
: string
| set a.hash = 1
< error
| 1193 HY000 Unknown system variable 'a.hash'
> exec ^
: string
| set a.having = 1
< error
| 1193 HY000 Unknown system variable 'a.having'
> exec ^
: string
| set a.header = 1
< error
| 1193 HY000 Unknown system variable 'a.header'
> exec ^
: string
| set a.help = 1
< error
| 1193 HY000 Unknown system variable 'a.help'
> exec ^
: string
| set a.high_priority = 1
< error
| 1193 HY000 Unknown system variable 'a.high_priority'
> exec ^
: string
| set a.histogram = 1
< error
| 1193 HY000 Unknown system variable 'a.histogram'
> exec ^
: string
| set a.history = 1
< error
| 1193 HY000 Unknown system variable 'a.history'
> exec ^
: string
| set a.host = 1
< error
| 1193 HY000 Unknown system variable 'a.host'
> exec ^
: string
| set a.hosts = 1
< error
| 1193 HY000 Unknown system variable 'a.hosts'
> exec ^
: string
| set a.hour = 1
< error
| 1193 HY000 Unknown system variable 'a.hour'
> exec ^
: string
| set a.hour_microsecond = 1
< error
| 1193 HY000 Unknown system variable 'a.hour_microsecond'
> exec ^
: string
| set a.hour_minute = 1
< error
| 1193 HY000 Unknown system variable 'a.hour_minute'
> exec ^
: string
| set a.hour_second = 1
< error
| 1193 HY000 Unknown system variable 'a.hour_second'
> exec ^
: string
| set a.identified = 1
< error
| 1193 HY000 Unknown system variable 'a.identified'
> exec ^
: string
| set a.if = 1
< error
| 1193 HY000 Unknown system variable 'a.if'
> exec ^
: string
| set a.ignore = 1
< error
| 1193 HY000 Unknown system variable 'a.ignore'
> exec ^
: string
| set a.ignore_server_ids = 1
< error
| 1193 HY000 Unknown system variable 'a.ignore_server_ids'
> exec ^
: string
| set a.import = 1
< error
| 1193 HY000 Unknown system variable 'a.import'
> exec ^
: string
| set a.in = 1
< error
| 1193 HY000 Unknown system variable 'a.in'
> exec ^
: string
| set a.inactive = 1
< error
| 1193 HY000 Unknown system variable 'a.inactive'
> exec ^
: string
| set a.index = 1
< error
| 1193 HY000 Unknown system variable 'a.index'
> exec ^
: string
| set a.indexes = 1
< error
| 1193 HY000 Unknown system variable 'a.indexes'
> exec ^
: string
| set a.infile = 1
< error
| 1193 HY000 Unknown system variable 'a.infile'
> exec ^
: string
| set a.initial = 1
< error
| 1193 HY000 Unknown system variable 'a.initial'
> exec ^
: string
| set a.initial_size = 1
< error
| 1193 HY000 Unknown system variable 'a.initial_size'
> exec ^
: string
| set a.initiate = 1
< error
| 1193 HY000 Unknown system variable 'a.initiate'
> exec ^
: string
| set a.inner = 1
< error
| 1193 HY000 Unknown system variable 'a.inner'
> exec ^
: string
| set a.inout = 1
< error
| 1193 HY000 Unknown system variable 'a.inout'
> exec ^
: string
| set a.insensitive = 1
< error
| 1193 HY000 Unknown system variable 'a.insensitive'
> exec ^
: string
| set a.insert = 1
< error
| 1193 HY000 Unknown system variable 'a.insert'
> exec ^
: string
| set a.insert_method = 1
< error
| 1193 HY000 Unknown system variable 'a.insert_method'
> exec ^
: string
| set a.install = 1
< error
| 1193 HY000 Unknown system variable 'a.install'
> exec ^
: string
| set a.instance = 1
< error
| 1193 HY000 Unknown system variable 'a.instance'
> exec ^
: string
| set a.int = 1
< error
| 1193 HY000 Unknown system variable 'a.int'
> exec ^
: string
| set a.int1 = 1
< error
| 1193 HY000 Unknown system variable 'a.int1'
> exec ^
: string
| set a.int2 = 1
< error
| 1193 HY000 Unknown system variable 'a.int2'
> exec ^
: string
| set a.int3 = 1
< error
| 1193 HY000 Unknown system variable 'a.int3'
> exec ^
: string
| set a.int4 = 1
< error
| 1193 HY000 Unknown system variable 'a.int4'
> exec ^
: string
| set a.int8 = 1
< error
| 1193 HY000 Unknown system variable 'a.int8'
> exec ^
: string
| set a.integer = 1
< error
| 1193 HY000 Unknown system variable 'a.integer'
> exec ^
: string
| set a.intersect = 1
< error
| 1193 HY000 Unknown system variable 'a.intersect'
> exec ^
: string
| set a.interval = 1
< error
| 1193 HY000 Unknown system variable 'a.interval'
> exec ^
: string
| set a.into = 1
< error
| 1193 HY000 Unknown system variable 'a.into'
> exec ^
: string
| set a.invisible = 1
< error
| 1193 HY000 Unknown system variable 'a.invisible'
> exec ^
: string
| set a.invoker = 1
< error
| 1193 HY000 Unknown system variable 'a.invoker'
> exec ^
: string
| set a.io = 1
< error
| 1193 HY000 Unknown system variable 'a.io'
> exec ^
: string
| set a.io_after_gtids = 1
< error
| 1193 HY000 Unknown system variable 'a.io_after_gtids'
> exec ^
: string
| set a.io_before_gtids = 1
< error
| 1193 HY000 Unknown system variable 'a.io_before_gtids'
> exec ^
: string
| set a.io_thread = 1
< error
| 1193 HY000 Unknown system variable 'a.io_thread'
> exec ^
: string
| set a.ipc = 1
< error
| 1193 HY000 Unknown system variable 'a.ipc'
> exec ^
: string
| set a.is = 1
< error
| 1193 HY000 Unknown system variable 'a.is'
> exec ^
: string
| set a.isolation = 1
< error
| 1193 HY000 Unknown system variable 'a.isolation'
> exec ^
: string
| set a.issuer = 1
< error
| 1193 HY000 Unknown system variable 'a.issuer'
> exec ^
: string
| set a.iterate = 1
< error
| 1193 HY000 Unknown system variable 'a.iterate'
> exec ^
: string
| set a.join = 1
< error
| 1193 HY000 Unknown system variable 'a.join'
> exec ^
: string
| set a.json = 1
< error
| 1193 HY000 Unknown system variable 'a.json'
> exec ^
: string
| set a.json_arrayagg = 1
< error
| 1193 HY000 Unknown system variable 'a.json_arrayagg'
> exec ^
: string
| set a.json_duality_object = 1
< error
| 1193 HY000 Unknown system variable 'a.json_duality_object'
> exec ^
: string
| set a.json_objectagg = 1
< error
| 1193 HY000 Unknown system variable 'a.json_objectagg'
> exec ^
: string
| set a.json_table = 1
< error
| 1193 HY000 Unknown system variable 'a.json_table'
> exec ^
: string
| set a.json_value = 1
< error
| 1193 HY000 Unknown system variable 'a.json_value'
> exec ^
: string
| set a.key = 1
< error
| 1193 HY000 Unknown system variable 'a.key'
> exec ^
: string
| set a.key_block_size = 1
< error
| 1193 HY000 Unknown system variable 'a.key_block_size'
> exec ^
: string
| set a.keyring = 1
< error
| 1193 HY000 Unknown system variable 'a.keyring'
> exec ^
: string
| set a.keys = 1
< error
| 1193 HY000 Unknown system variable 'a.keys'
> exec ^
: string
| set a.kill = 1
< error
| 1193 HY000 Unknown system variable 'a.kill'
> exec ^
: string
| set a.lag = 1
< error
| 1193 HY000 Unknown system variable 'a.lag'
> exec ^
: string
| set a.language = 1
< error
| 1193 HY000 Unknown system variable 'a.language'
> exec ^
: string
| set a.last = 1
< error
| 1193 HY000 Unknown system variable 'a.last'
> exec ^
: string
| set a.last_value = 1
< error
| 1193 HY000 Unknown system variable 'a.last_value'
> exec ^
: string
| set a.lateral = 1
< error
| 1193 HY000 Unknown system variable 'a.lateral'
> exec ^
: string
| set a.lead = 1
< error
| 1193 HY000 Unknown system variable 'a.lead'
> exec ^
: string
| set a.leading = 1
< error
| 1193 HY000 Unknown system variable 'a.leading'
> exec ^
: string
| set a.least = 1
< error
| 1193 HY000 Unknown system variable 'a.least'
> exec ^
: string
| set a.leave = 1
< error
| 1193 HY000 Unknown system variable 'a.leave'
> exec ^
: string
| set a.leaves = 1
< error
| 1193 HY000 Unknown system variable 'a.leaves'
> exec ^
: string
| set a.left = 1
< error
| 1193 HY000 Unknown system variable 'a.left'
> exec ^
: string
| set a.less = 1
< error
| 1193 HY000 Unknown system variable 'a.less'
> exec ^
: string
| set a.level = 1
< error
| 1193 HY000 Unknown system variable 'a.level'
> exec ^
: string
| set a.library = 1
< error
| 1193 HY000 Unknown system variable 'a.library'
> exec ^
: string
| set a.like = 1
< error
| 1193 HY000 Unknown system variable 'a.like'
> exec ^
: string
| set a.limit = 1
< error
| 1193 HY000 Unknown system variable 'a.limit'
> exec ^
: string
| set a.linear = 1
< error
| 1193 HY000 Unknown system variable 'a.linear'
> exec ^
: string
| set a.lines = 1
< error
| 1193 HY000 Unknown system variable 'a.lines'
> exec ^
: string
| set a.linestring = 1
< error
| 1193 HY000 Unknown system variable 'a.linestring'
> exec ^
: string
| set a.list = 1
< error
| 1193 HY000 Unknown system variable 'a.list'
> exec ^
: string
| set a.load = 1
< error
| 1193 HY000 Unknown system variable 'a.load'
> exec ^
: string
| set a.local = 1
< error
| 1193 HY000 Unknown system variable 'a.local'
> exec ^
: string
| set a.localtime = 1
< error
| 1193 HY000 Unknown system variable 'a.localtime'
> exec ^
: string
| set a.localtimestamp = 1
< error
| 1193 HY000 Unknown system variable 'a.localtimestamp'
> exec ^
: string
| set a.lock = 1
< error
| 1193 HY000 Unknown system variable 'a.lock'
> exec ^
: string
| set a.locked = 1
< error
| 1193 HY000 Unknown system variable 'a.locked'
> exec ^
: string
| set a.locks = 1
< error
| 1193 HY000 Unknown system variable 'a.locks'
> exec ^
: string
| set a.log = 1
< error
| 1193 HY000 Unknown system variable 'a.log'
> exec ^
: string
| set a.logfile = 1
< error
| 1193 HY000 Unknown system variable 'a.logfile'
> exec ^
: string
| set a.logs = 1
< error
| 1193 HY000 Unknown system variable 'a.logs'
> exec ^
: string
| set a.long = 1
< error
| 1193 HY000 Unknown system variable 'a.long'
> exec ^
: string
| set a.longblob = 1
< error
| 1193 HY000 Unknown system variable 'a.longblob'
> exec ^
: string
| set a.longtext = 1
< error
| 1193 HY000 Unknown system variable 'a.longtext'
> exec ^
: string
| set a.loop = 1
< error
| 1193 HY000 Unknown system variable 'a.loop'
> exec ^
: string
| set a.low_priority = 1
< error
| 1193 HY000 Unknown system variable 'a.low_priority'
> exec ^
: string
| set a.manifest = 1
< error
| 1193 HY000 Unknown system variable 'a.manifest'
> exec ^
: string
| set a.manual = 1
< error
| 1193 HY000 Unknown system variable 'a.manual'
> exec ^
: string
| set a.master = 1
< error
| 1193 HY000 Unknown system variable 'a.master'
> exec ^
: string
| set a.master_auto_position = 1
< error
| 1193 HY000 Unknown system variable 'a.master_auto_position'
> exec ^
: string
| set a.master_bind = 1
< error
| 1193 HY000 Unknown system variable 'a.master_bind'
> exec ^
: string
| set a.master_compression_algorithms = 1
< error
| 1193 HY000 Unknown system variable 'a.master_compression_algorithms'
> exec ^
: string
| set a.master_connect_retry = 1
< error
| 1193 HY000 Unknown system variable 'a.master_connect_retry'
> exec ^
: string
| set a.master_delay = 1
< error
| 1193 HY000 Unknown system variable 'a.master_delay'
> exec ^
: string
| set a.master_heartbeat_period = 1
< error
| 1193 HY000 Unknown system variable 'a.master_heartbeat_period'
> exec ^
: string
| set a.master_host = 1
< error
| 1193 HY000 Unknown system variable 'a.master_host'
> exec ^
: string
| set a.master_log_file = 1
< error
| 1193 HY000 Unknown system variable 'a.master_log_file'
> exec ^
: string
| set a.master_log_pos = 1
< error
| 1193 HY000 Unknown system variable 'a.master_log_pos'
> exec ^
: string
| set a.master_password = 1
< error
| 1193 HY000 Unknown system variable 'a.master_password'
> exec ^
: string
| set a.master_port = 1
< error
| 1193 HY000 Unknown system variable 'a.master_port'
> exec ^
: string
| set a.master_public_key_path = 1
< error
| 1193 HY000 Unknown system variable 'a.master_public_key_path'
> exec ^
: string
| set a.master_retry_count = 1
< error
| 1193 HY000 Unknown system variable 'a.master_retry_count'
> exec ^
: string
| set a.master_ssl = 1
< error
| 1193 HY000 Unknown system variable 'a.master_ssl'
> exec ^
: string
| set a.master_ssl_ca = 1
< error
| 1193 HY000 Unknown system variable 'a.master_ssl_ca'
> exec ^
: string
| set a.master_ssl_capath = 1
< error
| 1193 HY000 Unknown system variable 'a.master_ssl_capath'
> exec ^
: string
| set a.master_ssl_cert = 1
< error
| 1193 HY000 Unknown system variable 'a.master_ssl_cert'
> exec ^
: string
| set a.master_ssl_cipher = 1
< error
| 1193 HY000 Unknown system variable 'a.master_ssl_cipher'
> exec ^
: string
| set a.master_ssl_crl = 1
< error
| 1193 HY000 Unknown system variable 'a.master_ssl_crl'
> exec ^
: string
| set a.master_ssl_crlpath = 1
< error
| 1193 HY000 Unknown system variable 'a.master_ssl_crlpath'
> exec ^
: string
| set a.master_ssl_key = 1
< error
| 1193 HY000 Unknown system variable 'a.master_ssl_key'
> exec ^
: string
| set a.master_ssl_verify_server_cert = 1
< error
| 1193 HY000 Unknown system variable 'a.master_ssl_verify_server_cert'
> exec ^
: string
| set a.master_tls_ciphersuites = 1
< error
| 1193 HY000 Unknown system variable 'a.master_tls_ciphersuites'
> exec ^
: string
| set a.master_tls_version = 1
< error
| 1193 HY000 Unknown system variable 'a.master_tls_version'
> exec ^
: string
| set a.master_user = 1
< error
| 1193 HY000 Unknown system variable 'a.master_user'
> exec ^
: string
| set a.master_zstd_compression_level = 1
< error
| 1193 HY000 Unknown system variable 'a.master_zstd_compression_level'
> exec ^
: string
| set a.match = 1
< error
| 1193 HY000 Unknown system variable 'a.match'
> exec ^
: string
| set a.materialized = 1
< error
| 1193 HY000 Unknown system variable 'a.materialized'
> exec ^
: string
| set a.max = 1
< error
| 1193 HY000 Unknown system variable 'a.max'
> exec ^
: string
| set a.max_connections_per_hour = 1
< error
| 1193 HY000 Unknown system variable 'a.max_connections_per_hour'
> exec ^
: string
| set a.max_queries_per_hour = 1
< error
| 1193 HY000 Unknown system variable 'a.max_queries_per_hour'
> exec ^
: string
| set a.max_rows = 1
< error
| 1193 HY000 Unknown system variable 'a.max_rows'
> exec ^
: string
| set a.max_size = 1
< error
| 1193 HY000 Unknown system variable 'a.max_size'
> exec ^
: string
| set a.max_updates_per_hour = 1
< error
| 1193 HY000 Unknown system variable 'a.max_updates_per_hour'
> exec ^
: string
| set a.max_user_connections = 1
< error
| 1193 HY000 Unknown system variable 'a.max_user_connections'
> exec ^
: string
| set a.maxvalue = 1
< error
| 1193 HY000 Unknown system variable 'a.maxvalue'
> exec ^
: string
| set a.medium = 1
< error
| 1193 HY000 Unknown system variable 'a.medium'
> exec ^
: string
| set a.mediumblob = 1
< error
| 1193 HY000 Unknown system variable 'a.mediumblob'
> exec ^
: string
| set a.mediumint = 1
< error
| 1193 HY000 Unknown system variable 'a.mediumint'
> exec ^
: string
| set a.mediumtext = 1
< error
| 1193 HY000 Unknown system variable 'a.mediumtext'
> exec ^
: string
| set a.member = 1
< error
| 1193 HY000 Unknown system variable 'a.member'
> exec ^
: string
| set a.memory = 1
< error
| 1193 HY000 Unknown system variable 'a.memory'
> exec ^
: string
| set a.merge = 1
< error
| 1193 HY000 Unknown system variable 'a.merge'
> exec ^
: string
| set a.message_text = 1
< error
| 1193 HY000 Unknown system variable 'a.message_text'
> exec ^
: string
| set a.microsecond = 1
< error
| 1193 HY000 Unknown system variable 'a.microsecond'
> exec ^
: string
| set a.mid = 1
< error
| 1193 HY000 Unknown system variable 'a.mid'
> exec ^
: string
| set a.middleint = 1
< error
| 1193 HY000 Unknown system variable 'a.middleint'
> exec ^
: string
| set a.migrate = 1
< error
| 1193 HY000 Unknown system variable 'a.migrate'
> exec ^
: string
| set a.min = 1
< error
| 1193 HY000 Unknown system variable 'a.min'
> exec ^
: string
| set a.min_rows = 1
< error
| 1193 HY000 Unknown system variable 'a.min_rows'
> exec ^
: string
| set a.minute = 1
< error
| 1193 HY000 Unknown system variable 'a.minute'
> exec ^
: string
| set a.minute_microsecond = 1
< error
| 1193 HY000 Unknown system variable 'a.minute_microsecond'
> exec ^
: string
| set a.minute_second = 1
< error
| 1193 HY000 Unknown system variable 'a.minute_second'
> exec ^
: string
| set a.mod = 1
< error
| 1193 HY000 Unknown system variable 'a.mod'
> exec ^
: string
| set a.mode = 1
< error
| 1193 HY000 Unknown system variable 'a.mode'
> exec ^
: string
| set a.modifies = 1
< error
| 1193 HY000 Unknown system variable 'a.modifies'
> exec ^
: string
| set a.modify = 1
< error
| 1193 HY000 Unknown system variable 'a.modify'
> exec ^
: string
| set a.month = 1
< error
| 1193 HY000 Unknown system variable 'a.month'
> exec ^
: string
| set a.multilinestring = 1
< error
| 1193 HY000 Unknown system variable 'a.multilinestring'
> exec ^
: string
| set a.multipoint = 1
< error
| 1193 HY000 Unknown system variable 'a.multipoint'
> exec ^
: string
| set a.multipolygon = 1
< error
| 1193 HY000 Unknown system variable 'a.multipolygon'
> exec ^
: string
| set a.mutex = 1
< error
| 1193 HY000 Unknown system variable 'a.mutex'
> exec ^
: string
| set a.mysql_errno = 1
< error
| 1193 HY000 Unknown system variable 'a.mysql_errno'
> exec ^
: string
| set a.name = 1
< error
| 1193 HY000 Unknown system variable 'a.name'
> exec ^
: string
| set a.names = 1
< error
| 1193 HY000 Unknown system variable 'a.names'
> exec ^
: string
| set a.national = 1
< error
| 1193 HY000 Unknown system variable 'a.national'
> exec ^
: string
| set a.natural = 1
< error
| 1193 HY000 Unknown system variable 'a.natural'
> exec ^
: string
| set a.nchar = 1
< error
| 1193 HY000 Unknown system variable 'a.nchar'
> exec ^
: string
| set a.ndb = 1
< error
| 1193 HY000 Unknown system variable 'a.ndb'
> exec ^
: string
| set a.ndbcluster = 1
< error
| 1193 HY000 Unknown system variable 'a.ndbcluster'
> exec ^
: string
| set a.nested = 1
< error
| 1193 HY000 Unknown system variable 'a.nested'
> exec ^
: string
| set a.network_namespace = 1
< error
| 1193 HY000 Unknown system variable 'a.network_namespace'
> exec ^
: string
| set a.never = 1
< error
| 1193 HY000 Unknown system variable 'a.never'
> exec ^
: string
| set a.new = 1
< error
| 1193 HY000 Unknown system variable 'a.new'
> exec ^
: string
| set a.next = 1
< error
| 1193 HY000 Unknown system variable 'a.next'
> exec ^
: string
| set a.no = 1
< error
| 1193 HY000 Unknown system variable 'a.no'
> exec ^
: string
| set a.no_wait = 1
< error
| 1193 HY000 Unknown system variable 'a.no_wait'
> exec ^
: string
| set a.no_write_to_binlog = 1
< error
| 1193 HY000 Unknown system variable 'a.no_write_to_binlog'
> exec ^
: string
| set a.nodegroup = 1
< error
| 1193 HY000 Unknown system variable 'a.nodegroup'
> exec ^
: string
| set a.none = 1
< error
| 1193 HY000 Unknown system variable 'a.none'
> exec ^
: string
| set a.not = 1
< error
| 1193 HY000 Unknown system variable 'a.not'
> exec ^
: string
| set a.now = 1
< error
| 1193 HY000 Unknown system variable 'a.now'
> exec ^
: string
| set a.nowait = 1
< error
| 1193 HY000 Unknown system variable 'a.nowait'
> exec ^
: string
| set a.nth_value = 1
< error
| 1193 HY000 Unknown system variable 'a.nth_value'
> exec ^
: string
| set a.ntile = 1
< error
| 1193 HY000 Unknown system variable 'a.ntile'
> exec ^
: string
| set a.null = 1
< error
| 1193 HY000 Unknown system variable 'a.null'
> exec ^
: string
| set a.nullif = 1
< error
| 1193 HY000 Unknown system variable 'a.nullif'
> exec ^
: string
| set a.nulls = 1
< error
| 1193 HY000 Unknown system variable 'a.nulls'
> exec ^
: string
| set a.number = 1
< error
| 1193 HY000 Unknown system variable 'a.number'
> exec ^
: string
| set a.numeric = 1
< error
| 1193 HY000 Unknown system variable 'a.numeric'
> exec ^
: string
| set a.nvarchar = 1
< error
| 1193 HY000 Unknown system variable 'a.nvarchar'
> exec ^
: string
| set a.of = 1
< error
| 1193 HY000 Unknown system variable 'a.of'
> exec ^
: string
| set a.off = 1
< error
| 1193 HY000 Unknown system variable 'a.off'
> exec ^
: string
| set a.offset = 1
< error
| 1193 HY000 Unknown system variable 'a.offset'
> exec ^
: string
| set a.oj = 1
< error
| 1193 HY000 Unknown system variable 'a.oj'
> exec ^
: string
| set a.old = 1
< error
| 1193 HY000 Unknown system variable 'a.old'
> exec ^
: string
| set a.on = 1
< error
| 1193 HY000 Unknown system variable 'a.on'
> exec ^
: string
| set a.one = 1
< error
| 1193 HY000 Unknown system variable 'a.one'
> exec ^
: string
| set a.only = 1
< error
| 1193 HY000 Unknown system variable 'a.only'
> exec ^
: string
| set a.open = 1
< error
| 1193 HY000 Unknown system variable 'a.open'
> exec ^
: string
| set a.optimize = 1
< error
| 1193 HY000 Unknown system variable 'a.optimize'
> exec ^
: string
| set a.optimizer_costs = 1
< error
| 1193 HY000 Unknown system variable 'a.optimizer_costs'
> exec ^
: string
| set a.option = 1
< error
| 1193 HY000 Unknown system variable 'a.option'
> exec ^
: string
| set a.optional = 1
< error
| 1193 HY000 Unknown system variable 'a.optional'
> exec ^
: string
| set a.optionally = 1
< error
| 1193 HY000 Unknown system variable 'a.optionally'
> exec ^
: string
| set a.options = 1
< error
| 1193 HY000 Unknown system variable 'a.options'
> exec ^
: string
| set a.or = 1
< error
| 1193 HY000 Unknown system variable 'a.or'
> exec ^
: string
| set a.order = 1
< error
| 1193 HY000 Unknown system variable 'a.order'
> exec ^
: string
| set a.ordinality = 1
< error
| 1193 HY000 Unknown system variable 'a.ordinality'
> exec ^
: string
| set a.organization = 1
< error
| 1193 HY000 Unknown system variable 'a.organization'
> exec ^
: string
| set a.others = 1
< error
| 1193 HY000 Unknown system variable 'a.others'
> exec ^
: string
| set a.out = 1
< error
| 1193 HY000 Unknown system variable 'a.out'
> exec ^
: string
| set a.outer = 1
< error
| 1193 HY000 Unknown system variable 'a.outer'
> exec ^
: string
| set a.outfile = 1
< error
| 1193 HY000 Unknown system variable 'a.outfile'
> exec ^
: string
| set a.over = 1
< error
| 1193 HY000 Unknown system variable 'a.over'
> exec ^
: string
| set a.overlay = 1
< error
| 1193 HY000 Unknown system variable 'a.overlay'
> exec ^
: string
| set a.overwrite = 1
< error
| 1193 HY000 Unknown system variable 'a.overwrite'
> exec ^
: string
| set a.owner = 1
< error
| 1193 HY000 Unknown system variable 'a.owner'
> exec ^
: string
| set a.pack_keys = 1
< error
| 1193 HY000 Unknown system variable 'a.pack_keys'
> exec ^
: string
| set a.page = 1
< error
| 1193 HY000 Unknown system variable 'a.page'
> exec ^
: string
| set a.parallel = 1
< error
| 1193 HY000 Unknown system variable 'a.parallel'
> exec ^
: string
| set a.parameters = 1
< error
| 1193 HY000 Unknown system variable 'a.parameters'
> exec ^
: string
| set a.parse_tree = 1
< error
| 1193 HY000 Unknown system variable 'a.parse_tree'
> exec ^
: string
| set a.parser = 1
< error
| 1193 HY000 Unknown system variable 'a.parser'
> exec ^
: string
| set a.partial = 1
< error
| 1193 HY000 Unknown system variable 'a.partial'
> exec ^
: string
| set a.partition = 1
< error
| 1193 HY000 Unknown system variable 'a.partition'
> exec ^
: string
| set a.partitioning = 1
< error
| 1193 HY000 Unknown system variable 'a.partitioning'
> exec ^
: string
| set a.partitions = 1
< error
| 1193 HY000 Unknown system variable 'a.partitions'
> exec ^
: string
| set a.password = 1
< error
| 1193 HY000 Unknown system variable 'a.password'
> exec ^
: string
| set a.password_lock_time = 1
< error
| 1193 HY000 Unknown system variable 'a.password_lock_time'
> exec ^
: string
| set a.path = 1
< error
| 1193 HY000 Unknown system variable 'a.path'
> exec ^
: string
| set a.percent_rank = 1
< error
| 1193 HY000 Unknown system variable 'a.percent_rank'
> exec ^
: string
| set a.persist = 1
< error
| 1193 HY000 Unknown system variable 'a.persist'
> exec ^
: string
| set a.persist_only = 1
< error
| 1193 HY000 Unknown system variable 'a.persist_only'
> exec ^
: string
| set a.phase = 1
< error
| 1193 HY000 Unknown system variable 'a.phase'
> exec ^
: string
| set a.placing = 1
< error
| 1193 HY000 Unknown system variable 'a.placing'
> exec ^
: string
| set a.plugin = 1
< error
| 1193 HY000 Unknown system variable 'a.plugin'
> exec ^
: string
| set a.plugin_dir = 1
< error
| 1193 HY000 Unknown system variable 'a.plugin_dir'
> exec ^
: string
| set a.plugins = 1
< error
| 1193 HY000 Unknown system variable 'a.plugins'
> exec ^
: string
| set a.point = 1
< error
| 1193 HY000 Unknown system variable 'a.point'
> exec ^
: string
| set a.polygon = 1
< error
| 1193 HY000 Unknown system variable 'a.polygon'
> exec ^
: string
| set a.port = 1
< error
| 1193 HY000 Unknown system variable 'a.port'
> exec ^
: string
| set a.position = 1
< error
| 1193 HY000 Unknown system variable 'a.position'
> exec ^
: string
| set a.precedes = 1
< error
| 1193 HY000 Unknown system variable 'a.precedes'
> exec ^
: string
| set a.preceding = 1
< error
| 1193 HY000 Unknown system variable 'a.preceding'
> exec ^
: string
| set a.precision = 1
< error
| 1193 HY000 Unknown system variable 'a.precision'
> exec ^
: string
| set a.prefix = 1
< error
| 1193 HY000 Unknown system variable 'a.prefix'
> exec ^
: string
| set a.prepare = 1
< error
| 1193 HY000 Unknown system variable 'a.prepare'
> exec ^
: string
| set a.preserve = 1
< error
| 1193 HY000 Unknown system variable 'a.preserve'
> exec ^
: string
| set a.prev = 1
< error
| 1193 HY000 Unknown system variable 'a.prev'
> exec ^
: string
| set a.primary = 1
< error
| 1193 HY000 Unknown system variable 'a.primary'
> exec ^
: string
| set a.privilege_checks_user = 1
< error
| 1193 HY000 Unknown system variable 'a.privilege_checks_user'
> exec ^
: string
| set a.privileges = 1
< error
| 1193 HY000 Unknown system variable 'a.privileges'
> exec ^
: string
| set a.procedure = 1
< error
| 1193 HY000 Unknown system variable 'a.procedure'
> exec ^
: string
| set a.process = 1
< error
| 1193 HY000 Unknown system variable 'a.process'
> exec ^
: string
| set a.processlist = 1
< error
| 1193 HY000 Unknown system variable 'a.processlist'
> exec ^
: string
| set a.profile = 1
< error
| 1193 HY000 Unknown system variable 'a.profile'
> exec ^
: string
| set a.profiles = 1
< error
| 1193 HY000 Unknown system variable 'a.profiles'
> exec ^
: string
| set a.proxy = 1
< error
| 1193 HY000 Unknown system variable 'a.proxy'
> exec ^
: string
| set a.purge = 1
< error
| 1193 HY000 Unknown system variable 'a.purge'
> exec ^
: string
| set a.qualify = 1
< error
| 1193 HY000 Unknown system variable 'a.qualify'
> exec ^
: string
| set a.quarter = 1
< error
| 1193 HY000 Unknown system variable 'a.quarter'
> exec ^
: string
| set a.query = 1
< error
| 1193 HY000 Unknown system variable 'a.query'
> exec ^
: string
| set a.quick = 1
< error
| 1193 HY000 Unknown system variable 'a.quick'
> exec ^
: string
| set a.random = 1
< error
| 1193 HY000 Unknown system variable 'a.random'
> exec ^
: string
| set a.range = 1
< error
| 1193 HY000 Unknown system variable 'a.range'
> exec ^
: string
| set a.rank = 1
< error
| 1193 HY000 Unknown system variable 'a.rank'
> exec ^
: string
| set a.read = 1
< error
| 1193 HY000 Unknown system variable 'a.read'
> exec ^
: string
| set a.read_only = 1
< error
| 1193 HY000 Unknown system variable 'a.read_only'
> exec ^
: string
| set a.read_write = 1
< error
| 1193 HY000 Unknown system variable 'a.read_write'
> exec ^
: string
| set a.reads = 1
< error
| 1193 HY000 Unknown system variable 'a.reads'
> exec ^
: string
| set a.real = 1
< error
| 1193 HY000 Unknown system variable 'a.real'
> exec ^
: string
| set a.rebuild = 1
< error
| 1193 HY000 Unknown system variable 'a.rebuild'
> exec ^
: string
| set a.recover = 1
< error
| 1193 HY000 Unknown system variable 'a.recover'
> exec ^
: string
| set a.recursive = 1
< error
| 1193 HY000 Unknown system variable 'a.recursive'
> exec ^
: string
| set a.redo_buffer_size = 1
< error
| 1193 HY000 Unknown system variable 'a.redo_buffer_size'
> exec ^
: string
| set a.redundant = 1
< error
| 1193 HY000 Unknown system variable 'a.redundant'
> exec ^
: string
| set a.reference = 1
< error
| 1193 HY000 Unknown system variable 'a.reference'
> exec ^
: string
| set a.references = 1
< error
| 1193 HY000 Unknown system variable 'a.references'
> exec ^
: string
| set a.regexp = 1
< error
| 1193 HY000 Unknown system variable 'a.regexp'
> exec ^
: string
| set a.registration = 1
< error
| 1193 HY000 Unknown system variable 'a.registration'
> exec ^
: string
| set a.relational = 1
< error
| 1193 HY000 Unknown system variable 'a.relational'
> exec ^
: string
| set a.relay = 1
< error
| 1193 HY000 Unknown system variable 'a.relay'
> exec ^
: string
| set a.relay_log_file = 1
< error
| 1193 HY000 Unknown system variable 'a.relay_log_file'
> exec ^
: string
| set a.relay_log_pos = 1
< error
| 1193 HY000 Unknown system variable 'a.relay_log_pos'
> exec ^
: string
| set a.relay_thread = 1
< error
| 1193 HY000 Unknown system variable 'a.relay_thread'
> exec ^
: string
| set a.relaylog = 1
< error
| 1193 HY000 Unknown system variable 'a.relaylog'
> exec ^
: string
| set a.release = 1
< error
| 1193 HY000 Unknown system variable 'a.release'
> exec ^
: string
| set a.reload = 1
< error
| 1193 HY000 Unknown system variable 'a.reload'
> exec ^
: string
| set a.remove = 1
< error
| 1193 HY000 Unknown system variable 'a.remove'
> exec ^
: string
| set a.rename = 1
< error
| 1193 HY000 Unknown system variable 'a.rename'
> exec ^
: string
| set a.reorganize = 1
< error
| 1193 HY000 Unknown system variable 'a.reorganize'
> exec ^
: string
| set a.repair = 1
< error
| 1193 HY000 Unknown system variable 'a.repair'
> exec ^
: string
| set a.repeat = 1
< error
| 1193 HY000 Unknown system variable 'a.repeat'
> exec ^
: string
| set a.repeatable = 1
< error
| 1193 HY000 Unknown system variable 'a.repeatable'
> exec ^
: string
| set a.replace = 1
< error
| 1193 HY000 Unknown system variable 'a.replace'
> exec ^
: string
| set a.replica = 1
< error
| 1193 HY000 Unknown system variable 'a.replica'
> exec ^
: string
| set a.replicas = 1
< error
| 1193 HY000 Unknown system variable 'a.replicas'
> exec ^
: string
| set a.replicate_do_db = 1
< error
| 1193 HY000 Unknown system variable 'a.replicate_do_db'
> exec ^
: string
| set a.replicate_do_table = 1
< error
| 1193 HY000 Unknown system variable 'a.replicate_do_table'
> exec ^
: string
| set a.replicate_ignore_db = 1
< error
| 1193 HY000 Unknown system variable 'a.replicate_ignore_db'
> exec ^
: string
| set a.replicate_ignore_table = 1
< error
| 1193 HY000 Unknown system variable 'a.replicate_ignore_table'
> exec ^
: string
| set a.replicate_rewrite_db = 1
< error
| 1193 HY000 Unknown system variable 'a.replicate_rewrite_db'
> exec ^
: string
| set a.replicate_wild_do_table = 1
< error
| 1193 HY000 Unknown system variable 'a.replicate_wild_do_table'
> exec ^
: string
| set a.replicate_wild_ignore_table = 1
< error
| 1193 HY000 Unknown system variable 'a.replicate_wild_ignore_table'
> exec ^
: string
| set a.replication = 1
< error
| 1193 HY000 Unknown system variable 'a.replication'
> exec ^
: string
| set a.require = 1
< error
| 1193 HY000 Unknown system variable 'a.require'
> exec ^
: string
| set a.require_row_format = 1
< error
| 1193 HY000 Unknown system variable 'a.require_row_format'
> exec ^
: string
| set a.require_table_primary_key_check = 1
< error
| 1193 HY000 Unknown system variable 'a.require_table_primary_key_check'
> exec ^
: string
| set a.reset = 1
< error
| 1193 HY000 Unknown system variable 'a.reset'
> exec ^
: string
| set a.resignal = 1
< error
| 1193 HY000 Unknown system variable 'a.resignal'
> exec ^
: string
| set a.resource = 1
< error
| 1193 HY000 Unknown system variable 'a.resource'
> exec ^
: string
| set a.respect = 1
< error
| 1193 HY000 Unknown system variable 'a.respect'
> exec ^
: string
| set a.restart = 1
< error
| 1193 HY000 Unknown system variable 'a.restart'
> exec ^
: string
| set a.restore = 1
< error
| 1193 HY000 Unknown system variable 'a.restore'
> exec ^
: string
| set a.restrict = 1
< error
| 1193 HY000 Unknown system variable 'a.restrict'
> exec ^
: string
| set a.resume = 1
< error
| 1193 HY000 Unknown system variable 'a.resume'
> exec ^
: string
| set a.retain = 1
< error
| 1193 HY000 Unknown system variable 'a.retain'
> exec ^
: string
| set a.return = 1
< error
| 1193 HY000 Unknown system variable 'a.return'
> exec ^
: string
| set a.returned_sqlstate = 1
< error
| 1193 HY000 Unknown system variable 'a.returned_sqlstate'
> exec ^
: string
| set a.returning = 1
< error
| 1193 HY000 Unknown system variable 'a.returning'
> exec ^
: string
| set a.returns = 1
< error
| 1193 HY000 Unknown system variable 'a.returns'
> exec ^
: string
| set a.reuse = 1
< error
| 1193 HY000 Unknown system variable 'a.reuse'
> exec ^
: string
| set a.reverse = 1
< error
| 1193 HY000 Unknown system variable 'a.reverse'
> exec ^
: string
| set a.revoke = 1
< error
| 1193 HY000 Unknown system variable 'a.revoke'
> exec ^
: string
| set a.right = 1
< error
| 1193 HY000 Unknown system variable 'a.right'
> exec ^
: string
| set a.rlike = 1
< error
| 1193 HY000 Unknown system variable 'a.rlike'
> exec ^
: string
| set a.role = 1
< error
| 1193 HY000 Unknown system variable 'a.role'
> exec ^
: string
| set a.rollback = 1
< error
| 1193 HY000 Unknown system variable 'a.rollback'
> exec ^
: string
| set a.rollup = 1
< error
| 1193 HY000 Unknown system variable 'a.rollup'
> exec ^
: string
| set a.rotate = 1
< error
| 1193 HY000 Unknown system variable 'a.rotate'
> exec ^
: string
| set a.routine = 1
< error
| 1193 HY000 Unknown system variable 'a.routine'
> exec ^
: string
| set a.row = 1
< error
| 1193 HY000 Unknown system variable 'a.row'
> exec ^
: string
| set a.row_count = 1
< error
| 1193 HY000 Unknown system variable 'a.row_count'
> exec ^
: string
| set a.row_format = 1
< error
| 1193 HY000 Unknown system variable 'a.row_format'
> exec ^
: string
| set a.row_number = 1
< error
| 1193 HY000 Unknown system variable 'a.row_number'
> exec ^
: string
| set a.rows = 1
< error
| 1193 HY000 Unknown system variable 'a.rows'
> exec ^
: string
| set a.rtree = 1
< error
| 1193 HY000 Unknown system variable 'a.rtree'
> exec ^
: string
| set a.s3 = 1
< error
| 1193 HY000 Unknown system variable 'a.s3'
> exec ^
: string
| set a.savepoint = 1
< error
| 1193 HY000 Unknown system variable 'a.savepoint'
> exec ^
: string
| set a.schedule = 1
< error
| 1193 HY000 Unknown system variable 'a.schedule'
> exec ^
: string
| set a.schema = 1
< error
| 1193 HY000 Unknown system variable 'a.schema'
> exec ^
: string
| set a.schema_name = 1
< error
| 1193 HY000 Unknown system variable 'a.schema_name'
> exec ^
: string
| set a.schemas = 1
< error
| 1193 HY000 Unknown system variable 'a.schemas'
> exec ^
: string
| set a.second = 1
< error
| 1193 HY000 Unknown system variable 'a.second'
> exec ^
: string
| set a.second_microsecond = 1
< error
| 1193 HY000 Unknown system variable 'a.second_microsecond'
> exec ^
: string
| set a.secondary = 1
< error
| 1193 HY000 Unknown system variable 'a.secondary'
> exec ^
: string
| set a.secondary_engine = 1
< error
| 1193 HY000 Unknown system variable 'a.secondary_engine'
> exec ^
: string
| set a.secondary_engine_attribute = 1
< error
| 1193 HY000 Unknown system variable 'a.secondary_engine_attribute'
> exec ^
: string
| set a.secondary_load = 1
< error
| 1193 HY000 Unknown system variable 'a.secondary_load'
> exec ^
: string
| set a.secondary_unload = 1
< error
| 1193 HY000 Unknown system variable 'a.secondary_unload'
> exec ^
: string
| set a.security = 1
< error
| 1193 HY000 Unknown system variable 'a.security'
> exec ^
: string
| set a.select = 1
< error
| 1193 HY000 Unknown system variable 'a.select'
> exec ^
: string
| set a.sensitive = 1
< error
| 1193 HY000 Unknown system variable 'a.sensitive'
> exec ^
: string
| set a.separator = 1
< error
| 1193 HY000 Unknown system variable 'a.separator'
> exec ^
: string
| set a.serial = 1
< error
| 1193 HY000 Unknown system variable 'a.serial'
> exec ^
: string
| set a.serializable = 1
< error
| 1193 HY000 Unknown system variable 'a.serializable'
> exec ^
: string
| set a.server = 1
< error
| 1193 HY000 Unknown system variable 'a.server'
> exec ^
: string
| set a.session = 1
< error
| 1193 HY000 Unknown system variable 'a.session'
> exec ^
: string
| set a.session_user = 1
< error
| 1193 HY000 Unknown system variable 'a.session_user'
> exec ^
: string
| set a.set = 1
< error
| 1193 HY000 Unknown system variable 'a.set'
> exec ^
: string
| set a.sets = 1
< error
| 1193 HY000 Unknown system variable 'a.sets'
> exec ^
: string
| set a.share = 1
< error
| 1193 HY000 Unknown system variable 'a.share'
> exec ^
: string
| set a.show = 1
< error
| 1193 HY000 Unknown system variable 'a.show'
> exec ^
: string
| set a.shutdown = 1
< error
| 1193 HY000 Unknown system variable 'a.shutdown'
> exec ^
: string
| set a.signal = 1
< error
| 1193 HY000 Unknown system variable 'a.signal'
> exec ^
: string
| set a.signed = 1
< error
| 1193 HY000 Unknown system variable 'a.signed'
> exec ^
: string
| set a.simple = 1
< error
| 1193 HY000 Unknown system variable 'a.simple'
> exec ^
: string
| set a.skip = 1
< error
| 1193 HY000 Unknown system variable 'a.skip'
> exec ^
: string
| set a.slave = 1
< error
| 1193 HY000 Unknown system variable 'a.slave'
> exec ^
: string
| set a.slow = 1
< error
| 1193 HY000 Unknown system variable 'a.slow'
> exec ^
: string
| set a.smallint = 1
< error
| 1193 HY000 Unknown system variable 'a.smallint'
> exec ^
: string
| set a.snapshot = 1
< error
| 1193 HY000 Unknown system variable 'a.snapshot'
> exec ^
: string
| set a.socket = 1
< error
| 1193 HY000 Unknown system variable 'a.socket'
> exec ^
: string
| set a.some = 1
< error
| 1193 HY000 Unknown system variable 'a.some'
> exec ^
: string
| set a.soname = 1
< error
| 1193 HY000 Unknown system variable 'a.soname'
> exec ^
: string
| set a.sounds = 1
< error
| 1193 HY000 Unknown system variable 'a.sounds'
> exec ^
: string
| set a.source = 1
< error
| 1193 HY000 Unknown system variable 'a.source'
> exec ^
: string
| set a.source_auto_position = 1
< error
| 1193 HY000 Unknown system variable 'a.source_auto_position'
> exec ^
: string
| set a.source_bind = 1
< error
| 1193 HY000 Unknown system variable 'a.source_bind'
> exec ^
: string
| set a.source_compression_algorithms = 1
< error
| 1193 HY000 Unknown system variable 'a.source_compression_algorithms'
> exec ^
: string
| set a.source_connect_retry = 1
< error
| 1193 HY000 Unknown system variable 'a.source_connect_retry'
> exec ^
: string
| set a.source_connection_auto_failover = 1
< error
| 1193 HY000 Unknown system variable 'a.source_connection_auto_failover'
> exec ^
: string
| set a.source_delay = 1
< error
| 1193 HY000 Unknown system variable 'a.source_delay'
> exec ^
: string
| set a.source_heartbeat_period = 1
< error
| 1193 HY000 Unknown system variable 'a.source_heartbeat_period'
> exec ^
: string
| set a.source_host = 1
< error
| 1193 HY000 Unknown system variable 'a.source_host'
> exec ^
: string
| set a.source_log_file = 1
< error
| 1193 HY000 Unknown system variable 'a.source_log_file'
> exec ^
: string
| set a.source_log_pos = 1
< error
| 1193 HY000 Unknown system variable 'a.source_log_pos'
> exec ^
: string
| set a.source_password = 1
< error
| 1193 HY000 Unknown system variable 'a.source_password'
> exec ^
: string
| set a.source_port = 1
< error
| 1193 HY000 Unknown system variable 'a.source_port'
> exec ^
: string
| set a.source_public_key_path = 1
< error
| 1193 HY000 Unknown system variable 'a.source_public_key_path'
> exec ^
: string
| set a.source_retry_count = 1
< error
| 1193 HY000 Unknown system variable 'a.source_retry_count'
> exec ^
: string
| set a.source_ssl = 1
< error
| 1193 HY000 Unknown system variable 'a.source_ssl'
> exec ^
: string
| set a.source_ssl_ca = 1
< error
| 1193 HY000 Unknown system variable 'a.source_ssl_ca'
> exec ^
: string
| set a.source_ssl_capath = 1
< error
| 1193 HY000 Unknown system variable 'a.source_ssl_capath'
> exec ^
: string
| set a.source_ssl_cert = 1
< error
| 1193 HY000 Unknown system variable 'a.source_ssl_cert'
> exec ^
: string
| set a.source_ssl_cipher = 1
< error
| 1193 HY000 Unknown system variable 'a.source_ssl_cipher'
> exec ^
: string
| set a.source_ssl_crl = 1
< error
| 1193 HY000 Unknown system variable 'a.source_ssl_crl'
> exec ^
: string
| set a.source_ssl_crlpath = 1
< error
| 1193 HY000 Unknown system variable 'a.source_ssl_crlpath'
> exec ^
: string
| set a.source_ssl_key = 1
< error
| 1193 HY000 Unknown system variable 'a.source_ssl_key'
> exec ^
: string
| set a.source_ssl_verify_server_cert = 1
< error
| 1193 HY000 Unknown system variable 'a.source_ssl_verify_server_cert'
> exec ^
: string
| set a.source_tls_ciphersuites = 1
< error
| 1193 HY000 Unknown system variable 'a.source_tls_ciphersuites'
> exec ^
: string
| set a.source_tls_version = 1
< error
| 1193 HY000 Unknown system variable 'a.source_tls_version'
> exec ^
: string
| set a.source_user = 1
< error
| 1193 HY000 Unknown system variable 'a.source_user'
> exec ^
: string
| set a.source_zstd_compression_level = 1
< error
| 1193 HY000 Unknown system variable 'a.source_zstd_compression_level'
> exec ^
: string
| set a.spatial = 1
< error
| 1193 HY000 Unknown system variable 'a.spatial'
> exec ^
: string
| set a.specific = 1
< error
| 1193 HY000 Unknown system variable 'a.specific'
> exec ^
: string
| set a.sql = 1
< error
| 1193 HY000 Unknown system variable 'a.sql'
> exec ^
: string
| set a.sql_after_gtids = 1
< error
| 1193 HY000 Unknown system variable 'a.sql_after_gtids'
> exec ^
: string
| set a.sql_after_mts_gaps = 1
< error
| 1193 HY000 Unknown system variable 'a.sql_after_mts_gaps'
> exec ^
: string
| set a.sql_before_gtids = 1
< error
| 1193 HY000 Unknown system variable 'a.sql_before_gtids'
> exec ^
: string
| set a.sql_big_result = 1
< error
| 1193 HY000 Unknown system variable 'a.sql_big_result'
> exec ^
: string
| set a.sql_buffer_result = 1
< error
| 1193 HY000 Unknown system variable 'a.sql_buffer_result'
> exec ^
: string
| set a.sql_calc_found_rows = 1
< error
| 1193 HY000 Unknown system variable 'a.sql_calc_found_rows'
> exec ^
: string
| set a.sql_no_cache = 1
< error
| 1193 HY000 Unknown system variable 'a.sql_no_cache'
> exec ^
: string
| set a.sql_small_result = 1
< error
| 1193 HY000 Unknown system variable 'a.sql_small_result'
> exec ^
: string
| set a.sql_thread = 1
< error
| 1193 HY000 Unknown system variable 'a.sql_thread'
> exec ^
: string
| set a.sql_tsi_day = 1
< error
| 1193 HY000 Unknown system variable 'a.sql_tsi_day'
> exec ^
: string
| set a.sql_tsi_hour = 1
< error
| 1193 HY000 Unknown system variable 'a.sql_tsi_hour'
> exec ^
: string
| set a.sql_tsi_minute = 1
< error
| 1193 HY000 Unknown system variable 'a.sql_tsi_minute'
> exec ^
: string
| set a.sql_tsi_month = 1
< error
| 1193 HY000 Unknown system variable 'a.sql_tsi_month'
> exec ^
: string
| set a.sql_tsi_quarter = 1
< error
| 1193 HY000 Unknown system variable 'a.sql_tsi_quarter'
> exec ^
: string
| set a.sql_tsi_second = 1
< error
| 1193 HY000 Unknown system variable 'a.sql_tsi_second'
> exec ^
: string
| set a.sql_tsi_week = 1
< error
| 1193 HY000 Unknown system variable 'a.sql_tsi_week'
> exec ^
: string
| set a.sql_tsi_year = 1
< error
| 1193 HY000 Unknown system variable 'a.sql_tsi_year'
> exec ^
: string
| set a.sqlexception = 1
< error
| 1193 HY000 Unknown system variable 'a.sqlexception'
> exec ^
: string
| set a.sqlstate = 1
< error
| 1193 HY000 Unknown system variable 'a.sqlstate'
> exec ^
: string
| set a.sqlwarning = 1
< error
| 1193 HY000 Unknown system variable 'a.sqlwarning'
> exec ^
: string
| set a.srid = 1
< error
| 1193 HY000 Unknown system variable 'a.srid'
> exec ^
: string
| set a.sse_kms = 1
< error
| 1193 HY000 Unknown system variable 'a.sse_kms'
> exec ^
: string
| set a.sse_s3 = 1
< error
| 1193 HY000 Unknown system variable 'a.sse_s3'
> exec ^
: string
| set a.ssl = 1
< error
| 1193 HY000 Unknown system variable 'a.ssl'
> exec ^
: string
| set a.st_collect = 1
< error
| 1193 HY000 Unknown system variable 'a.st_collect'
> exec ^
: string
| set a.stacked = 1
< error
| 1193 HY000 Unknown system variable 'a.stacked'
> exec ^
: string
| set a.start = 1
< error
| 1193 HY000 Unknown system variable 'a.start'
> exec ^
: string
| set a.starting = 1
< error
| 1193 HY000 Unknown system variable 'a.starting'
> exec ^
: string
| set a.starts = 1
< error
| 1193 HY000 Unknown system variable 'a.starts'
> exec ^
: string
| set a.stats_auto_recalc = 1
< error
| 1193 HY000 Unknown system variable 'a.stats_auto_recalc'
> exec ^
: string
| set a.stats_persistent = 1
< error
| 1193 HY000 Unknown system variable 'a.stats_persistent'
> exec ^
: string
| set a.stats_sample_pages = 1
< error
| 1193 HY000 Unknown system variable 'a.stats_sample_pages'
> exec ^
: string
| set a.status = 1
< error
| 1193 HY000 Unknown system variable 'a.status'
> exec ^
: string
| set a.std = 1
< error
| 1193 HY000 Unknown system variable 'a.std'
> exec ^
: string
| set a.stddev = 1
< error
| 1193 HY000 Unknown system variable 'a.stddev'
> exec ^
: string
| set a.stddev_pop = 1
< error
| 1193 HY000 Unknown system variable 'a.stddev_pop'
> exec ^
: string
| set a.stddev_samp = 1
< error
| 1193 HY000 Unknown system variable 'a.stddev_samp'
> exec ^
: string
| set a.stop = 1
< error
| 1193 HY000 Unknown system variable 'a.stop'
> exec ^
: string
| set a.storage = 1
< error
| 1193 HY000 Unknown system variable 'a.storage'
> exec ^
: string
| set a.stored = 1
< error
| 1193 HY000 Unknown system variable 'a.stored'
> exec ^
: string
| set a.straight_join = 1
< error
| 1193 HY000 Unknown system variable 'a.straight_join'
> exec ^
: string
| set a.stream = 1
< error
| 1193 HY000 Unknown system variable 'a.stream'
> exec ^
: string
| set a.strict_load = 1
< error
| 1193 HY000 Unknown system variable 'a.strict_load'
> exec ^
: string
| set a.string = 1
< error
| 1193 HY000 Unknown system variable 'a.string'
> exec ^
: string
| set a.subclass_origin = 1
< error
| 1193 HY000 Unknown system variable 'a.subclass_origin'
> exec ^
: string
| set a.subdate = 1
< error
| 1193 HY000 Unknown system variable 'a.subdate'
> exec ^
: string
| set a.subject = 1
< error
| 1193 HY000 Unknown system variable 'a.subject'
> exec ^
: string
| set a.subpartition = 1
< error
| 1193 HY000 Unknown system variable 'a.subpartition'
> exec ^
: string
| set a.subpartitions = 1
< error
| 1193 HY000 Unknown system variable 'a.subpartitions'
> exec ^
: string
| set a.substr = 1
< error
| 1193 HY000 Unknown system variable 'a.substr'
> exec ^
: string
| set a.substring = 1
< error
| 1193 HY000 Unknown system variable 'a.substring'
> exec ^
: string
| set a.sum = 1
< error
| 1193 HY000 Unknown system variable 'a.sum'
> exec ^
: string
| set a.super = 1
< error
| 1193 HY000 Unknown system variable 'a.super'
> exec ^
: string
| set a.suspend = 1
< error
| 1193 HY000 Unknown system variable 'a.suspend'
> exec ^
: string
| set a.swaps = 1
< error
| 1193 HY000 Unknown system variable 'a.swaps'
> exec ^
: string
| set a.switches = 1
< error
| 1193 HY000 Unknown system variable 'a.switches'
> exec ^
: string
| set a.sysdate = 1
< error
| 1193 HY000 Unknown system variable 'a.sysdate'
> exec ^
: string
| set a.system = 1
< error
| 1193 HY000 Unknown system variable 'a.system'
> exec ^
: string
| set a.system_user = 1
< error
| 1193 HY000 Unknown system variable 'a.system_user'
> exec ^
: string
| set a.table = 1
< error
| 1193 HY000 Unknown system variable 'a.table'
> exec ^
: string
| set a.table_checksum = 1
< error
| 1193 HY000 Unknown system variable 'a.table_checksum'
> exec ^
: string
| set a.table_name = 1
< error
| 1193 HY000 Unknown system variable 'a.table_name'
> exec ^
: string
| set a.tables = 1
< error
| 1193 HY000 Unknown system variable 'a.tables'
> exec ^
: string
| set a.tablesample = 1
< error
| 1193 HY000 Unknown system variable 'a.tablesample'
> exec ^
: string
| set a.tablespace = 1
< error
| 1193 HY000 Unknown system variable 'a.tablespace'
> exec ^
: string
| set a.temporary = 1
< error
| 1193 HY000 Unknown system variable 'a.temporary'
> exec ^
: string
| set a.temptable = 1
< error
| 1193 HY000 Unknown system variable 'a.temptable'
> exec ^
: string
| set a.terminated = 1
< error
| 1193 HY000 Unknown system variable 'a.terminated'
> exec ^
: string
| set a.text = 1
< error
| 1193 HY000 Unknown system variable 'a.text'
> exec ^
: string
| set a.than = 1
< error
| 1193 HY000 Unknown system variable 'a.than'
> exec ^
: string
| set a.then = 1
< error
| 1193 HY000 Unknown system variable 'a.then'
> exec ^
: string
| set a.thread_priority = 1
< error
| 1193 HY000 Unknown system variable 'a.thread_priority'
> exec ^
: string
| set a.ties = 1
< error
| 1193 HY000 Unknown system variable 'a.ties'
> exec ^
: string
| set a.time = 1
< error
| 1193 HY000 Unknown system variable 'a.time'
> exec ^
: string
| set a.timestamp = 1
< error
| 1193 HY000 Unknown system variable 'a.timestamp'
> exec ^
: string
| set a.timestampadd = 1
< error
| 1193 HY000 Unknown system variable 'a.timestampadd'
> exec ^
: string
| set a.timestampdiff = 1
< error
| 1193 HY000 Unknown system variable 'a.timestampdiff'
> exec ^
: string
| set a.tinyblob = 1
< error
| 1193 HY000 Unknown system variable 'a.tinyblob'
> exec ^
: string
| set a.tinyint = 1
< error
| 1193 HY000 Unknown system variable 'a.tinyint'
> exec ^
: string
| set a.tinytext = 1
< error
| 1193 HY000 Unknown system variable 'a.tinytext'
> exec ^
: string
| set a.tls = 1
< error
| 1193 HY000 Unknown system variable 'a.tls'
> exec ^
: string
| set a.to = 1
< error
| 1193 HY000 Unknown system variable 'a.to'
> exec ^
: string
| set a.trailing = 1
< error
| 1193 HY000 Unknown system variable 'a.trailing'
> exec ^
: string
| set a.transaction = 1
< error
| 1193 HY000 Unknown system variable 'a.transaction'
> exec ^
: string
| set a.trigger = 1
< error
| 1193 HY000 Unknown system variable 'a.trigger'
> exec ^
: string
| set a.triggers = 1
< error
| 1193 HY000 Unknown system variable 'a.triggers'
> exec ^
: string
| set a.trim = 1
< error
| 1193 HY000 Unknown system variable 'a.trim'
> exec ^
: string
| set a.true = 1
< error
| 1193 HY000 Unknown system variable 'a.true'
> exec ^
: string
| set a.truncate = 1
< error
| 1193 HY000 Unknown system variable 'a.truncate'
> exec ^
: string
| set a.type = 1
< error
| 1193 HY000 Unknown system variable 'a.type'
> exec ^
: string
| set a.types = 1
< error
| 1193 HY000 Unknown system variable 'a.types'
> exec ^
: string
| set a.unbounded = 1
< error
| 1193 HY000 Unknown system variable 'a.unbounded'
> exec ^
: string
| set a.uncommitted = 1
< error
| 1193 HY000 Unknown system variable 'a.uncommitted'
> exec ^
: string
| set a.undefined = 1
< error
| 1193 HY000 Unknown system variable 'a.undefined'
> exec ^
: string
| set a.undo = 1
< error
| 1193 HY000 Unknown system variable 'a.undo'
> exec ^
: string
| set a.undo_buffer_size = 1
< error
| 1193 HY000 Unknown system variable 'a.undo_buffer_size'
> exec ^
: string
| set a.undofile = 1
< error
| 1193 HY000 Unknown system variable 'a.undofile'
> exec ^
: string
| set a.unicode = 1
< error
| 1193 HY000 Unknown system variable 'a.unicode'
> exec ^
: string
| set a.uninstall = 1
< error
| 1193 HY000 Unknown system variable 'a.uninstall'
> exec ^
: string
| set a.union = 1
< error
| 1193 HY000 Unknown system variable 'a.union'
> exec ^
: string
| set a.unique = 1
< error
| 1193 HY000 Unknown system variable 'a.unique'
> exec ^
: string
| set a.unknown = 1
< error
| 1193 HY000 Unknown system variable 'a.unknown'
> exec ^
: string
| set a.unlock = 1
< error
| 1193 HY000 Unknown system variable 'a.unlock'
> exec ^
: string
| set a.unregister = 1
< error
| 1193 HY000 Unknown system variable 'a.unregister'
> exec ^
: string
| set a.unsigned = 1
< error
| 1193 HY000 Unknown system variable 'a.unsigned'
> exec ^
: string
| set a.until = 1
< error
| 1193 HY000 Unknown system variable 'a.until'
> exec ^
: string
| set a.update = 1
< error
| 1193 HY000 Unknown system variable 'a.update'
> exec ^
: string
| set a.upgrade = 1
< error
| 1193 HY000 Unknown system variable 'a.upgrade'
> exec ^
: string
| set a.uri = 1
< error
| 1193 HY000 Unknown system variable 'a.uri'
> exec ^
: string
| set a.url = 1
< error
| 1193 HY000 Unknown system variable 'a.url'
> exec ^
: string
| set a.usage = 1
< error
| 1193 HY000 Unknown system variable 'a.usage'
> exec ^
: string
| set a.use = 1
< error
| 1193 HY000 Unknown system variable 'a.use'
> exec ^
: string
| set a.use_frm = 1
< error
| 1193 HY000 Unknown system variable 'a.use_frm'
> exec ^
: string
| set a.user = 1
< error
| 1193 HY000 Unknown system variable 'a.user'
> exec ^
: string
| set a.user_resources = 1
< error
| 1193 HY000 Unknown system variable 'a.user_resources'
> exec ^
: string
| set a.using = 1
< error
| 1193 HY000 Unknown system variable 'a.using'
> exec ^
: string
| set a.utc_date = 1
< error
| 1193 HY000 Unknown system variable 'a.utc_date'
> exec ^
: string
| set a.utc_time = 1
< error
| 1193 HY000 Unknown system variable 'a.utc_time'
> exec ^
: string
| set a.utc_timestamp = 1
< error
| 1193 HY000 Unknown system variable 'a.utc_timestamp'
> exec ^
: string
| set a.validate = 1
< error
| 1193 HY000 Unknown system variable 'a.validate'
> exec ^
: string
| set a.validation = 1
< error
| 1193 HY000 Unknown system variable 'a.validation'
> exec ^
: string
| set a.value = 1
< error
| 1193 HY000 Unknown system variable 'a.value'
> exec ^
: string
| set a.values = 1
< error
| 1193 HY000 Unknown system variable 'a.values'
> exec ^
: string
| set a.var_pop = 1
< error
| 1193 HY000 Unknown system variable 'a.var_pop'
> exec ^
: string
| set a.var_samp = 1
< error
| 1193 HY000 Unknown system variable 'a.var_samp'
> exec ^
: string
| set a.varbinary = 1
< error
| 1193 HY000 Unknown system variable 'a.varbinary'
> exec ^
: string
| set a.varchar = 1
< error
| 1193 HY000 Unknown system variable 'a.varchar'
> exec ^
: string
| set a.varcharacter = 1
< error
| 1193 HY000 Unknown system variable 'a.varcharacter'
> exec ^
: string
| set a.variables = 1
< error
| 1193 HY000 Unknown system variable 'a.variables'
> exec ^
: string
| set a.variance = 1
< error
| 1193 HY000 Unknown system variable 'a.variance'
> exec ^
: string
| set a.varying = 1
< error
| 1193 HY000 Unknown system variable 'a.varying'
> exec ^
: string
| set a.vcpu = 1
< error
| 1193 HY000 Unknown system variable 'a.vcpu'
> exec ^
: string
| set a.vector = 1
< error
| 1193 HY000 Unknown system variable 'a.vector'
> exec ^
: string
| set a.verify_key_constraints = 1
< error
| 1193 HY000 Unknown system variable 'a.verify_key_constraints'
> exec ^
: string
| set a.view = 1
< error
| 1193 HY000 Unknown system variable 'a.view'
> exec ^
: string
| set a.virtual = 1
< error
| 1193 HY000 Unknown system variable 'a.virtual'
> exec ^
: string
| set a.visible = 1
< error
| 1193 HY000 Unknown system variable 'a.visible'
> exec ^
: string
| set a.wait = 1
< error
| 1193 HY000 Unknown system variable 'a.wait'
> exec ^
: string
| set a.warnings = 1
< error
| 1193 HY000 Unknown system variable 'a.warnings'
> exec ^
: string
| set a.week = 1
< error
| 1193 HY000 Unknown system variable 'a.week'
> exec ^
: string
| set a.weight_string = 1
< error
| 1193 HY000 Unknown system variable 'a.weight_string'
> exec ^
: string
| set a.when = 1
< error
| 1193 HY000 Unknown system variable 'a.when'
> exec ^
: string
| set a.where = 1
< error
| 1193 HY000 Unknown system variable 'a.where'
> exec ^
: string
| set a.while = 1
< error
| 1193 HY000 Unknown system variable 'a.while'
> exec ^
: string
| set a.window = 1
< error
| 1193 HY000 Unknown system variable 'a.window'
> exec ^
: string
| set a.with = 1
< error
| 1193 HY000 Unknown system variable 'a.with'
> exec ^
: string
| set a.without = 1
< error
| 1193 HY000 Unknown system variable 'a.without'
> exec ^
: string
| set a.work = 1
< error
| 1193 HY000 Unknown system variable 'a.work'
> exec ^
: string
| set a.wrapper = 1
< error
| 1193 HY000 Unknown system variable 'a.wrapper'
> exec ^
: string
| set a.write = 1
< error
| 1193 HY000 Unknown system variable 'a.write'
> exec ^
: string
| set a.x509 = 1
< error
| 1193 HY000 Unknown system variable 'a.x509'
> exec ^
: string
| set a.xa = 1
< error
| 1193 HY000 Unknown system variable 'a.xa'
> exec ^
: string
| set a.xid = 1
< error
| 1193 HY000 Unknown system variable 'a.xid'
> exec ^
: string
| set a.xml = 1
< error
| 1193 HY000 Unknown system variable 'a.xml'
> exec ^
: string
| set a.xor = 1
< error
| 1193 HY000 Unknown system variable 'a.xor'
> exec ^
: string
| set a.year = 1
< error
| 1193 HY000 Unknown system variable 'a.year'
> exec ^
: string
| set a.year_month = 1
< error
| 1193 HY000 Unknown system variable 'a.year_month'
> exec ^
: string
| set a.zerofill = 1
< error
| 1193 HY000 Unknown system variable 'a.zerofill'
> exec ^
: string
| set a.zone = 1
< error
| 1193 HY000 Unknown system variable 'a.zone'
> exec ^
: string
| set @a = 1, absent = 2
< error
| 1193 HY000 Unknown system variable 'absent'
> exec ^
: string
| set @a = 1, accessible = 2
< error
| 1064 42000 …near 'accessible = 2' at line 1
> exec ^
: string
| set @a = 1, account = 2
< error
| 1193 HY000 Unknown system variable 'account'
> exec ^
: string
| set @a = 1, action = 2
< error
| 1193 HY000 Unknown system variable 'action'
> exec ^
: string
| set @a = 1, active = 2
< error
| 1193 HY000 Unknown system variable 'active'
> exec ^
: string
| set @a = 1, add = 2
< error
| 1064 42000 …near 'add = 2' at line 1
> exec ^
: string
| set @a = 1, adddate = 2
< error
| 1193 HY000 Unknown system variable 'adddate'
> exec ^
: string
| set @a = 1, admin = 2
< error
| 1193 HY000 Unknown system variable 'admin'
> exec ^
: string
| set @a = 1, after = 2
< error
| 1193 HY000 Unknown system variable 'after'
> exec ^
: string
| set @a = 1, against = 2
< error
| 1193 HY000 Unknown system variable 'against'
> exec ^
: string
| set @a = 1, aggregate = 2
< error
| 1193 HY000 Unknown system variable 'aggregate'
> exec ^
: string
| set @a = 1, algorithm = 2
< error
| 1193 HY000 Unknown system variable 'algorithm'
> exec ^
: string
| set @a = 1, all = 2
< error
| 1064 42000 …near 'all = 2' at line 1
> exec ^
: string
| set @a = 1, allow_missing_files = 2
< error
| 1193 HY000 Unknown system variable 'allow_missing_files'
> exec ^
: string
| set @a = 1, alter = 2
< error
| 1064 42000 …near 'alter = 2' at line 1
> exec ^
: string
| set @a = 1, always = 2
< error
| 1193 HY000 Unknown system variable 'always'
> exec ^
: string
| set @a = 1, analyze = 2
< error
| 1064 42000 …near 'analyze = 2' at line 1
> exec ^
: string
| set @a = 1, and = 2
< error
| 1064 42000 …near 'and = 2' at line 1
> exec ^
: string
| set @a = 1, any = 2
< error
| 1193 HY000 Unknown system variable 'any'
> exec ^
: string
| set @a = 1, array = 2
< error
| 1193 HY000 Unknown system variable 'array'
> exec ^
: string
| set @a = 1, as = 2
< error
| 1064 42000 …near 'as = 2' at line 1
> exec ^
: string
| set @a = 1, asc = 2
< error
| 1064 42000 …near 'asc = 2' at line 1
> exec ^
: string
| set @a = 1, ascii = 2
< error
| 1193 HY000 Unknown system variable 'ascii'
> exec ^
: string
| set @a = 1, asensitive = 2
< error
| 1064 42000 …near 'asensitive = 2' at line 1
> exec ^
: string
| set @a = 1, assign_gtids_to_anonymous_transactions = 2
< error
| 1193 HY000 Unknown system variable 'assign_gtids_to_anonymous_transactions'
> exec ^
: string
| set @a = 1, at = 2
< error
| 1193 HY000 Unknown system variable 'at'
> exec ^
: string
| set @a = 1, attribute = 2
< error
| 1193 HY000 Unknown system variable 'attribute'
> exec ^
: string
| set @a = 1, authentication = 2
< error
| 1193 HY000 Unknown system variable 'authentication'
> exec ^
: string
| set @a = 1, auto = 2
< error
| 1193 HY000 Unknown system variable 'auto'
> exec ^
: string
| set @a = 1, auto_increment = 2
< error
| 1193 HY000 Unknown system variable 'auto_increment'
> exec ^
: string
| set @a = 1, auto_refresh = 2
< error
| 1193 HY000 Unknown system variable 'auto_refresh'
> exec ^
: string
| set @a = 1, auto_refresh_source = 2
< error
| 1193 HY000 Unknown system variable 'auto_refresh_source'
> exec ^
: string
| set @a = 1, autoextend_size = 2
< error
| 1193 HY000 Unknown system variable 'autoextend_size'
> exec ^
: string
| set @a = 1, avg = 2
< error
| 1193 HY000 Unknown system variable 'avg'
> exec ^
: string
| set @a = 1, avg_row_length = 2
< error
| 1193 HY000 Unknown system variable 'avg_row_length'
> exec ^
: string
| set @a = 1, backup = 2
< error
| 1193 HY000 Unknown system variable 'backup'
> exec ^
: string
| set @a = 1, before = 2
< error
| 1064 42000 …near 'before = 2' at line 1
> exec ^
: string
| set @a = 1, begin = 2
< error
| 1193 HY000 Unknown system variable 'begin'
> exec ^
: string
| set @a = 1, bernoulli = 2
< error
| 1193 HY000 Unknown system variable 'bernoulli'
> exec ^
: string
| set @a = 1, between = 2
< error
| 1064 42000 …near 'between = 2' at line 1
> exec ^
: string
| set @a = 1, bigint = 2
< error
| 1064 42000 …near 'bigint = 2' at line 1
> exec ^
: string
| set @a = 1, binary = 2
< error
| 1064 42000 …near 'binary = 2' at line 1
> exec ^
: string
| set @a = 1, binlog = 2
< error
| 1193 HY000 Unknown system variable 'binlog'
> exec ^
: string
| set @a = 1, bit = 2
< error
| 1193 HY000 Unknown system variable 'bit'
> exec ^
: string
| set @a = 1, bit_and = 2
< error
| 1193 HY000 Unknown system variable 'bit_and'
> exec ^
: string
| set @a = 1, bit_or = 2
< error
| 1193 HY000 Unknown system variable 'bit_or'
> exec ^
: string
| set @a = 1, bit_xor = 2
< error
| 1193 HY000 Unknown system variable 'bit_xor'
> exec ^
: string
| set @a = 1, blob = 2
< error
| 1064 42000 …near 'blob = 2' at line 1
> exec ^
: string
| set @a = 1, block = 2
< error
| 1193 HY000 Unknown system variable 'block'
> exec ^
: string
| set @a = 1, bool = 2
< error
| 1193 HY000 Unknown system variable 'bool'
> exec ^
: string
| set @a = 1, boolean = 2
< error
| 1193 HY000 Unknown system variable 'boolean'
> exec ^
: string
| set @a = 1, both = 2
< error
| 1064 42000 …near 'both = 2' at line 1
> exec ^
: string
| set @a = 1, btree = 2
< error
| 1193 HY000 Unknown system variable 'btree'
> exec ^
: string
| set @a = 1, buckets = 2
< error
| 1193 HY000 Unknown system variable 'buckets'
> exec ^
: string
| set @a = 1, bulk = 2
< error
| 1193 HY000 Unknown system variable 'bulk'
> exec ^
: string
| set @a = 1, by = 2
< error
| 1064 42000 …near 'by = 2' at line 1
> exec ^
: string
| set @a = 1, byte = 2
< error
| 1193 HY000 Unknown system variable 'byte'
> exec ^
: string
| set @a = 1, cache = 2
< error
| 1193 HY000 Unknown system variable 'cache'
> exec ^
: string
| set @a = 1, call = 2
< error
| 1064 42000 …near 'call = 2' at line 1
> exec ^
: string
| set @a = 1, cascade = 2
< error
| 1064 42000 …near 'cascade = 2' at line 1
> exec ^
: string
| set @a = 1, cascaded = 2
< error
| 1193 HY000 Unknown system variable 'cascaded'
> exec ^
: string
| set @a = 1, case = 2
< error
| 1064 42000 …near 'case = 2' at line 1
> exec ^
: string
| set @a = 1, cast = 2
< error
| 1193 HY000 Unknown system variable 'cast'
> exec ^
: string
| set @a = 1, catalog_name = 2
< error
| 1193 HY000 Unknown system variable 'catalog_name'
> exec ^
: string
| set @a = 1, chain = 2
< error
| 1193 HY000 Unknown system variable 'chain'
> exec ^
: string
| set @a = 1, challenge_response = 2
< error
| 1193 HY000 Unknown system variable 'challenge_response'
> exec ^
: string
| set @a = 1, change = 2
< error
| 1064 42000 …near 'change = 2' at line 1
> exec ^
: string
| set @a = 1, changed = 2
< error
| 1193 HY000 Unknown system variable 'changed'
> exec ^
: string
| set @a = 1, channel = 2
< error
| 1193 HY000 Unknown system variable 'channel'
> exec ^
: string
| set @a = 1, char = 2
< error
| 1064 42000 …near '= 2' at line 1
> exec ^
: string
| set @a = 1, character = 2
< error
| 1064 42000 …near '= 2' at line 1
> exec ^
: string
| set @a = 1, charset = 2
< error
| 1193 HY000 Unknown system variable 'charset'
> exec ^
: string
| set @a = 1, check = 2
< error
| 1064 42000 …near 'check = 2' at line 1
> exec ^
: string
| set @a = 1, checksum = 2
< error
| 1193 HY000 Unknown system variable 'checksum'
> exec ^
: string
| set @a = 1, cipher = 2
< error
| 1193 HY000 Unknown system variable 'cipher'
> exec ^
: string
| set @a = 1, class_origin = 2
< error
| 1193 HY000 Unknown system variable 'class_origin'
> exec ^
: string
| set @a = 1, client = 2
< error
| 1193 HY000 Unknown system variable 'client'
> exec ^
: string
| set @a = 1, clone = 2
< error
| 1193 HY000 Unknown system variable 'clone'
> exec ^
: string
| set @a = 1, close = 2
< error
| 1193 HY000 Unknown system variable 'close'
> exec ^
: string
| set @a = 1, coalesce = 2
< error
| 1193 HY000 Unknown system variable 'coalesce'
> exec ^
: string
| set @a = 1, code = 2
< error
| 1193 HY000 Unknown system variable 'code'
> exec ^
: string
| set @a = 1, collate = 2
< error
| 1064 42000 …near 'collate = 2' at line 1
> exec ^
: string
| set @a = 1, collation = 2
< error
| 1193 HY000 Unknown system variable 'collation'
> exec ^
: string
| set @a = 1, column = 2
< error
| 1064 42000 …near 'column = 2' at line 1
> exec ^
: string
| set @a = 1, column_format = 2
< error
| 1193 HY000 Unknown system variable 'column_format'
> exec ^
: string
| set @a = 1, column_name = 2
< error
| 1193 HY000 Unknown system variable 'column_name'
> exec ^
: string
| set @a = 1, columns = 2
< error
| 1193 HY000 Unknown system variable 'columns'
> exec ^
: string
| set @a = 1, comment = 2
< error
| 1193 HY000 Unknown system variable 'comment'
> exec ^
: string
| set @a = 1, commit = 2
< error
| 1193 HY000 Unknown system variable 'commit'
> exec ^
: string
| set @a = 1, committed = 2
< error
| 1193 HY000 Unknown system variable 'committed'
> exec ^
: string
| set @a = 1, compact = 2
< error
| 1193 HY000 Unknown system variable 'compact'
> exec ^
: string
| set @a = 1, completion = 2
< error
| 1193 HY000 Unknown system variable 'completion'
> exec ^
: string
| set @a = 1, component = 2
< error
| 1193 HY000 Unknown system variable 'component'
> exec ^
: string
| set @a = 1, compressed = 2
< error
| 1193 HY000 Unknown system variable 'compressed'
> exec ^
: string
| set @a = 1, compression = 2
< error
| 1193 HY000 Unknown system variable 'compression'
> exec ^
: string
| set @a = 1, concat = 2
< error
| 1193 HY000 Unknown system variable 'concat'
> exec ^
: string
| set @a = 1, concurrent = 2
< error
| 1193 HY000 Unknown system variable 'concurrent'
> exec ^
: string
| set @a = 1, condition = 2
< error
| 1064 42000 …near 'condition = 2' at line 1
> exec ^
: string
| set @a = 1, connection = 2
< error
| 1193 HY000 Unknown system variable 'connection'
> exec ^
: string
| set @a = 1, consistent = 2
< error
| 1193 HY000 Unknown system variable 'consistent'
> exec ^
: string
| set @a = 1, constraint = 2
< error
| 1064 42000 …near 'constraint = 2' at line 1
> exec ^
: string
| set @a = 1, constraint_catalog = 2
< error
| 1193 HY000 Unknown system variable 'constraint_catalog'
> exec ^
: string
| set @a = 1, constraint_name = 2
< error
| 1193 HY000 Unknown system variable 'constraint_name'
> exec ^
: string
| set @a = 1, constraint_schema = 2
< error
| 1193 HY000 Unknown system variable 'constraint_schema'
> exec ^
: string
| set @a = 1, contains = 2
< error
| 1193 HY000 Unknown system variable 'contains'
> exec ^
: string
| set @a = 1, context = 2
< error
| 1193 HY000 Unknown system variable 'context'
> exec ^
: string
| set @a = 1, continue = 2
< error
| 1064 42000 …near 'continue = 2' at line 1
> exec ^
: string
| set @a = 1, convert = 2
< error
| 1064 42000 …near 'convert = 2' at line 1
> exec ^
: string
| set @a = 1, count = 2
< error
| 1193 HY000 Unknown system variable 'count'
> exec ^
: string
| set @a = 1, cpu = 2
< error
| 1193 HY000 Unknown system variable 'cpu'
> exec ^
: string
| set @a = 1, create = 2
< error
| 1064 42000 …near 'create = 2' at line 1
> exec ^
: string
| set @a = 1, cross = 2
< error
| 1064 42000 …near 'cross = 2' at line 1
> exec ^
: string
| set @a = 1, csv = 2
< error
| 1193 HY000 Unknown system variable 'csv'
> exec ^
: string
| set @a = 1, cube = 2
< error
| 1064 42000 …near 'cube = 2' at line 1
> exec ^
: string
| set @a = 1, cume_dist = 2
< error
| 1064 42000 …near 'cume_dist = 2' at line 1
> exec ^
: string
| set @a = 1, curdate = 2
< error
| 1193 HY000 Unknown system variable 'curdate'
> exec ^
: string
| set @a = 1, current = 2
< error
| 1193 HY000 Unknown system variable 'current'
> exec ^
: string
| set @a = 1, current_date = 2
< error
| 1064 42000 …near 'current_date = 2' at line 1
> exec ^
: string
| set @a = 1, current_time = 2
< error
| 1064 42000 …near 'current_time = 2' at line 1
> exec ^
: string
| set @a = 1, current_timestamp = 2
< error
| 1064 42000 …near 'current_timestamp = 2' at line 1
> exec ^
: string
| set @a = 1, current_user = 2
< error
| 1064 42000 …near 'current_user = 2' at line 1
> exec ^
: string
| set @a = 1, cursor = 2
< error
| 1064 42000 …near 'cursor = 2' at line 1
> exec ^
: string
| set @a = 1, cursor_name = 2
< error
| 1193 HY000 Unknown system variable 'cursor_name'
> exec ^
: string
| set @a = 1, curtime = 2
< error
| 1193 HY000 Unknown system variable 'curtime'
> exec ^
: string
| set @a = 1, data = 2
< error
| 1193 HY000 Unknown system variable 'data'
> exec ^
: string
| set @a = 1, database = 2
< error
| 1064 42000 …near 'database = 2' at line 1
> exec ^
: string
| set @a = 1, databases = 2
< error
| 1064 42000 …near 'databases = 2' at line 1
> exec ^
: string
| set @a = 1, datafile = 2
< error
| 1193 HY000 Unknown system variable 'datafile'
> exec ^
: string
| set @a = 1, date = 2
< error
| 1193 HY000 Unknown system variable 'date'
> exec ^
: string
| set @a = 1, date_add = 2
< error
| 1193 HY000 Unknown system variable 'date_add'
> exec ^
: string
| set @a = 1, date_sub = 2
< error
| 1193 HY000 Unknown system variable 'date_sub'
> exec ^
: string
| set @a = 1, datetime = 2
< error
| 1193 HY000 Unknown system variable 'datetime'
> exec ^
: string
| set @a = 1, day = 2
< error
| 1193 HY000 Unknown system variable 'day'
> exec ^
: string
| set @a = 1, day_hour = 2
< error
| 1064 42000 …near 'day_hour = 2' at line 1
> exec ^
: string
| set @a = 1, day_microsecond = 2
< error
| 1064 42000 …near 'day_microsecond = 2' at line 1
> exec ^
: string
| set @a = 1, day_minute = 2
< error
| 1064 42000 …near 'day_minute = 2' at line 1
> exec ^
: string
| set @a = 1, day_second = 2
< error
| 1064 42000 …near 'day_second = 2' at line 1
> exec ^
: string
| set @a = 1, deallocate = 2
< error
| 1193 HY000 Unknown system variable 'deallocate'
> exec ^
: string
| set @a = 1, dec = 2
< error
| 1064 42000 …near 'dec = 2' at line 1
> exec ^
: string
| set @a = 1, decimal = 2
< error
| 1064 42000 …near 'decimal = 2' at line 1
> exec ^
: string
| set @a = 1, declare = 2
< error
| 1064 42000 …near 'declare = 2' at line 1
> exec ^
: string
| set @a = 1, default = 2
< error
| 1064 42000 …near '= 2' at line 1
> exec ^
: string
| set @a = 1, default_auth = 2
< error
| 1193 HY000 Unknown system variable 'default_auth'
> exec ^
: string
| set @a = 1, definer = 2
< error
| 1193 HY000 Unknown system variable 'definer'
> exec ^
: string
| set @a = 1, definition = 2
< error
| 1193 HY000 Unknown system variable 'definition'
> exec ^
: string
| set @a = 1, delay_key_write = 2
< error
| 1229 HY000 Variable 'delay_key_write' is a GLOBAL variable and should be set with SET GLOBAL
> exec ^
: string
| set @a = 1, delay_key_write = 2
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set @a = 1, delayed = 2
< error
| 1064 42000 …near 'delayed = 2' at line 1
> exec ^
: string
| set @a = 1, delete = 2
< error
| 1064 42000 …near 'delete = 2' at line 1
> exec ^
: string
| set @a = 1, dense_rank = 2
< error
| 1064 42000 …near 'dense_rank = 2' at line 1
> exec ^
: string
| set @a = 1, desc = 2
< error
| 1064 42000 …near 'desc = 2' at line 1
> exec ^
: string
| set @a = 1, describe = 2
< error
| 1064 42000 …near 'describe = 2' at line 1
> exec ^
: string
| set @a = 1, description = 2
< error
| 1193 HY000 Unknown system variable 'description'
> exec ^
: string
| set @a = 1, deterministic = 2
< error
| 1064 42000 …near 'deterministic = 2' at line 1
> exec ^
: string
| set @a = 1, diagnostics = 2
< error
| 1193 HY000 Unknown system variable 'diagnostics'
> exec ^
: string
| set @a = 1, directory = 2
< error
| 1193 HY000 Unknown system variable 'directory'
> exec ^
: string
| set @a = 1, disable = 2
< error
| 1193 HY000 Unknown system variable 'disable'
> exec ^
: string
| set @a = 1, discard = 2
< error
| 1193 HY000 Unknown system variable 'discard'
> exec ^
: string
| set @a = 1, disk = 2
< error
| 1193 HY000 Unknown system variable 'disk'
> exec ^
: string
| set @a = 1, distinct = 2
< error
| 1064 42000 …near 'distinct = 2' at line 1
> exec ^
: string
| set @a = 1, distinctrow = 2
< error
| 1064 42000 …near 'distinctrow = 2' at line 1
> exec ^
: string
| set @a = 1, div = 2
< error
| 1064 42000 …near 'div = 2' at line 1
> exec ^
: string
| set @a = 1, do = 2
< error
| 1193 HY000 Unknown system variable 'do'
> exec ^
: string
| set @a = 1, double = 2
< error
| 1064 42000 …near 'double = 2' at line 1
> exec ^
: string
| set @a = 1, drop = 2
< error
| 1064 42000 …near 'drop = 2' at line 1
> exec ^
: string
| set @a = 1, dual = 2
< error
| 1064 42000 …near 'dual = 2' at line 1
> exec ^
: string
| set @a = 1, duality = 2
< error
| 1193 HY000 Unknown system variable 'duality'
> exec ^
: string
| set @a = 1, dumpfile = 2
< error
| 1193 HY000 Unknown system variable 'dumpfile'
> exec ^
: string
| set @a = 1, duplicate = 2
< error
| 1193 HY000 Unknown system variable 'duplicate'
> exec ^
: string
| set @a = 1, dynamic = 2
< error
| 1193 HY000 Unknown system variable 'dynamic'
> exec ^
: string
| set @a = 1, each = 2
< error
| 1064 42000 …near 'each = 2' at line 1
> exec ^
: string
| set @a = 1, else = 2
< error
| 1064 42000 …near 'else = 2' at line 1
> exec ^
: string
| set @a = 1, elseif = 2
< error
| 1064 42000 …near 'elseif = 2' at line 1
> exec ^
: string
| set @a = 1, empty = 2
< error
| 1064 42000 …near 'empty = 2' at line 1
> exec ^
: string
| set @a = 1, enable = 2
< error
| 1193 HY000 Unknown system variable 'enable'
> exec ^
: string
| set @a = 1, enclosed = 2
< error
| 1064 42000 …near 'enclosed = 2' at line 1
> exec ^
: string
| set @a = 1, encryption = 2
< error
| 1193 HY000 Unknown system variable 'encryption'
> exec ^
: string
| set @a = 1, end = 2
< error
| 1193 HY000 Unknown system variable 'end'
> exec ^
: string
| set @a = 1, ends = 2
< error
| 1193 HY000 Unknown system variable 'ends'
> exec ^
: string
| set @a = 1, enforced = 2
< error
| 1193 HY000 Unknown system variable 'enforced'
> exec ^
: string
| set @a = 1, engine = 2
< error
| 1193 HY000 Unknown system variable 'engine'
> exec ^
: string
| set @a = 1, engine_attribute = 2
< error
| 1193 HY000 Unknown system variable 'engine_attribute'
> exec ^
: string
| set @a = 1, engines = 2
< error
| 1193 HY000 Unknown system variable 'engines'
> exec ^
: string
| set @a = 1, enum = 2
< error
| 1193 HY000 Unknown system variable 'enum'
> exec ^
: string
| set @a = 1, error = 2
< error
| 1193 HY000 Unknown system variable 'error'
> exec ^
: string
| set @a = 1, errors = 2
< error
| 1193 HY000 Unknown system variable 'errors'
> exec ^
: string
| set @a = 1, escape = 2
< error
| 1193 HY000 Unknown system variable 'escape'
> exec ^
: string
| set @a = 1, escaped = 2
< error
| 1064 42000 …near 'escaped = 2' at line 1
> exec ^
: string
| set @a = 1, event = 2
< error
| 1193 HY000 Unknown system variable 'event'
> exec ^
: string
| set @a = 1, events = 2
< error
| 1193 HY000 Unknown system variable 'events'
> exec ^
: string
| set @a = 1, every = 2
< error
| 1193 HY000 Unknown system variable 'every'
> exec ^
: string
| set @a = 1, except = 2
< error
| 1064 42000 …near 'except = 2' at line 1
> exec ^
: string
| set @a = 1, exchange = 2
< error
| 1193 HY000 Unknown system variable 'exchange'
> exec ^
: string
| set @a = 1, exclude = 2
< error
| 1193 HY000 Unknown system variable 'exclude'
> exec ^
: string
| set @a = 1, execute = 2
< error
| 1193 HY000 Unknown system variable 'execute'
> exec ^
: string
| set @a = 1, exists = 2
< error
| 1064 42000 …near 'exists = 2' at line 1
> exec ^
: string
| set @a = 1, exit = 2
< error
| 1064 42000 …near 'exit = 2' at line 1
> exec ^
: string
| set @a = 1, expansion = 2
< error
| 1193 HY000 Unknown system variable 'expansion'
> exec ^
: string
| set @a = 1, expire = 2
< error
| 1193 HY000 Unknown system variable 'expire'
> exec ^
: string
| set @a = 1, explain = 2
< error
| 1064 42000 …near 'explain = 2' at line 1
> exec ^
: string
| set @a = 1, export = 2
< error
| 1193 HY000 Unknown system variable 'export'
> exec ^
: string
| set @a = 1, extended = 2
< error
| 1193 HY000 Unknown system variable 'extended'
> exec ^
: string
| set @a = 1, extent_size = 2
< error
| 1193 HY000 Unknown system variable 'extent_size'
> exec ^
: string
| set @a = 1, external = 2
< error
| 1193 HY000 Unknown system variable 'external'
> exec ^
: string
| set @a = 1, external_format = 2
< error
| 1193 HY000 Unknown system variable 'external_format'
> exec ^
: string
| set @a = 1, extract = 2
< error
| 1193 HY000 Unknown system variable 'extract'
> exec ^
: string
| set @a = 1, factor = 2
< error
| 1193 HY000 Unknown system variable 'factor'
> exec ^
: string
| set @a = 1, failed_login_attempts = 2
< error
| 1193 HY000 Unknown system variable 'failed_login_attempts'
> exec ^
: string
| set @a = 1, false = 2
< error
| 1064 42000 …near 'false = 2' at line 1
> exec ^
: string
| set @a = 1, fast = 2
< error
| 1193 HY000 Unknown system variable 'fast'
> exec ^
: string
| set @a = 1, faults = 2
< error
| 1193 HY000 Unknown system variable 'faults'
> exec ^
: string
| set @a = 1, fetch = 2
< error
| 1064 42000 …near 'fetch = 2' at line 1
> exec ^
: string
| set @a = 1, fields = 2
< error
| 1193 HY000 Unknown system variable 'fields'
> exec ^
: string
| set @a = 1, file = 2
< error
| 1193 HY000 Unknown system variable 'file'
> exec ^
: string
| set @a = 1, file_block_size = 2
< error
| 1193 HY000 Unknown system variable 'file_block_size'
> exec ^
: string
| set @a = 1, file_format = 2
< error
| 1193 HY000 Unknown system variable 'file_format'
> exec ^
: string
| set @a = 1, file_name = 2
< error
| 1193 HY000 Unknown system variable 'file_name'
> exec ^
: string
| set @a = 1, file_pattern = 2
< error
| 1193 HY000 Unknown system variable 'file_pattern'
> exec ^
: string
| set @a = 1, file_prefix = 2
< error
| 1193 HY000 Unknown system variable 'file_prefix'
> exec ^
: string
| set @a = 1, files = 2
< error
| 1193 HY000 Unknown system variable 'files'
> exec ^
: string
| set @a = 1, filter = 2
< error
| 1193 HY000 Unknown system variable 'filter'
> exec ^
: string
| set @a = 1, finish = 2
< error
| 1193 HY000 Unknown system variable 'finish'
> exec ^
: string
| set @a = 1, first = 2
< error
| 1193 HY000 Unknown system variable 'first'
> exec ^
: string
| set @a = 1, first_value = 2
< error
| 1064 42000 …near 'first_value = 2' at line 1
> exec ^
: string
| set @a = 1, fixed = 2
< error
| 1193 HY000 Unknown system variable 'fixed'
> exec ^
: string
| set @a = 1, float = 2
< error
| 1064 42000 …near 'float = 2' at line 1
> exec ^
: string
| set @a = 1, float4 = 2
< error
| 1064 42000 …near 'float4 = 2' at line 1
> exec ^
: string
| set @a = 1, float8 = 2
< error
| 1064 42000 …near 'float8 = 2' at line 1
> exec ^
: string
| set @a = 1, flush = 2
< error
| 1229 HY000 Variable 'flush' is a GLOBAL variable and should be set with SET GLOBAL
> exec ^
: string
| set @a = 1, flush = 2
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set @a = 1, following = 2
< error
| 1193 HY000 Unknown system variable 'following'
> exec ^
: string
| set @a = 1, follows = 2
< error
| 1193 HY000 Unknown system variable 'follows'
> exec ^
: string
| set @a = 1, for = 2
< error
| 1064 42000 …near 'for = 2' at line 1
> exec ^
: string
| set @a = 1, force = 2
< error
| 1064 42000 …near 'force = 2' at line 1
> exec ^
: string
| set @a = 1, foreign = 2
< error
| 1064 42000 …near 'foreign = 2' at line 1
> exec ^
: string
| set @a = 1, format = 2
< error
| 1193 HY000 Unknown system variable 'format'
> exec ^
: string
| set @a = 1, found = 2
< error
| 1193 HY000 Unknown system variable 'found'
> exec ^
: string
| set @a = 1, from = 2
< error
| 1064 42000 …near 'from = 2' at line 1
> exec ^
: string
| set @a = 1, full = 2
< error
| 1193 HY000 Unknown system variable 'full'
> exec ^
: string
| set @a = 1, fulltext = 2
< error
| 1064 42000 …near 'fulltext = 2' at line 1
> exec ^
: string
| set @a = 1, function = 2
< error
| 1064 42000 …near 'function = 2' at line 1
> exec ^
: string
| set @a = 1, general = 2
< error
| 1193 HY000 Unknown system variable 'general'
> exec ^
: string
| set @a = 1, generate = 2
< error
| 1193 HY000 Unknown system variable 'generate'
> exec ^
: string
| set @a = 1, generated = 2
< error
| 1064 42000 …near 'generated = 2' at line 1
> exec ^
: string
| set @a = 1, geomcollection = 2
< error
| 1193 HY000 Unknown system variable 'geomcollection'
> exec ^
: string
| set @a = 1, geometry = 2
< error
| 1193 HY000 Unknown system variable 'geometry'
> exec ^
: string
| set @a = 1, geometrycollection = 2
< error
| 1193 HY000 Unknown system variable 'geometrycollection'
> exec ^
: string
| set @a = 1, get = 2
< error
| 1064 42000 …near 'get = 2' at line 1
> exec ^
: string
| set @a = 1, get_format = 2
< error
| 1193 HY000 Unknown system variable 'get_format'
> exec ^
: string
| set @a = 1, get_master_public_key = 2
< error
| 1193 HY000 Unknown system variable 'get_master_public_key'
> exec ^
: string
| set @a = 1, get_source_public_key = 2
< error
| 1193 HY000 Unknown system variable 'get_source_public_key'
> exec ^
: string
| set @a = 1, global = 2
< error
| 1064 42000 …near '= 2' at line 1
> exec ^
: string
| set @a = 1, grant = 2
< error
| 1064 42000 …near 'grant = 2' at line 1
> exec ^
: string
| set @a = 1, grants = 2
< error
| 1193 HY000 Unknown system variable 'grants'
> exec ^
: string
| set @a = 1, greatest = 2
< error
| 1193 HY000 Unknown system variable 'greatest'
> exec ^
: string
| set @a = 1, group = 2
< error
| 1064 42000 …near 'group = 2' at line 1
> exec ^
: string
| set @a = 1, group_concat = 2
< error
| 1193 HY000 Unknown system variable 'group_concat'
> exec ^
: string
| set @a = 1, group_replication = 2
< error
| 1193 HY000 Unknown system variable 'group_replication'
> exec ^
: string
| set @a = 1, grouping = 2
< error
| 1064 42000 …near 'grouping = 2' at line 1
> exec ^
: string
| set @a = 1, groups = 2
< error
| 1064 42000 …near 'groups = 2' at line 1
> exec ^
: string
| set @a = 1, gtid_only = 2
< error
| 1193 HY000 Unknown system variable 'gtid_only'
> exec ^
: string
| set @a = 1, gtids = 2
< error
| 1193 HY000 Unknown system variable 'gtids'
> exec ^
: string
| set @a = 1, guided = 2
< error
| 1193 HY000 Unknown system variable 'guided'
> exec ^
: string
| set @a = 1, handler = 2
< error
| 1193 HY000 Unknown system variable 'handler'
> exec ^
: string
| set @a = 1, hash = 2
< error
| 1193 HY000 Unknown system variable 'hash'
> exec ^
: string
| set @a = 1, having = 2
< error
| 1064 42000 …near 'having = 2' at line 1
> exec ^
: string
| set @a = 1, header = 2
< error
| 1193 HY000 Unknown system variable 'header'
> exec ^
: string
| set @a = 1, help = 2
< error
| 1193 HY000 Unknown system variable 'help'
> exec ^
: string
| set @a = 1, high_priority = 2
< error
| 1064 42000 …near 'high_priority = 2' at line 1
> exec ^
: string
| set @a = 1, histogram = 2
< error
| 1193 HY000 Unknown system variable 'histogram'
> exec ^
: string
| set @a = 1, history = 2
< error
| 1193 HY000 Unknown system variable 'history'
> exec ^
: string
| set @a = 1, host = 2
< error
| 1193 HY000 Unknown system variable 'host'
> exec ^
: string
| set @a = 1, hosts = 2
< error
| 1193 HY000 Unknown system variable 'hosts'
> exec ^
: string
| set @a = 1, hour = 2
< error
| 1193 HY000 Unknown system variable 'hour'
> exec ^
: string
| set @a = 1, hour_microsecond = 2
< error
| 1064 42000 …near 'hour_microsecond = 2' at line 1
> exec ^
: string
| set @a = 1, hour_minute = 2
< error
| 1064 42000 …near 'hour_minute = 2' at line 1
> exec ^
: string
| set @a = 1, hour_second = 2
< error
| 1064 42000 …near 'hour_second = 2' at line 1
> exec ^
: string
| set @a = 1, identified = 2
< error
| 1193 HY000 Unknown system variable 'identified'
> exec ^
: string
| set @a = 1, if = 2
< error
| 1064 42000 …near 'if = 2' at line 1
> exec ^
: string
| set @a = 1, ignore = 2
< error
| 1064 42000 …near 'ignore = 2' at line 1
> exec ^
: string
| set @a = 1, ignore_server_ids = 2
< error
| 1193 HY000 Unknown system variable 'ignore_server_ids'
> exec ^
: string
| set @a = 1, import = 2
< error
| 1193 HY000 Unknown system variable 'import'
> exec ^
: string
| set @a = 1, in = 2
< error
| 1064 42000 …near 'in = 2' at line 1
> exec ^
: string
| set @a = 1, inactive = 2
< error
| 1193 HY000 Unknown system variable 'inactive'
> exec ^
: string
| set @a = 1, index = 2
< error
| 1064 42000 …near 'index = 2' at line 1
> exec ^
: string
| set @a = 1, indexes = 2
< error
| 1193 HY000 Unknown system variable 'indexes'
> exec ^
: string
| set @a = 1, infile = 2
< error
| 1064 42000 …near 'infile = 2' at line 1
> exec ^
: string
| set @a = 1, initial = 2
< error
| 1193 HY000 Unknown system variable 'initial'
> exec ^
: string
| set @a = 1, initial_size = 2
< error
| 1193 HY000 Unknown system variable 'initial_size'
> exec ^
: string
| set @a = 1, initiate = 2
< error
| 1193 HY000 Unknown system variable 'initiate'
> exec ^
: string
| set @a = 1, inner = 2
< error
| 1064 42000 …near 'inner = 2' at line 1
> exec ^
: string
| set @a = 1, inout = 2
< error
| 1064 42000 …near 'inout = 2' at line 1
> exec ^
: string
| set @a = 1, insensitive = 2
< error
| 1064 42000 …near 'insensitive = 2' at line 1
> exec ^
: string
| set @a = 1, insert = 2
< error
| 1064 42000 …near 'insert = 2' at line 1
> exec ^
: string
| set @a = 1, insert_method = 2
< error
| 1193 HY000 Unknown system variable 'insert_method'
> exec ^
: string
| set @a = 1, install = 2
< error
| 1193 HY000 Unknown system variable 'install'
> exec ^
: string
| set @a = 1, instance = 2
< error
| 1193 HY000 Unknown system variable 'instance'
> exec ^
: string
| set @a = 1, int = 2
< error
| 1064 42000 …near 'int = 2' at line 1
> exec ^
: string
| set @a = 1, int1 = 2
< error
| 1064 42000 …near 'int1 = 2' at line 1
> exec ^
: string
| set @a = 1, int2 = 2
< error
| 1064 42000 …near 'int2 = 2' at line 1
> exec ^
: string
| set @a = 1, int3 = 2
< error
| 1064 42000 …near 'int3 = 2' at line 1
> exec ^
: string
| set @a = 1, int4 = 2
< error
| 1064 42000 …near 'int4 = 2' at line 1
> exec ^
: string
| set @a = 1, int8 = 2
< error
| 1064 42000 …near 'int8 = 2' at line 1
> exec ^
: string
| set @a = 1, integer = 2
< error
| 1064 42000 …near 'integer = 2' at line 1
> exec ^
: string
| set @a = 1, intersect = 2
< error
| 1064 42000 …near 'intersect = 2' at line 1
> exec ^
: string
| set @a = 1, interval = 2
< error
| 1064 42000 …near 'interval = 2' at line 1
> exec ^
: string
| set @a = 1, into = 2
< error
| 1064 42000 …near 'into = 2' at line 1
> exec ^
: string
| set @a = 1, invisible = 2
< error
| 1193 HY000 Unknown system variable 'invisible'
> exec ^
: string
| set @a = 1, invoker = 2
< error
| 1193 HY000 Unknown system variable 'invoker'
> exec ^
: string
| set @a = 1, io = 2
< error
| 1193 HY000 Unknown system variable 'io'
> exec ^
: string
| set @a = 1, io_after_gtids = 2
< error
| 1064 42000 …near 'io_after_gtids = 2' at line 1
> exec ^
: string
| set @a = 1, io_before_gtids = 2
< error
| 1064 42000 …near 'io_before_gtids = 2' at line 1
> exec ^
: string
| set @a = 1, io_thread = 2
< error
| 1193 HY000 Unknown system variable 'io_thread'
> exec ^
: string
| set @a = 1, ipc = 2
< error
| 1193 HY000 Unknown system variable 'ipc'
> exec ^
: string
| set @a = 1, is = 2
< error
| 1064 42000 …near 'is = 2' at line 1
> exec ^
: string
| set @a = 1, isolation = 2
< error
| 1193 HY000 Unknown system variable 'isolation'
> exec ^
: string
| set @a = 1, issuer = 2
< error
| 1193 HY000 Unknown system variable 'issuer'
> exec ^
: string
| set @a = 1, iterate = 2
< error
| 1064 42000 …near 'iterate = 2' at line 1
> exec ^
: string
| set @a = 1, join = 2
< error
| 1064 42000 …near 'join = 2' at line 1
> exec ^
: string
| set @a = 1, json = 2
< error
| 1193 HY000 Unknown system variable 'json'
> exec ^
: string
| set @a = 1, json_arrayagg = 2
< error
| 1193 HY000 Unknown system variable 'json_arrayagg'
> exec ^
: string
| set @a = 1, json_duality_object = 2
< error
| 1193 HY000 Unknown system variable 'json_duality_object'
> exec ^
: string
| set @a = 1, json_objectagg = 2
< error
| 1193 HY000 Unknown system variable 'json_objectagg'
> exec ^
: string
| set @a = 1, json_table = 2
< error
| 1064 42000 …near 'json_table = 2' at line 1
> exec ^
: string
| set @a = 1, json_value = 2
< error
| 1193 HY000 Unknown system variable 'json_value'
> exec ^
: string
| set @a = 1, key = 2
< error
| 1064 42000 …near 'key = 2' at line 1
> exec ^
: string
| set @a = 1, key_block_size = 2
< error
| 1193 HY000 Unknown system variable 'key_block_size'
> exec ^
: string
| set @a = 1, keyring = 2
< error
| 1193 HY000 Unknown system variable 'keyring'
> exec ^
: string
| set @a = 1, keys = 2
< error
| 1064 42000 …near 'keys = 2' at line 1
> exec ^
: string
| set @a = 1, kill = 2
< error
| 1064 42000 …near 'kill = 2' at line 1
> exec ^
: string
| set @a = 1, lag = 2
< error
| 1064 42000 …near 'lag = 2' at line 1
> exec ^
: string
| set @a = 1, language = 2
< error
| 1193 HY000 Unknown system variable 'language'
> exec ^
: string
| set @a = 1, last = 2
< error
| 1193 HY000 Unknown system variable 'last'
> exec ^
: string
| set @a = 1, last_value = 2
< error
| 1064 42000 …near 'last_value = 2' at line 1
> exec ^
: string
| set @a = 1, lateral = 2
< error
| 1064 42000 …near 'lateral = 2' at line 1
> exec ^
: string
| set @a = 1, lead = 2
< error
| 1064 42000 …near 'lead = 2' at line 1
> exec ^
: string
| set @a = 1, leading = 2
< error
| 1064 42000 …near 'leading = 2' at line 1
> exec ^
: string
| set @a = 1, least = 2
< error
| 1193 HY000 Unknown system variable 'least'
> exec ^
: string
| set @a = 1, leave = 2
< error
| 1064 42000 …near 'leave = 2' at line 1
> exec ^
: string
| set @a = 1, leaves = 2
< error
| 1193 HY000 Unknown system variable 'leaves'
> exec ^
: string
| set @a = 1, left = 2
< error
| 1064 42000 …near 'left = 2' at line 1
> exec ^
: string
| set @a = 1, less = 2
< error
| 1193 HY000 Unknown system variable 'less'
> exec ^
: string
| set @a = 1, level = 2
< error
| 1193 HY000 Unknown system variable 'level'
> exec ^
: string
| set @a = 1, library = 2
< error
| 1193 HY000 Unknown system variable 'library'
> exec ^
: string
| set @a = 1, like = 2
< error
| 1064 42000 …near 'like = 2' at line 1
> exec ^
: string
| set @a = 1, limit = 2
< error
| 1064 42000 …near 'limit = 2' at line 1
> exec ^
: string
| set @a = 1, linear = 2
< error
| 1064 42000 …near 'linear = 2' at line 1
> exec ^
: string
| set @a = 1, lines = 2
< error
| 1064 42000 …near 'lines = 2' at line 1
> exec ^
: string
| set @a = 1, linestring = 2
< error
| 1193 HY000 Unknown system variable 'linestring'
> exec ^
: string
| set @a = 1, list = 2
< error
| 1193 HY000 Unknown system variable 'list'
> exec ^
: string
| set @a = 1, load = 2
< error
| 1064 42000 …near 'load = 2' at line 1
> exec ^
: string
| set @a = 1, local = 2
< error
| 1064 42000 …near '= 2' at line 1
> exec ^
: string
| set @a = 1, localtime = 2
< error
| 1064 42000 …near 'localtime = 2' at line 1
> exec ^
: string
| set @a = 1, localtimestamp = 2
< error
| 1064 42000 …near 'localtimestamp = 2' at line 1
> exec ^
: string
| set @a = 1, lock = 2
< error
| 1064 42000 …near 'lock = 2' at line 1
> exec ^
: string
| set @a = 1, locked = 2
< error
| 1193 HY000 Unknown system variable 'locked'
> exec ^
: string
| set @a = 1, locks = 2
< error
| 1193 HY000 Unknown system variable 'locks'
> exec ^
: string
| set @a = 1, log = 2
< error
| 1193 HY000 Unknown system variable 'log'
> exec ^
: string
| set @a = 1, logfile = 2
< error
| 1193 HY000 Unknown system variable 'logfile'
> exec ^
: string
| set @a = 1, logs = 2
< error
| 1193 HY000 Unknown system variable 'logs'
> exec ^
: string
| set @a = 1, long = 2
< error
| 1064 42000 …near 'long = 2' at line 1
> exec ^
: string
| set @a = 1, longblob = 2
< error
| 1064 42000 …near 'longblob = 2' at line 1
> exec ^
: string
| set @a = 1, longtext = 2
< error
| 1064 42000 …near 'longtext = 2' at line 1
> exec ^
: string
| set @a = 1, loop = 2
< error
| 1064 42000 …near 'loop = 2' at line 1
> exec ^
: string
| set @a = 1, low_priority = 2
< error
| 1064 42000 …near 'low_priority = 2' at line 1
> exec ^
: string
| set @a = 1, manifest = 2
< error
| 1193 HY000 Unknown system variable 'manifest'
> exec ^
: string
| set @a = 1, manual = 2
< error
| 1193 HY000 Unknown system variable 'manual'
> exec ^
: string
| set @a = 1, master = 2
< error
| 1193 HY000 Unknown system variable 'master'
> exec ^
: string
| set @a = 1, master_auto_position = 2
< error
| 1193 HY000 Unknown system variable 'master_auto_position'
> exec ^
: string
| set @a = 1, master_bind = 2
< error
| 1064 42000 …near 'master_bind = 2' at line 1
> exec ^
: string
| set @a = 1, master_compression_algorithms = 2
< error
| 1193 HY000 Unknown system variable 'master_compression_algorithms'
> exec ^
: string
| set @a = 1, master_connect_retry = 2
< error
| 1193 HY000 Unknown system variable 'master_connect_retry'
> exec ^
: string
| set @a = 1, master_delay = 2
< error
| 1193 HY000 Unknown system variable 'master_delay'
> exec ^
: string
| set @a = 1, master_heartbeat_period = 2
< error
| 1193 HY000 Unknown system variable 'master_heartbeat_period'
> exec ^
: string
| set @a = 1, master_host = 2
< error
| 1193 HY000 Unknown system variable 'master_host'
> exec ^
: string
| set @a = 1, master_log_file = 2
< error
| 1193 HY000 Unknown system variable 'master_log_file'
> exec ^
: string
| set @a = 1, master_log_pos = 2
< error
| 1193 HY000 Unknown system variable 'master_log_pos'
> exec ^
: string
| set @a = 1, master_password = 2
< error
| 1193 HY000 Unknown system variable 'master_password'
> exec ^
: string
| set @a = 1, master_port = 2
< error
| 1193 HY000 Unknown system variable 'master_port'
> exec ^
: string
| set @a = 1, master_public_key_path = 2
< error
| 1193 HY000 Unknown system variable 'master_public_key_path'
> exec ^
: string
| set @a = 1, master_retry_count = 2
< error
| 1193 HY000 Unknown system variable 'master_retry_count'
> exec ^
: string
| set @a = 1, master_ssl = 2
< error
| 1193 HY000 Unknown system variable 'master_ssl'
> exec ^
: string
| set @a = 1, master_ssl_ca = 2
< error
| 1193 HY000 Unknown system variable 'master_ssl_ca'
> exec ^
: string
| set @a = 1, master_ssl_capath = 2
< error
| 1193 HY000 Unknown system variable 'master_ssl_capath'
> exec ^
: string
| set @a = 1, master_ssl_cert = 2
< error
| 1193 HY000 Unknown system variable 'master_ssl_cert'
> exec ^
: string
| set @a = 1, master_ssl_cipher = 2
< error
| 1193 HY000 Unknown system variable 'master_ssl_cipher'
> exec ^
: string
| set @a = 1, master_ssl_crl = 2
< error
| 1193 HY000 Unknown system variable 'master_ssl_crl'
> exec ^
: string
| set @a = 1, master_ssl_crlpath = 2
< error
| 1193 HY000 Unknown system variable 'master_ssl_crlpath'
> exec ^
: string
| set @a = 1, master_ssl_key = 2
< error
| 1193 HY000 Unknown system variable 'master_ssl_key'
> exec ^
: string
| set @a = 1, master_ssl_verify_server_cert = 2
< error
| 1064 42000 …near 'master_ssl_verify_server_cert = 2' at line 1
> exec ^
: string
| set @a = 1, master_tls_ciphersuites = 2
< error
| 1193 HY000 Unknown system variable 'master_tls_ciphersuites'
> exec ^
: string
| set @a = 1, master_tls_version = 2
< error
| 1193 HY000 Unknown system variable 'master_tls_version'
> exec ^
: string
| set @a = 1, master_user = 2
< error
| 1193 HY000 Unknown system variable 'master_user'
> exec ^
: string
| set @a = 1, master_zstd_compression_level = 2
< error
| 1193 HY000 Unknown system variable 'master_zstd_compression_level'
> exec ^
: string
| set @a = 1, match = 2
< error
| 1064 42000 …near 'match = 2' at line 1
> exec ^
: string
| set @a = 1, materialized = 2
< error
| 1193 HY000 Unknown system variable 'materialized'
> exec ^
: string
| set @a = 1, max = 2
< error
| 1193 HY000 Unknown system variable 'max'
> exec ^
: string
| set @a = 1, max_connections_per_hour = 2
< error
| 1193 HY000 Unknown system variable 'max_connections_per_hour'
> exec ^
: string
| set @a = 1, max_queries_per_hour = 2
< error
| 1193 HY000 Unknown system variable 'max_queries_per_hour'
> exec ^
: string
| set @a = 1, max_rows = 2
< error
| 1193 HY000 Unknown system variable 'max_rows'
> exec ^
: string
| set @a = 1, max_size = 2
< error
| 1193 HY000 Unknown system variable 'max_size'
> exec ^
: string
| set @a = 1, max_updates_per_hour = 2
< error
| 1193 HY000 Unknown system variable 'max_updates_per_hour'
> exec ^
: string
| set @a = 1, max_user_connections = 2
< ok 0 0
> exec ^
: string
| set @a = 1, maxvalue = 2
< error
| 1064 42000 …near 'maxvalue = 2' at line 1
> exec ^
: string
| set @a = 1, medium = 2
< error
| 1193 HY000 Unknown system variable 'medium'
> exec ^
: string
| set @a = 1, mediumblob = 2
< error
| 1064 42000 …near 'mediumblob = 2' at line 1
> exec ^
: string
| set @a = 1, mediumint = 2
< error
| 1064 42000 …near 'mediumint = 2' at line 1
> exec ^
: string
| set @a = 1, mediumtext = 2
< error
| 1064 42000 …near 'mediumtext = 2' at line 1
> exec ^
: string
| set @a = 1, member = 2
< error
| 1193 HY000 Unknown system variable 'member'
> exec ^
: string
| set @a = 1, memory = 2
< error
| 1193 HY000 Unknown system variable 'memory'
> exec ^
: string
| set @a = 1, merge = 2
< error
| 1193 HY000 Unknown system variable 'merge'
> exec ^
: string
| set @a = 1, message_text = 2
< error
| 1193 HY000 Unknown system variable 'message_text'
> exec ^
: string
| set @a = 1, microsecond = 2
< error
| 1193 HY000 Unknown system variable 'microsecond'
> exec ^
: string
| set @a = 1, mid = 2
< error
| 1193 HY000 Unknown system variable 'mid'
> exec ^
: string
| set @a = 1, middleint = 2
< error
| 1064 42000 …near 'middleint = 2' at line 1
> exec ^
: string
| set @a = 1, migrate = 2
< error
| 1193 HY000 Unknown system variable 'migrate'
> exec ^
: string
| set @a = 1, min = 2
< error
| 1193 HY000 Unknown system variable 'min'
> exec ^
: string
| set @a = 1, min_rows = 2
< error
| 1193 HY000 Unknown system variable 'min_rows'
> exec ^
: string
| set @a = 1, minute = 2
< error
| 1193 HY000 Unknown system variable 'minute'
> exec ^
: string
| set @a = 1, minute_microsecond = 2
< error
| 1064 42000 …near 'minute_microsecond = 2' at line 1
> exec ^
: string
| set @a = 1, minute_second = 2
< error
| 1064 42000 …near 'minute_second = 2' at line 1
> exec ^
: string
| set @a = 1, mod = 2
< error
| 1064 42000 …near 'mod = 2' at line 1
> exec ^
: string
| set @a = 1, mode = 2
< error
| 1193 HY000 Unknown system variable 'mode'
> exec ^
: string
| set @a = 1, modifies = 2
< error
| 1064 42000 …near 'modifies = 2' at line 1
> exec ^
: string
| set @a = 1, modify = 2
< error
| 1193 HY000 Unknown system variable 'modify'
> exec ^
: string
| set @a = 1, month = 2
< error
| 1193 HY000 Unknown system variable 'month'
> exec ^
: string
| set @a = 1, multilinestring = 2
< error
| 1193 HY000 Unknown system variable 'multilinestring'
> exec ^
: string
| set @a = 1, multipoint = 2
< error
| 1193 HY000 Unknown system variable 'multipoint'
> exec ^
: string
| set @a = 1, multipolygon = 2
< error
| 1193 HY000 Unknown system variable 'multipolygon'
> exec ^
: string
| set @a = 1, mutex = 2
< error
| 1193 HY000 Unknown system variable 'mutex'
> exec ^
: string
| set @a = 1, mysql_errno = 2
< error
| 1193 HY000 Unknown system variable 'mysql_errno'
> exec ^
: string
| set @a = 1, name = 2
< error
| 1193 HY000 Unknown system variable 'name'
> exec ^
: string
| set @a = 1, names = 2
< error
| 1064 42000 …near '= 2' at line 1
> exec ^
: string
| set @a = 1, national = 2
< error
| 1193 HY000 Unknown system variable 'national'
> exec ^
: string
| set @a = 1, natural = 2
< error
| 1064 42000 …near 'natural = 2' at line 1
> exec ^
: string
| set @a = 1, nchar = 2
< error
| 1193 HY000 Unknown system variable 'nchar'
> exec ^
: string
| set @a = 1, ndb = 2
< error
| 1193 HY000 Unknown system variable 'ndb'
> exec ^
: string
| set @a = 1, ndbcluster = 2
< error
| 1193 HY000 Unknown system variable 'ndbcluster'
> exec ^
: string
| set @a = 1, nested = 2
< error
| 1193 HY000 Unknown system variable 'nested'
> exec ^
: string
| set @a = 1, network_namespace = 2
< error
| 1193 HY000 Unknown system variable 'network_namespace'
> exec ^
: string
| set @a = 1, never = 2
< error
| 1193 HY000 Unknown system variable 'never'
> exec ^
: string
| set @a = 1, new = 2
< ok 0 0
> exec ^
: string
| set @a = 1, next = 2
< error
| 1193 HY000 Unknown system variable 'next'
> exec ^
: string
| set @a = 1, no = 2
< error
| 1193 HY000 Unknown system variable 'no'
> exec ^
: string
| set @a = 1, no_wait = 2
< error
| 1193 HY000 Unknown system variable 'no_wait'
> exec ^
: string
| set @a = 1, no_write_to_binlog = 2
< error
| 1064 42000 …near 'no_write_to_binlog = 2' at line 1
> exec ^
: string
| set @a = 1, nodegroup = 2
< error
| 1193 HY000 Unknown system variable 'nodegroup'
> exec ^
: string
| set @a = 1, none = 2
< error
| 1193 HY000 Unknown system variable 'none'
> exec ^
: string
| set @a = 1, not = 2
< error
| 1064 42000 …near 'not = 2' at line 1
> exec ^
: string
| set @a = 1, now = 2
< error
| 1193 HY000 Unknown system variable 'now'
> exec ^
: string
| set @a = 1, nowait = 2
< error
| 1193 HY000 Unknown system variable 'nowait'
> exec ^
: string
| set @a = 1, nth_value = 2
< error
| 1064 42000 …near 'nth_value = 2' at line 1
> exec ^
: string
| set @a = 1, ntile = 2
< error
| 1064 42000 …near 'ntile = 2' at line 1
> exec ^
: string
| set @a = 1, null = 2
< error
| 1064 42000 …near 'null = 2' at line 1
> exec ^
: string
| set @a = 1, nullif = 2
< error
| 1193 HY000 Unknown system variable 'nullif'
> exec ^
: string
| set @a = 1, nulls = 2
< error
| 1193 HY000 Unknown system variable 'nulls'
> exec ^
: string
| set @a = 1, number = 2
< error
| 1193 HY000 Unknown system variable 'number'
> exec ^
: string
| set @a = 1, numeric = 2
< error
| 1064 42000 …near 'numeric = 2' at line 1
> exec ^
: string
| set @a = 1, nvarchar = 2
< error
| 1193 HY000 Unknown system variable 'nvarchar'
> exec ^
: string
| set @a = 1, of = 2
< error
| 1064 42000 …near 'of = 2' at line 1
> exec ^
: string
| set @a = 1, off = 2
< error
| 1193 HY000 Unknown system variable 'off'
> exec ^
: string
| set @a = 1, offset = 2
< error
| 1193 HY000 Unknown system variable 'offset'
> exec ^
: string
| set @a = 1, oj = 2
< error
| 1193 HY000 Unknown system variable 'oj'
> exec ^
: string
| set @a = 1, old = 2
< error
| 1229 HY000 Variable 'old' is a GLOBAL variable and should be set with SET GLOBAL
> exec ^
: string
| set @a = 1, old = 2
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set @a = 1, on = 2
< error
| 1064 42000 …near 'on = 2' at line 1
> exec ^
: string
| set @a = 1, one = 2
< error
| 1193 HY000 Unknown system variable 'one'
> exec ^
: string
| set @a = 1, only = 2
< error
| 1193 HY000 Unknown system variable 'only'
> exec ^
: string
| set @a = 1, open = 2
< error
| 1193 HY000 Unknown system variable 'open'
> exec ^
: string
| set @a = 1, optimize = 2
< error
| 1064 42000 …near 'optimize = 2' at line 1
> exec ^
: string
| set @a = 1, optimizer_costs = 2
< error
| 1064 42000 …near 'optimizer_costs = 2' at line 1
> exec ^
: string
| set @a = 1, option = 2
< error
| 1064 42000 …near 'option = 2' at line 1
> exec ^
: string
| set @a = 1, optional = 2
< error
| 1193 HY000 Unknown system variable 'optional'
> exec ^
: string
| set @a = 1, optionally = 2
< error
| 1064 42000 …near 'optionally = 2' at line 1
> exec ^
: string
| set @a = 1, options = 2
< error
| 1193 HY000 Unknown system variable 'options'
> exec ^
: string
| set @a = 1, or = 2
< error
| 1064 42000 …near 'or = 2' at line 1
> exec ^
: string
| set @a = 1, order = 2
< error
| 1064 42000 …near 'order = 2' at line 1
> exec ^
: string
| set @a = 1, ordinality = 2
< error
| 1193 HY000 Unknown system variable 'ordinality'
> exec ^
: string
| set @a = 1, organization = 2
< error
| 1193 HY000 Unknown system variable 'organization'
> exec ^
: string
| set @a = 1, others = 2
< error
| 1193 HY000 Unknown system variable 'others'
> exec ^
: string
| set @a = 1, out = 2
< error
| 1064 42000 …near 'out = 2' at line 1
> exec ^
: string
| set @a = 1, outer = 2
< error
| 1064 42000 …near 'outer = 2' at line 1
> exec ^
: string
| set @a = 1, outfile = 2
< error
| 1064 42000 …near 'outfile = 2' at line 1
> exec ^
: string
| set @a = 1, over = 2
< error
| 1064 42000 …near 'over = 2' at line 1
> exec ^
: string
| set @a = 1, overlay = 2
< error
| 1193 HY000 Unknown system variable 'overlay'
> exec ^
: string
| set @a = 1, overwrite = 2
< error
| 1193 HY000 Unknown system variable 'overwrite'
> exec ^
: string
| set @a = 1, owner = 2
< error
| 1193 HY000 Unknown system variable 'owner'
> exec ^
: string
| set @a = 1, pack_keys = 2
< error
| 1193 HY000 Unknown system variable 'pack_keys'
> exec ^
: string
| set @a = 1, page = 2
< error
| 1193 HY000 Unknown system variable 'page'
> exec ^
: string
| set @a = 1, parallel = 2
< error
| 1193 HY000 Unknown system variable 'parallel'
> exec ^
: string
| set @a = 1, parameters = 2
< error
| 1193 HY000 Unknown system variable 'parameters'
> exec ^
: string
| set @a = 1, parse_tree = 2
< error
| 1193 HY000 Unknown system variable 'parse_tree'
> exec ^
: string
| set @a = 1, parser = 2
< error
| 1193 HY000 Unknown system variable 'parser'
> exec ^
: string
| set @a = 1, partial = 2
< error
| 1193 HY000 Unknown system variable 'partial'
> exec ^
: string
| set @a = 1, partition = 2
< error
| 1064 42000 …near 'partition = 2' at line 1
> exec ^
: string
| set @a = 1, partitioning = 2
< error
| 1193 HY000 Unknown system variable 'partitioning'
> exec ^
: string
| set @a = 1, partitions = 2
< error
| 1193 HY000 Unknown system variable 'partitions'
> exec ^
: string
| set @a = 1, password = 2
< error
| 1193 HY000 Unknown system variable 'password'
> exec ^
: string
| set @a = 1, password_lock_time = 2
< error
| 1193 HY000 Unknown system variable 'password_lock_time'
> exec ^
: string
| set @a = 1, path = 2
< error
| 1193 HY000 Unknown system variable 'path'
> exec ^
: string
| set @a = 1, percent_rank = 2
< error
| 1064 42000 …near 'percent_rank = 2' at line 1
> exec ^
: string
| set @a = 1, persist = 2
< error
| 1064 42000 …near '= 2' at line 1
> exec ^
: string
| set @a = 1, persist_only = 2
< error
| 1064 42000 …near '= 2' at line 1
> exec ^
: string
| set @a = 1, phase = 2
< error
| 1193 HY000 Unknown system variable 'phase'
> exec ^
: string
| set @a = 1, placing = 2
< error
| 1193 HY000 Unknown system variable 'placing'
> exec ^
: string
| set @a = 1, plugin = 2
< error
| 1193 HY000 Unknown system variable 'plugin'
> exec ^
: string
| set @a = 1, plugin_dir = 2
< error
| 1229 HY000 Variable 'plugin_dir' is a GLOBAL variable and should be set with SET GLOBAL
> exec ^
: string
| set @a = 1, plugin_dir = 2
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set @a = 1, plugins = 2
< error
| 1193 HY000 Unknown system variable 'plugins'
> exec ^
: string
| set @a = 1, point = 2
< error
| 1193 HY000 Unknown system variable 'point'
> exec ^
: string
| set @a = 1, polygon = 2
< error
| 1193 HY000 Unknown system variable 'polygon'
> exec ^
: string
| set @a = 1, port = 2
< error
| 1229 HY000 Variable 'port' is a GLOBAL variable and should be set with SET GLOBAL
> exec ^
: string
| set @a = 1, port = 2
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set @a = 1, position = 2
< error
| 1193 HY000 Unknown system variable 'position'
> exec ^
: string
| set @a = 1, precedes = 2
< error
| 1193 HY000 Unknown system variable 'precedes'
> exec ^
: string
| set @a = 1, preceding = 2
< error
| 1193 HY000 Unknown system variable 'preceding'
> exec ^
: string
| set @a = 1, precision = 2
< error
| 1064 42000 …near 'precision = 2' at line 1
> exec ^
: string
| set @a = 1, prefix = 2
< error
| 1193 HY000 Unknown system variable 'prefix'
> exec ^
: string
| set @a = 1, prepare = 2
< error
| 1193 HY000 Unknown system variable 'prepare'
> exec ^
: string
| set @a = 1, preserve = 2
< error
| 1193 HY000 Unknown system variable 'preserve'
> exec ^
: string
| set @a = 1, prev = 2
< error
| 1193 HY000 Unknown system variable 'prev'
> exec ^
: string
| set @a = 1, primary = 2
< error
| 1064 42000 …near 'primary = 2' at line 1
> exec ^
: string
| set @a = 1, privilege_checks_user = 2
< error
| 1193 HY000 Unknown system variable 'privilege_checks_user'
> exec ^
: string
| set @a = 1, privileges = 2
< error
| 1193 HY000 Unknown system variable 'privileges'
> exec ^
: string
| set @a = 1, procedure = 2
< error
| 1064 42000 …near 'procedure = 2' at line 1
> exec ^
: string
| set @a = 1, process = 2
< error
| 1193 HY000 Unknown system variable 'process'
> exec ^
: string
| set @a = 1, processlist = 2
< error
| 1193 HY000 Unknown system variable 'processlist'
> exec ^
: string
| set @a = 1, profile = 2
< error
| 1193 HY000 Unknown system variable 'profile'
> exec ^
: string
| set @a = 1, profiles = 2
< error
| 1193 HY000 Unknown system variable 'profiles'
> exec ^
: string
| set @a = 1, proxy = 2
< error
| 1193 HY000 Unknown system variable 'proxy'
> exec ^
: string
| set @a = 1, purge = 2
< error
| 1064 42000 …near 'purge = 2' at line 1
> exec ^
: string
| set @a = 1, qualify = 2
< error
| 1193 HY000 Unknown system variable 'qualify'
> exec ^
: string
| set @a = 1, quarter = 2
< error
| 1193 HY000 Unknown system variable 'quarter'
> exec ^
: string
| set @a = 1, query = 2
< error
| 1193 HY000 Unknown system variable 'query'
> exec ^
: string
| set @a = 1, quick = 2
< error
| 1193 HY000 Unknown system variable 'quick'
> exec ^
: string
| set @a = 1, random = 2
< error
| 1193 HY000 Unknown system variable 'random'
> exec ^
: string
| set @a = 1, range = 2
< error
| 1064 42000 …near 'range = 2' at line 1
> exec ^
: string
| set @a = 1, rank = 2
< error
| 1064 42000 …near 'rank = 2' at line 1
> exec ^
: string
| set @a = 1, read = 2
< error
| 1064 42000 …near 'read = 2' at line 1
> exec ^
: string
| set @a = 1, read_only = 2
< error
| 1229 HY000 Variable 'read_only' is a GLOBAL variable and should be set with SET GLOBAL
> exec ^
: string
| set @a = 1, read_only = 2
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set @a = 1, read_write = 2
< error
| 1064 42000 …near 'read_write = 2' at line 1
> exec ^
: string
| set @a = 1, reads = 2
< error
| 1064 42000 …near 'reads = 2' at line 1
> exec ^
: string
| set @a = 1, real = 2
< error
| 1064 42000 …near 'real = 2' at line 1
> exec ^
: string
| set @a = 1, rebuild = 2
< error
| 1193 HY000 Unknown system variable 'rebuild'
> exec ^
: string
| set @a = 1, recover = 2
< error
| 1193 HY000 Unknown system variable 'recover'
> exec ^
: string
| set @a = 1, recursive = 2
< error
| 1064 42000 …near 'recursive = 2' at line 1
> exec ^
: string
| set @a = 1, redo_buffer_size = 2
< error
| 1193 HY000 Unknown system variable 'redo_buffer_size'
> exec ^
: string
| set @a = 1, redundant = 2
< error
| 1193 HY000 Unknown system variable 'redundant'
> exec ^
: string
| set @a = 1, reference = 2
< error
| 1193 HY000 Unknown system variable 'reference'
> exec ^
: string
| set @a = 1, references = 2
< error
| 1064 42000 …near 'references = 2' at line 1
> exec ^
: string
| set @a = 1, regexp = 2
< error
| 1064 42000 …near 'regexp = 2' at line 1
> exec ^
: string
| set @a = 1, registration = 2
< error
| 1193 HY000 Unknown system variable 'registration'
> exec ^
: string
| set @a = 1, relational = 2
< error
| 1193 HY000 Unknown system variable 'relational'
> exec ^
: string
| set @a = 1, relay = 2
< error
| 1193 HY000 Unknown system variable 'relay'
> exec ^
: string
| set @a = 1, relay_log_file = 2
< error
| 1193 HY000 Unknown system variable 'relay_log_file'
> exec ^
: string
| set @a = 1, relay_log_pos = 2
< error
| 1193 HY000 Unknown system variable 'relay_log_pos'
> exec ^
: string
| set @a = 1, relay_thread = 2
< error
| 1193 HY000 Unknown system variable 'relay_thread'
> exec ^
: string
| set @a = 1, relaylog = 2
< error
| 1193 HY000 Unknown system variable 'relaylog'
> exec ^
: string
| set @a = 1, release = 2
< error
| 1064 42000 …near 'release = 2' at line 1
> exec ^
: string
| set @a = 1, reload = 2
< error
| 1193 HY000 Unknown system variable 'reload'
> exec ^
: string
| set @a = 1, remove = 2
< error
| 1193 HY000 Unknown system variable 'remove'
> exec ^
: string
| set @a = 1, rename = 2
< error
| 1064 42000 …near 'rename = 2' at line 1
> exec ^
: string
| set @a = 1, reorganize = 2
< error
| 1193 HY000 Unknown system variable 'reorganize'
> exec ^
: string
| set @a = 1, repair = 2
< error
| 1193 HY000 Unknown system variable 'repair'
> exec ^
: string
| set @a = 1, repeat = 2
< error
| 1064 42000 …near 'repeat = 2' at line 1
> exec ^
: string
| set @a = 1, repeatable = 2
< error
| 1193 HY000 Unknown system variable 'repeatable'
> exec ^
: string
| set @a = 1, replace = 2
< error
| 1064 42000 …near 'replace = 2' at line 1
> exec ^
: string
| set @a = 1, replica = 2
< error
| 1193 HY000 Unknown system variable 'replica'
> exec ^
: string
| set @a = 1, replicas = 2
< error
| 1193 HY000 Unknown system variable 'replicas'
> exec ^
: string
| set @a = 1, replicate_do_db = 2
< error
| 1193 HY000 Unknown system variable 'replicate_do_db'
> exec ^
: string
| set @a = 1, replicate_do_table = 2
< error
| 1193 HY000 Unknown system variable 'replicate_do_table'
> exec ^
: string
| set @a = 1, replicate_ignore_db = 2
< error
| 1193 HY000 Unknown system variable 'replicate_ignore_db'
> exec ^
: string
| set @a = 1, replicate_ignore_table = 2
< error
| 1193 HY000 Unknown system variable 'replicate_ignore_table'
> exec ^
: string
| set @a = 1, replicate_rewrite_db = 2
< error
| 1193 HY000 Unknown system variable 'replicate_rewrite_db'
> exec ^
: string
| set @a = 1, replicate_wild_do_table = 2
< error
| 1193 HY000 Unknown system variable 'replicate_wild_do_table'
> exec ^
: string
| set @a = 1, replicate_wild_ignore_table = 2
< error
| 1193 HY000 Unknown system variable 'replicate_wild_ignore_table'
> exec ^
: string
| set @a = 1, replication = 2
< error
| 1193 HY000 Unknown system variable 'replication'
> exec ^
: string
| set @a = 1, require = 2
< error
| 1064 42000 …near 'require = 2' at line 1
> exec ^
: string
| set @a = 1, require_row_format = 2
< ok 0 0
> exec ^
: string
| set @a = 1, require_table_primary_key_check = 2
< error
| 1193 HY000 Unknown system variable 'require_table_primary_key_check'
> exec ^
: string
| set @a = 1, reset = 2
< error
| 1193 HY000 Unknown system variable 'reset'
> exec ^
: string
| set @a = 1, resignal = 2
< error
| 1064 42000 …near 'resignal = 2' at line 1
> exec ^
: string
| set @a = 1, resource = 2
< error
| 1193 HY000 Unknown system variable 'resource'
> exec ^
: string
| set @a = 1, respect = 2
< error
| 1193 HY000 Unknown system variable 'respect'
> exec ^
: string
| set @a = 1, restart = 2
< error
| 1193 HY000 Unknown system variable 'restart'
> exec ^
: string
| set @a = 1, restore = 2
< error
| 1193 HY000 Unknown system variable 'restore'
> exec ^
: string
| set @a = 1, restrict = 2
< error
| 1064 42000 …near 'restrict = 2' at line 1
> exec ^
: string
| set @a = 1, resume = 2
< error
| 1193 HY000 Unknown system variable 'resume'
> exec ^
: string
| set @a = 1, retain = 2
< error
| 1193 HY000 Unknown system variable 'retain'
> exec ^
: string
| set @a = 1, return = 2
< error
| 1064 42000 …near 'return = 2' at line 1
> exec ^
: string
| set @a = 1, returned_sqlstate = 2
< error
| 1193 HY000 Unknown system variable 'returned_sqlstate'
> exec ^
: string
| set @a = 1, returning = 2
< error
| 1193 HY000 Unknown system variable 'returning'
> exec ^
: string
| set @a = 1, returns = 2
< error
| 1193 HY000 Unknown system variable 'returns'
> exec ^
: string
| set @a = 1, reuse = 2
< error
| 1193 HY000 Unknown system variable 'reuse'
> exec ^
: string
| set @a = 1, reverse = 2
< error
| 1193 HY000 Unknown system variable 'reverse'
> exec ^
: string
| set @a = 1, revoke = 2
< error
| 1064 42000 …near 'revoke = 2' at line 1
> exec ^
: string
| set @a = 1, right = 2
< error
| 1064 42000 …near 'right = 2' at line 1
> exec ^
: string
| set @a = 1, rlike = 2
< error
| 1064 42000 …near 'rlike = 2' at line 1
> exec ^
: string
| set @a = 1, role = 2
< error
| 1193 HY000 Unknown system variable 'role'
> exec ^
: string
| set @a = 1, rollback = 2
< error
| 1193 HY000 Unknown system variable 'rollback'
> exec ^
: string
| set @a = 1, rollup = 2
< error
| 1193 HY000 Unknown system variable 'rollup'
> exec ^
: string
| set @a = 1, rotate = 2
< error
| 1193 HY000 Unknown system variable 'rotate'
> exec ^
: string
| set @a = 1, routine = 2
< error
| 1193 HY000 Unknown system variable 'routine'
> exec ^
: string
| set @a = 1, row = 2
< error
| 1064 42000 …near 'row = 2' at line 1
> exec ^
: string
| set @a = 1, row_count = 2
< error
| 1193 HY000 Unknown system variable 'row_count'
> exec ^
: string
| set @a = 1, row_format = 2
< error
| 1193 HY000 Unknown system variable 'row_format'
> exec ^
: string
| set @a = 1, row_number = 2
< error
| 1064 42000 …near 'row_number = 2' at line 1
> exec ^
: string
| set @a = 1, rows = 2
< error
| 1064 42000 …near 'rows = 2' at line 1
> exec ^
: string
| set @a = 1, rtree = 2
< error
| 1193 HY000 Unknown system variable 'rtree'
> exec ^
: string
| set @a = 1, s3 = 2
< error
| 1193 HY000 Unknown system variable 's3'
> exec ^
: string
| set @a = 1, savepoint = 2
< error
| 1193 HY000 Unknown system variable 'savepoint'
> exec ^
: string
| set @a = 1, schedule = 2
< error
| 1193 HY000 Unknown system variable 'schedule'
> exec ^
: string
| set @a = 1, schema = 2
< error
| 1064 42000 …near 'schema = 2' at line 1
> exec ^
: string
| set @a = 1, schema_name = 2
< error
| 1193 HY000 Unknown system variable 'schema_name'
> exec ^
: string
| set @a = 1, schemas = 2
< error
| 1064 42000 …near 'schemas = 2' at line 1
> exec ^
: string
| set @a = 1, second = 2
< error
| 1193 HY000 Unknown system variable 'second'
> exec ^
: string
| set @a = 1, second_microsecond = 2
< error
| 1064 42000 …near 'second_microsecond = 2' at line 1
> exec ^
: string
| set @a = 1, secondary = 2
< error
| 1193 HY000 Unknown system variable 'secondary'
> exec ^
: string
| set @a = 1, secondary_engine = 2
< error
| 1193 HY000 Unknown system variable 'secondary_engine'
> exec ^
: string
| set @a = 1, secondary_engine_attribute = 2
< error
| 1193 HY000 Unknown system variable 'secondary_engine_attribute'
> exec ^
: string
| set @a = 1, secondary_load = 2
< error
| 1193 HY000 Unknown system variable 'secondary_load'
> exec ^
: string
| set @a = 1, secondary_unload = 2
< error
| 1193 HY000 Unknown system variable 'secondary_unload'
> exec ^
: string
| set @a = 1, security = 2
< error
| 1193 HY000 Unknown system variable 'security'
> exec ^
: string
| set @a = 1, select = 2
< error
| 1064 42000 …near 'select = 2' at line 1
> exec ^
: string
| set @a = 1, sensitive = 2
< error
| 1064 42000 …near 'sensitive = 2' at line 1
> exec ^
: string
| set @a = 1, separator = 2
< error
| 1064 42000 …near 'separator = 2' at line 1
> exec ^
: string
| set @a = 1, serial = 2
< error
| 1193 HY000 Unknown system variable 'serial'
> exec ^
: string
| set @a = 1, serializable = 2
< error
| 1193 HY000 Unknown system variable 'serializable'
> exec ^
: string
| set @a = 1, server = 2
< error
| 1193 HY000 Unknown system variable 'server'
> exec ^
: string
| set @a = 1, session = 2
< error
| 1064 42000 …near '= 2' at line 1
> exec ^
: string
| set @a = 1, session_user = 2
< error
| 1193 HY000 Unknown system variable 'session_user'
> exec ^
: string
| set @a = 1, set = 2
< error
| 1064 42000 …near 'set = 2' at line 1
> exec ^
: string
| set @a = 1, sets = 2
< error
| 1193 HY000 Unknown system variable 'sets'
> exec ^
: string
| set @a = 1, share = 2
< error
| 1193 HY000 Unknown system variable 'share'
> exec ^
: string
| set @a = 1, show = 2
< error
| 1064 42000 …near 'show = 2' at line 1
> exec ^
: string
| set @a = 1, shutdown = 2
< error
| 1193 HY000 Unknown system variable 'shutdown'
> exec ^
: string
| set @a = 1, signal = 2
< error
| 1064 42000 …near 'signal = 2' at line 1
> exec ^
: string
| set @a = 1, signed = 2
< error
| 1193 HY000 Unknown system variable 'signed'
> exec ^
: string
| set @a = 1, simple = 2
< error
| 1193 HY000 Unknown system variable 'simple'
> exec ^
: string
| set @a = 1, skip = 2
< error
| 1193 HY000 Unknown system variable 'skip'
> exec ^
: string
| set @a = 1, slave = 2
< error
| 1193 HY000 Unknown system variable 'slave'
> exec ^
: string
| set @a = 1, slow = 2
< error
| 1193 HY000 Unknown system variable 'slow'
> exec ^
: string
| set @a = 1, smallint = 2
< error
| 1064 42000 …near 'smallint = 2' at line 1
> exec ^
: string
| set @a = 1, snapshot = 2
< error
| 1193 HY000 Unknown system variable 'snapshot'
> exec ^
: string
| set @a = 1, socket = 2
< error
| 1229 HY000 Variable 'socket' is a GLOBAL variable and should be set with SET GLOBAL
> exec ^
: string
| set @a = 1, socket = 2
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set @a = 1, some = 2
< error
| 1193 HY000 Unknown system variable 'some'
> exec ^
: string
| set @a = 1, soname = 2
< error
| 1193 HY000 Unknown system variable 'soname'
> exec ^
: string
| set @a = 1, sounds = 2
< error
| 1193 HY000 Unknown system variable 'sounds'
> exec ^
: string
| set @a = 1, source = 2
< error
| 1193 HY000 Unknown system variable 'source'
> exec ^
: string
| set @a = 1, source_auto_position = 2
< error
| 1193 HY000 Unknown system variable 'source_auto_position'
> exec ^
: string
| set @a = 1, source_bind = 2
< error
| 1193 HY000 Unknown system variable 'source_bind'
> exec ^
: string
| set @a = 1, source_compression_algorithms = 2
< error
| 1193 HY000 Unknown system variable 'source_compression_algorithms'
> exec ^
: string
| set @a = 1, source_connect_retry = 2
< error
| 1193 HY000 Unknown system variable 'source_connect_retry'
> exec ^
: string
| set @a = 1, source_connection_auto_failover = 2
< error
| 1193 HY000 Unknown system variable 'source_connection_auto_failover'
> exec ^
: string
| set @a = 1, source_delay = 2
< error
| 1193 HY000 Unknown system variable 'source_delay'
> exec ^
: string
| set @a = 1, source_heartbeat_period = 2
< error
| 1193 HY000 Unknown system variable 'source_heartbeat_period'
> exec ^
: string
| set @a = 1, source_host = 2
< error
| 1193 HY000 Unknown system variable 'source_host'
> exec ^
: string
| set @a = 1, source_log_file = 2
< error
| 1193 HY000 Unknown system variable 'source_log_file'
> exec ^
: string
| set @a = 1, source_log_pos = 2
< error
| 1193 HY000 Unknown system variable 'source_log_pos'
> exec ^
: string
| set @a = 1, source_password = 2
< error
| 1193 HY000 Unknown system variable 'source_password'
> exec ^
: string
| set @a = 1, source_port = 2
< error
| 1193 HY000 Unknown system variable 'source_port'
> exec ^
: string
| set @a = 1, source_public_key_path = 2
< error
| 1193 HY000 Unknown system variable 'source_public_key_path'
> exec ^
: string
| set @a = 1, source_retry_count = 2
< error
| 1193 HY000 Unknown system variable 'source_retry_count'
> exec ^
: string
| set @a = 1, source_ssl = 2
< error
| 1193 HY000 Unknown system variable 'source_ssl'
> exec ^
: string
| set @a = 1, source_ssl_ca = 2
< error
| 1193 HY000 Unknown system variable 'source_ssl_ca'
> exec ^
: string
| set @a = 1, source_ssl_capath = 2
< error
| 1193 HY000 Unknown system variable 'source_ssl_capath'
> exec ^
: string
| set @a = 1, source_ssl_cert = 2
< error
| 1193 HY000 Unknown system variable 'source_ssl_cert'
> exec ^
: string
| set @a = 1, source_ssl_cipher = 2
< error
| 1193 HY000 Unknown system variable 'source_ssl_cipher'
> exec ^
: string
| set @a = 1, source_ssl_crl = 2
< error
| 1193 HY000 Unknown system variable 'source_ssl_crl'
> exec ^
: string
| set @a = 1, source_ssl_crlpath = 2
< error
| 1193 HY000 Unknown system variable 'source_ssl_crlpath'
> exec ^
: string
| set @a = 1, source_ssl_key = 2
< error
| 1193 HY000 Unknown system variable 'source_ssl_key'
> exec ^
: string
| set @a = 1, source_ssl_verify_server_cert = 2
< error
| 1193 HY000 Unknown system variable 'source_ssl_verify_server_cert'
> exec ^
: string
| set @a = 1, source_tls_ciphersuites = 2
< error
| 1193 HY000 Unknown system variable 'source_tls_ciphersuites'
> exec ^
: string
| set @a = 1, source_tls_version = 2
< error
| 1193 HY000 Unknown system variable 'source_tls_version'
> exec ^
: string
| set @a = 1, source_user = 2
< error
| 1193 HY000 Unknown system variable 'source_user'
> exec ^
: string
| set @a = 1, source_zstd_compression_level = 2
< error
| 1193 HY000 Unknown system variable 'source_zstd_compression_level'
> exec ^
: string
| set @a = 1, spatial = 2
< error
| 1064 42000 …near 'spatial = 2' at line 1
> exec ^
: string
| set @a = 1, specific = 2
< error
| 1064 42000 …near 'specific = 2' at line 1
> exec ^
: string
| set @a = 1, sql = 2
< error
| 1064 42000 …near 'sql = 2' at line 1
> exec ^
: string
| set @a = 1, sql_after_gtids = 2
< error
| 1193 HY000 Unknown system variable 'sql_after_gtids'
> exec ^
: string
| set @a = 1, sql_after_mts_gaps = 2
< error
| 1193 HY000 Unknown system variable 'sql_after_mts_gaps'
> exec ^
: string
| set @a = 1, sql_before_gtids = 2
< error
| 1193 HY000 Unknown system variable 'sql_before_gtids'
> exec ^
: string
| set @a = 1, sql_big_result = 2
< error
| 1064 42000 …near 'sql_big_result = 2' at line 1
> exec ^
: string
| set @a = 1, sql_buffer_result = 2
< ok 0 0
> exec ^
: string
| set @a = 1, sql_calc_found_rows = 2
< error
| 1064 42000 …near 'sql_calc_found_rows = 2' at line 1
> exec ^
: string
| set @a = 1, sql_no_cache = 2
< error
| 1193 HY000 Unknown system variable 'sql_no_cache'
> exec ^
: string
| set @a = 1, sql_small_result = 2
< error
| 1064 42000 …near 'sql_small_result = 2' at line 1
> exec ^
: string
| set @a = 1, sql_thread = 2
< error
| 1193 HY000 Unknown system variable 'sql_thread'
> exec ^
: string
| set @a = 1, sql_tsi_day = 2
< error
| 1193 HY000 Unknown system variable 'sql_tsi_day'
> exec ^
: string
| set @a = 1, sql_tsi_hour = 2
< error
| 1193 HY000 Unknown system variable 'sql_tsi_hour'
> exec ^
: string
| set @a = 1, sql_tsi_minute = 2
< error
| 1193 HY000 Unknown system variable 'sql_tsi_minute'
> exec ^
: string
| set @a = 1, sql_tsi_month = 2
< error
| 1193 HY000 Unknown system variable 'sql_tsi_month'
> exec ^
: string
| set @a = 1, sql_tsi_quarter = 2
< error
| 1193 HY000 Unknown system variable 'sql_tsi_quarter'
> exec ^
: string
| set @a = 1, sql_tsi_second = 2
< error
| 1193 HY000 Unknown system variable 'sql_tsi_second'
> exec ^
: string
| set @a = 1, sql_tsi_week = 2
< error
| 1193 HY000 Unknown system variable 'sql_tsi_week'
> exec ^
: string
| set @a = 1, sql_tsi_year = 2
< error
| 1193 HY000 Unknown system variable 'sql_tsi_year'
> exec ^
: string
| set @a = 1, sqlexception = 2
< error
| 1064 42000 …near 'sqlexception = 2' at line 1
> exec ^
: string
| set @a = 1, sqlstate = 2
< error
| 1064 42000 …near 'sqlstate = 2' at line 1
> exec ^
: string
| set @a = 1, sqlwarning = 2
< error
| 1064 42000 …near 'sqlwarning = 2' at line 1
> exec ^
: string
| set @a = 1, srid = 2
< error
| 1193 HY000 Unknown system variable 'srid'
> exec ^
: string
| set @a = 1, sse_kms = 2
< error
| 1193 HY000 Unknown system variable 'sse_kms'
> exec ^
: string
| set @a = 1, sse_s3 = 2
< error
| 1193 HY000 Unknown system variable 'sse_s3'
> exec ^
: string
| set @a = 1, ssl = 2
< error
| 1064 42000 …near 'ssl = 2' at line 1
> exec ^
: string
| set @a = 1, st_collect = 2
< error
| 1193 HY000 Unknown system variable 'st_collect'
> exec ^
: string
| set @a = 1, stacked = 2
< error
| 1193 HY000 Unknown system variable 'stacked'
> exec ^
: string
| set @a = 1, start = 2
< error
| 1193 HY000 Unknown system variable 'start'
> exec ^
: string
| set @a = 1, starting = 2
< error
| 1064 42000 …near 'starting = 2' at line 1
> exec ^
: string
| set @a = 1, starts = 2
< error
| 1193 HY000 Unknown system variable 'starts'
> exec ^
: string
| set @a = 1, stats_auto_recalc = 2
< error
| 1193 HY000 Unknown system variable 'stats_auto_recalc'
> exec ^
: string
| set @a = 1, stats_persistent = 2
< error
| 1193 HY000 Unknown system variable 'stats_persistent'
> exec ^
: string
| set @a = 1, stats_sample_pages = 2
< error
| 1193 HY000 Unknown system variable 'stats_sample_pages'
> exec ^
: string
| set @a = 1, status = 2
< error
| 1193 HY000 Unknown system variable 'status'
> exec ^
: string
| set @a = 1, std = 2
< error
| 1193 HY000 Unknown system variable 'std'
> exec ^
: string
| set @a = 1, stddev = 2
< error
| 1193 HY000 Unknown system variable 'stddev'
> exec ^
: string
| set @a = 1, stddev_pop = 2
< error
| 1193 HY000 Unknown system variable 'stddev_pop'
> exec ^
: string
| set @a = 1, stddev_samp = 2
< error
| 1193 HY000 Unknown system variable 'stddev_samp'
> exec ^
: string
| set @a = 1, stop = 2
< error
| 1193 HY000 Unknown system variable 'stop'
> exec ^
: string
| set @a = 1, storage = 2
< error
| 1193 HY000 Unknown system variable 'storage'
> exec ^
: string
| set @a = 1, stored = 2
< error
| 1064 42000 …near 'stored = 2' at line 1
> exec ^
: string
| set @a = 1, straight_join = 2
< error
| 1064 42000 …near 'straight_join = 2' at line 1
> exec ^
: string
| set @a = 1, stream = 2
< error
| 1193 HY000 Unknown system variable 'stream'
> exec ^
: string
| set @a = 1, strict_load = 2
< error
| 1193 HY000 Unknown system variable 'strict_load'
> exec ^
: string
| set @a = 1, string = 2
< error
| 1193 HY000 Unknown system variable 'string'
> exec ^
: string
| set @a = 1, subclass_origin = 2
< error
| 1193 HY000 Unknown system variable 'subclass_origin'
> exec ^
: string
| set @a = 1, subdate = 2
< error
| 1193 HY000 Unknown system variable 'subdate'
> exec ^
: string
| set @a = 1, subject = 2
< error
| 1193 HY000 Unknown system variable 'subject'
> exec ^
: string
| set @a = 1, subpartition = 2
< error
| 1193 HY000 Unknown system variable 'subpartition'
> exec ^
: string
| set @a = 1, subpartitions = 2
< error
| 1193 HY000 Unknown system variable 'subpartitions'
> exec ^
: string
| set @a = 1, substr = 2
< error
| 1193 HY000 Unknown system variable 'substr'
> exec ^
: string
| set @a = 1, substring = 2
< error
| 1193 HY000 Unknown system variable 'substring'
> exec ^
: string
| set @a = 1, sum = 2
< error
| 1193 HY000 Unknown system variable 'sum'
> exec ^
: string
| set @a = 1, super = 2
< error
| 1193 HY000 Unknown system variable 'super'
> exec ^
: string
| set @a = 1, suspend = 2
< error
| 1193 HY000 Unknown system variable 'suspend'
> exec ^
: string
| set @a = 1, swaps = 2
< error
| 1193 HY000 Unknown system variable 'swaps'
> exec ^
: string
| set @a = 1, switches = 2
< error
| 1193 HY000 Unknown system variable 'switches'
> exec ^
: string
| set @a = 1, sysdate = 2
< error
| 1193 HY000 Unknown system variable 'sysdate'
> exec ^
: string
| set @a = 1, system = 2
< error
| 1064 42000 …near 'system = 2' at line 1
> exec ^
: string
| set @a = 1, system_user = 2
< error
| 1193 HY000 Unknown system variable 'system_user'
> exec ^
: string
| set @a = 1, table = 2
< error
| 1064 42000 …near 'table = 2' at line 1
> exec ^
: string
| set @a = 1, table_checksum = 2
< error
| 1193 HY000 Unknown system variable 'table_checksum'
> exec ^
: string
| set @a = 1, table_name = 2
< error
| 1193 HY000 Unknown system variable 'table_name'
> exec ^
: string
| set @a = 1, tables = 2
< error
| 1193 HY000 Unknown system variable 'tables'
> exec ^
: string
| set @a = 1, tablesample = 2
< error
| 1193 HY000 Unknown system variable 'tablesample'
> exec ^
: string
| set @a = 1, tablespace = 2
< error
| 1193 HY000 Unknown system variable 'tablespace'
> exec ^
: string
| set @a = 1, temporary = 2
< error
| 1193 HY000 Unknown system variable 'temporary'
> exec ^
: string
| set @a = 1, temptable = 2
< error
| 1193 HY000 Unknown system variable 'temptable'
> exec ^
: string
| set @a = 1, terminated = 2
< error
| 1064 42000 …near 'terminated = 2' at line 1
> exec ^
: string
| set @a = 1, text = 2
< error
| 1193 HY000 Unknown system variable 'text'
> exec ^
: string
| set @a = 1, than = 2
< error
| 1193 HY000 Unknown system variable 'than'
> exec ^
: string
| set @a = 1, then = 2
< error
| 1064 42000 …near 'then = 2' at line 1
> exec ^
: string
| set @a = 1, thread_priority = 2
< error
| 1193 HY000 Unknown system variable 'thread_priority'
> exec ^
: string
| set @a = 1, ties = 2
< error
| 1193 HY000 Unknown system variable 'ties'
> exec ^
: string
| set @a = 1, time = 2
< error
| 1193 HY000 Unknown system variable 'time'
> exec ^
: string
| set @a = 1, timestamp = 2
< ok 0 0
> exec ^
: string
| set @a = 1, timestampadd = 2
< error
| 1193 HY000 Unknown system variable 'timestampadd'
> exec ^
: string
| set @a = 1, timestampdiff = 2
< error
| 1193 HY000 Unknown system variable 'timestampdiff'
> exec ^
: string
| set @a = 1, tinyblob = 2
< error
| 1064 42000 …near 'tinyblob = 2' at line 1
> exec ^
: string
| set @a = 1, tinyint = 2
< error
| 1064 42000 …near 'tinyint = 2' at line 1
> exec ^
: string
| set @a = 1, tinytext = 2
< error
| 1064 42000 …near 'tinytext = 2' at line 1
> exec ^
: string
| set @a = 1, tls = 2
< error
| 1193 HY000 Unknown system variable 'tls'
> exec ^
: string
| set @a = 1, to = 2
< error
| 1064 42000 …near 'to = 2' at line 1
> exec ^
: string
| set @a = 1, trailing = 2
< error
| 1064 42000 …near 'trailing = 2' at line 1
> exec ^
: string
| set @a = 1, transaction = 2
< error
| 1193 HY000 Unknown system variable 'transaction'
> exec ^
: string
| set @a = 1, trigger = 2
< error
| 1064 42000 …near 'trigger = 2' at line 1
> exec ^
: string
| set @a = 1, triggers = 2
< error
| 1193 HY000 Unknown system variable 'triggers'
> exec ^
: string
| set @a = 1, trim = 2
< error
| 1193 HY000 Unknown system variable 'trim'
> exec ^
: string
| set @a = 1, true = 2
< error
| 1064 42000 …near 'true = 2' at line 1
> exec ^
: string
| set @a = 1, truncate = 2
< error
| 1193 HY000 Unknown system variable 'truncate'
> exec ^
: string
| set @a = 1, type = 2
< error
| 1193 HY000 Unknown system variable 'type'
> exec ^
: string
| set @a = 1, types = 2
< error
| 1193 HY000 Unknown system variable 'types'
> exec ^
: string
| set @a = 1, unbounded = 2
< error
| 1193 HY000 Unknown system variable 'unbounded'
> exec ^
: string
| set @a = 1, uncommitted = 2
< error
| 1193 HY000 Unknown system variable 'uncommitted'
> exec ^
: string
| set @a = 1, undefined = 2
< error
| 1193 HY000 Unknown system variable 'undefined'
> exec ^
: string
| set @a = 1, undo = 2
< error
| 1064 42000 …near 'undo = 2' at line 1
> exec ^
: string
| set @a = 1, undo_buffer_size = 2
< error
| 1193 HY000 Unknown system variable 'undo_buffer_size'
> exec ^
: string
| set @a = 1, undofile = 2
< error
| 1193 HY000 Unknown system variable 'undofile'
> exec ^
: string
| set @a = 1, unicode = 2
< error
| 1193 HY000 Unknown system variable 'unicode'
> exec ^
: string
| set @a = 1, uninstall = 2
< error
| 1193 HY000 Unknown system variable 'uninstall'
> exec ^
: string
| set @a = 1, union = 2
< error
| 1064 42000 …near 'union = 2' at line 1
> exec ^
: string
| set @a = 1, unique = 2
< error
| 1064 42000 …near 'unique = 2' at line 1
> exec ^
: string
| set @a = 1, unknown = 2
< error
| 1193 HY000 Unknown system variable 'unknown'
> exec ^
: string
| set @a = 1, unlock = 2
< error
| 1064 42000 …near 'unlock = 2' at line 1
> exec ^
: string
| set @a = 1, unregister = 2
< error
| 1193 HY000 Unknown system variable 'unregister'
> exec ^
: string
| set @a = 1, unsigned = 2
< error
| 1064 42000 …near 'unsigned = 2' at line 1
> exec ^
: string
| set @a = 1, until = 2
< error
| 1193 HY000 Unknown system variable 'until'
> exec ^
: string
| set @a = 1, update = 2
< error
| 1064 42000 …near 'update = 2' at line 1
> exec ^
: string
| set @a = 1, upgrade = 2
< error
| 1193 HY000 Unknown system variable 'upgrade'
> exec ^
: string
| set @a = 1, uri = 2
< error
| 1193 HY000 Unknown system variable 'uri'
> exec ^
: string
| set @a = 1, url = 2
< error
| 1193 HY000 Unknown system variable 'url'
> exec ^
: string
| set @a = 1, usage = 2
< error
| 1064 42000 …near 'usage = 2' at line 1
> exec ^
: string
| set @a = 1, use = 2
< error
| 1064 42000 …near 'use = 2' at line 1
> exec ^
: string
| set @a = 1, use_frm = 2
< error
| 1193 HY000 Unknown system variable 'use_frm'
> exec ^
: string
| set @a = 1, user = 2
< error
| 1193 HY000 Unknown system variable 'user'
> exec ^
: string
| set @a = 1, user_resources = 2
< error
| 1193 HY000 Unknown system variable 'user_resources'
> exec ^
: string
| set @a = 1, using = 2
< error
| 1064 42000 …near 'using = 2' at line 1
> exec ^
: string
| set @a = 1, utc_date = 2
< error
| 1064 42000 …near 'utc_date = 2' at line 1
> exec ^
: string
| set @a = 1, utc_time = 2
< error
| 1064 42000 …near 'utc_time = 2' at line 1
> exec ^
: string
| set @a = 1, utc_timestamp = 2
< error
| 1064 42000 …near 'utc_timestamp = 2' at line 1
> exec ^
: string
| set @a = 1, validate = 2
< error
| 1193 HY000 Unknown system variable 'validate'
> exec ^
: string
| set @a = 1, validation = 2
< error
| 1193 HY000 Unknown system variable 'validation'
> exec ^
: string
| set @a = 1, value = 2
< error
| 1193 HY000 Unknown system variable 'value'
> exec ^
: string
| set @a = 1, values = 2
< error
| 1064 42000 …near 'values = 2' at line 1
> exec ^
: string
| set @a = 1, var_pop = 2
< error
| 1193 HY000 Unknown system variable 'var_pop'
> exec ^
: string
| set @a = 1, var_samp = 2
< error
| 1193 HY000 Unknown system variable 'var_samp'
> exec ^
: string
| set @a = 1, varbinary = 2
< error
| 1064 42000 …near 'varbinary = 2' at line 1
> exec ^
: string
| set @a = 1, varchar = 2
< error
| 1064 42000 …near 'varchar = 2' at line 1
> exec ^
: string
| set @a = 1, varcharacter = 2
< error
| 1064 42000 …near 'varcharacter = 2' at line 1
> exec ^
: string
| set @a = 1, variables = 2
< error
| 1193 HY000 Unknown system variable 'variables'
> exec ^
: string
| set @a = 1, variance = 2
< error
| 1193 HY000 Unknown system variable 'variance'
> exec ^
: string
| set @a = 1, varying = 2
< error
| 1064 42000 …near 'varying = 2' at line 1
> exec ^
: string
| set @a = 1, vcpu = 2
< error
| 1193 HY000 Unknown system variable 'vcpu'
> exec ^
: string
| set @a = 1, vector = 2
< error
| 1193 HY000 Unknown system variable 'vector'
> exec ^
: string
| set @a = 1, verify_key_constraints = 2
< error
| 1193 HY000 Unknown system variable 'verify_key_constraints'
> exec ^
: string
| set @a = 1, view = 2
< error
| 1193 HY000 Unknown system variable 'view'
> exec ^
: string
| set @a = 1, virtual = 2
< error
| 1064 42000 …near 'virtual = 2' at line 1
> exec ^
: string
| set @a = 1, visible = 2
< error
| 1193 HY000 Unknown system variable 'visible'
> exec ^
: string
| set @a = 1, wait = 2
< error
| 1193 HY000 Unknown system variable 'wait'
> exec ^
: string
| set @a = 1, warnings = 2
< error
| 1193 HY000 Unknown system variable 'warnings'
> exec ^
: string
| set @a = 1, week = 2
< error
| 1193 HY000 Unknown system variable 'week'
> exec ^
: string
| set @a = 1, weight_string = 2
< error
| 1193 HY000 Unknown system variable 'weight_string'
> exec ^
: string
| set @a = 1, when = 2
< error
| 1064 42000 …near 'when = 2' at line 1
> exec ^
: string
| set @a = 1, where = 2
< error
| 1064 42000 …near 'where = 2' at line 1
> exec ^
: string
| set @a = 1, while = 2
< error
| 1064 42000 …near 'while = 2' at line 1
> exec ^
: string
| set @a = 1, window = 2
< error
| 1064 42000 …near 'window = 2' at line 1
> exec ^
: string
| set @a = 1, with = 2
< error
| 1064 42000 …near 'with = 2' at line 1
> exec ^
: string
| set @a = 1, without = 2
< error
| 1193 HY000 Unknown system variable 'without'
> exec ^
: string
| set @a = 1, work = 2
< error
| 1193 HY000 Unknown system variable 'work'
> exec ^
: string
| set @a = 1, wrapper = 2
< error
| 1193 HY000 Unknown system variable 'wrapper'
> exec ^
: string
| set @a = 1, write = 2
< error
| 1064 42000 …near 'write = 2' at line 1
> exec ^
: string
| set @a = 1, x509 = 2
< error
| 1193 HY000 Unknown system variable 'x509'
> exec ^
: string
| set @a = 1, xa = 2
< error
| 1193 HY000 Unknown system variable 'xa'
> exec ^
: string
| set @a = 1, xid = 2
< error
| 1193 HY000 Unknown system variable 'xid'
> exec ^
: string
| set @a = 1, xml = 2
< error
| 1193 HY000 Unknown system variable 'xml'
> exec ^
: string
| set @a = 1, xor = 2
< error
| 1064 42000 …near 'xor = 2' at line 1
> exec ^
: string
| set @a = 1, year = 2
< error
| 1193 HY000 Unknown system variable 'year'
> exec ^
: string
| set @a = 1, year_month = 2
< error
| 1064 42000 …near 'year_month = 2' at line 1
> exec ^
: string
| set @a = 1, zerofill = 2
< error
| 1064 42000 …near 'zerofill = 2' at line 1
> exec ^
: string
| set @a = 1, zone = 2
< error
| 1193 HY000 Unknown system variable 'zone'
