package utils

import "time"

func GetChineseTime() (string, error) {
	//加载中国时区
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return "", err
	}
	T := time.Now().In(loc)
	t := T.Format("2006-01-02 15:04:05")
	return t, nil
}
