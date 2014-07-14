package main

import (
	"github.com/eknkc/amber"
	"github.com/fmd/gin"
)

func main() {
	r := gin.Default()

	templates, err := amber.CompileDir("templates/", amber.DefaultDirOptions, amber.DefaultOptions)
	if err != nil {
		panic(err)
	}

	r.GET("/", func(c *gin.Context) {
		c.ExecHTML(200, templates["home"], nil)
	})

	r.GET("/menu", func(c *gin.Context) {
		c.ExecHTML(200, templates["menu"], nil)

	})

	r.Run(":5000")
}
