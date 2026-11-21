// Redis command-line call-shape grammar as ANTLR 4, generated from the
// official Redis 7.4.12 COMMAND table by redis/parser/cmd/commandgrammar.
// Do not edit by hand: rerun `go run ./redis/parser/cmd/commandgrammar`.
//
// Scope (§2.2 of docs/specs/2026-11-11-engine-compatibility-without-servers-design.md):
// covers every Redis command name, container subcommand, exact/minimum arity
// call shape, and shell-like command-line argument token form.
//
// pin: default redis:7.4.12-alpine sha256:7f7d30b7fdd09467bac58198a1f166dd69789d348d4fe859cccb20e0cf287fbe 7.4.12

grammar RedisCommands;

commandLine
    : command EOF
    ;

command
    : standaloneCommand
    | aclCommand
    | clientCommand
    | clusterCommand
    | commandCommand
    | configCommand
    | functionCommand
    | latencyCommand
    | memoryCommand
    | moduleCommand
    | objectCommand
    | pubsubCommand
    | scriptCommand
    | slowlogCommand
    | xgroupCommand
    | xinfoCommand
    ;

// Top-level Redis commands without subcommands.
standaloneCommand
    : 'APPEND' arg arg // COMMAND append (arity 3)
    | 'ASKING' // COMMAND asking (arity 1)
    | 'AUTH' arg variadicArgs // COMMAND auth (arity -2)
    | 'BGREWRITEAOF' // COMMAND bgrewriteaof (arity 1)
    | 'BGSAVE' variadicArgs // COMMAND bgsave (arity -1)
    | 'BITCOUNT' arg variadicArgs // COMMAND bitcount (arity -2)
    | 'BITFIELD' arg variadicArgs // COMMAND bitfield (arity -2)
    | 'BITFIELD_RO' arg variadicArgs // COMMAND bitfield_ro (arity -2)
    | 'BITOP' arg arg arg variadicArgs // COMMAND bitop (arity -4)
    | 'BITPOS' arg arg variadicArgs // COMMAND bitpos (arity -3)
    | 'BLMOVE' arg arg arg arg arg // COMMAND blmove (arity 6)
    | 'BLMPOP' arg arg arg arg variadicArgs // COMMAND blmpop (arity -5)
    | 'BLPOP' arg arg variadicArgs // COMMAND blpop (arity -3)
    | 'BRPOP' arg arg variadicArgs // COMMAND brpop (arity -3)
    | 'BRPOPLPUSH' arg arg arg // COMMAND brpoplpush (arity 4)
    | 'BZMPOP' arg arg arg arg variadicArgs // COMMAND bzmpop (arity -5)
    | 'BZPOPMAX' arg arg variadicArgs // COMMAND bzpopmax (arity -3)
    | 'BZPOPMIN' arg arg variadicArgs // COMMAND bzpopmin (arity -3)
    | 'COPY' arg arg variadicArgs // COMMAND copy (arity -3)
    | 'DBSIZE' // COMMAND dbsize (arity 1)
    | 'DEBUG' arg variadicArgs // COMMAND debug (arity -2)
    | 'DECR' arg // COMMAND decr (arity 2)
    | 'DECRBY' arg arg // COMMAND decrby (arity 3)
    | 'DEL' arg variadicArgs // COMMAND del (arity -2)
    | 'DISCARD' // COMMAND discard (arity 1)
    | 'DUMP' arg // COMMAND dump (arity 2)
    | 'ECHO' arg // COMMAND echo (arity 2)
    | 'EVAL' arg arg variadicArgs // COMMAND eval (arity -3)
    | 'EVAL_RO' arg arg variadicArgs // COMMAND eval_ro (arity -3)
    | 'EVALSHA' arg arg variadicArgs // COMMAND evalsha (arity -3)
    | 'EVALSHA_RO' arg arg variadicArgs // COMMAND evalsha_ro (arity -3)
    | 'EXEC' // COMMAND exec (arity 1)
    | 'EXISTS' arg variadicArgs // COMMAND exists (arity -2)
    | 'EXPIRE' arg arg variadicArgs // COMMAND expire (arity -3)
    | 'EXPIREAT' arg arg variadicArgs // COMMAND expireat (arity -3)
    | 'EXPIRETIME' arg // COMMAND expiretime (arity 2)
    | 'FAILOVER' variadicArgs // COMMAND failover (arity -1)
    | 'FCALL' arg arg variadicArgs // COMMAND fcall (arity -3)
    | 'FCALL_RO' arg arg variadicArgs // COMMAND fcall_ro (arity -3)
    | 'FLUSHALL' variadicArgs // COMMAND flushall (arity -1)
    | 'FLUSHDB' variadicArgs // COMMAND flushdb (arity -1)
    | 'GEOADD' arg arg arg arg variadicArgs // COMMAND geoadd (arity -5)
    | 'GEODIST' arg arg arg variadicArgs // COMMAND geodist (arity -4)
    | 'GEOHASH' arg variadicArgs // COMMAND geohash (arity -2)
    | 'GEOPOS' arg variadicArgs // COMMAND geopos (arity -2)
    | 'GEORADIUS' arg arg arg arg arg variadicArgs // COMMAND georadius (arity -6)
    | 'GEORADIUS_RO' arg arg arg arg arg variadicArgs // COMMAND georadius_ro (arity -6)
    | 'GEORADIUSBYMEMBER' arg arg arg arg variadicArgs // COMMAND georadiusbymember (arity -5)
    | 'GEORADIUSBYMEMBER_RO' arg arg arg arg variadicArgs // COMMAND georadiusbymember_ro (arity -5)
    | 'GEOSEARCH' arg arg arg arg arg arg variadicArgs // COMMAND geosearch (arity -7)
    | 'GEOSEARCHSTORE' arg arg arg arg arg arg arg variadicArgs // COMMAND geosearchstore (arity -8)
    | 'GET' arg // COMMAND get (arity 2)
    | 'GETBIT' arg arg // COMMAND getbit (arity 3)
    | 'GETDEL' arg // COMMAND getdel (arity 2)
    | 'GETEX' arg variadicArgs // COMMAND getex (arity -2)
    | 'GETRANGE' arg arg arg // COMMAND getrange (arity 4)
    | 'GETSET' arg arg // COMMAND getset (arity 3)
    | 'HDEL' arg arg variadicArgs // COMMAND hdel (arity -3)
    | 'HELLO' variadicArgs // COMMAND hello (arity -1)
    | 'HEXISTS' arg arg // COMMAND hexists (arity 3)
    | 'HEXPIRE' arg arg arg arg arg variadicArgs // COMMAND hexpire (arity -6)
    | 'HEXPIREAT' arg arg arg arg arg variadicArgs // COMMAND hexpireat (arity -6)
    | 'HEXPIRETIME' arg arg arg arg variadicArgs // COMMAND hexpiretime (arity -5)
    | 'HGET' arg arg // COMMAND hget (arity 3)
    | 'HGETALL' arg // COMMAND hgetall (arity 2)
    | 'HINCRBY' arg arg arg // COMMAND hincrby (arity 4)
    | 'HINCRBYFLOAT' arg arg arg // COMMAND hincrbyfloat (arity 4)
    | 'HKEYS' arg // COMMAND hkeys (arity 2)
    | 'HLEN' arg // COMMAND hlen (arity 2)
    | 'HMGET' arg arg variadicArgs // COMMAND hmget (arity -3)
    | 'HMSET' arg arg arg variadicArgs // COMMAND hmset (arity -4)
    | 'HPERSIST' arg arg arg arg variadicArgs // COMMAND hpersist (arity -5)
    | 'HPEXPIRE' arg arg arg arg arg variadicArgs // COMMAND hpexpire (arity -6)
    | 'HPEXPIREAT' arg arg arg arg arg variadicArgs // COMMAND hpexpireat (arity -6)
    | 'HPEXPIRETIME' arg arg arg arg variadicArgs // COMMAND hpexpiretime (arity -5)
    | 'HPTTL' arg arg arg arg variadicArgs // COMMAND hpttl (arity -5)
    | 'HRANDFIELD' arg variadicArgs // COMMAND hrandfield (arity -2)
    | 'HSCAN' arg arg variadicArgs // COMMAND hscan (arity -3)
    | 'HSET' arg arg arg variadicArgs // COMMAND hset (arity -4)
    | 'HSETNX' arg arg arg // COMMAND hsetnx (arity 4)
    | 'HSTRLEN' arg arg // COMMAND hstrlen (arity 3)
    | 'HTTL' arg arg arg arg variadicArgs // COMMAND httl (arity -5)
    | 'HVALS' arg // COMMAND hvals (arity 2)
    | 'INCR' arg // COMMAND incr (arity 2)
    | 'INCRBY' arg arg // COMMAND incrby (arity 3)
    | 'INCRBYFLOAT' arg arg // COMMAND incrbyfloat (arity 3)
    | 'INFO' variadicArgs // COMMAND info (arity -1)
    | 'KEYS' arg // COMMAND keys (arity 2)
    | 'LASTSAVE' // COMMAND lastsave (arity 1)
    | 'LCS' arg arg variadicArgs // COMMAND lcs (arity -3)
    | 'LINDEX' arg arg // COMMAND lindex (arity 3)
    | 'LINSERT' arg arg arg arg // COMMAND linsert (arity 5)
    | 'LLEN' arg // COMMAND llen (arity 2)
    | 'LMOVE' arg arg arg arg // COMMAND lmove (arity 5)
    | 'LMPOP' arg arg arg variadicArgs // COMMAND lmpop (arity -4)
    | 'LOLWUT' variadicArgs // COMMAND lolwut (arity -1)
    | 'LPOP' arg variadicArgs // COMMAND lpop (arity -2)
    | 'LPOS' arg arg variadicArgs // COMMAND lpos (arity -3)
    | 'LPUSH' arg arg variadicArgs // COMMAND lpush (arity -3)
    | 'LPUSHX' arg arg variadicArgs // COMMAND lpushx (arity -3)
    | 'LRANGE' arg arg arg // COMMAND lrange (arity 4)
    | 'LREM' arg arg arg // COMMAND lrem (arity 4)
    | 'LSET' arg arg arg // COMMAND lset (arity 4)
    | 'LTRIM' arg arg arg // COMMAND ltrim (arity 4)
    | 'MGET' arg variadicArgs // COMMAND mget (arity -2)
    | 'MIGRATE' arg arg arg arg arg variadicArgs // COMMAND migrate (arity -6)
    | 'MONITOR' // COMMAND monitor (arity 1)
    | 'MOVE' arg arg // COMMAND move (arity 3)
    | 'MSET' arg arg variadicArgs // COMMAND mset (arity -3)
    | 'MSETNX' arg arg variadicArgs // COMMAND msetnx (arity -3)
    | 'MULTI' // COMMAND multi (arity 1)
    | 'PERSIST' arg // COMMAND persist (arity 2)
    | 'PEXPIRE' arg arg variadicArgs // COMMAND pexpire (arity -3)
    | 'PEXPIREAT' arg arg variadicArgs // COMMAND pexpireat (arity -3)
    | 'PEXPIRETIME' arg // COMMAND pexpiretime (arity 2)
    | 'PFADD' arg variadicArgs // COMMAND pfadd (arity -2)
    | 'PFCOUNT' arg variadicArgs // COMMAND pfcount (arity -2)
    | 'PFDEBUG' arg arg // COMMAND pfdebug (arity 3)
    | 'PFMERGE' arg variadicArgs // COMMAND pfmerge (arity -2)
    | 'PFSELFTEST' // COMMAND pfselftest (arity 1)
    | 'PING' variadicArgs // COMMAND ping (arity -1)
    | 'PSETEX' arg arg arg // COMMAND psetex (arity 4)
    | 'PSUBSCRIBE' arg variadicArgs // COMMAND psubscribe (arity -2)
    | 'PSYNC' arg arg variadicArgs // COMMAND psync (arity -3)
    | 'PTTL' arg // COMMAND pttl (arity 2)
    | 'PUBLISH' arg arg // COMMAND publish (arity 3)
    | 'PUNSUBSCRIBE' variadicArgs // COMMAND punsubscribe (arity -1)
    | 'QUIT' variadicArgs // COMMAND quit (arity -1)
    | 'RANDOMKEY' // COMMAND randomkey (arity 1)
    | 'READONLY' // COMMAND readonly (arity 1)
    | 'READWRITE' // COMMAND readwrite (arity 1)
    | 'RENAME' arg arg // COMMAND rename (arity 3)
    | 'RENAMENX' arg arg // COMMAND renamenx (arity 3)
    | 'REPLCONF' variadicArgs // COMMAND replconf (arity -1)
    | 'REPLICAOF' arg arg // COMMAND replicaof (arity 3)
    | 'RESET' // COMMAND reset (arity 1)
    | 'RESTORE' arg arg arg variadicArgs // COMMAND restore (arity -4)
    | 'RESTORE-ASKING' arg arg arg variadicArgs // COMMAND restore-asking (arity -4)
    | 'ROLE' // COMMAND role (arity 1)
    | 'RPOP' arg variadicArgs // COMMAND rpop (arity -2)
    | 'RPOPLPUSH' arg arg // COMMAND rpoplpush (arity 3)
    | 'RPUSH' arg arg variadicArgs // COMMAND rpush (arity -3)
    | 'RPUSHX' arg arg variadicArgs // COMMAND rpushx (arity -3)
    | 'SADD' arg arg variadicArgs // COMMAND sadd (arity -3)
    | 'SAVE' // COMMAND save (arity 1)
    | 'SCAN' arg variadicArgs // COMMAND scan (arity -2)
    | 'SCARD' arg // COMMAND scard (arity 2)
    | 'SDIFF' arg variadicArgs // COMMAND sdiff (arity -2)
    | 'SDIFFSTORE' arg arg variadicArgs // COMMAND sdiffstore (arity -3)
    | 'SELECT' arg // COMMAND select (arity 2)
    | 'SET' arg arg variadicArgs // COMMAND set (arity -3)
    | 'SETBIT' arg arg arg // COMMAND setbit (arity 4)
    | 'SETEX' arg arg arg // COMMAND setex (arity 4)
    | 'SETNX' arg arg // COMMAND setnx (arity 3)
    | 'SETRANGE' arg arg arg // COMMAND setrange (arity 4)
    | 'SHUTDOWN' variadicArgs // COMMAND shutdown (arity -1)
    | 'SINTER' arg variadicArgs // COMMAND sinter (arity -2)
    | 'SINTERCARD' arg arg variadicArgs // COMMAND sintercard (arity -3)
    | 'SINTERSTORE' arg arg variadicArgs // COMMAND sinterstore (arity -3)
    | 'SISMEMBER' arg arg // COMMAND sismember (arity 3)
    | 'SLAVEOF' arg arg // COMMAND slaveof (arity 3)
    | 'SMEMBERS' arg // COMMAND smembers (arity 2)
    | 'SMISMEMBER' arg arg variadicArgs // COMMAND smismember (arity -3)
    | 'SMOVE' arg arg arg // COMMAND smove (arity 4)
    | 'SORT' arg variadicArgs // COMMAND sort (arity -2)
    | 'SORT_RO' arg variadicArgs // COMMAND sort_ro (arity -2)
    | 'SPOP' arg variadicArgs // COMMAND spop (arity -2)
    | 'SPUBLISH' arg arg // COMMAND spublish (arity 3)
    | 'SRANDMEMBER' arg variadicArgs // COMMAND srandmember (arity -2)
    | 'SREM' arg arg variadicArgs // COMMAND srem (arity -3)
    | 'SSCAN' arg arg variadicArgs // COMMAND sscan (arity -3)
    | 'SSUBSCRIBE' arg variadicArgs // COMMAND ssubscribe (arity -2)
    | 'STRLEN' arg // COMMAND strlen (arity 2)
    | 'SUBSCRIBE' arg variadicArgs // COMMAND subscribe (arity -2)
    | 'SUBSTR' arg arg arg // COMMAND substr (arity 4)
    | 'SUNION' arg variadicArgs // COMMAND sunion (arity -2)
    | 'SUNIONSTORE' arg arg variadicArgs // COMMAND sunionstore (arity -3)
    | 'SUNSUBSCRIBE' variadicArgs // COMMAND sunsubscribe (arity -1)
    | 'SWAPDB' arg arg // COMMAND swapdb (arity 3)
    | 'SYNC' // COMMAND sync (arity 1)
    | 'TIME' // COMMAND time (arity 1)
    | 'TOUCH' arg variadicArgs // COMMAND touch (arity -2)
    | 'TTL' arg // COMMAND ttl (arity 2)
    | 'TYPE' arg // COMMAND type (arity 2)
    | 'UNLINK' arg variadicArgs // COMMAND unlink (arity -2)
    | 'UNSUBSCRIBE' variadicArgs // COMMAND unsubscribe (arity -1)
    | 'UNWATCH' // COMMAND unwatch (arity 1)
    | 'WAIT' arg arg // COMMAND wait (arity 3)
    | 'WAITAOF' arg arg arg // COMMAND waitaof (arity 4)
    | 'WATCH' arg variadicArgs // COMMAND watch (arity -2)
    | 'XACK' arg arg arg variadicArgs // COMMAND xack (arity -4)
    | 'XADD' arg arg arg arg variadicArgs // COMMAND xadd (arity -5)
    | 'XAUTOCLAIM' arg arg arg arg arg variadicArgs // COMMAND xautoclaim (arity -6)
    | 'XCLAIM' arg arg arg arg arg variadicArgs // COMMAND xclaim (arity -6)
    | 'XDEL' arg arg variadicArgs // COMMAND xdel (arity -3)
    | 'XLEN' arg // COMMAND xlen (arity 2)
    | 'XPENDING' arg arg variadicArgs // COMMAND xpending (arity -3)
    | 'XRANGE' arg arg arg variadicArgs // COMMAND xrange (arity -4)
    | 'XREAD' arg arg arg variadicArgs // COMMAND xread (arity -4)
    | 'XREADGROUP' arg arg arg arg arg arg variadicArgs // COMMAND xreadgroup (arity -7)
    | 'XREVRANGE' arg arg arg variadicArgs // COMMAND xrevrange (arity -4)
    | 'XSETID' arg arg variadicArgs // COMMAND xsetid (arity -3)
    | 'XTRIM' arg arg arg variadicArgs // COMMAND xtrim (arity -4)
    | 'ZADD' arg arg arg variadicArgs // COMMAND zadd (arity -4)
    | 'ZCARD' arg // COMMAND zcard (arity 2)
    | 'ZCOUNT' arg arg arg // COMMAND zcount (arity 4)
    | 'ZDIFF' arg arg variadicArgs // COMMAND zdiff (arity -3)
    | 'ZDIFFSTORE' arg arg arg variadicArgs // COMMAND zdiffstore (arity -4)
    | 'ZINCRBY' arg arg arg // COMMAND zincrby (arity 4)
    | 'ZINTER' arg arg variadicArgs // COMMAND zinter (arity -3)
    | 'ZINTERCARD' arg arg variadicArgs // COMMAND zintercard (arity -3)
    | 'ZINTERSTORE' arg arg arg variadicArgs // COMMAND zinterstore (arity -4)
    | 'ZLEXCOUNT' arg arg arg // COMMAND zlexcount (arity 4)
    | 'ZMPOP' arg arg arg variadicArgs // COMMAND zmpop (arity -4)
    | 'ZMSCORE' arg arg variadicArgs // COMMAND zmscore (arity -3)
    | 'ZPOPMAX' arg variadicArgs // COMMAND zpopmax (arity -2)
    | 'ZPOPMIN' arg variadicArgs // COMMAND zpopmin (arity -2)
    | 'ZRANDMEMBER' arg variadicArgs // COMMAND zrandmember (arity -2)
    | 'ZRANGE' arg arg arg variadicArgs // COMMAND zrange (arity -4)
    | 'ZRANGEBYLEX' arg arg arg variadicArgs // COMMAND zrangebylex (arity -4)
    | 'ZRANGEBYSCORE' arg arg arg variadicArgs // COMMAND zrangebyscore (arity -4)
    | 'ZRANGESTORE' arg arg arg arg variadicArgs // COMMAND zrangestore (arity -5)
    | 'ZRANK' arg arg variadicArgs // COMMAND zrank (arity -3)
    | 'ZREM' arg arg variadicArgs // COMMAND zrem (arity -3)
    | 'ZREMRANGEBYLEX' arg arg arg // COMMAND zremrangebylex (arity 4)
    | 'ZREMRANGEBYRANK' arg arg arg // COMMAND zremrangebyrank (arity 4)
    | 'ZREMRANGEBYSCORE' arg arg arg // COMMAND zremrangebyscore (arity 4)
    | 'ZREVRANGE' arg arg arg variadicArgs // COMMAND zrevrange (arity -4)
    | 'ZREVRANGEBYLEX' arg arg arg variadicArgs // COMMAND zrevrangebylex (arity -4)
    | 'ZREVRANGEBYSCORE' arg arg arg variadicArgs // COMMAND zrevrangebyscore (arity -4)
    | 'ZREVRANK' arg arg variadicArgs // COMMAND zrevrank (arity -3)
    | 'ZSCAN' arg arg variadicArgs // COMMAND zscan (arity -3)
    | 'ZSCORE' arg arg // COMMAND zscore (arity 3)
    | 'ZUNION' arg arg variadicArgs // COMMAND zunion (arity -3)
    | 'ZUNIONSTORE' arg arg arg variadicArgs // COMMAND zunionstore (arity -4)
    ;

