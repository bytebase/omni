// ObjectId
db.users.find({ _id: ObjectId("507f1f77bcf86cd799439011") })
db.users.find({ _id: ObjectId() })

// ISODate
db.events.find({ createdAt: ISODate("2024-01-15T00:00:00.000Z") })
db.events.find({ createdAt: { $gt: ISODate() } })

// Date
db.events.find({ timestamp: Date() })
db.events.find({ timestamp: Date("2024-01-15") })
db.events.find({ timestamp: Date(1705276800000) })

// UUID
db.sessions.find({ sessionId: UUID("550e8400-e29b-41d4-a716-446655440000") })

// Long / NumberLong
db.stats.find({ count: Long(9007199254740993) })
db.stats.find({ count: Long("9007199254740993") })
db.stats.find({ count: NumberLong(123456789012345) })
db.stats.find({ count: NumberLong("123456789012345") })

// Int32 / NumberInt
db.items.find({ quantity: Int32(100) })
db.items.find({ quantity: NumberInt(100) })

// Double
db.measurements.find({ value: Double(3.14159) })
db.measurements.find({ value: Double(-273.15) })

// Decimal128 / NumberDecimal
db.financial.find({ amount: Decimal128("1234567890.123456789") })
db.financial.find({ amount: NumberDecimal("99.99") })

// Timestamp
db.oplog.find({ ts: Timestamp(1627811580, 1) })
db.oplog.find({ ts: Timestamp({ t: 1627811580, i: 1 }) })

// Mixed helper functions in complex queries
db.users.find({
    _id: ObjectId("507f1f77bcf86cd799439011"),
    createdAt: { $gt: ISODate("2024-01-01T00:00:00Z") },
    lastLogin: { $lt: Date() },
    sessionId: UUID("550e8400-e29b-41d4-a716-446655440000"),
    loginCount: NumberLong(1000)
})
