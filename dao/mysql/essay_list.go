package mysql

import (
	"blog/models"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

func GetRecommendEssayList(data *[]models.EssayData) error {
	sqlStr := `
		SELECT e.id, e.name, e.created_time, e.img_url 
		FROM essay e 
		WHERE if_recommend = true
		ORDER BY e.id DESC 
		LIMIT 5
	`
	return db.Select(data, sqlStr)
}

func GetEssayList(data *models.EssayListAndPage, query models.EssayQuery) error {
	var wg sync.WaitGroup
	var errChan = make(chan error, 2)
	// 计算偏移量
	offset := (query.Page - 1) * query.PageSize
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

	args = append(args, query.PageSize, offset)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := getList(data, whereClause, args); err != nil {
			errChan <- fmt.Errorf("getList failed,err:%w", err)
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := getCount(data, query.PageSize, whereClause, args); err != nil {
			errChan <- fmt.Errorf("getList failed,err:%w", err)
			return
		}
	}()

	wg.Wait()

	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func getList(data *models.EssayListAndPage, whereClause string, args []interface{}) error {
	rawDataList := make([]rawData, 0, 5)
	baseSelect := `
        SELECT e.id, e.name, e.kind_id, e.if_recommend, e.if_top, e.introduction, e.created_time, e.visited_times, e.img_url,
            k.name AS kind_name,
            GROUP_CONCAT(el.label_id) AS label_ids,
            GROUP_CONCAT(el.label_name) AS label_names
        FROM essay e
        LEFT JOIN kind k ON e.kind_id = k.id
        LEFT JOIN essay_label el ON e.id = el.essay_id
        `

	groupBy := "GROUP BY e.id"

	// 优先排序ifTop为true的记录
	orderBy := "ORDER BY e.if_top DESC, e.id DESC"

	sqlStr := fmt.Sprintf("%s %s %s %s LIMIT ? OFFSET ?",
		baseSelect, whereClause, groupBy, orderBy)

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
	return nil
}

func getCount(data *models.EssayListAndPage, PageSize int, whereClause string, args []interface{}) error {
	baseSql := `
        SELECT COUNT(DISTINCT e.id)
        FROM essay e 
        LEFT JOIN kind k ON e.kind_id = k.id
        LEFT JOIN essay_label el ON e.id = el.essay_id
   	`

	var totalCount int
	sqlStr := fmt.Sprintf("%s %s", baseSql, whereClause)
	if err := db.Get(&totalCount, sqlStr, args[:len(args)-2]...); err != nil {
		return err
	}

	data.TotalPage = (totalCount + PageSize - 1) / PageSize
	return nil
}

func GetAllEssay(data *[]models.EssayData) error {
	sqlStr := `
		SELECT e.id, e.name, e.created_time, e.img_url,e.kind_id,
		       k.name AS kind_name,
		       GROUP_CONCAT(el.label_id) AS label_ids ,GROUP_CONCAT(el.label_name) AS label_names
		FROM essay e
		LEFT JOIN kind k on e.kind_id = k.id
		LEFT JOIN essay_label el on e.id = el.essay_id
		GROUP BY e.id
		ORDER BY e.id DESC
	`
	var err error
	var rawDataList = new([]rawData)
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
