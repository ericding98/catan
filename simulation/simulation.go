package simulation

import (
	"github.com/ericding98/catan"
)

type Result struct {
	Game catan.Game
}

func Run(g catan.Game) (*Result, error) {
	if g.Done() {
		g.Reset(nil)
	}

	if err := g.Setup(); err != nil {
		return nil, err
	}

	for !g.Done() {
		if err := g.Step(); err != nil {
			return nil, err
		}
	}

	res := &Result{
		Game: g,
	}

	return res, nil
}
