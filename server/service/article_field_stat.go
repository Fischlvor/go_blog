package service

import (
	"fmt"
	"github.com/go-redis/redis"
	"server/global"
	"strconv"
)

var ArticleFieldCache *ArticleFieldCacheDB

func NewArticleFieldCacheDB() *ArticleFieldCacheDB {
	return &ArticleFieldCacheDB{
		Index: map[string]string{
			"views":    "article_views",
			"comments": "article_comments",
			"likes":    "article_likes",
		},
	}
}

type ArticleFieldCacheDB struct {
	Index map[string]string
}

func init() {
	ArticleFieldCache = NewArticleFieldCacheDB()
}

// Set 在原有基础上加一
func (c ArticleFieldCacheDB) Add(field, id string) error {
	err := global.Redis.HIncrBy(c.Index[field], id, 1).Err()
	return err
}

func (c ArticleFieldCacheDB) Get(field, id string) (int, error) {
	// 从 Redis 获取 JSON 字符串
	num, err := global.Redis.HGet(c.Index[field], id).Int()
	if err == redis.Nil {
		return 0, fmt.Errorf("field of article %s not found: %w", id, err)
	}
	return num, nil
}

func (c ArticleFieldCacheDB) Set(field, id string, num int) error {
	err := global.Redis.HSet(c.Index[field], id, num).Err()
	if err != nil {
		return fmt.Errorf("failed to set article field %s in redis: %v", field, err)
	}
	return nil
}

func (c ArticleFieldCacheDB) Delete(field, id string) error {
	// 检查缓存中是否存在该文章
	exists, err := global.Redis.HExists(c.Index[field], id).Result()
	if err != nil {
		return fmt.Errorf("failed to check article field existence in redis: %v", err)
	}
	// 如果存在则删除
	if exists {
		err = global.Redis.HDel(c.Index[field], id).Err()
		if err != nil {
			return fmt.Errorf("failed to delete article field %s from redis: %v", field, err)
		}
	}
	return nil
}

// GetInfo 取出数据
func (c ArticleFieldCacheDB) GetAllInfo(field string) map[string]int {
	var Info = map[string]int{}
	maps := global.Redis.HGetAll(c.Index[field]).Val()
	for id, val := range maps {
		num, _ := strconv.Atoi(val)
		Info[id] = num
	}
	return Info
}

func (c ArticleFieldCacheDB) getIndex(field, id string) string {
	return c.Index[field] + ":" + id
}

// Clear 清除数据
func (c ArticleFieldCacheDB) Clear(field string) {
	global.Redis.Del(c.Index[field])
}
