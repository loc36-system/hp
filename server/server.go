package server

import (
	"database/sql"
	"errors"
	"fmt"
	"gopkg.in/gorilla/mux.v1"
	"gopkg.in/qamarian-dtp/err.v0" // v0.1.1
	"gopkg.in/qamarian-lib/str.v2"
	"net/http"
	"strings"
	_ "gopkg.in/go-sql-driver/mysql.v1"
)

func init () {
	var errX error
	dayMonthYear, errX := regexp.Compile ("^20\d{2}(0[1-9]|1[0-2)(0[1-9]|[1-2]\d|3[0-1])$")
	if errX != nil {
		str.PrintEtr ("Regular expression compilation failed.", "err", "server.init ()")
		panic ("Regular expression compilation failed.")
	}
}

func serviceRequestServer (w http.ResponseWriter, r *http.Request) {
	defer func () { // Error handling.
		//
	}

	// Fetch data from database. ...1...  {

	// Data retrieval. ...2... {
	request, okX := mex.Vars ()["locations"]
	if okX == false || request == "" {
		panic (err.New ("No request data was provided.", 1, 1))
	}
	// ...2... }

	// Data validation. ...2... {
	locations := strings.Split (request, "_")
	for _, location := range locations {
		if location == "" {
			panic (err.New ("A location data was omited.", 1, 2))
		}
		data := strings.Split (location, "-")
		if len (data) == 1 {
			errMssg := fmt.Sprintf ("No day was specified for location '%s'.", data [0])
			panic (err.New (errMssg, 1, 3))
		}
		for index := 1; index <= len (data) - 1; index ++ {
			if dayMonthYear.Match (data [index]) == false {
				errMssg := fmt.Sprintf ("Data of invalid day requested for location '%s'.", data [0])
				panic (err.New (errMssg, 1, 4))
			}
		}
	}
	if len (locations) > 256 {
		panic (err.New ("Data of more than 256 locations may not be requested at a time."), 1, 5)
	}
	// ...2... }

	// Validating existence of all locations. ...2... {
	// Constructing query required for validation. ...3... { ...
	query := `
		SELECT COUNT (id)
		FROM location
		WHERE id IN (?
	`)
	for index := 1; index <= len (locations) - 1, index ++ {
		query = query + ", ?"
	}
	query = query + ")"
	// ...3... }

	errX := db.Ping ()
	if errX != nil {
		panic (err.New ("Database unreachable.", 2, 1, errX))
	}

	var noOfValidLocations int
	errY := db.QueryRow (query, locations ...).Scan (&noOfValidLocations)
	if errY != nil {
		panic (err.New ("Unable to check if all locations exist.", 2, 2, errY))
	}

	if noOfValidLocations != len (locations) {
		panic (err.New ("Some locations do not exist.", 2, 3))
	}
	// ...2... }

	// Querying data. ...2... {
	// Constructing query required to retrieve the sensor IDs of all locations. ...3... {
	queryB := `
		SELECT id, sensor
		FROM location
		WHERE id IN (?
	`)
	for index := 1; index <= len (locations) - 1, index ++ {
		queryB = queryB + ", ?"
	}
	queryB = queryB + ")"
	// ...3... }

	// Retrieving the sensor IDs of all locations. ...3... {
	var (
		locationIDs []string
		sensors []string
	)
	resultSet, errZ := db.Query (query, locations ...)
	if errZ != nil {
		panic (err.New ("Unable to fetch the sensor IDs of all locations.", 2, 4, errZ))
	}
	for resultSet.Next () {
		var (
			locationID string
			sensorID string
		)
		errA := resultSet.Scan (&locationID, &sensorID)
		if errA != nil {
			panic (err.New ("Unable to fetch the sensor ID of a location.", 2, 5, errA))
		}
		sensors = append (sensors, sensorID)
	}
	// ...3... }

	// Constructing query required to retrieve states from the database. ...3... {
	queryC := `
		SELECT UNIQUE state, day, time, sensor
		FROM state
		WHERE sensor IN (?
	`
	for index := 1; index <= len (sensors) - 1, index ++ {
		queryC = queryC + ", ?"
	}
	queryC = queryC + ")"
	// ...3... }

	// Retrieving the states of all locations. { ...3...
	var (
		states []state
	)
	resultSetB, errB := db.Query (queryC, sensors)
	if errB != nil {
		panic (err.New ("Unable to fetch the states of all locations.", 2, 6, errB))
	}
	for resultSetB.Next () {
		someState := state {}
		errC := resultSetB.Scan (&someState.state, &someState.day, &someState.time, &someState.sensor)
		if errC != nil {
			panic (err.New ("Unable to fetch a state data.", 2, 7, errC))
		}
		states := append (states, someState)
	}
	// ...3... }

	// ...2... }

	// ...1... }

	// Present fetched data. ...1... {}

	// Send data to user. ...1... {}
}

var (
	dayMonthYear *Regexp
	db *sql.DB
)

type state struct {
	state  string
	day    string
	time   string
	sensor string
}
