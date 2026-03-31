// rs.remove() - Remove a member from the replica set

// Remove member by hostname
rs.remove("mongo4:27017")

// Remove member with full hostname
rs.remove("mongo4.example.com:27017")
