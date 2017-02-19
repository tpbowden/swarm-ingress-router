package collector

import (
	"encoding/json"
	"fmt"

	"github.com/tpbowden/swarm-ingress-router/types"
)

func serializeServices(services []types.Service) string {
	serializedServices, err := json.Marshal(services)

	if err != nil {
		panic(fmt.Sprintf("Failed to serialize services: %v", err))
	}

	return string(serializedServices)
}
