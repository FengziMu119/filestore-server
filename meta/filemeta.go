package meta

import (
	mydb "filestore-server/db"
)

// FileMate : 文件元信息
type FileMate struct {
	FileSha1 string // sha1加密
	FileName string // 文件名
	FileSize int64  // 文件大小
	Location string // 文件路径
	UploadAt string // 上传时间
}

var fileMetas map[string]FileMate

func init() {
	fileMetas = make(map[string]FileMate)
}

//UploadFileMate: 新增/更新文件元信息
func UploadFileMate(fmeta FileMate) {
	fileMetas[fmeta.FileSha1] = fmeta
}

//UploadFileMetaDB: 保存数据到mysql
func UploadFileMetaDB(fmeta FileMate) bool {
	return mydb.OnFileUploadFinished(fmeta.FileSha1, fmeta.FileName, fmeta.FileSize, fmeta.Location)
}

//GetFileMateDB: 获取数据信息
func GetFileMateDB(filesha1 string) FileMate {
	var (
		fmeta FileMate
		tfile *mydb.TableFile
		err   error
	)
	if tfile, err = mydb.GetFileMeta(filesha1); err != nil {
		return FileMate{}
	}
	fmeta = FileMate{
		FileSha1: tfile.FileHash,
		FileName: tfile.FileName.String,
		FileSize: tfile.FileSize.Int64,
		Location: tfile.FileAddr.String,
	}
	return fmeta
}

// GetFileMeta : 通过sha1获取元信息
func GetFileMeta(fileSha1 string) FileMate {
	return fileMetas[fileSha1]
}

// RemoveFileMeta : 删除元信息
func RemoveFileMeta(fileSha1 string) {
	delete(fileMetas, fileSha1)
}
