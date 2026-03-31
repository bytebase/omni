// sh.configureQueryAnalyzer() - Configure the query analyzer for a collection

// Enable query analysis
sh.configureQueryAnalyzer("mydb.users", { mode: "full", sampleRate: 1.0 })

// Disable query analysis
sh.configureQueryAnalyzer("mydb.users", { mode: "off" })

// Sample a portion of queries
sh.configureQueryAnalyzer("mydb.orders", { mode: "full", sampleRate: 0.1 })
