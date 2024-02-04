package models

type DataAboutIndex struct {
	DataAboutIndexMenu []DataAboutIndexMenu `json:"dataAboutIndexMenu"`
	UserInfo           `json:"userInfo"`
}

type DataAboutIndexMenu struct {
	ClassifyKind    string                     `json:"classifyKind"`
	ClassifyDetails []DataAboutClassifyDetails `json:"classifyDetails"`
}

type DataAboutClassifyDetails struct {
	Kind   string                 `json:"Kind" db:"kind"`
	Name   string                 `json:"name"  db:"name"`
	Router string                 `json:"router" db:"router"`
	Essay  []ClassifyIncludeEssay `json:"essay" db:"name"`
}

type ClassifyIncludeEssay struct {
	Name   string `json:"name" db:"name"`
	Kind   string `json:"kind" db:"kind"`
	Router string `json:"router" db:"router"`
	//文章内容不传过去了 单独写一个接口来获取单个文章的数据
}

type UpdateEssay struct {
	OldName string `json:"oldName" db:"name" binding:"required"`
	Name    string `json:"name" db:"name" binding:"required"`
	Kind    string `json:"kind" db:"kind"`
	Router  string `json:"router" db:"route"`
	Content string `json:"content" db:"content" binding:"required"`
}
