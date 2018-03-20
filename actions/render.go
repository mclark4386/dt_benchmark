package actions

import (
	"cpsg-git.mattclark.guru/highlands/dt_benchmark/helpers"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/packr"
)

var r *render.Engine
var assetsBox = packr.NewBox("../public/assets")

func init() {
	r = render.New(render.Options{
		// HTML layout to be used for all HTML requests:
		HTMLLayout: "application.html",

		// Box containing all of the templates:
		TemplatesBox: packr.NewBox("../templates"),
		AssetsBox:    assetsBox,

		// Add template helpers here:
		Helpers: render.Helpers{
			"isCurrentUserSuperAdmin":       helpers.IsCurrentUserSuperAdmin,
			"getCurrentUser":                helpers.GetCurrentUserInTemplate,
			"isLoggedIn":                    helpers.IsLoggedInInTemplate,
			"isCurrentUserTeamOrSuperAdmin": helpers.IsCurrentUserTeamOrSuperAdmin,
		},
	})
}
