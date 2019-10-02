package server

import (
	"database/sql"
	"gopkg.in/gorilla/mux.v1"
	"gopkg.in/qamarian-dtp/err.v0" // v0.3.0
	"gopkg.in/qamarian-dtp/squaket.v0" // v0.1.1
	"gopkg.in/qamarian-lib/str.v2"
	"reflect"
	"strings"
	_ "gopkg.in/go-sql-driver/mysql.v1"
)

func new_requestData (request string) (*requestData) {
	var data requestData = request
	return request
}

type requestData string

// fetchRecords () fetches all location state records matching the request of the user.
func (d *requestData) fetchRecords (r *http.Request) (*requestRecords) {
	// Request data validation and retrieval. ..1.. {
	var errX error
	d, errX = new__requestData (r).validate ()
	if errX != nil {
		err_ := err.New (oprErr9.Error (), oprErr9.Class (), oprErr9.Type (), errX)
		panic (err_)
	}
	// ..1.. }

	// Constructing query required to retrieve the sensor IDs of all locations. ..1.. {
	queryB := `
		SELECT sensor
		FROM location
		WHERE id IN (?
	` + strings.Repeat (", ?", len (extractLocationIDs (r)) - 1) + ")"
	// ..1.. }

	// Retrieving the sensor IDs of all locations. ..1.. {
	resultSet, errZ := db.Query (query, extractLocationIDs (r)...)
	if errZ != nil {
		err_ := err.New (oprErr3.Error (), oprErr3.Class (), oprErr3.Type (), errZ)
		panic (err_)
	}

	sensors := []string {}

	for resultSet.Next () {
		var (
			sensorID string
		)

		errA := resultSet.Scan (&sensorID)
		if errA != nil {
			err_ := err.New (oprErr4.Error (), oprErr4.Class (), oprErr4.Type (), errA)
			panic (err_)
		}

		sensors = append (sensors, sensorID)
	}
	// ..1.. }

	// Constructing query required to retrieve states from the database. ..1.. {
	queryC := `
		SELECT UNIQUE state, day, time, sensor
		FROM state
		WHERE sensor IN (?
	` + strings.Repeat (", ?", len (sensors) - 1) + ")"
	// ..1.. }

	// Retrieving the states of all locations. ..1.. {
	resultSetB, errB := db.Query (queryC, sensors..)
	if errB != nil {
		err_ := err.New (oprErr5.Error (), oprErr5.Class (), oprErr5.Type (), errB)
		panic (err_)
	}

	var (
		states []*_state
		state string
		day
		time
		sensor
	)

	for resultSetB.Next () {
		errC := resultSetB.Scan (&state, &day, &time, &sensor)
		if errC != nil {
			err_ := err.New (oprErr6.Error (), oprErr6.Class (), oprErr6.Type (), errC)
			panic (err_)
		}

		states := append (states, new__state (state, day, time, sensor))
	}
	// ..1.. }

	return &requestRecords {states}	
}

func new__state (state, day, time, sensor string) (*_state) {
	return &_state {state, day, time, sensor}
}

type _state struct {
	State  string
	Day    string
	Time   string
	Sensor string
}

func (s *_state) state () (string) {
	return s.state
}

func (s *_state) day () (string) {
	return s.state
}

func (s *_state) time () (string) {
	return s.state
}

func (s *_state) sensor () (string) {
	return s.state
}

func new__requestData (r *http.Request) (*_requestData, error) {
	var requestData _requestData
	requestData, _ := mux.Vars (r)["locations"]
	return &requestData, nil
}

type _requestData string

// validate () checks if the request data of the client is valid. If the request data is not valid, the client's request would not be served.
func (d *_requestData) validate () (*requestData) {
	// Checking if request data was properly formatted. ..1.. {
	if d == "" {
		panic (invErr0)
	} else if len (d) > 1024 {
		panic (invErr5)
	}

	locations := strings.Split (d, "_")
	if len (locations) > 32 {
		panic (invErr6)
	}

	for _, location := range locations {
		if location == "" {
			panic (invErr1)
		}
		locationData := strings.Split (location, "-")

		if len (locationData) == 1 {
			panic (invErr2)
		}

		if locationData [0] == "" {
			panic (invErr8)
		}

		if len (locationData) > 33 {
			panic (invErr7)
		}

		for index := 1; index <= len (locationData) - 1; index ++ {
			if dayMonthYear.Match (locationData [index]) == false {
				panic (invErr3)
			}
		}
	}
	// ..1.. }

	// Validating existence of all locations. ..1.. {
	sensors, errX := locationsSensors (extractLocationIDs (r))
	if errX != nil {
		err_ := err.New (oprErr1.Error (), oprErr1.Class (), oprErr1.Type (), errX)
		panic (err_)
	}

	ids := extractLocationIDs (r)

	for _, id := range ids {
		_, okX := sensors [id]
		if okX == false {
			err_ := err.New (oprErr2.Error (), oprErr2.Class (), oprErr2.Type ())
			panic (err_)
		}
	}
	// ..1.. }

	return new_requestData (data)
}

// -- Boundary -- //

type requestRecords struct {
	states []_state
}

