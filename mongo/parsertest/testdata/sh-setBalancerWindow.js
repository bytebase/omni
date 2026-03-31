// sh.setBalancerWindow() - Set the balancer window

// Set balancer window with start and stop times
sh.setBalancerWindow("23:00", "06:00")

// Set overnight balancing window
sh.setBalancerWindow("02:00", "04:00")

// Set with options
sh.setBalancerWindow({ start: "00:00", stop: "05:00" })
