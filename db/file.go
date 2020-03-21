package db

import (
	"database/sql"
	mydb "filestore-server/db/mysql"
	"fmt"
)

// OnFileUploadFinished:文件上传完成，保存到数据库
func OnFileUploadFinished(filehash string, filename string, filesize int64, fileaddr string) bool {
	var (
		sqlStr string
		stmt   *sql.Stmt
		err    error
	)
	sqlStr = "insert into tbl_file (file_sha1,file_name,file_size,file_addr,status) VALUES (?,?,?,?,1)"
	if stmt, err = mydb.DBconn().Prepare(sqlStr); err != nil {
		fmt.Println("Failed to prepare statement, err" + err.Error())
		return false
	}
	defer stmt.Close()
	if _, err = stmt.Exec(filehash, filename, filesize, fileaddr); err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

//GetFileMeta : 获取上传的文件
func GetFileMeta(filesha1 string) (*TableFile, error) {
	var (
		sqlStr string
		stmt   *sql.Stmt
		err    error
		file   TableFile
	)
	sqlStr = "select file_sha1, file_name, file_size from tbl_file where file_hash = ? and status = 1"
	if stmt, err = mydb.DBconn().Prepare(sqlStr); err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer stmt.Close()

	if err = stmt.QueryRow(filesha1).Scan(&file.FileHash, &file.FileName, &file.FileSize, &file.FileAddr); err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return &file, nil
}
