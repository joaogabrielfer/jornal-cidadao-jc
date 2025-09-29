package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

const PORT string = "8080"

type App struct {
	DB *sql.DB
}

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type ChargesInfo struct{
	Filename string
	ModTime time.Time
}

func initialize_database(db *sql.DB) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS users (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"username" TEXT UNIQUE,
		"email" TEXT UNIQUE,
		"password_hash" TEXT
	  );`

	statement, err := db.Prepare(createTableSQL)
	if err != nil {
		log.Fatal("Erro preparando statement de criar tabela", err)
	}
	statement.Exec()
	log.Println("Tabela 'users' foi criada com sucesso ou ja existe")
}

func main() {
	db, err := sql.Open("sqlite3", "../db/users.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	initialize_database(db)

	app := &App{
		DB: db,
	}

	router := gin.Default()
	router.Static("/static", "../static/")

	router.POST("/signup", app.create_user)
	router.GET("/users", app.get_users)
	router.GET("/charge/:name", app.get_charges)
	router.GET("/charges", app.get_charges_list)
	router.GET("/signup", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`
			<form action="/signup" method="post">
				<input type="text" name="username" placeholder="Username" required><br>
				<input type="email" name="email" placeholder="Email" required><br>
				<input type="password" name="password" placeholder="Password" required><br>
				<button type="submit">Sign Up</button>
			</form>
		`))
	})

	router.Run(":" + PORT)
}

func (app *App) create_user(c *gin.Context) {
	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")

	if username == "" || email == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Todos os campos sao requeridos"})
		return
	}

	hashed_password, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Erro fazendo hash da senha: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falhou em criar conta"})
		return
	}

	insert_sql := `INSERT INTO users(username, email, password_hash) VALUES (?, ?, ?)`
	statement, err := app.DB.Prepare(insert_sql)
	if err != nil {
		log.Println("Erro preparando statement de insert", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falhou em criar a conta"})
		return
	}

	_, err = statement.Exec(username, email, string(hashed_password))
	if err != nil {
		log.Println("Erro executando statement", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nome ou email ja podem estar em uso"})
		return
	}

	c.Redirect(http.StatusSeeOther, "/users")
}

func (app *App) get_users(c *gin.Context){
	rows, err := app.DB.Query("SELECT username, email FROM users")
	if err != nil {
		log.Println("Erro buscando usuarios: ", err)
		c.HTML(http.StatusInternalServerError, "Erro buscando usuarios", nil)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User

		if err := rows.Scan(&user.Username, &user.Email); err != nil {
			log.Println("Erro escaneando usuario", err)
			continue 
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		log.Println("Erro durante iteraçao de linhas", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro processando lista de usuarios"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (app *App) get_charges(c *gin.Context) {
	charge_name := c.Param("name")
    charge_url := filepath.Join("../static", "images", "charges",  charge_name)

    c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`
    <!DOCTYPE html>
    <html lang="en">
        <head>
            <meta charset="UTF-8">
            <title>Charge</title>
        </head>
        <body>
            <img src="`+charge_url+`" alt="Charge">    
        </body>
    </html>
    `))
}

func (app *App) get_charges_list(c *gin.Context){
	charges_dir := "../static/images/charges" 

	files, err := os.ReadDir(charges_dir)
	if err != nil {
		log.Println("Erro lendo diretorio de charges: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Não foi possível ler diretório de charges."})
		return
	}
	var charges []ChargesInfo

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		file_path := charges_dir + "/" + file.Name()
		file_info, err := os.Stat(file_path)
		if err != nil {
			log.Println("Erro obtendo informação da charge: ", file.Name(), err)
			continue 
		}

		charges = append(charges, ChargesInfo{
			Filename: file.Name(),
			ModTime:  file_info.ModTime(),
		})
	}

	if len(charges) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Nenhuma charge encontrada"})
		return
	}

	sort.Slice(charges, func(i, j int) bool {
		return charges[i].ModTime.After(charges[j].ModTime)
	})

	type ChargeResponse struct {
		Filename string `json:"filename"`
		Date     string `json:"date"`
	}
	var responseData []ChargeResponse
	for _, charge := range charges {
		responseData = append(responseData, ChargeResponse{
			Filename: charge.Filename,
			Date: charge.ModTime.Format("02-01-2006 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, responseData)	
}
