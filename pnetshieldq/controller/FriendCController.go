package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pnetshieldq/log"
	"strconv"

	"pnetshieldq/dao"
	"pnetshieldq/entity"
)

type FriendCController struct {
	friendcDao *dao.FriendCDAO
}

func NewFriendCController(friendcDao *dao.FriendCDAO) *FriendCController {
	return &FriendCController{
		friendcDao: friendcDao,
	}
}

type FriendCRequest struct {
	ID       string `json:"UserID"`
	FriendID uint   `json:"FriendID"`
	C        string `json:"C"`
}

func (c *FriendCController) CreateFriendC(w http.ResponseWriter, r *http.Request) {
	var friendC entity.FriendC
	var friendCRequest FriendCRequest
	if err := json.NewDecoder(r.Body).Decode(&friendCRequest); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseUint(friendCRequest.ID, 10, 32)
	if err != nil {
		fmt.Println("解析问题")
	}
	friendC.UserID = uint(id)
	friendC.FriendID = friendCRequest.FriendID
	friendC.C = friendCRequest.C
	log := log.GetLogger()
	log.Info("User_" + friendCRequest.ID + " Generate C_value")
	if err := c.friendcDao.Create(&friendC); err != nil {
		http.Error(w, "Failed to create FriendC record: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(friendC)
}

func (c *FriendCController) GetFriendC(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	log := log.GetLogger()
	log.Info("User_" + idStr + " Get C_value")
	friendC, err := c.friendcDao.GetByID(uint(id))
	if err != nil {
		http.Error(w, "FriendC not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(friendC)
}
func (c *FriendCController) GetFriendCByFriendIDAndUserID(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("userID")
	friendIDStr := r.URL.Query().Get("friendID")

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid userID", http.StatusBadRequest)
		return
	}
	log := log.GetLogger()
	log.Info("User_" + friendIDStr + " Get" + userIDStr + " C_value")
	friendID, err := strconv.ParseUint(friendIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid friendID", http.StatusBadRequest)
		return
	}

	friendC, err := c.friendcDao.GetByFriendIDAndUserID(uint(userID), uint(friendID))
	if err != nil {
		http.Error(w, "FriendC not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(friendC.C)
}
