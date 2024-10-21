package models

type EssayContent struct {
	Name          string      `json:"name" db:"name"`
	KindID        int         `json:"kind_id,omitempty" db:"kind_id"`
	KindName      string      `json:"kindName" db:"kind_name"`
	LabelList     []LabelData `json:"labelList"`
	Id            int         `json:"id" db:"id"`
	Introduction  string      `json:"introduction" db:"introduction"`
	Content       string      `json:"content" db:"content"`
	VisitedTimes  int64       `json:"visitedTimes" db:"visited_times"`
	CreatedTime   string      `json:"createdTime" db:"created_time"`
	Keywords      []string    `json:"keywords,omitempty"`
	NearEssayList []EssayData `json:"nearEssayList,omitempty"`
}

type EssayParams struct {
	ID           int      `json:"id" db:"id"`
	Name         string   `json:"name" binding:"required" db:"name"`
	KindID       int      `json:"kindID" binding:"required" db:"kind_id"`
	LabelIds     []int    `json:"labelIds" bind:"required"`
	Introduction string   `json:"introduction" binding:"required" db:"introduction"`
	CreatedTime  string   `json:"createdTime" db:"created_time"`
	Content      string   `json:"content" binding:"required" db:"content"`
	ImgUrl       string   `json:"imgUrl" binging:"required" db:"img_url"`
	IfTop        bool     `json:"ifTop" binging:"required" db:"if_top"`
	IfRecommend  bool     `json:"ifRecommend"  binging:"required" db:"if_recommend"`
	Keywords     []string `json:"keywords"`
}

type EssayUpdateParams struct {
	EssayParams
	OldKindID   int   `json:"oldKindID" binding:"required"`
	OldLabelIds []int `json:"OldLabelIds" binding:"required"`
}
