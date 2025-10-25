package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	"github/jornal-cidadao-jc/internal/handlers"
	"github/jornal-cidadao-jc/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if err := godotenv.Load(); err != nil { log.Println("Aviso: Arquivo .env não encontrado.") }

	port := os.Getenv("PORT")
	if port == "" { port = "8000" }

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" { log.Fatal("DB_PATH não definido no arquivo .env") }

	staticPath := os.Getenv("STATIC_PATH")
	templatesPath := os.Getenv("TEMPLATES_PATH")

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil { log.Fatal("Erro ao conectar com o banco de dados: ", err) }
	defer db.Close()

	storageLayer := storage.NewStorage(db)
	storageLayer.InitializeDatabase()

	chargesDir := filepath.Join(staticPath, "images", "charges")
	httpHandler := handlers.NewHandler(storageLayer, chargesDir)

	router := gin.Default()
	router.Static("/static", staticPath)
	router.LoadHTMLGlob(templatesPath)

	setupRoutes(router, httpHandler)

	log.Printf("Servidor iniciando na porta %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Erro ao iniciar o servidor: ", err)
	}
}