// Container command ACL and its subcommands.
aclCommand
    : 'ACL' 'CAT' variadicArgs // COMMAND acl|cat (arity -2)
    | 'ACL' 'DELUSER' arg variadicArgs // COMMAND acl|deluser (arity -3)
    | 'ACL' 'DRYRUN' arg arg variadicArgs // COMMAND acl|dryrun (arity -4)
    | 'ACL' 'GENPASS' variadicArgs // COMMAND acl|genpass (arity -2)
    | 'ACL' 'GETUSER' arg // COMMAND acl|getuser (arity 3)
    | 'ACL' 'HELP' // COMMAND acl|help (arity 2)
    | 'ACL' 'LIST' // COMMAND acl|list (arity 2)
    | 'ACL' 'LOAD' // COMMAND acl|load (arity 2)
    | 'ACL' 'LOG' variadicArgs // COMMAND acl|log (arity -2)
    | 'ACL' 'SAVE' // COMMAND acl|save (arity 2)
    | 'ACL' 'SETUSER' arg variadicArgs // COMMAND acl|setuser (arity -3)
    | 'ACL' 'USERS' // COMMAND acl|users (arity 2)
    | 'ACL' 'WHOAMI' // COMMAND acl|whoami (arity 2)
    ;

// Container command CLIENT and its subcommands.
clientCommand
    : 'CLIENT' 'CACHING' arg // COMMAND client|caching (arity 3)
    | 'CLIENT' 'GETNAME' // COMMAND client|getname (arity 2)
    | 'CLIENT' 'GETREDIR' // COMMAND client|getredir (arity 2)
    | 'CLIENT' 'HELP' // COMMAND client|help (arity 2)
    | 'CLIENT' 'ID' // COMMAND client|id (arity 2)
    | 'CLIENT' 'INFO' // COMMAND client|info (arity 2)
    | 'CLIENT' 'KILL' arg variadicArgs // COMMAND client|kill (arity -3)
    | 'CLIENT' 'LIST' variadicArgs // COMMAND client|list (arity -2)
    | 'CLIENT' 'NO-EVICT' arg // COMMAND client|no-evict (arity 3)
    | 'CLIENT' 'NO-TOUCH' arg // COMMAND client|no-touch (arity 3)
    | 'CLIENT' 'PAUSE' arg variadicArgs // COMMAND client|pause (arity -3)
    | 'CLIENT' 'REPLY' arg // COMMAND client|reply (arity 3)
    | 'CLIENT' 'SETINFO' arg arg // COMMAND client|setinfo (arity 4)
    | 'CLIENT' 'SETNAME' arg // COMMAND client|setname (arity 3)
    | 'CLIENT' 'TRACKING' arg variadicArgs // COMMAND client|tracking (arity -3)
    | 'CLIENT' 'TRACKINGINFO' // COMMAND client|trackinginfo (arity 2)
    | 'CLIENT' 'UNBLOCK' arg variadicArgs // COMMAND client|unblock (arity -3)
    | 'CLIENT' 'UNPAUSE' // COMMAND client|unpause (arity 2)
    ;

