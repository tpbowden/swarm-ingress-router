package service

import (
	"encoding/json"
	"fmt"
)

func Serialize(services []Service) string {
	serializedServices, err := json.Marshal(services)
	if err != nil {
		panic(fmt.Sprintf("Failed to serialize services: %v", err))
	}

	return string(serializedServices)
}

func Deserialize(serializedServices string) []Service {
	var services []Service

	if err := json.Unmarshal([]byte(serializedServices), &services); err != nil {
		panic(fmt.Sprintf("Failed to deserialize services: %v", err))
	}

	return services
}
