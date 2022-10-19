package models

import (
	"time"

	"github.com/google/uuid"
)

type GraphQuery struct {
	Query string `json:"query"`
}

type GraphQueryID struct {
	ID int `json:"query"`
}

type DataForGraph struct {
	UrlFrom   string `json:"url_from" ch:"url_from"`
	UrlFromID int64  `json:"url_from_id" ch:"url_from_id"`
	UrlTo     string `json:"url_to" ch:"url_to"`
	UrlToID   int64  `json:"url_to_id" ch:"url_to_id"`
	Cnt       int32  `json:"cnt" ch:"cnt_links"`
}

type GraphNodeColorStyle struct {
	Background string `json:"background,omitempty"`
	Border     string `json:"border,omitempty"`
}

type GraphNodeColor struct {
	Background string              `json:"background"`
	Border     string              `json:"border,omitempty"`
	Highlight  GraphNodeColorStyle `json:"highlight,omitempty"`
	Hover      GraphNodeColorStyle `json:"hover,omitempty"`
}

type GraphNode struct {
	ID    int64       `json:"id" db:"id"`
	Label string      `json:"title" db:"url"`
	Title string      `json:"label" db:"title"`
	Color interface{} `json:"color,omitempty"`
	Value int32       `json:"value" db:"links"`
}

type NodeDescription struct {
	Value int32
	Label string
}

type GraphEdge struct {
	From int64 `json:"from" db:"id_from"`
	To   int64 `json:"to" db:"id_to"`
	//Value  int32 `json:"value" db:"links"`
	Dashes bool `json:"dashes"`
}

type Graph struct {
	Nodes []GraphNode `json:"nodes"`
	Edges []GraphEdge `json:"edges"`
}

type GraphExtended struct {
	Nodes     []GraphNode `json:"nodes"`
	Edges     []GraphEdge `json:"edges"`
	NodesList []GraphNode `json:"nodes_list"`
}

type NewGraphElement struct {
	ID  int `json:"id" db:"node"`
	Num int `json:"num" db:"num"`
}

type NewGraph struct {
	Description string            `json:"description" db:"description"`
	Sources     []NewGraphElement `json:"sources"`
	Cnt         int               `db:"cnt_elements"`
	GraphID     uuid.UUID         `json:"graph_id" db:"graph_id"`
	UserID      uuid.UUID         `json:"user_id" db:"user_id"`
}

type NewGraphResp struct {
	UUID        uuid.UUID `json:"graph_id" db:"graph_id"`
	Description string    `json:"description" db:"description"`
	Created     time.Time `json:"created" db:"created"`
}

type GraphCard struct {
	GraphID     uuid.UUID `json:"graph_id" db:"graph_id"`
	Cnt         int       `json:"cnt" db:"cnt_elements"`
	Description string    `json:"description" db:"description"`
	Created     time.Time `json:"created" db:"created"`
}

type GraphDel struct {
	GraphID uuid.UUID `json:"graph_id" db:"graph_id"`
}

type GraphUUID struct {
	GraphID uuid.UUID `json:"graph_id" db:"graph_id"`
}
