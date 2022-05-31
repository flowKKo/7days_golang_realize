package gee

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

func trace(message string) string{
	var pcs [32]uintptr

	// Callers() is used to return the program counter of the calling stack
	// the 0th caller is callers itself, the 1st is trace, the 2nd is defer func
	// so we skip the front 3 caller to simplify the log
	n := runtime.Callers(3, pcs[:])
	var str strings.Builder
	str.WriteString(message + "\nTraceback:")

	for _, pc := range pcs[:n]{
		// get pc's corresponding function
		fn := runtime.FuncForPC(pc)
		//get the filename nad line calling the function
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}


func Recovery() HandlerFunc{
	return func(c *Context){
		defer func(){
			if err := recover(); err != nil{
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				c.Status(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}


