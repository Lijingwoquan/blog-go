package models

type User struct {
	UserID   int64  `db:"user_id" json:"userID string"`
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"password"`
	Email    string `db:"email" json:"email"`
	Token    string
}
