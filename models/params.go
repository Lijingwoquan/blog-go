package models

type UserParams struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"rePassword" binding:"required,eqfield=Password"`
	Email      string `json:"email" binding:"required"`
}

type ClassifyParams struct {
	Kind string `json:"kind" binding:"required"`
	Name string `json:"name" binding:"required"`
	Icon string `json:"icon" binding:"omitempty"`
}

type UpdateEssayMsgParams struct {
	Name          string   `json:"name" db:"name" binding:"required"`
	Introduction  string   `json:"introduction" db:"introduction" binding:"required"`
	Kind          string   `json:"kind" db:"kind" binding:"required"`
	Content       string   `json:"content" db:"content" binding:"required"`
	Id            int      `json:"id" db:"id" binding:"required"`
	Keywords      []string `json:"keywords"`
	ImgUrl        string   `json:"imgUrl" binging:"required"`
	AdvertiseMsg  string   `json:"advertiseMsg" db:"advertiseMsg"`
	AdvertiseImg  string   `json:"advertiseImg"  db:"advertiseImg"`
	AdvertiseHref string   `json:"advertiseHref" db:"advertiseHref"`
}

type UpdateKindParams struct {
	Name string `json:"name" binding:"required" db:"name"`
	Icon string `json:"icon" binding:"required" db:"icon"`
	ID   int    `json:"id" binding:"required" db:"id"`
}

type UpdateClassifyParams struct {
	Name string `json:"name" binding:"required" db:"name"`
	ID   int    `json:"id" binding:"required" db:"id"`
}

type EssayQuery struct {
	Page     int
	PageSize int
	LabelID  int
	KindID   int
}

type SearchParam struct {
	Keyword string `json:"keyword" binging:"required"`
	IfAdd   bool   `json:"ifAdd"`
}

type EssayIdAndKeyword struct {
	EssayId  int      `json:"essayId" binging:"required"`
	Keywords []string `json:"keywords"`
}
