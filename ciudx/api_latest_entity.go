/*
 * IUDX Resource Server APIs
 *
 * The Resource Server is IUDX's data store which allows publication, subscription and discovery of data. For search and discovery, it allows users to search through temporal, geo-based and attribute queries.  For publication and subscription, it allows users to use AMQP streaming protocol over TLS. It enables *Providers* of datasources to publish data as per the IUDX data descriptor. It enables *Consumers* of datasources to search and query for data using HTTPs APIs. It enables *Subscribers* a.k.a [Streaming Consumer] of datasources to stream data using AMQP streaming protocol over TLS.
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package ciudx

import (
	"context"
	"encoding/json"
	"net/http"

	models "github.com/dataspace-mobility/rs-iudx/ciudx/models"

	"github.com/gin-gonic/gin"
	redis "github.com/go-redis/redis/v8"
)

// LatestEntities - Latest Data
func LatestEntities(c *gin.Context) {

	ctx := context.Background()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	val, err := redisClient.Get(ctx, c.Param("id")).Result()
	var res models.ModelEntityresponse
	if err == redis.Nil {
		res = models.ModelEntityresponse{
			Type:  "404",
			Title: "not found",
		}
	} else if err != nil {
		res = models.ModelEntityresponse{
			Type:  "400",
			Title: "Bad Request",
		}
	} else {
		var jsonResponse map[string]interface{}
		err := json.Unmarshal([]byte(val), &jsonResponse)
		if err != nil {
			res = models.ModelEntityresponse{
				Type:  "400",
				Title: "Bad Request",
			}
		} else {
			res = models.ModelEntityresponse{
				Type:    "200",
				Title:   "success",
				Results: []map[string]interface{}{jsonResponse},
			}
		}
	}

	c.JSON(http.StatusOK, res)
}
