package catan

type Resource interface {
	IsResource()
}

// Static type check
var _ Resource = resource("")

type resource string

func (resource) IsResource() {}

const (
	WoodResource  resource = resource("wood")
	BrickResource resource = resource("brick")
	WoolResource  resource = resource("wool")
	OreResource   resource = resource("ore")
	GrainResource resource = resource("grain")
)
