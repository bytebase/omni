-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestSetVariableNameSweepOceanBaseOracle default#1
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| set absent = 1
< ok 0 0
> exec ^
: string
| set accessible = 1
< error
| 1064 42000 …near 'accessible = 1' at line 1
> exec ^
: string
| set account = 1
< ok 0 0
> exec ^
: string
| set action = 1
< ok 0 0
> exec ^
: string
| set active = 1
< ok 0 0
> exec ^
: string
| set add = 1
< error
| 1064 42000 …near 'add = 1' at line 1
> exec ^
: string
| set adddate = 1
< ok 0 0
> exec ^
: string
| set admin = 1
< ok 0 0
> exec ^
: string
| set after = 1
< ok 0 0
> exec ^
: string
| set against = 1
< ok 0 0
> exec ^
: string
| set aggregate = 1
< ok 0 0
> exec ^
: string
| set algorithm = 1
< ok 0 0
> exec ^
: string
| set all = 1
< error
| 1064 42000 …near 'all = 1' at line 1
> exec ^
: string
| set allow_missing_files = 1
< ok 0 0
> exec ^
: string
| set alter = 1
< error
| 1064 42000 …near 'alter = 1' at line 1
> exec ^
: string
| set always = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set array = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set asensitive = 1
< error
| 1064 42000 …near 'asensitive = 1' at line 1
> exec ^
: string
| set assign_gtids_to_anonymous_transactions = 1
< ok 0 0
> exec ^
: string
| set at = 1
< ok 0 0
> exec ^
: string
| set attribute = 1
< ok 0 0
> exec ^
: string
| set authentication = 1
< ok 0 0
> exec ^
: string
| set auto = 1
< ok 0 0
> exec ^
: string
| set auto_increment = 1
< ok 0 0
> exec ^
: string
| set auto_increment_mode = 1
< ok 0 0
> exec ^
: string
| set auto_refresh = 1
< ok 0 0
> exec ^
: string
| set auto_refresh_source = 1
< ok 0 0
> exec ^
: string
| set autoextend_size = 1
< ok 0 0
> exec ^
: string
| set avg = 1
< ok 0 0
> exec ^
: string
| set avg_row_length = 1
< ok 0 0
> exec ^
: string
| set backup = 1
< ok 0 0
> exec ^
: string
| set before = 1
< error
| 1064 42000 …near 'before = 1' at line 1
> exec ^
: string
| set begin = 1
< ok 0 0
> exec ^
: string
| set bernoulli = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set bit = 1
< ok 0 0
> exec ^
: string
| set bit_and = 1
< ok 0 0
> exec ^
: string
| set bit_or = 1
< ok 0 0
> exec ^
: string
| set bit_xor = 1
< ok 0 0
> exec ^
: string
| set blob = 1
< error
| 1064 42000 …near 'blob = 1' at line 1
> exec ^
: string
| set block = 1
< ok 0 0
> exec ^
: string
| set block_size = 1
< ok 0 0
> exec ^
: string
| set bool = 1
< ok 0 0
> exec ^
: string
| set boolean = 1
< ok 0 0
> exec ^
: string
| set both = 1
< error
| 1064 42000 …near 'both = 1' at line 1
> exec ^
: string
| set btree = 1
< ok 0 0
> exec ^
: string
| set buckets = 1
< ok 0 0
> exec ^
: string
| set bulk = 1
< ok 0 0
> exec ^
: string
| set by = 1
< error
| 1064 42000 …near 'by = 1' at line 1
> exec ^
: string
| set byte = 1
< ok 0 0
> exec ^
: string
| set cache = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set case = 1
< error
| 1064 42000 …near 'case = 1' at line 1
> exec ^
: string
| set cast = 1
< ok 0 0
> exec ^
: string
| set catalog_name = 1
< ok 0 0
> exec ^
: string
| set chain = 1
< ok 0 0
> exec ^
: string
| set challenge_response = 1
< ok 0 0
> exec ^
: string
| set change = 1
< error
| 1064 42000 …near 'change = 1' at line 1
> exec ^
: string
| set changed = 1
< ok 0 0
> exec ^
: string
| set channel = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set check = 1
< error
| 1064 42000 …near 'check = 1' at line 1
> exec ^
: string
| set checksum = 1
< ok 0 0
> exec ^
: string
| set cipher = 1
< ok 0 0
> exec ^
: string
| set class_origin = 1
< ok 0 0
> exec ^
: string
| set client = 1
< ok 0 0
> exec ^
: string
| set clone = 1
< ok 0 0
> exec ^
: string
| set close = 1
< ok 0 0
> exec ^
: string
| set coalesce = 1
< ok 0 0
> exec ^
: string
| set code = 1
< ok 0 0
> exec ^
: string
| set collate = 1
< error
| 1064 42000 …near 'collate = 1' at line 1
> exec ^
: string
| set collation = 1
< ok 0 0
> exec ^
: string
| set column = 1
< error
| 1064 42000 …near 'column = 1' at line 1
> exec ^
: string
| set column_format = 1
< ok 0 0
> exec ^
: string
| set column_name = 1
< ok 0 0
> exec ^
: string
| set columns = 1
< ok 0 0
> exec ^
: string
| set comment = 1
< ok 0 0
> exec ^
: string
| set commit = 1
< ok 0 0
> exec ^
: string
| set committed = 1
< ok 0 0
> exec ^
: string
| set compact = 1
< ok 0 0
> exec ^
: string
| set completion = 1
< ok 0 0
> exec ^
: string
| set component = 1
< ok 0 0
> exec ^
: string
| set compressed = 1
< ok 0 0
> exec ^
: string
| set compression = 1
< ok 0 0
> exec ^
: string
| set concat = 1
< ok 0 0
> exec ^
: string
| set concurrent = 1
< ok 0 0
> exec ^
: string
| set condition = 1
< error
| 1064 42000 …near 'condition = 1' at line 1
> exec ^
: string
| set connection = 1
< ok 0 0
> exec ^
: string
| set consistent = 1
< ok 0 0
> exec ^
: string
| set constraint = 1
< error
| 1064 42000 …near 'constraint = 1' at line 1
> exec ^
: string
| set constraint_catalog = 1
< ok 0 0
> exec ^
: string
| set constraint_name = 1
< ok 0 0
> exec ^
: string
| set constraint_schema = 1
< ok 0 0
> exec ^
: string
| set contains = 1
< ok 0 0
> exec ^
: string
| set context = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set cpu = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set cube = 1
< ok 0 0
> exec ^
: string
| set cume_dist = 1
< ok 0 0
> exec ^
: string
| set curdate = 1
< ok 0 0
> exec ^
: string
| set current = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set curtime = 1
< ok 0 0
> exec ^
: string
| set data = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set date = 1
< ok 0 0
> exec ^
: string
| set date_add = 1
< ok 0 0
> exec ^
: string
| set date_sub = 1
< ok 0 0
> exec ^
: string
| set datetime = 1
< ok 0 0
> exec ^
: string
| set day = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set dec = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set definer = 1
< ok 0 0
> exec ^
: string
| set definition = 1
< ok 0 0
> exec ^
: string
| set delay_key_write = 1
< ok 0 0
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
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set deterministic = 1
< error
| 1064 42000 …near 'deterministic = 1' at line 1
> exec ^
: string
| set diagnostics = 1
< ok 0 0
> exec ^
: string
| set directory = 1
< ok 0 0
> exec ^
: string
| set disable = 1
< ok 0 0
> exec ^
: string
| set discard = 1
< ok 0 0
> exec ^
: string
| set disk = 1
< ok 0 0
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
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set dumpfile = 1
< ok 0 0
> exec ^
: string
| set duplicate = 1
< ok 0 0
> exec ^
: string
| set dynamic = 1
< ok 0 0
> exec ^
: string
| set dynamic_partition_policy = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set enable = 1
< ok 0 0
> exec ^
: string
| set enable_macro_block_bloom_filter = 1
< ok 0 0
> exec ^
: string
| set enclosed = 1
< error
| 1064 42000 …near 'enclosed = 1' at line 1
> exec ^
: string
| set encryption = 1
< ok 0 0
> exec ^
: string
| set end = 1
< ok 0 0
> exec ^
: string
| set ends = 1
< ok 0 0
> exec ^
: string
| set enforced = 1
< ok 0 0
> exec ^
: string
| set engine = 1
< ok 0 0
> exec ^
: string
| set engine_attribute = 1
< ok 0 0
> exec ^
: string
| set engines = 1
< ok 0 0
> exec ^
: string
| set enum = 1
< ok 0 0
> exec ^
: string
| set error = 1
< ok 0 0
> exec ^
: string
| set errors = 1
< ok 0 0
> exec ^
: string
| set escape = 1
< ok 0 0
> exec ^
: string
| set escaped = 1
< error
| 1064 42000 …near 'escaped = 1' at line 1
> exec ^
: string
| set event = 1
< error
| 1064 42000 …near 'event = 1' at line 1
> exec ^
: string
| set events = 1
< ok 0 0
> exec ^
: string
| set every = 1
< ok 0 0
> exec ^
: string
| set except = 1
< ok 0 0
> exec ^
: string
| set exchange = 1
< ok 0 0
> exec ^
: string
| set exclude = 1
< ok 0 0
> exec ^
: string
| set execute = 1
< error
| 1064 42000 …near 'execute = 1' at line 1
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
< ok 0 0
> exec ^
: string
| set expire = 1
< ok 0 0
> exec ^
: string
| set explain = 1
< error
| 1064 42000 …near 'explain = 1' at line 1
> exec ^
: string
| set export = 1
< ok 0 0
> exec ^
: string
| set extended = 1
< ok 0 0
> exec ^
: string
| set extent_size = 1
< ok 0 0
> exec ^
: string
| set external = 1
< ok 0 0
> exec ^
: string
| set external_format = 1
< ok 0 0
> exec ^
: string
| set extract = 1
< ok 0 0
> exec ^
: string
| set factor = 1
< ok 0 0
> exec ^
: string
| set failed_login_attempts = 1
< ok 0 0
> exec ^
: string
| set false = 1
< error
| 1064 42000 …near 'false = 1' at line 1
> exec ^
: string
| set fast = 1
< ok 0 0
> exec ^
: string
| set faults = 1
< ok 0 0
> exec ^
: string
| set fetch = 1
< error
| 1064 42000 …near 'fetch = 1' at line 1
> exec ^
: string
| set fields = 1
< ok 0 0
> exec ^
: string
| set file = 1
< error
| 1064 42000 …near 'file = 1' at line 1
> exec ^
: string
| set file_block_size = 1
< ok 0 0
> exec ^
: string
| set file_format = 1
< ok 0 0
> exec ^
: string
| set file_name = 1
< ok 0 0
> exec ^
: string
| set file_pattern = 1
< ok 0 0
> exec ^
: string
| set file_prefix = 1
< ok 0 0
> exec ^
: string
| set files = 1
< ok 0 0
> exec ^
: string
| set filter = 1
< ok 0 0
> exec ^
: string
| set finish = 1
< ok 0 0
> exec ^
: string
| set first = 1
< ok 0 0
> exec ^
: string
| set first_value = 1
< ok 0 0
> exec ^
: string
| set fixed = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set following = 1
< ok 0 0
> exec ^
: string
| set follows = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set found = 1
< ok 0 0
> exec ^
: string
| set from = 1
< error
| 1064 42000 …near 'from = 1' at line 1
> exec ^
: string
| set full = 1
< ok 0 0
> exec ^
: string
| set fulltext = 1
< error
| 1064 42000 …near 'fulltext = 1' at line 1
> exec ^
: string
| set function = 1
< ok 0 0
> exec ^
: string
| set general = 1
< ok 0 0
> exec ^
: string
| set generate = 1
< ok 0 0
> exec ^
: string
| set generated = 1
< error
| 1064 42000 …near 'generated = 1' at line 1
> exec ^
: string
| set geomcollection = 1
< ok 0 0
> exec ^
: string
| set geometry = 1
< ok 0 0
> exec ^
: string
| set geometrycollection = 1
< ok 0 0
> exec ^
: string
| set get = 1
< error
| 1064 42000 …near 'get = 1' at line 1
> exec ^
: string
| set get_format = 1
< ok 0 0
> exec ^
: string
| set get_master_public_key = 1
< ok 0 0
> exec ^
: string
| set get_source_public_key = 1
< ok 0 0
> exec ^
: string
| set global = 1
< ok 0 0
> exec ^
: string
| set grant = 1
< error
| 1064 42000 …near 'grant = 1' at line 1
> exec ^
: string
| set grants = 1
< ok 0 0
> exec ^
: string
| set greatest = 1
< ok 0 0
> exec ^
: string
| set group = 1
< error
| 1064 42000 …near 'group = 1' at line 1
> exec ^
: string
| set group_concat = 1
< ok 0 0
> exec ^
: string
| set group_replication = 1
< ok 0 0
> exec ^
: string
| set grouping = 1
< ok 0 0
> exec ^
: string
| set groups = 1
< ok 0 0
> exec ^
: string
| set gtid_only = 1
< ok 0 0
> exec ^
: string
| set gtids = 1
< ok 0 0
> exec ^
: string
| set guided = 1
< ok 0 0
> exec ^
: string
| set handler = 1
< ok 0 0
> exec ^
: string
| set hash = 1
< ok 0 0
> exec ^
: string
| set having = 1
< error
| 1064 42000 …near 'having = 1' at line 1
> exec ^
: string
| set header = 1
< ok 0 0
> exec ^
: string
| set help = 1
< ok 0 0
> exec ^
: string
| set high_priority = 1
< error
| 1064 42000 …near 'high_priority = 1' at line 1
> exec ^
: string
| set histogram = 1
< ok 0 0
> exec ^
: string
| set history = 1
< ok 0 0
> exec ^
: string
| set host = 1
< ok 0 0
> exec ^
: string
| set hosts = 1
< ok 0 0
> exec ^
: string
| set hour = 1
< ok 0 0
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
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set import = 1
< ok 0 0
> exec ^
: string
| set in = 1
< error
| 1064 42000 …near 'in = 1' at line 1
> exec ^
: string
| set inactive = 1
< ok 0 0
> exec ^
: string
| set index = 1
< error
| 1064 42000 …near 'index = 1' at line 1
> exec ^
: string
| set indexes = 1
< ok 0 0
> exec ^
: string
| set infile = 1
< error
| 1064 42000 …near 'infile = 1' at line 1
> exec ^
: string
| set initial = 1
< ok 0 0
> exec ^
: string
| set initial_size = 1
< ok 0 0
> exec ^
: string
| set initiate = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set install = 1
< ok 0 0
> exec ^
: string
| set instance = 1
< ok 0 0
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
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set invoker = 1
< ok 0 0
> exec ^
: string
| set io = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set ipc = 1
< ok 0 0
> exec ^
: string
| set is = 1
< error
| 1064 42000 …near 'is = 1' at line 1
> exec ^
: string
| set isolation = 1
< ok 0 0
> exec ^
: string
| set issuer = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set json_arrayagg = 1
< ok 0 0
> exec ^
: string
| set json_duality_object = 1
< ok 0 0
> exec ^
: string
| set json_objectagg = 1
< ok 0 0
> exec ^
: string
| set json_table = 1
< ok 0 0
> exec ^
: string
| set json_value = 1
< ok 0 0
> exec ^
: string
| set key = 1
< error
| 1064 42000 …near 'key = 1' at line 1
> exec ^
: string
| set key_block_size = 1
< ok 0 0
> exec ^
: string
| set keyring = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set language = 1
< ok 0 0
> exec ^
: string
| set last = 1
< ok 0 0
> exec ^
: string
| set last_value = 1
< ok 0 0
> exec ^
: string
| set lateral = 1
< ok 0 0
> exec ^
: string
| set lead = 1
< ok 0 0
> exec ^
: string
| set leading = 1
< error
| 1064 42000 …near 'leading = 1' at line 1
> exec ^
: string
| set least = 1
< ok 0 0
> exec ^
: string
| set leave = 1
< error
| 1064 42000 …near 'leave = 1' at line 1
> exec ^
: string
| set leaves = 1
< ok 0 0
> exec ^
: string
| set left = 1
< error
| 1064 42000 …near 'left = 1' at line 1
> exec ^
: string
| set less = 1
< ok 0 0
> exec ^
: string
| set level = 1
< ok 0 0
> exec ^
: string
| set library = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set list = 1
< ok 0 0
> exec ^
: string
| set load = 1
< error
| 1064 42000 …near 'load = 1' at line 1
> exec ^
: string
| set local = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set locks = 1
< ok 0 0
> exec ^
: string
| set log = 1
< ok 0 0
> exec ^
: string
| set logfile = 1
< ok 0 0
> exec ^
: string
| set logs = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set manual = 1
< ok 0 0
> exec ^
: string
| set master = 1
< ok 0 0
> exec ^
: string
| set master_auto_position = 1
< ok 0 0
> exec ^
: string
| set master_bind = 1
< error
| 1064 42000 …near 'master_bind = 1' at line 1
> exec ^
: string
| set master_compression_algorithms = 1
< ok 0 0
> exec ^
: string
| set master_connect_retry = 1
< ok 0 0
> exec ^
: string
| set master_delay = 1
< ok 0 0
> exec ^
: string
| set master_heartbeat_period = 1
< ok 0 0
> exec ^
: string
| set master_host = 1
< ok 0 0
> exec ^
: string
| set master_log_file = 1
< ok 0 0
> exec ^
: string
| set master_log_pos = 1
< ok 0 0
> exec ^
: string
| set master_password = 1
< ok 0 0
> exec ^
: string
| set master_port = 1
< ok 0 0
> exec ^
: string
| set master_public_key_path = 1
< ok 0 0
> exec ^
: string
| set master_retry_count = 1
< ok 0 0
> exec ^
: string
| set master_ssl = 1
< ok 0 0
> exec ^
: string
| set master_ssl_ca = 1
< ok 0 0
> exec ^
: string
| set master_ssl_capath = 1
< ok 0 0
> exec ^
: string
| set master_ssl_cert = 1
< ok 0 0
> exec ^
: string
| set master_ssl_cipher = 1
< ok 0 0
> exec ^
: string
| set master_ssl_crl = 1
< ok 0 0
> exec ^
: string
| set master_ssl_crlpath = 1
< ok 0 0
> exec ^
: string
| set master_ssl_key = 1
< ok 0 0
> exec ^
: string
| set master_ssl_verify_server_cert = 1
< error
| 1064 42000 …near 'master_ssl_verify_server_cert = 1' at line 1
> exec ^
: string
| set master_tls_ciphersuites = 1
< ok 0 0
> exec ^
: string
| set master_tls_version = 1
< ok 0 0
> exec ^
: string
| set master_user = 1
< ok 0 0
> exec ^
: string
| set master_zstd_compression_level = 1
< ok 0 0
> exec ^
: string
| set match = 1
< error
| 1064 42000 …near 'match = 1' at line 1
> exec ^
: string
| set materialized = 1
< ok 0 0
> exec ^
: string
| set max = 1
< ok 0 0
> exec ^
: string
| set max_connections_per_hour = 1
< ok 0 0
> exec ^
: string
| set max_queries_per_hour = 1
< ok 0 0
> exec ^
: string
| set max_rows = 1
< ok 0 0
> exec ^
: string
| set max_size = 1
< ok 0 0
> exec ^
: string
| set max_updates_per_hour = 1
< ok 0 0
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
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set memory = 1
< ok 0 0
> exec ^
: string
| set merge = 1
< ok 0 0
> exec ^
: string
| set message_text = 1
< ok 0 0
> exec ^
: string
| set microsecond = 1
< ok 0 0
> exec ^
: string
| set mid = 1
< ok 0 0
> exec ^
: string
| set middleint = 1
< error
| 1064 42000 …near 'middleint = 1' at line 1
> exec ^
: string
| set migrate = 1
< ok 0 0
> exec ^
: string
| set min = 1
< ok 0 0
> exec ^
: string
| set min_rows = 1
< ok 0 0
> exec ^
: string
| set minute = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set modifies = 1
< error
| 1064 42000 …near 'modifies = 1' at line 1
> exec ^
: string
| set modify = 1
< ok 0 0
> exec ^
: string
| set month = 1
< ok 0 0
> exec ^
: string
| set multilinestring = 1
< ok 0 0
> exec ^
: string
| set multipoint = 1
< ok 0 0
> exec ^
: string
| set multipolygon = 1
< ok 0 0
> exec ^
: string
| set mutex = 1
< ok 0 0
> exec ^
: string
| set mysql_errno = 1
< ok 0 0
> exec ^
: string
| set name = 1
< ok 0 0
> exec ^
: string
| set names = 1
< ok 0 0
> exec ^
: string
| set national = 1
< ok 0 0
> exec ^
: string
| set natural = 1
< error
| 1064 42000 …near 'natural = 1' at line 1
> exec ^
: string
| set nchar = 1
< ok 0 0
> exec ^
: string
| set ndb = 1
< ok 0 0
> exec ^
: string
| set ndbcluster = 1
< ok 0 0
> exec ^
: string
| set nested = 1
< ok 0 0
> exec ^
: string
| set network_namespace = 1
< ok 0 0
> exec ^
: string
| set never = 1
< ok 0 0
> exec ^
: string
| set new = 1
< ok 0 0
> exec ^
: string
| set next = 1
< ok 0 0
> exec ^
: string
| set no = 1
< ok 0 0
> exec ^
: string
| set no_wait = 1
< ok 0 0
> exec ^
: string
| set no_write_to_binlog = 1
< error
| 1064 42000 …near 'no_write_to_binlog = 1' at line 1
> exec ^
: string
| set nodegroup = 1
< ok 0 0
> exec ^
: string
| set none = 1
< error
| 1064 42000 …near 'none = 1' at line 1
> exec ^
: string
| set not = 1
< error
| 1064 42000 …near 'not = 1' at line 1
> exec ^
: string
| set now = 1
< ok 0 0
> exec ^
: string
| set nowait = 1
< ok 0 0
> exec ^
: string
| set nth_value = 1
< ok 0 0
> exec ^
: string
| set ntile = 1
< ok 0 0
> exec ^
: string
| set null = 1
< error
| 1064 42000 …near 'null = 1' at line 1
> exec ^
: string
| set nullif = 1
< ok 0 0
> exec ^
: string
| set nulls = 1
< ok 0 0
> exec ^
: string
| set number = 1
< ok 0 0
> exec ^
: string
| set numeric = 1
< error
| 1064 42000 …near 'numeric = 1' at line 1
> exec ^
: string
| set nvarchar = 1
< ok 0 0
> exec ^
: string
| set of = 1
< ok 0 0
> exec ^
: string
| set off = 1
< ok 0 0
> exec ^
: string
| set offset = 1
< ok 0 0
> exec ^
: string
| set oj = 1
< ok 0 0
> exec ^
: string
| set old = 1
< ok 0 0
> exec ^
: string
| set on = 1
< error
| 1064 42000 …near 'on = 1' at line 1
> exec ^
: string
| set one = 1
< ok 0 0
> exec ^
: string
| set only = 1
< ok 0 0
> exec ^
: string
| set open = 1
< ok 0 0
> exec ^
: string
| set optimize = 1
< error
| 1064 42000 …near 'optimize = 1' at line 1
> exec ^
: string
| set optimizer_costs = 1
< ok 0 0
> exec ^
: string
| set option = 1
< error
| 1064 42000 …near '= 1' at line 1
> exec ^
: string
| set optional = 1
< ok 0 0
> exec ^
: string
| set optionally = 1
< error
| 1064 42000 …near 'optionally = 1' at line 1
> exec ^
: string
| set options = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set organization = 1
< ok 0 0
> exec ^
: string
| set others = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set overlay = 1
< ok 0 0
> exec ^
: string
| set overwrite = 1
< ok 0 0
> exec ^
: string
| set owner = 1
< ok 0 0
> exec ^
: string
| set pack_keys = 1
< ok 0 0
> exec ^
: string
| set page = 1
< ok 0 0
> exec ^
: string
| set parallel = 1
< ok 0 0
> exec ^
: string
| set parameters = 1
< ok 0 0
> exec ^
: string
| set parse_tree = 1
< ok 0 0
> exec ^
: string
| set parser = 1
< ok 0 0
> exec ^
: string
| set partial = 1
< ok 0 0
> exec ^
: string
| set partition = 1
< error
| 1064 42000 …near 'partition = 1' at line 1
> exec ^
: string
| set partitioning = 1
< ok 0 0
> exec ^
: string
| set partitions = 1
< ok 0 0
> exec ^
: string
| set password = 1
< error
| 1064 42000 …near '1' at line 1
> exec ^
: string
| set password_lock_time = 1
< ok 0 0
> exec ^
: string
| set path = 1
< ok 0 0
> exec ^
: string
| set pctfree = 1
< ok 0 0
> exec ^
: string
| set percent_rank = 1
< ok 0 0
> exec ^
: string
| set persist = 1
< ok 0 0
> exec ^
: string
| set persist_only = 1
< ok 0 0
> exec ^
: string
| set phase = 1
< ok 0 0
> exec ^
: string
| set placing = 1
< ok 0 0
> exec ^
: string
| set plugin = 1
< ok 0 0
> exec ^
: string
| set plugin_dir = 1
< ok 0 0
> exec ^
: string
| set plugins = 1
< ok 0 0
> exec ^
: string
| set point = 1
< ok 0 0
> exec ^
: string
| set polygon = 1
< ok 0 0
> exec ^
: string
| set port = 1
< ok 0 0
> exec ^
: string
| set position = 1
< ok 0 0
> exec ^
: string
| set precedes = 1
< ok 0 0
> exec ^
: string
| set preceding = 1
< ok 0 0
> exec ^
: string
| set precision = 1
< error
| 1064 42000 …near 'precision = 1' at line 1
> exec ^
: string
| set prefix = 1
< ok 0 0
> exec ^
: string
| set prepare = 1
< ok 0 0
> exec ^
: string
| set preserve = 1
< ok 0 0
> exec ^
: string
| set prev = 1
< ok 0 0
> exec ^
: string
| set primary = 1
< error
| 1064 42000 …near 'primary = 1' at line 1
> exec ^
: string
| set privilege_checks_user = 1
< ok 0 0
> exec ^
: string
| set privileges = 1
< ok 0 0
> exec ^
: string
| set procedure = 1
< error
| 1064 42000 …near 'procedure = 1' at line 1
> exec ^
: string
| set process = 1
< error
| 1064 42000 …near 'process = 1' at line 1
> exec ^
: string
| set processlist = 1
< ok 0 0
> exec ^
: string
| set profile = 1
< ok 0 0
> exec ^
: string
| set profiles = 1
< ok 0 0
> exec ^
: string
| set proxy = 1
< error
| 1064 42000 …near 'proxy = 1' at line 1
> exec ^
: string
| set purge = 1
< error
| 1064 42000 …near 'purge = 1' at line 1
> exec ^
: string
| set qualify = 1
< ok 0 0
> exec ^
: string
| set quarter = 1
< ok 0 0
> exec ^
: string
| set query = 1
< ok 0 0
> exec ^
: string
| set quick = 1
< ok 0 0
> exec ^
: string
| set random = 1
< ok 0 0
> exec ^
: string
| set range = 1
< error
| 1064 42000 …near 'range = 1' at line 1
> exec ^
: string
| set rank = 1
< ok 0 0
> exec ^
: string
| set read = 1
< error
| 1064 42000 …near 'read = 1' at line 1
> exec ^
: string
| set read_only = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set recover = 1
< ok 0 0
> exec ^
: string
| set recursive = 1
< ok 0 0
> exec ^
: string
| set redo_buffer_size = 1
< ok 0 0
> exec ^
: string
| set redundant = 1
< ok 0 0
> exec ^
: string
| set reference = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set relational = 1
< ok 0 0
> exec ^
: string
| set relay = 1
< ok 0 0
> exec ^
: string
| set relay_log_file = 1
< ok 0 0
> exec ^
: string
| set relay_log_pos = 1
< ok 0 0
> exec ^
: string
| set relay_thread = 1
< ok 0 0
> exec ^
: string
| set relaylog = 1
< ok 0 0
> exec ^
: string
| set release = 1
< error
| 1064 42000 …near 'release = 1' at line 1
> exec ^
: string
| set reload = 1
< error
| 1064 42000 …near 'reload = 1' at line 1
> exec ^
: string
| set remove = 1
< ok 0 0
> exec ^
: string
| set rename = 1
< error
| 1064 42000 …near 'rename = 1' at line 1
> exec ^
: string
| set reorganize = 1
< ok 0 0
> exec ^
: string
| set repair = 1
< ok 0 0
> exec ^
: string
| set repeat = 1
< error
| 1064 42000 …near 'repeat = 1' at line 1
> exec ^
: string
| set repeatable = 1
< ok 0 0
> exec ^
: string
| set replace = 1
< error
| 1064 42000 …near 'replace = 1' at line 1
> exec ^
: string
| set replica = 1
< ok 0 0
> exec ^
: string
| set replica_num = 1
< ok 0 0
> exec ^
: string
| set replicas = 1
< ok 0 0
> exec ^
: string
| set replicate_do_db = 1
< ok 0 0
> exec ^
: string
| set replicate_do_table = 1
< ok 0 0
> exec ^
: string
| set replicate_ignore_db = 1
< ok 0 0
> exec ^
: string
| set replicate_ignore_table = 1
< ok 0 0
> exec ^
: string
| set replicate_rewrite_db = 1
< ok 0 0
> exec ^
: string
| set replicate_wild_do_table = 1
< ok 0 0
> exec ^
: string
| set replicate_wild_ignore_table = 1
< ok 0 0
> exec ^
: string
| set replication = 1
< error
| 1064 42000 …near 'replication = 1' at line 1
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
< ok 0 0
> exec ^
: string
| set reset = 1
< ok 0 0
> exec ^
: string
| set resignal = 1
< error
| 1064 42000 …near 'resignal = 1' at line 1
> exec ^
: string
| set resource = 1
< error
| 1064 42000 …near 'resource = 1' at line 1
> exec ^
: string
| set respect = 1
< ok 0 0
> exec ^
: string
| set restart = 1
< ok 0 0
> exec ^
: string
| set restore = 1
< ok 0 0
> exec ^
: string
| set restrict = 1
< error
| 1064 42000 …near 'restrict = 1' at line 1
> exec ^
: string
| set resume = 1
< ok 0 0
> exec ^
: string
| set retain = 1
< ok 0 0
> exec ^
: string
| set return = 1
< error
| 1064 42000 …near 'return = 1' at line 1
> exec ^
: string
| set returned_sqlstate = 1
< ok 0 0
> exec ^
: string
| set returning = 1
< ok 0 0
> exec ^
: string
| set returns = 1
< ok 0 0
> exec ^
: string
| set reuse = 1
< ok 0 0
> exec ^
: string
| set reverse = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set rollback = 1
< ok 0 0
> exec ^
: string
| set rollup = 1
< ok 0 0
> exec ^
: string
| set rotate = 1
< ok 0 0
> exec ^
: string
| set routine = 1
< ok 0 0
> exec ^
: string
| set row = 1
< ok 0 0
> exec ^
: string
| set row_count = 1
< ok 0 0
> exec ^
: string
| set row_format = 1
< ok 0 0
> exec ^
: string
| set row_number = 1
< ok 0 0
> exec ^
: string
| set rows = 1
< ok 0 0
> exec ^
: string
| set rtree = 1
< ok 0 0
> exec ^
: string
| set s3 = 1
< ok 0 0
> exec ^
: string
| set savepoint = 1
< ok 0 0
> exec ^
: string
| set schedule = 1
< ok 0 0
> exec ^
: string
| set schema = 1
< error
| 1064 42000 …near 'schema = 1' at line 1
> exec ^
: string
| set schema_name = 1
< ok 0 0
> exec ^
: string
| set schemas = 1
< error
| 1064 42000 …near 'schemas = 1' at line 1
> exec ^
: string
| set second = 1
< ok 0 0
> exec ^
: string
| set second_microsecond = 1
< error
| 1064 42000 …near 'second_microsecond = 1' at line 1
> exec ^
: string
| set secondary = 1
< ok 0 0
> exec ^
: string
| set secondary_engine = 1
< ok 0 0
> exec ^
: string
| set secondary_engine_attribute = 1
< ok 0 0
> exec ^
: string
| set secondary_load = 1
< ok 0 0
> exec ^
: string
| set secondary_unload = 1
< ok 0 0
> exec ^
: string
| set security = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set serializable = 1
< ok 0 0
> exec ^
: string
| set server = 1
< ok 0 0
> exec ^
: string
| set session = 1
< ok 0 0
> exec ^
: string
| set session_user = 1
< ok 0 0
> exec ^
: string
| set set = 1
< error
| 1064 42000 …near 'set = 1' at line 1
> exec ^
: string
| set sets = 1
< ok 0 0
> exec ^
: string
| set share = 1
< ok 0 0
> exec ^
: string
| set show = 1
< error
| 1064 42000 …near 'show = 1' at line 1
> exec ^
: string
| set shutdown = 1
< error
| 1064 42000 …near 'shutdown = 1' at line 1
> exec ^
: string
| set signal = 1
< error
| 1064 42000 …near 'signal = 1' at line 1
> exec ^
: string
| set signed = 1
< ok 0 0
> exec ^
: string
| set simple = 1
< ok 0 0
> exec ^
: string
| set skip = 1
< error
| 1064 42000 …near 'skip = 1' at line 1
> exec ^
: string
| set slave = 1
< ok 0 0
> exec ^
: string
| set slow = 1
< ok 0 0
> exec ^
: string
| set smallint = 1
< error
| 1064 42000 …near 'smallint = 1' at line 1
> exec ^
: string
| set snapshot = 1
< ok 0 0
> exec ^
: string
| set socket = 1
< ok 0 0
> exec ^
: string
| set some = 1
< ok 0 0
> exec ^
: string
| set soname = 1
< ok 0 0
> exec ^
: string
| set sounds = 1
< ok 0 0
> exec ^
: string
| set source = 1
< ok 0 0
> exec ^
: string
| set source_auto_position = 1
< ok 0 0
> exec ^
: string
| set source_bind = 1
< ok 0 0
> exec ^
: string
| set source_compression_algorithms = 1
< ok 0 0
> exec ^
: string
| set source_connect_retry = 1
< ok 0 0
> exec ^
: string
| set source_connection_auto_failover = 1
< ok 0 0
> exec ^
: string
| set source_delay = 1
< ok 0 0
> exec ^
: string
| set source_heartbeat_period = 1
< ok 0 0
> exec ^
: string
| set source_host = 1
< ok 0 0
> exec ^
: string
| set source_log_file = 1
< ok 0 0
> exec ^
: string
| set source_log_pos = 1
< ok 0 0
> exec ^
: string
| set source_password = 1
< ok 0 0
> exec ^
: string
| set source_port = 1
< ok 0 0
> exec ^
: string
| set source_public_key_path = 1
< ok 0 0
> exec ^
: string
| set source_retry_count = 1
< ok 0 0
> exec ^
: string
| set source_ssl = 1
< ok 0 0
> exec ^
: string
| set source_ssl_ca = 1
< ok 0 0
> exec ^
: string
| set source_ssl_capath = 1
< ok 0 0
> exec ^
: string
| set source_ssl_cert = 1
< ok 0 0
> exec ^
: string
| set source_ssl_cipher = 1
< ok 0 0
> exec ^
: string
| set source_ssl_crl = 1
< ok 0 0
> exec ^
: string
| set source_ssl_crlpath = 1
< ok 0 0
> exec ^
: string
| set source_ssl_key = 1
< ok 0 0
> exec ^
: string
| set source_ssl_verify_server_cert = 1
< ok 0 0
> exec ^
: string
| set source_tls_ciphersuites = 1
< ok 0 0
> exec ^
: string
| set source_tls_version = 1
< ok 0 0
> exec ^
: string
| set source_user = 1
< ok 0 0
> exec ^
: string
| set source_zstd_compression_level = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set sql_after_mts_gaps = 1
< ok 0 0
> exec ^
: string
| set sql_before_gtids = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set sql_small_result = 1
< error
| 1064 42000 …near 'sql_small_result = 1' at line 1
> exec ^
: string
| set sql_thread = 1
< ok 0 0
> exec ^
: string
| set sql_tsi_day = 1
< ok 0 0
> exec ^
: string
| set sql_tsi_hour = 1
< ok 0 0
> exec ^
: string
| set sql_tsi_minute = 1
< ok 0 0
> exec ^
: string
| set sql_tsi_month = 1
< ok 0 0
> exec ^
: string
| set sql_tsi_quarter = 1
< ok 0 0
> exec ^
: string
| set sql_tsi_second = 1
< ok 0 0
> exec ^
: string
| set sql_tsi_week = 1
< ok 0 0
> exec ^
: string
| set sql_tsi_year = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set sse_kms = 1
< ok 0 0
> exec ^
: string
| set sse_s3 = 1
< ok 0 0
> exec ^
: string
| set ssl = 1
< error
| 1064 42000 …near 'ssl = 1' at line 1
> exec ^
: string
| set st_collect = 1
< ok 0 0
> exec ^
: string
| set stacked = 1
< ok 0 0
> exec ^
: string
| set start = 1
< ok 0 0
> exec ^
: string
| set starting = 1
< error
| 1064 42000 …near 'starting = 1' at line 1
> exec ^
: string
| set starts = 1
< ok 0 0
> exec ^
: string
| set stats_auto_recalc = 1
< ok 0 0
> exec ^
: string
| set stats_persistent = 1
< ok 0 0
> exec ^
: string
| set stats_sample_pages = 1
< ok 0 0
> exec ^
: string
| set status = 1
< ok 0 0
> exec ^
: string
| set std = 1
< ok 0 0
> exec ^
: string
| set stddev = 1
< ok 0 0
> exec ^
: string
| set stddev_pop = 1
< ok 0 0
> exec ^
: string
| set stddev_samp = 1
< ok 0 0
> exec ^
: string
| set stop = 1
< ok 0 0
> exec ^
: string
| set storage = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set strict_load = 1
< ok 0 0
> exec ^
: string
| set string = 1
< ok 0 0
> exec ^
: string
| set subclass_origin = 1
< ok 0 0
> exec ^
: string
| set subdate = 1
< ok 0 0
> exec ^
: string
| set subject = 1
< ok 0 0
> exec ^
: string
| set subpartition = 1
< ok 0 0
> exec ^
: string
| set subpartitions = 1
< ok 0 0
> exec ^
: string
| set substr = 1
< ok 0 0
> exec ^
: string
| set substring = 1
< ok 0 0
> exec ^
: string
| set sum = 1
< ok 0 0
> exec ^
: string
| set super = 1
< error
| 1064 42000 …near 'super = 1' at line 1
> exec ^
: string
| set suspend = 1
< ok 0 0
> exec ^
: string
| set swaps = 1
< ok 0 0
> exec ^
: string
| set switches = 1
< ok 0 0
> exec ^
: string
| set sysdate = 1
< ok 0 0
> exec ^
: string
| set system = 1
< ok 0 0
> exec ^
: string
| set system_user = 1
< ok 0 0
> exec ^
: string
| set table = 1
< error
| 1064 42000 …near 'table = 1' at line 1
> exec ^
: string
| set table_checksum = 1
< ok 0 0
> exec ^
: string
| set table_name = 1
< ok 0 0
> exec ^
: string
| set tables = 1
< ok 0 0
> exec ^
: string
| set tablesample = 1
< ok 0 0
> exec ^
: string
| set tablespace = 1
< ok 0 0
> exec ^
: string
| set tablet_size = 1
< ok 0 0
> exec ^
: string
| set template = 1
< ok 0 0
> exec ^
: string
| set temporary = 1
< ok 0 0
> exec ^
: string
| set temptable = 1
< ok 0 0
> exec ^
: string
| set terminated = 1
< error
| 1064 42000 …near 'terminated = 1' at line 1
> exec ^
: string
| set text = 1
< ok 0 0
> exec ^
: string
| set than = 1
< ok 0 0
> exec ^
: string
| set then = 1
< error
| 1064 42000 …near 'then = 1' at line 1
> exec ^
: string
| set thread_priority = 1
< ok 0 0
> exec ^
: string
| set ties = 1
< ok 0 0
> exec ^
: string
| set time = 1
< ok 0 0
> exec ^
: string
| set timestamp = 1
< ok 0 0
> exec ^
: string
| set timestampadd = 1
< ok 0 0
> exec ^
: string
| set timestampdiff = 1
< ok 0 0
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
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set trigger = 1
< error
| 1064 42000 …near 'trigger = 1' at line 1
> exec ^
: string
| set triggers = 1
< ok 0 0
> exec ^
: string
| set trim = 1
< ok 0 0
> exec ^
: string
| set true = 1
< error
| 1064 42000 …near 'true = 1' at line 1
> exec ^
: string
| set truncate = 1
< ok 0 0
> exec ^
: string
| set type = 1
< ok 0 0
> exec ^
: string
| set types = 1
< ok 0 0
> exec ^
: string
| set unbounded = 1
< ok 0 0
> exec ^
: string
| set uncommitted = 1
< ok 0 0
> exec ^
: string
| set undefined = 1
< ok 0 0
> exec ^
: string
| set undo = 1
< error
| 1064 42000 …near 'undo = 1' at line 1
> exec ^
: string
| set undo_buffer_size = 1
< ok 0 0
> exec ^
: string
| set undofile = 1
< ok 0 0
> exec ^
: string
| set unicode = 1
< ok 0 0
> exec ^
: string
| set uninstall = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set unlock = 1
< error
| 1064 42000 …near 'unlock = 1' at line 1
> exec ^
: string
| set unregister = 1
< ok 0 0
> exec ^
: string
| set unsigned = 1
< error
| 1064 42000 …near 'unsigned = 1' at line 1
> exec ^
: string
| set until = 1
< ok 0 0
> exec ^
: string
| set update = 1
< error
| 1064 42000 …near 'update = 1' at line 1
> exec ^
: string
| set upgrade = 1
< ok 0 0
> exec ^
: string
| set uri = 1
< ok 0 0
> exec ^
: string
| set url = 1
< ok 0 0
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
| set use_bloom_filter = 1
< ok 0 0
> exec ^
: string
| set use_frm = 1
< ok 0 0
> exec ^
: string
| set user = 1
< ok 0 0
> exec ^
: string
| set user_resources = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set validation = 1
< ok 0 0
> exec ^
: string
| set value = 1
< ok 0 0
> exec ^
: string
| set values = 1
< error
| 1064 42000 …near 'values = 1' at line 1
> exec ^
: string
| set var_pop = 1
< ok 0 0
> exec ^
: string
| set var_samp = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set variance = 1
< ok 0 0
> exec ^
: string
| set varying = 1
< error
| 1064 42000 …near 'varying = 1' at line 1
> exec ^
: string
| set vcpu = 1
< ok 0 0
> exec ^
: string
| set vector = 1
< ok 0 0
> exec ^
: string
| set verify_key_constraints = 1
< ok 0 0
> exec ^
: string
| set view = 1
< ok 0 0
> exec ^
: string
| set virtual = 1
< error
| 1064 42000 …near 'virtual = 1' at line 1
> exec ^
: string
| set visible = 1
< ok 0 0
> exec ^
: string
| set wait = 1
< ok 0 0
> exec ^
: string
| set warnings = 1
< ok 0 0
> exec ^
: string
| set week = 1
< ok 0 0
> exec ^
: string
| set weight_string = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set with = 1
< error
| 1064 42000 …near 'with = 1' at line 1
> exec ^
: string
| set without = 1
< ok 0 0
> exec ^
: string
| set work = 1
< ok 0 0
> exec ^
: string
| set wrapper = 1
< ok 0 0
> exec ^
: string
| set write = 1
< error
| 1064 42000 …near 'write = 1' at line 1
> exec ^
: string
| set x509 = 1
< ok 0 0
> exec ^
: string
| set xa = 1
< ok 0 0
> exec ^
: string
| set xid = 1
< ok 0 0
> exec ^
: string
| set xml = 1
< ok 0 0
> exec ^
: string
| set xor = 1
< error
| 1064 42000 …near 'xor = 1' at line 1
> exec ^
: string
| set year = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set global absent = 1
< ok 0 0
> exec ^
: string
| set global accessible = 1
< error
| 1064 42000 …near 'accessible = 1' at line 1
> exec ^
: string
| set global account = 1
< ok 0 0
> exec ^
: string
| set global action = 1
< ok 0 0
> exec ^
: string
| set global active = 1
< ok 0 0
> exec ^
: string
| set global add = 1
< error
| 1064 42000 …near 'add = 1' at line 1
> exec ^
: string
| set global adddate = 1
< ok 0 0
> exec ^
: string
| set global admin = 1
< ok 0 0
> exec ^
: string
| set global after = 1
< ok 0 0
> exec ^
: string
| set global against = 1
< ok 0 0
> exec ^
: string
| set global aggregate = 1
< ok 0 0
> exec ^
: string
| set global algorithm = 1
< ok 0 0
> exec ^
: string
| set global all = 1
< error
| 1064 42000 …near 'all = 1' at line 1
> exec ^
: string
| set global allow_missing_files = 1
< ok 0 0
> exec ^
: string
| set global alter = 1
< error
| 1064 42000 …near 'alter = 1' at line 1
> exec ^
: string
| set global always = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set global array = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set global asensitive = 1
< error
| 1064 42000 …near 'asensitive = 1' at line 1
> exec ^
: string
| set global assign_gtids_to_anonymous_transactions = 1
< ok 0 0
> exec ^
: string
| set global at = 1
< ok 0 0
> exec ^
: string
| set global attribute = 1
< ok 0 0
> exec ^
: string
| set global authentication = 1
< ok 0 0
> exec ^
: string
| set global auto = 1
< ok 0 0
> exec ^
: string
| set global auto_increment = 1
< ok 0 0
> exec ^
: string
| set global auto_increment_mode = 1
< ok 0 0
> exec ^
: string
| set global auto_refresh = 1
< ok 0 0
> exec ^
: string
| set global auto_refresh_source = 1
< ok 0 0
> exec ^
: string
| set global autoextend_size = 1
< ok 0 0
> exec ^
: string
| set global avg = 1
< ok 0 0
> exec ^
: string
| set global avg_row_length = 1
< ok 0 0
> exec ^
: string
| set global backup = 1
< ok 0 0
> exec ^
: string
| set global before = 1
< error
| 1064 42000 …near 'before = 1' at line 1
> exec ^
: string
| set global begin = 1
< ok 0 0
> exec ^
: string
| set global bernoulli = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set global bit = 1
< ok 0 0
> exec ^
: string
| set global bit_and = 1
< ok 0 0
> exec ^
: string
| set global bit_or = 1
< ok 0 0
> exec ^
: string
| set global bit_xor = 1
< ok 0 0
> exec ^
: string
| set global blob = 1
< error
| 1064 42000 …near 'blob = 1' at line 1
> exec ^
: string
| set global block = 1
< ok 0 0
> exec ^
: string
| set global block_size = 1
< ok 0 0
> exec ^
: string
| set global bool = 1
< ok 0 0
> exec ^
: string
| set global boolean = 1
< ok 0 0
> exec ^
: string
| set global both = 1
< error
| 1064 42000 …near 'both = 1' at line 1
> exec ^
: string
| set global btree = 1
< ok 0 0
> exec ^
: string
| set global buckets = 1
< ok 0 0
> exec ^
: string
| set global bulk = 1
< ok 0 0
> exec ^
: string
| set global by = 1
< error
| 1064 42000 …near 'by = 1' at line 1
> exec ^
: string
| set global byte = 1
< ok 0 0
> exec ^
: string
| set global cache = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set global case = 1
< error
| 1064 42000 …near 'case = 1' at line 1
> exec ^
: string
| set global cast = 1
< ok 0 0
> exec ^
: string
| set global catalog_name = 1
< ok 0 0
> exec ^
: string
| set global chain = 1
< ok 0 0
> exec ^
: string
| set global challenge_response = 1
< ok 0 0
> exec ^
: string
| set global change = 1
< error
| 1064 42000 …near 'change = 1' at line 1
> exec ^
: string
| set global changed = 1
< ok 0 0
> exec ^
: string
| set global channel = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set global check = 1
< error
| 1064 42000 …near 'check = 1' at line 1
> exec ^
: string
| set global checksum = 1
< ok 0 0
> exec ^
: string
| set global cipher = 1
< ok 0 0
> exec ^
: string
| set global class_origin = 1
< ok 0 0
> exec ^
: string
| set global client = 1
< ok 0 0
> exec ^
: string
| set global clone = 1
< ok 0 0
> exec ^
: string
| set global close = 1
< ok 0 0
> exec ^
: string
| set global coalesce = 1
< ok 0 0
> exec ^
: string
| set global code = 1
< ok 0 0
> exec ^
: string
| set global collate = 1
< error
| 1064 42000 …near 'collate = 1' at line 1
> exec ^
: string
| set global collation = 1
< ok 0 0
> exec ^
: string
| set global column = 1
< error
| 1064 42000 …near 'column = 1' at line 1
> exec ^
: string
| set global column_format = 1
< ok 0 0
> exec ^
: string
| set global column_name = 1
< ok 0 0
> exec ^
: string
| set global columns = 1
< ok 0 0
> exec ^
: string
| set global comment = 1
< ok 0 0
> exec ^
: string
| set global commit = 1
< ok 0 0
> exec ^
: string
| set global committed = 1
< ok 0 0
> exec ^
: string
| set global compact = 1
< ok 0 0
> exec ^
: string
| set global completion = 1
< ok 0 0
> exec ^
: string
| set global component = 1
< ok 0 0
> exec ^
: string
| set global compressed = 1
< ok 0 0
> exec ^
: string
| set global compression = 1
< ok 0 0
> exec ^
: string
| set global concat = 1
< ok 0 0
> exec ^
: string
| set global concurrent = 1
< ok 0 0
> exec ^
: string
| set global condition = 1
< error
| 1064 42000 …near 'condition = 1' at line 1
> exec ^
: string
| set global connection = 1
< ok 0 0
> exec ^
: string
| set global consistent = 1
< ok 0 0
> exec ^
: string
| set global constraint = 1
< error
| 1064 42000 …near 'constraint = 1' at line 1
> exec ^
: string
| set global constraint_catalog = 1
< ok 0 0
> exec ^
: string
| set global constraint_name = 1
< ok 0 0
> exec ^
: string
| set global constraint_schema = 1
< ok 0 0
> exec ^
: string
| set global contains = 1
< ok 0 0
> exec ^
: string
| set global context = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set global cpu = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set global cube = 1
< ok 0 0
> exec ^
: string
| set global cume_dist = 1
< ok 0 0
> exec ^
: string
| set global curdate = 1
< ok 0 0
> exec ^
: string
| set global current = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set global curtime = 1
< ok 0 0
> exec ^
: string
| set global data = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set global date = 1
< ok 0 0
> exec ^
: string
| set global date_add = 1
< ok 0 0
> exec ^
: string
| set global date_sub = 1
< ok 0 0
> exec ^
: string
| set global datetime = 1
< ok 0 0
> exec ^
: string
| set global day = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set global dec = 1
< ok 0 0
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
| 1064 42000 …near 'default = 1' at line 1
> exec ^
: string
| set global default_auth = 1
< ok 0 0
> exec ^
: string
| set global definer = 1
< ok 0 0
> exec ^
: string
| set global definition = 1
< ok 0 0
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
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set global deterministic = 1
< error
| 1064 42000 …near 'deterministic = 1' at line 1
> exec ^
: string
| set global diagnostics = 1
< ok 0 0
> exec ^
: string
| set global directory = 1
< ok 0 0
> exec ^
: string
| set global disable = 1
< ok 0 0
> exec ^
: string
| set global discard = 1
< ok 0 0
> exec ^
: string
| set global disk = 1
< ok 0 0
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
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set global dumpfile = 1
< ok 0 0
> exec ^
: string
| set global duplicate = 1
< ok 0 0
> exec ^
: string
| set global dynamic = 1
< ok 0 0
> exec ^
: string
| set global dynamic_partition_policy = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set global enable = 1
< ok 0 0
> exec ^
: string
| set global enable_macro_block_bloom_filter = 1
< ok 0 0
> exec ^
: string
| set global enclosed = 1
< error
| 1064 42000 …near 'enclosed = 1' at line 1
> exec ^
: string
| set global encryption = 1
< ok 0 0
> exec ^
: string
| set global end = 1
< ok 0 0
> exec ^
: string
| set global ends = 1
< ok 0 0
> exec ^
: string
| set global enforced = 1
< ok 0 0
> exec ^
: string
| set global engine = 1
< ok 0 0
> exec ^
: string
| set global engine_attribute = 1
< ok 0 0
> exec ^
: string
| set global engines = 1
< ok 0 0
> exec ^
: string
| set global enum = 1
< ok 0 0
> exec ^
: string
| set global error = 1
< ok 0 0
> exec ^
: string
| set global errors = 1
< ok 0 0
> exec ^
: string
| set global escape = 1
< ok 0 0
> exec ^
: string
| set global escaped = 1
< error
| 1064 42000 …near 'escaped = 1' at line 1
> exec ^
: string
| set global event = 1
< ok 0 0
> exec ^
: string
| set global events = 1
< ok 0 0
> exec ^
: string
| set global every = 1
< ok 0 0
> exec ^
: string
| set global except = 1
< ok 0 0
> exec ^
: string
| set global exchange = 1
< ok 0 0
> exec ^
: string
| set global exclude = 1
< ok 0 0
> exec ^
: string
| set global execute = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set global expire = 1
< ok 0 0
> exec ^
: string
| set global explain = 1
< error
| 1064 42000 …near 'explain = 1' at line 1
> exec ^
: string
| set global export = 1
< ok 0 0
> exec ^
: string
| set global extended = 1
< ok 0 0
> exec ^
: string
| set global extent_size = 1
< ok 0 0
> exec ^
: string
| set global external = 1
< ok 0 0
> exec ^
: string
| set global external_format = 1
< ok 0 0
> exec ^
: string
| set global extract = 1
< ok 0 0
> exec ^
: string
| set global factor = 1
< ok 0 0
> exec ^
: string
| set global failed_login_attempts = 1
< ok 0 0
> exec ^
: string
| set global false = 1
< error
| 1064 42000 …near 'false = 1' at line 1
> exec ^
: string
| set global fast = 1
< ok 0 0
> exec ^
: string
| set global faults = 1
< ok 0 0
> exec ^
: string
| set global fetch = 1
< error
| 1064 42000 …near 'fetch = 1' at line 1
> exec ^
: string
| set global fields = 1
< ok 0 0
> exec ^
: string
| set global file = 1
< ok 0 0
> exec ^
: string
| set global file_block_size = 1
< ok 0 0
> exec ^
: string
| set global file_format = 1
< ok 0 0
> exec ^
: string
| set global file_name = 1
< ok 0 0
> exec ^
: string
| set global file_pattern = 1
< ok 0 0
> exec ^
: string
| set global file_prefix = 1
< ok 0 0
> exec ^
: string
| set global files = 1
< ok 0 0
> exec ^
: string
| set global filter = 1
< ok 0 0
> exec ^
: string
| set global finish = 1
< ok 0 0
> exec ^
: string
| set global first = 1
< ok 0 0
> exec ^
: string
| set global first_value = 1
< ok 0 0
> exec ^
: string
| set global fixed = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set global follows = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set global found = 1
< ok 0 0
> exec ^
: string
| set global from = 1
< error
| 1064 42000 …near 'from = 1' at line 1
> exec ^
: string
| set global full = 1
< ok 0 0
> exec ^
: string
| set global fulltext = 1
< error
| 1064 42000 …near 'fulltext = 1' at line 1
> exec ^
: string
| set global function = 1
< ok 0 0
> exec ^
: string
| set global general = 1
< ok 0 0
> exec ^
: string
| set global generate = 1
< ok 0 0
> exec ^
: string
| set global generated = 1
< error
| 1064 42000 …near 'generated = 1' at line 1
> exec ^
: string
| set global geomcollection = 1
< ok 0 0
> exec ^
: string
| set global geometry = 1
< ok 0 0
> exec ^
: string
| set global geometrycollection = 1
< ok 0 0
> exec ^
: string
| set global get = 1
< error
| 1064 42000 …near 'get = 1' at line 1
> exec ^
: string
| set global get_format = 1
< ok 0 0
> exec ^
: string
| set global get_master_public_key = 1
< ok 0 0
> exec ^
: string
| set global get_source_public_key = 1
< ok 0 0
> exec ^
: string
| set global global = 1
< ok 0 0
> exec ^
: string
| set global grant = 1
< error
| 1064 42000 …near 'grant = 1' at line 1
> exec ^
: string
| set global grants = 1
< ok 0 0
> exec ^
: string
| set global greatest = 1
< ok 0 0
> exec ^
: string
| set global group = 1
< error
| 1064 42000 …near 'group = 1' at line 1
> exec ^
: string
| set global group_concat = 1
< ok 0 0
> exec ^
: string
| set global group_replication = 1
< ok 0 0
> exec ^
: string
| set global grouping = 1
< ok 0 0
> exec ^
: string
| set global groups = 1
< ok 0 0
> exec ^
: string
| set global gtid_only = 1
< ok 0 0
> exec ^
: string
| set global gtids = 1
< ok 0 0
> exec ^
: string
| set global guided = 1
< ok 0 0
> exec ^
: string
| set global handler = 1
< ok 0 0
> exec ^
: string
| set global hash = 1
< ok 0 0
> exec ^
: string
| set global having = 1
< error
| 1064 42000 …near 'having = 1' at line 1
> exec ^
: string
| set global header = 1
< ok 0 0
> exec ^
: string
| set global help = 1
< ok 0 0
> exec ^
: string
| set global high_priority = 1
< error
| 1064 42000 …near 'high_priority = 1' at line 1
> exec ^
: string
| set global histogram = 1
< ok 0 0
> exec ^
: string
| set global history = 1
< ok 0 0
> exec ^
: string
| set global host = 1
< ok 0 0
> exec ^
: string
| set global hosts = 1
< ok 0 0
> exec ^
: string
| set global hour = 1
< ok 0 0
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
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set global import = 1
< ok 0 0
> exec ^
: string
| set global in = 1
< error
| 1064 42000 …near 'in = 1' at line 1
> exec ^
: string
| set global inactive = 1
< ok 0 0
> exec ^
: string
| set global index = 1
< error
| 1064 42000 …near 'index = 1' at line 1
> exec ^
: string
| set global indexes = 1
< ok 0 0
> exec ^
: string
| set global infile = 1
< error
| 1064 42000 …near 'infile = 1' at line 1
> exec ^
: string
| set global initial = 1
< ok 0 0
> exec ^
: string
| set global initial_size = 1
< ok 0 0
> exec ^
: string
| set global initiate = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set global install = 1
< ok 0 0
> exec ^
: string
| set global instance = 1
< ok 0 0
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
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set global invoker = 1
< ok 0 0
> exec ^
: string
| set global io = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set global ipc = 1
< ok 0 0
> exec ^
: string
| set global is = 1
< error
| 1064 42000 …near 'is = 1' at line 1
> exec ^
: string
| set global isolation = 1
< ok 0 0
> exec ^
: string
| set global issuer = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set global json_arrayagg = 1
< ok 0 0
> exec ^
: string
| set global json_duality_object = 1
< ok 0 0
> exec ^
: string
| set global json_objectagg = 1
< ok 0 0
> exec ^
: string
| set global json_table = 1
< ok 0 0
> exec ^
: string
| set global json_value = 1
< ok 0 0
> exec ^
: string
| set global key = 1
< error
| 1064 42000 …near 'key = 1' at line 1
> exec ^
: string
| set global key_block_size = 1
< ok 0 0
> exec ^
: string
| set global keyring = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set global language = 1
< ok 0 0
> exec ^
: string
| set global last = 1
< ok 0 0
> exec ^
: string
| set global last_value = 1
< ok 0 0
> exec ^
: string
| set global lateral = 1
< ok 0 0
> exec ^
: string
| set global lead = 1
< ok 0 0
> exec ^
: string
| set global leading = 1
< error
| 1064 42000 …near 'leading = 1' at line 1
> exec ^
: string
| set global least = 1
< ok 0 0
> exec ^
: string
| set global leave = 1
< error
| 1064 42000 …near 'leave = 1' at line 1
> exec ^
: string
| set global leaves = 1
< ok 0 0
> exec ^
: string
| set global left = 1
< error
| 1064 42000 …near 'left = 1' at line 1
> exec ^
: string
| set global less = 1
< ok 0 0
> exec ^
: string
| set global level = 1
< ok 0 0
> exec ^
: string
| set global library = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set global list = 1
< ok 0 0
> exec ^
: string
| set global load = 1
< error
| 1064 42000 …near 'load = 1' at line 1
> exec ^
: string
| set global local = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set global locks = 1
< ok 0 0
> exec ^
: string
| set global log = 1
< ok 0 0
> exec ^
: string
| set global logfile = 1
< ok 0 0
> exec ^
: string
| set global logs = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set global manual = 1
< ok 0 0
> exec ^
: string
| set global master = 1
< ok 0 0
> exec ^
: string
| set global master_auto_position = 1
< ok 0 0
> exec ^
: string
| set global master_bind = 1
< error
| 1064 42000 …near 'master_bind = 1' at line 1
> exec ^
: string
| set global master_compression_algorithms = 1
< ok 0 0
> exec ^
: string
| set global master_connect_retry = 1
< ok 0 0
> exec ^
: string
| set global master_delay = 1
< ok 0 0
> exec ^
: string
| set global master_heartbeat_period = 1
< ok 0 0
> exec ^
: string
| set global master_host = 1
< ok 0 0
> exec ^
: string
| set global master_log_file = 1
< ok 0 0
> exec ^
: string
| set global master_log_pos = 1
< ok 0 0
> exec ^
: string
| set global master_password = 1
< ok 0 0
> exec ^
: string
| set global master_port = 1
< ok 0 0
> exec ^
: string
| set global master_public_key_path = 1
< ok 0 0
> exec ^
: string
| set global master_retry_count = 1
< ok 0 0
> exec ^
: string
| set global master_ssl = 1
< ok 0 0
> exec ^
: string
| set global master_ssl_ca = 1
< ok 0 0
> exec ^
: string
| set global master_ssl_capath = 1
< ok 0 0
> exec ^
: string
| set global master_ssl_cert = 1
< ok 0 0
> exec ^
: string
| set global master_ssl_cipher = 1
< ok 0 0
> exec ^
: string
| set global master_ssl_crl = 1
< ok 0 0
> exec ^
: string
| set global master_ssl_crlpath = 1
< ok 0 0
> exec ^
: string
| set global master_ssl_key = 1
< ok 0 0
> exec ^
: string
| set global master_ssl_verify_server_cert = 1
< error
| 1064 42000 …near 'master_ssl_verify_server_cert = 1' at line 1
> exec ^
: string
| set global master_tls_ciphersuites = 1
< ok 0 0
> exec ^
: string
| set global master_tls_version = 1
< ok 0 0
> exec ^
: string
| set global master_user = 1
< ok 0 0
> exec ^
: string
| set global master_zstd_compression_level = 1
< ok 0 0
> exec ^
: string
| set global match = 1
< error
| 1064 42000 …near 'match = 1' at line 1
> exec ^
: string
| set global materialized = 1
< ok 0 0
> exec ^
: string
| set global max = 1
< ok 0 0
> exec ^
: string
| set global max_connections_per_hour = 1
< ok 0 0
> exec ^
: string
| set global max_queries_per_hour = 1
< ok 0 0
> exec ^
: string
| set global max_rows = 1
< ok 0 0
> exec ^
: string
| set global max_size = 1
< ok 0 0
> exec ^
: string
| set global max_updates_per_hour = 1
< ok 0 0
> exec ^
: string
| set global max_user_connections = 1
< ok 0 0
> exec ^
: string
| set global maxvalue = 1
< ok 0 0
> exec ^
: string
| set global medium = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set global memory = 1
< ok 0 0
> exec ^
: string
| set global merge = 1
< ok 0 0
> exec ^
: string
| set global message_text = 1
< ok 0 0
> exec ^
: string
| set global microsecond = 1
< ok 0 0
> exec ^
: string
| set global mid = 1
< ok 0 0
> exec ^
: string
| set global middleint = 1
< error
| 1064 42000 …near 'middleint = 1' at line 1
> exec ^
: string
| set global migrate = 1
< ok 0 0
> exec ^
: string
| set global min = 1
< ok 0 0
> exec ^
: string
| set global min_rows = 1
< ok 0 0
> exec ^
: string
| set global minute = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set global modifies = 1
< error
| 1064 42000 …near 'modifies = 1' at line 1
> exec ^
: string
| set global modify = 1
< ok 0 0
> exec ^
: string
| set global month = 1
< ok 0 0
> exec ^
: string
| set global multilinestring = 1
< ok 0 0
> exec ^
: string
| set global multipoint = 1
< ok 0 0
> exec ^
: string
| set global multipolygon = 1
< ok 0 0
> exec ^
: string
| set global mutex = 1
< ok 0 0
> exec ^
: string
| set global mysql_errno = 1
< ok 0 0
> exec ^
: string
| set global name = 1
< ok 0 0
> exec ^
: string
| set global names = 1
< ok 0 0
> exec ^
: string
| set global national = 1
< ok 0 0
> exec ^
: string
| set global natural = 1
< error
| 1064 42000 …near 'natural = 1' at line 1
> exec ^
: string
| set global nchar = 1
< ok 0 0
> exec ^
: string
| set global ndb = 1
< ok 0 0
> exec ^
: string
| set global ndbcluster = 1
< ok 0 0
> exec ^
: string
| set global nested = 1
< ok 0 0
> exec ^
: string
| set global network_namespace = 1
< ok 0 0
> exec ^
: string
| set global never = 1
< ok 0 0
> exec ^
: string
| set global new = 1
< ok 0 0
> exec ^
: string
| set global next = 1
< ok 0 0
> exec ^
: string
| set global no = 1
< ok 0 0
> exec ^
: string
| set global no_wait = 1
< ok 0 0
> exec ^
: string
| set global no_write_to_binlog = 1
< error
| 1064 42000 …near 'no_write_to_binlog = 1' at line 1
> exec ^
: string
| set global nodegroup = 1
< ok 0 0
> exec ^
: string
| set global none = 1
< ok 0 0
> exec ^
: string
| set global not = 1
< error
| 1064 42000 …near 'not = 1' at line 1
> exec ^
: string
| set global now = 1
< ok 0 0
> exec ^
: string
| set global nowait = 1
< ok 0 0
> exec ^
: string
| set global nth_value = 1
< ok 0 0
> exec ^
: string
| set global ntile = 1
< ok 0 0
> exec ^
: string
| set global null = 1
< error
| 1064 42000 …near 'null = 1' at line 1
> exec ^
: string
| set global nullif = 1
< ok 0 0
> exec ^
: string
| set global nulls = 1
< ok 0 0
> exec ^
: string
| set global number = 1
< ok 0 0
> exec ^
: string
| set global numeric = 1
< error
| 1064 42000 …near 'numeric = 1' at line 1
> exec ^
: string
| set global nvarchar = 1
< ok 0 0
> exec ^
: string
| set global of = 1
< ok 0 0
> exec ^
: string
| set global off = 1
< ok 0 0
> exec ^
: string
| set global offset = 1
< ok 0 0
> exec ^
: string
| set global oj = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set global only = 1
< ok 0 0
> exec ^
: string
| set global open = 1
< ok 0 0
> exec ^
: string
| set global optimize = 1
< error
| 1064 42000 …near 'optimize = 1' at line 1
> exec ^
: string
| set global optimizer_costs = 1
< ok 0 0
> exec ^
: string
| set global option = 1
< error
| 1064 42000 …near 'option = 1' at line 1
> exec ^
: string
| set global optional = 1
< ok 0 0
> exec ^
: string
| set global optionally = 1
< error
| 1064 42000 …near 'optionally = 1' at line 1
> exec ^
: string
| set global options = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set global organization = 1
< ok 0 0
> exec ^
: string
| set global others = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set global overlay = 1
< ok 0 0
> exec ^
: string
| set global overwrite = 1
< ok 0 0
> exec ^
: string
| set global owner = 1
< ok 0 0
> exec ^
: string
| set global pack_keys = 1
< ok 0 0
> exec ^
: string
| set global page = 1
< ok 0 0
> exec ^
: string
| set global parallel = 1
< ok 0 0
> exec ^
: string
| set global parameters = 1
< ok 0 0
> exec ^
: string
| set global parse_tree = 1
< ok 0 0
> exec ^
: string
| set global parser = 1
< ok 0 0
> exec ^
: string
| set global partial = 1
< ok 0 0
> exec ^
: string
| set global partition = 1
< error
| 1064 42000 …near 'partition = 1' at line 1
> exec ^
: string
| set global partitioning = 1
< ok 0 0
> exec ^
: string
| set global partitions = 1
< ok 0 0
> exec ^
: string
| set global password = 1
< ok 0 0
> exec ^
: string
| set global password_lock_time = 1
< ok 0 0
> exec ^
: string
| set global path = 1
< ok 0 0
> exec ^
: string
| set global pctfree = 1
< ok 0 0
> exec ^
: string
| set global percent_rank = 1
< ok 0 0
> exec ^
: string
| set global persist = 1
< ok 0 0
> exec ^
: string
| set global persist_only = 1
< ok 0 0
> exec ^
: string
| set global phase = 1
< ok 0 0
> exec ^
: string
| set global placing = 1
< ok 0 0
> exec ^
: string
| set global plugin = 1
< ok 0 0
> exec ^
: string
| set global plugin_dir = 1
< ok 0 0
> exec ^
: string
| set global plugins = 1
< ok 0 0
> exec ^
: string
| set global point = 1
< ok 0 0
> exec ^
: string
| set global polygon = 1
< ok 0 0
> exec ^
: string
| set global port = 1
< ok 0 0
> exec ^
: string
| set global position = 1
< ok 0 0
> exec ^
: string
| set global precedes = 1
< ok 0 0
> exec ^
: string
| set global preceding = 1
< ok 0 0
> exec ^
: string
| set global precision = 1
< error
| 1064 42000 …near 'precision = 1' at line 1
> exec ^
: string
| set global prefix = 1
< ok 0 0
> exec ^
: string
| set global prepare = 1
< ok 0 0
> exec ^
: string
| set global preserve = 1
< ok 0 0
> exec ^
: string
| set global prev = 1
< ok 0 0
> exec ^
: string
| set global primary = 1
< error
| 1064 42000 …near 'primary = 1' at line 1
> exec ^
: string
| set global privilege_checks_user = 1
< ok 0 0
> exec ^
: string
| set global privileges = 1
< ok 0 0
> exec ^
: string
| set global procedure = 1
< error
| 1064 42000 …near 'procedure = 1' at line 1
> exec ^
: string
| set global process = 1
< ok 0 0
> exec ^
: string
| set global processlist = 1
< ok 0 0
> exec ^
: string
| set global profile = 1
< ok 0 0
> exec ^
: string
| set global profiles = 1
< ok 0 0
> exec ^
: string
| set global proxy = 1
< ok 0 0
> exec ^
: string
| set global purge = 1
< error
| 1064 42000 …near 'purge = 1' at line 1
> exec ^
: string
| set global qualify = 1
< ok 0 0
> exec ^
: string
| set global quarter = 1
< ok 0 0
> exec ^
: string
| set global query = 1
< ok 0 0
> exec ^
: string
| set global quick = 1
< ok 0 0
> exec ^
: string
| set global random = 1
< ok 0 0
> exec ^
: string
| set global range = 1
< ok 0 0
> exec ^
: string
| set global rank = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set global recover = 1
< ok 0 0
> exec ^
: string
| set global recursive = 1
< ok 0 0
> exec ^
: string
| set global redo_buffer_size = 1
< ok 0 0
> exec ^
: string
| set global redundant = 1
< ok 0 0
> exec ^
: string
| set global reference = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set global relational = 1
< ok 0 0
> exec ^
: string
| set global relay = 1
< ok 0 0
> exec ^
: string
| set global relay_log_file = 1
< ok 0 0
> exec ^
: string
| set global relay_log_pos = 1
< ok 0 0
> exec ^
: string
| set global relay_thread = 1
< ok 0 0
> exec ^
: string
| set global relaylog = 1
< ok 0 0
> exec ^
: string
| set global release = 1
< error
| 1064 42000 …near 'release = 1' at line 1
> exec ^
: string
| set global reload = 1
< ok 0 0
> exec ^
: string
| set global remove = 1
< ok 0 0
> exec ^
: string
| set global rename = 1
< error
| 1064 42000 …near 'rename = 1' at line 1
> exec ^
: string
| set global reorganize = 1
< ok 0 0
> exec ^
: string
| set global repair = 1
< ok 0 0
> exec ^
: string
| set global repeat = 1
< error
| 1064 42000 …near 'repeat = 1' at line 1
> exec ^
: string
| set global repeatable = 1
< ok 0 0
> exec ^
: string
| set global replace = 1
< error
| 1064 42000 …near 'replace = 1' at line 1
> exec ^
: string
| set global replica = 1
< ok 0 0
> exec ^
: string
| set global replica_num = 1
< ok 0 0
> exec ^
: string
| set global replicas = 1
< ok 0 0
> exec ^
: string
| set global replicate_do_db = 1
< ok 0 0
> exec ^
: string
| set global replicate_do_table = 1
< ok 0 0
> exec ^
: string
| set global replicate_ignore_db = 1
< ok 0 0
> exec ^
: string
| set global replicate_ignore_table = 1
< ok 0 0
> exec ^
: string
| set global replicate_rewrite_db = 1
< ok 0 0
> exec ^
: string
| set global replicate_wild_do_table = 1
< ok 0 0
> exec ^
: string
| set global replicate_wild_ignore_table = 1
< ok 0 0
> exec ^
: string
| set global replication = 1
< ok 0 0
> exec ^
: string
| set global require = 1
< error
| 1064 42000 …near 'require = 1' at line 1
> exec ^
: string
| set global require_row_format = 1
< ok 0 0
> exec ^
: string
| set global require_table_primary_key_check = 1
< ok 0 0
> exec ^
: string
| set global reset = 1
< ok 0 0
> exec ^
: string
| set global resignal = 1
< error
| 1064 42000 …near 'resignal = 1' at line 1
> exec ^
: string
| set global resource = 1
< ok 0 0
> exec ^
: string
| set global respect = 1
< ok 0 0
> exec ^
: string
| set global restart = 1
< ok 0 0
> exec ^
: string
| set global restore = 1
< ok 0 0
> exec ^
: string
| set global restrict = 1
< error
| 1064 42000 …near 'restrict = 1' at line 1
> exec ^
: string
| set global resume = 1
< ok 0 0
> exec ^
: string
| set global retain = 1
< ok 0 0
> exec ^
: string
| set global return = 1
< error
| 1064 42000 …near 'return = 1' at line 1
> exec ^
: string
| set global returned_sqlstate = 1
< ok 0 0
> exec ^
: string
| set global returning = 1
< ok 0 0
> exec ^
: string
| set global returns = 1
< ok 0 0
> exec ^
: string
| set global reuse = 1
< ok 0 0
> exec ^
: string
| set global reverse = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set global rollback = 1
< ok 0 0
> exec ^
: string
| set global rollup = 1
< ok 0 0
> exec ^
: string
| set global rotate = 1
< ok 0 0
> exec ^
: string
| set global routine = 1
< ok 0 0
> exec ^
: string
| set global row = 1
< ok 0 0
> exec ^
: string
| set global row_count = 1
< ok 0 0
> exec ^
: string
| set global row_format = 1
< ok 0 0
> exec ^
: string
| set global row_number = 1
< ok 0 0
> exec ^
: string
| set global rows = 1
< ok 0 0
> exec ^
: string
| set global rtree = 1
< ok 0 0
> exec ^
: string
| set global s3 = 1
< ok 0 0
> exec ^
: string
| set global savepoint = 1
< ok 0 0
> exec ^
: string
| set global schedule = 1
< ok 0 0
> exec ^
: string
| set global schema = 1
< error
| 1064 42000 …near 'schema = 1' at line 1
> exec ^
: string
| set global schema_name = 1
< ok 0 0
> exec ^
: string
| set global schemas = 1
< error
| 1064 42000 …near 'schemas = 1' at line 1
> exec ^
: string
| set global second = 1
< ok 0 0
> exec ^
: string
| set global second_microsecond = 1
< error
| 1064 42000 …near 'second_microsecond = 1' at line 1
> exec ^
: string
| set global secondary = 1
< ok 0 0
> exec ^
: string
| set global secondary_engine = 1
< ok 0 0
> exec ^
: string
| set global secondary_engine_attribute = 1
< ok 0 0
> exec ^
: string
| set global secondary_load = 1
< ok 0 0
> exec ^
: string
| set global secondary_unload = 1
< ok 0 0
> exec ^
: string
| set global security = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set global serializable = 1
< ok 0 0
> exec ^
: string
| set global server = 1
< ok 0 0
> exec ^
: string
| set global session = 1
< ok 0 0
> exec ^
: string
| set global session_user = 1
< ok 0 0
> exec ^
: string
| set global set = 1
< error
| 1064 42000 …near 'set = 1' at line 1
> exec ^
: string
| set global sets = 1
< ok 0 0
> exec ^
: string
| set global share = 1
< ok 0 0
> exec ^
: string
| set global show = 1
< error
| 1064 42000 …near 'show = 1' at line 1
> exec ^
: string
| set global shutdown = 1
< ok 0 0
> exec ^
: string
| set global signal = 1
< error
| 1064 42000 …near 'signal = 1' at line 1
> exec ^
: string
| set global signed = 1
< ok 0 0
> exec ^
: string
| set global simple = 1
< ok 0 0
> exec ^
: string
| set global skip = 1
< error
| 1064 42000 …near 'skip = 1' at line 1
> exec ^
: string
| set global slave = 1
< ok 0 0
> exec ^
: string
| set global slow = 1
< ok 0 0
> exec ^
: string
| set global smallint = 1
< error
| 1064 42000 …near 'smallint = 1' at line 1
> exec ^
: string
| set global snapshot = 1
< ok 0 0
> exec ^
: string
| set global socket = 1
< ok 0 0
> exec ^
: string
| set global some = 1
< ok 0 0
> exec ^
: string
| set global soname = 1
< ok 0 0
> exec ^
: string
| set global sounds = 1
< ok 0 0
> exec ^
: string
| set global source = 1
< ok 0 0
> exec ^
: string
| set global source_auto_position = 1
< ok 0 0
> exec ^
: string
| set global source_bind = 1
< ok 0 0
> exec ^
: string
| set global source_compression_algorithms = 1
< ok 0 0
> exec ^
: string
| set global source_connect_retry = 1
< ok 0 0
> exec ^
: string
| set global source_connection_auto_failover = 1
< ok 0 0
> exec ^
: string
| set global source_delay = 1
< ok 0 0
> exec ^
: string
| set global source_heartbeat_period = 1
< ok 0 0
> exec ^
: string
| set global source_host = 1
< ok 0 0
> exec ^
: string
| set global source_log_file = 1
< ok 0 0
> exec ^
: string
| set global source_log_pos = 1
< ok 0 0
> exec ^
: string
| set global source_password = 1
< ok 0 0
> exec ^
: string
| set global source_port = 1
< ok 0 0
> exec ^
: string
| set global source_public_key_path = 1
< ok 0 0
> exec ^
: string
| set global source_retry_count = 1
< ok 0 0
> exec ^
: string
| set global source_ssl = 1
< ok 0 0
> exec ^
: string
| set global source_ssl_ca = 1
< ok 0 0
> exec ^
: string
| set global source_ssl_capath = 1
< ok 0 0
> exec ^
: string
| set global source_ssl_cert = 1
< ok 0 0
> exec ^
: string
| set global source_ssl_cipher = 1
< ok 0 0
> exec ^
: string
| set global source_ssl_crl = 1
< ok 0 0
> exec ^
: string
| set global source_ssl_crlpath = 1
< ok 0 0
> exec ^
: string
| set global source_ssl_key = 1
< ok 0 0
> exec ^
: string
| set global source_ssl_verify_server_cert = 1
< ok 0 0
> exec ^
: string
| set global source_tls_ciphersuites = 1
< ok 0 0
> exec ^
: string
| set global source_tls_version = 1
< ok 0 0
> exec ^
: string
| set global source_user = 1
< ok 0 0
> exec ^
: string
| set global source_zstd_compression_level = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set global sql_after_mts_gaps = 1
< ok 0 0
> exec ^
: string
| set global sql_before_gtids = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set global sql_small_result = 1
< error
| 1064 42000 …near 'sql_small_result = 1' at line 1
> exec ^
: string
| set global sql_thread = 1
< ok 0 0
> exec ^
: string
| set global sql_tsi_day = 1
< ok 0 0
> exec ^
: string
| set global sql_tsi_hour = 1
< ok 0 0
> exec ^
: string
| set global sql_tsi_minute = 1
< ok 0 0
> exec ^
: string
| set global sql_tsi_month = 1
< ok 0 0
> exec ^
: string
| set global sql_tsi_quarter = 1
< ok 0 0
> exec ^
: string
| set global sql_tsi_second = 1
< ok 0 0
> exec ^
: string
| set global sql_tsi_week = 1
< ok 0 0
> exec ^
: string
| set global sql_tsi_year = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set global sse_kms = 1
< ok 0 0
> exec ^
: string
| set global sse_s3 = 1
< ok 0 0
> exec ^
: string
| set global ssl = 1
< error
| 1064 42000 …near 'ssl = 1' at line 1
> exec ^
: string
| set global st_collect = 1
< ok 0 0
> exec ^
: string
| set global stacked = 1
< ok 0 0
> exec ^
: string
| set global start = 1
< ok 0 0
> exec ^
: string
| set global starting = 1
< error
| 1064 42000 …near 'starting = 1' at line 1
> exec ^
: string
| set global starts = 1
< ok 0 0
> exec ^
: string
| set global stats_auto_recalc = 1
< ok 0 0
> exec ^
: string
| set global stats_persistent = 1
< ok 0 0
> exec ^
: string
| set global stats_sample_pages = 1
< ok 0 0
> exec ^
: string
| set global status = 1
< ok 0 0
> exec ^
: string
| set global std = 1
< ok 0 0
> exec ^
: string
| set global stddev = 1
< ok 0 0
> exec ^
: string
| set global stddev_pop = 1
< ok 0 0
> exec ^
: string
| set global stddev_samp = 1
< ok 0 0
> exec ^
: string
| set global stop = 1
< ok 0 0
> exec ^
: string
| set global storage = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set global strict_load = 1
< ok 0 0
> exec ^
: string
| set global string = 1
< ok 0 0
> exec ^
: string
| set global subclass_origin = 1
< ok 0 0
> exec ^
: string
| set global subdate = 1
< ok 0 0
> exec ^
: string
| set global subject = 1
< ok 0 0
> exec ^
: string
| set global subpartition = 1
< ok 0 0
> exec ^
: string
| set global subpartitions = 1
< ok 0 0
> exec ^
: string
| set global substr = 1
< ok 0 0
> exec ^
: string
| set global substring = 1
< ok 0 0
> exec ^
: string
| set global sum = 1
< ok 0 0
> exec ^
: string
| set global super = 1
< ok 0 0
> exec ^
: string
| set global suspend = 1
< ok 0 0
> exec ^
: string
| set global swaps = 1
< ok 0 0
> exec ^
: string
| set global switches = 1
< ok 0 0
> exec ^
: string
| set global sysdate = 1
< ok 0 0
> exec ^
: string
| set global system = 1
< ok 0 0
> exec ^
: string
| set global system_user = 1
< ok 0 0
> exec ^
: string
| set global table = 1
< error
| 1064 42000 …near 'table = 1' at line 1
> exec ^
: string
| set global table_checksum = 1
< ok 0 0
> exec ^
: string
| set global table_name = 1
< ok 0 0
> exec ^
: string
| set global tables = 1
< ok 0 0
> exec ^
: string
| set global tablesample = 1
< ok 0 0
> exec ^
: string
| set global tablespace = 1
< ok 0 0
> exec ^
: string
| set global tablet_size = 1
< ok 0 0
> exec ^
: string
| set global template = 1
< ok 0 0
> exec ^
: string
| set global temporary = 1
< ok 0 0
> exec ^
: string
| set global temptable = 1
< ok 0 0
> exec ^
: string
| set global terminated = 1
< error
| 1064 42000 …near 'terminated = 1' at line 1
> exec ^
: string
| set global text = 1
< ok 0 0
> exec ^
: string
| set global than = 1
< ok 0 0
> exec ^
: string
| set global then = 1
< error
| 1064 42000 …near 'then = 1' at line 1
> exec ^
: string
| set global thread_priority = 1
< ok 0 0
> exec ^
: string
| set global ties = 1
< ok 0 0
> exec ^
: string
| set global time = 1
< ok 0 0
> exec ^
: string
| set global timestamp = 1
< ok 0 0
> exec ^
: string
| set global timestampadd = 1
< ok 0 0
> exec ^
: string
| set global timestampdiff = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set global to = 1
< error
| 1064 42000 …near '= 1' at line 1
> exec ^
: string
| set global trailing = 1
< error
| 1064 42000 …near 'trailing = 1' at line 1
> exec ^
: string
| set global transaction = 1
< error
| 1064 42000 …near '= 1' at line 1
> exec ^
: string
| set global trigger = 1
< error
| 1064 42000 …near 'trigger = 1' at line 1
> exec ^
: string
| set global triggers = 1
< ok 0 0
> exec ^
: string
| set global trim = 1
< ok 0 0
> exec ^
: string
| set global true = 1
< error
| 1064 42000 …near 'true = 1' at line 1
> exec ^
: string
| set global truncate = 1
< ok 0 0
> exec ^
: string
| set global type = 1
< ok 0 0
> exec ^
: string
| set global types = 1
< ok 0 0
> exec ^
: string
| set global unbounded = 1
< ok 0 0
> exec ^
: string
| set global uncommitted = 1
< ok 0 0
> exec ^
: string
| set global undefined = 1
< ok 0 0
> exec ^
: string
| set global undo = 1
< error
| 1064 42000 …near 'undo = 1' at line 1
> exec ^
: string
| set global undo_buffer_size = 1
< ok 0 0
> exec ^
: string
| set global undofile = 1
< ok 0 0
> exec ^
: string
| set global unicode = 1
< ok 0 0
> exec ^
: string
| set global uninstall = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set global unlock = 1
< error
| 1064 42000 …near 'unlock = 1' at line 1
> exec ^
: string
| set global unregister = 1
< ok 0 0
> exec ^
: string
| set global unsigned = 1
< error
| 1064 42000 …near 'unsigned = 1' at line 1
> exec ^
: string
| set global until = 1
< ok 0 0
> exec ^
: string
| set global update = 1
< error
| 1064 42000 …near 'update = 1' at line 1
> exec ^
: string
| set global upgrade = 1
< ok 0 0
> exec ^
: string
| set global uri = 1
< ok 0 0
> exec ^
: string
| set global url = 1
< ok 0 0
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
| set global use_bloom_filter = 1
< ok 0 0
> exec ^
: string
| set global use_frm = 1
< ok 0 0
> exec ^
: string
| set global user = 1
< ok 0 0
> exec ^
: string
| set global user_resources = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set global validation = 1
< ok 0 0
> exec ^
: string
| set global value = 1
< ok 0 0
> exec ^
: string
| set global values = 1
< error
| 1064 42000 …near 'values = 1' at line 1
> exec ^
: string
| set global var_pop = 1
< ok 0 0
> exec ^
: string
| set global var_samp = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set global variance = 1
< ok 0 0
> exec ^
: string
| set global varying = 1
< error
| 1064 42000 …near 'varying = 1' at line 1
> exec ^
: string
| set global vcpu = 1
< ok 0 0
> exec ^
: string
| set global vector = 1
< ok 0 0
> exec ^
: string
| set global verify_key_constraints = 1
< ok 0 0
> exec ^
: string
| set global view = 1
< ok 0 0
> exec ^
: string
| set global virtual = 1
< error
| 1064 42000 …near 'virtual = 1' at line 1
> exec ^
: string
| set global visible = 1
< ok 0 0
> exec ^
: string
| set global wait = 1
< ok 0 0
> exec ^
: string
| set global warnings = 1
< ok 0 0
> exec ^
: string
| set global week = 1
< ok 0 0
> exec ^
: string
| set global weight_string = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set global with = 1
< error
| 1064 42000 …near 'with = 1' at line 1
> exec ^
: string
| set global without = 1
< ok 0 0
> exec ^
: string
| set global work = 1
< ok 0 0
> exec ^
: string
| set global wrapper = 1
< ok 0 0
> exec ^
: string
| set global write = 1
< error
| 1064 42000 …near 'write = 1' at line 1
> exec ^
: string
| set global x509 = 1
< ok 0 0
> exec ^
: string
| set global xa = 1
< ok 0 0
> exec ^
: string
| set global xid = 1
< ok 0 0
> exec ^
: string
| set global xml = 1
< ok 0 0
> exec ^
: string
| set global xor = 1
< error
| 1064 42000 …near 'xor = 1' at line 1
> exec ^
: string
| set global year = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set absent.a = 1
< error
| 1064 42000 …near 'absent.a = 1' at line 1
> exec ^
: string
| set accessible.a = 1
< error
| 1064 42000 …near 'accessible.a = 1' at line 1
> exec ^
: string
| set account.a = 1
< error
| 1064 42000 …near 'account.a = 1' at line 1
> exec ^
: string
| set action.a = 1
< error
| 1064 42000 …near 'action.a = 1' at line 1
> exec ^
: string
| set active.a = 1
< error
| 1064 42000 …near 'active.a = 1' at line 1
> exec ^
: string
| set add.a = 1
< error
| 1064 42000 …near 'add.a = 1' at line 1
> exec ^
: string
| set adddate.a = 1
< error
| 1064 42000 …near 'adddate.a = 1' at line 1
> exec ^
: string
| set admin.a = 1
< error
| 1064 42000 …near 'admin.a = 1' at line 1
> exec ^
: string
| set after.a = 1
< error
| 1064 42000 …near 'after.a = 1' at line 1
> exec ^
: string
| set against.a = 1
< error
| 1064 42000 …near 'against.a = 1' at line 1
> exec ^
: string
| set aggregate.a = 1
< error
| 1064 42000 …near 'aggregate.a = 1' at line 1
> exec ^
: string
| set algorithm.a = 1
< error
| 1064 42000 …near 'algorithm.a = 1' at line 1
> exec ^
: string
| set all.a = 1
< error
| 1064 42000 …near 'all.a = 1' at line 1
> exec ^
: string
| set allow_missing_files.a = 1
< error
| 1064 42000 …near 'allow_missing_files.a = 1' at line 1
> exec ^
: string
| set alter.a = 1
< error
| 1064 42000 …near 'alter.a = 1' at line 1
> exec ^
: string
| set always.a = 1
< error
| 1064 42000 …near 'always.a = 1' at line 1
> exec ^
: string
| set analyze.a = 1
< error
| 1064 42000 …near 'analyze.a = 1' at line 1
> exec ^
: string
| set and.a = 1
< error
| 1064 42000 …near 'and.a = 1' at line 1
> exec ^
: string
| set any.a = 1
< error
| 1064 42000 …near 'any.a = 1' at line 1
> exec ^
: string
| set array.a = 1
< error
| 1064 42000 …near 'array.a = 1' at line 1
> exec ^
: string
| set as.a = 1
< error
| 1064 42000 …near 'as.a = 1' at line 1
> exec ^
: string
| set asc.a = 1
< error
| 1064 42000 …near 'asc.a = 1' at line 1
> exec ^
: string
| set ascii.a = 1
< error
| 1064 42000 …near 'ascii.a = 1' at line 1
> exec ^
: string
| set asensitive.a = 1
< error
| 1064 42000 …near 'asensitive.a = 1' at line 1
> exec ^
: string
| set assign_gtids_to_anonymous_transactions.a = 1
< error
| 1064 42000 …near 'assign_gtids_to_anonymous_transactions.a = 1' at line 1
> exec ^
: string
| set at.a = 1
< error
| 1064 42000 …near 'at.a = 1' at line 1
> exec ^
: string
| set attribute.a = 1
< error
| 1064 42000 …near 'attribute.a = 1' at line 1
> exec ^
: string
| set authentication.a = 1
< error
| 1064 42000 …near 'authentication.a = 1' at line 1
> exec ^
: string
| set auto.a = 1
< error
| 1064 42000 …near 'auto.a = 1' at line 1
> exec ^
: string
| set auto_increment.a = 1
< error
| 1064 42000 …near 'auto_increment.a = 1' at line 1
> exec ^
: string
| set auto_increment_mode.a = 1
< error
| 1064 42000 …near 'auto_increment_mode.a = 1' at line 1
> exec ^
: string
| set auto_refresh.a = 1
< error
| 1064 42000 …near 'auto_refresh.a = 1' at line 1
> exec ^
: string
| set auto_refresh_source.a = 1
< error
| 1064 42000 …near 'auto_refresh_source.a = 1' at line 1
> exec ^
: string
| set autoextend_size.a = 1
< error
| 1064 42000 …near 'autoextend_size.a = 1' at line 1
> exec ^
: string
| set avg.a = 1
< error
| 1064 42000 …near 'avg.a = 1' at line 1
> exec ^
: string
| set avg_row_length.a = 1
< error
| 1064 42000 …near 'avg_row_length.a = 1' at line 1
> exec ^
: string
| set backup.a = 1
< error
| 1064 42000 …near 'backup.a = 1' at line 1
> exec ^
: string
| set before.a = 1
< error
| 1064 42000 …near 'before.a = 1' at line 1
> exec ^
: string
| set begin.a = 1
< error
| 1064 42000 …near 'begin.a = 1' at line 1
> exec ^
: string
| set bernoulli.a = 1
< error
| 1064 42000 …near 'bernoulli.a = 1' at line 1
> exec ^
: string
| set between.a = 1
< error
| 1064 42000 …near 'between.a = 1' at line 1
> exec ^
: string
| set bigint.a = 1
< error
| 1064 42000 …near 'bigint.a = 1' at line 1
> exec ^
: string
| set binary.a = 1
< error
| 1064 42000 …near 'binary.a = 1' at line 1
> exec ^
: string
| set binlog.a = 1
< error
| 1064 42000 …near 'binlog.a = 1' at line 1
> exec ^
: string
| set bit.a = 1
< error
| 1064 42000 …near 'bit.a = 1' at line 1
> exec ^
: string
| set bit_and.a = 1
< error
| 1064 42000 …near 'bit_and.a = 1' at line 1
> exec ^
: string
| set bit_or.a = 1
< error
| 1064 42000 …near 'bit_or.a = 1' at line 1
> exec ^
: string
| set bit_xor.a = 1
< error
| 1064 42000 …near 'bit_xor.a = 1' at line 1
> exec ^
: string
| set blob.a = 1
< error
| 1064 42000 …near 'blob.a = 1' at line 1
> exec ^
: string
| set block.a = 1
< error
| 1064 42000 …near 'block.a = 1' at line 1
> exec ^
: string
| set block_size.a = 1
< error
| 1064 42000 …near 'block_size.a = 1' at line 1
> exec ^
: string
| set bool.a = 1
< error
| 1064 42000 …near 'bool.a = 1' at line 1
> exec ^
: string
| set boolean.a = 1
< error
| 1064 42000 …near 'boolean.a = 1' at line 1
> exec ^
: string
| set both.a = 1
< error
| 1064 42000 …near 'both.a = 1' at line 1
> exec ^
: string
| set btree.a = 1
< error
| 1064 42000 …near 'btree.a = 1' at line 1
> exec ^
: string
| set buckets.a = 1
< error
| 1064 42000 …near 'buckets.a = 1' at line 1
> exec ^
: string
| set bulk.a = 1
< error
| 1064 42000 …near 'bulk.a = 1' at line 1
> exec ^
: string
| set by.a = 1
< error
| 1064 42000 …near 'by.a = 1' at line 1
> exec ^
: string
| set byte.a = 1
< error
| 1064 42000 …near 'byte.a = 1' at line 1
> exec ^
: string
| set cache.a = 1
< error
| 1064 42000 …near 'cache.a = 1' at line 1
> exec ^
: string
| set call.a = 1
< error
| 1064 42000 …near 'call.a = 1' at line 1
> exec ^
: string
| set cascade.a = 1
< error
| 1064 42000 …near 'cascade.a = 1' at line 1
> exec ^
: string
| set cascaded.a = 1
< error
| 1064 42000 …near 'cascaded.a = 1' at line 1
> exec ^
: string
| set case.a = 1
< error
| 1064 42000 …near 'case.a = 1' at line 1
> exec ^
: string
| set cast.a = 1
< error
| 1064 42000 …near 'cast.a = 1' at line 1
> exec ^
: string
| set catalog_name.a = 1
< error
| 1064 42000 …near 'catalog_name.a = 1' at line 1
> exec ^
: string
| set chain.a = 1
< error
| 1064 42000 …near 'chain.a = 1' at line 1
> exec ^
: string
| set challenge_response.a = 1
< error
| 1064 42000 …near 'challenge_response.a = 1' at line 1
> exec ^
: string
| set change.a = 1
< error
| 1064 42000 …near 'change.a = 1' at line 1
> exec ^
: string
| set changed.a = 1
< error
| 1064 42000 …near 'changed.a = 1' at line 1
> exec ^
: string
| set channel.a = 1
< error
| 1064 42000 …near 'channel.a = 1' at line 1
> exec ^
: string
| set char.a = 1
< error
| 1064 42000 …near 'char.a = 1' at line 1
> exec ^
: string
| set character.a = 1
< error
| 1064 42000 …near 'character.a = 1' at line 1
> exec ^
: string
| set charset.a = 1
< error
| 1064 42000 …near 'charset.a = 1' at line 1
> exec ^
: string
| set check.a = 1
< error
| 1064 42000 …near 'check.a = 1' at line 1
> exec ^
: string
| set checksum.a = 1
< error
| 1064 42000 …near 'checksum.a = 1' at line 1
> exec ^
: string
| set cipher.a = 1
< error
| 1064 42000 …near 'cipher.a = 1' at line 1
> exec ^
: string
| set class_origin.a = 1
< error
| 1064 42000 …near 'class_origin.a = 1' at line 1
> exec ^
: string
| set client.a = 1
< error
| 1064 42000 …near 'client.a = 1' at line 1
> exec ^
: string
| set clone.a = 1
< error
| 1064 42000 …near 'clone.a = 1' at line 1
> exec ^
: string
| set close.a = 1
< error
| 1064 42000 …near 'close.a = 1' at line 1
> exec ^
: string
| set coalesce.a = 1
< error
| 1064 42000 …near 'coalesce.a = 1' at line 1
> exec ^
: string
| set code.a = 1
< error
| 1064 42000 …near 'code.a = 1' at line 1
> exec ^
: string
| set collate.a = 1
< error
| 1064 42000 …near 'collate.a = 1' at line 1
> exec ^
: string
| set collation.a = 1
< error
| 1064 42000 …near 'collation.a = 1' at line 1
> exec ^
: string
| set column.a = 1
< error
| 1064 42000 …near 'column.a = 1' at line 1
> exec ^
: string
| set column_format.a = 1
< error
| 1064 42000 …near 'column_format.a = 1' at line 1
> exec ^
: string
| set column_name.a = 1
< error
| 1064 42000 …near 'column_name.a = 1' at line 1
> exec ^
: string
| set columns.a = 1
< error
| 1064 42000 …near 'columns.a = 1' at line 1
> exec ^
: string
| set comment.a = 1
< error
| 1064 42000 …near 'comment.a = 1' at line 1
> exec ^
: string
| set commit.a = 1
< error
| 1064 42000 …near 'commit.a = 1' at line 1
> exec ^
: string
| set committed.a = 1
< error
| 1064 42000 …near 'committed.a = 1' at line 1
> exec ^
: string
| set compact.a = 1
< error
| 1064 42000 …near 'compact.a = 1' at line 1
> exec ^
: string
| set completion.a = 1
< error
| 1064 42000 …near 'completion.a = 1' at line 1
> exec ^
: string
| set component.a = 1
< error
| 1064 42000 …near 'component.a = 1' at line 1
> exec ^
: string
| set compressed.a = 1
< error
| 1064 42000 …near 'compressed.a = 1' at line 1
> exec ^
: string
| set compression.a = 1
< error
| 1064 42000 …near 'compression.a = 1' at line 1
> exec ^
: string
| set concat.a = 1
< error
| 1064 42000 …near 'concat.a = 1' at line 1
> exec ^
: string
| set concurrent.a = 1
< error
| 1064 42000 …near 'concurrent.a = 1' at line 1
> exec ^
: string
| set condition.a = 1
< error
| 1064 42000 …near 'condition.a = 1' at line 1
> exec ^
: string
| set connection.a = 1
< error
| 1064 42000 …near 'connection.a = 1' at line 1
> exec ^
: string
| set consistent.a = 1
< error
| 1064 42000 …near 'consistent.a = 1' at line 1
> exec ^
: string
| set constraint.a = 1
< error
| 1064 42000 …near 'constraint.a = 1' at line 1
> exec ^
: string
| set constraint_catalog.a = 1
< error
| 1064 42000 …near 'constraint_catalog.a = 1' at line 1
> exec ^
: string
| set constraint_name.a = 1
< error
| 1064 42000 …near 'constraint_name.a = 1' at line 1
> exec ^
: string
| set constraint_schema.a = 1
< error
| 1064 42000 …near 'constraint_schema.a = 1' at line 1
> exec ^
: string
| set contains.a = 1
< error
| 1064 42000 …near 'contains.a = 1' at line 1
> exec ^
: string
| set context.a = 1
< error
| 1064 42000 …near 'context.a = 1' at line 1
> exec ^
: string
| set continue.a = 1
< error
| 1064 42000 …near 'continue.a = 1' at line 1
> exec ^
: string
| set convert.a = 1
< error
| 1064 42000 …near 'convert.a = 1' at line 1
> exec ^
: string
| set count.a = 1
< error
| 1064 42000 …near 'count.a = 1' at line 1
> exec ^
: string
| set cpu.a = 1
< error
| 1064 42000 …near 'cpu.a = 1' at line 1
> exec ^
: string
| set create.a = 1
< error
| 1064 42000 …near 'create.a = 1' at line 1
> exec ^
: string
| set cross.a = 1
< error
| 1064 42000 …near 'cross.a = 1' at line 1
> exec ^
: string
| set csv.a = 1
< error
| 1064 42000 …near 'csv.a = 1' at line 1
> exec ^
: string
| set cube.a = 1
< error
| 1064 42000 …near 'cube.a = 1' at line 1
> exec ^
: string
| set cume_dist.a = 1
< error
| 1064 42000 …near 'cume_dist.a = 1' at line 1
> exec ^
: string
| set curdate.a = 1
< error
| 1064 42000 …near 'curdate.a = 1' at line 1
> exec ^
: string
| set current.a = 1
< error
| 1064 42000 …near 'current.a = 1' at line 1
> exec ^
: string
| set current_date.a = 1
< error
| 1064 42000 …near 'current_date.a = 1' at line 1
> exec ^
: string
| set current_time.a = 1
< error
| 1064 42000 …near 'current_time.a = 1' at line 1
> exec ^
: string
| set current_timestamp.a = 1
< error
| 1064 42000 …near 'current_timestamp.a = 1' at line 1
> exec ^
: string
| set current_user.a = 1
< error
| 1064 42000 …near 'current_user.a = 1' at line 1
> exec ^
: string
| set cursor.a = 1
< error
| 1064 42000 …near 'cursor.a = 1' at line 1
> exec ^
: string
| set cursor_name.a = 1
< error
| 1064 42000 …near 'cursor_name.a = 1' at line 1
> exec ^
: string
| set curtime.a = 1
< error
| 1064 42000 …near 'curtime.a = 1' at line 1
> exec ^
: string
| set data.a = 1
< error
| 1064 42000 …near 'data.a = 1' at line 1
> exec ^
: string
| set database.a = 1
< error
| 1064 42000 …near 'database.a = 1' at line 1
> exec ^
: string
| set databases.a = 1
< error
| 1064 42000 …near 'databases.a = 1' at line 1
> exec ^
: string
| set datafile.a = 1
< error
| 1064 42000 …near 'datafile.a = 1' at line 1
> exec ^
: string
| set date.a = 1
< error
| 1064 42000 …near 'date.a = 1' at line 1
> exec ^
: string
| set date_add.a = 1
< error
| 1064 42000 …near 'date_add.a = 1' at line 1
> exec ^
: string
| set date_sub.a = 1
< error
| 1064 42000 …near 'date_sub.a = 1' at line 1
> exec ^
: string
| set datetime.a = 1
< error
| 1064 42000 …near 'datetime.a = 1' at line 1
> exec ^
: string
| set day.a = 1
< error
| 1064 42000 …near 'day.a = 1' at line 1
> exec ^
: string
| set day_hour.a = 1
< error
| 1064 42000 …near 'day_hour.a = 1' at line 1
> exec ^
: string
| set day_microsecond.a = 1
< error
| 1064 42000 …near 'day_microsecond.a = 1' at line 1
> exec ^
: string
| set day_minute.a = 1
< error
| 1064 42000 …near 'day_minute.a = 1' at line 1
> exec ^
: string
| set day_second.a = 1
< error
| 1064 42000 …near 'day_second.a = 1' at line 1
> exec ^
: string
| set deallocate.a = 1
< error
| 1064 42000 …near 'deallocate.a = 1' at line 1
> exec ^
: string
| set dec.a = 1
< error
| 1064 42000 …near 'dec.a = 1' at line 1
> exec ^
: string
| set decimal.a = 1
< error
| 1064 42000 …near 'decimal.a = 1' at line 1
> exec ^
: string
| set declare.a = 1
< error
| 1064 42000 …near 'declare.a = 1' at line 1
> exec ^
: string
| set default.a = 1
< error
| 1064 42000 …near 'default.a = 1' at line 1
> exec ^
: string
| set default_auth.a = 1
< error
| 1064 42000 …near 'default_auth.a = 1' at line 1
> exec ^
: string
| set definer.a = 1
< error
| 1064 42000 …near 'definer.a = 1' at line 1
> exec ^
: string
| set definition.a = 1
< error
| 1064 42000 …near 'definition.a = 1' at line 1
> exec ^
: string
| set delay_key_write.a = 1
< error
| 1064 42000 …near 'delay_key_write.a = 1' at line 1
> exec ^
: string
| set delayed.a = 1
< error
| 1064 42000 …near 'delayed.a = 1' at line 1
> exec ^
: string
| set delete.a = 1
< error
| 1064 42000 …near 'delete.a = 1' at line 1
> exec ^
: string
| set dense_rank.a = 1
< error
| 1064 42000 …near 'dense_rank.a = 1' at line 1
> exec ^
: string
| set desc.a = 1
< error
| 1064 42000 …near 'desc.a = 1' at line 1
> exec ^
: string
| set describe.a = 1
< error
| 1064 42000 …near 'describe.a = 1' at line 1
> exec ^
: string
| set description.a = 1
< error
| 1064 42000 …near 'description.a = 1' at line 1
> exec ^
: string
| set deterministic.a = 1
< error
| 1064 42000 …near 'deterministic.a = 1' at line 1
> exec ^
: string
| set diagnostics.a = 1
< error
| 1064 42000 …near 'diagnostics.a = 1' at line 1
> exec ^
: string
| set directory.a = 1
< error
| 1064 42000 …near 'directory.a = 1' at line 1
> exec ^
: string
| set disable.a = 1
< error
| 1064 42000 …near 'disable.a = 1' at line 1
> exec ^
: string
| set discard.a = 1
< error
| 1064 42000 …near 'discard.a = 1' at line 1
> exec ^
: string
| set disk.a = 1
< error
| 1064 42000 …near 'disk.a = 1' at line 1
> exec ^
: string
| set distinct.a = 1
< error
| 1064 42000 …near 'distinct.a = 1' at line 1
> exec ^
: string
| set distinctrow.a = 1
< error
| 1064 42000 …near 'distinctrow.a = 1' at line 1
> exec ^
: string
| set div.a = 1
< error
| 1064 42000 …near 'div.a = 1' at line 1
> exec ^
: string
| set do.a = 1
< error
| 1064 42000 …near 'do.a = 1' at line 1
> exec ^
: string
| set double.a = 1
< error
| 1064 42000 …near 'double.a = 1' at line 1
> exec ^
: string
| set drop.a = 1
< error
| 1064 42000 …near 'drop.a = 1' at line 1
> exec ^
: string
| set dual.a = 1
< error
| 1064 42000 …near 'dual.a = 1' at line 1
> exec ^
: string
| set duality.a = 1
< error
| 1064 42000 …near 'duality.a = 1' at line 1
> exec ^
: string
| set dumpfile.a = 1
< error
| 1064 42000 …near 'dumpfile.a = 1' at line 1
> exec ^
: string
| set duplicate.a = 1
< error
| 1064 42000 …near 'duplicate.a = 1' at line 1
> exec ^
: string
| set dynamic.a = 1
< error
| 1064 42000 …near 'dynamic.a = 1' at line 1
> exec ^
: string
| set dynamic_partition_policy.a = 1
< error
| 1064 42000 …near 'dynamic_partition_policy.a = 1' at line 1
> exec ^
: string
| set each.a = 1
< error
| 1064 42000 …near 'each.a = 1' at line 1
> exec ^
: string
| set else.a = 1
< error
| 1064 42000 …near 'else.a = 1' at line 1
> exec ^
: string
| set elseif.a = 1
< error
| 1064 42000 …near 'elseif.a = 1' at line 1
> exec ^
: string
| set empty.a = 1
< error
| 1064 42000 …near 'empty.a = 1' at line 1
> exec ^
: string
| set enable.a = 1
< error
| 1064 42000 …near 'enable.a = 1' at line 1
> exec ^
: string
| set enable_macro_block_bloom_filter.a = 1
< error
| 1064 42000 …near 'enable_macro_block_bloom_filter.a = 1' at line 1
> exec ^
: string
| set enclosed.a = 1
< error
| 1064 42000 …near 'enclosed.a = 1' at line 1
> exec ^
: string
| set encryption.a = 1
< error
| 1064 42000 …near 'encryption.a = 1' at line 1
> exec ^
: string
| set end.a = 1
< error
| 1064 42000 …near 'end.a = 1' at line 1
> exec ^
: string
| set ends.a = 1
< error
| 1064 42000 …near 'ends.a = 1' at line 1
> exec ^
: string
| set enforced.a = 1
< error
| 1064 42000 …near 'enforced.a = 1' at line 1
> exec ^
: string
| set engine.a = 1
< error
| 1064 42000 …near 'engine.a = 1' at line 1
> exec ^
: string
| set engine_attribute.a = 1
< error
| 1064 42000 …near 'engine_attribute.a = 1' at line 1
> exec ^
: string
| set engines.a = 1
< error
| 1064 42000 …near 'engines.a = 1' at line 1
> exec ^
: string
| set enum.a = 1
< error
| 1064 42000 …near 'enum.a = 1' at line 1
> exec ^
: string
| set error.a = 1
< error
| 1064 42000 …near 'error.a = 1' at line 1
> exec ^
: string
| set errors.a = 1
< error
| 1064 42000 …near 'errors.a = 1' at line 1
> exec ^
: string
| set escape.a = 1
< error
| 1064 42000 …near 'escape.a = 1' at line 1
> exec ^
: string
| set escaped.a = 1
< error
| 1064 42000 …near 'escaped.a = 1' at line 1
> exec ^
: string
| set event.a = 1
< error
| 1064 42000 …near 'event.a = 1' at line 1
> exec ^
: string
| set events.a = 1
< error
| 1064 42000 …near 'events.a = 1' at line 1
> exec ^
: string
| set every.a = 1
< error
| 1064 42000 …near 'every.a = 1' at line 1
> exec ^
: string
| set except.a = 1
< error
| 1064 42000 …near 'except.a = 1' at line 1
> exec ^
: string
| set exchange.a = 1
< error
| 1064 42000 …near 'exchange.a = 1' at line 1
> exec ^
: string
| set exclude.a = 1
< error
| 1064 42000 …near 'exclude.a = 1' at line 1
> exec ^
: string
| set execute.a = 1
< error
| 1064 42000 …near 'execute.a = 1' at line 1
> exec ^
: string
| set exists.a = 1
< error
| 1064 42000 …near 'exists.a = 1' at line 1
> exec ^
: string
| set exit.a = 1
< error
| 1064 42000 …near 'exit.a = 1' at line 1
> exec ^
: string
| set expansion.a = 1
< error
| 1064 42000 …near 'expansion.a = 1' at line 1
> exec ^
: string
| set expire.a = 1
< error
| 1064 42000 …near 'expire.a = 1' at line 1
> exec ^
: string
| set explain.a = 1
< error
| 1064 42000 …near 'explain.a = 1' at line 1
> exec ^
: string
| set export.a = 1
< error
| 1064 42000 …near 'export.a = 1' at line 1
> exec ^
: string
| set extended.a = 1
< error
| 1064 42000 …near 'extended.a = 1' at line 1
> exec ^
: string
| set extent_size.a = 1
< error
| 1064 42000 …near 'extent_size.a = 1' at line 1
> exec ^
: string
| set external.a = 1
< error
| 1064 42000 …near 'external.a = 1' at line 1
> exec ^
: string
| set external_format.a = 1
< error
| 1064 42000 …near 'external_format.a = 1' at line 1
> exec ^
: string
| set extract.a = 1
< error
| 1064 42000 …near 'extract.a = 1' at line 1
> exec ^
: string
| set factor.a = 1
< error
| 1064 42000 …near 'factor.a = 1' at line 1
> exec ^
: string
| set failed_login_attempts.a = 1
< error
| 1064 42000 …near 'failed_login_attempts.a = 1' at line 1
> exec ^
: string
| set false.a = 1
< error
| 1064 42000 …near 'false.a = 1' at line 1
> exec ^
: string
| set fast.a = 1
< error
| 1064 42000 …near 'fast.a = 1' at line 1
> exec ^
: string
| set faults.a = 1
< error
| 1064 42000 …near 'faults.a = 1' at line 1
> exec ^
: string
| set fetch.a = 1
< error
| 1064 42000 …near 'fetch.a = 1' at line 1
> exec ^
: string
| set fields.a = 1
< error
| 1064 42000 …near 'fields.a = 1' at line 1
> exec ^
: string
| set file.a = 1
< error
| 1064 42000 …near 'file.a = 1' at line 1
> exec ^
: string
| set file_block_size.a = 1
< error
| 1064 42000 …near 'file_block_size.a = 1' at line 1
> exec ^
: string
| set file_format.a = 1
< error
| 1064 42000 …near 'file_format.a = 1' at line 1
> exec ^
: string
| set file_name.a = 1
< error
| 1064 42000 …near 'file_name.a = 1' at line 1
> exec ^
: string
| set file_pattern.a = 1
< error
| 1064 42000 …near 'file_pattern.a = 1' at line 1
> exec ^
: string
| set file_prefix.a = 1
< error
| 1064 42000 …near 'file_prefix.a = 1' at line 1
> exec ^
: string
| set files.a = 1
< error
| 1064 42000 …near 'files.a = 1' at line 1
> exec ^
: string
| set filter.a = 1
< error
| 1064 42000 …near 'filter.a = 1' at line 1
> exec ^
: string
| set finish.a = 1
< error
| 1064 42000 …near 'finish.a = 1' at line 1
> exec ^
: string
| set first.a = 1
< error
| 1064 42000 …near 'first.a = 1' at line 1
> exec ^
: string
| set first_value.a = 1
< error
| 1064 42000 …near 'first_value.a = 1' at line 1
> exec ^
: string
| set fixed.a = 1
< error
| 1064 42000 …near 'fixed.a = 1' at line 1
> exec ^
: string
| set float.a = 1
< error
| 1064 42000 …near 'float.a = 1' at line 1
> exec ^
: string
| set float4.a = 1
< error
| 1064 42000 …near 'float4.a = 1' at line 1
> exec ^
: string
| set float8.a = 1
< error
| 1064 42000 …near 'float8.a = 1' at line 1
> exec ^
: string
| set flush.a = 1
< error
| 1064 42000 …near 'flush.a = 1' at line 1
> exec ^
: string
| set following.a = 1
< error
| 1064 42000 …near 'following.a = 1' at line 1
> exec ^
: string
| set follows.a = 1
< error
| 1064 42000 …near 'follows.a = 1' at line 1
> exec ^
: string
| set for.a = 1
< error
| 1064 42000 …near 'for.a = 1' at line 1
> exec ^
: string
| set force.a = 1
< error
| 1064 42000 …near 'force.a = 1' at line 1
> exec ^
: string
| set foreign.a = 1
< error
| 1064 42000 …near 'foreign.a = 1' at line 1
> exec ^
: string
| set format.a = 1
< error
| 1064 42000 …near 'format.a = 1' at line 1
> exec ^
: string
| set found.a = 1
< error
| 1064 42000 …near 'found.a = 1' at line 1
> exec ^
: string
| set from.a = 1
< error
| 1064 42000 …near 'from.a = 1' at line 1
> exec ^
: string
| set full.a = 1
< error
| 1064 42000 …near 'full.a = 1' at line 1
> exec ^
: string
| set fulltext.a = 1
< error
| 1064 42000 …near 'fulltext.a = 1' at line 1
> exec ^
: string
| set function.a = 1
< error
| 1064 42000 …near 'function.a = 1' at line 1
> exec ^
: string
| set general.a = 1
< error
| 1064 42000 …near 'general.a = 1' at line 1
> exec ^
: string
| set generate.a = 1
< error
| 1064 42000 …near 'generate.a = 1' at line 1
> exec ^
: string
| set generated.a = 1
< error
| 1064 42000 …near 'generated.a = 1' at line 1
> exec ^
: string
| set geomcollection.a = 1
< error
| 1064 42000 …near 'geomcollection.a = 1' at line 1
> exec ^
: string
| set geometry.a = 1
< error
| 1064 42000 …near 'geometry.a = 1' at line 1
> exec ^
: string
| set geometrycollection.a = 1
< error
| 1064 42000 …near 'geometrycollection.a = 1' at line 1
> exec ^
: string
| set get.a = 1
< error
| 1064 42000 …near 'get.a = 1' at line 1
> exec ^
: string
| set get_format.a = 1
< error
| 1064 42000 …near 'get_format.a = 1' at line 1
> exec ^
: string
| set get_master_public_key.a = 1
< error
| 1064 42000 …near 'get_master_public_key.a = 1' at line 1
> exec ^
: string
| set get_source_public_key.a = 1
< error
| 1064 42000 …near 'get_source_public_key.a = 1' at line 1
> exec ^
: string
| set global.a = 1
< error
| 1064 42000 …near 'global.a = 1' at line 1
> exec ^
: string
| set grant.a = 1
< error
| 1064 42000 …near 'grant.a = 1' at line 1
> exec ^
: string
| set grants.a = 1
< error
| 1064 42000 …near 'grants.a = 1' at line 1
> exec ^
: string
| set greatest.a = 1
< error
| 1064 42000 …near 'greatest.a = 1' at line 1
> exec ^
: string
| set group.a = 1
< error
| 1064 42000 …near 'group.a = 1' at line 1
> exec ^
: string
| set group_concat.a = 1
< error
| 1064 42000 …near 'group_concat.a = 1' at line 1
> exec ^
: string
| set group_replication.a = 1
< error
| 1064 42000 …near 'group_replication.a = 1' at line 1
> exec ^
: string
| set grouping.a = 1
< error
| 1064 42000 …near 'grouping.a = 1' at line 1
> exec ^
: string
| set groups.a = 1
< error
| 1064 42000 …near 'groups.a = 1' at line 1
> exec ^
: string
| set gtid_only.a = 1
< error
| 1064 42000 …near 'gtid_only.a = 1' at line 1
> exec ^
: string
| set gtids.a = 1
< error
| 1064 42000 …near 'gtids.a = 1' at line 1
> exec ^
: string
| set guided.a = 1
< error
| 1064 42000 …near 'guided.a = 1' at line 1
> exec ^
: string
| set handler.a = 1
< error
| 1064 42000 …near 'handler.a = 1' at line 1
> exec ^
: string
| set hash.a = 1
< error
| 1064 42000 …near 'hash.a = 1' at line 1
> exec ^
: string
| set having.a = 1
< error
| 1064 42000 …near 'having.a = 1' at line 1
> exec ^
: string
| set header.a = 1
< error
| 1064 42000 …near 'header.a = 1' at line 1
> exec ^
: string
| set help.a = 1
< error
| 1064 42000 …near 'help.a = 1' at line 1
> exec ^
: string
| set high_priority.a = 1
< error
| 1064 42000 …near 'high_priority.a = 1' at line 1
> exec ^
: string
| set histogram.a = 1
< error
| 1064 42000 …near 'histogram.a = 1' at line 1
> exec ^
: string
| set history.a = 1
< error
| 1064 42000 …near 'history.a = 1' at line 1
> exec ^
: string
| set host.a = 1
< error
| 1064 42000 …near 'host.a = 1' at line 1
> exec ^
: string
| set hosts.a = 1
< error
| 1064 42000 …near 'hosts.a = 1' at line 1
> exec ^
: string
| set hour.a = 1
< error
| 1064 42000 …near 'hour.a = 1' at line 1
> exec ^
: string
| set hour_microsecond.a = 1
< error
| 1064 42000 …near 'hour_microsecond.a = 1' at line 1
> exec ^
: string
| set hour_minute.a = 1
< error
| 1064 42000 …near 'hour_minute.a = 1' at line 1
> exec ^
: string
| set hour_second.a = 1
< error
| 1064 42000 …near 'hour_second.a = 1' at line 1
> exec ^
: string
| set identified.a = 1
< error
| 1064 42000 …near 'identified.a = 1' at line 1
> exec ^
: string
| set if.a = 1
< error
| 1064 42000 …near 'if.a = 1' at line 1
> exec ^
: string
| set ignore.a = 1
< error
| 1064 42000 …near 'ignore.a = 1' at line 1
> exec ^
: string
| set ignore_server_ids.a = 1
< error
| 1064 42000 …near 'ignore_server_ids.a = 1' at line 1
> exec ^
: string
| set import.a = 1
< error
| 1064 42000 …near 'import.a = 1' at line 1
> exec ^
: string
| set in.a = 1
< error
| 1064 42000 …near 'in.a = 1' at line 1
> exec ^
: string
| set inactive.a = 1
< error
| 1064 42000 …near 'inactive.a = 1' at line 1
> exec ^
: string
| set index.a = 1
< error
| 1064 42000 …near 'index.a = 1' at line 1
> exec ^
: string
| set indexes.a = 1
< error
| 1064 42000 …near 'indexes.a = 1' at line 1
> exec ^
: string
| set infile.a = 1
< error
| 1064 42000 …near 'infile.a = 1' at line 1
> exec ^
: string
| set initial.a = 1
< error
| 1064 42000 …near 'initial.a = 1' at line 1
> exec ^
: string
| set initial_size.a = 1
< error
| 1064 42000 …near 'initial_size.a = 1' at line 1
> exec ^
: string
| set initiate.a = 1
< error
| 1064 42000 …near 'initiate.a = 1' at line 1
> exec ^
: string
| set inner.a = 1
< error
| 1064 42000 …near 'inner.a = 1' at line 1
> exec ^
: string
| set inout.a = 1
< error
| 1064 42000 …near 'inout.a = 1' at line 1
> exec ^
: string
| set insensitive.a = 1
< error
| 1064 42000 …near 'insensitive.a = 1' at line 1
> exec ^
: string
| set insert.a = 1
< error
| 1064 42000 …near 'insert.a = 1' at line 1
> exec ^
: string
| set insert_method.a = 1
< error
| 1064 42000 …near 'insert_method.a = 1' at line 1
> exec ^
: string
| set install.a = 1
< error
| 1064 42000 …near 'install.a = 1' at line 1
> exec ^
: string
| set instance.a = 1
< error
| 1064 42000 …near 'instance.a = 1' at line 1
> exec ^
: string
| set int.a = 1
< error
| 1064 42000 …near 'int.a = 1' at line 1
> exec ^
: string
| set int1.a = 1
< error
| 1064 42000 …near 'int1.a = 1' at line 1
> exec ^
: string
| set int2.a = 1
< error
| 1064 42000 …near 'int2.a = 1' at line 1
> exec ^
: string
| set int3.a = 1
< error
| 1064 42000 …near 'int3.a = 1' at line 1
> exec ^
: string
| set int4.a = 1
< error
| 1064 42000 …near 'int4.a = 1' at line 1
> exec ^
: string
| set int8.a = 1
< error
| 1064 42000 …near 'int8.a = 1' at line 1
> exec ^
: string
| set integer.a = 1
< error
| 1064 42000 …near 'integer.a = 1' at line 1
> exec ^
: string
| set intersect.a = 1
< error
| 1064 42000 …near 'intersect.a = 1' at line 1
> exec ^
: string
| set interval.a = 1
< error
| 1064 42000 …near 'interval.a = 1' at line 1
> exec ^
: string
| set into.a = 1
< error
| 1064 42000 …near 'into.a = 1' at line 1
> exec ^
: string
| set invisible.a = 1
< error
| 1064 42000 …near 'invisible.a = 1' at line 1
> exec ^
: string
| set invoker.a = 1
< error
| 1064 42000 …near 'invoker.a = 1' at line 1
> exec ^
: string
| set io.a = 1
< error
| 1064 42000 …near 'io.a = 1' at line 1
> exec ^
: string
| set io_after_gtids.a = 1
< error
| 1064 42000 …near 'io_after_gtids.a = 1' at line 1
> exec ^
: string
| set io_before_gtids.a = 1
< error
| 1064 42000 …near 'io_before_gtids.a = 1' at line 1
> exec ^
: string
| set io_thread.a = 1
< error
| 1064 42000 …near 'io_thread.a = 1' at line 1
> exec ^
: string
| set ipc.a = 1
< error
| 1064 42000 …near 'ipc.a = 1' at line 1
> exec ^
: string
| set is.a = 1
< error
| 1064 42000 …near 'is.a = 1' at line 1
> exec ^
: string
| set isolation.a = 1
< error
| 1064 42000 …near 'isolation.a = 1' at line 1
> exec ^
: string
| set issuer.a = 1
< error
| 1064 42000 …near 'issuer.a = 1' at line 1
> exec ^
: string
| set iterate.a = 1
< error
| 1064 42000 …near 'iterate.a = 1' at line 1
> exec ^
: string
| set join.a = 1
< error
| 1064 42000 …near 'join.a = 1' at line 1
> exec ^
: string
| set json.a = 1
< error
| 1064 42000 …near 'json.a = 1' at line 1
> exec ^
: string
| set json_arrayagg.a = 1
< error
| 1064 42000 …near 'json_arrayagg.a = 1' at line 1
> exec ^
: string
| set json_duality_object.a = 1
< error
| 1064 42000 …near 'json_duality_object.a = 1' at line 1
> exec ^
: string
| set json_objectagg.a = 1
< error
| 1064 42000 …near 'json_objectagg.a = 1' at line 1
> exec ^
: string
| set json_table.a = 1
< error
| 1064 42000 …near 'json_table.a = 1' at line 1
> exec ^
: string
| set json_value.a = 1
< error
| 1064 42000 …near 'json_value.a = 1' at line 1
> exec ^
: string
| set key.a = 1
< error
| 1064 42000 …near 'key.a = 1' at line 1
> exec ^
: string
| set key_block_size.a = 1
< error
| 1064 42000 …near 'key_block_size.a = 1' at line 1
> exec ^
: string
| set keyring.a = 1
< error
| 1064 42000 …near 'keyring.a = 1' at line 1
> exec ^
: string
| set keys.a = 1
< error
| 1064 42000 …near 'keys.a = 1' at line 1
> exec ^
: string
| set kill.a = 1
< error
| 1064 42000 …near 'kill.a = 1' at line 1
> exec ^
: string
| set lag.a = 1
< error
| 1064 42000 …near 'lag.a = 1' at line 1
> exec ^
: string
| set language.a = 1
< error
| 1064 42000 …near 'language.a = 1' at line 1
> exec ^
: string
| set last.a = 1
< error
| 1064 42000 …near 'last.a = 1' at line 1
> exec ^
: string
| set last_value.a = 1
< error
| 1064 42000 …near 'last_value.a = 1' at line 1
> exec ^
: string
| set lateral.a = 1
< error
| 1064 42000 …near 'lateral.a = 1' at line 1
> exec ^
: string
| set lead.a = 1
< error
| 1064 42000 …near 'lead.a = 1' at line 1
> exec ^
: string
| set leading.a = 1
< error
| 1064 42000 …near 'leading.a = 1' at line 1
> exec ^
: string
| set least.a = 1
< error
| 1064 42000 …near 'least.a = 1' at line 1
> exec ^
: string
| set leave.a = 1
< error
| 1064 42000 …near 'leave.a = 1' at line 1
> exec ^
: string
| set leaves.a = 1
< error
| 1064 42000 …near 'leaves.a = 1' at line 1
> exec ^
: string
| set left.a = 1
< error
| 1064 42000 …near 'left.a = 1' at line 1
> exec ^
: string
| set less.a = 1
< error
| 1064 42000 …near 'less.a = 1' at line 1
> exec ^
: string
| set level.a = 1
< error
| 1064 42000 …near 'level.a = 1' at line 1
> exec ^
: string
| set library.a = 1
< error
| 1064 42000 …near 'library.a = 1' at line 1
> exec ^
: string
| set like.a = 1
< error
| 1064 42000 …near 'like.a = 1' at line 1
> exec ^
: string
| set limit.a = 1
< error
| 1064 42000 …near 'limit.a = 1' at line 1
> exec ^
: string
| set linear.a = 1
< error
| 1064 42000 …near 'linear.a = 1' at line 1
> exec ^
: string
| set lines.a = 1
< error
| 1064 42000 …near 'lines.a = 1' at line 1
> exec ^
: string
| set linestring.a = 1
< error
| 1064 42000 …near 'linestring.a = 1' at line 1
> exec ^
: string
| set list.a = 1
< error
| 1064 42000 …near 'list.a = 1' at line 1
> exec ^
: string
| set load.a = 1
< error
| 1064 42000 …near 'load.a = 1' at line 1
> exec ^
: string
| set local.a = 1
< error
| 1064 42000 …near 'local.a = 1' at line 1
> exec ^
: string
| set localtime.a = 1
< error
| 1064 42000 …near 'localtime.a = 1' at line 1
> exec ^
: string
| set localtimestamp.a = 1
< error
| 1064 42000 …near 'localtimestamp.a = 1' at line 1
> exec ^
: string
| set lock.a = 1
< error
| 1064 42000 …near 'lock.a = 1' at line 1
> exec ^
: string
| set locked.a = 1
< error
| 1064 42000 …near 'locked.a = 1' at line 1
> exec ^
: string
| set locks.a = 1
< error
| 1064 42000 …near 'locks.a = 1' at line 1
> exec ^
: string
| set log.a = 1
< error
| 1064 42000 …near 'log.a = 1' at line 1
> exec ^
: string
| set logfile.a = 1
< error
| 1064 42000 …near 'logfile.a = 1' at line 1
> exec ^
: string
| set logs.a = 1
< error
| 1064 42000 …near 'logs.a = 1' at line 1
> exec ^
: string
| set long.a = 1
< error
| 1064 42000 …near 'long.a = 1' at line 1
> exec ^
: string
| set longblob.a = 1
< error
| 1064 42000 …near 'longblob.a = 1' at line 1
> exec ^
: string
| set longtext.a = 1
< error
| 1064 42000 …near 'longtext.a = 1' at line 1
> exec ^
: string
| set loop.a = 1
< error
| 1064 42000 …near 'loop.a = 1' at line 1
> exec ^
: string
| set low_priority.a = 1
< error
| 1064 42000 …near 'low_priority.a = 1' at line 1
> exec ^
: string
| set manifest.a = 1
< error
| 1064 42000 …near 'manifest.a = 1' at line 1
> exec ^
: string
| set manual.a = 1
< error
| 1064 42000 …near 'manual.a = 1' at line 1
> exec ^
: string
| set master.a = 1
< error
| 1064 42000 …near 'master.a = 1' at line 1
> exec ^
: string
| set master_auto_position.a = 1
< error
| 1064 42000 …near 'master_auto_position.a = 1' at line 1
> exec ^
: string
| set master_bind.a = 1
< error
| 1064 42000 …near 'master_bind.a = 1' at line 1
> exec ^
: string
| set master_compression_algorithms.a = 1
< error
| 1064 42000 …near 'master_compression_algorithms.a = 1' at line 1
> exec ^
: string
| set master_connect_retry.a = 1
< error
| 1064 42000 …near 'master_connect_retry.a = 1' at line 1
> exec ^
: string
| set master_delay.a = 1
< error
| 1064 42000 …near 'master_delay.a = 1' at line 1
> exec ^
: string
| set master_heartbeat_period.a = 1
< error
| 1064 42000 …near 'master_heartbeat_period.a = 1' at line 1
> exec ^
: string
| set master_host.a = 1
< error
| 1064 42000 …near 'master_host.a = 1' at line 1
> exec ^
: string
| set master_log_file.a = 1
< error
| 1064 42000 …near 'master_log_file.a = 1' at line 1
> exec ^
: string
| set master_log_pos.a = 1
< error
| 1064 42000 …near 'master_log_pos.a = 1' at line 1
> exec ^
: string
| set master_password.a = 1
< error
| 1064 42000 …near 'master_password.a = 1' at line 1
> exec ^
: string
| set master_port.a = 1
< error
| 1064 42000 …near 'master_port.a = 1' at line 1
> exec ^
: string
| set master_public_key_path.a = 1
< error
| 1064 42000 …near 'master_public_key_path.a = 1' at line 1
> exec ^
: string
| set master_retry_count.a = 1
< error
| 1064 42000 …near 'master_retry_count.a = 1' at line 1
> exec ^
: string
| set master_ssl.a = 1
< error
| 1064 42000 …near 'master_ssl.a = 1' at line 1
> exec ^
: string
| set master_ssl_ca.a = 1
< error
| 1064 42000 …near 'master_ssl_ca.a = 1' at line 1
> exec ^
: string
| set master_ssl_capath.a = 1
< error
| 1064 42000 …near 'master_ssl_capath.a = 1' at line 1
> exec ^
: string
| set master_ssl_cert.a = 1
< error
| 1064 42000 …near 'master_ssl_cert.a = 1' at line 1
> exec ^
: string
| set master_ssl_cipher.a = 1
< error
| 1064 42000 …near 'master_ssl_cipher.a = 1' at line 1
> exec ^
: string
| set master_ssl_crl.a = 1
< error
| 1064 42000 …near 'master_ssl_crl.a = 1' at line 1
> exec ^
: string
| set master_ssl_crlpath.a = 1
< error
| 1064 42000 …near 'master_ssl_crlpath.a = 1' at line 1
> exec ^
: string
| set master_ssl_key.a = 1
< error
| 1064 42000 …near 'master_ssl_key.a = 1' at line 1
> exec ^
: string
| set master_ssl_verify_server_cert.a = 1
< error
| 1064 42000 …near 'master_ssl_verify_server_cert.a = 1' at line 1
> exec ^
: string
| set master_tls_ciphersuites.a = 1
< error
| 1064 42000 …near 'master_tls_ciphersuites.a = 1' at line 1
> exec ^
: string
| set master_tls_version.a = 1
< error
| 1064 42000 …near 'master_tls_version.a = 1' at line 1
> exec ^
: string
| set master_user.a = 1
< error
| 1064 42000 …near 'master_user.a = 1' at line 1
> exec ^
: string
| set master_zstd_compression_level.a = 1
< error
| 1064 42000 …near 'master_zstd_compression_level.a = 1' at line 1
> exec ^
: string
| set match.a = 1
< error
| 1064 42000 …near 'match.a = 1' at line 1
> exec ^
: string
| set materialized.a = 1
< error
| 1064 42000 …near 'materialized.a = 1' at line 1
> exec ^
: string
| set max.a = 1
< error
| 1064 42000 …near 'max.a = 1' at line 1
> exec ^
: string
| set max_connections_per_hour.a = 1
< error
| 1064 42000 …near 'max_connections_per_hour.a = 1' at line 1
> exec ^
: string
| set max_queries_per_hour.a = 1
< error
| 1064 42000 …near 'max_queries_per_hour.a = 1' at line 1
> exec ^
: string
| set max_rows.a = 1
< error
| 1064 42000 …near 'max_rows.a = 1' at line 1
> exec ^
: string
| set max_size.a = 1
< error
| 1064 42000 …near 'max_size.a = 1' at line 1
> exec ^
: string
| set max_updates_per_hour.a = 1
< error
| 1064 42000 …near 'max_updates_per_hour.a = 1' at line 1
> exec ^
: string
| set max_user_connections.a = 1
< error
| 1064 42000 …near 'max_user_connections.a = 1' at line 1
> exec ^
: string
| set maxvalue.a = 1
< error
| 1064 42000 …near 'maxvalue.a = 1' at line 1
> exec ^
: string
| set medium.a = 1
< error
| 1064 42000 …near 'medium.a = 1' at line 1
> exec ^
: string
| set mediumblob.a = 1
< error
| 1064 42000 …near 'mediumblob.a = 1' at line 1
> exec ^
: string
| set mediumint.a = 1
< error
| 1064 42000 …near 'mediumint.a = 1' at line 1
> exec ^
: string
| set mediumtext.a = 1
< error
| 1064 42000 …near 'mediumtext.a = 1' at line 1
> exec ^
: string
| set member.a = 1
< error
| 1064 42000 …near 'member.a = 1' at line 1
> exec ^
: string
| set memory.a = 1
< error
| 1064 42000 …near 'memory.a = 1' at line 1
> exec ^
: string
| set merge.a = 1
< error
| 1064 42000 …near 'merge.a = 1' at line 1
> exec ^
: string
| set message_text.a = 1
< error
| 1064 42000 …near 'message_text.a = 1' at line 1
> exec ^
: string
| set microsecond.a = 1
< error
| 1064 42000 …near 'microsecond.a = 1' at line 1
> exec ^
: string
| set mid.a = 1
< error
| 1064 42000 …near 'mid.a = 1' at line 1
> exec ^
: string
| set middleint.a = 1
< error
| 1064 42000 …near 'middleint.a = 1' at line 1
> exec ^
: string
| set migrate.a = 1
< error
| 1064 42000 …near 'migrate.a = 1' at line 1
> exec ^
: string
| set min.a = 1
< error
| 1064 42000 …near 'min.a = 1' at line 1
> exec ^
: string
| set min_rows.a = 1
< error
| 1064 42000 …near 'min_rows.a = 1' at line 1
> exec ^
: string
| set minute.a = 1
< error
| 1064 42000 …near 'minute.a = 1' at line 1
> exec ^
: string
| set minute_microsecond.a = 1
< error
| 1064 42000 …near 'minute_microsecond.a = 1' at line 1
> exec ^
: string
| set minute_second.a = 1
< error
| 1064 42000 …near 'minute_second.a = 1' at line 1
> exec ^
: string
| set mod.a = 1
< error
| 1064 42000 …near 'mod.a = 1' at line 1
> exec ^
: string
| set mode.a = 1
< error
| 1064 42000 …near 'mode.a = 1' at line 1
> exec ^
: string
| set modifies.a = 1
< error
| 1064 42000 …near 'modifies.a = 1' at line 1
> exec ^
: string
| set modify.a = 1
< error
| 1064 42000 …near 'modify.a = 1' at line 1
> exec ^
: string
| set month.a = 1
< error
| 1064 42000 …near 'month.a = 1' at line 1
> exec ^
: string
| set multilinestring.a = 1
< error
| 1064 42000 …near 'multilinestring.a = 1' at line 1
> exec ^
: string
| set multipoint.a = 1
< error
| 1064 42000 …near 'multipoint.a = 1' at line 1
> exec ^
: string
| set multipolygon.a = 1
< error
| 1064 42000 …near 'multipolygon.a = 1' at line 1
> exec ^
: string
| set mutex.a = 1
< error
| 1064 42000 …near 'mutex.a = 1' at line 1
> exec ^
: string
| set mysql_errno.a = 1
< error
| 1064 42000 …near 'mysql_errno.a = 1' at line 1
> exec ^
: string
| set name.a = 1
< error
| 1064 42000 …near 'name.a = 1' at line 1
> exec ^
: string
| set names.a = 1
< error
| 1064 42000 …near 'names.a = 1' at line 1
> exec ^
: string
| set national.a = 1
< error
| 1064 42000 …near 'national.a = 1' at line 1
> exec ^
: string
| set natural.a = 1
< error
| 1064 42000 …near 'natural.a = 1' at line 1
> exec ^
: string
| set nchar.a = 1
< error
| 1064 42000 …near 'nchar.a = 1' at line 1
> exec ^
: string
| set ndb.a = 1
< error
| 1064 42000 …near 'ndb.a = 1' at line 1
> exec ^
: string
| set ndbcluster.a = 1
< error
| 1064 42000 …near 'ndbcluster.a = 1' at line 1
> exec ^
: string
| set nested.a = 1
< error
| 1064 42000 …near 'nested.a = 1' at line 1
> exec ^
: string
| set network_namespace.a = 1
< error
| 1064 42000 …near 'network_namespace.a = 1' at line 1
> exec ^
: string
| set never.a = 1
< error
| 1064 42000 …near 'never.a = 1' at line 1
> exec ^
: string
| set new.a = 1
< error
| 1064 42000 …near 'new.a = 1' at line 1
> exec ^
: string
| set next.a = 1
< error
| 1064 42000 …near 'next.a = 1' at line 1
> exec ^
: string
| set no.a = 1
< error
| 1064 42000 …near 'no.a = 1' at line 1
> exec ^
: string
| set no_wait.a = 1
< error
| 1064 42000 …near 'no_wait.a = 1' at line 1
> exec ^
: string
| set no_write_to_binlog.a = 1
< error
| 1064 42000 …near 'no_write_to_binlog.a = 1' at line 1
> exec ^
: string
| set nodegroup.a = 1
< error
| 1064 42000 …near 'nodegroup.a = 1' at line 1
> exec ^
: string
| set none.a = 1
< error
| 1064 42000 …near 'none.a = 1' at line 1
> exec ^
: string
| set not.a = 1
< error
| 1064 42000 …near 'not.a = 1' at line 1
> exec ^
: string
| set now.a = 1
< error
| 1064 42000 …near 'now.a = 1' at line 1
> exec ^
: string
| set nowait.a = 1
< error
| 1064 42000 …near 'nowait.a = 1' at line 1
> exec ^
: string
| set nth_value.a = 1
< error
| 1064 42000 …near 'nth_value.a = 1' at line 1
> exec ^
: string
| set ntile.a = 1
< error
| 1064 42000 …near 'ntile.a = 1' at line 1
> exec ^
: string
| set null.a = 1
< error
| 1064 42000 …near 'null.a = 1' at line 1
> exec ^
: string
| set nullif.a = 1
< error
| 1064 42000 …near 'nullif.a = 1' at line 1
> exec ^
: string
| set nulls.a = 1
< error
| 1064 42000 …near 'nulls.a = 1' at line 1
> exec ^
: string
| set number.a = 1
< error
| 1064 42000 …near 'number.a = 1' at line 1
> exec ^
: string
| set numeric.a = 1
< error
| 1064 42000 …near 'numeric.a = 1' at line 1
> exec ^
: string
| set nvarchar.a = 1
< error
| 1064 42000 …near 'nvarchar.a = 1' at line 1
> exec ^
: string
| set of.a = 1
< error
| 1064 42000 …near 'of.a = 1' at line 1
> exec ^
: string
| set off.a = 1
< error
| 1064 42000 …near 'off.a = 1' at line 1
> exec ^
: string
| set offset.a = 1
< error
| 1064 42000 …near 'offset.a = 1' at line 1
> exec ^
: string
| set oj.a = 1
< error
| 1064 42000 …near 'oj.a = 1' at line 1
> exec ^
: string
| set old.a = 1
< error
| 1064 42000 …near 'old.a = 1' at line 1
> exec ^
: string
| set on.a = 1
< error
| 1064 42000 …near 'on.a = 1' at line 1
> exec ^
: string
| set one.a = 1
< error
| 1064 42000 …near 'one.a = 1' at line 1
> exec ^
: string
| set only.a = 1
< error
| 1064 42000 …near 'only.a = 1' at line 1
> exec ^
: string
| set open.a = 1
< error
| 1064 42000 …near 'open.a = 1' at line 1
> exec ^
: string
| set optimize.a = 1
< error
| 1064 42000 …near 'optimize.a = 1' at line 1
> exec ^
: string
| set optimizer_costs.a = 1
< error
| 1064 42000 …near 'optimizer_costs.a = 1' at line 1
> exec ^
: string
| set option.a = 1
< error
| 1064 42000 …near 'option.a = 1' at line 1
> exec ^
: string
| set optional.a = 1
< error
| 1064 42000 …near 'optional.a = 1' at line 1
> exec ^
: string
| set optionally.a = 1
< error
| 1064 42000 …near 'optionally.a = 1' at line 1
> exec ^
: string
| set options.a = 1
< error
| 1064 42000 …near 'options.a = 1' at line 1
> exec ^
: string
| set or.a = 1
< error
| 1064 42000 …near 'or.a = 1' at line 1
> exec ^
: string
| set order.a = 1
< error
| 1064 42000 …near 'order.a = 1' at line 1
> exec ^
: string
| set ordinality.a = 1
< error
| 1064 42000 …near 'ordinality.a = 1' at line 1
> exec ^
: string
| set organization.a = 1
< error
| 1064 42000 …near 'organization.a = 1' at line 1
> exec ^
: string
| set others.a = 1
< error
| 1064 42000 …near 'others.a = 1' at line 1
> exec ^
: string
| set out.a = 1
< error
| 1064 42000 …near 'out.a = 1' at line 1
> exec ^
: string
| set outer.a = 1
< error
| 1064 42000 …near 'outer.a = 1' at line 1
> exec ^
: string
| set outfile.a = 1
< error
| 1064 42000 …near 'outfile.a = 1' at line 1
> exec ^
: string
| set over.a = 1
< error
| 1064 42000 …near 'over.a = 1' at line 1
> exec ^
: string
| set overlay.a = 1
< error
| 1064 42000 …near 'overlay.a = 1' at line 1
> exec ^
: string
| set overwrite.a = 1
< error
| 1064 42000 …near 'overwrite.a = 1' at line 1
> exec ^
: string
| set owner.a = 1
< error
| 1064 42000 …near 'owner.a = 1' at line 1
> exec ^
: string
| set pack_keys.a = 1
< error
| 1064 42000 …near 'pack_keys.a = 1' at line 1
> exec ^
: string
| set page.a = 1
< error
| 1064 42000 …near 'page.a = 1' at line 1
> exec ^
: string
| set parallel.a = 1
< error
| 1064 42000 …near 'parallel.a = 1' at line 1
> exec ^
: string
| set parameters.a = 1
< error
| 1064 42000 …near 'parameters.a = 1' at line 1
> exec ^
: string
| set parse_tree.a = 1
< error
| 1064 42000 …near 'parse_tree.a = 1' at line 1
> exec ^
: string
| set parser.a = 1
< error
| 1064 42000 …near 'parser.a = 1' at line 1
> exec ^
: string
| set partial.a = 1
< error
| 1064 42000 …near 'partial.a = 1' at line 1
> exec ^
: string
| set partition.a = 1
< error
| 1064 42000 …near 'partition.a = 1' at line 1
> exec ^
: string
| set partitioning.a = 1
< error
| 1064 42000 …near 'partitioning.a = 1' at line 1
> exec ^
: string
| set partitions.a = 1
< error
| 1064 42000 …near 'partitions.a = 1' at line 1
> exec ^
: string
| set password.a = 1
< error
| 1064 42000 …near 'password.a = 1' at line 1
> exec ^
: string
| set password_lock_time.a = 1
< error
| 1064 42000 …near 'password_lock_time.a = 1' at line 1
> exec ^
: string
| set path.a = 1
< error
| 1064 42000 …near 'path.a = 1' at line 1
> exec ^
: string
| set pctfree.a = 1
< error
| 1064 42000 …near 'pctfree.a = 1' at line 1
> exec ^
: string
| set percent_rank.a = 1
< error
| 1064 42000 …near 'percent_rank.a = 1' at line 1
> exec ^
: string
| set persist.a = 1
< error
| 1064 42000 …near 'persist.a = 1' at line 1
> exec ^
: string
| set persist_only.a = 1
< error
| 1064 42000 …near 'persist_only.a = 1' at line 1
> exec ^
: string
| set phase.a = 1
< error
| 1064 42000 …near 'phase.a = 1' at line 1
> exec ^
: string
| set placing.a = 1
< error
| 1064 42000 …near 'placing.a = 1' at line 1
> exec ^
: string
| set plugin.a = 1
< error
| 1064 42000 …near 'plugin.a = 1' at line 1
> exec ^
: string
| set plugin_dir.a = 1
< error
| 1064 42000 …near 'plugin_dir.a = 1' at line 1
> exec ^
: string
| set plugins.a = 1
< error
| 1064 42000 …near 'plugins.a = 1' at line 1
> exec ^
: string
| set point.a = 1
< error
| 1064 42000 …near 'point.a = 1' at line 1
> exec ^
: string
| set polygon.a = 1
< error
| 1064 42000 …near 'polygon.a = 1' at line 1
> exec ^
: string
| set port.a = 1
< error
| 1064 42000 …near 'port.a = 1' at line 1
> exec ^
: string
| set position.a = 1
< error
| 1064 42000 …near 'position.a = 1' at line 1
> exec ^
: string
| set precedes.a = 1
< error
| 1064 42000 …near 'precedes.a = 1' at line 1
> exec ^
: string
| set preceding.a = 1
< error
| 1064 42000 …near 'preceding.a = 1' at line 1
> exec ^
: string
| set precision.a = 1
< error
| 1064 42000 …near 'precision.a = 1' at line 1
> exec ^
: string
| set prefix.a = 1
< error
| 1064 42000 …near 'prefix.a = 1' at line 1
> exec ^
: string
| set prepare.a = 1
< error
| 1064 42000 …near 'prepare.a = 1' at line 1
> exec ^
: string
| set preserve.a = 1
< error
| 1064 42000 …near 'preserve.a = 1' at line 1
> exec ^
: string
| set prev.a = 1
< error
| 1064 42000 …near 'prev.a = 1' at line 1
> exec ^
: string
| set primary.a = 1
< error
| 1064 42000 …near 'primary.a = 1' at line 1
> exec ^
: string
| set privilege_checks_user.a = 1
< error
| 1064 42000 …near 'privilege_checks_user.a = 1' at line 1
> exec ^
: string
| set privileges.a = 1
< error
| 1064 42000 …near 'privileges.a = 1' at line 1
> exec ^
: string
| set procedure.a = 1
< error
| 1064 42000 …near 'procedure.a = 1' at line 1
> exec ^
: string
| set process.a = 1
< error
| 1064 42000 …near 'process.a = 1' at line 1
> exec ^
: string
| set processlist.a = 1
< error
| 1064 42000 …near 'processlist.a = 1' at line 1
> exec ^
: string
| set profile.a = 1
< error
| 1064 42000 …near 'profile.a = 1' at line 1
> exec ^
: string
| set profiles.a = 1
< error
| 1064 42000 …near 'profiles.a = 1' at line 1
> exec ^
: string
| set proxy.a = 1
< error
| 1064 42000 …near 'proxy.a = 1' at line 1
> exec ^
: string
| set purge.a = 1
< error
| 1064 42000 …near 'purge.a = 1' at line 1
> exec ^
: string
| set qualify.a = 1
< error
| 1064 42000 …near 'qualify.a = 1' at line 1
> exec ^
: string
| set quarter.a = 1
< error
| 1064 42000 …near 'quarter.a = 1' at line 1
> exec ^
: string
| set query.a = 1
< error
| 1064 42000 …near 'query.a = 1' at line 1
> exec ^
: string
| set quick.a = 1
< error
| 1064 42000 …near 'quick.a = 1' at line 1
> exec ^
: string
| set random.a = 1
< error
| 1064 42000 …near 'random.a = 1' at line 1
> exec ^
: string
| set range.a = 1
< error
| 1064 42000 …near 'range.a = 1' at line 1
> exec ^
: string
| set rank.a = 1
< error
| 1064 42000 …near 'rank.a = 1' at line 1
> exec ^
: string
| set read.a = 1
< error
| 1064 42000 …near 'read.a = 1' at line 1
> exec ^
: string
| set read_only.a = 1
< error
| 1064 42000 …near 'read_only.a = 1' at line 1
> exec ^
: string
| set read_write.a = 1
< error
| 1064 42000 …near 'read_write.a = 1' at line 1
> exec ^
: string
| set reads.a = 1
< error
| 1064 42000 …near 'reads.a = 1' at line 1
> exec ^
: string
| set real.a = 1
< error
| 1064 42000 …near 'real.a = 1' at line 1
> exec ^
: string
| set rebuild.a = 1
< error
| 1064 42000 …near 'rebuild.a = 1' at line 1
> exec ^
: string
| set recover.a = 1
< error
| 1064 42000 …near 'recover.a = 1' at line 1
> exec ^
: string
| set recursive.a = 1
< error
| 1064 42000 …near 'recursive.a = 1' at line 1
> exec ^
: string
| set redo_buffer_size.a = 1
< error
| 1064 42000 …near 'redo_buffer_size.a = 1' at line 1
> exec ^
: string
| set redundant.a = 1
< error
| 1064 42000 …near 'redundant.a = 1' at line 1
> exec ^
: string
| set reference.a = 1
< error
| 1064 42000 …near 'reference.a = 1' at line 1
> exec ^
: string
| set references.a = 1
< error
| 1064 42000 …near 'references.a = 1' at line 1
> exec ^
: string
| set regexp.a = 1
< error
| 1064 42000 …near 'regexp.a = 1' at line 1
> exec ^
: string
| set registration.a = 1
< error
| 1064 42000 …near 'registration.a = 1' at line 1
> exec ^
: string
| set relational.a = 1
< error
| 1064 42000 …near 'relational.a = 1' at line 1
> exec ^
: string
| set relay.a = 1
< error
| 1064 42000 …near 'relay.a = 1' at line 1
> exec ^
: string
| set relay_log_file.a = 1
< error
| 1064 42000 …near 'relay_log_file.a = 1' at line 1
> exec ^
: string
| set relay_log_pos.a = 1
< error
| 1064 42000 …near 'relay_log_pos.a = 1' at line 1
> exec ^
: string
| set relay_thread.a = 1
< error
| 1064 42000 …near 'relay_thread.a = 1' at line 1
> exec ^
: string
| set relaylog.a = 1
< error
| 1064 42000 …near 'relaylog.a = 1' at line 1
> exec ^
: string
| set release.a = 1
< error
| 1064 42000 …near 'release.a = 1' at line 1
> exec ^
: string
| set reload.a = 1
< error
| 1064 42000 …near 'reload.a = 1' at line 1
> exec ^
: string
| set remove.a = 1
< error
| 1064 42000 …near 'remove.a = 1' at line 1
> exec ^
: string
| set rename.a = 1
< error
| 1064 42000 …near 'rename.a = 1' at line 1
> exec ^
: string
| set reorganize.a = 1
< error
| 1064 42000 …near 'reorganize.a = 1' at line 1
> exec ^
: string
| set repair.a = 1
< error
| 1064 42000 …near 'repair.a = 1' at line 1
> exec ^
: string
| set repeat.a = 1
< error
| 1064 42000 …near 'repeat.a = 1' at line 1
> exec ^
: string
| set repeatable.a = 1
< error
| 1064 42000 …near 'repeatable.a = 1' at line 1
> exec ^
: string
| set replace.a = 1
< error
| 1064 42000 …near 'replace.a = 1' at line 1
> exec ^
: string
| set replica.a = 1
< error
| 1064 42000 …near 'replica.a = 1' at line 1
> exec ^
: string
| set replica_num.a = 1
< error
| 1064 42000 …near 'replica_num.a = 1' at line 1
> exec ^
: string
| set replicas.a = 1
< error
| 1064 42000 …near 'replicas.a = 1' at line 1
> exec ^
: string
| set replicate_do_db.a = 1
< error
| 1064 42000 …near 'replicate_do_db.a = 1' at line 1
> exec ^
: string
| set replicate_do_table.a = 1
< error
| 1064 42000 …near 'replicate_do_table.a = 1' at line 1
> exec ^
: string
| set replicate_ignore_db.a = 1
< error
| 1064 42000 …near 'replicate_ignore_db.a = 1' at line 1
> exec ^
: string
| set replicate_ignore_table.a = 1
< error
| 1064 42000 …near 'replicate_ignore_table.a = 1' at line 1
> exec ^
: string
| set replicate_rewrite_db.a = 1
< error
| 1064 42000 …near 'replicate_rewrite_db.a = 1' at line 1
> exec ^
: string
| set replicate_wild_do_table.a = 1
< error
| 1064 42000 …near 'replicate_wild_do_table.a = 1' at line 1
> exec ^
: string
| set replicate_wild_ignore_table.a = 1
< error
| 1064 42000 …near 'replicate_wild_ignore_table.a = 1' at line 1
> exec ^
: string
| set replication.a = 1
< error
| 1064 42000 …near 'replication.a = 1' at line 1
> exec ^
: string
| set require.a = 1
< error
| 1064 42000 …near 'require.a = 1' at line 1
> exec ^
: string
| set require_row_format.a = 1
< error
| 1064 42000 …near 'require_row_format.a = 1' at line 1
> exec ^
: string
| set require_table_primary_key_check.a = 1
< error
| 1064 42000 …near 'require_table_primary_key_check.a = 1' at line 1
> exec ^
: string
| set reset.a = 1
< error
| 1064 42000 …near 'reset.a = 1' at line 1
> exec ^
: string
| set resignal.a = 1
< error
| 1064 42000 …near 'resignal.a = 1' at line 1
> exec ^
: string
| set resource.a = 1
< error
| 1064 42000 …near 'resource.a = 1' at line 1
> exec ^
: string
| set respect.a = 1
< error
| 1064 42000 …near 'respect.a = 1' at line 1
> exec ^
: string
| set restart.a = 1
< error
| 1064 42000 …near 'restart.a = 1' at line 1
> exec ^
: string
| set restore.a = 1
< error
| 1064 42000 …near 'restore.a = 1' at line 1
> exec ^
: string
| set restrict.a = 1
< error
| 1064 42000 …near 'restrict.a = 1' at line 1
> exec ^
: string
| set resume.a = 1
< error
| 1064 42000 …near 'resume.a = 1' at line 1
> exec ^
: string
| set retain.a = 1
< error
| 1064 42000 …near 'retain.a = 1' at line 1
> exec ^
: string
| set return.a = 1
< error
| 1064 42000 …near 'return.a = 1' at line 1
> exec ^
: string
| set returned_sqlstate.a = 1
< error
| 1064 42000 …near 'returned_sqlstate.a = 1' at line 1
> exec ^
: string
| set returning.a = 1
< error
| 1064 42000 …near 'returning.a = 1' at line 1
> exec ^
: string
| set returns.a = 1
< error
| 1064 42000 …near 'returns.a = 1' at line 1
> exec ^
: string
| set reuse.a = 1
< error
| 1064 42000 …near 'reuse.a = 1' at line 1
> exec ^
: string
| set reverse.a = 1
< error
| 1064 42000 …near 'reverse.a = 1' at line 1
> exec ^
: string
| set revoke.a = 1
< error
| 1064 42000 …near 'revoke.a = 1' at line 1
> exec ^
: string
| set right.a = 1
< error
| 1064 42000 …near 'right.a = 1' at line 1
> exec ^
: string
| set rlike.a = 1
< error
| 1064 42000 …near 'rlike.a = 1' at line 1
> exec ^
: string
| set role.a = 1
< error
| 1064 42000 …near 'role.a = 1' at line 1
> exec ^
: string
| set rollback.a = 1
< error
| 1064 42000 …near 'rollback.a = 1' at line 1
> exec ^
: string
| set rollup.a = 1
< error
| 1064 42000 …near 'rollup.a = 1' at line 1
> exec ^
: string
| set rotate.a = 1
< error
| 1064 42000 …near 'rotate.a = 1' at line 1
> exec ^
: string
| set routine.a = 1
< error
| 1064 42000 …near 'routine.a = 1' at line 1
> exec ^
: string
| set row.a = 1
< error
| 1064 42000 …near 'row.a = 1' at line 1
> exec ^
: string
| set row_count.a = 1
< error
| 1064 42000 …near 'row_count.a = 1' at line 1
> exec ^
: string
| set row_format.a = 1
< error
| 1064 42000 …near 'row_format.a = 1' at line 1
> exec ^
: string
| set row_number.a = 1
< error
| 1064 42000 …near 'row_number.a = 1' at line 1
> exec ^
: string
| set rows.a = 1
< error
| 1064 42000 …near 'rows.a = 1' at line 1
> exec ^
: string
| set rtree.a = 1
< error
| 1064 42000 …near 'rtree.a = 1' at line 1
> exec ^
: string
| set s3.a = 1
< error
| 1064 42000 …near 's3.a = 1' at line 1
> exec ^
: string
| set savepoint.a = 1
< error
| 1064 42000 …near 'savepoint.a = 1' at line 1
> exec ^
: string
| set schedule.a = 1
< error
| 1064 42000 …near 'schedule.a = 1' at line 1
> exec ^
: string
| set schema.a = 1
< error
| 1064 42000 …near 'schema.a = 1' at line 1
> exec ^
: string
| set schema_name.a = 1
< error
| 1064 42000 …near 'schema_name.a = 1' at line 1
> exec ^
: string
| set schemas.a = 1
< error
| 1064 42000 …near 'schemas.a = 1' at line 1
> exec ^
: string
| set second.a = 1
< error
| 1064 42000 …near 'second.a = 1' at line 1
> exec ^
: string
| set second_microsecond.a = 1
< error
| 1064 42000 …near 'second_microsecond.a = 1' at line 1
> exec ^
: string
| set secondary.a = 1
< error
| 1064 42000 …near 'secondary.a = 1' at line 1
> exec ^
: string
| set secondary_engine.a = 1
< error
| 1064 42000 …near 'secondary_engine.a = 1' at line 1
> exec ^
: string
| set secondary_engine_attribute.a = 1
< error
| 1064 42000 …near 'secondary_engine_attribute.a = 1' at line 1
> exec ^
: string
| set secondary_load.a = 1
< error
| 1064 42000 …near 'secondary_load.a = 1' at line 1
> exec ^
: string
| set secondary_unload.a = 1
< error
| 1064 42000 …near 'secondary_unload.a = 1' at line 1
> exec ^
: string
| set security.a = 1
< error
| 1064 42000 …near 'security.a = 1' at line 1
> exec ^
: string
| set select.a = 1
< error
| 1064 42000 …near 'select.a = 1' at line 1
> exec ^
: string
| set sensitive.a = 1
< error
| 1064 42000 …near 'sensitive.a = 1' at line 1
> exec ^
: string
| set separator.a = 1
< error
| 1064 42000 …near 'separator.a = 1' at line 1
> exec ^
: string
| set serial.a = 1
< error
| 1064 42000 …near 'serial.a = 1' at line 1
> exec ^
: string
| set serializable.a = 1
< error
| 1064 42000 …near 'serializable.a = 1' at line 1
> exec ^
: string
| set server.a = 1
< error
| 1064 42000 …near 'server.a = 1' at line 1
> exec ^
: string
| set session.a = 1
< error
| 1064 42000 …near 'session.a = 1' at line 1
> exec ^
: string
| set session_user.a = 1
< error
| 1064 42000 …near 'session_user.a = 1' at line 1
> exec ^
: string
| set set.a = 1
< error
| 1064 42000 …near 'set.a = 1' at line 1
> exec ^
: string
| set sets.a = 1
< error
| 1064 42000 …near 'sets.a = 1' at line 1
> exec ^
: string
| set share.a = 1
< error
| 1064 42000 …near 'share.a = 1' at line 1
> exec ^
: string
| set show.a = 1
< error
| 1064 42000 …near 'show.a = 1' at line 1
> exec ^
: string
| set shutdown.a = 1
< error
| 1064 42000 …near 'shutdown.a = 1' at line 1
> exec ^
: string
| set signal.a = 1
< error
| 1064 42000 …near 'signal.a = 1' at line 1
> exec ^
: string
| set signed.a = 1
< error
| 1064 42000 …near 'signed.a = 1' at line 1
> exec ^
: string
| set simple.a = 1
< error
| 1064 42000 …near 'simple.a = 1' at line 1
> exec ^
: string
| set skip.a = 1
< error
| 1064 42000 …near 'skip.a = 1' at line 1
> exec ^
: string
| set slave.a = 1
< error
| 1064 42000 …near 'slave.a = 1' at line 1
> exec ^
: string
| set slow.a = 1
< error
| 1064 42000 …near 'slow.a = 1' at line 1
> exec ^
: string
| set smallint.a = 1
< error
| 1064 42000 …near 'smallint.a = 1' at line 1
> exec ^
: string
| set snapshot.a = 1
< error
| 1064 42000 …near 'snapshot.a = 1' at line 1
> exec ^
: string
| set socket.a = 1
< error
| 1064 42000 …near 'socket.a = 1' at line 1
> exec ^
: string
| set some.a = 1
< error
| 1064 42000 …near 'some.a = 1' at line 1
> exec ^
: string
| set soname.a = 1
< error
| 1064 42000 …near 'soname.a = 1' at line 1
> exec ^
: string
| set sounds.a = 1
< error
| 1064 42000 …near 'sounds.a = 1' at line 1
> exec ^
: string
| set source.a = 1
< error
| 1064 42000 …near 'source.a = 1' at line 1
> exec ^
: string
| set source_auto_position.a = 1
< error
| 1064 42000 …near 'source_auto_position.a = 1' at line 1
> exec ^
: string
| set source_bind.a = 1
< error
| 1064 42000 …near 'source_bind.a = 1' at line 1
> exec ^
: string
| set source_compression_algorithms.a = 1
< error
| 1064 42000 …near 'source_compression_algorithms.a = 1' at line 1
> exec ^
: string
| set source_connect_retry.a = 1
< error
| 1064 42000 …near 'source_connect_retry.a = 1' at line 1
> exec ^
: string
| set source_connection_auto_failover.a = 1
< error
| 1064 42000 …near 'source_connection_auto_failover.a = 1' at line 1
> exec ^
: string
| set source_delay.a = 1
< error
| 1064 42000 …near 'source_delay.a = 1' at line 1
> exec ^
: string
| set source_heartbeat_period.a = 1
< error
| 1064 42000 …near 'source_heartbeat_period.a = 1' at line 1
> exec ^
: string
| set source_host.a = 1
< error
| 1064 42000 …near 'source_host.a = 1' at line 1
> exec ^
: string
| set source_log_file.a = 1
< error
| 1064 42000 …near 'source_log_file.a = 1' at line 1
> exec ^
: string
| set source_log_pos.a = 1
< error
| 1064 42000 …near 'source_log_pos.a = 1' at line 1
> exec ^
: string
| set source_password.a = 1
< error
| 1064 42000 …near 'source_password.a = 1' at line 1
> exec ^
: string
| set source_port.a = 1
< error
| 1064 42000 …near 'source_port.a = 1' at line 1
> exec ^
: string
| set source_public_key_path.a = 1
< error
| 1064 42000 …near 'source_public_key_path.a = 1' at line 1
> exec ^
: string
| set source_retry_count.a = 1
< error
| 1064 42000 …near 'source_retry_count.a = 1' at line 1
> exec ^
: string
| set source_ssl.a = 1
< error
| 1064 42000 …near 'source_ssl.a = 1' at line 1
> exec ^
: string
| set source_ssl_ca.a = 1
< error
| 1064 42000 …near 'source_ssl_ca.a = 1' at line 1
> exec ^
: string
| set source_ssl_capath.a = 1
< error
| 1064 42000 …near 'source_ssl_capath.a = 1' at line 1
> exec ^
: string
| set source_ssl_cert.a = 1
< error
| 1064 42000 …near 'source_ssl_cert.a = 1' at line 1
> exec ^
: string
| set source_ssl_cipher.a = 1
< error
| 1064 42000 …near 'source_ssl_cipher.a = 1' at line 1
> exec ^
: string
| set source_ssl_crl.a = 1
< error
| 1064 42000 …near 'source_ssl_crl.a = 1' at line 1
> exec ^
: string
| set source_ssl_crlpath.a = 1
< error
| 1064 42000 …near 'source_ssl_crlpath.a = 1' at line 1
> exec ^
: string
| set source_ssl_key.a = 1
< error
| 1064 42000 …near 'source_ssl_key.a = 1' at line 1
> exec ^
: string
| set source_ssl_verify_server_cert.a = 1
< error
| 1064 42000 …near 'source_ssl_verify_server_cert.a = 1' at line 1
> exec ^
: string
| set source_tls_ciphersuites.a = 1
< error
| 1064 42000 …near 'source_tls_ciphersuites.a = 1' at line 1
> exec ^
: string
| set source_tls_version.a = 1
< error
| 1064 42000 …near 'source_tls_version.a = 1' at line 1
> exec ^
: string
| set source_user.a = 1
< error
| 1064 42000 …near 'source_user.a = 1' at line 1
> exec ^
: string
| set source_zstd_compression_level.a = 1
< error
| 1064 42000 …near 'source_zstd_compression_level.a = 1' at line 1
> exec ^
: string
| set spatial.a = 1
< error
| 1064 42000 …near 'spatial.a = 1' at line 1
> exec ^
: string
| set specific.a = 1
< error
| 1064 42000 …near 'specific.a = 1' at line 1
> exec ^
: string
| set sql.a = 1
< error
| 1064 42000 …near 'sql.a = 1' at line 1
> exec ^
: string
| set sql_after_gtids.a = 1
< error
| 1064 42000 …near 'sql_after_gtids.a = 1' at line 1
> exec ^
: string
| set sql_after_mts_gaps.a = 1
< error
| 1064 42000 …near 'sql_after_mts_gaps.a = 1' at line 1
> exec ^
: string
| set sql_before_gtids.a = 1
< error
| 1064 42000 …near 'sql_before_gtids.a = 1' at line 1
> exec ^
: string
| set sql_big_result.a = 1
< error
| 1064 42000 …near 'sql_big_result.a = 1' at line 1
> exec ^
: string
| set sql_buffer_result.a = 1
< error
| 1064 42000 …near 'sql_buffer_result.a = 1' at line 1
> exec ^
: string
| set sql_calc_found_rows.a = 1
< error
| 1064 42000 …near 'sql_calc_found_rows.a = 1' at line 1
> exec ^
: string
| set sql_no_cache.a = 1
< error
| 1064 42000 …near 'sql_no_cache.a = 1' at line 1
> exec ^
: string
| set sql_small_result.a = 1
< error
| 1064 42000 …near 'sql_small_result.a = 1' at line 1
> exec ^
: string
| set sql_thread.a = 1
< error
| 1064 42000 …near 'sql_thread.a = 1' at line 1
> exec ^
: string
| set sql_tsi_day.a = 1
< error
| 1064 42000 …near 'sql_tsi_day.a = 1' at line 1
> exec ^
: string
| set sql_tsi_hour.a = 1
< error
| 1064 42000 …near 'sql_tsi_hour.a = 1' at line 1
> exec ^
: string
| set sql_tsi_minute.a = 1
< error
| 1064 42000 …near 'sql_tsi_minute.a = 1' at line 1
> exec ^
: string
| set sql_tsi_month.a = 1
< error
| 1064 42000 …near 'sql_tsi_month.a = 1' at line 1
> exec ^
: string
| set sql_tsi_quarter.a = 1
< error
| 1064 42000 …near 'sql_tsi_quarter.a = 1' at line 1
> exec ^
: string
| set sql_tsi_second.a = 1
< error
| 1064 42000 …near 'sql_tsi_second.a = 1' at line 1
> exec ^
: string
| set sql_tsi_week.a = 1
< error
| 1064 42000 …near 'sql_tsi_week.a = 1' at line 1
> exec ^
: string
| set sql_tsi_year.a = 1
< error
| 1064 42000 …near 'sql_tsi_year.a = 1' at line 1
> exec ^
: string
| set sqlexception.a = 1
< error
| 1064 42000 …near 'sqlexception.a = 1' at line 1
> exec ^
: string
| set sqlstate.a = 1
< error
| 1064 42000 …near 'sqlstate.a = 1' at line 1
> exec ^
: string
| set sqlwarning.a = 1
< error
| 1064 42000 …near 'sqlwarning.a = 1' at line 1
> exec ^
: string
| set srid.a = 1
< error
| 1064 42000 …near 'srid.a = 1' at line 1
> exec ^
: string
| set sse_kms.a = 1
< error
| 1064 42000 …near 'sse_kms.a = 1' at line 1
> exec ^
: string
| set sse_s3.a = 1
< error
| 1064 42000 …near 'sse_s3.a = 1' at line 1
> exec ^
: string
| set ssl.a = 1
< error
| 1064 42000 …near 'ssl.a = 1' at line 1
> exec ^
: string
| set st_collect.a = 1
< error
| 1064 42000 …near 'st_collect.a = 1' at line 1
> exec ^
: string
| set stacked.a = 1
< error
| 1064 42000 …near 'stacked.a = 1' at line 1
> exec ^
: string
| set start.a = 1
< error
| 1064 42000 …near 'start.a = 1' at line 1
> exec ^
: string
| set starting.a = 1
< error
| 1064 42000 …near 'starting.a = 1' at line 1
> exec ^
: string
| set starts.a = 1
< error
| 1064 42000 …near 'starts.a = 1' at line 1
> exec ^
: string
| set stats_auto_recalc.a = 1
< error
| 1064 42000 …near 'stats_auto_recalc.a = 1' at line 1
> exec ^
: string
| set stats_persistent.a = 1
< error
| 1064 42000 …near 'stats_persistent.a = 1' at line 1
> exec ^
: string
| set stats_sample_pages.a = 1
< error
| 1064 42000 …near 'stats_sample_pages.a = 1' at line 1
> exec ^
: string
| set status.a = 1
< error
| 1064 42000 …near 'status.a = 1' at line 1
> exec ^
: string
| set std.a = 1
< error
| 1064 42000 …near 'std.a = 1' at line 1
> exec ^
: string
| set stddev.a = 1
< error
| 1064 42000 …near 'stddev.a = 1' at line 1
> exec ^
: string
| set stddev_pop.a = 1
< error
| 1064 42000 …near 'stddev_pop.a = 1' at line 1
> exec ^
: string
| set stddev_samp.a = 1
< error
| 1064 42000 …near 'stddev_samp.a = 1' at line 1
> exec ^
: string
| set stop.a = 1
< error
| 1064 42000 …near 'stop.a = 1' at line 1
> exec ^
: string
| set storage.a = 1
< error
| 1064 42000 …near 'storage.a = 1' at line 1
> exec ^
: string
| set stored.a = 1
< error
| 1064 42000 …near 'stored.a = 1' at line 1
> exec ^
: string
| set straight_join.a = 1
< error
| 1064 42000 …near 'straight_join.a = 1' at line 1
> exec ^
: string
| set stream.a = 1
< error
| 1064 42000 …near 'stream.a = 1' at line 1
> exec ^
: string
| set strict_load.a = 1
< error
| 1064 42000 …near 'strict_load.a = 1' at line 1
> exec ^
: string
| set string.a = 1
< error
| 1064 42000 …near 'string.a = 1' at line 1
> exec ^
: string
| set subclass_origin.a = 1
< error
| 1064 42000 …near 'subclass_origin.a = 1' at line 1
> exec ^
: string
| set subdate.a = 1
< error
| 1064 42000 …near 'subdate.a = 1' at line 1
> exec ^
: string
| set subject.a = 1
< error
| 1064 42000 …near 'subject.a = 1' at line 1
> exec ^
: string
| set subpartition.a = 1
< error
| 1064 42000 …near 'subpartition.a = 1' at line 1
> exec ^
: string
| set subpartitions.a = 1
< error
| 1064 42000 …near 'subpartitions.a = 1' at line 1
> exec ^
: string
| set substr.a = 1
< error
| 1064 42000 …near 'substr.a = 1' at line 1
> exec ^
: string
| set substring.a = 1
< error
| 1064 42000 …near 'substring.a = 1' at line 1
> exec ^
: string
| set sum.a = 1
< error
| 1064 42000 …near 'sum.a = 1' at line 1
> exec ^
: string
| set super.a = 1
< error
| 1064 42000 …near 'super.a = 1' at line 1
> exec ^
: string
| set suspend.a = 1
< error
| 1064 42000 …near 'suspend.a = 1' at line 1
> exec ^
: string
| set swaps.a = 1
< error
| 1064 42000 …near 'swaps.a = 1' at line 1
> exec ^
: string
| set switches.a = 1
< error
| 1064 42000 …near 'switches.a = 1' at line 1
> exec ^
: string
| set sysdate.a = 1
< error
| 1064 42000 …near 'sysdate.a = 1' at line 1
> exec ^
: string
| set system.a = 1
< error
| 1064 42000 …near 'system.a = 1' at line 1
> exec ^
: string
| set system_user.a = 1
< error
| 1064 42000 …near 'system_user.a = 1' at line 1
> exec ^
: string
| set table.a = 1
< error
| 1064 42000 …near 'table.a = 1' at line 1
> exec ^
: string
| set table_checksum.a = 1
< error
| 1064 42000 …near 'table_checksum.a = 1' at line 1
> exec ^
: string
| set table_name.a = 1
< error
| 1064 42000 …near 'table_name.a = 1' at line 1
> exec ^
: string
| set tables.a = 1
< error
| 1064 42000 …near 'tables.a = 1' at line 1
> exec ^
: string
| set tablesample.a = 1
< error
| 1064 42000 …near 'tablesample.a = 1' at line 1
> exec ^
: string
| set tablespace.a = 1
< error
| 1064 42000 …near 'tablespace.a = 1' at line 1
> exec ^
: string
| set tablet_size.a = 1
< error
| 1064 42000 …near 'tablet_size.a = 1' at line 1
> exec ^
: string
| set template.a = 1
< error
| 1064 42000 …near 'template.a = 1' at line 1
> exec ^
: string
| set temporary.a = 1
< error
| 1064 42000 …near 'temporary.a = 1' at line 1
> exec ^
: string
| set temptable.a = 1
< error
| 1064 42000 …near 'temptable.a = 1' at line 1
> exec ^
: string
| set terminated.a = 1
< error
| 1064 42000 …near 'terminated.a = 1' at line 1
> exec ^
: string
| set text.a = 1
< error
| 1064 42000 …near 'text.a = 1' at line 1
> exec ^
: string
| set than.a = 1
< error
| 1064 42000 …near 'than.a = 1' at line 1
> exec ^
: string
| set then.a = 1
< error
| 1064 42000 …near 'then.a = 1' at line 1
> exec ^
: string
| set thread_priority.a = 1
< error
| 1064 42000 …near 'thread_priority.a = 1' at line 1
> exec ^
: string
| set ties.a = 1
< error
| 1064 42000 …near 'ties.a = 1' at line 1
> exec ^
: string
| set time.a = 1
< error
| 1064 42000 …near 'time.a = 1' at line 1
> exec ^
: string
| set timestamp.a = 1
< error
| 1064 42000 …near 'timestamp.a = 1' at line 1
> exec ^
: string
| set timestampadd.a = 1
< error
| 1064 42000 …near 'timestampadd.a = 1' at line 1
> exec ^
: string
| set timestampdiff.a = 1
< error
| 1064 42000 …near 'timestampdiff.a = 1' at line 1
> exec ^
: string
| set tinyblob.a = 1
< error
| 1064 42000 …near 'tinyblob.a = 1' at line 1
> exec ^
: string
| set tinyint.a = 1
< error
| 1064 42000 …near 'tinyint.a = 1' at line 1
> exec ^
: string
| set tinytext.a = 1
< error
| 1064 42000 …near 'tinytext.a = 1' at line 1
> exec ^
: string
| set tls.a = 1
< error
| 1064 42000 …near 'tls.a = 1' at line 1
> exec ^
: string
| set to.a = 1
< error
| 1064 42000 …near 'to.a = 1' at line 1
> exec ^
: string
| set trailing.a = 1
< error
| 1064 42000 …near 'trailing.a = 1' at line 1
> exec ^
: string
| set transaction.a = 1
< error
| 1064 42000 …near 'transaction.a = 1' at line 1
> exec ^
: string
| set trigger.a = 1
< error
| 1064 42000 …near 'trigger.a = 1' at line 1
> exec ^
: string
| set triggers.a = 1
< error
| 1064 42000 …near 'triggers.a = 1' at line 1
> exec ^
: string
| set trim.a = 1
< error
| 1064 42000 …near 'trim.a = 1' at line 1
> exec ^
: string
| set true.a = 1
< error
| 1064 42000 …near 'true.a = 1' at line 1
> exec ^
: string
| set truncate.a = 1
< error
| 1064 42000 …near 'truncate.a = 1' at line 1
> exec ^
: string
| set type.a = 1
< error
| 1064 42000 …near 'type.a = 1' at line 1
> exec ^
: string
| set types.a = 1
< error
| 1064 42000 …near 'types.a = 1' at line 1
> exec ^
: string
| set unbounded.a = 1
< error
| 1064 42000 …near 'unbounded.a = 1' at line 1
> exec ^
: string
| set uncommitted.a = 1
< error
| 1064 42000 …near 'uncommitted.a = 1' at line 1
> exec ^
: string
| set undefined.a = 1
< error
| 1064 42000 …near 'undefined.a = 1' at line 1
> exec ^
: string
| set undo.a = 1
< error
| 1064 42000 …near 'undo.a = 1' at line 1
> exec ^
: string
| set undo_buffer_size.a = 1
< error
| 1064 42000 …near 'undo_buffer_size.a = 1' at line 1
> exec ^
: string
| set undofile.a = 1
< error
| 1064 42000 …near 'undofile.a = 1' at line 1
> exec ^
: string
| set unicode.a = 1
< error
| 1064 42000 …near 'unicode.a = 1' at line 1
> exec ^
: string
| set uninstall.a = 1
< error
| 1064 42000 …near 'uninstall.a = 1' at line 1
> exec ^
: string
| set union.a = 1
< error
| 1064 42000 …near 'union.a = 1' at line 1
> exec ^
: string
| set unique.a = 1
< error
| 1064 42000 …near 'unique.a = 1' at line 1
> exec ^
: string
| set unknown.a = 1
< error
| 1064 42000 …near 'unknown.a = 1' at line 1
> exec ^
: string
| set unlock.a = 1
< error
| 1064 42000 …near 'unlock.a = 1' at line 1
> exec ^
: string
| set unregister.a = 1
< error
| 1064 42000 …near 'unregister.a = 1' at line 1
> exec ^
: string
| set unsigned.a = 1
< error
| 1064 42000 …near 'unsigned.a = 1' at line 1
> exec ^
: string
| set until.a = 1
< error
| 1064 42000 …near 'until.a = 1' at line 1
> exec ^
: string
| set update.a = 1
< error
| 1064 42000 …near 'update.a = 1' at line 1
> exec ^
: string
| set upgrade.a = 1
< error
| 1064 42000 …near 'upgrade.a = 1' at line 1
> exec ^
: string
| set uri.a = 1
< error
| 1064 42000 …near 'uri.a = 1' at line 1
> exec ^
: string
| set url.a = 1
< error
| 1064 42000 …near 'url.a = 1' at line 1
> exec ^
: string
| set usage.a = 1
< error
| 1064 42000 …near 'usage.a = 1' at line 1
> exec ^
: string
| set use.a = 1
< error
| 1064 42000 …near 'use.a = 1' at line 1
> exec ^
: string
| set use_bloom_filter.a = 1
< error
| 1064 42000 …near 'use_bloom_filter.a = 1' at line 1
> exec ^
: string
| set use_frm.a = 1
< error
| 1064 42000 …near 'use_frm.a = 1' at line 1
> exec ^
: string
| set user.a = 1
< error
| 1064 42000 …near 'user.a = 1' at line 1
> exec ^
: string
| set user_resources.a = 1
< error
| 1064 42000 …near 'user_resources.a = 1' at line 1
> exec ^
: string
| set using.a = 1
< error
| 1064 42000 …near 'using.a = 1' at line 1
> exec ^
: string
| set utc_date.a = 1
< error
| 1064 42000 …near 'utc_date.a = 1' at line 1
> exec ^
: string
| set utc_time.a = 1
< error
| 1064 42000 …near 'utc_time.a = 1' at line 1
> exec ^
: string
| set utc_timestamp.a = 1
< error
| 1064 42000 …near 'utc_timestamp.a = 1' at line 1
> exec ^
: string
| set validate.a = 1
< error
| 1064 42000 …near 'validate.a = 1' at line 1
> exec ^
: string
| set validation.a = 1
< error
| 1064 42000 …near 'validation.a = 1' at line 1
> exec ^
: string
| set value.a = 1
< error
| 1064 42000 …near 'value.a = 1' at line 1
> exec ^
: string
| set values.a = 1
< error
| 1064 42000 …near 'values.a = 1' at line 1
> exec ^
: string
| set var_pop.a = 1
< error
| 1064 42000 …near 'var_pop.a = 1' at line 1
> exec ^
: string
| set var_samp.a = 1
< error
| 1064 42000 …near 'var_samp.a = 1' at line 1
> exec ^
: string
| set varbinary.a = 1
< error
| 1064 42000 …near 'varbinary.a = 1' at line 1
> exec ^
: string
| set varchar.a = 1
< error
| 1064 42000 …near 'varchar.a = 1' at line 1
> exec ^
: string
| set varcharacter.a = 1
< error
| 1064 42000 …near 'varcharacter.a = 1' at line 1
> exec ^
: string
| set variables.a = 1
< error
| 1064 42000 …near 'variables.a = 1' at line 1
> exec ^
: string
| set variance.a = 1
< error
| 1064 42000 …near 'variance.a = 1' at line 1
> exec ^
: string
| set varying.a = 1
< error
| 1064 42000 …near 'varying.a = 1' at line 1
> exec ^
: string
| set vcpu.a = 1
< error
| 1064 42000 …near 'vcpu.a = 1' at line 1
> exec ^
: string
| set vector.a = 1
< error
| 1064 42000 …near 'vector.a = 1' at line 1
> exec ^
: string
| set verify_key_constraints.a = 1
< error
| 1064 42000 …near 'verify_key_constraints.a = 1' at line 1
> exec ^
: string
| set view.a = 1
< error
| 1064 42000 …near 'view.a = 1' at line 1
> exec ^
: string
| set virtual.a = 1
< error
| 1064 42000 …near 'virtual.a = 1' at line 1
> exec ^
: string
| set visible.a = 1
< error
| 1064 42000 …near 'visible.a = 1' at line 1
> exec ^
: string
| set wait.a = 1
< error
| 1064 42000 …near 'wait.a = 1' at line 1
> exec ^
: string
| set warnings.a = 1
< error
| 1064 42000 …near 'warnings.a = 1' at line 1
> exec ^
: string
| set week.a = 1
< error
| 1064 42000 …near 'week.a = 1' at line 1
> exec ^
: string
| set weight_string.a = 1
< error
| 1064 42000 …near 'weight_string.a = 1' at line 1
> exec ^
: string
| set when.a = 1
< error
| 1064 42000 …near 'when.a = 1' at line 1
> exec ^
: string
| set where.a = 1
< error
| 1064 42000 …near 'where.a = 1' at line 1
> exec ^
: string
| set while.a = 1
< error
| 1064 42000 …near 'while.a = 1' at line 1
> exec ^
: string
| set window.a = 1
< error
| 1064 42000 …near 'window.a = 1' at line 1
> exec ^
: string
| set with.a = 1
< error
| 1064 42000 …near 'with.a = 1' at line 1
> exec ^
: string
| set without.a = 1
< error
| 1064 42000 …near 'without.a = 1' at line 1
> exec ^
: string
| set work.a = 1
< error
| 1064 42000 …near 'work.a = 1' at line 1
> exec ^
: string
| set wrapper.a = 1
< error
| 1064 42000 …near 'wrapper.a = 1' at line 1
> exec ^
: string
| set write.a = 1
< error
| 1064 42000 …near 'write.a = 1' at line 1
> exec ^
: string
| set x509.a = 1
< error
| 1064 42000 …near 'x509.a = 1' at line 1
> exec ^
: string
| set xa.a = 1
< error
| 1064 42000 …near 'xa.a = 1' at line 1
> exec ^
: string
| set xid.a = 1
< error
| 1064 42000 …near 'xid.a = 1' at line 1
> exec ^
: string
| set xml.a = 1
< error
| 1064 42000 …near 'xml.a = 1' at line 1
> exec ^
: string
| set xor.a = 1
< error
| 1064 42000 …near 'xor.a = 1' at line 1
> exec ^
: string
| set year.a = 1
< error
| 1064 42000 …near 'year.a = 1' at line 1
> exec ^
: string
| set year_month.a = 1
< error
| 1064 42000 …near 'year_month.a = 1' at line 1
> exec ^
: string
| set zerofill.a = 1
< error
| 1064 42000 …near 'zerofill.a = 1' at line 1
> exec ^
: string
| set zone.a = 1
< error
| 1064 42000 …near 'zone.a = 1' at line 1
> exec ^
: string
| set a.absent = 1
< error
| 1064 42000 …near 'a.absent = 1' at line 1
> exec ^
: string
| set a.accessible = 1
< error
| 1064 42000 …near 'a.accessible = 1' at line 1
> exec ^
: string
| set a.account = 1
< error
| 1064 42000 …near 'a.account = 1' at line 1
> exec ^
: string
| set a.action = 1
< error
| 1064 42000 …near 'a.action = 1' at line 1
> exec ^
: string
| set a.active = 1
< error
| 1064 42000 …near 'a.active = 1' at line 1
> exec ^
: string
| set a.add = 1
< error
| 1064 42000 …near 'a.add = 1' at line 1
> exec ^
: string
| set a.adddate = 1
< error
| 1064 42000 …near 'a.adddate = 1' at line 1
> exec ^
: string
| set a.admin = 1
< error
| 1064 42000 …near 'a.admin = 1' at line 1
> exec ^
: string
| set a.after = 1
< error
| 1064 42000 …near 'a.after = 1' at line 1
> exec ^
: string
| set a.against = 1
< error
| 1064 42000 …near 'a.against = 1' at line 1
> exec ^
: string
| set a.aggregate = 1
< error
| 1064 42000 …near 'a.aggregate = 1' at line 1
> exec ^
: string
| set a.algorithm = 1
< error
| 1064 42000 …near 'a.algorithm = 1' at line 1
> exec ^
: string
| set a.all = 1
< error
| 1064 42000 …near 'a.all = 1' at line 1
> exec ^
: string
| set a.allow_missing_files = 1
< error
| 1064 42000 …near 'a.allow_missing_files = 1' at line 1
> exec ^
: string
| set a.alter = 1
< error
| 1064 42000 …near 'a.alter = 1' at line 1
> exec ^
: string
| set a.always = 1
< error
| 1064 42000 …near 'a.always = 1' at line 1
> exec ^
: string
| set a.analyze = 1
< error
| 1064 42000 …near 'a.analyze = 1' at line 1
> exec ^
: string
| set a.and = 1
< error
| 1064 42000 …near 'a.and = 1' at line 1
> exec ^
: string
| set a.any = 1
< error
| 1064 42000 …near 'a.any = 1' at line 1
> exec ^
: string
| set a.array = 1
< error
| 1064 42000 …near 'a.array = 1' at line 1
> exec ^
: string
| set a.as = 1
< error
| 1064 42000 …near 'a.as = 1' at line 1
> exec ^
: string
| set a.asc = 1
< error
| 1064 42000 …near 'a.asc = 1' at line 1
> exec ^
: string
| set a.ascii = 1
< error
| 1064 42000 …near 'a.ascii = 1' at line 1
> exec ^
: string
| set a.asensitive = 1
< error
| 1064 42000 …near 'a.asensitive = 1' at line 1
> exec ^
: string
| set a.assign_gtids_to_anonymous_transactions = 1
< error
| 1064 42000 …near 'a.assign_gtids_to_anonymous_transactions = 1' at line 1
> exec ^
: string
| set a.at = 1
< error
| 1064 42000 …near 'a.at = 1' at line 1
> exec ^
: string
| set a.attribute = 1
< error
| 1064 42000 …near 'a.attribute = 1' at line 1
> exec ^
: string
| set a.authentication = 1
< error
| 1064 42000 …near 'a.authentication = 1' at line 1
> exec ^
: string
| set a.auto = 1
< error
| 1064 42000 …near 'a.auto = 1' at line 1
> exec ^
: string
| set a.auto_increment = 1
< error
| 1064 42000 …near 'a.auto_increment = 1' at line 1
> exec ^
: string
| set a.auto_increment_mode = 1
< error
| 1064 42000 …near 'a.auto_increment_mode = 1' at line 1
> exec ^
: string
| set a.auto_refresh = 1
< error
| 1064 42000 …near 'a.auto_refresh = 1' at line 1
> exec ^
: string
| set a.auto_refresh_source = 1
< error
| 1064 42000 …near 'a.auto_refresh_source = 1' at line 1
> exec ^
: string
| set a.autoextend_size = 1
< error
| 1064 42000 …near 'a.autoextend_size = 1' at line 1
> exec ^
: string
| set a.avg = 1
< error
| 1064 42000 …near 'a.avg = 1' at line 1
> exec ^
: string
| set a.avg_row_length = 1
< error
| 1064 42000 …near 'a.avg_row_length = 1' at line 1
> exec ^
: string
| set a.backup = 1
< error
| 1064 42000 …near 'a.backup = 1' at line 1
> exec ^
: string
| set a.before = 1
< error
| 1064 42000 …near 'a.before = 1' at line 1
> exec ^
: string
| set a.begin = 1
< error
| 1064 42000 …near 'a.begin = 1' at line 1
> exec ^
: string
| set a.bernoulli = 1
< error
| 1064 42000 …near 'a.bernoulli = 1' at line 1
> exec ^
: string
| set a.between = 1
< error
| 1064 42000 …near 'a.between = 1' at line 1
> exec ^
: string
| set a.bigint = 1
< error
| 1064 42000 …near 'a.bigint = 1' at line 1
> exec ^
: string
| set a.binary = 1
< error
| 1064 42000 …near 'a.binary = 1' at line 1
> exec ^
: string
| set a.binlog = 1
< error
| 1064 42000 …near 'a.binlog = 1' at line 1
> exec ^
: string
| set a.bit = 1
< error
| 1064 42000 …near 'a.bit = 1' at line 1
> exec ^
: string
| set a.bit_and = 1
< error
| 1064 42000 …near 'a.bit_and = 1' at line 1
> exec ^
: string
| set a.bit_or = 1
< error
| 1064 42000 …near 'a.bit_or = 1' at line 1
> exec ^
: string
| set a.bit_xor = 1
< error
| 1064 42000 …near 'a.bit_xor = 1' at line 1
> exec ^
: string
| set a.blob = 1
< error
| 1064 42000 …near 'a.blob = 1' at line 1
> exec ^
: string
| set a.block = 1
< error
| 1064 42000 …near 'a.block = 1' at line 1
> exec ^
: string
| set a.block_size = 1
< error
| 1064 42000 …near 'a.block_size = 1' at line 1
> exec ^
: string
| set a.bool = 1
< error
| 1064 42000 …near 'a.bool = 1' at line 1
> exec ^
: string
| set a.boolean = 1
< error
| 1064 42000 …near 'a.boolean = 1' at line 1
> exec ^
: string
| set a.both = 1
< error
| 1064 42000 …near 'a.both = 1' at line 1
> exec ^
: string
| set a.btree = 1
< error
| 1064 42000 …near 'a.btree = 1' at line 1
> exec ^
: string
| set a.buckets = 1
< error
| 1064 42000 …near 'a.buckets = 1' at line 1
> exec ^
: string
| set a.bulk = 1
< error
| 1064 42000 …near 'a.bulk = 1' at line 1
> exec ^
: string
| set a.by = 1
< error
| 1064 42000 …near 'a.by = 1' at line 1
> exec ^
: string
| set a.byte = 1
< error
| 1064 42000 …near 'a.byte = 1' at line 1
> exec ^
: string
| set a.cache = 1
< error
| 1064 42000 …near 'a.cache = 1' at line 1
> exec ^
: string
| set a.call = 1
< error
| 1064 42000 …near 'a.call = 1' at line 1
> exec ^
: string
| set a.cascade = 1
< error
| 1064 42000 …near 'a.cascade = 1' at line 1
> exec ^
: string
| set a.cascaded = 1
< error
| 1064 42000 …near 'a.cascaded = 1' at line 1
> exec ^
: string
| set a.case = 1
< error
| 1064 42000 …near 'a.case = 1' at line 1
> exec ^
: string
| set a.cast = 1
< error
| 1064 42000 …near 'a.cast = 1' at line 1
> exec ^
: string
| set a.catalog_name = 1
< error
| 1064 42000 …near 'a.catalog_name = 1' at line 1
> exec ^
: string
| set a.chain = 1
< error
| 1064 42000 …near 'a.chain = 1' at line 1
> exec ^
: string
| set a.challenge_response = 1
< error
| 1064 42000 …near 'a.challenge_response = 1' at line 1
> exec ^
: string
| set a.change = 1
< error
| 1064 42000 …near 'a.change = 1' at line 1
> exec ^
: string
| set a.changed = 1
< error
| 1064 42000 …near 'a.changed = 1' at line 1
> exec ^
: string
| set a.channel = 1
< error
| 1064 42000 …near 'a.channel = 1' at line 1
> exec ^
: string
| set a.char = 1
< error
| 1064 42000 …near 'a.char = 1' at line 1
> exec ^
: string
| set a.character = 1
< error
| 1064 42000 …near 'a.character = 1' at line 1
> exec ^
: string
| set a.charset = 1
< error
| 1064 42000 …near 'a.charset = 1' at line 1
> exec ^
: string
| set a.check = 1
< error
| 1064 42000 …near 'a.check = 1' at line 1
> exec ^
: string
| set a.checksum = 1
< error
| 1064 42000 …near 'a.checksum = 1' at line 1
> exec ^
: string
| set a.cipher = 1
< error
| 1064 42000 …near 'a.cipher = 1' at line 1
> exec ^
: string
| set a.class_origin = 1
< error
| 1064 42000 …near 'a.class_origin = 1' at line 1
> exec ^
: string
| set a.client = 1
< error
| 1064 42000 …near 'a.client = 1' at line 1
> exec ^
: string
| set a.clone = 1
< error
| 1064 42000 …near 'a.clone = 1' at line 1
> exec ^
: string
| set a.close = 1
< error
| 1064 42000 …near 'a.close = 1' at line 1
> exec ^
: string
| set a.coalesce = 1
< error
| 1064 42000 …near 'a.coalesce = 1' at line 1
> exec ^
: string
| set a.code = 1
< error
| 1064 42000 …near 'a.code = 1' at line 1
> exec ^
: string
| set a.collate = 1
< error
| 1064 42000 …near 'a.collate = 1' at line 1
> exec ^
: string
| set a.collation = 1
< error
| 1064 42000 …near 'a.collation = 1' at line 1
> exec ^
: string
| set a.column = 1
< error
| 1064 42000 …near 'a.column = 1' at line 1
> exec ^
: string
| set a.column_format = 1
< error
| 1064 42000 …near 'a.column_format = 1' at line 1
> exec ^
: string
| set a.column_name = 1
< error
| 1064 42000 …near 'a.column_name = 1' at line 1
> exec ^
: string
| set a.columns = 1
< error
| 1064 42000 …near 'a.columns = 1' at line 1
> exec ^
: string
| set a.comment = 1
< error
| 1064 42000 …near 'a.comment = 1' at line 1
> exec ^
: string
| set a.commit = 1
< error
| 1064 42000 …near 'a.commit = 1' at line 1
> exec ^
: string
| set a.committed = 1
< error
| 1064 42000 …near 'a.committed = 1' at line 1
> exec ^
: string
| set a.compact = 1
< error
| 1064 42000 …near 'a.compact = 1' at line 1
> exec ^
: string
| set a.completion = 1
< error
| 1064 42000 …near 'a.completion = 1' at line 1
> exec ^
: string
| set a.component = 1
< error
| 1064 42000 …near 'a.component = 1' at line 1
> exec ^
: string
| set a.compressed = 1
< error
| 1064 42000 …near 'a.compressed = 1' at line 1
> exec ^
: string
| set a.compression = 1
< error
| 1064 42000 …near 'a.compression = 1' at line 1
> exec ^
: string
| set a.concat = 1
< error
| 1064 42000 …near 'a.concat = 1' at line 1
> exec ^
: string
| set a.concurrent = 1
< error
| 1064 42000 …near 'a.concurrent = 1' at line 1
> exec ^
: string
| set a.condition = 1
< error
| 1064 42000 …near 'a.condition = 1' at line 1
> exec ^
: string
| set a.connection = 1
< error
| 1064 42000 …near 'a.connection = 1' at line 1
> exec ^
: string
| set a.consistent = 1
< error
| 1064 42000 …near 'a.consistent = 1' at line 1
> exec ^
: string
| set a.constraint = 1
< error
| 1064 42000 …near 'a.constraint = 1' at line 1
> exec ^
: string
| set a.constraint_catalog = 1
< error
| 1064 42000 …near 'a.constraint_catalog = 1' at line 1
> exec ^
: string
| set a.constraint_name = 1
< error
| 1064 42000 …near 'a.constraint_name = 1' at line 1
> exec ^
: string
| set a.constraint_schema = 1
< error
| 1064 42000 …near 'a.constraint_schema = 1' at line 1
> exec ^
: string
| set a.contains = 1
< error
| 1064 42000 …near 'a.contains = 1' at line 1
> exec ^
: string
| set a.context = 1
< error
| 1064 42000 …near 'a.context = 1' at line 1
> exec ^
: string
| set a.continue = 1
< error
| 1064 42000 …near 'a.continue = 1' at line 1
> exec ^
: string
| set a.convert = 1
< error
| 1064 42000 …near 'a.convert = 1' at line 1
> exec ^
: string
| set a.count = 1
< error
| 1064 42000 …near 'a.count = 1' at line 1
> exec ^
: string
| set a.cpu = 1
< error
| 1064 42000 …near 'a.cpu = 1' at line 1
> exec ^
: string
| set a.create = 1
< error
| 1064 42000 …near 'a.create = 1' at line 1
> exec ^
: string
| set a.cross = 1
< error
| 1064 42000 …near 'a.cross = 1' at line 1
> exec ^
: string
| set a.csv = 1
< error
| 1064 42000 …near 'a.csv = 1' at line 1
> exec ^
: string
| set a.cube = 1
< error
| 1064 42000 …near 'a.cube = 1' at line 1
> exec ^
: string
| set a.cume_dist = 1
< error
| 1064 42000 …near 'a.cume_dist = 1' at line 1
> exec ^
: string
| set a.curdate = 1
< error
| 1064 42000 …near 'a.curdate = 1' at line 1
> exec ^
: string
| set a.current = 1
< error
| 1064 42000 …near 'a.current = 1' at line 1
> exec ^
: string
| set a.current_date = 1
< error
| 1064 42000 …near 'a.current_date = 1' at line 1
> exec ^
: string
| set a.current_time = 1
< error
| 1064 42000 …near 'a.current_time = 1' at line 1
> exec ^
: string
| set a.current_timestamp = 1
< error
| 1064 42000 …near 'a.current_timestamp = 1' at line 1
> exec ^
: string
| set a.current_user = 1
< error
| 1064 42000 …near 'a.current_user = 1' at line 1
> exec ^
: string
| set a.cursor = 1
< error
| 1064 42000 …near 'a.cursor = 1' at line 1
> exec ^
: string
| set a.cursor_name = 1
< error
| 1064 42000 …near 'a.cursor_name = 1' at line 1
> exec ^
: string
| set a.curtime = 1
< error
| 1064 42000 …near 'a.curtime = 1' at line 1
> exec ^
: string
| set a.data = 1
< error
| 1064 42000 …near 'a.data = 1' at line 1
> exec ^
: string
| set a.database = 1
< error
| 1064 42000 …near 'a.database = 1' at line 1
> exec ^
: string
| set a.databases = 1
< error
| 1064 42000 …near 'a.databases = 1' at line 1
> exec ^
: string
| set a.datafile = 1
< error
| 1064 42000 …near 'a.datafile = 1' at line 1
> exec ^
: string
| set a.date = 1
< error
| 1064 42000 …near 'a.date = 1' at line 1
> exec ^
: string
| set a.date_add = 1
< error
| 1064 42000 …near 'a.date_add = 1' at line 1
> exec ^
: string
| set a.date_sub = 1
< error
| 1064 42000 …near 'a.date_sub = 1' at line 1
> exec ^
: string
| set a.datetime = 1
< error
| 1064 42000 …near 'a.datetime = 1' at line 1
> exec ^
: string
| set a.day = 1
< error
| 1064 42000 …near 'a.day = 1' at line 1
> exec ^
: string
| set a.day_hour = 1
< error
| 1064 42000 …near 'a.day_hour = 1' at line 1
> exec ^
: string
| set a.day_microsecond = 1
< error
| 1064 42000 …near 'a.day_microsecond = 1' at line 1
> exec ^
: string
| set a.day_minute = 1
< error
| 1064 42000 …near 'a.day_minute = 1' at line 1
> exec ^
: string
| set a.day_second = 1
< error
| 1064 42000 …near 'a.day_second = 1' at line 1
> exec ^
: string
| set a.deallocate = 1
< error
| 1064 42000 …near 'a.deallocate = 1' at line 1
> exec ^
: string
| set a.dec = 1
< error
| 1064 42000 …near 'a.dec = 1' at line 1
> exec ^
: string
| set a.decimal = 1
< error
| 1064 42000 …near 'a.decimal = 1' at line 1
> exec ^
: string
| set a.declare = 1
< error
| 1064 42000 …near 'a.declare = 1' at line 1
> exec ^
: string
| set a.default = 1
< error
| 1064 42000 …near 'a.default = 1' at line 1
> exec ^
: string
| set a.default_auth = 1
< error
| 1064 42000 …near 'a.default_auth = 1' at line 1
> exec ^
: string
| set a.definer = 1
< error
| 1064 42000 …near 'a.definer = 1' at line 1
> exec ^
: string
| set a.definition = 1
< error
| 1064 42000 …near 'a.definition = 1' at line 1
> exec ^
: string
| set a.delay_key_write = 1
< error
| 1064 42000 …near 'a.delay_key_write = 1' at line 1
> exec ^
: string
| set a.delayed = 1
< error
| 1064 42000 …near 'a.delayed = 1' at line 1
> exec ^
: string
| set a.delete = 1
< error
| 1064 42000 …near 'a.delete = 1' at line 1
> exec ^
: string
| set a.dense_rank = 1
< error
| 1064 42000 …near 'a.dense_rank = 1' at line 1
> exec ^
: string
| set a.desc = 1
< error
| 1064 42000 …near 'a.desc = 1' at line 1
> exec ^
: string
| set a.describe = 1
< error
| 1064 42000 …near 'a.describe = 1' at line 1
> exec ^
: string
| set a.description = 1
< error
| 1064 42000 …near 'a.description = 1' at line 1
> exec ^
: string
| set a.deterministic = 1
< error
| 1064 42000 …near 'a.deterministic = 1' at line 1
> exec ^
: string
| set a.diagnostics = 1
< error
| 1064 42000 …near 'a.diagnostics = 1' at line 1
> exec ^
: string
| set a.directory = 1
< error
| 1064 42000 …near 'a.directory = 1' at line 1
> exec ^
: string
| set a.disable = 1
< error
| 1064 42000 …near 'a.disable = 1' at line 1
> exec ^
: string
| set a.discard = 1
< error
| 1064 42000 …near 'a.discard = 1' at line 1
> exec ^
: string
| set a.disk = 1
< error
| 1064 42000 …near 'a.disk = 1' at line 1
> exec ^
: string
| set a.distinct = 1
< error
| 1064 42000 …near 'a.distinct = 1' at line 1
> exec ^
: string
| set a.distinctrow = 1
< error
| 1064 42000 …near 'a.distinctrow = 1' at line 1
> exec ^
: string
| set a.div = 1
< error
| 1064 42000 …near 'a.div = 1' at line 1
> exec ^
: string
| set a.do = 1
< error
| 1064 42000 …near 'a.do = 1' at line 1
> exec ^
: string
| set a.double = 1
< error
| 1064 42000 …near 'a.double = 1' at line 1
> exec ^
: string
| set a.drop = 1
< error
| 1064 42000 …near 'a.drop = 1' at line 1
> exec ^
: string
| set a.dual = 1
< error
| 1064 42000 …near 'a.dual = 1' at line 1
> exec ^
: string
| set a.duality = 1
< error
| 1064 42000 …near 'a.duality = 1' at line 1
> exec ^
: string
| set a.dumpfile = 1
< error
| 1064 42000 …near 'a.dumpfile = 1' at line 1
> exec ^
: string
| set a.duplicate = 1
< error
| 1064 42000 …near 'a.duplicate = 1' at line 1
> exec ^
: string
| set a.dynamic = 1
< error
| 1064 42000 …near 'a.dynamic = 1' at line 1
> exec ^
: string
| set a.dynamic_partition_policy = 1
< error
| 1064 42000 …near 'a.dynamic_partition_policy = 1' at line 1
> exec ^
: string
| set a.each = 1
< error
| 1064 42000 …near 'a.each = 1' at line 1
> exec ^
: string
| set a.else = 1
< error
| 1064 42000 …near 'a.else = 1' at line 1
> exec ^
: string
| set a.elseif = 1
< error
| 1064 42000 …near 'a.elseif = 1' at line 1
> exec ^
: string
| set a.empty = 1
< error
| 1064 42000 …near 'a.empty = 1' at line 1
> exec ^
: string
| set a.enable = 1
< error
| 1064 42000 …near 'a.enable = 1' at line 1
> exec ^
: string
| set a.enable_macro_block_bloom_filter = 1
< error
| 1064 42000 …near 'a.enable_macro_block_bloom_filter = 1' at line 1
> exec ^
: string
| set a.enclosed = 1
< error
| 1064 42000 …near 'a.enclosed = 1' at line 1
> exec ^
: string
| set a.encryption = 1
< error
| 1064 42000 …near 'a.encryption = 1' at line 1
> exec ^
: string
| set a.end = 1
< error
| 1064 42000 …near 'a.end = 1' at line 1
> exec ^
: string
| set a.ends = 1
< error
| 1064 42000 …near 'a.ends = 1' at line 1
> exec ^
: string
| set a.enforced = 1
< error
| 1064 42000 …near 'a.enforced = 1' at line 1
> exec ^
: string
| set a.engine = 1
< error
| 1064 42000 …near 'a.engine = 1' at line 1
> exec ^
: string
| set a.engine_attribute = 1
< error
| 1064 42000 …near 'a.engine_attribute = 1' at line 1
> exec ^
: string
| set a.engines = 1
< error
| 1064 42000 …near 'a.engines = 1' at line 1
> exec ^
: string
| set a.enum = 1
< error
| 1064 42000 …near 'a.enum = 1' at line 1
> exec ^
: string
| set a.error = 1
< error
| 1064 42000 …near 'a.error = 1' at line 1
> exec ^
: string
| set a.errors = 1
< error
| 1064 42000 …near 'a.errors = 1' at line 1
> exec ^
: string
| set a.escape = 1
< error
| 1064 42000 …near 'a.escape = 1' at line 1
> exec ^
: string
| set a.escaped = 1
< error
| 1064 42000 …near 'a.escaped = 1' at line 1
> exec ^
: string
| set a.event = 1
< error
| 1064 42000 …near 'a.event = 1' at line 1
> exec ^
: string
| set a.events = 1
< error
| 1064 42000 …near 'a.events = 1' at line 1
> exec ^
: string
| set a.every = 1
< error
| 1064 42000 …near 'a.every = 1' at line 1
> exec ^
: string
| set a.except = 1
< error
| 1064 42000 …near 'a.except = 1' at line 1
> exec ^
: string
| set a.exchange = 1
< error
| 1064 42000 …near 'a.exchange = 1' at line 1
> exec ^
: string
| set a.exclude = 1
< error
| 1064 42000 …near 'a.exclude = 1' at line 1
> exec ^
: string
| set a.execute = 1
< error
| 1064 42000 …near 'a.execute = 1' at line 1
> exec ^
: string
| set a.exists = 1
< error
| 1064 42000 …near 'a.exists = 1' at line 1
> exec ^
: string
| set a.exit = 1
< error
| 1064 42000 …near 'a.exit = 1' at line 1
> exec ^
: string
| set a.expansion = 1
< error
| 1064 42000 …near 'a.expansion = 1' at line 1
> exec ^
: string
| set a.expire = 1
< error
| 1064 42000 …near 'a.expire = 1' at line 1
> exec ^
: string
| set a.explain = 1
< error
| 1064 42000 …near 'a.explain = 1' at line 1
> exec ^
: string
| set a.export = 1
< error
| 1064 42000 …near 'a.export = 1' at line 1
> exec ^
: string
| set a.extended = 1
< error
| 1064 42000 …near 'a.extended = 1' at line 1
> exec ^
: string
| set a.extent_size = 1
< error
| 1064 42000 …near 'a.extent_size = 1' at line 1
> exec ^
: string
| set a.external = 1
< error
| 1064 42000 …near 'a.external = 1' at line 1
> exec ^
: string
| set a.external_format = 1
< error
| 1064 42000 …near 'a.external_format = 1' at line 1
> exec ^
: string
| set a.extract = 1
< error
| 1064 42000 …near 'a.extract = 1' at line 1
> exec ^
: string
| set a.factor = 1
< error
| 1064 42000 …near 'a.factor = 1' at line 1
> exec ^
: string
| set a.failed_login_attempts = 1
< error
| 1064 42000 …near 'a.failed_login_attempts = 1' at line 1
> exec ^
: string
| set a.false = 1
< error
| 1064 42000 …near 'a.false = 1' at line 1
> exec ^
: string
| set a.fast = 1
< error
| 1064 42000 …near 'a.fast = 1' at line 1
> exec ^
: string
| set a.faults = 1
< error
| 1064 42000 …near 'a.faults = 1' at line 1
> exec ^
: string
| set a.fetch = 1
< error
| 1064 42000 …near 'a.fetch = 1' at line 1
> exec ^
: string
| set a.fields = 1
< error
| 1064 42000 …near 'a.fields = 1' at line 1
> exec ^
: string
| set a.file = 1
< error
| 1064 42000 …near 'a.file = 1' at line 1
> exec ^
: string
| set a.file_block_size = 1
< error
| 1064 42000 …near 'a.file_block_size = 1' at line 1
> exec ^
: string
| set a.file_format = 1
< error
| 1064 42000 …near 'a.file_format = 1' at line 1
> exec ^
: string
| set a.file_name = 1
< error
| 1064 42000 …near 'a.file_name = 1' at line 1
> exec ^
: string
| set a.file_pattern = 1
< error
| 1064 42000 …near 'a.file_pattern = 1' at line 1
> exec ^
: string
| set a.file_prefix = 1
< error
| 1064 42000 …near 'a.file_prefix = 1' at line 1
> exec ^
: string
| set a.files = 1
< error
| 1064 42000 …near 'a.files = 1' at line 1
> exec ^
: string
| set a.filter = 1
< error
| 1064 42000 …near 'a.filter = 1' at line 1
> exec ^
: string
| set a.finish = 1
< error
| 1064 42000 …near 'a.finish = 1' at line 1
> exec ^
: string
| set a.first = 1
< error
| 1064 42000 …near 'a.first = 1' at line 1
> exec ^
: string
| set a.first_value = 1
< error
| 1064 42000 …near 'a.first_value = 1' at line 1
> exec ^
: string
| set a.fixed = 1
< error
| 1064 42000 …near 'a.fixed = 1' at line 1
> exec ^
: string
| set a.float = 1
< error
| 1064 42000 …near 'a.float = 1' at line 1
> exec ^
: string
| set a.float4 = 1
< error
| 1064 42000 …near 'a.float4 = 1' at line 1
> exec ^
: string
| set a.float8 = 1
< error
| 1064 42000 …near 'a.float8 = 1' at line 1
> exec ^
: string
| set a.flush = 1
< error
| 1064 42000 …near 'a.flush = 1' at line 1
> exec ^
: string
| set a.following = 1
< error
| 1064 42000 …near 'a.following = 1' at line 1
> exec ^
: string
| set a.follows = 1
< error
| 1064 42000 …near 'a.follows = 1' at line 1
> exec ^
: string
| set a.for = 1
< error
| 1064 42000 …near 'a.for = 1' at line 1
> exec ^
: string
| set a.force = 1
< error
| 1064 42000 …near 'a.force = 1' at line 1
> exec ^
: string
| set a.foreign = 1
< error
| 1064 42000 …near 'a.foreign = 1' at line 1
> exec ^
: string
| set a.format = 1
< error
| 1064 42000 …near 'a.format = 1' at line 1
> exec ^
: string
| set a.found = 1
< error
| 1064 42000 …near 'a.found = 1' at line 1
> exec ^
: string
| set a.from = 1
< error
| 1064 42000 …near 'a.from = 1' at line 1
> exec ^
: string
| set a.full = 1
< error
| 1064 42000 …near 'a.full = 1' at line 1
> exec ^
: string
| set a.fulltext = 1
< error
| 1064 42000 …near 'a.fulltext = 1' at line 1
> exec ^
: string
| set a.function = 1
< error
| 1064 42000 …near 'a.function = 1' at line 1
> exec ^
: string
| set a.general = 1
< error
| 1064 42000 …near 'a.general = 1' at line 1
> exec ^
: string
| set a.generate = 1
< error
| 1064 42000 …near 'a.generate = 1' at line 1
> exec ^
: string
| set a.generated = 1
< error
| 1064 42000 …near 'a.generated = 1' at line 1
> exec ^
: string
| set a.geomcollection = 1
< error
| 1064 42000 …near 'a.geomcollection = 1' at line 1
> exec ^
: string
| set a.geometry = 1
< error
| 1064 42000 …near 'a.geometry = 1' at line 1
> exec ^
: string
| set a.geometrycollection = 1
< error
| 1064 42000 …near 'a.geometrycollection = 1' at line 1
> exec ^
: string
| set a.get = 1
< error
| 1064 42000 …near 'a.get = 1' at line 1
> exec ^
: string
| set a.get_format = 1
< error
| 1064 42000 …near 'a.get_format = 1' at line 1
> exec ^
: string
| set a.get_master_public_key = 1
< error
| 1064 42000 …near 'a.get_master_public_key = 1' at line 1
> exec ^
: string
| set a.get_source_public_key = 1
< error
| 1064 42000 …near 'a.get_source_public_key = 1' at line 1
> exec ^
: string
| set a.global = 1
< error
| 1064 42000 …near 'a.global = 1' at line 1
> exec ^
: string
| set a.grant = 1
< error
| 1064 42000 …near 'a.grant = 1' at line 1
> exec ^
: string
| set a.grants = 1
< error
| 1064 42000 …near 'a.grants = 1' at line 1
> exec ^
: string
| set a.greatest = 1
< error
| 1064 42000 …near 'a.greatest = 1' at line 1
> exec ^
: string
| set a.group = 1
< error
| 1064 42000 …near 'a.group = 1' at line 1
> exec ^
: string
| set a.group_concat = 1
< error
| 1064 42000 …near 'a.group_concat = 1' at line 1
> exec ^
: string
| set a.group_replication = 1
< error
| 1064 42000 …near 'a.group_replication = 1' at line 1
> exec ^
: string
| set a.grouping = 1
< error
| 1064 42000 …near 'a.grouping = 1' at line 1
> exec ^
: string
| set a.groups = 1
< error
| 1064 42000 …near 'a.groups = 1' at line 1
> exec ^
: string
| set a.gtid_only = 1
< error
| 1064 42000 …near 'a.gtid_only = 1' at line 1
> exec ^
: string
| set a.gtids = 1
< error
| 1064 42000 …near 'a.gtids = 1' at line 1
> exec ^
: string
| set a.guided = 1
< error
| 1064 42000 …near 'a.guided = 1' at line 1
> exec ^
: string
| set a.handler = 1
< error
| 1064 42000 …near 'a.handler = 1' at line 1
> exec ^
: string
| set a.hash = 1
< error
| 1064 42000 …near 'a.hash = 1' at line 1
> exec ^
: string
| set a.having = 1
< error
| 1064 42000 …near 'a.having = 1' at line 1
> exec ^
: string
| set a.header = 1
< error
| 1064 42000 …near 'a.header = 1' at line 1
> exec ^
: string
| set a.help = 1
< error
| 1064 42000 …near 'a.help = 1' at line 1
> exec ^
: string
| set a.high_priority = 1
< error
| 1064 42000 …near 'a.high_priority = 1' at line 1
> exec ^
: string
| set a.histogram = 1
< error
| 1064 42000 …near 'a.histogram = 1' at line 1
> exec ^
: string
| set a.history = 1
< error
| 1064 42000 …near 'a.history = 1' at line 1
> exec ^
: string
| set a.host = 1
< error
| 1064 42000 …near 'a.host = 1' at line 1
> exec ^
: string
| set a.hosts = 1
< error
| 1064 42000 …near 'a.hosts = 1' at line 1
> exec ^
: string
| set a.hour = 1
< error
| 1064 42000 …near 'a.hour = 1' at line 1
> exec ^
: string
| set a.hour_microsecond = 1
< error
| 1064 42000 …near 'a.hour_microsecond = 1' at line 1
> exec ^
: string
| set a.hour_minute = 1
< error
| 1064 42000 …near 'a.hour_minute = 1' at line 1
> exec ^
: string
| set a.hour_second = 1
< error
| 1064 42000 …near 'a.hour_second = 1' at line 1
> exec ^
: string
| set a.identified = 1
< error
| 1064 42000 …near 'a.identified = 1' at line 1
> exec ^
: string
| set a.if = 1
< error
| 1064 42000 …near 'a.if = 1' at line 1
> exec ^
: string
| set a.ignore = 1
< error
| 1064 42000 …near 'a.ignore = 1' at line 1
> exec ^
: string
| set a.ignore_server_ids = 1
< error
| 1064 42000 …near 'a.ignore_server_ids = 1' at line 1
> exec ^
: string
| set a.import = 1
< error
| 1064 42000 …near 'a.import = 1' at line 1
> exec ^
: string
| set a.in = 1
< error
| 1064 42000 …near 'a.in = 1' at line 1
> exec ^
: string
| set a.inactive = 1
< error
| 1064 42000 …near 'a.inactive = 1' at line 1
> exec ^
: string
| set a.index = 1
< error
| 1064 42000 …near 'a.index = 1' at line 1
> exec ^
: string
| set a.indexes = 1
< error
| 1064 42000 …near 'a.indexes = 1' at line 1
> exec ^
: string
| set a.infile = 1
< error
| 1064 42000 …near 'a.infile = 1' at line 1
> exec ^
: string
| set a.initial = 1
< error
| 1064 42000 …near 'a.initial = 1' at line 1
> exec ^
: string
| set a.initial_size = 1
< error
| 1064 42000 …near 'a.initial_size = 1' at line 1
> exec ^
: string
| set a.initiate = 1
< error
| 1064 42000 …near 'a.initiate = 1' at line 1
> exec ^
: string
| set a.inner = 1
< error
| 1064 42000 …near 'a.inner = 1' at line 1
> exec ^
: string
| set a.inout = 1
< error
| 1064 42000 …near 'a.inout = 1' at line 1
> exec ^
: string
| set a.insensitive = 1
< error
| 1064 42000 …near 'a.insensitive = 1' at line 1
> exec ^
: string
| set a.insert = 1
< error
| 1064 42000 …near 'a.insert = 1' at line 1
> exec ^
: string
| set a.insert_method = 1
< error
| 1064 42000 …near 'a.insert_method = 1' at line 1
> exec ^
: string
| set a.install = 1
< error
| 1064 42000 …near 'a.install = 1' at line 1
> exec ^
: string
| set a.instance = 1
< error
| 1064 42000 …near 'a.instance = 1' at line 1
> exec ^
: string
| set a.int = 1
< error
| 1064 42000 …near 'a.int = 1' at line 1
> exec ^
: string
| set a.int1 = 1
< error
| 1064 42000 …near 'a.int1 = 1' at line 1
> exec ^
: string
| set a.int2 = 1
< error
| 1064 42000 …near 'a.int2 = 1' at line 1
> exec ^
: string
| set a.int3 = 1
< error
| 1064 42000 …near 'a.int3 = 1' at line 1
> exec ^
: string
| set a.int4 = 1
< error
| 1064 42000 …near 'a.int4 = 1' at line 1
> exec ^
: string
| set a.int8 = 1
< error
| 1064 42000 …near 'a.int8 = 1' at line 1
> exec ^
: string
| set a.integer = 1
< error
| 1064 42000 …near 'a.integer = 1' at line 1
> exec ^
: string
| set a.intersect = 1
< error
| 1064 42000 …near 'a.intersect = 1' at line 1
> exec ^
: string
| set a.interval = 1
< error
| 1064 42000 …near 'a.interval = 1' at line 1
> exec ^
: string
| set a.into = 1
< error
| 1064 42000 …near 'a.into = 1' at line 1
> exec ^
: string
| set a.invisible = 1
< error
| 1064 42000 …near 'a.invisible = 1' at line 1
> exec ^
: string
| set a.invoker = 1
< error
| 1064 42000 …near 'a.invoker = 1' at line 1
> exec ^
: string
| set a.io = 1
< error
| 1064 42000 …near 'a.io = 1' at line 1
> exec ^
: string
| set a.io_after_gtids = 1
< error
| 1064 42000 …near 'a.io_after_gtids = 1' at line 1
> exec ^
: string
| set a.io_before_gtids = 1
< error
| 1064 42000 …near 'a.io_before_gtids = 1' at line 1
> exec ^
: string
| set a.io_thread = 1
< error
| 1064 42000 …near 'a.io_thread = 1' at line 1
> exec ^
: string
| set a.ipc = 1
< error
| 1064 42000 …near 'a.ipc = 1' at line 1
> exec ^
: string
| set a.is = 1
< error
| 1064 42000 …near 'a.is = 1' at line 1
> exec ^
: string
| set a.isolation = 1
< error
| 1064 42000 …near 'a.isolation = 1' at line 1
> exec ^
: string
| set a.issuer = 1
< error
| 1064 42000 …near 'a.issuer = 1' at line 1
> exec ^
: string
| set a.iterate = 1
< error
| 1064 42000 …near 'a.iterate = 1' at line 1
> exec ^
: string
| set a.join = 1
< error
| 1064 42000 …near 'a.join = 1' at line 1
> exec ^
: string
| set a.json = 1
< error
| 1064 42000 …near 'a.json = 1' at line 1
> exec ^
: string
| set a.json_arrayagg = 1
< error
| 1064 42000 …near 'a.json_arrayagg = 1' at line 1
> exec ^
: string
| set a.json_duality_object = 1
< error
| 1064 42000 …near 'a.json_duality_object = 1' at line 1
> exec ^
: string
| set a.json_objectagg = 1
< error
| 1064 42000 …near 'a.json_objectagg = 1' at line 1
> exec ^
: string
| set a.json_table = 1
< error
| 1064 42000 …near 'a.json_table = 1' at line 1
> exec ^
: string
| set a.json_value = 1
< error
| 1064 42000 …near 'a.json_value = 1' at line 1
> exec ^
: string
| set a.key = 1
< error
| 1064 42000 …near 'a.key = 1' at line 1
> exec ^
: string
| set a.key_block_size = 1
< error
| 1064 42000 …near 'a.key_block_size = 1' at line 1
> exec ^
: string
| set a.keyring = 1
< error
| 1064 42000 …near 'a.keyring = 1' at line 1
> exec ^
: string
| set a.keys = 1
< error
| 1064 42000 …near 'a.keys = 1' at line 1
> exec ^
: string
| set a.kill = 1
< error
| 1064 42000 …near 'a.kill = 1' at line 1
> exec ^
: string
| set a.lag = 1
< error
| 1064 42000 …near 'a.lag = 1' at line 1
> exec ^
: string
| set a.language = 1
< error
| 1064 42000 …near 'a.language = 1' at line 1
> exec ^
: string
| set a.last = 1
< error
| 1064 42000 …near 'a.last = 1' at line 1
> exec ^
: string
| set a.last_value = 1
< error
| 1064 42000 …near 'a.last_value = 1' at line 1
> exec ^
: string
| set a.lateral = 1
< error
| 1064 42000 …near 'a.lateral = 1' at line 1
> exec ^
: string
| set a.lead = 1
< error
| 1064 42000 …near 'a.lead = 1' at line 1
> exec ^
: string
| set a.leading = 1
< error
| 1064 42000 …near 'a.leading = 1' at line 1
> exec ^
: string
| set a.least = 1
< error
| 1064 42000 …near 'a.least = 1' at line 1
> exec ^
: string
| set a.leave = 1
< error
| 1064 42000 …near 'a.leave = 1' at line 1
> exec ^
: string
| set a.leaves = 1
< error
| 1064 42000 …near 'a.leaves = 1' at line 1
> exec ^
: string
| set a.left = 1
< error
| 1064 42000 …near 'a.left = 1' at line 1
> exec ^
: string
| set a.less = 1
< error
| 1064 42000 …near 'a.less = 1' at line 1
> exec ^
: string
| set a.level = 1
< error
| 1064 42000 …near 'a.level = 1' at line 1
> exec ^
: string
| set a.library = 1
< error
| 1064 42000 …near 'a.library = 1' at line 1
> exec ^
: string
| set a.like = 1
< error
| 1064 42000 …near 'a.like = 1' at line 1
> exec ^
: string
| set a.limit = 1
< error
| 1064 42000 …near 'a.limit = 1' at line 1
> exec ^
: string
| set a.linear = 1
< error
| 1064 42000 …near 'a.linear = 1' at line 1
> exec ^
: string
| set a.lines = 1
< error
| 1064 42000 …near 'a.lines = 1' at line 1
> exec ^
: string
| set a.linestring = 1
< error
| 1064 42000 …near 'a.linestring = 1' at line 1
> exec ^
: string
| set a.list = 1
< error
| 1064 42000 …near 'a.list = 1' at line 1
> exec ^
: string
| set a.load = 1
< error
| 1064 42000 …near 'a.load = 1' at line 1
> exec ^
: string
| set a.local = 1
< error
| 1064 42000 …near 'a.local = 1' at line 1
> exec ^
: string
| set a.localtime = 1
< error
| 1064 42000 …near 'a.localtime = 1' at line 1
> exec ^
: string
| set a.localtimestamp = 1
< error
| 1064 42000 …near 'a.localtimestamp = 1' at line 1
> exec ^
: string
| set a.lock = 1
< error
| 1064 42000 …near 'a.lock = 1' at line 1
> exec ^
: string
| set a.locked = 1
< error
| 1064 42000 …near 'a.locked = 1' at line 1
> exec ^
: string
| set a.locks = 1
< error
| 1064 42000 …near 'a.locks = 1' at line 1
> exec ^
: string
| set a.log = 1
< error
| 1064 42000 …near 'a.log = 1' at line 1
> exec ^
: string
| set a.logfile = 1
< error
| 1064 42000 …near 'a.logfile = 1' at line 1
> exec ^
: string
| set a.logs = 1
< error
| 1064 42000 …near 'a.logs = 1' at line 1
> exec ^
: string
| set a.long = 1
< error
| 1064 42000 …near 'a.long = 1' at line 1
> exec ^
: string
| set a.longblob = 1
< error
| 1064 42000 …near 'a.longblob = 1' at line 1
> exec ^
: string
| set a.longtext = 1
< error
| 1064 42000 …near 'a.longtext = 1' at line 1
> exec ^
: string
| set a.loop = 1
< error
| 1064 42000 …near 'a.loop = 1' at line 1
> exec ^
: string
| set a.low_priority = 1
< error
| 1064 42000 …near 'a.low_priority = 1' at line 1
> exec ^
: string
| set a.manifest = 1
< error
| 1064 42000 …near 'a.manifest = 1' at line 1
> exec ^
: string
| set a.manual = 1
< error
| 1064 42000 …near 'a.manual = 1' at line 1
> exec ^
: string
| set a.master = 1
< error
| 1064 42000 …near 'a.master = 1' at line 1
> exec ^
: string
| set a.master_auto_position = 1
< error
| 1064 42000 …near 'a.master_auto_position = 1' at line 1
> exec ^
: string
| set a.master_bind = 1
< error
| 1064 42000 …near 'a.master_bind = 1' at line 1
> exec ^
: string
| set a.master_compression_algorithms = 1
< error
| 1064 42000 …near 'a.master_compression_algorithms = 1' at line 1
> exec ^
: string
| set a.master_connect_retry = 1
< error
| 1064 42000 …near 'a.master_connect_retry = 1' at line 1
> exec ^
: string
| set a.master_delay = 1
< error
| 1064 42000 …near 'a.master_delay = 1' at line 1
> exec ^
: string
| set a.master_heartbeat_period = 1
< error
| 1064 42000 …near 'a.master_heartbeat_period = 1' at line 1
> exec ^
: string
| set a.master_host = 1
< error
| 1064 42000 …near 'a.master_host = 1' at line 1
> exec ^
: string
| set a.master_log_file = 1
< error
| 1064 42000 …near 'a.master_log_file = 1' at line 1
> exec ^
: string
| set a.master_log_pos = 1
< error
| 1064 42000 …near 'a.master_log_pos = 1' at line 1
> exec ^
: string
| set a.master_password = 1
< error
| 1064 42000 …near 'a.master_password = 1' at line 1
> exec ^
: string
| set a.master_port = 1
< error
| 1064 42000 …near 'a.master_port = 1' at line 1
> exec ^
: string
| set a.master_public_key_path = 1
< error
| 1064 42000 …near 'a.master_public_key_path = 1' at line 1
> exec ^
: string
| set a.master_retry_count = 1
< error
| 1064 42000 …near 'a.master_retry_count = 1' at line 1
> exec ^
: string
| set a.master_ssl = 1
< error
| 1064 42000 …near 'a.master_ssl = 1' at line 1
> exec ^
: string
| set a.master_ssl_ca = 1
< error
| 1064 42000 …near 'a.master_ssl_ca = 1' at line 1
> exec ^
: string
| set a.master_ssl_capath = 1
< error
| 1064 42000 …near 'a.master_ssl_capath = 1' at line 1
> exec ^
: string
| set a.master_ssl_cert = 1
< error
| 1064 42000 …near 'a.master_ssl_cert = 1' at line 1
> exec ^
: string
| set a.master_ssl_cipher = 1
< error
| 1064 42000 …near 'a.master_ssl_cipher = 1' at line 1
> exec ^
: string
| set a.master_ssl_crl = 1
< error
| 1064 42000 …near 'a.master_ssl_crl = 1' at line 1
> exec ^
: string
| set a.master_ssl_crlpath = 1
< error
| 1064 42000 …near 'a.master_ssl_crlpath = 1' at line 1
> exec ^
: string
| set a.master_ssl_key = 1
< error
| 1064 42000 …near 'a.master_ssl_key = 1' at line 1
> exec ^
: string
| set a.master_ssl_verify_server_cert = 1
< error
| 1064 42000 …near 'a.master_ssl_verify_server_cert = 1' at line 1
> exec ^
: string
| set a.master_tls_ciphersuites = 1
< error
| 1064 42000 …near 'a.master_tls_ciphersuites = 1' at line 1
> exec ^
: string
| set a.master_tls_version = 1
< error
| 1064 42000 …near 'a.master_tls_version = 1' at line 1
> exec ^
: string
| set a.master_user = 1
< error
| 1064 42000 …near 'a.master_user = 1' at line 1
> exec ^
: string
| set a.master_zstd_compression_level = 1
< error
| 1064 42000 …near 'a.master_zstd_compression_level = 1' at line 1
> exec ^
: string
| set a.match = 1
< error
| 1064 42000 …near 'a.match = 1' at line 1
> exec ^
: string
| set a.materialized = 1
< error
| 1064 42000 …near 'a.materialized = 1' at line 1
> exec ^
: string
| set a.max = 1
< error
| 1064 42000 …near 'a.max = 1' at line 1
> exec ^
: string
| set a.max_connections_per_hour = 1
< error
| 1064 42000 …near 'a.max_connections_per_hour = 1' at line 1
> exec ^
: string
| set a.max_queries_per_hour = 1
< error
| 1064 42000 …near 'a.max_queries_per_hour = 1' at line 1
> exec ^
: string
| set a.max_rows = 1
< error
| 1064 42000 …near 'a.max_rows = 1' at line 1
> exec ^
: string
| set a.max_size = 1
< error
| 1064 42000 …near 'a.max_size = 1' at line 1
> exec ^
: string
| set a.max_updates_per_hour = 1
< error
| 1064 42000 …near 'a.max_updates_per_hour = 1' at line 1
> exec ^
: string
| set a.max_user_connections = 1
< error
| 1064 42000 …near 'a.max_user_connections = 1' at line 1
> exec ^
: string
| set a.maxvalue = 1
< error
| 1064 42000 …near 'a.maxvalue = 1' at line 1
> exec ^
: string
| set a.medium = 1
< error
| 1064 42000 …near 'a.medium = 1' at line 1
> exec ^
: string
| set a.mediumblob = 1
< error
| 1064 42000 …near 'a.mediumblob = 1' at line 1
> exec ^
: string
| set a.mediumint = 1
< error
| 1064 42000 …near 'a.mediumint = 1' at line 1
> exec ^
: string
| set a.mediumtext = 1
< error
| 1064 42000 …near 'a.mediumtext = 1' at line 1
> exec ^
: string
| set a.member = 1
< error
| 1064 42000 …near 'a.member = 1' at line 1
> exec ^
: string
| set a.memory = 1
< error
| 1064 42000 …near 'a.memory = 1' at line 1
> exec ^
: string
| set a.merge = 1
< error
| 1064 42000 …near 'a.merge = 1' at line 1
> exec ^
: string
| set a.message_text = 1
< error
| 1064 42000 …near 'a.message_text = 1' at line 1
> exec ^
: string
| set a.microsecond = 1
< error
| 1064 42000 …near 'a.microsecond = 1' at line 1
> exec ^
: string
| set a.mid = 1
< error
| 1064 42000 …near 'a.mid = 1' at line 1
> exec ^
: string
| set a.middleint = 1
< error
| 1064 42000 …near 'a.middleint = 1' at line 1
> exec ^
: string
| set a.migrate = 1
< error
| 1064 42000 …near 'a.migrate = 1' at line 1
> exec ^
: string
| set a.min = 1
< error
| 1064 42000 …near 'a.min = 1' at line 1
> exec ^
: string
| set a.min_rows = 1
< error
| 1064 42000 …near 'a.min_rows = 1' at line 1
> exec ^
: string
| set a.minute = 1
< error
| 1064 42000 …near 'a.minute = 1' at line 1
> exec ^
: string
| set a.minute_microsecond = 1
< error
| 1064 42000 …near 'a.minute_microsecond = 1' at line 1
> exec ^
: string
| set a.minute_second = 1
< error
| 1064 42000 …near 'a.minute_second = 1' at line 1
> exec ^
: string
| set a.mod = 1
< error
| 1064 42000 …near 'a.mod = 1' at line 1
> exec ^
: string
| set a.mode = 1
< error
| 1064 42000 …near 'a.mode = 1' at line 1
> exec ^
: string
| set a.modifies = 1
< error
| 1064 42000 …near 'a.modifies = 1' at line 1
> exec ^
: string
| set a.modify = 1
< error
| 1064 42000 …near 'a.modify = 1' at line 1
> exec ^
: string
| set a.month = 1
< error
| 1064 42000 …near 'a.month = 1' at line 1
> exec ^
: string
| set a.multilinestring = 1
< error
| 1064 42000 …near 'a.multilinestring = 1' at line 1
> exec ^
: string
| set a.multipoint = 1
< error
| 1064 42000 …near 'a.multipoint = 1' at line 1
> exec ^
: string
| set a.multipolygon = 1
< error
| 1064 42000 …near 'a.multipolygon = 1' at line 1
> exec ^
: string
| set a.mutex = 1
< error
| 1064 42000 …near 'a.mutex = 1' at line 1
> exec ^
: string
| set a.mysql_errno = 1
< error
| 1064 42000 …near 'a.mysql_errno = 1' at line 1
> exec ^
: string
| set a.name = 1
< error
| 1064 42000 …near 'a.name = 1' at line 1
> exec ^
: string
| set a.names = 1
< error
| 1064 42000 …near 'a.names = 1' at line 1
> exec ^
: string
| set a.national = 1
< error
| 1064 42000 …near 'a.national = 1' at line 1
> exec ^
: string
| set a.natural = 1
< error
| 1064 42000 …near 'a.natural = 1' at line 1
> exec ^
: string
| set a.nchar = 1
< error
| 1064 42000 …near 'a.nchar = 1' at line 1
> exec ^
: string
| set a.ndb = 1
< error
| 1064 42000 …near 'a.ndb = 1' at line 1
> exec ^
: string
| set a.ndbcluster = 1
< error
| 1064 42000 …near 'a.ndbcluster = 1' at line 1
> exec ^
: string
| set a.nested = 1
< error
| 1064 42000 …near 'a.nested = 1' at line 1
> exec ^
: string
| set a.network_namespace = 1
< error
| 1064 42000 …near 'a.network_namespace = 1' at line 1
> exec ^
: string
| set a.never = 1
< error
| 1064 42000 …near 'a.never = 1' at line 1
> exec ^
: string
| set a.new = 1
< error
| 1064 42000 …near 'a.new = 1' at line 1
> exec ^
: string
| set a.next = 1
< error
| 1064 42000 …near 'a.next = 1' at line 1
> exec ^
: string
| set a.no = 1
< error
| 1064 42000 …near 'a.no = 1' at line 1
> exec ^
: string
| set a.no_wait = 1
< error
| 1064 42000 …near 'a.no_wait = 1' at line 1
> exec ^
: string
| set a.no_write_to_binlog = 1
< error
| 1064 42000 …near 'a.no_write_to_binlog = 1' at line 1
> exec ^
: string
| set a.nodegroup = 1
< error
| 1064 42000 …near 'a.nodegroup = 1' at line 1
> exec ^
: string
| set a.none = 1
< error
| 1064 42000 …near 'a.none = 1' at line 1
> exec ^
: string
| set a.not = 1
< error
| 1064 42000 …near 'a.not = 1' at line 1
> exec ^
: string
| set a.now = 1
< error
| 1064 42000 …near 'a.now = 1' at line 1
> exec ^
: string
| set a.nowait = 1
< error
| 1064 42000 …near 'a.nowait = 1' at line 1
> exec ^
: string
| set a.nth_value = 1
< error
| 1064 42000 …near 'a.nth_value = 1' at line 1
> exec ^
: string
| set a.ntile = 1
< error
| 1064 42000 …near 'a.ntile = 1' at line 1
> exec ^
: string
| set a.null = 1
< error
| 1064 42000 …near 'a.null = 1' at line 1
> exec ^
: string
| set a.nullif = 1
< error
| 1064 42000 …near 'a.nullif = 1' at line 1
> exec ^
: string
| set a.nulls = 1
< error
| 1064 42000 …near 'a.nulls = 1' at line 1
> exec ^
: string
| set a.number = 1
< error
| 1064 42000 …near 'a.number = 1' at line 1
> exec ^
: string
| set a.numeric = 1
< error
| 1064 42000 …near 'a.numeric = 1' at line 1
> exec ^
: string
| set a.nvarchar = 1
< error
| 1064 42000 …near 'a.nvarchar = 1' at line 1
> exec ^
: string
| set a.of = 1
< error
| 1064 42000 …near 'a.of = 1' at line 1
> exec ^
: string
| set a.off = 1
< error
| 1064 42000 …near 'a.off = 1' at line 1
> exec ^
: string
| set a.offset = 1
< error
| 1064 42000 …near 'a.offset = 1' at line 1
> exec ^
: string
| set a.oj = 1
< error
| 1064 42000 …near 'a.oj = 1' at line 1
> exec ^
: string
| set a.old = 1
< error
| 1064 42000 …near 'a.old = 1' at line 1
> exec ^
: string
| set a.on = 1
< error
| 1064 42000 …near 'a.on = 1' at line 1
> exec ^
: string
| set a.one = 1
< error
| 1064 42000 …near 'a.one = 1' at line 1
> exec ^
: string
| set a.only = 1
< error
| 1064 42000 …near 'a.only = 1' at line 1
> exec ^
: string
| set a.open = 1
< error
| 1064 42000 …near 'a.open = 1' at line 1
> exec ^
: string
| set a.optimize = 1
< error
| 1064 42000 …near 'a.optimize = 1' at line 1
> exec ^
: string
| set a.optimizer_costs = 1
< error
| 1064 42000 …near 'a.optimizer_costs = 1' at line 1
> exec ^
: string
| set a.option = 1
< error
| 1064 42000 …near 'a.option = 1' at line 1
> exec ^
: string
| set a.optional = 1
< error
| 1064 42000 …near 'a.optional = 1' at line 1
> exec ^
: string
| set a.optionally = 1
< error
| 1064 42000 …near 'a.optionally = 1' at line 1
> exec ^
: string
| set a.options = 1
< error
| 1064 42000 …near 'a.options = 1' at line 1
> exec ^
: string
| set a.or = 1
< error
| 1064 42000 …near 'a.or = 1' at line 1
> exec ^
: string
| set a.order = 1
< error
| 1064 42000 …near 'a.order = 1' at line 1
> exec ^
: string
| set a.ordinality = 1
< error
| 1064 42000 …near 'a.ordinality = 1' at line 1
> exec ^
: string
| set a.organization = 1
< error
| 1064 42000 …near 'a.organization = 1' at line 1
> exec ^
: string
| set a.others = 1
< error
| 1064 42000 …near 'a.others = 1' at line 1
> exec ^
: string
| set a.out = 1
< error
| 1064 42000 …near 'a.out = 1' at line 1
> exec ^
: string
| set a.outer = 1
< error
| 1064 42000 …near 'a.outer = 1' at line 1
> exec ^
: string
| set a.outfile = 1
< error
| 1064 42000 …near 'a.outfile = 1' at line 1
> exec ^
: string
| set a.over = 1
< error
| 1064 42000 …near 'a.over = 1' at line 1
> exec ^
: string
| set a.overlay = 1
< error
| 1064 42000 …near 'a.overlay = 1' at line 1
> exec ^
: string
| set a.overwrite = 1
< error
| 1064 42000 …near 'a.overwrite = 1' at line 1
> exec ^
: string
| set a.owner = 1
< error
| 1064 42000 …near 'a.owner = 1' at line 1
> exec ^
: string
| set a.pack_keys = 1
< error
| 1064 42000 …near 'a.pack_keys = 1' at line 1
> exec ^
: string
| set a.page = 1
< error
| 1064 42000 …near 'a.page = 1' at line 1
> exec ^
: string
| set a.parallel = 1
< error
| 1064 42000 …near 'a.parallel = 1' at line 1
> exec ^
: string
| set a.parameters = 1
< error
| 1064 42000 …near 'a.parameters = 1' at line 1
> exec ^
: string
| set a.parse_tree = 1
< error
| 1064 42000 …near 'a.parse_tree = 1' at line 1
> exec ^
: string
| set a.parser = 1
< error
| 1064 42000 …near 'a.parser = 1' at line 1
> exec ^
: string
| set a.partial = 1
< error
| 1064 42000 …near 'a.partial = 1' at line 1
> exec ^
: string
| set a.partition = 1
< error
| 1064 42000 …near 'a.partition = 1' at line 1
> exec ^
: string
| set a.partitioning = 1
< error
| 1064 42000 …near 'a.partitioning = 1' at line 1
> exec ^
: string
| set a.partitions = 1
< error
| 1064 42000 …near 'a.partitions = 1' at line 1
> exec ^
: string
| set a.password = 1
< error
| 1064 42000 …near 'a.password = 1' at line 1
> exec ^
: string
| set a.password_lock_time = 1
< error
| 1064 42000 …near 'a.password_lock_time = 1' at line 1
> exec ^
: string
| set a.path = 1
< error
| 1064 42000 …near 'a.path = 1' at line 1
> exec ^
: string
| set a.pctfree = 1
< error
| 1064 42000 …near 'a.pctfree = 1' at line 1
> exec ^
: string
| set a.percent_rank = 1
< error
| 1064 42000 …near 'a.percent_rank = 1' at line 1
> exec ^
: string
| set a.persist = 1
< error
| 1064 42000 …near 'a.persist = 1' at line 1
> exec ^
: string
| set a.persist_only = 1
< error
| 1064 42000 …near 'a.persist_only = 1' at line 1
> exec ^
: string
| set a.phase = 1
< error
| 1064 42000 …near 'a.phase = 1' at line 1
> exec ^
: string
| set a.placing = 1
< error
| 1064 42000 …near 'a.placing = 1' at line 1
> exec ^
: string
| set a.plugin = 1
< error
| 1064 42000 …near 'a.plugin = 1' at line 1
> exec ^
: string
| set a.plugin_dir = 1
< error
| 1064 42000 …near 'a.plugin_dir = 1' at line 1
> exec ^
: string
| set a.plugins = 1
< error
| 1064 42000 …near 'a.plugins = 1' at line 1
> exec ^
: string
| set a.point = 1
< error
| 1064 42000 …near 'a.point = 1' at line 1
> exec ^
: string
| set a.polygon = 1
< error
| 1064 42000 …near 'a.polygon = 1' at line 1
> exec ^
: string
| set a.port = 1
< error
| 1064 42000 …near 'a.port = 1' at line 1
> exec ^
: string
| set a.position = 1
< error
| 1064 42000 …near 'a.position = 1' at line 1
> exec ^
: string
| set a.precedes = 1
< error
| 1064 42000 …near 'a.precedes = 1' at line 1
> exec ^
: string
| set a.preceding = 1
< error
| 1064 42000 …near 'a.preceding = 1' at line 1
> exec ^
: string
| set a.precision = 1
< error
| 1064 42000 …near 'a.precision = 1' at line 1
> exec ^
: string
| set a.prefix = 1
< error
| 1064 42000 …near 'a.prefix = 1' at line 1
> exec ^
: string
| set a.prepare = 1
< error
| 1064 42000 …near 'a.prepare = 1' at line 1
> exec ^
: string
| set a.preserve = 1
< error
| 1064 42000 …near 'a.preserve = 1' at line 1
> exec ^
: string
| set a.prev = 1
< error
| 1064 42000 …near 'a.prev = 1' at line 1
> exec ^
: string
| set a.primary = 1
< error
| 1064 42000 …near 'a.primary = 1' at line 1
> exec ^
: string
| set a.privilege_checks_user = 1
< error
| 1064 42000 …near 'a.privilege_checks_user = 1' at line 1
> exec ^
: string
| set a.privileges = 1
< error
| 1064 42000 …near 'a.privileges = 1' at line 1
> exec ^
: string
| set a.procedure = 1
< error
| 1064 42000 …near 'a.procedure = 1' at line 1
> exec ^
: string
| set a.process = 1
< error
| 1064 42000 …near 'a.process = 1' at line 1
> exec ^
: string
| set a.processlist = 1
< error
| 1064 42000 …near 'a.processlist = 1' at line 1
> exec ^
: string
| set a.profile = 1
< error
| 1064 42000 …near 'a.profile = 1' at line 1
> exec ^
: string
| set a.profiles = 1
< error
| 1064 42000 …near 'a.profiles = 1' at line 1
> exec ^
: string
| set a.proxy = 1
< error
| 1064 42000 …near 'a.proxy = 1' at line 1
> exec ^
: string
| set a.purge = 1
< error
| 1064 42000 …near 'a.purge = 1' at line 1
> exec ^
: string
| set a.qualify = 1
< error
| 1064 42000 …near 'a.qualify = 1' at line 1
> exec ^
: string
| set a.quarter = 1
< error
| 1064 42000 …near 'a.quarter = 1' at line 1
> exec ^
: string
| set a.query = 1
< error
| 1064 42000 …near 'a.query = 1' at line 1
> exec ^
: string
| set a.quick = 1
< error
| 1064 42000 …near 'a.quick = 1' at line 1
> exec ^
: string
| set a.random = 1
< error
| 1064 42000 …near 'a.random = 1' at line 1
> exec ^
: string
| set a.range = 1
< error
| 1064 42000 …near 'a.range = 1' at line 1
> exec ^
: string
| set a.rank = 1
< error
| 1064 42000 …near 'a.rank = 1' at line 1
> exec ^
: string
| set a.read = 1
< error
| 1064 42000 …near 'a.read = 1' at line 1
> exec ^
: string
| set a.read_only = 1
< error
| 1064 42000 …near 'a.read_only = 1' at line 1
> exec ^
: string
| set a.read_write = 1
< error
| 1064 42000 …near 'a.read_write = 1' at line 1
> exec ^
: string
| set a.reads = 1
< error
| 1064 42000 …near 'a.reads = 1' at line 1
> exec ^
: string
| set a.real = 1
< error
| 1064 42000 …near 'a.real = 1' at line 1
> exec ^
: string
| set a.rebuild = 1
< error
| 1064 42000 …near 'a.rebuild = 1' at line 1
> exec ^
: string
| set a.recover = 1
< error
| 1064 42000 …near 'a.recover = 1' at line 1
> exec ^
: string
| set a.recursive = 1
< error
| 1064 42000 …near 'a.recursive = 1' at line 1
> exec ^
: string
| set a.redo_buffer_size = 1
< error
| 1064 42000 …near 'a.redo_buffer_size = 1' at line 1
> exec ^
: string
| set a.redundant = 1
< error
| 1064 42000 …near 'a.redundant = 1' at line 1
> exec ^
: string
| set a.reference = 1
< error
| 1064 42000 …near 'a.reference = 1' at line 1
> exec ^
: string
| set a.references = 1
< error
| 1064 42000 …near 'a.references = 1' at line 1
> exec ^
: string
| set a.regexp = 1
< error
| 1064 42000 …near 'a.regexp = 1' at line 1
> exec ^
: string
| set a.registration = 1
< error
| 1064 42000 …near 'a.registration = 1' at line 1
> exec ^
: string
| set a.relational = 1
< error
| 1064 42000 …near 'a.relational = 1' at line 1
> exec ^
: string
| set a.relay = 1
< error
| 1064 42000 …near 'a.relay = 1' at line 1
> exec ^
: string
| set a.relay_log_file = 1
< error
| 1064 42000 …near 'a.relay_log_file = 1' at line 1
> exec ^
: string
| set a.relay_log_pos = 1
< error
| 1064 42000 …near 'a.relay_log_pos = 1' at line 1
> exec ^
: string
| set a.relay_thread = 1
< error
| 1064 42000 …near 'a.relay_thread = 1' at line 1
> exec ^
: string
| set a.relaylog = 1
< error
| 1064 42000 …near 'a.relaylog = 1' at line 1
> exec ^
: string
| set a.release = 1
< error
| 1064 42000 …near 'a.release = 1' at line 1
> exec ^
: string
| set a.reload = 1
< error
| 1064 42000 …near 'a.reload = 1' at line 1
> exec ^
: string
| set a.remove = 1
< error
| 1064 42000 …near 'a.remove = 1' at line 1
> exec ^
: string
| set a.rename = 1
< error
| 1064 42000 …near 'a.rename = 1' at line 1
> exec ^
: string
| set a.reorganize = 1
< error
| 1064 42000 …near 'a.reorganize = 1' at line 1
> exec ^
: string
| set a.repair = 1
< error
| 1064 42000 …near 'a.repair = 1' at line 1
> exec ^
: string
| set a.repeat = 1
< error
| 1064 42000 …near 'a.repeat = 1' at line 1
> exec ^
: string
| set a.repeatable = 1
< error
| 1064 42000 …near 'a.repeatable = 1' at line 1
> exec ^
: string
| set a.replace = 1
< error
| 1064 42000 …near 'a.replace = 1' at line 1
> exec ^
: string
| set a.replica = 1
< error
| 1064 42000 …near 'a.replica = 1' at line 1
> exec ^
: string
| set a.replica_num = 1
< error
| 1064 42000 …near 'a.replica_num = 1' at line 1
> exec ^
: string
| set a.replicas = 1
< error
| 1064 42000 …near 'a.replicas = 1' at line 1
> exec ^
: string
| set a.replicate_do_db = 1
< error
| 1064 42000 …near 'a.replicate_do_db = 1' at line 1
> exec ^
: string
| set a.replicate_do_table = 1
< error
| 1064 42000 …near 'a.replicate_do_table = 1' at line 1
> exec ^
: string
| set a.replicate_ignore_db = 1
< error
| 1064 42000 …near 'a.replicate_ignore_db = 1' at line 1
> exec ^
: string
| set a.replicate_ignore_table = 1
< error
| 1064 42000 …near 'a.replicate_ignore_table = 1' at line 1
> exec ^
: string
| set a.replicate_rewrite_db = 1
< error
| 1064 42000 …near 'a.replicate_rewrite_db = 1' at line 1
> exec ^
: string
| set a.replicate_wild_do_table = 1
< error
| 1064 42000 …near 'a.replicate_wild_do_table = 1' at line 1
> exec ^
: string
| set a.replicate_wild_ignore_table = 1
< error
| 1064 42000 …near 'a.replicate_wild_ignore_table = 1' at line 1
> exec ^
: string
| set a.replication = 1
< error
| 1064 42000 …near 'a.replication = 1' at line 1
> exec ^
: string
| set a.require = 1
< error
| 1064 42000 …near 'a.require = 1' at line 1
> exec ^
: string
| set a.require_row_format = 1
< error
| 1064 42000 …near 'a.require_row_format = 1' at line 1
> exec ^
: string
| set a.require_table_primary_key_check = 1
< error
| 1064 42000 …near 'a.require_table_primary_key_check = 1' at line 1
> exec ^
: string
| set a.reset = 1
< error
| 1064 42000 …near 'a.reset = 1' at line 1
> exec ^
: string
| set a.resignal = 1
< error
| 1064 42000 …near 'a.resignal = 1' at line 1
> exec ^
: string
| set a.resource = 1
< error
| 1064 42000 …near 'a.resource = 1' at line 1
> exec ^
: string
| set a.respect = 1
< error
| 1064 42000 …near 'a.respect = 1' at line 1
> exec ^
: string
| set a.restart = 1
< error
| 1064 42000 …near 'a.restart = 1' at line 1
> exec ^
: string
| set a.restore = 1
< error
| 1064 42000 …near 'a.restore = 1' at line 1
> exec ^
: string
| set a.restrict = 1
< error
| 1064 42000 …near 'a.restrict = 1' at line 1
> exec ^
: string
| set a.resume = 1
< error
| 1064 42000 …near 'a.resume = 1' at line 1
> exec ^
: string
| set a.retain = 1
< error
| 1064 42000 …near 'a.retain = 1' at line 1
> exec ^
: string
| set a.return = 1
< error
| 1064 42000 …near 'a.return = 1' at line 1
> exec ^
: string
| set a.returned_sqlstate = 1
< error
| 1064 42000 …near 'a.returned_sqlstate = 1' at line 1
> exec ^
: string
| set a.returning = 1
< error
| 1064 42000 …near 'a.returning = 1' at line 1
> exec ^
: string
| set a.returns = 1
< error
| 1064 42000 …near 'a.returns = 1' at line 1
> exec ^
: string
| set a.reuse = 1
< error
| 1064 42000 …near 'a.reuse = 1' at line 1
> exec ^
: string
| set a.reverse = 1
< error
| 1064 42000 …near 'a.reverse = 1' at line 1
> exec ^
: string
| set a.revoke = 1
< error
| 1064 42000 …near 'a.revoke = 1' at line 1
> exec ^
: string
| set a.right = 1
< error
| 1064 42000 …near 'a.right = 1' at line 1
> exec ^
: string
| set a.rlike = 1
< error
| 1064 42000 …near 'a.rlike = 1' at line 1
> exec ^
: string
| set a.role = 1
< error
| 1064 42000 …near 'a.role = 1' at line 1
> exec ^
: string
| set a.rollback = 1
< error
| 1064 42000 …near 'a.rollback = 1' at line 1
> exec ^
: string
| set a.rollup = 1
< error
| 1064 42000 …near 'a.rollup = 1' at line 1
> exec ^
: string
| set a.rotate = 1
< error
| 1064 42000 …near 'a.rotate = 1' at line 1
> exec ^
: string
| set a.routine = 1
< error
| 1064 42000 …near 'a.routine = 1' at line 1
> exec ^
: string
| set a.row = 1
< error
| 1064 42000 …near 'a.row = 1' at line 1
> exec ^
: string
| set a.row_count = 1
< error
| 1064 42000 …near 'a.row_count = 1' at line 1
> exec ^
: string
| set a.row_format = 1
< error
| 1064 42000 …near 'a.row_format = 1' at line 1
> exec ^
: string
| set a.row_number = 1
< error
| 1064 42000 …near 'a.row_number = 1' at line 1
> exec ^
: string
| set a.rows = 1
< error
| 1064 42000 …near 'a.rows = 1' at line 1
> exec ^
: string
| set a.rtree = 1
< error
| 1064 42000 …near 'a.rtree = 1' at line 1
> exec ^
: string
| set a.s3 = 1
< error
| 1064 42000 …near 'a.s3 = 1' at line 1
> exec ^
: string
| set a.savepoint = 1
< error
| 1064 42000 …near 'a.savepoint = 1' at line 1
> exec ^
: string
| set a.schedule = 1
< error
| 1064 42000 …near 'a.schedule = 1' at line 1
> exec ^
: string
| set a.schema = 1
< error
| 1064 42000 …near 'a.schema = 1' at line 1
> exec ^
: string
| set a.schema_name = 1
< error
| 1064 42000 …near 'a.schema_name = 1' at line 1
> exec ^
: string
| set a.schemas = 1
< error
| 1064 42000 …near 'a.schemas = 1' at line 1
> exec ^
: string
| set a.second = 1
< error
| 1064 42000 …near 'a.second = 1' at line 1
> exec ^
: string
| set a.second_microsecond = 1
< error
| 1064 42000 …near 'a.second_microsecond = 1' at line 1
> exec ^
: string
| set a.secondary = 1
< error
| 1064 42000 …near 'a.secondary = 1' at line 1
> exec ^
: string
| set a.secondary_engine = 1
< error
| 1064 42000 …near 'a.secondary_engine = 1' at line 1
> exec ^
: string
| set a.secondary_engine_attribute = 1
< error
| 1064 42000 …near 'a.secondary_engine_attribute = 1' at line 1
> exec ^
: string
| set a.secondary_load = 1
< error
| 1064 42000 …near 'a.secondary_load = 1' at line 1
> exec ^
: string
| set a.secondary_unload = 1
< error
| 1064 42000 …near 'a.secondary_unload = 1' at line 1
> exec ^
: string
| set a.security = 1
< error
| 1064 42000 …near 'a.security = 1' at line 1
> exec ^
: string
| set a.select = 1
< error
| 1064 42000 …near 'a.select = 1' at line 1
> exec ^
: string
| set a.sensitive = 1
< error
| 1064 42000 …near 'a.sensitive = 1' at line 1
> exec ^
: string
| set a.separator = 1
< error
| 1064 42000 …near 'a.separator = 1' at line 1
> exec ^
: string
| set a.serial = 1
< error
| 1064 42000 …near 'a.serial = 1' at line 1
> exec ^
: string
| set a.serializable = 1
< error
| 1064 42000 …near 'a.serializable = 1' at line 1
> exec ^
: string
| set a.server = 1
< error
| 1064 42000 …near 'a.server = 1' at line 1
> exec ^
: string
| set a.session = 1
< error
| 1064 42000 …near 'a.session = 1' at line 1
> exec ^
: string
| set a.session_user = 1
< error
| 1064 42000 …near 'a.session_user = 1' at line 1
> exec ^
: string
| set a.set = 1
< error
| 1064 42000 …near 'a.set = 1' at line 1
> exec ^
: string
| set a.sets = 1
< error
| 1064 42000 …near 'a.sets = 1' at line 1
> exec ^
: string
| set a.share = 1
< error
| 1064 42000 …near 'a.share = 1' at line 1
> exec ^
: string
| set a.show = 1
< error
| 1064 42000 …near 'a.show = 1' at line 1
> exec ^
: string
| set a.shutdown = 1
< error
| 1064 42000 …near 'a.shutdown = 1' at line 1
> exec ^
: string
| set a.signal = 1
< error
| 1064 42000 …near 'a.signal = 1' at line 1
> exec ^
: string
| set a.signed = 1
< error
| 1064 42000 …near 'a.signed = 1' at line 1
> exec ^
: string
| set a.simple = 1
< error
| 1064 42000 …near 'a.simple = 1' at line 1
> exec ^
: string
| set a.skip = 1
< error
| 1064 42000 …near 'a.skip = 1' at line 1
> exec ^
: string
| set a.slave = 1
< error
| 1064 42000 …near 'a.slave = 1' at line 1
> exec ^
: string
| set a.slow = 1
< error
| 1064 42000 …near 'a.slow = 1' at line 1
> exec ^
: string
| set a.smallint = 1
< error
| 1064 42000 …near 'a.smallint = 1' at line 1
> exec ^
: string
| set a.snapshot = 1
< error
| 1064 42000 …near 'a.snapshot = 1' at line 1
> exec ^
: string
| set a.socket = 1
< error
| 1064 42000 …near 'a.socket = 1' at line 1
> exec ^
: string
| set a.some = 1
< error
| 1064 42000 …near 'a.some = 1' at line 1
> exec ^
: string
| set a.soname = 1
< error
| 1064 42000 …near 'a.soname = 1' at line 1
> exec ^
: string
| set a.sounds = 1
< error
| 1064 42000 …near 'a.sounds = 1' at line 1
> exec ^
: string
| set a.source = 1
< error
| 1064 42000 …near 'a.source = 1' at line 1
> exec ^
: string
| set a.source_auto_position = 1
< error
| 1064 42000 …near 'a.source_auto_position = 1' at line 1
> exec ^
: string
| set a.source_bind = 1
< error
| 1064 42000 …near 'a.source_bind = 1' at line 1
> exec ^
: string
| set a.source_compression_algorithms = 1
< error
| 1064 42000 …near 'a.source_compression_algorithms = 1' at line 1
> exec ^
: string
| set a.source_connect_retry = 1
< error
| 1064 42000 …near 'a.source_connect_retry = 1' at line 1
> exec ^
: string
| set a.source_connection_auto_failover = 1
< error
| 1064 42000 …near 'a.source_connection_auto_failover = 1' at line 1
> exec ^
: string
| set a.source_delay = 1
< error
| 1064 42000 …near 'a.source_delay = 1' at line 1
> exec ^
: string
| set a.source_heartbeat_period = 1
< error
| 1064 42000 …near 'a.source_heartbeat_period = 1' at line 1
> exec ^
: string
| set a.source_host = 1
< error
| 1064 42000 …near 'a.source_host = 1' at line 1
> exec ^
: string
| set a.source_log_file = 1
< error
| 1064 42000 …near 'a.source_log_file = 1' at line 1
> exec ^
: string
| set a.source_log_pos = 1
< error
| 1064 42000 …near 'a.source_log_pos = 1' at line 1
> exec ^
: string
| set a.source_password = 1
< error
| 1064 42000 …near 'a.source_password = 1' at line 1
> exec ^
: string
| set a.source_port = 1
< error
| 1064 42000 …near 'a.source_port = 1' at line 1
> exec ^
: string
| set a.source_public_key_path = 1
< error
| 1064 42000 …near 'a.source_public_key_path = 1' at line 1
> exec ^
: string
| set a.source_retry_count = 1
< error
| 1064 42000 …near 'a.source_retry_count = 1' at line 1
> exec ^
: string
| set a.source_ssl = 1
< error
| 1064 42000 …near 'a.source_ssl = 1' at line 1
> exec ^
: string
| set a.source_ssl_ca = 1
< error
| 1064 42000 …near 'a.source_ssl_ca = 1' at line 1
> exec ^
: string
| set a.source_ssl_capath = 1
< error
| 1064 42000 …near 'a.source_ssl_capath = 1' at line 1
> exec ^
: string
| set a.source_ssl_cert = 1
< error
| 1064 42000 …near 'a.source_ssl_cert = 1' at line 1
> exec ^
: string
| set a.source_ssl_cipher = 1
< error
| 1064 42000 …near 'a.source_ssl_cipher = 1' at line 1
> exec ^
: string
| set a.source_ssl_crl = 1
< error
| 1064 42000 …near 'a.source_ssl_crl = 1' at line 1
> exec ^
: string
| set a.source_ssl_crlpath = 1
< error
| 1064 42000 …near 'a.source_ssl_crlpath = 1' at line 1
> exec ^
: string
| set a.source_ssl_key = 1
< error
| 1064 42000 …near 'a.source_ssl_key = 1' at line 1
> exec ^
: string
| set a.source_ssl_verify_server_cert = 1
< error
| 1064 42000 …near 'a.source_ssl_verify_server_cert = 1' at line 1
> exec ^
: string
| set a.source_tls_ciphersuites = 1
< error
| 1064 42000 …near 'a.source_tls_ciphersuites = 1' at line 1
> exec ^
: string
| set a.source_tls_version = 1
< error
| 1064 42000 …near 'a.source_tls_version = 1' at line 1
> exec ^
: string
| set a.source_user = 1
< error
| 1064 42000 …near 'a.source_user = 1' at line 1
> exec ^
: string
| set a.source_zstd_compression_level = 1
< error
| 1064 42000 …near 'a.source_zstd_compression_level = 1' at line 1
> exec ^
: string
| set a.spatial = 1
< error
| 1064 42000 …near 'a.spatial = 1' at line 1
> exec ^
: string
| set a.specific = 1
< error
| 1064 42000 …near 'a.specific = 1' at line 1
> exec ^
: string
| set a.sql = 1
< error
| 1064 42000 …near 'a.sql = 1' at line 1
> exec ^
: string
| set a.sql_after_gtids = 1
< error
| 1064 42000 …near 'a.sql_after_gtids = 1' at line 1
> exec ^
: string
| set a.sql_after_mts_gaps = 1
< error
| 1064 42000 …near 'a.sql_after_mts_gaps = 1' at line 1
> exec ^
: string
| set a.sql_before_gtids = 1
< error
| 1064 42000 …near 'a.sql_before_gtids = 1' at line 1
> exec ^
: string
| set a.sql_big_result = 1
< error
| 1064 42000 …near 'a.sql_big_result = 1' at line 1
> exec ^
: string
| set a.sql_buffer_result = 1
< error
| 1064 42000 …near 'a.sql_buffer_result = 1' at line 1
> exec ^
: string
| set a.sql_calc_found_rows = 1
< error
| 1064 42000 …near 'a.sql_calc_found_rows = 1' at line 1
> exec ^
: string
| set a.sql_no_cache = 1
< error
| 1064 42000 …near 'a.sql_no_cache = 1' at line 1
> exec ^
: string
| set a.sql_small_result = 1
< error
| 1064 42000 …near 'a.sql_small_result = 1' at line 1
> exec ^
: string
| set a.sql_thread = 1
< error
| 1064 42000 …near 'a.sql_thread = 1' at line 1
> exec ^
: string
| set a.sql_tsi_day = 1
< error
| 1064 42000 …near 'a.sql_tsi_day = 1' at line 1
> exec ^
: string
| set a.sql_tsi_hour = 1
< error
| 1064 42000 …near 'a.sql_tsi_hour = 1' at line 1
> exec ^
: string
| set a.sql_tsi_minute = 1
< error
| 1064 42000 …near 'a.sql_tsi_minute = 1' at line 1
> exec ^
: string
| set a.sql_tsi_month = 1
< error
| 1064 42000 …near 'a.sql_tsi_month = 1' at line 1
> exec ^
: string
| set a.sql_tsi_quarter = 1
< error
| 1064 42000 …near 'a.sql_tsi_quarter = 1' at line 1
> exec ^
: string
| set a.sql_tsi_second = 1
< error
| 1064 42000 …near 'a.sql_tsi_second = 1' at line 1
> exec ^
: string
| set a.sql_tsi_week = 1
< error
| 1064 42000 …near 'a.sql_tsi_week = 1' at line 1
> exec ^
: string
| set a.sql_tsi_year = 1
< error
| 1064 42000 …near 'a.sql_tsi_year = 1' at line 1
> exec ^
: string
| set a.sqlexception = 1
< error
| 1064 42000 …near 'a.sqlexception = 1' at line 1
> exec ^
: string
| set a.sqlstate = 1
< error
| 1064 42000 …near 'a.sqlstate = 1' at line 1
> exec ^
: string
| set a.sqlwarning = 1
< error
| 1064 42000 …near 'a.sqlwarning = 1' at line 1
> exec ^
: string
| set a.srid = 1
< error
| 1064 42000 …near 'a.srid = 1' at line 1
> exec ^
: string
| set a.sse_kms = 1
< error
| 1064 42000 …near 'a.sse_kms = 1' at line 1
> exec ^
: string
| set a.sse_s3 = 1
< error
| 1064 42000 …near 'a.sse_s3 = 1' at line 1
> exec ^
: string
| set a.ssl = 1
< error
| 1064 42000 …near 'a.ssl = 1' at line 1
> exec ^
: string
| set a.st_collect = 1
< error
| 1064 42000 …near 'a.st_collect = 1' at line 1
> exec ^
: string
| set a.stacked = 1
< error
| 1064 42000 …near 'a.stacked = 1' at line 1
> exec ^
: string
| set a.start = 1
< error
| 1064 42000 …near 'a.start = 1' at line 1
> exec ^
: string
| set a.starting = 1
< error
| 1064 42000 …near 'a.starting = 1' at line 1
> exec ^
: string
| set a.starts = 1
< error
| 1064 42000 …near 'a.starts = 1' at line 1
> exec ^
: string
| set a.stats_auto_recalc = 1
< error
| 1064 42000 …near 'a.stats_auto_recalc = 1' at line 1
> exec ^
: string
| set a.stats_persistent = 1
< error
| 1064 42000 …near 'a.stats_persistent = 1' at line 1
> exec ^
: string
| set a.stats_sample_pages = 1
< error
| 1064 42000 …near 'a.stats_sample_pages = 1' at line 1
> exec ^
: string
| set a.status = 1
< error
| 1064 42000 …near 'a.status = 1' at line 1
> exec ^
: string
| set a.std = 1
< error
| 1064 42000 …near 'a.std = 1' at line 1
> exec ^
: string
| set a.stddev = 1
< error
| 1064 42000 …near 'a.stddev = 1' at line 1
> exec ^
: string
| set a.stddev_pop = 1
< error
| 1064 42000 …near 'a.stddev_pop = 1' at line 1
> exec ^
: string
| set a.stddev_samp = 1
< error
| 1064 42000 …near 'a.stddev_samp = 1' at line 1
> exec ^
: string
| set a.stop = 1
< error
| 1064 42000 …near 'a.stop = 1' at line 1
> exec ^
: string
| set a.storage = 1
< error
| 1064 42000 …near 'a.storage = 1' at line 1
> exec ^
: string
| set a.stored = 1
< error
| 1064 42000 …near 'a.stored = 1' at line 1
> exec ^
: string
| set a.straight_join = 1
< error
| 1064 42000 …near 'a.straight_join = 1' at line 1
> exec ^
: string
| set a.stream = 1
< error
| 1064 42000 …near 'a.stream = 1' at line 1
> exec ^
: string
| set a.strict_load = 1
< error
| 1064 42000 …near 'a.strict_load = 1' at line 1
> exec ^
: string
| set a.string = 1
< error
| 1064 42000 …near 'a.string = 1' at line 1
> exec ^
: string
| set a.subclass_origin = 1
< error
| 1064 42000 …near 'a.subclass_origin = 1' at line 1
> exec ^
: string
| set a.subdate = 1
< error
| 1064 42000 …near 'a.subdate = 1' at line 1
> exec ^
: string
| set a.subject = 1
< error
| 1064 42000 …near 'a.subject = 1' at line 1
> exec ^
: string
| set a.subpartition = 1
< error
| 1064 42000 …near 'a.subpartition = 1' at line 1
> exec ^
: string
| set a.subpartitions = 1
< error
| 1064 42000 …near 'a.subpartitions = 1' at line 1
> exec ^
: string
| set a.substr = 1
< error
| 1064 42000 …near 'a.substr = 1' at line 1
> exec ^
: string
| set a.substring = 1
< error
| 1064 42000 …near 'a.substring = 1' at line 1
> exec ^
: string
| set a.sum = 1
< error
| 1064 42000 …near 'a.sum = 1' at line 1
> exec ^
: string
| set a.super = 1
< error
| 1064 42000 …near 'a.super = 1' at line 1
> exec ^
: string
| set a.suspend = 1
< error
| 1064 42000 …near 'a.suspend = 1' at line 1
> exec ^
: string
| set a.swaps = 1
< error
| 1064 42000 …near 'a.swaps = 1' at line 1
> exec ^
: string
| set a.switches = 1
< error
| 1064 42000 …near 'a.switches = 1' at line 1
> exec ^
: string
| set a.sysdate = 1
< error
| 1064 42000 …near 'a.sysdate = 1' at line 1
> exec ^
: string
| set a.system = 1
< error
| 1064 42000 …near 'a.system = 1' at line 1
> exec ^
: string
| set a.system_user = 1
< error
| 1064 42000 …near 'a.system_user = 1' at line 1
> exec ^
: string
| set a.table = 1
< error
| 1064 42000 …near 'a.table = 1' at line 1
> exec ^
: string
| set a.table_checksum = 1
< error
| 1064 42000 …near 'a.table_checksum = 1' at line 1
> exec ^
: string
| set a.table_name = 1
< error
| 1064 42000 …near 'a.table_name = 1' at line 1
> exec ^
: string
| set a.tables = 1
< error
| 1064 42000 …near 'a.tables = 1' at line 1
> exec ^
: string
| set a.tablesample = 1
< error
| 1064 42000 …near 'a.tablesample = 1' at line 1
> exec ^
: string
| set a.tablespace = 1
< error
| 1064 42000 …near 'a.tablespace = 1' at line 1
> exec ^
: string
| set a.tablet_size = 1
< error
| 1064 42000 …near 'a.tablet_size = 1' at line 1
> exec ^
: string
| set a.template = 1
< error
| 1064 42000 …near 'a.template = 1' at line 1
> exec ^
: string
| set a.temporary = 1
< error
| 1064 42000 …near 'a.temporary = 1' at line 1
> exec ^
: string
| set a.temptable = 1
< error
| 1064 42000 …near 'a.temptable = 1' at line 1
> exec ^
: string
| set a.terminated = 1
< error
| 1064 42000 …near 'a.terminated = 1' at line 1
> exec ^
: string
| set a.text = 1
< error
| 1064 42000 …near 'a.text = 1' at line 1
> exec ^
: string
| set a.than = 1
< error
| 1064 42000 …near 'a.than = 1' at line 1
> exec ^
: string
| set a.then = 1
< error
| 1064 42000 …near 'a.then = 1' at line 1
> exec ^
: string
| set a.thread_priority = 1
< error
| 1064 42000 …near 'a.thread_priority = 1' at line 1
> exec ^
: string
| set a.ties = 1
< error
| 1064 42000 …near 'a.ties = 1' at line 1
> exec ^
: string
| set a.time = 1
< error
| 1064 42000 …near 'a.time = 1' at line 1
> exec ^
: string
| set a.timestamp = 1
< error
| 1064 42000 …near 'a.timestamp = 1' at line 1
> exec ^
: string
| set a.timestampadd = 1
< error
| 1064 42000 …near 'a.timestampadd = 1' at line 1
> exec ^
: string
| set a.timestampdiff = 1
< error
| 1064 42000 …near 'a.timestampdiff = 1' at line 1
> exec ^
: string
| set a.tinyblob = 1
< error
| 1064 42000 …near 'a.tinyblob = 1' at line 1
> exec ^
: string
| set a.tinyint = 1
< error
| 1064 42000 …near 'a.tinyint = 1' at line 1
> exec ^
: string
| set a.tinytext = 1
< error
| 1064 42000 …near 'a.tinytext = 1' at line 1
> exec ^
: string
| set a.tls = 1
< error
| 1064 42000 …near 'a.tls = 1' at line 1
> exec ^
: string
| set a.to = 1
< error
| 1064 42000 …near 'a.to = 1' at line 1
> exec ^
: string
| set a.trailing = 1
< error
| 1064 42000 …near 'a.trailing = 1' at line 1
> exec ^
: string
| set a.transaction = 1
< error
| 1064 42000 …near 'a.transaction = 1' at line 1
> exec ^
: string
| set a.trigger = 1
< error
| 1064 42000 …near 'a.trigger = 1' at line 1
> exec ^
: string
| set a.triggers = 1
< error
| 1064 42000 …near 'a.triggers = 1' at line 1
> exec ^
: string
| set a.trim = 1
< error
| 1064 42000 …near 'a.trim = 1' at line 1
> exec ^
: string
| set a.true = 1
< error
| 1064 42000 …near 'a.true = 1' at line 1
> exec ^
: string
| set a.truncate = 1
< error
| 1064 42000 …near 'a.truncate = 1' at line 1
> exec ^
: string
| set a.type = 1
< error
| 1064 42000 …near 'a.type = 1' at line 1
> exec ^
: string
| set a.types = 1
< error
| 1064 42000 …near 'a.types = 1' at line 1
> exec ^
: string
| set a.unbounded = 1
< error
| 1064 42000 …near 'a.unbounded = 1' at line 1
> exec ^
: string
| set a.uncommitted = 1
< error
| 1064 42000 …near 'a.uncommitted = 1' at line 1
> exec ^
: string
| set a.undefined = 1
< error
| 1064 42000 …near 'a.undefined = 1' at line 1
> exec ^
: string
| set a.undo = 1
< error
| 1064 42000 …near 'a.undo = 1' at line 1
> exec ^
: string
| set a.undo_buffer_size = 1
< error
| 1064 42000 …near 'a.undo_buffer_size = 1' at line 1
> exec ^
: string
| set a.undofile = 1
< error
| 1064 42000 …near 'a.undofile = 1' at line 1
> exec ^
: string
| set a.unicode = 1
< error
| 1064 42000 …near 'a.unicode = 1' at line 1
> exec ^
: string
| set a.uninstall = 1
< error
| 1064 42000 …near 'a.uninstall = 1' at line 1
> exec ^
: string
| set a.union = 1
< error
| 1064 42000 …near 'a.union = 1' at line 1
> exec ^
: string
| set a.unique = 1
< error
| 1064 42000 …near 'a.unique = 1' at line 1
> exec ^
: string
| set a.unknown = 1
< error
| 1064 42000 …near 'a.unknown = 1' at line 1
> exec ^
: string
| set a.unlock = 1
< error
| 1064 42000 …near 'a.unlock = 1' at line 1
> exec ^
: string
| set a.unregister = 1
< error
| 1064 42000 …near 'a.unregister = 1' at line 1
> exec ^
: string
| set a.unsigned = 1
< error
| 1064 42000 …near 'a.unsigned = 1' at line 1
> exec ^
: string
| set a.until = 1
< error
| 1064 42000 …near 'a.until = 1' at line 1
> exec ^
: string
| set a.update = 1
< error
| 1064 42000 …near 'a.update = 1' at line 1
> exec ^
: string
| set a.upgrade = 1
< error
| 1064 42000 …near 'a.upgrade = 1' at line 1
> exec ^
: string
| set a.uri = 1
< error
| 1064 42000 …near 'a.uri = 1' at line 1
> exec ^
: string
| set a.url = 1
< error
| 1064 42000 …near 'a.url = 1' at line 1
> exec ^
: string
| set a.usage = 1
< error
| 1064 42000 …near 'a.usage = 1' at line 1
> exec ^
: string
| set a.use = 1
< error
| 1064 42000 …near 'a.use = 1' at line 1
> exec ^
: string
| set a.use_bloom_filter = 1
< error
| 1064 42000 …near 'a.use_bloom_filter = 1' at line 1
> exec ^
: string
| set a.use_frm = 1
< error
| 1064 42000 …near 'a.use_frm = 1' at line 1
> exec ^
: string
| set a.user = 1
< error
| 1064 42000 …near 'a.user = 1' at line 1
> exec ^
: string
| set a.user_resources = 1
< error
| 1064 42000 …near 'a.user_resources = 1' at line 1
> exec ^
: string
| set a.using = 1
< error
| 1064 42000 …near 'a.using = 1' at line 1
> exec ^
: string
| set a.utc_date = 1
< error
| 1064 42000 …near 'a.utc_date = 1' at line 1
> exec ^
: string
| set a.utc_time = 1
< error
| 1064 42000 …near 'a.utc_time = 1' at line 1
> exec ^
: string
| set a.utc_timestamp = 1
< error
| 1064 42000 …near 'a.utc_timestamp = 1' at line 1
> exec ^
: string
| set a.validate = 1
< error
| 1064 42000 …near 'a.validate = 1' at line 1
> exec ^
: string
| set a.validation = 1
< error
| 1064 42000 …near 'a.validation = 1' at line 1
> exec ^
: string
| set a.value = 1
< error
| 1064 42000 …near 'a.value = 1' at line 1
> exec ^
: string
| set a.values = 1
< error
| 1064 42000 …near 'a.values = 1' at line 1
> exec ^
: string
| set a.var_pop = 1
< error
| 1064 42000 …near 'a.var_pop = 1' at line 1
> exec ^
: string
| set a.var_samp = 1
< error
| 1064 42000 …near 'a.var_samp = 1' at line 1
> exec ^
: string
| set a.varbinary = 1
< error
| 1064 42000 …near 'a.varbinary = 1' at line 1
> exec ^
: string
| set a.varchar = 1
< error
| 1064 42000 …near 'a.varchar = 1' at line 1
> exec ^
: string
| set a.varcharacter = 1
< error
| 1064 42000 …near 'a.varcharacter = 1' at line 1
> exec ^
: string
| set a.variables = 1
< error
| 1064 42000 …near 'a.variables = 1' at line 1
> exec ^
: string
| set a.variance = 1
< error
| 1064 42000 …near 'a.variance = 1' at line 1
> exec ^
: string
| set a.varying = 1
< error
| 1064 42000 …near 'a.varying = 1' at line 1
> exec ^
: string
| set a.vcpu = 1
< error
| 1064 42000 …near 'a.vcpu = 1' at line 1
> exec ^
: string
| set a.vector = 1
< error
| 1064 42000 …near 'a.vector = 1' at line 1
> exec ^
: string
| set a.verify_key_constraints = 1
< error
| 1064 42000 …near 'a.verify_key_constraints = 1' at line 1
> exec ^
: string
| set a.view = 1
< error
| 1064 42000 …near 'a.view = 1' at line 1
> exec ^
: string
| set a.virtual = 1
< error
| 1064 42000 …near 'a.virtual = 1' at line 1
> exec ^
: string
| set a.visible = 1
< error
| 1064 42000 …near 'a.visible = 1' at line 1
> exec ^
: string
| set a.wait = 1
< error
| 1064 42000 …near 'a.wait = 1' at line 1
> exec ^
: string
| set a.warnings = 1
< error
| 1064 42000 …near 'a.warnings = 1' at line 1
> exec ^
: string
| set a.week = 1
< error
| 1064 42000 …near 'a.week = 1' at line 1
> exec ^
: string
| set a.weight_string = 1
< error
| 1064 42000 …near 'a.weight_string = 1' at line 1
> exec ^
: string
| set a.when = 1
< error
| 1064 42000 …near 'a.when = 1' at line 1
> exec ^
: string
| set a.where = 1
< error
| 1064 42000 …near 'a.where = 1' at line 1
> exec ^
: string
| set a.while = 1
< error
| 1064 42000 …near 'a.while = 1' at line 1
> exec ^
: string
| set a.window = 1
< error
| 1064 42000 …near 'a.window = 1' at line 1
> exec ^
: string
| set a.with = 1
< error
| 1064 42000 …near 'a.with = 1' at line 1
> exec ^
: string
| set a.without = 1
< error
| 1064 42000 …near 'a.without = 1' at line 1
> exec ^
: string
| set a.work = 1
< error
| 1064 42000 …near 'a.work = 1' at line 1
> exec ^
: string
| set a.wrapper = 1
< error
| 1064 42000 …near 'a.wrapper = 1' at line 1
> exec ^
: string
| set a.write = 1
< error
| 1064 42000 …near 'a.write = 1' at line 1
> exec ^
: string
| set a.x509 = 1
< error
| 1064 42000 …near 'a.x509 = 1' at line 1
> exec ^
: string
| set a.xa = 1
< error
| 1064 42000 …near 'a.xa = 1' at line 1
> exec ^
: string
| set a.xid = 1
< error
| 1064 42000 …near 'a.xid = 1' at line 1
> exec ^
: string
| set a.xml = 1
< error
| 1064 42000 …near 'a.xml = 1' at line 1
> exec ^
: string
| set a.xor = 1
< error
| 1064 42000 …near 'a.xor = 1' at line 1
> exec ^
: string
| set a.year = 1
< error
| 1064 42000 …near 'a.year = 1' at line 1
> exec ^
: string
| set a.year_month = 1
< error
| 1064 42000 …near 'a.year_month = 1' at line 1
> exec ^
: string
| set a.zerofill = 1
< error
| 1064 42000 …near 'a.zerofill = 1' at line 1
> exec ^
: string
| set a.zone = 1
< error
| 1064 42000 …near 'a.zone = 1' at line 1
> exec ^
: string
| set @a = 1, absent = 2
< ok 0 0
> exec ^
: string
| set @a = 1, accessible = 2
< error
| 1064 42000 …near 'accessible = 2' at line 1
> exec ^
: string
| set @a = 1, account = 2
< ok 0 0
> exec ^
: string
| set @a = 1, action = 2
< ok 0 0
> exec ^
: string
| set @a = 1, active = 2
< ok 0 0
> exec ^
: string
| set @a = 1, add = 2
< error
| 1064 42000 …near 'add = 2' at line 1
> exec ^
: string
| set @a = 1, adddate = 2
< ok 0 0
> exec ^
: string
| set @a = 1, admin = 2
< ok 0 0
> exec ^
: string
| set @a = 1, after = 2
< ok 0 0
> exec ^
: string
| set @a = 1, against = 2
< ok 0 0
> exec ^
: string
| set @a = 1, aggregate = 2
< ok 0 0
> exec ^
: string
| set @a = 1, algorithm = 2
< ok 0 0
> exec ^
: string
| set @a = 1, all = 2
< error
| 1064 42000 …near 'all = 2' at line 1
> exec ^
: string
| set @a = 1, allow_missing_files = 2
< ok 0 0
> exec ^
: string
| set @a = 1, alter = 2
< error
| 1064 42000 …near 'alter = 2' at line 1
> exec ^
: string
| set @a = 1, always = 2
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set @a = 1, array = 2
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set @a = 1, asensitive = 2
< error
| 1064 42000 …near 'asensitive = 2' at line 1
> exec ^
: string
| set @a = 1, assign_gtids_to_anonymous_transactions = 2
< ok 0 0
> exec ^
: string
| set @a = 1, at = 2
< ok 0 0
> exec ^
: string
| set @a = 1, attribute = 2
< ok 0 0
> exec ^
: string
| set @a = 1, authentication = 2
< ok 0 0
> exec ^
: string
| set @a = 1, auto = 2
< ok 0 0
> exec ^
: string
| set @a = 1, auto_increment = 2
< ok 0 0
> exec ^
: string
| set @a = 1, auto_increment_mode = 2
< ok 0 0
> exec ^
: string
| set @a = 1, auto_refresh = 2
< ok 0 0
> exec ^
: string
| set @a = 1, auto_refresh_source = 2
< ok 0 0
> exec ^
: string
| set @a = 1, autoextend_size = 2
< ok 0 0
> exec ^
: string
| set @a = 1, avg = 2
< ok 0 0
> exec ^
: string
| set @a = 1, avg_row_length = 2
< ok 0 0
> exec ^
: string
| set @a = 1, backup = 2
< ok 0 0
> exec ^
: string
| set @a = 1, before = 2
< error
| 1064 42000 …near 'before = 2' at line 1
> exec ^
: string
| set @a = 1, begin = 2
< ok 0 0
> exec ^
: string
| set @a = 1, bernoulli = 2
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set @a = 1, bit = 2
< ok 0 0
> exec ^
: string
| set @a = 1, bit_and = 2
< ok 0 0
> exec ^
: string
| set @a = 1, bit_or = 2
< ok 0 0
> exec ^
: string
| set @a = 1, bit_xor = 2
< ok 0 0
> exec ^
: string
| set @a = 1, blob = 2
< error
| 1064 42000 …near 'blob = 2' at line 1
> exec ^
: string
| set @a = 1, block = 2
< ok 0 0
> exec ^
: string
| set @a = 1, block_size = 2
< ok 0 0
> exec ^
: string
| set @a = 1, bool = 2
< ok 0 0
> exec ^
: string
| set @a = 1, boolean = 2
< ok 0 0
> exec ^
: string
| set @a = 1, both = 2
< error
| 1064 42000 …near 'both = 2' at line 1
> exec ^
: string
| set @a = 1, btree = 2
< ok 0 0
> exec ^
: string
| set @a = 1, buckets = 2
< ok 0 0
> exec ^
: string
| set @a = 1, bulk = 2
< ok 0 0
> exec ^
: string
| set @a = 1, by = 2
< error
| 1064 42000 …near 'by = 2' at line 1
> exec ^
: string
| set @a = 1, byte = 2
< ok 0 0
> exec ^
: string
| set @a = 1, cache = 2
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set @a = 1, case = 2
< error
| 1064 42000 …near 'case = 2' at line 1
> exec ^
: string
| set @a = 1, cast = 2
< ok 0 0
> exec ^
: string
| set @a = 1, catalog_name = 2
< ok 0 0
> exec ^
: string
| set @a = 1, chain = 2
< ok 0 0
> exec ^
: string
| set @a = 1, challenge_response = 2
< ok 0 0
> exec ^
: string
| set @a = 1, change = 2
< error
| 1064 42000 …near 'change = 2' at line 1
> exec ^
: string
| set @a = 1, changed = 2
< ok 0 0
> exec ^
: string
| set @a = 1, channel = 2
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set @a = 1, check = 2
< error
| 1064 42000 …near 'check = 2' at line 1
> exec ^
: string
| set @a = 1, checksum = 2
< ok 0 0
> exec ^
: string
| set @a = 1, cipher = 2
< ok 0 0
> exec ^
: string
| set @a = 1, class_origin = 2
< ok 0 0
> exec ^
: string
| set @a = 1, client = 2
< ok 0 0
> exec ^
: string
| set @a = 1, clone = 2
< ok 0 0
> exec ^
: string
| set @a = 1, close = 2
< ok 0 0
> exec ^
: string
| set @a = 1, coalesce = 2
< ok 0 0
> exec ^
: string
| set @a = 1, code = 2
< ok 0 0
> exec ^
: string
| set @a = 1, collate = 2
< error
| 1064 42000 …near 'collate = 2' at line 1
> exec ^
: string
| set @a = 1, collation = 2
< ok 0 0
> exec ^
: string
| set @a = 1, column = 2
< error
| 1064 42000 …near 'column = 2' at line 1
> exec ^
: string
| set @a = 1, column_format = 2
< ok 0 0
> exec ^
: string
| set @a = 1, column_name = 2
< ok 0 0
> exec ^
: string
| set @a = 1, columns = 2
< ok 0 0
> exec ^
: string
| set @a = 1, comment = 2
< ok 0 0
> exec ^
: string
| set @a = 1, commit = 2
< ok 0 0
> exec ^
: string
| set @a = 1, committed = 2
< ok 0 0
> exec ^
: string
| set @a = 1, compact = 2
< ok 0 0
> exec ^
: string
| set @a = 1, completion = 2
< ok 0 0
> exec ^
: string
| set @a = 1, component = 2
< ok 0 0
> exec ^
: string
| set @a = 1, compressed = 2
< ok 0 0
> exec ^
: string
| set @a = 1, compression = 2
< ok 0 0
> exec ^
: string
| set @a = 1, concat = 2
< ok 0 0
> exec ^
: string
| set @a = 1, concurrent = 2
< ok 0 0
> exec ^
: string
| set @a = 1, condition = 2
< error
| 1064 42000 …near 'condition = 2' at line 1
> exec ^
: string
| set @a = 1, connection = 2
< ok 0 0
> exec ^
: string
| set @a = 1, consistent = 2
< ok 0 0
> exec ^
: string
| set @a = 1, constraint = 2
< error
| 1064 42000 …near 'constraint = 2' at line 1
> exec ^
: string
| set @a = 1, constraint_catalog = 2
< ok 0 0
> exec ^
: string
| set @a = 1, constraint_name = 2
< ok 0 0
> exec ^
: string
| set @a = 1, constraint_schema = 2
< ok 0 0
> exec ^
: string
| set @a = 1, contains = 2
< ok 0 0
> exec ^
: string
| set @a = 1, context = 2
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set @a = 1, cpu = 2
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set @a = 1, cube = 2
< ok 0 0
> exec ^
: string
| set @a = 1, cume_dist = 2
< ok 0 0
> exec ^
: string
| set @a = 1, curdate = 2
< ok 0 0
> exec ^
: string
| set @a = 1, current = 2
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set @a = 1, curtime = 2
< ok 0 0
> exec ^
: string
| set @a = 1, data = 2
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set @a = 1, date = 2
< ok 0 0
> exec ^
: string
| set @a = 1, date_add = 2
< ok 0 0
> exec ^
: string
| set @a = 1, date_sub = 2
< ok 0 0
> exec ^
: string
| set @a = 1, datetime = 2
< ok 0 0
> exec ^
: string
| set @a = 1, day = 2
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set @a = 1, dec = 2
< ok 0 0
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
| 1064 42000 …near 'default = 2' at line 1
> exec ^
: string
| set @a = 1, default_auth = 2
< ok 0 0
> exec ^
: string
| set @a = 1, definer = 2
< ok 0 0
> exec ^
: string
| set @a = 1, definition = 2
< ok 0 0
> exec ^
: string
| set @a = 1, delay_key_write = 2
< ok 0 0
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
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set @a = 1, deterministic = 2
< error
| 1064 42000 …near 'deterministic = 2' at line 1
> exec ^
: string
| set @a = 1, diagnostics = 2
< ok 0 0
> exec ^
: string
| set @a = 1, directory = 2
< ok 0 0
> exec ^
: string
| set @a = 1, disable = 2
< ok 0 0
> exec ^
: string
| set @a = 1, discard = 2
< ok 0 0
> exec ^
: string
| set @a = 1, disk = 2
< ok 0 0
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
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set @a = 1, dumpfile = 2
< ok 0 0
> exec ^
: string
| set @a = 1, duplicate = 2
< ok 0 0
> exec ^
: string
| set @a = 1, dynamic = 2
< ok 0 0
> exec ^
: string
| set @a = 1, dynamic_partition_policy = 2
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set @a = 1, enable = 2
< ok 0 0
> exec ^
: string
| set @a = 1, enable_macro_block_bloom_filter = 2
< ok 0 0
> exec ^
: string
| set @a = 1, enclosed = 2
< error
| 1064 42000 …near 'enclosed = 2' at line 1
> exec ^
: string
| set @a = 1, encryption = 2
< ok 0 0
> exec ^
: string
| set @a = 1, end = 2
< ok 0 0
> exec ^
: string
| set @a = 1, ends = 2
< ok 0 0
> exec ^
: string
| set @a = 1, enforced = 2
< ok 0 0
> exec ^
: string
| set @a = 1, engine = 2
< ok 0 0
> exec ^
: string
| set @a = 1, engine_attribute = 2
< ok 0 0
> exec ^
: string
| set @a = 1, engines = 2
< ok 0 0
> exec ^
: string
| set @a = 1, enum = 2
< ok 0 0
> exec ^
: string
| set @a = 1, error = 2
< ok 0 0
> exec ^
: string
| set @a = 1, errors = 2
< ok 0 0
> exec ^
: string
| set @a = 1, escape = 2
< ok 0 0
> exec ^
: string
| set @a = 1, escaped = 2
< error
| 1064 42000 …near 'escaped = 2' at line 1
> exec ^
: string
| set @a = 1, event = 2
< error
| 1064 42000 …near 'event = 2' at line 1
> exec ^
: string
| set @a = 1, events = 2
< ok 0 0
> exec ^
: string
| set @a = 1, every = 2
< ok 0 0
> exec ^
: string
| set @a = 1, except = 2
< ok 0 0
> exec ^
: string
| set @a = 1, exchange = 2
< ok 0 0
> exec ^
: string
| set @a = 1, exclude = 2
< ok 0 0
> exec ^
: string
| set @a = 1, execute = 2
< error
| 1064 42000 …near 'execute = 2' at line 1
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
< ok 0 0
> exec ^
: string
| set @a = 1, expire = 2
< ok 0 0
> exec ^
: string
| set @a = 1, explain = 2
< error
| 1064 42000 …near 'explain = 2' at line 1
> exec ^
: string
| set @a = 1, export = 2
< ok 0 0
> exec ^
: string
| set @a = 1, extended = 2
< ok 0 0
> exec ^
: string
| set @a = 1, extent_size = 2
< ok 0 0
> exec ^
: string
| set @a = 1, external = 2
< ok 0 0
> exec ^
: string
| set @a = 1, external_format = 2
< ok 0 0
> exec ^
: string
| set @a = 1, extract = 2
< ok 0 0
> exec ^
: string
| set @a = 1, factor = 2
< ok 0 0
> exec ^
: string
| set @a = 1, failed_login_attempts = 2
< ok 0 0
> exec ^
: string
| set @a = 1, false = 2
< error
| 1064 42000 …near 'false = 2' at line 1
> exec ^
: string
| set @a = 1, fast = 2
< ok 0 0
> exec ^
: string
| set @a = 1, faults = 2
< ok 0 0
> exec ^
: string
| set @a = 1, fetch = 2
< error
| 1064 42000 …near 'fetch = 2' at line 1
> exec ^
: string
| set @a = 1, fields = 2
< ok 0 0
> exec ^
: string
| set @a = 1, file = 2
< error
| 1064 42000 …near 'file = 2' at line 1
> exec ^
: string
| set @a = 1, file_block_size = 2
< ok 0 0
> exec ^
: string
| set @a = 1, file_format = 2
< ok 0 0
> exec ^
: string
| set @a = 1, file_name = 2
< ok 0 0
> exec ^
: string
| set @a = 1, file_pattern = 2
< ok 0 0
> exec ^
: string
| set @a = 1, file_prefix = 2
< ok 0 0
> exec ^
: string
| set @a = 1, files = 2
< ok 0 0
> exec ^
: string
| set @a = 1, filter = 2
< ok 0 0
> exec ^
: string
| set @a = 1, finish = 2
< ok 0 0
> exec ^
: string
| set @a = 1, first = 2
< ok 0 0
> exec ^
: string
| set @a = 1, first_value = 2
< ok 0 0
> exec ^
: string
| set @a = 1, fixed = 2
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set @a = 1, following = 2
< ok 0 0
> exec ^
: string
| set @a = 1, follows = 2
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set @a = 1, found = 2
< ok 0 0
> exec ^
: string
| set @a = 1, from = 2
< error
| 1064 42000 …near 'from = 2' at line 1
> exec ^
: string
| set @a = 1, full = 2
< ok 0 0
> exec ^
: string
| set @a = 1, fulltext = 2
< error
| 1064 42000 …near 'fulltext = 2' at line 1
> exec ^
: string
| set @a = 1, function = 2
< ok 0 0
> exec ^
: string
| set @a = 1, general = 2
< ok 0 0
> exec ^
: string
| set @a = 1, generate = 2
< ok 0 0
> exec ^
: string
| set @a = 1, generated = 2
< error
| 1064 42000 …near 'generated = 2' at line 1
> exec ^
: string
| set @a = 1, geomcollection = 2
< ok 0 0
> exec ^
: string
| set @a = 1, geometry = 2
< ok 0 0
> exec ^
: string
| set @a = 1, geometrycollection = 2
< ok 0 0
> exec ^
: string
| set @a = 1, get = 2
< error
| 1064 42000 …near 'get = 2' at line 1
> exec ^
: string
| set @a = 1, get_format = 2
< ok 0 0
> exec ^
: string
| set @a = 1, get_master_public_key = 2
< ok 0 0
> exec ^
: string
| set @a = 1, get_source_public_key = 2
< ok 0 0
> exec ^
: string
| set @a = 1, global = 2
< ok 0 0
> exec ^
: string
| set @a = 1, grant = 2
< error
| 1064 42000 …near 'grant = 2' at line 1
> exec ^
: string
| set @a = 1, grants = 2
< ok 0 0
> exec ^
: string
| set @a = 1, greatest = 2
< ok 0 0
> exec ^
: string
| set @a = 1, group = 2
< error
| 1064 42000 …near 'group = 2' at line 1
> exec ^
: string
| set @a = 1, group_concat = 2
< ok 0 0
> exec ^
: string
| set @a = 1, group_replication = 2
< ok 0 0
> exec ^
: string
| set @a = 1, grouping = 2
< ok 0 0
> exec ^
: string
| set @a = 1, groups = 2
< ok 0 0
> exec ^
: string
| set @a = 1, gtid_only = 2
< ok 0 0
> exec ^
: string
| set @a = 1, gtids = 2
< ok 0 0
> exec ^
: string
| set @a = 1, guided = 2
< ok 0 0
> exec ^
: string
| set @a = 1, handler = 2
< ok 0 0
> exec ^
: string
| set @a = 1, hash = 2
< ok 0 0
> exec ^
: string
| set @a = 1, having = 2
< error
| 1064 42000 …near 'having = 2' at line 1
> exec ^
: string
| set @a = 1, header = 2
< ok 0 0
> exec ^
: string
| set @a = 1, help = 2
< ok 0 0
> exec ^
: string
| set @a = 1, high_priority = 2
< error
| 1064 42000 …near 'high_priority = 2' at line 1
> exec ^
: string
| set @a = 1, histogram = 2
< ok 0 0
> exec ^
: string
| set @a = 1, history = 2
< ok 0 0
> exec ^
: string
| set @a = 1, host = 2
< ok 0 0
> exec ^
: string
| set @a = 1, hosts = 2
< ok 0 0
> exec ^
: string
| set @a = 1, hour = 2
< ok 0 0
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
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set @a = 1, import = 2
< ok 0 0
> exec ^
: string
| set @a = 1, in = 2
< error
| 1064 42000 …near 'in = 2' at line 1
> exec ^
: string
| set @a = 1, inactive = 2
< ok 0 0
> exec ^
: string
| set @a = 1, index = 2
< error
| 1064 42000 …near 'index = 2' at line 1
> exec ^
: string
| set @a = 1, indexes = 2
< ok 0 0
> exec ^
: string
| set @a = 1, infile = 2
< error
| 1064 42000 …near 'infile = 2' at line 1
> exec ^
: string
| set @a = 1, initial = 2
< ok 0 0
> exec ^
: string
| set @a = 1, initial_size = 2
< ok 0 0
> exec ^
: string
| set @a = 1, initiate = 2
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set @a = 1, install = 2
< ok 0 0
> exec ^
: string
| set @a = 1, instance = 2
< ok 0 0
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
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set @a = 1, invoker = 2
< ok 0 0
> exec ^
: string
| set @a = 1, io = 2
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set @a = 1, ipc = 2
< ok 0 0
> exec ^
: string
| set @a = 1, is = 2
< error
| 1064 42000 …near 'is = 2' at line 1
> exec ^
: string
| set @a = 1, isolation = 2
< ok 0 0
> exec ^
: string
| set @a = 1, issuer = 2
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set @a = 1, json_arrayagg = 2
< ok 0 0
> exec ^
: string
| set @a = 1, json_duality_object = 2
< ok 0 0
> exec ^
: string
| set @a = 1, json_objectagg = 2
< ok 0 0
> exec ^
: string
| set @a = 1, json_table = 2
< ok 0 0
> exec ^
: string
| set @a = 1, json_value = 2
< ok 0 0
> exec ^
: string
| set @a = 1, key = 2
< error
| 1064 42000 …near 'key = 2' at line 1
> exec ^
: string
| set @a = 1, key_block_size = 2
< ok 0 0
> exec ^
: string
| set @a = 1, keyring = 2
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set @a = 1, language = 2
< ok 0 0
> exec ^
: string
| set @a = 1, last = 2
< ok 0 0
> exec ^
: string
| set @a = 1, last_value = 2
< ok 0 0
> exec ^
: string
| set @a = 1, lateral = 2
< ok 0 0
> exec ^
: string
| set @a = 1, lead = 2
< ok 0 0
> exec ^
: string
| set @a = 1, leading = 2
< error
| 1064 42000 …near 'leading = 2' at line 1
> exec ^
: string
| set @a = 1, least = 2
< ok 0 0
> exec ^
: string
| set @a = 1, leave = 2
< error
| 1064 42000 …near 'leave = 2' at line 1
> exec ^
: string
| set @a = 1, leaves = 2
< ok 0 0
> exec ^
: string
| set @a = 1, left = 2
< error
| 1064 42000 …near 'left = 2' at line 1
> exec ^
: string
| set @a = 1, less = 2
< ok 0 0
> exec ^
: string
| set @a = 1, level = 2
< ok 0 0
> exec ^
: string
| set @a = 1, library = 2
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set @a = 1, list = 2
< ok 0 0
> exec ^
: string
| set @a = 1, load = 2
< error
| 1064 42000 …near 'load = 2' at line 1
> exec ^
: string
| set @a = 1, local = 2
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set @a = 1, locks = 2
< ok 0 0
> exec ^
: string
| set @a = 1, log = 2
< ok 0 0
> exec ^
: string
| set @a = 1, logfile = 2
< ok 0 0
> exec ^
: string
| set @a = 1, logs = 2
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set @a = 1, manual = 2
< ok 0 0
> exec ^
: string
| set @a = 1, master = 2
< ok 0 0
> exec ^
: string
| set @a = 1, master_auto_position = 2
< ok 0 0
> exec ^
: string
| set @a = 1, master_bind = 2
< error
| 1064 42000 …near 'master_bind = 2' at line 1
> exec ^
: string
| set @a = 1, master_compression_algorithms = 2
< ok 0 0
> exec ^
: string
| set @a = 1, master_connect_retry = 2
< ok 0 0
> exec ^
: string
| set @a = 1, master_delay = 2
< ok 0 0
> exec ^
: string
| set @a = 1, master_heartbeat_period = 2
< ok 0 0
> exec ^
: string
| set @a = 1, master_host = 2
< ok 0 0
> exec ^
: string
| set @a = 1, master_log_file = 2
< ok 0 0
> exec ^
: string
| set @a = 1, master_log_pos = 2
< ok 0 0
> exec ^
: string
| set @a = 1, master_password = 2
< ok 0 0
> exec ^
: string
| set @a = 1, master_port = 2
< ok 0 0
> exec ^
: string
| set @a = 1, master_public_key_path = 2
< ok 0 0
> exec ^
: string
| set @a = 1, master_retry_count = 2
< ok 0 0
> exec ^
: string
| set @a = 1, master_ssl = 2
< ok 0 0
> exec ^
: string
| set @a = 1, master_ssl_ca = 2
< ok 0 0
> exec ^
: string
| set @a = 1, master_ssl_capath = 2
< ok 0 0
> exec ^
: string
| set @a = 1, master_ssl_cert = 2
< ok 0 0
> exec ^
: string
| set @a = 1, master_ssl_cipher = 2
< ok 0 0
> exec ^
: string
| set @a = 1, master_ssl_crl = 2
< ok 0 0
> exec ^
: string
| set @a = 1, master_ssl_crlpath = 2
< ok 0 0
> exec ^
: string
| set @a = 1, master_ssl_key = 2
< ok 0 0
> exec ^
: string
| set @a = 1, master_ssl_verify_server_cert = 2
< error
| 1064 42000 …near 'master_ssl_verify_server_cert = 2' at line 1
> exec ^
: string
| set @a = 1, master_tls_ciphersuites = 2
< ok 0 0
> exec ^
: string
| set @a = 1, master_tls_version = 2
< ok 0 0
> exec ^
: string
| set @a = 1, master_user = 2
< ok 0 0
> exec ^
: string
| set @a = 1, master_zstd_compression_level = 2
< ok 0 0
> exec ^
: string
| set @a = 1, match = 2
< error
| 1064 42000 …near 'match = 2' at line 1
> exec ^
: string
| set @a = 1, materialized = 2
< ok 0 0
> exec ^
: string
| set @a = 1, max = 2
< ok 0 0
> exec ^
: string
| set @a = 1, max_connections_per_hour = 2
< ok 0 0
> exec ^
: string
| set @a = 1, max_queries_per_hour = 2
< ok 0 0
> exec ^
: string
| set @a = 1, max_rows = 2
< ok 0 0
> exec ^
: string
| set @a = 1, max_size = 2
< ok 0 0
> exec ^
: string
| set @a = 1, max_updates_per_hour = 2
< ok 0 0
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
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set @a = 1, memory = 2
< ok 0 0
> exec ^
: string
| set @a = 1, merge = 2
< ok 0 0
> exec ^
: string
| set @a = 1, message_text = 2
< ok 0 0
> exec ^
: string
| set @a = 1, microsecond = 2
< ok 0 0
> exec ^
: string
| set @a = 1, mid = 2
< ok 0 0
> exec ^
: string
| set @a = 1, middleint = 2
< error
| 1064 42000 …near 'middleint = 2' at line 1
> exec ^
: string
| set @a = 1, migrate = 2
< ok 0 0
> exec ^
: string
| set @a = 1, min = 2
< ok 0 0
> exec ^
: string
| set @a = 1, min_rows = 2
< ok 0 0
> exec ^
: string
| set @a = 1, minute = 2
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set @a = 1, modifies = 2
< error
| 1064 42000 …near 'modifies = 2' at line 1
> exec ^
: string
| set @a = 1, modify = 2
< ok 0 0
> exec ^
: string
| set @a = 1, month = 2
< ok 0 0
> exec ^
: string
| set @a = 1, multilinestring = 2
< ok 0 0
> exec ^
: string
| set @a = 1, multipoint = 2
< ok 0 0
> exec ^
: string
| set @a = 1, multipolygon = 2
< ok 0 0
> exec ^
: string
| set @a = 1, mutex = 2
< ok 0 0
> exec ^
: string
| set @a = 1, mysql_errno = 2
< ok 0 0
> exec ^
: string
| set @a = 1, name = 2
< ok 0 0
> exec ^
: string
| set @a = 1, names = 2
< ok 0 0
> exec ^
: string
| set @a = 1, national = 2
< ok 0 0
> exec ^
: string
| set @a = 1, natural = 2
< error
| 1064 42000 …near 'natural = 2' at line 1
> exec ^
: string
| set @a = 1, nchar = 2
< ok 0 0
> exec ^
: string
| set @a = 1, ndb = 2
< ok 0 0
> exec ^
: string
| set @a = 1, ndbcluster = 2
< ok 0 0
> exec ^
: string
| set @a = 1, nested = 2
< ok 0 0
> exec ^
: string
| set @a = 1, network_namespace = 2
< ok 0 0
> exec ^
: string
| set @a = 1, never = 2
< ok 0 0
> exec ^
: string
| set @a = 1, new = 2
< ok 0 0
> exec ^
: string
| set @a = 1, next = 2
< ok 0 0
> exec ^
: string
| set @a = 1, no = 2
< ok 0 0
> exec ^
: string
| set @a = 1, no_wait = 2
< ok 0 0
> exec ^
: string
| set @a = 1, no_write_to_binlog = 2
< error
| 1064 42000 …near 'no_write_to_binlog = 2' at line 1
> exec ^
: string
| set @a = 1, nodegroup = 2
< ok 0 0
> exec ^
: string
| set @a = 1, none = 2
< error
| 1064 42000 …near 'none = 2' at line 1
> exec ^
: string
| set @a = 1, not = 2
< error
| 1064 42000 …near 'not = 2' at line 1
> exec ^
: string
| set @a = 1, now = 2
< ok 0 0
> exec ^
: string
| set @a = 1, nowait = 2
< ok 0 0
> exec ^
: string
| set @a = 1, nth_value = 2
< ok 0 0
> exec ^
: string
| set @a = 1, ntile = 2
< ok 0 0
> exec ^
: string
| set @a = 1, null = 2
< error
| 1064 42000 …near 'null = 2' at line 1
> exec ^
: string
| set @a = 1, nullif = 2
< ok 0 0
> exec ^
: string
| set @a = 1, nulls = 2
< ok 0 0
> exec ^
: string
| set @a = 1, number = 2
< ok 0 0
> exec ^
: string
| set @a = 1, numeric = 2
< error
| 1064 42000 …near 'numeric = 2' at line 1
> exec ^
: string
| set @a = 1, nvarchar = 2
< ok 0 0
> exec ^
: string
| set @a = 1, of = 2
< ok 0 0
> exec ^
: string
| set @a = 1, off = 2
< ok 0 0
> exec ^
: string
| set @a = 1, offset = 2
< ok 0 0
> exec ^
: string
| set @a = 1, oj = 2
< ok 0 0
> exec ^
: string
| set @a = 1, old = 2
< ok 0 0
> exec ^
: string
| set @a = 1, on = 2
< error
| 1064 42000 …near 'on = 2' at line 1
> exec ^
: string
| set @a = 1, one = 2
< ok 0 0
> exec ^
: string
| set @a = 1, only = 2
< ok 0 0
> exec ^
: string
| set @a = 1, open = 2
< ok 0 0
> exec ^
: string
| set @a = 1, optimize = 2
< error
| 1064 42000 …near 'optimize = 2' at line 1
> exec ^
: string
| set @a = 1, optimizer_costs = 2
< ok 0 0
> exec ^
: string
| set @a = 1, option = 2
< error
| 1064 42000 …near 'option = 2' at line 1
> exec ^
: string
| set @a = 1, optional = 2
< ok 0 0
> exec ^
: string
| set @a = 1, optionally = 2
< error
| 1064 42000 …near 'optionally = 2' at line 1
> exec ^
: string
| set @a = 1, options = 2
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set @a = 1, organization = 2
< ok 0 0
> exec ^
: string
| set @a = 1, others = 2
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set @a = 1, overlay = 2
< ok 0 0
> exec ^
: string
| set @a = 1, overwrite = 2
< ok 0 0
> exec ^
: string
| set @a = 1, owner = 2
< ok 0 0
> exec ^
: string
| set @a = 1, pack_keys = 2
< ok 0 0
> exec ^
: string
| set @a = 1, page = 2
< ok 0 0
> exec ^
: string
| set @a = 1, parallel = 2
< ok 0 0
> exec ^
: string
| set @a = 1, parameters = 2
< ok 0 0
> exec ^
: string
| set @a = 1, parse_tree = 2
< ok 0 0
> exec ^
: string
| set @a = 1, parser = 2
< ok 0 0
> exec ^
: string
| set @a = 1, partial = 2
< ok 0 0
> exec ^
: string
| set @a = 1, partition = 2
< error
| 1064 42000 …near 'partition = 2' at line 1
> exec ^
: string
| set @a = 1, partitioning = 2
< ok 0 0
> exec ^
: string
| set @a = 1, partitions = 2
< ok 0 0
> exec ^
: string
| set @a = 1, password = 2
< error
| 1064 42000 …near 'password = 2' at line 1
> exec ^
: string
| set @a = 1, password_lock_time = 2
< ok 0 0
> exec ^
: string
| set @a = 1, path = 2
< ok 0 0
> exec ^
: string
| set @a = 1, pctfree = 2
< ok 0 0
> exec ^
: string
| set @a = 1, percent_rank = 2
< ok 0 0
> exec ^
: string
| set @a = 1, persist = 2
< ok 0 0
> exec ^
: string
| set @a = 1, persist_only = 2
< ok 0 0
> exec ^
: string
| set @a = 1, phase = 2
< ok 0 0
> exec ^
: string
| set @a = 1, placing = 2
< ok 0 0
> exec ^
: string
| set @a = 1, plugin = 2
< ok 0 0
> exec ^
: string
| set @a = 1, plugin_dir = 2
< ok 0 0
> exec ^
: string
| set @a = 1, plugins = 2
< ok 0 0
> exec ^
: string
| set @a = 1, point = 2
< ok 0 0
> exec ^
: string
| set @a = 1, polygon = 2
< ok 0 0
> exec ^
: string
| set @a = 1, port = 2
< ok 0 0
> exec ^
: string
| set @a = 1, position = 2
< ok 0 0
> exec ^
: string
| set @a = 1, precedes = 2
< ok 0 0
> exec ^
: string
| set @a = 1, preceding = 2
< ok 0 0
> exec ^
: string
| set @a = 1, precision = 2
< error
| 1064 42000 …near 'precision = 2' at line 1
> exec ^
: string
| set @a = 1, prefix = 2
< ok 0 0
> exec ^
: string
| set @a = 1, prepare = 2
< ok 0 0
> exec ^
: string
| set @a = 1, preserve = 2
< ok 0 0
> exec ^
: string
| set @a = 1, prev = 2
< ok 0 0
> exec ^
: string
| set @a = 1, primary = 2
< error
| 1064 42000 …near 'primary = 2' at line 1
> exec ^
: string
| set @a = 1, privilege_checks_user = 2
< ok 0 0
> exec ^
: string
| set @a = 1, privileges = 2
< ok 0 0
> exec ^
: string
| set @a = 1, procedure = 2
< error
| 1064 42000 …near 'procedure = 2' at line 1
> exec ^
: string
| set @a = 1, process = 2
< error
| 1064 42000 …near 'process = 2' at line 1
> exec ^
: string
| set @a = 1, processlist = 2
< ok 0 0
> exec ^
: string
| set @a = 1, profile = 2
< ok 0 0
> exec ^
: string
| set @a = 1, profiles = 2
< ok 0 0
> exec ^
: string
| set @a = 1, proxy = 2
< error
| 1064 42000 …near 'proxy = 2' at line 1
> exec ^
: string
| set @a = 1, purge = 2
< error
| 1064 42000 …near 'purge = 2' at line 1
> exec ^
: string
| set @a = 1, qualify = 2
< ok 0 0
> exec ^
: string
| set @a = 1, quarter = 2
< ok 0 0
> exec ^
: string
| set @a = 1, query = 2
< ok 0 0
> exec ^
: string
| set @a = 1, quick = 2
< ok 0 0
> exec ^
: string
| set @a = 1, random = 2
< ok 0 0
> exec ^
: string
| set @a = 1, range = 2
< error
| 1064 42000 …near 'range = 2' at line 1
> exec ^
: string
| set @a = 1, rank = 2
< ok 0 0
> exec ^
: string
| set @a = 1, read = 2
< error
| 1064 42000 …near 'read = 2' at line 1
> exec ^
: string
| set @a = 1, read_only = 2
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set @a = 1, recover = 2
< ok 0 0
> exec ^
: string
| set @a = 1, recursive = 2
< ok 0 0
> exec ^
: string
| set @a = 1, redo_buffer_size = 2
< ok 0 0
> exec ^
: string
| set @a = 1, redundant = 2
< ok 0 0
> exec ^
: string
| set @a = 1, reference = 2
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set @a = 1, relational = 2
< ok 0 0
> exec ^
: string
| set @a = 1, relay = 2
< ok 0 0
> exec ^
: string
| set @a = 1, relay_log_file = 2
< ok 0 0
> exec ^
: string
| set @a = 1, relay_log_pos = 2
< ok 0 0
> exec ^
: string
| set @a = 1, relay_thread = 2
< ok 0 0
> exec ^
: string
| set @a = 1, relaylog = 2
< ok 0 0
> exec ^
: string
| set @a = 1, release = 2
< error
| 1064 42000 …near 'release = 2' at line 1
> exec ^
: string
| set @a = 1, reload = 2
< error
| 1064 42000 …near 'reload = 2' at line 1
> exec ^
: string
| set @a = 1, remove = 2
< ok 0 0
> exec ^
: string
| set @a = 1, rename = 2
< error
| 1064 42000 …near 'rename = 2' at line 1
> exec ^
: string
| set @a = 1, reorganize = 2
< ok 0 0
> exec ^
: string
| set @a = 1, repair = 2
< ok 0 0
> exec ^
: string
| set @a = 1, repeat = 2
< error
| 1064 42000 …near 'repeat = 2' at line 1
> exec ^
: string
| set @a = 1, repeatable = 2
< ok 0 0
> exec ^
: string
| set @a = 1, replace = 2
< error
| 1064 42000 …near 'replace = 2' at line 1
> exec ^
: string
| set @a = 1, replica = 2
< ok 0 0
> exec ^
: string
| set @a = 1, replica_num = 2
< ok 0 0
> exec ^
: string
| set @a = 1, replicas = 2
< ok 0 0
> exec ^
: string
| set @a = 1, replicate_do_db = 2
< ok 0 0
> exec ^
: string
| set @a = 1, replicate_do_table = 2
< ok 0 0
> exec ^
: string
| set @a = 1, replicate_ignore_db = 2
< ok 0 0
> exec ^
: string
| set @a = 1, replicate_ignore_table = 2
< ok 0 0
> exec ^
: string
| set @a = 1, replicate_rewrite_db = 2
< ok 0 0
> exec ^
: string
| set @a = 1, replicate_wild_do_table = 2
< ok 0 0
> exec ^
: string
| set @a = 1, replicate_wild_ignore_table = 2
< ok 0 0
> exec ^
: string
| set @a = 1, replication = 2
< error
| 1064 42000 …near 'replication = 2' at line 1
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
< ok 0 0
> exec ^
: string
| set @a = 1, reset = 2
< ok 0 0
> exec ^
: string
| set @a = 1, resignal = 2
< error
| 1064 42000 …near 'resignal = 2' at line 1
> exec ^
: string
| set @a = 1, resource = 2
< error
| 1064 42000 …near 'resource = 2' at line 1
> exec ^
: string
| set @a = 1, respect = 2
< ok 0 0
> exec ^
: string
| set @a = 1, restart = 2
< ok 0 0
> exec ^
: string
| set @a = 1, restore = 2
< ok 0 0
> exec ^
: string
| set @a = 1, restrict = 2
< error
| 1064 42000 …near 'restrict = 2' at line 1
> exec ^
: string
| set @a = 1, resume = 2
< ok 0 0
> exec ^
: string
| set @a = 1, retain = 2
< ok 0 0
> exec ^
: string
| set @a = 1, return = 2
< error
| 1064 42000 …near 'return = 2' at line 1
> exec ^
: string
| set @a = 1, returned_sqlstate = 2
< ok 0 0
> exec ^
: string
| set @a = 1, returning = 2
< ok 0 0
> exec ^
: string
| set @a = 1, returns = 2
< ok 0 0
> exec ^
: string
| set @a = 1, reuse = 2
< ok 0 0
> exec ^
: string
| set @a = 1, reverse = 2
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set @a = 1, rollback = 2
< ok 0 0
> exec ^
: string
| set @a = 1, rollup = 2
< ok 0 0
> exec ^
: string
| set @a = 1, rotate = 2
< ok 0 0
> exec ^
: string
| set @a = 1, routine = 2
< ok 0 0
> exec ^
: string
| set @a = 1, row = 2
< ok 0 0
> exec ^
: string
| set @a = 1, row_count = 2
< ok 0 0
> exec ^
: string
| set @a = 1, row_format = 2
< ok 0 0
> exec ^
: string
| set @a = 1, row_number = 2
< ok 0 0
> exec ^
: string
| set @a = 1, rows = 2
< ok 0 0
> exec ^
: string
| set @a = 1, rtree = 2
< ok 0 0
> exec ^
: string
| set @a = 1, s3 = 2
< ok 0 0
> exec ^
: string
| set @a = 1, savepoint = 2
< ok 0 0
> exec ^
: string
| set @a = 1, schedule = 2
< ok 0 0
> exec ^
: string
| set @a = 1, schema = 2
< error
| 1064 42000 …near 'schema = 2' at line 1
> exec ^
: string
| set @a = 1, schema_name = 2
< ok 0 0
> exec ^
: string
| set @a = 1, schemas = 2
< error
| 1064 42000 …near 'schemas = 2' at line 1
> exec ^
: string
| set @a = 1, second = 2
< ok 0 0
> exec ^
: string
| set @a = 1, second_microsecond = 2
< error
| 1064 42000 …near 'second_microsecond = 2' at line 1
> exec ^
: string
| set @a = 1, secondary = 2
< ok 0 0
> exec ^
: string
| set @a = 1, secondary_engine = 2
< ok 0 0
> exec ^
: string
| set @a = 1, secondary_engine_attribute = 2
< ok 0 0
> exec ^
: string
| set @a = 1, secondary_load = 2
< ok 0 0
> exec ^
: string
| set @a = 1, secondary_unload = 2
< ok 0 0
> exec ^
: string
| set @a = 1, security = 2
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set @a = 1, serializable = 2
< ok 0 0
> exec ^
: string
| set @a = 1, server = 2
< ok 0 0
> exec ^
: string
| set @a = 1, session = 2
< ok 0 0
> exec ^
: string
| set @a = 1, session_user = 2
< ok 0 0
> exec ^
: string
| set @a = 1, set = 2
< error
| 1064 42000 …near 'set = 2' at line 1
> exec ^
: string
| set @a = 1, sets = 2
< ok 0 0
> exec ^
: string
| set @a = 1, share = 2
< ok 0 0
> exec ^
: string
| set @a = 1, show = 2
< error
| 1064 42000 …near 'show = 2' at line 1
> exec ^
: string
| set @a = 1, shutdown = 2
< error
| 1064 42000 …near 'shutdown = 2' at line 1
> exec ^
: string
| set @a = 1, signal = 2
< error
| 1064 42000 …near 'signal = 2' at line 1
> exec ^
: string
| set @a = 1, signed = 2
< ok 0 0
> exec ^
: string
| set @a = 1, simple = 2
< ok 0 0
> exec ^
: string
| set @a = 1, skip = 2
< error
| 1064 42000 …near 'skip = 2' at line 1
> exec ^
: string
| set @a = 1, slave = 2
< ok 0 0
> exec ^
: string
| set @a = 1, slow = 2
< ok 0 0
> exec ^
: string
| set @a = 1, smallint = 2
< error
| 1064 42000 …near 'smallint = 2' at line 1
> exec ^
: string
| set @a = 1, snapshot = 2
< ok 0 0
> exec ^
: string
| set @a = 1, socket = 2
< ok 0 0
> exec ^
: string
| set @a = 1, some = 2
< ok 0 0
> exec ^
: string
| set @a = 1, soname = 2
< ok 0 0
> exec ^
: string
| set @a = 1, sounds = 2
< ok 0 0
> exec ^
: string
| set @a = 1, source = 2
< ok 0 0
> exec ^
: string
| set @a = 1, source_auto_position = 2
< ok 0 0
> exec ^
: string
| set @a = 1, source_bind = 2
< ok 0 0
> exec ^
: string
| set @a = 1, source_compression_algorithms = 2
< ok 0 0
> exec ^
: string
| set @a = 1, source_connect_retry = 2
< ok 0 0
> exec ^
: string
| set @a = 1, source_connection_auto_failover = 2
< ok 0 0
> exec ^
: string
| set @a = 1, source_delay = 2
< ok 0 0
> exec ^
: string
| set @a = 1, source_heartbeat_period = 2
< ok 0 0
> exec ^
: string
| set @a = 1, source_host = 2
< ok 0 0
> exec ^
: string
| set @a = 1, source_log_file = 2
< ok 0 0
> exec ^
: string
| set @a = 1, source_log_pos = 2
< ok 0 0
> exec ^
: string
| set @a = 1, source_password = 2
< ok 0 0
> exec ^
: string
| set @a = 1, source_port = 2
< ok 0 0
> exec ^
: string
| set @a = 1, source_public_key_path = 2
< ok 0 0
> exec ^
: string
| set @a = 1, source_retry_count = 2
< ok 0 0
> exec ^
: string
| set @a = 1, source_ssl = 2
< ok 0 0
> exec ^
: string
| set @a = 1, source_ssl_ca = 2
< ok 0 0
> exec ^
: string
| set @a = 1, source_ssl_capath = 2
< ok 0 0
> exec ^
: string
| set @a = 1, source_ssl_cert = 2
< ok 0 0
> exec ^
: string
| set @a = 1, source_ssl_cipher = 2
< ok 0 0
> exec ^
: string
| set @a = 1, source_ssl_crl = 2
< ok 0 0
> exec ^
: string
| set @a = 1, source_ssl_crlpath = 2
< ok 0 0
> exec ^
: string
| set @a = 1, source_ssl_key = 2
< ok 0 0
> exec ^
: string
| set @a = 1, source_ssl_verify_server_cert = 2
< ok 0 0
> exec ^
: string
| set @a = 1, source_tls_ciphersuites = 2
< ok 0 0
> exec ^
: string
| set @a = 1, source_tls_version = 2
< ok 0 0
> exec ^
: string
| set @a = 1, source_user = 2
< ok 0 0
> exec ^
: string
| set @a = 1, source_zstd_compression_level = 2
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set @a = 1, sql_after_mts_gaps = 2
< ok 0 0
> exec ^
: string
| set @a = 1, sql_before_gtids = 2
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set @a = 1, sql_small_result = 2
< error
| 1064 42000 …near 'sql_small_result = 2' at line 1
> exec ^
: string
| set @a = 1, sql_thread = 2
< ok 0 0
> exec ^
: string
| set @a = 1, sql_tsi_day = 2
< ok 0 0
> exec ^
: string
| set @a = 1, sql_tsi_hour = 2
< ok 0 0
> exec ^
: string
| set @a = 1, sql_tsi_minute = 2
< ok 0 0
> exec ^
: string
| set @a = 1, sql_tsi_month = 2
< ok 0 0
> exec ^
: string
| set @a = 1, sql_tsi_quarter = 2
< ok 0 0
> exec ^
: string
| set @a = 1, sql_tsi_second = 2
< ok 0 0
> exec ^
: string
| set @a = 1, sql_tsi_week = 2
< ok 0 0
> exec ^
: string
| set @a = 1, sql_tsi_year = 2
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set @a = 1, sse_kms = 2
< ok 0 0
> exec ^
: string
| set @a = 1, sse_s3 = 2
< ok 0 0
> exec ^
: string
| set @a = 1, ssl = 2
< error
| 1064 42000 …near 'ssl = 2' at line 1
> exec ^
: string
| set @a = 1, st_collect = 2
< ok 0 0
> exec ^
: string
| set @a = 1, stacked = 2
< ok 0 0
> exec ^
: string
| set @a = 1, start = 2
< ok 0 0
> exec ^
: string
| set @a = 1, starting = 2
< error
| 1064 42000 …near 'starting = 2' at line 1
> exec ^
: string
| set @a = 1, starts = 2
< ok 0 0
> exec ^
: string
| set @a = 1, stats_auto_recalc = 2
< ok 0 0
> exec ^
: string
| set @a = 1, stats_persistent = 2
< ok 0 0
> exec ^
: string
| set @a = 1, stats_sample_pages = 2
< ok 0 0
> exec ^
: string
| set @a = 1, status = 2
< ok 0 0
> exec ^
: string
| set @a = 1, std = 2
< ok 0 0
> exec ^
: string
| set @a = 1, stddev = 2
< ok 0 0
> exec ^
: string
| set @a = 1, stddev_pop = 2
< ok 0 0
> exec ^
: string
| set @a = 1, stddev_samp = 2
< ok 0 0
> exec ^
: string
| set @a = 1, stop = 2
< ok 0 0
> exec ^
: string
| set @a = 1, storage = 2
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set @a = 1, strict_load = 2
< ok 0 0
> exec ^
: string
| set @a = 1, string = 2
< ok 0 0
> exec ^
: string
| set @a = 1, subclass_origin = 2
< ok 0 0
> exec ^
: string
| set @a = 1, subdate = 2
< ok 0 0
> exec ^
: string
| set @a = 1, subject = 2
< ok 0 0
> exec ^
: string
| set @a = 1, subpartition = 2
< ok 0 0
> exec ^
: string
| set @a = 1, subpartitions = 2
< ok 0 0
> exec ^
: string
| set @a = 1, substr = 2
< ok 0 0
> exec ^
: string
| set @a = 1, substring = 2
< ok 0 0
> exec ^
: string
| set @a = 1, sum = 2
< ok 0 0
> exec ^
: string
| set @a = 1, super = 2
< error
| 1064 42000 …near 'super = 2' at line 1
> exec ^
: string
| set @a = 1, suspend = 2
< ok 0 0
> exec ^
: string
| set @a = 1, swaps = 2
< ok 0 0
> exec ^
: string
| set @a = 1, switches = 2
< ok 0 0
> exec ^
: string
| set @a = 1, sysdate = 2
< ok 0 0
> exec ^
: string
| set @a = 1, system = 2
< ok 0 0
> exec ^
: string
| set @a = 1, system_user = 2
< ok 0 0
> exec ^
: string
| set @a = 1, table = 2
< error
| 1064 42000 …near 'table = 2' at line 1
> exec ^
: string
| set @a = 1, table_checksum = 2
< ok 0 0
> exec ^
: string
| set @a = 1, table_name = 2
< ok 0 0
> exec ^
: string
| set @a = 1, tables = 2
< ok 0 0
> exec ^
: string
| set @a = 1, tablesample = 2
< ok 0 0
> exec ^
: string
| set @a = 1, tablespace = 2
< ok 0 0
> exec ^
: string
| set @a = 1, tablet_size = 2
< ok 0 0
> exec ^
: string
| set @a = 1, template = 2
< ok 0 0
> exec ^
: string
| set @a = 1, temporary = 2
< ok 0 0
> exec ^
: string
| set @a = 1, temptable = 2
< ok 0 0
> exec ^
: string
| set @a = 1, terminated = 2
< error
| 1064 42000 …near 'terminated = 2' at line 1
> exec ^
: string
| set @a = 1, text = 2
< ok 0 0
> exec ^
: string
| set @a = 1, than = 2
< ok 0 0
> exec ^
: string
| set @a = 1, then = 2
< error
| 1064 42000 …near 'then = 2' at line 1
> exec ^
: string
| set @a = 1, thread_priority = 2
< ok 0 0
> exec ^
: string
| set @a = 1, ties = 2
< ok 0 0
> exec ^
: string
| set @a = 1, time = 2
< ok 0 0
> exec ^
: string
| set @a = 1, timestamp = 2
< ok 0 0
> exec ^
: string
| set @a = 1, timestampadd = 2
< ok 0 0
> exec ^
: string
| set @a = 1, timestampdiff = 2
< ok 0 0
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
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set @a = 1, trigger = 2
< error
| 1064 42000 …near 'trigger = 2' at line 1
> exec ^
: string
| set @a = 1, triggers = 2
< ok 0 0
> exec ^
: string
| set @a = 1, trim = 2
< ok 0 0
> exec ^
: string
| set @a = 1, true = 2
< error
| 1064 42000 …near 'true = 2' at line 1
> exec ^
: string
| set @a = 1, truncate = 2
< ok 0 0
> exec ^
: string
| set @a = 1, type = 2
< ok 0 0
> exec ^
: string
| set @a = 1, types = 2
< ok 0 0
> exec ^
: string
| set @a = 1, unbounded = 2
< ok 0 0
> exec ^
: string
| set @a = 1, uncommitted = 2
< ok 0 0
> exec ^
: string
| set @a = 1, undefined = 2
< ok 0 0
> exec ^
: string
| set @a = 1, undo = 2
< error
| 1064 42000 …near 'undo = 2' at line 1
> exec ^
: string
| set @a = 1, undo_buffer_size = 2
< ok 0 0
> exec ^
: string
| set @a = 1, undofile = 2
< ok 0 0
> exec ^
: string
| set @a = 1, unicode = 2
< ok 0 0
> exec ^
: string
| set @a = 1, uninstall = 2
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set @a = 1, unlock = 2
< error
| 1064 42000 …near 'unlock = 2' at line 1
> exec ^
: string
| set @a = 1, unregister = 2
< ok 0 0
> exec ^
: string
| set @a = 1, unsigned = 2
< error
| 1064 42000 …near 'unsigned = 2' at line 1
> exec ^
: string
| set @a = 1, until = 2
< ok 0 0
> exec ^
: string
| set @a = 1, update = 2
< error
| 1064 42000 …near 'update = 2' at line 1
> exec ^
: string
| set @a = 1, upgrade = 2
< ok 0 0
> exec ^
: string
| set @a = 1, uri = 2
< ok 0 0
> exec ^
: string
| set @a = 1, url = 2
< ok 0 0
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
| set @a = 1, use_bloom_filter = 2
< ok 0 0
> exec ^
: string
| set @a = 1, use_frm = 2
< ok 0 0
> exec ^
: string
| set @a = 1, user = 2
< ok 0 0
> exec ^
: string
| set @a = 1, user_resources = 2
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set @a = 1, validation = 2
< ok 0 0
> exec ^
: string
| set @a = 1, value = 2
< ok 0 0
> exec ^
: string
| set @a = 1, values = 2
< error
| 1064 42000 …near 'values = 2' at line 1
> exec ^
: string
| set @a = 1, var_pop = 2
< ok 0 0
> exec ^
: string
| set @a = 1, var_samp = 2
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set @a = 1, variance = 2
< ok 0 0
> exec ^
: string
| set @a = 1, varying = 2
< error
| 1064 42000 …near 'varying = 2' at line 1
> exec ^
: string
| set @a = 1, vcpu = 2
< ok 0 0
> exec ^
: string
| set @a = 1, vector = 2
< ok 0 0
> exec ^
: string
| set @a = 1, verify_key_constraints = 2
< ok 0 0
> exec ^
: string
| set @a = 1, view = 2
< ok 0 0
> exec ^
: string
| set @a = 1, virtual = 2
< error
| 1064 42000 …near 'virtual = 2' at line 1
> exec ^
: string
| set @a = 1, visible = 2
< ok 0 0
> exec ^
: string
| set @a = 1, wait = 2
< ok 0 0
> exec ^
: string
| set @a = 1, warnings = 2
< ok 0 0
> exec ^
: string
| set @a = 1, week = 2
< ok 0 0
> exec ^
: string
| set @a = 1, weight_string = 2
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set @a = 1, with = 2
< error
| 1064 42000 …near 'with = 2' at line 1
> exec ^
: string
| set @a = 1, without = 2
< ok 0 0
> exec ^
: string
| set @a = 1, work = 2
< ok 0 0
> exec ^
: string
| set @a = 1, wrapper = 2
< ok 0 0
> exec ^
: string
| set @a = 1, write = 2
< error
| 1064 42000 …near 'write = 2' at line 1
> exec ^
: string
| set @a = 1, x509 = 2
< ok 0 0
> exec ^
: string
| set @a = 1, xa = 2
< ok 0 0
> exec ^
: string
| set @a = 1, xid = 2
< ok 0 0
> exec ^
: string
| set @a = 1, xml = 2
< ok 0 0
> exec ^
: string
| set @a = 1, xor = 2
< error
| 1064 42000 …near 'xor = 2' at line 1
> exec ^
: string
| set @a = 1, year = 2
< ok 0 0
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
< ok 0 0
