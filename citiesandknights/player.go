package citiesandknights

import "github.com/ericding98/catan"

type Strategy interface {
	PlaceInitialSettlementOn(catan.Board) *catan.Settlement
	PlaceInitialCityOn(catan.Board) *City
	PlaceRoadAnywhereNear(catan.Piece) *catan.Road
}

type Player interface {
	catan.UnimplementedPlayer
	Strategy
	Points() int
	Reset(Strategy, catan.Board)
	NewSettlement() *catan.Settlement
	NewCity() *City
	NewRoad() *catan.Road
	HasLongestRoad() bool
	LongestRoadLength() int
	HasLargestArmy() bool
	ArmySize() int
}

type player struct {
	Strategy
	catan.UnimplementedPlayer
	vps               int
	hasLongestRoad    bool
	longestRoadLength int
	hasLargestArmy    bool
	armySize          int
}

func NewPlayer(p catan.UnimplementedPlayer, s Strategy) Player {
	return &player{
		Strategy:            s,
		UnimplementedPlayer: p,
		vps:                 0,
		hasLongestRoad:      false,
		longestRoadLength:   0,
		hasLargestArmy:      false,
		armySize:            0,
	}
}

func (p *player) Points() int {
	var settlementCount, cityCount int

	p.UnimplementedPlayer.Board().IterV(func(v *catan.Vertex) bool {
		if v.Piece.Owner().ID() != p.ID() {
			return true
		}

		switch v.Piece.(type) {
		case *catan.Settlement:
			settlementCount++
		case *City:
			cityCount++
		}
		return true
	})

	vps := p.vps

	lrb, lab := p.calculateLongestRoadBonus(), p.calculateLargestArmyBonus()

	return settlementCount + cityCount*2 + vps + lrb + lab
}

func (p *player) calculateLongestRoadBonus() (n int) {
	if p.hasLongestRoad {
		n = 2
	}
	return
}

func (p *player) calculateLargestArmyBonus() (n int) {
	if p.hasLargestArmy {
		n = 2
	}
	return
}

func (p *player) Reset(s Strategy, b catan.Board) {
	p.Strategy = s
	p.vps = 0
}

func (p *player) NewSettlement() *catan.Settlement {
	return catan.NewSettlement(p)
}

func (p *player) NewCity() *City {
	return &City{
		id:    UID(),
		owner: p,
		b:     p.UnimplementedPlayer.Board(),
	}
}

func (p *player) NewRoad() *catan.Road {
	return catan.NewRoad(p)
}

func (p *player) HasLongestRoad() bool {
	return p.hasLongestRoad
}

func (p *player) LongestRoadLength() int {
	return p.longestRoadLength
}

func (p *player) HasLargestArmy() bool {
	return p.hasLargestArmy
}

func (p *player) ArmySize() int {
	return p.armySize
}
