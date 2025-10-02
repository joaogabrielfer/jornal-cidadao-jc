package model

import "time"

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type ChargesInfo struct {
	URL 		string		`json:"url"`
	Filename 	string		`json:"filename"`
	Title 		string		`json:"title"`
	ModTime  	time.Time 	`json:"modtime"`
}

type ChargeResponse struct {
	Filename string `json:"filename"`
	Date     string `json:"date"`
}
