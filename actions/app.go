package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/middleware"
	"github.com/gobuffalo/envy"

	"github.com/gobuffalo/buffalo/middleware/csrf"
	"github.com/gobuffalo/buffalo/middleware/i18n"
	"github.com/gobuffalo/packr"
	"github.com/mclark4386/buffalo_helpers"
	"github.com/mclark4386/dt_benchmark/models"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App
var T *i18n.Translator

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:         ENV,
			SessionName: "_dt_benchmark_session",
			PreWares:    []buffalo.PreWare{buffalo_helpers.AutoSetContentType},
		})

		if ENV == "development" {
			app.Use(middleware.ParameterLogger)
		}

		// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
		// Remove to disable this.
		app.Use(csrf.New)

		// Wraps each request in a transaction.
		//  c.Value("tx").(*pop.PopTransaction)
		// Remove to disable this.
		app.Use(middleware.PopTransaction(models.DB))
		app.Use(SetCurrentUser)

		// Setup and use translations:
		var err error
		if T, err = i18n.New(packr.NewBox("../locales"), "en-US"); err != nil {
			app.Stop(err)
		}
		app.Use(T.Middleware())
		app.Use(SetupNavbar)
		app.Use(Authorize)

		app.Middleware.Skip(Authorize, HomeHandler, AuthCreate, AuthNew, AuthDestroy)
		app.GET("/", HomeHandler)

		app.ServeFiles("/assets", assetsBox)
		app.GET("/login", AuthNew)
		app.POST("/login", AuthCreate)
		app.DELETE("/logout", AuthDestroy)

		usr := &UsersResource{}
		app.Middleware.Skip(Authorize, usr.Create, usr.New)
		app.Resource("/users", UsersResource{&buffalo.BaseResource{}})
		page := &PagesResource{}
		app.Middleware.Skip(Authorize, page.Show)
		app.Resource("/pages", PagesResource{&buffalo.BaseResource{}})
		app.Resource("/teams", TeamsResource{})
		app.Resource("/resources", ResourcesResource{})
		app.Resource("/benchmarks", BenchmarksResource{})
		app.Resource("/benchmark_items", BenchmarkItemsResource{})
		camp := &CampusesResource{}
		app.Middleware.Skip(Authorize, camp.Show, camp.Index)
		app.Resource("/campuses", CampusesResource{&buffalo.BaseResource{}})
		app.Resource("/team_positions", TeamPositionsResource{})
	}

	return app
}
