package model

import "time"

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type ChargesInfo struct {
	Filename string
	ModTime  time.Time
}

type ChargeResponse struct {
	Filename string `json:"filename"`
	Date     string `json:"date"`
}
