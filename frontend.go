package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	//"github.com/gorilla/context"
)

type Frontend struct {

}

func NewFrontend(r *httprouter.Router) *Frontend {
	f := &Frontend{}

	r.GET("/", f.Home)
	r.GET("/menu", f.Menu)
	r.GET("/contact", f.Contact)
	r.GET("/book", f.Book)

	return f
}

func (f *Frontend) Home(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {

}

func (f *Frontend) Menu(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
	
}

func (f *Frontend) Contact(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
	
}

func (f *Frontend) Book(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
	
}