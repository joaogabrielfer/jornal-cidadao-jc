package storage

import (
	"database/sql"
	"log"
	"github/jornal-cidadao-jc/internal/model"
)

type Storage struct {
	DB *sql.DB
}

func New_storage(db *sql.DB) *Storage {
	return &Storage{DB: db}
}

func (s *Storage) Initialize_database() {
	createTableSQL := `CREATE TABLE IF NOT EXISTS users (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"username" TEXT UNIQUE,
		"email" TEXT UNIQUE,
		"password_hash" TEXT
	  );`

	statement, err := s.DB.Prepare(createTableSQL)
	if err != nil {
		log.Fatal("Erro preparando statement de criar tabela", err)
	}
	statement.Exec()
	log.Println("Tabela 'users' foi criada com sucesso ou ja existe")
}

func (s *Storage) Create_user(username, email, password_hash string) error {
	insert_sql := `INSERT INTO users(username, email, password_hash) VALUES (?, ?, ?)`
	statement, err := s.DB.Prepare(insert_sql)
	if err != nil {
		log.Println("Erro preparando statement de insert", err)
		return err
	}

	_, err = statement.Exec(username, email, password_hash)
	if err != nil {
		log.Println("Erro executando statement", err)
		return err
	}
	return nil
}

func (s *Storage) Get_users() ([]model.User, error) {
	rows, err := s.DB.Query("SELECT username, email FROM users")
	if err != nil {
		log.Println("Erro buscando usuarios: ", err)
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.Username, &user.Email); err != nil {
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
