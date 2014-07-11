package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
    "github.com/eknkc/amber"
)

func main() {
	m := martini.Classic()
	m.Use(render.Renderer())

	m.Get("/", func() (int, string) {
		return 200, 
	})

	m.Run()
}
