package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"mime"
	"path/filepath"
	"path"
	"fmt"
	"os"
)

func StaticWithGzipAndHead(r *gin.Engine, p string, root string) {
	p = path.Join(p, "/*filepath")
	fileServer := http.FileServer(http.Dir(root))

	r.GET(p, func(c *gin.Context) {
		original := c.Request.URL.Path
		newPath := c.Params.ByName("filepath")
		gzPath := fmt.Sprintf("%s.gz",newPath)

		_, err := os.Stat(path.Join(root, gzPath))

		if err != nil {
			c.Request.URL.Path = newPath
			fileServer.ServeHTTP(c.Writer, c.Request)
			c.Request.URL.Path = original
			return
		}

		c.Request.URL.Path = gzPath
		ctype := mime.TypeByExtension(filepath.Ext(original))
		c.Writer.Header().Set("Content-Type", ctype)
		c.Writer.Header().Set("Content-Encoding", "gzip")

		fileServer.ServeHTTP(c.Writer, c.Request)
		c.Request.URL.Path = original
	})

	r.HEAD(p, func (c *gin.Context) {
		fp := path.Join(root, c.Params.ByName("filepath"))

		info, err := os.Stat(fp)
		if err != nil || info == nil {
			c.Abort(404)
			return
		}

		c.Writer.Header().Set("Content-Type",  mime.TypeByExtension(filepath.Ext(fp)))
		c.Writer.Header().Set("Last-Modified", info.ModTime().Format(http.TimeFormat))
		c.Abort(200)
		return
	})
}

func main() {
	r := gin.Default()

	publicDir := "public"
	templateDir := "templates"
	staticUrl := "/assets"

	//gzip := NewGzipGin(publicDir, "zipped", staticUrl)
	amber, err := NewAmberGin(templateDir)

	if err != nil {
		panic(err)
	}

	r.Use(amber.DevMiddleware())
	//r.Use(gzip.Middleware())

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

	StaticWithGzipAndHead(r, staticUrl, publicDir)

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "5000"
	}

	r.Run(fmt.Sprintf(":%s", port))
}
