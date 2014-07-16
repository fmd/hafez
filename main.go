package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
	"fmt"
	"os"
)

func StaticHead(staticUrl string, publicDir string) (string, func(*gin.Context)) {
	p := path.Join(staticUrl, "/*filepath")
	return p, func (c *gin.Context) {
		fp := c.Params.ByName("filepath")

		info, err := os.Stat(path.Join(publicDir, fp))
		if err != nil || info == nil {
			c.Abort(404)
			return
		}

		c.Writer.Header().Set("Last-Modified", info.ModTime().Format(http.TimeFormat))
		c.Abort(200)
		return
	}
}

func main() {
	r := gin.Default()

	publicDir := "public"
	templateDir := "templates"
	staticUrl := "/assets"

	gzip := NewGzipGin(path.Join(publicDir,"zipped"), publicDir, staticUrl)
	amber, err := NewAmberGin(templateDir)

	if err != nil {
		panic(err)
	}

	r.Use(amber.DevMiddleware())
	r.Use(gzip.Middleware())

	r.GET("/", func(c *gin.Context) {

		data := struct {
			StaticUrl string
		} {
			staticUrl,
		}

		c.Render(200, amber, "home", data)
	})

	r.GET("/menu", func(c *gin.Context) {

		data := struct {
			StaticUrl string
		} {
			staticUrl,
		}

		c.Render(200, amber, "menu", data)
	})

	r.Static(staticUrl, publicDir)
	r.HEAD(StaticHead(staticUrl, publicDir))

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "5000"
	}

	r.Run(fmt.Sprintf(":%s", port))
}
