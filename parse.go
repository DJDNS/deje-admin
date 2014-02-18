package main

import (
	"github.com/campadrenalin/go-deje/model"
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
