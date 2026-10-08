package main

import (
	"log"

	"POS/internal/configs"
	"POS/internal/container"
	"POS/internal/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("File .env tidak ditemukan, pakai environment system")
	}

	gin.ForceConsoleColor()

	// 2. Koneksi database & migrasi
	configs.ConnectDatabase()

	// 3. Rakit semua dependency (repo → service → handler)
	c := container.NewContainer(configs.DB)

	// 4. Setup router (middleware + route diurus di sini)
	r := routes.SetupRouter(&routes.Handlers{
		Auth:     c.Auth,
		Category: c.Category,
		Product:  c.Product,
	})

	// 5. Jalankan server
	if err := r.Run(); err != nil {
		log.Fatal("Gagal menjalankan server:", err)
	}
}
