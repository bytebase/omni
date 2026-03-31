// sh.addShard() - Add a shard to a sharded cluster

// Add a standalone mongod
sh.addShard("hostname:27017")

// Add a replica set shard
sh.addShard("rs0/hostname1:27017,hostname2:27017,hostname3:27017")

// Add with replica set name prefix
sh.addShard("shard0001/mongodb0.example.net:27017")

// Add multiple hosts for replica set
sh.addShard("rs1/host1:27017,host2:27017")
