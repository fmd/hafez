package main

import (
	//"log"
	"net/http"
	"html/template"
	"github.com/unrolled/render"
)

type Renderer struct {
	render.Head
	Template *template.Template
}

func (r Renderer) Render(w http.ResponseWriter, v interface{}) error {
	r.Head.Write(w)
	return r.Template.Execute(w, v)
}

func NewRenderer(c string, s int, t *template.Template) Renderer {
	return Renderer{
		Head: render.Head{
			ContentType: c,
			Status:      s,
		},
		Template: t,
	}
}