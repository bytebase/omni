-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestSetVariableNameSweepTiDBOracle default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| set absent = 1
< ok 0 0
> exec ^
: string
| set accessible = 1
< ok 0 0
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
| 1064 42000 …line 1 column 7 near "add = 1" 
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
| 1064 42000 …line 1 column 7 near "all = 1" 
> exec ^
: string
| set allow_missing_files = 1
< ok 0 0
> exec ^
: string
| set alter = 1
< error
| 1064 42000 …line 1 column 9 near "alter = 1" 
> exec ^
: string
| set always = 1
< ok 0 0
> exec ^
: string
| set analyze = 1
< error
| 1064 42000 …line 1 column 11 near "analyze = 1" 
> exec ^
: string
| set and = 1
< error
| 1064 42000 …line 1 column 7 near "and = 1" 
> exec ^
: string
| set any = 1
< ok 0 0
> exec ^
: string
| set array = 1
< error
| 1064 42000 …line 1 column 9 near "array = 1" 
> exec ^
: string
| set as = 1
< error
| 1064 42000 …line 1 column 6 near "as = 1" 
> exec ^
: string
| set asc = 1
< error
| 1064 42000 …line 1 column 7 near "asc = 1" 
> exec ^
: string
| set ascii = 1
< ok 0 0
> exec ^
: string
| set asensitive = 1
< ok 0 0
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
| set auto_id_cache = 1
< ok 0 0
> exec ^
: string
| set auto_increment = 1
< ok 0 0
> exec ^
: string
| set auto_random = 1
< ok 0 0
> exec ^
: string
| set auto_random_base = 1
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
| set background = 1
< ok 0 0
> exec ^
: string
| set backup = 1
< ok 0 0
> exec ^
: string
| set backups = 1
< ok 0 0
> exec ^
: string
| set batch = 1
< ok 0 0
> exec ^
: string
| set before = 1
< ok 0 0
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
| 1064 42000 …line 1 column 11 near "between = 1" 
> exec ^
: string
| set bigint = 1
< error
| 1064 42000 …line 1 column 10 near "bigint = 1" 
> exec ^
: string
| set binary = 1
< error
| 1064 42000 …line 1 column 10 near "binary = 1" 
> exec ^
: string
| set binding = 1
< ok 0 0
> exec ^
: string
| set binding_cache = 1
< ok 0 0
> exec ^
: string
| set bindings = 1
< ok 0 0
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
| 1064 42000 …line 1 column 8 near "blob = 1" 
> exec ^
: string
| set block = 1
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
| 1064 42000 …line 1 column 8 near "both = 1" 
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
| set builtins = 1
< ok 0 0
> exec ^
: string
| set bulk = 1
< ok 0 0
> exec ^
: string
| set burstable = 1
< ok 0 0
> exec ^
: string
| set by = 1
< error
| 1064 42000 …line 1 column 6 near "by = 1" 
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
| set calibrate = 1
< ok 0 0
> exec ^
: string
| set call = 1
< error
| 1064 42000 …line 1 column 8 near "call = 1" 
> exec ^
: string
| set cascade = 1
< error
| 1064 42000 …line 1 column 11 near "cascade = 1" 
> exec ^
: string
| set cascaded = 1
< ok 0 0
> exec ^
: string
| set case = 1
< error
| 1064 42000 …line 1 column 8 near "case = 1" 
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
| 1064 42000 …line 1 column 10 near "change = 1" 
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
| 1064 42000 …line 1 column 10 near "= 1" 
> exec ^
: string
| set character = 1
< error
| 1064 42000 …line 1 column 15 near "= 1" 
> exec ^
: string
| set charset = 1
< ok 0 0
> exec ^
: string
| set check = 1
< error
| 1064 42000 …line 1 column 9 near "check = 1" 
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
| set clustered = 1
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
| 1064 42000 …line 1 column 11 near "collate = 1" 
> exec ^
: string
| set collation = 1
< ok 0 0
> exec ^
: string
| set column = 1
< error
| 1064 42000 …line 1 column 10 near "column = 1" 
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
| set column_stats_usage = 1
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
< ok 0 0
> exec ^
: string
| set config = 1
< ok 0 0
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
| 1064 42000 …line 1 column 14 near "constraint = 1" 
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
| set constraints = 1
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
| 1064 42000 …line 1 column 12 near "continue = 1" 
> exec ^
: string
| set convert = 1
< error
| 1064 42000 …line 1 column 11 near "convert = 1" 
> exec ^
: string
| set cooldown = 1
< ok 0 0
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
| 1064 42000 …line 1 column 10 near "create = 1" 
> exec ^
: string
| set cross = 1
< error
| 1064 42000 …line 1 column 9 near "cross = 1" 
> exec ^
: string
| set cube = 1
< ok 0 0
> exec ^
: string
| set cume_dist = 1
< error
| 1064 42000 …line 1 column 13 near "cume_dist = 1" 
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
| 1064 42000 …line 1 column 16 near "current_date = 1" 
> exec ^
: string
| set current_time = 1
< error
| 1064 42000 …line 1 column 16 near "current_time = 1" 
> exec ^
: string
| set current_timestamp = 1
< error
| 1064 42000 …line 1 column 21 near "current_timestamp = 1" 
> exec ^
: string
| set current_user = 1
< error
| 1064 42000 …line 1 column 16 near "current_user = 1" 
> exec ^
: string
| set cursor = 1
< error
| 1064 42000 …line 1 column 10 near "cursor = 1" 
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
| set cycle = 1
< ok 0 0
> exec ^
: string
| set data = 1
< ok 0 0
> exec ^
: string
| set database = 1
< error
| 1064 42000 …line 1 column 12 near "database = 1" 
> exec ^
: string
| set databases = 1
< error
| 1064 42000 …line 1 column 13 near "databases = 1" 
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
| 1064 42000 …line 1 column 12 near "day_hour = 1" 
> exec ^
: string
| set day_microsecond = 1
< error
| 1064 42000 …line 1 column 19 near "day_microsecond = 1" 
> exec ^
: string
| set day_minute = 1
< error
| 1064 42000 …line 1 column 14 near "day_minute = 1" 
> exec ^
: string
| set day_second = 1
< error
| 1064 42000 …line 1 column 14 near "day_second = 1" 
> exec ^
: string
| set deallocate = 1
< ok 0 0
> exec ^
: string
| set dec = 1
< error
| 1064 42000 …line 1 column 7 near "dec = 1" 
> exec ^
: string
| set decimal = 1
< error
| 1064 42000 …line 1 column 11 near "decimal = 1" 
> exec ^
: string
| set declare = 1
< ok 0 0
> exec ^
: string
| set default = 1
< error
| 1064 42000 …line 1 column 13 near "= 1" 
> exec ^
: string
| set default_auth = 1
< ok 0 0
> exec ^
: string
| set defined = 1
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
| 1064 42000 …line 1 column 11 near "delayed = 1" 
> exec ^
: string
| set delete = 1
< error
| 1064 42000 …line 1 column 10 near "delete = 1" 
> exec ^
: string
| set dense_rank = 1
< error
| 1064 42000 …line 1 column 14 near "dense_rank = 1" 
> exec ^
: string
| set desc = 1
< error
| 1064 42000 …line 1 column 8 near "desc = 1" 
> exec ^
: string
| set describe = 1
< error
| 1064 42000 …line 1 column 12 near "describe = 1" 
> exec ^
: string
| set description = 1
< ok 0 0
> exec ^
: string
| set deterministic = 1
< ok 0 0
> exec ^
: string
| set diagnostics = 1
< ok 0 0
> exec ^
: string
| set digest = 1
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
| set disabled = 1
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
| 1064 42000 …line 1 column 12 near "distinct = 1" 
> exec ^
: string
| set distinctrow = 1
< error
| 1064 42000 …line 1 column 15 near "distinctrow = 1" 
> exec ^
: string
| set distribute = 1
< ok 0 0
> exec ^
: string
| set distribution = 1
< ok 0 0
> exec ^
: string
| set distributions = 1
< ok 0 0
> exec ^
: string
| set div = 1
< error
| 1064 42000 …line 1 column 7 near "div = 1" 
> exec ^
: string
| set do = 1
< ok 0 0
> exec ^
: string
| set double = 1
< error
| 1064 42000 …line 1 column 10 near "double = 1" 
> exec ^
: string
| set drop = 1
< error
| 1064 42000 …line 1 column 8 near "drop = 1" 
> exec ^
: string
| set dry = 1
< ok 0 0
> exec ^
: string
| set dryrun = 1
< ok 0 0
> exec ^
: string
| set dual = 1
< error
| 1064 42000 …line 1 column 8 near "dual = 1" 
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
| set duration = 1
< ok 0 0
> exec ^
: string
| set dynamic = 1
< ok 0 0
> exec ^
: string
| set each = 1
< ok 0 0
> exec ^
: string
| set else = 1
< error
| 1064 42000 …line 1 column 8 near "else = 1" 
> exec ^
: string
| set elseif = 1
< error
| 1064 42000 …line 1 column 10 near "elseif = 1" 
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
| set enabled = 1
< ok 0 0
> exec ^
: string
| set enclosed = 1
< error
| 1064 42000 …line 1 column 12 near "enclosed = 1" 
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
| set end_time = 1
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
| 1064 42000 …line 1 column 11 near "escaped = 1" 
> exec ^
: string
| set event = 1
< ok 0 0
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
| set exact = 1
< ok 0 0
> exec ^
: string
| set except = 1
< error
| 1064 42000 …line 1 column 10 near "except = 1" 
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
| set exec_elapsed = 1
< ok 0 0
> exec ^
: string
| set execute = 1
< ok 0 0
> exec ^
: string
| set exists = 1
< error
| 1064 42000 …line 1 column 10 near "exists = 1" 
> exec ^
: string
| set exit = 1
< error
| 1064 42000 …line 1 column 8 near "exit = 1" 
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
| 1064 42000 …line 1 column 11 near "explain = 1" 
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
| 1064 42000 …line 1 column 9 near "false = 1" 
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
| 1064 42000 …line 1 column 9 near "fetch = 1" 
> exec ^
: string
| set fields = 1
< ok 0 0
> exec ^
: string
| set file = 1
< ok 0 0
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
< error
| 1064 42000 …line 1 column 15 near "first_value = 1" 
> exec ^
: string
| set fixed = 1
< ok 0 0
> exec ^
: string
| set flashback = 1
< ok 0 0
> exec ^
: string
| set float = 1
< error
| 1064 42000 …line 1 column 9 near "float = 1" 
> exec ^
: string
| set float4 = 1
< error
| 1064 42000 …line 1 column 10 near "float4 = 1" 
> exec ^
: string
| set float8 = 1
< error
| 1064 42000 …line 1 column 10 near "float8 = 1" 
> exec ^
: string
| set flush = 1
< ok 0 0
> exec ^
: string
| set follower_constraints = 1
< ok 0 0
> exec ^
: string
| set followers = 1
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
| 1064 42000 …line 1 column 7 near "for = 1" 
> exec ^
: string
| set force = 1
< error
| 1064 42000 …line 1 column 9 near "force = 1" 
> exec ^
: string
| set foreign = 1
< error
| 1064 42000 …line 1 column 11 near "foreign = 1" 
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
| 1064 42000 …line 1 column 8 near "from = 1" 
> exec ^
: string
| set full = 1
< ok 0 0
> exec ^
: string
| set fulltext = 1
< error
| 1064 42000 …line 1 column 12 near "fulltext = 1" 
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
| 1064 42000 …line 1 column 13 near "generated = 1" 
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
< ok 0 0
> exec ^
: string
| set get_format = 1
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
| 1064 42000 …line 1 column 9 near "grant = 1" 
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
| 1064 42000 …line 1 column 9 near "group = 1" 
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
< error
| 1064 42000 …line 1 column 10 near "groups = 1" 
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
| 1064 42000 …line 1 column 10 near "having = 1" 
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
| set high = 1
< ok 0 0
> exec ^
: string
| set high_priority = 1
< error
| 1064 42000 …line 1 column 17 near "high_priority = 1" 
> exec ^
: string
| set histogram = 1
< ok 0 0
> exec ^
: string
| set histograms_in_flight = 1
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
| 1064 42000 …line 1 column 20 near "hour_microsecond = 1" 
> exec ^
: string
| set hour_minute = 1
< error
| 1064 42000 …line 1 column 15 near "hour_minute = 1" 
> exec ^
: string
| set hour_second = 1
< error
| 1064 42000 …line 1 column 15 near "hour_second = 1" 
> exec ^
: string
| set identified = 1
< ok 0 0
> exec ^
: string
| set if = 1
< error
| 1064 42000 …line 1 column 6 near "if = 1" 
> exec ^
: string
| set ignore = 1
< error
| 1064 42000 …line 1 column 10 near "ignore = 1" 
> exec ^
: string
| set ignore_server_ids = 1
< ok 0 0
> exec ^
: string
| set ilike = 1
< error
| 1064 42000 …line 1 column 9 near "ilike = 1" 
> exec ^
: string
| set import = 1
< ok 0 0
> exec ^
: string
| set in = 1
< error
| 1064 42000 …line 1 column 6 near "in = 1" 
> exec ^
: string
| set inactive = 1
< ok 0 0
> exec ^
: string
| set increment = 1
< ok 0 0
> exec ^
: string
| set index = 1
< error
| 1064 42000 …line 1 column 9 near "index = 1" 
> exec ^
: string
| set indexes = 1
< ok 0 0
> exec ^
: string
| set infile = 1
< error
| 1064 42000 …line 1 column 10 near "infile = 1" 
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
| 1064 42000 …line 1 column 9 near "inner = 1" 
> exec ^
: string
| set inout = 1
< error
| 1064 42000 …line 1 column 9 near "inout = 1" 
> exec ^
: string
| set insensitive = 1
< ok 0 0
> exec ^
: string
| set insert = 1
< error
| 1064 42000 …line 1 column 10 near "insert = 1" 
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
| 1064 42000 …line 1 column 7 near "int = 1" 
> exec ^
: string
| set int1 = 1
< error
| 1064 42000 …line 1 column 8 near "int1 = 1" 
> exec ^
: string
| set int2 = 1
< error
| 1064 42000 …line 1 column 8 near "int2 = 1" 
> exec ^
: string
| set int3 = 1
< error
| 1064 42000 …line 1 column 8 near "int3 = 1" 
> exec ^
: string
| set int4 = 1
< error
| 1064 42000 …line 1 column 8 near "int4 = 1" 
> exec ^
: string
| set int8 = 1
< error
| 1064 42000 …line 1 column 8 near "int8 = 1" 
> exec ^
: string
| set integer = 1
< error
| 1064 42000 …line 1 column 11 near "integer = 1" 
> exec ^
: string
| set intersect = 1
< error
| 1064 42000 …line 1 column 13 near "intersect = 1" 
> exec ^
: string
| set interval = 1
< error
| 1064 42000 …line 1 column 12 near "interval = 1" 
> exec ^
: string
| set into = 1
< error
| 1064 42000 …line 1 column 8 near "into = 1" 
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
< ok 0 0
> exec ^
: string
| set io_before_gtids = 1
< ok 0 0
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
| 1064 42000 …line 1 column 6 near "is = 1" 
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
| 1064 42000 …line 1 column 11 near "iterate = 1" 
> exec ^
: string
| set job = 1
< ok 0 0
> exec ^
: string
| set jobs = 1
< ok 0 0
> exec ^
: string
| set join = 1
< error
| 1064 42000 …line 1 column 8 near "join = 1" 
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
| 1064 42000 …line 1 column 7 near "key = 1" 
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
| 1064 42000 …line 1 column 8 near "keys = 1" 
> exec ^
: string
| set kill = 1
< error
| 1064 42000 …line 1 column 8 near "kill = 1" 
> exec ^
: string
| set labels = 1
< ok 0 0
> exec ^
: string
| set lag = 1
< error
| 1064 42000 …line 1 column 7 near "lag = 1" 
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
< error
| 1064 42000 …line 1 column 14 near "last_value = 1" 
> exec ^
: string
| set lateral = 1
< ok 0 0
> exec ^
: string
| set lead = 1
< error
| 1064 42000 …line 1 column 8 near "lead = 1" 
> exec ^
: string
| set leader_constraints = 1
< ok 0 0
> exec ^
: string
| set leading = 1
< error
| 1064 42000 …line 1 column 11 near "leading = 1" 
> exec ^
: string
| set learner_constraints = 1
< ok 0 0
> exec ^
: string
| set learners = 1
< ok 0 0
> exec ^
: string
| set least = 1
< ok 0 0
> exec ^
: string
| set leave = 1
< error
| 1064 42000 …line 1 column 9 near "leave = 1" 
> exec ^
: string
| set leaves = 1
< ok 0 0
> exec ^
: string
| set left = 1
< error
| 1064 42000 …line 1 column 8 near "left = 1" 
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
| 1064 42000 …line 1 column 8 near "like = 1" 
> exec ^
: string
| set limit = 1
< error
| 1064 42000 …line 1 column 9 near "limit = 1" 
> exec ^
: string
| set linear = 1
< error
| 1064 42000 …line 1 column 10 near "linear = 1" 
> exec ^
: string
| set lines = 1
< error
| 1064 42000 …line 1 column 9 near "lines = 1" 
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
| 1064 42000 …line 1 column 8 near "load = 1" 
> exec ^
: string
| set local = 1
< ok 0 0
> exec ^
: string
| set localtime = 1
< error
| 1064 42000 …line 1 column 13 near "localtime = 1" 
> exec ^
: string
| set localtimestamp = 1
< error
| 1064 42000 …line 1 column 18 near "localtimestamp = 1" 
> exec ^
: string
| set location = 1
< ok 0 0
> exec ^
: string
| set lock = 1
< error
| 1064 42000 …line 1 column 8 near "lock = 1" 
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
| 1064 42000 …line 1 column 8 near "long = 1" 
> exec ^
: string
| set longblob = 1
< error
| 1064 42000 …line 1 column 12 near "longblob = 1" 
> exec ^
: string
| set longtext = 1
< error
| 1064 42000 …line 1 column 12 near "longtext = 1" 
> exec ^
: string
| set loop = 1
< ok 0 0
> exec ^
: string
| set low = 1
< ok 0 0
> exec ^
: string
| set low_priority = 1
< error
| 1064 42000 …line 1 column 16 near "low_priority = 1" 
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
| set master_log_file = 1
< ok 0 0
> exec ^
: string
| set match = 1
< error
| 1064 42000 …line 1 column 9 near "match = 1" 
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
| 1064 42000 …line 1 column 12 near "maxvalue = 1" 
> exec ^
: string
| set medium = 1
< ok 0 0
> exec ^
: string
| set mediumblob = 1
< error
| 1064 42000 …line 1 column 14 near "mediumblob = 1" 
> exec ^
: string
| set mediumint = 1
< error
| 1064 42000 …line 1 column 13 near "mediumint = 1" 
> exec ^
: string
| set mediumtext = 1
< error
| 1064 42000 …line 1 column 14 near "mediumtext = 1" 
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
| 1064 42000 …line 1 column 13 near "middleint = 1" 
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
| 1064 42000 …line 1 column 22 near "minute_microsecond = 1" 
> exec ^
: string
| set minute_second = 1
< error
| 1064 42000 …line 1 column 17 near "minute_second = 1" 
> exec ^
: string
| set minvalue = 1
< ok 0 0
> exec ^
: string
| set mod = 1
< error
| 1064 42000 …line 1 column 7 near "mod = 1" 
> exec ^
: string
| set mode = 1
< ok 0 0
> exec ^
: string
| set modifies = 1
< ok 0 0
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
| 1064 42000 …line 1 column 11 near "natural = 1" 
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
| set next_row_id = 1
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
| 1064 42000 …line 1 column 22 near "no_write_to_binlog = 1" 
> exec ^
: string
| set nocache = 1
< ok 0 0
> exec ^
: string
| set nocycle = 1
< ok 0 0
> exec ^
: string
| set nodegroup = 1
< ok 0 0
> exec ^
: string
| set nomaxvalue = 1
< ok 0 0
> exec ^
: string
| set nominvalue = 1
< ok 0 0
> exec ^
: string
| set nonclustered = 1
< ok 0 0
> exec ^
: string
| set none = 1
< ok 0 0
> exec ^
: string
| set not = 1
< error
| 1064 42000 …line 1 column 7 near "not = 1" 
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
< error
| 1064 42000 …line 1 column 13 near "nth_value = 1" 
> exec ^
: string
| set ntile = 1
< error
| 1064 42000 …line 1 column 9 near "ntile = 1" 
> exec ^
: string
| set null = 1
< error
| 1064 42000 …line 1 column 8 near "null = 1" 
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
| 1064 42000 …line 1 column 11 near "numeric = 1" 
> exec ^
: string
| set nvarchar = 1
< ok 0 0
> exec ^
: string
| set of = 1
< error
| 1064 42000 …line 1 column 6 near "of = 1" 
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
| set oltp_read_only = 1
< ok 0 0
> exec ^
: string
| set oltp_read_write = 1
< ok 0 0
> exec ^
: string
| set oltp_write_only = 1
< ok 0 0
> exec ^
: string
| set on = 1
< error
| 1064 42000 …line 1 column 6 near "on = 1" 
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
| 1064 42000 …line 1 column 12 near "optimize = 1" 
> exec ^
: string
| set optimizer_costs = 1
< ok 0 0
> exec ^
: string
| set option = 1
< error
| 1064 42000 …line 1 column 10 near "option = 1" 
> exec ^
: string
| set optional = 1
< ok 0 0
> exec ^
: string
| set optionally = 1
< error
| 1064 42000 …line 1 column 14 near "optionally = 1" 
> exec ^
: string
| set options = 1
< ok 0 0
> exec ^
: string
| set or = 1
< error
| 1064 42000 …line 1 column 6 near "or = 1" 
> exec ^
: string
| set order = 1
< error
| 1064 42000 …line 1 column 9 near "order = 1" 
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
| 1064 42000 …line 1 column 7 near "out = 1" 
> exec ^
: string
| set outer = 1
< error
| 1064 42000 …line 1 column 9 near "outer = 1" 
> exec ^
: string
| set outfile = 1
< error
| 1064 42000 …line 1 column 11 near "outfile = 1" 
> exec ^
: string
| set over = 1
< error
| 1064 42000 …line 1 column 8 near "over = 1" 
> exec ^
: string
| set overlay = 1
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
| 1064 42000 …line 1 column 13 near "partition = 1" 
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
| 1064 42000 …line 1 column 16 near "1" 
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
| set percent_rank = 1
< error
| 1064 42000 …line 1 column 16 near "percent_rank = 1" 
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
| set placement = 1
< ok 0 0
> exec ^
: string
| set placing = 1
< ok 0 0
> exec ^
: string
| set plan = 1
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
| set policy = 1
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
| set pre_split_regions = 1
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
| 1064 42000 …line 1 column 13 near "precision = 1" 
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
| 1064 42000 …line 1 column 11 near "primary = 1" 
> exec ^
: string
| set primary_region = 1
< ok 0 0
> exec ^
: string
| set priority = 1
< ok 0 0
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
| 1064 42000 …line 1 column 13 near "procedure = 1" 
> exec ^
: string
| set process = 1
< ok 0 0
> exec ^
: string
| set processed_keys = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set purge = 1
< ok 0 0
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
| set query_limit = 1
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
| 1064 42000 …line 1 column 9 near "range = 1" 
> exec ^
: string
| set rank = 1
< error
| 1064 42000 …line 1 column 8 near "rank = 1" 
> exec ^
: string
| set read = 1
< error
| 1064 42000 …line 1 column 8 near "read = 1" 
> exec ^
: string
| set read_only = 1
< ok 0 0
> exec ^
: string
| set read_write = 1
< ok 0 0
> exec ^
: string
| set reads = 1
< ok 0 0
> exec ^
: string
| set real = 1
< error
| 1064 42000 …line 1 column 8 near "real = 1" 
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
< error
| 1064 42000 …line 1 column 13 near "recursive = 1" 
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
| 1064 42000 …line 1 column 14 near "references = 1" 
> exec ^
: string
| set regexp = 1
< error
| 1064 42000 …line 1 column 10 near "regexp = 1" 
> exec ^
: string
| set region = 1
< ok 0 0
> exec ^
: string
| set regions = 1
< ok 0 0
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
| 1064 42000 …line 1 column 11 near "release = 1" 
> exec ^
: string
| set reload = 1
< ok 0 0
> exec ^
: string
| set remove = 1
< ok 0 0
> exec ^
: string
| set rename = 1
< error
| 1064 42000 …line 1 column 10 near "rename = 1" 
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
| 1064 42000 …line 1 column 10 near "repeat = 1" 
> exec ^
: string
| set repeatable = 1
< ok 0 0
> exec ^
: string
| set replace = 1
< error
| 1064 42000 …line 1 column 11 near "replace = 1" 
> exec ^
: string
| set replica = 1
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
< ok 0 0
> exec ^
: string
| set require = 1
< error
| 1064 42000 …line 1 column 11 near "require = 1" 
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
< ok 0 0
> exec ^
: string
| set resource = 1
< ok 0 0
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
| set restores = 1
< ok 0 0
> exec ^
: string
| set restrict = 1
< error
| 1064 42000 …line 1 column 12 near "restrict = 1" 
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
< ok 0 0
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
| 1064 42000 …line 1 column 10 near "revoke = 1" 
> exec ^
: string
| set right = 1
< error
| 1064 42000 …line 1 column 9 near "right = 1" 
> exec ^
: string
| set rlike = 1
< error
| 1064 42000 …line 1 column 9 near "rlike = 1" 
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
< error
| 1064 42000 …line 1 column 7 near "row = 1" 
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
< error
| 1064 42000 …line 1 column 14 near "row_number = 1" 
> exec ^
: string
| set rows = 1
< error
| 1064 42000 …line 1 column 8 near "rows = 1" 
> exec ^
: string
| set rtree = 1
< ok 0 0
> exec ^
: string
| set ru = 1
< ok 0 0
> exec ^
: string
| set ru_per_sec = 1
< ok 0 0
> exec ^
: string
| set rule = 1
< ok 0 0
> exec ^
: string
| set run = 1
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
| 1064 42000 …line 1 column 10 near "schema = 1" 
> exec ^
: string
| set schema_name = 1
< ok 0 0
> exec ^
: string
| set schemas = 1
< error
| 1064 42000 …line 1 column 11 near "schemas = 1" 
> exec ^
: string
| set second = 1
< ok 0 0
> exec ^
: string
| set second_microsecond = 1
< error
| 1064 42000 …line 1 column 22 near "second_microsecond = 1" 
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
| 1064 42000 …line 1 column 10 near "select = 1" 
> exec ^
: string
| set sensitive = 1
< ok 0 0
> exec ^
: string
| set separator = 1
< ok 0 0
> exec ^
: string
| set sequence = 1
< ok 0 0
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
| set session_states = 1
< ok 0 0
> exec ^
: string
| set session_user = 1
< ok 0 0
> exec ^
: string
| set set = 1
< error
| 1064 42000 …line 1 column 7 near "set = 1" 
> exec ^
: string
| set sets = 1
< ok 0 0
> exec ^
: string
| set shard_row_id_bits = 1
< ok 0 0
> exec ^
: string
| set share = 1
< ok 0 0
> exec ^
: string
| set show = 1
< error
| 1064 42000 …line 1 column 8 near "show = 1" 
> exec ^
: string
| set shutdown = 1
< ok 0 0
> exec ^
: string
| set signal = 1
< ok 0 0
> exec ^
: string
| set signed = 1
< ok 0 0
> exec ^
: string
| set similar = 1
< ok 0 0
> exec ^
: string
| set simple = 1
< ok 0 0
> exec ^
: string
| set skip = 1
< ok 0 0
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
| 1064 42000 …line 1 column 12 near "smallint = 1" 
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
| 1064 42000 …line 1 column 11 near "spatial = 1" 
> exec ^
: string
| set specific = 1
< ok 0 0
> exec ^
: string
| set split = 1
< ok 0 0
> exec ^
: string
| set sql = 1
< error
| 1064 42000 …line 1 column 7 near "sql = 1" 
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
| 1064 42000 …line 1 column 18 near "sql_big_result = 1" 
> exec ^
: string
| set sql_buffer_result = 1
< ok 0 0
> exec ^
: string
| set sql_calc_found_rows = 1
< error
| 1064 42000 …line 1 column 23 near "sql_calc_found_rows = 1" 
> exec ^
: string
| set sql_no_cache = 1
< ok 0 0
> exec ^
: string
| set sql_small_result = 1
< error
| 1064 42000 …line 1 column 20 near "sql_small_result = 1" 
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
| 1064 42000 …line 1 column 16 near "sqlexception = 1" 
> exec ^
: string
| set sqlstate = 1
< error
| 1064 42000 …line 1 column 12 near "sqlstate = 1" 
> exec ^
: string
| set sqlwarning = 1
< error
| 1064 42000 …line 1 column 14 near "sqlwarning = 1" 
> exec ^
: string
| set srid = 1
< ok 0 0
> exec ^
: string
| set ssl = 1
< error
| 1064 42000 …line 1 column 7 near "ssl = 1" 
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
| set start_time = 1
< ok 0 0
> exec ^
: string
| set starting = 1
< error
| 1064 42000 …line 1 column 12 near "starting = 1" 
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
| set stats_buckets = 1
< ok 0 0
> exec ^
: string
| set stats_col_choice = 1
< ok 0 0
> exec ^
: string
| set stats_col_list = 1
< ok 0 0
> exec ^
: string
| set stats_extended = 1
< ok 0 0
> exec ^
: string
| set stats_healthy = 1
< ok 0 0
> exec ^
: string
| set stats_histograms = 1
< ok 0 0
> exec ^
: string
| set stats_locked = 1
< ok 0 0
> exec ^
: string
| set stats_meta = 1
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
| set stats_sample_rate = 1
< ok 0 0
> exec ^
: string
| set stats_topn = 1
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
| 1064 42000 …line 1 column 10 near "stored = 1" 
> exec ^
: string
| set straight_join = 1
< error
| 1064 42000 …line 1 column 17 near "straight_join = 1" 
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
< ok 0 0
> exec ^
: string
| set survival_preferences = 1
< ok 0 0
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
| set switch_group = 1
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
| 1064 42000 …line 1 column 9 near "table = 1" 
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
< error
| 1064 42000 …line 1 column 15 near "tablesample = 1" 
> exec ^
: string
| set tablespace = 1
< ok 0 0
> exec ^
: string
| set task_types = 1
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
| 1064 42000 …line 1 column 14 near "terminated = 1" 
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
| 1064 42000 …line 1 column 8 near "then = 1" 
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
| set tiflash = 1
< ok 0 0
> exec ^
: string
| set time = 1
< ok 0 0
> exec ^
: string
| set timeout = 1
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
| 1064 42000 …line 1 column 12 near "tinyblob = 1" 
> exec ^
: string
| set tinyint = 1
< error
| 1064 42000 …line 1 column 11 near "tinyint = 1" 
> exec ^
: string
| set tinytext = 1
< error
| 1064 42000 …line 1 column 12 near "tinytext = 1" 
> exec ^
: string
| set tls = 1
< ok 0 0
> exec ^
: string
| set to = 1
< error
| 1064 42000 …line 1 column 6 near "to = 1" 
> exec ^
: string
| set tpcc = 1
< ok 0 0
> exec ^
: string
| set tpch_10 = 1
< ok 0 0
> exec ^
: string
| set trailing = 1
< error
| 1064 42000 …line 1 column 12 near "trailing = 1" 
> exec ^
: string
| set transaction = 1
< ok 0 0
> exec ^
: string
| set trigger = 1
< error
| 1064 42000 …line 1 column 11 near "trigger = 1" 
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
| 1064 42000 …line 1 column 8 near "true = 1" 
> exec ^
: string
| set truncate = 1
< ok 0 0
> exec ^
: string
| set ttl = 1
< ok 0 0
> exec ^
: string
| set ttl_enable = 1
< ok 0 0
> exec ^
: string
| set ttl_job_interval = 1
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
< ok 0 0
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
| 1064 42000 …line 1 column 9 near "union = 1" 
> exec ^
: string
| set unique = 1
< error
| 1064 42000 …line 1 column 10 near "unique = 1" 
> exec ^
: string
| set unknown = 1
< ok 0 0
> exec ^
: string
| set unlimited = 1
< ok 0 0
> exec ^
: string
| set unlock = 1
< error
| 1064 42000 …line 1 column 10 near "unlock = 1" 
> exec ^
: string
| set unregister = 1
< ok 0 0
> exec ^
: string
| set unsigned = 1
< error
| 1064 42000 …line 1 column 12 near "unsigned = 1" 
> exec ^
: string
| set until = 1
< error
| 1064 42000 …line 1 column 9 near "until = 1" 
> exec ^
: string
| set update = 1
< error
| 1064 42000 …line 1 column 10 near "update = 1" 
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
| 1064 42000 …line 1 column 9 near "usage = 1" 
> exec ^
: string
| set use = 1
< error
| 1064 42000 …line 1 column 7 near "use = 1" 
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
| 1064 42000 …line 1 column 9 near "using = 1" 
> exec ^
: string
| set utc_date = 1
< error
| 1064 42000 …line 1 column 12 near "utc_date = 1" 
> exec ^
: string
| set utc_time = 1
< error
| 1064 42000 …line 1 column 12 near "utc_time = 1" 
> exec ^
: string
| set utc_timestamp = 1
< error
| 1064 42000 …line 1 column 17 near "utc_timestamp = 1" 
> exec ^
: string
| set utilization_limit = 1
< ok 0 0
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
| 1064 42000 …line 1 column 10 near "values = 1" 
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
| 1064 42000 …line 1 column 13 near "varbinary = 1" 
> exec ^
: string
| set varchar = 1
< error
| 1064 42000 …line 1 column 11 near "varchar = 1" 
> exec ^
: string
| set varcharacter = 1
< error
| 1064 42000 …line 1 column 16 near "varcharacter = 1" 
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
| 1064 42000 …line 1 column 11 near "varying = 1" 
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
| 1064 42000 …line 1 column 11 near "virtual = 1" 
> exec ^
: string
| set visible = 1
< ok 0 0
> exec ^
: string
| set voter_constraints = 1
< ok 0 0
> exec ^
: string
| set voters = 1
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
| set watch = 1
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
| 1064 42000 …line 1 column 8 near "when = 1" 
> exec ^
: string
| set where = 1
< error
| 1064 42000 …line 1 column 9 near "where = 1" 
> exec ^
: string
| set while = 1
< error
| 1064 42000 …line 1 column 9 near "while = 1" 
> exec ^
: string
| set window = 1
< error
| 1064 42000 …line 1 column 10 near "window = 1" 
> exec ^
: string
| set with = 1
< error
| 1064 42000 …line 1 column 8 near "with = 1" 
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
| set workload = 1
< ok 0 0
> exec ^
: string
| set wrapper = 1
< ok 0 0
> exec ^
: string
| set write = 1
< error
| 1064 42000 …line 1 column 9 near "write = 1" 
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
| 1064 42000 …line 1 column 7 near "xor = 1" 
> exec ^
: string
| set year = 1
< ok 0 0
> exec ^
: string
| set year_month = 1
< error
| 1064 42000 …line 1 column 14 near "year_month = 1" 
> exec ^
: string
| set zerofill = 1
< error
| 1064 42000 …line 1 column 12 near "zerofill = 1" 
> exec ^
: string
| set zone = 1
< ok 0 0
> exec ^
: string
| set absent := 1
< ok 0 0
> exec ^
: string
| set accessible := 1
< ok 0 0
> exec ^
: string
| set account := 1
< ok 0 0
> exec ^
: string
| set action := 1
< ok 0 0
> exec ^
: string
| set active := 1
< ok 0 0
> exec ^
: string
| set add := 1
< error
| 1064 42000 …line 1 column 7 near "add := 1" 
> exec ^
: string
| set adddate := 1
< ok 0 0
> exec ^
: string
| set admin := 1
< ok 0 0
> exec ^
: string
| set after := 1
< ok 0 0
> exec ^
: string
| set against := 1
< ok 0 0
> exec ^
: string
| set aggregate := 1
< ok 0 0
> exec ^
: string
| set algorithm := 1
< ok 0 0
> exec ^
: string
| set all := 1
< error
| 1064 42000 …line 1 column 7 near "all := 1" 
> exec ^
: string
| set allow_missing_files := 1
< ok 0 0
> exec ^
: string
| set alter := 1
< error
| 1064 42000 …line 1 column 9 near "alter := 1" 
> exec ^
: string
| set always := 1
< ok 0 0
> exec ^
: string
| set analyze := 1
< error
| 1064 42000 …line 1 column 11 near "analyze := 1" 
> exec ^
: string
| set and := 1
< error
| 1064 42000 …line 1 column 7 near "and := 1" 
> exec ^
: string
| set any := 1
< ok 0 0
> exec ^
: string
| set array := 1
< error
| 1064 42000 …line 1 column 9 near "array := 1" 
> exec ^
: string
| set as := 1
< error
| 1064 42000 …line 1 column 6 near "as := 1" 
> exec ^
: string
| set asc := 1
< error
| 1064 42000 …line 1 column 7 near "asc := 1" 
> exec ^
: string
| set ascii := 1
< ok 0 0
> exec ^
: string
| set asensitive := 1
< ok 0 0
> exec ^
: string
| set assign_gtids_to_anonymous_transactions := 1
< ok 0 0
> exec ^
: string
| set at := 1
< ok 0 0
> exec ^
: string
| set attribute := 1
< ok 0 0
> exec ^
: string
| set authentication := 1
< ok 0 0
> exec ^
: string
| set auto := 1
< ok 0 0
> exec ^
: string
| set auto_id_cache := 1
< ok 0 0
> exec ^
: string
| set auto_increment := 1
< ok 0 0
> exec ^
: string
| set auto_random := 1
< ok 0 0
> exec ^
: string
| set auto_random_base := 1
< ok 0 0
> exec ^
: string
| set auto_refresh := 1
< ok 0 0
> exec ^
: string
| set auto_refresh_source := 1
< ok 0 0
> exec ^
: string
| set autoextend_size := 1
< ok 0 0
> exec ^
: string
| set avg := 1
< ok 0 0
> exec ^
: string
| set avg_row_length := 1
< ok 0 0
> exec ^
: string
| set background := 1
< ok 0 0
> exec ^
: string
| set backup := 1
< ok 0 0
> exec ^
: string
| set backups := 1
< ok 0 0
> exec ^
: string
| set batch := 1
< ok 0 0
> exec ^
: string
| set before := 1
< ok 0 0
> exec ^
: string
| set begin := 1
< ok 0 0
> exec ^
: string
| set bernoulli := 1
< ok 0 0
> exec ^
: string
| set between := 1
< error
| 1064 42000 …line 1 column 11 near "between := 1" 
> exec ^
: string
| set bigint := 1
< error
| 1064 42000 …line 1 column 10 near "bigint := 1" 
> exec ^
: string
| set binary := 1
< error
| 1064 42000 …line 1 column 10 near "binary := 1" 
> exec ^
: string
| set binding := 1
< ok 0 0
> exec ^
: string
| set binding_cache := 1
< ok 0 0
> exec ^
: string
| set bindings := 1
< ok 0 0
> exec ^
: string
| set binlog := 1
< ok 0 0
> exec ^
: string
| set bit := 1
< ok 0 0
> exec ^
: string
| set bit_and := 1
< ok 0 0
> exec ^
: string
| set bit_or := 1
< ok 0 0
> exec ^
: string
| set bit_xor := 1
< ok 0 0
> exec ^
: string
| set blob := 1
< error
| 1064 42000 …line 1 column 8 near "blob := 1" 
> exec ^
: string
| set block := 1
< ok 0 0
> exec ^
: string
| set bool := 1
< ok 0 0
> exec ^
: string
| set boolean := 1
< ok 0 0
> exec ^
: string
| set both := 1
< error
| 1064 42000 …line 1 column 8 near "both := 1" 
> exec ^
: string
| set btree := 1
< ok 0 0
> exec ^
: string
| set buckets := 1
< ok 0 0
> exec ^
: string
| set builtins := 1
< ok 0 0
> exec ^
: string
| set bulk := 1
< ok 0 0
> exec ^
: string
| set burstable := 1
< ok 0 0
> exec ^
: string
| set by := 1
< error
| 1064 42000 …line 1 column 6 near "by := 1" 
> exec ^
: string
| set byte := 1
< ok 0 0
> exec ^
: string
| set cache := 1
< ok 0 0
> exec ^
: string
| set calibrate := 1
< ok 0 0
> exec ^
: string
| set call := 1
< error
| 1064 42000 …line 1 column 8 near "call := 1" 
> exec ^
: string
| set cascade := 1
< error
| 1064 42000 …line 1 column 11 near "cascade := 1" 
> exec ^
: string
| set cascaded := 1
< ok 0 0
> exec ^
: string
| set case := 1
< error
| 1064 42000 …line 1 column 8 near "case := 1" 
> exec ^
: string
| set cast := 1
< ok 0 0
> exec ^
: string
| set catalog_name := 1
< ok 0 0
> exec ^
: string
| set chain := 1
< ok 0 0
> exec ^
: string
| set challenge_response := 1
< ok 0 0
> exec ^
: string
| set change := 1
< error
| 1064 42000 …line 1 column 10 near "change := 1" 
> exec ^
: string
| set changed := 1
< ok 0 0
> exec ^
: string
| set channel := 1
< ok 0 0
> exec ^
: string
| set char := 1
< error
| 1064 42000 …line 1 column 11 near ":= 1" 
> exec ^
: string
| set character := 1
< error
| 1064 42000 …line 1 column 16 near ":= 1" 
> exec ^
: string
| set charset := 1
< ok 0 0
> exec ^
: string
| set check := 1
< error
| 1064 42000 …line 1 column 9 near "check := 1" 
> exec ^
: string
| set checksum := 1
< ok 0 0
> exec ^
: string
| set cipher := 1
< ok 0 0
> exec ^
: string
| set class_origin := 1
< ok 0 0
> exec ^
: string
| set client := 1
< ok 0 0
> exec ^
: string
| set clone := 1
< ok 0 0
> exec ^
: string
| set close := 1
< ok 0 0
> exec ^
: string
| set clustered := 1
< ok 0 0
> exec ^
: string
| set coalesce := 1
< ok 0 0
> exec ^
: string
| set code := 1
< ok 0 0
> exec ^
: string
| set collate := 1
< error
| 1064 42000 …line 1 column 11 near "collate := 1" 
> exec ^
: string
| set collation := 1
< ok 0 0
> exec ^
: string
| set column := 1
< error
| 1064 42000 …line 1 column 10 near "column := 1" 
> exec ^
: string
| set column_format := 1
< ok 0 0
> exec ^
: string
| set column_name := 1
< ok 0 0
> exec ^
: string
| set column_stats_usage := 1
< ok 0 0
> exec ^
: string
| set columns := 1
< ok 0 0
> exec ^
: string
| set comment := 1
< ok 0 0
> exec ^
: string
| set commit := 1
< ok 0 0
> exec ^
: string
| set committed := 1
< ok 0 0
> exec ^
: string
| set compact := 1
< ok 0 0
> exec ^
: string
| set completion := 1
< ok 0 0
> exec ^
: string
| set component := 1
< ok 0 0
> exec ^
: string
| set compressed := 1
< ok 0 0
> exec ^
: string
| set compression := 1
< ok 0 0
> exec ^
: string
| set concat := 1
< ok 0 0
> exec ^
: string
| set concurrent := 1
< ok 0 0
> exec ^
: string
| set condition := 1
< ok 0 0
> exec ^
: string
| set config := 1
< ok 0 0
> exec ^
: string
| set connection := 1
< ok 0 0
> exec ^
: string
| set consistent := 1
< ok 0 0
> exec ^
: string
| set constraint := 1
< error
| 1064 42000 …line 1 column 14 near "constraint := 1" 
> exec ^
: string
| set constraint_catalog := 1
< ok 0 0
> exec ^
: string
| set constraint_name := 1
< ok 0 0
> exec ^
: string
| set constraint_schema := 1
< ok 0 0
> exec ^
: string
| set constraints := 1
< ok 0 0
> exec ^
: string
| set contains := 1
< ok 0 0
> exec ^
: string
| set context := 1
< ok 0 0
> exec ^
: string
| set continue := 1
< error
| 1064 42000 …line 1 column 12 near "continue := 1" 
> exec ^
: string
| set convert := 1
< error
| 1064 42000 …line 1 column 11 near "convert := 1" 
> exec ^
: string
| set cooldown := 1
< ok 0 0
> exec ^
: string
| set count := 1
< ok 0 0
> exec ^
: string
| set cpu := 1
< ok 0 0
> exec ^
: string
| set create := 1
< error
| 1064 42000 …line 1 column 10 near "create := 1" 
> exec ^
: string
| set cross := 1
< error
| 1064 42000 …line 1 column 9 near "cross := 1" 
> exec ^
: string
| set cube := 1
< ok 0 0
> exec ^
: string
| set cume_dist := 1
< error
| 1064 42000 …line 1 column 13 near "cume_dist := 1" 
> exec ^
: string
| set curdate := 1
< ok 0 0
> exec ^
: string
| set current := 1
< ok 0 0
> exec ^
: string
| set current_date := 1
< error
| 1064 42000 …line 1 column 16 near "current_date := 1" 
> exec ^
: string
| set current_time := 1
< error
| 1064 42000 …line 1 column 16 near "current_time := 1" 
> exec ^
: string
| set current_timestamp := 1
< error
| 1064 42000 …line 1 column 21 near "current_timestamp := 1" 
> exec ^
: string
| set current_user := 1
< error
| 1064 42000 …line 1 column 16 near "current_user := 1" 
> exec ^
: string
| set cursor := 1
< error
| 1064 42000 …line 1 column 10 near "cursor := 1" 
> exec ^
: string
| set cursor_name := 1
< ok 0 0
> exec ^
: string
| set curtime := 1
< ok 0 0
> exec ^
: string
| set cycle := 1
< ok 0 0
> exec ^
: string
| set data := 1
< ok 0 0
> exec ^
: string
| set database := 1
< error
| 1064 42000 …line 1 column 12 near "database := 1" 
> exec ^
: string
| set databases := 1
< error
| 1064 42000 …line 1 column 13 near "databases := 1" 
> exec ^
: string
| set datafile := 1
< ok 0 0
> exec ^
: string
| set date := 1
< ok 0 0
> exec ^
: string
| set date_add := 1
< ok 0 0
> exec ^
: string
| set date_sub := 1
< ok 0 0
> exec ^
: string
| set datetime := 1
< ok 0 0
> exec ^
: string
| set day := 1
< ok 0 0
> exec ^
: string
| set day_hour := 1
< error
| 1064 42000 …line 1 column 12 near "day_hour := 1" 
> exec ^
: string
| set day_microsecond := 1
< error
| 1064 42000 …line 1 column 19 near "day_microsecond := 1" 
> exec ^
: string
| set day_minute := 1
< error
| 1064 42000 …line 1 column 14 near "day_minute := 1" 
> exec ^
: string
| set day_second := 1
< error
| 1064 42000 …line 1 column 14 near "day_second := 1" 
> exec ^
: string
| set deallocate := 1
< ok 0 0
> exec ^
: string
| set dec := 1
< error
| 1064 42000 …line 1 column 7 near "dec := 1" 
> exec ^
: string
| set decimal := 1
< error
| 1064 42000 …line 1 column 11 near "decimal := 1" 
> exec ^
: string
| set declare := 1
< ok 0 0
> exec ^
: string
| set default := 1
< error
| 1064 42000 …line 1 column 14 near ":= 1" 
> exec ^
: string
| set default_auth := 1
< ok 0 0
> exec ^
: string
| set defined := 1
< ok 0 0
> exec ^
: string
| set definer := 1
< ok 0 0
> exec ^
: string
| set definition := 1
< ok 0 0
> exec ^
: string
| set delay_key_write := 1
< ok 0 0
> exec ^
: string
| set delayed := 1
< error
| 1064 42000 …line 1 column 11 near "delayed := 1" 
> exec ^
: string
| set delete := 1
< error
| 1064 42000 …line 1 column 10 near "delete := 1" 
> exec ^
: string
| set dense_rank := 1
< error
| 1064 42000 …line 1 column 14 near "dense_rank := 1" 
> exec ^
: string
| set desc := 1
< error
| 1064 42000 …line 1 column 8 near "desc := 1" 
> exec ^
: string
| set describe := 1
< error
| 1064 42000 …line 1 column 12 near "describe := 1" 
> exec ^
: string
| set description := 1
< ok 0 0
> exec ^
: string
| set deterministic := 1
< ok 0 0
> exec ^
: string
| set diagnostics := 1
< ok 0 0
> exec ^
: string
| set digest := 1
< ok 0 0
> exec ^
: string
| set directory := 1
< ok 0 0
> exec ^
: string
| set disable := 1
< ok 0 0
> exec ^
: string
| set disabled := 1
< ok 0 0
> exec ^
: string
| set discard := 1
< ok 0 0
> exec ^
: string
| set disk := 1
< ok 0 0
> exec ^
: string
| set distinct := 1
< error
| 1064 42000 …line 1 column 12 near "distinct := 1" 
> exec ^
: string
| set distinctrow := 1
< error
| 1064 42000 …line 1 column 15 near "distinctrow := 1" 
> exec ^
: string
| set distribute := 1
< ok 0 0
> exec ^
: string
| set distribution := 1
< ok 0 0
> exec ^
: string
| set distributions := 1
< ok 0 0
> exec ^
: string
| set div := 1
< error
| 1064 42000 …line 1 column 7 near "div := 1" 
> exec ^
: string
| set do := 1
< ok 0 0
> exec ^
: string
| set double := 1
< error
| 1064 42000 …line 1 column 10 near "double := 1" 
> exec ^
: string
| set drop := 1
< error
| 1064 42000 …line 1 column 8 near "drop := 1" 
> exec ^
: string
| set dry := 1
< ok 0 0
> exec ^
: string
| set dryrun := 1
< ok 0 0
> exec ^
: string
| set dual := 1
< error
| 1064 42000 …line 1 column 8 near "dual := 1" 
> exec ^
: string
| set duality := 1
< ok 0 0
> exec ^
: string
| set dumpfile := 1
< ok 0 0
> exec ^
: string
| set duplicate := 1
< ok 0 0
> exec ^
: string
| set duration := 1
< ok 0 0
> exec ^
: string
| set dynamic := 1
< ok 0 0
> exec ^
: string
| set each := 1
< ok 0 0
> exec ^
: string
| set else := 1
< error
| 1064 42000 …line 1 column 8 near "else := 1" 
> exec ^
: string
| set elseif := 1
< error
| 1064 42000 …line 1 column 10 near "elseif := 1" 
> exec ^
: string
| set empty := 1
< ok 0 0
> exec ^
: string
| set enable := 1
< ok 0 0
> exec ^
: string
| set enabled := 1
< ok 0 0
> exec ^
: string
| set enclosed := 1
< error
| 1064 42000 …line 1 column 12 near "enclosed := 1" 
> exec ^
: string
| set encryption := 1
< ok 0 0
> exec ^
: string
| set end := 1
< ok 0 0
> exec ^
: string
| set end_time := 1
< ok 0 0
> exec ^
: string
| set ends := 1
< ok 0 0
> exec ^
: string
| set enforced := 1
< ok 0 0
> exec ^
: string
| set engine := 1
< ok 0 0
> exec ^
: string
| set engine_attribute := 1
< ok 0 0
> exec ^
: string
| set engines := 1
< ok 0 0
> exec ^
: string
| set enum := 1
< ok 0 0
> exec ^
: string
| set error := 1
< ok 0 0
> exec ^
: string
| set errors := 1
< ok 0 0
> exec ^
: string
| set escape := 1
< ok 0 0
> exec ^
: string
| set escaped := 1
< error
| 1064 42000 …line 1 column 11 near "escaped := 1" 
> exec ^
: string
| set event := 1
< ok 0 0
> exec ^
: string
| set events := 1
< ok 0 0
> exec ^
: string
| set every := 1
< ok 0 0
> exec ^
: string
| set exact := 1
< ok 0 0
> exec ^
: string
| set except := 1
< error
| 1064 42000 …line 1 column 10 near "except := 1" 
> exec ^
: string
| set exchange := 1
< ok 0 0
> exec ^
: string
| set exclude := 1
< ok 0 0
> exec ^
: string
| set exec_elapsed := 1
< ok 0 0
> exec ^
: string
| set execute := 1
< ok 0 0
> exec ^
: string
| set exists := 1
< error
| 1064 42000 …line 1 column 10 near "exists := 1" 
> exec ^
: string
| set exit := 1
< error
| 1064 42000 …line 1 column 8 near "exit := 1" 
> exec ^
: string
| set expansion := 1
< ok 0 0
> exec ^
: string
| set expire := 1
< ok 0 0
> exec ^
: string
| set explain := 1
< error
| 1064 42000 …line 1 column 11 near "explain := 1" 
> exec ^
: string
| set export := 1
< ok 0 0
> exec ^
: string
| set extended := 1
< ok 0 0
> exec ^
: string
| set extent_size := 1
< ok 0 0
> exec ^
: string
| set external := 1
< ok 0 0
> exec ^
: string
| set external_format := 1
< ok 0 0
> exec ^
: string
| set extract := 1
< ok 0 0
> exec ^
: string
| set factor := 1
< ok 0 0
> exec ^
: string
| set failed_login_attempts := 1
< ok 0 0
> exec ^
: string
| set false := 1
< error
| 1064 42000 …line 1 column 9 near "false := 1" 
> exec ^
: string
| set fast := 1
< ok 0 0
> exec ^
: string
| set faults := 1
< ok 0 0
> exec ^
: string
| set fetch := 1
< error
| 1064 42000 …line 1 column 9 near "fetch := 1" 
> exec ^
: string
| set fields := 1
< ok 0 0
> exec ^
: string
| set file := 1
< ok 0 0
> exec ^
: string
| set file_block_size := 1
< ok 0 0
> exec ^
: string
| set file_format := 1
< ok 0 0
> exec ^
: string
| set file_name := 1
< ok 0 0
> exec ^
: string
| set file_pattern := 1
< ok 0 0
> exec ^
: string
| set file_prefix := 1
< ok 0 0
> exec ^
: string
| set files := 1
< ok 0 0
> exec ^
: string
| set filter := 1
< ok 0 0
> exec ^
: string
| set finish := 1
< ok 0 0
> exec ^
: string
| set first := 1
< ok 0 0
> exec ^
: string
| set first_value := 1
< error
| 1064 42000 …line 1 column 15 near "first_value := 1" 
> exec ^
: string
| set fixed := 1
< ok 0 0
> exec ^
: string
| set flashback := 1
< ok 0 0
> exec ^
: string
| set float := 1
< error
| 1064 42000 …line 1 column 9 near "float := 1" 
> exec ^
: string
| set float4 := 1
< error
| 1064 42000 …line 1 column 10 near "float4 := 1" 
> exec ^
: string
| set float8 := 1
< error
| 1064 42000 …line 1 column 10 near "float8 := 1" 
> exec ^
: string
| set flush := 1
< ok 0 0
> exec ^
: string
| set follower_constraints := 1
< ok 0 0
> exec ^
: string
| set followers := 1
< ok 0 0
> exec ^
: string
| set following := 1
< ok 0 0
> exec ^
: string
| set follows := 1
< ok 0 0
> exec ^
: string
| set for := 1
< error
| 1064 42000 …line 1 column 7 near "for := 1" 
> exec ^
: string
| set force := 1
< error
| 1064 42000 …line 1 column 9 near "force := 1" 
> exec ^
: string
| set foreign := 1
< error
| 1064 42000 …line 1 column 11 near "foreign := 1" 
> exec ^
: string
| set format := 1
< ok 0 0
> exec ^
: string
| set found := 1
< ok 0 0
> exec ^
: string
| set from := 1
< error
| 1064 42000 …line 1 column 8 near "from := 1" 
> exec ^
: string
| set full := 1
< ok 0 0
> exec ^
: string
| set fulltext := 1
< error
| 1064 42000 …line 1 column 12 near "fulltext := 1" 
> exec ^
: string
| set function := 1
< ok 0 0
> exec ^
: string
| set general := 1
< ok 0 0
> exec ^
: string
| set generate := 1
< ok 0 0
> exec ^
: string
| set generated := 1
< error
| 1064 42000 …line 1 column 13 near "generated := 1" 
> exec ^
: string
| set geomcollection := 1
< ok 0 0
> exec ^
: string
| set geometry := 1
< ok 0 0
> exec ^
: string
| set geometrycollection := 1
< ok 0 0
> exec ^
: string
| set get := 1
< ok 0 0
> exec ^
: string
| set get_format := 1
< ok 0 0
> exec ^
: string
| set get_source_public_key := 1
< ok 0 0
> exec ^
: string
| set global := 1
< ok 0 0
> exec ^
: string
| set grant := 1
< error
| 1064 42000 …line 1 column 9 near "grant := 1" 
> exec ^
: string
| set grants := 1
< ok 0 0
> exec ^
: string
| set greatest := 1
< ok 0 0
> exec ^
: string
| set group := 1
< error
| 1064 42000 …line 1 column 9 near "group := 1" 
> exec ^
: string
| set group_concat := 1
< ok 0 0
> exec ^
: string
| set group_replication := 1
< ok 0 0
> exec ^
: string
| set grouping := 1
< ok 0 0
> exec ^
: string
| set groups := 1
< error
| 1064 42000 …line 1 column 10 near "groups := 1" 
> exec ^
: string
| set gtid_only := 1
< ok 0 0
> exec ^
: string
| set gtids := 1
< ok 0 0
> exec ^
: string
| set guided := 1
< ok 0 0
> exec ^
: string
| set handler := 1
< ok 0 0
> exec ^
: string
| set hash := 1
< ok 0 0
> exec ^
: string
| set having := 1
< error
| 1064 42000 …line 1 column 10 near "having := 1" 
> exec ^
: string
| set header := 1
< ok 0 0
> exec ^
: string
| set help := 1
< ok 0 0
> exec ^
: string
| set high := 1
< ok 0 0
> exec ^
: string
| set high_priority := 1
< error
| 1064 42000 …line 1 column 17 near "high_priority := 1" 
> exec ^
: string
| set histogram := 1
< ok 0 0
> exec ^
: string
| set histograms_in_flight := 1
< ok 0 0
> exec ^
: string
| set history := 1
< ok 0 0
> exec ^
: string
| set host := 1
< ok 0 0
> exec ^
: string
| set hosts := 1
< ok 0 0
> exec ^
: string
| set hour := 1
< ok 0 0
> exec ^
: string
| set hour_microsecond := 1
< error
| 1064 42000 …line 1 column 20 near "hour_microsecond := 1" 
> exec ^
: string
| set hour_minute := 1
< error
| 1064 42000 …line 1 column 15 near "hour_minute := 1" 
> exec ^
: string
| set hour_second := 1
< error
| 1064 42000 …line 1 column 15 near "hour_second := 1" 
> exec ^
: string
| set identified := 1
< ok 0 0
> exec ^
: string
| set if := 1
< error
| 1064 42000 …line 1 column 6 near "if := 1" 
> exec ^
: string
| set ignore := 1
< error
| 1064 42000 …line 1 column 10 near "ignore := 1" 
> exec ^
: string
| set ignore_server_ids := 1
< ok 0 0
> exec ^
: string
| set ilike := 1
< error
| 1064 42000 …line 1 column 9 near "ilike := 1" 
> exec ^
: string
| set import := 1
< ok 0 0
> exec ^
: string
| set in := 1
< error
| 1064 42000 …line 1 column 6 near "in := 1" 
> exec ^
: string
| set inactive := 1
< ok 0 0
> exec ^
: string
| set increment := 1
< ok 0 0
> exec ^
: string
| set index := 1
< error
| 1064 42000 …line 1 column 9 near "index := 1" 
> exec ^
: string
| set indexes := 1
< ok 0 0
> exec ^
: string
| set infile := 1
< error
| 1064 42000 …line 1 column 10 near "infile := 1" 
> exec ^
: string
| set initial := 1
< ok 0 0
> exec ^
: string
| set initial_size := 1
< ok 0 0
> exec ^
: string
| set initiate := 1
< ok 0 0
> exec ^
: string
| set inner := 1
< error
| 1064 42000 …line 1 column 9 near "inner := 1" 
> exec ^
: string
| set inout := 1
< error
| 1064 42000 …line 1 column 9 near "inout := 1" 
> exec ^
: string
| set insensitive := 1
< ok 0 0
> exec ^
: string
| set insert := 1
< error
| 1064 42000 …line 1 column 10 near "insert := 1" 
> exec ^
: string
| set insert_method := 1
< ok 0 0
> exec ^
: string
| set install := 1
< ok 0 0
> exec ^
: string
| set instance := 1
< ok 0 0
> exec ^
: string
| set int := 1
< error
| 1064 42000 …line 1 column 7 near "int := 1" 
> exec ^
: string
| set int1 := 1
< error
| 1064 42000 …line 1 column 8 near "int1 := 1" 
> exec ^
: string
| set int2 := 1
< error
| 1064 42000 …line 1 column 8 near "int2 := 1" 
> exec ^
: string
| set int3 := 1
< error
| 1064 42000 …line 1 column 8 near "int3 := 1" 
> exec ^
: string
| set int4 := 1
< error
| 1064 42000 …line 1 column 8 near "int4 := 1" 
> exec ^
: string
| set int8 := 1
< error
| 1064 42000 …line 1 column 8 near "int8 := 1" 
> exec ^
: string
| set integer := 1
< error
| 1064 42000 …line 1 column 11 near "integer := 1" 
> exec ^
: string
| set intersect := 1
< error
| 1064 42000 …line 1 column 13 near "intersect := 1" 
> exec ^
: string
| set interval := 1
< error
| 1064 42000 …line 1 column 12 near "interval := 1" 
> exec ^
: string
| set into := 1
< error
| 1064 42000 …line 1 column 8 near "into := 1" 
> exec ^
: string
| set invisible := 1
< ok 0 0
> exec ^
: string
| set invoker := 1
< ok 0 0
> exec ^
: string
| set io := 1
< ok 0 0
> exec ^
: string
| set io_after_gtids := 1
< ok 0 0
> exec ^
: string
| set io_before_gtids := 1
< ok 0 0
> exec ^
: string
| set io_thread := 1
< ok 0 0
> exec ^
: string
| set ipc := 1
< ok 0 0
> exec ^
: string
| set is := 1
< error
| 1064 42000 …line 1 column 6 near "is := 1" 
> exec ^
: string
| set isolation := 1
< ok 0 0
> exec ^
: string
| set issuer := 1
< ok 0 0
> exec ^
: string
| set iterate := 1
< error
| 1064 42000 …line 1 column 11 near "iterate := 1" 
> exec ^
: string
| set job := 1
< ok 0 0
> exec ^
: string
| set jobs := 1
< ok 0 0
> exec ^
: string
| set join := 1
< error
| 1064 42000 …line 1 column 8 near "join := 1" 
> exec ^
: string
| set json := 1
< ok 0 0
> exec ^
: string
| set json_arrayagg := 1
< ok 0 0
> exec ^
: string
| set json_duality_object := 1
< ok 0 0
> exec ^
: string
| set json_objectagg := 1
< ok 0 0
> exec ^
: string
| set json_table := 1
< ok 0 0
> exec ^
: string
| set json_value := 1
< ok 0 0
> exec ^
: string
| set key := 1
< error
| 1064 42000 …line 1 column 7 near "key := 1" 
> exec ^
: string
| set key_block_size := 1
< ok 0 0
> exec ^
: string
| set keyring := 1
< ok 0 0
> exec ^
: string
| set keys := 1
< error
| 1064 42000 …line 1 column 8 near "keys := 1" 
> exec ^
: string
| set kill := 1
< error
| 1064 42000 …line 1 column 8 near "kill := 1" 
> exec ^
: string
| set labels := 1
< ok 0 0
> exec ^
: string
| set lag := 1
< error
| 1064 42000 …line 1 column 7 near "lag := 1" 
> exec ^
: string
| set language := 1
< ok 0 0
> exec ^
: string
| set last := 1
< ok 0 0
> exec ^
: string
| set last_value := 1
< error
| 1064 42000 …line 1 column 14 near "last_value := 1" 
> exec ^
: string
| set lateral := 1
< ok 0 0
> exec ^
: string
| set lead := 1
< error
| 1064 42000 …line 1 column 8 near "lead := 1" 
> exec ^
: string
| set leader_constraints := 1
< ok 0 0
> exec ^
: string
| set leading := 1
< error
| 1064 42000 …line 1 column 11 near "leading := 1" 
> exec ^
: string
| set learner_constraints := 1
< ok 0 0
> exec ^
: string
| set learners := 1
< ok 0 0
> exec ^
: string
| set least := 1
< ok 0 0
> exec ^
: string
| set leave := 1
< error
| 1064 42000 …line 1 column 9 near "leave := 1" 
> exec ^
: string
| set leaves := 1
< ok 0 0
> exec ^
: string
| set left := 1
< error
| 1064 42000 …line 1 column 8 near "left := 1" 
> exec ^
: string
| set less := 1
< ok 0 0
> exec ^
: string
| set level := 1
< ok 0 0
> exec ^
: string
| set library := 1
< ok 0 0
> exec ^
: string
| set like := 1
< error
| 1064 42000 …line 1 column 8 near "like := 1" 
> exec ^
: string
| set limit := 1
< error
| 1064 42000 …line 1 column 9 near "limit := 1" 
> exec ^
: string
| set linear := 1
< error
| 1064 42000 …line 1 column 10 near "linear := 1" 
> exec ^
: string
| set lines := 1
< error
| 1064 42000 …line 1 column 9 near "lines := 1" 
> exec ^
: string
| set linestring := 1
< ok 0 0
> exec ^
: string
| set list := 1
< ok 0 0
> exec ^
: string
| set load := 1
< error
| 1064 42000 …line 1 column 8 near "load := 1" 
> exec ^
: string
| set local := 1
< ok 0 0
> exec ^
: string
| set localtime := 1
< error
| 1064 42000 …line 1 column 13 near "localtime := 1" 
> exec ^
: string
| set localtimestamp := 1
< error
| 1064 42000 …line 1 column 18 near "localtimestamp := 1" 
> exec ^
: string
| set location := 1
< ok 0 0
> exec ^
: string
| set lock := 1
< error
| 1064 42000 …line 1 column 8 near "lock := 1" 
> exec ^
: string
| set locked := 1
< ok 0 0
> exec ^
: string
| set locks := 1
< ok 0 0
> exec ^
: string
| set log := 1
< ok 0 0
> exec ^
: string
| set logfile := 1
< ok 0 0
> exec ^
: string
| set logs := 1
< ok 0 0
> exec ^
: string
| set long := 1
< error
| 1064 42000 …line 1 column 8 near "long := 1" 
> exec ^
: string
| set longblob := 1
< error
| 1064 42000 …line 1 column 12 near "longblob := 1" 
> exec ^
: string
| set longtext := 1
< error
| 1064 42000 …line 1 column 12 near "longtext := 1" 
> exec ^
: string
| set loop := 1
< ok 0 0
> exec ^
: string
| set low := 1
< ok 0 0
> exec ^
: string
| set low_priority := 1
< error
| 1064 42000 …line 1 column 16 near "low_priority := 1" 
> exec ^
: string
| set manual := 1
< ok 0 0
> exec ^
: string
| set master := 1
< ok 0 0
> exec ^
: string
| set master_log_file := 1
< ok 0 0
> exec ^
: string
| set match := 1
< error
| 1064 42000 …line 1 column 9 near "match := 1" 
> exec ^
: string
| set materialized := 1
< ok 0 0
> exec ^
: string
| set max := 1
< ok 0 0
> exec ^
: string
| set max_connections_per_hour := 1
< ok 0 0
> exec ^
: string
| set max_queries_per_hour := 1
< ok 0 0
> exec ^
: string
| set max_rows := 1
< ok 0 0
> exec ^
: string
| set max_size := 1
< ok 0 0
> exec ^
: string
| set max_updates_per_hour := 1
< ok 0 0
> exec ^
: string
| set max_user_connections := 1
< ok 0 0
> exec ^
: string
| set maxvalue := 1
< error
| 1064 42000 …line 1 column 12 near "maxvalue := 1" 
> exec ^
: string
| set medium := 1
< ok 0 0
> exec ^
: string
| set mediumblob := 1
< error
| 1064 42000 …line 1 column 14 near "mediumblob := 1" 
> exec ^
: string
| set mediumint := 1
< error
| 1064 42000 …line 1 column 13 near "mediumint := 1" 
> exec ^
: string
| set mediumtext := 1
< error
| 1064 42000 …line 1 column 14 near "mediumtext := 1" 
> exec ^
: string
| set member := 1
< ok 0 0
> exec ^
: string
| set memory := 1
< ok 0 0
> exec ^
: string
| set merge := 1
< ok 0 0
> exec ^
: string
| set message_text := 1
< ok 0 0
> exec ^
: string
| set microsecond := 1
< ok 0 0
> exec ^
: string
| set mid := 1
< ok 0 0
> exec ^
: string
| set middleint := 1
< error
| 1064 42000 …line 1 column 13 near "middleint := 1" 
> exec ^
: string
| set migrate := 1
< ok 0 0
> exec ^
: string
| set min := 1
< ok 0 0
> exec ^
: string
| set min_rows := 1
< ok 0 0
> exec ^
: string
| set minute := 1
< ok 0 0
> exec ^
: string
| set minute_microsecond := 1
< error
| 1064 42000 …line 1 column 22 near "minute_microsecond := 1" 
> exec ^
: string
| set minute_second := 1
< error
| 1064 42000 …line 1 column 17 near "minute_second := 1" 
> exec ^
: string
| set minvalue := 1
< ok 0 0
> exec ^
: string
| set mod := 1
< error
| 1064 42000 …line 1 column 7 near "mod := 1" 
> exec ^
: string
| set mode := 1
< ok 0 0
> exec ^
: string
| set modifies := 1
< ok 0 0
> exec ^
: string
| set modify := 1
< ok 0 0
> exec ^
: string
| set month := 1
< ok 0 0
> exec ^
: string
| set multilinestring := 1
< ok 0 0
> exec ^
: string
| set multipoint := 1
< ok 0 0
> exec ^
: string
| set multipolygon := 1
< ok 0 0
> exec ^
: string
| set mutex := 1
< ok 0 0
> exec ^
: string
| set mysql_errno := 1
< ok 0 0
> exec ^
: string
| set name := 1
< ok 0 0
> exec ^
: string
| set names := 1
< ok 0 0
> exec ^
: string
| set national := 1
< ok 0 0
> exec ^
: string
| set natural := 1
< error
| 1064 42000 …line 1 column 11 near "natural := 1" 
> exec ^
: string
| set nchar := 1
< ok 0 0
> exec ^
: string
| set ndb := 1
< ok 0 0
> exec ^
: string
| set ndbcluster := 1
< ok 0 0
> exec ^
: string
| set nested := 1
< ok 0 0
> exec ^
: string
| set network_namespace := 1
< ok 0 0
> exec ^
: string
| set never := 1
< ok 0 0
> exec ^
: string
| set new := 1
< ok 0 0
> exec ^
: string
| set next := 1
< ok 0 0
> exec ^
: string
| set next_row_id := 1
< ok 0 0
> exec ^
: string
| set no := 1
< ok 0 0
> exec ^
: string
| set no_wait := 1
< ok 0 0
> exec ^
: string
| set no_write_to_binlog := 1
< error
| 1064 42000 …line 1 column 22 near "no_write_to_binlog := 1" 
> exec ^
: string
| set nocache := 1
< ok 0 0
> exec ^
: string
| set nocycle := 1
< ok 0 0
> exec ^
: string
| set nodegroup := 1
< ok 0 0
> exec ^
: string
| set nomaxvalue := 1
< ok 0 0
> exec ^
: string
| set nominvalue := 1
< ok 0 0
> exec ^
: string
| set nonclustered := 1
< ok 0 0
> exec ^
: string
| set none := 1
< ok 0 0
> exec ^
: string
| set not := 1
< error
| 1064 42000 …line 1 column 7 near "not := 1" 
> exec ^
: string
| set now := 1
< ok 0 0
> exec ^
: string
| set nowait := 1
< ok 0 0
> exec ^
: string
| set nth_value := 1
< error
| 1064 42000 …line 1 column 13 near "nth_value := 1" 
> exec ^
: string
| set ntile := 1
< error
| 1064 42000 …line 1 column 9 near "ntile := 1" 
> exec ^
: string
| set null := 1
< error
| 1064 42000 …line 1 column 8 near "null := 1" 
> exec ^
: string
| set nullif := 1
< ok 0 0
> exec ^
: string
| set nulls := 1
< ok 0 0
> exec ^
: string
| set number := 1
< ok 0 0
> exec ^
: string
| set numeric := 1
< error
| 1064 42000 …line 1 column 11 near "numeric := 1" 
> exec ^
: string
| set nvarchar := 1
< ok 0 0
> exec ^
: string
| set of := 1
< error
| 1064 42000 …line 1 column 6 near "of := 1" 
> exec ^
: string
| set off := 1
< ok 0 0
> exec ^
: string
| set offset := 1
< ok 0 0
> exec ^
: string
| set oj := 1
< ok 0 0
> exec ^
: string
| set old := 1
< ok 0 0
> exec ^
: string
| set oltp_read_only := 1
< ok 0 0
> exec ^
: string
| set oltp_read_write := 1
< ok 0 0
> exec ^
: string
| set oltp_write_only := 1
< ok 0 0
> exec ^
: string
| set on := 1
< error
| 1064 42000 …line 1 column 6 near "on := 1" 
> exec ^
: string
| set one := 1
< ok 0 0
> exec ^
: string
| set only := 1
< ok 0 0
> exec ^
: string
| set open := 1
< ok 0 0
> exec ^
: string
| set optimize := 1
< error
| 1064 42000 …line 1 column 12 near "optimize := 1" 
> exec ^
: string
| set optimizer_costs := 1
< ok 0 0
> exec ^
: string
| set option := 1
< error
| 1064 42000 …line 1 column 10 near "option := 1" 
> exec ^
: string
| set optional := 1
< ok 0 0
> exec ^
: string
| set optionally := 1
< error
| 1064 42000 …line 1 column 14 near "optionally := 1" 
> exec ^
: string
| set options := 1
< ok 0 0
> exec ^
: string
| set or := 1
< error
| 1064 42000 …line 1 column 6 near "or := 1" 
> exec ^
: string
| set order := 1
< error
| 1064 42000 …line 1 column 9 near "order := 1" 
> exec ^
: string
| set ordinality := 1
< ok 0 0
> exec ^
: string
| set organization := 1
< ok 0 0
> exec ^
: string
| set others := 1
< ok 0 0
> exec ^
: string
| set out := 1
< error
| 1064 42000 …line 1 column 7 near "out := 1" 
> exec ^
: string
| set outer := 1
< error
| 1064 42000 …line 1 column 9 near "outer := 1" 
> exec ^
: string
| set outfile := 1
< error
| 1064 42000 …line 1 column 11 near "outfile := 1" 
> exec ^
: string
| set over := 1
< error
| 1064 42000 …line 1 column 8 near "over := 1" 
> exec ^
: string
| set overlay := 1
< ok 0 0
> exec ^
: string
| set owner := 1
< ok 0 0
> exec ^
: string
| set pack_keys := 1
< ok 0 0
> exec ^
: string
| set page := 1
< ok 0 0
> exec ^
: string
| set parallel := 1
< ok 0 0
> exec ^
: string
| set parameters := 1
< ok 0 0
> exec ^
: string
| set parse_tree := 1
< ok 0 0
> exec ^
: string
| set parser := 1
< ok 0 0
> exec ^
: string
| set partial := 1
< ok 0 0
> exec ^
: string
| set partition := 1
< error
| 1064 42000 …line 1 column 13 near "partition := 1" 
> exec ^
: string
| set partitioning := 1
< ok 0 0
> exec ^
: string
| set partitions := 1
< ok 0 0
> exec ^
: string
| set password := 1
< ok 0 0
> exec ^
: string
| set password_lock_time := 1
< ok 0 0
> exec ^
: string
| set path := 1
< ok 0 0
> exec ^
: string
| set percent_rank := 1
< error
| 1064 42000 …line 1 column 16 near "percent_rank := 1" 
> exec ^
: string
| set persist := 1
< ok 0 0
> exec ^
: string
| set persist_only := 1
< ok 0 0
> exec ^
: string
| set phase := 1
< ok 0 0
> exec ^
: string
| set placement := 1
< ok 0 0
> exec ^
: string
| set placing := 1
< ok 0 0
> exec ^
: string
| set plan := 1
< ok 0 0
> exec ^
: string
| set plugin := 1
< ok 0 0
> exec ^
: string
| set plugin_dir := 1
< ok 0 0
> exec ^
: string
| set plugins := 1
< ok 0 0
> exec ^
: string
| set point := 1
< ok 0 0
> exec ^
: string
| set policy := 1
< ok 0 0
> exec ^
: string
| set polygon := 1
< ok 0 0
> exec ^
: string
| set port := 1
< ok 0 0
> exec ^
: string
| set position := 1
< ok 0 0
> exec ^
: string
| set pre_split_regions := 1
< ok 0 0
> exec ^
: string
| set precedes := 1
< ok 0 0
> exec ^
: string
| set preceding := 1
< ok 0 0
> exec ^
: string
| set precision := 1
< error
| 1064 42000 …line 1 column 13 near "precision := 1" 
> exec ^
: string
| set prepare := 1
< ok 0 0
> exec ^
: string
| set preserve := 1
< ok 0 0
> exec ^
: string
| set prev := 1
< ok 0 0
> exec ^
: string
| set primary := 1
< error
| 1064 42000 …line 1 column 11 near "primary := 1" 
> exec ^
: string
| set primary_region := 1
< ok 0 0
> exec ^
: string
| set priority := 1
< ok 0 0
> exec ^
: string
| set privilege_checks_user := 1
< ok 0 0
> exec ^
: string
| set privileges := 1
< ok 0 0
> exec ^
: string
| set procedure := 1
< error
| 1064 42000 …line 1 column 13 near "procedure := 1" 
> exec ^
: string
| set process := 1
< ok 0 0
> exec ^
: string
| set processed_keys := 1
< ok 0 0
> exec ^
: string
| set processlist := 1
< ok 0 0
> exec ^
: string
| set profile := 1
< ok 0 0
> exec ^
: string
| set profiles := 1
< ok 0 0
> exec ^
: string
| set proxy := 1
< ok 0 0
> exec ^
: string
| set purge := 1
< ok 0 0
> exec ^
: string
| set qualify := 1
< ok 0 0
> exec ^
: string
| set quarter := 1
< ok 0 0
> exec ^
: string
| set query := 1
< ok 0 0
> exec ^
: string
| set query_limit := 1
< ok 0 0
> exec ^
: string
| set quick := 1
< ok 0 0
> exec ^
: string
| set random := 1
< ok 0 0
> exec ^
: string
| set range := 1
< error
| 1064 42000 …line 1 column 9 near "range := 1" 
> exec ^
: string
| set rank := 1
< error
| 1064 42000 …line 1 column 8 near "rank := 1" 
> exec ^
: string
| set read := 1
< error
| 1064 42000 …line 1 column 8 near "read := 1" 
> exec ^
: string
| set read_only := 1
< ok 0 0
> exec ^
: string
| set read_write := 1
< ok 0 0
> exec ^
: string
| set reads := 1
< ok 0 0
> exec ^
: string
| set real := 1
< error
| 1064 42000 …line 1 column 8 near "real := 1" 
> exec ^
: string
| set rebuild := 1
< ok 0 0
> exec ^
: string
| set recover := 1
< ok 0 0
> exec ^
: string
| set recursive := 1
< error
| 1064 42000 …line 1 column 13 near "recursive := 1" 
> exec ^
: string
| set redo_buffer_size := 1
< ok 0 0
> exec ^
: string
| set redundant := 1
< ok 0 0
> exec ^
: string
| set reference := 1
< ok 0 0
> exec ^
: string
| set references := 1
< error
| 1064 42000 …line 1 column 14 near "references := 1" 
> exec ^
: string
| set regexp := 1
< error
| 1064 42000 …line 1 column 10 near "regexp := 1" 
> exec ^
: string
| set region := 1
< ok 0 0
> exec ^
: string
| set regions := 1
< ok 0 0
> exec ^
: string
| set registration := 1
< ok 0 0
> exec ^
: string
| set relational := 1
< ok 0 0
> exec ^
: string
| set relay := 1
< ok 0 0
> exec ^
: string
| set relay_log_file := 1
< ok 0 0
> exec ^
: string
| set relay_log_pos := 1
< ok 0 0
> exec ^
: string
| set relay_thread := 1
< ok 0 0
> exec ^
: string
| set relaylog := 1
< ok 0 0
> exec ^
: string
| set release := 1
< error
| 1064 42000 …line 1 column 11 near "release := 1" 
> exec ^
: string
| set reload := 1
< ok 0 0
> exec ^
: string
| set remove := 1
< ok 0 0
> exec ^
: string
| set rename := 1
< error
| 1064 42000 …line 1 column 10 near "rename := 1" 
> exec ^
: string
| set reorganize := 1
< ok 0 0
> exec ^
: string
| set repair := 1
< ok 0 0
> exec ^
: string
| set repeat := 1
< error
| 1064 42000 …line 1 column 10 near "repeat := 1" 
> exec ^
: string
| set repeatable := 1
< ok 0 0
> exec ^
: string
| set replace := 1
< error
| 1064 42000 …line 1 column 11 near "replace := 1" 
> exec ^
: string
| set replica := 1
< ok 0 0
> exec ^
: string
| set replicas := 1
< ok 0 0
> exec ^
: string
| set replicate_do_db := 1
< ok 0 0
> exec ^
: string
| set replicate_do_table := 1
< ok 0 0
> exec ^
: string
| set replicate_ignore_db := 1
< ok 0 0
> exec ^
: string
| set replicate_ignore_table := 1
< ok 0 0
> exec ^
: string
| set replicate_rewrite_db := 1
< ok 0 0
> exec ^
: string
| set replicate_wild_do_table := 1
< ok 0 0
> exec ^
: string
| set replicate_wild_ignore_table := 1
< ok 0 0
> exec ^
: string
| set replication := 1
< ok 0 0
> exec ^
: string
| set require := 1
< error
| 1064 42000 …line 1 column 11 near "require := 1" 
> exec ^
: string
| set require_row_format := 1
< ok 0 0
> exec ^
: string
| set require_table_primary_key_check := 1
< ok 0 0
> exec ^
: string
| set reset := 1
< ok 0 0
> exec ^
: string
| set resignal := 1
< ok 0 0
> exec ^
: string
| set resource := 1
< ok 0 0
> exec ^
: string
| set respect := 1
< ok 0 0
> exec ^
: string
| set restart := 1
< ok 0 0
> exec ^
: string
| set restore := 1
< ok 0 0
> exec ^
: string
| set restores := 1
< ok 0 0
> exec ^
: string
| set restrict := 1
< error
| 1064 42000 …line 1 column 12 near "restrict := 1" 
> exec ^
: string
| set resume := 1
< ok 0 0
> exec ^
: string
| set retain := 1
< ok 0 0
> exec ^
: string
| set return := 1
< ok 0 0
> exec ^
: string
| set returned_sqlstate := 1
< ok 0 0
> exec ^
: string
| set returning := 1
< ok 0 0
> exec ^
: string
| set returns := 1
< ok 0 0
> exec ^
: string
| set reuse := 1
< ok 0 0
> exec ^
: string
| set reverse := 1
< ok 0 0
> exec ^
: string
| set revoke := 1
< error
| 1064 42000 …line 1 column 10 near "revoke := 1" 
> exec ^
: string
| set right := 1
< error
| 1064 42000 …line 1 column 9 near "right := 1" 
> exec ^
: string
| set rlike := 1
< error
| 1064 42000 …line 1 column 9 near "rlike := 1" 
> exec ^
: string
| set role := 1
< ok 0 0
> exec ^
: string
| set rollback := 1
< ok 0 0
> exec ^
: string
| set rollup := 1
< ok 0 0
> exec ^
: string
| set rotate := 1
< ok 0 0
> exec ^
: string
| set routine := 1
< ok 0 0
> exec ^
: string
| set row := 1
< error
| 1064 42000 …line 1 column 7 near "row := 1" 
> exec ^
: string
| set row_count := 1
< ok 0 0
> exec ^
: string
| set row_format := 1
< ok 0 0
> exec ^
: string
| set row_number := 1
< error
| 1064 42000 …line 1 column 14 near "row_number := 1" 
> exec ^
: string
| set rows := 1
< error
| 1064 42000 …line 1 column 8 near "rows := 1" 
> exec ^
: string
| set rtree := 1
< ok 0 0
> exec ^
: string
| set ru := 1
< ok 0 0
> exec ^
: string
| set ru_per_sec := 1
< ok 0 0
> exec ^
: string
| set rule := 1
< ok 0 0
> exec ^
: string
| set run := 1
< ok 0 0
> exec ^
: string
| set s3 := 1
< ok 0 0
> exec ^
: string
| set savepoint := 1
< ok 0 0
> exec ^
: string
| set schedule := 1
< ok 0 0
> exec ^
: string
| set schema := 1
< error
| 1064 42000 …line 1 column 10 near "schema := 1" 
> exec ^
: string
| set schema_name := 1
< ok 0 0
> exec ^
: string
| set schemas := 1
< error
| 1064 42000 …line 1 column 11 near "schemas := 1" 
> exec ^
: string
| set second := 1
< ok 0 0
> exec ^
: string
| set second_microsecond := 1
< error
| 1064 42000 …line 1 column 22 near "second_microsecond := 1" 
> exec ^
: string
| set secondary := 1
< ok 0 0
> exec ^
: string
| set secondary_engine := 1
< ok 0 0
> exec ^
: string
| set secondary_engine_attribute := 1
< ok 0 0
> exec ^
: string
| set secondary_load := 1
< ok 0 0
> exec ^
: string
| set secondary_unload := 1
< ok 0 0
> exec ^
: string
| set security := 1
< ok 0 0
> exec ^
: string
| set select := 1
< error
| 1064 42000 …line 1 column 10 near "select := 1" 
> exec ^
: string
| set sensitive := 1
< ok 0 0
> exec ^
: string
| set separator := 1
< ok 0 0
> exec ^
: string
| set sequence := 1
< ok 0 0
> exec ^
: string
| set serial := 1
< ok 0 0
> exec ^
: string
| set serializable := 1
< ok 0 0
> exec ^
: string
| set server := 1
< ok 0 0
> exec ^
: string
| set session := 1
< ok 0 0
> exec ^
: string
| set session_states := 1
< ok 0 0
> exec ^
: string
| set session_user := 1
< ok 0 0
> exec ^
: string
| set set := 1
< error
| 1064 42000 …line 1 column 7 near "set := 1" 
> exec ^
: string
| set sets := 1
< ok 0 0
> exec ^
: string
| set shard_row_id_bits := 1
< ok 0 0
> exec ^
: string
| set share := 1
< ok 0 0
> exec ^
: string
| set show := 1
< error
| 1064 42000 …line 1 column 8 near "show := 1" 
> exec ^
: string
| set shutdown := 1
< ok 0 0
> exec ^
: string
| set signal := 1
< ok 0 0
> exec ^
: string
| set signed := 1
< ok 0 0
> exec ^
: string
| set similar := 1
< ok 0 0
> exec ^
: string
| set simple := 1
< ok 0 0
> exec ^
: string
| set skip := 1
< ok 0 0
> exec ^
: string
| set slave := 1
< ok 0 0
> exec ^
: string
| set slow := 1
< ok 0 0
> exec ^
: string
| set smallint := 1
< error
| 1064 42000 …line 1 column 12 near "smallint := 1" 
> exec ^
: string
| set snapshot := 1
< ok 0 0
> exec ^
: string
| set socket := 1
< ok 0 0
> exec ^
: string
| set some := 1
< ok 0 0
> exec ^
: string
| set soname := 1
< ok 0 0
> exec ^
: string
| set sounds := 1
< ok 0 0
> exec ^
: string
| set source := 1
< ok 0 0
> exec ^
: string
| set source_auto_position := 1
< ok 0 0
> exec ^
: string
| set source_bind := 1
< ok 0 0
> exec ^
: string
| set source_compression_algorithms := 1
< ok 0 0
> exec ^
: string
| set source_connect_retry := 1
< ok 0 0
> exec ^
: string
| set source_connection_auto_failover := 1
< ok 0 0
> exec ^
: string
| set source_delay := 1
< ok 0 0
> exec ^
: string
| set source_heartbeat_period := 1
< ok 0 0
> exec ^
: string
| set source_host := 1
< ok 0 0
> exec ^
: string
| set source_log_file := 1
< ok 0 0
> exec ^
: string
| set source_log_pos := 1
< ok 0 0
> exec ^
: string
| set source_password := 1
< ok 0 0
> exec ^
: string
| set source_port := 1
< ok 0 0
> exec ^
: string
| set source_public_key_path := 1
< ok 0 0
> exec ^
: string
| set source_retry_count := 1
< ok 0 0
> exec ^
: string
| set source_ssl := 1
< ok 0 0
> exec ^
: string
| set source_ssl_ca := 1
< ok 0 0
> exec ^
: string
| set source_ssl_capath := 1
< ok 0 0
> exec ^
: string
| set source_ssl_cert := 1
< ok 0 0
> exec ^
: string
| set source_ssl_cipher := 1
< ok 0 0
> exec ^
: string
| set source_ssl_crl := 1
< ok 0 0
> exec ^
: string
| set source_ssl_crlpath := 1
< ok 0 0
> exec ^
: string
| set source_ssl_key := 1
< ok 0 0
> exec ^
: string
| set source_ssl_verify_server_cert := 1
< ok 0 0
> exec ^
: string
| set source_tls_ciphersuites := 1
< ok 0 0
> exec ^
: string
| set source_tls_version := 1
< ok 0 0
> exec ^
: string
| set source_user := 1
< ok 0 0
> exec ^
: string
| set source_zstd_compression_level := 1
< ok 0 0
> exec ^
: string
| set spatial := 1
< error
| 1064 42000 …line 1 column 11 near "spatial := 1" 
> exec ^
: string
| set specific := 1
< ok 0 0
> exec ^
: string
| set split := 1
< ok 0 0
> exec ^
: string
| set sql := 1
< error
| 1064 42000 …line 1 column 7 near "sql := 1" 
> exec ^
: string
| set sql_after_gtids := 1
< ok 0 0
> exec ^
: string
| set sql_after_mts_gaps := 1
< ok 0 0
> exec ^
: string
| set sql_before_gtids := 1
< ok 0 0
> exec ^
: string
| set sql_big_result := 1
< error
| 1064 42000 …line 1 column 18 near "sql_big_result := 1" 
> exec ^
: string
| set sql_buffer_result := 1
< ok 0 0
> exec ^
: string
| set sql_calc_found_rows := 1
< error
| 1064 42000 …line 1 column 23 near "sql_calc_found_rows := 1" 
> exec ^
: string
| set sql_no_cache := 1
< ok 0 0
> exec ^
: string
| set sql_small_result := 1
< error
| 1064 42000 …line 1 column 20 near "sql_small_result := 1" 
> exec ^
: string
| set sql_thread := 1
< ok 0 0
> exec ^
: string
| set sql_tsi_day := 1
< ok 0 0
> exec ^
: string
| set sql_tsi_hour := 1
< ok 0 0
> exec ^
: string
| set sql_tsi_minute := 1
< ok 0 0
> exec ^
: string
| set sql_tsi_month := 1
< ok 0 0
> exec ^
: string
| set sql_tsi_quarter := 1
< ok 0 0
> exec ^
: string
| set sql_tsi_second := 1
< ok 0 0
> exec ^
: string
| set sql_tsi_week := 1
< ok 0 0
> exec ^
: string
| set sql_tsi_year := 1
< ok 0 0
> exec ^
: string
| set sqlexception := 1
< error
| 1064 42000 …line 1 column 16 near "sqlexception := 1" 
> exec ^
: string
| set sqlstate := 1
< error
| 1064 42000 …line 1 column 12 near "sqlstate := 1" 
> exec ^
: string
| set sqlwarning := 1
< error
| 1064 42000 …line 1 column 14 near "sqlwarning := 1" 
> exec ^
: string
| set srid := 1
< ok 0 0
> exec ^
: string
| set ssl := 1
< error
| 1064 42000 …line 1 column 7 near "ssl := 1" 
> exec ^
: string
| set st_collect := 1
< ok 0 0
> exec ^
: string
| set stacked := 1
< ok 0 0
> exec ^
: string
| set start := 1
< ok 0 0
> exec ^
: string
| set start_time := 1
< ok 0 0
> exec ^
: string
| set starting := 1
< error
| 1064 42000 …line 1 column 12 near "starting := 1" 
> exec ^
: string
| set starts := 1
< ok 0 0
> exec ^
: string
| set stats_auto_recalc := 1
< ok 0 0
> exec ^
: string
| set stats_buckets := 1
< ok 0 0
> exec ^
: string
| set stats_col_choice := 1
< ok 0 0
> exec ^
: string
| set stats_col_list := 1
< ok 0 0
> exec ^
: string
| set stats_extended := 1
< ok 0 0
> exec ^
: string
| set stats_healthy := 1
< ok 0 0
> exec ^
: string
| set stats_histograms := 1
< ok 0 0
> exec ^
: string
| set stats_locked := 1
< ok 0 0
> exec ^
: string
| set stats_meta := 1
< ok 0 0
> exec ^
: string
| set stats_persistent := 1
< ok 0 0
> exec ^
: string
| set stats_sample_pages := 1
< ok 0 0
> exec ^
: string
| set stats_sample_rate := 1
< ok 0 0
> exec ^
: string
| set stats_topn := 1
< ok 0 0
> exec ^
: string
| set status := 1
< ok 0 0
> exec ^
: string
| set std := 1
< ok 0 0
> exec ^
: string
| set stddev := 1
< ok 0 0
> exec ^
: string
| set stddev_pop := 1
< ok 0 0
> exec ^
: string
| set stddev_samp := 1
< ok 0 0
> exec ^
: string
| set stop := 1
< ok 0 0
> exec ^
: string
| set storage := 1
< ok 0 0
> exec ^
: string
| set stored := 1
< error
| 1064 42000 …line 1 column 10 near "stored := 1" 
> exec ^
: string
| set straight_join := 1
< error
| 1064 42000 …line 1 column 17 near "straight_join := 1" 
> exec ^
: string
| set stream := 1
< ok 0 0
> exec ^
: string
| set strict_load := 1
< ok 0 0
> exec ^
: string
| set string := 1
< ok 0 0
> exec ^
: string
| set subclass_origin := 1
< ok 0 0
> exec ^
: string
| set subdate := 1
< ok 0 0
> exec ^
: string
| set subject := 1
< ok 0 0
> exec ^
: string
| set subpartition := 1
< ok 0 0
> exec ^
: string
| set subpartitions := 1
< ok 0 0
> exec ^
: string
| set substr := 1
< ok 0 0
> exec ^
: string
| set substring := 1
< ok 0 0
> exec ^
: string
| set sum := 1
< ok 0 0
> exec ^
: string
| set super := 1
< ok 0 0
> exec ^
: string
| set survival_preferences := 1
< ok 0 0
> exec ^
: string
| set suspend := 1
< ok 0 0
> exec ^
: string
| set swaps := 1
< ok 0 0
> exec ^
: string
| set switch_group := 1
< ok 0 0
> exec ^
: string
| set switches := 1
< ok 0 0
> exec ^
: string
| set sysdate := 1
< ok 0 0
> exec ^
: string
| set system := 1
< ok 0 0
> exec ^
: string
| set system_user := 1
< ok 0 0
> exec ^
: string
| set table := 1
< error
| 1064 42000 …line 1 column 9 near "table := 1" 
> exec ^
: string
| set table_checksum := 1
< ok 0 0
> exec ^
: string
| set table_name := 1
< ok 0 0
> exec ^
: string
| set tables := 1
< ok 0 0
> exec ^
: string
| set tablesample := 1
< error
| 1064 42000 …line 1 column 15 near "tablesample := 1" 
> exec ^
: string
| set tablespace := 1
< ok 0 0
> exec ^
: string
| set task_types := 1
< ok 0 0
> exec ^
: string
| set temporary := 1
< ok 0 0
> exec ^
: string
| set temptable := 1
< ok 0 0
> exec ^
: string
| set terminated := 1
< error
| 1064 42000 …line 1 column 14 near "terminated := 1" 
> exec ^
: string
| set text := 1
< ok 0 0
> exec ^
: string
| set than := 1
< ok 0 0
> exec ^
: string
| set then := 1
< error
| 1064 42000 …line 1 column 8 near "then := 1" 
> exec ^
: string
| set thread_priority := 1
< ok 0 0
> exec ^
: string
| set ties := 1
< ok 0 0
> exec ^
: string
| set tiflash := 1
< ok 0 0
> exec ^
: string
| set time := 1
< ok 0 0
> exec ^
: string
| set timeout := 1
< ok 0 0
> exec ^
: string
| set timestamp := 1
< ok 0 0
> exec ^
: string
| set timestampadd := 1
< ok 0 0
> exec ^
: string
| set timestampdiff := 1
< ok 0 0
> exec ^
: string
| set tinyblob := 1
< error
| 1064 42000 …line 1 column 12 near "tinyblob := 1" 
> exec ^
: string
| set tinyint := 1
< error
| 1064 42000 …line 1 column 11 near "tinyint := 1" 
> exec ^
: string
| set tinytext := 1
< error
| 1064 42000 …line 1 column 12 near "tinytext := 1" 
> exec ^
: string
| set tls := 1
< ok 0 0
> exec ^
: string
| set to := 1
< error
| 1064 42000 …line 1 column 6 near "to := 1" 
> exec ^
: string
| set tpcc := 1
< ok 0 0
> exec ^
: string
| set tpch_10 := 1
< ok 0 0
> exec ^
: string
| set trailing := 1
< error
| 1064 42000 …line 1 column 12 near "trailing := 1" 
> exec ^
: string
| set transaction := 1
< ok 0 0
> exec ^
: string
| set trigger := 1
< error
| 1064 42000 …line 1 column 11 near "trigger := 1" 
> exec ^
: string
| set triggers := 1
< ok 0 0
> exec ^
: string
| set trim := 1
< ok 0 0
> exec ^
: string
| set true := 1
< error
| 1064 42000 …line 1 column 8 near "true := 1" 
> exec ^
: string
| set truncate := 1
< ok 0 0
> exec ^
: string
| set ttl := 1
< ok 0 0
> exec ^
: string
| set ttl_enable := 1
< ok 0 0
> exec ^
: string
| set ttl_job_interval := 1
< ok 0 0
> exec ^
: string
| set type := 1
< ok 0 0
> exec ^
: string
| set types := 1
< ok 0 0
> exec ^
: string
| set unbounded := 1
< ok 0 0
> exec ^
: string
| set uncommitted := 1
< ok 0 0
> exec ^
: string
| set undefined := 1
< ok 0 0
> exec ^
: string
| set undo := 1
< ok 0 0
> exec ^
: string
| set undo_buffer_size := 1
< ok 0 0
> exec ^
: string
| set undofile := 1
< ok 0 0
> exec ^
: string
| set unicode := 1
< ok 0 0
> exec ^
: string
| set uninstall := 1
< ok 0 0
> exec ^
: string
| set union := 1
< error
| 1064 42000 …line 1 column 9 near "union := 1" 
> exec ^
: string
| set unique := 1
< error
| 1064 42000 …line 1 column 10 near "unique := 1" 
> exec ^
: string
| set unknown := 1
< ok 0 0
> exec ^
: string
| set unlimited := 1
< ok 0 0
> exec ^
: string
| set unlock := 1
< error
| 1064 42000 …line 1 column 10 near "unlock := 1" 
> exec ^
: string
| set unregister := 1
< ok 0 0
> exec ^
: string
| set unsigned := 1
< error
| 1064 42000 …line 1 column 12 near "unsigned := 1" 
> exec ^
: string
| set until := 1
< error
| 1064 42000 …line 1 column 9 near "until := 1" 
> exec ^
: string
| set update := 1
< error
| 1064 42000 …line 1 column 10 near "update := 1" 
> exec ^
: string
| set upgrade := 1
< ok 0 0
> exec ^
: string
| set uri := 1
< ok 0 0
> exec ^
: string
| set url := 1
< ok 0 0
> exec ^
: string
| set usage := 1
< error
| 1064 42000 …line 1 column 9 near "usage := 1" 
> exec ^
: string
| set use := 1
< error
| 1064 42000 …line 1 column 7 near "use := 1" 
> exec ^
: string
| set use_frm := 1
< ok 0 0
> exec ^
: string
| set user := 1
< ok 0 0
> exec ^
: string
| set user_resources := 1
< ok 0 0
> exec ^
: string
| set using := 1
< error
| 1064 42000 …line 1 column 9 near "using := 1" 
> exec ^
: string
| set utc_date := 1
< error
| 1064 42000 …line 1 column 12 near "utc_date := 1" 
> exec ^
: string
| set utc_time := 1
< error
| 1064 42000 …line 1 column 12 near "utc_time := 1" 
> exec ^
: string
| set utc_timestamp := 1
< error
| 1064 42000 …line 1 column 17 near "utc_timestamp := 1" 
> exec ^
: string
| set utilization_limit := 1
< ok 0 0
> exec ^
: string
| set validate := 1
< ok 0 0
> exec ^
: string
| set validation := 1
< ok 0 0
> exec ^
: string
| set value := 1
< ok 0 0
> exec ^
: string
| set values := 1
< error
| 1064 42000 …line 1 column 10 near "values := 1" 
> exec ^
: string
| set var_pop := 1
< ok 0 0
> exec ^
: string
| set var_samp := 1
< ok 0 0
> exec ^
: string
| set varbinary := 1
< error
| 1064 42000 …line 1 column 13 near "varbinary := 1" 
> exec ^
: string
| set varchar := 1
< error
| 1064 42000 …line 1 column 11 near "varchar := 1" 
> exec ^
: string
| set varcharacter := 1
< error
| 1064 42000 …line 1 column 16 near "varcharacter := 1" 
> exec ^
: string
| set variables := 1
< ok 0 0
> exec ^
: string
| set variance := 1
< ok 0 0
> exec ^
: string
| set varying := 1
< error
| 1064 42000 …line 1 column 11 near "varying := 1" 
> exec ^
: string
| set vcpu := 1
< ok 0 0
> exec ^
: string
| set vector := 1
< ok 0 0
> exec ^
: string
| set verify_key_constraints := 1
< ok 0 0
> exec ^
: string
| set view := 1
< ok 0 0
> exec ^
: string
| set virtual := 1
< error
| 1064 42000 …line 1 column 11 near "virtual := 1" 
> exec ^
: string
| set visible := 1
< ok 0 0
> exec ^
: string
| set voter_constraints := 1
< ok 0 0
> exec ^
: string
| set voters := 1
< ok 0 0
> exec ^
: string
| set wait := 1
< ok 0 0
> exec ^
: string
| set warnings := 1
< ok 0 0
> exec ^
: string
| set watch := 1
< ok 0 0
> exec ^
: string
| set week := 1
< ok 0 0
> exec ^
: string
| set weight_string := 1
< ok 0 0
> exec ^
: string
| set when := 1
< error
| 1064 42000 …line 1 column 8 near "when := 1" 
> exec ^
: string
| set where := 1
< error
| 1064 42000 …line 1 column 9 near "where := 1" 
> exec ^
: string
| set while := 1
< error
| 1064 42000 …line 1 column 9 near "while := 1" 
> exec ^
: string
| set window := 1
< error
| 1064 42000 …line 1 column 10 near "window := 1" 
> exec ^
: string
| set with := 1
< error
| 1064 42000 …line 1 column 8 near "with := 1" 
> exec ^
: string
| set without := 1
< ok 0 0
> exec ^
: string
| set work := 1
< ok 0 0
> exec ^
: string
| set workload := 1
< ok 0 0
> exec ^
: string
| set wrapper := 1
< ok 0 0
> exec ^
: string
| set write := 1
< error
| 1064 42000 …line 1 column 9 near "write := 1" 
> exec ^
: string
| set x509 := 1
< ok 0 0
> exec ^
: string
| set xa := 1
< ok 0 0
> exec ^
: string
| set xid := 1
< ok 0 0
> exec ^
: string
| set xml := 1
< ok 0 0
> exec ^
: string
| set xor := 1
< error
| 1064 42000 …line 1 column 7 near "xor := 1" 
> exec ^
: string
| set year := 1
< ok 0 0
> exec ^
: string
| set year_month := 1
< error
| 1064 42000 …line 1 column 14 near "year_month := 1" 
> exec ^
: string
| set zerofill := 1
< error
| 1064 42000 …line 1 column 12 near "zerofill := 1" 
> exec ^
: string
| set zone := 1
< ok 0 0
> exec ^
: string
| set session absent = 1
< ok 0 0
> exec ^
: string
| set session accessible = 1
< ok 0 0
> exec ^
: string
| set session account = 1
< ok 0 0
> exec ^
: string
| set session action = 1
< ok 0 0
> exec ^
: string
| set session active = 1
< ok 0 0
> exec ^
: string
| set session add = 1
< error
| 1064 42000 …line 1 column 15 near "add = 1" 
> exec ^
: string
| set session adddate = 1
< ok 0 0
> exec ^
: string
| set session admin = 1
< ok 0 0
> exec ^
: string
| set session after = 1
< ok 0 0
> exec ^
: string
| set session against = 1
< ok 0 0
> exec ^
: string
| set session aggregate = 1
< ok 0 0
> exec ^
: string
| set session algorithm = 1
< ok 0 0
> exec ^
: string
| set session all = 1
< error
| 1064 42000 …line 1 column 15 near "all = 1" 
> exec ^
: string
| set session allow_missing_files = 1
< ok 0 0
> exec ^
: string
| set session alter = 1
< error
| 1064 42000 …line 1 column 17 near "alter = 1" 
> exec ^
: string
| set session always = 1
< ok 0 0
> exec ^
: string
| set session analyze = 1
< error
| 1064 42000 …line 1 column 19 near "analyze = 1" 
> exec ^
: string
| set session and = 1
< error
| 1064 42000 …line 1 column 15 near "and = 1" 
> exec ^
: string
| set session any = 1
< ok 0 0
> exec ^
: string
| set session array = 1
< error
| 1064 42000 …line 1 column 17 near "array = 1" 
> exec ^
: string
| set session as = 1
< error
| 1064 42000 …line 1 column 14 near "as = 1" 
> exec ^
: string
| set session asc = 1
< error
| 1064 42000 …line 1 column 15 near "asc = 1" 
> exec ^
: string
| set session ascii = 1
< ok 0 0
> exec ^
: string
| set session asensitive = 1
< ok 0 0
> exec ^
: string
| set session assign_gtids_to_anonymous_transactions = 1
< ok 0 0
> exec ^
: string
| set session at = 1
< ok 0 0
> exec ^
: string
| set session attribute = 1
< ok 0 0
> exec ^
: string
| set session authentication = 1
< ok 0 0
> exec ^
: string
| set session auto = 1
< ok 0 0
> exec ^
: string
| set session auto_id_cache = 1
< ok 0 0
> exec ^
: string
| set session auto_increment = 1
< ok 0 0
> exec ^
: string
| set session auto_random = 1
< ok 0 0
> exec ^
: string
| set session auto_random_base = 1
< ok 0 0
> exec ^
: string
| set session auto_refresh = 1
< ok 0 0
> exec ^
: string
| set session auto_refresh_source = 1
< ok 0 0
> exec ^
: string
| set session autoextend_size = 1
< ok 0 0
> exec ^
: string
| set session avg = 1
< ok 0 0
> exec ^
: string
| set session avg_row_length = 1
< ok 0 0
> exec ^
: string
| set session background = 1
< ok 0 0
> exec ^
: string
| set session backup = 1
< ok 0 0
> exec ^
: string
| set session backups = 1
< ok 0 0
> exec ^
: string
| set session batch = 1
< ok 0 0
> exec ^
: string
| set session before = 1
< ok 0 0
> exec ^
: string
| set session begin = 1
< ok 0 0
> exec ^
: string
| set session bernoulli = 1
< ok 0 0
> exec ^
: string
| set session between = 1
< error
| 1064 42000 …line 1 column 19 near "between = 1" 
> exec ^
: string
| set session bigint = 1
< error
| 1064 42000 …line 1 column 18 near "bigint = 1" 
> exec ^
: string
| set session binary = 1
< error
| 1064 42000 …line 1 column 18 near "binary = 1" 
> exec ^
: string
| set session binding = 1
< ok 0 0
> exec ^
: string
| set session binding_cache = 1
< ok 0 0
> exec ^
: string
| set session bindings = 1
< ok 0 0
> exec ^
: string
| set session binlog = 1
< ok 0 0
> exec ^
: string
| set session bit = 1
< ok 0 0
> exec ^
: string
| set session bit_and = 1
< ok 0 0
> exec ^
: string
| set session bit_or = 1
< ok 0 0
> exec ^
: string
| set session bit_xor = 1
< ok 0 0
> exec ^
: string
| set session blob = 1
< error
| 1064 42000 …line 1 column 16 near "blob = 1" 
> exec ^
: string
| set session block = 1
< ok 0 0
> exec ^
: string
| set session bool = 1
< ok 0 0
> exec ^
: string
| set session boolean = 1
< ok 0 0
> exec ^
: string
| set session both = 1
< error
| 1064 42000 …line 1 column 16 near "both = 1" 
> exec ^
: string
| set session btree = 1
< ok 0 0
> exec ^
: string
| set session buckets = 1
< ok 0 0
> exec ^
: string
| set session builtins = 1
< ok 0 0
> exec ^
: string
| set session bulk = 1
< ok 0 0
> exec ^
: string
| set session burstable = 1
< ok 0 0
> exec ^
: string
| set session by = 1
< error
| 1064 42000 …line 1 column 14 near "by = 1" 
> exec ^
: string
| set session byte = 1
< ok 0 0
> exec ^
: string
| set session cache = 1
< ok 0 0
> exec ^
: string
| set session calibrate = 1
< ok 0 0
> exec ^
: string
| set session call = 1
< error
| 1064 42000 …line 1 column 16 near "call = 1" 
> exec ^
: string
| set session cascade = 1
< error
| 1064 42000 …line 1 column 19 near "cascade = 1" 
> exec ^
: string
| set session cascaded = 1
< ok 0 0
> exec ^
: string
| set session case = 1
< error
| 1064 42000 …line 1 column 16 near "case = 1" 
> exec ^
: string
| set session cast = 1
< ok 0 0
> exec ^
: string
| set session catalog_name = 1
< ok 0 0
> exec ^
: string
| set session chain = 1
< ok 0 0
> exec ^
: string
| set session challenge_response = 1
< ok 0 0
> exec ^
: string
| set session change = 1
< error
| 1064 42000 …line 1 column 18 near "change = 1" 
> exec ^
: string
| set session changed = 1
< ok 0 0
> exec ^
: string
| set session channel = 1
< ok 0 0
> exec ^
: string
| set session char = 1
< error
| 1064 42000 …line 1 column 16 near "char = 1" 
> exec ^
: string
| set session character = 1
< error
| 1064 42000 …line 1 column 21 near "character = 1" 
> exec ^
: string
| set session charset = 1
< ok 0 0
> exec ^
: string
| set session check = 1
< error
| 1064 42000 …line 1 column 17 near "check = 1" 
> exec ^
: string
| set session checksum = 1
< ok 0 0
> exec ^
: string
| set session cipher = 1
< ok 0 0
> exec ^
: string
| set session class_origin = 1
< ok 0 0
> exec ^
: string
| set session client = 1
< ok 0 0
> exec ^
: string
| set session clone = 1
< ok 0 0
> exec ^
: string
| set session close = 1
< ok 0 0
> exec ^
: string
| set session clustered = 1
< ok 0 0
> exec ^
: string
| set session coalesce = 1
< ok 0 0
> exec ^
: string
| set session code = 1
< ok 0 0
> exec ^
: string
| set session collate = 1
< error
| 1064 42000 …line 1 column 19 near "collate = 1" 
> exec ^
: string
| set session collation = 1
< ok 0 0
> exec ^
: string
| set session column = 1
< error
| 1064 42000 …line 1 column 18 near "column = 1" 
> exec ^
: string
| set session column_format = 1
< ok 0 0
> exec ^
: string
| set session column_name = 1
< ok 0 0
> exec ^
: string
| set session column_stats_usage = 1
< ok 0 0
> exec ^
: string
| set session columns = 1
< ok 0 0
> exec ^
: string
| set session comment = 1
< ok 0 0
> exec ^
: string
| set session commit = 1
< ok 0 0
> exec ^
: string
| set session committed = 1
< ok 0 0
> exec ^
: string
| set session compact = 1
< ok 0 0
> exec ^
: string
| set session completion = 1
< ok 0 0
> exec ^
: string
| set session component = 1
< ok 0 0
> exec ^
: string
| set session compressed = 1
< ok 0 0
> exec ^
: string
| set session compression = 1
< ok 0 0
> exec ^
: string
| set session concat = 1
< ok 0 0
> exec ^
: string
| set session concurrent = 1
< ok 0 0
> exec ^
: string
| set session condition = 1
< ok 0 0
> exec ^
: string
| set session config = 1
< ok 0 0
> exec ^
: string
| set session connection = 1
< ok 0 0
> exec ^
: string
| set session consistent = 1
< ok 0 0
> exec ^
: string
| set session constraint = 1
< error
| 1064 42000 …line 1 column 22 near "constraint = 1" 
> exec ^
: string
| set session constraint_catalog = 1
< ok 0 0
> exec ^
: string
| set session constraint_name = 1
< ok 0 0
> exec ^
: string
| set session constraint_schema = 1
< ok 0 0
> exec ^
: string
| set session constraints = 1
< ok 0 0
> exec ^
: string
| set session contains = 1
< ok 0 0
> exec ^
: string
| set session context = 1
< ok 0 0
> exec ^
: string
| set session continue = 1
< error
| 1064 42000 …line 1 column 20 near "continue = 1" 
> exec ^
: string
| set session convert = 1
< error
| 1064 42000 …line 1 column 19 near "convert = 1" 
> exec ^
: string
| set session cooldown = 1
< ok 0 0
> exec ^
: string
| set session count = 1
< ok 0 0
> exec ^
: string
| set session cpu = 1
< ok 0 0
> exec ^
: string
| set session create = 1
< error
| 1064 42000 …line 1 column 18 near "create = 1" 
> exec ^
: string
| set session cross = 1
< error
| 1064 42000 …line 1 column 17 near "cross = 1" 
> exec ^
: string
| set session cube = 1
< ok 0 0
> exec ^
: string
| set session cume_dist = 1
< error
| 1064 42000 …line 1 column 21 near "cume_dist = 1" 
> exec ^
: string
| set session curdate = 1
< ok 0 0
> exec ^
: string
| set session current = 1
< ok 0 0
> exec ^
: string
| set session current_date = 1
< error
| 1064 42000 …line 1 column 24 near "current_date = 1" 
> exec ^
: string
| set session current_time = 1
< error
| 1064 42000 …line 1 column 24 near "current_time = 1" 
> exec ^
: string
| set session current_timestamp = 1
< error
| 1064 42000 …line 1 column 29 near "current_timestamp = 1" 
> exec ^
: string
| set session current_user = 1
< error
| 1064 42000 …line 1 column 24 near "current_user = 1" 
> exec ^
: string
| set session cursor = 1
< error
| 1064 42000 …line 1 column 18 near "cursor = 1" 
> exec ^
: string
| set session cursor_name = 1
< ok 0 0
> exec ^
: string
| set session curtime = 1
< ok 0 0
> exec ^
: string
| set session cycle = 1
< ok 0 0
> exec ^
: string
| set session data = 1
< ok 0 0
> exec ^
: string
| set session database = 1
< error
| 1064 42000 …line 1 column 20 near "database = 1" 
> exec ^
: string
| set session databases = 1
< error
| 1064 42000 …line 1 column 21 near "databases = 1" 
> exec ^
: string
| set session datafile = 1
< ok 0 0
> exec ^
: string
| set session date = 1
< ok 0 0
> exec ^
: string
| set session date_add = 1
< ok 0 0
> exec ^
: string
| set session date_sub = 1
< ok 0 0
> exec ^
: string
| set session datetime = 1
< ok 0 0
> exec ^
: string
| set session day = 1
< ok 0 0
> exec ^
: string
| set session day_hour = 1
< error
| 1064 42000 …line 1 column 20 near "day_hour = 1" 
> exec ^
: string
| set session day_microsecond = 1
< error
| 1064 42000 …line 1 column 27 near "day_microsecond = 1" 
> exec ^
: string
| set session day_minute = 1
< error
| 1064 42000 …line 1 column 22 near "day_minute = 1" 
> exec ^
: string
| set session day_second = 1
< error
| 1064 42000 …line 1 column 22 near "day_second = 1" 
> exec ^
: string
| set session deallocate = 1
< ok 0 0
> exec ^
: string
| set session dec = 1
< error
| 1064 42000 …line 1 column 15 near "dec = 1" 
> exec ^
: string
| set session decimal = 1
< error
| 1064 42000 …line 1 column 19 near "decimal = 1" 
> exec ^
: string
| set session declare = 1
< ok 0 0
> exec ^
: string
| set session default = 1
< error
| 1064 42000 …line 1 column 19 near "default = 1" 
> exec ^
: string
| set session default_auth = 1
< ok 0 0
> exec ^
: string
| set session defined = 1
< ok 0 0
> exec ^
: string
| set session definer = 1
< ok 0 0
> exec ^
: string
| set session definition = 1
< ok 0 0
> exec ^
: string
| set session delay_key_write = 1
< ok 0 0
> exec ^
: string
| set session delayed = 1
< error
| 1064 42000 …line 1 column 19 near "delayed = 1" 
> exec ^
: string
| set session delete = 1
< error
| 1064 42000 …line 1 column 18 near "delete = 1" 
> exec ^
: string
| set session dense_rank = 1
< error
| 1064 42000 …line 1 column 22 near "dense_rank = 1" 
> exec ^
: string
| set session desc = 1
< error
| 1064 42000 …line 1 column 16 near "desc = 1" 
> exec ^
: string
| set session describe = 1
< error
| 1064 42000 …line 1 column 20 near "describe = 1" 
> exec ^
: string
| set session description = 1
< ok 0 0
> exec ^
: string
| set session deterministic = 1
< ok 0 0
> exec ^
: string
| set session diagnostics = 1
< ok 0 0
> exec ^
: string
| set session digest = 1
< ok 0 0
> exec ^
: string
| set session directory = 1
< ok 0 0
> exec ^
: string
| set session disable = 1
< ok 0 0
> exec ^
: string
| set session disabled = 1
< ok 0 0
> exec ^
: string
| set session discard = 1
< ok 0 0
> exec ^
: string
| set session disk = 1
< ok 0 0
> exec ^
: string
| set session distinct = 1
< error
| 1064 42000 …line 1 column 20 near "distinct = 1" 
> exec ^
: string
| set session distinctrow = 1
< error
| 1064 42000 …line 1 column 23 near "distinctrow = 1" 
> exec ^
: string
| set session distribute = 1
< ok 0 0
> exec ^
: string
| set session distribution = 1
< ok 0 0
> exec ^
: string
| set session distributions = 1
< ok 0 0
> exec ^
: string
| set session div = 1
< error
| 1064 42000 …line 1 column 15 near "div = 1" 
> exec ^
: string
| set session do = 1
< ok 0 0
> exec ^
: string
| set session double = 1
< error
| 1064 42000 …line 1 column 18 near "double = 1" 
> exec ^
: string
| set session drop = 1
< error
| 1064 42000 …line 1 column 16 near "drop = 1" 
> exec ^
: string
| set session dry = 1
< ok 0 0
> exec ^
: string
| set session dryrun = 1
< ok 0 0
> exec ^
: string
| set session dual = 1
< error
| 1064 42000 …line 1 column 16 near "dual = 1" 
> exec ^
: string
| set session duality = 1
< ok 0 0
> exec ^
: string
| set session dumpfile = 1
< ok 0 0
> exec ^
: string
| set session duplicate = 1
< ok 0 0
> exec ^
: string
| set session duration = 1
< ok 0 0
> exec ^
: string
| set session dynamic = 1
< ok 0 0
> exec ^
: string
| set session each = 1
< ok 0 0
> exec ^
: string
| set session else = 1
< error
| 1064 42000 …line 1 column 16 near "else = 1" 
> exec ^
: string
| set session elseif = 1
< error
| 1064 42000 …line 1 column 18 near "elseif = 1" 
> exec ^
: string
| set session empty = 1
< ok 0 0
> exec ^
: string
| set session enable = 1
< ok 0 0
> exec ^
: string
| set session enabled = 1
< ok 0 0
> exec ^
: string
| set session enclosed = 1
< error
| 1064 42000 …line 1 column 20 near "enclosed = 1" 
> exec ^
: string
| set session encryption = 1
< ok 0 0
> exec ^
: string
| set session end = 1
< ok 0 0
> exec ^
: string
| set session end_time = 1
< ok 0 0
> exec ^
: string
| set session ends = 1
< ok 0 0
> exec ^
: string
| set session enforced = 1
< ok 0 0
> exec ^
: string
| set session engine = 1
< ok 0 0
> exec ^
: string
| set session engine_attribute = 1
< ok 0 0
> exec ^
: string
| set session engines = 1
< ok 0 0
> exec ^
: string
| set session enum = 1
< ok 0 0
> exec ^
: string
| set session error = 1
< ok 0 0
> exec ^
: string
| set session errors = 1
< ok 0 0
> exec ^
: string
| set session escape = 1
< ok 0 0
> exec ^
: string
| set session escaped = 1
< error
| 1064 42000 …line 1 column 19 near "escaped = 1" 
> exec ^
: string
| set session event = 1
< ok 0 0
> exec ^
: string
| set session events = 1
< ok 0 0
> exec ^
: string
| set session every = 1
< ok 0 0
> exec ^
: string
| set session exact = 1
< ok 0 0
> exec ^
: string
| set session except = 1
< error
| 1064 42000 …line 1 column 18 near "except = 1" 
> exec ^
: string
| set session exchange = 1
< ok 0 0
> exec ^
: string
| set session exclude = 1
< ok 0 0
> exec ^
: string
| set session exec_elapsed = 1
< ok 0 0
> exec ^
: string
| set session execute = 1
< ok 0 0
> exec ^
: string
| set session exists = 1
< error
| 1064 42000 …line 1 column 18 near "exists = 1" 
> exec ^
: string
| set session exit = 1
< error
| 1064 42000 …line 1 column 16 near "exit = 1" 
> exec ^
: string
| set session expansion = 1
< ok 0 0
> exec ^
: string
| set session expire = 1
< ok 0 0
> exec ^
: string
| set session explain = 1
< error
| 1064 42000 …line 1 column 19 near "explain = 1" 
> exec ^
: string
| set session export = 1
< ok 0 0
> exec ^
: string
| set session extended = 1
< ok 0 0
> exec ^
: string
| set session extent_size = 1
< ok 0 0
> exec ^
: string
| set session external = 1
< ok 0 0
> exec ^
: string
| set session external_format = 1
< ok 0 0
> exec ^
: string
| set session extract = 1
< ok 0 0
> exec ^
: string
| set session factor = 1
< ok 0 0
> exec ^
: string
| set session failed_login_attempts = 1
< ok 0 0
> exec ^
: string
| set session false = 1
< error
| 1064 42000 …line 1 column 17 near "false = 1" 
> exec ^
: string
| set session fast = 1
< ok 0 0
> exec ^
: string
| set session faults = 1
< ok 0 0
> exec ^
: string
| set session fetch = 1
< error
| 1064 42000 …line 1 column 17 near "fetch = 1" 
> exec ^
: string
| set session fields = 1
< ok 0 0
> exec ^
: string
| set session file = 1
< ok 0 0
> exec ^
: string
| set session file_block_size = 1
< ok 0 0
> exec ^
: string
| set session file_format = 1
< ok 0 0
> exec ^
: string
| set session file_name = 1
< ok 0 0
> exec ^
: string
| set session file_pattern = 1
< ok 0 0
> exec ^
: string
| set session file_prefix = 1
< ok 0 0
> exec ^
: string
| set session files = 1
< ok 0 0
> exec ^
: string
| set session filter = 1
< ok 0 0
> exec ^
: string
| set session finish = 1
< ok 0 0
> exec ^
: string
| set session first = 1
< ok 0 0
> exec ^
: string
| set session first_value = 1
< error
| 1064 42000 …line 1 column 23 near "first_value = 1" 
> exec ^
: string
| set session fixed = 1
< ok 0 0
> exec ^
: string
| set session flashback = 1
< ok 0 0
> exec ^
: string
| set session float = 1
< error
| 1064 42000 …line 1 column 17 near "float = 1" 
> exec ^
: string
| set session float4 = 1
< error
| 1064 42000 …line 1 column 18 near "float4 = 1" 
> exec ^
: string
| set session float8 = 1
< error
| 1064 42000 …line 1 column 18 near "float8 = 1" 
> exec ^
: string
| set session flush = 1
< ok 0 0
> exec ^
: string
| set session follower_constraints = 1
< ok 0 0
> exec ^
: string
| set session followers = 1
< ok 0 0
> exec ^
: string
| set session following = 1
< ok 0 0
> exec ^
: string
| set session follows = 1
< ok 0 0
> exec ^
: string
| set session for = 1
< error
| 1064 42000 …line 1 column 15 near "for = 1" 
> exec ^
: string
| set session force = 1
< error
| 1064 42000 …line 1 column 17 near "force = 1" 
> exec ^
: string
| set session foreign = 1
< error
| 1064 42000 …line 1 column 19 near "foreign = 1" 
> exec ^
: string
| set session format = 1
< ok 0 0
> exec ^
: string
| set session found = 1
< ok 0 0
> exec ^
: string
| set session from = 1
< error
| 1064 42000 …line 1 column 16 near "from = 1" 
> exec ^
: string
| set session full = 1
< ok 0 0
> exec ^
: string
| set session fulltext = 1
< error
| 1064 42000 …line 1 column 20 near "fulltext = 1" 
> exec ^
: string
| set session function = 1
< ok 0 0
> exec ^
: string
| set session general = 1
< ok 0 0
> exec ^
: string
| set session generate = 1
< ok 0 0
> exec ^
: string
| set session generated = 1
< error
| 1064 42000 …line 1 column 21 near "generated = 1" 
> exec ^
: string
| set session geomcollection = 1
< ok 0 0
> exec ^
: string
| set session geometry = 1
< ok 0 0
> exec ^
: string
| set session geometrycollection = 1
< ok 0 0
> exec ^
: string
| set session get = 1
< ok 0 0
> exec ^
: string
| set session get_format = 1
< ok 0 0
> exec ^
: string
| set session get_source_public_key = 1
< ok 0 0
> exec ^
: string
| set session global = 1
< ok 0 0
> exec ^
: string
| set session grant = 1
< error
| 1064 42000 …line 1 column 17 near "grant = 1" 
> exec ^
: string
| set session grants = 1
< ok 0 0
> exec ^
: string
| set session greatest = 1
< ok 0 0
> exec ^
: string
| set session group = 1
< error
| 1064 42000 …line 1 column 17 near "group = 1" 
> exec ^
: string
| set session group_concat = 1
< ok 0 0
> exec ^
: string
| set session group_replication = 1
< ok 0 0
> exec ^
: string
| set session grouping = 1
< ok 0 0
> exec ^
: string
| set session groups = 1
< error
| 1064 42000 …line 1 column 18 near "groups = 1" 
> exec ^
: string
| set session gtid_only = 1
< ok 0 0
> exec ^
: string
| set session gtids = 1
< ok 0 0
> exec ^
: string
| set session guided = 1
< ok 0 0
> exec ^
: string
| set session handler = 1
< ok 0 0
> exec ^
: string
| set session hash = 1
< ok 0 0
> exec ^
: string
| set session having = 1
< error
| 1064 42000 …line 1 column 18 near "having = 1" 
> exec ^
: string
| set session header = 1
< ok 0 0
> exec ^
: string
| set session help = 1
< ok 0 0
> exec ^
: string
| set session high = 1
< ok 0 0
> exec ^
: string
| set session high_priority = 1
< error
| 1064 42000 …line 1 column 25 near "high_priority = 1" 
> exec ^
: string
| set session histogram = 1
< ok 0 0
> exec ^
: string
| set session histograms_in_flight = 1
< ok 0 0
> exec ^
: string
| set session history = 1
< ok 0 0
> exec ^
: string
| set session host = 1
< ok 0 0
> exec ^
: string
| set session hosts = 1
< ok 0 0
> exec ^
: string
| set session hour = 1
< ok 0 0
> exec ^
: string
| set session hour_microsecond = 1
< error
| 1064 42000 …line 1 column 28 near "hour_microsecond = 1" 
> exec ^
: string
| set session hour_minute = 1
< error
| 1064 42000 …line 1 column 23 near "hour_minute = 1" 
> exec ^
: string
| set session hour_second = 1
< error
| 1064 42000 …line 1 column 23 near "hour_second = 1" 
> exec ^
: string
| set session identified = 1
< ok 0 0
> exec ^
: string
| set session if = 1
< error
| 1064 42000 …line 1 column 14 near "if = 1" 
> exec ^
: string
| set session ignore = 1
< error
| 1064 42000 …line 1 column 18 near "ignore = 1" 
> exec ^
: string
| set session ignore_server_ids = 1
< ok 0 0
> exec ^
: string
| set session ilike = 1
< error
| 1064 42000 …line 1 column 17 near "ilike = 1" 
> exec ^
: string
| set session import = 1
< ok 0 0
> exec ^
: string
| set session in = 1
< error
| 1064 42000 …line 1 column 14 near "in = 1" 
> exec ^
: string
| set session inactive = 1
< ok 0 0
> exec ^
: string
| set session increment = 1
< ok 0 0
> exec ^
: string
| set session index = 1
< error
| 1064 42000 …line 1 column 17 near "index = 1" 
> exec ^
: string
| set session indexes = 1
< ok 0 0
> exec ^
: string
| set session infile = 1
< error
| 1064 42000 …line 1 column 18 near "infile = 1" 
> exec ^
: string
| set session initial = 1
< ok 0 0
> exec ^
: string
| set session initial_size = 1
< ok 0 0
> exec ^
: string
| set session initiate = 1
< ok 0 0
> exec ^
: string
| set session inner = 1
< error
| 1064 42000 …line 1 column 17 near "inner = 1" 
> exec ^
: string
| set session inout = 1
< error
| 1064 42000 …line 1 column 17 near "inout = 1" 
> exec ^
: string
| set session insensitive = 1
< ok 0 0
> exec ^
: string
| set session insert = 1
< error
| 1064 42000 …line 1 column 18 near "insert = 1" 
> exec ^
: string
| set session insert_method = 1
< ok 0 0
> exec ^
: string
| set session install = 1
< ok 0 0
> exec ^
: string
| set session instance = 1
< ok 0 0
> exec ^
: string
| set session int = 1
< error
| 1064 42000 …line 1 column 15 near "int = 1" 
> exec ^
: string
| set session int1 = 1
< error
| 1064 42000 …line 1 column 16 near "int1 = 1" 
> exec ^
: string
| set session int2 = 1
< error
| 1064 42000 …line 1 column 16 near "int2 = 1" 
> exec ^
: string
| set session int3 = 1
< error
| 1064 42000 …line 1 column 16 near "int3 = 1" 
> exec ^
: string
| set session int4 = 1
< error
| 1064 42000 …line 1 column 16 near "int4 = 1" 
> exec ^
: string
| set session int8 = 1
< error
| 1064 42000 …line 1 column 16 near "int8 = 1" 
> exec ^
: string
| set session integer = 1
< error
| 1064 42000 …line 1 column 19 near "integer = 1" 
> exec ^
: string
| set session intersect = 1
< error
| 1064 42000 …line 1 column 21 near "intersect = 1" 
> exec ^
: string
| set session interval = 1
< error
| 1064 42000 …line 1 column 20 near "interval = 1" 
> exec ^
: string
| set session into = 1
< error
| 1064 42000 …line 1 column 16 near "into = 1" 
> exec ^
: string
| set session invisible = 1
< ok 0 0
> exec ^
: string
| set session invoker = 1
< ok 0 0
> exec ^
: string
| set session io = 1
< ok 0 0
> exec ^
: string
| set session io_after_gtids = 1
< ok 0 0
> exec ^
: string
| set session io_before_gtids = 1
< ok 0 0
> exec ^
: string
| set session io_thread = 1
< ok 0 0
> exec ^
: string
| set session ipc = 1
< ok 0 0
> exec ^
: string
| set session is = 1
< error
| 1064 42000 …line 1 column 14 near "is = 1" 
> exec ^
: string
| set session isolation = 1
< ok 0 0
> exec ^
: string
| set session issuer = 1
< ok 0 0
> exec ^
: string
| set session iterate = 1
< error
| 1064 42000 …line 1 column 19 near "iterate = 1" 
> exec ^
: string
| set session job = 1
< ok 0 0
> exec ^
: string
| set session jobs = 1
< ok 0 0
> exec ^
: string
| set session join = 1
< error
| 1064 42000 …line 1 column 16 near "join = 1" 
> exec ^
: string
| set session json = 1
< ok 0 0
> exec ^
: string
| set session json_arrayagg = 1
< ok 0 0
> exec ^
: string
| set session json_duality_object = 1
< ok 0 0
> exec ^
: string
| set session json_objectagg = 1
< ok 0 0
> exec ^
: string
| set session json_table = 1
< ok 0 0
> exec ^
: string
| set session json_value = 1
< ok 0 0
> exec ^
: string
| set session key = 1
< error
| 1064 42000 …line 1 column 15 near "key = 1" 
> exec ^
: string
| set session key_block_size = 1
< ok 0 0
> exec ^
: string
| set session keyring = 1
< ok 0 0
> exec ^
: string
| set session keys = 1
< error
| 1064 42000 …line 1 column 16 near "keys = 1" 
> exec ^
: string
| set session kill = 1
< error
| 1064 42000 …line 1 column 16 near "kill = 1" 
> exec ^
: string
| set session labels = 1
< ok 0 0
> exec ^
: string
| set session lag = 1
< error
| 1064 42000 …line 1 column 15 near "lag = 1" 
> exec ^
: string
| set session language = 1
< ok 0 0
> exec ^
: string
| set session last = 1
< ok 0 0
> exec ^
: string
| set session last_value = 1
< error
| 1064 42000 …line 1 column 22 near "last_value = 1" 
> exec ^
: string
| set session lateral = 1
< ok 0 0
> exec ^
: string
| set session lead = 1
< error
| 1064 42000 …line 1 column 16 near "lead = 1" 
> exec ^
: string
| set session leader_constraints = 1
< ok 0 0
> exec ^
: string
| set session leading = 1
< error
| 1064 42000 …line 1 column 19 near "leading = 1" 
> exec ^
: string
| set session learner_constraints = 1
< ok 0 0
> exec ^
: string
| set session learners = 1
< ok 0 0
> exec ^
: string
| set session least = 1
< ok 0 0
> exec ^
: string
| set session leave = 1
< error
| 1064 42000 …line 1 column 17 near "leave = 1" 
> exec ^
: string
| set session leaves = 1
< ok 0 0
> exec ^
: string
| set session left = 1
< error
| 1064 42000 …line 1 column 16 near "left = 1" 
> exec ^
: string
| set session less = 1
< ok 0 0
> exec ^
: string
| set session level = 1
< ok 0 0
> exec ^
: string
| set session library = 1
< ok 0 0
> exec ^
: string
| set session like = 1
< error
| 1064 42000 …line 1 column 16 near "like = 1" 
> exec ^
: string
| set session limit = 1
< error
| 1064 42000 …line 1 column 17 near "limit = 1" 
> exec ^
: string
| set session linear = 1
< error
| 1064 42000 …line 1 column 18 near "linear = 1" 
> exec ^
: string
| set session lines = 1
< error
| 1064 42000 …line 1 column 17 near "lines = 1" 
> exec ^
: string
| set session linestring = 1
< ok 0 0
> exec ^
: string
| set session list = 1
< ok 0 0
> exec ^
: string
| set session load = 1
< error
| 1064 42000 …line 1 column 16 near "load = 1" 
> exec ^
: string
| set session local = 1
< ok 0 0
> exec ^
: string
| set session localtime = 1
< error
| 1064 42000 …line 1 column 21 near "localtime = 1" 
> exec ^
: string
| set session localtimestamp = 1
< error
| 1064 42000 …line 1 column 26 near "localtimestamp = 1" 
> exec ^
: string
| set session location = 1
< ok 0 0
> exec ^
: string
| set session lock = 1
< error
| 1064 42000 …line 1 column 16 near "lock = 1" 
> exec ^
: string
| set session locked = 1
< ok 0 0
> exec ^
: string
| set session locks = 1
< ok 0 0
> exec ^
: string
| set session log = 1
< ok 0 0
> exec ^
: string
| set session logfile = 1
< ok 0 0
> exec ^
: string
| set session logs = 1
< ok 0 0
> exec ^
: string
| set session long = 1
< error
| 1064 42000 …line 1 column 16 near "long = 1" 
> exec ^
: string
| set session longblob = 1
< error
| 1064 42000 …line 1 column 20 near "longblob = 1" 
> exec ^
: string
| set session longtext = 1
< error
| 1064 42000 …line 1 column 20 near "longtext = 1" 
> exec ^
: string
| set session loop = 1
< ok 0 0
> exec ^
: string
| set session low = 1
< ok 0 0
> exec ^
: string
| set session low_priority = 1
< error
| 1064 42000 …line 1 column 24 near "low_priority = 1" 
> exec ^
: string
| set session manual = 1
< ok 0 0
> exec ^
: string
| set session master = 1
< ok 0 0
> exec ^
: string
| set session master_log_file = 1
< ok 0 0
> exec ^
: string
| set session match = 1
< error
| 1064 42000 …line 1 column 17 near "match = 1" 
> exec ^
: string
| set session materialized = 1
< ok 0 0
> exec ^
: string
| set session max = 1
< ok 0 0
> exec ^
: string
| set session max_connections_per_hour = 1
< ok 0 0
> exec ^
: string
| set session max_queries_per_hour = 1
< ok 0 0
> exec ^
: string
| set session max_rows = 1
< ok 0 0
> exec ^
: string
| set session max_size = 1
< ok 0 0
> exec ^
: string
| set session max_updates_per_hour = 1
< ok 0 0
> exec ^
: string
| set session max_user_connections = 1
< ok 0 0
> exec ^
: string
| set session maxvalue = 1
< error
| 1064 42000 …line 1 column 20 near "maxvalue = 1" 
> exec ^
: string
| set session medium = 1
< ok 0 0
> exec ^
: string
| set session mediumblob = 1
< error
| 1064 42000 …line 1 column 22 near "mediumblob = 1" 
> exec ^
: string
| set session mediumint = 1
< error
| 1064 42000 …line 1 column 21 near "mediumint = 1" 
> exec ^
: string
| set session mediumtext = 1
< error
| 1064 42000 …line 1 column 22 near "mediumtext = 1" 
> exec ^
: string
| set session member = 1
< ok 0 0
> exec ^
: string
| set session memory = 1
< ok 0 0
> exec ^
: string
| set session merge = 1
< ok 0 0
> exec ^
: string
| set session message_text = 1
< ok 0 0
> exec ^
: string
| set session microsecond = 1
< ok 0 0
> exec ^
: string
| set session mid = 1
< ok 0 0
> exec ^
: string
| set session middleint = 1
< error
| 1064 42000 …line 1 column 21 near "middleint = 1" 
> exec ^
: string
| set session migrate = 1
< ok 0 0
> exec ^
: string
| set session min = 1
< ok 0 0
> exec ^
: string
| set session min_rows = 1
< ok 0 0
> exec ^
: string
| set session minute = 1
< ok 0 0
> exec ^
: string
| set session minute_microsecond = 1
< error
| 1064 42000 …line 1 column 30 near "minute_microsecond = 1" 
> exec ^
: string
| set session minute_second = 1
< error
| 1064 42000 …line 1 column 25 near "minute_second = 1" 
> exec ^
: string
| set session minvalue = 1
< ok 0 0
> exec ^
: string
| set session mod = 1
< error
| 1064 42000 …line 1 column 15 near "mod = 1" 
> exec ^
: string
| set session mode = 1
< ok 0 0
> exec ^
: string
| set session modifies = 1
< ok 0 0
> exec ^
: string
| set session modify = 1
< ok 0 0
> exec ^
: string
| set session month = 1
< ok 0 0
> exec ^
: string
| set session multilinestring = 1
< ok 0 0
> exec ^
: string
| set session multipoint = 1
< ok 0 0
> exec ^
: string
| set session multipolygon = 1
< ok 0 0
> exec ^
: string
| set session mutex = 1
< ok 0 0
> exec ^
: string
| set session mysql_errno = 1
< ok 0 0
> exec ^
: string
| set session name = 1
< ok 0 0
> exec ^
: string
| set session names = 1
< ok 0 0
> exec ^
: string
| set session national = 1
< ok 0 0
> exec ^
: string
| set session natural = 1
< error
| 1064 42000 …line 1 column 19 near "natural = 1" 
> exec ^
: string
| set session nchar = 1
< ok 0 0
> exec ^
: string
| set session ndb = 1
< ok 0 0
> exec ^
: string
| set session ndbcluster = 1
< ok 0 0
> exec ^
: string
| set session nested = 1
< ok 0 0
> exec ^
: string
| set session network_namespace = 1
< ok 0 0
> exec ^
: string
| set session never = 1
< ok 0 0
> exec ^
: string
| set session new = 1
< ok 0 0
> exec ^
: string
| set session next = 1
< ok 0 0
> exec ^
: string
| set session next_row_id = 1
< ok 0 0
> exec ^
: string
| set session no = 1
< ok 0 0
> exec ^
: string
| set session no_wait = 1
< ok 0 0
> exec ^
: string
| set session no_write_to_binlog = 1
< error
| 1064 42000 …line 1 column 30 near "no_write_to_binlog = 1" 
> exec ^
: string
| set session nocache = 1
< ok 0 0
> exec ^
: string
| set session nocycle = 1
< ok 0 0
> exec ^
: string
| set session nodegroup = 1
< ok 0 0
> exec ^
: string
| set session nomaxvalue = 1
< ok 0 0
> exec ^
: string
| set session nominvalue = 1
< ok 0 0
> exec ^
: string
| set session nonclustered = 1
< ok 0 0
> exec ^
: string
| set session none = 1
< ok 0 0
> exec ^
: string
| set session not = 1
< error
| 1064 42000 …line 1 column 15 near "not = 1" 
> exec ^
: string
| set session now = 1
< ok 0 0
> exec ^
: string
| set session nowait = 1
< ok 0 0
> exec ^
: string
| set session nth_value = 1
< error
| 1064 42000 …line 1 column 21 near "nth_value = 1" 
> exec ^
: string
| set session ntile = 1
< error
| 1064 42000 …line 1 column 17 near "ntile = 1" 
> exec ^
: string
| set session null = 1
< error
| 1064 42000 …line 1 column 16 near "null = 1" 
> exec ^
: string
| set session nullif = 1
< ok 0 0
> exec ^
: string
| set session nulls = 1
< ok 0 0
> exec ^
: string
| set session number = 1
< ok 0 0
> exec ^
: string
| set session numeric = 1
< error
| 1064 42000 …line 1 column 19 near "numeric = 1" 
> exec ^
: string
| set session nvarchar = 1
< ok 0 0
> exec ^
: string
| set session of = 1
< error
| 1064 42000 …line 1 column 14 near "of = 1" 
> exec ^
: string
| set session off = 1
< ok 0 0
> exec ^
: string
| set session offset = 1
< ok 0 0
> exec ^
: string
| set session oj = 1
< ok 0 0
> exec ^
: string
| set session old = 1
< ok 0 0
> exec ^
: string
| set session oltp_read_only = 1
< ok 0 0
> exec ^
: string
| set session oltp_read_write = 1
< ok 0 0
> exec ^
: string
| set session oltp_write_only = 1
< ok 0 0
> exec ^
: string
| set session on = 1
< error
| 1064 42000 …line 1 column 14 near "on = 1" 
> exec ^
: string
| set session one = 1
< ok 0 0
> exec ^
: string
| set session only = 1
< ok 0 0
> exec ^
: string
| set session open = 1
< ok 0 0
> exec ^
: string
| set session optimize = 1
< error
| 1064 42000 …line 1 column 20 near "optimize = 1" 
> exec ^
: string
| set session optimizer_costs = 1
< ok 0 0
> exec ^
: string
| set session option = 1
< error
| 1064 42000 …line 1 column 18 near "option = 1" 
> exec ^
: string
| set session optional = 1
< ok 0 0
> exec ^
: string
| set session optionally = 1
< error
| 1064 42000 …line 1 column 22 near "optionally = 1" 
> exec ^
: string
| set session options = 1
< ok 0 0
> exec ^
: string
| set session or = 1
< error
| 1064 42000 …line 1 column 14 near "or = 1" 
> exec ^
: string
| set session order = 1
< error
| 1064 42000 …line 1 column 17 near "order = 1" 
> exec ^
: string
| set session ordinality = 1
< ok 0 0
> exec ^
: string
| set session organization = 1
< ok 0 0
> exec ^
: string
| set session others = 1
< ok 0 0
> exec ^
: string
| set session out = 1
< error
| 1064 42000 …line 1 column 15 near "out = 1" 
> exec ^
: string
| set session outer = 1
< error
| 1064 42000 …line 1 column 17 near "outer = 1" 
> exec ^
: string
| set session outfile = 1
< error
| 1064 42000 …line 1 column 19 near "outfile = 1" 
> exec ^
: string
| set session over = 1
< error
| 1064 42000 …line 1 column 16 near "over = 1" 
> exec ^
: string
| set session overlay = 1
< ok 0 0
> exec ^
: string
| set session owner = 1
< ok 0 0
> exec ^
: string
| set session pack_keys = 1
< ok 0 0
> exec ^
: string
| set session page = 1
< ok 0 0
> exec ^
: string
| set session parallel = 1
< ok 0 0
> exec ^
: string
| set session parameters = 1
< ok 0 0
> exec ^
: string
| set session parse_tree = 1
< ok 0 0
> exec ^
: string
| set session parser = 1
< ok 0 0
> exec ^
: string
| set session partial = 1
< ok 0 0
> exec ^
: string
| set session partition = 1
< error
| 1064 42000 …line 1 column 21 near "partition = 1" 
> exec ^
: string
| set session partitioning = 1
< ok 0 0
> exec ^
: string
| set session partitions = 1
< ok 0 0
> exec ^
: string
| set session password = 1
< ok 0 0
> exec ^
: string
| set session password_lock_time = 1
< ok 0 0
> exec ^
: string
| set session path = 1
< ok 0 0
> exec ^
: string
| set session percent_rank = 1
< error
| 1064 42000 …line 1 column 24 near "percent_rank = 1" 
> exec ^
: string
| set session persist = 1
< ok 0 0
> exec ^
: string
| set session persist_only = 1
< ok 0 0
> exec ^
: string
| set session phase = 1
< ok 0 0
> exec ^
: string
| set session placement = 1
< ok 0 0
> exec ^
: string
| set session placing = 1
< ok 0 0
> exec ^
: string
| set session plan = 1
< ok 0 0
> exec ^
: string
| set session plugin = 1
< ok 0 0
> exec ^
: string
| set session plugin_dir = 1
< ok 0 0
> exec ^
: string
| set session plugins = 1
< ok 0 0
> exec ^
: string
| set session point = 1
< ok 0 0
> exec ^
: string
| set session policy = 1
< ok 0 0
> exec ^
: string
| set session polygon = 1
< ok 0 0
> exec ^
: string
| set session port = 1
< ok 0 0
> exec ^
: string
| set session position = 1
< ok 0 0
> exec ^
: string
| set session pre_split_regions = 1
< ok 0 0
> exec ^
: string
| set session precedes = 1
< ok 0 0
> exec ^
: string
| set session preceding = 1
< ok 0 0
> exec ^
: string
| set session precision = 1
< error
| 1064 42000 …line 1 column 21 near "precision = 1" 
> exec ^
: string
| set session prepare = 1
< ok 0 0
> exec ^
: string
| set session preserve = 1
< ok 0 0
> exec ^
: string
| set session prev = 1
< ok 0 0
> exec ^
: string
| set session primary = 1
< error
| 1064 42000 …line 1 column 19 near "primary = 1" 
> exec ^
: string
| set session primary_region = 1
< ok 0 0
> exec ^
: string
| set session priority = 1
< ok 0 0
> exec ^
: string
| set session privilege_checks_user = 1
< ok 0 0
> exec ^
: string
| set session privileges = 1
< ok 0 0
> exec ^
: string
| set session procedure = 1
< error
| 1064 42000 …line 1 column 21 near "procedure = 1" 
> exec ^
: string
| set session process = 1
< ok 0 0
> exec ^
: string
| set session processed_keys = 1
< ok 0 0
> exec ^
: string
| set session processlist = 1
< ok 0 0
> exec ^
: string
| set session profile = 1
< ok 0 0
> exec ^
: string
| set session profiles = 1
< ok 0 0
> exec ^
: string
| set session proxy = 1
< ok 0 0
> exec ^
: string
| set session purge = 1
< ok 0 0
> exec ^
: string
| set session qualify = 1
< ok 0 0
> exec ^
: string
| set session quarter = 1
< ok 0 0
> exec ^
: string
| set session query = 1
< ok 0 0
> exec ^
: string
| set session query_limit = 1
< ok 0 0
> exec ^
: string
| set session quick = 1
< ok 0 0
> exec ^
: string
| set session random = 1
< ok 0 0
> exec ^
: string
| set session range = 1
< error
| 1064 42000 …line 1 column 17 near "range = 1" 
> exec ^
: string
| set session rank = 1
< error
| 1064 42000 …line 1 column 16 near "rank = 1" 
> exec ^
: string
| set session read = 1
< error
| 1064 42000 …line 1 column 16 near "read = 1" 
> exec ^
: string
| set session read_only = 1
< ok 0 0
> exec ^
: string
| set session read_write = 1
< ok 0 0
> exec ^
: string
| set session reads = 1
< ok 0 0
> exec ^
: string
| set session real = 1
< error
| 1064 42000 …line 1 column 16 near "real = 1" 
> exec ^
: string
| set session rebuild = 1
< ok 0 0
> exec ^
: string
| set session recover = 1
< ok 0 0
> exec ^
: string
| set session recursive = 1
< error
| 1064 42000 …line 1 column 21 near "recursive = 1" 
> exec ^
: string
| set session redo_buffer_size = 1
< ok 0 0
> exec ^
: string
| set session redundant = 1
< ok 0 0
> exec ^
: string
| set session reference = 1
< ok 0 0
> exec ^
: string
| set session references = 1
< error
| 1064 42000 …line 1 column 22 near "references = 1" 
> exec ^
: string
| set session regexp = 1
< error
| 1064 42000 …line 1 column 18 near "regexp = 1" 
> exec ^
: string
| set session region = 1
< ok 0 0
> exec ^
: string
| set session regions = 1
< ok 0 0
> exec ^
: string
| set session registration = 1
< ok 0 0
> exec ^
: string
| set session relational = 1
< ok 0 0
> exec ^
: string
| set session relay = 1
< ok 0 0
> exec ^
: string
| set session relay_log_file = 1
< ok 0 0
> exec ^
: string
| set session relay_log_pos = 1
< ok 0 0
> exec ^
: string
| set session relay_thread = 1
< ok 0 0
> exec ^
: string
| set session relaylog = 1
< ok 0 0
> exec ^
: string
| set session release = 1
< error
| 1064 42000 …line 1 column 19 near "release = 1" 
> exec ^
: string
| set session reload = 1
< ok 0 0
> exec ^
: string
| set session remove = 1
< ok 0 0
> exec ^
: string
| set session rename = 1
< error
| 1064 42000 …line 1 column 18 near "rename = 1" 
> exec ^
: string
| set session reorganize = 1
< ok 0 0
> exec ^
: string
| set session repair = 1
< ok 0 0
> exec ^
: string
| set session repeat = 1
< error
| 1064 42000 …line 1 column 18 near "repeat = 1" 
> exec ^
: string
| set session repeatable = 1
< ok 0 0
> exec ^
: string
| set session replace = 1
< error
| 1064 42000 …line 1 column 19 near "replace = 1" 
> exec ^
: string
| set session replica = 1
< ok 0 0
> exec ^
: string
| set session replicas = 1
< ok 0 0
> exec ^
: string
| set session replicate_do_db = 1
< ok 0 0
> exec ^
: string
| set session replicate_do_table = 1
< ok 0 0
> exec ^
: string
| set session replicate_ignore_db = 1
< ok 0 0
> exec ^
: string
| set session replicate_ignore_table = 1
< ok 0 0
> exec ^
: string
| set session replicate_rewrite_db = 1
< ok 0 0
> exec ^
: string
| set session replicate_wild_do_table = 1
< ok 0 0
> exec ^
: string
| set session replicate_wild_ignore_table = 1
< ok 0 0
> exec ^
: string
| set session replication = 1
< ok 0 0
> exec ^
: string
| set session require = 1
< error
| 1064 42000 …line 1 column 19 near "require = 1" 
> exec ^
: string
| set session require_row_format = 1
< ok 0 0
> exec ^
: string
| set session require_table_primary_key_check = 1
< ok 0 0
> exec ^
: string
| set session reset = 1
< ok 0 0
> exec ^
: string
| set session resignal = 1
< ok 0 0
> exec ^
: string
| set session resource = 1
< ok 0 0
> exec ^
: string
| set session respect = 1
< ok 0 0
> exec ^
: string
| set session restart = 1
< ok 0 0
> exec ^
: string
| set session restore = 1
< ok 0 0
> exec ^
: string
| set session restores = 1
< ok 0 0
> exec ^
: string
| set session restrict = 1
< error
| 1064 42000 …line 1 column 20 near "restrict = 1" 
> exec ^
: string
| set session resume = 1
< ok 0 0
> exec ^
: string
| set session retain = 1
< ok 0 0
> exec ^
: string
| set session return = 1
< ok 0 0
> exec ^
: string
| set session returned_sqlstate = 1
< ok 0 0
> exec ^
: string
| set session returning = 1
< ok 0 0
> exec ^
: string
| set session returns = 1
< ok 0 0
> exec ^
: string
| set session reuse = 1
< ok 0 0
> exec ^
: string
| set session reverse = 1
< ok 0 0
> exec ^
: string
| set session revoke = 1
< error
| 1064 42000 …line 1 column 18 near "revoke = 1" 
> exec ^
: string
| set session right = 1
< error
| 1064 42000 …line 1 column 17 near "right = 1" 
> exec ^
: string
| set session rlike = 1
< error
| 1064 42000 …line 1 column 17 near "rlike = 1" 
> exec ^
: string
| set session role = 1
< ok 0 0
> exec ^
: string
| set session rollback = 1
< ok 0 0
> exec ^
: string
| set session rollup = 1
< ok 0 0
> exec ^
: string
| set session rotate = 1
< ok 0 0
> exec ^
: string
| set session routine = 1
< ok 0 0
> exec ^
: string
| set session row = 1
< error
| 1064 42000 …line 1 column 15 near "row = 1" 
> exec ^
: string
| set session row_count = 1
< ok 0 0
> exec ^
: string
| set session row_format = 1
< ok 0 0
> exec ^
: string
| set session row_number = 1
< error
| 1064 42000 …line 1 column 22 near "row_number = 1" 
> exec ^
: string
| set session rows = 1
< error
| 1064 42000 …line 1 column 16 near "rows = 1" 
> exec ^
: string
| set session rtree = 1
< ok 0 0
> exec ^
: string
| set session ru = 1
< ok 0 0
> exec ^
: string
| set session ru_per_sec = 1
< ok 0 0
> exec ^
: string
| set session rule = 1
< ok 0 0
> exec ^
: string
| set session run = 1
< ok 0 0
> exec ^
: string
| set session s3 = 1
< ok 0 0
> exec ^
: string
| set session savepoint = 1
< ok 0 0
> exec ^
: string
| set session schedule = 1
< ok 0 0
> exec ^
: string
| set session schema = 1
< error
| 1064 42000 …line 1 column 18 near "schema = 1" 
> exec ^
: string
| set session schema_name = 1
< ok 0 0
> exec ^
: string
| set session schemas = 1
< error
| 1064 42000 …line 1 column 19 near "schemas = 1" 
> exec ^
: string
| set session second = 1
< ok 0 0
> exec ^
: string
| set session second_microsecond = 1
< error
| 1064 42000 …line 1 column 30 near "second_microsecond = 1" 
> exec ^
: string
| set session secondary = 1
< ok 0 0
> exec ^
: string
| set session secondary_engine = 1
< ok 0 0
> exec ^
: string
| set session secondary_engine_attribute = 1
< ok 0 0
> exec ^
: string
| set session secondary_load = 1
< ok 0 0
> exec ^
: string
| set session secondary_unload = 1
< ok 0 0
> exec ^
: string
| set session security = 1
< ok 0 0
> exec ^
: string
| set session select = 1
< error
| 1064 42000 …line 1 column 18 near "select = 1" 
> exec ^
: string
| set session sensitive = 1
< ok 0 0
> exec ^
: string
| set session separator = 1
< ok 0 0
> exec ^
: string
| set session sequence = 1
< ok 0 0
> exec ^
: string
| set session serial = 1
< ok 0 0
> exec ^
: string
| set session serializable = 1
< ok 0 0
> exec ^
: string
| set session server = 1
< ok 0 0
> exec ^
: string
| set session session = 1
< ok 0 0
> exec ^
: string
| set session session_states = 1
< ok 0 0
> exec ^
: string
| set session session_user = 1
< ok 0 0
> exec ^
: string
| set session set = 1
< error
| 1064 42000 …line 1 column 15 near "set = 1" 
> exec ^
: string
| set session sets = 1
< ok 0 0
> exec ^
: string
| set session shard_row_id_bits = 1
< ok 0 0
> exec ^
: string
| set session share = 1
< ok 0 0
> exec ^
: string
| set session show = 1
< error
| 1064 42000 …line 1 column 16 near "show = 1" 
> exec ^
: string
| set session shutdown = 1
< ok 0 0
> exec ^
: string
| set session signal = 1
< ok 0 0
> exec ^
: string
| set session signed = 1
< ok 0 0
> exec ^
: string
| set session similar = 1
< ok 0 0
> exec ^
: string
| set session simple = 1
< ok 0 0
> exec ^
: string
| set session skip = 1
< ok 0 0
> exec ^
: string
| set session slave = 1
< ok 0 0
> exec ^
: string
| set session slow = 1
< ok 0 0
> exec ^
: string
| set session smallint = 1
< error
| 1064 42000 …line 1 column 20 near "smallint = 1" 
> exec ^
: string
| set session snapshot = 1
< ok 0 0
> exec ^
: string
| set session socket = 1
< ok 0 0
> exec ^
: string
| set session some = 1
< ok 0 0
> exec ^
: string
| set session soname = 1
< ok 0 0
> exec ^
: string
| set session sounds = 1
< ok 0 0
> exec ^
: string
| set session source = 1
< ok 0 0
> exec ^
: string
| set session source_auto_position = 1
< ok 0 0
> exec ^
: string
| set session source_bind = 1
< ok 0 0
> exec ^
: string
| set session source_compression_algorithms = 1
< ok 0 0
> exec ^
: string
| set session source_connect_retry = 1
< ok 0 0
> exec ^
: string
| set session source_connection_auto_failover = 1
< ok 0 0
> exec ^
: string
| set session source_delay = 1
< ok 0 0
> exec ^
: string
| set session source_heartbeat_period = 1
< ok 0 0
> exec ^
: string
| set session source_host = 1
< ok 0 0
> exec ^
: string
| set session source_log_file = 1
< ok 0 0
> exec ^
: string
| set session source_log_pos = 1
< ok 0 0
> exec ^
: string
| set session source_password = 1
< ok 0 0
> exec ^
: string
| set session source_port = 1
< ok 0 0
> exec ^
: string
| set session source_public_key_path = 1
< ok 0 0
> exec ^
: string
| set session source_retry_count = 1
< ok 0 0
> exec ^
: string
| set session source_ssl = 1
< ok 0 0
> exec ^
: string
| set session source_ssl_ca = 1
< ok 0 0
> exec ^
: string
| set session source_ssl_capath = 1
< ok 0 0
> exec ^
: string
| set session source_ssl_cert = 1
< ok 0 0
> exec ^
: string
| set session source_ssl_cipher = 1
< ok 0 0
> exec ^
: string
| set session source_ssl_crl = 1
< ok 0 0
> exec ^
: string
| set session source_ssl_crlpath = 1
< ok 0 0
> exec ^
: string
| set session source_ssl_key = 1
< ok 0 0
> exec ^
: string
| set session source_ssl_verify_server_cert = 1
< ok 0 0
> exec ^
: string
| set session source_tls_ciphersuites = 1
< ok 0 0
> exec ^
: string
| set session source_tls_version = 1
< ok 0 0
> exec ^
: string
| set session source_user = 1
< ok 0 0
> exec ^
: string
| set session source_zstd_compression_level = 1
< ok 0 0
> exec ^
: string
| set session spatial = 1
< error
| 1064 42000 …line 1 column 19 near "spatial = 1" 
> exec ^
: string
| set session specific = 1
< ok 0 0
> exec ^
: string
| set session split = 1
< ok 0 0
> exec ^
: string
| set session sql = 1
< error
| 1064 42000 …line 1 column 15 near "sql = 1" 
> exec ^
: string
| set session sql_after_gtids = 1
< ok 0 0
> exec ^
: string
| set session sql_after_mts_gaps = 1
< ok 0 0
> exec ^
: string
| set session sql_before_gtids = 1
< ok 0 0
> exec ^
: string
| set session sql_big_result = 1
< error
| 1064 42000 …line 1 column 26 near "sql_big_result = 1" 
> exec ^
: string
| set session sql_buffer_result = 1
< ok 0 0
> exec ^
: string
| set session sql_calc_found_rows = 1
< error
| 1064 42000 …line 1 column 31 near "sql_calc_found_rows = 1" 
> exec ^
: string
| set session sql_no_cache = 1
< ok 0 0
> exec ^
: string
| set session sql_small_result = 1
< error
| 1064 42000 …line 1 column 28 near "sql_small_result = 1" 
> exec ^
: string
| set session sql_thread = 1
< ok 0 0
> exec ^
: string
| set session sql_tsi_day = 1
< ok 0 0
> exec ^
: string
| set session sql_tsi_hour = 1
< ok 0 0
> exec ^
: string
| set session sql_tsi_minute = 1
< ok 0 0
> exec ^
: string
| set session sql_tsi_month = 1
< ok 0 0
> exec ^
: string
| set session sql_tsi_quarter = 1
< ok 0 0
> exec ^
: string
| set session sql_tsi_second = 1
< ok 0 0
> exec ^
: string
| set session sql_tsi_week = 1
< ok 0 0
> exec ^
: string
| set session sql_tsi_year = 1
< ok 0 0
> exec ^
: string
| set session sqlexception = 1
< error
| 1064 42000 …line 1 column 24 near "sqlexception = 1" 
> exec ^
: string
| set session sqlstate = 1
< error
| 1064 42000 …line 1 column 20 near "sqlstate = 1" 
> exec ^
: string
| set session sqlwarning = 1
< error
| 1064 42000 …line 1 column 22 near "sqlwarning = 1" 
> exec ^
: string
| set session srid = 1
< ok 0 0
> exec ^
: string
| set session ssl = 1
< error
| 1064 42000 …line 1 column 15 near "ssl = 1" 
> exec ^
: string
| set session st_collect = 1
< ok 0 0
> exec ^
: string
| set session stacked = 1
< ok 0 0
> exec ^
: string
| set session start = 1
< ok 0 0
> exec ^
: string
| set session start_time = 1
< ok 0 0
> exec ^
: string
| set session starting = 1
< error
| 1064 42000 …line 1 column 20 near "starting = 1" 
> exec ^
: string
| set session starts = 1
< ok 0 0
> exec ^
: string
| set session stats_auto_recalc = 1
< ok 0 0
> exec ^
: string
| set session stats_buckets = 1
< ok 0 0
> exec ^
: string
| set session stats_col_choice = 1
< ok 0 0
> exec ^
: string
| set session stats_col_list = 1
< ok 0 0
> exec ^
: string
| set session stats_extended = 1
< ok 0 0
> exec ^
: string
| set session stats_healthy = 1
< ok 0 0
> exec ^
: string
| set session stats_histograms = 1
< ok 0 0
> exec ^
: string
| set session stats_locked = 1
< ok 0 0
> exec ^
: string
| set session stats_meta = 1
< ok 0 0
> exec ^
: string
| set session stats_persistent = 1
< ok 0 0
> exec ^
: string
| set session stats_sample_pages = 1
< ok 0 0
> exec ^
: string
| set session stats_sample_rate = 1
< ok 0 0
> exec ^
: string
| set session stats_topn = 1
< ok 0 0
> exec ^
: string
| set session status = 1
< ok 0 0
> exec ^
: string
| set session std = 1
< ok 0 0
> exec ^
: string
| set session stddev = 1
< ok 0 0
> exec ^
: string
| set session stddev_pop = 1
< ok 0 0
> exec ^
: string
| set session stddev_samp = 1
< ok 0 0
> exec ^
: string
| set session stop = 1
< ok 0 0
> exec ^
: string
| set session storage = 1
< ok 0 0
> exec ^
: string
| set session stored = 1
< error
| 1064 42000 …line 1 column 18 near "stored = 1" 
> exec ^
: string
| set session straight_join = 1
< error
| 1064 42000 …line 1 column 25 near "straight_join = 1" 
> exec ^
: string
| set session stream = 1
< ok 0 0
> exec ^
: string
| set session strict_load = 1
< ok 0 0
> exec ^
: string
| set session string = 1
< ok 0 0
> exec ^
: string
| set session subclass_origin = 1
< ok 0 0
> exec ^
: string
| set session subdate = 1
< ok 0 0
> exec ^
: string
| set session subject = 1
< ok 0 0
> exec ^
: string
| set session subpartition = 1
< ok 0 0
> exec ^
: string
| set session subpartitions = 1
< ok 0 0
> exec ^
: string
| set session substr = 1
< ok 0 0
> exec ^
: string
| set session substring = 1
< ok 0 0
> exec ^
: string
| set session sum = 1
< ok 0 0
> exec ^
: string
| set session super = 1
< ok 0 0
> exec ^
: string
| set session survival_preferences = 1
< ok 0 0
> exec ^
: string
| set session suspend = 1
< ok 0 0
> exec ^
: string
| set session swaps = 1
< ok 0 0
> exec ^
: string
| set session switch_group = 1
< ok 0 0
> exec ^
: string
| set session switches = 1
< ok 0 0
> exec ^
: string
| set session sysdate = 1
< ok 0 0
> exec ^
: string
| set session system = 1
< ok 0 0
> exec ^
: string
| set session system_user = 1
< ok 0 0
> exec ^
: string
| set session table = 1
< error
| 1064 42000 …line 1 column 17 near "table = 1" 
> exec ^
: string
| set session table_checksum = 1
< ok 0 0
> exec ^
: string
| set session table_name = 1
< ok 0 0
> exec ^
: string
| set session tables = 1
< ok 0 0
> exec ^
: string
| set session tablesample = 1
< error
| 1064 42000 …line 1 column 23 near "tablesample = 1" 
> exec ^
: string
| set session tablespace = 1
< ok 0 0
> exec ^
: string
| set session task_types = 1
< ok 0 0
> exec ^
: string
| set session temporary = 1
< ok 0 0
> exec ^
: string
| set session temptable = 1
< ok 0 0
> exec ^
: string
| set session terminated = 1
< error
| 1064 42000 …line 1 column 22 near "terminated = 1" 
> exec ^
: string
| set session text = 1
< ok 0 0
> exec ^
: string
| set session than = 1
< ok 0 0
> exec ^
: string
| set session then = 1
< error
| 1064 42000 …line 1 column 16 near "then = 1" 
> exec ^
: string
| set session thread_priority = 1
< ok 0 0
> exec ^
: string
| set session ties = 1
< ok 0 0
> exec ^
: string
| set session tiflash = 1
< ok 0 0
> exec ^
: string
| set session time = 1
< ok 0 0
> exec ^
: string
| set session timeout = 1
< ok 0 0
> exec ^
: string
| set session timestamp = 1
< ok 0 0
> exec ^
: string
| set session timestampadd = 1
< ok 0 0
> exec ^
: string
| set session timestampdiff = 1
< ok 0 0
> exec ^
: string
| set session tinyblob = 1
< error
| 1064 42000 …line 1 column 20 near "tinyblob = 1" 
> exec ^
: string
| set session tinyint = 1
< error
| 1064 42000 …line 1 column 19 near "tinyint = 1" 
> exec ^
: string
| set session tinytext = 1
< error
| 1064 42000 …line 1 column 20 near "tinytext = 1" 
> exec ^
: string
| set session tls = 1
< ok 0 0
> exec ^
: string
| set session to = 1
< error
| 1064 42000 …line 1 column 14 near "to = 1" 
> exec ^
: string
| set session tpcc = 1
< ok 0 0
> exec ^
: string
| set session tpch_10 = 1
< ok 0 0
> exec ^
: string
| set session trailing = 1
< error
| 1064 42000 …line 1 column 20 near "trailing = 1" 
> exec ^
: string
| set session transaction = 1
< ok 0 0
> exec ^
: string
| set session trigger = 1
< error
| 1064 42000 …line 1 column 19 near "trigger = 1" 
> exec ^
: string
| set session triggers = 1
< ok 0 0
> exec ^
: string
| set session trim = 1
< ok 0 0
> exec ^
: string
| set session true = 1
< error
| 1064 42000 …line 1 column 16 near "true = 1" 
> exec ^
: string
| set session truncate = 1
< ok 0 0
> exec ^
: string
| set session ttl = 1
< ok 0 0
> exec ^
: string
| set session ttl_enable = 1
< ok 0 0
> exec ^
: string
| set session ttl_job_interval = 1
< ok 0 0
> exec ^
: string
| set session type = 1
< ok 0 0
> exec ^
: string
| set session types = 1
< ok 0 0
> exec ^
: string
| set session unbounded = 1
< ok 0 0
> exec ^
: string
| set session uncommitted = 1
< ok 0 0
> exec ^
: string
| set session undefined = 1
< ok 0 0
> exec ^
: string
| set session undo = 1
< ok 0 0
> exec ^
: string
| set session undo_buffer_size = 1
< ok 0 0
> exec ^
: string
| set session undofile = 1
< ok 0 0
> exec ^
: string
| set session unicode = 1
< ok 0 0
> exec ^
: string
| set session uninstall = 1
< ok 0 0
> exec ^
: string
| set session union = 1
< error
| 1064 42000 …line 1 column 17 near "union = 1" 
> exec ^
: string
| set session unique = 1
< error
| 1064 42000 …line 1 column 18 near "unique = 1" 
> exec ^
: string
| set session unknown = 1
< ok 0 0
> exec ^
: string
| set session unlimited = 1
< ok 0 0
> exec ^
: string
| set session unlock = 1
< error
| 1064 42000 …line 1 column 18 near "unlock = 1" 
> exec ^
: string
| set session unregister = 1
< ok 0 0
> exec ^
: string
| set session unsigned = 1
< error
| 1064 42000 …line 1 column 20 near "unsigned = 1" 
> exec ^
: string
| set session until = 1
< error
| 1064 42000 …line 1 column 17 near "until = 1" 
> exec ^
: string
| set session update = 1
< error
| 1064 42000 …line 1 column 18 near "update = 1" 
> exec ^
: string
| set session upgrade = 1
< ok 0 0
> exec ^
: string
| set session uri = 1
< ok 0 0
> exec ^
: string
| set session url = 1
< ok 0 0
> exec ^
: string
| set session usage = 1
< error
| 1064 42000 …line 1 column 17 near "usage = 1" 
> exec ^
: string
| set session use = 1
< error
| 1064 42000 …line 1 column 15 near "use = 1" 
> exec ^
: string
| set session use_frm = 1
< ok 0 0
> exec ^
: string
| set session user = 1
< ok 0 0
> exec ^
: string
| set session user_resources = 1
< ok 0 0
> exec ^
: string
| set session using = 1
< error
| 1064 42000 …line 1 column 17 near "using = 1" 
> exec ^
: string
| set session utc_date = 1
< error
| 1064 42000 …line 1 column 20 near "utc_date = 1" 
> exec ^
: string
| set session utc_time = 1
< error
| 1064 42000 …line 1 column 20 near "utc_time = 1" 
> exec ^
: string
| set session utc_timestamp = 1
< error
| 1064 42000 …line 1 column 25 near "utc_timestamp = 1" 
> exec ^
: string
| set session utilization_limit = 1
< ok 0 0
> exec ^
: string
| set session validate = 1
< ok 0 0
> exec ^
: string
| set session validation = 1
< ok 0 0
> exec ^
: string
| set session value = 1
< ok 0 0
> exec ^
: string
| set session values = 1
< error
| 1064 42000 …line 1 column 18 near "values = 1" 
> exec ^
: string
| set session var_pop = 1
< ok 0 0
> exec ^
: string
| set session var_samp = 1
< ok 0 0
> exec ^
: string
| set session varbinary = 1
< error
| 1064 42000 …line 1 column 21 near "varbinary = 1" 
> exec ^
: string
| set session varchar = 1
< error
| 1064 42000 …line 1 column 19 near "varchar = 1" 
> exec ^
: string
| set session varcharacter = 1
< error
| 1064 42000 …line 1 column 24 near "varcharacter = 1" 
> exec ^
: string
| set session variables = 1
< ok 0 0
> exec ^
: string
| set session variance = 1
< ok 0 0
> exec ^
: string
| set session varying = 1
< error
| 1064 42000 …line 1 column 19 near "varying = 1" 
> exec ^
: string
| set session vcpu = 1
< ok 0 0
> exec ^
: string
| set session vector = 1
< ok 0 0
> exec ^
: string
| set session verify_key_constraints = 1
< ok 0 0
> exec ^
: string
| set session view = 1
< ok 0 0
> exec ^
: string
| set session virtual = 1
< error
| 1064 42000 …line 1 column 19 near "virtual = 1" 
> exec ^
: string
| set session visible = 1
< ok 0 0
> exec ^
: string
| set session voter_constraints = 1
< ok 0 0
> exec ^
: string
| set session voters = 1
< ok 0 0
> exec ^
: string
| set session wait = 1
< ok 0 0
> exec ^
: string
| set session warnings = 1
< ok 0 0
> exec ^
: string
| set session watch = 1
< ok 0 0
> exec ^
: string
| set session week = 1
< ok 0 0
> exec ^
: string
| set session weight_string = 1
< ok 0 0
> exec ^
: string
| set session when = 1
< error
| 1064 42000 …line 1 column 16 near "when = 1" 
> exec ^
: string
| set session where = 1
< error
| 1064 42000 …line 1 column 17 near "where = 1" 
> exec ^
: string
| set session while = 1
< error
| 1064 42000 …line 1 column 17 near "while = 1" 
> exec ^
: string
| set session window = 1
< error
| 1064 42000 …line 1 column 18 near "window = 1" 
> exec ^
: string
| set session with = 1
< error
| 1064 42000 …line 1 column 16 near "with = 1" 
> exec ^
: string
| set session without = 1
< ok 0 0
> exec ^
: string
| set session work = 1
< ok 0 0
> exec ^
: string
| set session workload = 1
< ok 0 0
> exec ^
: string
| set session wrapper = 1
< ok 0 0
> exec ^
: string
| set session write = 1
< error
| 1064 42000 …line 1 column 17 near "write = 1" 
> exec ^
: string
| set session x509 = 1
< ok 0 0
> exec ^
: string
| set session xa = 1
< ok 0 0
> exec ^
: string
| set session xid = 1
< ok 0 0
> exec ^
: string
| set session xml = 1
< ok 0 0
> exec ^
: string
| set session xor = 1
< error
| 1064 42000 …line 1 column 15 near "xor = 1" 
> exec ^
: string
| set session year = 1
< ok 0 0
> exec ^
: string
| set session year_month = 1
< error
| 1064 42000 …line 1 column 22 near "year_month = 1" 
> exec ^
: string
| set session zerofill = 1
< error
| 1064 42000 …line 1 column 20 near "zerofill = 1" 
> exec ^
: string
| set session zone = 1
< ok 0 0
> exec ^
: string
| set global absent = 1
< ok 0 0
> exec ^
: string
| set global accessible = 1
< ok 0 0
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
| 1064 42000 …line 1 column 14 near "add = 1" 
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
| 1064 42000 …line 1 column 14 near "all = 1" 
> exec ^
: string
| set global allow_missing_files = 1
< ok 0 0
> exec ^
: string
| set global alter = 1
< error
| 1064 42000 …line 1 column 16 near "alter = 1" 
> exec ^
: string
| set global always = 1
< ok 0 0
> exec ^
: string
| set global analyze = 1
< error
| 1064 42000 …line 1 column 18 near "analyze = 1" 
> exec ^
: string
| set global and = 1
< error
| 1064 42000 …line 1 column 14 near "and = 1" 
> exec ^
: string
| set global any = 1
< ok 0 0
> exec ^
: string
| set global array = 1
< error
| 1064 42000 …line 1 column 16 near "array = 1" 
> exec ^
: string
| set global as = 1
< error
| 1064 42000 …line 1 column 13 near "as = 1" 
> exec ^
: string
| set global asc = 1
< error
| 1064 42000 …line 1 column 14 near "asc = 1" 
> exec ^
: string
| set global ascii = 1
< ok 0 0
> exec ^
: string
| set global asensitive = 1
< ok 0 0
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
| set global auto_id_cache = 1
< ok 0 0
> exec ^
: string
| set global auto_increment = 1
< ok 0 0
> exec ^
: string
| set global auto_random = 1
< ok 0 0
> exec ^
: string
| set global auto_random_base = 1
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
| set global background = 1
< ok 0 0
> exec ^
: string
| set global backup = 1
< ok 0 0
> exec ^
: string
| set global backups = 1
< ok 0 0
> exec ^
: string
| set global batch = 1
< ok 0 0
> exec ^
: string
| set global before = 1
< ok 0 0
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
| 1064 42000 …line 1 column 18 near "between = 1" 
> exec ^
: string
| set global bigint = 1
< error
| 1064 42000 …line 1 column 17 near "bigint = 1" 
> exec ^
: string
| set global binary = 1
< error
| 1064 42000 …line 1 column 17 near "binary = 1" 
> exec ^
: string
| set global binding = 1
< ok 0 0
> exec ^
: string
| set global binding_cache = 1
< ok 0 0
> exec ^
: string
| set global bindings = 1
< ok 0 0
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
| 1064 42000 …line 1 column 15 near "blob = 1" 
> exec ^
: string
| set global block = 1
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
| 1064 42000 …line 1 column 15 near "both = 1" 
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
| set global builtins = 1
< ok 0 0
> exec ^
: string
| set global bulk = 1
< ok 0 0
> exec ^
: string
| set global burstable = 1
< ok 0 0
> exec ^
: string
| set global by = 1
< error
| 1064 42000 …line 1 column 13 near "by = 1" 
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
| set global calibrate = 1
< ok 0 0
> exec ^
: string
| set global call = 1
< error
| 1064 42000 …line 1 column 15 near "call = 1" 
> exec ^
: string
| set global cascade = 1
< error
| 1064 42000 …line 1 column 18 near "cascade = 1" 
> exec ^
: string
| set global cascaded = 1
< ok 0 0
> exec ^
: string
| set global case = 1
< error
| 1064 42000 …line 1 column 15 near "case = 1" 
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
| 1064 42000 …line 1 column 17 near "change = 1" 
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
| 1064 42000 …line 1 column 15 near "char = 1" 
> exec ^
: string
| set global character = 1
< error
| 1064 42000 …line 1 column 20 near "character = 1" 
> exec ^
: string
| set global charset = 1
< ok 0 0
> exec ^
: string
| set global check = 1
< error
| 1064 42000 …line 1 column 16 near "check = 1" 
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
| set global clustered = 1
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
| 1064 42000 …line 1 column 18 near "collate = 1" 
> exec ^
: string
| set global collation = 1
< ok 0 0
> exec ^
: string
| set global column = 1
< error
| 1064 42000 …line 1 column 17 near "column = 1" 
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
| set global column_stats_usage = 1
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
< ok 0 0
> exec ^
: string
| set global config = 1
< ok 0 0
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
| 1064 42000 …line 1 column 21 near "constraint = 1" 
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
| set global constraints = 1
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
| 1064 42000 …line 1 column 19 near "continue = 1" 
> exec ^
: string
| set global convert = 1
< error
| 1064 42000 …line 1 column 18 near "convert = 1" 
> exec ^
: string
| set global cooldown = 1
< ok 0 0
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
| 1064 42000 …line 1 column 17 near "create = 1" 
> exec ^
: string
| set global cross = 1
< error
| 1064 42000 …line 1 column 16 near "cross = 1" 
> exec ^
: string
| set global cube = 1
< ok 0 0
> exec ^
: string
| set global cume_dist = 1
< error
| 1064 42000 …line 1 column 20 near "cume_dist = 1" 
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
| 1064 42000 …line 1 column 23 near "current_date = 1" 
> exec ^
: string
| set global current_time = 1
< error
| 1064 42000 …line 1 column 23 near "current_time = 1" 
> exec ^
: string
| set global current_timestamp = 1
< error
| 1064 42000 …line 1 column 28 near "current_timestamp = 1" 
> exec ^
: string
| set global current_user = 1
< error
| 1064 42000 …line 1 column 23 near "current_user = 1" 
> exec ^
: string
| set global cursor = 1
< error
| 1064 42000 …line 1 column 17 near "cursor = 1" 
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
| set global cycle = 1
< ok 0 0
> exec ^
: string
| set global data = 1
< ok 0 0
> exec ^
: string
| set global database = 1
< error
| 1064 42000 …line 1 column 19 near "database = 1" 
> exec ^
: string
| set global databases = 1
< error
| 1064 42000 …line 1 column 20 near "databases = 1" 
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
| 1064 42000 …line 1 column 19 near "day_hour = 1" 
> exec ^
: string
| set global day_microsecond = 1
< error
| 1064 42000 …line 1 column 26 near "day_microsecond = 1" 
> exec ^
: string
| set global day_minute = 1
< error
| 1064 42000 …line 1 column 21 near "day_minute = 1" 
> exec ^
: string
| set global day_second = 1
< error
| 1064 42000 …line 1 column 21 near "day_second = 1" 
> exec ^
: string
| set global deallocate = 1
< ok 0 0
> exec ^
: string
| set global dec = 1
< error
| 1064 42000 …line 1 column 14 near "dec = 1" 
> exec ^
: string
| set global decimal = 1
< error
| 1064 42000 …line 1 column 18 near "decimal = 1" 
> exec ^
: string
| set global declare = 1
< ok 0 0
> exec ^
: string
| set global default = 1
< error
| 1064 42000 …line 1 column 18 near "default = 1" 
> exec ^
: string
| set global default_auth = 1
< ok 0 0
> exec ^
: string
| set global defined = 1
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
| 1064 42000 …line 1 column 18 near "delayed = 1" 
> exec ^
: string
| set global delete = 1
< error
| 1064 42000 …line 1 column 17 near "delete = 1" 
> exec ^
: string
| set global dense_rank = 1
< error
| 1064 42000 …line 1 column 21 near "dense_rank = 1" 
> exec ^
: string
| set global desc = 1
< error
| 1064 42000 …line 1 column 15 near "desc = 1" 
> exec ^
: string
| set global describe = 1
< error
| 1064 42000 …line 1 column 19 near "describe = 1" 
> exec ^
: string
| set global description = 1
< ok 0 0
> exec ^
: string
| set global deterministic = 1
< ok 0 0
> exec ^
: string
| set global diagnostics = 1
< ok 0 0
> exec ^
: string
| set global digest = 1
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
| set global disabled = 1
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
| 1064 42000 …line 1 column 19 near "distinct = 1" 
> exec ^
: string
| set global distinctrow = 1
< error
| 1064 42000 …line 1 column 22 near "distinctrow = 1" 
> exec ^
: string
| set global distribute = 1
< ok 0 0
> exec ^
: string
| set global distribution = 1
< ok 0 0
> exec ^
: string
| set global distributions = 1
< ok 0 0
> exec ^
: string
| set global div = 1
< error
| 1064 42000 …line 1 column 14 near "div = 1" 
> exec ^
: string
| set global do = 1
< ok 0 0
> exec ^
: string
| set global double = 1
< error
| 1064 42000 …line 1 column 17 near "double = 1" 
> exec ^
: string
| set global drop = 1
< error
| 1064 42000 …line 1 column 15 near "drop = 1" 
> exec ^
: string
| set global dry = 1
< ok 0 0
> exec ^
: string
| set global dryrun = 1
< ok 0 0
> exec ^
: string
| set global dual = 1
< error
| 1064 42000 …line 1 column 15 near "dual = 1" 
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
| set global duration = 1
< ok 0 0
> exec ^
: string
| set global dynamic = 1
< ok 0 0
> exec ^
: string
| set global each = 1
< ok 0 0
> exec ^
: string
| set global else = 1
< error
| 1064 42000 …line 1 column 15 near "else = 1" 
> exec ^
: string
| set global elseif = 1
< error
| 1064 42000 …line 1 column 17 near "elseif = 1" 
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
| set global enabled = 1
< ok 0 0
> exec ^
: string
| set global enclosed = 1
< error
| 1064 42000 …line 1 column 19 near "enclosed = 1" 
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
| set global end_time = 1
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
| 1064 42000 …line 1 column 18 near "escaped = 1" 
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
| set global exact = 1
< ok 0 0
> exec ^
: string
| set global except = 1
< error
| 1064 42000 …line 1 column 17 near "except = 1" 
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
| set global exec_elapsed = 1
< ok 0 0
> exec ^
: string
| set global execute = 1
< ok 0 0
> exec ^
: string
| set global exists = 1
< error
| 1064 42000 …line 1 column 17 near "exists = 1" 
> exec ^
: string
| set global exit = 1
< error
| 1064 42000 …line 1 column 15 near "exit = 1" 
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
| 1064 42000 …line 1 column 18 near "explain = 1" 
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
| 1064 42000 …line 1 column 16 near "false = 1" 
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
| 1064 42000 …line 1 column 16 near "fetch = 1" 
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
< error
| 1064 42000 …line 1 column 22 near "first_value = 1" 
> exec ^
: string
| set global fixed = 1
< ok 0 0
> exec ^
: string
| set global flashback = 1
< ok 0 0
> exec ^
: string
| set global float = 1
< error
| 1064 42000 …line 1 column 16 near "float = 1" 
> exec ^
: string
| set global float4 = 1
< error
| 1064 42000 …line 1 column 17 near "float4 = 1" 
> exec ^
: string
| set global float8 = 1
< error
| 1064 42000 …line 1 column 17 near "float8 = 1" 
> exec ^
: string
| set global flush = 1
< ok 0 0
> exec ^
: string
| set global follower_constraints = 1
< ok 0 0
> exec ^
: string
| set global followers = 1
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
| 1064 42000 …line 1 column 14 near "for = 1" 
> exec ^
: string
| set global force = 1
< error
| 1064 42000 …line 1 column 16 near "force = 1" 
> exec ^
: string
| set global foreign = 1
< error
| 1064 42000 …line 1 column 18 near "foreign = 1" 
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
| 1064 42000 …line 1 column 15 near "from = 1" 
> exec ^
: string
| set global full = 1
< ok 0 0
> exec ^
: string
| set global fulltext = 1
< error
| 1064 42000 …line 1 column 19 near "fulltext = 1" 
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
| 1064 42000 …line 1 column 20 near "generated = 1" 
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
< ok 0 0
> exec ^
: string
| set global get_format = 1
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
| 1064 42000 …line 1 column 16 near "grant = 1" 
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
| 1064 42000 …line 1 column 16 near "group = 1" 
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
< error
| 1064 42000 …line 1 column 17 near "groups = 1" 
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
| 1064 42000 …line 1 column 17 near "having = 1" 
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
| set global high = 1
< ok 0 0
> exec ^
: string
| set global high_priority = 1
< error
| 1064 42000 …line 1 column 24 near "high_priority = 1" 
> exec ^
: string
| set global histogram = 1
< ok 0 0
> exec ^
: string
| set global histograms_in_flight = 1
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
| 1064 42000 …line 1 column 27 near "hour_microsecond = 1" 
> exec ^
: string
| set global hour_minute = 1
< error
| 1064 42000 …line 1 column 22 near "hour_minute = 1" 
> exec ^
: string
| set global hour_second = 1
< error
| 1064 42000 …line 1 column 22 near "hour_second = 1" 
> exec ^
: string
| set global identified = 1
< ok 0 0
> exec ^
: string
| set global if = 1
< error
| 1064 42000 …line 1 column 13 near "if = 1" 
> exec ^
: string
| set global ignore = 1
< error
| 1064 42000 …line 1 column 17 near "ignore = 1" 
> exec ^
: string
| set global ignore_server_ids = 1
< ok 0 0
> exec ^
: string
| set global ilike = 1
< error
| 1064 42000 …line 1 column 16 near "ilike = 1" 
> exec ^
: string
| set global import = 1
< ok 0 0
> exec ^
: string
| set global in = 1
< error
| 1064 42000 …line 1 column 13 near "in = 1" 
> exec ^
: string
| set global inactive = 1
< ok 0 0
> exec ^
: string
| set global increment = 1
< ok 0 0
> exec ^
: string
| set global index = 1
< error
| 1064 42000 …line 1 column 16 near "index = 1" 
> exec ^
: string
| set global indexes = 1
< ok 0 0
> exec ^
: string
| set global infile = 1
< error
| 1064 42000 …line 1 column 17 near "infile = 1" 
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
| 1064 42000 …line 1 column 16 near "inner = 1" 
> exec ^
: string
| set global inout = 1
< error
| 1064 42000 …line 1 column 16 near "inout = 1" 
> exec ^
: string
| set global insensitive = 1
< ok 0 0
> exec ^
: string
| set global insert = 1
< error
| 1064 42000 …line 1 column 17 near "insert = 1" 
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
| 1064 42000 …line 1 column 14 near "int = 1" 
> exec ^
: string
| set global int1 = 1
< error
| 1064 42000 …line 1 column 15 near "int1 = 1" 
> exec ^
: string
| set global int2 = 1
< error
| 1064 42000 …line 1 column 15 near "int2 = 1" 
> exec ^
: string
| set global int3 = 1
< error
| 1064 42000 …line 1 column 15 near "int3 = 1" 
> exec ^
: string
| set global int4 = 1
< error
| 1064 42000 …line 1 column 15 near "int4 = 1" 
> exec ^
: string
| set global int8 = 1
< error
| 1064 42000 …line 1 column 15 near "int8 = 1" 
> exec ^
: string
| set global integer = 1
< error
| 1064 42000 …line 1 column 18 near "integer = 1" 
> exec ^
: string
| set global intersect = 1
< error
| 1064 42000 …line 1 column 20 near "intersect = 1" 
> exec ^
: string
| set global interval = 1
< error
| 1064 42000 …line 1 column 19 near "interval = 1" 
> exec ^
: string
| set global into = 1
< error
| 1064 42000 …line 1 column 15 near "into = 1" 
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
< ok 0 0
> exec ^
: string
| set global io_before_gtids = 1
< ok 0 0
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
| 1064 42000 …line 1 column 13 near "is = 1" 
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
| 1064 42000 …line 1 column 18 near "iterate = 1" 
> exec ^
: string
| set global job = 1
< ok 0 0
> exec ^
: string
| set global jobs = 1
< ok 0 0
> exec ^
: string
| set global join = 1
< error
| 1064 42000 …line 1 column 15 near "join = 1" 
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
| 1064 42000 …line 1 column 14 near "key = 1" 
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
| 1064 42000 …line 1 column 15 near "keys = 1" 
> exec ^
: string
| set global kill = 1
< error
| 1064 42000 …line 1 column 15 near "kill = 1" 
> exec ^
: string
| set global labels = 1
< ok 0 0
> exec ^
: string
| set global lag = 1
< error
| 1064 42000 …line 1 column 14 near "lag = 1" 
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
< error
| 1064 42000 …line 1 column 21 near "last_value = 1" 
> exec ^
: string
| set global lateral = 1
< ok 0 0
> exec ^
: string
| set global lead = 1
< error
| 1064 42000 …line 1 column 15 near "lead = 1" 
> exec ^
: string
| set global leader_constraints = 1
< ok 0 0
> exec ^
: string
| set global leading = 1
< error
| 1064 42000 …line 1 column 18 near "leading = 1" 
> exec ^
: string
| set global learner_constraints = 1
< ok 0 0
> exec ^
: string
| set global learners = 1
< ok 0 0
> exec ^
: string
| set global least = 1
< ok 0 0
> exec ^
: string
| set global leave = 1
< error
| 1064 42000 …line 1 column 16 near "leave = 1" 
> exec ^
: string
| set global leaves = 1
< ok 0 0
> exec ^
: string
| set global left = 1
< error
| 1064 42000 …line 1 column 15 near "left = 1" 
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
| 1064 42000 …line 1 column 15 near "like = 1" 
> exec ^
: string
| set global limit = 1
< error
| 1064 42000 …line 1 column 16 near "limit = 1" 
> exec ^
: string
| set global linear = 1
< error
| 1064 42000 …line 1 column 17 near "linear = 1" 
> exec ^
: string
| set global lines = 1
< error
| 1064 42000 …line 1 column 16 near "lines = 1" 
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
| 1064 42000 …line 1 column 15 near "load = 1" 
> exec ^
: string
| set global local = 1
< ok 0 0
> exec ^
: string
| set global localtime = 1
< error
| 1064 42000 …line 1 column 20 near "localtime = 1" 
> exec ^
: string
| set global localtimestamp = 1
< error
| 1064 42000 …line 1 column 25 near "localtimestamp = 1" 
> exec ^
: string
| set global location = 1
< ok 0 0
> exec ^
: string
| set global lock = 1
< error
| 1064 42000 …line 1 column 15 near "lock = 1" 
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
| 1064 42000 …line 1 column 15 near "long = 1" 
> exec ^
: string
| set global longblob = 1
< error
| 1064 42000 …line 1 column 19 near "longblob = 1" 
> exec ^
: string
| set global longtext = 1
< error
| 1064 42000 …line 1 column 19 near "longtext = 1" 
> exec ^
: string
| set global loop = 1
< ok 0 0
> exec ^
: string
| set global low = 1
< ok 0 0
> exec ^
: string
| set global low_priority = 1
< error
| 1064 42000 …line 1 column 23 near "low_priority = 1" 
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
| set global master_log_file = 1
< ok 0 0
> exec ^
: string
| set global match = 1
< error
| 1064 42000 …line 1 column 16 near "match = 1" 
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
< error
| 1064 42000 …line 1 column 19 near "maxvalue = 1" 
> exec ^
: string
| set global medium = 1
< ok 0 0
> exec ^
: string
| set global mediumblob = 1
< error
| 1064 42000 …line 1 column 21 near "mediumblob = 1" 
> exec ^
: string
| set global mediumint = 1
< error
| 1064 42000 …line 1 column 20 near "mediumint = 1" 
> exec ^
: string
| set global mediumtext = 1
< error
| 1064 42000 …line 1 column 21 near "mediumtext = 1" 
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
| 1064 42000 …line 1 column 20 near "middleint = 1" 
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
| 1064 42000 …line 1 column 29 near "minute_microsecond = 1" 
> exec ^
: string
| set global minute_second = 1
< error
| 1064 42000 …line 1 column 24 near "minute_second = 1" 
> exec ^
: string
| set global minvalue = 1
< ok 0 0
> exec ^
: string
| set global mod = 1
< error
| 1064 42000 …line 1 column 14 near "mod = 1" 
> exec ^
: string
| set global mode = 1
< ok 0 0
> exec ^
: string
| set global modifies = 1
< ok 0 0
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
| 1064 42000 …line 1 column 18 near "natural = 1" 
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
| set global next_row_id = 1
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
| 1064 42000 …line 1 column 29 near "no_write_to_binlog = 1" 
> exec ^
: string
| set global nocache = 1
< ok 0 0
> exec ^
: string
| set global nocycle = 1
< ok 0 0
> exec ^
: string
| set global nodegroup = 1
< ok 0 0
> exec ^
: string
| set global nomaxvalue = 1
< ok 0 0
> exec ^
: string
| set global nominvalue = 1
< ok 0 0
> exec ^
: string
| set global nonclustered = 1
< ok 0 0
> exec ^
: string
| set global none = 1
< ok 0 0
> exec ^
: string
| set global not = 1
< error
| 1064 42000 …line 1 column 14 near "not = 1" 
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
< error
| 1064 42000 …line 1 column 20 near "nth_value = 1" 
> exec ^
: string
| set global ntile = 1
< error
| 1064 42000 …line 1 column 16 near "ntile = 1" 
> exec ^
: string
| set global null = 1
< error
| 1064 42000 …line 1 column 15 near "null = 1" 
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
| 1064 42000 …line 1 column 18 near "numeric = 1" 
> exec ^
: string
| set global nvarchar = 1
< ok 0 0
> exec ^
: string
| set global of = 1
< error
| 1064 42000 …line 1 column 13 near "of = 1" 
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
| set global oltp_read_only = 1
< ok 0 0
> exec ^
: string
| set global oltp_read_write = 1
< ok 0 0
> exec ^
: string
| set global oltp_write_only = 1
< ok 0 0
> exec ^
: string
| set global on = 1
< error
| 1064 42000 …line 1 column 13 near "on = 1" 
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
| 1064 42000 …line 1 column 19 near "optimize = 1" 
> exec ^
: string
| set global optimizer_costs = 1
< ok 0 0
> exec ^
: string
| set global option = 1
< error
| 1064 42000 …line 1 column 17 near "option = 1" 
> exec ^
: string
| set global optional = 1
< ok 0 0
> exec ^
: string
| set global optionally = 1
< error
| 1064 42000 …line 1 column 21 near "optionally = 1" 
> exec ^
: string
| set global options = 1
< ok 0 0
> exec ^
: string
| set global or = 1
< error
| 1064 42000 …line 1 column 13 near "or = 1" 
> exec ^
: string
| set global order = 1
< error
| 1064 42000 …line 1 column 16 near "order = 1" 
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
| 1064 42000 …line 1 column 14 near "out = 1" 
> exec ^
: string
| set global outer = 1
< error
| 1064 42000 …line 1 column 16 near "outer = 1" 
> exec ^
: string
| set global outfile = 1
< error
| 1064 42000 …line 1 column 18 near "outfile = 1" 
> exec ^
: string
| set global over = 1
< error
| 1064 42000 …line 1 column 15 near "over = 1" 
> exec ^
: string
| set global overlay = 1
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
| 1064 42000 …line 1 column 20 near "partition = 1" 
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
| set global percent_rank = 1
< error
| 1064 42000 …line 1 column 23 near "percent_rank = 1" 
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
| set global placement = 1
< ok 0 0
> exec ^
: string
| set global placing = 1
< ok 0 0
> exec ^
: string
| set global plan = 1
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
| set global policy = 1
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
| set global pre_split_regions = 1
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
| 1064 42000 …line 1 column 20 near "precision = 1" 
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
| 1064 42000 …line 1 column 18 near "primary = 1" 
> exec ^
: string
| set global primary_region = 1
< ok 0 0
> exec ^
: string
| set global priority = 1
< ok 0 0
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
| 1064 42000 …line 1 column 20 near "procedure = 1" 
> exec ^
: string
| set global process = 1
< ok 0 0
> exec ^
: string
| set global processed_keys = 1
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
< ok 0 0
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
| set global query_limit = 1
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
< error
| 1064 42000 …line 1 column 16 near "range = 1" 
> exec ^
: string
| set global rank = 1
< error
| 1064 42000 …line 1 column 15 near "rank = 1" 
> exec ^
: string
| set global read = 1
< error
| 1064 42000 …line 1 column 15 near "read = 1" 
> exec ^
: string
| set global read_only = 1
< ok 0 0
> exec ^
: string
| set global read_write = 1
< ok 0 0
> exec ^
: string
| set global reads = 1
< ok 0 0
> exec ^
: string
| set global real = 1
< error
| 1064 42000 …line 1 column 15 near "real = 1" 
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
< error
| 1064 42000 …line 1 column 20 near "recursive = 1" 
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
| 1064 42000 …line 1 column 21 near "references = 1" 
> exec ^
: string
| set global regexp = 1
< error
| 1064 42000 …line 1 column 17 near "regexp = 1" 
> exec ^
: string
| set global region = 1
< ok 0 0
> exec ^
: string
| set global regions = 1
< ok 0 0
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
| 1064 42000 …line 1 column 18 near "release = 1" 
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
| 1064 42000 …line 1 column 17 near "rename = 1" 
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
| 1064 42000 …line 1 column 17 near "repeat = 1" 
> exec ^
: string
| set global repeatable = 1
< ok 0 0
> exec ^
: string
| set global replace = 1
< error
| 1064 42000 …line 1 column 18 near "replace = 1" 
> exec ^
: string
| set global replica = 1
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
| 1064 42000 …line 1 column 18 near "require = 1" 
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
< ok 0 0
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
| set global restores = 1
< ok 0 0
> exec ^
: string
| set global restrict = 1
< error
| 1064 42000 …line 1 column 19 near "restrict = 1" 
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
< ok 0 0
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
| 1064 42000 …line 1 column 17 near "revoke = 1" 
> exec ^
: string
| set global right = 1
< error
| 1064 42000 …line 1 column 16 near "right = 1" 
> exec ^
: string
| set global rlike = 1
< error
| 1064 42000 …line 1 column 16 near "rlike = 1" 
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
< error
| 1064 42000 …line 1 column 14 near "row = 1" 
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
< error
| 1064 42000 …line 1 column 21 near "row_number = 1" 
> exec ^
: string
| set global rows = 1
< error
| 1064 42000 …line 1 column 15 near "rows = 1" 
> exec ^
: string
| set global rtree = 1
< ok 0 0
> exec ^
: string
| set global ru = 1
< ok 0 0
> exec ^
: string
| set global ru_per_sec = 1
< ok 0 0
> exec ^
: string
| set global rule = 1
< ok 0 0
> exec ^
: string
| set global run = 1
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
| 1064 42000 …line 1 column 17 near "schema = 1" 
> exec ^
: string
| set global schema_name = 1
< ok 0 0
> exec ^
: string
| set global schemas = 1
< error
| 1064 42000 …line 1 column 18 near "schemas = 1" 
> exec ^
: string
| set global second = 1
< ok 0 0
> exec ^
: string
| set global second_microsecond = 1
< error
| 1064 42000 …line 1 column 29 near "second_microsecond = 1" 
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
| 1064 42000 …line 1 column 17 near "select = 1" 
> exec ^
: string
| set global sensitive = 1
< ok 0 0
> exec ^
: string
| set global separator = 1
< ok 0 0
> exec ^
: string
| set global sequence = 1
< ok 0 0
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
| set global session_states = 1
< ok 0 0
> exec ^
: string
| set global session_user = 1
< ok 0 0
> exec ^
: string
| set global set = 1
< error
| 1064 42000 …line 1 column 14 near "set = 1" 
> exec ^
: string
| set global sets = 1
< ok 0 0
> exec ^
: string
| set global shard_row_id_bits = 1
< ok 0 0
> exec ^
: string
| set global share = 1
< ok 0 0
> exec ^
: string
| set global show = 1
< error
| 1064 42000 …line 1 column 15 near "show = 1" 
> exec ^
: string
| set global shutdown = 1
< ok 0 0
> exec ^
: string
| set global signal = 1
< ok 0 0
> exec ^
: string
| set global signed = 1
< ok 0 0
> exec ^
: string
| set global similar = 1
< ok 0 0
> exec ^
: string
| set global simple = 1
< ok 0 0
> exec ^
: string
| set global skip = 1
< ok 0 0
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
| 1064 42000 …line 1 column 19 near "smallint = 1" 
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
| 1064 42000 …line 1 column 18 near "spatial = 1" 
> exec ^
: string
| set global specific = 1
< ok 0 0
> exec ^
: string
| set global split = 1
< ok 0 0
> exec ^
: string
| set global sql = 1
< error
| 1064 42000 …line 1 column 14 near "sql = 1" 
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
| 1064 42000 …line 1 column 25 near "sql_big_result = 1" 
> exec ^
: string
| set global sql_buffer_result = 1
< ok 0 0
> exec ^
: string
| set global sql_calc_found_rows = 1
< error
| 1064 42000 …line 1 column 30 near "sql_calc_found_rows = 1" 
> exec ^
: string
| set global sql_no_cache = 1
< ok 0 0
> exec ^
: string
| set global sql_small_result = 1
< error
| 1064 42000 …line 1 column 27 near "sql_small_result = 1" 
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
| 1064 42000 …line 1 column 23 near "sqlexception = 1" 
> exec ^
: string
| set global sqlstate = 1
< error
| 1064 42000 …line 1 column 19 near "sqlstate = 1" 
> exec ^
: string
| set global sqlwarning = 1
< error
| 1064 42000 …line 1 column 21 near "sqlwarning = 1" 
> exec ^
: string
| set global srid = 1
< ok 0 0
> exec ^
: string
| set global ssl = 1
< error
| 1064 42000 …line 1 column 14 near "ssl = 1" 
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
| set global start_time = 1
< ok 0 0
> exec ^
: string
| set global starting = 1
< error
| 1064 42000 …line 1 column 19 near "starting = 1" 
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
| set global stats_buckets = 1
< ok 0 0
> exec ^
: string
| set global stats_col_choice = 1
< ok 0 0
> exec ^
: string
| set global stats_col_list = 1
< ok 0 0
> exec ^
: string
| set global stats_extended = 1
< ok 0 0
> exec ^
: string
| set global stats_healthy = 1
< ok 0 0
> exec ^
: string
| set global stats_histograms = 1
< ok 0 0
> exec ^
: string
| set global stats_locked = 1
< ok 0 0
> exec ^
: string
| set global stats_meta = 1
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
| set global stats_sample_rate = 1
< ok 0 0
> exec ^
: string
| set global stats_topn = 1
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
| 1064 42000 …line 1 column 17 near "stored = 1" 
> exec ^
: string
| set global straight_join = 1
< error
| 1064 42000 …line 1 column 24 near "straight_join = 1" 
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
| set global survival_preferences = 1
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
| set global switch_group = 1
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
| 1064 42000 …line 1 column 16 near "table = 1" 
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
< error
| 1064 42000 …line 1 column 22 near "tablesample = 1" 
> exec ^
: string
| set global tablespace = 1
< ok 0 0
> exec ^
: string
| set global task_types = 1
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
| 1064 42000 …line 1 column 21 near "terminated = 1" 
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
| 1064 42000 …line 1 column 15 near "then = 1" 
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
| set global tiflash = 1
< ok 0 0
> exec ^
: string
| set global time = 1
< ok 0 0
> exec ^
: string
| set global timeout = 1
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
| 1064 42000 …line 1 column 19 near "tinyblob = 1" 
> exec ^
: string
| set global tinyint = 1
< error
| 1064 42000 …line 1 column 18 near "tinyint = 1" 
> exec ^
: string
| set global tinytext = 1
< error
| 1064 42000 …line 1 column 19 near "tinytext = 1" 
> exec ^
: string
| set global tls = 1
< ok 0 0
> exec ^
: string
| set global to = 1
< error
| 1064 42000 …line 1 column 13 near "to = 1" 
> exec ^
: string
| set global tpcc = 1
< ok 0 0
> exec ^
: string
| set global tpch_10 = 1
< ok 0 0
> exec ^
: string
| set global trailing = 1
< error
| 1064 42000 …line 1 column 19 near "trailing = 1" 
> exec ^
: string
| set global transaction = 1
< ok 0 0
> exec ^
: string
| set global trigger = 1
< error
| 1064 42000 …line 1 column 18 near "trigger = 1" 
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
| 1064 42000 …line 1 column 15 near "true = 1" 
> exec ^
: string
| set global truncate = 1
< ok 0 0
> exec ^
: string
| set global ttl = 1
< ok 0 0
> exec ^
: string
| set global ttl_enable = 1
< ok 0 0
> exec ^
: string
| set global ttl_job_interval = 1
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
< ok 0 0
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
| 1064 42000 …line 1 column 16 near "union = 1" 
> exec ^
: string
| set global unique = 1
< error
| 1064 42000 …line 1 column 17 near "unique = 1" 
> exec ^
: string
| set global unknown = 1
< ok 0 0
> exec ^
: string
| set global unlimited = 1
< ok 0 0
> exec ^
: string
| set global unlock = 1
< error
| 1064 42000 …line 1 column 17 near "unlock = 1" 
> exec ^
: string
| set global unregister = 1
< ok 0 0
> exec ^
: string
| set global unsigned = 1
< error
| 1064 42000 …line 1 column 19 near "unsigned = 1" 
> exec ^
: string
| set global until = 1
< error
| 1064 42000 …line 1 column 16 near "until = 1" 
> exec ^
: string
| set global update = 1
< error
| 1064 42000 …line 1 column 17 near "update = 1" 
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
| 1064 42000 …line 1 column 16 near "usage = 1" 
> exec ^
: string
| set global use = 1
< error
| 1064 42000 …line 1 column 14 near "use = 1" 
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
| 1064 42000 …line 1 column 16 near "using = 1" 
> exec ^
: string
| set global utc_date = 1
< error
| 1064 42000 …line 1 column 19 near "utc_date = 1" 
> exec ^
: string
| set global utc_time = 1
< error
| 1064 42000 …line 1 column 19 near "utc_time = 1" 
> exec ^
: string
| set global utc_timestamp = 1
< error
| 1064 42000 …line 1 column 24 near "utc_timestamp = 1" 
> exec ^
: string
| set global utilization_limit = 1
< ok 0 0
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
| 1064 42000 …line 1 column 17 near "values = 1" 
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
| 1064 42000 …line 1 column 20 near "varbinary = 1" 
> exec ^
: string
| set global varchar = 1
< error
| 1064 42000 …line 1 column 18 near "varchar = 1" 
> exec ^
: string
| set global varcharacter = 1
< error
| 1064 42000 …line 1 column 23 near "varcharacter = 1" 
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
| 1064 42000 …line 1 column 18 near "varying = 1" 
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
| 1064 42000 …line 1 column 18 near "virtual = 1" 
> exec ^
: string
| set global visible = 1
< ok 0 0
> exec ^
: string
| set global voter_constraints = 1
< ok 0 0
> exec ^
: string
| set global voters = 1
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
| set global watch = 1
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
| 1064 42000 …line 1 column 15 near "when = 1" 
> exec ^
: string
| set global where = 1
< error
| 1064 42000 …line 1 column 16 near "where = 1" 
> exec ^
: string
| set global while = 1
< error
| 1064 42000 …line 1 column 16 near "while = 1" 
> exec ^
: string
| set global window = 1
< error
| 1064 42000 …line 1 column 17 near "window = 1" 
> exec ^
: string
| set global with = 1
< error
| 1064 42000 …line 1 column 15 near "with = 1" 
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
| set global workload = 1
< ok 0 0
> exec ^
: string
| set global wrapper = 1
< ok 0 0
> exec ^
: string
| set global write = 1
< error
| 1064 42000 …line 1 column 16 near "write = 1" 
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
| 1064 42000 …line 1 column 14 near "xor = 1" 
> exec ^
: string
| set global year = 1
< ok 0 0
> exec ^
: string
| set global year_month = 1
< error
| 1064 42000 …line 1 column 21 near "year_month = 1" 
> exec ^
: string
| set global zerofill = 1
< error
| 1064 42000 …line 1 column 19 near "zerofill = 1" 
> exec ^
: string
| set global zone = 1
< ok 0 0
> exec ^
: string
| set local absent = 1
< ok 0 0
> exec ^
: string
| set local accessible = 1
< ok 0 0
> exec ^
: string
| set local account = 1
< ok 0 0
> exec ^
: string
| set local action = 1
< ok 0 0
> exec ^
: string
| set local active = 1
< ok 0 0
> exec ^
: string
| set local add = 1
< error
| 1064 42000 …line 1 column 13 near "add = 1" 
> exec ^
: string
| set local adddate = 1
< ok 0 0
> exec ^
: string
| set local admin = 1
< ok 0 0
> exec ^
: string
| set local after = 1
< ok 0 0
> exec ^
: string
| set local against = 1
< ok 0 0
> exec ^
: string
| set local aggregate = 1
< ok 0 0
> exec ^
: string
| set local algorithm = 1
< ok 0 0
> exec ^
: string
| set local all = 1
< error
| 1064 42000 …line 1 column 13 near "all = 1" 
> exec ^
: string
| set local allow_missing_files = 1
< ok 0 0
> exec ^
: string
| set local alter = 1
< error
| 1064 42000 …line 1 column 15 near "alter = 1" 
> exec ^
: string
| set local always = 1
< ok 0 0
> exec ^
: string
| set local analyze = 1
< error
| 1064 42000 …line 1 column 17 near "analyze = 1" 
> exec ^
: string
| set local and = 1
< error
| 1064 42000 …line 1 column 13 near "and = 1" 
> exec ^
: string
| set local any = 1
< ok 0 0
> exec ^
: string
| set local array = 1
< error
| 1064 42000 …line 1 column 15 near "array = 1" 
> exec ^
: string
| set local as = 1
< error
| 1064 42000 …line 1 column 12 near "as = 1" 
> exec ^
: string
| set local asc = 1
< error
| 1064 42000 …line 1 column 13 near "asc = 1" 
> exec ^
: string
| set local ascii = 1
< ok 0 0
> exec ^
: string
| set local asensitive = 1
< ok 0 0
> exec ^
: string
| set local assign_gtids_to_anonymous_transactions = 1
< ok 0 0
> exec ^
: string
| set local at = 1
< ok 0 0
> exec ^
: string
| set local attribute = 1
< ok 0 0
> exec ^
: string
| set local authentication = 1
< ok 0 0
> exec ^
: string
| set local auto = 1
< ok 0 0
> exec ^
: string
| set local auto_id_cache = 1
< ok 0 0
> exec ^
: string
| set local auto_increment = 1
< ok 0 0
> exec ^
: string
| set local auto_random = 1
< ok 0 0
> exec ^
: string
| set local auto_random_base = 1
< ok 0 0
> exec ^
: string
| set local auto_refresh = 1
< ok 0 0
> exec ^
: string
| set local auto_refresh_source = 1
< ok 0 0
> exec ^
: string
| set local autoextend_size = 1
< ok 0 0
> exec ^
: string
| set local avg = 1
< ok 0 0
> exec ^
: string
| set local avg_row_length = 1
< ok 0 0
> exec ^
: string
| set local background = 1
< ok 0 0
> exec ^
: string
| set local backup = 1
< ok 0 0
> exec ^
: string
| set local backups = 1
< ok 0 0
> exec ^
: string
| set local batch = 1
< ok 0 0
> exec ^
: string
| set local before = 1
< ok 0 0
> exec ^
: string
| set local begin = 1
< ok 0 0
> exec ^
: string
| set local bernoulli = 1
< ok 0 0
> exec ^
: string
| set local between = 1
< error
| 1064 42000 …line 1 column 17 near "between = 1" 
> exec ^
: string
| set local bigint = 1
< error
| 1064 42000 …line 1 column 16 near "bigint = 1" 
> exec ^
: string
| set local binary = 1
< error
| 1064 42000 …line 1 column 16 near "binary = 1" 
> exec ^
: string
| set local binding = 1
< ok 0 0
> exec ^
: string
| set local binding_cache = 1
< ok 0 0
> exec ^
: string
| set local bindings = 1
< ok 0 0
> exec ^
: string
| set local binlog = 1
< ok 0 0
> exec ^
: string
| set local bit = 1
< ok 0 0
> exec ^
: string
| set local bit_and = 1
< ok 0 0
> exec ^
: string
| set local bit_or = 1
< ok 0 0
> exec ^
: string
| set local bit_xor = 1
< ok 0 0
> exec ^
: string
| set local blob = 1
< error
| 1064 42000 …line 1 column 14 near "blob = 1" 
> exec ^
: string
| set local block = 1
< ok 0 0
> exec ^
: string
| set local bool = 1
< ok 0 0
> exec ^
: string
| set local boolean = 1
< ok 0 0
> exec ^
: string
| set local both = 1
< error
| 1064 42000 …line 1 column 14 near "both = 1" 
> exec ^
: string
| set local btree = 1
< ok 0 0
> exec ^
: string
| set local buckets = 1
< ok 0 0
> exec ^
: string
| set local builtins = 1
< ok 0 0
> exec ^
: string
| set local bulk = 1
< ok 0 0
> exec ^
: string
| set local burstable = 1
< ok 0 0
> exec ^
: string
| set local by = 1
< error
| 1064 42000 …line 1 column 12 near "by = 1" 
> exec ^
: string
| set local byte = 1
< ok 0 0
> exec ^
: string
| set local cache = 1
< ok 0 0
> exec ^
: string
| set local calibrate = 1
< ok 0 0
> exec ^
: string
| set local call = 1
< error
| 1064 42000 …line 1 column 14 near "call = 1" 
> exec ^
: string
| set local cascade = 1
< error
| 1064 42000 …line 1 column 17 near "cascade = 1" 
> exec ^
: string
| set local cascaded = 1
< ok 0 0
> exec ^
: string
| set local case = 1
< error
| 1064 42000 …line 1 column 14 near "case = 1" 
> exec ^
: string
| set local cast = 1
< ok 0 0
> exec ^
: string
| set local catalog_name = 1
< ok 0 0
> exec ^
: string
| set local chain = 1
< ok 0 0
> exec ^
: string
| set local challenge_response = 1
< ok 0 0
> exec ^
: string
| set local change = 1
< error
| 1064 42000 …line 1 column 16 near "change = 1" 
> exec ^
: string
| set local changed = 1
< ok 0 0
> exec ^
: string
| set local channel = 1
< ok 0 0
> exec ^
: string
| set local char = 1
< error
| 1064 42000 …line 1 column 14 near "char = 1" 
> exec ^
: string
| set local character = 1
< error
| 1064 42000 …line 1 column 19 near "character = 1" 
> exec ^
: string
| set local charset = 1
< ok 0 0
> exec ^
: string
| set local check = 1
< error
| 1064 42000 …line 1 column 15 near "check = 1" 
> exec ^
: string
| set local checksum = 1
< ok 0 0
> exec ^
: string
| set local cipher = 1
< ok 0 0
> exec ^
: string
| set local class_origin = 1
< ok 0 0
> exec ^
: string
| set local client = 1
< ok 0 0
> exec ^
: string
| set local clone = 1
< ok 0 0
> exec ^
: string
| set local close = 1
< ok 0 0
> exec ^
: string
| set local clustered = 1
< ok 0 0
> exec ^
: string
| set local coalesce = 1
< ok 0 0
> exec ^
: string
| set local code = 1
< ok 0 0
> exec ^
: string
| set local collate = 1
< error
| 1064 42000 …line 1 column 17 near "collate = 1" 
> exec ^
: string
| set local collation = 1
< ok 0 0
> exec ^
: string
| set local column = 1
< error
| 1064 42000 …line 1 column 16 near "column = 1" 
> exec ^
: string
| set local column_format = 1
< ok 0 0
> exec ^
: string
| set local column_name = 1
< ok 0 0
> exec ^
: string
| set local column_stats_usage = 1
< ok 0 0
> exec ^
: string
| set local columns = 1
< ok 0 0
> exec ^
: string
| set local comment = 1
< ok 0 0
> exec ^
: string
| set local commit = 1
< ok 0 0
> exec ^
: string
| set local committed = 1
< ok 0 0
> exec ^
: string
| set local compact = 1
< ok 0 0
> exec ^
: string
| set local completion = 1
< ok 0 0
> exec ^
: string
| set local component = 1
< ok 0 0
> exec ^
: string
| set local compressed = 1
< ok 0 0
> exec ^
: string
| set local compression = 1
< ok 0 0
> exec ^
: string
| set local concat = 1
< ok 0 0
> exec ^
: string
| set local concurrent = 1
< ok 0 0
> exec ^
: string
| set local condition = 1
< ok 0 0
> exec ^
: string
| set local config = 1
< ok 0 0
> exec ^
: string
| set local connection = 1
< ok 0 0
> exec ^
: string
| set local consistent = 1
< ok 0 0
> exec ^
: string
| set local constraint = 1
< error
| 1064 42000 …line 1 column 20 near "constraint = 1" 
> exec ^
: string
| set local constraint_catalog = 1
< ok 0 0
> exec ^
: string
| set local constraint_name = 1
< ok 0 0
> exec ^
: string
| set local constraint_schema = 1
< ok 0 0
> exec ^
: string
| set local constraints = 1
< ok 0 0
> exec ^
: string
| set local contains = 1
< ok 0 0
> exec ^
: string
| set local context = 1
< ok 0 0
> exec ^
: string
| set local continue = 1
< error
| 1064 42000 …line 1 column 18 near "continue = 1" 
> exec ^
: string
| set local convert = 1
< error
| 1064 42000 …line 1 column 17 near "convert = 1" 
> exec ^
: string
| set local cooldown = 1
< ok 0 0
> exec ^
: string
| set local count = 1
< ok 0 0
> exec ^
: string
| set local cpu = 1
< ok 0 0
> exec ^
: string
| set local create = 1
< error
| 1064 42000 …line 1 column 16 near "create = 1" 
> exec ^
: string
| set local cross = 1
< error
| 1064 42000 …line 1 column 15 near "cross = 1" 
> exec ^
: string
| set local cube = 1
< ok 0 0
> exec ^
: string
| set local cume_dist = 1
< error
| 1064 42000 …line 1 column 19 near "cume_dist = 1" 
> exec ^
: string
| set local curdate = 1
< ok 0 0
> exec ^
: string
| set local current = 1
< ok 0 0
> exec ^
: string
| set local current_date = 1
< error
| 1064 42000 …line 1 column 22 near "current_date = 1" 
> exec ^
: string
| set local current_time = 1
< error
| 1064 42000 …line 1 column 22 near "current_time = 1" 
> exec ^
: string
| set local current_timestamp = 1
< error
| 1064 42000 …line 1 column 27 near "current_timestamp = 1" 
> exec ^
: string
| set local current_user = 1
< error
| 1064 42000 …line 1 column 22 near "current_user = 1" 
> exec ^
: string
| set local cursor = 1
< error
| 1064 42000 …line 1 column 16 near "cursor = 1" 
> exec ^
: string
| set local cursor_name = 1
< ok 0 0
> exec ^
: string
| set local curtime = 1
< ok 0 0
> exec ^
: string
| set local cycle = 1
< ok 0 0
> exec ^
: string
| set local data = 1
< ok 0 0
> exec ^
: string
| set local database = 1
< error
| 1064 42000 …line 1 column 18 near "database = 1" 
> exec ^
: string
| set local databases = 1
< error
| 1064 42000 …line 1 column 19 near "databases = 1" 
> exec ^
: string
| set local datafile = 1
< ok 0 0
> exec ^
: string
| set local date = 1
< ok 0 0
> exec ^
: string
| set local date_add = 1
< ok 0 0
> exec ^
: string
| set local date_sub = 1
< ok 0 0
> exec ^
: string
| set local datetime = 1
< ok 0 0
> exec ^
: string
| set local day = 1
< ok 0 0
> exec ^
: string
| set local day_hour = 1
< error
| 1064 42000 …line 1 column 18 near "day_hour = 1" 
> exec ^
: string
| set local day_microsecond = 1
< error
| 1064 42000 …line 1 column 25 near "day_microsecond = 1" 
> exec ^
: string
| set local day_minute = 1
< error
| 1064 42000 …line 1 column 20 near "day_minute = 1" 
> exec ^
: string
| set local day_second = 1
< error
| 1064 42000 …line 1 column 20 near "day_second = 1" 
> exec ^
: string
| set local deallocate = 1
< ok 0 0
> exec ^
: string
| set local dec = 1
< error
| 1064 42000 …line 1 column 13 near "dec = 1" 
> exec ^
: string
| set local decimal = 1
< error
| 1064 42000 …line 1 column 17 near "decimal = 1" 
> exec ^
: string
| set local declare = 1
< ok 0 0
> exec ^
: string
| set local default = 1
< error
| 1064 42000 …line 1 column 17 near "default = 1" 
> exec ^
: string
| set local default_auth = 1
< ok 0 0
> exec ^
: string
| set local defined = 1
< ok 0 0
> exec ^
: string
| set local definer = 1
< ok 0 0
> exec ^
: string
| set local definition = 1
< ok 0 0
> exec ^
: string
| set local delay_key_write = 1
< ok 0 0
> exec ^
: string
| set local delayed = 1
< error
| 1064 42000 …line 1 column 17 near "delayed = 1" 
> exec ^
: string
| set local delete = 1
< error
| 1064 42000 …line 1 column 16 near "delete = 1" 
> exec ^
: string
| set local dense_rank = 1
< error
| 1064 42000 …line 1 column 20 near "dense_rank = 1" 
> exec ^
: string
| set local desc = 1
< error
| 1064 42000 …line 1 column 14 near "desc = 1" 
> exec ^
: string
| set local describe = 1
< error
| 1064 42000 …line 1 column 18 near "describe = 1" 
> exec ^
: string
| set local description = 1
< ok 0 0
> exec ^
: string
| set local deterministic = 1
< ok 0 0
> exec ^
: string
| set local diagnostics = 1
< ok 0 0
> exec ^
: string
| set local digest = 1
< ok 0 0
> exec ^
: string
| set local directory = 1
< ok 0 0
> exec ^
: string
| set local disable = 1
< ok 0 0
> exec ^
: string
| set local disabled = 1
< ok 0 0
> exec ^
: string
| set local discard = 1
< ok 0 0
> exec ^
: string
| set local disk = 1
< ok 0 0
> exec ^
: string
| set local distinct = 1
< error
| 1064 42000 …line 1 column 18 near "distinct = 1" 
> exec ^
: string
| set local distinctrow = 1
< error
| 1064 42000 …line 1 column 21 near "distinctrow = 1" 
> exec ^
: string
| set local distribute = 1
< ok 0 0
> exec ^
: string
| set local distribution = 1
< ok 0 0
> exec ^
: string
| set local distributions = 1
< ok 0 0
> exec ^
: string
| set local div = 1
< error
| 1064 42000 …line 1 column 13 near "div = 1" 
> exec ^
: string
| set local do = 1
< ok 0 0
> exec ^
: string
| set local double = 1
< error
| 1064 42000 …line 1 column 16 near "double = 1" 
> exec ^
: string
| set local drop = 1
< error
| 1064 42000 …line 1 column 14 near "drop = 1" 
> exec ^
: string
| set local dry = 1
< ok 0 0
> exec ^
: string
| set local dryrun = 1
< ok 0 0
> exec ^
: string
| set local dual = 1
< error
| 1064 42000 …line 1 column 14 near "dual = 1" 
> exec ^
: string
| set local duality = 1
< ok 0 0
> exec ^
: string
| set local dumpfile = 1
< ok 0 0
> exec ^
: string
| set local duplicate = 1
< ok 0 0
> exec ^
: string
| set local duration = 1
< ok 0 0
> exec ^
: string
| set local dynamic = 1
< ok 0 0
> exec ^
: string
| set local each = 1
< ok 0 0
> exec ^
: string
| set local else = 1
< error
| 1064 42000 …line 1 column 14 near "else = 1" 
> exec ^
: string
| set local elseif = 1
< error
| 1064 42000 …line 1 column 16 near "elseif = 1" 
> exec ^
: string
| set local empty = 1
< ok 0 0
> exec ^
: string
| set local enable = 1
< ok 0 0
> exec ^
: string
| set local enabled = 1
< ok 0 0
> exec ^
: string
| set local enclosed = 1
< error
| 1064 42000 …line 1 column 18 near "enclosed = 1" 
> exec ^
: string
| set local encryption = 1
< ok 0 0
> exec ^
: string
| set local end = 1
< ok 0 0
> exec ^
: string
| set local end_time = 1
< ok 0 0
> exec ^
: string
| set local ends = 1
< ok 0 0
> exec ^
: string
| set local enforced = 1
< ok 0 0
> exec ^
: string
| set local engine = 1
< ok 0 0
> exec ^
: string
| set local engine_attribute = 1
< ok 0 0
> exec ^
: string
| set local engines = 1
< ok 0 0
> exec ^
: string
| set local enum = 1
< ok 0 0
> exec ^
: string
| set local error = 1
< ok 0 0
> exec ^
: string
| set local errors = 1
< ok 0 0
> exec ^
: string
| set local escape = 1
< ok 0 0
> exec ^
: string
| set local escaped = 1
< error
| 1064 42000 …line 1 column 17 near "escaped = 1" 
> exec ^
: string
| set local event = 1
< ok 0 0
> exec ^
: string
| set local events = 1
< ok 0 0
> exec ^
: string
| set local every = 1
< ok 0 0
> exec ^
: string
| set local exact = 1
< ok 0 0
> exec ^
: string
| set local except = 1
< error
| 1064 42000 …line 1 column 16 near "except = 1" 
> exec ^
: string
| set local exchange = 1
< ok 0 0
> exec ^
: string
| set local exclude = 1
< ok 0 0
> exec ^
: string
| set local exec_elapsed = 1
< ok 0 0
> exec ^
: string
| set local execute = 1
< ok 0 0
> exec ^
: string
| set local exists = 1
< error
| 1064 42000 …line 1 column 16 near "exists = 1" 
> exec ^
: string
| set local exit = 1
< error
| 1064 42000 …line 1 column 14 near "exit = 1" 
> exec ^
: string
| set local expansion = 1
< ok 0 0
> exec ^
: string
| set local expire = 1
< ok 0 0
> exec ^
: string
| set local explain = 1
< error
| 1064 42000 …line 1 column 17 near "explain = 1" 
> exec ^
: string
| set local export = 1
< ok 0 0
> exec ^
: string
| set local extended = 1
< ok 0 0
> exec ^
: string
| set local extent_size = 1
< ok 0 0
> exec ^
: string
| set local external = 1
< ok 0 0
> exec ^
: string
| set local external_format = 1
< ok 0 0
> exec ^
: string
| set local extract = 1
< ok 0 0
> exec ^
: string
| set local factor = 1
< ok 0 0
> exec ^
: string
| set local failed_login_attempts = 1
< ok 0 0
> exec ^
: string
| set local false = 1
< error
| 1064 42000 …line 1 column 15 near "false = 1" 
> exec ^
: string
| set local fast = 1
< ok 0 0
> exec ^
: string
| set local faults = 1
< ok 0 0
> exec ^
: string
| set local fetch = 1
< error
| 1064 42000 …line 1 column 15 near "fetch = 1" 
> exec ^
: string
| set local fields = 1
< ok 0 0
> exec ^
: string
| set local file = 1
< ok 0 0
> exec ^
: string
| set local file_block_size = 1
< ok 0 0
> exec ^
: string
| set local file_format = 1
< ok 0 0
> exec ^
: string
| set local file_name = 1
< ok 0 0
> exec ^
: string
| set local file_pattern = 1
< ok 0 0
> exec ^
: string
| set local file_prefix = 1
< ok 0 0
> exec ^
: string
| set local files = 1
< ok 0 0
> exec ^
: string
| set local filter = 1
< ok 0 0
> exec ^
: string
| set local finish = 1
< ok 0 0
> exec ^
: string
| set local first = 1
< ok 0 0
> exec ^
: string
| set local first_value = 1
< error
| 1064 42000 …line 1 column 21 near "first_value = 1" 
> exec ^
: string
| set local fixed = 1
< ok 0 0
> exec ^
: string
| set local flashback = 1
< ok 0 0
> exec ^
: string
| set local float = 1
< error
| 1064 42000 …line 1 column 15 near "float = 1" 
> exec ^
: string
| set local float4 = 1
< error
| 1064 42000 …line 1 column 16 near "float4 = 1" 
> exec ^
: string
| set local float8 = 1
< error
| 1064 42000 …line 1 column 16 near "float8 = 1" 
> exec ^
: string
| set local flush = 1
< ok 0 0
> exec ^
: string
| set local follower_constraints = 1
< ok 0 0
> exec ^
: string
| set local followers = 1
< ok 0 0
> exec ^
: string
| set local following = 1
< ok 0 0
> exec ^
: string
| set local follows = 1
< ok 0 0
> exec ^
: string
| set local for = 1
< error
| 1064 42000 …line 1 column 13 near "for = 1" 
> exec ^
: string
| set local force = 1
< error
| 1064 42000 …line 1 column 15 near "force = 1" 
> exec ^
: string
| set local foreign = 1
< error
| 1064 42000 …line 1 column 17 near "foreign = 1" 
> exec ^
: string
| set local format = 1
< ok 0 0
> exec ^
: string
| set local found = 1
< ok 0 0
> exec ^
: string
| set local from = 1
< error
| 1064 42000 …line 1 column 14 near "from = 1" 
> exec ^
: string
| set local full = 1
< ok 0 0
> exec ^
: string
| set local fulltext = 1
< error
| 1064 42000 …line 1 column 18 near "fulltext = 1" 
> exec ^
: string
| set local function = 1
< ok 0 0
> exec ^
: string
| set local general = 1
< ok 0 0
> exec ^
: string
| set local generate = 1
< ok 0 0
> exec ^
: string
| set local generated = 1
< error
| 1064 42000 …line 1 column 19 near "generated = 1" 
> exec ^
: string
| set local geomcollection = 1
< ok 0 0
> exec ^
: string
| set local geometry = 1
< ok 0 0
> exec ^
: string
| set local geometrycollection = 1
< ok 0 0
> exec ^
: string
| set local get = 1
< ok 0 0
> exec ^
: string
| set local get_format = 1
< ok 0 0
> exec ^
: string
| set local get_source_public_key = 1
< ok 0 0
> exec ^
: string
| set local global = 1
< ok 0 0
> exec ^
: string
| set local grant = 1
< error
| 1064 42000 …line 1 column 15 near "grant = 1" 
> exec ^
: string
| set local grants = 1
< ok 0 0
> exec ^
: string
| set local greatest = 1
< ok 0 0
> exec ^
: string
| set local group = 1
< error
| 1064 42000 …line 1 column 15 near "group = 1" 
> exec ^
: string
| set local group_concat = 1
< ok 0 0
> exec ^
: string
| set local group_replication = 1
< ok 0 0
> exec ^
: string
| set local grouping = 1
< ok 0 0
> exec ^
: string
| set local groups = 1
< error
| 1064 42000 …line 1 column 16 near "groups = 1" 
> exec ^
: string
| set local gtid_only = 1
< ok 0 0
> exec ^
: string
| set local gtids = 1
< ok 0 0
> exec ^
: string
| set local guided = 1
< ok 0 0
> exec ^
: string
| set local handler = 1
< ok 0 0
> exec ^
: string
| set local hash = 1
< ok 0 0
> exec ^
: string
| set local having = 1
< error
| 1064 42000 …line 1 column 16 near "having = 1" 
> exec ^
: string
| set local header = 1
< ok 0 0
> exec ^
: string
| set local help = 1
< ok 0 0
> exec ^
: string
| set local high = 1
< ok 0 0
> exec ^
: string
| set local high_priority = 1
< error
| 1064 42000 …line 1 column 23 near "high_priority = 1" 
> exec ^
: string
| set local histogram = 1
< ok 0 0
> exec ^
: string
| set local histograms_in_flight = 1
< ok 0 0
> exec ^
: string
| set local history = 1
< ok 0 0
> exec ^
: string
| set local host = 1
< ok 0 0
> exec ^
: string
| set local hosts = 1
< ok 0 0
> exec ^
: string
| set local hour = 1
< ok 0 0
> exec ^
: string
| set local hour_microsecond = 1
< error
| 1064 42000 …line 1 column 26 near "hour_microsecond = 1" 
> exec ^
: string
| set local hour_minute = 1
< error
| 1064 42000 …line 1 column 21 near "hour_minute = 1" 
> exec ^
: string
| set local hour_second = 1
< error
| 1064 42000 …line 1 column 21 near "hour_second = 1" 
> exec ^
: string
| set local identified = 1
< ok 0 0
> exec ^
: string
| set local if = 1
< error
| 1064 42000 …line 1 column 12 near "if = 1" 
> exec ^
: string
| set local ignore = 1
< error
| 1064 42000 …line 1 column 16 near "ignore = 1" 
> exec ^
: string
| set local ignore_server_ids = 1
< ok 0 0
> exec ^
: string
| set local ilike = 1
< error
| 1064 42000 …line 1 column 15 near "ilike = 1" 
> exec ^
: string
| set local import = 1
< ok 0 0
> exec ^
: string
| set local in = 1
< error
| 1064 42000 …line 1 column 12 near "in = 1" 
> exec ^
: string
| set local inactive = 1
< ok 0 0
> exec ^
: string
| set local increment = 1
< ok 0 0
> exec ^
: string
| set local index = 1
< error
| 1064 42000 …line 1 column 15 near "index = 1" 
> exec ^
: string
| set local indexes = 1
< ok 0 0
> exec ^
: string
| set local infile = 1
< error
| 1064 42000 …line 1 column 16 near "infile = 1" 
> exec ^
: string
| set local initial = 1
< ok 0 0
> exec ^
: string
| set local initial_size = 1
< ok 0 0
> exec ^
: string
| set local initiate = 1
< ok 0 0
> exec ^
: string
| set local inner = 1
< error
| 1064 42000 …line 1 column 15 near "inner = 1" 
> exec ^
: string
| set local inout = 1
< error
| 1064 42000 …line 1 column 15 near "inout = 1" 
> exec ^
: string
| set local insensitive = 1
< ok 0 0
> exec ^
: string
| set local insert = 1
< error
| 1064 42000 …line 1 column 16 near "insert = 1" 
> exec ^
: string
| set local insert_method = 1
< ok 0 0
> exec ^
: string
| set local install = 1
< ok 0 0
> exec ^
: string
| set local instance = 1
< ok 0 0
> exec ^
: string
| set local int = 1
< error
| 1064 42000 …line 1 column 13 near "int = 1" 
> exec ^
: string
| set local int1 = 1
< error
| 1064 42000 …line 1 column 14 near "int1 = 1" 
> exec ^
: string
| set local int2 = 1
< error
| 1064 42000 …line 1 column 14 near "int2 = 1" 
> exec ^
: string
| set local int3 = 1
< error
| 1064 42000 …line 1 column 14 near "int3 = 1" 
> exec ^
: string
| set local int4 = 1
< error
| 1064 42000 …line 1 column 14 near "int4 = 1" 
> exec ^
: string
| set local int8 = 1
< error
| 1064 42000 …line 1 column 14 near "int8 = 1" 
> exec ^
: string
| set local integer = 1
< error
| 1064 42000 …line 1 column 17 near "integer = 1" 
> exec ^
: string
| set local intersect = 1
< error
| 1064 42000 …line 1 column 19 near "intersect = 1" 
> exec ^
: string
| set local interval = 1
< error
| 1064 42000 …line 1 column 18 near "interval = 1" 
> exec ^
: string
| set local into = 1
< error
| 1064 42000 …line 1 column 14 near "into = 1" 
> exec ^
: string
| set local invisible = 1
< ok 0 0
> exec ^
: string
| set local invoker = 1
< ok 0 0
> exec ^
: string
| set local io = 1
< ok 0 0
> exec ^
: string
| set local io_after_gtids = 1
< ok 0 0
> exec ^
: string
| set local io_before_gtids = 1
< ok 0 0
> exec ^
: string
| set local io_thread = 1
< ok 0 0
> exec ^
: string
| set local ipc = 1
< ok 0 0
> exec ^
: string
| set local is = 1
< error
| 1064 42000 …line 1 column 12 near "is = 1" 
> exec ^
: string
| set local isolation = 1
< ok 0 0
> exec ^
: string
| set local issuer = 1
< ok 0 0
> exec ^
: string
| set local iterate = 1
< error
| 1064 42000 …line 1 column 17 near "iterate = 1" 
> exec ^
: string
| set local job = 1
< ok 0 0
> exec ^
: string
| set local jobs = 1
< ok 0 0
> exec ^
: string
| set local join = 1
< error
| 1064 42000 …line 1 column 14 near "join = 1" 
> exec ^
: string
| set local json = 1
< ok 0 0
> exec ^
: string
| set local json_arrayagg = 1
< ok 0 0
> exec ^
: string
| set local json_duality_object = 1
< ok 0 0
> exec ^
: string
| set local json_objectagg = 1
< ok 0 0
> exec ^
: string
| set local json_table = 1
< ok 0 0
> exec ^
: string
| set local json_value = 1
< ok 0 0
> exec ^
: string
| set local key = 1
< error
| 1064 42000 …line 1 column 13 near "key = 1" 
> exec ^
: string
| set local key_block_size = 1
< ok 0 0
> exec ^
: string
| set local keyring = 1
< ok 0 0
> exec ^
: string
| set local keys = 1
< error
| 1064 42000 …line 1 column 14 near "keys = 1" 
> exec ^
: string
| set local kill = 1
< error
| 1064 42000 …line 1 column 14 near "kill = 1" 
> exec ^
: string
| set local labels = 1
< ok 0 0
> exec ^
: string
| set local lag = 1
< error
| 1064 42000 …line 1 column 13 near "lag = 1" 
> exec ^
: string
| set local language = 1
< ok 0 0
> exec ^
: string
| set local last = 1
< ok 0 0
> exec ^
: string
| set local last_value = 1
< error
| 1064 42000 …line 1 column 20 near "last_value = 1" 
> exec ^
: string
| set local lateral = 1
< ok 0 0
> exec ^
: string
| set local lead = 1
< error
| 1064 42000 …line 1 column 14 near "lead = 1" 
> exec ^
: string
| set local leader_constraints = 1
< ok 0 0
> exec ^
: string
| set local leading = 1
< error
| 1064 42000 …line 1 column 17 near "leading = 1" 
> exec ^
: string
| set local learner_constraints = 1
< ok 0 0
> exec ^
: string
| set local learners = 1
< ok 0 0
> exec ^
: string
| set local least = 1
< ok 0 0
> exec ^
: string
| set local leave = 1
< error
| 1064 42000 …line 1 column 15 near "leave = 1" 
> exec ^
: string
| set local leaves = 1
< ok 0 0
> exec ^
: string
| set local left = 1
< error
| 1064 42000 …line 1 column 14 near "left = 1" 
> exec ^
: string
| set local less = 1
< ok 0 0
> exec ^
: string
| set local level = 1
< ok 0 0
> exec ^
: string
| set local library = 1
< ok 0 0
> exec ^
: string
| set local like = 1
< error
| 1064 42000 …line 1 column 14 near "like = 1" 
> exec ^
: string
| set local limit = 1
< error
| 1064 42000 …line 1 column 15 near "limit = 1" 
> exec ^
: string
| set local linear = 1
< error
| 1064 42000 …line 1 column 16 near "linear = 1" 
> exec ^
: string
| set local lines = 1
< error
| 1064 42000 …line 1 column 15 near "lines = 1" 
> exec ^
: string
| set local linestring = 1
< ok 0 0
> exec ^
: string
| set local list = 1
< ok 0 0
> exec ^
: string
| set local load = 1
< error
| 1064 42000 …line 1 column 14 near "load = 1" 
> exec ^
: string
| set local local = 1
< ok 0 0
> exec ^
: string
| set local localtime = 1
< error
| 1064 42000 …line 1 column 19 near "localtime = 1" 
> exec ^
: string
| set local localtimestamp = 1
< error
| 1064 42000 …line 1 column 24 near "localtimestamp = 1" 
> exec ^
: string
| set local location = 1
< ok 0 0
> exec ^
: string
| set local lock = 1
< error
| 1064 42000 …line 1 column 14 near "lock = 1" 
> exec ^
: string
| set local locked = 1
< ok 0 0
> exec ^
: string
| set local locks = 1
< ok 0 0
> exec ^
: string
| set local log = 1
< ok 0 0
> exec ^
: string
| set local logfile = 1
< ok 0 0
> exec ^
: string
| set local logs = 1
< ok 0 0
> exec ^
: string
| set local long = 1
< error
| 1064 42000 …line 1 column 14 near "long = 1" 
> exec ^
: string
| set local longblob = 1
< error
| 1064 42000 …line 1 column 18 near "longblob = 1" 
> exec ^
: string
| set local longtext = 1
< error
| 1064 42000 …line 1 column 18 near "longtext = 1" 
> exec ^
: string
| set local loop = 1
< ok 0 0
> exec ^
: string
| set local low = 1
< ok 0 0
> exec ^
: string
| set local low_priority = 1
< error
| 1064 42000 …line 1 column 22 near "low_priority = 1" 
> exec ^
: string
| set local manual = 1
< ok 0 0
> exec ^
: string
| set local master = 1
< ok 0 0
> exec ^
: string
| set local master_log_file = 1
< ok 0 0
> exec ^
: string
| set local match = 1
< error
| 1064 42000 …line 1 column 15 near "match = 1" 
> exec ^
: string
| set local materialized = 1
< ok 0 0
> exec ^
: string
| set local max = 1
< ok 0 0
> exec ^
: string
| set local max_connections_per_hour = 1
< ok 0 0
> exec ^
: string
| set local max_queries_per_hour = 1
< ok 0 0
> exec ^
: string
| set local max_rows = 1
< ok 0 0
> exec ^
: string
| set local max_size = 1
< ok 0 0
> exec ^
: string
| set local max_updates_per_hour = 1
< ok 0 0
> exec ^
: string
| set local max_user_connections = 1
< ok 0 0
> exec ^
: string
| set local maxvalue = 1
< error
| 1064 42000 …line 1 column 18 near "maxvalue = 1" 
> exec ^
: string
| set local medium = 1
< ok 0 0
> exec ^
: string
| set local mediumblob = 1
< error
| 1064 42000 …line 1 column 20 near "mediumblob = 1" 
> exec ^
: string
| set local mediumint = 1
< error
| 1064 42000 …line 1 column 19 near "mediumint = 1" 
> exec ^
: string
| set local mediumtext = 1
< error
| 1064 42000 …line 1 column 20 near "mediumtext = 1" 
> exec ^
: string
| set local member = 1
< ok 0 0
> exec ^
: string
| set local memory = 1
< ok 0 0
> exec ^
: string
| set local merge = 1
< ok 0 0
> exec ^
: string
| set local message_text = 1
< ok 0 0
> exec ^
: string
| set local microsecond = 1
< ok 0 0
> exec ^
: string
| set local mid = 1
< ok 0 0
> exec ^
: string
| set local middleint = 1
< error
| 1064 42000 …line 1 column 19 near "middleint = 1" 
> exec ^
: string
| set local migrate = 1
< ok 0 0
> exec ^
: string
| set local min = 1
< ok 0 0
> exec ^
: string
| set local min_rows = 1
< ok 0 0
> exec ^
: string
| set local minute = 1
< ok 0 0
> exec ^
: string
| set local minute_microsecond = 1
< error
| 1064 42000 …line 1 column 28 near "minute_microsecond = 1" 
> exec ^
: string
| set local minute_second = 1
< error
| 1064 42000 …line 1 column 23 near "minute_second = 1" 
> exec ^
: string
| set local minvalue = 1
< ok 0 0
> exec ^
: string
| set local mod = 1
< error
| 1064 42000 …line 1 column 13 near "mod = 1" 
> exec ^
: string
| set local mode = 1
< ok 0 0
> exec ^
: string
| set local modifies = 1
< ok 0 0
> exec ^
: string
| set local modify = 1
< ok 0 0
> exec ^
: string
| set local month = 1
< ok 0 0
> exec ^
: string
| set local multilinestring = 1
< ok 0 0
> exec ^
: string
| set local multipoint = 1
< ok 0 0
> exec ^
: string
| set local multipolygon = 1
< ok 0 0
> exec ^
: string
| set local mutex = 1
< ok 0 0
> exec ^
: string
| set local mysql_errno = 1
< ok 0 0
> exec ^
: string
| set local name = 1
< ok 0 0
> exec ^
: string
| set local names = 1
< ok 0 0
> exec ^
: string
| set local national = 1
< ok 0 0
> exec ^
: string
| set local natural = 1
< error
| 1064 42000 …line 1 column 17 near "natural = 1" 
> exec ^
: string
| set local nchar = 1
< ok 0 0
> exec ^
: string
| set local ndb = 1
< ok 0 0
> exec ^
: string
| set local ndbcluster = 1
< ok 0 0
> exec ^
: string
| set local nested = 1
< ok 0 0
> exec ^
: string
| set local network_namespace = 1
< ok 0 0
> exec ^
: string
| set local never = 1
< ok 0 0
> exec ^
: string
| set local new = 1
< ok 0 0
> exec ^
: string
| set local next = 1
< ok 0 0
> exec ^
: string
| set local next_row_id = 1
< ok 0 0
> exec ^
: string
| set local no = 1
< ok 0 0
> exec ^
: string
| set local no_wait = 1
< ok 0 0
> exec ^
: string
| set local no_write_to_binlog = 1
< error
| 1064 42000 …line 1 column 28 near "no_write_to_binlog = 1" 
> exec ^
: string
| set local nocache = 1
< ok 0 0
> exec ^
: string
| set local nocycle = 1
< ok 0 0
> exec ^
: string
| set local nodegroup = 1
< ok 0 0
> exec ^
: string
| set local nomaxvalue = 1
< ok 0 0
> exec ^
: string
| set local nominvalue = 1
< ok 0 0
> exec ^
: string
| set local nonclustered = 1
< ok 0 0
> exec ^
: string
| set local none = 1
< ok 0 0
> exec ^
: string
| set local not = 1
< error
| 1064 42000 …line 1 column 13 near "not = 1" 
> exec ^
: string
| set local now = 1
< ok 0 0
> exec ^
: string
| set local nowait = 1
< ok 0 0
> exec ^
: string
| set local nth_value = 1
< error
| 1064 42000 …line 1 column 19 near "nth_value = 1" 
> exec ^
: string
| set local ntile = 1
< error
| 1064 42000 …line 1 column 15 near "ntile = 1" 
> exec ^
: string
| set local null = 1
< error
| 1064 42000 …line 1 column 14 near "null = 1" 
> exec ^
: string
| set local nullif = 1
< ok 0 0
> exec ^
: string
| set local nulls = 1
< ok 0 0
> exec ^
: string
| set local number = 1
< ok 0 0
> exec ^
: string
| set local numeric = 1
< error
| 1064 42000 …line 1 column 17 near "numeric = 1" 
> exec ^
: string
| set local nvarchar = 1
< ok 0 0
> exec ^
: string
| set local of = 1
< error
| 1064 42000 …line 1 column 12 near "of = 1" 
> exec ^
: string
| set local off = 1
< ok 0 0
> exec ^
: string
| set local offset = 1
< ok 0 0
> exec ^
: string
| set local oj = 1
< ok 0 0
> exec ^
: string
| set local old = 1
< ok 0 0
> exec ^
: string
| set local oltp_read_only = 1
< ok 0 0
> exec ^
: string
| set local oltp_read_write = 1
< ok 0 0
> exec ^
: string
| set local oltp_write_only = 1
< ok 0 0
> exec ^
: string
| set local on = 1
< error
| 1064 42000 …line 1 column 12 near "on = 1" 
> exec ^
: string
| set local one = 1
< ok 0 0
> exec ^
: string
| set local only = 1
< ok 0 0
> exec ^
: string
| set local open = 1
< ok 0 0
> exec ^
: string
| set local optimize = 1
< error
| 1064 42000 …line 1 column 18 near "optimize = 1" 
> exec ^
: string
| set local optimizer_costs = 1
< ok 0 0
> exec ^
: string
| set local option = 1
< error
| 1064 42000 …line 1 column 16 near "option = 1" 
> exec ^
: string
| set local optional = 1
< ok 0 0
> exec ^
: string
| set local optionally = 1
< error
| 1064 42000 …line 1 column 20 near "optionally = 1" 
> exec ^
: string
| set local options = 1
< ok 0 0
> exec ^
: string
| set local or = 1
< error
| 1064 42000 …line 1 column 12 near "or = 1" 
> exec ^
: string
| set local order = 1
< error
| 1064 42000 …line 1 column 15 near "order = 1" 
> exec ^
: string
| set local ordinality = 1
< ok 0 0
> exec ^
: string
| set local organization = 1
< ok 0 0
> exec ^
: string
| set local others = 1
< ok 0 0
> exec ^
: string
| set local out = 1
< error
| 1064 42000 …line 1 column 13 near "out = 1" 
> exec ^
: string
| set local outer = 1
< error
| 1064 42000 …line 1 column 15 near "outer = 1" 
> exec ^
: string
| set local outfile = 1
< error
| 1064 42000 …line 1 column 17 near "outfile = 1" 
> exec ^
: string
| set local over = 1
< error
| 1064 42000 …line 1 column 14 near "over = 1" 
> exec ^
: string
| set local overlay = 1
< ok 0 0
> exec ^
: string
| set local owner = 1
< ok 0 0
> exec ^
: string
| set local pack_keys = 1
< ok 0 0
> exec ^
: string
| set local page = 1
< ok 0 0
> exec ^
: string
| set local parallel = 1
< ok 0 0
> exec ^
: string
| set local parameters = 1
< ok 0 0
> exec ^
: string
| set local parse_tree = 1
< ok 0 0
> exec ^
: string
| set local parser = 1
< ok 0 0
> exec ^
: string
| set local partial = 1
< ok 0 0
> exec ^
: string
| set local partition = 1
< error
| 1064 42000 …line 1 column 19 near "partition = 1" 
> exec ^
: string
| set local partitioning = 1
< ok 0 0
> exec ^
: string
| set local partitions = 1
< ok 0 0
> exec ^
: string
| set local password = 1
< ok 0 0
> exec ^
: string
| set local password_lock_time = 1
< ok 0 0
> exec ^
: string
| set local path = 1
< ok 0 0
> exec ^
: string
| set local percent_rank = 1
< error
| 1064 42000 …line 1 column 22 near "percent_rank = 1" 
> exec ^
: string
| set local persist = 1
< ok 0 0
> exec ^
: string
| set local persist_only = 1
< ok 0 0
> exec ^
: string
| set local phase = 1
< ok 0 0
> exec ^
: string
| set local placement = 1
< ok 0 0
> exec ^
: string
| set local placing = 1
< ok 0 0
> exec ^
: string
| set local plan = 1
< ok 0 0
> exec ^
: string
| set local plugin = 1
< ok 0 0
> exec ^
: string
| set local plugin_dir = 1
< ok 0 0
> exec ^
: string
| set local plugins = 1
< ok 0 0
> exec ^
: string
| set local point = 1
< ok 0 0
> exec ^
: string
| set local policy = 1
< ok 0 0
> exec ^
: string
| set local polygon = 1
< ok 0 0
> exec ^
: string
| set local port = 1
< ok 0 0
> exec ^
: string
| set local position = 1
< ok 0 0
> exec ^
: string
| set local pre_split_regions = 1
< ok 0 0
> exec ^
: string
| set local precedes = 1
< ok 0 0
> exec ^
: string
| set local preceding = 1
< ok 0 0
> exec ^
: string
| set local precision = 1
< error
| 1064 42000 …line 1 column 19 near "precision = 1" 
> exec ^
: string
| set local prepare = 1
< ok 0 0
> exec ^
: string
| set local preserve = 1
< ok 0 0
> exec ^
: string
| set local prev = 1
< ok 0 0
> exec ^
: string
| set local primary = 1
< error
| 1064 42000 …line 1 column 17 near "primary = 1" 
> exec ^
: string
| set local primary_region = 1
< ok 0 0
> exec ^
: string
| set local priority = 1
< ok 0 0
> exec ^
: string
| set local privilege_checks_user = 1
< ok 0 0
> exec ^
: string
| set local privileges = 1
< ok 0 0
> exec ^
: string
| set local procedure = 1
< error
| 1064 42000 …line 1 column 19 near "procedure = 1" 
> exec ^
: string
| set local process = 1
< ok 0 0
> exec ^
: string
| set local processed_keys = 1
< ok 0 0
> exec ^
: string
| set local processlist = 1
< ok 0 0
> exec ^
: string
| set local profile = 1
< ok 0 0
> exec ^
: string
| set local profiles = 1
< ok 0 0
> exec ^
: string
| set local proxy = 1
< ok 0 0
> exec ^
: string
| set local purge = 1
< ok 0 0
> exec ^
: string
| set local qualify = 1
< ok 0 0
> exec ^
: string
| set local quarter = 1
< ok 0 0
> exec ^
: string
| set local query = 1
< ok 0 0
> exec ^
: string
| set local query_limit = 1
< ok 0 0
> exec ^
: string
| set local quick = 1
< ok 0 0
> exec ^
: string
| set local random = 1
< ok 0 0
> exec ^
: string
| set local range = 1
< error
| 1064 42000 …line 1 column 15 near "range = 1" 
> exec ^
: string
| set local rank = 1
< error
| 1064 42000 …line 1 column 14 near "rank = 1" 
> exec ^
: string
| set local read = 1
< error
| 1064 42000 …line 1 column 14 near "read = 1" 
> exec ^
: string
| set local read_only = 1
< ok 0 0
> exec ^
: string
| set local read_write = 1
< ok 0 0
> exec ^
: string
| set local reads = 1
< ok 0 0
> exec ^
: string
| set local real = 1
< error
| 1064 42000 …line 1 column 14 near "real = 1" 
> exec ^
: string
| set local rebuild = 1
< ok 0 0
> exec ^
: string
| set local recover = 1
< ok 0 0
> exec ^
: string
| set local recursive = 1
< error
| 1064 42000 …line 1 column 19 near "recursive = 1" 
> exec ^
: string
| set local redo_buffer_size = 1
< ok 0 0
> exec ^
: string
| set local redundant = 1
< ok 0 0
> exec ^
: string
| set local reference = 1
< ok 0 0
> exec ^
: string
| set local references = 1
< error
| 1064 42000 …line 1 column 20 near "references = 1" 
> exec ^
: string
| set local regexp = 1
< error
| 1064 42000 …line 1 column 16 near "regexp = 1" 
> exec ^
: string
| set local region = 1
< ok 0 0
> exec ^
: string
| set local regions = 1
< ok 0 0
> exec ^
: string
| set local registration = 1
< ok 0 0
> exec ^
: string
| set local relational = 1
< ok 0 0
> exec ^
: string
| set local relay = 1
< ok 0 0
> exec ^
: string
| set local relay_log_file = 1
< ok 0 0
> exec ^
: string
| set local relay_log_pos = 1
< ok 0 0
> exec ^
: string
| set local relay_thread = 1
< ok 0 0
> exec ^
: string
| set local relaylog = 1
< ok 0 0
> exec ^
: string
| set local release = 1
< error
| 1064 42000 …line 1 column 17 near "release = 1" 
> exec ^
: string
| set local reload = 1
< ok 0 0
> exec ^
: string
| set local remove = 1
< ok 0 0
> exec ^
: string
| set local rename = 1
< error
| 1064 42000 …line 1 column 16 near "rename = 1" 
> exec ^
: string
| set local reorganize = 1
< ok 0 0
> exec ^
: string
| set local repair = 1
< ok 0 0
> exec ^
: string
| set local repeat = 1
< error
| 1064 42000 …line 1 column 16 near "repeat = 1" 
> exec ^
: string
| set local repeatable = 1
< ok 0 0
> exec ^
: string
| set local replace = 1
< error
| 1064 42000 …line 1 column 17 near "replace = 1" 
> exec ^
: string
| set local replica = 1
< ok 0 0
> exec ^
: string
| set local replicas = 1
< ok 0 0
> exec ^
: string
| set local replicate_do_db = 1
< ok 0 0
> exec ^
: string
| set local replicate_do_table = 1
< ok 0 0
> exec ^
: string
| set local replicate_ignore_db = 1
< ok 0 0
> exec ^
: string
| set local replicate_ignore_table = 1
< ok 0 0
> exec ^
: string
| set local replicate_rewrite_db = 1
< ok 0 0
> exec ^
: string
| set local replicate_wild_do_table = 1
< ok 0 0
> exec ^
: string
| set local replicate_wild_ignore_table = 1
< ok 0 0
> exec ^
: string
| set local replication = 1
< ok 0 0
> exec ^
: string
| set local require = 1
< error
| 1064 42000 …line 1 column 17 near "require = 1" 
> exec ^
: string
| set local require_row_format = 1
< ok 0 0
> exec ^
: string
| set local require_table_primary_key_check = 1
< ok 0 0
> exec ^
: string
| set local reset = 1
< ok 0 0
> exec ^
: string
| set local resignal = 1
< ok 0 0
> exec ^
: string
| set local resource = 1
< ok 0 0
> exec ^
: string
| set local respect = 1
< ok 0 0
> exec ^
: string
| set local restart = 1
< ok 0 0
> exec ^
: string
| set local restore = 1
< ok 0 0
> exec ^
: string
| set local restores = 1
< ok 0 0
> exec ^
: string
| set local restrict = 1
< error
| 1064 42000 …line 1 column 18 near "restrict = 1" 
> exec ^
: string
| set local resume = 1
< ok 0 0
> exec ^
: string
| set local retain = 1
< ok 0 0
> exec ^
: string
| set local return = 1
< ok 0 0
> exec ^
: string
| set local returned_sqlstate = 1
< ok 0 0
> exec ^
: string
| set local returning = 1
< ok 0 0
> exec ^
: string
| set local returns = 1
< ok 0 0
> exec ^
: string
| set local reuse = 1
< ok 0 0
> exec ^
: string
| set local reverse = 1
< ok 0 0
> exec ^
: string
| set local revoke = 1
< error
| 1064 42000 …line 1 column 16 near "revoke = 1" 
> exec ^
: string
| set local right = 1
< error
| 1064 42000 …line 1 column 15 near "right = 1" 
> exec ^
: string
| set local rlike = 1
< error
| 1064 42000 …line 1 column 15 near "rlike = 1" 
> exec ^
: string
| set local role = 1
< ok 0 0
> exec ^
: string
| set local rollback = 1
< ok 0 0
> exec ^
: string
| set local rollup = 1
< ok 0 0
> exec ^
: string
| set local rotate = 1
< ok 0 0
> exec ^
: string
| set local routine = 1
< ok 0 0
> exec ^
: string
| set local row = 1
< error
| 1064 42000 …line 1 column 13 near "row = 1" 
> exec ^
: string
| set local row_count = 1
< ok 0 0
> exec ^
: string
| set local row_format = 1
< ok 0 0
> exec ^
: string
| set local row_number = 1
< error
| 1064 42000 …line 1 column 20 near "row_number = 1" 
> exec ^
: string
| set local rows = 1
< error
| 1064 42000 …line 1 column 14 near "rows = 1" 
> exec ^
: string
| set local rtree = 1
< ok 0 0
> exec ^
: string
| set local ru = 1
< ok 0 0
> exec ^
: string
| set local ru_per_sec = 1
< ok 0 0
> exec ^
: string
| set local rule = 1
< ok 0 0
> exec ^
: string
| set local run = 1
< ok 0 0
> exec ^
: string
| set local s3 = 1
< ok 0 0
> exec ^
: string
| set local savepoint = 1
< ok 0 0
> exec ^
: string
| set local schedule = 1
< ok 0 0
> exec ^
: string
| set local schema = 1
< error
| 1064 42000 …line 1 column 16 near "schema = 1" 
> exec ^
: string
| set local schema_name = 1
< ok 0 0
> exec ^
: string
| set local schemas = 1
< error
| 1064 42000 …line 1 column 17 near "schemas = 1" 
> exec ^
: string
| set local second = 1
< ok 0 0
> exec ^
: string
| set local second_microsecond = 1
< error
| 1064 42000 …line 1 column 28 near "second_microsecond = 1" 
> exec ^
: string
| set local secondary = 1
< ok 0 0
> exec ^
: string
| set local secondary_engine = 1
< ok 0 0
> exec ^
: string
| set local secondary_engine_attribute = 1
< ok 0 0
> exec ^
: string
| set local secondary_load = 1
< ok 0 0
> exec ^
: string
| set local secondary_unload = 1
< ok 0 0
> exec ^
: string
| set local security = 1
< ok 0 0
> exec ^
: string
| set local select = 1
< error
| 1064 42000 …line 1 column 16 near "select = 1" 
> exec ^
: string
| set local sensitive = 1
< ok 0 0
> exec ^
: string
| set local separator = 1
< ok 0 0
> exec ^
: string
| set local sequence = 1
< ok 0 0
> exec ^
: string
| set local serial = 1
< ok 0 0
> exec ^
: string
| set local serializable = 1
< ok 0 0
> exec ^
: string
| set local server = 1
< ok 0 0
> exec ^
: string
| set local session = 1
< ok 0 0
> exec ^
: string
| set local session_states = 1
< ok 0 0
> exec ^
: string
| set local session_user = 1
< ok 0 0
> exec ^
: string
| set local set = 1
< error
| 1064 42000 …line 1 column 13 near "set = 1" 
> exec ^
: string
| set local sets = 1
< ok 0 0
> exec ^
: string
| set local shard_row_id_bits = 1
< ok 0 0
> exec ^
: string
| set local share = 1
< ok 0 0
> exec ^
: string
| set local show = 1
< error
| 1064 42000 …line 1 column 14 near "show = 1" 
> exec ^
: string
| set local shutdown = 1
< ok 0 0
> exec ^
: string
| set local signal = 1
< ok 0 0
> exec ^
: string
| set local signed = 1
< ok 0 0
> exec ^
: string
| set local similar = 1
< ok 0 0
> exec ^
: string
| set local simple = 1
< ok 0 0
> exec ^
: string
| set local skip = 1
< ok 0 0
> exec ^
: string
| set local slave = 1
< ok 0 0
> exec ^
: string
| set local slow = 1
< ok 0 0
> exec ^
: string
| set local smallint = 1
< error
| 1064 42000 …line 1 column 18 near "smallint = 1" 
> exec ^
: string
| set local snapshot = 1
< ok 0 0
> exec ^
: string
| set local socket = 1
< ok 0 0
> exec ^
: string
| set local some = 1
< ok 0 0
> exec ^
: string
| set local soname = 1
< ok 0 0
> exec ^
: string
| set local sounds = 1
< ok 0 0
> exec ^
: string
| set local source = 1
< ok 0 0
> exec ^
: string
| set local source_auto_position = 1
< ok 0 0
> exec ^
: string
| set local source_bind = 1
< ok 0 0
> exec ^
: string
| set local source_compression_algorithms = 1
< ok 0 0
> exec ^
: string
| set local source_connect_retry = 1
< ok 0 0
> exec ^
: string
| set local source_connection_auto_failover = 1
< ok 0 0
> exec ^
: string
| set local source_delay = 1
< ok 0 0
> exec ^
: string
| set local source_heartbeat_period = 1
< ok 0 0
> exec ^
: string
| set local source_host = 1
< ok 0 0
> exec ^
: string
| set local source_log_file = 1
< ok 0 0
> exec ^
: string
| set local source_log_pos = 1
< ok 0 0
> exec ^
: string
| set local source_password = 1
< ok 0 0
> exec ^
: string
| set local source_port = 1
< ok 0 0
> exec ^
: string
| set local source_public_key_path = 1
< ok 0 0
> exec ^
: string
| set local source_retry_count = 1
< ok 0 0
> exec ^
: string
| set local source_ssl = 1
< ok 0 0
> exec ^
: string
| set local source_ssl_ca = 1
< ok 0 0
> exec ^
: string
| set local source_ssl_capath = 1
< ok 0 0
> exec ^
: string
| set local source_ssl_cert = 1
< ok 0 0
> exec ^
: string
| set local source_ssl_cipher = 1
< ok 0 0
> exec ^
: string
| set local source_ssl_crl = 1
< ok 0 0
> exec ^
: string
| set local source_ssl_crlpath = 1
< ok 0 0
> exec ^
: string
| set local source_ssl_key = 1
< ok 0 0
> exec ^
: string
| set local source_ssl_verify_server_cert = 1
< ok 0 0
> exec ^
: string
| set local source_tls_ciphersuites = 1
< ok 0 0
> exec ^
: string
| set local source_tls_version = 1
< ok 0 0
> exec ^
: string
| set local source_user = 1
< ok 0 0
> exec ^
: string
| set local source_zstd_compression_level = 1
< ok 0 0
> exec ^
: string
| set local spatial = 1
< error
| 1064 42000 …line 1 column 17 near "spatial = 1" 
> exec ^
: string
| set local specific = 1
< ok 0 0
> exec ^
: string
| set local split = 1
< ok 0 0
> exec ^
: string
| set local sql = 1
< error
| 1064 42000 …line 1 column 13 near "sql = 1" 
> exec ^
: string
| set local sql_after_gtids = 1
< ok 0 0
> exec ^
: string
| set local sql_after_mts_gaps = 1
< ok 0 0
> exec ^
: string
| set local sql_before_gtids = 1
< ok 0 0
> exec ^
: string
| set local sql_big_result = 1
< error
| 1064 42000 …line 1 column 24 near "sql_big_result = 1" 
> exec ^
: string
| set local sql_buffer_result = 1
< ok 0 0
> exec ^
: string
| set local sql_calc_found_rows = 1
< error
| 1064 42000 …line 1 column 29 near "sql_calc_found_rows = 1" 
> exec ^
: string
| set local sql_no_cache = 1
< ok 0 0
> exec ^
: string
| set local sql_small_result = 1
< error
| 1064 42000 …line 1 column 26 near "sql_small_result = 1" 
> exec ^
: string
| set local sql_thread = 1
< ok 0 0
> exec ^
: string
| set local sql_tsi_day = 1
< ok 0 0
> exec ^
: string
| set local sql_tsi_hour = 1
< ok 0 0
> exec ^
: string
| set local sql_tsi_minute = 1
< ok 0 0
> exec ^
: string
| set local sql_tsi_month = 1
< ok 0 0
> exec ^
: string
| set local sql_tsi_quarter = 1
< ok 0 0
> exec ^
: string
| set local sql_tsi_second = 1
< ok 0 0
> exec ^
: string
| set local sql_tsi_week = 1
< ok 0 0
> exec ^
: string
| set local sql_tsi_year = 1
< ok 0 0
> exec ^
: string
| set local sqlexception = 1
< error
| 1064 42000 …line 1 column 22 near "sqlexception = 1" 
> exec ^
: string
| set local sqlstate = 1
< error
| 1064 42000 …line 1 column 18 near "sqlstate = 1" 
> exec ^
: string
| set local sqlwarning = 1
< error
| 1064 42000 …line 1 column 20 near "sqlwarning = 1" 
> exec ^
: string
| set local srid = 1
< ok 0 0
> exec ^
: string
| set local ssl = 1
< error
| 1064 42000 …line 1 column 13 near "ssl = 1" 
> exec ^
: string
| set local st_collect = 1
< ok 0 0
> exec ^
: string
| set local stacked = 1
< ok 0 0
> exec ^
: string
| set local start = 1
< ok 0 0
> exec ^
: string
| set local start_time = 1
< ok 0 0
> exec ^
: string
| set local starting = 1
< error
| 1064 42000 …line 1 column 18 near "starting = 1" 
> exec ^
: string
| set local starts = 1
< ok 0 0
> exec ^
: string
| set local stats_auto_recalc = 1
< ok 0 0
> exec ^
: string
| set local stats_buckets = 1
< ok 0 0
> exec ^
: string
| set local stats_col_choice = 1
< ok 0 0
> exec ^
: string
| set local stats_col_list = 1
< ok 0 0
> exec ^
: string
| set local stats_extended = 1
< ok 0 0
> exec ^
: string
| set local stats_healthy = 1
< ok 0 0
> exec ^
: string
| set local stats_histograms = 1
< ok 0 0
> exec ^
: string
| set local stats_locked = 1
< ok 0 0
> exec ^
: string
| set local stats_meta = 1
< ok 0 0
> exec ^
: string
| set local stats_persistent = 1
< ok 0 0
> exec ^
: string
| set local stats_sample_pages = 1
< ok 0 0
> exec ^
: string
| set local stats_sample_rate = 1
< ok 0 0
> exec ^
: string
| set local stats_topn = 1
< ok 0 0
> exec ^
: string
| set local status = 1
< ok 0 0
> exec ^
: string
| set local std = 1
< ok 0 0
> exec ^
: string
| set local stddev = 1
< ok 0 0
> exec ^
: string
| set local stddev_pop = 1
< ok 0 0
> exec ^
: string
| set local stddev_samp = 1
< ok 0 0
> exec ^
: string
| set local stop = 1
< ok 0 0
> exec ^
: string
| set local storage = 1
< ok 0 0
> exec ^
: string
| set local stored = 1
< error
| 1064 42000 …line 1 column 16 near "stored = 1" 
> exec ^
: string
| set local straight_join = 1
< error
| 1064 42000 …line 1 column 23 near "straight_join = 1" 
> exec ^
: string
| set local stream = 1
< ok 0 0
> exec ^
: string
| set local strict_load = 1
< ok 0 0
> exec ^
: string
| set local string = 1
< ok 0 0
> exec ^
: string
| set local subclass_origin = 1
< ok 0 0
> exec ^
: string
| set local subdate = 1
< ok 0 0
> exec ^
: string
| set local subject = 1
< ok 0 0
> exec ^
: string
| set local subpartition = 1
< ok 0 0
> exec ^
: string
| set local subpartitions = 1
< ok 0 0
> exec ^
: string
| set local substr = 1
< ok 0 0
> exec ^
: string
| set local substring = 1
< ok 0 0
> exec ^
: string
| set local sum = 1
< ok 0 0
> exec ^
: string
| set local super = 1
< ok 0 0
> exec ^
: string
| set local survival_preferences = 1
< ok 0 0
> exec ^
: string
| set local suspend = 1
< ok 0 0
> exec ^
: string
| set local swaps = 1
< ok 0 0
> exec ^
: string
| set local switch_group = 1
< ok 0 0
> exec ^
: string
| set local switches = 1
< ok 0 0
> exec ^
: string
| set local sysdate = 1
< ok 0 0
> exec ^
: string
| set local system = 1
< ok 0 0
> exec ^
: string
| set local system_user = 1
< ok 0 0
> exec ^
: string
| set local table = 1
< error
| 1064 42000 …line 1 column 15 near "table = 1" 
> exec ^
: string
| set local table_checksum = 1
< ok 0 0
> exec ^
: string
| set local table_name = 1
< ok 0 0
> exec ^
: string
| set local tables = 1
< ok 0 0
> exec ^
: string
| set local tablesample = 1
< error
| 1064 42000 …line 1 column 21 near "tablesample = 1" 
> exec ^
: string
| set local tablespace = 1
< ok 0 0
> exec ^
: string
| set local task_types = 1
< ok 0 0
> exec ^
: string
| set local temporary = 1
< ok 0 0
> exec ^
: string
| set local temptable = 1
< ok 0 0
> exec ^
: string
| set local terminated = 1
< error
| 1064 42000 …line 1 column 20 near "terminated = 1" 
> exec ^
: string
| set local text = 1
< ok 0 0
> exec ^
: string
| set local than = 1
< ok 0 0
> exec ^
: string
| set local then = 1
< error
| 1064 42000 …line 1 column 14 near "then = 1" 
> exec ^
: string
| set local thread_priority = 1
< ok 0 0
> exec ^
: string
| set local ties = 1
< ok 0 0
> exec ^
: string
| set local tiflash = 1
< ok 0 0
> exec ^
: string
| set local time = 1
< ok 0 0
> exec ^
: string
| set local timeout = 1
< ok 0 0
> exec ^
: string
| set local timestamp = 1
< ok 0 0
> exec ^
: string
| set local timestampadd = 1
< ok 0 0
> exec ^
: string
| set local timestampdiff = 1
< ok 0 0
> exec ^
: string
| set local tinyblob = 1
< error
| 1064 42000 …line 1 column 18 near "tinyblob = 1" 
> exec ^
: string
| set local tinyint = 1
< error
| 1064 42000 …line 1 column 17 near "tinyint = 1" 
> exec ^
: string
| set local tinytext = 1
< error
| 1064 42000 …line 1 column 18 near "tinytext = 1" 
> exec ^
: string
| set local tls = 1
< ok 0 0
> exec ^
: string
| set local to = 1
< error
| 1064 42000 …line 1 column 12 near "to = 1" 
> exec ^
: string
| set local tpcc = 1
< ok 0 0
> exec ^
: string
| set local tpch_10 = 1
< ok 0 0
> exec ^
: string
| set local trailing = 1
< error
| 1064 42000 …line 1 column 18 near "trailing = 1" 
> exec ^
: string
| set local transaction = 1
< ok 0 0
> exec ^
: string
| set local trigger = 1
< error
| 1064 42000 …line 1 column 17 near "trigger = 1" 
> exec ^
: string
| set local triggers = 1
< ok 0 0
> exec ^
: string
| set local trim = 1
< ok 0 0
> exec ^
: string
| set local true = 1
< error
| 1064 42000 …line 1 column 14 near "true = 1" 
> exec ^
: string
| set local truncate = 1
< ok 0 0
> exec ^
: string
| set local ttl = 1
< ok 0 0
> exec ^
: string
| set local ttl_enable = 1
< ok 0 0
> exec ^
: string
| set local ttl_job_interval = 1
< ok 0 0
> exec ^
: string
| set local type = 1
< ok 0 0
> exec ^
: string
| set local types = 1
< ok 0 0
> exec ^
: string
| set local unbounded = 1
< ok 0 0
> exec ^
: string
| set local uncommitted = 1
< ok 0 0
> exec ^
: string
| set local undefined = 1
< ok 0 0
> exec ^
: string
| set local undo = 1
< ok 0 0
> exec ^
: string
| set local undo_buffer_size = 1
< ok 0 0
> exec ^
: string
| set local undofile = 1
< ok 0 0
> exec ^
: string
| set local unicode = 1
< ok 0 0
> exec ^
: string
| set local uninstall = 1
< ok 0 0
> exec ^
: string
| set local union = 1
< error
| 1064 42000 …line 1 column 15 near "union = 1" 
> exec ^
: string
| set local unique = 1
< error
| 1064 42000 …line 1 column 16 near "unique = 1" 
> exec ^
: string
| set local unknown = 1
< ok 0 0
> exec ^
: string
| set local unlimited = 1
< ok 0 0
> exec ^
: string
| set local unlock = 1
< error
| 1064 42000 …line 1 column 16 near "unlock = 1" 
> exec ^
: string
| set local unregister = 1
< ok 0 0
> exec ^
: string
| set local unsigned = 1
< error
| 1064 42000 …line 1 column 18 near "unsigned = 1" 
> exec ^
: string
| set local until = 1
< error
| 1064 42000 …line 1 column 15 near "until = 1" 
> exec ^
: string
| set local update = 1
< error
| 1064 42000 …line 1 column 16 near "update = 1" 
> exec ^
: string
| set local upgrade = 1
< ok 0 0
> exec ^
: string
| set local uri = 1
< ok 0 0
> exec ^
: string
| set local url = 1
< ok 0 0
> exec ^
: string
| set local usage = 1
< error
| 1064 42000 …line 1 column 15 near "usage = 1" 
> exec ^
: string
| set local use = 1
< error
| 1064 42000 …line 1 column 13 near "use = 1" 
> exec ^
: string
| set local use_frm = 1
< ok 0 0
> exec ^
: string
| set local user = 1
< ok 0 0
> exec ^
: string
| set local user_resources = 1
< ok 0 0
> exec ^
: string
| set local using = 1
< error
| 1064 42000 …line 1 column 15 near "using = 1" 
> exec ^
: string
| set local utc_date = 1
< error
| 1064 42000 …line 1 column 18 near "utc_date = 1" 
> exec ^
: string
| set local utc_time = 1
< error
| 1064 42000 …line 1 column 18 near "utc_time = 1" 
> exec ^
: string
| set local utc_timestamp = 1
< error
| 1064 42000 …line 1 column 23 near "utc_timestamp = 1" 
> exec ^
: string
| set local utilization_limit = 1
< ok 0 0
> exec ^
: string
| set local validate = 1
< ok 0 0
> exec ^
: string
| set local validation = 1
< ok 0 0
> exec ^
: string
| set local value = 1
< ok 0 0
> exec ^
: string
| set local values = 1
< error
| 1064 42000 …line 1 column 16 near "values = 1" 
> exec ^
: string
| set local var_pop = 1
< ok 0 0
> exec ^
: string
| set local var_samp = 1
< ok 0 0
> exec ^
: string
| set local varbinary = 1
< error
| 1064 42000 …line 1 column 19 near "varbinary = 1" 
> exec ^
: string
| set local varchar = 1
< error
| 1064 42000 …line 1 column 17 near "varchar = 1" 
> exec ^
: string
| set local varcharacter = 1
< error
| 1064 42000 …line 1 column 22 near "varcharacter = 1" 
> exec ^
: string
| set local variables = 1
< ok 0 0
> exec ^
: string
| set local variance = 1
< ok 0 0
> exec ^
: string
| set local varying = 1
< error
| 1064 42000 …line 1 column 17 near "varying = 1" 
> exec ^
: string
| set local vcpu = 1
< ok 0 0
> exec ^
: string
| set local vector = 1
< ok 0 0
> exec ^
: string
| set local verify_key_constraints = 1
< ok 0 0
> exec ^
: string
| set local view = 1
< ok 0 0
> exec ^
: string
| set local virtual = 1
< error
| 1064 42000 …line 1 column 17 near "virtual = 1" 
> exec ^
: string
| set local visible = 1
< ok 0 0
> exec ^
: string
| set local voter_constraints = 1
< ok 0 0
> exec ^
: string
| set local voters = 1
< ok 0 0
> exec ^
: string
| set local wait = 1
< ok 0 0
> exec ^
: string
| set local warnings = 1
< ok 0 0
> exec ^
: string
| set local watch = 1
< ok 0 0
> exec ^
: string
| set local week = 1
< ok 0 0
> exec ^
: string
| set local weight_string = 1
< ok 0 0
> exec ^
: string
| set local when = 1
< error
| 1064 42000 …line 1 column 14 near "when = 1" 
> exec ^
: string
| set local where = 1
< error
| 1064 42000 …line 1 column 15 near "where = 1" 
> exec ^
: string
| set local while = 1
< error
| 1064 42000 …line 1 column 15 near "while = 1" 
> exec ^
: string
| set local window = 1
< error
| 1064 42000 …line 1 column 16 near "window = 1" 
> exec ^
: string
| set local with = 1
< error
| 1064 42000 …line 1 column 14 near "with = 1" 
> exec ^
: string
| set local without = 1
< ok 0 0
> exec ^
: string
| set local work = 1
< ok 0 0
> exec ^
: string
| set local workload = 1
< ok 0 0
> exec ^
: string
| set local wrapper = 1
< ok 0 0
> exec ^
: string
| set local write = 1
< error
| 1064 42000 …line 1 column 15 near "write = 1" 
> exec ^
: string
| set local x509 = 1
< ok 0 0
> exec ^
: string
| set local xa = 1
< ok 0 0
> exec ^
: string
| set local xid = 1
< ok 0 0
> exec ^
: string
| set local xml = 1
< ok 0 0
> exec ^
: string
| set local xor = 1
< error
| 1064 42000 …line 1 column 13 near "xor = 1" 
> exec ^
: string
| set local year = 1
< ok 0 0
> exec ^
: string
| set local year_month = 1
< error
| 1064 42000 …line 1 column 20 near "year_month = 1" 
> exec ^
: string
| set local zerofill = 1
< error
| 1064 42000 …line 1 column 18 near "zerofill = 1" 
> exec ^
: string
| set local zone = 1
< ok 0 0
> exec ^
: string
| set absent.a = 1
< ok 0 0
> exec ^
: string
| set accessible.a = 1
< ok 0 0
> exec ^
: string
| set account.a = 1
< ok 0 0
> exec ^
: string
| set action.a = 1
< ok 0 0
> exec ^
: string
| set active.a = 1
< ok 0 0
> exec ^
: string
| set add.a = 1
< ok 0 0
> exec ^
: string
| set adddate.a = 1
< ok 0 0
> exec ^
: string
| set admin.a = 1
< ok 0 0
> exec ^
: string
| set after.a = 1
< ok 0 0
> exec ^
: string
| set against.a = 1
< ok 0 0
> exec ^
: string
| set aggregate.a = 1
< ok 0 0
> exec ^
: string
| set algorithm.a = 1
< ok 0 0
> exec ^
: string
| set all.a = 1
< ok 0 0
> exec ^
: string
| set allow_missing_files.a = 1
< ok 0 0
> exec ^
: string
| set alter.a = 1
< ok 0 0
> exec ^
: string
| set always.a = 1
< ok 0 0
> exec ^
: string
| set analyze.a = 1
< ok 0 0
> exec ^
: string
| set and.a = 1
< ok 0 0
> exec ^
: string
| set any.a = 1
< ok 0 0
> exec ^
: string
| set array.a = 1
< ok 0 0
> exec ^
: string
| set as.a = 1
< ok 0 0
> exec ^
: string
| set asc.a = 1
< ok 0 0
> exec ^
: string
| set ascii.a = 1
< ok 0 0
> exec ^
: string
| set asensitive.a = 1
< ok 0 0
> exec ^
: string
| set assign_gtids_to_anonymous_transactions.a = 1
< ok 0 0
> exec ^
: string
| set at.a = 1
< ok 0 0
> exec ^
: string
| set attribute.a = 1
< ok 0 0
> exec ^
: string
| set authentication.a = 1
< ok 0 0
> exec ^
: string
| set auto.a = 1
< ok 0 0
> exec ^
: string
| set auto_id_cache.a = 1
< ok 0 0
> exec ^
: string
| set auto_increment.a = 1
< ok 0 0
> exec ^
: string
| set auto_random.a = 1
< ok 0 0
> exec ^
: string
| set auto_random_base.a = 1
< ok 0 0
> exec ^
: string
| set auto_refresh.a = 1
< ok 0 0
> exec ^
: string
| set auto_refresh_source.a = 1
< ok 0 0
> exec ^
: string
| set autoextend_size.a = 1
< ok 0 0
> exec ^
: string
| set avg.a = 1
< ok 0 0
> exec ^
: string
| set avg_row_length.a = 1
< ok 0 0
> exec ^
: string
| set background.a = 1
< ok 0 0
> exec ^
: string
| set backup.a = 1
< ok 0 0
> exec ^
: string
| set backups.a = 1
< ok 0 0
> exec ^
: string
| set batch.a = 1
< ok 0 0
> exec ^
: string
| set before.a = 1
< ok 0 0
> exec ^
: string
| set begin.a = 1
< ok 0 0
> exec ^
: string
| set bernoulli.a = 1
< ok 0 0
> exec ^
: string
| set between.a = 1
< ok 0 0
> exec ^
: string
| set bigint.a = 1
< ok 0 0
> exec ^
: string
| set binary.a = 1
< ok 0 0
> exec ^
: string
| set binding.a = 1
< ok 0 0
> exec ^
: string
| set binding_cache.a = 1
< ok 0 0
> exec ^
: string
| set bindings.a = 1
< ok 0 0
> exec ^
: string
| set binlog.a = 1
< ok 0 0
> exec ^
: string
| set bit.a = 1
< ok 0 0
> exec ^
: string
| set bit_and.a = 1
< ok 0 0
> exec ^
: string
| set bit_or.a = 1
< ok 0 0
> exec ^
: string
| set bit_xor.a = 1
< ok 0 0
> exec ^
: string
| set blob.a = 1
< ok 0 0
> exec ^
: string
| set block.a = 1
< ok 0 0
> exec ^
: string
| set bool.a = 1
< ok 0 0
> exec ^
: string
| set boolean.a = 1
< ok 0 0
> exec ^
: string
| set both.a = 1
< ok 0 0
> exec ^
: string
| set btree.a = 1
< ok 0 0
> exec ^
: string
| set buckets.a = 1
< ok 0 0
> exec ^
: string
| set builtins.a = 1
< ok 0 0
> exec ^
: string
| set bulk.a = 1
< ok 0 0
> exec ^
: string
| set burstable.a = 1
< ok 0 0
> exec ^
: string
| set by.a = 1
< ok 0 0
> exec ^
: string
| set byte.a = 1
< ok 0 0
> exec ^
: string
| set cache.a = 1
< ok 0 0
> exec ^
: string
| set calibrate.a = 1
< ok 0 0
> exec ^
: string
| set call.a = 1
< ok 0 0
> exec ^
: string
| set cascade.a = 1
< ok 0 0
> exec ^
: string
| set cascaded.a = 1
< ok 0 0
> exec ^
: string
| set case.a = 1
< ok 0 0
> exec ^
: string
| set cast.a = 1
< ok 0 0
> exec ^
: string
| set catalog_name.a = 1
< ok 0 0
> exec ^
: string
| set chain.a = 1
< ok 0 0
> exec ^
: string
| set challenge_response.a = 1
< ok 0 0
> exec ^
: string
| set change.a = 1
< ok 0 0
> exec ^
: string
| set changed.a = 1
< ok 0 0
> exec ^
: string
| set channel.a = 1
< ok 0 0
> exec ^
: string
| set char.a = 1
< ok 0 0
> exec ^
: string
| set character.a = 1
< ok 0 0
> exec ^
: string
| set charset.a = 1
< ok 0 0
> exec ^
: string
| set check.a = 1
< ok 0 0
> exec ^
: string
| set checksum.a = 1
< ok 0 0
> exec ^
: string
| set cipher.a = 1
< ok 0 0
> exec ^
: string
| set class_origin.a = 1
< ok 0 0
> exec ^
: string
| set client.a = 1
< ok 0 0
> exec ^
: string
| set clone.a = 1
< ok 0 0
> exec ^
: string
| set close.a = 1
< ok 0 0
> exec ^
: string
| set clustered.a = 1
< ok 0 0
> exec ^
: string
| set coalesce.a = 1
< ok 0 0
> exec ^
: string
| set code.a = 1
< ok 0 0
> exec ^
: string
| set collate.a = 1
< ok 0 0
> exec ^
: string
| set collation.a = 1
< ok 0 0
> exec ^
: string
| set column.a = 1
< ok 0 0
> exec ^
: string
| set column_format.a = 1
< ok 0 0
> exec ^
: string
| set column_name.a = 1
< ok 0 0
> exec ^
: string
| set column_stats_usage.a = 1
< ok 0 0
> exec ^
: string
| set columns.a = 1
< ok 0 0
> exec ^
: string
| set comment.a = 1
< ok 0 0
> exec ^
: string
| set commit.a = 1
< ok 0 0
> exec ^
: string
| set committed.a = 1
< ok 0 0
> exec ^
: string
| set compact.a = 1
< ok 0 0
> exec ^
: string
| set completion.a = 1
< ok 0 0
> exec ^
: string
| set component.a = 1
< ok 0 0
> exec ^
: string
| set compressed.a = 1
< ok 0 0
> exec ^
: string
| set compression.a = 1
< ok 0 0
> exec ^
: string
| set concat.a = 1
< ok 0 0
> exec ^
: string
| set concurrent.a = 1
< ok 0 0
> exec ^
: string
| set condition.a = 1
< ok 0 0
> exec ^
: string
| set config.a = 1
< ok 0 0
> exec ^
: string
| set connection.a = 1
< ok 0 0
> exec ^
: string
| set consistent.a = 1
< ok 0 0
> exec ^
: string
| set constraint.a = 1
< ok 0 0
> exec ^
: string
| set constraint_catalog.a = 1
< ok 0 0
> exec ^
: string
| set constraint_name.a = 1
< ok 0 0
> exec ^
: string
| set constraint_schema.a = 1
< ok 0 0
> exec ^
: string
| set constraints.a = 1
< ok 0 0
> exec ^
: string
| set contains.a = 1
< ok 0 0
> exec ^
: string
| set context.a = 1
< ok 0 0
> exec ^
: string
| set continue.a = 1
< ok 0 0
> exec ^
: string
| set convert.a = 1
< ok 0 0
> exec ^
: string
| set cooldown.a = 1
< ok 0 0
> exec ^
: string
| set count.a = 1
< ok 0 0
> exec ^
: string
| set cpu.a = 1
< ok 0 0
> exec ^
: string
| set create.a = 1
< ok 0 0
> exec ^
: string
| set cross.a = 1
< ok 0 0
> exec ^
: string
| set cube.a = 1
< ok 0 0
> exec ^
: string
| set cume_dist.a = 1
< ok 0 0
> exec ^
: string
| set curdate.a = 1
< ok 0 0
> exec ^
: string
| set current.a = 1
< ok 0 0
> exec ^
: string
| set current_date.a = 1
< ok 0 0
> exec ^
: string
| set current_time.a = 1
< ok 0 0
> exec ^
: string
| set current_timestamp.a = 1
< ok 0 0
> exec ^
: string
| set current_user.a = 1
< ok 0 0
> exec ^
: string
| set cursor.a = 1
< ok 0 0
> exec ^
: string
| set cursor_name.a = 1
< ok 0 0
> exec ^
: string
| set curtime.a = 1
< ok 0 0
> exec ^
: string
| set cycle.a = 1
< ok 0 0
> exec ^
: string
| set data.a = 1
< ok 0 0
> exec ^
: string
| set database.a = 1
< ok 0 0
> exec ^
: string
| set databases.a = 1
< ok 0 0
> exec ^
: string
| set datafile.a = 1
< ok 0 0
> exec ^
: string
| set date.a = 1
< ok 0 0
> exec ^
: string
| set date_add.a = 1
< ok 0 0
> exec ^
: string
| set date_sub.a = 1
< ok 0 0
> exec ^
: string
| set datetime.a = 1
< ok 0 0
> exec ^
: string
| set day.a = 1
< ok 0 0
> exec ^
: string
| set day_hour.a = 1
< ok 0 0
> exec ^
: string
| set day_microsecond.a = 1
< ok 0 0
> exec ^
: string
| set day_minute.a = 1
< ok 0 0
> exec ^
: string
| set day_second.a = 1
< ok 0 0
> exec ^
: string
| set deallocate.a = 1
< ok 0 0
> exec ^
: string
| set dec.a = 1
< ok 0 0
> exec ^
: string
| set decimal.a = 1
< ok 0 0
> exec ^
: string
| set declare.a = 1
< ok 0 0
> exec ^
: string
| set default.a = 1
< ok 0 0
> exec ^
: string
| set default_auth.a = 1
< ok 0 0
> exec ^
: string
| set defined.a = 1
< ok 0 0
> exec ^
: string
| set definer.a = 1
< ok 0 0
> exec ^
: string
| set definition.a = 1
< ok 0 0
> exec ^
: string
| set delay_key_write.a = 1
< ok 0 0
> exec ^
: string
| set delayed.a = 1
< ok 0 0
> exec ^
: string
| set delete.a = 1
< ok 0 0
> exec ^
: string
| set dense_rank.a = 1
< ok 0 0
> exec ^
: string
| set desc.a = 1
< ok 0 0
> exec ^
: string
| set describe.a = 1
< ok 0 0
> exec ^
: string
| set description.a = 1
< ok 0 0
> exec ^
: string
| set deterministic.a = 1
< ok 0 0
> exec ^
: string
| set diagnostics.a = 1
< ok 0 0
> exec ^
: string
| set digest.a = 1
< ok 0 0
> exec ^
: string
| set directory.a = 1
< ok 0 0
> exec ^
: string
| set disable.a = 1
< ok 0 0
> exec ^
: string
| set disabled.a = 1
< ok 0 0
> exec ^
: string
| set discard.a = 1
< ok 0 0
> exec ^
: string
| set disk.a = 1
< ok 0 0
> exec ^
: string
| set distinct.a = 1
< ok 0 0
> exec ^
: string
| set distinctrow.a = 1
< ok 0 0
> exec ^
: string
| set distribute.a = 1
< ok 0 0
> exec ^
: string
| set distribution.a = 1
< ok 0 0
> exec ^
: string
| set distributions.a = 1
< ok 0 0
> exec ^
: string
| set div.a = 1
< ok 0 0
> exec ^
: string
| set do.a = 1
< ok 0 0
> exec ^
: string
| set double.a = 1
< ok 0 0
> exec ^
: string
| set drop.a = 1
< ok 0 0
> exec ^
: string
| set dry.a = 1
< ok 0 0
> exec ^
: string
| set dryrun.a = 1
< ok 0 0
> exec ^
: string
| set dual.a = 1
< ok 0 0
> exec ^
: string
| set duality.a = 1
< ok 0 0
> exec ^
: string
| set dumpfile.a = 1
< ok 0 0
> exec ^
: string
| set duplicate.a = 1
< ok 0 0
> exec ^
: string
| set duration.a = 1
< ok 0 0
> exec ^
: string
| set dynamic.a = 1
< ok 0 0
> exec ^
: string
| set each.a = 1
< ok 0 0
> exec ^
: string
| set else.a = 1
< ok 0 0
> exec ^
: string
| set elseif.a = 1
< ok 0 0
> exec ^
: string
| set empty.a = 1
< ok 0 0
> exec ^
: string
| set enable.a = 1
< ok 0 0
> exec ^
: string
| set enabled.a = 1
< ok 0 0
> exec ^
: string
| set enclosed.a = 1
< ok 0 0
> exec ^
: string
| set encryption.a = 1
< ok 0 0
> exec ^
: string
| set end.a = 1
< ok 0 0
> exec ^
: string
| set end_time.a = 1
< ok 0 0
> exec ^
: string
| set ends.a = 1
< ok 0 0
> exec ^
: string
| set enforced.a = 1
< ok 0 0
> exec ^
: string
| set engine.a = 1
< ok 0 0
> exec ^
: string
| set engine_attribute.a = 1
< ok 0 0
> exec ^
: string
| set engines.a = 1
< ok 0 0
> exec ^
: string
| set enum.a = 1
< ok 0 0
> exec ^
: string
| set error.a = 1
< ok 0 0
> exec ^
: string
| set errors.a = 1
< ok 0 0
> exec ^
: string
| set escape.a = 1
< ok 0 0
> exec ^
: string
| set escaped.a = 1
< ok 0 0
> exec ^
: string
| set event.a = 1
< ok 0 0
> exec ^
: string
| set events.a = 1
< ok 0 0
> exec ^
: string
| set every.a = 1
< ok 0 0
> exec ^
: string
| set exact.a = 1
< ok 0 0
> exec ^
: string
| set except.a = 1
< ok 0 0
> exec ^
: string
| set exchange.a = 1
< ok 0 0
> exec ^
: string
| set exclude.a = 1
< ok 0 0
> exec ^
: string
| set exec_elapsed.a = 1
< ok 0 0
> exec ^
: string
| set execute.a = 1
< ok 0 0
> exec ^
: string
| set exists.a = 1
< ok 0 0
> exec ^
: string
| set exit.a = 1
< ok 0 0
> exec ^
: string
| set expansion.a = 1
< ok 0 0
> exec ^
: string
| set expire.a = 1
< ok 0 0
> exec ^
: string
| set explain.a = 1
< ok 0 0
> exec ^
: string
| set export.a = 1
< ok 0 0
> exec ^
: string
| set extended.a = 1
< ok 0 0
> exec ^
: string
| set extent_size.a = 1
< ok 0 0
> exec ^
: string
| set external.a = 1
< ok 0 0
> exec ^
: string
| set external_format.a = 1
< ok 0 0
> exec ^
: string
| set extract.a = 1
< ok 0 0
> exec ^
: string
| set factor.a = 1
< ok 0 0
> exec ^
: string
| set failed_login_attempts.a = 1
< ok 0 0
> exec ^
: string
| set false.a = 1
< ok 0 0
> exec ^
: string
| set fast.a = 1
< ok 0 0
> exec ^
: string
| set faults.a = 1
< ok 0 0
> exec ^
: string
| set fetch.a = 1
< ok 0 0
> exec ^
: string
| set fields.a = 1
< ok 0 0
> exec ^
: string
| set file.a = 1
< ok 0 0
> exec ^
: string
| set file_block_size.a = 1
< ok 0 0
> exec ^
: string
| set file_format.a = 1
< ok 0 0
> exec ^
: string
| set file_name.a = 1
< ok 0 0
> exec ^
: string
| set file_pattern.a = 1
< ok 0 0
> exec ^
: string
| set file_prefix.a = 1
< ok 0 0
> exec ^
: string
| set files.a = 1
< ok 0 0
> exec ^
: string
| set filter.a = 1
< ok 0 0
> exec ^
: string
| set finish.a = 1
< ok 0 0
> exec ^
: string
| set first.a = 1
< ok 0 0
> exec ^
: string
| set first_value.a = 1
< ok 0 0
> exec ^
: string
| set fixed.a = 1
< ok 0 0
> exec ^
: string
| set flashback.a = 1
< ok 0 0
> exec ^
: string
| set float.a = 1
< ok 0 0
> exec ^
: string
| set float4.a = 1
< ok 0 0
> exec ^
: string
| set float8.a = 1
< ok 0 0
> exec ^
: string
| set flush.a = 1
< ok 0 0
> exec ^
: string
| set follower_constraints.a = 1
< ok 0 0
> exec ^
: string
| set followers.a = 1
< ok 0 0
> exec ^
: string
| set following.a = 1
< ok 0 0
> exec ^
: string
| set follows.a = 1
< ok 0 0
> exec ^
: string
| set for.a = 1
< ok 0 0
> exec ^
: string
| set force.a = 1
< ok 0 0
> exec ^
: string
| set foreign.a = 1
< ok 0 0
> exec ^
: string
| set format.a = 1
< ok 0 0
> exec ^
: string
| set found.a = 1
< ok 0 0
> exec ^
: string
| set from.a = 1
< ok 0 0
> exec ^
: string
| set full.a = 1
< ok 0 0
> exec ^
: string
| set fulltext.a = 1
< ok 0 0
> exec ^
: string
| set function.a = 1
< ok 0 0
> exec ^
: string
| set general.a = 1
< ok 0 0
> exec ^
: string
| set generate.a = 1
< ok 0 0
> exec ^
: string
| set generated.a = 1
< ok 0 0
> exec ^
: string
| set geomcollection.a = 1
< ok 0 0
> exec ^
: string
| set geometry.a = 1
< ok 0 0
> exec ^
: string
| set geometrycollection.a = 1
< ok 0 0
> exec ^
: string
| set get.a = 1
< ok 0 0
> exec ^
: string
| set get_format.a = 1
< ok 0 0
> exec ^
: string
| set get_source_public_key.a = 1
< ok 0 0
> exec ^
: string
| set global.a = 1
< ok 0 0
> exec ^
: string
| set grant.a = 1
< ok 0 0
> exec ^
: string
| set grants.a = 1
< ok 0 0
> exec ^
: string
| set greatest.a = 1
< ok 0 0
> exec ^
: string
| set group.a = 1
< ok 0 0
> exec ^
: string
| set group_concat.a = 1
< ok 0 0
> exec ^
: string
| set group_replication.a = 1
< ok 0 0
> exec ^
: string
| set grouping.a = 1
< ok 0 0
> exec ^
: string
| set groups.a = 1
< ok 0 0
> exec ^
: string
| set gtid_only.a = 1
< ok 0 0
> exec ^
: string
| set gtids.a = 1
< ok 0 0
> exec ^
: string
| set guided.a = 1
< ok 0 0
> exec ^
: string
| set handler.a = 1
< ok 0 0
> exec ^
: string
| set hash.a = 1
< ok 0 0
> exec ^
: string
| set having.a = 1
< ok 0 0
> exec ^
: string
| set header.a = 1
< ok 0 0
> exec ^
: string
| set help.a = 1
< ok 0 0
> exec ^
: string
| set high.a = 1
< ok 0 0
> exec ^
: string
| set high_priority.a = 1
< ok 0 0
> exec ^
: string
| set histogram.a = 1
< ok 0 0
> exec ^
: string
| set histograms_in_flight.a = 1
< ok 0 0
> exec ^
: string
| set history.a = 1
< ok 0 0
> exec ^
: string
| set host.a = 1
< ok 0 0
> exec ^
: string
| set hosts.a = 1
< ok 0 0
> exec ^
: string
| set hour.a = 1
< ok 0 0
> exec ^
: string
| set hour_microsecond.a = 1
< ok 0 0
> exec ^
: string
| set hour_minute.a = 1
< ok 0 0
> exec ^
: string
| set hour_second.a = 1
< ok 0 0
> exec ^
: string
| set identified.a = 1
< ok 0 0
> exec ^
: string
| set if.a = 1
< ok 0 0
> exec ^
: string
| set ignore.a = 1
< ok 0 0
> exec ^
: string
| set ignore_server_ids.a = 1
< ok 0 0
> exec ^
: string
| set ilike.a = 1
< ok 0 0
> exec ^
: string
| set import.a = 1
< ok 0 0
> exec ^
: string
| set in.a = 1
< ok 0 0
> exec ^
: string
| set inactive.a = 1
< ok 0 0
> exec ^
: string
| set increment.a = 1
< ok 0 0
> exec ^
: string
| set index.a = 1
< ok 0 0
> exec ^
: string
| set indexes.a = 1
< ok 0 0
> exec ^
: string
| set infile.a = 1
< ok 0 0
> exec ^
: string
| set initial.a = 1
< ok 0 0
> exec ^
: string
| set initial_size.a = 1
< ok 0 0
> exec ^
: string
| set initiate.a = 1
< ok 0 0
> exec ^
: string
| set inner.a = 1
< ok 0 0
> exec ^
: string
| set inout.a = 1
< ok 0 0
> exec ^
: string
| set insensitive.a = 1
< ok 0 0
> exec ^
: string
| set insert.a = 1
< ok 0 0
> exec ^
: string
| set insert_method.a = 1
< ok 0 0
> exec ^
: string
| set install.a = 1
< ok 0 0
> exec ^
: string
| set instance.a = 1
< ok 0 0
> exec ^
: string
| set int.a = 1
< ok 0 0
> exec ^
: string
| set int1.a = 1
< ok 0 0
> exec ^
: string
| set int2.a = 1
< ok 0 0
> exec ^
: string
| set int3.a = 1
< ok 0 0
> exec ^
: string
| set int4.a = 1
< ok 0 0
> exec ^
: string
| set int8.a = 1
< ok 0 0
> exec ^
: string
| set integer.a = 1
< ok 0 0
> exec ^
: string
| set intersect.a = 1
< ok 0 0
> exec ^
: string
| set interval.a = 1
< ok 0 0
> exec ^
: string
| set into.a = 1
< ok 0 0
> exec ^
: string
| set invisible.a = 1
< ok 0 0
> exec ^
: string
| set invoker.a = 1
< ok 0 0
> exec ^
: string
| set io.a = 1
< ok 0 0
> exec ^
: string
| set io_after_gtids.a = 1
< ok 0 0
> exec ^
: string
| set io_before_gtids.a = 1
< ok 0 0
> exec ^
: string
| set io_thread.a = 1
< ok 0 0
> exec ^
: string
| set ipc.a = 1
< ok 0 0
> exec ^
: string
| set is.a = 1
< ok 0 0
> exec ^
: string
| set isolation.a = 1
< ok 0 0
> exec ^
: string
| set issuer.a = 1
< ok 0 0
> exec ^
: string
| set iterate.a = 1
< ok 0 0
> exec ^
: string
| set job.a = 1
< ok 0 0
> exec ^
: string
| set jobs.a = 1
< ok 0 0
> exec ^
: string
| set join.a = 1
< ok 0 0
> exec ^
: string
| set json.a = 1
< ok 0 0
> exec ^
: string
| set json_arrayagg.a = 1
< ok 0 0
> exec ^
: string
| set json_duality_object.a = 1
< ok 0 0
> exec ^
: string
| set json_objectagg.a = 1
< ok 0 0
> exec ^
: string
| set json_table.a = 1
< ok 0 0
> exec ^
: string
| set json_value.a = 1
< ok 0 0
> exec ^
: string
| set key.a = 1
< ok 0 0
> exec ^
: string
| set key_block_size.a = 1
< ok 0 0
> exec ^
: string
| set keyring.a = 1
< ok 0 0
> exec ^
: string
| set keys.a = 1
< ok 0 0
> exec ^
: string
| set kill.a = 1
< ok 0 0
> exec ^
: string
| set labels.a = 1
< ok 0 0
> exec ^
: string
| set lag.a = 1
< ok 0 0
> exec ^
: string
| set language.a = 1
< ok 0 0
> exec ^
: string
| set last.a = 1
< ok 0 0
> exec ^
: string
| set last_value.a = 1
< ok 0 0
> exec ^
: string
| set lateral.a = 1
< ok 0 0
> exec ^
: string
| set lead.a = 1
< ok 0 0
> exec ^
: string
| set leader_constraints.a = 1
< ok 0 0
> exec ^
: string
| set leading.a = 1
< ok 0 0
> exec ^
: string
| set learner_constraints.a = 1
< ok 0 0
> exec ^
: string
| set learners.a = 1
< ok 0 0
> exec ^
: string
| set least.a = 1
< ok 0 0
> exec ^
: string
| set leave.a = 1
< ok 0 0
> exec ^
: string
| set leaves.a = 1
< ok 0 0
> exec ^
: string
| set left.a = 1
< ok 0 0
> exec ^
: string
| set less.a = 1
< ok 0 0
> exec ^
: string
| set level.a = 1
< ok 0 0
> exec ^
: string
| set library.a = 1
< ok 0 0
> exec ^
: string
| set like.a = 1
< ok 0 0
> exec ^
: string
| set limit.a = 1
< ok 0 0
> exec ^
: string
| set linear.a = 1
< ok 0 0
> exec ^
: string
| set lines.a = 1
< ok 0 0
> exec ^
: string
| set linestring.a = 1
< ok 0 0
> exec ^
: string
| set list.a = 1
< ok 0 0
> exec ^
: string
| set load.a = 1
< ok 0 0
> exec ^
: string
| set local.a = 1
< ok 0 0
> exec ^
: string
| set localtime.a = 1
< ok 0 0
> exec ^
: string
| set localtimestamp.a = 1
< ok 0 0
> exec ^
: string
| set location.a = 1
< ok 0 0
> exec ^
: string
| set lock.a = 1
< ok 0 0
> exec ^
: string
| set locked.a = 1
< ok 0 0
> exec ^
: string
| set locks.a = 1
< ok 0 0
> exec ^
: string
| set log.a = 1
< ok 0 0
> exec ^
: string
| set logfile.a = 1
< ok 0 0
> exec ^
: string
| set logs.a = 1
< ok 0 0
> exec ^
: string
| set long.a = 1
< ok 0 0
> exec ^
: string
| set longblob.a = 1
< ok 0 0
> exec ^
: string
| set longtext.a = 1
< ok 0 0
> exec ^
: string
| set loop.a = 1
< ok 0 0
> exec ^
: string
| set low.a = 1
< ok 0 0
> exec ^
: string
| set low_priority.a = 1
< ok 0 0
> exec ^
: string
| set manual.a = 1
< ok 0 0
> exec ^
: string
| set master.a = 1
< ok 0 0
> exec ^
: string
| set master_log_file.a = 1
< ok 0 0
> exec ^
: string
| set match.a = 1
< ok 0 0
> exec ^
: string
| set materialized.a = 1
< ok 0 0
> exec ^
: string
| set max.a = 1
< ok 0 0
> exec ^
: string
| set max_connections_per_hour.a = 1
< ok 0 0
> exec ^
: string
| set max_queries_per_hour.a = 1
< ok 0 0
> exec ^
: string
| set max_rows.a = 1
< ok 0 0
> exec ^
: string
| set max_size.a = 1
< ok 0 0
> exec ^
: string
| set max_updates_per_hour.a = 1
< ok 0 0
> exec ^
: string
| set max_user_connections.a = 1
< ok 0 0
> exec ^
: string
| set maxvalue.a = 1
< ok 0 0
> exec ^
: string
| set medium.a = 1
< ok 0 0
> exec ^
: string
| set mediumblob.a = 1
< ok 0 0
> exec ^
: string
| set mediumint.a = 1
< ok 0 0
> exec ^
: string
| set mediumtext.a = 1
< ok 0 0
> exec ^
: string
| set member.a = 1
< ok 0 0
> exec ^
: string
| set memory.a = 1
< ok 0 0
> exec ^
: string
| set merge.a = 1
< ok 0 0
> exec ^
: string
| set message_text.a = 1
< ok 0 0
> exec ^
: string
| set microsecond.a = 1
< ok 0 0
> exec ^
: string
| set mid.a = 1
< ok 0 0
> exec ^
: string
| set middleint.a = 1
< ok 0 0
> exec ^
: string
| set migrate.a = 1
< ok 0 0
> exec ^
: string
| set min.a = 1
< ok 0 0
> exec ^
: string
| set min_rows.a = 1
< ok 0 0
> exec ^
: string
| set minute.a = 1
< ok 0 0
> exec ^
: string
| set minute_microsecond.a = 1
< ok 0 0
> exec ^
: string
| set minute_second.a = 1
< ok 0 0
> exec ^
: string
| set minvalue.a = 1
< ok 0 0
> exec ^
: string
| set mod.a = 1
< ok 0 0
> exec ^
: string
| set mode.a = 1
< ok 0 0
> exec ^
: string
| set modifies.a = 1
< ok 0 0
> exec ^
: string
| set modify.a = 1
< ok 0 0
> exec ^
: string
| set month.a = 1
< ok 0 0
> exec ^
: string
| set multilinestring.a = 1
< ok 0 0
> exec ^
: string
| set multipoint.a = 1
< ok 0 0
> exec ^
: string
| set multipolygon.a = 1
< ok 0 0
> exec ^
: string
| set mutex.a = 1
< ok 0 0
> exec ^
: string
| set mysql_errno.a = 1
< ok 0 0
> exec ^
: string
| set name.a = 1
< ok 0 0
> exec ^
: string
| set names.a = 1
< ok 0 0
> exec ^
: string
| set national.a = 1
< ok 0 0
> exec ^
: string
| set natural.a = 1
< ok 0 0
> exec ^
: string
| set nchar.a = 1
< ok 0 0
> exec ^
: string
| set ndb.a = 1
< ok 0 0
> exec ^
: string
| set ndbcluster.a = 1
< ok 0 0
> exec ^
: string
| set nested.a = 1
< ok 0 0
> exec ^
: string
| set network_namespace.a = 1
< ok 0 0
> exec ^
: string
| set never.a = 1
< ok 0 0
> exec ^
: string
| set new.a = 1
< ok 0 0
> exec ^
: string
| set next.a = 1
< ok 0 0
> exec ^
: string
| set next_row_id.a = 1
< ok 0 0
> exec ^
: string
| set no.a = 1
< ok 0 0
> exec ^
: string
| set no_wait.a = 1
< ok 0 0
> exec ^
: string
| set no_write_to_binlog.a = 1
< ok 0 0
> exec ^
: string
| set nocache.a = 1
< ok 0 0
> exec ^
: string
| set nocycle.a = 1
< ok 0 0
> exec ^
: string
| set nodegroup.a = 1
< ok 0 0
> exec ^
: string
| set nomaxvalue.a = 1
< ok 0 0
> exec ^
: string
| set nominvalue.a = 1
< ok 0 0
> exec ^
: string
| set nonclustered.a = 1
< ok 0 0
> exec ^
: string
| set none.a = 1
< ok 0 0
> exec ^
: string
| set not.a = 1
< ok 0 0
> exec ^
: string
| set now.a = 1
< ok 0 0
> exec ^
: string
| set nowait.a = 1
< ok 0 0
> exec ^
: string
| set nth_value.a = 1
< ok 0 0
> exec ^
: string
| set ntile.a = 1
< ok 0 0
> exec ^
: string
| set null.a = 1
< ok 0 0
> exec ^
: string
| set nullif.a = 1
< ok 0 0
> exec ^
: string
| set nulls.a = 1
< ok 0 0
> exec ^
: string
| set number.a = 1
< ok 0 0
> exec ^
: string
| set numeric.a = 1
< ok 0 0
> exec ^
: string
| set nvarchar.a = 1
< ok 0 0
> exec ^
: string
| set of.a = 1
< ok 0 0
> exec ^
: string
| set off.a = 1
< ok 0 0
> exec ^
: string
| set offset.a = 1
< ok 0 0
> exec ^
: string
| set oj.a = 1
< ok 0 0
> exec ^
: string
| set old.a = 1
< ok 0 0
> exec ^
: string
| set oltp_read_only.a = 1
< ok 0 0
> exec ^
: string
| set oltp_read_write.a = 1
< ok 0 0
> exec ^
: string
| set oltp_write_only.a = 1
< ok 0 0
> exec ^
: string
| set on.a = 1
< ok 0 0
> exec ^
: string
| set one.a = 1
< ok 0 0
> exec ^
: string
| set only.a = 1
< ok 0 0
> exec ^
: string
| set open.a = 1
< ok 0 0
> exec ^
: string
| set optimize.a = 1
< ok 0 0
> exec ^
: string
| set optimizer_costs.a = 1
< ok 0 0
> exec ^
: string
| set option.a = 1
< ok 0 0
> exec ^
: string
| set optional.a = 1
< ok 0 0
> exec ^
: string
| set optionally.a = 1
< ok 0 0
> exec ^
: string
| set options.a = 1
< ok 0 0
> exec ^
: string
| set or.a = 1
< ok 0 0
> exec ^
: string
| set order.a = 1
< ok 0 0
> exec ^
: string
| set ordinality.a = 1
< ok 0 0
> exec ^
: string
| set organization.a = 1
< ok 0 0
> exec ^
: string
| set others.a = 1
< ok 0 0
> exec ^
: string
| set out.a = 1
< ok 0 0
> exec ^
: string
| set outer.a = 1
< ok 0 0
> exec ^
: string
| set outfile.a = 1
< ok 0 0
> exec ^
: string
| set over.a = 1
< ok 0 0
> exec ^
: string
| set overlay.a = 1
< ok 0 0
> exec ^
: string
| set owner.a = 1
< ok 0 0
> exec ^
: string
| set pack_keys.a = 1
< ok 0 0
> exec ^
: string
| set page.a = 1
< ok 0 0
> exec ^
: string
| set parallel.a = 1
< ok 0 0
> exec ^
: string
| set parameters.a = 1
< ok 0 0
> exec ^
: string
| set parse_tree.a = 1
< ok 0 0
> exec ^
: string
| set parser.a = 1
< ok 0 0
> exec ^
: string
| set partial.a = 1
< ok 0 0
> exec ^
: string
| set partition.a = 1
< ok 0 0
> exec ^
: string
| set partitioning.a = 1
< ok 0 0
> exec ^
: string
| set partitions.a = 1
< ok 0 0
> exec ^
: string
| set password.a = 1
< ok 0 0
> exec ^
: string
| set password_lock_time.a = 1
< ok 0 0
> exec ^
: string
| set path.a = 1
< ok 0 0
> exec ^
: string
| set percent_rank.a = 1
< ok 0 0
> exec ^
: string
| set persist.a = 1
< ok 0 0
> exec ^
: string
| set persist_only.a = 1
< ok 0 0
> exec ^
: string
| set phase.a = 1
< ok 0 0
> exec ^
: string
| set placement.a = 1
< ok 0 0
> exec ^
: string
| set placing.a = 1
< ok 0 0
> exec ^
: string
| set plan.a = 1
< ok 0 0
> exec ^
: string
| set plugin.a = 1
< ok 0 0
> exec ^
: string
| set plugin_dir.a = 1
< ok 0 0
> exec ^
: string
| set plugins.a = 1
< ok 0 0
> exec ^
: string
| set point.a = 1
< ok 0 0
> exec ^
: string
| set policy.a = 1
< ok 0 0
> exec ^
: string
| set polygon.a = 1
< ok 0 0
> exec ^
: string
| set port.a = 1
< ok 0 0
> exec ^
: string
| set position.a = 1
< ok 0 0
> exec ^
: string
| set pre_split_regions.a = 1
< ok 0 0
> exec ^
: string
| set precedes.a = 1
< ok 0 0
> exec ^
: string
| set preceding.a = 1
< ok 0 0
> exec ^
: string
| set precision.a = 1
< ok 0 0
> exec ^
: string
| set prepare.a = 1
< ok 0 0
> exec ^
: string
| set preserve.a = 1
< ok 0 0
> exec ^
: string
| set prev.a = 1
< ok 0 0
> exec ^
: string
| set primary.a = 1
< ok 0 0
> exec ^
: string
| set primary_region.a = 1
< ok 0 0
> exec ^
: string
| set priority.a = 1
< ok 0 0
> exec ^
: string
| set privilege_checks_user.a = 1
< ok 0 0
> exec ^
: string
| set privileges.a = 1
< ok 0 0
> exec ^
: string
| set procedure.a = 1
< ok 0 0
> exec ^
: string
| set process.a = 1
< ok 0 0
> exec ^
: string
| set processed_keys.a = 1
< ok 0 0
> exec ^
: string
| set processlist.a = 1
< ok 0 0
> exec ^
: string
| set profile.a = 1
< ok 0 0
> exec ^
: string
| set profiles.a = 1
< ok 0 0
> exec ^
: string
| set proxy.a = 1
< ok 0 0
> exec ^
: string
| set purge.a = 1
< ok 0 0
> exec ^
: string
| set qualify.a = 1
< ok 0 0
> exec ^
: string
| set quarter.a = 1
< ok 0 0
> exec ^
: string
| set query.a = 1
< ok 0 0
> exec ^
: string
| set query_limit.a = 1
< ok 0 0
> exec ^
: string
| set quick.a = 1
< ok 0 0
> exec ^
: string
| set random.a = 1
< ok 0 0
> exec ^
: string
| set range.a = 1
< ok 0 0
> exec ^
: string
| set rank.a = 1
< ok 0 0
> exec ^
: string
| set read.a = 1
< ok 0 0
> exec ^
: string
| set read_only.a = 1
< ok 0 0
> exec ^
: string
| set read_write.a = 1
< ok 0 0
> exec ^
: string
| set reads.a = 1
< ok 0 0
> exec ^
: string
| set real.a = 1
< ok 0 0
> exec ^
: string
| set rebuild.a = 1
< ok 0 0
> exec ^
: string
| set recover.a = 1
< ok 0 0
> exec ^
: string
| set recursive.a = 1
< ok 0 0
> exec ^
: string
| set redo_buffer_size.a = 1
< ok 0 0
> exec ^
: string
| set redundant.a = 1
< ok 0 0
> exec ^
: string
| set reference.a = 1
< ok 0 0
> exec ^
: string
| set references.a = 1
< ok 0 0
> exec ^
: string
| set regexp.a = 1
< ok 0 0
> exec ^
: string
| set region.a = 1
< ok 0 0
> exec ^
: string
| set regions.a = 1
< ok 0 0
> exec ^
: string
| set registration.a = 1
< ok 0 0
> exec ^
: string
| set relational.a = 1
< ok 0 0
> exec ^
: string
| set relay.a = 1
< ok 0 0
> exec ^
: string
| set relay_log_file.a = 1
< ok 0 0
> exec ^
: string
| set relay_log_pos.a = 1
< ok 0 0
> exec ^
: string
| set relay_thread.a = 1
< ok 0 0
> exec ^
: string
| set relaylog.a = 1
< ok 0 0
> exec ^
: string
| set release.a = 1
< ok 0 0
> exec ^
: string
| set reload.a = 1
< ok 0 0
> exec ^
: string
| set remove.a = 1
< ok 0 0
> exec ^
: string
| set rename.a = 1
< ok 0 0
> exec ^
: string
| set reorganize.a = 1
< ok 0 0
> exec ^
: string
| set repair.a = 1
< ok 0 0
> exec ^
: string
| set repeat.a = 1
< ok 0 0
> exec ^
: string
| set repeatable.a = 1
< ok 0 0
> exec ^
: string
| set replace.a = 1
< ok 0 0
> exec ^
: string
| set replica.a = 1
< ok 0 0
> exec ^
: string
| set replicas.a = 1
< ok 0 0
> exec ^
: string
| set replicate_do_db.a = 1
< ok 0 0
> exec ^
: string
| set replicate_do_table.a = 1
< ok 0 0
> exec ^
: string
| set replicate_ignore_db.a = 1
< ok 0 0
> exec ^
: string
| set replicate_ignore_table.a = 1
< ok 0 0
> exec ^
: string
| set replicate_rewrite_db.a = 1
< ok 0 0
> exec ^
: string
| set replicate_wild_do_table.a = 1
< ok 0 0
> exec ^
: string
| set replicate_wild_ignore_table.a = 1
< ok 0 0
> exec ^
: string
| set replication.a = 1
< ok 0 0
> exec ^
: string
| set require.a = 1
< ok 0 0
> exec ^
: string
| set require_row_format.a = 1
< ok 0 0
> exec ^
: string
| set require_table_primary_key_check.a = 1
< ok 0 0
> exec ^
: string
| set reset.a = 1
< ok 0 0
> exec ^
: string
| set resignal.a = 1
< ok 0 0
> exec ^
: string
| set resource.a = 1
< ok 0 0
> exec ^
: string
| set respect.a = 1
< ok 0 0
> exec ^
: string
| set restart.a = 1
< ok 0 0
> exec ^
: string
| set restore.a = 1
< ok 0 0
> exec ^
: string
| set restores.a = 1
< ok 0 0
> exec ^
: string
| set restrict.a = 1
< ok 0 0
> exec ^
: string
| set resume.a = 1
< ok 0 0
> exec ^
: string
| set retain.a = 1
< ok 0 0
> exec ^
: string
| set return.a = 1
< ok 0 0
> exec ^
: string
| set returned_sqlstate.a = 1
< ok 0 0
> exec ^
: string
| set returning.a = 1
< ok 0 0
> exec ^
: string
| set returns.a = 1
< ok 0 0
> exec ^
: string
| set reuse.a = 1
< ok 0 0
> exec ^
: string
| set reverse.a = 1
< ok 0 0
> exec ^
: string
| set revoke.a = 1
< ok 0 0
> exec ^
: string
| set right.a = 1
< ok 0 0
> exec ^
: string
| set rlike.a = 1
< ok 0 0
> exec ^
: string
| set role.a = 1
< ok 0 0
> exec ^
: string
| set rollback.a = 1
< ok 0 0
> exec ^
: string
| set rollup.a = 1
< ok 0 0
> exec ^
: string
| set rotate.a = 1
< ok 0 0
> exec ^
: string
| set routine.a = 1
< ok 0 0
> exec ^
: string
| set row.a = 1
< ok 0 0
> exec ^
: string
| set row_count.a = 1
< ok 0 0
> exec ^
: string
| set row_format.a = 1
< ok 0 0
> exec ^
: string
| set row_number.a = 1
< ok 0 0
> exec ^
: string
| set rows.a = 1
< ok 0 0
> exec ^
: string
| set rtree.a = 1
< ok 0 0
> exec ^
: string
| set ru.a = 1
< ok 0 0
> exec ^
: string
| set ru_per_sec.a = 1
< ok 0 0
> exec ^
: string
| set rule.a = 1
< ok 0 0
> exec ^
: string
| set run.a = 1
< ok 0 0
> exec ^
: string
| set s3.a = 1
< ok 0 0
> exec ^
: string
| set savepoint.a = 1
< ok 0 0
> exec ^
: string
| set schedule.a = 1
< ok 0 0
> exec ^
: string
| set schema.a = 1
< ok 0 0
> exec ^
: string
| set schema_name.a = 1
< ok 0 0
> exec ^
: string
| set schemas.a = 1
< ok 0 0
> exec ^
: string
| set second.a = 1
< ok 0 0
> exec ^
: string
| set second_microsecond.a = 1
< ok 0 0
> exec ^
: string
| set secondary.a = 1
< ok 0 0
> exec ^
: string
| set secondary_engine.a = 1
< ok 0 0
> exec ^
: string
| set secondary_engine_attribute.a = 1
< ok 0 0
> exec ^
: string
| set secondary_load.a = 1
< ok 0 0
> exec ^
: string
| set secondary_unload.a = 1
< ok 0 0
> exec ^
: string
| set security.a = 1
< ok 0 0
> exec ^
: string
| set select.a = 1
< ok 0 0
> exec ^
: string
| set sensitive.a = 1
< ok 0 0
> exec ^
: string
| set separator.a = 1
< ok 0 0
> exec ^
: string
| set sequence.a = 1
< ok 0 0
> exec ^
: string
| set serial.a = 1
< ok 0 0
> exec ^
: string
| set serializable.a = 1
< ok 0 0
> exec ^
: string
| set server.a = 1
< ok 0 0
> exec ^
: string
| set session.a = 1
< ok 0 0
> exec ^
: string
| set session_states.a = 1
< ok 0 0
> exec ^
: string
| set session_user.a = 1
< ok 0 0
> exec ^
: string
| set set.a = 1
< ok 0 0
> exec ^
: string
| set sets.a = 1
< ok 0 0
> exec ^
: string
| set shard_row_id_bits.a = 1
< ok 0 0
> exec ^
: string
| set share.a = 1
< ok 0 0
> exec ^
: string
| set show.a = 1
< ok 0 0
> exec ^
: string
| set shutdown.a = 1
< ok 0 0
> exec ^
: string
| set signal.a = 1
< ok 0 0
> exec ^
: string
| set signed.a = 1
< ok 0 0
> exec ^
: string
| set similar.a = 1
< ok 0 0
> exec ^
: string
| set simple.a = 1
< ok 0 0
> exec ^
: string
| set skip.a = 1
< ok 0 0
> exec ^
: string
| set slave.a = 1
< ok 0 0
> exec ^
: string
| set slow.a = 1
< ok 0 0
> exec ^
: string
| set smallint.a = 1
< ok 0 0
> exec ^
: string
| set snapshot.a = 1
< ok 0 0
> exec ^
: string
| set socket.a = 1
< ok 0 0
> exec ^
: string
| set some.a = 1
< ok 0 0
> exec ^
: string
| set soname.a = 1
< ok 0 0
> exec ^
: string
| set sounds.a = 1
< ok 0 0
> exec ^
: string
| set source.a = 1
< ok 0 0
> exec ^
: string
| set source_auto_position.a = 1
< ok 0 0
> exec ^
: string
| set source_bind.a = 1
< ok 0 0
> exec ^
: string
| set source_compression_algorithms.a = 1
< ok 0 0
> exec ^
: string
| set source_connect_retry.a = 1
< ok 0 0
> exec ^
: string
| set source_connection_auto_failover.a = 1
< ok 0 0
> exec ^
: string
| set source_delay.a = 1
< ok 0 0
> exec ^
: string
| set source_heartbeat_period.a = 1
< ok 0 0
> exec ^
: string
| set source_host.a = 1
< ok 0 0
> exec ^
: string
| set source_log_file.a = 1
< ok 0 0
> exec ^
: string
| set source_log_pos.a = 1
< ok 0 0
> exec ^
: string
| set source_password.a = 1
< ok 0 0
> exec ^
: string
| set source_port.a = 1
< ok 0 0
> exec ^
: string
| set source_public_key_path.a = 1
< ok 0 0
> exec ^
: string
| set source_retry_count.a = 1
< ok 0 0
> exec ^
: string
| set source_ssl.a = 1
< ok 0 0
> exec ^
: string
| set source_ssl_ca.a = 1
< ok 0 0
> exec ^
: string
| set source_ssl_capath.a = 1
< ok 0 0
> exec ^
: string
| set source_ssl_cert.a = 1
< ok 0 0
> exec ^
: string
| set source_ssl_cipher.a = 1
< ok 0 0
> exec ^
: string
| set source_ssl_crl.a = 1
< ok 0 0
> exec ^
: string
| set source_ssl_crlpath.a = 1
< ok 0 0
> exec ^
: string
| set source_ssl_key.a = 1
< ok 0 0
> exec ^
: string
| set source_ssl_verify_server_cert.a = 1
< ok 0 0
> exec ^
: string
| set source_tls_ciphersuites.a = 1
< ok 0 0
> exec ^
: string
| set source_tls_version.a = 1
< ok 0 0
> exec ^
: string
| set source_user.a = 1
< ok 0 0
> exec ^
: string
| set source_zstd_compression_level.a = 1
< ok 0 0
> exec ^
: string
| set spatial.a = 1
< ok 0 0
> exec ^
: string
| set specific.a = 1
< ok 0 0
> exec ^
: string
| set split.a = 1
< ok 0 0
> exec ^
: string
| set sql.a = 1
< ok 0 0
> exec ^
: string
| set sql_after_gtids.a = 1
< ok 0 0
> exec ^
: string
| set sql_after_mts_gaps.a = 1
< ok 0 0
> exec ^
: string
| set sql_before_gtids.a = 1
< ok 0 0
> exec ^
: string
| set sql_big_result.a = 1
< ok 0 0
> exec ^
: string
| set sql_buffer_result.a = 1
< ok 0 0
> exec ^
: string
| set sql_calc_found_rows.a = 1
< ok 0 0
> exec ^
: string
| set sql_no_cache.a = 1
< ok 0 0
> exec ^
: string
| set sql_small_result.a = 1
< ok 0 0
> exec ^
: string
| set sql_thread.a = 1
< ok 0 0
> exec ^
: string
| set sql_tsi_day.a = 1
< ok 0 0
> exec ^
: string
| set sql_tsi_hour.a = 1
< ok 0 0
> exec ^
: string
| set sql_tsi_minute.a = 1
< ok 0 0
> exec ^
: string
| set sql_tsi_month.a = 1
< ok 0 0
> exec ^
: string
| set sql_tsi_quarter.a = 1
< ok 0 0
> exec ^
: string
| set sql_tsi_second.a = 1
< ok 0 0
> exec ^
: string
| set sql_tsi_week.a = 1
< ok 0 0
> exec ^
: string
| set sql_tsi_year.a = 1
< ok 0 0
> exec ^
: string
| set sqlexception.a = 1
< ok 0 0
> exec ^
: string
| set sqlstate.a = 1
< ok 0 0
> exec ^
: string
| set sqlwarning.a = 1
< ok 0 0
> exec ^
: string
| set srid.a = 1
< ok 0 0
> exec ^
: string
| set ssl.a = 1
< ok 0 0
> exec ^
: string
| set st_collect.a = 1
< ok 0 0
> exec ^
: string
| set stacked.a = 1
< ok 0 0
> exec ^
: string
| set start.a = 1
< ok 0 0
> exec ^
: string
| set start_time.a = 1
< ok 0 0
> exec ^
: string
| set starting.a = 1
< ok 0 0
> exec ^
: string
| set starts.a = 1
< ok 0 0
> exec ^
: string
| set stats_auto_recalc.a = 1
< ok 0 0
> exec ^
: string
| set stats_buckets.a = 1
< ok 0 0
> exec ^
: string
| set stats_col_choice.a = 1
< ok 0 0
> exec ^
: string
| set stats_col_list.a = 1
< ok 0 0
> exec ^
: string
| set stats_extended.a = 1
< ok 0 0
> exec ^
: string
| set stats_healthy.a = 1
< ok 0 0
> exec ^
: string
| set stats_histograms.a = 1
< ok 0 0
> exec ^
: string
| set stats_locked.a = 1
< ok 0 0
> exec ^
: string
| set stats_meta.a = 1
< ok 0 0
> exec ^
: string
| set stats_persistent.a = 1
< ok 0 0
> exec ^
: string
| set stats_sample_pages.a = 1
< ok 0 0
> exec ^
: string
| set stats_sample_rate.a = 1
< ok 0 0
> exec ^
: string
| set stats_topn.a = 1
< ok 0 0
> exec ^
: string
| set status.a = 1
< ok 0 0
> exec ^
: string
| set std.a = 1
< ok 0 0
> exec ^
: string
| set stddev.a = 1
< ok 0 0
> exec ^
: string
| set stddev_pop.a = 1
< ok 0 0
> exec ^
: string
| set stddev_samp.a = 1
< ok 0 0
> exec ^
: string
| set stop.a = 1
< ok 0 0
> exec ^
: string
| set storage.a = 1
< ok 0 0
> exec ^
: string
| set stored.a = 1
< ok 0 0
> exec ^
: string
| set straight_join.a = 1
< ok 0 0
> exec ^
: string
| set stream.a = 1
< ok 0 0
> exec ^
: string
| set strict_load.a = 1
< ok 0 0
> exec ^
: string
| set string.a = 1
< ok 0 0
> exec ^
: string
| set subclass_origin.a = 1
< ok 0 0
> exec ^
: string
| set subdate.a = 1
< ok 0 0
> exec ^
: string
| set subject.a = 1
< ok 0 0
> exec ^
: string
| set subpartition.a = 1
< ok 0 0
> exec ^
: string
| set subpartitions.a = 1
< ok 0 0
> exec ^
: string
| set substr.a = 1
< ok 0 0
> exec ^
: string
| set substring.a = 1
< ok 0 0
> exec ^
: string
| set sum.a = 1
< ok 0 0
> exec ^
: string
| set super.a = 1
< ok 0 0
> exec ^
: string
| set survival_preferences.a = 1
< ok 0 0
> exec ^
: string
| set suspend.a = 1
< ok 0 0
> exec ^
: string
| set swaps.a = 1
< ok 0 0
> exec ^
: string
| set switch_group.a = 1
< ok 0 0
> exec ^
: string
| set switches.a = 1
< ok 0 0
> exec ^
: string
| set sysdate.a = 1
< ok 0 0
> exec ^
: string
| set system.a = 1
< ok 0 0
> exec ^
: string
| set system_user.a = 1
< ok 0 0
> exec ^
: string
| set table.a = 1
< ok 0 0
> exec ^
: string
| set table_checksum.a = 1
< ok 0 0
> exec ^
: string
| set table_name.a = 1
< ok 0 0
> exec ^
: string
| set tables.a = 1
< ok 0 0
> exec ^
: string
| set tablesample.a = 1
< ok 0 0
> exec ^
: string
| set tablespace.a = 1
< ok 0 0
> exec ^
: string
| set task_types.a = 1
< ok 0 0
> exec ^
: string
| set temporary.a = 1
< ok 0 0
> exec ^
: string
| set temptable.a = 1
< ok 0 0
> exec ^
: string
| set terminated.a = 1
< ok 0 0
> exec ^
: string
| set text.a = 1
< ok 0 0
> exec ^
: string
| set than.a = 1
< ok 0 0
> exec ^
: string
| set then.a = 1
< ok 0 0
> exec ^
: string
| set thread_priority.a = 1
< ok 0 0
> exec ^
: string
| set ties.a = 1
< ok 0 0
> exec ^
: string
| set tiflash.a = 1
< ok 0 0
> exec ^
: string
| set time.a = 1
< ok 0 0
> exec ^
: string
| set timeout.a = 1
< ok 0 0
> exec ^
: string
| set timestamp.a = 1
< ok 0 0
> exec ^
: string
| set timestampadd.a = 1
< ok 0 0
> exec ^
: string
| set timestampdiff.a = 1
< ok 0 0
> exec ^
: string
| set tinyblob.a = 1
< ok 0 0
> exec ^
: string
| set tinyint.a = 1
< ok 0 0
> exec ^
: string
| set tinytext.a = 1
< ok 0 0
> exec ^
: string
| set tls.a = 1
< ok 0 0
> exec ^
: string
| set to.a = 1
< ok 0 0
> exec ^
: string
| set tpcc.a = 1
< ok 0 0
> exec ^
: string
| set tpch_10.a = 1
< ok 0 0
> exec ^
: string
| set trailing.a = 1
< ok 0 0
> exec ^
: string
| set transaction.a = 1
< ok 0 0
> exec ^
: string
| set trigger.a = 1
< ok 0 0
> exec ^
: string
| set triggers.a = 1
< ok 0 0
> exec ^
: string
| set trim.a = 1
< ok 0 0
> exec ^
: string
| set true.a = 1
< ok 0 0
> exec ^
: string
| set truncate.a = 1
< ok 0 0
> exec ^
: string
| set ttl.a = 1
< ok 0 0
> exec ^
: string
| set ttl_enable.a = 1
< ok 0 0
> exec ^
: string
| set ttl_job_interval.a = 1
< ok 0 0
> exec ^
: string
| set type.a = 1
< ok 0 0
> exec ^
: string
| set types.a = 1
< ok 0 0
> exec ^
: string
| set unbounded.a = 1
< ok 0 0
> exec ^
: string
| set uncommitted.a = 1
< ok 0 0
> exec ^
: string
| set undefined.a = 1
< ok 0 0
> exec ^
: string
| set undo.a = 1
< ok 0 0
> exec ^
: string
| set undo_buffer_size.a = 1
< ok 0 0
> exec ^
: string
| set undofile.a = 1
< ok 0 0
> exec ^
: string
| set unicode.a = 1
< ok 0 0
> exec ^
: string
| set uninstall.a = 1
< ok 0 0
> exec ^
: string
| set union.a = 1
< ok 0 0
> exec ^
: string
| set unique.a = 1
< ok 0 0
> exec ^
: string
| set unknown.a = 1
< ok 0 0
> exec ^
: string
| set unlimited.a = 1
< ok 0 0
> exec ^
: string
| set unlock.a = 1
< ok 0 0
> exec ^
: string
| set unregister.a = 1
< ok 0 0
> exec ^
: string
| set unsigned.a = 1
< ok 0 0
> exec ^
: string
| set until.a = 1
< ok 0 0
> exec ^
: string
| set update.a = 1
< ok 0 0
> exec ^
: string
| set upgrade.a = 1
< ok 0 0
> exec ^
: string
| set uri.a = 1
< ok 0 0
> exec ^
: string
| set url.a = 1
< ok 0 0
> exec ^
: string
| set usage.a = 1
< ok 0 0
> exec ^
: string
| set use.a = 1
< ok 0 0
> exec ^
: string
| set use_frm.a = 1
< ok 0 0
> exec ^
: string
| set user.a = 1
< ok 0 0
> exec ^
: string
| set user_resources.a = 1
< ok 0 0
> exec ^
: string
| set using.a = 1
< ok 0 0
> exec ^
: string
| set utc_date.a = 1
< ok 0 0
> exec ^
: string
| set utc_time.a = 1
< ok 0 0
> exec ^
: string
| set utc_timestamp.a = 1
< ok 0 0
> exec ^
: string
| set utilization_limit.a = 1
< ok 0 0
> exec ^
: string
| set validate.a = 1
< ok 0 0
> exec ^
: string
| set validation.a = 1
< ok 0 0
> exec ^
: string
| set value.a = 1
< ok 0 0
> exec ^
: string
| set values.a = 1
< ok 0 0
> exec ^
: string
| set var_pop.a = 1
< ok 0 0
> exec ^
: string
| set var_samp.a = 1
< ok 0 0
> exec ^
: string
| set varbinary.a = 1
< ok 0 0
> exec ^
: string
| set varchar.a = 1
< ok 0 0
> exec ^
: string
| set varcharacter.a = 1
< ok 0 0
> exec ^
: string
| set variables.a = 1
< ok 0 0
> exec ^
: string
| set variance.a = 1
< ok 0 0
> exec ^
: string
| set varying.a = 1
< ok 0 0
> exec ^
: string
| set vcpu.a = 1
< ok 0 0
> exec ^
: string
| set vector.a = 1
< ok 0 0
> exec ^
: string
| set verify_key_constraints.a = 1
< ok 0 0
> exec ^
: string
| set view.a = 1
< ok 0 0
> exec ^
: string
| set virtual.a = 1
< ok 0 0
> exec ^
: string
| set visible.a = 1
< ok 0 0
> exec ^
: string
| set voter_constraints.a = 1
< ok 0 0
> exec ^
: string
| set voters.a = 1
< ok 0 0
> exec ^
: string
| set wait.a = 1
< ok 0 0
> exec ^
: string
| set warnings.a = 1
< ok 0 0
> exec ^
: string
| set watch.a = 1
< ok 0 0
> exec ^
: string
| set week.a = 1
< ok 0 0
> exec ^
: string
| set weight_string.a = 1
< ok 0 0
> exec ^
: string
| set when.a = 1
< ok 0 0
> exec ^
: string
| set where.a = 1
< ok 0 0
> exec ^
: string
| set while.a = 1
< ok 0 0
> exec ^
: string
| set window.a = 1
< ok 0 0
> exec ^
: string
| set with.a = 1
< ok 0 0
> exec ^
: string
| set without.a = 1
< ok 0 0
> exec ^
: string
| set work.a = 1
< ok 0 0
> exec ^
: string
| set workload.a = 1
< ok 0 0
> exec ^
: string
| set wrapper.a = 1
< ok 0 0
> exec ^
: string
| set write.a = 1
< ok 0 0
> exec ^
: string
| set x509.a = 1
< ok 0 0
> exec ^
: string
| set xa.a = 1
< ok 0 0
> exec ^
: string
| set xid.a = 1
< ok 0 0
> exec ^
: string
| set xml.a = 1
< ok 0 0
> exec ^
: string
| set xor.a = 1
< ok 0 0
> exec ^
: string
| set year.a = 1
< ok 0 0
> exec ^
: string
| set year_month.a = 1
< ok 0 0
> exec ^
: string
| set zerofill.a = 1
< ok 0 0
> exec ^
: string
| set zone.a = 1
< ok 0 0
> exec ^
: string
| set a.absent = 1
< ok 0 0
> exec ^
: string
| set a.accessible = 1
< ok 0 0
> exec ^
: string
| set a.account = 1
< ok 0 0
> exec ^
: string
| set a.action = 1
< ok 0 0
> exec ^
: string
| set a.active = 1
< ok 0 0
> exec ^
: string
| set a.add = 1
< ok 0 0
> exec ^
: string
| set a.adddate = 1
< ok 0 0
> exec ^
: string
| set a.admin = 1
< ok 0 0
> exec ^
: string
| set a.after = 1
< ok 0 0
> exec ^
: string
| set a.against = 1
< ok 0 0
> exec ^
: string
| set a.aggregate = 1
< ok 0 0
> exec ^
: string
| set a.algorithm = 1
< ok 0 0
> exec ^
: string
| set a.all = 1
< ok 0 0
> exec ^
: string
| set a.allow_missing_files = 1
< ok 0 0
> exec ^
: string
| set a.alter = 1
< ok 0 0
> exec ^
: string
| set a.always = 1
< ok 0 0
> exec ^
: string
| set a.analyze = 1
< ok 0 0
> exec ^
: string
| set a.and = 1
< ok 0 0
> exec ^
: string
| set a.any = 1
< ok 0 0
> exec ^
: string
| set a.array = 1
< ok 0 0
> exec ^
: string
| set a.as = 1
< ok 0 0
> exec ^
: string
| set a.asc = 1
< ok 0 0
> exec ^
: string
| set a.ascii = 1
< ok 0 0
> exec ^
: string
| set a.asensitive = 1
< ok 0 0
> exec ^
: string
| set a.assign_gtids_to_anonymous_transactions = 1
< ok 0 0
> exec ^
: string
| set a.at = 1
< ok 0 0
> exec ^
: string
| set a.attribute = 1
< ok 0 0
> exec ^
: string
| set a.authentication = 1
< ok 0 0
> exec ^
: string
| set a.auto = 1
< ok 0 0
> exec ^
: string
| set a.auto_id_cache = 1
< ok 0 0
> exec ^
: string
| set a.auto_increment = 1
< ok 0 0
> exec ^
: string
| set a.auto_random = 1
< ok 0 0
> exec ^
: string
| set a.auto_random_base = 1
< ok 0 0
> exec ^
: string
| set a.auto_refresh = 1
< ok 0 0
> exec ^
: string
| set a.auto_refresh_source = 1
< ok 0 0
> exec ^
: string
| set a.autoextend_size = 1
< ok 0 0
> exec ^
: string
| set a.avg = 1
< ok 0 0
> exec ^
: string
| set a.avg_row_length = 1
< ok 0 0
> exec ^
: string
| set a.background = 1
< ok 0 0
> exec ^
: string
| set a.backup = 1
< ok 0 0
> exec ^
: string
| set a.backups = 1
< ok 0 0
> exec ^
: string
| set a.batch = 1
< ok 0 0
> exec ^
: string
| set a.before = 1
< ok 0 0
> exec ^
: string
| set a.begin = 1
< ok 0 0
> exec ^
: string
| set a.bernoulli = 1
< ok 0 0
> exec ^
: string
| set a.between = 1
< ok 0 0
> exec ^
: string
| set a.bigint = 1
< ok 0 0
> exec ^
: string
| set a.binary = 1
< ok 0 0
> exec ^
: string
| set a.binding = 1
< ok 0 0
> exec ^
: string
| set a.binding_cache = 1
< ok 0 0
> exec ^
: string
| set a.bindings = 1
< ok 0 0
> exec ^
: string
| set a.binlog = 1
< ok 0 0
> exec ^
: string
| set a.bit = 1
< ok 0 0
> exec ^
: string
| set a.bit_and = 1
< ok 0 0
> exec ^
: string
| set a.bit_or = 1
< ok 0 0
> exec ^
: string
| set a.bit_xor = 1
< ok 0 0
> exec ^
: string
| set a.blob = 1
< ok 0 0
> exec ^
: string
| set a.block = 1
< ok 0 0
> exec ^
: string
| set a.bool = 1
< ok 0 0
> exec ^
: string
| set a.boolean = 1
< ok 0 0
> exec ^
: string
| set a.both = 1
< ok 0 0
> exec ^
: string
| set a.btree = 1
< ok 0 0
> exec ^
: string
| set a.buckets = 1
< ok 0 0
> exec ^
: string
| set a.builtins = 1
< ok 0 0
> exec ^
: string
| set a.bulk = 1
< ok 0 0
> exec ^
: string
| set a.burstable = 1
< ok 0 0
> exec ^
: string
| set a.by = 1
< ok 0 0
> exec ^
: string
| set a.byte = 1
< ok 0 0
> exec ^
: string
| set a.cache = 1
< ok 0 0
> exec ^
: string
| set a.calibrate = 1
< ok 0 0
> exec ^
: string
| set a.call = 1
< ok 0 0
> exec ^
: string
| set a.cascade = 1
< ok 0 0
> exec ^
: string
| set a.cascaded = 1
< ok 0 0
> exec ^
: string
| set a.case = 1
< ok 0 0
> exec ^
: string
| set a.cast = 1
< ok 0 0
> exec ^
: string
| set a.catalog_name = 1
< ok 0 0
> exec ^
: string
| set a.chain = 1
< ok 0 0
> exec ^
: string
| set a.challenge_response = 1
< ok 0 0
> exec ^
: string
| set a.change = 1
< ok 0 0
> exec ^
: string
| set a.changed = 1
< ok 0 0
> exec ^
: string
| set a.channel = 1
< ok 0 0
> exec ^
: string
| set a.char = 1
< ok 0 0
> exec ^
: string
| set a.character = 1
< ok 0 0
> exec ^
: string
| set a.charset = 1
< ok 0 0
> exec ^
: string
| set a.check = 1
< ok 0 0
> exec ^
: string
| set a.checksum = 1
< ok 0 0
> exec ^
: string
| set a.cipher = 1
< ok 0 0
> exec ^
: string
| set a.class_origin = 1
< ok 0 0
> exec ^
: string
| set a.client = 1
< ok 0 0
> exec ^
: string
| set a.clone = 1
< ok 0 0
> exec ^
: string
| set a.close = 1
< ok 0 0
> exec ^
: string
| set a.clustered = 1
< ok 0 0
> exec ^
: string
| set a.coalesce = 1
< ok 0 0
> exec ^
: string
| set a.code = 1
< ok 0 0
> exec ^
: string
| set a.collate = 1
< ok 0 0
> exec ^
: string
| set a.collation = 1
< ok 0 0
> exec ^
: string
| set a.column = 1
< ok 0 0
> exec ^
: string
| set a.column_format = 1
< ok 0 0
> exec ^
: string
| set a.column_name = 1
< ok 0 0
> exec ^
: string
| set a.column_stats_usage = 1
< ok 0 0
> exec ^
: string
| set a.columns = 1
< ok 0 0
> exec ^
: string
| set a.comment = 1
< ok 0 0
> exec ^
: string
| set a.commit = 1
< ok 0 0
> exec ^
: string
| set a.committed = 1
< ok 0 0
> exec ^
: string
| set a.compact = 1
< ok 0 0
> exec ^
: string
| set a.completion = 1
< ok 0 0
> exec ^
: string
| set a.component = 1
< ok 0 0
> exec ^
: string
| set a.compressed = 1
< ok 0 0
> exec ^
: string
| set a.compression = 1
< ok 0 0
> exec ^
: string
| set a.concat = 1
< ok 0 0
> exec ^
: string
| set a.concurrent = 1
< ok 0 0
> exec ^
: string
| set a.condition = 1
< ok 0 0
> exec ^
: string
| set a.config = 1
< ok 0 0
> exec ^
: string
| set a.connection = 1
< ok 0 0
> exec ^
: string
| set a.consistent = 1
< ok 0 0
> exec ^
: string
| set a.constraint = 1
< ok 0 0
> exec ^
: string
| set a.constraint_catalog = 1
< ok 0 0
> exec ^
: string
| set a.constraint_name = 1
< ok 0 0
> exec ^
: string
| set a.constraint_schema = 1
< ok 0 0
> exec ^
: string
| set a.constraints = 1
< ok 0 0
> exec ^
: string
| set a.contains = 1
< ok 0 0
> exec ^
: string
| set a.context = 1
< ok 0 0
> exec ^
: string
| set a.continue = 1
< ok 0 0
> exec ^
: string
| set a.convert = 1
< ok 0 0
> exec ^
: string
| set a.cooldown = 1
< ok 0 0
> exec ^
: string
| set a.count = 1
< ok 0 0
> exec ^
: string
| set a.cpu = 1
< ok 0 0
> exec ^
: string
| set a.create = 1
< ok 0 0
> exec ^
: string
| set a.cross = 1
< ok 0 0
> exec ^
: string
| set a.cube = 1
< ok 0 0
> exec ^
: string
| set a.cume_dist = 1
< ok 0 0
> exec ^
: string
| set a.curdate = 1
< ok 0 0
> exec ^
: string
| set a.current = 1
< ok 0 0
> exec ^
: string
| set a.current_date = 1
< ok 0 0
> exec ^
: string
| set a.current_time = 1
< ok 0 0
> exec ^
: string
| set a.current_timestamp = 1
< ok 0 0
> exec ^
: string
| set a.current_user = 1
< ok 0 0
> exec ^
: string
| set a.cursor = 1
< ok 0 0
> exec ^
: string
| set a.cursor_name = 1
< ok 0 0
> exec ^
: string
| set a.curtime = 1
< ok 0 0
> exec ^
: string
| set a.cycle = 1
< ok 0 0
> exec ^
: string
| set a.data = 1
< ok 0 0
> exec ^
: string
| set a.database = 1
< ok 0 0
> exec ^
: string
| set a.databases = 1
< ok 0 0
> exec ^
: string
| set a.datafile = 1
< ok 0 0
> exec ^
: string
| set a.date = 1
< ok 0 0
> exec ^
: string
| set a.date_add = 1
< ok 0 0
> exec ^
: string
| set a.date_sub = 1
< ok 0 0
> exec ^
: string
| set a.datetime = 1
< ok 0 0
> exec ^
: string
| set a.day = 1
< ok 0 0
> exec ^
: string
| set a.day_hour = 1
< ok 0 0
> exec ^
: string
| set a.day_microsecond = 1
< ok 0 0
> exec ^
: string
| set a.day_minute = 1
< ok 0 0
> exec ^
: string
| set a.day_second = 1
< ok 0 0
> exec ^
: string
| set a.deallocate = 1
< ok 0 0
> exec ^
: string
| set a.dec = 1
< ok 0 0
> exec ^
: string
| set a.decimal = 1
< ok 0 0
> exec ^
: string
| set a.declare = 1
< ok 0 0
> exec ^
: string
| set a.default = 1
< ok 0 0
> exec ^
: string
| set a.default_auth = 1
< ok 0 0
> exec ^
: string
| set a.defined = 1
< ok 0 0
> exec ^
: string
| set a.definer = 1
< ok 0 0
> exec ^
: string
| set a.definition = 1
< ok 0 0
> exec ^
: string
| set a.delay_key_write = 1
< ok 0 0
> exec ^
: string
| set a.delayed = 1
< ok 0 0
> exec ^
: string
| set a.delete = 1
< ok 0 0
> exec ^
: string
| set a.dense_rank = 1
< ok 0 0
> exec ^
: string
| set a.desc = 1
< ok 0 0
> exec ^
: string
| set a.describe = 1
< ok 0 0
> exec ^
: string
| set a.description = 1
< ok 0 0
> exec ^
: string
| set a.deterministic = 1
< ok 0 0
> exec ^
: string
| set a.diagnostics = 1
< ok 0 0
> exec ^
: string
| set a.digest = 1
< ok 0 0
> exec ^
: string
| set a.directory = 1
< ok 0 0
> exec ^
: string
| set a.disable = 1
< ok 0 0
> exec ^
: string
| set a.disabled = 1
< ok 0 0
> exec ^
: string
| set a.discard = 1
< ok 0 0
> exec ^
: string
| set a.disk = 1
< ok 0 0
> exec ^
: string
| set a.distinct = 1
< ok 0 0
> exec ^
: string
| set a.distinctrow = 1
< ok 0 0
> exec ^
: string
| set a.distribute = 1
< ok 0 0
> exec ^
: string
| set a.distribution = 1
< ok 0 0
> exec ^
: string
| set a.distributions = 1
< ok 0 0
> exec ^
: string
| set a.div = 1
< ok 0 0
> exec ^
: string
| set a.do = 1
< ok 0 0
> exec ^
: string
| set a.double = 1
< ok 0 0
> exec ^
: string
| set a.drop = 1
< ok 0 0
> exec ^
: string
| set a.dry = 1
< ok 0 0
> exec ^
: string
| set a.dryrun = 1
< ok 0 0
> exec ^
: string
| set a.dual = 1
< ok 0 0
> exec ^
: string
| set a.duality = 1
< ok 0 0
> exec ^
: string
| set a.dumpfile = 1
< ok 0 0
> exec ^
: string
| set a.duplicate = 1
< ok 0 0
> exec ^
: string
| set a.duration = 1
< ok 0 0
> exec ^
: string
| set a.dynamic = 1
< ok 0 0
> exec ^
: string
| set a.each = 1
< ok 0 0
> exec ^
: string
| set a.else = 1
< ok 0 0
> exec ^
: string
| set a.elseif = 1
< ok 0 0
> exec ^
: string
| set a.empty = 1
< ok 0 0
> exec ^
: string
| set a.enable = 1
< ok 0 0
> exec ^
: string
| set a.enabled = 1
< ok 0 0
> exec ^
: string
| set a.enclosed = 1
< ok 0 0
> exec ^
: string
| set a.encryption = 1
< ok 0 0
> exec ^
: string
| set a.end = 1
< ok 0 0
> exec ^
: string
| set a.end_time = 1
< ok 0 0
> exec ^
: string
| set a.ends = 1
< ok 0 0
> exec ^
: string
| set a.enforced = 1
< ok 0 0
> exec ^
: string
| set a.engine = 1
< ok 0 0
> exec ^
: string
| set a.engine_attribute = 1
< ok 0 0
> exec ^
: string
| set a.engines = 1
< ok 0 0
> exec ^
: string
| set a.enum = 1
< ok 0 0
> exec ^
: string
| set a.error = 1
< ok 0 0
> exec ^
: string
| set a.errors = 1
< ok 0 0
> exec ^
: string
| set a.escape = 1
< ok 0 0
> exec ^
: string
| set a.escaped = 1
< ok 0 0
> exec ^
: string
| set a.event = 1
< ok 0 0
> exec ^
: string
| set a.events = 1
< ok 0 0
> exec ^
: string
| set a.every = 1
< ok 0 0
> exec ^
: string
| set a.exact = 1
< ok 0 0
> exec ^
: string
| set a.except = 1
< ok 0 0
> exec ^
: string
| set a.exchange = 1
< ok 0 0
> exec ^
: string
| set a.exclude = 1
< ok 0 0
> exec ^
: string
| set a.exec_elapsed = 1
< ok 0 0
> exec ^
: string
| set a.execute = 1
< ok 0 0
> exec ^
: string
| set a.exists = 1
< ok 0 0
> exec ^
: string
| set a.exit = 1
< ok 0 0
> exec ^
: string
| set a.expansion = 1
< ok 0 0
> exec ^
: string
| set a.expire = 1
< ok 0 0
> exec ^
: string
| set a.explain = 1
< ok 0 0
> exec ^
: string
| set a.export = 1
< ok 0 0
> exec ^
: string
| set a.extended = 1
< ok 0 0
> exec ^
: string
| set a.extent_size = 1
< ok 0 0
> exec ^
: string
| set a.external = 1
< ok 0 0
> exec ^
: string
| set a.external_format = 1
< ok 0 0
> exec ^
: string
| set a.extract = 1
< ok 0 0
> exec ^
: string
| set a.factor = 1
< ok 0 0
> exec ^
: string
| set a.failed_login_attempts = 1
< ok 0 0
> exec ^
: string
| set a.false = 1
< ok 0 0
> exec ^
: string
| set a.fast = 1
< ok 0 0
> exec ^
: string
| set a.faults = 1
< ok 0 0
> exec ^
: string
| set a.fetch = 1
< ok 0 0
> exec ^
: string
| set a.fields = 1
< ok 0 0
> exec ^
: string
| set a.file = 1
< ok 0 0
> exec ^
: string
| set a.file_block_size = 1
< ok 0 0
> exec ^
: string
| set a.file_format = 1
< ok 0 0
> exec ^
: string
| set a.file_name = 1
< ok 0 0
> exec ^
: string
| set a.file_pattern = 1
< ok 0 0
> exec ^
: string
| set a.file_prefix = 1
< ok 0 0
> exec ^
: string
| set a.files = 1
< ok 0 0
> exec ^
: string
| set a.filter = 1
< ok 0 0
> exec ^
: string
| set a.finish = 1
< ok 0 0
> exec ^
: string
| set a.first = 1
< ok 0 0
> exec ^
: string
| set a.first_value = 1
< ok 0 0
> exec ^
: string
| set a.fixed = 1
< ok 0 0
> exec ^
: string
| set a.flashback = 1
< ok 0 0
> exec ^
: string
| set a.float = 1
< ok 0 0
> exec ^
: string
| set a.float4 = 1
< ok 0 0
> exec ^
: string
| set a.float8 = 1
< ok 0 0
> exec ^
: string
| set a.flush = 1
< ok 0 0
> exec ^
: string
| set a.follower_constraints = 1
< ok 0 0
> exec ^
: string
| set a.followers = 1
< ok 0 0
> exec ^
: string
| set a.following = 1
< ok 0 0
> exec ^
: string
| set a.follows = 1
< ok 0 0
> exec ^
: string
| set a.for = 1
< ok 0 0
> exec ^
: string
| set a.force = 1
< ok 0 0
> exec ^
: string
| set a.foreign = 1
< ok 0 0
> exec ^
: string
| set a.format = 1
< ok 0 0
> exec ^
: string
| set a.found = 1
< ok 0 0
> exec ^
: string
| set a.from = 1
< ok 0 0
> exec ^
: string
| set a.full = 1
< ok 0 0
> exec ^
: string
| set a.fulltext = 1
< ok 0 0
> exec ^
: string
| set a.function = 1
< ok 0 0
> exec ^
: string
| set a.general = 1
< ok 0 0
> exec ^
: string
| set a.generate = 1
< ok 0 0
> exec ^
: string
| set a.generated = 1
< ok 0 0
> exec ^
: string
| set a.geomcollection = 1
< ok 0 0
> exec ^
: string
| set a.geometry = 1
< ok 0 0
> exec ^
: string
| set a.geometrycollection = 1
< ok 0 0
> exec ^
: string
| set a.get = 1
< ok 0 0
> exec ^
: string
| set a.get_format = 1
< ok 0 0
> exec ^
: string
| set a.get_source_public_key = 1
< ok 0 0
> exec ^
: string
| set a.global = 1
< ok 0 0
> exec ^
: string
| set a.grant = 1
< ok 0 0
> exec ^
: string
| set a.grants = 1
< ok 0 0
> exec ^
: string
| set a.greatest = 1
< ok 0 0
> exec ^
: string
| set a.group = 1
< ok 0 0
> exec ^
: string
| set a.group_concat = 1
< ok 0 0
> exec ^
: string
| set a.group_replication = 1
< ok 0 0
> exec ^
: string
| set a.grouping = 1
< ok 0 0
> exec ^
: string
| set a.groups = 1
< ok 0 0
> exec ^
: string
| set a.gtid_only = 1
< ok 0 0
> exec ^
: string
| set a.gtids = 1
< ok 0 0
> exec ^
: string
| set a.guided = 1
< ok 0 0
> exec ^
: string
| set a.handler = 1
< ok 0 0
> exec ^
: string
| set a.hash = 1
< ok 0 0
> exec ^
: string
| set a.having = 1
< ok 0 0
> exec ^
: string
| set a.header = 1
< ok 0 0
> exec ^
: string
| set a.help = 1
< ok 0 0
> exec ^
: string
| set a.high = 1
< ok 0 0
> exec ^
: string
| set a.high_priority = 1
< ok 0 0
> exec ^
: string
| set a.histogram = 1
< ok 0 0
> exec ^
: string
| set a.histograms_in_flight = 1
< ok 0 0
> exec ^
: string
| set a.history = 1
< ok 0 0
> exec ^
: string
| set a.host = 1
< ok 0 0
> exec ^
: string
| set a.hosts = 1
< ok 0 0
> exec ^
: string
| set a.hour = 1
< ok 0 0
> exec ^
: string
| set a.hour_microsecond = 1
< ok 0 0
> exec ^
: string
| set a.hour_minute = 1
< ok 0 0
> exec ^
: string
| set a.hour_second = 1
< ok 0 0
> exec ^
: string
| set a.identified = 1
< ok 0 0
> exec ^
: string
| set a.if = 1
< ok 0 0
> exec ^
: string
| set a.ignore = 1
< ok 0 0
> exec ^
: string
| set a.ignore_server_ids = 1
< ok 0 0
> exec ^
: string
| set a.ilike = 1
< ok 0 0
> exec ^
: string
| set a.import = 1
< ok 0 0
> exec ^
: string
| set a.in = 1
< ok 0 0
> exec ^
: string
| set a.inactive = 1
< ok 0 0
> exec ^
: string
| set a.increment = 1
< ok 0 0
> exec ^
: string
| set a.index = 1
< ok 0 0
> exec ^
: string
| set a.indexes = 1
< ok 0 0
> exec ^
: string
| set a.infile = 1
< ok 0 0
> exec ^
: string
| set a.initial = 1
< ok 0 0
> exec ^
: string
| set a.initial_size = 1
< ok 0 0
> exec ^
: string
| set a.initiate = 1
< ok 0 0
> exec ^
: string
| set a.inner = 1
< ok 0 0
> exec ^
: string
| set a.inout = 1
< ok 0 0
> exec ^
: string
| set a.insensitive = 1
< ok 0 0
> exec ^
: string
| set a.insert = 1
< ok 0 0
> exec ^
: string
| set a.insert_method = 1
< ok 0 0
> exec ^
: string
| set a.install = 1
< ok 0 0
> exec ^
: string
| set a.instance = 1
< ok 0 0
> exec ^
: string
| set a.int = 1
< ok 0 0
> exec ^
: string
| set a.int1 = 1
< ok 0 0
> exec ^
: string
| set a.int2 = 1
< ok 0 0
> exec ^
: string
| set a.int3 = 1
< ok 0 0
> exec ^
: string
| set a.int4 = 1
< ok 0 0
> exec ^
: string
| set a.int8 = 1
< ok 0 0
> exec ^
: string
| set a.integer = 1
< ok 0 0
> exec ^
: string
| set a.intersect = 1
< ok 0 0
> exec ^
: string
| set a.interval = 1
< ok 0 0
> exec ^
: string
| set a.into = 1
< ok 0 0
> exec ^
: string
| set a.invisible = 1
< ok 0 0
> exec ^
: string
| set a.invoker = 1
< ok 0 0
> exec ^
: string
| set a.io = 1
< ok 0 0
> exec ^
: string
| set a.io_after_gtids = 1
< ok 0 0
> exec ^
: string
| set a.io_before_gtids = 1
< ok 0 0
> exec ^
: string
| set a.io_thread = 1
< ok 0 0
> exec ^
: string
| set a.ipc = 1
< ok 0 0
> exec ^
: string
| set a.is = 1
< ok 0 0
> exec ^
: string
| set a.isolation = 1
< ok 0 0
> exec ^
: string
| set a.issuer = 1
< ok 0 0
> exec ^
: string
| set a.iterate = 1
< ok 0 0
> exec ^
: string
| set a.job = 1
< ok 0 0
> exec ^
: string
| set a.jobs = 1
< ok 0 0
> exec ^
: string
| set a.join = 1
< ok 0 0
> exec ^
: string
| set a.json = 1
< ok 0 0
> exec ^
: string
| set a.json_arrayagg = 1
< ok 0 0
> exec ^
: string
| set a.json_duality_object = 1
< ok 0 0
> exec ^
: string
| set a.json_objectagg = 1
< ok 0 0
> exec ^
: string
| set a.json_table = 1
< ok 0 0
> exec ^
: string
| set a.json_value = 1
< ok 0 0
> exec ^
: string
| set a.key = 1
< ok 0 0
> exec ^
: string
| set a.key_block_size = 1
< ok 0 0
> exec ^
: string
| set a.keyring = 1
< ok 0 0
> exec ^
: string
| set a.keys = 1
< ok 0 0
> exec ^
: string
| set a.kill = 1
< ok 0 0
> exec ^
: string
| set a.labels = 1
< ok 0 0
> exec ^
: string
| set a.lag = 1
< ok 0 0
> exec ^
: string
| set a.language = 1
< ok 0 0
> exec ^
: string
| set a.last = 1
< ok 0 0
> exec ^
: string
| set a.last_value = 1
< ok 0 0
> exec ^
: string
| set a.lateral = 1
< ok 0 0
> exec ^
: string
| set a.lead = 1
< ok 0 0
> exec ^
: string
| set a.leader_constraints = 1
< ok 0 0
> exec ^
: string
| set a.leading = 1
< ok 0 0
> exec ^
: string
| set a.learner_constraints = 1
< ok 0 0
> exec ^
: string
| set a.learners = 1
< ok 0 0
> exec ^
: string
| set a.least = 1
< ok 0 0
> exec ^
: string
| set a.leave = 1
< ok 0 0
> exec ^
: string
| set a.leaves = 1
< ok 0 0
> exec ^
: string
| set a.left = 1
< ok 0 0
> exec ^
: string
| set a.less = 1
< ok 0 0
> exec ^
: string
| set a.level = 1
< ok 0 0
> exec ^
: string
| set a.library = 1
< ok 0 0
> exec ^
: string
| set a.like = 1
< ok 0 0
> exec ^
: string
| set a.limit = 1
< ok 0 0
> exec ^
: string
| set a.linear = 1
< ok 0 0
> exec ^
: string
| set a.lines = 1
< ok 0 0
> exec ^
: string
| set a.linestring = 1
< ok 0 0
> exec ^
: string
| set a.list = 1
< ok 0 0
> exec ^
: string
| set a.load = 1
< ok 0 0
> exec ^
: string
| set a.local = 1
< ok 0 0
> exec ^
: string
| set a.localtime = 1
< ok 0 0
> exec ^
: string
| set a.localtimestamp = 1
< ok 0 0
> exec ^
: string
| set a.location = 1
< ok 0 0
> exec ^
: string
| set a.lock = 1
< ok 0 0
> exec ^
: string
| set a.locked = 1
< ok 0 0
> exec ^
: string
| set a.locks = 1
< ok 0 0
> exec ^
: string
| set a.log = 1
< ok 0 0
> exec ^
: string
| set a.logfile = 1
< ok 0 0
> exec ^
: string
| set a.logs = 1
< ok 0 0
> exec ^
: string
| set a.long = 1
< ok 0 0
> exec ^
: string
| set a.longblob = 1
< ok 0 0
> exec ^
: string
| set a.longtext = 1
< ok 0 0
> exec ^
: string
| set a.loop = 1
< ok 0 0
> exec ^
: string
| set a.low = 1
< ok 0 0
> exec ^
: string
| set a.low_priority = 1
< ok 0 0
> exec ^
: string
| set a.manual = 1
< ok 0 0
> exec ^
: string
| set a.master = 1
< ok 0 0
> exec ^
: string
| set a.master_log_file = 1
< ok 0 0
> exec ^
: string
| set a.match = 1
< ok 0 0
> exec ^
: string
| set a.materialized = 1
< ok 0 0
> exec ^
: string
| set a.max = 1
< ok 0 0
> exec ^
: string
| set a.max_connections_per_hour = 1
< ok 0 0
> exec ^
: string
| set a.max_queries_per_hour = 1
< ok 0 0
> exec ^
: string
| set a.max_rows = 1
< ok 0 0
> exec ^
: string
| set a.max_size = 1
< ok 0 0
> exec ^
: string
| set a.max_updates_per_hour = 1
< ok 0 0
> exec ^
: string
| set a.max_user_connections = 1
< ok 0 0
> exec ^
: string
| set a.maxvalue = 1
< ok 0 0
> exec ^
: string
| set a.medium = 1
< ok 0 0
> exec ^
: string
| set a.mediumblob = 1
< ok 0 0
> exec ^
: string
| set a.mediumint = 1
< ok 0 0
> exec ^
: string
| set a.mediumtext = 1
< ok 0 0
> exec ^
: string
| set a.member = 1
< ok 0 0
> exec ^
: string
| set a.memory = 1
< ok 0 0
> exec ^
: string
| set a.merge = 1
< ok 0 0
> exec ^
: string
| set a.message_text = 1
< ok 0 0
> exec ^
: string
| set a.microsecond = 1
< ok 0 0
> exec ^
: string
| set a.mid = 1
< ok 0 0
> exec ^
: string
| set a.middleint = 1
< ok 0 0
> exec ^
: string
| set a.migrate = 1
< ok 0 0
> exec ^
: string
| set a.min = 1
< ok 0 0
> exec ^
: string
| set a.min_rows = 1
< ok 0 0
> exec ^
: string
| set a.minute = 1
< ok 0 0
> exec ^
: string
| set a.minute_microsecond = 1
< ok 0 0
> exec ^
: string
| set a.minute_second = 1
< ok 0 0
> exec ^
: string
| set a.minvalue = 1
< ok 0 0
> exec ^
: string
| set a.mod = 1
< ok 0 0
> exec ^
: string
| set a.mode = 1
< ok 0 0
> exec ^
: string
| set a.modifies = 1
< ok 0 0
> exec ^
: string
| set a.modify = 1
< ok 0 0
> exec ^
: string
| set a.month = 1
< ok 0 0
> exec ^
: string
| set a.multilinestring = 1
< ok 0 0
> exec ^
: string
| set a.multipoint = 1
< ok 0 0
> exec ^
: string
| set a.multipolygon = 1
< ok 0 0
> exec ^
: string
| set a.mutex = 1
< ok 0 0
> exec ^
: string
| set a.mysql_errno = 1
< ok 0 0
> exec ^
: string
| set a.name = 1
< ok 0 0
> exec ^
: string
| set a.names = 1
< ok 0 0
> exec ^
: string
| set a.national = 1
< ok 0 0
> exec ^
: string
| set a.natural = 1
< ok 0 0
> exec ^
: string
| set a.nchar = 1
< ok 0 0
> exec ^
: string
| set a.ndb = 1
< ok 0 0
> exec ^
: string
| set a.ndbcluster = 1
< ok 0 0
> exec ^
: string
| set a.nested = 1
< ok 0 0
> exec ^
: string
| set a.network_namespace = 1
< ok 0 0
> exec ^
: string
| set a.never = 1
< ok 0 0
> exec ^
: string
| set a.new = 1
< ok 0 0
> exec ^
: string
| set a.next = 1
< ok 0 0
> exec ^
: string
| set a.next_row_id = 1
< ok 0 0
> exec ^
: string
| set a.no = 1
< ok 0 0
> exec ^
: string
| set a.no_wait = 1
< ok 0 0
> exec ^
: string
| set a.no_write_to_binlog = 1
< ok 0 0
> exec ^
: string
| set a.nocache = 1
< ok 0 0
> exec ^
: string
| set a.nocycle = 1
< ok 0 0
> exec ^
: string
| set a.nodegroup = 1
< ok 0 0
> exec ^
: string
| set a.nomaxvalue = 1
< ok 0 0
> exec ^
: string
| set a.nominvalue = 1
< ok 0 0
> exec ^
: string
| set a.nonclustered = 1
< ok 0 0
> exec ^
: string
| set a.none = 1
< ok 0 0
> exec ^
: string
| set a.not = 1
< ok 0 0
> exec ^
: string
| set a.now = 1
< ok 0 0
> exec ^
: string
| set a.nowait = 1
< ok 0 0
> exec ^
: string
| set a.nth_value = 1
< ok 0 0
> exec ^
: string
| set a.ntile = 1
< ok 0 0
> exec ^
: string
| set a.null = 1
< ok 0 0
> exec ^
: string
| set a.nullif = 1
< ok 0 0
> exec ^
: string
| set a.nulls = 1
< ok 0 0
> exec ^
: string
| set a.number = 1
< ok 0 0
> exec ^
: string
| set a.numeric = 1
< ok 0 0
> exec ^
: string
| set a.nvarchar = 1
< ok 0 0
> exec ^
: string
| set a.of = 1
< ok 0 0
> exec ^
: string
| set a.off = 1
< ok 0 0
> exec ^
: string
| set a.offset = 1
< ok 0 0
> exec ^
: string
| set a.oj = 1
< ok 0 0
> exec ^
: string
| set a.old = 1
< ok 0 0
> exec ^
: string
| set a.oltp_read_only = 1
< ok 0 0
> exec ^
: string
| set a.oltp_read_write = 1
< ok 0 0
> exec ^
: string
| set a.oltp_write_only = 1
< ok 0 0
> exec ^
: string
| set a.on = 1
< ok 0 0
> exec ^
: string
| set a.one = 1
< ok 0 0
> exec ^
: string
| set a.only = 1
< ok 0 0
> exec ^
: string
| set a.open = 1
< ok 0 0
> exec ^
: string
| set a.optimize = 1
< ok 0 0
> exec ^
: string
| set a.optimizer_costs = 1
< ok 0 0
> exec ^
: string
| set a.option = 1
< ok 0 0
> exec ^
: string
| set a.optional = 1
< ok 0 0
> exec ^
: string
| set a.optionally = 1
< ok 0 0
> exec ^
: string
| set a.options = 1
< ok 0 0
> exec ^
: string
| set a.or = 1
< ok 0 0
> exec ^
: string
| set a.order = 1
< ok 0 0
> exec ^
: string
| set a.ordinality = 1
< ok 0 0
> exec ^
: string
| set a.organization = 1
< ok 0 0
> exec ^
: string
| set a.others = 1
< ok 0 0
> exec ^
: string
| set a.out = 1
< ok 0 0
> exec ^
: string
| set a.outer = 1
< ok 0 0
> exec ^
: string
| set a.outfile = 1
< ok 0 0
> exec ^
: string
| set a.over = 1
< ok 0 0
> exec ^
: string
| set a.overlay = 1
< ok 0 0
> exec ^
: string
| set a.owner = 1
< ok 0 0
> exec ^
: string
| set a.pack_keys = 1
< ok 0 0
> exec ^
: string
| set a.page = 1
< ok 0 0
> exec ^
: string
| set a.parallel = 1
< ok 0 0
> exec ^
: string
| set a.parameters = 1
< ok 0 0
> exec ^
: string
| set a.parse_tree = 1
< ok 0 0
> exec ^
: string
| set a.parser = 1
< ok 0 0
> exec ^
: string
| set a.partial = 1
< ok 0 0
> exec ^
: string
| set a.partition = 1
< ok 0 0
> exec ^
: string
| set a.partitioning = 1
< ok 0 0
> exec ^
: string
| set a.partitions = 1
< ok 0 0
> exec ^
: string
| set a.password = 1
< ok 0 0
> exec ^
: string
| set a.password_lock_time = 1
< ok 0 0
> exec ^
: string
| set a.path = 1
< ok 0 0
> exec ^
: string
| set a.percent_rank = 1
< ok 0 0
> exec ^
: string
| set a.persist = 1
< ok 0 0
> exec ^
: string
| set a.persist_only = 1
< ok 0 0
> exec ^
: string
| set a.phase = 1
< ok 0 0
> exec ^
: string
| set a.placement = 1
< ok 0 0
> exec ^
: string
| set a.placing = 1
< ok 0 0
> exec ^
: string
| set a.plan = 1
< ok 0 0
> exec ^
: string
| set a.plugin = 1
< ok 0 0
> exec ^
: string
| set a.plugin_dir = 1
< ok 0 0
> exec ^
: string
| set a.plugins = 1
< ok 0 0
> exec ^
: string
| set a.point = 1
< ok 0 0
> exec ^
: string
| set a.policy = 1
< ok 0 0
> exec ^
: string
| set a.polygon = 1
< ok 0 0
> exec ^
: string
| set a.port = 1
< ok 0 0
> exec ^
: string
| set a.position = 1
< ok 0 0
> exec ^
: string
| set a.pre_split_regions = 1
< ok 0 0
> exec ^
: string
| set a.precedes = 1
< ok 0 0
> exec ^
: string
| set a.preceding = 1
< ok 0 0
> exec ^
: string
| set a.precision = 1
< ok 0 0
> exec ^
: string
| set a.prepare = 1
< ok 0 0
> exec ^
: string
| set a.preserve = 1
< ok 0 0
> exec ^
: string
| set a.prev = 1
< ok 0 0
> exec ^
: string
| set a.primary = 1
< ok 0 0
> exec ^
: string
| set a.primary_region = 1
< ok 0 0
> exec ^
: string
| set a.priority = 1
< ok 0 0
> exec ^
: string
| set a.privilege_checks_user = 1
< ok 0 0
> exec ^
: string
| set a.privileges = 1
< ok 0 0
> exec ^
: string
| set a.procedure = 1
< ok 0 0
> exec ^
: string
| set a.process = 1
< ok 0 0
> exec ^
: string
| set a.processed_keys = 1
< ok 0 0
> exec ^
: string
| set a.processlist = 1
< ok 0 0
> exec ^
: string
| set a.profile = 1
< ok 0 0
> exec ^
: string
| set a.profiles = 1
< ok 0 0
> exec ^
: string
| set a.proxy = 1
< ok 0 0
> exec ^
: string
| set a.purge = 1
< ok 0 0
> exec ^
: string
| set a.qualify = 1
< ok 0 0
> exec ^
: string
| set a.quarter = 1
< ok 0 0
> exec ^
: string
| set a.query = 1
< ok 0 0
> exec ^
: string
| set a.query_limit = 1
< ok 0 0
> exec ^
: string
| set a.quick = 1
< ok 0 0
> exec ^
: string
| set a.random = 1
< ok 0 0
> exec ^
: string
| set a.range = 1
< ok 0 0
> exec ^
: string
| set a.rank = 1
< ok 0 0
> exec ^
: string
| set a.read = 1
< ok 0 0
> exec ^
: string
| set a.read_only = 1
< ok 0 0
> exec ^
: string
| set a.read_write = 1
< ok 0 0
> exec ^
: string
| set a.reads = 1
< ok 0 0
> exec ^
: string
| set a.real = 1
< ok 0 0
> exec ^
: string
| set a.rebuild = 1
< ok 0 0
> exec ^
: string
| set a.recover = 1
< ok 0 0
> exec ^
: string
| set a.recursive = 1
< ok 0 0
> exec ^
: string
| set a.redo_buffer_size = 1
< ok 0 0
> exec ^
: string
| set a.redundant = 1
< ok 0 0
> exec ^
: string
| set a.reference = 1
< ok 0 0
> exec ^
: string
| set a.references = 1
< ok 0 0
> exec ^
: string
| set a.regexp = 1
< ok 0 0
> exec ^
: string
| set a.region = 1
< ok 0 0
> exec ^
: string
| set a.regions = 1
< ok 0 0
> exec ^
: string
| set a.registration = 1
< ok 0 0
> exec ^
: string
| set a.relational = 1
< ok 0 0
> exec ^
: string
| set a.relay = 1
< ok 0 0
> exec ^
: string
| set a.relay_log_file = 1
< ok 0 0
> exec ^
: string
| set a.relay_log_pos = 1
< ok 0 0
> exec ^
: string
| set a.relay_thread = 1
< ok 0 0
> exec ^
: string
| set a.relaylog = 1
< ok 0 0
> exec ^
: string
| set a.release = 1
< ok 0 0
> exec ^
: string
| set a.reload = 1
< ok 0 0
> exec ^
: string
| set a.remove = 1
< ok 0 0
> exec ^
: string
| set a.rename = 1
< ok 0 0
> exec ^
: string
| set a.reorganize = 1
< ok 0 0
> exec ^
: string
| set a.repair = 1
< ok 0 0
> exec ^
: string
| set a.repeat = 1
< ok 0 0
> exec ^
: string
| set a.repeatable = 1
< ok 0 0
> exec ^
: string
| set a.replace = 1
< ok 0 0
> exec ^
: string
| set a.replica = 1
< ok 0 0
> exec ^
: string
| set a.replicas = 1
< ok 0 0
> exec ^
: string
| set a.replicate_do_db = 1
< ok 0 0
> exec ^
: string
| set a.replicate_do_table = 1
< ok 0 0
> exec ^
: string
| set a.replicate_ignore_db = 1
< ok 0 0
> exec ^
: string
| set a.replicate_ignore_table = 1
< ok 0 0
> exec ^
: string
| set a.replicate_rewrite_db = 1
< ok 0 0
> exec ^
: string
| set a.replicate_wild_do_table = 1
< ok 0 0
> exec ^
: string
| set a.replicate_wild_ignore_table = 1
< ok 0 0
> exec ^
: string
| set a.replication = 1
< ok 0 0
> exec ^
: string
| set a.require = 1
< ok 0 0
> exec ^
: string
| set a.require_row_format = 1
< ok 0 0
> exec ^
: string
| set a.require_table_primary_key_check = 1
< ok 0 0
> exec ^
: string
| set a.reset = 1
< ok 0 0
> exec ^
: string
| set a.resignal = 1
< ok 0 0
> exec ^
: string
| set a.resource = 1
< ok 0 0
> exec ^
: string
| set a.respect = 1
< ok 0 0
> exec ^
: string
| set a.restart = 1
< ok 0 0
> exec ^
: string
| set a.restore = 1
< ok 0 0
> exec ^
: string
| set a.restores = 1
< ok 0 0
> exec ^
: string
| set a.restrict = 1
< ok 0 0
> exec ^
: string
| set a.resume = 1
< ok 0 0
> exec ^
: string
| set a.retain = 1
< ok 0 0
> exec ^
: string
| set a.return = 1
< ok 0 0
> exec ^
: string
| set a.returned_sqlstate = 1
< ok 0 0
> exec ^
: string
| set a.returning = 1
< ok 0 0
> exec ^
: string
| set a.returns = 1
< ok 0 0
> exec ^
: string
| set a.reuse = 1
< ok 0 0
> exec ^
: string
| set a.reverse = 1
< ok 0 0
> exec ^
: string
| set a.revoke = 1
< ok 0 0
> exec ^
: string
| set a.right = 1
< ok 0 0
> exec ^
: string
| set a.rlike = 1
< ok 0 0
> exec ^
: string
| set a.role = 1
< ok 0 0
> exec ^
: string
| set a.rollback = 1
< ok 0 0
> exec ^
: string
| set a.rollup = 1
< ok 0 0
> exec ^
: string
| set a.rotate = 1
< ok 0 0
> exec ^
: string
| set a.routine = 1
< ok 0 0
> exec ^
: string
| set a.row = 1
< ok 0 0
> exec ^
: string
| set a.row_count = 1
< ok 0 0
> exec ^
: string
| set a.row_format = 1
< ok 0 0
> exec ^
: string
| set a.row_number = 1
< ok 0 0
> exec ^
: string
| set a.rows = 1
< ok 0 0
> exec ^
: string
| set a.rtree = 1
< ok 0 0
> exec ^
: string
| set a.ru = 1
< ok 0 0
> exec ^
: string
| set a.ru_per_sec = 1
< ok 0 0
> exec ^
: string
| set a.rule = 1
< ok 0 0
> exec ^
: string
| set a.run = 1
< ok 0 0
> exec ^
: string
| set a.s3 = 1
< ok 0 0
> exec ^
: string
| set a.savepoint = 1
< ok 0 0
> exec ^
: string
| set a.schedule = 1
< ok 0 0
> exec ^
: string
| set a.schema = 1
< ok 0 0
> exec ^
: string
| set a.schema_name = 1
< ok 0 0
> exec ^
: string
| set a.schemas = 1
< ok 0 0
> exec ^
: string
| set a.second = 1
< ok 0 0
> exec ^
: string
| set a.second_microsecond = 1
< ok 0 0
> exec ^
: string
| set a.secondary = 1
< ok 0 0
> exec ^
: string
| set a.secondary_engine = 1
< ok 0 0
> exec ^
: string
| set a.secondary_engine_attribute = 1
< ok 0 0
> exec ^
: string
| set a.secondary_load = 1
< ok 0 0
> exec ^
: string
| set a.secondary_unload = 1
< ok 0 0
> exec ^
: string
| set a.security = 1
< ok 0 0
> exec ^
: string
| set a.select = 1
< ok 0 0
> exec ^
: string
| set a.sensitive = 1
< ok 0 0
> exec ^
: string
| set a.separator = 1
< ok 0 0
> exec ^
: string
| set a.sequence = 1
< ok 0 0
> exec ^
: string
| set a.serial = 1
< ok 0 0
> exec ^
: string
| set a.serializable = 1
< ok 0 0
> exec ^
: string
| set a.server = 1
< ok 0 0
> exec ^
: string
| set a.session = 1
< ok 0 0
> exec ^
: string
| set a.session_states = 1
< ok 0 0
> exec ^
: string
| set a.session_user = 1
< ok 0 0
> exec ^
: string
| set a.set = 1
< ok 0 0
> exec ^
: string
| set a.sets = 1
< ok 0 0
> exec ^
: string
| set a.shard_row_id_bits = 1
< ok 0 0
> exec ^
: string
| set a.share = 1
< ok 0 0
> exec ^
: string
| set a.show = 1
< ok 0 0
> exec ^
: string
| set a.shutdown = 1
< ok 0 0
> exec ^
: string
| set a.signal = 1
< ok 0 0
> exec ^
: string
| set a.signed = 1
< ok 0 0
> exec ^
: string
| set a.similar = 1
< ok 0 0
> exec ^
: string
| set a.simple = 1
< ok 0 0
> exec ^
: string
| set a.skip = 1
< ok 0 0
> exec ^
: string
| set a.slave = 1
< ok 0 0
> exec ^
: string
| set a.slow = 1
< ok 0 0
> exec ^
: string
| set a.smallint = 1
< ok 0 0
> exec ^
: string
| set a.snapshot = 1
< ok 0 0
> exec ^
: string
| set a.socket = 1
< ok 0 0
> exec ^
: string
| set a.some = 1
< ok 0 0
> exec ^
: string
| set a.soname = 1
< ok 0 0
> exec ^
: string
| set a.sounds = 1
< ok 0 0
> exec ^
: string
| set a.source = 1
< ok 0 0
> exec ^
: string
| set a.source_auto_position = 1
< ok 0 0
> exec ^
: string
| set a.source_bind = 1
< ok 0 0
> exec ^
: string
| set a.source_compression_algorithms = 1
< ok 0 0
> exec ^
: string
| set a.source_connect_retry = 1
< ok 0 0
> exec ^
: string
| set a.source_connection_auto_failover = 1
< ok 0 0
> exec ^
: string
| set a.source_delay = 1
< ok 0 0
> exec ^
: string
| set a.source_heartbeat_period = 1
< ok 0 0
> exec ^
: string
| set a.source_host = 1
< ok 0 0
> exec ^
: string
| set a.source_log_file = 1
< ok 0 0
> exec ^
: string
| set a.source_log_pos = 1
< ok 0 0
> exec ^
: string
| set a.source_password = 1
< ok 0 0
> exec ^
: string
| set a.source_port = 1
< ok 0 0
> exec ^
: string
| set a.source_public_key_path = 1
< ok 0 0
> exec ^
: string
| set a.source_retry_count = 1
< ok 0 0
> exec ^
: string
| set a.source_ssl = 1
< ok 0 0
> exec ^
: string
| set a.source_ssl_ca = 1
< ok 0 0
> exec ^
: string
| set a.source_ssl_capath = 1
< ok 0 0
> exec ^
: string
| set a.source_ssl_cert = 1
< ok 0 0
> exec ^
: string
| set a.source_ssl_cipher = 1
< ok 0 0
> exec ^
: string
| set a.source_ssl_crl = 1
< ok 0 0
> exec ^
: string
| set a.source_ssl_crlpath = 1
< ok 0 0
> exec ^
: string
| set a.source_ssl_key = 1
< ok 0 0
> exec ^
: string
| set a.source_ssl_verify_server_cert = 1
< ok 0 0
> exec ^
: string
| set a.source_tls_ciphersuites = 1
< ok 0 0
> exec ^
: string
| set a.source_tls_version = 1
< ok 0 0
> exec ^
: string
| set a.source_user = 1
< ok 0 0
> exec ^
: string
| set a.source_zstd_compression_level = 1
< ok 0 0
> exec ^
: string
| set a.spatial = 1
< ok 0 0
> exec ^
: string
| set a.specific = 1
< ok 0 0
> exec ^
: string
| set a.split = 1
< ok 0 0
> exec ^
: string
| set a.sql = 1
< ok 0 0
> exec ^
: string
| set a.sql_after_gtids = 1
< ok 0 0
> exec ^
: string
| set a.sql_after_mts_gaps = 1
< ok 0 0
> exec ^
: string
| set a.sql_before_gtids = 1
< ok 0 0
> exec ^
: string
| set a.sql_big_result = 1
< ok 0 0
> exec ^
: string
| set a.sql_buffer_result = 1
< ok 0 0
> exec ^
: string
| set a.sql_calc_found_rows = 1
< ok 0 0
> exec ^
: string
| set a.sql_no_cache = 1
< ok 0 0
> exec ^
: string
| set a.sql_small_result = 1
< ok 0 0
> exec ^
: string
| set a.sql_thread = 1
< ok 0 0
> exec ^
: string
| set a.sql_tsi_day = 1
< ok 0 0
> exec ^
: string
| set a.sql_tsi_hour = 1
< ok 0 0
> exec ^
: string
| set a.sql_tsi_minute = 1
< ok 0 0
> exec ^
: string
| set a.sql_tsi_month = 1
< ok 0 0
> exec ^
: string
| set a.sql_tsi_quarter = 1
< ok 0 0
> exec ^
: string
| set a.sql_tsi_second = 1
< ok 0 0
> exec ^
: string
| set a.sql_tsi_week = 1
< ok 0 0
> exec ^
: string
| set a.sql_tsi_year = 1
< ok 0 0
> exec ^
: string
| set a.sqlexception = 1
< ok 0 0
> exec ^
: string
| set a.sqlstate = 1
< ok 0 0
> exec ^
: string
| set a.sqlwarning = 1
< ok 0 0
> exec ^
: string
| set a.srid = 1
< ok 0 0
> exec ^
: string
| set a.ssl = 1
< ok 0 0
> exec ^
: string
| set a.st_collect = 1
< ok 0 0
> exec ^
: string
| set a.stacked = 1
< ok 0 0
> exec ^
: string
| set a.start = 1
< ok 0 0
> exec ^
: string
| set a.start_time = 1
< ok 0 0
> exec ^
: string
| set a.starting = 1
< ok 0 0
> exec ^
: string
| set a.starts = 1
< ok 0 0
> exec ^
: string
| set a.stats_auto_recalc = 1
< ok 0 0
> exec ^
: string
| set a.stats_buckets = 1
< ok 0 0
> exec ^
: string
| set a.stats_col_choice = 1
< ok 0 0
> exec ^
: string
| set a.stats_col_list = 1
< ok 0 0
> exec ^
: string
| set a.stats_extended = 1
< ok 0 0
> exec ^
: string
| set a.stats_healthy = 1
< ok 0 0
> exec ^
: string
| set a.stats_histograms = 1
< ok 0 0
> exec ^
: string
| set a.stats_locked = 1
< ok 0 0
> exec ^
: string
| set a.stats_meta = 1
< ok 0 0
> exec ^
: string
| set a.stats_persistent = 1
< ok 0 0
> exec ^
: string
| set a.stats_sample_pages = 1
< ok 0 0
> exec ^
: string
| set a.stats_sample_rate = 1
< ok 0 0
> exec ^
: string
| set a.stats_topn = 1
< ok 0 0
> exec ^
: string
| set a.status = 1
< ok 0 0
> exec ^
: string
| set a.std = 1
< ok 0 0
> exec ^
: string
| set a.stddev = 1
< ok 0 0
> exec ^
: string
| set a.stddev_pop = 1
< ok 0 0
> exec ^
: string
| set a.stddev_samp = 1
< ok 0 0
> exec ^
: string
| set a.stop = 1
< ok 0 0
> exec ^
: string
| set a.storage = 1
< ok 0 0
> exec ^
: string
| set a.stored = 1
< ok 0 0
> exec ^
: string
| set a.straight_join = 1
< ok 0 0
> exec ^
: string
| set a.stream = 1
< ok 0 0
> exec ^
: string
| set a.strict_load = 1
< ok 0 0
> exec ^
: string
| set a.string = 1
< ok 0 0
> exec ^
: string
| set a.subclass_origin = 1
< ok 0 0
> exec ^
: string
| set a.subdate = 1
< ok 0 0
> exec ^
: string
| set a.subject = 1
< ok 0 0
> exec ^
: string
| set a.subpartition = 1
< ok 0 0
> exec ^
: string
| set a.subpartitions = 1
< ok 0 0
> exec ^
: string
| set a.substr = 1
< ok 0 0
> exec ^
: string
| set a.substring = 1
< ok 0 0
> exec ^
: string
| set a.sum = 1
< ok 0 0
> exec ^
: string
| set a.super = 1
< ok 0 0
> exec ^
: string
| set a.survival_preferences = 1
< ok 0 0
> exec ^
: string
| set a.suspend = 1
< ok 0 0
> exec ^
: string
| set a.swaps = 1
< ok 0 0
> exec ^
: string
| set a.switch_group = 1
< ok 0 0
> exec ^
: string
| set a.switches = 1
< ok 0 0
> exec ^
: string
| set a.sysdate = 1
< ok 0 0
> exec ^
: string
| set a.system = 1
< ok 0 0
> exec ^
: string
| set a.system_user = 1
< ok 0 0
> exec ^
: string
| set a.table = 1
< ok 0 0
> exec ^
: string
| set a.table_checksum = 1
< ok 0 0
> exec ^
: string
| set a.table_name = 1
< ok 0 0
> exec ^
: string
| set a.tables = 1
< ok 0 0
> exec ^
: string
| set a.tablesample = 1
< ok 0 0
> exec ^
: string
| set a.tablespace = 1
< ok 0 0
> exec ^
: string
| set a.task_types = 1
< ok 0 0
> exec ^
: string
| set a.temporary = 1
< ok 0 0
> exec ^
: string
| set a.temptable = 1
< ok 0 0
> exec ^
: string
| set a.terminated = 1
< ok 0 0
> exec ^
: string
| set a.text = 1
< ok 0 0
> exec ^
: string
| set a.than = 1
< ok 0 0
> exec ^
: string
| set a.then = 1
< ok 0 0
> exec ^
: string
| set a.thread_priority = 1
< ok 0 0
> exec ^
: string
| set a.ties = 1
< ok 0 0
> exec ^
: string
| set a.tiflash = 1
< ok 0 0
> exec ^
: string
| set a.time = 1
< ok 0 0
> exec ^
: string
| set a.timeout = 1
< ok 0 0
> exec ^
: string
| set a.timestamp = 1
< ok 0 0
> exec ^
: string
| set a.timestampadd = 1
< ok 0 0
> exec ^
: string
| set a.timestampdiff = 1
< ok 0 0
> exec ^
: string
| set a.tinyblob = 1
< ok 0 0
> exec ^
: string
| set a.tinyint = 1
< ok 0 0
> exec ^
: string
| set a.tinytext = 1
< ok 0 0
> exec ^
: string
| set a.tls = 1
< ok 0 0
> exec ^
: string
| set a.to = 1
< ok 0 0
> exec ^
: string
| set a.tpcc = 1
< ok 0 0
> exec ^
: string
| set a.tpch_10 = 1
< ok 0 0
> exec ^
: string
| set a.trailing = 1
< ok 0 0
> exec ^
: string
| set a.transaction = 1
< ok 0 0
> exec ^
: string
| set a.trigger = 1
< ok 0 0
> exec ^
: string
| set a.triggers = 1
< ok 0 0
> exec ^
: string
| set a.trim = 1
< ok 0 0
> exec ^
: string
| set a.true = 1
< ok 0 0
> exec ^
: string
| set a.truncate = 1
< ok 0 0
> exec ^
: string
| set a.ttl = 1
< ok 0 0
> exec ^
: string
| set a.ttl_enable = 1
< ok 0 0
> exec ^
: string
| set a.ttl_job_interval = 1
< ok 0 0
> exec ^
: string
| set a.type = 1
< ok 0 0
> exec ^
: string
| set a.types = 1
< ok 0 0
> exec ^
: string
| set a.unbounded = 1
< ok 0 0
> exec ^
: string
| set a.uncommitted = 1
< ok 0 0
> exec ^
: string
| set a.undefined = 1
< ok 0 0
> exec ^
: string
| set a.undo = 1
< ok 0 0
> exec ^
: string
| set a.undo_buffer_size = 1
< ok 0 0
> exec ^
: string
| set a.undofile = 1
< ok 0 0
> exec ^
: string
| set a.unicode = 1
< ok 0 0
> exec ^
: string
| set a.uninstall = 1
< ok 0 0
> exec ^
: string
| set a.union = 1
< ok 0 0
> exec ^
: string
| set a.unique = 1
< ok 0 0
> exec ^
: string
| set a.unknown = 1
< ok 0 0
> exec ^
: string
| set a.unlimited = 1
< ok 0 0
> exec ^
: string
| set a.unlock = 1
< ok 0 0
> exec ^
: string
| set a.unregister = 1
< ok 0 0
> exec ^
: string
| set a.unsigned = 1
< ok 0 0
> exec ^
: string
| set a.until = 1
< ok 0 0
> exec ^
: string
| set a.update = 1
< ok 0 0
> exec ^
: string
| set a.upgrade = 1
< ok 0 0
> exec ^
: string
| set a.uri = 1
< ok 0 0
> exec ^
: string
| set a.url = 1
< ok 0 0
> exec ^
: string
| set a.usage = 1
< ok 0 0
> exec ^
: string
| set a.use = 1
< ok 0 0
> exec ^
: string
| set a.use_frm = 1
< ok 0 0
> exec ^
: string
| set a.user = 1
< ok 0 0
> exec ^
: string
| set a.user_resources = 1
< ok 0 0
> exec ^
: string
| set a.using = 1
< ok 0 0
> exec ^
: string
| set a.utc_date = 1
< ok 0 0
> exec ^
: string
| set a.utc_time = 1
< ok 0 0
> exec ^
: string
| set a.utc_timestamp = 1
< ok 0 0
> exec ^
: string
| set a.utilization_limit = 1
< ok 0 0
> exec ^
: string
| set a.validate = 1
< ok 0 0
> exec ^
: string
| set a.validation = 1
< ok 0 0
> exec ^
: string
| set a.value = 1
< ok 0 0
> exec ^
: string
| set a.values = 1
< ok 0 0
> exec ^
: string
| set a.var_pop = 1
< ok 0 0
> exec ^
: string
| set a.var_samp = 1
< ok 0 0
> exec ^
: string
| set a.varbinary = 1
< ok 0 0
> exec ^
: string
| set a.varchar = 1
< ok 0 0
> exec ^
: string
| set a.varcharacter = 1
< ok 0 0
> exec ^
: string
| set a.variables = 1
< ok 0 0
> exec ^
: string
| set a.variance = 1
< ok 0 0
> exec ^
: string
| set a.varying = 1
< ok 0 0
> exec ^
: string
| set a.vcpu = 1
< ok 0 0
> exec ^
: string
| set a.vector = 1
< ok 0 0
> exec ^
: string
| set a.verify_key_constraints = 1
< ok 0 0
> exec ^
: string
| set a.view = 1
< ok 0 0
> exec ^
: string
| set a.virtual = 1
< ok 0 0
> exec ^
: string
| set a.visible = 1
< ok 0 0
> exec ^
: string
| set a.voter_constraints = 1
< ok 0 0
> exec ^
: string
| set a.voters = 1
< ok 0 0
> exec ^
: string
| set a.wait = 1
< ok 0 0
> exec ^
: string
| set a.warnings = 1
< ok 0 0
> exec ^
: string
| set a.watch = 1
< ok 0 0
> exec ^
: string
| set a.week = 1
< ok 0 0
> exec ^
: string
| set a.weight_string = 1
< ok 0 0
> exec ^
: string
| set a.when = 1
< ok 0 0
> exec ^
: string
| set a.where = 1
< ok 0 0
> exec ^
: string
| set a.while = 1
< ok 0 0
> exec ^
: string
| set a.window = 1
< ok 0 0
> exec ^
: string
| set a.with = 1
< ok 0 0
> exec ^
: string
| set a.without = 1
< ok 0 0
> exec ^
: string
| set a.work = 1
< ok 0 0
> exec ^
: string
| set a.workload = 1
< ok 0 0
> exec ^
: string
| set a.wrapper = 1
< ok 0 0
> exec ^
: string
| set a.write = 1
< ok 0 0
> exec ^
: string
| set a.x509 = 1
< ok 0 0
> exec ^
: string
| set a.xa = 1
< ok 0 0
> exec ^
: string
| set a.xid = 1
< ok 0 0
> exec ^
: string
| set a.xml = 1
< ok 0 0
> exec ^
: string
| set a.xor = 1
< ok 0 0
> exec ^
: string
| set a.year = 1
< ok 0 0
> exec ^
: string
| set a.year_month = 1
< ok 0 0
> exec ^
: string
| set a.zerofill = 1
< ok 0 0
> exec ^
: string
| set a.zone = 1
< ok 0 0
> exec ^
: string
| set a = 1, absent = 2
< ok 0 0
> exec ^
: string
| set a = 1, accessible = 2
< ok 0 0
> exec ^
: string
| set a = 1, account = 2
< ok 0 0
> exec ^
: string
| set a = 1, action = 2
< ok 0 0
> exec ^
: string
| set a = 1, active = 2
< ok 0 0
> exec ^
: string
| set a = 1, add = 2
< error
| 1064 42000 …line 1 column 14 near "add = 2" 
> exec ^
: string
| set a = 1, adddate = 2
< ok 0 0
> exec ^
: string
| set a = 1, admin = 2
< ok 0 0
> exec ^
: string
| set a = 1, after = 2
< ok 0 0
> exec ^
: string
| set a = 1, against = 2
< ok 0 0
> exec ^
: string
| set a = 1, aggregate = 2
< ok 0 0
> exec ^
: string
| set a = 1, algorithm = 2
< ok 0 0
> exec ^
: string
| set a = 1, all = 2
< error
| 1064 42000 …line 1 column 14 near "all = 2" 
> exec ^
: string
| set a = 1, allow_missing_files = 2
< ok 0 0
> exec ^
: string
| set a = 1, alter = 2
< error
| 1064 42000 …line 1 column 16 near "alter = 2" 
> exec ^
: string
| set a = 1, always = 2
< ok 0 0
> exec ^
: string
| set a = 1, analyze = 2
< error
| 1064 42000 …line 1 column 18 near "analyze = 2" 
> exec ^
: string
| set a = 1, and = 2
< error
| 1064 42000 …line 1 column 14 near "and = 2" 
> exec ^
: string
| set a = 1, any = 2
< ok 0 0
> exec ^
: string
| set a = 1, array = 2
< error
| 1064 42000 …line 1 column 16 near "array = 2" 
> exec ^
: string
| set a = 1, as = 2
< error
| 1064 42000 …line 1 column 13 near "as = 2" 
> exec ^
: string
| set a = 1, asc = 2
< error
| 1064 42000 …line 1 column 14 near "asc = 2" 
> exec ^
: string
| set a = 1, ascii = 2
< ok 0 0
> exec ^
: string
| set a = 1, asensitive = 2
< ok 0 0
> exec ^
: string
| set a = 1, assign_gtids_to_anonymous_transactions = 2
< ok 0 0
> exec ^
: string
| set a = 1, at = 2
< ok 0 0
> exec ^
: string
| set a = 1, attribute = 2
< ok 0 0
> exec ^
: string
| set a = 1, authentication = 2
< ok 0 0
> exec ^
: string
| set a = 1, auto = 2
< ok 0 0
> exec ^
: string
| set a = 1, auto_id_cache = 2
< ok 0 0
> exec ^
: string
| set a = 1, auto_increment = 2
< ok 0 0
> exec ^
: string
| set a = 1, auto_random = 2
< ok 0 0
> exec ^
: string
| set a = 1, auto_random_base = 2
< ok 0 0
> exec ^
: string
| set a = 1, auto_refresh = 2
< ok 0 0
> exec ^
: string
| set a = 1, auto_refresh_source = 2
< ok 0 0
> exec ^
: string
| set a = 1, autoextend_size = 2
< ok 0 0
> exec ^
: string
| set a = 1, avg = 2
< ok 0 0
> exec ^
: string
| set a = 1, avg_row_length = 2
< ok 0 0
> exec ^
: string
| set a = 1, background = 2
< ok 0 0
> exec ^
: string
| set a = 1, backup = 2
< ok 0 0
> exec ^
: string
| set a = 1, backups = 2
< ok 0 0
> exec ^
: string
| set a = 1, batch = 2
< ok 0 0
> exec ^
: string
| set a = 1, before = 2
< ok 0 0
> exec ^
: string
| set a = 1, begin = 2
< ok 0 0
> exec ^
: string
| set a = 1, bernoulli = 2
< ok 0 0
> exec ^
: string
| set a = 1, between = 2
< error
| 1064 42000 …line 1 column 18 near "between = 2" 
> exec ^
: string
| set a = 1, bigint = 2
< error
| 1064 42000 …line 1 column 17 near "bigint = 2" 
> exec ^
: string
| set a = 1, binary = 2
< error
| 1064 42000 …line 1 column 17 near "binary = 2" 
> exec ^
: string
| set a = 1, binding = 2
< ok 0 0
> exec ^
: string
| set a = 1, binding_cache = 2
< ok 0 0
> exec ^
: string
| set a = 1, bindings = 2
< ok 0 0
> exec ^
: string
| set a = 1, binlog = 2
< ok 0 0
> exec ^
: string
| set a = 1, bit = 2
< ok 0 0
> exec ^
: string
| set a = 1, bit_and = 2
< ok 0 0
> exec ^
: string
| set a = 1, bit_or = 2
< ok 0 0
> exec ^
: string
| set a = 1, bit_xor = 2
< ok 0 0
> exec ^
: string
| set a = 1, blob = 2
< error
| 1064 42000 …line 1 column 15 near "blob = 2" 
> exec ^
: string
| set a = 1, block = 2
< ok 0 0
> exec ^
: string
| set a = 1, bool = 2
< ok 0 0
> exec ^
: string
| set a = 1, boolean = 2
< ok 0 0
> exec ^
: string
| set a = 1, both = 2
< error
| 1064 42000 …line 1 column 15 near "both = 2" 
> exec ^
: string
| set a = 1, btree = 2
< ok 0 0
> exec ^
: string
| set a = 1, buckets = 2
< ok 0 0
> exec ^
: string
| set a = 1, builtins = 2
< ok 0 0
> exec ^
: string
| set a = 1, bulk = 2
< ok 0 0
> exec ^
: string
| set a = 1, burstable = 2
< ok 0 0
> exec ^
: string
| set a = 1, by = 2
< error
| 1064 42000 …line 1 column 13 near "by = 2" 
> exec ^
: string
| set a = 1, byte = 2
< ok 0 0
> exec ^
: string
| set a = 1, cache = 2
< ok 0 0
> exec ^
: string
| set a = 1, calibrate = 2
< ok 0 0
> exec ^
: string
| set a = 1, call = 2
< error
| 1064 42000 …line 1 column 15 near "call = 2" 
> exec ^
: string
| set a = 1, cascade = 2
< error
| 1064 42000 …line 1 column 18 near "cascade = 2" 
> exec ^
: string
| set a = 1, cascaded = 2
< ok 0 0
> exec ^
: string
| set a = 1, case = 2
< error
| 1064 42000 …line 1 column 15 near "case = 2" 
> exec ^
: string
| set a = 1, cast = 2
< ok 0 0
> exec ^
: string
| set a = 1, catalog_name = 2
< ok 0 0
> exec ^
: string
| set a = 1, chain = 2
< ok 0 0
> exec ^
: string
| set a = 1, challenge_response = 2
< ok 0 0
> exec ^
: string
| set a = 1, change = 2
< error
| 1064 42000 …line 1 column 17 near "change = 2" 
> exec ^
: string
| set a = 1, changed = 2
< ok 0 0
> exec ^
: string
| set a = 1, channel = 2
< ok 0 0
> exec ^
: string
| set a = 1, char = 2
< error
| 1064 42000 …line 1 column 17 near "= 2" 
> exec ^
: string
| set a = 1, character = 2
< error
| 1064 42000 …line 1 column 22 near "= 2" 
> exec ^
: string
| set a = 1, charset = 2
< ok 0 0
> exec ^
: string
| set a = 1, check = 2
< error
| 1064 42000 …line 1 column 16 near "check = 2" 
> exec ^
: string
| set a = 1, checksum = 2
< ok 0 0
> exec ^
: string
| set a = 1, cipher = 2
< ok 0 0
> exec ^
: string
| set a = 1, class_origin = 2
< ok 0 0
> exec ^
: string
| set a = 1, client = 2
< ok 0 0
> exec ^
: string
| set a = 1, clone = 2
< ok 0 0
> exec ^
: string
| set a = 1, close = 2
< ok 0 0
> exec ^
: string
| set a = 1, clustered = 2
< ok 0 0
> exec ^
: string
| set a = 1, coalesce = 2
< ok 0 0
> exec ^
: string
| set a = 1, code = 2
< ok 0 0
> exec ^
: string
| set a = 1, collate = 2
< error
| 1064 42000 …line 1 column 18 near "collate = 2" 
> exec ^
: string
| set a = 1, collation = 2
< ok 0 0
> exec ^
: string
| set a = 1, column = 2
< error
| 1064 42000 …line 1 column 17 near "column = 2" 
> exec ^
: string
| set a = 1, column_format = 2
< ok 0 0
> exec ^
: string
| set a = 1, column_name = 2
< ok 0 0
> exec ^
: string
| set a = 1, column_stats_usage = 2
< ok 0 0
> exec ^
: string
| set a = 1, columns = 2
< ok 0 0
> exec ^
: string
| set a = 1, comment = 2
< ok 0 0
> exec ^
: string
| set a = 1, commit = 2
< ok 0 0
> exec ^
: string
| set a = 1, committed = 2
< ok 0 0
> exec ^
: string
| set a = 1, compact = 2
< ok 0 0
> exec ^
: string
| set a = 1, completion = 2
< ok 0 0
> exec ^
: string
| set a = 1, component = 2
< ok 0 0
> exec ^
: string
| set a = 1, compressed = 2
< ok 0 0
> exec ^
: string
| set a = 1, compression = 2
< ok 0 0
> exec ^
: string
| set a = 1, concat = 2
< ok 0 0
> exec ^
: string
| set a = 1, concurrent = 2
< ok 0 0
> exec ^
: string
| set a = 1, condition = 2
< ok 0 0
> exec ^
: string
| set a = 1, config = 2
< ok 0 0
> exec ^
: string
| set a = 1, connection = 2
< ok 0 0
> exec ^
: string
| set a = 1, consistent = 2
< ok 0 0
> exec ^
: string
| set a = 1, constraint = 2
< error
| 1064 42000 …line 1 column 21 near "constraint = 2" 
> exec ^
: string
| set a = 1, constraint_catalog = 2
< ok 0 0
> exec ^
: string
| set a = 1, constraint_name = 2
< ok 0 0
> exec ^
: string
| set a = 1, constraint_schema = 2
< ok 0 0
> exec ^
: string
| set a = 1, constraints = 2
< ok 0 0
> exec ^
: string
| set a = 1, contains = 2
< ok 0 0
> exec ^
: string
| set a = 1, context = 2
< ok 0 0
> exec ^
: string
| set a = 1, continue = 2
< error
| 1064 42000 …line 1 column 19 near "continue = 2" 
> exec ^
: string
| set a = 1, convert = 2
< error
| 1064 42000 …line 1 column 18 near "convert = 2" 
> exec ^
: string
| set a = 1, cooldown = 2
< ok 0 0
> exec ^
: string
| set a = 1, count = 2
< ok 0 0
> exec ^
: string
| set a = 1, cpu = 2
< ok 0 0
> exec ^
: string
| set a = 1, create = 2
< error
| 1064 42000 …line 1 column 17 near "create = 2" 
> exec ^
: string
| set a = 1, cross = 2
< error
| 1064 42000 …line 1 column 16 near "cross = 2" 
> exec ^
: string
| set a = 1, cube = 2
< ok 0 0
> exec ^
: string
| set a = 1, cume_dist = 2
< error
| 1064 42000 …line 1 column 20 near "cume_dist = 2" 
> exec ^
: string
| set a = 1, curdate = 2
< ok 0 0
> exec ^
: string
| set a = 1, current = 2
< ok 0 0
> exec ^
: string
| set a = 1, current_date = 2
< error
| 1064 42000 …line 1 column 23 near "current_date = 2" 
> exec ^
: string
| set a = 1, current_time = 2
< error
| 1064 42000 …line 1 column 23 near "current_time = 2" 
> exec ^
: string
| set a = 1, current_timestamp = 2
< error
| 1064 42000 …line 1 column 28 near "current_timestamp = 2" 
> exec ^
: string
| set a = 1, current_user = 2
< error
| 1064 42000 …line 1 column 23 near "current_user = 2" 
> exec ^
: string
| set a = 1, cursor = 2
< error
| 1064 42000 …line 1 column 17 near "cursor = 2" 
> exec ^
: string
| set a = 1, cursor_name = 2
< ok 0 0
> exec ^
: string
| set a = 1, curtime = 2
< ok 0 0
> exec ^
: string
| set a = 1, cycle = 2
< ok 0 0
> exec ^
: string
| set a = 1, data = 2
< ok 0 0
> exec ^
: string
| set a = 1, database = 2
< error
| 1064 42000 …line 1 column 19 near "database = 2" 
> exec ^
: string
| set a = 1, databases = 2
< error
| 1064 42000 …line 1 column 20 near "databases = 2" 
> exec ^
: string
| set a = 1, datafile = 2
< ok 0 0
> exec ^
: string
| set a = 1, date = 2
< ok 0 0
> exec ^
: string
| set a = 1, date_add = 2
< ok 0 0
> exec ^
: string
| set a = 1, date_sub = 2
< ok 0 0
> exec ^
: string
| set a = 1, datetime = 2
< ok 0 0
> exec ^
: string
| set a = 1, day = 2
< ok 0 0
> exec ^
: string
| set a = 1, day_hour = 2
< error
| 1064 42000 …line 1 column 19 near "day_hour = 2" 
> exec ^
: string
| set a = 1, day_microsecond = 2
< error
| 1064 42000 …line 1 column 26 near "day_microsecond = 2" 
> exec ^
: string
| set a = 1, day_minute = 2
< error
| 1064 42000 …line 1 column 21 near "day_minute = 2" 
> exec ^
: string
| set a = 1, day_second = 2
< error
| 1064 42000 …line 1 column 21 near "day_second = 2" 
> exec ^
: string
| set a = 1, deallocate = 2
< ok 0 0
> exec ^
: string
| set a = 1, dec = 2
< error
| 1064 42000 …line 1 column 14 near "dec = 2" 
> exec ^
: string
| set a = 1, decimal = 2
< error
| 1064 42000 …line 1 column 18 near "decimal = 2" 
> exec ^
: string
| set a = 1, declare = 2
< ok 0 0
> exec ^
: string
| set a = 1, default = 2
< error
| 1064 42000 …line 1 column 18 near "default = 2" 
> exec ^
: string
| set a = 1, default_auth = 2
< ok 0 0
> exec ^
: string
| set a = 1, defined = 2
< ok 0 0
> exec ^
: string
| set a = 1, definer = 2
< ok 0 0
> exec ^
: string
| set a = 1, definition = 2
< ok 0 0
> exec ^
: string
| set a = 1, delay_key_write = 2
< ok 0 0
> exec ^
: string
| set a = 1, delayed = 2
< error
| 1064 42000 …line 1 column 18 near "delayed = 2" 
> exec ^
: string
| set a = 1, delete = 2
< error
| 1064 42000 …line 1 column 17 near "delete = 2" 
> exec ^
: string
| set a = 1, dense_rank = 2
< error
| 1064 42000 …line 1 column 21 near "dense_rank = 2" 
> exec ^
: string
| set a = 1, desc = 2
< error
| 1064 42000 …line 1 column 15 near "desc = 2" 
> exec ^
: string
| set a = 1, describe = 2
< error
| 1064 42000 …line 1 column 19 near "describe = 2" 
> exec ^
: string
| set a = 1, description = 2
< ok 0 0
> exec ^
: string
| set a = 1, deterministic = 2
< ok 0 0
> exec ^
: string
| set a = 1, diagnostics = 2
< ok 0 0
> exec ^
: string
| set a = 1, digest = 2
< ok 0 0
> exec ^
: string
| set a = 1, directory = 2
< ok 0 0
> exec ^
: string
| set a = 1, disable = 2
< ok 0 0
> exec ^
: string
| set a = 1, disabled = 2
< ok 0 0
> exec ^
: string
| set a = 1, discard = 2
< ok 0 0
> exec ^
: string
| set a = 1, disk = 2
< ok 0 0
> exec ^
: string
| set a = 1, distinct = 2
< error
| 1064 42000 …line 1 column 19 near "distinct = 2" 
> exec ^
: string
| set a = 1, distinctrow = 2
< error
| 1064 42000 …line 1 column 22 near "distinctrow = 2" 
> exec ^
: string
| set a = 1, distribute = 2
< ok 0 0
> exec ^
: string
| set a = 1, distribution = 2
< ok 0 0
> exec ^
: string
| set a = 1, distributions = 2
< ok 0 0
> exec ^
: string
| set a = 1, div = 2
< error
| 1064 42000 …line 1 column 14 near "div = 2" 
> exec ^
: string
| set a = 1, do = 2
< ok 0 0
> exec ^
: string
| set a = 1, double = 2
< error
| 1064 42000 …line 1 column 17 near "double = 2" 
> exec ^
: string
| set a = 1, drop = 2
< error
| 1064 42000 …line 1 column 15 near "drop = 2" 
> exec ^
: string
| set a = 1, dry = 2
< ok 0 0
> exec ^
: string
| set a = 1, dryrun = 2
< ok 0 0
> exec ^
: string
| set a = 1, dual = 2
< error
| 1064 42000 …line 1 column 15 near "dual = 2" 
> exec ^
: string
| set a = 1, duality = 2
< ok 0 0
> exec ^
: string
| set a = 1, dumpfile = 2
< ok 0 0
> exec ^
: string
| set a = 1, duplicate = 2
< ok 0 0
> exec ^
: string
| set a = 1, duration = 2
< ok 0 0
> exec ^
: string
| set a = 1, dynamic = 2
< ok 0 0
> exec ^
: string
| set a = 1, each = 2
< ok 0 0
> exec ^
: string
| set a = 1, else = 2
< error
| 1064 42000 …line 1 column 15 near "else = 2" 
> exec ^
: string
| set a = 1, elseif = 2
< error
| 1064 42000 …line 1 column 17 near "elseif = 2" 
> exec ^
: string
| set a = 1, empty = 2
< ok 0 0
> exec ^
: string
| set a = 1, enable = 2
< ok 0 0
> exec ^
: string
| set a = 1, enabled = 2
< ok 0 0
> exec ^
: string
| set a = 1, enclosed = 2
< error
| 1064 42000 …line 1 column 19 near "enclosed = 2" 
> exec ^
: string
| set a = 1, encryption = 2
< ok 0 0
> exec ^
: string
| set a = 1, end = 2
< ok 0 0
> exec ^
: string
| set a = 1, end_time = 2
< ok 0 0
> exec ^
: string
| set a = 1, ends = 2
< ok 0 0
> exec ^
: string
| set a = 1, enforced = 2
< ok 0 0
> exec ^
: string
| set a = 1, engine = 2
< ok 0 0
> exec ^
: string
| set a = 1, engine_attribute = 2
< ok 0 0
> exec ^
: string
| set a = 1, engines = 2
< ok 0 0
> exec ^
: string
| set a = 1, enum = 2
< ok 0 0
> exec ^
: string
| set a = 1, error = 2
< ok 0 0
> exec ^
: string
| set a = 1, errors = 2
< ok 0 0
> exec ^
: string
| set a = 1, escape = 2
< ok 0 0
> exec ^
: string
| set a = 1, escaped = 2
< error
| 1064 42000 …line 1 column 18 near "escaped = 2" 
> exec ^
: string
| set a = 1, event = 2
< ok 0 0
> exec ^
: string
| set a = 1, events = 2
< ok 0 0
> exec ^
: string
| set a = 1, every = 2
< ok 0 0
> exec ^
: string
| set a = 1, exact = 2
< ok 0 0
> exec ^
: string
| set a = 1, except = 2
< error
| 1064 42000 …line 1 column 17 near "except = 2" 
> exec ^
: string
| set a = 1, exchange = 2
< ok 0 0
> exec ^
: string
| set a = 1, exclude = 2
< ok 0 0
> exec ^
: string
| set a = 1, exec_elapsed = 2
< ok 0 0
> exec ^
: string
| set a = 1, execute = 2
< ok 0 0
> exec ^
: string
| set a = 1, exists = 2
< error
| 1064 42000 …line 1 column 17 near "exists = 2" 
> exec ^
: string
| set a = 1, exit = 2
< error
| 1064 42000 …line 1 column 15 near "exit = 2" 
> exec ^
: string
| set a = 1, expansion = 2
< ok 0 0
> exec ^
: string
| set a = 1, expire = 2
< ok 0 0
> exec ^
: string
| set a = 1, explain = 2
< error
| 1064 42000 …line 1 column 18 near "explain = 2" 
> exec ^
: string
| set a = 1, export = 2
< ok 0 0
> exec ^
: string
| set a = 1, extended = 2
< ok 0 0
> exec ^
: string
| set a = 1, extent_size = 2
< ok 0 0
> exec ^
: string
| set a = 1, external = 2
< ok 0 0
> exec ^
: string
| set a = 1, external_format = 2
< ok 0 0
> exec ^
: string
| set a = 1, extract = 2
< ok 0 0
> exec ^
: string
| set a = 1, factor = 2
< ok 0 0
> exec ^
: string
| set a = 1, failed_login_attempts = 2
< ok 0 0
> exec ^
: string
| set a = 1, false = 2
< error
| 1064 42000 …line 1 column 16 near "false = 2" 
> exec ^
: string
| set a = 1, fast = 2
< ok 0 0
> exec ^
: string
| set a = 1, faults = 2
< ok 0 0
> exec ^
: string
| set a = 1, fetch = 2
< error
| 1064 42000 …line 1 column 16 near "fetch = 2" 
> exec ^
: string
| set a = 1, fields = 2
< ok 0 0
> exec ^
: string
| set a = 1, file = 2
< ok 0 0
> exec ^
: string
| set a = 1, file_block_size = 2
< ok 0 0
> exec ^
: string
| set a = 1, file_format = 2
< ok 0 0
> exec ^
: string
| set a = 1, file_name = 2
< ok 0 0
> exec ^
: string
| set a = 1, file_pattern = 2
< ok 0 0
> exec ^
: string
| set a = 1, file_prefix = 2
< ok 0 0
> exec ^
: string
| set a = 1, files = 2
< ok 0 0
> exec ^
: string
| set a = 1, filter = 2
< ok 0 0
> exec ^
: string
| set a = 1, finish = 2
< ok 0 0
> exec ^
: string
| set a = 1, first = 2
< ok 0 0
> exec ^
: string
| set a = 1, first_value = 2
< error
| 1064 42000 …line 1 column 22 near "first_value = 2" 
> exec ^
: string
| set a = 1, fixed = 2
< ok 0 0
> exec ^
: string
| set a = 1, flashback = 2
< ok 0 0
> exec ^
: string
| set a = 1, float = 2
< error
| 1064 42000 …line 1 column 16 near "float = 2" 
> exec ^
: string
| set a = 1, float4 = 2
< error
| 1064 42000 …line 1 column 17 near "float4 = 2" 
> exec ^
: string
| set a = 1, float8 = 2
< error
| 1064 42000 …line 1 column 17 near "float8 = 2" 
> exec ^
: string
| set a = 1, flush = 2
< ok 0 0
> exec ^
: string
| set a = 1, follower_constraints = 2
< ok 0 0
> exec ^
: string
| set a = 1, followers = 2
< ok 0 0
> exec ^
: string
| set a = 1, following = 2
< ok 0 0
> exec ^
: string
| set a = 1, follows = 2
< ok 0 0
> exec ^
: string
| set a = 1, for = 2
< error
| 1064 42000 …line 1 column 14 near "for = 2" 
> exec ^
: string
| set a = 1, force = 2
< error
| 1064 42000 …line 1 column 16 near "force = 2" 
> exec ^
: string
| set a = 1, foreign = 2
< error
| 1064 42000 …line 1 column 18 near "foreign = 2" 
> exec ^
: string
| set a = 1, format = 2
< ok 0 0
> exec ^
: string
| set a = 1, found = 2
< ok 0 0
> exec ^
: string
| set a = 1, from = 2
< error
| 1064 42000 …line 1 column 15 near "from = 2" 
> exec ^
: string
| set a = 1, full = 2
< ok 0 0
> exec ^
: string
| set a = 1, fulltext = 2
< error
| 1064 42000 …line 1 column 19 near "fulltext = 2" 
> exec ^
: string
| set a = 1, function = 2
< ok 0 0
> exec ^
: string
| set a = 1, general = 2
< ok 0 0
> exec ^
: string
| set a = 1, generate = 2
< ok 0 0
> exec ^
: string
| set a = 1, generated = 2
< error
| 1064 42000 …line 1 column 20 near "generated = 2" 
> exec ^
: string
| set a = 1, geomcollection = 2
< ok 0 0
> exec ^
: string
| set a = 1, geometry = 2
< ok 0 0
> exec ^
: string
| set a = 1, geometrycollection = 2
< ok 0 0
> exec ^
: string
| set a = 1, get = 2
< ok 0 0
> exec ^
: string
| set a = 1, get_format = 2
< ok 0 0
> exec ^
: string
| set a = 1, get_source_public_key = 2
< ok 0 0
> exec ^
: string
| set a = 1, global = 2
< ok 0 0
> exec ^
: string
| set a = 1, grant = 2
< error
| 1064 42000 …line 1 column 16 near "grant = 2" 
> exec ^
: string
| set a = 1, grants = 2
< ok 0 0
> exec ^
: string
| set a = 1, greatest = 2
< ok 0 0
> exec ^
: string
| set a = 1, group = 2
< error
| 1064 42000 …line 1 column 16 near "group = 2" 
> exec ^
: string
| set a = 1, group_concat = 2
< ok 0 0
> exec ^
: string
| set a = 1, group_replication = 2
< ok 0 0
> exec ^
: string
| set a = 1, grouping = 2
< ok 0 0
> exec ^
: string
| set a = 1, groups = 2
< error
| 1064 42000 …line 1 column 17 near "groups = 2" 
> exec ^
: string
| set a = 1, gtid_only = 2
< ok 0 0
> exec ^
: string
| set a = 1, gtids = 2
< ok 0 0
> exec ^
: string
| set a = 1, guided = 2
< ok 0 0
> exec ^
: string
| set a = 1, handler = 2
< ok 0 0
> exec ^
: string
| set a = 1, hash = 2
< ok 0 0
> exec ^
: string
| set a = 1, having = 2
< error
| 1064 42000 …line 1 column 17 near "having = 2" 
> exec ^
: string
| set a = 1, header = 2
< ok 0 0
> exec ^
: string
| set a = 1, help = 2
< ok 0 0
> exec ^
: string
| set a = 1, high = 2
< ok 0 0
> exec ^
: string
| set a = 1, high_priority = 2
< error
| 1064 42000 …line 1 column 24 near "high_priority = 2" 
> exec ^
: string
| set a = 1, histogram = 2
< ok 0 0
> exec ^
: string
| set a = 1, histograms_in_flight = 2
< ok 0 0
> exec ^
: string
| set a = 1, history = 2
< ok 0 0
> exec ^
: string
| set a = 1, host = 2
< ok 0 0
> exec ^
: string
| set a = 1, hosts = 2
< ok 0 0
> exec ^
: string
| set a = 1, hour = 2
< ok 0 0
> exec ^
: string
| set a = 1, hour_microsecond = 2
< error
| 1064 42000 …line 1 column 27 near "hour_microsecond = 2" 
> exec ^
: string
| set a = 1, hour_minute = 2
< error
| 1064 42000 …line 1 column 22 near "hour_minute = 2" 
> exec ^
: string
| set a = 1, hour_second = 2
< error
| 1064 42000 …line 1 column 22 near "hour_second = 2" 
> exec ^
: string
| set a = 1, identified = 2
< ok 0 0
> exec ^
: string
| set a = 1, if = 2
< error
| 1064 42000 …line 1 column 13 near "if = 2" 
> exec ^
: string
| set a = 1, ignore = 2
< error
| 1064 42000 …line 1 column 17 near "ignore = 2" 
> exec ^
: string
| set a = 1, ignore_server_ids = 2
< ok 0 0
> exec ^
: string
| set a = 1, ilike = 2
< error
| 1064 42000 …line 1 column 16 near "ilike = 2" 
> exec ^
: string
| set a = 1, import = 2
< ok 0 0
> exec ^
: string
| set a = 1, in = 2
< error
| 1064 42000 …line 1 column 13 near "in = 2" 
> exec ^
: string
| set a = 1, inactive = 2
< ok 0 0
> exec ^
: string
| set a = 1, increment = 2
< ok 0 0
> exec ^
: string
| set a = 1, index = 2
< error
| 1064 42000 …line 1 column 16 near "index = 2" 
> exec ^
: string
| set a = 1, indexes = 2
< ok 0 0
> exec ^
: string
| set a = 1, infile = 2
< error
| 1064 42000 …line 1 column 17 near "infile = 2" 
> exec ^
: string
| set a = 1, initial = 2
< ok 0 0
> exec ^
: string
| set a = 1, initial_size = 2
< ok 0 0
> exec ^
: string
| set a = 1, initiate = 2
< ok 0 0
> exec ^
: string
| set a = 1, inner = 2
< error
| 1064 42000 …line 1 column 16 near "inner = 2" 
> exec ^
: string
| set a = 1, inout = 2
< error
| 1064 42000 …line 1 column 16 near "inout = 2" 
> exec ^
: string
| set a = 1, insensitive = 2
< ok 0 0
> exec ^
: string
| set a = 1, insert = 2
< error
| 1064 42000 …line 1 column 17 near "insert = 2" 
> exec ^
: string
| set a = 1, insert_method = 2
< ok 0 0
> exec ^
: string
| set a = 1, install = 2
< ok 0 0
> exec ^
: string
| set a = 1, instance = 2
< ok 0 0
> exec ^
: string
| set a = 1, int = 2
< error
| 1064 42000 …line 1 column 14 near "int = 2" 
> exec ^
: string
| set a = 1, int1 = 2
< error
| 1064 42000 …line 1 column 15 near "int1 = 2" 
> exec ^
: string
| set a = 1, int2 = 2
< error
| 1064 42000 …line 1 column 15 near "int2 = 2" 
> exec ^
: string
| set a = 1, int3 = 2
< error
| 1064 42000 …line 1 column 15 near "int3 = 2" 
> exec ^
: string
| set a = 1, int4 = 2
< error
| 1064 42000 …line 1 column 15 near "int4 = 2" 
> exec ^
: string
| set a = 1, int8 = 2
< error
| 1064 42000 …line 1 column 15 near "int8 = 2" 
> exec ^
: string
| set a = 1, integer = 2
< error
| 1064 42000 …line 1 column 18 near "integer = 2" 
> exec ^
: string
| set a = 1, intersect = 2
< error
| 1064 42000 …line 1 column 20 near "intersect = 2" 
> exec ^
: string
| set a = 1, interval = 2
< error
| 1064 42000 …line 1 column 19 near "interval = 2" 
> exec ^
: string
| set a = 1, into = 2
< error
| 1064 42000 …line 1 column 15 near "into = 2" 
> exec ^
: string
| set a = 1, invisible = 2
< ok 0 0
> exec ^
: string
| set a = 1, invoker = 2
< ok 0 0
> exec ^
: string
| set a = 1, io = 2
< ok 0 0
> exec ^
: string
| set a = 1, io_after_gtids = 2
< ok 0 0
> exec ^
: string
| set a = 1, io_before_gtids = 2
< ok 0 0
> exec ^
: string
| set a = 1, io_thread = 2
< ok 0 0
> exec ^
: string
| set a = 1, ipc = 2
< ok 0 0
> exec ^
: string
| set a = 1, is = 2
< error
| 1064 42000 …line 1 column 13 near "is = 2" 
> exec ^
: string
| set a = 1, isolation = 2
< ok 0 0
> exec ^
: string
| set a = 1, issuer = 2
< ok 0 0
> exec ^
: string
| set a = 1, iterate = 2
< error
| 1064 42000 …line 1 column 18 near "iterate = 2" 
> exec ^
: string
| set a = 1, job = 2
< ok 0 0
> exec ^
: string
| set a = 1, jobs = 2
< ok 0 0
> exec ^
: string
| set a = 1, join = 2
< error
| 1064 42000 …line 1 column 15 near "join = 2" 
> exec ^
: string
| set a = 1, json = 2
< ok 0 0
> exec ^
: string
| set a = 1, json_arrayagg = 2
< ok 0 0
> exec ^
: string
| set a = 1, json_duality_object = 2
< ok 0 0
> exec ^
: string
| set a = 1, json_objectagg = 2
< ok 0 0
> exec ^
: string
| set a = 1, json_table = 2
< ok 0 0
> exec ^
: string
| set a = 1, json_value = 2
< ok 0 0
> exec ^
: string
| set a = 1, key = 2
< error
| 1064 42000 …line 1 column 14 near "key = 2" 
> exec ^
: string
| set a = 1, key_block_size = 2
< ok 0 0
> exec ^
: string
| set a = 1, keyring = 2
< ok 0 0
> exec ^
: string
| set a = 1, keys = 2
< error
| 1064 42000 …line 1 column 15 near "keys = 2" 
> exec ^
: string
| set a = 1, kill = 2
< error
| 1064 42000 …line 1 column 15 near "kill = 2" 
> exec ^
: string
| set a = 1, labels = 2
< ok 0 0
> exec ^
: string
| set a = 1, lag = 2
< error
| 1064 42000 …line 1 column 14 near "lag = 2" 
> exec ^
: string
| set a = 1, language = 2
< ok 0 0
> exec ^
: string
| set a = 1, last = 2
< ok 0 0
> exec ^
: string
| set a = 1, last_value = 2
< error
| 1064 42000 …line 1 column 21 near "last_value = 2" 
> exec ^
: string
| set a = 1, lateral = 2
< ok 0 0
> exec ^
: string
| set a = 1, lead = 2
< error
| 1064 42000 …line 1 column 15 near "lead = 2" 
> exec ^
: string
| set a = 1, leader_constraints = 2
< ok 0 0
> exec ^
: string
| set a = 1, leading = 2
< error
| 1064 42000 …line 1 column 18 near "leading = 2" 
> exec ^
: string
| set a = 1, learner_constraints = 2
< ok 0 0
> exec ^
: string
| set a = 1, learners = 2
< ok 0 0
> exec ^
: string
| set a = 1, least = 2
< ok 0 0
> exec ^
: string
| set a = 1, leave = 2
< error
| 1064 42000 …line 1 column 16 near "leave = 2" 
> exec ^
: string
| set a = 1, leaves = 2
< ok 0 0
> exec ^
: string
| set a = 1, left = 2
< error
| 1064 42000 …line 1 column 15 near "left = 2" 
> exec ^
: string
| set a = 1, less = 2
< ok 0 0
> exec ^
: string
| set a = 1, level = 2
< ok 0 0
> exec ^
: string
| set a = 1, library = 2
< ok 0 0
> exec ^
: string
| set a = 1, like = 2
< error
| 1064 42000 …line 1 column 15 near "like = 2" 
> exec ^
: string
| set a = 1, limit = 2
< error
| 1064 42000 …line 1 column 16 near "limit = 2" 
> exec ^
: string
| set a = 1, linear = 2
< error
| 1064 42000 …line 1 column 17 near "linear = 2" 
> exec ^
: string
| set a = 1, lines = 2
< error
| 1064 42000 …line 1 column 16 near "lines = 2" 
> exec ^
: string
| set a = 1, linestring = 2
< ok 0 0
> exec ^
: string
| set a = 1, list = 2
< ok 0 0
> exec ^
: string
| set a = 1, load = 2
< error
| 1064 42000 …line 1 column 15 near "load = 2" 
> exec ^
: string
| set a = 1, local = 2
< ok 0 0
> exec ^
: string
| set a = 1, localtime = 2
< error
| 1064 42000 …line 1 column 20 near "localtime = 2" 
> exec ^
: string
| set a = 1, localtimestamp = 2
< error
| 1064 42000 …line 1 column 25 near "localtimestamp = 2" 
> exec ^
: string
| set a = 1, location = 2
< ok 0 0
> exec ^
: string
| set a = 1, lock = 2
< error
| 1064 42000 …line 1 column 15 near "lock = 2" 
> exec ^
: string
| set a = 1, locked = 2
< ok 0 0
> exec ^
: string
| set a = 1, locks = 2
< ok 0 0
> exec ^
: string
| set a = 1, log = 2
< ok 0 0
> exec ^
: string
| set a = 1, logfile = 2
< ok 0 0
> exec ^
: string
| set a = 1, logs = 2
< ok 0 0
> exec ^
: string
| set a = 1, long = 2
< error
| 1064 42000 …line 1 column 15 near "long = 2" 
> exec ^
: string
| set a = 1, longblob = 2
< error
| 1064 42000 …line 1 column 19 near "longblob = 2" 
> exec ^
: string
| set a = 1, longtext = 2
< error
| 1064 42000 …line 1 column 19 near "longtext = 2" 
> exec ^
: string
| set a = 1, loop = 2
< ok 0 0
> exec ^
: string
| set a = 1, low = 2
< ok 0 0
> exec ^
: string
| set a = 1, low_priority = 2
< error
| 1064 42000 …line 1 column 23 near "low_priority = 2" 
> exec ^
: string
| set a = 1, manual = 2
< ok 0 0
> exec ^
: string
| set a = 1, master = 2
< ok 0 0
> exec ^
: string
| set a = 1, master_log_file = 2
< ok 0 0
> exec ^
: string
| set a = 1, match = 2
< error
| 1064 42000 …line 1 column 16 near "match = 2" 
> exec ^
: string
| set a = 1, materialized = 2
< ok 0 0
> exec ^
: string
| set a = 1, max = 2
< ok 0 0
> exec ^
: string
| set a = 1, max_connections_per_hour = 2
< ok 0 0
> exec ^
: string
| set a = 1, max_queries_per_hour = 2
< ok 0 0
> exec ^
: string
| set a = 1, max_rows = 2
< ok 0 0
> exec ^
: string
| set a = 1, max_size = 2
< ok 0 0
> exec ^
: string
| set a = 1, max_updates_per_hour = 2
< ok 0 0
> exec ^
: string
| set a = 1, max_user_connections = 2
< ok 0 0
> exec ^
: string
| set a = 1, maxvalue = 2
< error
| 1064 42000 …line 1 column 19 near "maxvalue = 2" 
> exec ^
: string
| set a = 1, medium = 2
< ok 0 0
> exec ^
: string
| set a = 1, mediumblob = 2
< error
| 1064 42000 …line 1 column 21 near "mediumblob = 2" 
> exec ^
: string
| set a = 1, mediumint = 2
< error
| 1064 42000 …line 1 column 20 near "mediumint = 2" 
> exec ^
: string
| set a = 1, mediumtext = 2
< error
| 1064 42000 …line 1 column 21 near "mediumtext = 2" 
> exec ^
: string
| set a = 1, member = 2
< ok 0 0
> exec ^
: string
| set a = 1, memory = 2
< ok 0 0
> exec ^
: string
| set a = 1, merge = 2
< ok 0 0
> exec ^
: string
| set a = 1, message_text = 2
< ok 0 0
> exec ^
: string
| set a = 1, microsecond = 2
< ok 0 0
> exec ^
: string
| set a = 1, mid = 2
< ok 0 0
> exec ^
: string
| set a = 1, middleint = 2
< error
| 1064 42000 …line 1 column 20 near "middleint = 2" 
> exec ^
: string
| set a = 1, migrate = 2
< ok 0 0
> exec ^
: string
| set a = 1, min = 2
< ok 0 0
> exec ^
: string
| set a = 1, min_rows = 2
< ok 0 0
> exec ^
: string
| set a = 1, minute = 2
< ok 0 0
> exec ^
: string
| set a = 1, minute_microsecond = 2
< error
| 1064 42000 …line 1 column 29 near "minute_microsecond = 2" 
> exec ^
: string
| set a = 1, minute_second = 2
< error
| 1064 42000 …line 1 column 24 near "minute_second = 2" 
> exec ^
: string
| set a = 1, minvalue = 2
< ok 0 0
> exec ^
: string
| set a = 1, mod = 2
< error
| 1064 42000 …line 1 column 14 near "mod = 2" 
> exec ^
: string
| set a = 1, mode = 2
< ok 0 0
> exec ^
: string
| set a = 1, modifies = 2
< ok 0 0
> exec ^
: string
| set a = 1, modify = 2
< ok 0 0
> exec ^
: string
| set a = 1, month = 2
< ok 0 0
> exec ^
: string
| set a = 1, multilinestring = 2
< ok 0 0
> exec ^
: string
| set a = 1, multipoint = 2
< ok 0 0
> exec ^
: string
| set a = 1, multipolygon = 2
< ok 0 0
> exec ^
: string
| set a = 1, mutex = 2
< ok 0 0
> exec ^
: string
| set a = 1, mysql_errno = 2
< ok 0 0
> exec ^
: string
| set a = 1, name = 2
< ok 0 0
> exec ^
: string
| set a = 1, names = 2
< ok 0 0
> exec ^
: string
| set a = 1, national = 2
< ok 0 0
> exec ^
: string
| set a = 1, natural = 2
< error
| 1064 42000 …line 1 column 18 near "natural = 2" 
> exec ^
: string
| set a = 1, nchar = 2
< ok 0 0
> exec ^
: string
| set a = 1, ndb = 2
< ok 0 0
> exec ^
: string
| set a = 1, ndbcluster = 2
< ok 0 0
> exec ^
: string
| set a = 1, nested = 2
< ok 0 0
> exec ^
: string
| set a = 1, network_namespace = 2
< ok 0 0
> exec ^
: string
| set a = 1, never = 2
< ok 0 0
> exec ^
: string
| set a = 1, new = 2
< ok 0 0
> exec ^
: string
| set a = 1, next = 2
< ok 0 0
> exec ^
: string
| set a = 1, next_row_id = 2
< ok 0 0
> exec ^
: string
| set a = 1, no = 2
< ok 0 0
> exec ^
: string
| set a = 1, no_wait = 2
< ok 0 0
> exec ^
: string
| set a = 1, no_write_to_binlog = 2
< error
| 1064 42000 …line 1 column 29 near "no_write_to_binlog = 2" 
> exec ^
: string
| set a = 1, nocache = 2
< ok 0 0
> exec ^
: string
| set a = 1, nocycle = 2
< ok 0 0
> exec ^
: string
| set a = 1, nodegroup = 2
< ok 0 0
> exec ^
: string
| set a = 1, nomaxvalue = 2
< ok 0 0
> exec ^
: string
| set a = 1, nominvalue = 2
< ok 0 0
> exec ^
: string
| set a = 1, nonclustered = 2
< ok 0 0
> exec ^
: string
| set a = 1, none = 2
< ok 0 0
> exec ^
: string
| set a = 1, not = 2
< error
| 1064 42000 …line 1 column 14 near "not = 2" 
> exec ^
: string
| set a = 1, now = 2
< ok 0 0
> exec ^
: string
| set a = 1, nowait = 2
< ok 0 0
> exec ^
: string
| set a = 1, nth_value = 2
< error
| 1064 42000 …line 1 column 20 near "nth_value = 2" 
> exec ^
: string
| set a = 1, ntile = 2
< error
| 1064 42000 …line 1 column 16 near "ntile = 2" 
> exec ^
: string
| set a = 1, null = 2
< error
| 1064 42000 …line 1 column 15 near "null = 2" 
> exec ^
: string
| set a = 1, nullif = 2
< ok 0 0
> exec ^
: string
| set a = 1, nulls = 2
< ok 0 0
> exec ^
: string
| set a = 1, number = 2
< ok 0 0
> exec ^
: string
| set a = 1, numeric = 2
< error
| 1064 42000 …line 1 column 18 near "numeric = 2" 
> exec ^
: string
| set a = 1, nvarchar = 2
< ok 0 0
> exec ^
: string
| set a = 1, of = 2
< error
| 1064 42000 …line 1 column 13 near "of = 2" 
> exec ^
: string
| set a = 1, off = 2
< ok 0 0
> exec ^
: string
| set a = 1, offset = 2
< ok 0 0
> exec ^
: string
| set a = 1, oj = 2
< ok 0 0
> exec ^
: string
| set a = 1, old = 2
< ok 0 0
> exec ^
: string
| set a = 1, oltp_read_only = 2
< ok 0 0
> exec ^
: string
| set a = 1, oltp_read_write = 2
< ok 0 0
> exec ^
: string
| set a = 1, oltp_write_only = 2
< ok 0 0
> exec ^
: string
| set a = 1, on = 2
< error
| 1064 42000 …line 1 column 13 near "on = 2" 
> exec ^
: string
| set a = 1, one = 2
< ok 0 0
> exec ^
: string
| set a = 1, only = 2
< ok 0 0
> exec ^
: string
| set a = 1, open = 2
< ok 0 0
> exec ^
: string
| set a = 1, optimize = 2
< error
| 1064 42000 …line 1 column 19 near "optimize = 2" 
> exec ^
: string
| set a = 1, optimizer_costs = 2
< ok 0 0
> exec ^
: string
| set a = 1, option = 2
< error
| 1064 42000 …line 1 column 17 near "option = 2" 
> exec ^
: string
| set a = 1, optional = 2
< ok 0 0
> exec ^
: string
| set a = 1, optionally = 2
< error
| 1064 42000 …line 1 column 21 near "optionally = 2" 
> exec ^
: string
| set a = 1, options = 2
< ok 0 0
> exec ^
: string
| set a = 1, or = 2
< error
| 1064 42000 …line 1 column 13 near "or = 2" 
> exec ^
: string
| set a = 1, order = 2
< error
| 1064 42000 …line 1 column 16 near "order = 2" 
> exec ^
: string
| set a = 1, ordinality = 2
< ok 0 0
> exec ^
: string
| set a = 1, organization = 2
< ok 0 0
> exec ^
: string
| set a = 1, others = 2
< ok 0 0
> exec ^
: string
| set a = 1, out = 2
< error
| 1064 42000 …line 1 column 14 near "out = 2" 
> exec ^
: string
| set a = 1, outer = 2
< error
| 1064 42000 …line 1 column 16 near "outer = 2" 
> exec ^
: string
| set a = 1, outfile = 2
< error
| 1064 42000 …line 1 column 18 near "outfile = 2" 
> exec ^
: string
| set a = 1, over = 2
< error
| 1064 42000 …line 1 column 15 near "over = 2" 
> exec ^
: string
| set a = 1, overlay = 2
< ok 0 0
> exec ^
: string
| set a = 1, owner = 2
< ok 0 0
> exec ^
: string
| set a = 1, pack_keys = 2
< ok 0 0
> exec ^
: string
| set a = 1, page = 2
< ok 0 0
> exec ^
: string
| set a = 1, parallel = 2
< ok 0 0
> exec ^
: string
| set a = 1, parameters = 2
< ok 0 0
> exec ^
: string
| set a = 1, parse_tree = 2
< ok 0 0
> exec ^
: string
| set a = 1, parser = 2
< ok 0 0
> exec ^
: string
| set a = 1, partial = 2
< ok 0 0
> exec ^
: string
| set a = 1, partition = 2
< error
| 1064 42000 …line 1 column 20 near "partition = 2" 
> exec ^
: string
| set a = 1, partitioning = 2
< ok 0 0
> exec ^
: string
| set a = 1, partitions = 2
< ok 0 0
> exec ^
: string
| set a = 1, password = 2
< ok 0 0
> exec ^
: string
| set a = 1, password_lock_time = 2
< ok 0 0
> exec ^
: string
| set a = 1, path = 2
< ok 0 0
> exec ^
: string
| set a = 1, percent_rank = 2
< error
| 1064 42000 …line 1 column 23 near "percent_rank = 2" 
> exec ^
: string
| set a = 1, persist = 2
< ok 0 0
> exec ^
: string
| set a = 1, persist_only = 2
< ok 0 0
> exec ^
: string
| set a = 1, phase = 2
< ok 0 0
> exec ^
: string
| set a = 1, placement = 2
< ok 0 0
> exec ^
: string
| set a = 1, placing = 2
< ok 0 0
> exec ^
: string
| set a = 1, plan = 2
< ok 0 0
> exec ^
: string
| set a = 1, plugin = 2
< ok 0 0
> exec ^
: string
| set a = 1, plugin_dir = 2
< ok 0 0
> exec ^
: string
| set a = 1, plugins = 2
< ok 0 0
> exec ^
: string
| set a = 1, point = 2
< ok 0 0
> exec ^
: string
| set a = 1, policy = 2
< ok 0 0
> exec ^
: string
| set a = 1, polygon = 2
< ok 0 0
> exec ^
: string
| set a = 1, port = 2
< ok 0 0
> exec ^
: string
| set a = 1, position = 2
< ok 0 0
> exec ^
: string
| set a = 1, pre_split_regions = 2
< ok 0 0
> exec ^
: string
| set a = 1, precedes = 2
< ok 0 0
> exec ^
: string
| set a = 1, preceding = 2
< ok 0 0
> exec ^
: string
| set a = 1, precision = 2
< error
| 1064 42000 …line 1 column 20 near "precision = 2" 
> exec ^
: string
| set a = 1, prepare = 2
< ok 0 0
> exec ^
: string
| set a = 1, preserve = 2
< ok 0 0
> exec ^
: string
| set a = 1, prev = 2
< ok 0 0
> exec ^
: string
| set a = 1, primary = 2
< error
| 1064 42000 …line 1 column 18 near "primary = 2" 
> exec ^
: string
| set a = 1, primary_region = 2
< ok 0 0
> exec ^
: string
| set a = 1, priority = 2
< ok 0 0
> exec ^
: string
| set a = 1, privilege_checks_user = 2
< ok 0 0
> exec ^
: string
| set a = 1, privileges = 2
< ok 0 0
> exec ^
: string
| set a = 1, procedure = 2
< error
| 1064 42000 …line 1 column 20 near "procedure = 2" 
> exec ^
: string
| set a = 1, process = 2
< ok 0 0
> exec ^
: string
| set a = 1, processed_keys = 2
< ok 0 0
> exec ^
: string
| set a = 1, processlist = 2
< ok 0 0
> exec ^
: string
| set a = 1, profile = 2
< ok 0 0
> exec ^
: string
| set a = 1, profiles = 2
< ok 0 0
> exec ^
: string
| set a = 1, proxy = 2
< ok 0 0
> exec ^
: string
| set a = 1, purge = 2
< ok 0 0
> exec ^
: string
| set a = 1, qualify = 2
< ok 0 0
> exec ^
: string
| set a = 1, quarter = 2
< ok 0 0
> exec ^
: string
| set a = 1, query = 2
< ok 0 0
> exec ^
: string
| set a = 1, query_limit = 2
< ok 0 0
> exec ^
: string
| set a = 1, quick = 2
< ok 0 0
> exec ^
: string
| set a = 1, random = 2
< ok 0 0
> exec ^
: string
| set a = 1, range = 2
< error
| 1064 42000 …line 1 column 16 near "range = 2" 
> exec ^
: string
| set a = 1, rank = 2
< error
| 1064 42000 …line 1 column 15 near "rank = 2" 
> exec ^
: string
| set a = 1, read = 2
< error
| 1064 42000 …line 1 column 15 near "read = 2" 
> exec ^
: string
| set a = 1, read_only = 2
< ok 0 0
> exec ^
: string
| set a = 1, read_write = 2
< ok 0 0
> exec ^
: string
| set a = 1, reads = 2
< ok 0 0
> exec ^
: string
| set a = 1, real = 2
< error
| 1064 42000 …line 1 column 15 near "real = 2" 
> exec ^
: string
| set a = 1, rebuild = 2
< ok 0 0
> exec ^
: string
| set a = 1, recover = 2
< ok 0 0
> exec ^
: string
| set a = 1, recursive = 2
< error
| 1064 42000 …line 1 column 20 near "recursive = 2" 
> exec ^
: string
| set a = 1, redo_buffer_size = 2
< ok 0 0
> exec ^
: string
| set a = 1, redundant = 2
< ok 0 0
> exec ^
: string
| set a = 1, reference = 2
< ok 0 0
> exec ^
: string
| set a = 1, references = 2
< error
| 1064 42000 …line 1 column 21 near "references = 2" 
> exec ^
: string
| set a = 1, regexp = 2
< error
| 1064 42000 …line 1 column 17 near "regexp = 2" 
> exec ^
: string
| set a = 1, region = 2
< ok 0 0
> exec ^
: string
| set a = 1, regions = 2
< ok 0 0
> exec ^
: string
| set a = 1, registration = 2
< ok 0 0
> exec ^
: string
| set a = 1, relational = 2
< ok 0 0
> exec ^
: string
| set a = 1, relay = 2
< ok 0 0
> exec ^
: string
| set a = 1, relay_log_file = 2
< ok 0 0
> exec ^
: string
| set a = 1, relay_log_pos = 2
< ok 0 0
> exec ^
: string
| set a = 1, relay_thread = 2
< ok 0 0
> exec ^
: string
| set a = 1, relaylog = 2
< ok 0 0
> exec ^
: string
| set a = 1, release = 2
< error
| 1064 42000 …line 1 column 18 near "release = 2" 
> exec ^
: string
| set a = 1, reload = 2
< ok 0 0
> exec ^
: string
| set a = 1, remove = 2
< ok 0 0
> exec ^
: string
| set a = 1, rename = 2
< error
| 1064 42000 …line 1 column 17 near "rename = 2" 
> exec ^
: string
| set a = 1, reorganize = 2
< ok 0 0
> exec ^
: string
| set a = 1, repair = 2
< ok 0 0
> exec ^
: string
| set a = 1, repeat = 2
< error
| 1064 42000 …line 1 column 17 near "repeat = 2" 
> exec ^
: string
| set a = 1, repeatable = 2
< ok 0 0
> exec ^
: string
| set a = 1, replace = 2
< error
| 1064 42000 …line 1 column 18 near "replace = 2" 
> exec ^
: string
| set a = 1, replica = 2
< ok 0 0
> exec ^
: string
| set a = 1, replicas = 2
< ok 0 0
> exec ^
: string
| set a = 1, replicate_do_db = 2
< ok 0 0
> exec ^
: string
| set a = 1, replicate_do_table = 2
< ok 0 0
> exec ^
: string
| set a = 1, replicate_ignore_db = 2
< ok 0 0
> exec ^
: string
| set a = 1, replicate_ignore_table = 2
< ok 0 0
> exec ^
: string
| set a = 1, replicate_rewrite_db = 2
< ok 0 0
> exec ^
: string
| set a = 1, replicate_wild_do_table = 2
< ok 0 0
> exec ^
: string
| set a = 1, replicate_wild_ignore_table = 2
< ok 0 0
> exec ^
: string
| set a = 1, replication = 2
< ok 0 0
> exec ^
: string
| set a = 1, require = 2
< error
| 1064 42000 …line 1 column 18 near "require = 2" 
> exec ^
: string
| set a = 1, require_row_format = 2
< ok 0 0
> exec ^
: string
| set a = 1, require_table_primary_key_check = 2
< ok 0 0
> exec ^
: string
| set a = 1, reset = 2
< ok 0 0
> exec ^
: string
| set a = 1, resignal = 2
< ok 0 0
> exec ^
: string
| set a = 1, resource = 2
< ok 0 0
> exec ^
: string
| set a = 1, respect = 2
< ok 0 0
> exec ^
: string
| set a = 1, restart = 2
< ok 0 0
> exec ^
: string
| set a = 1, restore = 2
< ok 0 0
> exec ^
: string
| set a = 1, restores = 2
< ok 0 0
> exec ^
: string
| set a = 1, restrict = 2
< error
| 1064 42000 …line 1 column 19 near "restrict = 2" 
> exec ^
: string
| set a = 1, resume = 2
< ok 0 0
> exec ^
: string
| set a = 1, retain = 2
< ok 0 0
> exec ^
: string
| set a = 1, return = 2
< ok 0 0
> exec ^
: string
| set a = 1, returned_sqlstate = 2
< ok 0 0
> exec ^
: string
| set a = 1, returning = 2
< ok 0 0
> exec ^
: string
| set a = 1, returns = 2
< ok 0 0
> exec ^
: string
| set a = 1, reuse = 2
< ok 0 0
> exec ^
: string
| set a = 1, reverse = 2
< ok 0 0
> exec ^
: string
| set a = 1, revoke = 2
< error
| 1064 42000 …line 1 column 17 near "revoke = 2" 
> exec ^
: string
| set a = 1, right = 2
< error
| 1064 42000 …line 1 column 16 near "right = 2" 
> exec ^
: string
| set a = 1, rlike = 2
< error
| 1064 42000 …line 1 column 16 near "rlike = 2" 
> exec ^
: string
| set a = 1, role = 2
< ok 0 0
> exec ^
: string
| set a = 1, rollback = 2
< ok 0 0
> exec ^
: string
| set a = 1, rollup = 2
< ok 0 0
> exec ^
: string
| set a = 1, rotate = 2
< ok 0 0
> exec ^
: string
| set a = 1, routine = 2
< ok 0 0
> exec ^
: string
| set a = 1, row = 2
< error
| 1064 42000 …line 1 column 14 near "row = 2" 
> exec ^
: string
| set a = 1, row_count = 2
< ok 0 0
> exec ^
: string
| set a = 1, row_format = 2
< ok 0 0
> exec ^
: string
| set a = 1, row_number = 2
< error
| 1064 42000 …line 1 column 21 near "row_number = 2" 
> exec ^
: string
| set a = 1, rows = 2
< error
| 1064 42000 …line 1 column 15 near "rows = 2" 
> exec ^
: string
| set a = 1, rtree = 2
< ok 0 0
> exec ^
: string
| set a = 1, ru = 2
< ok 0 0
> exec ^
: string
| set a = 1, ru_per_sec = 2
< ok 0 0
> exec ^
: string
| set a = 1, rule = 2
< ok 0 0
> exec ^
: string
| set a = 1, run = 2
< ok 0 0
> exec ^
: string
| set a = 1, s3 = 2
< ok 0 0
> exec ^
: string
| set a = 1, savepoint = 2
< ok 0 0
> exec ^
: string
| set a = 1, schedule = 2
< ok 0 0
> exec ^
: string
| set a = 1, schema = 2
< error
| 1064 42000 …line 1 column 17 near "schema = 2" 
> exec ^
: string
| set a = 1, schema_name = 2
< ok 0 0
> exec ^
: string
| set a = 1, schemas = 2
< error
| 1064 42000 …line 1 column 18 near "schemas = 2" 
> exec ^
: string
| set a = 1, second = 2
< ok 0 0
> exec ^
: string
| set a = 1, second_microsecond = 2
< error
| 1064 42000 …line 1 column 29 near "second_microsecond = 2" 
> exec ^
: string
| set a = 1, secondary = 2
< ok 0 0
> exec ^
: string
| set a = 1, secondary_engine = 2
< ok 0 0
> exec ^
: string
| set a = 1, secondary_engine_attribute = 2
< ok 0 0
> exec ^
: string
| set a = 1, secondary_load = 2
< ok 0 0
> exec ^
: string
| set a = 1, secondary_unload = 2
< ok 0 0
> exec ^
: string
| set a = 1, security = 2
< ok 0 0
> exec ^
: string
| set a = 1, select = 2
< error
| 1064 42000 …line 1 column 17 near "select = 2" 
> exec ^
: string
| set a = 1, sensitive = 2
< ok 0 0
> exec ^
: string
| set a = 1, separator = 2
< ok 0 0
> exec ^
: string
| set a = 1, sequence = 2
< ok 0 0
> exec ^
: string
| set a = 1, serial = 2
< ok 0 0
> exec ^
: string
| set a = 1, serializable = 2
< ok 0 0
> exec ^
: string
| set a = 1, server = 2
< ok 0 0
> exec ^
: string
| set a = 1, session = 2
< ok 0 0
> exec ^
: string
| set a = 1, session_states = 2
< ok 0 0
> exec ^
: string
| set a = 1, session_user = 2
< ok 0 0
> exec ^
: string
| set a = 1, set = 2
< error
| 1064 42000 …line 1 column 14 near "set = 2" 
> exec ^
: string
| set a = 1, sets = 2
< ok 0 0
> exec ^
: string
| set a = 1, shard_row_id_bits = 2
< ok 0 0
> exec ^
: string
| set a = 1, share = 2
< ok 0 0
> exec ^
: string
| set a = 1, show = 2
< error
| 1064 42000 …line 1 column 15 near "show = 2" 
> exec ^
: string
| set a = 1, shutdown = 2
< ok 0 0
> exec ^
: string
| set a = 1, signal = 2
< ok 0 0
> exec ^
: string
| set a = 1, signed = 2
< ok 0 0
> exec ^
: string
| set a = 1, similar = 2
< ok 0 0
> exec ^
: string
| set a = 1, simple = 2
< ok 0 0
> exec ^
: string
| set a = 1, skip = 2
< ok 0 0
> exec ^
: string
| set a = 1, slave = 2
< ok 0 0
> exec ^
: string
| set a = 1, slow = 2
< ok 0 0
> exec ^
: string
| set a = 1, smallint = 2
< error
| 1064 42000 …line 1 column 19 near "smallint = 2" 
> exec ^
: string
| set a = 1, snapshot = 2
< ok 0 0
> exec ^
: string
| set a = 1, socket = 2
< ok 0 0
> exec ^
: string
| set a = 1, some = 2
< ok 0 0
> exec ^
: string
| set a = 1, soname = 2
< ok 0 0
> exec ^
: string
| set a = 1, sounds = 2
< ok 0 0
> exec ^
: string
| set a = 1, source = 2
< ok 0 0
> exec ^
: string
| set a = 1, source_auto_position = 2
< ok 0 0
> exec ^
: string
| set a = 1, source_bind = 2
< ok 0 0
> exec ^
: string
| set a = 1, source_compression_algorithms = 2
< ok 0 0
> exec ^
: string
| set a = 1, source_connect_retry = 2
< ok 0 0
> exec ^
: string
| set a = 1, source_connection_auto_failover = 2
< ok 0 0
> exec ^
: string
| set a = 1, source_delay = 2
< ok 0 0
> exec ^
: string
| set a = 1, source_heartbeat_period = 2
< ok 0 0
> exec ^
: string
| set a = 1, source_host = 2
< ok 0 0
> exec ^
: string
| set a = 1, source_log_file = 2
< ok 0 0
> exec ^
: string
| set a = 1, source_log_pos = 2
< ok 0 0
> exec ^
: string
| set a = 1, source_password = 2
< ok 0 0
> exec ^
: string
| set a = 1, source_port = 2
< ok 0 0
> exec ^
: string
| set a = 1, source_public_key_path = 2
< ok 0 0
> exec ^
: string
| set a = 1, source_retry_count = 2
< ok 0 0
> exec ^
: string
| set a = 1, source_ssl = 2
< ok 0 0
> exec ^
: string
| set a = 1, source_ssl_ca = 2
< ok 0 0
> exec ^
: string
| set a = 1, source_ssl_capath = 2
< ok 0 0
> exec ^
: string
| set a = 1, source_ssl_cert = 2
< ok 0 0
> exec ^
: string
| set a = 1, source_ssl_cipher = 2
< ok 0 0
> exec ^
: string
| set a = 1, source_ssl_crl = 2
< ok 0 0
> exec ^
: string
| set a = 1, source_ssl_crlpath = 2
< ok 0 0
> exec ^
: string
| set a = 1, source_ssl_key = 2
< ok 0 0
> exec ^
: string
| set a = 1, source_ssl_verify_server_cert = 2
< ok 0 0
> exec ^
: string
| set a = 1, source_tls_ciphersuites = 2
< ok 0 0
> exec ^
: string
| set a = 1, source_tls_version = 2
< ok 0 0
> exec ^
: string
| set a = 1, source_user = 2
< ok 0 0
> exec ^
: string
| set a = 1, source_zstd_compression_level = 2
< ok 0 0
> exec ^
: string
| set a = 1, spatial = 2
< error
| 1064 42000 …line 1 column 18 near "spatial = 2" 
> exec ^
: string
| set a = 1, specific = 2
< ok 0 0
> exec ^
: string
| set a = 1, split = 2
< ok 0 0
> exec ^
: string
| set a = 1, sql = 2
< error
| 1064 42000 …line 1 column 14 near "sql = 2" 
> exec ^
: string
| set a = 1, sql_after_gtids = 2
< ok 0 0
> exec ^
: string
| set a = 1, sql_after_mts_gaps = 2
< ok 0 0
> exec ^
: string
| set a = 1, sql_before_gtids = 2
< ok 0 0
> exec ^
: string
| set a = 1, sql_big_result = 2
< error
| 1064 42000 …line 1 column 25 near "sql_big_result = 2" 
> exec ^
: string
| set a = 1, sql_buffer_result = 2
< ok 0 0
> exec ^
: string
| set a = 1, sql_calc_found_rows = 2
< error
| 1064 42000 …line 1 column 30 near "sql_calc_found_rows = 2" 
> exec ^
: string
| set a = 1, sql_no_cache = 2
< ok 0 0
> exec ^
: string
| set a = 1, sql_small_result = 2
< error
| 1064 42000 …line 1 column 27 near "sql_small_result = 2" 
> exec ^
: string
| set a = 1, sql_thread = 2
< ok 0 0
> exec ^
: string
| set a = 1, sql_tsi_day = 2
< ok 0 0
> exec ^
: string
| set a = 1, sql_tsi_hour = 2
< ok 0 0
> exec ^
: string
| set a = 1, sql_tsi_minute = 2
< ok 0 0
> exec ^
: string
| set a = 1, sql_tsi_month = 2
< ok 0 0
> exec ^
: string
| set a = 1, sql_tsi_quarter = 2
< ok 0 0
> exec ^
: string
| set a = 1, sql_tsi_second = 2
< ok 0 0
> exec ^
: string
| set a = 1, sql_tsi_week = 2
< ok 0 0
> exec ^
: string
| set a = 1, sql_tsi_year = 2
< ok 0 0
> exec ^
: string
| set a = 1, sqlexception = 2
< error
| 1064 42000 …line 1 column 23 near "sqlexception = 2" 
> exec ^
: string
| set a = 1, sqlstate = 2
< error
| 1064 42000 …line 1 column 19 near "sqlstate = 2" 
> exec ^
: string
| set a = 1, sqlwarning = 2
< error
| 1064 42000 …line 1 column 21 near "sqlwarning = 2" 
> exec ^
: string
| set a = 1, srid = 2
< ok 0 0
> exec ^
: string
| set a = 1, ssl = 2
< error
| 1064 42000 …line 1 column 14 near "ssl = 2" 
> exec ^
: string
| set a = 1, st_collect = 2
< ok 0 0
> exec ^
: string
| set a = 1, stacked = 2
< ok 0 0
> exec ^
: string
| set a = 1, start = 2
< ok 0 0
> exec ^
: string
| set a = 1, start_time = 2
< ok 0 0
> exec ^
: string
| set a = 1, starting = 2
< error
| 1064 42000 …line 1 column 19 near "starting = 2" 
> exec ^
: string
| set a = 1, starts = 2
< ok 0 0
> exec ^
: string
| set a = 1, stats_auto_recalc = 2
< ok 0 0
> exec ^
: string
| set a = 1, stats_buckets = 2
< ok 0 0
> exec ^
: string
| set a = 1, stats_col_choice = 2
< ok 0 0
> exec ^
: string
| set a = 1, stats_col_list = 2
< ok 0 0
> exec ^
: string
| set a = 1, stats_extended = 2
< ok 0 0
> exec ^
: string
| set a = 1, stats_healthy = 2
< ok 0 0
> exec ^
: string
| set a = 1, stats_histograms = 2
< ok 0 0
> exec ^
: string
| set a = 1, stats_locked = 2
< ok 0 0
> exec ^
: string
| set a = 1, stats_meta = 2
< ok 0 0
> exec ^
: string
| set a = 1, stats_persistent = 2
< ok 0 0
> exec ^
: string
| set a = 1, stats_sample_pages = 2
< ok 0 0
> exec ^
: string
| set a = 1, stats_sample_rate = 2
< ok 0 0
> exec ^
: string
| set a = 1, stats_topn = 2
< ok 0 0
> exec ^
: string
| set a = 1, status = 2
< ok 0 0
> exec ^
: string
| set a = 1, std = 2
< ok 0 0
> exec ^
: string
| set a = 1, stddev = 2
< ok 0 0
> exec ^
: string
| set a = 1, stddev_pop = 2
< ok 0 0
> exec ^
: string
| set a = 1, stddev_samp = 2
< ok 0 0
> exec ^
: string
| set a = 1, stop = 2
< ok 0 0
> exec ^
: string
| set a = 1, storage = 2
< ok 0 0
> exec ^
: string
| set a = 1, stored = 2
< error
| 1064 42000 …line 1 column 17 near "stored = 2" 
> exec ^
: string
| set a = 1, straight_join = 2
< error
| 1064 42000 …line 1 column 24 near "straight_join = 2" 
> exec ^
: string
| set a = 1, stream = 2
< ok 0 0
> exec ^
: string
| set a = 1, strict_load = 2
< ok 0 0
> exec ^
: string
| set a = 1, string = 2
< ok 0 0
> exec ^
: string
| set a = 1, subclass_origin = 2
< ok 0 0
> exec ^
: string
| set a = 1, subdate = 2
< ok 0 0
> exec ^
: string
| set a = 1, subject = 2
< ok 0 0
> exec ^
: string
| set a = 1, subpartition = 2
< ok 0 0
> exec ^
: string
| set a = 1, subpartitions = 2
< ok 0 0
> exec ^
: string
| set a = 1, substr = 2
< ok 0 0
> exec ^
: string
| set a = 1, substring = 2
< ok 0 0
> exec ^
: string
| set a = 1, sum = 2
< ok 0 0
> exec ^
: string
| set a = 1, super = 2
< ok 0 0
> exec ^
: string
| set a = 1, survival_preferences = 2
< ok 0 0
> exec ^
: string
| set a = 1, suspend = 2
< ok 0 0
> exec ^
: string
| set a = 1, swaps = 2
< ok 0 0
> exec ^
: string
| set a = 1, switch_group = 2
< ok 0 0
> exec ^
: string
| set a = 1, switches = 2
< ok 0 0
> exec ^
: string
| set a = 1, sysdate = 2
< ok 0 0
> exec ^
: string
| set a = 1, system = 2
< ok 0 0
> exec ^
: string
| set a = 1, system_user = 2
< ok 0 0
> exec ^
: string
| set a = 1, table = 2
< error
| 1064 42000 …line 1 column 16 near "table = 2" 
> exec ^
: string
| set a = 1, table_checksum = 2
< ok 0 0
> exec ^
: string
| set a = 1, table_name = 2
< ok 0 0
> exec ^
: string
| set a = 1, tables = 2
< ok 0 0
> exec ^
: string
| set a = 1, tablesample = 2
< error
| 1064 42000 …line 1 column 22 near "tablesample = 2" 
> exec ^
: string
| set a = 1, tablespace = 2
< ok 0 0
> exec ^
: string
| set a = 1, task_types = 2
< ok 0 0
> exec ^
: string
| set a = 1, temporary = 2
< ok 0 0
> exec ^
: string
| set a = 1, temptable = 2
< ok 0 0
> exec ^
: string
| set a = 1, terminated = 2
< error
| 1064 42000 …line 1 column 21 near "terminated = 2" 
> exec ^
: string
| set a = 1, text = 2
< ok 0 0
> exec ^
: string
| set a = 1, than = 2
< ok 0 0
> exec ^
: string
| set a = 1, then = 2
< error
| 1064 42000 …line 1 column 15 near "then = 2" 
> exec ^
: string
| set a = 1, thread_priority = 2
< ok 0 0
> exec ^
: string
| set a = 1, ties = 2
< ok 0 0
> exec ^
: string
| set a = 1, tiflash = 2
< ok 0 0
> exec ^
: string
| set a = 1, time = 2
< ok 0 0
> exec ^
: string
| set a = 1, timeout = 2
< ok 0 0
> exec ^
: string
| set a = 1, timestamp = 2
< ok 0 0
> exec ^
: string
| set a = 1, timestampadd = 2
< ok 0 0
> exec ^
: string
| set a = 1, timestampdiff = 2
< ok 0 0
> exec ^
: string
| set a = 1, tinyblob = 2
< error
| 1064 42000 …line 1 column 19 near "tinyblob = 2" 
> exec ^
: string
| set a = 1, tinyint = 2
< error
| 1064 42000 …line 1 column 18 near "tinyint = 2" 
> exec ^
: string
| set a = 1, tinytext = 2
< error
| 1064 42000 …line 1 column 19 near "tinytext = 2" 
> exec ^
: string
| set a = 1, tls = 2
< ok 0 0
> exec ^
: string
| set a = 1, to = 2
< error
| 1064 42000 …line 1 column 13 near "to = 2" 
> exec ^
: string
| set a = 1, tpcc = 2
< ok 0 0
> exec ^
: string
| set a = 1, tpch_10 = 2
< ok 0 0
> exec ^
: string
| set a = 1, trailing = 2
< error
| 1064 42000 …line 1 column 19 near "trailing = 2" 
> exec ^
: string
| set a = 1, transaction = 2
< ok 0 0
> exec ^
: string
| set a = 1, trigger = 2
< error
| 1064 42000 …line 1 column 18 near "trigger = 2" 
> exec ^
: string
| set a = 1, triggers = 2
< ok 0 0
> exec ^
: string
| set a = 1, trim = 2
< ok 0 0
> exec ^
: string
| set a = 1, true = 2
< error
| 1064 42000 …line 1 column 15 near "true = 2" 
> exec ^
: string
| set a = 1, truncate = 2
< ok 0 0
> exec ^
: string
| set a = 1, ttl = 2
< ok 0 0
> exec ^
: string
| set a = 1, ttl_enable = 2
< ok 0 0
> exec ^
: string
| set a = 1, ttl_job_interval = 2
< ok 0 0
> exec ^
: string
| set a = 1, type = 2
< ok 0 0
> exec ^
: string
| set a = 1, types = 2
< ok 0 0
> exec ^
: string
| set a = 1, unbounded = 2
< ok 0 0
> exec ^
: string
| set a = 1, uncommitted = 2
< ok 0 0
> exec ^
: string
| set a = 1, undefined = 2
< ok 0 0
> exec ^
: string
| set a = 1, undo = 2
< ok 0 0
> exec ^
: string
| set a = 1, undo_buffer_size = 2
< ok 0 0
> exec ^
: string
| set a = 1, undofile = 2
< ok 0 0
> exec ^
: string
| set a = 1, unicode = 2
< ok 0 0
> exec ^
: string
| set a = 1, uninstall = 2
< ok 0 0
> exec ^
: string
| set a = 1, union = 2
< error
| 1064 42000 …line 1 column 16 near "union = 2" 
> exec ^
: string
| set a = 1, unique = 2
< error
| 1064 42000 …line 1 column 17 near "unique = 2" 
> exec ^
: string
| set a = 1, unknown = 2
< ok 0 0
> exec ^
: string
| set a = 1, unlimited = 2
< ok 0 0
> exec ^
: string
| set a = 1, unlock = 2
< error
| 1064 42000 …line 1 column 17 near "unlock = 2" 
> exec ^
: string
| set a = 1, unregister = 2
< ok 0 0
> exec ^
: string
| set a = 1, unsigned = 2
< error
| 1064 42000 …line 1 column 19 near "unsigned = 2" 
> exec ^
: string
| set a = 1, until = 2
< error
| 1064 42000 …line 1 column 16 near "until = 2" 
> exec ^
: string
| set a = 1, update = 2
< error
| 1064 42000 …line 1 column 17 near "update = 2" 
> exec ^
: string
| set a = 1, upgrade = 2
< ok 0 0
> exec ^
: string
| set a = 1, uri = 2
< ok 0 0
> exec ^
: string
| set a = 1, url = 2
< ok 0 0
> exec ^
: string
| set a = 1, usage = 2
< error
| 1064 42000 …line 1 column 16 near "usage = 2" 
> exec ^
: string
| set a = 1, use = 2
< error
| 1064 42000 …line 1 column 14 near "use = 2" 
> exec ^
: string
| set a = 1, use_frm = 2
< ok 0 0
> exec ^
: string
| set a = 1, user = 2
< ok 0 0
> exec ^
: string
| set a = 1, user_resources = 2
< ok 0 0
> exec ^
: string
| set a = 1, using = 2
< error
| 1064 42000 …line 1 column 16 near "using = 2" 
> exec ^
: string
| set a = 1, utc_date = 2
< error
| 1064 42000 …line 1 column 19 near "utc_date = 2" 
> exec ^
: string
| set a = 1, utc_time = 2
< error
| 1064 42000 …line 1 column 19 near "utc_time = 2" 
> exec ^
: string
| set a = 1, utc_timestamp = 2
< error
| 1064 42000 …line 1 column 24 near "utc_timestamp = 2" 
> exec ^
: string
| set a = 1, utilization_limit = 2
< ok 0 0
> exec ^
: string
| set a = 1, validate = 2
< ok 0 0
> exec ^
: string
| set a = 1, validation = 2
< ok 0 0
> exec ^
: string
| set a = 1, value = 2
< ok 0 0
> exec ^
: string
| set a = 1, values = 2
< error
| 1064 42000 …line 1 column 17 near "values = 2" 
> exec ^
: string
| set a = 1, var_pop = 2
< ok 0 0
> exec ^
: string
| set a = 1, var_samp = 2
< ok 0 0
> exec ^
: string
| set a = 1, varbinary = 2
< error
| 1064 42000 …line 1 column 20 near "varbinary = 2" 
> exec ^
: string
| set a = 1, varchar = 2
< error
| 1064 42000 …line 1 column 18 near "varchar = 2" 
> exec ^
: string
| set a = 1, varcharacter = 2
< error
| 1064 42000 …line 1 column 23 near "varcharacter = 2" 
> exec ^
: string
| set a = 1, variables = 2
< ok 0 0
> exec ^
: string
| set a = 1, variance = 2
< ok 0 0
> exec ^
: string
| set a = 1, varying = 2
< error
| 1064 42000 …line 1 column 18 near "varying = 2" 
> exec ^
: string
| set a = 1, vcpu = 2
< ok 0 0
> exec ^
: string
| set a = 1, vector = 2
< ok 0 0
> exec ^
: string
| set a = 1, verify_key_constraints = 2
< ok 0 0
> exec ^
: string
| set a = 1, view = 2
< ok 0 0
> exec ^
: string
| set a = 1, virtual = 2
< error
| 1064 42000 …line 1 column 18 near "virtual = 2" 
> exec ^
: string
| set a = 1, visible = 2
< ok 0 0
> exec ^
: string
| set a = 1, voter_constraints = 2
< ok 0 0
> exec ^
: string
| set a = 1, voters = 2
< ok 0 0
> exec ^
: string
| set a = 1, wait = 2
< ok 0 0
> exec ^
: string
| set a = 1, warnings = 2
< ok 0 0
> exec ^
: string
| set a = 1, watch = 2
< ok 0 0
> exec ^
: string
| set a = 1, week = 2
< ok 0 0
> exec ^
: string
| set a = 1, weight_string = 2
< ok 0 0
> exec ^
: string
| set a = 1, when = 2
< error
| 1064 42000 …line 1 column 15 near "when = 2" 
> exec ^
: string
| set a = 1, where = 2
< error
| 1064 42000 …line 1 column 16 near "where = 2" 
> exec ^
: string
| set a = 1, while = 2
< error
| 1064 42000 …line 1 column 16 near "while = 2" 
> exec ^
: string
| set a = 1, window = 2
< error
| 1064 42000 …line 1 column 17 near "window = 2" 
> exec ^
: string
| set a = 1, with = 2
< error
| 1064 42000 …line 1 column 15 near "with = 2" 
> exec ^
: string
| set a = 1, without = 2
< ok 0 0
> exec ^
: string
| set a = 1, work = 2
< ok 0 0
> exec ^
: string
| set a = 1, workload = 2
< ok 0 0
> exec ^
: string
| set a = 1, wrapper = 2
< ok 0 0
> exec ^
: string
| set a = 1, write = 2
< error
| 1064 42000 …line 1 column 16 near "write = 2" 
> exec ^
: string
| set a = 1, x509 = 2
< ok 0 0
> exec ^
: string
| set a = 1, xa = 2
< ok 0 0
> exec ^
: string
| set a = 1, xid = 2
< ok 0 0
> exec ^
: string
| set a = 1, xml = 2
< ok 0 0
> exec ^
: string
| set a = 1, xor = 2
< error
| 1064 42000 …line 1 column 14 near "xor = 2" 
> exec ^
: string
| set a = 1, year = 2
< ok 0 0
> exec ^
: string
| set a = 1, year_month = 2
< error
| 1064 42000 …line 1 column 21 near "year_month = 2" 
> exec ^
: string
| set a = 1, zerofill = 2
< error
| 1064 42000 …line 1 column 19 near "zerofill = 2" 
> exec ^
: string
| set a = 1, zone = 2
< ok 0 0
> exec ^
: string
| set absent = 1, a = 2
< ok 0 0
> exec ^
: string
| set accessible = 1, a = 2
< ok 0 0
> exec ^
: string
| set account = 1, a = 2
< ok 0 0
> exec ^
: string
| set action = 1, a = 2
< ok 0 0
> exec ^
: string
| set active = 1, a = 2
< ok 0 0
> exec ^
: string
| set add = 1, a = 2
< error
| 1064 42000 …line 1 column 7 near "add = 1, a = 2" 
> exec ^
: string
| set adddate = 1, a = 2
< ok 0 0
> exec ^
: string
| set admin = 1, a = 2
< ok 0 0
> exec ^
: string
| set after = 1, a = 2
< ok 0 0
> exec ^
: string
| set against = 1, a = 2
< ok 0 0
> exec ^
: string
| set aggregate = 1, a = 2
< ok 0 0
> exec ^
: string
| set algorithm = 1, a = 2
< ok 0 0
> exec ^
: string
| set all = 1, a = 2
< error
| 1064 42000 …line 1 column 7 near "all = 1, a = 2" 
> exec ^
: string
| set allow_missing_files = 1, a = 2
< ok 0 0
> exec ^
: string
| set alter = 1, a = 2
< error
| 1064 42000 …line 1 column 9 near "alter = 1, a = 2" 
> exec ^
: string
| set always = 1, a = 2
< ok 0 0
> exec ^
: string
| set analyze = 1, a = 2
< error
| 1064 42000 …line 1 column 11 near "analyze = 1, a = 2" 
> exec ^
: string
| set and = 1, a = 2
< error
| 1064 42000 …line 1 column 7 near "and = 1, a = 2" 
> exec ^
: string
| set any = 1, a = 2
< ok 0 0
> exec ^
: string
| set array = 1, a = 2
< error
| 1064 42000 …line 1 column 9 near "array = 1, a = 2" 
> exec ^
: string
| set as = 1, a = 2
< error
| 1064 42000 …line 1 column 6 near "as = 1, a = 2" 
> exec ^
: string
| set asc = 1, a = 2
< error
| 1064 42000 …line 1 column 7 near "asc = 1, a = 2" 
> exec ^
: string
| set ascii = 1, a = 2
< ok 0 0
> exec ^
: string
| set asensitive = 1, a = 2
< ok 0 0
> exec ^
: string
| set assign_gtids_to_anonymous_transactions = 1, a = 2
< ok 0 0
> exec ^
: string
| set at = 1, a = 2
< ok 0 0
> exec ^
: string
| set attribute = 1, a = 2
< ok 0 0
> exec ^
: string
| set authentication = 1, a = 2
< ok 0 0
> exec ^
: string
| set auto = 1, a = 2
< ok 0 0
> exec ^
: string
| set auto_id_cache = 1, a = 2
< ok 0 0
> exec ^
: string
| set auto_increment = 1, a = 2
< ok 0 0
> exec ^
: string
| set auto_random = 1, a = 2
< ok 0 0
> exec ^
: string
| set auto_random_base = 1, a = 2
< ok 0 0
> exec ^
: string
| set auto_refresh = 1, a = 2
< ok 0 0
> exec ^
: string
| set auto_refresh_source = 1, a = 2
< ok 0 0
> exec ^
: string
| set autoextend_size = 1, a = 2
< ok 0 0
> exec ^
: string
| set avg = 1, a = 2
< ok 0 0
> exec ^
: string
| set avg_row_length = 1, a = 2
< ok 0 0
> exec ^
: string
| set background = 1, a = 2
< ok 0 0
> exec ^
: string
| set backup = 1, a = 2
< ok 0 0
> exec ^
: string
| set backups = 1, a = 2
< ok 0 0
> exec ^
: string
| set batch = 1, a = 2
< ok 0 0
> exec ^
: string
| set before = 1, a = 2
< ok 0 0
> exec ^
: string
| set begin = 1, a = 2
< ok 0 0
> exec ^
: string
| set bernoulli = 1, a = 2
< ok 0 0
> exec ^
: string
| set between = 1, a = 2
< error
| 1064 42000 …line 1 column 11 near "between = 1, a = 2" 
> exec ^
: string
| set bigint = 1, a = 2
< error
| 1064 42000 …line 1 column 10 near "bigint = 1, a = 2" 
> exec ^
: string
| set binary = 1, a = 2
< error
| 1064 42000 …line 1 column 10 near "binary = 1, a = 2" 
> exec ^
: string
| set binding = 1, a = 2
< ok 0 0
> exec ^
: string
| set binding_cache = 1, a = 2
< ok 0 0
> exec ^
: string
| set bindings = 1, a = 2
< ok 0 0
> exec ^
: string
| set binlog = 1, a = 2
< ok 0 0
> exec ^
: string
| set bit = 1, a = 2
< ok 0 0
> exec ^
: string
| set bit_and = 1, a = 2
< ok 0 0
> exec ^
: string
| set bit_or = 1, a = 2
< ok 0 0
> exec ^
: string
| set bit_xor = 1, a = 2
< ok 0 0
> exec ^
: string
| set blob = 1, a = 2
< error
| 1064 42000 …line 1 column 8 near "blob = 1, a = 2" 
> exec ^
: string
| set block = 1, a = 2
< ok 0 0
> exec ^
: string
| set bool = 1, a = 2
< ok 0 0
> exec ^
: string
| set boolean = 1, a = 2
< ok 0 0
> exec ^
: string
| set both = 1, a = 2
< error
| 1064 42000 …line 1 column 8 near "both = 1, a = 2" 
> exec ^
: string
| set btree = 1, a = 2
< ok 0 0
> exec ^
: string
| set buckets = 1, a = 2
< ok 0 0
> exec ^
: string
| set builtins = 1, a = 2
< ok 0 0
> exec ^
: string
| set bulk = 1, a = 2
< ok 0 0
> exec ^
: string
| set burstable = 1, a = 2
< ok 0 0
> exec ^
: string
| set by = 1, a = 2
< error
| 1064 42000 …line 1 column 6 near "by = 1, a = 2" 
> exec ^
: string
| set byte = 1, a = 2
< ok 0 0
> exec ^
: string
| set cache = 1, a = 2
< ok 0 0
> exec ^
: string
| set calibrate = 1, a = 2
< ok 0 0
> exec ^
: string
| set call = 1, a = 2
< error
| 1064 42000 …line 1 column 8 near "call = 1, a = 2" 
> exec ^
: string
| set cascade = 1, a = 2
< error
| 1064 42000 …line 1 column 11 near "cascade = 1, a = 2" 
> exec ^
: string
| set cascaded = 1, a = 2
< ok 0 0
> exec ^
: string
| set case = 1, a = 2
< error
| 1064 42000 …line 1 column 8 near "case = 1, a = 2" 
> exec ^
: string
| set cast = 1, a = 2
< ok 0 0
> exec ^
: string
| set catalog_name = 1, a = 2
< ok 0 0
> exec ^
: string
| set chain = 1, a = 2
< ok 0 0
> exec ^
: string
| set challenge_response = 1, a = 2
< ok 0 0
> exec ^
: string
| set change = 1, a = 2
< error
| 1064 42000 …line 1 column 10 near "change = 1, a = 2" 
> exec ^
: string
| set changed = 1, a = 2
< ok 0 0
> exec ^
: string
| set channel = 1, a = 2
< ok 0 0
> exec ^
: string
| set char = 1, a = 2
< error
| 1064 42000 …line 1 column 10 near "= 1, a = 2" 
> exec ^
: string
| set character = 1, a = 2
< error
| 1064 42000 …line 1 column 15 near "= 1, a = 2" 
> exec ^
: string
| set charset = 1, a = 2
< ok 0 0
> exec ^
: string
| set check = 1, a = 2
< error
| 1064 42000 …line 1 column 9 near "check = 1, a = 2" 
> exec ^
: string
| set checksum = 1, a = 2
< ok 0 0
> exec ^
: string
| set cipher = 1, a = 2
< ok 0 0
> exec ^
: string
| set class_origin = 1, a = 2
< ok 0 0
> exec ^
: string
| set client = 1, a = 2
< ok 0 0
> exec ^
: string
| set clone = 1, a = 2
< ok 0 0
> exec ^
: string
| set close = 1, a = 2
< ok 0 0
> exec ^
: string
| set clustered = 1, a = 2
< ok 0 0
> exec ^
: string
| set coalesce = 1, a = 2
< ok 0 0
> exec ^
: string
| set code = 1, a = 2
< ok 0 0
> exec ^
: string
| set collate = 1, a = 2
< error
| 1064 42000 …line 1 column 11 near "collate = 1, a = 2" 
> exec ^
: string
| set collation = 1, a = 2
< ok 0 0
> exec ^
: string
| set column = 1, a = 2
< error
| 1064 42000 …line 1 column 10 near "column = 1, a = 2" 
> exec ^
: string
| set column_format = 1, a = 2
< ok 0 0
> exec ^
: string
| set column_name = 1, a = 2
< ok 0 0
> exec ^
: string
| set column_stats_usage = 1, a = 2
< ok 0 0
> exec ^
: string
| set columns = 1, a = 2
< ok 0 0
> exec ^
: string
| set comment = 1, a = 2
< ok 0 0
> exec ^
: string
| set commit = 1, a = 2
< ok 0 0
> exec ^
: string
| set committed = 1, a = 2
< ok 0 0
> exec ^
: string
| set compact = 1, a = 2
< ok 0 0
> exec ^
: string
| set completion = 1, a = 2
< ok 0 0
> exec ^
: string
| set component = 1, a = 2
< ok 0 0
> exec ^
: string
| set compressed = 1, a = 2
< ok 0 0
> exec ^
: string
| set compression = 1, a = 2
< ok 0 0
> exec ^
: string
| set concat = 1, a = 2
< ok 0 0
> exec ^
: string
| set concurrent = 1, a = 2
< ok 0 0
> exec ^
: string
| set condition = 1, a = 2
< ok 0 0
> exec ^
: string
| set config = 1, a = 2
< ok 0 0
> exec ^
: string
| set connection = 1, a = 2
< ok 0 0
> exec ^
: string
| set consistent = 1, a = 2
< ok 0 0
> exec ^
: string
| set constraint = 1, a = 2
< error
| 1064 42000 …line 1 column 14 near "constraint = 1, a = 2" 
> exec ^
: string
| set constraint_catalog = 1, a = 2
< ok 0 0
> exec ^
: string
| set constraint_name = 1, a = 2
< ok 0 0
> exec ^
: string
| set constraint_schema = 1, a = 2
< ok 0 0
> exec ^
: string
| set constraints = 1, a = 2
< ok 0 0
> exec ^
: string
| set contains = 1, a = 2
< ok 0 0
> exec ^
: string
| set context = 1, a = 2
< ok 0 0
> exec ^
: string
| set continue = 1, a = 2
< error
| 1064 42000 …line 1 column 12 near "continue = 1, a = 2" 
> exec ^
: string
| set convert = 1, a = 2
< error
| 1064 42000 …line 1 column 11 near "convert = 1, a = 2" 
> exec ^
: string
| set cooldown = 1, a = 2
< ok 0 0
> exec ^
: string
| set count = 1, a = 2
< ok 0 0
> exec ^
: string
| set cpu = 1, a = 2
< ok 0 0
> exec ^
: string
| set create = 1, a = 2
< error
| 1064 42000 …line 1 column 10 near "create = 1, a = 2" 
> exec ^
: string
| set cross = 1, a = 2
< error
| 1064 42000 …line 1 column 9 near "cross = 1, a = 2" 
> exec ^
: string
| set cube = 1, a = 2
< ok 0 0
> exec ^
: string
| set cume_dist = 1, a = 2
< error
| 1064 42000 …line 1 column 13 near "cume_dist = 1, a = 2" 
> exec ^
: string
| set curdate = 1, a = 2
< ok 0 0
> exec ^
: string
| set current = 1, a = 2
< ok 0 0
> exec ^
: string
| set current_date = 1, a = 2
< error
| 1064 42000 …line 1 column 16 near "current_date = 1, a = 2" 
> exec ^
: string
| set current_time = 1, a = 2
< error
| 1064 42000 …line 1 column 16 near "current_time = 1, a = 2" 
> exec ^
: string
| set current_timestamp = 1, a = 2
< error
| 1064 42000 …line 1 column 21 near "current_timestamp = 1, a = 2" 
> exec ^
: string
| set current_user = 1, a = 2
< error
| 1064 42000 …line 1 column 16 near "current_user = 1, a = 2" 
> exec ^
: string
| set cursor = 1, a = 2
< error
| 1064 42000 …line 1 column 10 near "cursor = 1, a = 2" 
> exec ^
: string
| set cursor_name = 1, a = 2
< ok 0 0
> exec ^
: string
| set curtime = 1, a = 2
< ok 0 0
> exec ^
: string
| set cycle = 1, a = 2
< ok 0 0
> exec ^
: string
| set data = 1, a = 2
< ok 0 0
> exec ^
: string
| set database = 1, a = 2
< error
| 1064 42000 …line 1 column 12 near "database = 1, a = 2" 
> exec ^
: string
| set databases = 1, a = 2
< error
| 1064 42000 …line 1 column 13 near "databases = 1, a = 2" 
> exec ^
: string
| set datafile = 1, a = 2
< ok 0 0
> exec ^
: string
| set date = 1, a = 2
< ok 0 0
> exec ^
: string
| set date_add = 1, a = 2
< ok 0 0
> exec ^
: string
| set date_sub = 1, a = 2
< ok 0 0
> exec ^
: string
| set datetime = 1, a = 2
< ok 0 0
> exec ^
: string
| set day = 1, a = 2
< ok 0 0
> exec ^
: string
| set day_hour = 1, a = 2
< error
| 1064 42000 …line 1 column 12 near "day_hour = 1, a = 2" 
> exec ^
: string
| set day_microsecond = 1, a = 2
< error
| 1064 42000 …line 1 column 19 near "day_microsecond = 1, a = 2" 
> exec ^
: string
| set day_minute = 1, a = 2
< error
| 1064 42000 …line 1 column 14 near "day_minute = 1, a = 2" 
> exec ^
: string
| set day_second = 1, a = 2
< error
| 1064 42000 …line 1 column 14 near "day_second = 1, a = 2" 
> exec ^
: string
| set deallocate = 1, a = 2
< ok 0 0
> exec ^
: string
| set dec = 1, a = 2
< error
| 1064 42000 …line 1 column 7 near "dec = 1, a = 2" 
> exec ^
: string
| set decimal = 1, a = 2
< error
| 1064 42000 …line 1 column 11 near "decimal = 1, a = 2" 
> exec ^
: string
| set declare = 1, a = 2
< ok 0 0
> exec ^
: string
| set default = 1, a = 2
< error
| 1064 42000 …line 1 column 13 near "= 1, a = 2" 
> exec ^
: string
| set default_auth = 1, a = 2
< ok 0 0
> exec ^
: string
| set defined = 1, a = 2
< ok 0 0
> exec ^
: string
| set definer = 1, a = 2
< ok 0 0
> exec ^
: string
| set definition = 1, a = 2
< ok 0 0
> exec ^
: string
| set delay_key_write = 1, a = 2
< ok 0 0
> exec ^
: string
| set delayed = 1, a = 2
< error
| 1064 42000 …line 1 column 11 near "delayed = 1, a = 2" 
> exec ^
: string
| set delete = 1, a = 2
< error
| 1064 42000 …line 1 column 10 near "delete = 1, a = 2" 
> exec ^
: string
| set dense_rank = 1, a = 2
< error
| 1064 42000 …line 1 column 14 near "dense_rank = 1, a = 2" 
> exec ^
: string
| set desc = 1, a = 2
< error
| 1064 42000 …line 1 column 8 near "desc = 1, a = 2" 
> exec ^
: string
| set describe = 1, a = 2
< error
| 1064 42000 …line 1 column 12 near "describe = 1, a = 2" 
> exec ^
: string
| set description = 1, a = 2
< ok 0 0
> exec ^
: string
| set deterministic = 1, a = 2
< ok 0 0
> exec ^
: string
| set diagnostics = 1, a = 2
< ok 0 0
> exec ^
: string
| set digest = 1, a = 2
< ok 0 0
> exec ^
: string
| set directory = 1, a = 2
< ok 0 0
> exec ^
: string
| set disable = 1, a = 2
< ok 0 0
> exec ^
: string
| set disabled = 1, a = 2
< ok 0 0
> exec ^
: string
| set discard = 1, a = 2
< ok 0 0
> exec ^
: string
| set disk = 1, a = 2
< ok 0 0
> exec ^
: string
| set distinct = 1, a = 2
< error
| 1064 42000 …line 1 column 12 near "distinct = 1, a = 2" 
> exec ^
: string
| set distinctrow = 1, a = 2
< error
| 1064 42000 …line 1 column 15 near "distinctrow = 1, a = 2" 
> exec ^
: string
| set distribute = 1, a = 2
< ok 0 0
> exec ^
: string
| set distribution = 1, a = 2
< ok 0 0
> exec ^
: string
| set distributions = 1, a = 2
< ok 0 0
> exec ^
: string
| set div = 1, a = 2
< error
| 1064 42000 …line 1 column 7 near "div = 1, a = 2" 
> exec ^
: string
| set do = 1, a = 2
< ok 0 0
> exec ^
: string
| set double = 1, a = 2
< error
| 1064 42000 …line 1 column 10 near "double = 1, a = 2" 
> exec ^
: string
| set drop = 1, a = 2
< error
| 1064 42000 …line 1 column 8 near "drop = 1, a = 2" 
> exec ^
: string
| set dry = 1, a = 2
< ok 0 0
> exec ^
: string
| set dryrun = 1, a = 2
< ok 0 0
> exec ^
: string
| set dual = 1, a = 2
< error
| 1064 42000 …line 1 column 8 near "dual = 1, a = 2" 
> exec ^
: string
| set duality = 1, a = 2
< ok 0 0
> exec ^
: string
| set dumpfile = 1, a = 2
< ok 0 0
> exec ^
: string
| set duplicate = 1, a = 2
< ok 0 0
> exec ^
: string
| set duration = 1, a = 2
< ok 0 0
> exec ^
: string
| set dynamic = 1, a = 2
< ok 0 0
> exec ^
: string
| set each = 1, a = 2
< ok 0 0
> exec ^
: string
| set else = 1, a = 2
< error
| 1064 42000 …line 1 column 8 near "else = 1, a = 2" 
> exec ^
: string
| set elseif = 1, a = 2
< error
| 1064 42000 …line 1 column 10 near "elseif = 1, a = 2" 
> exec ^
: string
| set empty = 1, a = 2
< ok 0 0
> exec ^
: string
| set enable = 1, a = 2
< ok 0 0
> exec ^
: string
| set enabled = 1, a = 2
< ok 0 0
> exec ^
: string
| set enclosed = 1, a = 2
< error
| 1064 42000 …line 1 column 12 near "enclosed = 1, a = 2" 
> exec ^
: string
| set encryption = 1, a = 2
< ok 0 0
> exec ^
: string
| set end = 1, a = 2
< ok 0 0
> exec ^
: string
| set end_time = 1, a = 2
< ok 0 0
> exec ^
: string
| set ends = 1, a = 2
< ok 0 0
> exec ^
: string
| set enforced = 1, a = 2
< ok 0 0
> exec ^
: string
| set engine = 1, a = 2
< ok 0 0
> exec ^
: string
| set engine_attribute = 1, a = 2
< ok 0 0
> exec ^
: string
| set engines = 1, a = 2
< ok 0 0
> exec ^
: string
| set enum = 1, a = 2
< ok 0 0
> exec ^
: string
| set error = 1, a = 2
< ok 0 0
> exec ^
: string
| set errors = 1, a = 2
< ok 0 0
> exec ^
: string
| set escape = 1, a = 2
< ok 0 0
> exec ^
: string
| set escaped = 1, a = 2
< error
| 1064 42000 …line 1 column 11 near "escaped = 1, a = 2" 
> exec ^
: string
| set event = 1, a = 2
< ok 0 0
> exec ^
: string
| set events = 1, a = 2
< ok 0 0
> exec ^
: string
| set every = 1, a = 2
< ok 0 0
> exec ^
: string
| set exact = 1, a = 2
< ok 0 0
> exec ^
: string
| set except = 1, a = 2
< error
| 1064 42000 …line 1 column 10 near "except = 1, a = 2" 
> exec ^
: string
| set exchange = 1, a = 2
< ok 0 0
> exec ^
: string
| set exclude = 1, a = 2
< ok 0 0
> exec ^
: string
| set exec_elapsed = 1, a = 2
< ok 0 0
> exec ^
: string
| set execute = 1, a = 2
< ok 0 0
> exec ^
: string
| set exists = 1, a = 2
< error
| 1064 42000 …line 1 column 10 near "exists = 1, a = 2" 
> exec ^
: string
| set exit = 1, a = 2
< error
| 1064 42000 …line 1 column 8 near "exit = 1, a = 2" 
> exec ^
: string
| set expansion = 1, a = 2
< ok 0 0
> exec ^
: string
| set expire = 1, a = 2
< ok 0 0
> exec ^
: string
| set explain = 1, a = 2
< error
| 1064 42000 …line 1 column 11 near "explain = 1, a = 2" 
> exec ^
: string
| set export = 1, a = 2
< ok 0 0
> exec ^
: string
| set extended = 1, a = 2
< ok 0 0
> exec ^
: string
| set extent_size = 1, a = 2
< ok 0 0
> exec ^
: string
| set external = 1, a = 2
< ok 0 0
> exec ^
: string
| set external_format = 1, a = 2
< ok 0 0
> exec ^
: string
| set extract = 1, a = 2
< ok 0 0
> exec ^
: string
| set factor = 1, a = 2
< ok 0 0
> exec ^
: string
| set failed_login_attempts = 1, a = 2
< ok 0 0
> exec ^
: string
| set false = 1, a = 2
< error
| 1064 42000 …line 1 column 9 near "false = 1, a = 2" 
> exec ^
: string
| set fast = 1, a = 2
< ok 0 0
> exec ^
: string
| set faults = 1, a = 2
< ok 0 0
> exec ^
: string
| set fetch = 1, a = 2
< error
| 1064 42000 …line 1 column 9 near "fetch = 1, a = 2" 
> exec ^
: string
| set fields = 1, a = 2
< ok 0 0
> exec ^
: string
| set file = 1, a = 2
< ok 0 0
> exec ^
: string
| set file_block_size = 1, a = 2
< ok 0 0
> exec ^
: string
| set file_format = 1, a = 2
< ok 0 0
> exec ^
: string
| set file_name = 1, a = 2
< ok 0 0
> exec ^
: string
| set file_pattern = 1, a = 2
< ok 0 0
> exec ^
: string
| set file_prefix = 1, a = 2
< ok 0 0
> exec ^
: string
| set files = 1, a = 2
< ok 0 0
> exec ^
: string
| set filter = 1, a = 2
< ok 0 0
> exec ^
: string
| set finish = 1, a = 2
< ok 0 0
> exec ^
: string
| set first = 1, a = 2
< ok 0 0
> exec ^
: string
| set first_value = 1, a = 2
< error
| 1064 42000 …line 1 column 15 near "first_value = 1, a = 2" 
> exec ^
: string
| set fixed = 1, a = 2
< ok 0 0
> exec ^
: string
| set flashback = 1, a = 2
< ok 0 0
> exec ^
: string
| set float = 1, a = 2
< error
| 1064 42000 …line 1 column 9 near "float = 1, a = 2" 
> exec ^
: string
| set float4 = 1, a = 2
< error
| 1064 42000 …line 1 column 10 near "float4 = 1, a = 2" 
> exec ^
: string
| set float8 = 1, a = 2
< error
| 1064 42000 …line 1 column 10 near "float8 = 1, a = 2" 
> exec ^
: string
| set flush = 1, a = 2
< ok 0 0
> exec ^
: string
| set follower_constraints = 1, a = 2
< ok 0 0
> exec ^
: string
| set followers = 1, a = 2
< ok 0 0
> exec ^
: string
| set following = 1, a = 2
< ok 0 0
> exec ^
: string
| set follows = 1, a = 2
< ok 0 0
> exec ^
: string
| set for = 1, a = 2
< error
| 1064 42000 …line 1 column 7 near "for = 1, a = 2" 
> exec ^
: string
| set force = 1, a = 2
< error
| 1064 42000 …line 1 column 9 near "force = 1, a = 2" 
> exec ^
: string
| set foreign = 1, a = 2
< error
| 1064 42000 …line 1 column 11 near "foreign = 1, a = 2" 
> exec ^
: string
| set format = 1, a = 2
< ok 0 0
> exec ^
: string
| set found = 1, a = 2
< ok 0 0
> exec ^
: string
| set from = 1, a = 2
< error
| 1064 42000 …line 1 column 8 near "from = 1, a = 2" 
> exec ^
: string
| set full = 1, a = 2
< ok 0 0
> exec ^
: string
| set fulltext = 1, a = 2
< error
| 1064 42000 …line 1 column 12 near "fulltext = 1, a = 2" 
> exec ^
: string
| set function = 1, a = 2
< ok 0 0
> exec ^
: string
| set general = 1, a = 2
< ok 0 0
> exec ^
: string
| set generate = 1, a = 2
< ok 0 0
> exec ^
: string
| set generated = 1, a = 2
< error
| 1064 42000 …line 1 column 13 near "generated = 1, a = 2" 
> exec ^
: string
| set geomcollection = 1, a = 2
< ok 0 0
> exec ^
: string
| set geometry = 1, a = 2
< ok 0 0
> exec ^
: string
| set geometrycollection = 1, a = 2
< ok 0 0
> exec ^
: string
| set get = 1, a = 2
< ok 0 0
> exec ^
: string
| set get_format = 1, a = 2
< ok 0 0
> exec ^
: string
| set get_source_public_key = 1, a = 2
< ok 0 0
> exec ^
: string
| set global = 1, a = 2
< ok 0 0
> exec ^
: string
| set grant = 1, a = 2
< error
| 1064 42000 …line 1 column 9 near "grant = 1, a = 2" 
> exec ^
: string
| set grants = 1, a = 2
< ok 0 0
> exec ^
: string
| set greatest = 1, a = 2
< ok 0 0
> exec ^
: string
| set group = 1, a = 2
< error
| 1064 42000 …line 1 column 9 near "group = 1, a = 2" 
> exec ^
: string
| set group_concat = 1, a = 2
< ok 0 0
> exec ^
: string
| set group_replication = 1, a = 2
< ok 0 0
> exec ^
: string
| set grouping = 1, a = 2
< ok 0 0
> exec ^
: string
| set groups = 1, a = 2
< error
| 1064 42000 …line 1 column 10 near "groups = 1, a = 2" 
> exec ^
: string
| set gtid_only = 1, a = 2
< ok 0 0
> exec ^
: string
| set gtids = 1, a = 2
< ok 0 0
> exec ^
: string
| set guided = 1, a = 2
< ok 0 0
> exec ^
: string
| set handler = 1, a = 2
< ok 0 0
> exec ^
: string
| set hash = 1, a = 2
< ok 0 0
> exec ^
: string
| set having = 1, a = 2
< error
| 1064 42000 …line 1 column 10 near "having = 1, a = 2" 
> exec ^
: string
| set header = 1, a = 2
< ok 0 0
> exec ^
: string
| set help = 1, a = 2
< ok 0 0
> exec ^
: string
| set high = 1, a = 2
< ok 0 0
> exec ^
: string
| set high_priority = 1, a = 2
< error
| 1064 42000 …line 1 column 17 near "high_priority = 1, a = 2" 
> exec ^
: string
| set histogram = 1, a = 2
< ok 0 0
> exec ^
: string
| set histograms_in_flight = 1, a = 2
< ok 0 0
> exec ^
: string
| set history = 1, a = 2
< ok 0 0
> exec ^
: string
| set host = 1, a = 2
< ok 0 0
> exec ^
: string
| set hosts = 1, a = 2
< ok 0 0
> exec ^
: string
| set hour = 1, a = 2
< ok 0 0
> exec ^
: string
| set hour_microsecond = 1, a = 2
< error
| 1064 42000 …line 1 column 20 near "hour_microsecond = 1, a = 2" 
> exec ^
: string
| set hour_minute = 1, a = 2
< error
| 1064 42000 …line 1 column 15 near "hour_minute = 1, a = 2" 
> exec ^
: string
| set hour_second = 1, a = 2
< error
| 1064 42000 …line 1 column 15 near "hour_second = 1, a = 2" 
> exec ^
: string
| set identified = 1, a = 2
< ok 0 0
> exec ^
: string
| set if = 1, a = 2
< error
| 1064 42000 …line 1 column 6 near "if = 1, a = 2" 
> exec ^
: string
| set ignore = 1, a = 2
< error
| 1064 42000 …line 1 column 10 near "ignore = 1, a = 2" 
> exec ^
: string
| set ignore_server_ids = 1, a = 2
< ok 0 0
> exec ^
: string
| set ilike = 1, a = 2
< error
| 1064 42000 …line 1 column 9 near "ilike = 1, a = 2" 
> exec ^
: string
| set import = 1, a = 2
< ok 0 0
> exec ^
: string
| set in = 1, a = 2
< error
| 1064 42000 …line 1 column 6 near "in = 1, a = 2" 
> exec ^
: string
| set inactive = 1, a = 2
< ok 0 0
> exec ^
: string
| set increment = 1, a = 2
< ok 0 0
> exec ^
: string
| set index = 1, a = 2
< error
| 1064 42000 …line 1 column 9 near "index = 1, a = 2" 
> exec ^
: string
| set indexes = 1, a = 2
< ok 0 0
> exec ^
: string
| set infile = 1, a = 2
< error
| 1064 42000 …line 1 column 10 near "infile = 1, a = 2" 
> exec ^
: string
| set initial = 1, a = 2
< ok 0 0
> exec ^
: string
| set initial_size = 1, a = 2
< ok 0 0
> exec ^
: string
| set initiate = 1, a = 2
< ok 0 0
> exec ^
: string
| set inner = 1, a = 2
< error
| 1064 42000 …line 1 column 9 near "inner = 1, a = 2" 
> exec ^
: string
| set inout = 1, a = 2
< error
| 1064 42000 …line 1 column 9 near "inout = 1, a = 2" 
> exec ^
: string
| set insensitive = 1, a = 2
< ok 0 0
> exec ^
: string
| set insert = 1, a = 2
< error
| 1064 42000 …line 1 column 10 near "insert = 1, a = 2" 
> exec ^
: string
| set insert_method = 1, a = 2
< ok 0 0
> exec ^
: string
| set install = 1, a = 2
< ok 0 0
> exec ^
: string
| set instance = 1, a = 2
< ok 0 0
> exec ^
: string
| set int = 1, a = 2
< error
| 1064 42000 …line 1 column 7 near "int = 1, a = 2" 
> exec ^
: string
| set int1 = 1, a = 2
< error
| 1064 42000 …line 1 column 8 near "int1 = 1, a = 2" 
> exec ^
: string
| set int2 = 1, a = 2
< error
| 1064 42000 …line 1 column 8 near "int2 = 1, a = 2" 
> exec ^
: string
| set int3 = 1, a = 2
< error
| 1064 42000 …line 1 column 8 near "int3 = 1, a = 2" 
> exec ^
: string
| set int4 = 1, a = 2
< error
| 1064 42000 …line 1 column 8 near "int4 = 1, a = 2" 
> exec ^
: string
| set int8 = 1, a = 2
< error
| 1064 42000 …line 1 column 8 near "int8 = 1, a = 2" 
> exec ^
: string
| set integer = 1, a = 2
< error
| 1064 42000 …line 1 column 11 near "integer = 1, a = 2" 
> exec ^
: string
| set intersect = 1, a = 2
< error
| 1064 42000 …line 1 column 13 near "intersect = 1, a = 2" 
> exec ^
: string
| set interval = 1, a = 2
< error
| 1064 42000 …line 1 column 12 near "interval = 1, a = 2" 
> exec ^
: string
| set into = 1, a = 2
< error
| 1064 42000 …line 1 column 8 near "into = 1, a = 2" 
> exec ^
: string
| set invisible = 1, a = 2
< ok 0 0
> exec ^
: string
| set invoker = 1, a = 2
< ok 0 0
> exec ^
: string
| set io = 1, a = 2
< ok 0 0
> exec ^
: string
| set io_after_gtids = 1, a = 2
< ok 0 0
> exec ^
: string
| set io_before_gtids = 1, a = 2
< ok 0 0
> exec ^
: string
| set io_thread = 1, a = 2
< ok 0 0
> exec ^
: string
| set ipc = 1, a = 2
< ok 0 0
> exec ^
: string
| set is = 1, a = 2
< error
| 1064 42000 …line 1 column 6 near "is = 1, a = 2" 
> exec ^
: string
| set isolation = 1, a = 2
< ok 0 0
> exec ^
: string
| set issuer = 1, a = 2
< ok 0 0
> exec ^
: string
| set iterate = 1, a = 2
< error
| 1064 42000 …line 1 column 11 near "iterate = 1, a = 2" 
> exec ^
: string
| set job = 1, a = 2
< ok 0 0
> exec ^
: string
| set jobs = 1, a = 2
< ok 0 0
> exec ^
: string
| set join = 1, a = 2
< error
| 1064 42000 …line 1 column 8 near "join = 1, a = 2" 
> exec ^
: string
| set json = 1, a = 2
< ok 0 0
> exec ^
: string
| set json_arrayagg = 1, a = 2
< ok 0 0
> exec ^
: string
| set json_duality_object = 1, a = 2
< ok 0 0
> exec ^
: string
| set json_objectagg = 1, a = 2
< ok 0 0
> exec ^
: string
| set json_table = 1, a = 2
< ok 0 0
> exec ^
: string
| set json_value = 1, a = 2
< ok 0 0
> exec ^
: string
| set key = 1, a = 2
< error
| 1064 42000 …line 1 column 7 near "key = 1, a = 2" 
> exec ^
: string
| set key_block_size = 1, a = 2
< ok 0 0
> exec ^
: string
| set keyring = 1, a = 2
< ok 0 0
> exec ^
: string
| set keys = 1, a = 2
< error
| 1064 42000 …line 1 column 8 near "keys = 1, a = 2" 
> exec ^
: string
| set kill = 1, a = 2
< error
| 1064 42000 …line 1 column 8 near "kill = 1, a = 2" 
> exec ^
: string
| set labels = 1, a = 2
< ok 0 0
> exec ^
: string
| set lag = 1, a = 2
< error
| 1064 42000 …line 1 column 7 near "lag = 1, a = 2" 
> exec ^
: string
| set language = 1, a = 2
< ok 0 0
> exec ^
: string
| set last = 1, a = 2
< ok 0 0
> exec ^
: string
| set last_value = 1, a = 2
< error
| 1064 42000 …line 1 column 14 near "last_value = 1, a = 2" 
> exec ^
: string
| set lateral = 1, a = 2
< ok 0 0
> exec ^
: string
| set lead = 1, a = 2
< error
| 1064 42000 …line 1 column 8 near "lead = 1, a = 2" 
> exec ^
: string
| set leader_constraints = 1, a = 2
< ok 0 0
> exec ^
: string
| set leading = 1, a = 2
< error
| 1064 42000 …line 1 column 11 near "leading = 1, a = 2" 
> exec ^
: string
| set learner_constraints = 1, a = 2
< ok 0 0
> exec ^
: string
| set learners = 1, a = 2
< ok 0 0
> exec ^
: string
| set least = 1, a = 2
< ok 0 0
> exec ^
: string
| set leave = 1, a = 2
< error
| 1064 42000 …line 1 column 9 near "leave = 1, a = 2" 
> exec ^
: string
| set leaves = 1, a = 2
< ok 0 0
> exec ^
: string
| set left = 1, a = 2
< error
| 1064 42000 …line 1 column 8 near "left = 1, a = 2" 
> exec ^
: string
| set less = 1, a = 2
< ok 0 0
> exec ^
: string
| set level = 1, a = 2
< ok 0 0
> exec ^
: string
| set library = 1, a = 2
< ok 0 0
> exec ^
: string
| set like = 1, a = 2
< error
| 1064 42000 …line 1 column 8 near "like = 1, a = 2" 
> exec ^
: string
| set limit = 1, a = 2
< error
| 1064 42000 …line 1 column 9 near "limit = 1, a = 2" 
> exec ^
: string
| set linear = 1, a = 2
< error
| 1064 42000 …line 1 column 10 near "linear = 1, a = 2" 
> exec ^
: string
| set lines = 1, a = 2
< error
| 1064 42000 …line 1 column 9 near "lines = 1, a = 2" 
> exec ^
: string
| set linestring = 1, a = 2
< ok 0 0
> exec ^
: string
| set list = 1, a = 2
< ok 0 0
> exec ^
: string
| set load = 1, a = 2
< error
| 1064 42000 …line 1 column 8 near "load = 1, a = 2" 
> exec ^
: string
| set local = 1, a = 2
< ok 0 0
> exec ^
: string
| set localtime = 1, a = 2
< error
| 1064 42000 …line 1 column 13 near "localtime = 1, a = 2" 
> exec ^
: string
| set localtimestamp = 1, a = 2
< error
| 1064 42000 …line 1 column 18 near "localtimestamp = 1, a = 2" 
> exec ^
: string
| set location = 1, a = 2
< ok 0 0
> exec ^
: string
| set lock = 1, a = 2
< error
| 1064 42000 …line 1 column 8 near "lock = 1, a = 2" 
> exec ^
: string
| set locked = 1, a = 2
< ok 0 0
> exec ^
: string
| set locks = 1, a = 2
< ok 0 0
> exec ^
: string
| set log = 1, a = 2
< ok 0 0
> exec ^
: string
| set logfile = 1, a = 2
< ok 0 0
> exec ^
: string
| set logs = 1, a = 2
< ok 0 0
> exec ^
: string
| set long = 1, a = 2
< error
| 1064 42000 …line 1 column 8 near "long = 1, a = 2" 
> exec ^
: string
| set longblob = 1, a = 2
< error
| 1064 42000 …line 1 column 12 near "longblob = 1, a = 2" 
> exec ^
: string
| set longtext = 1, a = 2
< error
| 1064 42000 …line 1 column 12 near "longtext = 1, a = 2" 
> exec ^
: string
| set loop = 1, a = 2
< ok 0 0
> exec ^
: string
| set low = 1, a = 2
< ok 0 0
> exec ^
: string
| set low_priority = 1, a = 2
< error
| 1064 42000 …line 1 column 16 near "low_priority = 1, a = 2" 
> exec ^
: string
| set manual = 1, a = 2
< ok 0 0
> exec ^
: string
| set master = 1, a = 2
< ok 0 0
> exec ^
: string
| set master_log_file = 1, a = 2
< ok 0 0
> exec ^
: string
| set match = 1, a = 2
< error
| 1064 42000 …line 1 column 9 near "match = 1, a = 2" 
> exec ^
: string
| set materialized = 1, a = 2
< ok 0 0
> exec ^
: string
| set max = 1, a = 2
< ok 0 0
> exec ^
: string
| set max_connections_per_hour = 1, a = 2
< ok 0 0
> exec ^
: string
| set max_queries_per_hour = 1, a = 2
< ok 0 0
> exec ^
: string
| set max_rows = 1, a = 2
< ok 0 0
> exec ^
: string
| set max_size = 1, a = 2
< ok 0 0
> exec ^
: string
| set max_updates_per_hour = 1, a = 2
< ok 0 0
> exec ^
: string
| set max_user_connections = 1, a = 2
< ok 0 0
> exec ^
: string
| set maxvalue = 1, a = 2
< error
| 1064 42000 …line 1 column 12 near "maxvalue = 1, a = 2" 
> exec ^
: string
| set medium = 1, a = 2
< ok 0 0
> exec ^
: string
| set mediumblob = 1, a = 2
< error
| 1064 42000 …line 1 column 14 near "mediumblob = 1, a = 2" 
> exec ^
: string
| set mediumint = 1, a = 2
< error
| 1064 42000 …line 1 column 13 near "mediumint = 1, a = 2" 
> exec ^
: string
| set mediumtext = 1, a = 2
< error
| 1064 42000 …line 1 column 14 near "mediumtext = 1, a = 2" 
> exec ^
: string
| set member = 1, a = 2
< ok 0 0
> exec ^
: string
| set memory = 1, a = 2
< ok 0 0
> exec ^
: string
| set merge = 1, a = 2
< ok 0 0
> exec ^
: string
| set message_text = 1, a = 2
< ok 0 0
> exec ^
: string
| set microsecond = 1, a = 2
< ok 0 0
> exec ^
: string
| set mid = 1, a = 2
< ok 0 0
> exec ^
: string
| set middleint = 1, a = 2
< error
| 1064 42000 …line 1 column 13 near "middleint = 1, a = 2" 
> exec ^
: string
| set migrate = 1, a = 2
< ok 0 0
> exec ^
: string
| set min = 1, a = 2
< ok 0 0
> exec ^
: string
| set min_rows = 1, a = 2
< ok 0 0
> exec ^
: string
| set minute = 1, a = 2
< ok 0 0
> exec ^
: string
| set minute_microsecond = 1, a = 2
< error
| 1064 42000 …line 1 column 22 near "minute_microsecond = 1, a = 2" 
> exec ^
: string
| set minute_second = 1, a = 2
< error
| 1064 42000 …line 1 column 17 near "minute_second = 1, a = 2" 
> exec ^
: string
| set minvalue = 1, a = 2
< ok 0 0
> exec ^
: string
| set mod = 1, a = 2
< error
| 1064 42000 …line 1 column 7 near "mod = 1, a = 2" 
> exec ^
: string
| set mode = 1, a = 2
< ok 0 0
> exec ^
: string
| set modifies = 1, a = 2
< ok 0 0
> exec ^
: string
| set modify = 1, a = 2
< ok 0 0
> exec ^
: string
| set month = 1, a = 2
< ok 0 0
> exec ^
: string
| set multilinestring = 1, a = 2
< ok 0 0
> exec ^
: string
| set multipoint = 1, a = 2
< ok 0 0
> exec ^
: string
| set multipolygon = 1, a = 2
< ok 0 0
> exec ^
: string
| set mutex = 1, a = 2
< ok 0 0
> exec ^
: string
| set mysql_errno = 1, a = 2
< ok 0 0
> exec ^
: string
| set name = 1, a = 2
< ok 0 0
> exec ^
: string
| set names = 1, a = 2
< ok 0 0
> exec ^
: string
| set national = 1, a = 2
< ok 0 0
> exec ^
: string
| set natural = 1, a = 2
< error
| 1064 42000 …line 1 column 11 near "natural = 1, a = 2" 
> exec ^
: string
| set nchar = 1, a = 2
< ok 0 0
> exec ^
: string
| set ndb = 1, a = 2
< ok 0 0
> exec ^
: string
| set ndbcluster = 1, a = 2
< ok 0 0
> exec ^
: string
| set nested = 1, a = 2
< ok 0 0
> exec ^
: string
| set network_namespace = 1, a = 2
< ok 0 0
> exec ^
: string
| set never = 1, a = 2
< ok 0 0
> exec ^
: string
| set new = 1, a = 2
< ok 0 0
> exec ^
: string
| set next = 1, a = 2
< ok 0 0
> exec ^
: string
| set next_row_id = 1, a = 2
< ok 0 0
> exec ^
: string
| set no = 1, a = 2
< ok 0 0
> exec ^
: string
| set no_wait = 1, a = 2
< ok 0 0
> exec ^
: string
| set no_write_to_binlog = 1, a = 2
< error
| 1064 42000 …line 1 column 22 near "no_write_to_binlog = 1, a = 2" 
> exec ^
: string
| set nocache = 1, a = 2
< ok 0 0
> exec ^
: string
| set nocycle = 1, a = 2
< ok 0 0
> exec ^
: string
| set nodegroup = 1, a = 2
< ok 0 0
> exec ^
: string
| set nomaxvalue = 1, a = 2
< ok 0 0
> exec ^
: string
| set nominvalue = 1, a = 2
< ok 0 0
> exec ^
: string
| set nonclustered = 1, a = 2
< ok 0 0
> exec ^
: string
| set none = 1, a = 2
< ok 0 0
> exec ^
: string
| set not = 1, a = 2
< error
| 1064 42000 …line 1 column 7 near "not = 1, a = 2" 
> exec ^
: string
| set now = 1, a = 2
< ok 0 0
> exec ^
: string
| set nowait = 1, a = 2
< ok 0 0
> exec ^
: string
| set nth_value = 1, a = 2
< error
| 1064 42000 …line 1 column 13 near "nth_value = 1, a = 2" 
> exec ^
: string
| set ntile = 1, a = 2
< error
| 1064 42000 …line 1 column 9 near "ntile = 1, a = 2" 
> exec ^
: string
| set null = 1, a = 2
< error
| 1064 42000 …line 1 column 8 near "null = 1, a = 2" 
> exec ^
: string
| set nullif = 1, a = 2
< ok 0 0
> exec ^
: string
| set nulls = 1, a = 2
< ok 0 0
> exec ^
: string
| set number = 1, a = 2
< ok 0 0
> exec ^
: string
| set numeric = 1, a = 2
< error
| 1064 42000 …line 1 column 11 near "numeric = 1, a = 2" 
> exec ^
: string
| set nvarchar = 1, a = 2
< ok 0 0
> exec ^
: string
| set of = 1, a = 2
< error
| 1064 42000 …line 1 column 6 near "of = 1, a = 2" 
> exec ^
: string
| set off = 1, a = 2
< ok 0 0
> exec ^
: string
| set offset = 1, a = 2
< ok 0 0
> exec ^
: string
| set oj = 1, a = 2
< ok 0 0
> exec ^
: string
| set old = 1, a = 2
< ok 0 0
> exec ^
: string
| set oltp_read_only = 1, a = 2
< ok 0 0
> exec ^
: string
| set oltp_read_write = 1, a = 2
< ok 0 0
> exec ^
: string
| set oltp_write_only = 1, a = 2
< ok 0 0
> exec ^
: string
| set on = 1, a = 2
< error
| 1064 42000 …line 1 column 6 near "on = 1, a = 2" 
> exec ^
: string
| set one = 1, a = 2
< ok 0 0
> exec ^
: string
| set only = 1, a = 2
< ok 0 0
> exec ^
: string
| set open = 1, a = 2
< ok 0 0
> exec ^
: string
| set optimize = 1, a = 2
< error
| 1064 42000 …line 1 column 12 near "optimize = 1, a = 2" 
> exec ^
: string
| set optimizer_costs = 1, a = 2
< ok 0 0
> exec ^
: string
| set option = 1, a = 2
< error
| 1064 42000 …line 1 column 10 near "option = 1, a = 2" 
> exec ^
: string
| set optional = 1, a = 2
< ok 0 0
> exec ^
: string
| set optionally = 1, a = 2
< error
| 1064 42000 …line 1 column 14 near "optionally = 1, a = 2" 
> exec ^
: string
| set options = 1, a = 2
< ok 0 0
> exec ^
: string
| set or = 1, a = 2
< error
| 1064 42000 …line 1 column 6 near "or = 1, a = 2" 
> exec ^
: string
| set order = 1, a = 2
< error
| 1064 42000 …line 1 column 9 near "order = 1, a = 2" 
> exec ^
: string
| set ordinality = 1, a = 2
< ok 0 0
> exec ^
: string
| set organization = 1, a = 2
< ok 0 0
> exec ^
: string
| set others = 1, a = 2
< ok 0 0
> exec ^
: string
| set out = 1, a = 2
< error
| 1064 42000 …line 1 column 7 near "out = 1, a = 2" 
> exec ^
: string
| set outer = 1, a = 2
< error
| 1064 42000 …line 1 column 9 near "outer = 1, a = 2" 
> exec ^
: string
| set outfile = 1, a = 2
< error
| 1064 42000 …line 1 column 11 near "outfile = 1, a = 2" 
> exec ^
: string
| set over = 1, a = 2
< error
| 1064 42000 …line 1 column 8 near "over = 1, a = 2" 
> exec ^
: string
| set overlay = 1, a = 2
< ok 0 0
> exec ^
: string
| set owner = 1, a = 2
< ok 0 0
> exec ^
: string
| set pack_keys = 1, a = 2
< ok 0 0
> exec ^
: string
| set page = 1, a = 2
< ok 0 0
> exec ^
: string
| set parallel = 1, a = 2
< ok 0 0
> exec ^
: string
| set parameters = 1, a = 2
< ok 0 0
> exec ^
: string
| set parse_tree = 1, a = 2
< ok 0 0
> exec ^
: string
| set parser = 1, a = 2
< ok 0 0
> exec ^
: string
| set partial = 1, a = 2
< ok 0 0
> exec ^
: string
| set partition = 1, a = 2
< error
| 1064 42000 …line 1 column 13 near "partition = 1, a = 2" 
> exec ^
: string
| set partitioning = 1, a = 2
< ok 0 0
> exec ^
: string
| set partitions = 1, a = 2
< ok 0 0
> exec ^
: string
| set password = 1, a = 2
< error
| 1064 42000 …line 1 column 16 near "1, a = 2" 
> exec ^
: string
| set password_lock_time = 1, a = 2
< ok 0 0
> exec ^
: string
| set path = 1, a = 2
< ok 0 0
> exec ^
: string
| set percent_rank = 1, a = 2
< error
| 1064 42000 …line 1 column 16 near "percent_rank = 1, a = 2" 
> exec ^
: string
| set persist = 1, a = 2
< ok 0 0
> exec ^
: string
| set persist_only = 1, a = 2
< ok 0 0
> exec ^
: string
| set phase = 1, a = 2
< ok 0 0
> exec ^
: string
| set placement = 1, a = 2
< ok 0 0
> exec ^
: string
| set placing = 1, a = 2
< ok 0 0
> exec ^
: string
| set plan = 1, a = 2
< ok 0 0
> exec ^
: string
| set plugin = 1, a = 2
< ok 0 0
> exec ^
: string
| set plugin_dir = 1, a = 2
< ok 0 0
> exec ^
: string
| set plugins = 1, a = 2
< ok 0 0
> exec ^
: string
| set point = 1, a = 2
< ok 0 0
> exec ^
: string
| set policy = 1, a = 2
< ok 0 0
> exec ^
: string
| set polygon = 1, a = 2
< ok 0 0
> exec ^
: string
| set port = 1, a = 2
< ok 0 0
> exec ^
: string
| set position = 1, a = 2
< ok 0 0
> exec ^
: string
| set pre_split_regions = 1, a = 2
< ok 0 0
> exec ^
: string
| set precedes = 1, a = 2
< ok 0 0
> exec ^
: string
| set preceding = 1, a = 2
< ok 0 0
> exec ^
: string
| set precision = 1, a = 2
< error
| 1064 42000 …line 1 column 13 near "precision = 1, a = 2" 
> exec ^
: string
| set prepare = 1, a = 2
< ok 0 0
> exec ^
: string
| set preserve = 1, a = 2
< ok 0 0
> exec ^
: string
| set prev = 1, a = 2
< ok 0 0
> exec ^
: string
| set primary = 1, a = 2
< error
| 1064 42000 …line 1 column 11 near "primary = 1, a = 2" 
> exec ^
: string
| set primary_region = 1, a = 2
< ok 0 0
> exec ^
: string
| set priority = 1, a = 2
< ok 0 0
> exec ^
: string
| set privilege_checks_user = 1, a = 2
< ok 0 0
> exec ^
: string
| set privileges = 1, a = 2
< ok 0 0
> exec ^
: string
| set procedure = 1, a = 2
< error
| 1064 42000 …line 1 column 13 near "procedure = 1, a = 2" 
> exec ^
: string
| set process = 1, a = 2
< ok 0 0
> exec ^
: string
| set processed_keys = 1, a = 2
< ok 0 0
> exec ^
: string
| set processlist = 1, a = 2
< ok 0 0
> exec ^
: string
| set profile = 1, a = 2
< ok 0 0
> exec ^
: string
| set profiles = 1, a = 2
< ok 0 0
> exec ^
: string
| set proxy = 1, a = 2
< ok 0 0
> exec ^
: string
| set purge = 1, a = 2
< ok 0 0
> exec ^
: string
| set qualify = 1, a = 2
< ok 0 0
> exec ^
: string
| set quarter = 1, a = 2
< ok 0 0
> exec ^
: string
| set query = 1, a = 2
< ok 0 0
> exec ^
: string
| set query_limit = 1, a = 2
< ok 0 0
> exec ^
: string
| set quick = 1, a = 2
< ok 0 0
> exec ^
: string
| set random = 1, a = 2
< ok 0 0
> exec ^
: string
| set range = 1, a = 2
< error
| 1064 42000 …line 1 column 9 near "range = 1, a = 2" 
> exec ^
: string
| set rank = 1, a = 2
< error
| 1064 42000 …line 1 column 8 near "rank = 1, a = 2" 
> exec ^
: string
| set read = 1, a = 2
< error
| 1064 42000 …line 1 column 8 near "read = 1, a = 2" 
> exec ^
: string
| set read_only = 1, a = 2
< ok 0 0
> exec ^
: string
| set read_write = 1, a = 2
< ok 0 0
> exec ^
: string
| set reads = 1, a = 2
< ok 0 0
> exec ^
: string
| set real = 1, a = 2
< error
| 1064 42000 …line 1 column 8 near "real = 1, a = 2" 
> exec ^
: string
| set rebuild = 1, a = 2
< ok 0 0
> exec ^
: string
| set recover = 1, a = 2
< ok 0 0
> exec ^
: string
| set recursive = 1, a = 2
< error
| 1064 42000 …line 1 column 13 near "recursive = 1, a = 2" 
> exec ^
: string
| set redo_buffer_size = 1, a = 2
< ok 0 0
> exec ^
: string
| set redundant = 1, a = 2
< ok 0 0
> exec ^
: string
| set reference = 1, a = 2
< ok 0 0
> exec ^
: string
| set references = 1, a = 2
< error
| 1064 42000 …line 1 column 14 near "references = 1, a = 2" 
> exec ^
: string
| set regexp = 1, a = 2
< error
| 1064 42000 …line 1 column 10 near "regexp = 1, a = 2" 
> exec ^
: string
| set region = 1, a = 2
< ok 0 0
> exec ^
: string
| set regions = 1, a = 2
< ok 0 0
> exec ^
: string
| set registration = 1, a = 2
< ok 0 0
> exec ^
: string
| set relational = 1, a = 2
< ok 0 0
> exec ^
: string
| set relay = 1, a = 2
< ok 0 0
> exec ^
: string
| set relay_log_file = 1, a = 2
< ok 0 0
> exec ^
: string
| set relay_log_pos = 1, a = 2
< ok 0 0
> exec ^
: string
| set relay_thread = 1, a = 2
< ok 0 0
> exec ^
: string
| set relaylog = 1, a = 2
< ok 0 0
> exec ^
: string
| set release = 1, a = 2
< error
| 1064 42000 …line 1 column 11 near "release = 1, a = 2" 
> exec ^
: string
| set reload = 1, a = 2
< ok 0 0
> exec ^
: string
| set remove = 1, a = 2
< ok 0 0
> exec ^
: string
| set rename = 1, a = 2
< error
| 1064 42000 …line 1 column 10 near "rename = 1, a = 2" 
> exec ^
: string
| set reorganize = 1, a = 2
< ok 0 0
> exec ^
: string
| set repair = 1, a = 2
< ok 0 0
> exec ^
: string
| set repeat = 1, a = 2
< error
| 1064 42000 …line 1 column 10 near "repeat = 1, a = 2" 
> exec ^
: string
| set repeatable = 1, a = 2
< ok 0 0
> exec ^
: string
| set replace = 1, a = 2
< error
| 1064 42000 …line 1 column 11 near "replace = 1, a = 2" 
> exec ^
: string
| set replica = 1, a = 2
< ok 0 0
> exec ^
: string
| set replicas = 1, a = 2
< ok 0 0
> exec ^
: string
| set replicate_do_db = 1, a = 2
< ok 0 0
> exec ^
: string
| set replicate_do_table = 1, a = 2
< ok 0 0
> exec ^
: string
| set replicate_ignore_db = 1, a = 2
< ok 0 0
> exec ^
: string
| set replicate_ignore_table = 1, a = 2
< ok 0 0
> exec ^
: string
| set replicate_rewrite_db = 1, a = 2
< ok 0 0
> exec ^
: string
| set replicate_wild_do_table = 1, a = 2
< ok 0 0
> exec ^
: string
| set replicate_wild_ignore_table = 1, a = 2
< ok 0 0
> exec ^
: string
| set replication = 1, a = 2
< ok 0 0
> exec ^
: string
| set require = 1, a = 2
< error
| 1064 42000 …line 1 column 11 near "require = 1, a = 2" 
> exec ^
: string
| set require_row_format = 1, a = 2
< ok 0 0
> exec ^
: string
| set require_table_primary_key_check = 1, a = 2
< ok 0 0
> exec ^
: string
| set reset = 1, a = 2
< ok 0 0
> exec ^
: string
| set resignal = 1, a = 2
< ok 0 0
> exec ^
: string
| set resource = 1, a = 2
< ok 0 0
> exec ^
: string
| set respect = 1, a = 2
< ok 0 0
> exec ^
: string
| set restart = 1, a = 2
< ok 0 0
> exec ^
: string
| set restore = 1, a = 2
< ok 0 0
> exec ^
: string
| set restores = 1, a = 2
< ok 0 0
> exec ^
: string
| set restrict = 1, a = 2
< error
| 1064 42000 …line 1 column 12 near "restrict = 1, a = 2" 
> exec ^
: string
| set resume = 1, a = 2
< ok 0 0
> exec ^
: string
| set retain = 1, a = 2
< ok 0 0
> exec ^
: string
| set return = 1, a = 2
< ok 0 0
> exec ^
: string
| set returned_sqlstate = 1, a = 2
< ok 0 0
> exec ^
: string
| set returning = 1, a = 2
< ok 0 0
> exec ^
: string
| set returns = 1, a = 2
< ok 0 0
> exec ^
: string
| set reuse = 1, a = 2
< ok 0 0
> exec ^
: string
| set reverse = 1, a = 2
< ok 0 0
> exec ^
: string
| set revoke = 1, a = 2
< error
| 1064 42000 …line 1 column 10 near "revoke = 1, a = 2" 
> exec ^
: string
| set right = 1, a = 2
< error
| 1064 42000 …line 1 column 9 near "right = 1, a = 2" 
> exec ^
: string
| set rlike = 1, a = 2
< error
| 1064 42000 …line 1 column 9 near "rlike = 1, a = 2" 
> exec ^
: string
| set role = 1, a = 2
< ok 0 0
> exec ^
: string
| set rollback = 1, a = 2
< ok 0 0
> exec ^
: string
| set rollup = 1, a = 2
< ok 0 0
> exec ^
: string
| set rotate = 1, a = 2
< ok 0 0
> exec ^
: string
| set routine = 1, a = 2
< ok 0 0
> exec ^
: string
| set row = 1, a = 2
< error
| 1064 42000 …line 1 column 7 near "row = 1, a = 2" 
> exec ^
: string
| set row_count = 1, a = 2
< ok 0 0
> exec ^
: string
| set row_format = 1, a = 2
< ok 0 0
> exec ^
: string
| set row_number = 1, a = 2
< error
| 1064 42000 …line 1 column 14 near "row_number = 1, a = 2" 
> exec ^
: string
| set rows = 1, a = 2
< error
| 1064 42000 …line 1 column 8 near "rows = 1, a = 2" 
> exec ^
: string
| set rtree = 1, a = 2
< ok 0 0
> exec ^
: string
| set ru = 1, a = 2
< ok 0 0
> exec ^
: string
| set ru_per_sec = 1, a = 2
< ok 0 0
> exec ^
: string
| set rule = 1, a = 2
< ok 0 0
> exec ^
: string
| set run = 1, a = 2
< ok 0 0
> exec ^
: string
| set s3 = 1, a = 2
< ok 0 0
> exec ^
: string
| set savepoint = 1, a = 2
< ok 0 0
> exec ^
: string
| set schedule = 1, a = 2
< ok 0 0
> exec ^
: string
| set schema = 1, a = 2
< error
| 1064 42000 …line 1 column 10 near "schema = 1, a = 2" 
> exec ^
: string
| set schema_name = 1, a = 2
< ok 0 0
> exec ^
: string
| set schemas = 1, a = 2
< error
| 1064 42000 …line 1 column 11 near "schemas = 1, a = 2" 
> exec ^
: string
| set second = 1, a = 2
< ok 0 0
> exec ^
: string
| set second_microsecond = 1, a = 2
< error
| 1064 42000 …line 1 column 22 near "second_microsecond = 1, a = 2" 
> exec ^
: string
| set secondary = 1, a = 2
< ok 0 0
> exec ^
: string
| set secondary_engine = 1, a = 2
< ok 0 0
> exec ^
: string
| set secondary_engine_attribute = 1, a = 2
< ok 0 0
> exec ^
: string
| set secondary_load = 1, a = 2
< ok 0 0
> exec ^
: string
| set secondary_unload = 1, a = 2
< ok 0 0
> exec ^
: string
| set security = 1, a = 2
< ok 0 0
> exec ^
: string
| set select = 1, a = 2
< error
| 1064 42000 …line 1 column 10 near "select = 1, a = 2" 
> exec ^
: string
| set sensitive = 1, a = 2
< ok 0 0
> exec ^
: string
| set separator = 1, a = 2
< ok 0 0
> exec ^
: string
| set sequence = 1, a = 2
< ok 0 0
> exec ^
: string
| set serial = 1, a = 2
< ok 0 0
> exec ^
: string
| set serializable = 1, a = 2
< ok 0 0
> exec ^
: string
| set server = 1, a = 2
< ok 0 0
> exec ^
: string
| set session = 1, a = 2
< ok 0 0
> exec ^
: string
| set session_states = 1, a = 2
< ok 0 0
> exec ^
: string
| set session_user = 1, a = 2
< ok 0 0
> exec ^
: string
| set set = 1, a = 2
< error
| 1064 42000 …line 1 column 7 near "set = 1, a = 2" 
> exec ^
: string
| set sets = 1, a = 2
< ok 0 0
> exec ^
: string
| set shard_row_id_bits = 1, a = 2
< ok 0 0
> exec ^
: string
| set share = 1, a = 2
< ok 0 0
> exec ^
: string
| set show = 1, a = 2
< error
| 1064 42000 …line 1 column 8 near "show = 1, a = 2" 
> exec ^
: string
| set shutdown = 1, a = 2
< ok 0 0
> exec ^
: string
| set signal = 1, a = 2
< ok 0 0
> exec ^
: string
| set signed = 1, a = 2
< ok 0 0
> exec ^
: string
| set similar = 1, a = 2
< ok 0 0
> exec ^
: string
| set simple = 1, a = 2
< ok 0 0
> exec ^
: string
| set skip = 1, a = 2
< ok 0 0
> exec ^
: string
| set slave = 1, a = 2
< ok 0 0
> exec ^
: string
| set slow = 1, a = 2
< ok 0 0
> exec ^
: string
| set smallint = 1, a = 2
< error
| 1064 42000 …line 1 column 12 near "smallint = 1, a = 2" 
> exec ^
: string
| set snapshot = 1, a = 2
< ok 0 0
> exec ^
: string
| set socket = 1, a = 2
< ok 0 0
> exec ^
: string
| set some = 1, a = 2
< ok 0 0
> exec ^
: string
| set soname = 1, a = 2
< ok 0 0
> exec ^
: string
| set sounds = 1, a = 2
< ok 0 0
> exec ^
: string
| set source = 1, a = 2
< ok 0 0
> exec ^
: string
| set source_auto_position = 1, a = 2
< ok 0 0
> exec ^
: string
| set source_bind = 1, a = 2
< ok 0 0
> exec ^
: string
| set source_compression_algorithms = 1, a = 2
< ok 0 0
> exec ^
: string
| set source_connect_retry = 1, a = 2
< ok 0 0
> exec ^
: string
| set source_connection_auto_failover = 1, a = 2
< ok 0 0
> exec ^
: string
| set source_delay = 1, a = 2
< ok 0 0
> exec ^
: string
| set source_heartbeat_period = 1, a = 2
< ok 0 0
> exec ^
: string
| set source_host = 1, a = 2
< ok 0 0
> exec ^
: string
| set source_log_file = 1, a = 2
< ok 0 0
> exec ^
: string
| set source_log_pos = 1, a = 2
< ok 0 0
> exec ^
: string
| set source_password = 1, a = 2
< ok 0 0
> exec ^
: string
| set source_port = 1, a = 2
< ok 0 0
> exec ^
: string
| set source_public_key_path = 1, a = 2
< ok 0 0
> exec ^
: string
| set source_retry_count = 1, a = 2
< ok 0 0
> exec ^
: string
| set source_ssl = 1, a = 2
< ok 0 0
> exec ^
: string
| set source_ssl_ca = 1, a = 2
< ok 0 0
> exec ^
: string
| set source_ssl_capath = 1, a = 2
< ok 0 0
> exec ^
: string
| set source_ssl_cert = 1, a = 2
< ok 0 0
> exec ^
: string
| set source_ssl_cipher = 1, a = 2
< ok 0 0
> exec ^
: string
| set source_ssl_crl = 1, a = 2
< ok 0 0
> exec ^
: string
| set source_ssl_crlpath = 1, a = 2
< ok 0 0
> exec ^
: string
| set source_ssl_key = 1, a = 2
< ok 0 0
> exec ^
: string
| set source_ssl_verify_server_cert = 1, a = 2
< ok 0 0
> exec ^
: string
| set source_tls_ciphersuites = 1, a = 2
< ok 0 0
> exec ^
: string
| set source_tls_version = 1, a = 2
< ok 0 0
> exec ^
: string
| set source_user = 1, a = 2
< ok 0 0
> exec ^
: string
| set source_zstd_compression_level = 1, a = 2
< ok 0 0
> exec ^
: string
| set spatial = 1, a = 2
< error
| 1064 42000 …line 1 column 11 near "spatial = 1, a = 2" 
> exec ^
: string
| set specific = 1, a = 2
< ok 0 0
> exec ^
: string
| set split = 1, a = 2
< ok 0 0
> exec ^
: string
| set sql = 1, a = 2
< error
| 1064 42000 …line 1 column 7 near "sql = 1, a = 2" 
> exec ^
: string
| set sql_after_gtids = 1, a = 2
< ok 0 0
> exec ^
: string
| set sql_after_mts_gaps = 1, a = 2
< ok 0 0
> exec ^
: string
| set sql_before_gtids = 1, a = 2
< ok 0 0
> exec ^
: string
| set sql_big_result = 1, a = 2
< error
| 1064 42000 …line 1 column 18 near "sql_big_result = 1, a = 2" 
> exec ^
: string
| set sql_buffer_result = 1, a = 2
< ok 0 0
> exec ^
: string
| set sql_calc_found_rows = 1, a = 2
< error
| 1064 42000 …line 1 column 23 near "sql_calc_found_rows = 1, a = 2" 
> exec ^
: string
| set sql_no_cache = 1, a = 2
< ok 0 0
> exec ^
: string
| set sql_small_result = 1, a = 2
< error
| 1064 42000 …line 1 column 20 near "sql_small_result = 1, a = 2" 
> exec ^
: string
| set sql_thread = 1, a = 2
< ok 0 0
> exec ^
: string
| set sql_tsi_day = 1, a = 2
< ok 0 0
> exec ^
: string
| set sql_tsi_hour = 1, a = 2
< ok 0 0
> exec ^
: string
| set sql_tsi_minute = 1, a = 2
< ok 0 0
> exec ^
: string
| set sql_tsi_month = 1, a = 2
< ok 0 0
> exec ^
: string
| set sql_tsi_quarter = 1, a = 2
< ok 0 0
> exec ^
: string
| set sql_tsi_second = 1, a = 2
< ok 0 0
> exec ^
: string
| set sql_tsi_week = 1, a = 2
< ok 0 0
> exec ^
: string
| set sql_tsi_year = 1, a = 2
< ok 0 0
> exec ^
: string
| set sqlexception = 1, a = 2
< error
| 1064 42000 …line 1 column 16 near "sqlexception = 1, a = 2" 
> exec ^
: string
| set sqlstate = 1, a = 2
< error
| 1064 42000 …line 1 column 12 near "sqlstate = 1, a = 2" 
> exec ^
: string
| set sqlwarning = 1, a = 2
< error
| 1064 42000 …line 1 column 14 near "sqlwarning = 1, a = 2" 
> exec ^
: string
| set srid = 1, a = 2
< ok 0 0
> exec ^
: string
| set ssl = 1, a = 2
< error
| 1064 42000 …line 1 column 7 near "ssl = 1, a = 2" 
> exec ^
: string
| set st_collect = 1, a = 2
< ok 0 0
> exec ^
: string
| set stacked = 1, a = 2
< ok 0 0
> exec ^
: string
| set start = 1, a = 2
< ok 0 0
> exec ^
: string
| set start_time = 1, a = 2
< ok 0 0
> exec ^
: string
| set starting = 1, a = 2
< error
| 1064 42000 …line 1 column 12 near "starting = 1, a = 2" 
> exec ^
: string
| set starts = 1, a = 2
< ok 0 0
> exec ^
: string
| set stats_auto_recalc = 1, a = 2
< ok 0 0
> exec ^
: string
| set stats_buckets = 1, a = 2
< ok 0 0
> exec ^
: string
| set stats_col_choice = 1, a = 2
< ok 0 0
> exec ^
: string
| set stats_col_list = 1, a = 2
< ok 0 0
> exec ^
: string
| set stats_extended = 1, a = 2
< ok 0 0
> exec ^
: string
| set stats_healthy = 1, a = 2
< ok 0 0
> exec ^
: string
| set stats_histograms = 1, a = 2
< ok 0 0
> exec ^
: string
| set stats_locked = 1, a = 2
< ok 0 0
> exec ^
: string
| set stats_meta = 1, a = 2
< ok 0 0
> exec ^
: string
| set stats_persistent = 1, a = 2
< ok 0 0
> exec ^
: string
| set stats_sample_pages = 1, a = 2
< ok 0 0
> exec ^
: string
| set stats_sample_rate = 1, a = 2
< ok 0 0
> exec ^
: string
| set stats_topn = 1, a = 2
< ok 0 0
> exec ^
: string
| set status = 1, a = 2
< ok 0 0
> exec ^
: string
| set std = 1, a = 2
< ok 0 0
> exec ^
: string
| set stddev = 1, a = 2
< ok 0 0
> exec ^
: string
| set stddev_pop = 1, a = 2
< ok 0 0
> exec ^
: string
| set stddev_samp = 1, a = 2
< ok 0 0
> exec ^
: string
| set stop = 1, a = 2
< ok 0 0
> exec ^
: string
| set storage = 1, a = 2
< ok 0 0
> exec ^
: string
| set stored = 1, a = 2
< error
| 1064 42000 …line 1 column 10 near "stored = 1, a = 2" 
> exec ^
: string
| set straight_join = 1, a = 2
< error
| 1064 42000 …line 1 column 17 near "straight_join = 1, a = 2" 
> exec ^
: string
| set stream = 1, a = 2
< ok 0 0
> exec ^
: string
| set strict_load = 1, a = 2
< ok 0 0
> exec ^
: string
| set string = 1, a = 2
< ok 0 0
> exec ^
: string
| set subclass_origin = 1, a = 2
< ok 0 0
> exec ^
: string
| set subdate = 1, a = 2
< ok 0 0
> exec ^
: string
| set subject = 1, a = 2
< ok 0 0
> exec ^
: string
| set subpartition = 1, a = 2
< ok 0 0
> exec ^
: string
| set subpartitions = 1, a = 2
< ok 0 0
> exec ^
: string
| set substr = 1, a = 2
< ok 0 0
> exec ^
: string
| set substring = 1, a = 2
< ok 0 0
> exec ^
: string
| set sum = 1, a = 2
< ok 0 0
> exec ^
: string
| set super = 1, a = 2
< ok 0 0
> exec ^
: string
| set survival_preferences = 1, a = 2
< ok 0 0
> exec ^
: string
| set suspend = 1, a = 2
< ok 0 0
> exec ^
: string
| set swaps = 1, a = 2
< ok 0 0
> exec ^
: string
| set switch_group = 1, a = 2
< ok 0 0
> exec ^
: string
| set switches = 1, a = 2
< ok 0 0
> exec ^
: string
| set sysdate = 1, a = 2
< ok 0 0
> exec ^
: string
| set system = 1, a = 2
< ok 0 0
> exec ^
: string
| set system_user = 1, a = 2
< ok 0 0
> exec ^
: string
| set table = 1, a = 2
< error
| 1064 42000 …line 1 column 9 near "table = 1, a = 2" 
> exec ^
: string
| set table_checksum = 1, a = 2
< ok 0 0
> exec ^
: string
| set table_name = 1, a = 2
< ok 0 0
> exec ^
: string
| set tables = 1, a = 2
< ok 0 0
> exec ^
: string
| set tablesample = 1, a = 2
< error
| 1064 42000 …line 1 column 15 near "tablesample = 1, a = 2" 
> exec ^
: string
| set tablespace = 1, a = 2
< ok 0 0
> exec ^
: string
| set task_types = 1, a = 2
< ok 0 0
> exec ^
: string
| set temporary = 1, a = 2
< ok 0 0
> exec ^
: string
| set temptable = 1, a = 2
< ok 0 0
> exec ^
: string
| set terminated = 1, a = 2
< error
| 1064 42000 …line 1 column 14 near "terminated = 1, a = 2" 
> exec ^
: string
| set text = 1, a = 2
< ok 0 0
> exec ^
: string
| set than = 1, a = 2
< ok 0 0
> exec ^
: string
| set then = 1, a = 2
< error
| 1064 42000 …line 1 column 8 near "then = 1, a = 2" 
> exec ^
: string
| set thread_priority = 1, a = 2
< ok 0 0
> exec ^
: string
| set ties = 1, a = 2
< ok 0 0
> exec ^
: string
| set tiflash = 1, a = 2
< ok 0 0
> exec ^
: string
| set time = 1, a = 2
< ok 0 0
> exec ^
: string
| set timeout = 1, a = 2
< ok 0 0
> exec ^
: string
| set timestamp = 1, a = 2
< ok 0 0
> exec ^
: string
| set timestampadd = 1, a = 2
< ok 0 0
> exec ^
: string
| set timestampdiff = 1, a = 2
< ok 0 0
> exec ^
: string
| set tinyblob = 1, a = 2
< error
| 1064 42000 …line 1 column 12 near "tinyblob = 1, a = 2" 
> exec ^
: string
| set tinyint = 1, a = 2
< error
| 1064 42000 …line 1 column 11 near "tinyint = 1, a = 2" 
> exec ^
: string
| set tinytext = 1, a = 2
< error
| 1064 42000 …line 1 column 12 near "tinytext = 1, a = 2" 
> exec ^
: string
| set tls = 1, a = 2
< ok 0 0
> exec ^
: string
| set to = 1, a = 2
< error
| 1064 42000 …line 1 column 6 near "to = 1, a = 2" 
> exec ^
: string
| set tpcc = 1, a = 2
< ok 0 0
> exec ^
: string
| set tpch_10 = 1, a = 2
< ok 0 0
> exec ^
: string
| set trailing = 1, a = 2
< error
| 1064 42000 …line 1 column 12 near "trailing = 1, a = 2" 
> exec ^
: string
| set transaction = 1, a = 2
< ok 0 0
> exec ^
: string
| set trigger = 1, a = 2
< error
| 1064 42000 …line 1 column 11 near "trigger = 1, a = 2" 
> exec ^
: string
| set triggers = 1, a = 2
< ok 0 0
> exec ^
: string
| set trim = 1, a = 2
< ok 0 0
> exec ^
: string
| set true = 1, a = 2
< error
| 1064 42000 …line 1 column 8 near "true = 1, a = 2" 
> exec ^
: string
| set truncate = 1, a = 2
< ok 0 0
> exec ^
: string
| set ttl = 1, a = 2
< ok 0 0
> exec ^
: string
| set ttl_enable = 1, a = 2
< ok 0 0
> exec ^
: string
| set ttl_job_interval = 1, a = 2
< ok 0 0
> exec ^
: string
| set type = 1, a = 2
< ok 0 0
> exec ^
: string
| set types = 1, a = 2
< ok 0 0
> exec ^
: string
| set unbounded = 1, a = 2
< ok 0 0
> exec ^
: string
| set uncommitted = 1, a = 2
< ok 0 0
> exec ^
: string
| set undefined = 1, a = 2
< ok 0 0
> exec ^
: string
| set undo = 1, a = 2
< ok 0 0
> exec ^
: string
| set undo_buffer_size = 1, a = 2
< ok 0 0
> exec ^
: string
| set undofile = 1, a = 2
< ok 0 0
> exec ^
: string
| set unicode = 1, a = 2
< ok 0 0
> exec ^
: string
| set uninstall = 1, a = 2
< ok 0 0
> exec ^
: string
| set union = 1, a = 2
< error
| 1064 42000 …line 1 column 9 near "union = 1, a = 2" 
> exec ^
: string
| set unique = 1, a = 2
< error
| 1064 42000 …line 1 column 10 near "unique = 1, a = 2" 
> exec ^
: string
| set unknown = 1, a = 2
< ok 0 0
> exec ^
: string
| set unlimited = 1, a = 2
< ok 0 0
> exec ^
: string
| set unlock = 1, a = 2
< error
| 1064 42000 …line 1 column 10 near "unlock = 1, a = 2" 
> exec ^
: string
| set unregister = 1, a = 2
< ok 0 0
> exec ^
: string
| set unsigned = 1, a = 2
< error
| 1064 42000 …line 1 column 12 near "unsigned = 1, a = 2" 
> exec ^
: string
| set until = 1, a = 2
< error
| 1064 42000 …line 1 column 9 near "until = 1, a = 2" 
> exec ^
: string
| set update = 1, a = 2
< error
| 1064 42000 …line 1 column 10 near "update = 1, a = 2" 
> exec ^
: string
| set upgrade = 1, a = 2
< ok 0 0
> exec ^
: string
| set uri = 1, a = 2
< ok 0 0
> exec ^
: string
| set url = 1, a = 2
< ok 0 0
> exec ^
: string
| set usage = 1, a = 2
< error
| 1064 42000 …line 1 column 9 near "usage = 1, a = 2" 
> exec ^
: string
| set use = 1, a = 2
< error
| 1064 42000 …line 1 column 7 near "use = 1, a = 2" 
> exec ^
: string
| set use_frm = 1, a = 2
< ok 0 0
> exec ^
: string
| set user = 1, a = 2
< ok 0 0
> exec ^
: string
| set user_resources = 1, a = 2
< ok 0 0
> exec ^
: string
| set using = 1, a = 2
< error
| 1064 42000 …line 1 column 9 near "using = 1, a = 2" 
> exec ^
: string
| set utc_date = 1, a = 2
< error
| 1064 42000 …line 1 column 12 near "utc_date = 1, a = 2" 
> exec ^
: string
| set utc_time = 1, a = 2
< error
| 1064 42000 …line 1 column 12 near "utc_time = 1, a = 2" 
> exec ^
: string
| set utc_timestamp = 1, a = 2
< error
| 1064 42000 …line 1 column 17 near "utc_timestamp = 1, a = 2" 
> exec ^
: string
| set utilization_limit = 1, a = 2
< ok 0 0
> exec ^
: string
| set validate = 1, a = 2
< ok 0 0
> exec ^
: string
| set validation = 1, a = 2
< ok 0 0
> exec ^
: string
| set value = 1, a = 2
< ok 0 0
> exec ^
: string
| set values = 1, a = 2
< error
| 1064 42000 …line 1 column 10 near "values = 1, a = 2" 
> exec ^
: string
| set var_pop = 1, a = 2
< ok 0 0
> exec ^
: string
| set var_samp = 1, a = 2
< ok 0 0
> exec ^
: string
| set varbinary = 1, a = 2
< error
| 1064 42000 …line 1 column 13 near "varbinary = 1, a = 2" 
> exec ^
: string
| set varchar = 1, a = 2
< error
| 1064 42000 …line 1 column 11 near "varchar = 1, a = 2" 
> exec ^
: string
| set varcharacter = 1, a = 2
< error
| 1064 42000 …line 1 column 16 near "varcharacter = 1, a = 2" 
> exec ^
: string
| set variables = 1, a = 2
< ok 0 0
> exec ^
: string
| set variance = 1, a = 2
< ok 0 0
> exec ^
: string
| set varying = 1, a = 2
< error
| 1064 42000 …line 1 column 11 near "varying = 1, a = 2" 
> exec ^
: string
| set vcpu = 1, a = 2
< ok 0 0
> exec ^
: string
| set vector = 1, a = 2
< ok 0 0
> exec ^
: string
| set verify_key_constraints = 1, a = 2
< ok 0 0
> exec ^
: string
| set view = 1, a = 2
< ok 0 0
> exec ^
: string
| set virtual = 1, a = 2
< error
| 1064 42000 …line 1 column 11 near "virtual = 1, a = 2" 
> exec ^
: string
| set visible = 1, a = 2
< ok 0 0
> exec ^
: string
| set voter_constraints = 1, a = 2
< ok 0 0
> exec ^
: string
| set voters = 1, a = 2
< ok 0 0
> exec ^
: string
| set wait = 1, a = 2
< ok 0 0
> exec ^
: string
| set warnings = 1, a = 2
< ok 0 0
> exec ^
: string
| set watch = 1, a = 2
< ok 0 0
> exec ^
: string
| set week = 1, a = 2
< ok 0 0
> exec ^
: string
| set weight_string = 1, a = 2
< ok 0 0
> exec ^
: string
| set when = 1, a = 2
< error
| 1064 42000 …line 1 column 8 near "when = 1, a = 2" 
> exec ^
: string
| set where = 1, a = 2
< error
| 1064 42000 …line 1 column 9 near "where = 1, a = 2" 
> exec ^
: string
| set while = 1, a = 2
< error
| 1064 42000 …line 1 column 9 near "while = 1, a = 2" 
> exec ^
: string
| set window = 1, a = 2
< error
| 1064 42000 …line 1 column 10 near "window = 1, a = 2" 
> exec ^
: string
| set with = 1, a = 2
< error
| 1064 42000 …line 1 column 8 near "with = 1, a = 2" 
> exec ^
: string
| set without = 1, a = 2
< ok 0 0
> exec ^
: string
| set work = 1, a = 2
< ok 0 0
> exec ^
: string
| set workload = 1, a = 2
< ok 0 0
> exec ^
: string
| set wrapper = 1, a = 2
< ok 0 0
> exec ^
: string
| set write = 1, a = 2
< error
| 1064 42000 …line 1 column 9 near "write = 1, a = 2" 
> exec ^
: string
| set x509 = 1, a = 2
< ok 0 0
> exec ^
: string
| set xa = 1, a = 2
< ok 0 0
> exec ^
: string
| set xid = 1, a = 2
< ok 0 0
> exec ^
: string
| set xml = 1, a = 2
< ok 0 0
> exec ^
: string
| set xor = 1, a = 2
< error
| 1064 42000 …line 1 column 7 near "xor = 1, a = 2" 
> exec ^
: string
| set year = 1, a = 2
< ok 0 0
> exec ^
: string
| set year_month = 1, a = 2
< error
| 1064 42000 …line 1 column 14 near "year_month = 1, a = 2" 
> exec ^
: string
| set zerofill = 1, a = 2
< error
| 1064 42000 …line 1 column 12 near "zerofill = 1, a = 2" 
> exec ^
: string
| set zone = 1, a = 2
< ok 0 0
