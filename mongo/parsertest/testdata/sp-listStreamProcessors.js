// sp.listStreamProcessors() - List all stream processors

// List all stream processors
sp.listStreamProcessors()

// List stream processors with filter
sp.listStreamProcessors({ name: "myProcessor" })

// List stream processors by state
sp.listStreamProcessors({ state: "RUNNING" })
sp.listStreamProcessors({ state: "STOPPED" })
sp.listStreamProcessors({ state: "CREATED" })

// List stream processors with multiple filters
sp.listStreamProcessors({ state: "RUNNING", name: "analyticsProcessor" })

// List stream processors with regex pattern
sp.listStreamProcessors({ name: /^analytics.*/ })
