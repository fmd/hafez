package main

import (
    "github.com/gin-gonic/gin"
    "net/url"
    "path"
    "log"
    "fmt"
    "os"
)

type GzipGin struct {
    ZippedDir string
    UnzippedDir string
    StaticUrl string
}

func NewGzipGin(zippedDir string, unzippedDir string, staticUrl string) *GzipGin {
    return &GzipGin{zippedDir, unzippedDir, staticUrl}
}

func (g *GzipGin) Middleware() gin.HandlerFunc {
    return func (c *gin.Context) {
            file := c.Params.ByName("filepath")

            if len(file) == 0 {
                return
            }

            filegz := fmt.Sprintf("%s.gz", file)
            fp := path.Join(g.ZippedDir, filegz)
            up := fmt.Sprintf("%s%s", g.StaticUrl, filegz)

            url, err := url.Parse(up)
            
            if err != nil {
                log.Println(err.Error())
                return
            }

            info, err := os.Stat(fp)
            
            if err == nil && info != nil {
                c.Writer.Header().Set("Content-Encoding", "gzip")
                c.Request.URL = url
            }

            return
    }
}