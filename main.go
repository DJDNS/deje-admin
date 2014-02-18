package main

import (
	"github.com/campadrenalin/go-deje"
	"github.com/campadrenalin/go-deje/model"
	"github.com/codegangsta/martini"
	"log"
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

func do_home() string {
	return `<!DOCTYPE html>
<html>
<head>
    <title>DEJE Admin Interface</title>
</head>
<body>

<h1>DEJE Admin Interface</h1>
Open document by IRC location:
<form action="/open" method="get">
    <label>IRC server host<input name="host"/></label><br/>
    <label>IRC server port<input name="port"/></label><br/>
    <label>#channel<input name="channel"/></label><br/>
    <input type="submit"/>
</form>

</body>
</html>
`
}

func do_open(r *http.Request, l *log.Logger, c *deje.DEJEController) (int, string) {
	r.ParseForm()
	port, err := strconv.ParseUint(get_form(r, "port"), 10, 32)
	if err != nil {
		return 500, err.Error()
	}
	location := model.IRCLocation{
		Host:    get_form(r, "host"),
		Port:    uint32(port),
		Channel: get_form(r, "channel"),
	}
	doc := c.GetDocument(location)

	length := strconv.FormatUint(uint64(doc.Events.Length()), 10)
	return 200, length
}

func main() {
	controller := deje.NewDEJEController()
	m := martini.Classic()
	m.Map(controller)

	m.Get("/", do_home)
	m.Get("/open", do_open)
	m.Run()
}
