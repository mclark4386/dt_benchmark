package helpers

import (
	"cpsg-git.mattclark.guru/highlands/dt_benchmark/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/uuid"
)

func IsSuperAdmin(user *models.User) bool {
	return user.IsSuperAdmin
}

func IsTeamAdminOrBetter(user *models.User, team_id uuid.UUID) bool {
	temp_id := uuid.Must(uuid.NewV4())
	return user.IsSuperAdmin && team_id != temp_id
}

// Template Helpers
func GetCurrentUserInTemplate(c plush.HelperContext) *models.User {
	return c.Value("current_user").(*models.User)
}

func IsLoggedInInTemplate(c plush.HelperContext) bool {
	return c.Value("current_user") != nil
}

func IsCurrentUserSuperAdmin(c plush.HelperContext) bool {
	return GetCurrentUserInTemplate(c).IsSuperAdmin
}

func IsCurrentUserTeamOrSuperAdmin(team_id uuid.UUID, c plush.HelperContext) bool {
	return GetCurrentUserInTemplate(c).IsSuperAdmin
}

// Action Helpers

func GetCurrentUser(c buffalo.Context) *models.User {
	return c.Value("current_user").(*models.User)
}

func IsSuperAdminOrRedirect(c buffalo.Context) error {
	user := GetCurrentUser(c)
	if IsSuperAdmin(user) {
		return nil
	} else {
		c.Flash().Add("danger", "You don't have permissions for that!")
		return c.Redirect(302, "/")
	}
}

func IsTeamAdminBetterOrRedirect(c buffalo.Context, team_id uuid.UUID) error {
	return IsSuperAdminOrRedirect(c)
}

func IsCampusAdminBetterOrRedirect(c buffalo.Context, campus_id uuid.UUID) error {
	return IsSuperAdminOrRedirect(c)
}

func IsTeamOrCampusAdminBetterOrRedirect(c buffalo.Context, team_id uuid.UUID, campus_id uuid.UUID) error {
	return IsSuperAdminOrRedirect(c)
}

func IsLeaderBetterOrRedirect(c buffalo.Context, user_id uuid.UUID) error {
	return IsSuperAdminOrRedirect(c)
}
