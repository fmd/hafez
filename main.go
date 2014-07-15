package main

import (
	"github.com/eknkc/amber"
	"github.com/fmd/gin"
	"html/template"
	"path/filepath"
	"net/http"
	"path"
	"time"
	"fmt"
	"os"
)

type AmberGin struct {
	TemplateDir string
	Templates map[string]*template.Template
	
	modTime time.Time
}

func (m *AmberGin) WalkTemplates(path string, info os.FileInfo, err error) error {
	if info.ModTime().After(m.modTime) {
		m.modTime = info.ModTime()
		return nil
	}
	return nil
}

func (m *AmberGin) TemplatesChanged() bool {
	if m.modTime.IsZero() {
		m.modTime = time.Now()
		return true
	}

	oldMod := m.modTime
	err := filepath.Walk(m.TemplateDir, m.WalkTemplates)

	if !oldMod.Equal(m.modTime) {
		return true
	}

	if err != nil {
		m.modTime = time.Now()
		return true
	}

	return false
}

func (m *AmberGin) Compile() error {
	t, err := amber.CompileDir(m.TemplateDir, amber.DefaultDirOptions, amber.DefaultOptions)
	if err != nil {
		return err
	}

	m.Templates = t
	return nil
}

func (m *AmberGin) DevMiddleware() gin.HandlerFunc {
	return func (c *gin.Context) {
		if (m.TemplatesChanged()) {
			err := m.Compile()	
			if err != nil {
				c.String(200, err.Error())
			}

			c.Writer.Header().Set("Last-Modified", m.modTime.Format(http.TimeFormat))
		}
	}
}

func NewAmberGin(templateDir string) (*AmberGin, error) {
	m := &AmberGin{}
	m.TemplateDir = templateDir
	err := m.Compile()
	if err != nil {
		return nil, err
	}

	return m, nil
}

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

	m, err := NewAmberGin(templateDir)
	if err != nil {
		panic(err)
	}

	r.Use(m.DevMiddleware())

	r.GET("/", func(c *gin.Context) {

		data := struct {
			StaticUrl string
		} {
			staticUrl,
		}

		c.ExecHTML(200, m.Templates["home"], data)
	})

	r.GET("/menu", func(c *gin.Context) {

		data := struct {
			StaticUrl string
		} {
			staticUrl,
		}

		c.ExecHTML(200, m.Templates["menu"], data)
	})

	r.Static(staticUrl, publicDir)
	r.HEAD(StaticHead(staticUrl, publicDir))

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "5000"
	}

	r.Run(fmt.Sprintf(":%s", port))
}