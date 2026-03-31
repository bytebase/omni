// $function operator - custom JavaScript functions in aggregation
// The JavaScript code is passed as a string in the "body" field

// Basic $function usage
db.players.aggregate([
    { $addFields: {
        isFound: {
            $function: {
                body: "function(name) { return name.length > 5; }",
                args: ["$name"],
                lang: "js"
            }
        }
    }}
])

// $function with multiple arguments
db.orders.aggregate([
    { $addFields: {
        discount: {
            $function: {
                body: "function(price, quantity) { return price * quantity * 0.1; }",
                args: ["$price", "$quantity"],
                lang: "js"
            }
        }
    }}
])

// $function with complex logic
db.events.aggregate([
    { $addFields: {
        category: {
            $function: {
                body: "function(score) { if (score >= 90) return 'A'; else if (score >= 80) return 'B'; else if (score >= 70) return 'C'; else return 'F'; }",
                args: ["$score"],
                lang: "js"
            }
        }
    }}
])

// $function with no arguments
db.items.aggregate([
    { $addFields: {
        timestamp: {
            $function: {
                body: "function() { return new Date().toISOString(); }",
                args: [],
                lang: "js"
            }
        }
    }}
])

// $function with array argument
db.data.aggregate([
    { $addFields: {
        processed: {
            $function: {
                body: "function(arr) { return arr.map(x => x * 2); }",
                args: ["$values"],
                lang: "js"
            }
        }
    }}
])

// $function with multiline string (escaped)
db.documents.aggregate([
    { $addFields: {
        result: {
            $function: {
                body: "function(doc) {\n  var sum = 0;\n  for (var i = 0; i < doc.length; i++) {\n    sum += doc[i];\n  }\n  return sum;\n}",
                args: ["$numbers"],
                lang: "js"
            }
        }
    }}
])

// $accumulator operator - custom accumulator with JavaScript
db.sales.aggregate([
    { $group: {
        _id: "$category",
        customSum: {
            $accumulator: {
                init: "function() { return { sum: 0, count: 0 }; }",
                accumulate: "function(state, value) { return { sum: state.sum + value, count: state.count + 1 }; }",
                accumulateArgs: ["$amount"],
                merge: "function(state1, state2) { return { sum: state1.sum + state2.sum, count: state1.count + state2.count }; }",
                finalize: "function(state) { return state.sum / state.count; }",
                lang: "js"
            }
        }
    }}
])

// $accumulator with initArgs
db.transactions.aggregate([
    { $group: {
        _id: "$accountId",
        balance: {
            $accumulator: {
                init: "function(initial) { return initial; }",
                initArgs: [0],
                accumulate: "function(state, amount, type) { return type === 'credit' ? state + amount : state - amount; }",
                accumulateArgs: ["$amount", "$type"],
                merge: "function(s1, s2) { return s1 + s2; }",
                lang: "js"
            }
        }
    }}
])

// $accumulator without finalize
db.metrics.aggregate([
    { $group: {
        _id: "$source",
        values: {
            $accumulator: {
                init: "function() { return []; }",
                accumulate: "function(state, val) { state.push(val); return state; }",
                accumulateArgs: ["$value"],
                merge: "function(s1, s2) { return s1.concat(s2); }",
                lang: "js"
            }
        }
    }}
])

// Combined $function and standard operators
db.products.aggregate([
    { $match: { inStock: true } },
    { $addFields: {
        priceCategory: {
            $function: {
                body: "function(p) { return p < 10 ? 'cheap' : p < 100 ? 'moderate' : 'expensive'; }",
                args: ["$price"],
                lang: "js"
            }
        }
    }},
    { $group: { _id: "$priceCategory", count: { $sum: 1 } } },
    { $sort: { count: -1 } }
])

// $function in $project stage
db.users.aggregate([
    { $project: {
        name: 1,
        maskedEmail: {
            $function: {
                body: "function(email) { var parts = email.split('@'); return parts[0].substring(0, 2) + '***@' + parts[1]; }",
                args: ["$email"],
                lang: "js"
            }
        }
    }}
])

// $function with object argument
db.records.aggregate([
    { $addFields: {
        serialized: {
            $function: {
                body: "function(obj) { return JSON.stringify(obj); }",
                args: ["$metadata"],
                lang: "js"
            }
        }
    }}
])

// $function returning object
db.coordinates.aggregate([
    { $addFields: {
        point: {
            $function: {
                body: "function(lat, lng) { return { type: 'Point', coordinates: [lng, lat] }; }",
                args: ["$latitude", "$longitude"],
                lang: "js"
            }
        }
    }}
])

// Nested $function in $facet
db.analytics.aggregate([
    { $facet: {
        processed: [
            { $addFields: {
                normalized: {
                    $function: {
                        body: "function(v, min, max) { return (v - min) / (max - min); }",
                        args: ["$value", "$minValue", "$maxValue"],
                        lang: "js"
                    }
                }
            }}
        ],
        summary: [
            { $group: { _id: null, total: { $sum: "$value" } } }
        ]
    }}
])
