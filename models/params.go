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
	Icon   string `json:"icon" binding:"omitempty"`
}

type EssayParams struct {
	Kind         string `json:"kind" binding:"required" `
	Name         string `json:"name" binding:"required"`
	Introduction string `json:"introduction" binding:"required"`
	Content      string `json:"content" binding:"required"`
	Router       string `json:"router" binding:"required"`
}

type UpdateKindParams struct {
	Name string `json:"name" binding:"required" db:"name"`
	Icon string `json:"icon" binding:"required" db:"icon"`
	ID   int    `json:"id" binding:"required" db:"id"`
}

type UpdateClassifyParams struct {
	Name   string `json:"name" binding:"required" db:"name"`
	Router string `json:"router" binding:"required" db:"router"`
	ID     int    `json:"id" binding:"required" db:"id"`
}
