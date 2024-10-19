package mysql

import "blog/models"

type rawData struct {
	models.EssayData
	LabelIDs   string `db:"label_ids"`
	LabelNames string `db:"label_names"`
}
