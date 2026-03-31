// sh.abortReshardCollection() - Abort an in-progress resharding operation

// Basic usage
sh.abortReshardCollection("mydb.users")

// Abort resharding for specific namespace
sh.abortReshardCollection("test.orders")
