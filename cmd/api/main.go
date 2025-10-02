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

	db_path := os.Getenv("DB_PATH")
	if db_path == "" {
		log.Fatal("DB_PATH não definido no arquivo .env")
	}

	static_path := os.Getenv("STATIC_PATH")
	templates_path := os.Getenv("TEMPLATES_PATH")

	db, err := sql.Open("sqlite3", db_path)
	if err != nil {
		log.Fatal("Erro ao conectar com o banco de dados: ", err)
	}
	defer db.Close()

	storage_layer := storage.New_storage(db)
	storage_layer.Initialize_database()

	charges_dir := filepath.Join(static_path, "images", "charges")
	http_handler := handlers.New_handler(storage_layer, charges_dir)

	router := gin.Default()
	router.Static("/static", static_path)
	router.LoadHTMLGlob(templates_path)

	router.POST("/api/cadastro", http_handler.Create_user)
	router.GET("/api/users", http_handler.Get_users)
	router.GET("/api/charges", http_handler.Get_charges_list)
	router.GET("/api/charges/random", http_handler.Get_random_charge)
	
	router.GET("/cadastro", http_handler.Get_signup_page)
	router.GET("/", http_handler.Get_index_page)
	
	log.Printf("Servidor iniciando na porta %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Erro ao iniciar o servidor: ", err)
	}
}

