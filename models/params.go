package models

type UserParams struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"rePassword" binding:"required,eqfield=Password"`
	Email      string `json:"email" binding:"required"`
}

type ClassifyParams struct {
	ClassifyKind  string `json:"classifyKind" binding:"required"`
	ClassifyName  string `json:"classifyName" binding:"required"`
	ClassifyRoute string `json:"classifyRoute" binding:"required"`
}

type EssayParams struct {
	EssayKind    string `json:"essayKind" binding:"required"`
	EssayName    string `json:"essayName" binding:"required"`
	EssayContent string `json:"essayContent" binding:"required"`
	EssayRoute   string `json:"essayRoute" binding:"required"`
}
