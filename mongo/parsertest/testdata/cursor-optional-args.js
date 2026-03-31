// Test cursor methods with optional arguments

// find and findOne with multiple args
db.coll.find({}, {a: 1})
db.coll.findOne({name: "test"}, {_id: 0})

// Empty cursor methods
db.coll.find().sort()
db.coll.find().collation()
db.coll.find().comment()
db.coll.find().hint()
db.coll.find().max()
db.coll.find().min()
db.coll.find().readConcern()
db.coll.find().returnKey()
db.coll.find().showRecordId()
db.coll.find().projection()

// aggregate with empty pipeline
db.coll.aggregate()
