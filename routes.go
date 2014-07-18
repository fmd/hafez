package main

import (
    "net/http"
    "github.com/gorilla/context"
    "github.com/unrolled/render"
    "github.com/julienschmidt/httprouter"
)

func Home(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    fmt.Println(w, context.Get(AppInfo))
}