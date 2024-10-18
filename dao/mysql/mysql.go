package mysql

import (
	"blog/pkg/snowflake"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var db *sqlx.DB

func Init() (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		viper.GetString("mysql.username"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.dbname"),
	)
	//MustConnect--> 如果没有连接上就panic掉
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(viper.GetInt("mysql.max_open_con"))
	db.SetMaxIdleConns(viper.GetInt("mysql.max_idle_con"))

	//建表操作
	if err = createUserTale(); err != nil {
		zap.L().Error("createUserTale(db) failed,err:%v", zap.Error(err))
		return err
	}
	if err = createKindTable(); err != nil {
		zap.L().Error("createKindTable(db) failed,err:%v", zap.Error(err))
		return err
	}
<<<<<<< HEAD
	err = createClassifyTable(db)
	if err != nil {
		zap.L().Error(" createClassifyTable(db) failed,err:%v", zap.Error(err))
=======
	if err = createLabelTable(); err != nil {
		zap.L().Error(" createLabelTable(db) failed,err:%v", zap.Error(err))
>>>>>>> dev
		return err
	}
	if err = createEssayTable(); err != nil {
		zap.L().Error("createEssayTable(db) failed,err:%v", zap.Error(err))
		return err
	}
	if err = createInvalidToken(); err != nil {
		zap.L().Error("createInvalidToken(db) failed,err:%v", zap.Error(err))
		return err
	}
	if err = createEssayLabelTable(); err != nil {
		zap.L().Error("createEssayLabelTable(db) failed,err:%v", zap.Error(err))
		return err
	}
	return
}

func createUserTale() (err error) {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	//建表操作
	sqlStr1 := `CREATE TABLE IF NOT EXISTS users (
			id Int AUTO_INCREMENT PRIMARY KEY,
			user_id BIGINT NOT NULL ,
			username VARCHAR(24) NOT NULL,
			password VARCHAR(96) NOT NULL,
			email VARCHAR(32) NOT NULL,
			create_time timestamp default CURRENT_TIMESTAMP NULL,
			update_time timestamp default NULL ON UPDATE CURRENT_TIMESTAMP)`
	if _, err = tx.Exec(sqlStr1); err != nil {
		return err
	}

	username := viper.GetString("manager.username")
	password := encryptPassword(viper.GetString("manager.password"))
	email := viper.GetString("manager.email")
	uid := snowflake.GenID()

	//插入管理员
	sqlStr2 := `INSERT INTO users (username,password,email,user_id) SELECT ?,?,?,? WHERE NOT EXISTS(SELECT 1 FROM users WHERE username = ?)`
	_, err = tx.Exec(sqlStr2, username, password, email, uid, username)
	return err
}

func createInvalidToken() (err error) {
	sqlStr := `CREATE TABLE IF NOT EXISTS tokenInvalid(
    id INT AUTO_INCREMENT PRIMARY KEY,
	token text NOT NULL,
	expiration  INT NOT NULL)`
	_, err = db.Exec(sqlStr)
	return err
}

func createKindTable() (err error) {
	sqlStr := `CREATE TABLE IF NOT EXISTS kind(
	id INT AUTO_INCREMENT PRIMARY KEY,
	name VARCHAR(60) NOT NULL,
	icon VARCHAR(60) NOT NULL)`
	_, err = db.Exec(sqlStr)
	return err
}

<<<<<<< HEAD
func createClassifyTable(db *sqlx.DB) (err error) {
	sqlStr := `CREATE TABLE IF NOT EXISTS classify(
	id INT AUTO_INCREMENT PRIMARY KEY,
	kind VARCHAR(60) NOT NULL,
	name VARCHAR(60) NOT NULL,
	router VARCHAR(60) NOT NULL )`
=======
func createLabelTable() (err error) {
	sqlStr := `CREATE TABLE IF NOT EXISTS label(
	id INT AUTO_INCREMENT PRIMARY KEY,
	name VARCHAR(60) NOT NULL
	)`
>>>>>>> dev
	_, err = db.Exec(sqlStr)
	return err
}

func createEssayTable() (err error) {
	sqlStr := `CREATE TABLE IF NOT EXISTS essay(
<<<<<<< HEAD
    id INT AUTO_INCREMENT PRIMARY KEY,
    eid BIGINT NOT NULL ,
    kind VARCHAR(60) NOT NULL,
	name VARCHAR(60) NOT NULL,
	content TEXT NOT NULL,
	introduction VARCHAR(180) NOT NULL,
    router VARCHAR(60) NOT NULL ,
	visitedTimes BIGINT NOT NULL DEFAULT 0 ,
    imgUrl varchar(100) NOT NULL,
    createdTime TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updatedTime TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	advertiseMsg VARCHAR(30),
	advertiseImg VARCHAR(100),
	advertiseHref VARCHAR(100)
=======
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(60) NOT NULL,
		eid BIGINT NOT NULL ,
		kind_id INT NOT NULL,
		content TEXT NOT NULL,
		introduction VARCHAR(180) NOT NULL,
		visitedTimes BIGINT NOT NULL DEFAULT 0 ,
		imgUrl varchar(100) NOT NULL,
		createdTime TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updatedTime TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		ifRecommend BOOL NOT NULL DEFAULT FALSE,
		FOREIGN KEY (kind_id) REFERENCES kind(id)
            ON DELETE RESTRICT
            ON UPDATE CASCADE
>>>>>>> dev
    )`
	_, err = db.Exec(sqlStr)
	return err
}

func createEssayLabelTable() (err error) {
	sqlStr := `CREATE TABLE IF NOT EXISTS essay_label(
    essay_id INT NOT NULL PRIMARY KEY,
    label_id INT NOT NULL PRIMARY KEY,
    FOREIGN KEY (essay_id) REFERENCES essay(id)
        	ON DELETE RESTRICT
            ON UPDATE CASCADE,
    FOREIGN KEY (label_id) REFERENCES label(id)  
        	ON DELETE RESTRICT
            ON UPDATE CASCADE
    )`
	_, err = db.Exec(sqlStr)
	return err
}
