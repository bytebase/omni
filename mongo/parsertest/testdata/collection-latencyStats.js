// db.collection.latencyStats() - Get latency statistics for a collection

// Basic latency stats
db.users.latencyStats()
db.orders.latencyStats()

// With histogram options
db.users.latencyStats({ histograms: true })

// Without histogram details
db.orders.latencyStats({ histograms: false })

// Collection access patterns
db["users"].latencyStats()
db.getCollection("users").latencyStats()
db["high-traffic"].latencyStats({ histograms: true })
db.getCollection("production.api").latencyStats()
