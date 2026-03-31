// db.collection.configureQueryAnalyzer() - Configure query analyzer for a collection

// Enable query analyzer with mode
db.users.configureQueryAnalyzer({ mode: "full" })
db.orders.configureQueryAnalyzer({ mode: "off" })

// Configure with sample rate
db.products.configureQueryAnalyzer({ mode: "full", sampleRate: 100 })
db.events.configureQueryAnalyzer({ mode: "full", sampleRate: 50 })
db.logs.configureQueryAnalyzer({ mode: "full", sampleRate: 10 })

// Disable query analyzer
db.users.configureQueryAnalyzer({ mode: "off" })
db.analytics.configureQueryAnalyzer({ mode: "off" })

// Configure for specific analysis
db.orders.configureQueryAnalyzer({
    mode: "full",
    sampleRate: 200
})

db.transactions.configureQueryAnalyzer({
    mode: "full",
    sampleRate: 1000
})

// Collection access patterns
db["users"].configureQueryAnalyzer({ mode: "full" })
db["users"].configureQueryAnalyzer({ mode: "full", sampleRate: 100 })
db.getCollection("users").configureQueryAnalyzer({ mode: "off" })
db.getCollection("orders").configureQueryAnalyzer({ mode: "full", sampleRate: 50 })
db["query-logs"].configureQueryAnalyzer({ mode: "full", sampleRate: 25 })
