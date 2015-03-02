package pg

import (
	"fmt"
	"reflect"
)

type Entity interface {
	Result
	Table() string
}

type EntityHandler interface {
	Create(entity Entity) error
	FindOne(entity Entity, where string, whereParams ...interface{}) (Entity, error)
	FindAll(entity Entity, where string, whereParams ...interface{}) ([]Entity, error)
	Update(entity Entity) error
	Delete(entity Entity) error
}

func NewEntity(e Entity) Entity {
	mustBeStructPtr(e)
	entityType := reflect.TypeOf(e).Elem()
	newEntity := reflect.New(entityType)
	return newEntity.Interface().(Entity)
}

func mustBeStructPtr(e interface{}) {
	v := reflect.ValueOf(e)
	t := v.Type()
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		panic(fmt.Errorf("entity must be pointer to struct; got %T", v))
	}
}
