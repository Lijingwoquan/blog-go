package models

type UserParams struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"rePassword" binding:"required,eqfield=Password"`
	Email      string `json:"email" binding:"required"`
}

type ClassifyParams struct {
	Kind   string `json:"kind" binding:"required"`
	Name   string `json:"name" binding:"required"`
	Router string `json:"router" binding:"required"`
}

type EssayParams struct {
	Kind    string `json:"kind" binding:"required"`
	Name    string `json:"name" binding:"required"`
	Content string `json:"content" binding:"required"`
	Router  string `json:"router" binding:"required"`
}
