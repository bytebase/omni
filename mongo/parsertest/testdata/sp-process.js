// sp.process() - Run an inline stream processing pipeline

// Basic process with pipeline
sp.process([
    { $source: { connectionName: "sample_stream_solar" } },
    { $match: { status: "active" } },
    { $emit: { connectionName: "output_stream" } }
])

// Process with aggregation
sp.process([
    { $source: { connectionName: "events" } },
    { $group: { _id: "$category", count: { $sum: 1 } } },
    { $merge: { into: "event_counts" } }
])

// Process with window function
sp.process([
    { $source: { connectionName: "sensor_data" } },
    { $tumblingWindow: {
        interval: { size: 5, unit: "second" },
        pipeline: [
            { $group: { _id: null, avgTemp: { $avg: "$temperature" } } }
        ]
    }},
    { $emit: { connectionName: "temp_averages" } }
])

// Process with complex transformations
sp.process([
    { $source: { connectionName: "raw_logs" } },
    { $match: { level: { $in: ["ERROR", "WARN"] } } },
    { $project: {
        timestamp: 1,
        level: 1,
        message: 1,
        processed: true
    }},
    { $merge: { into: "filtered_logs" } }
])
