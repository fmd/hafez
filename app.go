package main

import (
	"github.com/eknkc/amber"
	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
	"os"
	"strconv"
)

//AppOptions Constants
const (
	DefaultTemplateDir string = "templates/"
	DefaultPublicDir   string = "public/"
	DefaultStaticUrl   string = "/assets"
	DefaultPort        int    = 5000
)

type AppOptions struct {
	TemplateDir string
	PublicDir   string
	StaticUrl   string
	Port        int
}

func (opt AppOptions) Process() AppOptions {
	if len(opt.PublicDir) == 0 {
		opt.PublicDir = DefaultPublicDir
	}

	if len(opt.StaticUrl) == 0 {
		opt.StaticUrl = DefaultStaticUrl
	}

	if len(opt.TemplateDir) == 0 {
		opt.TemplateDir = DefaultTemplateDir
	}

	if opt.Port == 0 {
		var err error

		opt.Port, err = strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			opt.Port = DefaultPort
		}
	}

	return opt
}

type App struct {

	//TODO: Change string for templateDir and publicDir to http.FileSystem? 
	templateDir string
	publicDir   string
	staticUrl   string
	port        int

	negroni   *negroni.Negroni
	router    *httprouter.Router
	frontend  *Frontend
	templates map[string]*template.Template
}

func (a *App) Templates() (map[string]*template.Template, error) {
	return amber.CompileDir(a.templateDir, amber.DefaultDirOptions, amber.DefaultOptions)
}

func (a *App) Frontend() *Frontend {
	return NewFrontend(a.router, a.templates)
}

func NewApp(opts AppOptions) *App {
	opts = opts.Process()

	a := &App{
		templateDir: opts.TemplateDir,
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

	//Set up templates
	var err error
	a.templates, err = a.Templates()
	if err != nil {
		panic(err)
	}

	//Set up router
	r := httprouter.New()
	a.negroni.UseHandler(r)

	//Set up frontend routing
	a.frontend = a.Frontend()

	return a
}

func (a *App) Run() {
	a.negroni.Run(":" + strconv.Itoa(a.port))
}