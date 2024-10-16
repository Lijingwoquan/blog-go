package models

// DataAboutEssayListAndPage 分页查询相关
type DataAboutEssayListAndPage struct {
	EssayList  *[]DataAboutEssay `json:"list"`
	TotalPages int               `json:"totalPages"`
}

//文章查询相关

type EssayData struct {
	Name          string        `json:"name" db:"name"`
	Kind          string        `json:"kind" db:"kind"`
	Id            int           `json:"id" db:"id"`
	Introduction  string        `json:"introduction" db:"introduction"`
	Content       string        `json:"content" db:"content"`
	VisitedTimes  int64         `json:"visitedTimes" db:"visitedTimes"`
	CreatedTime   string        `json:"createdTime" db:"createdTime"`
	UpdatedTime   string        `json:"updatedTime" db:"updatedTime"`
	Keywords      []string      `json:"keywords"`
	Eid           int64         `json:"eid" db:"eid"`
	ImgUrl        string        `json:"imgUrl" binging:"required"  db:"imgUrl"`
	AdvertiseMsg  string        `json:"advertiseMsg" db:"advertiseMsg"`
	AdvertiseImg  string        `json:"advertiseImg" db:"advertiseImg"`
	AdvertiseHref string        `json:"advertiseHref" db:"advertiseHref"`
	Next          AdjacentEssay `json:"next"`
	Last          AdjacentEssay `json:"last"`
}

type AdjacentEssay struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
