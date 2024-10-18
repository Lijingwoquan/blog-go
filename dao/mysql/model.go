package mysql

import "blog/models"

type rawData struct {
	models.DataAboutEssay
	LabelIDs   string `db:"label_ids"`
	LabelNames string `db:"label_names"`
}
