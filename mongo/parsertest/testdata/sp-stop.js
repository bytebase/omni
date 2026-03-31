// sp.<processor>.stop() - Stop a stream processor

// Stop processor via direct access
sp.myProcessor.stop()

// Stop various processors
sp.analyticsProcessor.stop()
sp.dataProcessor.stop()
sp.eventProcessor.stop()
sp.sensorProcessor.stop()

// Stop processors with different naming patterns
sp.my_stream_processor.stop()
sp.processor1.stop()
sp.stream_v2.stop()

// Stop pipeline processors
sp.ingestPipeline.stop()
sp.transformPipeline.stop()
sp.outputPipeline.stop()

// Stop for maintenance
sp.productionProcessor.stop()
sp.stagingProcessor.stop()
