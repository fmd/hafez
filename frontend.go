package main

import (
	"net/http"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
)

//TODO: Move render.Render into templates.go for easy rendering across multiple 'apps'?
type Frontend struct {
	render    *render.Render
	templates *Templates
	showError bool
}

func NewFrontend(r *httprouter.Router, t *Templates, showError bool) *Frontend {
	f := &Frontend{
		showError: showError,
	}

	//Set up the route(s).
	r.GET("/",          f.Page("home",    http.StatusOK))

	//Set up router 404 and 500 pages
	r.NotFound     = f.NotFound()
	r.PanicHandler = f.ErrorPage()

	//Create our render.Render.
	n := render.New(render.Options{})
	f.render = n
	f.templates = t

	return f
}

func (f *Frontend) NotFound() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		//Attempt to get the template. If we can't find it, 500 out.
		t, err := f.templates.GetTemplate("404")
		if err != nil {
			panic(err)
		}

		//Template is OK. Render the template.
		r := NewRenderer(render.ContentHTML, 404, t)
		r.Render(rw, nil)
	}
}

func (f *Frontend) ErrorPage() func(http.ResponseWriter, *http.Request, interface{}) {
	return func(rw http.ResponseWriter, req *http.Request, _ interface{}) {
		//Attempt to get the template. If we can't find it, 500 out.
		t, err := f.templates.GetTemplate("500")
		if err != nil {
			panic(err)
		}

		//Template is OK. Render the template.
		r := NewRenderer(render.ContentHTML, 500, t)
		r.Render(rw, nil)
	}
}

func (f *Frontend) Page(name string, code int) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {

		//Attempt to get the template. If we can't find it, 500 out.
		t, err := f.templates.GetTemplate(name)
		if err != nil {
			if !f.showError {
				panic(err)
			}
			
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		//Get app context
		context := context.Get(req, AppContextKey)

		//Template is OK. Render the template.
		r := NewRenderer(render.ContentHTML, code, t)
		r.Render(rw, context)
	}
}