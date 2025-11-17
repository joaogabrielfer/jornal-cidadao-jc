package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github/jornal-cidadao-jc/internal/handlers"
	"github/jornal-cidadao-jc/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: Arquivo .env não encontrado.")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		log.Fatal("DB_PATH não definido no arquivo .env")
	}

	staticPath := os.Getenv("STATIC_PATH")
	templatesPath := os.Getenv("TEMPLATES_PATH")
	if templatesPath == "" {
		templatesPath = "../../static/templates"
	}
	// Convert to absolute path to avoid issues with relative paths
	if !filepath.IsAbs(templatesPath) {
		absPath, err := filepath.Abs(templatesPath)
		if err == nil {
			templatesPath = absPath
		}
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("Erro ao conectar com o banco de dados: ", err)
	}
	defer db.Close()

	storageLayer := storage.NewStorage(db)
	storageLayer.InitializeDatabase()

	chargesDir := filepath.Join(staticPath, "images", "charges")
	httpHandler := handlers.NewHandler(storageLayer, chargesDir)

	router := gin.Default()
	router.Static("/static", staticPath)
	// Load only .tmpl files to avoid trying to load directories like Styles/
	// Clean the path and ensure no trailing slash or existing glob patterns
	templatesPath = filepath.Clean(templatesPath)
	templatesPath = strings.TrimSuffix(templatesPath, string(filepath.Separator))
	// Remove any existing glob patterns (* or **) from the path
	templatesPath = strings.TrimSuffix(templatesPath, "/*")
	templatesPath = strings.TrimSuffix(templatesPath, "/**")
	// Use forward slashes for glob pattern (works on all platforms)
	templatesPattern := fmt.Sprintf("%s/*.tmpl", strings.ReplaceAll(templatesPath, "\\", "/"))
	log.Printf("Loading templates from pattern: %s", templatesPattern)
	router.LoadHTMLGlob(templatesPattern)

	setupRoutes(router, httpHandler)

	log.Printf("Servidor iniciando na porta %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Erro ao iniciar o servidor: ", err)
	}
}
