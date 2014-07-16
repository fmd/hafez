package main

import (
    //"github.com/gin-gonic/gin/render"
    "github.com/gin-gonic/gin"
    "github.com/eknkc/amber"
    "html/template"
    "path/filepath"
    "net/http"
    "errors"
    "time"
    "fmt"
    "os"
)

type AmberGin struct {
    TemplateDir string
    Templates map[string]*template.Template
    modTime time.Time
}

func (a *AmberGin) Render(w http.ResponseWriter, code int, data ...interface{}) error {
    tmpl := a.Templates[data[0].(string)]
    if tmpl == nil {
        return errors.New(fmt.Sprintf("Template \"%s\" does not exist!", data[0].(string)))
    }

    if code >= 0 {
        w.Header().Set("Content-Type", "text/html")
        w.WriteHeader(code)
    }

    return tmpl.Execute(w, data[1])
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