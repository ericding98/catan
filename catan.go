package catan

import (
	"crypto/rand"
	Hexpkg "encoding/Hex"
)

func UID() string {
	buf := make([]byte, 8)
	rand.Read(buf)
	return Hexpkg.EncodeToString(buf)
}

type Game interface {
	GameState
	NewBoard() Board
	Populate(...Player)
	MaxPoints() int
	Done() bool
	Setup() error
	Step() error
}

type game struct {
	GameState
	maxPoints int
}

func NewGame(state GameState, maxPoints int) Game {
	return &game{
		GameState: state,
		maxPoints: maxPoints,
	}
}

func (g *game) NewBoard() Board {
	grid := make([][]*Hex, 5)
	cols := 3
	dir := 1
	for i := range grid {
		grid[i] = make([]*Hex, cols)
		for j := range grid[i] {
			grid[i][j] = new(Hex)
		}

		if cols == 5 {
			dir = -1
		}

		cols += dir
	}

	b := &board{
		g:    g,
		grid: grid,
	}

	b.joinNodes()

	g.GameState.Reset(b)

	return b
}

func (g *game) Populate(ps ...Player) {
	g.GameState.Populate(ps...)
}

func (g *game) MaxPoints() int {
	return g.maxPoints
}

func (g *game) Done() bool {
	for _, ps := range g.GameState.Players() {
		if ps.Points() == g.MaxPoints() {
			return true
		}
	}
	return false
}

type Node interface {
	Hexes() []*Hex
	Piece() Piece
	SetPiece(Piece)
	isNode()
}

type Vertex struct {
	p  Piece
	hs []*Hex
}

func (v *Vertex) Hexes() []*Hex {
	return v.hs
}

func (v *Vertex) Piece() Piece {
	return v.p
}

func (v *Vertex) SetPiece(p Piece) {
	v.p = p
}

func (*Vertex) isNode() {}

type Edge struct {
	p  Piece
	hs []*Hex
}

func (e *Edge) Hexes() []*Hex {
	return e.hs
}

func (e *Edge) Piece() Piece {
	return e.p
}

func (e *Edge) SetPiece(p Piece) {
	e.p = p
}

func (*Edge) isNode() {}

// vs[0] is the top-left-most vertex, vs[1] is the next clockwise vertex, etc.
// es[0] is the edge connecting v1 to vs[1], es[1] is the next clockwise edge, etc.
type Hex struct {
	vs [6]Vertex
	es [6]Edge
}

// Position on a game board is encoded as 4 int:
// 0 - row index
// 1 - column index
// 2 - edge or vertex
// 3 - edge or vertex index
type Position [4]uint8

type Board interface {
	Game() Game
	NewPlayer() UnimplementedPlayer
	Setup() error
	Contains(Piece) bool
	At(Position) Node
	PlacePieceAt(Piece, Position)
	IterV(func(*Vertex) bool)
	IterE(func(*Edge) bool)
	IterAll(func(Node) bool)
}

type board struct {
	g    Game
	grid [][]*Hex
}

func (b *board) joinNodes() {
	for i, row := range b.grid {
		for j, col := range row {
			if i == 0 && j == 0 {
				right := row[j+1]
				bLeft := b.grid[i+1][j]
				bRight := b.grid[i+1][j+1]
				col.vs[2].hs = append(col.vs[2].hs, right)
				col.vs[3].hs = append(col.vs[3].hs, right)
				col.vs[3].hs = append(col.vs[3].hs, bRight)
				col.vs[4].hs = append(col.vs[4].hs, bRight)
				col.vs[4].hs = append(col.vs[4].hs, bLeft)
				col.vs[5].hs = append(col.vs[5].hs, bLeft)
			}

			if i == 0 {
			}

			if i > 0 && i <= 2 {
				if j == 0 {
				}

				if j == i+2 {
				}

				// return h.vs[3:5]
			}

			if j == 0 {
			}

			// return h.vs[3:]
		}
	}
}

func (b *board) Game() Game {
	return b.g
}

func (b *board) NewPlayer() UnimplementedPlayer {
	return &unimplementedPlayer{
		id: UID(),
		b:  b,
	}
}

func (b *board) Setup() error {
	return nil
}

func (b *board) Contains(p Piece) bool {
	return false
}

func (b *board) At(pos Position) Node {
	ref := b.grid[pos[0]][pos[1]]

	if pos[2] == 0 {
		// vertex
		return &ref.vs[pos[3]]
	} else {
		// edge
		return &ref.es[pos[3]]
	}
}

func (b *board) PlacePieceAt(p Piece, pos Position) {
	ref := b.grid[pos[0]][pos[1]]

	if pos[2] == 0 {
		// vertex
		ref.vs[pos[3]].SetPiece(p)
	} else {
		// edge
		ref.es[pos[3]].SetPiece(p)
	}
}

func (b *board) collectVInspectionSet(h *Hex, i, j int) []Vertex {
	if i == 0 && j == 0 {
		return h.vs[:]
	}

	if i == 0 {
		return h.vs[1:5]
	}

	if i > 0 && i <= 2 {
		if j == 0 {
			return append(make([]Vertex, 0), h.vs[3], h.vs[4], h.vs[5], h.vs[0])
		}

		if j == i+2 {
			return h.vs[2:5]
		}

		return h.vs[3:5]
	}

	if j == 0 {
		return h.vs[3:]
	}

	return h.vs[3:]
}

