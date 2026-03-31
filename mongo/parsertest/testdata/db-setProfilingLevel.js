// db.setProfilingLevel() - Set the database profiler level

// Turn off profiling
db.setProfilingLevel(0)

// Profile slow operations only (default threshold 100ms)
db.setProfilingLevel(1)

// Profile all operations
db.setProfilingLevel(2)

// Set custom slow operation threshold
db.setProfilingLevel(1, 50)
db.setProfilingLevel(1, 200)

// With options object
db.setProfilingLevel(1, { slowms: 100 })
db.setProfilingLevel(1, { slowms: 50, sampleRate: 0.5 })

// Profile with sample rate
db.setProfilingLevel(2, { sampleRate: 0.1 })
