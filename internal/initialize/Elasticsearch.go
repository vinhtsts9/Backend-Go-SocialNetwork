package initialize

import (
	"log"

	"go-ecommerce-backend-api/m/v2/global"

	"github.com/elastic/go-elasticsearch/v8"
)

func InitElasticSearch() {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %s0", err)
	}
	global.Elastic = es
	log.Println("Connected to Elasticsearch")
}
