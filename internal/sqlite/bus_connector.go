package sqlite

import (
	"database/sql"
	"fmt"
	"math"

	"github.com/eryalito/vigo-bus-core/internal/config"
	"github.com/eryalito/vigo-bus-core/pkg/api"

	_ "github.com/mattn/go-sqlite3"
)

// BusConnector is a struct that holds the database connection
type BusConnector struct {
	DB *sql.DB
}

// NewBusConnector initializes a new database given a path
func NewBusConnector() (*BusConnector, error) {
	db, err := sql.Open("sqlite3", config.StopsDBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	connector := &BusConnector{DB: db}
	if err := connector.initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize database: %v", err)
	}

	return connector, nil
}

// initialize sets up the initial database schema
func (c *BusConnector) initialize() error {
	createLinesTable := `
    CREATE TABLE IF NOT EXISTS lines (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL
    );`

	createStopsTable := `
    CREATE TABLE IF NOT EXISTS stops (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        stop_number INTEGER NOT NULL,
        stop_id INTEGER NOT NULL,
        name TEXT NOT NULL,
        lat REAL,
        lon REAL
    );`

	createLineStopsTable := `
    CREATE TABLE IF NOT EXISTS line_stops (
        line_id INTEGER,
        stop_id INTEGER,
        PRIMARY KEY (line_id, stop_id),
        FOREIGN KEY (line_id) REFERENCES lines(id),
        FOREIGN KEY (stop_id) REFERENCES stops(id)
    );`

	_, err := c.DB.Exec(createLinesTable)
	if err != nil {
		return fmt.Errorf("failed to create lines table: %v", err)
	}

	_, err = c.DB.Exec(createStopsTable)
	if err != nil {
		return fmt.Errorf("failed to create stops table: %v", err)
	}

	_, err = c.DB.Exec(createLineStopsTable)
	if err != nil {
		return fmt.Errorf("failed to create line_stops table: %v", err)
	}

	return nil
}

// InsertLine inserts a new line into the lines table
func (c *BusConnector) InsertLine(name string) (int64, error) {
	insertQuery := `INSERT INTO lines (name) VALUES (?)`
	result, err := c.DB.Exec(insertQuery, name)
	if err != nil {
		return 0, fmt.Errorf("failed to insert line: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %v", err)
	}

	return id, nil
}

// InsertStop inserts a new stop into the stops table
func (c *BusConnector) InsertStop(stop api.Stop) (int64, error) {
	insertQuery := `INSERT INTO stops (stop_number, stop_id, name, lat, lon) VALUES (?, ?, ?, ?, ?, ?)`
	result, err := c.DB.Exec(insertQuery, stop.StopNumber, stop.StopID, stop.Name, stop.Location.Lat, stop.Location.Lon)
	if err != nil {
		return 0, fmt.Errorf("failed to insert stop: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %v", err)
	}

	return id, nil
}

// AddStopToLine adds a stop to a line in the line_stops table
func (c *BusConnector) AddStopToLine(lineID, stopID int) error {
	insertQuery := `INSERT INTO line_stops (line_id, stop_id) VALUES (?, ?)`
	_, err := c.DB.Exec(insertQuery, lineID, stopID)
	if err != nil {
		return fmt.Errorf("failed to add stop to line: %v", err)
	}

	return nil
}

// GetLines retrieves all lines from the lines table
func (c *BusConnector) GetLines() ([]api.Line, error) {
	query := `SELECT id, name FROM lines`
	rows, err := c.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query lines: %v", err)
	}
	defer rows.Close()

	var lines []api.Line
	for rows.Next() {
		var line api.Line
		if err := rows.Scan(&line.ID, &line.Name); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		lines = append(lines, line)
	}

	return lines, nil
}

// GetLineByName retrieves a line from the lines table by name
func (c *BusConnector) GetLineByName(name string) (api.Line, error) {
	query := `SELECT id, name FROM lines WHERE name = ?`
	row := c.DB.QueryRow(query, name)

	var line api.Line
	if err := row.Scan(&line.ID, &line.Name); err != nil {
		return api.Line{}, fmt.Errorf("failed to scan row: %v", err)
	}

	return line, nil
}

// GetStops retrieves all stops from the stops table
func (c *BusConnector) GetStops() ([]api.Stop, error) {
	query := `SELECT id, stop_number, stop_id, name, lat, lon FROM stops`
	rows, err := c.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query stops: %v", err)
	}
	defer rows.Close()

	var stops []api.Stop
	for rows.Next() {
		var stop api.Stop
		if err := rows.Scan(&stop.ID, &stop.StopNumber, &stop.StopID, &stop.Name, &stop.Location.Lat, &stop.Location.Lon); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		stops = append(stops, stop)
	}

	return stops, nil
}

// GetStopByNumber retrieves a stop from the stops table by stop number
func (c *BusConnector) GetStopByNumber(stopNumber int) (api.Stop, error) {
	query := `SELECT id, stop_number, stop_id, name, lat, lon FROM stops WHERE stop_number = ?`
	row := c.DB.QueryRow(query, stopNumber)

	var stop api.Stop
	if err := row.Scan(&stop.ID, &stop.StopNumber, &stop.StopID, &stop.Name, &stop.Location.Lat, &stop.Location.Lon); err != nil {
		return api.Stop{}, fmt.Errorf("failed to scan row: %v", err)
	}

	return stop, nil
}

// FindStopsByText retrieves a stop from the stops table by text matching a stop name
func (c *BusConnector) FindStopsByText(text string) ([]api.Stop, error) {
	query := `SELECT id, stop_number, stop_id, name, lat, lon FROM stops WHERE name LIKE ?`
	rows, err := c.DB.Query(query, "%"+text+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to query stops: %v", err)
	}
	defer rows.Close()

	var stops []api.Stop
	for rows.Next() {
		var stop api.Stop
		if err := rows.Scan(&stop.ID, &stop.StopNumber, &stop.StopID, &stop.Name, &stop.Location.Lat, &stop.Location.Lon); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		stops = append(stops, stop)
	}

	return stops, nil
}

// FindStopsByLocation retrieves stops from the stops table within a given radius around a location in meters
func (c *BusConnector) FindStopsByLocation(lat, lon, radius float64) ([]api.Stop, error) {
	query := `SELECT id, stop_number, stop_id, name, lat, lon FROM stops`
	rows, err := c.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query stops: %v", err)
	}
	defer rows.Close()

	var stops []api.Stop
	for rows.Next() {
		var stop api.Stop
		if err := rows.Scan(&stop.ID, &stop.StopNumber, &stop.StopID, &stop.Name, &stop.Location.Lat, &stop.Location.Lon); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		// Calculate the distance using the Haversine formula
		distance := haversine(lat, lon, stop.Location.Lat, stop.Location.Lon)
		if distance <= radius {
			stops = append(stops, stop)
		}
	}

	return stops, nil
}

// Close closes the database connection
func (c *BusConnector) Close() error {
	if err := c.DB.Close(); err != nil {
		return fmt.Errorf("failed to close database: %v", err)
	}
	return nil
}

// Haversine formula to calculate the distance between two points in meters
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371000 // Radius of the Earth in meters
	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	dlat := lat2Rad - lat1Rad
	dlon := lon2Rad - lon1Rad

	a := math.Sin(dlat/2)*math.Sin(dlat/2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}