// Container command CLUSTER and its subcommands.
clusterCommand
    : 'CLUSTER' 'ADDSLOTS' arg variadicArgs // COMMAND cluster|addslots (arity -3)
    | 'CLUSTER' 'ADDSLOTSRANGE' arg arg variadicArgs // COMMAND cluster|addslotsrange (arity -4)
    | 'CLUSTER' 'BUMPEPOCH' // COMMAND cluster|bumpepoch (arity 2)
    | 'CLUSTER' 'COUNT-FAILURE-REPORTS' arg // COMMAND cluster|count-failure-reports (arity 3)
    | 'CLUSTER' 'COUNTKEYSINSLOT' arg // COMMAND cluster|countkeysinslot (arity 3)
    | 'CLUSTER' 'DELSLOTS' arg variadicArgs // COMMAND cluster|delslots (arity -3)
    | 'CLUSTER' 'DELSLOTSRANGE' arg arg variadicArgs // COMMAND cluster|delslotsrange (arity -4)
    | 'CLUSTER' 'FAILOVER' variadicArgs // COMMAND cluster|failover (arity -2)
    | 'CLUSTER' 'FLUSHSLOTS' // COMMAND cluster|flushslots (arity 2)
    | 'CLUSTER' 'FORGET' arg // COMMAND cluster|forget (arity 3)
    | 'CLUSTER' 'GETKEYSINSLOT' arg arg // COMMAND cluster|getkeysinslot (arity 4)
    | 'CLUSTER' 'HELP' // COMMAND cluster|help (arity 2)
    | 'CLUSTER' 'INFO' // COMMAND cluster|info (arity 2)
    | 'CLUSTER' 'KEYSLOT' arg // COMMAND cluster|keyslot (arity 3)
    | 'CLUSTER' 'LINKS' // COMMAND cluster|links (arity 2)
    | 'CLUSTER' 'MEET' arg arg variadicArgs // COMMAND cluster|meet (arity -4)
    | 'CLUSTER' 'MYID' // COMMAND cluster|myid (arity 2)
    | 'CLUSTER' 'MYSHARDID' // COMMAND cluster|myshardid (arity 2)
    | 'CLUSTER' 'NODES' // COMMAND cluster|nodes (arity 2)
    | 'CLUSTER' 'REPLICAS' arg // COMMAND cluster|replicas (arity 3)
    | 'CLUSTER' 'REPLICATE' arg // COMMAND cluster|replicate (arity 3)
    | 'CLUSTER' 'RESET' variadicArgs // COMMAND cluster|reset (arity -2)
    | 'CLUSTER' 'SAVECONFIG' // COMMAND cluster|saveconfig (arity 2)
    | 'CLUSTER' 'SET-CONFIG-EPOCH' arg // COMMAND cluster|set-config-epoch (arity 3)
    | 'CLUSTER' 'SETSLOT' arg arg variadicArgs // COMMAND cluster|setslot (arity -4)
    | 'CLUSTER' 'SHARDS' // COMMAND cluster|shards (arity 2)
    | 'CLUSTER' 'SLAVES' arg // COMMAND cluster|slaves (arity 3)
    | 'CLUSTER' 'SLOTS' // COMMAND cluster|slots (arity 2)
    ;

