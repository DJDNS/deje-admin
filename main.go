package main

import (
    //"github.com/campadrenalin/go-deje"
    "github.com/codegangsta/martini"
)

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

func main() {
    //controller := deje.NewDEJEController()
    m := martini.Classic()

    m.Get("/", do_home)
    m.Run()
}
