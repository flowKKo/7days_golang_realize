package session

import (
	"geeorm/log"
	"reflect"
)

const (
	BeforeQuery  = "BeforeQuery"
	AfterQuery   = "AfterQuery"
	BeforeUpdate = "BeforeUpdate"
	AfterUpdate  = "AfterUpdate"
	BeforeDelete = "BeforeDelete"
	AfterDelete  = "AfterDelete"
	BeforeInsert = "BeforeInsert"
	AfterInsert  = "AfterInsert"
)

//CallMethod calls the registered hooks
func (s *Session) CallMethod(method string, value interface{}){

	// use MethodByName to get function pointer
	// if we want to use hooks, then we should realize target struct's hooks method
	// such as func (account *Account) BeforeInsert(s *session.Session) error {}
	fm := reflect.ValueOf(s.RefTable().Model).MethodByName(method)

	if value != nil{
		fm = reflect.ValueOf(value).MethodByName(method)
	}

	param := []reflect.Value{reflect.ValueOf(s)}
	if fm.IsValid(){
		if v := fm.Call(param); len(v) > 0{
			if err, ok := v[0].Interface().(error); ok{
				log.Error(err)
			}
		}
	}
	return
}