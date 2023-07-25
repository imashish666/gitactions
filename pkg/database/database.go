package database

import "github.com/jmoiron/sqlx"

type DatabaseOps interface {
	Select(query string, data interface{}, args ...interface{}) error
	Insert(query string, args ...interface{}) error
	Get(query string, data interface{}, args ...interface{}) error
}

type Database struct {
	DB *sqlx.DB
}

// NewDatabase returns an instance of Database struct
func NewDatabase(db *sqlx.DB) Database {
	return Database{DB: db}
}

// Select is used for fetching data from db
func (m Database) Select(query string, data interface{}, args ...interface{}) error {
	return m.DB.Select(data, query, args...)
}

// Insert is used for adding data to db
func (m Database) Insert(query string, args ...interface{}) error {
	_, err := m.DB.Exec(query, args)
	return err
}

// Get is used for fetching specific data from db
func (m Database) Get(query string, data interface{}, args ...interface{}) error {
	return m.DB.Get(data, query, args...)
}
