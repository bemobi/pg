package pg

import "database/sql"

// Tx implements the EntityHandler interface
type Tx struct {
	tx *sql.Tx
}

func (t *Tx) Commit() error {
	return t.tx.Commit()
}

func (t *Tx) Rollback() error {
	return t.tx.Rollback()
}

func (t *Tx) Create(entity Entity) error {
	return create(t.tx, entity)
}

func (t *Tx) FindOne(entity Entity, where string, whereParams ...interface{}) (Entity, error) {
	return findOne(t.tx, entity, where, whereParams...)
}

func (t *Tx) FindAll(entity Entity, where string, whereParams ...interface{}) ([]Entity, error) {
	return findAll(t.tx, entity, where, whereParams...)
}

func (t *Tx) Update(entity Entity) error {
	return update(t.tx, entity)
}

func (t *Tx) Delete(entity Entity) error {
	return delete(t.tx, entity)
}
