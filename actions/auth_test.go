package actions

import "cpsg-git.mattclark.guru/highlands/dt_benchmark/models"

func (as *ActionSuite) Test_Auth_New() {
	res := as.HTML("/login").Get()
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), "Sign In")
}

func (as *ActionSuite) Test_Auth_Create() {
	u := &models.User{
		Email:                "mark@example.com",
		Password:             "password",
		PasswordConfirmation: "password",
	}
	verrs, err := u.Create(as.DB)
	as.NoError(err)
	as.False(verrs.HasAny())

	res := as.HTML("/login").Post(u)
	as.Equal(302, res.Code)
	as.Equal("/", res.Location())
}

func (as *ActionSuite) Test_Auth_Create_UnknownUser() {
	u := &models.User{
		Email:    "mark@example.com",
		Password: "password",
	}
	res := as.HTML("/login").Post(u)
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
	res := as.HTML("/login").Post(u)
	as.Equal(422, res.Code)
	as.Contains(res.Body.String(), "invalid email/password")
}
