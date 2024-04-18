package ktwingraph

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin"
	log "github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/logger"
)

var logger = log.NewLogger()

func LoadTwinGraphByInstances(twinInstances []string) (ktwin.TwinGraph, error) {
	var ktwinGraph ktwin.TwinGraph
	var ktwinGraphList []ktwin.TwinInstanceGraph
	for _, twinInstance := range twinInstances {
		ktwinGraph, err := getTwinGraphInstance(twinInstance)
		if err != nil {
			logger.Error("Error getting Twin Graph instance", err)
		}
		ktwinGraphList = append(ktwinGraphList, ktwinGraph.TwinInstancesGraph...)
	}

	if len(ktwinGraphList) == 0 {
		writeTwinGraph(ktwinGraph)
		return ktwinGraph, nil
	}

	ktwinGraph.TwinInstancesGraph = ktwinGraphList

	if os.Getenv("ENV") != "local" && os.Getenv("ENV") != "test" {
		writeTwinGraph(ktwinGraph)
	}
	return ktwinGraph, nil
}

func getTwinGraphInstance(twinInstance string) (*ktwin.TwinGraph, error) {
	var ktwinGraph ktwin.TwinGraph

	if os.Getenv("ENV") == "local" || os.Getenv("ENV") == "test" {
		ktwinGraph, err := loadLocalTwinGraph()
		return ktwinGraph, err
	}

	ktwinGraphStoreURL := os.Getenv("KTWIN_GRAPH_URL")
	response, err := http.Get(ktwinGraphStoreURL + "/" + twinInstance)
	if err != nil {
		fmt.Println("Error while calling service:", err)
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		message := fmt.Sprintf("Error while calling service, status code: %d", response.StatusCode)
		logger.Error(message, nil)
		return nil, errors.New(message)
	}

	if json.NewDecoder(response.Body).Decode(&ktwinGraph); err != nil {
		logger.Error("Error parsing response body:", err)
		return nil, err
	}

	return &ktwinGraph, nil
}

func loadLocalTwinGraph() (*ktwin.TwinGraph, error) {
	jsonStr := os.Getenv("KTWIN_GRAPH")
	var result ktwin.TwinGraph
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		logger.Error("Error parsing JSON:", err)
		return nil, err
	}
	return &result, nil
}

func writeTwinGraph(ktwinGraph ktwin.TwinGraph) error {
	graphByteArray, err := json.Marshal(ktwinGraph)

	if err != nil {
		return err
	}

	err = os.WriteFile("ktwin_graph.json", graphByteArray, 0644)

	if err != nil {
		return err
	}

	return nil
}

// Get the Graph relationship by name and instance
func GetRelationshipFromGraph(twinInstance, relationshipName string, twinGraph ktwin.TwinGraph) *ktwin.TwinInstanceReference {
	for _, instance := range twinGraph.TwinInstancesGraph {
		if instance.Name == twinInstance {
			for _, relationship := range instance.Relationships {
				if relationship.Name == relationshipName {
					return &relationship
				}
			}
		}
	}

	return nil
}

// Get Twin Graph Node by twin instance and interface
func GetTwinGraphByRelation(targetTwinInterface, sourceTwinInstance string, twinGraph ktwin.TwinGraph) *ktwin.TwinInstanceReference {
	for _, sourceTwinGraph := range twinGraph.TwinInstancesGraph {
		for _, relationship := range sourceTwinGraph.Relationships {
			if relationship.Interface == targetTwinInterface {
				return &relationship
			}
		}
	}

	return nil
}
