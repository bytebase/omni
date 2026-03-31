// sh.addTagRange() - Associate a range of shard key values with a tag (deprecated, use updateZoneKeyRange)

// Basic usage with numeric range
sh.addTagRange("mydb.users", { zipcode: "10001" }, { zipcode: "10099" }, "NYC")

// Range with MinKey/MaxKey
sh.addTagRange("mydb.orders", { region: "A" }, { region: "M" }, "REGION_A_M")

// Compound shard key range
sh.addTagRange(
    "mydb.events",
    { date: ISODate("2024-01-01"), region: "US" },
    { date: ISODate("2024-12-31"), region: "US" },
    "US_2024"
)
