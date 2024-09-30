package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/eryalito/vigo-bus-core/internal/config"
	"github.com/eryalito/vigo-bus-core/pkg/api"

	_ "github.com/mattn/go-sqlite3"
)

type IdentityConnector struct {
	DB *sql.DB
}

// NewIdentityConnector creates a new IdentityConnector and initializes the database
func NewIdentityConnector() (*IdentityConnector, error) {
	db, err := sql.Open("sqlite3", config.IdentityDBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	connector := &IdentityConnector{DB: db}
	if err := connector.createTables(); err != nil {
		return nil, err
	}

	return connector, nil
}

// createTables creates the identities and favorite_stops tables if they don't exist
func (c *IdentityConnector) createTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS identities (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            uuid TEXT NOT NULL,
            provider TEXT NOT NULL,
            metadata TEXT
        );`,
		`CREATE TABLE IF NOT EXISTS favorite_stops (
            identity_id INTEGER,
            stop_number INTEGER,
            FOREIGN KEY(identity_id) REFERENCES identities(id)
        );`,
	}

	for _, query := range queries {
		if _, err := c.DB.Exec(query); err != nil {
			return fmt.Errorf("failed to create table: %v", err)
		}
	}
	return nil
}

// InsertIdentity inserts a new identity into the database
func (c *IdentityConnector) InsertIdentity(identity *api.Identity) error {
	tx, err := c.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	query := `INSERT INTO identities (metadata, uuid, provider) VALUES (?, ?, ?)`
	result, err := tx.Exec(query, identity.Metadata, identity.UUID, identity.Provider)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert identity: %v", err)
	}

	identityID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get last insert id: %v", err)
	}

	for _, stop := range identity.FavoriteStops {
		query := `INSERT INTO favorite_stops (identity_id, stop_number) VALUES (?, ?)`
		if _, err := tx.Exec(query, identityID, stop.StopNumber); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert favorite stop: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

// GetIdentity retrieves an identity by ID
func (c *IdentityConnector) GetIdentity(id int) (*api.Identity, error) {
	query := `SELECT id, metadata, uuid, provider FROM identities WHERE id = ?`
	row := c.DB.QueryRow(query, id)

	var identity api.Identity
	if err := row.Scan(&identity.ID, &identity.Metadata, &identity.UUID, &identity.Provider); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No identity found
		}
		return nil, fmt.Errorf("failed to get identity: %v", err)
	}

	query = `SELECT stop_number FROM favorite_stops WHERE identity_id = ?`
	rows, err := c.DB.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get favorite stops: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var stop api.Stop
		if err := rows.Scan(&stop.StopNumber); err != nil {
			return nil, fmt.Errorf("failed to scan favorite stop: %v", err)
		}
		identity.FavoriteStops = append(identity.FavoriteStops, stop)
	}

	return &identity, nil
}

// UpdateIdentity updates an existing identity in the database
func (c *IdentityConnector) UpdateIdentity(identity *api.Identity) error {
	tx, err := c.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	query := `UPDATE identities SET metadata = ?, uuid = ?, provider = ? WHERE id = ?`
	if _, err := tx.Exec(query, identity.Metadata, identity.UUID, identity.Provider, identity.ID); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update identity: %v", err)
	}

	query = `DELETE FROM favorite_stops WHERE identity_id = ?`
	if _, err := tx.Exec(query, identity.ID); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete favorite stops: %v", err)
	}

	for _, stop := range identity.FavoriteStops {
		query := `INSERT INTO favorite_stops (identity_id, stop_number) VALUES (?, ?)`
		if _, err := tx.Exec(query, identity.ID, stop.StopNumber); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert favorite stop: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

// DeleteIdentity deletes an identity by ID
func (c *IdentityConnector) DeleteIdentity(id int) error {
	tx, err := c.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	query := `DELETE FROM favorite_stops WHERE identity_id = ?`
	if _, err := tx.Exec(query, id); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete favorite stops: %v", err)
	}

	query = `DELETE FROM identities WHERE id = ?`
	if _, err := tx.Exec(query, id); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete identity: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

// GetUserByUUID retrieves an identity by UUID and provider
func (c *IdentityConnector) GetUserByUUID(provider, uuid string) (*api.Identity, error) {
	query := `SELECT id FROM identities WHERE uuid = ? AND provider = ?`
	row := c.DB.QueryRow(query, uuid, provider)

	var id int
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No identity found
		}
		return nil, fmt.Errorf("failed to get identity: %v", err)
	}

	return c.GetIdentity(id)
}

// Close closes the database connection
func (c *IdentityConnector) Close() error {
	if err := c.DB.Close(); err != nil {
		return fmt.Errorf("failed to close database: %v", err)
	}
	return nil
}
