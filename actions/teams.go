package actions

import (
	"cpsg-git.mattclark.guru/highlands/dt_benchmark/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (Team)
// DB Table: Plural (teams)
// Resource: Plural (Teams)
// Path: Plural (/teams)
// View Template Folder: Plural (/templates/teams/)

// TeamsResource is the resource for the Team model
type TeamsResource struct {
	buffalo.Resource
}

// List gets all Teams. This function is mapped to the path
// GET /teams
func (v TeamsResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	teams := &models.Teams{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all Teams from the DB
	if err := q.All(teams); err != nil {
		return errors.WithStack(err)
	}

	// Make Teams available inside the html template
	c.Set("teams", teams)

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

	return c.Render(200, r.HTML("teams/index.html"))
}

// Show gets the data for one Team. This function is mapped to
// the path GET /teams/{team_id}
func (v TeamsResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Team
	team := &models.Team{}

	// To find the Team the parameter team_id is used.
	if err := tx.Find(team, c.Param("team_id")); err != nil {
		return c.Error(404, err)
	}

	// Make team available inside the html template
	c.Set("team", team)

	return c.Render(200, r.HTML("teams/show.html"))
}

// New renders the form for creating a new Team.
// This function is mapped to the path GET /teams/new
func (v TeamsResource) New(c buffalo.Context) error {
	// Make team available inside the html template
	c.Set("team", &models.Team{})

	return c.Render(200, r.HTML("teams/new.html"))
}

// Create adds a Team to the DB. This function is mapped to the
// path POST /teams
func (v TeamsResource) Create(c buffalo.Context) error {
	// Allocate an empty Team
	team := &models.Team{}

	// Bind team to the html form elements
	if err := c.Bind(team); err != nil {
		return errors.WithStack(err)
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(team)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make team available inside the html template
		c.Set("team", team)

		// Make the errors available inside the html template
		c.Set("errors", verrs)

		// Render again the new.html template that the user can
		// correct the input.
		return c.Render(422, r.HTML("teams/new.html"))
	}

	// If there are no errors set a success message
	c.Flash().Add("success", "Team was created successfully")

	// and redirect to the teams index page
	return c.Redirect(302, "/teams/%s", team.ID)
}

// Edit renders a edit form for a Team. This function is
// mapped to the path GET /teams/{team_id}/edit
func (v TeamsResource) Edit(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Team
	team := &models.Team{}

	if err := tx.Find(team, c.Param("team_id")); err != nil {
		return c.Error(404, err)
	}

	// Make team available inside the html template
	c.Set("team", team)
	return c.Render(200, r.HTML("teams/edit.html"))
}

// Update changes a Team in the DB. This function is mapped to
// the path PUT /teams/{team_id}
func (v TeamsResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Team
	team := &models.Team{}

	if err := tx.Find(team, c.Param("team_id")); err != nil {
		return c.Error(404, err)
	}

	// Bind Team to the html form elements
	if err := c.Bind(team); err != nil {
		return errors.WithStack(err)
	}

	verrs, err := tx.ValidateAndUpdate(team)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make team available inside the html template
		c.Set("team", team)

		// Make the errors available inside the html template
		c.Set("errors", verrs)

		// Render again the edit.html template that the user can
		// correct the input.
		return c.Render(422, r.HTML("teams/edit.html"))
	}

	// If there are no errors set a success message
	c.Flash().Add("success", "Team was updated successfully")

	// and redirect to the teams index page
	return c.Redirect(302, "/teams/%s", team.ID)
}

// Destroy deletes a Team from the DB. This function is mapped
// to the path DELETE /teams/{team_id}
func (v TeamsResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Team
	team := &models.Team{}

	// To find the Team the parameter team_id is used.
	if err := tx.Find(team, c.Param("team_id")); err != nil {
		return c.Error(404, err)
	}

	if err := tx.Destroy(team); err != nil {
		return errors.WithStack(err)
	}

	// If there are no errors set a flash message
	c.Flash().Add("success", "Team was destroyed successfully")

	// Redirect to the teams index page
	return c.Redirect(302, "/teams")
}
