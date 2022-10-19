package localstorage

import (
	"AlexSarva/media/models"
	"errors"
	"log"
	"sync"
)

var ErrNodeNotFound = errors.New("node not found")

type NodeStorage struct {
	NodeList map[int64]*models.NodeDescription
	mutex    *sync.Mutex
}

func NewNodeLocalStorage() *NodeStorage {
	return &NodeStorage{
		NodeList: make(map[int64]*models.NodeDescription),
		mutex:    new(sync.Mutex),
	}
}

func (s *NodeStorage) Upsert(id int64, desc *models.NodeDescription) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	node, ok := s.NodeList[id]
	if !ok {
		log.Println(desc.Value)
		s.NodeList[id] = desc
	} else {
		log.Printf("New %d", desc.Value)
		newNode := &models.NodeDescription{
			Value: node.Value + desc.Value,
			Label: node.Label,
		}
		s.NodeList[id] = newNode
	}

	return nil
}

func (s *NodeStorage) Get(id int64) (*models.NodeDescription, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	node, ok := s.NodeList[id]
	if !ok {
		return &models.NodeDescription{}, ErrNodeNotFound
	}
	return node, nil
}

func (s *NodeStorage) GenerateNodes(query string) ([]models.GraphNode, error) {
	s.mutex.Lock()
	var nodes []models.GraphNode
	for key, element := range s.NodeList {

		var color = "#41e0c9"

		if element.Label == query {
			color = "#e04141"
		}

		node := models.GraphNode{
			ID:    key,
			Label: element.Label,
			Color: color,
			Value: element.Value,
		}
		nodes = append(nodes, node)
	}

	s.mutex.Unlock()

	return nodes, nil
}
