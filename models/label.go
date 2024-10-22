package models

type LabelData struct {
	ID    int    `json:"id" db:"id"`
	Name  string `json:"name"  db:"name"`
	Color string `json:"color,omitempty"`
}

type LabelParams struct {
	Name string `json:"name" binding:"required" db:"name"`
}

type LabelUpdateParams struct {
	LabelParams
	ID int `json:"id" binding:"required" db:"id"`
}
