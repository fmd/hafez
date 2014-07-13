package main

import (
	"bytes"
	"github.com/eknkc/amber"
	"github.com/fmd/gin"
)

func main() {
	r := gin.Default()

	home, err := amber.CompileFile("templates/home.tmpl", amber.Options{})
	if err != nil {
		panic(err)
	}

	r.GET("/", func(c *gin.Context) {
		var out bytes.Buffer

		err := home.Execute(&out, nil)
		if err != nil {
			c.String(500, "Internal Server Error.")
		}

		c.ExecHTML(200, home, nil)
	})

	r.Run(":5000")
}
