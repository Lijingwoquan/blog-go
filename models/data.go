package models

// index返回数据相关

type DataAboutIndex struct {
	Menu []DataAboutIndexMenu `json:"menu"`
}
type DataAboutIndexMenu struct {
	Kind         DataAboutKind       `json:"kind"`
	ClassifyList []DataAboutClassify `json:"classifyList"`
}

type DataAboutKind struct {
	Name string `json:"name" db:"name"`
	Icon string `json:"icon" db:"icon"`
	Id   int    `json:"id" db:"id"`
}
type DataAboutClassify struct {
	Kind   string `json:"kind" db:"kind"`
	Name   string `json:"name"  db:"name"`
	Router string `json:"router" db:"router"`
	ID     int    `json:"id" db:"id"`
}

// 分页查询相关

type DataAboutEssayListAndPage struct {
	EssayList  *[]DataAboutEssay `json:"list"`
	TotalPages int               `json:"totalPages"`
}
type DataAboutEssay struct {
	Name         string   `json:"name" db:"name"`
	Kind         string   `json:"kind" db:"kind"`
	Router       string   `json:"router" db:"router"`
	KindRouter   string   `json:"kindRouter"`
	Introduction string   `json:"introduction" db:"introduction"`
	ID           int      `json:"id" db:"id"`
	CreatedTime  string   `json:"createdTime" db:"createdTime"`
	Keywords     []string `json:"keywords"`
	VisitedTimes int64    `json:"visitedTimes"`
	ImgUrl       string   `json:"imgUrl" db:"imgUrl"`
	Eid          int64    `db:"eid"`
}

//文章查询相关

type EssayData struct {
	Name          string   `json:"name" db:"name"`
	Kind          string   `json:"kind" db:"kind"`
	Id            int      `json:"id" db:"id"`
	Introduction  string   `json:"introduction" db:"introduction"`
	Router        string   `json:"router"`
	KindRouter    string   `json:"kindRouter"`
	Content       string   `json:"content" db:"content"`
	VisitedTimes  int64    `json:"visitedTimes" db:"visitedTimes"`
	CreatedTime   string   `json:"createdTime" db:"createdTime"`
	UpdatedTime   string   `json:"updatedTime" db:"updatedTime"`
	Keywords      []string `json:"keywords"`
	Eid           int64    `json:"eid" db:"eid"`
	ImgUrl        string   `json:"imgUrl" binging:"required"  db:"imgUrl"`
	AdvertiseMsg  string   `json:"advertiseMsg" db:"advertiseMsg"`
	AdvertiseImg  string   `json:"advertiseImg" db:"advertiseImg"`
	AdvertiseHref string   `json:"advertiseHref" db:"advertiseHref"`
}
