package helpers

import (
	"cpsg-git.mattclark.guru/highlands/dt_benchmark/models"
	"github.com/gobuffalo/uuid"
)

func IsSuperAdmin(user *models.User) bool {
	return user.IsSuperAdmin
}

func IsTeamAdminOrBetter(user *models.User, team_id uuid.UUID) bool {
	temp_id := uuid.Must(uuid.NewV4())
	return user.IsSuperAdmin && team_id != temp_id
}
