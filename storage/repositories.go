package storage

import (
	"AlexSarva/media/models"

	"github.com/google/uuid"
)

// Repo primary interface for all types of databases
type Repo interface {
	Ping() bool
	GetSearch(text string) ([]models.SearchRes, error)
	GetGraphByURL(text string) (models.Graph, error)
	GetGraphByID(id int) (models.Graph, error)
	GetFullGraph() (models.Graph, error)
	GetSourceInfoByURL(text string) (models.GraphNode, error)
	GetSourceInfoByID(id int) (models.GraphNode, error)
	AddNewGraph(graphInfo models.NewGraph) (models.NewGraphResp, error)
	GetGraphCards(userID uuid.UUID) ([]models.GraphCard, error)
	DeleteGraphCard(userID, graphID uuid.UUID) ([]models.GraphCard, error)
	GetGraphByUUID(GraphID uuid.UUID) (models.GraphExtended, error)
	//GetGraphData(url string) ([]models.DataForGraph, error)
	//NewUser(user *models.User) error
	//GetUser(username string) (*models.User, error)
}
