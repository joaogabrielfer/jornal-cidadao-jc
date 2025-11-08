package storage

import (
	"database/sql"
	"log"
	"math"
	"strings"
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

	createEnquetesTableSQL := `CREATE TABLE IF NOT EXISTS poll (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		question TEXT NOT NULL,
		article_id INTEGER NOT NULL UNIQUE,
		FOREIGN KEY (article_id) REFERENCES article(id) ON DELETE CASCADE
	);`

	statement, err = s.DB.Prepare(createEnquetesTableSQL)
	if err != nil {
		log.Fatal("Erro preparando statement de criar tabela de enquetes", err)
	}
	statement.Exec()
	log.Println("Tabela 'poll' foi criada com sucesso ou ja existe")

	createOpcoesTableSQL := `CREATE TABLE IF NOT EXISTS poll_options (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		option_text TEXT NOT NULL,
		votes INTEGER NOT NULL DEFAULT 0,
		poll_id INTEGER NOT NULL,
		FOREIGN KEY (poll_id) REFERENCES poll(id) ON DELETE CASCADE
	);`
	statement, err = s.DB.Prepare(createOpcoesTableSQL)
	if err != nil {
		log.Fatal("Erro preparando statement de criar tabela de opções de enquete", err)
	}
	statement.Exec()
	log.Println("Tabela 'poll_options' foi criada com sucesso ou ja existe")

	alterPostsTableSQL := `ALTER TABLE posts ADD COLUMN ultima_atualizacao DATETIME;`
	s.DB.Exec(alterPostsTableSQL) 

	createPostsTableSQL := ` CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT,
		media_url TEXT NOT NULL,
		author_id INTEGER NOT NULL,
		
		-- 1: Em análise, 2: Aprovada, 3: Rejeitada
		status INTEGER NOT NULL DEFAULT 1, 
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		ultima_atualizacao DATETIME,
		
		CHECK (status IN (1, 2, 3)), 
		
		FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
	);`

	statement, err = s.DB.Prepare(createPostsTableSQL)
	if err != nil {
		log.Fatal("Erro preparando statement de criar tabela de posts", err)
	}
	statement.Exec()
	log.Println("Tabela 'posts' foi criada com sucesso ou ja existe")

	createLogsTableSQL := `
	CREATE TABLE IF NOT EXISTS post_status_logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER NOT NULL,
		old_status INTEGER NOT NULL,
		new_status INTEGER NOT NULL,
		changed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
	);`
	statement, err = s.DB.Prepare(createLogsTableSQL)
	if err != nil {
		log.Fatal("Erro preparando statement de criar tabela de logs de status", err)
	}
	statement.Exec()
	log.Println("Tabela 'post_status_logs' foi criada com sucesso ou ja existe")
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

func (s *Storage) CreateArticleWithPoll(title, author, body, pollQuestion string, pollOptions []string) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	insertArticleSQL := `INSERT INTO article (title, author, body) VALUES (?, ?, ?)`
	res, err := tx.Exec(insertArticleSQL, title, author, body)
	if err != nil {
		return err
	}
	articleID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	if pollQuestion != "" && len(pollOptions) > 0 {
		var validOptions []string
		for _, opt := range pollOptions {
			if strings.TrimSpace(opt) != "" {
				validOptions = append(validOptions, opt)
			}
		}

		if len(validOptions) > 0 {
			insertPollSQL := `INSERT INTO poll (question, article_id) VALUES (?, ?)`
			res, err = tx.Exec(insertPollSQL, pollQuestion, articleID)
			if err != nil {
				return err
			}
			pollID, err := res.LastInsertId()
			if err != nil {
				return err
			}

			insertOptionSQL := `INSERT INTO poll_options (option_text, poll_id) VALUES (?, ?)`
			for _, optionText := range validOptions {
				if _, err := tx.Exec(insertOptionSQL, optionText, pollID); err != nil {
					return err
				}
			}
		}
	}
	return tx.Commit()
}

