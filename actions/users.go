package actions

import (
	"fmt"

	"cpsg-git.mattclark.guru/highlands/dt_benchmark/models"
	"github.com/gobuffalo/buffalo"
	"github.com/markbates/pop"
	"github.com/pkg/errors"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (User)
// DB Table: Plural (users)
// Resource: Plural (Users)
// Path: Plural (/users)
// View Template Folder: Plural (/templates/users/)

// UsersResource is the resource for the User model
type UsersResource struct {
	buffalo.Resource
}

// List gets all Users. This function is mapped to the path
// GET /users
func (v UsersResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	users := &models.Users{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all Users from the DB
	if err := q.All(users); err != nil {
		return errors.WithStack(err)
	}

	// Make Users available inside the html template
	c.Set("users", users)

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

	return c.Render(200, r.HTML("users/index.html"))
}

// Show gets the data for one User. This function is mapped to
// the path GET /users/{user_id}
func (v UsersResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty User
	user := &models.User{}

	// To find the User the parameter user_id is used.
	if err := tx.Find(user, c.Param("user_id")); err != nil {
		return c.Error(404, err)
	}

	// Make user available inside the html template
	c.Set("user", user)

	return c.Render(200, r.HTML("users/show.html"))
}

// New renders the form for creating a new User.
// This function is mapped to the path GET /users/new
func (v UsersResource) New(c buffalo.Context) error {
	// Make user available inside the html template
	c.Set("user", &models.User{})

	return c.Render(200, r.HTML("users/new.html"))
}

// Create adds a User to the DB. This function is mapped to the
// path POST /users
func (v UsersResource) Create(c buffalo.Context) error {
	// Allocate an empty User
	user := &models.User{}

	// Bind user to the html form elements
	if err := c.Bind(user); err != nil {
		return errors.WithStack(err)
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Validate the data from the html form
	verrs, err := user.Create(tx)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make user available inside the html template
		c.Set("user", user)

		// Make the errors available inside the html template
		c.Set("errors", verrs)

		fmt.Printf("%v:%v", user, verrs)
		// Render again the new.html template that the user can
		// correct the input.
		return c.Render(422, r.HTML("users/new.html"))
	}

	// If there are no errors set a success message
	c.Session().Set("current_user_id", user.ID)
	c.Flash().Add("success", "User was created successfully")

	// and redirect to the users index page
	return c.Redirect(302, "/users/%s", user.ID)
}

// Edit renders a edit form for a User. This function is
// mapped to the path GET /users/{user_id}/edit
func (v UsersResource) Edit(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty User
	user := &models.User{}

	if err := tx.Find(user, c.Param("user_id")); err != nil {
		return c.Error(404, err)
	}

	// Make user available inside the html template
	c.Set("user", user)
	return c.Render(200, r.HTML("users/edit.html"))
}

// Update changes a User in the DB. This function is mapped to
// the path PUT /users/{user_id}
func (v UsersResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty User
	user := &models.User{}

	if err := tx.Find(user, c.Param("user_id")); err != nil {
		return c.Error(404, err)
	}

	// Bind User to the html form elements
	if err := c.Bind(user); err != nil {
		return errors.WithStack(err)
	}

	verrs, err := tx.ValidateAndUpdate(user)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make user available inside the html template
		c.Set("user", user)

		// Make the errors available inside the html template
		c.Set("errors", verrs)

		// Render again the edit.html template that the user can
		// correct the input.
		return c.Render(422, r.HTML("users/edit.html"))
	}

	// If there are no errors set a success message
	c.Flash().Add("success", "User was updated successfully")

	// and redirect to the users index page
	return c.Redirect(302, "/users/%s", user.ID)
}

// Destroy deletes a User from the DB. This function is mapped
// to the path DELETE /users/{user_id}
func (v UsersResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty User
	user := &models.User{}

	// To find the User the parameter user_id is used.
	if err := tx.Find(user, c.Param("user_id")); err != nil {
		return c.Error(404, err)
	}

	if err := tx.Destroy(user); err != nil {
		return errors.WithStack(err)
	}

	if current_user, cuok := c.Value("current_user").(*models.User); cuok && user.ID == current_user.ID {
		c.Flash().Add("success", "User was destroyed successfully")
		return AuthDestroy(c)
	}

	// If there are no errors set a flash message
	c.Flash().Add("success", "User was destroyed successfully")

	// Redirect to the users index page
	return c.Redirect(302, "/users")
}

// SetCurrentUser attempts to find a user based on the current_user_id
// in the session. If one is found it is set on the context.
func SetCurrentUser(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uid := c.Session().Get("current_user_id"); uid != nil {
			u := &models.User{}
			tx, ok := c.Value("tx").(*pop.Connection)
			if !ok {
				return errors.WithStack(errors.New("no transaction found"))
			}

			err := tx.Find(u, uid)
			if err != nil {
				return errors.WithStack(err)
			}
			c.Set("current_user", u)
		}
		return next(c)
	}
}

// Authorize require a user be logged in before accessing a route
func Authorize(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uid := c.Session().Get("current_user_id"); uid == nil {
			c.Flash().Add("danger", "You must be authorized to see that page")
			return c.Redirect(302, "/")
		}
		return next(c)
	}
}
