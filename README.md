# mongoutils.go
Helper Functions for using MongoDB and Wrapper Functions Around the mgo Driver

This package is intended to be used to connect to a Mongo Database and simply some functionality. It is not meant for power users or for production since there is no authenticaiton and configuration is done by editting this file.  But then again, power users would make their own wrapper functions around mgo.

---

###Usage

This library will connect to your MongoDB(s) and then store the session connection in a global variable "MGO_SESSION".  Anytime you want to use your MongoDB, you will need to include this file and copy the MGO_SESSION.  Copying the session results in using a connection from the "pool" instead of relying on a single connections which could create latency and slow performance.

There is some required setup for your environment. Please see below.

###Setup

To use this library, the user needs to do a bit of setup first.  This is done so all MongoDB configuration is in one place instead of in separate files.

- SERVERS:
	- List of server(s) to connect to. If you are using a replica set, provide as many servers as possible.
	- Single host: "localhost:27017/"
	- Replica Set: "localhost:27017, host2:27017, host3:27017/"

- DATABASE:
	- Simple, the name of the database to connect to.

- COLL(ections):
	- A list of collections stores as constants.  This makes maintaining your collection names easier since changes only occur in one place.
	- Example: COLL_USERS = "users"

- READ_PREFERENCE:
	- This is a more advanced feature.  You can read about it in the MongoDB and mgo docs.
	- Basically, do all reads go to the "master" or are the reads spread across all servers available.
	- Default: Monotonic aka reads are spread over slaves.

- WRITE_CONCERN:
	- Another more advanced feature.
	- Basically, when do you consider a write successful? When it is aknowledged, written to disk, written to many disks?
	- Default: Majority & Fsync aka a majority of the available servers must have written the write to disk.

