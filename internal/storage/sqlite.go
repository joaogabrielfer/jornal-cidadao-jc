package storage

import (
	"database/sql"
	"log"
	"time"

	"github/jornal-cidadao-jc/internal/model"
)

type Storage struct {
	DB *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{DB: db}
}

func (s *Storage) InitializeDatabase() {
	createUserTableSQL := `CREATE TABLE IF NOT EXISTS users (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"username" TEXT UNIQUE,
		"email" TEXT UNIQUE,
		"password_hash" TEXT
	  );`
	statement, err := s.DB.Prepare(createUserTableSQL)
	if err != nil {
		log.Fatal("Erro preparando statement de criar tabela de usuários", err)
	}
	statement.Exec()
	log.Println("Tabela 'users' foi criada com sucesso ou ja existe")

	createChargeTableSQL := `CREATE TABLE IF NOT EXISTS charges (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		filename TEXT NOT NULL UNIQUE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	statement, err = s.DB.Prepare(createChargeTableSQL)
	if err != nil {
		log.Fatal("Erro preparando statement de criar tabela de charges", err)
	}
	statement.Exec()
	log.Println("Tabela 'charges' foi criada com sucesso ou ja existe")

		createArticleTableSQL := `CREATE TABLE IF NOT EXISTS article (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		author TEXT NOT NULL,
		body TEXT NOT NULL
	);`
	statement, err = s.DB.Prepare(createArticleTableSQL)
	if err != nil {
		log.Fatal("Erro preparando statement de criar tabela de artigos", err)
	}
	statement.Exec()
	log.Println("Tabela 'article' foi criada com sucesso ou ja existe")

}

func (s *Storage) CreateUser(username, email, passwordHash string) error {
	insertSQL := `INSERT INTO users(username, email, password_hash) VALUES (?, ?, ?)`
	statement, err := s.DB.Prepare(insertSQL)
	if err != nil {
		log.Println("Erro preparando statement de insert", err)
		return err
	}
	_, err = statement.Exec(username, email, passwordHash)
	if err != nil {
		log.Println("Erro executando statement", err)
		return err
	}
	return nil
}

func (s *Storage) GetUsers() ([]model.User, error) {
	rows, err := s.DB.Query("SELECT id, username, email FROM users")
	if err != nil {
		log.Println("Erro buscando usuarios: ", err)
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email); err != nil {
			log.Println("Erro escaneando usuario", err)
			continue
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		log.Println("Erro durante iteraçao de linhas", err)
		return nil, err
	}
	return users, nil
}

func (s *Storage) DeleteUser(id int) error {
	_, err := s.DB.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}

func (s *Storage) CreateCharge(title, filename string) error {
	insertSQL := `INSERT INTO charges (title, filename) VALUES (?, ?)`
	_, err := s.DB.Exec(insertSQL, title, filename)
	return err
}

func (s *Storage) GetAllCharges() ([]model.Charge, error) {
	query := `SELECT id, title, filename, created_at FROM charges ORDER BY created_at DESC`
	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var charges []model.Charge
	for rows.Next() {
		var charge model.Charge
		var createdAt time.Time
		if err := rows.Scan(&charge.ID, &charge.Title, &charge.Filename, &createdAt); err != nil {
			return nil, err
		}
		charge.Date = model.FormattedTime(createdAt)
		charges = append(charges, charge)
	}
	return charges, nil
}

func (s *Storage) GetChargeByID(id int) (model.Charge, error) {
	var charge model.Charge
	var createdAt time.Time

	query := `SELECT id, title, filename, created_at FROM charges WHERE id = ?`
	
	err := s.DB.QueryRow(query, id).Scan(&charge.ID, &charge.Title, &charge.Filename, &createdAt)
	if err != nil {
		return model.Charge{}, err
	}
	
	charge.Date = model.FormattedTime(createdAt)
	return charge, nil
}

func (s *Storage) DeleteCharge(id int) (string, error) {
	var filename string
	err := s.DB.QueryRow("SELECT filename FROM charges WHERE id = ?", id).Scan(&filename)
	if err != nil {
		return "", err
	}

	_, err = s.DB.Exec("DELETE FROM charges WHERE id = ?", id)
	if err != nil {
		return "", err
	}
	return filename, nil
}

func (s *Storage) GetUserByID(id int) (model.User, error) {
	var user model.User
	query := `SELECT id, username, email FROM users WHERE id = ?`

	err := s.DB.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (s *Storage) CreateArticle(title, author, body string) error{
	insertSQL := `INSERT INTO article (title, author, body) VALUES (?, ?, ?)`
	
	_, err := s.DB.Exec(insertSQL, title, author, body)
	return err
}

func (s *Storage) UpdateArticle(id int, title, author, body string) error {
	updateSQL := `UPDATE article SET title = ?, author = ?, body = ? WHERE id = ?`
	
	_, err := s.DB.Exec(updateSQL, title, author, body, id)
	return err
}

func (s *Storage) GetArticles() ([]model.Article, error) {
	query := `SELECT id, title, author, body FROM article`
	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []model.Article
	for rows.Next() {
		var article model.Article
		if err := rows.Scan(&article.ID, &article.Title, &article.Author, &article.Body); err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}
	return articles, nil
}

func (s *Storage) GetArticleByID(id int) (model.Article, error) {
	var article model.Article

	query := `SELECT id, title, author, body FROM article WHERE id = ?`
	
	err := s.DB.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Author, &article.Body)
	if err != nil {
		return model.Article{}, err
	}
	
	return article, nil
}

func (s *Storage) DeleteArticle(id int) error {
	_, err := s.DB.Exec("DELETE FROM article WHERE id = ?", id)
	return err
}
