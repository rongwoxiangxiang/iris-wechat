package elasticalSearch

import (
	"context"
	"encoding/json"
	"fmt"
	"gopkg.in/olivere/elastic.v5"
	"iris/common"
	"iris/config"
	"iris/models"
	"log"
	"strconv"
)
var client *elastic.Client

func init()  {
	var err error
	host := config.GetConfigs().OthersConfig.ElasticSearchServerConfig.Host
	if host == "" {
		log.Fatalf("es init fail [0]: config host never gotten")
	}
	client, err = elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(host))
	if err != nil {
		log.Fatalf("es init fail [1]: %v", err)
	}
	exists, err := client.IndexExists("replies").Do(context.Background())
	if err != nil {
		log.Fatalf("es init fail [2]: %v", err)
	}
	if !exists {
		createIndex()
	}
}

func UpdateData(reply models.ReplyModel) {
	_, err := client.Index().Index("replies").Type("reply").Id(strconv.FormatInt(reply.Id, 10)).BodyJson(reply).Do(context.Background())
	if err != nil {
		log.Fatalf("es add doc fail:reply:%v err:[1]: %v", reply, err)
	}
}

func Match(alias string, wid int64) (reply models.ReplyModel) {
	var res *elastic.SearchResult
	var err error

	boolQuery := elastic.NewBoolQuery()
	boolQuery.Must(elastic.NewTermQuery("Disabled", common.NO_VALUE))
	boolQuery.Must(elastic.NewTermQuery("Wid", wid))
	boolQuery.Should(elastic.NewTermQuery("ClickKey", alias))
	boolQuery.Should(elastic.NewMatchQuery("Alias", alias))
	res, err = client.Search("replies").Type("reply").Query(boolQuery).Size(1).Pretty(true).Do(context.Background())
	if err != nil {
		log.Println(err.Error())
	}

	if res.Hits.TotalHits > 0 {
		err := json.Unmarshal(*res.Hits.Hits[0].Source, &reply)
		if err != nil {
			log.Println("es hit json to struct fail: res:%v, err:%v", res, err)
		}
	}
	return
}

func createIndex() bool {
	response, err := client.CreateIndex("replies").Body(GetMapping()).Do(context.Background())
	if err != nil {
		log.Fatalf("es create index fail [1]: %v", err)
	}
	return response.Acknowledged
}

func DeleteIndex(index... string) bool {
	response, err := client.DeleteIndex(index...).Do(context.Background())
	if err != nil {
		fmt.Printf("es delete index failed, err: %v\n", err)
	}
	return response.Acknowledged
}