package models

type KindData struct {
	ID         int    `json:"id" db:"id"`
	Name       string `json:"name" db:"name"`
	Icon       string `json:"icon" db:"icon"`
	EssayCount int8   `json:"essayCount" db:"essay_count"`
}

type KindParams struct {
	Name string `json:"name" binding:"required" db:"name"`
	Icon string `json:"icon" binding:"required" db:"icon"`
}

type KindUpdateParams struct {
	KindParams
	ID int `json:"id" binding:"required" db:"id"`
}
