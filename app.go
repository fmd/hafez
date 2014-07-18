package main

import (
    "os"
    "strconv"
    "net/http"
    "github.com/gorilla/context"
    "github.com/fmd/hafez/routes"
    "github.com/codegangsta/negroni"
    "github.com/julienschmidt/httprouter"
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

//Context keys
const (
    AppInfo = iota
    UserInfo
)

type App struct {
    publicDir string
    staticUrl string
    port      int

    router  *httprouter.Router
    negroni *negroni.Negroni
}

func (a *App) Middleware(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
    context.Set(req, AppInfo, "Info!")
}

func NewApp(opts AppOptions) *App {
    opts = opts.Process()

    router := httprouter.New()

    n := negroni.New()
    n.Use(negroni.NewRecovery())
    n.Use(negroni.NewLogger())

    a := &App{
        publicDir: opts.PublicDir,
        staticUrl: opts.StaticUrl,
        port:      opts.Port,

        router:    router,
        negroni:   n,
    }

    a.UseHandler(a.Middleware)
    a.Routes()

    return a
}

func (a *App) Run() {
    a.negroni.Run(":" + strconv.Itoa(a.port))
}

func (a *App) Routes() {
    a.router.GET("/", routes.Home)
    a.negroni.UseHandler(a.router)
}