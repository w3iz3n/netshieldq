package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"pnetshieldq/dao"
	"pnetshieldq/entity"
	"pnetshieldq/utils"
	"strconv"
	"time"
)

type UserController struct {
	userDao       *dao.UserDao
	userAvatarDao *dao.UserAvatarDAO
}
type UserPublic struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	LastSeen string `json:"last_seen"`
}

func NewUserController(userDao *dao.UserDao, userAvatarDao *dao.UserAvatarDAO) *UserController {
	return &UserController{
		userDao:       userDao,
		userAvatarDao: userAvatarDao,
	}
}

func (c *UserController) Register(w http.ResponseWriter, r *http.Request) {
	var registrationData struct {
		Username string
		Password string
		Email    string
	}
	if err := json.NewDecoder(r.Body).Decode(&registrationData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registrationData.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// 创建通用用户记录
	user := entity.User{
		Username: registrationData.Username,
		Password: string(hashedPassword),
		Email:    registrationData.Email,
	}

	if err := c.userDao.Create(&user); err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}
	defaultAvatarBase64, err := utils.GetDefaultAvatarBase64()
	if err != nil {
		http.Error(w, "Failed to get default avatar", http.StatusInternalServerError)
		return
	}

	if _, err := c.userAvatarDao.InsertUserAvatar(defaultAvatarBase64); err != nil {
		fmt.Println("Failed to insert avatar:", err)
		http.Error(w, "Failed to save avatar", http.StatusInternalServerError)
		return
	}
	//// 创建 MQTT 用户记录
	//mqttUser := entity.MqttUser{
	//	Username: registrationData.Username,
	//	Password: string(hashedPassword), // 注意：这里应根据实际需求决定是否再次加密或使用不同的密码
	//}
	//
	//if err := c.mqttUserDao.Create(&mqttUser); err != nil {
	//	http.Error(w, "Failed to create MQTT user", http.StatusInternalServerError)
	//	return
	//}
	//
	//// 创建 MQTT ACL 条目
	//acl := entity.MQTTACL{
	//	Allow:    true,
	//	Username: &user.Username,
	//	ClientID: "user_" + user.Username,        // 如果需要，设置正确的 ClientID
	//	Access:   3,                              // 订阅和发布权限
	//	Topic:    "user/" + user.Username + "/#", // 为用户分配专属 topic
	//}
	//
	//if err := c.aclDao.Insert(&acl); err != nil {
	//	http.Error(w, "Failed to set MQTT ACL", http.StatusInternalServerError)
	//	return
	//}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
	fmt.Println("注册成功")
}

func (c *UserController) Login(w http.ResponseWriter, r *http.Request) {
	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	storedUser, err := c.userDao.FindByUsername(loginData.Username)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	error := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(loginData.Password))
	if error != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}
	tokenString, err := utils.GenerateJWT(storedUser.UserID, storedUser.Username, loginData.Password)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	if err != nil {
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	// 发送包含 token, userID 和 username 的 JSON 响应
	responseData := map[string]interface{}{
		"token":    tokenString,
		"userid":   storedUser.UserID,
		"username": storedUser.Username,
	}
	json.NewEncoder(w).Encode(responseData)
}

func (c *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	userId, _ := strconv.Atoi(r.URL.Query().Get("id"))
	user, err := c.userDao.FindByID(userId)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}

// @Summary Update User
// @Description update user by id
// @Param   id path int true "User ID"
// @Param   user body entity.User true "User"
// @Success 200 {string} string	"User updated successfully"
// @Router /users/{id} [put]
func (c *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var updateData struct {
		Email    *string `json:"email"`
		Password *string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 获取当前用户信息
	_, err = c.userDao.FindByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// 创建一个更新字段的 map
	updates := make(map[string]interface{})

	// 检查并更新邮箱
	if updateData.Email != nil {
		updates["email"] = *updateData.Email
	}

	// 检查并更新密码
	if updateData.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*updateData.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		updates["password"] = string(hashedPassword)
	}

	// 更新 LastSeen 为当前时间
	updates["last_seen"] = time.Now()

	// 调用更新方法
	if err := c.userDao.UpdateFields(userID, updates); err != nil {
		http.Error(w, "Error while updating user", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("User updated successfully"))
}

// @Summary Delete User
// @Description delete user by id
// @Param   id path int true "User ID"
// @Success 200 {string} string	"User deleted successfully"
// @Router /users/{id} [delete]
func (c *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := c.userDao.DeleteByID(userId); err != nil {
		http.Error(w, "Error while deleting user", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("User deleted successfully"))
}

// GetAllUsers retrieves all users without password
func (c *UserController) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := c.userDao.FindAll()
	if err != nil {
		http.Error(w, "Error while getting users", http.StatusInternalServerError)
		return
	}

	var usersPublic []UserPublic
	for _, user := range users {
		usersPublic = append(usersPublic, UserPublic{
			UserID:   user.UserID,
			Username: user.Username,
			LastSeen: user.LastSeen.Format("2006-01-02 15:04:05"),
		})
	}

	json.NewEncoder(w).Encode(usersPublic)
}
func (uc *UserController) GetUserCount(w http.ResponseWriter, r *http.Request) {
	count, err := uc.userDao.GetMaxUserID()
	if err != nil {
		http.Error(w, "Failed to get user count: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		UserCount int `json:"user_count"`
	}{
		UserCount: count,
	}

	// Set the response header and encode the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
	}
}
