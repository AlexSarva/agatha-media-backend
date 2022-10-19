package storagepg

import (
	"AlexSarva/media/models"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// ErrDuplicatePK error that occurs when adding exists user or order number
var ErrDuplicatePK = errors.New("duplicate PK")

// ErrNoData error that occurs when no values selected from database
var ErrNoData = errors.New("no values to delete")

type PostgresDB struct {
	database *sqlx.DB
}

func NewPostgresDBConnection(config string) *PostgresDB {
	db, err := sqlx.Connect("postgres", config)
	if err != nil {
		log.Fatalln(err)
	}
	return &PostgresDB{
		database: db,
	}
}

func (d *PostgresDB) Ping() bool {
	return d.database.Ping() == nil
}

//func (d *PostgresDB) InsertURL(id, rawURL, shortURL, userID string) error {
//	URLData := &models.URL{
//		ID:       id,
//		RawURL:   rawURL,
//		ShortURL: shortURL,
//		Created:  time.Now(),
//		UserID:   userID,
//	}
//
//	tx := d.database.MustBegin()
//	resInsert, resErr := tx.NamedExec("INSERT INTO public.urls (id, short_url, raw_url, user_id, created) VALUES (:id, :short_url, :raw_url, :user_id, :created) on conflict (raw_url) do nothing ", &URLData)
//	affectedRows, _ := resInsert.RowsAffected()
//	if affectedRows == 0 {
//		return storage.ErrDuplicatePK
//	}
//	if resErr != nil {
//		log.Println(resErr)
//	}
//	commitErr := tx.Commit()
//	if commitErr != nil {
//		log.Println(commitErr)
//	}
//	return nil
//}

//func (d *PostgresDB) InsertMany(bathURL []models.URL) error {
//	_, err := d.database.NamedExec(`INSERT INTO public.urls (id, short_url, raw_url, user_id, created)
//        VALUES (:id, :short_url, :raw_url, :user_id, :created) on conflict (raw_url) do nothing`, bathURL)
//	if err != nil {
//		log.Println(err)
//	}
//	return nil
//}

func (d *PostgresDB) GetSearch(text string) ([]models.SearchRes, error) {
	var srcs []models.SearchRes
	log.Println(text)
	query := "select id, url, coalesce(title, url) title from analytics.graph_nodes where search_field ilike '%" + text + "%' order by id limit 5;"
	err := d.database.Select(&srcs, query)
	if err != nil {
		log.Println(err)
	}
	return srcs, err
}

func (d *PostgresDB) GetGraphByURL(text string) (models.Graph, error) {
	var graph models.Graph
	var mainNode models.GraphNode
	var graphSubNodes []models.GraphNode
	var graphNodes []models.GraphNode
	var graphRawEdges []models.GraphEdge
	var graphEdges []models.GraphEdge

	errNode := d.database.Get(&mainNode, "SELECT id, url, links, coalesce(title, url) title FROM analytics.graph_nodes WHERE url=$1", text)
	if errNode != nil {
		log.Println("errNode: ", errNode)
		return models.Graph{}, errNode
	}
	mainNode.Color = models.GraphNodeColor{
		Background: "rgba(8, 217, 174, 0.9)",
		Border:     "rgba(96, 169, 191, 0.8)",
		Highlight: models.GraphNodeColorStyle{
			Background: "rgb(187, 163, 217)",
			Border:     "rgb(187, 163, 217)",
		},
		Hover: models.GraphNodeColorStyle{
			Background: "rgba(8, 217, 174, 0.9)",
			Border:     "rgb(211, 114, 214)",
		},
	}
	graphNodes = append(graphNodes, mainNode)

	errSubNodes := d.database.Select(&graphSubNodes, `
with main_id as (
    select id from analytics.graph_nodes where url = $1),
    all_nodes as (
select distinct unnest(array[id_from, id_to]) ids from analytics.graph_edges
where exists(select 1 from main_id where main_id.id = graph_edges.id_from)
and links >= 5)
select id, url, links, coalesce(title, url) title FROM analytics.graph_nodes
where 1=1
and not exists(select 1 from main_id where main_id.id = graph_nodes.id)
and exists(select 1 from all_nodes where all_nodes.ids = graph_nodes.id);`, text)
	if errSubNodes != nil {
		log.Println("errSubNode: ", errSubNodes)
		return models.Graph{}, errSubNodes
	}

	for _, node := range graphSubNodes {
		node.Color = models.GraphNodeColor{
			Background: "rgba(252, 213, 173, 0.9)",
			Border:     "rgb(252, 213, 173)",
			Highlight: models.GraphNodeColorStyle{
				Background: "rgb(187, 163, 217)",
				Border:     "rgb(187, 163, 217)",
			},
			Hover: models.GraphNodeColorStyle{
				Background: "rgba(252, 213, 173, 0.9)",
				Border:     "rgb(211, 114, 214)",
			},
		}
		graphNodes = append(graphNodes, node)
	}

	graph.Nodes = graphNodes

	errEdges := d.database.Select(&graphRawEdges, `
with main_id as (
select id from analytics.graph_nodes where url = $1)
select id_from, id_to from analytics.graph_edges
where exists(select 1 from main_id where main_id.id = graph_edges.id_from)
and links >= 5;`, text)
	if errEdges != nil {
		log.Println("errEdges: ", errEdges)
		return models.Graph{}, errEdges
	}

	for _, edge := range graphRawEdges {
		edge.Dashes = true
		graphEdges = append(graphEdges, edge)
	}

	graph.Edges = graphEdges

	return graph, nil
}

func (d *PostgresDB) GetGraphByUUID(graphID uuid.UUID) (models.GraphExtended, error) {
	var graph models.GraphExtended
	var graphRawEdges []models.GraphEdge
	var graphEdges []models.GraphEdge
	var listNodes []models.GraphNode
	var mainRawNodes []models.GraphNode
	var mainNodes []models.GraphNode
	var subRawNodes []models.GraphNode

	errEdges := d.database.Select(&graphRawEdges, `with nodes as (
select node from media.graphs_elements
            where graph_id = $1)
select id_from, id_to from analytics.graph_edges
where 1=1
and links >= 5
and exists(select 1 from nodes where nodes.node = graph_edges.id_from);`, graphID)
	if errEdges != nil {
		log.Println("errEdges: ", errEdges)
		return models.GraphExtended{}, errEdges
	}

	for _, edge := range graphRawEdges {
		edge.Dashes = true
		graphEdges = append(graphEdges, edge)
	}

	errNode := d.database.Select(&mainRawNodes, `with nodes as (
    select node, num from media.graphs_elements
    where graph_id = $1)
SELECT id, url, links, coalesce(title, url) title FROM analytics.graph_nodes
inner join nodes on nodes.node = graph_nodes.id
where 1=1
and exists(select 1 from nodes where nodes.node = graph_nodes.id)
order by nodes.num;`, graphID)
	if errNode != nil {
		log.Println("errNode: ", errNode)
		return models.GraphExtended{}, errNode
	}

	for _, node := range mainRawNodes {
		node.Color = models.GraphNodeColor{
			Background: "rgba(8, 217, 174, 0.9)",
			Border:     "rgba(96, 169, 191, 0.8)",
			Highlight: models.GraphNodeColorStyle{
				Background: "rgb(187, 163, 217)",
				Border:     "rgb(187, 163, 217)",
			},
			Hover: models.GraphNodeColorStyle{
				Background: "rgba(8, 217, 174, 0.9)",
				Border:     "rgb(211, 114, 214)",
			},
		}
		mainNodes = append(mainNodes, node)
		listNodes = append(listNodes, node)
	}

	errSubNode := d.database.Select(&subRawNodes, `
with nodes as (
    select node from media.graphs_elements
    where graph_id = $1),
    all_nodes as (
        select distinct unnest(array[id_from, id_to]) ids from analytics.graph_edges
        where 1=1
        and exists(select 1 from nodes where nodes.node = graph_edges.id_from)
        and links >= 5)
select id, url, links, coalesce(title, url) title FROM analytics.graph_nodes
where 1=1
  and not exists(select 1 from nodes where nodes.node = graph_nodes.id)
  and exists(select 1 from all_nodes where all_nodes.ids = graph_nodes.id);`, graphID)
	if errSubNode != nil {
		log.Println("errSubNode: ", errSubNode)
		return models.GraphExtended{}, errSubNode
	}

	for _, node := range subRawNodes {
		node.Color = models.GraphNodeColor{
			Background: "rgba(252, 213, 173, 0.9)",
			Border:     "rgb(252, 213, 173)",
			Highlight: models.GraphNodeColorStyle{
				Background: "rgb(187, 163, 217)",
				Border:     "rgb(187, 163, 217)",
			},
			Hover: models.GraphNodeColorStyle{
				Background: "rgba(252, 213, 173, 0.9)",
				Border:     "rgb(211, 114, 214)",
			},
		}
		mainNodes = append(mainNodes, node)
	}

	graph.Edges = graphEdges
	graph.Nodes = mainNodes
	graph.NodesList = listNodes

	return graph, nil
}

func (d *PostgresDB) GetGraphByID(id int) (models.Graph, error) {
	var graph models.Graph
	var mainNode models.GraphNode
	var graphSubNodes []models.GraphNode
	var graphNodes []models.GraphNode
	var graphRawEdges []models.GraphEdge
	var graphEdges []models.GraphEdge

	errEdges := d.database.Select(&graphRawEdges, `
select id_from, id_to from analytics.graph_edges
where id_from = $1 and links >= 5;`, id)
	if errEdges != nil {
		log.Println("errEdges: ", errEdges)
		return models.Graph{}, errEdges
	}

	for _, edge := range graphRawEdges {
		edge.Dashes = true
		graphEdges = append(graphEdges, edge)
	}

	if len(graphRawEdges) == 0 {
		graphEdges = append(graphEdges, models.GraphEdge{
			From: int64(id),
			To:   int64(id),
			//Value:  0,
			Dashes: false,
		})
	}

	errNode := d.database.Get(&mainNode, "SELECT id, url, links, coalesce(title, url) title FROM analytics.graph_nodes WHERE id=$1", id)
	if errNode != nil {
		log.Println("errNode: ", errNode)
		return models.Graph{}, errNode
	}
	if len(graphRawEdges) == 0 {
		mainNode.Color = models.GraphNodeColor{
			Background: "rgba(155, 168, 171, 0.9)",
			Border:     "rgba(155, 168, 171, 0.9)",
			Highlight: models.GraphNodeColorStyle{
				Background: "rgba(155, 168, 171, 0.9)",
				Border:     "rgba(155, 168, 171, 0.9)",
			},
			Hover: models.GraphNodeColorStyle{
				Background: "rgba(8, 217, 174, 0.9)",
				Border:     "rgb(211, 114, 214)",
			},
		}
	} else {
		mainNode.Color = models.GraphNodeColor{
			Background: "rgba(8, 217, 174, 0.9)",
			Border:     "rgba(96, 169, 191, 0.8)",
			Highlight: models.GraphNodeColorStyle{
				Background: "rgb(187, 163, 217)",
				Border:     "rgb(187, 163, 217)",
			},
			Hover: models.GraphNodeColorStyle{
				Background: "rgba(8, 217, 174, 0.9)",
				Border:     "rgb(211, 114, 214)",
			},
		}
	}

	graphNodes = append(graphNodes, mainNode)

	errSubNodes := d.database.Select(&graphSubNodes, `
with 
    all_nodes as (
select distinct unnest(array[id_from, id_to]) ids from analytics.graph_edges
where graph_edges.id_from = $1 and links >= 5)
select id, url, links, coalesce(title, url) title FROM analytics.graph_nodes
where 1=1
and id != $1
and exists(select 1 from all_nodes where all_nodes.ids = graph_nodes.id);`, id)
	if errSubNodes != nil {
		log.Println("errSubNode: ", errSubNodes)
		return models.Graph{}, errSubNodes
	}

	for _, node := range graphSubNodes {
		node.Color = models.GraphNodeColor{
			Background: "rgba(252, 213, 173, 0.9)",
			Border:     "rgb(252, 213, 173)",
			Highlight: models.GraphNodeColorStyle{
				Background: "rgb(187, 163, 217)",
				Border:     "rgb(187, 163, 217)",
			},
			Hover: models.GraphNodeColorStyle{
				Background: "rgba(252, 213, 173, 0.9)",
				Border:     "rgb(211, 114, 214)",
			},
		}
		graphNodes = append(graphNodes, node)
	}

	graph.Nodes = graphNodes
	graph.Edges = graphEdges

	return graph, nil
}

func (d *PostgresDB) GetSourceInfoByURL(text string) (models.GraphNode, error) {
	var srcs models.GraphNode
	errNode := d.database.Get(&srcs, "SELECT id, url, links, coalesce(title, url) title FROM analytics.graph_nodes WHERE url=$1", text)
	if errNode != nil {
		log.Println("errNode: ", errNode)
		return models.GraphNode{}, errNode
	}
	return srcs, nil
}

func (d *PostgresDB) GetSourceInfoByID(id int) (models.GraphNode, error) {
	var srcs models.GraphNode
	errNode := d.database.Get(&srcs, "SELECT id, url, links, coalesce(title, url) title FROM analytics.graph_nodes WHERE id=$1", id)
	if errNode != nil {
		log.Println("errNode: ", errNode)
		return models.GraphNode{}, errNode
	}
	return srcs, nil
}

func (d *PostgresDB) AddNewGraph(graphInfo models.NewGraph) (models.NewGraphResp, error) {
	log.Println("Работаем с базой")
	tx := d.database.MustBegin()
	resInsert, resErr := tx.NamedExec("INSERT INTO media.graphs (user_id, graph_id, cnt_elements, description) VALUES (:user_id, :graph_id, :cnt_elements, :description) on conflict(graph_id) do nothing", &graphInfo)
	if resErr != nil {
		tx.Commit()
		return models.NewGraphResp{}, resErr
	}

	affectedMainRows, _ := resInsert.RowsAffected()
	if affectedMainRows == 0 {
		tx.Commit()
		return models.NewGraphResp{}, ErrDuplicatePK
	}

	affectedRows, _ := resInsert.RowsAffected()
	log.Println("Загружено строк в media.graphs: ", affectedRows)

	nodesRows := 0
	for _, item := range graphInfo.Sources {
		res := tx.MustExec("INSERT INTO media.graphs_elements (graph_id, node, num) VALUES ($1, $2, $3)", graphInfo.GraphID, item.ID, item.Num)
		rows, _ := res.RowsAffected()
		nodesRows += int(rows)
	}
	log.Println("Загружено строк в media.graphs_elements: ", nodesRows)

	commitErr := tx.Commit()
	if commitErr != nil {
		return models.NewGraphResp{}, commitErr
	}

	var srcs models.NewGraphResp
	errCreated := d.database.Get(&srcs, "SELECT graph_id, description, created FROM media.graphs WHERE graph_id=$1", graphInfo.GraphID)
	if errCreated != nil {
		return models.NewGraphResp{}, errCreated
	}
	return srcs, nil
}

func (d *PostgresDB) GetGraphCards(userID uuid.UUID) ([]models.GraphCard, error) {
	var graphCards []models.GraphCard
	graphCardsErr := d.database.Select(&graphCards, `select graph_id, cnt_elements, description, created from media.graphs
where 1=1
and is_del = 0
and user_id = $1
order by created desc;
`, userID)
	if graphCardsErr != nil {
		log.Println("Нет загруженных графов ", graphCardsErr)
		return []models.GraphCard{}, graphCardsErr
	}

	return graphCards, nil
}

func (d *PostgresDB) DeleteGraphCard(userID, graphID uuid.UUID) ([]models.GraphCard, error) {
	sqlStr := fmt.Sprintf("update media.graphs set is_del=1 where user_id = '%s' and graph_id = '%s'", userID.String(), graphID.String())
	log.Println(sqlStr)
	ret, err := d.database.Exec(sqlStr)
	if err != nil {
		log.Printf("update failed, err:%v\n", err)
		return []models.GraphCard{}, err
	}
	affectedMainRows, _ := ret.RowsAffected()
	if affectedMainRows == 0 {
		return []models.GraphCard{}, ErrNoData
	}
	log.Printf("update success, affected rows:%d\n", affectedMainRows)

	var graphCards []models.GraphCard
	graphCardsErr := d.database.Select(&graphCards, `select graph_id, cnt_elements, description, created from media.graphs
where 1=1
and is_del = 0
and user_id = $1
order by created desc;
`, userID)
	if graphCardsErr != nil {
		log.Println("Нет загруженных графов ", graphCardsErr)
		return []models.GraphCard{}, graphCardsErr
	}

	return graphCards, nil
}

func (d *PostgresDB) GetFullGraph() (models.Graph, error) {
	var graph models.Graph
	//var mainNode models.GraphNode
	var graphSubNodes []models.GraphNode
	var graphNodes []models.GraphNode
	var graphRawEdges []models.GraphEdge
	var graphEdges []models.GraphEdge

	errEdges := d.database.Select(&graphRawEdges, `
select id_from, id_to from analytics.graph_edges
where links >= 5;`)
	if errEdges != nil {
		log.Println("errEdges: ", errEdges)
		return models.Graph{}, errEdges
	}

	for _, edge := range graphRawEdges {
		edge.Dashes = true
		graphEdges = append(graphEdges, edge)
	}

	errSubNodes := d.database.Select(&graphSubNodes, `
with 
    all_nodes as (
select distinct unnest(array[id_from, id_to]) ids from analytics.graph_edges
where links >= 5)
select id, url, links, coalesce(title, url) title FROM analytics.graph_nodes
where 1=1
and exists(select 1 from all_nodes where all_nodes.ids = graph_nodes.id);`)
	if errSubNodes != nil {
		log.Println("errSubNode: ", errSubNodes)
		return models.Graph{}, errSubNodes
	}

	for _, node := range graphSubNodes {
		node.Color = models.GraphNodeColor{
			Background: "rgba(252, 213, 173, 0.9)",
			Border:     "rgb(252, 213, 173)",
			Highlight: models.GraphNodeColorStyle{
				Background: "rgb(187, 163, 217)",
				Border:     "rgb(187, 163, 217)",
			},
			Hover: models.GraphNodeColorStyle{
				Background: "rgba(252, 213, 173, 0.9)",
				Border:     "rgb(211, 114, 214)",
			},
		}
		graphNodes = append(graphNodes, node)
	}

	graph.Nodes = graphNodes
	graph.Edges = graphEdges

	return graph, nil
}

//func (d *PostgresDB) GetURLByRaw(rawURL string) (*models.URL, error) {
//	var getURL models.URL
//	err := d.database.Get(&getURL, "SELECT id, short_url, raw_url, user_id, deleted, created FROM public.urls WHERE raw_url=$1", rawURL)
//	if err != nil {
//		log.Println(err)
//	}
//	return &getURL, err
//}
//
//func (d *PostgresDB) GetUserURLs(userID string) ([]models.UserURL, error) {
//	var allURLs []models.UserURL
//	log.Println(userID)
//	err := d.database.Select(&allURLs, "SELECT short_url, raw_url FROM public.urls where user_id=$1", userID)
//	if err != nil {
//		log.Println(err)
//	}
//	return allURLs, err
//}
//
//func (d *PostgresDB) Delete(userID string, shortURLs []string) error {
//	query := "UPDATE public.urls SET deleted=1 WHERE user_id=? AND id IN (?)"
//	qry, args, err := sqlx.In(query, userID, shortURLs)
//	if err != nil {
//		return err
//	}
//
//	if _, execErr := d.database.Exec(d.database.Rebind(qry), args...); execErr != nil {
//		log.Println(execErr)
//		return execErr
//	}
//	return nil
//}
