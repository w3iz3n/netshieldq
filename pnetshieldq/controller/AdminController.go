package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pnetshieldq/dao"
	"pnetshieldq/utils"
	"strconv"
)

type AdminController struct {
	AdminDao *dao.AdminDao
	UserDao  *dao.UserDao
}

func NewAdminController(adminDao *dao.AdminDao, userDao *dao.UserDao) *AdminController {
	return &AdminController{
		AdminDao: adminDao,
		UserDao:  userDao,
	}
}
func (mc *AdminController) AdminLogin(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	fmt.Println("AdminLogin")
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//log := log.GetLogger()
	//log.Info("AdminLogin")
	isAdmin, err := mc.AdminDao.CheckAdmin(requestData.Username, requestData.Password)
	if err != nil {
		http.Error(w, "Database error while checking admin", http.StatusInternalServerError)
		return
	}
	if !isAdmin {
		http.Error(w, "User is not an admin", http.StatusForbidden)
		return
	}

	token, err := utils.GenerateJWT(1, requestData.Username, requestData.Password)
	if err != nil {
		http.Error(w, "Failed to generate JWT", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
func (ac *AdminController) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := ac.UserDao.FindAll()
	if err != nil {
		http.Error(w, "Error while getting users", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

func (ac *AdminController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userId, _ := strconv.Atoi(r.URL.Query().Get("id"))
	if err := ac.UserDao.DeleteUser(userId); err != nil {
		http.Error(w, "Error while deleting user", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("User deleted successfully"))
}
