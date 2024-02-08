package mysql

import (
	"blog/models"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
)

const (
	secret         = "liuzihao.online"
	tokenIsInvalid = "token存在于黑名单中"
)

func CheckUserExist(username string, email string) (err error) {
	//用户名
	sqlStr := `select count(username) from users where username = ? `
	var count int8
	//db.Get(&count, sqlStr, username) --> 将
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return errors.New("用户已经存在")
	}

	//邮箱
	sqlStr = `select count(user_id) from users where email =?`
	count = 0
	if err := db.Get(&count, sqlStr, email); err != nil {
		return err
	}
	if count > 0 {
		return errors.New("用户已经存在")
	}
	return err
}

func CheckUserExist2(username string, email string) (err error) {
	//用户名
	sqlStr := `select count(username) from users where username = ? `
	var count int8
	//db.Get(&count, sqlStr, username) --> 将
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 1 {
		return errors.New("用户已经存在")
	}

	//邮箱
	sqlStr = `select count(user_id) from users where email =?`
	count = 0
	if err := db.Get(&count, sqlStr, email); err != nil {
		return err
	}
	if count > 1 {
		return errors.New("用户已经存在")
	}
	return err
}

func InsertUser(user *models.User) (err error) {
	sqlStr := `INSERT INTO users (username,password,email,user_id)  values(?,?,?,?)`
	user.Password = encryptPassword(user.Password)
	_, err = db.Exec(sqlStr, user.Username, user.Password, user.Email, user.UserID)
	if err != nil {
		return err
	}
	return err
}

func encryptPassword(oPassword string) string {
	// 创建一个 SHA-256 哈希对象
	h := sha256.New()

	// 将秘密值写入哈希对象
	h.Write([]byte(secret))

	// 将原始密码写入哈希对象，计算哈希值，返回十六进制表示的哈希字符串
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}

func Login(u *models.User) (err error) {
	oldPassword := *(&u.Password)
	sqlStr := `select user_id,username,password from users where username = ?`
	err = db.Get(u, sqlStr, u.Username)
	if errors.Is(err, sql.ErrNoRows) {
		return err
	} else if err != nil {
		return err
	}
	encryptedPassword := encryptPassword(oldPassword)

	if encryptedPassword != u.Password {
		return errors.New("登陆失败")
	}
	return err
}

func Logout(token string, remain int64) (err error) {
	//1.将该token存放在黑名单
	sqlStr := `INSERT INTO tokenInvalid(token, expiration) VALUES(?,?)`
	_, err = db.Exec(sqlStr, token, remain)
	return err
}

func CheckTokenIfInvalid(token string) (err error) {
	sqlStr := `SELECT count(token) FROM tokenInvalid where token = ?`
	var count int8
	err = db.Get(&count, sqlStr, token)
	if count > 0 {
		return errors.New(tokenIsInvalid)
	}
	return nil
}

func UpdateUserMsg(user *models.UserParams, id int64) (err error) {
	err = CheckUserExist2(user.Username, user.Email)
	if err != nil {
		return err
	}
	sqlStr := `UPDATE users SET username = ?,password = ?,email = ? where user_id = ?`
	_, err = db.Exec(sqlStr, user.Username, user.Password, user.Email, id)
	return err
}

func GetUserMsg(user *models.UserParams, id int64) (err error) {
	sqlStr := `SELECT username,email FROM users where user_id = ?`
	err = db.Get(user, sqlStr, id)
	return
}
