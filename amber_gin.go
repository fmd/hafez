package main

import (
    "github.com/eknkc/amber"
    "github.com/fmd/gin"
    "html/template"
    "path/filepath"
    "net/http"
    "time"
    "os"
)

type AmberGin struct {
    TemplateDir string
    Templates map[string]*template.Template
    modTime time.Time
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