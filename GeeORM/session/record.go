package session

import (
	"errors"
	"geeorm/clause"
	"reflect"
)

func (s *Session) Insert(values ...interface{}) (int64, error) {

	// firstly, call clause.Set() many times to build every sub sentence
	recordValues := make([]interface{}, 0)
	for _, value := range values {
		s.CallMethod(BeforeInsert, value)
		table := s.Model(value).RefTable()
		s.clause.Set(clause.INSERT, table.Name, table.FieldNames)
		// call RecordValues() to get values in the struct
		recordValues = append(recordValues, table.RecordValues(value))
	}

	// then call clause.Build(), build the final sql sentence at the sequence of passing
	s.clause.Set(clause.VALUES, recordValues...)
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)

	// after building, calling exec to execute the sql sentence
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	s.CallMethod(AfterInsert, nil)
	return result.RowsAffected()
}

func (s *Session) Find(values interface{}) error {
	// use reflect to transform the values into element type
	s.CallMethod(BeforeQuery, nil)
	destSlice := reflect.Indirect(reflect.ValueOf(values))
	destType := destSlice.Type().Elem()
	table := s.Model(reflect.New(destType).Elem().Interface()).RefTable()

	// then create select clause
	s.clause.Set(clause.SELECT, table.Name, table.FieldNames)
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)

	// execute sql query
	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		return err
	}

	// transform the query result into element type
	for rows.Next() {
		dest := reflect.New(destType).Elem()
		var values []interface{}

		for _, name := range table.FieldNames {
			//fmt.Println(reflect.TypeOf(dest.FieldByName(name).Addr()))
			//fmt.Println(reflect.TypeOf(dest.FieldByName(name).Addr().Interface()))

			// Addr() returns a pointer value representing the address
			// Addr().interface() is just the type of *int, *string
			// which means that it is a pointer of the element type
			values = append(values, dest.FieldByName(name).Addr().Interface())
		}

		// call rows.Scan() to assign each row's value to corresponding member in values
		if err := rows.Scan(values...); err != nil {
			return err
		}

		s.CallMethod(AfterQuery, dest.Addr().Interface())

		// add dest into slice destSlice
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return rows.Close()
}

// support map[string]interface{}
// also support kv list: "Name", "Tom", "Age", 18, ...

func (s *Session) Update(kv ...interface{}) (int64, error) {
	// judge if kv[0]'s type is map[string]interface{}
	// if the type isn't map, then transform it automatically
	s.CallMethod(BeforeUpdate, nil)
	m, ok := kv[0].(map[string]interface{})
	if !ok {
		m = make(map[string]interface{})
		for i := 0; i < len(kv); i += 2 {
			m[kv[i].(string)] = kv[i+1]
		}
	}
	s.clause.Set(clause.UPDATE, s.refTable.Name, m)
	sql, vars := s.clause.Build(clause.UPDATE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	s.CallMethod(AfterUpdate, nil)
	return result.RowsAffected()
}

func (s *Session) Delete() (int64, error){
	s.CallMethod(BeforeDelete, nil)
	s.clause.Set(clause.DELETE, s.refTable.Name)

	sql, vars := s.clause.Build(clause.DELETE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	s.CallMethod(AfterDelete, nil)
	return result.RowsAffected()
}

func (s *Session) Count() (int64, error){
	s.clause.Set(clause.COUNT, s.refTable.Name)

	sql, vars := s.clause.Build(clause.COUNT, clause.WHERE)
	row := s.Raw(sql, vars...).QueryRow()
	var tmp int64


	if err := row.Scan(&tmp); err != nil{
		return 0, err
	}
	return tmp, nil
}

// Limit adds limit condition to clause
func (s *Session) Limit(num int) *Session{
	s.clause.Set(clause.LIMIT, num)
	return s
}

// Where adds limit condition to clause
func (s *Session) Where (desc string, args ...interface{}) *Session{
	var vars []interface{}
	s.clause.Set(clause.WHERE, append(append(vars, desc), args...)...)
	return s
}

// OrderBy adds order by condition to clause
func (s *Session) OrderBy(desc string) *Session{
	s.clause.Set(clause.ORDERBY, desc)
	return s
}

// First only return one record
func (s*Session) First(value interface{}) error{
	dest := reflect.Indirect(reflect.ValueOf(value))
	destSlice := reflect.New(reflect.SliceOf(dest.Type())).Elem()
	if err := s.Limit(1).Find(destSlice.Addr().Interface()); err != nil{
		return err
	}
	if destSlice.Len() == 0 {
		return errors.New("NOT FOUND")
	}
	dest.Set(destSlice.Index(0))
	return nil
}