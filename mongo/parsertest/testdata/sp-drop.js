// sp.<processor>.drop() - Drop (delete) a stream processor

// Drop processor via direct access
sp.myProcessor.drop()

// Drop various processors
sp.analyticsProcessor.drop()
sp.dataProcessor.drop()
sp.eventProcessor.drop()
sp.sensorProcessor.drop()

// Drop processors with different naming patterns
sp.my_stream_processor.drop()
sp.processor1.drop()
sp.stream_v2.drop()

// Drop old or deprecated processors
sp.oldProcessor.drop()
sp.legacyProcessor.drop()
sp.testProcessor.drop()

// Drop pipeline processors
sp.deprecatedPipeline.drop()
sp.temporaryProcessor.drop()