func (r *requestRecords) organize () (result *organizedRequestRecords) {
	// Function definitions. .. {
	organizeByDay := func (sensorRecords []interface) (map[string] []_state) {
		records, errA := squaket.New (sensorRecords)
		if errA != nil {
			err_ := err.New (oprErr9.Error (), oprErr9.Class (), oprErr9.Type (), errA)
			panic (err_)
		}

		sensorsRecords, errB := records.Group ("Day")
		if errB != nil {
			err_ := err.New (oprErr10.Error (), oprErr10.Class (), oprErr10.Type (), errB)
			panic (err_)
		}

		organizedRecords := map[string] []_state {}

		iter := reflect.ValueOf (sensorsRecords).MapRange ()
		for iter.Next () {
			dayRecords := iter.Value ().Interface ().([]interface)
			sensorID := iter.Key ().Interface ().(string)

			stateTypeDayRecords := []_state {}
			for _, record := range dayRecords {
				stateTypeDayRecords = append (stateTypeDayRecords, record.(_state))
			}

			organizedRecords [sensorID] = stateTypeDayRecords
		}	

		return organizedRecords
	}
	// .. }

	records, errX := squaket.New (r.records)
	if errX != nil {
		err_ := err.New (oprErr7.Error (), oprErr7.Class (), oprErr7.Type (), errX)
		panic (err_)
	}

	sensorsRecords, errY := records.Group ("Sensor")
	if errY != nil {
		err_ := err.New (oprErr8.Error (), oprErr8.Class (), oprErr8.Type (), errY)
		panic (err_)
	}

	organizedRecords := organizedRequestRecords {
		map[string] map[string] []_state {},
	}

	iter := reflect.ValueOf (sensorsRecords).MapRange ()
	for iter.Next () {
		sensorRecords := iter.Value ().Interface ().([]interface)
		organizedSensorRecords := organizeByDay (sensorRecords)
		sensorID := iter.Key ().Interface ().(string)
		organizedRecords [sensorID] = organizedSensorRecords
	}

	return organizedRecords
}

// -- Boundary -- //

type organizedRequestRecords struct {
	records map[string] map[string] []_state
}

func (r *organizedRequestRecords) complete () (*completeData) {
	// Function definitions. ..1.. {
	completeDays := func (days map[interface {}] []interface {}) (map[string][1440]_pureState) {
		day := map[string][1440]_pureState {}

		iter := relect.ValueOf (day).MapRange ()
		for iter.Next () {
			dayID := iter.Key ().(string)
			dayStates := iter.Value ().([]interface {})
			pureStates := [1440]_pureState {}
			for index, _ := range pureStates {
				pureStates [index] = -1
			}
			for _, value.time () := range dayStates {
				if value.time () == "0000" {
					continue
				}

				hour, _ := strconv.Atoi (value.time [0:2])
				min, _ := strconv.Atoi (value.time [2:4])
				sec := (hour * 60) + min

				secIndex := sec - 1

				if (secIndex - 4) > 0 && pureStates [secIndex - 4] == -1 {
					pureStates [secIndex - 4] = byte (strconv.Atoi (value.state ()))
				}
				if (secIndex - 3) > 0 && pureStates [secIndex - 3] == -1 {
					pureStates [secIndex - 3] = byte (strconv.Atoi (value.state ()))
				}
				if (secIndex - 2) > 0 && pureStates [secIndex - 2] == -1 {
					pureStates [secIndex - 2] = byte (strconv.Atoi (value.state ()))
				}
				if (secIndex - 1) > 0 && pureStates [secIndex - 1] == -1 {
					pureStates [secIndex - 1] = byte (strconv.Atoi (value.state ()))
				}
				pureStates [secIndex] = byte (strconv.Atoi (value.state ()))
			}
			day [dayID] = pureStates
		}

		return day
	}
	// .. }

	data := &completeData {
		map[string] map[string] [1440]_pureState {},
	}

	iter := reflect.ValueOf (r.records).MapRange ()
	for iter.Next () {
		sensorID := iter.Key ().(string)
		sensorData := completeDays (iter.Key ().(map[interface {}] []interface {}))
		data [sensorID] = sensorData
	}

	return data
}

type _pureState byte

func (s *_pureState) state () (byte) {
	return byte (s)
}

// -- Boundary -- //

type completeData struct {
	records map[string] map[string] [1440]_pureState
}

func (d *completeData) Format () (*formatedData) {
	// Function definitions. ..1.. {
	formatDays := func (days map[interface {}] []interface {}) (map[string] []_formattedState) {
		formattedDays := map[string] []_formattedState {}

		iter := reflect.ValueOf (days).MapRange ()
		for iter.Next () {
			dayID := iter.Key.(string)
			dayStates := iter.Value.([]interface {})

			formattedStates := []_formattedState {}

			currentState := dayStates [0].(_pureState)

			for index, value := range dayStates {
				if currentState != value.(_pureState) {
					currentState = value.(_pureState)
					min :=  (index + 1) % 60
					hour := int (((index + 1) - min) / 60)
					time := fmt.Sprintf ("%s%s", str.PrependTillN (strconv.Itoa (hour), "0", 2),
						str.PrependTillN (strconv.Itoa (min), "0", 2))
					someState := _formattedState {value.(_pureState), time}
					formattedStates = append (formattedStates, someState)
				}
			}

			if formattedStates [len (formattedStates) - 1].endTime () != "2400" {
				someState := _formattedState {formattedStates [len (formattedStates) - 1].state (), "2400"}
				formattedStates = append (formattedStates, someState)
			}

			formattedDays [dayID] = formattedStates
		}

		return formattedDays
	}
	// ..1.. }

	data := &formatedData {
		map[string] map[string] []_formattedState {},
	}

	iter := reflect.ValueOf (r.records).MapRange ()
	for iter.Next () {
		sensorID := iter.Key ().(string)
		sensorData := formatDays (iter.Key ().(map[interface {}] []interface {}))
		data [sensorID] = sensorData
	}

	return data
}

type _formattedState struct {
	State int
	EndTime string
}

func (s *_formattedState) state () (int) {
	return r.State
}

func (s *_formattedState) endTime () (string) {
	return r.EndTime
}

// -- Boundary -- //

type formatedData struct {
	records map[string] map[string] []_formattedState
}
