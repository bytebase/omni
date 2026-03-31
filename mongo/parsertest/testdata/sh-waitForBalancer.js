// sh.waitForBalancer() - Wait for a balancing round to complete

// Basic usage (no arguments)
sh.waitForBalancer()

// With timeout
sh.waitForBalancer(true)

// With timeout in milliseconds
sh.waitForBalancer(true, 60000)
