// sh.moveRange() - Move a range of shard key values to a shard

// Basic usage
sh.moveRange("mydb.users", { min: { zipcode: "10001" }, max: { zipcode: "20000" } }, "shard0001")

// Move range with specific bounds
sh.moveRange(
    "test.orders",
    { min: { customerId: 1000 }, max: { customerId: 2000 } },
    "shard0002"
)

// Move with forceJumbo option
sh.moveRange(
    "mydb.events",
    { min: { region: "US" }, max: { region: "UZ" } },
    "shard0003",
    { forceJumbo: true }
)
