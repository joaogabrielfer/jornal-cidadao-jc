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
	router.LoadHTMLGlob(templatesPath)

	api := router.Group("/api")
	{
		api.POST("/users", httpHandler.CreateUser)
		api.GET("/users", httpHandler.GetUsers)

		api.GET("/charges", httpHandler.GetChargesList)
		api.GET("/charges/random", httpHandler.GetRandomCharge)

		api.GET("/materias", httpHandler.GetArticles)
		api.GET("/materia/:id", httpHandler.GetArticleByID)
	}

	router.GET("/", httpHandler.GetIndexPage)
	router.GET("/charge", httpHandler.GetNoIdChargePage)
	router.GET("/charge/:id", httpHandler.GetChargePage)
	router.GET("/cadastro", httpHandler.GetSignupPage)
	router.GET("/login", httpHandler.GetLoginPage)

	admin := router.Group("/admin")
	{
		admin.GET("/", httpHandler.GetAdminPage)

		admin.GET("/adicionar-charge", httpHandler.GetUploadChargePage)
		admin.GET("/charges", httpHandler.GetDeleteChargePage)
		admin.POST("/charge", httpHandler.UploadCharge)
		admin.DELETE("/charge/:id", httpHandler.DeleteCharge)

		admin.DELETE("/user/:id", httpHandler.DeleteUser)
		admin.GET("/users", httpHandler.GetUsersAdminPage)

		admin.GET("/materia", httpHandler.GetUploadArticlePage)
		admin.GET("/materias", httpHandler.GetArticlesPage)
		admin.POST("/materia", httpHandler.UploadArticle)
		admin.PUT("/materia/:id", httpHandler.UpdateArticle)
		admin.DELETE("/materia/:id", httpHandler.DeleteArticle)
	}

	log.Printf("Servidor iniciando na porta %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Erro ao iniciar o servidor: ", err)
	}
}
