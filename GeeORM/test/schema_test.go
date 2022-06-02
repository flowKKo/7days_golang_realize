package test

import (
	"geeorm/dialect"
	"geeorm/schema"
	"testing"
)

type User struct {
	Name string `geeorm:"primary key"`
	Age  int
}

var TestDial, _ = dialect.GetDialect("sqlite3")

func TestParse(t *testing.T) {
	schema := schema.Parse(&User{}, TestDial)
	if schema.Name != "User" || len(schema.Fields) != 2{
		t.Fatal("failed to parse User struct")
	}
	if schema.GetField("Name").Tag != "primary key"{
		t.Fatal("failed to parse primary key")
	}
}
