package ipresolver

import (
	"usepolymer.co/infrastructure/ipresolver/maxmind"
	"usepolymer.co/infrastructure/ipresolver/types"
)

var IPResolverInstance types.IPResolver = &maxmind.MaxMindIPResolver{}
