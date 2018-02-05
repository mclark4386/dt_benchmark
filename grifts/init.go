package grifts

import (
	"cpsg-git.mattclark.guru/highlands/dt_benchmark/actions"
	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}
