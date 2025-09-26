package main

import (
	"log"
	"net/http"
	"strconv"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/gin-gonic/gin"
)

const PORT string = "8080"

type Task struct{
	Id int
	Name string
	Done bool
}

type App struct{
	DB *sql.DB
}

func createTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS tasks (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"name" TEXT,
		"done" INTEGER
	  );`

	statement, err := db.Prepare(createTableSQL)
	if err != nil {
		log.Fatal("Error preparing create table statement:", err)
	}
	statement.Exec()
	log.Println("Table 'tasks' created or already exists.")
}

func main(){
	db, err := sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() 
	createTable(db)

	app := &App{
		DB: db,
	}

	router := gin.Default()
	router.LoadHTMLGlob("../frontend/templates/*")

	router.GET("/", app.get_index)
	router.POST("/submit", app.create_task)
	router.POST("/complete/:id", app.complete_task)
	router.POST("/uncomplete/:id", app.uncomplete_task)
	router.POST("/delete/:id", app.delete_task)

	router.Run(":"+PORT)
}

func (app *App) get_index(c *gin.Context){
	tasks := []Task{}

	rows, err := app.DB.Query("SELECT id, name, done FROM tasks")
	if err != nil {
		log.Println("Error querying tasks:", err)
		c.HTML(http.StatusInternalServerError, "Error querying tasks", nil)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var task Task

		if err := rows.Scan(&task.Id, &task.Name, &task.Done); err != nil {
			log.Println("Error scanning task:", err)
			continue 
		}
		tasks = append(tasks, task)
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"tasks": tasks,
	})
}

func (app *App) create_task(c *gin.Context) {
	name := c.PostForm("name")
	done_str := c.PostForm("done")

	var done bool
	if done_str == "on"{
		done = true
	} 

	insert_sql := `INSERT INTO tasks(name, done) VALUES (?, ?)`
	statement, err := app.DB.Prepare(insert_sql)
	if err != nil {
        log.Println("Error preparing insert:", err)
        c.Redirect(http.StatusSeeOther, "/")
        return
    }

	statement.Exec(name, done)
	c.Redirect(http.StatusSeeOther, "/")
}

func (app *App) complete_task(c *gin.Context){
	id_str := c.Param("id")	
	id, err := strconv.Atoi(id_str)
	if err != nil{
		c.String(http.StatusInternalServerError, "Erro convertendo o id")
	}

	update_sql := `UPDATE tasks SET done = 1 WHERE id = ?`
	statement, err := app.DB.Prepare(update_sql)
	if err != nil{
		log.Println("Error preparing update statement:", err)
        return
	}

	_, err = statement.Exec(id)
	if err != nil{
		log.Println("Error executing update statement:", err)
	}

	c.Redirect(http.StatusSeeOther, "/")
}

func (app *App) uncomplete_task(c *gin.Context){
	id_str := c.Param("id")	
	id, err := strconv.Atoi(id_str)
	if err != nil{
		c.String(http.StatusInternalServerError, "Erro convertendo o id")
	}

	update_sql := `UPDATE tasks SET done = 0 WHERE id = ?`
	statement, err := app.DB.Prepare(update_sql)
	if err != nil{
		log.Println("Error preparing update statement:", err)
		return
	}

	_, err = statement.Exec(id)
	if err != nil{
		log.Println("Error executing update statement:", err)
	}

	c.Redirect(http.StatusSeeOther, "/")
}

func (app *App) delete_task(c *gin.Context){
	id_str := c.Param("id")
	id, err := strconv.Atoi(id_str)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid ID format")
		return
	}

	delete_sql := `DELETE FROM tasks WHERE id = ?`
	statement, err := app.DB.Prepare(delete_sql)
	if err != nil{
		log.Println("Error preparing delete statement: ", err)
		return
	}

	_, err = statement.Exec(id)
	if err != nil{
		log.Println("Error executing delete statement: ", err)
	}

	c.Redirect(http.StatusSeeOther, "/")
}