// Container command COMMAND and its subcommands.
commandCommand
    : 'COMMAND' // COMMAND command (arity -1)
    | 'COMMAND' 'COUNT' // COMMAND command|count (arity 2)
    | 'COMMAND' 'DOCS' variadicArgs // COMMAND command|docs (arity -2)
    | 'COMMAND' 'GETKEYS' arg variadicArgs // COMMAND command|getkeys (arity -3)
    | 'COMMAND' 'GETKEYSANDFLAGS' arg variadicArgs // COMMAND command|getkeysandflags (arity -3)
    | 'COMMAND' 'HELP' // COMMAND command|help (arity 2)
    | 'COMMAND' 'INFO' variadicArgs // COMMAND command|info (arity -2)
    | 'COMMAND' 'LIST' variadicArgs // COMMAND command|list (arity -2)
    ;

// Container command CONFIG and its subcommands.
configCommand
    : 'CONFIG' 'GET' arg variadicArgs // COMMAND config|get (arity -3)
    | 'CONFIG' 'HELP' // COMMAND config|help (arity 2)
    | 'CONFIG' 'RESETSTAT' // COMMAND config|resetstat (arity 2)
    | 'CONFIG' 'REWRITE' // COMMAND config|rewrite (arity 2)
    | 'CONFIG' 'SET' arg arg variadicArgs // COMMAND config|set (arity -4)
    ;

