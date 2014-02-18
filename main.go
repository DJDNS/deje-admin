package main

import (
	"github.com/campadrenalin/go-deje"
	"github.com/codegangsta/martini"
	"github.com/martini-contrib/render"
	"net/http"
)

func do_home(r render.Render) {
	r.HTML(200, "root", nil)
}

func do_open(req *http.Request, c *deje.DEJEController, r render.Render) {
	location, err := get_location(req)
	if err != nil {
		r.HTML(500, "error", err)
		return
	}

	doc := c.GetDocument(*location)
	r.HTML(200, "console", doc)
}

func do_notfound(r render.Render) {
	r.HTML(404, "404", nil)
}

func main() {
	controller := deje.NewDEJEController()
	m := martini.Classic()
	m.Map(controller)
	m.Use(martini.Static("static"))
	m.Use(render.Renderer(render.Options{
		Layout: "layout",
	}))

	m.Get("/", do_home)
	m.Get("/open", do_open)
	m.NotFound(do_notfound)

	m.Run()
}
