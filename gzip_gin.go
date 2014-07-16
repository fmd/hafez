package main

import (
    "github.com/fmd/gin"
    "net/url"
    "path"
    "log"
    "fmt"
    "os"
)

type GzipGin struct {
    ZippedFolder string
    UnzippedFolder string
}

func NewGzipGin(zippedFolder string, unzippedFolder string) *GzipGin {
    return &GzipGin{zippedFolder, unzippedFolder}
}

func (g *GzipGin) Middleware() gin.HandlerFunc {
    return func (c *gin.Context) {
            fp := c.Params.ByName("filepath")

            if len(fp) == 0 {
                return
            }

            np := fmt.Sprintf("%s.gz",path.Join(g.ZippedFolder, fp))
            url, err := url.Parse(np)
            if err != nil {
                log.Println(err.Error())
                return
            }

            info, err := os.Stat(np)
            
            if err == nil && info != nil {
                c.Writer.Header().Set("Content-Encoding", "gzip")
                c.Req.URL = url
                log.Println("GZIP!")
            }

            return
    }
}