package main

import (
	//"github.com/eknkc/amber"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// This handler will match /user/john but will not match neither /user/ or /user
	r.GET("/user/:name", func(c *gin.Context) {
		name := c.Params.ByName("name")
		message := "Hello " + name
		c.String(200, message)
	})

	// However, this one will match /user/john and also /user/john/send
	r.GET("/user/:name/*action", func(c *gin.Context) {
		name := c.Params.ByName("name")
		action := c.Params.ByName("action")
		message := name + " is " + action
		c.String(200, message)
	})

	// Listen and server on 0.0.0.0:8080
	r.Run(":8080")
}
