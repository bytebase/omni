// sp.createStreamProcessor() - Create a stream processor

// Basic stream processor
sp.createStreamProcessor("myProcessor", {
    pipeline: [
        { $source: { connectionName: "connection1" } },
        { $match: { status: "active" } },
        { $emit: { connectionName: "connection2" } }
    ]
})

// With options
sp.createStreamProcessor("analyticsProcessor", {
    pipeline: [
        { $source: { connectionName: "logs" } },
        { $merge: { into: "processed_logs" } }
    ],
    options: { dlq: { connectionName: "dlq_connection" } }
})

// Stream processor with aggregation stages
sp.createStreamProcessor("aggregateProcessor", {
    pipeline: [
        { $source: { connectionName: "events" } },
        { $match: { type: "purchase" } },
        { $group: { _id: "$userId", total: { $sum: "$amount" } } },
        { $emit: { connectionName: "aggregated_events" } }
    ]
})

// Stream processor with window operations
sp.createStreamProcessor("windowProcessor", {
    pipeline: [
        { $source: { connectionName: "metrics" } },
        { $tumblingWindow: {
            interval: { size: 1, unit: "minute" },
            pipeline: [
                { $group: { _id: "$sensor", avg: { $avg: "$value" } } }
            ]
        }},
        { $merge: { into: "metric_averages" } }
    ]
})

// Stream processor with dead letter queue
sp.createStreamProcessor("dlqProcessor", {
    pipeline: [
        { $source: { connectionName: "input_stream" } },
        { $match: { valid: true } },
        { $emit: { connectionName: "output_stream" } }
    ],
    options: {
        dlq: {
            connectionName: "dlq_connection",
            db: "errors",
            coll: "failed_records"
        }
    }
})
