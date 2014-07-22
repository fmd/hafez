package main

import (
	"path"
	"errors"
	"html/template"
	"github.com/eknkc/amber"
)

type Templates struct {
	directory string
	extension string
	recompile bool                          //Recompile on refresh? (disable in production).
	templates map[string]*template.Template
}

func NewTemplates(directory string, extension string, recompile bool) (*Templates, error) {
	t := &Templates{
		directory: directory,
		extension: extension,
		recompile: recompile,
	}

	err := t.CompileDir()
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (t *Templates) CompileDir() error {
	var err error
	t.templates, err = amber.CompileDir(t.directory, amber.DefaultDirOptions, amber.DefaultOptions)
	return err
}

func (t *Templates) CompileFile(name string) (*template.Template, error) {
	return amber.CompileFile(path.Join(t.directory, name) + t.extension, amber.DefaultOptions)
}

func (t *Templates) GetTemplate(name string) (*template.Template, error) {
	if t.recompile {
		return t.CompileFile(name)
	}

	_, ok := t.templates[name]
	if ok {
		return t.templates[name], nil
	}
	
	return nil, errors.New("Could not find template: " + name)
}