func (s *Storage) UpdateArticleWithPoll(articleID int, title, author, body, pollQuestion string, pollOptions []string) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	updateArticleSQL := `UPDATE article SET title = ?, author = ?, body = ? WHERE id = ?`
	if _, err := tx.Exec(updateArticleSQL, title, author, body, articleID); err != nil {
		tx.Rollback()
		return err
	}

	deleteEnqueteSQL := `DELETE FROM poll WHERE article_id = ?`
	if _, err := tx.Exec(deleteEnqueteSQL, articleID); err != nil {
		tx.Rollback()
		return err
	}

	if pollQuestion != "" && len(pollOptions) > 0 {
		insertEnqueteSQL := `INSERT INTO poll (question, article_id) VALUES (?, ?)`
		res, err := tx.Exec(insertEnqueteSQL, pollQuestion, articleID)
		if err != nil {
			tx.Rollback()
			return err
		}
		enqueteID, err := res.LastInsertId()
		if err != nil {
			tx.Rollback()
			return err
		}

		insertOpcaoSQL := `INSERT INTO poll_options (option_text, poll_id) VALUES (?, ?)`
		for _, optionText := range pollOptions {
			if _, err := tx.Exec(insertOpcaoSQL, optionText, enqueteID); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit()
}

func (s *Storage) GetArticles() ([]model.Article, error) {
	query := `
		SELECT
			a.id, a.title, a.author, a.body,
			p.id, p.question,
			po.id, po.option_text, po.votes
		FROM article AS a
		LEFT JOIN poll AS p ON a.id = p.article_id
		LEFT JOIN poll_options AS po ON p.id = po.poll_id
		ORDER BY a.id DESC, po.id ASC;
	`
	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	articlesMap := make(map[int]*model.Article)
	for rows.Next() {
		var articleID int
		var articleTitle, articleAuthor, articleBody string
		var pollID sql.NullInt64
		var pollQuestion sql.NullString
		var optionID sql.NullInt64
		var optionText sql.NullString
		var optionVotes sql.NullInt64

		if err := rows.Scan(
			&articleID, &articleTitle, &articleAuthor, &articleBody,
			&pollID, &pollQuestion,
			&optionID, &optionText, &optionVotes,
		); err != nil {
			return nil, err
		}

		if _, exists := articlesMap[articleID]; !exists {
			articlesMap[articleID] = &model.Article{
				ID:     articleID,
				Title:  articleTitle,
				Author: articleAuthor,
				Body:   articleBody,
			}
		}

		article := articlesMap[articleID]
		if pollID.Valid && article.Poll == nil {
			article.Poll = &model.Poll{
				ID:        int(pollID.Int64),
				Question:  pollQuestion.String,
				ArticleID: articleID,
				Options:   []model.PollOption{},
			}
		}

		if optionID.Valid && article.Poll != nil {
			article.Poll.Options = append(article.Poll.Options, model.PollOption{
				ID:         int(optionID.Int64),
				OptionText: optionText.String,
				Votes:      int(optionVotes.Int64),
			})
		}
	}

	var articles []model.Article
	for i := range articlesMap {
		articles = append(articles, *articlesMap[i])
	}
	return articles, nil
}

func (s *Storage) GetArticleByID(id int) (model.Article, error) {
	query := `
		SELECT
			a.id, a.title, a.author, a.body,
			p.id, p.question,
			po.id, po.option_text, po.votes
		FROM article AS a
		LEFT JOIN poll AS p ON a.id = p.article_id
		LEFT JOIN poll_options AS po ON p.id = po.poll_id
		WHERE a.id = ?
		ORDER BY po.id ASC;
	`
	rows, err := s.DB.Query(query, id)
	if err != nil {
		return model.Article{}, err
	}
	defer rows.Close()

	var article model.Article
	var hasData bool = false

	for rows.Next() {
		if !hasData {
			hasData = true
		}

		var pollID sql.NullInt64
		var pollQuestion sql.NullString
		var optionID sql.NullInt64
		var optionText sql.NullString
		var optionVotes sql.NullInt64

		if article.ID == 0 {
			if err := rows.Scan(
				&article.ID, &article.Title, &article.Author, &article.Body,
				&pollID, &pollQuestion,
				&optionID, &optionText, &optionVotes,
			); err != nil {
				return model.Article{}, err
			}
		} else {
			if err := rows.Scan(
				new(int), new(string), new(string), new(string),
				&pollID, &pollQuestion,
				&optionID, &optionText, &optionVotes,
			); err != nil {
				return model.Article{}, err
			}
		}

		if pollID.Valid && article.Poll == nil {
			article.Poll = &model.Poll{
				ID:        int(pollID.Int64),
				Question:  pollQuestion.String,
				ArticleID: article.ID,
				Options:   []model.PollOption{},
			}
		}
		if optionID.Valid && article.Poll != nil {
			article.Poll.Options = append(article.Poll.Options, model.PollOption{
				ID:         int(optionID.Int64),
				OptionText: optionText.String,
				Votes:      int(optionVotes.Int64),
			})
		}
	}

	if !hasData {
		return model.Article{}, sql.ErrNoRows
	}

	return article, nil
}

func (s *Storage) DeleteArticle(id int) error {
	_, err := s.DB.Exec("DELETE FROM article WHERE id = ?", id)
	return err
}

func (s *Storage) VotePoll(optionID int) error {
	updateVoteSQL := `UPDATE poll_options SET votes = votes + 1 WHERE id = ?`
	res, err := s.DB.Exec(updateVoteSQL, optionID)
	if err != nil {
		log.Println("Erro ao atualizar os votos:", err)
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (s *Storage) CreatePost(title, description, mediaUrl string, authorID int, status model.PostStatus) error {
	insertSQL := `INSERT INTO posts (title, description, media_url, author_id, status) VALUES (?, ?, ?, ?, ?)`
	_, err := s.DB.Exec(insertSQL, title, description, mediaUrl, authorID, status)
	return err
}

func (s *Storage) GetPostByID(id int) (model.Post, error) {
	var post model.Post
	var createdAt, updatedAt time.Time
	query := `SELECT id, title, description, media_url, author_id, status, created_at, ultima_atualizacao FROM posts WHERE id = ?`

	err := s.DB.QueryRow(query, id).Scan(
		&post.ID,
		&post.Title,
		&post.Description,
		&post.MediaURL,
		&post.AuthorID,
		&post.Status,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		log.Println("Erro ao escanear um post:", err)
		return model.Post{}, err
	}

	post.Date = model.FormattedTime(createdAt)
	post.UltimaAtualizacao = model.FormattedTime(updatedAt)

	return post, nil
}

func (s *Storage) GetApprovedPostsPaginated(page, pageSize int) ([]model.Post, model.Metadata, error) {
	statusAprovada := model.StatusAprovado

	countQuery := `SELECT COUNT(*) FROM posts WHERE status = ?`
	
	var totalRecords int
	err := s.DB.QueryRow(countQuery, statusAprovada).Scan(&totalRecords)
	if err != nil {
		return nil, model.Metadata{}, err
	}

	query := `
		SELECT id, title, description, media_url, author_id, status, created_at, ultima_atualizacao
		FROM posts 
		WHERE status = ? 
		ORDER BY created_at DESC 
		LIMIT ? OFFSET ?`

	limit := pageSize
	offset := (page - 1) * pageSize

	rows, err := s.DB.Query(query, statusAprovada, limit, offset)
	if err != nil {
		return nil, model.Metadata{}, err
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var post model.Post
		var createdAt, updatedAt time.Time
		if err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Description,
			&post.MediaURL,
			&post.AuthorID,
			&post.Status,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, model.Metadata{}, err
		}
		post.Date = model.FormattedTime(createdAt)
		post.UltimaAtualizacao = model.FormattedTime(updatedAt)
		posts = append(posts, post)
	}
	if err = rows.Err(); err != nil {
		return nil, model.Metadata{}, err
	}

	lastPage := 0
	if totalRecords > 0 {
		lastPage = int(math.Ceil(float64(totalRecords) / float64(pageSize)))
	}
	
	metadata := model.Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     lastPage,
		TotalRecords: totalRecords,
	}

	return posts, metadata, nil
}

func (s *Storage) GetPostsByAuthorID(authorID int) ([]model.Post, error) {
	query := `SELECT id, title, description, media_url, author_id, status, ultima_atualizacao, created_at FROM posts WHERE author_id = ? ORDER BY created_at DESC`

	rows, err := s.DB.Query(query, authorID)
	if err != nil {
		log.Println("Erro ao buscar posts por autor:", err)
		return nil, err
	}
	defer rows.Close()

	var posts []model.Post
	
	for rows.Next() {
		var post model.Post
		var createdAt, updatedAt time.Time

		if err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Description,
			&post.MediaURL,
			&post.AuthorID,
			&post.Status,
			&createdAt,
			&updatedAt,
		); err != nil {
			log.Println("Erro ao escanear um post:", err)
			return nil, err 
		}

		post.Date = model.FormattedTime(createdAt)
		post.UltimaAtualizacao = model.FormattedTime(updatedAt)
		
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		log.Println("Erro durante a iteração dos posts por autor:", err)
		return nil, err
	}

	return posts, nil
}

func (s *Storage) UpdatePostStatus(postID int, newStatus model.PostStatus) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var oldStatus model.PostStatus
	queryOldStatus := `SELECT status FROM posts WHERE id = ?`
	err = tx.QueryRow(queryOldStatus, postID).Scan(&oldStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			return sql.ErrNoRows 
		}
		return err
	}

	if oldStatus == newStatus {
		return nil 
	}

	updateTime := time.Now()
	updatePostSQL := `UPDATE posts SET status = ?, ultima_atualizacao = ? WHERE id = ?`
	_, err = tx.Exec(updatePostSQL, newStatus, updateTime, postID)
	if err != nil {
		return err
	}

	insertLogSQL := `INSERT INTO post_status_logs (post_id, old_status, new_status) VALUES (?, ?, ?)`
	_, err = tx.Exec(insertLogSQL, postID, oldStatus, newStatus)
	if err != nil {
		return err
	}

	return tx.Commit()
}