// Container command FUNCTION and its subcommands.
functionCommand
    : 'FUNCTION' 'DELETE' arg // COMMAND function|delete (arity 3)
    | 'FUNCTION' 'DUMP' // COMMAND function|dump (arity 2)
    | 'FUNCTION' 'FLUSH' variadicArgs // COMMAND function|flush (arity -2)
    | 'FUNCTION' 'HELP' // COMMAND function|help (arity 2)
    | 'FUNCTION' 'KILL' // COMMAND function|kill (arity 2)
    | 'FUNCTION' 'LIST' variadicArgs // COMMAND function|list (arity -2)
    | 'FUNCTION' 'LOAD' arg variadicArgs // COMMAND function|load (arity -3)
    | 'FUNCTION' 'RESTORE' arg variadicArgs // COMMAND function|restore (arity -3)
    | 'FUNCTION' 'STATS' // COMMAND function|stats (arity 2)
    ;

// Container command LATENCY and its subcommands.
latencyCommand
    : 'LATENCY' 'DOCTOR' // COMMAND latency|doctor (arity 2)
    | 'LATENCY' 'GRAPH' arg // COMMAND latency|graph (arity 3)
    | 'LATENCY' 'HELP' // COMMAND latency|help (arity 2)
    | 'LATENCY' 'HISTOGRAM' variadicArgs // COMMAND latency|histogram (arity -2)
    | 'LATENCY' 'HISTORY' arg // COMMAND latency|history (arity 3)
    | 'LATENCY' 'LATEST' // COMMAND latency|latest (arity 2)
    | 'LATENCY' 'RESET' variadicArgs // COMMAND latency|reset (arity -2)
    ;

