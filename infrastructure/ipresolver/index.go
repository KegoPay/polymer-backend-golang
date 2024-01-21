package ipresolver

import (
	"kego.com/infrastructure/ipresolver/maxmind"
	"kego.com/infrastructure/ipresolver/types"
)

var IPResolverInstance types.IPResolver =  &maxmind.MaxMindIPResolver{}