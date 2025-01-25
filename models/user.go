package models

type User struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	IDCard int    `json:"idcard"`
}

type UpdateIDCardRequest struct {
	IDCard string `json:"idcard"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

type ResponseMessage struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}