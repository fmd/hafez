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

	//Start off like r.Static
	p = path.Join(p, "/*filepath")
	fileServer := http.FileServer(http.Dir(root))

	//Start GET request (as usual)
	r.GET(p, func(c *gin.Context) {

		original := c.Request.URL.Path
		newPath := c.Params.ByName("filepath")

		//This path will be called by os.Stat to see if there's a gzipped version of the file.
		gzPath := fmt.Sprintf("%s.gz",newPath)

		//See if the gzipped version exists.
		_, err := os.Stat(path.Join(root, gzPath))

		//If the gzipped version doesn't exist, serve the file as usual.
		if err != nil {
			c.Request.URL.Path = newPath
			fileServer.ServeHTTP(c.Writer, c.Request)
		} else {
			//If the gzipped file does exist, preserve the original content type, and set encoding to gzip.
			c.Request.URL.Path = gzPath
			ctype := mime.TypeByExtension(filepath.Ext(original))
			c.Writer.Header().Set("Content-Type", ctype)
			c.Writer.Header().Set("Content-Encoding", "gzip")

			//Serve the file.
			fileServer.ServeHTTP(c.Writer, c.Request)
		}

		//Revert the path as usual!
		c.Request.URL.Path = original
	})

	//I use live.js, which makes HEAD requests to the static files (as do a lot of tools and other things).
	//This method responds to HEAD requests with a correct content type and last-modified date.
	r.HEAD(p, func (c *gin.Context) {
		fp := path.Join(root, c.Params.ByName("filepath"))

		info, err := os.Stat(fp)
		if err != nil || info == nil {
			c.Abort(404)
		} else {
			c.Writer.Header().Set("Content-Type",  mime.TypeByExtension(filepath.Ext(fp)))
			c.Writer.Header().Set("Last-Modified", info.ModTime().Format(http.TimeFormat))
			c.Abort(200)
		}
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
