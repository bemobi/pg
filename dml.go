package pg

import (
	"bytes"
	"database/sql"
	"fmt"
	"strings"
)

type driverWrapper interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

func query(driver driverWrapper, result Result, sql string, params ...interface{}) ([]Result, error) {
	rows, err := driver.Query(sql, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]Result, 0)
	for rows.Next() {
		instance := NewResult(result)
		err := fields{instance}.Scan(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, instance)
	}
	return list, nil
}

func exec(driver driverWrapper, sql string, params ...interface{}) (sql.Result, error) {
	return driver.Exec(sql, params...)
}

func create(driver driverWrapper, entity Entity) error {
	ef := fields{entity}

	updatableFields := ef.updatable()
	updatableFieldsStr := strings.Join(updatableFields, ",")

	notUpdatableFields := ef.notUpdatable()
	notUpdatableFieldsStr := strings.Join(notUpdatableFields, ",")

	placeholders := make([]string, 0)
	for index := range updatableFields {
		placeholders = append(placeholders, fmt.Sprintf("$%d", index+1))
	}
	placeholderList := strings.Join(placeholders, ",")

	query := fmt.Sprintf(
		"insert into %s (%s) values (%s) returning %s",
		entity.Table(), updatableFieldsStr, placeholderList, notUpdatableFieldsStr,
	)

	updatableValues := ef.updatableValues()
	notUpdatableValues := ef.notUpdatableValues()

	err := driver.QueryRow(query, updatableValues...).Scan(notUpdatableValues...)
	switch {
	case err == sql.ErrNoRows:
		return fmt.Errorf("Could not insert row: %s", err)
	case err != nil:
		return fmt.Errorf("Failed to execute statement: %s", err)
	default:
		return nil
	}
}

func delete(driver driverWrapper, entity Entity) error {
	ef := fields{entity}

	pkFields := ef.pk(true)
	pkValues := ef.pkValues(true)

	pkList := make([]string, 0)
	for index, field := range pkFields {
		pkList = append(
			pkList,
			fmt.Sprintf("%s = $%d", field, index+1),
		)
	}
	pkListStr := strings.Join(pkList, " and ")

	var query bytes.Buffer
	query.WriteString("delete from ")
	query.WriteString(entity.Table())
	query.WriteString(" where ")
	query.WriteString(pkListStr)

	_, err := driver.Exec(query.String(), pkValues...)
	if err != nil {
		return fmt.Errorf("Failed to execute statement: %s", err)
	}

	return nil
}

func update(driver driverWrapper, entity Entity) error {
	ef := fields{entity}

	updatableFields := ef.updatable()
	updatableValues := ef.updatableValues()

	notUpdatableFields := ef.notUpdatable()
	notUpdatableValues := ef.notUpdatableValues()
	notUpdatableFieldsStr := strings.Join(notUpdatableFields, ",")

	pkFields := ef.pk(true)
	pkValues := ef.pkValues(true)

	index := 0
	fieldValueList := make([]string, 0)
	for _, field := range updatableFields {
		fieldValueList = append(
			fieldValueList,
			fmt.Sprintf("%s = $%d", field, index+1),
		)
		index += 1
	}
	fieldValueListStr := strings.Join(fieldValueList, ",")

	pkList := make([]string, 0)
	for _, field := range pkFields {
		pkList = append(
			pkList,
			fmt.Sprintf("%s = $%d", field, index+1),
		)
		index += 1
	}
	pkListStr := strings.Join(pkList, " and ")

	var query bytes.Buffer
	query.WriteString("update ")
	query.WriteString(entity.Table())
	query.WriteString(" set ")
	query.WriteString(fieldValueListStr)
	query.WriteString(" where ")
	query.WriteString(pkListStr)
	query.WriteString(" returning ")
	query.WriteString(notUpdatableFieldsStr)

	params := make([]interface{}, 0)
	for _, param := range updatableValues {
		params = append(params, param)
	}
	for _, param := range pkValues {
		params = append(params, param)
	}

	err := driver.QueryRow(query.String(), params...).Scan(notUpdatableValues...)
	switch {
	case err == sql.ErrNoRows:
		return fmt.Errorf("Could not insert row: %s", err)
	case err != nil:
		return fmt.Errorf("Failed to execute statement: %s", err)
	default:
		return nil
	}
}

func findOne(driver driverWrapper, entity Entity, where string, whereParams ...interface{}) (Entity, error) {
	entities, err := findAll(driver, entity, where, whereParams...)
	if err != nil {
		return nil, err
	}
	switch len(entities) {
	case 1:
		return entities[0], nil
	case 0:
		return nil, ERecordNotFound
	default:
		return nil, EMultipleResults
	}
}

func findAll(driver driverWrapper, entity Entity, where string, whereParams ...interface{}) ([]Entity, error) {
	fieldsStr := strings.Join(fields{entity}.all(), ",")

	var sql bytes.Buffer
	sql.WriteString(fmt.Sprintf("select %s from %s", fieldsStr, entity.Table()))

	if where != "" {
		sql.WriteString(" where ")
		sql.WriteString(where)
	}

	sql.WriteString(" order by ")
	sql.WriteString(entity.OrderBy())

	rows, err := driver.Query(sql.String(), whereParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]Entity, 0)
	for rows.Next() {
		instance := NewEntity(entity)
		err := fields{instance}.Scan(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, instance)
	}
	return list, nil
}
