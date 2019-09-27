package actions

import "github.com/mclark4386/dt_benchmark/models"

func (as *ActionSuite) Test_UsersResource_List() {
	//	as.Fail("Not Implemented!")
}

func (as *ActionSuite) Test_UsersResource_Show() {
	//	as.Fail("Not Implemented!")
}

func (as *ActionSuite) Test_UsersResource_Edit() {
	//	as.Fail("Not Implemented!")
}

func (as *ActionSuite) Test_UsersResource_Update() {
	//	as.Fail("Not Implemented!")
}

func (as *ActionSuite) Test_UsersResource_Destroy() {
	//	as.Fail("Not Implemented!")
}

func (as *ActionSuite) Test_Users_Create() {
	count, err := as.DB.Count("users")
	as.NoError(err)
	as.Equal(0, count)

	u := &models.User{
		Email:                "mark@example.com",
		Password:             "password",
		PasswordConfirmation: "password",
	}

	res := as.JSON("/api/v1/users").Post(u)
	as.Equal(302, res.Code)

	count, err = as.DB.Count("users")
	as.NoError(err)
	as.Equal(1, count)
}
