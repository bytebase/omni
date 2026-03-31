// sh.removeTagRange() - Remove a tag range (deprecated, use updateZoneKeyRange)

// Basic usage
sh.removeTagRange("mydb.users", { zipcode: "10001" }, { zipcode: "10099" }, "NYC")

// Remove compound key range
sh.removeTagRange(
    "mydb.events",
    { date: ISODate("2024-01-01"), region: "US" },
    { date: ISODate("2024-12-31"), region: "US" },
    "US_2024"
)
