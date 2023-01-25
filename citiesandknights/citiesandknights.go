package citiesandknights

import (
	"crypto/rand"
	"math/big"

	"github.com/ericding98/catan"
)

const minPoints int = 13

func NewGame() catan.Game {
	return catan.NewGame(NewGameState(), minPoints)
}

type GameState struct {
	catan.UnimplementedGameState
}

func NewGameState() catan.GameState {
	return &GameState{
		UnimplementedGameState: catan.NewGameState(),
	}
}

func (s *GameState) initializePlayerOrder() error {
	ps := s.Players()
	pCount := len(ps)
	for range ps {
		i, err := rand.Int(rand.Reader, big.NewInt(int64(pCount)))
		if err != nil {
			return err
		}

		idx := int(i.Int64())
		s.PopulateInOrder(ps[idx])
	}
	return nil
}

func (s *GameState) initializePlayerPlacements() error {
	ps := s.Players()
	for _, p := range ps {
		s := p.PlaceInitialSettlementOn(s.Board())
		p.PlaceRoadAnywhereNear(s)
	}

	for i := len(ps) - 1; i >= 0; i-- {
		p := ps[i]
		c := p.PlaceInitialCityOn(s.Board())
		if _, err := c.ProduceSurroundingResources(); err != nil {
			return err
		}

		p.PlaceRoadAnywhereNear(c)
	}
	return nil
}

func (s *GameState) Setup() error {
	if err := s.initializePlayerOrder(); err != nil {
		return err
	}
	if err := s.initializePlayerPlacements(); err != nil {
		return err
	}
	return nil
}

func (s *GameState) Step() error {
	return nil
}

func (s *GameState) Reset(b catan.Board) {
	s.UnimplementedGameState.Reset(b)

	// TODO: implement reset
}
