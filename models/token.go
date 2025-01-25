package models

type Token struct {
	UserID int    `json:"userID"`
	Name   string `json:"name"`
	Role   string `json:"role"`
	Token  string `json:"token"`
}
