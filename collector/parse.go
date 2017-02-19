package collector

import (
	"github.com/docker/docker/api/types/swarm"

	"github.com/tpbowden/swarm-ingress-router/types"
)

func parseServices(swarmServices []swarm.Service) []types.Service {
	result := []types.Service{}
	return result
}
