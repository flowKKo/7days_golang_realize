package test

import (
	"geeorm/clause"
	"reflect"
	"testing"
)

func TestSelect(t *testing.T) {
	var cls clause.Clause
	cls.Set(clause.LIMIT, 3)
	cls.Set(clause.SELECT, "User", []string{"*"})
	cls.Set(clause.WHERE, "Name = ?", "Tom")
	cls.Set(clause.ORDERBY, "Age ASC")
	sql, vars := cls.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	t.Log(sql, vars)
	if sql != "SELECT * FROM User WHERE Name = ? ORDER BY Age ASC LIMIT ?" {
		t.Fatal("failed to build SQL")
	}
	if !reflect.DeepEqual(vars, []interface{}{"Tom", 3}) {
		t.Fatal("failed to build SQLVars")
	}
}

func TestClause_Build(t *testing.T) {
	t.Run("select", func(t *testing.T) {
		TestSelect(t)
	})
}