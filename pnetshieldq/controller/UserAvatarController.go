package controller

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"pnetshieldq/dao"
	"pnetshieldq/utils"
	"strconv"
	"strings"
)

type UserAvatarController struct {
	AvatarDAO *dao.UserAvatarDAO
}

func NewUserAvatarController(avatarDAO *dao.UserAvatarDAO) *UserAvatarController {
	return &UserAvatarController{AvatarDAO: avatarDAO}
}

type AvatarResponse struct {
	Avatar string `json:"avatar"`
}

func (c *UserAvatarController) UpdateAvatarHandler(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserIDFromJWT(r)
	//fmt.Println("头像更新成功")
	// 解析 multipart/form-data 请求
	err := r.ParseMultipartForm(10 << 20) // 限制上传文件的大小
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// 获取文件
	file, _, err := r.FormFile("avatar")
	if err != nil {
		http.Error(w, "Failed to retrieve file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 读取文件内容
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	// 确定文件的 MIME 类型
	contentType := http.DetectContentType(fileBytes)
	if !strings.HasPrefix(contentType, "image/") {
		http.Error(w, "Invalid image type", http.StatusBadRequest)
		return
	}

	// 将文件内容编码为 Base64
	base64Image := base64.StdEncoding.EncodeToString(fileBytes)
	//base64Image = "data:" + contentType + ";base64," + base64Image
	//fmt.Println(base64Image)

	// 更新 Base64 编码的图像到数据库
	err = c.AvatarDAO.UpdateUserAvatar(userID, base64Image)
	if err != nil {
		http.Error(w, "Failed to update avatar", http.StatusInternalServerError)
		return
	}

	// 返回成功响应
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Avatar updated successfully"))
}

func (c *UserAvatarController) GetAvatarHandler(w http.ResponseWriter, r *http.Request) {
	// 获取用户ID
	userID := utils.GetUserIDFromJWT(r)
	// 从数据库获取头像数据
	avatarData, err := c.AvatarDAO.GetUserAvatar(userID)
	//fmt.Println(avatarData)
	if err != nil {
		if err.Error() == "user avatar not found" {
			http.Error(w, "User avatar not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to get avatar", http.StatusInternalServerError)
		}
		return
	}
	base64Avatar := string(avatarData)

	// 返回 JSON 格式的数据给前端
	response := AvatarResponse{
		Avatar: base64Avatar,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (c *UserAvatarController) GetAvatarHandlerFriend(w http.ResponseWriter, r *http.Request) {
	// 从URL路径中获取用户ID
	vars := mux.Vars(r)

	userIDStr := vars["userID"]

	// 将用户ID从字符串转换为整数
	userID, err := strconv.Atoi(userIDStr)
	fmt.Println(userIDStr)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	// 从数据库获取头像数据
	avatarData, err := c.AvatarDAO.GetUserAvatar(userID)
	if err != nil {
		if err.Error() == "user avatar not found" {
			http.Error(w, "User avatar not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to get avatar", http.StatusInternalServerError)
		}
		return
	}
	base64Avatar := string(avatarData)

	// 返回 JSON 格式的数据给前端
	response := AvatarResponse{
		Avatar: base64Avatar,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
