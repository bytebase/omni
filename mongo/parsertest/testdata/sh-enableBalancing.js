// sh.enableBalancing() - Enable balancing for a specific collection

// Basic usage
sh.enableBalancing("mydb.users")

// Enable for specific namespace
sh.enableBalancing("test.orders")

// Enable after maintenance window
sh.enableBalancing("analytics.events")
