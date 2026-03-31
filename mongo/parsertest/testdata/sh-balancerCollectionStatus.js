// sh.balancerCollectionStatus() - Get balancer status for a specific collection

// Basic usage
sh.balancerCollectionStatus("mydb.users")

// Check status for sharded collection
sh.balancerCollectionStatus("test.orders")

// Namespace with dot notation
sh.balancerCollectionStatus("analytics.events")
