package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/mclark4386/dt_benchmark/models"
	"github.com/pkg/errors"
)

func SetupNavbar(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		//Load Teams
		teams := &models.Teams{}
		tx, ok := c.Value("tx").(*pop.Connection)
		if !ok {
			return errors.WithStack(errors.New("no transaction found"))
		}

		err := tx.All(teams)
		if err != nil {
			return errors.WithStack(err)
		}
		c.Set("teams", teams)

		//Load Campuses
		campuses := &models.Campuses{}

		err = tx.All(campuses)
		if err != nil {
			return errors.WithStack(err)
		}
		c.Set("campuses", campuses)

		return next(c)
	}
}