// Container command MEMORY and its subcommands.
memoryCommand
    : 'MEMORY' 'DOCTOR' // COMMAND memory|doctor (arity 2)
    | 'MEMORY' 'HELP' // COMMAND memory|help (arity 2)
    | 'MEMORY' 'MALLOC-STATS' // COMMAND memory|malloc-stats (arity 2)
    | 'MEMORY' 'PURGE' // COMMAND memory|purge (arity 2)
    | 'MEMORY' 'STATS' // COMMAND memory|stats (arity 2)
    | 'MEMORY' 'USAGE' arg variadicArgs // COMMAND memory|usage (arity -3)
    ;

// Container command MODULE and its subcommands.
moduleCommand
    : 'MODULE' 'HELP' // COMMAND module|help (arity 2)
    | 'MODULE' 'LIST' // COMMAND module|list (arity 2)
    | 'MODULE' 'LOAD' arg variadicArgs // COMMAND module|load (arity -3)
    | 'MODULE' 'LOADEX' arg variadicArgs // COMMAND module|loadex (arity -3)
    | 'MODULE' 'UNLOAD' arg // COMMAND module|unload (arity 3)
    ;

// Container command OBJECT and its subcommands.
objectCommand
    : 'OBJECT' 'ENCODING' arg // COMMAND object|encoding (arity 3)
    | 'OBJECT' 'FREQ' arg // COMMAND object|freq (arity 3)
    | 'OBJECT' 'HELP' // COMMAND object|help (arity 2)
    | 'OBJECT' 'IDLETIME' arg // COMMAND object|idletime (arity 3)
    | 'OBJECT' 'REFCOUNT' arg // COMMAND object|refcount (arity 3)
    ;

