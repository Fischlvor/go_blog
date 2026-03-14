package main

import (
	"log"
	"os"
	"path/filepath"

	"server-blog-v2/config"
	"server-blog-v2/pkg/postgres"

	"gorm.io/gen"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Gen: config error: %s", err)
	}

	pg, err := postgres.New(
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DBName,
		cfg.Postgres.SSLMode,
		cfg.Postgres.TimeZone,
		cfg.Postgres.MaxIdleConns,
		cfg.Postgres.MaxOpenConns,
	)
	if err != nil {
		log.Fatalf("Gen: postgres error: %s", err)
	}

	wd, _ := os.Getwd()

	g := gen.NewGenerator(gen.Config{
		OutPath:          filepath.Join(wd, "internal", "repo", "persistence", "gen", "query"),
		ModelPkgPath:     filepath.Join(wd, "internal", "repo", "persistence", "gen", "model"),
		Mode:             gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable:    true,
		FieldWithTypeTag: true,
	})

	g.UseDB(pg.DB)
	// ensure generated models import gorm for gorm.DeletedAt
	g.WithImportPkgPath("gorm.io/gorm")

	log.Printf("Gen: generating models and queries")

	// Generate models with deleted_at mapped to gorm.DeletedAt for soft delete
	g.ApplyBasic(
		// 用户相关
		g.GenerateModel("users", gen.FieldType("deleted_at", "gorm.DeletedAt")),
		g.GenerateModel("logins", gen.FieldType("deleted_at", "gorm.DeletedAt")),

		// 内容相关
		g.GenerateModel("article_categories", gen.FieldType("deleted_at", "gorm.DeletedAt")),
		g.GenerateModel("article_tags", gen.FieldType("deleted_at", "gorm.DeletedAt")),
		g.GenerateModel("articles", gen.FieldType("deleted_at", "gorm.DeletedAt")),
		g.GenerateModel("article_likes", gen.FieldType("deleted_at", "gorm.DeletedAt")),
		g.GenerateModel("article_views"),

		// 评论相关
		g.GenerateModel("comments", gen.FieldType("deleted_at", "gorm.DeletedAt")),
		g.GenerateModel("comment_likes"),

		// AI 相关
		g.GenerateModel("ai_chat_sessions", gen.FieldType("deleted_at", "gorm.DeletedAt")),
		g.GenerateModel("ai_chat_messages", gen.FieldType("deleted_at", "gorm.DeletedAt")),
		g.GenerateModel("ai_models", gen.FieldType("deleted_at", "gorm.DeletedAt")),

		// 资源相关
		g.GenerateModel("files", gen.FieldType("deleted_at", "gorm.DeletedAt")),
		g.GenerateModel("resources", gen.FieldType("deleted_at", "gorm.DeletedAt")),
		g.GenerateModel("resource_upload_tasks"),

		// 其他
		g.GenerateModel("feedbacks", gen.FieldType("deleted_at", "gorm.DeletedAt")),
		g.GenerateModel("links", gen.FieldType("deleted_at", "gorm.DeletedAt")),
		g.GenerateModel("advertisements"),
		g.GenerateModel("footer_links"),
		g.GenerateModel("emoji_groups"),
		g.GenerateModel("emoji_sprites"),
	)

	g.Execute()
	log.Printf("Gen: generate success, wd=%s", wd)
}
