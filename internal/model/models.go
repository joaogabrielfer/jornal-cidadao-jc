package model

import(
	"time"
	"fmt"
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
	ID 		 int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type Charge struct {
	ID 			int 			`json:"id"`
	URL 		string			`json:"url"`
	Filename 	string			`json:"filename"`
	Title 		string			`json:"title"`
	Date	 	FormattedTime 	`json:"date"`
}

type Article struct{
	ID 		int		`json:"id"`
	Title 	string	`json:"title"`
	Author 	string	`json:"author"`
	Body 	string	`json:"body"`
	Poll 	*Poll	`json:"poll,omitempty"`
}

type PollOption struct {
	ID         int    `json:"id"`
	OptionText string `json:"option_text"`
	Votes      int    `json:"votes"`
	PollID     int    `json:"-"` 
}

type Poll struct{
	ID 			int				`json:"id"`
	Question 	string			`json:"question"`
	ArticleID 	int 			`json:"article_id"`
	Options 	[]PollOption	`json:"options"`
}
