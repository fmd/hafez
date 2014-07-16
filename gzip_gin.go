package main

import (
    "github.com/gin-gonic/gin"
    //"io/ioutil"
    //"net/http"
    "net/url"
    "strings"
    "path"
    "fmt"
    "os"
)

type GzipGin struct {
    PublicDir string
    ZippedDir string
    StaticUrl string
}

func NewGzipGin(publicDir string, zippedDir string, staticUrl string) *GzipGin {

    if strings.TrimSuffix(staticUrl, "/") == staticUrl {
        staticUrl = fmt.Sprintf("%s/", staticUrl)
    }

    return &GzipGin{publicDir, zippedDir, staticUrl}
}

func (g *GzipGin) Middleware() gin.HandlerFunc {
    return func (c *gin.Context) {
            file := c.Params.ByName("filepath")

            if len(file) == 0 {
                return
            }

            filegz := fmt.Sprintf("%s.gz", file)
            fp := path.Join(g.PublicDir, g.ZippedDir, filegz)
            up := fmt.Sprintf("%s%s%s", g.StaticUrl, g.ZippedDir, filegz)
            url, err := url.Parse(up)

            if err != nil {
                return
            }

            info, err := os.Stat(fp)

            if err == nil && info != nil {
                //bytes, err := ioutil.ReadFile(path.Join(g.PublicDir, file))
                if err != nil {
                    return
                }

                for i := range c.Params {
                    if c.Params[i].Key == "filepath" {
                        c.Params[i].Value = path.Join(g.ZippedDir, filegz)
                        break
                    }
                }

                c.Request.URL.Path = url.Path

                c.Next()
                c.Writer.Header().Set("Content-Encoding", "gzip")
                c.Writer.Header().Set("Content-Type", "poop")
            }

            return
    }
}