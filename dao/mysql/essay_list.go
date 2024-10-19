package mysql

import (
	"blog/models"
	"fmt"
	"strconv"
	"strings"
)

func GetRecommendEssayList(data *[]models.EssayData) error {
	sqlStr := `
		SELECT e.id, e.name, e.createdTime, e.imgUrl 
		FROM essay e 
		WHERE ifRecommend = true
	`
	return db.Select(data, sqlStr)
}

func GetEssayList(data *models.EssayListAndPage, query models.EssayQuery) error {
	// 计算偏移量
	offset := (query.Page - 1) * query.PageSize
	baseSelect := `
		 SELECT  e.id, e.name, e.kind_id, e.ifRecommend, e.introduction, e.createdTime, e.visitedTimes, e.imgUrl,
            k.name AS kind_name,
            GROUP_CONCAT(el.label_id) AS label_ids,
            GROUP_CONCAT(el.label_name) AS label_names
        FROM essay e
        LEFT JOIN kind k ON e.kind_id = k.id
        LEFT JOIN essay_label el ON e.id = el.essay_id`

	baseCount := `
		SELECT COUNT(DISTINCT e.id)
        FROM essay e 
        LEFT JOIN kind k ON e.kind_id = k.id
        LEFT JOIN essay_label el ON e.id = el.essay_id
		`

	where := make([]string, 0)
	args := make([]interface{}, 0)
	if query.LabelID != 0 {
		where = append(where, "el.label_id = ?")
		args = append(args, query.LabelID)
	}
	if query.KindID != 0 {
		where = append(where, "e.kind_id = ?")
		args = append(args, query.KindID)
	}
	whereClause := ""
	if len(where) > 0 {
		whereClause = "WHERE " + strings.Join(where, " AND ")
	}

	// 添加GROUP BY子句
	groupBy := "GROUP BY e.id, e.name, e.kind_id, e.ifRecommend, e.introduction, e.createdTime, e.visitedTimes, e.imgUrl, k.name"

	// 构建完整的SQL语句
	sqlStr := fmt.Sprintf("%s %s %s ORDER BY e.id DESC LIMIT ? OFFSET ?",
		baseSelect, whereClause, groupBy)
	countSqlStr := fmt.Sprintf("%s %s", baseCount, whereClause)
	args = append(args, query.PageSize, offset)

	var rawDataList []rawData

	// 执行分页查询
	if err := db.Select(&rawDataList, sqlStr, args...); err != nil {
		return err
	}

	// 处理查询结果
	data.EssayList = make([]models.EssayData, len(rawDataList))
	for i, raw := range rawDataList {
		data.EssayList[i] = raw.EssayData
		// 处理标签数据
		if raw.LabelIDs != "" && raw.LabelNames != "" {
			ids := strings.Split(raw.LabelIDs, ",")
			names := strings.Split(raw.LabelNames, ",")
			data.EssayList[i].LabelList = make([]models.LabelData, len(ids))

			for j := range ids {
				id, _ := strconv.Atoi(ids[j])
				name := names[j]
				data.EssayList[i].LabelList[j] = models.LabelData{
					ID:   id,
					Name: name,
				}
			}
		}
	}

	// 执行总记录统计
	var totalCount int
	if err := db.Get(&totalCount, countSqlStr, args[:len(args)-2]...); err != nil {
		return err
	}

	// 计算总页数
	data.TotalPage = (totalCount + query.PageSize - 1) / query.PageSize

	return nil
}

func GetAllEssay(data *[]models.EssayData) error {
	var rawDataList = new([]rawData)
	sqlStr := `
		SELECT e.id, e.name, e.createdTime, e.imgUrl,e.kind_id,
		       k.name AS kind_name,
		       GROUP_CONCAT(el.label_id) AS label_ids ,GROUP_CONCAT(el.label_name) AS label_names
		FROM essay e
		LEFT JOIN kind k on e.kind_id = k.id
		LEFT JOIN essay_label el on e.id = el.essay_id
		GROUP BY e.id, e.name, e.createdTime, e.imgUrl, e.kind_id,  k.name
		ORDER BY id DESC
	`
	var err error
	if err = db.Select(rawDataList, sqlStr); err != nil {
		return err
	}
	*data = make([]models.EssayData, len(*rawDataList))
	for i, raw := range *rawDataList {
		(*data)[i] = raw.EssayData
		if raw.LabelNames != "" && raw.LabelIDs != "" {
			ids := strings.Split(raw.LabelIDs, ",")
			names := strings.Split(raw.LabelNames, ",")
			(*data)[i].LabelList = make([]models.LabelData, len(ids))
			for j := range ids {
				id, _ := strconv.Atoi(ids[j])
				(*data)[i].LabelList[j] = models.LabelData{
					ID:   id,
					Name: names[j],
				}
			}
		}
	}
	return err
}
