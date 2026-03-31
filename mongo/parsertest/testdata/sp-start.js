// sp.<processor>.start() - Start a stream processor

// Start processor via direct access
sp.myProcessor.start()

// Start various processors
sp.analyticsProcessor.start()
sp.dataProcessor.start()
sp.eventProcessor.start()
sp.sensorProcessor.start()

// Start processors with different naming patterns
sp.my_stream_processor.start()
sp.processor1.start()
sp.stream_v2.start()

// Start a specific processor for data pipeline
sp.ingestPipeline.start()
sp.transformPipeline.start()
sp.outputPipeline.start()
