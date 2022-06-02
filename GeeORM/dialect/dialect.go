package dialect

import "reflect"

// use dialect to isolate the difference of database
// and make it easy to expand

var dialectsMap = map[string]Dialect{}

type Dialect interface{
	// DataTypeOf is used to transform go datatype to db datatype
	DataTypeOf(typ reflect.Value) string
	// TableExistSQL is used to check if a table exists
	TableExistSQL(tableName string) (string, []interface{})
}

func RegisterDialect(name string, dialect Dialect){
	dialectsMap[name] = dialect
}

func GetDialect(name string) (dialect Dialect, ok bool){
	dialect, ok = dialectsMap[name]
	return
}



