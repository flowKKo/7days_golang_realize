package clause

import (
	"strings"
)

type Clause struct {
	sql     map[Type]string
	sqlVars map[Type][]interface{}
}

type Type int

const (
	INSERT Type = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBY
	UPDATE
	DELETE
	COUNT
)

// Set is used to call generator to generate sql sentence's
// sub sentence and store it in Clause body
func (c *Clause) Set(name Type, vars ...interface{}){
	if c.sql == nil{
		c.sql = make(map[Type]string)
		c.sqlVars = make(map[Type][]interface{})
	}
	sql, vars := generators[name](vars...)
	c.sql[name] = sql
	c.sqlVars[name] = vars
}

// Build is used to combine several sub sentences together into a complete one
func (c *Clause) Build(orders ...Type)(string, []interface{}){
	var sqls []string
	var vars []interface{}
	for _, order := range orders{
		if sql, ok := c.sql[order]; ok{
			sqls = append(sqls, sql)
			vars = append(vars, c.sqlVars[order]...)
		}
	}
	return strings.Join(sqls, " "), vars
}


