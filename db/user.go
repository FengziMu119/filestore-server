package db

import (
	"database/sql"
	mydb "filestore-server/db/mysql"
	"fmt"
)

//UserSignup : 通过用户名密码注册
func UserSignup(username string, password string) bool {
	var (
		stmt         *sql.Stmt
		err          error
		sqlStr       string
		ret          sql.Result
		rowsAffected int64
	)
	sqlStr = "insert into tbl_user (`user_name`,`user_pwd`)VALUES(?,?)"
	if stmt, err = mydb.DBconn().Prepare(sqlStr); err != nil {
		fmt.Printf("Failed to insert err,err:%s\n", err.Error())
		return false
	}
	if ret, err = stmt.Exec(username, password); err != nil {
		fmt.Printf("Failed to insert err,err:%s\n", err.Error())
		return false
	}
	if rowsAffected, err = ret.RowsAffected(); err == nil && rowsAffected > 0 {
		return true
	}
	return false
}

//UserSignin : 判断密码是否一致
func UserSignin(username string, encpwd string) bool {
	var (
		sqlStr string
		stmt   *sql.Stmt
		rows   *sql.Rows
		pRows  []map[string]interface{}
		err    error
	)
	sqlStr = "select * from tbl_user where user_name = ? limit 1"
	if stmt, err = mydb.DBconn().Prepare(sqlStr); err != nil {
		fmt.Printf("stmt err %s\n", err.Error())
		return false
	}
	defer stmt.Close()
	rows, err = stmt.Query(username)
	if err != nil {
		fmt.Println(err.Error())
		return false
	} else if rows == nil {
		fmt.Printf("Username unfind err %s\n", err.Error())
		return false
	}
	pRows = mydb.ParseRows(rows)
	if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == encpwd {
		return true
	}
	return false
}

// UpdateToken : 刷新用户登录的token
func UpdateToken(username string, token string) bool {
	var (
		sqlStr string
		stmt   *sql.Stmt
		//ret sql.Result
		err error
	)
	sqlStr = "replace into tbl_user_token(`user_name`,`user_token`)VALUES(?,?)"
	if stmt, err = mydb.DBconn().Prepare(sqlStr); err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()
	if _, err = stmt.Exec(username, token); err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

type UserInfo struct {
	Username     string
	Email        string
	Phone        string
	SignupAt     string
	LastActiveAt string
	Status       int
}

//GetUserInfo 用户信息查询
func GetUserInfo(username string) (UserInfo, error) {
	var (
		user   UserInfo
		sqlStr string
		stmt   *sql.Stmt
		err    error
	)
	sqlStr = "select user_name, signup_at from tbl_user where user_name = ? limit 1"
	if stmt, err = mydb.DBconn().Prepare(sqlStr); err != nil {
		fmt.Println(err.Error())
		return user, err
	}
	defer stmt.Close()
	if err = stmt.QueryRow(username).Scan(&user.Username, &user.SignupAt); err != nil {
		return user, err
	}
	return user, nil
}
