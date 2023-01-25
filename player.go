package catan

type UnimplementedPlayer interface {
	ID() string
	Board() Board
}

type unimplementedPlayer struct {
	id string
	b  Board
}

func (p *unimplementedPlayer) ID() string {
	return p.id
}

func (p *unimplementedPlayer) Board() Board {
	return p.b
}

type Strategy interface {
	PlaceInitialSettlementOn(Board) *Settlement
	PlaceRoadAnywhereNear(Piece) *Road
}

type Player interface {
	UnimplementedPlayer
	Strategy
	Points() int
	Reset(Strategy, Board)
	NewSettlement() *Settlement
	NewCity() *City
	NewRoad() *Road
	HasLongestRoad() bool
	LongestRoadLength() int
	HasLargestArmy() bool
	ArmySize() int
}

type player struct {
	Strategy
	UnimplementedPlayer
	vps               int
	hasLongestRoad    bool
	longestRoadLength int
	hasLargestArmy    bool
	armySize          int
}

func NewPlayer(p UnimplementedPlayer, s Strategy) Player {
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

	p.UnimplementedPlayer.Board().IterV(func(v *Vertex) bool {
		p := v.Piece()

		if p.Owner().ID() != p.ID() {
			return true
		}

		switch p.(type) {
		case *Settlement:
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

func (p *player) Reset(s Strategy, b Board) {
	p.Strategy = s
	p.vps = 0
}

func (p *player) NewSettlement() *Settlement {
	return NewSettlement(p)
}

func (p *player) NewCity() *City {
	return NewCity(p)
}

func (p *player) NewRoad() *Road {
	return NewRoad(p)
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
