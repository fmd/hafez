package main

import (
	"github.com/eknkc/amber"
	"github.com/fmd/gin"
	"net/http"
	"time"
	"os"
)

var amberDir string = "templates/"
var modTime time.Time

func TemplateMod() gin.HandlerFunc {
	return func (c *gin.Context) {
		if modTime.IsZero() {
			modTime = time.Now()
		}

		info, err := os.Stat(amberDir)
		if err != nil {
			modTime = time.Now()
		}

		if info.ModTime().After(modTime) {
			modTime = info.ModTime()
		}

		c.Writer.Header().Set("Last-Modified", modTime.Format(http.TimeFormat))
	}
}

func main() {
	r := gin.Default()
	r.Use(TemplateMod())

	templates, err := amber.CompileDir(amberDir, amber.DefaultDirOptions, amber.DefaultOptions)
	if err != nil {
		panic(err)
	}

	r.GET("/", func(c *gin.Context) {
		c.ExecHTML(200, templates["home"], nil)
	})

	r.GET("/menu", func(c *gin.Context) {
		c.ExecHTML(200, templates["menu"], nil)
	})

	r.Static("public/", "./public")

	r.Run(":5000")
}
