package citiesandknights

import "github.com/ericding98/catan"

type Commodity interface {
	catan.Resource
	isCommodity()
}

// Static type check
var _ Commodity = commodity("")

type commodity string

func (commodity) IsResource() {}

func (commodity) isCommodity() {}

const (
	PaperCommodity commodity = commodity("paper")
	CoinCommodity  commodity = commodity("coin")
	ClothCommodity commodity = commodity("cloth")
)
