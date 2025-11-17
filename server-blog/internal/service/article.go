package service

import (
	"context"
	"encoding/json"
	"errors"
	"server/internal/model/appTypes"
	"server/internal/model/database"
	"server/internal/model/elasticsearch"
	"server/internal/model/other"
	"server/internal/model/request"
	"server/pkg/global"
	"server/pkg/utils"
	"strconv"
	"time"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/scriptlanguage"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"
	"gorm.io/gorm"
)

type ArticleService struct {
}

// processCoverURLs 处理ES查询结果中的封面URL，将相对路径转换为完整URL
func (articleService *ArticleService) processCoverURLs(hits []types.Hit) {
	for i := range hits {
		var source map[string]interface{}
		if err := json.Unmarshal(hits[i].Source_, &source); err == nil {
			if cover, exists := source["cover"].(string); exists {
				source["cover"] = utils.PublicURLFromDB(cover)
				// 重新序列化回去
				if newSource, err := json.Marshal(source); err == nil {
					hits[i].Source_ = newSource
				}
			}
		}
	}
}

func (articleService *ArticleService) ArticleInfoByID(id string) (elasticsearch.Article, error) {
	article, err := articleService.Get(id)
	if err != nil {
		return article, errors.New("failed to get article")
	}
	if cacheViews, ok := ArticleFieldCache.Get("views", id); ok == nil {
		article.Views = max(cacheViews, article.Views)
	} else {
		ArticleFieldCache.Set("views", id, article.Views)
	}
	//异步更新浏览量
	go func() {
		_ = ArticleFieldCache.Add("views", id)
	}()
	// 拼接封面URL
	article.Cover = utils.PublicURLFromDB(article.Cover)
	return article, nil
}

// 在 ArticleSearch 函数中添加默认按自定义分数排序的逻辑
func (articleService *ArticleService) ArticleSearch(info request.ArticleSearch) (interface{}, int64, error) {
	req := &search.Request{
		Query: &types.Query{},
	}

	boolQuery := &types.BoolQuery{}

	// 根据查询字段查询
	if info.Query != "" {
		boolQuery.Should = []types.Query{
			{Match: map[string]types.MatchQuery{"title": {Query: info.Query}}},
			{Match: map[string]types.MatchQuery{"keyword": {Query: info.Query}}},
			{Match: map[string]types.MatchQuery{"abstract": {Query: info.Query}}},
			{Match: map[string]types.MatchQuery{"content": {Query: info.Query}}},
		}
	}

	// 根据标签筛选
	if info.Tag != "" {
		boolQuery.Must = []types.Query{
			{Match: map[string]types.MatchQuery{"tags": {Query: info.Tag}}},
		}
	}

	// 根据类别筛选
	if info.Category != "" {
		boolQuery.Filter = []types.Query{
			{Term: map[string]types.TermQuery{"category": {Value: info.Category}}},
		}
	}

	// 如果有查询条件，则使用 Bool 查询，否则使用 MatchAll 查询
	if boolQuery.Should != nil || boolQuery.Must != nil || boolQuery.Filter != nil {
		req.Query.Bool = boolQuery
	} else {
		req.Query.MatchAll = &types.MatchAllQuery{}
	}

	// 设置排序字段
	var sortField string
	var order sortorder.SortOrder
	if info.Sort != "" {
		switch info.Sort {
		case "time":
			sortField = "created_at"
		case "view":
			sortField = "views"
		case "comment":
			sortField = "comments"
		case "like":
			sortField = "likes"
		default:
			sortField = "created_at"
		}

		if info.Order != "asc" {
			order = sortorder.Desc
		} else {
			order = sortorder.Asc
		}
	} else {
		// 如果没有指定排序条件，使用自定义分数排序
		sortField = "score"
		if info.Order != "asc" {
			order = sortorder.Desc
		} else {
			order = sortorder.Asc
		}
	}

	req.Sort = []types.SortCombinations{
		types.SortOptions{
			SortOptions: map[string]types.FieldSort{
				sortField: {Order: &order},
			},
		},
	}

	option := other.EsOption{
		PageInfo:       info.PageInfo,
		Index:          elasticsearch.ArticleIndex(),
		Request:        req,
		SourceIncludes: []string{"created_at", "cover", "title", "abstract", "category", "tags", "views", "comments", "likes", "score"},
	}
	list, total, err := utils.EsPagination(context.TODO(), option)
	if err != nil {
		return nil, 0, err
	}
	// 拼接封面URL
	articleService.processCoverURLs(list)
	return list, total, nil
}

