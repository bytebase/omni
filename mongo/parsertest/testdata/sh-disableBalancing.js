// sh.disableBalancing() - Disable balancing for a specific collection

// Basic usage
sh.disableBalancing("mydb.users")

// Disable for specific namespace
sh.disableBalancing("test.orders")

// Disable for collection in analytics database
sh.disableBalancing("analytics.events")
