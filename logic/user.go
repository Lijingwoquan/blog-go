package logic

import (
	"blog/dao/mysql"
	"blog/models"
	"blog/pkg/jwt"
	"blog/pkg/snowflake"
)

func Signup(u *models.UserParams) (err error) {
	//1.判断用户是否存在 --> 判断username和email
	if err = mysql.CheckUserExist(u.Username, u.Email); err != nil {
		return err
	}
	//2.生成uid并保存相关信息
	uid := snowflake.GenID()
	user := models.User{
		UserID:   uid,
		Username: u.Username,
		Password: u.Password,
		Email:    u.Email,
		Token:    "",
	}
	//3.将用户储存在数据库
	return mysql.InsertUser(&user)
}

func Login(u *models.User) (err error) {
	//1.判断账号密码是否正确
	if err = mysql.Login(u); err != nil {
		return err
	}
	//2. jwt生成token
	var token string
	if token, err = jwt.GenToken(u); err != nil {
		return err
	}
	//将token保存
	u.Token = token
	return nil
}

//func Logout(token string) error {
//	//1.得到token还剩余的时间
//	MyClaims, err := jwt.ParseToken(token)
//	if err != nil {
//		return err
//	}
//	//2.将该token储存在数据库中
//	return mysql.Logout(token, MyClaims.ExpiresAt)
//}

func UpdateUserMsg(user *models.UserParams, id int64) error {
	//从数据库中修改数据
	return mysql.UpdateUserMsg(user, id)
}
