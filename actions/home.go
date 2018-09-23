package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/mclark4386/dt_benchmark/models"
)

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {
	homeSlug := "home"
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if ok {
		// Allocate an empty Page
		page := &models.Page{}

		// To find the Page the parameter page_id is used.
		if err := tx.Where("slug = ?", homeSlug).First(page); err == nil {
			c.Set("home", page)
		}
	}

	return c.Render(200, r.HTML("index.html"))
}
