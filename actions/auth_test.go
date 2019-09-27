package actions

import (
	"fmt"

	"github.com/mclark4386/dt_benchmark/models"
)

func (as *ActionSuite) Test_Auth_Create() {
	u := &models.User{
		Email:                "mark@example.com",
		Password:             "password",
		PasswordConfirmation: "password",
	}
	pass := u.Password
	verrs, err := u.Create(as.DB)
	as.NoError(err)
	as.False(verrs.HasAny())

	u.Password = pass
	res := as.JSON("/api/v1/login").Post(u)
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), "mark@example.com")
	headerKeys := make([]string, 0, len(res.Header()))
	for k, v := range res.Header() {
		fmt.Printf("\nk:%v  v:%v\n", k, v)
		headerKeys = append(headerKeys, k)
	}
	as.Contains(headerKeys, "Access-Token")
}

func (as *ActionSuite) Test_Auth_Create_UnknownUser() {
	u := &models.User{
		Email:    "mark@example.com",
		Password: "password",
	}
	res := as.JSON("/api/v1/login").Post(u)
	as.Equal(422, res.Code)
	as.Contains(res.Body.String(), "invalid email/password")
}

func (as *ActionSuite) Test_Auth_Create_BadPassword() {
	u := &models.User{
		Email:                "mark@example.com",
		Password:             "password",
		PasswordConfirmation: "password",
	}
	verrs, err := u.Create(as.DB)
	as.NoError(err)
	as.False(verrs.HasAny())

	u.Password = "bad"
	res := as.JSON("/api/v1/login").Post(u)
	as.Equal(422, res.Code)
	as.Contains(res.Body.String(), "invalid email/password")
}
