package main

import (
	"github.com/campadrenalin/go-deje"
	"github.com/campadrenalin/go-deje/model"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	"net/http"
	"strconv"
)

func get_form(r *http.Request, key string) string {
	values := r.Form[key]
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func get_location(r *http.Request) (*model.IRCLocation, error) {
	r.ParseForm()
	port, err := strconv.ParseUint(get_form(r, "port"), 10, 32)
	if err != nil {
		return nil, err
	}

	location := &model.IRCLocation{
		Host:    get_form(r, "host"),
		Port:    uint32(port),
		Channel: get_form(r, "channel"),
	}

    return location, nil
}

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