// Container command PUBSUB and its subcommands.
pubsubCommand
    : 'PUBSUB' 'CHANNELS' variadicArgs // COMMAND pubsub|channels (arity -2)
    | 'PUBSUB' 'HELP' // COMMAND pubsub|help (arity 2)
    | 'PUBSUB' 'NUMPAT' // COMMAND pubsub|numpat (arity 2)
    | 'PUBSUB' 'NUMSUB' variadicArgs // COMMAND pubsub|numsub (arity -2)
    | 'PUBSUB' 'SHARDCHANNELS' variadicArgs // COMMAND pubsub|shardchannels (arity -2)
    | 'PUBSUB' 'SHARDNUMSUB' variadicArgs // COMMAND pubsub|shardnumsub (arity -2)
    ;

// Container command SCRIPT and its subcommands.
scriptCommand
    : 'SCRIPT' 'DEBUG' arg // COMMAND script|debug (arity 3)
    | 'SCRIPT' 'EXISTS' arg variadicArgs // COMMAND script|exists (arity -3)
    | 'SCRIPT' 'FLUSH' variadicArgs // COMMAND script|flush (arity -2)
    | 'SCRIPT' 'HELP' // COMMAND script|help (arity 2)
    | 'SCRIPT' 'KILL' // COMMAND script|kill (arity 2)
    | 'SCRIPT' 'LOAD' arg // COMMAND script|load (arity 3)
    ;

