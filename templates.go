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
	dir 	:= t.directory
	dirOpts := amber.DefaultDirOptions
	opts    := amber.DefaultOptions
	t.templates, err = amber.CompileDir(dir, dirOpts, opts)
	return err
}

func (t *Templates) CompileFile(name string) (*template.Template, error) {
	path := path.Join(t.directory, name) + t.extension
	opts := amber.DefaultOptions
	return amber.CompileFile(path, opts)
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