package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/kimgabz/booking-app-go/pkg/config"
	"github.com/kimgabz/booking-app-go/pkg/handlers"
	"github.com/kimgabz/booking-app-go/pkg/render"
)

const portNumber = ":8080"

var app config.AppConfig

// Not the same as "session" initilized in main() below.
var session *scs.SessionManager

func main() {
	// Change this value to true in production.
	app.InProduction = false

	// Session information.
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// Template cache
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	// Change "UseCache" value to "true" in production.
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	// Start web server.
	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
