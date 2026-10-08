package configs

import (
	"fmt"
	"log"
	"os"
	"time"

	"POS/internal/models"

	"github.com/thienel/tlog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Pakai logger tlog untuk SQL logging
		Logger: tlog.NewGormLogger(
			tlog.WithSlowThreshold(200*time.Millisecond),
			tlog.WithIgnoreRecordNotFound(true),
		),
	})
	if err != nil {
		log.Fatal("Gagal koneksi database:", err)
	}

	// Auto migrate
	if err := db.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Product{},
	); err != nil {
		log.Fatal("Gagal migrasi:", err)
	}

	// Partial unique index
	indexes := []string{
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_active
			ON users (email) WHERE deleted_at IS NULL;`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_categories_name_active
			ON categories (name) WHERE deleted_at IS NULL;`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_products_sku_active
			ON products (sku) WHERE deleted_at IS NULL;`,
	}
	for _, idx := range indexes {
		if err := db.Exec(idx).Error; err != nil {
			log.Fatal("Gagal buat index:", err)
		}
	}

	DB = db
	log.Println("Database connected & migrated")
}
