package main

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"os"
	"strconv"
)

//AppOptions Constants
const (
	DefaultPublicDir string = "public/"
	DefaultStaticUrl string = "/assets"
	DefaultPort      int    = 5000
)

type AppOptions struct {
	PublicDir string
	StaticUrl string
	Port      int
}

func (opt AppOptions) Process() AppOptions {
	if len(opt.PublicDir) == 0 {
		opt.PublicDir = DefaultPublicDir
	}

	if len(opt.StaticUrl) == 0 {
		opt.StaticUrl = DefaultStaticUrl
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
	publicDir string
	staticUrl string
	port      int

	negroni *negroni.Negroni
}

func NewApp(opts AppOptions) *App {
	opts = opts.Process()

	a := &App{
		publicDir: opts.PublicDir,
		staticUrl: opts.StaticUrl,
		port:      opts.Port,
	}

	//Set up Negroni
	a.negroni = negroni.New()
	a.negroni.Use(negroni.NewRecovery())
	a.negroni.Use(negroni.NewLogger())

	//Set up static fileserver
	s := negroni.NewStatic(http.Dir(a.publicDir))
	s.Prefix = a.staticUrl
	a.negroni.Use(s)

	return a
}

func (a *App) Run() {
	a.negroni.Run(":" + strconv.Itoa(a.port))
}

func (a *App) Routes() {
	r := httprouter.New()

	r.GET("/", a.frontend.Home)
	r.GET("/menu", a.frontend.Menu)
	r.GET("/gallery", a.frontend.Gallery)
	r.GET("/book", a.frontend.Book)

	a.negroni.UseHandler(r)
}
