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

type UpdateEssayMsgParams struct {
	Name         string   `json:"name" db:"name" binding:"required"`
	Introduction string   `json:"introduction" db:"introduction" binding:"required"`
	Kind         string   `json:"kind" db:"kind" binding:"required"`
	Router       string   `json:"router" db:"route" binding:"required"`
	Content      string   `json:"content" db:"content" binding:"required"`
	Id           int      `json:"id" db:"id" binding:"required"`
	Keywords     []string `json:"keywords" binging:"omitempty"`
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

type EssayQuery struct {
	Page     int
	PageSize int
	Classify string
}