func (articleService *ArticleService) ArticleCategory() ([]database.ArticleCategory, error) {
	var category []database.ArticleCategory
	if err := global.DB.Find(&category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (articleService *ArticleService) ArticleTags() ([]database.ArticleTag, error) {
	var tags []database.ArticleTag
	if err := global.DB.Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func (articleService *ArticleService) ArticleLike(req request.ArticleLike) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		var al database.ArticleLike
		var num int

		// 如果用户未收藏，则创建收藏记录
		if errors.Is(tx.Where("user_uuid = ? AND article_id = ?", req.UserUUID, req.ArticleID).First(&al).Error, gorm.ErrRecordNotFound) {
			if err := tx.Create(&database.ArticleLike{UserUUID: req.UserUUID, ArticleID: req.ArticleID}).Error; err != nil {
				return err
			}
			num = 1
		} else { // 如果用户已经收藏，则取消收藏
			if err := tx.Delete(&al).Error; err != nil {
				return err
			}
			num = -1
		}

		// 更新文章收藏数
		source := "ctx._source.likes += " + strconv.Itoa(num)
		script := types.Script{Source: &source, Lang: &scriptlanguage.Painless}
		_, err := global.ESClient.Update(elasticsearch.ArticleIndex(), req.ArticleID).Script(&script).Do(context.TODO())
		return err
	})
}

func (articleService *ArticleService) ArticleIsLike(req request.ArticleLike) (bool, error) {
	return !errors.Is(global.DB.Where("user_uuid = ? AND article_id = ?", req.UserUUID, req.ArticleID).First(&database.ArticleLike{}).Error, gorm.ErrRecordNotFound), nil
}

func (articleService *ArticleService) ArticleLikesList(info request.ArticleLikesList) (interface{}, int64, error) {
	db := global.DB.Where("user_uuid = ?", info.UserUUID)
	option := other.MySQLOption{
		PageInfo: info.PageInfo,
		Where:    db,
	}

	l, total, err := utils.MySQLPagination(&database.ArticleLike{}, option)
	if err != nil {
		return nil, 0, err
	}
	var list []struct {
		Id_     string                `json:"_id"`
		Source_ elasticsearch.Article `json:"_source"`
	}

	for _, articleLike := range l {
		article, err := articleService.Get(articleLike.ArticleID)
		if err != nil {
			return nil, 0, err
		}
		article.UpdatedAt = ""
		article.Keyword = ""
		article.Content = ""
		// 拼接封面URL
		article.Cover = utils.PublicURLFromDB(article.Cover)
		list = append(list, struct {
			Id_     string                `json:"_id"`
			Source_ elasticsearch.Article `json:"_source"`
		}{
			Id_:     articleLike.ArticleID,
			Source_: article,
		})
	}
	return list, total, nil
}

