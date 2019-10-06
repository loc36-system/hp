package server

import (
	"database/sql"
	"gopkg.in/gorilla/mux.v1"
	"gopkg.in/qamarian-dtp/err.v0" // v0.4.0
	"gopkg.in/qamarian-lib/str.v2" // v2.0.0
	"math/big"
	"net/http"
	"regexp"
	"strings"
)

// -- Boundary -- //

func init () {
	var errX error
	dayMonthYear, errX = regexp.Compile (`^20\d{2}(0[1-9]|1[0-2)(0[1-9]|[1-2]\d|3[0-1])$`)
	if errX != nil {
		str.PrintEtr ("Regular expression compilation failed.", "err", "init ()")
		panic ("Regular expression compilation failed.")
	}
}

var (
	dayMonthYear *regexp.Regexp // Cache
	db *sql.DB           // Cache
)

// -- Boundary -- //

func extractLocationIDs (r *http.Request) ([]string) {
	output := []string {}

	data, _ := mux.Vars (r)["locations"]
	locationsData := strings.Split (data, "_")

	for _, locationData := range locationsData {
		segments := strings.Split (locationData, "-")
		output = append (output, segments [0])
	}

	return output
}

// -- Boundary -- //

func locationsSensors (locations []string) (*_locationsSensors, error) {
	query := `SELECT id, sensor
	FROM location
	WHERE id IN (?` + strings.Repeat (", ?", len (locations) - 1) + ")"

	errX := db.Ping ()
	if errX != nil {
		errMssg := "Database unreachable."
		return nil, err.New (errMssg, big.NewInt (0), big.NewInt (0), errX)
	}

	var (
		location string
		sensor string
	)

	result, errY := db.Query (query, locations...)
	if errY != nil {
		errMssg := "Unable to successfully query database for locations sensors."
		return nil, err.New (errMssg, big.NewInt (0), big.NewInt (0), errY)
	}

	output := _locationsSensors {}

	for result.Next () {
		errZ := result.Scan (&location, &sensor)
		if errZ != nil {
			errMssg := "Unable to fetch the sensor of a location."
			return nil, err.New (errMssg, big.NewInt (0), big.NewInt (0), errZ)
		}

		output.add (location, sensor)
	}

	return output, nil
}

type _locationsSensors map[string]string

func (l *_locationsSensors) add (location, sensor string) {
	l [location] = sensor
}

func (l *_locationsSensors) getLocationSensor (location string) (string, bool) {
	return l [location]
}

// -- Boundary -- //
