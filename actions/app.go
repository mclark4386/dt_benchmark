package actions

import (
	"net/http"
	"strings"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo-pop/pop/popmw"
	contenttype "github.com/gobuffalo/mw-contenttype"
	"github.com/gobuffalo/x/sessions"
	"github.com/rs/cors"

	// popmw "github.com/gobuffalo/buffalo-pop/pop/popmw"
	"github.com/gobuffalo/envy"
	i18n "github.com/gobuffalo/mw-i18n"
	paramlogger "github.com/gobuffalo/mw-paramlogger"

	"github.com/mclark4386/dt_benchmark/middleware"
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
//
// Routing, middleware, groups, etc... are declared TOP -> DOWN.
// This means if you add a middleware to `app` *after* declaring a
// group, that group will NOT have that new middleware. The same
// is true of resource declarations as well.
//
// It also means that routes are checked in the order they are declared.
// `ServeFiles` is a CATCH-ALL route, so it should always be
// placed last in the route declarations, as it will prevent routes
// declared after it to never be called.
func App() *buffalo.App {
	if app == nil {
		c := cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"HEAD", "GET", "POST", "PUT", "PATCH", "DELETE"},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: false,
			ExposedHeaders:   []string{"Access-Token"},
		})

		app = buffalo.New(buffalo.Options{
			Env:          ENV,
			SessionName:  "_dt_benchmark_session",
			SessionStore: sessions.Null{},
			PreWares: []buffalo.PreWare{
				func(h http.Handler) http.Handler {
					return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
						path := strings.ToLower(req.URL.Path)
						_, hasCT := req.Header["Content-Type"]
						if strings.HasSuffix(path, ".json") {
							req.URL.Path = path[:len(path)-5]
							if !hasCT {
								req.Header["Content-Type"] = []string{"json"}
							}
						} else if strings.HasSuffix(path, ".xml") {
							req.URL.Path = path[:len(path)-4]
							if !hasCT {
								req.Header["Content-Type"] = []string{"xml"}
							}
						}
					})
				},
				c.Handler,
			},
		})
		// app.Muxer().StrictSlash(false)

		if ENV == "development" {
			app.Use(paramlogger.ParameterLogger)
		}

		// Set the request content type to JSON
		app.Use(contenttype.Set("application/json"))

		// Wraps each request in a transaction.
		//  c.Value("tx").(*pop.PopTransaction)
		// Remove to disable this.
		app.Use(popmw.Transaction(models.DB))

		api := app.Group("/api/v1")

		api.Use(middleware.TokenMiddleware)

		api.Middleware.Skip(middleware.TokenMiddleware, AuthCreate, AuthDestroy)

		api.POST("/login", AuthCreate)
		api.DELETE("/logout", AuthDestroy)

		usr := &UsersResource{}
		api.Middleware.Skip(middleware.TokenMiddleware, usr.Create, usr.New)
		api.Resource("/users", UsersResource{&buffalo.BaseResource{}})
		page := &PagesResource{}
		api.Middleware.Skip(middleware.TokenMiddleware, page.Show)
		api.Resource("/pages", PagesResource{&buffalo.BaseResource{}})
		api.Resource("/teams", TeamsResource{})
		api.Resource("/resources", ResourcesResource{})
		api.Resource("/benchmarks", BenchmarksResource{})
		api.Resource("/benchmark_items", BenchmarkItemsResource{})
		camp := &CampusesResource{}
		api.Middleware.Skip(middleware.TokenMiddleware, camp.Show, camp.List)
		api.Resource("/campuses", CampusesResource{&buffalo.BaseResource{}})
		api.Resource("/team_positions", TeamPositionsResource{})
		//end api

		// app.ServeFiles("/", assetsBox)
		// Custom static file server (replaces app.ServeFiles)
		configureAssetsBoxSeparator()
		app.GET("/{asset:.*}", StaticFileGet)
		app.HEAD("/{asset:.*}", StaticFileGet)
	}

	return app
}
