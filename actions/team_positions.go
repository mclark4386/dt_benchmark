package actions

import (
	"cpsg-git.mattclark.guru/highlands/dt_benchmark/helpers"
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
// Model: Singular (TeamPosition)
// DB Table: Plural (team_positions)
// Resource: Plural (TeamPositions)
// Path: Plural (/team_positions)
// View Template Folder: Plural (/templates/team_positions/)

// TeamPositionsResource is the resource for the TeamPosition model
type TeamPositionsResource struct {
	buffalo.Resource
}

// List gets all TeamPositions. This function is mapped to the path
// GET /team_positions
func (v TeamPositionsResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	teamPositions := &models.TeamPositions{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all TeamPositions from the DB
	if err := q.Eager().All(teamPositions); err != nil {
		return errors.WithStack(err)
	}

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

	return c.Render(200, r.Auto(c, teamPositions))
}

// Show gets the data for one TeamPosition. This function is mapped to
// the path GET /team_positions/{team_position_id}
func (v TeamPositionsResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty TeamPosition
	teamPosition := &models.TeamPosition{}

	// To find the TeamPosition the parameter team_position_id is used.
	if err := tx.Eager().Find(teamPosition, c.Param("team_position_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.Auto(c, teamPosition))
}

// New renders the form for creating a new TeamPosition.
// This function is mapped to the path GET /team_positions/new
func (v TeamPositionsResource) New(c buffalo.Context) error {
	if helpers.IsAnyTeamAdminBetterOrRedirect(c) != nil {
		return nil
	}
	v.SetupForms(c)
	return c.Render(200, r.Auto(c, &models.TeamPosition{}))
}

// Create adds a TeamPosition to the DB. This function is mapped to the
// path POST /team_positions
func (v TeamPositionsResource) Create(c buffalo.Context) error {
	// Allocate an empty TeamPosition
	teamPosition := &models.TeamPosition{}

	// Bind teamPosition to the html form elements
	if err := c.Bind(teamPosition); err != nil {
		return errors.WithStack(err)
	}

	// Make sure the current_user is either a super admin or a team admin for the team
	// we are trying to add this position to.
	if helpers.IsTeamAdminBetterOrRedirect(c, teamPosition.TeamID) != nil {
		return nil
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(teamPosition)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)

		// Render again the new.html template that the user can
		// correct the input.
		return c.Render(422, r.Auto(c, teamPosition))
	}

	// If there are no errors set a success message
	c.Flash().Add("success", "TeamPosition was created successfully")

	// and redirect to the team_positions index page
	return c.Render(201, r.Auto(c, teamPosition))
}

// Edit renders a edit form for a TeamPosition. This function is
// mapped to the path GET /team_positions/{team_position_id}/edit
func (v TeamPositionsResource) Edit(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty TeamPosition
	teamPosition := &models.TeamPosition{}

	if err := tx.Eager().Find(teamPosition, c.Param("team_position_id")); err != nil {
		return c.Error(404, err)
	}

	// Make sure the current_user is either a super admin or a team admin for the team
	// we are part of.
	if helpers.IsTeamAdminBetterOrRedirect(c, teamPosition.TeamID) != nil {
		return nil
	}

	v.SetupForms(c)
	return c.Render(200, r.Auto(c, teamPosition))
}

// Update changes a TeamPosition in the DB. This function is mapped to
// the path PUT /team_positions/{team_position_id}
func (v TeamPositionsResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty TeamPosition
	teamPosition := &models.TeamPosition{}

	if err := tx.Find(teamPosition, c.Param("team_position_id")); err != nil {
		return c.Error(404, err)
	}

	// Make sure the current_user is either a super admin or a team admin for the team
	// we are part of.
	if helpers.IsTeamAdminBetterOrRedirect(c, teamPosition.TeamID) != nil {
		return nil
	}

	// Bind TeamPosition to the html form elements
	if err := c.Bind(teamPosition); err != nil {
		return errors.WithStack(err)
	}

	verrs, err := tx.ValidateAndUpdate(teamPosition)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)

		// Render again the edit.html template that the user can
		// correct the input.
		return c.Render(422, r.Auto(c, teamPosition))
	}

	// If there are no errors set a success message
	c.Flash().Add("success", "TeamPosition was updated successfully")

	// and redirect to the team_positions index page
	return c.Render(200, r.Auto(c, teamPosition))
}

// Destroy deletes a TeamPosition from the DB. This function is mapped
// to the path DELETE /team_positions/{team_position_id}
func (v TeamPositionsResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty TeamPosition
	teamPosition := &models.TeamPosition{}

	// To find the TeamPosition the parameter team_position_id is used.
	if err := tx.Find(teamPosition, c.Param("team_position_id")); err != nil {
		return c.Error(404, err)
	}

	// Make sure the current_user is either a super admin or a team admin for the team
	// we are part of.
	if helpers.IsTeamAdminBetterOrRedirect(c, teamPosition.TeamID) != nil {
		return nil
	}

	if err := tx.Destroy(teamPosition); err != nil {
		return errors.WithStack(err)
	}

	// If there are no errors set a flash message
	c.Flash().Add("success", "TeamPosition was destroyed successfully")

	// Redirect to the team_positions index page
	return c.Render(200, r.Auto(c, teamPosition))
}

func (v TeamPositionsResource) SetupForms(c buffalo.Context) {
	teams := models.Teams{}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		c.Set("teams", teams)
		return
	}

	if err := tx.All(&teams); err != nil {
		c.Set("teams", teams)
		return
	}
	c.Set("teams", teams)
}
