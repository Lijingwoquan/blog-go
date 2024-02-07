package models

type DataAboutIndex struct {
	DataAboutIndexMenu []DataAboutIndexMenu `json:"dataAboutIndexMenu"`
	UserInfo           `json:"userInfo"`
}

type DataAboutIndexMenu struct {
	DataAboutClassify
	ClassifyDetails []DataAboutClassifyDetails `json:"classifyDetails"`
}

type DataAboutClassify struct {
	ClassifyKind string `json:"classifyKind" db:"name"`
	Icon         string `json:"icon" db:"icon"`
}

type DataAboutClassifyDetails struct {
	Kind   string           `json:"Kind" db:"kind"`
	Name   string           `json:"name"  db:"name"`
	Router string           `json:"router" db:"router"`
	ID     int              `json:"id" db:"id"`
	Essay  []DataAboutEssay `json:"essay" db:"name"`
}

type DataAboutEssay struct {
	Name         string `json:"name" db:"name"`
	Kind         string `json:"kind" db:"kind"`
	Router       string `json:"router" db:"router"`
	Introduction string `json:"introduction" db:"introduction"`
	ID           int    `json:"id" db:"id"`
	CreatedTime  string `json:"createdTime" db:"createdTime"`
	UpdatedTime  string `json:"updatedTime" db:"updatedTime"`
}

type UpdateEssay struct {
	OldName string `json:"oldName" db:"name" binding:"required"`
	Name    string `json:"name" db:"name" binding:"required"`
	Kind    string `json:"kind" db:"kind"`
	Router  string `json:"router" db:"route"`
	Content string `json:"content" db:"content" binding:"required"`
}

type EssayData struct {
	Name         string `json:"name" db:"name"`
	Kind         string `json:"kind" db:"kind"`
	Introduction string `json:"introduction" db:"introduction"`
	Content      string `json:"content" db:"content"`
	CreatedTime  string `json:"createdTime" db:"createdTime"`
	UpdatedTime  string `json:"updatedTime" db:"updatedTime"`
}
