package main

import (
	"github.com/gorilla/context"
	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

type key int

const AppContextKey key = 0

type AppContext struct {
	StaticUrl   string
	Development bool
}

func (a *App) ContextMiddleware(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	c := AppContext{
		StaticUrl:   a.staticUrl,
		Development: a.development,
	}

	context.Set(req, AppContextKey, c)
	next(rw, req)
}

type AppOptions struct {
	Development bool
	TemplateDir string
	PublicDir   string
	StaticUrl   string
	Port        int
}

type App struct {
	development bool

	//TODO: Change string for templateDir and publicDir to http.FileSystem?
	publicDir   string
	staticUrl   string
	port        int

	negroni   *negroni.Negroni
	router    *httprouter.Router
	templates *Templates
	frontend  *Frontend
}

func NewApp(opts AppOptions) *App {
	var err error

	//Create the app instance
	a := &App{
		development: opts.Development,
		publicDir:   opts.PublicDir,
		staticUrl:   opts.StaticUrl,
		port:        opts.Port,
	}

	//Set up Negroni
	a.negroni = negroni.New()
	a.negroni.Use(negroni.NewRecovery())
	a.negroni.Use(negroni.NewLogger())

	//Set up static fileserver
	s := negroni.NewStatic(http.Dir(a.publicDir))
	s.Prefix = a.staticUrl
	a.negroni.Use(s)

	//Use app context middleware
	a.negroni.Use(negroni.HandlerFunc(a.ContextMiddleware))

	//Set up router
	a.router = httprouter.New()
	a.negroni.UseHandler(a.router)

	//Set up templates. TODO: Base recompile bool on development flag.
	a.templates, err = NewTemplates(opts.TemplateDir, ".amber", a.development)
	if err != nil {
		panic(err)
	}

	//Set up frontend routing
	a.frontend = NewFrontend(a.router, a.templates, a.development)
	return a
}

func (a *App) Run() {
	a.negroni.Run(":" + strconv.Itoa(a.port))
}