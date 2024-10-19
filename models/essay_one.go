package models

type EssayContent struct {
	Name          string      `json:"name" db:"name"`
	KindID        int         `json:"kind_id,omitempty" db:"kind_id"`
	Id            int         `json:"id" db:"id"`
	Introduction  string      `json:"introduction" db:"introduction"`
	Content       string      `json:"content" db:"content"`
	VisitedTimes  int64       `json:"visitedTimes" db:"visitedTimes"`
	CreatedTime   string      `json:"createdTime" db:"createdTime"`
	Keywords      []string    `json:"keywords,omitempty"`
	NearEssayList []EssayData `json:"nearEssayList,omitempty"`
}
