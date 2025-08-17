package service

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"server/global"
	"server/model/elasticsearch"
)

var ArticleCache *ArticleCacheDB

func NewArticleCacheDB() *ArticleCacheDB {
	return &ArticleCacheDB{
		Index: "article",
	}
}

type ArticleCacheDB struct {
	Index string
}

func init() {
	ArticleCache = NewArticleCacheDB()
}

// Set 在原有基础上加一
//func (c ArticleCacheDB) Add(id string) error {
//	num, _ := global.Redis.HGet(c.Index, id).Int()
//	num++
//	err := global.Redis.HSet(c.Index, id, num).Err()
//	return err
//}

func (c ArticleCacheDB) Get(id string, article *elasticsearch.Article) error {
	// 从 Redis 获取 JSON 字符串
	jsonData, err := global.Redis.HGet(c.Index, id).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("article not found: %w", err)
		}
		return fmt.Errorf("failed to get article from redis: %w", err)
	}

	// 解析 JSON 到 Article 结构体
	if err := json.Unmarshal([]byte(jsonData), &article); err != nil {
		return fmt.Errorf("failed to unmarshal article: %w", err)
	}

	return nil
}

func (c ArticleCacheDB) Set(id string, article *elasticsearch.Article) error {
	jsonData, err := json.Marshal(article)
	if err != nil {
		return fmt.Errorf("failed to marshal article: %v", err)
	}

	err = global.Redis.HSet(c.Index, id, jsonData).Err()
	if err != nil {
		return fmt.Errorf("failed to set article in redis: %v", err)
	}
	return nil
}

func (c ArticleCacheDB) Delete(id string) error {
	// 检查缓存中是否存在该文章
	exists, err := global.Redis.HExists(c.Index, id).Result()
	if err != nil {
		return fmt.Errorf("failed to check article existence in redis: %v", err)
	}
	// 如果存在则删除
	if exists {
		err = global.Redis.HDel(c.Index, id).Err()
		if err != nil {
			return fmt.Errorf("failed to delete article from redis: %v", err)
		}
	}
	return nil
}

// GetAllArticles 取出数据
//func (c ArticleCacheDB) GetAllArticles() ([]elasticsearch.Article, error) {
//	articles := make([]elasticsearch.Article, 0)
//	maps := global.Redis.HGetAll(c.Index).Val()
//	for _, jsonData := range maps {
//		var article elasticsearch.Article
//		if err := json.Unmarshal([]byte(jsonData), &article); err != nil {
//			return []elasticsearch.Article{}, fmt.Errorf("failed to unmarshal article: %w", err)
//		}
//		articles = append(articles, article)
//	}
//	return articles, nil
//}

//func (c ArticleCacheDB) getIndex(id string) string {
//	return c.Index + ":" + id
//}

// Clear 清除数据
func (c ArticleCacheDB) Clear() {
	global.Redis.Del(c.Index)
}
