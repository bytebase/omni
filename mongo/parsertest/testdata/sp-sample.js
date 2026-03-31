// sp.<processor>.sample() - Sample documents from a stream processor

// Sample from processor via direct access
sp.myProcessor.sample()

// Sample from various processors
sp.analyticsProcessor.sample()
sp.dataProcessor.sample()
sp.eventProcessor.sample()
sp.sensorProcessor.sample()

// Sample from processors with different naming patterns
sp.my_stream_processor.sample()
sp.processor1.sample()
sp.stream_v2.sample()

// Sample from pipeline processors
sp.ingestPipeline.sample()
sp.transformPipeline.sample()
sp.outputPipeline.sample()

// Sample for debugging
sp.debugProcessor.sample()
sp.testProcessor.sample()
