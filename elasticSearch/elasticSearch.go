package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"go-ecommerce-backend-api/m/v2/global"
	model "go-ecommerce-backend-api/m/v2/internal/models"
	"log"

	"github.com/gin-gonic/gin"
)

func SearchUser(ctx *gin.Context) {
	query := ctx.DefaultQuery("name", "")
	if query == "" {
		ctx.JSON(400, gin.H{"error": "Name parameter is required"})
		return
	}

	// Tạo query Elasticsearch
	esQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []interface{}{
					map[string]interface{}{
						"wildcard": map[string]interface{}{
							"after.user_nickname": map[string]interface{}{
								"value":            "*" + query + "*",
								"case_insensitive": true,
							},
						},
					},
					map[string]interface{}{
						"match_phrase": map[string]interface{}{
							"after.user_nickname": map[string]interface{}{
								"query":    query,
								"slop":     1,
								"analyzer": "standard",
							},
						},
					},
					// Nếu muốn dùng query_string cho wildcard phức tạp hơn, có thể thêm:
					// map[string]interface{}{
					// 	"query_string": map[string]interface{}{
					// 		"default_field": "after.user_nickname",
					// 		"query":         "*" + query + "*",
					// 		"analyze_wildcard": true,
					// 	},
					// },
				},
				"minimum_should_match": 1,
			},
		},
		"size":    5,
		"_source": []string{"after.user_nickname", "after.user_avatar"},
	}

	queryJSON, err := json.Marshal(esQuery)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to build query"})
		return
	}

	// Gửi request tìm kiếm tới Elasticsearch
	res, err := global.Elastic.Search(
		global.Elastic.Search.WithContext(context.Background()),
		global.Elastic.Search.WithIndex("dbserver1.test.user_info"),
		global.Elastic.Search.WithBody(bytes.NewReader(queryJSON)),
		global.Elastic.Search.WithPretty(),
	)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Search failed"})
		log.Println("Elasticsearch search error:", err)
		return
	}
	defer res.Body.Close()

	// Decode phản hồi từ Elasticsearch
	var esResponse map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&esResponse); err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to parse response"})
		log.Println("Elasticsearch response parsing error:", err)
		return
	}

	// Khởi tạo danh sách người dùng
	users := []model.UserSearch{}
	hits, ok := esResponse["hits"].(map[string]interface{})
	if !ok || hits == nil {
		ctx.JSON(500, gin.H{"error": "Invalid or empty hits in response"})
		log.Println("Invalid or empty hits in response")
		return
	}
	// Lấy danh sách hits từ phản hồi
	hitsData, ok := esResponse["hits"].(map[string]interface{})["hits"].([]interface{})
	if !ok {
		ctx.JSON(500, gin.H{"error": "Unexpected response structure"})
		log.Println("Invalid response structure for hits")
		return
	}

	// Duyệt qua các hits để xây dựng danh sách người dùng
	for _, hit := range hitsData {
		hitMap, ok := hit.(map[string]interface{})
		if !ok {
			log.Println("Invalid hit structure")
			continue
		}

		source, ok := hitMap["_source"].(map[string]interface{})
		if !ok {
			log.Println("Missing _source in hit")
			continue
		}

		after, ok := source["after"].(map[string]interface{})
		if !ok {
			log.Println("Missing after in _source")
			continue
		}

		// Lấy thông tin từ "after"
		user := model.UserSearch{
			UserNickname: getString(after, "user_nickname"),
			UserAvatar:   getString(after, "user_avatar"),
		}
		users = append(users, user)
	}

	// Ghi log phản hồi
	log.Printf("Elasticsearch response: %+v\n", esResponse)

	// Trả kết quả
	ctx.JSON(200, users)
}

// Hàm tiện ích để lấy chuỗi an toàn từ map
func getString(data map[string]interface{}, key string) string {
	if val, ok := data[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}
