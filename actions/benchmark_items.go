package actions

import (
	"fmt"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/mclark4386/dt_benchmark/helpers"
	"github.com/mclark4386/dt_benchmark/models"
	"github.com/pkg/errors"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (BenchmarkItem)
// DB Table: Plural (benchmark_items)
// Resource: Plural (BenchmarkItems)
// Path: Plural (/benchmark_items)
// View Template Folder: Plural (/templates/benchmark_items/)

// BenchmarkItemsResource is the resource for the BenchmarkItem model
type BenchmarkItemsResource struct {
	buffalo.Resource
}

// List gets all BenchmarkItems. This function is mapped to the path
// GET /benchmark_items
func (v BenchmarkItemsResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	benchmarkItems := &models.BenchmarkItems{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all BenchmarkItems from the DB
	if err := q.All(benchmarkItems); err != nil {
		return errors.WithStack(err)
	}

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

	return c.Render(200, r.Auto(c, benchmarkItems))
}

// Show gets the data for one BenchmarkItem. This function is mapped to
// the path GET /benchmark_items/{benchmark_item_id}
func (v BenchmarkItemsResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty BenchmarkItem
	benchmarkItem := &models.BenchmarkItem{}

	// To find the BenchmarkItem the parameter benchmark_item_id is used.
	if err := tx.Eager().Find(benchmarkItem, c.Param("benchmark_item_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.Auto(c, benchmarkItem))
}

// New renders the form for creating a new BenchmarkItem.
// This function is mapped to the path GET /benchmark_items/new
func (v BenchmarkItemsResource) New(c buffalo.Context) error {
	if helpers.IsSuperAdminOrRedirect(c) != nil {
		return nil
	}
	return c.Render(200, r.Auto(c, &models.BenchmarkItem{}))
}

// Create adds a BenchmarkItem to the DB. This function is mapped to the
// path POST /benchmark_items
func (v BenchmarkItemsResource) Create(c buffalo.Context) error {
	if helpers.IsSuperAdminOrRedirect(c) != nil {
		return nil
	}
	// Allocate an empty BenchmarkItem
	benchmarkItem := &models.BenchmarkItem{}

	// Bind benchmarkItem to the html form elements
	if err := c.Bind(benchmarkItem); err != nil {
		return errors.WithStack(err)
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(benchmarkItem)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)

		// Render again the new.html template that the user can
		// correct the input.
		return c.Render(422, r.Auto(c, benchmarkItem))
	}

	// If there are no errors set a success message
	c.Flash().Add("success", "BenchmarkItem was created successfully")

	// and redirect to the benchmark_items index page
	return c.Render(201, r.Auto(c, benchmarkItem))
}

// Edit renders a edit form for a BenchmarkItem. This function is
// mapped to the path GET /benchmark_items/{benchmark_item_id}/edit
func (v BenchmarkItemsResource) Edit(c buffalo.Context) error {
	if helpers.IsSuperAdminOrRedirect(c) != nil {
		return nil
	}
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty BenchmarkItem
	benchmarkItem := &models.BenchmarkItem{}

	if err := tx.Find(benchmarkItem, c.Param("benchmark_item_id")); err != nil {
		return c.Error(404, err)
	} else {
		err = tx.Load(benchmarkItem)
		if err != nil {
			return c.Error(404, err)
		}
	}

	resources := models.Resources{}

	if err := tx.All(&resources); err != nil {
		fmt.Printf("ERROR pulling resources: %v\n", err)
	}

	c.Set("resources", resources)

	bi_resources := []string{}
	for i := range benchmarkItem.Resources {
		bi_resources = append(bi_resources, benchmarkItem.Resources[i].ID.String())
	}
	c.Set("bi_resources", bi_resources)

	return c.Render(200, r.Auto(c, benchmarkItem))
}

// Update changes a BenchmarkItem in the DB. This function is mapped to
// the path PUT /benchmark_items/{benchmark_item_id}
func (v BenchmarkItemsResource) Update(c buffalo.Context) error {
	if helpers.IsSuperAdminOrRedirect(c) != nil {
		return nil
	}
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	type UpdateElements struct {
		Resources []string `json:"resources"`
	}

	elements := UpdateElements{}

	if err := c.Bind(&elements); err != nil {
		return errors.WithStack(err)
	}

	// Allocate an empty BenchmarkItem
	benchmarkItem := &models.BenchmarkItem{}

	if err := tx.Find(benchmarkItem, c.Param("benchmark_item_id")); err != nil {
		return c.Error(404, err)
	}

	// Bind BenchmarkItem to the html form elements
	if err := c.Bind(benchmarkItem); err != nil {
		return errors.WithStack(err)
	}

	fmt.Printf("================\nupdate elements: %v\n", elements)

	benchmarkItem.UpdateResources(tx, elements.Resources)

	verrs, err := tx.ValidateAndUpdate(benchmarkItem)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)

		// Render again the edit.html template that the user can
		// correct the input.
		return c.Render(422, r.Auto(c, benchmarkItem))
	}

	// If there are no errors set a success message
	c.Flash().Add("success", "BenchmarkItem was updated successfully")

	// and redirect to the benchmark_items index page
	return c.Render(200, r.Auto(c, benchmarkItem))
}

// Destroy deletes a BenchmarkItem from the DB. This function is mapped
// to the path DELETE /benchmark_items/{benchmark_item_id}
func (v BenchmarkItemsResource) Destroy(c buffalo.Context) error {
	if helpers.IsSuperAdminOrRedirect(c) != nil {
		return nil
	}
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty BenchmarkItem
	benchmarkItem := &models.BenchmarkItem{}

	// To find the BenchmarkItem the parameter benchmark_item_id is used.
	if err := tx.Find(benchmarkItem, c.Param("benchmark_item_id")); err != nil {
		return c.Error(404, err)
	}

	if err := tx.Destroy(benchmarkItem); err != nil {
		return errors.WithStack(err)
	}

	// If there are no errors set a flash message
	c.Flash().Add("success", "BenchmarkItem was destroyed successfully")

	// Redirect to the benchmark_items index page
	return c.Render(200, r.Auto(c, benchmarkItem))
}
