package helpers_test

import (
	"github.com/mclark4386/dt_benchmark/models"

	"testing"

	"github.com/mclark4386/dt_benchmark/helpers"
)

func Test_IsSuperAdmin_Admin(t *testing.T) {
	tu := &models.User{IsSuperAdmin: true}
	if !helpers.IsSuperAdmin(tu) {
		t.Error("expected user to be super admin")
	}
}

func Test_IsSuperAdmin_NotAdmin(t *testing.T) {
	tu := &models.User{IsSuperAdmin: false}
	if helpers.IsSuperAdmin(tu) {
		t.Error("expected user to NOT be super admin")
	}
}
