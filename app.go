package main

import (
    "os"
    "fmt"
    "strconv"
    //"github.com/unrolled/render"
    "github.com/codegangsta/negroni"
    "github.com/julienschmidt/httprouter"
    "net/http"
)

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

    router  *httprouter.Router
    negroni *negroni.Negroni
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

    return a
}

func (a *App) Run() {
    a.negroni.Run(":" + strconv.Itoa(a.port))
}

func (a *App) Routes() {
    a.router.GET("/", func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
        fmt.Println(w, "Hello world!")
    })
}