func (b *board) collectEInspectionSet(h *Hex, i, j int) []Edge {
	if i == 0 && j == 0 {
		return h.es[:]
	}

	if i == 0 {
		return h.es[:5]
	}

	if i > 0 && i <= 2 {
		if j == 0 {
			return append(make([]Edge, 0), h.es[0], h.es[2], h.es[3], h.es[4], h.es[5])
		}

		if j == i+2 {
			return h.es[1:5]
		}

		return h.es[2:5]
	}

	if j == 0 {
		return h.es[2:]
	}

	return h.es[2:5]
}

func (b *board) collectAllInspectionSet(h *Hex, i, j int) []Node {
	es, vs := b.collectEInspectionSet(h, i, j), b.collectVInspectionSet(h, i, j)
	ns := make([]Node, len(es)+len(vs))

	isNextEdge := true
	ej := 0
	vj := 0
	for i := range ns {
		if isNextEdge {
			ns[i] = &es[ej]
			ej++
		} else {
			ns[i] = &vs[vj]
			vj++
		}
		isNextEdge = !isNextEdge
	}

	return ns
}

func (b *board) IterV(fn func(*Vertex) bool) {
	for i, row := range b.grid {
		for j, col := range row {
			vs := b.collectVInspectionSet(col, i, j)
			for _, v := range vs {
				if !fn(&v) {
					return
				}
			}
		}
	}
}

func (b *board) IterE(fn func(*Edge) bool) {
	for i, row := range b.grid {
		for j, col := range row {
			es := b.collectEInspectionSet(col, i, j)
			for _, e := range es {
				if !fn(&e) {
					return
				}
			}
		}
	}
}

func (b *board) IterAll(fn func(Node) bool) {
	for i, row := range b.grid {
		for j, col := range row {
			ns := b.collectAllInspectionSet(col, i, j)
			for _, n := range ns {
				if !fn(n) {
					return
				}
			}
		}
	}
}

type UnimplementedGameState interface {
	Board() Board
	Reset(b Board)
	Populate(...Player)
	Players() []Player
	PopulateInOrder(...Player)
}

type GameState interface {
	UnimplementedGameState
	Setup() error
	Step() error
}

type gameState struct {
	b     Board
	ps    []Player
	order []Player
}

func NewGameState() UnimplementedGameState {
	return &gameState{
		b:     nil,
		ps:    make([]Player, 0),
		order: make([]Player, 0),
	}
}

func (s *gameState) Board() Board {
	return s.b
}

func (s *gameState) Reset(b Board) {
	if b == nil {
		b = s.b.Game().NewBoard()
	}

	for _, p := range s.ps {
		p.Reset(p, b)
	}
}

func (s *gameState) Populate(ps ...Player) {
	if s.ps == nil {
		s.ps = make([]Player, 0)
	}

	s.ps = append(s.ps, ps...)
}

func (s *gameState) PopulateInOrder(ps ...Player) {
	if s.order == nil {
		s.order = make([]Player, 0)
	}

	s.order = append(s.ps, ps...)
}

func (s *gameState) Players() []Player {
	if len(s.order) == len(s.ps) {
		return s.order
	}
	return s.ps
}

type Piece interface {
	ID() string
	Owner() Player
	Pos() Position
	PlaceAt(Position)
	IsPlacedOnBoard() bool
	IsPiece()
}

type House interface {
	Piece
	ProduceSurroundingResources() ([]Resource, error)
	isHouse()
}

var (
	_ House = new(Settlement)
	_ House = new(City)
)

type Settlement struct {
	id    string
	owner Player
	b     Board
	pos   Position
}

func NewSettlement(p Player) *Settlement {
	return &Settlement{
		id:    UID(),
		owner: p,
		b:     p.Board(),
	}
}

func (*Settlement) isHouse() {}
func (*Settlement) IsPiece() {}

func (s *Settlement) ID() string {
	return s.id
}

func (s *Settlement) Owner() Player {
	return s.owner
}

func (s *Settlement) Pos() Position {
	return s.pos
}

func (s *Settlement) PlaceAt(pos Position) {
	s.pos = pos
	s.b.PlacePieceAt(s, pos)
}

func (s *Settlement) IsPlacedOnBoard() bool {
	return s.b.Contains(s)
}

func (s *Settlement) ProduceSurroundingResources() ([]Resource, error) {
	n := s.b.At(s.Pos())

	return nil, nil
}

type City struct {
	id    string
	owner Player
	b     Board
	pos   Position
}

func NewCity(p Player) *City {
	return &City{
		id:    UID(),
		owner: p,
		b:     p.Board(),
	}
}

func (*City) isHouse() {}
func (*City) IsPiece() {}

func (c *City) ID() string {
	return c.id
}

func (c *City) Owner() Player {
	return c.owner
}

func (c *City) Pos() Position {
	return c.pos
}

func (c *City) PlaceAt(pos Position) {
	c.pos = pos
	c.b.PlacePieceAt(c, pos)
}

func (c *City) IsPlacedOnBoard() bool {
	return c.b.Contains(c)
}

func (c *City) ProduceSurroundingResources() ([]Resource, error) {
	return nil, nil
}

type Road struct {
	id    string
	owner Player
	b     Board
	pos   Position
}

func NewRoad(p Player) *Road {
	return &Road{
		id:    UID(),
		owner: p,
		b:     p.Board(),
	}
}

func (*Road) IsPiece() {}

func (r *Road) ID() string {
	return r.id
}

func (r *Road) Owner() Player {
	return r.owner
}

func (r *Road) Pos() Position {
	return r.pos
}

func (r *Road) PlaceAt(pos Position) {
	r.pos = pos
	r.b.PlacePieceAt(r, pos)
}

func (r *Road) IsPlacedOnBoard() bool {
	return r.b.Contains(r)
}
