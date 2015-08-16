/*
Package mongoutils is used to simplify some usage of the mgo MongoDB driver.

Use this package to connect, disconnect, host mgo session data (global variable), and perform some simple functions.
Basically, this library helps clean up your code base elsewhere.

This package uses and requires the mgo driver for MongoDB. No other drivers are supported.

When connecting to a MongoDB, this library will store the connection data in a global variable saved in this file.
Include this file wherever you need to use your DB.
However, you must copy the session (per mgo documents) in order to use different connections to the database (aka pooling connections instead of using only one connection).

It is highly suggested that you create another file in your project for storing your servers, database, and collection names as constants.
The result is easier maintenance of your project since these values are saved in one location for easy editing.
You can also store your global session data in this file instead of relying on the global variable below.

Note: this package is not meant to meant for production environments.
*/

package mongoutils

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	ID_LENGTH           = 24
	LIMIT_DEFAULT_VALUE = 5
	LIMIT_RETURN_ALL    = 0
	SORT_DEFAULT        = "_id"
)

var (
	//GLOBAL SESSION DATA
	SESSION *mgo.Session

	//ERROR MESSAGES
	ErrIdBadLength = errors.New("idMustBe24CharactersLong")
	ErrIdNotHex    = errors.New("idNotHexadecimal")
	ErrNoResults   = errors.New("noResultsFound")
)

//*********************************************************************************************************************************
//CONNECT & DISCONNECT

//CONNECT TO DB AND SETUP WRITE CONCERN AND READ PREFERENCE
//servers is string of at least one server to connect to ("localhost:27017/") and should be many hosts if using a replica set
//database is the name of your database to connect to
//readPreference is an mgo consistency constant (Eventual, Monotonic, Strong)
//writeConcern is an mgo *Safe type
//saves the connected session pool to a global variable.
func Connect(servers string, database string, readPreference int, writeConcern *mgo.Safe) {
	//connection uri
	uri := servers + database

	//connect to db
	session, err := mgo.Dial(uri)
	if err != nil {
		log.Println("mongoutils.go-Connect error")
		log.Panicln(err)
		return
	}

	//set db consistency
	//read preference
	setReadPreference(session, readPreference)

	//set safety mode
	//write concern
	session.SetSafe(writeConcern)

	//store session in global variable
	//access this session by importing this file
	log.Println("mongoutils.go-Connect okay")
	SESSION = session
	return
}

//CLOSE THE DB CONNECTION
//using the mgo Close() function
func Close() {
	SESSION.Close()
	return
}

//SET READ PREFERENCE
//0, 1, 2 are from mgo documents
//need to do this since mgo.mode is not exported
func setReadPreference(s *mgo.Session, input int) {
	if input == 0 {
		s.SetMode(mgo.Eventual, true)
	} else if input == 1 {
		s.SetMode(mgo.Monotonic, true)
	} else if input == 2 {
		s.SetMode(mgo.Strong, true)
	} else {
		log.Println("mongoutils.go-SetReadPreference error")
	}

	return
}

//*********************************************************************************************************************************
//ERROR HANDLING

//CHECK IF A FIND ONE RETURNED NO RESULTS
func NoResult(input error) (bool, error) {
	if input == mgo.ErrNotFound {
		return true, ErrNoResults
	}

	return false, nil
}

//*********************************************************************************************************************************
//OBJECT ID VALIDATION

//CHECK IF AN ID IS THE CORRECT LENGTH
//mongo ids are exactly 24 characters
//in: objectId as a string
//out: boolean and error if the input is not exactly 24 characters long
func isIdCorrectLength(inputId string) (bool, error) {
	if len(inputId) != ID_LENGTH {
		return false, ErrIdBadLength
	}

	return true, nil
}

//CHECK IF A STRING CAN BE A VALID MONGO ID
//mongo ids are hexidecimal characters
//in: objectId as a string
//out: boolean and error if the input is not hexidecimal
func isValidHexString(inputId string) (bool, error) {
	if bson.IsObjectIdHex(inputId) == false {
		return false, ErrIdNotHex
	}

	return true, nil
}

//CHECK IF ID IS VALID
//wrapper around the above functions
//in: objectId as a string
//out: boolean and error if the input is not a valid string representation of an objectId
func isValidId(inputId string) (bool, error) {
	if yes, err := isIdCorrectLength(inputId); yes == false {
		return false, err
	}
	if yes, err := isValidHexString(inputId); yes == false {
		return false, err
	}

	return true, nil
}

//*********************************************************************************************************************************
//OBJECT ID CONVERSION

//CONVERT A STRING INTO AN OBJECT ID
//validates the input string first and returns an error if input is not a valid string to convert
//in: objectId as a string
//out: mongo objectId and error if the input is not valid
func GetObjectIdFromString(inputId string) (bson.ObjectId, error) {
	//validate input
	if yes, err := isValidId(inputId); yes == false {
		return bson.NewObjectId(), err
	}

	return bson.ObjectIdHex(inputId), nil
}

//CONVERT AN OBJECT ID INTO A STRING
//in: mongo objectId
//out: string exactly 24 characters long and hexidecimal
func GetStringFromObjectId(input bson.ObjectId) string {
	return input.Hex()
}

//*********************************************************************************************************************************
//QUERIES

//GET A LIMIT FOR NUMBER FOR RESULTS TO RETURN FROM GET VARIABLE
//return the limit as an integer to use in db query
//5 is the default if the limit form value is not understood
//a limit of 0 (zero) actually returns all results, not none
//gets the limit value from an http GET form value i.e. example.com?limit=10
func Limit(r *http.Request) int {
	//get value from get variable
	limit := r.FormValue("limit")

	//if no limit was set in form value, set limit to default
	if len(limit) == 0 {
		return LIMIT_DEFAULT_VALUE
	}

	//if limit was set to a keyword, return all docs
	if limit == "none" || limit == "all" {
		return LIMIT_RETURN_ALL
	}

	//limit was given as a number in form value
	//convert form value to int
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return LIMIT_DEFAULT_VALUE
	}

	//no error, limit was an integer
	return limitInt
}

//GET A FIELD TO SORT FIND RESULTS BY FROM GET VARIABLE
//you can sort by one or many fields, each field name separated by a comma without whitespace
//you can prepend a (-) minus sign to sort in decending order
//example.log/?sort=birthday,-username
//make sure to use the value this function returns as sortOrder... (note three periods) in mgo Sort()
//this way mgo will apply all sorts to your query
func Sort(r *http.Request) []string {
	//parse the request form values
	r.ParseForm()
	sort := r.FormValue("sort")

	//check if there is a value set
	if len(sort) == 0 {
		return []string{SORT_DEFAULT}
	}

	//split the results to sort by many fields
	split := strings.Split(sort, ",")

	return split
}
