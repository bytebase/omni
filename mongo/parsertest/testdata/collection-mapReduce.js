// db.collection.mapReduce() - Perform map-reduce aggregation (deprecated - use aggregation pipeline)
// Note: JavaScript function expressions are not supported by this parser.
// The examples below show the method structure with string placeholders.

// Basic mapReduce with string function representations
db.orders.mapReduce(
    "function() { emit(this.customerId, this.amount); }",
    "function(key, values) { return Array.sum(values); }",
    { out: "order_totals" }
)

// mapReduce with query
db.orders.mapReduce(
    "function() { emit(this.product, 1); }",
    "function(key, values) { return Array.sum(values); }",
    {
        query: { status: "completed" },
        out: "product_counts"
    }
)

// mapReduce with inline output
db.products.mapReduce(
    "function() { emit(this.category, this.price); }",
    "function(key, values) { return Array.avg(values); }",
    { out: { inline: 1 } }
)

// mapReduce with finalize function
db.scores.mapReduce(
    "function() { emit(this.playerId, { sum: this.score, count: 1 }); }",
    "function(key, values) { return { sum: 0, count: 0 }; }",
    {
        finalize: "function(key, reduced) { return reduced.sum / reduced.count; }",
        out: "player_averages"
    }
)

// mapReduce with sort and limit
db.logs.mapReduce(
    "function() { emit(this.level, 1); }",
    "function(key, values) { return Array.sum(values); }",
    {
        sort: { timestamp: -1 },
        limit: 10000,
        out: { inline: 1 }
    }
)

// mapReduce with scope
db.orders.mapReduce(
    "function() { emit(this.customerId, this.amount * factor); }",
    "function(key, values) { return Array.sum(values); }",
    {
        scope: { factor: 1.1 },
        out: "adjusted_totals"
    }
)

// mapReduce with jsMode
db.events.mapReduce(
    "function() { emit(this.type, 1); }",
    "function(key, values) { return Array.sum(values); }",
    {
        jsMode: true,
        out: { inline: 1 }
    }
)

// mapReduce with output options
db.sales.mapReduce(
    "function() { emit(this.region, this.revenue); }",
    "function(key, values) { return Array.sum(values); }",
    { out: { replace: "region_revenue" } }
)

db.sales.mapReduce(
    "function() { emit(this.region, this.revenue); }",
    "function(key, values) { return Array.sum(values); }",
    { out: { merge: "region_revenue" } }
)

db.sales.mapReduce(
    "function() { emit(this.region, this.revenue); }",
    "function(key, values) { return Array.sum(values); }",
    { out: { reduce: "region_revenue" } }
)

// mapReduce with output to different database
db.orders.mapReduce(
    "function() { emit(this.product, this.qty); }",
    "function(key, values) { return Array.sum(values); }",
    { out: { replace: "product_totals", db: "reports" } }
)

// mapReduce with verbose output
db.logs.mapReduce(
    "function() { emit(this.source, 1); }",
    "function(key, values) { return Array.sum(values); }",
    {
        out: { inline: 1 },
        verbose: true
    }
)

// mapReduce with bypassDocumentValidation
db.data.mapReduce(
    "function() { emit(this.key, this.value); }",
    "function(key, values) { return values.join(','); }",
    {
        out: "combined_data",
        bypassDocumentValidation: true
    }
)

// mapReduce with collation
db.products.mapReduce(
    "function() { emit(this.name.toLowerCase(), 1); }",
    "function(key, values) { return Array.sum(values); }",
    {
        out: { inline: 1 },
        collation: { locale: "en", strength: 2 }
    }
)

// Collection access patterns
db["orders"].mapReduce("function() { emit(this.x, 1); }", "function(k, v) { return Array.sum(v); }", { out: { inline: 1 } })
db.getCollection("orders").mapReduce("function() { emit(this.y, 1); }", "function(k, v) { return Array.sum(v); }", { out: "result" })
db["order-items"].mapReduce("function() { emit(this.orderId, this.total); }", "function(k, v) { return Array.sum(v); }", { out: { inline: 1 } })
