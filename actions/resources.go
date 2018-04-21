package actions

import (
	"cpsg-git.mattclark.guru/highlands/dt_benchmark/helpers"
	"cpsg-git.mattclark.guru/highlands/dt_benchmark/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/pkg/errors"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (Resource)
// DB Table: Plural (resources)
// Resource: Plural (Resources)
// Path: Plural (/resources)
// View Template Folder: Plural (/templates/resources/)

// ResourcesResource is the resource for the Resource model
type ResourcesResource struct {
	buffalo.Resource
}

// List gets all Resources. This function is mapped to the path
// GET /resources
func (v ResourcesResource) List(c buffalo.Context) error {
	//TODO: Should this be team or campus or better? would need to be able to pass a list of teams

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	resources := &models.Resources{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all Resources from the DB
	if err := q.All(resources); err != nil {
		return errors.WithStack(err)
	}

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

	return c.Render(200, r.Auto(c, resources))
}

// Show gets the data for one Resource. This function is mapped to
// the path GET /resources/{resource_id}
func (v ResourcesResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Resource
	resource := &models.Resource{}

	// To find the Resource the parameter resource_id is used.
	if err := tx.Find(resource, c.Param("resource_id")); err != nil {
		return c.Error(404, err)
	}

	//TODO: Should this be team or campus or better? would need to be able to pass a list of teams

	return c.Render(200, r.Auto(c, resource))
}

// New renders the form for creating a new Resource.
// This function is mapped to the path GET /resources/new
func (v ResourcesResource) New(c buffalo.Context) error {
	//TODO: Should this be team or campus or better? would need to be able to pass a list of teams
	if helpers.IsSuperAdminOrRedirect(c) != nil {
		return nil
	}

	return c.Render(200, r.Auto(c, &models.Resource{}))
}

// Create adds a Resource to the DB. This function is mapped to the
// path POST /resources
func (v ResourcesResource) Create(c buffalo.Context) error {
	//TODO: Should this be team or campus or better? would need to be able to pass a list of teams
	if helpers.IsTeamOrCampusAdminBetterOrRedirect(c, uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4())) != nil {
		return nil
	}

	// Allocate an empty Resource
	resource := &models.Resource{}

	// Bind resource to the html form elements
	if err := c.Bind(resource); err != nil {
		return errors.WithStack(err)
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(resource)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)

		// Render again the new.html template that the user can
		// correct the input.
		return c.Render(422, r.Auto(c, resource))
	}

	// If there are no errors set a success message
	c.Flash().Add("success", "Resource was created successfully")

	// and redirect to the resources index page
	return c.Render(201, r.Auto(c, resource))
}

// Edit renders a edit form for a Resource. This function is
// mapped to the path GET /resources/{resource_id}/edit
func (v ResourcesResource) Edit(c buffalo.Context) error {
	//TODO: Should this be team or campus or better? would need to be able to pass a list of teams
	if helpers.IsSuperAdminOrRedirect(c) != nil {
		return nil
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Resource
	resource := &models.Resource{}

	if err := tx.Find(resource, c.Param("resource_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.Auto(c, resource))
}

// Update changes a Resource in the DB. This function is mapped to
// the path PUT /resources/{resource_id}
func (v ResourcesResource) Update(c buffalo.Context) error {
	//TODO: Should this be team or campus or better? would need to be able to pass a list of teams
	if helpers.IsSuperAdminOrRedirect(c) != nil {
		return nil
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Resource
	resource := &models.Resource{}

	if err := tx.Find(resource, c.Param("resource_id")); err != nil {
		return c.Error(404, err)
	}

	// Bind Resource to the html form elements
	if err := c.Bind(resource); err != nil {
		return errors.WithStack(err)
	}

	verrs, err := tx.ValidateAndUpdate(resource)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)

		// Render again the edit.html template that the user can
		// correct the input.
		return c.Render(422, r.Auto(c, resource))
	}

	// If there are no errors set a success message
	c.Flash().Add("success", "Resource was updated successfully")

	// and redirect to the resources index page
	return c.Render(200, r.Auto(c, resource))
}

// Destroy deletes a Resource from the DB. This function is mapped
// to the path DELETE /resources/{resource_id}
func (v ResourcesResource) Destroy(c buffalo.Context) error {
	//TODO: Should this be team or campus or better? would need to be able to pass a list of teams
	if helpers.IsSuperAdminOrRedirect(c) != nil {
		return nil
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Resource
	resource := &models.Resource{}

	// To find the Resource the parameter resource_id is used.
	if err := tx.Find(resource, c.Param("resource_id")); err != nil {
		return c.Error(404, err)
	}

	if err := tx.Destroy(resource); err != nil {
		return errors.WithStack(err)
	}

	// If there are no errors set a flash message
	c.Flash().Add("success", "Resource was destroyed successfully")

	// Redirect to the resources index page
	return c.Render(200, r.Auto(c, resource))
}
