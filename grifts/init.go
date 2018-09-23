package grifts

import (
	"github.com/gobuffalo/buffalo"
	"github.com/mclark4386/dt_benchmark/actions"
)

func init() {
	buffalo.Grifts(actions.App())
}
