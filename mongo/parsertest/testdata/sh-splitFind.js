// sh.splitFind() - Split the chunk that contains a document matching the query

// Basic usage
sh.splitFind("mydb.users", { zipcode: "10001" })

// Split chunk containing document
sh.splitFind("test.orders", { customerId: 12345 })

// Split with compound key query
sh.splitFind("analytics.events", { region: "EU", timestamp: ISODate("2024-03-15") })