func (articleService *ArticleService) ArticleCreate(req request.ArticleCreate) error {
	b, err := articleService.Exits(req.Title)
	if err != nil {
		return err
	}
	if b {
		return errors.New("the article already exists")
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	dbCover := utils.DBURLFromPublic(req.Cover)
	articleToCreate := elasticsearch.Article{
		CreatedAt: now,
		UpdatedAt: now,
		Cover:     dbCover,
		Title:     req.Title,
		Keyword:   req.Title,
		Category:  req.Category,
		Tags:      req.Tags,
		Abstract:  req.Abstract,
		Content:   req.Content,
	}
	err = global.DB.Transaction(func(tx *gorm.DB) error {
		// 同时更新文章类别表中的数据
		if err := articleService.UpdateCategoryCount(tx, "", articleToCreate.Category); err != nil {
			return err
		}

		// 同时更新文章标签表中的数据
		if err := articleService.UpdateTagsCount(tx, []string{}, articleToCreate.Tags); err != nil {
			return err
		}

		// 同时更新图片表中的图片类别
		if err := utils.ChangeImagesCategory(tx, []string{dbCover}, appTypes.Cover); err != nil {
			return err
		}
		// 拿到 Text 中所有插图的 URL
		illustrations, err := utils.FindIllustrations(articleToCreate.Content)
		if err != nil {
			return err
		}
		for i := range illustrations {
			illustrations[i] = utils.DBURLFromPublic(illustrations[i])
		}
		if err := utils.ChangeImagesCategory(tx, illustrations, appTypes.Illustration); err != nil {
			return err
		}

		return articleService.Create(&articleToCreate)
	})

	return err
}

func (articleService *ArticleService) ArticleDelete(req request.ArticleDelete) error {
	if len(req.IDs) == 0 {
		return nil
	}
	return global.DB.Transaction(func(tx *gorm.DB) error {
		for _, id := range req.IDs {
			articleToDelete, err := articleService.Get(id)
			if err != nil {
				return err
			}

			// 同时更新文章类别表中的数据
			if err := articleService.UpdateCategoryCount(tx, articleToDelete.Category, ""); err != nil {
				return err
			}

			// 同时更新文章标签表中的数据
			if err := articleService.UpdateTagsCount(tx, articleToDelete.Tags, []string{}); err != nil {
				return err
			}

			// 同时更新图片表中的图片类别
			if err := utils.InitImagesCategory(tx, []string{articleToDelete.Cover}); err != nil {
				return err
			}
			illustrations, err := utils.FindIllustrations(articleToDelete.Content)
			if err != nil {
				return err
			}
			if err := utils.InitImagesCategory(tx, illustrations); err != nil {
				return err
			}
			// 同时删除该文章下的所有评论
			comments, err := ServiceGroupApp.CommentService.CommentInfoByArticleID(request.CommentInfoByArticleID{ArticleID: id})
			if err != nil {
				return err
			}
			for _, comment := range comments {
				if err := ServiceGroupApp.CommentService.DeleteCommentAndChildren(tx, comment.ID); err != nil {
					return err
				}
			}
		}
		return articleService.Delete(req.IDs)
	})
}

func (articleService *ArticleService) ArticleUpdate(req request.ArticleUpdate) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	articleToUpdate := struct {
		UpdatedAt string   `json:"updated_at"`
		Cover     string   `json:"cover"`
		Title     string   `json:"title"`
		Keyword   string   `json:"keyword"`
		Category  string   `json:"category"`
		Tags      []string `json:"tags"`
		Abstract  string   `json:"abstract"`
		Content   string   `json:"content"`
	}{
		UpdatedAt: now,
		Cover:     utils.DBURLFromPublic(req.Cover),
		Title:     req.Title,
		Keyword:   req.Title,
		Category:  req.Category,
		Tags:      req.Tags,
		Abstract:  req.Abstract,
		Content:   req.Content,
	}
	return global.DB.Transaction(func(tx *gorm.DB) error {
		oldArticle, err := articleService.Get(req.ID)
		if err != nil {
			return err
		}

		// 同时更新文章类别表中的数据
		if err := articleService.UpdateCategoryCount(tx, oldArticle.Category, articleToUpdate.Category); err != nil {
			return err
		}

		// 同时更新文章标签表中的数据
		if err := articleService.UpdateTagsCount(tx, oldArticle.Tags, articleToUpdate.Tags); err != nil {
			return err
		}

		// 同时更新图片表中的图片类别
		if articleToUpdate.Cover != oldArticle.Cover {
			if err := utils.InitImagesCategory(tx, []string{utils.DBURLFromPublic(oldArticle.Cover)}); err != nil {
				return err
			}
			if err := utils.ChangeImagesCategory(tx, []string{utils.DBURLFromPublic(articleToUpdate.Cover)}, appTypes.Cover); err != nil {
				return err
			}
		}
		oldIllustrations, err := utils.FindIllustrations(oldArticle.Content)
		if err != nil {
			return err
		}
		newIllustrations, err := utils.FindIllustrations(articleToUpdate.Content)
		if err != nil {
			return err
		}
		addedIllustrations, removedIllustrations := utils.DiffArrays(oldIllustrations, newIllustrations)
		// 归一化插图 URL
		for i := range removedIllustrations {
			removedIllustrations[i] = utils.DBURLFromPublic(removedIllustrations[i])
		}
		for i := range addedIllustrations {
			addedIllustrations[i] = utils.DBURLFromPublic(addedIllustrations[i])
		}
		if err := utils.InitImagesCategory(tx, removedIllustrations); err != nil {
			return err
		}
		if err := utils.ChangeImagesCategory(tx, addedIllustrations, appTypes.Illustration); err != nil {
			return err
		}

		return articleService.Update(req.ID, articleToUpdate)
	})
}

func (articleService *ArticleService) ArticleList(info request.ArticleList) (interface{}, int64, error) {
	req := &search.Request{
		Query: &types.Query{},
	}

	boolQuery := &types.BoolQuery{}

	// 根据标题查询
	if info.Title != nil {
		boolQuery.Must = append(boolQuery.Must, types.Query{Match: map[string]types.MatchQuery{"title": {Query: *info.Title}}})
	}

	// 根据简介查询
	if info.Abstract != nil {
		boolQuery.Must = append(boolQuery.Must, types.Query{Match: map[string]types.MatchQuery{"abstract": {Query: *info.Abstract}}})
	}

	// 根据类别筛选
	if info.Category != nil {
		boolQuery.Filter = []types.Query{
			{
				Term: map[string]types.TermQuery{
					"category": {Value: info.Category},
				},
			},
		}
	}

	// 根据条件执行查询
	if boolQuery.Must != nil || boolQuery.Filter != nil {
		req.Query.Bool = boolQuery
	} else {
		req.Query.MatchAll = &types.MatchAllQuery{}
		req.Sort = []types.SortCombinations{
			types.SortOptions{
				SortOptions: map[string]types.FieldSort{
					"created_at": {Order: &sortorder.Desc},
				},
			},
		}
	}

	option := other.EsOption{
		PageInfo: info.PageInfo,
		Index:    elasticsearch.ArticleIndex(),
		Request:  req,
	}
	list, total, err := utils.EsPagination(context.TODO(), option)
	if err != nil {
		return nil, 0, err
	}
	// 拼接封面URL
	articleService.processCoverURLs(list)
	return list, total, nil
}
