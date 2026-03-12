package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"server-blog-v2/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	// 命令行参数
	action := flag.String("action", "up", "Migration action: up, down, version, force, drop")
	steps := flag.Int("steps", 0, "Number of migrations to apply (for up/down with steps)")
	version := flag.Int("version", 0, "Force migration version (for force action)")
	flag.Parse()

	// 加载配置
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	// 构建数据库 URL
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.DBName,
		cfg.Postgres.SSLMode,
	)

	// 获取迁移文件路径
	wd, _ := os.Getwd()
	migrationsPath := filepath.Join(wd, "migrations")
	sourceURL := "file://" + migrationsPath

	// 创建迁移实例
	m, err := migrate.New(sourceURL, dbURL)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}
	defer m.Close()

	switch *action {
	case "up":
		if *steps > 0 {
			err = m.Steps(*steps)
		} else {
			err = m.Up()
		}
		if err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("Migration up failed: %v", err)
		}
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("✅ No changes to apply")
		} else {
			log.Println("✅ Migration up completed successfully")
		}

	case "down":
		if *steps > 0 {
			err = m.Steps(-*steps)
		} else {
			err = m.Down()
		}
		if err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("Migration down failed: %v", err)
		}
		log.Println("✅ Migration down completed successfully")

	case "version":
		ver, dirty, err := m.Version()
		if err != nil {
			if errors.Is(err, migrate.ErrNilVersion) {
				log.Println("No migrations applied yet")
			} else {
				log.Fatalf("Failed to get version: %v", err)
			}
		} else {
			log.Printf("Current version: %d (dirty: %v)", ver, dirty)
		}

	case "force":
		if *version == 0 {
			log.Fatalf("Please specify -version for force action")
		}
		if err := m.Force(*version); err != nil {
			log.Fatalf("Force version failed: %v", err)
		}
		log.Printf("✅ Forced version to %d", *version)

	case "drop":
		if err := m.Drop(); err != nil {
			log.Fatalf("Drop failed: %v", err)
		}
		log.Println("✅ All tables dropped")

	case "status":
		showTableStatus(dbURL)

	default:
		log.Fatalf("Unknown action: %s. Use: up, down, version, force, drop, status", *action)
	}
}

func showTableStatus(dbURL string) {
	// 使用 database/sql 连接
	db, err := openDB(dbURL)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer db.Close()

	tables := []string{
		"users", "categories", "tags", "articles", "article_tags",
		"article_likes", "article_views", "comments", "comment_likes",
		"chat_sessions", "chat_messages", "feedbacks", "links", "files",
		"advertisements", "footer_links", "emoji_groups", "emojis",
		"emoji_sprites", "emoji_tasks", "resources", "resource_upload_tasks", "logins",
	}

	log.Println("Database tables status:")
	log.Println("========================")

	for _, table := range tables {
		var exists bool
		query := `SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' AND table_name = $1
		)`
		if err := db.QueryRow(query, table).Scan(&exists); err != nil {
			log.Printf("  ❌ %s: error checking", table)
			continue
		}

		if exists {
			log.Printf("  ✅ %s", table)
		} else {
			log.Printf("  ❌ %s: not exists", table)
		}
	}
}

func openDB(dbURL string) (*sql.DB, error) {
	return sql.Open("postgres", dbURL)
}
