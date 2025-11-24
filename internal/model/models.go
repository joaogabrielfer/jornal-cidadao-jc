package model

import (
	"fmt"
	"time"
)

type PostStatus int

const (
	StatusEmAnalise PostStatus = iota + 1
	StatusAprovado
	StatusRejeitado
	StatusDesconhecido
)

type FormattedTime time.Time

func (ft FormattedTime) MarshalJSON() ([]byte, error) {
	t := time.Time(ft)
	formatted := t.Format("02-01-2006 15:04:05")
	return []byte(fmt.Sprintf(`"%s"`, formatted)), nil
}

func (ft FormattedTime) Format() string {
	t := time.Time(ft)
	return t.Format("02/01/2006")
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type Charge struct {
	ID       int           `json:"id"`
	URL      string        `json:"url"`
	Filename string        `json:"filename"`
	Title    string        `json:"title"`
	Date     FormattedTime `json:"date"`
}

type Article struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Body   string `json:"body"`
	Poll   *Poll  `json:"poll,omitempty"`
}

type PollOption struct {
	ID         int    `json:"id"`
	OptionText string `json:"option_text"`
	Votes      int    `json:"votes"`
	PollID     int    `json:"-"`
}

type Poll struct {
	ID        int          `json:"id"`
	Question  string       `json:"question"`
	ArticleID int          `json:"article_id"`
	Options   []PollOption `json:"options"`
}

func (s PostStatus) String() string {
	switch s {
	case StatusEmAnalise:
		return "Em análise"
	case StatusAprovado:
		return "Aprovado"
	case StatusRejeitado:
		return "Rejeitada"
	default:
		return "Desconhecido"
	}
}

func ToPostStatus(s string) PostStatus {
	switch s {
	case "analise":
		return StatusEmAnalise
	case "aprovado":
		return StatusAprovado
	case "rejeitado":
		return StatusRejeitado
	default:
		return StatusDesconhecido
	}
}

type Post struct {
	ID                int           `json:"id"`
	Title             string        `json:"titulo"`
	Description       string        `json:"description"`
	MediaURL          string        `json:"media_url"`
	AuthorID          int           `json:"author_id"`
	Status            PostStatus    `json:"status"`
	Date              FormattedTime `json:"date"`
	UltimaAtualizacao *FormattedTime `json:"ultima_atualizacao,omitempty"`
}

type PostStatusLog struct {
	ID         int        `json:"id"`
	PostID     int        `json:"post_id"`
	OldStatus  PostStatus `json:"old_status"`
	NewStatus  PostStatus `json:"new_status"`
	ChangedAt  time.Time  `json:"changed_at"`
}

type PaginatedPosts struct {
	Posts    []Post   `json:"posts"`
	Metadata Metadata `json:"metadata"`
}

type Metadata struct {
	CurrentPage  int `json:"current_page"`
	PageSize     int `json:"page_size"`
	FirstPage    int `json:"first_page"`
	LastPage     int `json:"last_page"`
	TotalRecords int `json:"total_records"`
}

type PostReport struct {
	ID        int           `json:"id"`
	PostID    int           `json:"post_id"`
	Reason    string        `json:"reason"`
	CreatedAt FormattedTime `json:"created_at"`
}
