// Line comment
db.users.find({ name: "alice" }) // inline comment

// Block comment
/* This is a block comment */
db.users.find({ age: 25 })

/*
 * Multi-line
 * block comment
 */
db.users.find({
    /* comment inside document */
    name: "test",
    // another comment
    age: 30
})

// Multiple statements with comments
show dbs // list all databases
show collections /* list collections */
db.users.find() // find all users
