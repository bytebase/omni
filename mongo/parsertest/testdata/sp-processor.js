// sp.<processor-name> - Access a specific stream processor

// Access processor by name and check stats
sp.myProcessor.stats()

// Access processor and start it
sp.analyticsProcessor.start()

// Access processor and stop it
sp.dataProcessor.stop()

// Access processor and drop it
sp.oldProcessor.drop()

// Access processor and get sample data
sp.sensorProcessor.sample()

// Chain operations on processor
sp.eventProcessor.stats()
sp.eventProcessor.start()
sp.eventProcessor.stop()

// Access processor with underscore in name
sp.my_stream_processor.stats()

// Access processor with number in name
sp.processor1.stats()
sp.stream_v2.start()
