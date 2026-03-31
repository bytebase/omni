// sp.getProcessor() - Get a stream processor by name

// Get processor by name
sp.getProcessor("myProcessor")

// Get processor and check its state
sp.getProcessor("analyticsProcessor")

// Get processor with different names
sp.getProcessor("dataProcessor")
sp.getProcessor("eventProcessor")
sp.getProcessor("sensorProcessor")

// Get processor with special naming patterns
sp.getProcessor("my_stream_processor")
sp.getProcessor("processor-v1")
sp.getProcessor("stream_2024")
