package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/syke99/oasis"
)

var addr = ":3000"

type Person struct {
	Name     string
	Age      int
	Birthday time.Time
}

func main() {
	// create a new Chi Router
	r := chi.NewRouter()

	// add data-only endpoints here

	// upgrade your Chi Router to an Oasis Router
	router := oasis.UpgradeRouter(r)

	// add Oasis Endpoints. Handlers are mapped to
	// specific HTTP Methods, and the HandlerFunc
	// for each Handler can be used to create
	// payloads to be added to the Island for
	// rendering, i.e. calling a database, doing
	// operations on the Handler-specific Island
	// (example below), etc.
	router.AddEndpoint(oasis.Endpoint{
		Route: "/greeting",
		Handlers: map[oasis.HTTPMethod]oasis.HandlerWithMiddleware{
			oasis.MethodGet: {
				HandlerFunc: func(w http.ResponseWriter, r *http.Request) {

				},
				Island: helloIsland(),
			},
		},
	}).AddEndpoint(oasis.Endpoint{
		Route: "/greeting/hi",
		Handlers: map[oasis.HTTPMethod]oasis.HandlerWithMiddleware{
			oasis.MethodGet: {
				HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
					// you can obtain the props
					// for an Island for a specific
					// route by passing the request
					// to PropsForRequest. You can
					// then access each prop by the
					// key you assigned it whenever
					// calling (Island).AddProp()
					// or (Island).AddProps()
					p := oasis.PropsForRequest(r)["person"].(Person)

					isBday := isBirthday(time.Now(), p.Birthday)

					// payloads are additional data
					// points to be rendered to an
					// Island besides values stored
					// in the Island's props. This is
					// useful for things like making
					// calls to a database for additional
					// data to render, etc.
					payload := oasis.NewPayload()

					if isBday {
						payload.Set("age", p.Age+1)
						payload.Set("isBirthday", isBday)
					} else {
						payload.Set("isBirthday", isBday)
					}

					// after adding data to the payload,
					// marshal it so it can be written
					// back to the Island
					bt, _ := payload.Marshal()

					// simply write your response as normal
					_, _ = w.Write(bt)
				},
				Island: greetingIsland(),
			},
		},
	})

	// server your Oasis Router
	err := http.ListenAndServe(addr, router)
	if err != nil {
		log.Fatal(err)
	}
}

func isBirthday(today time.Time, birthday time.Time) bool {
	currentMonth := today.Month()
	currentDay := today.Day()

	bdayMonth := birthday.Month()
	bdayDay := birthday.Day()

	if (currentMonth == bdayMonth) && (currentDay == bdayDay) {
		return true
	}

	return false
}

func helloIsland() oasis.Island {
	helloDiv := oasis.NewIsland("navbar", hello)

	helloDiv.AddProps(map[string]any{
		"ID":      "helloID",
		"classes": "hello greeting",
		"url":     "/greeting/hi",
		"target":  "this",
		"swap":    "innerHTML",
		"trigger": "load",
	})

	return helloDiv
}

func greetingIsland() oasis.Island {
	greetingDiv := oasis.NewIsland("greeting", hi)

	greetingDiv.AddProp("person", Person{
		Name:     "Jane Doe",
		Age:      28,
		Birthday: time.Date(1995, 12, 31, 0, 0, 0, 0, time.Local),
	})

	return greetingDiv
}
