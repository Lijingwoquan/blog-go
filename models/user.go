package models

// User 储存在数据量的信息
type User struct {
	UserID   int64  `db:"user_id" json:"userID string"`
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"password"`
	Email    string `db:"email" json:"email"`
	Token    string `json:"token"`
}

// UserInfo 返回给前端的信息
type UserInfo struct {
	UserID   int64  `db:"user_id" json:"userID,string"`
	Username string `db:"username" json:"username"`
	Email    string `db:"email" json:"email"`
}

type UserParams struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"rePassword" binding:"required,eqfield=Password"`
	Email      string `json:"email" binding:"required"`
}
