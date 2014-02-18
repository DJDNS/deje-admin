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
	//length := strconv.FormatUint(uint64(doc.Events.Length()), 10)
	r.HTML(200, "console", doc)
}

func main() {
	controller := deje.NewDEJEController()
	m := martini.Classic()
	m.Map(controller)
	m.Use(render.Renderer(render.Options{
		Layout: "layout",
	}))

	m.Get("/", do_home)
	m.Get("/open", do_open)
	m.Run()
}
