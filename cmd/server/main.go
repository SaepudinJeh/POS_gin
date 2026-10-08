package main

import (
	"log"
	"os"

	"POS/internal/configs"
	"POS/internal/container"
	"POS/internal/routes"

	"github.com/joho/godotenv"
	"github.com/thienel/tlog"
)

func main() {
	// 1. Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("File .env tidak ditemukan, pakai environment system")
	}

	// 2. Init logger (sebelum apapun)
	initLogger()
	defer tlog.Sync() // flush buffer saat shutdown

	// 3. Koneksi database
	configs.ConnectDatabase()

	// 4. Rakit dependency
	c := container.NewContainer(configs.DB)

	// 5. Setup router
	r := routes.SetupRouter(&routes.Handlers{
		Auth:     c.Auth,
		Category: c.Category,
		Product:  c.Product,
	})

	// 6. Jalankan
	if err := r.Run(); err != nil {
		tlog.Fatal("Gagal menjalankan server")
	}
}

func initLogger() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	cfg := tlog.DefaultConfig().
		WithEnvironment(env).
		WithLevel(os.Getenv("LOG_LEVEL")).
		WithAppName("pos-api").
		WithVersion("1.0.0")

	// Kalau production, tulis juga ke file
	if env == "production" {
		cfg = cfg.
			WithFile("logs/app.log").
			WithFileRotation(100, 5, 30, true) // 100MB, 5 backups, 30 hari, compress
	}

	if err := tlog.Init(cfg); err != nil {
		log.Fatal("Gagal init logger:", err)
	}
}
