// sh.updateZoneKeyRange() - Associate a range of shard key values with a zone

// Basic usage
sh.updateZoneKeyRange("mydb.users", { zipcode: "10001" }, { zipcode: "10099" }, "NYC")

// Remove zone assignment (null zone)
sh.updateZoneKeyRange("mydb.users", { zipcode: "10001" }, { zipcode: "10099" }, null)

// Compound shard key range
sh.updateZoneKeyRange(
    "mydb.events",
    { region: "US", timestamp: ISODate("2024-01-01") },
    { region: "US", timestamp: ISODate("2025-01-01") },
    "US_2024"
)

// Geographic zone assignment
sh.updateZoneKeyRange(
    "mydb.customers",
    { country: "FR" },
    { country: "FZ" },
    "EU-WEST"
)
