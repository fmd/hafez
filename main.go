package main

import (
	"os"
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func main() {
	router := httprouter.New()
	router.GET("/", func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		fmt.Fprintln(w, "Welcome to the homepage using http router")
	})

	n := negroni.Classic()
	n.UseHandler(router)

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "5000"
	}

	n.Run(":" + port)
}
