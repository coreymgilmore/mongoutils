# mongoutils.go
Helper Functions for using the Golang mgo Driver for MongoDB

The intended use of this package is to simplify some functionality.  It is not meant for power users or for production.
Basically, this package can be used to connect or disconnect to a MongoDB (or replica set), perform ObjectId validation, convert between ObjectId and a string, and grab limit and sort values from a form value.  This package will also save the connection data to a global variable that is accessible wherever this file is imported.

Of note: when you use the connection data, you will need to "copy" the connection in order to use the pool of available connections.  More about this can be found in the mgo docs.

---

###Functions

####Connect(servers string, database string, readPreference int, writeConcern *mgo.Safe)
- Connect to your DB(s).  Stores the session connection data into a global variable.
- `servers` is a single or list of servers in "servername:port" format separated by commas and ending with a "/".
- `database` is the name of your database.
- `readPreference` is an mgo constant that defines how you want reads to be spread across the available servers.
` `writeConcern` is the safety level at which MongoDB will determine if a write was successful or not.  More can be found in the mgo docs.

####NoResult(input error)
Checks the "error" returned from a "Find One" (mgo .One()) func to see if the error means no data was found for the query.

Ex.:
```go
data := struct{}
err  := db.DB(DATABASE).C(COLLECTION).Find(bson.M{"username":"test@test.com"}).One(&data)
if mongoutils.NoResult(err) {
	//no results
} else {
	//handle data found or other errors
}
```
####GetObjectIdFromString()
Gets the MongoDB BSON ObjectId representation of a hexidecimal string.  Returns an error if a string cannot be converted.

####GetStringFromObjectId()
Self explainatory.  Use the string representation of an ObjectId for user facing tasks.

####Limit()
Used if you are making API calls.  Grabs a limit from a form value variable, "limit".  Allows users to set limits instead of hardcoded on the back end.

####Sort()
Same as `Limit` above but for sort fields.  Only allows one field to be sorted.  Name of field can start with (-) to sort in decending order.
