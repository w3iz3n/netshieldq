package controller

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"pnetshieldq/dao"
	"pnetshieldq/entity"
	"pnetshieldq/utils"
	"strconv"
	"time"
)

type FileController struct {
	FileDao   *dao.FileDao
	OSSClient utils.OSSClient
}

func NewFileController(fileDao *dao.FileDao, ossClient utils.OSSClient) *FileController {
	return &FileController{
		FileDao:   fileDao,
		OSSClient: ossClient,
	}
}

// @Summary Upload a file
// @Description Upload a file to the server
// @Accept  mpfd
// @Produce  json
// @Param file formData file true "The file to upload"
// @Success 200 {string} string "File uploaded successfully"
// @Router /file/upload [post]
func (fc *FileController) UploadFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(60 << 20)
	if err != nil {
		http.Error(w, "Error parsing form: "+err.Error(), http.StatusInternalServerError)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	receiverID := r.FormValue("receiverId")
	if receiverID == "" {
		http.Error(w, "Receiver ID is required", http.StatusBadRequest)
		return
	}

	senderID := utils.GetUserIDFromJWT(r)
	if err != nil {
		http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
		return
	}

	tempFile, err := ioutil.TempFile("temp", "upload-*")
	if err != nil {
		http.Error(w, "Failed to create temporary file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	_, err = io.Copy(tempFile, file)
	if err != nil {
		http.Error(w, "Failed to save file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	folderName := "netshieldq"
	filePath := path.Join(folderName, handler.Filename)

	fileURL, err := fc.OSSClient.UploadFile(filePath, tempFile.Name())
	if err != nil {
		http.Error(w, "Failed to upload file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fileRecord := &entity.File{
		FileName:   handler.Filename,
		FileURL:    fileURL,
		SenderID:   strconv.Itoa(senderID),
		ReceiverID: receiverID,
		Timestamp:  time.Now(),
		FileType:   "file",
	}

	if err := fc.FileDao.Create(fileRecord); err != nil {
		http.Error(w, "Failed to save file record: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("File uploaded successfully. URL: " + fileURL))
}

// @Summary Delete a file
// @Description Delete a file from the server
// @Param fileName query string true "The name of the file to delete"
// @Success 200 {string} string "File deleted successfully"
// @Router /file/delete [get]
func (fc *FileController) DeleteFile(w http.ResponseWriter, r *http.Request) {
	// 使用文件ID作为参数
	fileName := r.URL.Query().Get("fileName")

	// 查找文件以确保它存在
	file, err := fc.FileDao.FindByName(fileName)
	if err != nil {
		http.Error(w, "File not found: "+err.Error(), http.StatusNotFound)
		return
	}

	// 构建文件路径
	folderName := "netshieldq"
	filePath := path.Join(folderName, file.FileName)

	// 删除对象存储中的文件
	err = fc.OSSClient.DeleteFile(filePath)
	if err != nil {
		http.Error(w, "Failed to delete file in storage: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 从数据库中删除文件记录
	err = fc.FileDao.Delete(file)
	if err != nil {
		http.Error(w, "Failed to delete file record: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("File deleted successfully"))
}

// @Summary Get user's files
// @Description Get all files related to the current user
// @Accept  json
// @Produce  json
// @Success 200 {array} entity.File
// @Router /file/userfiles [get]
func (fc *FileController) GetUserFiles(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserIDFromJWT(r)
	if userID == -1 {
		http.Error(w, "Invalid token or no token provided", http.StatusUnauthorized)
		return
	}
	files, err := fc.FileDao.GetByUserID(strconv.Itoa(userID))
	if err != nil {
		http.Error(w, "Failed to retrieve files", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(files); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

// @Summary Download a file
// @Description Download a file by ID
// @Param fileID query string true "The ID of the file to download"
// @Produce octet-stream
// @Success 200 "File downloaded successfully"
// @Failure 404 "File not found"
// @Failure 500 "Internal server error"
// @Router /file/download [get]
func (fc *FileController) DownloadFile(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("fileName")
	if filename == "" {
		http.Error(w, "File ID is required", http.StatusBadRequest)
		return
	}

	// 查找文件以确保它存在
	file, err := fc.FileDao.FindByName(filename)
	if err != nil {
		http.Error(w, "File not found: "+err.Error(), http.StatusNotFound)
		return
	}

	// 从对象存储服务中获取文件
	filePath := path.Join("netshieldq", file.FileName)
	fileContent, err := fc.OSSClient.GetFile(filePath)
	if err != nil {
		http.Error(w, "Failed to retrieve file from storage: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 设置响应头以允许文件下载
	w.Header().Set("Content-Disposition", "attachment; filename="+file.FileName)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)

	// 发送文件内容
	if _, err := w.Write(fileContent); err != nil {
		http.Error(w, "Failed to send file content: "+err.Error(), http.StatusInternalServerError)
	}
}
