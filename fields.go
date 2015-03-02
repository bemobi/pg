package pg

import (
	"database/sql"
	"reflect"
)

type fields struct {
	Result
}

func (f fields) extract() []reflect.StructField {
	mustBeStructPtr(f.Result)
	entityType := reflect.TypeOf(f.Result).Elem()
	list := make([]reflect.StructField, 0)
	for i := 0; i < entityType.NumField(); i++ {
		field := entityType.Field(i)
		t := newTag(field)
		if !t.omitfield {
			list = append(list, field)
		}
	}
	return list
}

func (f fields) all() []string {
	list := make([]string, 0)
	for _, field := range f.extract() {
		t := newTag(field)
		list = append(list, t.fieldname)
	}
	return list
}

func (f fields) autofilled(autofilled bool) []string {
	list := make([]string, 0)
	for _, field := range f.extract() {
		t := newTag(field)
		if t.autofilled == autofilled {
			list = append(list, t.fieldname)
		}
	}
	return list
}

func (f fields) autofilledValues(autofilled bool) []interface{} {
	entityValue := reflect.ValueOf(f.Result).Elem()
	values := make([]interface{}, 0)
	for _, field := range f.extract() {
		t := newTag(field)
		if t.autofilled == autofilled {
			value := entityValue.FieldByName(field.Name).Addr().Interface()
			values = append(values, value)
		}
	}
	return values
}

func (f fields) pk(pk bool) []string {
	list := make([]string, 0)
	for _, field := range f.extract() {
		t := newTag(field)
		if t.pk == pk {
			list = append(list, t.fieldname)
		}
	}
	return list
}

func (f fields) pkValues(pk bool) []interface{} {
	entityValue := reflect.ValueOf(f.Result).Elem()
	values := make([]interface{}, 0)
	for _, field := range f.extract() {
		t := newTag(field)
		if t.pk == pk {
			value := entityValue.FieldByName(field.Name).Addr().Interface()
			values = append(values, value)
		}
	}
	return values
}

func (f fields) updatable() []string {
	return f.autofilled(false)
}

func (f fields) notUpdatable() []string {
	return f.autofilled(true)
}

func (f fields) values() []interface{} {
	entityValue := reflect.ValueOf(f.Result).Elem()
	list := make([]interface{}, 0)
	for _, field := range f.extract() {
		value := entityValue.FieldByName(field.Name).Addr().Interface()
		list = append(list, value)
	}
	return list
}

func (f fields) updatableValues() []interface{} {
	return f.autofilledValues(false)
}

func (f fields) notUpdatableValues() []interface{} {
	return f.autofilledValues(true)
}

func (f fields) Scan(rows *sql.Rows) error {
	values := f.values()
	return rows.Scan(values...)
}
