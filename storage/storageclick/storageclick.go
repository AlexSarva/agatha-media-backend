package clickhousestorage

import (
	"AlexSarva/media/models"
	"AlexSarva/media/utils/dbutils"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type ClickHouse struct {
	Database driver.Conn
	ctx      context.Context
}

func MyClickHouseDB(path string) *ClickHouse {
	parsedCfg, parseCfgErr := dbutils.ParseConfigDB(path)
	if parseCfgErr != nil {
		log.Println("Проблема с конфиогом для ClickHouse")
	}
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"127.0.0.1:9000"},
		Auth: clickhouse.Auth{
			Database: parsedCfg.DatabaseName,
			Username: parsedCfg.User,
			Password: parsedCfg.Password,
		},
		//Debug:           true,
		DialTimeout:     time.Second,
		MaxOpenConns:    50,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour * 24,
	})
	if err != nil {
		log.Println("НЕТ подключения к ClickHouse: ", err)
	}
	return &ClickHouse{
		Database: conn,
		ctx:      context.Background(),
	}
}

func (c *ClickHouse) Ping() bool {
	pingErr := c.Database.Ping(c.ctx)
	if pingErr != nil {
		log.Println("НЕ пингуется", pingErr)
		return false
	}
	return true
}

func (c *ClickHouse) GetGraphData(url string) ([]models.DataForGraph, error) {
	log.Printf("%T", url)
	log.Printf("%v", url)
	ctx := clickhouse.Context(context.Background(), clickhouse.WithSettings(clickhouse.Settings{
		"max_block_size": 1000,
	}), clickhouse.WithProgress(func(p *clickhouse.Progress) {
		fmt.Println("progress: ", p)
	}))
	var graphInfo []models.DataForGraph
	err := c.Database.Select(ctx, &graphInfo, `
select url_from,url_from_id, url_to, url_to_id, cnt_links from crawler.graphs
    where 1=1
    and (url_from = $1 or url_to = $1);
`, url)
	if err != nil {
		log.Println(err)
	}
	//log.Printf("%+v\n", graphInfo)
	return graphInfo, nil
}
