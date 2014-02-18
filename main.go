package main

import (
	"github.com/campadrenalin/go-deje"
	"github.com/codegangsta/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"strconv"
)

func do_home(r render.Render) {
    r.HTML(200, "root", nil)
}

func do_open(r *http.Request, c *deje.DEJEController) (int, string) {
    location, err := get_location(r)
    if err != nil {
        return 500, err.Error()
    }

	doc := c.GetDocument(*location)
	length := strconv.FormatUint(uint64(doc.Events.Length()), 10)
	return 200, length
}

func main() {
	controller := deje.NewDEJEController()
	m := martini.Classic()
	m.Map(controller)
    m.Use(render.Renderer())

	m.Get("/", do_home)
	m.Get("/open", do_open)
	m.Run()
}
