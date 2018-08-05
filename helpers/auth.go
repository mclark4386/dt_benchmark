package helpers

import (
	"cpsg-git.mattclark.guru/highlands/dt_benchmark/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/uuid"
)

// IsSuperAdmin will return if the user is a super admin or not
func IsSuperAdmin(user *models.User) bool {
	return user.IsSuperAdmin
}

// IsTeamAdminOrBetter will return if the user is a super admin or
// the admin for a specific team or not
func IsTeamAdminOrBetter(user *models.User, team_id uuid.UUID) bool {
	isTeamAdmin := false
	for i := range user.TeamsIAdmin {
		if user.TeamsIAdmin[i].ID == team_id {
			isTeamAdmin = true
			break
		}
	}
	return user.IsSuperAdmin || isTeamAdmin
}

func IsCampusAdminOrBetter(user *models.User, campus_id uuid.UUID) bool {
	isCampusAdmin := false
	for _, campus := range user.CampusesIAdmin {
		if campus.ID == campus_id {
			isCampusAdmin = true
			break
		}
	}
	return user.IsSuperAdmin || isCampusAdmin
}

// IsAnyTeamAdminOrBetter will return if the user is a super admin or
// the admin for any team or not
func IsAnyTeamAdminOrBetter(user *models.User) bool {
	return user.IsSuperAdmin || len(user.TeamsIAdmin) > 0
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
	return IsTeamAdminOrBetter(GetCurrentUserInTemplate(c), team_id)
}

func IsCurrentUserCampusOrSuperAdmin(campus_id uuid.UUID, c plush.HelperContext) bool {
	return IsCampusAdminOrBetter(GetCurrentUserInTemplate(c), campus_id)
}

func IsCurrentUserAnyTeamOrSuperAdmin(c plush.HelperContext) bool {
	return IsAnyTeamAdminOrBetter(GetCurrentUserInTemplate(c))
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
	user := GetCurrentUser(c)
	if IsTeamAdminOrBetter(user, team_id) {
		return nil
	} else {
		c.Flash().Add("danger", "You don't have permissions for that!")
		return c.Redirect(302, "/")
	}
}

func IsAnyTeamAdminBetterOrRedirect(c buffalo.Context) error {
	user := GetCurrentUser(c)
	if IsAnyTeamAdminOrBetter(user) {
		return nil
	} else {
		c.Flash().Add("danger", "You don't have permissions for that!")
		return c.Redirect(302, "/")
	}
}

func IsCampusAdminBetterOrRedirect(c buffalo.Context, campus_id uuid.UUID) error {
	user := GetCurrentUser(c)
	if IsCampusAdminOrBetter(user, campus_id) {
		return nil
	} else {
		c.Flash().Add("danger", "You don't have permissions for that!")
		return c.Redirect(302, "/")
	}
}

func IsTeamOrCampusAdminBetterOrRedirect(c buffalo.Context, team_id, campus_id uuid.UUID) error {
	return IsTeamAdminBetterOrRedirect(c, team_id)
}

func IsLeaderBetterOrRedirect(c buffalo.Context, user_id uuid.UUID) error {
	return IsSuperAdminOrRedirect(c)
}
