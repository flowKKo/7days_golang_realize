package schema

import (
	"geeorm/dialect"
	"go/ast"
	"reflect"
)

type Field struct {
	Name string
	Type string
	Tag  string
}

type Schema struct {
	// the mapped object
	Model interface{}
	// table name
	Name   string
	Fields []*Field
	// contain fields' name: the row name
	FieldNames []string
	fieldMap   map[string]*Field
}

func (schema *Schema) GetField(name string) *Field {
	return schema.fieldMap[name]
}

// Parse is used to parse any struct object to Schema object
func Parse(dest interface{}, d dialect.Dialect) *Schema {
	// if dest is a pointer, through reflect.Indirect we can get pointer's
	// element type's value, and use .Type() method to reflect.Type obj
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model:    dest,
		Name:     modelType.Name(),
		fieldMap: make(map[string]*Field),
	}
	// use modelType.NumField() to get struct's member num
	for i := 0; i < modelType.NumField(); i++ {
		// use Type.Filed(i) to access struct's member variants
		p := modelType.Field(i)
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				// get member variants' name and Type in target database
				Name: p.Name,
				// use reflect.Indirect() to get the object of the given pointer
				// use DataTypeOf to transform the datatype into a database one
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			// use Lookup to get tag's content
			if v, ok := p.Tag.Lookup("geeorm"); ok {
				field.Tag = v
			}
			// add current Field to schema
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = field
		}
	}
	return schema
}

func (schema *Schema) RecordValues(dest interface{})[]interface{}{
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	var fieldValues []interface{}
	for _, field := range schema.Fields{
		fieldValues = append(fieldValues, destValue.FieldByName(field.Name).Interface())
	}
	return fieldValues
}