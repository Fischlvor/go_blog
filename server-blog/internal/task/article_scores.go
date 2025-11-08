package task

import (
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/scriptlanguage"
	"server/pkg/global"
	"server/internal/model/elasticsearch"
	"server/pkg/utils"
	"strconv"
	"time"
)

// 更新文章分数的定时任务
func UpdateArticleScores() error {
	req := &search.Request{
		Query: &types.Query{
			MatchAll: &types.MatchAllQuery{},
		},
		// 正确的字段名是Source_，不是SourceIncludes
		Source_: []string{"created_at", "views", "comments", "likes"},
	}

	res, err := global.ESClient.Search().Index(elasticsearch.ArticleIndex()).Request(req).Do(context.TODO())
	if err != nil {
		return err
	}

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

		createTime, err := time.Parse("2006-01-02 15:04:05", article.CreatedAt)
		if err != nil {
			return err
		}

		score := utils.CalculateArticleScore(createTime, article.Views, article.Comments, article.Likes)

		source := "ctx._source.score = " + strconv.FormatFloat(score, 'f', -1, 64)
		// 使用正确的Script构造方式
		script := types.Script{
			Source: &source,
			Lang:   &scriptlanguage.Painless,
		}

		_, err = global.ESClient.Update(elasticsearch.ArticleIndex(), *hit.Id_).Script(&script).Do(context.TODO())

		if err != nil {
			return err
		}
	}
	return nil
}
