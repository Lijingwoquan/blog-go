package logic

import (
	"blog/dao/mysql"
	"blog/models"
	"blog/pkg/jwt"
	"blog/pkg/snowflake"
	"go.uber.org/zap"
)

func Signup(u *models.UserParams) (err error) {
	//1.判断用户是否存在 --> 判断username和email
	err = mysql.CheckUserExist(u.Username, u.Email)
	if err != nil {
		zap.L().Info("mysql.CheckUserExist(u.Username) failed", zap.Error(err))
		return err
	}

	//2.生成uid
	uid := snowflake.GenID()

	//3.将用户储存在数据库
	user := new(models.User)
	user.UserID = uid
	user.Username = u.Username
	user.Password = u.Password
	user.Email = u.Email
	err = mysql.InsertUser(user)
	if err != nil {
		zap.L().Error("mysql.InsertUser() failed", zap.Error(err))
		return err
	}
	return err
}

func Login(u *models.User) (err error) {
	//1.判断账号密码是否正确
	if err = mysql.Login(u); err != nil {
		zap.L().Info("mysql.Login() failed", zap.Error(err))
		return err
	}

	//2. jwt生成token
	var token string
	token, err = jwt.GenToken(u)
	if err != nil {
		zap.L().Error("jwt.GenToken(u) failed", zap.Error(err))
		return err
	}

	//将token保存
	u.Token = token
	return err
}
