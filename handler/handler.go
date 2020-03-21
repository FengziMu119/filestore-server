package handler

import (
	"encoding/json"
	"filestore-server/meta"
	"filestore-server/util"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

const pwd_salt = "sftfdsa"

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	var (
		data     []byte
		err      error
		file     multipart.File
		head     *multipart.FileHeader
		path     string
		fileMeta meta.FileMate
		newFile  *os.File
	)
	if r.Method == "GET" {
		// 返回上传html页面
		if data, err = ioutil.ReadFile("./static/view/index.html"); err != nil {
			io.WriteString(w, "Internal server err")
			return
		}
		io.WriteString(w, string(data))

	} else if r.Method == "POST" {
		// 接收文件流及保存到本地目录
		if file, head, err = r.FormFile("file"); err != nil {
			fmt.Printf("Failed to get data err %s\n", err.Error())
			return
		}
		defer file.Close()
		path = "G:/workspace/src/github.com/FengziMu119/filestore-server/tmp/" + head.Filename
		fileMeta = meta.FileMate{
			FileName: head.Filename,
			Location: path,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		if newFile, err = os.Create(fileMeta.Location); err != nil {
			fmt.Printf("Failed create file,err %s\n", err.Error())
			return
		}
		defer newFile.Close()
		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("Failed to save file ,err %s\n", err.Error())
			return
		}
		newFile.Seek(0, 0)
		fileMeta.FileSha1 = util.FileSha1(newFile)
		fmt.Println(fileMeta.FileSha1)
		//meta.UploadFileMate(fileMeta)
		_ = meta.UploadFileMetaDB(fileMeta)
		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}
}

//UploadSucHandler: 上传完成
func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload file success!")
}

//GetFileMetaHandler: 获取文件元信息
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	var (
		filehash string
		fMeta    meta.FileMate
		data     []byte
		err      error
	)
	r.ParseForm()
	filehash = r.Form["filehash"][0]
	//fMeta = meta.GetFileMeta(filehash)
	fMeta = meta.GetFileMateDB(filehash)
	if data, err = json.Marshal(fMeta); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

//DownloadHandler: 下载文件
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	var (
		fsha1 string
		fm    meta.FileMate
		err   error
		f     *os.File
		data  []byte
	)
	r.ParseForm()
	fsha1 = r.Form.Get("filehash")
	fm = meta.GetFileMeta(fsha1)
	if f, err = os.Open(fm.Location); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	if data, err = ioutil.ReadAll(f); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/octect-stream")
	w.Header().Set("Content-Descrption", "attachment;filename=\""+fm.FileName+"\"")
	w.Write(data)
}

//FileMataUpdateHandler: 重命名
func FileMataUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var (
		opType      string
		fileSha1    string
		newFileName string
		curFileMata meta.FileMate
		data        []byte
		err         error
	)
	r.ParseForm()
	opType = r.Form.Get("op")
	fileSha1 = r.Form.Get("filehash")
	newFileName = r.Form.Get("filename")
	if opType != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	curFileMata = meta.GetFileMeta(fileSha1)
	curFileMata.FileName = newFileName
	meta.UploadFileMate(curFileMata)
	if data, err = json.Marshal(curFileMata); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
	w.WriteHeader(http.StatusOK)
}

//FileDeleteHandler: 删除元信息
func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	var (
		fileSha1 string
		fMeta    meta.FileMate
	)
	r.ParseForm()
	fileSha1 = r.Form.Get("filesha1")
	fMeta = meta.GetFileMeta(fileSha1)
	os.Remove(fMeta.Location)
	meta.RemoveFileMeta(fileSha1)
	w.WriteHeader(http.StatusOK)
}