// Container command SLOWLOG and its subcommands.
slowlogCommand
    : 'SLOWLOG' 'GET' variadicArgs // COMMAND slowlog|get (arity -2)
    | 'SLOWLOG' 'HELP' // COMMAND slowlog|help (arity 2)
    | 'SLOWLOG' 'LEN' // COMMAND slowlog|len (arity 2)
    | 'SLOWLOG' 'RESET' // COMMAND slowlog|reset (arity 2)
    ;

// Container command XGROUP and its subcommands.
xgroupCommand
    : 'XGROUP' 'CREATE' arg arg arg variadicArgs // COMMAND xgroup|create (arity -5)
    | 'XGROUP' 'CREATECONSUMER' arg arg arg // COMMAND xgroup|createconsumer (arity 5)
    | 'XGROUP' 'DELCONSUMER' arg arg arg // COMMAND xgroup|delconsumer (arity 5)
    | 'XGROUP' 'DESTROY' arg arg // COMMAND xgroup|destroy (arity 4)
    | 'XGROUP' 'HELP' // COMMAND xgroup|help (arity 2)
    | 'XGROUP' 'SETID' arg arg arg variadicArgs // COMMAND xgroup|setid (arity -5)
    ;

// Container command XINFO and its subcommands.
xinfoCommand
    : 'XINFO' 'CONSUMERS' arg arg // COMMAND xinfo|consumers (arity 4)
    | 'XINFO' 'GROUPS' arg // COMMAND xinfo|groups (arity 3)
    | 'XINFO' 'HELP' // COMMAND xinfo|help (arity 2)
    | 'XINFO' 'STREAM' arg variadicArgs // COMMAND xinfo|stream (arity -3)
    ;

// Zero or more extra arguments for variadic (negative-arity) commands.
variadicArgs
    :
    | arg variadicArgs
    ;

// Command-line argument token forms (see redis/testdata/tokens.tsv).
arg
    : WORD
    | SINGLE_QUOTED
    | DOUBLE_QUOTED
    | HASH_ARG
    | INT_LITERAL
    | FLOAT_LITERAL
    ;

// Lexer token classes for command-line arguments.
SINGLE_QUOTED : '\'' (~['\r\n])* '\'' ;
DOUBLE_QUOTED : '"' ('\\' . | ~["\\\r\n])* '"' ;
HASH_ARG      : '#' ~[ \t\r\n]+ ;
FLOAT_LITERAL : '-'? [0-9]+ '.' [0-9]+ ;
INT_LITERAL   : '-'? [0-9]+ ;
WORD          : ~[ \t\r\n'"#]+ ;
WS            : [ \t\r]+ -> skip ;
