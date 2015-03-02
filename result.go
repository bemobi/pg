package pg

import "reflect"

type Result interface{}

type ResultHandler interface {
	Query(result Result, sql string, params ...interface{}) ([]Result, error)
}

func NewResult(r Result) Result {
	mustBeStructPtr(r)
	resultType := reflect.TypeOf(r).Elem()
	newResult := reflect.New(resultType)
	return newResult.Interface().(Result)
}
