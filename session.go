package pg

import (
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
)

var (
	ERecordNotFound  = errors.New("Record not found")
	EMultipleResults = errors.New("Unexpected multiple results from query")
)

// Session implements the EntityHandler and the ResultHandler interface
type Session struct {
	DB *sql.DB
}

func NewSession(connectionString string) (*Session, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	return &Session{db}, nil
}

func (s *Session) Create(entity Entity) error {
	return create(s.DB, entity)
}

func (s *Session) FindOne(entity Entity, where string, whereParams ...interface{}) (Entity, error) {
	return findOne(s.DB, entity, where, whereParams...)
}

func (s *Session) FindAll(entity Entity, where string, whereParams ...interface{}) ([]Entity, error) {
	return findAll(s.DB, entity, where, whereParams...)
}

func (s *Session) Update(entity Entity) error {
	return update(s.DB, entity)
}

func (s *Session) Delete(entity Entity) error {
	return delete(s.DB, entity)
}

func (s *Session) Query(result Result, sql string, params ...interface{}) ([]Result, error) {
	return query(s.DB, result, sql, params...)
}

func (s *Session) Tx() (*Tx, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	return &Tx{tx}, nil
}
