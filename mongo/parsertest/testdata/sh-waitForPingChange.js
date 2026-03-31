// sh.waitForPingChange() - Wait for a change in ping times

// Basic usage with activePings array
sh.waitForPingChange([{ _id: "shard0000" }, { _id: "shard0001" }])

// With timeout
sh.waitForPingChange([{ _id: "shard0000" }], 60000)

// With options
sh.waitForPingChange([{ _id: "shard0000" }], { timeout: 60000 })
