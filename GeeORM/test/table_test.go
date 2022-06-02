package test

import (
	"geeorm"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

func TestSession_CreateTable(t *testing.T) {
	engine, _ := geeorm.NewEngine("sqlite3", "gee.db")
	defer engine.Close()

	s := engine.NewSession().Model(&User{})
	_ = s.DropTable()
	_ = s.CreateTable()
	if !s.HasTable() {
		t.Fatal("Failed to create table User")
	}
}