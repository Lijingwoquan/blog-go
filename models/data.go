package models

type DataAboutIndex struct {
	ClassifyKindName string                     `json:"classifyKind"`
	ClassifyDetails  []DataAboutClassifyDetails `json:"classifyDetails"`
}

type DataAboutClassifyDetails struct {
	ClassifyKindName string                 `json:"classifyKind" db:"classifyKindName"`
	ClassifyName     string                 `json:"classifyName"  db:"classifyName"`
	ClassifyRoute    string                 `json:"classifyRoute" db:"classifyRoute"`
	ClassifyEssay    []ClassifyIncludeEssay `json:"classifyEssay" db:"essayName"`
}

type ClassifyIncludeEssay struct {
	EssayName  string `json:"essayName" db:"essayName"`
	EssayKind  string `json:"essayKind" db:"essayKind"`
	EssayRoute string `json:"essayRoute" db:"essayRoute"`
}

type UpdateEssay struct {
	EssayOldName string `json:"essayOldName" db:"essayName" binding:"required"`
	EssayName    string `json:"essayName" db:"essayName" binding:"required"`
	EssayKind    string `json:"essayKind" db:"essayKind"`
	EssayRoute   string `json:"essayRoute" db:"essayRoute"`
	EssayContent string `json:"essayContent" db:"essayContent" binding:"required"`
}
