package actions

import (
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/packr/v2"
	"github.com/mclark4386/dt_benchmark/helpers"
)

var r *render.Engine
var assetsBox = packr.New("../public/assets", "../public/assets")

func init() {
	r = render.New(render.Options{
		// HTML layout to be used for all HTML requests:
		HTMLLayout: "application.html",

		// Box containing all of the templates:
		TemplatesBox: packr.New("../templates", "../templates"),
		AssetsBox:    assetsBox,

		// Add template helpers here:
		Helpers: render.Helpers{
			"isCurrentUserSuperAdmin":       helpers.IsCurrentUserSuperAdmin,
			"getCurrentUser":                helpers.GetCurrentUserInTemplate,
			"isLoggedIn":                    helpers.IsLoggedInInTemplate,
			"isCurrentUserTeamOrSuperAdmin": helpers.IsCurrentUserTeamOrSuperAdmin,
			"arrayContains":                 helpers.StringArrayContains,
			"append":                        helpers.ArrayAppend,
		},
	})
}
