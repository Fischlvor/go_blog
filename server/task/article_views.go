package task

import (
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/scriptlanguage"
	"server/global"
	"server/model/elasticsearch"
	"server/service"
	"server/utils"
	"strconv"
	"time"
)

// UpdateArticleMetricsTask combines views and scores update functionality.
func UpdateArticleMetricsTask() error {
	// Step 1: Get the views total from Redis cache
	viewsInfo := service.ArticleFieldCache.GetAllInfo("views")

	// Step 2: Query Elasticsearch for article data
	req := &search.Request{
		Query: &types.Query{
			MatchAll: &types.MatchAllQuery{},
		},
		Source_: []string{"created_at", "views", "comments", "likes"},
	}

	res, err := global.ESClient.Search().Index(elasticsearch.ArticleIndex()).Request(req).Do(context.TODO())
	if err != nil {
		return err
	}

	// Step 3: Process each article and update both views and score
	for _, hit := range res.Hits.Hits {
		var article struct {
			CreatedAt string `json:"created_at"`
			Views     int64  `json:"views"`
			Comments  int64  `json:"comments"`
			Likes     int64  `json:"likes"`
		}

		if err := json.Unmarshal(hit.Source_, &article); err != nil {
			return err
		}

		// Use the total views from the cache (if available) instead of the increment
		if cachedViews, cacheExists := viewsInfo[*hit.Id_]; cacheExists {
			// Directly use the cached views, without adding to the existing value
			article.Views = max(article.Views, int64(cachedViews))
		}

		// Step 4: Calculate the score
		createTime, err := time.Parse("2006-01-02 15:04:05", article.CreatedAt)
		if err != nil {
			return err
		}

		score := utils.CalculateArticleScore(createTime, article.Views, article.Comments, article.Likes)

		// Update Elasticsearch document with new views and score
		source := "ctx._source.views = " + strconv.Itoa(int(article.Views)) + " ; ctx._source.score = " + strconv.FormatFloat(score, 'f', -1, 64)
		script := types.Script{
			Source: &source,
			Lang:   &scriptlanguage.Painless,
		}

		if _, err := global.ESClient.Update(elasticsearch.ArticleIndex(), *hit.Id_).Script(&script).Do(context.TODO()); err != nil {
			return err
		}
	}

	// Step 5: Clear views from Redis cache after update
	service.ArticleFieldCache.Clear("views")
	service.ArticleCache.Clear()

	return nil
}
