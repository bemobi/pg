package pg

import (
	"reflect"
	"strings"
)

const tagName = "db"

type tag struct {
	fieldname  string
	omitfield  bool
	autofilled bool
	pk         bool
}

func newTag(field reflect.StructField) *tag {
	values := strings.Split(field.Tag.Get(tagName), ",")
	tag := tag{omitfield: true}
	for _, value := range values {
		if value == "autofilled" {
			tag.autofilled = true
		} else if value == "pk" {
			tag.pk = true
		} else if value != "-" {
			tag.fieldname = value
			tag.omitfield = false
		}
	}
	return &tag
}
