// Package elasticsearch Elasticsearch 客户端封装。
package elasticsearch

import (
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
)

// New 创建 Elasticsearch 客户端。
func New(addresses []string, username, password string) (*elasticsearch.TypedClient, error) {
	cfg := elasticsearch.Config{
		Addresses: addresses,
	}

	if username != "" && password != "" {
		cfg.Username = username
		cfg.Password = password
	}

	client, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("elasticsearch connection error: %w", err)
	}

	return client, nil
}

// ArticleIndex 返回文章索引名称。
func ArticleIndex() string {
	return "blog_articles"
}
